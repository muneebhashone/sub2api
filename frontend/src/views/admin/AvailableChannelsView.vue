<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex flex-col justify-between gap-4 lg:flex-row lg:items-start">
          <div class="flex flex-1 flex-wrap items-center gap-3">
            <div class="relative w-full sm:w-80">
              <Icon
                name="search"
                size="md"
                class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"
              />
              <input
                v-model="searchQuery"
                type="text"
                :placeholder="t('admin.availableChannels.searchPlaceholder')"
                class="input pl-10"
              />
            </div>
          </div>

          <div class="flex w-full flex-shrink-0 flex-wrap items-center justify-end gap-3 lg:w-auto">
            <button
              @click="loadChannels"
              :disabled="loading"
              class="btn btn-secondary"
              :title="t('common.refresh', 'Refresh')"
            >
              <Icon name="refresh" size="md" :class="loading ? 'animate-spin' : ''" />
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <AvailableChannelsTable
          :columns="columns"
          :rows="filteredChannels"
          :loading="loading"
          pricing-key-prefix="admin.availableChannels.pricing"
          :no-pricing-label="t('admin.availableChannels.noPricing')"
          :no-models-label="t('admin.availableChannels.noModels')"
          :empty-label="t('admin.availableChannels.empty')"
        >
          <template #empty-groups>{{ t('admin.availableChannels.noGroups') placeholderplaceholder</template>

          <template #cell-status="{ row placeholder">
            <span
              :class="
                row.status === CHANNEL_STATUS_ACTIVE
                  ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400'
                  : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-400'
              "
              class="inline-flex items-center rounded px-2 py-0.5 text-xs font-medium"
            >
              {{ statusLabel(row.status) placeholderplaceholder
            </span>
          </template>

          <template #cell-billing_model_source="{ row placeholder">
            <span class="text-xs text-gray-700 dark:text-gray-300">
              {{ t(`admin.availableChannels.billingSource.${row.billing_model_sourceplaceholder`) placeholderplaceholder
            </span>
          </template>
        </AvailableChannelsTable>
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import AvailableChannelsTable from '@/components/channels/AvailableChannelsTable.vue'
import channelsAPI, { type AvailableChannel placeholder from '@/api/admin/channels'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import { CHANNEL_STATUS_ACTIVE, type ChannelStatus placeholder from '@/constants/channel'

const { t placeholder = useI18n()
const appStore = useAppStore()

const channels = ref<AvailableChannel[]>([])
const loading = ref(false)
const searchQuery = ref('')

const columns = computed(() => [
  { key: 'name', label: t('admin.availableChannels.columns.name') placeholder,
  { key: 'status', label: t('admin.availableChannels.columns.status') placeholder,
  { key: 'billing_model_source', label: t('admin.availableChannels.columns.billingSource') placeholder,
  { key: 'groups', label: t('admin.availableChannels.columns.groups') placeholder,
  { key: 'supported_models', label: t('admin.availableChannels.columns.supportedModels') placeholder
])

function statusLabel(status: ChannelStatus): string {
  return status === CHANNEL_STATUS_ACTIVE
    ? t('admin.availableChannels.statusActive')
    : t('admin.availableChannels.statusDisabled')
placeholder

const filteredChannels = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return channels.value
  return channels.value.filter((ch) => {
    if (ch.name.toLowerCase().includes(q)) return true
    if ((ch.description || '').toLowerCase().includes(q)) return true
    if (ch.groups.some((g) => g.name.toLowerCase().includes(q))) return true
    if (ch.supported_models.some((m) => m.name.toLowerCase().includes(q))) return true
    return false
  placeholder)
placeholder)

async function loadChannels() {
  loading.value = true
  try {
    channels.value = await channelsAPI.listAvailable()
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  placeholder finally {
    loading.value = false
  placeholder
placeholder

onMounted(loadChannels)
</script>
