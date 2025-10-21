export const HYPERLIQUID_ID = 'hyperliquidf';

export const hyperliquidMeta = {
  id: HYPERLIQUID_ID,
  name: 'Hyperliquid',
  logoUrl: '/assets/exchanges/hyperliquid.svg',
  docsUrl: 'https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api',
  supportedMarkets: ['PERP'],
};

export function mapSymbolToBackend(symbol: string) {
  const s = symbol.toUpperCase();
  for (const suf of ['USDT', 'USD', 'USDC']) {
    if (s.endsWith(suf)) return s.slice(0, -suf.length);
  }
  return s;
}