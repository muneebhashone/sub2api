import { mount placeholder from '@vue/test-utils'
import { describe, expect, it, vi placeholder from 'vitest'

import type { MonitorQuotaSnapshot placeholder from '@/api/admin/channelMonitor'
import MonitorQuotaView from '../MonitorQuotaView.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    // te() 恒真：已知 token 直接返回 i18n key，便于断言 window/label 映射。
    useI18n: () => ({ t: (key: string) => key, te: () => true placeholder),
  placeholder
placeholder)

function makeSnapshot(overrides: Partial<MonitorQuotaSnapshot> = {placeholder): MonitorQuotaSnapshot {
  return {
    source: 'usage',
    success: true,
    fetched_at: '2026-08-18T00:00:00Z',
    ...overrides,
  placeholder
placeholder

describe('MonitorQuotaView', () => {
  it('renders nothing without a snapshot', () => {
    const wrapper = mount(MonitorQuotaView, { props: { snapshot: null placeholder placeholder)
    expect(wrapper.find('[data-testid="monitor-quota-view"]').exists()).toBe(false)
    expect(wrapper.text()).toBe('')
  placeholder)

  it('maps tier window/label tokens through i18n and colors utilization', () => {
    const wrapper = mount(MonitorQuotaView, {
      props: {
        snapshot: makeSnapshot({
          tiers: [
            { window: '5h', used_percent: 42.4 placeholder,
            { window: '7d', label: 'pro', used_percent: 80 placeholder,
            { window: 'weekly', label: 'unknown-token', used_percent: 95 placeholder,
          ],
        placeholder),
      placeholder,
    placeholder)

    const rows = wrapper.findAll('[data-testid="monitor-quota-view"] .flex.items-center')
    // tier 行只在 success 且有数据时渲染
    expect(rows.length).toBeGreaterThanOrEqual(3)
    const text = wrapper.text()
    // 已知 window token 走 i18n
    expect(text).toContain('monitorCommon.quota.windows.5h')
    expect(text).toContain('monitorCommon.quota.windows.7d')
    // 已知 label token 拼成 label/window
    expect(text).toContain('monitorCommon.quota.labels.pro/monitorCommon.quota.windows.7d')
    // 未知 label 原样透出（前向兼容）
    expect(text).toContain('unknown-token/monitorCommon.quota.windows.weekly')
    // 百分比取整
    expect(text).toContain('42%')
    expect(text).toContain('80%')
    expect(text).toContain('95%')

    const html = wrapper.html()
    // 阈值配色：≥90 红 / ≥75 黄 / 其余绿（与账号页 CNProviderQuotaCell 对齐）
    expect(html).toContain('bg-emerald-500')
    expect(html).toContain('bg-amber-500')
    expect(html).toContain('bg-red-500')
  placeholder)

  it('clamps the tier bar width into 0-100', () => {
    const wrapper = mount(MonitorQuotaView, {
      props: {
        snapshot: makeSnapshot({ tiers: [{ window: '5h', used_percent: 240 placeholder] placeholder),
      placeholder,
    placeholder)
    expect(wrapper.html()).toContain('width: 100%')
  placeholder)

  it('shows the plan level badge and multi-currency balances', () => {
    const wrapper = mount(MonitorQuotaView, {
      props: {
        snapshot: makeSnapshot({
          plan_level: 'Max20',
          balances: [
            { currency: 'CNY', balance: 12.5 placeholder,
            { currency: 'USD', balance: 0 placeholder,
          ],
        placeholder),
      placeholder,
    placeholder)

    expect(wrapper.text()).toContain('Max20')
    expect(wrapper.text()).toContain('12.50 CNY')
    expect(wrapper.text()).toContain('0.00 USD')
    // 余额为 0 用红色警示
    expect(wrapper.html()).toContain('text-red-600')
  placeholder)

  it('falls back to the single balance + currency pair', () => {
    const wrapper = mount(MonitorQuotaView, {
      props: { snapshot: makeSnapshot({ balance: 3.2, currency: 'CNY' placeholder) placeholder,
    placeholder)
    expect(wrapper.text()).toContain('3.20 CNY')
  placeholder)

  it('renders a truncated error state when the fetch failed', () => {
    const longError = 'x'.repeat(60)
    const wrapper = mount(MonitorQuotaView, {
      props: {
        snapshot: makeSnapshot({ success: false, error: longError placeholder),
      placeholder,
    placeholder)

    const error = wrapper.get('[data-testid="monitor-quota-error"]')
    expect(error.text()).toBe(`${'x'.repeat(48)placeholder…`)
    expect(error.attributes('title')).toBe(longError)
  placeholder)

  it('keeps failed snapshots from rendering tier rows', () => {
    const wrapper = mount(MonitorQuotaView, {
      props: {
        snapshot: makeSnapshot({
          success: false,
          tiers: [{ window: '5h', used_percent: 10 placeholder],
        placeholder),
      placeholder,
    placeholder)
    expect(wrapper.text()).not.toContain('10%')
  placeholder)
placeholder)
