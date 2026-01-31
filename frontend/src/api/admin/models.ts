import { apiClient placeholder from '@/api/client'

export async function getPlatformModels(platform: string): Promise<string[]> {
  const { data placeholder = await apiClient.get<string[]>('/admin/models', {
    params: { platform placeholder
  placeholder)
  return data
placeholder

export const modelsAPI = {
  getPlatformModels
placeholder

export default modelsAPI
