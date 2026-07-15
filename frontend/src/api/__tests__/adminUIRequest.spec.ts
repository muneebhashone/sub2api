import { describe, expect, it placeholder from 'vitest'

import {
  ADMIN_UI_REQUEST_HEADER,
  USER_UI_REQUEST_HEADER,
  isUserTimingAPIPath,
  shouldMarkAdminUIRequest,
  shouldMarkUserUIRequest,
placeholder from '@/api/adminUIRequest'

describe('Admin UI request marker', () => {
  it('uses the stable request header name', () => {
    expect(ADMIN_UI_REQUEST_HEADER).toBe('X-Admin-UI-Request')
  placeholder)

  it.each([
    '/admin',
    '/admin/users',
    '/api/v1/admin',
    '/api/v1/admin/accounts?status=active',
    'https://api.example.test/api/v1/admin/dashboard',
  ])('marks Admin API request %s before page navigation', (requestURL) => {
    expect(shouldMarkAdminUIRequest(requestURL, '/login')).toBe(true)
  placeholder)

  it.each(['/keys', '/groups/available', '/auth/me', '/announcements'])(
    'marks shared request %s while an Admin page is active',
    (requestURL) => {
      expect(shouldMarkAdminUIRequest(requestURL, '/admin/dashboard')).toBe(true)
    placeholder
  )

  it.each([
    ['/keys', '/dashboard'],
    ['/api/v1/administer', '/dashboard'],
    ['/keys', '/administrator'],
    ['', '/'],
  ])('does not mark request %s on page %s', (requestURL, pagePath) => {
    expect(shouldMarkAdminUIRequest(requestURL, pagePath)).toBe(false)
  placeholder)
placeholder)

describe('User UI request marker', () => {
  it('uses the stable request header name', () => {
    expect(USER_UI_REQUEST_HEADER).toBe('X-User-UI-Request')
  placeholder)

  it.each([
    '/auth/me',
    '/auth/revoke-all-sessions',
    '/auth/oauth/bind-token',
    '/user',
    '/user/profile',
    '/user/password',
    '/user/notify-email/send-code',
    '/user/totp/status',
    '/user/aff',
    '/user/platform-quotas',
    '/keys',
    '/keys/12',
    '/groups/available',
    '/groups/rates',
    '/channels/available',
    '/usage',
    '/usage/stats',
    '/usage/dashboard/snapshot-v2',
    '/announcements',
    '/announcements/3/read',
    '/redeem',
    '/redeem/history',
    '/subscriptions',
    '/subscriptions/active',
    '/channel-monitors',
    '/channel-monitors/9/status',
    '/payment/config',
    '/payment/plans',
    '/payment/orders',
    '/payment/orders/my',
    '/api/v1/auth/me',
    '/api/v1/keys?page=1',
    'https://api.example.test/api/v1/payment/orders/1',
  ])('marks user timing API %s', (requestURL) => {
    expect(shouldMarkUserUIRequest(requestURL)).toBe(true)
    expect(isUserTimingAPIPath(requestURL)).toBe(true)
  placeholder)

  it.each([
    '/auth/login',
    '/settings/public',
    '/admin/users',
    '/groups',
    '/channels',
    '/payment/public/orders/verify',
    '/payment/webhook/stripe',
    '/api/v1/payment/public/orders/resolve',
    '',
  ])('does not mark non-user timing API %s', (requestURL) => {
    expect(shouldMarkUserUIRequest(requestURL)).toBe(false)
  placeholder)
placeholder)
