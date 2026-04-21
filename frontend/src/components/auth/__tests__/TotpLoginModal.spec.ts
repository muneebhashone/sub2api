import { mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import TotpLoginModal from '@/components/auth/TotpLoginModal.vue'

const { showErrorMock placeholder = vi.hoisted(() => ({
  showErrorMock: vi.fn(),
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  placeholder),
placeholder))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    showError: (...args: any[]) => showErrorMock(...args),
  placeholder),
placeholder))

describe('TotpLoginModal', () => {
  beforeEach(() => {
    showErrorMock.mockReset()
  placeholder)

  it('sends verification errors to toast and does not render inline red text', async () => {
    const wrapper = mount(TotpLoginModal, {
      props: {
        tempToken: 'temp-token',
        userEmailMasked: 'u***@example.com',
      placeholder,
    placeholder)

    ;(wrapper.vm as unknown as { setError: (message: string) => void placeholder).setError('Invalid code')
    await wrapper.vm.$nextTick()

    expect(showErrorMock).toHaveBeenCalledWith('Invalid code')
    expect(wrapper.text()).not.toContain('Invalid code')
    expect(wrapper.find('.bg-red-50').exists()).toBe(false)
  placeholder)
placeholder)
