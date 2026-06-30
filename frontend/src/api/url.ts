const DEFAULT_API_BASE_URL = '/api/v1'
const API_BASE_URL = normalizeAPIBaseURL(import.meta.env.VITE_API_BASE_URL)

function normalizePath(path: string): string {
  return path.startsWith('/') ? path : `/${pathplaceholder`
placeholder

function normalizeAPIBaseURL(value: unknown): string {
  const raw = String(value || DEFAULT_API_BASE_URL).trim() || DEFAULT_API_BASE_URL
  const withoutTrailingSlash = raw.replace(/\/+$/, '')
  if (/^[a-z][a-z\d+.-]*:\/\//i.test(withoutTrailingSlash) || withoutTrailingSlash.startsWith('//')) {
    return withoutTrailingSlash
  placeholder
  return normalizePath(withoutTrailingSlash)
placeholder

export function getAPIBaseURL(): string {
  return API_BASE_URL
placeholder

export function buildApiUrl(path: string): string {
  const base = getAPIBaseURL().replace(/\/+$/, '')
  let suffix = normalizePath(path)
  if (suffix === DEFAULT_API_BASE_URL) {
    suffix = ''
  placeholder else if (suffix.startsWith(`${DEFAULT_API_BASE_URLplaceholder/`)) {
    suffix = suffix.slice(DEFAULT_API_BASE_URL.length)
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
