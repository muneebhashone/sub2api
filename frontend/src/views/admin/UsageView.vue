<template>
  <AppLayout>
    <div class="space-y-6">
      <UsageStatsCards :stats="usageStats" />
      <UsageFilters v-model="filters" v-model:startDate="startDate" v-model:endDate="endDate" :exporting="exporting" @change="applyFilters" @reset="resetFilters" @export="exportToExcel" />
      <UsageTable :data="usageLogs" :loading="loading" />
      <Pagination v-if="pagination.total > 0" :page="pagination.page" :total="pagination.total" :page-size="pagination.page_size" @update:page="handlePageChange" @update:pageSize="handlePageSizeChange" />
    </div>
  </AppLayout>
  <UsageExportProgress :show="exportProgress.show" :progress="exportProgress.progress" :current="exportProgress.current" :total="exportProgress.total" :estimated-time="exportProgress.estimatedTime" @cancel="cancelExport" />
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted placeholder from 'vue'
import { saveAs placeholder from 'file-saver'
import { useAppStore placeholder from '@/stores/app'; import { adminAPI placeholder from '@/api/admin'; import { adminUsageAPI placeholder from '@/api/admin/usage'
import AppLayout from '@/components/layout/AppLayout.vue'; import Pagination from '@/components/common/Pagination.vue'
import UsageStatsCards from '@/components/admin/usage/UsageStatsCards.vue'; import UsageFilters from '@/components/admin/usage/UsageFilters.vue'
import UsageTable from '@/components/admin/usage/UsageTable.vue'; import UsageExportProgress from '@/components/admin/usage/UsageExportProgress.vue'
import type { UsageLog placeholder from '@/types'; import type { AdminUsageStatsResponse, AdminUsageQueryParams placeholder from '@/api/admin/usage'

const appStore = useAppStore()
const usageStats = ref<AdminUsageStatsResponse | null>(null); const usageLogs = ref<UsageLog[]>([]); const loading = ref(false); const exporting = ref(false)
let abortController: AbortController | null = null; let exportAbortController: AbortController | null = null
const exportProgress = reactive({ show: false, progress: 0, current: 0, total: 0, estimatedTime: '' placeholder)

const formatLD = (d: Date) => d.toISOString().split('T')[0]
const now = new Date(); const weekAgo = new Date(Date.now() - 6 * 86400000)
const startDate = ref(formatLD(weekAgo)); const endDate = ref(formatLD(now))
const filters = ref<AdminUsageQueryParams>({ user_id: undefined, model: undefined, group_id: undefined, start_date: startDate.value, end_date: endDate.value placeholder)
const pagination = reactive({ page: 1, page_size: 20, total: 0 placeholder)

const loadLogs = async () => {
  abortController?.abort(); const c = new AbortController(); abortController = c; loading.value = true
  try {
    const res = await adminAPI.usage.list({ page: pagination.page, page_size: pagination.page_size, ...filters.value placeholder, { signal: c.signal placeholder)
    if(!c.signal.aborted) { usageLogs.value = res.items; pagination.total = res.total placeholder
  placeholder catch (error: any) { if(error?.name !== 'AbortError') console.error('Failed to load usage logs:', error) placeholder finally { if(abortController === c) loading.value = false placeholder
placeholder
const loadStats = async () => { try { const s = await adminAPI.usage.getStats(filters.value); usageStats.value = s placeholder catch (error) { console.error('Failed to load usage stats:', error) placeholder placeholder
const applyFilters = () => { pagination.page = 1; loadLogs(); loadStats() placeholder
const resetFilters = () => { startDate.value = formatLD(weekAgo); endDate.value = formatLD(now); filters.value = { start_date: startDate.value, end_date: endDate.value placeholder; applyFilters() placeholder
const handlePageChange = (p: number) => { pagination.page = p; loadLogs() placeholder
const handlePageSizeChange = (s: number) => { pagination.page_size = s; pagination.page = 1; loadLogs() placeholder
const cancelExport = () => exportAbortController?.abort()

const exportToExcel = async () => {
  if (exporting.value) return; exporting.value = true; exportProgress.show = true
  const c = new AbortController(); exportAbortController = c
  try {
    const all: UsageLog[] = []; let p = 1; let total = pagination.total
    while (true) {
      const res = await adminUsageAPI.list({ page: p, page_size: 100, ...filters.value placeholder, { signal: c.signal placeholder)
      if (c.signal.aborted) break; if (p === 1) { total = res.total; exportProgress.total = total placeholder
      if (res.items?.length) all.push(...res.items)
      exportProgress.current = all.length; exportProgress.progress = total > 0 ? Math.min(100, Math.round(all.length/total*100)) : 0
      if (all.length >= total || res.items.length < 100) break; p++
    placeholder
    if(!c.signal.aborted) {
      // 动态加载 xlsx，降低首屏包体并减少高危依赖的常驻暴露面。
      const XLSX = await import('xlsx')
      const ws = XLSX.utils.json_to_sheet(all); const wb = XLSX.utils.book_new(); XLSX.utils.book_append_sheet(wb, ws, 'Usage')
      saveAs(new Blob([XLSX.write(wb, { bookType: 'xlsx', type: 'array' placeholder)], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' placeholder), `usage_${Date.now()placeholder.xlsx`)
      appStore.showSuccess('Export Success')
    placeholder
  placeholder catch { appStore.showError('Export Failed') placeholder
  finally { if(exportAbortController === c) { exportAbortController = null; exporting.value = false; exportProgress.show = false placeholder placeholder
placeholder

onMounted(() => { loadLogs(); loadStats() placeholder)
onUnmounted(() => { abortController?.abort(); exportAbortController?.abort() placeholder)
</script>
