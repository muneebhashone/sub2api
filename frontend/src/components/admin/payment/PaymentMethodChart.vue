<template>
  <div class="card p-4">
    <h3 class="mb-4 text-sm font-semibold text-gray-900 dark:text-white">
      {{ t('payment.admin.paymentDistribution') placeholderplaceholder
    </h3>
    <div
      v-if="!methods?.length"
      class="flex h-32 items-center justify-center text-sm text-gray-500 dark:text-gray-400"
    >
      {{ t('payment.admin.noData') placeholderplaceholder
    </div>
    <div v-else class="space-y-3">
      <div v-for="method in methods" :key="method.type" class="space-y-1">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span :class="['inline-block h-3 w-3 rounded-full', colorMap[method.type] || 'bg-gray-400']"></span>
            <span class="text-sm text-gray-700 dark:text-gray-300">
              {{ t('payment.methods.' + method.type, method.type) placeholderplaceholder
            </span>
          </div>
          <div class="text-right">
            <span class="text-sm font-medium text-gray-900 dark:text-white">
              ${{ method.amount.toFixed(2) placeholderplaceholder
            </span>
            <span class="ml-2 text-xs text-gray-500 dark:text-gray-400">
              ({{ method.count placeholderplaceholder)
            </span>
          </div>
        </div>
        <div class="h-2 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-dark-700">
          <div
            :class="['h-full rounded-full transition-all', barColorMap[method.type] || 'bg-gray-400']"
            :style="{ width: barWidth(method.amount) + '%' placeholder"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'

const { t placeholder = useI18n()

const props = defineProps<{
  methods: { type: string; amount: number; count: number placeholder[]
placeholder>()

const colorMap: Record<string, string> = {
  alipay: 'bg-blue-500',
  wxpay: 'bg-green-500',
  alipay_direct: 'bg-blue-400',
  wxpay_direct: 'bg-green-400',
  stripe: 'bg-purple-500',
placeholder

const barColorMap: Record<string, string> = {
  alipay: 'bg-blue-500',
  wxpay: 'bg-green-500',
  alipay_direct: 'bg-blue-400',
  wxpay_direct: 'bg-green-400',
  stripe: 'bg-purple-500',
placeholder

const maxAmount = computed(() => {
  if (!props.methods?.length) return 1
  return Math.max(...props.methods.map(m => m.amount), 1)
placeholder)

function barWidth(amount: number): number {
  return Math.min((amount / maxAmount.value) * 100, 100)
placeholder
</script>
