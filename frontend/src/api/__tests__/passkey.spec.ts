import { afterEach, beforeEach, describe, expect, it, vi placeholder from 'vitest'

const { get, post, patch, remove, credentialGet placeholder = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  patch: vi.fn(),
  remove: vi.fn(),
  credentialGet: vi.fn()
placeholder))

vi.mock('@/api/client', () => ({
  apiClient: {
    get,
    post,
    patch,
    delete: remove
  placeholder
placeholder))

import { passkeyAPI placeholder from '@/api/passkey'

class FakePublicKeyCredential {
  id = 'credential-id'
  rawId = Uint8Array.from([1, 2, 3]).buffer
  type = 'public-key'
  authenticatorAttachment = 'platform'
  response = {
    authenticatorData: Uint8Array.from([4, 5]).buffer,
    clientDataJSON: Uint8Array.from([6, 7]).buffer,
    signature: Uint8Array.from([8, 9]).buffer,
    userHandle: Uint8Array.from([10, 11]).buffer
  placeholder

  getClientExtensionResults(): AuthenticationExtensionsClientOutputs {
    return {placeholder
  placeholder
placeholder

describe('passkey api', () => {
  beforeEach(() => {
    get.mockReset()
    post.mockReset()
    patch.mockReset()
    remove.mockReset()
    credentialGet.mockReset()

    vi.stubGlobal('PublicKeyCredential', FakePublicKeyCredential)
    Object.defineProperty(window, 'PublicKeyCredential', {
      configurable: true,
      value: FakePublicKeyCredential
    placeholder)
    Object.defineProperty(navigator, 'credentials', {
      configurable: true,
      value: { get: credentialGet placeholder
    placeholder)
  placeholder)

  afterEach(() => {
    vi.unstubAllGlobals()
  placeholder)

  it('converts assertion options and response bytes to WebAuthn JSON', async () => {
    post
      .mockResolvedValueOnce({
        data: {
          session_token: 'one-time-session',
          options: {
            publicKey: {
              challenge: 'AQID',
              rpId: 'sub2api.example.com',
              userVerification: 'required'
            placeholder
          placeholder
        placeholder
      placeholder)
      .mockResolvedValueOnce({
        data: {
          access_token: 'access',
          token_type: 'Bearer',
          user: { id: 1 placeholder
        placeholder
      placeholder)
    credentialGet.mockResolvedValue(new FakePublicKeyCredential())

    await passkeyAPI.login()

    const request = credentialGet.mock.calls[0][0] as CredentialRequestOptions
    expect(Array.from(new Uint8Array(request.publicKey!.challenge))).toEqual([1, 2, 3])
    expect(request.publicKey!.userVerification).toBe('required')

    expect(post).toHaveBeenNthCalledWith(2, '/auth/passkey/login/finish', {
      session_token: 'one-time-session',
      credential: {
        id: 'credential-id',
        rawId: 'AQID',
        type: 'public-key',
        authenticatorAttachment: 'platform',
        clientExtensionResults: {placeholder,
        response: {
          authenticatorData: 'BAU',
          clientDataJSON: 'Bgc',
          signature: 'CAk',
          userHandle: 'Cgs'
        placeholder
      placeholder
    placeholder)
  placeholder)
placeholder)
