<template>
  <AuthLayout>
    <div class="space-y-6">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.oidc.callbackTitle', { providerName placeholder) placeholderplaceholder
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{
            isProcessing
              ? t('auth.oidc.callbackProcessing', { providerName placeholder)
              : t('auth.oidc.callbackHint')
          placeholderplaceholder
        </p>
      </div>

      <transition name="fade">
        <div
          v-if="
            needsInvitation ||
            needsAdoptionConfirmation ||
            needsCreateAccount ||
            needsBindLogin ||
            needsTotpChallenge
          "
          class="space-y-4"
        >
          <div
            v-if="adoptionRequired && (suggestedDisplayName || suggestedAvatarUrl)"
            class="rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-600 dark:bg-dark-800/60"
          >
            <div class="space-y-3">
              <div class="space-y-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  Use {{ providerName placeholderplaceholder profile details
                </p>
                <p class="text-xs text-gray-500 dark:text-dark-400">
                  Choose whether to apply the nickname or avatar from {{ providerName placeholderplaceholder to this account.
                </p>
              </div>

              <label
                v-if="suggestedDisplayName"
                class="flex items-start gap-3 rounded-lg border border-gray-200 bg-white p-3 text-sm dark:border-dark-600 dark:bg-dark-900/50"
              >
                <input v-model="adoptDisplayName" type="checkbox" class="mt-1 h-4 w-4" />
                <span class="space-y-1">
                  <span class="block font-medium text-gray-900 dark:text-white">
                    Use display name
                  </span>
                  <span class="block text-gray-500 dark:text-dark-400">
                    {{ suggestedDisplayName placeholderplaceholder
                  </span>
                </span>
              </label>

              <label
                v-if="suggestedAvatarUrl"
                class="flex items-start gap-3 rounded-lg border border-gray-200 bg-white p-3 text-sm dark:border-dark-600 dark:bg-dark-900/50"
              >
                <input v-model="adoptAvatar" type="checkbox" class="mt-1 h-4 w-4" />
                <img
                  :src="suggestedAvatarUrl"
                  :alt="`${providerNameplaceholder avatar`"
                  class="h-10 w-10 rounded-full border border-gray-200 object-cover dark:border-dark-600"
                />
                <span class="space-y-1">
                  <span class="block font-medium text-gray-900 dark:text-white">
                    Use avatar
                  </span>
                  <span class="block break-all text-gray-500 dark:text-dark-400">
                    {{ suggestedAvatarUrl placeholderplaceholder
                  </span>
                </span>
              </label>
            </div>
          </div>

          <template v-if="needsInvitation">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              {{ t('auth.oidc.invitationRequired', { providerName placeholder) placeholderplaceholder
            </p>
            <div>
              <input
                v-model="invitationCode"
                type="text"
                class="input w-full"
                :placeholder="t('auth.invitationCodePlaceholder')"
                :disabled="isSubmitting"
                @keyup.enter="handleSubmitInvitation"
              />
            </div>
            <transition name="fade">
              <p v-if="invitationError" class="text-sm text-red-600 dark:text-red-400">
                {{ invitationError placeholderplaceholder
              </p>
            </transition>
            <button
              class="btn btn-primary w-full"
              :disabled="isSubmitting || !invitationCode.trim()"
              @click="handleSubmitInvitation"
            >
              {{
                isSubmitting
                  ? t('auth.oidc.completing')
                  : t('auth.oidc.completeRegistration')
              placeholderplaceholder
            </button>

            <div
              class="rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-600 dark:bg-dark-800/60"
            >
              <div class="space-y-3">
                <div class="space-y-1">
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ t('auth.alreadyHaveAccount') placeholderplaceholder
                  </p>
                  <p class="text-xs text-gray-500 dark:text-dark-400">
                    {{
                      hasCurrentAuthToken
                        ? 'Bind this WeChat identity to the account currently signed in on this browser.'
                        : 'Sign in to an existing account, then bind this WeChat identity to it.'
                    placeholderplaceholder
                  </p>
                </div>

                <input
                  v-if="!hasCurrentAuthToken"
                  v-model="existingAccountEmail"
                  data-testid="existing-account-email"
                  type="email"
                  class="input w-full"
                  :placeholder="t('auth.emailPlaceholder')"
                  :disabled="isSubmitting"
                />

                <button
                  data-testid="existing-account-submit"
                  type="button"
                  class="btn btn-secondary w-full"
                  :disabled="isSubmitting"
                  @click="handleExistingAccountBinding"
                >
                  {{ hasCurrentAuthToken ? 'Bind current account' : t('auth.signIn') placeholderplaceholder
                </button>
              </div>
            </div>
          </template>

          <template v-else-if="needsAdoptionConfirmation">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Review the {{ providerName placeholderplaceholder profile details before continuing.
            </p>
            <button class="btn btn-primary w-full" :disabled="isSubmitting" @click="handleContinueLogin">
              {{ isSubmitting ? t('common.processing') : 'Continue' placeholderplaceholder
            </button>
          </template>

          <template v-else-if="needsCreateAccount">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Enter an email address to create your account and continue.
            </p>
            <PendingOAuthCreateAccountForm
              test-id-prefix="wechat"
              :initial-email="pendingAccountEmail"
              :is-submitting="isSubmitting"
              :error-message="accountActionError"
              @submit="handleCreateAccount"
              @switch-to-bind="switchToBindLoginMode"
            />
          </template>

          <template v-else-if="needsBindLogin">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Log in to an existing account to bind this {{ providerName placeholderplaceholder sign-in.
            </p>
            <div class="space-y-3">
              <input
                v-model="bindLoginEmail"
                data-testid="wechat-bind-login-email"
                type="email"
                class="input w-full"
                placeholder="you@example.com"
                :disabled="isSubmitting"
                @keyup.enter="handleBindLogin"
              />
              <input
                v-model="bindLoginPassword"
                data-testid="wechat-bind-login-password"
                type="password"
                class="input w-full"
                placeholder="Password"
                :disabled="isSubmitting"
                @keyup.enter="handleBindLogin"
              />
              <button
                data-testid="wechat-bind-login-submit"
                class="btn btn-primary w-full"
                :disabled="isSubmitting || !bindLoginEmail.trim() || !bindLoginPassword"
                @click="handleBindLogin"
              >
                {{ isSubmitting ? t('common.processing') : 'Log in and bind' placeholderplaceholder
              </button>
              <button
                v-if="canReturnToCreateAccount"
                class="btn btn-secondary w-full"
                :disabled="isSubmitting"
                @click="switchToCreateAccountMode"
              >
                Use a different email
              </button>
            </div>
            <transition name="fade">
              <p v-if="accountActionError" class="text-sm text-red-600 dark:text-red-400">
                {{ accountActionError placeholderplaceholder
              </p>
            </transition>
          </template>

          <template v-else-if="needsTotpChallenge">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Enter the 6-digit verification code for
              <span class="font-medium">{{ totpUserEmailMasked || 'your account' placeholderplaceholder</span>
              to finish binding this {{ providerName placeholderplaceholder sign-in.
            </p>
            <div class="space-y-3">
              <input
                v-model="totpCode"
                data-testid="wechat-bind-login-totp"
                type="text"
                inputmode="numeric"
                maxlength="6"
                class="input w-full"
                placeholder="123456"
                :disabled="isSubmitting"
                @keyup.enter="handleSubmitTotpChallenge"
              />
              <button
                data-testid="wechat-bind-login-totp-submit"
                class="btn btn-primary w-full"
                :disabled="isSubmitting || totpCode.trim().length !== 6"
                @click="handleSubmitTotpChallenge"
              >
                {{ isSubmitting ? t('common.processing') : 'Verify and continue' placeholderplaceholder
              </button>
            </div>
            <transition name="fade">
              <p v-if="totpError" class="text-sm text-red-600 dark:text-red-400">
                {{ totpError placeholderplaceholder
              </p>
            </transition>
          </template>
        </div>
      </transition>

      <transition name="fade">
        <div
          v-if="errorMessage"
          class="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-800/50 dark:bg-red-900/20"
        >
          <div class="flex items-start gap-3">
            <div class="flex-shrink-0">
              <Icon name="exclamationCircle" size="md" class="text-red-500" />
            </div>
            <div class="space-y-2">
              <p class="text-sm text-red-700 dark:text-red-400">
                {{ errorMessage placeholderplaceholder
              </p>
              <router-link to="/login" class="btn btn-primary">
                {{ t('auth.oidc.backToLogin') placeholderplaceholder
              </router-link>
            </div>
          </div>
        </div>
      </transition>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref placeholder from 'vue'
import { useRoute, useRouter placeholder from 'vue-router'
import { useI18n placeholder from 'vue-i18n'
import { AuthLayout placeholder from '@/components/layout'
import PendingOAuthCreateAccountForm, {
  type PendingOAuthCreateAccountPayload
placeholder from '@/components/auth/PendingOAuthCreateAccountForm.vue'
import Icon from '@/components/icons/Icon.vue'
import { apiClient placeholder from '@/api/client'
import { useAuthStore, useAppStore placeholder from '@/stores'
import {
  completeWeChatOAuthRegistration,
  exchangePendingOAuthCompletion,
  getAuthToken,
  hasExplicitWeChatOAuthCapabilities,
  getOAuthCompletionKind,
  isOAuthLoginCompletion,
  login2FA,
  prepareOAuthBindAccessTokenCookie,
  persistOAuthTokenContext,
  resolveWeChatOAuthStartStrict,
  type OAuthAdoptionDecision,
  type OAuthTokenResponse,
  type PendingOAuthExchangeResponse
placeholder from '@/api/auth'

const route = useRoute()
const router = useRouter()
const { t placeholder = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const isProcessing = ref(true)
const errorMessage = ref('')
const needsInvitation = ref(false)
const invitationCode = ref('')
const isSubmitting = ref(false)
const invitationError = ref('')
const redirectTo = ref('/dashboard')
const adoptionRequired = ref(false)
const suggestedDisplayName = ref('')
const suggestedAvatarUrl = ref('')
const existingAccountEmail = ref('')
const adoptDisplayName = ref(true)
const adoptAvatar = ref(true)
const needsAdoptionConfirmation = ref(false)
const pendingAccountAction = ref<'none' | 'create_account' | 'bind_login'>('none')
const pendingAccountEmail = ref('')
const bindLoginEmail = ref('')
const bindLoginPassword = ref('')
const legacyPendingOAuthToken = ref('')
const accountActionError = ref('')
const canReturnToCreateAccount = ref(false)
const needsTotpChallenge = ref(false)
const totpTempToken = ref('')
const totpCode = ref('')
const totpError = ref('')
const totpUserEmailMasked = ref('')
const bindSuccessMessage = t('profile.authBindings.bindSuccess')

const providerName = 'WeChat'
const needsCreateAccount = computed(() => pendingAccountAction.value === 'create_account')
const needsBindLogin = computed(() => pendingAccountAction.value === 'bind_login')
const hasCurrentAuthToken = computed(() => Boolean(getAuthToken()))

type PendingWeChatCompletion = PendingOAuthExchangeResponse & {
  step?: string
  pending_email?: string
  resolved_email?: string
  existing_account_email?: string
  email?: string
  intent?: string
  requires_2fa?: boolean
  temp_token?: string
  user_email_masked?: string
placeholder

function persistPendingAuthSession(redirect?: string) {
  authStore.setPendingAuthSession({
    token: '',
    token_field: 'pending_oauth_token',
    provider: 'wechat',
    redirect: sanitizeRedirectPath(redirect || redirectTo.value)
  placeholder)
placeholder

function clearPendingAuthSession() {
  authStore.clearPendingAuthSession()
placeholder

function parseFragmentParams(): URLSearchParams {
  const raw = typeof window !== 'undefined' ? window.location.hash : ''
  const hash = raw.startsWith('#') ? raw.slice(1) : raw
  return new URLSearchParams(hash)
placeholder

function readLegacyFragmentLogin(params: URLSearchParams): OAuthTokenResponse | null {
  const accessToken = params.get('access_token')?.trim() || ''
  if (!accessToken) {
    return null
  placeholder

  const completion: OAuthTokenResponse = {
    access_token: accessToken
  placeholder
  const refreshToken = params.get('refresh_token')?.trim() || ''
  if (refreshToken) {
    completion.refresh_token = refreshToken
  placeholder
  const expiresIn = Number.parseInt(params.get('expires_in')?.trim() || '', 10)
  if (Number.isFinite(expiresIn) && expiresIn > 0) {
    completion.expires_in = expiresIn
  placeholder
  const tokenType = params.get('token_type')?.trim() || ''
  if (tokenType) {
    completion.token_type = tokenType
  placeholder
  return completion
placeholder

function sanitizeRedirectPath(path: string | null | undefined): string {
  if (!path) return '/dashboard'
  if (!path.startsWith('/')) return '/dashboard'
  if (path.startsWith('//')) return '/dashboard'
  if (path.includes('://')) return '/dashboard'
  if (path.includes('\n') || path.includes('\r')) return '/dashboard'
  return path
placeholder

async function ensurePublicSettingsLoaded(): Promise<void> {
  if (hasExplicitWeChatOAuthCapabilities(appStore.cachedPublicSettings)) {
    return
  placeholder

  if (appStore.publicSettingsLoaded) {
    return
  placeholder

  await appStore.fetchPublicSettings()
placeholder

function resolveConfiguredWeChatOAuthMode(): 'open' | 'mp' | null {
  if (!hasExplicitWeChatOAuthCapabilities(appStore.cachedPublicSettings)) {
    return null
  placeholder

  return resolveWeChatOAuthStartStrict(appStore.cachedPublicSettings).mode
placeholder

function resolveWeChatOAuthUnavailableMessage(): string {
  const resolved = resolveWeChatOAuthStartStrict(appStore.cachedPublicSettings)

  switch (resolved.unavailableReason) {
    case 'capability_unknown':
      return 'WeChat sign-in availability could not be confirmed. Refresh and retry.'
    case 'external_browser_required':
      return 'This WeChat sign-in flow is only available in your system browser.'
    case 'wechat_browser_required':
      return 'This WeChat sign-in flow is only available inside the WeChat browser.'
    case 'not_configured':
      return 'WeChat sign-in is not configured yet.'
    default:
      return t('auth.loginFailed')
  placeholder
placeholder

function resolveRuntimeWeChatOAuthMode(): 'open' | 'mp' {
  if (typeof navigator === 'undefined') {
    return 'open'
  placeholder
  return /MicroMessenger/i.test(navigator.userAgent) ? 'mp' : 'open'
placeholder

function normalizeWeChatOAuthMode(value: unknown): 'open' | 'mp' | null {
  return value === 'open' || value === 'mp' ? value : null
placeholder

function resolveRequestedWeChatOAuthMode(): 'open' | 'mp' | null {
  const configuredMode = resolveConfiguredWeChatOAuthMode()
  if (configuredMode) {
    return configuredMode
  placeholder

  const queryMode = normalizeWeChatOAuthMode(route.query.mode)
  if (queryMode) {
    return queryMode
  placeholder

  return resolveRuntimeWeChatOAuthMode()
placeholder

function resolveRedirectTarget(): string {
  return sanitizeRedirectPath(
    (route.query.redirect as string | undefined) || redirectTo.value || '/dashboard'
  )
placeholder

function resolveWeChatStartURL(intent: 'bind_current_user' | 'adopt_existing_user_by_email'): string | null {
  const apiBase = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api/v1'
  const normalized = apiBase.replace(/\/$/, '')
  const mode = resolveRequestedWeChatOAuthMode()
  if (!mode) {
    return null
  placeholder
  const params = new URLSearchParams({
    mode,
    redirect: resolveRedirectTarget(),
    intent,
  placeholder)

  const email = existingAccountEmail.value.trim()
  if (email) {
    params.set('email', email)
  placeholder

  return `${normalizedplaceholder/auth/oauth/wechat/start?${params.toString()placeholder`
placeholder

function buildExistingAccountResumePath(): string | null {
  const mode = resolveRequestedWeChatOAuthMode()
  if (!mode) {
    return null
  placeholder
  const params = new URLSearchParams({
    wechat_bind_existing: '1',
    redirect: resolveRedirectTarget(),
    mode,
  placeholder)

  const email = existingAccountEmail.value.trim()
  if (email) {
    params.set('email', email)
  placeholder

  return `/auth/wechat/callback?${params.toString()placeholder`
placeholder

function currentAdoptionDecision(): OAuthAdoptionDecision {
  return {
    adoptDisplayName: adoptDisplayName.value,
    adoptAvatar: adoptAvatar.value
  placeholder
placeholder

function resolveResumeEmail(): string {
  return typeof route.query.email === 'string' ? route.query.email.trim() : ''
placeholder

function serializeAdoptionDecision(decision: OAuthAdoptionDecision): Record<string, boolean> {
  const payload: Record<string, boolean> = {placeholder
  if (typeof decision.adoptDisplayName === 'boolean') {
    payload.adopt_display_name = decision.adoptDisplayName
  placeholder
  if (typeof decision.adoptAvatar === 'boolean') {
    payload.adopt_avatar = decision.adoptAvatar
  placeholder
  return payload
placeholder

async function handleExistingAccountBinding() {
  const unavailableMessage = resolveConfiguredWeChatOAuthMode() === null
    ? resolveWeChatOAuthUnavailableMessage()
    : ''

  if (getAuthToken()) {
    const startURL = resolveWeChatStartURL('bind_current_user')
    if (!startURL) {
      errorMessage.value = unavailableMessage || resolveWeChatOAuthUnavailableMessage()
      appStore.showError(errorMessage.value)
      return
    placeholder
    prepareOAuthBindAccessTokenCookie()
    window.location.href = startURL
    return
  placeholder

  const resumePath = buildExistingAccountResumePath()
  if (!resumePath) {
    errorMessage.value = unavailableMessage || resolveWeChatOAuthUnavailableMessage()
    appStore.showError(errorMessage.value)
    return
  placeholder

  const params = new URLSearchParams({
    redirect: resumePath,
  placeholder)
  const email = existingAccountEmail.value.trim()
  if (email) {
    params.set('email', email)
  placeholder
  await router.replace(`/login?${params.toString()placeholder`)
placeholder

function applyAdoptionSuggestionState(completion: PendingOAuthExchangeResponse) {
  adoptionRequired.value = completion.adoption_required === true
  suggestedDisplayName.value = completion.suggested_display_name || ''
  suggestedAvatarUrl.value = completion.suggested_avatar_url || ''

  if (!suggestedDisplayName.value) {
    adoptDisplayName.value = false
  placeholder
  if (!suggestedAvatarUrl.value) {
    adoptAvatar.value = false
  placeholder
placeholder

function hasSuggestedProfile(completion: PendingOAuthExchangeResponse): boolean {
  return Boolean(completion.suggested_display_name || completion.suggested_avatar_url)
placeholder

function normalizedPendingState(value: string | null | undefined): string {
  return value?.trim().toLowerCase() || ''
placeholder

function extractPendingAccountEmail(completion: PendingWeChatCompletion): string {
  return (
    completion.pending_email ||
    completion.existing_account_email ||
    completion.resolved_email ||
    completion.email ||
    resolveResumeEmail() ||
    ''
  ).trim()
placeholder

function resolvePendingAccountAction(
  completion: PendingWeChatCompletion
): 'none' | 'create_account' | 'bind_login' {
  const raw = normalizedPendingState(completion.step || completion.error || completion.intent)
  if (raw === 'email_required' || raw === 'create_account_required' || raw === 'create_account') {
    return 'create_account'
  placeholder
  if (
    raw === 'bind_login_required' ||
    raw === 'bind_login' ||
    raw === 'existing_account_binding_required' ||
    raw === 'existing_account_required' ||
    raw === 'adopt_existing_user_by_email'
  ) {
    return 'bind_login'
  placeholder
  return 'none'
placeholder

function applyPendingAccountAction(completion: PendingWeChatCompletion) {
  const action = resolvePendingAccountAction(completion)
  pendingAccountAction.value = action
  accountActionError.value = ''
  needsTotpChallenge.value = false
  totpTempToken.value = ''
  totpCode.value = ''
  totpError.value = ''
  totpUserEmailMasked.value = ''

  const email = extractPendingAccountEmail(completion)
  if (action === 'create_account') {
    pendingAccountEmail.value = email
    canReturnToCreateAccount.value = true
    return
  placeholder

  if (action === 'bind_login') {
    bindLoginEmail.value = email
    bindLoginPassword.value = ''
    canReturnToCreateAccount.value = true
    return
  placeholder

  canReturnToCreateAccount.value = false
placeholder

function applyTotpChallenge(completion: PendingWeChatCompletion): boolean {
  if (completion.requires_2fa !== true || !completion.temp_token) {
    return false
  placeholder

  pendingAccountAction.value = 'none'
  needsInvitation.value = false
  needsAdoptionConfirmation.value = false
  needsTotpChallenge.value = true
  totpTempToken.value = completion.temp_token
  totpCode.value = ''
  totpError.value = ''
  totpUserEmailMasked.value = completion.user_email_masked || ''
  isProcessing.value = false
  return true
placeholder

function switchToBindLoginMode(nextEmail?: string) {
  pendingAccountAction.value = 'bind_login'
  bindLoginEmail.value = bindLoginEmail.value.trim() || nextEmail?.trim() || pendingAccountEmail.value.trim()
  bindLoginPassword.value = ''
  accountActionError.value = ''
  canReturnToCreateAccount.value = true
placeholder

function switchToCreateAccountMode() {
  pendingAccountAction.value = 'create_account'
  pendingAccountEmail.value = pendingAccountEmail.value.trim() || bindLoginEmail.value.trim()
  accountActionError.value = ''
placeholder

function getRequestErrorMessage(error: unknown, fallback: string): string {
  const err = error as { message?: string; response?: { data?: { detail?: string; message?: string placeholder placeholder placeholder
  return err.response?.data?.detail || err.response?.data?.message || err.message || fallback
placeholder

function isCreateAccountRecoveryError(error: unknown): boolean {
  const data = (error as {
    response?: {
      data?: {
        reason?: string
        error?: string
        code?: string
        step?: string
        intent?: string
      placeholder
    placeholder
  placeholder).response?.data
  const states = [data?.reason, data?.error, data?.code, data?.step, data?.intent]
    .map(value => value?.trim().toLowerCase())
    .filter((value): value is string => Boolean(value))

  return states.includes('email_exists') ||
    states.includes('bind_login_required') ||
    states.includes('bind_login') ||
    states.includes('adopt_existing_user_by_email')
placeholder

async function finalizeCompletion(completion: PendingOAuthExchangeResponse, redirect: string) {
  if (getOAuthCompletionKind(completion) === 'bind') {
    const bindRedirect = sanitizeRedirectPath(completion.redirect || '/profile')
    clearPendingAuthSession()
    appStore.showSuccess(bindSuccessMessage)
    await router.replace(bindRedirect)
    return
  placeholder

  if (!isOAuthLoginCompletion(completion)) {
    throw new Error(t('auth.oidc.callbackMissingToken'))
  placeholder

  persistOAuthTokenContext(completion)
  await authStore.setToken(completion.access_token)
  appStore.showSuccess(t('auth.loginSuccess'))
  await router.replace(redirect)
placeholder

async function finalizePendingAccountResponse(completion: PendingWeChatCompletion) {
  applyAdoptionSuggestionState(completion)
  const redirect = sanitizeRedirectPath(completion.redirect || redirectTo.value)

  if (completion.error === 'invitation_required') {
    pendingAccountAction.value = 'none'
    needsInvitation.value = true
    needsAdoptionConfirmation.value = false
    isProcessing.value = false
    persistPendingAuthSession(redirect)
    return
  placeholder

  if (applyTotpChallenge(completion)) {
    persistPendingAuthSession(redirect)
    return
  placeholder

  applyPendingAccountAction(completion)
  if (pendingAccountAction.value !== 'none') {
    needsInvitation.value = false
    needsAdoptionConfirmation.value = false
    isProcessing.value = false
    persistPendingAuthSession(redirect)
    return
  placeholder

  await finalizeCompletion(completion, redirect)
placeholder

async function handleSubmitInvitation() {
  invitationError.value = ''
  if (!invitationCode.value.trim()) return

  isSubmitting.value = true
  try {
    const tokenData = legacyPendingOAuthToken.value
      ? (
          await apiClient.post<OAuthTokenResponse>('/auth/oauth/wechat/complete-registration', {
            pending_oauth_token: legacyPendingOAuthToken.value,
            invitation_code: invitationCode.value.trim(),
            ...serializeAdoptionDecision(currentAdoptionDecision())
          placeholder)
        ).data
      : await completeWeChatOAuthRegistration(
          invitationCode.value.trim(),
          currentAdoptionDecision()
        )
    persistOAuthTokenContext(tokenData)
    await authStore.setToken(tokenData.access_token)
    appStore.showSuccess(t('auth.loginSuccess'))
    await router.replace(redirectTo.value)
  placeholder catch (e: unknown) {
    const err = e as { message?: string; response?: { data?: { message?: string placeholder placeholder placeholder
    invitationError.value =
      err.response?.data?.message || err.message || t('auth.oidc.completeRegistrationFailed')
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

async function handleContinueLogin() {
  isSubmitting.value = true
  try {
    const completion = await exchangePendingOAuthCompletion(currentAdoptionDecision()) as PendingWeChatCompletion
    await finalizePendingAccountResponse(completion)
  placeholder catch (e: unknown) {
    errorMessage.value = getRequestErrorMessage(e, t('auth.loginFailed'))
    appStore.showError(errorMessage.value)
    needsAdoptionConfirmation.value = false
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

async function handleCreateAccount(payload: PendingOAuthCreateAccountPayload) {
  accountActionError.value = ''
  if (!payload.email || !payload.password) return

  isSubmitting.value = true
  try {
    const { data placeholder = await apiClient.post<PendingWeChatCompletion>('/auth/oauth/pending/create-account', {
      email: payload.email,
      password: payload.password,
      verify_code: payload.verifyCode || undefined,
      invitation_code: payload.invitationCode || undefined,
      ...serializeAdoptionDecision(currentAdoptionDecision())
    placeholder)
    await finalizePendingAccountResponse(data)
  placeholder catch (e: unknown) {
    if (isCreateAccountRecoveryError(e)) {
      switchToBindLoginMode(payload.email)
      return
    placeholder
    accountActionError.value = getRequestErrorMessage(e, t('auth.loginFailed'))
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

async function handleBindLogin() {
  accountActionError.value = ''
  const email = bindLoginEmail.value.trim()
  const password = bindLoginPassword.value
  if (!email || !password) return

  isSubmitting.value = true
  try {
    const { data placeholder = await apiClient.post<PendingWeChatCompletion>('/auth/oauth/pending/bind-login', {
      email,
      password,
      ...serializeAdoptionDecision(currentAdoptionDecision())
    placeholder)
    await finalizePendingAccountResponse(data)
  placeholder catch (e: unknown) {
    accountActionError.value = getRequestErrorMessage(e, t('auth.loginFailed'))
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

async function handleSubmitTotpChallenge() {
  totpError.value = ''
  const code = totpCode.value.trim()
  if (!totpTempToken.value || code.length !== 6) return

  isSubmitting.value = true
  try {
    const completion = await login2FA({
      temp_token: totpTempToken.value,
      totp_code: code
    placeholder)
    persistOAuthTokenContext(completion)
    await authStore.setToken(completion.access_token)
    appStore.showSuccess(t('auth.loginSuccess'))
    await router.replace(redirectTo.value)
  placeholder catch (e: unknown) {
    totpError.value = getRequestErrorMessage(e, t('auth.loginFailed'))
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

onMounted(async () => {
  try {
    await ensurePublicSettingsLoaded()
  placeholder catch {
    // Binding recovery requires confirmed capability flags. Use the strict guard below.
  placeholder

  if (typeof route.query.email === 'string') {
    existingAccountEmail.value = route.query.email
  placeholder

  if (route.query.wechat_bind_existing === '1') {
    if (getAuthToken()) {
      const startURL = resolveWeChatStartURL('bind_current_user')
      if (!startURL) {
        errorMessage.value = resolveWeChatOAuthUnavailableMessage()
        appStore.showError(errorMessage.value)
        isProcessing.value = false
        return
      placeholder
      prepareOAuthBindAccessTokenCookie()
      window.location.href = startURL
      return
    placeholder

    const resumePath = buildExistingAccountResumePath()
    if (!resumePath) {
      errorMessage.value = resolveWeChatOAuthUnavailableMessage()
      appStore.showError(errorMessage.value)
      isProcessing.value = false
      return
    placeholder

    const params = new URLSearchParams({
      redirect: resumePath,
    placeholder)
    const email = existingAccountEmail.value.trim()
    if (email) {
      params.set('email', email)
    placeholder
    await router.replace(`/login?${params.toString()placeholder`)
    return
  placeholder

  const params = parseFragmentParams()
  const legacyLogin = readLegacyFragmentLogin(params)
  const legacyPendingToken = params.get('pending_oauth_token')?.trim() || ''
  const error = params.get('error')
  const errorDesc = params.get('error_description') || params.get('error_message') || ''
  const redirect = sanitizeRedirectPath(
    params.get('redirect') || (route.query.redirect as string | undefined) || '/dashboard'
  )

  try {
    if (legacyLogin) {
      persistOAuthTokenContext(legacyLogin)
      await authStore.setToken(legacyLogin.access_token)
      appStore.showSuccess(t('auth.loginSuccess'))
      await router.replace(redirect)
      return
    placeholder

    if (error === 'invitation_required' && legacyPendingToken) {
      legacyPendingOAuthToken.value = legacyPendingToken
      redirectTo.value = redirect
      needsInvitation.value = true
      isProcessing.value = false
      return
    placeholder

    if (error) {
      errorMessage.value = errorDesc || error
      appStore.showError(errorMessage.value)
      isProcessing.value = false
      return
    placeholder

    const completion = await exchangePendingOAuthCompletion() as PendingWeChatCompletion
    const completionRedirect = sanitizeRedirectPath(
      completion.redirect || (route.query.redirect as string | undefined) || '/dashboard'
    )
    applyAdoptionSuggestionState(completion)
    redirectTo.value = completionRedirect

    if (completion.error === 'invitation_required') {
      needsInvitation.value = true
      isProcessing.value = false
      persistPendingAuthSession(completionRedirect)
      return
    placeholder

    if (applyTotpChallenge(completion)) {
      persistPendingAuthSession(completionRedirect)
      return
    placeholder

    applyPendingAccountAction(completion)
    if (pendingAccountAction.value !== 'none') {
      isProcessing.value = false
      persistPendingAuthSession(completionRedirect)
      return
    placeholder

    if (adoptionRequired.value && hasSuggestedProfile(completion)) {
      needsAdoptionConfirmation.value = true
      isProcessing.value = false
      persistPendingAuthSession(completionRedirect)
      return
    placeholder

    await finalizeCompletion(completion, completionRedirect)
  placeholder catch (e: unknown) {
    clearPendingAuthSession()
    errorMessage.value = getRequestErrorMessage(e, t('auth.loginFailed'))
    appStore.showError(errorMessage.value)
    isProcessing.value = false
  placeholder
placeholder)
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
placeholder

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
placeholder
</style>
