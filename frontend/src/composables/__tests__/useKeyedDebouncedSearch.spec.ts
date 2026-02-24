import { beforeEach, afterEach, describe, expect, it, vi placeholder from 'vitest'
import { useKeyedDebouncedSearch placeholder from '@/composables/useKeyedDebouncedSearch'

const flushPromises = () => Promise.resolve()

describe('useKeyedDebouncedSearch', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  placeholder)

  afterEach(() => {
    vi.useRealTimers()
  placeholder)

  it('为不同 key 独立防抖触发搜索', async () => {
    const search = vi.fn().mockResolvedValue([])
    const onSuccess = vi.fn()

    const searcher = useKeyedDebouncedSearch<string[]>({
      delay: 100,
      search,
      onSuccess
    placeholder)

    searcher.trigger('a', 'foo')
    searcher.trigger('b', 'bar')

    expect(search).not.toHaveBeenCalled()

    vi.advanceTimersByTime(100)
    await flushPromises()

    expect(search).toHaveBeenCalledTimes(2)
    expect(search).toHaveBeenNthCalledWith(
      1,
      'foo',
      expect.objectContaining({ key: 'a', signal: expect.any(AbortSignal) placeholder)
    )
    expect(search).toHaveBeenNthCalledWith(
      2,
      'bar',
      expect.objectContaining({ key: 'b', signal: expect.any(AbortSignal) placeholder)
    )
    expect(onSuccess).toHaveBeenCalledTimes(2)
  placeholder)

  it('同 key 新请求会取消旧请求并忽略过期响应', async () => {
    const resolves: Array<(value: string[]) => void> = []
    const search = vi.fn().mockImplementation(
      () => new Promise<string[]>((resolve) => {
        resolves.push(resolve)
      placeholder)
    )
    const onSuccess = vi.fn()

    const searcher = useKeyedDebouncedSearch<string[]>({
      delay: 50,
      search,
      onSuccess
    placeholder)

    searcher.trigger('rule-1', 'first')
    vi.advanceTimersByTime(50)
    await flushPromises()

    searcher.trigger('rule-1', 'second')
    vi.advanceTimersByTime(50)
    await flushPromises()

    expect(search).toHaveBeenCalledTimes(2)

    resolves[1](['second'])
    await flushPromises()
    expect(onSuccess).toHaveBeenCalledTimes(1)
    expect(onSuccess).toHaveBeenLastCalledWith('rule-1', ['second'])

    resolves[0](['first'])
    await flushPromises()
    expect(onSuccess).toHaveBeenCalledTimes(1)
  placeholder)

  it('clearKey 会取消未执行任务', () => {
    const search = vi.fn().mockResolvedValue([])
    const onSuccess = vi.fn()

    const searcher = useKeyedDebouncedSearch<string[]>({
      delay: 100,
      search,
      onSuccess
    placeholder)

    searcher.trigger('a', 'foo')
    searcher.clearKey('a')

    vi.advanceTimersByTime(100)

    expect(search).not.toHaveBeenCalled()
    expect(onSuccess).not.toHaveBeenCalled()
  placeholder)
placeholder)
