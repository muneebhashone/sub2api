<template>
  <div class="space-y-4">
    <button type="button" :disabled="disabled" class="btn btn-secondary w-full" @click="startLogin">
      <span
        class="mr-2 inline-flex h-5 w-5 items-center justify-center rounded-full bg-primary-100 text-xs font-semibold text-primary-700 dark:bg-primary-900/30 dark:text-primary-300"
      >
        {{ providerInitial placeholderplaceholder
      </span>
      {{ t('auth.oidc.signIn', { providerName: normalizedProviderName placeholder) placeholderplaceholder
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
import { computed placeholder from 'vue'
import { useRoute placeholder from 'vue-router'
import { useI18n placeholder from 'vue-i18n'
import type { OAuthLoginStart placeholder from '@/api/auth'
import { resolveAffiliateReferralCode, storeOAuthAffiliateCode placeholder from '@/utils/oauthAffiliate'

const props = withDefaults(defineProps<{
  disabled?: boolean
  affCode?: string
  providerName?: string
  showDivider?: boolean
placeholder>(), {
  providerName: 'OIDC',
  showDivider: true
placeholder)
const emit = defineEmits<{
  start: [request: OAuthLoginStart]
placeholder>()

const route = useRoute()
const { t placeholder = useI18n()

const normalizedProviderName = computed(() => {
  const name = props.providerName?.trim()
  return name || 'OIDC'
placeholder)

const providerInitial = computed(() => normalizedProviderName.value.charAt(0).toUpperCase() || 'O')

function startLogin(): void {
  const redirectTo = (route.query.redirect as string) || '/dashboard'
  storeOAuthAffiliateCode(resolveAffiliateReferralCode(props.affCode, route.query.aff, route.query.aff_code))
  emit('start', { provider: 'oidc', params: { redirect: redirectTo placeholder placeholder)
placeholder
</script>
