/**
 * Admin Usage API endpoints
 * Handles admin-level usage logs and statistics retrieval
 */

import { apiClient placeholder from '../client'
import type { UsageLog, UsageQueryParams, PaginatedResponse placeholder from '@/types'

// ==================== Types ====================

export interface AdminUsageStatsResponse {
  total_requests: number
  total_input_tokens: number
  total_output_tokens: number
  total_cache_tokens: number
  total_tokens: number
  total_cost: number
  total_actual_cost: number
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
): Promise<PaginatedResponse<UsageLog>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<UsageLog>>('/admin/usage', {
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
  billing_type?: number
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

export const adminUsageAPI = {
  list,
  getStats,
  searchUsers,
  searchApiKeys
placeholder

export default adminUsageAPI
