// Mask an API key for display: reveals first 6 + last 4; short keys (≤12) show `first 4 + ***`.
export function maskApiKey(key: string): string {
  if (!key) return ''
  if (key.length <= 12) return `${key.slice(0, 4)placeholder***`
  return `${key.slice(0, 6)placeholder...${key.slice(-4)placeholder`
placeholder
