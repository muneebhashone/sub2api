import { describe, it, expect, vi, beforeEach placeholder from 'vitest'
import { setActivePinia, createPinia placeholder from 'pinia'
import { useForm placeholder from '@/composables/useForm'
import { useAppStore placeholder from '@/stores/app'

// Mock API 依赖（app store 内部引用了这些）
vi.mock('@/api/admin/system', () => ({
  checkUpdates: vi.fn(),
placeholder))
vi.mock('@/api/auth', () => ({
  getPublicSettings: vi.fn(),
placeholder))

describe('useForm', () => {
  let appStore: ReturnType<typeof useAppStore>

  beforeEach(() => {
    setActivePinia(createPinia())
    appStore = useAppStore()
    vi.clearAllMocks()
  placeholder)

  it('submit 期间 loading 为 true，完成后为 false', async () => {
    let resolveSubmit: () => void
    const submitFn = vi.fn(
      () => new Promise<void>((resolve) => { resolveSubmit = resolve placeholder)
    )

    const { loading, submit placeholder = useForm({
      form: { name: 'test' placeholder,
      submitFn,
    placeholder)

    expect(loading.value).toBe(false)

    const submitPromise = submit()
    // 提交中
    expect(loading.value).toBe(true)

    resolveSubmit!()
    await submitPromise

    expect(loading.value).toBe(false)
  placeholder)

  it('submit 成功时显示成功消息', async () => {
    const submitFn = vi.fn().mockResolvedValue(undefined)
    const showSuccessSpy = vi.spyOn(appStore, 'showSuccess')

    const { submit placeholder = useForm({
      form: { name: 'test' placeholder,
      submitFn,
      successMsg: '保存成功',
    placeholder)

    await submit()

    expect(showSuccessSpy).toHaveBeenCalledWith('保存成功')
  placeholder)

  it('submit 成功但无 successMsg 时不调用 showSuccess', async () => {
    const submitFn = vi.fn().mockResolvedValue(undefined)
    const showSuccessSpy = vi.spyOn(appStore, 'showSuccess')

    const { submit placeholder = useForm({
      form: { name: 'test' placeholder,
      submitFn,
    placeholder)

    await submit()

    expect(showSuccessSpy).not.toHaveBeenCalled()
  placeholder)

  it('submit 失败时显示错误消息并抛出错误', async () => {
    const error = Object.assign(new Error('提交失败'), {
      response: { data: { message: '服务器错误' placeholder placeholder,
    placeholder)
    const submitFn = vi.fn().mockRejectedValue(error)
    const showErrorSpy = vi.spyOn(appStore, 'showError')

    const { submit, loading placeholder = useForm({
      form: { name: 'test' placeholder,
      submitFn,
    placeholder)

    await expect(submit()).rejects.toThrow('提交失败')

    expect(showErrorSpy).toHaveBeenCalled()
    expect(loading.value).toBe(false)
  placeholder)

  it('submit 失败时使用自定义 errorMsg', async () => {
    const submitFn = vi.fn().mockRejectedValue(new Error('network'))
    const showErrorSpy = vi.spyOn(appStore, 'showError')

    const { submit placeholder = useForm({
      form: { name: 'test' placeholder,
      submitFn,
      errorMsg: '自定义错误提示',
    placeholder)

    await expect(submit()).rejects.toThrow()

    expect(showErrorSpy).toHaveBeenCalledWith('自定义错误提示')
  placeholder)

  it('loading 中不会重复提交', async () => {
    let resolveSubmit: () => void
    const submitFn = vi.fn(
      () => new Promise<void>((resolve) => { resolveSubmit = resolve placeholder)
    )

    const { submit placeholder = useForm({
      form: { name: 'test' placeholder,
      submitFn,
    placeholder)

    // 第一次提交
    const p1 = submit()
    // 第二次提交（应被忽略，因为 loading=true）
    submit()

    expect(submitFn).toHaveBeenCalledTimes(1)

    resolveSubmit!()
    await p1
  placeholder)

  it('传递 form 数据到 submitFn', async () => {
    const formData = { name: 'test', email: 'test@example.com' placeholder
    const submitFn = vi.fn().mockResolvedValue(undefined)

    const { submit placeholder = useForm({
      form: formData,
      submitFn,
    placeholder)

    await submit()

    expect(submitFn).toHaveBeenCalledWith(formData)
  placeholder)
placeholder)
