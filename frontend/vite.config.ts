import { defineConfig, loadEnv, Plugin placeholder from 'vite'
import vue from '@vitejs/plugin-vue'
import checker from 'vite-plugin-checker'
import { resolve placeholder from 'path'
import { Buffer placeholder from 'node:buffer'
import { lookup placeholder from 'node:dns/promises'
import type { IncomingMessage, ServerResponse placeholder from 'node:http'
import { isIP placeholder from 'node:net'

/**
 * Vite 插件：开发模式下注入公开配置到 index.html
 * 与生产模式的后端注入行为保持一致，消除闪烁
 */
function injectPublicSettings(backendUrl: string): Plugin {
  return {
    name: 'inject-public-settings',
    apply: 'serve',
    transformIndexHtml: {
      order: 'pre',
      async handler(html) {
        try {
          const response = await fetch(`${backendUrlplaceholder/api/v1/settings/public`, {
            signal: AbortSignal.timeout(2000)
          placeholder)
          if (response.ok) {
            const data = await response.json()
            if (data.code === 0 && data.data) {
              const script = `<script>window.__APP_CONFIG__=${JSON.stringify(data.data)placeholder;</script>`
              return html.replace('</head>', `${scriptplaceholder\n</head>`)
            placeholder
          placeholder
        placeholder catch (e) {
          console.warn('[vite] 无法获取公开配置，将回退到 API 调用:', (e as Error).message)
        placeholder
        return html
      placeholder
    placeholder
  placeholder
placeholder

const LLM_TESTER_MAX_BODY_BYTES = 12 * 1024 * 1024
const LLM_TESTER_TIMEOUT_MS = 300000

function llmTesterDevProxy(): Plugin {
  return {
    name: 'llm-tester-dev-proxy',
    apply: 'serve',
    configureServer(server) {
      server.middlewares.use(async (req, res, next) => {
        const pathname = new URL(req.url || '/', 'http://localhost').pathname
        if (req.method !== 'POST' || !pathname.startsWith('/api/v1/llm-tester/')) {
          next()
          return
        placeholder

        try {
          const body = await readDevProxyJson(req)
          const route = pathname.slice('/api/v1/llm-tester/'.length)
          if (route === 'models') {
            await forwardDevLLMTesterRequest(res, body, 'GET', 'models')
            return
          placeholder
          if (route === 'chat/completions') {
            await forwardDevLLMTesterRequest(res, body, 'POST', 'chat/completions')
            return
          placeholder
          if (route === 'images/generations') {
            await forwardDevLLMTesterRequest(res, body, 'POST', 'images/generations')
            return
          placeholder
          if (route === 'responses') {
            await forwardDevLLMTesterRequest(res, body, 'POST', 'responses')
            return
          placeholder
          next()
        placeholder catch (error) {
          writeDevProxyError(res, 502, devProxyErrorMessage(error))
        placeholder
      placeholder)
    placeholder
  placeholder
placeholder

async function readDevProxyJson(req: IncomingMessage): Promise<Record<string, any>> {
  const chunks: Buffer[] = []
  let total = 0
  for await (const chunk of req) {
    const buffer = Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk)
    total += buffer.length
    if (total > LLM_TESTER_MAX_BODY_BYTES) {
      throw new Error('request body is too large')
    placeholder
    chunks.push(buffer)
  placeholder

  try {
    const parsed = JSON.parse(Buffer.concat(chunks).toString('utf8'))
    return parsed && typeof parsed === 'object' ? parsed : {placeholder
  placeholder catch {
    throw new Error('invalid request body')
  placeholder
placeholder

async function forwardDevLLMTesterRequest(
  res: ServerResponse,
  body: Record<string, any>,
  method: 'GET' | 'POST',
  resource: 'models' | 'chat/completions' | 'images/generations' | 'responses'
) {
  const baseUrl = String(body.base_url || '').trim()
  const apiKey = String(body.api_key || '').trim()
  if (!baseUrl) {
    writeDevProxyError(res, 400, 'base_url is required')
    return
  placeholder
  if (!apiKey) {
    writeDevProxyError(res, 400, 'api_key is required')
    return
  placeholder
  if (apiKey.length > 8192) {
    writeDevProxyError(res, 400, 'api_key is too long')
    return
  placeholder
  if (method === 'POST' && !body.payload) {
    writeDevProxyError(res, 400, 'payload is required')
    return
  placeholder

  const endpoint = await buildDevLLMTesterEndpoint(baseUrl, resource)
  const upstream = await fetch(endpoint, {
    method,
    headers: {
      Authorization: `Bearer ${apiKeyplaceholder`,
      Accept: 'application/json',
      'Content-Type': 'application/json',
      'User-Agent': 'Sub2API-LLM-Tester/1.0',
      'X-Title': 'Sub2API LLM Tester'
    placeholder,
    body: method === 'POST' ? JSON.stringify(body.payload || {placeholder) : undefined,
    signal: AbortSignal.timeout(LLM_TESTER_TIMEOUT_MS)
  placeholder)
  const payload = Buffer.from(await upstream.arrayBuffer())
  if (payload.length > LLM_TESTER_MAX_BODY_BYTES) {
    writeDevProxyError(res, 502, 'upstream response is too large')
    return
  placeholder

  res.statusCode = upstream.status
  res.setHeader('Content-Type', upstream.headers.get('content-type') || 'application/json')
  res.end(payload)
placeholder

async function buildDevLLMTesterEndpoint(baseUrl: string, resource: 'models' | 'chat/completions' | 'images/generations' | 'responses'): Promise<string> {
  const url = new URL(baseUrl.replace(/\/+$/, ''))
  if (url.protocol !== 'https:') {
    throw new Error('base_url must use https')
  placeholder
  if (url.username || url.password) {
    throw new Error('base_url must not include user info')
  placeholder
  await assertDevProxyPublicHost(url.hostname)
  url.search = ''
  url.hash = ''
  if (!/\/v\d+$/i.test(url.pathname)) {
    url.pathname = `${url.pathname.replace(/\/+$/, '')placeholder/v1`
  placeholder
  url.pathname = `${url.pathname.replace(/\/+$/, '')placeholder/${resourceplaceholder`
  return url.toString()
placeholder

async function assertDevProxyPublicHost(hostname: string) {
  const host = hostname.trim().toLowerCase()
  if (isBlockedDevProxyHost(host)) {
    throw new Error(`host is not allowed: ${hostnameplaceholder`)
  placeholder
  if (isIP(host)) {
    if (isBlockedDevProxyIP(host)) throw new Error(`host is not allowed: ${hostnameplaceholder`)
    return
  placeholder
  const addrs = await lookup(host, { all: true, verbatim: false placeholder)
  if (!addrs.length) {
    throw new Error(`host did not resolve: ${hostnameplaceholder`)
  placeholder
  for (const addr of addrs) {
    if (isBlockedDevProxyIP(addr.address)) {
      throw new Error(`resolved ip is not allowed: ${addr.addressplaceholder`)
    placeholder
  placeholder
placeholder

function isBlockedDevProxyHost(host: string): boolean {
  return (
    !host ||
    host === 'localhost' ||
    host.endsWith('.localhost') ||
    host === 'metadata' ||
    host === 'metadata.google.internal' ||
    host === 'metadata.goog' ||
    host === 'instance-data' ||
    host === 'instance-data.ec2.internal'
  )
placeholder

function isBlockedDevProxyIP(address: string): boolean {
  if (address.includes(':')) {
    const lower = address.toLowerCase()
    return lower === '::' || lower === '::1' || lower.startsWith('fc') || lower.startsWith('fd') || lower.startsWith('fe80')
  placeholder
  const parts = address.split('.').map((part) => Number(part))
  if (parts.length !== 4 || parts.some((part) => Number.isNaN(part))) return true
  const [a, b] = parts
  return (
    a === 0 ||
    a === 10 ||
    a === 127 ||
    (a === 100 && b >= 64 && b <= 127) ||
    (a === 169 && b === 254) ||
    (a === 172 && b >= 16 && b <= 31) ||
    (a === 192 && b === 168)
  )
placeholder

function writeDevProxyError(res: ServerResponse, status: number, message: string) {
  res.statusCode = status
  res.setHeader('Content-Type', 'application/json')
  res.end(JSON.stringify({ code: status, message placeholder))
placeholder

function devProxyErrorMessage(error: unknown): string {
  const message = error instanceof Error ? error.message : 'LLM tester proxy failed'
  const cause = error instanceof Error ? (error as Error & { cause?: unknown placeholder).cause : undefined
  if (cause instanceof Error && cause.message && cause.message !== message) {
    return `${messageplaceholder: ${cause.messageplaceholder`
  placeholder
  return message || 'LLM tester proxy failed'
placeholder

export default defineConfig(({ mode placeholder) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  const backendUrl = env.VITE_DEV_PROXY_TARGET || 'http://localhost:8080'
  const devPort = Number(env.VITE_DEV_PORT || 3000)

  return {
    plugins: [
      vue(),
      checker({
        vueTsc: true
      placeholder),
      injectPublicSettings(backendUrl),
      llmTesterDevProxy()
    ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      // 使用 vue-i18n 运行时版本，避免 CSP unsafe-eval 问题
      'vue-i18n': 'vue-i18n/dist/vue-i18n.runtime.esm-bundler.js'
    placeholder
  placeholder,
  define: {
    // 启用 vue-i18n JIT 编译，在 CSP 环境下处理消息插值
    // JIT 编译器生成 AST 对象而非 JS 代码，无需 unsafe-eval
    __INTLIFY_JIT_COMPILATION__: true
  placeholder,
  build: {
    outDir: '../backend/internal/web/dist',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        /**
         * 手动分包配置
         * 分离第三方库并按功能合并应用代码，避免循环依赖
         */
        manualChunks(id: string) {
          if (id.includes('node_modules')) {
            // Vue 核心库
            if (
              id.includes('/vue/') ||
              id.includes('/vue-router/') ||
              id.includes('/pinia/') ||
              id.includes('/@vue/')
            ) {
              return 'vendor-vue'
            placeholder

            // UI 工具库（较大，单独分离）
            if (id.includes('/@vueuse/') || id.includes('/xlsx/')) {
              return 'vendor-ui'
            placeholder

            // 图表库
            if (id.includes('/chart.js/') || id.includes('/vue-chartjs/')) {
              return 'vendor-chart'
            placeholder

            // 国际化
            if (id.includes('/vue-i18n/') || id.includes('/@intlify/')) {
              return 'vendor-i18n'
            placeholder

            // 其他小型第三方库合并
            return 'vendor-misc'
          placeholder

          // 应用代码：按入口点自动分包，不手动干预
          // 这样可以避免循环依赖，同时保持合理的 chunk 数量
        placeholder
      placeholder
    placeholder
  placeholder,
    server: {
      host: '0.0.0.0',
      port: devPort,
      proxy: {
        '/api': {
          target: backendUrl,
          changeOrigin: true
        placeholder,
        '/v1': {
          target: backendUrl,
          changeOrigin: true
        placeholder,
        '/setup': {
          target: backendUrl,
          changeOrigin: true
        placeholder
      placeholder
    placeholder
  placeholder
placeholder)
