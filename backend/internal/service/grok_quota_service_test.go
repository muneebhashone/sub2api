//go:build unit

package service

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type grokQuotaAccountRepo struct {
	*mockAccountRepoForPlatform
	updates               map[int64]map[string]any
	updateCalls           int
	rateLimitedCalls      int
	lastRateLimitedID     int64
	lastRateLimitResetAt  time.Time
	tempUnschedCalls      int
	lastTempUnschedID     int64
	lastTempUnschedUntil  time.Time
	lastTempUnschedReason string
placeholder

func (r *grokQuotaAccountRepo) UpdateExtra(_ context.Context, id int64, updates map[string]any) error {
	r.updateCalls++
	if r.updates == nil {
		r.updates = make(map[int64]map[string]any)
placeholder
	r.updates[id] = updates
	return nil
placeholder

func (r *grokQuotaAccountRepo) SetRateLimited(_ context.Context, id int64, resetAt time.Time) error {
	r.rateLimitedCalls++
	r.lastRateLimitedID = id
	r.lastRateLimitResetAt = resetAt
	return nil
placeholder

func (r *grokQuotaAccountRepo) SetRateLimitedIfLater(ctx context.Context, id int64, resetAt time.Time) error {
	return r.SetRateLimited(ctx, id, resetAt)
placeholder

func (r *grokQuotaAccountRepo) SetTempUnschedulable(_ context.Context, id int64, until time.Time, reason string) error {
	r.tempUnschedCalls++
	r.lastTempUnschedID = id
	r.lastTempUnschedUntil = until
	r.lastTempUnschedReason = reason
	return nil
placeholder

type grokQuotaProxyRepo struct {
	proxyRepoStub
	proxies map[int64]*Proxy
	calls   int
placeholder

func (r *grokQuotaProxyRepo) GetByID(_ context.Context, id int64) (*Proxy, error) {
	r.calls++
	return r.proxies[id], nil
placeholder

func TestGrokQuotaServiceProbeUsageStoresHeaders(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:          42,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{42: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Ratelimit-Limit-Requests":     []string{"10"placeholder,
			"X-Ratelimit-Remaining-Requests": []string{"7"placeholder,
			"X-Ratelimit-Reset-Requests":     []string{"2000000000"placeholder,
			"X-Ratelimit-Limit-Tokens":       []string{"1000"placeholder,
			"X-Ratelimit-Remaining-Tokens":   []string{"900"placeholder,
	placeholder,
		Body: io.NopCloser(strings.NewReader(`{"id":"resp_probe"placeholder`)),
placeholderplaceholder
	svc := NewGrokQuotaService(repo, nil, NewGrokTokenProvider(repo, nil), upstream)

	result, err := svc.ProbeUsage(context.Background(), 42)
placeholder
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.Equal(t, "grok-4.3", result.Model)
	require.True(t, result.HeadersObserved)
	require.NotNil(t, result.Snapshot)
	require.True(t, result.Snapshot.HeadersObserved)
	require.Equal(t, "active_probe", result.Snapshot.ObservationSource)
	require.NotEmpty(t, result.Snapshot.LastProbeAt)
	require.NotEmpty(t, result.Snapshot.LastHeadersSeenAt)
	require.NotNil(t, result.Snapshot.Requests)
	require.EqualValues(t, 10, *result.Snapshot.Requests.Limit)
	require.EqualValues(t, 7, *result.Snapshot.Requests.Remaining)
	require.Equal(t, "https://cli-chat-proxy.grok.com/v1/responses", upstream.lastReq.URL.String())
	require.Equal(t, "Bearer access-token", upstream.lastReq.Header.Get("Authorization"))
	require.Equal(t, grokCLIVersion, upstream.lastReq.Header.Get("X-Grok-Client-Version"))
	require.Equal(t, "grok-4.3", gjson.GetBytes(upstream.lastBody, "model").String())
	require.Contains(t, string(upstream.lastBody), `"max_output_tokens":1`)
	require.Contains(t, string(upstream.lastBody), `"store":false`)
	require.NotNil(t, repo.updates[42][grokQuotaSnapshotExtraKey])
placeholder

func TestGrokQuotaServiceProbeUsageIgnoresAccountGrokMapping(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:          47,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"model_mapping": map[string]any{
				"grok":          "grok-composer",
				"grok-composer": "grok-composer-2.5-fast",
		placeholder,
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{47: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"resp_probe"placeholder`)),
placeholderplaceholder
	svc := NewGrokQuotaService(repo, nil, NewGrokTokenProvider(repo, nil), upstream)

	result, err := svc.ProbeUsage(context.Background(), 47)
placeholder
	require.Equal(t, "grok-4.3", result.Model)
	require.Equal(t, "grok-4.3", gjson.GetBytes(upstream.lastBody, "model").String())
	require.NotContains(t, string(upstream.lastBody), "grok-composer")
placeholder

func TestGrokQuotaServiceProbeUsageReportsProbeModelOnUpstreamError(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:          48,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{48: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(`{"code":"invalid-argument","error":"Model not found"placeholder`)),
placeholderplaceholder
	svc := NewGrokQuotaService(repo, nil, NewGrokTokenProvider(repo, nil), upstream)

	_, err := svc.ProbeUsage(context.Background(), 48)
placeholder
	require.Equal(t, "GROK_QUOTA_PROBE_UPSTREAM_ERROR", infraerrors.Reason(err))
	require.Contains(t, infraerrors.Message(err), `probe model "grok-4.3"`)
placeholder

func TestGrokQuotaServiceProbeUsageLoadsProxyWhenAccountEdgeMissing(t *testing.T) {
	t.Parallel()

	proxyID := int64(7)
	account := &Account{
		ID:          46,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
		ProxyID:     &proxyID,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{46: accountplaceholder,
	placeholder,
placeholder
	proxyRepo := &grokQuotaProxyRepo{
		proxies: map[int64]*Proxy{
			proxyID: {
				ID:       proxyID,
				Protocol: "http",
				Host:     "proxy.test",
				Port:     3128,
		placeholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"resp_probe"placeholder`)),
placeholderplaceholder
	svc := NewGrokQuotaService(repo, proxyRepo, NewGrokTokenProvider(repo, nil), upstream)

	_, err := svc.ProbeUsage(context.Background(), 46)
placeholder
	require.Equal(t, 1, proxyRepo.calls)
	require.Equal(t, "http://proxy.test:3128", upstream.lastProxyURL)
placeholder

func TestGrokQuotaServiceProbeUsageStoresNoHeadersState(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:          45,
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
		Concurrency: 1,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{45: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(`{"id":"resp_probe"placeholder`)),
placeholderplaceholder
	svc := NewGrokQuotaService(repo, nil, NewGrokTokenProvider(repo, nil), upstream)

	result, err := svc.ProbeUsage(context.Background(), 45)
placeholder
	require.Equal(t, http.StatusOK, result.StatusCode)
	require.False(t, result.HeadersObserved)
	require.NotNil(t, result.Snapshot)
	require.False(t, result.Snapshot.HeadersObserved)
	require.Equal(t, "active_probe", result.Snapshot.ObservationSource)
	require.NotEmpty(t, result.Snapshot.LastProbeAt)
	require.Empty(t, result.Snapshot.LastHeadersSeenAt)

	stored, ok := repo.updates[45][grokQuotaSnapshotExtraKey].(*xai.QuotaSnapshot)
	require.True(t, ok)
	require.False(t, stored.HeadersObserved)
	require.Equal(t, http.StatusOK, stored.StatusCode)
placeholder

func TestGrokQuotaServiceProbeUsageReturnsRateLimitedSnapshot(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:       43,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "access-token",
			"expires_at":   time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	placeholder,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{43: accountplaceholder,
	placeholder,
placeholder
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"45"placeholderplaceholder,
		Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"rate limited"placeholderplaceholder`)),
placeholderplaceholder
	svc := NewGrokQuotaService(repo, nil, NewGrokTokenProvider(repo, nil), upstream)

	result, err := svc.ProbeUsage(context.Background(), 43)
placeholder
	require.Equal(t, http.StatusTooManyRequests, result.StatusCode)
	require.NotNil(t, result.Snapshot)
	require.NotNil(t, result.Snapshot.RetryAfterSeconds)
	require.Equal(t, 45, *result.Snapshot.RetryAfterSeconds)
	require.Equal(t, 1, repo.rateLimitedCalls)
	require.Equal(t, account.ID, repo.lastRateLimitedID)
	require.WithinDuration(t, time.Now().Add(45*time.Second), repo.lastRateLimitResetAt, time.Second)
	require.Zero(t, repo.tempUnschedCalls)
placeholder

func TestGrokQuotaServiceResetQuotaUnsupported(t *testing.T) {
	t.Parallel()

	account := &Account{
		ID:       44,
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
	repo := &grokQuotaAccountRepo{
		mockAccountRepoForPlatform: &mockAccountRepoForPlatform{
			accountsByID: map[int64]*Account{44: accountplaceholder,
	placeholder,
placeholder
	svc := NewGrokQuotaService(repo, nil, nil, nil)

	_, err := svc.ResetQuota(context.Background(), 44)
placeholder
	require.Equal(t, http.StatusNotImplemented, infraerrors.Code(err))
	require.Equal(t, "GROK_QUOTA_RESET_UNSUPPORTED", infraerrors.Reason(err))
placeholder

func TestShouldAutoPauseGrokAccountByQuota(t *testing.T) {
	t.Parallel()

	zero := int64(0)
	limit := int64(10)
	resetFuture := time.Now().Add(time.Minute).Unix()
	retryAfter := 30
	tests := []struct {
		name     string
		snapshot xai.QuotaSnapshot
		want     bool
placeholder{
		{
			name: "remaining requests exhausted",
			snapshot: xai.QuotaSnapshot{
				Requests:  &xai.QuotaWindow{Limit: &limit, Remaining: &zero, ResetUnix: &resetFutureplaceholder,
				UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		placeholder,
			want: true,
	placeholder,
		{
			name: "retry after active",
			snapshot: xai.QuotaSnapshot{
				RetryAfterSeconds: &retryAfter,
				UpdatedAt:         time.Now().UTC().Format(time.RFC3339),
		placeholder,
			want: true,
	placeholder,
		{
			name: "retry after expired",
			snapshot: xai.QuotaSnapshot{
				RetryAfterSeconds: &retryAfter,
				UpdatedAt:         time.Now().Add(-time.Duration(retryAfter+1) * time.Second).UTC().Format(time.RFC3339),
		placeholder,
			want: false,
	placeholder,
		{
			name: "stale snapshot ignored",
			snapshot: xai.QuotaSnapshot{
				Requests:  &xai.QuotaWindow{Limit: &limit, Remaining: &zero, ResetUnix: &resetFutureplaceholder,
				UpdatedAt: time.Now().Add(-3 * time.Hour).UTC().Format(time.RFC3339),
		placeholder,
			want: false,
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			account := &Account{
				Platform: PlatformGrok,
				Type:     AccountTypeOAuth,
				Extra: map[string]any{
					grokQuotaSnapshotExtraKey: tt.snapshot,
			placeholder,
		placeholder
			got, _ := shouldAutoPauseGrokAccountByQuota(account)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder
