/**
 * User-facing Channel Monitor API endpoints
 * Read-only views for end users to inspect channel availability/status.
 */

import { apiClient placeholder from './client'
import type { Provider, MonitorStatus placeholder from './admin/channelMonitor'

export type { Provider, MonitorStatus placeholder from './admin/channelMonitor'

export interface UserMonitorExtraModel {
  model: string
  status: MonitorStatus
  latency_ms: number | null
placeholder

export interface UserMonitorView {
  id: number
  name: string
  provider: Provider
  group_name: string
  primary_model: string
  primary_status: MonitorStatus
  primary_latency_ms: number | null
  availability_7d: number
  extra_models: UserMonitorExtraModel[]
placeholder

export interface UserMonitorListResponse {
  items: UserMonitorView[]
placeholder

export interface UserMonitorModelDetail {
  model: string
  latest_status: MonitorStatus
  latest_latency_ms: number | null
  availability_7d: number
  availability_15d: number
  availability_30d: number
  avg_latency_7d_ms: number | null
placeholder

export interface UserMonitorDetail {
  id: number
  name: string
  provider: Provider
  group_name: string
  models: UserMonitorModelDetail[]
placeholder

/**
 * List all monitor views available to the current user.
 */
export async function list(options?: { signal?: AbortSignal placeholder): Promise<UserMonitorListResponse> {
  const { data placeholder = await apiClient.get<UserMonitorListResponse>('/channel-monitors', {
    signal: options?.signal,
  placeholder)
  return data
placeholder

/**
 * Get detailed status (multi-window availability + latency) for a single monitor.
 */
export async function status(id: number): Promise<UserMonitorDetail> {
  const { data placeholder = await apiClient.get<UserMonitorDetail>(`/channel-monitors/${idplaceholder/status`)
  return data
placeholder

export const channelMonitorUserAPI = {
  list,
  status,
placeholder

export default channelMonitorUserAPI
