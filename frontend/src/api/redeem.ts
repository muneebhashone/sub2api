/**
 * Redeem code API endpoints
 * Handles redeem code redemption for users
 */

import { apiClient placeholder from './client'
import type { RedeemCodeRequest placeholder from '@/types'

export interface RedeemHistoryItem {
  id: number
  code: string
  type: string
  value: number
  status: string
  used_at: string
  created_at: string
  // 订阅类型专用字段
  group_id?: number
  validity_days?: number
  group?: {
    id: number
    name: string
  placeholder
placeholder

/**
 * Redeem a code
 * @param code - Redeem code string
 * @returns Redemption result with updated balance or concurrency
 */
export async function redeem(code: string): Promise<{
  message: string
  type: string
  value: number
  new_balance?: number
  new_concurrency?: number
placeholder> {
  const payload: RedeemCodeRequest = { code placeholder

  const { data placeholder = await apiClient.post<{
    message: string
    type: string
    value: number
    new_balance?: number
    new_concurrency?: number
  placeholder>('/redeem', payload)

  return data
placeholder

/**
 * Get user's redemption history
 * @returns List of redeemed codes
 */
export async function getHistory(): Promise<RedeemHistoryItem[]> {
  const { data placeholder = await apiClient.get<RedeemHistoryItem[]>('/redeem/history')
  return data
placeholder

export const redeemAPI = {
  redeem,
  getHistory
placeholder

export default redeemAPI
