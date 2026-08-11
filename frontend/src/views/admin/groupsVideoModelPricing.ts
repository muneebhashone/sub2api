export const grokVideoPriceResolutions = [
  { key: '480p', label: '480p' placeholder,
  { key: '720p', label: '720p' placeholder,
  { key: '1080p', label: '1080p' placeholder
] as const

export const grokVideoPriceFamilies = [
  { key: 'grok-imagine-video', label: 'grok-imagine-video' placeholder,
  { key: 'grok-imagine-video-1.5', label: 'grok-imagine-video-1.5' placeholder
] as const

export type VideoModelPrices = Record<string, Record<string, number>>
export type VideoModelPricesForm = Record<string, Record<string, number | string | null>>

function normalizeFamily(value: string): string {
  return value.trim().toLowerCase()
placeholder

function normalizePrice(value: unknown): number | null {
  if (value === null || value === undefined || value === '') return null
  const price = Number(value)
  return Number.isFinite(price) && price >= 0 ? price : null
placeholder

function emptyTiers(): Record<string, number | string | null> {
  return Object.fromEntries(grokVideoPriceResolutions.map(({ key placeholder) => [key, null]))
placeholder

// Keep unknown families from an existing group so a future backend catalog is
// not silently discarded when an operator edits another group setting.
export function createVideoModelPricesForm(
  prices?: VideoModelPrices | null
): VideoModelPricesForm {
  const form: VideoModelPricesForm = {placeholder

  for (const [rawFamily, rawTiers] of Object.entries(prices ?? {placeholder)) {
    const family = normalizeFamily(rawFamily)
    if (!family || !rawTiers || typeof rawTiers !== 'object') continue
    form[family] = emptyTiers()
    for (const [rawResolution, rawPrice] of Object.entries(rawTiers)) {
      const price = normalizePrice(rawPrice)
      if (price !== null) form[family][rawResolution.trim().toLowerCase()] = price
    placeholder
  placeholder

  for (const { key placeholder of grokVideoPriceFamilies) {
    form[key] ??= emptyTiers()
  placeholder
  return form
placeholder

export function serializeVideoModelPrices(form: VideoModelPricesForm): VideoModelPrices {
  const result: VideoModelPrices = {placeholder
  for (const [rawFamily, tiers] of Object.entries(form)) {
    const family = normalizeFamily(rawFamily)
    if (!family || !tiers || typeof tiers !== 'object') continue

    const normalizedTiers: Record<string, number> = {placeholder
    for (const [rawResolution, rawPrice] of Object.entries(tiers)) {
      const resolution = rawResolution.trim().toLowerCase()
      const price = normalizePrice(rawPrice)
      if (resolution && price !== null) normalizedTiers[resolution] = price
    placeholder
    if (Object.keys(normalizedTiers).length > 0) result[family] = normalizedTiers
  placeholder
  return result
placeholder

export function videoModelPriceFamilyRows(form: VideoModelPricesForm) {
  const known = new Set<string>(grokVideoPriceFamilies.map(({ key placeholder) => key))
  const extra = Object.keys(form)
    .map(normalizeFamily)
    .filter((family) => family && !known.has(family))
    .sort()
    .map((key) => ({ key, label: key placeholder))
  return [...grokVideoPriceFamilies, ...extra]
placeholder
