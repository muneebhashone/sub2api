const OAUTH_AFFILIATE_CODE_KEY = 'oauth_aff_code'
const AFFILIATE_REFERRAL_CODE_KEY = 'affiliate_referral_code'
const AFFILIATE_REFERRAL_TTL_MS = 30 * 24 * 60 * 60 * 1000

interface StoredAffiliateReferralCode {
  code: string
  expiresAt: number
placeholder

export function normalizeOAuthAffiliateCode(value?: unknown): string {
  const raw = Array.isArray(value) ? value[0] : value
  return typeof raw === 'string' ? raw.trim() : ''
placeholder

export function pickOAuthAffiliateCode(...values: unknown[]): string {
  for (const value of values) {
    const code = normalizeOAuthAffiliateCode(value)
    if (code) {
      return code
    placeholder
  placeholder
  return ''
placeholder

export function storeAffiliateReferralCode(value?: unknown, now = Date.now()): void {
  if (typeof window === 'undefined') {
    return
  placeholder
  const code = normalizeOAuthAffiliateCode(value)
  if (!code) {
    return
  placeholder
  try {
    const payload: StoredAffiliateReferralCode = {
      code,
      expiresAt: now + AFFILIATE_REFERRAL_TTL_MS
    placeholder
    window.localStorage.setItem(AFFILIATE_REFERRAL_CODE_KEY, JSON.stringify(payload))
  placeholder catch {
    // 忽略浏览器存储异常。
  placeholder
placeholder

export function loadAffiliateReferralCode(now = Date.now()): string {
  if (typeof window === 'undefined') {
    return ''
  placeholder
  try {
    const raw = window.localStorage.getItem(AFFILIATE_REFERRAL_CODE_KEY)
    if (!raw) {
      return ''
    placeholder
    const parsed = JSON.parse(raw) as Partial<StoredAffiliateReferralCode>
    const code = normalizeOAuthAffiliateCode(parsed.code)
    const expiresAt = Number(parsed.expiresAt) || 0
    if (!code || expiresAt <= now) {
      clearAffiliateReferralCode()
      return ''
    placeholder
    return code
  placeholder catch {
    clearAffiliateReferralCode()
    return ''
  placeholder
placeholder

export function clearAffiliateReferralCode(): void {
  if (typeof window === 'undefined') {
    return
  placeholder
  try {
    window.localStorage.removeItem(AFFILIATE_REFERRAL_CODE_KEY)
  placeholder catch {
    // 忽略浏览器存储异常。
  placeholder
placeholder

export function resolveAffiliateReferralCode(...values: unknown[]): string {
  const code = pickOAuthAffiliateCode(...values)
  if (code) {
    storeAffiliateReferralCode(code)
    return code
  placeholder
  return loadAffiliateReferralCode()
placeholder

export function storeOAuthAffiliateCode(value?: unknown): void {
  if (typeof window === 'undefined') {
    return
  placeholder
  const code = normalizeOAuthAffiliateCode(value)
  try {
    if (code) {
      window.sessionStorage.setItem(OAUTH_AFFILIATE_CODE_KEY, code)
    placeholder else {
      window.sessionStorage.removeItem(OAUTH_AFFILIATE_CODE_KEY)
    placeholder
  placeholder catch {
    // 忽略浏览器存储异常。
  placeholder
placeholder

export function loadOAuthAffiliateCode(): string {
  if (typeof window === 'undefined') {
    return ''
  placeholder
  try {
    return normalizeOAuthAffiliateCode(window.sessionStorage.getItem(OAUTH_AFFILIATE_CODE_KEY))
  placeholder catch {
    return ''
  placeholder
placeholder

export function clearOAuthAffiliateCode(): void {
  if (typeof window === 'undefined') {
    return
  placeholder
  try {
    window.sessionStorage.removeItem(OAUTH_AFFILIATE_CODE_KEY)
  placeholder catch {
    // 忽略浏览器存储异常。
  placeholder
placeholder

export function clearAllAffiliateReferralCodes(): void {
  clearOAuthAffiliateCode()
  clearAffiliateReferralCode()
placeholder

export function oauthAffiliatePayload(value?: unknown): { aff_code?: string placeholder {
  const code = normalizeOAuthAffiliateCode(value)
  return code ? { aff_code: code placeholder : {placeholder
placeholder
