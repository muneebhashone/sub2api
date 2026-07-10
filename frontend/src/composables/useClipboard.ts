import { ref placeholder from 'vue'
import { useAppStore placeholder from '@/stores/app'
import { i18n placeholder from '@/i18n'

const { t placeholder = i18n.global

/**
 * 检测是否支持 Clipboard API（需要安全上下文：HTTPS/localhost）
 */
function isClipboardSupported(): boolean {
  return !!(navigator.clipboard && window.isSecureContext)
placeholder

/**
 * 降级方案：使用 textarea + execCommand
 * 使用 textarea 而非 input，以正确处理多行文本
 */
function fallbackCopy(text: string): boolean {
  const textarea = document.createElement('textarea')
  textarea.value = text
  textarea.setAttribute('readonly', 'true')
  textarea.style.cssText = 'position:fixed;left:0;top:0;width:1px;height:1px;opacity:0;pointer-events:none'
  document.body.appendChild(textarea)
  textarea.focus({ preventScroll: true placeholder)
  textarea.select()
  textarea.setSelectionRange(0, textarea.value.length)
  try {
    return document.execCommand('copy')
  placeholder finally {
    document.body.removeChild(textarea)
  placeholder
placeholder

export function useClipboard() {
  const appStore = useAppStore()
  const copied = ref(false)

  const copyToClipboard = async (
    text: string,
    successMessage?: string
  ): Promise<boolean> => {
    if (!text) return false

    let success = false

    if (isClipboardSupported()) {
      try {
        await navigator.clipboard.writeText(text)
        success = true
      placeholder catch {
        success = fallbackCopy(text)
      placeholder
    placeholder else {
      success = fallbackCopy(text)
    placeholder

    if (success) {
      copied.value = true
      appStore.showSuccess(successMessage || t('common.copiedToClipboard'))
      setTimeout(() => {
        copied.value = false
      placeholder, 2000)
    placeholder else {
      appStore.showError(t('common.copyFailed'))
    placeholder

    return success
  placeholder

  return { copied, copyToClipboard placeholder
placeholder
