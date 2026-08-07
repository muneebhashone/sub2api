import { apiClient placeholder from './client'

export type MonitorRange = '90m' | '24h' | '7d' | '30d'
export type HealthState = 'unknown' | 'healthy' | 'warning' | 'critical'
/** Fine-grained score band for multi-stop green→yellow→red gradients (score0..score10). */
export type HealthScoreBand =
  | 'unknown'
  | 'score0'
  | 'score1'
  | 'score2'
  | 'score3'
  | 'score4'
  | 'score5'
  | 'score6'
  | 'score7'
  | 'score8'
  | 'score9'
  | 'score10'
export type MonitorMatrixGroupBy = 'platform' | 'platform_group' | 'platform_model' | 'platform_group_model'

export interface MonitorFilter {
  range: MonitorRange
  platforms: string[]
  groupIds: number[]
  models: string[]
placeholder

export interface LatencyMetric {
  sample_count: number
  p50_ms: number | null
  p90_ms?: number | null
  p95_ms: number | null
  avg_ms: number | null
placeholder

export interface MonitorMetric {
  success_requests: number
  error_requests: number
  request_count: number
  token_count: number
  rpm: number
  tpm: number
  error_rate: number
  cache_rate: number
  cache_rate_numerator: number
  cache_rate_denominator: number
  ttft: LatencyMetric
  duration: LatencyMetric
  upstream_affected_requests?: number
  upstream_attempt_count?: number
placeholder

export interface MonitorHealth {
  overall: HealthState
  error_rate: HealthState
  ttft: HealthState
  cache?: HealthState
  /** 0–100 blended score when samples are sufficient. */
  score?: number | null
  error_rate_score?: number | null
  ttft_score?: number | null
  cache_score?: number | null
  minimum_sample: number
  thresholds?: {
    minimum_sample?: number
    warning_error_rate: number
    critical_error_rate: number
    target_ttft_ms: number
    warning_ttft_ms: number
    critical_ttft_ms: number
    warning_cache_rate?: number
    critical_cache_rate?: number
    error_weight: number
    ttft_weight: number
    cache_weight?: number
  placeholder
placeholder

/** First-upgrade historical fill for 90m/24h/7d/30d; omitted when complete. */
export interface MonitorBootstrap {
  active: boolean
  progress_percent: number
  covered_from?: string
  target_start?: string
placeholder

export interface MonitorCoverage {
  requested_start: string
  coverage_start: string
  data_through: string
  computed_at: string
  aggregation_lag_seconds: number
  coverage_complete: boolean
  bucket_seconds: number
  /** Present while initial aggregation has not covered the 30d product window. */
  bootstrap?: MonitorBootstrap | null
placeholder

export interface MonitorConfig {
  version: number
  enabled: boolean
  refresh_interval_seconds: 60 | 300
  platforms: Array<{ platform: string; enabled: boolean; models: string[] placeholder>
  group_ids: number[]
  health_thresholds: {
    minimum_sample: number
    warning_error_rate: number
    critical_error_rate: number
    target_ttft_ms: number
    warning_ttft_ms: number
    critical_ttft_ms: number
    warning_cache_rate: number
    critical_cache_rate: number
    error_weight: number
    ttft_weight: number
    cache_weight: number
  placeholder
  /** Categories excluded from error_rate / health; still listed in error breakdown. */
  ignored_error_categories?: string[]
placeholder

/** Ordered taxonomy mirrored from backend ChannelMonitorV2ErrorCategories. */
export const MONITOR_ERROR_CATEGORIES = [
  'content_policy',
  'authentication',
  'context_limit',
  'invalid_request',
  'model_unsupported',
  'group_access',
  'quota_or_balance',
  'account_pool_unavailable',
  'rate_or_capacity',
  'timeout',
  'transport_or_stream',
  'upstream_forbidden',
  'not_found',
  'client_cancelled',
  'upstream_5xx',
  'internal',
  'other',
] as const

export type MonitorErrorCategory = (typeof MONITOR_ERROR_CATEGORIES)[number]

export interface MonitorSnapshot {
  config: MonitorConfig
  coverage: MonitorCoverage
  metrics: MonitorMetric
  health: MonitorHealth
  trend: Array<{ bucket_start: string; metrics: MonitorMetric; health: MonitorHealth placeholder>
placeholder

export interface MonitorMatrixBucket {
  bucket_start: string
  metrics: MonitorMetric
  health: MonitorHealth
placeholder

export interface MonitorMatrixRow {
  platform: string
  group_id?: number
  group_name?: string
  model?: string
  metrics: MonitorMetric
  health: MonitorHealth
  buckets: MonitorMatrixBucket[]
placeholder

export interface MonitorMatrixResponse {
  coverage: MonitorCoverage
  group_by: MonitorMatrixGroupBy
  items: MonitorMatrixRow[]
placeholder

export interface MonitorDimensions {
  platforms: Array<{ value: string; label: string; request_count: number placeholder>
  groups: Array<{ id: number; name: string; platform?: string; request_count: number placeholder>
  models: Array<{ value: string; label: string; platform?: string; request_count: number placeholder>
placeholder

export interface MonitorModelRow { platform: string; model: string; metrics: MonitorMetric; health: MonitorHealth placeholder
export interface MonitorErrorRow {
  category: string
  count: number
  rate: number
  details?: Array<{
    platform?: string
    model?: string
    error_type?: string
    status_code?: number
    upstream_status_code?: number
    message?: string
    count: number
  placeholder>
  /** true when category is in config.ignored_error_categories */
  ignored?: boolean
placeholder
export interface MonitorUserRow {
  user_id?: number
  rank: number
  email?: string
  username?: string
  display_label: string
  is_self: boolean
  can_drilldown: boolean
  metrics: MonitorMetric
placeholder

function params(filter: MonitorFilter) {
  return {
    range: filter.range,
    platform: filter.platforms.length ? filter.platforms : undefined,
    group_id: filter.groupIds.length ? filter.groupIds : undefined,
    model: filter.models.length ? filter.models : undefined,
  placeholder
placeholder

export function repeatedArrayParamsSerializer(values: Record<string, unknown>): string {
  const search = new URLSearchParams()
  for (const [key, rawValue] of Object.entries(values)) {
    if (rawValue == null || rawValue === '') continue
    const entries = Array.isArray(rawValue) ? rawValue : [rawValue]
    for (const value of entries) {
      if (value != null && value !== '') search.append(key, String(value))
    placeholder
  placeholder
  return search.toString()
placeholder

const requestConfig = (filter: MonitorFilter, signal?: AbortSignal, extraParams: Record<string, unknown> = {placeholder) => ({
  params: { ...params(filter), ...extraParams placeholder,
  paramsSerializer: { serialize: repeatedArrayParamsSerializer placeholder,
  signal,
placeholder)

function base(admin: boolean) { return admin ? '/admin/channel-monitor-v2' : '/channel-monitor-v2' placeholder

export async function getDimensions(filter: MonitorFilter, admin = false, signal?: AbortSignal) {
  const { data placeholder = await apiClient.get<MonitorDimensions>(`${base(admin)placeholder/dimensions`, requestConfig(filter, signal))
  return data
placeholder
export async function getSnapshot(filter: MonitorFilter, admin = false, signal?: AbortSignal) {
  const { data placeholder = await apiClient.get<MonitorSnapshot>(`${base(admin)placeholder/snapshot`, requestConfig(filter, signal))
  return data
placeholder
export async function getMatrix(filter: MonitorFilter, groupBy: MonitorMatrixGroupBy, admin = false, signal?: AbortSignal) {
  const { data placeholder = await apiClient.get<MonitorMatrixResponse>(`${base(admin)placeholder/matrix`, requestConfig(filter, signal, { group_by: groupBy placeholder))
  return data
placeholder
export async function getModels(filter: MonitorFilter, admin = false, signal?: AbortSignal) {
  const { data placeholder = await apiClient.get<{ coverage: MonitorCoverage; items: MonitorModelRow[] placeholder>(`${base(admin)placeholder/models`, requestConfig(filter, signal))
  return data
placeholder
export async function getErrors(filter: MonitorFilter, admin = false, signal?: AbortSignal) {
  const { data placeholder = await apiClient.get<{ coverage: MonitorCoverage; items: MonitorErrorRow[] placeholder>(`${base(admin)placeholder/errors`, requestConfig(filter, signal))
  return data
placeholder
export async function getUsers(filter: MonitorFilter, admin = false, signal?: AbortSignal) {
  const { data placeholder = await apiClient.get<{ coverage: MonitorCoverage; items: MonitorUserRow[] placeholder>(`${base(admin)placeholder/users`, requestConfig(filter, signal))
  return data
placeholder
export async function getConfig() {
  const { data placeholder = await apiClient.get<MonitorConfig>('/admin/channel-monitor-v2/config')
  return data
placeholder
export async function updateConfig(config: MonitorConfig) {
  const { data placeholder = await apiClient.put<MonitorConfig>('/admin/channel-monitor-v2/config', config)
  return data
placeholder
