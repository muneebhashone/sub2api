import { ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import type { GeminiOAuthCapabilities placeholder from '@/api/admin/gemini'

export interface GeminiTokenInfo {
  access_token?: string
  refresh_token?: string
  token_type?: string
  scope?: string
  expires_at?: number | string
  project_id?: string
  oauth_type?: string
  [key: string]: unknown
placeholder

export function useGeminiOAuth() {
  const appStore = useAppStore()
  const { t placeholder = useI18n()

  const authUrl = ref('')
  const sessionId = ref('')
  const state = ref('')
  const loading = ref(false)
  const error = ref('')

  const resetState = () => {
    authUrl.value = ''
    sessionId.value = ''
    state.value = ''
    loading.value = false
    error.value = ''
  placeholder

  const generateAuthUrl = async (
    proxyId: number | null | undefined,
    projectId?: string | null,
    oauthType?: string
  ): Promise<boolean> => {
    loading.value = true
    authUrl.value = ''
    sessionId.value = ''
    state.value = ''
    error.value = ''

    try {
      const payload: Record<string, unknown> = {placeholder
      if (proxyId) payload.proxy_id = proxyId
      const trimmedProjectID = projectId?.trim()
      if (trimmedProjectID) payload.project_id = trimmedProjectID
      if (oauthType) payload.oauth_type = oauthType

      const response = await adminAPI.gemini.generateAuthUrl(payload as any)
      authUrl.value = response.auth_url
      sessionId.value = response.session_id
      state.value = response.state
      return true
    placeholder catch (err: any) {
      error.value = err.response?.data?.detail || t('admin.accounts.oauth.gemini.failedToGenerateUrl')
      appStore.showError(error.value)
      return false
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  const exchangeAuthCode = async (params: {
    code: string
    sessionId: string
    state: string
    proxyId?: number | null
    oauthType?: string
  placeholder): Promise<GeminiTokenInfo | null> => {
    const code = params.code?.trim()
    if (!code || !params.sessionId || !params.state) {
      error.value = t('admin.accounts.oauth.gemini.missingExchangeParams')
      return null
    placeholder

    loading.value = true
    error.value = ''

    try {
      const payload: Record<string, unknown> = {
        session_id: params.sessionId,
        state: params.state,
        code
      placeholder
      if (params.proxyId) payload.proxy_id = params.proxyId
      if (params.oauthType) payload.oauth_type = params.oauthType

      const tokenInfo = await adminAPI.gemini.exchangeCode(payload as any)
      return tokenInfo as GeminiTokenInfo
    placeholder catch (err: any) {
      // Check for specific missing project_id error
      const errorMessage = err.message || err.response?.data?.message || ''
      if (errorMessage.includes('missing project_id')) {
        error.value = t('admin.accounts.oauth.gemini.missingProjectId')
      placeholder else {
        error.value = errorMessage || t('admin.accounts.oauth.gemini.failedToExchangeCode')
      placeholder
      appStore.showError(error.value)
      return null
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  const buildCredentials = (tokenInfo: GeminiTokenInfo): Record<string, unknown> => {
    let expiresAt: string | undefined
    if (typeof tokenInfo.expires_at === 'number' && Number.isFinite(tokenInfo.expires_at)) {
      expiresAt = Math.floor(tokenInfo.expires_at).toString()
    placeholder else if (typeof tokenInfo.expires_at === 'string' && tokenInfo.expires_at.trim()) {
      expiresAt = tokenInfo.expires_at.trim()
    placeholder

    return {
      access_token: tokenInfo.access_token,
      refresh_token: tokenInfo.refresh_token,
      token_type: tokenInfo.token_type,
      expires_at: expiresAt,
      scope: tokenInfo.scope,
      project_id: tokenInfo.project_id,
      oauth_type: tokenInfo.oauth_type
    placeholder
  placeholder

  const getCapabilities = async (): Promise<GeminiOAuthCapabilities | null> => {
    try {
      return await adminAPI.gemini.getCapabilities()
    placeholder catch (err: any) {
      // Capabilities are optional for older servers; don't block the UI.
      return null
    placeholder
  placeholder

  return {
    authUrl,
    sessionId,
    state,
    loading,
    error,
    resetState,
    generateAuthUrl,
    exchangeAuthCode,
    buildCredentials,
    getCapabilities
  placeholder
placeholder
