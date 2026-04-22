import { describe, expect, it placeholder from 'vitest'
import {
  buildPaymentErrorToastMessage,
  describePaymentScenarioError,
  normalizePaymentMethodForDisplay,
placeholder from '../paymentUx'

describe('normalizePaymentMethodForDisplay', () => {
  it('collapses visible payment aliases to canonical method ids', () => {
    expect(normalizePaymentMethodForDisplay(' alipay_direct ')).toBe('alipay')
    expect(normalizePaymentMethodForDisplay('wxpay_direct')).toBe('wxpay')
    expect(normalizePaymentMethodForDisplay('wechat_pay')).toBe('wxpay')
  placeholder)

  it('leaves non-aliased methods untouched', () => {
    expect(normalizePaymentMethodForDisplay('stripe')).toBe('stripe')
  placeholder)
placeholder)

describe('describePaymentScenarioError', () => {
  it('maps WeChat H5 authorization errors to explicit in-app guidance', () => {
    expect(describePaymentScenarioError(
      { reason: 'WECHAT_H5_NOT_AUTHORIZED' placeholder,
      { paymentMethod: 'wxpay', isMobile: true, isWechatBrowser: false placeholder,
    )).toEqual({
      messageKey: 'payment.errors.wechatH5NotAuthorized',
      hintKey: 'payment.errors.wechatOpenInWeChatHint',
    placeholder)
  placeholder)

  it('maps WeChat H5 authorization errors when provider aliases use wxpay_direct', () => {
    expect(describePaymentScenarioError(
      { reason: 'WECHAT_H5_NOT_AUTHORIZED' placeholder,
      { paymentMethod: 'wxpay_direct', isMobile: true, isWechatBrowser: false placeholder,
    )).toEqual({
      messageKey: 'payment.errors.wechatH5NotAuthorized',
      hintKey: 'payment.errors.wechatOpenInWeChatHint',
    placeholder)
  placeholder)

  it('maps missing WeixinJSBridge to a JSAPI-specific prompt', () => {
    expect(describePaymentScenarioError(
      new Error('WeixinJSBridge is unavailable'),
      { paymentMethod: 'wxpay', isMobile: true, isWechatBrowser: true placeholder,
    )).toEqual({
      messageKey: 'payment.errors.wechatJsapiUnavailable',
      hintKey: 'payment.errors.wechatOpenInWeChatHint',
    placeholder)
  placeholder)

  it('maps the internal JSAPI unavailable marker to the same prompt', () => {
    expect(describePaymentScenarioError(
      new Error('WECHAT_JSAPI_UNAVAILABLE'),
      { paymentMethod: 'wxpay', isMobile: true, isWechatBrowser: true placeholder,
    )).toEqual({
      messageKey: 'payment.errors.wechatJsapiUnavailable',
      hintKey: 'payment.errors.wechatOpenInWeChatHint',
    placeholder)
  placeholder)

  it('maps generic desktop Alipay failures to QR guidance', () => {
    expect(describePaymentScenarioError(
      { reason: 'PAYMENT_GATEWAY_ERROR' placeholder,
      { paymentMethod: 'alipay', isMobile: false, isWechatBrowser: false placeholder,
    )).toEqual({
      messageKey: 'payment.errors.alipayDesktopUnavailable',
      hintKey: 'payment.errors.alipayDesktopQrHint',
    placeholder)
  placeholder)
placeholder)

describe('buildPaymentErrorToastMessage', () => {
  it('returns the main message when no hint is present', () => {
    expect(buildPaymentErrorToastMessage('Payment failed')).toBe('Payment failed')
  placeholder)

  it('appends the hint to the toast body when present', () => {
    expect(buildPaymentErrorToastMessage('Payment failed', 'Open WeChat to continue.')).toBe(
      'Payment failed Open WeChat to continue.'
    )
  placeholder)
placeholder)
