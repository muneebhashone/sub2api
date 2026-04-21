<template>
  <div class="card overflow-hidden">
    <div
      class="border-b border-gray-100 bg-gradient-to-r from-primary-500/10 to-primary-600/5 px-6 py-5 dark:border-dark-700 dark:from-primary-500/20 dark:to-primary-600/10"
    >
      <div class="flex items-center gap-4">
        <div
          class="flex h-16 w-16 items-center justify-center overflow-hidden rounded-2xl bg-gradient-to-br from-primary-500 to-primary-600 text-2xl font-bold text-white shadow-lg shadow-primary-500/20"
        >
          <img
            v-if="avatarUrl"
            :src="avatarUrl"
            :alt="displayName"
            class="h-full w-full object-cover"
          >
          <span v-else>{{ avatarInitial placeholderplaceholder</span>
        </div>
        <div class="min-w-0 flex-1">
          <h2 class="truncate text-lg font-semibold text-gray-900 dark:text-white">
            {{ user?.email placeholderplaceholder
          </h2>
          <div class="mt-1 flex items-center gap-2">
            <span :class="['badge', user?.role === 'admin' ? 'badge-primary' : 'badge-gray']">
              {{ user?.role === 'admin' ? t('profile.administrator') : t('profile.user') placeholderplaceholder
            </span>
            <span
              :class="['badge', user?.status === 'active' ? 'badge-success' : 'badge-danger']"
            >
              {{ user?.status placeholderplaceholder
            </span>
          </div>
        </div>
      </div>
    </div>
    <div class="px-6 py-4">
      <div class="space-y-3">
        <div class="flex items-center gap-3 text-sm text-gray-600 dark:text-gray-400">
          <Icon name="mail" size="sm" class="text-gray-400 dark:text-gray-500" />
          <span class="truncate">{{ user?.email placeholderplaceholder</span>
        </div>
        <div
          v-if="user?.username"
          class="flex items-center gap-3 text-sm text-gray-600 dark:text-gray-400"
        >
          <Icon name="user" size="sm" class="text-gray-400 dark:text-gray-500" />
          <span class="truncate">{{ user.username placeholderplaceholder</span>
        </div>
      </div>

      <div
        v-if="sourceHints.length"
        class="mt-4 grid gap-2 rounded-2xl border border-gray-100 bg-gray-50/80 p-3 text-xs text-gray-500 dark:border-dark-700 dark:bg-dark-900/30 dark:text-gray-400"
      >
        <div
          v-for="hint in sourceHints"
          :key="hint.key"
          class="flex items-start gap-2"
        >
          <Icon name="link" size="sm" class="mt-0.5 text-gray-400 dark:text-gray-500" />
          <span>{{ hint.text placeholderplaceholder</span>
        </div>
      </div>

      <div
        class="mt-4 rounded-2xl border border-gray-100 bg-white/90 p-4 dark:border-dark-700 dark:bg-dark-900/30"
      >
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('profile.avatar.title') placeholderplaceholder
            </h3>
            <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
              {{ t('profile.avatar.description') placeholderplaceholder
            </p>
          </div>
          <button
            data-testid="profile-avatar-delete"
            type="button"
            class="btn btn-secondary btn-sm"
            :disabled="avatarSaving"
            @click="handleAvatarDelete"
          >
            {{ t('common.delete') placeholderplaceholder
          </button>
        </div>

        <div class="mt-3 space-y-3">
          <label
            for="profile-avatar-input"
            class="text-xs font-medium uppercase tracking-wide text-gray-500 dark:text-gray-400"
          >
            {{ t('profile.avatar.inputLabel') placeholderplaceholder
          </label>
          <textarea
            id="profile-avatar-input"
            data-testid="profile-avatar-input"
            v-model="avatarDraft"
            rows="3"
            class="input min-h-[88px]"
            :placeholder="t('profile.avatar.inputPlaceholder')"
          />
          <div class="flex flex-wrap items-center gap-2">
            <label class="btn btn-secondary btn-sm cursor-pointer">
              <input
                data-testid="profile-avatar-file-input"
                type="file"
                accept="image/*"
                class="hidden"
                @change="handleAvatarFileChange"
              >
              {{ t('profile.avatar.uploadAction') placeholderplaceholder
            </label>
            <button
              data-testid="profile-avatar-save"
              type="button"
              class="btn btn-primary btn-sm"
              :disabled="avatarSaving"
              @click="handleAvatarSave"
            >
              {{ t('common.save') placeholderplaceholder
            </button>
            <span class="text-xs text-gray-400 dark:text-gray-500">
              {{ t('profile.avatar.uploadHint') placeholderplaceholder
            </span>
          </div>
        </div>
      </div>

      <ProfileIdentityBindingsSection
        class="mt-4"
        :user="user"
        :linuxdo-enabled="linuxdoEnabled"
        :oidc-enabled="oidcEnabled"
        :oidc-provider-name="oidcProviderName"
        :wechat-enabled="wechatEnabled"
        :wechat-open-enabled="wechatOpenEnabled"
        :wechat-mp-enabled="wechatMpEnabled"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { userAPI placeholder from '@/api'
import Icon from '@/components/icons/Icon.vue'
import ProfileIdentityBindingsSection from '@/components/user/profile/ProfileIdentityBindingsSection.vue'
import { useAppStore placeholder from '@/stores/app'
import { useAuthStore placeholder from '@/stores/auth'
import type { User, UserAuthProvider, UserProfileSourceContext placeholder from '@/types'
import { extractApiErrorMessage placeholder from '@/utils/apiError'

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
const authStore = useAuthStore()
const appStore = useAppStore()

const maxAvatarBytes = 100 * 1024
const targetAvatarUploadBytes = 20 * 1024
const avatarScaleSteps = [1, 0.92, 0.84, 0.76, 0.68, 0.6, 0.52, 0.44, 0.36]
const avatarQualitySteps = [0.92, 0.84, 0.76, 0.68, 0.6, 0.52, 0.44, 0.36]
const avatarDraft = ref(props.user?.avatar_url?.trim() || '')
const avatarSaving = ref(false)

const providerLabels = computed<Record<UserAuthProvider, string>>(() => ({
  email: t('profile.authBindings.providers.email'),
  linuxdo: t('profile.authBindings.providers.linuxdo'),
  oidc: t('profile.authBindings.providers.oidc', { providerName: props.oidcProviderName placeholder),
  wechat: t('profile.authBindings.providers.wechat'),
placeholder))

const avatarUrl = computed(() => props.user?.avatar_url?.trim() || '')
const displayName = computed(() => props.user?.username?.trim() || props.user?.email?.trim() || 'User')
const avatarInitial = computed(() => displayName.value.charAt(0).toUpperCase() || 'U')

watch(
  () => props.user?.avatar_url,
  (value) => {
    avatarDraft.value = value?.trim() || ''
  placeholder
)

function normalizeProvider(value: string): UserAuthProvider | null {
  const normalized = value.trim().toLowerCase()
  if (normalized === 'email' || normalized === 'linuxdo' || normalized === 'wechat') {
    return normalized
  placeholder
  if (normalized === 'oidc' || normalized.startsWith('oidc:') || normalized.startsWith('oidc/')) {
    return 'oidc'
  placeholder
  return null
placeholder

function readObjectString(source: Record<string, unknown>, ...keys: string[]): string {
  for (const key of keys) {
    const value = source[key]
    if (typeof value === 'string' && value.trim()) {
      return value.trim()
    placeholder
  placeholder
  return ''
placeholder

function resolveThirdPartySource(
  rawSource: string | UserProfileSourceContext | null | undefined
): { provider: UserAuthProvider; label: string placeholder | null {
  if (!rawSource) {
    return null
  placeholder

  if (typeof rawSource === 'string') {
    const provider = normalizeProvider(rawSource)
    if (!provider || provider === 'email') {
      return null
    placeholder
    return {
      provider,
      label: providerLabels.value[provider],
    placeholder
  placeholder

  const sourceRecord = rawSource as Record<string, unknown>
  const provider = normalizeProvider(
    readObjectString(sourceRecord, 'provider', 'source', 'provider_type', 'auth_provider')
  )
  if (!provider || provider === 'email') {
    return null
  placeholder

  const explicitLabel = readObjectString(
    sourceRecord,
    'provider_label',
    'label',
    'provider_name',
    'providerName'
  )

  return {
    provider,
    label: explicitLabel || providerLabels.value[provider],
  placeholder
placeholder

const sourceHints = computed(() => {
  const currentUser = props.user
  if (!currentUser) {
    return []
  placeholder

  const hints: Array<{ key: string; text: string placeholder> = []
  const avatarSource = resolveThirdPartySource(
    currentUser.profile_sources?.avatar ?? currentUser.avatar_source
  )
  const usernameSource = resolveThirdPartySource(
    currentUser.profile_sources?.username ??
      currentUser.profile_sources?.display_name ??
      currentUser.profile_sources?.nickname ??
      currentUser.display_name_source ??
      currentUser.username_source ??
      currentUser.nickname_source
  )

  if (avatarSource) {
    hints.push({
      key: 'avatar',
      text: t('profile.authBindings.source.avatar', { providerName: avatarSource.label placeholder),
    placeholder)
  placeholder

  if (usernameSource) {
    hints.push({
      key: 'username',
      text: t('profile.authBindings.source.username', { providerName: usernameSource.label placeholder),
    placeholder)
  placeholder

  return hints
placeholder)

function estimateDataURLByteSize(value: string): number {
  const [, encoded = ''] = value.split(',', 2)
  const sanitized = encoded.replace(/\s+/g, '')
  const padding = sanitized.endsWith('==') ? 2 : sanitized.endsWith('=') ? 1 : 0
  return Math.max(0, Math.floor((sanitized.length * 3) / 4) - padding)
placeholder

function validateAvatarInput(value: string): string | null {
  const normalized = value.trim()
  if (!normalized) {
    return null
  placeholder

  if (normalized.startsWith('data:')) {
    if (!/^data:image\/[a-zA-Z0-9.+-]+;base64,/i.test(normalized)) {
      appStore.showError(t('profile.avatar.invalidValue'))
      return null
    placeholder
    if (estimateDataURLByteSize(normalized) > maxAvatarBytes) {
      appStore.showError(t('profile.avatar.fileTooLarge'))
      return null
    placeholder
    return normalized
  placeholder

  try {
    const parsed = new URL(normalized)
    if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
      return normalized
    placeholder
  placeholder catch {
    // Invalid URL is handled below.
  placeholder

  appStore.showError(t('profile.avatar.invalidValue'))
  return null
placeholder

function readFileAsDataURL(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(typeof reader.result === 'string' ? reader.result : '')
    reader.onerror = () => reject(reader.error ?? new Error('avatar_read_failed'))
    reader.readAsDataURL(file)
  placeholder)
placeholder

function loadImage(dataURL: string): Promise<HTMLImageElement> {
  return new Promise((resolve, reject) => {
    const image = new Image()
    image.onload = () => resolve(image)
    image.onerror = () => reject(new Error(t('profile.avatar.readFailed')))
    image.src = dataURL
  placeholder)
placeholder

function canvasToBlob(canvas: HTMLCanvasElement, type: string, quality: number): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob((blob) => {
      if (!blob) {
        reject(new Error(t('profile.avatar.compressFailed')))
        return
      placeholder
      resolve(blob)
    placeholder, type, quality)
  placeholder)
placeholder

async function compressAvatarFile(file: File): Promise<File> {
  const sourceDataURL = await readFileAsDataURL(file)
  const image = await loadImage(sourceDataURL)
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error(t('profile.avatar.compressFailed'))
  placeholder

  for (const scale of avatarScaleSteps) {
    const width = Math.max(1, Math.round(image.naturalWidth * scale))
    const height = Math.max(1, Math.round(image.naturalHeight * scale))
    canvas.width = width
    canvas.height = height
    ctx.clearRect(0, 0, width, height)
    ctx.drawImage(image, 0, 0, width, height)

    for (const quality of avatarQualitySteps) {
      const blob = await canvasToBlob(canvas, 'image/webp', quality)
      if (blob.size <= targetAvatarUploadBytes) {
        const fileName = file.name.replace(/\.[^.]+$/, '') || 'avatar'
        return new File([blob], `${fileNameplaceholder.webp`, { type: 'image/webp' placeholder)
      placeholder
    placeholder
  placeholder

  throw new Error(t('profile.avatar.compressTooLarge'))
placeholder

async function prepareAvatarUpload(file: File): Promise<File> {
  if (!file.type.startsWith('image/')) {
    throw new Error(t('profile.avatar.invalidType'))
  placeholder
  if (file.type === 'image/gif') {
    if (file.size > targetAvatarUploadBytes) {
      throw new Error(t('profile.avatar.gifTooLarge'))
    placeholder
    return file
  placeholder
  if (file.size <= targetAvatarUploadBytes) {
    return file
  placeholder
  return compressAvatarFile(file)
placeholder

async function handleAvatarFileChange(event: Event) {
  const input = event.target as HTMLInputElement | null
  const file = input?.files?.[0]
  if (input) {
    input.value = ''
  placeholder
  if (!file) {
    return
  placeholder
  if (!file.type.startsWith('image/')) {
    appStore.showError(t('profile.avatar.invalidType'))
    return
  placeholder
  if (file.size > maxAvatarBytes) {
    appStore.showError(t('profile.avatar.fileTooLarge'))
    return
  placeholder

  try {
    const preparedFile = await prepareAvatarUpload(file)
    const dataURL = await readFileAsDataURL(preparedFile)
    const normalized = validateAvatarInput(dataURL)
    if (!normalized) {
      return
    placeholder
    avatarDraft.value = normalized
  placeholder catch (error: unknown) {
    appStore.showError(extractApiErrorMessage(error, t('common.error')))
  placeholder
placeholder

async function handleAvatarSave() {
  const normalized = validateAvatarInput(avatarDraft.value)
  if (!normalized) {
    return
  placeholder

  avatarSaving.value = true
  try {
    const updated = await userAPI.updateProfile({ avatar_url: normalized placeholder)
    authStore.user = updated
    avatarDraft.value = updated.avatar_url?.trim() || ''
    appStore.showSuccess(t('profile.avatar.saveSuccess'))
  placeholder catch (error: unknown) {
    appStore.showError(extractApiErrorMessage(error, t('common.error')))
  placeholder finally {
    avatarSaving.value = false
  placeholder
placeholder

async function handleAvatarDelete() {
  if (avatarSaving.value) {
    return
  placeholder
  if (!avatarDraft.value.trim() && !props.user?.avatar_url?.trim()) {
    appStore.showError(t('profile.avatar.emptyDeleteHint'))
    return
  placeholder

  avatarSaving.value = true
  try {
    const updated = await userAPI.updateProfile({ avatar_url: '' placeholder)
    authStore.user = updated
    avatarDraft.value = ''
    appStore.showSuccess(t('profile.avatar.deleteSuccess'))
  placeholder catch (error: unknown) {
    appStore.showError(extractApiErrorMessage(error, t('common.error')))
  placeholder finally {
    avatarSaving.value = false
  placeholder
placeholder
</script>
