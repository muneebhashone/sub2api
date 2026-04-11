<template>
  <BaseDialog
    :show="show"
    :title="t('admin.announcements.readStatus')"
    width="extra-wide"
    @close="handleClose"
  >
    <div class="space-y-4">
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex-1">
          <input
            v-model="search"
            type="text"
            class="input"
            :placeholder="t('admin.announcements.searchUsers')"
            @input="handleSearch"
          />
        </div>
        <button @click="load" :disabled="loading" class="btn btn-secondary" :title="t('common.refresh')">
          <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
        </button>
      </div>

      <DataTable
        :columns="columns"
        :data="items"
        :loading="loading"
        :server-side-sort="true"
        default-sort-key="email"
        default-sort-order="asc"
        @sort="handleSort"
      >
        <template #cell-email="{ value placeholder">
          <span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
        </template>

        <template #cell-balance="{ value placeholder">
          <span class="font-medium text-gray-900 dark:text-white">${{ Number(value ?? 0).toFixed(2) placeholderplaceholder</span>
        </template>

        <template #cell-eligible="{ value placeholder">
          <span :class="['badge', value ? 'badge-success' : 'badge-gray']">
            {{ value ? t('admin.announcements.eligible') : t('common.no') placeholderplaceholder
          </span>
        </template>

        <template #cell-read_at="{ value placeholder">
          <span class="text-sm text-gray-500 dark:text-dark-400">
            {{ value ? formatDateTime(value) : t('admin.announcements.unread') placeholderplaceholder
          </span>
        </template>
      </DataTable>

      <Pagination
        v-if="pagination.total > 0"
        :page="pagination.page"
        :total="pagination.total"
        :page-size="pagination.page_size"
        @update:page="handlePageChange"
        @update:pageSize="handlePageSizeChange"
      />
    </div>

    <template #footer>
      <div class="flex justify-end">
        <button type="button" class="btn btn-secondary" @click="handleClose">{{ t('common.close') placeholderplaceholder</button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, onUnmounted, reactive, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import { formatDateTime placeholder from '@/utils/format'
import type { AnnouncementUserReadStatus placeholder from '@/types'
import type { Column placeholder from '@/components/common/types'
import { getPersistedPageSize placeholder from '@/composables/usePersistedPageSize'

import BaseDialog from '@/components/common/BaseDialog.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import Icon from '@/components/icons/Icon.vue'

const { t placeholder = useI18n()
const appStore = useAppStore()

const props = defineProps<{
  show: boolean
  announcementId: number | null
placeholder>()

const emit = defineEmits<{
  (e: 'close'): void
placeholder>()

const loading = ref(false)
const search = ref('')

const pagination = reactive({
  page: 1,
  page_size: getPersistedPageSize(),
  total: 0,
  pages: 0
placeholder)

const sortState = reactive({
  sort_by: 'email',
  sort_order: 'asc' as 'asc' | 'desc'
placeholder)

const items = ref<AnnouncementUserReadStatus[]>([])

const columns = computed<Column[]>(() => [
  { key: 'email', label: t('common.email'), sortable: true placeholder,
  { key: 'username', label: t('admin.users.columns.username'), sortable: true placeholder,
  { key: 'balance', label: t('common.balance'), sortable: true placeholder,
  { key: 'eligible', label: t('admin.announcements.eligible') placeholder,
  { key: 'read_at', label: t('admin.announcements.readAt') placeholder
])

let currentController: AbortController | null = null
let searchDebounceTimer: number | null = null

function resetDialogState() {
  loading.value = false
  search.value = ''
  items.value = []
  pagination.page = 1
  pagination.total = 0
  pagination.pages = 0
  sortState.sort_by = 'email'
  sortState.sort_order = 'asc'
placeholder

function cancelPendingLoad(resetState = false) {
  if (searchDebounceTimer) {
    window.clearTimeout(searchDebounceTimer)
    searchDebounceTimer = null
  placeholder
  currentController?.abort()
  currentController = null
  if (resetState) {
    resetDialogState()
  placeholder
placeholder

async function load() {
  if (!props.show || !props.announcementId) return

  currentController?.abort()
  const requestController = new AbortController()
  currentController = requestController
  const { signal placeholder = requestController

  try {
    loading.value = true
    const res = await adminAPI.announcements.getReadStatus(
      props.announcementId,
      pagination.page,
      pagination.page_size,
      {
        search: search.value,
        sort_by: sortState.sort_by,
        sort_order: sortState.sort_order
      placeholder,
      { signal placeholder
    )

    if (signal.aborted || currentController !== requestController) return

    items.value = res.items
    pagination.total = res.total
    pagination.pages = res.pages
    pagination.page = res.page
    pagination.page_size = res.page_size
  placeholder catch (error: any) {
    if (
      signal.aborted ||
      currentController !== requestController ||
      error?.name === 'AbortError' ||
      error?.code === 'ERR_CANCELED'
    ) {
      return
    placeholder
    console.error('Failed to load read status:', error)
    appStore.showError(error.response?.data?.detail || t('admin.announcements.failedToLoadReadStatus'))
  placeholder finally {
    if (currentController === requestController) {
      loading.value = false
      currentController = null
    placeholder
  placeholder
placeholder

function handlePageChange(page: number) {
  pagination.page = page
  load()
placeholder

function handlePageSizeChange(pageSize: number) {
  pagination.page_size = pageSize
  pagination.page = 1
  load()
placeholder

function handleSort(key: string, order: 'asc' | 'desc') {
  sortState.sort_by = key
  sortState.sort_order = order
  pagination.page = 1
  load()
placeholder

function handleSearch() {
  if (searchDebounceTimer) window.clearTimeout(searchDebounceTimer)
  searchDebounceTimer = window.setTimeout(() => {
    pagination.page = 1
    load()
  placeholder, 300)
placeholder

function handleClose() {
  cancelPendingLoad(true)
  emit('close')
placeholder

watch(
  () => props.show,
  (v) => {
    if (!v) {
      cancelPendingLoad(true)
      return
    placeholder
    pagination.page = 1
    load()
  placeholder
)

watch(
  () => props.announcementId,
  () => {
    if (!props.show) return
    pagination.page = 1
    load()
  placeholder
)

onUnmounted(() => {
  cancelPendingLoad()
placeholder)
</script>
