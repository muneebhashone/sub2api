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
  const { data placeholder = await apiClient.post<BackupRecord>('/admin/backups', req || {placeholder, { timeout: 600000 placeholder)
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
export async function restoreBackup(id: string, password: string): Promise<void> {
  await apiClient.post(`/admin/backups/${idplaceholder/restore`, { password placeholder, { timeout: 600000 placeholder)
placeholder

export const backupAPI = {
  getS3Config,
  updateS3Config,
  testS3Connection,
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
