<template>
  <AppLayout>
    <div class="mx-auto max-w-[1600px] pb-28">
      <header class="mb-6 flex flex-wrap items-end justify-between gap-4">
        <div>
          <p class="text-xs font-semibold uppercase tracking-[0.16em] text-primary-600 dark:text-primary-400">{{ t('nav.securityAudit') placeholderplaceholder</p>
          <h1 class="mt-1 text-2xl font-semibold tracking-tight text-gray-950 dark:text-white">{{ t('admin.promptAudit.title') placeholderplaceholder</h1>
          <p class="mt-2 max-w-3xl text-sm text-gray-500 dark:text-dark-300">{{ t('admin.promptAudit.description') placeholderplaceholder</p>
        </div>
        <div v-if="draft" class="text-right text-xs text-gray-500 dark:text-dark-400">
          <p>{{ t('admin.promptAudit.configVersion', { version: draft.config_version placeholder) placeholderplaceholder</p>
          <p v-if="draft.updated_at" class="mt-1">{{ formatDate(draft.updated_at) placeholderplaceholder</p>
        </div>
      </header>

      <div v-if="loadErrors.config && !draft" role="alert" class="rounded-xl border border-red-200 bg-red-50 p-5 dark:border-red-900 dark:bg-red-950/30">
        <p class="text-sm text-red-700 dark:text-red-300">{{ loadErrors.config placeholderplaceholder</p>
        <button type="button" class="btn btn-secondary btn-sm mt-3" @click="loadConfig">{{ t('admin.promptAudit.actions.retry') placeholderplaceholder</button>
      </div>

      <main v-else class="rounded-2xl border border-gray-200 bg-white px-4 shadow-sm dark:border-dark-700 dark:bg-dark-850 sm:px-6 lg:px-8">
        <RuntimeOverview :runtime="runtime" :loading="loading.runtime" :error="loadErrors.runtime" @refresh="loadRuntime" />

        <template v-if="draft">
          <EndpointPool
            :endpoints="draft.endpoints"
            :probe-results="probeResults"
            :probing-ids="probingIds"
            @update:endpoints="updateEndpoints"
            @probe="runProbe"
          />
          <div v-if="loadErrors.groups" role="alert" class="mt-5 rounded-lg bg-amber-50 px-4 py-3 text-sm text-amber-800 dark:bg-amber-950/30 dark:text-amber-200">{{ loadErrors.groups placeholderplaceholder</div>
          <PolicyPanel :draft="draft" :groups="groups" @update:draft="replaceDraft" />
        </template>

        <EventWorkspace
          :events="events.items"
          :total="events.total"
          :page="events.page"
          :page-size="events.page_size"
          :filters="filters"
          :selected-ids="selectedEventIds"
          :loading="loading.events"
          :error="loadErrors.events"
          @filters-change="handleFiltersChanged"
          @search="applyEventFilters"
          @selection="selectedEventIds = $event"
          @page="changePage"
          @page-size="changePageSize"
          @view="openEvent"
          @delete="requestSingleDelete"
          @batch-delete="requestBatchDelete"
          @preview-delete="requestFilterDeletePreview"
        />
      </main>
    </div>

    <div v-if="draft" class="fixed inset-x-0 bottom-0 z-30 border-t border-gray-200 bg-white/95 px-4 py-3 shadow-[0_-12px_35px_rgba(15,23,42,0.08)] backdrop-blur dark:border-dark-700 dark:bg-dark-900/95 lg:left-64">
      <div class="mx-auto flex max-w-[1600px] flex-wrap items-center justify-between gap-3">
        <div class="flex flex-wrap items-center gap-x-5 gap-y-2">
          <SaveToggle :label="t('admin.promptAudit.saveBar.enabled')" :model-value="draft.enabled" data-test="enabled-toggle" @update:model-value="setEnabled" />
          <SaveToggle :label="t('admin.promptAudit.saveBar.blocking')" :model-value="draft.blocking_enabled" :disabled="!draft.enabled" data-test="blocking-toggle" @update:model-value="setBlocking" />
          <SaveToggle :label="t('admin.promptAudit.saveBar.storePass')" :model-value="draft.store_pass_events" data-test="store-pass-toggle" @update:model-value="replaceDraft({ ...draft!, store_pass_events: $event placeholder)" />
        </div>
        <div class="flex items-center gap-3">
          <span class="text-sm" :class="dirty ? 'text-amber-700 dark:text-amber-300' : 'text-gray-500 dark:text-dark-400'">
            {{ dirty ? t('admin.promptAudit.saveBar.dirty') : t('admin.promptAudit.saveBar.synced') placeholderplaceholder
          </span>
          <button type="button" class="btn btn-secondary" :disabled="!dirty || loading.saving" @click="resetDraft">{{ t('common.reset') placeholderplaceholder</button>
          <button type="button" class="btn btn-primary" :disabled="!dirty || loading.saving" data-test="save-config" @click="saveConfig">
            {{ loading.saving ? t('common.saving') : t('common.save') placeholderplaceholder
          </button>
        </div>
      </div>
    </div>

    <ConfirmDialog
      :show="showBlockingConfirmation"
      :title="t('admin.promptAudit.blockingConfirm.title')"
      :message="t('admin.promptAudit.blockingConfirm.message')"
      :confirm-text="t('admin.promptAudit.blockingConfirm.confirm')"
      danger
      @confirm="confirmBlocking"
      @cancel="showBlockingConfirmation = false"
    />
    <ConfirmDialog
      :show="deleteRequest.mode !== ''"
      :title="t('admin.promptAudit.events.deleteConfirmTitle')"
      :message="t('admin.promptAudit.events.deleteConfirmMessage', { count: deleteRequest.ids.length placeholder)"
      :confirm-text="t('common.delete')"
      danger
      @confirm="confirmIDDelete"
      @cancel="clearDeleteRequest"
    />
    <BaseDialog :show="Boolean(deletePreview)" :title="t('admin.promptAudit.events.filterDeleteTitle')" width="normal" @close="deletePreview = null">
      <div v-if="deletePreview" class="space-y-4 text-sm text-gray-600 dark:text-dark-300">
        <p>{{ t('admin.promptAudit.events.filterDeleteCount', { count: deletePreview.matched_count placeholder) placeholderplaceholder</p>
        <dl class="grid grid-cols-[auto_1fr] gap-x-3 gap-y-2">
          <dt>{{ t('admin.promptAudit.events.snapshotMax') placeholderplaceholder</dt><dd>{{ deletePreview.snapshot_max_id placeholderplaceholder</dd>
          <dt>Filter SHA-256</dt><dd class="break-all font-mono text-xs">{{ deletePreview.filter_hash placeholderplaceholder</dd>
          <dt>{{ t('admin.promptAudit.events.expiresAt') placeholderplaceholder</dt><dd>{{ formatDate(deletePreview.expires_at) placeholderplaceholder</dd>
        </dl>
        <p class="rounded-lg bg-amber-50 px-3 py-2 text-amber-800 dark:bg-amber-950/30 dark:text-amber-200">{{ t('admin.promptAudit.events.filterDeleteWarning') placeholderplaceholder</p>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button type="button" class="btn btn-secondary" @click="deletePreview = null">{{ t('common.cancel') placeholderplaceholder</button>
          <button type="button" class="btn btn-danger" :disabled="loading.deleting" data-test="confirm-filter-delete" @click="confirmFilterDelete">{{ t('admin.promptAudit.events.confirmFilterDelete') placeholderplaceholder</button>
        </div>
      </template>
    </BaseDialog>
    <EventDetailDialog :show="showEventDetail" :event="activeEvent" :loading="loading.detail" @close="closeEventDetail" />
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, defineComponent, h, onMounted, reactive, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorCode, extractApiErrorMessage placeholder from '@/utils/apiError'
import RuntimeOverview from './components/RuntimeOverview.vue'
import EndpointPool from './components/EndpointPool.vue'
import PolicyPanel from './components/PolicyPanel.vue'
import EventWorkspace from './components/EventWorkspace.vue'
import EventDetailDialog from './components/EventDetailDialog.vue'
import promptAuditAPI from './api'
import type {
  PromptAuditDraft,
  PromptAuditEndpointDraft,
  PromptAuditEvent,
  PromptAuditGroup,
  PromptAuditRuntime,
  PromptDeletePreview,
  PromptEventFilters,
  PromptEventPage,
  PromptLoadErrors,
  PromptProbeResult,
placeholder from './types'
import { buildUpdateRequest, cloneData, configToDraft, draftFingerprint, emptyEventFilters placeholder from './viewModel'

const { t, locale placeholder = useI18n()
const appStore = useAppStore()
const serverConfig = ref<PromptAuditDraft | null>(null)
const draft = ref<PromptAuditDraft | null>(null)
const runtime = ref<PromptAuditRuntime | null>(null)
const groups = ref<PromptAuditGroup[]>([])
const events = reactive<PromptEventPage>({ items: [], total: 0, page: 1, page_size: 20, pages: 0 placeholder)
const filters = ref<PromptEventFilters>(emptyEventFilters())
const appliedFilters = ref<PromptEventFilters>(emptyEventFilters())
const selectedEventIds = ref<number[]>([])
const activeEvent = ref<PromptAuditEvent | null>(null)
const showEventDetail = ref(false)
const probeResults = reactive<Record<string, PromptProbeResult>>({placeholder)
const probingIds = ref<string[]>([])
const deletePreview = ref<PromptDeletePreview | null>(null)
const showBlockingConfirmation = ref(false)
const deleteRequest = reactive<{ mode: '' | 'single' | 'batch'; ids: number[] placeholder>({ mode: '', ids: [] placeholder)
const loading = reactive({ config: false, runtime: false, groups: false, events: false, saving: false, detail: false, deleting: false placeholder)
const loadErrors = reactive<PromptLoadErrors>({ config: '', runtime: '', groups: '', events: '' placeholder)
const dirty = computed(() => draftFingerprint(draft.value) !== draftFingerprint(serverConfig.value))

const SaveToggle = defineComponent({
  inheritAttrs: false,
  props: { label: { type: String, required: true placeholder, modelValue: { type: Boolean, required: true placeholder, disabled: { type: Boolean, default: false placeholder placeholder,
  emits: ['update:modelValue'],
  setup(props, { emit, attrs placeholder) {
    return () => h('label', { class: ['flex items-center gap-2 text-sm', props.disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'] placeholder, [
      h('button', {
        ...attrs, type: 'button', role: 'switch', 'aria-checked': props.modelValue, 'aria-label': props.label, disabled: props.disabled,
        class: ['relative h-6 w-11 rounded-full transition-colors', props.modelValue ? 'bg-primary-600' : 'bg-gray-300 dark:bg-dark-600'],
        onClick: () => !props.disabled && emit('update:modelValue', !props.modelValue),
      placeholder, [h('span', { class: ['absolute top-0.5 h-5 w-5 rounded-full bg-white shadow transition-transform', props.modelValue ? 'translate-x-5' : 'translate-x-0.5'] placeholder)]),
      h('span', { class: 'text-gray-700 dark:text-dark-200' placeholder, props.label),
    ])
  placeholder,
placeholder)

function errorMessage(error: unknown, fallbackKey: string): string {
  const code = extractApiErrorCode(error)
  if (code) {
    const key = `admin.promptAudit.errors.${codeplaceholder`
    const translated = t(key)
    if (translated !== key) return translated
  placeholder
  return extractApiErrorMessage(error, t(fallbackKey))
placeholder

async function loadConfig() {
  loading.config = true
  loadErrors.config = ''
  try {
    const config = await promptAuditAPI.getConfig()
    serverConfig.value = configToDraft(config)
    draft.value = configToDraft(config)
  placeholder catch (error) {
    loadErrors.config = errorMessage(error, 'admin.promptAudit.errors.loadConfig')
  placeholder finally {
    loading.config = false
  placeholder
placeholder
async function loadRuntime() {
  loading.runtime = true
  loadErrors.runtime = ''
  try { runtime.value = await promptAuditAPI.getRuntime() placeholder
  catch (error) { loadErrors.runtime = errorMessage(error, 'admin.promptAudit.errors.loadRuntime') placeholder
  finally { loading.runtime = false placeholder
placeholder
async function loadGroups() {
  loading.groups = true
  loadErrors.groups = ''
  try { groups.value = await promptAuditAPI.listGroups() placeholder
  catch (error) { loadErrors.groups = errorMessage(error, 'admin.promptAudit.errors.loadGroups') placeholder
  finally { loading.groups = false placeholder
placeholder
async function loadEvents() {
  loading.events = true
  loadErrors.events = ''
  try {
    const result = await promptAuditAPI.listEvents(appliedFilters.value, events.page, events.page_size)
    Object.assign(events, result)
    selectedEventIds.value = []
  placeholder catch (error) {
    loadErrors.events = errorMessage(error, 'admin.promptAudit.errors.loadEvents')
  placeholder finally {
    loading.events = false
  placeholder
placeholder
async function loadInitial() {
  await Promise.allSettled([loadConfig(), loadRuntime(), loadGroups(), loadEvents()])
placeholder

function replaceDraft(value: PromptAuditDraft) { draft.value = cloneData(value) placeholder
function updateEndpoints(value: PromptAuditEndpointDraft[]) {
  if (!draft.value) return
  replaceDraft({ ...draft.value, endpoints: value placeholder)
placeholder
function setEnabled(value: boolean) {
  if (!draft.value) return
  replaceDraft({ ...draft.value, enabled: value, blocking_enabled: value ? draft.value.blocking_enabled : false placeholder)
placeholder
function setBlocking(value: boolean) {
  if (!draft.value || !draft.value.enabled) return
  if (value && !draft.value.blocking_enabled) { showBlockingConfirmation.value = true; return placeholder
  replaceDraft({ ...draft.value, blocking_enabled: value placeholder)
placeholder
function confirmBlocking() {
  showBlockingConfirmation.value = false
  if (draft.value) replaceDraft({ ...draft.value, blocking_enabled: true placeholder)
placeholder
function resetDraft() {
  if (serverConfig.value) draft.value = cloneData(serverConfig.value)
placeholder
async function saveConfig() {
  if (!draft.value || !dirty.value) return
  loading.saving = true
  try {
    const saved = await promptAuditAPI.updateConfig(buildUpdateRequest(draft.value))
    serverConfig.value = configToDraft(saved)
    draft.value = configToDraft(saved)
    appStore.showSuccess(t('admin.promptAudit.messages.saved'))
    await loadRuntime()
  placeholder catch (error) {
    const code = extractApiErrorCode(error)
    appStore.showError(errorMessage(error, code === 'prompt_audit_config_conflict' ? 'admin.promptAudit.errors.prompt_audit_config_conflict' : 'admin.promptAudit.errors.saveConfig'))
  placeholder finally {
    loading.saving = false
  placeholder
placeholder
async function runProbe(endpoint: PromptAuditEndpointDraft) {
  if (probingIds.value.includes(endpoint.id)) return
  probingIds.value = [...probingIds.value, endpoint.id]
  try {
    const result = await promptAuditAPI.probeEndpoint(endpoint)
    probeResults[endpoint.id] = result
    if (result.ok) appStore.showSuccess(t('admin.promptAudit.messages.probeSucceeded'))
    else appStore.showError(`${result.error_code || result.statusplaceholder: ${result.messageplaceholder`)
  placeholder catch (error) {
    appStore.showError(errorMessage(error, 'admin.promptAudit.errors.probe'))
  placeholder finally {
    probingIds.value = probingIds.value.filter((id) => id !== endpoint.id)
  placeholder
placeholder

function handleFiltersChanged(value: PromptEventFilters) {
  filters.value = cloneData(value)
  deletePreview.value = null
placeholder
function applyEventFilters(value: PromptEventFilters) {
  filters.value = cloneData(value)
  appliedFilters.value = cloneData(value)
  events.page = 1
  deletePreview.value = null
  void loadEvents()
placeholder
function changePage(value: number) { events.page = value; void loadEvents() placeholder
function changePageSize(value: number) { events.page_size = value; events.page = 1; void loadEvents() placeholder
async function openEvent(id: number) {
  showEventDetail.value = true
  loading.detail = true
  activeEvent.value = null
  try { activeEvent.value = await promptAuditAPI.getEvent(id) placeholder
  catch (error) { appStore.showError(errorMessage(error, 'admin.promptAudit.errors.loadDetail')); showEventDetail.value = false placeholder
  finally { loading.detail = false placeholder
placeholder
function closeEventDetail() { showEventDetail.value = false; activeEvent.value = null placeholder
function requestSingleDelete(id: number) { deleteRequest.mode = 'single'; deleteRequest.ids = [id] placeholder
function requestBatchDelete() { if (selectedEventIds.value.length) { deleteRequest.mode = 'batch'; deleteRequest.ids = [...selectedEventIds.value] placeholder placeholder
function clearDeleteRequest() { deleteRequest.mode = ''; deleteRequest.ids = [] placeholder
async function confirmIDDelete() {
  const mode = deleteRequest.mode
  const ids = [...deleteRequest.ids]
  clearDeleteRequest()
  if (!mode || ids.length === 0) return
  loading.deleting = true
  try {
    const result = mode === 'single' ? await promptAuditAPI.deleteEvent(ids[0]) : await promptAuditAPI.batchDeleteEvents(ids)
    appStore.showSuccess(t('admin.promptAudit.messages.deleted', { count: result.deleted_events placeholder))
    await Promise.allSettled([loadEvents(), loadRuntime()])
  placeholder catch (error) { appStore.showError(errorMessage(error, 'admin.promptAudit.errors.delete')) placeholder
  finally { loading.deleting = false placeholder
placeholder
async function requestFilterDeletePreview() {
  loading.deleting = true
  try { deletePreview.value = await promptAuditAPI.previewDelete(filters.value) placeholder
  catch (error) { appStore.showError(errorMessage(error, 'admin.promptAudit.errors.previewDelete')) placeholder
  finally { loading.deleting = false placeholder
placeholder
async function confirmFilterDelete() {
  if (!deletePreview.value) return
  const preview = deletePreview.value
  loading.deleting = true
  try {
    const result = await promptAuditAPI.deleteEventsByFilter(filters.value, preview)
    deletePreview.value = null
    appStore.showSuccess(t('admin.promptAudit.messages.deleted', { count: result.deleted_events placeholder))
    await Promise.allSettled([loadEvents(), loadRuntime()])
  placeholder catch (error) {
    deletePreview.value = null
    appStore.showError(errorMessage(error, 'admin.promptAudit.errors.deleteConfirmation'))
  placeholder finally { loading.deleting = false placeholder
placeholder
function formatDate(value: string): string {
  return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'medium' placeholder).format(new Date(value))
placeholder

onMounted(loadInitial)
</script>
