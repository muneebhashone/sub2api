<template>
  <AuthLayout>
    <div class="space-y-6">
      <!-- Title -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">
          {{ t('auth.resetPasswordTitle') placeholderplaceholder
        </h2>
        <p class="mt-2 text-sm text-gray-500 dark:text-dark-400">
          {{ t('auth.resetPasswordHint') placeholderplaceholder
        </p>
      </div>

      <!-- Invalid Link State -->
      <div v-if="isInvalidLink" class="space-y-6">
        <div class="rounded-xl border border-red-200 bg-red-50 p-6 dark:border-red-800/50 dark:bg-red-900/20">
          <div class="flex flex-col items-center gap-4 text-center">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-800/50">
              <Icon name="exclamationCircle" size="lg" class="text-red-600 dark:text-red-400" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-red-800 dark:text-red-200">
                {{ t('auth.invalidResetLink') placeholderplaceholder
              </h3>
              <p class="mt-2 text-sm text-red-700 dark:text-red-300">
                {{ t('auth.invalidResetLinkHint') placeholderplaceholder
              </p>
            </div>
          </div>
        </div>

        <div class="text-center">
          <router-link
            to="/forgot-password"
            class="inline-flex items-center gap-2 font-medium text-primary-600 transition-colors hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300"
          >
            {{ t('auth.requestNewResetLink') placeholderplaceholder
          </router-link>
        </div>
      </div>

      <!-- Success State -->
      <div v-else-if="isSuccess" class="space-y-6">
        <div class="rounded-xl border border-green-200 bg-green-50 p-6 dark:border-green-800/50 dark:bg-green-900/20">
          <div class="flex flex-col items-center gap-4 text-center">
            <div class="flex h-12 w-12 items-center justify-center rounded-full bg-green-100 dark:bg-green-800/50">
              <Icon name="checkCircle" size="lg" class="text-green-600 dark:text-green-400" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-green-800 dark:text-green-200">
                {{ t('auth.passwordResetSuccess') placeholderplaceholder
              </h3>
              <p class="mt-2 text-sm text-green-700 dark:text-green-300">
                {{ t('auth.passwordResetSuccessHint') placeholderplaceholder
              </p>
            </div>
          </div>
        </div>

        <div class="text-center">
          <router-link
            to="/login"
            class="btn btn-primary inline-flex items-center gap-2"
          >
            <Icon name="login" size="md" />
            {{ t('auth.signIn') placeholderplaceholder
          </router-link>
        </div>
      </div>

      <!-- Form State -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-5">
        <!-- Email (readonly) -->
        <div>
          <label for="email" class="input-label">
            {{ t('auth.emailLabel') placeholderplaceholder
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="mail" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="email"
              :value="email"
              type="email"
              readonly
              disabled
              class="input pl-11 bg-gray-50 dark:bg-dark-700"
            />
          </div>
        </div>

        <!-- New Password Input -->
        <div>
          <label for="password" class="input-label">
            {{ t('auth.newPassword') placeholderplaceholder
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="password"
              v-model="formData.password"
              :type="showPassword ? 'text' : 'password'"
              required
              autocomplete="new-password"
              :disabled="isLoading"
              class="input pl-11 pr-11"
              :class="{ 'input-error': errors.password placeholder"
              :placeholder="t('auth.newPasswordPlaceholder')"
            />
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300"
            >
              <Icon v-if="showPassword" name="eyeOff" size="md" />
              <Icon v-else name="eye" size="md" />
            </button>
          </div>
          <p v-if="errors.password" class="input-error-text">
            {{ errors.password placeholderplaceholder
          </p>
        </div>

        <!-- Confirm Password Input -->
        <div>
          <label for="confirmPassword" class="input-label">
            {{ t('auth.confirmPassword') placeholderplaceholder
          </label>
          <div class="relative">
            <div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
              <Icon name="lock" size="md" class="text-gray-400 dark:text-dark-500" />
            </div>
            <input
              id="confirmPassword"
              v-model="formData.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              required
              autocomplete="new-password"
              :disabled="isLoading"
              class="input pl-11 pr-11"
              :class="{ 'input-error': errors.confirmPassword placeholder"
              :placeholder="t('auth.confirmPasswordPlaceholder')"
            />
            <button
              type="button"
              @click="showConfirmPassword = !showConfirmPassword"
              class="absolute inset-y-0 right-0 flex items-center pr-3.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-dark-300"
            >
              <Icon v-if="showConfirmPassword" name="eyeOff" size="md" />
              <Icon v-else name="eye" size="md" />
            </button>
          </div>
          <p v-if="errors.confirmPassword" class="input-error-text">
            {{ errors.confirmPassword placeholderplaceholder
          </p>
        </div>

        <!-- Error Message -->
        <transition name="fade">
          <div
            v-if="errorMessage"
            class="rounded-xl border border-red-200 bg-red-50 p-4 dark:border-red-800/50 dark:bg-red-900/20"
          >
            <div class="flex items-start gap-3">
              <div class="flex-shrink-0">
                <Icon name="exclamationCircle" size="md" class="text-red-500" />
              </div>
              <p class="text-sm text-red-700 dark:text-red-400">
                {{ errorMessage placeholderplaceholder
              </p>
            </div>
          </div>
        </transition>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isLoading"
          class="btn btn-primary w-full"
        >
          <svg
            v-if="isLoading"
            class="-ml-1 mr-2 h-4 w-4 animate-spin text-white"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <Icon v-else name="checkCircle" size="md" class="mr-2" />
          {{ isLoading ? t('auth.resettingPassword') : t('auth.resetPassword') placeholderplaceholder
        </button>
      </form>
    </div>

    <!-- Footer -->
    <template #footer>
      <p class="text-gray-500 dark:text-dark-400">
        {{ t('auth.rememberedPassword') placeholderplaceholder
        <router-link
          to="/login"
          class="font-medium text-primary-600 transition-colors hover:text-primary-500 dark:text-primary-400 dark:hover:text-primary-300"
        >
          {{ t('auth.signIn') placeholderplaceholder
        </router-link>
      </p>
    </template>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted placeholder from 'vue'
import { useRoute placeholder from 'vue-router'
import { useI18n placeholder from 'vue-i18n'
import { AuthLayout placeholder from '@/components/layout'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore placeholder from '@/stores'
import { resetPassword placeholder from '@/api/auth'

const { t placeholder = useI18n()

// ==================== Router & Stores ====================

const route = useRoute()
const appStore = useAppStore()

// ==================== State ====================

const isLoading = ref<boolean>(false)
const isSuccess = ref<boolean>(false)
const errorMessage = ref<string>('')
const showPassword = ref<boolean>(false)
const showConfirmPassword = ref<boolean>(false)

// URL parameters
const email = ref<string>('')
const token = ref<string>('')

const formData = reactive({
  password: '',
  confirmPassword: ''
placeholder)

const errors = reactive({
  password: '',
  confirmPassword: ''
placeholder)

// Check if the reset link is valid (has email and token)
const isInvalidLink = computed(() => !email.value || !token.value)

// ==================== Lifecycle ====================

onMounted(() => {
  // Get email and token from URL query parameters
  email.value = (route.query.email as string) || ''
  token.value = (route.query.token as string) || ''
placeholder)

// ==================== Validation ====================

function validateForm(): boolean {
  errors.password = ''
  errors.confirmPassword = ''

  let isValid = true

  // Password validation
  if (!formData.password) {
    errors.password = t('auth.passwordRequired')
    isValid = false
  placeholder else if (formData.password.length < 6) {
    errors.password = t('auth.passwordMinLength')
    isValid = false
  placeholder

  // Confirm password validation
  if (!formData.confirmPassword) {
    errors.confirmPassword = t('auth.confirmPasswordRequired')
    isValid = false
  placeholder else if (formData.password !== formData.confirmPassword) {
    errors.confirmPassword = t('auth.passwordsDoNotMatch')
    isValid = false
  placeholder

  return isValid
placeholder

// ==================== Form Handlers ====================

async function handleSubmit(): Promise<void> {
  errorMessage.value = ''

  if (!validateForm()) {
    return
  placeholder

  isLoading.value = true

  try {
    await resetPassword({
      email: email.value,
      token: token.value,
      new_password: formData.password
    placeholder)

    isSuccess.value = true
    appStore.showSuccess(t('auth.passwordResetSuccess'))
  placeholder catch (error: unknown) {
    const err = error as { message?: string; response?: { data?: { detail?: string; code?: string placeholder placeholder placeholder

    // Check for invalid/expired token error
    if (err.response?.data?.code === 'INVALID_RESET_TOKEN') {
      errorMessage.value = t('auth.invalidOrExpiredToken')
    placeholder else if (err.response?.data?.detail) {
      errorMessage.value = err.response.data.detail
    placeholder else if (err.message) {
      errorMessage.value = err.message
    placeholder else {
      errorMessage.value = t('auth.resetPasswordFailed')
    placeholder

    appStore.showError(errorMessage.value)
  placeholder finally {
    isLoading.value = false
  placeholder
placeholder
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
