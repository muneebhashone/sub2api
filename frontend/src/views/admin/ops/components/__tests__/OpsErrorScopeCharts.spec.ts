import { mount placeholder from '@vue/test-utils'
import { describe, expect, it, vi placeholder from 'vitest'
import { defineComponent placeholder from 'vue'
import OpsErrorDistributionChart from '../OpsErrorDistributionChart.vue'
import OpsErrorTrendChart from '../OpsErrorTrendChart.vue'

vi.mock('chart.js', () => ({
  Chart: { register: vi.fn() placeholder,
  ArcElement: {placeholder,
  CategoryScale: {placeholder,
  Filler: {placeholder,
  Legend: {placeholder,
  LineElement: {placeholder,
  LinearScale: {placeholder,
  PointElement: {placeholder,
  Title: {placeholder,
  Tooltip: {placeholder,
placeholder))

vi.mock('vue-chartjs', async () => {
  const { defineComponent placeholder = await import('vue')

  return {
    Doughnut: defineComponent({
      name: 'Doughnut',
      props: {
        data: { type: Object, required: true placeholder,
        options: { type: Object, default: () => ({placeholder) placeholder,
      placeholder,
      template: '<div class="doughnut-stub" />',
    placeholder),
    Line: defineComponent({
      name: 'Line',
      props: {
        data: { type: Object, required: true placeholder,
        options: { type: Object, default: () => ({placeholder) placeholder,
      placeholder,
      template: '<div class="line-stub" />',
    placeholder),
  placeholder
placeholder)

vi.mock('../../utils/opsFormatters', () => ({
  formatHistoryLabel: (date: string | undefined) => date ?? '',
  sumNumbers: (values: Array<number | null | undefined>) =>
    values.reduce<number>((total, value) => total + (typeof value === 'number' && Number.isFinite(value) ? value : 0), 0),
placeholder))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()

  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    placeholder),
  placeholder
placeholder)

const HelpTooltipStub = defineComponent({
  name: 'HelpTooltip',
  props: {
    content: { type: String, default: '' placeholder,
  placeholder,
  template: '<span class="help-tooltip-stub" />',
placeholder)

const EmptyStateStub = defineComponent({
  name: 'EmptyState',
  props: {
    title: { type: String, default: '' placeholder,
    description: { type: String, default: '' placeholder,
  placeholder,
  template: '<div class="empty-state-stub" />',
placeholder)

const globalStubs = {
  stubs: {
    HelpTooltip: HelpTooltipStub,
    EmptyState: EmptyStateStub,
  placeholder,
placeholder

describe('Ops SLA-scoped error charts', () => {
  it('错误分布图按 SLA 错误数统计，不把业务限制错误算进请求错误分布', () => {
    const wrapper = mount(OpsErrorDistributionChart, {
      props: {
        loading: false,
        data: {
          total: 10,
          items: [
            { status_code: 400, total: 7, sla: 2, business_limited: 5 placeholder,
            { status_code: 503, total: 3, sla: 0, business_limited: 3 placeholder,
          ],
        placeholder,
      placeholder,
      global: globalStubs,
    placeholder)

    const doughnut = wrapper.findComponent({ name: 'Doughnut' placeholder)
    expect(doughnut.exists()).toBe(true)
    expect(doughnut.props('data')).toMatchObject({
      labels: ['admin.ops.client'],
      datasets: [{ data: [2] placeholder],
    placeholder)
  placeholder)

  it('错误分布图在只有业务限制错误时显示为空态', () => {
    const wrapper = mount(OpsErrorDistributionChart, {
      props: {
        loading: false,
        data: {
          total: 4,
          items: [{ status_code: 500, total: 4, sla: 0, business_limited: 4 placeholder],
        placeholder,
      placeholder,
      global: globalStubs,
    placeholder)

    expect(wrapper.findComponent({ name: 'Doughnut' placeholder).exists()).toBe(false)
    expect(wrapper.find('.empty-state-stub').exists()).toBe(true)
  placeholder)

  it('错误趋势图的请求错误详情按钮只按 SLA 错误启用', () => {
    const wrapper = mount(OpsErrorTrendChart, {
      props: {
        loading: false,
        timeRange: '1h',
        points: [
          {
            bucket_start: '2026-05-18T00:00:00Z',
            error_count_total: 5,
            business_limited_count: 5,
            error_count_sla: 0,
            upstream_error_count_excl_429_529: 0,
            upstream_429_count: 0,
            upstream_529_count: 0,
          placeholder,
        ],
      placeholder,
      global: globalStubs,
    placeholder)

    const requestErrorsButton = wrapper.findAll('button')[0]
    expect(requestErrorsButton.attributes('disabled')).toBeDefined()
  placeholder)
placeholder)
