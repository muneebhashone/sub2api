<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <!-- Top Toolbar: Left (search + filters) / Right (actions) -->
        <div class="flex flex-wrap items-start justify-between gap-4">
          <!-- Left: Fuzzy search + filters (wrap to multiple lines) -->
          <div class="flex flex-1 flex-wrap items-center gap-3">
            <!-- Search -->
            <div class="relative w-full sm:w-64">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
              />
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('admin.proxies.searchProxies')"
                class="input pl-10"
                @input="handleSearch"
              />
            </div>

            <!-- Filters -->
            <div class="w-full sm:w-40">
              <Select
                v-model="filters.protocol"
                :options="protocolOptions"
                :placeholder="t('admin.proxies.allProtocols')"
                @change="loadProxies"
              />
            </div>
            <div class="w-full sm:w-36">
              <Select
                v-model="filters.status"
                :options="statusOptions"
                :placeholder="t('admin.proxies.allStatus')"
                @change="loadProxies"
              />
            </div>
          </div>

          <!-- Right: Actions -->
          <div class="ml-auto flex flex-wrap items-center justify-end gap-3">
            <button
              @click="loadProxies"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
            <button
              @click="handleBatchTest"
              :disabled="batchTesting || loading"
              class="btn btn-secondary"
              :title="t('admin.proxies.testConnection')"
            >
              <Icon name="play" size="md" class="mr-2" />
              {{ t('admin.proxies.testConnection') placeholderplaceholder
            </button>
            <button
              @click="openBatchDelete"
              :disabled="selectedCount === 0"
              class="btn btn-danger"
              :title="t('admin.proxies.batchDeleteAction')"
            >
              <Icon name="trash" size="md" class="mr-2" />
              {{ t('admin.proxies.batchDeleteAction') placeholderplaceholder
            </button>
            <button @click="showImportData = true" class="btn btn-secondary">
              {{ t('admin.proxies.dataImport') placeholderplaceholder
            </button>
            <button @click="showExportDataDialog = true" class="btn btn-secondary">
              {{ t('admin.proxies.dataExport') placeholderplaceholder
            </button>
            <button @click="showCreateModal = true" class="btn btn-primary">
              <Icon name="plus" size="md" class="mr-2" />
              {{ t('admin.proxies.createProxy') placeholderplaceholder
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="proxies" :loading="loading">
          <template #header-select>
            <input
              type="checkbox"
              class="h-4 w-4 cursor-pointer rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              :checked="allVisibleSelected"
              @click.stop
              @change="toggleSelectAllVisible($event)"
            />
          </template>

          <template #cell-select="{ row placeholder">
            <input
              type="checkbox"
              class="h-4 w-4 cursor-pointer rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              :checked="selectedProxyIds.has(row.id)"
              @click.stop
              @change="toggleSelectRow(row.id, $event)"
            />
          </template>

          <template #cell-name="{ value placeholder">
            <span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
          </template>

          <template #cell-protocol="{ value placeholder">
            <span
              v-if="value"
              :class="['badge', value.startsWith('socks5') ? 'badge-primary' : 'badge-gray']"
            >
              {{ value.toUpperCase() placeholderplaceholder
            </span>
            <span v-else class="text-sm text-gray-400">-</span>
          </template>

          <template #cell-address="{ row placeholder">
            <code class="code text-xs">{{ row.host placeholderplaceholder:{{ row.port placeholderplaceholder</code>
          </template>

          <template #cell-location="{ row placeholder">
            <div class="flex items-center gap-2">
              <img
                v-if="row.country_code"
                :src="flagUrl(row.country_code)"
                :alt="row.country || row.country_code"
                class="h-4 w-6 rounded-sm"
              />
              <span v-if="formatLocation(row)" class="text-sm text-gray-700 dark:text-gray-200">
                {{ formatLocation(row) placeholderplaceholder
              </span>
              <span v-else class="text-sm text-gray-400">-</span>
            </div>
          </template>

          <template #cell-account_count="{ row, value placeholder">
            <button
              v-if="(value || 0) > 0"
              type="button"
              class="inline-flex items-center rounded bg-gray-100 px-2 py-0.5 text-xs font-medium text-primary-700 hover:bg-gray-200 dark:bg-dark-600 dark:text-primary-300 dark:hover:bg-dark-500"
              @click="openAccountsModal(row)"
            >
              {{ t('admin.groups.accountsCount', { count: value || 0 placeholder) placeholderplaceholder
            </button>
            <span
              v-else
              class="inline-flex items-center rounded bg-gray-100 px-2 py-0.5 text-xs font-medium text-gray-800 dark:bg-dark-600 dark:text-gray-300"
            >
              {{ t('admin.groups.accountsCount', { count: 0 placeholder) placeholderplaceholder
            </span>
          </template>

          <template #cell-latency="{ row placeholder">
            <span
              v-if="row.latency_status === 'failed'"
              class="badge badge-danger"
              :title="row.latency_message || undefined"
            >
              {{ t('admin.proxies.latencyFailed') placeholderplaceholder
            </span>
            <span
              v-else-if="typeof row.latency_ms === 'number'"
              :class="['badge', row.latency_ms < 200 ? 'badge-success' : 'badge-warning']"
            >
              {{ row.latency_ms placeholderplaceholderms
            </span>
            <span v-else class="text-sm text-gray-400">-</span>
          </template>

          <template #cell-status="{ value placeholder">
            <span :class="['badge', value === 'active' ? 'badge-success' : 'badge-danger']">
              {{ t('admin.accounts.status.' + value) placeholderplaceholder
            </span>
          </template>

          <template #cell-actions="{ row placeholder">
            <div class="flex items-center gap-1">
              <button
                @click="handleTestConnection(row)"
                :disabled="testingProxyIds.has(row.id)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-emerald-50 hover:text-emerald-600 disabled:cursor-not-allowed disabled:opacity-50 dark:hover:bg-emerald-900/20 dark:hover:text-emerald-400"
              >
                <svg
                  v-if="testingProxyIds.has(row.id)"
                  class="h-4 w-4 animate-spin"
                  fill="none"
                  viewBox="0 0 24 24"
                >
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"
                  ></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                <Icon v-else name="checkCircle" size="sm" />
                <span class="text-xs">{{ t('admin.proxies.testConnection') placeholderplaceholder</span>
              </button>
              <button
                @click="handleEdit(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-gray-100 hover:text-primary-600 dark:hover:bg-dark-700 dark:hover:text-primary-400"
              >
                <Icon name="edit" size="sm" />
                <span class="text-xs">{{ t('common.edit') placeholderplaceholder</span>
              </button>
              <button
                @click="handleDelete(row)"
                class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400"
              >
                <Icon name="trash" size="sm" />
                <span class="text-xs">{{ t('common.delete') placeholderplaceholder</span>
              </button>
            </div>
          </template>

          <template #empty>
            <EmptyState
              :title="t('admin.proxies.noProxiesYet')"
              :description="t('admin.proxies.createFirstProxy')"
              :action-text="t('admin.proxies.createProxy')"
              @action="showCreateModal = true"
            />
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

    <!-- Create Proxy Modal -->
    <BaseDialog
      :show="showCreateModal"
      :title="t('admin.proxies.createProxy')"
      width="normal"
      @close="closeCreateModal"
    >
      <!-- Tab Switch -->
      <div class="mb-6 flex border-b border-gray-200 dark:border-dark-600">
        <button
          type="button"
          @click="createMode = 'standard'"
          :class="[
            '-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-colors',
            createMode === 'standard'
              ? 'border-primary-500 text-primary-600 dark:text-primary-400'
              : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
          ]"
        >
          <Icon name="plus" size="sm" class="mr-1.5 inline" />
          {{ t('admin.proxies.standardAdd') placeholderplaceholder
        </button>
        <button
          type="button"
          @click="createMode = 'batch'"
          :class="[
            '-mb-px border-b-2 px-4 py-2 text-sm font-medium transition-colors',
            createMode === 'batch'
              ? 'border-primary-500 text-primary-600 dark:text-primary-400'
              : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'
          ]"
        >
          <svg
            class="mr-1.5 inline h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="1.5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M3.75 12h16.5m-16.5 3.75h16.5M3.75 19.5h16.5M5.625 4.5h12.75a1.875 1.875 0 010 3.75H5.625a1.875 1.875 0 010-3.75z"
            />
          </svg>
          {{ t('admin.proxies.batchAdd') placeholderplaceholder
        </button>
      </div>

      <!-- Standard Add Form -->
      <form
        v-if="createMode === 'standard'"
        id="create-proxy-form"
        @submit.prevent="handleCreateProxy"
        class="space-y-5"
      >
        <div>
          <label class="input-label">{{ t('admin.proxies.name') placeholderplaceholder</label>
          <input
            v-model="createForm.name"
            type="text"
            required
            class="input"
            :placeholder="t('admin.proxies.enterProxyName')"
          />
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.protocol') placeholderplaceholder</label>
          <Select v-model="createForm.protocol" :options="protocolSelectOptions" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="input-label">{{ t('admin.proxies.host') placeholderplaceholder</label>
            <input
              v-model="createForm.host"
              type="text"
              required
              :placeholder="t('admin.proxies.form.hostPlaceholder')"
              class="input"
            />
          </div>
          <div>
            <label class="input-label">{{ t('admin.proxies.port') placeholderplaceholder</label>
            <input
              v-model.number="createForm.port"
              type="number"
              required
              min="1"
              max="65535"
              :placeholder="t('admin.proxies.form.portPlaceholder')"
              class="input"
            />
          </div>
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.username') placeholderplaceholder</label>
          <input
            v-model="createForm.username"
            type="text"
            class="input"
            :placeholder="t('admin.proxies.optionalAuth')"
          />
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.password') placeholderplaceholder</label>
          <input
            v-model="createForm.password"
            type="password"
            class="input"
            :placeholder="t('admin.proxies.optionalAuth')"
          />
        </div>

      </form>

      <!-- Batch Add Form -->
      <div v-else class="space-y-5">
        <div>
          <label class="input-label">{{ t('admin.proxies.batchInput') placeholderplaceholder</label>
          <textarea
            v-model="batchInput"
            rows="10"
            class="input font-mono text-sm"
            :placeholder="t('admin.proxies.batchInputPlaceholder')"
            @input="parseBatchInput"
          ></textarea>
          <p class="input-hint mt-2">
            {{ t('admin.proxies.batchInputHint') placeholderplaceholder
          </p>
        </div>

        <!-- Parse Result -->
        <div v-if="batchParseResult.total > 0" class="rounded-lg bg-gray-50 p-4 dark:bg-dark-700">
            <div class="flex items-center gap-4 text-sm">
              <div class="flex items-center gap-1.5">
              <Icon name="checkCircle" size="sm" :stroke-width="2" class="text-primary-500" />
              <span class="text-gray-700 dark:text-gray-300">
                {{ t('admin.proxies.parsedCount', { count: batchParseResult.valid placeholder) placeholderplaceholder
              </span>
            </div>
            <div v-if="batchParseResult.invalid > 0" class="flex items-center gap-1.5">
              <Icon
                name="exclamationCircle"
                size="sm"
                :stroke-width="2"
                class="text-amber-500"
              />
              <span class="text-amber-600 dark:text-amber-400">
                {{ t('admin.proxies.invalidCount', { count: batchParseResult.invalid placeholder) placeholderplaceholder
              </span>
            </div>
            <div v-if="batchParseResult.duplicate > 0" class="flex items-center gap-1.5">
              <svg
                class="h-4 w-4 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 01-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 011.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 00-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 01-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 00-3.375-3.375h-1.5a1.125 1.125 0 01-1.125-1.125v-1.5a3.375 3.375 0 00-3.375-3.375H9.75"
                />
              </svg>
              <span class="text-gray-500 dark:text-gray-400">
                {{ t('admin.proxies.duplicateCount', { count: batchParseResult.duplicate placeholder) placeholderplaceholder
              </span>
            </div>
          </div>
        </div>

      </div>

      <template #footer>
        <div class="flex justify-end gap-3">
          <button @click="closeCreateModal" type="button" class="btn btn-secondary">
            {{ t('common.cancel') placeholderplaceholder
          </button>
          <button
            v-if="createMode === 'standard'"
            type="submit"
            form="create-proxy-form"
            :disabled="submitting"
            class="btn btn-primary"
          >
            <svg
              v-if="submitting"
              class="-ml-1 mr-2 h-4 w-4 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{ submitting ? t('admin.proxies.creating') : t('common.create') placeholderplaceholder
          </button>
          <button
            v-else
            @click="handleBatchCreate"
            type="button"
            :disabled="submitting || batchParseResult.valid === 0"
            class="btn btn-primary"
          >
            <svg
              v-if="submitting"
              class="-ml-1 mr-2 h-4 w-4 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{
              submitting
                ? t('admin.proxies.importing')
                : t('admin.proxies.importProxies', { count: batchParseResult.valid placeholder)
            placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Edit Proxy Modal -->
    <BaseDialog
      :show="showEditModal"
      :title="t('admin.proxies.editProxy')"
      width="normal"
      @close="closeEditModal"
    >
      <form
        v-if="editingProxy"
        id="edit-proxy-form"
        @submit.prevent="handleUpdateProxy"
        class="space-y-5"
      >
        <div>
          <label class="input-label">{{ t('admin.proxies.name') placeholderplaceholder</label>
          <input v-model="editForm.name" type="text" required class="input" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.protocol') placeholderplaceholder</label>
          <Select v-model="editForm.protocol" :options="protocolSelectOptions" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="input-label">{{ t('admin.proxies.host') placeholderplaceholder</label>
            <input v-model="editForm.host" type="text" required class="input" />
          </div>
          <div>
            <label class="input-label">{{ t('admin.proxies.port') placeholderplaceholder</label>
            <input
              v-model.number="editForm.port"
              type="number"
              required
              min="1"
              max="65535"
              class="input"
            />
          </div>
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.username') placeholderplaceholder</label>
          <input v-model="editForm.username" type="text" class="input" />
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.password') placeholderplaceholder</label>
          <input
            v-model="editForm.password"
            type="password"
            :placeholder="t('admin.proxies.leaveEmptyToKeep')"
            class="input"
          />
        </div>
        <div>
          <label class="input-label">{{ t('admin.proxies.status') placeholderplaceholder</label>
          <Select v-model="editForm.status" :options="editStatusOptions" />
        </div>

      </form>

      <template #footer>
        <div class="flex justify-end gap-3">
          <button @click="closeEditModal" type="button" class="btn btn-secondary">
            {{ t('common.cancel') placeholderplaceholder
          </button>
          <button
            v-if="editingProxy"
            type="submit"
            form="edit-proxy-form"
            :disabled="submitting"
            class="btn btn-primary"
          >
            <svg
              v-if="submitting"
              class="-ml-1 mr-2 h-4 w-4 animate-spin"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {{ submitting ? t('admin.proxies.updating') : t('common.update') placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDialog
      :show="showDeleteDialog"
      :title="t('admin.proxies.deleteProxy')"
      :message="t('admin.proxies.deleteConfirm', { name: deletingProxy?.name placeholder)"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="confirmDelete"
      @cancel="showDeleteDialog = false"
    />

    <!-- Batch Delete Confirmation Dialog -->
    <ConfirmDialog
      :show="showBatchDeleteDialog"
      :title="t('admin.proxies.batchDelete')"
      :message="t('admin.proxies.batchDeleteConfirm', { count: selectedCount placeholder)"
      :confirm-text="t('common.delete')"
      :cancel-text="t('common.cancel')"
      :danger="true"
      @confirm="confirmBatchDelete"
      @cancel="showBatchDeleteDialog = false"
    />
    <ConfirmDialog
      :show="showExportDataDialog"
      :title="t('admin.proxies.dataExport')"
      :message="t('admin.proxies.dataExportConfirmMessage')"
      :confirm-text="t('admin.proxies.dataExportConfirm')"
      :cancel-text="t('common.cancel')"
      @confirm="handleExportData"
      @cancel="showExportDataDialog = false"
    />

    <ImportDataModal
      :show="showImportData"
      @close="showImportData = false"
      @imported="handleDataImported"
    />

    <!-- Proxy Accounts Dialog -->
    <BaseDialog
      :show="showAccountsModal"
      :title="t('admin.proxies.accountsTitle', { name: accountsProxy?.name || '' placeholder)"
      width="normal"
      @close="closeAccountsModal"
    >
      <div v-if="accountsLoading" class="flex items-center justify-center py-8 text-sm text-gray-500">
        <Icon name="refresh" size="md" class="mr-2 animate-spin" />
        {{ t('common.loading') placeholderplaceholder
      </div>
      <div v-else-if="proxyAccounts.length === 0" class="py-6 text-center text-sm text-gray-500">
        {{ t('admin.proxies.accountsEmpty') placeholderplaceholder
      </div>
      <div v-else class="max-h-80 overflow-auto">
        <table class="min-w-full divide-y divide-gray-200 text-sm dark:divide-dark-700">
          <thead class="bg-gray-50 text-xs uppercase text-gray-500 dark:bg-dark-800 dark:text-dark-400">
            <tr>
              <th class="px-4 py-2 text-left">{{ t('admin.proxies.accountName') placeholderplaceholder</th>
              <th class="px-4 py-2 text-left">{{ t('admin.accounts.columns.platformType') placeholderplaceholder</th>
              <th class="px-4 py-2 text-left">{{ t('admin.proxies.accountNotes') placeholderplaceholder</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200 bg-white dark:divide-dark-700 dark:bg-dark-900">
            <tr v-for="account in proxyAccounts" :key="account.id">
              <td class="px-4 py-2 font-medium text-gray-900 dark:text-white">{{ account.name placeholderplaceholder</td>
              <td class="px-4 py-2">
                <PlatformTypeBadge :platform="account.platform" :type="account.type" />
              </td>
              <td class="px-4 py-2 text-gray-600 dark:text-gray-300">
                {{ account.notes || '-' placeholderplaceholder
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <button @click="closeAccountsModal" class="btn btn-secondary">
            {{ t('common.close') placeholderplaceholder
          </button>
        </div>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import type { Proxy, ProxyAccountSummary, ProxyProtocol placeholder from '@/types'
import type { Column placeholder from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import Pagination from '@/components/common/Pagination.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ImportDataModal from '@/components/admin/proxy/ImportDataModal.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import PlatformTypeBadge from '@/components/common/PlatformTypeBadge.vue'

const { t placeholder = useI18n()
const appStore = useAppStore()

const columns = computed<Column[]>(() => [
  { key: 'select', label: '', sortable: false placeholder,
  { key: 'name', label: t('admin.proxies.columns.name'), sortable: true placeholder,
  { key: 'protocol', label: t('admin.proxies.columns.protocol'), sortable: true placeholder,
  { key: 'address', label: t('admin.proxies.columns.address'), sortable: false placeholder,
  { key: 'location', label: t('admin.proxies.columns.location'), sortable: false placeholder,
  { key: 'account_count', label: t('admin.proxies.columns.accounts'), sortable: true placeholder,
  { key: 'latency', label: t('admin.proxies.columns.latency'), sortable: false placeholder,
  { key: 'status', label: t('admin.proxies.columns.status'), sortable: true placeholder,
  { key: 'actions', label: t('admin.proxies.columns.actions'), sortable: false placeholder
])

// Filter options
const protocolOptions = computed(() => [
  { value: '', label: t('admin.proxies.allProtocols') placeholder,
  { value: 'http', label: 'HTTP' placeholder,
  { value: 'https', label: 'HTTPS' placeholder,
  { value: 'socks5', label: 'SOCKS5' placeholder,
  { value: 'socks5h', label: 'SOCKS5H' placeholder
])

const statusOptions = computed(() => [
  { value: '', label: t('admin.proxies.allStatus') placeholder,
  { value: 'active', label: t('admin.accounts.status.active') placeholder,
  { value: 'inactive', label: t('admin.accounts.status.inactive') placeholder
])

// Form options
const protocolSelectOptions = computed(() => [
  { value: 'http', label: t('admin.proxies.protocols.http') placeholder,
  { value: 'https', label: t('admin.proxies.protocols.https') placeholder,
  { value: 'socks5', label: t('admin.proxies.protocols.socks5') placeholder,
  { value: 'socks5h', label: t('admin.proxies.protocols.socks5h') placeholder
])

const editStatusOptions = computed(() => [
  { value: 'active', label: t('admin.accounts.status.active') placeholder,
  { value: 'inactive', label: t('admin.accounts.status.inactive') placeholder
])

const proxies = ref<Proxy[]>([])
const loading = ref(false)
const searchQuery = ref('')
const filters = reactive({
  protocol: '',
  status: ''
placeholder)
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0,
  pages: 0
placeholder)

const showCreateModal = ref(false)
const showEditModal = ref(false)
const showImportData = ref(false)
const showDeleteDialog = ref(false)
const showBatchDeleteDialog = ref(false)
const showExportDataDialog = ref(false)
const showAccountsModal = ref(false)
const submitting = ref(false)
const exportingData = ref(false)
const testingProxyIds = ref<Set<number>>(new Set())
const batchTesting = ref(false)
const selectedProxyIds = ref<Set<number>>(new Set())
const accountsProxy = ref<Proxy | null>(null)
const proxyAccounts = ref<ProxyAccountSummary[]>([])
const accountsLoading = ref(false)
const editingProxy = ref<Proxy | null>(null)
const deletingProxy = ref<Proxy | null>(null)

const selectedCount = computed(() => selectedProxyIds.value.size)
const allVisibleSelected = computed(() => {
  if (proxies.value.length === 0) return false
  return proxies.value.every((proxy) => selectedProxyIds.value.has(proxy.id))
placeholder)

// Batch import state
const createMode = ref<'standard' | 'batch'>('standard')
const batchInput = ref('')
const batchParseResult = reactive({
  total: 0,
  valid: 0,
  invalid: 0,
  duplicate: 0,
  proxies: [] as Array<{
    protocol: ProxyProtocol
    host: string
    port: number
    username: string
    password: string
  placeholder>
placeholder)

const createForm = reactive({
  name: '',
  protocol: 'http' as ProxyProtocol,
  host: '',
  port: 8080,
  username: '',
  password: ''
placeholder)

const editForm = reactive({
  name: '',
  protocol: 'http' as ProxyProtocol,
  host: '',
  port: 8080,
  username: '',
  password: '',
  status: 'active' as 'active' | 'inactive'
placeholder)

let abortController: AbortController | null = null

const isAbortError = (error: unknown) => {
  if (!error || typeof error !== 'object') return false
  const maybeError = error as { name?: string; code?: string placeholder
  return maybeError.name === 'AbortError' || maybeError.code === 'ERR_CANCELED'
placeholder

const toggleSelectRow = (id: number, event: Event) => {
  const target = event.target as HTMLInputElement
  const next = new Set(selectedProxyIds.value)
  if (target.checked) {
    next.add(id)
  placeholder else {
    next.delete(id)
  placeholder
  selectedProxyIds.value = next
placeholder

const toggleSelectAllVisible = (event: Event) => {
  const target = event.target as HTMLInputElement
  const next = new Set(selectedProxyIds.value)
  for (const proxy of proxies.value) {
    if (target.checked) {
      next.add(proxy.id)
    placeholder else {
      next.delete(proxy.id)
    placeholder
  placeholder
  selectedProxyIds.value = next
placeholder

const loadProxies = async () => {
  if (abortController) {
    abortController.abort()
  placeholder
  const currentAbortController = new AbortController()
  abortController = currentAbortController
  loading.value = true
  try {
    const response = await adminAPI.proxies.list(pagination.page, pagination.page_size, {
      protocol: filters.protocol || undefined,
      status: filters.status as any,
      search: searchQuery.value || undefined
    placeholder, { signal: currentAbortController.signal placeholder)
    if (currentAbortController.signal.aborted || abortController !== currentAbortController) {
      return
    placeholder
    proxies.value = response.items
    pagination.total = response.total
    pagination.pages = response.pages
  placeholder catch (error) {
    if (isAbortError(error)) {
      return
    placeholder
    appStore.showError(t('admin.proxies.failedToLoad'))
    console.error('Error loading proxies:', error)
  placeholder finally {
    if (abortController === currentAbortController) {
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
    loadProxies()
  placeholder, 300)
placeholder

const handlePageChange = (page: number) => {
  pagination.page = page
  loadProxies()
placeholder

const handlePageSizeChange = (pageSize: number) => {
  pagination.page_size = pageSize
  pagination.page = 1
  loadProxies()
placeholder

const closeCreateModal = () => {
  showCreateModal.value = false
  createMode.value = 'standard'
  createForm.name = ''
  createForm.protocol = 'http'
  createForm.host = ''
  createForm.port = 8080
  createForm.username = ''
  createForm.password = ''
  batchInput.value = ''
  batchParseResult.total = 0
  batchParseResult.valid = 0
  batchParseResult.invalid = 0
  batchParseResult.duplicate = 0
  batchParseResult.proxies = []
placeholder

const handleDataImported = () => {
  showImportData.value = false
  loadProxies()
placeholder

// Parse proxy URL: protocol://user:pass@host:port or protocol://host:port
const parseProxyUrl = (
  line: string
): {
  protocol: ProxyProtocol
  host: string
  port: number
  username: string
  password: string
placeholder | null => {
  const trimmed = line.trim()
  if (!trimmed) return null

  // Regex to parse proxy URL (supports http, https, socks5, socks5h)
  const regex = /^(https?|socks5h?):\/\/(?:([^:@]+):([^@]+)@)?([^:]+):(\d+)$/i
  const match = trimmed.match(regex)

  if (!match) return null

  const [, protocol, username, password, host, port] = match
  const portNum = parseInt(port, 10)

  if (portNum < 1 || portNum > 65535) return null

  return {
    protocol: protocol.toLowerCase() as ProxyProtocol,
    host: host.trim(),
    port: portNum,
    username: username?.trim() || '',
    password: password?.trim() || ''
  placeholder
placeholder

const parseBatchInput = () => {
  const lines = batchInput.value.split('\n').filter((l) => l.trim())
  const seen = new Set<string>()
  const proxies: typeof batchParseResult.proxies = []
  let invalid = 0
  let duplicate = 0

  for (const line of lines) {
    const parsed = parseProxyUrl(line)
    if (!parsed) {
      invalid++
      continue
    placeholder

    // Check for duplicates (same host:port:username:password)
    const key = `${parsed.hostplaceholder:${parsed.portplaceholder:${parsed.usernameplaceholder:${parsed.passwordplaceholder`
    if (seen.has(key)) {
      duplicate++
      continue
    placeholder
    seen.add(key)
    proxies.push(parsed)
  placeholder

  batchParseResult.total = lines.length
  batchParseResult.valid = proxies.length
  batchParseResult.invalid = invalid
  batchParseResult.duplicate = duplicate
  batchParseResult.proxies = proxies
placeholder

const handleBatchCreate = async () => {
  if (batchParseResult.valid === 0) return

  submitting.value = true
  try {
    const result = await adminAPI.proxies.batchCreate(batchParseResult.proxies)
    const created = result.created || 0
    const skipped = result.skipped || 0

    if (created > 0) {
      appStore.showSuccess(t('admin.proxies.batchImportSuccess', { created, skipped placeholder))
    placeholder else {
      appStore.showInfo(t('admin.proxies.batchImportAllSkipped', { skipped placeholder))
    placeholder

    closeCreateModal()
    loadProxies()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.failedToImport'))
    console.error('Error batch creating proxies:', error)
  placeholder finally {
    submitting.value = false
  placeholder
placeholder

const handleCreateProxy = async () => {
  if (!createForm.name.trim()) {
    appStore.showError(t('admin.proxies.nameRequired'))
    return
  placeholder
  if (!createForm.host.trim()) {
    appStore.showError(t('admin.proxies.hostRequired'))
    return
  placeholder
  if (createForm.port < 1 || createForm.port > 65535) {
    appStore.showError(t('admin.proxies.portInvalid'))
    return
  placeholder
  submitting.value = true
  try {
    await adminAPI.proxies.create({
      name: createForm.name.trim(),
      protocol: createForm.protocol,
      host: createForm.host.trim(),
      port: createForm.port,
      username: createForm.username.trim() || null,
      password: createForm.password.trim() || null
    placeholder)
    appStore.showSuccess(t('admin.proxies.proxyCreated'))
    closeCreateModal()
    loadProxies()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.failedToCreate'))
    console.error('Error creating proxy:', error)
  placeholder finally {
    submitting.value = false
  placeholder
placeholder

const handleEdit = (proxy: Proxy) => {
  editingProxy.value = proxy
  editForm.name = proxy.name
  editForm.protocol = proxy.protocol
  editForm.host = proxy.host
  editForm.port = proxy.port
  editForm.username = proxy.username || ''
  editForm.password = ''
  editForm.status = proxy.status
  showEditModal.value = true
placeholder

const closeEditModal = () => {
  showEditModal.value = false
  editingProxy.value = null
placeholder

const handleUpdateProxy = async () => {
  if (!editingProxy.value) return
  if (!editForm.name.trim()) {
    appStore.showError(t('admin.proxies.nameRequired'))
    return
  placeholder
  if (!editForm.host.trim()) {
    appStore.showError(t('admin.proxies.hostRequired'))
    return
  placeholder
  if (editForm.port < 1 || editForm.port > 65535) {
    appStore.showError(t('admin.proxies.portInvalid'))
    return
  placeholder

  submitting.value = true
  try {
    const updateData: any = {
      name: editForm.name.trim(),
      protocol: editForm.protocol,
      host: editForm.host.trim(),
      port: editForm.port,
      username: editForm.username.trim() || null,
      status: editForm.status
    placeholder

    // Only include password if it was changed
    const trimmedPassword = editForm.password.trim()
    if (trimmedPassword) {
      updateData.password = trimmedPassword
    placeholder

    await adminAPI.proxies.update(editingProxy.value.id, updateData)
    appStore.showSuccess(t('admin.proxies.proxyUpdated'))
    closeEditModal()
    loadProxies()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.failedToUpdate'))
    console.error('Error updating proxy:', error)
  placeholder finally {
    submitting.value = false
  placeholder
placeholder

const applyLatencyResult = (
  proxyId: number,
  result: {
    success: boolean
    latency_ms?: number
    message?: string
    ip_address?: string
    country?: string
    country_code?: string
    region?: string
    city?: string
  placeholder
) => {
  const target = proxies.value.find((proxy) => proxy.id === proxyId)
  if (!target) return
  if (result.success) {
    target.latency_status = 'success'
    target.latency_ms = result.latency_ms
    target.ip_address = result.ip_address
    target.country = result.country
    target.country_code = result.country_code
    target.region = result.region
    target.city = result.city
  placeholder else {
    target.latency_status = 'failed'
    target.latency_ms = undefined
    target.ip_address = undefined
    target.country = undefined
    target.country_code = undefined
    target.region = undefined
    target.city = undefined
  placeholder
  target.latency_message = result.message
placeholder

const formatLocation = (proxy: Proxy) => {
  const parts = [proxy.country, proxy.city].filter(Boolean) as string[]
  return parts.join(' · ')
placeholder

const flagUrl = (code: string) =>
  `https://unpkg.com/flag-icons/flags/4x3/${code.toLowerCase()placeholder.svg`

const startTestingProxy = (proxyId: number) => {
  testingProxyIds.value = new Set([...testingProxyIds.value, proxyId])
placeholder

const stopTestingProxy = (proxyId: number) => {
  const next = new Set(testingProxyIds.value)
  next.delete(proxyId)
  testingProxyIds.value = next
placeholder

const runProxyTest = async (proxyId: number, notify: boolean) => {
  startTestingProxy(proxyId)
  try {
    const result = await adminAPI.proxies.testProxy(proxyId)
    applyLatencyResult(proxyId, result)
    if (notify) {
      if (result.success) {
        const message = result.latency_ms
          ? t('admin.proxies.proxyWorkingWithLatency', { latency: result.latency_ms placeholder)
          : t('admin.proxies.proxyWorking')
        appStore.showSuccess(message)
      placeholder else {
        appStore.showError(result.message || t('admin.proxies.proxyTestFailed'))
      placeholder
    placeholder
    return result
  placeholder catch (error: any) {
    const message = error.response?.data?.detail || t('admin.proxies.failedToTest')
    applyLatencyResult(proxyId, { success: false, message placeholder)
    if (notify) {
      appStore.showError(message)
    placeholder
    console.error('Error testing proxy:', error)
    return null
  placeholder finally {
    stopTestingProxy(proxyId)
  placeholder
placeholder

const handleTestConnection = async (proxy: Proxy) => {
  await runProxyTest(proxy.id, true)
placeholder

const fetchAllProxiesForBatch = async (): Promise<Proxy[]> => {
  const pageSize = 200
  const result: Proxy[] = []
  let page = 1
  let totalPages = 1

  while (page <= totalPages) {
    const response = await adminAPI.proxies.list(
      page,
      pageSize,
      {
        protocol: filters.protocol || undefined,
        status: filters.status as any,
        search: searchQuery.value || undefined
      placeholder
    )
    result.push(...response.items)
    totalPages = response.pages || 1
    page++
  placeholder

  return result
placeholder

const runBatchProxyTests = async (ids: number[]) => {
  if (ids.length === 0) return
  const concurrency = 5
  let index = 0

  const worker = async () => {
    while (index < ids.length) {
      const current = ids[index]
      index++
      await runProxyTest(current, false)
    placeholder
  placeholder

  const workers = Array.from({ length: Math.min(concurrency, ids.length) placeholder, () => worker())
  await Promise.all(workers)
placeholder

const handleBatchTest = async () => {
  if (batchTesting.value) return

  batchTesting.value = true
  try {
    let ids: number[] = []
    if (selectedCount.value > 0) {
      ids = Array.from(selectedProxyIds.value)
    placeholder else {
      const allProxies = await fetchAllProxiesForBatch()
      ids = allProxies.map((proxy) => proxy.id)
    placeholder

    if (ids.length === 0) {
      appStore.showInfo(t('admin.proxies.batchTestEmpty'))
      return
    placeholder

    await runBatchProxyTests(ids)
    appStore.showSuccess(t('admin.proxies.batchTestDone', { count: ids.length placeholder))
    loadProxies()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.batchTestFailed'))
    console.error('Error batch testing proxies:', error)
  placeholder finally {
    batchTesting.value = false
  placeholder
placeholder

const formatExportTimestamp = () => {
  const now = new Date()
  const pad2 = (value: number) => String(value).padStart(2, '0')
  return `${now.getFullYear()placeholder${pad2(now.getMonth() + 1)placeholder${pad2(now.getDate())placeholder${pad2(now.getHours())placeholder${pad2(now.getMinutes())placeholder${pad2(now.getSeconds())placeholder`
placeholder

const handleExportData = async () => {
  if (exportingData.value) return
  exportingData.value = true
  try {
    const dataPayload = await adminAPI.proxies.exportData({
      protocol: filters.protocol || undefined,
      status: (filters.status || undefined) as 'active' | 'inactive' | undefined,
      search: searchQuery.value || undefined
    placeholder)
    const timestamp = formatExportTimestamp()
    const filename = `sub2api-proxy-${timestampplaceholder.json`
    const blob = new Blob([JSON.stringify(dataPayload, null, 2)], { type: 'application/json' placeholder)
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    link.click()
    URL.revokeObjectURL(url)
    appStore.showSuccess(t('admin.proxies.dataExported'))
  placeholder catch (error: any) {
    appStore.showError(error?.message || t('admin.proxies.dataExportFailed'))
  placeholder finally {
    exportingData.value = false
    showExportDataDialog.value = false
  placeholder
placeholder

const handleDelete = (proxy: Proxy) => {
  if ((proxy.account_count || 0) > 0) {
    appStore.showError(t('admin.proxies.deleteBlockedInUse'))
    return
  placeholder
  deletingProxy.value = proxy
  showDeleteDialog.value = true
placeholder

const openBatchDelete = () => {
  if (selectedCount.value === 0) {
    return
  placeholder
  showBatchDeleteDialog.value = true
placeholder

const confirmDelete = async () => {
  if (!deletingProxy.value) return

  try {
    await adminAPI.proxies.delete(deletingProxy.value.id)
    appStore.showSuccess(t('admin.proxies.proxyDeleted'))
    showDeleteDialog.value = false
    if (selectedProxyIds.value.has(deletingProxy.value.id)) {
      const next = new Set(selectedProxyIds.value)
      next.delete(deletingProxy.value.id)
      selectedProxyIds.value = next
    placeholder
    deletingProxy.value = null
    loadProxies()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.failedToDelete'))
    console.error('Error deleting proxy:', error)
  placeholder
placeholder

const confirmBatchDelete = async () => {
  const ids = Array.from(selectedProxyIds.value)
  if (ids.length === 0) {
    showBatchDeleteDialog.value = false
    return
  placeholder

  try {
    const result = await adminAPI.proxies.batchDelete(ids)
    const deleted = result.deleted_ids?.length || 0
    const skipped = result.skipped?.length || 0

    if (deleted > 0) {
      appStore.showSuccess(t('admin.proxies.batchDeleteDone', { deleted, skipped placeholder))
    placeholder else if (skipped > 0) {
      appStore.showInfo(t('admin.proxies.batchDeleteSkipped', { skipped placeholder))
    placeholder

    selectedProxyIds.value = new Set()
    showBatchDeleteDialog.value = false
    loadProxies()
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.batchDeleteFailed'))
    console.error('Error batch deleting proxies:', error)
  placeholder
placeholder

const openAccountsModal = async (proxy: Proxy) => {
  accountsProxy.value = proxy
  proxyAccounts.value = []
  accountsLoading.value = true
  showAccountsModal.value = true

  try {
    proxyAccounts.value = await adminAPI.proxies.getProxyAccounts(proxy.id)
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('admin.proxies.accountsFailed'))
    console.error('Error loading proxy accounts:', error)
  placeholder finally {
    accountsLoading.value = false
  placeholder
placeholder

const closeAccountsModal = () => {
  showAccountsModal.value = false
  accountsProxy.value = null
  proxyAccounts.value = []
placeholder

onMounted(() => {
  loadProxies()
placeholder)

onUnmounted(() => {
  clearTimeout(searchTimeout)
  abortController?.abort()
placeholder)
</script>
