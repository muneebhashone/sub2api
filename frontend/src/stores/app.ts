/**
 * Application State Store
 * Manages global UI state including sidebar, loading indicators, and toast notifications
 */

import { defineStore placeholder from 'pinia';
import { ref, computed placeholder from 'vue';
import type { Toast, ToastType placeholder from '@/types';

export const useAppStore = defineStore('app', () => {
  // ==================== State ====================
  
  const sidebarCollapsed = ref<boolean>(false);
  const loading = ref<boolean>(false);
  const toasts = ref<Toast[]>([]);

  // Auto-incrementing ID for toasts
  let toastIdCounter = 0;

  // ==================== Computed ====================
  
  const hasActiveToasts = computed(() => toasts.value.length > 0);
  
  const loadingCount = ref<number>(0);

  // ==================== Actions ====================

  /**
   * Toggle sidebar collapsed state
   */
  function toggleSidebar(): void {
    sidebarCollapsed.value = !sidebarCollapsed.value;
  placeholder

  /**
   * Set sidebar collapsed state explicitly
   * @param collapsed - Whether sidebar should be collapsed
   */
  function setSidebarCollapsed(collapsed: boolean): void {
    sidebarCollapsed.value = collapsed;
  placeholder

  /**
   * Set global loading state
   * @param isLoading - Whether app is in loading state
   */
  function setLoading(isLoading: boolean): void {
    if (isLoading) {
      loadingCount.value++;
    placeholder else {
      loadingCount.value = Math.max(0, loadingCount.value - 1);
    placeholder
    loading.value = loadingCount.value > 0;
  placeholder

  /**
   * Show a toast notification
   * @param type - Type of toast (success, error, info, warning)
   * @param message - Toast message content
   * @param duration - Auto-dismiss duration in ms (undefined = no auto-dismiss)
   * @returns Toast ID for manual dismissal
   */
  function showToast(
    type: ToastType,
    message: string,
    duration?: number
  ): string {
    const id = `toast-${++toastIdCounterplaceholder`;
    const toast: Toast = {
      id,
      type,
      message,
      duration,
      startTime: duration !== undefined ? Date.now() : undefined,
    placeholder;

    toasts.value.push(toast);

    // Auto-dismiss if duration is specified
    if (duration !== undefined) {
      setTimeout(() => {
        hideToast(id);
      placeholder, duration);
    placeholder

    return id;
  placeholder

  /**
   * Show a success toast
   * @param message - Success message
   * @param duration - Auto-dismiss duration in ms (default: 3000)
   */
  function showSuccess(message: string, duration: number = 3000): string {
    return showToast('success', message, duration);
  placeholder

  /**
   * Show an error toast
   * @param message - Error message
   * @param duration - Auto-dismiss duration in ms (default: 5000)
   */
  function showError(message: string, duration: number = 5000): string {
    return showToast('error', message, duration);
  placeholder

  /**
   * Show an info toast
   * @param message - Info message
   * @param duration - Auto-dismiss duration in ms (default: 3000)
   */
  function showInfo(message: string, duration: number = 3000): string {
    return showToast('info', message, duration);
  placeholder

  /**
   * Show a warning toast
   * @param message - Warning message
   * @param duration - Auto-dismiss duration in ms (default: 4000)
   */
  function showWarning(message: string, duration: number = 4000): string {
    return showToast('warning', message, duration);
  placeholder

  /**
   * Hide a specific toast by ID
   * @param id - Toast ID to hide
   */
  function hideToast(id: string): void {
    const index = toasts.value.findIndex((t) => t.id === id);
    if (index !== -1) {
      toasts.value.splice(index, 1);
    placeholder
  placeholder

  /**
   * Clear all toasts
   */
  function clearAllToasts(): void {
    toasts.value = [];
  placeholder

  /**
   * Execute an async operation with loading state
   * Automatically manages loading indicator
   * @param operation - Async operation to execute
   * @returns Promise resolving to operation result
   */
  async function withLoading<T>(operation: () => Promise<T>): Promise<T> {
    setLoading(true);
    try {
      return await operation();
    placeholder finally {
      setLoading(false);
    placeholder
  placeholder

  /**
   * Execute an async operation with loading and error handling
   * Shows error toast on failure
   * @param operation - Async operation to execute
   * @param errorMessage - Custom error message (optional)
   * @returns Promise resolving to operation result or null on error
   */
  async function withLoadingAndError<T>(
    operation: () => Promise<T>,
    errorMessage?: string
  ): Promise<T | null> {
    setLoading(true);
    try {
      return await operation();
    placeholder catch (error) {
      const message =
        errorMessage ||
        (error as { message?: string placeholder).message ||
        'An error occurred';
      showError(message);
      return null;
    placeholder finally {
      setLoading(false);
    placeholder
  placeholder

  /**
   * Reset app state to defaults
   * Useful for cleanup or testing
   */
  function reset(): void {
    sidebarCollapsed.value = false;
    loading.value = false;
    loadingCount.value = 0;
    toasts.value = [];
  placeholder

  // ==================== Return Store API ====================

  return {
    // State
    sidebarCollapsed,
    loading,
    toasts,
    
    // Computed
    hasActiveToasts,
    
    // Actions
    toggleSidebar,
    setSidebarCollapsed,
    setLoading,
    showToast,
    showSuccess,
    showError,
    showInfo,
    showWarning,
    hideToast,
    clearAllToasts,
    withLoading,
    withLoadingAndError,
    reset,
  placeholder;
placeholder);
