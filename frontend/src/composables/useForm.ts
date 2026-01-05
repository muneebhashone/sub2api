import { ref placeholder from 'vue'
import { useAppStore placeholder from '@/stores/app'

interface UseFormOptions<T> {
  form: T
  submitFn: (data: T) => Promise<void>
  successMsg?: string
  errorMsg?: string
placeholder

/**
 * 统一表单提交逻辑
 * 管理加载状态、错误捕获及通知
 */
export function useForm<T>(options: UseFormOptions<T>) {
  const { form, submitFn, successMsg, errorMsg placeholder = options
  const loading = ref(false)
  const appStore = useAppStore()

  const submit = async () => {
    if (loading.value) return
    
    loading.value = true
    try {
      await submitFn(form)
      if (successMsg) {
        appStore.showSuccess(successMsg)
      placeholder
    placeholder catch (error: any) {
      const detail = error.response?.data?.detail || error.response?.data?.message || error.message
      appStore.showError(errorMsg || detail)
      // 继续抛出错误，让组件有机会进行局部处理（如验证错误显示）
      throw error
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  return {
    loading,
    submit
  placeholder
placeholder
