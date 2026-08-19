<template>
  <BaseDialog
    :show="show"
    :title="t('admin.channelMonitor.runResultTitle')"
    width="normal"
    @close="$emit('close')"
  >
    <div class="space-y-2">
      <div
        v-for="r in results"
        :key="r.model"
        class="flex items-center justify-between rounded-lg border border-gray-200 px-3 py-2 text-sm dark:border-dark-600"
      >
        <div class="flex flex-col">
          <span class="font-medium text-gray-900 dark:text-white">{{ formatMonitorModel(r.model) placeholderplaceholder</span>
          <span v-if="r.message" class="text-xs text-gray-500 dark:text-gray-400">{{ r.message placeholderplaceholder</span>
          <MonitorQuotaView :snapshot="r.quota" class="mt-1" />
        </div>
        <div class="flex items-center gap-2">
          <span
            class="inline-flex items-center rounded-full px-2 py-0.5 text-[11px]"
            :class="statusBadgeClass(r.status)"
          >
            {{ statusLabel(r.status) placeholderplaceholder
          </span>
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatLatency(r.latency_ms) placeholderplaceholder ms</span>
        </div>
      </div>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <button @click="$emit('close')" class="btn btn-primary">
          {{ t('common.close') placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { useI18n placeholder from 'vue-i18n'
import type { CheckResult placeholder from '@/api/admin/channelMonitor'
import BaseDialog from '@/components/common/BaseDialog.vue'
import MonitorQuotaView from '@/components/common/MonitorQuotaView.vue'
import { useChannelMonitorFormat placeholder from '@/composables/useChannelMonitorFormat'

defineProps<{
  show: boolean
  results: CheckResult[]
placeholder>()

defineEmits<{
  (e: 'close'): void
placeholder>()

const { t placeholder = useI18n()
const { statusLabel, statusBadgeClass, formatLatency, formatMonitorModel placeholder = useChannelMonitorFormat()
</script>
