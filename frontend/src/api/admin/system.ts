/**
 * System API endpoints for admin operations
 */

import { apiClient placeholder from '../client'

export interface ReleaseInfo {
  name: string
  body: string
  published_at: string
  html_url: string
placeholder

export interface VersionInfo {
  current_version: string
  latest_version: string
  has_update: boolean
  release_info?: ReleaseInfo
  cached: boolean
  warning?: string
  build_type: string // "source" for manual builds, "release" for CI builds
placeholder

/**
 * Get current version
 */
export async function getVersion(): Promise<{ version: string placeholder> {
  const { data placeholder = await apiClient.get<{ version: string placeholder>('/admin/system/version')
  return data
placeholder

/**
 * Check for updates
 * @param force - Force refresh from GitHub API
 */
export async function checkUpdates(force = false): Promise<VersionInfo> {
  const { data placeholder = await apiClient.get<VersionInfo>('/admin/system/check-updates', {
    params: force ? { force: 'true' placeholder : undefined
  placeholder)
  return data
placeholder

export interface UpdateResult {
  message: string
  need_restart: boolean
placeholder

/**
 * Perform system update
 * Downloads and applies the latest version
 */
export async function performUpdate(): Promise<UpdateResult> {
  const { data placeholder = await apiClient.post<UpdateResult>('/admin/system/update')
  return data
placeholder

/**
 * Rollback to previous version
 */
export async function rollback(): Promise<UpdateResult> {
  const { data placeholder = await apiClient.post<UpdateResult>('/admin/system/rollback')
  return data
placeholder

/**
 * Restart the service
 */
export async function restartService(): Promise<{ message: string placeholder> {
  const { data placeholder = await apiClient.post<{ message: string placeholder>('/admin/system/restart')
  return data
placeholder

export const systemAPI = {
  getVersion,
  checkUpdates,
  performUpdate,
  rollback,
  restartService
placeholder

export default systemAPI
