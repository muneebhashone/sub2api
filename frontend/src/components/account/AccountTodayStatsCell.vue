<template>
  <div>
    <!-- Loading state -->
    <div v-if="loading" class="space-y-0.5">
      <div class="h-3 w-12 animate-pulse rounded bg-gray-200 dark:bg-gray-700"></div>
      <div class="h-3 w-16 animate-pulse rounded bg-gray-200 dark:bg-gray-700"></div>
      <div class="h-3 w-10 animate-pulse rounded bg-gray-200 dark:bg-gray-700"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="text-xs text-red-500">
      {{ error placeholderplaceholder
    </div>

    <!-- Stats data -->
    <div v-else-if="stats" class="space-y-0.5 text-xs">
      <!-- Requests -->
      <div class="flex items-center gap-1">
        <span class="text-gray-500 dark:text-gray-400"
          >{{ t('admin.accounts.stats.requests') placeholderplaceholder:</span
        >
        <span class="font-medium text-gray-700 dark:text-gray-300">{{
          formatNumber(stats.requests)
        placeholderplaceholder</span>
      </div>
      <!-- Tokens -->
      <div class="flex items-center gap-1">
        <span class="text-gray-500 dark:text-gray-400"
          >{{ t('admin.accounts.stats.tokens') placeholderplaceholder:</span
        >
        <span class="font-medium text-gray-700 dark:text-gray-300">{{
          formatTokens(stats.tokens)
        placeholderplaceholder</span>
      </div>
      <!-- Cost (Account) -->
      <div class="flex items-center gap-1">
        <span class="text-gray-500 dark:text-gray-400">{{ t('usage.accountBilled') placeholderplaceholder:</span>
        <span class="font-medium text-emerald-600 dark:text-emerald-400">{{
          formatCurrency(stats.cost)
        placeholderplaceholder</span>
      </div>
      <!-- Cost (User/API Key) -->
      <div v-if="stats.user_cost != null" class="flex items-center gap-1">
        <span class="text-gray-500 dark:text-gray-400">{{ t('usage.userBilled') placeholderplaceholder:</span>
        <span class="font-medium text-gray-700 dark:text-gray-300">{{
          formatCurrency(stats.user_cost)
        placeholderplaceholder</span>
      </div>
    </div>

    <!-- No data -->
    <div v-else class="text-xs text-gray-400">-</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { adminAPI placeholder from '@/api/admin'
import type { Account, WindowStats placeholder from '@/types'
import { formatNumber, formatCurrency placeholder from '@/utils/format'

const props = defineProps<{
  account: Account
placeholder>()

const { t placeholder = useI18n()

const loading = ref(false)
const error = ref<string | null>(null)
const stats = ref<WindowStats | null>(null)

// Format large token numbers (e.g., 1234567 -> 1.23M)
const formatTokens = (tokens: number): string => {
  if (tokens >= 1000000) {
    return `${(tokens / 1000000).toFixed(2)placeholderM`
  placeholder else if (tokens >= 1000) {
    return `${(tokens / 1000).toFixed(1)placeholderK`
  placeholder
  return tokens.toString()
placeholder

const loadStats = async () => {
  loading.value = true
  error.value = null

  try {
    stats.value = await adminAPI.accounts.getTodayStats(props.account.id)
  placeholder catch (e: any) {
    error.value = 'Failed'
    console.error('Failed to load today stats:', e)
  placeholder finally {
    loading.value = false
  placeholder
placeholder

onMounted(() => {
  loadStats()
placeholder)
</script>
