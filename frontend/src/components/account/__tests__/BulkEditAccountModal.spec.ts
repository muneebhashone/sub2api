import { describe, expect, it, vi placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'
import BulkEditAccountModal from '../BulkEditAccountModal.vue'

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn(),
    showInfo: vi.fn()
  placeholder)
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      bulkEdit: vi.fn()
    placeholder
  placeholder
placeholder))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    placeholder)
  placeholder
placeholder)

function mountModal() {
  return mount(BulkEditAccountModal, {
    props: {
      show: true,
      accountIds: [1, 2],
      selectedPlatforms: ['antigravity'],
      proxies: [],
      groups: []
    placeholder as any,
    global: {
      stubs: {
        BaseDialog: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
        Select: true,
        ProxySelector: true,
        GroupSelector: true,
        Icon: true
      placeholder
    placeholder
  placeholder)
placeholder

describe('BulkEditAccountModal', () => {
  it('antigravity 白名单包含 Gemini 图片模型且过滤掉普通 GPT 模型', () => {
    const wrapper = mountModal()

    expect(wrapper.text()).toContain('Gemini 3.1 Flash Image')
    expect(wrapper.text()).toContain('Gemini 3 Pro Image (Legacy)')
    expect(wrapper.text()).not.toContain('GPT-5.3 Codex')
  placeholder)

  it('antigravity 映射预设包含图片映射并过滤 OpenAI 预设', async () => {
    const wrapper = mountModal()

    const mappingTab = wrapper.findAll('button').find((btn) => btn.text().includes('admin.accounts.modelMapping'))
    expect(mappingTab).toBeTruthy()
    await mappingTab!.trigger('click')

    expect(wrapper.text()).toContain('Gemini 3.1 Image')
    expect(wrapper.text()).toContain('G3 Image→3.1')
    expect(wrapper.text()).not.toContain('GPT-5.3 Codex')
  placeholder)
placeholder)
