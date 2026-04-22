<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 dark:bg-dark-900">
    <div class="w-full max-w-md space-y-6">
      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
      </div>
      <template v-else>
        <!-- Status Icon -->
        <div class="text-center">
          <div v-if="isSuccess"
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
            <svg class="h-10 w-10 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor"
              stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <div v-else-if="isPending"
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-yellow-100 dark:bg-yellow-900/30">
            <div class="h-10 w-10 animate-spin rounded-full border-4 border-yellow-500 border-t-transparent"></div>
          </div>
          <div v-else
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
            <svg class="h-10 w-10 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
          <h2 class="mt-4 text-2xl font-bold text-gray-900 dark:text-white">
            {{ statusTitle placeholderplaceholder
          </h2>
          <p v-if="isPending" class="mt-2 text-sm text-gray-500 dark:text-gray-400">
            {{ t('payment.result.processingHint') placeholderplaceholder
          </p>
        </div>
        <!-- Order Info -->
        <div v-if="order" class="rounded-xl bg-white p-5 shadow-sm dark:bg-dark-800">
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">#{{ order.id placeholderplaceholder</span>
            </div>
            <div v-if="order.out_trade_no" class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderNo') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ order.out_trade_no placeholderplaceholder</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.baseAmount') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">&#165;{{ baseAmount.toFixed(2) placeholderplaceholder</span>
            </div>
            <div v-if="order.fee_rate > 0" class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.fee') placeholderplaceholder ({{ order.fee_rate placeholderplaceholder%)</span>
              <span class="font-medium text-gray-900 dark:text-white">&#165;{{ feeAmount.toFixed(2) placeholderplaceholder</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') placeholderplaceholder</span>
              <span class="font-bold text-primary-600 dark:text-primary-400">&#165;{{ order.pay_amount.toFixed(2) placeholderplaceholder</span>
            </div>
            <div v-if="order.amount !== order.pay_amount" class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.creditedAmount') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ order.order_type === 'balance' ? '$' : '¥' placeholderplaceholder{{ order.amount.toFixed(2) placeholderplaceholder</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.paymentMethod') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ t(paymentMethodI18nKey(order.payment_type), normalizedOrderPaymentType(order.payment_type)) placeholderplaceholder</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.status') placeholderplaceholder</span>
              <OrderStatusBadge :status="order.status" />
            </div>
          </div>
        </div>
        <!-- EasyPay return info (when no order loaded) -->
        <div v-else-if="returnInfo" class="rounded-xl bg-white p-5 shadow-sm dark:bg-dark-800">
          <div class="space-y-3 text-sm">
            <div v-if="returnInfo.outTradeNo" class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ returnInfo.outTradeNo placeholderplaceholder</span>
            </div>
            <div v-if="returnInfo.money" class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">&#165;{{ returnInfo.money placeholderplaceholder</span>
            </div>
            <div v-if="returnInfo.type" class="flex justify-between">
              <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.paymentMethod') placeholderplaceholder</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ t(paymentMethodI18nKey(returnInfo.type), normalizedOrderPaymentType(returnInfo.type)) placeholderplaceholder</span>
            </div>
          </div>
        </div>
        <!-- Actions -->
        <div class="flex gap-3">
          <button class="btn btn-secondary flex-1" @click="router.push('/purchase')">{{ t('payment.result.backToRecharge') placeholderplaceholder</button>
          <button class="btn btn-primary flex-1" @click="router.push('/orders')">{{ t('payment.result.viewOrders') placeholderplaceholder</button>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRoute, useRouter placeholder from 'vue-router'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'
import {
  PAYMENT_RECOVERY_STORAGE_KEY,
  clearPaymentRecoverySnapshot,
  readPaymentRecoverySnapshot,
placeholder from '@/components/payment/paymentFlow'
import { usePaymentStore placeholder from '@/stores/payment'
import { paymentAPI placeholder from '@/api/payment'
import type { PaymentOrder placeholder from '@/types/payment'
import { normalizePaymentMethodForDisplay, paymentMethodI18nKey placeholder from './paymentUx'

const { t placeholder = useI18n()
const route = useRoute()
const router = useRouter()
const paymentStore = usePaymentStore()

const order = ref<PaymentOrder | null>(null)
const loading = ref(true)

interface ReturnInfo {
  outTradeNo: string
  money: string
  type: string
  tradeStatus: string
placeholder
const returnInfo = ref<ReturnInfo | null>(null)

const SUCCESS_STATUSES = new Set(['COMPLETED', 'PAID', 'RECHARGING'])
const PENDING_STATUSES = new Set(['PENDING', 'CREATED', 'WAITING', 'PROCESSING'])
const STATUS_REFRESH_INTERVAL_MS = 2000
const STATUS_REFRESH_MAX_ATTEMPTS = 15

let statusRefreshTimer: ReturnType<typeof setTimeout> | null = null
const refreshAttempts = ref(0)

/** 充值金额 = pay_amount / (1 + fee_rate/100)，fee_rate=0 时等于 pay_amount */
const baseAmount = computed(() => {
  if (!order.value || order.value.fee_rate <= 0) return order.value?.pay_amount ?? 0
  return Math.round((order.value.pay_amount / (1 + order.value.fee_rate / 100)) * 100) / 100
placeholder)

/** 手续费 = pay_amount - baseAmount */
const feeAmount = computed(() => {
  if (!order.value || order.value.fee_rate <= 0) return 0
  return Math.round((order.value.pay_amount - baseAmount.value) * 100) / 100
placeholder)

const isSuccess = computed(() => {
  return isSuccessStatus(order.value?.status)
placeholder)

const isPending = computed(() => {
  return isPendingStatus(order.value?.status)
placeholder)

const statusTitle = computed(() => {
  if (isSuccess.value) {
    return t('payment.result.success')
  placeholder
  if (isPending.value) {
    return t('payment.result.processing')
  placeholder
  return t('payment.result.failed')
placeholder)

function normalizedOrderPaymentType(paymentType: string): string {
  return normalizePaymentMethodForDisplay(paymentType) || paymentType
placeholder

function normalizeOrderStatus(status: string | null | undefined): string {
  return String(status || '').trim().toUpperCase()
placeholder

function isSuccessStatus(status: string | null | undefined): boolean {
  return SUCCESS_STATUSES.has(normalizeOrderStatus(status))
placeholder

function isPendingStatus(status: string | null | undefined): boolean {
  return PENDING_STATUSES.has(normalizeOrderStatus(status))
placeholder

async function resolveOrderFromResumeToken(resumeToken: string): Promise<PaymentOrder | null> {
  try {
    const result = await paymentAPI.resolveOrderPublicByResumeToken(resumeToken)
    return result.data
  placeholder catch (_err: unknown) {
    return null
  placeholder
placeholder

function clearStatusRefreshTimer(): void {
  if (statusRefreshTimer !== null) {
    clearTimeout(statusRefreshTimer)
    statusRefreshTimer = null
  placeholder
placeholder

function clearRecoverySnapshot(): void {
  if (typeof window === 'undefined') return
  clearPaymentRecoverySnapshot(window.localStorage, PAYMENT_RECOVERY_STORAGE_KEY)
placeholder

function clearRecoverySnapshotForTerminalStatus(status: string | null | undefined): void {
  if (!status) return
  if (!isPendingStatus(status)) {
    clearRecoverySnapshot()
  placeholder
placeholder

function scheduleStatusRefresh(refreshOrder: (() => Promise<PaymentOrder | null>) | null): void {
  clearStatusRefreshTimer()
  if (!refreshOrder || !isPending.value || refreshAttempts.value >= STATUS_REFRESH_MAX_ATTEMPTS) {
    return
  placeholder

  statusRefreshTimer = setTimeout(async () => {
    refreshAttempts.value += 1
    const refreshedOrder = await refreshOrder()
    if (refreshedOrder) {
      order.value = refreshedOrder
      clearRecoverySnapshotForTerminalStatus(refreshedOrder.status)
    placeholder

    if (isPendingStatus(order.value?.status)) {
      scheduleStatusRefresh(refreshOrder)
    placeholder
  placeholder, STATUS_REFRESH_INTERVAL_MS)
placeholder

onMounted(async () => {
  const resumeToken = typeof route.query.resume_token === 'string'
    ? route.query.resume_token
    : ''
  const routeOrderId = Number(route.query.order_id) || 0
  const outTradeNo = String(route.query.out_trade_no || '')
  let orderId = 0

  if (resumeToken && typeof window !== 'undefined') {
    const restored = readPaymentRecoverySnapshot(
      window.localStorage.getItem(PAYMENT_RECOVERY_STORAGE_KEY),
      { resumeToken placeholder,
    )
    if (restored?.orderId) {
      orderId = restored.orderId
    placeholder
  placeholder

  if (!order.value && resumeToken && orderId) {
    try {
      order.value = await paymentStore.pollOrderStatus(orderId)
    placeholder catch (_err: unknown) {
      // Fall through to signed resume-token recovery below.
    placeholder
  placeholder

  if (resumeToken) {
    const resolvedOrder = await resolveOrderFromResumeToken(resumeToken)
    if (resolvedOrder) {
      order.value = resolvedOrder
      if (!orderId) {
        orderId = resolvedOrder.id
      placeholder
    placeholder
  placeholder

  if (!resumeToken) {
    orderId = routeOrderId
  placeholder

  if (!order.value && !resumeToken && orderId) {
    try {
      order.value = await paymentStore.pollOrderStatus(orderId)
    placeholder catch (_err: unknown) {
      // Order lookup failed, will try legacy fallback below when possible.
    placeholder
  placeholder

  const hasLegacyFallbackContext = typeof route.query.trade_status === 'string'
    && route.query.trade_status.trim() !== ''
  if (!order.value && !resumeToken && !orderId && outTradeNo && hasLegacyFallbackContext) {
    returnInfo.value = {
      outTradeNo,
      money: String(route.query.money || ''),
      type: String(route.query.type || ''),
      tradeStatus: String(route.query.trade_status || ''),
    placeholder
  placeholder

  const refreshOrder = async (): Promise<PaymentOrder | null> => {
    if (resumeToken) {
      return await resolveOrderFromResumeToken(resumeToken)
    placeholder

    if (orderId) {
      return await paymentStore.pollOrderStatus(orderId)
    placeholder

    return null
  placeholder

  if (isPendingStatus(order.value?.status)) {
    scheduleStatusRefresh(refreshOrder)
  placeholder else if (order.value) {
    clearRecoverySnapshotForTerminalStatus(order.value.status)
  placeholder else if (returnInfo.value) {
    clearRecoverySnapshot()
  placeholder
  loading.value = false
placeholder)

onBeforeUnmount(() => {
  clearStatusRefreshTimer()
placeholder)
</script>
