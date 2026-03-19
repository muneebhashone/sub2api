/**
 * Admin Groups API endpoints
 * Handles API key group management for administrators
 */

import { apiClient placeholder from '../client'
import type {
  AdminGroup,
  GroupPlatform,
  CreateGroupRequest,
  UpdateGroupRequest,
  PaginatedResponse
placeholder from '@/types'

/**
 * List all groups with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters (platform, status, is_exclusive, search)
 * @returns Paginated list of groups
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    platform?: GroupPlatform
    status?: 'active' | 'inactive'
    is_exclusive?: boolean
    search?: string
  placeholder,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<AdminGroup>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<AdminGroup>>('/admin/groups', {
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
 * Get all active groups (without pagination)
 * @param platform - Optional platform filter
 * @returns List of all active groups
 */
export async function getAll(platform?: GroupPlatform): Promise<AdminGroup[]> {
  const { data placeholder = await apiClient.get<AdminGroup[]>('/admin/groups/all', {
    params: platform ? { platform placeholder : undefined
  placeholder)
  return data
placeholder

/**
 * Get active groups by platform
 * @param platform - Platform to filter by
 * @returns List of groups for the specified platform
 */
export async function getByPlatform(platform: GroupPlatform): Promise<AdminGroup[]> {
  return getAll(platform)
placeholder

/**
 * Get group by ID
 * @param id - Group ID
 * @returns Group details
 */
export async function getById(id: number): Promise<AdminGroup> {
  const { data placeholder = await apiClient.get<AdminGroup>(`/admin/groups/${idplaceholder`)
  return data
placeholder

/**
 * Create new group
 * @param groupData - Group data
 * @returns Created group
 */
export async function create(groupData: CreateGroupRequest): Promise<AdminGroup> {
  const { data placeholder = await apiClient.post<AdminGroup>('/admin/groups', groupData)
  return data
placeholder

/**
 * Update group
 * @param id - Group ID
 * @param updates - Fields to update
 * @returns Updated group
 */
export async function update(id: number, updates: UpdateGroupRequest): Promise<AdminGroup> {
  const { data placeholder = await apiClient.put<AdminGroup>(`/admin/groups/${idplaceholder`, updates)
  return data
placeholder

/**
 * Delete group
 * @param id - Group ID
 * @returns Success confirmation
 */
export async function deleteGroup(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/groups/${idplaceholder`)
  return data
placeholder

/**
 * Toggle group status
 * @param id - Group ID
 * @param status - New status
 * @returns Updated group
 */
export async function toggleStatus(id: number, status: 'active' | 'inactive'): Promise<AdminGroup> {
  return update(id, { status placeholder)
placeholder

/**
 * Get group statistics
 * @param id - Group ID
 * @returns Group usage statistics
 */
export async function getStats(id: number): Promise<{
  total_api_keys: number
  active_api_keys: number
  total_requests: number
  total_cost: number
placeholder> {
  const { data placeholder = await apiClient.get<{
    total_api_keys: number
    active_api_keys: number
    total_requests: number
    total_cost: number
  placeholder>(`/admin/groups/${idplaceholder/stats`)
  return data
placeholder

/**
 * Get API keys in a group
 * @param id - Group ID
 * @param page - Page number
 * @param pageSize - Items per page
 * @returns Paginated list of API keys in the group
 */
export async function getGroupApiKeys(
  id: number,
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<any>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<any>>(`/admin/groups/${idplaceholder/api-keys`, {
    params: { page, page_size: pageSize placeholder
  placeholder)
  return data
placeholder

/**
 * Rate multiplier entry for a user in a group
 */
export interface GroupRateMultiplierEntry {
  user_id: number
  user_name: string
  user_email: string
  user_notes: string
  user_status: string
  rate_multiplier: number
placeholder

/**
 * Get rate multipliers for users in a group
 * @param id - Group ID
 * @returns List of user rate multiplier entries
 */
export async function getGroupRateMultipliers(id: number): Promise<GroupRateMultiplierEntry[]> {
  const { data placeholder = await apiClient.get<GroupRateMultiplierEntry[]>(
    `/admin/groups/${idplaceholder/rate-multipliers`
  )
  return data
placeholder

/**
 * Update group sort orders
 * @param updates - Array of { id, sort_order placeholder objects
 * @returns Success confirmation
 */
export async function updateSortOrder(
  updates: Array<{ id: number; sort_order: number placeholder>
): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.put<{ message: string placeholder>('/admin/groups/sort-order', {
    updates
  placeholder)
  return data
placeholder

/**
 * Clear all rate multipliers for a group
 * @param id - Group ID
 * @returns Success confirmation
 */
export async function clearGroupRateMultipliers(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/groups/${idplaceholder/rate-multipliers`)
  return data
placeholder

/**
 * Batch set rate multipliers for users in a group
 * @param id - Group ID
 * @param entries - Array of { user_id, rate_multiplier placeholder
 * @returns Success confirmation
 */
export async function batchSetGroupRateMultipliers(
  id: number,
  entries: Array<{ user_id: number; rate_multiplier: number placeholder>
): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.put<{ message: string placeholder>(
    `/admin/groups/${idplaceholder/rate-multipliers`,
    { entries placeholder
  )
  return data
placeholder

/**
 * Get usage summary (today + cumulative cost) for all groups
 * @param timezone - IANA timezone string (e.g. "Asia/Shanghai")
 * @returns Array of group usage summaries
 */
export async function getUsageSummary(
  timezone?: string
): Promise<{ group_id: number; today_cost: number; total_cost: number placeholder[]> {
  const { data placeholder = await apiClient.get<
    { group_id: number; today_cost: number; total_cost: number placeholder[]
  >('/admin/groups/usage-summary', {
    params: timezone ? { timezone placeholder : undefined
  placeholder)
  return data
placeholder

/**
 * Get capacity summary (concurrency/sessions/RPM) for all active groups
 */
export async function getCapacitySummary(): Promise<
  { group_id: number; concurrency_used: number; concurrency_max: number; sessions_used: number; sessions_max: number; rpm_used: number; rpm_max: number placeholder[]
> {
  const { data placeholder = await apiClient.get<
    { group_id: number; concurrency_used: number; concurrency_max: number; sessions_used: number; sessions_max: number; rpm_used: number; rpm_max: number placeholder[]
  >('/admin/groups/capacity-summary')
  return data
placeholder

export const groupsAPI = {
  list,
  getAll,
  getByPlatform,
  getById,
  create,
  update,
  delete: deleteGroup,
  toggleStatus,
  getStats,
  getGroupApiKeys,
  getGroupRateMultipliers,
  clearGroupRateMultipliers,
  batchSetGroupRateMultipliers,
  updateSortOrder,
  getUsageSummary,
  getCapacitySummary
placeholder

export default groupsAPI
