<template>
  <div class="card p-4">
    <div class="mb-4 flex items-center justify-between gap-3">
      <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
        {{ t('admin.dashboard.groupDistribution') placeholderplaceholder
      </h3>
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
    </div>
    <div v-if="loading" class="flex h-48 items-center justify-center">
      <LoadingSpinner />
    </div>
    <div v-else-if="displayGroupStats.length > 0 && chartData" class="flex items-center gap-6">
      <div class="h-48 w-48">
        <Doughnut :data="chartData" :options="doughnutOptions" />
      </div>
      <div class="max-h-48 flex-1 overflow-y-auto">
        <table class="w-full text-xs">
          <thead>
            <tr class="text-gray-500 dark:text-gray-400">
              <th class="pb-2 text-left">{{ t('admin.dashboard.group') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.requests') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.tokens') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.actual') placeholderplaceholder</th>
              <th class="pb-2 text-right">{{ t('admin.dashboard.standard') placeholderplaceholder</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="group in displayGroupStats"
              :key="group.group_id"
              class="border-t border-gray-100 dark:border-gray-700"
            >
              <td
                class="max-w-[100px] truncate py-1.5 font-medium text-gray-900 dark:text-white"
                :title="group.group_name || String(group.group_id)"
              >
                {{ group.group_name || t('admin.dashboard.noGroup') placeholderplaceholder
              </td>
              <td class="py-1.5 text-right text-gray-600 dark:text-gray-400">
                {{ formatNumber(group.requests) placeholderplaceholder
              </td>
              <td class="py-1.5 text-right text-gray-600 dark:text-gray-400">
                {{ formatTokens(group.total_tokens) placeholderplaceholder
              </td>
              <td class="py-1.5 text-right text-green-600 dark:text-green-400">
                ${{ formatCost(group.actual_cost) placeholderplaceholder
              </td>
              <td class="py-1.5 text-right text-gray-400 dark:text-gray-500">
                ${{ formatCost(group.cost) placeholderplaceholder
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
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { Chart as ChartJS, ArcElement, Tooltip, Legend placeholder from 'chart.js'
import { Doughnut placeholder from 'vue-chartjs'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { GroupStat placeholder from '@/types'

ChartJS.register(ArcElement, Tooltip, Legend)

const { t placeholder = useI18n()

type DistributionMetric = 'tokens' | 'actual_cost'

const props = withDefaults(defineProps<{
  groupStats: GroupStat[]
  loading?: boolean
  metric?: DistributionMetric
  showMetricToggle?: boolean
placeholder>(), {
  loading: false,
  metric: 'tokens',
  showMetricToggle: false,
placeholder)

const emit = defineEmits<{
  'update:metric': [value: DistributionMetric]
placeholder>()

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
  '#84cc16'
]

const displayGroupStats = computed(() => {
  if (!props.groupStats?.length) return []

  const metricKey = props.metric === 'actual_cost' ? 'actual_cost' : 'total_tokens'
  return [...props.groupStats].sort((a, b) => b[metricKey] - a[metricKey])
placeholder)

const chartData = computed(() => {
  if (!props.groupStats?.length) return null

  return {
    labels: displayGroupStats.value.map((g) => g.group_name || String(g.group_id)),
    datasets: [
      {
        data: displayGroupStats.value.map((g) => props.metric === 'actual_cost' ? g.actual_cost : g.total_tokens),
        backgroundColor: chartColors.slice(0, displayGroupStats.value.length),
        borderWidth: 0
      placeholder
    ]
  placeholder
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
