<template>
  <DataTable :columns="columns" :data="orders" :loading="loading">
    <template #cell-id="{ value placeholder">
      <span class="font-mono text-sm">#{{ value placeholderplaceholder</span>
    </template>
    <template #cell-out_trade_no="{ value placeholder">
      <span class="text-sm text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
    </template>
    <template v-if="showUser" #cell-user_email="{ value, row placeholder">
      <div class="text-sm">
        <span class="text-gray-900 dark:text-white">{{ value || row.user_name || '#' + row.user_id placeholderplaceholder</span>
        <span v-if="row.user_notes" class="ml-1 text-xs text-gray-400">({{ row.user_notes placeholderplaceholder)</span>
      </div>
    </template>
    <template #cell-pay_amount="{ value, row placeholder">
      <div class="text-sm">
        <span class="font-medium text-gray-900 dark:text-white">¥{{ value.toFixed(2) placeholderplaceholder</span>
        <span v-if="row.fee_rate > 0" class="ml-1 text-xs text-gray-400" :title="t('payment.orders.fee') + ': ' + row.fee_rate + '%'">
          ({{ t('payment.orders.fee') placeholderplaceholder {{ row.fee_rate placeholderplaceholder%)
        </span>
        <div v-if="row.amount !== row.pay_amount" class="text-xs text-gray-500">
          {{ t('payment.orders.creditedAmount') placeholderplaceholder: {{ row.order_type === 'balance' ? '$' : '¥' placeholderplaceholder{{ row.amount.toFixed(2) placeholderplaceholder
        </div>
      </div>
    </template>
    <template #cell-payment_type="{ value placeholder">
      <span class="text-sm text-gray-700 dark:text-gray-300">{{ t('payment.methods.' + value, value) placeholderplaceholder</span>
    </template>
    <template #cell-status="{ value placeholder">
      <OrderStatusBadge :status="value" />
    </template>
    <template #cell-created_at="{ value placeholder">
      <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDate(value) placeholderplaceholder</span>
    </template>
    <template #cell-actions="{ row placeholder">
      <slot name="actions" :row="row" />
    </template>
  </DataTable>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import type { PaymentOrder placeholder from '@/types/payment'
import type { Column placeholder from '@/components/common/types'
import DataTable from '@/components/common/DataTable.vue'
import OrderStatusBadge from '@/components/payment/OrderStatusBadge.vue'

const { t placeholder = useI18n()

const props = defineProps<{
  orders: PaymentOrder[]
  loading: boolean
  showUser?: boolean
placeholder>()

function formatDate(dateStr: string) { return new Date(dateStr).toLocaleString() placeholder

const columns = computed((): Column[] => {
  const cols: Column[] = [
    { key: 'id', label: t('payment.orders.orderId') placeholder,
    { key: 'out_trade_no', label: t('payment.orders.orderNo') placeholder,
  ]
  if (props.showUser) {
    cols.push({ key: 'user_email', label: t('payment.admin.colUser') placeholder)
  placeholder
  cols.push(
    { key: 'pay_amount', label: t('payment.orders.payAmount') placeholder,
    { key: 'payment_type', label: t('payment.orders.paymentMethod') placeholder,
    { key: 'status', label: t('payment.orders.status') placeholder,
    { key: 'created_at', label: t('payment.orders.createdAt') placeholder,
    { key: 'actions', label: t('common.actions') placeholder,
  )
  return cols
placeholder)
</script>
