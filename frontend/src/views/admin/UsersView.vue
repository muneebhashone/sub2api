<template>
  <AppLayout>
    <TablePageLayout>
      <!-- Single Row: Search, Filters, and Actions -->
      <template #filters>
        <div class="flex w-full flex-wrap-reverse items-center justify-between gap-4">
          <!-- Left: Search + Active Filters -->
          <div class="flex min-w-[280px] flex-1 flex-wrap content-start items-center gap-3">
            <!-- Search Box -->
            <div class="relative w-full sm:w-64">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
              />
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('admin.users.searchUsers')"
                class="input pl-10"
                @input="handleSearch"
              />
            </div>

            <!-- Role Filter (visible when enabled) -->
            <div v-if="visibleFilters.has('role')" class="w-full sm:w-32">
              <Select
                v-model="filters.role"
                :options="[
                  { value: '', label: t('admin.users.allRoles') placeholder,
                  { value: 'admin', label: t('admin.users.admin') placeholder,
                  { value: 'user', label: t('admin.users.user') placeholder
                ]"
                @change="applyFilter"
              />
            </div>

            <!-- Status Filter (visible when enabled) -->
            <div v-if="visibleFilters.has('status')" class="w-full sm:w-32">
              <Select
                v-model="filters.status"
                :options="[
                  { value: '', label: t('admin.users.allStatus') placeholder,
                  { value: 'active', label: t('common.active') placeholder,
                  { value: 'disabled', label: t('admin.users.disabled') placeholder
                ]"
                @change="applyFilter"
              />
            </div>

            <!-- Dynamic Attribute Filters -->
            <template v-for="(value, attrId) in activeAttributeFilters" :key="attrId">
              <div
                v-if="visibleFilters.has(`attr_${attrIdplaceholder`)"
                class="relative w-full sm:w-36"
              >
                <!-- Text/Email/URL/Textarea/Date type: styled input -->
                <input
                  v-if="['text', 'textarea', 'email', 'url', 'date'].includes(getAttributeDefinition(Number(attrId))?.type || 'text')"
                  :value="value"
                  @input="(e) => updateAttributeFilter(Number(attrId), (e.target as HTMLInputElement).value)"
                  @keyup.enter="applyFilter"
                  :placeholder="getAttributeDefinitionName(Number(attrId))"
                  class="input w-full"
                />
                <!-- Number type: number input -->
                <input
                  v-else-if="getAttributeDefinition(Number(attrId))?.type === 'number'"
                  :value="value"
                  type="number"
                  @input="(e) => updateAttributeFilter(Number(attrId), (e.target as HTMLInputElement).value)"
                  @keyup.enter="applyFilter"
                  :placeholder="getAttributeDefinitionName(Number(attrId))"
                  class="input w-full"
                />
                <!-- Select/Multi-select type -->
                <template v-else-if="['select', 'multi_select'].includes(getAttributeDefinition(Number(attrId))?.type || '')">
                  <div class="w-full">
                    <Select
                      :model-value="value"
                      :options="[
                        { value: '', label: getAttributeDefinitionName(Number(attrId)) placeholder,
                        ...(getAttributeDefinition(Number(attrId))?.options || [])
                      ]"
                      @update:model-value="(val) => { updateAttributeFilter(Number(attrId), String(val ?? '')); applyFilter() placeholder"
                    />
                  </div>
                </template>
                <!-- Fallback -->
                <input
                  v-else
                  :value="value"
                  @input="(e) => updateAttributeFilter(Number(attrId), (e.target as HTMLInputElement).value)"
                  @keyup.enter="applyFilter"
                  :placeholder="getAttributeDefinitionName(Number(attrId))"
                  class="input w-full"
                />
              </div>
            </template>
          </div>

          <!-- Right: Actions and Settings -->
          <div class="ml-auto flex max-w-full flex-wrap items-center justify-end gap-3">
            <!-- Refresh Button -->
            <button
              @click="loadUsers"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <!-- Filter Settings Dropdown -->
            <div class="relative" ref="filterDropdownRef">
              <button
                @click="showFilterDropdown = !showFilterDropdown"
                class="btn btn-secondary"
              >
                <Icon name="filter" size="sm" class="mr-1.5" />
                {{ t('admin.users.filterSettings') placeholderplaceholder
              </button>
              <!-- Dropdown menu -->
              <div
                v-if="showFilterDropdown"
                class="absolute right-0 top-full z-50 mt-1 w-48 rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800"
              >
                <!-- Built-in filters -->
                <button
                  v-for="filter in builtInFilters"
                  :key="filter.key"
                  @click="toggleBuiltInFilter(filter.key)"
                  class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
                >
                  <span>{{ filter.name placeholderplaceholder</span>
                  <Icon
                    v-if="visibleFilters.has(filter.key)"
                    name="check"
                    size="sm"
                    class="text-primary-500"
                    :stroke-width="2"
                  />
                </button>
                <!-- Divider if custom attributes exist -->
                <div
                  v-if="filterableAttributes.length > 0"
                  class="my-1 border-t border-gray-100 dark:border-dark-700"
                ></div>
                <!-- Custom attribute filters -->
                <button
                  v-for="attr in filterableAttributes"
                  :key="attr.id"
                  @click="toggleAttributeFilter(attr)"
                  class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
                >
                  <span>{{ attr.name placeholderplaceholder</span>
                  <Icon
                    v-if="visibleFilters.has(`attr_${attr.idplaceholder`)"
                    name="check"
                    size="sm"
                    class="text-primary-500"
                    :stroke-width="2"
                  />
                </button>
              </div>
            </div>
            <!-- Column Settings Dropdown -->
            <div class="relative" ref="columnDropdownRef">
              <button
                @click="showColumnDropdown = !showColumnDropdown"
                class="btn btn-secondary"
              >
                <svg class="mr-1.5 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="1.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 4.5v15m6-15v15m-10.875 0h15.75c.621 0 1.125-.504 1.125-1.125V5.625c0-.621-.504-1.125-1.125-1.125H4.125C3.504 4.5 3 5.004 3 5.625v12.75c0 .621.504 1.125 1.125 1.125z" />
                </svg>
                {{ t('admin.users.columnSettings') placeholderplaceholder
              </button>
              <!-- Dropdown menu -->
              <div
                v-if="showColumnDropdown"
                class="absolute right-0 top-full z-50 mt-1 max-h-80 w-48 overflow-y-auto rounded-lg border border-gray-200 bg-white py-1 shadow-lg dark:border-dark-600 dark:bg-dark-800"
              >
                <button
                  v-for="col in toggleableColumns"
                  :key="col.key"
                  @click="toggleColumn(col.key)"
                  class="flex w-full items-center justify-between px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
                >
                  <span>{{ col.label placeholderplaceholder</span>
                  <Icon
                    v-if="isColumnVisible(col.key)"
                    name="check"
                    size="sm"
                    class="text-primary-500"
                    :stroke-width="2"
                  />
                </button>
              </div>
            </div>
            <!-- Attributes Config Button -->
            <button @click="showAttributesModal = true" class="btn btn-secondary">
              <Icon name="cog" size="sm" class="mr-1.5" />
              {{ t('admin.users.attributes.configButton') placeholderplaceholder
            </button>
            <!-- Create User Button -->
            <button @click="showCreateModal = true" class="btn btn-primary">
              <Icon name="plus" size="md" class="mr-2" />
              {{ t('admin.users.createUser') placeholderplaceholder
            </button>
          </div>
        </div>
      </template>

      <!-- Users Table -->
      <template #table>
        <DataTable :columns="columns" :data="users" :loading="loading" :actions-count="7">
          <template #cell-email="{ value placeholder">
            <div class="flex items-center gap-2">
              <div
                class="flex h-8 w-8 items-center justify-center rounded-full bg-primary-100 dark:bg-primary-900/30"
              >
                <span class="text-sm font-medium text-primary-700 dark:text-primary-300">
                  {{ value.charAt(0).toUpperCase() placeholderplaceholder
                </span>
              </div>
              <span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
            </div>
          </template>

          <template #cell-username="{ value placeholder">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value || '-' placeholderplaceholder</span>
          </template>

          <template #cell-notes="{ value placeholder">
            <div class="max-w-xs">
              <span
                v-if="value"
                :title="value.length > 30 ? value : undefined"
                class="block truncate text-sm text-gray-600 dark:text-gray-400"
              >
                {{ value.length > 30 ? value.substring(0, 25) + '...' : value placeholderplaceholder
              </span>
              <span v-else class="text-sm text-gray-400">-</span>
            </div>
          </template>

          <!-- Dynamic attribute columns -->
          <template
            v-for="def in attributeDefinitions.filter(d => d.enabled)"
            :key="def.id"
            #[`cell-attr_${def.idplaceholder`]="{ row placeholder"
          >
            <div class="max-w-xs">
              <span
                class="block truncate text-sm text-gray-700 dark:text-gray-300"
                :title="getAttributeValue(row.id, def.id)"
              >
                {{ getAttributeValue(row.id, def.id) placeholderplaceholder
              </span>
            </div>
          </template>

          <template #cell-role="{ value placeholder">
            <span :class="['badge', value === 'admin' ? 'badge-purple' : 'badge-gray']">
              {{ t('admin.users.roles.' + value) placeholderplaceholder
            </span>
          </template>

          <template #cell-subscriptions="{ row placeholder">
            <div
              v-if="row.subscriptions && row.subscriptions.length > 0"
              class="flex flex-wrap gap-1.5"
            >
              <GroupBadge
                v-for="sub in row.subscriptions"
                :key="sub.id"
                :name="sub.group?.name || ''"
                :platform="sub.group?.platform"
                :subscription-type="sub.group?.subscription_type"
                :rate-multiplier="sub.group?.rate_multiplier"
                :days-remaining="sub.expires_at ? getDaysRemaining(sub.expires_at) : null"
                :title="sub.expires_at ? formatDateTime(sub.expires_at) : ''"
              />
            </div>
            <span
              v-else
              class="inline-flex items-center gap-1.5 rounded-md bg-gray-50 px-2 py-1 text-xs text-gray-400 dark:bg-dark-700/50 dark:text-dark-500"
            >
              <Icon name="ban" size="xs" class="h-3.5 w-3.5" />
              <span>{{ t('admin.users.noSubscription') placeholderplaceholder</span>
            </span>
          </template>

          <template #cell-balance="{ value placeholder">
            <span class="font-medium text-gray-900 dark:text-white">${{ value.toFixed(2) placeholderplaceholder</span>
          </template>

          <template #cell-usage="{ row placeholder">
            <div class="text-sm">
              <div class="flex items-center gap-1.5">
                <span class="text-gray-500 dark:text-gray-400">{{ t('admin.users.today') placeholderplaceholder:</span>
                <span class="font-medium text-gray-900 dark:text-white">
                  ${{ (usageStats[row.id]?.today_actual_cost ?? 0).toFixed(4) placeholderplaceholder
                </span>
              </div>
              <div class="mt-0.5 flex items-center gap-1.5">
                <span class="text-gray-500 dark:text-gray-400">{{ t('admin.users.total') placeholderplaceholder:</span>
                <span class="font-medium text-gray-900 dark:text-white">
                  ${{ (usageStats[row.id]?.total_actual_cost ?? 0).toFixed(4) placeholderplaceholder
                </span>
              </div>
            </div>
          </template>

          <template #cell-concurrency="{ value placeholder">
            <span class="text-sm text-gray-700 dark:text-gray-300">{{ value placeholderplaceholder</span>
          </template>

          <template #cell-status="{ value placeholder">
            <div class="flex items-center gap-1.5">
              <span
                :class="[
                  'inline-block h-2 w-2 rounded-full',
                  value === 'active' ? 'bg-green-500' : 'bg-red-500'
                ]"
              ></span>
              <span class="text-sm text-gray-700 dark:text-gray-300">
                {{ value === 'active' ? t('common.active') : t('admin.users.disabled') placeholderplaceholder
              </span>
            </div>
          </template>

          <template #cell-created_at="{ value placeholder">
            <span class="text-sm text-gray-500 dark:text-dark-400">{{ formatDateTime(value) placeholderplaceholder</span>
          </template>

          <template #cell-actions="{ row placeholder">
            <div class="flex items-center gap-1">
              <!-- Edit Button -->
              <button
                @click="handleEdit(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-primary-600 dark:hover:bg-dark-700 dark:hover:text-primary-400"
              >
                <Icon name="edit" size="sm" />
                <span class="text-xs">{{ t('common.edit') placeholderplaceholder</span>
              </button>

              <!-- Toggle Status Button (not for admin) -->
              <button
                v-if="row.role !== 'admin'"
                @click="handleToggleStatus(row)"
                :class="[
                  'flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors',
                  row.status === 'active'
                    ? 'hover:bg-orange-50 hover:text-orange-600 dark:hover:bg-orange-900/20 dark:hover:text-orange-400'
                    : 'hover:bg-green-50 hover:text-green-600 dark:hover:bg-green-900/20 dark:hover:text-green-400'
                ]"
              >
                <Icon v-if="row.status === 'active'" name="ban" size="sm" />
                <Icon v-else name="checkCircle" size="sm" />
                <span class="text-xs">{{ row.status === 'active' ? t('admin.users.disable') : t('admin.users.enable') placeholderplaceholder</span>
              </button>

              <!-- More Actions Menu Trigger -->
              <button
                :ref="(el) => setActionButtonRef(row.id, el)"
                @click="openActionMenu(row)"
                class="action-menu-trigger flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-dark-700 dark:hover:text-white"
                :class="{ 'bg-gray-100 text-gray-900 dark:bg-dark-700 dark:text-white': activeMenuId === row.id placeholder"
              >
                <Icon name="more" size="sm" />
                <span class="text-xs">{{ t('common.more') placeholderplaceholder</span>
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState
              :title="t('admin.users.noUsersYet')"
              :description="t('admin.users.createFirstUser')"
              :action-text="t('admin.users.createUser')"
              @action="showCreateModal = true"
            />
          </template>
        </DataTable>
      </template>

      <!-- Pagination -->
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

    <!-- Action Menu (Teleported) -->
    <Teleport to="body">
      <div
        v-if="activeMenuId !== null && menuPosition"
        class="action-menu-content fixed z-[9999] w-48 overflow-hidden rounded-xl bg-white shadow-lg ring-1 ring-black/5 dark:bg-dark-800 dark:ring-white/10"
        :style="{ top: menuPosition.top + 'px', left: menuPosition.left + 'px' placeholder"
      >
        <div class="py-1">
          <template v-for="user in users" :key="user.id">
            <template v-if="user.id === activeMenuId">
              <!-- View API Keys -->
              <button
                @click="handleViewApiKeys(user); closeActionMenu()"
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
              >
                <Icon name="key" size="sm" class="text-gray-400" :stroke-width="2" />
                {{ t('admin.users.apiKeys') placeholderplaceholder
              </button>

              <!-- Allowed Groups -->
              <button
                @click="handleAllowedGroups(user); closeActionMenu()"
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
              >
                <Icon name="users" size="sm" class="text-gray-400" :stroke-width="2" />
                {{ t('admin.users.groups') placeholderplaceholder
              </button>

              <div class="my-1 border-t border-gray-100 dark:border-dark-700"></div>

              <!-- Deposit -->
              <button
                @click="handleDeposit(user); closeActionMenu()"
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
              >
                <Icon name="plus" size="sm" class="text-emerald-500" :stroke-width="2" />
                {{ t('admin.users.deposit') placeholderplaceholder
              </button>

              <!-- Withdraw -->
              <button
                @click="handleWithdraw(user); closeActionMenu()"
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-dark-700"
              >
                <svg class="h-4 w-4 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
                </svg>
                {{ t('admin.users.withdraw') placeholderplaceholder
              </button>

              <div class="my-1 border-t border-gray-100 dark:border-dark-700"></div>

              <!-- Delete (not for admin) -->
              <button
                v-if="user.role !== 'admin'"
                @click="handleDelete(user); closeActionMenu()"
                class="flex w-full items-center gap-2 px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:text-red-400 dark:hover:bg-red-900/20"
              >
                <Icon name="trash" size="sm" :stroke-width="2" />
                {{ t('common.delete') placeholderplaceholder
              </button>
            </template>
          </template>
        </div>
      </div>
    </Teleport>

    <ConfirmDialog :show="showDeleteDialog" :title="t('admin.users.deleteUser')" :message="t('admin.users.deleteConfirm', { email: deletingUser?.email placeholder)" :danger="true" @confirm="confirmDelete" @cancel="showDeleteDialog = false" />
    <UserCreateModal :show="showCreateModal" @close="showCreateModal = false" @success="loadUsers" />
    <UserEditModal :show="showEditModal" :user="editingUser" @close="closeEditModal" @success="loadUsers" />
    <UserApiKeysModal :show="showApiKeysModal" :user="viewingUser" @close="closeApiKeysModal" />
    <UserAllowedGroupsModal :show="showAllowedGroupsModal" :user="allowedGroupsUser" @close="closeAllowedGroupsModal" @success="loadUsers" />
    <UserBalanceModal :show="showBalanceModal" :user="balanceUser" :operation="balanceOperation" @close="closeBalanceModal" @success="loadUsers" />
    <UserAttributesConfigModal :show="showAttributesModal" @close="handleAttributesModalClose" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, type ComponentPublicInstance placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { formatDateTime placeholder from '@/utils/format'
import Icon from '@/components/icons/Icon.vue'

const { t placeholder = useI18n()
import { adminAPI placeholder from '@/api/admin'
import type { User, UserAttributeDefinition placeholder from '@/types'
import type { BatchUserUsageStats placeholder from '@/api/admin/dashboard'
import type { Column placeholder from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import Select from '@/components/common/Select.vue'
import UserAttributesConfigModal from '@/components/user/UserAttributesConfigModal.vue'
import UserCreateModal from '@/components/admin/user/UserCreateModal.vue'
import UserEditModal from '@/components/admin/user/UserEditModal.vue'
import UserApiKeysModal from '@/components/admin/user/UserApiKeysModal.vue'
import UserAllowedGroupsModal from '@/components/admin/user/UserAllowedGroupsModal.vue'
import UserBalanceModal from '@/components/admin/user/UserBalanceModal.vue'

const appStore = useAppStore()

// Generate dynamic attribute columns from enabled definitions
const attributeColumns = computed<Column[]>(() =>
  attributeDefinitions.value
    .filter(def => def.enabled)
    .map(def => ({
      key: `attr_${def.idplaceholder`,
      label: def.name,
      sortable: false
    placeholder))
)

// Get formatted attribute value for display in table
const getAttributeValue = (userId: number, attrId: number): string => {
  const userAttrs = userAttributeValues.value[userId]
  if (!userAttrs) return '-'
  const value = userAttrs[attrId]
  if (!value) return '-'

  // Find definition for this attribute
  const def = attributeDefinitions.value.find(d => d.id === attrId)
  if (!def) return value

  // Format based on type
  if (def.type === 'multi_select' && value) {
    try {
      const arr = JSON.parse(value)
      if (Array.isArray(arr)) {
        // Map values to labels
        return arr.map(v => {
          const opt = def.options?.find(o => o.value === v)
          return opt?.label || v
        placeholder).join(', ')
      placeholder
    placeholder catch {
      return value
    placeholder
  placeholder

  if (def.type === 'select' && value && def.options) {
    const opt = def.options.find(o => o.value === value)
    return opt?.label || value
  placeholder

  return value
placeholder

// All possible columns (for column settings)
const allColumns = computed<Column[]>(() => [
  { key: 'email', label: t('admin.users.columns.user'), sortable: true placeholder,
  { key: 'id', label: 'ID', sortable: true placeholder,
  { key: 'username', label: t('admin.users.columns.username'), sortable: true placeholder,
  { key: 'notes', label: t('admin.users.columns.notes'), sortable: false placeholder,
  // Dynamic attribute columns
  ...attributeColumns.value,
  { key: 'role', label: t('admin.users.columns.role'), sortable: true placeholder,
  { key: 'subscriptions', label: t('admin.users.columns.subscriptions'), sortable: false placeholder,
  { key: 'balance', label: t('admin.users.columns.balance'), sortable: true placeholder,
  { key: 'usage', label: t('admin.users.columns.usage'), sortable: false placeholder,
  { key: 'concurrency', label: t('admin.users.columns.concurrency'), sortable: true placeholder,
  { key: 'status', label: t('admin.users.columns.status'), sortable: true placeholder,
  { key: 'created_at', label: t('admin.users.columns.created'), sortable: true placeholder,
  { key: 'actions', label: t('admin.users.columns.actions'), sortable: false placeholder
])

// Columns that can be toggled (exclude email and actions which are always visible)
const toggleableColumns = computed(() =>
  allColumns.value.filter(col => col.key !== 'email' && col.key !== 'actions')
)

// Hidden columns (stored in Set - columns NOT in this set are visible)
// This way, new columns are visible by default
const hiddenColumns = reactive<Set<string>>(new Set())

// Default hidden columns (columns hidden by default on first load)
const DEFAULT_HIDDEN_COLUMNS = ['notes', 'subscriptions', 'usage', 'concurrency']

// localStorage key for column settings
const HIDDEN_COLUMNS_KEY = 'user-hidden-columns'

// Load saved column settings
const loadSavedColumns = () => {
  try {
    const saved = localStorage.getItem(HIDDEN_COLUMNS_KEY)
    if (saved) {
      const parsed = JSON.parse(saved) as string[]
      parsed.forEach(key => hiddenColumns.add(key))
    placeholder else {
      // Use default hidden columns on first load
      DEFAULT_HIDDEN_COLUMNS.forEach(key => hiddenColumns.add(key))
    placeholder
  placeholder catch (e) {
    console.error('Failed to load saved columns:', e)
    DEFAULT_HIDDEN_COLUMNS.forEach(key => hiddenColumns.add(key))
  placeholder
placeholder

// Save column settings to localStorage
const saveColumnsToStorage = () => {
  try {
    localStorage.setItem(HIDDEN_COLUMNS_KEY, JSON.stringify([...hiddenColumns]))
  placeholder catch (e) {
    console.error('Failed to save columns:', e)
  placeholder
placeholder

// Toggle column visibility
const toggleColumn = (key: string) => {
  if (hiddenColumns.has(key)) {
    hiddenColumns.delete(key)
  placeholder else {
    hiddenColumns.add(key)
  placeholder
  saveColumnsToStorage()
placeholder

// Check if column is visible (not in hidden set)
const isColumnVisible = (key: string) => !hiddenColumns.has(key)

// Filtered columns based on visibility
const columns = computed<Column[]>(() =>
  allColumns.value.filter(col =>
    col.key === 'email' || col.key === 'actions' || !hiddenColumns.has(col.key)
  )
)

const users = ref<User[]>([])
const loading = ref(false)
const searchQuery = ref('')

// Filter values (role, status, and custom attributes)
const filters = reactive({
  role: '',
  status: ''
placeholder)
const activeAttributeFilters = reactive<Record<number, string>>({placeholder)

// Visible filters tracking (which filters are shown in the UI)
// Keys: 'role', 'status', 'attr_${idplaceholder'
const visibleFilters = reactive<Set<string>>(new Set())

// Dropdown states
const showFilterDropdown = ref(false)
const showColumnDropdown = ref(false)

// Dropdown refs for click outside detection
const filterDropdownRef = ref<HTMLElement | null>(null)
const columnDropdownRef = ref<HTMLElement | null>(null)

// localStorage keys
const FILTER_VALUES_KEY = 'user-filter-values'
const VISIBLE_FILTERS_KEY = 'user-visible-filters'

// All filterable attribute definitions (enabled attributes)
const filterableAttributes = computed(() =>
  attributeDefinitions.value.filter(def => def.enabled)
)

// Built-in filter definitions
const builtInFilters = computed(() => [
  { key: 'role', name: t('admin.users.columns.role'), type: 'select' as const placeholder,
  { key: 'status', name: t('admin.users.columns.status'), type: 'select' as const placeholder
])

// Load saved filters from localStorage
const loadSavedFilters = () => {
  try {
    // Load visible filters
    const savedVisible = localStorage.getItem(VISIBLE_FILTERS_KEY)
    if (savedVisible) {
      const parsed = JSON.parse(savedVisible) as string[]
      parsed.forEach(key => visibleFilters.add(key))
    placeholder
    // Load filter values
    const savedValues = localStorage.getItem(FILTER_VALUES_KEY)
    if (savedValues) {
      const parsed = JSON.parse(savedValues)
      if (parsed.role) filters.role = parsed.role
      if (parsed.status) filters.status = parsed.status
      if (parsed.attributes) {
        Object.assign(activeAttributeFilters, parsed.attributes)
      placeholder
    placeholder
  placeholder catch (e) {
    console.error('Failed to load saved filters:', e)
  placeholder
placeholder

// Save filters to localStorage
const saveFiltersToStorage = () => {
  try {
    // Save visible filters
    localStorage.setItem(VISIBLE_FILTERS_KEY, JSON.stringify([...visibleFilters]))
    // Save filter values
    const values = {
      role: filters.role,
      status: filters.status,
      attributes: activeAttributeFilters
    placeholder
    localStorage.setItem(FILTER_VALUES_KEY, JSON.stringify(values))
  placeholder catch (e) {
    console.error('Failed to save filters:', e)
  placeholder
placeholder

// Get attribute definition by ID
const getAttributeDefinition = (attrId: number): UserAttributeDefinition | undefined => {
  return attributeDefinitions.value.find(d => d.id === attrId)
placeholder
const usageStats = ref<Record<string, BatchUserUsageStats>>({placeholder)
// User attribute definitions and values
const attributeDefinitions = ref<UserAttributeDefinition[]>([])
const userAttributeValues = ref<Record<number, Record<number, string>>>({placeholder)
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  pages: 0
placeholder)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteDialog = ref(false)
const showApiKeysModal = ref(false)
const showAttributesModal = ref(false)
const editingUser = ref<User | null>(null)
const deletingUser = ref<User | null>(null)
const viewingUser = ref<User | null>(null)
let abortController: AbortController | null = null

// Action Menu State
const activeMenuId = ref<number | null>(null)
const menuPosition = ref<{ top: number; left: number placeholder | null>(null)
const actionButtonRefs = ref<Map<number, HTMLElement>>(new Map())

const setActionButtonRef = (userId: number, el: Element | ComponentPublicInstance | null) => {
  if (el instanceof HTMLElement) {
    actionButtonRefs.value.set(userId, el)
  placeholder else {
    actionButtonRefs.value.delete(userId)
  placeholder
placeholder

const openActionMenu = (user: User) => {
  if (activeMenuId.value === user.id) {
    closeActionMenu()
  placeholder else {
    const buttonEl = actionButtonRefs.value.get(user.id)
    if (buttonEl) {
      const rect = buttonEl.getBoundingClientRect()
      const menuWidth = 192
      const menuHeight = 240
      const padding = 8
      const viewportWidth = window.innerWidth
      const viewportHeight = window.innerHeight
      const left = Math.min(
        Math.max(rect.right - menuWidth, padding),
        Math.max(viewportWidth - menuWidth - padding, padding)
      )
      let top = rect.bottom + 4
      if (top + menuHeight > viewportHeight - padding) {
        top = Math.max(rect.top - menuHeight - 4, padding)
      placeholder
      // Position menu near the trigger, clamped to viewport
      menuPosition.value = {
        top,
        left
      placeholder
    placeholder
    activeMenuId.value = user.id
  placeholder
placeholder

const closeActionMenu = () => {
  activeMenuId.value = null
  menuPosition.value = null
placeholder

// Close menu when clicking outside
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.action-menu-trigger') && !target.closest('.action-menu-content')) {
    closeActionMenu()
  placeholder
  // Close filter dropdown when clicking outside
  if (filterDropdownRef.value && !filterDropdownRef.value.contains(target)) {
    showFilterDropdown.value = false
  placeholder
  // Close column dropdown when clicking outside
  if (columnDropdownRef.value && !columnDropdownRef.value.contains(target)) {
    showColumnDropdown.value = false
  placeholder
placeholder

// Allowed groups modal state
const showAllowedGroupsModal = ref(false)
const allowedGroupsUser = ref<User | null>(null)

// Balance (Deposit/Withdraw) modal state
const showBalanceModal = ref(false)
const balanceUser = ref<User | null>(null)
const balanceOperation = ref<'add' | 'subtract'>('add')

// 计算剩余天数
const getDaysRemaining = (expiresAt: string): number => {
  const now = new Date()
  const expires = new Date(expiresAt)
  const diffMs = expires.getTime() - now.getTime()
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24))
placeholder

const loadAttributeDefinitions = async () => {
  try {
    attributeDefinitions.value = await adminAPI.userAttributes.listEnabledDefinitions()
  placeholder catch (e) {
    console.error('Failed to load attribute definitions:', e)
  placeholder
placeholder

// Handle attributes modal close - reload definitions and users
const handleAttributesModalClose = async () => {
  showAttributesModal.value = false
  await loadAttributeDefinitions()
  loadUsers()
placeholder

const loadUsers = async () => {
  abortController?.abort()
  const currentAbortController = new AbortController()
  abortController = currentAbortController
  const { signal placeholder = currentAbortController
  loading.value = true
  try {
    // Build attribute filters from active filters
    const attrFilters: Record<number, string> = {placeholder
    for (const [attrId, value] of Object.entries(activeAttributeFilters)) {
      if (value) {
        attrFilters[Number(attrId)] = value
      placeholder
    placeholder

    const response = await adminAPI.users.list(
      pagination.page,
      pagination.page_size,
      {
        role: filters.role as any,
        status: filters.status as any,
        search: searchQuery.value || undefined,
        attributes: Object.keys(attrFilters).length > 0 ? attrFilters : undefined
      placeholder,
      { signal placeholder
    )
    if (signal.aborted) {
      return
    placeholder
    users.value = response.items
    pagination.total = response.total
    pagination.pages = response.pages

    // Load usage stats and attribute values for all users in the list
    if (response.items.length > 0) {
      const userIds = response.items.map((u) => u.id)
      // Load usage stats
      try {
        const usageResponse = await adminAPI.dashboard.getBatchUsersUsage(userIds)
        if (signal.aborted) {
          return
        placeholder
        usageStats.value = usageResponse.stats
      placeholder catch (e) {
        if (signal.aborted) {
          return
        placeholder
        console.error('Failed to load usage stats:', e)
      placeholder
      // Load attribute values
      if (attributeDefinitions.value.length > 0) {
        try {
          const attrResponse = await adminAPI.userAttributes.getBatchUserAttributes(userIds)
          if (signal.aborted) {
            return
          placeholder
          userAttributeValues.value = attrResponse.attributes
        placeholder catch (e) {
          if (signal.aborted) {
            return
          placeholder
          console.error('Failed to load user attribute values:', e)
        placeholder
      placeholder
    placeholder
  placeholder catch (error) {
    const errorInfo = error as { name?: string; code?: string placeholder
    if (errorInfo?.name === 'AbortError' || errorInfo?.name === 'CanceledError' || errorInfo?.code === 'ERR_CANCELED') {
      return
    placeholder
    appStore.showError(t('admin.users.failedToLoad'))
    console.error('Error loading users:', error)
  placeholder finally {
    if (abortController === currentAbortController) {
      loading.value = false
    placeholder
  placeholder
placeholder

let searchTimeout: ReturnType<typeof setTimeout>
const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.page = 1
    loadUsers()
  placeholder, 300)
placeholder

const handlePageChange = (page: number) => {
  pagination.page = page
  loadUsers()
placeholder

const handlePageSizeChange = (pageSize: number) => {
  pagination.page_size = pageSize
  pagination.page = 1
  loadUsers()
placeholder

// Filter helpers
const getAttributeDefinitionName = (attrId: number): string => {
  const def = attributeDefinitions.value.find(d => d.id === attrId)
  return def?.name || String(attrId)
placeholder

// Toggle a built-in filter (role/status)
const toggleBuiltInFilter = (key: string) => {
  if (visibleFilters.has(key)) {
    visibleFilters.delete(key)
    if (key === 'role') filters.role = ''
    if (key === 'status') filters.status = ''
  placeholder else {
    visibleFilters.add(key)
  placeholder
  saveFiltersToStorage()
  pagination.page = 1
  loadUsers()
placeholder

// Toggle a custom attribute filter
const toggleAttributeFilter = (attr: UserAttributeDefinition) => {
  const key = `attr_${attr.idplaceholder`
  if (visibleFilters.has(key)) {
    visibleFilters.delete(key)
    delete activeAttributeFilters[attr.id]
  placeholder else {
    visibleFilters.add(key)
    activeAttributeFilters[attr.id] = ''
  placeholder
  saveFiltersToStorage()
  pagination.page = 1
  loadUsers()
placeholder

const updateAttributeFilter = (attrId: number, value: string) => {
  activeAttributeFilters[attrId] = value
placeholder

// Apply filter and save to localStorage
const applyFilter = () => {
  saveFiltersToStorage()
  loadUsers()
placeholder

const handleEdit = (user: User) => {
  editingUser.value = user
  showEditModal.value = true
placeholder

const closeEditModal = () => {
  showEditModal.value = false
  editingUser.value = null
placeholder

const handleToggleStatus = async (user: User) => {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  try {
    await adminAPI.users.toggleStatus(user.id, newStatus)
    appStore.showSuccess(
      newStatus === 'active' ? t('admin.users.userEnabled') : t('admin.users.userDisabled')
    )
    loadUsers()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.users.failedToToggle'))
    console.error('Error toggling user status:', error)
  placeholder
placeholder

const handleViewApiKeys = (user: User) => {
  viewingUser.value = user
  showApiKeysModal.value = true
placeholder

const closeApiKeysModal = () => {
  showApiKeysModal.value = false
  viewingUser.value = null
placeholder

const handleAllowedGroups = (user: User) => {
  allowedGroupsUser.value = user
  showAllowedGroupsModal.value = true
placeholder

const closeAllowedGroupsModal = () => {
  showAllowedGroupsModal.value = false
  allowedGroupsUser.value = null
placeholder

const handleDelete = (user: User) => {
  deletingUser.value = user
  showDeleteDialog.value = true
placeholder

const confirmDelete = async () => {
  if (!deletingUser.value) return
  try {
    await adminAPI.users.delete(deletingUser.value.id)
    appStore.showSuccess(t('common.success'))
    showDeleteDialog.value = false
    deletingUser.value = null
    loadUsers()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.users.failedToDelete'))
    console.error('Error deleting user:', error)
  placeholder
placeholder

const handleDeposit = (user: User) => {
  balanceUser.value = user
  balanceOperation.value = 'add'
  showBalanceModal.value = true
placeholder

const handleWithdraw = (user: User) => {
  balanceUser.value = user
  balanceOperation.value = 'subtract'
  showBalanceModal.value = true
placeholder

const closeBalanceModal = () => {
  showBalanceModal.value = false
  balanceUser.value = null
placeholder
onMounted(async () => {
  await loadAttributeDefinitions()
  loadSavedFilters()
  loadSavedColumns()
  loadUsers()
  document.addEventListener('click', handleClickOutside)
placeholder)

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  clearTimeout(searchTimeout)
  abortController?.abort()
placeholder)
</script>
