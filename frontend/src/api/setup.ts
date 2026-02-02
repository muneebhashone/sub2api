/**
 * Setup API endpoints
 */
import axios from 'axios'

// Create a separate client for setup endpoints (not under /api/v1)
const setupClient = axios.create({
  baseURL: '',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  placeholder
placeholder)

export interface SetupStatus {
  needs_setup: boolean
  step: string
placeholder

export interface DatabaseConfig {
  host: string
  port: number
  user: string
  password: string
  dbname: string
  sslmode: string
placeholder

export interface RedisConfig {
  host: string
  port: number
  password: string
  db: number
  enable_tls: boolean
placeholder

export interface AdminConfig {
  email: string
  password: string
placeholder

export interface ServerConfig {
  host: string
  port: number
  mode: string
placeholder

export interface InstallRequest {
  database: DatabaseConfig
  redis: RedisConfig
  admin: AdminConfig
  server: ServerConfig
placeholder

export interface InstallResponse {
  message: string
  restart: boolean
placeholder

/**
 * Get setup status
 */
export async function getSetupStatus(): Promise<SetupStatus> {
  const response = await setupClient.get('/setup/status')
  return response.data.data
placeholder

/**
 * Test database connection
 */
export async function testDatabase(config: DatabaseConfig): Promise<void> {
  await setupClient.post('/setup/test-db', config)
placeholder

/**
 * Test Redis connection
 */
export async function testRedis(config: RedisConfig): Promise<void> {
  await setupClient.post('/setup/test-redis', config)
placeholder

/**
 * Perform installation
 */
export async function install(config: InstallRequest): Promise<InstallResponse> {
  const response = await setupClient.post('/setup/install', config)
  return response.data.data
placeholder
