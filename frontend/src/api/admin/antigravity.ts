/**
 * Admin Antigravity API endpoints
 * Handles Antigravity (Google Cloud AI Companion) OAuth flows for administrators
 */

import { apiClient placeholder from '../client'

export interface AntigravityAuthUrlResponse {
  auth_url: string
  session_id: string
  state: string
placeholder

export interface AntigravityAuthUrlRequest {
  proxy_id?: number
placeholder

export interface AntigravityExchangeCodeRequest {
  session_id: string
  state: string
  code: string
  proxy_id?: number
placeholder

export interface AntigravityTokenInfo {
  access_token?: string
  refresh_token?: string
  token_type?: string
  expires_at?: number | string
  expires_in?: number
  project_id?: string
  email?: string
  [key: string]: unknown
placeholder

export async function generateAuthUrl(
  payload: AntigravityAuthUrlRequest
): Promise<AntigravityAuthUrlResponse> {
  const { data placeholder = await apiClient.post<AntigravityAuthUrlResponse>(
    '/admin/antigravity/oauth/auth-url',
    payload
  )
  return data
placeholder

export async function exchangeCode(
  payload: AntigravityExchangeCodeRequest
): Promise<AntigravityTokenInfo> {
  const { data placeholder = await apiClient.post<AntigravityTokenInfo>(
    '/admin/antigravity/oauth/exchange-code',
    payload
  )
  return data
placeholder

export default { generateAuthUrl, exchangeCode placeholder
