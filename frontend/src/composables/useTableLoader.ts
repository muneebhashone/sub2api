import { ref, reactive, onUnmounted, toRaw placeholder from 'vue'
import { useDebounceFn placeholder from '@vueuse/core'
import type { BasePaginationResponse, FetchOptions placeholder from '@/types'

interface PaginationState {
  page: number
  page_size: number
  total: number
  pages: number
placeholder

interface TableLoaderOptions<T, P> {
  fetchFn: (page: number, pageSize: number, params: P, options?: FetchOptions) => Promise<BasePaginationResponse<T>>
  initialParams?: P
  pageSize?: number
  debounceMs?: number
placeholder

/**
 * 通用表格数据加载 Composable
 * 统一处理分页、筛选、搜索防抖和请求取消
 */
export function useTableLoader<T, P extends Record<string, any>>(options: TableLoaderOptions<T, P>) {
  const { fetchFn, initialParams, pageSize = 20, debounceMs = 300 placeholder = options

  const items = ref<T[]>([])
  const loading = ref(false)
  const params = reactive<P>({ ...(initialParams || {placeholder) placeholder as P)
  const pagination = reactive<PaginationState>({
    page: 1,
    page_size: pageSize,
    total: 0,
    pages: 0
  placeholder)

  let abortController: AbortController | null = null

  const isAbortError = (error: any) => {
    return error?.name === 'AbortError' || error?.code === 'ERR_CANCELED' || error?.name === 'CanceledError'
  placeholder

  const load = async () => {
    if (abortController) {
      abortController.abort()
    placeholder
    abortController = new AbortController()
    loading.value = true

    try {
      const response = await fetchFn(
        pagination.page,
        pagination.page_size,
        toRaw(params) as P,
        { signal: abortController.signal placeholder
      )
      
      items.value = response.items || []
      pagination.total = response.total || 0
      pagination.pages = response.pages || 0
    placeholder catch (error) {
      if (!isAbortError(error)) {
        console.error('Table load error:', error)
        throw error
      placeholder
    placeholder finally {
      if (abortController && !abortController.signal.aborted) {
        loading.value = false
      placeholder
    placeholder
  placeholder

  const reload = () => {
    pagination.page = 1
    return load()
  placeholder

  const debouncedReload = useDebounceFn(reload, debounceMs)

  const handlePageChange = (page: number) => {
    pagination.page = page
    load()
  placeholder

  const handlePageSizeChange = (size: number) => {
    pagination.page_size = size
    pagination.page = 1
    load()
  placeholder

  onUnmounted(() => {
    abortController?.abort()
  placeholder)

  return {
    items,
    loading,
    params,
    pagination,
    load,
    reload,
    debouncedReload,
    handlePageChange,
    handlePageSizeChange
  placeholder
placeholder
