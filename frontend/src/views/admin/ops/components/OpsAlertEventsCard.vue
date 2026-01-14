<script setup lang="ts">
import { computed, onMounted, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import Select from '@/components/common/Select.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import { opsAPI, type AlertEventsQuery placeholder from '@/api/admin/ops'
import type { AlertEvent placeholder from '../types'
import { formatDateTime placeholder from '../utils/opsFormatters'

const { t placeholder = useI18n()
const appStore = useAppStore()

const PAGE_SIZE = 20

const loading = ref(false)
const loadingMore = ref(false)
const events = ref<AlertEvent[]>([])
const hasMore = ref(true)

// Detail modal
const showDetail = ref(false)
const selected = ref<AlertEvent | null>(null)
const detailLoading = ref(false)
const detailActionLoading = ref(false)
const historyLoading = ref(false)
const history = ref<AlertEvent[]>([])
const historyRange = ref('7d')
const historyRangeOptions = computed(() => [
  { value: '7d', label: t('admin.ops.timeRange.7d') placeholder,
  { value: '30d', label: t('admin.ops.timeRange.30d') placeholder
])

const silenceDuration = ref('1h')
const silenceDurationOptions = computed(() => [
  { value: '1h', label: t('admin.ops.timeRange.1h') placeholder,
  { value: '24h', label: t('admin.ops.timeRange.24h') placeholder,
  { value: '7d', label: t('admin.ops.timeRange.7d') placeholder
])

// Filters
const timeRange = ref('24h')
const timeRangeOptions = computed(() => [
  { value: '5m', label: t('admin.ops.timeRange.5m') placeholder,
  { value: '30m', label: t('admin.ops.timeRange.30m') placeholder,
  { value: '1h', label: t('admin.ops.timeRange.1h') placeholder,
  { value: '6h', label: t('admin.ops.timeRange.6h') placeholder,
  { value: '24h', label: t('admin.ops.timeRange.24h') placeholder,
  { value: '7d', label: t('admin.ops.timeRange.7d') placeholder,
  { value: '30d', label: t('admin.ops.timeRange.30d') placeholder
])

const severity = ref<string>('')
const severityOptions = computed(() => [
  { value: '', label: t('common.all') placeholder,
  { value: 'P0', label: 'P0' placeholder,
  { value: 'P1', label: 'P1' placeholder,
  { value: 'P2', label: 'P2' placeholder,
  { value: 'P3', label: 'P3' placeholder
])

const status = ref<string>('')
const statusOptions = computed(() => [
  { value: '', label: t('common.all') placeholder,
  { value: 'firing', label: t('admin.ops.alertEvents.status.firing') placeholder,
  { value: 'resolved', label: t('admin.ops.alertEvents.status.resolved') placeholder,
  { value: 'manual_resolved', label: t('admin.ops.alertEvents.status.manualResolved') placeholder
])

const emailSent = ref<string>('')
const emailSentOptions = computed(() => [
  { value: '', label: t('common.all') placeholder,
  { value: 'true', label: t('admin.ops.alertEvents.table.emailSent') placeholder,
  { value: 'false', label: t('admin.ops.alertEvents.table.emailIgnored') placeholder
])

function buildQuery(overrides: Partial<AlertEventsQuery> = {placeholder): AlertEventsQuery {
  const q: AlertEventsQuery = {
    limit: PAGE_SIZE,
    time_range: timeRange.value
  placeholder
  if (severity.value) q.severity = severity.value
  if (status.value) q.status = status.value
  if (emailSent.value === 'true') q.email_sent = true
  if (emailSent.value === 'false') q.email_sent = false
  return { ...q, ...overrides placeholder
placeholder

async function loadFirstPage() {
  loading.value = true
  try {
    const data = await opsAPI.listAlertEvents(buildQuery())
    events.value = data
    hasMore.value = data.length === PAGE_SIZE
  placeholder catch (err: any) {
    console.error('[OpsAlertEventsCard] Failed to load alert events', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertEvents.loadFailed'))
    events.value = []
    hasMore.value = false
  placeholder finally {
    loading.value = false
  placeholder
placeholder

async function loadMore() {
  if (loadingMore.value || loading.value) return
  if (!hasMore.value) return
  const last = events.value[events.value.length - 1]
  if (!last) return

  loadingMore.value = true
  try {
    const data = await opsAPI.listAlertEvents(
      buildQuery({ before_fired_at: last.fired_at || last.created_at, before_id: last.id placeholder)
    )
    if (!data.length) {
      hasMore.value = false
      return
    placeholder
    events.value = [...events.value, ...data]
    if (data.length < PAGE_SIZE) hasMore.value = false
  placeholder catch (err: any) {
    console.error('[OpsAlertEventsCard] Failed to load more alert events', err)
    hasMore.value = false
  placeholder finally {
    loadingMore.value = false
  placeholder
placeholder

function onScroll(e: Event) {
  const el = e.target as HTMLElement | null
  if (!el) return
  const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 120
  if (nearBottom) loadMore()
placeholder

function getDimensionString(event: AlertEvent | null | undefined, key: string): string {
  const v = event?.dimensions?.[key]
  if (v == null) return ''
  if (typeof v === 'string') return v
  if (typeof v === 'number' || typeof v === 'boolean') return String(v)
  return ''
placeholder

function formatDurationMs(ms: number): string {
  const safe = Math.max(0, Math.floor(ms))
  const sec = Math.floor(safe / 1000)
  if (sec < 60) return `${secplaceholders`
  const min = Math.floor(sec / 60)
  if (min < 60) return `${minplaceholderm`
  const hr = Math.floor(min / 60)
  if (hr < 24) return `${hrplaceholderh`
  const day = Math.floor(hr / 24)
  return `${dayplaceholderd`
placeholder

function formatDurationLabel(event: AlertEvent): string {
  const firedAt = new Date(event.fired_at || event.created_at)
  if (Number.isNaN(firedAt.getTime())) return '-'
  const resolvedAtStr = event.resolved_at || null
  const status = String(event.status || '').trim().toLowerCase()

  if (resolvedAtStr) {
    const resolvedAt = new Date(resolvedAtStr)
    if (!Number.isNaN(resolvedAt.getTime())) {
      const ms = resolvedAt.getTime() - firedAt.getTime()
      const prefix = status === 'manual_resolved'
        ? t('admin.ops.alertEvents.status.manualResolved')
        : t('admin.ops.alertEvents.status.resolved')
      return `${prefixplaceholder ${formatDurationMs(ms)placeholder`
    placeholder
  placeholder

  const now = Date.now()
  const ms = now - firedAt.getTime()
  return `${t('admin.ops.alertEvents.status.firing')placeholder ${formatDurationMs(ms)placeholder`
placeholder

function formatDimensionsSummary(event: AlertEvent): string {
  const parts: string[] = []
  const platform = getDimensionString(event, 'platform')
  if (platform) parts.push(`platform=${platformplaceholder`)
  const groupId = event.dimensions?.group_id
  if (groupId != null && groupId !== '') parts.push(`group_id=${String(groupId)placeholder`)
  const region = getDimensionString(event, 'region')
  if (region) parts.push(`region=${regionplaceholder`)
  return parts.length ? parts.join(' ') : '-'
placeholder

function closeDetail() {
  showDetail.value = false
  selected.value = null
  history.value = []
placeholder

async function openDetail(row: AlertEvent) {
  showDetail.value = true
  selected.value = row
  detailLoading.value = true
  historyLoading.value = true

  try {
    const detail = await opsAPI.getAlertEvent(row.id)
    selected.value = detail
  placeholder catch (err: any) {
    console.error('[OpsAlertEventsCard] Failed to load alert detail', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertEvents.detail.loadFailed'))
  placeholder finally {
    detailLoading.value = false
  placeholder

  await loadHistory()
placeholder

async function loadHistory() {
  const ev = selected.value
  if (!ev) {
    history.value = []
    historyLoading.value = false
    return
  placeholder

  historyLoading.value = true
  try {
    const platform = getDimensionString(ev, 'platform')
    const groupIdRaw = ev.dimensions?.group_id
    const groupId = typeof groupIdRaw === 'number' ? groupIdRaw : undefined

    const items = await opsAPI.listAlertEvents({
      limit: 20,
      time_range: historyRange.value,
      platform: platform || undefined,
      group_id: groupId,
      status: ''
    placeholder)

    // Best-effort: narrow to same rule_id + dimensions
    history.value = items.filter((it) => {
      if (it.rule_id !== ev.rule_id) return false
      const p1 = getDimensionString(it, 'platform')
      const p2 = getDimensionString(ev, 'platform')
      if ((p1 || '') !== (p2 || '')) return false
      const g1 = it.dimensions?.group_id
      const g2 = ev.dimensions?.group_id
      return (g1 ?? null) === (g2 ?? null)
    placeholder)
  placeholder catch (err: any) {
    console.error('[OpsAlertEventsCard] Failed to load alert history', err)
    history.value = []
  placeholder finally {
    historyLoading.value = false
  placeholder
placeholder

function durationToUntilRFC3339(duration: string): string {
  const now = Date.now()
  if (duration === '1h') return new Date(now + 60 * 60 * 1000).toISOString()
  if (duration === '24h') return new Date(now + 24 * 60 * 60 * 1000).toISOString()
  if (duration === '7d') return new Date(now + 7 * 24 * 60 * 60 * 1000).toISOString()
  return new Date(now + 60 * 60 * 1000).toISOString()
placeholder

async function silenceAlert() {
  const ev = selected.value
  if (!ev) return
  if (detailActionLoading.value) return
  detailActionLoading.value = true
  try {
    const platform = getDimensionString(ev, 'platform')
    const groupIdRaw = ev.dimensions?.group_id
    const groupId = typeof groupIdRaw === 'number' ? groupIdRaw : null
    const region = getDimensionString(ev, 'region') || null

    await opsAPI.createAlertSilence({
      rule_id: ev.rule_id,
      platform: platform || '',
      group_id: groupId ?? undefined,
      region: region ?? undefined,
      until: durationToUntilRFC3339(silenceDuration.value),
      reason: `silence from UI (${silenceDuration.valueplaceholder)`
    placeholder)

    appStore.showSuccess(t('admin.ops.alertEvents.detail.silenceSuccess'))
  placeholder catch (err: any) {
    console.error('[OpsAlertEventsCard] Failed to silence alert', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertEvents.detail.silenceFailed'))
  placeholder finally {
    detailActionLoading.value = false
  placeholder
placeholder

async function manualResolve() {
  if (!selected.value) return
  if (detailActionLoading.value) return
  detailActionLoading.value = true
  try {
    await opsAPI.updateAlertEventStatus(selected.value.id, 'manual_resolved')
    appStore.showSuccess(t('admin.ops.alertEvents.detail.manualResolvedSuccess'))

    // Refresh detail + first page to reflect new status
    const detail = await opsAPI.getAlertEvent(selected.value.id)
    selected.value = detail
    await loadFirstPage()
    await loadHistory()
  placeholder catch (err: any) {
    console.error('[OpsAlertEventsCard] Failed to resolve alert', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertEvents.detail.manualResolvedFailed'))
  placeholder finally {
    detailActionLoading.value = false
  placeholder
placeholder

onMounted(() => {
  loadFirstPage()
placeholder)

watch([timeRange, severity, status, emailSent], () => {
  events.value = []
  hasMore.value = true
  loadFirstPage()
placeholder)

watch(historyRange, () => {
  if (showDetail.value) loadHistory()
placeholder)

function severityBadgeClass(severity: string | undefined): string {
  const s = String(severity || '').trim().toLowerCase()
  if (s === 'p0' || s === 'critical') return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
  if (s === 'p1' || s === 'warning') return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
  if (s === 'p2' || s === 'info') return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
  if (s === 'p3') return 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300'
  return 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300'
placeholder

function statusBadgeClass(status: string | undefined): string {
  const s = String(status || '').trim().toLowerCase()
  if (s === 'firing') return 'bg-red-50 text-red-700 ring-red-600/20 dark:bg-red-900/30 dark:text-red-300 dark:ring-red-500/30'
  if (s === 'resolved') return 'bg-green-50 text-green-700 ring-green-600/20 dark:bg-green-900/30 dark:text-green-300 dark:ring-green-500/30'
  if (s === 'manual_resolved') return 'bg-slate-50 text-slate-700 ring-slate-600/20 dark:bg-slate-900/30 dark:text-slate-300 dark:ring-slate-500/30'
  return 'bg-gray-50 text-gray-700 ring-gray-600/20 dark:bg-gray-900/30 dark:text-gray-300 dark:ring-gray-500/30'
placeholder

function formatStatusLabel(status: string | undefined): string {
  const s = String(status || '').trim().toLowerCase()
  if (!s) return '-'
  if (s === 'firing') return t('admin.ops.alertEvents.status.firing')
  if (s === 'resolved') return t('admin.ops.alertEvents.status.resolved')
  if (s === 'manual_resolved') return t('admin.ops.alertEvents.status.manualResolved')
  return s.toUpperCase()
placeholder

const empty = computed(() => events.value.length === 0 && !loading.value)
</script>

<template>
  <div class="rounded-3xl bg-white p-6 shadow-sm ring-1 ring-gray-900/5 dark:bg-dark-800 dark:ring-dark-700">
    <div class="mb-4 flex items-start justify-between gap-4">
      <div>
        <h3 class="text-sm font-bold text-gray-900 dark:text-white">{{ t('admin.ops.alertEvents.title') placeholderplaceholder</h3>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.ops.alertEvents.description') placeholderplaceholder</p>
      </div>

      <div class="flex items-center gap-2">
        <Select :model-value="timeRange" :options="timeRangeOptions" class="w-[120px]" @change="timeRange = String($event || '24h')" />
        <Select :model-value="severity" :options="severityOptions" class="w-[88px]" @change="severity = String($event || '')" />
        <Select :model-value="status" :options="statusOptions" class="w-[110px]" @change="status = String($event || '')" />
        <Select :model-value="emailSent" :options="emailSentOptions" class="w-[110px]" @change="emailSent = String($event || '')" />
        <button
          class="flex items-center gap-1.5 rounded-lg bg-gray-100 px-3 py-1.5 text-xs font-bold text-gray-700 transition-colors hover:bg-gray-200 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-dark-700 dark:text-gray-300 dark:hover:bg-dark-600"
          :disabled="loading"
          @click="loadFirstPage"
        >
          <svg class="h-3.5 w-3.5" :class="{ 'animate-spin': loading placeholder" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ t('common.refresh') placeholderplaceholder
        </button>
      </div>
    </div>

    <div v-if="loading" class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
      <svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      {{ t('admin.ops.alertEvents.loading') placeholderplaceholder
    </div>

    <div v-else-if="empty" class="rounded-xl border border-dashed border-gray-200 p-8 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400">
      {{ t('admin.ops.alertEvents.empty') placeholderplaceholder
    </div>

    <div v-else class="overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700">
      <div class="max-h-[600px] overflow-y-auto" @scroll="onScroll">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
          <thead class="sticky top-0 z-10 bg-gray-50 dark:bg-dark-900">
            <tr>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.time') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.severity') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.platform') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.ruleId') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.title') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.duration') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.dimensions') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-right text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertEvents.table.email') placeholderplaceholder
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white dark:divide-dark-700 dark:bg-dark-800">
            <tr
              v-for="row in events"
              :key="row.id"
              class="cursor-pointer hover:bg-gray-50 dark:hover:bg-dark-700/50"
              @click="openDetail(row)"
              :title="row.title || ''"
            >
              <td class="whitespace-nowrap px-4 py-3 text-xs text-gray-600 dark:text-gray-300">
                {{ formatDateTime(row.fired_at || row.created_at) placeholderplaceholder
              </td>
              <td class="whitespace-nowrap px-4 py-3">
                <div class="flex items-center gap-2">
                  <span class="rounded-full px-2 py-1 text-[10px] font-bold" :class="severityBadgeClass(String(row.severity || ''))">
                    {{ row.severity || '-' placeholderplaceholder
                  </span>
                  <span class="inline-flex items-center rounded-full px-2 py-1 text-[10px] font-bold ring-1 ring-inset" :class="statusBadgeClass(row.status)">
                    {{ formatStatusLabel(row.status) placeholderplaceholder
                  </span>
                </div>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-xs text-gray-600 dark:text-gray-300">
                {{ getDimensionString(row, 'platform') || '-' placeholderplaceholder
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-xs text-gray-600 dark:text-gray-300">
                <span class="font-mono">#{{ row.rule_id placeholderplaceholder</span>
              </td>
              <td class="min-w-[260px] px-4 py-3 text-xs text-gray-700 dark:text-gray-200">
                <div class="font-semibold truncate max-w-[360px]">{{ row.title || '-' placeholderplaceholder</div>
                <div v-if="row.description" class="mt-0.5 line-clamp-2 text-[11px] text-gray-500 dark:text-gray-400">
                  {{ row.description placeholderplaceholder
                </div>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-xs text-gray-600 dark:text-gray-300">
                {{ formatDurationLabel(row) placeholderplaceholder
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-[11px] text-gray-500 dark:text-gray-400">
                {{ formatDimensionsSummary(row) placeholderplaceholder
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-right text-xs">
                <span
                  class="inline-flex items-center justify-end gap-1.5"
                  :title="row.email_sent ? t('admin.ops.alertEvents.table.emailSent') : t('admin.ops.alertEvents.table.emailIgnored')"
                >
                  <Icon
                    v-if="row.email_sent"
                    name="checkCircle"
                    size="sm"
                    class="text-green-600 dark:text-green-400"
                  />
                  <Icon
                    v-else
                    name="ban"
                    size="sm"
                    class="text-gray-400 dark:text-gray-500"
                  />
                  <span class="text-[11px] font-bold text-gray-600 dark:text-gray-300">
                    {{ row.email_sent ? t('admin.ops.alertEvents.table.emailSent') : t('admin.ops.alertEvents.table.emailIgnored') placeholderplaceholder
                  </span>
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        <div v-if="loadingMore" class="flex items-center justify-center gap-2 py-3 text-xs text-gray-500 dark:text-gray-400">
          <svg class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ t('admin.ops.alertEvents.loading') placeholderplaceholder
        </div>
        <div v-else-if="!hasMore && events.length > 0" class="py-3 text-center text-xs text-gray-400">
          -
        </div>
      </div>
    </div>

    <BaseDialog
      :show="showDetail"
      :title="t('admin.ops.alertEvents.detail.title')"
      width="wide"
      :close-on-click-outside="true"
      @close="closeDetail"
    >
      <div v-if="detailLoading" class="flex items-center justify-center py-10 text-sm text-gray-500 dark:text-gray-400">
        {{ t('admin.ops.alertEvents.detail.loading') placeholderplaceholder
      </div>

      <div v-else-if="!selected" class="py-10 text-center text-sm text-gray-500 dark:text-gray-400">
        {{ t('admin.ops.alertEvents.detail.empty') placeholderplaceholder
      </div>

      <div v-else class="space-y-5">
        <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <span class="inline-flex items-center rounded-full px-2 py-1 text-[10px] font-bold" :class="severityBadgeClass(String(selected.severity || ''))">
                  {{ selected.severity || '-' placeholderplaceholder
                </span>
                <span class="inline-flex items-center rounded-full px-2 py-1 text-[10px] font-bold ring-1 ring-inset" :class="statusBadgeClass(selected.status)">
                  {{ formatStatusLabel(selected.status) placeholderplaceholder
                </span>
              </div>
              <div class="mt-2 text-sm font-semibold text-gray-900 dark:text-white">
                {{ selected.title || '-' placeholderplaceholder
              </div>
              <div v-if="selected.description" class="mt-1 whitespace-pre-wrap text-xs text-gray-600 dark:text-gray-300">
                {{ selected.description placeholderplaceholder
              </div>
            </div>

            <div class="flex flex-wrap gap-2">
              <div class="flex items-center gap-2 rounded-lg bg-white px-2 py-1 ring-1 ring-gray-200 dark:bg-dark-800 dark:ring-dark-700">
                <span class="text-[11px] font-bold text-gray-600 dark:text-gray-300">{{ t('admin.ops.alertEvents.detail.silence') placeholderplaceholder</span>
                <Select
                  :model-value="silenceDuration"
                  :options="silenceDurationOptions"
                  class="w-[110px]"
                  @change="silenceDuration = String($event || '1h')"
                />
                <button type="button" class="btn btn-secondary btn-sm" :disabled="detailActionLoading" @click="silenceAlert">
                  <Icon name="ban" size="sm" />
                  {{ t('common.apply') placeholderplaceholder
                </button>
              </div>

              <button type="button" class="btn btn-secondary btn-sm" :disabled="detailActionLoading" @click="manualResolve">
                <Icon name="checkCircle" size="sm" />
                {{ t('admin.ops.alertEvents.detail.manualResolve') placeholderplaceholder
              </button>
            </div>
          </div>
        </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.alertEvents.detail.firedAt') placeholderplaceholder</div>
              <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ formatDateTime(selected.fired_at || selected.created_at) placeholderplaceholder</div>
            </div>
            <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.alertEvents.detail.resolvedAt') placeholderplaceholder</div>
              <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ selected.resolved_at ? formatDateTime(selected.resolved_at) : '-' placeholderplaceholder</div>
            </div>
            <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.alertEvents.detail.ruleId') placeholderplaceholder</div>
              <div class="mt-1 flex flex-wrap items-center gap-2">
                <div class="font-mono text-sm font-bold text-gray-900 dark:text-white">#{{ selected.rule_id placeholderplaceholder</div>
                <a
                  class="inline-flex items-center gap-1 rounded-md bg-white px-2 py-1 text-[11px] font-bold text-gray-700 ring-1 ring-gray-200 hover:bg-gray-50 dark:bg-dark-800 dark:text-gray-200 dark:ring-dark-700 dark:hover:bg-dark-700"
                  :href="`/admin/ops?open_alert_rules=1&alert_rule_id=${selected.rule_idplaceholder`"
                >
                  <Icon name="externalLink" size="xs" />
                  {{ t('admin.ops.alertEvents.detail.viewRule') placeholderplaceholder
                </a>
                <a
                  class="inline-flex items-center gap-1 rounded-md bg-white px-2 py-1 text-[11px] font-bold text-gray-700 ring-1 ring-gray-200 hover:bg-gray-50 dark:bg-dark-800 dark:text-gray-200 dark:ring-dark-700 dark:hover:bg-dark-700"
                  :href="`/admin/ops?platform=${encodeURIComponent(getDimensionString(selected,'platform')||'')placeholder&group_id=${selected.dimensions?.group_id || ''placeholder&error_type=request&open_error_details=1`"
                >
                  <Icon name="externalLink" size="xs" />
                  {{ t('admin.ops.alertEvents.detail.viewLogs') placeholderplaceholder
                </a>
              </div>
            </div>
            <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.alertEvents.detail.dimensions') placeholderplaceholder</div>
              <div class="mt-1 text-sm text-gray-900 dark:text-white">
                <div v-if="getDimensionString(selected, 'platform')">platform={{ getDimensionString(selected, 'platform') placeholderplaceholder</div>
                <div v-if="selected.dimensions?.group_id">group_id={{ selected.dimensions.group_id placeholderplaceholder</div>
                <div v-if="getDimensionString(selected, 'region')">region={{ getDimensionString(selected, 'region') placeholderplaceholder</div>
              </div>
            </div>
          </div>


        <div class="rounded-xl border border-gray-200 bg-white p-4 dark:border-dark-700 dark:bg-dark-800">
          <div class="mb-3 flex flex-wrap items-center justify-between gap-3">
            <div>
              <div class="text-sm font-bold text-gray-900 dark:text-white">{{ t('admin.ops.alertEvents.detail.historyTitle') placeholderplaceholder</div>
              <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.ops.alertEvents.detail.historyHint') placeholderplaceholder</div>
            </div>
            <Select :model-value="historyRange" :options="historyRangeOptions" class="w-[140px]" @change="historyRange = String($event || '7d')" />
          </div>

          <div v-if="historyLoading" class="py-6 text-center text-xs text-gray-500 dark:text-gray-400">
            {{ t('admin.ops.alertEvents.detail.historyLoading') placeholderplaceholder
          </div>
          <div v-else-if="history.length === 0" class="py-6 text-center text-xs text-gray-500 dark:text-gray-400">
            {{ t('admin.ops.alertEvents.detail.historyEmpty') placeholderplaceholder
          </div>
          <div v-else class="overflow-hidden rounded-lg border border-gray-100 dark:border-dark-700">
            <table class="min-w-full divide-y divide-gray-100 dark:divide-dark-700">
              <thead class="bg-gray-50 dark:bg-dark-900">
                <tr>
                  <th class="px-3 py-2 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.ops.alertEvents.table.time') placeholderplaceholder</th>
                  <th class="px-3 py-2 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.ops.alertEvents.table.status') placeholderplaceholder</th>
                  <th class="px-3 py-2 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('admin.ops.alertEvents.table.metric') placeholderplaceholder</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-dark-700">
                <tr v-for="it in history" :key="it.id" class="hover:bg-gray-50 dark:hover:bg-dark-700/50">
                  <td class="px-3 py-2 text-xs text-gray-600 dark:text-gray-300">{{ formatDateTime(it.fired_at || it.created_at) placeholderplaceholder</td>
                  <td class="px-3 py-2 text-xs">
                    <span class="inline-flex items-center rounded-full px-2 py-1 text-[10px] font-bold ring-1 ring-inset" :class="statusBadgeClass(it.status)">
                      {{ formatStatusLabel(it.status) placeholderplaceholder
                    </span>
                  </td>
                  <td class="px-3 py-2 text-xs text-gray-600 dark:text-gray-300">
                    <span v-if="typeof it.metric_value === 'number' && typeof it.threshold_value === 'number'">
                      {{ it.metric_value.toFixed(2) placeholderplaceholder / {{ it.threshold_value.toFixed(2) placeholderplaceholder
                    </span>
                    <span v-else>-</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

