<template>
  <AppLayout>
    <MonitorHero
      :overall-status="overallStatus"
      :updated-at="updatedAt"
      :interval-seconds="DEFAULT_INTERVAL_SECONDS"
      :window="currentWindow"
      :loading="loading"
      @update:window="handleWindowChange"
      @refresh="manualReload"
    />

    <MonitorCardGrid
      :items="items"
      :window="currentWindow"
      :countdown-seconds="countdown"
      :loading="loading"
      :detail-cache="detailCache"
      @card-click="openDetail"
    />

    <MonitorDetailDialog
      :show="showDetail"
      :monitor-id="detailTarget?.id ?? null"
      :title="detailTitle"
      @close="closeDetail"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch placeholder from 'vue'
import { useI18n placeholder from 'vue-i18n'
import { useAppStore placeholder from '@/stores/app'
import { extractApiErrorMessage placeholder from '@/utils/apiError'
import {
  list as listChannelMonitorViews,
  status as fetchChannelMonitorDetail,
  type UserMonitorView,
  type UserMonitorDetail,
placeholder from '@/api/channelMonitor'
import AppLayout from '@/components/layout/AppLayout.vue'
import MonitorHero, {
  type MonitorWindow,
  type OverallStatus,
placeholder from '@/components/user/monitor/MonitorHero.vue'
import MonitorCardGrid from '@/components/user/monitor/MonitorCardGrid.vue'
import MonitorDetailDialog from '@/components/user/MonitorDetailDialog.vue'
import { DEFAULT_INTERVAL_SECONDS, STATUS_OPERATIONAL placeholder from '@/constants/channelMonitor'

const { t placeholder = useI18n()
const appStore = useAppStore()

// ── State ──
const items = ref<UserMonitorView[]>([])
const loading = ref(false)
const updatedAt = ref<string | null>(null)
const currentWindow = ref<MonitorWindow>('7d')
const detailCache = reactive<Record<number, UserMonitorDetail>>({placeholder)
const countdown = ref(DEFAULT_INTERVAL_SECONDS)

const showDetail = ref(false)
const detailTarget = ref<UserMonitorView | null>(null)

let countdownTimer: number | undefined
let abortController: AbortController | null = null

// ── Computed ──
const overallStatus = computed<OverallStatus>(() => {
  if (items.value.length === 0) return 'operational'
  let hasFailure = false
  let hasDegraded = false
  for (const it of items.value) {
    if (it.primary_status === 'failed' || it.primary_status === 'error') hasFailure = true
    else if (it.primary_status !== STATUS_OPERATIONAL) hasDegraded = true
  placeholder
  if (hasFailure) return 'unavailable'
  if (hasDegraded) return 'degraded'
  return 'operational'
placeholder)

const detailTitle = computed(() => {
  return detailTarget.value?.name || t('channelStatus.detailTitle')
placeholder)

// ── Loaders ──
async function reload(silent = false) {
  if (abortController) abortController.abort()
  const ctrl = new AbortController()
  abortController = ctrl
  if (!silent) loading.value = true
  try {
    const res = await listChannelMonitorViews({ signal: ctrl.signal placeholder)
    if (ctrl.signal.aborted || abortController !== ctrl) return
    items.value = res.items || []
    updatedAt.value = new Date().toISOString()
  placeholder catch (err: unknown) {
    const e = err as { name?: string; code?: string placeholder
    if (e?.name === 'AbortError' || e?.code === 'ERR_CANCELED') return
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.loadError')))
  placeholder finally {
    if (abortController === ctrl) {
      if (!silent) loading.value = false
      countdown.value = DEFAULT_INTERVAL_SECONDS
      abortController = null
    placeholder
  placeholder
placeholder

async function manualReload() {
  await reload(false)
  // After base reload, refresh any cached detail records so non-7d availability
  // values stay in sync without forcing the user to switch tabs again.
  if (currentWindow.value !== '7d') {
    await Promise.all(items.value.map(it => loadDetail(it.id, true)))
  placeholder
placeholder

async function loadDetail(id: number, force = false) {
  if (!force && detailCache[id]) return
  try {
    detailCache[id] = await fetchChannelMonitorDetail(id)
  placeholder catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.detailLoadError')))
  placeholder
placeholder

async function ensureDetailsForWindow() {
  if (currentWindow.value === '7d') return
  await Promise.all(items.value.map(it => loadDetail(it.id)))
placeholder

// ── Handlers ──
async function handleWindowChange(value: MonitorWindow) {
  currentWindow.value = value
  await ensureDetailsForWindow()
placeholder

function openDetail(row: UserMonitorView) {
  detailTarget.value = row
  showDetail.value = true
placeholder

function closeDetail() {
  showDetail.value = false
  detailTarget.value = null
placeholder

// ── Polling ──
function tick() {
  if (countdown.value <= 1) {
    void reload(true)
    return
  placeholder
  countdown.value -= 1
placeholder

watch(items, () => {
  // Lazily load detail entries when window requires it and the list refreshes.
  void ensureDetailsForWindow()
placeholder)

function startTimer() {
  if (countdownTimer !== undefined) return
  countdownTimer = setInterval(tick, 1000) as unknown as number
placeholder

function stopTimer() {
  if (countdownTimer !== undefined) {
    clearInterval(countdownTimer)
    countdownTimer = undefined
  placeholder
placeholder

watch(
  () => appStore.cachedPublicSettings?.channel_monitor_enabled,
  (enabled) => {
    if (enabled === false) stopTimer()
    else startTimer()
  placeholder,
)

// ── Lifecycle ──
onMounted(() => {
  void reload(false)
  if (appStore.cachedPublicSettings?.channel_monitor_enabled !== false) startTimer()
placeholder)

onBeforeUnmount(() => {
  stopTimer()
  if (abortController) abortController.abort()
placeholder)
</script>
