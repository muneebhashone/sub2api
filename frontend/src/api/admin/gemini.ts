/**
 * Admin Gemini API endpoints
 * Handles Gemini OAuth flows for administrators
 */

import { apiClient placeholder from '../client'

export interface GeminiAuthUrlResponse {
  auth_url: string
  session_id: string
  state: string
placeholder

export interface GeminiOAuthCapabilities {
  ai_studio_oauth_enabled: boolean
  required_redirect_uris: string[]
placeholder

export interface GeminiAuthUrlRequest {
  proxy_id?: number
  project_id?: string
  oauth_type?: 'code_assist' | 'google_one' | 'ai_studio'
  tier_id?: string
placeholder

export interface GeminiExchangeCodeRequest {
  session_id: string
  state: string
  code: string
  proxy_id?: number
  oauth_type?: 'code_assist' | 'google_one' | 'ai_studio'
  tier_id?: string
placeholder

export type GeminiTokenInfo = {
  access_token?: string
  refresh_token?: string
  token_type?: string
  scope?: string
  expires_in?: number
  expires_at?: number
  project_id?: string
  oauth_type?: string
  tier_id?: string
  extra?: Record<string, unknown>
  [key: string]: unknown
placeholder

export async function generateAuthUrl(
  payload: GeminiAuthUrlRequest
): Promise<GeminiAuthUrlResponse> {
  const { data placeholder = await apiClient.post<GeminiAuthUrlResponse>(
    '/admin/gemini/oauth/auth-url',
    payload
  )
  return data
placeholder

export async function exchangeCode(payload: GeminiExchangeCodeRequest): Promise<GeminiTokenInfo> {
  const { data placeholder = await apiClient.post<GeminiTokenInfo>(
    '/admin/gemini/oauth/exchange-code',
    payload
  )
  return data
placeholder

export async function getCapabilities(): Promise<GeminiOAuthCapabilities> {
  const { data placeholder = await apiClient.get<GeminiOAuthCapabilities>('/admin/gemini/oauth/capabilities')
  return data
placeholder

export default { generateAuthUrl, exchangeCode, getCapabilities placeholder
