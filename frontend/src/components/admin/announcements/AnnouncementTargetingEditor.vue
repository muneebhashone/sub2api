<template>
  <div class="rounded-2xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <div class="text-sm font-medium text-gray-900 dark:text-white">
          {{ t('admin.announcements.form.targetingMode') placeholderplaceholder
        </div>
        <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">
          {{ mode === 'all' ? t('admin.announcements.form.targetingAll') : t('admin.announcements.form.targetingCustom') placeholderplaceholder
        </div>
      </div>

      <div class="flex items-center gap-3">
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input
            type="radio"
            name="announcement-targeting-mode"
            value="all"
            :checked="mode === 'all'"
            @change="setMode('all')"
            class="h-4 w-4"
          />
          {{ t('admin.announcements.form.targetingAll') placeholderplaceholder
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input
            type="radio"
            name="announcement-targeting-mode"
            value="custom"
            :checked="mode === 'custom'"
            @change="setMode('custom')"
            class="h-4 w-4"
          />
          {{ t('admin.announcements.form.targetingCustom') placeholderplaceholder
        </label>
      </div>
    </div>

    <div v-if="mode === 'custom'" class="mt-4 space-y-4">
      <div class="flex items-center justify-between">
        <div class="text-sm font-medium text-gray-900 dark:text-white">
          OR
          <span class="ml-1 text-xs font-normal text-gray-500 dark:text-dark-400">
            ({{ anyOf.length placeholderplaceholder/50)
          </span>
        </div>
        <button
          type="button"
          class="btn btn-secondary"
          :disabled="anyOf.length >= 50"
          @click="addOrGroup"
        >
          <Icon name="plus" size="sm" class="mr-1" />
          {{ t('admin.announcements.form.addOrGroup') placeholderplaceholder
        </button>
      </div>

      <div v-if="anyOf.length === 0" class="rounded-xl border border-dashed border-gray-300 p-4 text-sm text-gray-500 dark:border-dark-600 dark:text-dark-400">
        {{ t('admin.announcements.form.targetingCustom') placeholderplaceholder: {{ t('admin.announcements.form.addOrGroup') placeholderplaceholder
      </div>

      <div
        v-for="(group, groupIndex) in anyOf"
        :key="groupIndex"
        class="rounded-2xl border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-800"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <div class="text-sm font-medium text-gray-900 dark:text-white">
              {{ t('admin.announcements.form.targetingCustom') placeholderplaceholder #{{ groupIndex + 1 placeholderplaceholder
              <span class="ml-2 text-xs font-normal text-gray-500 dark:text-dark-400">AND ({{ (group.all_of?.length || 0) placeholderplaceholder/50)</span>
            </div>
            <div class="mt-1 text-xs text-gray-500 dark:text-dark-400">
              {{ t('admin.announcements.form.addAndCondition') placeholderplaceholder
            </div>
          </div>

          <button
            type="button"
            class="btn btn-secondary"
            @click="removeOrGroup(groupIndex)"
          >
            <Icon name="trash" size="sm" class="mr-1" />
            {{ t('common.delete') placeholderplaceholder
          </button>
        </div>

        <div class="mt-4 space-y-3">
          <div
            v-for="(cond, condIndex) in (group.all_of || [])"
            :key="condIndex"
            class="rounded-xl border border-gray-200 bg-gray-50 p-3 dark:border-dark-700 dark:bg-dark-900/30"
          >
            <div class="flex flex-col gap-3 md:flex-row md:items-end">
              <div class="w-full md:w-52">
                <label class="input-label">{{ t('admin.announcements.form.conditionType') placeholderplaceholder</label>
                <Select
                  :model-value="cond.type"
                  :options="conditionTypeOptions"
                  @update:model-value="(v) => setConditionType(groupIndex, condIndex, v as any)"
                />
              </div>

              <div v-if="cond.type === 'subscription'" class="flex-1">
                <label class="input-label">{{ t('admin.announcements.form.selectPackages') placeholderplaceholder</label>
                <GroupSelector
                  v-model="subscriptionSelections[groupIndex][condIndex]"
                  :groups="groups"
                />
              </div>

              <div v-else class="flex flex-1 flex-col gap-3 sm:flex-row">
                <div class="w-full sm:w-44">
                  <label class="input-label">{{ t('admin.announcements.form.operator') placeholderplaceholder</label>
                  <Select
                    :model-value="cond.operator"
                    :options="balanceOperatorOptions"
                    @update:model-value="(v) => setOperator(groupIndex, condIndex, v as any)"
                  />
                </div>
                <div class="w-full sm:flex-1">
                  <label class="input-label">{{ t('admin.announcements.form.balanceValue') placeholderplaceholder</label>
                  <input
                    :value="String(cond.value ?? '')"
                    type="number"
                    step="any"
                    class="input"
                    @input="(e) => setBalanceValue(groupIndex, condIndex, (e.target as HTMLInputElement).value)"
                  />
                </div>
              </div>

              <div class="flex justify-end">
                <button
                  type="button"
                  class="btn btn-secondary"
                  @click="removeAndCondition(groupIndex, condIndex)"
                >
                  <Icon name="trash" size="sm" class="mr-1" />
                  {{ t('common.delete') placeholderplaceholder
                </button>
              </div>
            </div>
          </div>

          <div class="flex justify-end">
            <button
              type="button"
              class="btn btn-secondary"
              :disabled="(group.all_of?.length || 0) >= 50"
              @click="addAndCondition(groupIndex)"
            >
              <Icon name="plus" size="sm" class="mr-1" />
              {{ t('admin.announcements.form.addAndCondition') placeholderplaceholder
            </button>
          </div>
        </div>
      </div>

      <div v-if="validationError" class="rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700 dark:border-red-900/30 dark:bg-red-900/10 dark:text-red-300">
        {{ validationError placeholderplaceholder
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import type {
  AdminGroup,
  AnnouncementTargeting,
  AnnouncementCondition,
  AnnouncementConditionGroup,
  AnnouncementConditionType,
  AnnouncementOperator
placeholder from '@/types'

import Select from '@/components/common/Select.vue'
import GroupSelector from '@/components/common/GroupSelector.vue'
import Icon from '@/components/icons/Icon.vue'

const { t placeholder = useI18n()

const props = defineProps<{
  modelValue: AnnouncementTargeting
  groups: AdminGroup[]
placeholder>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: AnnouncementTargeting): void
placeholder>()

const anyOf = computed(() => props.modelValue?.any_of ?? [])

type Mode = 'all' | 'custom'
const mode = computed<Mode>(() => (anyOf.value.length === 0 ? 'all' : 'custom'))

const conditionTypeOptions = computed(() => [
  { value: 'subscription', label: t('admin.announcements.form.conditionSubscription') placeholder,
  { value: 'balance', label: t('admin.announcements.form.conditionBalance') placeholder
])

const balanceOperatorOptions = computed(() => [
  { value: 'gt', label: t('admin.announcements.operators.gt') placeholder,
  { value: 'gte', label: t('admin.announcements.operators.gte') placeholder,
  { value: 'lt', label: t('admin.announcements.operators.lt') placeholder,
  { value: 'lte', label: t('admin.announcements.operators.lte') placeholder,
  { value: 'eq', label: t('admin.announcements.operators.eq') placeholder
])

function setMode(next: Mode) {
  if (next === 'all') {
    emit('update:modelValue', { any_of: [] placeholder)
    return
  placeholder
  if (anyOf.value.length === 0) {
    emit('update:modelValue', { any_of: [{ all_of: [defaultSubscriptionCondition()] placeholder] placeholder)
  placeholder
placeholder

function defaultSubscriptionCondition(): AnnouncementCondition {
  return {
    type: 'subscription' as AnnouncementConditionType,
    operator: 'in' as AnnouncementOperator,
    group_ids: []
  placeholder
placeholder

function defaultBalanceCondition(): AnnouncementCondition {
  return {
    type: 'balance' as AnnouncementConditionType,
    operator: 'gte' as AnnouncementOperator,
    value: 0
  placeholder
placeholder

type TargetingDraft = {
  any_of: AnnouncementConditionGroup[]
placeholder

function updateTargeting(mutator: (draft: TargetingDraft) => void) {
  const draft: TargetingDraft = JSON.parse(JSON.stringify(props.modelValue ?? { any_of: [] placeholder))
  if (!draft.any_of) draft.any_of = []
  mutator(draft)
  emit('update:modelValue', draft)
placeholder

function addOrGroup() {
  updateTargeting((draft) => {
    if (draft.any_of.length >= 50) return
    draft.any_of.push({ all_of: [defaultSubscriptionCondition()] placeholder)
  placeholder)
placeholder

function removeOrGroup(groupIndex: number) {
  updateTargeting((draft) => {
    draft.any_of.splice(groupIndex, 1)
  placeholder)
placeholder

function addAndCondition(groupIndex: number) {
  updateTargeting((draft) => {
    const group = draft.any_of[groupIndex]
    if (!group.all_of) group.all_of = []
    if (group.all_of.length >= 50) return
    group.all_of.push(defaultSubscriptionCondition())
  placeholder)
placeholder

function removeAndCondition(groupIndex: number, condIndex: number) {
  updateTargeting((draft) => {
    const group = draft.any_of[groupIndex]
    if (!group?.all_of) return
    group.all_of.splice(condIndex, 1)
  placeholder)
placeholder

function setConditionType(groupIndex: number, condIndex: number, nextType: AnnouncementConditionType) {
  updateTargeting((draft) => {
    const group = draft.any_of[groupIndex]
    if (!group?.all_of) return

    if (nextType === 'subscription') {
      group.all_of[condIndex] = defaultSubscriptionCondition()
    placeholder else {
      group.all_of[condIndex] = defaultBalanceCondition()
    placeholder
  placeholder)
placeholder

function setOperator(groupIndex: number, condIndex: number, op: AnnouncementOperator) {
  updateTargeting((draft) => {
    const group = draft.any_of[groupIndex]
    if (!group?.all_of) return

    const cond = group.all_of[condIndex]
    if (!cond) return

    cond.operator = op
  placeholder)
placeholder

function setBalanceValue(groupIndex: number, condIndex: number, raw: string) {
  const n = raw === '' ? 0 : Number(raw)
  updateTargeting((draft) => {
    const group = draft.any_of[groupIndex]
    if (!group?.all_of) return

    const cond = group.all_of[condIndex]
    if (!cond) return

    cond.value = Number.isFinite(n) ? n : 0
  placeholder)
placeholder

// We keep group_ids selection in a parallel reactive map because GroupSelector is numeric list.
// Then we mirror it back to targeting.group_ids via a watcher.
const subscriptionSelections = reactive<Record<number, Record<number, number[]>>>({placeholder)

function ensureSelectionPath(groupIndex: number, condIndex: number) {
  if (!subscriptionSelections[groupIndex]) subscriptionSelections[groupIndex] = {placeholder
  if (!subscriptionSelections[groupIndex][condIndex]) subscriptionSelections[groupIndex][condIndex] = []
placeholder

watch(
  () => props.modelValue,
  (v) => {
    const groups = v?.any_of ?? []
    for (let gi = 0; gi < groups.length; gi++) {
      const allOf = groups[gi]?.all_of ?? []
      for (let ci = 0; ci < allOf.length; ci++) {
        const c = allOf[ci]
        if (c?.type === 'subscription') {
          ensureSelectionPath(gi, ci)
          subscriptionSelections[gi][ci] = (c.group_ids ?? []).slice()
        placeholder
      placeholder
    placeholder
  placeholder,
  { immediate: true, deep: true placeholder
)

watch(
  () => subscriptionSelections,
  () => {
    // sync back to targeting
    updateTargeting((draft) => {
      const groups = draft.any_of ?? []
      for (let gi = 0; gi < groups.length; gi++) {
        const allOf = groups[gi]?.all_of ?? []
        for (let ci = 0; ci < allOf.length; ci++) {
          const c = allOf[ci]
          if (c?.type === 'subscription') {
            ensureSelectionPath(gi, ci)
            c.operator = 'in' as AnnouncementOperator
            c.group_ids = (subscriptionSelections[gi]?.[ci] ?? []).slice()
          placeholder
        placeholder
      placeholder
    placeholder)
  placeholder,
  { deep: true placeholder
)

const validationError = computed(() => {
  if (mode.value !== 'custom') return ''

  const groups = anyOf.value
  if (groups.length === 0) return t('admin.announcements.form.addOrGroup')

  if (groups.length > 50) return 'any_of > 50'

  for (const g of groups) {
    const allOf = g?.all_of ?? []
    if (allOf.length === 0) return t('admin.announcements.form.addAndCondition')
    if (allOf.length > 50) return 'all_of > 50'

    for (const c of allOf) {
      if (c.type === 'subscription') {
        if (!c.group_ids || c.group_ids.length === 0) return t('admin.announcements.form.selectPackages')
      placeholder
    placeholder
  placeholder

  return ''
placeholder)
</script>
