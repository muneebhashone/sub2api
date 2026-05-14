import { describe, expect, it placeholder from 'vitest'
import { validateIntervals, type IntervalFormEntry placeholder from '../types'

function makeInterval(over: Partial<IntervalFormEntry>): IntervalFormEntry {
  return {
    min_tokens: 0,
    max_tokens: null,
    tier_label: '',
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    per_request_price: null,
    sort_order: 0,
    ...over,
  placeholder
placeholder

describe('validateIntervals', () => {
  describe('token mode', () => {
    it('rejects unbounded interval that is not last', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: null, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 200000, max_tokens: 500000, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token')).toMatch(/无上限/)
    placeholder)

    it('accepts unbounded interval at the end', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: 200000, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 200000, max_tokens: null, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token')).toBeNull()
    placeholder)

    it('rejects overlapping intervals', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: 250000, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 200000, max_tokens: 500000, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token')).toMatch(/重叠/)
    placeholder)

    it('defaults mode to token when omitted', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: null, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 100, max_tokens: 200, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals)).toMatch(/无上限/)
    placeholder)
  placeholder)

  describe('image / per_request mode', () => {
    it('allows multiple unbounded tiers identified by label', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', per_request_price: 0.04 placeholder),
        makeInterval({ tier_label: '2K', per_request_price: 0.06 placeholder),
        makeInterval({ tier_label: '4K', per_request_price: 0.08 placeholder),
      ]
      expect(validateIntervals(intervals, 'image')).toBeNull()
      expect(validateIntervals(intervals, 'per_request')).toBeNull()
    placeholder)

    it('still rejects negative prices', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', per_request_price: -1 placeholder),
      ]
      expect(validateIntervals(intervals, 'image')).toMatch(/不能为负数/)
    placeholder)

    it('still rejects max <= min on a single tier', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', min_tokens: 100, max_tokens: 50, per_request_price: 0.04 placeholder),
      ]
      expect(validateIntervals(intervals, 'image')).toMatch(/必须大于/)
    placeholder)
  placeholder)
placeholder)
