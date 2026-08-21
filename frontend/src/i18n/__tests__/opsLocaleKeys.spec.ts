import { describe, expect, it placeholder from 'vitest'
import en from '@/i18n/locales/en'
import zh from '@/i18n/locales/zh'

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
    'admin.ops.errorDetail.upstreamStatus',
    'admin.ops.errorDetail.rootCause',
    'admin.ops.errorDetail.diagnosticPayloads',
    'admin.ops.errorDetail.payloads.client',
    'admin.ops.errorDetail.payloads.upstream_message',
    'admin.ops.errorDetail.payloads.upstream_detail',
    'admin.ops.errorDetail.payloads.upstream_events',
  ]

  for (const key of requiredKeys) {
    it(`en locale has ${keyplaceholder`, () => {
      const enKeys = flattenKeys(en)
      expect(enKeys).toContain(key)
    placeholder)
  placeholder

  for (const key of requiredKeys) {
    it(`zh locale has ${keyplaceholder`, () => {
      const zhKeys = flattenKeys(zh)
      expect(zhKeys).toContain(key)
    placeholder)
  placeholder
placeholder)

describe('groups locale key completeness', () => {
  it('en locale has admin.groups.failedToSave', () => {
    const enKeys = flattenKeys(en)
    expect(enKeys).toContain('admin.groups.failedToSave')
  placeholder)

  const webSearchPricingKeys = [
    'admin.groups.webSearchPricing.title',
    'admin.groups.webSearchPricing.pricePerCall',
    'admin.groups.webSearchPricing.pricePerCallHint',
    'admin.groups.webSearchPricing.finalPricePreview',
  ]

  for (const key of webSearchPricingKeys) {
    it(`en and zh locales both have ${keyplaceholder`, () => {
      expect(flattenKeys(en)).toContain(key)
      expect(flattenKeys(zh)).toContain(key)
    placeholder)
  placeholder
placeholder)
