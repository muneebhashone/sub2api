import { describe, it, expect placeholder from 'vitest'
import {
  ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY,
  applyAntigravityProjectID,
  applyInterceptWarmup
placeholder from '../credentialsBuilder'

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

describe('applyAntigravityProjectID', () => {
  it('create + project id: trims and stores configured project fallback', () => {
    const creds: Record<string, unknown> = { access_token: 'tok' placeholder
    applyAntigravityProjectID(creds, '  configured-project  ', 'create')
    expect(creds[ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY]).toBe('configured-project')
  placeholder)

  it('create + empty project id: should not add the field', () => {
    const creds: Record<string, unknown> = { access_token: 'tok' placeholder
    applyAntigravityProjectID(creds, '   ', 'create')
    expect(ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY in creds).toBe(false)
  placeholder)

  it('edit + empty project id: deletes existing fallback', () => {
    const creds: Record<string, unknown> = {
      access_token: 'tok',
      [ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY]: 'old-project'
    placeholder
    applyAntigravityProjectID(creds, '', 'edit')
    expect(ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY in creds).toBe(false)
  placeholder)

  it('does not affect onboard project_id or other credentials', () => {
    const creds: Record<string, unknown> = {
      project_id: 'onboard-project',
      model_mapping: { 'gemini-*': 'gemini-2.5-flash' placeholder
    placeholder
    applyAntigravityProjectID(creds, 'configured-project', 'edit')
    expect(creds.project_id).toBe('onboard-project')
    expect(creds.model_mapping).toEqual({ 'gemini-*': 'gemini-2.5-flash' placeholder)
    expect(creds[ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY]).toBe('configured-project')
  placeholder)
placeholder)
