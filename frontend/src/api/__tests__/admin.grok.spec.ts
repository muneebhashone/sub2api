import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { post placeholder = vi.hoisted(() => ({
  post: vi.fn(),
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: { post placeholder,
placeholder))

import { authorizePassword, createFromSSO, getGrokSSOImportTimeout placeholder from '@/api/admin/grok'

describe('admin Grok SSO import API', () => {
  beforeEach(() => {
    post.mockReset()
    post.mockResolvedValue({ data: { created: [], failed: [] placeholder placeholder)
  placeholder)

  it.each([
    [1, 180_000],
    [3, 180_000],
    [4, 270_000],
    [7, 360_000],
  ])('uses a timeout sized for %i keys', async (keyCount, expectedTimeout) => {
    expect(getGrokSSOImportTimeout(keyCount)).toBe(expectedTimeout)

    await createFromSSO({
      sso_tokens: Array.from({ length: keyCount placeholder, (_, index) => `sso-${index + 1placeholder`),
    placeholder)

    expect(post).toHaveBeenCalledWith(
      '/admin/grok/sso-to-oauth',
      expect.objectContaining({ sso_tokens: expect.any(Array) placeholder),
      { timeout: expectedTimeout placeholder,
    )
  placeholder)

  it('preserves password whitespace and applies the authorization timeout', async () => {
    post.mockResolvedValueOnce({ data: { access_token: 'access-token' placeholder placeholder)

    await authorizePassword(' user@example.com ----  password with spaces  ', 7)

    expect(post).toHaveBeenCalledWith(
      '/admin/grok/oauth/password',
      {
        email: 'user@example.com',
        password: '  password with spaces  ',
        proxy_id: 7,
      placeholder,
      { timeout: 120_000 placeholder,
    )
  placeholder)
placeholder)
