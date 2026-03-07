import { defineStore placeholder from 'pinia'
import { ref, computed placeholder from 'vue'
import { announcementsAPI placeholder from '@/api'
import type { UserAnnouncement placeholder from '@/types'

const THROTTLE_MS = 20 * 60 * 1000 // 20 minutes

export const useAnnouncementStore = defineStore('announcements', () => {
  // State
  const announcements = ref<UserAnnouncement[]>([])
  const loading = ref(false)
  const lastFetchTime = ref(0)
  const popupQueue = ref<UserAnnouncement[]>([])
  const currentPopup = ref<UserAnnouncement | null>(null)

  // Session-scoped dedup set — not reactive, used as plain lookup only
  let shownPopupIds = new Set<number>()

  // Getters
  const unreadCount = computed(() =>
    announcements.value.filter((a) => !a.read_at).length
  )

  // Actions
  async function fetchAnnouncements(force = false) {
    const now = Date.now()
    if (!force && lastFetchTime.value > 0 && now - lastFetchTime.value < THROTTLE_MS) {
      return
    placeholder

    // Set immediately to prevent concurrent duplicate requests
    lastFetchTime.value = now

    try {
      loading.value = true
      const all = await announcementsAPI.list(false)
      announcements.value = all.slice(0, 20)
      enqueueNewPopups()
    placeholder catch (err: any) {
      // Revert throttle timestamp on failure so retry is allowed
      lastFetchTime.value = 0
      console.error('Failed to fetch announcements:', err)
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  function enqueueNewPopups() {
    const newPopups = announcements.value.filter(
      (a) => a.notify_mode === 'popup' && !a.read_at && !shownPopupIds.has(a.id)
    )
    if (newPopups.length === 0) return

    for (const p of newPopups) {
      if (!popupQueue.value.some((q) => q.id === p.id)) {
        popupQueue.value.push(p)
      placeholder
    placeholder

    if (!currentPopup.value) {
      showNextPopup()
    placeholder
  placeholder

  function showNextPopup() {
    if (popupQueue.value.length === 0) {
      currentPopup.value = null
      return
    placeholder
    currentPopup.value = popupQueue.value.shift()!
    shownPopupIds.add(currentPopup.value.id)
  placeholder

  async function dismissPopup() {
    if (!currentPopup.value) return
    const id = currentPopup.value.id
    currentPopup.value = null

    // Mark as read (fire-and-forget, UI already updated)
    markAsRead(id)

    // Show next popup after a short delay
    if (popupQueue.value.length > 0) {
      setTimeout(() => showNextPopup(), 300)
    placeholder
  placeholder

  async function markAsRead(id: number) {
    try {
      await announcementsAPI.markRead(id)
      const ann = announcements.value.find((a) => a.id === id)
      if (ann) {
        ann.read_at = new Date().toISOString()
      placeholder
    placeholder catch (err: any) {
      console.error('Failed to mark announcement as read:', err)
    placeholder
  placeholder

  async function markAllAsRead() {
    const unread = announcements.value.filter((a) => !a.read_at)
    if (unread.length === 0) return

    try {
      loading.value = true
      await Promise.all(unread.map((a) => announcementsAPI.markRead(a.id)))
      announcements.value.forEach((a) => {
        if (!a.read_at) {
          a.read_at = new Date().toISOString()
        placeholder
      placeholder)
    placeholder catch (err: any) {
      console.error('Failed to mark all as read:', err)
      throw err
    placeholder finally {
      loading.value = false
    placeholder
  placeholder

  function reset() {
    announcements.value = []
    lastFetchTime.value = 0
    shownPopupIds = new Set()
    popupQueue.value = []
    currentPopup.value = null
    loading.value = false
  placeholder

  return {
    // State
    announcements,
    loading,
    currentPopup,
    // Getters
    unreadCount,
    // Actions
    fetchAnnouncements,
    dismissPopup,
    markAsRead,
    markAllAsRead,
    reset,
  placeholder
placeholder)
