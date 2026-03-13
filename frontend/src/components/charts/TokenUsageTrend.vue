<template>
  <div class="card p-4">
    <h3 class="mb-4 text-sm font-semibold text-gray-900 dark:text-white">
      {{ t('admin.dashboard.tokenUsageTrend') placeholderplaceholder
    </h3>
    <div v-if="loading" class="flex h-48 items-center justify-center">
      <LoadingSpinner />
    </div>
    <div v-else-if="trendData.length > 0 && chartData" class="h-48">
      <Line :data="chartData" :options="lineOptions" />
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
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
placeholder from 'chart.js'
import { Line placeholder from 'vue-chartjs'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { TrendDataPoint placeholder from '@/types'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const { t placeholder = useI18n()

const props = defineProps<{
  trendData: TrendDataPoint[]
  loading?: boolean
placeholder>()

const isDarkMode = computed(() => {
  return document.documentElement.classList.contains('dark')
placeholder)

const chartColors = computed(() => ({
  text: isDarkMode.value ? '#e5e7eb' : '#374151',
  grid: isDarkMode.value ? '#374151' : '#e5e7eb',
  input: '#3b82f6',
  output: '#10b981',
  cacheCreation: '#f59e0b',
  cacheRead: '#06b6d4'
placeholder))

const chartData = computed(() => {
  if (!props.trendData?.length) return null

  return {
    labels: props.trendData.map((d) => d.date),
    datasets: [
      {
        label: 'Input',
        data: props.trendData.map((d) => d.input_tokens),
        borderColor: chartColors.value.input,
        backgroundColor: `${chartColors.value.inputplaceholder20`,
        fill: true,
        tension: 0.3
      placeholder,
      {
        label: 'Output',
        data: props.trendData.map((d) => d.output_tokens),
        borderColor: chartColors.value.output,
        backgroundColor: `${chartColors.value.outputplaceholder20`,
        fill: true,
        tension: 0.3
      placeholder,
      {
        label: 'Cache Creation',
        data: props.trendData.map((d) => d.cache_creation_tokens),
        borderColor: chartColors.value.cacheCreation,
        backgroundColor: `${chartColors.value.cacheCreationplaceholder20`,
        fill: true,
        tension: 0.3
      placeholder,
      {
        label: 'Cache Read',
        data: props.trendData.map((d) => d.cache_read_tokens),
        borderColor: chartColors.value.cacheRead,
        backgroundColor: `${chartColors.value.cacheReadplaceholder20`,
        fill: true,
        tension: 0.3
      placeholder
    ]
  placeholder
placeholder)

const lineOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index' as const
  placeholder,
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        color: chartColors.value.text,
        usePointStyle: true,
        pointStyle: 'circle',
        padding: 15,
        font: {
          size: 11
        placeholder
      placeholder
    placeholder,
    tooltip: {
      callbacks: {
        label: (context: any) => {
          return `${context.dataset.labelplaceholder: ${formatTokens(context.raw)placeholder`
        placeholder,
        footer: (tooltipItems: any) => {
          const dataIndex = tooltipItems[0]?.dataIndex
          if (dataIndex !== undefined && props.trendData[dataIndex]) {
            const data = props.trendData[dataIndex]
            return `Actual: $${formatCost(data.actual_cost)placeholder | Standard: $${formatCost(data.cost)placeholder`
          placeholder
          return ''
        placeholder
      placeholder
    placeholder
  placeholder,
  scales: {
    x: {
      grid: {
        color: chartColors.value.grid
      placeholder,
      ticks: {
        color: chartColors.value.text,
        font: {
          size: 10
        placeholder
      placeholder
    placeholder,
    y: {
      grid: {
        color: chartColors.value.grid
      placeholder,
      ticks: {
        color: chartColors.value.text,
        font: {
          size: 10
        placeholder,
        callback: (value: string | number) => formatTokens(Number(value))
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
