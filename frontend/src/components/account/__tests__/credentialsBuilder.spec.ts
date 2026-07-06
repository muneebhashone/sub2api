import { describe, it, expect placeholder from 'vitest'
import {
  ANTIGRAVITY_PROJECT_ID_CREDENTIAL_KEY,
  HEADER_OVERRIDE_ENABLED_CREDENTIAL_KEY,
  HEADER_OVERRIDES_CREDENTIAL_KEY,
  applyAntigravityProjectID,
  applyHeaderOverride,
  applyInterceptWarmup,
  buildHeaderOverridesObject,
  getHeaderOverrideTemplate,
  isHeaderOverridePlatform,
  splitHeaderOverridesObject,
  validateHeaderOverrideRows
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

describe('isHeaderOverridePlatform', () => {
  it('only anthropic and openai are supported', () => {
    expect(isHeaderOverridePlatform('anthropic')).toBe(true)
    expect(isHeaderOverridePlatform('openai')).toBe(true)
    expect(isHeaderOverridePlatform('gemini')).toBe(false)
    expect(isHeaderOverridePlatform('grok')).toBe(false)
    expect(isHeaderOverridePlatform('antigravity')).toBe(false)
    expect(isHeaderOverridePlatform('')).toBe(false)
  placeholder)
placeholder)

describe('validateHeaderOverrideRows', () => {
  it('accepts valid rows and empty placeholder rows', () => {
    expect(
      validateHeaderOverrideRows([
        { name: 'user-agent', value: 'my-agent/1.0' placeholder,
        { name: 'x-app', value: '' placeholder,
        { name: '', value: '' placeholder
      ])
    ).toBeNull()
  placeholder)

  it('rejects empty name with non-empty value', () => {
    expect(validateHeaderOverrideRows([{ name: '', value: 'v' placeholder])).toBe('invalidName')
  placeholder)

  it('rejects invalid header names', () => {
    expect(validateHeaderOverrideRows([{ name: 'bad name', value: '' placeholder])).toBe('invalidName')
    expect(validateHeaderOverrideRows([{ name: 'bad:name', value: '' placeholder])).toBe('invalidName')
    expect(validateHeaderOverrideRows([{ name: '名称', value: '' placeholder])).toBe('invalidName')
  placeholder)

  it('rejects blocked header names case-insensitively', () => {
    expect(validateHeaderOverrideRows([{ name: 'Authorization', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'X-Api-Key', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'host', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'Content-Length', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'Content-Type', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'Cookie', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'x-goog-api-key', value: '' placeholder])).toBe('blockedName')
  placeholder)

  it('rejects duplicate names case-insensitively', () => {
    expect(
      validateHeaderOverrideRows([
        { name: 'User-Agent', value: 'a' placeholder,
        { name: 'user-agent', value: 'b' placeholder
      ])
    ).toBe('duplicateName')
  placeholder)
placeholder)

describe('buildHeaderOverridesObject / splitHeaderOverridesObject', () => {
  it('lowercases names, trims values and drops empty-name rows', () => {
    expect(
      buildHeaderOverridesObject([
        { name: ' User-Agent ', value: ' my-agent ' placeholder,
        { name: 'X-App', value: '' placeholder,
        { name: '', value: 'ignored' placeholder
      ])
    ).toEqual({ 'user-agent': 'my-agent', 'x-app': '' placeholder)
  placeholder)

  it('splits an object into sorted rows and ignores non-string values', () => {
    expect(
      splitHeaderOverridesObject({ 'x-app': 'cli', 'user-agent': 'ua', bogus: 42 placeholder)
    ).toEqual([
      { name: 'user-agent', value: 'ua' placeholder,
      { name: 'x-app', value: 'cli' placeholder
    ])
    expect(splitHeaderOverridesObject(null)).toEqual([])
    expect(splitHeaderOverridesObject(['a'])).toEqual([])
    expect(splitHeaderOverridesObject('str')).toEqual([])
  placeholder)

  it('roundtrips through build and split', () => {
    const rows = [
      { name: 'user-agent', value: 'ua' placeholder,
      { name: 'x-app', value: 'cli' placeholder
    ]
    expect(splitHeaderOverridesObject(buildHeaderOverridesObject(rows))).toEqual(rows)
  placeholder)
placeholder)

describe('getHeaderOverrideTemplate', () => {
  it('returns Claude Code CLI headers with empty values for anthropic', () => {
    const rows = getHeaderOverrideTemplate('anthropic')
    expect(rows.every((r) => r.value === '')).toBe(true)
    const names = rows.map((r) => r.name)
    expect(names).toContain('user-agent')
    expect(names).toContain('x-app')
    expect(names).toContain('anthropic-beta')
    expect(names).toContain('x-stainless-lang')
    expect(validateHeaderOverrideRows(rows)).toBeNull()
  placeholder)

  it('returns Codex CLI headers with empty values for openai', () => {
    const rows = getHeaderOverrideTemplate('openai')
    expect(rows.every((r) => r.value === '')).toBe(true)
    const names = rows.map((r) => r.name)
    expect(names).toContain('user-agent')
    expect(names).toContain('originator')
    expect(names).toContain('openai-beta')
    expect(validateHeaderOverrideRows(rows)).toBeNull()
  placeholder)
placeholder)

describe('applyHeaderOverride', () => {
  it('create + enabled: writes enabled flag and overrides object', () => {
    const creds: Record<string, unknown> = { api_key: 'sk' placeholder
    applyHeaderOverride(creds, true, [{ name: 'User-Agent', value: 'ua' placeholder], 'create')
    expect(creds[HEADER_OVERRIDE_ENABLED_CREDENTIAL_KEY]).toBe(true)
    expect(creds[HEADER_OVERRIDES_CREDENTIAL_KEY]).toEqual({ 'user-agent': 'ua' placeholder)
  placeholder)

  it('create + disabled: does not add fields', () => {
    const creds: Record<string, unknown> = { api_key: 'sk' placeholder
    applyHeaderOverride(creds, false, [{ name: 'user-agent', value: 'ua' placeholder], 'create')
    expect(HEADER_OVERRIDE_ENABLED_CREDENTIAL_KEY in creds).toBe(false)
    expect(HEADER_OVERRIDES_CREDENTIAL_KEY in creds).toBe(false)
  placeholder)

  it('edit + disabled: deletes existing fields', () => {
    const creds: Record<string, unknown> = {
      api_key: 'sk',
      [HEADER_OVERRIDE_ENABLED_CREDENTIAL_KEY]: true,
      [HEADER_OVERRIDES_CREDENTIAL_KEY]: { 'user-agent': 'ua' placeholder
    placeholder
    applyHeaderOverride(creds, false, [], 'edit')
    expect(HEADER_OVERRIDE_ENABLED_CREDENTIAL_KEY in creds).toBe(false)
    expect(HEADER_OVERRIDES_CREDENTIAL_KEY in creds).toBe(false)
    expect(creds.api_key).toBe('sk')
  placeholder)

  it('edit + enabled: replaces overrides object wholesale', () => {
    const creds: Record<string, unknown> = {
      [HEADER_OVERRIDE_ENABLED_CREDENTIAL_KEY]: true,
      [HEADER_OVERRIDES_CREDENTIAL_KEY]: { 'x-old': 'old' placeholder
    placeholder
    applyHeaderOverride(creds, true, [{ name: 'x-new', value: 'new' placeholder], 'edit')
    expect(creds[HEADER_OVERRIDES_CREDENTIAL_KEY]).toEqual({ 'x-new': 'new' placeholder)
  placeholder)
placeholder)

describe('validateHeaderOverrideRows value/entry limits', () => {
  it('rejects websocket handshake headers', () => {
    expect(validateHeaderOverrideRows([{ name: 'Sec-WebSocket-Key', value: '' placeholder])).toBe(
      'blockedName'
    )
  placeholder)

  it('rejects control characters in values', () => {
    expect(validateHeaderOverrideRows([{ name: 'x-app', value: 'a\x0bb' placeholder])).toBe('invalidValue')
  placeholder)

  it('rejects oversized values', () => {
    expect(validateHeaderOverrideRows([{ name: 'x-app', value: 'a'.repeat(8193) placeholder])).toBe(
      'invalidValue'
    )
  placeholder)

  it('measures value length in UTF-8 bytes to match backend', () => {
    // 3000 个 CJK 字符 = 3000 UTF-16 code units，但 9000 UTF-8 字节 > 8192
    expect(validateHeaderOverrideRows([{ name: 'x-app', value: '测'.repeat(3000) placeholder])).toBe(
      'invalidValue'
    )
    expect(validateHeaderOverrideRows([{ name: 'x-app', value: '测'.repeat(2000) placeholder])).toBeNull()
  placeholder)

  it('rejects too many entries', () => {
    const rows = Array.from({ length: 65 placeholder, (_, i) => ({ name: `x-h-${iplaceholder`, value: 'v' placeholder))
    expect(validateHeaderOverrideRows(rows)).toBe('tooManyEntries')
  placeholder)
placeholder)

describe('validateHeaderOverrideRows session isolation headers', () => {
  it('rejects per-request session headers', () => {
    expect(validateHeaderOverrideRows([{ name: 'session_id', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'Conversation_ID', value: '' placeholder])).toBe('blockedName')
    expect(validateHeaderOverrideRows([{ name: 'x-codex-turn-state', value: '' placeholder])).toBe(
      'blockedName'
    )
    expect(validateHeaderOverrideRows([{ name: 'X-Claude-Code-Session-Id', value: '' placeholder])).toBe(
      'blockedName'
    )
    expect(validateHeaderOverrideRows([{ name: 'x-client-request-id', value: '' placeholder])).toBe(
      'blockedName'
    )
  placeholder)

  it('allows tab inside value', () => {
    expect(validateHeaderOverrideRows([{ name: 'x-app', value: 'a\tb' placeholder])).toBeNull()
  placeholder)

  it('rejects oversized names', () => {
    expect(validateHeaderOverrideRows([{ name: 'x'.repeat(201), value: 'v' placeholder])).toBe('invalidName')
  placeholder)
placeholder)
