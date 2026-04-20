<template>
  <AuthLayout>
    <div class="space-y-6">
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.linuxdo.callbackTitle') placeholderplaceholder
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ isProcessing ? t('auth.linuxdo.callbackProcessing') : t('auth.linuxdo.callbackHint') placeholderplaceholder
        </p>
      </div>

      <transition name="fade">
        <div
          v-if="needsInvitation || needsAdoptionConfirmation || needsCreateAccount || needsBindLogin"
          class="space-y-4"
        >
          <div
            v-if="adoptionRequired && (suggestedDisplayName || suggestedAvatarUrl)"
            class="rounded-xl border border-gray-200 bg-gray-50 p-4 dark:border-dark-600 dark:bg-dark-800/60"
          >
            <div class="space-y-3">
              <div class="space-y-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  Use LinuxDo profile details
                </p>
                <p class="text-xs text-gray-500 dark:text-dark-400">
                  Choose whether to apply the nickname or avatar from LinuxDo to this account.
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
                  alt="LinuxDo avatar"
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
              {{ t('auth.linuxdo.invitationRequired') placeholderplaceholder
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
              {{ isSubmitting ? t('auth.linuxdo.completing') : t('auth.linuxdo.completeRegistration') placeholderplaceholder
            </button>
          </template>

          <template v-else-if="needsAdoptionConfirmation">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Review the LinuxDo profile details before continuing.
            </p>
            <button class="btn btn-primary w-full" :disabled="isSubmitting" @click="handleContinueLogin">
              {{ isSubmitting ? t('common.processing') : 'Continue' placeholderplaceholder
            </button>
          </template>

          <template v-else-if="needsCreateAccount">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Enter an email address to create your account and continue.
            </p>
            <div class="space-y-3">
              <input
                v-model="pendingAccountEmail"
                data-testid="linuxdo-create-account-email"
                type="email"
                class="input w-full"
                placeholder="you@example.com"
                :disabled="isSubmitting"
                @keyup.enter="handleCreateAccount"
              />
              <button
                data-testid="linuxdo-create-account-submit"
                class="btn btn-primary w-full"
                :disabled="isSubmitting || !pendingAccountEmail.trim()"
                @click="handleCreateAccount"
              >
                {{ isSubmitting ? t('common.processing') : 'Create account' placeholderplaceholder
              </button>
              <button
                class="btn btn-secondary w-full"
                :disabled="isSubmitting"
                @click="switchToBindLoginMode"
              >
                I already have an account
              </button>
            </div>
            <transition name="fade">
              <p v-if="accountActionError" class="text-sm text-red-600 dark:text-red-400">
                {{ accountActionError placeholderplaceholder
              </p>
            </transition>
          </template>

          <template v-else-if="needsBindLogin">
            <p class="text-sm text-gray-700 dark:text-gray-300">
              Log in to an existing account to bind this LinuxDo sign-in.
            </p>
            <div class="space-y-3">
              <input
                v-model="bindLoginEmail"
                data-testid="linuxdo-bind-login-email"
                type="email"
                class="input w-full"
                placeholder="you@example.com"
                :disabled="isSubmitting"
                @keyup.enter="handleBindLogin"
              />
              <input
                v-model="bindLoginPassword"
                data-testid="linuxdo-bind-login-password"
                type="password"
                class="input w-full"
                placeholder="Password"
                :disabled="isSubmitting"
                @keyup.enter="handleBindLogin"
              />
              <button
                data-testid="linuxdo-bind-login-submit"
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
                {{ t('auth.linuxdo.backToLogin') placeholderplaceholder
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
import Icon from '@/components/icons/Icon.vue'
import { apiClient placeholder from '@/api/client'
import { useAuthStore, useAppStore placeholder from '@/stores'
import {
  completeLinuxDoOAuthRegistration,
  exchangePendingOAuthCompletion,
  getOAuthCompletionKind,
  isOAuthLoginCompletion,
  persistOAuthTokenContext,
  type OAuthAdoptionDecision,
  type PendingOAuthExchangeResponse
placeholder from '@/api/auth'

const route = useRoute()
const router = useRouter()
const { t placeholder = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

const isProcessing = ref(true)
const errorMessage = ref('')

// Invitation code flow state
const needsInvitation = ref(false)
const invitationCode = ref('')
const isSubmitting = ref(false)
const invitationError = ref('')
const redirectTo = ref('/dashboard')
const adoptionRequired = ref(false)
const suggestedDisplayName = ref('')
const suggestedAvatarUrl = ref('')
const adoptDisplayName = ref(true)
const adoptAvatar = ref(true)
const needsAdoptionConfirmation = ref(false)
const pendingAccountAction = ref<'none' | 'create_account' | 'bind_login'>('none')
const pendingAccountEmail = ref('')
const bindLoginEmail = ref('')
const bindLoginPassword = ref('')
const accountActionError = ref('')
const canReturnToCreateAccount = ref(false)
const bindSuccessMessage = t('profile.authBindings.bindSuccess')

const needsCreateAccount = computed(() => pendingAccountAction.value === 'create_account')
const needsBindLogin = computed(() => pendingAccountAction.value === 'bind_login')

type LinuxDoPendingActionResponse = PendingOAuthExchangeResponse & {
  step?: string
  email?: string
  resolved_email?: string
placeholder

function parseFragmentParams(): URLSearchParams {
  const raw = typeof window !== 'undefined' ? window.location.hash : ''
  const hash = raw.startsWith('#') ? raw.slice(1) : raw
  return new URLSearchParams(hash)
placeholder

function sanitizeRedirectPath(path: string | null | undefined): string {
  if (!path) return '/dashboard'
  if (!path.startsWith('/')) return '/dashboard'
  if (path.startsWith('//')) return '/dashboard'
  if (path.includes('://')) return '/dashboard'
  if (path.includes('\n') || path.includes('\r')) return '/dashboard'
  return path
placeholder

function currentAdoptionDecision(): OAuthAdoptionDecision {
  return {
    adoptDisplayName: adoptDisplayName.value,
    adoptAvatar: adoptAvatar.value
  placeholder
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

function applyAdoptionSuggestionState(completion: {
  adoption_required?: boolean
  suggested_display_name?: string
  suggested_avatar_url?: string
placeholder) {
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

function hasSuggestedProfile(completion: {
  suggested_display_name?: string
  suggested_avatar_url?: string
placeholder): boolean {
  return Boolean(completion.suggested_display_name || completion.suggested_avatar_url)
placeholder

function extractPendingAccountEmail(completion: LinuxDoPendingActionResponse): string {
  return (completion.email || completion.resolved_email || '').trim()
placeholder

function resolvePendingAccountAction(completion: LinuxDoPendingActionResponse): 'none' | 'create_account' | 'bind_login' {
  const raw = (completion.step || completion.error || '').trim().toLowerCase()
  if (raw === 'email_required' || raw === 'create_account_required' || raw === 'create_account') {
    return 'create_account'
  placeholder
  if (raw === 'bind_login_required' || raw === 'bind_login') {
    return 'bind_login'
  placeholder
  return 'none'
placeholder

function applyPendingAccountAction(completion: LinuxDoPendingActionResponse) {
  const action = resolvePendingAccountAction(completion)
  pendingAccountAction.value = action
  accountActionError.value = ''

  const email = extractPendingAccountEmail(completion)
  if (action === 'create_account') {
    pendingAccountEmail.value = email
    canReturnToCreateAccount.value = true
    return
  placeholder

  if (action === 'bind_login') {
    bindLoginEmail.value = email
    bindLoginPassword.value = ''
    canReturnToCreateAccount.value = false
    return
  placeholder

  canReturnToCreateAccount.value = false
placeholder

function switchToBindLoginMode() {
  pendingAccountAction.value = 'bind_login'
  bindLoginEmail.value = bindLoginEmail.value.trim() || pendingAccountEmail.value.trim()
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

async function finalizeCompletion(completion: PendingOAuthExchangeResponse, redirect: string) {
  if (getOAuthCompletionKind(completion) === 'bind') {
    const bindRedirect = sanitizeRedirectPath(completion.redirect || '/profile')
    appStore.showSuccess(bindSuccessMessage)
    await router.replace(bindRedirect)
    return
  placeholder

  if (!isOAuthLoginCompletion(completion)) {
    throw new Error(t('auth.linuxdo.callbackMissingToken'))
  placeholder

  persistOAuthTokenContext(completion)
  await authStore.setToken(completion.access_token)
  appStore.showSuccess(t('auth.loginSuccess'))
  await router.replace(redirect)
placeholder

async function finalizePendingAccountResponse(completion: LinuxDoPendingActionResponse) {
  applyAdoptionSuggestionState(completion)

  if (completion.error === 'invitation_required') {
    pendingAccountAction.value = 'none'
    needsInvitation.value = true
    needsAdoptionConfirmation.value = false
    isProcessing.value = false
    return
  placeholder

  applyPendingAccountAction(completion)
  if (pendingAccountAction.value !== 'none') {
    needsInvitation.value = false
    needsAdoptionConfirmation.value = false
    isProcessing.value = false
    return
  placeholder

  const redirect = sanitizeRedirectPath(completion.redirect || redirectTo.value)
  await finalizeCompletion(completion, redirect)
placeholder

async function handleSubmitInvitation() {
  invitationError.value = ''
  if (!invitationCode.value.trim()) return

  isSubmitting.value = true
  try {
    const tokenData = await completeLinuxDoOAuthRegistration(
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
      err.response?.data?.message || err.message || t('auth.linuxdo.completeRegistrationFailed')
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

async function handleContinueLogin() {
  isSubmitting.value = true
  try {
    const completion = await exchangePendingOAuthCompletion(currentAdoptionDecision())
    await finalizeCompletion(completion, redirectTo.value)
  placeholder catch (e: unknown) {
    errorMessage.value = getRequestErrorMessage(e, t('auth.loginFailed'))
    appStore.showError(errorMessage.value)
    needsAdoptionConfirmation.value = false
  placeholder finally {
    isSubmitting.value = false
  placeholder
placeholder

async function handleCreateAccount() {
  accountActionError.value = ''
  const email = pendingAccountEmail.value.trim()
  if (!email) return

  isSubmitting.value = true
  try {
    const { data placeholder = await apiClient.post<LinuxDoPendingActionResponse>('/auth/oauth/pending/create-account', {
      email,
      ...serializeAdoptionDecision(currentAdoptionDecision())
    placeholder)
    await finalizePendingAccountResponse(data)
  placeholder catch (e: unknown) {
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
    const { data placeholder = await apiClient.post<LinuxDoPendingActionResponse>('/auth/oauth/pending/bind-login', {
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

onMounted(async () => {
  const params = parseFragmentParams()
  const error = params.get('error')
  const errorDesc = params.get('error_description') || params.get('error_message') || ''

  if (error) {
    errorMessage.value = errorDesc || error
    appStore.showError(errorMessage.value)
    isProcessing.value = false
    return
  placeholder

  try {
    const completion = await exchangePendingOAuthCompletion()
    const redirect = sanitizeRedirectPath(
      completion.redirect || (route.query.redirect as string | undefined) || '/dashboard'
    )
    applyAdoptionSuggestionState(completion)
    redirectTo.value = redirect

    if (completion.error === 'invitation_required') {
      needsInvitation.value = true
      isProcessing.value = false
      return
    placeholder

    applyPendingAccountAction(completion as LinuxDoPendingActionResponse)
    if (pendingAccountAction.value !== 'none') {
      isProcessing.value = false
      return
    placeholder

    if (adoptionRequired.value && hasSuggestedProfile(completion)) {
      needsAdoptionConfirmation.value = true
      isProcessing.value = false
      return
    placeholder

    await finalizeCompletion(completion, redirect)
  placeholder catch (e: unknown) {
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
