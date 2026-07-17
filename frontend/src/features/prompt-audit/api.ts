import { apiClient placeholder from '@/api/client'
import type {
  PromptAuditConfig,
  PromptAuditEvent,
  PromptAuditGroup,
  PromptAuditRuntime,
  PromptAuditUpdateRequest,
  PromptDeletePreview,
  PromptDeleteResult,
  PromptEventFilters,
  PromptEventPage,
  PromptProbeResult,
  PromptAuditEndpointDraft,
placeholder from './types'
import { eventFilterPayload, eventQueryParams placeholder from './viewModel'

const basePath = '/admin/prompt-audit'

export async function getConfig(): Promise<PromptAuditConfig> {
  const { data placeholder = await apiClient.get<PromptAuditConfig>(`${basePathplaceholder/config`)
  return data
placeholder

export async function updateConfig(payload: PromptAuditUpdateRequest): Promise<PromptAuditConfig> {
  const { data placeholder = await apiClient.put<PromptAuditConfig>(`${basePathplaceholder/config`, payload)
  return data
placeholder

export async function probeEndpoint(endpoint: PromptAuditEndpointDraft): Promise<PromptProbeResult> {
  const { data placeholder = await apiClient.post<PromptProbeResult>(`${basePathplaceholder/endpoints/probe`, {
    endpoint: {
      id: endpoint.id,
      name: endpoint.name,
      protocol: 'openai_compatible',
      base_url: endpoint.base_url,
      model: endpoint.model,
      token: endpoint.token || undefined,
      timeout_ms: endpoint.timeout_ms,
      input_limit: endpoint.input_limit,
      enabled: endpoint.enabled,
    placeholder,
  placeholder)
  return data
placeholder

export async function getRuntime(): Promise<PromptAuditRuntime> {
  const { data placeholder = await apiClient.get<PromptAuditRuntime>(`${basePathplaceholder/runtime`)
  return data
placeholder

export async function listEvents(
  filters: PromptEventFilters,
  page: number,
  pageSize: number,
): Promise<PromptEventPage> {
  const { data placeholder = await apiClient.get<PromptEventPage>(`${basePathplaceholder/events`, {
    params: { page, page_size: pageSize, ...eventQueryParams(filters) placeholder,
  placeholder)
  return data
placeholder

export async function getEvent(id: number): Promise<PromptAuditEvent> {
  const { data placeholder = await apiClient.get<PromptAuditEvent>(`${basePathplaceholder/events/${idplaceholder`)
  return data
placeholder

export async function deleteEvent(id: number): Promise<PromptDeleteResult> {
  const { data placeholder = await apiClient.delete<PromptDeleteResult>(`${basePathplaceholder/events/${idplaceholder`)
  return data
placeholder

export async function batchDeleteEvents(ids: number[]): Promise<PromptDeleteResult> {
  const { data placeholder = await apiClient.post<PromptDeleteResult>(`${basePathplaceholder/events/batch-delete`, { ids placeholder)
  return data
placeholder

export async function previewDelete(filters: PromptEventFilters): Promise<PromptDeletePreview> {
  const { data placeholder = await apiClient.post<PromptDeletePreview>(
    `${basePathplaceholder/events/delete-preview`,
    eventFilterPayload(filters),
  )
  return data
placeholder

export async function deleteEventsByFilter(
  filters: PromptEventFilters,
  preview: PromptDeletePreview,
): Promise<PromptDeleteResult> {
  const { data placeholder = await apiClient.post<PromptDeleteResult>(`${basePathplaceholder/events/delete-by-filter`, {
    filter: eventFilterPayload(filters),
    snapshot_max_id: preview.snapshot_max_id,
    filter_hash: preview.filter_hash,
    confirmation_token: preview.confirmation_token,
    confirm: true,
  placeholder)
  return data
placeholder

export async function listGroups(): Promise<PromptAuditGroup[]> {
  const { data placeholder = await apiClient.get<PromptAuditGroup[]>('/admin/groups/all', {
    params: { include_inactive: true placeholder,
  placeholder)
  return data
placeholder

export const promptAuditAPI = {
  getConfig,
  updateConfig,
  probeEndpoint,
  getRuntime,
  listEvents,
  getEvent,
  deleteEvent,
  batchDeleteEvents,
  previewDelete,
  deleteEventsByFilter,
  listGroups,
placeholder

export default promptAuditAPI
