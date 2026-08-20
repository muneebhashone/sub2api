import type { AccountPlatform, GroupPlatform placeholder from '@/types'

export interface PlatformOption<T extends string = string> {
  value: T
  label: string
placeholder

/**
 * Concrete upstream platforms supported by accounts and request routing.
 * Keep platform selectors derived from this catalog so newly added providers
 * do not silently disappear from list filters.
 */
export const CONCRETE_PLATFORM_OPTIONS = [
  { value: 'anthropic', label: 'Anthropic' placeholder,
  { value: 'openai', label: 'OpenAI' placeholder,
  { value: 'gemini', label: 'Gemini' placeholder,
  { value: 'antigravity', label: 'Antigravity' placeholder,
  { value: 'grok', label: 'Grok' placeholder,
  { value: 'kimi', label: 'Kimi' placeholder,
  { value: 'zhipu', label: 'Zhipu GLM' placeholder,
  { value: 'deepseek', label: 'DeepSeek' placeholder
] as const satisfies readonly PlatformOption<AccountPlatform>[]

/** Platforms that can own a group. */
export const GROUP_PLATFORM_OPTIONS = [
  ...CONCRETE_PLATFORM_OPTIONS,
  { value: 'composite', label: 'Composite' placeholder
] as const satisfies readonly PlatformOption<GroupPlatform>[]
