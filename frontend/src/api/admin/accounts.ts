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
  CheckMixedChannelRequest,
  CheckMixedChannelResponse
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
    search?: string
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
export async function getUsage(id: number): Promise<AccountUsageInfo> {
  const { data placeholder = await apiClient.get<AccountUsageInfo>(`/admin/accounts/${idplaceholder/usage`)
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
  accountIds: number[],
  updates: Record<string, unknown>
): Promise<{
  success: number
  failed: number
  success_ids?: number[]
  failed_ids?: number[]
  results: Array<{ account_id: number; success: boolean; error?: string placeholder>
  placeholder> {
  const { data placeholder = await apiClient.post<{
    success: number
    failed: number
    success_ids?: number[]
    failed_ids?: number[]
    results: Array<{ account_id: number; success: boolean; error?: string placeholder>
  placeholder>('/admin/accounts/bulk-update', {
    account_ids: accountIds,
    ...updates
  placeholder)
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
  placeholder>('/admin/accounts/sync/crs', params)
  return data
placeholder

export async function exportData(options?: {
  ids?: number[]
  filters?: {
    platform?: string
    type?: string
    status?: string
    search?: string
  placeholder
  includeProxies?: boolean
placeholder): Promise<AdminDataPayload> {
  const params: Record<string, string> = {placeholder
  if (options?.ids && options.ids.length > 0) {
    params.ids = options.ids.join(',')
  placeholder else if (options?.filters) {
    const { platform, type, status, search placeholder = options.filters
    if (platform) params.platform = platform
    if (type) params.type = type
    if (status) params.status = status
    if (search) params.search = search
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
  endpoint: string = '/admin/openai/refresh-token'
): Promise<Record<string, unknown>> {
  const payload: { refresh_token: string; proxy_id?: number placeholder = {
    refresh_token: refreshToken
  placeholder
  if (proxyId) {
    payload.proxy_id = proxyId
  placeholder
  const { data placeholder = await apiClient.post<Record<string, unknown>>(endpoint, payload)
  return data
placeholder

/**
 * Validate Sora session token and exchange to access token
 * @param sessionToken - Sora session token
 * @param proxyId - Optional proxy ID
 * @param endpoint - API endpoint path
 * @returns Token information including access_token
 */
export async function validateSoraSessionToken(
  sessionToken: string,
  proxyId?: number | null,
  endpoint: string = '/admin/sora/st2at'
): Promise<Record<string, unknown>> {
  const payload: { session_token: string; proxy_id?: number placeholder = {
    session_token: sessionToken
  placeholder
  if (proxyId) {
    payload.proxy_id = proxyId
  placeholder
  const { data placeholder = await apiClient.post<Record<string, unknown>>(endpoint, payload)
  return data
placeholder

export const accountsAPI = {
  list,
  listWithEtag,
  getById,
  create,
  update,
  checkMixedChannelRisk,
  delete: deleteAccount,
  toggleStatus,
  testAccount,
  refreshCredentials,
  getStats,
  clearError,
  getUsage,
  getTodayStats,
  clearRateLimit,
  getTempUnschedulableStatus,
  resetTempUnschedulable,
  setSchedulable,
  getAvailableModels,
  generateAuthUrl,
  exchangeCode,
  refreshOpenAIToken,
  validateSoraSessionToken,
  batchCreate,
  batchUpdateCredentials,
  bulkUpdate,
  previewFromCrs,
  syncFromCrs,
  exportData,
  importData,
  getAntigravityDefaultModelMapping
placeholder

export default accountsAPI
