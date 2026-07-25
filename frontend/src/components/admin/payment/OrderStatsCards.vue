<template>
  <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
    <!-- Today Revenue -->
    <div class="card p-4">
      <div class="flex items-center gap-3">
        <div class="rounded-lg bg-green-100 p-2 dark:bg-green-900/30">
          <Icon name="dollar" size="md" class="text-green-600 dark:text-green-400" :stroke-width="2" />
        </div>
        <div>
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('payment.admin.todayRevenue') placeholderplaceholder</p>
          <p v-for="[currency, amount] in sortedAmounts(stats.today_amount)" :key="currency" class="text-xl font-bold text-gray-900 dark:text-white">
            {{ formatMoney(currency, amount) placeholderplaceholder
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-400">
            {{ stats.today_count placeholderplaceholder {{ t('payment.admin.orders') placeholderplaceholder
          </p>
        </div>
      </div>
    </div>

    <!-- Total Revenue -->
    <div class="card p-4">
      <div class="flex items-center gap-3">
        <div class="rounded-lg bg-blue-100 p-2 dark:bg-blue-900/30">
          <Icon name="creditCard" size="md" class="text-blue-600 dark:text-blue-400" :stroke-width="2" />
        </div>
        <div>
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('payment.admin.totalRevenue') placeholderplaceholder</p>
          <p v-for="[currency, amount] in sortedAmounts(stats.total_amount)" :key="currency" class="text-xl font-bold text-gray-900 dark:text-white">
            {{ formatMoney(currency, amount) placeholderplaceholder
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-400">
            {{ stats.total_count placeholderplaceholder {{ t('payment.admin.orders') placeholderplaceholder
          </p>
        </div>
      </div>
    </div>

    <!-- Today Orders -->
    <div class="card p-4">
      <div class="flex items-center gap-3">
        <div class="rounded-lg bg-purple-100 p-2 dark:bg-purple-900/30">
          <Icon name="chart" size="md" class="text-purple-600 dark:text-purple-400" :stroke-width="2" />
        </div>
        <div>
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('payment.admin.todayOrders') placeholderplaceholder</p>
          <p class="text-xl font-bold text-gray-900 dark:text-white">{{ stats.today_count placeholderplaceholder</p>
        </div>
      </div>
    </div>

    <!-- Average Amount -->
    <div class="card p-4">
      <div class="flex items-center gap-3">
        <div class="rounded-lg bg-amber-100 p-2 dark:bg-amber-900/30">
          <Icon name="chart" size="md" class="text-amber-600 dark:text-amber-400" :stroke-width="2" />
        </div>
        <div>
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('payment.admin.avgAmount') placeholderplaceholder</p>
          <p v-for="[currency, amount] in sortedAmounts(stats.avg_amount)" :key="currency" class="text-xl font-bold text-gray-900 dark:text-white">
            {{ formatMoney(currency, amount) placeholderplaceholder
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n placeholder from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import type { CurrencyAmounts, DashboardStats placeholder from '@/types/payment'

const { t placeholder = useI18n()

defineProps<{
  stats: DashboardStats
placeholder>()

function sortedAmounts(amounts: CurrencyAmounts): [string, number][] {
  return Object.entries(amounts).sort(([left], [right]) => left.localeCompare(right))
placeholder

function formatMoney(currency: string, amount: number): string {
  return new Intl.NumberFormat(undefined, { style: 'currency', currency placeholder).format(amount)
placeholder
</script>
