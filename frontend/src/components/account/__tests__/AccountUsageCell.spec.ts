import { describe, expect, it, vi, beforeEach placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'
import AccountUsageCell from '../AccountUsageCell.vue'

const { getUsage placeholder = vi.hoisted(() => ({
  getUsage: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      getUsage
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

describe('AccountUsageCell', () => {
  beforeEach(() => {
    getUsage.mockReset()
  placeholder)

  it('Antigravity 图片用量会聚合新旧 image 模型', async () => {
    getUsage.mockResolvedValue({
      antigravity_quota: {
        'gemini-3.1-flash-image': {
          utilization: 20,
          reset_time: '2026-03-01T10:00:00Z'
        placeholder,
        'gemini-3-pro-image': {
          utilization: 70,
          reset_time: '2026-03-01T09:00:00Z'
        placeholder
      placeholder
    placeholder)

    const wrapper = mount(AccountUsageCell, {
      props: {
        account: {
          id: 1001,
          platform: 'antigravity',
          type: 'oauth',
          extra: {placeholder
        placeholder as any
      placeholder,
      global: {
        stubs: {
          UsageProgressBar: {
            props: ['label', 'utilization', 'resetsAt', 'color'],
            template: '<div class="usage-bar">{{ label placeholderplaceholder|{{ utilization placeholderplaceholder|{{ resetsAt placeholderplaceholder</div>'
          placeholder,
          AccountQuotaInfo: true
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(wrapper.text()).toContain('admin.accounts.usageWindow.gemini3Image|70|2026-03-01T09:00:00Z')
  placeholder)
placeholder)
