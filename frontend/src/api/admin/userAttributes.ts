/**
 * Admin User Attributes API endpoints
 * Handles user custom attribute definitions and values
 */

import { apiClient placeholder from '../client'
import type {
  UserAttributeDefinition,
  UserAttributeValue,
  CreateUserAttributeRequest,
  UpdateUserAttributeRequest,
  UserAttributeValuesMap
placeholder from '@/types'

/**
 * Get all attribute definitions
 */
export async function listDefinitions(): Promise<UserAttributeDefinition[]> {
  const { data placeholder = await apiClient.get<UserAttributeDefinition[]>('/admin/user-attributes')
  return data
placeholder

/**
 * Get enabled attribute definitions only
 */
export async function listEnabledDefinitions(): Promise<UserAttributeDefinition[]> {
  const { data placeholder = await apiClient.get<UserAttributeDefinition[]>('/admin/user-attributes', {
    params: { enabled: true placeholder
  placeholder)
  return data
placeholder

/**
 * Create a new attribute definition
 */
export async function createDefinition(
  request: CreateUserAttributeRequest
): Promise<UserAttributeDefinition> {
  const { data placeholder = await apiClient.post<UserAttributeDefinition>('/admin/user-attributes', request)
  return data
placeholder

/**
 * Update an attribute definition
 */
export async function updateDefinition(
  id: number,
  request: UpdateUserAttributeRequest
): Promise<UserAttributeDefinition> {
  const { data placeholder = await apiClient.put<UserAttributeDefinition>(
    `/admin/user-attributes/${idplaceholder`,
    request
  )
  return data
placeholder

/**
 * Delete an attribute definition
 */
export async function deleteDefinition(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/user-attributes/${idplaceholder`)
  return data
placeholder

/**
 * Reorder attribute definitions
 */
export async function reorderDefinitions(ids: number[]): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.put<{ message: string placeholder>('/admin/user-attributes/reorder', {
    ids
  placeholder)
  return data
placeholder

/**
 * Get user's attribute values
 */
export async function getUserAttributeValues(userId: number): Promise<UserAttributeValue[]> {
  const { data placeholder = await apiClient.get<UserAttributeValue[]>(
    `/admin/users/${userIdplaceholder/attributes`
  )
  return data
placeholder

/**
 * Update user's attribute values (batch)
 */
export async function updateUserAttributeValues(
  userId: number,
  values: UserAttributeValuesMap
): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.put<{ message: string placeholder>(
    `/admin/users/${userIdplaceholder/attributes`,
    { values placeholder
  )
  return data
placeholder

/**
 * Batch response type
 */
export interface BatchUserAttributesResponse {
  attributes: Record<number, Record<number, string>>
placeholder

/**
 * Get attribute values for multiple users
 */
export async function getBatchUserAttributes(
  userIds: number[]
): Promise<BatchUserAttributesResponse> {
  const { data placeholder = await apiClient.post<BatchUserAttributesResponse>(
    '/admin/user-attributes/batch',
    { user_ids: userIds placeholder
  )
  return data
placeholder

export const userAttributesAPI = {
  listDefinitions,
  listEnabledDefinitions,
  createDefinition,
  updateDefinition,
  deleteDefinition,
  reorderDefinitions,
  getUserAttributeValues,
  updateUserAttributeValues,
  getBatchUserAttributes
placeholder

export default userAttributesAPI
