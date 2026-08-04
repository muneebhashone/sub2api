import { defineComponent, h placeholder from 'vue'
import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import LoginView from '@/views/auth/LoginView.vue'

const loginMock = vi.fn()
const loginWithPasskeyMock = vi.fn()
const getPublicSettingsMock = vi.fn()
const startOAuthLoginMock = vi.fn()
const verifyTencentMock = vi.fn()
const captchaResetMock = vi.fn()
const locationState = { href: 'http://localhost/login' placeholder

vi.mock('vue-router', () => ({
  useRouter: () => ({
    currentRoute: { value: { query: {placeholder placeholder placeholder,
    push: vi.fn()
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

vi.mock('@/stores', () => ({
  useAuthStore: () => ({
    login: (...args: unknown[]) => loginMock(...args),
    loginWithPasskey: (...args: unknown[]) => loginWithPasskeyMock(...args)
  placeholder),
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn(),
    showWarning: vi.fn()
  placeholder)
placeholder))

vi.mock('@/api/auth', async () => {
  const actual = await vi.importActual<typeof import('@/api/auth')>('@/api/auth')
  return {
    ...actual,
    getPublicSettings: (...args: unknown[]) => getPublicSettingsMock(...args),
    startOAuthLogin: (...args: unknown[]) => startOAuthLoginMock(...args),
    isTotp2FARequired: () => false,
    isWeChatWebOAuthEnabled: () => false
  placeholder
placeholder)

const CaptchaChallengeStub = defineComponent({
  setup(_, { expose placeholder) {
    expose({
      verifyTencent: verifyTencentMock,
      reset: captchaResetMock
    placeholder)
    return () => h('div')
  placeholder
placeholder)

const OAuthButtonStub = defineComponent({
  emits: ['start'],
  setup(_, { emit placeholder) {
    return () => h('button', {
      type: 'button',
      'data-testid': 'oauth-start',
      onClick: () => emit('start', {
        provider: 'github',
        params: { redirect: '/dashboard' placeholder
      placeholder)
    placeholder)
  placeholder
placeholder)

function mountLogin() {
  return mount(LoginView, {
    global: {
      stubs: {
        AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
        RouterLink: true,
        TurnstileWidget: CaptchaChallengeStub,
        Icon: true,
        LoginAgreementPrompt: true,
        TotpLoginModal: true,
        EmailOAuthButtons: OAuthButtonStub,
        LinuxDoOAuthSection: true,
        DingTalkOAuthSection: true,
        OidcOAuthSection: true,
        WechatOAuthSection: true
      placeholder
    placeholder
  placeholder)
placeholder

describe('Tencent captcha action gate', () => {
  beforeEach(() => {
    loginMock.mockReset()
    loginWithPasskeyMock.mockReset()
    getPublicSettingsMock.mockReset()
    startOAuthLoginMock.mockReset()
    verifyTencentMock.mockReset()
    captchaResetMock.mockReset()
    getPublicSettingsMock.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: '',
      tencent_captcha_enabled: true,
      tencent_captcha_app_id: 'tencent-app-id',
      backend_mode_enabled: false,
      password_reset_enabled: false,
      passkey_enabled: true,
      github_oauth_enabled: true,
      google_oauth_enabled: false
    placeholder)
    loginMock.mockResolvedValue({placeholder)
    loginWithPasskeyMock.mockResolvedValue({placeholder)
    startOAuthLoginMock.mockResolvedValue({ authorize_url: 'https://github.example/authorize' placeholder)
    verifyTencentMock.mockResolvedValue({ ticket: 'ticket-1', randstr: '@rand-1' placeholder)
    Object.defineProperty(window, 'PublicKeyCredential', {
      configurable: true,
      value: class PublicKeyCredential {placeholder
    placeholder)
    locationState.href = 'http://localhost/login'
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: locationState
    placeholder)
  placeholder)

  it('clicking login opens Tencent captcha before calling login', async () => {
    const wrapper = mountLogin()
    await flushPromises()
    await wrapper.get('#email').setValue('user@example.com')
    await wrapper.get('#password').setValue('secret-123')

    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(verifyTencentMock).toHaveBeenCalledOnce()
    expect(loginMock).toHaveBeenCalledWith(expect.objectContaining({
      tencent_captcha_ticket: 'ticket-1',
      tencent_captcha_randstr: '@rand-1'
    placeholder))
  placeholder)

  it('does not call login when Tencent captcha is closed', async () => {
    verifyTencentMock.mockResolvedValue(null)
    const wrapper = mountLogin()
    await flushPromises()
    await wrapper.get('#email').setValue('user@example.com')
    await wrapper.get('#password').setValue('secret-123')

    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(verifyTencentMock).toHaveBeenCalledOnce()
    expect(loginMock).not.toHaveBeenCalled()
  placeholder)

  it('does not open Tencent captcha when login form validation fails', async () => {
    const wrapper = mountLogin()
    await flushPromises()

    await wrapper.get('form').trigger('submit')
    await flushPromises()

    expect(verifyTencentMock).not.toHaveBeenCalled()
    expect(loginMock).not.toHaveBeenCalled()
  placeholder)

  it('starts OAuth through the Tencent gate before navigating', async () => {
    const wrapper = mountLogin()
    await flushPromises()

    await wrapper.get('[data-testid="oauth-start"]').trigger('click')
    await flushPromises()

    expect(verifyTencentMock).toHaveBeenCalledOnce()
    expect(startOAuthLoginMock).toHaveBeenCalledWith(
      { provider: 'github', params: { redirect: '/dashboard' placeholder placeholder,
      {
        tencent_captcha_ticket: 'ticket-1',
        tencent_captcha_randstr: '@rand-1'
      placeholder
    )
    expect(locationState.href).toBe('https://github.example/authorize')
    expect(captchaResetMock).toHaveBeenCalledOnce()
  placeholder)

  it('does not start OAuth when Tencent captcha is closed', async () => {
    verifyTencentMock.mockResolvedValue(null)
    const wrapper = mountLogin()
    await flushPromises()

    await wrapper.get('[data-testid="oauth-start"]').trigger('click')
    await flushPromises()

    expect(startOAuthLoginMock).not.toHaveBeenCalled()
    expect(locationState.href).toBe('http://localhost/login')
  placeholder)

  it('passes a fresh Tencent proof to Passkey login', async () => {
    const wrapper = mountLogin()
    await flushPromises()

    await wrapper.get('button.btn-secondary.w-full').trigger('click')
    await flushPromises()

    expect(verifyTencentMock).toHaveBeenCalledOnce()
    expect(loginWithPasskeyMock).toHaveBeenCalledWith({
      tencent_captcha_ticket: 'ticket-1',
      tencent_captcha_randstr: '@rand-1'
    placeholder)
    expect(captchaResetMock).toHaveBeenCalledOnce()
  placeholder)

  it('does not invoke Passkey when Tencent captcha is closed', async () => {
    verifyTencentMock.mockResolvedValue(null)
    const wrapper = mountLogin()
    await flushPromises()

    await wrapper.get('button.btn-secondary.w-full').trigger('click')
    await flushPromises()

    expect(loginWithPasskeyMock).not.toHaveBeenCalled()
  placeholder)
placeholder)
