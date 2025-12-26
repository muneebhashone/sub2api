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

export interface GeminiAuthUrlRequest {
  proxy_id?: number
  project_id?: string
  oauth_type?: 'code_assist' | 'ai_studio'
placeholder

export interface GeminiExchangeCodeRequest {
  session_id: string
  state: string
  code: string
  proxy_id?: number
  oauth_type?: 'code_assist' | 'ai_studio'
placeholder

export type GeminiTokenInfo = Record<string, unknown>

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

export default { generateAuthUrl, exchangeCode placeholder
