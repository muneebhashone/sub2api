import { describe, it, expect, vi, beforeEach placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'
import ImportDataModal from '@/components/admin/proxy/ImportDataModal.vue'

const showError = vi.fn()
const showSuccess = vi.fn()

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
    showSuccess
  placeholder)
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    proxies: {
      importData: vi.fn()
    placeholder
  placeholder
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  placeholder)
placeholder))

describe('Proxy ImportDataModal', () => {
  beforeEach(() => {
    showError.mockReset()
    showSuccess.mockReset()
  placeholder)

  it('未选择文件时提示错误', async () => {
    const wrapper = mount(ImportDataModal, {
      props: { show: true placeholder,
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' placeholder
        placeholder
      placeholder
    placeholder)

    await wrapper.find('form').trigger('submit')
    expect(showError).toHaveBeenCalledWith('admin.proxies.dataImportSelectFile')
  placeholder)

  it('无效 JSON 时提示解析失败', async () => {
    const wrapper = mount(ImportDataModal, {
      props: { show: true placeholder,
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' placeholder
        placeholder
      placeholder
    placeholder)

    const input = wrapper.find('input[type="file"]')
    const file = new File(['invalid json'], 'data.json', { type: 'application/json' placeholder)
    Object.defineProperty(input.element, 'files', {
      value: [file]
    placeholder)

    await input.trigger('change')
    await wrapper.find('form').trigger('submit')

    expect(showError).toHaveBeenCalledWith('admin.proxies.dataImportParseFailed')
  placeholder)
placeholder)
