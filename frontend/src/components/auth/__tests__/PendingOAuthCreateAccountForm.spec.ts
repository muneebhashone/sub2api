import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import PendingOAuthCreateAccountForm from '../PendingOAuthCreateAccountForm.vue'

const sendVerifyCode = vi.fn()
const getPublicSettings = vi.fn()

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    placeholder)
  placeholder
placeholder)

vi.mock('@/api/auth', async () => {
  const actual = await vi.importActual<typeof import('@/api/auth')>('@/api/auth')
  return {
    ...actual,
    sendVerifyCode: (...args: any[]) => sendVerifyCode(...args),
    getPublicSettings: (...args: any[]) => getPublicSettings(...args)
  placeholder
placeholder)

describe('PendingOAuthCreateAccountForm', () => {
  beforeEach(() => {
    sendVerifyCode.mockReset()
    getPublicSettings.mockReset()
    getPublicSettings.mockResolvedValue({
      turnstile_enabled: false,
      turnstile_site_key: ''
    placeholder)
  placeholder)

  it('emits trimmed email, password, and verify code on submit', async () => {
    const wrapper = mount(PendingOAuthCreateAccountForm, {
      props: {
        providerName: 'LinuxDo',
        testIdPrefix: 'linuxdo',
        initialEmail: 'prefill@example.com',
        isSubmitting: false
      placeholder
    placeholder)

    await wrapper.get('[data-testid="linuxdo-create-account-email"]').setValue('  user@example.com  ')
    await wrapper.get('[data-testid="linuxdo-create-account-password"]').setValue('secret-123')
    await wrapper.get('[data-testid="linuxdo-create-account-verify-code"]').setValue(' 246810 ')
    await wrapper.get('form').trigger('submit.prevent')

    expect(wrapper.emitted('submit')).toEqual([
      [
        {
          email: 'user@example.com',
          password: 'secret-123',
          verifyCode: '246810'
        placeholder
      ]
    ])
  placeholder)

  it('sends a verify code for the trimmed email value', async () => {
    sendVerifyCode.mockResolvedValue({
      message: 'sent',
      countdown: 60
    placeholder)

    const wrapper = mount(PendingOAuthCreateAccountForm, {
      props: {
        providerName: 'LinuxDo',
        testIdPrefix: 'linuxdo',
        initialEmail: '',
        isSubmitting: false
      placeholder
    placeholder)

    await wrapper.get('[data-testid="linuxdo-create-account-email"]').setValue('  user@example.com  ')
    await wrapper.get('[data-testid="linuxdo-create-account-send-code"]').trigger('click')
    await flushPromises()

    expect(sendVerifyCode).toHaveBeenCalledWith({
      email: 'user@example.com'
    placeholder)
  placeholder)

  it('requires a turnstile token before sending a verify code when turnstile is enabled', async () => {
    getPublicSettings.mockResolvedValue({
      turnstile_enabled: true,
      turnstile_site_key: 'site-key'
    placeholder)
    sendVerifyCode.mockResolvedValue({
      message: 'sent',
      countdown: 60
    placeholder)

    const wrapper = mount(PendingOAuthCreateAccountForm, {
      props: {
        providerName: 'LinuxDo',
        testIdPrefix: 'linuxdo',
        initialEmail: '',
        isSubmitting: false
      placeholder,
      global: {
        stubs: {
          TurnstileWidget: {
            template: '<button data-testid="turnstile-verify" @click="$emit(\'verify\', \'turnstile-token\')">verify</button>'
          placeholder
        placeholder
      placeholder
    placeholder)

    await flushPromises()
    await wrapper.get('[data-testid="linuxdo-create-account-email"]').setValue('  user@example.com  ')

    expect(wrapper.get('[data-testid="linuxdo-create-account-send-code"]').attributes('disabled')).toBeDefined()

    await wrapper.get('[data-testid="turnstile-verify"]').trigger('click')
    await wrapper.get('[data-testid="linuxdo-create-account-send-code"]').trigger('click')
    await flushPromises()

    expect(sendVerifyCode).toHaveBeenCalledWith({
      email: 'user@example.com',
      turnstile_token: 'turnstile-token'
    placeholder)
  placeholder)
placeholder)
