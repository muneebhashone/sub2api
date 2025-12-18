/**
 * User API endpoints
 * Handles user profile management and password changes
 */

import { apiClient placeholder from './client';
import type { User, ChangePasswordRequest placeholder from '@/types';

/**
 * Get current user profile
 * @returns User profile data
 */
export async function getProfile(): Promise<User> {
  const { data placeholder = await apiClient.get<User>('/users/me');
  return data;
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
    new_password: newPassword,
  placeholder;

  const { data placeholder = await apiClient.post<{ message: string placeholder>('/users/me/password', payload);
  return data;
placeholder

export const userAPI = {
  getProfile,
  changePassword,
placeholder;

export default userAPI;
