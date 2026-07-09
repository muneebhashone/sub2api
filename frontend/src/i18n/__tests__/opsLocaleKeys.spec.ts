import { describe, expect, it placeholder from 'vitest'
import en from '@/i18n/locales/en'

function flattenKeys(obj: Record<string, any>, prefix = ''): string[] {
  const keys: string[] = []
  for (const [k, v] of Object.entries(obj)) {
    const fullKey = prefix ? `${prefixplaceholder.${kplaceholder` : k
    if (typeof v === 'object' && v !== null && !Array.isArray(v)) {
      keys.push(...flattenKeys(v, fullKey))
    placeholder else {
      keys.push(fullKey)
    placeholder
  placeholder
  return keys
placeholder

describe('ops locale key completeness', () => {
  const requiredKeys = [
    'admin.ops.result',
    'admin.ops.timeRange.custom',
    'admin.ops.customTimeRange.startTime',
    'admin.ops.customTimeRange.endTime',
  ]

  for (const key of requiredKeys) {
    it(`en locale has ${keyplaceholder`, () => {
      const enKeys = flattenKeys(en)
      expect(enKeys).toContain(key)
    placeholder)
  placeholder
placeholder)

describe('groups locale key completeness', () => {
  it('en locale has admin.groups.failedToSave', () => {
    const enKeys = flattenKeys(en)
    expect(enKeys).toContain('admin.groups.failedToSave')
  placeholder)
placeholder)
