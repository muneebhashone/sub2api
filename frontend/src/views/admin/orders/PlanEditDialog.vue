<template>
  <BaseDialog :show="show" :title="plan ? t('payment.admin.editPlan') : t('payment.admin.createPlan')" width="wide" @close="emit('close')">
    <form id="plan-form" @submit.prevent="handleSavePlan" class="space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="input-label">{{ t('payment.admin.planName') placeholderplaceholder <span class="text-red-500">*</span></label>
          <input v-model="planForm.name" type="text" class="input" required />
        </div>
        <div>
          <label class="input-label">{{ t('payment.admin.group') placeholderplaceholder <span class="text-red-500">*</span></label>
          <Select v-model="planForm.group_id" :options="groupOptions" :placeholder="t('payment.admin.selectGroup')" class="w-full">
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

      <div><label class="input-label">{{ t('payment.admin.planDescription') placeholderplaceholder <span class="text-red-500">*</span></label><textarea v-model="planForm.description" rows="2" class="input" required></textarea></div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="input-label">{{ t('payment.admin.price') placeholderplaceholder <span class="text-red-500">*</span></label>
          <input v-model.number="planForm.price" type="number" step="0.01" min="0.01" class="input" required />
          <p v-if="subscriptionCnyPreview" class="mt-1 text-xs font-medium text-primary-600 dark:text-primary-400">
            {{ t('payment.admin.subscriptionCnyPayPreview', { amount: subscriptionCnyPreview.amount placeholder) placeholderplaceholder
            <span v-if="subscriptionCnyPreview.feeRate > 0">
              {{ t('payment.admin.subscriptionCnyPayPreviewWithFee', { feeRate: subscriptionCnyPreview.feeRate, total: subscriptionCnyPreview.total placeholder) placeholderplaceholder
            </span>
          </p>
        </div>
        <div><label class="input-label">{{ t('payment.admin.originalPrice') placeholderplaceholder</label><input v-model.number="planForm.original_price" type="number" step="0.01" min="0" class="input" /></div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div><label class="input-label">{{ t('payment.admin.validity') placeholderplaceholder <span class="text-red-500">*</span></label><input v-model.number="planForm.validity_days" type="number" min="1" class="input" required /></div>
        <div><label class="input-label">{{ t('payment.admin.validityUnit') placeholderplaceholder <span class="text-red-500">*</span></label><Select v-model="planForm.validity_unit" :options="validityUnitOptions" /></div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div><label class="input-label">{{ t('payment.admin.sortOrder') placeholderplaceholder</label><input v-model.number="planForm.sort_order" type="number" min="0" class="input" /></div>
        <div>
          <label class="input-label">{{ t('payment.admin.currency') placeholderplaceholder</label>
          <input v-model="planForm.currency" type="text" maxlength="3" class="input uppercase" :placeholder="t('payment.admin.currencyPlaceholder')" />
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t('payment.admin.currencyHint') placeholderplaceholder</p>
        </div>
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
        <button type="button" @click="emit('close')" class="btn btn-secondary">{{ t('common.cancel') placeholderplaceholder</button>
        <button type="submit" form="plan-form" :disabled="saving" class="btn btn-primary">{{ saving ? t('common.saving') : t('common.save') placeholderplaceholder</button>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminPaymentAPI placeholder from '@/api/admin/payment'
import type { AdminPaymentConfig placeholder from '@/api/admin/payment'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import { formatPaymentAmount placeholder from '@/components/payment/currency'
import type { SubscriptionPlan placeholder from '@/types/payment'
import type { AdminGroup placeholder from '@/types'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import GroupBadge from '@/components/common/GroupBadge.vue'
import { platformTextClass placeholder from '@/utils/platformColors'

const props = defineProps<{
  show: boolean
  plan: SubscriptionPlan | null
  groups: AdminGroup[]
  paymentConfig?: AdminPaymentConfig | null
placeholder>()

const emit = defineEmits<{
  close: []
  saved: []
placeholder>()

const { t placeholder = useI18n()
const appStore = useAppStore()

const saving = ref(false)
const planForm = reactive({ name: '', group_id: null as number | null, description: '', price: 0, original_price: 0, currency: '', validity_days: 30, validity_unit: 'days', sort_order: 0, for_sale: true placeholder)
const planFeaturesText = ref('')

const validityUnitOptions = computed(() => [
  { value: 'days', label: t('payment.admin.days') placeholder,
  { value: 'weeks', label: t('payment.admin.weeks') placeholder,
  { value: 'months', label: t('payment.admin.months') placeholder,
])

const groupOptions = computed(() =>
  props.groups
    .filter(g => g.subscription_type === 'subscription')
    .map(g => ({
      value: g.id,
      label: `${g.nameplaceholder — ${g.platformplaceholder (${g.rate_multiplierplaceholderx)`,
      platform: g.platform,
    placeholder)),
)

const selectedGroupInfo = computed(() => {
  if (!planForm.group_id) return null
  return props.groups.find(g => g.id === planForm.group_id) || null
placeholder)

function roundCnyAmount(value: number): number {
  return Math.round(value * 100) / 100
placeholder

function ceilCnyAmount(value: number): number {
  return Math.ceil(value * 100) / 100
placeholder

const subscriptionCnyPreview = computed(() => {
  const price = Number(planForm.price) || 0
  const rate = Number(props.paymentConfig?.subscription_usd_to_cny_rate) || 0
  if (price <= 0 || rate <= 0) return null

  const amount = roundCnyAmount(price * rate)
  const feeRate = Number(props.paymentConfig?.recharge_fee_rate) || 0
  const fee = feeRate > 0 ? ceilCnyAmount((amount * feeRate) / 100) : 0
  const total = feeRate > 0 ? roundCnyAmount(amount + fee) : amount

  return {
    amount: formatPaymentAmount(amount, 'CNY'),
    feeRate,
    total: formatPaymentAmount(total, 'CNY'),
  placeholder
placeholder)

// Reset form when dialog opens
watch(() => props.show, (visible) => {
  if (!visible) return
  if (props.plan) {
    Object.assign(planForm, { name: props.plan.name, group_id: props.plan.group_id, description: props.plan.description, price: props.plan.price, original_price: props.plan.original_price || 0, currency: props.plan.currency || '', validity_days: props.plan.validity_days, validity_unit: props.plan.validity_unit || 'days', sort_order: props.plan.sort_order || 0, for_sale: props.plan.for_sale placeholder)
    planFeaturesText.value = (props.plan.features || []).join('\n')
  placeholder else {
    Object.assign(planForm, { name: '', group_id: null, description: '', price: 0, original_price: 0, currency: '', validity_days: 30, validity_unit: 'days', sort_order: 0, for_sale: true placeholder)
    planFeaturesText.value = ''
  placeholder
placeholder)

/** Build request payload with snake_case keys matching backend JSON tags */
function buildPlanPayload() {
  const features = planFeaturesText.value.split('\n').map(f => f.trim()).filter(Boolean).join('\n')
  return {
    name: planForm.name,
    group_id: planForm.group_id,
    description: planForm.description,
    price: planForm.price,
    original_price: planForm.original_price || 0,
    currency: planForm.currency.trim().toUpperCase(),
    validity_days: planForm.validity_days,
    validity_unit: planForm.validity_unit,
    sort_order: planForm.sort_order,
    for_sale: planForm.for_sale,
    features,
  placeholder
placeholder

async function handleSavePlan() {
  if (!planForm.group_id) {
    appStore.showError(t('payment.admin.groupRequired'))
    return
  placeholder
  if (!planForm.price || planForm.price <= 0) {
    appStore.showError(t('payment.admin.priceRequired'))
    return
  placeholder
  if (!planForm.validity_days || planForm.validity_days < 1) {
    appStore.showError(t('payment.admin.validityRequired'))
    return
  placeholder
  saving.value = true
  try {
    const data = buildPlanPayload()
    if (props.plan) { await adminPaymentAPI.updatePlan(props.plan.id, data) placeholder
    else { await adminPaymentAPI.createPlan(data) placeholder
    appStore.showSuccess(t('common.saved'))
    emit('close')
    emit('saved')
  placeholder catch (err: unknown) { appStore.showError(extractApiErrorMessage(err, t('common.error'))) placeholder
  finally { saving.value = false placeholder
placeholder
</script>
