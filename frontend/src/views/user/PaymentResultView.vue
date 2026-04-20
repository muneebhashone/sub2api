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
          <div v-else
            class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/30">
            <svg class="h-10 w-10 text-red-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
          <h2 class="mt-4 text-2xl font-bold text-gray-900 dark:text-white">
            {{ isSuccess ? t('payment.result.success') : t('payment.result.failed') placeholderplaceholder
          </h2>
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
import { ref, computed, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRoute, useRouter placeholder from 'vue-router'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'
import { PAYMENT_RECOVERY_STORAGE_KEY, readPaymentRecoverySnapshot placeholder from '@/components/payment/paymentFlow'
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
  return !!order.value && SUCCESS_STATUSES.has(order.value.status)
placeholder)

function normalizedOrderPaymentType(paymentType: string): string {
  return normalizePaymentMethodForDisplay(paymentType) || paymentType
placeholder

onMounted(async () => {
  const resumeToken = typeof route.query.resume_token === 'string'
    ? route.query.resume_token
    : ''
  let orderId = Number(route.query.order_id) || 0
  const outTradeNo = String(route.query.out_trade_no || '')

  if (!orderId && resumeToken && typeof window !== 'undefined') {
    const restored = readPaymentRecoverySnapshot(
      window.localStorage.getItem(PAYMENT_RECOVERY_STORAGE_KEY),
      { resumeToken placeholder,
    )
    if (restored?.orderId) {
      orderId = restored.orderId
    placeholder
  placeholder

  if (!order.value && !orderId && resumeToken) {
    try {
      const result = await paymentAPI.resolveOrderPublicByResumeToken(resumeToken)
      order.value = result.data
      orderId = result.data.id
    placeholder catch (_err: unknown) {
      // Resume token recovery failed, continue to legacy fallback paths.
    placeholder
  placeholder

  if (!order.value && orderId) {
    try {
      order.value = await paymentStore.pollOrderStatus(orderId)
    placeholder catch (_err: unknown) {
      // Order lookup failed, will try legacy fallback below when possible.
    placeholder
  placeholder

  if (!order.value && outTradeNo) {
    returnInfo.value = {
      outTradeNo,
      money: String(route.query.money || ''),
      type: String(route.query.type || ''),
      tradeStatus: String(route.query.trade_status || ''),
    placeholder

    try {
      const result = await paymentAPI.verifyOrderPublic(outTradeNo)
      order.value = result.data
    placeholder catch (_err: unknown) {
      try {
        const result = await paymentAPI.verifyOrder(outTradeNo)
        order.value = result.data
      placeholder catch (_e: unknown) { /* fall through */ placeholder
    placeholder
  placeholder

  if (!order.value && orderId) {
    try {
      order.value = await paymentStore.pollOrderStatus(orderId)
    placeholder catch (_err: unknown) {
      // Order lookup failed, will show returnInfo fallback.
    placeholder
  placeholder
  loading.value = false
placeholder)
</script>
