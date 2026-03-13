import { describe, expect, it placeholder from 'vitest'
import { buildAuthErrorMessage placeholder from '@/utils/authError'

describe('buildAuthErrorMessage', () => {
  it('prefers response detail message when available', () => {
    const message = buildAuthErrorMessage(
      {
        response: {
          data: {
            detail: 'detailed message',
            message: 'plain message'
          placeholder
        placeholder,
      placeholder,
      { fallback: 'fallback' placeholder
    )
    expect(message).toBe('detailed message')
  placeholder)

  it('falls back to response message when detail is unavailable', () => {
    const message = buildAuthErrorMessage(
      {
        response: {
          data: {
            message: 'plain message'
          placeholder
        placeholder,
      placeholder,
      { fallback: 'fallback' placeholder
    )
    expect(message).toBe('plain message')
  placeholder)

  it('falls back to error.message when response payload is unavailable', () => {
    const message = buildAuthErrorMessage(
      {
        message: 'error message'
      placeholder,
      { fallback: 'fallback' placeholder
    )
    expect(message).toBe('error message')
  placeholder)

  it('uses fallback when no message can be extracted', () => {
    expect(buildAuthErrorMessage({placeholder, { fallback: 'fallback' placeholder)).toBe('fallback')
  placeholder)
placeholder)
