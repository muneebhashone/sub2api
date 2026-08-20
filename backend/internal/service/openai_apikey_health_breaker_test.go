package service

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type openAIAPIKeyHealthSettingRepo struct {
	SettingRepository
	value string
placeholder

func (r *openAIAPIKeyHealthSettingRepo) GetValue(context.Context, string) (string, error) {
	return r.value, nil
placeholder

type openAIAPIKeyHealthAccountRepo struct {
	AccountRepository
	setCalls int
	reason   string
placeholder

func (r *openAIAPIKeyHealthAccountRepo) SetTempUnschedulable(_ context.Context, _ int64, _ time.Time, reason string) error {
	r.setCalls++
	r.reason = reason
	return nil
placeholder

type openAIAPIKeyHealthCacheStub struct {
	TempUnschedCache
	recordCalls int
	resetCalls  int
	setCalls    int
	tripped     bool
placeholder

func (c *openAIAPIKeyHealthCacheStub) RecordOpenAIAPIKeyHealthFailure(context.Context, int64, int, int) (int64, bool, error) {
	c.recordCalls++
	return 3, c.tripped, nil
placeholder

func (c *openAIAPIKeyHealthCacheStub) ResetOpenAIAPIKeyHealthFailures(context.Context, int64) error {
	c.resetCalls++
	return nil
placeholder

func (c *openAIAPIKeyHealthCacheStub) SetTempUnsched(context.Context, int64, *TempUnschedState) error {
	c.setCalls++
	return nil
placeholder

type openAIAPIKeyHealthRuntimeBlocker struct{ calls int placeholder

func (b *openAIAPIKeyHealthRuntimeBlocker) BlockAccountScheduling(*Account, time.Time, string) {
	b.calls++
placeholder
func (*openAIAPIKeyHealthRuntimeBlocker) ClearAccountSchedulingBlock(int64) {placeholder

func openAIHealthPoolAccount() *Account {
placeholder
		ID:       42,
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
placeholder
			"pool_mode": true,
	placeholder,
placeholder
placeholder

func TestClassifyOpenAIAPIKeyHealthFailureExclusions(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		eligible bool
placeholder{
		{name: "account attributed 502", err: &UpstreamFailoverError{StatusCode: http.StatusBadGatewayplaceholder, eligible: trueplaceholder,
		{name: "request scoped capacity", err: &UpstreamFailoverError{StatusCode: 529, RequestScopedTransient: trueplaceholderplaceholder,
		{name: "provider scoped overload", err: &UpstreamFailoverError{StatusCode: 529, Scope: GatewayFailureScopeProviderplaceholderplaceholder,
		{name: "dedicated same account retry", err: &UpstreamFailoverError{StatusCode: http.StatusTooManyRequests, RetryableOnSameAccount: trueplaceholderplaceholder,
		{name: "credential disable path", err: &UpstreamFailoverError{StatusCode: http.StatusUnauthorized, Stage: GatewayFailureStageAccountAuth, Scope: GatewayFailureScopeAccountplaceholderplaceholder,
		{name: "client request", err: &UpstreamFailoverError{StatusCode: http.StatusBadRequestplaceholderplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, eligible := classifyOpenAIAPIKeyHealthFailure(tt.err)
			require.Equal(t, tt.eligible, eligible)
	placeholder)
placeholder
placeholder

func TestOpenAIAPIKeyHealthBreakerDefaultDisabled(t *testing.T) {
	settings := NewSettingService(&openAIAPIKeyHealthSettingRepo{placeholder, &config.Config{placeholder)
	cache := &openAIAPIKeyHealthCacheStub{tripped: trueplaceholder
	svc := NewRateLimitService(&openAIAPIKeyHealthAccountRepo{placeholder, nil, &config.Config{placeholder, nil, cache)
	svc.SetSettingService(settings)
	svc.SetOpenAIAPIKeyHealthCache(cache)

	require.False(t, svc.ObserveOpenAIAPIKeyHealthFailure(context.Background(), openAIHealthPoolAccount(), &UpstreamFailoverError{StatusCode: http.StatusBadGatewayplaceholder))
	require.Zero(t, cache.recordCalls)
placeholder

func TestOpenAIAPIKeyHealthBreakerTripsPersistedAndRuntimeState(t *testing.T) {
	encoded, err := json.Marshal(OpenAIAPIKeyHealthBreakerSettings{Enabled: true, WindowMinutes: 1, FailureThreshold: 3, CooldownMinutes: 5placeholder)
placeholder
	settings := NewSettingService(&openAIAPIKeyHealthSettingRepo{value: string(encoded)placeholder, &config.Config{placeholder)
	cache := &openAIAPIKeyHealthCacheStub{tripped: trueplaceholder
	repo := &openAIAPIKeyHealthAccountRepo{placeholder
	blocker := &openAIAPIKeyHealthRuntimeBlocker{placeholder
	svc := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, cache)
	svc.SetSettingService(settings)
	svc.SetOpenAIAPIKeyHealthCache(cache)
	svc.SetAccountRuntimeBlocker(blocker)
	account := openAIHealthPoolAccount()

	require.True(t, svc.ObserveOpenAIAPIKeyHealthFailure(context.Background(), account, &UpstreamFailoverError{StatusCode: http.StatusBadGateway, ResponseBody: []byte(`{"error":"upstream"placeholder`)placeholder))
	require.Equal(t, 1, cache.recordCalls)
	require.Equal(t, 1, cache.setCalls)
	require.Equal(t, 1, repo.setCalls)
	require.Equal(t, 1, blocker.calls)
	require.NotNil(t, account.TempUnschedulableUntil)
	require.Contains(t, repo.reason, openAIAPIKeyHealthBreakerReason)
placeholder

func TestOpenAIAPIKeyHealthSuccessResetsOnlyEligiblePoolAccount(t *testing.T) {
	encoded, err := json.Marshal(OpenAIAPIKeyHealthBreakerSettings{Enabled: true, WindowMinutes: 1, FailureThreshold: 3, CooldownMinutes: 5placeholder)
placeholder
	settings := NewSettingService(&openAIAPIKeyHealthSettingRepo{value: string(encoded)placeholder, &config.Config{placeholder)
	cache := &openAIAPIKeyHealthCacheStub{placeholder
	svc := NewRateLimitService(&openAIAPIKeyHealthAccountRepo{placeholder, nil, &config.Config{placeholder, nil, cache)
	svc.SetSettingService(settings)
	svc.SetOpenAIAPIKeyHealthCache(cache)

	svc.ObserveOpenAIAPIKeyHealthSuccess(context.Background(), openAIHealthPoolAccount())
	svc.ObserveOpenAIAPIKeyHealthSuccess(context.Background(), &Account{ID: 43, Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder)
	require.Equal(t, 1, cache.resetCalls)
placeholder
