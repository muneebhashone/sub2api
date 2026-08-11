const ipGeoMocks = vi.hoisted(() => ({
  getEntry: vi.fn(() => ({ status: 'idle' as const placeholder)),
  fetchOne: vi.fn(),
  fetchBatch: vi.fn(),
placeholder))

const appStoreMocks = vi.hoisted(() => ({
  showSuccess: vi.fn(),
  showError: vi.fn(),
placeholder))

vi.mock('@/utils/ipGeoLookup', () => ipGeoMocks)
vi.mock('@/stores/app', () => ({ useAppStore: () => appStoreMocks placeholder))

import { describe, expect, it, vi, beforeEach placeholder from 'vitest'
import { mount placeholder from '@vue/test-utils'
import { nextTick placeholder from 'vue'

import UsageTable from '../UsageTable.vue'

const messages: Record<string, string> = {
  'admin.usage.userDeletedBadge': 'Deleted',
  'usage.costDetails': 'Cost Breakdown',
  'admin.usage.inputCost': 'Input Cost',
  'admin.usage.outputCost': 'Output Cost',
  'admin.usage.cacheCreationCost': 'Cache Creation Cost',
  'admin.usage.cacheReadCost': 'Cache Read Cost',
  'usage.inputTokenPrice': 'Input price',
  'usage.outputTokenPrice': 'Output price',
  'usage.perMillionTokens': '/ 1M tokens',
  'usage.serviceTier': 'Service tier',
  'usage.serviceTierPriority': 'Fast',
  'usage.serviceTierFlex': 'Flex',
  'usage.serviceTierStandard': 'Standard',
  'usage.rate': 'Rate',
  'usage.accountMultiplier': 'Account rate',
  'usage.original': 'Original',
  'usage.userBilled': 'User billed',
  'usage.accountBilled': 'Account billed',
  'usage.imageUnit': ' images',
  'usage.imageCount': 'Image count',
  'usage.imageBillingSize': 'Billing size',
  'usage.imageInputSize': 'Input size',
  'usage.imageOutputSize': 'Output size',
  'usage.imageSizeSource': 'Size source',
  'usage.imageSizeBreakdown': 'Size breakdown',
  'usage.imageSizeSourceOutput': 'Upstream output',
  'usage.imageSizeSourceInput': 'Request input',
  'usage.imageSizeSourceDefault': 'Default billing tier',
  'usage.imageSizeSourceLegacy': 'Legacy record',
  'usage.imageSizeSourceMissing': 'Not recorded',
  'usage.imageSizeNotRecorded': 'not recorded',
  'usage.imageSizeLegacyUnstandardized': 'legacy unstandardized',
  'usage.imageSizeUnknown': 'unknown',
  'usage.imageUnitPrice': 'Per-image price',
  'usage.imageTotalPrice': 'Image total price',
  'admin.usage.billingModeToken': 'Token',
  'admin.usage.billingModePerRequest': 'Per request',
  'admin.usage.billingModeImage': 'Image',
	'admin.usage.requestIdCopied': 'Request ID copied',
	'keys.copied': 'Copied',
	'keys.copyToClipboard': 'Copy to clipboard',
	'common.copyFailed': 'Copy failed',
	'usage.requestedModel': 'Requested',
	'usage.sentUpstreamModel': 'Sent upstream',
	'usage.upstreamResponseModel': 'Upstream response',
	'usage.modelVariant': 'Possible version variant',
	'usage.modelMismatch': 'Different model',
placeholder

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => messages[key] ?? key,
    placeholder),
  placeholder
placeholder)

const DataTableStub = {
  props: ['data'],
  template: `
    <div>
      <div v-for="row in data" :key="row.request_id">
        <slot name="cell-model" :row="row" :value="row.model" />
        <slot name="cell-billing_mode" :row="row" />
        <slot name="cell-tokens" :row="row" />
        <slot name="cell-cost" :row="row" />
        <slot name="cell-request_id" :row="row" />
      </div>
    </div>
  `,
placeholder

const baseImageRow = {
  request_id: 'req-admin-image',
  model: 'gpt-image-2',
  actual_cost: 0.4,
  total_cost: 0.4,
  account_rate_multiplier: 1,
  rate_multiplier: 1,
  service_tier: null,
  input_cost: 0,
  output_cost: 0,
  cache_creation_cost: 0,
  cache_read_cost: 0,
  input_tokens: 0,
  output_tokens: 0,
  cache_creation_tokens: 0,
  cache_read_tokens: 0,
  cache_creation_5m_tokens: 0,
  cache_creation_1h_tokens: 0,
  cache_ttl_overridden: false,
  billing_mode: 'image',
  image_count: 2,
  image_size: '2K',
  image_input_size: null,
  image_output_size: null,
  image_size_source: null,
  image_size_breakdown: null,
placeholder

describe('admin UsageTable tooltip', () => {
  beforeEach(() => {
    vi.spyOn(HTMLElement.prototype, 'getBoundingClientRect').mockReturnValue({
      x: 0,
      y: 0,
      top: 20,
      left: 20,
      right: 120,
      bottom: 40,
      width: 100,
      height: 20,
      toJSON: () => ({placeholder),
    placeholder as DOMRect)
  placeholder)

  it('marks only usage rows that actually applied long-context billing', () => {
    const wrapper = mount(UsageTable, {
      props: {
        data: [
          {
            ...baseImageRow,
            request_id: 'req-long-context-enabled',
            long_context_billing_applied: true,
          placeholder,
          {
            ...baseImageRow,
            request_id: 'req-long-context-disabled',
            long_context_billing_applied: false,
          placeholder,
        ],
        loading: false,
        columns: [],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStub,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    expect(wrapper.findAll('[data-testid="long-context-billing-marker"]')).toHaveLength(1)
    expect(wrapper.get('[data-testid="long-context-billing-marker"]').text()).toBe('x2')
  placeholder)

  it('shows service tier and billing breakdown in cost tooltip', async () => {
    const row = {
      request_id: 'req-admin-1',
      actual_cost: 0.092883,
      total_cost: 0.092883,
      account_rate_multiplier: 1,
      rate_multiplier: 1,
      service_tier: 'priority',
      input_cost: 0.020285,
      output_cost: 0.00303,
      cache_creation_cost: 0,
      cache_read_cost: 0.069568,
      input_tokens: 4057,
      output_tokens: 101,
    placeholder

    const wrapper = mount(UsageTable, {
      props: {
        data: [row],
        loading: false,
        columns: [],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStub,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    const tooltipTriggers = wrapper.findAll('.group.relative')
    await tooltipTriggers[tooltipTriggers.length - 1].trigger('mouseenter')
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('Service tier')
    expect(text).toContain('Fast')
    expect(text).toContain('Rate')
    expect(text).toContain('1.00x')
    expect(text).toContain('Account rate')
    expect(text).toContain('User billed')
    expect(text).toContain('Account billed')
    expect(text).toContain('$0.092883')
    expect(text).toContain('$5.0000 / 1M tokens')
    expect(text).toContain('$30.0000 / 1M tokens')
    expect(text).toContain('$0.069568')
  placeholder)

  it('shows requested and upstream models separately for admin rows', () => {
    const row = {
      request_id: 'req-admin-model-1',
      model: 'claude-sonnet-4',
      upstream_model: 'claude-sonnet-4-20250514',
      actual_cost: 0,
      total_cost: 0,
      account_rate_multiplier: 1,
      rate_multiplier: 1,
      input_cost: 0,
      output_cost: 0,
      cache_creation_cost: 0,
      cache_read_cost: 0,
      input_tokens: 0,
      output_tokens: 0,
    placeholder

    const wrapper = mount(UsageTable, {
      props: {
        data: [row],
        loading: false,
        columns: [],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStub,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    const text = wrapper.text()
    expect(text).toContain('claude-sonnet-4')
    expect(text).toContain('claude-sonnet-4-20250514')
  placeholder)

	it.each([
		{
			name: 'possible version variant',
			responseModel: 'gpt-5.5-2026-08-01',
			expectedBadge: 'Possible version variant',
	placeholder,
		{
			name: 'different upstream model',
			responseModel: 'gpt-5.4',
			expectedBadge: 'Different model',
	placeholder,
	])('shows a compact upstream response audit marker for $name', ({ responseModel, expectedBadge placeholder) => {
		const wrapper = mount(UsageTable, {
			props: {
				data: [{
					request_id: `req-${responseModelplaceholder`,
					model: 'gpt-5.6-sol',
					upstream_model: 'gpt-5.5',
					model_mapping_chain: 'gpt-5.6-sol→gpt-5.5',
					upstream_response_model: responseModel,
					upstream_model_mismatch: true,
			placeholder],
				loading: false,
				columns: [],
		placeholder,
			global: {
				stubs: {
					DataTable: DataTableStub,
					EmptyState: true,
					Icon: true,
					Teleport: true,
			placeholder,
		placeholder,
	placeholder)

		const text = wrapper.text()
		expect(text).toContain('gpt-5.6-sol')
		expect(text).toContain('gpt-5.5')
		expect(text).toContain(responseModel)
		expect(text).toContain(expectedBadge)
placeholder)

  it.each([
    {
      name: 'defaulted row',
      row: {
        ...baseImageRow,
        request_id: 'req-admin-default-image',
        image_size: '2K',
        image_input_size: 'auto',
        image_output_size: null,
        image_size_source: 'default',
      placeholder,
      expected: ['2K', 'Default billing tier', 'auto', 'unknown'],
    placeholder,
    {
      name: 'output-sourced row',
      row: {
        ...baseImageRow,
        request_id: 'req-admin-output-image',
        image_size: '4K',
        image_input_size: '1024x1024',
        image_output_size: '3840x2160',
        image_size_source: 'output',
        image_size_breakdown: { '4K': 1 placeholder,
      placeholder,
      expected: ['4K', 'Upstream output', '1024x1024', '3840x2160', '4K x 1'],
    placeholder,
    {
      name: 'input-sourced row',
      row: {
        ...baseImageRow,
        request_id: 'req-admin-input-image',
        image_size: '1K',
        image_input_size: '1024x1024',
        image_output_size: null,
        image_size_source: 'input',
      placeholder,
      expected: ['1K', 'Request input', '1024x1024', 'unknown'],
    placeholder,
    {
      name: 'legacy unstandardized row',
      row: {
        ...baseImageRow,
        request_id: 'req-admin-legacy-unstandardized-image',
        image_size: '512x512',
        image_input_size: null,
        image_output_size: null,
        image_size_source: null,
      placeholder,
      expected: ['legacy unstandardized: 512x512', 'Legacy record', 'unknown'],
    placeholder,
  ])('shows image usage metadata for $name', async ({ row, expected placeholder) => {
    const wrapper = mount(UsageTable, {
      props: {
        data: [row],
        loading: false,
        columns: [],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStub,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    await wrapper.find('.group.relative').trigger('mouseenter')
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('Image count')
    expect(text).toContain('Billing size')
    expect(text).toContain('Size source')
    expect(text).toContain('Input size')
    expect(text).toContain('Output size')
    expect(text).toContain('Per-image price')
    expect(text).toContain('Image total price')
    for (const value of expected) {
      expect(text).toContain(value)
    placeholder
  placeholder)

  it('displays historical image rows with missing billing_mode as image usage without a 2K fallback', async () => {
    const wrapper = mount(UsageTable, {
      props: {
        data: [
          {
            ...baseImageRow,
            request_id: 'req-admin-legacy-missing-image',
            billing_mode: null,
            image_size: null,
            image_input_size: null,
            image_output_size: null,
            image_size_source: null,
            image_size_breakdown: null,
          placeholder,
        ],
        loading: false,
        columns: [],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStub,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    await wrapper.find('.group.relative').trigger('mouseenter')
    await nextTick()

    const text = wrapper.text()
    expect(text).toContain('Image')
    expect(text).toContain('Image count')
    expect(text).toContain('Per-image price')
    expect(text).toContain('not recorded')
    expect(text).not.toContain('(2K)')
  placeholder)
placeholder)

describe('admin UsageTable request ID column', () => {
  beforeEach(() => {
    appStoreMocks.showSuccess.mockReset()
    appStoreMocks.showError.mockReset()
  placeholder)

  it('renders and copies the request ID', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined)
    vi.stubGlobal('navigator', { clipboard: { writeText placeholder placeholder)

    const wrapper = mount(UsageTable, {
      props: {
        data: [{ ...baseImageRow, request_id: 'req-admin-visible-id' placeholder],
        loading: false,
        columns: [{ key: 'request_id', label: 'Request ID' placeholder],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStub,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    expect(wrapper.text()).toContain('req-admin-visible-id')
    await wrapper.get('button[title="Copy to clipboard"]').trigger('click')

    expect(writeText).toHaveBeenCalledWith('req-admin-visible-id')
    expect(appStoreMocks.showSuccess).toHaveBeenCalledWith('Request ID copied')
  placeholder)
placeholder)

describe('admin UsageTable IP geolocation batch toolbar', () => {
  const DataTableStubWithIp = {
    props: ['data'],
    template: `
      <div>
        <div v-for="row in data" :key="row.request_id">
          <slot name="cell-ip_address" :row="row" />
        </div>
      </div>
    `,
  placeholder

  beforeEach(() => {
    ipGeoMocks.getEntry.mockReset()
    ipGeoMocks.fetchOne.mockReset()
    ipGeoMocks.fetchBatch.mockReset()
    ipGeoMocks.getEntry.mockReturnValue({ status: 'idle' placeholder)
  placeholder)

  it('does not render the batch toolbar when the ip_address column is not visible', () => {
    const wrapper = mount(UsageTable, {
      props: {
        data: [{ request_id: 'r1', ip_address: '8.8.8.8' placeholder],
        loading: false,
        columns: [],
      placeholder,
      global: { stubs: { DataTable: DataTableStubWithIp, EmptyState: true, Teleport: true placeholder placeholder,
    placeholder)
    expect(wrapper.text()).not.toContain('usage.ipGeo.batchFetch')
  placeholder)

  it('renders the batch toolbar with a pending count when the ip_address column is visible', () => {
    const wrapper = mount(UsageTable, {
      props: {
        data: [
          { request_id: 'r1', ip_address: '8.8.8.8' placeholder,
          { request_id: 'r2', ip_address: '8.8.8.8' placeholder,
          { request_id: 'r3', ip_address: '1.1.1.1' placeholder,
        ],
        loading: false,
        columns: [{ key: 'ip_address', label: 'IP' placeholder],
      placeholder,
      global: { stubs: { DataTable: DataTableStubWithIp, EmptyState: true, Teleport: true placeholder placeholder,
    placeholder)
    expect(wrapper.text()).toContain('usage.ipGeo.pending')
    const button = wrapper.find('button')
    expect(button.exists()).toBe(true)
    expect((button.element as HTMLButtonElement).disabled).toBe(false)
  placeholder)

  it('fetches deduplicated IPs from the current page when the batch button is clicked', async () => {
    ipGeoMocks.fetchBatch.mockResolvedValue(true)
    const wrapper = mount(UsageTable, {
      props: {
        data: [
          { request_id: 'r1', ip_address: '8.8.8.8' placeholder,
          { request_id: 'r2', ip_address: '8.8.8.8' placeholder,
          { request_id: 'r3', ip_address: '1.1.1.1' placeholder,
        ],
        loading: false,
        columns: [{ key: 'ip_address', label: 'IP' placeholder],
      placeholder,
      global: { stubs: { DataTable: DataTableStubWithIp, EmptyState: true, Teleport: true placeholder placeholder,
    placeholder)
    await wrapper.find('button').trigger('click')
    expect(ipGeoMocks.fetchBatch).toHaveBeenCalledWith(['8.8.8.8', '1.1.1.1'])
    expect(wrapper.emitted('ipGeoBatchFailed')).toBeUndefined()
  placeholder)

  it('emits ipGeoBatchFailed when the batch request reports a network-level failure', async () => {
    ipGeoMocks.fetchBatch.mockResolvedValue(false)
    const wrapper = mount(UsageTable, {
      props: {
        data: [{ request_id: 'r1', ip_address: '8.8.8.8' placeholder],
        loading: false,
        columns: [{ key: 'ip_address', label: 'IP' placeholder],
      placeholder,
      global: { stubs: { DataTable: DataTableStubWithIp, EmptyState: true, Teleport: true placeholder placeholder,
    placeholder)
    await wrapper.find('button').trigger('click')
    expect(wrapper.emitted('ipGeoBatchFailed')).toHaveLength(1)
  placeholder)

  it('renders IpGeoCell content for ip_address cells', () => {
    ipGeoMocks.getEntry.mockReturnValue({ status: 'success', label: 'CN · Guangdong · Shenzhen', detail: {placeholder placeholder)
    const wrapper = mount(UsageTable, {
      props: {
        data: [{ request_id: 'r1', ip_address: '121.35.47.43' placeholder],
        loading: false,
        columns: [{ key: 'ip_address', label: 'IP' placeholder],
      placeholder,
      global: { stubs: { DataTable: DataTableStubWithIp, EmptyState: true, Teleport: true placeholder placeholder,
    placeholder)
    expect(wrapper.text()).toContain('121.35.47.43')
    expect(wrapper.text()).toContain('CN · Guangdong · Shenzhen')
  placeholder)
placeholder)

// A DataTable stub that also renders cell-user, so the deleted badge can be asserted.
const DataTableStubWithUser = {
  props: ['data'],
  template: `
    <div>
      <div v-for="row in data" :key="row.request_id">
        <slot name="cell-user" :row="row" />
        <slot name="cell-model" :row="row" :value="row.model" />
        <slot name="cell-billing_mode" :row="row" />
        <slot name="cell-tokens" :row="row" />
        <slot name="cell-cost" :row="row" />
      </div>
    </div>
  `,
placeholder

describe('admin UsageTable deleted-user badge', () => {
  it('renders deleted badge for a soft-deleted user row', () => {
    const row = {
      request_id: 'req-deleted-user-1',
      model: 'claude-3',
      user_id: 2,
      user: { id: 2, email: 'd@test.com', deleted_at: '2026-05-28T00:00:00Z' placeholder,
      actual_cost: 0,
      total_cost: 0,
      input_cost: 0,
      output_cost: 0,
      rate_multiplier: 1,
      input_tokens: 1,
      output_tokens: 1,
    placeholder

    const wrapper = mount(UsageTable, {
      props: {
        data: [row],
        loading: false,
        columns: [{ key: 'user', label: 'User' placeholder],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStubWithUser,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    expect(wrapper.text()).toContain('Deleted')
    expect(wrapper.text()).toContain('d@test.com')
  placeholder)

  it('does NOT render deleted badge for an active user row', () => {
    const row = {
      request_id: 'req-active-user-1',
      model: 'claude-3',
      user_id: 3,
      user: { id: 3, email: 'active@test.com', deleted_at: null placeholder,
      actual_cost: 0,
      total_cost: 0,
      input_cost: 0,
      output_cost: 0,
      rate_multiplier: 1,
      input_tokens: 1,
      output_tokens: 1,
    placeholder

    const wrapper = mount(UsageTable, {
      props: {
        data: [row],
        loading: false,
        columns: [{ key: 'user', label: 'User' placeholder],
      placeholder,
      global: {
        stubs: {
          DataTable: DataTableStubWithUser,
          EmptyState: true,
          Icon: true,
          Teleport: true,
        placeholder,
      placeholder,
    placeholder)

    expect(wrapper.text()).not.toContain('Deleted')
    expect(wrapper.text()).toContain('active@test.com')
  placeholder)
placeholder)
