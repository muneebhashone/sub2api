import { afterEach, describe, expect, it, vi placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'
import { nextTick placeholder from 'vue'
import BaseDialog from '../BaseDialog.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (key: string) => key placeholder)
placeholder))

describe('BaseDialog', () => {
  afterEach(() => {
    document.body.innerHTML = ''
    document.body.classList.remove('modal-open')
  placeholder)

  it('resets body scroll position when reopened', async () => {
    const wrapper = mount(BaseDialog, {
      attachTo: document.body,
      props: { show: false, title: 'Details' placeholder,
      slots: { default: '<div style="height: 2000px">content</div>' placeholder,
      global: { stubs: { Icon: true placeholder placeholder
    placeholder)

    await wrapper.setProps({ show: true placeholder)
    await nextTick()
    const body = document.body.querySelector<HTMLElement>('.modal-body')
    expect(body).not.toBeNull()
    body!.scrollTop = 480

    await wrapper.setProps({ show: false placeholder)
    await wrapper.setProps({ show: true placeholder)
    await nextTick()

    expect(document.body.querySelector<HTMLElement>('.modal-body')?.scrollTop).toBe(0)
    wrapper.unmount()
  placeholder)
placeholder)
