import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { get, post, put, del placeholder = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  del: vi.fn()
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: { get, post, put, delete: del placeholder
placeholder))

import {
  deleteOllamaCloudUsageSession,
  getOllamaCloudUsage,
  getOllamaCloudUsageSettings,
  refreshOllamaCloudUsage,
  saveOllamaCloudUsageSession,
  setOllamaCloudUsageAutoRefresh,
  updateOllamaCloudUsageSettings
placeholder from '@/api/admin/accounts'

const state = {
  account_id: 7,
  eligible: true,
  configured: true,
  auto_refresh_enabled: false,
  encryption_key_configured: true
placeholder

describe('admin Ollama Cloud usage API', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    put.mockReset()
    del.mockReset()
  placeholder)

  it('uses dedicated global settings endpoints', async () => {
    const settings = { enabled: false, interval_minutes: 60, debounce_minutes: 1 placeholder
    get.mockResolvedValueOnce({ data: settings placeholder)
    put.mockResolvedValueOnce({ data: settings placeholder)

    await expect(getOllamaCloudUsageSettings()).resolves.toEqual(settings)
    await expect(updateOllamaCloudUsageSettings(settings)).resolves.toEqual(settings)
    expect(get).toHaveBeenCalledWith('/admin/accounts/ollama-cloud-usage/settings')
    expect(put).toHaveBeenCalledWith('/admin/accounts/ollama-cloud-usage/settings', settings)
  placeholder)

  it('keeps session configuration write-only and separate from account updates', async () => {
    get.mockResolvedValueOnce({ data: state placeholder)
    put.mockResolvedValueOnce({ data: state placeholder).mockResolvedValueOnce({ data: state placeholder)
    del.mockResolvedValueOnce({ data: { ...state, configured: false placeholder placeholder)
    post.mockResolvedValueOnce({ data: state placeholder)

    await expect(getOllamaCloudUsage(7)).resolves.toEqual(state)
    await expect(saveOllamaCloudUsageSession(7, 'wos-session=secret')).resolves.toEqual(state)
    await expect(setOllamaCloudUsageAutoRefresh(7, true)).resolves.toEqual(state)
    await expect(refreshOllamaCloudUsage(7)).resolves.toEqual(state)
    await expect(deleteOllamaCloudUsageSession(7)).resolves.toMatchObject({ configured: false placeholder)

    expect(put).toHaveBeenNthCalledWith(1, '/admin/accounts/7/ollama-cloud-usage/session', { session: 'wos-session=secret' placeholder)
    expect(put).toHaveBeenNthCalledWith(2, '/admin/accounts/7/ollama-cloud-usage/auto-refresh', { enabled: true placeholder)
    expect(post).toHaveBeenCalledWith('/admin/accounts/7/ollama-cloud-usage/refresh')
    expect(del).toHaveBeenCalledWith('/admin/accounts/7/ollama-cloud-usage/session')
  placeholder)
placeholder)
