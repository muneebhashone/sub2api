<template>
  <AppLayout>
    <div class="space-y-6">
      <UsageStatsCards :stats="usageStats" />
      <!-- Charts Section -->
      <div class="space-y-4">
        <div class="card p-4">
          <div class="flex flex-wrap items-center gap-4">
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.dashboard.timeRange') placeholderplaceholder:</span>
              <DateRangePicker
                v-model:start-date="startDate"
                v-model:end-date="endDate"
                @change="onDateRangeChange"
              />
            </div>
            <div class="ml-auto flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('admin.dashboard.granularity') placeholderplaceholder:</span>
              <div class="w-28">
                <Select v-model="granularity" :options="granularityOptions" @change="loadChartData" />
              </div>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <ModelDistributionChart
            v-model:source="modelDistributionSource"
            v-model:metric="modelDistributionMetric"
            :model-stats="requestedModelStats"
            :upstream-model-stats="upstreamModelStats"
            :mapping-model-stats="mappingModelStats"
            :loading="modelStatsLoading"
            :show-source-toggle="true"
            :show-metric-toggle="true"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
          />
          <GroupDistributionChart
            v-model:metric="groupDistributionMetric"
            :group-stats="groupStats"
            :loading="chartsLoading"
            :show-metric-toggle="true"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
          />
        </div>
        <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
          <EndpointDistributionChart
            v-model:source="endpointDistributionSource"
            v-model:metric="endpointDistributionMetric"
            :endpoint-stats="inboundEndpointStats"
            :upstream-endpoint-stats="upstreamEndpointStats"
            :endpoint-path-stats="endpointPathStats"
            :loading="endpointStatsLoading"
            :show-source-toggle="true"
            :show-metric-toggle="true"
            :title="t('usage.endpointDistribution')"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
          />
          <TokenUsageTrend :trend-data="trendData" :loading="chartsLoading" />
        </div>
      </div>
      <!-- 明细区：tab 栏 + 筛选 + 内容收进同一张卡片，消除割裂感 -->
      <div class="card">
        <div class="flex flex-wrap items-center border-b border-gray-200 px-2 dark:border-dark-700 sm:px-4">
          <button
            v-for="tab in detailTabs"
            :key="tab.key"
            type="button"
            data-testid="usage-detail-tab"
            class="-mb-px inline-flex items-center gap-1.5 border-b-2 px-3 py-3 text-sm font-medium transition-colors sm:px-4"
            :class="activeTab === tab.key
              ? 'border-primary-500 text-primary-600 dark:text-primary-400'
              : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-gray-400 dark:hover:border-dark-500 dark:hover:text-gray-200'"
            @click="switchTab(tab.key)"
          >
            <Icon :name="tab.icon" size="sm" />
            {{ tab.label placeholderplaceholder
          </button>
        </div>

        <UsageFilters v-model="filters" ref="usageFiltersRef" flat :mode="activeTab" class="border-b border-gray-100 dark:border-dark-700/50" :start-date="startDate" :end-date="endDate" :exporting="exporting" :model-options="modelNameOptions" @change="applyFilters" @refresh="refreshData" @reset="resetFilters" @cleanup="openCleanupDialog" @export="exportToExcel">
          <template #after-reset>
            <div v-if="activeTab !== 'ranking'" class="relative" ref="columnDropdownRef">
              <button
                @click="showColumnDropdown = !showColumnDropdown"
                class="btn btn-secondary px-2 md:px-3"
                :title="t('admin.users.columnSettings')"
              >
                <svg class="h-4 w-4 md:mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 4.5v15m6-15v15m-10.875 0h15.75c.621 0 1.125-.504 1.125-1.125V5.625c0-.621-.504-1.125-1.125-1.125H4.125C3.504 4.5 3 5.004 3 5.625v12.75c0 .621.504 1.125 1.125 1.125z" />
                </svg>
                <span class="hidden md:inline">{{ t('admin.users.columnSettings') placeholderplaceholder</span>
              </button>
              <div
                v-if="showColumnDropdown"
                class="absolute right-0 top-full z-50 mt-1 max-h-80 w-48 overflow-y-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800"
              >
                <button
                  v-for="col in currentToggleableColumns"
                  :key="col.key"
                  @click="toggleCurrentColumn(col.key)"
                  class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
                >
                  <span>{{ col.label placeholderplaceholder</span>
                  <Icon
                    v-if="isCurrentColumnVisible(col.key)"
                    name="check"
                    size="sm"
                    class="text-primary-500"
                    :stroke-width="2"
                  />
                </button>
              </div>
            </div>
          </template>
        </UsageFilters>

        <div v-show="activeTab === 'usage'" class="overflow-hidden rounded-b-2xl">
          <UsageTable
            flat
            :data="usageLogs"
            :loading="loading"
            :columns="visibleColumns"
            :server-side-sort="true"
            :default-sort-key="'created_at'"
            :default-sort-order="'desc'"
            @sort="handleSort"
            @userClick="handleUserClick"
            @ipGeoBatchFailed="handleIpGeoBatchFailed"
          />
          <Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="handlePageChange" @update:pageSize="handlePageSizeChange" />
        </div>
        <div v-show="activeTab === 'errors'" class="overflow-hidden rounded-b-2xl">
          <OpsErrorLogTable
            flat
            :rows="errRows" :total="errTotal" :loading="errLoading"
            :page="errPage" :page-size="errPageSize"
            :visible-column-keys="errVisibleColumnKeys"
            user-clickable
            @userClick="handleUserClick"
            @openErrorDetail="openError"
            @sort="onErrSort"
            @update:page="onErrPage"
            @update:pageSize="onErrPageSize"
            @ipGeoBatchFailed="handleIpGeoBatchFailed" />
        </div>
        <!-- 懒挂载：首次切到该 tab 才请求排行数据，之后随筛选自动刷新 -->
        <div v-if="rankingMounted" v-show="activeTab === 'ranking'" class="overflow-hidden rounded-b-2xl">
          <UserTokenRanking
            ref="rankingRef"
            :start-date="startDate"
            :end-date="endDate"
            :filters="breakdownFilters"
            :model="filters.model"
            @select-user="handleRankingSelectUser"
          />
        </div>
      </div>
      <OpsErrorDetailModal v-model:show="showErrorModal" :error-id="selectedErrorId" :error-type="'request'" />
    </div>
  </AppLayout>
  <UsageExportProgress :show="exportProgress.show" :progress="exportProgress.progress" :current="exportProgress.current" :total="exportProgress.total" :estimated-time="exportProgress.estimatedTime" @cancel="cancelExport" />
  <UsageCleanupDialog
    :show="cleanupDialogVisible"
    :filters="filters"
    :start-date="startDate"
    :end-date="endDate"
    @close="cleanupDialogVisible = false"
  />
  <!-- Balance history modal triggered from usage table user click -->
  <UserBalanceHistoryModal
    :show="showBalanceHistoryModal"
    :user="balanceHistoryUser"
    :hide-actions="true"
    @close="showBalanceHistoryModal = false; balanceHistoryUser = null"
  />
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { saveAs placeholder from 'file-saver'
import { useRoute placeholder from 'vue-router'
import { useAppStore placeholder from '@/stores/app'; import { adminAPI placeholder from '@/api/admin'; import { adminUsageAPI placeholder from '@/api/admin/usage'
import { getPersistedPageSize placeholder from '@/composables/usePersistedPageSize'
import { formatReasoningEffort placeholder from '@/utils/format'
import { resolveUsageRequestType, requestTypeToLegacyStream placeholder from '@/utils/usageRequestType'
import AppLayout from '@/components/layout/AppLayout.vue'; import Pagination from '@/components/common/Pagination.vue'; import Select from '@/components/common/Select.vue'; import DateRangePicker from '@/components/common/DateRangePicker.vue'
import UsageStatsCards from '@/components/admin/usage/UsageStatsCards.vue'; import UsageFilters from '@/components/admin/usage/UsageFilters.vue'
import UsageTable from '@/components/admin/usage/UsageTable.vue'; import UsageExportProgress from '@/components/admin/usage/UsageExportProgress.vue'
import UserTokenRanking from '@/components/admin/usage/UserTokenRanking.vue'
import UsageCleanupDialog from '@/components/admin/usage/UsageCleanupDialog.vue'
import UserBalanceHistoryModal from '@/components/admin/user/UserBalanceHistoryModal.vue'
import OpsErrorLogTable from '@/views/admin/ops/components/OpsErrorLogTable.vue'
import OpsErrorDetailModal from '@/views/admin/ops/components/OpsErrorDetailModal.vue'
import { listErrorLogs placeholder from '@/api/admin/ops'
import type { OpsErrorLog placeholder from '@/api/admin/ops'
import ModelDistributionChart from '@/components/charts/ModelDistributionChart.vue'; import GroupDistributionChart from '@/components/charts/GroupDistributionChart.vue'; import TokenUsageTrend from '@/components/charts/TokenUsageTrend.vue'
import EndpointDistributionChart from '@/components/charts/EndpointDistributionChart.vue'
import Icon from '@/components/icons/Icon.vue'
import type { AdminUsageLog, TrendDataPoint, ModelStat, GroupStat, EndpointStat, AdminUser placeholder from '@/types'; import type { AdminUsageStatsResponse, AdminUsageQueryParams placeholder from '@/api/admin/usage'

const { t placeholder = useI18n()
const appStore = useAppStore()
type DistributionMetric = 'tokens' | 'actual_cost'
type EndpointSource = 'inbound' | 'upstream' | 'path'
type ModelDistributionSource = 'requested' | 'upstream' | 'mapping'
const route = useRoute()
const usageStats = ref<AdminUsageStatsResponse | null>(null); const usageLogs = ref<AdminUsageLog[]>([]); const loading = ref(false); const exporting = ref(false)
const trendData = ref<TrendDataPoint[]>([]); const requestedModelStats = ref<ModelStat[]>([]); const upstreamModelStats = ref<ModelStat[]>([]); const mappingModelStats = ref<ModelStat[]>([]); const groupStats = ref<GroupStat[]>([]); const chartsLoading = ref(false); const modelStatsLoading = ref(false); const granularity = ref<'day' | 'hour'>('hour')
const modelDistributionMetric = ref<DistributionMetric>('tokens')
const modelDistributionSource = ref<ModelDistributionSource>('requested')
const loadedModelSources = reactive<Record<ModelDistributionSource, boolean>>({
  requested: false,
  upstream: false,
  mapping: false,
placeholder)
const groupDistributionMetric = ref<DistributionMetric>('tokens')
const endpointDistributionMetric = ref<DistributionMetric>('tokens')
const endpointDistributionSource = ref<EndpointSource>('inbound')
const inboundEndpointStats = ref<EndpointStat[]>([])
const upstreamEndpointStats = ref<EndpointStat[]>([])
const endpointPathStats = ref<EndpointStat[]>([])
const endpointStatsLoading = ref(false)
let abortController: AbortController | null = null; let exportAbortController: AbortController | null = null
let chartReqSeq = 0
let statsReqSeq = 0
let modelStatsReqSeq = 0
const exportProgress = reactive({ show: false, progress: 0, current: 0, total: 0, estimatedTime: '' placeholder)
const cleanupDialogVisible = ref(false)
// Balance history modal state
const showBalanceHistoryModal = ref(false)
const balanceHistoryUser = ref<AdminUser | null>(null)

const breakdownFilters = computed(() => {
  const f: Record<string, any> = {placeholder
  if (filters.value.user_id) f.user_id = filters.value.user_id
  if (filters.value.api_key_id) f.api_key_id = filters.value.api_key_id
  if (filters.value.account_id) f.account_id = filters.value.account_id
  if (filters.value.group_id) f.group_id = filters.value.group_id
  if (filters.value.request_type != null) f.request_type = filters.value.request_type
  if (filters.value.billing_type != null) f.billing_type = filters.value.billing_type
  return f
placeholder)

const modelNameOptions = computed(() =>
  Array.from(new Set(requestedModelStats.value.map((m) => m.model).filter(Boolean))).sort()
)

const handleUserClick = async (userId: number) => {
  try {
    const user = await adminAPI.users.getById(userId, true)
    balanceHistoryUser.value = user
    showBalanceHistoryModal.value = true
  placeholder catch {
    appStore.showError(t('admin.usage.failedToLoadUser'))
  placeholder
placeholder

// Drill down from the per-user token ranking: scope the whole usage view to
// that user and jump to the usage-detail tab so the drill-down is visible.
const handleRankingSelectUser = (userId: number, email: string) => {
  filters.value = { ...filters.value, user_id: userId placeholder
  usageFiltersRef.value?.setUserKeyword?.(email || '')
  activeTab.value = 'usage'
  applyFilters()
placeholder

const granularityOptions = computed(() => [{ value: 'day', label: t('admin.dashboard.day') placeholder, { value: 'hour', label: t('admin.dashboard.hour') placeholder])
// Use local timezone to avoid UTC timezone issues
const formatLD = (d: Date) => {
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${yearplaceholder-${monthplaceholder-${dayplaceholder`
placeholder
const getLast24HoursRangeDates = (): { start: string; end: string placeholder => {
  const end = new Date()
  const start = new Date(end.getTime() - 24 * 60 * 60 * 1000)
  return {
    start: formatLD(start),
    end: formatLD(end)
  placeholder
placeholder
const getGranularityForRange = (start: string, end: string): 'day' | 'hour' => {
  const startTime = new Date(`${startplaceholderT00:00:00`).getTime()
  const endTime = new Date(`${endplaceholderT00:00:00`).getTime()
  const daysDiff = Math.ceil((endTime - startTime) / (1000 * 60 * 60 * 24))
  return daysDiff <= 1 ? 'hour' : 'day'
placeholder
const defaultRange = getLast24HoursRangeDates()
const startDate = ref(defaultRange.start); const endDate = ref(defaultRange.end)
const filters = ref<AdminUsageQueryParams>({ user_id: undefined, model: undefined, group_id: undefined, request_type: undefined, billing_type: null, start_date: startDate.value, end_date: endDate.value placeholder)
const pagination = reactive({ page: 1, page_size: getPersistedPageSize(), total: 0 placeholder)
const sortState = reactive({
  sort_by: 'created_at',
  sort_order: 'desc' as 'asc' | 'desc'
placeholder)

const getSingleQueryValue = (value: string | null | Array<string | null> | undefined): string | undefined => {
  if (Array.isArray(value)) return value.find((item): item is string => typeof item === 'string' && item.length > 0)
  return typeof value === 'string' && value.length > 0 ? value : undefined
placeholder

const getNumericQueryValue = (value: string | null | Array<string | null> | undefined): number | undefined => {
  const raw = getSingleQueryValue(value)
  if (!raw) return undefined
  const parsed = Number(raw)
  return Number.isFinite(parsed) ? parsed : undefined
placeholder

const applyRouteQueryFilters = () => {
  const queryStartDate = getSingleQueryValue(route.query.start_date)
  const queryEndDate = getSingleQueryValue(route.query.end_date)
  const queryUserId = getNumericQueryValue(route.query.user_id)

  if (queryStartDate) {
    startDate.value = queryStartDate
  placeholder
  if (queryEndDate) {
    endDate.value = queryEndDate
  placeholder

  filters.value = {
    ...filters.value,
    user_id: queryUserId,
    start_date: startDate.value,
    end_date: endDate.value
  placeholder
  granularity.value = getGranularityForRange(startDate.value, endDate.value)
placeholder

const loadRouteUserFilterLabel = async () => {
  const requestedUserId = filters.value.user_id
  if (!requestedUserId) return
  const userSearchRevision = usageFiltersRef.value?.getUserSearchRevision?.()

  const routeUserFilterIsCurrent = () => (
    filters.value.user_id === requestedUserId
    && usageFiltersRef.value?.getUserSearchRevision?.() === userSearchRevision
  )

  try {
    const user = await adminAPI.users.getById(requestedUserId, true)
    if (!routeUserFilterIsCurrent()) return
    usageFiltersRef.value?.setUserKeyword?.(user.email || String(requestedUserId))
  placeholder catch {
    if (!routeUserFilterIsCurrent()) return
    usageFiltersRef.value?.setUserKeyword?.(String(requestedUserId))
  placeholder
placeholder

const onDateRangeChange = (range: { startDate: string; endDate: string; preset: string | null placeholder) => {
  startDate.value = range.startDate
  endDate.value = range.endDate
  filters.value = {
    ...filters.value,
    start_date: range.startDate,
    end_date: range.endDate
  placeholder
  granularity.value = getGranularityForRange(range.startDate, range.endDate)
  applyFilters()
placeholder

const buildUsageListParams = (
  page: number,
  pageSize: number,
  exactTotal: boolean
): AdminUsageQueryParams => {
  const requestType = filters.value.request_type
  const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
  return {
    page,
    page_size: pageSize,
    exact_total: exactTotal,
    ...filters.value,
    stream: legacyStream === null ? undefined : legacyStream,
    sort_by: sortState.sort_by,
    sort_order: sortState.sort_order
  placeholder
placeholder

const loadLogs = async () => {
  abortController?.abort(); const c = new AbortController(); abortController = c; loading.value = true
  try {
    const res = await adminAPI.usage.list(
      buildUsageListParams(pagination.page, pagination.page_size, false),
      { signal: c.signal placeholder
    )
    if(!c.signal.aborted) { usageLogs.value = res.items; pagination.total = res.total placeholder
  placeholder catch (error: any) { if(error?.name !== 'AbortError') console.error('Failed to load usage logs:', error) placeholder finally { if(abortController === c) loading.value = false placeholder
placeholder
const loadStats = async (force = false) => {
  const seq = ++statsReqSeq
  endpointStatsLoading.value = true
  try {
    const requestType = filters.value.request_type
    const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
    const s = await adminAPI.usage.getStats({
      ...filters.value,
      stream: legacyStream === null ? undefined : legacyStream,
      ...(force ? { nocache: 1 placeholder : {placeholder),
    placeholder)
    if (seq !== statsReqSeq) return
    usageStats.value = s
    inboundEndpointStats.value = s.endpoints || []
    upstreamEndpointStats.value = s.upstream_endpoints || []
    endpointPathStats.value = s.endpoint_paths || []
  placeholder catch (error) {
    if (seq !== statsReqSeq) return
    console.error('Failed to load usage stats:', error)
    inboundEndpointStats.value = []
    upstreamEndpointStats.value = []
    endpointPathStats.value = []
  placeholder finally {
    if (seq === statsReqSeq) endpointStatsLoading.value = false
  placeholder
placeholder

// 失效模型统计缓存:仅标记需要重取,保留旧数据直到新数据到达(避免刷新时图表闪空)。
const invalidateModelStatsCache = () => {
  loadedModelSources.requested = false
  loadedModelSources.upstream = false
  loadedModelSources.mapping = false
placeholder

const loadModelStats = async (source: ModelDistributionSource, force = false) => {
  if (!force && loadedModelSources[source]) {
    return
  placeholder

  const seq = ++modelStatsReqSeq
  modelStatsLoading.value = true
  try {
    const requestType = filters.value.request_type
    const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
    const baseParams = {
      start_date: filters.value.start_date || startDate.value,
      end_date: filters.value.end_date || endDate.value,
      user_id: filters.value.user_id,
      model: filters.value.model,
      api_key_id: filters.value.api_key_id,
      account_id: filters.value.account_id,
      group_id: filters.value.group_id,
      request_type: requestType,
      stream: legacyStream === null ? undefined : legacyStream,
      billing_type: filters.value.billing_type,
	  upstream_model_mismatch: filters.value.upstream_model_mismatch,
    placeholder

    const response = await adminAPI.dashboard.getModelStats({ ...baseParams, model_source: source placeholder)

    if (seq !== modelStatsReqSeq) return

    const models = response.models || []
    if (source === 'requested') {
      requestedModelStats.value = models
    placeholder else if (source === 'upstream') {
      upstreamModelStats.value = models
    placeholder else {
      mappingModelStats.value = models
    placeholder
    loadedModelSources[source] = true
  placeholder catch (error) {
    if (seq !== modelStatsReqSeq) return
    console.error('Failed to load model stats:', error)
    if (source === 'requested') {
      requestedModelStats.value = []
    placeholder else if (source === 'upstream') {
      upstreamModelStats.value = []
    placeholder else {
      mappingModelStats.value = []
    placeholder
    loadedModelSources[source] = false
  placeholder finally {
    if (seq === modelStatsReqSeq) modelStatsLoading.value = false
  placeholder
placeholder

const loadChartData = async () => {
  const seq = ++chartReqSeq
  chartsLoading.value = true
  try {
    const requestType = filters.value.request_type
    const legacyStream = requestType ? requestTypeToLegacyStream(requestType) : filters.value.stream
    const snapshot = await adminAPI.dashboard.getSnapshotV2({
      start_date: filters.value.start_date || startDate.value,
      end_date: filters.value.end_date || endDate.value,
      granularity: granularity.value,
      user_id: filters.value.user_id,
      model: filters.value.model,
      api_key_id: filters.value.api_key_id,
      account_id: filters.value.account_id,
      group_id: filters.value.group_id,
      request_type: requestType,
      stream: legacyStream === null ? undefined : legacyStream,
      billing_type: filters.value.billing_type,
	  upstream_model_mismatch: filters.value.upstream_model_mismatch,
      include_stats: false,
      include_trend: true,
      include_model_stats: false,
      include_group_stats: true,
      include_users_trend: false
    placeholder)
    if (seq !== chartReqSeq) return
    trendData.value = snapshot.trend || []
    groupStats.value = snapshot.groups || []
  placeholder catch (error) { console.error('Failed to load chart data:', error) placeholder finally { if (seq === chartReqSeq) chartsLoading.value = false placeholder
placeholder
const applyFilters = () => {
  pagination.page = 1
  invalidateModelStatsCache()
  loadLogs()
  loadStats()
  loadModelStats(modelDistributionSource.value, true)
  loadChartData()
  errPage.value = 1
  if (activeTab.value === 'errors') {
    loadAdminErrors()
  placeholder else {
    errRows.value = []
  placeholder
placeholder
const refreshData = () => {
  invalidateModelStatsCache()
  loadLogs()
  loadStats(true)
  loadModelStats(modelDistributionSource.value, true)
  loadChartData()
  if (activeTab.value === 'errors') loadAdminErrors()
  if (rankingMounted.value) rankingRef.value?.reload()
placeholder
const resetFilters = () => {
  const range = getLast24HoursRangeDates()
  startDate.value = range.start
  endDate.value = range.end
  filters.value = { start_date: startDate.value, end_date: endDate.value, request_type: undefined, billing_type: null, billing_mode: undefined placeholder
  granularity.value = getGranularityForRange(startDate.value, endDate.value)
  applyFilters()
placeholder
const handlePageChange = (p: number) => { pagination.page = p; loadLogs() placeholder
const handlePageSizeChange = (s: number) => { pagination.page_size = s; pagination.page = 1; loadLogs() placeholder
const handleSort = (key: string, order: 'asc' | 'desc') => {
  sortState.sort_by = key
  sortState.sort_order = order
  pagination.page = 1
  loadLogs()
placeholder

const handleIpGeoBatchFailed = () => {
  appStore.showError(t('usage.ipGeo.batchFailed'))
placeholder
const cancelExport = () => exportAbortController?.abort()
const openCleanupDialog = () => { cleanupDialogVisible.value = true placeholder
const getRequestTypeLabel = (log: AdminUsageLog): string => {
  const requestType = resolveUsageRequestType(log)
  if (requestType === 'cyber') return t('usage.cyber')
  if (requestType === 'live') return t('usage.live')
  if (requestType === 'ws_v2') return t('usage.ws')
  if (requestType === 'stream') return t('usage.stream')
  if (requestType === 'sync') return t('usage.sync')
  return t('usage.unknown')
placeholder

const exportToExcel = async () => {
  if (exporting.value) return; exporting.value = true; exportProgress.show = true
  const c = new AbortController(); exportAbortController = c
  try {
    let p = 1; let total = pagination.total; let exportedCount = 0
    const XLSX = await import('xlsx')
    const headers = [
      t('usage.time'), t('admin.usage.user'), t('usage.apiKeyFilter'),
      t('admin.usage.account'), t('usage.requestedModel'), t('usage.sentUpstreamModel'), t('usage.upstreamResponseModel'), t('usage.upstreamModelMismatch'), t('usage.reasoningEffort'), t('admin.usage.group'),
      t('usage.inboundEndpoint'), t('usage.upstreamEndpoint'),
      t('usage.type'),
      t('admin.usage.inputTokens'), t('admin.usage.outputTokens'),
      t('admin.usage.cacheReadTokens'), t('admin.usage.cacheCreationTokens'),
      t('admin.usage.inputCost'), t('admin.usage.outputCost'),
      t('admin.usage.cacheReadCost'), t('admin.usage.cacheCreationCost'),
      t('usage.rate'), t('usage.accountMultiplier'), t('usage.original'), t('usage.userBilled'), t('usage.accountBilled'),
      t('usage.firstToken'), t('usage.duration'),
      t('admin.usage.requestId'), t('usage.userAgent'), t('admin.usage.ipAddress')
    ]
    const ws = XLSX.utils.aoa_to_sheet([headers])
    while (true) {
      const res = await adminUsageAPI.list(
        buildUsageListParams(p, 100, true),
        { signal: c.signal placeholder
      )
      if (c.signal.aborted) break; if (p === 1) { total = res.total; exportProgress.total = total placeholder
      const rows = (res.items || []).map((log: AdminUsageLog) => [
        log.created_at, log.user?.email || '', log.api_key?.name || '', log.account?.name || '', log.model,
        log.upstream_model || log.model, log.upstream_response_model || '', log.upstream_model_mismatch == null ? '' : t(log.upstream_model_mismatch ? 'common.yes' : 'common.no'), formatReasoningEffort(log.reasoning_effort), log.group?.name || '',
        log.inbound_endpoint || '', log.upstream_endpoint || '', getRequestTypeLabel(log),
        log.input_tokens, log.output_tokens, log.cache_read_tokens, log.cache_creation_tokens,
        log.input_cost?.toFixed(6) || '0.000000', log.output_cost?.toFixed(6) || '0.000000',
        log.cache_read_cost?.toFixed(6) || '0.000000', log.cache_creation_cost?.toFixed(6) || '0.000000',
        log.rate_multiplier?.toPrecision(4) || '1.00', (log.account_rate_multiplier ?? 1).toPrecision(4),
        log.total_cost?.toFixed(6) || '0.000000', log.actual_cost?.toFixed(6) || '0.000000',
        ((log.account_stats_cost ?? log.total_cost) * (log.account_rate_multiplier ?? 1)).toFixed(6), log.first_token_ms ?? '', log.duration_ms,
        log.request_id || '', log.user_agent || '', log.ip_address || ''
      ])
      if (rows.length) {
        XLSX.utils.sheet_add_aoa(ws, rows, { origin: -1 placeholder)
      placeholder
      exportedCount += rows.length
      exportProgress.current = exportedCount
      exportProgress.progress = total > 0 ? Math.min(100, Math.round(exportedCount / total * 100)) : 0
      if (exportedCount >= total || res.items.length < 100) break; p++
    placeholder
    if(!c.signal.aborted) {
      const wb = XLSX.utils.book_new()
      XLSX.utils.book_append_sheet(wb, ws, 'Usage')
      saveAs(new Blob([XLSX.write(wb, { bookType: 'xlsx', type: 'array' placeholder)], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' placeholder), `usage_${filters.value.start_dateplaceholder_to_${filters.value.end_dateplaceholder.xlsx`)
      appStore.showSuccess(t('usage.exportSuccess'))
    placeholder
  placeholder catch (error) { console.error('Failed to export:', error); appStore.showError('Export Failed') placeholder
  finally { if(exportAbortController === c) { exportAbortController = null; exporting.value = false; exportProgress.show = false placeholder placeholder
placeholder

// Column visibility
const ALWAYS_VISIBLE = ['user', 'created_at']
const DEFAULT_HIDDEN_COLUMNS = ['reasoning_effort', 'request_id', 'user_agent']
const HIDDEN_COLUMNS_KEY = 'usage-hidden-columns'
const HIDDEN_COLUMNS_VERSION_KEY = 'usage-hidden-columns-version'
const HIDDEN_COLUMNS_CURRENT_VERSION = 'request-id-hidden-by-default'

const allColumns = computed(() => [
  { key: 'user', label: t('admin.usage.user'), sortable: false placeholder,
  { key: 'api_key', label: t('usage.apiKeyFilter'), sortable: false placeholder,
  { key: 'account', label: t('admin.usage.account'), sortable: false placeholder,
  { key: 'model', label: t('usage.model'), sortable: true placeholder,
  { key: 'reasoning_effort', label: t('usage.reasoningEffort'), sortable: false placeholder,
  { key: 'endpoint', label: t('usage.endpoint'), sortable: false placeholder,
  { key: 'group', label: t('admin.usage.group'), sortable: false placeholder,
  { key: 'stream', label: t('usage.type'), sortable: false placeholder,
  { key: 'billing_mode', label: t('admin.usage.billingMode'), sortable: false placeholder,
  { key: 'tokens', label: t('usage.tokens'), sortable: false placeholder,
  { key: 'cost', label: t('usage.cost'), sortable: false placeholder,
  { key: 'latency', label: t('usage.latency'), sortable: false placeholder,
  { key: 'created_at', label: t('usage.time'), sortable: true placeholder,
  { key: 'request_id', label: t('admin.usage.requestId'), sortable: false placeholder,
  { key: 'user_agent', label: t('usage.userAgent'), sortable: false placeholder,
  { key: 'ip_address', label: t('admin.usage.ipAddress'), sortable: false placeholder
])

const hiddenColumns = reactive<Set<string>>(new Set())

const toggleableColumns = computed(() =>
  allColumns.value.filter(col => !ALWAYS_VISIBLE.includes(col.key))
)

const visibleColumns = computed(() =>
  allColumns.value.filter(col =>
    ALWAYS_VISIBLE.includes(col.key) || !hiddenColumns.has(col.key)
  )
)

const isColumnVisible = (key: string) => !hiddenColumns.has(key)

const toggleColumn = (key: string) => {
  if (hiddenColumns.has(key)) {
    hiddenColumns.delete(key)
  placeholder else {
    hiddenColumns.add(key)
  placeholder
  try {
    localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
    localStorage.setItem(HIDDEN_COLUMNS_VERSION_KEY, HIDDEN_COLUMNS_CURRENT_VERSION)
  placeholder catch (e) {
    console.error('Failed to save columns:', e)
  placeholder
placeholder

// ---- 错误请求 tab 列设置(与用量明细同机制,独立存储) ----
const ERR_ALWAYS_VISIBLE = ['user', 'status', 'created_at', 'actions']
const ERR_DEFAULT_HIDDEN_COLUMNS = ['user_agent']
const ERR_HIDDEN_COLUMNS_KEY = 'usage-error-hidden-columns'

// key 集合须与 OpsErrorLogTable 内部 allColumns 一致
const errAllColumns = computed(() => [
  { key: 'user', label: t('admin.ops.errorLog.user') placeholder,
  { key: 'api_key', label: t('admin.ops.errorLog.apiKey') placeholder,
  { key: 'account', label: t('admin.ops.errorLog.account') placeholder,
  { key: 'platform', label: t('admin.ops.errorLog.platform') placeholder,
  { key: 'model', label: t('admin.ops.errorLog.model') placeholder,
  { key: 'endpoint', label: t('admin.ops.errorLog.endpoint') placeholder,
  { key: 'group', label: t('admin.ops.errorLog.group') placeholder,
  { key: 'type', label: t('admin.ops.errorLog.type') placeholder,
  { key: 'category', label: t('usage.errors.category') placeholder,
  { key: 'status', label: t('admin.ops.errorLog.status') placeholder,
  { key: 'message', label: t('admin.ops.errorLog.message') placeholder,
  { key: 'created_at', label: t('admin.ops.errorLog.time') placeholder,
  { key: 'user_agent', label: t('usage.userAgent') placeholder,
  { key: 'client_ip', label: t('admin.ops.errorLog.ip') placeholder,
  { key: 'actions', label: t('admin.ops.errorLog.action') placeholder,
])

const errHiddenColumns = reactive<Set<string>>(new Set())

const errToggleableColumns = computed(() =>
  errAllColumns.value.filter(col => !ERR_ALWAYS_VISIBLE.includes(col.key))
)

const errVisibleColumnKeys = computed(() =>
  errAllColumns.value
    .filter(col => ERR_ALWAYS_VISIBLE.includes(col.key) || !errHiddenColumns.has(col.key))
    .map(col => col.key)
)

const toggleErrColumn = (key: string) => {
  if (errHiddenColumns.has(key)) {
    errHiddenColumns.delete(key)
  placeholder else {
    errHiddenColumns.add(key)
  placeholder
  try {
    localStorage.setItem(ERR_HIDDEN_COLUMNS_KEY, JSON.stringify([...errHiddenColumns]))
  placeholder catch (e) {
    console.error('Failed to save error columns:', e)
  placeholder
placeholder

const loadSavedErrColumns = () => {
  try {
    const saved = localStorage.getItem(ERR_HIDDEN_COLUMNS_KEY)
    const keys = saved ? (JSON.parse(saved) as string[]) : ERR_DEFAULT_HIDDEN_COLUMNS
    keys.forEach((key) => errHiddenColumns.add(key))
  placeholder catch {
    ERR_DEFAULT_HIDDEN_COLUMNS.forEach((key) => errHiddenColumns.add(key))
  placeholder
placeholder

// 列设置下拉按当前 tab 分发
const currentToggleableColumns = computed(() =>
  activeTab.value === 'errors' ? errToggleableColumns.value : toggleableColumns.value
)
const isCurrentColumnVisible = (key: string) =>
  activeTab.value === 'errors' ? !errHiddenColumns.has(key) : isColumnVisible(key)
const toggleCurrentColumn = (key: string) =>
  activeTab.value === 'errors' ? toggleErrColumn(key) : toggleColumn(key)

const loadSavedColumns = () => {
  try {
    const saved = localStorage.getItem(HIDDEN_COLUMNS_KEY)
    if (saved) {
      (JSON.parse(saved) as string[]).forEach((key) => {
        hiddenColumns.add(key)
      placeholder)
      if (localStorage.getItem(HIDDEN_COLUMNS_VERSION_KEY) !== HIDDEN_COLUMNS_CURRENT_VERSION) {
        hiddenColumns.add('request_id')
        localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
        localStorage.setItem(HIDDEN_COLUMNS_VERSION_KEY, HIDDEN_COLUMNS_CURRENT_VERSION)
      placeholder
    placeholder else {
      DEFAULT_HIDDEN_COLUMNS.forEach((key) => {
        hiddenColumns.add(key)
      placeholder)
      localStorage.setItem(HIDDEN_COLUMNS_VERSION_KEY, HIDDEN_COLUMNS_CURRENT_VERSION)
    placeholder
  placeholder catch {
    DEFAULT_HIDDEN_COLUMNS.forEach((key) => {
      hiddenColumns.add(key)
    placeholder)
  placeholder
placeholder

// Detail tabs
type DetailTab = 'usage' | 'errors' | 'ranking'
const activeTab = ref<DetailTab>('usage')
const detailTabs = computed(() => [
  { key: 'usage' as const, label: t('usage.tabs.usage'), icon: 'document' as const placeholder,
  { key: 'errors' as const, label: t('usage.tabs.errors'), icon: 'exclamationTriangle' as const placeholder,
  { key: 'ranking' as const, label: t('usage.tabs.ranking'), icon: 'chart' as const placeholder,
])
const usageFiltersRef = ref<InstanceType<typeof UsageFilters> | null>(null)
const rankingMounted = ref(false)
const rankingRef = ref<InstanceType<typeof UserTokenRanking> | null>(null)

const switchTab = (tab: DetailTab) => {
  activeTab.value = tab
  if (tab === 'errors' && errRows.value.length === 0) loadAdminErrors()
  if (tab === 'ranking') rankingMounted.value = true
placeholder

// Error tab state
const errRows = ref<OpsErrorLog[]>([])
const errLoading = ref(false)
const errPage = ref(1)
const errPageSize = ref(20)
const errTotal = ref(0)
const errSortBy = ref('created_at')
const errSortOrder = ref<'asc' | 'desc'>('desc')
const showErrorModal = ref(false)
const selectedErrorId = ref<number | null>(null)

// 注意：'YYYY-MM-DDT00:00:00' 无时区后缀，按本地时区解析后再转 UTC——与页面其它日期处理语义一致，刻意如此，勿改成 'T00:00:00Z'
const toRFC3339 = (d: string | undefined, endOfDay = false): string | undefined =>
  d ? new Date(d + (endOfDay ? 'T23:59:59.999' : 'T00:00:00')).toISOString() : undefined

const loadAdminErrors = async () => {
  errLoading.value = true
  try {
    const resp = await listErrorLogs({
      page: errPage.value,
      page_size: errPageSize.value,
      view: 'all',
      start_time: toRFC3339(filters.value.start_date),
      end_time: toRFC3339(filters.value.end_date, true),
      user_id: filters.value.user_id ?? undefined,
      api_key_id: filters.value.api_key_id ?? undefined,
      account_id: filters.value.account_id ?? undefined,
      group_id: filters.value.group_id ?? undefined,
      model: filters.value.model || undefined,
      phase: filters.value.error_phase || undefined,
      category: filters.value.error_category || undefined,
      status_codes: filters.value.status_code != null ? String(filters.value.status_code) : undefined,
      sort_by: errSortBy.value,
      sort_order: errSortOrder.value,
    placeholder)
    errRows.value = resp.items
    errTotal.value = resp.total
  placeholder catch (error) {
    console.error('Failed to load admin errors:', error)
    appStore.showError(t('usage.errors.failedToLoad'))
  placeholder finally {
    errLoading.value = false
  placeholder
placeholder

const onErrSort = (sortBy: string, sortOrder: 'asc' | 'desc') => {
  errSortBy.value = sortBy
  errSortOrder.value = sortOrder
  errPage.value = 1
  loadAdminErrors()
placeholder
const onErrPage = (p: number) => { errPage.value = p; loadAdminErrors() placeholder
const onErrPageSize = (s: number) => { errPageSize.value = s; errPage.value = 1; loadAdminErrors() placeholder
const openError = (id: number) => { selectedErrorId.value = id; showErrorModal.value = true placeholder

const showColumnDropdown = ref(false)
const columnDropdownRef = ref<HTMLElement | null>(null)

const handleColumnClickOutside = (event: MouseEvent) => {
  if (columnDropdownRef.value && !columnDropdownRef.value.contains(event.target as HTMLElement)) {
    showColumnDropdown.value = false
  placeholder
placeholder

onMounted(() => {
  applyRouteQueryFilters()
  void loadRouteUserFilterLabel()
  loadLogs()
  loadStats()
  loadModelStats(modelDistributionSource.value, true)
  window.setTimeout(() => {
    void loadChartData()
  placeholder, 120)
  loadSavedColumns()
  loadSavedErrColumns()
  document.addEventListener('click', handleColumnClickOutside)
placeholder)
onUnmounted(() => { abortController?.abort(); exportAbortController?.abort(); document.removeEventListener('click', handleColumnClickOutside) placeholder)

watch(modelDistributionSource, (source) => {
  void loadModelStats(source)
placeholder)

defineExpose({ requestedModelStats, refreshData placeholder)
</script>
