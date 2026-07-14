import { mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

import DataTable from '../DataTable.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  placeholder)
placeholder))

const stubDesktopMatchMedia = () => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation((query: string) => ({
      matches: true,
      media: query,
      onchange: null,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      addListener: vi.fn(),
      removeListener: vi.fn(),
      dispatchEvent: vi.fn()
    placeholder))
  placeholder)
placeholder

describe('DataTable', () => {
  beforeEach(() => {
    stubDesktopMatchMedia()
    localStorage.clear()
  placeholder)

  it('renders paired sort arrows and highlights the active direction', async () => {
    const wrapper = mount(DataTable, {
      props: {
        columns: [
          { key: 'name', label: 'Name', sortable: true placeholder,
          { key: 'created_at', label: 'Created', sortable: true placeholder
        ],
        data: [
          { id: 1, name: 'Beta', created_at: '2026-01-02T00:00:00Z' placeholder,
          { id: 2, name: 'Alpha', created_at: '2026-01-01T00:00:00Z' placeholder
        ],
        defaultSortKey: 'name',
        defaultSortOrder: 'asc'
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const nameHeader = wrapper.findAll('th')[0]
    expect(nameHeader.attributes('aria-sort')).toBe('ascending')
    expect(nameHeader.findAll('svg')).toHaveLength(2)
    expect(nameHeader.findAll('svg')[0].classes()).toContain('text-primary-600')
    expect(nameHeader.findAll('svg')[1].classes()).toContain('text-gray-300')

    await nameHeader.trigger('click')
    await wrapper.vm.$nextTick()

    expect(nameHeader.attributes('aria-sort')).toBe('descending')
    expect(nameHeader.findAll('svg')[0].classes()).toContain('text-gray-300')
    expect(nameHeader.findAll('svg')[1].classes()).toContain('text-primary-600')
  placeholder)

  it('renders every row with no virtual padding spacer for small datasets (virtualization off)', async () => {
    const data = Array.from({ length: 8 placeholder, (_, i) => ({ id: i + 1, name: `Row ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    // Virtualization is OFF for a small list…
    expect((wrapper.vm as any).shouldVirtualize).toBe(false)
    // …every row is in the DOM…
    expect(wrapper.findAll('tbody tr[data-index]')).toHaveLength(data.length)
    // …and there are no aria-hidden virtual padding spacer rows.
    expect(wrapper.findAll('tbody tr[aria-hidden="true"]')).toHaveLength(0)
  placeholder)

  it('switches to windowed rendering once row count exceeds virtualizeThreshold', async () => {
    const data = Array.from({ length: 12 placeholder, (_, i) => ({ id: i + 1, name: `Row ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data,
        virtualizeThreshold: 3
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    // Virtualization is ON: the mode-switch decision flipped…
    expect((wrapper.vm as any).shouldVirtualize).toBe(true)
    // …and the virtualizer drives off the full row count.
    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    expect(instance.options.count).toBe(data.length)
  placeholder)

  it('keys the virtualizer size cache by row identity, not index (avoids stale heights on sort/filter)', async () => {
    const data = Array.from({ length: 12 placeholder, (_, i) => ({ id: 100 + i, name: `Row ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data,
        rowKey: 'id',
        virtualizeThreshold: 3
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    // getItemKey must resolve to the row's stable key (id), not the positional index.
    expect(instance.options.getItemKey(0)).toBe(100)
    expect(instance.options.getItemKey(5)).toBe(105)
  placeholder)

  it('clears stale row and element caches when pagination replaces the row ID set', async () => {
    const firstPage = Array.from({ length: 100 placeholder, (_, i) => ({ id: i + 1, name: `First ${i + 1placeholder` placeholder))
    const secondPage = Array.from({ length: 100 placeholder, (_, i) => ({ id: i + 101, name: `Second ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data: firstPage,
        rowKey: 'id',
        virtualizeThreshold: 1
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    const firstPageIDs = firstPage.map(row => row.id)
    ;(instance as any).itemSizeCache = new Map(firstPageIDs.map(id => [id, 156]))
    instance.elementsCache.clear()
    for (const id of firstPageIDs) {
      instance.elementsCache.set(id, document.createElement('tr'))
    placeholder
    const measureElementSpy = vi.spyOn(instance, 'measureElement')

    await wrapper.setProps({ data: secondPage placeholder)
    await wrapper.vm.$nextTick()

    const sizeCache = (instance as any).itemSizeCache as Map<number, number>
    expect(sizeCache.size).toBeLessThanOrEqual(secondPage.length)
    expect(instance.elementsCache.size).toBeLessThanOrEqual(secondPage.length)
    expect(firstPageIDs.some(id => sizeCache.has(id))).toBe(false)
    expect(firstPageIDs.some(id => instance.elementsCache.has(id))).toBe(false)
    expect(measureElementSpy.mock.calls.some(([node]) => node === null)).toBe(true)
  placeholder)

  it('clears stale caches when equal-length pages replace rows without stable keys', async () => {
    const firstPage = Array.from({ length: 12 placeholder, (_, i) => ({ name: `First ${i + 1placeholder` placeholder))
    const secondPage = Array.from({ length: 12 placeholder, (_, i) => ({ name: `Second ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data: firstPage,
        virtualizeThreshold: 1
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    const measureElementSpy = vi.spyOn(instance, 'measureElement')

    await wrapper.setProps({ data: secondPage placeholder)
    await wrapper.vm.$nextTick()

    expect(measureElementSpy.mock.calls.some(([node]) => node === null)).toBe(true)
  placeholder)

  it('conservatively clears caches when duplicate row-key multiplicity changes', async () => {
    const firstPage = [
      { id: 1, name: 'First A' placeholder,
      { id: 1, name: 'First B' placeholder,
      { id: 2, name: 'First C' placeholder
    ]
    const secondPage = [
      { id: 1, name: 'Second A' placeholder,
      { id: 2, name: 'Second B' placeholder,
      { id: 2, name: 'Second C' placeholder
    ]
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data: firstPage,
        rowKey: 'id',
        virtualizeThreshold: 1
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    const measureElementSpy = vi.spyOn(instance, 'measureElement')

    await wrapper.setProps({ data: secondPage placeholder)
    await wrapper.vm.$nextTick()

    expect(measureElementSpy.mock.calls.some(([node]) => node === null)).toBe(true)
  placeholder)

  it('preserves cache when rows without stable keys only reorder the same objects', async () => {
    const data = Array.from({ length: 12 placeholder, (_, i) => ({ name: `Row ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data,
        virtualizeThreshold: 1
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    const measureSpy = vi.spyOn(instance, 'measure')

    await wrapper.setProps({ data: [...data].reverse() placeholder)
    await wrapper.vm.$nextTick()

    expect(measureSpy).not.toHaveBeenCalled()
  placeholder)

  it('preserves stable row height cache when the same row IDs are only reordered', async () => {
    const data = Array.from({ length: 100 placeholder, (_, i) => ({ id: i + 1, name: `Row ${i + 1placeholder` placeholder))
    const wrapper = mount(DataTable, {
      props: {
        columns: [{ key: 'name', label: 'Name' placeholder],
        data,
        rowKey: 'id',
        virtualizeThreshold: 1
      placeholder
    placeholder)

    await wrapper.vm.$nextTick()

    const exposed = (wrapper.vm as any).virtualizer
    const instance = exposed?.value ?? exposed
    ;(instance as any).itemSizeCache = new Map(data.map(row => [row.id, 156]))
    const measureSpy = vi.spyOn(instance, 'measure')

    await wrapper.setProps({ data: [...data].reverse() placeholder)
    await wrapper.vm.$nextTick()

    const sizeCache = (instance as any).itemSizeCache as Map<number, number>
    expect(measureSpy).not.toHaveBeenCalled()
    expect(sizeCache.size).toBe(100)
  placeholder)
placeholder)
