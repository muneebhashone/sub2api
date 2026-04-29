export const OPENAI_WS_MODE_OFF = 'off'
export const OPENAI_WS_MODE_CTX_POOL = 'ctx_pool'
export const OPENAI_WS_MODE_PASSTHROUGH = 'passthrough'
export const OPENAI_WS_MODE_HTTP_BRIDGE = 'http_bridge'

export type OpenAIWSMode =
  | typeof OPENAI_WS_MODE_OFF
  | typeof OPENAI_WS_MODE_CTX_POOL
  | typeof OPENAI_WS_MODE_PASSTHROUGH
  | typeof OPENAI_WS_MODE_HTTP_BRIDGE

const OPENAI_WS_MODES = new Set<OpenAIWSMode>([
  OPENAI_WS_MODE_OFF,
  OPENAI_WS_MODE_CTX_POOL,
  OPENAI_WS_MODE_PASSTHROUGH,
  OPENAI_WS_MODE_HTTP_BRIDGE
])

export interface ResolveOpenAIWSModeOptions {
  modeKey: string
  enabledKey: string
  fallbackEnabledKeys?: string[]
  defaultMode?: OpenAIWSMode
placeholder

export const normalizeOpenAIWSMode = (mode: unknown): OpenAIWSMode | null => {
  if (typeof mode !== 'string') return null
  const normalized = mode.trim().toLowerCase()
  if (normalized === 'shared' || normalized === 'dedicated') {
    return OPENAI_WS_MODE_CTX_POOL
  placeholder
  if (OPENAI_WS_MODES.has(normalized as OpenAIWSMode)) {
    return normalized as OpenAIWSMode
  placeholder
  return null
placeholder

export const openAIWSModeFromEnabled = (enabled: unknown): OpenAIWSMode | null => {
  if (typeof enabled !== 'boolean') return null
  return enabled ? OPENAI_WS_MODE_CTX_POOL : OPENAI_WS_MODE_OFF
placeholder

export const isOpenAIWSModeEnabled = (mode: OpenAIWSMode): boolean => {
  return mode !== OPENAI_WS_MODE_OFF
placeholder

export const resolveOpenAIWSModeConcurrencyHintKey = (
  mode: OpenAIWSMode
): 'admin.accounts.openai.wsModeConcurrencyHint' | 'admin.accounts.openai.wsModePassthroughHint' => {
  if (mode === OPENAI_WS_MODE_PASSTHROUGH || mode === OPENAI_WS_MODE_HTTP_BRIDGE) {
    return 'admin.accounts.openai.wsModePassthroughHint'
  placeholder
  return 'admin.accounts.openai.wsModeConcurrencyHint'
placeholder

export const resolveOpenAIWSModeFromExtra = (
  extra: Record<string, unknown> | null | undefined,
  options: ResolveOpenAIWSModeOptions
): OpenAIWSMode => {
  const fallback = options.defaultMode ?? OPENAI_WS_MODE_OFF
  if (!extra) return fallback

  const mode = normalizeOpenAIWSMode(extra[options.modeKey])
  if (mode) return mode

  const enabledMode = openAIWSModeFromEnabled(extra[options.enabledKey])
  if (enabledMode) return enabledMode

  const fallbackKeys = options.fallbackEnabledKeys ?? []
  for (const key of fallbackKeys) {
    const modeFromFallbackKey = openAIWSModeFromEnabled(extra[key])
    if (modeFromFallbackKey) return modeFromFallbackKey
  placeholder

  return fallback
placeholder
