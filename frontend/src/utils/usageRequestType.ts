import type { UsageRequestType placeholder from '@/types'

export interface UsageRequestTypeLike {
  request_type?: string | null
  stream?: boolean | null
  openai_ws_mode?: boolean | null
placeholder

const VALID_REQUEST_TYPES = new Set<UsageRequestType>(['unknown', 'sync', 'stream', 'ws_v2'])

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
  if (!requestType || requestType === 'unknown') {
    return null
  placeholder
  if (requestType === 'sync') {
    return false
  placeholder
  return true
placeholder
