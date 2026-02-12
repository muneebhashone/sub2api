import { ref placeholder from 'vue'

export interface MixedChannelWarningDetails {
  groupName: string
  currentPlatform: string
  otherPlatform: string
placeholder

function isMixedChannelWarningError(error: any): boolean {
  return error?.response?.status === 409 && error?.response?.data?.error === 'mixed_channel_warning'
placeholder

function extractMixedChannelWarningDetails(error: any): MixedChannelWarningDetails {
  const details = error?.response?.data?.details || {placeholder
  return {
    groupName: details.group_name || 'Unknown',
    currentPlatform: details.current_platform || 'Unknown',
    otherPlatform: details.other_platform || 'Unknown'
  placeholder
placeholder

export function useMixedChannelWarning() {
  const show = ref(false)
  const details = ref<MixedChannelWarningDetails | null>(null)

  const pendingPayload = ref<any | null>(null)
  const pendingRequest = ref<((payload: any) => Promise<any>) | null>(null)
  const pendingOnSuccess = ref<(() => void) | null>(null)
  const pendingOnError = ref<((error: any) => void) | null>(null)

  const clearPending = () => {
    pendingPayload.value = null
    pendingRequest.value = null
    pendingOnSuccess.value = null
    pendingOnError.value = null
    details.value = null
  placeholder

  const tryRequest = async (
    payload: any,
    request: (payload: any) => Promise<any>,
    opts?: {
      onSuccess?: () => void
      onError?: (error: any) => void
    placeholder
  ): Promise<boolean> => {
    try {
      await request(payload)
      opts?.onSuccess?.()
      return true
    placeholder catch (error: any) {
      if (isMixedChannelWarningError(error)) {
        details.value = extractMixedChannelWarningDetails(error)
        pendingPayload.value = payload
        pendingRequest.value = request
        pendingOnSuccess.value = opts?.onSuccess || null
        pendingOnError.value = opts?.onError || null
        show.value = true
        return false
      placeholder

      if (opts?.onError) {
        opts.onError(error)
        return false
      placeholder
      throw error
    placeholder
  placeholder

  const confirm = async (): Promise<boolean> => {
    show.value = false
    if (!pendingPayload.value || !pendingRequest.value) {
      clearPending()
      return false
    placeholder

    pendingPayload.value.confirm_mixed_channel_risk = true

    try {
      await pendingRequest.value(pendingPayload.value)
      pendingOnSuccess.value?.()
      return true
    placeholder catch (error: any) {
      if (pendingOnError.value) {
        pendingOnError.value(error)
        return false
      placeholder
      throw error
    placeholder finally {
      clearPending()
    placeholder
  placeholder

  const cancel = () => {
    show.value = false
    clearPending()
  placeholder

  return {
    show,
    details,
    tryRequest,
    confirm,
    cancel
  placeholder
placeholder

