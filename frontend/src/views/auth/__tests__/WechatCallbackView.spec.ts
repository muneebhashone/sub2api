import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import WechatCallbackView from '@/views/auth/WechatCallbackView.vue'

const {
  postMock,
  replaceMock,
  setTokenMock,
  showSuccessMock,
  showErrorMock,
  routeState,
placeholder = vi.hoisted(() => ({
  postMock: vi.fn(),
  replaceMock: vi.fn(),
  setTokenMock: vi.fn(),
  showSuccessMock: vi.fn(),
  showErrorMock: vi.fn(),
  routeState: {
    query: {placeholder as Record<string, unknown>,
  placeholder,
placeholder))

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
  useRouter: () => ({
    replace: replaceMock,
  placeholder),
placeholder))

vi.mock('vue-i18n', () => ({
  createI18n: () => ({
    global: {
      t: (key: string) => key,
    placeholder,
  placeholder),
  useI18n: () => ({
    t: (key: string, params?: Record<string, string>) => {
      if (key === 'auth.oidc.callbackTitle') {
        return `Signing you in with ${params?.providerName ?? ''placeholder`.trim()
      placeholder
      if (key === 'auth.oidc.callbackProcessing') {
        return `Completing login with ${params?.providerName ?? ''placeholder`.trim()
      placeholder
      if (key === 'auth.oidc.invitationRequired') {
        return `${params?.providerName ?? ''placeholder invitation required`.trim()
      placeholder
      if (key === 'auth.oidc.completeRegistration') {
        return 'Complete registration'
      placeholder
      if (key === 'auth.oidc.completing') {
        return 'Completing'
      placeholder
      if (key === 'auth.oidc.backToLogin') {
        return 'Back to login'
      placeholder
      if (key === 'auth.invitationCodePlaceholder') {
        return 'Invitation code'
      placeholder
      if (key === 'auth.loginSuccess') {
        return 'Login success'
      placeholder
      if (key === 'auth.loginFailed') {
        return 'Login failed'
      placeholder
      if (key === 'auth.oidc.callbackHint') {
        return 'Callback hint'
      placeholder
      if (key === 'auth.oidc.callbackMissingToken') {
        return 'Missing login token'
      placeholder
      if (key === 'auth.oidc.completeRegistrationFailed') {
        return 'Complete registration failed'
      placeholder
      return key
    placeholder,
  placeholder),
placeholder))

vi.mock('@/stores', () => ({
  useAuthStore: () => ({
    setToken: setTokenMock,
  placeholder),
  useAppStore: () => ({
    showSuccess: showSuccessMock,
    showError: showErrorMock,
  placeholder),
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: {
    post: postMock,
  placeholder,
placeholder))

describe('WechatCallbackView', () => {
  beforeEach(() => {
    postMock.mockReset()
    replaceMock.mockReset()
    setTokenMock.mockReset()
    showSuccessMock.mockReset()
    showErrorMock.mockReset()
    routeState.query = {placeholder
    localStorage.clear()
  placeholder)

  it('does not send adoption decisions during the initial exchange', async () => {
    postMock.mockResolvedValueOnce({
      data: {
        access_token: 'access-token',
        refresh_token: 'refresh-token',
        expires_in: 3600,
        redirect: '/dashboard',
        adoption_required: true,
      placeholder,
    placeholder)
    setTokenMock.mockResolvedValue({placeholder)

    mount(WechatCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(postMock).toHaveBeenCalledWith('/auth/oauth/pending/exchange', {placeholder)
    expect(postMock).toHaveBeenCalledTimes(1)
  placeholder)

  it('waits for explicit adoption confirmation before finishing a non-invitation login', async () => {
    postMock
      .mockResolvedValueOnce({
        data: {
          redirect: '/dashboard',
          adoption_required: true,
          suggested_display_name: 'WeChat Nick',
          suggested_avatar_url: 'https://cdn.example/wechat.png',
        placeholder,
      placeholder)
      .mockResolvedValueOnce({
        data: {
          access_token: 'wechat-access-token',
          refresh_token: 'wechat-refresh-token',
          expires_in: 3600,
          token_type: 'Bearer',
          redirect: '/dashboard',
        placeholder,
      placeholder)
    setTokenMock.mockResolvedValue({placeholder)

    const wrapper = mount(WechatCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(wrapper.text()).toContain('WeChat Nick')
    expect(setTokenMock).not.toHaveBeenCalled()
    expect(replaceMock).not.toHaveBeenCalled()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)
    await checkboxes[1].setValue(false)

    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(1)
    await buttons[0].trigger('click')
    await flushPromises()

    expect(postMock).toHaveBeenNthCalledWith(1, '/auth/oauth/pending/exchange', {placeholder)
    expect(postMock).toHaveBeenNthCalledWith(2, '/auth/oauth/pending/exchange', {
      adopt_display_name: true,
      adopt_avatar: false,
    placeholder)
    expect(setTokenMock).toHaveBeenCalledWith('wechat-access-token')
    expect(replaceMock).toHaveBeenCalledWith('/dashboard')
    expect(localStorage.getItem('refresh_token')).toBe('wechat-refresh-token')
  placeholder)

  it('renders adoption choices for invitation flow and submits the selected values', async () => {
    postMock
      .mockResolvedValueOnce({
        data: {
          error: 'invitation_required',
          redirect: '/subscriptions',
          adoption_required: true,
          suggested_display_name: 'WeChat Nick',
          suggested_avatar_url: 'https://cdn.example/wechat.png',
        placeholder,
      placeholder)
      .mockResolvedValueOnce({
        data: {
          access_token: 'wechat-invite-token',
          refresh_token: 'wechat-invite-refresh',
          expires_in: 600,
          token_type: 'Bearer',
        placeholder,
      placeholder)

    const wrapper = mount(WechatCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false,
        placeholder,
      placeholder,
    placeholder)

    await flushPromises()

    expect(wrapper.text()).toContain('WeChat Nick')
    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)
    await checkboxes[0].setValue(false)
    await wrapper.get('input[type="text"]').setValue(' INVITE-CODE ')
    await wrapper.get('button').trigger('click')
    await flushPromises()

    expect(postMock).toHaveBeenNthCalledWith(2, '/auth/oauth/wechat/complete-registration', {
      invitation_code: 'INVITE-CODE',
      adopt_display_name: false,
      adopt_avatar: true,
    placeholder)
    expect(setTokenMock).toHaveBeenCalledWith('wechat-invite-token')
    expect(replaceMock).toHaveBeenCalledWith('/subscriptions')
  placeholder)
placeholder)
