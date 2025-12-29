/**
 * Onboarding Store
 * Manages onboarding tour state and control methods
 */

import { defineStore placeholder from 'pinia'
import { markRaw, ref, shallowRef placeholder from 'vue'
import type { Driver placeholder from 'driver.js'

type VoidCallback = () => void
type NextStepCallback = (delay?: number) => Promise<void>
type IsCurrentStepCallback = (selector: string) => boolean

export const useOnboardingStore = defineStore('onboarding', () => {
  const replayCallback = ref<VoidCallback | null>(null)
  const nextStepCallback = ref<NextStepCallback | null>(null)
  const isCurrentStepCallback = ref<IsCurrentStepCallback | null>(null)

  // 全局 driver 实例，跨组件保持
  const driverInstance = shallowRef<Driver | null>(null)

  function setReplayCallback(callback: VoidCallback | null): void {
    replayCallback.value = callback
  placeholder

  function setControlMethods(methods: {
    nextStep: NextStepCallback,
    isCurrentStep: IsCurrentStepCallback
  placeholder): void {
    nextStepCallback.value = methods.nextStep
    isCurrentStepCallback.value = methods.isCurrentStep
  placeholder

  function clearControlMethods(): void {
    nextStepCallback.value = null
    isCurrentStepCallback.value = null
  placeholder

  function setDriverInstance(driver: Driver | null): void {
    driverInstance.value = driver ? markRaw(driver) : null
  placeholder

  function getDriverInstance(): Driver | null {
    return driverInstance.value
  placeholder

  function isDriverActive(): boolean {
    return driverInstance.value?.isActive?.() ?? false
  placeholder

  function replay(): void {
    if (replayCallback.value) {
      replayCallback.value()
    placeholder
  placeholder

  /**
   * Manually advance to the next step
   * @param delay Optional delay in ms (useful for waiting for animations)
   */
  async function nextStep(delay = 0): Promise<void> {
    if (nextStepCallback.value) {
      await nextStepCallback.value(delay)
    placeholder
  placeholder

  /**
   * Check if the tour is currently highlighting a specific element
   */
  function isCurrentStep(selector: string): boolean {
    if (isCurrentStepCallback.value) {
      return isCurrentStepCallback.value(selector)
    placeholder
    return false
  placeholder

  return {
    setReplayCallback,
    setControlMethods,
    clearControlMethods,
    setDriverInstance,
    getDriverInstance,
    isDriverActive,
    replay,
    nextStep,
    isCurrentStep
  placeholder
placeholder)
