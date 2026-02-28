<template>
  <AppLayout>
    <div class="purchase-page-layout">
      <div class="card flex-1 min-h-0 overflow-hidden">
        <div v-if="loading" class="flex h-full items-center justify-center py-12">
          <div
            class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"
          ></div>
        </div>

        <div
          v-else-if="!purchaseEnabled"
          class="flex h-full items-center justify-center p-10 text-center"
        >
          <div class="max-w-md">
            <div
              class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-dark-700"
            >
              <Icon name="creditCard" size="lg" class="text-gray-400" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('purchase.notEnabledTitle') placeholderplaceholder
            </h3>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              {{ t('purchase.notEnabledDesc') placeholderplaceholder
            </p>
          </div>
        </div>

        <div
          v-else-if="!isValidUrl"
          class="flex h-full items-center justify-center p-10 text-center"
        >
          <div class="max-w-md">
            <div
              class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-dark-700"
            >
              <Icon name="link" size="lg" class="text-gray-400" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('purchase.notConfiguredTitle') placeholderplaceholder
            </h3>
            <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
              {{ t('purchase.notConfiguredDesc') placeholderplaceholder
            </p>
          </div>
        </div>

        <div v-else class="purchase-embed-shell">
          <a
            :href="purchaseUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="btn btn-secondary btn-sm purchase-open-fab"
          >
            <Icon name="externalLink" size="sm" class="mr-1.5" :stroke-width="2" />
            {{ t('purchase.openInNewTab') placeholderplaceholder
          </a>
          <iframe
            :src="purchaseUrl"
            class="purchase-embed-frame"
            allowfullscreen
          ></iframe>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores'
import { useAuthStore placeholder from '@/stores/auth'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'

const { t placeholder = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const PURCHASE_USER_ID_QUERY_KEY = 'user_id'
const PURCHASE_THEME_QUERY_KEY = 'theme'
const PURCHASE_UI_MODE_QUERY_KEY = 'ui_mode'
const PURCHASE_UI_MODE_EMBEDDED = 'embedded'

const loading = ref(false)
const purchaseTheme = ref<'light' | 'dark'>('light')
let themeObserver: MutationObserver | null = null

const purchaseEnabled = computed(() => {
  return appStore.cachedPublicSettings?.purchase_subscription_enabled ?? false
placeholder)

function detectTheme(): 'light' | 'dark' {
  if (typeof document === 'undefined') return 'light'
  return document.documentElement.classList.contains('dark') ? 'dark' : 'light'
placeholder

function buildPurchaseUrl(baseUrl: string, userId?: number, theme: 'light' | 'dark' = 'light'): string {
  if (!baseUrl) return baseUrl
  try {
    const url = new URL(baseUrl)
    if (userId) {
      url.searchParams.set(PURCHASE_USER_ID_QUERY_KEY, String(userId))
    placeholder
    url.searchParams.set(PURCHASE_THEME_QUERY_KEY, theme)
    url.searchParams.set(PURCHASE_UI_MODE_QUERY_KEY, PURCHASE_UI_MODE_EMBEDDED)
    return url.toString()
  placeholder catch {
    const params: string[] = []
    if (userId) {
      params.push(`${PURCHASE_USER_ID_QUERY_KEYplaceholder=${encodeURIComponent(String(userId))placeholder`)
    placeholder
    params.push(`${PURCHASE_THEME_QUERY_KEYplaceholder=${encodeURIComponent(theme)placeholder`)
    params.push(`${PURCHASE_UI_MODE_QUERY_KEYplaceholder=${encodeURIComponent(PURCHASE_UI_MODE_EMBEDDED)placeholder`)
    const separator = baseUrl.includes('?') ? '&' : '?'
    return `${baseUrlplaceholder${separatorplaceholder${params.join('&')placeholder`
  placeholder
placeholder

const purchaseUrl = computed(() => {
  const baseUrl = (appStore.cachedPublicSettings?.purchase_subscription_url || '').trim()
  return buildPurchaseUrl(baseUrl, authStore.user?.id, purchaseTheme.value)
placeholder)

const isValidUrl = computed(() => {
  const url = purchaseUrl.value
  return url.startsWith('http://') || url.startsWith('https://')
placeholder)

onMounted(async () => {
  purchaseTheme.value = detectTheme()

  if (typeof document !== 'undefined') {
    themeObserver = new MutationObserver(() => {
      purchaseTheme.value = detectTheme()
    placeholder)
    themeObserver.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class'],
    placeholder)
  placeholder

  if (appStore.publicSettingsLoaded) return
  loading.value = true
  try {
    await appStore.fetchPublicSettings()
  placeholder finally {
    loading.value = false
  placeholder
placeholder)

onUnmounted(() => {
  if (themeObserver) {
    themeObserver.disconnect()
    themeObserver = null
  placeholder
placeholder)
</script>

<style scoped>
.purchase-page-layout {
  @apply flex flex-col;
  height: calc(100vh - 64px - 4rem);
placeholder

.purchase-embed-shell {
  @apply relative;
  @apply h-full w-full overflow-auto rounded-2xl;
  @apply bg-gradient-to-b from-gray-50 to-white dark:from-dark-900 dark:to-dark-950;
  @apply p-3 sm:p-4;
placeholder

.purchase-open-fab {
  @apply absolute right-3 top-3 z-10;
  @apply shadow-sm backdrop-blur supports-[backdrop-filter]:bg-white/80;
placeholder

.purchase-embed-frame {
  display: block;
  margin: 0 auto;
  width: min(100%, 440px);
  height: 840px;
  border: 0;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.14);
  background: transparent;
placeholder

@media (max-width: 640px) {
  .purchase-embed-frame {
    width: 100%;
    height: 780px;
    border-radius: 12px;
  placeholder
placeholder
</style>
