import { describe, expect, it placeholder from 'vitest'
import { baseCompile placeholder from '@intlify/message-compiler'

import en from '../locales/en'
import zh from '../locales/zh'

// vue-i18n 在运行时才编译消息：文案里未转义的花括号（如内嵌 JSON 示例
// "{\"user-agent\": ...placeholder"）会在渲染时抛 "Invalid token in placeholder"，
// 直接炸掉整个组件树，且构建期完全无感。本测试把全部文案预编译一遍，
// 将该类问题固化为显式失败。字面量花括号请用 {'{'placeholder / {'placeholder'placeholder 转义，
// 或将语言中立的示例文本（如 JSON）移出 i18n。
function collectCompileErrors(node: unknown, path: string, out: string[]): void {
  if (typeof node === 'string') {
    baseCompile(node, {
      onError: (err) => {
        out.push(`${pathplaceholder: ${err.messageplaceholder`)
      placeholder
    placeholder)
    return
  placeholder
  if (Array.isArray(node)) {
    node.forEach((item, index) => collectCompileErrors(item, `${pathplaceholder[${indexplaceholder]`, out))
    return
  placeholder
  if (node && typeof node === 'object') {
    for (const [key, value] of Object.entries(node as Record<string, unknown>)) {
      collectCompileErrors(value, path ? `${pathplaceholder.${keyplaceholder` : key, out)
    placeholder
  placeholder
placeholder

describe('locale messages compile', () => {
  it.each([
    ['zh', zh],
    ['en', en]
  ] as const)('%s messages all compile without placeholder errors', (locale, messages) => {
    const errors: string[] = []
    collectCompileErrors(messages, locale, errors)
    expect(errors).toEqual([])
  placeholder)
placeholder)
