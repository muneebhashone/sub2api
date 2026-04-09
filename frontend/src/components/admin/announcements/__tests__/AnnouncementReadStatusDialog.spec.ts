import { describe, it, expect, vi, beforeEach placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import AnnouncementReadStatusDialog from '../AnnouncementReadStatusDialog.vue'

const { getReadStatus, showError placeholder = vi.hoisted(() => ({
  getReadStatus: vi.fn(),
  showError: vi.fn(),
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    announcements: {
      getReadStatus,
    placeholder,
  placeholder,
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError,
  placeholder),
placeholder))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    placeholder),
  placeholder
placeholder)

vi.mock('@/composables/usePersistedPageSize', () => ({
  getPersistedPageSize: () => 20,
placeholder))

const BaseDialogStub = {
  props: ['show', 'title', 'width'],
  emits: ['close'],
  template: '<div><slot /><slot name="footer" /></div>',
placeholder

describe('AnnouncementReadStatusDialog', () => {
  beforeEach(() => {
    getReadStatus.mockReset()
    showError.mockReset()
    vi.useFakeTimers()
  placeholder)

  it('closes by aborting active requests and clearing debounced reloads', async () => {
    let activeSignal: AbortSignal | undefined
    getReadStatus.mockImplementation(async (...args: any[]) => {
      activeSignal = args[4]?.signal
      return new Promise(() => {placeholder)
    placeholder)

    const wrapper = mount(AnnouncementReadStatusDialog, {
      props: {
        show: false,
        announcementId: 1,
      placeholder,
      global: {
        stubs: {
          BaseDialog: BaseDialogStub,
          DataTable: true,
          Pagination: true,
          Icon: true,
        placeholder,
      placeholder,
    placeholder)

    await wrapper.setProps({ show: true placeholder)
    await flushPromises()

    expect(getReadStatus).toHaveBeenCalledTimes(1)
    expect(activeSignal?.aborted).toBe(false)

    const setupState = (wrapper.vm as any).$?.setupState
    setupState.search = 'alice'
    setupState.handleSearch()

    setupState.handleClose()
    await flushPromises()

    expect(activeSignal?.aborted).toBe(true)
    expect(wrapper.emitted('close')).toHaveLength(1)

    vi.advanceTimersByTime(350)
    await flushPromises()

    expect(getReadStatus).toHaveBeenCalledTimes(1)
  placeholder)
placeholder)
