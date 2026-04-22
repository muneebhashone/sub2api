<template>
  <div :class="props.embedded ? 'space-y-4' : 'card overflow-hidden'">
    <div
      v-if="!props.embedded"
      class="border-b border-gray-100 px-6 py-4 dark:border-dark-700"
    >
      <h2 class="text-lg font-medium text-gray-900 dark:text-white">
        {{ t('profile.authBindings.title') placeholderplaceholder
      </h2>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        {{ t('profile.authBindings.description') placeholderplaceholder
      </p>
    </div>

    <div :class="props.embedded ? 'space-y-4' : 'divide-y divide-gray-100 dark:divide-dark-700'">
      <div v-if="props.embedded">
        <p class="text-sm font-semibold text-gray-900 dark:text-white">
          {{ t('profile.authBindings.title') placeholderplaceholder
        </p>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          {{ t('profile.authBindings.description') placeholderplaceholder
        </p>
      </div>

      <div
        v-for="item in providerItems"
        :key="item.provider"
        :class="rowClass"
      >
        <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div class="flex min-w-0 flex-1 items-start gap-4">
            <div
              :class="providerIconClass(item.provider)"
              class="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl text-sm font-semibold"
            >
              <Icon
                v-if="item.provider === 'email'"
                name="mail"
                size="sm"
                class="text-current"
              />
              <span v-else>{{ providerInitial(item.provider) placeholderplaceholder</span>
            </div>

            <div class="min-w-0 flex-1 space-y-3">
              <div class="flex flex-wrap items-center gap-2">
                <h3 class="font-medium text-gray-900 dark:text-white">
                  {{ item.label placeholderplaceholder
                </h3>
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

              <p
                v-if="providerSummary(item.provider)"
                class="text-sm text-gray-600 dark:text-gray-300"
              >
                {{ providerSummary(item.provider) placeholderplaceholder
              </p>

              <div
                v-if="item.details && (item.details.display_name || item.details.subject_hint || bindingCountLabel(item.details) || item.details.note)"
                class="grid gap-1 text-sm text-gray-500 dark:text-gray-400"
              >
                <p
                  v-if="item.details.display_name"
                  class="font-medium text-gray-700 dark:text-gray-200"
                >
                  {{ item.details.display_name placeholderplaceholder
                </p>
                <p v-if="item.details.subject_hint">
                  {{ item.details.subject_hint placeholderplaceholder
                </p>
                <p v-if="bindingCountLabel(item.details)">
                  {{ bindingCountLabel(item.details) placeholderplaceholder
                </p>
                <p v-if="item.details.note">
                  {{ item.details.note placeholderplaceholder
                </p>
              </div>

              <div
                v-if="item.provider === 'email' && showEmailForm"
                data-testid="profile-binding-email-form"
                class="grid gap-2 sm:grid-cols-[minmax(0,1.4fr)_auto]"
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
                  :placeholder="emailPasswordPlaceholder"
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
                      : emailSubmitActionLabel
                  placeholderplaceholder
                </button>
              </div>
            </div>
          </div>

          <div class="flex shrink-0 flex-wrap items-center gap-3">
            <button
              v-if="item.provider === 'email' && compact"
              data-testid="profile-binding-email-toggle"
              type="button"
              class="btn btn-secondary btn-sm"
              @click="toggleEmailForm"
            >
              {{
                showEmailForm
                  ? t('profile.authBindings.hideEmailFormAction')
                  : t('profile.authBindings.manageEmailAction')
              placeholderplaceholder
            </button>
            <button
              v-if="item.canBind"
              :data-testid="`profile-binding-${item.providerplaceholder-action`"
              type="button"
              class="btn btn-primary btn-sm"
              @click="startBinding(item.provider)"
            >
              {{ t('profile.authBindings.bindAction', { providerName: item.label placeholder) placeholderplaceholder
            </button>
            <button
              v-if="item.canUnbind"
              :data-testid="`profile-binding-${item.providerplaceholder-unbind`"
              type="button"
              class="btn btn-secondary btn-sm"
              :disabled="unbindingProvider === item.provider"
              @click="handleUnbindForItem(item.provider, item.label)"
            >
              {{
                unbindingProvider === item.provider
                  ? t('common.loading')
                  : t('profile.authBindings.unbindAction')
              placeholderplaceholder
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
import {
  bindEmailIdentity,
  sendEmailBindingCode,
  startOAuthBinding,
  unbindAuthIdentity,
placeholder from '@/api/user'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore, useAuthStore placeholder from '@/stores'
import type { User, UserAuthBindingStatus, UserAuthProvider placeholder from '@/types'

type BindableProvider = Exclude<UserAuthProvider, 'email'>

const props = withDefaults(
  defineProps<{
    user: User | null
    linuxdoEnabled?: boolean
    oidcEnabled?: boolean
    oidcProviderName?: string
    wechatEnabled?: boolean
    wechatOpenEnabled?: boolean
    wechatMpEnabled?: boolean
    embedded?: boolean
    compact?: boolean
  placeholder>(),
  {
    linuxdoEnabled: false,
    oidcEnabled: false,
    oidcProviderName: 'OIDC',
    wechatEnabled: false,
    wechatOpenEnabled: undefined,
    wechatMpEnabled: undefined,
    embedded: false,
    compact: false,
  placeholder
)

const { t placeholder = useI18n()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

const localUser = ref<User | null>(null)
const isSendingEmailCode = ref(false)
const isBindingEmail = ref(false)
const isEmailFormExpanded = ref(!props.compact)
const unbindingProvider = ref<BindableProvider | null>(null)
const emailBindingForm = reactive({
  email: '',
  verifyCode: '',
  password: '',
placeholder)

watch(
  () => props.user,
  (user) => {
    localUser.value = null
    if (!user) {
      return
    placeholder
    if (typeof user.email === 'string' && !user.email.endsWith('.invalid')) {
      emailBindingForm.email = user.email
    placeholder
  placeholder,
  { immediate: true placeholder
)

watch(
  () => props.compact,
  (value) => {
    if (!value) {
      isEmailFormExpanded.value = true
    placeholder
  placeholder,
  { immediate: true placeholder
)

const currentUser = computed(() => localUser.value ?? props.user)
const compact = computed(() => props.compact)
const rowClass = computed(() =>
  props.embedded
    ? compact.value
      ? 'rounded-2xl border border-gray-100 bg-white p-4 shadow-sm dark:border-dark-700 dark:bg-dark-900/40'
      : 'rounded-2xl border border-gray-100 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-900/30'
    : 'px-6 py-5'
)
const emailBound = computed(() => getBindingStatus('email'))
const showEmailForm = computed(() => !compact.value || isEmailFormExpanded.value)
const emailPasswordPlaceholder = computed(() =>
  emailBound.value
    ? t('profile.authBindings.replaceEmailPasswordPlaceholder')
    : t('profile.authBindings.passwordPlaceholder')
)
const emailSubmitActionLabel = computed(() =>
  emailBound.value
    ? t('profile.authBindings.confirmEmailReplaceAction')
    : t('profile.authBindings.confirmEmailBindAction')
)

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
    if (typeof user?.email_bound === 'boolean') {
      return user.email_bound
    placeholder
    const nested = user?.auth_bindings?.email ?? user?.identity_bindings?.email
    const normalized = normalizeBindingStatus(nested)
    return normalized ?? false
  placeholder

  const directFlag = user?.[`${providerplaceholder_bound` as keyof User]
  if (typeof directFlag === 'boolean') {
    return directFlag
  placeholder

  const nested = user?.auth_bindings?.[provider] ?? user?.identity_bindings?.[provider]
  const normalized = normalizeBindingStatus(nested)
  return normalized ?? false
placeholder

function getBindingDetails(provider: UserAuthProvider): UserAuthBindingStatus | null {
  const binding = currentUser.value?.auth_bindings?.[provider] ?? currentUser.value?.identity_bindings?.[provider]
  if (!binding || typeof binding === 'boolean') {
    return null
  placeholder
  return binding
placeholder

function isProviderEnabledForBinding(provider: BindableProvider): boolean {
  if (provider === 'linuxdo') {
    return props.linuxdoEnabled
  placeholder
  if (provider === 'oidc') {
    return props.oidcEnabled
  placeholder
  return resolvedWeChatBinding.value.mode !== null
placeholder

const providerItems = computed(() => [
  {
    provider: 'email' as const,
    label: t('profile.authBindings.providers.email'),
    bound: getBindingStatus('email'),
    canBind: false,
    canUnbind: false,
    details: getBindingDetails('email'),
  placeholder,
  {
    provider: 'linuxdo' as const,
    label: t('profile.authBindings.providers.linuxdo'),
    bound: getBindingStatus('linuxdo'),
    canBind:
      !getBindingStatus('linuxdo') &&
      isProviderEnabledForBinding('linuxdo') &&
      (getBindingDetails('linuxdo')?.can_bind ?? true),
    canUnbind: Boolean(getBindingStatus('linuxdo') && getBindingDetails('linuxdo')?.can_unbind),
    details: getBindingDetails('linuxdo'),
  placeholder,
  {
    provider: 'oidc' as const,
    label: t('profile.authBindings.providers.oidc', { providerName: props.oidcProviderName placeholder),
    bound: getBindingStatus('oidc'),
    canBind:
      !getBindingStatus('oidc') &&
      isProviderEnabledForBinding('oidc') &&
      (getBindingDetails('oidc')?.can_bind ?? true),
    canUnbind: Boolean(getBindingStatus('oidc') && getBindingDetails('oidc')?.can_unbind),
    details: getBindingDetails('oidc'),
  placeholder,
  {
    provider: 'wechat' as const,
    label: t('profile.authBindings.providers.wechat'),
    bound: getBindingStatus('wechat'),
    canBind:
      !getBindingStatus('wechat') &&
      isProviderEnabledForBinding('wechat') &&
      (getBindingDetails('wechat')?.can_bind ?? true),
    canUnbind: Boolean(getBindingStatus('wechat') && getBindingDetails('wechat')?.can_unbind),
    details: getBindingDetails('wechat'),
  placeholder,
])

function providerInitial(provider: UserAuthProvider): string {
  if (provider === 'linuxdo') {
    return 'L'
  placeholder
  if (provider === 'wechat') {
    return 'W'
  placeholder
  if (provider === 'oidc') {
    return 'O'
  placeholder
  return 'E'
placeholder

function providerIconClass(provider: UserAuthProvider): string {
  if (provider === 'linuxdo') {
    return 'bg-orange-100 text-orange-600 dark:bg-orange-900/20 dark:text-orange-300'
  placeholder
  if (provider === 'wechat') {
    return 'bg-green-100 text-green-600 dark:bg-green-900/20 dark:text-green-300'
  placeholder
  if (provider === 'oidc') {
    return 'bg-sky-100 text-sky-600 dark:bg-sky-900/20 dark:text-sky-300'
  placeholder
  return 'bg-primary-100 text-primary-600 dark:bg-primary-900/20 dark:text-primary-300'
placeholder

function providerSummary(provider: UserAuthProvider): string {
  if (provider === 'email') {
    return currentUser.value?.email || ''
  placeholder
  return ''
placeholder

function bindingCountLabel(details: UserAuthBindingStatus | null): string {
  if (!details || typeof details.bound_count !== 'number' || details.bound_count <= 1) {
    return ''
  placeholder
  return t('profile.authBindings.boundCount', { count: details.bound_count placeholder)
placeholder

function toggleEmailForm(): void {
  isEmailFormExpanded.value = !isEmailFormExpanded.value
placeholder

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

async function handleUnbind(provider: BindableProvider, providerLabel: string): Promise<void> {
  unbindingProvider.value = provider
  try {
    const user = await unbindAuthIdentity(provider)
    applyUpdatedUser(user)
    appStore.showSuccess(t('profile.authBindings.unbindSuccess', { providerName: providerLabel placeholder))
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder).message || t('common.tryAgain'))
  placeholder finally {
    unbindingProvider.value = null
  placeholder
placeholder

function handleUnbindForItem(provider: UserAuthProvider, providerLabel: string): void {
  if (provider === 'email') {
    return
  placeholder
  void handleUnbind(provider, providerLabel)
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
  if (requireCode && !emailBound.value && emailBindingForm.password.length < 6) {
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
    const replacingBoundEmail = emailBound.value
    applyUpdatedUser(user)
    emailBindingForm.verifyCode = ''
    emailBindingForm.password = ''
    if (compact.value) {
      isEmailFormExpanded.value = false
    placeholder
    appStore.showSuccess(
      replacingBoundEmail
        ? t('profile.authBindings.replaceSuccess')
        : t('profile.authBindings.bindSuccess')
    )
  placeholder catch (error) {
    appStore.showError((error as { message?: string placeholder).message || t('common.tryAgain'))
  placeholder finally {
    isBindingEmail.value = false
  placeholder
placeholder
</script>
