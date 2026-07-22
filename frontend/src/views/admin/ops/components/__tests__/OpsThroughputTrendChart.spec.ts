import { describe, expect, it, vi placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'

import OpsThroughputTrendChart from '../OpsThroughputTrendChart.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    placeholder),
  placeholder
placeholder)

vi.mock('vue-chartjs', () => ({
  Line: {
    props: ['data', 'options'],
    template: '<div class="line-chart" />',
  placeholder,
placeholder))

describe('OpsThroughputTrendChart', () => {
  it('allows the header controls to wrap on narrow screens', () => {
    const wrapper = mount(OpsThroughputTrendChart, {
      props: {
        points: [],
        loading: false,
        timeRange: '1h',
      placeholder,
      global: {
        stubs: {
          EmptyState: true,
          HelpTooltip: true,
        placeholder,
      placeholder,
    placeholder)

    const header = wrapper.get('[data-testid="throughput-chart-header"]')
    expect(header.classes()).toEqual(expect.arrayContaining(['flex-col', 'sm:flex-row']))

    const toolbar = wrapper.get('[data-testid="throughput-chart-toolbar"]')
    expect(toolbar.classes()).toEqual(expect.arrayContaining(['w-full', 'flex-wrap', 'sm:w-auto']))
    expect(toolbar.findAll('button')).toHaveLength(3)
    toolbar.findAll('button').forEach((button) => {
      expect(button.classes()).toContain('shrink-0')
      expect(button.classes()).not.toContain('ml-2')
    placeholder)
  placeholder)
placeholder)
