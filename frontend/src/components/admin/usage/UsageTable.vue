<template>
  <div class="card overflow-hidden">
    <div class="overflow-auto">
      <DataTable :columns="cols" :data="data" :loading="loading">
        <template #cell-user="{ row placeholder">
          <div class="text-sm">
            <span class="font-medium text-gray-900 dark:text-white">{{ row.user?.email || '-' placeholderplaceholder</span>
            <span class="ml-1 text-gray-500 dark:text-gray-400">#{{ row.user_id placeholderplaceholder</span>
          </div>
        </template>

        <template #cell-api_key="{ row placeholder">
          <span class="text-sm text-gray-900 dark:text-white">{{ row.api_key?.name || '-' placeholderplaceholder</span>
        </template>

        <template #cell-account="{ row placeholder">
          <span class="text-sm text-gray-900 dark:text-white">{{ row.account?.name || '-' placeholderplaceholder</span>
        </template>

        <template #cell-model="{ value placeholder">
          <span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
        </template>

        <template #cell-group="{ row placeholder">
          <span v-if="row.group" class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200">
            {{ row.group.name placeholderplaceholder
          </span>
          <span v-else class="text-sm text-gray-400 dark:text-gray-500">-</span>
        </template>

        <template #cell-stream="{ row placeholder">
          <span class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium" :class="row.stream ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200' : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'">
            {{ row.stream ? t('usage.stream') : t('usage.sync') placeholderplaceholder
          </span>
        </template>

        <template #cell-tokens="{ row placeholder">
          <div class="space-y-1 text-sm">
            <div class="flex items-center gap-2">
              <div class="inline-flex items-center gap-1">
                <Icon name="arrowDown" size="sm" class="h-3.5 w-3.5 text-emerald-500" />
                <span class="font-medium text-gray-900 dark:text-white">{{ row.input_tokens?.toLocaleString() || 0 placeholderplaceholder</span>
              </div>
              <div class="inline-flex items-center gap-1">
                <Icon name="arrowUp" size="sm" class="h-3.5 w-3.5 text-violet-500" />
                <span class="font-medium text-gray-900 dark:text-white">{{ row.output_tokens?.toLocaleString() || 0 placeholderplaceholder</span>
              </div>
            </div>
            <div v-if="row.cache_read_tokens > 0 || row.cache_creation_tokens > 0" class="flex items-center gap-2">
              <div v-if="row.cache_read_tokens > 0" class="inline-flex items-center gap-1">
                <svg class="h-3.5 w-3.5 text-sky-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" /></svg>
                <span class="font-medium text-sky-600 dark:text-sky-400">{{ formatCacheTokens(row.cache_read_tokens) placeholderplaceholder</span>
              </div>
              <div v-if="row.cache_creation_tokens > 0" class="inline-flex items-center gap-1">
                <svg class="h-3.5 w-3.5 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                <span class="font-medium text-amber-600 dark:text-amber-400">{{ formatCacheTokens(row.cache_creation_tokens) placeholderplaceholder</span>
              </div>
            </div>
          </div>
        </template>

        <template #cell-cost="{ row placeholder">
          <span class="font-medium text-green-600 dark:text-green-400">${{ row.actual_cost?.toFixed(6) || '0.000000' placeholderplaceholder</span>
        </template>

        <template #cell-billing_type="{ row placeholder">
          <span class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium" :class="row.billing_type === 1 ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200' : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-200'">
            {{ row.billing_type === 1 ? t('usage.subscription') : t('usage.balance') placeholderplaceholder
          </span>
        </template>

        <template #cell-first_token="{ row placeholder">
          <span v-if="row.first_token_ms != null" class="text-sm text-gray-600 dark:text-gray-400">{{ formatDuration(row.first_token_ms) placeholderplaceholder</span>
          <span v-else class="text-sm text-gray-400 dark:text-gray-500">-</span>
        </template>

        <template #cell-duration="{ row placeholder">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ formatDuration(row.duration_ms) placeholderplaceholder</span>
        </template>

        <template #cell-created_at="{ value placeholder">
          <span class="text-sm text-gray-600 dark:text-gray-400">{{ formatDateTime(value) placeholderplaceholder</span>
        </template>

        <template #cell-request_id="{ row placeholder">
          <div v-if="row.request_id" class="flex items-center gap-1.5 max-w-[120px]">
            <span class="font-mono text-xs text-gray-500 dark:text-gray-400 truncate" :title="row.request_id">{{ row.request_id placeholderplaceholder</span>
            <button @click="copyRequestId(row.request_id)" class="flex-shrink-0 rounded p-0.5 transition-colors hover:bg-gray-100 dark:hover:bg-dark-700" :class="copiedRequestId === row.request_id ? 'text-green-500' : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'" :title="copiedRequestId === row.request_id ? t('keys.copied') : t('keys.copyToClipboard')">
              <svg v-if="copiedRequestId === row.request_id" class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" /></svg>
              <Icon v-else name="copy" size="sm" class="h-3.5 w-3.5" />
            </button>
          </div>
          <span v-else class="text-gray-400 dark:text-gray-500">-</span>
        </template>

        <template #empty><EmptyState :message="t('usage.noRecords')" /></template>
      </DataTable>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { formatDateTime placeholder from '@/utils/format'
import { useAppStore placeholder from '@/stores/app'
import DataTable from '@/components/common/DataTable.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Icon from '@/components/icons/Icon.vue'

defineProps(['data', 'loading'])
const { t placeholder = useI18n()
const appStore = useAppStore()
const copiedRequestId = ref<string | null>(null)

const cols = computed(() => [
  { key: 'user', label: t('admin.usage.user'), sortable: false placeholder,
  { key: 'api_key', label: t('usage.apiKeyFilter'), sortable: false placeholder,
  { key: 'account', label: t('admin.usage.account'), sortable: false placeholder,
  { key: 'model', label: t('usage.model'), sortable: true placeholder,
  { key: 'group', label: t('admin.usage.group'), sortable: false placeholder,
  { key: 'stream', label: t('usage.type'), sortable: false placeholder,
  { key: 'tokens', label: t('usage.tokens'), sortable: false placeholder,
  { key: 'cost', label: t('usage.cost'), sortable: false placeholder,
  { key: 'billing_type', label: t('usage.billingType'), sortable: false placeholder,
  { key: 'first_token', label: t('usage.firstToken'), sortable: false placeholder,
  { key: 'duration', label: t('usage.duration'), sortable: false placeholder,
  { key: 'created_at', label: t('usage.time'), sortable: true placeholder,
  { key: 'request_id', label: t('admin.usage.requestId'), sortable: false placeholder
])

const formatCacheTokens = (tokens: number): string => {
  if (tokens >= 1000000) return `${(tokens / 1000000).toFixed(1)placeholderM`
  if (tokens >= 1000) return `${(tokens / 1000).toFixed(1)placeholderK`
  return tokens.toString()
placeholder

const formatDuration = (ms: number | null | undefined): string => {
  if (ms == null) return '-'
  if (ms < 1000) return `${msplaceholderms`
  return `${(ms / 1000).toFixed(2)placeholders`
placeholder

const copyRequestId = async (requestId: string) => {
  try {
    await navigator.clipboard.writeText(requestId)
    copiedRequestId.value = requestId
    appStore.showSuccess(t('admin.usage.requestIdCopied'))
    setTimeout(() => { copiedRequestId.value = null placeholder, 2000)
  placeholder catch {
    appStore.showError(t('common.copyFailed'))
  placeholder
placeholder
</script>
