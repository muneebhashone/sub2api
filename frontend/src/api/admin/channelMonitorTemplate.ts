/**
 * Admin Channel Monitor Request Template API.
 *
 * 模板 = 一组可复用的 headers + 可选 body 覆盖配置。
 * 应用到监控 = 拷贝快照；模板后续变动不自动同步，需手动点「应用到关联监控」刷新。
 */

import { apiClient placeholder from '../client'
import type { APIMode, BodyOverrideMode, Provider placeholder from './channelMonitor'

export interface ChannelMonitorTemplate {
  id: number
  name: string
  provider: Provider
  api_mode: APIMode
  description: string
  extra_headers: Record<string, string>
  body_override_mode: BodyOverrideMode
  body_override: Record<string, unknown> | null
  created_at: string
  updated_at: string
  /** 关联的监控数量（快照来自此模板，仅 template_id 匹配即可） */
  associated_monitors: number
placeholder

export interface ListParams {
  provider?: Provider
  api_mode?: APIMode
placeholder

export interface ListResponse {
  items: ChannelMonitorTemplate[]
placeholder

export interface CreateParams {
  name: string
  provider: Provider
  api_mode?: APIMode
  description?: string
  extra_headers?: Record<string, string>
  body_override_mode?: BodyOverrideMode
  body_override?: Record<string, unknown> | null
placeholder

export interface UpdateParams {
  name?: string
  api_mode?: APIMode
  description?: string
  extra_headers?: Record<string, string>
  body_override_mode?: BodyOverrideMode
  body_override?: Record<string, unknown> | null
placeholder

export interface ApplyResponse {
  affected: number
placeholder

export interface AssociatedMonitorBrief {
  id: number
  name: string
  provider: Provider
  api_mode: APIMode
  enabled: boolean
placeholder

export interface AssociatedMonitorsResponse {
  items: AssociatedMonitorBrief[]
placeholder

export async function list(params: ListParams = {placeholder): Promise<ListResponse> {
  const { data placeholder = await apiClient.get<ListResponse>('/admin/channel-monitor-templates', {
    params,
  placeholder)
  return data
placeholder

export async function get(id: number): Promise<ChannelMonitorTemplate> {
  const { data placeholder = await apiClient.get<ChannelMonitorTemplate>(
    `/admin/channel-monitor-templates/${idplaceholder`,
  )
  return data
placeholder

export async function create(params: CreateParams): Promise<ChannelMonitorTemplate> {
  const { data placeholder = await apiClient.post<ChannelMonitorTemplate>(
    '/admin/channel-monitor-templates',
    params,
  )
  return data
placeholder

export async function update(id: number, params: UpdateParams): Promise<ChannelMonitorTemplate> {
  const { data placeholder = await apiClient.put<ChannelMonitorTemplate>(
    `/admin/channel-monitor-templates/${idplaceholder`,
    params,
  )
  return data
placeholder

export async function del(id: number): Promise<void> {
  await apiClient.delete(`/admin/channel-monitor-templates/${idplaceholder`)
placeholder

/**
 * Apply the template to the specified associated monitors (overwrite snapshot fields).
 * monitorIds must be a non-empty subset of the template's associated monitors.
 * Returns count of actually affected monitors.
 */
export async function apply(id: number, monitorIds: number[]): Promise<ApplyResponse> {
  const { data placeholder = await apiClient.post<ApplyResponse>(
    `/admin/channel-monitor-templates/${idplaceholder/apply`,
    { monitor_ids: monitorIds placeholder,
  )
  return data
placeholder

/**
 * List monitors currently associated to this template (used by apply picker).
 */
export async function listAssociatedMonitors(id: number): Promise<AssociatedMonitorsResponse> {
  const { data placeholder = await apiClient.get<AssociatedMonitorsResponse>(
    `/admin/channel-monitor-templates/${idplaceholder/monitors`,
  )
  return data
placeholder

export const channelMonitorTemplateAPI = {
  list,
  get,
  create,
  update,
  del,
  apply,
  listAssociatedMonitors,
placeholder

export default channelMonitorTemplateAPI
