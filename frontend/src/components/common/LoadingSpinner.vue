<template>
  <div
    :class="['spinner', sizeClasses, colorClass]"
    role="status"
    :aria-label="t('common.loading')"
  >
    <span class="sr-only">{{ t('common.loading') placeholderplaceholder</span>
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'

const { t placeholder = useI18n()

type SpinnerSize = 'sm' | 'md' | 'lg' | 'xl'
type SpinnerColor = 'primary' | 'secondary' | 'white' | 'gray'

interface Props {
  size?: SpinnerSize
  color?: SpinnerColor
placeholder

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'primary'
placeholder)

const sizeClasses = computed(() => {
  const sizes: Record<SpinnerSize, string> = {
    sm: 'w-4 h-4 border-2',
    md: 'w-8 h-8 border-2',
    lg: 'w-12 h-12 border-[3px]',
    xl: 'w-16 h-16 border-4'
  placeholder
  return sizes[props.size]
placeholder)

const colorClass = computed(() => {
  const colors: Record<SpinnerColor, string> = {
    primary: 'text-primary-500',
    secondary: 'text-gray-500 dark:text-dark-400',
    white: 'text-white',
    gray: 'text-gray-400 dark:text-dark-500'
  placeholder
  return colors[props.color]
placeholder)
</script>

<style scoped>
.spinner {
  @apply inline-block rounded-full border-solid border-current border-r-transparent;
  animation: spin 0.75s linear infinite;
placeholder

@keyframes spin {
  from {
    transform: rotate(0deg);
  placeholder
  to {
    transform: rotate(360deg);
  placeholder
placeholder
</style>
