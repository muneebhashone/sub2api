import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { post placeholder = vi.hoisted(() => ({
  post: vi.fn(),
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: { post placeholder,
placeholder))

import { createFromSSO, getGrokSSOImportTimeout placeholder from '@/api/admin/grok'

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
placeholder)
