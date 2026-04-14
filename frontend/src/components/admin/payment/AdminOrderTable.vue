<template>
  <div class="space-y-4">
    <div class="card p-4">
      <div class="flex flex-wrap items-center gap-3">
        <div class="flex-1 sm:max-w-64">
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t('payment.admin.searchOrders')"
            class="input"
            @input="handleSearch"
          />
        </div>
        <Select
          v-model="filters.status"
          :options="statusFilterOptions"
          class="w-36"
          @change="emitFiltersChanged"
        />
        <Select
          v-model="filters.payment_type"
          :options="paymentTypeFilterOptions"
          class="w-40"
          @change="emitFiltersChanged"
        />
        <Select
          v-model="filters.order_type"
          :options="orderTypeFilterOptions"
          class="w-36"
          @change="emitFiltersChanged"
        />
        <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
          <button
            @click="emit('refresh')"
            :disabled="loading"
            class="btn btn-secondary"
            :title="t('common.refresh')"
          >
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </div>
    </div>

    <DataTable :columns="columns" :data="orders" :loading="loading">
      <template #cell-id="{ value placeholder">
        <span class="font-mono text-sm">#{{ value placeholderplaceholder</span>
      </template>

      <template #cell-user_id="{ value placeholder">
        <span class="text-sm text-gray-600 dark:text-gray-400">#{{ value placeholderplaceholder</span>
      </template>

      <template #cell-amount="{ value, row placeholder">
        <div class="text-sm">
          <span class="font-medium text-gray-900 dark:text-white">{{ row.order_type === 'balance' ? '$' : '¥' placeholderplaceholder{{ value.toFixed(2) placeholderplaceholder</span>
          <span v-if="row.pay_amount !== value" class="ml-1 text-xs text-gray-500">
            ({{ t('payment.orders.payAmount') placeholderplaceholder: ¥{{ row.pay_amount.toFixed(2) placeholderplaceholder)
          </span>
        </div>
      </template>

      <template #cell-payment_type="{ value placeholder">
        <span class="text-sm text-gray-700 dark:text-gray-300">
          {{ t('payment.methods.' + value, value) placeholderplaceholder
        </span>
      </template>

      <template #cell-status="{ value placeholder">
        <span :class="['badge', statusBadgeClass(value)]">
          {{ t('payment.status.' + value.toLowerCase(), value) placeholderplaceholder
        </span>
      </template>

      <template #cell-order_type="{ value placeholder">
        <span class="text-sm text-gray-700 dark:text-gray-300">
          {{ t('payment.admin.' + value + 'Order', value) placeholderplaceholder
        </span>
      </template>

      <template #cell-created_at="{ value placeholder">
        <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDateTime(value) placeholderplaceholder</span>
      </template>

      <template #cell-actions="{ row placeholder">
        <div class="flex items-center gap-2">
          <button
            @click="emit('detail', row)"
            class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-50 hover:text-gray-700 dark:hover:bg-gray-800/50 dark:hover:text-gray-300"
          >
            <Icon name="eye" size="sm" />
            <span class="text-xs">{{ t('common.view') placeholderplaceholder</span>
          </button>
          <button
            v-if="row.status === 'PENDING'"
            @click="emit('cancel', row)"
            class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-yellow-50 hover:text-yellow-600 dark:hover:bg-yellow-900/20 dark:hover:text-yellow-400"
          >
            <Icon name="x" size="sm" />
            <span class="text-xs">{{ t('payment.orders.cancel') placeholderplaceholder</span>
          </button>
          <button
            v-if="row.status === 'FAILED'"
            @click="emit('retry', row)"
            class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-blue-900/20 dark:hover:text-blue-400"
          >
            <Icon name="refresh" size="sm" />
            <span class="text-xs">{{ t('payment.admin.retry') placeholderplaceholder</span>
          </button>
          <button
            v-if="canRefundRow(row)"
            @click="emit('refund', row)"
            class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
          >
            <Icon name="dollar" size="sm" />
            <span class="text-xs">{{ t('payment.admin.refund') placeholderplaceholder</span>
          </button>
        </div>
      </template>
    </DataTable>

    <Pagination
      v-if="total > 0"
      :page="page"
      :total="total"
      :page-size="pageSize"
      @update:page="emit('update:page', $event)"
      @update:pageSize="emit('update:pageSize', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import type { PaymentOrder placeholder from '@/types/payment'
import type { Column placeholder from '@/components/common/types'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import { statusBadgeClass, canRefund, formatOrderDateTime placeholder from '@/components/payment/orderUtils'

const { t placeholder = useI18n()

defineProps<{
  orders: PaymentOrder[]
  loading: boolean
  page: number
  pageSize: number
  total: number
placeholder>()

const emit = defineEmits<{
  (e: 'detail', order: PaymentOrder): void
  (e: 'cancel', order: PaymentOrder): void
  (e: 'retry', order: PaymentOrder): void
  (e: 'refund', order: PaymentOrder): void
  (e: 'refresh'): void
  (e: 'update:page', page: number): void
  (e: 'update:pageSize', size: number): void
  (e: 'filter', filters: { keyword?: string; status?: string; payment_type?: string; order_type?: string placeholder): void
placeholder>()

const searchQuery = ref('')
const filters = reactive({ status: '', payment_type: '', order_type: '' placeholder)

let debounceTimer: ReturnType<typeof setTimeout> | null = null
function handleSearch() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => emitFiltersChanged(), 300)
placeholder

function emitFiltersChanged() {
  emit('filter', {
    keyword: searchQuery.value || undefined,
    status: filters.status || undefined,
    payment_type: filters.payment_type || undefined,
    order_type: filters.order_type || undefined,
  placeholder)
placeholder

const columns = computed<Column[]>(() => [
  { key: 'id', label: t('payment.orders.orderId') placeholder,
  { key: 'user_id', label: t('payment.orders.userId') placeholder,
  { key: 'amount', label: t('payment.orders.amount') placeholder,
  { key: 'payment_type', label: t('payment.orders.paymentMethod') placeholder,
  { key: 'status', label: t('payment.orders.status') placeholder,
  { key: 'order_type', label: t('payment.orders.orderType') placeholder,
  { key: 'created_at', label: t('payment.orders.createdAt') placeholder,
  { key: 'actions', label: t('payment.orders.actions') placeholder,
])

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
  { value: 'REFUND_FAILED', label: t('payment.status.refund_failed') placeholder,
])

const paymentTypeFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allPaymentTypes') placeholder,
  { value: 'alipay', label: t('payment.methods.alipay') placeholder,
  { value: 'wxpay', label: t('payment.methods.wxpay') placeholder,
  { value: 'stripe', label: t('payment.methods.stripe') placeholder,
])

const orderTypeFilterOptions = computed(() => [
  { value: '', label: t('payment.admin.allOrderTypes') placeholder,
  { value: 'balance', label: t('payment.admin.balanceOrder') placeholder,
  { value: 'subscription', label: t('payment.admin.subscriptionOrder') placeholder,
])

function canRefundRow(order: PaymentOrder): boolean {
  return canRefund(order.status)
placeholder

function formatDateTime(dateStr: string): string {
  return formatOrderDateTime(dateStr)
placeholder
</script>
