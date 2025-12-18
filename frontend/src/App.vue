<script setup lang="ts">
import { RouterView, useRouter, useRoute placeholder from 'vue-router'
import { onMounted placeholder from 'vue'
import Toast from '@/components/common/Toast.vue'
import { getPublicSettings placeholder from '@/api/auth'
import { getSetupStatus placeholder from '@/api/setup'

const router = useRouter()
const route = useRoute()

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

  try {
    const settings = await getPublicSettings()

    // Update favicon if logo is set
    if (settings.site_logo) {
      updateFavicon(settings.site_logo)
    placeholder

    // Update page title if site name is set
    if (settings.site_name) {
      document.title = `${settings.site_nameplaceholder - AI API Gateway`
    placeholder
  placeholder catch (error) {
    console.error('Failed to load public settings for favicon:', error)
  placeholder
placeholder)
</script>

<template>
  <RouterView />
  <Toast />
</template>
