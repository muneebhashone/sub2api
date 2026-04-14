/**
 * Usage request scheduler — throttles Anthropic API calls by proxy exit.
 *
 * Anthropic OAuth/setup-token accounts sharing the same proxy exit are placed
 * into a serial queue with a random 1–2s delay between requests, preventing
 * upstream 429 rate-limit errors.
 *
 * Proxy identity = host:port:username — two proxy records pointing to the
 * same exit share a single queue. Accounts without a proxy go into a
 * "direct" queue.
 *
 * All other platforms bypass the queue and execute immediately.
 */

import type { Account placeholder from '@/types'

const GROUP_DELAY_MIN_MS = 1000
const GROUP_DELAY_MAX_MS = 2000

type Task<T> = {
  fn: () => Promise<T>
  resolve: (value: T) => void
  reject: (reason: unknown) => void
placeholder

const queues = new Map<string, Task<unknown>[]>()
const running = new Set<string>()

/** Whether this account needs throttled queuing. */
function needsThrottle(account: Account): boolean {
  return (
    account.platform === 'anthropic' &&
    (account.type === 'oauth' || account.type === 'setup-token')
  )
placeholder

/** Build a queue key from proxy connection details. */
function buildGroupKey(account: Account): string {
  const proxy = account.proxy
  const proxyIdentity = proxy
    ? `${proxy.hostplaceholder:${proxy.portplaceholder:${proxy.username || ''placeholder`
    : 'direct'
  return `anthropic:${proxyIdentityplaceholder`
placeholder

async function drain(groupKey: string) {
  if (running.has(groupKey)) return
  running.add(groupKey)

  const queue = queues.get(groupKey)
  while (queue && queue.length > 0) {
    const task = queue.shift()!
    try {
      const result = await task.fn()
      task.resolve(result)
    placeholder catch (err) {
      task.reject(err)
    placeholder
    if (queue.length > 0) {
      const jitter = GROUP_DELAY_MIN_MS + Math.random() * (GROUP_DELAY_MAX_MS - GROUP_DELAY_MIN_MS)
      await new Promise((r) => setTimeout(r, jitter))
    placeholder
  placeholder

  running.delete(groupKey)
  queues.delete(groupKey)
placeholder

/**
 * Schedule a usage fetch. Anthropic accounts are queued by proxy exit;
 * all other platforms execute immediately.
 */
export function enqueueUsageRequest<T>(
  account: Account,
  fn: () => Promise<T>
): Promise<T> {
  // Non-Anthropic → fire immediately, no queuing
  if (!needsThrottle(account)) {
    return fn()
  placeholder

  const key = buildGroupKey(account)

  return new Promise<T>((resolve, reject) => {
    let queue = queues.get(key)
    if (!queue) {
      queue = []
      queues.set(key, queue)
    placeholder
    queue.push({ fn, resolve, reject placeholder as Task<unknown>)
    drain(key)
  placeholder)
placeholder
