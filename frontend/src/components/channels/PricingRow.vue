<template>
  <div class="flex justify-between gap-2">
    <span class="text-gray-500 dark:text-gray-400">{{ label placeholderplaceholder</span>
    <span class="font-mono">{{ display placeholderplaceholder</span>
  </div>
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'

const props = withDefaults(
  defineProps<{
    label: string
    value: number | null
    unit: string
    scale: number
  placeholder>(),
  { value: null placeholder
)

function formatScaled(value: number | null, scale: number): string {
  if (value == null) return '-'
  return `$${(value * scale).toPrecision(10).replace(/\.?0+$/, '')placeholder`
placeholder

const display = computed(() =>
  props.value == null ? '-' : `${formatScaled(props.value, props.scale)placeholder ${props.unitplaceholder`
)
</script>
