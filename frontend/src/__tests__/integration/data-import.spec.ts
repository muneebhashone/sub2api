import { describe, it, expect, vi, beforeEach placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'
import ImportDataModal from '@/components/admin/account/ImportDataModal.vue'

const showError = vi.fn()
const showSuccess = vi.fn()
const showWarning = vi.fn()

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
    showSuccess,
    showWarning
  placeholder)
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      importData: vi.fn()
    placeholder
  placeholder
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  placeholder)
placeholder))

const mountModal = () =>
  mount(ImportDataModal, {
    props: { show: true placeholder,
    global: {
      stubs: {
        BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' placeholder
      placeholder
    placeholder
  placeholder)

const makeJsonFile = (name: string, content: string, type = 'application/json') => {
  const file = new File([content], name, { type placeholder)
  Object.defineProperty(file, 'text', {
    value: () => Promise.resolve(content)
  placeholder)
  return file
placeholder

const setInputFiles = (element: Element, files: File[]) => {
  Object.defineProperty(element, 'files', {
    value: files,
    configurable: true
  placeholder)
placeholder

describe('ImportDataModal', () => {
  beforeEach(async () => {
    showError.mockReset()
    showSuccess.mockReset()
    showWarning.mockReset()
    const { adminAPI placeholder = await import('@/api/admin')
    vi.mocked(adminAPI.accounts.importData).mockReset()
  placeholder)

  it('未选择文件时提示错误', async () => {
    const wrapper = mountModal()

    await wrapper.find('form').trigger('submit')
    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportSelectFile')
  placeholder)

  it('无效 JSON 时按文件名提示解析失败', async () => {
    const { adminAPI placeholder = await import('@/api/admin')
    const wrapper = mountModal()

    const input = wrapper.find('input[type="file"]')
    setInputFiles(input.element, [makeJsonFile('data.json', 'invalid json')])

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportParseFailedFile')
    expect(adminAPI.accounts.importData).not.toHaveBeenCalled()
  placeholder)

  it('不是导出数据的 JSON 按文件名拒绝', async () => {
    const { adminAPI placeholder = await import('@/api/admin')
    const wrapper = mountModal()

    const input = wrapper.find('input[type="file"]')
    setInputFiles(input.element, [makeJsonFile('random.json', JSON.stringify({ name: 'test' placeholder))])

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportInvalidFile')
    expect(adminAPI.accounts.importData).not.toHaveBeenCalled()
  placeholder)

  it('无有效 JSON 的选择不清空已有选择', async () => {
    const { adminAPI placeholder = await import('@/api/admin')
    vi.mocked(adminAPI.accounts.importData).mockResolvedValue({
      proxy_created: 0,
      proxy_reused: 0,
      proxy_failed: 0,
      account_created: 1,
      account_failed: 0
    placeholder)

    const wrapper = mountModal()
    const input = wrapper.find('input[type="file"]')

    const valid = makeJsonFile(
      'valid.json',
      JSON.stringify({ exported_at: '2026-07-05T00:00:00Z', proxies: [], accounts: [{ name: 'a' placeholder] placeholder)
    )
    setInputFiles(input.element, [valid])
    await input.trigger('change')

    setInputFiles(input.element, [new File(['hello'], 'notes.txt', { type: 'text/plain' placeholder)])
    await input.trigger('change')
    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportSelectFile')

    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(adminAPI.accounts.importData).toHaveBeenCalledWith({
      data: expect.objectContaining({
        accounts: [{ name: 'a' placeholder]
      placeholder),
      skip_default_group_bind: true
    placeholder)
  placeholder)

  it('merges multiple selected JSON files before importing', async () => {
    const { adminAPI placeholder = await import('@/api/admin')
    vi.mocked(adminAPI.accounts.importData).mockResolvedValue({
      proxy_created: 0,
      proxy_reused: 0,
      proxy_failed: 0,
      account_created: 2,
      account_failed: 0
    placeholder)

    const wrapper = mountModal()

    const input = wrapper.find('input[type="file"]')
    const first = makeJsonFile(
      'first.json',
      JSON.stringify({ exported_at: '2026-07-05T00:00:00Z', proxies: [], accounts: [{ name: 'a' placeholder] placeholder)
    )
    const second = makeJsonFile(
      'second.json',
      JSON.stringify({
        exported_at: '2026-07-05T00:00:01Z',
        proxies: [{ proxy_key: 'p' placeholder],
        accounts: [{ name: 'b' placeholder]
      placeholder)
    )
    setInputFiles(input.element, [first, second])

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(adminAPI.accounts.importData).toHaveBeenCalledWith({
      data: expect.objectContaining({
        proxies: [{ proxy_key: 'p' placeholder],
        accounts: [{ name: 'a' placeholder, { name: 'b' placeholder]
      placeholder),
      skip_default_group_bind: true
    placeholder)
    expect(showSuccess).toHaveBeenCalledWith('admin.accounts.dataImportSuccess')
  placeholder)

  it('部分成功时关闭弹窗仍通知父组件刷新', async () => {
    const { adminAPI placeholder = await import('@/api/admin')
    vi.mocked(adminAPI.accounts.importData).mockResolvedValue({
      proxy_created: 0,
      proxy_reused: 0,
      proxy_failed: 0,
      account_created: 1,
      account_failed: 1
    placeholder)

    const wrapper = mountModal()
    const input = wrapper.find('input[type="file"]')
    setInputFiles(input.element, [
      makeJsonFile(
        'mixed.json',
        JSON.stringify({
          exported_at: '2026-07-05T00:00:00Z',
          proxies: [],
          accounts: [{ name: 'a' placeholder, { name: 'b' placeholder]
        placeholder)
      )
    ])

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('admin.accounts.dataImportCompletedWithErrors')
    expect(wrapper.emitted('imported')).toBeUndefined()

    // 第二个 btn-secondary 是 footer 的取消按钮(第一个是选择文件)
    await wrapper.findAll('button.btn-secondary')[1]!.trigger('click')

    expect(wrapper.emitted('imported')).toHaveLength(1)
    expect(wrapper.emitted('close')).toHaveLength(1)
  placeholder)
placeholder)
