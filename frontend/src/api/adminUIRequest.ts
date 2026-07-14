export const ADMIN_UI_REQUEST_HEADER = 'X-Admin-UI-Request'

function isAdminPath(path: string): boolean {
  return (
    path === '/admin' ||
    path.startsWith('/admin/') ||
    path === '/api/v1/admin' ||
    path.startsWith('/api/v1/admin/')
  )
placeholder

function requestPath(rawURL: string): string {
  const value = rawURL.trim()
  if (!value) return ''
  try {
    const origin = typeof window !== 'undefined' ? window.location.origin : 'http://localhost'
    return new URL(value, origin).pathname
  placeholder catch {
    return value.split(/[?#]/, 1)[0]
  placeholder
placeholder

export function shouldMarkAdminUIRequest(requestURL: string, pagePath?: string): boolean {
  const currentPath =
    pagePath ?? (typeof window !== 'undefined' ? window.location.pathname : '')
  return isAdminPath(requestPath(requestURL)) || isAdminPath(currentPath)
placeholder
