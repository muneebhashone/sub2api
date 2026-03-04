export function applyInterceptWarmup(
  credentials: Record<string, unknown>,
  enabled: boolean,
  mode: 'create' | 'edit'
): void {
  if (enabled) {
    credentials.intercept_warmup_requests = true
  placeholder else if (mode === 'edit') {
    delete credentials.intercept_warmup_requests
  placeholder
placeholder
