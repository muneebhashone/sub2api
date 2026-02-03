/**
 * User Announcements API endpoints
 */

import { apiClient placeholder from './client'
import type { UserAnnouncement placeholder from '@/types'

export async function list(unreadOnly: boolean = false): Promise<UserAnnouncement[]> {
  const { data placeholder = await apiClient.get<UserAnnouncement[]>('/announcements', {
    params: unreadOnly ? { unread_only: 1 placeholder : {placeholder
  placeholder)
  return data
placeholder

export async function markRead(id: number): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.post<{ message: string placeholder>(`/announcements/${idplaceholder/read`)
  return data
placeholder

const announcementsAPI = {
  list,
  markRead
placeholder

export default announcementsAPI

