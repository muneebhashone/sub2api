import { describe, expect, it placeholder from 'vitest'
import {
  apiIntervalsToForm,
  apiTimePricingToForm,
  createDefaultTimePricingForm,
  formIntervalsToAPI,
  formTimePricingToAPI,
  isValidPositiveMultiplier,
  validateIntervals,
  validateTimePricing,
  type IntervalFormEntry,
  type TimePricingFormEntry,
  type TimePricingPeriodFormEntry,
placeholder from '../types'

describe('interval multiplier conversion', () => {
  it('preserves component multipliers without MTok conversion', () => {
    const form = apiIntervalsToForm([{
      min_tokens: 272000,
      max_tokens: null,
      tier_label: '',
      input_price: null,
      output_price: null,
      cache_write_price: null,
      cache_read_price: null,
      input_multiplier: 2,
      output_multiplier: 1.5,
      cache_write_multiplier: 2,
      cache_read_multiplier: 2,
      per_request_price: null,
      sort_order: 0,
    placeholder])

    expect(form[0].input_multiplier).toBe(2)
    expect(form[0].output_multiplier).toBe(1.5)
    expect(formIntervalsToAPI(form)[0]).toMatchObject({
      input_multiplier: 2,
      output_multiplier: 1.5,
      cache_write_multiplier: 2,
      cache_read_multiplier: 2,
    placeholder)
  placeholder)
placeholder)

describe('positive multiplier validation', () => {
  it('accepts empty and positive values but rejects zero and negative values', () => {
    expect(isValidPositiveMultiplier(null)).toBe(true)
    expect(isValidPositiveMultiplier('')).toBe(true)
    expect(isValidPositiveMultiplier('0.5')).toBe(true)
    expect(isValidPositiveMultiplier(0)).toBe(false)
    expect(isValidPositiveMultiplier(-1)).toBe(false)
  placeholder)

  it('rejects a zero interval multiplier', () => {
    expect(validateIntervals([
      makeInterval({ min_tokens: 100, input_multiplier: 0 placeholder),
    ], 'token', t)).toContain('multiplierPositive')
  placeholder)
placeholder)

function makeInterval(over: Partial<IntervalFormEntry>): IntervalFormEntry {
  return {
    min_tokens: 0,
    max_tokens: null,
    tier_label: '',
    input_price: null,
    output_price: null,
    cache_write_price: null,
    cache_read_price: null,
    input_multiplier: null,
    output_multiplier: null,
    cache_write_multiplier: null,
    cache_read_multiplier: null,
    per_request_price: null,
    sort_order: 0,
    ...over,
  placeholder
placeholder

function t(key: string, params?: Record<string, unknown>): string {
  return `${keyplaceholder${params ? ` ${JSON.stringify(params)placeholder` : ''placeholder`
placeholder

describe('validateIntervals', () => {
  describe('token mode', () => {
    it('rejects unbounded interval that is not last', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: null, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 200000, max_tokens: 500000, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token', t)).toContain('unboundedLast')
    placeholder)

    it('accepts unbounded interval at the end', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: 200000, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 200000, max_tokens: null, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token', t)).toBeNull()
    placeholder)

    it('rejects overlapping intervals', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: 250000, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 200000, max_tokens: 500000, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token', t)).toContain('overlap')
    placeholder)

    it('rejects unbounded interval in token mode', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ min_tokens: 0, max_tokens: null, input_price: 1, output_price: 1 placeholder),
        makeInterval({ min_tokens: 100, max_tokens: 200, input_price: 2, output_price: 2 placeholder),
      ]
      expect(validateIntervals(intervals, 'token', t)).toContain('unboundedLast')
    placeholder)
  placeholder)

  describe('image / per_request mode', () => {
    it('allows multiple unbounded tiers identified by label', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', per_request_price: 0.04 placeholder),
        makeInterval({ tier_label: '2K', per_request_price: 0.06 placeholder),
        makeInterval({ tier_label: '4K', per_request_price: 0.08 placeholder),
      ]
      expect(validateIntervals(intervals, 'image', t)).toBeNull()
      expect(validateIntervals(intervals, 'per_request', t)).toBeNull()
    placeholder)

    it('still rejects negative prices', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', per_request_price: -1 placeholder),
      ]
      expect(validateIntervals(intervals, 'image', t)).toContain('negativePrice')
    placeholder)

    it('still rejects max <= min on a single tier', () => {
      const intervals: IntervalFormEntry[] = [
        makeInterval({ tier_label: '1K', min_tokens: 100, max_tokens: 50, per_request_price: 0.04 placeholder),
      ]
      expect(validateIntervals(intervals, 'image', t)).toContain('maxGreaterThanMin')
    placeholder)
  placeholder)
placeholder)

describe('time pricing', () => {
  it('uses a disabled Shanghai default', () => {
    const form = createDefaultTimePricingForm()
    expect(form).toEqual({ timezone: 'Asia/Shanghai', periods: [], weekdays_only: false placeholder)
    expect(formTimePricingToAPI(form)).toBeNull()
  placeholder)

  it('defaults missing API day scope to every day', () => {
    expect(apiTimePricingToForm({
      timezone: 'Asia/Shanghai',
      periods: [{ start_time: '09:00', end_time: '12:00', multiplier: 2 placeholder],
    placeholder).weekdays_only).toBe(false)
  placeholder)

  it('round-trips day scope and formats multiplier', () => {
    const form = apiTimePricingToForm({
      timezone: 'Asia/Shanghai',
      weekdays_only: true,
      periods: [{ start_time: '09:00', end_time: '12:00', multiplier: 2 placeholder],
    placeholder)
    expect(form.weekdays_only).toBe(true)
    expect(form.periods[0]).toEqual({
      start_time: '09:00:00',
      end_time: '12:00:00',
      multiplier: '2.00',
    placeholder)
    expect(formTimePricingToAPI(form)).toEqual({
      timezone: 'Asia/Shanghai',
      weekdays_only: true,
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: 2 placeholder],
    placeholder)
  placeholder)

  it.each([
    ['separated', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder, { start_time: '14:00:00', end_time: '18:00:00', multiplier: '2.00' placeholder], null],
    ['adjacent', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder, { start_time: '12:00:00', end_time: '14:00:00', multiplier: '1.50' placeholder], null],
    ['midnight split', [{ start_time: '22:00:00', end_time: '00:00:00', multiplier: '2.00' placeholder, { start_time: '00:00:00', end_time: '02:00:00', multiplier: '2.00' placeholder], null],
    ['overlap by one second', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder, { start_time: '11:59:59', end_time: '14:00:00', multiplier: '2.00' placeholder], 'overlap'],
    ['cross midnight', [{ start_time: '22:00:00', end_time: '02:00:00', multiplier: '2.00' placeholder], 'range'],
    ['equal midnight', [{ start_time: '00:00:00', end_time: '00:00:00', multiplier: '2.00' placeholder], 'range'],
    ['missing seconds', [{ start_time: '09:00', end_time: '12:00', multiplier: '2.00' placeholder], 'format'],
    ['zero', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '0.00' placeholder], 'multiplier'],
    ['three decimals', [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '1.001' placeholder], 'multiplier'],
  ])('%s', (_name, periods, errorKey) => {
    const result = validateTimePricing({
      timezone: 'Asia/Shanghai',
      periods: periods as TimePricingPeriodFormEntry[],
    placeholder, t)
    if (errorKey === null) expect(result).toBeNull()
    else expect(result).toContain(String(errorKey))
  placeholder)

  it('rejects non-IANA timezone', () => {
    expect(validateTimePricing({
      timezone: 'UTC+8',
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder, t)).toContain('timezone')
  placeholder)

  it.each([
    ['missing', undefined],
    ['blank', '   '],
  ])('rejects a %s timezone without throwing during conversion', (_name, timezone) => {
    const form = {
      timezone,
      periods: [{ start_time: '09:00:00', end_time: '12:00:00', multiplier: '2.00' placeholder],
    placeholder as unknown as TimePricingFormEntry

    expect(validateTimePricing(form, t)).toContain('timezone')
    expect(() => formTimePricingToAPI(form)).not.toThrow()
    expect(formTimePricingToAPI(form)?.timezone).toBe('')
  placeholder)
placeholder)
