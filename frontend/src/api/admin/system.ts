/**
 * System API endpoints for admin operations
 */

import { apiClient placeholder from '../client';

export interface ReleaseInfo {
  name: string;
  body: string;
  published_at: string;
  html_url: string;
placeholder

export interface VersionInfo {
  current_version: string;
  latest_version: string;
  has_update: boolean;
  release_info?: ReleaseInfo;
  cached: boolean;
  warning?: string;
  build_type: string; // "source" for manual builds, "release" for CI builds
placeholder

/**
 * Get current version
 */
export async function getVersion(): Promise<{ version: string placeholder> {
  const { data placeholder = await apiClient.get<{ version: string placeholder>('/admin/system/version');
  return data;
placeholder

/**
 * Check for updates
 * @param force - Force refresh from GitHub API
 */
export async function checkUpdates(force = false): Promise<VersionInfo> {
  const { data placeholder = await apiClient.get<VersionInfo>('/admin/system/check-updates', {
    params: force ? { force: 'true' placeholder : undefined,
  placeholder);
  return data;
placeholder

export const systemAPI = {
  getVersion,
  checkUpdates,
placeholder;

export default systemAPI;
