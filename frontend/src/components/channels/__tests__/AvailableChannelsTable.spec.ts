import { readFileSync placeholder from 'node:fs'
import { dirname, resolve placeholder from 'node:path'
import { fileURLToPath placeholder from 'node:url'

import { describe, expect, it placeholder from 'vitest'

const componentPath = resolve(dirname(fileURLToPath(import.meta.url)), '../AvailableChannelsTable.vue')
const componentSource = readFileSync(componentPath, 'utf8')

describe('AvailableChannelsTable scroll integration', () => {
  // #4555：根元素必须是 TablePageLayout 滚动链约定的 .table-wrapper，
  // 否则内容超出视口高度时被外层 overflow-hidden 裁剪且没有滚动条。
  it('mounts the table on the .table-wrapper scroll hook', () => {
    expect(componentSource).toMatch(/<div class="table-wrapper">\s*<table/)
  placeholder)

  it('does not clip content with its own overflow-hidden card wrapper', () => {
    expect(componentSource).not.toMatch(/<div class="card overflow-hidden">/)
  placeholder)
placeholder)
