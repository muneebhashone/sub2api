import { describe, expect, it, vi, beforeEach, afterEach placeholder from 'vitest'
import { flushPromises, mount placeholder from '@vue/test-utils'
import { defineComponent placeholder from 'vue'
import AccountTestModal from '../AccountTestModal.vue'

const { getAvailableModelsMock placeholder = vi.hoisted(() => ({
  getAvailableModelsMock: vi.fn()
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      getAvailableModels: getAvailableModelsMock
    placeholder
  placeholder
placeholder))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: vi.fn()
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

const BaseDialogStub = defineComponent({
  name: 'BaseDialog',
  props: { show: { type: Boolean, default: false placeholder placeholder,
  template: '<div v-if="show"><slot /><slot name="footer" /></div>'
placeholder)

const SelectStub = defineComponent({
  name: 'SelectStub',
  props: {
    modelValue: { type: [String, Number, Boolean, null], default: '' placeholder,
    options: { type: Array, default: () => [] placeholder,
    valueKey: { type: String, default: 'value' placeholder,
    labelKey: { type: String, default: 'label' placeholder
  placeholder,
  emits: ['update:modelValue'],
  template: `
    <select
      v-bind="$attrs"
      :value="modelValue"
      @change="$emit('update:modelValue', $event.target.value)"
    >
      <option
        v-for="option in options"
        :key="option[valueKey]"
        :value="option[valueKey]"
      >
        {{ option[labelKey] placeholderplaceholder
      </option>
    </select>
  `
placeholder)

const TextAreaStub = defineComponent({
  name: 'TextArea',
  props: {
    modelValue: { type: String, default: '' placeholder
  placeholder,
  emits: ['update:modelValue'],
  template: `
    <textarea
      v-bind="$attrs"
      :value="modelValue"
      @input="$emit('update:modelValue', $event.target.value)"
    />
  `
placeholder)

function buildAccount() {
  return {
    id: 1,
    name: 'OpenAI OAuth',
    platform: 'openai',
    type: 'oauth',
    status: 'active',
    credentials: {placeholder,
    extra: {placeholder,
    concurrency: 1,
    priority: 1,
    proxy_id: null,
    auto_pause_on_expired: false
  placeholder as any
placeholder

describe('AccountTestModal', () => {
  const originalFetch = global.fetch

  beforeEach(() => {
    getAvailableModelsMock.mockReset()
    getAvailableModelsMock.mockResolvedValue([
      { id: 'gpt-5.4', display_name: 'GPT-5.4' placeholder
    ])
    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      body: {
        getReader: () => ({
          read: vi.fn().mockResolvedValue({ done: true, value: undefined placeholder)
        placeholder)
      placeholder
    placeholder as any)
    localStorage.setItem('auth_token', 'test-token')
  placeholder)

  afterEach(() => {
    global.fetch = originalFetch
    localStorage.clear()
  placeholder)

  it('posts compact mode for OpenAI compact probe', async () => {
    const wrapper = mount(AccountTestModal, {
      props: {
        show: true,
        account: buildAccount()
      placeholder,
      global: {
        stubs: {
          BaseDialog: BaseDialogStub,
          Select: SelectStub,
          TextArea: TextAreaStub,
          Icon: true
        placeholder
      placeholder
    placeholder)

    await flushPromises()
    ;(wrapper.vm as any).selectedModelId = 'gpt-5.4'
    ;(wrapper.vm as any).testMode = 'compact'
    await (wrapper.vm as any).startTest()
    await flushPromises()

    expect(global.fetch).toHaveBeenCalledTimes(1)
    const [, options] = (global.fetch as any).mock.calls[0]
    expect(JSON.parse(options.body)).toMatchObject({
      model_id: 'gpt-5.4',
      mode: 'compact'
    placeholder)
  placeholder)
placeholder)
