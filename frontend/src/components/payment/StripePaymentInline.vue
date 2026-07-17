<template>
  <div class="space-y-4">
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="h-8 w-8 animate-spin rounded-full border-4 border-primary-500 border-t-transparent"></div>
    </div>
    <div v-else-if="initError" class="card p-6 text-center">
      <p class="text-sm text-red-600 dark:text-red-400">{{ initError placeholderplaceholder</p>
      <button class="btn btn-secondary mt-4" @click="$emit('back')">{{ t('payment.result.backToRecharge') placeholderplaceholder</button>
    </div>
    <!-- Success -->
    <template v-else-if="success">
      <div class="card p-6">
        <div class="flex flex-col items-center space-y-4 py-4">
          <div class="flex h-16 w-16 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
            <Icon name="check" size="lg" class="text-green-500" />
          </div>
          <p class="text-lg font-bold text-gray-900 dark:text-white">{{ t('payment.result.success') placeholderplaceholder</p>
          <div class="w-full rounded-xl bg-gray-50 p-4 dark:bg-dark-800">
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') placeholderplaceholder</span>
                <span class="font-medium text-gray-900 dark:text-white">#{{ orderId placeholderplaceholder</span>
              </div>
              <div v-if="amount > 0" class="flex justify-between">
                <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.amount') placeholderplaceholder</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ creditedAmountSymbol placeholderplaceholder{{ amount.toFixed(2) placeholderplaceholder</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') placeholderplaceholder</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ paymentAmountSymbol placeholderplaceholder{{ payAmount.toFixed(2) placeholderplaceholder</span>
              </div>
            </div>
          </div>
          <button class="btn btn-primary" @click="$emit('done')">{{ t('common.confirm') placeholderplaceholder</button>
        </div>
      </div>
    </template>
    <template v-else>
      <!-- Amount -->
      <div class="card overflow-hidden">
        <div class="bg-gradient-to-br from-[#635bff] to-[#4f46e5] px-6 py-5 text-center">
          <p class="text-sm font-medium text-indigo-200">{{ t('payment.actualPay') placeholderplaceholder</p>
          <p class="mt-1 text-3xl font-bold text-white">{{ paymentAmountSymbol placeholderplaceholder{{ payAmount.toFixed(2) placeholderplaceholder</p>
        </div>
      </div>
      <!-- Stripe Payment Element -->
      <div class="card p-6">
        <div ref="stripeMount" class="min-h-[200px]"></div>
        <p v-if="error" class="mt-4 text-sm text-red-600 dark:text-red-400">{{ error placeholderplaceholder</p>
        <button class="btn btn-stripe mt-6 w-full py-3 text-base" :disabled="submitting || !ready" @click="handlePay">
          <span v-if="submitting" class="flex items-center justify-center gap-2">
            <span class="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent"></span>
            {{ t('common.processing') placeholderplaceholder
          </span>
          <span v-else>{{ t('payment.stripePay') placeholderplaceholder</span>
        </button>
      </div>
      <!-- Cancel order -->
      <button class="btn btn-secondary w-full" :disabled="cancelling" @click="handleCancel">
        {{ cancelling ? t('common.processing') : t('payment.qr.cancelOrder') placeholderplaceholder
      </button>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, nextTick placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRouter placeholder from 'vue-router'
import { extractI18nErrorMessage placeholder from '@/utils/apiError'
import { paymentAPI placeholder from '@/api/payment'
import { useAppStore placeholder from '@/stores'
import { getPaymentPopupFeatures placeholder from '@/components/payment/providerConfig'
import { currencySymbol placeholder from '@/components/payment/currency'
import type { Stripe, StripeElements placeholder from '@stripe/stripe-js'
import Icon from '@/components/icons/Icon.vue'

// Stripe payment methods that open a popup (redirect or QR code)
const POPUP_METHODS = new Set(['alipay', 'wechat_pay'])

const props = defineProps<{
  orderId: number
  amount: number
  clientSecret: string
  orderType?: 'balance' | 'subscription'
  publishableKey: string
  payAmount: number
  currency?: string
placeholder>()

const emit = defineEmits<{ success: []; done: []; back: []; redirect: [orderId: number, payUrl: string] placeholder>()

const { t placeholder = useI18n()
const router = useRouter()
const appStore = useAppStore()

const stripeMount = ref<HTMLElement | null>(null)
const loading = ref(true)
const initError = ref('')
const error = ref('')
const submitting = ref(false)
const cancelling = ref(false)
const success = ref(false)
const ready = ref(false)
const selectedType = ref('')
const creditedAmountSymbol = currencySymbol('USD')
const paymentAmountSymbol = computed(() => currencySymbol(props.currency))

let stripeInstance: Stripe | null = null
let elementsInstance: StripeElements | null = null

onMounted(async () => {
  try {
    const { loadStripe placeholder = await import('@stripe/stripe-js/pure')
    const stripe = await loadStripe(props.publishableKey)
    if (!stripe) { initError.value = t('payment.stripeLoadFailed'); return placeholder

    stripeInstance = stripe
    loading.value = false
    await nextTick()
    if (!stripeMount.value) return

    const isDark = document.documentElement.classList.contains('dark')
    const elements = stripe.elements({
      clientSecret: props.clientSecret,
      appearance: { theme: isDark ? 'night' : 'stripe', variables: { borderRadius: '8px' placeholder placeholder,
    placeholder)
    elementsInstance = elements
    const paymentElement = elements.create('payment', {
      layout: 'tabs',
      paymentMethodOrder: ['alipay', 'wechat_pay', 'card', 'link'],
    placeholder as Record<string, unknown>)
    paymentElement.mount(stripeMount.value)
    paymentElement.on('ready', () => { ready.value = true placeholder)
    paymentElement.on('change', (event: { value: { type: string placeholder placeholder) => {
      selectedType.value = event.value.type
    placeholder)
  placeholder catch (err: unknown) {
    initError.value = extractI18nErrorMessage(err, t, 'payment.errors', t('payment.stripeLoadFailed'))
  placeholder finally {
    loading.value = false
  placeholder
placeholder)

async function handlePay() {
  if (!stripeInstance || !elementsInstance || submitting.value) return

  // Alipay / WeChat Pay: open popup for redirect or QR display
  if (POPUP_METHODS.has(selectedType.value)) {
    const popupUrl = router.resolve({
      path: '/payment/stripe-popup',
      query: {
        order_id: String(props.orderId),
        method: selectedType.value,
        amount: String(props.payAmount),
      placeholder,
    placeholder).href
    const popup = window.open(popupUrl, 'paymentPopup', getPaymentPopupFeatures())

    const onReady = (event: MessageEvent) => {
      if (event.source !== popup || event.data?.type !== 'STRIPE_POPUP_READY') return
      window.removeEventListener('message', onReady)
      popup?.postMessage({
        type: 'STRIPE_POPUP_INIT',
        clientSecret: props.clientSecret,
        publishableKey: props.publishableKey,
      placeholder, window.location.origin)
    placeholder
    window.addEventListener('message', onReady)

    emit('redirect', props.orderId, popupUrl)
    return
  placeholder

  // Card / Link: confirm inline
  submitting.value = true
  error.value = ''
  try {
    const { error: stripeError placeholder = await stripeInstance.confirmPayment({
      elements: elementsInstance,
      confirmParams: {
        return_url: window.location.origin + '/payment/result?order_id=' + props.orderId + '&status=success',
      placeholder,
      redirect: 'if_required',
    placeholder)
    if (stripeError) {
      error.value = stripeError.message || t('payment.result.failed')
    placeholder else {
      success.value = true
      emit('success')
    placeholder
  placeholder catch (err: unknown) {
    error.value = extractI18nErrorMessage(err, t, 'payment.errors', t('payment.result.failed'))
  placeholder finally {
    submitting.value = false
  placeholder
placeholder

async function handleCancel() {
  if (!props.orderId || cancelling.value) return
  cancelling.value = true
  try {
    await paymentAPI.cancelOrder(props.orderId)
    emit('back')
  placeholder catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  placeholder finally {
    cancelling.value = false
  placeholder
placeholder
</script>
