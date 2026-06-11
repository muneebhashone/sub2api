<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-wrap items-center gap-3">
          <!-- Right: Action buttons -->
          <div class="flex flex-1 flex-wrap items-center justify-end gap-2">
            <button
              @click="loadModules"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <DataTable :columns="columns" :data="modules" :loading="loading" row-key="id">
          <template #cell-id="{ value placeholder">
            <span class="font-medium text-gray-900 dark:text-white">{{ value placeholderplaceholder</span>
          </template>

          <template #cell-namespace="{ value placeholder">
            <span class="text-sm text-gray-600 dark:text-gray-300">{{ value placeholderplaceholder</span>
          </template>

          <template #cell-enabled="{ value placeholder">
            <span :class="['badge', value ? 'badge-primary' : 'badge-gray']">
              {{
                value
                  ? t('admin.modules.enabledLabels.enabled')
                  : t('admin.modules.enabledLabels.disabled')
              placeholderplaceholder
            </span>
          </template>

          <template #cell-state="{ value placeholder">
            <span :class="['badge', stateBadgeClass(value)]">
              {{ stateLabel(value) placeholderplaceholder
            </span>
          </template>

          <template #cell-error="{ value placeholder">
            <span
              v-if="value"
              class="block max-w-xs truncate text-sm text-red-600 dark:text-red-400"
              :title="value"
            >
              {{ value placeholderplaceholder
            </span>
            <span v-else class="text-sm text-gray-400 dark:text-dark-500">-</span>
          </template>

          <template #empty>
            <EmptyState
              :title="t('admin.modules.noModules')"
              :description="t('admin.modules.noModulesDescription')"
            />
          </template>
        </DataTable>
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { storeToRefs placeholder from 'pinia'
import { useAppStore placeholder from '@/stores/app'
import { useModulesStore placeholder from '@/stores/modules'
import type { ModuleState placeholder from '@/types'
import type { Column placeholder from '@/components/common/types'

import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import Icon from '@/components/icons/Icon.vue'

const { t placeholder = useI18n()
const appStore = useAppStore()
const modulesStore = useModulesStore()

const { modules, loading placeholder = storeToRefs(modulesStore)

const columns = computed<Column[]>(() => [
  { key: 'id', label: t('admin.modules.columns.id') placeholder,
  { key: 'namespace', label: t('admin.modules.columns.namespace') placeholder,
  { key: 'name', label: t('admin.modules.columns.name') placeholder,
  { key: 'enabled', label: t('admin.modules.columns.enabled') placeholder,
  { key: 'state', label: t('admin.modules.columns.state') placeholder,
  { key: 'error', label: t('admin.modules.columns.error') placeholder
])

const stateBadgeClass = (state: ModuleState) => {
  if (state === 'running') return 'badge-success'
  if (state === 'errored') return 'badge-danger'
  if (state === 'registered') return 'badge-gray'
  // stopped / provisioned: secondary tone
  return 'badge-primary'
placeholder

const stateLabel = (state: ModuleState) => {
  if (state === 'registered') return t('admin.modules.stateLabels.registered')
  if (state === 'provisioned') return t('admin.modules.stateLabels.provisioned')
  if (state === 'running') return t('admin.modules.stateLabels.running')
  if (state === 'stopped') return t('admin.modules.stateLabels.stopped')
  if (state === 'errored') return t('admin.modules.stateLabels.errored')
  return state
placeholder

async function loadModules() {
  try {
    await modulesStore.fetchModules()
  placeholder catch (error: any) {
    console.error('Error loading modules:', error)
    appStore.showError(error?.response?.data?.message || t('admin.modules.failedToLoad'))
  placeholder
placeholder

onMounted(() => {
  loadModules()
placeholder)
</script>
