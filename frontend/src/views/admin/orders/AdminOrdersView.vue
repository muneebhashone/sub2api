<template>
  <AppLayout>
    <div class="space-y-4">
      <!-- Filters -->
      <div class="card p-4">
        <div class="flex flex-wrap items-center gap-3">
          <div class="flex-1 sm:max-w-64">
            <input v-model="orderSearch" type="text" :placeholder="t('payment.admin.searchOrders')" class="input" @input="debounceLoadOrders" />
          </div>
          <Select v-model="orderFilters.status" :options="statusFilterOptions" class="w-36" @change="loadOrders" />
          <Select v-model="orderFilters.payment_type" :options="paymentTypeFilterOptions" class="w-40" @change="loadOrders" />
          <Select v-model="orderFilters.order_type" :options="orderTypeFilterOptions" class="w-36" @change="loadOrders" />
          <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
            <button @click="loadOrders" :disabled="ordersLoading" class="btn btn-secondary" :title="t('common.refresh')">
              <Icon name="refresh" size="md" :class="ordersLoading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </div>

      <!-- Table -->
      <OrderTable :orders="orders" :loading="ordersLoading" show-user>
        <template #actions="{ row placeholder">
          <div class="flex items-center gap-1">
            <button @click="showOrderDetail(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-gray-600 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-dark-600">
              <Icon name="eye" size="sm" />
              {{ t('common.view') placeholderplaceholder
            </button>
            <button v-if="row.status === 'PENDING'" @click="handleCancelOrder(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-yellow-600 hover:bg-yellow-50 dark:text-yellow-400 dark:hover:bg-yellow-900/20">
              <Icon name="x" size="sm" />
              {{ t('payment.orders.cancel') placeholderplaceholder
            </button>
            <button v-if="row.status === 'FAILED'" @click="handleRetryOrder(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-blue-600 hover:bg-blue-50 dark:text-blue-400 dark:hover:bg-blue-900/20">
              <Icon name="refresh" size="sm" />
              {{ t('payment.admin.retry') placeholderplaceholder
            </button>
            <template v-if="row.status === 'REFUND_REQUESTED'">
              <span v-if="row.refund_amount" class="rounded-full bg-purple-100 px-1.5 py-0.5 text-xs font-medium text-purple-700 dark:bg-purple-900/30 dark:text-purple-300">{{ creditedAmountSymbol placeholderplaceholder{{ row.refund_amount.toFixed(2) placeholderplaceholder</span>
              <button @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-purple-600 hover:bg-purple-50 dark:text-purple-400 dark:hover:bg-purple-900/20">
                <Icon name="check" size="sm" />
                {{ t('payment.admin.approveRefund') placeholderplaceholder
              </button>
            </template>
            <button v-else-if="row.status === 'REFUND_FAILED'" @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-purple-600 hover:bg-purple-50 dark:text-purple-400 dark:hover:bg-purple-900/20">
              <Icon name="refresh" size="sm" />
              {{ t('payment.admin.retryRefund') placeholderplaceholder
            </button>
            <button v-else-if="row.status === 'REFUND_PENDING'" :disabled="refundQueryingIds.has(row.id)" @click="handleQueryRefund(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-orange-600 hover:bg-orange-50 disabled:opacity-60 dark:text-orange-400 dark:hover:bg-orange-900/20">
              <Icon name="refresh" size="sm" :class="refundQueryingIds.has(row.id) ? 'animate-spin' : ''" />
              {{ t('payment.admin.queryRefundStatus') placeholderplaceholder
            </button>
            <button v-else-if="row.status === 'COMPLETED' || row.status === 'PARTIALLY_REFUNDED'" @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20">
              <Icon name="dollar" size="sm" />
              {{ t('payment.admin.refund') placeholderplaceholder
            </button>
          </div>
        </template>
      </OrderTable>
      <Pagination v-if="orderPagination.total > 0" :page="orderPagination.page" :total="orderPagination.total" :page-size="orderPagination.page_size" @update:page="handleOrderPageChange" @update:pageSize="handleOrderPageSizeChange" />
    </div>

    <!-- Order Detail Dialog -->
    <BaseDialog :show="showDetailDialog" :title="t('payment.admin.orderDetail')" width="wide" @close="showDetailDialog = false">
      <div v-if="selectedOrder" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') placeholderplaceholder</p><p class="font-mono text-sm font-medium text-gray-900 dark:text-white">#{{ selectedOrder.id placeholderplaceholder</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderNo') placeholderplaceholder</p><p class="text-sm font-medium text-gray-900 dark:text-white">{{ selectedOrder.out_trade_no placeholderplaceholder</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.status') placeholderplaceholder</p><OrderStatusBadge :status="selectedOrder.status" /></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.amount') placeholderplaceholder</p><p class="text-sm font-medium text-gray-900 dark:text-white">{{ creditedAmountSymbol placeholderplaceholder{{ selectedOrder.amount.toFixed(2) placeholderplaceholder</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.payAmount') placeholderplaceholder</p><p class="text-sm font-medium text-gray-900 dark:text-white">{{ paymentAmountSymbol(selectedOrder) placeholderplaceholder{{ selectedOrder.pay_amount.toFixed(2) placeholderplaceholder</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.paymentMethod') placeholderplaceholder</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ t('payment.methods.' + selectedOrder.payment_type, selectedOrder.payment_type) placeholderplaceholder</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.feeRate') placeholderplaceholder</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ selectedOrder.fee_rate placeholderplaceholder%</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.orders.createdAt') placeholderplaceholder</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.created_at) placeholderplaceholder</p></div>
          <div><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.expiresAt') placeholderplaceholder</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.expires_at) placeholderplaceholder</p></div>
          <div v-if="selectedOrder.paid_at"><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.paidAt') placeholderplaceholder</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.paid_at) placeholderplaceholder</p></div>
          <div v-if="selectedOrder.refund_amount"><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundAmount') placeholderplaceholder</p><p class="text-sm font-medium text-red-600 dark:text-red-400">{{ creditedAmountSymbol placeholderplaceholder{{ selectedOrder.refund_amount.toFixed(2) placeholderplaceholder</p></div>
          <div v-if="selectedOrder.refund_reason" class="col-span-2"><p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundReason') placeholderplaceholder</p><p class="text-sm text-gray-700 dark:text-gray-300">{{ selectedOrder.refund_reason placeholderplaceholder</p></div>
          <!-- Refund request info -->
          <div v-if="selectedOrder.refund_requested_at" class="col-span-2 border-t border-gray-200 pt-3 dark:border-dark-600">
            <p class="mb-2 text-xs font-medium text-purple-600 dark:text-purple-400">{{ t('payment.admin.refundRequestInfo') placeholderplaceholder</p>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundRequestedAt') placeholderplaceholder</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ formatDateTime(selectedOrder.refund_requested_at) placeholderplaceholder</p>
              </div>
              <div>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundRequestedBy') placeholderplaceholder</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">#{{ selectedOrder.refund_requested_by placeholderplaceholder</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.refundRequestReason') placeholderplaceholder</p>
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ selectedOrder.refund_request_reason placeholderplaceholder</p>
              </div>
            </div>
          </div>
        </div>
        <!-- Audit Logs -->
        <div v-if="orderAuditLogs.length > 0" class="border-t border-gray-200 pt-4 dark:border-dark-600">
          <p class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('payment.admin.auditLogs') placeholderplaceholder</p>
          <div class="max-h-48 space-y-2 overflow-y-auto">
            <div v-for="log in orderAuditLogs" :key="log.id" class="rounded-lg border border-gray-100 bg-gray-50 p-2.5 dark:border-dark-600 dark:bg-dark-800">
              <div class="flex items-center justify-between">
                <span class="text-xs font-medium text-gray-700 dark:text-gray-300">{{ log.action placeholderplaceholder</span>
                <span class="text-xs text-gray-400">{{ formatDateTime(log.created_at) placeholderplaceholder</span>
              </div>
              <div v-if="log.detail" class="mt-1 break-all text-xs text-gray-500 dark:text-gray-400">{{ log.detail placeholderplaceholder</div>
              <div v-if="log.operator" class="mt-1 text-xs text-gray-400">{{ t('payment.admin.operator') placeholderplaceholder: {{ log.operator placeholderplaceholder</div>
            </div>
          </div>
        </div>
      </div>
    </BaseDialog>

    <AdminRefundDialog :show="showRefundDialog" :order="selectedOrder" :submitting="refundSubmitting" :require-force="refundRequireForce" :warning="refundWarning" @confirm="handleRefund" @cancel="closeRefundDialog" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminPaymentAPI placeholder from '@/api/admin/payment'
import { extractI18nErrorMessage placeholder from '@/utils/apiError'
import { formatOrderDateTime placeholder from '@/components/payment/orderUtils'
import type { PaymentOrder placeholder from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import AdminRefundDialog from '@/components/admin/payment/AdminRefundDialog.vue'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'
import OrderTable from '@/components/payment/OrderTable.vue'
import { currencySymbol placeholder from '@/components/payment/currency'

interface AuditLog {
  id: number
  action: string
  detail: string | null
  operator: string | null
  created_at: string
placeholder

const { t placeholder = useI18n()
const appStore = useAppStore()

const ordersLoading = ref(false)
const orders = ref<PaymentOrder[]>([])
const orderSearch = ref('')
const orderFilters = reactive({ status: '', payment_type: '', order_type: '' placeholder)
const orderPagination = reactive({ page: 1, page_size: 20, total: 0 placeholder)
const selectedOrder = ref<PaymentOrder | null>(null)
const showDetailDialog = ref(false)
const showRefundDialog = ref(false)
const refundSubmitting = ref(false)
const refundRequireForce = ref(false)
const refundWarning = ref('')
const refundQueryingIds = ref(new Set<number>())
const orderAuditLogs = ref<AuditLog[]>([])
const creditedAmountSymbol = currencySymbol('USD')

function paymentAmountSymbol(order: PaymentOrder | null | undefined): string {
  return currencySymbol(order?.currency)
placeholder

let debounceTimer: ReturnType<typeof setTimeout> | null = null
function debounceLoadOrders() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => loadOrders(), 300)
placeholder

async function loadOrders() {
  ordersLoading.value = true
  try {
    const res = await adminPaymentAPI.getOrders({
      page: orderPagination.page, page_size: orderPagination.page_size,
      keyword: orderSearch.value || undefined, status: orderFilters.status || undefined,
      payment_type: orderFilters.payment_type || undefined, order_type: orderFilters.order_type || undefined,
    placeholder)
    orders.value = res.data.items || []
    orderPagination.total = res.data.total || 0
  placeholder catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  placeholder finally { ordersLoading.value = false placeholder
placeholder

function handleOrderPageChange(page: number) { orderPagination.page = page; loadOrders() placeholder
function handleOrderPageSizeChange(size: number) { orderPagination.page_size = size; orderPagination.page = 1; loadOrders() placeholder

const statusFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allStatuses') placeholder,
  { value: 'PENDING', label: t('payment.status.pending') placeholder,
  { value: 'PAID', label: t('payment.status.paid') placeholder,
  { value: 'COMPLETED', label: t('payment.status.completed') placeholder,
  { value: 'EXPIRED', label: t('payment.status.expired') placeholder,
  { value: 'CANCELLED', label: t('payment.status.cancelled') placeholder,
  { value: 'FAILED', label: t('payment.status.failed') placeholder,
  { value: 'REFUNDED', label: t('payment.status.refunded') placeholder,
  { value: 'REFUND_REQUESTED', label: t('payment.status.refund_requested') placeholder,
  { value: 'REFUND_PENDING', label: t('payment.status.refund_pending') placeholder,
  { value: 'REFUND_FAILED', label: t('payment.status.refund_failed') placeholder,
])

const paymentTypeFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allPaymentTypes') placeholder,
  { value: 'alipay', label: t('payment.methods.alipay') placeholder,
  { value: 'wxpay', label: t('payment.methods.wxpay') placeholder,
  { value: 'stripe', label: t('payment.methods.stripe') placeholder,
  { value: 'airwallex', label: t('payment.methods.airwallex') placeholder,
])

const orderTypeFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allOrderTypes') placeholder,
  { value: 'balance', label: t('payment.admin.balanceOrder') placeholder,
  { value: 'subscription', label: t('payment.admin.subscriptionOrder') placeholder,
])

async function showOrderDetail(order: PaymentOrder) {
  selectedOrder.value = order
  orderAuditLogs.value = []
  showDetailDialog.value = true
  try {
    const res = await adminPaymentAPI.getOrder(order.id)
    const data = res.data as unknown as Record<string, unknown>
    if (data.order) selectedOrder.value = data.order as PaymentOrder
    orderAuditLogs.value = ((data.auditLogs || data.audit_logs || []) as unknown) as AuditLog[]
  placeholder catch (_err: unknown) { /* keep cached order data */ placeholder
placeholder

async function handleCancelOrder(order: PaymentOrder) {
  try { await adminPaymentAPI.cancelOrder(order.id); appStore.showSuccess(t('payment.admin.orderCancelled')); loadOrders() placeholder
  catch (err: unknown) { appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error'))) placeholder
placeholder

async function handleRetryOrder(order: PaymentOrder) {
  try { await adminPaymentAPI.retryRecharge(order.id); appStore.showSuccess(t('payment.admin.retrySuccess')); loadOrders() placeholder
  catch (err: unknown) { appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error'))) placeholder
placeholder

function openRefundDialog(order: PaymentOrder) {
  selectedOrder.value = order
  refundRequireForce.value = false
  refundWarning.value = ''
  showRefundDialog.value = true
placeholder

function closeRefundDialog() {
  showRefundDialog.value = false
  refundRequireForce.value = false
  refundWarning.value = ''
placeholder

function isRefundPendingWarning(warning: string | undefined): boolean {
  return /pending|处理中|待/.test(String(warning || '').toLowerCase())
placeholder

async function handleRefund(data: { amount: number; reason: string; deduct_balance: boolean; force: boolean placeholder) {
  if (!selectedOrder.value) return
  refundSubmitting.value = true
  try {
    const res = await adminPaymentAPI.refundOrder(selectedOrder.value.id, { amount: data.amount, reason: data.reason, deduct_balance: data.deduct_balance, force: data.force placeholder)
    if (res.data.success) {
      appStore.showSuccess(t('payment.admin.refundSuccess'))
      closeRefundDialog()
      loadOrders()
      return
    placeholder
    if (isRefundPendingWarning(res.data.warning)) {
      appStore.showSuccess(t('payment.admin.refundPending'))
      closeRefundDialog()
      loadOrders()
      return
    placeholder
    if (res.data.require_force) {
      // Backend needs an explicit force confirmation (e.g. the user spent their
      // balance after requesting the refund). Keep the dialog open and surface
      // the force checkbox instead of dropping the admin back to the list.
      refundRequireForce.value = true
      refundWarning.value = res.data.warning || ''
      return
    placeholder
    appStore.showError(res.data.warning || t('common.error'))
  placeholder catch (err: unknown) { appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error'))) placeholder
  finally { refundSubmitting.value = false placeholder
placeholder

async function handleQueryRefund(order: PaymentOrder) {
  refundQueryingIds.value = new Set(refundQueryingIds.value).add(order.id)
  try {
    const res = await adminPaymentAPI.queryRefund(order.id)
    if (res.data.success) {
      appStore.showSuccess(t('payment.admin.refundSuccess'))
    placeholder else if (isRefundPendingWarning(res.data.warning)) {
      appStore.showSuccess(t('payment.admin.refundPending'))
    placeholder else {
      appStore.showError(res.data.warning || t('common.error'))
    placeholder
    loadOrders()
  placeholder catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  placeholder finally {
    const next = new Set(refundQueryingIds.value)
    next.delete(order.id)
    refundQueryingIds.value = next
  placeholder
placeholder

function formatDateTime(dateStr: string): string { return formatOrderDateTime(dateStr) placeholder

onMounted(() => loadOrders())
</script>
