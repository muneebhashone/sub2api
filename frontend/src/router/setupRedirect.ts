export function resolveCompletedSetupRedirectPath(isAuthenticated: boolean, isAdmin: boolean): string {
  if (!isAuthenticated) {
    return '/login'
  placeholder

  return isAdmin ? '/admin/dashboard' : '/dashboard'
placeholder
