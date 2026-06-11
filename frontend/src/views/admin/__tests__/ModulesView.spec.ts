import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'
import type { DOMWrapper, VueWrapper placeholder from '@vue/test-utils'
import { createPinia placeholder from 'pinia'

import ModulesView from '../ModulesView.vue'
import type { Module placeholder from '@/types'

const { listModules, showError, showSuccess placeholder = vi.hoisted(() => ({
  listModules: vi.fn(),
  showError: vi.fn(),
  showSuccess: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    modules: {
      list: listModules
    placeholder
  placeholder
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
    showSuccess
  placeholder)
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

const createModules = (): Module[] => [
  {
    id: 'core.gateway',
    namespace: 'core',
    name: 'gateway',
    enabled: true,
    state: 'running',
    error: ''
  placeholder,
  {
    id: 'demo.broken',
    namespace: 'demo',
    name: 'broken',
    enabled: true,
    state: 'errored',
    error: 'provision failed: dependency missing'
  placeholder
]

const AppLayoutStub = { template: '<div><slot /></div>' placeholder
const TablePageLayoutStub = {
  template: '<div><slot name="filters" /><slot name="table" /><slot name="pagination" /></div>'
placeholder
const DataTableStub = {
  props: ['columns', 'data', 'loading'],
  template: `
    <div>
      <div data-test="columns">{{ columns.map(col => col.key).join(',') placeholderplaceholder</div>
      <div v-if="loading" data-test="loading">loading</div>
      <div v-else-if="!data || data.length === 0" data-test="empty"><slot name="empty" /></div>
      <template v-else>
        <div v-for="row in data" :key="row.id" data-test="row">
          <slot name="cell-id" :value="row.id" :row="row" />
          <slot name="cell-namespace" :value="row.namespace" :row="row" />
          <slot name="cell-enabled" :value="row.enabled" :row="row" />
          <slot name="cell-state" :value="row.state" :row="row" />
          <slot name="cell-error" :value="row.error" :row="row" />
        </div>
      </template>
    </div>
  `
placeholder

function mountView() {
  return mount(ModulesView, {
    global: {
      plugins: [createPinia()],
      stubs: {
        AppLayout: AppLayoutStub,
        TablePageLayout: TablePageLayoutStub,
        DataTable: DataTableStub,
        Icon: true
      placeholder
    placeholder
  placeholder)
placeholder

function findRefreshButton(wrapper: VueWrapper): DOMWrapper<HTMLButtonElement> {
  const button = wrapper
    .findAll<HTMLButtonElement>('button')
    .find((item) => item.attributes('title') === 'common.refresh')
  if (!button) {
    throw new Error('refresh button not found')
  placeholder
  return button
placeholder

describe('admin ModulesView', () => {
  beforeEach(() => {
    listModules.mockReset()
    showError.mockReset()
    showSuccess.mockReset()

    listModules.mockResolvedValue({ modules: createModules() placeholder)
  placeholder)

  it('fetches modules on mount and renders the module list', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(listModules).toHaveBeenCalledTimes(1)

    const columns = wrapper.get('[data-test="columns"]').text()
    expect(columns.split(',')).toEqual(['id', 'namespace', 'name', 'enabled', 'state', 'error'])

    const rows = wrapper.findAll('[data-test="row"]')
    expect(rows).toHaveLength(2)
    expect(wrapper.text()).toContain('core.gateway')
    expect(wrapper.text()).toContain('demo.broken')
    expect(showError).not.toHaveBeenCalled()
  placeholder)

  it('renders semantic state badges and error text', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.find('.badge-success').text()).toBe('admin.modules.stateLabels.running')
    expect(wrapper.find('.badge-danger').text()).toBe('admin.modules.stateLabels.errored')
    expect(wrapper.text()).toContain('provision failed: dependency missing')
  placeholder)

  it('shows the loading state while the request is in flight', async () => {
    let resolveList!: (value: { modules: Module[] placeholder) => void
    listModules.mockImplementation(
      () =>
        new Promise<{ modules: Module[] placeholder>((resolve) => {
          resolveList = resolve
        placeholder)
    )

    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.find('[data-test="loading"]').exists()).toBe(true)

    resolveList({ modules: [] placeholder)
    await flushPromises()

    expect(wrapper.find('[data-test="loading"]').exists()).toBe(false)
  placeholder)

  it('renders the empty state when no modules are registered', async () => {
    listModules.mockResolvedValue({ modules: [] placeholder)

    const wrapper = mountView()
    await flushPromises()

    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('admin.modules.noModules')
    expect(wrapper.text()).toContain('admin.modules.noModulesDescription')
  placeholder)

  it('reloads modules when the refresh button is clicked', async () => {
    const wrapper = mountView()
    await flushPromises()

    expect(listModules).toHaveBeenCalledTimes(1)

    await findRefreshButton(wrapper).trigger('click')
    await flushPromises()

    expect(listModules).toHaveBeenCalledTimes(2)
  placeholder)

  it('shows an error toast when loading fails', async () => {
    listModules.mockRejectedValue(new Error('network down'))

    const wrapper = mountView()
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('admin.modules.failedToLoad')
    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  placeholder)
placeholder)
