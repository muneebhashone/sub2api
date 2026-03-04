/**
 * Admin Users API endpoints
 * Handles user management for administrators
 */

import { apiClient placeholder from '../client'
import type { AdminUser, UpdateUserRequest, PaginatedResponse, ApiKey placeholder from '@/types'

/**
 * List all users with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters (status, role, search, attributes)
 * @param options - Optional request options (signal)
 * @returns Paginated list of users
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    status?: 'active' | 'disabled'
    role?: 'admin' | 'user'
    search?: string
    attributes?: Record<number, string>  // attributeId -> value
    include_subscriptions?: boolean
  placeholder,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<AdminUser>> {
  // Build params with attribute filters in attr[id]=value format
  const params: Record<string, any> = {
    page,
    page_size: pageSize,
    status: filters?.status,
    role: filters?.role,
    search: filters?.search,
    include_subscriptions: filters?.include_subscriptions
  placeholder

  // Add attribute filters as attr[id]=value
  if (filters?.attributes) {
    for (const [attrId, value] of Object.entries(filters.attributes)) {
      if (value) {
        params[`attr[${attrIdplaceholder]`] = value
      placeholder
    placeholder
  placeholder
  const { data placeholder = await apiClient.get<PaginatedResponse<AdminUser>>('/admin/users', {
    params,
    signal: options?.signal
  placeholder)
  return data
placeholder

/**
 * Get user by ID
 * @param id - User ID
 * @returns User details
 */
export async function getById(id: number): Promise<AdminUser> {
  const { data placeholder = await apiClient.get<AdminUser>(`/admin/users/${idplaceholder`)
  return data
placeholder

/**
 * Create new user
 * @param userData - User data (email, password, etc.)
 * @returns Created user
 */
export async function create(userData: {
  email: string
  password: string
  balance?: number
  concurrency?: number
  allowed_groups?: number[] | null
placeholder): Promise<AdminUser> {
  const { data placeholder = await apiClient.post<AdminUser>('/admin/users', userData)
  return data
placeholder

/**
 * Update user
 * @param id - User ID
 * @param updates - Fields to update
 * @returns Updated user
 */
export async function update(id: number, updates: UpdateUserRequest): Promise<AdminUser> {
  const { data placeholder = await apiClient.put<AdminUser>(`/admin/users/${idplaceholder`, updates)
  return data
placeholder

/**
 * Delete user
 * @param id - User ID
 * @returns Success confirmation
 */
export async function deleteUser(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/users/${idplaceholder`)
  return data
placeholder

/**
 * Update user balance
 * @param id - User ID
 * @param balance - New balance
 * @param operation - Operation type ('set', 'add', 'subtract')
 * @param notes - Optional notes for the balance adjustment
 * @returns Updated user
 */
export async function updateBalance(
  id: number,
  balance: number,
  operation: 'set' | 'add' | 'subtract' = 'set',
  notes?: string
): Promise<AdminUser> {
  const { data placeholder = await apiClient.post<AdminUser>(`/admin/users/${idplaceholder/balance`, {
    balance,
    operation,
    notes: notes || ''
  placeholder)
  return data
placeholder

/**
 * Update user concurrency
 * @param id - User ID
 * @param concurrency - New concurrency limit
 * @returns Updated user
 */
export async function updateConcurrency(id: number, concurrency: number): Promise<AdminUser> {
  return update(id, { concurrency placeholder)
placeholder

/**
 * Toggle user status
 * @param id - User ID
 * @param status - New status
 * @returns Updated user
 */
export async function toggleStatus(id: number, status: 'active' | 'disabled'): Promise<AdminUser> {
  return update(id, { status placeholder)
placeholder

/**
 * Get user's API keys
 * @param id - User ID
 * @returns List of user's API keys
 */
export async function getUserApiKeys(id: number): Promise<PaginatedResponse<ApiKey>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<ApiKey>>(`/admin/users/${idplaceholder/api-keys`)
  return data
placeholder

/**
 * Get user's usage statistics
 * @param id - User ID
 * @param period - Time period
 * @returns User usage statistics
 */
export async function getUserUsageStats(
  id: number,
  period: string = 'month'
): Promise<{
  total_requests: number
  total_cost: number
  total_tokens: number
placeholder> {
  const { data placeholder = await apiClient.get<{
    total_requests: number
    total_cost: number
    total_tokens: number
  placeholder>(`/admin/users/${idplaceholder/usage`, {
    params: { period placeholder
  placeholder)
  return data
placeholder

/**
 * Balance history item returned from the API
 */
export interface BalanceHistoryItem {
  id: number
  code: string
  type: string
  value: number
  status: string
  used_by: number | null
  used_at: string | null
  created_at: string
  group_id: number | null
  validity_days: number
  notes: string
  user?: { id: number; email: string placeholder | null
  group?: { id: number; name: string placeholder | null
placeholder

// Balance history response extends pagination with total_recharged summary
export interface BalanceHistoryResponse extends PaginatedResponse<BalanceHistoryItem> {
  total_recharged: number
placeholder

/**
 * Get user's balance/concurrency change history
 * @param id - User ID
 * @param page - Page number
 * @param pageSize - Items per page
 * @param type - Optional type filter (balance, admin_balance, concurrency, admin_concurrency, subscription)
 * @returns Paginated balance history with total_recharged
 */
export async function getUserBalanceHistory(
  id: number,
  page: number = 1,
  pageSize: number = 20,
  type?: string
): Promise<BalanceHistoryResponse> {
  const params: Record<string, any> = { page, page_size: pageSize placeholder
  if (type) params.type = type
  const { data placeholder = await apiClient.get<BalanceHistoryResponse>(
    `/admin/users/${idplaceholder/balance-history`,
    { params placeholder
  )
  return data
placeholder

export const usersAPI = {
  list,
  getById,
  create,
  update,
  delete: deleteUser,
  updateBalance,
  updateConcurrency,
  toggleStatus,
  getUserApiKeys,
  getUserUsageStats,
  getUserBalanceHistory
placeholder

export default usersAPI
