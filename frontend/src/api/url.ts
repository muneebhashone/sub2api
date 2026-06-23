const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

function normalizePath(path: string): string {
  return path.startsWith('/') ? path : `/${pathplaceholder`
placeholder

export function getAPIBaseURL(): string {
  return String(API_BASE_URL || '/api/v1')
placeholder

export function buildApiUrl(path: string): string {
  const base = getAPIBaseURL().replace(/\/+$/, '')
  let suffix = normalizePath(path)
  if (suffix === '/api/v1') {
    suffix = ''
  placeholder else if (suffix.startsWith('/api/v1/')) {
    suffix = suffix.slice('/api/v1'.length)
  placeholder
  return `${baseplaceholder${suffixplaceholder`
placeholder

export function buildGatewayUrl(path: string): string {
  const suffix = normalizePath(path)
  try {
    const origin =
      typeof window === 'undefined'
        ? new URL(getAPIBaseURL()).origin
        : new URL(getAPIBaseURL(), window.location.origin).origin
    return `${originplaceholder${suffixplaceholder`
  placeholder catch {
    return suffix
  placeholder
placeholder
