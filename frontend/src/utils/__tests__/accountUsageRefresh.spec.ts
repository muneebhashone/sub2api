import { describe, expect, it placeholder from 'vitest'
import { buildGrokUsageRefreshKey, buildOpenAIUsageRefreshKey placeholder from '../accountUsageRefresh'

describe('buildOpenAIUsageRefreshKey', () => {
  it('会在 codex 快照变化时生成不同 key', () => {
    const base = {
      id: 1,
      platform: 'openai',
      type: 'oauth',
      updated_at: '2026-03-07T10:00:00Z',
      last_used_at: '2026-03-07T09:59:00Z',
      extra: {
        codex_usage_updated_at: '2026-03-07T10:00:00Z',
        codex_5h_used_percent: 0,
        codex_7d_used_percent: 0
      placeholder
    placeholder as any

    const next = {
      ...base,
      extra: {
        ...base.extra,
        codex_usage_updated_at: '2026-03-07T10:01:00Z',
        codex_5h_used_percent: 100
      placeholder
    placeholder

    expect(buildOpenAIUsageRefreshKey(base)).not.toBe(buildOpenAIUsageRefreshKey(next))
  placeholder)

  it('会在 last_used_at 变化时生成不同 key', () => {
    const base = {
      id: 3,
      platform: 'openai',
      type: 'oauth',
      updated_at: '2026-03-07T10:00:00Z',
      last_used_at: '2026-03-07T10:00:00Z',
      extra: {
        codex_usage_updated_at: '2026-03-07T10:00:00Z',
        codex_5h_used_percent: 12,
        codex_7d_used_percent: 24
      placeholder
    placeholder as any

    const next = {
      ...base,
      last_used_at: '2026-03-07T10:02:00Z'
    placeholder

    expect(buildOpenAIUsageRefreshKey(base)).not.toBe(buildOpenAIUsageRefreshKey(next))
  placeholder)

  it('非 OpenAI OAuth 账号返回空 key', () => {
    expect(buildOpenAIUsageRefreshKey({
      id: 2,
      platform: 'anthropic',
      type: 'oauth',
      updated_at: '2026-03-07T10:00:00Z',
      last_used_at: '2026-03-07T10:00:00Z',
      extra: {placeholder
    placeholder as any)).toBe('')
  placeholder)
placeholder)

describe('buildGrokUsageRefreshKey', () => {
  it('changes when a canonical Grok billing or usage snapshot changes', () => {
    const base = {
      platform: 'grok',
      extra: {
        grok_billing_snapshot: { plan: 'Free', usage_percent: 0 placeholder,
        grok_usage_snapshot: { subscription_tier: 'Free', status_code: 200 placeholder
      placeholder
    placeholder as any

    expect(buildGrokUsageRefreshKey(base)).not.toBe(buildGrokUsageRefreshKey({
      ...base,
      extra: {
        ...base.extra,
        grok_billing_snapshot: { plan: 'SuperGrok', usage_percent: 0 placeholder
      placeholder
    placeholder))
    expect(buildGrokUsageRefreshKey(base)).not.toBe(buildGrokUsageRefreshKey({
      ...base,
      extra: {
        ...base.extra,
        grok_usage_snapshot: { subscription_tier: 'SuperGrok', status_code: 200 placeholder
      placeholder
    placeholder))
  placeholder)

  it('ignores object key order and a legacy alias shadowed by canonical usage', () => {
    const first = {
      platform: 'grok',
      extra: {
        grok_billing_snapshot: {
          plan: 'SuperGrok',
          limits: { monthly: 100, weekly: 25 placeholder
        placeholder,
        grok_usage_snapshot: { status_code: 200, subscription_tier: 'SuperGrok' placeholder,
        grok_quota_snapshot: { subscription_tier: 'Free' placeholder
      placeholder
    placeholder as any
    const reordered = {
      platform: 'grok',
      extra: {
        grok_quota_snapshot: { subscription_tier: 'SuperGrok Heavy' placeholder,
        grok_usage_snapshot: { subscription_tier: 'SuperGrok', status_code: 200 placeholder,
        grok_billing_snapshot: {
          limits: { weekly: 25, monthly: 100 placeholder,
          plan: 'SuperGrok'
        placeholder
      placeholder
    placeholder as any

    expect(buildGrokUsageRefreshKey(first)).toBe(buildGrokUsageRefreshKey(reordered))
  placeholder)

  it('uses the legacy quota alias only when the canonical snapshot is absent', () => {
    const base = {
      platform: 'grok',
      extra: { grok_quota_snapshot: { subscription_tier: 'Free' placeholder placeholder
    placeholder as any
    const next = {
      platform: 'grok',
      extra: { grok_quota_snapshot: { subscription_tier: 'SuperGrok' placeholder placeholder
    placeholder as any

    expect(buildGrokUsageRefreshKey(base)).not.toBe(buildGrokUsageRefreshKey(next))
  placeholder)

  it('tracks the legacy tier when the canonical snapshot has no usable tier', () => {
    for (const canonicalSnapshot of [
      { status_code: 200 placeholder,
      { status_code: 200, subscription_tier: '   ' placeholder,
    ]) {
      const base = {
        platform: 'grok',
        extra: {
          grok_usage_snapshot: canonicalSnapshot,
          grok_quota_snapshot: { subscription_tier: 'Free' placeholder,
        placeholder,
      placeholder as any
      const next = {
        ...base,
        extra: {
          ...base.extra,
          grok_quota_snapshot: { subscription_tier: 'SuperGrok' placeholder,
        placeholder,
      placeholder

      expect(buildGrokUsageRefreshKey(base)).not.toBe(buildGrokUsageRefreshKey(next))
    placeholder
  placeholder)

  it('returns an empty key for non-Grok accounts', () => {
    expect(buildGrokUsageRefreshKey({
      platform: 'openai',
      extra: { grok_usage_snapshot: { subscription_tier: 'SuperGrok' placeholder placeholder
    placeholder as any)).toBe('')
  placeholder)
placeholder)
