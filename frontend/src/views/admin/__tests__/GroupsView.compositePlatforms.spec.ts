import { describe, expect, it placeholder from 'vitest'
import { CONCRETE_PLATFORM_OPTIONS placeholder from '@/constants/platforms'

describe('GroupsView Composite route options', () => {
  it('offers Kimi, Zhipu GLM, and DeepSeek as route targets', () => {
    expect(CONCRETE_PLATFORM_OPTIONS.map((option) => option.value)).toEqual(
      expect.arrayContaining(['kimi', 'zhipu', 'deepseek'])
    )
  placeholder)
placeholder)
