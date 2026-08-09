/**
 * Admin Accounts API endpoints
 * Handles AI platform account management for administrators
 */

import { apiClient placeholder from '../client'
import type {
  Account,
  CreateAccountRequest,
  UpdateAccountRequest,
  PaginatedResponse,
  AccountUsageInfo,
  WindowStats,
  ClaudeModel,
  AccountUsageStatsResponse,
  TempUnschedulableStatus,
  AdminDataPayload,
  AdminDataImportResult,
  CodexSessionImportRequest,
  CodexSessionImportResult,
  OpenAICodexPATCreateRequest,
  CheckMixedChannelRequest,
  CheckMixedChannelResponse,
  UpstreamBillingProbeResult,
  UpstreamBillingProbeSettings,
  OllamaCloudUsageSettings,
  OllamaCloudUsageState
placeholder from '@/types'

/**
 * List all accounts with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters
 * @returns Paginated list of accounts
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    platform?: string
    type?: string
    status?: string
    group?: string
    search?: string
    privacy_mode?: string
    lite?: string
    include_scheduler_score?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  placeholder,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<Account>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<Account>>('/admin/accounts', {
    params: {
      page,
      page_size: pageSize,
      ...filters
    placeholder,
    signal: options?.signal
  placeholder)
  return data
placeholder

export interface AccountListWithEtagResult {
  notModified: boolean
  etag: string | null
  data: PaginatedResponse<Account> | null
placeholder

export async function listWithEtag(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    platform?: string
    type?: string
    status?: string
    group?: string
    search?: string
    privacy_mode?: string
    lite?: string
    include_scheduler_score?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  placeholder,
  options?: {
    signal?: AbortSignal
    etag?: string | null
  placeholder
): Promise<AccountListWithEtagResult> {
  const headers: Record<string, string> = {placeholder
  if (options?.etag) {
    headers['If-None-Match'] = options.etag
  placeholder

  const response = await apiClient.get<PaginatedResponse<Account>>('/admin/accounts', {
    params: {
      page,
      page_size: pageSize,
      ...filters
    placeholder,
    headers,
    signal: options?.signal,
    validateStatus: (status) => (status >= 200 && status < 300) || status === 304
  placeholder)

  const etagHeader = typeof response.headers?.etag === 'string' ? response.headers.etag : null
  if (response.status === 304) {
    return {
      notModified: true,
      etag: etagHeader,
      data: null
    placeholder
  placeholder

  return {
    notModified: false,
    etag: etagHeader,
    data: response.data
  placeholder
placeholder

/**
 * Get account by ID
 * @param id - Account ID
 * @returns Account details
 */
export async function getById(id: number): Promise<Account> {
  const { data placeholder = await apiClient.get<Account>(`/admin/accounts/${idplaceholder`)
  return data
placeholder

/**
 * Create new account
 * @param accountData - Account data
 * @returns Created account
 */
export async function create(accountData: CreateAccountRequest): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>('/admin/accounts', accountData)
  return data
placeholder

/**
 * Duplicate an account while keeping credentials on the server.
 * @param id - Source account ID
 * @returns Newly created account
 */
const duplicateOperationKeys = new Map<number, string>()

function duplicateOperationStorageKey(id: number): string {
  return `sub2api:admin:account-duplicate:${idplaceholder`
placeholder

function getStoredDuplicateOperationKey(id: number): string | null {
  try {
    return globalThis.sessionStorage?.getItem(duplicateOperationStorageKey(id)) ?? null
  placeholder catch {
    return null
  placeholder
placeholder

function storeDuplicateOperationKey(id: number, key: string | null): void {
  try {
    if (key) globalThis.sessionStorage?.setItem(duplicateOperationStorageKey(id), key)
    else globalThis.sessionStorage?.removeItem(duplicateOperationStorageKey(id))
  placeholder catch {
    // In-memory retry protection still works when browser storage is unavailable.
  placeholder
placeholder

export async function duplicate(id: number): Promise<Account> {
  let idempotencyKey = duplicateOperationKeys.get(id) ?? getStoredDuplicateOperationKey(id)
  if (!idempotencyKey) {
    const requestID = globalThis.crypto?.randomUUID?.() ?? `${Date.now()placeholder-${Math.random().toString(36).slice(2)placeholder`
    idempotencyKey = `account-duplicate-${idplaceholder-${requestIDplaceholder`
  placeholder
  duplicateOperationKeys.set(id, idempotencyKey)
  storeDuplicateOperationKey(id, idempotencyKey)
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${idplaceholder/duplicate`, undefined, {
    headers: { 'Idempotency-Key': idempotencyKey placeholder
  placeholder)
  duplicateOperationKeys.delete(id)
  storeDuplicateOperationKey(id, null)
  return data
placeholder

/**
 * Update account
 * @param id - Account ID
 * @param updates - Fields to update
 * @returns Updated account
 */
export async function update(id: number, updates: UpdateAccountRequest): Promise<Account> {
  const { data placeholder = await apiClient.put<Account>(`/admin/accounts/${idplaceholder`, updates)
  return data
placeholder

/**
 * Check mixed-channel risk for account-group binding.
 */
export async function checkMixedChannelRisk(
  payload: CheckMixedChannelRequest
): Promise<CheckMixedChannelResponse> {
  const { data placeholder = await apiClient.post<CheckMixedChannelResponse>('/admin/accounts/check-mixed-channel', payload)
  return data
placeholder

/**
 * Delete account
 * @param id - Account ID
 * @returns Success confirmation
 */
export async function deleteAccount(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/accounts/${idplaceholder`)
  return data
placeholder

/**
 * Toggle account status
 * @param id - Account ID
 * @param status - New status
 * @returns Updated account
 */
export async function toggleStatus(id: number, status: 'active' | 'inactive'): Promise<Account> {
  return update(id, { status placeholder)
placeholder

/**
 * Test account connectivity
 * @param id - Account ID
 * @returns Test result
 */
export async function testAccount(id: number): Promise<{
  success: boolean
  message: string
  latency_ms?: number
placeholder> {
  const { data placeholder = await apiClient.post<{
    success: boolean
    message: string
    latency_ms?: number
  placeholder>(`/admin/accounts/${idplaceholder/test`)
  return data
placeholder

/**
 * Refresh account credentials
 * @param id - Account ID
 * @returns Updated account
 */
export async function refreshCredentials(id: number): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${idplaceholder/refresh`)
  return data
placeholder

/**
 * Apply OAuth credentials after re-authorization.
 *
 * Unlike `update()`, this endpoint:
 * - never overwrites the whole `extra` JSONB (merges incrementally instead),
 *   so persistent settings like `base_rpm`, `window_cost_limit`, `max_sessions`,
 *   `quota_*` and `privacy_mode` are preserved
 * - clears the account error and invalidates the token cache server-side
 */
export async function applyOAuthCredentials(
  id: number,
  payload: {
    type: 'oauth' | 'setup-token'
    credentials: Record<string, unknown>
    extra?: Record<string, unknown>
  placeholder
): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(
    `/admin/accounts/${idplaceholder/apply-oauth-credentials`,
    payload
  )
  return data
placeholder

/**
 * Get account usage statistics
 * @param id - Account ID
 * @param days - Number of days (default: 30)
 * @returns Account usage statistics with history, summary, and models
 */
export async function getStats(id: number, days: number = 30): Promise<AccountUsageStatsResponse> {
  const { data placeholder = await apiClient.get<AccountUsageStatsResponse>(`/admin/accounts/${idplaceholder/stats`, {
    params: { days placeholder
  placeholder)
  return data
placeholder

/**
 * Clear account error
 * @param id - Account ID
 * @returns Updated account
 */
export async function clearError(id: number): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${idplaceholder/clear-error`)
  return data
placeholder

/**
 * Get account usage information (5h/7d window)
 * @param id - Account ID
 * @returns Account usage info
 */
export async function getUsage(id: number, source?: 'passive' | 'active', force?: boolean): Promise<AccountUsageInfo> {
  const params: Record<string, string> = {placeholder
  if (source) params.source = source
  if (force) params.force = 'true'
  const { data placeholder = await apiClient.get<AccountUsageInfo>(`/admin/accounts/${idplaceholder/usage`, {
    params: Object.keys(params).length > 0 ? params : undefined
  placeholder)
  return data
placeholder

export interface BatchAccountUsageResponse {
  usage: Record<string, AccountUsageInfo>
  errors: Record<string, string>
placeholder

export async function getBatchUsage(accountIds: number[], force?: boolean): Promise<BatchAccountUsageResponse> {
  const { data placeholder = await apiClient.post<BatchAccountUsageResponse>('/admin/accounts/usage/batch', {
    account_ids: accountIds,
    force: force === true
  placeholder)
  return data
placeholder

/**
 * Clear account rate limit status
 * @param id - Account ID
 * @returns Updated account
 */
export async function clearRateLimit(id: number): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(
    `/admin/accounts/${idplaceholder/clear-rate-limit`
  )
  return data
placeholder

/**
 * Recover account runtime state in one call
 * @param id - Account ID
 * @returns Updated account
 */
export async function recoverState(id: number): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${idplaceholder/recover-state`)
  return data
placeholder

/**
 * Reset account quota usage
 * @param id - Account ID
 * @returns Updated account
 */
export async function resetAccountQuota(id: number): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(
    `/admin/accounts/${idplaceholder/reset-quota`
  )
  return data
placeholder

/**
 * Get temporary unschedulable status
 * @param id - Account ID
 * @returns Status with detail state if active
 */
export async function getTempUnschedulableStatus(id: number): Promise<TempUnschedulableStatus> {
  const { data placeholder = await apiClient.get<TempUnschedulableStatus>(
    `/admin/accounts/${idplaceholder/temp-unschedulable`
  )
  return data
placeholder

/**
 * Reset temporary unschedulable status
 * @param id - Account ID
 * @returns Success confirmation
 */
export async function resetTempUnschedulable(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(
    `/admin/accounts/${idplaceholder/temp-unschedulable`
  )
  return data
placeholder

/**
 * Generate OAuth authorization URL
 * @param endpoint - API endpoint path
 * @param config - Proxy configuration
 * @returns Auth URL and session ID
 */
export async function generateAuthUrl(
  endpoint: string,
  config: { proxy_id?: number placeholder
): Promise<{ auth_url: string; session_id: string placeholder> {
  const { data placeholder = await apiClient.post<{ auth_url: string; session_id: string placeholder>(endpoint, config)
  return data
placeholder

/**
 * Exchange authorization code for tokens
 * @param endpoint - API endpoint path
 * @param exchangeData - Session ID, code, and optional proxy config
 * @returns Token information
 */
export async function exchangeCode(
  endpoint: string,
  exchangeData: { session_id: string; code: string; state?: string; proxy_id?: number placeholder
): Promise<Record<string, unknown>> {
  const { data placeholder = await apiClient.post<Record<string, unknown>>(endpoint, exchangeData)
  return data
placeholder

/**
 * Batch create accounts
 * @param accounts - Array of account data
 * @returns Results of batch creation
 */
export async function batchCreate(accounts: CreateAccountRequest[]): Promise<{
  success: number
  failed: number
  results: Array<{ success: boolean; account?: Account; error?: string placeholder>
placeholder> {
  const { data placeholder = await apiClient.post<{
    success: number
    failed: number
    results: Array<{ success: boolean; account?: Account; error?: string placeholder>
  placeholder>('/admin/accounts/batch', { accounts placeholder)
  return data
placeholder

/**
 * Batch update credentials fields for multiple accounts
 * @param request - Batch update request containing account IDs, field name, and value
 * @returns Results of batch update
 */
export async function batchUpdateCredentials(request: {
  account_ids: number[]
  field: string
  value: any
placeholder): Promise<{
  success: number
  failed: number
  results: Array<{ account_id: number; success: boolean; error?: string placeholder>
placeholder> {
  const { data placeholder = await apiClient.post<{
    success: number
    failed: number
    results: Array<{ account_id: number; success: boolean; error?: string placeholder>
  placeholder>('/admin/accounts/batch-update-credentials', request)
  return data
placeholder

/**
 * Bulk update multiple accounts
 * @param accountIds - Array of account IDs
 * @param updates - Fields to update
 * @returns Success confirmation
 */
export async function bulkUpdate(
  accountIdsOrPayload: number[] | Record<string, unknown>,
  updates?: Record<string, unknown>
): Promise<{
  success: number
  failed: number
  success_ids?: number[]
  failed_ids?: number[]
  results: Array<{ account_id: number; success: boolean; error?: string placeholder>
  placeholder> {
  const payload = Array.isArray(accountIdsOrPayload)
    ? {
        account_ids: accountIdsOrPayload,
        ...(updates ?? {placeholder)
      placeholder
    : accountIdsOrPayload
  const { data placeholder = await apiClient.post<{
    success: number
    failed: number
    success_ids?: number[]
    failed_ids?: number[]
    results: Array<{ account_id: number; success: boolean; error?: string placeholder>
  placeholder>('/admin/accounts/bulk-update', payload)
  return data
placeholder

/**
 * Get account today statistics
 * @param id - Account ID
 * @returns Today's stats (requests, tokens, cost)
 */
export async function getTodayStats(id: number): Promise<WindowStats> {
  const { data placeholder = await apiClient.get<WindowStats>(`/admin/accounts/${idplaceholder/today-stats`)
  return data
placeholder

export interface BatchTodayStatsResponse {
  stats: Record<string, WindowStats>
placeholder

/**
 * 批量获取多个账号的今日统计
 * @param accountIds - 账号 ID 列表
 * @returns 以账号 ID（字符串）为键的统计映射
 */
export async function getBatchTodayStats(accountIds: number[]): Promise<BatchTodayStatsResponse> {
  const { data placeholder = await apiClient.post<BatchTodayStatsResponse>('/admin/accounts/today-stats/batch', {
    account_ids: accountIds
  placeholder)
  return data
placeholder

/**
 * Set account schedulable status
 * @param id - Account ID
 * @param schedulable - Whether the account should participate in scheduling
 * @returns Updated account
 */
export async function setSchedulable(id: number, schedulable: boolean): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${idplaceholder/schedulable`, {
    schedulable
  placeholder)
  return data
placeholder

/**
 * Get available models for an account
 * @param id - Account ID
 * @returns List of available models for this account
 */
export async function getAvailableModels(id: number): Promise<ClaudeModel[]> {
  const { data placeholder = await apiClient.get<ClaudeModel[]>(`/admin/accounts/${idplaceholder/models`)
  return data
placeholder

export interface SyncUpstreamModelsResult {
  models: string[]
placeholder

/**
 * Sync live supported models from the account's upstream model-list endpoint
 * @param id - Account ID
 * @returns List of model IDs returned by the upstream
 */
export async function syncUpstreamModels(id: number): Promise<SyncUpstreamModelsResult> {
  const { data placeholder = await apiClient.post<SyncUpstreamModelsResult>(`/admin/accounts/${idplaceholder/models/sync-upstream`)
  return data
placeholder

export interface SyncUpstreamPreviewParams {
  platform: string
  type: string
  base_url?: string
  api_key: string
placeholder

/**
 * Preview upstream models without a saved account (create-flow)
 * @param params - Connection credentials
 * @returns List of model IDs returned by the upstream
 */
export async function syncUpstreamModelsPreview(params: SyncUpstreamPreviewParams): Promise<SyncUpstreamModelsResult> {
  const { data placeholder = await apiClient.post<SyncUpstreamModelsResult>('/admin/accounts/models/sync-upstream-preview', params)
  return data
placeholder

export interface CRSPreviewAccount {
  crs_account_id: string
  kind: string
  name: string
  platform: string
  type: string
placeholder

export interface PreviewFromCRSResult {
  new_accounts: CRSPreviewAccount[]
  existing_accounts: CRSPreviewAccount[]
placeholder

export async function previewFromCrs(params: {
  base_url: string
  username: string
  password: string
placeholder): Promise<PreviewFromCRSResult> {
  const { data placeholder = await apiClient.post<PreviewFromCRSResult>('/admin/accounts/sync/crs/preview', params)
  return data
placeholder

export async function syncFromCrs(params: {
  base_url: string
  username: string
  password: string
  sync_proxies?: boolean
  selected_account_ids?: string[]
placeholder): Promise<{
  created: number
  updated: number
  skipped: number
  failed: number
  items: Array<{
    crs_account_id: string
    kind: string
    name: string
    action: string
    error?: string
  placeholder>
placeholder> {
  const { data placeholder = await apiClient.post<{
    created: number
    updated: number
    skipped: number
    failed: number
    items: Array<{
      crs_account_id: string
      kind: string
      name: string
      action: string
      error?: string
    placeholder>
  placeholder>('/admin/accounts/sync/crs', params, {
    timeout: 180000 // 180s timeout: sync refreshes each existing account's OAuth token serially
  placeholder)
  return data
placeholder

export async function exportData(options?: {
  ids?: number[]
  filters?: {
    platform?: string
    type?: string
    status?: string
    group?: string
    privacy_mode?: string
    search?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  placeholder
  includeProxies?: boolean
placeholder): Promise<AdminDataPayload> {
  const params: Record<string, string> = {placeholder
  if (options?.ids && options.ids.length > 0) {
    params.ids = options.ids.join(',')
  placeholder else if (options?.filters) {
    const { platform, type, status, group, privacy_mode, search, sort_by, sort_order placeholder = options.filters
    if (platform) params.platform = platform
    if (type) params.type = type
    if (status) params.status = status
    if (group) params.group = group
    if (privacy_mode) params.privacy_mode = privacy_mode
    if (search) params.search = search
    if (sort_by) params.sort_by = sort_by
    if (sort_order) params.sort_order = sort_order
  placeholder
  if (options?.includeProxies === false) {
    params.include_proxies = 'false'
  placeholder
  const { data placeholder = await apiClient.get<AdminDataPayload>('/admin/accounts/data', { params placeholder)
  return data
placeholder

export async function importData(payload: {
  data: AdminDataPayload
  skip_default_group_bind?: boolean
placeholder): Promise<AdminDataImportResult> {
  const { data placeholder = await apiClient.post<AdminDataImportResult>('/admin/accounts/data', {
    data: payload.data,
    skip_default_group_bind: payload.skip_default_group_bind
  placeholder)
  return data
placeholder

export async function importCodexSession(payload: CodexSessionImportRequest): Promise<CodexSessionImportResult> {
  const { data placeholder = await apiClient.post<CodexSessionImportResult>('/admin/accounts/import/codex-session', payload, {
    timeout: 120000 // 120s timeout for large session imports
  placeholder)
  return data
placeholder

export async function createOpenAICodexPAT(payload: OpenAICodexPATCreateRequest): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>('/admin/openai/create-from-codex-pat', payload)
  return data
placeholder

/**
 * Get Antigravity default model mapping from backend
 * @returns Default model mapping (from -> to)
 */
export async function getAntigravityDefaultModelMapping(): Promise<Record<string, string>> {
  const { data placeholder = await apiClient.get<Record<string, string>>(
    '/admin/accounts/antigravity/default-model-mapping'
  )
  return data
placeholder

/**
 * Refresh OpenAI token using refresh token
 * @param refreshToken - The refresh token
 * @param proxyId - Optional proxy ID
 * @returns Token information including access_token, email, etc.
 */
export async function refreshOpenAIToken(
  refreshToken: string,
  proxyId?: number | null,
  endpoint: string = '/admin/openai/refresh-token',
  clientId?: string
): Promise<Record<string, unknown>> {
  const payload: { refresh_token: string; proxy_id?: number; client_id?: string placeholder = {
    refresh_token: refreshToken
  placeholder
  if (proxyId) {
    payload.proxy_id = proxyId
  placeholder
  if (clientId) {
    payload.client_id = clientId
  placeholder
  const { data placeholder = await apiClient.post<Record<string, unknown>>(endpoint, payload)
  return data
placeholder

/**
 * Batch operation result type
 */
export interface BatchOperationResult {
  total: number
  success: number
  failed: number
  success_ids?: number[]
  failed_ids?: number[]
  errors?: Array<{ account_id: number; error: string placeholder>
  warnings?: Array<{ account_id: number; warning: string placeholder>
placeholder

/**
 * Revert account proxy to original before fallback
 * @param id - Account ID
 * @returns Success confirmation
 */
export async function revertProxyFallback(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.post<{ message: string placeholder>(`/admin/accounts/${idplaceholder/revert-proxy-fallback`)
  return data
placeholder

/**
 * Delete multiple accounts with bounded server-side concurrency.
 */
export async function batchDelete(accountIds: number[]): Promise<BatchOperationResult> {
  const { data placeholder = await apiClient.post<BatchOperationResult>('/admin/accounts/batch-delete', {
    account_ids: accountIds
  placeholder)
  return data
placeholder

/**
 * Batch clear account errors
 * @param accountIds - Array of account IDs
 * @returns Batch operation result
 */
export async function batchClearError(accountIds: number[]): Promise<BatchOperationResult> {
  const { data placeholder = await apiClient.post<BatchOperationResult>('/admin/accounts/batch-clear-error', {
    account_ids: accountIds
  placeholder)
  return data
placeholder

/**
 * Batch refresh account credentials
 * @param accountIds - Array of account IDs
 * @returns Batch operation result
 */
export async function batchRefresh(accountIds: number[]): Promise<BatchOperationResult> {
  const { data placeholder = await apiClient.post<BatchOperationResult>('/admin/accounts/batch-refresh', {
    account_ids: accountIds,
  placeholder, {
    timeout: 120000  // 120s timeout for large batch refreshes
  placeholder)
  return data
placeholder

/**
 * Set privacy for an Antigravity OAuth account
 * @param id - Account ID
 * @returns Updated account
 */
export async function setPrivacy(id: number): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${idplaceholder/set-privacy`)
  return data
placeholder

/**
 * OpenAI / Codex rate-limit reset feature: query and reset upstream usage.
 */
export interface OpenAIRateLimitWindow {
  used_percent: number
  limit_window_seconds: number
  reset_after_seconds: number
  reset_at: number
placeholder

export interface OpenAIRateLimit {
  allowed: boolean
  limit_reached: boolean
  primary_window?: OpenAIRateLimitWindow | null
  secondary_window?: OpenAIRateLimitWindow | null
placeholder

export interface OpenAIAdditionalRateLimit {
  limit_name: string
  metered_feature: string
  rate_limit?: OpenAIRateLimit | null
placeholder

export interface OpenAIRateLimitResetCreditDetail {
  expires_at?: string
placeholder

export interface OpenAIRateLimitResetCredits {
  available_count: number
  credits?: OpenAIRateLimitResetCreditDetail[]
placeholder

export interface OpenAIQuotaUsage {
  user_id?: string
  account_id?: string
  email?: string
  plan_type?: string
  rate_limit?: OpenAIRateLimit | null
  additional_rate_limits?: OpenAIAdditionalRateLimit[]
  rate_limit_reset_credits?: OpenAIRateLimitResetCredits | null
  fetched_at: number
placeholder

export interface OpenAIQuotaResetCredit {
  id?: string
  reset_type?: string
  status?: string
  granted_at?: string
  expires_at?: string
  redeem_started_at?: string
  redeemed_at?: string
placeholder

export interface OpenAIQuotaResetResult {
  code: string
  credit?: OpenAIQuotaResetCredit | null
  windows_reset: number
  quota?: OpenAIQuotaUsage | null
  account?: Account | null
  cache_refreshed: boolean
  account_state_recovered: boolean
  warning_code?:
    | 'reset_credit_cache_refresh_failed'
    | 'account_state_recovery_failed'
    | 'account_state_refresh_failed'
placeholder

/** Usage payload plus whether the reset-credit snapshot was persisted. */
export interface OpenAIQuotaRefreshResult extends OpenAIQuotaUsage {
  cache_persisted: boolean
placeholder

/**
 * Query the upstream quota AND persist the reset-credit snapshot on the account
 * so the card can be rehydrated without an upstream round-trip. It is a POST
 * because it writes account state (and must therefore be audited).
 *
 * The read-only `GET /admin/openai/accounts/:id/quota` endpoint still exists for
 * API consumers; the panel always wants the snapshot persisted, so it has no
 * client binding here.
 */
export async function refreshOpenAIQuota(id: number): Promise<OpenAIQuotaRefreshResult> {
  const { data placeholder = await apiClient.post<OpenAIQuotaRefreshResult>(
    `/admin/openai/accounts/${idplaceholder/quota/refresh`
  )
  return data
placeholder

/**
 * Consume one rate-limit-reset credit for an OpenAI/Codex OAuth account.
 *
 * The credit is non-refundable and the endpoint chains an upstream reset with an
 * upstream re-query, so it needs a larger budget than the default client
 * timeout: aborting locally would report a successful consumption as a failure
 * and invite a retry that spends a second credit.
 */
export async function resetOpenAIQuota(id: number): Promise<OpenAIQuotaResetResult> {
  const { data placeholder = await apiClient.post<OpenAIQuotaResetResult>(
    `/admin/openai/accounts/${idplaceholder/reset-quota`,
    undefined,
    { timeout: 90_000 placeholder
  )
  return data
placeholder

export interface SparkShadowCreatePayload {
  name?: string
  priority?: number
  concurrency?: number
  group_ids?: number[]
placeholder

export async function createSparkShadow(parentId: number, payload: SparkShadowCreatePayload): Promise<Account> {
  const { data placeholder = await apiClient.post<Account>(`/admin/accounts/${parentIdplaceholder/shadow`, payload)
  return data
placeholder

export async function getUpstreamBillingProbeSettings(): Promise<UpstreamBillingProbeSettings> {
  const { data placeholder = await apiClient.get<UpstreamBillingProbeSettings>('/admin/accounts/upstream-billing-probe/settings')
  return data
placeholder

export async function updateUpstreamBillingProbeSettings(
  settings: UpstreamBillingProbeSettings
): Promise<UpstreamBillingProbeSettings> {
  const { data placeholder = await apiClient.put<UpstreamBillingProbeSettings>(
    '/admin/accounts/upstream-billing-probe/settings',
    settings
  )
  return data
placeholder

export async function setUpstreamBillingProbeEnabled(id: number, enabled: boolean): Promise<void> {
  await apiClient.put(`/admin/accounts/${idplaceholder/upstream-billing-probe`, { enabled placeholder)
placeholder

export async function probeUpstreamBilling(id: number): Promise<UpstreamBillingProbeResult> {
  const { data placeholder = await apiClient.post<UpstreamBillingProbeResult>(`/admin/accounts/${idplaceholder/upstream-billing-probe`)
  return data
placeholder

export async function probeUpstreamBillingBatch(accountIds: number[]): Promise<UpstreamBillingProbeResult[]> {
  const { data placeholder = await apiClient.post<{ results: UpstreamBillingProbeResult[] placeholder>(
    '/admin/accounts/upstream-billing-probe/batch',
    { account_ids: accountIds placeholder
  )
  return data.results
placeholder

export async function getOllamaCloudUsageSettings(): Promise<OllamaCloudUsageSettings> {
  const { data placeholder = await apiClient.get<OllamaCloudUsageSettings>('/admin/accounts/ollama-cloud-usage/settings')
  return data
placeholder

export async function updateOllamaCloudUsageSettings(
  settings: OllamaCloudUsageSettings
): Promise<OllamaCloudUsageSettings> {
  const { data placeholder = await apiClient.put<OllamaCloudUsageSettings>(
    '/admin/accounts/ollama-cloud-usage/settings',
    settings
  )
  return data
placeholder

export async function getOllamaCloudUsage(id: number): Promise<OllamaCloudUsageState> {
  const { data placeholder = await apiClient.get<OllamaCloudUsageState>(`/admin/accounts/${idplaceholder/ollama-cloud-usage`)
  return data
placeholder

export async function saveOllamaCloudUsageSession(id: number, session: string): Promise<OllamaCloudUsageState> {
  const { data placeholder = await apiClient.put<OllamaCloudUsageState>(`/admin/accounts/${idplaceholder/ollama-cloud-usage/session`, {
    session
  placeholder)
  return data
placeholder

export async function deleteOllamaCloudUsageSession(id: number): Promise<OllamaCloudUsageState> {
  const { data placeholder = await apiClient.delete<OllamaCloudUsageState>(`/admin/accounts/${idplaceholder/ollama-cloud-usage/session`)
  return data
placeholder

export async function setOllamaCloudUsageAutoRefresh(id: number, enabled: boolean): Promise<OllamaCloudUsageState> {
  const { data placeholder = await apiClient.put<OllamaCloudUsageState>(`/admin/accounts/${idplaceholder/ollama-cloud-usage/auto-refresh`, {
    enabled
  placeholder)
  return data
placeholder

export async function refreshOllamaCloudUsage(id: number): Promise<OllamaCloudUsageState> {
  const { data placeholder = await apiClient.post<OllamaCloudUsageState>(`/admin/accounts/${idplaceholder/ollama-cloud-usage/refresh`)
  return data
placeholder

export const accountsAPI = {
  list,
  listWithEtag,
  getById,
  create,
  duplicate,
  update,
  checkMixedChannelRisk,
  delete: deleteAccount,
  toggleStatus,
  testAccount,
  refreshCredentials,
  applyOAuthCredentials,
  getStats,
  clearError,
  getUsage,
  getBatchUsage,
  getTodayStats,
  getBatchTodayStats,
  clearRateLimit,
  recoverState,
  resetAccountQuota,
  getTempUnschedulableStatus,
  resetTempUnschedulable,
  setSchedulable,
  getAvailableModels,
  syncUpstreamModels,
  syncUpstreamModelsPreview,
  generateAuthUrl,
  exchangeCode,
  refreshOpenAIToken,
  batchCreate,
  batchUpdateCredentials,
  bulkUpdate,
  previewFromCrs,
  syncFromCrs,
  exportData,
  importData,
  importCodexSession,
  createOpenAICodexPAT,
  getAntigravityDefaultModelMapping,
  batchDelete,
  batchClearError,
  batchRefresh,
  setPrivacy,
  revertProxyFallback,
  refreshOpenAIQuota,
  resetOpenAIQuota,
  createSparkShadow,
  getUpstreamBillingProbeSettings,
  updateUpstreamBillingProbeSettings,
  setUpstreamBillingProbeEnabled,
  probeUpstreamBilling,
  probeUpstreamBillingBatch,
  getOllamaCloudUsageSettings,
  updateOllamaCloudUsageSettings,
  getOllamaCloudUsage,
  saveOllamaCloudUsageSession,
  deleteOllamaCloudUsageSession,
  setOllamaCloudUsageAutoRefresh,
  refreshOllamaCloudUsage
placeholder

export default accountsAPI
