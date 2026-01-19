/**
 * Admin Usage API endpoints
 * Handles admin-level usage logs and statistics retrieval
 */

import { apiClient placeholder from '../client'
import type { AdminUsageLog, UsageQueryParams, PaginatedResponse placeholder from '@/types'

// ==================== Types ====================

export interface AdminUsageStatsResponse {
  total_requests: number
  total_input_tokens: number
  total_output_tokens: number
  total_cache_tokens: number
  total_tokens: number
  total_cost: number
  total_actual_cost: number
  total_account_cost?: number
  average_duration_ms: number
placeholder

export interface SimpleUser {
  id: number
  email: string
placeholder

export interface SimpleApiKey {
  id: number
  name: string
  user_id: number
placeholder

export interface UsageCleanupFilters {
  start_time: string
  end_time: string
  user_id?: number
  api_key_id?: number
  account_id?: number
  group_id?: number
  model?: string | null
  stream?: boolean | null
  billing_type?: number | null
placeholder

export interface UsageCleanupTask {
  id: number
  status: string
  filters: UsageCleanupFilters
  created_by: number
  deleted_rows: number
  error_message?: string | null
  canceled_by?: number | null
  canceled_at?: string | null
  started_at?: string | null
  finished_at?: string | null
  created_at: string
  updated_at: string
placeholder

export interface CreateUsageCleanupTaskRequest {
  start_date: string
  end_date: string
  user_id?: number
  api_key_id?: number
  account_id?: number
  group_id?: number
  model?: string | null
  stream?: boolean | null
  billing_type?: number | null
  timezone?: string
placeholder

export interface AdminUsageQueryParams extends UsageQueryParams {
  user_id?: number
placeholder

// ==================== API Functions ====================

/**
 * List all usage logs with optional filters (admin only)
 * @param params - Query parameters for filtering and pagination
 * @returns Paginated list of usage logs
 */
export async function list(
  params: AdminUsageQueryParams,
  options?: { signal?: AbortSignal placeholder
): Promise<PaginatedResponse<AdminUsageLog>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<AdminUsageLog>>('/admin/usage', {
    params,
    signal: options?.signal
  placeholder)
  return data
placeholder

/**
 * Get usage statistics with optional filters (admin only)
 * @param params - Query parameters for filtering
 * @returns Usage statistics
 */
export async function getStats(params: {
  user_id?: number
  api_key_id?: number
  account_id?: number
  group_id?: number
  model?: string
  stream?: boolean
  period?: string
  start_date?: string
  end_date?: string
  timezone?: string
placeholder): Promise<AdminUsageStatsResponse> {
  const { data placeholder = await apiClient.get<AdminUsageStatsResponse>('/admin/usage/stats', {
    params
  placeholder)
  return data
placeholder

/**
 * Search users by email keyword (admin only)
 * @param keyword - Email keyword to search
 * @returns List of matching users (max 30)
 */
export async function searchUsers(keyword: string): Promise<SimpleUser[]> {
  const { data placeholder = await apiClient.get<SimpleUser[]>('/admin/usage/search-users', {
    params: { q: keyword placeholder
  placeholder)
  return data
placeholder

/**
 * Search API keys by user ID and/or keyword (admin only)
 * @param userId - Optional user ID to filter by
 * @param keyword - Optional keyword to search in key name
 * @returns List of matching API keys (max 30)
 */
export async function searchApiKeys(userId?: number, keyword?: string): Promise<SimpleApiKey[]> {
  const params: Record<string, unknown> = {placeholder
  if (userId !== undefined) {
    params.user_id = userId
  placeholder
  if (keyword) {
    params.q = keyword
  placeholder
  const { data placeholder = await apiClient.get<SimpleApiKey[]>('/admin/usage/search-api-keys', {
    params
  placeholder)
  return data
placeholder

/**
 * List usage cleanup tasks (admin only)
 * @param params - Query parameters for pagination
 * @returns Paginated list of cleanup tasks
 */
export async function listCleanupTasks(
  params: { page?: number; page_size?: number placeholder,
  options?: { signal?: AbortSignal placeholder
): Promise<PaginatedResponse<UsageCleanupTask>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<UsageCleanupTask>>('/admin/usage/cleanup-tasks', {
    params,
    signal: options?.signal
  placeholder)
  return data
placeholder

/**
 * Create a usage cleanup task (admin only)
 * @param payload - Cleanup task parameters
 * @returns Created cleanup task
 */
export async function createCleanupTask(payload: CreateUsageCleanupTaskRequest): Promise<UsageCleanupTask> {
  const { data placeholder = await apiClient.post<UsageCleanupTask>('/admin/usage/cleanup-tasks', payload)
  return data
placeholder

/**
 * Cancel a usage cleanup task (admin only)
 * @param taskId - Task ID to cancel
 */
export async function cancelCleanupTask(taskId: number): Promise<{ id: number; status: string placeholder> {
  const { data placeholder = await apiClient.post<{ id: number; status: string placeholder>(
    `/admin/usage/cleanup-tasks/${taskIdplaceholder/cancel`
  )
  return data
placeholder

export const adminUsageAPI = {
  list,
  getStats,
  searchUsers,
  searchApiKeys,
  listCleanupTasks,
  createCleanupTask,
  cancelCleanupTask
placeholder

export default adminUsageAPI
