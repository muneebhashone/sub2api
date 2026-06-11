import { defineStore placeholder from 'pinia'
import { ref placeholder from 'vue'
import { adminAPI placeholder from '@/api/admin'
import type { Module placeholder from '@/types'

export const useModulesStore = defineStore('modules', () => {
  // State
  const modules = ref<Module[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Actions
  async function fetchModules() {
    loading.value = true
    error.value = null
    try {
      const res = await adminAPI.modules.list()
      modules.value = res.modules ?? []
    placeholder catch (err: any) {
      error.value = err?.response?.data?.message || err?.message || 'Failed to load modules'
      throw err
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  return {
    // State
    modules,
    loading,
    error,
    // Actions
    fetchModules,
  placeholder
placeholder)
