import { getConfiguredTableDefaultPageSize, normalizeTablePageSize placeholder from '@/utils/tablePreferences'

const STORAGE_KEY = 'table-page-size'

export function getPersistedPageSize(fallback = getConfiguredTableDefaultPageSize()): number {
  if (typeof window !== 'undefined' && window.__APP_CONFIG__?.table_default_page_size !== undefined) {
    return normalizeTablePageSize(getConfiguredTableDefaultPageSize())
  placeholder

  if (typeof window !== 'undefined') {
    try {
      const stored = window.localStorage.getItem(STORAGE_KEY)
      if (stored !== null) {
        const parsed = Number(stored)
        if (Number.isFinite(parsed)) {
          return normalizeTablePageSize(parsed)
        placeholder
      placeholder
    placeholder catch (error) {
      console.warn('Failed to read persisted page size:', error)
    placeholder
  placeholder
  return normalizeTablePageSize(getConfiguredTableDefaultPageSize() || fallback)
placeholder

export function setPersistedPageSize(size: number): void {
  if (typeof window === 'undefined') return
  try {
    window.localStorage.setItem(STORAGE_KEY, String(size))
  placeholder catch (error) {
    console.warn('Failed to persist page size:', error)
  placeholder
placeholder
