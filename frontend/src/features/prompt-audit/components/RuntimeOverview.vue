<template>
  <section aria-labelledby="prompt-runtime-title" class="border-b border-gray-200 pb-6 dark:border-dark-700">
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div>
        <h2 id="prompt-runtime-title" class="text-base font-semibold text-gray-950 dark:text-white">
          {{ t('admin.promptAudit.runtime.title') placeholderplaceholder
        </h2>
        <p class="mt-1 text-sm text-gray-500 dark:text-dark-300">
          {{ t('admin.promptAudit.runtime.description') placeholderplaceholder
        </p>
      </div>
      <button type="button" class="btn btn-secondary btn-sm" :disabled="loading" @click="$emit('refresh')">
        {{ t('admin.promptAudit.actions.refresh') placeholderplaceholder
      </button>
    </div>

    <div v-if="error" role="alert" class="mt-4 rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700 dark:bg-red-950/30 dark:text-red-300">
      {{ error placeholderplaceholder
    </div>
    <div v-else-if="loading && !runtime" class="mt-5 grid grid-cols-2 gap-4 lg:grid-cols-6" aria-busy="true">
      <div v-for="index in 6" :key="index" class="h-14 animate-pulse rounded-lg bg-gray-100 dark:bg-dark-800" />
    </div>
    <template v-else-if="runtime">
      <dl class="mt-5 grid grid-cols-2 gap-x-6 gap-y-5 lg:grid-cols-6">
        <div>
          <dt class="text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.runtime.process') placeholderplaceholder</dt>
          <dd class="mt-1 flex items-center gap-2 text-sm font-semibold text-gray-900 dark:text-white">
            <span class="h-2 w-2 rounded-full" :class="statusDot(runtime.process_status)" />
            {{ t(`admin.promptAudit.status.${runtime.process_statusplaceholder`) placeholderplaceholder
          </dd>
        </div>
        <div>
          <dt class="text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.runtime.mode') placeholderplaceholder</dt>
          <dd class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">{{ t(`admin.promptAudit.mode.${runtime.effective_modeplaceholder`) placeholderplaceholder</dd>
        </div>
        <div>
          <dt class="text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.runtime.version') placeholderplaceholder</dt>
          <dd class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">
            {{ runtime.active_config_version placeholderplaceholder / {{ runtime.expected_config_version placeholderplaceholder
          </dd>
        </div>
        <div>
          <dt class="text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.runtime.workers') placeholderplaceholder</dt>
          <dd class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">{{ runtime.worker_active placeholderplaceholder / {{ runtime.worker_total placeholderplaceholder</dd>
        </div>
        <div>
          <dt class="text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.runtime.queue') placeholderplaceholder</dt>
          <dd class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">{{ runtime.queue.active placeholderplaceholder / {{ runtime.queue_capacity placeholderplaceholder</dd>
        </div>
        <div>
          <dt class="text-xs uppercase tracking-wide text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.runtime.dependencies') placeholderplaceholder</dt>
          <dd class="mt-1 text-sm font-semibold text-gray-900 dark:text-white">DB {{ runtime.database_status placeholderplaceholder · Redis {{ runtime.redis_status placeholderplaceholder</dd>
        </div>
      </dl>

      <div class="mt-5 grid gap-4 border-t border-gray-100 pt-5 dark:border-dark-800 lg:grid-cols-[1.2fr_1fr]">
        <div class="min-w-0">
          <h3 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.promptAudit.runtime.guardMetrics') placeholderplaceholder</h3>
          <div class="mt-3 flex flex-wrap gap-x-5 gap-y-2 text-sm text-gray-600 dark:text-dark-300">
            <span>{{ t('admin.promptAudit.metrics.total') placeholderplaceholder <strong>{{ runtime.guard_metrics.total placeholderplaceholder</strong></span>
            <span>{{ t('admin.promptAudit.metrics.allowed') placeholderplaceholder <strong>{{ runtime.guard_metrics.allowed placeholderplaceholder</strong></span>
            <span>{{ t('admin.promptAudit.metrics.flagged') placeholderplaceholder <strong>{{ runtime.guard_metrics.flagged placeholderplaceholder</strong></span>
            <span>{{ t('admin.promptAudit.metrics.blocked') placeholderplaceholder <strong>{{ runtime.guard_metrics.blocked placeholderplaceholder</strong></span>
            <span>{{ t('admin.promptAudit.metrics.unavailable') placeholderplaceholder <strong>{{ runtime.guard_metrics.unavailable placeholderplaceholder</strong></span>
            <span>{{ t('admin.promptAudit.metrics.timeouts') placeholderplaceholder <strong>{{ runtime.guard_metrics.timeouts placeholderplaceholder</strong></span>
            <span>{{ t('admin.promptAudit.metrics.failovers') placeholderplaceholder <strong>{{ runtime.guard_metrics.failovers placeholderplaceholder</strong></span>
            <span v-if="runtime.guard_metrics.latency_p95_ms != null">P95 <strong>{{ runtime.guard_metrics.latency_p95_ms placeholderplaceholder ms</strong></span>
          </div>
          <p class="mt-3 text-xs text-gray-500 dark:text-dark-400">
            {{ t('admin.promptAudit.runtime.queueBreakdown', {
              queued: runtime.queue.queued,
              processing: runtime.queue.processing,
              retry: runtime.queue.retry,
              done: runtime.queue.done,
              failed: runtime.queue.failed,
            placeholder) placeholderplaceholder
          </p>
          <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
            {{ t('admin.promptAudit.runtime.deliveryTotals', { enqueued: runtime.enqueued_total, dropped: runtime.dropped_total, processed: runtime.processed_total, failed: runtime.failed_total placeholder) placeholderplaceholder
          </p>
        </div>
        <div>
          <h3 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.promptAudit.runtime.latest') placeholderplaceholder</h3>
          <p class="mt-2 text-sm text-gray-600 dark:text-dark-300">
            {{ runtime.last_processed_at ? formatDate(runtime.last_processed_at) : t('admin.promptAudit.common.never') placeholderplaceholder
          </p>
          <p v-if="runtime.last_error_code" class="mt-1 break-words text-sm text-red-600 dark:text-red-300">
            {{ runtime.last_error_code placeholderplaceholder<span v-if="runtime.last_error_message"> · {{ runtime.last_error_message placeholderplaceholder</span>
          </p>
          <div v-if="Object.keys(runtime.endpoints).length" class="mt-3 flex flex-wrap gap-2">
            <span v-for="(probe, id) in runtime.endpoints" :key="id" class="rounded-full px-2 py-1 text-xs" :class="probe.ok ? 'bg-emerald-50 text-emerald-700 dark:bg-emerald-950/40 dark:text-emerald-300' : 'bg-red-50 text-red-700 dark:bg-red-950/40 dark:text-red-300'">
              {{ id placeholderplaceholder · {{ probe.status placeholderplaceholder · {{ probe.latency_ms placeholderplaceholder ms
            </span>
          </div>
        </div>
      </div>
    </template>
  </section>
</template>

<script setup lang="ts">
import { useI18n placeholder from 'vue-i18n'
import type { PromptAuditRuntime placeholder from '../types'

defineProps<{ runtime: PromptAuditRuntime | null; loading: boolean; error: string placeholder>()
defineEmits<{ (event: 'refresh'): void placeholder>()
const { t, locale placeholder = useI18n()

function formatDate(value: string): string {
  return new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium', timeStyle: 'medium' placeholder).format(new Date(value))
placeholder

function statusDot(status: string): string {
  if (status === 'running') return 'bg-emerald-500'
  if (status === 'disabled') return 'bg-gray-400'
  if (status === 'degraded') return 'bg-amber-500'
  return 'bg-red-500'
placeholder
</script>
