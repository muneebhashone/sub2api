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

export interface ListAffiliateRecordsParams {
  page?: number
  page_size?: number
  search?: string
  start_at?: string
  end_at?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
  timezone?: string
placeholder

export interface AffiliateInviteRecord {
  inviter_id: number
  inviter_email: string
  inviter_username: string
  invitee_id: number
  invitee_email: string
  invitee_username: string
  aff_code: string
  total_rebate: number
  created_at: string
placeholder

export interface AffiliateRebateRecord {
  order_id: number
  out_trade_no: string
  inviter_id: number
  inviter_email: string
  inviter_username: string
  invitee_id: number
  invitee_email: string
  invitee_username: string
  order_amount: number
  pay_amount: number
  rebate_amount: number
  payment_type: string
  order_status: string
  created_at: string
placeholder

export interface AffiliateTransferRecord {
  ledger_id: number
  user_id: number
  user_email: string
  username: string
  amount: number
  balance_after?: number | null
  available_quota_after?: number | null
  frozen_quota_after?: number | null
  history_quota_after?: number | null
  snapshot_available: boolean
  created_at: string
placeholder

export interface AffiliateUserOverview {
  user_id: number
  email: string
  username: string
  aff_code: string
  rebate_rate_percent: number
  invited_count: number
  rebated_invitee_count: number
  available_quota: number
  history_quota: number
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

function recordParams(params: ListAffiliateRecordsParams = {placeholder) {
  return {
    page: params.page ?? 1,
    page_size: params.page_size ?? 20,
    search: params.search ?? '',
    start_at: params.start_at || undefined,
    end_at: params.end_at || undefined,
    sort_by: params.sort_by || undefined,
    sort_order: params.sort_order || undefined,
    timezone: params.timezone || undefined,
  placeholder
placeholder

export async function listInviteRecords(
  params: ListAffiliateRecordsParams = {placeholder,
): Promise<PaginatedResponse<AffiliateInviteRecord>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<AffiliateInviteRecord>>(
    '/admin/affiliates/invites',
    { params: recordParams(params) placeholder,
  )
  return data
placeholder

export async function listRebateRecords(
  params: ListAffiliateRecordsParams = {placeholder,
): Promise<PaginatedResponse<AffiliateRebateRecord>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<AffiliateRebateRecord>>(
    '/admin/affiliates/rebates',
    { params: recordParams(params) placeholder,
  )
  return data
placeholder

export async function listTransferRecords(
  params: ListAffiliateRecordsParams = {placeholder,
): Promise<PaginatedResponse<AffiliateTransferRecord>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<AffiliateTransferRecord>>(
    '/admin/affiliates/transfers',
    { params: recordParams(params) placeholder,
  )
  return data
placeholder

export async function getUserOverview(
  userId: number,
): Promise<AffiliateUserOverview> {
  const { data placeholder = await apiClient.get<AffiliateUserOverview>(
    `/admin/affiliates/users/${userIdplaceholder/overview`,
  )
  return data
placeholder

export const affiliatesAPI = {
  listUsers,
  lookupUsers,
  updateUserSettings,
  clearUserSettings,
  batchSetRate,
  listInviteRecords,
  listRebateRecords,
  listTransferRecords,
  getUserOverview,
placeholder

export default affiliatesAPI
