import { ref placeholder from 'vue'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'

export interface OpenAITokenInfo {
  access_token?: string
  refresh_token?: string
  client_id?: string
  id_token?: string
  token_type?: string
  expires_in?: number
  expires_at?: number
  scope?: string
  email?: string
  name?: string
  plan_type?: string
  privacy_mode?: string
  // OpenAI specific IDs (extracted from ID Token)
  chatgpt_account_id?: string
  chatgpt_user_id?: string
  organization_id?: string
  [key: string]: unknown
placeholder

export type OpenAIOAuthPlatform = 'openai'

interface UseOpenAIOAuthOptions {
  platform?: OpenAIOAuthPlatform
placeholder

export function useOpenAIOAuth(_options?: UseOpenAIOAuthOptions) {
  const appStore = useAppStore()
  const endpointPrefix = '/admin/openai'

  // State
  const authUrl = ref('')
  const sessionId = ref('')
  const oauthState = ref('')
  const loading = ref(false)
  const error = ref('')

  // Reset state
  const resetState = () => {
    authUrl.value = ''
    sessionId.value = ''
    oauthState.value = ''
    loading.value = false
    error.value = ''
  placeholder

  // Generate auth URL for OpenAI OAuth
  const generateAuthUrl = async (
    proxyId?: number | null,
    redirectUri?: string
  ): Promise<boolean> => {
    loading.value = true
    authUrl.value = ''
    sessionId.value = ''
    oauthState.value = ''
    error.value = ''

    try {
      const payload: Record<string, unknown> = {placeholder
      if (proxyId) {
        payload.proxy_id = proxyId
      placeholder
      if (redirectUri) {
        payload.redirect_uri = redirectUri
      placeholder

      const response = await adminAPI.accounts.generateAuthUrl(
        `${endpointPrefixplaceholder/generate-auth-url`,
        payload
      )
      authUrl.value = response.auth_url
      sessionId.value = response.session_id
      try {
        const parsed = new URL(response.auth_url)
        oauthState.value = parsed.searchParams.get('state') || ''
      placeholder catch {
        oauthState.value = ''
      placeholder
      return true
    placeholder catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to generate OpenAI auth URL'
      appStore.showError(error.value)
      return false
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  // Exchange auth code for tokens
  const exchangeAuthCode = async (
    code: string,
    currentSessionId: string,
    state: string,
    proxyId?: number | null
  ): Promise<OpenAITokenInfo | null> => {
    if (!code.trim() || !currentSessionId || !state.trim()) {
      error.value = 'Missing auth code, session ID, or state'
      return null
    placeholder

    loading.value = true
    error.value = ''

    try {
      const payload: { session_id: string; code: string; state: string; proxy_id?: number placeholder = {
        session_id: currentSessionId,
        code: code.trim(),
        state: state.trim()
      placeholder
      if (proxyId) {
        payload.proxy_id = proxyId
      placeholder

      const tokenInfo = await adminAPI.accounts.exchangeCode(`${endpointPrefixplaceholder/exchange-code`, payload)
      return tokenInfo as OpenAITokenInfo
    placeholder catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to exchange OpenAI auth code'
      appStore.showError(error.value)
      return null
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  // Validate refresh token and get full token info
  // clientId: 指定 OAuth client_id（用于第三方渠道获取的 RT，如 app_LlGpXReQgckcGGUo2JrYvtJK）
  const validateRefreshToken = async (
    refreshToken: string,
    proxyId?: number | null,
    clientId?: string
  ): Promise<OpenAITokenInfo | null> => {
    if (!refreshToken.trim()) {
      error.value = 'Missing refresh token'
      return null
    placeholder

    loading.value = true
    error.value = ''

    try {
      // Use dedicated refresh-token endpoint
      const tokenInfo = await adminAPI.accounts.refreshOpenAIToken(
        refreshToken.trim(),
        proxyId,
        `${endpointPrefixplaceholder/refresh-token`,
        clientId
      )
      return tokenInfo as OpenAITokenInfo
    placeholder catch (err: any) {
      error.value = err.response?.data?.detail || err.message || 'Failed to validate refresh token'
      appStore.showError(error.value)
      return null
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  // Build credentials for OpenAI OAuth account (aligned with backend BuildAccountCredentials)
  const buildCredentials = (tokenInfo: OpenAITokenInfo): Record<string, unknown> => {
    const creds: Record<string, unknown> = {
      access_token: tokenInfo.access_token,
      expires_at: tokenInfo.expires_at
    placeholder

    // 仅在返回了新的 refresh_token 时才写入，防止用空值覆盖已有令牌
    if (tokenInfo.refresh_token) {
      creds.refresh_token = tokenInfo.refresh_token
    placeholder
    if (tokenInfo.id_token) {
      creds.id_token = tokenInfo.id_token
    placeholder
    if (tokenInfo.email) {
      creds.email = tokenInfo.email
    placeholder
    if (tokenInfo.chatgpt_account_id) {
      creds.chatgpt_account_id = tokenInfo.chatgpt_account_id
    placeholder
    if (tokenInfo.chatgpt_user_id) {
      creds.chatgpt_user_id = tokenInfo.chatgpt_user_id
    placeholder
    if (tokenInfo.organization_id) {
      creds.organization_id = tokenInfo.organization_id
    placeholder
    if (tokenInfo.plan_type) {
      creds.plan_type = tokenInfo.plan_type
    placeholder
    if (tokenInfo.client_id) {
      creds.client_id = tokenInfo.client_id
    placeholder

    return creds
  placeholder

  // Build extra info from token response
  const buildExtraInfo = (tokenInfo: OpenAITokenInfo): Record<string, string> | undefined => {
    const extra: Record<string, string> = {placeholder
    if (tokenInfo.email) {
      extra.email = tokenInfo.email
    placeholder
    if (tokenInfo.name) {
      extra.name = tokenInfo.name
    placeholder
    if (tokenInfo.privacy_mode) {
      extra.privacy_mode = tokenInfo.privacy_mode
    placeholder
    return Object.keys(extra).length > 0 ? extra : undefined
  placeholder

  return {
    // State
    authUrl,
    sessionId,
    oauthState,
    loading,
    error,
    // Methods
    resetState,
    generateAuthUrl,
    exchangeAuthCode,
    validateRefreshToken,
    buildCredentials,
    buildExtraInfo
  placeholder
placeholder
