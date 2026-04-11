<template>
  <div class="card p-4">
    <h3 class="mb-4 text-sm font-semibold text-gray-900 dark:text-white">
      {{ t('payment.admin.dailyRevenue') placeholderplaceholder
    </h3>
    <div class="h-64">
      <div v-if="loading" class="flex h-full items-center justify-center">
        <LoadingSpinner size="md" />
      </div>
      <Line v-else-if="chartData" :data="chartData" :options="chartOptions" />
      <div
        v-else
        class="flex h-full items-center justify-center text-sm text-gray-500 dark:text-gray-400"
      >
        {{ t('payment.admin.noData') placeholderplaceholder
      </div>
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
  Tooltip,
  Legend,
  Filler
placeholder from 'chart.js'
import { Line placeholder from 'vue-chartjs'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend, Filler)

const { t placeholder = useI18n()

const props = defineProps<{
  data: { date: string; amount: number; count: number placeholder[]
  loading?: boolean
placeholder>()

const chartData = computed(() => {
  if (!props.data || props.data.length === 0) return null
  return {
    labels: props.data.map(d => d.date),
    datasets: [
      {
        label: t('payment.admin.revenue'),
        data: props.data.map(d => d.amount),
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        fill: true,
        tension: 0.3,
        pointRadius: 3,
        pointHoverRadius: 5,
      placeholder,
      {
        label: t('payment.admin.orderCount'),
        data: props.data.map(d => d.count),
        borderColor: 'rgb(16, 185, 129)',
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        fill: false,
        tension: 0.3,
        pointRadius: 3,
        pointHoverRadius: 5,
        yAxisID: 'y1',
      placeholder
    ]
  placeholder
placeholder)

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: { mode: 'index' as const, intersect: false placeholder,
  scales: {
    y: {
      type: 'linear' as const,
      display: true,
      position: 'left' as const,
      title: { display: true, text: t('payment.admin.revenue') placeholder,
    placeholder,
    y1: {
      type: 'linear' as const,
      display: true,
      position: 'right' as const,
      title: { display: true, text: t('payment.admin.orderCount') placeholder,
      grid: { drawOnChartArea: false placeholder,
    placeholder
  placeholder,
  plugins: {
    legend: { position: 'top' as const placeholder,
  placeholder
placeholder
</script>
