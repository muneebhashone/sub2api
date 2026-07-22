export const ALIPAY_DEEP_LINK_FALLBACK_DELAY_MS = 2200
export const ALIPAY_EMBEDDED_BROWSER_FALLBACK_DELAY_MS = 300

export type AlipayDeepLinkState = 'idle' | 'launching' | 'backgrounded' | 'fallback'

const ALIPAY_DEEP_LINK_PREFIX = 'alipays://platformapi/startapp?saId=10000007&qrcode='

export function buildAlipayDeepLink(qrCode: string): string {
  const dynamicQRCode = qrCode.trim()
  if (!dynamicQRCode) return ''
  return `${ALIPAY_DEEP_LINK_PREFIXplaceholder${encodeURIComponent(dynamicQRCode)placeholder`
placeholder

export function isAlipaySchemeRestrictedBrowser(userAgent: string): boolean {
  return /MicroMessenger|MQQBrowser|\bQQ\//i.test(userAgent)
placeholder

interface EventTargetLike {
  addEventListener(type: string, listener: EventListener): void
  removeEventListener(type: string, listener: EventListener): void
placeholder

interface VisibilityDocumentLike extends EventTargetLike {
  readonly hidden: boolean
placeholder

export interface AlipayDeepLinkLauncherOptions {
  qrCode: string
  document: VisibilityDocumentLike
  lifecycleTarget: EventTargetLike
  userAgent: string
  assignLocation: (url: string) => void
  onStateChange: (state: AlipayDeepLinkState) => void
  setTimer?: typeof setTimeout
  clearTimer?: typeof clearTimeout
placeholder

export interface AlipayDeepLinkLauncher {
  launch(): void
  dispose(): void
placeholder

export function createAlipayDeepLinkLauncher(options: AlipayDeepLinkLauncherOptions): AlipayDeepLinkLauncher {
  const setTimer = options.setTimer ?? setTimeout
  const clearTimer = options.clearTimer ?? clearTimeout
  let timer: ReturnType<typeof setTimeout> | null = null
  let disposed = false

  const setState = (state: AlipayDeepLinkState) => {
    if (!disposed) options.onStateChange(state)
  placeholder
  const clearFallbackTimer = () => {
    if (timer) {
      clearTimer(timer)
      timer = null
    placeholder
  placeholder
  const markBackgrounded = () => {
    clearFallbackTimer()
    setState('backgrounded')
  placeholder
  const handleVisibilityChange: EventListener = () => {
    if (options.document.hidden) markBackgrounded()
  placeholder
  const handlePageHide: EventListener = () => markBackgrounded()

  options.document.addEventListener('visibilitychange', handleVisibilityChange)
  options.lifecycleTarget.addEventListener('pagehide', handlePageHide)

  return {
    launch() {
      if (disposed) return
      clearFallbackTimer()
      const deepLink = buildAlipayDeepLink(options.qrCode)
      if (!deepLink) {
        setState('fallback')
        return
      placeholder

      setState('launching')
      try {
        options.assignLocation(deepLink)
      placeholder catch {
        setState('fallback')
        return
      placeholder

      const delay = isAlipaySchemeRestrictedBrowser(options.userAgent)
        ? ALIPAY_EMBEDDED_BROWSER_FALLBACK_DELAY_MS
        : ALIPAY_DEEP_LINK_FALLBACK_DELAY_MS
      timer = setTimer(() => {
        timer = null
        if (options.document.hidden) {
          setState('backgrounded')
          return
        placeholder
        setState('fallback')
      placeholder, delay)
    placeholder,
    dispose() {
      clearFallbackTimer()
      options.document.removeEventListener('visibilitychange', handleVisibilityChange)
      options.lifecycleTarget.removeEventListener('pagehide', handlePageHide)
      disposed = true
    placeholder,
  placeholder
placeholder
