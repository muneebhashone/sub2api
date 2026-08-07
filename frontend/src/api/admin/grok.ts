/**
 * Admin Grok/xAI API endpoints
 * Handles xAI OAuth flows for administrators.
 */

import { apiClient placeholder from '../client'
import type { GrokBillingSummary, GrokQuotaWindow, WindowStats placeholder from '@/types'

export type { GrokBillingSummary, GrokQuotaWindow placeholder from '@/types'

export interface GrokAuthUrlResponse {
  auth_url: string
  session_id: string
  state: string
placeholder

export interface GrokAuthUrlRequest {
  proxy_id?: number
  redirect_uri?: string
placeholder

export interface GrokExchangeCodeRequest {
  session_id: string
  state: string
  code: string
  proxy_id?: number
  redirect_uri?: string
placeholder

export interface GrokTokenInfo {
  access_token?: string
  refresh_token?: string
  token_type?: string
  id_token?: string
  expires_at?: number | string
  expires_in?: number
  scope?: string
  client_id?: string
  email?: string
  sub?: string
  team_id?: string
  subscription_tier?: string
  entitlement_status?: string
  [key: string]: unknown
placeholder

export interface GrokSSOToOAuthRequest {
  sso_tokens: string[]
  name?: string
  notes?: string | null
  proxy_id?: number | null
  group_ids?: number[]
  credentials?: Record<string, unknown>
  extra?: Record<string, unknown>
  concurrency?: number
  load_factor?: number
  priority?: number
  rate_multiplier?: number
  expires_at?: number | null
  auto_pause_on_expired?: boolean
placeholder

export interface GrokSSOToOAuthItemResult {
  index: number
  name?: string
  email?: string
  account?: unknown
  error?: string
placeholder

export interface GrokSSOToOAuthResponse {
  created: GrokSSOToOAuthItemResult[]
  failed: GrokSSOToOAuthItemResult[]
placeholder

const GROK_SSO_IMPORT_CONCURRENCY = 3
const GROK_SSO_IMPORT_TIMEOUT_PER_BATCH_MS = 90_000
const GROK_SSO_IMPORT_TIMEOUT_BUFFER_MS = 90_000

export function getGrokSSOImportTimeout(keyCount: number): number {
  const batches = Math.ceil(Math.max(1, keyCount) / GROK_SSO_IMPORT_CONCURRENCY)
  return batches * GROK_SSO_IMPORT_TIMEOUT_PER_BATCH_MS + GROK_SSO_IMPORT_TIMEOUT_BUFFER_MS
placeholder

export interface GrokQuotaSnapshot {
  requests?: GrokQuotaWindow | null
  tokens?: GrokQuotaWindow | null
  retry_after_seconds?: number | null
  subscription_tier?: string
  entitlement_status?: string
  status_code?: number
  headers?: Record<string, string>
  headers_observed: boolean
  observation_source?: string
  last_probe_at?: string
  last_headers_seen_at?: string
  updated_at: string
placeholder

export interface GrokQuotaProbeResult {
  source: 'active_probe' | 'billing_probe' | 'hybrid_probe'
  model?: string
  billing?: GrokBillingSummary | null
  snapshot?: GrokQuotaSnapshot | null
  local_usage_24h?: WindowStats | null
  local_usage_7d?: WindowStats | null
  local_usage_monthly?: WindowStats | null
  status_code?: number
  headers_observed: boolean
  reset_supported: boolean
  fetched_at: number
  persisted?: boolean
  probe_error?: string
placeholder

export interface GrokQuotaResetResult {
  supported: boolean
  code: string
  message: string
placeholder

export async function generateAuthUrl(
  payload: GrokAuthUrlRequest
): Promise<GrokAuthUrlResponse> {
  const { data placeholder = await apiClient.post<GrokAuthUrlResponse>(
    '/admin/grok/oauth/auth-url',
    payload
  )
  return data
placeholder

export async function exchangeCode(payload: GrokExchangeCodeRequest): Promise<GrokTokenInfo> {
  const { data placeholder = await apiClient.post<GrokTokenInfo>(
    '/admin/grok/oauth/exchange-code',
    payload
  )
  return data
placeholder

export async function refreshGrokToken(
  refreshToken: string,
  proxyId?: number | null
): Promise<GrokTokenInfo> {
  const payload: Record<string, unknown> = { refresh_token: refreshToken placeholder
  if (proxyId) payload.proxy_id = proxyId

  const { data placeholder = await apiClient.post<GrokTokenInfo>(
    '/admin/grok/oauth/refresh-token',
    payload
  )
  return data
placeholder

export async function queryQuota(id: number): Promise<GrokQuotaProbeResult> {
  const { data placeholder = await apiClient.get<GrokQuotaProbeResult>(`/admin/grok/accounts/${idplaceholder/quota`)
  return data
placeholder

export async function resetQuota(id: number): Promise<GrokQuotaResetResult> {
  const { data placeholder = await apiClient.post<GrokQuotaResetResult>(`/admin/grok/accounts/${idplaceholder/reset-quota`)
  return data
placeholder

export async function createFromSSO(payload: GrokSSOToOAuthRequest): Promise<GrokSSOToOAuthResponse> {
  const { data placeholder = await apiClient.post<GrokSSOToOAuthResponse>(
    '/admin/grok/sso-to-oauth',
    payload,
    { timeout: getGrokSSOImportTimeout(payload.sso_tokens.length) placeholder
  )
  return data
placeholder

/** Validate a browser SSO cookie and convert to Build OAuth tokens (no raw SSO stored). */
export async function validateSSOToken(
  ssoToken: string,
  proxyId?: number | null
): Promise<GrokTokenInfo> {
  const payload: Record<string, unknown> = { sso_token: ssoToken placeholder
  if (proxyId) payload.proxy_id = proxyId
  const { data placeholder = await apiClient.post<GrokTokenInfo>('/admin/grok/oauth/sso-token', payload)
  return data
placeholder

/**
 * Password login → ephemeral SSO → Build OAuth.
 * Password is only sent over the wire for this call; never persist it in credentials.
 */
export async function authorizePassword(
  emailAndPassword: string,
  proxyId?: number | null
): Promise<GrokTokenInfo> {
  // Format: email----password (password may contain dashes).
  const sep = '----'
  const idx = emailAndPassword.indexOf(sep)
  const email = (idx >= 0 ? emailAndPassword.slice(0, idx) : emailAndPassword).trim()
  const password = idx >= 0 ? emailAndPassword.slice(idx + sep.length) : ''
  const payload: Record<string, unknown> = { email, password placeholder
  if (proxyId) payload.proxy_id = proxyId
  const { data placeholder = await apiClient.post<GrokTokenInfo>('/admin/grok/oauth/password', payload)
  return data
placeholder

export default {
  generateAuthUrl,
  exchangeCode,
  refreshGrokToken,
  queryQuota,
  resetQuota,
  createFromSSO,
  validateSSOToken,
  authorizePassword,
placeholder
