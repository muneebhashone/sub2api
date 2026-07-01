import { buildApiUrl placeholder from '@/api/client'

export interface LLMTesterProfile {
  id: string
  name: string
  provider: 'openrouter' | 'sub2api' | 'custom'
  baseUrl: string
  apiKey: string
  selectedModel: string
  lastFetchedAt?: string
placeholder

export interface LLMTesterModel {
  id: string
  name: string
  ownedBy?: string
  contextLength?: number
  raw?: Record<string, unknown>
placeholder

export type LLMTesterModelCapability = 'chat' | 'vision' | 'image_generation' | 'video_generation'

export interface LLMTesterAttachment {
  id: string
  name: string
  type: string
  size: number
  kind: 'image' | 'text' | 'media' | 'file'
  dataUrl?: string
  text?: string
placeholder

export interface LLMTesterMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
  attachments?: LLMTesterAttachment[]
placeholder

export interface ChatCompletionOptions {
  baseUrl: string
  apiKey: string
  model: string
  messages: LLMTesterMessage[]
  systemInstruction?: string
  temperature?: number
  maxTokens?: number
  signal?: AbortSignal
placeholder

export interface ImageGenerationOptions {
  baseUrl: string
  apiKey: string
  model: string
  messages: LLMTesterMessage[]
  systemInstruction?: string
  signal?: AbortSignal
placeholder

export interface ImageGenerationResult {
  text: string
  attachments: LLMTesterAttachment[]
  raw: unknown
placeholder

export type MediaGenerationResult = ImageGenerationResult

interface OpenAIContentTextPart {
  type: 'text'
  text: string
placeholder

interface OpenAIContentImagePart {
  type: 'image_url'
  image_url: {
    url: string
  placeholder
placeholder

type OpenAIMessageContent = string | Array<OpenAIContentTextPart | OpenAIContentImagePart>

interface OpenAIChatMessage {
  role: 'system' | 'user' | 'assistant'
  content: OpenAIMessageContent
placeholder

export const OPENROUTER_BASE_URL = 'https://openrouter.ai/api/v1'

export function defaultSub2APIBaseUrl(): string {
  return '/v1'
placeholder

export function normalizeBaseUrl(input: string): string {
  const trimmed = input.trim().replace(/\/+$/, '')
  if (!trimmed) return ''
  if (/^https?:\/\//i.test(trimmed) || trimmed.startsWith('/')) return trimmed
  return `https://${trimmedplaceholder`
placeholder

export type LLMTesterProxyPath = 'models' | 'chat/completions' | 'images/generations' | 'videos/generations' | 'responses'

export function buildOpenAIEndpoint(baseUrl: string, path: LLMTesterProxyPath): string {
  const normalized = normalizeBaseUrl(baseUrl)
  if (!normalized) return ''
  const resource = path.replace(/^v\d+\//, '')
  if (/\/v\d+$/i.test(normalized)) return `${normalizedplaceholder/${resourceplaceholder`
  return `${normalizedplaceholder/v1/${resourceplaceholder`
placeholder

function getHeaderSafeSiteTitle(): string {
  if (typeof document === 'undefined') return 'Sub2API LLM Tester'
  return document.title || 'Sub2API LLM Tester'
placeholder

function buildHeaders(apiKey: string): HeadersInit {
  return {
    Authorization: `Bearer ${apiKeyplaceholder`,
    'Content-Type': 'application/json',
    'X-Title': getHeaderSafeSiteTitle(),
  placeholder
placeholder

function buildJsonHeaders(): HeadersInit {
  return {
    'Content-Type': 'application/json',
  placeholder
placeholder

function getObject(value: unknown): Record<string, unknown> | undefined {
  return value && typeof value === 'object' ? value as Record<string, unknown> : undefined
placeholder

function getString(value: unknown): string | undefined {
  return typeof value === 'string' && value.trim() ? value : undefined
placeholder

function getNumber(value: unknown): number | undefined {
  return typeof value === 'number' && Number.isFinite(value) ? value : undefined
placeholder

function getStringArray(value: unknown): string[] {
  if (!Array.isArray(value)) return []
  return value
    .map((item) => typeof item === 'string' ? item.trim().toLowerCase() : '')
    .filter(Boolean)
placeholder

export function isLikelyChatCompletionModelId(modelId: string): boolean {
  const id = modelId.trim().toLowerCase()
  if (!id) return false
  if (/(^|[/:-])(?:text-)?embedding/.test(id) || id.includes('embedding')) return false
  if (/(^|[/:-])(?:gpt-)?image(?:-|$)/.test(id) || id.includes('/image-')) return false
  if (isLikelyImageGenerationModelId(id) || isLikelyVideoGenerationModelId(id)) return false
  if (id.includes('dall-e') || id.includes('whisper') || id.includes('tts')) return false
  if (id.includes('moderation') || id.includes('omni-moderation')) return false
  if (id.includes('transcribe') || id.includes('realtime')) return false
  return true
placeholder

const GROK_IMAGE_MODEL_IDS = new Set([
  'grok-imagine',
  'grok-imagine-image',
  'grok-imagine-image-quality',
  'grok-imagine-edit',
])

const GROK_VIDEO_MODEL_IDS = new Set([
  'grok-imagine-video',
  'grok-imagine-video-1.5',
])

export function isLikelyImageGenerationModelId(modelId: string): boolean {
  const id = modelId.trim().toLowerCase()
  if (!id) return false
  return (
    GROK_IMAGE_MODEL_IDS.has(id) ||
    /(^|[/:-])(?:gpt-)?image(?:-|$)/.test(id) ||
    id.includes('/image-') ||
    id.includes('dall-e') ||
    id.includes('imagen')
  )
placeholder

export function isLikelyVideoGenerationModelId(modelId: string): boolean {
  const id = modelId.trim().toLowerCase()
  if (!id) return false
  return GROK_VIDEO_MODEL_IDS.has(id) || id.includes('video-generation') || /(^|[/:-])video(?:-|$)/.test(id)
placeholder

function splitModalities(value: string): string[] {
  return value
    .split(/[+,]/)
    .map((part) => part.trim().toLowerCase())
    .filter(Boolean)
placeholder

function getModelModalities(model: LLMTesterModel): { input: string[]; output: string[] placeholder {
  const architecture = getObject(model.raw?.architecture)
  const input = new Set(getStringArray(architecture?.input_modalities))
  const output = new Set(getStringArray(architecture?.output_modalities))

  const modality = getString(architecture?.modality)?.toLowerCase()
  if (modality?.includes('->')) {
    const [inputSide, outputSide] = modality.split('->')
    splitModalities(inputSide || '').forEach((item) => input.add(item))
    splitModalities(outputSide || '').forEach((item) => output.add(item))
  placeholder

  return {
    input: Array.from(input),
    output: Array.from(output),
  placeholder
placeholder

function isKnownUnsupportedModelId(modelId: string): boolean {
  const id = modelId.trim().toLowerCase()
  return (
    /(^|[/:-])(?:text-)?embedding/.test(id) ||
    id.includes('embedding') ||
    id.includes('moderation') ||
    id.includes('omni-moderation') ||
    id.includes('whisper') ||
    id.includes('tts') ||
    id.includes('transcribe') ||
    id.includes('realtime')
  )
placeholder

export function getLLMTesterModelCapabilities(model: LLMTesterModel): LLMTesterModelCapability[] {
  const capabilities = new Set<LLMTesterModelCapability>()
  const modalities = getModelModalities(model)
  const hasOutputMetadata = modalities.output.length > 0
  const outputsText = modalities.output.includes('text')
  const outputsImage = modalities.output.includes('image') || isLikelyImageGenerationModelId(model.id)
  const outputsVideo = modalities.output.includes('video') || isLikelyVideoGenerationModelId(model.id)
  const unsupportedByTester = isKnownUnsupportedModelId(model.id)

  if (outputsImage) {
    capabilities.add('image_generation')
  placeholder

  if (outputsVideo) {
    capabilities.add('video_generation')
  placeholder

  if (!unsupportedByTester && !outputsImage && !outputsVideo && (!hasOutputMetadata || outputsText)) {
    capabilities.add('chat')
  placeholder

  if (capabilities.has('chat') && modalities.input.includes('image')) {
    capabilities.add('vision')
  placeholder

  return Array.from(capabilities)
placeholder

export function isChatCompletionModel(model: LLMTesterModel): boolean {
  return getLLMTesterModelCapabilities(model).includes('chat')
placeholder

export function isImageGenerationModel(model: LLMTesterModel): boolean {
  return getLLMTesterModelCapabilities(model).includes('image_generation')
placeholder

export function isVideoGenerationModel(model: LLMTesterModel): boolean {
  return getLLMTesterModelCapabilities(model).includes('video_generation')
placeholder

export function isLLMTesterSupportedModel(model: LLMTesterModel): boolean {
  const capabilities = getLLMTesterModelCapabilities(model)
  return capabilities.includes('chat') || capabilities.includes('image_generation') || capabilities.includes('video_generation')
placeholder

function extractErrorMessage(payload: unknown, fallback: string): string {
  const obj = getObject(payload)
  const errorObj = getObject(obj?.error)
  return (
    getString(errorObj?.message) ||
    getString(obj?.message) ||
    getString(obj?.detail) ||
    fallback
  )
placeholder

async function parseResponsePayload(response: Response): Promise<unknown> {
  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) return response.json()
  const text = await response.text()
  try {
    return JSON.parse(text)
  placeholder catch {
    return text
  placeholder
placeholder

function unwrapApiEnvelope(payload: unknown): unknown {
  const obj = getObject(payload)
  if (!obj || !('code' in obj) || !('data' in obj)) return payload
  return obj.data
placeholder

function shouldUseTesterProxy(baseUrl: string): boolean {
  const normalized = normalizeBaseUrl(baseUrl)
  if (!normalized || normalized.startsWith('/')) return false
  if (typeof window === 'undefined') return true
  try {
    return new URL(normalized).origin !== window.location.origin
  placeholder catch {
    return true
  placeholder
placeholder

async function postTesterProxy(path: LLMTesterProxyPath, body: Record<string, unknown>, signal?: AbortSignal): Promise<unknown> {
  const response = await fetch(buildApiUrl(`/llm-tester/${pathplaceholder`), {
    method: 'POST',
    headers: buildJsonHeaders(),
    body: JSON.stringify(body),
    signal,
  placeholder)
  const payload = await parseResponsePayload(response)
  if (!response.ok) {
    const fallback = path === 'models'
      ? `Failed to fetch models (${response.statusplaceholder)`
      : path === 'videos/generations'
        ? `Video generation failed (${response.statusplaceholder)`
        : path === 'images/generations' || path === 'responses'
          ? `Image generation failed (${response.statusplaceholder)`
          : `Chat request failed (${response.statusplaceholder)`
    throw new Error(extractErrorMessage(payload, fallback))
  placeholder
  return unwrapApiEnvelope(payload)
placeholder

export function parseModelList(payload: unknown): LLMTesterModel[] {
  const obj = getObject(payload)
  const data = Array.isArray(obj?.data) ? obj.data : Array.isArray(payload) ? payload : []

  return data
    .map((item): LLMTesterModel | null => {
      const raw = getObject(item)
      if (!raw) return null

      const id = getString(raw.id) || getString(raw.name)
      if (!id) return null

      const topProvider = getObject(raw.top_provider)
      return {
        id,
        name: getString(raw.name) || id,
        ownedBy: getString(raw.owned_by) || getString(raw.ownedBy),
        contextLength: getNumber(raw.context_length) || getNumber(raw.contextLength) || getNumber(topProvider?.context_length),
        raw,
      placeholder
    placeholder)
    .filter((model): model is LLMTesterModel => model !== null)
    .filter(isLLMTesterSupportedModel)
    .sort((a, b) => a.id.localeCompare(b.id))
placeholder

export async function fetchLLMModels(baseUrl: string, apiKey: string, signal?: AbortSignal): Promise<LLMTesterModel[]> {
  const endpoint = buildOpenAIEndpoint(baseUrl, 'models')
  if (!endpoint) throw new Error('Base URL is required')

  if (shouldUseTesterProxy(baseUrl)) {
    const payload = await postTesterProxy('models', {
      base_url: normalizeBaseUrl(baseUrl),
      api_key: apiKey,
    placeholder, signal)
    return parseModelList(payload)
  placeholder

  const response = await fetch(endpoint, {
    method: 'GET',
    headers: buildHeaders(apiKey),
    signal,
  placeholder)
  const payload = await parseResponsePayload(response)
  if (!response.ok) {
    throw new Error(extractErrorMessage(payload, `Failed to fetch models (${response.statusplaceholder)`))
  placeholder

  return parseModelList(payload)
placeholder

function inferLanguage(filename: string, type: string): string {
  const lower = filename.toLowerCase()
  const ext = lower.includes('.') ? lower.split('.').pop() || '' : ''
  const byExt: Record<string, string> = {
    js: 'javascript',
    jsx: 'jsx',
    ts: 'typescript',
    tsx: 'tsx',
    vue: 'vue',
    py: 'python',
    go: 'go',
    rs: 'rust',
    java: 'java',
    c: 'c',
    cpp: 'cpp',
    cs: 'csharp',
    html: 'html',
    css: 'css',
    json: 'json',
    md: 'markdown',
    sh: 'bash',
    sql: 'sql',
    yml: 'yaml',
    yaml: 'yaml',
    xml: 'xml',
    toml: 'toml',
    csv: 'csv',
  placeholder
  if (byExt[ext]) return byExt[ext]
  if (type.includes('json')) return 'json'
  if (type.includes('markdown')) return 'markdown'
  if (type.includes('html')) return 'html'
  return ''
placeholder

function formatTextAttachment(attachment: LLMTesterAttachment): string {
  const language = inferLanguage(attachment.name, attachment.type)
  return [
    `Attached file: ${attachment.nameplaceholder`,
    `\`\`\`${languageplaceholder`,
    attachment.text || '',
    '```',
  ].join('\n')
placeholder

function buildImageGenerationPrompt(messages: LLMTesterMessage[], systemInstruction = ''): string {
  const latestUserMessage = [...messages].reverse().find((message) => message.role === 'user')
  const attachments = latestUserMessage?.attachments || []
  const textAttachments = attachments.filter((attachment) => attachment.kind === 'text' && attachment.text)
  const mediaAttachments = attachments.filter((attachment) => attachment.kind !== 'text')

  const sections = [
    systemInstruction.trim(),
    latestUserMessage?.content.trim() || '',
    ...textAttachments.map(formatTextAttachment),
    ...mediaAttachments.map((attachment) => `Attached reference file: ${attachment.nameplaceholder (${attachment.type || 'unknown type'placeholder, ${attachment.sizeplaceholder bytes).`),
  ].filter(Boolean)

  return sections.join('\n\n')
placeholder

function buildMediaGenerationPrompt(messages: LLMTesterMessage[], systemInstruction = ''): string {
  return buildImageGenerationPrompt(messages, systemInstruction)
placeholder

function buildUserContent(message: LLMTesterMessage): OpenAIMessageContent {
  const attachments = message.attachments || []
  const imageAttachments = attachments.filter((attachment) => attachment.kind === 'image' && attachment.dataUrl)
  const textAttachments = attachments.filter((attachment) => attachment.kind === 'text' && attachment.text)
  const otherAttachments = attachments.filter((attachment) => attachment.kind !== 'image' && attachment.kind !== 'text')

  const textParts = [
    message.content.trim(),
    ...textAttachments.map(formatTextAttachment),
    ...otherAttachments.map((attachment) => `Attached media: ${attachment.nameplaceholder (${attachment.type || 'unknown type'placeholder, ${attachment.sizeplaceholder bytes).`),
  ].filter(Boolean)

  if (imageAttachments.length === 0) return textParts.join('\n\n')

  const content: Array<OpenAIContentTextPart | OpenAIContentImagePart> = []
  content.push({
    type: 'text',
    text: textParts.join('\n\n') || 'Please analyze the attached image.',
  placeholder)

  for (const attachment of imageAttachments) {
    if (!attachment.dataUrl) continue
    content.push({
      type: 'image_url',
      image_url: { url: attachment.dataUrl placeholder,
    placeholder)
  placeholder

  return content
placeholder

export function buildChatCompletionMessages(messages: LLMTesterMessage[], systemInstruction = ''): OpenAIChatMessage[] {
  const out: OpenAIChatMessage[] = []
  const system = systemInstruction.trim()
  if (system) {
    out.push({ role: 'system', content: system placeholder)
  placeholder

  for (const message of messages) {
    out.push({
      role: message.role,
      content: message.role === 'user' ? buildUserContent(message) : message.content,
    placeholder)
  placeholder

  return out
placeholder

export function extractChatCompletionText(payload: unknown): string {
  const obj = getObject(payload)
  const choices = Array.isArray(obj?.choices) ? obj.choices : []
  const firstChoice = getObject(choices[0])
  const message = getObject(firstChoice?.message)
  const content = message?.content

  if (typeof content === 'string') return content
  if (Array.isArray(content)) {
    return content
      .map((part) => {
        const partObj = getObject(part)
        return getString(partObj?.text) || getString(partObj?.content) || ''
      placeholder)
      .filter(Boolean)
      .join('\n')
  placeholder

  const text = getString(firstChoice?.text)
  if (text) return text

  return JSON.stringify(payload, null, 2)
placeholder

export function extractImageGenerationResult(payload: unknown): ImageGenerationResult {
  const attachments: LLMTesterAttachment[] = []
  const lines: string[] = []

  const pushImageAttachment = (rawValue: unknown, index: number) => {
    const value = normalizeGeneratedImageValue(rawValue)
    if (!value) return
    attachments.push({
      id: `generated-image-${Date.now()placeholder-${indexplaceholder`,
      name: `generated-image-${index + 1placeholder.png`,
      type: 'image/png',
      size: 0,
      kind: 'image',
      dataUrl: value,
    placeholder)
  placeholder

  const explicitImageResult = (value: unknown): unknown => {
    const text = getString(value)
    if (!text) return value
    if (/^(?:data:image\/|https?:\/\/)/i.test(text)) return text
    return `data:image/png;base64,${textplaceholder`
  placeholder

  const processOutputItem = (item: unknown) => {
    const outputItem = getObject(item)
    if (!outputItem) return
    const type = getString(outputItem.type)

    if (type === 'image_generation_call') {
      const b64 = getString(outputItem.b64_json)
      pushImageAttachment(b64 ? `data:image/png;base64,${b64placeholder` : explicitImageResult(outputItem.result) || outputItem.image_url || outputItem.url, attachments.length)
      const revisedPrompt = getString(outputItem.revised_prompt)
      if (revisedPrompt) {
        lines.push(`Revised prompt: ${revisedPromptplaceholder`)
      placeholder
    placeholder

    const content = Array.isArray(outputItem.content) ? outputItem.content : []
    content.forEach((part) => {
      const partObj = getObject(part)
      if (!partObj) return
      const partType = getString(partObj.type)
      const text = getString(partObj.text)
      if (text && (partType === 'output_text' || partType === 'text')) {
        lines.push(text)
      placeholder
      const b64 = getString(partObj.b64_json)
      pushImageAttachment(b64 ? `data:image/png;base64,${b64placeholder` : explicitImageResult(partObj.result) || partObj.image_url || partObj.url, attachments.length)
    placeholder)

    const outputText = getString(outputItem.text)
    if (outputText && type !== 'image_generation_call') {
      lines.push(outputText)
    placeholder
  placeholder

  const processPayload = (rawPayload: unknown) => {
    const obj = getObject(rawPayload)
    if (!obj) return

    if (obj.item) {
      processOutputItem(obj.item)
    placeholder
    if (obj.response) {
      processPayload(obj.response)
    placeholder

    const data = Array.isArray(obj.data) ? obj.data : []
    data.forEach((item, index) => {
      const image = getObject(item)
      if (!image) return

      const revisedPrompt = getString(image.revised_prompt)
      if (revisedPrompt) {
        lines.push(`Revised prompt: ${revisedPromptplaceholder`)
      placeholder

      const b64 = getString(image.b64_json)
      const url = getString(image.url)
      pushImageAttachment(b64 ? `data:image/png;base64,${b64placeholder` : url, index)
    placeholder)

    const output = Array.isArray(obj.output) ? obj.output : []
    output.forEach(processOutputItem)
  placeholder

  const payloads = typeof payload === 'string' ? parseEventStreamPayload(payload) : [payload]
  payloads.forEach(processPayload)

  if (attachments.length > 0) {
    lines.unshift(`Generated ${attachments.lengthplaceholder image${attachments.length === 1 ? '' : 's'placeholder.`)
  placeholder

  return {
    text: lines.join('\n\n') || JSON.stringify(payload, null, 2),
    attachments,
    raw: payload,
  placeholder
placeholder

function parseEventStreamPayload(payload: string): unknown[] {
  const events: unknown[] = []
  const dataLines: string[] = []

  const flush = () => {
    const data = dataLines.join('\n').trim()
    dataLines.length = 0
    if (!data || data === '[DONE]') return
    try {
      events.push(JSON.parse(data))
    placeholder catch {
      events.push(data)
    placeholder
  placeholder

  for (const line of payload.split(/\r?\n/)) {
    if (line.startsWith('data:')) {
      dataLines.push(line.slice(5).trimStart())
      continue
    placeholder
    if (!line.trim()) {
      flush()
    placeholder
  placeholder
  flush()

  if (events.length > 0) return events
  try {
    return [JSON.parse(payload)]
  placeholder catch {
    return []
  placeholder
placeholder

function normalizeGeneratedImageValue(value: unknown): string {
  if (typeof value === 'object' && value !== null) {
    const obj = getObject(value)
    return normalizeGeneratedImageValue(obj?.url || obj?.b64_json || obj?.result)
  placeholder
  const text = getString(value)
  if (!text) return ''
  if (/^data:image\//i.test(text)) return text
  if (/^https?:\/\//i.test(text)) return text
  const compact = text.replace(/\s+/g, '')
  if (compact.length > 100 && /^[A-Za-z0-9+/=]+$/.test(compact)) {
    return `data:image/png;base64,${compactplaceholder`
  placeholder
  return ''
placeholder

function normalizeGeneratedMediaValue(value: unknown): string {
  if (typeof value === 'object' && value !== null) {
    const obj = getObject(value)
    return normalizeGeneratedMediaValue(
      obj?.url ||
      obj?.video_url ||
      obj?.download_url ||
      obj?.b64_json ||
      obj?.base64 ||
      obj?.result
    )
  placeholder
  const text = getString(value)
  if (!text) return ''
  if (/^data:video\//i.test(text)) return text
  if (/^https?:\/\//i.test(text)) return text
  const compact = text.replace(/\s+/g, '')
  if (compact.length > 100 && /^[A-Za-z0-9+/=]+$/.test(compact)) {
    return `data:video/mp4;base64,${compactplaceholder`
  placeholder
  return ''
placeholder

export function extractVideoGenerationResult(payload: unknown): MediaGenerationResult {
  const attachments: LLMTesterAttachment[] = []
  const lines: string[] = []

  const pushVideoAttachment = (rawValue: unknown, index: number) => {
    const value = normalizeGeneratedMediaValue(rawValue)
    if (!value) return
    attachments.push({
      id: `generated-video-${Date.now()placeholder-${indexplaceholder`,
      name: `generated-video-${index + 1placeholder.mp4`,
      type: 'video/mp4',
      size: 0,
      kind: 'media',
      dataUrl: value,
    placeholder)
  placeholder

  const processObject = (value: unknown) => {
    const obj = getObject(value)
    if (!obj) return

    const status = getString(obj.status)
    if (status) lines.push(`Status: ${statusplaceholder`)
    const id = getString(obj.id) || getString(obj.request_id)
    if (id) lines.push(`Request ID: ${idplaceholder`)
    const revisedPrompt = getString(obj.revised_prompt)
    if (revisedPrompt) lines.push(`Revised prompt: ${revisedPromptplaceholder`)

    pushVideoAttachment(obj, attachments.length)

    const data = Array.isArray(obj.data) ? obj.data : []
    data.forEach((item) => {
      processObject(item)
    placeholder)

    const output = Array.isArray(obj.output) ? obj.output : []
    output.forEach((item) => {
      processObject(item)
    placeholder)

    const content = Array.isArray(obj.content) ? obj.content : []
    content.forEach((item) => {
      const itemObj = getObject(item)
      const text = getString(itemObj?.text)
      if (text) lines.push(text)
      processObject(item)
    placeholder)
  placeholder

  const payloads = typeof payload === 'string' ? parseEventStreamPayload(payload) : [payload]
  payloads.forEach(processObject)

  const uniqueLines = Array.from(new Set(lines))
  if (attachments.length > 0) {
    uniqueLines.unshift(`Generated ${attachments.lengthplaceholder video${attachments.length === 1 ? '' : 's'placeholder.`)
  placeholder

  return {
    text: uniqueLines.join('\n\n') || JSON.stringify(payload, null, 2),
    attachments,
    raw: payload,
  placeholder
placeholder

function imageToolModelId(model: string): string {
  const trimmed = model.trim()
  if (!trimmed) return 'gpt-image-2'
  const parts = trimmed.split('/').filter(Boolean)
  return parts[parts.length - 1] || trimmed
placeholder

function imageResponsesDriverModel(model: string): string {
  return isLikelyImageGenerationModelId(model) ? 'gpt-5.4' : model
placeholder

function buildResponsesImageGenerationBody(model: string, prompt: string): Record<string, unknown> {
  return {
    model: imageResponsesDriverModel(model),
    stream: true,
    tools: [
      {
        type: 'image_generation',
        model: imageToolModelId(model),
      placeholder,
    ],
    input: [
      {
        role: 'user',
        content: [
          {
            type: 'input_text',
            text: prompt,
          placeholder,
        ],
      placeholder,
    ],
  placeholder
placeholder

function isAbortError(error: unknown): boolean {
  return error instanceof DOMException && error.name === 'AbortError'
placeholder

async function postOpenAIResource(
  baseUrl: string,
  apiKey: string,
  path: LLMTesterProxyPath,
  body: Record<string, unknown>,
  signal?: AbortSignal
): Promise<unknown> {
  const endpoint = buildOpenAIEndpoint(baseUrl, path)
  if (!endpoint) throw new Error('Base URL is required')

  if (shouldUseTesterProxy(baseUrl)) {
    return postTesterProxy(path, {
      base_url: normalizeBaseUrl(baseUrl),
      api_key: apiKey,
      payload: body,
    placeholder, signal)
  placeholder

  const response = await fetch(endpoint, {
    method: 'POST',
    headers: buildHeaders(apiKey),
    body: JSON.stringify(body),
    signal,
  placeholder)
  const payload = await parseResponsePayload(response)
  if (!response.ok) {
    const fallback = path === 'chat/completions'
      ? `Chat request failed (${response.statusplaceholder)`
      : path === 'videos/generations'
        ? `Video generation failed (${response.statusplaceholder)`
        : `Image generation failed (${response.statusplaceholder)`
    throw new Error(extractErrorMessage(payload, fallback))
  placeholder

  return payload
placeholder

export async function sendLLMChatCompletion(options: ChatCompletionOptions): Promise<{ text: string; raw: unknown placeholder> {
  const endpoint = buildOpenAIEndpoint(options.baseUrl, 'chat/completions')
  if (!endpoint) throw new Error('Base URL is required')

  const body: Record<string, unknown> = {
    model: options.model,
    messages: buildChatCompletionMessages(options.messages, options.systemInstruction),
    stream: false,
  placeholder

  if (typeof options.temperature === 'number' && Number.isFinite(options.temperature)) {
    body.temperature = options.temperature
  placeholder
  if (typeof options.maxTokens === 'number' && Number.isFinite(options.maxTokens) && options.maxTokens > 0) {
    body.max_tokens = Math.floor(options.maxTokens)
  placeholder

  if (shouldUseTesterProxy(options.baseUrl)) {
    const payload = await postTesterProxy('chat/completions', {
      base_url: normalizeBaseUrl(options.baseUrl),
      api_key: options.apiKey,
      payload: body,
    placeholder, options.signal)
    return {
      text: extractChatCompletionText(payload),
      raw: payload,
    placeholder
  placeholder

  const response = await fetch(endpoint, {
    method: 'POST',
    headers: buildHeaders(options.apiKey),
    body: JSON.stringify(body),
    signal: options.signal,
  placeholder)
  const payload = await parseResponsePayload(response)
  if (!response.ok) {
    throw new Error(extractErrorMessage(payload, `Chat request failed (${response.statusplaceholder)`))
  placeholder

  return {
    text: extractChatCompletionText(payload),
    raw: payload,
  placeholder
placeholder

export async function sendLLMImageGeneration(options: ImageGenerationOptions): Promise<ImageGenerationResult> {
  const prompt = buildImageGenerationPrompt(options.messages, options.systemInstruction)
  if (!prompt) throw new Error('Prompt is required for image generation')

  const body: Record<string, unknown> = {
    model: options.model,
    prompt,
    n: 1,
  placeholder
  if (/^gpt-image-/i.test(imageToolModelId(options.model))) {
    body.stream = true
  placeholder

  try {
    const payload = await postOpenAIResource(options.baseUrl, options.apiKey, 'images/generations', body, options.signal)
    return extractImageGenerationResult(payload)
  placeholder catch (primaryError) {
    if (isAbortError(primaryError)) throw primaryError

    try {
      const fallbackPayload = await postOpenAIResource(
        options.baseUrl,
        options.apiKey,
        'responses',
        buildResponsesImageGenerationBody(options.model, prompt),
        options.signal
      )
      const fallbackResult = extractImageGenerationResult(fallbackPayload)
      if (fallbackResult.attachments.length > 0) return fallbackResult
      throw new Error('Responses image tool returned no image output')
    placeholder catch (fallbackError) {
      if (isAbortError(fallbackError)) throw fallbackError
      const primaryMessage = primaryError instanceof Error ? primaryError.message : 'Image endpoint failed'
      const fallbackMessage = fallbackError instanceof Error ? fallbackError.message : 'Responses fallback failed'
      throw new Error(`${primaryMessageplaceholder; responses fallback failed: ${fallbackMessageplaceholder`)
    placeholder
  placeholder
placeholder

export async function sendLLMVideoGeneration(options: ImageGenerationOptions): Promise<MediaGenerationResult> {
  const prompt = buildMediaGenerationPrompt(options.messages, options.systemInstruction)
  if (!prompt) throw new Error('Prompt is required for video generation')

  const body: Record<string, unknown> = {
    model: options.model,
    prompt,
  placeholder

  const payload = await postOpenAIResource(options.baseUrl, options.apiKey, 'videos/generations', body, options.signal)
  return extractVideoGenerationResult(payload)
placeholder
