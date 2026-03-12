<template>
  <BaseDialog :show="show" :title="t('admin.groups.rateMultipliersTitle')" width="normal" @close="$emit('close')">
    <div v-if="group" class="space-y-5">
      <!-- 分组信息 -->
      <div class="flex flex-wrap items-center gap-3 rounded-lg bg-gray-50 px-4 py-3 text-sm dark:bg-dark-700">
        <span class="font-medium text-gray-900 dark:text-white">{{ group.name placeholderplaceholder</span>
        <span class="text-gray-400">|</span>
        <span class="text-gray-600 dark:text-gray-400">{{ t('admin.groups.platforms.' + group.platform) placeholderplaceholder</span>
        <span class="text-gray-400">|</span>
        <span class="text-gray-600 dark:text-gray-400">
          {{ t('admin.groups.columns.rateMultiplier') placeholderplaceholder: {{ group.rate_multiplier placeholderplaceholderx
        </span>
      </div>

      <!-- 添加用户 -->
      <div class="rounded-lg border border-gray-200 p-4 dark:border-dark-600">
        <h4 class="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('admin.groups.addUserRate') placeholderplaceholder
        </h4>
        <div class="flex items-end gap-3">
          <div class="relative flex-1">
            <input
              v-model="searchQuery"
              type="text"
              class="input w-full"
              :placeholder="t('admin.groups.searchUserPlaceholder')"
              @input="handleSearchUsers"
              @focus="showDropdown = true"
            />
            <!-- 搜索结果下拉 -->
            <div
              v-if="showDropdown && searchResults.length > 0"
              class="absolute left-0 right-0 top-full z-10 mt-1 max-h-48 overflow-y-auto rounded-lg border border-gray-200 bg-white shadow-lg dark:border-dark-500 dark:bg-dark-700"
            >
              <button
                v-for="user in searchResults"
                :key="user.id"
                type="button"
                class="flex w-full items-center px-3 py-2 text-left text-sm hover:bg-gray-50 dark:hover:bg-dark-600"
                @click="selectUser(user)"
              >
                <span class="text-gray-900 dark:text-white">{{ user.email placeholderplaceholder</span>
              </button>
            </div>
          </div>
          <div class="w-28">
            <input
              v-model.number="newRate"
              type="number"
              step="0.001"
              min="0"
              class="hide-spinner input w-full"
              placeholder="1.0"
            />
          </div>
          <button
            type="button"
            class="btn btn-primary shrink-0"
            :disabled="!selectedUser || !newRate || addingRate"
            @click="handleAddRate"
          >
            <Icon v-if="addingRate" name="refresh" size="sm" class="mr-1 animate-spin" />
            {{ t('common.add') placeholderplaceholder
          </button>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="flex justify-center py-8">
        <svg class="h-8 w-8 animate-spin text-primary-500" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- 已设置的用户列表 -->
      <div v-else>
        <h4 class="mb-3 text-sm font-medium text-gray-700 dark:text-gray-300">
          {{ t('admin.groups.rateMultipliers') placeholderplaceholder ({{ entries.length placeholderplaceholder)
        </h4>

        <div v-if="entries.length === 0" class="py-8 text-center text-sm text-gray-400 dark:text-gray-500">
          {{ t('admin.groups.noRateMultipliers') placeholderplaceholder
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="entry in entries"
            :key="entry.user_id"
            class="flex items-center gap-3 rounded-lg border border-gray-200 px-4 py-3 dark:border-dark-600"
          >
            <span class="flex-1 text-sm text-gray-900 dark:text-white">{{ entry.user_email placeholderplaceholder</span>
            <input
              type="number"
              step="0.001"
              min="0"
              :value="entry.rate_multiplier"
              class="hide-spinner w-24 rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm font-medium transition-colors focus:border-primary-500 focus:outline-none focus:ring-2 focus:ring-primary-500/20 dark:border-dark-500 dark:bg-dark-700 dark:focus:border-primary-500"
              @blur="handleUpdateRate(entry, ($event.target as HTMLInputElement).value)"
              @keydown.enter="($event.target as HTMLInputElement).blur()"
            />
            <button
              type="button"
              class="rounded-lg p-1.5 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
              @click="handleDeleteRate(entry)"
            >
              <Icon name="trash" size="sm" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import type { GroupRateMultiplierEntry placeholder from '@/api/admin/groups'
import type { AdminGroup, AdminUser placeholder from '@/types'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'

const props = defineProps<{
  show: boolean
  group: AdminGroup | null
placeholder>()

const emit = defineEmits<{
  close: []
  success: []
placeholder>()

const { t placeholder = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const entries = ref<GroupRateMultiplierEntry[]>([])
const searchQuery = ref('')
const searchResults = ref<AdminUser[]>([])
const showDropdown = ref(false)
const selectedUser = ref<AdminUser | null>(null)
const newRate = ref<number | null>(null)
const addingRate = ref(false)

let searchTimeout: ReturnType<typeof setTimeout>

const loadEntries = async () => {
  if (!props.group) return
  loading.value = true
  try {
    entries.value = await adminAPI.groups.getGroupRateMultipliers(props.group.id)
  placeholder catch (error) {
    appStore.showError(t('admin.groups.failedToLoad'))
    console.error('Error loading group rate multipliers:', error)
  placeholder finally {
    loading.value = false
  placeholder
placeholder

watch(() => props.show, (val) => {
  if (val && props.group) {
    loadEntries()
    searchQuery.value = ''
    searchResults.value = []
    selectedUser.value = null
    newRate.value = null
  placeholder
placeholder)

const handleSearchUsers = () => {
  clearTimeout(searchTimeout)
  selectedUser.value = null
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    showDropdown.value = false
    return
  placeholder
  searchTimeout = setTimeout(async () => {
    try {
      const res = await adminAPI.users.list(1, 10, { search: searchQuery.value.trim() placeholder)
      searchResults.value = res.items
      showDropdown.value = true
    placeholder catch {
      searchResults.value = []
    placeholder
  placeholder, 300)
placeholder

const selectUser = (user: AdminUser) => {
  selectedUser.value = user
  searchQuery.value = user.email
  showDropdown.value = false
  searchResults.value = []
placeholder

const handleAddRate = async () => {
  if (!selectedUser.value || !newRate.value || !props.group) return
  addingRate.value = true
  try {
    await adminAPI.users.update(selectedUser.value.id, {
      group_rates: { [props.group.id]: newRate.value placeholder
    placeholder)
    appStore.showSuccess(t('admin.groups.rateAdded'))
    searchQuery.value = ''
    selectedUser.value = null
    newRate.value = null
    await loadEntries()
    emit('success')
  placeholder catch (error) {
    appStore.showError(t('admin.groups.failedToSave'))
    console.error('Error adding rate multiplier:', error)
  placeholder finally {
    addingRate.value = false
  placeholder
placeholder

const handleUpdateRate = async (entry: GroupRateMultiplierEntry, value: string) => {
  if (!props.group) return
  const numValue = parseFloat(value)
  if (isNaN(numValue) || numValue === entry.rate_multiplier) return
  try {
    await adminAPI.users.update(entry.user_id, {
      group_rates: { [props.group.id]: numValue placeholder
    placeholder)
    appStore.showSuccess(t('admin.groups.rateUpdated'))
    await loadEntries()
    emit('success')
  placeholder catch (error) {
    appStore.showError(t('admin.groups.failedToSave'))
    console.error('Error updating rate multiplier:', error)
  placeholder
placeholder

const handleDeleteRate = async (entry: GroupRateMultiplierEntry) => {
  if (!props.group) return
  try {
    await adminAPI.users.update(entry.user_id, {
      group_rates: { [props.group.id]: null placeholder
    placeholder)
    appStore.showSuccess(t('admin.groups.rateDeleted'))
    await loadEntries()
    emit('success')
  placeholder catch (error) {
    appStore.showError(t('admin.groups.failedToSave'))
    console.error('Error deleting rate multiplier:', error)
  placeholder
placeholder

// 点击外部关闭下拉
const handleClickOutside = () => {
  showDropdown.value = false
placeholder

if (typeof document !== 'undefined') {
  document.addEventListener('click', handleClickOutside)
placeholder
</script>

<style scoped>
.hide-spinner::-webkit-outer-spin-button,
.hide-spinner::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
placeholder
.hide-spinner {
  -moz-appearance: textfield;
placeholder
</style>
