<script setup lang="ts">
import { RouterView, useRouter, useRoute placeholder from 'vue-router'
import { onMounted, watch placeholder from 'vue'
import Toast from '@/components/common/Toast.vue'
import { useAppStore placeholder from '@/stores'
import { getSetupStatus placeholder from '@/api/setup'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()

/**
 * Update favicon dynamically
 * @param logoUrl - URL of the logo to use as favicon
 */
function updateFavicon(logoUrl: string) {
  // Find existing favicon link or create new one
  let link = document.querySelector<HTMLLinkElement>('link[rel="icon"]')
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  placeholder
  link.type = logoUrl.endsWith('.svg') ? 'image/svg+xml' : 'image/x-icon'
  link.href = logoUrl
placeholder

// Watch for site settings changes and update favicon/title
watch(
  () => appStore.siteLogo,
  (newLogo) => {
    if (newLogo) {
      updateFavicon(newLogo)
    placeholder
  placeholder,
  { immediate: true placeholder
)

watch(
  () => appStore.siteName,
  (newName) => {
    if (newName) {
      document.title = `${newNameplaceholder - AI API Gateway`
    placeholder
  placeholder,
  { immediate: true placeholder
)

onMounted(async () => {
  // Check if setup is needed
  try {
    const status = await getSetupStatus()
    if (status.needs_setup && route.path !== '/setup') {
      router.replace('/setup')
      return
    placeholder
  placeholder catch {
    // If setup endpoint fails, assume normal mode and continue
  placeholder

  // Load public settings into appStore (will be cached for other components)
  await appStore.fetchPublicSettings()
placeholder)
</script>

<template>
  <RouterView />
  <Toast />
</template>
