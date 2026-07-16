<template>
    <div class="space-y-6">
      <!-- S3 Storage Config -->
      <div class="card p-6">
        <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">
              {{ t('admin.backup.s3.title') placeholderplaceholder
            </h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.backup.s3.descriptionPrefix') placeholderplaceholder
              <button type="button" class="text-primary-600 underline hover:text-primary-700 dark:text-primary-400 dark:hover:text-primary-300" @click="showR2Guide = true">Cloudflare R2</button>
              {{ t('admin.backup.s3.descriptionSuffix') placeholderplaceholder
            </p>
          </div>
        </div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.s3.endpoint') placeholderplaceholder</label>
            <input v-model="s3Form.endpoint" class="input w-full" placeholder="https://<account_id>.r2.cloudflarestorage.com" />
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.s3.region') placeholderplaceholder</label>
            <input v-model="s3Form.region" class="input w-full" placeholder="auto" />
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.s3.bucket') placeholderplaceholder</label>
            <input v-model="s3Form.bucket" class="input w-full" />
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.s3.prefix') placeholderplaceholder</label>
            <input v-model="s3Form.prefix" class="input w-full" placeholder="backups/" />
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.s3.accessKeyId') placeholderplaceholder</label>
            <input v-model="s3Form.access_key_id" class="input w-full" />
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.s3.secretAccessKey') placeholderplaceholder</label>
            <input v-model="s3Form.secret_access_key" type="password" class="input w-full" :placeholder="s3SecretConfigured ? t('admin.backup.s3.secretConfigured') : ''" />
          </div>
          <label class="inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 md:col-span-2">
            <input v-model="s3Form.force_path_style" type="checkbox" />
            <span>{{ t('admin.backup.s3.forcePathStyle') placeholderplaceholder</span>
          </label>
        </div>
        <div class="mt-4 flex flex-wrap gap-2">
          <button type="button" class="btn btn-secondary btn-sm" :disabled="testingS3" @click="testS3">
            {{ testingS3 ? t('common.loading') : t('admin.backup.s3.testConnection') placeholderplaceholder
          </button>
          <button type="button" class="btn btn-primary btn-sm" :disabled="savingS3" @click="saveS3Config">
            {{ savingS3 ? t('common.loading') : t('common.save') placeholderplaceholder
          </button>
        </div>
      </div>

      <!-- Schedule Config -->
      <div class="card p-6">
        <div class="mb-4">
          <h3 class="text-base font-semibold text-gray-900 dark:text-white">
            {{ t('admin.backup.schedule.title') placeholderplaceholder
          </h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ t('admin.backup.schedule.description') placeholderplaceholder
          </p>
        </div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <label class="inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 md:col-span-2">
            <input v-model="scheduleForm.enabled" type="checkbox" />
            <span>{{ t('admin.backup.schedule.enabled') placeholderplaceholder</span>
          </label>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.schedule.cronExpr') placeholderplaceholder</label>
            <input v-model="scheduleForm.cron_expr" class="input w-full" placeholder="0 2 * * *" />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.backup.schedule.cronHint') placeholderplaceholder</p>
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.schedule.retainDays') placeholderplaceholder</label>
            <input v-model.number="scheduleForm.retain_days" type="number" min="0" class="input w-full" />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.backup.schedule.retainDaysHint') placeholderplaceholder</p>
          </div>
          <div>
            <label class="mb-1 block text-xs font-medium text-gray-600 dark:text-gray-400">{{ t('admin.backup.schedule.retainCount') placeholderplaceholder</label>
            <input v-model.number="scheduleForm.retain_count" type="number" min="0" class="input w-full" />
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('admin.backup.schedule.retainCountHint') placeholderplaceholder</p>
          </div>
        </div>
        <div class="mt-4">
          <button type="button" class="btn btn-primary btn-sm" :disabled="savingSchedule" @click="saveSchedule">
            {{ savingSchedule ? t('common.loading') : t('common.save') placeholderplaceholder
          </button>
        </div>
      </div>

      <!-- Backup Operations -->
      <div class="card p-6">
        <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h3 class="text-base font-semibold text-gray-900 dark:text-white">
              {{ t('admin.backup.operations.title') placeholderplaceholder
            </h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
              {{ t('admin.backup.operations.description') placeholderplaceholder
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <div class="flex items-center gap-1">
              <label class="text-xs text-gray-600 dark:text-gray-400">{{ t('admin.backup.operations.expireDays') placeholderplaceholder</label>
              <input v-model.number="manualExpireDays" type="number" min="0" class="input w-20 text-xs" />
            </div>
            <button type="button" class="btn btn-primary btn-sm" :disabled="creatingBackup" @click="createBackup">
              {{ creatingBackup ? t('admin.backup.operations.backing') : t('admin.backup.operations.createBackup') placeholderplaceholder
            </button>
            <button type="button" class="btn btn-secondary btn-sm" :disabled="loadingBackups" @click="loadBackups">
              {{ loadingBackups ? t('common.loading') : t('common.refresh') placeholderplaceholder
            </button>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full min-w-[800px] text-sm">
            <thead>
              <tr class="border-b border-gray-200 text-left text-xs uppercase tracking-wide text-gray-500 dark:border-dark-700 dark:text-gray-400">
                <th class="py-2 pr-4">ID</th>
                <th class="py-2 pr-4">{{ t('admin.backup.columns.status') placeholderplaceholder</th>
                <th class="py-2 pr-4">{{ t('admin.backup.columns.fileName') placeholderplaceholder</th>
                <th class="py-2 pr-4">{{ t('admin.backup.columns.size') placeholderplaceholder</th>
                <th class="py-2 pr-4">{{ t('admin.backup.columns.expiresAt') placeholderplaceholder</th>
                <th class="py-2 pr-4">{{ t('admin.backup.columns.triggeredBy') placeholderplaceholder</th>
                <th class="py-2 pr-4">{{ t('admin.backup.columns.startedAt') placeholderplaceholder</th>
                <th class="py-2">{{ t('admin.backup.columns.actions') placeholderplaceholder</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in backups" :key="record.id" class="border-b border-gray-100 align-top dark:border-dark-800">
                <td class="py-3 pr-4 font-mono text-xs">{{ record.id placeholderplaceholder</td>
                <td class="py-3 pr-4">
                  <span
                    class="rounded px-2 py-0.5 text-xs"
                    :class="statusClass(record.status)"
                  >
                    {{ record.status === 'running' && record.progress
                      ? t(`admin.backup.progress.${record.progressplaceholder`)
                      : t(`admin.backup.status.${record.statusplaceholder`) placeholderplaceholder
                  </span>
                </td>
                <td class="py-3 pr-4 text-xs">{{ record.file_name placeholderplaceholder</td>
                <td class="py-3 pr-4 text-xs">{{ formatSize(record.size_bytes) placeholderplaceholder</td>
                <td class="py-3 pr-4 text-xs">
                  {{ record.expires_at ? formatDate(record.expires_at) : t('admin.backup.neverExpire') placeholderplaceholder
                </td>
                <td class="py-3 pr-4 text-xs">
                  {{ record.triggered_by === 'scheduled' ? t('admin.backup.trigger.scheduled') : t('admin.backup.trigger.manual') placeholderplaceholder
                </td>
                <td class="py-3 pr-4 text-xs">{{ formatDate(record.started_at) placeholderplaceholder</td>
                <td class="py-3 text-xs">
                  <div class="flex flex-wrap gap-1">
                    <button
                      v-if="record.status === 'completed'"
                      type="button"
                      class="btn btn-secondary btn-xs"
                      @click="downloadBackup(record.id)"
                    >
                      {{ t('admin.backup.actions.download') placeholderplaceholder
                    </button>
                    <button
                      v-if="record.status === 'completed'"
                      type="button"
                      class="btn btn-secondary btn-xs"
                      :disabled="restoringId === record.id"
                      @click="restoreBackup(record.id)"
                    >
                      {{ restoringId === record.id ? t('common.loading') : t('admin.backup.actions.restore') placeholderplaceholder
                    </button>
                    <button
                      type="button"
                      class="btn btn-danger btn-xs"
                      @click="removeBackup(record.id)"
                    >
                      {{ t('common.delete') placeholderplaceholder
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="backups.length === 0">
                <td colspan="8" class="py-6 text-center text-sm text-gray-500 dark:text-gray-400">
                  {{ t('admin.backup.empty') placeholderplaceholder
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Cloudflare R2 Setup Guide Modal -->
    <teleport to="body">
      <transition name="modal">
        <div v-if="showR2Guide" class="fixed inset-0 z-50 flex items-center justify-center p-4" @mousedown.self="showR2Guide = false">
          <div class="fixed inset-0 bg-black/50" @click="showR2Guide = false"></div>
          <div class="relative max-h-[85vh] w-full max-w-2xl overflow-y-auto rounded-xl bg-white p-6 shadow-2xl dark:bg-dark-800">
            <button type="button" class="absolute right-4 top-4 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200" @click="showR2Guide = false">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>

            <h2 class="mb-4 text-lg font-bold text-gray-900 dark:text-white">{{ t('admin.backup.r2Guide.title') placeholderplaceholder</h2>
            <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">{{ t('admin.backup.r2Guide.intro') placeholderplaceholder</p>

            <!-- Step 1 -->
            <div class="mb-5">
              <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-gray-900 dark:text-white">
                <span class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-100 text-xs font-bold text-primary-700 dark:bg-primary-900/40 dark:text-primary-300">1</span>
                {{ t('admin.backup.r2Guide.step1.title') placeholderplaceholder
              </h3>
              <ol class="ml-8 list-decimal space-y-1 text-sm text-gray-600 dark:text-gray-300">
                <li>{{ t('admin.backup.r2Guide.step1.line1') placeholderplaceholder</li>
                <li>{{ t('admin.backup.r2Guide.step1.line2') placeholderplaceholder</li>
                <li>{{ t('admin.backup.r2Guide.step1.line3') placeholderplaceholder</li>
              </ol>
            </div>

            <!-- Step 2 -->
            <div class="mb-5">
              <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-gray-900 dark:text-white">
                <span class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-100 text-xs font-bold text-primary-700 dark:bg-primary-900/40 dark:text-primary-300">2</span>
                {{ t('admin.backup.r2Guide.step2.title') placeholderplaceholder
              </h3>
              <ol class="ml-8 list-decimal space-y-1 text-sm text-gray-600 dark:text-gray-300">
                <li>{{ t('admin.backup.r2Guide.step2.line1') placeholderplaceholder</li>
                <li>{{ t('admin.backup.r2Guide.step2.line2') placeholderplaceholder</li>
                <li>{{ t('admin.backup.r2Guide.step2.line3') placeholderplaceholder</li>
                <li>{{ t('admin.backup.r2Guide.step2.line4') placeholderplaceholder</li>
              </ol>
              <div class="mt-2 rounded-lg bg-amber-50 p-3 text-xs text-amber-700 dark:bg-amber-900/20 dark:text-amber-300">
                {{ t('admin.backup.r2Guide.step2.warning') placeholderplaceholder
              </div>
            </div>

            <!-- Step 3 -->
            <div class="mb-5">
              <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-gray-900 dark:text-white">
                <span class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-100 text-xs font-bold text-primary-700 dark:bg-primary-900/40 dark:text-primary-300">3</span>
                {{ t('admin.backup.r2Guide.step3.title') placeholderplaceholder
              </h3>
              <p class="ml-8 text-sm text-gray-600 dark:text-gray-300">{{ t('admin.backup.r2Guide.step3.desc') placeholderplaceholder</p>
              <code class="ml-8 mt-1 block rounded bg-gray-100 px-3 py-2 text-xs text-gray-800 dark:bg-dark-700 dark:text-gray-200">https://&lt;{{ t('admin.backup.r2Guide.step3.accountId') placeholderplaceholder&gt;.r2.cloudflarestorage.com</code>
            </div>

            <!-- Step 4: Fill form -->
            <div class="mb-5">
              <h3 class="mb-2 flex items-center gap-2 text-sm font-semibold text-gray-900 dark:text-white">
                <span class="flex h-6 w-6 items-center justify-center rounded-full bg-primary-100 text-xs font-bold text-primary-700 dark:bg-primary-900/40 dark:text-primary-300">4</span>
                {{ t('admin.backup.r2Guide.step4.title') placeholderplaceholder
              </h3>
              <div class="ml-8 overflow-hidden rounded-lg border border-gray-200 dark:border-dark-600">
                <table class="w-full text-sm">
                  <tbody>
                    <tr v-for="(row, i) in r2ConfigRows" :key="i" class="border-b border-gray-100 dark:border-dark-700 last:border-0">
                      <td class="whitespace-nowrap bg-gray-50 px-3 py-2 font-medium text-gray-700 dark:bg-dark-700 dark:text-gray-300">{{ row.field placeholderplaceholder</td>
                      <td class="px-3 py-2 text-gray-600 dark:text-gray-400"><code class="text-xs">{{ row.value placeholderplaceholder</code></td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Free tier note -->
            <div class="rounded-lg bg-green-50 p-3 text-xs text-green-700 dark:bg-green-900/20 dark:text-green-300">
              {{ t('admin.backup.r2Guide.freeTier') placeholderplaceholder
            </div>

            <div class="mt-4 text-right">
              <button type="button" class="btn btn-primary btn-sm" @click="showR2Guide = false">{{ t('common.close') placeholderplaceholder</button>
            </div>
          </div>
        </div>
      </transition>
    </teleport>
    <TotpStepUpDialog :controller="backupStepUp" />
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { adminAPI placeholder from '@/api'
import { useAppStore placeholder from '@/stores'
import type { BackupS3Config, BackupScheduleConfig, BackupRecord placeholder from '@/api/admin/backup'
import { useStepUp, isStepUpBlocked, isStepUpCancelled, stepUpBlockReason placeholder from '@/composables/useStepUp'
import TotpStepUpDialog from '@/components/auth/TotpStepUpDialog.vue'

const { t placeholder = useI18n()
const appStore = useAppStore()
const backupStepUp = useStepUp()

// 敏感操作被 2FA 门控拦截时的统一提示。
function reportStepUpBlocked(error: unknown): boolean {
  if (!isStepUpBlocked(error)) return false
  appStore.showError(
    stepUpBlockReason(error) === 'STEP_UP_ADMIN_API_KEY_FORBIDDEN'
      ? t('stepUp.adminApiKeyForbidden')
      : t('stepUp.notEnabled')
  )
  return true
placeholder

// S3 config
const s3Form = ref<BackupS3Config>({
  endpoint: '',
  region: 'auto',
  bucket: '',
  access_key_id: '',
  secret_access_key: '',
  prefix: 'backups/',
  force_path_style: false,
placeholder)
const s3SecretConfigured = ref(false)
const savingS3 = ref(false)
const testingS3 = ref(false)

// Schedule config
const scheduleForm = ref<BackupScheduleConfig>({
  enabled: false,
  cron_expr: '0 2 * * *',
  retain_days: 14,
  retain_count: 10,
placeholder)
const savingSchedule = ref(false)

// Backups
const backups = ref<BackupRecord[]>([])
const loadingBackups = ref(false)
const creatingBackup = ref(false)
const restoringId = ref('')
const manualExpireDays = ref(14)

// Polling
const pollingTimer = ref<ReturnType<typeof setInterval> | null>(null)
const restoringPollingTimer = ref<ReturnType<typeof setInterval> | null>(null)
const MAX_POLL_COUNT = 900

function updateRecordInList(updated: BackupRecord) {
  const idx = backups.value.findIndex(r => r.id === updated.id)
  if (idx >= 0) {
    backups.value[idx] = updated
  placeholder
placeholder

function startPolling(backupId: string) {
  stopPolling()
  let count = 0
  pollingTimer.value = setInterval(async () => {
    if (count++ >= MAX_POLL_COUNT) {
      stopPolling()
      creatingBackup.value = false
      appStore.showWarning(t('admin.backup.operations.backupRunning'))
      return
    placeholder
    try {
      const record = await adminAPI.backup.getBackup(backupId)
      updateRecordInList(record)
      if (record.status === 'completed' || record.status === 'failed') {
        stopPolling()
        creatingBackup.value = false
        if (record.status === 'completed') {
          appStore.showSuccess(t('admin.backup.operations.backupCreated'))
        placeholder else {
          appStore.showError(record.error_message || t('admin.backup.operations.backupFailed'))
        placeholder
        await loadBackups()
      placeholder
    placeholder catch {
      // 轮询失败时不中断
    placeholder
  placeholder, 2000)
placeholder

function stopPolling() {
  if (pollingTimer.value) {
    clearInterval(pollingTimer.value)
    pollingTimer.value = null
  placeholder
placeholder

function startRestorePolling(backupId: string) {
  stopRestorePolling()
  let count = 0
  restoringPollingTimer.value = setInterval(async () => {
    if (count++ >= MAX_POLL_COUNT) {
      stopRestorePolling()
      restoringId.value = ''
      appStore.showWarning(t('admin.backup.operations.restoreRunning'))
      return
    placeholder
    try {
      const record = await adminAPI.backup.getBackup(backupId)
      updateRecordInList(record)
      if (record.restore_status === 'completed' || record.restore_status === 'failed') {
        stopRestorePolling()
        restoringId.value = ''
        if (record.restore_status === 'completed') {
          appStore.showSuccess(t('admin.backup.actions.restoreSuccess'))
        placeholder else {
          appStore.showError(record.restore_error || t('admin.backup.operations.restoreFailed'))
        placeholder
        await loadBackups()
      placeholder
    placeholder catch {
      // 轮询失败时不中断
    placeholder
  placeholder, 2000)
placeholder

function stopRestorePolling() {
  if (restoringPollingTimer.value) {
    clearInterval(restoringPollingTimer.value)
    restoringPollingTimer.value = null
  placeholder
placeholder

function handleVisibilityChange() {
  if (document.hidden) {
    stopPolling()
    stopRestorePolling()
  placeholder else {
    // 标签页恢复时刷新列表，检查是否仍有活跃操作
    loadBackups().then(() => {
      const running = backups.value.find(r => r.status === 'running')
      if (running) {
        creatingBackup.value = true
        startPolling(running.id)
      placeholder
      const restoring = backups.value.find(r => r.restore_status === 'running')
      if (restoring) {
        restoringId.value = restoring.id
        startRestorePolling(restoring.id)
      placeholder
    placeholder)
  placeholder
placeholder

// R2 guide
const showR2Guide = ref(false)
const r2ConfigRows = computed(() => [
  { field: t('admin.backup.s3.endpoint'), value: 'https://<account_id>.r2.cloudflarestorage.com' placeholder,
  { field: t('admin.backup.s3.region'), value: 'auto' placeholder,
  { field: t('admin.backup.s3.bucket'), value: t('admin.backup.r2Guide.step4.bucketValue') placeholder,
  { field: t('admin.backup.s3.prefix'), value: 'backups/' placeholder,
  { field: 'Access Key ID', value: t('admin.backup.r2Guide.step4.fromStep2') placeholder,
  { field: 'Secret Access Key', value: t('admin.backup.r2Guide.step4.fromStep2') placeholder,
  { field: t('admin.backup.s3.forcePathStyle'), value: t('admin.backup.r2Guide.step4.unchecked') placeholder,
])

async function loadS3Config() {
  try {
    const cfg = await adminAPI.backup.getS3Config()
    s3Form.value = {
      endpoint: cfg.endpoint || '',
      region: cfg.region || 'auto',
      bucket: cfg.bucket || '',
      access_key_id: cfg.access_key_id || '',
      secret_access_key: '',
      prefix: cfg.prefix || 'backups/',
      force_path_style: cfg.force_path_style,
    placeholder
    s3SecretConfigured.value = Boolean(cfg.access_key_id)
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder
placeholder

async function saveS3Config() {
  savingS3.value = true
  try {
    await adminAPI.backup.updateS3Config(s3Form.value)
    appStore.showSuccess(t('admin.backup.s3.saved'))
    await loadS3Config()
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder finally {
    savingS3.value = false
  placeholder
placeholder

async function testS3() {
  testingS3.value = true
  try {
    const result = await adminAPI.backup.testS3Connection(s3Form.value)
    if (result.ok) {
      appStore.showSuccess(result.message || t('admin.backup.s3.testSuccess'))
    placeholder else {
      appStore.showError(result.message || t('admin.backup.s3.testFailed'))
    placeholder
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder finally {
    testingS3.value = false
  placeholder
placeholder

async function loadSchedule() {
  try {
    const cfg = await adminAPI.backup.getSchedule()
    scheduleForm.value = {
      enabled: cfg.enabled,
      cron_expr: cfg.cron_expr || '0 2 * * *',
      retain_days: cfg.retain_days || 14,
      retain_count: cfg.retain_count || 10,
    placeholder
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder
placeholder

async function saveSchedule() {
  savingSchedule.value = true
  try {
    await adminAPI.backup.updateSchedule(scheduleForm.value)
    appStore.showSuccess(t('admin.backup.schedule.saved'))
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder finally {
    savingSchedule.value = false
  placeholder
placeholder

async function loadBackups() {
  loadingBackups.value = true
  try {
    const result = await adminAPI.backup.listBackups()
    backups.value = result.items || []
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder finally {
    loadingBackups.value = false
  placeholder
placeholder

async function createBackup() {
  creatingBackup.value = true
  try {
    const record = await backupStepUp.run(() => adminAPI.backup.createBackup({ expire_days: manualExpireDays.value placeholder))
    // 插入到列表顶部
    backups.value.unshift(record)
    startPolling(record.id)
  placeholder catch (error: any) {
    if (isStepUpCancelled(error)) {
      creatingBackup.value = false
      return
    placeholder
    if (reportStepUpBlocked(error)) {
      creatingBackup.value = false
      return
    placeholder
    if (error?.response?.status === 409) {
      appStore.showWarning(t('admin.backup.operations.alreadyInProgress'))
    placeholder else {
      appStore.showError(error?.message || t('errors.networkError'))
    placeholder
    creatingBackup.value = false
  placeholder
placeholder

async function downloadBackup(id: string) {
  try {
    const result = await backupStepUp.run(() => adminAPI.backup.getDownloadURL(id))
    // 预签名 URL 带 attachment disposition，同页 anchor 导航直接触发下载；
    // 不用 window.open：step-up 弹窗 await 会耗尽瞬态用户激活，新标签页会被浏览器拦截。
    const link = document.createElement('a')
    link.href = result.url
    link.rel = 'noopener'
    link.click()
  placeholder catch (error) {
    if (isStepUpCancelled(error)) return
    if (reportStepUpBlocked(error)) return
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder
placeholder

async function restoreBackup(id: string) {
  if (!window.confirm(t('admin.backup.actions.restoreConfirm'))) return
  const password = window.prompt(t('admin.backup.actions.restorePasswordPrompt'))
  if (!password) return
  restoringId.value = id
  try {
    const record = await adminAPI.backup.restoreBackup(id, password)
    updateRecordInList(record)
    startRestorePolling(id)
  placeholder catch (error: any) {
    if (error?.response?.status === 409) {
      appStore.showWarning(t('admin.backup.operations.restoreRunning'))
    placeholder else {
      appStore.showError(error?.message || t('errors.networkError'))
    placeholder
    restoringId.value = ''
  placeholder
placeholder

async function removeBackup(id: string) {
  if (!window.confirm(t('admin.backup.actions.deleteConfirm'))) return
  try {
    await adminAPI.backup.deleteBackup(id)
    appStore.showSuccess(t('admin.backup.actions.deleted'))
    await loadBackups()
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder)?.message || t('errors.networkError'))
  placeholder
placeholder

function statusClass(status: string): string {
  switch (status) {
    case 'completed':
      return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300'
    case 'running':
      return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
    case 'failed':
      return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300'
    default:
      return 'bg-gray-100 text-gray-700 dark:bg-dark-800 dark:text-gray-300'
  placeholder
placeholder

function formatSize(bytes: number): string {
  if (!bytes || bytes <= 0) return '-'
  if (bytes < 1024) return `${bytesplaceholder B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)placeholder KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)placeholder MB`
placeholder

function formatDate(value?: string): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString()
placeholder

onMounted(async () => {
  document.addEventListener('visibilitychange', handleVisibilityChange)
  await Promise.all([loadS3Config(), loadSchedule(), loadBackups()])

  // 如果有正在 running 的备份，恢复轮询
  const runningBackup = backups.value.find(r => r.status === 'running')
  if (runningBackup) {
    creatingBackup.value = true
    startPolling(runningBackup.id)
  placeholder
  const restoringBackup = backups.value.find(r => r.restore_status === 'running')
  if (restoringBackup) {
    restoringId.value = restoringBackup.id
    startRestorePolling(restoringBackup.id)
  placeholder
placeholder)

onBeforeUnmount(() => {
  stopPolling()
  stopRestorePolling()
  document.removeEventListener('visibilitychange', handleVisibilityChange)
placeholder)
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
placeholder
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
placeholder
</style>
