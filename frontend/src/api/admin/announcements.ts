/**
 * Admin Announcements API endpoints
 */

import { apiClient placeholder from '../client'
import type {
  Announcement,
  AnnouncementUserReadStatus,
  BasePaginationResponse,
  CreateAnnouncementRequest,
  UpdateAnnouncementRequest
placeholder from '@/types'

export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    status?: string
    search?: string
  placeholder
): Promise<BasePaginationResponse<Announcement>> {
  const { data placeholder = await apiClient.get<BasePaginationResponse<Announcement>>('/admin/announcements', {
    params: { page, page_size: pageSize, ...filters placeholder
  placeholder)
  return data
placeholder

export async function getById(id: number): Promise<Announcement> {
  const { data placeholder = await apiClient.get<Announcement>(`/admin/announcements/${idplaceholder`)
  return data
placeholder

export async function create(request: CreateAnnouncementRequest): Promise<Announcement> {
  const { data placeholder = await apiClient.post<Announcement>('/admin/announcements', request)
  return data
placeholder

export async function update(id: number, request: UpdateAnnouncementRequest): Promise<Announcement> {
  const { data placeholder = await apiClient.put<Announcement>(`/admin/announcements/${idplaceholder`, request)
  return data
placeholder

export async function deleteAnnouncement(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.delete<{ message: string placeholder>(`/admin/announcements/${idplaceholder`)
  return data
placeholder

export async function getReadStatus(
  id: number,
  page: number = 1,
  pageSize: number = 20,
  search: string = ''
): Promise<BasePaginationResponse<AnnouncementUserReadStatus>> {
  const { data placeholder = await apiClient.get<BasePaginationResponse<AnnouncementUserReadStatus>>(
    `/admin/announcements/${idplaceholder/read-status`,
    { params: { page, page_size: pageSize, search placeholder placeholder
  )
  return data
placeholder

const announcementsAPI = {
  list,
  getById,
  create,
  update,
  delete: deleteAnnouncement,
  getReadStatus
placeholder

export default announcementsAPI

