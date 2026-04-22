import { mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import OAuthCallbackView from '@/views/auth/OAuthCallbackView.vue'

const { routeState, showErrorMock, copyToClipboardMock placeholder = vi.hoisted(() => ({
  routeState: {
    query: {placeholder as Record<string, unknown>,
  placeholder,
  showErrorMock: vi.fn(),
  copyToClipboardMock: vi.fn(),
placeholder))

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
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

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: (...args: any[]) => copyToClipboardMock(...args),
  placeholder),
placeholder))

describe('OAuthCallbackView', () => {
  beforeEach(() => {
    routeState.query = {placeholder
    showErrorMock.mockReset()
    copyToClipboardMock.mockReset()
  placeholder)

  it('renders localized callback copy actions', () => {
    routeState.query = {
      code: 'oauth-code',
      state: 'oauth-state',
    placeholder

    const wrapper = mount(OAuthCallbackView)

    expect(wrapper.text()).toContain('auth.oauth.callbackTitle')
    expect(wrapper.text()).toContain('auth.oauth.callbackHint')
    expect(wrapper.text()).toContain('common.copy')
    expect(wrapper.find('input[value="oauth-code"]').exists()).toBe(true)
    expect(wrapper.find('input[value="oauth-state"]').exists()).toBe(true)
  placeholder)

  it('sends callback errors to toast instead of rendering inline red text', () => {
    routeState.query = {
      error: 'oauth failed',
    placeholder

    const wrapper = mount(OAuthCallbackView)

    expect(showErrorMock).toHaveBeenCalledWith('oauth failed')
    expect(wrapper.text()).not.toContain('oauth failed')
    expect(wrapper.find('.bg-red-50').exists()).toBe(false)
  placeholder)
placeholder)
