import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const post = vi.hoisted(() => vi.fn())

vi.mock('@/api/client', () => ({
  apiClient: { post placeholder
placeholder))

import {
  buildOAuthLoginStartURL,
  startOAuthLogin,
  type OAuthLoginStart
placeholder from '@/api/auth'

describe('OAuth captcha start API', () => {
  beforeEach(() => {
    post.mockReset()
  placeholder)

  it('posts Tencent captcha proof and preserves OAuth query parameters', async () => {
    const request: OAuthLoginStart = {
      provider: 'github',
      params: { redirect: '/dashboard', aff_code: 'AFF123' placeholder
    placeholder
    const proof = {
      tencent_captcha_ticket: 'ticket',
      tencent_captcha_randstr: '@rand'
    placeholder
    post.mockResolvedValue({ data: { authorize_url: 'https://github.com/login/oauth/authorize' placeholder placeholder)

    await expect(startOAuthLogin(request, proof)).resolves.toEqual({
      authorize_url: 'https://github.com/login/oauth/authorize'
    placeholder)
    expect(post).toHaveBeenCalledWith('/auth/oauth/github/start', proof, {
      params: request.params
    placeholder)
  placeholder)

  it('builds the legacy GET start URL when Tencent captcha is disabled', () => {
    expect(buildOAuthLoginStartURL({
      provider: 'wechat',
      params: { mode: 'open', redirect: '/billing?plan=pro' placeholder
    placeholder)).toBe('/api/v1/auth/oauth/wechat/start?mode=open&redirect=%2Fbilling%3Fplan%3Dpro')
  placeholder)
placeholder)
