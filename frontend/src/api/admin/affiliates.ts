/**
 * Admin Affiliate API endpoints
 * Manage per-user affiliate (邀请返利) configurations:
 * exclusive invite codes (overrides aff_code) and exclusive rebate rates.
 */

import { apiClient placeholder from '../client'
import type { PaginatedResponse placeholder from '@/types'

export interface AffiliateAdminEntry {
  user_id: number
  email: string
  username: string
  aff_code: string
  aff_code_custom: boolean
  aff_rebate_rate_percent?: number | null
  aff_count: number
placeholder

export interface ListAffiliateUsersParams {
  page?: number
  page_size?: number
  search?: string
placeholder

export interface UpdateAffiliateUserRequest {
  aff_code?: string
  aff_rebate_rate_percent?: number | null
  /** Set true to explicitly clear the per-user rate (sets it to NULL). */
  clear_rebate_rate?: boolean
placeholder

export interface BatchSetRateRequest {
  user_ids: number[]
  aff_rebate_rate_percent?: number | null
  /** Set true to clear rates instead of setting. */
  clear?: boolean
placeholder

export interface SimpleUser {
  id: number
  email: string
  username: string
placeholder

export async function listUsers(
  params: ListAffiliateUsersParams = {placeholder,
): Promise<PaginatedResponse<AffiliateAdminEntry>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<AffiliateAdminEntry>>(
    '/admin/affiliates/users',
    {
      params: {
        page: params.page ?? 1,
        page_size: params.page_size ?? 20,
        search: params.search ?? '',
      placeholder,
    placeholder,
  )
  return data
placeholder

export async function lookupUsers(q: string): Promise<SimpleUser[]> {
  const { data placeholder = await apiClient.get<SimpleUser[]>(
    '/admin/affiliates/users/lookup',
    { params: { q placeholder placeholder,
  )
  return data
placeholder

export async function updateUserSettings(
  userId: number,
  payload: UpdateAffiliateUserRequest,
): Promise<{ user_id: number placeholder> {
  const { data placeholder = await apiClient.put<{ user_id: number placeholder>(
    `/admin/affiliates/users/${userIdplaceholder`,
    payload,
  )
  return data
placeholder

export async function clearUserSettings(
  userId: number,
): Promise<{ user_id: number placeholder> {
  const { data placeholder = await apiClient.delete<{ user_id: number placeholder>(
    `/admin/affiliates/users/${userIdplaceholder`,
  )
  return data
placeholder

export async function batchSetRate(
  payload: BatchSetRateRequest,
): Promise<{ affected: number placeholder> {
  const { data placeholder = await apiClient.post<{ affected: number placeholder>(
    '/admin/affiliates/users/batch-rate',
    payload,
  )
  return data
placeholder

export const affiliatesAPI = {
  listUsers,
  lookupUsers,
  updateUserSettings,
  clearUserSettings,
  batchSetRate,
placeholder

export default affiliatesAPI
