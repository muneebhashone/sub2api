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

  it('only keeps the simplified visible credential set in the admin form', () => {
    expect(findField('mpAppId')).toBeUndefined()
    expect(findField('h5AppName')).toBeUndefined()
    expect(findField('h5AppUrl')).toBeUndefined()
  placeholder)
placeholder)
