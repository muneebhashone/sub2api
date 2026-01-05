<template>
  <div
    :class="[
      'animate-pulse bg-gray-200 dark:bg-dark-700',
      variant === 'circle' ? 'rounded-full' : 'rounded-lg',
      customClass
    ]"
    :style="style"
  ></div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'

interface Props {
  variant?: 'rect' | 'circle' | 'text'
  width?: string | number
  height?: string | number
  class?: string
placeholder

const props = withDefaults(defineProps<Props>(), {
  variant: 'rect',
  width: '100%'
placeholder)

const customClass = computed(() => props.class || '')

const style = computed(() => {
  const s: Record<string, string> = {placeholder
  
  if (props.width) {
    s.width = typeof props.width === 'number' ? `${props.widthplaceholderpx` : props.width
  placeholder
  
  if (props.height) {
    s.height = typeof props.height === 'number' ? `${props.heightplaceholderpx` : props.height
  placeholder else if (props.variant === 'text') {
    s.height = '1em'
    s.marginTop = '0.25em'
    s.marginBottom = '0.25em'
  placeholder
  
  return s
placeholder)
</script>
