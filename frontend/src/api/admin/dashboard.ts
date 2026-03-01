/**
 * Admin Dashboard API endpoints
 * Provides system-wide statistics and metrics
 */

import { apiClient placeholder from '../client'
import type {
  DashboardStats,
  TrendDataPoint,
  ModelStat,
  GroupStat,
  ApiKeyUsageTrendPoint,
  UserUsageTrendPoint,
  UsageRequestType
placeholder from '@/types'

/**
 * Get dashboard statistics
 * @returns Dashboard statistics including users, keys, accounts, and token usage
 */
export async function getStats(): Promise<DashboardStats> {
  const { data placeholder = await apiClient.get<DashboardStats>('/admin/dashboard/stats')
  return data
placeholder

/**
 * Get real-time metrics
 * @returns Real-time system metrics
 */
export async function getRealtimeMetrics(): Promise<{
  active_requests: number
  requests_per_minute: number
  average_response_time: number
  error_rate: number
placeholder> {
  const { data placeholder = await apiClient.get<{
    active_requests: number
    requests_per_minute: number
    average_response_time: number
    error_rate: number
  placeholder>('/admin/dashboard/realtime')
  return data
placeholder

export interface TrendParams {
  start_date?: string
  end_date?: string
  granularity?: 'day' | 'hour'
  user_id?: number
  api_key_id?: number
  model?: string
  account_id?: number
  group_id?: number
  request_type?: UsageRequestType
  stream?: boolean
  billing_type?: number | null
placeholder

export interface TrendResponse {
  trend: TrendDataPoint[]
  start_date: string
  end_date: string
  granularity: string
placeholder

/**
 * Get usage trend data
 * @param params - Query parameters for filtering
 * @returns Usage trend data
 */
export async function getUsageTrend(params?: TrendParams): Promise<TrendResponse> {
  const { data placeholder = await apiClient.get<TrendResponse>('/admin/dashboard/trend', { params placeholder)
  return data
placeholder

export interface ModelStatsParams {
  start_date?: string
  end_date?: string
  user_id?: number
  api_key_id?: number
  model?: string
  account_id?: number
  group_id?: number
  request_type?: UsageRequestType
  stream?: boolean
  billing_type?: number | null
placeholder

export interface ModelStatsResponse {
  models: ModelStat[]
  start_date: string
  end_date: string
placeholder

/**
 * Get model usage statistics
 * @param params - Query parameters for filtering
 * @returns Model usage statistics
 */
export async function getModelStats(params?: ModelStatsParams): Promise<ModelStatsResponse> {
  const { data placeholder = await apiClient.get<ModelStatsResponse>('/admin/dashboard/models', { params placeholder)
  return data
placeholder

export interface GroupStatsParams {
  start_date?: string
  end_date?: string
  user_id?: number
  api_key_id?: number
  account_id?: number
  group_id?: number
  request_type?: UsageRequestType
  stream?: boolean
  billing_type?: number | null
placeholder

export interface GroupStatsResponse {
  groups: GroupStat[]
  start_date: string
  end_date: string
placeholder

/**
 * Get group usage statistics
 * @param params - Query parameters for filtering
 * @returns Group usage statistics
 */
export async function getGroupStats(params?: GroupStatsParams): Promise<GroupStatsResponse> {
  const { data placeholder = await apiClient.get<GroupStatsResponse>('/admin/dashboard/groups', { params placeholder)
  return data
placeholder

export interface ApiKeyTrendParams extends TrendParams {
  limit?: number
placeholder

export interface ApiKeyTrendResponse {
  trend: ApiKeyUsageTrendPoint[]
  start_date: string
  end_date: string
  granularity: string
placeholder

/**
 * Get API key usage trend data
 * @param params - Query parameters for filtering
 * @returns API key usage trend data
 */
export async function getApiKeyUsageTrend(
  params?: ApiKeyTrendParams
): Promise<ApiKeyTrendResponse> {
  const { data placeholder = await apiClient.get<ApiKeyTrendResponse>('/admin/dashboard/api-keys-trend', {
    params
  placeholder)
  return data
placeholder

export interface UserTrendParams extends TrendParams {
  limit?: number
placeholder

export interface UserTrendResponse {
  trend: UserUsageTrendPoint[]
  start_date: string
  end_date: string
  granularity: string
placeholder

/**
 * Get user usage trend data
 * @param params - Query parameters for filtering
 * @returns User usage trend data
 */
export async function getUserUsageTrend(params?: UserTrendParams): Promise<UserTrendResponse> {
  const { data placeholder = await apiClient.get<UserTrendResponse>('/admin/dashboard/users-trend', {
    params
  placeholder)
  return data
placeholder

export interface BatchUserUsageStats {
  user_id: number
  today_actual_cost: number
  total_actual_cost: number
placeholder

export interface BatchUsersUsageResponse {
  stats: Record<string, BatchUserUsageStats>
placeholder

/**
 * Get batch usage stats for multiple users
 * @param userIds - Array of user IDs
 * @returns Usage stats map keyed by user ID
 */
export async function getBatchUsersUsage(userIds: number[]): Promise<BatchUsersUsageResponse> {
  const { data placeholder = await apiClient.post<BatchUsersUsageResponse>('/admin/dashboard/users-usage', {
    user_ids: userIds
  placeholder)
  return data
placeholder

export interface BatchApiKeyUsageStats {
  api_key_id: number
  today_actual_cost: number
  total_actual_cost: number
placeholder

export interface BatchApiKeysUsageResponse {
  stats: Record<string, BatchApiKeyUsageStats>
placeholder

/**
 * Get batch usage stats for multiple API keys
 * @param apiKeyIds - Array of API key IDs
 * @returns Usage stats map keyed by API key ID
 */
export async function getBatchApiKeysUsage(
  apiKeyIds: number[]
): Promise<BatchApiKeysUsageResponse> {
  const { data placeholder = await apiClient.post<BatchApiKeysUsageResponse>(
    '/admin/dashboard/api-keys-usage',
    {
      api_key_ids: apiKeyIds
    placeholder
  )
  return data
placeholder

export const dashboardAPI = {
  getStats,
  getRealtimeMetrics,
  getUsageTrend,
  getModelStats,
  getGroupStats,
  getApiKeyUsageTrend,
  getUserUsageTrend,
  getBatchUsersUsage,
  getBatchApiKeysUsage
placeholder

export default dashboardAPI
