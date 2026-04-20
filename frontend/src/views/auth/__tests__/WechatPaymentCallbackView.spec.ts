import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import WechatPaymentCallbackView from '@/views/auth/WechatPaymentCallbackView.vue'

const { replaceMock, routeState, locationState placeholder = vi.hoisted(() => ({
  replaceMock: vi.fn(),
  routeState: {
    query: {placeholder as Record<string, unknown>,
  placeholder,
  locationState: {
    current: {
      href: 'http://localhost/auth/wechat/payment/callback',
      hash: '',
      search: '',
      pathname: '/auth/wechat/payment/callback',
      origin: 'http://localhost',
    placeholder as Location & { origin: string placeholder,
  placeholder,
placeholder))

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
  useRouter: () => ({
    replace: replaceMock,
  placeholder),
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
    locale: { value: 'zh-CN' placeholder,
  placeholder),
placeholder))

describe('WechatPaymentCallbackView', () => {
  beforeEach(() => {
    replaceMock.mockReset()
    routeState.query = {placeholder
    locationState.current = {
      href: 'http://localhost/auth/wechat/payment/callback',
      hash: '',
      search: '',
      pathname: '/auth/wechat/payment/callback',
      origin: 'http://localhost',
    placeholder as Location & { origin: string placeholder
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: locationState.current,
    placeholder)
  placeholder)

  it('redirects back to purchase with an opaque resume token from hash fragment', async () => {
    locationState.current.hash = '#wechat_resume_token=resume-token-123&redirect=%2Fpurchase%3Ffrom%3Dwechat'

    mount(WechatPaymentCallbackView)
    await flushPromises()

    expect(replaceMock).toHaveBeenCalledWith({
      path: '/purchase',
      query: {
        from: 'wechat',
        wechat_resume: '1',
        wechat_resume_token: 'resume-token-123',
      placeholder,
    placeholder)
  placeholder)

  it('shows an error when the callback payload is missing the resume token', async () => {
    locationState.current.hash = '#payment_type=wxpay'

    const wrapper = mount(WechatPaymentCallbackView)
    await flushPromises()

    expect(replaceMock).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('微信支付回调缺少恢复令牌。')
  placeholder)
placeholder)
