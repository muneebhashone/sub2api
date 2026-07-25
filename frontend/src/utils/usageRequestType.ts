import type { UsageRequestType placeholder from '@/types'

export interface UsageRequestTypeLike {
  request_type?: string | null
  stream?: boolean | null
  openai_ws_mode?: boolean | null
placeholder

const VALID_REQUEST_TYPES = new Set<UsageRequestType>(['unknown', 'sync', 'stream', 'ws_v2', 'cyber', 'live'])

export const isUsageRequestType = (value: unknown): value is UsageRequestType => {
  return typeof value === 'string' && VALID_REQUEST_TYPES.has(value as UsageRequestType)
placeholder

export const resolveUsageRequestType = (value: UsageRequestTypeLike): UsageRequestType => {
  if (isUsageRequestType(value.request_type)) {
    return value.request_type
  placeholder
  if (value.openai_ws_mode) {
    return 'ws_v2'
  placeholder
  return value.stream ? 'stream' : 'sync'
placeholder

export const requestTypeToLegacyStream = (requestType?: UsageRequestType | null): boolean | null | undefined => {
  // cyber 与 stream 正交（cyber 可发生在 stream 或非 stream 请求），不映射到 legacy stream 维度。
  if (!requestType || requestType === 'unknown' || requestType === 'cyber' || requestType === 'live') {
    return null
  placeholder
  if (requestType === 'sync') {
    return false
  placeholder
  return true
placeholder
