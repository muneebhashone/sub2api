<template>
  <div class="flex flex-wrap items-center gap-3">
    <SearchInput
      :model-value="searchQuery"
      :placeholder="t('admin.accounts.searchAccounts')"
      class="w-full sm:w-64"
      @update:model-value="$emit('update:searchQuery', $event)"
      @search="$emit('change')"
    />
    <Select :model-value="filters.platform" class="w-40" :options="pOpts" @update:model-value="updatePlatform" @change="$emit('change')" />
    <Select :model-value="filters.type" class="w-40" :options="tOpts" @update:model-value="updateType" @change="$emit('change')" />
    <Select :model-value="filters.status" class="w-40" :options="sOpts" @update:model-value="updateStatus" @change="$emit('change')" />
    <Select :model-value="filters.group" class="w-40" :options="gOpts" @update:model-value="updateGroup" @change="$emit('change')" />
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'; import { useI18n placeholder from 'vue-i18n'; import Select from '@/components/common/Select.vue'; import SearchInput from '@/components/common/SearchInput.vue'
import type { AdminGroup placeholder from '@/types'
const props = defineProps<{ searchQuery: string; filters: Record<string, any>; groups?: AdminGroup[] placeholder>()
const emit = defineEmits(['update:searchQuery', 'update:filters', 'change']); const { t placeholder = useI18n()
const updatePlatform = (value: string | number | boolean | null) => { emit('update:filters', { ...props.filters, platform: value placeholder) placeholder
const updateType = (value: string | number | boolean | null) => { emit('update:filters', { ...props.filters, type: value placeholder) placeholder
const updateStatus = (value: string | number | boolean | null) => { emit('update:filters', { ...props.filters, status: value placeholder) placeholder
const updateGroup = (value: string | number | boolean | null) => { emit('update:filters', { ...props.filters, group: value placeholder) placeholder
const pOpts = computed(() => [{ value: '', label: t('admin.accounts.allPlatforms') placeholder, { value: 'anthropic', label: 'Anthropic' placeholder, { value: 'openai', label: 'OpenAI' placeholder, { value: 'gemini', label: 'Gemini' placeholder, { value: 'antigravity', label: 'Antigravity' placeholder, { value: 'sora', label: 'Sora' placeholder])
const tOpts = computed(() => [{ value: '', label: t('admin.accounts.allTypes') placeholder, { value: 'oauth', label: t('admin.accounts.oauthType') placeholder, { value: 'setup-token', label: t('admin.accounts.setupToken') placeholder, { value: 'apikey', label: t('admin.accounts.apiKey') placeholder])
const sOpts = computed(() => [{ value: '', label: t('admin.accounts.allStatus') placeholder, { value: 'active', label: t('admin.accounts.status.active') placeholder, { value: 'inactive', label: t('admin.accounts.status.inactive') placeholder, { value: 'error', label: t('admin.accounts.status.error') placeholder, { value: 'rate_limited', label: t('admin.accounts.status.rateLimited') placeholder, { value: 'temp_unschedulable', label: t('admin.accounts.status.tempUnschedulable') placeholder])
const gOpts = computed(() => [{ value: '', label: t('admin.accounts.allGroups') placeholder, ...(props.groups || []).map(g => ({ value: String(g.id), label: g.name placeholder))])
</script>
