import { describe, expect, it placeholder from 'vitest'
import { detectMobileDevice placeholder from '../device'

describe('detectMobileDevice', () => {
  it('prefers userAgentData.mobile when available', () => {
    expect(detectMobileDevice({
      navigator: {
        userAgent: 'Mozilla/5.0',
        userAgentData: { mobile: true placeholder,
      placeholder,
    placeholder)).toBe(true)
  placeholder)

  it('recognizes handheld browsers from the mobile UA token', () => {
    expect(detectMobileDevice({
      navigator: {
        userAgent: 'Mozilla/5.0 (Linux; Android 14; Pixel 8) AppleWebKit/537.36 Chrome/136.0 Mobile Safari/537.36',
        maxTouchPoints: 5,
      placeholder,
    placeholder)).toBe(true)
  placeholder)

  it('recognizes iPadOS desktop mode via touch capability', () => {
    expect(detectMobileDevice({
      navigator: {
        userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15) AppleWebKit/605.1.15 Version/17.0 Safari/605.1.15',
        platform: 'MacIntel',
        maxTouchPoints: 5,
      placeholder,
    placeholder)).toBe(true)
  placeholder)

  it('falls back to input capability detection for touch-first devices', () => {
    expect(detectMobileDevice({
      navigator: {
        userAgent: 'Mozilla/5.0',
        maxTouchPoints: 10,
      placeholder,
      matchMedia: (query) => ({
        matches: query === '(pointer: coarse)' || query === '(hover: none)',
      placeholder),
    placeholder)).toBe(true)
  placeholder)

  it('keeps desktop environments as non-mobile', () => {
    expect(detectMobileDevice({
      navigator: {
        userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Chrome/136.0 Safari/537.36',
        platform: 'MacIntel',
        maxTouchPoints: 0,
      placeholder,
      matchMedia: () => ({ matches: false placeholder),
    placeholder)).toBe(false)
  placeholder)
placeholder)
