<template>
  <AppLayout>
    <div class="space-y-4">
      <!-- Filters -->
      <div class="card p-4">
        <div class="flex flex-wrap items-center gap-3">
          <Select v-model="currentFilter" :options="statusFilters" class="w-36" @change="fetchOrders" />
          <div class="flex flex-1 items-center justify-end gap-2">
            <button @click="fetchOrders" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button class="btn btn-primary" @click="router.push('/purchase')">{{ t('payment.result.backToRecharge') placeholderplaceholder</button>
          </div>
        </div>
      </div>

      <!-- Table -->
      <OrderTable :orders="orders" :loading="loading">
        <template #actions="{ row placeholder">
          <div class="flex items-center gap-2">
            <button v-if="row.status === 'PENDING'" @click="handleCancel(row.id)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-yellow-600 hover:bg-yellow-50 dark:text-yellow-400 dark:hover:bg-yellow-900/20">
              <Icon name="x" size="sm" />
              <span>{{ t('payment.orders.cancel') placeholderplaceholder</span>
            </button>
            <button v-if="canRequestRefund(row)" @click="openRefundDialog(row)" class="inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium text-purple-600 hover:bg-purple-50 dark:text-purple-400 dark:hover:bg-purple-900/20">
              <Icon name="dollar" size="sm" />
              <span>{{ t('payment.orders.requestRefund') placeholderplaceholder</span>
            </button>
          </div>
        </template>
      </OrderTable>

      <!-- Pagination -->
      <Pagination
        v-if="pagination.total > 0"
        :page="pagination.page"
        :total="pagination.total"
        :page-size="pagination.page_size"
        @update:page="handlePageChange"
        @update:pageSize="handlePageSizeChange"
      />
    </div>

    <!-- Cancel Confirm Dialog -->
    <BaseDialog :show="!!cancelTargetId" :title="t('payment.orders.cancel')" width="narrow" @close="cancelTargetId = null">
      <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('payment.confirmCancel') placeholderplaceholder</p>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="cancelTargetId = null">{{ t('common.cancel') placeholderplaceholder</button>
          <button class="btn btn-danger" :disabled="actionLoading" @click="confirmCancel">{{ actionLoading ? t('common.processing') : t('payment.orders.cancel') placeholderplaceholder</button>
        </div>
      </template>
    </BaseDialog>

    <!-- Refund Dialog -->
    <BaseDialog :show="!!refundTarget" :title="t('payment.orders.requestRefund')" @close="refundTarget = null">
      <div v-if="refundTarget" class="space-y-4">
        <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-800">
          <div class="flex justify-between text-sm">
            <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.orderId') placeholderplaceholder</span>
            <span class="font-mono text-gray-900 dark:text-white">#{{ refundTarget.id placeholderplaceholder</span>
          </div>
          <div class="mt-2 flex justify-between text-sm">
            <span class="text-gray-500 dark:text-gray-400">{{ t('payment.orders.amount') placeholderplaceholder</span>
            <span class="text-gray-900 dark:text-white">${{ refundTarget.amount.toFixed(2) placeholderplaceholder</span>
          </div>
        </div>
        <div>
          <label class="input-label">{{ t('payment.refundReason') placeholderplaceholder</label>
          <textarea v-model="refundReason" rows="3" class="input mt-1 w-full" :placeholder="t('payment.refundReasonPlaceholder')" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button class="btn btn-secondary" @click="refundTarget = null">{{ t('common.cancel') placeholderplaceholder</button>
          <button class="btn btn-primary" :disabled="actionLoading || !refundReason.trim()" @click="confirmRefund">{{ actionLoading ? t('common.processing') : t('payment.orders.requestRefund') placeholderplaceholder</button>
        </div>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRouter placeholder from 'vue-router'
import { useAppStore placeholder from '@/stores'
import { paymentAPI placeholder from '@/api/payment'
import { extractI18nErrorMessage placeholder from '@/utils/apiError'
import type { PaymentOrder placeholder from '@/types/payment'
import AppLayout from '@/components/layout/AppLayout.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import OrderTable from '@/components/payment/OrderTable.vue'

const { t placeholder = useI18n()
const router = useRouter()
const appStore = useAppStore()

const loading = ref(false)
const actionLoading = ref(false)
const orders = ref<PaymentOrder[]>([])
const refundEligibleProviders = ref<Set<string>>(new Set())
const currentFilter = ref('')
const cancelTargetId = ref<number | null>(null)
const refundTarget = ref<PaymentOrder | null>(null)
const refundReason = ref('')
const pagination = reactive({ page: 1, page_size: 20, total: 0 placeholder)

const statusFilters = computed(() => [
  { value: '', label: t('common.all') placeholder,
  { value: 'PENDING', label: t('payment.status.pending') placeholder,
  { value: 'COMPLETED', label: t('payment.status.completed') placeholder,
  { value: 'FAILED', label: t('payment.status.failed') placeholder,
  { value: 'REFUNDED', label: t('payment.status.refunded') placeholder,
])

async function fetchOrders() {
  loading.value = true
  try {
    const res = await paymentAPI.getMyOrders({
      page: pagination.page,
      page_size: pagination.page_size,
      status: currentFilter.value || undefined,
    placeholder)
    orders.value = res.data.items || []
    pagination.total = res.data.total || 0
  placeholder catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

function handlePageChange(page: number) { pagination.page = page; fetchOrders() placeholder
function handlePageSizeChange(size: number) { pagination.page_size = size; pagination.page = 1; fetchOrders() placeholder

function handleCancel(orderId: number) { cancelTargetId.value = orderId placeholder

async function confirmCancel() {
  if (!cancelTargetId.value) return
  actionLoading.value = true
  try {
    await paymentAPI.cancelOrder(cancelTargetId.value)
    appStore.showSuccess(t('common.success'))
    cancelTargetId.value = null
    await fetchOrders()
  placeholder catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  placeholder finally {
    actionLoading.value = false
  placeholder
placeholder

function openRefundDialog(order: PaymentOrder) { refundTarget.value = order; refundReason.value = '' placeholder

async function confirmRefund() {
  if (!refundTarget.value || !refundReason.value.trim()) return
  actionLoading.value = true
  try {
    await paymentAPI.requestRefund(refundTarget.value.id, { reason: refundReason.value.trim() placeholder)
    appStore.showSuccess(t('common.success'))
    refundTarget.value = null
    refundReason.value = ''
    await fetchOrders()
  placeholder catch (err: unknown) {
    appStore.showError(extractI18nErrorMessage(err, t, 'payment.errors', t('common.error')))
  placeholder finally {
    actionLoading.value = false
  placeholder
placeholder

function canRequestRefund(order: PaymentOrder): boolean {
  if (order.status !== 'COMPLETED') return false
  if (!order.provider_instance_id) return false
  return refundEligibleProviders.value.has(order.provider_instance_id)
placeholder

async function loadRefundEligibility() {
  try {
    const res = await paymentAPI.getRefundEligibleProviders()
    refundEligibleProviders.value = new Set(res.data.provider_instance_ids || [])
  placeholder catch { /* ignore — default to hiding refund button */ placeholder
placeholder

onMounted(() => { fetchOrders(); loadRefundEligibility() placeholder)
</script>
