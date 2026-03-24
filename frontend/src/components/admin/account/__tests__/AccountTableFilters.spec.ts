import { describe, expect, it, vi placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'

import AccountTableFilters from '../AccountTableFilters.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    placeholder)
  placeholder
placeholder)

describe('AccountTableFilters', () => {
  it('renders privacy mode options and emits privacy_mode updates', async () => {
    const wrapper = mount(AccountTableFilters, {
      props: {
        searchQuery: '',
        filters: {
          platform: '',
          type: '',
          status: '',
          group: '',
          privacy_mode: ''
        placeholder,
        groups: []
      placeholder,
      global: {
        stubs: {
          SearchInput: {
            template: '<div />'
          placeholder,
          Select: {
            props: ['modelValue', 'options'],
            emits: ['update:modelValue', 'change'],
            template: '<div class="select-stub" :data-options="JSON.stringify(options)" />'
          placeholder
        placeholder
      placeholder
    placeholder)

    const selects = wrapper.findAll('.select-stub')
    expect(selects).toHaveLength(5)

    const privacyOptions = JSON.parse(selects[3].attributes('data-options'))
    expect(privacyOptions).toEqual([
      { value: '', label: 'admin.accounts.allPrivacyModes' placeholder,
      { value: '__unset__', label: 'admin.accounts.privacyUnset' placeholder,
      { value: 'training_off', label: 'Privacy' placeholder,
      { value: 'training_set_cf_blocked', label: 'CF' placeholder,
      { value: 'training_set_failed', label: 'Fail' placeholder
    ])
  placeholder)
placeholder)
