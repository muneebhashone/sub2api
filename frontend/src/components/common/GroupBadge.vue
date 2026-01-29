<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 rounded-md px-2 py-0.5 text-xs font-medium transition-colors',
      badgeClass
    ]"
  >
    <!-- Platform logo -->
    <PlatformIcon v-if="platform" :platform="platform" size="sm" />
    <!-- Group name -->
    <span class="truncate">{{ name placeholderplaceholder</span>
    <!-- Right side label -->
    <span v-if="showLabel" :class="labelClass">
      {{ labelText placeholderplaceholder
    </span>
  </span>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import type { SubscriptionType, GroupPlatform placeholder from '@/types'
import PlatformIcon from './PlatformIcon.vue'

interface Props {
  name: string
  platform?: GroupPlatform
  subscriptionType?: SubscriptionType
  rateMultiplier?: number
  showRate?: boolean
  daysRemaining?: number | null // 剩余天数（订阅类型时使用）
placeholder

const props = withDefaults(defineProps<Props>(), {
  subscriptionType: 'standard',
  showRate: true,
  daysRemaining: null
placeholder)

const { t placeholder = useI18n()

const isSubscription = computed(() => props.subscriptionType === 'subscription')

// 是否显示右侧标签
const showLabel = computed(() => {
  if (!props.showRate) return false
  // 订阅类型：显示天数或"订阅"
  if (isSubscription.value) return true
  // 标准类型：显示倍率
  return props.rateMultiplier !== undefined
placeholder)

// Label text
const labelText = computed(() => {
  if (isSubscription.value) {
    // 如果有剩余天数，显示天数
    if (props.daysRemaining !== null && props.daysRemaining !== undefined) {
      if (props.daysRemaining <= 0) {
        return t('admin.users.expired')
      placeholder
      return t('admin.users.daysRemaining', { days: props.daysRemaining placeholder)
    placeholder
    // 否则显示"订阅"
    return t('groups.subscription')
  placeholder
  return props.rateMultiplier !== undefined ? `${props.rateMultiplierplaceholderx` : ''
placeholder)

// Label style based on type and days remaining
const labelClass = computed(() => {
  const base = 'px-1.5 py-0.5 rounded text-[10px] font-semibold'

  if (!isSubscription.value) {
    // Standard: subtle background
    return `${baseplaceholder bg-black/10 dark:bg-white/10`
  placeholder

  // 订阅类型：根据剩余天数显示不同颜色
  if (props.daysRemaining !== null && props.daysRemaining !== undefined) {
    if (props.daysRemaining <= 0 || props.daysRemaining <= 3) {
      // 已过期或紧急（<=3天）：红色
      return `${baseplaceholder bg-red-200/80 text-red-800 dark:bg-red-800/50 dark:text-red-300`
    placeholder
    if (props.daysRemaining <= 7) {
      // 警告（<=7天）：橙色
      return `${baseplaceholder bg-amber-200/80 text-amber-800 dark:bg-amber-800/50 dark:text-amber-300`
    placeholder
  placeholder

  // 正常状态或无天数：根据平台显示主题色
  if (props.platform === 'anthropic') {
    return `${baseplaceholder bg-orange-200/60 text-orange-800 dark:bg-orange-800/40 dark:text-orange-300`
  placeholder
  if (props.platform === 'openai') {
    return `${baseplaceholder bg-emerald-200/60 text-emerald-800 dark:bg-emerald-800/40 dark:text-emerald-300`
  placeholder
  if (props.platform === 'gemini') {
    return `${baseplaceholder bg-blue-200/60 text-blue-800 dark:bg-blue-800/40 dark:text-blue-300`
  placeholder
  if (props.platform === 'sora') {
    return `${baseplaceholder bg-rose-200/60 text-rose-800 dark:bg-rose-800/40 dark:text-rose-300`
  placeholder
  return `${baseplaceholder bg-violet-200/60 text-violet-800 dark:bg-violet-800/40 dark:text-violet-300`
placeholder)

// Badge color based on platform and subscription type
const badgeClass = computed(() => {
  if (props.platform === 'anthropic') {
    // Claude: orange theme
    return isSubscription.value
      ? 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400'
      : 'bg-amber-50 text-amber-700 dark:bg-amber-900/20 dark:text-amber-400'
  placeholder else if (props.platform === 'openai') {
    // OpenAI: green theme
    return isSubscription.value
      ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
      : 'bg-green-50 text-green-700 dark:bg-green-900/20 dark:text-green-400'
  placeholder
  if (props.platform === 'gemini') {
    return isSubscription.value
      ? 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'
      : 'bg-sky-50 text-sky-700 dark:bg-sky-900/20 dark:text-sky-400'
  placeholder
  if (props.platform === 'sora') {
    return isSubscription.value
      ? 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-400'
      : 'bg-rose-50 text-rose-700 dark:bg-rose-900/20 dark:text-rose-400'
  placeholder
  // Fallback: original colors
  return isSubscription.value
    ? 'bg-violet-100 text-violet-700 dark:bg-violet-900/30 dark:text-violet-400'
    : 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400'
placeholder)
</script>
