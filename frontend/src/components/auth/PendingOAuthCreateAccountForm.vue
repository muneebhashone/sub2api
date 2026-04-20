<template>
  <form class="space-y-3" @submit.prevent="handleSubmit">
    <input
      v-model="email"
      :data-testid="`${testIdPrefixplaceholder-create-account-email`"
      type="email"
      class="input w-full"
      placeholder="you@example.com"
      :disabled="isSubmitting || isSendingCode"
    />
    <input
      v-model="password"
      :data-testid="`${testIdPrefixplaceholder-create-account-password`"
      type="password"
      class="input w-full"
      placeholder="Password"
      :disabled="isSubmitting"
    />
    <div v-if="turnstileEnabled && turnstileSiteKey" class="space-y-2">
      <TurnstileWidget
        ref="turnstileRef"
        :site-key="turnstileSiteKey"
        @verify="onTurnstileVerify"
        @expire="onTurnstileExpire"
        @error="onTurnstileError"
      />
    </div>
    <div class="flex gap-3">
      <input
        v-model="verifyCode"
        :data-testid="`${testIdPrefixplaceholder-create-account-verify-code`"
        type="text"
        inputmode="numeric"
        maxlength="6"
        class="input min-w-0 flex-1"
        placeholder="123456"
        :disabled="isSubmitting"
      />
      <button
        :data-testid="`${testIdPrefixplaceholder-create-account-send-code`"
        type="button"
        class="btn btn-secondary shrink-0"
        :disabled="isSubmitting || isSendingCode || countdown > 0 || !email.trim() || (turnstileEnabled && !turnstileToken)"
        @click="handleSendCode"
      >
        {{
          isSendingCode
            ? t('auth.sendingCode')
            : countdown > 0
              ? t('auth.resendCountdown', { countdown placeholder)
              : t('auth.sendCode')
        placeholderplaceholder
      </button>
    </div>
    <p v-if="sendCodeSuccess" class="text-sm text-green-600 dark:text-green-400">
      {{ t('auth.codeSentSuccess') placeholderplaceholder
    </p>
    <p v-else class="text-xs text-gray-500 dark:text-dark-400">
      {{ t('auth.verificationCodeHint') placeholderplaceholder
    </p>
    <button
      :data-testid="`${testIdPrefixplaceholder-create-account-submit`"
      type="button"
      class="btn btn-primary w-full"
      :disabled="isSubmitting || !email.trim() || password.length < 6"
      @click="handleSubmit"
    >
      {{ isSubmitting ? t('common.processing') : 'Create account' placeholderplaceholder
    </button>
    <button
      type="button"
      class="btn btn-secondary w-full"
      :disabled="isSubmitting"
      @click="emitSwitchToBind"
    >
      I already have an account
    </button>
    <transition name="fade">
      <p v-if="sendCodeError" class="text-sm text-red-600 dark:text-red-400">
        {{ sendCodeError placeholderplaceholder
      </p>
    </transition>
    <transition name="fade">
      <p v-if="errorMessage" class="text-sm text-red-600 dark:text-red-400">
        {{ errorMessage placeholderplaceholder
      </p>
    </transition>
  </form>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import TurnstileWidget from '@/components/TurnstileWidget.vue'
import { getPublicSettings, sendVerifyCode placeholder from '@/api/auth'

export type PendingOAuthCreateAccountPayload = {
  email: string
  password: string
  verifyCode: string
placeholder

const props = defineProps<{
  initialEmail: string
  testIdPrefix: string
  isSubmitting: boolean
  errorMessage?: string
placeholder>()

const emit = defineEmits<{
  submit: [payload: PendingOAuthCreateAccountPayload]
  switchToBind: [email: string]
placeholder>()

const { t placeholder = useI18n()

const email = ref('')
const password = ref('')
const verifyCode = ref('')
const isSendingCode = ref(false)
const sendCodeError = ref('')
const sendCodeSuccess = ref(false)
const countdown = ref(0)
const turnstileEnabled = ref(false)
const turnstileSiteKey = ref('')
const turnstileToken = ref('')
const turnstileRef = ref<InstanceType<typeof TurnstileWidget> | null>(null)

let countdownTimer: ReturnType<typeof setInterval> | null = null

watch(
  () => props.initialEmail,
  value => {
    email.value = value || ''
  placeholder,
  { immediate: true placeholder
)

function clearCountdown() {
  if (countdownTimer) {
    clearInterval(countdownTimer)
    countdownTimer = null
  placeholder
placeholder

function startCountdown(seconds: number) {
  clearCountdown()
  countdown.value = Math.max(0, seconds)

  if (countdown.value <= 0) {
    return
  placeholder

  countdownTimer = setInterval(() => {
    if (countdown.value <= 1) {
      countdown.value = 0
      clearCountdown()
      return
    placeholder

    countdown.value -= 1
  placeholder, 1000)
placeholder

function getRequestErrorMessage(error: unknown, fallback: string): string {
  const err = error as { message?: string; response?: { data?: { detail?: string; message?: string placeholder placeholder placeholder
  return err.response?.data?.detail || err.response?.data?.message || err.message || fallback
placeholder

function resetTurnstile() {
  turnstileToken.value = ''
  turnstileRef.value?.reset()
placeholder

function onTurnstileVerify(token: string) {
  turnstileToken.value = token
  sendCodeError.value = ''
placeholder

function onTurnstileExpire() {
  turnstileToken.value = ''
  sendCodeError.value = t('auth.turnstileExpired')
placeholder

function onTurnstileError() {
  turnstileToken.value = ''
  sendCodeError.value = t('auth.turnstileFailed')
placeholder

async function handleSendCode() {
  const trimmedEmail = email.value.trim()
  if (!trimmedEmail) {
    return
  placeholder

  if (turnstileEnabled.value && !turnstileToken.value) {
    sendCodeError.value = t('auth.completeVerification')
    return
  placeholder

  isSendingCode.value = true
  sendCodeError.value = ''
  sendCodeSuccess.value = false

  try {
    const response = await sendVerifyCode({
      email: trimmedEmail,
      turnstile_token: turnstileEnabled.value ? turnstileToken.value : undefined
    placeholder)
    sendCodeSuccess.value = true
    startCountdown(response.countdown)
    if (turnstileEnabled.value) {
      resetTurnstile()
    placeholder
  placeholder catch (error: unknown) {
    sendCodeError.value = getRequestErrorMessage(error, t('auth.sendCodeFailed'))
  placeholder finally {
    isSendingCode.value = false
  placeholder
placeholder

function handleSubmit() {
  const trimmedEmail = email.value.trim()
  if (!trimmedEmail || password.value.length < 6) {
    return
  placeholder

  emit('submit', {
    email: trimmedEmail,
    password: password.value,
    verifyCode: verifyCode.value.trim()
  placeholder)
placeholder

function emitSwitchToBind() {
  emit('switchToBind', email.value.trim())
placeholder

onMounted(async () => {
  try {
    const settings = await getPublicSettings()
    turnstileEnabled.value = settings.turnstile_enabled === true
    turnstileSiteKey.value = settings.turnstile_site_key || ''
  placeholder catch {
    turnstileEnabled.value = false
    turnstileSiteKey.value = ''
  placeholder
placeholder)

onUnmounted(() => {
  clearCountdown()
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
