import { flushPromises, mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import TencentCaptchaGate from '@/components/TencentCaptchaGate.vue'
import { resetTencentCaptchaLoaderForTest placeholder from '@/utils/tencentCaptcha'

const locale = { value: 'zh' placeholder

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ locale placeholder)
placeholder))

type CaptchaResult = {
  ret: number
  ticket?: string | null
  randstr?: string | null
  errorCode?: number
placeholder

describe('TencentCaptchaGate', () => {
  beforeEach(() => {
    locale.value = 'zh'
    delete window.TencentCaptcha
    document.head.querySelectorAll('script[src*="TJCaptcha.js"]').forEach((node) => node.remove())
    resetTencentCaptchaLoaderForTest()
  placeholder)

  it('does not render a visible verification button', () => {
    const wrapper = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)

    expect(wrapper.find('button').exists()).toBe(false)
  placeholder)

  it('resolves proof after Tencent SDK success', async () => {
    let callback: ((result: CaptchaResult) => void) | undefined
    window.TencentCaptcha = class {
      constructor(_appId: string, resultCallback: (result: CaptchaResult) => void) {
        callback = resultCallback
      placeholder
      show = vi.fn()
      destroy = vi.fn()
    placeholder
    const wrapper = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)

    const verification = wrapper.vm.verify()
    await flushPromises()
    callback?.({ ret: 0, ticket: 'ticket-value', randstr: 'rand-value' placeholder)

    await expect(verification).resolves.toEqual({ ticket: 'ticket-value', randstr: 'rand-value' placeholder)
  placeholder)

  it('resolves null when the user closes the popup', async () => {
    let callback: ((result: CaptchaResult) => void) | undefined
    window.TencentCaptcha = class {
      constructor(_appId: string, resultCallback: (result: CaptchaResult) => void) {
        callback = resultCallback
      placeholder
      show = vi.fn()
      destroy = vi.fn()
    placeholder
    const wrapper = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)

    const verification = wrapper.vm.verify()
    await flushPromises()
    callback?.({ ret: 2, ticket: null placeholder)

    await expect(verification).resolves.toBeNull()
  placeholder)

  it('rejects SDK load failures and disaster-recovery tickets', async () => {
    const failedLoad = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)
    const loadVerification = failedLoad.vm.verify()
    const script = document.head.querySelector<HTMLScriptElement>('script[src*="TJCaptcha.js"]')
    expect(script).not.toBeNull()
    script?.dispatchEvent(new Event('error'))
    await expect(loadVerification).rejects.toThrow('Failed to load Tencent Captcha SDK')

    let callback: ((result: CaptchaResult) => void) | undefined
    window.TencentCaptcha = class {
      constructor(_appId: string, resultCallback: (result: CaptchaResult) => void) {
        callback = resultCallback
      placeholder
      show = vi.fn()
      destroy = vi.fn()
    placeholder
    const failedResult = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)
    const resultVerification = failedResult.vm.verify()
    await flushPromises()
    callback?.({ ret: 0, ticket: 'trerror_1001_123456789', randstr: '@fallback', errorCode: 1001 placeholder)

    await expect(resultVerification).rejects.toThrow('Tencent Captcha verification failed')
  placeholder)

  it('reuses one pending promise for concurrent verify calls', async () => {
    const show = vi.fn()
    let callback: ((result: CaptchaResult) => void) | undefined
    window.TencentCaptcha = class {
      constructor(_appId: string, resultCallback: (result: CaptchaResult) => void) {
        callback = resultCallback
      placeholder
      show = show
      destroy = vi.fn()
    placeholder
    const wrapper = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)

    const first = wrapper.vm.verify()
    const second = wrapper.vm.verify()
    await flushPromises()
    callback?.({ ret: 0, ticket: 'ticket-value', randstr: 'rand-value' placeholder)

    await expect(first).resolves.toEqual({ ticket: 'ticket-value', randstr: 'rand-value' placeholder)
    await expect(second).resolves.toEqual({ ticket: 'ticket-value', randstr: 'rand-value' placeholder)
    expect(show).toHaveBeenCalledOnce()
  placeholder)

  it('settles a pending verification when reset', async () => {
    const destroy = vi.fn()
    window.TencentCaptcha = class {
      constructor(_appId: string, _callback: (result: CaptchaResult) => void) {placeholder
      show = vi.fn()
      destroy = destroy
    placeholder
    const wrapper = mount(TencentCaptchaGate, { props: { appId: '123456789' placeholder placeholder)

    const verification = wrapper.vm.verify()
    await flushPromises()
    wrapper.vm.reset()

    await expect(verification).resolves.toBeNull()
    expect(destroy).toHaveBeenCalledOnce()
  placeholder)
placeholder)
