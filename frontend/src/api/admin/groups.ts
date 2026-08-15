/**
 * Admin Groups API endpoints
 * Handles API key group management for administrators
 */

import { apiClient placeholder from '../client'
import type {
  AdminGroup,
  GroupPlatform,
  CompositeModelRoute,
  CompositeModelRouteInput,
  CompositeRoutePreviewRequest,
  CompositeRouteDecision,
  CreateGroupRequest,
  UpdateGroupRequest,
  PaginatedResponse
placeholder from '@/types'

export interface LiveCapability {
  supported: boolean
  reason?: string
placeholder

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
    sort_by?: string
    sort_order?: 'asc' | 'desc'
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
 * Get ALL groups including disabled ones — used by the API Key group filter so
 * that admins can filter users whose keys are still bound to a now-disabled group.
 */
export async function getAllIncludingInactive(): Promise<AdminGroup[]> {
  const { data placeholder = await apiClient.get<AdminGroup[]>('/admin/groups/all', {
    params: { include_inactive: true placeholder
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

/** 获取当前 Sub2API 服务端的 Live 运行环境能力。 */
export async function getLiveCapability(): Promise<LiveCapability> {
  const { data placeholder = await apiClient.get<LiveCapability>('/admin/groups/live-capability')
  return data
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
 * Get candidate models for custom /v1/models list.
 * id=0 returns platform default models for create flow.
 */
export async function getModelsListCandidates(
  id: number,
  platform?: GroupPlatform
): Promise<string[]> {
  const { data placeholder = await apiClient.get<{ models: string[] placeholder>(
    `/admin/groups/${idplaceholder/models-list-candidates`,
    {
      params: platform ? { platform placeholder : undefined
    placeholder
  )
  return data.models || []
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
 * Duplicate a group on the server so configuration that is not present in the
 * list response is preserved. Keep the operation key after ambiguous failures
 * so a retry replays the original operation instead of creating another group.
 */
const duplicateOperationKeys = new Map<string, string>()

interface DuplicateOperationScope {
  adminID: string
  key: string
placeholder

function getCurrentAdminID(): string | null {
  try {
    const rawUser = globalThis.localStorage?.getItem('auth_user')
    if (!rawUser) return null

    const user: unknown = JSON.parse(rawUser)
    if (typeof user !== 'object' || user === null) return null

    const id = (user as { id?: unknown placeholder).id
    if (typeof id !== 'number' || !Number.isSafeInteger(id) || id <= 0) return null
    return String(id)
  placeholder catch {
    return null
  placeholder
placeholder

function duplicateOperationScope(id: number): DuplicateOperationScope | null {
  const adminID = getCurrentAdminID()
  if (!adminID) return null

  return {
    adminID,
    key: `sub2api:admin:group-duplicate:${adminIDplaceholder:${idplaceholder`
  placeholder
placeholder

function getStoredDuplicateOperationKey(storageKey: string): string | null {
  try {
    return globalThis.sessionStorage?.getItem(storageKey) ?? null
  placeholder catch {
    return null
  placeholder
placeholder

function storeDuplicateOperationKey(storageKey: string, key: string | null): void {
  try {
    if (key) globalThis.sessionStorage?.setItem(storageKey, key)
    else globalThis.sessionStorage?.removeItem(storageKey)
  placeholder catch {
    // In-memory retry protection still works when browser storage is unavailable.
  placeholder
placeholder

export async function duplicate(id: number): Promise<AdminGroup> {
  const scope = duplicateOperationScope(id)
  let idempotencyKey = scope
    ? duplicateOperationKeys.get(scope.key) ?? getStoredDuplicateOperationKey(scope.key)
    : null
  if (!idempotencyKey) {
    const requestID = globalThis.crypto?.randomUUID?.() ?? `${Date.now()placeholder-${Math.random().toString(36).slice(2)placeholder`
    idempotencyKey = `group-duplicate-${scope?.adminID ?? 'unknown-admin'placeholder-${idplaceholder-${requestIDplaceholder`
  placeholder
  if (scope) {
    duplicateOperationKeys.set(scope.key, idempotencyKey)
    storeDuplicateOperationKey(scope.key, idempotencyKey)
  placeholder

  const { data placeholder = await apiClient.post<AdminGroup>(`/admin/groups/${idplaceholder/duplicate`, undefined, {
    headers: { 'Idempotency-Key': idempotencyKey placeholder
  placeholder)

  if (scope) {
    duplicateOperationKeys.delete(scope.key)
    storeDuplicateOperationKey(scope.key, null)
  placeholder
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

export async function listCompositeRoutes(id: number): Promise<CompositeModelRoute[]> {
  const { data placeholder = await apiClient.get<CompositeModelRoute[]>(`/admin/groups/${idplaceholder/composite-routes`)
  return data
placeholder

export async function createCompositeRoute(
  id: number,
  route: CompositeModelRouteInput
): Promise<CompositeModelRoute> {
  const { data placeholder = await apiClient.post<CompositeModelRoute>(
    `/admin/groups/${idplaceholder/composite-routes`,
    route
  )
  return data
placeholder

export async function updateCompositeRoute(
  id: number,
  routeId: number,
  route: CompositeModelRouteInput
): Promise<CompositeModelRoute> {
  const { data placeholder = await apiClient.put<CompositeModelRoute>(
    `/admin/groups/${idplaceholder/composite-routes/${routeIdplaceholder`,
    route
  )
  return data
placeholder

export async function deleteCompositeRoute(
  id: number,
  routeId: number
): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(
    `/admin/groups/${idplaceholder/composite-routes/${routeIdplaceholder`
  )
  return data
placeholder

export async function previewCompositeRoute(
  id: number,
  request: CompositeRoutePreviewRequest
): Promise<CompositeRouteDecision> {
  const { data placeholder = await apiClient.post<CompositeRouteDecision>(
    `/admin/groups/${idplaceholder/composite-routes/preview`,
    request
  )
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
  rate_multiplier?: number | null
  rpm_override?: number | null
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
 * Only touches rate_multiplier column; preserves rpm_override on existing rows.
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
 * RPM override entry for a user in a group
 */
export interface GroupRPMOverrideEntry {
  user_id: number
  user_name: string
  user_email: string
  user_notes: string
  user_status: string
  rpm_override: number
placeholder

/**
 * Get RPM overrides for users in a group (subset of rate-multipliers endpoint).
 */
export async function getGroupRPMOverrides(id: number): Promise<GroupRPMOverrideEntry[]> {
  const { data placeholder = await apiClient.get<GroupRateMultiplierEntry[]>(
    `/admin/groups/${idplaceholder/rate-multipliers`
  )
  return data
    .filter(e => e.rpm_override != null)
    .map(e => ({
      user_id: e.user_id,
      user_name: e.user_name,
      user_email: e.user_email,
      user_notes: e.user_notes,
      user_status: e.user_status,
      rpm_override: e.rpm_override as number
    placeholder))
placeholder

/**
 * Batch set RPM overrides for users in a group.
 * Only touches rpm_override column; preserves rate_multiplier on existing rows.
 */
export async function batchSetGroupRPMOverrides(
  id: number,
  entries: Array<{ user_id: number; rpm_override: number placeholder>
): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.put<{ message: string placeholder>(
    `/admin/groups/${idplaceholder/rpm-overrides`,
    { entries placeholder
  )
  return data
placeholder

/**
 * Clear all RPM overrides for a group (preserves rate_multiplier).
 */
export async function clearGroupRPMOverrides(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/groups/${idplaceholder/rpm-overrides`)
  return data
placeholder

/**
 * Get usage summary (today + yesterday + cumulative cost) for all groups
 * @returns Array of group usage summaries
 */
export async function getUsageSummary(): Promise<
  { group_id: number; today_cost: number; yesterday_cost: number; total_cost: number placeholder[]
> {
  const { data placeholder = await apiClient.get<
    { group_id: number; today_cost: number; yesterday_cost: number; total_cost: number placeholder[]
  >('/admin/groups/usage-summary')
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
  getAllIncludingInactive,
  getLiveCapability,
  getById,
  getModelsListCandidates,
  create,
  duplicate,
  update,
  delete: deleteGroup,
  toggleStatus,
  getStats,
  getGroupApiKeys,
  listCompositeRoutes,
  createCompositeRoute,
  updateCompositeRoute,
  deleteCompositeRoute,
  previewCompositeRoute,
  getGroupRateMultipliers,
  clearGroupRateMultipliers,
  batchSetGroupRateMultipliers,
  getGroupRPMOverrides,
  clearGroupRPMOverrides,
  batchSetGroupRPMOverrides,
  updateSortOrder,
  getUsageSummary,
  getCapacitySummary
placeholder

export default groupsAPI
