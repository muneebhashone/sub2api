import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { get, post placeholder = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: {
    get,
    post,
  placeholder,
placeholder))

import { paymentAPI placeholder from '@/api/payment'

describe('payment api', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    get.mockResolvedValue({ data: {placeholder placeholder)
    post.mockResolvedValue({ data: {placeholder placeholder)
  placeholder)

  it('keeps legacy public out_trade_no verification for upgrade compatibility', async () => {
    await paymentAPI.verifyOrderPublic('legacy-order-no')

    expect(post).toHaveBeenCalledWith('/payment/public/orders/verify', {
      out_trade_no: 'legacy-order-no',
    placeholder)
  placeholder)

  it('keeps signed public resume-token resolve endpoint', async () => {
    await paymentAPI.resolveOrderPublicByResumeToken('resume-token-123')

    expect(post).toHaveBeenCalledWith('/payment/public/orders/resolve', {
      resume_token: 'resume-token-123',
    placeholder)
  placeholder)
placeholder)
