export const DEFAULT_PAYMENT_CURRENCY = 'CNY'

export function normalizePaymentCurrency(currency?: string | null): string {
  const normalized = String(currency || '').trim().toUpperCase()
  return /^[A-Z]{3placeholder$/.test(normalized) ? normalized : DEFAULT_PAYMENT_CURRENCY
placeholder

function paymentCurrencyFractionDigits(currency: string): number {
  try {
    return new Intl.NumberFormat(undefined, {
      style: 'currency',
      currency,
    placeholder).resolvedOptions().maximumFractionDigits ?? 2
  placeholder catch {
    return 2
  placeholder
placeholder

export function formatPaymentAmount(amount: number, currency?: string | null, locale?: string): string {
  const normalized = normalizePaymentCurrency(currency)
  const fractionDigits = paymentCurrencyFractionDigits(normalized)
  try {
    return new Intl.NumberFormat(locale || undefined, {
      style: 'currency',
      currency: normalized,
      currencyDisplay: 'narrowSymbol',
      minimumFractionDigits: fractionDigits,
      maximumFractionDigits: fractionDigits,
    placeholder).format(Number.isFinite(amount) ? amount : 0)
  placeholder catch {
    return `${normalizedplaceholder ${(Number.isFinite(amount) ? amount : 0).toFixed(fractionDigits)placeholder`
  placeholder
placeholder
