import { readFileSync placeholder from 'node:fs'
import { resolve placeholder from 'node:path'
import { describe, expect, it placeholder from 'vitest'

describe('GroupsView Composite route options', () => {
  it('offers Kimi, Zhipu GLM, and DeepSeek as route targets', () => {
    const source = readFileSync(resolve('src/views/admin/GroupsView.vue'), 'utf8')
    const options = source.slice(
      source.indexOf('const compositeRoutePlatformOptions'),
      source.indexOf('const compositeRouteEndpointOptions')
    )

    expect(options).toContain('{ value: "kimi", label: "Kimi" placeholder')
    expect(options).toContain('{ value: "zhipu", label: "Zhipu GLM" placeholder')
    expect(options).toContain('{ value: "deepseek", label: "DeepSeek" placeholder')
  placeholder)
placeholder)
