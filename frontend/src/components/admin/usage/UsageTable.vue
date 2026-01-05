<template>
  <div class="card overflow-hidden"><div class="overflow-auto">
    <DataTable :columns="cols" :data="data" :loading="loading">
      <template #cell-user="{ row placeholder"><div class="text-sm"><span class="font-medium text-gray-900 dark:text-white">{{ row.user?.email || '-' placeholderplaceholder</span><span class="ml-1 text-xs text-gray-400">#{{ row.user_id placeholderplaceholder</span></div></template>
      <template #cell-model="{ value placeholder"><span class="font-medium">{{ value placeholderplaceholder</span></template>
      <template #cell-tokens="{ row placeholder"><div class="text-sm">In: {{ row.input_tokens.toLocaleString() placeholderplaceholder / Out: {{ row.output_tokens.toLocaleString() placeholderplaceholder</div></template>
      <template #cell-cost="{ row placeholder"><span class="font-medium text-green-600">${{ row.actual_cost.toFixed(6) placeholderplaceholder</span></template>
      <template #cell-created_at="{ value placeholder"><span class="text-sm text-gray-500">{{ formatDateTime(value) placeholderplaceholder</span></template>
      <template #empty><EmptyState :message="t('usage.noRecords')" /></template>
    </DataTable>
  </div></div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'; import { useI18n placeholder from 'vue-i18n'; import { formatDateTime placeholder from '@/utils/format'; import DataTable from '@/components/common/DataTable.vue'; import EmptyState from '@/components/common/EmptyState.vue'
defineProps(['data', 'loading']); const { t placeholder = useI18n()
const cols = computed(() => [
  { key: 'user', label: t('admin.usage.user') placeholder, { key: 'model', label: t('usage.model'), sortable: true placeholder,
  { key: 'tokens', label: t('usage.tokens') placeholder, { key: 'cost', label: t('usage.cost') placeholder,
  { key: 'created_at', label: t('usage.time'), sortable: true placeholder
])
</script>