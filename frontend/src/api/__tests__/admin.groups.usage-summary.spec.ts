import { beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { get placeholder = vi.hoisted(() => ({
  get: vi.fn(),
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: { get placeholder,
placeholder))

import { getUsageSummary placeholder from '@/api/admin/groups'

describe('admin group usage summary API', () => {
  beforeEach(() => {
    get.mockReset()
    get.mockResolvedValue({ data: [] placeholder)
  placeholder)

  it('does not send browser timezone parameters', async () => {
    const summary = [
      { group_id: 1, today_cost: 1.25, yesterday_cost: 2.5, total_cost: 9.75 placeholder,
    ]
    get.mockResolvedValue({ data: summary placeholder)

    await expect(getUsageSummary()).resolves.toEqual(summary)

    expect(get).toHaveBeenCalledWith('/admin/groups/usage-summary')
  placeholder)
placeholder)
