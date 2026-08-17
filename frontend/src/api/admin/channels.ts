/**
 * Admin Channels API endpoints
 * Handles channel management for administrators
 */

import { apiClient placeholder from '../client'
import type { BillingMode, ChannelStatus, BillingModelSource placeholder from '@/constants/channel'

export type { BillingMode placeholder from '@/constants/channel'

export interface PricingInterval {
  id?: number
  min_tokens: number
  max_tokens: number | null
  tier_label: string
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  per_request_price: number | null
  sort_order: number
placeholder

export interface ChannelTimePricingPeriod {
  start_time: string
  end_time: string
  multiplier: number
placeholder

export interface ChannelTimePricing {
  timezone: string
  periods: ChannelTimePricingPeriod[]
placeholder

export interface ChannelModelPricing {
  id?: number
  platform: string
  models: string[]
  billing_mode: BillingMode
  input_price: number | null
  output_price: number | null
  cache_write_price: number | null
  cache_read_price: number | null
  image_input_price: number | null
  image_output_price: number | null
  per_request_price: number | null
  intervals: PricingInterval[]
  time_pricing: ChannelTimePricing | null
placeholder

export interface AccountStatsPricingRule {
  id?: number
  name: string
  group_ids: number[]
  account_ids: number[]
  pricing: ChannelModelPricing[]
placeholder

export interface Channel {
  id: number
  name: string
  description: string
  status: ChannelStatus
  billing_model_source: BillingModelSource
  restrict_models: boolean
  features_config?: Record<string, unknown>
  group_ids: number[]
  model_pricing: ChannelModelPricing[]
  model_mapping: Record<string, Record<string, string>> // platform → {src→dstplaceholder
  apply_pricing_to_account_stats: boolean
  account_stats_pricing_rules: AccountStatsPricingRule[]
  created_at: string
  updated_at: string
placeholder

export interface CreateChannelRequest {
  name: string
  description?: string
  group_ids?: number[]
  model_pricing?: ChannelModelPricing[]
  model_mapping?: Record<string, Record<string, string>>
  billing_model_source?: string
  restrict_models?: boolean
  features_config?: Record<string, unknown>
  apply_pricing_to_account_stats?: boolean
  account_stats_pricing_rules?: AccountStatsPricingRule[]
placeholder

export interface UpdateChannelRequest {
  name?: string
  description?: string
  status?: string
  group_ids?: number[]
  model_pricing?: ChannelModelPricing[]
  model_mapping?: Record<string, Record<string, string>>
  billing_model_source?: string
  restrict_models?: boolean
  features_config?: Record<string, unknown>
  apply_pricing_to_account_stats?: boolean
  account_stats_pricing_rules?: AccountStatsPricingRule[]
placeholder

interface PaginatedResponse<T> {
  items: T[]
  total: number
placeholder

/**
 * List channels with pagination
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    status?: string
    search?: string
    sort_by?: string
    sort_order?: 'asc' | 'desc'
  placeholder,
  options?: { signal?: AbortSignal placeholder
): Promise<PaginatedResponse<Channel>> {
  const { data placeholder = await apiClient.get<PaginatedResponse<Channel>>('/admin/channels', {
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
 * Get channel by ID
 */
export async function getById(id: number): Promise<Channel> {
  const { data placeholder = await apiClient.get<Channel>(`/admin/channels/${idplaceholder`)
  return data
placeholder

/**
 * Create a new channel
 */
export async function create(req: CreateChannelRequest): Promise<Channel> {
  const { data placeholder = await apiClient.post<Channel>('/admin/channels', req)
  return data
placeholder

/**
 * Update a channel
 */
export async function update(id: number, req: UpdateChannelRequest): Promise<Channel> {
  const { data placeholder = await apiClient.put<Channel>(`/admin/channels/${idplaceholder`, req)
  return data
placeholder

/**
 * Delete a channel
 */
export async function remove(id: number): Promise<void> {
  await apiClient.delete(`/admin/channels/${idplaceholder`)
placeholder

export interface ModelDefaultPricing {
  found: boolean
  input_price?: number    // per-token price
  output_price?: number
  cache_write_price?: number
  cache_read_price?: number
  image_input_price?: number
  image_output_price?: number
placeholder

export async function getModelDefaultPricing(model: string): Promise<ModelDefaultPricing> {
  const { data placeholder = await apiClient.get<ModelDefaultPricing>('/admin/channels/model-pricing', {
    params: { model placeholder
  placeholder)
  return data
placeholder

export interface SyncPricingModelsResult {
  models: string[]
placeholder

/**
 * Fetch the latest model names from the LiteLLM pricing catalog for the given platform
 */
export async function syncPricingModels(platform: string): Promise<SyncPricingModelsResult> {
  const { data placeholder = await apiClient.get<SyncPricingModelsResult>('/admin/channels/pricing/sync-models', {
    params: { platform placeholder
  placeholder)
  return data
placeholder

const channelsAPI = { list, getById, create, update, remove, getModelDefaultPricing, syncPricingModels placeholder
export default channelsAPI
