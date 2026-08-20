import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { defineComponent placeholder from 'vue'
import { mount placeholder from '@vue/test-utils'

const { updateAccountMock, checkMixedChannelRiskMock, authIsSimpleMode placeholder = vi.hoisted(() => ({
  updateAccountMock: vi.fn(),
  checkMixedChannelRiskMock: vi.fn(),
  authIsSimpleMode: { value: true placeholder
placeholder))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn(),
    showInfo: vi.fn()
  placeholder)
placeholder))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    get isSimpleMode() {
      return authIsSimpleMode.value
    placeholder
  placeholder)
placeholder))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    accounts: {
      update: updateAccountMock,
      checkMixedChannelRisk: checkMixedChannelRiskMock
    placeholder,
    settings: {
      getWebSearchEmulationConfig: vi.fn().mockResolvedValue({ enabled: false, providers: [] placeholder),
      getSettings: vi.fn().mockResolvedValue({placeholder)
    placeholder,
    tlsFingerprintProfiles: {
      list: vi.fn().mockResolvedValue([])
    placeholder
  placeholder
placeholder))

vi.mock('@/api/admin/accounts', () => ({
  getAntigravityDefaultModelMapping: vi.fn()
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

import EditAccountModal from '../EditAccountModal.vue'

const BaseDialogStub = defineComponent({
  name: 'BaseDialog',
  props: {
    show: {
      type: Boolean,
      default: false
    placeholder
  placeholder,
  template: '<div v-if="show"><slot /><slot name="footer" /></div>'
placeholder)

const ModelWhitelistSelectorStub = defineComponent({
  name: 'ModelWhitelistSelector',
  props: {
    modelValue: {
      type: Array,
      default: () => []
    placeholder
  placeholder,
  emits: ['update:modelValue'],
  template: `
    <div>
      <button
        type="button"
        data-testid="rewrite-to-snapshot"
        @click="$emit('update:modelValue', ['gpt-5.2-2025-12-11'])"
      >
        rewrite
      </button>
      <span data-testid="model-whitelist-value">
        {{ Array.isArray(modelValue) ? modelValue.join(',') : '' placeholderplaceholder
      </span>
    </div>
  `
placeholder)

const SelectStub = defineComponent({
  name: 'SelectStub',
  props: {
    modelValue: {
      type: [String, Number, Boolean, null],
      default: ''
    placeholder,
    options: {
      type: Array,
      default: () => []
    placeholder
  placeholder,
  emits: ['update:modelValue'],
  template: `
    <select
      v-bind="$attrs"
      :value="modelValue"
      @change="$emit('update:modelValue', $event.target.value)"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label placeholderplaceholder
      </option>
    </select>
  `
placeholder)

const GroupSelectorStub = defineComponent({
  name: 'GroupSelector',
  props: {
    modelValue: {
      type: Array,
      default: () => []
    placeholder
  placeholder,
  emits: ['update:modelValue'],
  template: `
    <div data-testid="group-selector">
      <button
        type="button"
        data-testid="set-shadow-group"
        @click="$emit('update:modelValue', [7])"
      >
        group
      </button>
    </div>
  `
placeholder)

function buildAccount() {
  return {
    id: 1,
    name: 'OpenAI Key',
    notes: '',
    platform: 'openai',
    type: 'apikey',
    credentials: {
      api_key: 'sk-test',
      base_url: 'https://api.openai.com',
      model_mapping: {
        'gpt-5.2': 'gpt-5.2'
      placeholder
    placeholder,
    extra: {placeholder,
    proxy_id: null,
    concurrency: 1,
    priority: 1,
    rate_multiplier: 1,
    status: 'active',
    group_ids: [],
    expires_at: null,
    auto_pause_on_expired: false
  placeholder as any
placeholder

function buildOpenAISparkShadowAccount() {
  const account = buildAccount()
  return {
    ...account,
    id: 4,
    name: 'OpenAI Spark Shadow',
    type: 'oauth',
    parent_account_id: 1,
    credentials: {
      access_token: 'parent-access-token',
      refresh_token: 'parent-refresh-token',
      api_key: 'sk-parent',
      base_url: 'https://api.openai.com',
      model_mapping: {
        'gpt-5.3-codex-spark': 'gpt-5.3-codex-spark'
      placeholder,
      compact_model_mapping: {
        'gpt-5.3-codex-spark': 'gpt-5.3-codex-spark-compact'
      placeholder
    placeholder,
    group_ids: []
  placeholder as any
placeholder

function buildVertexAccount() {
  return {
    id: 2,
    name: 'Vertex SA',
    notes: '',
    platform: 'gemini',
    type: 'service_account',
    credentials: {
      service_account_json: '{"type":"service_account","client_email":"sa@example.iam.gserviceaccount.com","private_key":"-----BEGIN PRIVATE KEY-----\\nMIIE\\n-----END PRIVATE KEY-----\\n"placeholder',
      project_id: 'demo-project',
      client_email: 'sa@example.iam.gserviceaccount.com',
      location: 'us-central1',
      tier_id: 'vertex'
    placeholder,
    extra: {placeholder,
    proxy_id: null,
    concurrency: 1,
    priority: 1,
    rate_multiplier: 1,
    status: 'active',
    group_ids: [],
    expires_at: null,
    auto_pause_on_expired: false
  placeholder as any
placeholder

function buildAntigravityAccount(projectId = 'configured-project') {
  return {
    id: 3,
    name: 'Antigravity OAuth',
    notes: '',
    platform: 'antigravity',
    type: 'oauth',
    credentials: {
      antigravity_project_id: projectId,
      model_mapping: {
        'gemini-2.5-flash': 'gemini-2.5-flash'
      placeholder
    placeholder,
    extra: {placeholder,
    proxy_id: null,
    concurrency: 1,
    priority: 1,
    rate_multiplier: 1,
    status: 'active',
    group_ids: [],
    expires_at: null,
    auto_pause_on_expired: false
  placeholder as any
placeholder

function buildGrokOAuthAccount() {
  return {
    id: 5,
    name: 'Grok OAuth',
    notes: '',
    platform: 'grok',
    type: 'oauth',
    credentials: {
      refresh_token: 'grok-rt',
      base_url: 'https://api.x.ai/v1',
      model_mapping: {
        'grok-latest': 'grok-4.3'
      placeholder
    placeholder,
    extra: {placeholder,
    proxy_id: null,
    concurrency: 1,
    priority: 1,
    rate_multiplier: 1,
    status: 'active',
    group_ids: [],
    expires_at: null,
    auto_pause_on_expired: false
  placeholder as any
placeholder

function buildGrokAPIKeyAccount() {
  return {
    ...buildAccount(),
    id: 6,
    name: 'Grok API Key',
    platform: 'grok',
    credentials: {placeholder,
    credentials_status: { has_api_key: true placeholder,
    concurrency: 2
  placeholder as any
placeholder

function buildOpenAISetupTokenAccount() {
  return {
    ...buildAccount(),
    type: 'setup-token',
    extra: {
      openai_oauth_responses_websockets_v2_mode: 'ctx_pool',
      openai_oauth_responses_websockets_v2_enabled: true
    placeholder
  placeholder as any
placeholder

function mountModal(account = buildAccount()) {
  return mount(EditAccountModal, {
    props: {
      show: true,
      account,
      proxies: [],
      groups: []
    placeholder,
    global: {
      stubs: {
        BaseDialog: BaseDialogStub,
        Select: SelectStub,
        Icon: true,
        ProxySelector: true,
        GroupSelector: GroupSelectorStub,
        ModelWhitelistSelector: ModelWhitelistSelectorStub
      placeholder
    placeholder
  placeholder)
placeholder

describe('EditAccountModal', () => {
  beforeEach(() => {
    authIsSimpleMode.value = true
  placeholder)

  it('reopening the same account rehydrates the OpenAI whitelist from props', async () => {
    const account = buildAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect(wrapper.get('[data-testid="model-whitelist-value"]').text()).toBe('gpt-5.2')

    await wrapper.get('[data-testid="rewrite-to-snapshot"]').trigger('click')
    expect(wrapper.get('[data-testid="model-whitelist-value"]').text()).toBe('gpt-5.2-2025-12-11')

    await wrapper.setProps({ show: false placeholder)
    await wrapper.setProps({ show: true placeholder)

    expect(wrapper.get('[data-testid="model-whitelist-value"]').text()).toBe('gpt-5.2')

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.model_mapping).toEqual({
      'gpt-5.2': 'gpt-5.2'
    placeholder)
  placeholder)

  it('preserves adaptive GLM endpoints on submit', async () => {
    const account = buildAccount()
    account.platform = 'zhipu'
    account.credentials = {
      api_key: 'sk-glm',
      account_mode: 'coding',
      api_protocol: 'adaptive',
      base_url: 'https://open.bigmodel.cn/api/coding/paas/v4',
      api_base_urls: {
        chat_completions: 'https://open.bigmodel.cn/api/coding/paas/v4',
        anthropic: 'https://open.bigmodel.cn/api/anthropic'
      placeholder
    placeholder
    updateAccountMock.mockReset().mockResolvedValue(account)
    checkMixedChannelRiskMock.mockReset().mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials).toMatchObject({
      account_mode: 'coding',
      api_protocol: 'adaptive',
      base_url: 'https://open.bigmodel.cn/api/coding/paas/v4',
      api_base_urls: {
        chat_completions: 'https://open.bigmodel.cn/api/coding/paas/v4',
        anthropic: 'https://open.bigmodel.cn/api/anthropic'
      placeholder
    placeholder)
  placeholder)

  it.each([
    ['explicit Chat Completions', 'chat_completions'],
    ['legacy missing protocol', undefined]
  ])('preserves a custom CN relay for %s accounts', async (_name, storedProtocol) => {
    const account = buildAccount()
    account.platform = 'zhipu'
    account.credentials = {
      api_key: 'sk-glm',
      account_mode: 'payg',
      base_url: 'https://relay.example.com/v1'
    placeholder
    if (storedProtocol) {
      account.credentials.api_protocol = storedProtocol
    placeholder
    updateAccountMock.mockReset().mockResolvedValue(account)
    checkMixedChannelRiskMock.mockReset().mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    const submittedCredentials = updateAccountMock.mock.calls[0]?.[1]?.credentials
    expect(submittedCredentials).toMatchObject({
      account_mode: 'payg',
      api_protocol: 'chat_completions',
      base_url: 'https://relay.example.com/v1'
    placeholder)
    expect(submittedCredentials).not.toHaveProperty('api_base_urls')
  placeholder)

  it('uses the legacy base_url when adaptive endpoints are missing', async () => {
    const account = buildAccount()
    account.platform = 'zhipu'
    account.credentials = {
      api_key: 'sk-glm',
      account_mode: 'payg',
      api_protocol: 'adaptive',
      base_url: 'https://relay.example.com/v1',
      api_base_urls: {
        chat_completions: '   '
      placeholder
    placeholder
    updateAccountMock.mockReset().mockResolvedValue(account)
    checkMixedChannelRiskMock.mockReset().mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials).toMatchObject({
      api_protocol: 'adaptive',
      base_url: 'https://relay.example.com/v1',
      api_base_urls: {
        chat_completions: 'https://relay.example.com/v1',
        anthropic: 'https://open.bigmodel.cn/api/anthropic'
      placeholder
    placeholder)
  placeholder)

  it('carries a fixed Chat relay into Adaptive when the user switches protocols', async () => {
    const account = buildAccount()
    account.platform = 'zhipu'
    account.credentials = {
      api_key: 'sk-glm',
      account_mode: 'payg',
      api_protocol: 'chat_completions',
      base_url: 'https://relay.example.com/v1'
    placeholder
    updateAccountMock.mockReset().mockResolvedValue(account)
    checkMixedChannelRiskMock.mockReset().mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)
    const adaptiveButton = wrapper
      .findAll('button')
      .find(button => button.text().includes('admin.accounts.cnProviders.apiProtocol.adaptive'))
    expect(adaptiveButton).toBeDefined()
    await adaptiveButton!.trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials).toMatchObject({
      api_protocol: 'adaptive',
      base_url: 'https://relay.example.com/v1',
      api_base_urls: {
        chat_completions: 'https://relay.example.com/v1'
      placeholder
    placeholder)
  placeholder)

  it.each([
    {
      name: 'Anthropic',
      platform: 'zhipu',
      protocol: 'anthropic',
      baseUrl: 'https://relay.example.com/anthropic',
      expectedBaseUrl: 'https://open.bigmodel.cn/api/paas/v4',
      expectedProtocolUrls: {
        chat_completions: 'https://open.bigmodel.cn/api/paas/v4',
        anthropic: 'https://relay.example.com/anthropic'
      placeholder
    placeholder,
    {
      name: 'Responses',
      platform: 'deepseek',
      protocol: 'responses',
      baseUrl: 'https://relay.example.com/responses',
      expectedBaseUrl: 'https://api.deepseek.com',
      expectedProtocolUrls: {
        chat_completions: 'https://api.deepseek.com',
        anthropic: 'https://api.deepseek.com/anthropic',
        responses: 'https://relay.example.com/responses'
      placeholder
    placeholder
  ])('keeps a fixed $name relay in its protocol slot when switching to Adaptive', async (testCase) => {
    const account = buildAccount()
    account.platform = testCase.platform
    account.credentials = {
      api_key: 'sk-cn',
      account_mode: 'payg',
      api_protocol: testCase.protocol,
      base_url: testCase.baseUrl
    placeholder
    updateAccountMock.mockReset().mockResolvedValue(account)
    checkMixedChannelRiskMock.mockReset().mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)
    const adaptiveButton = wrapper
      .findAll('button')
      .find(button => button.text().includes('admin.accounts.cnProviders.apiProtocol.adaptive'))
    expect(adaptiveButton).toBeDefined()
    await adaptiveButton!.trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials).toMatchObject({
      api_protocol: 'adaptive',
      base_url: testCase.expectedBaseUrl,
      api_base_urls: testCase.expectedProtocolUrls
    placeholder)
  placeholder)

  it('preserves model mappings when editing the whitelist', async () => {
    const account = buildAccount()
    account.credentials.model_mapping = {
      'gpt-5.2': 'gpt-5.2',
      'gpt-latest': 'gpt-5.2'
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect(wrapper.get('[data-testid="model-whitelist-value"]').text()).toBe('gpt-5.2')

    await wrapper.get('[data-testid="rewrite-to-snapshot"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.model_mapping).toEqual({
      'gpt-5.2-2025-12-11': 'gpt-5.2-2025-12-11',
      'gpt-latest': 'gpt-5.2'
    placeholder)
  placeholder)

  it('submits OpenAI compact mode and compact-only model mapping', async () => {
    const account = buildAccount()
    account.extra = {
      openai_compact_mode: 'force_on'
    placeholder
    account.credentials = {
      ...account.credentials,
      compact_model_mapping: {
        'gpt-5.4': 'gpt-5.4-openai-compact'
      placeholder
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_compact_mode).toBe('force_on')
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.compact_model_mapping).toEqual({
      'gpt-5.4': 'gpt-5.4-openai-compact'
    placeholder)
  placeholder)

  it('loads and submits the per-account OpenAI long-context billing toggle', async () => {
    const account = buildAccount()
    account.extra = {
      openai_long_context_billing_enabled: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const toggle = wrapper.get('[data-testid="openai-long-context-billing-toggle"]')
    expect(toggle.attributes('aria-checked')).toBe('true')

    await toggle.trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_long_context_billing_enabled).toBe(false)
  placeholder)

  it('loads and clears the OAuth-only Codex namespace flatten toggle', async () => {
    const account = buildAccount()
    account.type = 'oauth'
    account.extra = {
      openai_responses_flatten_namespaces: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const toggle = wrapper.get('[data-testid="edit-openai-flatten-namespaces-toggle"]')

    // 关闭后应从 extra 中删除该键，而不是写入 false
    await toggle.trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty(
      'openai_responses_flatten_namespaces'
    )
  placeholder)

  it('submits the Codex namespace flatten toggle when switched on', async () => {
    const account = buildAccount()
    account.type = 'oauth'
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    await wrapper.get('[data-testid="edit-openai-flatten-namespaces-toggle"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_responses_flatten_namespaces).toBe(
      true
    )
  placeholder)

  it('hides the Codex namespace flatten toggle for non-OAuth OpenAI accounts', async () => {
    const account = buildAccount()
    const wrapper = mountModal(account)

    expect(wrapper.find('[data-testid="edit-openai-flatten-namespaces-toggle"]').exists()).toBe(
      false
    )
  placeholder)

  it('defaults legacy OpenAI accounts to long-context billing disabled', async () => {
    const account = buildAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const toggle = wrapper.get('[data-testid="openai-long-context-billing-toggle"]')
    expect(toggle.attributes('aria-checked')).toBe('false')

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_long_context_billing_enabled).toBe(false)
  placeholder)

  it('does not render or submit the long-context billing toggle for Spark shadow accounts', async () => {
    const account = buildOpenAISparkShadowAccount()
    account.extra = {
      openai_long_context_billing_enabled: false
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)
    const wrapper = mountModal(account)

    expect(wrapper.find('[data-testid="openai-long-context-billing-toggle"]').exists()).toBe(false)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty(
      'openai_long_context_billing_enabled'
    )
  placeholder)

  it('preserves an explicit OpenAI long-context billing opt-out', async () => {
    const account = buildAccount()
    account.extra = {
      openai_long_context_billing_enabled: false
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const toggle = wrapper.get('[data-testid="openai-long-context-billing-toggle"]')
    expect(toggle.attributes('aria-checked')).toBe('false')

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_long_context_billing_enabled).toBe(false)
  placeholder)

  it('fails closed for malformed OpenAI long-context billing values', async () => {
    const account = buildAccount()
    account.extra = {
      openai_long_context_billing_enabled: 'false'
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect(wrapper.get('[data-testid="openai-long-context-billing-toggle"]').attributes('aria-checked')).toBe('false')

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_long_context_billing_enabled).toBe(false)
  placeholder)

  it('loads and submits Grok OAuth model mapping edits', async () => {
    const account = buildGrokOAuthAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    expect(wrapper.text()).toContain('Imagine Image')
    expect(wrapper.text()).toContain('Imagine Video')

    const inputWithValue = (value: string) => {
      const input = wrapper
        .findAll('input')
        .find((input) => (input.element as HTMLInputElement).value === value)
      expect(input).toBeTruthy()
      return input!
    placeholder

    await inputWithValue('grok-latest').setValue('grok')
    await inputWithValue('grok-4.3').setValue('grok-build-0.1')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.model_mapping).toEqual({
      grok: 'grok-build-0.1'
    placeholder)
  placeholder)

  it('uses the official xAI base URL when a Grok API-key account omits base_url', async () => {
    const account = buildGrokAPIKeyAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect((wrapper.get('input[placeholder="https://api.x.ai/v1"]').element as HTMLInputElement).value)
      .toBe('https://api.x.ai/v1')

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.base_url).toBe('https://api.x.ai/v1')
  placeholder)

  it('only submits model mapping credentials when saving an OpenAI spark shadow account', async () => {
    authIsSimpleMode.value = false
    const account = buildOpenAISparkShadowAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('[data-testid="set-shadow-group"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    const payload = updateAccountMock.mock.calls[0]?.[1]
    expect(payload?.group_ids).toEqual([7])
    expect(payload?.credentials).toEqual({
      model_mapping: {
        'gpt-5.3-codex-spark': 'gpt-5.3-codex-spark'
      placeholder,
      compact_model_mapping: {
        'gpt-5.3-codex-spark': 'gpt-5.3-codex-spark-compact'
      placeholder
    placeholder)
  placeholder)

  it('submits OpenAI APIKey Responses support override mode', async () => {
    const account = buildAccount()
    account.extra = {
      openai_responses_mode: 'force_chat_completions',
      openai_responses_supported: false
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('[data-testid="openai-responses-mode-select"]').setValue('force_responses')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_responses_mode).toBe('force_responses')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_responses_supported).toBe(false)
  placeholder)

  it('submits the account upstream billing auto-probe setting', async () => {
    const account = buildAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const toggle = wrapper.get('[data-testid="upstream-billing-auto-probe"]')
    expect(toggle.attributes('aria-checked')).toBe('false')

    await toggle.trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.upstream_billing_probe_enabled).toBe(true)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty(
      'upstream_billing_probe_enabled'
    )
  placeholder)

  it('exposes the upstream billing auto-probe toggle for non-OpenAI API-key accounts', async () => {
    // 探测已放宽到全部 API-key 平台：grok 账号同样能开启并保存。
    const account = buildAccount()
    account.platform = 'grok'
    account.name = 'grok-relay'
    account.credentials = { api_key: 'sk-grok', base_url: 'https://relay.example/v1' placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const toggle = wrapper.get('[data-testid="upstream-billing-auto-probe"]')
    expect(toggle.attributes('aria-checked')).toBe('false')

    await toggle.trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.upstream_billing_probe_enabled).toBe(true)
  placeholder)

  it('enabling rate sync also enables probing and stops submitting a manual rate', async () => {
    const account = buildAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const syncToggle = wrapper.get('[data-testid="upstream-billing-rate-sync"]')
    const probeToggle = wrapper.get('[data-testid="upstream-billing-auto-probe"]')
    const rateInput = wrapper.get<HTMLInputElement>('[data-testid="account-rate-multiplier"]')
    expect(syncToggle.attributes('aria-checked')).toBe('false')
    expect(probeToggle.attributes('aria-checked')).toBe('false')
    expect(rateInput.element.disabled).toBe(false)
    expect(wrapper.text()).toContain('admin.accounts.billingRateMultiplierHint')
    expect(wrapper.text()).not.toContain('admin.accounts.upstreamBilling.syncRateManagedHint')

    await syncToggle.trigger('click')
    expect(syncToggle.attributes('aria-checked')).toBe('true')
    expect(probeToggle.attributes('aria-checked')).toBe('true')
    expect(rateInput.element.disabled).toBe(true)
    expect(wrapper.text()).toContain('admin.accounts.upstreamBilling.syncRateManagedHint')
    expect(wrapper.text()).not.toContain('admin.accounts.billingRateMultiplierHint')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    const payload = updateAccountMock.mock.calls[0]?.[1]
    expect(payload?.upstream_billing_probe_enabled).toBe(true)
    expect(payload?.upstream_billing_rate_sync_enabled).toBe(true)
    expect(payload).not.toHaveProperty('rate_multiplier')
  placeholder)

  it('disabling probing also disables rate sync and restores manual rate editing', async () => {
    const account = buildAccount()
    account.extra = {
      upstream_billing_probe_enabled: true,
      upstream_billing_rate_sync_enabled: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const syncToggle = wrapper.get('[data-testid="upstream-billing-rate-sync"]')
    const probeToggle = wrapper.get('[data-testid="upstream-billing-auto-probe"]')
    const rateInput = wrapper.get<HTMLInputElement>('[data-testid="account-rate-multiplier"]')
    expect(syncToggle.attributes('aria-checked')).toBe('true')
    expect(rateInput.element.disabled).toBe(true)

    await probeToggle.trigger('click')
    expect(probeToggle.attributes('aria-checked')).toBe('false')
    expect(syncToggle.attributes('aria-checked')).toBe('false')
    expect(rateInput.element.disabled).toBe(false)
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    const payload = updateAccountMock.mock.calls[0]?.[1]
    expect(payload?.upstream_billing_probe_enabled).toBe(false)
    expect(payload?.upstream_billing_rate_sync_enabled).toBe(false)
    expect(payload?.rate_multiplier).toBe(1)
  placeholder)

  it('disabling only rate sync keeps automatic probing enabled', async () => {
    const account = buildAccount()
    account.extra = {
      upstream_billing_probe_enabled: true,
      upstream_billing_rate_sync_enabled: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    await wrapper.get('[data-testid="upstream-billing-rate-sync"]').trigger('click')
    expect(wrapper.get('[data-testid="upstream-billing-auto-probe"]').attributes('aria-checked')).toBe(
      'true'
    )
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    const payload = updateAccountMock.mock.calls[0]?.[1]
    expect(payload?.upstream_billing_probe_enabled).toBe(true)
    expect(payload?.upstream_billing_rate_sync_enabled).toBe(false)
    expect(payload?.rate_multiplier).toBe(1)
  placeholder)

  it('clears OpenAI APIKey Responses override when set back to auto', async () => {
    const account = buildAccount()
    account.extra = {
      openai_responses_mode: 'force_chat_completions',
      openai_responses_supported: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('[data-testid="openai-responses-mode-select"]').setValue('auto')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('openai_responses_mode')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_responses_supported).toBe(true)
  placeholder)

  it('submits OpenAI APIKey endpoint capabilities from credentials', async () => {
    const account = buildAccount()
    account.credentials.openai_capabilities = ['chat_completions']
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect(wrapper.findAll('input[type="checkbox"]').some((input) => (input.element as HTMLInputElement).checked)).toBe(true)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.openai_capabilities).toEqual([
      'chat_completions'
    ])
  placeholder)

	it('submits OpenAI quota auto-pause thresholds in extra', async () => {
	  const account = buildAccount()
	  account.extra = {
		auto_pause_5h_threshold: 0.9,
		auto_pause_7d_threshold: 0.8
	  placeholder
	  updateAccountMock.mockReset()
	  checkMixedChannelRiskMock.mockReset()
	  checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
	  updateAccountMock.mockResolvedValue(account)

	  const wrapper = mountModal(account)

	  await wrapper.get('[data-testid="auto-pause-5h-threshold"]').setValue('95')
	  await wrapper.get('[data-testid="auto-pause-7d-threshold"]').setValue('96')
	  await wrapper.get('form#edit-account-form').trigger('submit.prevent')

	  expect(updateAccountMock).toHaveBeenCalledTimes(1)
	  expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.auto_pause_5h_threshold).toBe(0.95)
	  expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.auto_pause_7d_threshold).toBe(0.96)
placeholder)

	it('submits OpenAI quota auto-pause disable flag in extra', async () => {
	  // Toggling the per-account disable flag must persist as auto_pause_5h_disabled
	  // so an admin can exempt one account from auto-pause even when a global default
	  // threshold is configured (otherwise leaving the threshold blank would silently
	  // fall back to the global default).
	  const account = buildAccount()
	  updateAccountMock.mockReset()
	  checkMixedChannelRiskMock.mockReset()
	  checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
	  updateAccountMock.mockResolvedValue(account)

	  const wrapper = mountModal(account)

	  await wrapper.get('[data-testid="auto-pause-5h-disabled"]').trigger('click')
	  await wrapper.get('form#edit-account-form').trigger('submit.prevent')

	  expect(updateAccountMock).toHaveBeenCalledTimes(1)
	  expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.auto_pause_5h_disabled).toBe(true)
	  expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.auto_pause_7d_disabled).toBeUndefined()
placeholder)

  it('keeps at least one OpenAI APIKey endpoint capability selected', async () => {
    const account = buildAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    const chatCheckbox = wrapper.get<HTMLInputElement>(
      '[data-testid="openai-endpoint-capability-chat_completions"]'
    )
    const embeddingsCheckbox = wrapper.get<HTMLInputElement>(
      '[data-testid="openai-endpoint-capability-embeddings"]'
    )

    expect(chatCheckbox.element.checked).toBe(true)
    expect(embeddingsCheckbox.element.checked).toBe(true)

    await embeddingsCheckbox.setValue(false)

    expect(chatCheckbox.element.checked).toBe(true)
    expect(embeddingsCheckbox.element.checked).toBe(false)

    await chatCheckbox.setValue(false)

    expect(chatCheckbox.element.checked).toBe(true)
    expect(embeddingsCheckbox.element.checked).toBe(false)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.openai_capabilities).toEqual([
      'chat_completions'
    ])
  placeholder)

  it('disables text generation protocol when only embeddings requests are accepted', async () => {
    const account = buildAccount()
    account.credentials.openai_capabilities = ['embeddings']
    account.extra = {
      openai_responses_mode: 'force_responses',
      openai_responses_supported: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    const responsesModeSelect = wrapper.get<HTMLSelectElement>(
      '[data-testid="openai-responses-mode-select"]'
    )

    expect(responsesModeSelect.element.disabled).toBe(true)
    expect(wrapper.find('[data-testid="openai-responses-mode-not-applicable"]').exists()).toBe(true)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.openai_capabilities).toEqual([
      'embeddings'
    ])
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('openai_responses_mode')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_responses_supported).toBe(true)
  placeholder)

  it('submits Codex image tool force-inject mode as bridge override', async () => {
    const account = buildAccount()
    account.extra = {
      codex_image_generation_bridge: false,
      codex_image_generation_bridge_enabled: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect(wrapper.text()).toContain('admin.accounts.openai.codexImageTool')
    expect(wrapper.text()).toContain('admin.accounts.openai.codexImageToolDesc')
    expect(wrapper.text()).toContain('admin.accounts.openai.codexImageToolEnabledDesc')

    await wrapper.get('button[data-testid="codex-image-tool-enabled"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.codex_image_generation_bridge).toBe(true)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('codex_image_generation_bridge_enabled')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('codex_image_generation_explicit_tool_policy')
  placeholder)

  it('submits Codex image tool no-injection mode without strip policy', async () => {
    const account = buildAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('button[data-testid="codex-image-tool-disabled"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.codex_image_generation_bridge).toBe(false)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('codex_image_generation_explicit_tool_policy')
  placeholder)

  it('submits Codex image tool block mode as strip policy and clears bridge override', async () => {
    const account = buildAccount()
    account.extra = {
      codex_image_generation_bridge: true
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    expect(wrapper.text()).toContain('admin.accounts.openai.codexImageToolBlock')
    expect(wrapper.text()).toContain('admin.accounts.openai.codexImageToolBlockDesc')

    await wrapper.get('button[data-testid="codex-image-tool-block"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.codex_image_generation_explicit_tool_policy).toBe('strip')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('codex_image_generation_bridge')
  placeholder)

  it('loads strip policy as block mode and clears both keys when reset to inherit', async () => {
    const account = buildAccount()
    account.extra = {
      codex_image_generation_explicit_tool_policy: 'strip'
    placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('button[data-testid="codex-image-tool-inherit"]').trigger('click')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('codex_image_generation_explicit_tool_policy')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra).not.toHaveProperty('codex_image_generation_bridge')
  placeholder)

  it('setup-token account can select and submit OAuth WS mode', async () => {
    const account = buildOpenAISetupTokenAccount()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('[data-testid="edit-openai-ws-mode-select"]').setValue('http_bridge')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_oauth_responses_websockets_v2_mode).toBe('http_bridge')
    expect(updateAccountMock.mock.calls[0]?.[1]?.extra?.openai_oauth_responses_websockets_v2_enabled).toBe(true)
  placeholder)

  it('allows saving apikey account when backend redacted api_key but credentials_status reports it exists', async () => {
    // 新前端 + 新后端：响应已脱敏，credentials 里没有 api_key，credentials_status.has_api_key=true
    const account = buildAccount()
    account.credentials = {
      base_url: 'https://api.openai.com',
      model_mapping: { 'gpt-5.2': 'gpt-5.2' placeholder
    placeholder
    account.credentials_status = { has_api_key: true placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    // 用户未输入新 key 时，payload 不应带 api_key，由后端合并保留旧值
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials).not.toHaveProperty('api_key')
  placeholder)

  it('allows saving apikey account against legacy backend without credentials_status', async () => {
    // 新前端 + 旧后端：credentials_status 缺失，但 credentials.api_key 仍是明文，应允许保存
    const account = buildAccount()
    // 显式确保没有 credentials_status
    expect(account.credentials_status).toBeUndefined()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    // 旧后端响应未脱敏，原 api_key 会随 currentCredentials 一起传回去（旧行为，等价于无操作）
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.api_key).toBe('sk-test')
  placeholder)

  it('blocks apikey save when neither credentials_status nor legacy api_key indicates existence', async () => {
    const account = buildAccount()
    account.credentials = {
      base_url: 'https://api.openai.com'
    placeholder
    // 既没有 credentials_status 也没有旧的 api_key
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).not.toHaveBeenCalled()
  placeholder)

  it('allows saving Vertex SA account when backend redacted service_account_json but credentials_status reports it exists', async () => {
    // 新前端 + 新后端：响应已脱敏，credentials 里没有 service_account_json，credentials_status.has_service_account_json=true
    const account = buildVertexAccount()
    account.credentials = {
      project_id: 'demo-project',
      client_email: 'sa@example.iam.gserviceaccount.com',
      location: 'us-central1',
      tier_id: 'vertex'
    placeholder
    account.credentials_status = { has_service_account_json: true placeholder
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.project_id).toBe('demo-project')
  placeholder)

  it('allows saving Vertex SA account against legacy backend without credentials_status', async () => {
    // 新前端 + 旧后端：credentials_status 缺失，但 credentials.service_account_json 仍是明文，应允许保存
    const account = buildVertexAccount()
    expect(account.credentials_status).toBeUndefined()
    expect(account.credentials.service_account_json).toBeTruthy()
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
  placeholder)

  it('blocks Vertex SA save when neither credentials_status nor legacy json indicates existence', async () => {
    const account = buildVertexAccount()
    account.credentials = {
      project_id: 'demo-project',
      client_email: 'sa@example.iam.gserviceaccount.com',
      location: 'us-central1',
      tier_id: 'vertex'
    placeholder
    // 既没有 credentials_status 也没有旧的 service_account_json
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)

    const wrapper = mountModal(account)

    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).not.toHaveBeenCalled()
  placeholder)

  it('loads and submits Antigravity configured project fallback', async () => {
    const account = buildAntigravityAccount('configured-project')
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const input = wrapper.get<HTMLInputElement>('[data-testid="antigravity-project-id-input"]')
    expect(input.element.value).toBe('configured-project')

    await input.setValue('  updated-project  ')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials?.antigravity_project_id).toBe(
      'updated-project'
    )
  placeholder)

  it('clears Antigravity configured project fallback when input is empty', async () => {
    const account = buildAntigravityAccount('configured-project')
    updateAccountMock.mockReset()
    checkMixedChannelRiskMock.mockReset()
    checkMixedChannelRiskMock.mockResolvedValue({ has_risk: false placeholder)
    updateAccountMock.mockResolvedValue(account)

    const wrapper = mountModal(account)
    const input = wrapper.get<HTMLInputElement>('[data-testid="antigravity-project-id-input"]')

    await input.setValue('')
    await wrapper.get('form#edit-account-form').trigger('submit.prevent')

    expect(updateAccountMock).toHaveBeenCalledTimes(1)
    expect(updateAccountMock.mock.calls[0]?.[1]?.credentials).not.toHaveProperty(
      'antigravity_project_id'
    )
  placeholder)
placeholder)
