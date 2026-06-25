import { afterEach, beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'
import PaymentQRDialog from '../PaymentQRDialog.vue'

const pollOrderStatus = vi.hoisted(() => vi.fn())
const cancelOrder = vi.hoisted(() => vi.fn())
const verifyOrder = vi.hoisted(() => vi.fn())
const showError = vi.hoisted(() => vi.fn())
const toCanvas = vi.hoisted(() => vi.fn())

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    placeholder),
  placeholder
placeholder)

vi.mock('@/stores/payment', () => ({
  usePaymentStore: () => ({
    pollOrderStatus,
  placeholder),
placeholder))

vi.mock('@/stores', () => ({
  useAppStore: () => ({
    showError,
  placeholder),
placeholder))

vi.mock('@/api/payment', () => ({
  paymentAPI: {
    cancelOrder,
    verifyOrder,
  placeholder,
placeholder))

vi.mock('qrcode', () => ({
  default: {
    toCanvas,
  placeholder,
placeholder))

const paidOrder = {
  id: 42,
  user_id: 9,
  amount: 100,
  pay_amount: 108,
  currency: 'CNY',
  fee_rate: 8,
  payment_type: 'alipay',
  out_trade_no: 'sub2_202606250001',
  status: 'COMPLETED',
  order_type: 'subscription',
  created_at: '2026-06-25T10:00:00Z',
  expires_at: '2099-01-01T10:30:00Z',
  refund_amount: 0,
placeholder

describe('PaymentQRDialog currency display', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    pollOrderStatus.mockReset().mockResolvedValue(paidOrder)
    cancelOrder.mockReset()
    verifyOrder.mockReset()
    showError.mockReset()
    toCanvas.mockReset().mockResolvedValue(undefined)
  placeholder)

  afterEach(() => {
    vi.useRealTimers()
  placeholder)

  it('uses order currency for pay_amount and USD for credited amount', async () => {
    const wrapper = mount(PaymentQRDialog, {
      props: {
        show: false,
        orderId: 42,
        qrCode: '',
        expiresAt: '2099-01-01T10:30:00Z',
        paymentType: 'alipay',
      placeholder,
      global: {
        stubs: {
          BaseDialog: {
            props: ['show'],
            template: '<div v-if="show"><slot /><slot name="footer" /></div>',
          placeholder,
          Icon: true,
        placeholder,
      placeholder,
    placeholder)

    await wrapper.setProps({ show: true placeholder)
    await flushPromises()
    await vi.advanceTimersByTimeAsync(3000)
    await flushPromises()

    expect(pollOrderStatus).toHaveBeenCalledWith(42)
    expect(wrapper.text()).toContain('$100.00')
    expect(wrapper.text()).toContain('¥108.00')
  placeholder)
placeholder)
