import { apiClient placeholder from '@/api/client'

export interface AdminComplianceAcknowledgement {
  version: string
  document_zh: string
  document_en: string
  admin_user_id: number
  ip_address?: string
  user_agent?: string
  accepted_at: string
placeholder

export interface AdminComplianceStatus {
  required: boolean
  version: string
  document_path_zh: string
  document_path_en: string
  document_url_zh: string
  document_url_en: string
  ack_phrase_zh: string
  ack_phrase_en: string
  acknowledgement?: AdminComplianceAcknowledgement
placeholder

export interface AcceptAdminComplianceRequest {
  phrase: string
  language: string
placeholder

export const adminComplianceAPI = {
  async getStatus(): Promise<AdminComplianceStatus> {
    const { data placeholder = await apiClient.get<AdminComplianceStatus>('/admin/compliance')
    return data
  placeholder,

  async accept(payload: AcceptAdminComplianceRequest): Promise<AdminComplianceStatus> {
    const { data placeholder = await apiClient.post<AdminComplianceStatus>('/admin/compliance/accept', payload)
    return data
  placeholder
placeholder

export default adminComplianceAPI
