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
  PublicSettings
placeholder from '@/types'

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
 * @param credentials - Username and password
 * @returns Authentication response with token and user data
 */
export async function login(credentials: LoginRequest): Promise<AuthResponse> {
  const { data placeholder = await apiClient.post<AuthResponse>('/auth/login', credentials)

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

export const authAPI = {
  login,
  register,
  getCurrentUser,
  logout,
  isAuthenticated,
  setAuthToken,
  getAuthToken,
  clearAuthToken,
  getPublicSettings,
  sendVerifyCode,
  validatePromoCode
placeholder

export default authAPI
