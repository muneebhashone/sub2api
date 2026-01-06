<template>
  <div class="flex flex-wrap items-start gap-3">
    <div class="min-w-0 flex-1">
      <SearchInput
        :model-value="searchQuery"
        :placeholder="t('admin.accounts.searchAccounts')"
        @update:model-value="$emit('update:searchQuery', $event)"
        @search="$emit('change')"
      />
    </div>
    <div class="flex flex-wrap items-center gap-3">
      <Select :model-value="filters.platform" class="w-40 flex-shrink-0" :options="pOpts" @update:model-value="updatePlatform" @change="$emit('change')" />
      <Select :model-value="filters.status" class="w-40 flex-shrink-0" :options="sOpts" @update:model-value="updateStatus" @change="$emit('change')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'; import { useI18n placeholder from 'vue-i18n'; import Select from '@/components/common/Select.vue'; import SearchInput from '@/components/common/SearchInput.vue'
const props = defineProps(['searchQuery', 'filters']); const emit = defineEmits(['update:searchQuery', 'update:filters', 'change']); const { t placeholder = useI18n()
const updatePlatform = (value: string | number | boolean | null) => { emit('update:filters', { ...props.filters, platform: value placeholder) placeholder
const updateStatus = (value: string | number | boolean | null) => { emit('update:filters', { ...props.filters, status: value placeholder) placeholder
const pOpts = computed(() => [{ value: '', label: t('admin.accounts.allPlatforms') placeholder, { value: 'openai', label: 'OpenAI' placeholder, { value: 'anthropic', label: 'Anthropic' placeholder, { value: 'gemini', label: 'Gemini' placeholder])
const sOpts = computed(() => [{ value: '', label: t('admin.accounts.allStatus') placeholder, { value: 'active', label: t('admin.accounts.status.active') placeholder, { value: 'error', label: t('admin.accounts.status.error') placeholder])
</script>
