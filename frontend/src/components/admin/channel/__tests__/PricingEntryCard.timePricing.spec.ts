import { shallowMount placeholder from '@vue/test-utils'
import { describe, expect, it, vi placeholder from 'vitest'
import PricingEntryCard from '../PricingEntryCard.vue'
import type { PricingFormEntry placeholder from '../types'

vi.mock('vue-i18n', async importOriginal => ({
  ...await importOriginal<typeof import('vue-i18n')>(),
  useI18n: () => ({ t: (key: string) => key placeholder),
placeholder))

function createEntry(billingMode: PricingFormEntry['billing_mode'] = 'token'): PricingFormEntry {
  return {
    models: [],
    billing_mode: billingMode,
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    fast_multiplier: null,
    flex_multiplier: null,
    image_input_price: null,
    image_output_price: null,
    per_request_price: null,
    intervals: [],
    time_pricing: {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00', end_time: '12:00', multiplier: '2.00' placeholder],
    placeholder,
  placeholder
placeholder

describe('PricingEntryCard time pricing visibility', () => {
  it('is hidden by default', () => {
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry: createEntry() placeholder,
    placeholder)

    expect(wrapper.findComponent({ name: 'TimePricingSection' placeholder).exists()).toBe(false)
  placeholder)

  it('is shown for token pricing when explicitly enabled', () => {
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry: createEntry(), enableTimePricing: true placeholder,
    placeholder)

    expect(wrapper.findComponent({ name: 'TimePricingSection' placeholder).exists()).toBe(true)
  placeholder)

  it('is hidden for non-token pricing even when explicitly enabled', () => {
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry: createEntry('per_request'), enableTimePricing: true placeholder,
    placeholder)

    expect(wrapper.findComponent({ name: 'TimePricingSection' placeholder).exists()).toBe(false)
  placeholder)

  it('clears time periods when changing billing mode', () => {
    const entry = createEntry()
    const wrapper = shallowMount(PricingEntryCard, {
      props: { entry, enableTimePricing: true placeholder,
    placeholder)

    wrapper.findComponent({ name: 'Select' placeholder).vm.$emit('update:modelValue', 'image')

    expect(wrapper.emitted('update')?.[0]?.[0]).toEqual({
      ...entry,
      billing_mode: 'image',
      intervals: [],
      time_pricing: { timezone: 'Asia/Shanghai', periods: [] placeholder,
    placeholder)
    expect(entry.time_pricing.periods).toHaveLength(1)
  placeholder)
placeholder)

describe('PricingEntryCard service tier multipliers', () => {
  it('shows Fast and Flex controls only when explicitly enabled', () => {
    const hidden = shallowMount(PricingEntryCard, { props: { entry: createEntry() placeholder placeholder)
    expect(hidden.text()).not.toContain('admin.channels.form.fastMultiplier')

    const shown = shallowMount(PricingEntryCard, {
      props: { entry: createEntry(), enableTierMultipliers: true placeholder,
    placeholder)
    expect(shown.text()).toContain('admin.channels.form.fastMultiplier')
    expect(shown.text()).toContain('admin.channels.form.flexMultiplier')
  placeholder)
placeholder)
