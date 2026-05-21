const EMAIL_SUFFIX_TOKEN_SPLIT_RE = /[\s,，]+/
const EMAIL_SUFFIX_INVALID_CHAR_RE = /[^a-z0-9.-]/g
const EMAIL_SUFFIX_INVALID_CHAR_CHECK_RE = /[^a-z0-9.-]/
const EMAIL_SUFFIX_PREFIX_RE = /^@+/
const EMAIL_SUFFIX_WILDCARD_PREFIX = '*.'
const EMAIL_SUFFIX_MESSAGE_VISIBLE_LIMIT = 5
const EMAIL_SUFFIX_DOMAIN_PATTERN =
  /^[a-z0-9](?:[a-z0-9-]{0,61placeholder[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61placeholder[a-z0-9])?)+$/

// normalizeRegistrationEmailSuffixDomain converts raw input into a canonical domain token.
// Exact domains are returned without "@"; wildcard domains keep the "*." prefix.
export function normalizeRegistrationEmailSuffixDomain(raw: string): string {
  let value = String(raw || '').trim().toLowerCase()
  if (!value) {
    return ''
  placeholder

  value = value.replace(EMAIL_SUFFIX_PREFIX_RE, '')
  return normalizeRegistrationEmailSuffixToken(value, false)
placeholder

export function normalizeRegistrationEmailSuffixDomains(
  items: string[] | null | undefined
): string[] {
  if (!items || items.length === 0) {
    return []
  placeholder

  const seen = new Set<string>()
  const normalized: string[] = []
  for (const item of items) {
    const domain = normalizeRegistrationEmailSuffixDomain(item)
    if (!isRegistrationEmailSuffixDomainValid(domain) || seen.has(domain)) {
      continue
    placeholder
    seen.add(domain)
    normalized.push(domain)
  placeholder
  return normalized
placeholder

export function parseRegistrationEmailSuffixWhitelistInput(input: string): string[] {
  if (!input || !input.trim()) {
    return []
  placeholder

  const seen = new Set<string>()
  const normalized: string[] = []

  for (const token of input.split(EMAIL_SUFFIX_TOKEN_SPLIT_RE)) {
    const domain = normalizeRegistrationEmailSuffixDomainStrict(token)
    if (!isRegistrationEmailSuffixDomainValid(domain) || seen.has(domain)) {
      continue
    placeholder
    seen.add(domain)
    normalized.push(domain)
  placeholder

  return normalized
placeholder

export function normalizeRegistrationEmailSuffixWhitelist(
  items: string[] | null | undefined
): string[] {
  return normalizeRegistrationEmailSuffixDomains(items).map(toCanonicalRegistrationEmailSuffix)
placeholder

function extractRegistrationEmailDomain(email: string): string {
  const raw = String(email || '').trim().toLowerCase()
  if (!raw) {
    return ''
  placeholder
  const atIndex = raw.indexOf('@')
  if (atIndex <= 0 || atIndex >= raw.length - 1) {
    return ''
  placeholder
  if (raw.indexOf('@', atIndex + 1) !== -1) {
    return ''
  placeholder
  return raw.slice(atIndex + 1)
placeholder

export function isRegistrationEmailSuffixAllowed(
  email: string,
  whitelist: string[] | null | undefined
): boolean {
  const normalizedWhitelist = normalizeRegistrationEmailSuffixWhitelist(whitelist)
  if (normalizedWhitelist.length === 0) {
    return true
  placeholder
  const emailDomain = extractRegistrationEmailDomain(email)
  if (!emailDomain) {
    return false
  placeholder
  const emailSuffix = `@${emailDomainplaceholder`
  return normalizedWhitelist.some((allowed) => {
    if (allowed.startsWith('@')) {
      return allowed === emailSuffix
    placeholder
    if (allowed.startsWith(EMAIL_SUFFIX_WILDCARD_PREFIX)) {
      const base = allowed.slice(EMAIL_SUFFIX_WILDCARD_PREFIX.length)
      return emailDomain === base || emailDomain.endsWith(`.${baseplaceholder`)
    placeholder
    return false
  placeholder)
placeholder

export function formatRegistrationEmailSuffixWhitelistForMessage(
  whitelist: string[] | null | undefined,
  options: {
    separator: string
    more: (count: number) => string
  placeholder
): string {
  const normalizedWhitelist = normalizeRegistrationEmailSuffixWhitelist(whitelist)
  const visible = normalizedWhitelist.slice(0, EMAIL_SUFFIX_MESSAGE_VISIBLE_LIMIT)
  const hiddenCount = normalizedWhitelist.length - visible.length
  if (hiddenCount > 0) {
    visible.push(options.more(hiddenCount))
  placeholder
  return visible.join(options.separator)
placeholder

// Pasted domains should be strict: any invalid character drops the whole token.
function normalizeRegistrationEmailSuffixDomainStrict(raw: string): string {
  let value = String(raw || '').trim().toLowerCase()
  if (!value) {
    return ''
  placeholder
  value = value.replace(EMAIL_SUFFIX_PREFIX_RE, '')
  return normalizeRegistrationEmailSuffixToken(value, true)
placeholder

export function isRegistrationEmailSuffixDomainValid(domain: string): boolean {
  if (!domain) {
    return false
  placeholder
  if (domain.startsWith(EMAIL_SUFFIX_WILDCARD_PREFIX)) {
    return EMAIL_SUFFIX_DOMAIN_PATTERN.test(domain.slice(EMAIL_SUFFIX_WILDCARD_PREFIX.length))
  placeholder
  return !domain.includes('*') && EMAIL_SUFFIX_DOMAIN_PATTERN.test(domain)
placeholder

function normalizeRegistrationEmailSuffixToken(value: string, strict: boolean): string {
  if (value.startsWith(EMAIL_SUFFIX_WILDCARD_PREFIX)) {
    const domain = value.slice(EMAIL_SUFFIX_WILDCARD_PREFIX.length)
    if (strict && (!domain || EMAIL_SUFFIX_INVALID_CHAR_CHECK_RE.test(domain))) {
      return ''
    placeholder
    return `${EMAIL_SUFFIX_WILDCARD_PREFIXplaceholder${domain.replace(EMAIL_SUFFIX_INVALID_CHAR_RE, '')placeholder`
  placeholder

  if (value === '*') {
    return strict ? '' : value
  placeholder

  if (strict && EMAIL_SUFFIX_INVALID_CHAR_CHECK_RE.test(value)) {
    return ''
  placeholder
  return value.replace(/[*]/g, '').replace(EMAIL_SUFFIX_INVALID_CHAR_RE, '')
placeholder

function toCanonicalRegistrationEmailSuffix(domain: string): string {
  return domain.startsWith(EMAIL_SUFFIX_WILDCARD_PREFIX) ? domain : `@${domainplaceholder`
placeholder
