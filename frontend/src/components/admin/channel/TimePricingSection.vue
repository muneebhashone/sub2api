<template>
  <section class="mt-3 border-t border-gray-200 pt-3 dark:border-dark-600">
    <div class="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
      <div class="min-w-0 flex-1 sm:max-w-2xl">
        <label class="block text-xs font-medium text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.form.timePricing') placeholderplaceholder
        </label>
        <div class="mt-2 grid grid-cols-1 gap-2 sm:grid-cols-2">
          <div class="min-w-0">
            <label class="block text-xs text-gray-400">
              {{ t('admin.channels.form.timezone') placeholderplaceholder
            </label>
            <Select
              :model-value="modelValue.timezone"
              :options="timezoneOptions"
              :aria-label="t('admin.channels.form.timezone')"
              data-testid="time-pricing-timezone"
              searchable
              creatable
              class="mt-1 w-full"
              @update:model-value="updateTimezone"
            />
          </div>
          <div class="min-w-0">
            <label class="block text-xs text-gray-400">
              {{ t('admin.channels.form.timePricingDayScope') placeholderplaceholder
            </label>
            <Select
              :model-value="modelValue.weekdays_only"
              :options="dayScopeOptions"
              :aria-label="t('admin.channels.form.timePricingDayScope')"
              data-testid="time-pricing-day-scope"
              class="mt-1 w-full"
              @update:model-value="updateDayScope"
            />
          </div>
        </div>
      </div>
      <button
        type="button"
        class="self-start text-xs text-primary-600 hover:text-primary-700 sm:self-end sm:pb-2"
        data-testid="add-time-period"
        @click="addPeriod"
      >
        + {{ t('admin.channels.form.addTimePeriod') placeholderplaceholder
      </button>
    </div>

    <div v-if="modelValue.periods.length > 0" class="mt-3 space-y-3">
      <div
        v-for="(period, index) in modelValue.periods"
        :key="index"
        class="grid grid-cols-1 gap-2 border-t border-gray-200 pt-3 sm:grid-cols-[minmax(0,1fr)_minmax(0,1fr)_minmax(0,1fr)_2rem] sm:items-end dark:border-dark-600"
      >
        <div class="min-w-0">
          <label :for="`${inputIdPrefixplaceholder-start-${indexplaceholder`" class="block text-xs text-gray-400">
            {{ t('admin.channels.form.startTime') placeholderplaceholder
          </label>
          <input
            :id="`${inputIdPrefixplaceholder-start-${indexplaceholder`"
            :value="period.start_time"
            type="text"
            inputmode="numeric"
            maxlength="8"
            placeholder="HH:mm:ss"
            pattern="[0-9]{2placeholder:[0-9]{2placeholder:[0-9]{2placeholder"
            autocomplete="off"
            class="input mt-1 w-full text-sm"
            @input="updatePeriod(index, 'start_time', normalizeClockTime(($event.target as HTMLInputElement).value))"
          />
        </div>
        <div class="min-w-0">
          <label :for="`${inputIdPrefixplaceholder-end-${indexplaceholder`" class="block text-xs text-gray-400">
            {{ t('admin.channels.form.endTime') placeholderplaceholder
          </label>
          <input
            :id="`${inputIdPrefixplaceholder-end-${indexplaceholder`"
            :value="period.end_time"
            type="text"
            inputmode="numeric"
            maxlength="8"
            placeholder="HH:mm:ss"
            pattern="[0-9]{2placeholder:[0-9]{2placeholder:[0-9]{2placeholder"
            autocomplete="off"
            class="input mt-1 w-full text-sm"
            @input="updatePeriod(index, 'end_time', normalizeClockTime(($event.target as HTMLInputElement).value))"
          />
        </div>
        <div class="min-w-0">
          <label :for="`${inputIdPrefixplaceholder-multiplier-${indexplaceholder`" class="block text-xs text-gray-400">
            {{ t('admin.channels.form.multiplier') placeholderplaceholder
          </label>
          <input
            :id="`${inputIdPrefixplaceholder-multiplier-${indexplaceholder`"
            :value="period.multiplier"
            type="number"
            min="0.01"
            step="0.01"
            class="input mt-1 w-full text-sm"
            @input="updatePeriod(index, 'multiplier', ($event.target as HTMLInputElement).value)"
            @blur="formatMultiplier(index, ($event.target as HTMLInputElement).value)"
          />
        </div>
        <button
          type="button"
          class="flex h-8 w-8 items-center justify-center rounded text-gray-400 hover:text-red-500"
          :title="t('admin.channels.form.removeTimePeriod')"
          :aria-label="t('admin.channels.form.removeTimePeriod')"
          :data-testid="`remove-time-period-${indexplaceholder`"
          @click="removePeriod(index)"
        >
          <Icon name="trash" size="sm" />
        </button>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, getCurrentInstance placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import {
  COMMON_TIMEZONES,
  formatTimezoneOffset,
  isValidTimePricingMultiplier,
  type TimePricingFormEntry,
  type TimePricingPeriodFormEntry,
placeholder from './types'

const { t placeholder = useI18n()

const props = defineProps<{ modelValue: TimePricingFormEntry placeholder>()
const emit = defineEmits<{ 'update:modelValue': [value: TimePricingFormEntry] placeholder>()
const inputIdPrefix = `time-pricing-${getCurrentInstance()?.uidplaceholder`

const timezoneOptions = COMMON_TIMEZONES.map(value => {
  const offset = formatTimezoneOffset(value)
  return { value, label: offset ? `${valueplaceholder (${offsetplaceholder)` : value placeholder
placeholder)

const dayScopeOptions = computed(() => [
  { value: false, label: t('admin.channels.form.timePricingEveryDay') placeholder,
  { value: true, label: t('admin.channels.form.timePricingWeekdaysOnly') placeholder,
])

function updateTimezone(value: string | number | boolean | null) {
  emit('update:modelValue', { ...props.modelValue, timezone: String(value ?? '') placeholder)
placeholder

function updateDayScope(value: string | number | boolean | null) {
  emit('update:modelValue', { ...props.modelValue, weekdays_only: value === true placeholder)
placeholder

function normalizeClockTime(value: string): string {
  const normalized = value.replace(/：/g, ':')
  return normalized === '24:00:00' ? '00:00:00' : normalized
placeholder

function addPeriod() {
  emit('update:modelValue', {
    ...props.modelValue,
    periods: [
      ...props.modelValue.periods,
      { start_time: '', end_time: '', multiplier: '1.00' placeholder,
    ],
  placeholder)
placeholder

function updatePeriod(index: number, field: keyof TimePricingPeriodFormEntry, value: string) {
  const periods = props.modelValue.periods.map((period, current) =>
    current === index ? { ...period, [field]: value placeholder : period)
  emit('update:modelValue', { ...props.modelValue, periods placeholder)
placeholder

function formatMultiplier(index: number, value: string) {
  if (!isValidTimePricingMultiplier(value)) return
  updatePeriod(index, 'multiplier', Number(value).toFixed(2))
placeholder

function removePeriod(index: number) {
  emit('update:modelValue', {
    ...props.modelValue,
    periods: props.modelValue.periods.filter((_period, current) => current !== index),
  placeholder)
placeholder
</script>
