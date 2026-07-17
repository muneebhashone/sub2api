<template>
  <BaseDialog :show="show" :title="t('admin.promptAudit.events.detailTitle')" width="extra-wide" @close="$emit('close')">
    <div v-if="loading" class="py-12 text-center text-sm text-gray-500" aria-busy="true">{{ t('common.loading') placeholderplaceholder</div>
    <div v-else-if="event" class="flex flex-col">
      <div class="flex flex-wrap gap-2 border-b border-gray-200 pb-3 dark:border-dark-700" role="tablist">
        <button v-for="tab in tabs" :key="tab" type="button" role="tab" :aria-selected="activeTab === tab" class="rounded-md px-3 py-1.5 text-sm" :class="activeTab === tab ? 'bg-primary-50 text-primary-700 dark:bg-primary-950/40 dark:text-primary-300' : 'text-gray-600 dark:text-dark-300'" @click="activeTab = tab">
          {{ t(`admin.promptAudit.events.tabs.${tabplaceholder`) placeholderplaceholder
        </button>
      </div>

      <!-- Fixed panel height so switching tabs does not resize the dialog -->
      <div class="mt-5 h-[min(62vh,36rem)] overflow-y-auto" data-test="event-detail-tab-panel">
        <div v-show="activeTab === 'summary'" class="grid gap-5 lg:grid-cols-2" role="tabpanel">
          <div>
            <h4 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.promptAudit.events.promptFull') placeholderplaceholder</h4>
            <pre class="mt-2 max-h-[min(46vh,26rem)] overflow-auto whitespace-pre-wrap break-words rounded-lg bg-gray-50 p-4 text-sm text-gray-700 dark:bg-dark-900 dark:text-dark-200" data-test="summary-prompt-full">{{ displayPrompt(event) placeholderplaceholder</pre>
          </div>
          <dl class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.decision') placeholderplaceholder</dt><dd class="font-medium text-gray-900 dark:text-white">{{ formatDecisionAction(event.decision, event.action) placeholderplaceholder</dd>
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.user') placeholderplaceholder</dt><dd>{{ event.snapshot.username || '—' placeholderplaceholder</dd>
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.email') placeholderplaceholder</dt><dd>{{ event.snapshot.user_email || '—' placeholderplaceholder</dd>
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.apiKey') placeholderplaceholder</dt><dd>{{ event.snapshot.api_key_name || '—' placeholderplaceholder</dd>
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.group') placeholderplaceholder</dt><dd>{{ event.snapshot.group_name || '—' placeholderplaceholder</dd>
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.model') placeholderplaceholder</dt><dd>{{ event.snapshot.model || '—' placeholderplaceholder</dd>
            <dt class="text-gray-500">{{ t('admin.promptAudit.events.categories') placeholderplaceholder</dt><dd>{{ formatCategories(event.categories) placeholderplaceholder</dd>
          </dl>
        </div>

        <div v-show="activeTab === 'risks'" class="space-y-5" role="tabpanel">
          <div class="grid gap-4 lg:grid-cols-2">
            <section data-test="risk-prompt-preview">
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.promptAudit.events.promptFull') placeholderplaceholder</h4>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.events.promptFullHint') placeholderplaceholder</p>
              <pre class="mt-2 h-[min(46vh,26rem)] overflow-auto whitespace-pre-wrap break-words rounded-lg bg-gray-50 p-4 text-sm text-gray-700 dark:bg-dark-900 dark:text-dark-200" data-test="risk-prompt-full">{{ displayPrompt(event) placeholderplaceholder</pre>
            </section>
            <section data-test="risk-guard-return">
              <h4 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.promptAudit.events.guardReturn') placeholderplaceholder</h4>
              <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">{{ t('admin.promptAudit.events.guardReturnHint') placeholderplaceholder</p>
              <pre class="mt-2 h-[min(46vh,26rem)] overflow-auto whitespace-pre-wrap break-words rounded-lg bg-gray-50 p-4 font-mono text-xs text-gray-700 dark:bg-dark-900 dark:text-dark-200">{{ formatGuardReturn(event) placeholderplaceholder</pre>
            </section>
          </div>

          <div class="space-y-3">
            <h4 class="text-sm font-medium text-gray-900 dark:text-white">{{ t('admin.promptAudit.events.riskSummaries') placeholderplaceholder</h4>
            <article v-for="issue in event.issue_summaries" :key="`${issue.scanner_idplaceholder-${issue.codeplaceholder`" class="border-l-2 border-red-400 pl-4" data-test="risk-issue">
              <div class="flex flex-wrap items-center gap-2">
                <h5 class="font-medium text-gray-900 dark:text-white">{{ issueTitle(issue) placeholderplaceholder</h5>
                <span class="text-xs text-red-600 dark:text-red-300">{{ issueSeverity(issue) placeholderplaceholder · {{ issueAction(issue) placeholderplaceholder</span>
              </div>
              <p class="mt-1 text-sm text-gray-600 dark:text-dark-300">{{ issueDescription(issue) placeholderplaceholder</p>
              <dl class="mt-2 grid gap-1 text-xs text-gray-500 dark:text-dark-400 sm:grid-cols-2">
                <div><dt class="inline text-gray-400">{{ t('admin.promptAudit.events.categories') placeholderplaceholder · </dt><dd class="inline">{{ translateCategory(issue.category || issue.scanner_id) placeholderplaceholder</dd></div>
                <div><dt class="inline text-gray-400">{{ t('admin.promptAudit.events.score') placeholderplaceholder · </dt><dd class="inline">{{ issue.score placeholderplaceholder</dd></div>
                <div class="sm:col-span-2"><dt class="inline text-gray-400">{{ t('admin.promptAudit.events.evidence') placeholderplaceholder · </dt><dd class="inline break-words">{{ issue.evidence ? translateEvidence(issue.evidence) : '—' placeholderplaceholder</dd></div>
              </dl>
            </article>
            <p v-if="event.issue_summaries.length === 0" class="py-6 text-center text-sm text-gray-500">{{ t('admin.promptAudit.events.noRisks') placeholderplaceholder</p>
          </div>
        </div>

        <dl v-show="activeTab === 'technical'" class="grid grid-cols-[auto_minmax(0,1fr)] gap-x-4 gap-y-2 text-sm" role="tabpanel">
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.requestId') placeholderplaceholder</dt><dd class="break-all font-mono">{{ event.snapshot.request_id || '—' placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.promptHash') placeholderplaceholder</dt><dd class="break-all font-mono">{{ event.snapshot.prompt_hash placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.scanner') placeholderplaceholder</dt><dd>{{ event.scanner_backend placeholderplaceholder · {{ event.scanner_version placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.policy') placeholderplaceholder</dt><dd>{{ event.policy_id placeholderplaceholder · v{{ event.policy_version placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.guardEndpoint') placeholderplaceholder</dt><dd>{{ event.guard_endpoint_id placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.config') placeholderplaceholder</dt><dd>v{{ event.config_version placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.chunks') placeholderplaceholder</dt><dd>{{ event.chunk_total placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.latency') placeholderplaceholder</dt><dd>{{ event.latency_ms placeholderplaceholder ms</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.stage') placeholderplaceholder</dt><dd>{{ event.snapshot.stage || 'http' placeholderplaceholder</dd>
          <dt class="text-gray-500">{{ t('admin.promptAudit.events.technical.protocol') placeholderplaceholder</dt><dd>{{ event.snapshot.protocol placeholderplaceholder · {{ event.snapshot.endpoint placeholderplaceholder</dd>
        </dl>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import type { PromptAuditEvent, PromptIssueSummary placeholder from '../types'
import { SCANNER_CATALOG placeholder from '../viewModel'

const props = defineProps<{ show: boolean; event: PromptAuditEvent | null; loading: boolean placeholder>()
defineEmits<{ (event: 'close'): void placeholder>()
const { t placeholder = useI18n()
const tabs = ['summary', 'risks', 'technical'] as const
const activeTab = ref<(typeof tabs)[number]>('summary')
watch(() => props.event?.id, () => { activeTab.value = 'summary' placeholder)

const DECISIONS = new Set(['pass', 'flag', 'critical'])
const ACTIONS = new Set(['Allow', 'Warn', 'Block'])
const RISK_LEVELS = new Set(['low', 'medium', 'high', 'critical'])

function displayPrompt(event: PromptAuditEvent): string {
  return event.snapshot.full_prompt || event.snapshot.redacted_preview || '—'
placeholder

function formatDecisionAction(decision: string, action: string): string {
  const decisionLabel = DECISIONS.has(decision) ? t(`admin.promptAudit.decisions.${decisionplaceholder`) : decision
  const actionLabel = ACTIONS.has(action) ? t(`admin.promptAudit.actions.${actionplaceholder`) : action
  return `${decisionLabelplaceholder · ${actionLabelplaceholder`
placeholder
function translateCategory(category: string): string {
  return SCANNER_CATALOG.some((scanner) => scanner.id === category)
    ? t(`admin.promptAudit.scanners.${categoryplaceholder`)
    : category
placeholder
function formatCategories(categories: string[]): string {
  if (!categories.length) return '—'
  return categories.map(translateCategory).join(', ')
placeholder
function translateEvidence(value: string): string {
  const byId = SCANNER_CATALOG.find((scanner) => scanner.id === value)
  if (byId) return t(`admin.promptAudit.scanners.${byId.idplaceholder`)
  const byLabel = SCANNER_CATALOG.find((scanner) => scanner.label === value)
  if (byLabel) return t(`admin.promptAudit.scanners.${byLabel.idplaceholder`)
  return value
placeholder
function formatGuardReturn(event: PromptAuditEvent): string {
  const evidence: Record<string, string> = {placeholder
  for (const [key, value] of Object.entries(event.scanner_evidence || {placeholder)) {
    evidence[key] = translateEvidence(value)
  placeholder
  return JSON.stringify({
    decision: DECISIONS.has(event.decision) ? t(`admin.promptAudit.decisions.${event.decisionplaceholder`) : event.decision,
    risk_level: RISK_LEVELS.has(event.risk_level) ? t(`admin.promptAudit.riskLevels.${event.risk_levelplaceholder`) : event.risk_level,
    action: ACTIONS.has(event.action) ? t(`admin.promptAudit.actions.${event.actionplaceholder`) : event.action,
    categories: event.categories.map(translateCategory),
    matched_scanners: event.matched_scanners.map(translateCategory),
    scanner_scores: event.scanner_scores,
    scanner_evidence: evidence,
    scanner_backend: event.scanner_backend,
    scanner_version: event.scanner_version,
    guard_endpoint_id: event.guard_endpoint_id,
    chunk_total: event.chunk_total,
    latency_ms: event.latency_ms,
  placeholder, null, 2)
placeholder
function issueTitle(issue: PromptIssueSummary): string {
  return translateCategory(issue.category || issue.scanner_id) || issue.title
placeholder
function issueDescription(issue: PromptIssueSummary): string {
  const category = issue.category || issue.scanner_id
  const key = `admin.promptAudit.scannerDescriptions.${categoryplaceholder`
  const label = t(key)
  return label === key ? issue.description : label
placeholder
function issueSeverity(issue: PromptIssueSummary): string {
  return RISK_LEVELS.has(issue.severity) ? t(`admin.promptAudit.riskLevels.${issue.severityplaceholder`) : issue.severity_label || issue.severity
placeholder
function issueAction(issue: PromptIssueSummary): string {
  return ACTIONS.has(issue.action) ? t(`admin.promptAudit.actions.${issue.actionplaceholder`) : issue.action_label || issue.action
placeholder
</script>
