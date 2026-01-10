<template>
  <AppLayout>
    <TablePageLayout>
      <template #actions>
        <div class="flex justify-end gap-3">
          <button
            @click="loadCodes"
            :disabled="loading"
            class="btn btn-secondary"
            :title="t('common.refresh')"
          >
            <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
          </button>
          <button @click="showCreateDialog = true" class="btn btn-primary">
            <Icon name="plus" size="md" class="mr-1" />
            {{ t('admin.promo.createCode') placeholderplaceholder
          </button>
        </div>
      </template>

      <template #filters>
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div class="max-w-md flex-1">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('admin.promo.searchCodes')"
              class="input"
              @input="handleSearch"
            />
          </div>
          <div class="flex gap-2">
            <Select
              v-model="filters.status"
              :options="filterStatusOptions"
              class="w-36"
              @change="loadCodes"
            />
          </div>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="codes" :loading="loading">
          <template #cell-code="{ value placeholder">
            <div class="flex items-center space-x-2">
              <code class="font-mono text-sm text-gray-900 dark:text-gray-100">{{ value placeholderplaceholder</code>
              <button
                @click="copyToClipboard(value)"
                :class="[
                  'flex items-center transition-colors',
                  copiedCode === value
                    ? 'text-green-500'
                    : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'
                ]"
                :title="copiedCode === value ? t('admin.promo.copied') : t('keys.copyToClipboard')"
              >
                <Icon v-if="copiedCode !== value" name="copy" size="sm" :stroke-width="2" />
                <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </button>
            </div>
          </template>

          <template #cell-bonus_amount="{ value placeholder">
            <span class="text-sm font-medium text-gray-900 dark:text-white">
              ${{ value.toFixed(2) placeholderplaceholder
            </span>
          </template>

          <template #cell-usage="{ row placeholder">
            <span class="text-sm text-gray-600 dark:text-gray-300">
              {{ row.used_count placeholderplaceholder / {{ row.max_uses === 0 ? '∞' : row.max_uses placeholderplaceholder
            </span>
          </template>

          <template #cell-status="{ value, row placeholder">
            <span
              :class="[
                'badge',
                getStatusClass(value, row)
              ]"
            >
              {{ getStatusLabel(value, row) placeholderplaceholder
            </span>
          </template>

          <template #cell-expires_at="{ value placeholder">
            <span class="text-sm text-gray-500 dark:text-dark-400">
              {{ value ? formatDateTime(value) : t('admin.promo.neverExpires') placeholderplaceholder
            </span>
          </template>

          <template #cell-created_at="{ value placeholder">
            <span class="text-sm text-gray-500 dark:text-dark-400">
              {{ formatDateTime(value) placeholderplaceholder
            </span>
          </template>

          <template #cell-actions="{ row placeholder">
            <div class="flex items-center space-x-1">
              <button
                @click="copyRegisterLink(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-green-50 hover:text-green-600 dark:hover:bg-green-900/20 dark:hover:text-green-400"
                :title="t('admin.promo.copyRegisterLink')"
              >
                <Icon name="link" size="sm" />
              </button>
              <button
                @click="handleViewUsages(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-blue-900/20 dark:hover:text-blue-400"
                :title="t('admin.promo.viewUsages')"
              >
                <Icon name="eye" size="sm" />
              </button>
              <button
                @click="handleEdit(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:hover:bg-dark-600 dark:hover:text-gray-300"
                :title="t('common.edit')"
              >
                <Icon name="edit" size="sm" />
              </button>
              <button
                @click="handleDelete(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
                :title="t('common.delete')"
              >
                <Icon name="trash" size="sm" />
              </button>
            </div>
          </template>
        </DataTable>
      </template>

      <template #pagination>
        <Pagination
          v-if="pagination.total > 0"
          :page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.page_size"
          @update:page="handlePageChange"
          @update:pageSize="handlePageSizeChange"
        />
      </template>
    </TablePageLayout>

    <!-- Create Dialog -->
    <BaseDialog
      :show="showCreateDialog"
      :title="t('admin.promo.createCode')"
      width="normal"
      @close="showCreateDialog = false"
    >
      <form id="create-promo-form" @submit.prevent="handleCreate" class="space-y-4">
        <div>
          <label class="input-label">
            {{ t('admin.promo.code') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('admin.promo.autoGenerate') placeholderplaceholder)</span>
          </label>
          <input
            v-model="createForm.code"
            type="text"
            class="input font-mono uppercase"
            :placeholder="t('admin.promo.codePlaceholder')"
          />
        </div>
        <div>
          <label class="input-label">{{ t('admin.promo.bonusAmount') placeholderplaceholder</label>
          <input
            v-model.number="createForm.bonus_amount"
            type="number"
            step="0.01"
            min="0"
            required
            class="input"
          />
        </div>
        <div>
          <label class="input-label">
            {{ t('admin.promo.maxUses') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('admin.promo.zeroUnlimited') placeholderplaceholder)</span>
          </label>
          <input
            v-model.number="createForm.max_uses"
            type="number"
            min="0"
            class="input"
          />
        </div>
        <div>
          <label class="input-label">
            {{ t('admin.promo.expiresAt') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('common.optional') placeholderplaceholder)</span>
          </label>
          <input
            v-model="createForm.expires_at_str"
            type="datetime-local"
            class="input"
          />
        </div>
        <div>
          <label class="input-label">
            {{ t('admin.promo.notes') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('common.optional') placeholderplaceholder)</span>
          </label>
          <textarea
            v-model="createForm.notes"
            rows="2"
            class="input"
            :placeholder="t('admin.promo.notesPlaceholder')"
          ></textarea>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button type="button" @click="showCreateDialog = false" class="btn btn-secondary">
            {{ t('common.cancel') placeholderplaceholder
          </button>
          <button type="submit" form="create-promo-form" :disabled="creating" class="btn btn-primary">
            {{ creating ? t('common.creating') : t('common.create') placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Edit Dialog -->
    <BaseDialog
      :show="showEditDialog"
      :title="t('admin.promo.editCode')"
      width="normal"
      @close="closeEditDialog"
    >
      <form id="edit-promo-form" @submit.prevent="handleUpdate" class="space-y-4">
        <div>
          <label class="input-label">{{ t('admin.promo.code') placeholderplaceholder</label>
          <input
            v-model="editForm.code"
            type="text"
            class="input font-mono uppercase"
          />
        </div>
        <div>
          <label class="input-label">{{ t('admin.promo.bonusAmount') placeholderplaceholder</label>
          <input
            v-model.number="editForm.bonus_amount"
            type="number"
            step="0.01"
            min="0"
            required
            class="input"
          />
        </div>
        <div>
          <label class="input-label">
            {{ t('admin.promo.maxUses') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('admin.promo.zeroUnlimited') placeholderplaceholder)</span>
          </label>
          <input
            v-model.number="editForm.max_uses"
            type="number"
            min="0"
            class="input"
          />
        </div>
        <div>
          <label class="input-label">{{ t('admin.promo.status') placeholderplaceholder</label>
          <Select v-model="editForm.status" :options="statusOptions" />
        </div>
        <div>
          <label class="input-label">
            {{ t('admin.promo.expiresAt') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('common.optional') placeholderplaceholder)</span>
          </label>
          <input
            v-model="editForm.expires_at_str"
            type="datetime-local"
            class="input"
          />
        </div>
        <div>
          <label class="input-label">
            {{ t('admin.promo.notes') placeholderplaceholder
            <span class="ml-1 text-xs font-normal text-gray-400">({{ t('common.optional') placeholderplaceholder)</span>
          </label>
          <textarea
            v-model="editForm.notes"
            rows="2"
            class="input"
          ></textarea>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button type="button" @click="closeEditDialog" class="btn btn-secondary">
            {{ t('common.cancel') placeholderplaceholder
          </button>
          <button type="submit" form="edit-promo-form" :disabled="updating" class="btn btn-primary">
            {{ updating ? t('common.saving') : t('common.save') placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Usages Dialog -->
    <BaseDialog
      :show="showUsagesDialog"
      :title="t('admin.promo.usageRecords')"
      width="wide"
      @close="showUsagesDialog = false"
    >
      <div v-if="usagesLoading" class="flex items-center justify-center py-8">
        <Icon name="refresh" size="lg" class="animate-spin text-gray-400" />
      </div>
      <div v-else-if="usages.length === 0" class="py-8 text-center text-gray-500 dark:text-gray-400">
        {{ t('admin.promo.noUsages') placeholderplaceholder
      </div>
      <div v-else class="space-y-3">
        <div
          v-for="usage in usages"
          :key="usage.id"
          class="flex items-center justify-between rounded-lg border border-gray-200 p-3 dark:border-dark-600"
        >
          <div class="flex items-center gap-3">
            <div class="flex h-8 w-8 items-center justify-center rounded-full bg-green-100 dark:bg-green-900/30">
              <Icon name="user" size="sm" class="text-green-600 dark:text-green-400" />
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900 dark:text-white">
                {{ usage.user?.email || t('admin.promo.userPrefix', { id: usage.user_id placeholder) placeholderplaceholder
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400">
                {{ formatDateTime(usage.used_at) placeholderplaceholder
              </p>
            </div>
          </div>
          <div class="text-right">
            <span class="text-sm font-medium text-green-600 dark:text-green-400">
              +${{ usage.bonus_amount.toFixed(2) placeholderplaceholder
            </span>
          </div>
        </div>
        <!-- Usages Pagination -->
        <div v-if="usagesTotal > usagesPageSize" class="mt-4">
          <Pagination
            :page="usagesPage"
            :total="usagesTotal"
            :page-size="usagesPageSize"
            :page-size-options="[10, 20, 50]"
            @update:page="handleUsagesPageChange"
            @update:page-size="(size: number) => { usagesPageSize = size; usagesPage = 1; loadUsages() placeholder"
          />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <button type="button" @click="showUsagesDialog = false" class="btn btn-secondary">
            {{ t('common.close') placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      :show="showDeleteDialog"
      :title="t('admin.promo.deleteCode')"
      :message="t('admin.promo.deleteCodeConfirm')"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      danger
      @confirm="confirmDelete"
      @cancel="showDeleteDialog = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { useClipboard placeholder from '@/composables/useClipboard'
import { adminAPI placeholder from '@/api/admin'
import { formatDateTime placeholder from '@/utils/format'
import type { PromoCode, PromoCodeUsage placeholder from '@/types'
import type { Column placeholder from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'

const { t placeholder = useI18n()
const appStore = useAppStore()
const { copyToClipboard: clipboardCopy placeholder = useClipboard()

// State
const codes = ref<PromoCode[]>([])
const loading = ref(false)
const creating = ref(false)
const updating = ref(false)
const searchQuery = ref('')
const copiedCode = ref<string | null>(null)

const filters = reactive({
  status: ''
placeholder)

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
placeholder)

// Dialogs
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const showUsagesDialog = ref(false)

const editingCode = ref<PromoCode | null>(null)
const deletingCode = ref<PromoCode | null>(null)

// Usages
const usages = ref<PromoCodeUsage[]>([])
const usagesLoading = ref(false)
const currentViewingCode = ref<PromoCode | null>(null)
const usagesPage = ref(1)
const usagesPageSize = ref(20)
const usagesTotal = ref(0)

// Forms
const createForm = reactive({
  code: '',
  bonus_amount: 1,
  max_uses: 0,
  expires_at_str: '',
  notes: ''
placeholder)

const editForm = reactive({
  code: '',
  bonus_amount: 0,
  max_uses: 0,
  status: 'active' as 'active' | 'disabled',
  expires_at_str: '',
  notes: ''
placeholder)

// Options
const filterStatusOptions = computed(() => [
  { value: '', label: t('admin.promo.allStatus') placeholder,
  { value: 'active', label: t('admin.promo.statusActive') placeholder,
  { value: 'disabled', label: t('admin.promo.statusDisabled') placeholder
])

const statusOptions = computed(() => [
  { value: 'active', label: t('admin.promo.statusActive') placeholder,
  { value: 'disabled', label: t('admin.promo.statusDisabled') placeholder
])

const columns = computed<Column[]>(() => [
  { key: 'code', label: t('admin.promo.columns.code') placeholder,
  { key: 'bonus_amount', label: t('admin.promo.columns.bonusAmount'), sortable: true placeholder,
  { key: 'usage', label: t('admin.promo.columns.usage') placeholder,
  { key: 'status', label: t('admin.promo.columns.status'), sortable: true placeholder,
  { key: 'expires_at', label: t('admin.promo.columns.expiresAt'), sortable: true placeholder,
  { key: 'created_at', label: t('admin.promo.columns.createdAt'), sortable: true placeholder,
  { key: 'actions', label: t('admin.promo.columns.actions') placeholder
])

// Helpers
const getStatusClass = (status: string, row: PromoCode) => {
  if (row.expires_at && new Date(row.expires_at) < new Date()) {
    return 'badge-danger'
  placeholder
  if (row.max_uses > 0 && row.used_count >= row.max_uses) {
    return 'badge-gray'
  placeholder
  return status === 'active' ? 'badge-success' : 'badge-gray'
placeholder

const getStatusLabel = (status: string, row: PromoCode) => {
  if (row.expires_at && new Date(row.expires_at) < new Date()) {
    return t('admin.promo.statusExpired')
  placeholder
  if (row.max_uses > 0 && row.used_count >= row.max_uses) {
    return t('admin.promo.statusMaxUsed')
  placeholder
  return status === 'active' ? t('admin.promo.statusActive') : t('admin.promo.statusDisabled')
placeholder

// API calls
let abortController: AbortController | null = null

const loadCodes = async () => {
  if (abortController) {
    abortController.abort()
  placeholder
  const currentController = new AbortController()
  abortController = currentController
  loading.value = true

  try {
    const response = await adminAPI.promo.list(
      pagination.page,
      pagination.page_size,
      {
        status: filters.status || undefined,
        search: searchQuery.value || undefined
      placeholder
    )
    if (currentController.signal.aborted) return

    codes.value = response.items
    pagination.total = response.total
  placeholder catch (error: any) {
    if (currentController.signal.aborted || error?.name === 'AbortError') return
    appStore.showError(t('admin.promo.failedToLoad'))
    console.error('Error loading promo codes:', error)
  placeholder finally {
    if (abortController === currentController && !currentController.signal.aborted) {
      loading.value = false
      abortController = null
    placeholder
  placeholder
placeholder

let searchTimeout: ReturnType<typeof setTimeout>
const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.page = 1
    loadCodes()
  placeholder, 300)
placeholder

const handlePageChange = (page: number) => {
  pagination.page = page
  loadCodes()
placeholder

const handlePageSizeChange = (pageSize: number) => {
  pagination.page_size = pageSize
  pagination.page = 1
  loadCodes()
placeholder

const copyToClipboard = async (text: string) => {
  const success = await clipboardCopy(text, t('admin.promo.copied'))
  if (success) {
    copiedCode.value = text
    setTimeout(() => {
      copiedCode.value = null
    placeholder, 2000)
  placeholder
placeholder

// Create
const handleCreate = async () => {
  creating.value = true
  try {
    await adminAPI.promo.create({
      code: createForm.code || undefined,
      bonus_amount: createForm.bonus_amount,
      max_uses: createForm.max_uses,
      expires_at: createForm.expires_at_str ? Math.floor(new Date(createForm.expires_at_str).getTime() / 1000) : undefined,
      notes: createForm.notes || undefined
    placeholder)
    appStore.showSuccess(t('admin.promo.codeCreated'))
    showCreateDialog.value = false
    resetCreateForm()
    loadCodes()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.promo.failedToCreate'))
  placeholder finally {
    creating.value = false
  placeholder
placeholder

const resetCreateForm = () => {
  createForm.code = ''
  createForm.bonus_amount = 1
  createForm.max_uses = 0
  createForm.expires_at_str = ''
  createForm.notes = ''
placeholder

// Edit
const handleEdit = (code: PromoCode) => {
  editingCode.value = code
  editForm.code = code.code
  editForm.bonus_amount = code.bonus_amount
  editForm.max_uses = code.max_uses
  editForm.status = code.status
  editForm.expires_at_str = code.expires_at ? new Date(code.expires_at).toISOString().slice(0, 16) : ''
  editForm.notes = code.notes || ''
  showEditDialog.value = true
placeholder

const closeEditDialog = () => {
  showEditDialog.value = false
  editingCode.value = null
placeholder

const handleUpdate = async () => {
  if (!editingCode.value) return

  updating.value = true
  try {
    await adminAPI.promo.update(editingCode.value.id, {
      code: editForm.code,
      bonus_amount: editForm.bonus_amount,
      max_uses: editForm.max_uses,
      status: editForm.status,
      expires_at: editForm.expires_at_str ? Math.floor(new Date(editForm.expires_at_str).getTime() / 1000) : 0,
      notes: editForm.notes
    placeholder)
    appStore.showSuccess(t('admin.promo.codeUpdated'))
    closeEditDialog()
    loadCodes()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.promo.failedToUpdate'))
  placeholder finally {
    updating.value = false
  placeholder
placeholder

// Copy Register Link
const copyRegisterLink = async (code: PromoCode) => {
  const baseUrl = window.location.origin
  const registerLink = `${baseUrlplaceholder/register?promo=${encodeURIComponent(code.code)placeholder`

  try {
    await navigator.clipboard.writeText(registerLink)
    appStore.showSuccess(t('admin.promo.registerLinkCopied'))
  placeholder catch (error) {
    // Fallback for older browsers
    const textArea = document.createElement('textarea')
    textArea.value = registerLink
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    appStore.showSuccess(t('admin.promo.registerLinkCopied'))
  placeholder
placeholder

// Delete
const handleDelete = (code: PromoCode) => {
  deletingCode.value = code
  showDeleteDialog.value = true
placeholder

const confirmDelete = async () => {
  if (!deletingCode.value) return

  try {
    await adminAPI.promo.delete(deletingCode.value.id)
    appStore.showSuccess(t('admin.promo.codeDeleted'))
    showDeleteDialog.value = false
    deletingCode.value = null
    loadCodes()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.promo.failedToDelete'))
  placeholder
placeholder

// View Usages
const handleViewUsages = async (code: PromoCode) => {
  currentViewingCode.value = code
  showUsagesDialog.value = true
  usagesPage.value = 1
  await loadUsages()
placeholder

const loadUsages = async () => {
  if (!currentViewingCode.value) return
  usagesLoading.value = true
  usages.value = []

  try {
    const response = await adminAPI.promo.getUsages(
      currentViewingCode.value.id,
      usagesPage.value,
      usagesPageSize.value
    )
    usages.value = response.items
    usagesTotal.value = response.total
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.promo.failedToLoadUsages'))
  placeholder finally {
    usagesLoading.value = false
  placeholder
placeholder

const handleUsagesPageChange = (page: number) => {
  usagesPage.value = page
  loadUsages()
placeholder

onMounted(() => {
  loadCodes()
placeholder)

onUnmounted(() => {
  clearTimeout(searchTimeout)
  abortController?.abort()
placeholder)
</script>
