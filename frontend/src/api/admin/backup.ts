import { apiClient placeholder from '../client'

export interface BackupS3Config {
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  secret_access_key?: string
  prefix: string
  force_path_style: boolean
placeholder

export interface BackupScheduleConfig {
  enabled: boolean
  cron_expr: string
  retain_days: number
  retain_count: number
placeholder

export interface BackupRecord {
  id: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  backup_type: string
  file_name: string
  s3_key: string
  size_bytes: number
  triggered_by: string
  error_message?: string
  started_at: string
  finished_at?: string
  expires_at?: string
  progress?: string
  restore_status?: string
  restore_error?: string
  restored_at?: string
placeholder

export interface CreateBackupRequest {
  expire_days?: number
placeholder

export interface TestS3Response {
  ok: boolean
  message: string
placeholder

// S3 Config
export async function getS3Config(): Promise<BackupS3Config> {
  const { data placeholder = await apiClient.get<BackupS3Config>('/admin/backups/s3-config')
  return data
placeholder

export async function updateS3Config(config: BackupS3Config): Promise<BackupS3Config> {
  const { data placeholder = await apiClient.put<BackupS3Config>('/admin/backups/s3-config', config)
  return data
placeholder

export async function testS3Connection(config: BackupS3Config): Promise<TestS3Response> {
  const { data placeholder = await apiClient.post<TestS3Response>('/admin/backups/s3-config/test', config)
  return data
placeholder

// Async image object storage
//
// Shares the S3 client with backups, so `reuse_backup_s3` borrows the endpoint and
// credentials configured above and only keeps its own bucket/prefix.
export interface ImageStorageConfig {
  enabled: boolean
  reuse_backup_s3: boolean
  bucket: string
  prefix: string
  public_base_url: string
  presign_expiry_hours: number
  max_download_bytes: number
  endpoint: string
  region: string
  access_key_id: string
  secret_access_key?: string
  force_path_style: boolean
placeholder

export interface ImageStorageConfigResponse {
  config: ImageStorageConfig
  secret_configured: boolean
placeholder

export async function getImageStorageConfig(): Promise<ImageStorageConfigResponse> {
  const { data placeholder = await apiClient.get<ImageStorageConfigResponse>('/admin/backups/image-storage')
  return data
placeholder

export async function updateImageStorageConfig(
  config: ImageStorageConfig,
): Promise<ImageStorageConfig> {
  const { data placeholder = await apiClient.put<ImageStorageConfig>('/admin/backups/image-storage', config)
  return data
placeholder

export async function testImageStorageConnection(
  config: ImageStorageConfig,
): Promise<TestS3Response> {
  const { data placeholder = await apiClient.post<TestS3Response>(
    '/admin/backups/image-storage/test',
    config,
  )
  return data
placeholder

// Schedule
export async function getSchedule(): Promise<BackupScheduleConfig> {
  const { data placeholder = await apiClient.get<BackupScheduleConfig>('/admin/backups/schedule')
  return data
placeholder

export async function updateSchedule(config: BackupScheduleConfig): Promise<BackupScheduleConfig> {
  const { data placeholder = await apiClient.put<BackupScheduleConfig>('/admin/backups/schedule', config)
  return data
placeholder

// Backup operations
export async function createBackup(req?: CreateBackupRequest): Promise<BackupRecord> {
  const { data placeholder = await apiClient.post<BackupRecord>('/admin/backups', req || {placeholder)
  return data
placeholder

export async function listBackups(): Promise<{ items: BackupRecord[] placeholder> {
  const { data placeholder = await apiClient.get<{ items: BackupRecord[] placeholder>('/admin/backups')
  return data
placeholder

export async function getBackup(id: string): Promise<BackupRecord> {
  const { data placeholder = await apiClient.get<BackupRecord>(`/admin/backups/${idplaceholder`)
  return data
placeholder

export async function deleteBackup(id: string): Promise<void> {
  await apiClient.delete(`/admin/backups/${idplaceholder`)
placeholder

export async function getDownloadURL(id: string): Promise<{ url: string placeholder> {
  const { data placeholder = await apiClient.get<{ url: string placeholder>(`/admin/backups/${idplaceholder/download-url`)
  return data
placeholder

// Restore
export async function restoreBackup(id: string, password: string): Promise<BackupRecord> {
  const { data placeholder = await apiClient.post<BackupRecord>(`/admin/backups/${idplaceholder/restore`, { password placeholder)
  return data
placeholder

export const backupAPI = {
  getS3Config,
  updateS3Config,
  testS3Connection,
  getImageStorageConfig,
  updateImageStorageConfig,
  testImageStorageConnection,
  getSchedule,
  updateSchedule,
  createBackup,
  listBackups,
  getBackup,
  deleteBackup,
  getDownloadURL,
  restoreBackup,
placeholder

export default backupAPI
