<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.dataImportTitle')"
    width="normal"
    close-on-click-outside
    @close="handleClose"
  >
    <form id="import-data-form" class="space-y-4" @submit.prevent="handleImport">
      <div class="text-sm text-gray-600 dark:text-dark-300">
        {{ t('admin.accounts.dataImportHint') placeholderplaceholder
      </div>
      <div
        class="rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-600 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-400"
      >
        {{ t('admin.accounts.dataImportWarning') placeholderplaceholder
      </div>

      <div>
        <label class="input-label">{{ t('admin.accounts.dataImportFile') placeholderplaceholder</label>
        <div
          class="flex items-center justify-between gap-3 rounded-lg border border-dashed px-4 py-3 transition-colors"
          :class="dragActive
            ? 'border-primary-400 bg-primary-50/70 dark:border-primary-500 dark:bg-primary-900/20'
            : 'border-gray-300 bg-gray-50 dark:border-dark-600 dark:bg-dark-800'"
          @dragenter.prevent="handleDragEnter"
          @dragover.prevent="handleDragOver"
          @dragleave.prevent="handleDragLeave"
          @drop.prevent="handleDrop"
        >
          <div class="min-w-0">
            <div class="truncate text-sm text-gray-700 dark:text-dark-200" :title="fileListTitle">
              {{ selectedFilesLabel || t('admin.accounts.dataImportSelectFile') placeholderplaceholder
            </div>
            <div class="text-xs text-gray-500 dark:text-dark-400">
              JSON (.json)
              <span v-if="files.length > 1"> · {{ fileListTitle placeholderplaceholder</span>
            </div>
          </div>
          <button type="button" class="btn btn-secondary shrink-0" @click="openFilePicker">
            {{ t('common.chooseFile') placeholderplaceholder
          </button>
        </div>
        <input
          ref="fileInput"
          type="file"
          class="hidden"
          accept="application/json,.json"
          multiple
          @change="handleFileChange"
        />
      </div>

      <div
        v-if="result"
        class="space-y-2 rounded-xl border border-gray-200 p-4 dark:border-dark-700"
      >
        <div class="text-sm font-medium text-gray-900 dark:text-white">
          {{ t('admin.accounts.dataImportResult') placeholderplaceholder
        </div>
        <div class="text-sm text-gray-700 dark:text-dark-300">
          {{ t('admin.accounts.dataImportResultSummary', result) placeholderplaceholder
        </div>

        <div v-if="errorItems.length" class="mt-2">
          <div class="text-sm font-medium text-red-600 dark:text-red-400">
            {{ t('admin.accounts.dataImportErrors') placeholderplaceholder
          </div>
          <div
            class="mt-2 max-h-48 overflow-auto rounded-lg bg-gray-50 p-3 font-mono text-xs dark:bg-dark-800"
          >
            <div v-for="(item, idx) in errorItems" :key="idx" class="whitespace-pre-wrap">
              {{ item.kind placeholderplaceholder {{ item.name || item.proxy_key || '-' placeholderplaceholder — {{ item.message placeholderplaceholder
            </div>
          </div>
        </div>
      </div>
    </form>

    <template #footer>
      <div class="flex justify-end gap-3">
        <button class="btn btn-secondary" type="button" :disabled="importing" @click="handleClose">
          {{ t('common.cancel') placeholderplaceholder
        </button>
        <button
          class="btn btn-primary"
          type="submit"
          form="import-data-form"
          :disabled="importing"
        >
          {{ importing ? t('admin.accounts.dataImporting') : t('admin.accounts.dataImportButton') placeholderplaceholder
        </button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { adminAPI placeholder from '@/api/admin'
import { useAppStore placeholder from '@/stores/app'
import type { AdminDataImportResult placeholder from '@/types'

interface Props {
  show: boolean
placeholder

interface Emits {
  (e: 'close'): void
  (e: 'imported'): void
placeholder

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { t placeholder = useI18n()
const appStore = useAppStore()

const importing = ref(false)
const files = ref<File[]>([])
const dragActive = ref(false)
const dragDepth = ref(0)
const result = ref<AdminDataImportResult | null>(null)

const fileInput = ref<HTMLInputElement | null>(null)
const selectedFilesLabel = computed(() => {
  if (files.value.length === 0) return ''
  if (files.value.length === 1) return files.value[0]?.name || ''
  return t('admin.accounts.selectedCount', { count: files.value.length placeholder)
placeholder)
const fileListTitle = computed(() => files.value.map((item) => item.name).join(', '))

const errorItems = computed(() => result.value?.errors || [])

watch(
  () => props.show,
  (open) => {
    if (open) {
      files.value = []
      dragActive.value = false
      dragDepth.value = 0
      result.value = null
      if (fileInput.value) {
        fileInput.value.value = ''
      placeholder
    placeholder
  placeholder
)

const openFilePicker = () => {
  fileInput.value?.click()
placeholder

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  setSelectedFiles(target.files)
placeholder

const handleClose = () => {
  if (importing.value) return
  emit('close')
placeholder

const isJsonFile = (sourceFile: File) => {
  const name = sourceFile.name.toLowerCase()
  return name.endsWith('.json') || sourceFile.type === 'application/json'
placeholder

const setSelectedFiles = (sourceFiles: FileList | File[] | null | undefined) => {
  if (importing.value) return
  const picked = Array.from(sourceFiles || []).filter(isJsonFile)
  if (!picked.length) {
    files.value = []
    appStore.showError(t('admin.accounts.dataImportSelectFile'))
    return
  placeholder
  files.value = picked
  result.value = null
placeholder

const handleDragEnter = () => {
  if (importing.value) return
  dragDepth.value += 1
  dragActive.value = true
placeholder

const handleDragOver = () => {
  if (importing.value) return
  dragActive.value = true
placeholder

const handleDragLeave = () => {
  if (importing.value) return
  dragDepth.value = Math.max(0, dragDepth.value - 1)
  if (dragDepth.value === 0) {
    dragActive.value = false
  placeholder
placeholder

const handleDrop = (event: DragEvent) => {
  if (importing.value) return
  dragDepth.value = 0
  dragActive.value = false
  setSelectedFiles(event.dataTransfer?.files)
placeholder

const readFileAsText = async (sourceFile: File): Promise<string> => {
  if (typeof sourceFile.text === 'function') {
    return sourceFile.text()
  placeholder

  if (typeof sourceFile.arrayBuffer === 'function') {
    const buffer = await sourceFile.arrayBuffer()
    return new TextDecoder().decode(buffer)
  placeholder

  return await new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result ?? ''))
    reader.onerror = () => reject(reader.error || new Error('Failed to read file'))
    reader.readAsText(sourceFile)
  placeholder)
placeholder

const mergeDataPayloads = (payloads: any[]) => {
  if (payloads.length === 1) return payloads[0]

  return {
    type: payloads.find((item) => typeof item?.type === 'string')?.type,
    version: payloads.find((item) => typeof item?.version === 'number')?.version,
    exported_at: new Date().toISOString(),
    proxies: payloads.flatMap((item) => Array.isArray(item?.proxies) ? item.proxies : []),
    accounts: payloads.flatMap((item) => Array.isArray(item?.accounts) ? item.accounts : []),
    skipped_shadows: payloads.reduce((sum, item) => {
      const count = Number(item?.skipped_shadows || 0)
      return Number.isFinite(count) ? sum + count : sum
    placeholder, 0)
  placeholder
placeholder

const handleImport = async () => {
  if (files.value.length === 0) {
    appStore.showError(t('admin.accounts.dataImportSelectFile'))
    return
  placeholder

  importing.value = true
  try {
    const dataPayloads = []
    for (const sourceFile of files.value) {
      const text = await readFileAsText(sourceFile)
      dataPayloads.push(JSON.parse(text))
    placeholder
    const dataPayload = mergeDataPayloads(dataPayloads)

    const res = await adminAPI.accounts.importData({
      data: dataPayload,
      skip_default_group_bind: true
    placeholder)

    result.value = res

    const msgParams: Record<string, unknown> = {
      account_created: res.account_created,
      account_failed: res.account_failed,
      proxy_created: res.proxy_created,
      proxy_reused: res.proxy_reused,
      proxy_failed: res.proxy_failed,
    placeholder
    if (res.account_failed > 0 || res.proxy_failed > 0) {
      appStore.showError(t('admin.accounts.dataImportCompletedWithErrors', msgParams))
    placeholder else {
      appStore.showSuccess(t('admin.accounts.dataImportSuccess', msgParams))
      emit('imported')
    placeholder
  placeholder catch (error: any) {
    if (error instanceof SyntaxError) {
      appStore.showError(t('admin.accounts.dataImportParseFailed'))
    placeholder else {
      appStore.showError(error?.message || t('admin.accounts.dataImportFailed'))
    placeholder
  placeholder finally {
    importing.value = false
  placeholder
placeholder
</script>
