<template>
  <div class="card p-4">
    <div class="mb-4 flex items-center justify-between gap-3">
      <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
        {{ !enableRankingView || activeView === 'model_distribution'
          ? t('admin.dashboard.modelDistribution')
          : t('admin.dashboard.spendingRankingTitle') placeholderplaceholder
      </h3>
      <div class="flex items-center gap-2">
        <div
          v-if="showMetricToggle"
          class="inline-flex rounded-lg border border-gray-200 bg-gray-50 p-0.5 dark:border-gray-700 dark:bg-dark-800"
        >
          <button
            type="button"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition-colors"
            :class="metric === 'tokens'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="emit('update:metric', 'tokens')"
          >
            {{ t('admin.dashboard.metricTokens') placeholderplaceholder
          </button>
          <button
            type="button"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition-colors"
            :class="metric === 'actual_cost'
              ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
              : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'"
            @click="emit('update:metric', 'actual_cost')"
          >
            {{ t('admin.dashboard.metricActualCost') placeholderplaceholder
          </button>
        </div>
        <div v-if="enableRankingView" class="inline-flex rounded-lg bg-gray-100 p-1 dark:bg-dark-800">
          <button
            type="button"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition-colors"
            :class="
              activeView === 'model_distribution'
                ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
                : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
            "
            @click="activeView = 'model_distribution'"
          >
            {{ t('admin.dashboard.viewModelDistribution') placeholderplaceholder
          </button>
          <button
            type="button"
            class="rounded-md px-2.5 py-1 text-xs font-medium transition-colors"
            :class="
              activeView === 'spending_ranking'
                ? 'bg-white text-gray-900 shadow-sm dark:bg-dark-700 dark:text-white'
                : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
            "
            @click="activeView = 'spending_ranking'"
          >
            {{ t('admin.dashboard.viewSpendingRanking') placeholderplaceholder
          </button>
        </div>
      </div>
    </div>

    <div v-if="activeView === 'model_distribution' && loading" class="flex h-48 items-center justify-center">
      <LoadingSpinner />
    </div>
    <div
      v-else-if="activeView === 'model_distribution' && displayModelStats.length > 0 && chartData"
      class="flex items-center gap-6"
    >
      <div class="h-48 w-48">
        <Doughnut :data="chartData" :options="doughnutOptions" />
      </div>
      <div class="max-h-64 flex-1 overflow-y-auto">
        <table class="w-full text-xs">
          <thead>
            <tr class="text-gray-500 dark:text-gray-400">
              <th class="pb-2 text-left">{{ t('admin.dashboard.model') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.requests') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.tokens') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.actual') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.standard') placeholderplaceholder</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="model in displayModelStats" :key="model.model">
              <tr
                class="border-t border-gray-100 cursor-pointer transition-colors hover:bg-gray-50 dark:border-gray-700 dark:hover:bg-dark-700/40"
                @click="toggleBreakdown('model', model.model)"
              >
                <td
                  class="max-w-[100px] truncate py-1.5 font-medium text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300"
                  :title="model.model"
                >
                  <span class="inline-flex items-center gap-1">
                    <svg v-if="expandedKey === `model-${model.modelplaceholder`" class="h-3 w-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
                    <svg v-else class="h-3 w-3 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
                    {{ model.model placeholderplaceholder
                  </span>
                </td>
                <td class="py-1.5 text-right text-gray-600 dark:text-gray-400">
                  {{ formatNumber(model.requests) placeholderplaceholder
                </td>
                <td class="py-1.5 text-right text-gray-600 dark:text-gray-400">
                  {{ formatTokens(model.total_tokens) placeholderplaceholder
                </td>
                <td class="py-1.5 text-right text-green-600 dark:text-green-400">
                  ${{ formatCost(model.actual_cost) placeholderplaceholder
                </td>
                <td class="py-1.5 text-right text-gray-400 dark:text-gray-500">
                  ${{ formatCost(model.cost) placeholderplaceholder
                </td>
              </tr>
              <tr v-if="expandedKey === `model-${model.modelplaceholder`">
                <td colspan="5" class="p-0">
                  <UserBreakdownSubTable
                    :items="breakdownItems"
                    :loading="breakdownLoading"
                  />
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>
    <div
      v-else-if="activeView === 'model_distribution'"
      class="flex h-48 items-center justify-center text-sm text-gray-500 dark:text-gray-400"
    >
      {{ t('admin.dashboard.noDataAvailable') placeholderplaceholder
    </div>

    <div v-else-if="rankingLoading" class="flex h-48 items-center justify-center">
      <LoadingSpinner />
    </div>
    <div
      v-else-if="rankingError"
      class="flex h-48 items-center justify-center text-sm text-gray-500 dark:text-gray-400"
    >
      {{ t('admin.dashboard.failedToLoad') placeholderplaceholder
    </div>
    <div v-else-if="rankingDisplayItems.length > 0 && rankingChartData" class="flex items-center gap-6">
      <div class="h-48 w-48">
        <Doughnut :data="rankingChartData" :options="rankingDoughnutOptions" />
      </div>
      <div class="max-h-48 flex-1 overflow-y-auto">
        <table class="w-full text-xs">
          <thead>
            <tr class="text-gray-500 dark:text-gray-400">
              <th class="pb-2 text-left">{{ t('admin.dashboard.spendingRankingUser') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.spendingRankingRequests') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.spendingRankingTokens') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.spendingRankingSpend') placeholderplaceholder</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in rankingDisplayItems"
              :key="item.isOther ? 'others' : `${item.user_idplaceholder-${indexplaceholder`"
              class="border-t border-gray-100 transition-colors dark:border-gray-700"
              :class="item.isOther
                ? 'bg-gray-50/70 dark:bg-dark-700/20'
                : 'cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-700/40'"
              @click="item.isOther ? undefined : emit('ranking-click', item)"
            >
              <td class="py-1.5">
                <div class="flex min-w-0 items-center gap-2">
                  <span class="shrink-0 text-[11px] font-semibold text-gray-500 dark:text-gray-400">
                    {{ item.isOther ? 'Σ' : `#${index + 1placeholder` placeholderplaceholder
                  </span>
                  <span
                    class="block max-w-[140px] truncate font-medium text-gray-900 dark:text-white"
                    :title="getRankingRowLabel(item)"
                  >
                    {{ getRankingRowLabel(item) placeholderplaceholder
                  </span>
                </div>
              </td>
              <td class="py-1.5 text-right text-gray-600 dark:text-gray-400">
                {{ formatNumber(item.requests) placeholderplaceholder
              </td>
              <td class="py-1.5 text-right text-gray-600 dark:text-gray-400">
                {{ formatTokens(item.tokens) placeholderplaceholder
              </td>
              <td class="py-1.5 text-right text-green-600 dark:text-green-400">
                ${{ formatCost(item.actual_cost) placeholderplaceholder
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div
      v-else
      class="flex h-48 items-center justify-center text-sm text-gray-500 dark:text-gray-400"
    >
      {{ t('admin.dashboard.noDataAvailable') placeholderplaceholder
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { Chart as ChartJS, ArcElement, Tooltip, Legend placeholder from 'chart.js'
import { Doughnut placeholder from 'vue-chartjs'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import UserBreakdownSubTable from './UserBreakdownSubTable.vue'
import type { ModelStat, UserSpendingRankingItem, UserBreakdownItem placeholder from '@/types'
import { getUserBreakdown placeholder from '@/api/admin/dashboard'

ChartJS.register(ArcElement, Tooltip, Legend)

const { t placeholder = useI18n()

type DistributionMetric = 'tokens' | 'actual_cost'
type RankingDisplayItem = UserSpendingRankingItem & { isOther?: boolean placeholder
const props = withDefaults(defineProps<{
  modelStats: ModelStat[]
  enableRankingView?: boolean
  rankingItems?: UserSpendingRankingItem[]
  rankingTotalActualCost?: number
  rankingTotalRequests?: number
  rankingTotalTokens?: number
  loading?: boolean
  metric?: DistributionMetric
  showMetricToggle?: boolean
  rankingLoading?: boolean
  rankingError?: boolean
  startDate?: string
  endDate?: string
placeholder>(), {
  enableRankingView: false,
  rankingItems: () => [],
  rankingTotalActualCost: 0,
  rankingTotalRequests: 0,
  rankingTotalTokens: 0,
  loading: false,
  metric: 'tokens',
  showMetricToggle: false,
  rankingLoading: false,
  rankingError: false
placeholder)

const expandedKey = ref<string | null>(null)
const breakdownItems = ref<UserBreakdownItem[]>([])
const breakdownLoading = ref(false)

const toggleBreakdown = async (type: string, id: string) => {
  const key = `${typeplaceholder-${idplaceholder`
  if (expandedKey.value === key) {
    expandedKey.value = null
    return
  placeholder
  expandedKey.value = key
  breakdownLoading.value = true
  breakdownItems.value = []
  try {
    const res = await getUserBreakdown({
      start_date: props.startDate,
      end_date: props.endDate,
      model: id,
    placeholder)
    breakdownItems.value = res.users || []
  placeholder catch {
    breakdownItems.value = []
  placeholder finally {
    breakdownLoading.value = false
  placeholder
placeholder

const emit = defineEmits<{
  'update:metric': [value: DistributionMetric]
  'ranking-click': [item: UserSpendingRankingItem]
placeholder>()

const enableRankingView = computed(() => props.enableRankingView)
const activeView = ref<'model_distribution' | 'spending_ranking'>('model_distribution')

const chartColors = [
  '#3b82f6',
  '#10b981',
  '#f59e0b',
  '#ef4444',
  '#8b5cf6',
  '#ec4899',
  '#14b8a6',
  '#f97316',
  '#6366f1',
  '#84cc16',
  '#06b6d4',
  '#a855f7'
]

const displayModelStats = computed(() => {
  if (!props.modelStats?.length) return []

  const metricKey = props.metric === 'actual_cost' ? 'actual_cost' : 'total_tokens'
  return [...props.modelStats].sort((a, b) => b[metricKey] - a[metricKey])
placeholder)

const chartData = computed(() => {
  if (!props.modelStats?.length) return null

  return {
    labels: displayModelStats.value.map((m) => m.model),
    datasets: [
      {
        data: displayModelStats.value.map((m) => props.metric === 'actual_cost' ? m.actual_cost : m.total_tokens),
        backgroundColor: chartColors.slice(0, displayModelStats.value.length),
        borderWidth: 0
      placeholder
    ]
  placeholder
placeholder)

const rankingChartData = computed(() => {
  if (!props.rankingItems?.length) return null

  const labels = props.rankingItems.map((item, index) => `#${index + 1placeholder ${getRankingUserLabel(item)placeholder`)
  const data = props.rankingItems.map((item) => item.actual_cost)
  const backgroundColor = chartColors.slice(0, props.rankingItems.length)

  if (otherRankingItem.value) {
    labels.push(t('admin.dashboard.spendingRankingOther'))
    data.push(otherRankingItem.value.actual_cost)
    backgroundColor.push('#94a3b8')
  placeholder

  return {
    labels,
    datasets: [
      {
        data,
        backgroundColor,
        borderWidth: 0
      placeholder
    ]
  placeholder
placeholder)

const otherRankingItem = computed<RankingDisplayItem | null>(() => {
  if (!props.rankingItems?.length) return null

  const rankedActualCost = props.rankingItems.reduce((sum, item) => sum + item.actual_cost, 0)
  const rankedRequests = props.rankingItems.reduce((sum, item) => sum + item.requests, 0)
  const rankedTokens = props.rankingItems.reduce((sum, item) => sum + item.tokens, 0)

  const otherActualCost = Math.max((props.rankingTotalActualCost || 0) - rankedActualCost, 0)
  const otherRequests = Math.max((props.rankingTotalRequests || 0) - rankedRequests, 0)
  const otherTokens = Math.max((props.rankingTotalTokens || 0) - rankedTokens, 0)

  if (otherActualCost <= 0.000001 && otherRequests <= 0 && otherTokens <= 0) return null

  return {
    user_id: 0,
    email: '',
    actual_cost: otherActualCost,
    requests: otherRequests,
    tokens: otherTokens,
    isOther: true
  placeholder
placeholder)

const rankingDisplayItems = computed<RankingDisplayItem[]>(() => {
  if (!props.rankingItems?.length) return []
  return otherRankingItem.value
    ? [...props.rankingItems, otherRankingItem.value]
    : [...props.rankingItems]
placeholder)

const doughnutOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    placeholder,
    tooltip: {
      callbacks: {
        label: (context: any) => {
          const value = context.raw as number
          const total = context.dataset.data.reduce((a: number, b: number) => a + b, 0)
          const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0'
          const formattedValue = props.metric === 'actual_cost'
            ? `$${formatCost(value)placeholder`
            : formatTokens(value)
          return `${context.labelplaceholder: ${formattedValueplaceholder (${percentageplaceholder%)`
        placeholder
      placeholder
    placeholder
  placeholder
placeholder))

const rankingDoughnutOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    placeholder,
    tooltip: {
      callbacks: {
        label: (context: any) => {
          const value = context.raw as number
          const total = context.dataset.data.reduce((a: number, b: number) => a + b, 0)
          const percentage = total > 0 ? ((value / total) * 100).toFixed(1) : '0.0'
          return `${context.labelplaceholder: $${formatCost(value)placeholder (${percentageplaceholder%)`
        placeholder
      placeholder
    placeholder
  placeholder
placeholder))

const formatTokens = (value: number): string => {
  if (value >= 1_000_000_000) {
    return `${(value / 1_000_000_000).toFixed(2)placeholderB`
  placeholder else if (value >= 1_000_000) {
    return `${(value / 1_000_000).toFixed(2)placeholderM`
  placeholder else if (value >= 1_000) {
    return `${(value / 1_000).toFixed(2)placeholderK`
  placeholder
  return value.toLocaleString()
placeholder

const formatNumber = (value: number): string => {
  return value.toLocaleString()
placeholder

const getRankingUserLabel = (item: UserSpendingRankingItem): string => {
  if (item.email) return item.email
  return t('admin.redeem.userPrefix', { id: item.user_id placeholder)
placeholder

const getRankingRowLabel = (item: RankingDisplayItem): string => {
  if (item.isOther) return t('admin.dashboard.spendingRankingOther')
  return getRankingUserLabel(item)
placeholder

const formatCost = (value: number): string => {
  if (value >= 1000) {
    return (value / 1000).toFixed(2) + 'K'
  placeholder else if (value >= 1) {
    return value.toFixed(2)
  placeholder else if (value >= 0.01) {
    return value.toFixed(3)
  placeholder
  return value.toFixed(4)
placeholder
</script>
