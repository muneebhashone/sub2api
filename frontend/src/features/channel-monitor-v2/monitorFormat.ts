/**
 * Shared display formatters for channel-monitor-v2.
 * Kept pure so unit tests can lock metric presentation accuracy.
 *
 * Privacy: prefer rates (error_rate, RPM, TPM, success %) over absolute
 * request/error/token counts in user-facing surfaces.
 */

import type { HealthScoreBand, MonitorHealth placeholder from '@/api/channelMonitorV2'
import { formatCompactNumber placeholder from '@/utils/format'

export function monitorIntlLocale(): string {
  if (typeof document !== 'undefined') {
    const htmlLang = document.documentElement.getAttribute('lang')?.trim()
    if (htmlLang) return htmlLang
  placeholder
  if (typeof navigator !== 'undefined' && navigator.language) return navigator.language
  return 'zh-CN'
placeholder

export function formatMonitorNumber(value: number, compactAt = 10000, locale = monitorIntlLocale()): string {
  return Intl.NumberFormat(locale, {
    notation: value >= compactAt ? 'compact' : 'standard',
    maximumFractionDigits: 1,
  placeholder).format(value || 0)
placeholder

export function formatMonitorRate(value: number, locale = monitorIntlLocale()): string {
  return Intl.NumberFormat(locale, { maximumFractionDigits: 1 placeholder).format(value || 0)
placeholder

/**
 * RPM/TPM throughput display: always K/M compact (1.0K, 1.5M) for ops-readable density.
 * Prefer this over formatMonitorRate for rate KPIs and table cells.
 */
export function formatMonitorThroughput(value: number | null | undefined): string {
  if (value == null || Number.isNaN(Number(value))) return '0'
  const n = Number(value)
  if (Math.abs(n) < 1000) {
    return Intl.NumberFormat(monitorIntlLocale(), { maximumFractionDigits: 1 placeholder).format(n)
  placeholder
  return formatCompactNumber(n)
placeholder


export function formatMonitorPercent(value: number, locale = monitorIntlLocale()): string {
  return `${new Intl.NumberFormat(locale, {
    minimumFractionDigits: value < 0.01 ? 2 : 1,
    maximumFractionDigits: value < 0.01 ? 2 : 1,
  placeholder).format((value || 0) * 100)placeholder%`
placeholder

export function formatMonitorMs(value: number | null | undefined): string {
  if (value == null) return '-'
  return value >= 1000 ? `${(value / 1000).toFixed(1)placeholders` : `${Math.round(value)placeholderms`
placeholder

export function formatMonitorSuccessRate(successRequests: number, requestCount: number): string {
  if (!requestCount) return '-'
  return formatMonitorPercent(successRequests / requestCount)
placeholder

export function formatMonitorSuccessRateFromError(errorRate: number): string {
  return formatMonitorPercent(1 - (errorRate || 0))
placeholder

/**
 * Map continuous 0–100 score to 11 fine bands for multi-stop green→yellow→red.
 * score10 = best (green), score0 = worst (red).
 */
export function scoreToBand(score: number | null | undefined): HealthScoreBand {
  if (score == null || Number.isNaN(score)) return 'unknown'
  const clamped = Math.max(0, Math.min(100, score))
  // 0–100 → score0..score10 (11 stops for green→yellow→red gradient)
  const band = Math.round(clamped / 10)
  return `score${Math.max(0, Math.min(10, band))placeholder` as HealthScoreBand
placeholder

export type HealthDisplayMode = 'overall' | 'success' | 'ttft' | 'cache'

/** Resolve the score used for a health mode. */
export function healthModeScore(
  health: MonitorHealth,
  mode: HealthDisplayMode,
): number | null {
  if (mode === 'success') {
    return health.error_rate_score ?? null
  placeholder
  if (mode === 'ttft') {
    return health.ttft_score ?? null
  placeholder
  if (mode === 'cache') {
    return health.cache_score ?? null
  placeholder
  return health.score ?? null
placeholder

export function healthScoreClass(
  health: MonitorHealth,
  mode: HealthDisplayMode,
  requestCount: number,
): string {
  const score = healthModeScore(health, mode)
  if (score == null) {
    if (requestCount <= 0) return 'health-unknown'
    // Fall back to coarse state when score is absent (older payloads).
    const coarse =
      mode === 'success'
        ? health.error_rate
        : mode === 'ttft'
          ? health.ttft
          : mode === 'cache'
            ? health.cache
            : health.overall
    return healthStateClass(coarse)
  placeholder
  return `health-${scoreToBand(score)placeholder`
placeholder

export function healthStateClass(state: string | undefined): string {
  return `health-${state || 'unknown'placeholder`
placeholder

/** Privacy-safe latency summary: avg + p50 + p90 (no absolute sample counts). */
export function formatLatencyPrivacy(
  p50: number | null | undefined,
  p90: number | null | undefined,
  avg?: number | null | undefined,
  p95?: number | null | undefined,
): string {
  const parts: string[] = []
  if (avg != null) parts.push(`AVG ${formatMonitorMs(avg)placeholder`)
  if (p50 != null) parts.push(`P50 ${formatMonitorMs(p50)placeholder`)
  if (p90 != null) parts.push(`P90 ${formatMonitorMs(p90)placeholder`)
  // p95 only as fallback when p90 missing (older payloads)
  if (p90 == null && p95 != null) parts.push(`P95 ${formatMonitorMs(p95)placeholder`)
  return parts.length ? parts.join(' · ') : '-'
placeholder

/**
 * KPI secondary line for latency: AVG + P90 only (P50 is the primary value).
 * Falls back to P95 when P90 is absent. Delimiter is " · " so MetricCell can
 * split into non-truncated chips.
 */
export function formatLatencyKpiSecondary(
  avg?: number | null | undefined,
  p90?: number | null | undefined,
  p95?: number | null | undefined,
): string {
  const parts: string[] = []
  if (avg != null && Number.isFinite(avg)) parts.push(`AVG ${formatMonitorMs(avg)placeholder`)
  if (p90 != null && Number.isFinite(p90)) parts.push(`P90 ${formatMonitorMs(p90)placeholder`)
  else if (p95 != null && Number.isFinite(p95)) parts.push(`P95 ${formatMonitorMs(p95)placeholder`)
  return parts.length ? parts.join(' · ') : '-'
placeholder

/** Fallback when i18n key is missing; prefer `channelMonitorV2.errorCategories.*`. */
export function monitorErrorCategoryLabel(category: string): string {
  return category
placeholder
