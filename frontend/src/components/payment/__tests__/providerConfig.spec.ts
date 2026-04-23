import { describe, expect, it placeholder from 'vitest'
import { PROVIDER_CONFIG_FIELDS placeholder from '@/components/payment/providerConfig'

function findField(key: string) {
  const fields = PROVIDER_CONFIG_FIELDS.wxpay || []
  return fields.find(field => field.key === key)
placeholder

describe('PROVIDER_CONFIG_FIELDS.wxpay', () => {
  it('keeps admin form validation aligned with backend-required credentials', () => {
    expect(findField('publicKeyId')?.optional).toBeFalsy()
    expect(findField('certSerial')?.optional).toBeFalsy()
  placeholder)

  it('exposes optional mp and H5 metadata fields for WeChat-specific flows', () => {
    expect(findField('mpAppId')?.optional).toBe(true)
    expect(findField('h5AppName')?.optional).toBe(true)
    expect(findField('h5AppUrl')?.optional).toBe(true)
  placeholder)
placeholder)
