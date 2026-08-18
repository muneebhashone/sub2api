<template>
  <BaseDialog
    :show="show"
    :title="title"
    width="wide"
    @close="$emit('close')"
  >
    <div v-if="loading" class="py-8 text-center text-sm text-gray-500">
      {{ t('common.loading') placeholderplaceholder
    </div>
    <div v-else-if="!detail" class="py-8 text-center text-sm text-gray-500">
      {{ t('channelStatus.detailLoadError') placeholderplaceholder
    </div>
    <div v-else class="overflow-x-auto">
      <table class="w-full text-left text-sm">
        <thead class="border-b border-gray-200 dark:border-dark-700">
          <tr class="text-xs uppercase tracking-wider text-gray-500 dark:text-gray-400">
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.model') placeholderplaceholder</th>
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.latestStatus') placeholderplaceholder</th>
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.latestLatency') placeholderplaceholder</th>
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.availability7d') placeholderplaceholder</th>
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.availability15d') placeholderplaceholder</th>
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.availability30d') placeholderplaceholder</th>
            <th class="py-2 pr-3">{{ t('channelStatus.detailColumns.avgLatency7d') placeholderplaceholder</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="m in detail.models"
            :key="m.model"
            class="border-b border-gray-100 dark:border-dark-800"
          >
            <td class="py-2 pr-3 font-medium text-gray-900 dark:text-gray-100">{{ formatMonitorModel(m.model) placeholderplaceholder</td>
            <td class="py-2 pr-3">
              <span
                class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px]"
                :class="statusBadgeClass(m.latest_status)"
              >
                {{ statusLabel(m.latest_status) placeholderplaceholder
              </span>
            </td>
            <td class="py-2 pr-3 text-gray-700 dark:text-gray-300">{{ formatLatency(m.latest_latency_ms) placeholderplaceholder</td>
            <td class="py-2 pr-3 text-gray-700 dark:text-gray-300">{{ formatPercent(m.availability_7d) placeholderplaceholder</td>
            <td class="py-2 pr-3 text-gray-700 dark:text-gray-300">{{ formatPercent(m.availability_15d) placeholderplaceholder</td>
            <td class="py-2 pr-3 text-gray-700 dark:text-gray-300">{{ formatPercent(m.availability_30d) placeholderplaceholder</td>
            <td class="py-2 pr-3 text-gray-700 dark:text-gray-300">{{ formatLatency(m.avg_latency_7d_ms) placeholderplaceholder</td>
          </tr>
        </tbody>
      </table>
    </div>

    <template #footer>
      <div class="flex justify-end">
        <button @click="$emit('close')" class="btn btn-secondary">
          {{ t('channelStatus.closeDetail') placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import {
  status as fetchChannelMonitorDetail,
  type UserMonitorDetail,
placeholder from '@/api/channelMonitor'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { useChannelMonitorFormat placeholder from '@/composables/useChannelMonitorFormat'

const props = defineProps<{
  show: boolean
  monitorId: number | null
  title: string
placeholder>()

defineEmits<{
  (e: 'close'): void
placeholder>()

const { t placeholder = useI18n()
const appStore = useAppStore()
const { statusLabel, statusBadgeClass, formatLatency, formatPercent, formatMonitorModel placeholder = useChannelMonitorFormat()

const detail = ref<UserMonitorDetail | null>(null)
const loading = ref(false)

async function load(id: number) {
  detail.value = null
  loading.value = true
  try {
    detail.value = await fetchChannelMonitorDetail(id)
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.detailLoadError')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

watch(
  () => [props.show, props.monitorId] as const,
  ([show, id]) => {
    if (!show) {
      detail.value = null
      return
    placeholder
    if (id != null) void load(id)
  placeholder,
  { immediate: true placeholder,
)
</script>
