import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import OidcCallbackView from '../OidcCallbackView.vue'

const replace = vi.fn()
const showSuccess = vi.fn()
const showError = vi.fn()
const setToken = vi.fn()
const exchangePendingOAuthCompletion = vi.fn()
const completeOIDCOAuthRegistration = vi.fn()
const getPublicSettings = vi.fn()

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
    setToken
  placeholder),
  useAppStore: () => ({
    showSuccess,
    showError
  placeholder)
placeholder))

vi.mock('@/api/auth', async () => {
  const actual = await vi.importActual<typeof import('@/api/auth')>('@/api/auth')
  return {
    ...actual,
    exchangePendingOAuthCompletion: (...args: any[]) => exchangePendingOAuthCompletion(...args),
    completeOIDCOAuthRegistration: (...args: any[]) => completeOIDCOAuthRegistration(...args),
    getPublicSettings: (...args: any[]) => getPublicSettings(...args)
  placeholder
placeholder)

describe('OidcCallbackView', () => {
  beforeEach(() => {
    replace.mockReset()
    showSuccess.mockReset()
    showError.mockReset()
    setToken.mockReset()
    exchangePendingOAuthCompletion.mockReset()
    completeOIDCOAuthRegistration.mockReset()
    getPublicSettings.mockReset()
    getPublicSettings.mockResolvedValue({
      oidc_oauth_provider_name: 'ExampleID'
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

    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(1)
    await buttons[0].trigger('click')
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

  it('renders pending email collection ui and routes to register with the entered email', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'email_required',
      redirect: '/profile',
      provider_fallback: 'ExampleID'
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

    expect(setToken).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('Continue with email')

    await wrapper.get('input[type="email"]').setValue('alice@example.com')
    await wrapper.get('button').trigger('click')

    expect(replace).toHaveBeenCalledWith({
      path: '/register',
      query: {
        email: 'alice@example.com',
        redirect: '/profile',
        provider: 'ExampleID'
      placeholder
    placeholder)
  placeholder)

  it('renders existing-account binding ui and routes to login', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'existing_account_binding_required',
      redirect: '/profile',
      existing_account_email: 'alice@example.com'
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

    expect(wrapper.text()).toContain('alice@example.com')
    expect(wrapper.text()).toContain('Sign in to bind')

    await wrapper.get('button').trigger('click')

    expect(replace).toHaveBeenCalledWith({
      path: '/login',
      query: {
        email: 'alice@example.com',
        redirect: '/profile',
        provider: 'ExampleID'
      placeholder
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

    expect(wrapper.text()).toContain('OIDC Nick')
    expect(exchangePendingOAuthCompletion).toHaveBeenCalledTimes(1)
    expect(exchangePendingOAuthCompletion).toHaveBeenCalledWith()

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
placeholder)
