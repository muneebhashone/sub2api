<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch placeholder from 'vue'
import { useMediaQuery placeholder from '@vueuse/core'
import { useI18n placeholder from 'vue-i18n'
import { opsAPI, type OpsRuntimeLogConfig, type OpsSystemLog, type OpsSystemLogSinkHealth placeholder from '@/api/admin/ops'
import Pagination from '@/components/common/Pagination.vue'
import Select from '@/components/common/Select.vue'
import { useAppStore placeholder from '@/stores'

const appStore = useAppStore()
const { t placeholder = useI18n()

// 与 DataTable 一致：< 768px 切换为卡片视图，避免宽表在移动端被截断。
const isDesktopViewport = useMediaQuery('(min-width: 768px)')

const props = withDefaults(defineProps<{
  platformFilter?: string
  refreshToken?: number
placeholder>(), {
  platformFilter: '',
  refreshToken: 0
placeholder)

const loading = ref(false)
const logs = ref<OpsSystemLog[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const health = ref<OpsSystemLogSinkHealth>({
  queue_depth: 0,
  queue_capacity: 0,
  dropped_count: 0,
  write_failed_count: 0,
  written_count: 0,
  avg_write_delay_ms: 0
placeholder)

const runtimeLoading = ref(false)
const runtimeSaving = ref(false)
const runtimeConfig = reactive<OpsRuntimeLogConfig>({
  level: 'info',
  enable_sampling: false,
  sampling_initial: 100,
  sampling_thereafter: 100,
  caller: true,
  stacktrace_level: 'error',
  retention_days: 30
placeholder)

const filters = reactive({
  time_range: '1h' as '5m' | '30m' | '1h' | '6h' | '24h' | '7d' | '30d',
  start_time: '',
  end_time: '',
  host: '',
  level: '',
  component: '',
  request_id: '',
  client_request_id: '',
  user_id: '',
  api_key_id: '',
  account_id: '',
  platform: '',
  model: '',
  q: ''
placeholder)

const runtimeLevelOptions = [
  { value: 'debug', label: 'debug' placeholder,
  { value: 'info', label: 'info' placeholder,
  { value: 'warn', label: 'warn' placeholder,
  { value: 'error', label: 'error' placeholder
]

const stacktraceLevelOptions = [
  { value: 'none', label: 'none' placeholder,
  { value: 'error', label: 'error' placeholder,
  { value: 'fatal', label: 'fatal' placeholder
]

const timeRangeOptions = [
  { value: '5m', label: '5m' placeholder,
  { value: '30m', label: '30m' placeholder,
  { value: '1h', label: '1h' placeholder,
  { value: '6h', label: '6h' placeholder,
  { value: '24h', label: '24h' placeholder,
  { value: '7d', label: '7d' placeholder,
  { value: '30d', label: '30d' placeholder
]

const filterLevelOptions = computed(() => [
  { value: '', label: t('admin.ops.systemLogs.all') placeholder,
  { value: 'debug', label: 'debug' placeholder,
  { value: 'info', label: 'info' placeholder,
  { value: 'warn', label: 'warn' placeholder,
  { value: 'error', label: 'error' placeholder
])

const levelBadgeClass = (level: string) => {
  const v = String(level || '').toLowerCase()
  if (v === 'error' || v === 'fatal') return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
  if (v === 'warn' || v === 'warning') return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
  if (v === 'debug') return 'bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300'
  return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
placeholder

const formatTime = (value: string) => {
  if (!value) return '-'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return d.toLocaleString()
placeholder

const getExtraString = (extra: Record<string, any> | undefined, key: string) => {
  if (!extra) return ''
  const v = extra[key]
  if (v == null) return ''
  if (typeof v === 'string') return v.trim()
  if (typeof v === 'number' || typeof v === 'boolean') return String(v)
  return ''
placeholder

const formatSystemLogDetail = (row: OpsSystemLog) => {
  const parts: string[] = []
  const msg = String(row.message || '').trim()
  if (msg) parts.push(msg)

  const extra = row.extra || {placeholder
  const statusCode = getExtraString(extra, 'status_code')
  const latencyMs = getExtraString(extra, 'latency_ms')
  const method = getExtraString(extra, 'method')
  const path = getExtraString(extra, 'path')
  const clientIP = getExtraString(extra, 'client_ip')
  const protocol = getExtraString(extra, 'protocol')

  const accessParts: string[] = []
  if (statusCode) accessParts.push(`status=${statusCodeplaceholder`)
  if (latencyMs) accessParts.push(`latency_ms=${latencyMsplaceholder`)
  if (method) accessParts.push(`method=${methodplaceholder`)
  if (path) accessParts.push(`path=${pathplaceholder`)
  if (clientIP) accessParts.push(`ip=${clientIPplaceholder`)
  if (protocol) accessParts.push(`proto=${protocolplaceholder`)
  if (accessParts.length > 0) parts.push(accessParts.join(' '))

  const corrParts: string[] = []
  if (row.request_id) corrParts.push(`req=${row.request_idplaceholder`)
  if (row.client_request_id) corrParts.push(`client_req=${row.client_request_idplaceholder`)
  if (row.user_id != null) corrParts.push(`user=${row.user_idplaceholder`)
  if (row.api_key_id != null) corrParts.push(`key=${row.api_key_idplaceholder`)
  if (row.account_id != null) corrParts.push(`acc=${row.account_idplaceholder`)
  if (row.platform) corrParts.push(`platform=${row.platformplaceholder`)
  if (row.model) corrParts.push(`model=${row.modelplaceholder`)
  if (corrParts.length > 0) parts.push(corrParts.join(' '))

  const errors = getExtraString(extra, 'errors')
  if (errors) parts.push(`errors=${errorsplaceholder`)
  const err = getExtraString(extra, 'err') || getExtraString(extra, 'error')
  if (err) parts.push(`error=${errplaceholder`)

  // 用空格拼接，交给 CSS 自动换行，尽量“填满再换行”。
  return parts.join('  ')
placeholder

const toRFC3339 = (value: string) => {
  if (!value) return undefined
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return undefined
  return d.toISOString()
placeholder

const buildQuery = () => {
  const query: Record<string, any> = {
    page: page.value,
    page_size: pageSize.value,
    time_range: filters.time_range
  placeholder

  if (filters.time_range === '30d') {
    query.time_range = '30d'
  placeholder
  if (filters.start_time) query.start_time = toRFC3339(filters.start_time)
  if (filters.end_time) query.end_time = toRFC3339(filters.end_time)
  if (filters.host.trim()) query.host = filters.host.trim()
  if (filters.level.trim()) query.level = filters.level.trim()
  if (filters.component.trim()) query.component = filters.component.trim()
  if (filters.request_id.trim()) query.request_id = filters.request_id.trim()
  if (filters.client_request_id.trim()) query.client_request_id = filters.client_request_id.trim()
  if (filters.user_id.trim()) {
    const v = Number.parseInt(filters.user_id.trim(), 10)
    if (Number.isFinite(v) && v > 0) query.user_id = v
  placeholder
  if (filters.api_key_id.trim()) {
    const v = Number.parseInt(filters.api_key_id.trim(), 10)
    if (Number.isFinite(v) && v > 0) query.api_key_id = v
  placeholder
  if (filters.account_id.trim()) {
    const v = Number.parseInt(filters.account_id.trim(), 10)
    if (Number.isFinite(v) && v > 0) query.account_id = v
  placeholder
  if (filters.platform.trim()) query.platform = filters.platform.trim()
  if (filters.model.trim()) query.model = filters.model.trim()
  if (filters.q.trim()) query.q = filters.q.trim()
  return query
placeholder

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await opsAPI.listSystemLogs(buildQuery())
    logs.value = res.items || []
    total.value = res.total || 0
  placeholder catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to fetch logs', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.systemLogs.loadFailed'))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

const fetchHealth = async () => {
  try {
    health.value = await opsAPI.getSystemLogSinkHealth()
  placeholder catch {
    // 忽略健康数据读取失败，不影响主流程。
  placeholder
placeholder

const loadRuntimeConfig = async () => {
  runtimeLoading.value = true
  try {
    const cfg = await opsAPI.getRuntimeLogConfig()
    runtimeConfig.level = cfg.level
    runtimeConfig.enable_sampling = cfg.enable_sampling
    runtimeConfig.sampling_initial = cfg.sampling_initial
    runtimeConfig.sampling_thereafter = cfg.sampling_thereafter
    runtimeConfig.caller = cfg.caller
    runtimeConfig.stacktrace_level = cfg.stacktrace_level
    runtimeConfig.retention_days = cfg.retention_days
  placeholder catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to load runtime log config', err)
  placeholder finally {
    runtimeLoading.value = false
  placeholder
placeholder

const saveRuntimeConfig = async () => {
  runtimeSaving.value = true
  try {
    const saved = await opsAPI.updateRuntimeLogConfig({ ...runtimeConfig placeholder)
    runtimeConfig.level = saved.level
    runtimeConfig.enable_sampling = saved.enable_sampling
    runtimeConfig.sampling_initial = saved.sampling_initial
    runtimeConfig.sampling_thereafter = saved.sampling_thereafter
    runtimeConfig.caller = saved.caller
    runtimeConfig.stacktrace_level = saved.stacktrace_level
    runtimeConfig.retention_days = saved.retention_days
    appStore.showSuccess(t('admin.ops.systemLogs.runtimeConfigActive'))
  placeholder catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to save runtime log config', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.systemLogs.runtimeConfigSaveFailed'))
  placeholder finally {
    runtimeSaving.value = false
  placeholder
placeholder

const resetRuntimeConfig = async () => {
  const ok = window.confirm(t('admin.ops.systemLogs.resetRuntimeConfigConfirm'))
  if (!ok) return

  runtimeSaving.value = true
  try {
    const saved = await opsAPI.resetRuntimeLogConfig()
    runtimeConfig.level = saved.level
    runtimeConfig.enable_sampling = saved.enable_sampling
    runtimeConfig.sampling_initial = saved.sampling_initial
    runtimeConfig.sampling_thereafter = saved.sampling_thereafter
    runtimeConfig.caller = saved.caller
    runtimeConfig.stacktrace_level = saved.stacktrace_level
    runtimeConfig.retention_days = saved.retention_days
    appStore.showSuccess(t('admin.ops.systemLogs.runtimeConfigReset'))
    await fetchHealth()
  placeholder catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to reset runtime log config', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.systemLogs.runtimeConfigResetFailed'))
  placeholder finally {
    runtimeSaving.value = false
  placeholder
placeholder

const cleanupCurrentFilter = async () => {
  const ok = window.confirm(t('admin.ops.systemLogs.cleanupConfirm'))
  if (!ok) return
  try {
    const payload = {
      start_time: toRFC3339(filters.start_time),
      end_time: toRFC3339(filters.end_time),
      host: filters.host.trim() || undefined,
      level: filters.level.trim() || undefined,
      component: filters.component.trim() || undefined,
      request_id: filters.request_id.trim() || undefined,
      client_request_id: filters.client_request_id.trim() || undefined,
      user_id: filters.user_id.trim() ? Number.parseInt(filters.user_id.trim(), 10) : undefined,
      api_key_id: filters.api_key_id.trim() ? Number.parseInt(filters.api_key_id.trim(), 10) : undefined,
      account_id: filters.account_id.trim() ? Number.parseInt(filters.account_id.trim(), 10) : undefined,
      platform: filters.platform.trim() || undefined,
      model: filters.model.trim() || undefined,
      q: filters.q.trim() || undefined
    placeholder
    const res = await opsAPI.cleanupSystemLogs(payload)
    appStore.showSuccess(t('admin.ops.systemLogs.cleanupSuccess', { count: res.deleted || 0 placeholder))
    page.value = 1
    await Promise.all([fetchLogs(), fetchHealth()])
  placeholder catch (err: any) {
    console.error('[OpsSystemLogTable] Failed to cleanup logs', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.systemLogs.cleanupFailed'))
  placeholder
placeholder

const resetFilters = () => {
  filters.time_range = '1h'
  filters.start_time = ''
  filters.end_time = ''
  filters.host = ''
  filters.level = ''
  filters.component = ''
  filters.request_id = ''
  filters.client_request_id = ''
  filters.user_id = ''
  filters.api_key_id = ''
  filters.account_id = ''
  filters.platform = props.platformFilter || ''
  filters.model = ''
  filters.q = ''
  page.value = 1
  fetchLogs()
placeholder

watch(() => props.platformFilter, (v) => {
  if (v && !filters.platform) {
    filters.platform = v
    page.value = 1
    fetchLogs()
  placeholder
placeholder)

watch(() => props.refreshToken, () => {
  fetchLogs()
  fetchHealth()
placeholder)

const onPageChange = (next: number) => {
  page.value = next
  fetchLogs()
placeholder

const onPageSizeChange = (next: number) => {
  pageSize.value = next
  page.value = 1
  fetchLogs()
placeholder

const applyFilters = () => {
  page.value = 1
  fetchLogs()
placeholder

const hasData = computed(() => logs.value.length > 0)

onMounted(async () => {
  if (props.platformFilter) {
    filters.platform = props.platformFilter
  placeholder
  await Promise.all([fetchLogs(), fetchHealth(), loadRuntimeConfig()])
placeholder)
</script>

<template>
  <section class="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-900/60">
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <div>
        <h3 class="text-sm font-bold text-gray-900 dark:text-white">{{ t('admin.ops.systemLogs.title') placeholderplaceholder</h3>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.ops.systemLogs.description') placeholderplaceholder</p>
      </div>
      <div class="flex flex-wrap items-center gap-2 text-xs">
        <span class="rounded-md bg-gray-100 px-2 py-1 text-gray-700 dark:bg-dark-700 dark:text-gray-200">{{ t('admin.ops.systemLogs.queue') placeholderplaceholder {{ health.queue_depth placeholderplaceholder/{{ health.queue_capacity placeholderplaceholder</span>
        <span class="rounded-md bg-gray-100 px-2 py-1 text-gray-700 dark:bg-dark-700 dark:text-gray-200">{{ t('admin.ops.systemLogs.written') placeholderplaceholder {{ health.written_count placeholderplaceholder</span>
        <span class="rounded-md bg-amber-100 px-2 py-1 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300">{{ t('admin.ops.systemLogs.dropped') placeholderplaceholder {{ health.dropped_count placeholderplaceholder</span>
        <span class="rounded-md bg-red-100 px-2 py-1 text-red-700 dark:bg-red-900/30 dark:text-red-300">{{ t('admin.ops.systemLogs.failed') placeholderplaceholder {{ health.write_failed_count placeholderplaceholder</span>
      </div>
    </div>

    <div class="mb-4 rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-800/70">
      <div class="mb-2 flex items-center justify-between">
        <div class="text-xs font-semibold text-gray-700 dark:text-gray-200">{{ t('admin.ops.systemLogs.runtimeConfig') placeholderplaceholder</div>
        <span v-if="runtimeLoading" class="text-xs text-gray-500">{{ t('common.loading') placeholderplaceholder</span>
      </div>
      <div class="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-6">
        <label class="text-xs text-gray-600 dark:text-gray-300">
          {{ t('admin.ops.systemLogs.level') placeholderplaceholder
          <Select v-model="runtimeConfig.level" class="mt-1" :options="runtimeLevelOptions" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          {{ t('admin.ops.systemLogs.stacktraceThreshold') placeholderplaceholder
          <Select v-model="runtimeConfig.stacktrace_level" class="mt-1" :options="stacktraceLevelOptions" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          {{ t('admin.ops.systemLogs.samplingInitial') placeholderplaceholder
          <input v-model.number="runtimeConfig.sampling_initial" type="number" min="1" class="input mt-1" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          {{ t('admin.ops.systemLogs.samplingThereafter') placeholderplaceholder
          <input v-model.number="runtimeConfig.sampling_thereafter" type="number" min="1" class="input mt-1" />
        </label>
        <label class="text-xs text-gray-600 dark:text-gray-300">
          {{ t('admin.ops.systemLogs.retentionDays') placeholderplaceholder
          <input v-model.number="runtimeConfig.retention_days" type="number" min="1" max="3650" class="input mt-1" />
        </label>
        <div class="md:col-span-2 xl:col-span-6">
          <div class="grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto] lg:items-end">
            <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
              <label class="inline-flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                <input v-model="runtimeConfig.caller" type="checkbox" />
                {{ t('admin.ops.systemLogs.caller') placeholderplaceholder
              </label>
              <label class="inline-flex items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                <input v-model="runtimeConfig.enable_sampling" type="checkbox" />
                {{ t('admin.ops.systemLogs.sampling') placeholderplaceholder
              </label>
            </div>
            <div class="flex flex-wrap items-center gap-2 lg:justify-end">
              <button type="button" class="btn btn-primary btn-sm" :disabled="runtimeSaving" @click="saveRuntimeConfig">
                {{ runtimeSaving ? t('common.saving') : t('admin.ops.systemLogs.saveAndApply') placeholderplaceholder
              </button>
              <button type="button" class="btn btn-secondary btn-sm" :disabled="runtimeSaving" @click="resetRuntimeConfig">
                {{ t('admin.ops.systemLogs.resetDefaults') placeholderplaceholder
              </button>
            </div>
          </div>
        </div>
      </div>
      <p v-if="health.last_error" class="mt-2 text-xs text-red-600 dark:text-red-400">{{ t('admin.ops.systemLogs.latestWriteError') placeholderplaceholder {{ health.last_error placeholderplaceholder</p>
    </div>

    <div class="mb-4 grid grid-cols-1 gap-3 md:grid-cols-5">
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.timeRange') placeholderplaceholder
        <Select v-model="filters.time_range" class="mt-1" :options="timeRangeOptions" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.startTime') placeholderplaceholder
        <input v-model="filters.start_time" type="datetime-local" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.endTime') placeholderplaceholder
        <input v-model="filters.end_time" type="datetime-local" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.level') placeholderplaceholder
        <Select v-model="filters.level" class="mt-1" :options="filterLevelOptions" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.component') placeholderplaceholder
        <input v-model="filters.component" type="text" class="input mt-1" :placeholder="t('admin.ops.systemLogs.componentPlaceholder')" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.host') placeholderplaceholder
        <input v-model="filters.host" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        request_id
        <input v-model="filters.request_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        client_request_id
        <input v-model="filters.client_request_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        user_id
        <input v-model="filters.user_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.keyId') placeholderplaceholder
        <input v-model="filters.api_key_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        account_id
        <input v-model="filters.account_id" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.platform') placeholderplaceholder
        <input v-model="filters.platform" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.model') placeholderplaceholder
        <input v-model="filters.model" type="text" class="input mt-1" />
      </label>
      <label class="text-xs text-gray-600 dark:text-gray-300">
        {{ t('admin.ops.systemLogs.keyword') placeholderplaceholder
        <input v-model="filters.q" type="text" class="input mt-1" :placeholder="t('admin.ops.systemLogs.keywordPlaceholder')" />
      </label>
    </div>

    <div class="mb-3 flex flex-wrap gap-2">
      <button type="button" class="btn btn-primary btn-sm" @click="applyFilters">{{ t('admin.ops.systemLogs.search') placeholderplaceholder</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="resetFilters">{{ t('common.reset') placeholderplaceholder</button>
      <button type="button" class="btn btn-danger btn-sm" @click="cleanupCurrentFilter">{{ t('admin.ops.systemLogs.cleanCurrentFilters') placeholderplaceholder</button>
      <button type="button" class="btn btn-secondary btn-sm" @click="fetchHealth">{{ t('admin.ops.systemLogs.refreshHealth') placeholderplaceholder</button>
    </div>

    <div class="overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700">
      <div v-if="loading" class="px-4 py-8 text-center text-sm text-gray-500">{{ t('common.loading') placeholderplaceholder</div>
      <div v-else-if="!hasData" class="px-4 py-8 text-center text-sm text-gray-500">{{ t('admin.ops.systemLogs.empty') placeholderplaceholder</div>
      <div v-else-if="!isDesktopViewport" class="divide-y divide-gray-100 dark:divide-dark-800">
        <div v-for="row in logs" :key="row.id" class="space-y-1.5 p-3">
          <div class="flex items-center justify-between gap-2">
            <span class="inline-flex rounded-full px-2 py-0.5 text-xs font-semibold" :class="levelBadgeClass(row.level)">
              {{ row.level placeholderplaceholder
            </span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatTime(row.created_at) placeholderplaceholder</span>
          </div>
          <div v-if="row.host" class="truncate text-xs text-gray-500 dark:text-gray-400" :title="row.host">
            {{ row.host placeholderplaceholder
          </div>
          <div class="whitespace-normal break-all text-xs text-gray-700 dark:text-gray-300">
            {{ formatSystemLogDetail(row) placeholderplaceholder
          </div>
        </div>
      </div>
      <div v-else class="overflow-auto">
        <table class="min-w-full table-fixed divide-y divide-gray-200 dark:divide-dark-700">
          <thead class="bg-gray-50 dark:bg-dark-900">
            <tr>
              <th class="w-[170px] px-3 py-2 text-left text-[11px] font-semibold text-gray-500">{{ t('admin.ops.systemLogs.time') placeholderplaceholder</th>
              <th class="w-[160px] px-3 py-2 text-left text-[11px] font-semibold text-gray-500">{{ t('admin.ops.systemLogs.host') placeholderplaceholder</th>
              <th class="w-[80px] px-3 py-2 text-left text-[11px] font-semibold text-gray-500">{{ t('admin.ops.systemLogs.level') placeholderplaceholder</th>
              <th class="px-3 py-2 text-left text-[11px] font-semibold text-gray-500">{{ t('admin.ops.systemLogs.logDetails') placeholderplaceholder</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-dark-800">
            <tr v-for="row in logs" :key="row.id" class="align-top">
              <td class="px-3 py-2 text-xs text-gray-700 dark:text-gray-300">{{ formatTime(row.created_at) placeholderplaceholder</td>
              <td class="px-3 py-2 text-xs text-gray-700 dark:text-gray-300">
                <span class="block truncate" :title="row.host || '-'">{{ row.host || '-' placeholderplaceholder</span>
              </td>
              <td class="px-3 py-2 text-xs">
                <span class="inline-flex rounded-full px-2 py-0.5 font-semibold" :class="levelBadgeClass(row.level)">
                  {{ row.level placeholderplaceholder
                </span>
              </td>
              <td class="px-3 py-2 text-xs text-gray-700 dark:text-gray-300 whitespace-normal break-all">
                {{ formatSystemLogDetail(row) placeholderplaceholder
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination
        :total="total"
        :page="page"
        :page-size="pageSize"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </div>
  </section>
</template>
