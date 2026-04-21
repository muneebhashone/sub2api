import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import ProfileView from '@/views/user/ProfileView.vue'

const {
  fetchPublicSettingsMock,
  refreshUserMock,
  authState
placeholder = vi.hoisted(() => ({
  fetchPublicSettingsMock: vi.fn(),
  refreshUserMock: vi.fn(),
  authState: {
    user: null as Record<string, unknown> | null,
    refreshUser: vi.fn()
  placeholder
placeholder))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => authState
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    fetchPublicSettings: fetchPublicSettingsMock
  placeholder)
placeholder))

vi.mock('@/utils/format', () => ({
  formatDate: () => 'April 2026'
placeholder))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    placeholder)
  placeholder
placeholder)

describe('ProfileView', () => {
  beforeEach(() => {
    refreshUserMock.mockReset()
    fetchPublicSettingsMock.mockReset()
    refreshUserMock.mockResolvedValue(undefined)
    authState.refreshUser = refreshUserMock
    authState.user = {
      id: 1,
      username: 'alice',
      email: 'alice@example.com',
      role: 'user',
      balance: 10,
      concurrency: 2,
      status: 'active',
      allowed_groups: null,
      balance_notify_enabled: true,
      balance_notify_threshold: null,
      balance_notify_extra_emails: [],
      created_at: '2026-04-20T00:00:00Z',
      updated_at: '2026-04-20T00:00:00Z'
    placeholder
    fetchPublicSettingsMock.mockResolvedValue({
      contact_info: '',
      balance_low_notify_enabled: false,
      balance_low_notify_threshold: 0,
      linuxdo_oauth_enabled: true,
      wechat_oauth_enabled: true,
      wechat_oauth_open_enabled: true,
      wechat_oauth_mp_enabled: false,
      oidc_oauth_enabled: true,
      oidc_oauth_provider_name: 'OIDC'
    placeholder)
  placeholder)

  it('renders info, avatar, and account binding cards as separate sections', async () => {
    const wrapper = mount(ProfileView, {
      global: {
        stubs: {
          AppLayout: { template: '<div><slot /></div>' placeholder,
          StatCard: { template: '<div class="stat-card" />' placeholder,
          ProfileInfoCard: { template: '<div data-testid="profile-info-card" />' placeholder,
          ProfileAvatarCard: { template: '<div data-testid="profile-avatar-card" />' placeholder,
          ProfileAccountBindingsCard: { template: '<div data-testid="profile-account-bindings-card" />' placeholder,
          ProfileEditForm: true,
          ProfileBalanceNotifyCard: true,
          ProfilePasswordForm: true,
          ProfileTotpCard: true,
          Icon: true
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    const html = wrapper.html()
    expect(html.indexOf('profile-info-card')).toBeGreaterThan(-1)
    expect(html.indexOf('profile-avatar-card')).toBeGreaterThan(html.indexOf('profile-info-card'))
    expect(html.indexOf('profile-account-bindings-card')).toBeGreaterThan(html.indexOf('profile-avatar-card'))
  placeholder)
placeholder)
