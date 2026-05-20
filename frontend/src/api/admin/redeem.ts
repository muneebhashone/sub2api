/**
 * Admin Redeem Codes API endpoints
 * Handles redeem code generation and management for administrators
 */

import { apiClient placeholder from '../client'
import type {
  RedeemCode,
  GenerateRedeemCodesRequest,
  BatchUpdateRedeemCodeFields,
  RedeemCodeType,
  PaginatedResponse
placeholder from '@/types'

/**
 * List all redeem codes with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters
 * @returns Paginated list of redeem codes
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    type?: RedeemCodeType
    status?: 'active' | 'used' | 'expired' | 'unused' | 'disabled'
    search?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  placeholder,
  options?: {
    signal?: AbortSignal
  placeholder
): Promise<PaginatedResponse<RedeemCode>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<RedeemCode>>('/admin/redeem-codes', {
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
 * Get redeem code by ID
 * @param id - Redeem code ID
 * @returns Redeem code details
 */
export async function getById(id: number): Promise<RedeemCode> {
  const { data placeholder = await apiClient.get<RedeemCode>(`/admin/redeem-codes/${idplaceholder`)
  return data
placeholder

/**
 * Generate new redeem codes
 * @param count - Number of codes to generate
 * @param type - Type of redeem code
 * @param value - Value of the code
 * @param groupId - Group ID (required for subscription type)
 * @param validityDays - Validity days (for subscription type)
 * @param expiresInDays - Days before the code itself expires
 * @returns Array of generated redeem codes
 */
export async function generate(
  count: number,
  type: RedeemCodeType,
  value: number,
  groupId?: number | null,
  validityDays?: number,
  expiresInDays?: number | null
): Promise<RedeemCode[]> {
  const payload: GenerateRedeemCodesRequest = {
    count,
    type,
    value
  placeholder

  // 订阅类型专用字段
  if (type === 'subscription') {
    payload.group_id = groupId
    if (validityDays && validityDays > 0) {
      payload.validity_days = validityDays
    placeholder
  placeholder
  if (expiresInDays && expiresInDays > 0) {
    payload.expires_in_days = expiresInDays
  placeholder

  const { data placeholder = await apiClient.post<RedeemCode[]>('/admin/redeem-codes/generate', payload)
  return data
placeholder

/**
 * Delete redeem code
 * @param id - Redeem code ID
 * @returns Success confirmation
 */
export async function deleteCode(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/redeem-codes/${idplaceholder`)
  return data
placeholder

/**
 * Batch delete redeem codes
 * @param ids - Array of redeem code IDs
 * @returns Success confirmation
 */
export async function batchDelete(ids: number[]): Promise<{
  deleted: number
  message: string
placeholder> {
  const { data placeholder = await apiClient.post<{
    deleted: number
    message: string
  placeholder>('/admin/redeem-codes/batch-delete', { ids placeholder)
  return data
placeholder

/**
 * Batch update selected redeem code fields
 * @param ids - Array of redeem code IDs
 * @param fields - Field collection to update
 * @returns Updated count
 */
export async function batchUpdate(
  ids: number[],
  fields: BatchUpdateRedeemCodeFields
): Promise<{
  updated: number
  message: string
placeholder> {
  const { data placeholder = await apiClient.post<{
    updated: number
    message: string
  placeholder>('/admin/redeem-codes/batch-update', { ids, fields placeholder)
  return data
placeholder

/**
 * Expire redeem code
 * @param id - Redeem code ID
 * @returns Updated redeem code
 */
export async function expire(id: number): Promise<RedeemCode> {
  const { data placeholder = await apiClient.post<RedeemCode>(`/admin/redeem-codes/${idplaceholder/expire`)
  return data
placeholder

/**
 * Get redeem code statistics
 * @returns Statistics about redeem codes
 */
export async function getStats(): Promise<{
  total_codes: number
  active_codes: number
  used_codes: number
  expired_codes: number
  total_value_distributed: number
  by_type: Record<RedeemCodeType, number>
placeholder> {
  const { data placeholder = await apiClient.get<{
    total_codes: number
    active_codes: number
    used_codes: number
    expired_codes: number
    total_value_distributed: number
    by_type: Record<RedeemCodeType, number>
  placeholder>('/admin/redeem-codes/stats')
  return data
placeholder

/**
 * Export redeem codes to CSV
 * @param filters - Optional filters
 * @returns CSV data as blob
 */
export async function exportCodes(filters?: {
  type?: RedeemCodeType
  status?: 'used' | 'expired' | 'unused' | 'disabled'
  search?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
placeholder): Promise<Blob> {
  const response = await apiClient.get('/admin/redeem-codes/export', {
    params: filters,
    responseType: 'blob'
  placeholder)
  return response.data
placeholder

export const redeemAPI = {
  list,
  getById,
  generate,
  delete: deleteCode,
  batchDelete,
  batchUpdate,
  expire,
  getStats,
  exportCodes
placeholder

export default redeemAPI
