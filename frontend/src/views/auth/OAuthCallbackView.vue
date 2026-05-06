<template>
  <div class="min-h-screen bg-gray-50 px-4 py-10 dark:bg-dark-900">
    <div class="mx-auto max-w-2xl">
      <div v-if="isProcessing" class="card p-6 text-center">
        <div class="mx-auto h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent"></div>
        <h1 class="mt-4 text-lg font-semibold text-gray-900 dark:text-white">
          {{ t('auth.oauth.callbackTitle') placeholderplaceholder
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ t('auth.oauth.callbackHint') placeholderplaceholder
        </p>
      </div>

      <div v-else-if="needsInvitation" class="card p-6">
        <h1 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ t('auth.oidc.callbackTitle', { providerName placeholder) placeholderplaceholder
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ t('auth.oidc.invitationRequired', { providerName placeholder) placeholderplaceholder
        </p>

        <div class="mt-6 space-y-4">
          <input
            v-model="invitationCode"
            type="text"
            class="input w-full"
            :placeholder="t('auth.invitationCodePlaceholder')"
            :disabled="isSubmitting"
            @keyup.enter="handleSubmitInvitation"
          />
          <p v-if="invitationError" class="text-sm text-red-600 dark:text-red-400">
            {{ invitationError placeholderplaceholder
          </p>
          <button
            class="btn btn-primary w-full"
            type="button"
            :disabled="isSubmitting || !invitationCode.trim()"
            @click="handleSubmitInvitation"
          >
            {{ isSubmitting ? t('common.processing') : t('auth.oidc.completeRegistration') placeholderplaceholder
          </button>
        </div>
      </div>

      <div v-else-if="invalidCallback" class="card p-6 text-center">
        <h1 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ t('auth.oauth.invalidCallbackTitle') placeholderplaceholder
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ t('auth.oauth.invalidCallbackHint') placeholderplaceholder
        </p>
        <button class="btn btn-primary mt-6" type="button" @click="router.replace('/login')">
          {{ t('auth.backToLogin') placeholderplaceholder
        </button>
      </div>

      <div v-else class="card p-6">
        <h1 class="text-lg font-semibold text-gray-900 dark:text-white">
          {{ t('auth.oauth.callbackTitle') placeholderplaceholder
        </h1>
        <p class="mt-2 text-sm text-gray-600 dark:text-gray-400">
          {{ t('auth.oauth.callbackHint') placeholderplaceholder
        </p>

        <div class="mt-6 space-y-4">
          <div>
            <label class="input-label">{{ t('auth.oauth.code') placeholderplaceholder</label>
            <div class="flex gap-2">
              <input class="input flex-1 font-mono text-sm" :value="code" readonly />
              <button class="btn btn-secondary" type="button" :disabled="!code" @click="copy(code)">
                {{ t('common.copy') placeholderplaceholder
              </button>
            </div>
          </div>

          <div>
            <label class="input-label">{{ t('auth.oauth.state') placeholderplaceholder</label>
            <div class="flex gap-2">
              <input class="input flex-1 font-mono text-sm" :value="state" readonly />
              <button
                class="btn btn-secondary"
                type="button"
                :disabled="!state"
                @click="copy(state)"
              >
                {{ t('common.copy') placeholderplaceholder
              </button>
            </div>
          </div>

          <div>
            <label class="input-label">{{ t('auth.oauth.fullUrl') placeholderplaceholder</label>
            <div class="flex gap-2">
              <input class="input flex-1 font-mono text-xs" :value="fullUrl" readonly />
              <button
                class="btn btn-secondary"
                type="button"
                :disabled="!fullUrl"
                @click="copy(fullUrl)"
              >
                {{ t('common.copy') placeholderplaceholder
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRoute, useRouter placeholder from 'vue-router'
import { useClipboard placeholder from '@/composables/useClipboard'
import { useAppStore, useAuthStore placeholder from '@/stores'
import { apiClient placeholder from '@/api/client'
import {
  exchangePendingOAuthCompletion,
  persistOAuthTokenContext,
  type OAuthTokenResponse
placeholder from '@/api/auth'
import {
  clearAllAffiliateReferralCodes,
  loadOAuthAffiliateCode,
  oauthAffiliatePayload
placeholder from '@/utils/oauthAffiliate'

const route = useRoute()
const router = useRouter()
const { t placeholder = useI18n()
const { copyToClipboard placeholder = useClipboard()
const appStore = useAppStore()
const authStore = useAuthStore()
const isProcessing = ref(false)
const isSubmitting = ref(false)
const needsInvitation = ref(false)
const invitationCode = ref('')
const invitationError = ref('')
const pendingProvider = ref<'github' | 'google'>('github')
const redirectTo = ref('/dashboard')
const invalidCallback = ref(false)

type EmailOAuthPendingCompletion = Partial<OAuthTokenResponse> & {
  error?: string
  provider?: string
  redirect?: string
placeholder

const code = computed(() => (route.query.code as string) || '')
const state = computed(() => (route.query.state as string) || '')
const error = computed(
  () => (route.query.error as string) || (route.query.error_description as string) || ''
)

const fullUrl = computed(() => {
  if (typeof window === 'undefined') return ''
  return window.location.href
placeholder)
const providerName = computed(() =>
  pendingProvider.value === 'google' ? 'Google' : 'GitHub'
)

function parseFragmentParams(): URLSearchParams {
  const raw = typeof window !== 'undefined' ? window.location.hash : ''
  const hash = raw.startsWith('#') ? raw.slice(1) : raw
  return new URLSearchParams(hash)
placeholder

function readTokenResponse(params: URLSearchParams): OAuthTokenResponse | null {
  const accessToken = params.get('access_token')?.trim() || ''
  if (!accessToken) return null

  const response: OAuthTokenResponse = { access_token: accessToken placeholder
  const refreshToken = params.get('refresh_token')?.trim() || ''
  if (refreshToken) response.refresh_token = refreshToken
  const expiresIn = Number.parseInt(params.get('expires_in')?.trim() || '', 10)
  if (Number.isFinite(expiresIn) && expiresIn > 0) response.expires_in = expiresIn
  const tokenType = params.get('token_type')?.trim() || ''
  if (tokenType) response.token_type = tokenType
  return response
placeholder

function sanitizeRedirectPath(path: string | null | undefined): string {
  if (!path) return '/dashboard'
  if (!path.startsWith('/')) return '/dashboard'
  if (path.startsWith('//')) return '/dashboard'
  if (path.includes('://')) return '/dashboard'
  if (path.includes('\n') || path.includes('\r')) return '/dashboard'
  return path
placeholder

async function finalizeTokenResponse(tokenResponse: OAuthTokenResponse, redirect: string) {
  persistOAuthTokenContext(tokenResponse)
  await authStore.setToken(tokenResponse.access_token)
  clearAllAffiliateReferralCodes()
  appStore.showSuccess(t('auth.loginSuccess'))
  await router.replace(sanitizeRedirectPath(redirect))
placeholder

function hasOAuthTokenResponse(value: Partial<OAuthTokenResponse>): value is OAuthTokenResponse {
  return typeof value.access_token === 'string' && value.access_token.trim() !== ''
placeholder

async function resumePendingEmailOAuth() {
  isProcessing.value = true
  try {
    const completion = await exchangePendingOAuthCompletion() as EmailOAuthPendingCompletion
    const completionRedirect = completion.redirect || '/dashboard'
    if (hasOAuthTokenResponse(completion)) {
      await finalizeTokenResponse(completion, completionRedirect)
      return
    placeholder

    const provider = String(completion.provider || '').toLowerCase()
    if (provider === 'github' || provider === 'google') {
      pendingProvider.value = provider
    placeholder
    redirectTo.value = sanitizeRedirectPath(completionRedirect)

    if (completion.error === 'invitation_required') {
      needsInvitation.value = true
      isProcessing.value = false
      return
    placeholder

    appStore.showError(completion.error || t('auth.loginFailed'))
  placeholder catch (e: unknown) {
    const err = e as { message?: string; response?: { data?: { message?: string placeholder placeholder placeholder
    const message = err.response?.data?.message || err.message || t('auth.loginFailed')
    appStore.showError(message)
    invalidCallback.value = true
  placeholder finally {
    if (!needsInvitation.value) {
      isProcessing.value = false
    placeholder
  placeholder
placeholder

async function handleSubmitInvitation() {
  invitationError.value = ''
  const code = invitationCode.value.trim()
  if (!code) return

  isSubmitting.value = true
  try {
    const { data placeholder = await apiClient.post<OAuthTokenResponse>(
      `/auth/oauth/${pendingProvider.valueplaceholder/complete-registration`,
      {
        invitation_code: code,
        ...oauthAffiliatePayload(loadOAuthAffiliateCode())
      placeholder
    )
    await finalizeTokenResponse(data, redirectTo.value)
  placeholder catch (e: unknown) {
    const err = e as { message?: string; response?: { data?: { message?: string placeholder placeholder placeholder
    invitationError.value =
      err.response?.data?.message || err.message || t('auth.oidc.completeRegistrationFailed')
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

onMounted(async () => {
  const params = parseFragmentParams()
  const tokenResponse = readTokenResponse(params)
  const fragmentError = params.get('error') || ''
  const fragmentErrorDescription =
    params.get('error_description') || params.get('error_message') || ''

  if (fragmentError) {
    appStore.showError(fragmentErrorDescription || fragmentError)
    return
  placeholder
  if (!tokenResponse) {
    if (route.path === '/auth/oauth/callback') {
      await resumePendingEmailOAuth()
    placeholder
    return
  placeholder

  isProcessing.value = true
  try {
    await finalizeTokenResponse(tokenResponse, params.get('redirect') || '/dashboard')
  placeholder catch (error: unknown) {
    const message = (error as { message?: string placeholder)?.message || t('auth.loginFailed')
    appStore.showError(message)
    isProcessing.value = false
  placeholder
placeholder)

watch(
  error,
  (message) => {
    if (message) {
      appStore.showError(message)
    placeholder
  placeholder,
  { immediate: true placeholder
)

const copy = (value: string) => {
  if (!value) return
  copyToClipboard(value)
placeholder
</script>
