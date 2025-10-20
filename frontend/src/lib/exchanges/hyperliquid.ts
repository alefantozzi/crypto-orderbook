/**
 * Hyperliquid Exchange Metadata
 * 
 * This file contains metadata and helper functions for integrating
 * Hyperliquid perpetual futures exchange into the orderbook aggregator.
 */

export const HYPERLIQUID_METADATA = {
  id: 'hyperliquid',
  name: 'Hyperliquid',
  logoUrl: '/assets/exchanges/hyperliquid.svg',
  supportedMarkets: ['perps'] as const,
  websocketUrl: 'wss://api.hyperliquid.xyz/ws', // TODO: Verify actual WebSocket URL
  restUrl: 'https://api.hyperliquid.xyz/info', // TODO: Verify actual REST API URL
  documentationUrl: 'https://hyperliquid.gitbook.io/hyperliquid-docs', // TODO: Add actual docs URL
};

/**
 * Maps a standard symbol format (e.g., "BTCUSDT") to Hyperliquid's expected format
 * 
 * TODO: Implement actual symbol mapping based on Hyperliquid's API specification
 * Hyperliquid may use different symbol formats for perpetual futures.
 * Common formats seen in other exchanges:
 * - BTC-PERP
 * - BTCUSD-PERP
 * - BTC-USD
 * 
 * @param symbol - Standard symbol in format like "BTCUSDT"
 * @returns Hyperliquid-specific symbol format
 */
export function mapSymbolToBackend(symbol: string): string {
  // TODO: Replace this placeholder implementation with actual Hyperliquid symbol mapping
  
  // Example placeholder logic (adjust based on Hyperliquid's actual format):
  // If Hyperliquid uses BTC-PERP format for BTC perpetual:
  const upperSymbol = symbol.toUpperCase();
  if (upperSymbol.indexOf('USDT') === upperSymbol.length - 4) {
    const base = upperSymbol.substring(0, upperSymbol.length - 4);
    return `${base}-PERP`; // Placeholder format
  }
  
  if (upperSymbol.indexOf('USD') === upperSymbol.length - 3) {
    const base = upperSymbol.substring(0, upperSymbol.length - 3);
    return `${base}-PERP`; // Placeholder format
  }
  
  // Default: return as-is
  return upperSymbol;
}

/**
 * Symbol mapping table for common trading pairs
 * TODO: Populate with actual Hyperliquid symbol mappings
 */
export const SYMBOL_MAPPINGS: Record<string, string> = {
  'BTCUSDT': 'BTC-PERP',  // TODO: Verify
  'ETHUSDT': 'ETH-PERP',  // TODO: Verify
  'SOLUSDT': 'SOL-PERP',  // TODO: Verify
  // Add more mappings as needed
};

/**
 * Get Hyperliquid-specific symbol from standard format using the mapping table
 * @param symbol - Standard symbol format
 * @returns Hyperliquid symbol format
 */
export function getHyperliquidSymbol(symbol: string): string {
  const upperSymbol = symbol.toUpperCase();
  return SYMBOL_MAPPINGS[upperSymbol] || mapSymbolToBackend(symbol);
}

export default HYPERLIQUID_METADATA;
