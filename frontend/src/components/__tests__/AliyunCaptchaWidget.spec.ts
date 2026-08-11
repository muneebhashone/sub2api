import { describe, it, expect, vi, beforeEach, afterEach placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'
import AliyunCaptchaWidget from '../AliyunCaptchaWidget.vue'

interface CapturedInitOptions {
  SceneId: string
  prefix: string
  mode: string
  element: string
  button: string
  captchaVerifyCallback: (param: string) => { captchaResult: boolean placeholder
  language?: string
placeholder

const i18nStub = {
  install(app: { config: { globalProperties: Record<string, unknown> placeholder placeholder) {
    app.config.globalProperties.$t = (key: string) => key
  placeholder
placeholder

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ locale: { value: 'zh-CN' placeholder, t: (key: string) => key placeholder)
placeholder))

describe('AliyunCaptchaWidget', () => {
  let initOptions: CapturedInitOptions | null

  beforeEach(() => {
    initOptions = null
    window.initAliyunCaptcha = vi.fn((options: CapturedInitOptions) => {
      initOptions = options
    placeholder) as unknown as typeof window.initAliyunCaptcha
  placeholder)

  afterEach(() => {
    vi.useRealTimers()
    delete window.initAliyunCaptcha
    delete window.AliyunCaptchaConfig
    document.getElementById('aliyunCaptcha-window-popup')?.remove()
    document.getElementById('aliyunCaptcha-mask')?.remove()
  placeholder)

  function mountWidget() {
    return mount(AliyunCaptchaWidget, {
      props: { sceneId: 'scene-1', prefix: 'prefix-1', region: 'cn' as const placeholder,
      attachTo: document.body,
      global: { plugins: [i18nStub] placeholder
    placeholder)
  placeholder

  function createVisiblePopup(): HTMLElement {
    const popup = document.createElement('div')
    popup.id = 'aliyunCaptcha-window-popup'
    popup.style.display = 'block'
    document.body.appendChild(popup)
    return popup
  placeholder

  it('渲染可见验证按钮并以 popup 模式初始化，全局配置就位', async () => {
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    const button = wrapper.get('button')
    expect(button.text()).toContain('auth.captchaClickToVerify')
    expect(window.AliyunCaptchaConfig).toEqual({ region: 'cn', prefix: 'prefix-1' placeholder)
    expect(initOptions).not.toBeNull()
    expect(initOptions!.mode).toBe('popup')
    expect(initOptions!.SceneId).toBe('scene-1')
    expect(initOptions!.language).toBe('cn')

    wrapper.unmount()
  placeholder)

  it('用户点击按钮进入验证中，验证完成后 emit verify 并置已通过', async () => {
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    await wrapper.get('button').trigger('click')
    expect(wrapper.get('button').text()).toContain('auth.captchaVerifying')

    const result = initOptions!.captchaVerifyCallback('captcha-param-1')
    expect(result).toEqual({ captchaResult: true placeholder)
    await wrapper.vm.$nextTick()

    expect(wrapper.emitted('verify')).toEqual([['captcha-param-1']])
    expect(wrapper.get('button').text()).toContain('auth.captchaVerified')
    expect(wrapper.get('button').attributes('disabled')).toBeDefined()

    wrapper.unmount()
  placeholder)

  it('verify() 在已预验证时直接复用缓存的 captchaVerifyParam', async () => {
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    await wrapper.get('button').trigger('click')
    initOptions!.captchaVerifyCallback('captcha-param-2')

    const vm = wrapper.vm as unknown as { verify: () => Promise<string | null> placeholder
    await expect(vm.verify()).resolves.toBe('captcha-param-2')

    wrapper.unmount()
  placeholder)

  it('verify() 在未预验证时触发弹窗流程并等待结果', async () => {
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    const vm = wrapper.vm as unknown as { verify: () => Promise<string | null> placeholder
    const pending = vm.verify()
    await Promise.resolve()
    await Promise.resolve()
    await wrapper.vm.$nextTick()
    expect(wrapper.get('button').text()).toContain('auth.captchaVerifying')

    initOptions!.captchaVerifyCallback('captcha-param-3')
    await expect(pending).resolves.toBe('captcha-param-3')

    wrapper.unmount()
  placeholder)

  it('弹窗未出现前会按 tick 重试触发按钮（SDK 异步绑定兜底）', async () => {
    vi.useFakeTimers()
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    const vm = wrapper.vm as unknown as { verify: () => Promise<string | null> placeholder
    void vm.verify()
    await Promise.resolve()
    await Promise.resolve()

    const button = wrapper.get('button').element as HTMLButtonElement
    const clickSpy = vi.fn()
    button.addEventListener('click', clickSpy)

    await vi.advanceTimersByTimeAsync(1000)
    expect(clickSpy.mock.calls.length).toBeGreaterThanOrEqual(3)

    button.removeEventListener('click', clickSpy)
    wrapper.unmount()
  placeholder)

  it('弹窗出现后被用户关闭时 resolve null 并回到未验证态', async () => {
    vi.useFakeTimers()
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    const vm = wrapper.vm as unknown as { verify: () => Promise<string | null> placeholder
    const pending = vm.verify()
    await Promise.resolve()
    await Promise.resolve()

    const popup = createVisiblePopup()
    await vi.advanceTimersByTimeAsync(400)
    popup.remove()
    await vi.advanceTimersByTimeAsync(400)

    await expect(pending).resolves.toBeNull()
    expect(wrapper.get('button').text()).toContain('auth.captchaClickToVerify')

    wrapper.unmount()
  placeholder)

  it('reset() 清空缓存并取消进行中的验证', async () => {
    const wrapper = mountWidget()
    await Promise.resolve()
    await Promise.resolve()

    await wrapper.get('button').trigger('click')
    initOptions!.captchaVerifyCallback('captcha-param-4')

    const vm = wrapper.vm as unknown as {
      verify: () => Promise<string | null>
      reset: () => void
    placeholder
    vm.reset()
    await wrapper.vm.$nextTick()
    expect(wrapper.get('button').text()).toContain('auth.captchaClickToVerify')

    // reset 后缓存失效，verify() 重新走弹窗流程
    const pending = vm.verify()
    await Promise.resolve()
    initOptions!.captchaVerifyCallback('captcha-param-5')
    await expect(pending).resolves.toBe('captcha-param-5')

    wrapper.unmount()
  placeholder)
placeholder)

declare global {
  interface Window {
    initAliyunCaptcha?: (options: unknown) => void
    AliyunCaptchaConfig?: { region: string; prefix: string placeholder
  placeholder
placeholder
