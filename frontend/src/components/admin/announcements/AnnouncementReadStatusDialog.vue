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

      <DataTable :columns="columns" :data="items" :loading="loading">
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
import { computed, onMounted, reactive, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import { formatDateTime placeholder from '@/utils/format'
import type { AnnouncementUserReadStatus placeholder from '@/types'
import type { Column placeholder from '@/components/common/types'

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
  page_size: 20,
  total: 0,
  pages: 0
placeholder)

const items = ref<AnnouncementUserReadStatus[]>([])

const columns = computed<Column[]>(() => [
  { key: 'email', label: t('common.email') placeholder,
  { key: 'username', label: t('admin.users.columns.username') placeholder,
  { key: 'balance', label: t('common.balance') placeholder,
  { key: 'eligible', label: t('admin.announcements.eligible') placeholder,
  { key: 'read_at', label: t('admin.announcements.readAt') placeholder
])

let currentController: AbortController | null = null

async function load() {
  if (!props.show || !props.announcementId) return

  if (currentController) currentController.abort()
  currentController = new AbortController()

  try {
    loading.value = true
    const res = await adminAPI.announcements.getReadStatus(
      props.announcementId,
      pagination.page,
      pagination.page_size,
      search.value
    )

    items.value = res.items
    pagination.total = res.total
    pagination.pages = res.pages
    pagination.page = res.page
    pagination.page_size = res.page_size
  placeholder catch (error: any) {
    if (currentController.signal.aborted || error?.name === 'AbortError') return
    console.error('Failed to load read status:', error)
    appStore.showError(error.response?.data?.detail || t('admin.announcements.failedToLoadReadStatus'))
  placeholder finally {
    loading.value = false
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

let searchDebounceTimer: number | null = null
function handleSearch() {
  if (searchDebounceTimer) window.clearTimeout(searchDebounceTimer)
  searchDebounceTimer = window.setTimeout(() => {
    pagination.page = 1
    load()
  placeholder, 300)
placeholder

function handleClose() {
  emit('close')
placeholder

watch(
  () => props.show,
  (v) => {
    if (!v) return
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

onMounted(() => {
  // noop
placeholder)
</script>
