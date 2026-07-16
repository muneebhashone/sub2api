import { beforeEach, describe, expect, it, vi placeholder from 'vitest'
import { emptyEventFilters placeholder from '../viewModel'

const client = vi.hoisted(() => ({ get: vi.fn(), put: vi.fn(), post: vi.fn(), delete: vi.fn() placeholder))
vi.mock('@/api/client', () => ({ apiClient: client placeholder))

import promptAuditAPI from '../api'

describe('Prompt Audit API', () => {
  beforeEach(() => Object.values(client).forEach((mock) => mock.mockReset()))

  it('uses the independent admin route namespace', async () => {
    client.get.mockResolvedValue({ data: { config_version: 1 placeholder placeholder)
    await promptAuditAPI.getConfig()
    expect(client.get).toHaveBeenCalledWith('/admin/prompt-audit/config')

    client.get.mockResolvedValue({ data: { process_status: 'running' placeholder placeholder)
    await promptAuditAPI.getRuntime()
    expect(client.get).toHaveBeenCalledWith('/admin/prompt-audit/runtime')
  placeholder)

  it('sends a temporary probe token only in the request and never invents response credentials', async () => {
    client.post.mockResolvedValue({ data: { ok: true, token_applied: true placeholder placeholder)
    const result = await promptAuditAPI.probeEndpoint({
      id: 'guard-1', name: 'Guard', protocol: 'openai_compatible', base_url: 'http://127.0.0.1:8000', model: 'guard',
      token: 'api-canary-secret', clear_token: false, timeout_ms: 1000, input_limit: 1000, enabled: true, has_token: false, token_status: 'missing',
    placeholder)
    expect(client.post).toHaveBeenCalledWith('/admin/prompt-audit/endpoints/probe', expect.objectContaining({ endpoint: expect.objectContaining({ token: 'api-canary-secret' placeholder) placeholder))
    expect(JSON.stringify(result)).not.toContain('api-canary-secret')
  placeholder)

  it('passes a server preview token through the confirmed filter-delete contract', async () => {
    client.post.mockResolvedValue({ data: { deleted_events: 2, deleted_jobs: 2 placeholder placeholder)
    const filters = emptyEventFilters()
    filters.start_at = '2026-07-15T00:00'
    filters.end_at = '2026-07-16T00:00'
    await promptAuditAPI.deleteEventsByFilter(filters, {
      matched_count: 2, filter_summary: {placeholder, snapshot_max_id: 10, filter_hash: 'a'.repeat(64), confirmation_token: 'opaque-token', expires_at: '2026-07-16T00:05:00Z',
    placeholder)
    expect(client.post).toHaveBeenCalledWith('/admin/prompt-audit/events/delete-by-filter', expect.objectContaining({
      snapshot_max_id: 10, filter_hash: 'a'.repeat(64), confirmation_token: 'opaque-token', confirm: true,
    placeholder))
  placeholder)
placeholder)
