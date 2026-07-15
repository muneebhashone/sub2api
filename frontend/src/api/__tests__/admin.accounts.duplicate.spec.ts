import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { post placeholder = vi.hoisted(() => ({
  post: vi.fn()
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: { post placeholder
placeholder))

import { duplicate placeholder from '@/api/admin/accounts'

describe('admin account duplicate API', () => {
  beforeEach(() => {
    post.mockReset()
    post.mockResolvedValue({ data: { id: 43, name: 'primary (Copy)' placeholder placeholder)
    vi.spyOn(globalThis.crypto, 'randomUUID').mockReturnValue('11111111-1111-4111-8111-111111111111')
  placeholder)

  it('sends a stable idempotency key with the duplicate request', async () => {
    const account = await duplicate(42)

    expect(post).toHaveBeenCalledWith('/admin/accounts/42/duplicate', undefined, {
      headers: {
        'Idempotency-Key': 'account-duplicate-42-11111111-1111-4111-8111-111111111111'
      placeholder
    placeholder)
    expect(account).toEqual({ id: 43, name: 'primary (Copy)' placeholder)
  placeholder)

  it('reuses the operation key after an ambiguous failed request', async () => {
    post.mockRejectedValueOnce(new Error('network timeout'))
    await expect(duplicate(99)).rejects.toThrow('network timeout')

    post.mockResolvedValueOnce({ data: { id: 100, name: 'retry (Copy)' placeholder placeholder)
    await duplicate(99)

    expect(post).toHaveBeenCalledTimes(2)
    const firstHeaders = post.mock.calls[0][2].headers
    const secondHeaders = post.mock.calls[1][2].headers
    expect(secondHeaders).toEqual(firstHeaders)
  placeholder)
placeholder)
