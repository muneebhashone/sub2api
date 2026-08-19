<template>
  <AppLayout>
    <div class="w-full min-w-0 space-y-6 pb-8">
      <header
        class="page-header mb-0 rounded-3xl bg-white p-5 shadow-sm ring-1 ring-gray-900/5 dark:bg-dark-800 dark:ring-dark-700 sm:p-6"
      >
        <h1 class="page-title flex items-center gap-2 text-xl font-black text-gray-900 dark:text-white">
          <span class="inline-flex h-8 w-8 items-center justify-center rounded-xl bg-blue-50 text-blue-500 dark:bg-blue-900/30 dark:text-blue-400">
            <Icon name="chart" size="sm" />
          </span>
          {{ t('admin.channelMonitor.title') placeholderplaceholder
        </h1>
        <p class="page-description mt-1.5 text-xs text-gray-500 dark:text-gray-400">
          {{
            isV1Mode
              ? t('channelMonitorV2.admin.descriptionV1')
              : t('channelMonitorV2.admin.descriptionV2')
          placeholderplaceholder
        </p>
        <div class="mt-4 border-t border-gray-100 pt-4 dark:border-dark-700">
          <div
            class="tabs inline-flex w-full max-w-xl flex-wrap sm:w-auto"
            role="tablist"
            :aria-label="t('channelMonitorV2.admin.tabAria')"
          >
            <button
              type="button"
              role="tab"
              class="tab flex-1 sm:flex-none"
              :class="adminMonitorTab === 'v2' ? 'tab-active' : ''"
              :aria-selected="adminMonitorTab === 'v2'"
              @click="adminMonitorTab = 'v2'"
            >
              {{ t('channelMonitorV2.admin.tabV2') placeholderplaceholder
            </button>
            <button
              type="button"
              role="tab"
              class="tab flex-1 sm:flex-none"
              :class="adminMonitorTab === 'legacy' ? 'tab-active' : ''"
              :aria-selected="adminMonitorTab === 'legacy'"
              @click="adminMonitorTab = 'legacy'"
            >
              {{ isV1Mode ? t('channelMonitorV2.admin.tabV1Active') : t('channelMonitorV2.admin.tabV1History') placeholderplaceholder
            </button>
          </div>
        </div>
      </header>

      <MonitorSettingsPanel v-if="adminMonitorTab === 'v2'" />

      <TablePageLayout v-else>
      <template #filters>
        <MonitorFiltersBar
          v-model:search="searchQuery"
          v-model:provider="providerFilter"
          v-model:enabled="enabledFilter"
          :loading="loading"
          @reload="reload"
          @create="openCreateDialog"
          @manage-templates="showTemplateManager = true"
          @search-input="handleSearch"
        />
      </template>

      <template #table>
        <DataTable :columns="columns" :data="monitors" :loading="loading">
          <template #cell-name="{ row, value placeholder">
            <div class="flex items-center gap-1.5">
              <span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
              <HelpTooltip v-if="row.api_key_decrypt_failed" :content="t('admin.channelMonitor.apiKeyDecryptFailed')">
                <Icon name="exclamationTriangle" size="sm" class="text-red-500" />
              </HelpTooltip>
            </div>
          </template>

          <template #cell-provider="{ row placeholder">
            <span class="inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium" :class="providerBadgeClass(row.provider)">
              {{ providerLabel(row.provider) placeholderplaceholder
            </span>
            <!-- 三种检测模式并列展示，quota 系配额数据源与纯探活一眼可分 -->
            <span class="ml-1 inline-flex items-center rounded-md px-2 py-0.5 text-xs font-medium" :class="checkModeBadgeClass(row.check_mode)">
              {{ checkModeLabel(row.check_mode) placeholderplaceholder
            </span>
          </template>

          <template #cell-primary_model="{ row placeholder">
            <MonitorPrimaryModelCell :row="row" />
          </template>

          <template #cell-availability_7d="{ row placeholder">
            <span class="text-sm text-gray-900 dark:text-gray-100">{{ formatAvailability(row) placeholderplaceholder</span>
          </template>

          <template #cell-latency="{ row placeholder">
            <span class="text-sm text-gray-900 dark:text-gray-100">{{ formatLatency(row.primary_latency_ms) placeholderplaceholder</span>
          </template>

          <template #cell-enabled="{ row placeholder">
            <Toggle :modelValue="row.enabled" @update:modelValue="toggleEnabled(row)" />
          </template>

          <template #cell-actions="{ row placeholder">
            <MonitorActionsCell
              :row="row"
              :running="runningId === row.id"
              :duplicating="duplicatingIds.has(row.id)"
              @run="handleRunNow"
              @duplicate="handleDuplicate"
              @edit="openEditDialog"
              @delete="handleDelete"
            />
          </template>

          <template #empty>
            <EmptyState
              :title="t('admin.channelMonitor.noMonitorsYet')"
              :description="t('admin.channelMonitor.createFirstMonitor')"
              :action-text="t('admin.channelMonitor.createButton')"
              @action="openCreateDialog"
            />
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="onPageChange"
          @update:pageSize="onPageSizeChange"
        />
      </template>
      </TablePageLayout>
    </div>

    <MonitorFormDialog
      :show="showDialog"
      :monitor="editing"
      @close="closeDialog"
      @saved="reload"
    />

    <MonitorTemplateManagerDialog
      :show="showTemplateManager"
      @close="showTemplateManager = false"
      @updated="reload"
    />

    <MonitorRunResultDialog
      :show="showRunResult"
      :results="runResults"
      @close="showRunResult = false"
    />

    <ConfirmDialog
      :show="showDeleteDialog"
      :title="t('common.delete')"
      :message="deleteConfirmMessage"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="confirmDelete"
      @cancel="showDeleteDialog = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import { adminAPI placeholder from '@/api/admin'
import type {
  ChannelMonitor,
  CheckResult,
  ListParams,
  Provider,
placeholder from '@/api/admin/channelMonitor'
import type { Column placeholder from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import HelpTooltip from '@/components/common/HelpTooltip.vue'
import Icon from '@/components/icons/Icon.vue'
import Toggle from '@/components/common/Toggle.vue'
import MonitorFiltersBar from '@/components/admin/monitor/MonitorFiltersBar.vue'
import MonitorFormDialog from '@/components/admin/monitor/MonitorFormDialog.vue'
import MonitorTemplateManagerDialog from '@/components/admin/monitor/MonitorTemplateManagerDialog.vue'
import MonitorRunResultDialog from '@/components/admin/monitor/MonitorRunResultDialog.vue'
import MonitorPrimaryModelCell from '@/components/admin/monitor/MonitorPrimaryModelCell.vue'
import MonitorActionsCell from '@/components/admin/monitor/MonitorActionsCell.vue'
import { getPersistedPageSize placeholder from '@/composables/usePersistedPageSize'
import { useChannelMonitorFormat placeholder from '@/composables/useChannelMonitorFormat'
import MonitorSettingsPanel from '@/features/channel-monitor-v2/MonitorSettingsPanel.vue'
import { isChannelMonitorV1Mode placeholder from '@/utils/featureFlags'

const { t placeholder = useI18n()
const appStore = useAppStore()
const isV1Mode = computed(() => isChannelMonitorV1Mode())
const adminMonitorTab = ref<'v2' | 'legacy'>(isChannelMonitorV1Mode() ? 'legacy' : 'v2')
const {
  providerLabel,
  providerBadgeClass,
  checkModeLabel,
  checkModeBadgeClass,
  formatLatency,
  formatAvailability,
placeholder = useChannelMonitorFormat()

const monitors = ref<ChannelMonitor[]>([])
const loading = ref(false)
const runningId = ref<number | null>(null)
const searchQuery = ref('')
const providerFilter = ref<Provider | ''>('')
const enabledFilter = ref<'' | 'true' | 'false'>('')
const pagination = reactive({ page: 1, page_size: getPersistedPageSize(), total: 0 placeholder)

const showDialog = ref(false)
const showTemplateManager = ref(false)
const editing = ref<ChannelMonitor | null>(null)
const showDeleteDialog = ref(false)
const deleting = ref<ChannelMonitor | null>(null)
const showRunResult = ref(false)
const runResults = ref<CheckResult[]>([])
const duplicatingIds = reactive(new Set<number>())

let abortController: AbortController | null = null
let searchTimeout: ReturnType<typeof setTimeout> | null = null

const columns = computed<Column[]>(() => [
  { key: 'name', label: t('admin.channelMonitor.columns.name'), sortable: false placeholder,
  { key: 'provider', label: t('admin.channelMonitor.columns.provider'), sortable: false placeholder,
  { key: 'primary_model', label: t('admin.channelMonitor.columns.primaryModel'), sortable: false placeholder,
  { key: 'availability_7d', label: t('admin.channelMonitor.columns.availability7d'), sortable: false placeholder,
  { key: 'latency', label: t('admin.channelMonitor.columns.latency'), sortable: false placeholder,
  { key: 'enabled', label: t('admin.channelMonitor.columns.enabled'), sortable: false placeholder,
  { key: 'actions', label: t('admin.channelMonitor.columns.actions'), sortable: false placeholder,
])

const deleteConfirmMessage = computed(() => {
  const name = deleting.value?.name || ''
  return t('admin.channelMonitor.deleteConfirm', { name placeholder)
placeholder)

async function reload() {
  if (abortController) abortController.abort()
  const ctrl = new AbortController()
  abortController = ctrl
  loading.value = true
  try {
    const params: ListParams = {
      page: pagination.page,
      page_size: pagination.page_size,
    placeholder
    if (providerFilter.value) params.provider = providerFilter.value
    if (enabledFilter.value === 'true') params.enabled = true
    if (enabledFilter.value === 'false') params.enabled = false
    if (searchQuery.value.trim()) params.search = searchQuery.value.trim()

    const res = await adminAPI.channelMonitor.list(params, { signal: ctrl.signal placeholder)
    if (ctrl.signal.aborted || abortController !== ctrl) return
    monitors.value = res.items || []
    pagination.total = res.total
  placeholder catch (err: unknown) {
    const e = err as { name?: string; code?: string placeholder
    if (e?.name === 'AbortError' || e?.code === 'ERR_CANCELED') return
    appStore.showError(extractApiErrorMessage(err, t('admin.channelMonitor.loadError')))
  placeholder finally {
    if (abortController === ctrl) {
      loading.value = false
      abortController = null
    placeholder
  placeholder
placeholder

function handleSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.page = 1
    reload()
  placeholder, 300)
placeholder

function onPageChange(page: number) {
  pagination.page = page
  reload()
placeholder

function onPageSizeChange(size: number) {
  pagination.page_size = size
  pagination.page = 1
  reload()
placeholder

function openCreateDialog() {
  editing.value = null
  showDialog.value = true
placeholder

function openEditDialog(row: ChannelMonitor) {
  editing.value = row
  showDialog.value = true
placeholder

function closeDialog() {
  showDialog.value = false
  editing.value = null
placeholder

async function toggleEnabled(row: ChannelMonitor) {
  const next = !row.enabled
  try {
    await adminAPI.channelMonitor.update(row.id, { enabled: next placeholder)
    row.enabled = next
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  placeholder
placeholder

async function handleRunNow(row: ChannelMonitor) {
  if (!isV1Mode.value) {
    appStore.showError(t('admin.channelMonitor.runFailed'))
    return
  placeholder
  if (runningId.value != null) return
  runningId.value = row.id
  try {
    const res = await adminAPI.channelMonitor.runNow(row.id)
    runResults.value = res.results || []
    showRunResult.value = true
    appStore.showSuccess(t('admin.channelMonitor.runSuccess'))
    // Refresh row to get latest status from backend
    void reload()
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.channelMonitor.runFailed')))
  placeholder finally {
    runningId.value = null
  placeholder
placeholder

async function handleDuplicate(row: ChannelMonitor) {
  if (row.api_key_decrypt_failed) {
    appStore.showError(t('admin.channelMonitor.duplicateKeyUnavailable'))
    return
  placeholder
  if (duplicatingIds.has(row.id)) return

  duplicatingIds.add(row.id)
  try {
    const duplicate = await adminAPI.channelMonitor.duplicate(row.id)
    appStore.showSuccess(t('admin.channelMonitor.duplicateSuccess', { name: duplicate.name placeholder))
    await reload()
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('admin.channelMonitor.duplicateFailed')))
  placeholder finally {
    duplicatingIds.delete(row.id)
  placeholder
placeholder

function handleDelete(row: ChannelMonitor) {
  deleting.value = row
  showDeleteDialog.value = true
placeholder

async function confirmDelete() {
  if (!deleting.value) return
  try {
    await adminAPI.channelMonitor.del(deleting.value.id)
    appStore.showSuccess(t('admin.channelMonitor.deleteSuccess'))
    showDeleteDialog.value = false
    deleting.value = null
    reload()
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  placeholder
placeholder

watch(adminMonitorTab, (tab) => {
  if (tab === 'legacy' && monitors.value.length === 0) void reload()
placeholder)
onMounted(() => {
  if (adminMonitorTab.value === 'legacy') void reload()
placeholder)
onUnmounted(() => {
  if (searchTimeout) clearTimeout(searchTimeout)
  abortController?.abort()
placeholder)
</script>
