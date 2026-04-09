import { getConfiguredTableDefaultPageSize, normalizeTablePageSize placeholder from '@/utils/tablePreferences'

const STORAGE_KEY = 'table-page-size'
const SOURCE_KEY = 'table-page-size-source'

/**
 * 从 localStorage 读取/写入 pageSize
 * 全局共享一个 key，所有表格统一偏好
 */
export function getPersistedPageSize(fallback = getConfiguredTableDefaultPageSize()): number {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      return normalizeTablePageSize(stored)
    placeholder
  placeholder catch {
    // localStorage 不可用（隐私模式等）
  placeholder
  return normalizeTablePageSize(fallback)
placeholder

export function setPersistedPageSize(size: number): void {
  try {
    localStorage.setItem(STORAGE_KEY, String(normalizeTablePageSize(size)))
    localStorage.setItem(SOURCE_KEY, 'user')
  placeholder catch {
    // 静默失败
  placeholder
placeholder

export function syncPersistedPageSizeWithSystemDefault(defaultSize = getConfiguredTableDefaultPageSize()): void {
  try {
    const normalizedDefault = normalizeTablePageSize(defaultSize)
    const stored = localStorage.getItem(STORAGE_KEY)
    const source = localStorage.getItem(SOURCE_KEY)
    const normalizedStored = stored ? normalizeTablePageSize(stored) : null

    if ((source === 'user' || (source === null && stored !== null)) && stored) {
      localStorage.setItem(STORAGE_KEY, String(normalizedStored ?? normalizedDefault))
      localStorage.setItem(SOURCE_KEY, 'user')
      return
    placeholder

    localStorage.setItem(STORAGE_KEY, String(normalizedDefault))
    localStorage.setItem(SOURCE_KEY, 'system')
  placeholder catch {
    // 静默失败
  placeholder
placeholder
