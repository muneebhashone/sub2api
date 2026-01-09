<template>
  <BaseDialog :show="show" :title="title" width="full" :close-on-click-outside="true" @close="close">
    <div v-if="loading" class="flex items-center justify-center py-16">
      <div class="flex flex-col items-center gap-3">
        <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-primary-600"></div>
        <div class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ t('admin.ops.errorDetail.loading') placeholderplaceholder</div>
      </div>
    </div>

    <div v-else-if="!detail" class="py-10 text-center text-sm text-gray-500 dark:text-gray-400">
      {{ emptyText placeholderplaceholder
    </div>

    <div v-else class="space-y-6 p-6">
      <!-- Top Summary -->
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
        <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.errorDetail.requestId') placeholderplaceholder</div>
          <div class="mt-1 break-all font-mono text-sm font-medium text-gray-900 dark:text-white">
            {{ detail.request_id || detail.client_request_id || '—' placeholderplaceholder
          </div>
        </div>

        <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.errorDetail.time') placeholderplaceholder</div>
          <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">
            {{ formatDateTime(detail.created_at) placeholderplaceholder
          </div>
        </div>

        <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.errorDetail.phase') placeholderplaceholder</div>
          <div class="mt-1 text-sm font-bold uppercase text-gray-900 dark:text-white">
            {{ detail.phase || '—' placeholderplaceholder
          </div>
          <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
            {{ detail.type || '—' placeholderplaceholder
          </div>
        </div>

        <div class="rounded-xl bg-gray-50 p-4 dark:bg-dark-900">
          <div class="text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.errorDetail.status') placeholderplaceholder</div>
          <div class="mt-1 flex flex-wrap items-center gap-2">
            <span :class="['inline-flex items-center rounded-lg px-2 py-1 text-xs font-black ring-1 ring-inset shadow-sm', statusClass]">
              {{ detail.status_code placeholderplaceholder
            </span>
            <span
              v-if="detail.severity"
              :class="['rounded-md px-2 py-0.5 text-[10px] font-black shadow-sm', severityClass]"
            >
              {{ detail.severity placeholderplaceholder
            </span>
          </div>
        </div>
      </div>

      <!-- Message -->
      <div class="rounded-xl bg-gray-50 p-6 dark:bg-dark-900">
        <h3 class="mb-4 text-sm font-black uppercase tracking-wider text-gray-900 dark:text-white">{{ t('admin.ops.errorDetail.message') placeholderplaceholder</h3>
        <div class="text-sm font-medium text-gray-800 dark:text-gray-200 break-words">
          {{ detail.message || '—' placeholderplaceholder
        </div>
      </div>

      <!-- Basic Info -->
      <div class="rounded-xl bg-gray-50 p-6 dark:bg-dark-900">
        <h3 class="mb-4 text-sm font-black uppercase tracking-wider text-gray-900 dark:text-white">{{ t('admin.ops.errorDetail.basicInfo') placeholderplaceholder</h3>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <div>
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.platform') placeholderplaceholder</div>
            <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ detail.platform || '—' placeholderplaceholder</div>
          </div>
          <div>
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.model') placeholderplaceholder</div>
            <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">{{ detail.model || '—' placeholderplaceholder</div>
          </div>
          <div>
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.latency') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-sm font-bold text-gray-900 dark:text-white">
              {{ detail.latency_ms != null ? `${detail.latency_msplaceholderms` : '—' placeholderplaceholder
            </div>
          </div>
          <div>
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.ttft') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-sm font-bold text-gray-900 dark:text-white">
              {{ detail.time_to_first_token_ms != null ? `${detail.time_to_first_token_msplaceholderms` : '—' placeholderplaceholder
            </div>
          </div>
          <div>
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.businessLimited') placeholderplaceholder</div>
            <div class="mt-1 text-sm font-medium text-gray-900 dark:text-white">
              {{ detail.is_business_limited ? 'true' : 'false' placeholderplaceholder
            </div>
          </div>
          <div>
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.requestPath') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-xs text-gray-700 dark:text-gray-200 break-all">
              {{ detail.request_path || '—' placeholderplaceholder
            </div>
          </div>
        </div>
      </div>

      <!-- Timings (best-effort fields) -->
      <div class="rounded-xl bg-gray-50 p-6 dark:bg-dark-900">
        <h3 class="mb-4 text-sm font-black uppercase tracking-wider text-gray-900 dark:text-white">{{ t('admin.ops.errorDetail.timings') placeholderplaceholder</h3>
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <div class="rounded-lg bg-white p-4 shadow-sm dark:bg-dark-800">
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.auth') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-sm font-bold text-gray-900 dark:text-white">
              {{ detail.auth_latency_ms != null ? `${detail.auth_latency_msplaceholderms` : '—' placeholderplaceholder
            </div>
          </div>
          <div class="rounded-lg bg-white p-4 shadow-sm dark:bg-dark-800">
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.routing') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-sm font-bold text-gray-900 dark:text-white">
              {{ detail.routing_latency_ms != null ? `${detail.routing_latency_msplaceholderms` : '—' placeholderplaceholder
            </div>
          </div>
          <div class="rounded-lg bg-white p-4 shadow-sm dark:bg-dark-800">
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.upstream') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-sm font-bold text-gray-900 dark:text-white">
              {{ detail.upstream_latency_ms != null ? `${detail.upstream_latency_msplaceholderms` : '—' placeholderplaceholder
            </div>
          </div>
          <div class="rounded-lg bg-white p-4 shadow-sm dark:bg-dark-800">
            <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.response') placeholderplaceholder</div>
            <div class="mt-1 font-mono text-sm font-bold text-gray-900 dark:text-white">
              {{ detail.response_latency_ms != null ? `${detail.response_latency_msplaceholderms` : '—' placeholderplaceholder
            </div>
          </div>
        </div>
      </div>

      <!-- Retry -->
      <div class="rounded-xl bg-gray-50 p-6 dark:bg-dark-900">
        <div class="flex flex-col justify-between gap-4 md:flex-row md:items-start">
          <div class="space-y-1">
            <h3 class="text-sm font-black uppercase tracking-wider text-gray-900 dark:text-white">{{ t('admin.ops.errorDetail.retry') placeholderplaceholder</h3>
            <div class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.ops.errorDetail.retryNote1') placeholderplaceholder
            </div>
          </div>
          <div class="flex flex-wrap gap-2">
            <button type="button" class="btn btn-secondary btn-sm" :disabled="retrying" @click="openRetryConfirm('client')">
              {{ t('admin.ops.errorDetail.retryClient') placeholderplaceholder
            </button>
            <button
              type="button"
              class="btn btn-secondary btn-sm"
              :disabled="retrying || !pinnedAccountId"
              @click="openRetryConfirm('upstream')"
              :title="pinnedAccountId ? '' : t('admin.ops.errorDetail.retryUpstreamHint')"
            >
              {{ t('admin.ops.errorDetail.retryUpstream') placeholderplaceholder
            </button>
          </div>
        </div>

        <div class="mt-4 grid grid-cols-1 gap-4 md:grid-cols-3">
          <div class="md:col-span-1">
            <label class="mb-1 block text-xs font-bold uppercase tracking-wider text-gray-400">{{ t('admin.ops.errorDetail.pinnedAccountId') placeholderplaceholder</label>
            <input v-model="pinnedAccountIdInput" type="text" class="input font-mono text-sm" :placeholder="t('admin.ops.errorDetail.pinnedAccountIdHint')" />
            <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('admin.ops.errorDetail.retryNote2') placeholderplaceholder
            </div>
          </div>
          <div class="md:col-span-2">
            <div class="rounded-lg bg-white p-4 shadow-sm dark:bg-dark-800">
              <div class="text-xs font-bold uppercase text-gray-400">{{ t('admin.ops.errorDetail.retryNotes') placeholderplaceholder</div>
              <ul class="mt-2 list-disc space-y-1 pl-5 text-xs text-gray-600 dark:text-gray-300">
                <li>{{ t('admin.ops.errorDetail.retryNote3') placeholderplaceholder</li>
                <li>{{ t('admin.ops.errorDetail.retryNote4') placeholderplaceholder</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Request body -->
      <div class="rounded-xl bg-gray-50 p-6 dark:bg-dark-900">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-black uppercase tracking-wider text-gray-900 dark:text-white">{{ t('admin.ops.errorDetail.requestBody') placeholderplaceholder</h3>
          <div
            v-if="detail.request_body_truncated"
            class="rounded-full bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700 dark:bg-amber-900/30 dark:text-amber-300"
          >
            {{ t('admin.ops.errorDetail.trimmed') placeholderplaceholder
          </div>
        </div>
        <pre
          class="mt-4 max-h-[420px] overflow-auto rounded-xl border border-gray-200 bg-white p-4 text-xs text-gray-800 dark:border-dark-700 dark:bg-dark-800 dark:text-gray-100"
        ><code>{{ prettyJSON(detail.request_body) placeholderplaceholder</code></pre>
      </div>

      <!-- Error body -->
      <div class="rounded-xl bg-gray-50 p-6 dark:bg-dark-900">
        <h3 class="text-sm font-black uppercase tracking-wider text-gray-900 dark:text-white">{{ t('admin.ops.errorDetail.errorBody') placeholderplaceholder</h3>
        <pre
          class="mt-4 max-h-[420px] overflow-auto rounded-xl border border-gray-200 bg-white p-4 text-xs text-gray-800 dark:border-dark-700 dark:bg-dark-800 dark:text-gray-100"
        ><code>{{ prettyJSON(detail.error_body) placeholderplaceholder</code></pre>
      </div>
    </div>
  </BaseDialog>

  <ConfirmDialog
    :show="showRetryConfirm"
    :title="t('admin.ops.errorDetail.confirmRetry')"
    :message="retryConfirmMessage"
    @confirm="runConfirmedRetry"
    @cancel="cancelRetry"
  />
</template>

<script setup lang="ts">
import { computed, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { useAppStore placeholder from '@/stores'
import { opsAPI, type OpsErrorDetail, type OpsRetryMode placeholder from '@/api/admin/ops'
import { formatDateTime placeholder from '@/utils/format'
import { getSeverityClass placeholder from '../utils/opsFormatters'

interface Props {
  show: boolean
  errorId: number | null
placeholder

interface Emits {
  (e: 'update:show', value: boolean): void
placeholder

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const { t placeholder = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const detail = ref<OpsErrorDetail | null>(null)

const retrying = ref(false)
const showRetryConfirm = ref(false)
const pendingRetryMode = ref<OpsRetryMode>('client')

const pinnedAccountIdInput = ref('')
const pinnedAccountId = computed<number | null>(() => {
  const raw = String(pinnedAccountIdInput.value || '').trim()
  if (!raw) return null
  const n = Number.parseInt(raw, 10)
  return Number.isFinite(n) && n > 0 ? n : null
placeholder)

const title = computed(() => {
  if (!props.errorId) return 'Error Detail'
  return `Error #${props.errorIdplaceholder`
placeholder)

const emptyText = computed(() => 'No error selected.')

function close() {
  emit('update:show', false)
placeholder

function prettyJSON(raw?: string): string {
  if (!raw) return t('admin.ops.errorDetail.na')
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  placeholder catch {
    return raw
  placeholder
placeholder

async function fetchDetail(id: number) {
  loading.value = true
  try {
    const d = await opsAPI.getErrorLogDetail(id)
    detail.value = d

    // Default pinned account from error log if present.
    if (d.account_id && d.account_id > 0) {
      pinnedAccountIdInput.value = String(d.account_id)
    placeholder else {
      pinnedAccountIdInput.value = ''
    placeholder
  placeholder catch (err: any) {
    detail.value = null
    appStore.showError(err?.message || 'Failed to load error detail')
  placeholder finally {
    loading.value = false
  placeholder
placeholder

watch(
  () => [props.show, props.errorId] as const,
  ([show, id]) => {
    if (!show) {
      detail.value = null
      return
    placeholder
    if (typeof id === 'number' && id > 0) {
      fetchDetail(id)
    placeholder
  placeholder,
  { immediate: true placeholder
)

function openRetryConfirm(mode: OpsRetryMode) {
  pendingRetryMode.value = mode
  showRetryConfirm.value = true
placeholder

const retryConfirmMessage = computed(() => {
  const mode = pendingRetryMode.value
  if (mode === 'upstream') {
    return t('admin.ops.errorDetail.confirmRetryMessage')
  placeholder
  return t('admin.ops.errorDetail.confirmRetryHint')
placeholder)

const severityClass = computed(() => {
  if (!detail.value?.severity) return 'bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300'
  return getSeverityClass(detail.value.severity)
placeholder)

const statusClass = computed(() => {
  const code = detail.value?.status_code ?? 0
  if (code >= 500) return 'bg-red-50 text-red-700 ring-red-600/20 dark:bg-red-900/30 dark:text-red-400 dark:ring-red-500/30'
  if (code === 429) return 'bg-purple-50 text-purple-700 ring-purple-600/20 dark:bg-purple-900/30 dark:text-purple-400 dark:ring-purple-500/30'
  if (code >= 400) return 'bg-amber-50 text-amber-700 ring-amber-600/20 dark:bg-amber-900/30 dark:text-amber-400 dark:ring-amber-500/30'
  return 'bg-gray-50 text-gray-700 ring-gray-600/20 dark:bg-gray-900/30 dark:text-gray-400 dark:ring-gray-500/30'
placeholder)

async function runConfirmedRetry() {
  if (!props.errorId) return
  const mode = pendingRetryMode.value
  showRetryConfirm.value = false

  retrying.value = true
  try {
    const req =
      mode === 'upstream'
        ? { mode, pinned_account_id: pinnedAccountId.value ?? undefined placeholder
        : { mode placeholder

    const res = await opsAPI.retryErrorRequest(props.errorId, req)
    const summary = res.status === 'succeeded' ? t('admin.ops.errorDetail.retrySuccess') : t('admin.ops.errorDetail.retryFailed')
    appStore.showSuccess(summary)
  placeholder catch (err: any) {
    appStore.showError(err?.message || 'Retry failed')
  placeholder finally {
    retrying.value = false
  placeholder
placeholder

function cancelRetry() {
  showRetryConfirm.value = false
placeholder
</script>
