<template>
  <div class="rounded-lg border border-gray-200 bg-gray-50 p-3 dark:border-dark-600 dark:bg-dark-800">
    <!-- Header: Models + Billing Mode + Remove -->
    <div class="mb-3 flex items-start gap-2">
      <div class="flex-1">
        <label class="text-xs font-medium text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.form.models', '模型列表') placeholderplaceholder
        </label>
        <ModelTagInput
          :models="entry.models"
          @update:models="emit('update', { ...entry, models: $event placeholder)"
          :placeholder="t('admin.channels.form.modelsPlaceholder', '输入模型名后按回车添加，支持通配符 *')"
          class="mt-1"
        />
      </div>
      <div class="w-40">
        <label class="text-xs font-medium text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.form.billingMode', '计费模式') placeholderplaceholder
        </label>
        <Select
          :modelValue="entry.billing_mode"
          @update:modelValue="emit('update', { ...entry, billing_mode: $event as BillingMode, intervals: [] placeholder)"
          :options="billingModeOptions"
          class="mt-1"
        />
      </div>
      <button
        type="button"
        @click="emit('remove')"
        class="mt-5 rounded p-1 text-gray-400 hover:text-red-500"
      >
        <Icon name="trash" size="sm" />
      </button>
    </div>

    <!-- Token mode -->
    <div v-if="entry.billing_mode === 'token'">
      <!-- Default prices (fallback when no interval matches) -->
      <label class="text-xs font-medium text-gray-500 dark:text-gray-400">
        {{ t('admin.channels.form.defaultPrices', '默认价格（未命中区间时使用）') placeholderplaceholder
        <span class="ml-1 font-normal text-gray-400">$/MTok</span>
      </label>
      <div class="mt-1 grid grid-cols-2 gap-2 sm:grid-cols-5">
        <div>
          <label class="text-xs text-gray-400">{{ t('admin.channels.form.inputPrice', '输入') placeholderplaceholder</label>
          <input :value="entry.input_price" @input="emitField('input_price', ($event.target as HTMLInputElement).value)"
            type="number" step="any" min="0" class="input mt-0.5 text-sm" :placeholder="t('admin.channels.form.pricePlaceholder', '默认')" />
        </div>
        <div>
          <label class="text-xs text-gray-400">{{ t('admin.channels.form.outputPrice', '输出') placeholderplaceholder</label>
          <input :value="entry.output_price" @input="emitField('output_price', ($event.target as HTMLInputElement).value)"
            type="number" step="any" min="0" class="input mt-0.5 text-sm" :placeholder="t('admin.channels.form.pricePlaceholder', '默认')" />
        </div>
        <div>
          <label class="text-xs text-gray-400">{{ t('admin.channels.form.cacheWritePrice', '缓存写入') placeholderplaceholder</label>
          <input :value="entry.cache_write_price" @input="emitField('cache_write_price', ($event.target as HTMLInputElement).value)"
            type="number" step="any" min="0" class="input mt-0.5 text-sm" :placeholder="t('admin.channels.form.pricePlaceholder', '默认')" />
        </div>
        <div>
          <label class="text-xs text-gray-400">{{ t('admin.channels.form.cacheReadPrice', '缓存读取') placeholderplaceholder</label>
          <input :value="entry.cache_read_price" @input="emitField('cache_read_price', ($event.target as HTMLInputElement).value)"
            type="number" step="any" min="0" class="input mt-0.5 text-sm" :placeholder="t('admin.channels.form.pricePlaceholder', '默认')" />
        </div>
        <div>
          <label class="text-xs text-gray-400">{{ t('admin.channels.form.imageTokenPrice', '图片输出') placeholderplaceholder</label>
          <input :value="entry.image_output_price" @input="emitField('image_output_price', ($event.target as HTMLInputElement).value)"
            type="number" step="any" min="0" class="input mt-0.5 text-sm" :placeholder="t('admin.channels.form.pricePlaceholder', '默认')" />
        </div>
      </div>

      <!-- Token intervals -->
      <div class="mt-3">
        <div class="flex items-center justify-between">
          <label class="text-xs font-medium text-gray-500 dark:text-gray-400">
            {{ t('admin.channels.form.intervals', '上下文区间定价（可选）') placeholderplaceholder
            <span class="ml-1 font-normal text-gray-400">(min, max]</span>
          </label>
          <button type="button" @click="addInterval" class="text-xs text-primary-600 hover:text-primary-700">
            + {{ t('admin.channels.form.addInterval', '添加区间') placeholderplaceholder
          </button>
        </div>
        <div v-if="entry.intervals && entry.intervals.length > 0" class="mt-2 space-y-2">
          <IntervalRow
            v-for="(iv, idx) in entry.intervals"
            :key="idx"
            :interval="iv"
            :mode="entry.billing_mode"
            @update="updateInterval(idx, $event)"
            @remove="removeInterval(idx)"
          />
        </div>
      </div>
    </div>

    <!-- Per-request mode -->
    <div v-else-if="entry.billing_mode === 'per_request'">
      <div class="flex items-center justify-between">
        <label class="text-xs font-medium text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.form.requestTiers', '按次计费层级') placeholderplaceholder
        </label>
        <button type="button" @click="addInterval" class="text-xs text-primary-600 hover:text-primary-700">
          + {{ t('admin.channels.form.addTier', '添加层级') placeholderplaceholder
        </button>
      </div>
      <div v-if="entry.intervals && entry.intervals.length > 0" class="mt-2 space-y-2">
        <IntervalRow
          v-for="(iv, idx) in entry.intervals"
          :key="idx"
          :interval="iv"
          :mode="entry.billing_mode"
          @update="updateInterval(idx, $event)"
          @remove="removeInterval(idx)"
        />
      </div>
      <div v-else class="mt-2 rounded border border-dashed border-gray-300 p-3 text-center text-xs text-gray-400 dark:border-dark-500">
        {{ t('admin.channels.form.noTiersYet', '暂无层级，点击添加配置按次计费价格') placeholderplaceholder
      </div>
    </div>

    <!-- Image mode (legacy per-request) -->
    <div v-else-if="entry.billing_mode === 'image'">
      <div class="flex items-center justify-between">
        <label class="text-xs font-medium text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.form.imageTiers', '图片计费层级（按次）') placeholderplaceholder
        </label>
        <button type="button" @click="addImageTier" class="text-xs text-primary-600 hover:text-primary-700">
          + {{ t('admin.channels.form.addTier', '添加层级') placeholderplaceholder
        </button>
      </div>
      <div v-if="entry.intervals && entry.intervals.length > 0" class="mt-2 space-y-2">
        <IntervalRow
          v-for="(iv, idx) in entry.intervals"
          :key="idx"
          :interval="iv"
          :mode="entry.billing_mode"
          @update="updateInterval(idx, $event)"
          @remove="removeInterval(idx)"
        />
      </div>
      <div v-else>
        <div class="mt-2 grid grid-cols-2 gap-2 sm:grid-cols-4">
          <div>
            <label class="text-xs text-gray-400">{{ t('admin.channels.form.imageOutputPrice', '图片输出价格') placeholderplaceholder</label>
            <input :value="entry.image_output_price" @input="emitField('image_output_price', ($event.target as HTMLInputElement).value)"
              type="number" step="any" min="0" class="input mt-0.5 text-sm" :placeholder="t('admin.channels.form.pricePlaceholder', '默认')" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import Select from '@/components/common/Select.vue'
import Icon from '@/components/icons/Icon.vue'
import IntervalRow from './IntervalRow.vue'
import ModelTagInput from './ModelTagInput.vue'
import type { PricingFormEntry, IntervalFormEntry placeholder from './types'
import type { BillingMode placeholder from '@/api/admin/channels'

const { t placeholder = useI18n()

const props = defineProps<{
  entry: PricingFormEntry
placeholder>()

const emit = defineEmits<{
  update: [entry: PricingFormEntry]
  remove: []
placeholder>()

const billingModeOptions = computed(() => [
  { value: 'token', label: 'Token' placeholder,
  { value: 'per_request', label: t('admin.channels.billingMode.perRequest', '按次') placeholder,
  { value: 'image', label: t('admin.channels.billingMode.image', '图片（按次）') placeholder
])

function emitField(field: keyof PricingFormEntry, value: string) {
  emit('update', { ...props.entry, [field]: value === '' ? null : value placeholder)
placeholder

function addInterval() {
  const intervals = [...(props.entry.intervals || [])]
  intervals.push({
    min_tokens: 0, max_tokens: null, tier_label: '',
    input_price: null, output_price: null, cache_write_price: null,
    cache_read_price: null, per_request_price: null,
    sort_order: intervals.length
  placeholder)
  emit('update', { ...props.entry, intervals placeholder)
placeholder

function addImageTier() {
  const intervals = [...(props.entry.intervals || [])]
  const labels = ['1K', '2K', '4K', 'HD']
  intervals.push({
    min_tokens: 0, max_tokens: null, tier_label: labels[intervals.length] || '',
    input_price: null, output_price: null, cache_write_price: null,
    cache_read_price: null, per_request_price: null,
    sort_order: intervals.length
  placeholder)
  emit('update', { ...props.entry, intervals placeholder)
placeholder

function updateInterval(idx: number, updated: IntervalFormEntry) {
  const intervals = [...(props.entry.intervals || [])]
  intervals[idx] = updated
  emit('update', { ...props.entry, intervals placeholder)
placeholder

function removeInterval(idx: number) {
  const intervals = [...(props.entry.intervals || [])]
  intervals.splice(idx, 1)
  emit('update', { ...props.entry, intervals placeholder)
placeholder
</script>
