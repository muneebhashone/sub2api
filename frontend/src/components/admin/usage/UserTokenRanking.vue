<template>
  <div class="card p-4">
    <!-- Header -->
    <div class="mb-3 flex flex-wrap items-center gap-3">
      <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100">
        {{ t('admin.usage.tokenRanking.title') placeholderplaceholder
      </h3>
      <span class="text-xs text-gray-400">{{ t('admin.usage.tokenRanking.subtitle') placeholderplaceholder</span>
      <div class="ml-auto flex flex-wrap items-center gap-2">
        <input
          v-model="search"
          type="text"
          class="input h-8 w-44 text-sm"
          :placeholder="t('admin.usage.tokenRanking.searchPlaceholder')"
        />
        <div class="w-28">
          <Select v-model="limit" :options="limitOptions" @change="load" />
        </div>
        <button
          type="button"
          class="btn btn-secondary h-8 px-2"
          :title="t('common.refresh')"
          @click="load"
        >
          <Icon name="refresh" size="sm" :class="{ 'animate-spin': loading placeholder" />
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-gray-100 text-xs text-gray-500 dark:border-gray-700 dark:text-gray-400">
            <th class="w-10 py-2 pr-2 text-left font-medium">#</th>
            <th class="py-2 pr-2 text-left font-medium">{{ t('admin.usage.tokenRanking.columns.user') placeholderplaceholder</th>
            <th
              v-for="col in sortableColumns"
              :key="col.key"
              class="cursor-pointer select-none py-2 pl-2 text-right font-medium hover:text-primary-500"
              :class="{ 'text-primary-600 dark:text-primary-400': sortBy === col.key placeholder"
              @click="setSort(col.key)"
            >
              {{ t(col.label) placeholderplaceholder
              <span v-if="sortBy === col.key" aria-hidden="true">↓</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="loading">
            <td :colspan="sortableColumns.length + 2" class="py-8 text-center">
              <LoadingSpinner />
            </td>
          </tr>
          <tr v-else-if="displayItems.length === 0">
            <td :colspan="sortableColumns.length + 2" class="py-8 text-center text-sm text-gray-400">
              {{ t('admin.dashboard.noDataAvailable') placeholderplaceholder
            </td>
          </tr>
          <tr
            v-for="(item, index) in displayItems"
            v-else
            :key="item.user_id"
            class="cursor-pointer border-b border-gray-50 transition-colors hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-dark-700/40"
            @click="$emit('select-user', item.user_id, item.email)"
          >
            <td class="py-2 pr-2 text-gray-400">{{ index + 1 placeholderplaceholder</td>
            <td class="max-w-[220px] truncate py-2 pr-2 text-gray-700 dark:text-gray-200" :title="item.email">
              {{ item.email || `User #${item.user_idplaceholder` placeholderplaceholder
            </td>
            <td class="py-2 pl-2 text-right text-gray-500 dark:text-gray-400">{{ item.requests.toLocaleString() placeholderplaceholder</td>
            <td class="py-2 pl-2 text-right text-gray-500 dark:text-gray-400">{{ fmtTokens(item.input_tokens) placeholderplaceholder</td>
            <td class="py-2 pl-2 text-right text-gray-500 dark:text-gray-400">{{ fmtTokens(item.output_tokens) placeholderplaceholder</td>
            <td class="py-2 pl-2 text-right text-gray-500 dark:text-gray-400">{{ fmtTokens(item.cache_tokens) placeholderplaceholder</td>
            <td class="py-2 pl-2 text-right font-medium text-gray-800 dark:text-gray-100">{{ fmtTokens(item.total_tokens) placeholderplaceholder</td>
            <td class="py-2 pl-2 text-right text-green-600 dark:text-green-400">${{ fmtCost(item.actual_cost) placeholderplaceholder</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Footer -->
    <div v-if="!loading && items.length > 0" class="mt-2 text-right text-xs text-gray-400">
      {{ t('admin.usage.tokenRanking.userCount', { count: displayItems.length placeholder) placeholderplaceholder
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { getUserBreakdown, type UserBreakdownParams placeholder from '@/api/admin/dashboard'
import { formatCompactNumber, formatCostFixed placeholder from '@/utils/format'
import type { UserBreakdownItem placeholder from '@/types'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

const props = defineProps<{
  startDate: string
  endDate: string
  filters: Record<string, unknown>
  model?: string
placeholder>()

defineEmits<{ (e: 'select-user', userId: number, email: string): void placeholder>()

const { t placeholder = useI18n()

type SortKey = NonNullable<UserBreakdownParams['sort_by']>
const sortableColumns: { key: SortKey; label: string placeholder[] = [
  { key: 'requests', label: 'admin.usage.tokenRanking.columns.requests' placeholder,
  { key: 'input_tokens', label: 'admin.usage.tokenRanking.columns.inputTokens' placeholder,
  { key: 'output_tokens', label: 'admin.usage.tokenRanking.columns.outputTokens' placeholder,
  { key: 'cache_tokens', label: 'admin.usage.tokenRanking.columns.cacheTokens' placeholder,
  { key: 'total_tokens', label: 'admin.usage.tokenRanking.columns.totalTokens' placeholder,
  { key: 'actual_cost', label: 'admin.usage.tokenRanking.columns.cost' placeholder,
]

const limitOptions = [
  { value: 20, label: 'Top 20' placeholder,
  { value: 50, label: 'Top 50' placeholder,
  { value: 100, label: 'Top 100' placeholder,
  { value: 200, label: 'Top 200' placeholder,
]

const items = ref<UserBreakdownItem[]>([])
const loading = ref(false)
const sortBy = ref<SortKey>('total_tokens')
const limit = ref(50)
const search = ref('')
let reqSeq = 0

const displayItems = computed(() => {
  const kw = search.value.trim().toLowerCase()
  if (!kw) return items.value
  return items.value.filter((u) => (u.email || `user #${u.user_idplaceholder`).toLowerCase().includes(kw))
placeholder)

const fmtTokens = (v: number) => formatCompactNumber(v)
const fmtCost = (v: number) => formatCostFixed(v, 4)

const setSort = (key: SortKey) => {
  if (sortBy.value === key) return
  sortBy.value = key
  load()
placeholder

const load = async () => {
  const seq = ++reqSeq
  loading.value = true
  try {
    const params: UserBreakdownParams = {
      ...props.filters,
      start_date: props.startDate,
      end_date: props.endDate,
      sort_by: sortBy.value,
      limit: limit.value,
    placeholder
    if (props.model) params.model = props.model
    const res = await getUserBreakdown(params)
    if (seq !== reqSeq) return
    items.value = res.users || []
  placeholder catch {
    if (seq !== reqSeq) return
    items.value = []
  placeholder finally {
    if (seq === reqSeq) loading.value = false
  placeholder
placeholder

// Reload when the shared filters / date range / model change.
watch(
  () => [props.startDate, props.endDate, props.model, JSON.stringify(props.filters)],
  () => load(),
  { immediate: true placeholder
)

defineExpose({ reload: load placeholder)
</script>
