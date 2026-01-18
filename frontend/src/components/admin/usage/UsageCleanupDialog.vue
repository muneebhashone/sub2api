<template>
  <BaseDialog :show="show" :title="t('admin.usage.cleanup.title')" width="wide" @close="handleClose">
    <div class="space-y-4">
      <UsageFilters
        v-model="localFilters"
        v-model:startDate="localStartDate"
        v-model:endDate="localEndDate"
        :exporting="false"
        :show-actions="false"
        @change="noop"
      />

      <div class="rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-700 dark:border-amber-500/30 dark:bg-amber-500/10 dark:text-amber-200">
        {{ t('admin.usage.cleanup.warning') placeholderplaceholder
      </div>

      <div class="rounded-xl border border-gray-200 p-4 dark:border-dark-700">
        <div class="flex items-center justify-between">
          <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">
            {{ t('admin.usage.cleanup.recentTasks') placeholderplaceholder
          </h4>
          <button type="button" class="btn btn-ghost btn-sm" @click="loadTasks">
            {{ t('common.refresh') placeholderplaceholder
          </button>
        </div>

        <div class="mt-3 space-y-2">
          <div v-if="tasksLoading" class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('admin.usage.cleanup.loadingTasks') placeholderplaceholder
          </div>
          <div v-else-if="tasks.length === 0" class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('admin.usage.cleanup.noTasks') placeholderplaceholder
          </div>
          <div v-else class="space-y-2">
            <div
              v-for="task in tasks"
              :key="task.id"
              class="flex flex-col gap-2 rounded-lg border border-gray-100 px-3 py-2 text-sm text-gray-600 dark:border-dark-700 dark:text-gray-300"
            >
              <div class="flex flex-wrap items-center justify-between gap-2">
                <div class="flex items-center gap-2">
                  <span :class="statusClass(task.status)" class="rounded-full px-2 py-0.5 text-xs font-semibold">
                    {{ statusLabel(task.status) placeholderplaceholder
                  </span>
                  <span class="text-xs text-gray-400">#{{ task.id placeholderplaceholder</span>
                  <button
                    v-if="canCancel(task)"
                    type="button"
                    class="btn btn-ghost btn-xs text-rose-600 hover:text-rose-700 dark:text-rose-300"
                    @click="openCancelConfirm(task)"
                  >
                    {{ t('admin.usage.cleanup.cancel') placeholderplaceholder
                  </button>
                </div>
                <div class="text-xs text-gray-400">
                  {{ formatDateTime(task.created_at) placeholderplaceholder
                </div>
              </div>
              <div class="flex flex-wrap items-center gap-4 text-xs text-gray-500 dark:text-gray-400">
                <span>{{ t('admin.usage.cleanup.range') placeholderplaceholder: {{ formatRange(task) placeholderplaceholder</span>
                <span>{{ t('admin.usage.cleanup.deletedRows') placeholderplaceholder: {{ task.deleted_rows.toLocaleString() placeholderplaceholder</span>
              </div>
              <div v-if="task.error_message" class="text-xs text-rose-500">
                {{ task.error_message placeholderplaceholder
              </div>
            </div>
          </div>
        </div>

        <Pagination
          v-if="tasksTotal > tasksPageSize"
          class="mt-4"
          :total="tasksTotal"
          :page="tasksPage"
          :page-size="tasksPageSize"
          :page-size-options="[5]"
          :show-page-size-selector="false"
          :show-jump="true"
          @update:page="handleTaskPageChange"
          @update:pageSize="handleTaskPageSizeChange"
        />
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button type="button" class="btn btn-secondary" @click="handleClose">
          {{ t('common.cancel') placeholderplaceholder
        </button>
        <button type="button" class="btn btn-danger" :disabled="submitting" @click="openConfirm">
          {{ submitting ? t('admin.usage.cleanup.submitting') : t('admin.usage.cleanup.submit') placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>

  <ConfirmDialog
    :show="confirmVisible"
    :title="t('admin.usage.cleanup.confirmTitle')"
    :message="t('admin.usage.cleanup.confirmMessage')"
    :confirm-text="t('admin.usage.cleanup.confirmSubmit')"
    danger
    @confirm="submitCleanup"
    @cancel="confirmVisible = false"
  />

  <ConfirmDialog
    :show="cancelConfirmVisible"
    :title="t('admin.usage.cleanup.cancelConfirmTitle')"
    :message="t('admin.usage.cleanup.cancelConfirmMessage')"
    :confirm-text="t('admin.usage.cleanup.cancelConfirm')"
    danger
    @confirm="cancelTask"
    @cancel="cancelConfirmVisible = false"
  />
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Pagination from '@/components/common/Pagination.vue'
import UsageFilters from '@/components/admin/usage/UsageFilters.vue'
import { adminUsageAPI placeholder from '@/api/admin/usage'
import type { AdminUsageQueryParams, UsageCleanupTask, CreateUsageCleanupTaskRequest placeholder from '@/api/admin/usage'

interface Props {
  show: boolean
  filters: AdminUsageQueryParams
  startDate: string
  endDate: string
placeholder

const props = defineProps<Props>()
const emit = defineEmits(['close'])

const { t placeholder = useI18n()
const appStore = useAppStore()

const localFilters = ref<AdminUsageQueryParams>({placeholder)
const localStartDate = ref('')
const localEndDate = ref('')

const tasks = ref<UsageCleanupTask[]>([])
const tasksLoading = ref(false)
const tasksPage = ref(1)
const tasksPageSize = ref(5)
const tasksTotal = ref(0)
const submitting = ref(false)
const confirmVisible = ref(false)
const cancelConfirmVisible = ref(false)
const canceling = ref(false)
const cancelTarget = ref<UsageCleanupTask | null>(null)
let pollTimer: number | null = null

const noop = () => {placeholder

const resetFilters = () => {
  localFilters.value = { ...props.filters placeholder
  localStartDate.value = props.startDate
  localEndDate.value = props.endDate
  localFilters.value.start_date = localStartDate.value
  localFilters.value.end_date = localEndDate.value
  tasksPage.value = 1
  tasksTotal.value = 0
placeholder

const startPolling = () => {
  stopPolling()
  pollTimer = window.setInterval(() => {
    loadTasks()
  placeholder, 10000)
placeholder

const stopPolling = () => {
  if (pollTimer !== null) {
    window.clearInterval(pollTimer)
    pollTimer = null
  placeholder
placeholder

const handleClose = () => {
  stopPolling()
  confirmVisible.value = false
  cancelConfirmVisible.value = false
  canceling.value = false
  cancelTarget.value = null
  submitting.value = false
  emit('close')
placeholder

const statusLabel = (status: string) => {
  const map: Record<string, string> = {
    pending: t('admin.usage.cleanup.status.pending'),
    running: t('admin.usage.cleanup.status.running'),
    succeeded: t('admin.usage.cleanup.status.succeeded'),
    failed: t('admin.usage.cleanup.status.failed'),
    canceled: t('admin.usage.cleanup.status.canceled')
  placeholder
  return map[status] || status
placeholder

const statusClass = (status: string) => {
  const map: Record<string, string> = {
    pending: 'bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-200',
    running: 'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-200',
    succeeded: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-200',
    failed: 'bg-rose-100 text-rose-700 dark:bg-rose-500/20 dark:text-rose-200',
    canceled: 'bg-gray-200 text-gray-600 dark:bg-dark-600 dark:text-gray-300'
  placeholder
  return map[status] || 'bg-gray-100 text-gray-600'
placeholder

const formatDateTime = (value?: string | null) => {
  if (!value) return '--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
placeholder

const formatRange = (task: UsageCleanupTask) => {
  const start = formatDateTime(task.filters.start_time)
  const end = formatDateTime(task.filters.end_time)
  return `${startplaceholder ~ ${endplaceholder`
placeholder

const getUserTimezone = () => {
  try {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
  placeholder catch {
    return 'UTC'
  placeholder
placeholder

const loadTasks = async () => {
  if (!props.show) return
  tasksLoading.value = true
  try {
    const res = await adminUsageAPI.listCleanupTasks({
      page: tasksPage.value,
      page_size: tasksPageSize.value
    placeholder)
    tasks.value = res.items || []
    tasksTotal.value = res.total || 0
    if (res.page) {
      tasksPage.value = res.page
    placeholder
    if (res.page_size) {
      tasksPageSize.value = res.page_size
    placeholder
  placeholder catch (error) {
    console.error('Failed to load cleanup tasks:', error)
    appStore.showError(t('admin.usage.cleanup.loadFailed'))
  placeholder finally {
    tasksLoading.value = false
  placeholder
placeholder

const handleTaskPageChange = (page: number) => {
  tasksPage.value = page
  loadTasks()
placeholder

const handleTaskPageSizeChange = (size: number) => {
  if (!Number.isFinite(size) || size <= 0) return
  tasksPageSize.value = size
  tasksPage.value = 1
  loadTasks()
placeholder

const openConfirm = () => {
  confirmVisible.value = true
placeholder

const canCancel = (task: UsageCleanupTask) => {
  return task.status === 'pending' || task.status === 'running'
placeholder

const openCancelConfirm = (task: UsageCleanupTask) => {
  cancelTarget.value = task
  cancelConfirmVisible.value = true
placeholder

const buildPayload = (): CreateUsageCleanupTaskRequest | null => {
  if (!localStartDate.value || !localEndDate.value) {
    appStore.showError(t('admin.usage.cleanup.missingRange'))
    return null
  placeholder

  const payload: CreateUsageCleanupTaskRequest = {
    start_date: localStartDate.value,
    end_date: localEndDate.value,
    timezone: getUserTimezone()
  placeholder

  if (localFilters.value.user_id && localFilters.value.user_id > 0) {
    payload.user_id = localFilters.value.user_id
  placeholder
  if (localFilters.value.api_key_id && localFilters.value.api_key_id > 0) {
    payload.api_key_id = localFilters.value.api_key_id
  placeholder
  if (localFilters.value.account_id && localFilters.value.account_id > 0) {
    payload.account_id = localFilters.value.account_id
  placeholder
  if (localFilters.value.group_id && localFilters.value.group_id > 0) {
    payload.group_id = localFilters.value.group_id
  placeholder
  if (localFilters.value.model) {
    payload.model = localFilters.value.model
  placeholder
  if (localFilters.value.stream !== null && localFilters.value.stream !== undefined) {
    payload.stream = localFilters.value.stream
  placeholder
  if (localFilters.value.billing_type !== null && localFilters.value.billing_type !== undefined) {
    payload.billing_type = localFilters.value.billing_type
  placeholder

  return payload
placeholder

const submitCleanup = async () => {
  const payload = buildPayload()
  if (!payload) {
    confirmVisible.value = false
    return
  placeholder
  submitting.value = true
  confirmVisible.value = false
  try {
    await adminUsageAPI.createCleanupTask(payload)
    appStore.showSuccess(t('admin.usage.cleanup.submitSuccess'))
    loadTasks()
  placeholder catch (error) {
    console.error('Failed to create cleanup task:', error)
    appStore.showError(t('admin.usage.cleanup.submitFailed'))
  placeholder finally {
    submitting.value = false
  placeholder
placeholder

const cancelTask = async () => {
  const task = cancelTarget.value
  if (!task) {
    cancelConfirmVisible.value = false
    return
  placeholder
  canceling.value = true
  cancelConfirmVisible.value = false
  try {
    await adminUsageAPI.cancelCleanupTask(task.id)
    appStore.showSuccess(t('admin.usage.cleanup.cancelSuccess'))
    loadTasks()
  placeholder catch (error) {
    console.error('Failed to cancel cleanup task:', error)
    appStore.showError(t('admin.usage.cleanup.cancelFailed'))
  placeholder finally {
    canceling.value = false
    cancelTarget.value = null
  placeholder
placeholder

watch(
  () => props.show,
  (show) => {
    if (show) {
      resetFilters()
      loadTasks()
      startPolling()
    placeholder else {
      stopPolling()
    placeholder
  placeholder
)

onUnmounted(() => {
  stopPolling()
placeholder)
</script>
