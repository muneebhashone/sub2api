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
        class="flex items-center justify-between gap-3 rounded-xl bg-white/80 px-3 py-2.5 dark:bg-dark-800/70"
      >
        <div class="min-w-0">
          <div class="text-sm font-medium text-gray-900 dark:text-white">
            {{ item.label placeholderplaceholder
          </div>
        </div>

        <div class="flex shrink-0 items-center gap-2">
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
</template>

<script setup lang="ts">
import { computed placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useRoute placeholder from 'vue-router'
import {
  hasExplicitWeChatOAuthCapabilities,
  resolveWeChatOAuthStartStrict,
  type WeChatOAuthPublicSettings,
placeholder from '@/api/auth'
import { startOAuthBinding placeholder from '@/api/user'
import { useAppStore placeholder from '@/stores'
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
  const currentUser = props.user

  if (provider === 'email') {
    return typeof currentUser?.email_bound === 'boolean'
      ? currentUser.email_bound
      : Boolean(currentUser?.email)
  placeholder

  const directFlag = currentUser?.[`${providerplaceholder_bound` as keyof User]
  if (typeof directFlag === 'boolean') {
    return directFlag
  placeholder

  const nested = currentUser?.auth_bindings?.[provider] ?? currentUser?.identity_bindings?.[provider]
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
</script>
