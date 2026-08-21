import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import CNProviderBalanceCell from '../CNProviderBalanceCell.vue'
import type { Account placeholder from '@/types'

const { queryBalance placeholder = vi.hoisted(() => ({
  queryBalance: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    cnProviders: { queryBalance placeholder
  placeholder
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  placeholder)
placeholder))

const account = {
  id: 7,
  platform: 'kimi',
  type: 'apikey',
  credentials: { account_mode: 'payg' placeholder,
  extra: {
    kimi_balance: 12.5,
    kimi_balance_currency: 'CNY'
  placeholder
placeholder as Account

describe('CNProviderBalanceCell', () => {
  beforeEach(() => {
    queryBalance.mockReset()
  placeholder)

  it('renders the persisted balance as static text with an explicit query action', async () => {
    const wrapper = mount(CNProviderBalanceCell, { props: { account placeholder placeholder)
    await flushPromises()

    // Snapshot value renders without any probe.
    expect(queryBalance).not.toHaveBeenCalled()
    expect(wrapper.get('[data-test="cn-provider-balance-value"]').text()).toContain('CNY 12.50')

    // The control reads as an action; the i18n mock returns the key itself.
    const probeButton = wrapper.get('[data-test="cn-provider-balance-probe"]')
    expect(probeButton.text()).toBe('admin.accounts.cnProviders.probe')

    await probeButton.trigger('click')
    await flushPromises()
    expect(queryBalance).toHaveBeenCalledWith(account.id)
  placeholder)

  it('shows the low-balance badge from the snapshot marker', () => {
    const lowAccount = {
      ...account,
      extra: { kimi_balance: 0.4, kimi_balance_low: true placeholder
    placeholder as Account

    const wrapper = mount(CNProviderBalanceCell, { props: { account: lowAccount placeholder placeholder)

    expect(wrapper.text()).toContain('admin.accounts.cnProviders.balanceLow')
  placeholder)

  it('keeps the snapshot balance visible when a query fails', async () => {
    queryBalance.mockResolvedValue({ success: false, error: 'HTTP 401' placeholder)
    const wrapper = mount(CNProviderBalanceCell, { props: { account placeholder placeholder)

    await wrapper.get('[data-test="cn-provider-balance-probe"]').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('CNY 12.50')
    expect(wrapper.text()).toContain('HTTP 401')
  placeholder)
placeholder)
