import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import OidcCallbackView from '../OidcCallbackView.vue'

const replace = vi.fn()
const showSuccess = vi.fn()
const showError = vi.fn()
const setToken = vi.fn()
const setPendingAuthSession = vi.fn()
const clearPendingAuthSession = vi.fn()
const exchangePendingOAuthCompletion = vi.fn()
const completeOIDCOAuthRegistration = vi.fn()
const getPublicSettings = vi.fn()
const login2FA = vi.fn()
const apiClientPost = vi.fn()
const sendVerifyCode = vi.fn()

vi.mock('vue-router', () => ({
  useRoute: () => ({
    query: {placeholder
  placeholder),
  useRouter: () => ({
    replace
  placeholder)
placeholder))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, params?: Record<string, string>) => {
        if (!params?.providerName) {
          return key
        placeholder
        return `${keyplaceholder:${params.providerNameplaceholder`
      placeholder
    placeholder)
  placeholder
placeholder)

vi.mock('@/stores', () => ({
  useAuthStore: () => ({
    setToken,
    setPendingAuthSession,
    clearPendingAuthSession
  placeholder),
  useAppStore: () => ({
    showSuccess,
    showError
  placeholder)
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: {
    post: (...args: any[]) => apiClientPost(...args)
  placeholder
placeholder))

vi.mock('@/api/auth', async () => {
  const actual = await vi.importActual<typeof import('@/api/auth')>('@/api/auth')
  return {
    ...actual,
    exchangePendingOAuthCompletion: (...args: any[]) => exchangePendingOAuthCompletion(...args),
    completeOIDCOAuthRegistration: (...args: any[]) => completeOIDCOAuthRegistration(...args),
    getPublicSettings: (...args: any[]) => getPublicSettings(...args),
    login2FA: (...args: any[]) => login2FA(...args),
    sendVerifyCode: (...args: any[]) => sendVerifyCode(...args)
  placeholder
placeholder)

describe('OidcCallbackView', () => {
  beforeEach(() => {
    replace.mockReset()
    showSuccess.mockReset()
    showError.mockReset()
    setToken.mockReset()
    setPendingAuthSession.mockReset()
    clearPendingAuthSession.mockReset()
    exchangePendingOAuthCompletion.mockReset()
    completeOIDCOAuthRegistration.mockReset()
    getPublicSettings.mockReset()
    login2FA.mockReset()
    apiClientPost.mockReset()
    sendVerifyCode.mockReset()
    getPublicSettings.mockResolvedValue({
      oidc_oauth_provider_name: 'ExampleID',
      turnstile_enabled: false,
      turnstile_site_key: ''
    placeholder)
  placeholder)

  it('does not send adoption decisions during the initial exchange', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      expires_in: 3600,
      redirect: '/dashboard',
      adoption_required: true
    placeholder)
    setToken.mockResolvedValue({placeholder)

    mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(exchangePendingOAuthCompletion).toHaveBeenCalledTimes(1)
    expect(exchangePendingOAuthCompletion).toHaveBeenCalledWith()
  placeholder)

  it('waits for explicit adoption confirmation before finishing a non-invitation login', async () => {
    exchangePendingOAuthCompletion
      .mockResolvedValueOnce({
        redirect: '/dashboard',
        adoption_required: true,
        suggested_display_name: 'OIDC Nick',
        suggested_avatar_url: 'https://cdn.example/oidc.png'
      placeholder)
      .mockResolvedValueOnce({
        access_token: 'access-token',
        refresh_token: 'refresh-token',
        expires_in: 3600,
        redirect: '/dashboard'
      placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(wrapper.text()).toContain('OIDC Nick')
    expect(setToken).not.toHaveBeenCalled()
    expect(replace).not.toHaveBeenCalled()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    await checkboxes[0].setValue(false)

    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()

    expect(exchangePendingOAuthCompletion).toHaveBeenCalledTimes(2)
    expect(exchangePendingOAuthCompletion).toHaveBeenNthCalledWith(1)
    expect(exchangePendingOAuthCompletion).toHaveBeenNthCalledWith(2, {
      adoptDisplayName: false,
      adoptAvatar: true
    placeholder)
    expect(setToken).toHaveBeenCalledWith('access-token')
    expect(replace).toHaveBeenCalledWith('/dashboard')
  placeholder)

  it('supports bind completion after adoption confirmation', async () => {
    exchangePendingOAuthCompletion
      .mockResolvedValueOnce({
        redirect: '/dashboard',
        adoption_required: true,
        suggested_display_name: 'OIDC Nick',
        suggested_avatar_url: 'https://cdn.example/oidc.png'
      placeholder)
      .mockResolvedValueOnce({
        redirect: '/profile'
      placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()

    expect(exchangePendingOAuthCompletion).toHaveBeenNthCalledWith(2, {
      adoptDisplayName: true,
      adoptAvatar: true
    placeholder)
    expect(setToken).not.toHaveBeenCalled()
    expect(showSuccess).toHaveBeenCalledWith('profile.authBindings.bindSuccess')
    expect(replace).toHaveBeenCalledWith('/profile')
  placeholder)

  it('keeps rendering pending bind-login UI when adoption confirmation leads to another pending step', async () => {
    exchangePendingOAuthCompletion
      .mockResolvedValueOnce({
        redirect: '/profile',
        adoption_required: true,
        suggested_display_name: 'OIDC Nick',
        suggested_avatar_url: 'https://cdn.example/oidc.png'
      placeholder)
      .mockResolvedValueOnce({
        step: 'bind_login_required',
        redirect: '/profile',
        email: 'existing@example.com',
        adoption_required: true,
        suggested_display_name: 'OIDC Nick',
        suggested_avatar_url: 'https://cdn.example/oidc.png'
      placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()
    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()

    expect(showSuccess).not.toHaveBeenCalled()
    expect(replace).not.toHaveBeenCalled()
    expect((wrapper.get('[data-testid="oidc-bind-login-email"]').element as HTMLInputElement).value).toBe(
      'existing@example.com'
    )
  placeholder)

  it('persists a pending auth session when the oauth flow still needs account creation', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'email_required',
      redirect: '/welcome'
    placeholder)

    mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(setPendingAuthSession).toHaveBeenCalledWith({
      token: '',
      token_field: 'pending_oauth_token',
      provider: 'oidc',
      redirect: '/welcome'
    placeholder)
  placeholder)

  it('renders adoption choices for invitation flow and submits the selected values', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'invitation_required',
      redirect: '/dashboard',
      adoption_required: true,
      suggested_display_name: 'OIDC Nick',
      suggested_avatar_url: 'https://cdn.example/oidc.png'
    placeholder)
    completeOIDCOAuthRegistration.mockResolvedValue({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      expires_in: 3600,
      token_type: 'Bearer'
    placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)
    await checkboxes[1].setValue(false)
    await wrapper.find('input[type="text"]').setValue('invite-code')
    await wrapper.find('button').trigger('click')

    expect(completeOIDCOAuthRegistration).toHaveBeenCalledWith('invite-code', {
      adoptDisplayName: true,
      adoptAvatar: false
    placeholder)
  placeholder)

  it('collects email, password, and verify code for pending oauth account creation and submits adoption decisions', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'email_required',
      redirect: '/welcome',
      adoption_required: true,
      suggested_display_name: 'OIDC Nick',
      suggested_avatar_url: 'https://cdn.example/oidc.png'
    placeholder)
    apiClientPost.mockResolvedValue({
      data: {
        access_token: 'new-access-token',
        refresh_token: 'new-refresh-token',
        expires_in: 3600,
        token_type: 'Bearer'
      placeholder
    placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)
    await checkboxes[1].setValue(false)
    await wrapper.get('[data-testid="oidc-create-account-email"]').setValue('  new@example.com  ')
    await wrapper.get('[data-testid="oidc-create-account-password"]').setValue('secret-123')
    await wrapper.get('[data-testid="oidc-create-account-verify-code"]').setValue('246810')
    await wrapper.get('[data-testid="oidc-create-account-submit"]').trigger('click')
    await flushPromises()

    expect(apiClientPost).toHaveBeenCalledWith('/auth/oauth/pending/create-account', {
      email: 'new@example.com',
      password: 'secret-123',
      verify_code: '246810',
      adopt_display_name: true,
      adopt_avatar: false
    placeholder)
    expect(setToken).toHaveBeenCalledWith('new-access-token')
    expect(replace).toHaveBeenCalledWith('/welcome')
  placeholder)

  it('sends a verify code for pending oauth account creation', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'email_required',
      redirect: '/welcome'
    placeholder)
    sendVerifyCode.mockResolvedValue({
      message: 'sent',
      countdown: 60
    placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    await wrapper.get('[data-testid="oidc-create-account-email"]').setValue('  new@example.com  ')
    await wrapper.get('[data-testid="oidc-create-account-send-code"]').trigger('click')
    await flushPromises()

    expect(sendVerifyCode).toHaveBeenCalledWith({
      email: 'new@example.com'
    placeholder)
  placeholder)

  it('shows bind-login form for existing account binding and submits credentials with adoption decisions', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'adopt_existing_user_by_email',
      redirect: '/profile/security',
      email: 'existing@example.com',
      adoption_required: true,
      suggested_display_name: 'OIDC Nick',
      suggested_avatar_url: 'https://cdn.example/oidc.png'
    placeholder)
    apiClientPost.mockResolvedValue({
      data: {
        access_token: 'bind-access-token',
        refresh_token: 'bind-refresh-token',
        expires_in: 3600,
        token_type: 'Bearer'
      placeholder
    placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)
    await checkboxes[0].setValue(false)
    await wrapper.get('[data-testid="oidc-bind-login-email"]').setValue('existing@example.com')
    await wrapper.get('[data-testid="oidc-bind-login-password"]').setValue('secret-password')
    await wrapper.get('[data-testid="oidc-bind-login-submit"]').trigger('click')
    await flushPromises()

    expect(apiClientPost).toHaveBeenCalledWith('/auth/oauth/pending/bind-login', {
      email: 'existing@example.com',
      password: 'secret-password',
      adopt_display_name: false,
      adopt_avatar: true
    placeholder)
    expect(setToken).toHaveBeenCalledWith('bind-access-token')
    expect(replace).toHaveBeenCalledWith('/profile/security')
  placeholder)

  it('handles bind-login 2FA challenge before redirecting', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'adopt_existing_user_by_email',
      redirect: '/profile',
      email: 'existing@example.com',
      adoption_required: true,
      suggested_display_name: 'OIDC Nick',
      suggested_avatar_url: 'https://cdn.example/oidc.png'
    placeholder)
    apiClientPost.mockResolvedValue({
      data: {
        requires_2fa: true,
        temp_token: 'temp-123',
        user_email_masked: 'o***g@example.com'
      placeholder
    placeholder)
    login2FA.mockResolvedValue({
      access_token: '2fa-access-token'
    placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(OidcCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    await wrapper.get('[data-testid="oidc-bind-login-password"]').setValue('secret-password')
    await wrapper.get('[data-testid="oidc-bind-login-submit"]').trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('o***g@example.com')
    expect(login2FA).not.toHaveBeenCalled()

    await wrapper.get('[data-testid="oidc-bind-login-totp"]').setValue('123456')
    await wrapper.get('[data-testid="oidc-bind-login-totp-submit"]').trigger('click')
    await flushPromises()

    expect(login2FA).toHaveBeenCalledWith({
      temp_token: 'temp-123',
      totp_code: '123456'
    placeholder)
    expect(setToken).toHaveBeenCalledWith('2fa-access-token')
    expect(replace).toHaveBeenCalledWith('/profile')
  placeholder)
placeholder)
