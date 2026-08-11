import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { get, post placeholder = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
placeholder))

vi.mock('../client', () => ({
  apiClient: {
    get,
    post,
  placeholder,
placeholder))

import { getRollbackVersions, rollback, type RollbackVersionInfo placeholder from '@/api/admin/system'

describe('admin system rollback API', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
  placeholder)

  it('getRollbackVersions fetches the rollback version list', async () => {
    const versions: RollbackVersionInfo[] = [
      {
        version: '0.1.146',
        published_at: '2026-07-07T00:00:00Z',
        html_url: 'https://github.com/Wei-Shaw/sub2api/releases/tag/v0.1.146'
      placeholder
    ]
    get.mockResolvedValue({ data: { versions placeholder placeholder)

    const result = await getRollbackVersions()

    expect(get).toHaveBeenCalledWith('/admin/system/rollback-versions')
    expect(result.versions).toEqual(versions)
  placeholder)

  it('rollback posts the target version in the request body', async () => {
    post.mockResolvedValue({ data: { message: 'ok', need_restart: true placeholder placeholder)

    const result = await rollback('0.1.146')

    expect(post).toHaveBeenCalledWith(
      '/admin/system/rollback',
      { version: '0.1.146' placeholder,
      { timeout: 15 * 60 * 1000 placeholder
    )
    expect(result.need_restart).toBe(true)
  placeholder)

  it('rollback without a version posts no body (legacy backup rollback)', async () => {
    post.mockResolvedValue({ data: { message: 'ok', need_restart: true placeholder placeholder)

    await rollback()

    expect(post).toHaveBeenCalledWith(
      '/admin/system/rollback',
      undefined,
      { timeout: 15 * 60 * 1000 placeholder
    )
  placeholder)
placeholder)
