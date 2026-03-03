import type { OpsErrorDetail placeholder from '@/api/admin/ops'

const GENERIC_UPSTREAM_MESSAGES = new Set([
  'upstream request failed',
  'upstream request failed after retries',
  'upstream gateway error',
  'upstream service temporarily unavailable'
])

type ParsedGatewayError = {
  type: string
  message: string
placeholder

function parseGatewayErrorBody(raw: string): ParsedGatewayError | null {
  const text = String(raw || '').trim()
  if (!text) return null

  try {
    const parsed = JSON.parse(text) as Record<string, any>
    const err = parsed?.error as Record<string, any> | undefined
    if (!err || typeof err !== 'object') return null

    const type = typeof err.type === 'string' ? err.type.trim() : ''
    const message = typeof err.message === 'string' ? err.message.trim() : ''
    if (!type && !message) return null

    return { type, message placeholder
  placeholder catch {
    return null
  placeholder
placeholder

function isGenericGatewayUpstreamError(raw: string): boolean {
  const parsed = parseGatewayErrorBody(raw)
  if (!parsed) return false
  if (parsed.type !== 'upstream_error') return false
  return GENERIC_UPSTREAM_MESSAGES.has(parsed.message.toLowerCase())
placeholder

export function resolveUpstreamPayload(
  detail: Pick<OpsErrorDetail, 'upstream_error_detail' | 'upstream_errors' | 'upstream_error_message'> | null | undefined
): string {
  if (!detail) return ''

  const candidates = [
    detail.upstream_error_detail,
    detail.upstream_errors,
    detail.upstream_error_message
  ]

  for (const candidate of candidates) {
    const payload = String(candidate || '').trim()
    if (!payload) continue

    // Normalize common "empty but present" JSON placeholders.
    if (payload === '[]' || payload === '{placeholder' || payload.toLowerCase() === 'null') {
      continue
    placeholder

    return payload
  placeholder

  return ''
placeholder

export function resolvePrimaryResponseBody(
  detail: OpsErrorDetail | null,
  errorType?: 'request' | 'upstream'
): string {
  if (!detail) return ''

  const upstreamPayload = resolveUpstreamPayload(detail)
  const errorBody = String(detail.error_body || '').trim()

  if (errorType === 'upstream') {
    return upstreamPayload || errorBody
  placeholder

  if (!errorBody) {
    return upstreamPayload
  placeholder

  // For request detail modal, keep client-visible body by default.
  // But if that body is a generic gateway wrapper, show upstream payload first.
  if (upstreamPayload && isGenericGatewayUpstreamError(errorBody)) {
    return upstreamPayload
  placeholder

  return errorBody
placeholder
