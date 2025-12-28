import { ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { adminAPI placeholder from '@/api/admin'
import type { AntigravityTokenInfo placeholder from '@/api/admin/antigravity'

export function useAntigravityOAuth() {
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

  const generateAuthUrl = async (proxyId: number | null | undefined): Promise<boolean> => {
    loading.value = true
    authUrl.value = ''
    sessionId.value = ''
    state.value = ''
    error.value = ''

    try {
      const payload: Record<string, unknown> = {placeholder
      if (proxyId) payload.proxy_id = proxyId

      const response = await adminAPI.antigravity.generateAuthUrl(payload as any)
      authUrl.value = response.auth_url
      sessionId.value = response.session_id
      state.value = response.state
      return true
    placeholder catch (err: any) {
      error.value =
        err.response?.data?.detail || t('admin.accounts.oauth.antigravity.failedToGenerateUrl')
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
  placeholder): Promise<AntigravityTokenInfo | null> => {
    const code = params.code?.trim()
    if (!code || !params.sessionId || !params.state) {
      error.value = t('admin.accounts.oauth.antigravity.missingExchangeParams')
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

      const tokenInfo = await adminAPI.antigravity.exchangeCode(payload as any)
      return tokenInfo as AntigravityTokenInfo
    placeholder catch (err: any) {
      error.value =
        err.response?.data?.detail || t('admin.accounts.oauth.antigravity.failedToExchangeCode')
      appStore.showError(error.value)
      return null
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  const buildCredentials = (tokenInfo: AntigravityTokenInfo): Record<string, unknown> => {
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
      project_id: tokenInfo.project_id,
      email: tokenInfo.email
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
    buildCredentials
  placeholder
placeholder
