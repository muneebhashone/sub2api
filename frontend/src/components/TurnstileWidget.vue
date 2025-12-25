<template>
  <div v-if="siteKey" class="turnstile-wrapper">
    <div ref="containerRef" class="turnstile-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch placeholder from 'vue'

interface TurnstileRenderOptions {
  sitekey: string
  callback: (token: string) => void
  'expired-callback'?: () => void
  'error-callback'?: () => void
  theme?: 'light' | 'dark' | 'auto'
  size?: 'normal' | 'compact' | 'flexible'
placeholder

interface TurnstileAPI {
  render: (container: HTMLElement, options: TurnstileRenderOptions) => string
  reset: (widgetId?: string) => void
  remove: (widgetId?: string) => void
placeholder

declare global {
  interface Window {
    turnstile?: TurnstileAPI
    onTurnstileLoad?: () => void
  placeholder
placeholder

const props = withDefaults(
  defineProps<{
    siteKey: string
    theme?: 'light' | 'dark' | 'auto'
    size?: 'normal' | 'compact' | 'flexible'
  placeholder>(),
  {
    theme: 'auto',
    size: 'flexible'
  placeholder
)

const emit = defineEmits<{
  (e: 'verify', token: string): void
  (e: 'expire'): void
  (e: 'error'): void
placeholder>()

const containerRef = ref<HTMLElement | null>(null)
const widgetId = ref<string | null>(null)
const scriptLoaded = ref(false)

const loadScript = (): Promise<void> => {
  return new Promise((resolve, reject) => {
    if (window.turnstile) {
      scriptLoaded.value = true
      resolve()
      return
    placeholder

    // Check if script is already loading
    const existingScript = document.querySelector('script[src*="turnstile"]')
    if (existingScript) {
      window.onTurnstileLoad = () => {
        scriptLoaded.value = true
        resolve()
      placeholder
      return
    placeholder

    const script = document.createElement('script')
    script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js?onload=onTurnstileLoad'
    script.async = true
    script.defer = true

    window.onTurnstileLoad = () => {
      scriptLoaded.value = true
      resolve()
    placeholder

    script.onerror = () => {
      reject(new Error('Failed to load Turnstile script'))
    placeholder

    document.head.appendChild(script)
  placeholder)
placeholder

const renderWidget = () => {
  if (!window.turnstile || !containerRef.value || !props.siteKey) {
    return
  placeholder

  // Remove existing widget if any
  if (widgetId.value) {
    try {
      window.turnstile.remove(widgetId.value)
    placeholder catch {
      // Ignore errors when removing
    placeholder
    widgetId.value = null
  placeholder

  // Clear container
  containerRef.value.innerHTML = ''

  widgetId.value = window.turnstile.render(containerRef.value, {
    sitekey: props.siteKey,
    callback: (token: string) => {
      emit('verify', token)
    placeholder,
    'expired-callback': () => {
      emit('expire')
    placeholder,
    'error-callback': () => {
      emit('error')
    placeholder,
    theme: props.theme,
    size: props.size
  placeholder)
placeholder

const reset = () => {
  if (window.turnstile && widgetId.value) {
    window.turnstile.reset(widgetId.value)
  placeholder
placeholder

// Expose reset method to parent
defineExpose({ reset placeholder)

onMounted(async () => {
  if (!props.siteKey) {
    return
  placeholder

  try {
    await loadScript()
    renderWidget()
  placeholder catch (error) {
    console.error('Failed to initialize Turnstile:', error)
    emit('error')
  placeholder
placeholder)

onUnmounted(() => {
  if (window.turnstile && widgetId.value) {
    try {
      window.turnstile.remove(widgetId.value)
    placeholder catch {
      // Ignore errors when removing
    placeholder
  placeholder
placeholder)

// Re-render when siteKey changes
watch(
  () => props.siteKey,
  (newKey) => {
    if (newKey && scriptLoaded.value) {
      renderWidget()
    placeholder
  placeholder
)
</script>

<style scoped>
.turnstile-wrapper {
  width: 100%;
placeholder

.turnstile-container {
  width: 100%;
  min-height: 65px;
placeholder

/* Make the Turnstile iframe fill the container width */
.turnstile-container :deep(iframe) {
  width: 100% !important;
placeholder
</style>
