/**
 * Admin Subscriptions API endpoints
 * Handles user subscription management for administrators
 */

import { apiClient placeholder from '../client'
import type {
  UserSubscription,
  SubscriptionProgress,
  AssignSubscriptionRequest,
  BulkAssignSubscriptionRequest,
  ExtendSubscriptionRequest,
  PaginatedResponse
placeholder from '@/types'

/**
 * List all subscriptions with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters (status, user_id, group_id, sort_by, sort_order)
 * @returns Paginated list of subscriptions
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    status?: 'active' | 'expired' | 'revoked'
    user_id?: number
    group_id?: number
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  placeholder,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<UserSubscription>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<UserSubscription>>(
    '/admin/subscriptions',
    {
      params: {
        page,
        page_size: pageSize,
        ...filters
      placeholder,
      signal: options?.signal
    placeholder
  )
  return data
placeholder

/**
 * Get subscription by ID
 * @param id - Subscription ID
 * @returns Subscription details
 */
export async function getById(id: number): Promise<UserSubscription> {
  const { data placeholder = await apiClient.get<UserSubscription>(`/admin/subscriptions/${idplaceholder`)
  return data
placeholder

/**
 * Get subscription progress
 * @param id - Subscription ID
 * @returns Subscription progress with usage stats
 */
export async function getProgress(id: number): Promise<SubscriptionProgress> {
  const { data placeholder = await apiClient.get<SubscriptionProgress>(`/admin/subscriptions/${idplaceholder/progress`)
  return data
placeholder

/**
 * Assign subscription to user
 * @param request - Assignment request
 * @returns Created subscription
 */
export async function assign(request: AssignSubscriptionRequest): Promise<UserSubscription> {
  const { data placeholder = await apiClient.post<UserSubscription>('/admin/subscriptions/assign', request)
  return data
placeholder

/**
 * Bulk assign subscriptions to multiple users
 * @param request - Bulk assignment request
 * @returns Created subscriptions
 */
export async function bulkAssign(
  request: BulkAssignSubscriptionRequest
): Promise<UserSubscription[]> {
  const { data placeholder = await apiClient.post<UserSubscription[]>(
    '/admin/subscriptions/bulk-assign',
    request
  )
  return data
placeholder

/**
 * Extend subscription validity
 * @param id - Subscription ID
 * @param request - Extension request with days
 * @returns Updated subscription
 */
export async function extend(
  id: number,
  request: ExtendSubscriptionRequest
): Promise<UserSubscription> {
  const { data placeholder = await apiClient.post<UserSubscription>(
    `/admin/subscriptions/${idplaceholder/extend`,
    request
  )
  return data
placeholder

/**
 * Revoke subscription
 * @param id - Subscription ID
 * @returns Success confirmation
 */
export async function revoke(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/subscriptions/${idplaceholder`)
  return data
placeholder

/**
 * Reset daily and/or weekly usage quota for a subscription
 * @param id - Subscription ID
 * @param options - Which windows to reset
 * @returns Updated subscription
 */
export async function resetQuota(
  id: number,
  options: { daily: boolean; weekly: boolean placeholder
): Promise<UserSubscription> {
  const { data placeholder = await apiClient.post<UserSubscription>(
    `/admin/subscriptions/${idplaceholder/reset-quota`,
    options
  )
  return data
placeholder

/**
 * List subscriptions by group
 * @param groupId - Group ID
 * @param page - Page number
 * @param pageSize - Items per page
 * @returns Paginated list of subscriptions in the group
 */
export async function listByGroup(
  groupId: number,
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<UserSubscription>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<UserSubscription>>(
    `/admin/groups/${groupIdplaceholder/subscriptions`,
    {
      params: { page, page_size: pageSize placeholder
    placeholder
  )
  return data
placeholder

/**
 * List subscriptions by user
 * @param userId - User ID
 * @param page - Page number
 * @param pageSize - Items per page
 * @returns Paginated list of user's subscriptions
 */
export async function listByUser(
  userId: number,
  page: number = 1,
  pageSize: number = 20
): Promise<PaginatedResponse<UserSubscription>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<UserSubscription>>(
    `/admin/users/${userIdplaceholder/subscriptions`,
    {
      params: { page, page_size: pageSize placeholder
    placeholder
  )
  return data
placeholder

export const subscriptionsAPI = {
  list,
  getById,
  getProgress,
  assign,
  bulkAssign,
  extend,
  revoke,
  resetQuota,
  listByGroup,
  listByUser
placeholder

export default subscriptionsAPI
