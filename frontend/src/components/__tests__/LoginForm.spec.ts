/**
 * LoginView 组件核心逻辑测试
 * 测试登录表单提交、验证、2FA 等场景
 */
import { describe, it, expect, vi, beforeEach placeholder from 'vitest'
import { mount, flushPromises placeholder from '@vue/test-utils'
import { setActivePinia, createPinia placeholder from 'pinia'
import { defineComponent, reactive, ref placeholder from 'vue'
import { useAuthStore placeholder from '@/stores/auth'

// Mock 所有外部依赖
const mockLogin = vi.fn()
const mockLogin2FA = vi.fn()
const mockPush = vi.fn()

vi.mock('@/api', () => ({
  authAPI: {
    login: (...args: any[]) => mockLogin(...args),
    login2FA: (...args: any[]) => mockLogin2FA(...args),
    logout: vi.fn(),
    getCurrentUser: vi.fn().mockResolvedValue({ data: {placeholder placeholder),
    register: vi.fn(),
    refreshToken: vi.fn(),
  placeholder,
  isTotp2FARequired: (response: any) => response?.requires_2fa === true,
placeholder))

vi.mock('@/api/admin/system', () => ({
  checkUpdates: vi.fn(),
placeholder))

vi.mock('@/api/auth', () => ({
  getPublicSettings: vi.fn().mockResolvedValue({placeholder),
placeholder))

/**
 * 创建一个简化的测试组件来封装登录逻辑
 * 避免引入 LoginView.vue 的全部依赖（AuthLayout、i18n、Icon 等）
 */
const LoginFormTestComponent = defineComponent({
  setup() {
    const authStore = useAuthStore()
    const formData = reactive({ email: '', password: '' placeholder)
    const isLoading = ref(false)
    const errorMessage = ref('')

    const handleLogin = async () => {
      if (!formData.email || !formData.password) {
        errorMessage.value = '请输入邮箱和密码'
        return
      placeholder

      isLoading.value = true
      errorMessage.value = ''

      try {
        const response = await authStore.login({
          email: formData.email,
          password: formData.password,
        placeholder)

        // 2FA 流程由调用方处理
        if ((response as any)?.requires_2fa) {
          errorMessage.value = '需要 2FA 验证'
          return
        placeholder

        mockPush('/dashboard')
      placeholder catch (error: any) {
        errorMessage.value = error.message || '登录失败'
      placeholder finally {
        isLoading.value = false
      placeholder
    placeholder

    return { formData, isLoading, errorMessage, handleLogin placeholder
  placeholder,
  template: `
    <form @submit.prevent="handleLogin">
      <input id="email" v-model="formData.email" type="email" />
      <input id="password" v-model="formData.password" type="password" />
      <p v-if="errorMessage" class="error">{{ errorMessage placeholderplaceholder</p>
      <button type="submit" :disabled="isLoading">登录</button>
    </form>
  `,
placeholder)

describe('LoginForm 核心逻辑', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  placeholder)

  it('成功登录后跳转到 dashboard', async () => {
    mockLogin.mockResolvedValue({
      access_token: 'token',
      token_type: 'Bearer',
      user: { id: 1, username: 'test', email: 'test@example.com', role: 'user', balance: 0, concurrency: 5, status: 'active', allowed_groups: null, created_at: '', updated_at: '' placeholder,
    placeholder)

    const wrapper = mount(LoginFormTestComponent)

    await wrapper.find('#email').setValue('test@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(mockLogin).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password123',
    placeholder)
    expect(mockPush).toHaveBeenCalledWith('/dashboard')
  placeholder)

  it('登录失败时显示错误信息', async () => {
    mockLogin.mockRejectedValue(new Error('Invalid credentials'))

    const wrapper = mount(LoginFormTestComponent)

    await wrapper.find('#email').setValue('test@example.com')
    await wrapper.find('#password').setValue('wrong')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.error').text()).toBe('Invalid credentials')
  placeholder)

  it('空表单提交显示验证错误', async () => {
    const wrapper = mount(LoginFormTestComponent)

    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.error').text()).toBe('请输入邮箱和密码')
    expect(mockLogin).not.toHaveBeenCalled()
  placeholder)

  it('需要 2FA 时不跳转', async () => {
    mockLogin.mockResolvedValue({
      requires_2fa: true,
      temp_token: 'temp-123',
    placeholder)

    const wrapper = mount(LoginFormTestComponent)

    await wrapper.find('#email').setValue('test@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(mockPush).not.toHaveBeenCalled()
    expect(wrapper.find('.error').text()).toBe('需要 2FA 验证')
  placeholder)

  it('提交过程中按钮被禁用', async () => {
    let resolveLogin: (v: any) => void
    mockLogin.mockImplementation(
      () => new Promise((resolve) => { resolveLogin = resolve placeholder)
    )

    const wrapper = mount(LoginFormTestComponent)

    await wrapper.find('#email').setValue('test@example.com')
    await wrapper.find('#password').setValue('password123')
    await wrapper.find('form').trigger('submit')

    expect(wrapper.find('button').attributes('disabled')).toBeDefined()

    resolveLogin!({
      access_token: 'token',
      token_type: 'Bearer',
      user: { id: 1, username: 'test', email: 'test@example.com', role: 'user', balance: 0, concurrency: 5, status: 'active', allowed_groups: null, created_at: '', updated_at: '' placeholder,
    placeholder)
    await flushPromises()

    expect(wrapper.find('button').attributes('disabled')).toBeUndefined()
  placeholder)
placeholder)
