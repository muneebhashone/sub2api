/**
 * API Keys management endpoints
 * Handles CRUD operations for user API keys
 */

import { apiClient placeholder from './client'
import type { ApiKey, CreateApiKeyRequest, UpdateApiKeyRequest, PaginatedResponse placeholder from '@/types'

/**
 * List all API keys for current user
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 10)
 * @param options - Optional request options
 * @returns Paginated list of API keys
 */
export async function list(
  page: number = 1,
  pageSize: number = 10,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<ApiKey>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<ApiKey>>('/keys', {
    params: { page, page_size: pageSize placeholder,
    signal: options?.signal
  placeholder)
  return data
placeholder

/**
 * Get API key by ID
 * @param id - API key ID
 * @returns API key details
 */
export async function getById(id: number): Promise<ApiKey> {
  const { data placeholder = await apiClient.get<ApiKey>(`/keys/${idplaceholder`)
  return data
placeholder

/**
 * Create new API key
 * @param name - Key name
 * @param groupId - Optional group ID
 * @param customKey - Optional custom key value
 * @param ipWhitelist - Optional IP whitelist
 * @param ipBlacklist - Optional IP blacklist
 * @param quota - Optional quota limit in USD (0 = unlimited)
 * @param expiresInDays - Optional days until expiry (undefined = never expires)
 * @returns Created API key
 */
export async function create(
  name: string,
  groupId?: number | null,
  customKey?: string,
  ipWhitelist?: string[],
  ipBlacklist?: string[],
  quota?: number,
  expiresInDays?: number
): Promise<ApiKey> {
  const payload: CreateApiKeyRequest = { name placeholder
  if (groupId !== undefined) {
    payload.group_id = groupId
  placeholder
  if (customKey) {
    payload.custom_key = customKey
  placeholder
  if (ipWhitelist && ipWhitelist.length > 0) {
    payload.ip_whitelist = ipWhitelist
  placeholder
  if (ipBlacklist && ipBlacklist.length > 0) {
    payload.ip_blacklist = ipBlacklist
  placeholder
  if (quota !== undefined && quota > 0) {
    payload.quota = quota
  placeholder
  if (expiresInDays !== undefined && expiresInDays > 0) {
    payload.expires_in_days = expiresInDays
  placeholder

  const { data placeholder = await apiClient.post<ApiKey>('/keys', payload)
  return data
placeholder

/**
 * Update API key
 * @param id - API key ID
 * @param updates - Fields to update
 * @returns Updated API key
 */
export async function update(id: number, updates: UpdateApiKeyRequest): Promise<ApiKey> {
  const { data placeholder = await apiClient.put<ApiKey>(`/keys/${idplaceholder`, updates)
  return data
placeholder

/**
 * Delete API key
 * @param id - API key ID
 * @returns Success confirmation
 */
export async function deleteKey(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/keys/${idplaceholder`)
  return data
placeholder

/**
 * Toggle API key status (active/inactive)
 * @param id - API key ID
 * @param status - New status
 * @returns Updated API key
 */
export async function toggleStatus(id: number, status: 'active' | 'inactive'): Promise<ApiKey> {
  return update(id, { status placeholder)
placeholder

export const keysAPI = {
  list,
  getById,
  create,
  update,
  delete: deleteKey,
  toggleStatus
placeholder

export default keysAPI
