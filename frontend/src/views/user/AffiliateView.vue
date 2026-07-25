<template>
  <AppLayout>
    <div class="space-y-6">
      <div v-if="loading" class="flex justify-center py-12">
        <div
          class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"
        ></div>
      </div>

      <template v-else-if="detail">
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="card p-5">
            <p class="flex items-center gap-1.5 text-sm text-gray-500 dark:text-dark-400">
              <Icon name="dollar" size="sm" class="text-primary-500" />
              {{ t('affiliate.stats.rebateRate') placeholderplaceholder
            </p>
            <p class="mt-2 text-2xl font-semibold text-primary-600 dark:text-primary-400">
              {{ formattedRebateRate placeholderplaceholder<span class="ml-0.5 text-base font-medium">%</span>
            </p>
            <p class="mt-1 text-xs text-gray-400 dark:text-dark-500">
              {{ t('affiliate.stats.rebateRateHint') placeholderplaceholder
            </p>
          </div>
          <div class="card p-5">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.invitedUsers') placeholderplaceholder</p>
            <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">
              {{ formatCount(detail.aff_count) placeholderplaceholder
            </p>
          </div>
          <div class="card p-5">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.availableQuota') placeholderplaceholder</p>
            <p class="mt-2 text-2xl font-semibold text-emerald-600 dark:text-emerald-400">
              {{ formatCurrency(detail.aff_quota) placeholderplaceholder
            </p>
          </div>
          <div class="card p-5">
            <p class="text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.stats.totalQuota') placeholderplaceholder</p>
            <p class="mt-2 text-2xl font-semibold text-gray-900 dark:text-white">
              {{ formatCurrency(detail.aff_history_quota) placeholderplaceholder
            </p>
            <p v-if="detail.aff_frozen_quota > 0" class="mt-1 text-xs text-amber-600 dark:text-amber-400">
              {{ t('affiliate.stats.frozenQuota') placeholderplaceholder: {{ formatCurrency(detail.aff_frozen_quota) placeholderplaceholder
            </p>
          </div>
        </div>

        <div class="card p-6">
          <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.title') placeholderplaceholder</h3>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.description') placeholderplaceholder</p>

          <div class="mt-5 grid gap-4 md:grid-cols-2">
            <div class="space-y-2">
              <p class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('affiliate.yourCode') placeholderplaceholder</p>
              <div class="flex flex-col items-stretch gap-2 rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 dark:border-dark-700 dark:bg-dark-900 sm:flex-row sm:items-center">
                <code class="min-w-0 break-all text-sm font-semibold text-gray-900 dark:text-white sm:flex-1 sm:truncate">{{ detail.aff_code placeholderplaceholder</code>
                <button class="btn btn-secondary btn-sm w-full sm:w-auto sm:shrink-0" @click="copyCode">
                  <Icon name="copy" size="sm" />
                  <span>{{ t('affiliate.copyCode') placeholderplaceholder</span>
                </button>
              </div>
            </div>

            <div class="space-y-2">
              <p class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('affiliate.inviteLink') placeholderplaceholder</p>
              <div class="flex flex-col items-stretch gap-2 rounded-xl border border-gray-200 bg-gray-50 px-3 py-2 dark:border-dark-700 dark:bg-dark-900 sm:flex-row sm:items-center">
                <code class="min-w-0 break-all text-sm text-gray-700 dark:text-gray-300 sm:flex-1 sm:truncate">{{ inviteLink placeholderplaceholder</code>
                <button class="btn btn-secondary btn-sm w-full sm:w-auto sm:shrink-0" @click="copyInviteLink">
                  <Icon name="copy" size="sm" />
                  <span>{{ t('affiliate.copyLink') placeholderplaceholder</span>
                </button>
              </div>
            </div>
          </div>

          <div class="mt-5 rounded-xl border border-primary-200 bg-primary-50 p-4 dark:border-primary-900/40 dark:bg-primary-900/20">
            <p class="text-sm font-medium text-primary-800 dark:text-primary-200">{{ t('affiliate.tips.title') placeholderplaceholder</p>
            <ul class="mt-2 space-y-1 text-sm text-primary-700 dark:text-primary-300">
              <li>1. {{ t('affiliate.tips.line1') placeholderplaceholder</li>
              <li>2. {{ t('affiliate.tips.line2', { rate: `${formattedRebateRateplaceholder%` placeholder) placeholderplaceholder</li>
              <li>3. {{ t('affiliate.tips.line3') placeholderplaceholder</li>
              <li v-if="detail.aff_frozen_quota > 0">4. {{ t('affiliate.tips.line4') placeholderplaceholder</li>
            </ul>
          </div>
        </div>

        <div class="card p-6">
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.transfer.title') placeholderplaceholder</h3>
              <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">{{ t('affiliate.transfer.description') placeholderplaceholder</p>
            </div>
            <button
              class="btn btn-primary"
              :disabled="transferring || detail.aff_quota <= 0"
              @click="transferQuota"
            >
              <Icon v-if="transferring" name="refresh" size="sm" class="animate-spin" />
              <Icon v-else name="dollar" size="sm" />
              <span>{{ transferring ? t('affiliate.transfer.transferring') : t('affiliate.transfer.button') placeholderplaceholder</span>
            </button>
          </div>
          <p v-if="detail.aff_quota <= 0" class="mt-3 text-sm text-amber-600 dark:text-amber-400">
            {{ t('affiliate.transfer.empty') placeholderplaceholder
          </p>
        </div>

        <div class="card p-6">
          <h3 class="text-base font-semibold text-gray-900 dark:text-white">{{ t('affiliate.invitees.title') placeholderplaceholder</h3>
          <div v-if="detail.invitees.length === 0" class="mt-4 rounded-xl border border-dashed border-gray-300 p-6 text-center text-sm text-gray-500 dark:border-dark-700 dark:text-dark-400">
            {{ t('affiliate.invitees.empty') placeholderplaceholder
          </div>
          <div v-else class="mt-4 overflow-x-auto">
            <table class="w-full min-w-[560px] text-left text-sm">
              <thead>
                <tr class="border-b border-gray-200 text-gray-500 dark:border-dark-700 dark:text-dark-400">
                  <th class="px-3 py-2 font-medium">{{ t('affiliate.invitees.columns.email') placeholderplaceholder</th>
                  <th class="px-3 py-2 font-medium">{{ t('affiliate.invitees.columns.username') placeholderplaceholder</th>
                  <th class="px-3 py-2 font-medium text-right">{{ t('affiliate.invitees.columns.rebate') placeholderplaceholder</th>
                  <th class="px-3 py-2 font-medium">{{ t('affiliate.invitees.columns.joinedAt') placeholderplaceholder</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in detail.invitees"
                  :key="item.user_id"
                  class="border-b border-gray-100 last:border-b-0 dark:border-dark-800"
                >
                  <td class="px-3 py-3 text-gray-900 dark:text-white">{{ item.email || '-' placeholderplaceholder</td>
                  <td class="px-3 py-3 text-gray-700 dark:text-gray-300">{{ item.username || '-' placeholderplaceholder</td>
                  <td class="px-3 py-3 text-right font-medium text-emerald-600 dark:text-emerald-400">{{ formatCurrency(item.total_rebate) placeholderplaceholder</td>
                  <td class="px-3 py-3 text-gray-700 dark:text-gray-300">{{ formatDateTime(item.created_at) || '-' placeholderplaceholder</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import userAPI from '@/api/user'
import type { UserAffiliateDetail placeholder from '@/types'
import { useAppStore placeholder from '@/stores/app'
import { useAuthStore placeholder from '@/stores/auth'
import { useClipboard placeholder from '@/composables/useClipboard'
import { formatCurrency, formatDateTime placeholder from '@/utils/format'
import { extractApiErrorMessage placeholder from '@/utils/apiError'

const { t placeholder = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
const { copyToClipboard placeholder = useClipboard()

const loading = ref(true)
const transferring = ref(false)
const detail = ref<UserAffiliateDetail | null>(null)

const inviteLink = computed(() => {
  if (!detail.value) return ''
  if (typeof window === 'undefined') return `/register?aff=${encodeURIComponent(detail.value.aff_code)placeholder`
  return `${window.location.originplaceholder/register?aff=${encodeURIComponent(detail.value.aff_code)placeholder`
placeholder)

// Rebate rate is a percentage in the range [0, 100]; backend already clamps it.
// We trim trailing zeros (e.g. 20.00 → "20", 12.50 → "12.5") for a cleaner UI.
const formattedRebateRate = computed(() => {
  const v = detail.value?.effective_rebate_rate_percent ?? 0
  const rounded = Math.round(v * 100) / 100
  return Number.isInteger(rounded) ? String(rounded) : rounded.toString()
placeholder)

function formatCount(value: number): string {
  return value.toLocaleString()
placeholder

async function loadAffiliateDetail(silent = false): Promise<void> {
  if (!silent) {
    loading.value = true
  placeholder
  try {
    detail.value = await userAPI.getAffiliateDetail()
  placeholder catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.loadFailed')))
  placeholder finally {
    if (!silent) {
      loading.value = false
    placeholder
  placeholder
placeholder

async function copyCode(): Promise<void> {
  if (!detail.value?.aff_code) return
  await copyToClipboard(detail.value.aff_code, t('affiliate.codeCopied'))
placeholder

async function copyInviteLink(): Promise<void> {
  if (!inviteLink.value) return
  await copyToClipboard(inviteLink.value, t('affiliate.linkCopied'))
placeholder

async function transferQuota(): Promise<void> {
  if (!detail.value || detail.value.aff_quota <= 0 || transferring.value) return
  transferring.value = true
  try {
    const resp = await userAPI.transferAffiliateQuota()
    appStore.showSuccess(t('affiliate.transfer.success', { amount: formatCurrency(resp.transferred_quota) placeholder))
    await Promise.all([
      loadAffiliateDetail(true),
      authStore.refreshUser().catch(() => undefined),
    ])
  placeholder catch (error) {
    appStore.showError(extractApiErrorMessage(error, t('affiliate.transferFailed')))
  placeholder finally {
    transferring.value = false
  placeholder
placeholder

onMounted(() => {
  void loadAffiliateDetail()
placeholder)
</script>
