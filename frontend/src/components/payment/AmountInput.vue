<template>
  <div class="space-y-4">
    <!-- Quick Amount Buttons -->
    <div>
      <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t('payment.quickAmounts') placeholderplaceholder
      </label>
      <div class="grid grid-cols-3 gap-2">
        <button
          v-for="amt in filteredAmounts"
          :key="amt"
          type="button"
          :class="[
            'rounded-lg border-2 px-4 py-3 text-center font-medium transition-colors',
            modelValue === amt
              ? 'border-primary-500 bg-primary-50 text-primary-700 dark:border-primary-400 dark:bg-primary-900/40 dark:text-primary-300'
              : 'border-gray-200 bg-white text-gray-700 hover:border-gray-300 dark:border-dark-600 dark:bg-dark-800 dark:text-gray-200 dark:hover:border-dark-500',
          ]"
          @click="selectAmount(amt)"
        >
          {{ amt placeholderplaceholder
        </button>
      </div>
    </div>

    <!-- Custom Amount Input -->
    <div>
      <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ t('payment.customAmount') placeholderplaceholder
      </label>
      <div class="relative">
        <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 dark:text-dark-500">
          $
        </span>
        <input
          type="text"
          inputmode="decimal"
          :value="customText"
          :placeholder="placeholderText"
          class="input w-full py-3 pl-8 pr-4"
          @input="handleInput"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'

const props = withDefaults(defineProps<{
  amounts?: number[]
  modelValue: number | null
  min?: number
  max?: number
placeholder>(), {
  amounts: () => [10, 20, 50, 100, 200, 500, 1000, 2000, 5000],
  min: 0,
  max: 0,
placeholder)

const emit = defineEmits<{
  'update:modelValue': [value: number | null]
placeholder>()

const { t placeholder = useI18n()

const customText = ref('')

// 0 = no limit
const filteredAmounts = computed(() =>
  props.amounts.filter((a) => (props.min <= 0 || a >= props.min) && (props.max <= 0 || a <= props.max))
)

const placeholderText = computed(() => {
  if (props.min > 0 && props.max > 0) return `${props.minplaceholder - ${props.maxplaceholder`
  if (props.min > 0) return `≥ ${props.minplaceholder`
  if (props.max > 0) return `≤ ${props.maxplaceholder`
  return t('payment.enterAmount')
placeholder)

const AMOUNT_PATTERN = /^\d*(\.\d{0,2placeholder)?$/

function selectAmount(amt: number) {
  customText.value = String(amt)
  emit('update:modelValue', amt)
placeholder

function handleInput(e: Event) {
  const val = (e.target as HTMLInputElement).value
  if (!AMOUNT_PATTERN.test(val)) return
  customText.value = val
  if (val === '') {
    emit('update:modelValue', null)
    return
  placeholder
  const num = parseFloat(val)
  if (!isNaN(num) && num > 0) {
    emit('update:modelValue', num)
  placeholder else {
    emit('update:modelValue', null)
  placeholder
placeholder

watch(() => props.modelValue, (v) => {
  if (v !== null && String(v) !== customText.value) {
    customText.value = String(v)
  placeholder
placeholder, { immediate: true placeholder)
</script>
