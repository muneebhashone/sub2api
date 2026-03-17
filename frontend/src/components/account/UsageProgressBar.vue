<template>
  <div>
    <!-- Window stats row (above progress bar) -->
    <div
      v-if="windowStats && (windowStats.requests > 0 || windowStats.tokens > 0)"
      class="mb-0.5 flex items-center"
    >
      <div class="flex items-center gap-1.5 text-[9px] text-gray-500 dark:text-gray-400">
        <span class="rounded bg-gray-100 px-1.5 py-0.5 dark:bg-gray-800">
          {{ formatRequests placeholderplaceholder req
        </span>
        <span class="rounded bg-gray-100 px-1.5 py-0.5 dark:bg-gray-800">
          {{ formatTokens placeholderplaceholder
        </span>
        <span class="rounded bg-gray-100 px-1.5 py-0.5 dark:bg-gray-800" :title="t('usage.accountBilled')">
          A ${{ formatAccountCost placeholderplaceholder
        </span>
        <span
          v-if="windowStats?.user_cost != null"
          class="rounded bg-gray-100 px-1.5 py-0.5 dark:bg-gray-800"
          :title="t('usage.userBilled')"
        >
          U ${{ formatUserCost placeholderplaceholder
        </span>
      </div>
    </div>

    <!-- Progress bar row -->
    <div class="flex items-center gap-1">
      <!-- Label badge (fixed width for alignment) -->
      <span
        :class="['w-[32px] shrink-0 rounded px-1 text-center text-[10px] font-medium', labelClass]"
      >
        {{ label placeholderplaceholder
      </span>

      <!-- Progress bar container -->
      <div class="h-1.5 w-8 shrink-0 overflow-hidden rounded-full bg-gray-200 dark:bg-gray-700">
        <div
          :class="['h-full transition-all duration-300', barClass]"
          :style="{ width: barWidth placeholder"
        ></div>
      </div>

      <!-- Percentage -->
      <span :class="['w-[32px] shrink-0 text-right text-[10px] font-medium', textClass]">
        {{ displayPercent placeholderplaceholder
      </span>

      <!-- Reset time -->
      <span v-if="shouldShowResetTime" class="shrink-0 text-[10px] text-gray-400">
        {{ formatResetTime placeholderplaceholder
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch placeholder from 'vue'
import { useIntervalFn placeholder from '@vueuse/core'
import { useI18n placeholder from 'vue-i18n'
import type { WindowStats placeholder from '@/types'
import { formatCompactNumber placeholder from '@/utils/format'

const props = defineProps<{
  label: string
  utilization: number // Percentage (0-100+)
  resetsAt?: string | null
  color: 'indigo' | 'emerald' | 'purple' | 'amber'
  windowStats?: WindowStats | null
  showNowWhenIdle?: boolean
placeholder>()

const { t placeholder = useI18n()

// Reactive clock for countdown — only runs when a reset time is shown,
// to avoid creating many idle timers across large account lists.
const now = ref(new Date())
const { pause: pauseClock, resume: resumeClock placeholder = useIntervalFn(
  () => {
    now.value = new Date()
  placeholder,
  60_000,
  { immediate: false placeholder,
)
if (props.resetsAt) resumeClock()
watch(
  () => props.resetsAt,
  (val) => {
    if (val) {
      now.value = new Date()
      resumeClock()
    placeholder else {
      pauseClock()
    placeholder
  placeholder,
)

// Label background colors
const labelClass = computed(() => {
  const colors = {
    indigo: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300',
    emerald: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300',
    purple: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-300',
    amber: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300'
  placeholder
  return colors[props.color]
placeholder)

// Progress bar color based on utilization
const barClass = computed(() => {
  if (props.utilization >= 100) {
    return 'bg-red-500'
  placeholder else if (props.utilization >= 80) {
    return 'bg-amber-500'
  placeholder else {
    return 'bg-green-500'
  placeholder
placeholder)

// Text color based on utilization
const textClass = computed(() => {
  if (props.utilization >= 100) {
    return 'text-red-600 dark:text-red-400'
  placeholder else if (props.utilization >= 80) {
    return 'text-amber-600 dark:text-amber-400'
  placeholder else {
    return 'text-gray-600 dark:text-gray-400'
  placeholder
placeholder)

// Bar width (capped at 100%)
const barWidth = computed(() => {
  return `${Math.min(props.utilization, 100)placeholder%`
placeholder)

// Display percentage (cap at 999% for readability)
const displayPercent = computed(() => {
  const percent = Math.round(props.utilization)
  return percent > 999 ? '>999%' : `${percentplaceholder%`
placeholder)

const shouldShowResetTime = computed(() => {
  if (props.resetsAt) return true
  return Boolean(props.showNowWhenIdle && props.utilization <= 0)
placeholder)

// Format reset time
const formatResetTime = computed(() => {
  // For rolling windows, when utilization is 0%, treat as immediately available.
  if (props.showNowWhenIdle && props.utilization <= 0) {
    return '现在'
  placeholder

  if (!props.resetsAt) return '-'

  const date = new Date(props.resetsAt)
  const diffMs = date.getTime() - now.value.getTime()

  if (diffMs <= 0) return '现在'

  const diffHours = Math.floor(diffMs / (1000 * 60 * 60))
  const diffMins = Math.floor((diffMs % (1000 * 60 * 60)) / (1000 * 60))

  if (diffHours >= 24) {
    const days = Math.floor(diffHours / 24)
    return `${daysplaceholderd ${diffHours % 24placeholderh`
  placeholder else if (diffHours > 0) {
    return `${diffHoursplaceholderh ${diffMinsplaceholderm`
  placeholder else {
    return `${diffMinsplaceholderm`
  placeholder
placeholder)

// Window stats formatters
const formatRequests = computed(() => {
  if (!props.windowStats) return ''
  return formatCompactNumber(props.windowStats.requests, { allowBillions: false placeholder)
placeholder)

const formatTokens = computed(() => {
  if (!props.windowStats) return ''
  return formatCompactNumber(props.windowStats.tokens)
placeholder)

const formatAccountCost = computed(() => {
  if (!props.windowStats) return '0.00'
  return props.windowStats.cost.toFixed(2)
placeholder)

const formatUserCost = computed(() => {
  if (!props.windowStats || props.windowStats.user_cost == null) return '0.00'
  return props.windowStats.user_cost.toFixed(2)
placeholder)

</script>
