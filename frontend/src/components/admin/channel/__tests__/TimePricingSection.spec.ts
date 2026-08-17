import { mount placeholder from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import TimePricingSection from '../TimePricingSection.vue'
import { createDefaultTimePricingForm placeholder from '../types'

vi.mock('vue-i18n', async importOriginal => ({
  ...await importOriginal<typeof import('vue-i18n')>(),
  useI18n: () => ({ t: (key: string) => key placeholder),
placeholder))

const SelectStub = {
  name: 'Select',
  props: ['modelValue', 'options'],
  emits: ['update:modelValue'],
  template: '<button type="button" data-testid="timezone-select" />',
placeholder

describe('TimePricingSection', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  placeholder)

  it('adds a neutral period immutably', async () => {
    const value = createDefaultTimePricingForm()
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    await wrapper.get('[data-testid="add-time-period"]').trigger('click')

    expect(wrapper.emitted('update:modelValue')?.[0]?.[0]).toEqual({
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '', end_time: '', multiplier: '1.00' placeholder],
    placeholder)
    expect(value.periods).toHaveLength(0)
  placeholder)

  it('removes without mutating props', async () => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    const removeButton = wrapper.get('[data-testid="remove-time-period-0"]')
    expect(removeButton.attributes('title')).toBe('admin.channels.form.removeTimePeriod')
    expect(removeButton.attributes('aria-label')).toBe('admin.channels.form.removeTimePeriod')
    await removeButton.trigger('click')

    expect(wrapper.emitted('update:modelValue')?.[0]?.[0]).toEqual({
      timezone: 'Asia/Shanghai',
      periods: [],
    placeholder)
    expect(value.periods).toHaveLength(1)
  placeholder)

  it('uses locale-independent 24-hour second inputs', () => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    const timeInputs = wrapper.findAll('input[inputmode="numeric"]')
    expect(timeInputs).toHaveLength(2)
    expect(timeInputs.every(input => input.attributes('type') === 'text')).toBe(true)
    expect(timeInputs.every(input => input.attributes('maxlength') === '8')).toBe(true)
    expect(timeInputs.every(input => input.attributes('placeholder') === 'HH:mm:ss')).toBe(true)
  placeholder)

  it('updates the timezone and normalizes full-width colons immutably', async () => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    wrapper.getComponent(SelectStub).vm.$emit('update:modelValue', 'Europe/London')
    await wrapper.get('input').setValue('10：30：15')

    expect(wrapper.emitted('update:modelValue')?.[0]?.[0]).toEqual({
      ...value,
      timezone: 'Europe/London',
    placeholder)
    expect(wrapper.emitted('update:modelValue')?.[1]?.[0]).toEqual({
      ...value,
      periods: [{ ...value.periods[0], start_time: '10:30:15' placeholder],
    placeholder)
    expect(value).toEqual({
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder)
  placeholder)

  it('normalizes 24:00:00 to 00:00:00', async () => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    await wrapper.get('input').setValue('24:00:00')

    expect(wrapper.emitted('update:modelValue')?.[0]?.[0]).toEqual({
      ...value,
      periods: [{ ...value.periods[0], start_time: '00:00:00' placeholder],
    placeholder)
  placeholder)

  it('formats only a valid positive multiplier on blur', async () => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2' placeholder],
    placeholder
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    const multiplier = wrapper.get('input[type="number"]')
    expect(multiplier.attributes('min')).toBe('0.01')
    expect(multiplier.attributes('step')).toBe('0.01')
    await multiplier.trigger('blur')

    expect(wrapper.emitted('update:modelValue')?.[0]?.[0]).toEqual({
      ...value,
      periods: [{ ...value.periods[0], multiplier: '2.00' placeholder],
    placeholder)

  placeholder)

  it.each(['1.234', '1e2', '.5'])('keeps invalid multiplier %s unchanged on blur', async invalid => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: invalid placeholder],
    placeholder
    const wrapper = mount(TimePricingSection, {
      props: { modelValue: value placeholder,
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    await wrapper.get('input[type="number"]').trigger('blur')

    expect(wrapper.emitted('update:modelValue')).toBeUndefined()
    expect(value.periods[0].multiplier).toBe(invalid)
  placeholder)

  it('uses unique input ids across component instances', () => {
    const value = {
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder
    const wrapper = mount({
      components: { TimePricingSection placeholder,
      data: () => ({ first: value, second: value placeholder),
      template: `
        <div>
          <TimePricingSection v-model="first" />
          <TimePricingSection v-model="second" />
        </div>
      `,
    placeholder, {
      global: { stubs: { Select: SelectStub, Icon: true placeholder placeholder,
    placeholder)

    const ids = wrapper.findAll('input').map(input => input.attributes('id'))
    expect(ids).toHaveLength(6)
    expect(new Set(ids).size).toBe(ids.length)
  placeholder)
placeholder)
