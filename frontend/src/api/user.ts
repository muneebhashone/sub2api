/**
 * User API endpoints
 * Handles user profile management and password changes
 */

import { apiClient placeholder from './client'
import { prepareOAuthBindAccessTokenCookie placeholder from './auth'
import type { User, ChangePasswordRequest, NotifyEmailEntry, UserAuthProvider placeholder from '@/types'

/**
 * Get current user profile
 * @returns User profile data
 */
export async function getProfile(): Promise<User> {
  const { data placeholder = await apiClient.get<User>('/user/profile')
  return data
placeholder

/**
 * Update current user profile
 * @param profile - Profile data to update
 * @returns Updated user profile data
 */
export async function updateProfile(profile: {
  username?: string
  balance_notify_enabled?: boolean
  balance_notify_threshold?: number | null
  balance_notify_extra_emails?: NotifyEmailEntry[]
placeholder): Promise<User> {
  const { data placeholder = await apiClient.put<User>('/user', profile)
  return data
placeholder

/**
 * Change current user password
 * @param passwords - Old and new password
 * @returns Success message
 */
export async function changePassword(
  oldPassword: string,
  newPassword: string
): Promise<{ message: string placeholder> {
  const payload: ChangePasswordRequest = {
    old_password: oldPassword,
    new_password: newPassword
  placeholder

  const { data placeholder = await apiClient.put<{ message: string placeholder>('/user/password', payload)
  return data
placeholder

/**
 * Send verification code for adding a notify email
 * @param email - Email address to verify
 */
export async function sendNotifyEmailCode(email: string): Promise<void> {
  await apiClient.post('/user/notify-email/send-code', { email placeholder)
placeholder

/**
 * Verify and add a notify email
 * @param email - Email address to add
 * @param code - Verification code
 */
export async function verifyNotifyEmail(email: string, code: string): Promise<void> {
  await apiClient.post('/user/notify-email/verify', { email, code placeholder)
placeholder

/**
 * Remove a notify email
 * @param email - Email address to remove
 */
export async function removeNotifyEmail(email: string): Promise<void> {
  await apiClient.delete('/user/notify-email', { data: { email placeholder placeholder)
placeholder

/**
 * Toggle a notify email's disabled state
 * @param email - Email address (empty string for primary email placeholder)
 * @param disabled - Whether to disable the email
 */
export async function toggleNotifyEmail(email: string, disabled: boolean): Promise<User> {
  const { data placeholder = await apiClient.put<User>('/user/notify-email/toggle', { email, disabled placeholder)
  return data
placeholder

export type BindableOAuthProvider = Exclude<UserAuthProvider, 'email'>

interface BuildOAuthBindingStartURLOptions {
  redirectTo?: string
placeholder

export function resolveWeChatOAuthMode(): 'open' | 'mp' {
  if (typeof navigator === 'undefined') {
    return 'open'
  placeholder
  return /MicroMessenger/i.test(navigator.userAgent) ? 'mp' : 'open'
placeholder

export function buildOAuthBindingStartURL(
  provider: BindableOAuthProvider,
  options: BuildOAuthBindingStartURLOptions = {placeholder
): string {
  const redirectTo = options.redirectTo?.trim() || '/profile'
  const apiBase = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api/v1'
  const normalized = apiBase.replace(/\/$/, '')
  const params = new URLSearchParams({
    redirect: redirectTo,
    intent: 'bind_current_user'
  placeholder)

  if (provider === 'wechat') {
    params.set('mode', resolveWeChatOAuthMode())
  placeholder

  return `${normalizedplaceholder/auth/oauth/${providerplaceholder/start?${params.toString()placeholder`
placeholder

export function startOAuthBinding(
  provider: BindableOAuthProvider,
  options: BuildOAuthBindingStartURLOptions = {placeholder
): void {
  if (typeof window === 'undefined') {
    return
  placeholder
  prepareOAuthBindAccessTokenCookie()
  window.location.href = buildOAuthBindingStartURL(provider, options)
placeholder

export const userAPI = {
  getProfile,
  updateProfile,
  changePassword,
  sendNotifyEmailCode,
  verifyNotifyEmail,
  removeNotifyEmail,
  toggleNotifyEmail,
  buildOAuthBindingStartURL,
  startOAuthBinding
placeholder

export default userAPI
