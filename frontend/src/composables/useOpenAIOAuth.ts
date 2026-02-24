import { ref placeholder from 'vue'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'

export interface OpenAITokenInfo {
  access_token?: string
  refresh_token?: string
  id_token?: string
  token_type?: string
  expires_in?: number
  expires_at?: number
  scope?: string
  email?: string
  name?: string
  // OpenAI specific IDs (extracted from ID Token)
  chatgpt_account_id?: string
  chatgpt_user_id?: string
  organization_id?: string
  [key: string]: unknown
placeholder

export type OpenAIOAuthPlatform = 'openai' | 'sora'

interface UseOpenAIOAuthOptions {
  platform?: OpenAIOAuthPlatform
placeholder

export function useOpenAIOAuth(options?: UseOpenAIOAuthOptions) {
  const appStore = useAppStore()
  const oauthPlatform = options?.platform ?? 'openai'
  const endpointPrefix = oauthPlatform === 'sora' ? '/admin/sora' : '/admin/openai'

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
  const validateRefreshToken = async (
    refreshToken: string,
    proxyId?: number | null
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
        `${endpointPrefixplaceholder/refresh-token`
      )
      return tokenInfo as OpenAITokenInfo
    placeholder catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to validate refresh token'
      appStore.showError(error.value)
      return null
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  // Validate Sora session token and get access token
  const validateSessionToken = async (
    sessionToken: string,
    proxyId?: number | null
  ): Promise<OpenAITokenInfo | null> => {
    if (!sessionToken.trim()) {
      error.value = 'Missing session token'
      return null
    placeholder
    loading.value = true
    error.value = ''
    try {
      const tokenInfo = await adminAPI.accounts.validateSoraSessionToken(
        sessionToken.trim(),
        proxyId,
        `${endpointPrefixplaceholder/st2at`
      )
      return tokenInfo as OpenAITokenInfo
    placeholder catch (err: any) {
      error.value = err.response?.data?.detail || 'Failed to validate session token'
      appStore.showError(error.value)
      return null
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  // Build credentials for OpenAI OAuth account
  const buildCredentials = (tokenInfo: OpenAITokenInfo): Record<string, unknown> => {
    const creds: Record<string, unknown> = {
      access_token: tokenInfo.access_token,
      refresh_token: tokenInfo.refresh_token,
      token_type: tokenInfo.token_type,
      expires_in: tokenInfo.expires_in,
      expires_at: tokenInfo.expires_at,
      scope: tokenInfo.scope
    placeholder

    // Include OpenAI specific IDs (required for forwarding)
    if (tokenInfo.chatgpt_account_id) {
      creds.chatgpt_account_id = tokenInfo.chatgpt_account_id
    placeholder
    if (tokenInfo.chatgpt_user_id) {
      creds.chatgpt_user_id = tokenInfo.chatgpt_user_id
    placeholder
    if (tokenInfo.organization_id) {
      creds.organization_id = tokenInfo.organization_id
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
    validateSessionToken,
    buildCredentials,
    buildExtraInfo
  placeholder
placeholder
