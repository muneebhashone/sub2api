<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-medium text-gray-900 dark:text-white">
        {{ t('profile.editProfile') placeholderplaceholder
      </h2>
    </div>
    <div class="px-6 py-6">
      <form @submit.prevent="handleUpdateProfile" class="space-y-4">
        <div>
          <label for="username" class="input-label">
            {{ t('profile.username') placeholderplaceholder
          </label>
          <input
            id="username"
            v-model="username"
            type="text"
            class="input"
            :placeholder="t('profile.enterUsername')"
          />
        </div>

        <div class="flex justify-end pt-4">
          <button type="submit" :disabled="loading" class="btn btn-primary">
            {{ loading ? t('profile.updating') : t('profile.updateProfile') placeholderplaceholder
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAuthStore placeholder from '@/stores/auth'
import { useAppStore placeholder from '@/stores/app'
import { userAPI placeholder from '@/api'

const props = defineProps<{
  initialUsername: string
placeholder>()

const { t placeholder = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()

const username = ref(props.initialUsername)
const loading = ref(false)

watch(() => props.initialUsername, (val) => {
  username.value = val
placeholder)

const handleUpdateProfile = async () => {
  if (!username.value.trim()) {
    appStore.showError(t('profile.usernameRequired'))
    return
  placeholder

  loading.value = true
  try {
    const updatedUser = await userAPI.updateProfile({
      username: username.value
    placeholder)
    authStore.user = updatedUser
    appStore.showSuccess(t('profile.updateSuccess'))
  placeholder catch (error: any) {
    appStore.showError(error.response?.data?.detail || t('profile.updateFailed'))
  placeholder finally {
    loading.value = false
  placeholder
placeholder
</script>
