import { afterEach, beforeEach, describe, expect, it, vi placeholder from 'vitest'

describe('user api oauth binding urls', () => {
  beforeEach(() => {
    vi.resetModules()
    vi.stubEnv('VITE_API_BASE_URL', 'https://api.example.com/api/v1')
  placeholder)

  afterEach(() => {
    vi.unstubAllEnvs()
  placeholder)

  it('builds third-party bind urls against the bind start endpoint', async () => {
    const { buildOAuthBindingStartURL placeholder = await import('@/api/user')

    expect(buildOAuthBindingStartURL('linuxdo', { redirectTo: '/settings/profile' placeholder)).toBe(
      'https://api.example.com/api/v1/auth/oauth/linuxdo/bind/start?redirect=%2Fsettings%2Fprofile&intent=bind_current_user'
    )
    expect(
      buildOAuthBindingStartURL('wechat', {
        redirectTo: '/settings/profile',
        wechatOAuthSettings: {
          wechat_oauth_open_enabled: true,
          wechat_oauth_mp_enabled: false,
          wechat_oauth_mobile_enabled: false
        placeholder
      placeholder)
    ).toBe(
      'https://api.example.com/api/v1/auth/oauth/wechat/bind/start?redirect=%2Fsettings%2Fprofile&intent=bind_current_user&mode=open'
    )
  placeholder)
placeholder)
