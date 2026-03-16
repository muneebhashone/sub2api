<template>
  <div class="bg-gray-50/50 dark:bg-dark-700/30">
    <div v-if="loading" class="flex items-center justify-center py-3">
      <LoadingSpinner />
    </div>
    <div v-else-if="items.length === 0" class="py-2 text-center text-xs text-gray-400">
      {{ t('admin.dashboard.noDataAvailable') placeholderplaceholder
    </div>
    <table v-else class="w-full text-xs">
      <thead>
        <tr class="text-gray-400 dark:text-gray-500">
          <th class="py-1 pl-6 text-left">{{ t('admin.dashboard.spendingRankingUser') placeholderplaceholder</th>
          <th class="py-1 text-right">{{ t('admin.dashboard.requests') placeholderplaceholder</th>
          <th class="py-1 text-right">{{ t('admin.dashboard.tokens') placeholderplaceholder</th>
          <th class="py-1 text-right">{{ t('admin.dashboard.actual') placeholderplaceholder</th>
          <th class="py-1 pr-1 text-right">{{ t('admin.dashboard.standard') placeholderplaceholder</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="user in items"
          :key="user.user_id"
          class="border-t border-gray-100/50 dark:border-gray-700/50"
        >
          <td class="max-w-[120px] truncate py-1 pl-6 text-gray-600 dark:text-gray-300" :title="user.email">
            {{ user.email || `User #${user.user_idplaceholder` placeholderplaceholder
          </td>
          <td class="py-1 text-right text-gray-500 dark:text-gray-400">
            {{ user.requests.toLocaleString() placeholderplaceholder
          </td>
          <td class="py-1 text-right text-gray-500 dark:text-gray-400">
            {{ formatTokens(user.total_tokens) placeholderplaceholder
          </td>
          <td class="py-1 text-right text-green-600 dark:text-green-400">
            ${{ formatCost(user.actual_cost) placeholderplaceholder
          </td>
          <td class="py-1 pr-1 text-right text-gray-400 dark:text-gray-500">
            ${{ formatCost(user.cost) placeholderplaceholder
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { useI18n placeholder from 'vue-i18n'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { UserBreakdownItem placeholder from '@/types'

const { t placeholder = useI18n()

defineProps<{
  items: UserBreakdownItem[]
  loading?: boolean
placeholder>()

const formatTokens = (value: number): string => {
  if (value >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(2)placeholderB`
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(2)placeholderM`
  if (value >= 1_000) return `${(value / 1_000).toFixed(2)placeholderK`
  return value.toLocaleString()
placeholder

const formatCost = (value: number): string => {
  if (value >= 1000) return (value / 1000).toFixed(2) + 'K'
  if (value >= 1) return value.toFixed(2)
  if (value >= 0.01) return value.toFixed(3)
  return value.toFixed(4)
placeholder
</script>
