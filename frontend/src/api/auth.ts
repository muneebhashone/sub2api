/**
 * Authentication API endpoints
 * Handles user login, registration, and logout operations
 */

import { apiClient placeholder from './client'
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
 * Get authentication token from localStorage
 */
export function getAuthToken(): string | null {
  return localStorage.getItem('auth_token')
placeholder

/**
 * Clear authentication token from localStorage
 */
export function clearAuthToken(): void {
  localStorage.removeItem('auth_token')
  localStorage.removeItem('auth_user')
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
 */
export function logout(): void {
  clearAuthToken()
  // Optionally redirect to login page
  // window.location.href = '/login';
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

export const authAPI = {
  login,
  login2FA,
  isTotp2FARequired,
  register,
  getCurrentUser,
  logout,
  isAuthenticated,
  setAuthToken,
  getAuthToken,
  clearAuthToken,
  getPublicSettings,
  sendVerifyCode,
  validatePromoCode,
  forgotPassword,
  resetPassword
placeholder

export default authAPI
