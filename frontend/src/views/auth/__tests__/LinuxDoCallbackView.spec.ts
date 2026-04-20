import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'

import LinuxDoCallbackView from '../LinuxDoCallbackView.vue'

const replace = vi.fn()
const showSuccess = vi.fn()
const showError = vi.fn()
const setToken = vi.fn()
const exchangePendingOAuthCompletion = vi.fn()
const completeLinuxDoOAuthRegistration = vi.fn()

vi.mock('vue-router', () => ({
  useRoute: () => ({
    query: {placeholder
  placeholder),
  useRouter: () => ({
    replace
  placeholder)
placeholder))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    placeholder)
  placeholder
placeholder)

vi.mock('@/stores', () => ({
  useAuthStore: () => ({
    setToken
  placeholder),
  useAppStore: () => ({
    showSuccess,
    showError
  placeholder)
placeholder))

vi.mock('@/api/auth', async () => {
  const actual = await vi.importActual<typeof import('@/api/auth')>('@/api/auth')
  return {
    ...actual,
    exchangePendingOAuthCompletion: (...args: any[]) => exchangePendingOAuthCompletion(...args),
    completeLinuxDoOAuthRegistration: (...args: any[]) => completeLinuxDoOAuthRegistration(...args)
  placeholder
placeholder)

describe('LinuxDoCallbackView', () => {
  beforeEach(() => {
    replace.mockReset()
    showSuccess.mockReset()
    showError.mockReset()
    setToken.mockReset()
    exchangePendingOAuthCompletion.mockReset()
    completeLinuxDoOAuthRegistration.mockReset()
  placeholder)

  it('does not send adoption decisions during the initial exchange', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      expires_in: 3600,
      redirect: '/dashboard',
      adoption_required: true
    placeholder)
    setToken.mockResolvedValue({placeholder)

    mount(LinuxDoCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(exchangePendingOAuthCompletion).toHaveBeenCalledTimes(1)
    expect(exchangePendingOAuthCompletion).toHaveBeenCalledWith()
  placeholder)

  it('waits for explicit adoption confirmation before finishing a non-invitation login', async () => {
    exchangePendingOAuthCompletion
      .mockResolvedValueOnce({
        redirect: '/dashboard',
        adoption_required: true,
        suggested_display_name: 'LinuxDo Nick',
        suggested_avatar_url: 'https://cdn.example/linuxdo.png'
      placeholder)
      .mockResolvedValueOnce({
        access_token: 'access-token',
        refresh_token: 'refresh-token',
        expires_in: 3600,
        redirect: '/dashboard'
      placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(LinuxDoCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(wrapper.text()).toContain('LinuxDo Nick')
    expect(setToken).not.toHaveBeenCalled()
    expect(replace).not.toHaveBeenCalled()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    await checkboxes[1].setValue(false)

    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(1)
    await buttons[0].trigger('click')
    await flushPromises()

    expect(exchangePendingOAuthCompletion).toHaveBeenCalledTimes(2)
    expect(exchangePendingOAuthCompletion).toHaveBeenNthCalledWith(1)
    expect(exchangePendingOAuthCompletion).toHaveBeenNthCalledWith(2, {
      adoptDisplayName: true,
      adoptAvatar: false
    placeholder)
    expect(setToken).toHaveBeenCalledWith('access-token')
    expect(replace).toHaveBeenCalledWith('/dashboard')
  placeholder)

  it('treats a completion without token as bind success and returns to profile', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({placeholder)

    mount(LinuxDoCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(setToken).not.toHaveBeenCalled()
    expect(showSuccess).toHaveBeenCalledWith('profile.authBindings.bindSuccess')
    expect(replace).toHaveBeenCalledWith('/profile')
  placeholder)

  it('supports bind completion after adoption confirmation', async () => {
    exchangePendingOAuthCompletion
      .mockResolvedValueOnce({
        redirect: '/dashboard',
        adoption_required: true,
        suggested_display_name: 'LinuxDo Nick',
        suggested_avatar_url: 'https://cdn.example/linuxdo.png'
      placeholder)
      .mockResolvedValueOnce({
        redirect: '/profile/security'
      placeholder)

    const wrapper = mount(LinuxDoCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()

    expect(exchangePendingOAuthCompletion).toHaveBeenNthCalledWith(2, {
      adoptDisplayName: true,
      adoptAvatar: true
    placeholder)
    expect(setToken).not.toHaveBeenCalled()
    expect(showSuccess).toHaveBeenCalledWith('profile.authBindings.bindSuccess')
    expect(replace).toHaveBeenCalledWith('/profile/security')
  placeholder)

  it('renders adoption choices for invitation flow and submits the selected values', async () => {
    exchangePendingOAuthCompletion.mockResolvedValue({
      error: 'invitation_required',
      redirect: '/dashboard',
      adoption_required: true,
      suggested_display_name: 'LinuxDo Nick',
      suggested_avatar_url: 'https://cdn.example/linuxdo.png'
    placeholder)
    completeLinuxDoOAuthRegistration.mockResolvedValue({
      access_token: 'access-token',
      refresh_token: 'refresh-token',
      expires_in: 3600,
      token_type: 'Bearer'
    placeholder)
    setToken.mockResolvedValue({placeholder)

    const wrapper = mount(LinuxDoCallbackView, {
      global: {
        stubs: {
          AuthLayout: { template: '<div><slot /></div>' placeholder,
          Icon: true,
          RouterLink: { template: '<a><slot /></a>' placeholder,
          transition: false
        placeholder
      placeholder
    placeholder)

    await flushPromises()

    expect(wrapper.text()).toContain('LinuxDo Nick')
    expect(exchangePendingOAuthCompletion).toHaveBeenCalledTimes(1)
    expect(exchangePendingOAuthCompletion).toHaveBeenCalledWith()

    const checkboxes = wrapper.findAll('input[type="checkbox"]')
    expect(checkboxes).toHaveLength(2)

    await checkboxes[0].setValue(false)
    await wrapper.find('input[type="text"]').setValue('invite-code')
    await wrapper.find('button').trigger('click')

    expect(completeLinuxDoOAuthRegistration).toHaveBeenCalledWith('invite-code', {
      adoptDisplayName: false,
      adoptAvatar: true
    placeholder)
  placeholder)
placeholder)
