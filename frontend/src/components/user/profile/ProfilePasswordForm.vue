<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-medium text-gray-900 dark:text-white">
        {{ t('profile.changePassword') placeholderplaceholder
      </h2>
    </div>
    <div class="px-6 py-6">
      <form @submit.prevent="handleChangePassword" class="space-y-4">
        <div>
          <label for="old_password" class="input-label">
            {{ t('profile.currentPassword') placeholderplaceholder
          </label>
          <input
            id="old_password"
            v-model="form.old_password"
            type="password"
            required
            autocomplete="current-password"
            class="input"
          />
        </div>

        <div>
          <label for="new_password" class="input-label">
            {{ t('profile.newPassword') placeholderplaceholder
          </label>
          <input
            id="new_password"
            v-model="form.new_password"
            type="password"
            required
            autocomplete="new-password"
            class="input"
          />
          <p class="input-hint">
            {{ t('profile.passwordHint') placeholderplaceholder
          </p>
        </div>

        <div>
          <label for="confirm_password" class="input-label">
            {{ t('profile.confirmNewPassword') placeholderplaceholder
          </label>
          <input
            id="confirm_password"
            v-model="form.confirm_password"
            type="password"
            required
            autocomplete="new-password"
            class="input"
          />
          <p
            v-if="form.new_password && form.confirm_password && form.new_password !== form.confirm_password"
            class="input-error-text"
          >
            {{ t('profile.passwordsNotMatch') placeholderplaceholder
          </p>
        </div>

        <div class="flex justify-end pt-4">
          <button type="submit" :disabled="loading" class="btn btn-primary">
            {{ loading ? t('profile.changingPassword') : t('profile.changePasswordButton') placeholderplaceholder
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { userAPI placeholder from '@/api'

const { t placeholder = useI18n()
const appStore = useAppStore()

const loading = ref(false)
const form = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
placeholder)

const handleChangePassword = async () => {
  if (form.value.new_password !== form.value.confirm_password) {
    appStore.showError(t('profile.passwordsNotMatch'))
    return
  placeholder

  if (form.value.new_password.length < 8) {
    appStore.showError(t('profile.passwordTooShort'))
    return
  placeholder

  loading.value = true
  try {
    await userAPI.changePassword(form.value.old_password, form.value.new_password)
    form.value = { old_password: '', new_password: '', confirm_password: '' placeholder
    appStore.showSuccess(t('profile.passwordChangeSuccess'))
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('profile.passwordChangeFailed'))
  placeholder finally {
    loading.value = false
  placeholder
placeholder
</script>
