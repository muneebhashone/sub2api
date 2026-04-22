import type { LocationQuery, LocationQueryRaw placeholder from 'vue-router'
import type { SubscriptionPlan placeholder from '@/types/payment'
import { normalizeVisibleMethod placeholder from '@/components/payment/paymentFlow'

export interface ParsedWechatResumeRoute {
  orderAmount: number
  orderType: 'balance' | 'subscription'
  paymentType: string
  planId?: number
  openid?: string
  wechatResumeToken?: string
placeholder

function readQueryString(query: LocationQuery, key: string): string {
  const value = query[key]
  if (Array.isArray(value)) {
    return typeof value[0] === 'string' ? value[0] : ''
  placeholder
  return typeof value === 'string' ? value : ''
placeholder

export function hasWechatResumeQuery(query: LocationQuery): boolean {
  if (readQueryString(query, 'wechat_resume') === '1') {
    return true
  placeholder
  return readQueryString(query, 'wechat_resume_token') !== ''
    || readQueryString(query, 'openid') !== ''
placeholder

export function parseWechatResumeRoute(
  query: LocationQuery,
  plans: SubscriptionPlan[],
  fallbackBalanceAmount: number,
): ParsedWechatResumeRoute | null {
  if (!hasWechatResumeQuery(query)) {
    return null
  placeholder

  const wechatResumeToken = readQueryString(query, 'wechat_resume_token')
  if (wechatResumeToken) {
    return {
      wechatResumeToken,
      paymentType: 'wxpay',
      orderType: 'balance',
      orderAmount: 0,
    placeholder
  placeholder

  const openid = readQueryString(query, 'openid')
  if (!openid) {
    return null
  placeholder

  const paymentType = normalizeVisibleMethod(readQueryString(query, 'payment_type')) || 'wxpay'
  const orderType = readQueryString(query, 'order_type') === 'subscription' ? 'subscription' : 'balance'
  const planId = Number.parseInt(readQueryString(query, 'plan_id'), 10)
  const rawAmount = Number.parseFloat(readQueryString(query, 'amount'))
  const orderAmount = Number.isFinite(rawAmount) && rawAmount > 0
    ? rawAmount
    : (orderType === 'subscription'
      ? (plans.find(plan => plan.id === planId)?.price ?? 0)
      : fallbackBalanceAmount)

  return {
    openid,
    paymentType,
    orderType,
    orderAmount,
    planId: Number.isFinite(planId) && planId > 0 ? planId : undefined,
  placeholder
placeholder

export function stripWechatResumeQuery(query: LocationQuery): LocationQueryRaw {
  const nextQuery: LocationQueryRaw = { ...query placeholder
  delete nextQuery.wechat_resume
  delete nextQuery.wechat_resume_token
  delete nextQuery.openid
  delete nextQuery.state
  delete nextQuery.scope
  delete nextQuery.payment_type
  delete nextQuery.amount
  delete nextQuery.order_type
  delete nextQuery.plan_id
  return nextQuery
placeholder
