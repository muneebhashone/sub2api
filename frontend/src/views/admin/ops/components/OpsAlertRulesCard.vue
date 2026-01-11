<script setup lang="ts">
import { computed, onMounted, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Select from '@/components/common/Select.vue'
import { opsAPI placeholder from '@/api/admin/ops'
import type { AlertRule, MetricType, Operator placeholder from '../types'
import type { OpsSeverity placeholder from '@/api/admin/ops'
import { formatDateTime placeholder from '../utils/opsFormatters'

const { t placeholder = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const rules = ref<AlertRule[]>([])

async function load() {
  loading.value = true
  try {
    rules.value = await opsAPI.listAlertRules()
  placeholder catch (err: any) {
    console.error('[OpsAlertRulesCard] Failed to load rules', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertRules.loadFailed'))
    rules.value = []
  placeholder finally {
    loading.value = false
  placeholder
placeholder

onMounted(() => {
  load()
placeholder)

const sortedRules = computed(() => {
  return [...rules.value].sort((a, b) => (b.id || 0) - (a.id || 0))
placeholder)

const showEditor = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)
const draft = ref<AlertRule | null>(null)

const metricOptions = computed(() => {
  const items: Array<{ value: MetricType; label: string placeholder> = [
    { value: 'success_rate', label: t('admin.ops.alertRules.metrics.successRate') placeholder,
    { value: 'error_rate', label: t('admin.ops.alertRules.metrics.errorRate') placeholder,
    { value: 'p95_latency_ms', label: t('admin.ops.alertRules.metrics.p95') placeholder,
    { value: 'p99_latency_ms', label: t('admin.ops.alertRules.metrics.p99') placeholder,
    { value: 'cpu_usage_percent', label: t('admin.ops.alertRules.metrics.cpu') placeholder,
    { value: 'memory_usage_percent', label: t('admin.ops.alertRules.metrics.memory') placeholder,
    { value: 'concurrency_queue_depth', label: t('admin.ops.alertRules.metrics.queueDepth') placeholder
  ]
  return items
placeholder)

const operatorOptions = computed(() => {
  const ops: Operator[] = ['>', '>=', '<', '<=', '==', '!=']
  return ops.map((o) => ({ value: o, label: o placeholder))
placeholder)

const severityOptions = computed(() => {
  const sev: OpsSeverity[] = ['P0', 'P1', 'P2', 'P3']
  return sev.map((s) => ({ value: s, label: s placeholder))
placeholder)

const windowOptions = computed(() => {
  const windows = [1, 5, 60]
  return windows.map((m) => ({ value: m, label: `${mplaceholderm` placeholder))
placeholder)

function newRuleDraft(): AlertRule {
  return {
    name: '',
    description: '',
    enabled: true,
    metric_type: 'error_rate',
    operator: '>',
    threshold: 1,
    window_minutes: 1,
    sustained_minutes: 2,
    severity: 'P1',
    cooldown_minutes: 10,
    notify_email: true
  placeholder
placeholder

function openCreate() {
  editingId.value = null
  draft.value = newRuleDraft()
  showEditor.value = true
placeholder

function openEdit(rule: AlertRule) {
  editingId.value = rule.id ?? null
  draft.value = JSON.parse(JSON.stringify(rule))
  showEditor.value = true
placeholder

const editorValidation = computed(() => {
  const errors: string[] = []
  const r = draft.value
  if (!r) return { valid: true, errors placeholder
  if (!r.name || !r.name.trim()) errors.push(t('admin.ops.alertRules.validation.nameRequired'))
  if (!r.metric_type) errors.push(t('admin.ops.alertRules.validation.metricRequired'))
  if (!r.operator) errors.push(t('admin.ops.alertRules.validation.operatorRequired'))
  if (!(typeof r.threshold === 'number' && Number.isFinite(r.threshold)))
    errors.push(t('admin.ops.alertRules.validation.thresholdRequired'))
  if (!(typeof r.window_minutes === 'number' && Number.isFinite(r.window_minutes) && [1, 5, 60].includes(r.window_minutes))) {
    errors.push(t('admin.ops.alertRules.validation.windowRange'))
  placeholder
  if (!(typeof r.sustained_minutes === 'number' && Number.isFinite(r.sustained_minutes) && r.sustained_minutes >= 1 && r.sustained_minutes <= 1440)) {
    errors.push(t('admin.ops.alertRules.validation.sustainedRange'))
  placeholder
  if (!(typeof r.cooldown_minutes === 'number' && Number.isFinite(r.cooldown_minutes) && r.cooldown_minutes >= 0 && r.cooldown_minutes <= 1440)) {
    errors.push(t('admin.ops.alertRules.validation.cooldownRange'))
  placeholder
  return { valid: errors.length === 0, errors placeholder
placeholder)

async function save() {
  if (!draft.value) return
  if (!editorValidation.value.valid) {
    appStore.showError(editorValidation.value.errors[0] || t('admin.ops.alertRules.validation.invalid'))
    return
  placeholder
  saving.value = true
  try {
    if (editingId.value) {
      await opsAPI.updateAlertRule(editingId.value, draft.value)
    placeholder else {
      await opsAPI.createAlertRule(draft.value)
    placeholder
    showEditor.value = false
    draft.value = null
    editingId.value = null
    await load()
    appStore.showSuccess(t('admin.ops.alertRules.saveSuccess'))
  placeholder catch (err: any) {
    console.error('[OpsAlertRulesCard] Failed to save rule', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertRules.saveFailed'))
  placeholder finally {
    saving.value = false
  placeholder
placeholder

const showDeleteConfirm = ref(false)
const pendingDelete = ref<AlertRule | null>(null)

function requestDelete(rule: AlertRule) {
  pendingDelete.value = rule
  showDeleteConfirm.value = true
placeholder

async function confirmDelete() {
  if (!pendingDelete.value?.id) return
  try {
    await opsAPI.deleteAlertRule(pendingDelete.value.id)
    showDeleteConfirm.value = false
    pendingDelete.value = null
    await load()
    appStore.showSuccess(t('admin.ops.alertRules.deleteSuccess'))
  placeholder catch (err: any) {
    console.error('[OpsAlertRulesCard] Failed to delete rule', err)
    appStore.showError(err?.response?.data?.detail || t('admin.ops.alertRules.deleteFailed'))
  placeholder
placeholder

function cancelDelete() {
  showDeleteConfirm.value = false
  pendingDelete.value = null
placeholder
</script>

<template>
  <div class="rounded-3xl bg-white p-6 shadow-sm ring-1 ring-gray-900/5 dark:bg-dark-800 dark:ring-dark-700">
    <div class="mb-4 flex items-start justify-between gap-4">
      <div>
        <h3 class="text-sm font-bold text-gray-900 dark:text-white">{{ t('admin.ops.alertRules.title') placeholderplaceholder</h3>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.ops.alertRules.description') placeholderplaceholder</p>
      </div>

      <div class="flex items-center gap-2">
        <button class="btn btn-sm btn-primary" :disabled="loading" @click="openCreate">
          {{ t('admin.ops.alertRules.create') placeholderplaceholder
        </button>
        <button
          class="flex items-center gap-1.5 rounded-lg bg-gray-100 px-3 py-1.5 text-xs font-bold text-gray-700 transition-colors hover:bg-gray-200 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-dark-700 dark:text-gray-300 dark:hover:bg-dark-600"
          :disabled="loading"
          @click="load"
        >
          <svg class="h-3.5 w-3.5" :class="{ 'animate-spin': loading placeholder" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ t('common.refresh') placeholderplaceholder
        </button>
      </div>
    </div>

    <div v-if="loading" class="py-10 text-center text-sm text-gray-500 dark:text-gray-400">
      {{ t('admin.ops.alertRules.loading') placeholderplaceholder
    </div>

    <div v-else-if="sortedRules.length === 0" class="rounded-xl border border-dashed border-gray-200 p-8 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-gray-400">
      {{ t('admin.ops.alertRules.empty') placeholderplaceholder
    </div>

    <div v-else class="overflow-hidden rounded-xl border border-gray-200 dark:border-dark-700">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-dark-700">
          <thead class="bg-gray-50 dark:bg-dark-900">
            <tr>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertRules.table.name') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertRules.table.metric') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertRules.table.severity') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-left text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertRules.table.enabled') placeholderplaceholder
              </th>
              <th class="px-4 py-3 text-right text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-gray-400">
                {{ t('admin.ops.alertRules.table.actions') placeholderplaceholder
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white dark:divide-dark-700 dark:bg-dark-800">
            <tr v-for="row in sortedRules" :key="row.id" class="hover:bg-gray-50 dark:hover:bg-dark-700/50">
              <td class="px-4 py-3">
                <div class="text-xs font-bold text-gray-900 dark:text-white">{{ row.name placeholderplaceholder</div>
                <div v-if="row.description" class="mt-0.5 line-clamp-2 text-[11px] text-gray-500 dark:text-gray-400">
                  {{ row.description placeholderplaceholder
                </div>
                <div v-if="row.updated_at" class="mt-1 text-[10px] text-gray-400">
                  {{ formatDateTime(row.updated_at) placeholderplaceholder
                </div>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-xs text-gray-700 dark:text-gray-200">
                <span class="font-mono">{{ row.metric_type placeholderplaceholder</span>
                <span class="mx-1 text-gray-400">{{ row.operator placeholderplaceholder</span>
                <span class="font-mono">{{ row.threshold placeholderplaceholder</span>
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-xs font-bold text-gray-700 dark:text-gray-200">
                {{ row.severity placeholderplaceholder
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-xs text-gray-700 dark:text-gray-200">
                {{ row.enabled ? t('common.enabled') : t('common.disabled') placeholderplaceholder
              </td>
              <td class="whitespace-nowrap px-4 py-3 text-right text-xs">
                <button class="btn btn-sm btn-secondary" @click="openEdit(row)">{{ t('common.edit') placeholderplaceholder</button>
                <button class="ml-2 btn btn-sm btn-danger" @click="requestDelete(row)">{{ t('common.delete') placeholderplaceholder</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <BaseDialog
      :show="showEditor"
      :title="editingId ? t('admin.ops.alertRules.editTitle') : t('admin.ops.alertRules.createTitle')"
      width="wide"
      @close="showEditor = false"
    >
      <div class="space-y-4">
        <div v-if="!editorValidation.valid" class="rounded-xl bg-red-50 p-4 text-xs text-red-700 dark:bg-red-900/30 dark:text-red-300">
          <div class="font-bold">{{ t('admin.ops.alertRules.validation.title') placeholderplaceholder</div>
          <ul class="mt-1 list-disc pl-5">
            <li v-for="e in editorValidation.errors" :key="e">{{ e placeholderplaceholder</li>
          </ul>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div class="md:col-span-2">
            <label class="input-label">{{ t('admin.ops.alertRules.form.name') placeholderplaceholder</label>
            <input v-model="draft!.name" class="input" type="text" />
          </div>

          <div class="md:col-span-2">
            <label class="input-label">{{ t('admin.ops.alertRules.form.description') placeholderplaceholder</label>
            <input v-model="draft!.description" class="input" type="text" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.metric') placeholderplaceholder</label>
            <Select v-model="draft!.metric_type" :options="metricOptions" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.operator') placeholderplaceholder</label>
            <Select v-model="draft!.operator" :options="operatorOptions" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.threshold') placeholderplaceholder</label>
            <input v-model.number="draft!.threshold" class="input" type="number" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.severity') placeholderplaceholder</label>
            <Select v-model="draft!.severity" :options="severityOptions" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.window') placeholderplaceholder</label>
            <Select v-model="draft!.window_minutes" :options="windowOptions" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.sustained') placeholderplaceholder</label>
            <input v-model.number="draft!.sustained_minutes" class="input" type="number" min="1" max="1440" />
          </div>

          <div>
            <label class="input-label">{{ t('admin.ops.alertRules.form.cooldown') placeholderplaceholder</label>
            <input v-model.number="draft!.cooldown_minutes" class="input" type="number" min="0" max="1440" />
          </div>

          <div class="flex items-center justify-between rounded-xl bg-gray-50 px-4 py-3 dark:bg-dark-800/50 md:col-span-2">
            <span class="text-xs font-bold text-gray-700 dark:text-gray-200">{{ t('admin.ops.alertRules.form.enabled') placeholderplaceholder</span>
            <input v-model="draft!.enabled" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
          </div>

          <div class="flex items-center justify-between rounded-xl bg-gray-50 px-4 py-3 dark:bg-dark-800/50 md:col-span-2">
            <span class="text-xs font-bold text-gray-700 dark:text-gray-200">{{ t('admin.ops.alertRules.form.notifyEmail') placeholderplaceholder</span>
            <input v-model="draft!.notify_email" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <button class="btn btn-secondary" :disabled="saving" @click="showEditor = false">
            {{ t('common.cancel') placeholderplaceholder
          </button>
          <button class="btn btn-primary" :disabled="saving" @click="save">
            {{ saving ? t('common.saving') : t('common.save') placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>

    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('admin.ops.alertRules.deleteConfirmTitle')"
      :message="t('admin.ops.alertRules.deleteConfirmMessage')"
      :confirmText="t('common.delete')"
      :cancelText="t('common.cancel')"
      @confirm="confirmDelete"
      @cancel="cancelDelete"
    />
  </div>
</template>
