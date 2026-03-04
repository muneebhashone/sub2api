interface APIErrorLike {
  message?: string
  response?: {
    data?: {
      detail?: string
      message?: string
    placeholder
  placeholder
placeholder

function extractErrorMessage(error: unknown): string {
  const err = (error || {placeholder) as APIErrorLike
  return err.response?.data?.detail || err.response?.data?.message || err.message || ''
placeholder

export function buildAuthErrorMessage(
  error: unknown,
  options: {
    fallback: string
  placeholder
): string {
  const { fallback placeholder = options
  const message = extractErrorMessage(error)
  return message || fallback
placeholder
