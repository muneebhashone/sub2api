import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import UserEditModal from '../UserEditModal.vue'

const { update, updateUserAttributeValues, showSuccess, showError placeholder = vi.hoisted(() => ({
  update: vi.fn(),
  updateUserAttributeValues: vi.fn(),
  showSuccess: vi.fn(),
  showError: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    users: { update placeholder,
    userAttributes: { updateUserAttributeValues placeholder
  placeholder
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ showSuccess, showError placeholder)
placeholder))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({ copyToClipboard: vi.fn() placeholder)
placeholder))

// useStepUp pulls in the API client, which needs the real i18n instance.
vi.mock('vue-i18n', async (importOriginal) => ({
  ...(await importOriginal<typeof import('vue-i18n')>()),
  useI18n: () => ({
    t: (key: string, params?: Record<string, unknown>) =>
      params ? `${keyplaceholder:${JSON.stringify(params)placeholder` : key
  placeholder)
placeholder))

const mountModal = (concurrency: number) => mount(UserEditModal, {
  props: {
    show: true,
    user: { id: 7, email: 'user@example.test', username: 'user', notes: '', role: 'user', concurrency, rpm_limit: 0 placeholder as never
  placeholder,
  global: {
    stubs: {
      BaseDialog: {
        props: ['show', 'title'],
        template: '<div v-if="show"><slot /><slot name="footer" /></div>'
      placeholder,
      Select: true,
      Icon: true,
      UserAttributeForm: true,
      TotpStepUpDialog: true
    placeholder
  placeholder
placeholder)

describe('UserEditModal concurrency', () => {
  beforeEach(() => {
    update.mockReset()
    updateUserAttributeValues.mockReset()
    showSuccess.mockReset()
    showError.mockReset()
    update.mockResolvedValue({placeholder)
  placeholder)

  // Regression coverage for issue #5977: the gateway treats concurrency <= 0 as
  // unlimited (AcquireUserSlot) and both the batch limits endpoint and the bulk
  // edit modal accept 0, so this dialog must not be the only place that rejects
  // it — doing so blocked every other edit on such a user.
  it('saves an unlimited (0) concurrency instead of blocking the whole form', async () => {
    const wrapper = mountModal(0)

    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(showError).not.toHaveBeenCalled()
    expect(update).toHaveBeenCalledWith(7, expect.objectContaining({ concurrency: 0 placeholder))
    expect(wrapper.emitted('success')).toBeTruthy()
  placeholder)

  it('still rejects a negative concurrency', async () => {
    const wrapper = mountModal(3)

    await wrapper.get('[data-test="concurrency-input"]').setValue('-1')
    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(showError).toHaveBeenCalledWith('admin.users.concurrencyNonNegative')
    expect(update).not.toHaveBeenCalled()
  placeholder)
placeholder)
