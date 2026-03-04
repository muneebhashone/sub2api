import { describe, expect, it placeholder from 'vitest'
import type { OpsErrorDetail placeholder from '@/api/admin/ops'
import { resolvePrimaryResponseBody, resolveUpstreamPayload placeholder from '../errorDetailResponse'

function makeDetail(overrides: Partial<OpsErrorDetail>): OpsErrorDetail {
  return {
    id: 1,
    created_at: '2026-01-01T00:00:00Z',
    phase: 'request',
    type: 'api_error',
    error_owner: 'platform',
    error_source: 'gateway',
    severity: 'P2',
    status_code: 502,
    platform: 'openai',
    model: 'gpt-4o-mini',
    is_retryable: true,
    retry_count: 0,
    resolved: false,
    client_request_id: 'crid-1',
    request_id: 'rid-1',
    message: 'Upstream request failed',
    user_email: 'user@example.com',
    account_name: 'acc',
    group_name: 'group',
    error_body: '',
    user_agent: '',
    request_body: '',
    request_body_truncated: false,
    is_business_limited: false,
    ...overrides
  placeholder
placeholder

describe('errorDetailResponse', () => {
  it('prefers upstream payload for request modal when error_body is generic gateway wrapper', () => {
    const detail = makeDetail({
      error_body: JSON.stringify({
        type: 'error',
        error: {
          type: 'upstream_error',
          message: 'Upstream request failed'
        placeholder
      placeholder),
      upstream_error_detail: '{"provider_message":"real upstream detail"placeholder'
    placeholder)

    expect(resolvePrimaryResponseBody(detail, 'request')).toBe('{"provider_message":"real upstream detail"placeholder')
  placeholder)

  it('keeps error_body for request modal when body is not generic wrapper', () => {
    const detail = makeDetail({
      error_body: JSON.stringify({
        type: 'error',
        error: {
          type: 'upstream_error',
          message: 'Upstream authentication failed, please contact administrator'
        placeholder
      placeholder),
      upstream_error_detail: '{"provider_message":"real upstream detail"placeholder'
    placeholder)

    expect(resolvePrimaryResponseBody(detail, 'request')).toBe(detail.error_body)
  placeholder)

  it('uses upstream payload first in upstream modal', () => {
    const detail = makeDetail({
      phase: 'upstream',
      upstream_error_message: 'provider 503 overloaded',
      error_body: '{"type":"error","error":{"type":"upstream_error","message":"Upstream request failed"placeholderplaceholder'
    placeholder)

    expect(resolvePrimaryResponseBody(detail, 'upstream')).toBe('provider 503 overloaded')
  placeholder)

  it('falls back to upstream payload when request error_body is empty', () => {
    const detail = makeDetail({
      error_body: '',
      upstream_error_message: 'dial tcp timeout'
    placeholder)

    expect(resolvePrimaryResponseBody(detail, 'request')).toBe('dial tcp timeout')
  placeholder)

  it('resolves upstream payload by detail -> events -> message priority', () => {
    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: 'detail payload',
      upstream_errors: '[{"message":"event payload"placeholder]',
      upstream_error_message: 'message payload'
    placeholder))).toBe('detail payload')

    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: '',
      upstream_errors: '[{"message":"event payload"placeholder]',
      upstream_error_message: 'message payload'
    placeholder))).toBe('[{"message":"event payload"placeholder]')

    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: '',
      upstream_errors: '',
      upstream_error_message: 'message payload'
    placeholder))).toBe('message payload')
  placeholder)

  it('treats empty JSON placeholders in upstream payload as empty', () => {
    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: '',
      upstream_errors: '[]',
      upstream_error_message: ''
    placeholder))).toBe('')

    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: '',
      upstream_errors: '{placeholder',
      upstream_error_message: ''
    placeholder))).toBe('')

    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: '',
      upstream_errors: 'null',
      upstream_error_message: ''
    placeholder))).toBe('')
  placeholder)

  it('skips placeholder candidates and falls back to the next upstream field', () => {
    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: '',
      upstream_errors: '[]',
      upstream_error_message: 'fallback message'
    placeholder))).toBe('fallback message')

    expect(resolveUpstreamPayload(makeDetail({
      upstream_error_detail: 'null',
      upstream_errors: '',
      upstream_error_message: 'fallback message'
    placeholder))).toBe('fallback message')
  placeholder)
placeholder)
