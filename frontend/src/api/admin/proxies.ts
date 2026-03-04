/**
 * Admin Proxies API endpoints
 * Handles proxy server management for administrators
 */

import { apiClient placeholder from '../client'
import type {
  Proxy,
  ProxyAccountSummary,
  ProxyQualityCheckResult,
  CreateProxyRequest,
  UpdateProxyRequest,
  PaginatedResponse,
  AdminDataPayload,
  AdminDataImportResult
placeholder from '@/types'

/**
 * List all proxies with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters
 * @returns Paginated list of proxies
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    protocol?: string
    status?: 'active' | 'inactive'
    search?: string
  placeholder,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<Proxy>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<Proxy>>('/admin/proxies', {
    params: {
      page,
      page_size: pageSize,
      ...filters
    placeholder,
    signal: options?.signal
  placeholder)
  return data
placeholder

/**
 * Get all active proxies (without pagination)
 * @returns List of all active proxies
 */
export async function getAll(): Promise<Proxy[]> {
  const { data placeholder = await apiClient.get<Proxy[]>('/admin/proxies/all')
  return data
placeholder

/**
 * Get all active proxies with account count (sorted by creation time desc)
 * @returns List of all active proxies with account count
 */
export async function getAllWithCount(): Promise<Proxy[]> {
  const { data placeholder = await apiClient.get<Proxy[]>('/admin/proxies/all', {
    params: { with_count: 'true' placeholder
  placeholder)
  return data
placeholder

/**
 * Get proxy by ID
 * @param id - Proxy ID
 * @returns Proxy details
 */
export async function getById(id: number): Promise<Proxy> {
  const { data placeholder = await apiClient.get<Proxy>(`/admin/proxies/${idplaceholder`)
  return data
placeholder

/**
 * Create new proxy
 * @param proxyData - Proxy data
 * @returns Created proxy
 */
export async function create(proxyData: CreateProxyRequest): Promise<Proxy> {
  const { data placeholder = await apiClient.post<Proxy>('/admin/proxies', proxyData)
  return data
placeholder

/**
 * Update proxy
 * @param id - Proxy ID
 * @param updates - Fields to update
 * @returns Updated proxy
 */
export async function update(id: number, updates: UpdateProxyRequest): Promise<Proxy> {
  const { data placeholder = await apiClient.put<Proxy>(`/admin/proxies/${idplaceholder`, updates)
  return data
placeholder

/**
 * Delete proxy
 * @param id - Proxy ID
 * @returns Success confirmation
 */
export async function deleteProxy(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/proxies/${idplaceholder`)
  return data
placeholder

/**
 * Toggle proxy status
 * @param id - Proxy ID
 * @param status - New status
 * @returns Updated proxy
 */
export async function toggleStatus(id: number, status: 'active' | 'inactive'): Promise<Proxy> {
  return update(id, { status placeholder)
placeholder

/**
 * Test proxy connectivity
 * @param id - Proxy ID
 * @returns Test result with IP info
 */
export async function testProxy(id: number): Promise<{
  success: boolean
  message: string
  latency_ms?: number
  ip_address?: string
  city?: string
  region?: string
  country?: string
  country_code?: string
placeholder> {
  const { data placeholder = await apiClient.post<{
    success: boolean
    message: string
    latency_ms?: number
    ip_address?: string
    city?: string
    region?: string
    country?: string
    country_code?: string
  placeholder>(`/admin/proxies/${idplaceholder/test`)
  return data
placeholder

/**
 * Check proxy quality across common AI targets
 * @param id - Proxy ID
 * @returns Quality check result
 */
export async function checkProxyQuality(id: number): Promise<ProxyQualityCheckResult> {
  const { data placeholder = await apiClient.post<ProxyQualityCheckResult>(`/admin/proxies/${idplaceholder/quality-check`)
  return data
placeholder

/**
 * Get proxy usage statistics
 * @param id - Proxy ID
 * @returns Proxy usage statistics
 */
export async function getStats(id: number): Promise<{
  total_accounts: number
  active_accounts: number
  total_requests: number
  success_rate: number
  average_latency: number
placeholder> {
  const { data placeholder = await apiClient.get<{
    total_accounts: number
    active_accounts: number
    total_requests: number
    success_rate: number
    average_latency: number
  placeholder>(`/admin/proxies/${idplaceholder/stats`)
  return data
placeholder

/**
 * Get accounts using a proxy
 * @param id - Proxy ID
 * @returns List of accounts using the proxy
 */
export async function getProxyAccounts(id: number): Promise<ProxyAccountSummary[]> {
  const { data placeholder = await apiClient.get<ProxyAccountSummary[]>(`/admin/proxies/${idplaceholder/accounts`)
  return data
placeholder

/**
 * Batch create proxies
 * @param proxies - Array of proxy data to create
 * @returns Creation result with count of created and skipped
 */
export async function batchCreate(
  proxies: Array<{
    protocol: string
    host: string
    port: number
    username?: string
    password?: string
  placeholder>
): Promise<{
  created: number
  skipped: number
placeholder> {
  const { data placeholder = await apiClient.post<{
    created: number
    skipped: number
  placeholder>('/admin/proxies/batch', { proxies placeholder)
  return data
placeholder

export async function batchDelete(ids: number[]): Promise<{
  deleted_ids: number[]
  skipped: Array<{ id: number; reason: string placeholder>
placeholder> {
  const { data placeholder = await apiClient.post<{
    deleted_ids: number[]
    skipped: Array<{ id: number; reason: string placeholder>
  placeholder>('/admin/proxies/batch-delete', { ids placeholder)
  return data
placeholder

export async function exportData(options?: {
  ids?: number[]
  filters?: {
    protocol?: string
    status?: 'active' | 'inactive'
    search?: string
  placeholder
placeholder): Promise<AdminDataPayload> {
  const params: Record<string, string> = {placeholder
  if (options?.ids && options.ids.length > 0) {
    params.ids = options.ids.join(',')
  placeholder else if (options?.filters) {
    const { protocol, status, search placeholder = options.filters
    if (protocol) params.protocol = protocol
    if (status) params.status = status
    if (search) params.search = search
  placeholder
  const { data placeholder = await apiClient.get<AdminDataPayload>('/admin/proxies/data', { params placeholder)
  return data
placeholder

export async function importData(payload: {
  data: AdminDataPayload
placeholder): Promise<AdminDataImportResult> {
  const { data placeholder = await apiClient.post<AdminDataImportResult>('/admin/proxies/data', payload)
  return data
placeholder

export const proxiesAPI = {
  list,
  getAll,
  getAllWithCount,
  getById,
  create,
  update,
  delete: deleteProxy,
  toggleStatus,
  testProxy,
  checkProxyQuality,
  getStats,
  getProxyAccounts,
  batchCreate,
  batchDelete,
  exportData,
  importData
placeholder

export default proxiesAPI
