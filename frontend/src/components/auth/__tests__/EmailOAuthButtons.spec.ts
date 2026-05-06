import { mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import EmailOAuthButtons from '@/components/auth/EmailOAuthButtons.vue'

const routeState = vi.hoisted(() => ({
  query: {placeholder as Record<string, unknown>,
placeholder))

const locationState = vi.hoisted(() => ({
  current: { href: 'http://localhost/register?aff=AFF123' placeholder as { href: string placeholder,
placeholder))

vi.mock('vue-router', () => ({
  useRoute: () => routeState,
placeholder))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string, params?: Record<string, string>) => {
      if (key === 'auth.emailOAuth.signIn') {
        return `使用 ${params?.providerName ?? ''placeholder 登录`
      placeholder
      return key
    placeholder,
  placeholder),
placeholder))

describe('EmailOAuthButtons', () => {
  beforeEach(() => {
    routeState.query = { redirect: '/billing?plan=pro', aff: 'AFF123' placeholder
    locationState.current = { href: 'http://localhost/register?aff=AFF123' placeholder
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: locationState.current,
    placeholder)
    window.localStorage.clear()
    window.sessionStorage.clear()
  placeholder)

  it('passes the affiliate code to the email oauth start URL', async () => {
    const wrapper = mount(EmailOAuthButtons, {
      props: {
        githubEnabled: true,
        googleEnabled: false,
      placeholder,
      global: {
        stubs: {
          GitHubMark: true,
          GoogleMark: true,
        placeholder,
      placeholder,
    placeholder)

    await wrapper.get('button').trigger('click')

    expect(locationState.current.href).toBe(
      '/api/v1/auth/oauth/github/start?redirect=%2Fbilling%3Fplan%3Dpro&aff_code=AFF123'
    )
    expect(window.sessionStorage.getItem('oauth_aff_code')).toBe('AFF123')
  placeholder)

  it('uses a full-width descriptive button when only GitHub is enabled', () => {
    const wrapper = mount(EmailOAuthButtons, {
      props: {
        githubEnabled: true,
        googleEnabled: false,
      placeholder,
      global: {
        stubs: {
          GitHubMark: true,
          GoogleMark: true,
        placeholder,
      placeholder,
    placeholder)

    expect(wrapper.find('.grid').classes()).not.toContain('sm:grid-cols-2')
    expect(wrapper.get('button').text()).toContain('使用 GitHub 登录')
  placeholder)

  it('uses compact labels and two columns when GitHub and Google are both enabled', () => {
    const wrapper = mount(EmailOAuthButtons, {
      props: {
        githubEnabled: true,
        googleEnabled: true,
      placeholder,
      global: {
        stubs: {
          GitHubMark: true,
          GoogleMark: true,
        placeholder,
      placeholder,
    placeholder)

    expect(wrapper.find('.grid').classes()).toContain('sm:grid-cols-2')
    const buttons = wrapper.findAll('button')
    expect(buttons).toHaveLength(2)
    expect(buttons[0].text()).toContain('GitHub')
    expect(buttons[0].text()).not.toContain('使用 GitHub 登录')
    expect(buttons[1].text()).toContain('Google')
    expect(buttons[1].text()).not.toContain('使用 Google 登录')
  placeholder)
placeholder)
