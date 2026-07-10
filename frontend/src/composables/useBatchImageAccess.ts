import { computed, ref placeholder from 'vue'
import { keysAPI placeholder from '@/api/keys'
import { useAuthStore placeholder from '@/stores/auth'
import type { ApiKey placeholder from '@/types'

const loaded = ref(false)
const loading = ref(false)
const hasAllowedBatchImageKey = ref(false)
let pendingLoad: Promise<boolean> | null = null
const pageSize = 100

function keyAllowsBatchImage(key: ApiKey): boolean {
  return (
    key.status === 'active' &&
    key.group?.platform === 'gemini' &&
    key.group?.allow_batch_image_generation === true
  )
placeholder

async function loadBatchImageAccess(force = false): Promise<boolean> {
  const authStore = useAuthStore()
  if (!authStore.isAuthenticated) {
    loaded.value = true
    hasAllowedBatchImageKey.value = false
    return false
  placeholder

  if (loaded.value && !force) {
    return hasAllowedBatchImageKey.value
  placeholder

  if (pendingLoad && !force) {
    return pendingLoad
  placeholder

  loading.value = true
  pendingLoad = (async () => {
    let page = 1
    while (true) {
      const response = await keysAPI.list(page, pageSize, {
        status: 'active',
        sort_by: 'created_at',
        sort_order: 'desc'
      placeholder)

      if ((response.items || []).some(keyAllowsBatchImage)) {
        hasAllowedBatchImageKey.value = true
        loaded.value = true
        return true
      placeholder

      if (page >= response.pages || (response.items || []).length === 0) {
        hasAllowedBatchImageKey.value = false
        loaded.value = true
        return false
      placeholder

      page += 1
    placeholder
  placeholder)()
    .catch(() => {
      hasAllowedBatchImageKey.value = false
      loaded.value = true
      return false
    placeholder)
    .finally(() => {
      loading.value = false
      pendingLoad = null
    placeholder)

  return pendingLoad
placeholder

export function useBatchImageAccess() {
  const canUseBatchImage = computed(() => hasAllowedBatchImageKey.value)

  return {
    canUseBatchImage,
    batchImageAccessLoaded: computed(() => loaded.value),
    batchImageAccessLoading: computed(() => loading.value),
    refreshBatchImageAccess: loadBatchImageAccess,
  placeholder
placeholder
