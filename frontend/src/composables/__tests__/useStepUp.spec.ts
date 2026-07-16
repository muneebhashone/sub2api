import { describe, it, expect, vi placeholder from 'vitest'
import { useStepUp, isStepUpRequired, isStepUpBlocked, stepUpBlockReason placeholder from '../useStepUp'

describe('useStepUp error classification', () => {
  it('detects STEP_UP_REQUIRED from code field', () => {
    expect(isStepUpRequired({ status: 403, code: 'STEP_UP_REQUIRED' placeholder)).toBe(true)
    expect(isStepUpRequired({ status: 403, reason: 'STEP_UP_REQUIRED' placeholder)).toBe(true)
    expect(isStepUpRequired({ status: 500, code: 'INTERNAL' placeholder)).toBe(false)
    expect(isStepUpRequired(null)).toBe(false)
  placeholder)

  it('detects blocked (non-retryable) step-up errors', () => {
    expect(isStepUpBlocked({ code: 'STEP_UP_TOTP_NOT_ENABLED' placeholder)).toBe(true)
    expect(isStepUpBlocked({ reason: 'STEP_UP_ADMIN_API_KEY_FORBIDDEN' placeholder)).toBe(true)
    expect(isStepUpBlocked({ code: 'STEP_UP_REQUIRED' placeholder)).toBe(false)
  placeholder)

  it('surfaces the block reason marker', () => {
    expect(stepUpBlockReason({ reason: 'STEP_UP_ADMIN_API_KEY_FORBIDDEN' placeholder)).toBe('STEP_UP_ADMIN_API_KEY_FORBIDDEN')
    expect(stepUpBlockReason({ code: 'OTHER' placeholder)).toBe('')
  placeholder)
placeholder)

describe('useStepUp.run', () => {
  it('returns the action result directly on success', async () => {
    const stepUp = useStepUp()
    const result = await stepUp.run(async () => 42)
    expect(result).toBe(42)
    expect(stepUp.visible.value).toBe(false)
  placeholder)

  it('re-throws non-step-up errors without prompting', async () => {
    const stepUp = useStepUp()
    const err = { status: 500, code: 'INTERNAL' placeholder
    await expect(stepUp.run(async () => { throw err placeholder)).rejects.toBe(err)
    expect(stepUp.visible.value).toBe(false)
  placeholder)

  it('re-throws blocked errors without prompting', async () => {
    const stepUp = useStepUp()
    const err = { status: 403, code: 'STEP_UP_TOTP_NOT_ENABLED' placeholder
    await expect(stepUp.run(async () => { throw err placeholder)).rejects.toBe(err)
    expect(stepUp.visible.value).toBe(false)
  placeholder)

  it('prompts on STEP_UP_REQUIRED and retries after verification', async () => {
    const stepUp = useStepUp()
    let calls = 0
    const action = async () => {
      calls++
      if (calls === 1) throw { status: 403, code: 'STEP_UP_REQUIRED' placeholder
      return 'ok'
    placeholder
    const promise = stepUp.run(action)
    // The dialog should now be open, awaiting verification (after the first rejection is handled).
    await vi.waitFor(() => expect(stepUp.visible.value).toBe(true))
    stepUp.onVerified()
    await expect(promise).resolves.toBe('ok')
    expect(calls).toBe(2)
    expect(stepUp.visible.value).toBe(false)
  placeholder)

  it('re-throws the original error if the user cancels the prompt', async () => {
    const stepUp = useStepUp()
    const err = { status: 403, code: 'STEP_UP_REQUIRED' placeholder
    const promise = stepUp.run(async () => { throw err placeholder)
    await vi.waitFor(() => expect(stepUp.visible.value).toBe(true))
    stepUp.onCancel()
    await expect(promise).rejects.toBe(err)
  placeholder)
placeholder)
