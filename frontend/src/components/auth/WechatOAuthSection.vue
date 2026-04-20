<template>
  <div class="space-y-4">
    <button type="button" :disabled="disabled" class="btn btn-secondary w-full" @click="startLogin">
      <span
        class="mr-2 inline-flex h-5 w-5 items-center justify-center rounded-full bg-green-100 text-xs font-semibold text-green-700 dark:bg-green-900/30 dark:text-green-300"
      >
        W
      </span>
      {{ t('auth.oidc.signIn', { providerName placeholder) placeholderplaceholder
    </button>

    <div v-if="showDivider" class="flex items-center gap-3">
      <div class="h-px flex-1 bg-gray-200 dark:bg-dark-700"></div>
      <span class="text-xs text-gray-500 dark:text-dark-400">
        {{ t('auth.oauthOrContinue') placeholderplaceholder
      </span>
      <div class="h-px flex-1 bg-gray-200 dark:bg-dark-700"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute placeholder from 'vue-router'
import { useI18n placeholder from 'vue-i18n'

withDefaults(defineProps<{
  disabled?: boolean
  showDivider?: boolean
placeholder>(), {
  showDivider: true,
placeholder)

const route = useRoute()
const { t placeholder = useI18n()

const providerName = 'WeChat'

function resolveWeChatOAuthMode(): 'open' | 'mp' {
  if (typeof navigator === 'undefined') {
    return 'open'
  placeholder
  return /MicroMessenger/i.test(navigator.userAgent) ? 'mp' : 'open'
placeholder

function startLogin(): void {
  const redirectTo = (route.query.redirect as string) || '/dashboard'
  const apiBase = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api/v1'
  const normalized = apiBase.replace(/\/$/, '')
  const mode = resolveWeChatOAuthMode()
  const startURL = `${normalizedplaceholder/auth/oauth/wechat/start?mode=${modeplaceholder&redirect=${encodeURIComponent(redirectTo)placeholder`
  window.location.href = startURL
placeholder
</script>
