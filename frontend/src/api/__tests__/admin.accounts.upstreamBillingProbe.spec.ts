import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { get, post, put placeholder = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn()
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: { get, post, put placeholder
placeholder))

import {
  getUpstreamBillingProbeSettings,
  probeUpstreamBilling,
  probeUpstreamBillingBatch,
  setUpstreamBillingProbeEnabled,
  updateUpstreamBillingProbeSettings
placeholder from '@/api/admin/accounts'

describe('admin account upstream billing probe API', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    put.mockReset()
  placeholder)

  it('reads and updates global settings', async () => {
    const settings = { enabled: true, interval_minutes: 30 placeholder
    get.mockResolvedValueOnce({ data: settings placeholder)
    put.mockResolvedValueOnce({ data: settings placeholder)

    await expect(getUpstreamBillingProbeSettings()).resolves.toEqual(settings)
    await expect(updateUpstreamBillingProbeSettings(settings)).resolves.toEqual(settings)
    expect(get).toHaveBeenCalledWith('/admin/accounts/upstream-billing-probe/settings')
    expect(put).toHaveBeenCalledWith('/admin/accounts/upstream-billing-probe/settings', settings)
  placeholder)

  it('uses dedicated account and batch endpoints', async () => {
    const result = { account_id: 7, snapshot: { status: 'unsupported' placeholder placeholder
    put.mockResolvedValueOnce({ data: {placeholder placeholder)
    post.mockResolvedValueOnce({ data: result placeholder)
    post.mockResolvedValueOnce({ data: { results: [result] placeholder placeholder)

    await setUpstreamBillingProbeEnabled(7, true)
    await expect(probeUpstreamBilling(7)).resolves.toEqual(result)
    await expect(probeUpstreamBillingBatch([7])).resolves.toEqual([result])

    expect(put).toHaveBeenCalledWith('/admin/accounts/7/upstream-billing-probe', { enabled: true placeholder)
    expect(post).toHaveBeenNthCalledWith(1, '/admin/accounts/7/upstream-billing-probe')
    expect(post).toHaveBeenNthCalledWith(2, '/admin/accounts/upstream-billing-probe/batch', { account_ids: [7] placeholder)
  placeholder)
placeholder)
