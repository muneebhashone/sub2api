import { mount placeholder from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi placeholder from 'vitest'
import WechatOAuthSection from '@/components/auth/WechatOAuthSection.vue'

const routeState = vi.hoisted(() => ({
  query: {placeholder as Record<string, unknown>,
placeholder))

const locationState = vi.hoisted(() => ({
  current: { href: 'http://localhost/login' placeholder as { href: string placeholder,
placeholder))

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, string>) => {
      if (key === 'auth.oidc.signIn') {
        return `Continue with ${params?.providerName ?? ''placeholder`.trim()
      placeholder
      if (key === 'auth.oauthOrContinue') {
        return 'or continue'
      placeholder
      return key
    placeholder,
  placeholder),
placeholder))

describe('WechatOAuthSection', () => {
  beforeEach(() => {
    routeState.query = { redirect: '/billing?plan=pro' placeholder
    locationState.current = { href: 'http://localhost/login' placeholder
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: locationState.current,
    placeholder)
    Object.defineProperty(window.navigator, 'userAgent', {
      configurable: true,
      value: 'Mozilla/5.0',
    placeholder)
  placeholder)

  afterEach(() => {
    vi.unstubAllGlobals()
  placeholder)

  it('starts the open WeChat OAuth flow with the current redirect target', async () => {
    const wrapper = mount(WechatOAuthSection)

    expect(wrapper.text()).toContain('WeChat')

    await wrapper.get('button').trigger('click')

    expect(locationState.current.href).toContain(
      '/api/v1/auth/oauth/wechat/start?mode=open&redirect=%2Fbilling%3Fplan%3Dpro'
    )
  placeholder)

  it('uses mp mode inside the WeChat browser', async () => {
    Object.defineProperty(window.navigator, 'userAgent', {
      configurable: true,
      value: 'Mozilla/5.0 MicroMessenger',
    placeholder)
    const wrapper = mount(WechatOAuthSection)

    await wrapper.get('button').trigger('click')

    expect(locationState.current.href).toContain(
      '/api/v1/auth/oauth/wechat/start?mode=mp&redirect=%2Fbilling%3Fplan%3Dpro'
    )
  placeholder)
placeholder)
