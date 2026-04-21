/**
 * User Payment API endpoints
 * Handles payment operations for regular users
 */

import { apiClient placeholder from './client'
import type {
  PaymentConfig,
  SubscriptionPlan,
  PaymentChannel,
  MethodLimitsResponse,
  CheckoutInfoResponse,
  CreateOrderRequest,
  CreateOrderResult,
  PaymentOrder
placeholder from '@/types/payment'
import type { BasePaginationResponse placeholder from '@/types'

export const paymentAPI = {
  /** Get payment configuration (enabled types, limits, etc.) */
  getConfig() {
    return apiClient.get<PaymentConfig>('/payment/config')
  placeholder,

  /** Get available subscription plans */
  getPlans() {
    return apiClient.get<SubscriptionPlan[]>('/payment/plans')
  placeholder,

  /** Get available payment channels */
  getChannels() {
    return apiClient.get<PaymentChannel[]>('/payment/channels')
  placeholder,

  /** Get all checkout page data in a single call */
  getCheckoutInfo() {
    return apiClient.get<CheckoutInfoResponse>('/payment/checkout-info')
  placeholder,

  /** Get payment method limits and fee rates */
  getLimits() {
    return apiClient.get<MethodLimitsResponse>('/payment/limits')
  placeholder,

  /** Create a new payment order */
  createOrder(data: CreateOrderRequest) {
    return apiClient.post<CreateOrderResult>('/payment/orders', data)
  placeholder,

  /** Get current user's orders */
  getMyOrders(params?: { page?: number; page_size?: number; status?: string placeholder) {
    return apiClient.get<BasePaginationResponse<PaymentOrder>>('/payment/orders/my', { params placeholder)
  placeholder,

  /** Get a specific order by ID */
  getOrder(id: number) {
    return apiClient.get<PaymentOrder>(`/payment/orders/${idplaceholder`)
  placeholder,

  /** Cancel a pending order */
  cancelOrder(id: number) {
    return apiClient.post(`/payment/orders/${idplaceholder/cancel`)
  placeholder,

  /** Verify order payment status with upstream provider */
  verifyOrder(outTradeNo: string) {
    return apiClient.post<PaymentOrder>('/payment/orders/verify', { out_trade_no: outTradeNo placeholder)
  placeholder,

  /** Resolve an order from a signed resume token without auth */
  resolveOrderPublicByResumeToken(resumeToken: string) {
    return apiClient.post<PaymentOrder>('/payment/public/orders/resolve', { resume_token: resumeToken placeholder)
  placeholder,

  /** Request a refund for a completed order */
  requestRefund(id: number, data: { reason: string placeholder) {
    return apiClient.post(`/payment/orders/${idplaceholder/refund-request`, data)
  placeholder,

  /** Get provider instance IDs that allow user refund */
  getRefundEligibleProviders() {
    return apiClient.get<{ provider_instance_ids: string[] placeholder>('/payment/orders/refund-eligible-providers')
  placeholder
placeholder
