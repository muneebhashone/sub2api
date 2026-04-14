<template>
  <div class="card">
    <!-- Header -->
    <div class="border-b border-gray-100 px-4 py-3 dark:border-dark-700">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-base font-semibold text-gray-900 dark:text-white">
            {{ t('admin.settings.payment.providerManagement') placeholderplaceholder
          </h2>
          <p class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
            {{ t('admin.settings.payment.providerManagementDesc') placeholderplaceholder
          </p>
        </div>
        <div class="flex items-center gap-2">
          <button
            type="button"
            @click="emit('refresh')"
            :disabled="loading"
            class="btn btn-secondary btn-sm"
            :title="t('common.refresh')"
          >
            <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
          </button>
          <button
            type="button"
            @click="emit('create')"
            :disabled="!canCreate"
            :class="canCreate
              ? 'btn btn-primary btn-sm'
              : 'btn btn-secondary btn-sm cursor-not-allowed opacity-50'"
          >
            {{ t('admin.settings.payment.createProvider') placeholderplaceholder
          </button>
        </div>
      </div>
    </div>

    <!-- List -->
    <div class="p-4">
      <!-- Loading -->
      <div v-if="loading && !providers.length" class="flex items-center justify-center py-6">
        <div class="h-5 w-5 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
      </div>

      <!-- Provider cards (draggable) -->
      <VueDraggable
        v-if="providers.length"
        v-model="localProviders"
        :animation="200"
        handle=".drag-handle"
        class="space-y-3"
        @end="onDragEnd"
      >
        <div v-for="p in localProviders" :key="p.id" class="flex items-start gap-2">
          <div class="drag-handle mt-3 flex cursor-grab items-center text-gray-300 hover:text-gray-500 active:cursor-grabbing dark:text-dark-600 dark:hover:text-dark-400">
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path d="M7 2a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM13 2a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM7 8a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM13 8a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM7 14a2 2 0 1 0 0 4 2 2 0 0 0 0-4zM13 14a2 2 0 1 0 0 4 2 2 0 0 0 0-4z"/>
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <ProviderCard
              :provider="p"
              :enabled="isEnabled(p.provider_key)"
              :available-types="getTypes(p.provider_key)"
              @toggle-field="(field) => emit('toggleField', p, field)"
              @toggle-type="(type) => emit('toggleType', p, type)"
              @edit="emit('edit', p)"
              @delete="emit('delete', p)"
            />
          </div>
        </div>
      </VueDraggable>

      <!-- Empty -->
      <div v-else-if="!loading" class="py-6 text-center">
        <p class="text-sm text-gray-500 dark:text-gray-400">
          {{ canCreate
            ? t('admin.settings.payment.noProviders')
            : t('admin.settings.payment.enableTypesFirst') placeholderplaceholder
        </p>
        <button
          type="button"
          v-if="canCreate"
          @click="emit('create')"
          class="btn btn-primary btn-sm mt-2"
        >
          {{ t('admin.settings.payment.createProvider') placeholderplaceholder
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { VueDraggable placeholder from 'vue-draggable-plus'
import Icon from '@/components/icons/Icon.vue'
import ProviderCard from './ProviderCard.vue'
import type { ProviderInstance placeholder from '@/types/payment'
import type { TypeOption placeholder from './providerConfig'
import { getAvailableTypes placeholder from './providerConfig'

const props = defineProps<{
  providers: ProviderInstance[]
  loading: boolean
  canCreate: boolean
  enabledPaymentTypes: string[]
  allPaymentTypes: TypeOption[]
  redirectLabel: string
placeholder>()

const emit = defineEmits<{
  refresh: []
  create: []
  edit: [provider: ProviderInstance]
  delete: [provider: ProviderInstance]
  toggleField: [provider: ProviderInstance, field: 'enabled' | 'refund_enabled' | 'allow_user_refund']
  toggleType: [provider: ProviderInstance, type: string]
  reorder: [providers: { id: number; sort_order: number placeholder[]]
placeholder>()

const { t placeholder = useI18n()

const localProviders = ref<ProviderInstance[]>([])

watch(() => props.providers, (val) => {
  localProviders.value = [...val]
placeholder, { immediate: true placeholder)

function onDragEnd() {
  const updates = localProviders.value.map((p, idx) => ({
    id: p.id,
    sort_order: idx,
  placeholder))
  emit('reorder', updates)
placeholder

function isEnabled(providerKey: string): boolean {
  return props.enabledPaymentTypes.includes(providerKey)
placeholder

function getTypes(providerKey: string): TypeOption[] {
  return getAvailableTypes(providerKey, props.allPaymentTypes, props.redirectLabel)
    .map(opt => opt.label === opt.value
      ? { ...opt, label: t(`payment.methods.${opt.valueplaceholder`, opt.value) placeholder
      : opt,
    )
placeholder
</script>
