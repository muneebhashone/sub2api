import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const post = vi.fn()

vi.mock('@/api/client', () => ({
  apiClient: {
    post
  placeholder
placeholder))

describe('oauth adoption auth api', () => {
  beforeEach(() => {
    post.mockReset()
    post.mockResolvedValue({ data: {placeholder placeholder)
    localStorage.clear()
    document.cookie = 'oauth_bind_access_token=; Max-Age=0; path=/'
  placeholder)

  it('posts adoption decisions when exchanging pending oauth completion', async () => {
    const { exchangePendingOAuthCompletion placeholder = await import('@/api/auth')

    await exchangePendingOAuthCompletion({
      adoptDisplayName: false,
      adoptAvatar: true
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/pending/exchange', {
      adopt_display_name: false,
      adopt_avatar: true
    placeholder)
  placeholder)

  it('posts bind-login decisions when finalizing pending oauth bind flow', async () => {
    const { completePendingOAuthBindLogin placeholder = await import('@/api/auth')

    await completePendingOAuthBindLogin({
      adoptDisplayName: true,
      adoptAvatar: false
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/pending/exchange', {
      adopt_display_name: true,
      adopt_avatar: false
    placeholder)
  placeholder)

  it('posts linuxdo invitation completion with adoption decisions', async () => {
    const { completeLinuxDoOAuthRegistration placeholder = await import('@/api/auth')

    await completeLinuxDoOAuthRegistration('invite-code', {
      adoptDisplayName: true,
      adoptAvatar: false
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/linuxdo/complete-registration', {
      invitation_code: 'invite-code',
      adopt_display_name: true,
      adopt_avatar: false
    placeholder)
  placeholder)

  it('posts linuxdo create-account completion with adoption decisions', async () => {
    const { createPendingLinuxDoOAuthAccount placeholder = await import('@/api/auth')

    await createPendingLinuxDoOAuthAccount('invite-code', {
      adoptDisplayName: false,
      adoptAvatar: true
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/linuxdo/complete-registration', {
      invitation_code: 'invite-code',
      adopt_display_name: false,
      adopt_avatar: true
    placeholder)
  placeholder)

  it('posts oidc invitation completion with adoption decisions', async () => {
    const { completeOIDCOAuthRegistration placeholder = await import('@/api/auth')

    await completeOIDCOAuthRegistration('invite-code', {
      adoptDisplayName: false,
      adoptAvatar: true
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/oidc/complete-registration', {
      invitation_code: 'invite-code',
      adopt_display_name: false,
      adopt_avatar: true
    placeholder)
  placeholder)

  it('posts oidc create-account completion with adoption decisions', async () => {
    const { createPendingOIDCOAuthAccount placeholder = await import('@/api/auth')

    await createPendingOIDCOAuthAccount('invite-code', {
      adoptDisplayName: true,
      adoptAvatar: false
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/oidc/complete-registration', {
      invitation_code: 'invite-code',
      adopt_display_name: true,
      adopt_avatar: false
    placeholder)
  placeholder)

  it('posts wechat invitation completion with adoption decisions', async () => {
    const { completeWeChatOAuthRegistration placeholder = await import('@/api/auth')

    await completeWeChatOAuthRegistration('invite-code', {
      adoptDisplayName: true,
      adoptAvatar: true
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/wechat/complete-registration', {
      invitation_code: 'invite-code',
      adopt_display_name: true,
      adopt_avatar: true
    placeholder)
  placeholder)

  it('posts wechat create-account completion with adoption decisions', async () => {
    const { createPendingWeChatOAuthAccount placeholder = await import('@/api/auth')

    await createPendingWeChatOAuthAccount('invite-code', {
      adoptDisplayName: false,
      adoptAvatar: false
    placeholder)

    expect(post).toHaveBeenCalledWith('/auth/oauth/wechat/complete-registration', {
      invitation_code: 'invite-code',
      adopt_display_name: false,
      adopt_avatar: false
    placeholder)
  placeholder)

  it('classifies oauth completion results as login or bind', async () => {
    const { getOAuthCompletionKind placeholder = await import('@/api/auth')

    expect(getOAuthCompletionKind({ access_token: 'access-token' placeholder)).toBe('login')
    expect(getOAuthCompletionKind({ redirect: '/profile' placeholder)).toBe('bind')
  placeholder)

  it('provides bind-login utility helpers for invitation and suggested profile states', async () => {
    const {
      getPendingOAuthBindLoginKind,
      hasPendingOAuthSuggestedProfile,
      isPendingOAuthCreateAccountRequired
    placeholder = await import('@/api/auth')

    expect(getPendingOAuthBindLoginKind({ access_token: 'access-token' placeholder)).toBe('login')
    expect(getPendingOAuthBindLoginKind({ redirect: '/profile' placeholder)).toBe('bind')
    expect(
      isPendingOAuthCreateAccountRequired({
        error: 'invitation_required'
      placeholder)
    ).toBe(true)
    expect(
      isPendingOAuthCreateAccountRequired({
        error: 'other'
      placeholder)
    ).toBe(false)
    expect(
      hasPendingOAuthSuggestedProfile({
        suggested_display_name: 'OAuth Nick'
      placeholder)
    ).toBe(true)
    expect(
      hasPendingOAuthSuggestedProfile({
        suggested_avatar_url: 'https://cdn.example/avatar.png'
      placeholder)
    ).toBe(true)
    expect(hasPendingOAuthSuggestedProfile({placeholder)).toBe(false)
  placeholder)

  it('requests an HttpOnly oauth bind cookie before redirect binding', async () => {
    localStorage.setItem('auth_token', 'access-token-value')
    const { prepareOAuthBindAccessTokenCookie placeholder = await import('@/api/auth')

    await prepareOAuthBindAccessTokenCookie()

    expect(post).toHaveBeenCalledWith('/auth/oauth/bind-token')
  placeholder)
placeholder)
