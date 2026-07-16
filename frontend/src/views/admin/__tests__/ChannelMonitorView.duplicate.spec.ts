import { defineComponent placeholder from 'vue'
import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

import type { ChannelMonitor placeholder from '@/api/admin/channelMonitor'
import MonitorActionsCell from '@/components/admin/monitor/MonitorActionsCell.vue'
import ChannelMonitorView from '@/views/admin/ChannelMonitorView.vue'

const {
  listMonitors,
  duplicateMonitor,
  showSuccess,
  showError,
placeholder = vi.hoisted(() => ({
  listMonitors: vi.fn(),
  duplicateMonitor: vi.fn(),
  showSuccess: vi.fn(),
  showError: vi.fn(),
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    channelMonitor: {
      list: listMonitors,
      duplicate: duplicateMonitor,
      update: vi.fn(),
      runNow: vi.fn(),
      del: vi.fn(),
    placeholder,
  placeholder,
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showSuccess, showError placeholder),
placeholder))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key placeholder),
  placeholder
placeholder)

const AppLayoutStub = defineComponent({
  template: '<main><slot /></main>',
placeholder)

const TablePageLayoutStub = defineComponent({
  template: '<section><slot name="filters" /><slot name="table" /><slot name="pagination" /></section>',
placeholder)

const DataTableStub = defineComponent({
  props: {
    data: { type: Array, default: () => [] placeholder,
    columns: { type: Array, default: () => [] placeholder,
    loading: { type: Boolean, default: false placeholder,
  placeholder,
  template: '<div><div v-for="row in data" :key="row.id"><slot name="cell-actions" :row="row" /></div></div>',
placeholder)

function makeMonitor(overrides: Partial<ChannelMonitor> = {placeholder): ChannelMonitor {
  return {
    id: 42,
    name: 'primary',
    provider: 'openai',
    api_mode: 'chat_completions',
    endpoint: 'https://api.example.com',
    api_key_masked: 'sk-t***',
    primary_model: 'gpt-4o-mini',
    extra_models: [],
    group_name: '',
    enabled: true,
    interval_seconds: 60,
    jitter_seconds: 0,
    last_checked_at: null,
    created_by: 1,
    created_at: '2026-07-16T00:00:00Z',
    updated_at: '2026-07-16T00:00:00Z',
    primary_status: '',
    primary_latency_ms: null,
    availability_7d: 0,
    extra_models_status: [],
    template_id: null,
    extra_headers: {placeholder,
    body_override_mode: 'off',
    body_override: null,
    ...overrides,
  placeholder
placeholder

const monitor = makeMonitor()

function mountView() {
  return mount(ChannelMonitorView, {
    global: {
      stubs: {
        AppLayout: AppLayoutStub,
        TablePageLayout: TablePageLayoutStub,
        DataTable: DataTableStub,
        MonitorFiltersBar: true,
        Pagination: true,
        ConfirmDialog: true,
        EmptyState: true,
        HelpTooltip: true,
        Toggle: true,
        MonitorFormDialog: true,
        MonitorTemplateManagerDialog: true,
        MonitorRunResultDialog: true,
        MonitorPrimaryModelCell: true,
      placeholder,
    placeholder,
  placeholder)
placeholder

describe('ChannelMonitorView duplicate action', () => {
  beforeEach(() => {
    localStorage.clear()
    for (const fn of [listMonitors, duplicateMonitor, showSuccess, showError]) fn.mockReset()
    listMonitors.mockResolvedValue({
      items: [monitor],
      total: 1,
      page: 1,
      page_size: 20,
      pages: 1,
    placeholder)
    duplicateMonitor.mockResolvedValue(makeMonitor({ id: 43, name: 'primary (Copy)', enabled: false placeholder))
  placeholder)

  it('duplicates the selected monitor, reports success, and refreshes the list', async () => {
    const wrapper = mountView()
    await flushPromises()

    wrapper.findComponent(MonitorActionsCell).vm.$emit('duplicate', monitor)
    await flushPromises()

    expect(duplicateMonitor).toHaveBeenCalledTimes(1)
    expect(duplicateMonitor).toHaveBeenCalledWith(42)
    expect(showSuccess).toHaveBeenCalledWith('admin.channelMonitor.duplicateSuccess')
    expect(listMonitors.mock.calls.length).toBeGreaterThan(1)
    wrapper.unmount()
  placeholder)

  it('keeps a successful duplicate successful when the follow-up refresh fails', async () => {
    listMonitors
      .mockResolvedValueOnce({
        items: [monitor],
        total: 1,
        page: 1,
        page_size: 20,
        pages: 1,
      placeholder)
      .mockRejectedValueOnce(new Error('refresh failed'))
    const wrapper = mountView()
    await flushPromises()

    wrapper.findComponent(MonitorActionsCell).vm.$emit('duplicate', monitor)
    await flushPromises()

    expect(showSuccess).toHaveBeenCalledWith('admin.channelMonitor.duplicateSuccess')
    expect(showError).toHaveBeenCalledWith('refresh failed')
    expect(showError).not.toHaveBeenCalledWith('admin.channelMonitor.duplicateFailed')
    wrapper.unmount()
  placeholder)

  it('ignores repeated clicks while a duplicate request is in flight', async () => {
    let resolveDuplicate!: (value: ChannelMonitor) => void
    duplicateMonitor.mockImplementationOnce(() => new Promise(resolve => { resolveDuplicate = resolve placeholder))
    const wrapper = mountView()
    await flushPromises()

    const actions = wrapper.findComponent(MonitorActionsCell)
    actions.vm.$emit('duplicate', monitor)
    actions.vm.$emit('duplicate', monitor)
    await wrapper.vm.$nextTick()

    expect(duplicateMonitor).toHaveBeenCalledTimes(1)
    expect(actions.props('duplicating')).toBe(true)

    resolveDuplicate(makeMonitor({ id: 43, name: 'primary (Copy)', enabled: false placeholder))
    await flushPromises()
    expect(wrapper.findComponent(MonitorActionsCell).props('duplicating')).toBe(false)
    wrapper.unmount()
  placeholder)

  it('shows the API error when duplication fails', async () => {
    duplicateMonitor.mockRejectedValueOnce(new Error('duplicate failed'))
    const wrapper = mountView()
    await flushPromises()

    wrapper.findComponent(MonitorActionsCell).vm.$emit('duplicate', monitor)
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('duplicate failed')
    wrapper.unmount()
  placeholder)

  it('rejects a defensive duplicate event when the API key is unavailable', async () => {
    const unavailable = makeMonitor({ id: 99, api_key_decrypt_failed: true placeholder)
    listMonitors.mockResolvedValueOnce({
      items: [unavailable],
      total: 1,
      page: 1,
      page_size: 20,
      pages: 1,
    placeholder)
    const wrapper = mountView()
    await flushPromises()

    wrapper.findComponent(MonitorActionsCell).vm.$emit('duplicate', unavailable)
    await flushPromises()

    expect(duplicateMonitor).not.toHaveBeenCalled()
    expect(showError).toHaveBeenCalledWith('admin.channelMonitor.duplicateKeyUnavailable')
    wrapper.unmount()
  placeholder)
placeholder)
