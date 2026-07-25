import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import RegisterView from '@/views/auth/RegisterView.vue'

const { getPublicSettingsMock placeholder = vi.hoisted(() => ({
  getPublicSettingsMock: vi.fn()
placeholder))

const publicSettings = {
  registration_enabled: true,
  email_verify_enabled: false,
  promo_code_enabled: false,
  invitation_code_enabled: false,
  affiliate_enabled: true,
  turnstile_enabled: true,
  turnstile_site_key: 'site-key',
  site_name: 'Sub2API',
  registration_email_suffix_whitelist: [],
  linuxdo_oauth_enabled: false,
  wechat_oauth_enabled: false,
  oidc_oauth_enabled: false,
  github_oauth_enabled: false,
  google_oauth_enabled: false
placeholder

vi.mock('vue-router', () => ({
  useRouter: () => ({ push: vi.fn() placeholder),
  useRoute: () => ({ query: {placeholder placeholder)
placeholder))

vi.mock('vue-i18n', () => ({
  createI18n: () => ({
    global: {
      t: (key: string) => key
    placeholder
  placeholder),
  useI18n: () => ({
    t: (key: string) => key,
    locale: { value: 'en' placeholder
  placeholder)
placeholder))

vi.mock('@/stores', () => ({
  useAuthStore: () => ({ register: vi.fn() placeholder),
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
    getPublicSettings: (...args: unknown[]) => getPublicSettingsMock(...args)
  placeholder
placeholder)

function mountRegister() {
  return mount(RegisterView, {
    global: {
      stubs: {
        AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
        Icon: true,
        TurnstileWidget: { template: '<div data-testid="turnstile-widget" />' placeholder,
        LoginAgreementPrompt: true,
        EmailOAuthButtons: true,
        LinuxDoOAuthSection: true,
        WechatOAuthSection: true,
        OidcOAuthSection: true,
        RouterLink: true,
        transition: false
      placeholder
    placeholder
  placeholder)
placeholder

describe('RegisterView invitation layout', () => {
  beforeEach(() => {
    getPublicSettingsMock.mockReset()
    getPublicSettingsMock.mockResolvedValue(publicSettings)
  placeholder)

  it('keeps the optional affiliate invitation field before Turnstile', async () => {
    const wrapper = mountRegister()
    await flushPromises()

    const invitationField = wrapper.get('[data-testid="affiliate-invitation-field"]')
    const turnstile = wrapper.get('[data-testid="registration-turnstile"]')

    expect(invitationField.get('input').attributes('id')).toBe('affiliate_code')
    expect(invitationField.text()).toContain('common.optional')
    expect(
      invitationField.element.compareDocumentPosition(turnstile.element) &
        Node.DOCUMENT_POSITION_FOLLOWING
    ).toBeTruthy()
  placeholder)

  it('uses the mandatory invitation field without duplicating the affiliate field', async () => {
    getPublicSettingsMock.mockResolvedValueOnce({
      ...publicSettings,
      invitation_code_enabled: true
    placeholder)

    const wrapper = mountRegister()
    await flushPromises()

    expect(wrapper.find('[data-testid="affiliate-invitation-field"]').exists()).toBe(false)
    expect(wrapper.get('#invitation_code').exists()).toBe(true)
  placeholder)
placeholder)
