/**
 * Admin Modules API endpoints (plugin module observability, read-only)
 */

import { apiClient placeholder from '../client'
import type { Module placeholder from '@/types'

export async function list(): Promise<{ modules: Module[] placeholder> {
  const { data placeholder = await apiClient.get<{ modules: Module[] placeholder>('/admin/modules')
  return data
placeholder

const modulesAPI = {
  list
placeholder

export default modulesAPI
