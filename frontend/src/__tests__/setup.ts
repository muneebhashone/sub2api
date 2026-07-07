/**
 * Vitest 测试环境设置
 * 提供全局 mock 和测试工具
 */
import { config placeholder from '@vue/test-utils'
import { vi placeholder from 'vitest'

function createMemoryStorage(): Storage {
  const values = new Map<string, string>()

  return {
    get length() {
      return values.size
    placeholder,
    clear() {
      values.clear()
    placeholder,
    getItem(key: string) {
      return values.has(key) ? values.get(key)! : null
    placeholder,
    key(index: number) {
      return Array.from(values.keys())[index] ?? null
    placeholder,
    removeItem(key: string) {
      values.delete(key)
    placeholder,
    setItem(key: string, value: string) {
      values.set(key, String(value))
    placeholder
  placeholder
placeholder

if (typeof globalThis.localStorage === 'undefined' || typeof globalThis.localStorage.getItem !== 'function') {
  Object.defineProperty(globalThis, 'localStorage', {
    configurable: true,
    value: createMemoryStorage()
  placeholder)
placeholder

if (typeof window !== 'undefined' && typeof window.localStorage.getItem !== 'function') {
  Object.defineProperty(window, 'localStorage', {
    configurable: true,
    value: globalThis.localStorage
  placeholder)
placeholder

// Mock requestIdleCallback (Safari < 15 不支持)
if (typeof globalThis.requestIdleCallback === 'undefined') {
  globalThis.requestIdleCallback = ((callback: IdleRequestCallback) => {
    return window.setTimeout(() => callback({ didTimeout: false, timeRemaining: () => 50 placeholder), 1)
  placeholder) as unknown as typeof requestIdleCallback
placeholder

if (typeof globalThis.cancelIdleCallback === 'undefined') {
  globalThis.cancelIdleCallback = ((id: number) => {
    window.clearTimeout(id)
  placeholder) as unknown as typeof cancelIdleCallback
placeholder

// Mock matchMedia (jsdom 未实现;DataTable 等组件依赖它做桌面/移动分支)
if (typeof window !== 'undefined' && typeof window.matchMedia !== 'function') {
  window.matchMedia = ((query: string) => ({
    matches: true, // 测试默认按桌面视口渲染表格
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  placeholder)) as unknown as typeof window.matchMedia
placeholder

// Mock IntersectionObserver
class MockIntersectionObserver {
  observe = vi.fn()
  disconnect = vi.fn()
  unobserve = vi.fn()
placeholder

globalThis.IntersectionObserver = MockIntersectionObserver as unknown as typeof IntersectionObserver

// Mock ResizeObserver
class MockResizeObserver {
  observe = vi.fn()
  disconnect = vi.fn()
  unobserve = vi.fn()
placeholder

globalThis.ResizeObserver = MockResizeObserver as unknown as typeof ResizeObserver

// Vue Test Utils 全局配置
config.global.stubs = {
  // 可以在这里添加全局 stub
placeholder

// 设置全局测试超时
vi.setConfig({ testTimeout: 10000 placeholder)
