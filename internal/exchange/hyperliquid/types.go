package hyperliquid

// Config holds configuration for Hyperliquid exchange
type Config struct {
	Symbol string
}

// SubscriptionMessage represents the subscription request to Hyperliquid WebSocket
// TODO: Replace with actual Hyperliquid subscription format once API documentation is reviewed
type SubscriptionMessage struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// WSMessage represents a WebSocket message from Hyperliquid
// TODO: Replace with actual Hyperliquid message structure once API documentation is reviewed
type WSMessage struct {
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// DepthData represents orderbook depth data from Hyperliquid
// TODO: Replace with actual Hyperliquid depth data format (snapshot/update distinction)
type DepthData struct {
	Action       string     `json:"action"`       // TODO: Determine if Hyperliquid uses "snapshot"/"update" or similar
	UpdateID     int64      `json:"updateId"`     // TODO: Verify field name for update ID
	Bids         [][]string `json:"bids"`         // TODO: Verify if format is [["price", "quantity"]] or different
	Asks         [][]string `json:"asks"`         // TODO: Verify if format is [["price", "quantity"]] or different
	Timestamp    int64      `json:"ts,omitempty"` // TODO: Verify timestamp field
}

// TradeData represents trade data from Hyperliquid
// TODO: Replace with actual Hyperliquid trade data format
type TradeData struct {
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Timestamp int64  `json:"timestamp"`
	Side      string `json:"side"` // "buy" or "sell"
}
