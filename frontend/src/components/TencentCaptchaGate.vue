<template>
  <span v-if="false" />
</template>

<script setup lang="ts">
import { onBeforeUnmount placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import {
  loadTencentCaptcha,
  type TencentCaptchaProof,
  type TencentCaptchaResult
placeholder from '@/utils/tencentCaptcha'

const { locale placeholder = useI18n()
const props = defineProps<{ appId: string placeholder>()

let instance: { show(): void; destroy(): void placeholder | null = null
let pending: Promise<TencentCaptchaProof | null> | null = null
let cancelPending: (() => void) | null = null

function createVerificationPromise(): Promise<TencentCaptchaProof | null> {
  return new Promise((resolve, reject) => {
    let settled = false

    const finish = (callback: () => void): void => {
      if (settled) return
      settled = true
      if (cancelPending === cancel) cancelPending = null
      instance?.destroy()
      instance = null
      callback()
    placeholder
    const cancel = (): void => finish(() => resolve(null))

    cancelPending = cancel
    void loadTencentCaptcha()
      .then((TencentCaptcha) => {
        if (cancelPending !== cancel) return

        const userLanguage = locale.value.toLowerCase().startsWith('zh') ? 'zh-cn' : 'en'
        instance = new TencentCaptcha(props.appId, (result: TencentCaptchaResult) => {
          if (result.ret === 2) {
            finish(() => resolve(null))
            return
          placeholder

          const ticket = result.ticket?.trim() || ''
          const randstr = result.randstr?.trim() || ''
          if (!ticket || !randstr || ticket.startsWith('trerror_') || result.errorCode !== undefined) {
            finish(() => reject(new Error('Tencent Captcha verification failed')))
            return
          placeholder

          finish(() => resolve({ ticket, randstr placeholder))
        placeholder, { userLanguage placeholder)
        instance.show()
      placeholder)
      .catch((error: unknown) => finish(() => reject(error)))
  placeholder)
placeholder

function verify(): Promise<TencentCaptchaProof | null> {
  if (pending) return pending
  pending = createVerificationPromise().finally(() => {
    pending = null
  placeholder)
  return pending
placeholder

function reset(): void {
  instance?.destroy()
  instance = null
  cancelPending?.()
  cancelPending = null
placeholder

onBeforeUnmount(reset)
defineExpose({ verify, reset placeholder)
</script>
