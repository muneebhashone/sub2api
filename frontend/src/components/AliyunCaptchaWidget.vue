<template>
  <div v-if="sceneId && prefix" class="aliyun-captcha-wrapper">
    <button
      :id="buttonId"
      type="button"
      class="aliyun-captcha-button"
      :class="state === 'verified' ? 'aliyun-captcha-button--verified' : ''"
      :disabled="state === 'verified'"
    >
      <svg
        v-if="state === 'verified'"
        class="aliyun-captcha-icon"
        viewBox="0 0 20 20"
        fill="currentColor"
        aria-hidden="true"
      >
        <path
          fill-rule="evenodd"
          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
          clip-rule="evenodd"
        />
      </svg>
      <svg
        v-else
        class="aliyun-captcha-icon"
        viewBox="0 0 20 20"
        fill="currentColor"
        aria-hidden="true"
      >
        <path
          fill-rule="evenodd"
          d="M9.661 2.237a.531.531 0 01.678 0 11.947 11.947 0 007.078 2.749.5.5 0 01.479.425c.069.52.104 1.05.104 1.59 0 5.162-3.26 9.563-7.834 11.256a.48.48 0 01-.332 0C5.26 16.564 2 12.163 2 7c0-.538.035-1.069.104-1.589a.5.5 0 01.48-.425 11.947 11.947 0 007.077-2.75z"
          clip-rule="evenodd"
        />
      </svg>
      <span>{{ buttonText placeholderplaceholder</span>
    </button>
    <div :id="elementId"></div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'

interface AliyunCaptchaVerifyResult {
  captchaResult: boolean
  bizResult?: boolean
placeholder

interface AliyunCaptchaInitOptions {
  SceneId: string
  prefix: string
  mode: 'popup' | 'embed'
  element: string
  button: string
  captchaVerifyCallback: (
    captchaVerifyParam: string
  ) => AliyunCaptchaVerifyResult | Promise<AliyunCaptchaVerifyResult>
  onBizResultCallback: (bizResult: boolean) => void
  getInstance: (instance: unknown) => void
  slideStyle?: { width: number; height: number placeholder
  language?: string
placeholder

declare global {
  interface Window {
    initAliyunCaptcha?: (options: AliyunCaptchaInitOptions) => void
    AliyunCaptchaConfig?: { region: string; prefix: string placeholder
  placeholder
placeholder

const props = withDefaults(
  defineProps<{
    sceneId: string
    prefix: string
    region?: 'cn' | 'sgp'
  placeholder>(),
  {
    region: 'cn'
  placeholder
)

const emit = defineEmits<{
  (e: 'verify', param: string): void
  (e: 'expire'): void
  (e: 'error'): void
placeholder>()

const { t, locale placeholder = useI18n()

const uid = Math.random().toString(36).slice(2, 10)
const buttonId = `aliyun-captcha-button-${uidplaceholder`
const elementId = `aliyun-captcha-element-${uidplaceholder`

// idle: 未验证可点击；verifying: 弹窗已拉起（关闭后回到 idle 可重试）；verified: 已通过
const state = ref<'idle' | 'verifying' | 'verified'>('idle')

const buttonText = computed(() => {
  switch (state.value) {
    case 'verified':
      return t('auth.captchaVerified')
    case 'verifying':
      return t('auth.captchaVerifying')
    default:
      return t('auth.captchaClickToVerify')
  placeholder
placeholder)

const SCRIPT_SRC = 'https://o.alicdn.com/captcha-frontend/aliyunCaptcha/AliyunCaptcha.js'
const POPUP_ID = 'aliyunCaptcha-window-popup'
const MASK_ID = 'aliyunCaptcha-mask'
const POPUP_OPEN_TIMEOUT_MS = 8000
const POPUP_WATCH_INTERVAL_MS = 300

// captchaVerifyParam 是一次性参数：verified 后缓存于此，提交失败需 reset 后重新验证
let cachedParam: string | null = null
let pending: { resolve: (value: string | null) => void placeholder | null = null
let popupWatchTimer: number | null = null
let readyPromise: Promise<void> | null = null

const loadScript = (): Promise<void> => {
  return new Promise((resolve, reject) => {
    // 全局配置必须在脚本加载前就位（region/prefix 全站一致，重复赋值无副作用）
    window.AliyunCaptchaConfig = { region: props.region, prefix: props.prefix placeholder

    if (window.initAliyunCaptcha) {
      resolve()
      return
    placeholder

    const existingScript = document.querySelector<HTMLScriptElement>(
      'script[src*="aliyunCaptcha/AliyunCaptcha"]'
    )
    if (existingScript) {
      existingScript.addEventListener('load', () => resolve())
      existingScript.addEventListener('error', () =>
        reject(new Error('Failed to load Aliyun captcha script'))
      )
      return
    placeholder

    const script = document.createElement('script')
    script.src = SCRIPT_SRC
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => reject(new Error('Failed to load Aliyun captcha script'))
    document.head.appendChild(script)
  placeholder)
placeholder

function initCaptcha(): void {
  if (!window.initAliyunCaptcha) {
    throw new Error('Aliyun captcha script not ready')
  placeholder
  window.initAliyunCaptcha({
    SceneId: props.sceneId,
    prefix: props.prefix,
    mode: 'popup',
    element: `#${elementIdplaceholder`,
    button: `#${buttonIdplaceholder`,
    // 这里不发业务请求，只把 captchaVerifyParam 当 token 交给页面，随登录/注册等
    // 业务请求的 turnstile_token 字段提交，由后端在业务接口内调阿里云校验。
    captchaVerifyCallback: (captchaVerifyParam: string) => {
      onCaptchaParam(captchaVerifyParam)
      return { captchaResult: true placeholder
    placeholder,
    onBizResultCallback: () => {placeholder,
    getInstance: () => {placeholder,
    slideStyle: { width: 360, height: 40 placeholder,
    language: locale.value.toLowerCase().startsWith('zh') ? 'cn' : 'en'
  placeholder)
placeholder

function ensureReady(): Promise<void> {
  if (!readyPromise) {
    readyPromise = loadScript().then(() => initCaptcha())
    readyPromise.catch(() => {
      // 失败后允许下次重试（如网络恢复）
      readyPromise = null
    placeholder)
  placeholder
  return readyPromise
placeholder

function onCaptchaParam(param: string): void {
  stopPopupWatch()
  cachedParam = param
  state.value = 'verified'
  emit('verify', param)
  const current = pending
  pending = null
  current?.resolve(param)
placeholder

function settlePending(value: string | null): void {
  const current = pending
  pending = null
  current?.resolve(value)
placeholder

function isPopupVisible(): boolean {
  const popup = document.getElementById(POPUP_ID)
  if (!popup) return false
  return window.getComputedStyle(popup).display !== 'none'
placeholder

function stopPopupWatch(): void {
  if (popupWatchTimer !== null) {
    window.clearInterval(popupWatchTimer)
    popupWatchTimer = null
  placeholder
placeholder

// SDK 没有用户关闭弹窗的回调，靠轮询弹窗可见性兜底：
// 弹窗出现过又消失且未产生 param → 用户主动关闭；迟迟未出现 → 打开失败。
// initAliyunCaptcha 对触发按钮的事件绑定是异步完成的，首次 click 可能落空，
// 因此弹窗出现前每个 tick 重试触发一次。
function startPopupWatch(): void {
  stopPopupWatch()
  const startedAt = Date.now()
  let seen = false
  popupWatchTimer = window.setInterval(() => {
    if (state.value === 'verified') {
      stopPopupWatch()
      return
    placeholder
    if (isPopupVisible()) {
      seen = true
      return
    placeholder
    if (seen || Date.now() - startedAt > POPUP_OPEN_TIMEOUT_MS) {
      stopPopupWatch()
      state.value = 'idle'
      settlePending(null)
      return
    placeholder
    document.getElementById(buttonId)?.click()
  placeholder, POPUP_WATCH_INTERVAL_MS)
placeholder

// 用户点击与程序化触发共用：置 verifying 并启动弹窗监视（幂等，重试 click 不重置计时）
function handleTriggerClick(): void {
  if (state.value === 'verified') {
    return
  placeholder
  state.value = 'verifying'
  if (popupWatchTimer === null) {
    startPopupWatch()
  placeholder
placeholder

// 程序化触发验证（OAuth 启动、passkey、未预验证时的表单提交兜底）：
// 已通过预验证则直接复用缓存的 captchaVerifyParam；否则弹出验证码等待结果。
// 用户关闭/未能弹出 resolve null；脚本加载失败 reject。
async function verify(): Promise<string | null> {
  if (state.value === 'verified' && cachedParam) {
    return cachedParam
  placeholder
  settlePending(null)
  await ensureReady()
  return new Promise<string | null>((resolve) => {
    pending = { resolve placeholder
    document.getElementById(buttonId)?.click()
  placeholder)
placeholder

// 重置为未验证态；captchaVerifyParam 是一次性参数，服务端校验失败后需用户重新验证
function reset(): void {
  stopPopupWatch()
  settlePending(null)
  cachedParam = null
  state.value = 'idle'
placeholder

defineExpose({ verify, reset placeholder)

onMounted(async () => {
  if (!props.sceneId || !props.prefix) {
    return
  placeholder

  document.getElementById(buttonId)?.addEventListener('click', handleTriggerClick)
  try {
    await ensureReady()
  placeholder catch (error) {
    console.error('Failed to initialize Aliyun captcha:', error)
    emit('error')
  placeholder
placeholder)

onUnmounted(() => {
  document.getElementById(buttonId)?.removeEventListener('click', handleTriggerClick)
  stopPopupWatch()
  settlePending(null)
  // SDK 不会自清理弹窗 DOM，残留会导致下次挂载时回调重复触发
  document.getElementById(MASK_ID)?.remove()
  document.getElementById(POPUP_ID)?.remove()
placeholder)
</script>

<style scoped>
.aliyun-captcha-wrapper {
  width: 100%;
placeholder

.aliyun-captcha-button {
  display: flex;
  width: 100%;
  min-height: 44px;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  border-radius: 0.5rem;
  border: 1px solid rgb(209 213 219);
  background-color: rgb(249 250 251);
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgb(55 65 81);
  transition:
    border-color 0.15s ease,
    background-color 0.15s ease,
    color 0.15s ease;
placeholder

.aliyun-captcha-button:hover:not(:disabled) {
  border-color: rgb(156 163 175);
  background-color: rgb(243 244 246);
placeholder

.aliyun-captcha-button--verified {
  border-color: rgb(34 197 94);
  background-color: rgb(240 253 244);
  color: rgb(21 128 61);
placeholder

:root.dark .aliyun-captcha-button,
.dark .aliyun-captcha-button {
  border-color: rgb(55 65 81);
  background-color: rgb(31 41 55);
  color: rgb(209 213 219);
placeholder

:root.dark .aliyun-captcha-button:hover:not(:disabled),
.dark .aliyun-captcha-button:hover:not(:disabled) {
  border-color: rgb(75 85 99);
  background-color: rgb(55 65 81);
placeholder

:root.dark .aliyun-captcha-button--verified,
.dark .aliyun-captcha-button--verified {
  border-color: rgb(34 197 94);
  background-color: rgb(20 83 45 / 0.3);
  color: rgb(134 239 172);
placeholder

.aliyun-captcha-icon {
  height: 1rem;
  width: 1rem;
  flex-shrink: 0;
placeholder
</style>
