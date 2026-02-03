/**
 * TOTP (2FA) API endpoints
 * Handles Two-Factor Authentication with Google Authenticator
 */

import { apiClient placeholder from './client'
import type {
  TotpStatus,
  TotpSetupRequest,
  TotpSetupResponse,
  TotpEnableRequest,
  TotpEnableResponse,
  TotpDisableRequest,
  TotpVerificationMethod
placeholder from '@/types'

/**
 * Get TOTP status for current user
 * @returns TOTP status including enabled state and feature availability
 */
export async function getStatus(): Promise<TotpStatus> {
  const { data placeholder = await apiClient.get<TotpStatus>('/user/totp/status')
  return data
placeholder

/**
 * Get verification method for TOTP operations
 * @returns Method ('email' or 'password') required for setup/disable
 */
export async function getVerificationMethod(): Promise<TotpVerificationMethod> {
  const { data placeholder = await apiClient.get<TotpVerificationMethod>('/user/totp/verification-method')
  return data
placeholder

/**
 * Send email verification code for TOTP operations
 * @returns Success response
 */
export async function sendVerifyCode(): Promise<{ success: boolean placeholder> {
  const { data placeholder = await apiClient.post<{ success: boolean placeholder>('/user/totp/send-code')
  return data
placeholder

/**
 * Initiate TOTP setup - generates secret and QR code
 * @param request - Email code or password depending on verification method
 * @returns Setup response with secret, QR code URL, and setup token
 */
export async function initiateSetup(request?: TotpSetupRequest): Promise<TotpSetupResponse> {
  const { data placeholder = await apiClient.post<TotpSetupResponse>('/user/totp/setup', request || {placeholder)
  return data
placeholder

/**
 * Complete TOTP setup by verifying the code
 * @param request - TOTP code and setup token
 * @returns Enable response with success status and enabled timestamp
 */
export async function enable(request: TotpEnableRequest): Promise<TotpEnableResponse> {
  const { data placeholder = await apiClient.post<TotpEnableResponse>('/user/totp/enable', request)
  return data
placeholder

/**
 * Disable TOTP for current user
 * @param request - Email code or password depending on verification method
 * @returns Success response
 */
export async function disable(request: TotpDisableRequest): Promise<{ success: boolean placeholder> {
  const { data placeholder = await apiClient.post<{ success: boolean placeholder>('/user/totp/disable', request)
  return data
placeholder

export const totpAPI = {
  getStatus,
  getVerificationMethod,
  sendVerifyCode,
  initiateSetup,
  enable,
  disable
placeholder

export default totpAPI
