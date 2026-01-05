<template>
  <div class="flex flex-wrap items-center gap-3">
    <SearchInput
      :model-value="searchQuery"
      :placeholder="t('admin.accounts.searchAccounts')"
      class="w-full sm:w-64"
      @update:model-value="$emit('update:searchQuery', $event)"
      @search="$emit('change')"
    />
    <Select v-model="filters.platform" class="w-40" :options="pOpts" @change="$emit('change')" />
    <Select v-model="filters.type" class="w-40" :options="tOpts" @change="$emit('change')" />
    <Select v-model="filters.status" class="w-40" :options="sOpts" @change="$emit('change')" />
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'; import { useI18n placeholder from 'vue-i18n'; import Select from '@/components/common/Select.vue'; import SearchInput from '@/components/common/SearchInput.vue'
defineProps(['searchQuery', 'filters']); defineEmits(['update:searchQuery', 'change']); const { t placeholder = useI18n()
const pOpts = computed(() => [{ value: '', label: t('admin.accounts.allPlatforms') placeholder, { value: 'openai', label: 'OpenAI' placeholder, { value: 'anthropic', label: 'Anthropic' placeholder, { value: 'gemini', label: 'Gemini' placeholder, { value: 'antigravity', label: 'Antigravity' placeholder])
const tOpts = computed(() => [{ value: '', label: t('admin.accounts.allTypes') placeholder, { value: 'oauth', label: t('admin.accounts.oauthType') placeholder, { value: 'setup-token', label: t('admin.accounts.setupToken') placeholder, { value: 'apikey', label: t('admin.accounts.apiKey') placeholder])
const sOpts = computed(() => [{ value: '', label: t('admin.accounts.allStatus') placeholder, { value: 'active', label: t('admin.accounts.status.active') placeholder, { value: 'inactive', label: t('admin.accounts.status.inactive') placeholder, { value: 'error', label: t('admin.accounts.status.error') placeholder])
</script>
