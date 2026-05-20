<template>
  <BaseDialog
    :show="show"
    :title="t('admin.channelMonitor.template.applyPickerTitle', { name: templateName placeholder)"
    @close="$emit('close')"
  >
    <p class="mb-3 text-sm text-gray-600 dark:text-gray-400">
      {{ t('admin.channelMonitor.template.applyPickerHint') placeholderplaceholder
    </p>

    <div v-if="loading" class="py-6 text-center text-sm text-gray-400">
      {{ t('common.loading') placeholderplaceholder
    </div>

    <div v-else-if="monitors.length === 0" class="py-6 text-center text-sm text-gray-400">
      {{ t('admin.channelMonitor.template.applyPickerEmpty') placeholderplaceholder
    </div>

    <div v-else>
      <!-- 全选/全不选 -->
      <div class="mb-2 flex items-center gap-3 text-xs">
        <button
          type="button"
          class="text-primary-600 hover:underline dark:text-primary-400"
          @click="selectAll"
        >
          {{ t('common.selectAll') placeholderplaceholder
        </button>
        <button
          type="button"
          class="text-gray-500 hover:underline dark:text-gray-400"
          @click="selectNone"
        >
          {{ t('admin.channelMonitor.template.selectNone') placeholderplaceholder
        </button>
        <span class="ml-auto text-gray-500 dark:text-gray-400">
          {{ t('admin.channelMonitor.template.selectedCount', {
            n: selectedIds.length,
            total: monitors.length,
          placeholder) placeholderplaceholder
        </span>
      </div>

      <ul class="max-h-80 divide-y divide-gray-100 overflow-y-auto rounded-lg border border-gray-200 dark:divide-dark-700 dark:border-dark-700">
        <li
          v-for="m in monitors"
          :key="m.id"
          class="flex cursor-pointer items-center gap-3 px-3 py-2 hover:bg-gray-50 dark:hover:bg-dark-800"
          @click="toggle(m.id)"
        >
          <input
            type="checkbox"
            :checked="selectedSet.has(m.id)"
            class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
            @click.stop="toggle(m.id)"
          />
          <span class="font-medium text-gray-900 dark:text-white">{{ m.name placeholderplaceholder</span>
          <span class="text-xs text-gray-400">{{ m.provider placeholderplaceholder</span>
          <span v-if="m.provider === 'openai'" class="text-xs text-gray-400">{{ m.api_mode placeholderplaceholder</span>
          <span
            v-if="!m.enabled"
            class="ml-auto rounded bg-gray-100 px-1.5 py-0.5 text-xs text-gray-500 dark:bg-dark-700 dark:text-gray-400"
          >
            {{ t('admin.channelMonitor.onlyDisabled').replace(/^仅|^Only /, '') placeholderplaceholder
          </span>
        </li>
      </ul>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <button class="btn btn-secondary" @click="$emit('close')">
          {{ t('common.cancel') placeholderplaceholder
        </button>
        <button
          class="btn btn-primary"
          :disabled="submitting || selectedIds.length === 0"
          @click="handleApply"
        >
          {{ submitting
            ? t('common.submitting')
            : t('admin.channelMonitor.template.applyPickerConfirm', { n: selectedIds.length placeholder) placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import { adminAPI placeholder from '@/api/admin'
import type { AssociatedMonitorBrief placeholder from '@/api/admin/channelMonitorTemplate'
import BaseDialog from '@/components/common/BaseDialog.vue'

const props = defineProps<{
  show: boolean
  templateId: number | null
  templateName: string
placeholder>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'applied', affected: number): void
placeholder>()

const { t placeholder = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const submitting = ref(false)
const monitors = ref<AssociatedMonitorBrief[]>([])
const selectedIds = ref<number[]>([])

const selectedSet = computed(() => new Set(selectedIds.value))

watch(
  () => [props.show, props.templateId] as const,
  ([show, id]) => {
    if (!show || id == null) return
    void fetchMonitors(id)
  placeholder,
  { immediate: true placeholder,
)

async function fetchMonitors(id: number) {
  loading.value = true
  monitors.value = []
  selectedIds.value = []
  try {
    const { items placeholder = await adminAPI.channelMonitorTemplate.listAssociatedMonitors(id)
    monitors.value = items
    // 默认全选
    selectedIds.value = items.map((m) => m.id)
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

function toggle(id: number) {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) selectedIds.value.splice(idx, 1)
  else selectedIds.value.push(id)
placeholder

function selectAll() {
  selectedIds.value = monitors.value.map((m) => m.id)
placeholder

function selectNone() {
  selectedIds.value = []
placeholder

async function handleApply() {
  if (props.templateId == null || selectedIds.value.length === 0 || submitting.value) return
  submitting.value = true
  try {
    const { affected placeholder = await adminAPI.channelMonitorTemplate.apply(
      props.templateId,
      [...selectedIds.value],
    )
    appStore.showSuccess(t('admin.channelMonitor.template.applySuccess', { n: affected placeholder))
    emit('applied', affected)
    emit('close')
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  placeholder finally {
    submitting.value = false
  placeholder
placeholder
</script>
