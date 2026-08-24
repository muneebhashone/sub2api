import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import AccountsView from '../AccountsView.vue'

const { listAccounts placeholder = vi.hoisted(() => ({
  listAccounts: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      list: listAccounts,
      listWithEtag: vi.fn(),
      getBatchTodayStats: vi.fn().mockResolvedValue({ stats: {placeholder placeholder),
      getUpstreamBillingProbeSettings: vi.fn().mockResolvedValue({ enabled: true, interval_minutes: 30 placeholder),
      delete: vi.fn(),
      batchClearError: vi.fn(),
      batchRefresh: vi.fn(),
      toggleSchedulable: vi.fn()
    placeholder,
    proxies: { getAll: vi.fn().mockResolvedValue([]) placeholder,
    groups: { getAll: vi.fn().mockResolvedValue([]) placeholder
  placeholder
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showError: vi.fn(), showSuccess: vi.fn(), showInfo: vi.fn() placeholder)
placeholder))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({ token: 'test-token', isSimpleMode: false placeholder)
placeholder))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({ t: (key: string) => key placeholder)
  placeholder
placeholder)

const DataTableStub = {
  props: ['columns'],
  emits: ['sort'],
  template: `
    <div data-test="data-table">
      <span v-for="column in columns" :key="column.key" :data-column="column.key">
        {{ column.sortable ? 'sortable' : 'fixed' placeholderplaceholder
      </span>
      <button data-test="sort-priority" @click="$emit('sort', 'priority', 'desc')" />
    </div>
  `
placeholder

function mountView() {
  return mount(AccountsView, {
    global: {
      stubs: {
        AppLayout: { template: '<div><slot /></div>' placeholder,
        TablePageLayout: {
          template: '<div><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>'
        placeholder,
        DataTable: DataTableStub,
        AccountTableActions: { template: '<div><slot name="after" /></div>' placeholder,
        AccountTableFilters: true,
        AccountBulkActionsBar: true,
        Pagination: true,
        ConfirmDialog: true,
        AccountActionMenu: true,
        ImportDataModal: true,
        ReAuthAccountModal: true,
        AccountTestModal: true,
        AccountStatsModal: true,
        ScheduledTestsPanel: true,
        SyncFromCrsModal: true,
        TempUnschedStatusModal: true,
        ErrorPassthroughRulesModal: true,
        TLSFingerprintProfilesModal: true,
        CreateAccountModal: true,
        EditAccountModal: true,
        BulkEditAccountModal: true,
        PlatformTypeBadge: true,
        AccountCapacityCell: true,
        AccountStatusIndicator: true,
        AccountTodayStatsCell: true,
        AccountGroupsCell: true,
        AccountUsageCell: true,
        HelpTooltip: true,
        Icon: true,
        Teleport: true
      placeholder
    placeholder
  placeholder)
placeholder

describe('admin AccountsView priority column preferences', () => {
  beforeEach(() => {
    localStorage.clear()
    listAccounts.mockReset().mockResolvedValue({
      items: [],
      total: 0,
      page: 1,
      page_size: 20,
      pages: 0
    placeholder)
  placeholder)

  it('shows priority as a sortable column for fresh preferences', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.get('[data-column="priority"]').text()).toBe('sortable')

    await wrapper.get('[data-test="sort-priority"]').trigger('click')
    await flushPromises()

    expect(listAccounts).toHaveBeenLastCalledWith(
      1,
      20,
      expect.objectContaining({ sort_by: 'priority', sort_order: 'desc' placeholder),
      expect.objectContaining({ signal: expect.any(AbortSignal) placeholder)
    )
  placeholder)

  it('preserves an existing preference that explicitly hides priority', async () => {
    localStorage.setItem('account-hidden-columns', JSON.stringify(['priority', 'today_stats']))
    localStorage.setItem('account-hidden-columns-version', 'scheduler-score-hidden-by-default')

    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.find('[data-column="priority"]').exists()).toBe(false)
    expect(JSON.parse(localStorage.getItem('account-hidden-columns') || '[]')).toEqual([
      'priority',
      'today_stats'
    ])
  placeholder)

  it('keeps priority visible while migrating older saved preferences', async () => {
    localStorage.setItem('account-hidden-columns', JSON.stringify(['today_stats']))

    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.get('[data-column="priority"]').text()).toBe('sortable')
    expect(JSON.parse(localStorage.getItem('account-hidden-columns') || '[]')).toEqual(
      expect.arrayContaining(['today_stats', 'scheduler_score'])
    )
    expect(JSON.parse(localStorage.getItem('account-hidden-columns') || '[]')).not.toContain('priority')
  placeholder)
placeholder)
