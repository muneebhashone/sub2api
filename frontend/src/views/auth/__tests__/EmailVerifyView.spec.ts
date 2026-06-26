import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import EmailVerifyView from '@/views/auth/EmailVerifyView.vue'

const {
  pushMock,
  showSuccessMock,
  showErrorMock,
  registerMock,
  setTokenMock,
  setPendingAuthSessionMock,
  clearPendingAuthSessionMock,
  getPublicSettingsMock,
  sendVerifyCodeMock,
  sendPendingOAuthVerifyCodeMock,
  persistOAuthTokenContextMock,
  apiClientPostMock,
  authStoreState,
placeholder = vi.hoisted(() => ({
  pushMock: vi.fn(),
  showSuccessMock: vi.fn(),
  showErrorMock: vi.fn(),
  registerMock: vi.fn(),
  setTokenMock: vi.fn(),
  setPendingAuthSessionMock: vi.fn(),
  clearPendingAuthSessionMock: vi.fn(),
  getPublicSettingsMock: vi.fn(),
  sendVerifyCodeMock: vi.fn(),
  sendPendingOAuthVerifyCodeMock: vi.fn(),
  persistOAuthTokenContextMock: vi.fn(),
  apiClientPostMock: vi.fn(),
  authStoreState: {
    pendingAuthSession: null as null | {
      token: string
      token_field: 'pending_auth_token' | 'pending_oauth_token'
      provider: string
      redirect?: string
      adoption_required?: boolean
      suggested_display_name?: string
      suggested_avatar_url?: string
    placeholder
  placeholder,
placeholder))

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: pushMock,
  placeholder),
placeholder))

vi.mock('vue-i18n', () => ({
  createI18n: () => ({
    global: {
      t: (key: string) => key,
    placeholder,
  placeholder),
  useI18n: () => ({
    t: (key: string, params?: Record<string, string | number>) => {
      if (key === 'auth.accountCreatedSuccess') {
        return `Account created for ${params?.siteName ?? 'Sub2API'placeholder`
      placeholder
      return key
    placeholder,
    locale: { value: 'en' placeholder,
  placeholder),
placeholder))

vi.mock('@/stores', () => ({
  useAuthStore: () => ({
    pendingAuthSession: authStoreState.pendingAuthSession,
    register: (...args: any[]) => registerMock(...args),
    setToken: (...args: any[]) => setTokenMock(...args),
    setPendingAuthSession: (...args: any[]) => setPendingAuthSessionMock(...args),
    clearPendingAuthSession: (...args: any[]) => clearPendingAuthSessionMock(...args),
  placeholder),
  useAppStore: () => ({
    showSuccess: (...args: any[]) => showSuccessMock(...args),
    showError: (...args: any[]) => showErrorMock(...args),
  placeholder),
placeholder))

vi.mock('@/api/auth', async () => {
  const actual = await vi.importActual<typeof import('@/api/auth')>('@/api/auth')
  return {
    ...actual,
    getPublicSettings: (...args: any[]) => getPublicSettingsMock(...args),
    sendVerifyCode: (...args: any[]) => sendVerifyCodeMock(...args),
    sendPendingOAuthVerifyCode: (...args: any[]) => sendPendingOAuthVerifyCodeMock(...args),
    persistOAuthTokenContext: (...args: any[]) => persistOAuthTokenContextMock(...args),
  placeholder
placeholder)

vi.mock('@/api/client', () => ({
  apiClient: {
    post: (...args: any[]) => apiClientPostMock(...args),
  placeholder,
placeholder))

describe('EmailVerifyView', () => {
  beforeEach(() => {
    pushMock.mockReset()
    showSuccessMock.mockReset()
    showErrorMock.mockReset()
    registerMock.mockReset()
    setTokenMock.mockReset()
    setPendingAuthSessionMock.mockReset()
    clearPendingAuthSessionMock.mockReset()
    getPublicSettingsMock.mockReset()
    sendVerifyCodeMock.mockReset()
    sendPendingOAuthVerifyCodeMock.mockReset()
    persistOAuthTokenContextMock.mockReset()
    apiClientPostMock.mockReset()
    authStoreState.pendingAuthSession = null
    sessionStorage.clear()
    localStorage.clear()

    getPublicSettingsMock.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: '',
      site_name: 'Sub2API',
      registration_email_suffix_whitelist: [],
    placeholder)
    sendVerifyCodeMock.mockResolvedValue({ countdown: 60 placeholder)
    sendPendingOAuthVerifyCodeMock.mockResolvedValue({ countdown: 60 placeholder)
    setTokenMock.mockResolvedValue({placeholder)
  placeholder)

  it('uses the pending oauth verify-code endpoint when register data carries a pending auth session', async () => {
    authStoreState.pendingAuthSession = {
      token: 'pending-token-1',
      token_field: 'pending_auth_token',
      provider: 'wechat',
      redirect: '/profile',
    placeholder
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'fresh@example.com',
        password: 'secret-123',
      placeholder)
    )

    mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(sendPendingOAuthVerifyCodeMock).toHaveBeenCalledWith({
      email: 'fresh@example.com',
      pending_auth_token: 'pending-token-1',
    placeholder)
    expect(sendVerifyCodeMock).not.toHaveBeenCalled()
  placeholder)

  it('skips the registration email suffix whitelist for pending oauth verification', async () => {
    authStoreState.pendingAuthSession = {
      token: 'pending-token-2',
      token_field: 'pending_auth_token',
      provider: 'oidc',
      redirect: '/profile',
    placeholder
    getPublicSettingsMock.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: '',
      site_name: 'Sub2API',
      registration_email_suffix_whitelist: ['allowed.com'],
    placeholder)
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'fresh@example.com',
        password: 'secret-123',
      placeholder)
    )

    mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(sendPendingOAuthVerifyCodeMock).toHaveBeenCalledWith({
      email: 'fresh@example.com',
      pending_auth_token: 'pending-token-2',
    placeholder)
    expect(showErrorMock).not.toHaveBeenCalled()
  placeholder)

  it('uses the pending oauth verify-code endpoint when auth store only carries the pending provider', async () => {
    authStoreState.pendingAuthSession = {
      token: '',
      token_field: 'pending_oauth_token',
      provider: 'oidc',
      redirect: '/profile',
    placeholder
    getPublicSettingsMock.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: '',
      site_name: 'Sub2API',
      registration_email_suffix_whitelist: ['allowed.com'],
    placeholder)
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'fresh@example.com',
        password: 'secret-123',
      placeholder)
    )

    mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(sendPendingOAuthVerifyCodeMock).toHaveBeenCalledWith({
      email: 'fresh@example.com',
      pending_oauth_token: undefined,
    placeholder)
    expect(sendVerifyCodeMock).not.toHaveBeenCalled()
    expect(showErrorMock).not.toHaveBeenCalled()
  placeholder)

  it('returns to the oauth callback flow when pending send-code detects an existing account email', async () => {
    authStoreState.pendingAuthSession = {
      token: '',
      token_field: 'pending_oauth_token',
      provider: 'oidc',
      redirect: '/profile/security',
    placeholder
    getPublicSettingsMock.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: '',
      site_name: 'Sub2API',
      registration_email_suffix_whitelist: ['allowed.com'],
    placeholder)
    sendPendingOAuthVerifyCodeMock.mockResolvedValue({
      auth_result: 'pending_session',
      provider: 'oidc',
      redirect: '/profile/security',
    placeholder)
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'fresh@example.com',
        password: 'secret-123',
      placeholder)
    )

    mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(setPendingAuthSessionMock).toHaveBeenCalledWith({
      token: '',
      token_field: 'pending_oauth_token',
      provider: 'oidc',
      redirect: '/profile/security',
    placeholder)
    expect(pushMock).toHaveBeenCalledWith('/auth/oidc/callback')
    expect(showErrorMock).not.toHaveBeenCalled()
  placeholder)

  it('submits pending auth account creation when session storage has no pending metadata but auth store does', async () => {
    authStoreState.pendingAuthSession = {
      token: 'pending-token-1',
      token_field: 'pending_auth_token',
      provider: 'wechat',
      redirect: '/profile',
    placeholder
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'fresh@example.com',
        password: 'secret-123',
        aff_code: 'AFF123',
      placeholder)
    )
    apiClientPostMock.mockResolvedValue({
      data: {
        access_token: 'oauth-access-token',
        refresh_token: 'oauth-refresh-token',
        expires_in: 3600,
        token_type: 'Bearer',
      placeholder,
    placeholder)

    const wrapper = mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()
    await wrapper.get('#code').setValue('123456')
    await wrapper.get('form').trigger('submit.prevent')
    await flushPromises()

    expect(apiClientPostMock).toHaveBeenCalledWith('/auth/oauth/pending/create-account', {
      email: 'fresh@example.com',
      password: 'secret-123',
      verify_code: '123456',
      aff_code: 'AFF123',
    placeholder)
    expect(persistOAuthTokenContextMock).toHaveBeenCalledWith({
      access_token: 'oauth-access-token',
      refresh_token: 'oauth-refresh-token',
      expires_in: 3600,
      token_type: 'Bearer',
    placeholder)
    expect(setTokenMock).toHaveBeenCalledWith('oauth-access-token')
    expect(clearPendingAuthSessionMock).toHaveBeenCalled()
    expect(pushMock).toHaveBeenCalledWith('/profile')
    expect(registerMock).not.toHaveBeenCalled()
  placeholder)

  it('returns to the oauth callback flow when pending account creation becomes bind-login', async () => {
    authStoreState.pendingAuthSession = {
      token: '',
      token_field: 'pending_oauth_token',
      provider: 'oidc',
      redirect: '/profile/security',
    placeholder
    getPublicSettingsMock.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: '',
      site_name: 'Sub2API',
      registration_email_suffix_whitelist: ['allowed.com'],
    placeholder)
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'fresh@example.com',
        password: 'secret-123',
      placeholder)
    )
    apiClientPostMock.mockResolvedValue({
      data: {
        auth_result: 'pending_session',
        provider: 'oidc',
        step: 'bind_login_required',
        redirect: '/profile/security',
        email: 'fresh@example.com',
      placeholder,
    placeholder)

    const wrapper = mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()
    await wrapper.get('#code').setValue('123456')
    await wrapper.get('form').trigger('submit.prevent')
    await flushPromises()

    expect(apiClientPostMock).toHaveBeenCalledWith('/auth/oauth/pending/create-account', {
      email: 'fresh@example.com',
      password: 'secret-123',
      verify_code: '123456',
    placeholder)
    expect(setPendingAuthSessionMock).toHaveBeenCalledWith({
      token: '',
      token_field: 'pending_oauth_token',
      provider: 'oidc',
      redirect: '/profile/security',
    placeholder)
    expect(pushMock).toHaveBeenCalledWith('/auth/oidc/callback')
    expect(setTokenMock).not.toHaveBeenCalled()
    expect(persistOAuthTokenContextMock).not.toHaveBeenCalled()
    expect(clearPendingAuthSessionMock).not.toHaveBeenCalled()
    expect(showSuccessMock).not.toHaveBeenCalled()
  placeholder)

  it('keeps the normal email registration flow unchanged', async () => {
    sessionStorage.setItem(
      'register_data',
      JSON.stringify({
        email: 'normal@example.com',
        password: 'secret-456',
        promo_code: 'PROMO',
        invitation_code: 'INVITE',
      placeholder)
    )
    registerMock.mockResolvedValue({placeholder)

    const wrapper = mount(EmailVerifyView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /><slot name="footer" /></div>' placeholder,
          Icon: true,
          TurnstileWidget: true,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()
    await wrapper.get('#code').setValue('654321')
    await wrapper.get('form').trigger('submit.prevent')
    await flushPromises()

    expect(registerMock).toHaveBeenCalledWith({
      email: 'normal@example.com',
      password: 'secret-456',
      verify_code: '654321',
      turnstile_token: undefined,
      promo_code: 'PROMO',
      invitation_code: 'INVITE',
    placeholder)
    expect(apiClientPostMock).not.toHaveBeenCalled()
    expect(pushMock).toHaveBeenCalledWith('/dashboard')
  placeholder)
placeholder)
