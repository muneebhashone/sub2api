import { describe, expect, it, vi placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'
import { nextTick placeholder from 'vue'
import PaymentProviderDialog from '@/components/payment/PaymentProviderDialog.vue'

const messages: Record<string, string> = {
  'admin.settings.payment.providerConfig': 'Credentials',
  'admin.settings.payment.paymentGuideTrigger': 'View payment guide',
  'admin.settings.payment.alipayGuideSummary': 'Desktop prefers QR precreate and falls back to cashier; mobile prefers WAP checkout.',
  'admin.settings.payment.wxpayGuideSummary': 'Desktop prefers Native QR; mobile routes to JSAPI or H5 based on browser context.',
placeholder

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => messages[key] ?? key,
  placeholder),
placeholder))

function mountDialog() {
  return mount(PaymentProviderDialog, {
    props: {
      show: true,
      saving: false,
      editing: null,
      allKeyOptions: [
        { value: 'alipay', label: 'Alipay' placeholder,
        { value: 'wxpay', label: 'WeChat Pay' placeholder,
        { value: 'stripe', label: 'Stripe' placeholder,
      ],
      enabledKeyOptions: [
        { value: 'alipay', label: 'Alipay' placeholder,
        { value: 'wxpay', label: 'WeChat Pay' placeholder,
      ],
      allPaymentTypes: [
        { value: 'alipay', label: 'Alipay' placeholder,
        { value: 'wxpay', label: 'WeChat Pay' placeholder,
      ],
      redirectLabel: 'Redirect',
    placeholder,
    global: {
      stubs: {
        BaseDialog: {
          template: '<div><slot /><slot name="footer" /></div>',
        placeholder,
        Select: {
          props: ['modelValue', 'options', 'disabled'],
          template: '<div />',
        placeholder,
        ToggleSwitch: {
          template: '<div />',
        placeholder,
      placeholder,
    placeholder,
  placeholder)
placeholder

describe('PaymentProviderDialog payment guide', () => {
  it('shows no payment guide for providers without a flow guide', () => {
    const wrapper = mountDialog()

    expect(wrapper.text()).not.toContain(messages['admin.settings.payment.alipayGuideSummary'])
    expect(wrapper.text()).not.toContain(messages['admin.settings.payment.wxpayGuideSummary'])
    expect(wrapper.find('button[title="View payment guide"]').exists()).toBe(false)
  placeholder)

  it.each([
    ['alipay', 'admin.settings.payment.alipayGuideSummary'],
    ['wxpay', 'admin.settings.payment.wxpayGuideSummary'],
  ])('shows the payment guide summary for %s', async (providerKey, summaryKey) => {
    const wrapper = mountDialog()

    ;(wrapper.vm as unknown as { reset: (key: string) => void placeholder).reset(providerKey)
    await nextTick()

    expect(wrapper.text()).toContain(messages[summaryKey])
    expect(wrapper.find('button[title="View payment guide"]').exists()).toBe(true)
  placeholder)
placeholder)
