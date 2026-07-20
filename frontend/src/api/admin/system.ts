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

export interface RollbackVersionInfo {
  version: string
  published_at: string
  html_url: string
placeholder

/**
 * Get versions available for rollback (up to 3 versions older than current)
 */
export async function getRollbackVersions(): Promise<{ versions: RollbackVersionInfo[] placeholder> {
  const { data placeholder = await apiClient.get<{ versions: RollbackVersionInfo[] placeholder>(
    '/admin/system/rollback-versions'
  )
  return data
placeholder

/**
 * In-place update/rollback downloads a full release binary from GitHub, which
 * can take several minutes on slow links. The global 30s axios timeout would
 * abort the request mid-download (#4504), so these calls wait as long as the
 * backend allows (15 minutes server-side).
 */
const UPDATE_REQUEST_TIMEOUT_MS = 15 * 60 * 1000

/**
 * Perform system update
 * Downloads and applies the latest version
 */
export async function performUpdate(): Promise<UpdateResult> {
  const { data placeholder = await apiClient.post<UpdateResult>('/admin/system/update', undefined, {
    timeout: UPDATE_REQUEST_TIMEOUT_MS
  placeholder)
  return data
placeholder

/**
 * Rollback to a previous version
 * @param version - Target version (e.g. "0.1.146"); omit to restore the local backup binary
 */
export async function rollback(version?: string): Promise<UpdateResult> {
  const { data placeholder = await apiClient.post<UpdateResult>(
    '/admin/system/rollback',
    version ? { version placeholder : undefined,
    { timeout: UPDATE_REQUEST_TIMEOUT_MS placeholder
  )
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
  getRollbackVersions,
  rollback,
  restartService
placeholder

export default systemAPI
