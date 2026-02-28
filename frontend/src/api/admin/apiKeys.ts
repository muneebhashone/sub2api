/**
 * Admin API Keys API endpoints
 * Handles API key management for administrators
 */

import { apiClient placeholder from '../client'
import type { ApiKey placeholder from '@/types'

export interface UpdateApiKeyGroupResult {
  api_key: ApiKey
  auto_granted_group_access: boolean
  granted_group_id?: number
  granted_group_name?: string
placeholder

/**
 * Update an API key's group binding
 * @param id - API Key ID
 * @param groupId - Group ID (0 to unbind, positive to bind, null/undefined to skip)
 * @returns Updated API key with auto-grant info
 */
export async function updateApiKeyGroup(id: number, groupId: number | null): Promise<UpdateApiKeyGroupResult> {
  const { data placeholder = await apiClient.put<UpdateApiKeyGroupResult>(`/admin/api-keys/${idplaceholder`, {
    group_id: groupId === null ? 0 : groupId
  placeholder)
  return data
placeholder

export const apiKeysAPI = {
  updateApiKeyGroup
placeholder

export default apiKeysAPI
