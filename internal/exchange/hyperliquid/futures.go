package hyperliquid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"orderbook/internal/exchange"

	"github.com/gorilla/websocket"
)

const (
	hyperliquidWSURL = "wss://api.hyperliquid.xyz/ws"
	infoRESTEndpoint = "https://api.hyperliquid.xyz/info"
	pingInterval     = 50 * time.Second
)

type FuturesExchange struct {
	symbol     string
	wsConn     *websocket.Conn
	updateChan chan *exchange.DepthUpdate
	done       chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc

	health      atomic.Value
	snapshotMu  sync.Mutex
	snapshot    *exchange.Snapshot
	snapshotSet bool
}

func NewFuturesExchange(config Config) *FuturesExchange {
	ctx, cancel := context.WithCancel(context.Background())

	ex := &FuturesExchange{
		symbol:     config.Symbol,
		updateChan: make(chan *exchange.DepthUpdate, 1000),
		done:       make(chan struct{}),
		ctx:        ctx,
		cancel:     cancel,
	}

	ex.health.Store(exchange.HealthStatus{
		Connected:    false,
		LastPing:     time.Time{},
		MessageCount: 0,
		ErrorCount:   0,
	})

	return ex
}

func (e *FuturesExchange) GetName() exchange.ExchangeName { return exchange.Hyperliquidf }
func (e *FuturesExchange) GetSymbol() string              { return e.symbol }

func (e *FuturesExchange) Connect(ctx context.Context) error {
	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}

	conn, _, err := dialer.DialContext(ctx, hyperliquidWSURL, nil)
	if err != nil {
		e.incrementErrorCount()
		return fmt.Errorf("hyperliquid websocket dial failed: %w", err)
	}

	e.wsConn = conn
	e.updateConnectionStatus(true)
	log.Printf("[%s] WebSocket connected to %s", e.GetName(), hyperliquidWSURL)

	// helper to write subscribe
	subscribe := func(sub interface{}) error {
		msg := map[string]interface{}{"method": "subscribe", "subscription": sub}
		return conn.WriteJSON(msg)
	}

	coin := mapToHLCoin(e.symbol)

	// subscribe l2Book (orderbook)
	if err := subscribe(map[string]interface{}{"type": "l2Book", "coin": coin}); err != nil {
		e.incrementErrorCount()
		conn.Close()
		return fmt.Errorf("failed subscribe l2Book: %w", err)
	}

	// subscribe trades (optional)
	if err := subscribe(map[string]interface{}{"type": "trades", "coin": coin}); err != nil {
		// do not fail entire connector for trades subscription; log only
		log.Printf("[%s] warning: failed subscribe trades: %v", e.GetName(), err)
	}

	go e.pingLoop()
	go e.readMessages()

	return nil
}

func (e *FuturesExchange) Close() error {
	if e.cancel != nil {
		e.cancel()
	}

	if e.wsConn != nil {
		select {
		case <-e.done:
		default:
			close(e.done)
		}

		_ = e.wsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

		select {
		case <-time.After(time.Second):
		}

		e.updateConnectionStatus(false)
		return e.wsConn.Close()
	}
	return nil
}

func (e *FuturesExchange) GetSnapshot(ctx context.Context) (*exchange.Snapshot, error) {
	log.Printf("[%s] Fetching L2 snapshot via REST for %s", e.GetName(), e.symbol)

	reqBody := map[string]interface{}{"type": "l2Book", "coin": mapToHLCoin(e.symbol)}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", infoRESTEndpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		e.incrementErrorCount()
		return nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Coin   string      `json:"coin"`
		Time   int64       `json:"time"`
		Levels [][]WsLevel `json:"levels"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		e.incrementErrorCount()
		return nil, fmt.Errorf("decode snapshot failed: %w", err)
	}

	bids := make([]exchange.PriceLevel, 0, len(res.Levels[0]))
	for _, lvl := range res.Levels[0] {
		bids = append(bids, exchange.PriceLevel{Price: lvl.Px, Quantity: lvl.Sz})
	}
	asks := make([]exchange.PriceLevel, 0, len(res.Levels[1]))
	for _, lvl := range res.Levels[1] {
		asks = append(asks, exchange.PriceLevel{Price: lvl.Px, Quantity: lvl.Sz})
	}

	snapshot := &exchange.Snapshot{
		Exchange:     e.GetName(),
		Symbol:       res.Coin,
		LastUpdateID: 0,
		Bids:         bids,
		Asks:         asks,
		Timestamp:    time.UnixMilli(res.Time),
	}

	e.snapshotMu.Lock()
	e.snapshot = snapshot
	e.snapshotSet = true
	e.snapshotMu.Unlock()

	return snapshot, nil
}

func (e *FuturesExchange) Updates() <-chan *exchange.DepthUpdate { return e.updateChan }
func (e *FuturesExchange) IsConnected() bool                     { return e.wsConn != nil }
func (e *FuturesExchange) Health() exchange.HealthStatus {
	if s, ok := e.health.Load().(exchange.HealthStatus); ok {
		return s
	}
	return exchange.HealthStatus{}
}

func (e *FuturesExchange) pingLoop() {
	t := time.NewTicker(pingInterval)
	defer t.Stop()
	for {
		select {
		case <-e.ctx.Done():
			return
		case <-e.done:
			return
		case <-t.C:
			if e.wsConn != nil {
				_ = e.wsConn.WriteJSON(map[string]string{"method": "ping"})
			}
		}
	}
}

func (e *FuturesExchange) readMessages() {
	defer close(e.updateChan)
	defer e.updateConnectionStatus(false)

	for {
		select {
		case <-e.ctx.Done():
			log.Printf("[%s] context cancelled, stopping readMessages", e.GetName())
			return
		case <-e.done:
			return
		default:
			_, msgBytes, err := e.wsConn.ReadMessage()
			if err != nil {
				e.incrementErrorCount()
				log.Printf("[%s] WebSocket read error: %v", e.GetName(), err)
				return
			}

			var env struct {
				Channel string          `json:"channel"`
				Data    json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(msgBytes, &env); err != nil {
				log.Printf("[%s] failed to unmarshal envelope: %v", e.GetName(), err)
				continue
			}

			switch env.Channel {
			case "subscriptionResponse":
				// ack; ignore
				continue
			case "l2Book":
				var book WsBook
				if err := json.Unmarshal(env.Data, &book); err != nil {
					log.Printf("[%s] failed parse l2Book: %v", e.GetName(), err)
					continue
				}
				e.incrementMessageCount()
				e.updateLastPing()

				bids := make([]exchange.PriceLevel, 0, len(book.Levels[0]))
				for _, lvl := range book.Levels[0] {
					bids = append(bids, exchange.PriceLevel{Price: lvl.Px, Quantity: lvl.Sz})
				}
				asks := make([]exchange.PriceLevel, 0, len(book.Levels[1]))
				for _, lvl := range book.Levels[1] {
					asks = append(asks, exchange.PriceLevel{Price: lvl.Px, Quantity: lvl.Sz})
				}

				depth := &exchange.DepthUpdate{
					Exchange:      e.GetName(),
					Symbol:        book.Coin,
					EventTime:     time.UnixMilli(book.Time),
					FirstUpdateID: 0,
					FinalUpdateID: 0,
					PrevUpdateID:  0,
					Bids:          bids,
					Asks:          asks,
				}

				select {
				case e.updateChan <- depth:
				case <-e.ctx.Done():
					return
				case <-e.done:
					return
				default:
					log.Printf("[%s] Warning: update channel full, skipping update", e.GetName())
				}

			case "trades":
				var trades []WsTrade
				if err := json.Unmarshal(env.Data, &trades); err != nil {
					log.Printf("[%s] failed parse trades: %v", e.GetName(), err)
					continue
				}
				e.incrementMessageCount()
				e.updateLastPing()
				// Minimal behavior: log trades. If you want trades forwarded to UI, see notes below.
				for _, t := range trades {
					log.Printf("[%s] trade %s px=%s sz=%s side=%s time=%d", e.GetName(), t.Coin, t.Px, t.Sz, t.Side, t.Time)
				}
			default:
				// ignore other channels for now
				continue
			}
		}
	}
}

// health helpers
func (e *FuturesExchange) updateConnectionStatus(connected bool) {
	status := e.Health()
	status.Connected = connected
	if !connected {
		now := time.Now()
		status.ReconnectTime = &now
	}
	e.health.Store(status)
}

func (e *FuturesExchange) incrementMessageCount() {
	status := e.Health()
	status.MessageCount++
	e.health.Store(status)
}

func (e *FuturesExchange) incrementErrorCount() {
	status := e.Health()
	status.ErrorCount++
	e.health.Store(status)
}

func (e *FuturesExchange) updateLastPing() {
	status := e.Health()
	status.LastPing = time.Now()
	e.health.Store(status)
}

// mapToHLCoin converts symbols like BTCUSDT -> BTC
func mapToHLCoin(sym string) string {
	s := strings.ToUpper(sym)
	for _, suf := range []string{"USDT", "USD", "USDC"} {
		if strings.HasSuffix(s, suf) {
			return strings.TrimSuffix(s, suf)
		}
	}
	return s
}
