import { reactive, ref placeholder from 'vue'
import { adminAPI placeholder from '@/api/admin'
import { QUOTA_THRESHOLD_TYPE_FIXED placeholder from '@/constants/account'

export const QUOTA_NOTIFY_DIMS = ['daily', 'weekly', 'total'] as const
export type QuotaNotifyDim = (typeof QUOTA_NOTIFY_DIMS)[number]

interface DimState {
  enabled: boolean | null
  threshold: number | null
  thresholdType: string | null
placeholder

export function useQuotaNotifyState() {
  const globalEnabled = ref(false)
  const state = reactive<Record<QuotaNotifyDim, DimState>>({
    daily: { enabled: null, threshold: null, thresholdType: null placeholder,
    weekly: { enabled: null, threshold: null, thresholdType: null placeholder,
    total: { enabled: null, threshold: null, thresholdType: null placeholder,
  placeholder)

  function loadGlobalState() {
    adminAPI.settings
      .getSettings()
      .then((settings) => {
        globalEnabled.value = settings.account_quota_notify_enabled === true
      placeholder)
      .catch(() => {
        globalEnabled.value = false
      placeholder)
  placeholder

  function loadFromExtra(extra: Record<string, unknown> | null | undefined) {
    for (const d of QUOTA_NOTIFY_DIMS) {
      state[d].enabled = (extra?.[`quota_notify_${dplaceholder_enabled`] as boolean) ?? null
      state[d].threshold = (extra?.[`quota_notify_${dplaceholder_threshold`] as number) ?? null
      state[d].thresholdType = (extra?.[`quota_notify_${dplaceholder_threshold_type`] as string) ?? null
    placeholder
  placeholder

  function writeToExtra(extra: Record<string, unknown>, mode: 'create' | 'update') {
    for (const d of QUOTA_NOTIFY_DIMS) {
      const s = state[d]
      if (s.enabled) {
        extra[`quota_notify_${dplaceholder_enabled`] = true
        if (s.threshold != null) {
          extra[`quota_notify_${dplaceholder_threshold`] = s.threshold
        placeholder else if (mode === 'update') {
          delete extra[`quota_notify_${dplaceholder_threshold`]
        placeholder
        extra[`quota_notify_${dplaceholder_threshold_type`] = s.thresholdType || QUOTA_THRESHOLD_TYPE_FIXED
      placeholder else if (mode === 'update') {
        delete extra[`quota_notify_${dplaceholder_enabled`]
        delete extra[`quota_notify_${dplaceholder_threshold`]
        delete extra[`quota_notify_${dplaceholder_threshold_type`]
      placeholder
    placeholder
  placeholder

  function reset() {
    for (const d of QUOTA_NOTIFY_DIMS) {
      state[d].enabled = null
      state[d].threshold = null
      state[d].thresholdType = null
    placeholder
  placeholder

  return { globalEnabled, state, loadGlobalState, loadFromExtra, writeToExtra, reset placeholder
placeholder
