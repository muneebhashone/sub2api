<template>
  <div class="flex min-h-screen items-center justify-center bg-slate-50 p-4 dark:bg-slate-950">
    <div
      class="w-full max-w-md space-y-4 rounded-2xl border border-slate-200 bg-white p-6 shadow-lg dark:border-slate-700 dark:bg-slate-900"
    >
      <!-- Amount + Order ID -->
      <div v-if="amount" class="text-center">
        <p class="text-3xl font-bold" :style="{ color: methodColor placeholder">¥{{ amount placeholderplaceholder</p>
        <p v-if="orderId" class="mt-1 text-sm text-gray-500 dark:text-slate-400">
          {{ t('payment.orders.orderId') placeholderplaceholder: {{ orderId placeholderplaceholder
        </p>
      </div>

      <!-- Error -->
      <div v-if="error" class="space-y-3">
        <div
          class="rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-700 dark:bg-red-900/30 dark:text-red-400"
        >
          {{ error placeholderplaceholder
        </div>
        <button
          class="w-full text-sm underline dark:text-blue-400 dark:hover:text-blue-300"
          :style="{ color: methodColor placeholder"
          @click="closeWindow"
        >
          {{ t('common.close') placeholderplaceholder
        </button>
      </div>

      <!-- Success -->
      <div v-else-if="success" class="space-y-3 py-4 text-center">
        <div class="text-5xl text-green-600 dark:text-green-400">✓</div>
        <p class="text-sm text-gray-500 dark:text-slate-400">{{ t('payment.result.success') placeholderplaceholder</p>
        <button
          class="text-sm underline dark:text-blue-400 dark:hover:text-blue-300"
          :style="{ color: methodColor placeholder"
          @click="closeWindow"
        >
          {{ t('common.close') placeholderplaceholder
        </button>
      </div>

      <!-- Loading / Redirecting -->
      <div v-else class="flex items-center justify-center py-8">
        <div
          class="h-8 w-8 animate-spin rounded-full border-2 border-t-transparent"
          :style="{ borderColor: methodColor, borderTopColor: 'transparent' placeholder"
        />
        <span class="ml-3 text-sm text-gray-500 dark:text-slate-400">{{ hint placeholderplaceholder</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRoute placeholder from 'vue-router'
import { extractI18nErrorMessage placeholder from '@/utils/apiError'
import { isMobileDevice placeholder from '@/utils/device'
import { buildApiUrl placeholder from '@/api/client'

interface StripeWithWechatPay {
  confirmWechatPayPayment(clientSecret: string, options: Record<string, unknown>): Promise<{ error?: { message?: string placeholder; paymentIntent?: { status: string placeholder placeholder>
placeholder

const METHOD_COLORS: Record<string, string> = {
  alipay: '#00AEEF',
  wechat_pay: '#07C160',
placeholder
const DEFAULT_METHOD_COLOR = '#635bff'

const { t placeholder = useI18n()
const route = useRoute()

const orderId = String(route.query.order_id || '')
const method = String(route.query.method || 'alipay')
const amount = String(route.query.amount || '')

const methodColor = computed(() => METHOD_COLORS[method] || DEFAULT_METHOD_COLOR)

const error = ref('')
const success = ref(false)
const hint = ref(t('payment.stripePopup.redirecting'))

let pollTimer: ReturnType<typeof setInterval> | null = null
let initTimeoutTimer: ReturnType<typeof setTimeout> | null = null
let messageHandler: ((event: MessageEvent) => void) | null = null

function closeWindow() { window.close() placeholder

function clearInitTimeout() {
  if (initTimeoutTimer) {
    clearTimeout(initTimeoutTimer)
    initTimeoutTimer = null
  placeholder
placeholder

onMounted(() => {
  messageHandler = (event: MessageEvent) => {
    if (event.origin !== window.location.origin) return
    if (event.data?.type !== 'STRIPE_POPUP_INIT') return
    // INIT 已到达，取消兜底超时，避免长时间的扫码支付被误判为超时。
    clearInitTimeout()
    if (messageHandler) {
      window.removeEventListener('message', messageHandler)
      messageHandler = null
    placeholder
    initStripe(event.data.clientSecret, event.data.publishableKey)
  placeholder
  window.addEventListener('message', messageHandler)

  if (window.opener) {
    window.opener.postMessage({ type: 'STRIPE_POPUP_READY' placeholder, window.location.origin)
  placeholder

  // 仅兜底“父窗口始终未发 STRIPE_POPUP_INIT”的场景。
  initTimeoutTimer = setTimeout(() => {
    if (!error.value && !success.value) {
      error.value = t('payment.stripePopup.timeout')
    placeholder
  placeholder, 15000)
placeholder)

onUnmounted(() => {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null placeholder
  clearInitTimeout()
  if (messageHandler) {
    window.removeEventListener('message', messageHandler)
    messageHandler = null
  placeholder
placeholder)

async function initStripe(clientSecret: string, publishableKey: string) {
  if (!clientSecret || !publishableKey) {
    error.value = t('payment.stripeMissingParams')
    return
  placeholder
  try {
    const { loadStripe placeholder = await import('@stripe/stripe-js/pure')
    const stripe = await loadStripe(publishableKey)
    if (!stripe) { error.value = t('payment.stripeLoadFailed'); return placeholder

    const returnUrl = window.location.origin + '/payment/result?order_id=' + orderId + '&status=success'

    if (method === 'alipay') {
      // Alipay: redirect this popup to Alipay payment page
      const { error: err placeholder = await stripe.confirmAlipayPayment(clientSecret, { return_url: returnUrl placeholder)
      if (err) error.value = err.message || t('payment.result.failed')
    placeholder else if (method === 'wechat_pay') {
      // WeChat: Stripe shows its built-in QR dialog, user scans, promise resolves
      hint.value = t('payment.stripePopup.loadingQr')
      const result = await (stripe as unknown as StripeWithWechatPay).confirmWechatPayPayment(clientSecret, {
        payment_method_options: { wechat_pay: { client: isMobileDevice() ? 'mobile_web' : 'web' placeholder placeholder,
      placeholder)
      if (result.error) {
        error.value = result.error.message || t('payment.result.failed')
      placeholder else if (result.paymentIntent?.status === 'succeeded') {
        success.value = true
        setTimeout(closeWindow, 2000)
      placeholder else {
        // Payment not completed (user closed QR dialog)
        startPolling()
      placeholder
    placeholder
  placeholder catch (err: unknown) {
    error.value = extractI18nErrorMessage(err, t, 'payment.errors', t('payment.stripeLoadFailed'))
  placeholder
placeholder

function startPolling() {
  let inFlight = false
  pollTimer = setInterval(async () => {
    // 防重入：接口响应慢于轮询间隔时避免并发重叠请求。
    if (inFlight) return
    inFlight = true
    try {
      // access token 存储在 localStorage 的 'auth_token' 键下（见 api/client.ts），
      // 之前误读 'token' 导致轮询请求不带认证、永远 401，支付成功无法被检测到。
      const token = localStorage.getItem('auth_token') || ''
      const res = await fetch(buildApiUrl(`/payment/orders/${orderIdplaceholder`), {
        headers: token ? { Authorization: 'Bearer ' + token placeholder : {placeholder,
        credentials: 'include',
      placeholder)
      if (!res.ok) return
      const data = await res.json()
      const status = data?.data?.status
      if (status === 'COMPLETED' || status === 'PAID') {
        if (pollTimer) { clearInterval(pollTimer); pollTimer = null placeholder
        success.value = true
        setTimeout(closeWindow, 2000)
      placeholder
    placeholder catch { /* ignore */ placeholder finally {
      inFlight = false
    placeholder
  placeholder, 3000)
placeholder
</script>
