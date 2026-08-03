/**
 * Authentication API endpoints
 * Handles user login, registration, and logout operations
 */

import { apiClient placeholder from './client'
import { refreshAuthTokens, type RefreshTokenResponse placeholder from './tokenRefresh'
export type { RefreshTokenResponse placeholder from './tokenRefresh'
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  CurrentUserResponse,
  SendVerifyCodeRequest,
  SendVerifyCodeResponse,
  PublicSettings,
  TotpLoginResponse,
  TotpLogin2FARequest
placeholder from '@/types'

/**
 * Login response type - can be either full auth or 2FA required
 */
export type LoginResponse = AuthResponse | TotpLoginResponse

/**
 * Type guard to check if login response requires 2FA
 */
export function isTotp2FARequired(response: LoginResponse): response is TotpLoginResponse {
  return 'requires_2fa' in response && response.requires_2fa === true
placeholder

/**
 * Store authentication token in localStorage
 */
export function setAuthToken(token: string): void {
  localStorage.setItem('auth_token', token)
placeholder

/**
 * Store refresh token in localStorage
 */
export function setRefreshToken(token: string): void {
  localStorage.setItem('refresh_token', token)
placeholder

/**
 * Store token expiration timestamp in localStorage
 * Converts expires_in (seconds) to absolute timestamp (milliseconds)
 */
export function setTokenExpiresAt(expiresIn: number): void {
  const expiresAt = Date.now() + expiresIn * 1000
  localStorage.setItem('token_expires_at', String(expiresAt))
placeholder

/**
 * Get authentication token from localStorage
 */
export function getAuthToken(): string | null {
  return localStorage.getItem('auth_token')
placeholder

/**
 * Get refresh token from localStorage
 */
export function getRefreshToken(): string | null {
  return localStorage.getItem('refresh_token')
placeholder

/**
 * Get token expiration timestamp from localStorage
 */
export function getTokenExpiresAt(): number | null {
  const value = localStorage.getItem('token_expires_at')
  return value ? parseInt(value, 10) : null
placeholder

/**
 * Clear authentication token from localStorage
 */
export function clearAuthToken(): void {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('refresh_token')
  localStorage.removeItem('auth_user')
  localStorage.removeItem('token_expires_at')
placeholder

/**
 * User login
 * @param credentials - Email and password
 * @returns Authentication response with token and user data, or 2FA required response
 */
export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  const { data placeholder = await apiClient.post<LoginResponse>('/auth/login', credentials)

  // Only store token if 2FA is not required
  if (!isTotp2FARequired(data)) {
    setAuthToken(data.access_token)
    if (data.refresh_token) {
      setRefreshToken(data.refresh_token)
    placeholder
    if (data.expires_in) {
      setTokenExpiresAt(data.expires_in)
    placeholder
    localStorage.setItem('auth_user', JSON.stringify(data.user))
  placeholder

  return data
placeholder

/**
 * Complete login with 2FA code
 * @param request - Temp token and TOTP code
 * @returns Authentication response with token and user data
 */
export async function login2FA(request: TotpLogin2FARequest): Promise<AuthResponse> {
  const { data placeholder = await apiClient.post<AuthResponse>('/auth/login/2fa', request)

  // Store token and user data
  setAuthToken(data.access_token)
  if (data.refresh_token) {
    setRefreshToken(data.refresh_token)
  placeholder
  if (data.expires_in) {
    setTokenExpiresAt(data.expires_in)
  placeholder
  localStorage.setItem('auth_user', JSON.stringify(data.user))

  return data
placeholder

/**
 * User registration
 * @param userData - Registration data (username, email, password)
 * @returns Authentication response with token and user data
 */
export async function register(userData: RegisterRequest): Promise<AuthResponse> {
  const { data placeholder = await apiClient.post<AuthResponse>('/auth/register', userData)

  // Store token and user data
  setAuthToken(data.access_token)
  if (data.refresh_token) {
    setRefreshToken(data.refresh_token)
  placeholder
  if (data.expires_in) {
    setTokenExpiresAt(data.expires_in)
  placeholder
  localStorage.setItem('auth_user', JSON.stringify(data.user))

  return data
placeholder

/**
 * Get current authenticated user
 * @returns User profile data
 */
export async function getCurrentUser() {
  return apiClient.get<CurrentUserResponse>('/auth/me')
placeholder

/**
 * User logout
 * Clears authentication token and user data from localStorage
 * Optionally revokes the refresh token on the server
 */
export async function logout(): Promise<void> {
  const refreshToken = getRefreshToken()

  // Try to revoke the refresh token on the server
  if (refreshToken) {
    try {
      await apiClient.post('/auth/logout', { refresh_token: refreshToken placeholder)
    placeholder catch {
      // Ignore errors - we still want to clear local state
    placeholder
  placeholder

  clearAuthToken()
placeholder

/**
 * Refresh token response
 */
export interface OAuthTokenResponse {
  access_token: string
  refresh_token?: string
  expires_in?: number
  token_type?: string
placeholder

export interface PendingOAuthBindLoginResponse extends Partial<OAuthTokenResponse> {
  auth_result?: string
  redirect?: string
  error?: string
  requires_2fa?: boolean
  temp_token?: string
  user_email_masked?: string
  adoption_required?: boolean
  suggested_display_name?: string
  suggested_avatar_url?: string
placeholder

export type PendingOAuthExchangeResponse = PendingOAuthBindLoginResponse

export interface PendingOAuthCreateAccountResponse extends OAuthTokenResponse {
  auth_result?: string
placeholder

export interface PendingOAuthSendVerifyCodeResponse extends SendVerifyCodeResponse {
  auth_result?: string
  provider?: string
  redirect?: string
placeholder

export type OAuthCompletionKind = 'login' | 'bind'

export interface OAuthAdoptionDecision {
  adoptDisplayName?: boolean
  adoptAvatar?: boolean
placeholder

function serializeOAuthAdoptionDecision(
  decision?: OAuthAdoptionDecision
): Record<string, boolean> {
  const payload: Record<string, boolean> = {placeholder

  if (typeof decision?.adoptDisplayName === 'boolean') {
    payload.adopt_display_name = decision.adoptDisplayName
  placeholder
  if (typeof decision?.adoptAvatar === 'boolean') {
    payload.adopt_avatar = decision.adoptAvatar
  placeholder

  return payload
placeholder

export function isOAuthLoginCompletion(
  completion: Partial<OAuthTokenResponse>
): completion is OAuthTokenResponse {
  return typeof completion.access_token === 'string' && completion.access_token.trim().length > 0
placeholder

export function getOAuthCompletionKind(
  completion: Partial<OAuthTokenResponse>
): OAuthCompletionKind {
  return isOAuthLoginCompletion(completion) ? 'login' : 'bind'
placeholder

export function getPendingOAuthBindLoginKind(
  completion: PendingOAuthBindLoginResponse
): OAuthCompletionKind {
  return getOAuthCompletionKind(completion)
placeholder

export function isPendingOAuthCreateAccountRequired(
  completion: Pick<PendingOAuthBindLoginResponse, 'error'>
): boolean {
  return completion.error === 'invitation_required'
placeholder

export function hasPendingOAuthSuggestedProfile(
  completion: Pick<
    PendingOAuthBindLoginResponse,
    'suggested_display_name' | 'suggested_avatar_url'
  >
): boolean {
  return Boolean(completion.suggested_display_name || completion.suggested_avatar_url)
placeholder

export function persistOAuthTokenContext(tokens: Partial<OAuthTokenResponse>): void {
  if (tokens.refresh_token) {
    setRefreshToken(tokens.refresh_token)
  placeholder
  if (tokens.expires_in) {
    setTokenExpiresAt(tokens.expires_in)
  placeholder
placeholder

export async function prepareOAuthBindAccessTokenCookie(): Promise<void> {
  if (!getAuthToken()) {
    return
  placeholder
  await apiClient.post('/auth/oauth/bind-token')
placeholder

/**
 * Refresh the access token using the refresh token
 * @returns New token pair
 */
export async function refreshToken(): Promise<RefreshTokenResponse> {
  return refreshAuthTokens()
placeholder

/**
 * Revoke all sessions for the current user
 * @returns Response with message
 */
export async function revokeAllSessions(): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.post<{ message: string placeholder>('/auth/revoke-all-sessions')
  return data
placeholder

/**
 * Check if user is authenticated
 * @returns True if user has valid token
 */
export function isAuthenticated(): boolean {
  return getAuthToken() !== null
placeholder

/**
 * Get public settings (no auth required)
 * @returns Public settings including registration and Turnstile config
 */
export async function getPublicSettings(): Promise<PublicSettings> {
  const { data placeholder = await apiClient.get<PublicSettings>('/settings/public')
  return data
placeholder

export type WeChatOAuthMode = 'open' | 'mp'
export type WeChatOAuthUnavailableReason =
  | 'not_configured'
  | 'capability_unknown'
  | 'external_browser_required'
  | 'wechat_browser_required'
  | 'native_app_required'

export interface ResolvedWeChatOAuthStart {
  mode: WeChatOAuthMode | null
  openEnabled: boolean
  mpEnabled: boolean
  mobileEnabled: boolean
  isWeChatBrowser: boolean
  unavailableReason: WeChatOAuthUnavailableReason | null
placeholder

export type WeChatOAuthPublicSettings = {
  wechat_oauth_enabled?: boolean
  wechat_oauth_open_enabled?: boolean
  wechat_oauth_mp_enabled?: boolean
  wechat_oauth_mobile_enabled?: boolean
placeholder

export function isWeChatWebOAuthEnabled(
  settings: WeChatOAuthPublicSettings | null | undefined,
): boolean {
  const legacyEnabled = settings?.wechat_oauth_enabled ?? false
  const hasExplicitCapabilities =
    typeof settings?.wechat_oauth_open_enabled === 'boolean' ||
    typeof settings?.wechat_oauth_mp_enabled === 'boolean'

  if (!hasExplicitCapabilities) {
    return legacyEnabled
  placeholder

  return settings?.wechat_oauth_open_enabled === true || settings?.wechat_oauth_mp_enabled === true
placeholder

export function hasExplicitWeChatOAuthCapabilities(
  settings: WeChatOAuthPublicSettings | null | undefined,
): settings is WeChatOAuthPublicSettings & {
  wechat_oauth_open_enabled: boolean
  wechat_oauth_mp_enabled: boolean
placeholder {
  return typeof settings?.wechat_oauth_open_enabled === 'boolean'
    && typeof settings?.wechat_oauth_mp_enabled === 'boolean'
placeholder

export function resolveWeChatOAuthStart(
  settings: WeChatOAuthPublicSettings | null | undefined,
  userAgent?: string
): ResolvedWeChatOAuthStart {
  const normalizedUserAgent = (userAgent
    ?? (typeof navigator !== 'undefined' ? navigator.userAgent : '')
    ?? '').trim()
  const isWeChatBrowser = /MicroMessenger/i.test(normalizedUserAgent)
  const legacyEnabled = settings?.wechat_oauth_enabled ?? false
  const openEnabled = typeof settings?.wechat_oauth_open_enabled === 'boolean'
    ? settings.wechat_oauth_open_enabled
    : legacyEnabled
  const mpEnabled = typeof settings?.wechat_oauth_mp_enabled === 'boolean'
    ? settings.wechat_oauth_mp_enabled
    : legacyEnabled
  const mobileEnabled = typeof settings?.wechat_oauth_mobile_enabled === 'boolean'
    ? settings.wechat_oauth_mobile_enabled
    : false

  if (isWeChatBrowser) {
    if (mpEnabled) {
      return { mode: 'mp', openEnabled, mpEnabled, mobileEnabled, isWeChatBrowser, unavailableReason: null placeholder
    placeholder
    if (openEnabled) {
      return { mode: null, openEnabled, mpEnabled, mobileEnabled, isWeChatBrowser, unavailableReason: 'external_browser_required' placeholder
    placeholder
    return { mode: null, openEnabled, mpEnabled, mobileEnabled, isWeChatBrowser, unavailableReason: 'not_configured' placeholder
  placeholder

  if (openEnabled) {
    return { mode: 'open', openEnabled, mpEnabled, mobileEnabled, isWeChatBrowser, unavailableReason: null placeholder
  placeholder
  if (mpEnabled) {
    return { mode: null, openEnabled, mpEnabled, mobileEnabled, isWeChatBrowser, unavailableReason: 'wechat_browser_required' placeholder
  placeholder
  return { mode: null, openEnabled, mpEnabled, mobileEnabled, isWeChatBrowser, unavailableReason: 'not_configured' placeholder
placeholder

export function resolveWeChatOAuthStartStrict(
  settings: WeChatOAuthPublicSettings | null | undefined,
  userAgent?: string,
): ResolvedWeChatOAuthStart {
  const normalizedUserAgent = (userAgent
    ?? (typeof navigator !== 'undefined' ? navigator.userAgent : '')
    ?? '').trim()
  const isWeChatBrowser = /MicroMessenger/i.test(normalizedUserAgent)

  if (!hasExplicitWeChatOAuthCapabilities(settings)) {
    return {
      mode: null,
      openEnabled: false,
      mpEnabled: false,
      mobileEnabled: false,
      isWeChatBrowser,
      unavailableReason: 'capability_unknown',
    placeholder
  placeholder

  return resolveWeChatOAuthStart(settings, normalizedUserAgent)
placeholder

/**
 * Send verification code to email
 * @param request - Email and optional Turnstile token
 * @returns Response with countdown seconds
 */
export async function sendVerifyCode(
  request: SendVerifyCodeRequest
): Promise<SendVerifyCodeResponse> {
  const { data placeholder = await apiClient.post<SendVerifyCodeResponse>('/auth/send-verify-code', request)
  return data
placeholder

export async function sendPendingOAuthVerifyCode(
  request: SendVerifyCodeRequest
): Promise<PendingOAuthSendVerifyCodeResponse> {
  const { data placeholder = await apiClient.post<PendingOAuthSendVerifyCodeResponse>(
    '/auth/oauth/pending/send-verify-code',
    request
  )
  return data
placeholder

/**
 * Validate promo code response
 */
export interface ValidatePromoCodeResponse {
  valid: boolean
  bonus_amount?: number
  error_code?: string
  message?: string
placeholder

/**
 * Validate promo code (public endpoint, no auth required)
 * @param code - Promo code to validate
 * @returns Validation result with bonus amount if valid
 */
export async function validatePromoCode(code: string): Promise<ValidatePromoCodeResponse> {
  const { data placeholder = await apiClient.post<ValidatePromoCodeResponse>('/auth/validate-promo-code', { code placeholder)
  return data
placeholder

/**
 * Validate invitation code response
 */
export interface ValidateInvitationCodeResponse {
  valid: boolean
  error_code?: string
placeholder

/**
 * Validate invitation code (public endpoint, no auth required)
 * @param code - Invitation code to validate
 * @returns Validation result
 */
export async function validateInvitationCode(code: string): Promise<ValidateInvitationCodeResponse> {
  const { data placeholder = await apiClient.post<ValidateInvitationCodeResponse>('/auth/validate-invitation-code', { code placeholder)
  return data
placeholder

/**
 * Forgot password request
 */
export interface ForgotPasswordRequest {
  email: string
  turnstile_token?: string
placeholder

/**
 * Forgot password response
 */
export interface ForgotPasswordResponse {
  message: string
placeholder

/**
 * Request password reset link
 * @param request - Email and optional Turnstile token
 * @returns Response with message
 */
export async function forgotPassword(request: ForgotPasswordRequest): Promise<ForgotPasswordResponse> {
  const { data placeholder = await apiClient.post<ForgotPasswordResponse>('/auth/forgot-password', request)
  return data
placeholder

/**
 * Reset password request
 */
export interface ResetPasswordRequest {
  email: string
  token: string
  new_password: string
placeholder

/**
 * Reset password response
 */
export interface ResetPasswordResponse {
  message: string
placeholder

/**
 * Reset password with token
 * @param request - Email, token, and new password
 * @returns Response with message
 */
export async function resetPassword(request: ResetPasswordRequest): Promise<ResetPasswordResponse> {
  const { data placeholder = await apiClient.post<ResetPasswordResponse>('/auth/reset-password', request)
  return data
placeholder

/**
 * Complete LinuxDo OAuth registration by supplying an invitation code
 * @param invitationCode - Invitation code entered by the user
 * @returns Token pair on success
 */
export async function completeLinuxDoOAuthRegistration(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<OAuthTokenResponse> {
  return createPendingLinuxDoOAuthAccount(invitationCode, decision, affiliateCode)
placeholder

/**
 * Complete OIDC OAuth registration by supplying an invitation code
 * @param invitationCode - Invitation code entered by the user
 * @returns Token pair on success
 */
export async function completeOIDCOAuthRegistration(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<OAuthTokenResponse> {
  return createPendingOIDCOAuthAccount(invitationCode, decision, affiliateCode)
placeholder

export async function completeWeChatOAuthRegistration(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<OAuthTokenResponse> {
  return createPendingWeChatOAuthAccount(invitationCode, decision, affiliateCode)
placeholder

async function createPendingOAuthAccount(
  provider: 'linuxdo' | 'oidc' | 'wechat' | 'dingtalk',
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<PendingOAuthCreateAccountResponse> {
  const normalizedAffiliateCode = affiliateCode?.trim()
  const { data placeholder = await apiClient.post<PendingOAuthCreateAccountResponse>(
    `/auth/oauth/${providerplaceholder/complete-registration`,
    {
      invitation_code: invitationCode,
      ...(normalizedAffiliateCode ? { aff_code: normalizedAffiliateCode placeholder : {placeholder),
      ...serializeOAuthAdoptionDecision(decision)
    placeholder
  )
  return data
placeholder

export async function createPendingLinuxDoOAuthAccount(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<PendingOAuthCreateAccountResponse> {
  return createPendingOAuthAccount('linuxdo', invitationCode, decision, affiliateCode)
placeholder

export async function createPendingOIDCOAuthAccount(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<PendingOAuthCreateAccountResponse> {
  return createPendingOAuthAccount('oidc', invitationCode, decision, affiliateCode)
placeholder

export async function createPendingWeChatOAuthAccount(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<PendingOAuthCreateAccountResponse> {
  return createPendingOAuthAccount('wechat', invitationCode, decision, affiliateCode)
placeholder

export async function createPendingDingTalkOAuthAccount(
  invitationCode: string,
  decision?: OAuthAdoptionDecision,
  affiliateCode?: string
): Promise<PendingOAuthCreateAccountResponse> {
  return createPendingOAuthAccount('dingtalk', invitationCode, decision, affiliateCode)
placeholder

export async function completePendingOAuthBindLogin(
  decision?: OAuthAdoptionDecision
): Promise<PendingOAuthBindLoginResponse> {
  const { data placeholder = await apiClient.post<PendingOAuthBindLoginResponse>(
    '/auth/oauth/pending/exchange',
    serializeOAuthAdoptionDecision(decision)
  )
  return data
placeholder

export async function exchangePendingOAuthCompletion(
  decision?: OAuthAdoptionDecision
): Promise<PendingOAuthExchangeResponse> {
  return completePendingOAuthBindLogin(decision)
placeholder

export const authAPI = {
  login,
  login2FA,
  isTotp2FARequired,
  register,
  getCurrentUser,
  logout,
  isAuthenticated,
  setAuthToken,
  setRefreshToken,
  setTokenExpiresAt,
  getAuthToken,
  getRefreshToken,
  getTokenExpiresAt,
  clearAuthToken,
  getPublicSettings,
  sendVerifyCode,
  sendPendingOAuthVerifyCode,
  validatePromoCode,
  validateInvitationCode,
  forgotPassword,
  resetPassword,
  refreshToken,
  revokeAllSessions,
  getPendingOAuthBindLoginKind,
  isPendingOAuthCreateAccountRequired,
  hasPendingOAuthSuggestedProfile,
  completePendingOAuthBindLogin,
  createPendingLinuxDoOAuthAccount,
  createPendingOIDCOAuthAccount,
  createPendingWeChatOAuthAccount,
  exchangePendingOAuthCompletion,
  completeLinuxDoOAuthRegistration,
  completeOIDCOAuthRegistration,
  completeWeChatOAuthRegistration,
  createPendingDingTalkOAuthAccount
placeholder

export default authAPI
