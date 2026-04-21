<template>
  <div class="rounded-2xl border border-gray-100 bg-gray-50/80 p-4 dark:border-dark-700 dark:bg-dark-900/30">
    <div>
      <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
        {{ t('profile.authBindings.title') placeholderplaceholder
      </h3>
      <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
        {{ t('profile.authBindings.description') placeholderplaceholder
      </p>
    </div>

    <div class="mt-4 space-y-2">
      <div
        v-for="item in providerItems"
        :key="item.provider"
        class="rounded-xl bg-white/80 px-3 py-3 dark:bg-dark-800/70"
      >
        <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between sm:gap-4">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <div class="text-sm font-medium text-gray-900 dark:text-white">
                {{ item.label placeholderplaceholder
              </div>
              <span
                :data-testid="`profile-binding-${item.providerplaceholder-status`"
                :class="['badge', item.bound ? 'badge-success' : 'badge-gray']"
              >
                {{
                  item.bound
                    ? t('profile.authBindings.status.bound')
                    : t('profile.authBindings.status.notBound')
                placeholderplaceholder
              </span>
            </div>

            <div
              v-if="item.provider === 'email' && !item.bound"
              class="mt-3 grid gap-2 sm:grid-cols-[minmax(0,1.4fr)_auto]"
            >
              <input
                v-model.trim="emailBindingForm.email"
                data-testid="profile-binding-email-input"
                type="email"
                class="input"
                :placeholder="t('profile.authBindings.emailPlaceholder')"
                :disabled="isSendingEmailCode || isBindingEmail"
              />
              <button
                data-testid="profile-binding-email-send-code"
                type="button"
                class="btn btn-secondary btn-sm"
                :disabled="isSendingEmailCode || isBindingEmail"
                @click="sendEmailCode"
              >
                {{
                  isSendingEmailCode
                    ? t('common.loading')
                    : t('profile.authBindings.sendCodeAction')
                placeholderplaceholder
              </button>
              <input
                v-model.trim="emailBindingForm.verifyCode"
                data-testid="profile-binding-email-code-input"
                type="text"
                inputmode="numeric"
                maxlength="6"
                class="input"
                :placeholder="t('profile.authBindings.codePlaceholder')"
                :disabled="isBindingEmail"
              />
              <input
                v-model="emailBindingForm.password"
                data-testid="profile-binding-email-password-input"
                type="password"
                class="input"
                :placeholder="t('profile.authBindings.passwordPlaceholder')"
                :disabled="isBindingEmail"
              />
              <button
                data-testid="profile-binding-email-submit"
                type="button"
                class="btn btn-primary btn-sm sm:col-span-2"
                :disabled="isBindingEmail"
                @click="bindEmail"
              >
                {{
                  isBindingEmail
                    ? t('common.loading')
                    : t('profile.authBindings.confirmEmailBindAction')
                placeholderplaceholder
              </button>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-2">
            <button
              v-if="item.canBind"
              :data-testid="`profile-binding-${item.providerplaceholder-action`"
              type="button"
              class="btn btn-secondary btn-sm"
              @click="startBinding(item.provider)"
            >
              {{ t('profile.authBindings.bindAction', { providerName: item.label placeholder) placeholderplaceholder
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRoute placeholder from 'vue-router'
import {
  hasExplicitWeChatOAuthCapabilities,
  resolveWeChatOAuthStartStrict,
  type WeChatOAuthPublicSettings,
placeholder from '@/api/auth'
import { bindEmailIdentity, sendEmailBindingCode, startOAuthBinding placeholder from '@/api/user'
import { useAppStore, useAuthStore placeholder from '@/stores'
import type { User, UserAuthBindingStatus, UserAuthProvider placeholder from '@/types'

const props = withDefaults(
  defineProps<{
    user: User | null
    linuxdoEnabled?: boolean
    oidcEnabled?: boolean
    oidcProviderName?: string
    wechatEnabled?: boolean
    wechatOpenEnabled?: boolean
    wechatMpEnabled?: boolean
  placeholder>(),
  {
    linuxdoEnabled: false,
    oidcEnabled: false,
    oidcProviderName: 'OIDC',
    wechatEnabled: false,
    wechatOpenEnabled: undefined,
    wechatMpEnabled: undefined,
  placeholder
)

const { t placeholder = useI18n()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

const localUser = ref<User | null>(null)
const isSendingEmailCode = ref(false)
const isBindingEmail = ref(false)
const emailBindingForm = reactive({
  email: '',
  verifyCode: '',
  password: '',
placeholder)

watch(
  () => props.user,
  (user) => {
    localUser.value = null
    if (!user || getBindingStatusForUser(user, 'email')) {
      return
    placeholder
    if (typeof user.email === 'string' && !user.email.endsWith('.invalid')) {
      emailBindingForm.email = user.email
    placeholder
  placeholder,
  { immediate: true placeholder
)

const currentUser = computed(() => localUser.value ?? props.user)

const wechatOAuthSettings = computed<WeChatOAuthPublicSettings | null>(() => {
  if (hasExplicitWeChatOAuthCapabilities(appStore.cachedPublicSettings)) {
    return appStore.cachedPublicSettings
  placeholder

  if (typeof props.wechatOpenEnabled === 'boolean' && typeof props.wechatMpEnabled === 'boolean') {
    return {
      wechat_oauth_enabled: props.wechatEnabled,
      wechat_oauth_open_enabled: props.wechatOpenEnabled,
      wechat_oauth_mp_enabled: props.wechatMpEnabled,
    placeholder
  placeholder

  return null
placeholder)

const resolvedWeChatBinding = computed(() => resolveWeChatOAuthStartStrict(wechatOAuthSettings.value))

function normalizeBindingStatus(binding: boolean | UserAuthBindingStatus | undefined): boolean | null {
  if (typeof binding === 'boolean') {
    return binding
  placeholder
  if (!binding) {
    return null
  placeholder
  if (typeof binding.bound === 'boolean') {
    return binding.bound
  placeholder
  return Boolean(binding.provider_subject || binding.issuer || binding.provider_key)
placeholder

function getBindingStatus(provider: UserAuthProvider): boolean {
  return getBindingStatusForUser(currentUser.value, provider)
placeholder

function getBindingStatusForUser(user: User | null | undefined, provider: UserAuthProvider): boolean {
  if (provider === 'email') {
    return typeof user?.email_bound === 'boolean' ? user.email_bound : Boolean(user?.email)
  placeholder

  const directFlag = user?.[`${providerplaceholder_bound` as keyof User]
  if (typeof directFlag === 'boolean') {
    return directFlag
  placeholder

  const nested = user?.auth_bindings?.[provider] ?? user?.identity_bindings?.[provider]
  const normalized = normalizeBindingStatus(nested)
  return normalized ?? false
placeholder

const providerItems = computed(() => [
  {
    provider: 'email' as const,
    label: t('profile.authBindings.providers.email'),
    bound: getBindingStatus('email'),
    canBind: false,
  placeholder,
  {
    provider: 'linuxdo' as const,
    label: t('profile.authBindings.providers.linuxdo'),
    bound: getBindingStatus('linuxdo'),
    canBind: props.linuxdoEnabled && !getBindingStatus('linuxdo'),
  placeholder,
  {
    provider: 'oidc' as const,
    label: t('profile.authBindings.providers.oidc', { providerName: props.oidcProviderName placeholder),
    bound: getBindingStatus('oidc'),
    canBind: props.oidcEnabled && !getBindingStatus('oidc'),
  placeholder,
  {
    provider: 'wechat' as const,
    label: t('profile.authBindings.providers.wechat'),
    bound: getBindingStatus('wechat'),
    canBind: resolvedWeChatBinding.value.mode !== null && !getBindingStatus('wechat'),
  placeholder,
])

function startBinding(provider: UserAuthProvider): void {
  if (provider === 'email') {
    return
  placeholder
  startOAuthBinding(provider, {
    redirectTo: route.fullPath || '/profile',
    wechatOAuthSettings: provider === 'wechat' ? wechatOAuthSettings.value : null,
  placeholder)
placeholder

function applyUpdatedUser(user: User): void {
  localUser.value = user
  authStore.user = user
placeholder

function validateEmailBindingForm(requireCode: boolean): boolean {
  if (!emailBindingForm.email) {
    appStore.showError(t('auth.emailRequired'))
    return false
  placeholder
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(emailBindingForm.email)) {
    appStore.showError(t('auth.invalidEmail'))
    return false
  placeholder
  if (requireCode && !emailBindingForm.verifyCode) {
    appStore.showError(t('auth.codeRequired'))
    return false
  placeholder
  if (requireCode && !emailBindingForm.password) {
    appStore.showError(t('auth.passwordRequired'))
    return false
  placeholder
  if (requireCode && emailBindingForm.password.length < 6) {
    appStore.showError(t('auth.passwordMinLength'))
    return false
  placeholder
  return true
placeholder

async function sendEmailCode(): Promise<void> {
  if (!validateEmailBindingForm(false)) {
    return
  placeholder

  isSendingEmailCode.value = true
  try {
    await sendEmailBindingCode(emailBindingForm.email)
    appStore.showSuccess(t('profile.authBindings.codeSentTo', { email: emailBindingForm.email placeholder))
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder).message || t('auth.sendCodeFailed'))
  placeholder finally {
    isSendingEmailCode.value = false
  placeholder
placeholder

async function bindEmail(): Promise<void> {
  if (!validateEmailBindingForm(true)) {
    return
  placeholder

  isBindingEmail.value = true
  try {
    const user = await bindEmailIdentity({
      email: emailBindingForm.email,
      verify_code: emailBindingForm.verifyCode,
      password: emailBindingForm.password,
    placeholder)
    applyUpdatedUser(user)
    emailBindingForm.verifyCode = ''
    emailBindingForm.password = ''
    appStore.showSuccess(t('profile.authBindings.bindSuccess'))
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder).message || t('common.tryAgain'))
  placeholder finally {
    isBindingEmail.value = false
  placeholder
placeholder
</script>
