import { ref, reactive, onUnmounted placeholder from 'vue'
import { useDebounceFn placeholder from '@vueuse/core'

interface PaginationState {
  page: number
  page_size: number
  total: number
  pages: number
placeholder

interface TableLoaderOptions<T, P> {
  fetchFn: (page: number, pageSize: number, params: P, options?: { signal: AbortSignal placeholder) => Promise<{
    items: T[]
    total: number
    pages: number
  placeholder>
  initialParams?: P
  pageSize?: number
  debounceMs?: number
placeholder

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
    return error?.name === 'AbortError' || error?.code === 'ERR_CANCELED'
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
        params,
        { signal: abortController.signal placeholder
      )
      
      items.value = response.items
      pagination.total = response.total
      pagination.pages = response.pages
    placeholder catch (error) {
      if (!isAbortError(error)) {
        throw error
      placeholder
    placeholder finally {
      if (abortController?.signal.aborted === false) {
        loading.value = false
      placeholder
    placeholder
  placeholder

  const reload = () => {
    pagination.page = 1
    return load()
  placeholder

  const debouncedLoad = useDebounceFn(reload, debounceMs)

  const handlePageChange = (page: number) => {
    pagination.page = page
    load()
  placeholder

  const handlePageSizeChange = (size: number) => {
    pagination.page_size = size
    reload()
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
    debouncedLoad,
    handlePageChange,
    handlePageSizeChange
  placeholder
placeholder
