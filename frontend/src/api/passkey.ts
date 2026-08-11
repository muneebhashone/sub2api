import { apiClient placeholder from './client'
import type { ActionCaptchaRequestProof, AuthResponse placeholder from '@/types'

export interface PasskeyCredentialSummary {
  id: number
  name: string
  created_at: string
  last_used_at?: string
  backup: boolean
placeholder

interface CeremonyOptionsResponse {
  session_token: string
  options: {
    publicKey: Record<string, unknown>
  placeholder
placeholder

function requirePasskeySupport(): void {
  if (!window.PublicKeyCredential || !navigator.credentials) {
    throw new Error('Passkeys are not supported by this browser')
  placeholder
placeholder

function base64URLToBuffer(value: string): ArrayBuffer {
  const normalized = value.replace(/-/g, '+').replace(/_/g, '/')
  const padded = normalized + '='.repeat((4 - (normalized.length % 4)) % 4)
  const binary = atob(padded)
  const bytes = Uint8Array.from(binary, (character) => character.charCodeAt(0))
  return bytes.buffer
placeholder

function bufferToBase64URL(value: ArrayBuffer | null): string | null {
  if (value === null) return null
  const bytes = new Uint8Array(value)
  let binary = ''
  for (const byte of bytes) binary += String.fromCharCode(byte)
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
placeholder

function creationOptionsFromJSON(
  value: Record<string, unknown>
): PublicKeyCredentialCreationOptions {
  const options = { ...value placeholder as Record<string, unknown>
  options.challenge = base64URLToBuffer(String(options.challenge))

  const user = { ...(options.user as Record<string, unknown>) placeholder
  user.id = base64URLToBuffer(String(user.id))
  options.user = user

  if (Array.isArray(options.excludeCredentials)) {
    options.excludeCredentials = options.excludeCredentials.map((descriptor) => ({
      ...(descriptor as Record<string, unknown>),
      id: base64URLToBuffer(String((descriptor as Record<string, unknown>).id))
    placeholder))
  placeholder
  return options as unknown as PublicKeyCredentialCreationOptions
placeholder

function requestOptionsFromJSON(
  value: Record<string, unknown>
): PublicKeyCredentialRequestOptions {
  const options = { ...value placeholder as Record<string, unknown>
  options.challenge = base64URLToBuffer(String(options.challenge))
  if (Array.isArray(options.allowCredentials)) {
    options.allowCredentials = options.allowCredentials.map((descriptor) => ({
      ...(descriptor as Record<string, unknown>),
      id: base64URLToBuffer(String((descriptor as Record<string, unknown>).id))
    placeholder))
  placeholder
  return options as unknown as PublicKeyCredentialRequestOptions
placeholder

function serializeRegistrationCredential(credential: PublicKeyCredential): Record<string, unknown> {
  const response = credential.response as AuthenticatorAttestationResponse
  return {
    id: credential.id,
    rawId: bufferToBase64URL(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      attestationObject: bufferToBase64URL(response.attestationObject),
      clientDataJSON: bufferToBase64URL(response.clientDataJSON),
      transports: typeof response.getTransports === 'function' ? response.getTransports() : []
    placeholder
  placeholder
placeholder

function serializeAssertionCredential(credential: PublicKeyCredential): Record<string, unknown> {
  const response = credential.response as AuthenticatorAssertionResponse
  return {
    id: credential.id,
    rawId: bufferToBase64URL(credential.rawId),
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      authenticatorData: bufferToBase64URL(response.authenticatorData),
      clientDataJSON: bufferToBase64URL(response.clientDataJSON),
      signature: bufferToBase64URL(response.signature),
      userHandle: bufferToBase64URL(response.userHandle)
    placeholder
  placeholder
placeholder

async function login(proof?: ActionCaptchaRequestProof): Promise<AuthResponse> {
  requirePasskeySupport()
  const { data: begin placeholder = proof
    ? await apiClient.post<CeremonyOptionsResponse>('/auth/passkey/login/begin', proof)
    : await apiClient.post<CeremonyOptionsResponse>('/auth/passkey/login/begin')
  const credential = await navigator.credentials.get({
    publicKey: requestOptionsFromJSON(begin.options.publicKey)
  placeholder)
  if (!(credential instanceof PublicKeyCredential)) {
    throw new Error('Passkey sign-in was cancelled')
  placeholder
  const { data placeholder = await apiClient.post<AuthResponse>('/auth/passkey/login/finish', {
    session_token: begin.session_token,
    credential: serializeAssertionCredential(credential)
  placeholder)
  return data
placeholder

async function register(name: string, password: string): Promise<PasskeyCredentialSummary> {
  requirePasskeySupport()
  const { data: begin placeholder = await apiClient.post<CeremonyOptionsResponse>(
    '/user/passkeys/register/begin',
    { password placeholder
  )
  const credential = await navigator.credentials.create({
    publicKey: creationOptionsFromJSON(begin.options.publicKey)
  placeholder)
  if (!(credential instanceof PublicKeyCredential)) {
    throw new Error('Passkey creation was cancelled')
  placeholder
  const { data placeholder = await apiClient.post<PasskeyCredentialSummary>(
    '/user/passkeys/register/finish',
    {
      session_token: begin.session_token,
      name,
      credential: serializeRegistrationCredential(credential)
    placeholder
  )
  return data
placeholder

async function list(): Promise<PasskeyCredentialSummary[]> {
  const { data placeholder = await apiClient.get<PasskeyCredentialSummary[]>('/user/passkeys')
  return data
placeholder

async function rename(id: number, name: string): Promise<void> {
  await apiClient.patch(`/user/passkeys/${idplaceholder`, { name placeholder)
placeholder

async function remove(id: number, password: string): Promise<void> {
  await apiClient.delete(`/user/passkeys/${idplaceholder`, { data: { password placeholder placeholder)
placeholder

export const passkeyAPI = {
  isSupported: () => Boolean(window.PublicKeyCredential && navigator.credentials),
  login,
  register,
  list,
  rename,
  remove
placeholder
