import { apiClient placeholder from '../client'

export type BackupType = 'postgres' | 'redis' | 'full'
export type BackupJobStatus = 'queued' | 'running' | 'succeeded' | 'failed' | 'partial_succeeded'

export interface BackupAgentInfo {
  status: string
  version: string
  uptime_seconds: number
placeholder

export interface BackupAgentHealth {
  enabled: boolean
  reason: string
  socket_path: string
  agent?: BackupAgentInfo
placeholder

export interface DataManagementPostgresConfig {
  host: string
  port: number
  user: string
  password?: string
  password_configured?: boolean
  database: string
  ssl_mode: string
  container_name: string
placeholder

export interface DataManagementRedisConfig {
  addr: string
  username: string
  password?: string
  password_configured?: boolean
  db: number
  container_name: string
placeholder

export interface DataManagementS3Config {
  enabled: boolean
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  secret_access_key?: string
  secret_access_key_configured?: boolean
  prefix: string
  force_path_style: boolean
  use_ssl: boolean
placeholder

export interface DataManagementConfig {
  source_mode: 'direct' | 'docker_exec'
  backup_root: string
  sqlite_path?: string
  retention_days: number
  keep_last: number
  active_postgres_profile_id?: string
  active_redis_profile_id?: string
  active_s3_profile_id?: string
  postgres: DataManagementPostgresConfig
  redis: DataManagementRedisConfig
  s3: DataManagementS3Config
placeholder

export type SourceType = 'postgres' | 'redis'

export interface DataManagementSourceConfig {
  host: string
  port: number
  user: string
  password?: string
  database: string
  ssl_mode: string
  addr: string
  username: string
  db: number
  container_name: string
placeholder

export interface DataManagementSourceProfile {
  source_type: SourceType
  profile_id: string
  name: string
  is_active: boolean
  password_configured?: boolean
  config: DataManagementSourceConfig
  created_at?: string
  updated_at?: string
placeholder

export interface TestS3Request {
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  secret_access_key: string
  prefix?: string
  force_path_style?: boolean
  use_ssl?: boolean
placeholder

export interface TestS3Response {
  ok: boolean
  message: string
placeholder

export interface CreateBackupJobRequest {
  backup_type: BackupType
  upload_to_s3?: boolean
  s3_profile_id?: string
  postgres_profile_id?: string
  redis_profile_id?: string
  idempotency_key?: string
placeholder

export interface CreateBackupJobResponse {
  job_id: string
  status: BackupJobStatus
placeholder

export interface BackupArtifactInfo {
  local_path: string
  size_bytes: number
  sha256: string
placeholder

export interface BackupS3Info {
  bucket: string
  key: string
  etag: string
placeholder

export interface BackupJob {
  job_id: string
  backup_type: BackupType
  status: BackupJobStatus
  triggered_by: string
  s3_profile_id?: string
  postgres_profile_id?: string
  redis_profile_id?: string
  started_at?: string
  finished_at?: string
  error_message?: string
  artifact?: BackupArtifactInfo
  s3?: BackupS3Info
placeholder

export interface ListSourceProfilesResponse {
  items: DataManagementSourceProfile[]
placeholder

export interface CreateSourceProfileRequest {
  profile_id: string
  name: string
  config: DataManagementSourceConfig
  set_active?: boolean
placeholder

export interface UpdateSourceProfileRequest {
  name: string
  config: DataManagementSourceConfig
placeholder

export interface DataManagementS3Profile {
  profile_id: string
  name: string
  is_active: boolean
  s3: DataManagementS3Config
  secret_access_key_configured?: boolean
  created_at?: string
  updated_at?: string
placeholder

export interface ListS3ProfilesResponse {
  items: DataManagementS3Profile[]
placeholder

export interface CreateS3ProfileRequest {
  profile_id: string
  name: string
  enabled: boolean
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  secret_access_key?: string
  prefix?: string
  force_path_style?: boolean
  use_ssl?: boolean
  set_active?: boolean
placeholder

export interface UpdateS3ProfileRequest {
  name: string
  enabled: boolean
  endpoint: string
  region: string
  bucket: string
  access_key_id: string
  secret_access_key?: string
  prefix?: string
  force_path_style?: boolean
  use_ssl?: boolean
placeholder

export interface ListBackupJobsRequest {
  page_size?: number
  page_token?: string
  status?: BackupJobStatus
  backup_type?: BackupType
placeholder

export interface ListBackupJobsResponse {
  items: BackupJob[]
  next_page_token?: string
placeholder

export async function getAgentHealth(): Promise<BackupAgentHealth> {
  const { data placeholder = await apiClient.get<BackupAgentHealth>('/admin/data-management/agent/health')
  return data
placeholder

export async function getConfig(): Promise<DataManagementConfig> {
  const { data placeholder = await apiClient.get<DataManagementConfig>('/admin/data-management/config')
  return data
placeholder

export async function updateConfig(request: DataManagementConfig): Promise<DataManagementConfig> {
  const { data placeholder = await apiClient.put<DataManagementConfig>('/admin/data-management/config', request)
  return data
placeholder

export async function testS3(request: TestS3Request): Promise<TestS3Response> {
  const { data placeholder = await apiClient.post<TestS3Response>('/admin/data-management/s3/test', request)
  return data
placeholder

export async function listSourceProfiles(sourceType: SourceType): Promise<ListSourceProfilesResponse> {
  const { data placeholder = await apiClient.get<ListSourceProfilesResponse>(`/admin/data-management/sources/${sourceTypeplaceholder/profiles`)
  return data
placeholder

export async function createSourceProfile(sourceType: SourceType, request: CreateSourceProfileRequest): Promise<DataManagementSourceProfile> {
  const { data placeholder = await apiClient.post<DataManagementSourceProfile>(`/admin/data-management/sources/${sourceTypeplaceholder/profiles`, request)
  return data
placeholder

export async function updateSourceProfile(sourceType: SourceType, profileID: string, request: UpdateSourceProfileRequest): Promise<DataManagementSourceProfile> {
  const { data placeholder = await apiClient.put<DataManagementSourceProfile>(`/admin/data-management/sources/${sourceTypeplaceholder/profiles/${profileIDplaceholder`, request)
  return data
placeholder

export async function deleteSourceProfile(sourceType: SourceType, profileID: string): Promise<void> {
  await apiClient.delete(`/admin/data-management/sources/${sourceTypeplaceholder/profiles/${profileIDplaceholder`)
placeholder

export async function setActiveSourceProfile(sourceType: SourceType, profileID: string): Promise<DataManagementSourceProfile> {
  const { data placeholder = await apiClient.post<DataManagementSourceProfile>(`/admin/data-management/sources/${sourceTypeplaceholder/profiles/${profileIDplaceholder/activate`)
  return data
placeholder

export async function listS3Profiles(): Promise<ListS3ProfilesResponse> {
  const { data placeholder = await apiClient.get<ListS3ProfilesResponse>('/admin/data-management/s3/profiles')
  return data
placeholder

export async function createS3Profile(request: CreateS3ProfileRequest): Promise<DataManagementS3Profile> {
  const { data placeholder = await apiClient.post<DataManagementS3Profile>('/admin/data-management/s3/profiles', request)
  return data
placeholder

export async function updateS3Profile(profileID: string, request: UpdateS3ProfileRequest): Promise<DataManagementS3Profile> {
  const { data placeholder = await apiClient.put<DataManagementS3Profile>(`/admin/data-management/s3/profiles/${profileIDplaceholder`, request)
  return data
placeholder

export async function deleteS3Profile(profileID: string): Promise<void> {
  await apiClient.delete(`/admin/data-management/s3/profiles/${profileIDplaceholder`)
placeholder

export async function setActiveS3Profile(profileID: string): Promise<DataManagementS3Profile> {
  const { data placeholder = await apiClient.post<DataManagementS3Profile>(`/admin/data-management/s3/profiles/${profileIDplaceholder/activate`)
  return data
placeholder

export async function createBackupJob(request: CreateBackupJobRequest): Promise<CreateBackupJobResponse> {
  const headers = request.idempotency_key
    ? { 'X-Idempotency-Key': request.idempotency_key placeholder
    : undefined

  const { data placeholder = await apiClient.post<CreateBackupJobResponse>(
    '/admin/data-management/backups',
    request,
    { headers placeholder
  )
  return data
placeholder

export async function listBackupJobs(request?: ListBackupJobsRequest): Promise<ListBackupJobsResponse> {
  const { data placeholder = await apiClient.get<ListBackupJobsResponse>('/admin/data-management/backups', {
    params: request
  placeholder)
  return data
placeholder

export async function getBackupJob(jobID: string): Promise<BackupJob> {
  const { data placeholder = await apiClient.get<BackupJob>(`/admin/data-management/backups/${jobIDplaceholder`)
  return data
placeholder

export const dataManagementAPI = {
  getAgentHealth,
  getConfig,
  updateConfig,
  listSourceProfiles,
  createSourceProfile,
  updateSourceProfile,
  deleteSourceProfile,
  setActiveSourceProfile,
  testS3,
  listS3Profiles,
  createS3Profile,
  updateS3Profile,
  deleteS3Profile,
  setActiveS3Profile,
  createBackupJob,
  listBackupJobs,
  getBackupJob
placeholder

export default dataManagementAPI
