<template>
  <section v-if="state?.eligible" class="space-y-4 border-t border-gray-200 pt-4 dark:border-dark-600" data-testid="ollama-cloud-usage-settings">
    <div class="flex items-start justify-between gap-4">
      <div>
        <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ t('admin.accounts.ollamaCloud.title') placeholderplaceholder
        </h3>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {{ t('admin.accounts.ollamaCloud.sessionSecurityHint') placeholderplaceholder
        </p>
      </div>
      <span
        class="whitespace-nowrap rounded px-2 py-1 text-xs font-medium"
        :class="state.configured
          ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
          : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'"
      >
        {{ state.configured ? t('admin.accounts.ollamaCloud.configured') : t('admin.accounts.ollamaCloud.notConfigured') placeholderplaceholder
      </span>
    </div>

    <div v-if="loading" class="flex h-20 items-center justify-center text-gray-400">
      <Icon name="refresh" size="sm" class="animate-spin" />
    </div>
    <template v-else>
      <div v-if="!state.encryption_key_configured" class="rounded border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800 dark:border-amber-800/50 dark:bg-amber-900/20 dark:text-amber-200">
        {{ t('admin.accounts.ollamaCloud.encryptionKeyRequired') placeholderplaceholder
      </div>

      <div
        v-if="snapshot"
        class="border-y border-gray-100 py-3 dark:border-dark-700"
        data-testid="ollama-cloud-usage-details"
      >
        <div class="grid grid-cols-[minmax(4rem,auto)_minmax(0,1fr)] gap-x-3 gap-y-1.5 text-xs">
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.plan') placeholderplaceholder</span>
          <span class="break-words text-gray-900 dark:text-white">{{ snapshot.data?.plan || '-' placeholderplaceholder</span>
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.fiveHour') placeholderplaceholder</span>
          <span class="break-words text-gray-900 dark:text-white">{{ windowSummary(snapshot.data?.five_hour) placeholderplaceholder</span>
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.sevenDay') placeholderplaceholder</span>
          <span class="break-words text-gray-900 dark:text-white">{{ windowSummary(snapshot.data?.seven_day) placeholderplaceholder</span>
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.balance') placeholderplaceholder</span>
          <span class="break-words text-gray-900 dark:text-white">{{ snapshot.data?.balance || '-' placeholderplaceholder</span>
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.models') placeholderplaceholder</span>
          <span class="break-words text-gray-900 dark:text-white">{{ modelSummary placeholderplaceholder</span>
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.status') placeholderplaceholder</span>
          <span class="break-words font-medium text-gray-900 dark:text-white">{{ statusLabel placeholderplaceholder</span>
          <span class="text-gray-500 dark:text-gray-400">{{ t('admin.accounts.ollamaCloud.updatedAt') placeholderplaceholder</span>
          <span class="break-words text-gray-900 dark:text-white">{{ formatDate(snapshot.fetched_at || snapshot.last_attempt_at) placeholderplaceholder</span>
        </div>
        <p v-if="snapshot.last_error" class="mt-2 break-words border-t border-gray-100 pt-2 text-xs text-amber-700 dark:border-dark-700 dark:text-amber-300">
          {{ t(`admin.accounts.ollamaCloud.errors.${snapshot.last_errorplaceholder`, snapshot.last_error) placeholderplaceholder
        </p>
      </div>

      <div>
        <label class="input-label" for="ollama-cloud-session">{{ t('admin.accounts.ollamaCloud.sessionLabel') placeholderplaceholder</label>
        <textarea
          id="ollama-cloud-session"
          v-model="session"
          rows="3"
          class="input font-mono text-xs"
          autocomplete="new-password"
          data-1p-ignore
          data-lpignore="true"
          data-bwignore="true"
          :placeholder="t('admin.accounts.ollamaCloud.sessionPlaceholder')"
        />
        <p class="input-hint">{{ t('admin.accounts.ollamaCloud.writeOnlyHint') placeholderplaceholder</p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <button
          type="button"
          class="btn btn-primary btn-sm"
          :disabled="saving || !session.trim() || !state.encryption_key_configured"
          data-testid="ollama-cloud-session-save"
          @click="saveSession"
        >
          <Icon name="check" size="xs" class="mr-1.5" />
          {{ t('common.save') placeholderplaceholder
        </button>
        <button
          v-if="state.configured"
          type="button"
          class="btn btn-secondary btn-sm text-red-600 dark:text-red-400"
          :disabled="saving"
          data-testid="ollama-cloud-session-delete"
          @click="showDeleteConfirm = true"
        >
          <Icon name="trash" size="xs" class="mr-1.5" />
          {{ t('admin.accounts.ollamaCloud.deleteSession') placeholderplaceholder
        </button>
        <button
          v-if="state.configured"
          type="button"
          class="btn btn-secondary btn-sm"
          :disabled="refreshing"
          data-testid="ollama-cloud-refresh"
          @click="refreshUsage"
        >
          <Icon name="refresh" size="xs" class="mr-1.5" :class="{ 'animate-spin': refreshing placeholder" />
          {{ t('admin.accounts.ollamaCloud.refreshNow') placeholderplaceholder
        </button>
      </div>

      <div v-if="state.configured" class="flex items-center justify-between gap-4 border-t border-gray-100 pt-4 dark:border-dark-700">
        <div>
          <label class="text-sm font-medium text-gray-900 dark:text-white">
            {{ t('admin.accounts.ollamaCloud.autoRefresh') placeholderplaceholder
          </label>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            {{ t('admin.accounts.ollamaCloud.autoRefreshHint') placeholderplaceholder
          </p>
        </div>
        <Toggle
          :model-value="state.auto_refresh_enabled"
          :disabled="saving"
          data-testid="ollama-cloud-auto-refresh"
          @update:model-value="setAutoRefresh"
        />
      </div>
    </template>

    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('admin.accounts.ollamaCloud.deleteSession')"
      :message="t('admin.accounts.ollamaCloud.deleteConfirm')"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="deleteSession"
      @cancel="showDeleteConfirm = false"
    />
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { adminAPI placeholder from '@/api/admin'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage, extractI18nErrorMessage placeholder from '@/utils/apiError'
import type { Account, OllamaCloudUsageState, OllamaCloudUsageWindow placeholder from '@/types'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Toggle from '@/components/common/Toggle.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{ account: Account placeholder>()
const emit = defineEmits<{ updated: [state: OllamaCloudUsageState] placeholder>()
const { t placeholder = useI18n()
const appStore = useAppStore()
const state = ref<OllamaCloudUsageState | null>(props.account.ollama_cloud_usage ?? null)
const session = ref('')
const loading = ref(false)
const saving = ref(false)
const refreshing = ref(false)
const showDeleteConfirm = ref(false)
const snapshot = computed(() => state.value?.snapshot)
const statusLabel = computed(() => {
  if (!snapshot.value) return t('admin.accounts.ollamaCloud.notRefreshed')
  if (snapshot.value.status === 'unauthorized') return t('admin.accounts.ollamaCloud.unauthorized')
  if (snapshot.value.status === 'failed') return t('admin.accounts.ollamaCloud.failed')
  return t('admin.accounts.ollamaCloud.ok')
placeholder)
const modelSummary = computed(() => snapshot.value?.data?.models?.map(model => {
  const window = model.window === 'five_hour'
    ? t('admin.accounts.ollamaCloud.fiveHourShort')
    : t('admin.accounts.ollamaCloud.sevenDayShort')
  return `${windowplaceholder ${model.modelplaceholder: ${model.requestsplaceholder`
placeholder).join(', ') || '-')

const formatPercent = (value?: number) => typeof value === 'number' && Number.isFinite(value)
  ? `${value.toFixed(value % 1 ? 1 : 0)placeholder%`
  : '-'
const formatDate = (value?: string) => {
  if (!value) return '-'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : date.toLocaleString()
placeholder
const windowSummary = (window?: OllamaCloudUsageWindow) => {
  if (!window) return '-'
  const reset = window.reset_at ? formatDate(window.reset_at) : window.reset_text
  return reset
    ? t('admin.accounts.ollamaCloud.windowWithReset', { percent: formatPercent(window.used_percent), reset placeholder)
    : formatPercent(window.used_percent)
placeholder

const applyState = (next: OllamaCloudUsageState) => {
  state.value = next
  emit('updated', next)
placeholder

const load = async () => {
  loading.value = true
  try {
    applyState(await adminAPI.accounts.getOllamaCloudUsage(props.account.id))
  placeholder catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.accounts.ollamaCloud.loadFailed')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

const saveSession = async () => {
  if (!session.value.trim()) return
  saving.value = true
  try {
    applyState(await adminAPI.accounts.saveOllamaCloudUsageSession(props.account.id, session.value))
    session.value = ''
    appStore.showSuccess(t('admin.accounts.ollamaCloud.sessionSaved'))
  placeholder catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.accounts.ollamaCloud.sessionSaveFailed')))
  placeholder finally {
    saving.value = false
  placeholder
placeholder

const deleteSession = async () => {
  saving.value = true
  showDeleteConfirm.value = false
  try {
    applyState(await adminAPI.accounts.deleteOllamaCloudUsageSession(props.account.id))
    session.value = ''
    appStore.showSuccess(t('admin.accounts.ollamaCloud.sessionDeleted'))
  placeholder catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.accounts.ollamaCloud.sessionDeleteFailed')))
  placeholder finally {
    saving.value = false
  placeholder
placeholder

const setAutoRefresh = async (enabled: boolean) => {
  saving.value = true
  try {
    applyState(await adminAPI.accounts.setOllamaCloudUsageAutoRefresh(props.account.id, enabled))
  placeholder catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('admin.accounts.ollamaCloud.autoRefreshFailed')))
  placeholder finally {
    saving.value = false
  placeholder
placeholder

const refreshUsage = async () => {
  refreshing.value = true
  try {
    applyState(await adminAPI.accounts.refreshOllamaCloudUsage(props.account.id))
    appStore.showSuccess(t('admin.accounts.ollamaCloud.refreshSuccess'))
  placeholder catch (error) {
    appStore.showError(extractI18nErrorMessage(
      error,
      t,
      'admin.accounts.ollamaCloud.errors',
      t('admin.accounts.ollamaCloud.refreshFailed')
    ))
  placeholder finally {
    refreshing.value = false
  placeholder
placeholder

watch(() => props.account.id, () => {
  state.value = props.account.ollama_cloud_usage ?? null
  session.value = ''
  if (!state.value) void load()
placeholder)

onMounted(() => {
  if (!state.value) void load()
placeholder)
</script>
