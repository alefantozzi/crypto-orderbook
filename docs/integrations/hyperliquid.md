# Hyperliquid Exchange Integration

This document describes the integration of Hyperliquid perpetual futures exchange into the crypto orderbook aggregator.

## Overview

Hyperliquid is a decentralized perpetual futures exchange. This integration adds support for aggregating Hyperliquid orderbook data alongside other centralized and decentralized exchanges.

## Implementation Status

The connector has been scaffolded with placeholder implementations following the patterns used by existing exchange connectors (BingX, Bybit, etc.). The following components have been created:

### Backend Components

1. **internal/exchange/hyperliquid/types.go**
   - Config struct for exchange configuration
   - Message type definitions (placeholder structures)
   - DepthData and TradeData types

2. **internal/exchange/hyperliquid/futures.go**
   - FuturesExchange struct implementing the Exchange interface
   - WebSocket connection management
   - Message handling routines (with TODOs for parsing)
   - Snapshot and update conversion functions

3. **internal/factory/factory.go**
   - Updated to include Hyperliquid in the factory pattern
   - Added to list of implemented exchanges

### Frontend Components

1. **frontend/src/lib/exchanges/hyperliquid.ts**
   - Exchange metadata (id, name, URLs)
   - Symbol mapping functions (placeholder)
   - Helper utilities for symbol conversion

2. **frontend/src/components/ExchangeIcon.tsx**
   - Updated to include Hyperliquid icon component (placeholder SVG)

## TODOs and Next Steps

### 1. Fetch Hyperliquid API Documentation

- [ ] Obtain official Hyperliquid WebSocket API documentation
- [ ] Identify correct WebSocket URL for perpetual futures
- [ ] Identify REST API endpoint for orderbook snapshots (if available)
- [ ] Document authentication requirements (if any)

### 2. Implement WebSocket Message Parsing

#### Subscription Messages
- [ ] Replace placeholder subscription format in `futures.go` Connect() method
- [ ] Determine correct channel/topic naming for orderbook streams
- [ ] Implement trade subscription if needed
- [ ] Handle subscription confirmation messages

#### Response Messages
- [ ] Update `types.go` with actual Hyperliquid message structures
- [ ] Implement parsing for snapshot messages
- [ ] Implement parsing for incremental update messages
- [ ] Handle trade messages (if needed)
- [ ] Implement ping/pong handling based on Hyperliquid's requirements

#### Message Differentiation
- [ ] Determine how Hyperliquid distinguishes between snapshot and update messages
- [ ] Update `handleMessage()` in `futures.go` to properly detect message types
- [ ] Implement proper sequence ID tracking for update continuity

### 3. Symbol Mapping

- [ ] Document Hyperliquid's symbol naming conventions for perpetuals
- [ ] Update `SYMBOL_MAPPINGS` table in `frontend/src/lib/exchanges/hyperliquid.ts`
- [ ] Implement `mapSymbolToBackend()` function with correct conversion logic
- [ ] Add common trading pairs (BTC, ETH, SOL, etc.)

### 4. Testing

#### Local Development Testing
- [ ] Set up local development environment
- [ ] Configure environment variables for Hyperliquid connection
- [ ] Test WebSocket connection establishment
- [ ] Verify snapshot reception
- [ ] Verify incremental update processing
- [ ] Test reconnection logic

#### Integration Testing
- [ ] Test with live Hyperliquid API
- [ ] Verify orderbook accuracy against Hyperliquid's UI
- [ ] Test with multiple symbols
- [ ] Measure latency and performance
- [ ] Test error handling and edge cases

#### Frontend Testing
- [ ] Test exchange metadata display
- [ ] Verify symbol mapping in UI
- [ ] Test orderbook visualization
- [ ] Check statistics calculation

### 5. Frontend Assets

- [ ] Obtain official Hyperliquid logo (SVG format preferred)
- [ ] Add logo to `frontend/public/assets/exchanges/hyperliquid.svg`
- [ ] Update ExchangeIcon component to use actual logo
- [ ] Test icon rendering in different contexts (light/dark mode)

### 6. Documentation Updates

- [ ] Update main README.md to list Hyperliquid as supported exchange
- [ ] Document configuration options
- [ ] Add example usage
- [ ] Document any limitations or known issues

## Configuration

### Environment Variables

```bash
# Example configuration (update with actual values)
HYPERLIQUID_WS_URL="wss://api.hyperliquid.xyz/ws"
HYPERLIQUID_REST_URL="https://api.hyperliquid.xyz/info"
```

### Running the Connector

```bash
# Build the project
go build ./cmd/...

# Run with Hyperliquid
./orderbook --exchange hyperliquid --symbol BTCUSDT
```

## Code Structure

### Key Files

- `internal/exchange/hyperliquid/futures.go` - Main connector implementation
- `internal/exchange/hyperliquid/types.go` - Type definitions
- `internal/factory/factory.go` - Factory integration
- `frontend/src/lib/exchanges/hyperliquid.ts` - Frontend metadata

### Key Functions to Implement

1. **Connect()** - WebSocket connection and subscription
2. **handleMessage()** - Parse and route incoming messages
3. **convertSnapshot()** - Convert Hyperliquid snapshot to canonical format
4. **convertDepthUpdate()** - Convert Hyperliquid updates to canonical format
5. **mapSymbolToBackend()** - Frontend symbol conversion

## Testing Checklist

- [ ] Code compiles without errors ✅
- [ ] WebSocket connection establishes successfully
- [ ] Snapshot received and parsed correctly
- [ ] Incremental updates received and applied
- [ ] Orderbook maintains correct state
- [ ] Frontend displays Hyperliquid data
- [ ] Symbol mapping works correctly
- [ ] Reconnection logic functions properly
- [ ] Error handling covers edge cases
- [ ] Performance is acceptable (latency, memory)

## Resources

- Hyperliquid Website: https://hyperliquid.xyz
- Hyperliquid Documentation: (TODO: Add link)
- WebSocket API Docs: (TODO: Add link)
- REST API Docs: (TODO: Add link)

## Notes

- Hyperliquid is a decentralized exchange, so connection patterns may differ from CEX integrations
- Verify if authentication is required for market data feeds
- Check rate limits and connection policies
- Consider if L2 data depth is configurable

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Verify WebSocket URL is correct
   - Check if authentication is required
   - Ensure network connectivity

2. **No Snapshot Received**
   - Verify subscription message format
   - Check if symbol format is correct
   - Review WebSocket logs for errors

3. **Update Continuity Gaps**
   - Verify sequence ID tracking
   - Check update ID fields in messages
   - Ensure snapshot is fully received before processing updates

4. **Symbol Not Found**
   - Update symbol mapping table
   - Verify symbol format with Hyperliquid docs
   - Check if symbol exists on Hyperliquid

## Contributing

When implementing the TODOs:

1. Test each change incrementally
2. Add logging for debugging
3. Follow existing code patterns from other connectors
4. Update this documentation with findings
5. Add unit tests where applicable

## Support

For issues specific to this integration, refer to:
- Project repository issues
- Hyperliquid community channels
- API documentation (once available)
