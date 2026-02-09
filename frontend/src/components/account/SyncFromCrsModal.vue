<template>
  <BaseDialog
    :show="show"
    :title="t('admin.accounts.syncFromCrsTitle')"
    width="normal"
    close-on-click-outside
    @close="handleClose"
  >
    <!-- Step 1: Input credentials -->
    <form
      v-if="currentStep === 'input'"
      id="sync-from-crs-form"
      class="space-y-4"
      @submit.prevent="handlePreview"
    >
      <div class="text-sm text-gray-600 dark:text-dark-300">
        {{ t('admin.accounts.syncFromCrsDesc') placeholderplaceholder
      </div>
      <div
        class="rounded-lg bg-gray-50 p-3 text-xs text-gray-500 dark:bg-dark-700/60 dark:text-dark-400"
      >
        {{ t('admin.accounts.crsUpdateBehaviorNote') placeholderplaceholder
      </div>
      <div
        class="rounded-lg border border-amber-200 bg-amber-50 p-3 text-xs text-amber-600 dark:border-amber-800 dark:bg-amber-900/20 dark:text-amber-400"
      >
        {{ t('admin.accounts.crsVersionRequirement') placeholderplaceholder
      </div>

      <div class="grid grid-cols-1 gap-4">
        <div>
          <label for="crs-base-url" class="input-label">{{ t('admin.accounts.crsBaseUrl') placeholderplaceholder</label>
          <input
            id="crs-base-url"
            v-model="form.base_url"
            type="text"
            class="input"
            required
            :placeholder="t('admin.accounts.crsBaseUrlPlaceholder')"
          />
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <label for="crs-username" class="input-label">{{ t('admin.accounts.crsUsername') placeholderplaceholder</label>
            <input id="crs-username" v-model="form.username" type="text" class="input" required autocomplete="username" />
          </div>
          <div>
            <label for="crs-password" class="input-label">{{ t('admin.accounts.crsPassword') placeholderplaceholder</label>
            <input
              id="crs-password"
              v-model="form.password"
              type="password"
              class="input"
              required
              autocomplete="current-password"
            />
          </div>
        </div>

        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-dark-300">
          <input
            v-model="form.sync_proxies"
            type="checkbox"
            class="rounded border-gray-300 dark:border-dark-600"
          />
          {{ t('admin.accounts.syncProxies') placeholderplaceholder
        </label>
      </div>
    </form>

    <!-- Step 2: Preview & select -->
    <div v-else-if="currentStep === 'preview' && previewResult" class="space-y-4">
      <!-- Existing accounts (read-only info) -->
      <div
        v-if="previewResult.existing_accounts.length"
        class="rounded-lg bg-gray-50 p-3 dark:bg-dark-700/60"
      >
        <div class="mb-2 text-sm font-medium text-gray-700 dark:text-dark-300">
          {{ t('admin.accounts.crsExistingAccounts') placeholderplaceholder
          <span class="ml-1 text-xs text-gray-400">({{ previewResult.existing_accounts.length placeholderplaceholder)</span>
        </div>
        <div class="max-h-32 overflow-auto text-xs text-gray-500 dark:text-dark-400">
          <div
            v-for="acc in previewResult.existing_accounts"
            :key="acc.crs_account_id"
            class="flex items-center gap-2 py-0.5"
          >
            <span
              class="inline-block rounded bg-blue-100 px-1.5 py-0.5 text-[10px] font-medium text-blue-700 dark:bg-blue-900/30 dark:text-blue-400"
            >{{ acc.platform placeholderplaceholder / {{ acc.type placeholderplaceholder</span>
            <span class="truncate">{{ acc.name placeholderplaceholder</span>
          </div>
        </div>
      </div>

      <!-- New accounts (selectable) -->
      <div v-if="previewResult.new_accounts.length">
        <div class="mb-2 flex items-center justify-between">
          <div class="text-sm font-medium text-gray-900 dark:text-white">
            {{ t('admin.accounts.crsNewAccounts') placeholderplaceholder
            <span class="ml-1 text-xs text-gray-400">({{ previewResult.new_accounts.length placeholderplaceholder)</span>
          </div>
          <div class="flex gap-2">
            <button
              type="button"
              class="text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400"
              @click="selectAll"
            >{{ t('admin.accounts.crsSelectAll') placeholderplaceholder</button>
            <button
              type="button"
              class="text-xs text-gray-500 hover:text-gray-600 dark:text-gray-400"
              @click="selectNone"
            >{{ t('admin.accounts.crsSelectNone') placeholderplaceholder</button>
          </div>
        </div>
        <div
          class="max-h-48 overflow-auto rounded-lg border border-gray-200 p-2 dark:border-dark-600"
        >
          <label
            v-for="acc in previewResult.new_accounts"
            :key="acc.crs_account_id"
            class="flex cursor-pointer items-center gap-2 rounded px-2 py-1.5 hover:bg-gray-50 dark:hover:bg-dark-700/40"
          >
            <input
              type="checkbox"
              :checked="selectedIds.has(acc.crs_account_id)"
              class="rounded border-gray-300 dark:border-dark-600"
              @change="toggleSelect(acc.crs_account_id)"
            />
            <span
              class="inline-block rounded bg-green-100 px-1.5 py-0.5 text-[10px] font-medium text-green-700 dark:bg-green-900/30 dark:text-green-400"
            >{{ acc.platform placeholderplaceholder / {{ acc.type placeholderplaceholder</span>
            <span class="truncate text-sm text-gray-700 dark:text-dark-300">{{ acc.name placeholderplaceholder</span>
          </label>
        </div>
        <div class="mt-1 text-xs text-gray-400">
          {{ t('admin.accounts.crsSelectedCount', { count: selectedIds.size placeholder) placeholderplaceholder
        </div>
      </div>

      <!-- Sync options summary -->
      <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-dark-400">
        <span>{{ t('admin.accounts.syncProxies') placeholderplaceholder:</span>
        <span :class="form.sync_proxies ? 'text-green-600 dark:text-green-400' : 'text-gray-400 dark:text-dark-500'">
          {{ form.sync_proxies ? t('common.yes') : t('common.no') placeholderplaceholder
        </span>
      </div>

      <!-- No new accounts -->
      <div
        v-if="!previewResult.new_accounts.length"
        class="rounded-lg bg-gray-50 p-4 text-center text-sm text-gray-500 dark:bg-dark-700/60 dark:text-dark-400"
      >
        {{ t('admin.accounts.crsNoNewAccounts') placeholderplaceholder
        <span v-if="previewResult.existing_accounts.length">
          {{ t('admin.accounts.crsWillUpdate', { count: previewResult.existing_accounts.length placeholder) placeholderplaceholder
        </span>
      </div>
    </div>

    <!-- Step 3: Result -->
    <div v-else-if="currentStep === 'result' && result" class="space-y-4">
      <div
        class="space-y-2 rounded-xl border border-gray-200 p-4 dark:border-dark-700"
      >
        <div class="text-sm font-medium text-gray-900 dark:text-white">
          {{ t('admin.accounts.syncResult') placeholderplaceholder
        </div>
        <div class="text-sm text-gray-700 dark:text-dark-300">
          {{ t('admin.accounts.syncResultSummary', result) placeholderplaceholder
        </div>

        <div v-if="errorItems.length" class="mt-2">
          <div class="text-sm font-medium text-red-600 dark:text-red-400">
            {{ t('admin.accounts.syncErrors') placeholderplaceholder
          </div>
          <div
            class="mt-2 max-h-48 overflow-auto rounded-lg bg-gray-50 p-3 font-mono text-xs dark:bg-dark-800"
          >
            <div v-for="(item, idx) in errorItems" :key="idx" class="whitespace-pre-wrap">
              {{ item.kind placeholderplaceholder {{ item.crs_account_id placeholderplaceholder — {{ item.action
              placeholderplaceholder{{ item.error ? `: ${item.errorplaceholder` : '' placeholderplaceholder
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-3">
        <!-- Step 1: Input -->
        <template v-if="currentStep === 'input'">
          <button
            class="btn btn-secondary"
            type="button"
            :disabled="previewing"
            @click="handleClose"
          >
            {{ t('common.cancel') placeholderplaceholder
          </button>
          <button
            class="btn btn-primary"
            type="submit"
            form="sync-from-crs-form"
            :disabled="previewing"
          >
            {{ previewing ? t('admin.accounts.crsPreviewing') : t('admin.accounts.crsPreview') placeholderplaceholder
          </button>
        </template>

        <!-- Step 2: Preview -->
        <template v-else-if="currentStep === 'preview'">
          <button
            class="btn btn-secondary"
            type="button"
            :disabled="syncing"
            @click="handleBack"
          >
            {{ t('admin.accounts.crsBack') placeholderplaceholder
          </button>
          <button
            class="btn btn-primary"
            type="button"
            :disabled="syncing || hasNewButNoneSelected"
            @click="handleSync"
          >
            {{ syncing ? t('admin.accounts.syncing') : t('admin.accounts.syncNow') placeholderplaceholder
          </button>
        </template>

        <!-- Step 3: Result -->
        <template v-else-if="currentStep === 'result'">
          <button class="btn btn-secondary" type="button" @click="handleClose">
            {{ t('common.close') placeholderplaceholder
          </button>
        </template>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import type { PreviewFromCRSResult placeholder from '@/api/admin/accounts'

interface Props {
  show: boolean
placeholder

interface Emits {
  (e: 'close'): void
  (e: 'synced'): void
placeholder

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { t placeholder = useI18n()
const appStore = useAppStore()

type Step = 'input' | 'preview' | 'result'
const currentStep = ref<Step>('input')
const previewing = ref(false)
const syncing = ref(false)
const previewResult = ref<PreviewFromCRSResult | null>(null)
const selectedIds = ref(new Set<string>())
const result = ref<Awaited<ReturnType<typeof adminAPI.accounts.syncFromCrs>> | null>(null)

const form = reactive({
  base_url: '',
  username: '',
  password: '',
  sync_proxies: true
placeholder)

const hasNewButNoneSelected = computed(() => {
  if (!previewResult.value) return false
  return previewResult.value.new_accounts.length > 0 && selectedIds.value.size === 0
placeholder)

const errorItems = computed(() => {
  if (!result.value?.items) return []
  return result.value.items.filter(
    (i) => i.action === 'failed' || (i.action === 'skipped' && i.error !== 'not selected')
  )
placeholder)

watch(
  () => props.show,
  (open) => {
    if (open) {
      currentStep.value = 'input'
      previewResult.value = null
      selectedIds.value = new Set()
      result.value = null
      form.base_url = ''
      form.username = ''
      form.password = ''
      form.sync_proxies = true
    placeholder
  placeholder
)

const handleClose = () => {
  if (syncing.value || previewing.value) {
    return
  placeholder
  emit('close')
placeholder

const handleBack = () => {
  currentStep.value = 'input'
  previewResult.value = null
  selectedIds.value = new Set()
placeholder

const selectAll = () => {
  if (!previewResult.value) return
  selectedIds.value = new Set(previewResult.value.new_accounts.map((a) => a.crs_account_id))
placeholder

const selectNone = () => {
  selectedIds.value = new Set()
placeholder

const toggleSelect = (id: string) => {
  const s = new Set(selectedIds.value)
  if (s.has(id)) {
    s.delete(id)
  placeholder else {
    s.add(id)
  placeholder
  selectedIds.value = s
placeholder

const handlePreview = async () => {
  if (!form.base_url.trim() || !form.username.trim() || !form.password.trim()) {
    appStore.showError(t('admin.accounts.syncMissingFields'))
    return
  placeholder

  previewing.value = true
  try {
    const res = await adminAPI.accounts.previewFromCrs({
      base_url: form.base_url.trim(),
      username: form.username.trim(),
      password: form.password
    placeholder)
    previewResult.value = res
    // Auto-select all new accounts
    selectedIds.value = new Set(res.new_accounts.map((a) => a.crs_account_id))
    currentStep.value = 'preview'
  placeholder catch (error: any) {
    appStore.showError(error?.message || t('admin.accounts.crsPreviewFailed'))
  placeholder finally {
    previewing.value = false
  placeholder
placeholder

const handleSync = async () => {
  if (!form.base_url.trim() || !form.username.trim() || !form.password.trim()) {
    appStore.showError(t('admin.accounts.syncMissingFields'))
    return
  placeholder

  syncing.value = true
  try {
    const res = await adminAPI.accounts.syncFromCrs({
      base_url: form.base_url.trim(),
      username: form.username.trim(),
      password: form.password,
      sync_proxies: form.sync_proxies,
      selected_account_ids: [...selectedIds.value]
    placeholder)
    result.value = res
    currentStep.value = 'result'

    if (res.failed > 0) {
      appStore.showError(t('admin.accounts.syncCompletedWithErrors', res))
    placeholder else {
      appStore.showSuccess(t('admin.accounts.syncCompleted', res))
    placeholder
    emit('synced')
  placeholder catch (error: any) {
    appStore.showError(error?.message || t('admin.accounts.syncFailed'))
  placeholder finally {
    syncing.value = false
  placeholder
placeholder
</script>
