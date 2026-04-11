/**
 * Admin Payment API endpoints
 * Handles payment management operations for administrators
 */

import { apiClient placeholder from '../client'
import type {
  DashboardStats,
  PaymentOrder,
  PaymentChannel,
  SubscriptionPlan,
  ProviderInstance
placeholder from '@/types/payment'
import type { BasePaginationResponse placeholder from '@/types'

/** Admin-facing payment config returned by GET /admin/payment/config */
export interface AdminPaymentConfig {
  enabled: boolean
  min_amount: number
  max_amount: number
  daily_limit: number
  order_timeout_minutes: number
  max_pending_orders: number
  enabled_payment_types: string[]
  balance_disabled: boolean
  load_balance_strategy: string
  product_name_prefix: string
  product_name_suffix: string
  help_image_url: string
  help_text: string
placeholder

/** Fields accepted by PUT /admin/payment/config (all optional via pointer semantics) */
export interface UpdatePaymentConfigRequest {
  enabled?: boolean
  min_amount?: number
  max_amount?: number
  daily_limit?: number
  order_timeout_minutes?: number
  max_pending_orders?: number
  enabled_payment_types?: string[]
  balance_disabled?: boolean
  load_balance_strategy?: string
  product_name_prefix?: string
  product_name_suffix?: string
  help_image_url?: string
  help_text?: string
placeholder

export const adminPaymentAPI = {
  // ==================== Config ====================

  /** Get payment configuration (admin view) */
  getConfig() {
    return apiClient.get<AdminPaymentConfig>('/admin/payment/config')
  placeholder,

  /** Update payment configuration */
  updateConfig(data: UpdatePaymentConfigRequest) {
    return apiClient.put('/admin/payment/config', data)
  placeholder,

  // ==================== Dashboard ====================

  /** Get payment dashboard statistics */
  getDashboard(days?: number) {
    return apiClient.get<DashboardStats>('/admin/payment/dashboard', {
      params: days ? { days placeholder : undefined
    placeholder)
  placeholder,

  // ==================== Orders ====================

  /** Get all orders (paginated, with filters) */
  getOrders(params?: {
    page?: number
    page_size?: number
    status?: string
    payment_type?: string
    user_id?: number
    keyword?: string
    start_date?: string
    end_date?: string
    order_type?: string
  placeholder) {
    return apiClient.get<BasePaginationResponse<PaymentOrder>>('/admin/payment/orders', { params placeholder)
  placeholder,

  /** Get a specific order by ID */
  getOrder(id: number) {
    return apiClient.get<PaymentOrder>(`/admin/payment/orders/${idplaceholder`)
  placeholder,

  /** Cancel an order (admin) */
  cancelOrder(id: number) {
    return apiClient.post(`/admin/payment/orders/${idplaceholder/cancel`)
  placeholder,

  /** Retry recharge for a failed order */
  retryRecharge(id: number) {
    return apiClient.post(`/admin/payment/orders/${idplaceholder/retry`)
  placeholder,

  /** Process a refund */
  refundOrder(id: number, data: { amount: number; reason: string; deduct_balance?: boolean; force?: boolean placeholder) {
    return apiClient.post(`/admin/payment/orders/${idplaceholder/refund`, data)
  placeholder,

  // ==================== Channels ====================

  /** Get all payment channels */
  getChannels() {
    return apiClient.get<PaymentChannel[]>('/admin/payment/channels')
  placeholder,

  /** Create a payment channel */
  createChannel(data: Partial<PaymentChannel>) {
    return apiClient.post<PaymentChannel>('/admin/payment/channels', data)
  placeholder,

  /** Update a payment channel */
  updateChannel(id: number, data: Partial<PaymentChannel>) {
    return apiClient.put<PaymentChannel>(`/admin/payment/channels/${idplaceholder`, data)
  placeholder,

  /** Delete a payment channel */
  deleteChannel(id: number) {
    return apiClient.delete(`/admin/payment/channels/${idplaceholder`)
  placeholder,

  // ==================== Subscription Plans ====================

  /** Get all subscription plans */
  getPlans() {
    return apiClient.get<SubscriptionPlan[]>('/admin/payment/plans')
  placeholder,

  /** Create a subscription plan */
  createPlan(data: Record<string, unknown>) {
    return apiClient.post<SubscriptionPlan>('/admin/payment/plans', data)
  placeholder,

  /** Update a subscription plan */
  updatePlan(id: number, data: Record<string, unknown>) {
    return apiClient.put<SubscriptionPlan>(`/admin/payment/plans/${idplaceholder`, data)
  placeholder,

  /** Delete a subscription plan */
  deletePlan(id: number) {
    return apiClient.delete(`/admin/payment/plans/${idplaceholder`)
  placeholder,

  // ==================== Provider Instances ====================

  /** Get all provider instances */
  getProviders() {
    return apiClient.get<ProviderInstance[]>('/admin/payment/providers')
  placeholder,

  /** Create a provider instance */
  createProvider(data: Partial<ProviderInstance>) {
    return apiClient.post<ProviderInstance>('/admin/payment/providers', data)
  placeholder,

  /** Update a provider instance */
  updateProvider(id: number, data: Partial<ProviderInstance>) {
    return apiClient.put<ProviderInstance>(`/admin/payment/providers/${idplaceholder`, data)
  placeholder,

  /** Delete a provider instance */
  deleteProvider(id: number) {
    return apiClient.delete(`/admin/payment/providers/${idplaceholder`)
  placeholder
placeholder

export default adminPaymentAPI
