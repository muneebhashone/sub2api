import { describe, it, expect placeholder from 'vitest'
import { applyInterceptWarmup placeholder from '../credentialsBuilder'

describe('applyInterceptWarmup', () => {
  it('create + enabled=true: should set intercept_warmup_requests to true', () => {
    const creds: Record<string, unknown> = { access_token: 'tok' placeholder
    applyInterceptWarmup(creds, true, 'create')
    expect(creds.intercept_warmup_requests).toBe(true)
  placeholder)

  it('create + enabled=false: should not add the field', () => {
    const creds: Record<string, unknown> = { access_token: 'tok' placeholder
    applyInterceptWarmup(creds, false, 'create')
    expect('intercept_warmup_requests' in creds).toBe(false)
  placeholder)

  it('edit + enabled=true: should set intercept_warmup_requests to true', () => {
    const creds: Record<string, unknown> = { api_key: 'sk' placeholder
    applyInterceptWarmup(creds, true, 'edit')
    expect(creds.intercept_warmup_requests).toBe(true)
  placeholder)

  it('edit + enabled=false + field exists: should delete the field', () => {
    const creds: Record<string, unknown> = { api_key: 'sk', intercept_warmup_requests: true placeholder
    applyInterceptWarmup(creds, false, 'edit')
    expect('intercept_warmup_requests' in creds).toBe(false)
  placeholder)

  it('edit + enabled=false + field absent: should not throw', () => {
    const creds: Record<string, unknown> = { api_key: 'sk' placeholder
    applyInterceptWarmup(creds, false, 'edit')
    expect('intercept_warmup_requests' in creds).toBe(false)
  placeholder)

  it('should not affect other fields', () => {
    const creds: Record<string, unknown> = {
      api_key: 'sk',
      base_url: 'url',
      intercept_warmup_requests: true
    placeholder
    applyInterceptWarmup(creds, false, 'edit')
    expect(creds.api_key).toBe('sk')
    expect(creds.base_url).toBe('url')
    expect('intercept_warmup_requests' in creds).toBe(false)
  placeholder)
placeholder)
