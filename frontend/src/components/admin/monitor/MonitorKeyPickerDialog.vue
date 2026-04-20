<template>
  <BaseDialog
    :show="show"
    :title="t('admin.channelMonitor.form.selectKeyTitle')"
    width="normal"
    @close="$emit('close')"
  >
    <div class="space-y-3">
      <p class="text-xs text-gray-500 dark:text-gray-400">
        {{ t('admin.channelMonitor.form.selectKeyHint') placeholderplaceholder
      </p>
      <div v-if="loading" class="py-6 text-center text-sm text-gray-500">
        {{ t('common.loading') placeholderplaceholder
      </div>
      <div v-else-if="keys.length === 0" class="py-6 text-center text-sm text-gray-500">
        {{ t('admin.channelMonitor.form.noActiveKey') placeholderplaceholder
      </div>
      <div v-else class="max-h-72 space-y-1 overflow-auto">
        <button
          v-for="k in keys"
          :key="k.id"
          type="button"
          @click="$emit('pick', k)"
          class="block w-full rounded-lg border border-gray-200 px-3 py-2 text-left text-sm transition-colors hover:bg-gray-50 dark:border-dark-600 dark:hover:bg-dark-700"
        >
          <div class="font-medium text-gray-900 dark:text-white">{{ k.name placeholderplaceholder</div>
          <div class="font-mono text-xs text-gray-500 dark:text-gray-400">{{ maskKey(k.key) placeholderplaceholder</div>
        </button>
      </div>
    </div>
    <template #footer>
      <div class="flex justify-end">
        <button @click="$emit('close')" class="btn btn-secondary">
          {{ t('common.cancel') placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { useI18n placeholder from 'vue-i18n'
import type { ApiKey placeholder from '@/types'
import BaseDialog from '@/components/common/BaseDialog.vue'

defineProps<{
  show: boolean
  loading: boolean
  keys: ApiKey[]
placeholder>()

defineEmits<{
  (e: 'close'): void
  (e: 'pick', key: ApiKey): void
placeholder>()

const { t placeholder = useI18n()

function maskKey(key: string): string {
  if (!key) return ''
  if (key.length <= 12) return `${key.slice(0, 4)placeholder***`
  return `${key.slice(0, 6)placeholder...${key.slice(-4)placeholder`
placeholder
</script>
