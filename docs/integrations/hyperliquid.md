```markdown
# Hyperliquid integration notes

- WS: wss://api.hyperliquid.xyz/ws
- Subscribe: {"method":"subscribe","subscription": {"type":"l2Book","coin":"BTC"}}
- Ping: {"method":"ping"} every 50s
- REST snapshot: POST https://api.hyperliquid.xyz/info  { "type":"l2Book", "coin": "<coin>" }

Quick test (local)
1. Create branch and add files (see commands below).
2. Backend:
   go build -o crypto-orderbook ./cmd
   ./crypto-orderbook -symbol BTCUSDT
   - The connector will map BTCUSDT -> BTC for Hyperliquid.
   - Look for logs:
     [hyperliquidf] WebSocket connected to wss://api.hyperliquid.xyz/ws
     [hyperliquidf] Fetching L2 snapshot via REST for BTCUSDT
     [hyperliquidf] trade BTC px=... sz=... side=... time=...
3. Frontend:
   cd frontend
   npm install
   npm run dev
   - Open the Vite URL (usually http://localhost:5173) and confirm the UI is connected (status indicator) and Hyperliquid orderbook updates appear.

TODOs:
- Optionally add a trade-forwarding pipeline (add exchange.Trade, Trades() channel, server broadcast, frontend trade UI).
- Replace any additional brand-kit variants if you want theme-specific logos.
```