<template>
  <div v-if="visible" class="space-y-1">
    <div class="flex flex-wrap items-center gap-1.5">
      <button
        type="button"
        class="inline-flex items-center gap-0.5 rounded px-1.5 py-0.5 text-[10px] font-medium text-cyan-700 transition-colors hover:bg-cyan-50 disabled:cursor-not-allowed disabled:opacity-50 dark:text-cyan-300 dark:hover:bg-cyan-900/30"
        :disabled="loading"
        :title="t('admin.accounts.usageWindow.grokProbeTooltip')"
        @click="handleProbe"
      >
        <svg
          class="h-2.5 w-2.5"
          :class="{ 'animate-spin': loading placeholder"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        {{ t('admin.accounts.usageWindow.grokProbe') placeholderplaceholder
      </button>
    </div>

    <!-- Compact mode: parent already shows 7d/30d/prepaid or 24h — only surface errors. -->
    <div
      v-if="!compact && summary"
      class="text-[10px] text-gray-600 dark:text-gray-300"
    >
      {{ summary placeholderplaceholder
    </div>
    <div v-if="error" class="truncate text-[10px] text-red-600 dark:text-red-400" :title="error">
      {{ truncatedError placeholderplaceholder
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { adminAPI placeholder from '@/api/admin'
import type { GrokQuotaProbeResult placeholder from '@/api/admin/grok'
import type { Account placeholder from '@/types'

const props = withDefaults(
  defineProps<{
    account: Account
    /** When true, only show the probe button (+ errors). No duplicate weekly summary. */
    compact?: boolean
  placeholder>(),
  { compact: false placeholder
)

const emit = defineEmits<{ probed: [result: GrokQuotaProbeResult] placeholder>()

const { t placeholder = useI18n()

const visible = computed(() => props.account.platform === 'grok' && props.account.type === 'oauth')
const loading = ref(false)
const error = ref<string | null>(null)
const data = ref<GrokQuotaProbeResult | null>(null)

const extractErrorMessage = (e: unknown): string => {
  const err = e as {
    message?: string
    reason?: string
    response?: { data?: { message?: string; error?: string placeholder placeholder
  placeholder
  return (
    err?.message ||
    err?.reason ||
    err?.response?.data?.message ||
    err?.response?.data?.error ||
    t('common.error')
  )
placeholder

const summary = computed(() => {
  if (props.compact || !data.value) return ''
  // Non-compact fallback (rarely used): brief weekly percent if present.
  const billing = data.value.billing
  if (billing?.period_type?.toLowerCase() === 'weekly' && billing.usage_percent != null) {
    return t('admin.accounts.usageWindow.grokWeeklyUsage', {
      percent: Math.round(Math.min(100, Math.max(0, billing.usage_percent)))
    placeholder)
  placeholder
  return ''
placeholder)

const truncatedError = computed(() => {
  if (!error.value) return ''
  return error.value.length > 80 ? `${error.value.slice(0, 80)placeholder...` : error.value
placeholder)

const handleProbe = async () => {
  if (loading.value) return
  loading.value = true
  error.value = null
  try {
    data.value = await adminAPI.grok.queryQuota(props.account.id)
    error.value = data.value.probe_error || null
    emit('probed', data.value)
  placeholder catch (e) {
    error.value = extractErrorMessage(e)
  placeholder finally {
    loading.value = false
  placeholder
placeholder

watch(
  () => props.account.id,
  () => {
    data.value = null
    error.value = null
    loading.value = false
  placeholder
)
</script>
