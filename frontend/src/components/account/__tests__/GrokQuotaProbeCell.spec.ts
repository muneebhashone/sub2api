import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import GrokQuotaProbeCell from '../GrokQuotaProbeCell.vue'
import type { Account placeholder from '@/types'

const { queryQuota placeholder = vi.hoisted(() => ({
  queryQuota: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    grok: { queryQuota placeholder
  placeholder
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) =>
      params?.percent == null ? key : `${keyplaceholder:${params.percentplaceholder`
  placeholder)
placeholder))

const account = {
  id: 99,
  platform: 'grok',
  type: 'oauth'
placeholder as Account

describe('GrokQuotaProbeCell', () => {
  beforeEach(() => {
    queryQuota.mockReset()
  placeholder)

  it('keeps billing data while exposing a failed Free quota fallback', async () => {
    queryQuota.mockResolvedValue({
      source: 'hybrid_probe',
      billing: { period_type: 'weekly', usage_percent: null placeholder,
      headers_observed: false,
      reset_supported: false,
      fetched_at: 1,
      probe_error: 'upstream returned 402 for probe model "grok-4.5"'
    placeholder)
    const wrapper = mount(GrokQuotaProbeCell, { props: { account placeholder placeholder)

    await wrapper.get('button').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('upstream returned 402 for probe model "grok-4.5"')
    expect(wrapper.emitted('probed')?.[0]?.[0]).toMatchObject({
      billing: { period_type: 'weekly', usage_percent: null placeholder,
      probe_error: 'upstream returned 402 for probe model "grok-4.5"'
    placeholder)
  placeholder)
placeholder)
