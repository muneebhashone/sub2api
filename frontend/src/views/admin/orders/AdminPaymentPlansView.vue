<template>
  <AppLayout>
    <div class="space-y-4">
      <!-- Actions -->
      <div class="flex items-center justify-end gap-2">
        <button @click="loadPlans" :disabled="plansLoading" class="btn btn-secondary" :title="t('common.refresh')">
          <Icon name="refresh" size="md" :class="plansLoading ? 'animate-spin' : ''" />
        </button>
        <button @click="openPlanEdit(null)" class="btn btn-primary">{{ t('payment.admin.createPlan') placeholderplaceholder</button>
      </div>

      <!-- Plans Table -->
      <DataTable :columns="planColumns" :data="plans" :loading="plansLoading">
        <template #cell-name="{ value, row placeholder">
          <span class="text-sm font-medium" :class="getPlanNameClass(row.group_id)">{{ value placeholderplaceholder</span>
        </template>
        <template #cell-group_id="{ value placeholder">
          <span v-if="isGroupMissing(value)" class="text-sm">
            <span class="text-gray-400">#{{ value placeholderplaceholder</span>
            <span class="ml-1 badge badge-danger">{{ t('payment.admin.groupMissing') placeholderplaceholder</span>
          </span>
          <GroupBadge
            v-else-if="getGroup(value)"
            :name="getGroup(value)!.name"
            :platform="getGroup(value)!.platform"
            :rate-multiplier="getGroup(value)!.rate_multiplier"
          />
          <span v-else class="text-sm text-gray-400">-</span>
        </template>
        <template #cell-price="{ value, row placeholder">
          <div class="text-sm">
            <span class="font-medium text-gray-900 dark:text-white">${{ value.toFixed(2) placeholderplaceholder</span>
            <span v-if="row.original_price" class="ml-1 text-xs text-gray-400 line-through">${{ row.original_price.toFixed(2) placeholderplaceholder</span>
          </div>
        </template>
        <template #cell-validity_days="{ value, row placeholder">
          <span class="text-sm">{{ value placeholderplaceholder {{ t('payment.admin.' + (row.validity_unit || 'days')) placeholderplaceholder</span>
        </template>
        <template #cell-for_sale="{ value, row placeholder">
          <button
            type="button"
            :class="[
              'relative inline-flex h-5 w-9 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2',
              value ? 'bg-primary-500' : 'bg-gray-300 dark:bg-dark-600'
            ]"
            @click="toggleForSale(row)"
          >
            <span :class="[
              'pointer-events-none inline-block h-4 w-4 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
              value ? 'translate-x-4' : 'translate-x-0'
            ]" />
          </button>
        </template>
        <template #cell-actions="{ row placeholder">
          <div class="flex items-center gap-2">
            <button @click="openPlanEdit(row)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-blue-50 hover:text-blue-600 dark:hover:bg-blue-900/20 dark:hover:text-blue-400">
              <Icon name="edit" size="sm" />
              <span class="text-xs">{{ t('common.edit') placeholderplaceholder</span>
            </button>
            <button @click="confirmDeletePlan(row)" class="flex flex-col items-center gap-0.5 rounded-lg p-1.5 text-gray-500 transition-colors hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/20 dark:hover:text-red-400">
              <Icon name="trash" size="sm" />
              <span class="text-xs">{{ t('common.delete') placeholderplaceholder</span>
            </button>
          </div>
        </template>
      </DataTable>
    </div>

    <!-- Plan Edit Dialog -->
    <BaseDialog :show="showPlanDialog" :title="editingPlan ? t('payment.admin.editPlan') : t('payment.admin.createPlan')" width="wide" @close="showPlanDialog = false">
      <form id="plan-form" @submit.prevent="handleSavePlan" class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="input-label">{{ t('payment.admin.planName') placeholderplaceholder</label>
            <input v-model="planForm.name" type="text" class="input" required />
          </div>
          <div>
            <label class="input-label">{{ t('payment.admin.group') placeholderplaceholder</label>
            <Select v-model="planForm.group_id" :options="groupOptions" class="w-full">
              <template #selected="{ option placeholder">
                <span v-if="option?.platform" :class="platformTextClass(String(option.platform))">{{ option.label placeholderplaceholder</span>
                <span v-else>{{ option?.label || t('payment.admin.selectGroup') placeholderplaceholder</span>
              </template>
              <template #option="{ option, selected placeholder">
                <span class="flex-1 truncate text-left" :class="option.platform ? platformTextClass(String(option.platform)) : ''">{{ option.label placeholderplaceholder</span>
                <Icon v-if="selected" name="check" size="sm" class="text-primary-500" :stroke-width="2" />
              </template>
            </Select>
          </div>
        </div>

        <!-- Group Info Preview -->
        <div v-if="selectedGroupInfo" class="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-dark-600 dark:bg-dark-800">
          <div class="mb-2 flex items-center gap-2">
            <GroupBadge :name="selectedGroupInfo.name" :platform="selectedGroupInfo.platform" :rate-multiplier="selectedGroupInfo.rate_multiplier" />
          </div>
          <div class="grid grid-cols-2 gap-2 text-xs">
            <div><span class="text-gray-500">{{ t('payment.admin.dailyLimit') placeholderplaceholder:</span> <span class="ml-1 font-medium text-gray-700 dark:text-gray-300">{{ selectedGroupInfo.daily_limit_usd != null ? '$' + selectedGroupInfo.daily_limit_usd : t('payment.admin.unlimited') placeholderplaceholder</span></div>
            <div><span class="text-gray-500">{{ t('payment.admin.weeklyLimit') placeholderplaceholder:</span> <span class="ml-1 font-medium text-gray-700 dark:text-gray-300">{{ selectedGroupInfo.weekly_limit_usd != null ? '$' + selectedGroupInfo.weekly_limit_usd : t('payment.admin.unlimited') placeholderplaceholder</span></div>
            <div><span class="text-gray-500">{{ t('payment.admin.monthlyLimit') placeholderplaceholder:</span> <span class="ml-1 font-medium text-gray-700 dark:text-gray-300">{{ selectedGroupInfo.monthly_limit_usd != null ? '$' + selectedGroupInfo.monthly_limit_usd : t('payment.admin.unlimited') placeholderplaceholder</span></div>
          </div>
        </div>

        <div><label class="input-label">{{ t('payment.admin.planDescription') placeholderplaceholder</label><textarea v-model="planForm.description" rows="2" class="input"></textarea></div>
        <div class="grid grid-cols-3 gap-4">
          <div><label class="input-label">{{ t('payment.admin.price') placeholderplaceholder</label><input v-model.number="planForm.price" type="number" step="0.01" min="0" class="input" required /></div>
          <div><label class="input-label">{{ t('payment.admin.originalPrice') placeholderplaceholder</label><input v-model.number="planForm.original_price" type="number" step="0.01" min="0" class="input" /></div>
          <div><label class="input-label">{{ t('payment.admin.sortOrder') placeholderplaceholder</label><input v-model.number="planForm.sort_order" type="number" min="0" class="input" /></div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div><label class="input-label">{{ t('payment.admin.validityDays') placeholderplaceholder</label><input v-model.number="planForm.validity_days" type="number" min="1" class="input" required /></div>
          <div><label class="input-label">{{ t('payment.admin.validityUnit') placeholderplaceholder</label><Select v-model="planForm.validity_unit" :options="validityUnitOptions" /></div>
        </div>
        <div>
          <label class="input-label">{{ t('payment.admin.features') placeholderplaceholder</label>
          <textarea v-model="planFeaturesText" rows="3" class="input" :placeholder="t('payment.admin.featuresPlaceholder')"></textarea>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.featuresHint') placeholderplaceholder</p>
        </div>
        <div class="flex items-center gap-3">
          <label class="text-sm text-gray-700 dark:text-gray-300">{{ t('payment.admin.forSale') placeholderplaceholder</label>
          <button
            type="button"
            :class="[
              'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2',
              planForm.for_sale ? 'bg-primary-500' : 'bg-gray-300 dark:bg-dark-600'
            ]"
            @click="planForm.for_sale = !planForm.for_sale"
          >
            <span :class="[
              'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out',
              planForm.for_sale ? 'translate-x-5' : 'translate-x-0'
            ]" />
          </button>
        </div>
      </form>
      <template #footer>
        <div class="flex justify-end gap-3">
          <button type="button" @click="showPlanDialog = false" class="btn btn-secondary">{{ t('common.cancel') placeholderplaceholder</button>
          <button type="submit" form="plan-form" :disabled="planSaving" class="btn btn-primary">{{ planSaving ? t('common.saving') : t('common.save') placeholderplaceholder</button>
        </div>
      </template>
    </BaseDialog>

    <ConfirmDialog :show="showDeletePlanDialog" :title="t('payment.admin.deletePlan')" :message="t('payment.admin.deletePlanConfirm')" :confirm-text="t('common.delete')" danger @confirm="handleDeletePlan" @cancel="showDeletePlanDialog = false" />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminPaymentAPI placeholder from '@/api/admin/payment'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import adminAPI from '@/api/admin'
import type { SubscriptionPlan placeholder from '@/types/payment'
import type { AdminGroup placeholder from '@/types'
import type { Column placeholder from '@/components/common/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import DataTable from '@/components/common/DataTable.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import { platformTextClass placeholder from '@/utils/platformColors'

const { t placeholder = useI18n()
const appStore = useAppStore()

// ==================== Groups ====================

const groups = ref<AdminGroup[]>([])

async function loadGroups() {
  try {
    groups.value = await adminAPI.groups.getAll()
  placeholder catch { /* ignore */ placeholder
placeholder

function getGroup(id: number): AdminGroup | undefined {
  return groups.value.find(g => g.id === id)
placeholder

function isGroupMissing(id: number): boolean {
  return id > 0 && !groups.value.find(g => g.id === id)
placeholder

function getPlanNameClass(groupId: number): string {
  const group = getGroup(groupId)
  return group ? platformTextClass(group.platform) : 'text-gray-900 dark:text-white'
placeholder

const groupOptions = computed(() => [
  { value: 0, label: t('payment.admin.selectGroup'), platform: '' placeholder,
  ...groups.value
    .filter(g => g.subscription_type === 'subscription')
    .map(g => ({
      value: g.id,
      label: `${g.nameplaceholder — ${g.platformplaceholder (${g.rate_multiplierplaceholderx)`,
      platform: g.platform,
    placeholder)),
])

const selectedGroupInfo = computed(() => {
  if (!planForm.group_id) return null
  return groups.value.find(g => g.id === planForm.group_id) || null
placeholder)

// ==================== Plans ====================

const plansLoading = ref(false)
const plans = ref<SubscriptionPlan[]>([])
const showPlanDialog = ref(false)
const showDeletePlanDialog = ref(false)
const planSaving = ref(false)
const editingPlan = ref<SubscriptionPlan | null>(null)
const deletingPlanId = ref<number | null>(null)
const planForm = reactive({ name: '', group_id: 0, description: '', price: 0, original_price: 0, validity_days: 30, validity_unit: 'days', for_sale: true, sort_order: 0 placeholder)
const planFeaturesText = ref('')

const validityUnitOptions = computed(() => [
  { value: 'days', label: t('payment.admin.days') placeholder,
  { value: 'weeks', label: t('payment.admin.weeks') placeholder,
  { value: 'months', label: t('payment.admin.months') placeholder,
])

const planColumns = computed((): Column[] => [
  { key: 'id', label: 'ID' placeholder,
  { key: 'name', label: t('payment.admin.planName') placeholder,
  { key: 'group_id', label: t('payment.admin.group') placeholder,
  { key: 'price', label: t('payment.admin.price') placeholder,
  { key: 'validity_days', label: t('payment.admin.validityDays') placeholder,
  { key: 'for_sale', label: t('payment.admin.forSale') placeholder,
  { key: 'sort_order', label: t('payment.admin.sortOrder') placeholder,
  { key: 'actions', label: t('common.actions') placeholder,
])

async function loadPlans() {
  plansLoading.value = true
  try {
    const res = await adminPaymentAPI.getPlans()
    // Backend returns features as newline-separated string; parse to array
    plans.value = (res.data || []).map((p: Omit<SubscriptionPlan, 'features'> & { features: string | string[] placeholder) => ({
      ...p,
      features: typeof p.features === 'string'
        ? p.features.split('\n').map((f: string) => f.trim()).filter(Boolean)
        : (p.features || []),
    placeholder))
  placeholder
  catch (err: unknown) { appStore.showError(extractApiErrorMessage(err, t('common.error'))) placeholder
  finally { plansLoading.value = false placeholder
placeholder

function openPlanEdit(plan: SubscriptionPlan | null) {
  editingPlan.value = plan
  if (plan) {
    Object.assign(planForm, { name: plan.name, group_id: plan.group_id, description: plan.description, price: plan.price, original_price: plan.original_price || 0, validity_days: plan.validity_days, validity_unit: plan.validity_unit || 'days', for_sale: plan.for_sale, sort_order: plan.sort_order placeholder)
    planFeaturesText.value = (plan.features || []).join('\n')
  placeholder else {
    Object.assign(planForm, { name: '', group_id: 0, description: '', price: 0, original_price: 0, validity_days: 30, validity_unit: 'days', for_sale: true, sort_order: 0 placeholder)
    planFeaturesText.value = ''
  placeholder
  showPlanDialog.value = true
placeholder

/** Build request payload with snake_case keys matching backend JSON tags */
function buildPlanPayload() {
  const features = planFeaturesText.value.split('\n').map(f => f.trim()).filter(Boolean).join('\n')
  return {
    name: planForm.name,
    group_id: planForm.group_id,
    description: planForm.description,
    price: planForm.price,
    original_price: planForm.original_price || 0,
    validity_days: planForm.validity_days,
    validity_unit: planForm.validity_unit,
    for_sale: planForm.for_sale,
    sort_order: planForm.sort_order,
    features,
  placeholder
placeholder

async function handleSavePlan() {
  planSaving.value = true
  try {
    const data = buildPlanPayload()
    if (editingPlan.value) { await adminPaymentAPI.updatePlan(editingPlan.value.id, data) placeholder
    else { await adminPaymentAPI.createPlan(data) placeholder
    appStore.showSuccess(t('common.saved')); showPlanDialog.value = false; loadPlans()
  placeholder catch (err: unknown) { appStore.showError(extractApiErrorMessage(err, t('common.error'))) placeholder
  finally { planSaving.value = false placeholder
placeholder

/** Quick toggle for_sale from the list */
async function toggleForSale(plan: SubscriptionPlan) {
  try {
    await adminPaymentAPI.updatePlan(plan.id, { for_sale: !plan.for_sale placeholder)
    plan.for_sale = !plan.for_sale
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('common.error')))
  placeholder
placeholder

function confirmDeletePlan(plan: SubscriptionPlan) { deletingPlanId.value = plan.id; showDeletePlanDialog.value = true placeholder
async function handleDeletePlan() {
  if (!deletingPlanId.value) return
  try { await adminPaymentAPI.deletePlan(deletingPlanId.value); appStore.showSuccess(t('common.deleted')); showDeletePlanDialog.value = false; loadPlans() placeholder
  catch (err: unknown) { appStore.showError(extractApiErrorMessage(err, t('common.error'))) placeholder
placeholder

// ==================== Lifecycle ====================

onMounted(() => {
  loadGroups()
  loadPlans()
placeholder)
</script>
