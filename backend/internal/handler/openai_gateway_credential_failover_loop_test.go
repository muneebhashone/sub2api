//go:build unit

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type grokCredentialHandlerRepo struct {
	service.AccountRepository
	mu             sync.Mutex
	accounts       []service.Account
	setErrorIDs    []int64
	setTempIDs     []int64
	rateLimitIDs   []int64
	updateExtraIDs []int64
	selectionCalls int
	setErrorErr    error
	setTempErr     error
	missingOnGet   map[int64]bool
placeholder

func (r *grokCredentialHandlerRepo) ListSchedulableByPlatform(_ context.Context, platform string) ([]service.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.selectionCalls++
	out := make([]service.Account, 0, len(r.accounts))
	for _, account := range r.accounts {
		if account.Platform == platform && account.IsSchedulable() {
			out = append(out, account)
	placeholder
placeholder
	return out, nil
placeholder

func (r *grokCredentialHandlerRepo) ListSchedulableByGroupIDAndPlatform(ctx context.Context, _ int64, platform string) ([]service.Account, error) {
	return r.ListSchedulableByPlatform(ctx, platform)
placeholder

func (r *grokCredentialHandlerRepo) ListSchedulableUngroupedByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	return r.ListSchedulableByPlatform(ctx, platform)
placeholder

func (r *grokCredentialHandlerRepo) GetByID(_ context.Context, id int64) (*service.Account, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.missingOnGet[id] {
		return nil, nil
placeholder
	for _, account := range r.accounts {
		if account.ID == id {
			copy := account
			copy.Credentials = cloneCredentialMap(account.Credentials)
			return &copy, nil
	placeholder
placeholder
	return nil, nil
placeholder

func (r *grokCredentialHandlerRepo) SetError(_ context.Context, id int64, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setErrorIDs = append(r.setErrorIDs, id)
	if r.setErrorErr != nil {
		return r.setErrorErr
placeholder
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			r.accounts[i].Status = service.StatusError
			r.accounts[i].Schedulable = false
			r.accounts[i].ErrorMessage = message
	placeholder
placeholder
	return nil
placeholder

func (r *grokCredentialHandlerRepo) SetTempUnschedulable(_ context.Context, id int64, until time.Time, _ string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.setTempIDs = append(r.setTempIDs, id)
	if r.setTempErr != nil {
		return r.setTempErr
placeholder
	for i := range r.accounts {
		if r.accounts[i].ID == id {
			value := until
			r.accounts[i].TempUnschedulableUntil = &value
	placeholder
placeholder
	return nil
placeholder

func (r *grokCredentialHandlerRepo) SetRateLimited(_ context.Context, id int64, resetAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rateLimitIDs = append(r.rateLimitIDs, id)
	for i := range r.accounts {
		if r.accounts[i].ID != id {
			continue
	placeholder
		now := time.Now()
		r.accounts[i].RateLimitedAt = &now
		value := resetAt
		r.accounts[i].RateLimitResetAt = &value
placeholder
	return nil
placeholder

func (r *grokCredentialHandlerRepo) SetRateLimitedIfLater(ctx context.Context, id int64, resetAt time.Time) error {
	r.mu.Lock()
	for i := range r.accounts {
		if r.accounts[i].ID == id && r.accounts[i].RateLimitResetAt != nil && !resetAt.After(*r.accounts[i].RateLimitResetAt) {
			r.mu.Unlock()
			return nil
	placeholder
placeholder
	r.mu.Unlock()
	return r.SetRateLimited(ctx, id, resetAt)
placeholder

func (r *grokCredentialHandlerRepo) SetGrokCredentialErrorIfMatch(
	_ context.Context,
	id int64,
	snapshot service.GrokCredentialMutationSnapshot,
	message string,
) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.accounts {
		account := &r.accounts[i]
		if account.ID != id || !handlerGrokCredentialSnapshotMatches(account, snapshot) {
			continue
	placeholder
		r.setErrorIDs = append(r.setErrorIDs, id)
		if r.setErrorErr != nil {
			return false, r.setErrorErr
	placeholder
		account.Status = service.StatusError
		account.Schedulable = false
		account.ErrorMessage = message
		return true, nil
placeholder
	return false, nil
placeholder

func (r *grokCredentialHandlerRepo) SetGrokCredentialTempUnschedulableIfMatch(
	_ context.Context,
	id int64,
	snapshot service.GrokCredentialMutationSnapshot,
	until time.Time,
	_ string,
) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.accounts {
		account := &r.accounts[i]
		if account.ID != id || !handlerGrokCredentialSnapshotMatches(account, snapshot) {
			continue
	placeholder
		r.setTempIDs = append(r.setTempIDs, id)
		if r.setTempErr != nil {
			return false, r.setTempErr
	placeholder
		value := until
		account.TempUnschedulableUntil = &value
		return true, nil
placeholder
	return false, nil
placeholder

func handlerGrokCredentialSnapshotMatches(account *service.Account, snapshot service.GrokCredentialMutationSnapshot) bool {
	if account == nil {
		return false
placeholder
	credentialsJSON, err := json.Marshal(account.Credentials)
	return err == nil && account.IsGrokOAuth() && account.IsSchedulable() && string(credentialsJSON) == snapshot.CredentialsJSON &&
		handlerGrokCredentialProxyIDsEqual(account.ProxyID, snapshot.ProxyID)
placeholder

func handlerGrokCredentialProxyIDsEqual(left, right *int64) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
placeholder
	return *left == *right
placeholder

func (r *grokCredentialHandlerRepo) UpdateExtra(_ context.Context, id int64, updates map[string]any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.updateExtraIDs = append(r.updateExtraIDs, id)
	for i := range r.accounts {
		if r.accounts[i].ID != id {
			continue
	placeholder
		if r.accounts[i].Extra == nil {
			r.accounts[i].Extra = map[string]any{placeholder
	placeholder
		for key, value := range updates {
			r.accounts[i].Extra[key] = value
	placeholder
placeholder
	return nil
placeholder

func (r *grokCredentialHandlerRepo) errorIDs() []int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]int64(nil), r.setErrorIDs...)
placeholder

func (r *grokCredentialHandlerRepo) selectorCalls() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.selectionCalls
placeholder

func (r *grokCredentialHandlerRepo) rateLimitedAccountIDs() []int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]int64(nil), r.rateLimitIDs...)
placeholder

type grokCredentialHandlerTokenCache struct {
	service.GrokTokenCache
	mu        sync.Mutex
	deleteErr error
placeholder

func (c *grokCredentialHandlerTokenCache) GetAccessToken(context.Context, string) (string, error) {
	return "", errors.New("not cached")
placeholder

func (c *grokCredentialHandlerTokenCache) SetAccessToken(context.Context, string, string, time.Duration) error {
	return nil
placeholder

func (c *grokCredentialHandlerTokenCache) DeleteAccessToken(context.Context, string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.deleteErr
placeholder

func (c *grokCredentialHandlerTokenCache) AcquireRefreshLock(context.Context, string, time.Duration) (bool, error) {
	return true, nil
placeholder

func (c *grokCredentialHandlerTokenCache) ReleaseRefreshLock(context.Context, string) error {
	return nil
placeholder

func cloneCredentialMap(source map[string]any) map[string]any {
	cloned := make(map[string]any, len(source))
	for key, value := range source {
		cloned[key] = value
placeholder
	return cloned
placeholder

type grokCredentialHandlerRefresher struct {
	mode    string
	started chan struct{placeholder
	once    sync.Once
placeholder

func (r *grokCredentialHandlerRefresher) CacheKey(account *service.Account) string {
	return service.GrokTokenCacheKey(account)
placeholder

func (r *grokCredentialHandlerRefresher) CanRefresh(account *service.Account) bool {
	return account != nil && account.IsGrokOAuth()
placeholder

func (r *grokCredentialHandlerRefresher) NeedsRefresh(account *service.Account, _ time.Duration) bool {
	return account != nil && (account.ID == 801 || r.mode == "all_revoked")
placeholder

func (r *grokCredentialHandlerRefresher) Refresh(ctx context.Context, _ *service.Account) (map[string]any, error) {
	switch r.mode {
	case "revoked", "all_revoked", "mutation_set_error", "mutation_cache":
		return nil, infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_grant")
	case "provider":
		return nil, infraerrors.New(http.StatusBadGateway, "GROK_OAUTH_TOKEN_REFRESH_FAILED", "invalid_client")
	case "cancel":
		r.once.Do(func() { close(r.started) placeholder)
		<-ctx.Done()
		return nil, ctx.Err()
	case "transient", "mutation_temp":
		return nil, errors.New("temporary refresh transport failure")
	default:
		return nil, nil
placeholder
placeholder

type grokCredentialHandlerUpstream struct {
	service.HTTPUpstream
	mu            sync.Mutex
	hits          []int64
	requestURLs   []string
	authorization []string
	failAccountID int64
	rateLimitIDs  map[int64]bool
	failureStatus map[int64]int
	cancelRequest context.CancelFunc
placeholder

func (u *grokCredentialHandlerUpstream) Do(req *http.Request, _ string, accountID int64, _ int) (*http.Response, error) {
	var requestBody []byte
	if req.Body != nil {
		requestBody, _ = io.ReadAll(req.Body)
placeholder
	u.mu.Lock()
	u.hits = append(u.hits, accountID)
	u.requestURLs = append(u.requestURLs, req.URL.String())
	u.authorization = append(u.authorization, req.Header.Get("Authorization"))
	failAccountID := u.failAccountID
	rateLimited := u.rateLimitIDs[accountID]
	failureStatus := u.failureStatus[accountID]
	cancelRequest := u.cancelRequest
	u.mu.Unlock()
	if rateLimited {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Header: http.Header{
				"Content-Type": []string{"application/json"placeholder,
				"Retry-After":  []string{"60"placeholder,
		placeholder,
			Body: io.NopCloser(bytes.NewBufferString(`{"error":{"message":"rate limited"placeholderplaceholder`)),
	placeholder, nil
placeholder
	if failureStatus > 0 {
		return &http.Response{
			StatusCode: failureStatus,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"message":"upstream unavailable"placeholderplaceholder`)),
	placeholder, nil
placeholder
	if accountID == failAccountID {
		if cancelRequest != nil {
			cancelRequest()
	placeholder
		return &http.Response{
			StatusCode: http.StatusPaymentRequired,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"message":"payment required"placeholderplaceholder`)),
	placeholder, nil
placeholder
	if bytes.Contains(requestBody, []byte(`"stream":true`)) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
			Body: io.NopCloser(bytes.NewBufferString(
				"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_healthy\",\"model\":\"grok-4.5\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1placeholderplaceholderplaceholder\n\n",
			)),
	placeholder, nil
placeholder
	if strings.Contains(req.URL.Path, "/chat/completions") {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
			Body: io.NopCloser(bytes.NewBufferString(
				`{"id":"chatcmpl_healthy","object":"chat.completion","model":"grok-4.5","choices":[{"index":0,"message":{"role":"assistant","content":"ok"placeholder,"finish_reason":"stop"placeholder],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2placeholderplaceholder`,
			)),
	placeholder, nil
placeholder
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body: io.NopCloser(bytes.NewBufferString(
			`{"id":"resp_healthy","object":"response","model":"grok-4.5","status":"completed","output":[],"usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`,
		)),
placeholder, nil
placeholder

func (u *grokCredentialHandlerUpstream) accountHits() []int64 {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]int64(nil), u.hits...)
placeholder

func (u *grokCredentialHandlerUpstream) requests() ([]string, []string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	return append([]string(nil), u.requestURLs...), append([]string(nil), u.authorization...)
placeholder

func TestResponsesCredentialFailoverLoop(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("revoked account selects healthy account", func(t *testing.T) {
		h, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "revoked")
		defer cleanup()
		_ = h

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		require.Contains(t, recorder.Body.String(), "resp_healthy")
		require.Equal(t, []int64{801placeholder, repo.errorIDs())
		require.Equal(t, []int64{802placeholder, upstream.accountHits())
		requestURLs, authorization := upstream.requests()
		require.Equal(t, []string{xai.DefaultCLIBaseURL + "/responses"placeholder, requestURLs)
		require.Equal(t, []string{"Bearer healthy-access"placeholder, authorization)
placeholder)

	t.Run("provider configuration stops before healthy account", func(t *testing.T) {
		h, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "provider")
		defer cleanup()

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
		require.Contains(t, recorder.Body.String(), service.GrokCredentialUnavailableClientMessage)
		require.Empty(t, repo.errorIDs())
		require.Empty(t, upstream.accountHits())
		require.Equal(t, 1, repo.selectorCalls())
		require.Zero(t, h.gatewayService.SnapshotOpenAIAccountSchedulerMetrics().RuntimeStatsAccountCount,
			"provider-scoped auth failure must not penalize the selected account")
placeholder)

	t.Run("parent cancellation stops before healthy account", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "cancel")
		defer cleanup()

		ctx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		done := make(chan struct{placeholder)
		go func() {
			defer close(done)
			router.ServeHTTP(recorder, req)
	placeholder()

		select {
		case <-time.After(2 * time.Second):
			t.Fatal("credential refresh did not start")
		case <-findHandlerRefresherStarted(router):
			cancel()
	placeholder
		select {
		case <-time.After(2 * time.Second):
			t.Fatal("handler did not stop after cancellation")
		case <-done:
	placeholder

		require.Empty(t, repo.errorIDs())
		require.Empty(t, upstream.accountHits())
placeholder)

	t.Run("post-mapping cancellation stops before scheduler mutation or reselection", func(t *testing.T) {
		h, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "postmap_cancel")
		defer cleanup()
		ctx, cancel := context.WithCancel(context.Background())
		upstream.mu.Lock()
		upstream.failAccountID = 801
		upstream.cancelRequest = cancel
		upstream.mu.Unlock()

		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`)).WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(recorder, req)

		require.Equal(t, []int64{801placeholder, upstream.accountHits())
		require.Empty(t, repo.errorIDs())
		require.Equal(t, 1, repo.selectorCalls())
		require.Zero(t, h.gatewayService.SnapshotOpenAIAccountSchedulerMetrics().RuntimeStatsAccountCount)
placeholder)

	t.Run("pre-cancelled request never invokes an account selector", func(t *testing.T) {
		tests := []struct {
			name   string
			method string
			path   string
			body   string
	placeholder{
			{name: "responses", method: http.MethodPost, path: "/openai/v1/responses", body: `{"model":"grok","input":"hello","stream":falseplaceholder`placeholder,
			{name: "messages", method: http.MethodPost, path: "/openai/v1/messages", body: `{"model":"grok","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder]placeholder`placeholder,
			{name: "chat completions", method: http.MethodPost, path: "/openai/v1/chat/completions", body: `{"model":"grok","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`placeholder,
			{name: "grok media", method: http.MethodGet, path: "/openai/v1/videos/request-1"placeholder,
	placeholder
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "revoked")
				defer cleanup()
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				recorder := httptest.NewRecorder()
				req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body)).WithContext(ctx)
				req.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, req)

				require.Zero(t, repo.selectorCalls())
				require.Empty(t, upstream.accountHits())
		placeholder)
	placeholder
placeholder)

	t.Run("credential state mutation failures stop before reselection", func(t *testing.T) {
		for _, mode := range []string{"mutation_set_error", "mutation_temp", "mutation_cache"placeholder {
			t.Run(mode, func(t *testing.T) {
				_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, mode)
				defer cleanup()

				recorder := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
				req.Header.Set("Content-Type", "application/json")
				router.ServeHTTP(recorder, req)

				require.Equal(t, http.StatusServiceUnavailable, recorder.Code, recorder.Body.String())
				require.Contains(t, recorder.Body.String(), service.GrokCredentialUnavailableClientMessage)
				require.Empty(t, upstream.accountHits())
				require.Equal(t, 1, repo.selectorCalls())
		placeholder)
	placeholder
placeholder)

	t.Run("missing credential provider stops before upstream or reselection", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "nil_provider")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusServiceUnavailable, recorder.Code, recorder.Body.String())
		require.Contains(t, recorder.Body.String(), service.GrokCredentialUnavailableClientMessage)
		require.Equal(t, 1, repo.selectorCalls())
		require.Empty(t, upstream.accountHits())
		require.Empty(t, repo.errorIDs())
placeholder)
placeholder

func TestResponsesGrok429FailoverIsBounded(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("first rate limited account selects healthy account", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "first_429")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		require.Contains(t, recorder.Body.String(), "resp_healthy")
		require.Equal(t, []int64{801, 802placeholder, upstream.accountHits())
		require.Equal(t, []int64{801placeholder, repo.rateLimitedAccountIDs())
placeholder)

	t.Run("two rate limited accounts stop without sweeping the pool", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "all_429")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusTooManyRequests, recorder.Code, recorder.Body.String())
		require.Equal(t, []int64{801, 802placeholder, upstream.accountHits())
		require.Equal(t, []int64{801, 802placeholder, repo.rateLimitedAccountIDs())
		require.NotContains(t, recorder.Body.String(), "expired")
		require.NotContains(t, recorder.Body.String(), "healthy-access")
		require.NotContains(t, recorder.Body.String(), "rate limited")
placeholder)
placeholder

func TestResponsesGrok429FailoverHandlesMixedStatuses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("429 then 500 stops after the bounded followup", func(t *testing.T) {
		_, _, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "mixed_429_500")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusBadGateway, recorder.Code, recorder.Body.String())
		require.Equal(t, []int64{801, 802placeholder, upstream.accountHits())
		require.NotContains(t, recorder.Body.String(), "upstream unavailable")
placeholder)

	t.Run("500 then 429 permits one healthy followup", func(t *testing.T) {
		_, _, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "mixed_500_429")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		require.Equal(t, []int64{801, 802, 803placeholder, upstream.accountHits())
placeholder)

	t.Run("OAuth 429 then API-key failure cannot bypass the bound", func(t *testing.T) {
		_, _, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "oauth_429_apikey_500")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusBadGateway, recorder.Code, recorder.Body.String())
		require.Equal(t, []int64{801, 802placeholder, upstream.accountHits())
placeholder)
placeholder

func TestGrokMedia429FailoverIsBounded(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("first 429 selects one healthy followup", func(t *testing.T) {
		_, _, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "first_429")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/videos/generations", bytes.NewBufferString(`{"model":"grok-imagine-video","prompt":"waves"placeholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
		require.Equal(t, []int64{801, 802placeholder, upstream.accountHits())
placeholder)

	t.Run("second 429 stops without sweeping a third account", func(t *testing.T) {
		_, _, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "all_429")
		defer cleanup()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/videos/generations", bytes.NewBufferString(`{"model":"grok-imagine-video","prompt":"waves"placeholder`))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusTooManyRequests, recorder.Code, recorder.Body.String())
		require.Equal(t, []int64{801, 802placeholder, upstream.accountHits())
		require.NotContains(t, recorder.Body.String(), "rate limited")
placeholder)
placeholder

func TestGrokOAuthCredentialFailoverAcrossHTTPHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	endpoints := []struct {
		name   string
		method string
		path   string
		body   string
placeholder{
		{name: "messages", method: http.MethodPost, path: "/openai/v1/messages", body: `{"model":"grok","max_tokens":16,"messages":[{"role":"user","content":"hello"placeholder]placeholder`placeholder,
		{name: "chat completions", method: http.MethodPost, path: "/openai/v1/chat/completions", body: `{"model":"grok","messages":[{"role":"user","content":"hello"placeholder],"stream":falseplaceholder`placeholder,
		{name: "chat completions raw fallback", method: http.MethodPost, path: "/openai/v1/chat/completions", body: `{"model":"grok","messages":[{"role":"user","content":"hello"placeholder],"stop":["END"],"stream":falseplaceholder`placeholder,
		{name: "grok media", method: http.MethodPost, path: "/openai/v1/videos/generations", body: `{"model":"grok-imagine-video","prompt":"waves"placeholder`placeholder,
placeholder

	for _, endpoint := range endpoints {
		t.Run(endpoint.name+" revoked selects healthy", func(t *testing.T) {
			_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "revoked")
			defer cleanup()
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(endpoint.method, endpoint.path, bytes.NewBufferString(endpoint.body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
			require.Equal(t, []int64{801placeholder, repo.errorIDs())
			require.Equal(t, []int64{802placeholder, upstream.accountHits())
	placeholder)

		t.Run(endpoint.name+" all accounts exhausted safely", func(t *testing.T) {
			_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "all_revoked")
			defer cleanup()
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(endpoint.method, endpoint.path, bytes.NewBufferString(endpoint.body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusServiceUnavailable, recorder.Code, recorder.Body.String())
			require.Contains(t, recorder.Body.String(), service.GrokCredentialUnavailableClientMessage)
			require.NotContains(t, recorder.Body.String(), "revoked-refresh")
			require.NotContains(t, recorder.Body.String(), "healthy-refresh")
			require.Equal(t, []int64{801, 802placeholder, repo.errorIDs())
			require.Empty(t, upstream.accountHits())
	placeholder)
placeholder
placeholder

func TestGrokOAuthMissingSelectedRowRetriesHealthyAccountWithoutMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "missing_row")
	defer cleanup()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/openai/v1/responses", bytes.NewBufferString(`{"model":"grok","input":"hello","stream":falseplaceholder`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, req)

	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	require.Equal(t, []int64{802placeholder, upstream.accountHits())
	require.Empty(t, repo.errorIDs())
	require.Empty(t, repo.setTempIDs)
placeholder

func TestResponsesWebSocketCredentialFailoverLoop(t *testing.T) {
	gin.SetMode(gin.TestMode)
	dial := func(t *testing.T, router *gin.Engine) (*coderws.Conn, func()) {
	placeholder
		server := httptest.NewServer(router)
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, _, err := coderws.Dial(ctx, "ws"+strings.TrimPrefix(server.URL, "http")+"/openai/v1/responses", nil)
		cancel()
	placeholder
		return conn, func() {
			_ = conn.CloseNow()
			server.Close()
	placeholder
placeholder
	writeFirst := func(t *testing.T, conn *coderws.Conn) {
	placeholder
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		require.NoError(t, conn.Write(ctx, coderws.MessageText, []byte(`{"type":"response.create","model":"grok","input":"hello","stream":falseplaceholder`)))
placeholder

	t.Run("revoked account selects healthy account", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "revoked")
		defer cleanup()
		conn, closeConn := dial(t, router)
		defer closeConn()
		writeFirst(t, conn)

		readCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, payload, err := conn.Read(readCtx)
		cancel()
	placeholder
		require.Contains(t, string(payload), "resp_healthy")
		require.Equal(t, []int64{801placeholder, repo.errorIDs())
		require.Equal(t, 2, repo.selectorCalls())
		require.Equal(t, []int64{802placeholder, upstream.accountHits())
placeholder)

	t.Run("provider configuration stops", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "provider")
		defer cleanup()
		conn, closeConn := dial(t, router)
		defer closeConn()
		writeFirst(t, conn)

		readCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, _, err := conn.Read(readCtx)
		cancel()
		var closeErr coderws.CloseError
		require.ErrorAs(t, err, &closeErr)
		require.Contains(t, closeErr.Reason, service.GrokCredentialUnavailableClientMessage)
		require.Equal(t, 1, repo.selectorCalls())
		require.Empty(t, upstream.accountHits())
placeholder)

	t.Run("parent cancellation prevents reselection", func(t *testing.T) {
		_, repo, upstream, router, cleanup := newGrokCredentialFailoverHandler(t, "cancel")
		defer cleanup()
		conn, closeConn := dial(t, router)
		writeFirst(t, conn)
		select {
		case <-findHandlerRefresherStarted(router):
		case <-time.After(2 * time.Second):
			t.Fatal("credential refresh did not start")
	placeholder
		closeConn()

		require.Eventually(t, func() bool { return repo.selectorCalls() == 1 placeholder, 2*time.Second, 20*time.Millisecond)
		require.Empty(t, repo.errorIDs())
		require.Empty(t, upstream.accountHits())
placeholder)
placeholder

var handlerRefresherStarted sync.Map

func findHandlerRefresherStarted(router *gin.Engine) <-chan struct{placeholder {
	value, _ := handlerRefresherStarted.Load(router)
	return value.(chan struct{placeholder)
placeholder

func newGrokCredentialFailoverHandler(t *testing.T, mode string) (*OpenAIGatewayHandler, *grokCredentialHandlerRepo, *grokCredentialHandlerUpstream, *gin.Engine, func()) {
placeholder
	groupID := int64(901)
	accounts := []service.Account{
		{
			ID: 801, Name: "revoked", Platform: service.PlatformGrok, Type: service.AccountTypeOAuth,
			Status: service.StatusActive, Schedulable: true, Concurrency: 1, Priority: 1,
	placeholder
				"access_token": "expired", "refresh_token": "revoked-refresh",
				"expires_at": time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
		placeholder,
			Extra: map[string]any{service.GrokMediaEligibleExtraKey: trueplaceholder,
	placeholder,
		{
			ID: 802, Name: "healthy", Platform: service.PlatformGrok, Type: service.AccountTypeOAuth,
			Status: service.StatusActive, Schedulable: true, Concurrency: 1, Priority: 2,
	placeholder
				"access_token": "healthy-access", "refresh_token": "healthy-refresh",
				"expires_at": time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
		placeholder,
			Extra: map[string]any{service.GrokMediaEligibleExtraKey: trueplaceholder,
	placeholder,
placeholder
	if mode == "postmap_cancel" || mode == "first_429" || mode == "all_429" || mode == "mixed_429_500" || mode == "mixed_500_429" || mode == "oauth_429_apikey_500" {
		accounts[0].Credentials["expires_at"] = time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339)
placeholder
	if mode == "all_429" || mode == "mixed_429_500" || mode == "mixed_500_429" || mode == "oauth_429_apikey_500" {
		accounts = append(accounts, service.Account{
			ID: 803, Name: "untried-healthy", Platform: service.PlatformGrok, Type: service.AccountTypeOAuth,
			Status: service.StatusActive, Schedulable: true, Concurrency: 1, Priority: 3,
	placeholder
				"access_token": "untried-healthy-access", "refresh_token": "untried-healthy-refresh",
				"expires_at": time.Now().Add(2 * time.Hour).UTC().Format(time.RFC3339),
		placeholder,
			Extra: map[string]any{service.GrokMediaEligibleExtraKey: trueplaceholder,
	placeholder)
placeholder
	if mode == "oauth_429_apikey_500" {
		accounts[1].Type = service.AccountTypeAPIKey
		accounts[1].Credentials = map[string]any{"api_key": "third-party-key"placeholder
placeholder
	if mode == "all_revoked" {
		accounts[1].Credentials["expires_at"] = time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
placeholder
	repo := &grokCredentialHandlerRepo{accounts: accounts, missingOnGet: map[int64]bool{placeholderplaceholder
	if mode == "missing_row" {
		repo.missingOnGet[801] = true
placeholder
	if mode == "mutation_set_error" {
		repo.setErrorErr = errors.New("database write failed")
placeholder
	if mode == "mutation_temp" {
		repo.setTempErr = errors.New("database write failed")
placeholder
	refresher := &grokCredentialHandlerRefresher{mode: mode, started: make(chan struct{placeholder)placeholder
	tokenCache := &grokCredentialHandlerTokenCache{placeholder
	if mode == "mutation_cache" {
		tokenCache.deleteErr = errors.New("cache delete failed")
placeholder
	var provider *service.GrokTokenProvider
	if mode != "nil_provider" {
		provider = service.NewGrokTokenProvider(repo, tokenCache)
		provider.SetRefreshAPI(service.NewOAuthRefreshAPI(repo, tokenCache), refresher)
placeholder
	upstream := &grokCredentialHandlerUpstream{placeholder
	switch mode {
	case "first_429":
		upstream.rateLimitIDs = map[int64]bool{801: trueplaceholder
	case "all_429":
		upstream.rateLimitIDs = map[int64]bool{801: true, 802: trueplaceholder
	case "mixed_429_500":
		upstream.rateLimitIDs = map[int64]bool{801: trueplaceholder
		upstream.failureStatus = map[int64]int{802: http.StatusInternalServerErrorplaceholder
	case "mixed_500_429":
		upstream.failureStatus = map[int64]int{801: http.StatusInternalServerErrorplaceholder
		upstream.rateLimitIDs = map[int64]bool{802: trueplaceholder
	case "oauth_429_apikey_500":
		upstream.rateLimitIDs = map[int64]bool{801: trueplaceholder
		upstream.failureStatus = map[int64]int{802: http.StatusInternalServerErrorplaceholder
placeholder
	cfg := &config.Config{RunMode: config.RunModeSimpleplaceholder
	cfg.Gateway.MaxAccountSwitches = 3
	billingCache := service.NewBillingCacheService(nil, nil, nil, nil, nil, nil, cfg, nil)
	gateway := service.NewOpenAIGatewayService(
		repo, nil, nil, nil, nil, nil, nil, cfg, nil, nil,
		service.NewBillingService(cfg, nil), nil, billingCache, upstream,
		&service.DeferredService{placeholder, nil, provider, nil, nil, nil, nil, nil,
	)
	cache := &concurrencyCacheMock{
		acquireUserSlotFn:    func(context.Context, int64, int, string) (bool, error) { return true, nil placeholder,
		acquireAccountSlotFn: func(context.Context, int64, int, string) (bool, error) { return true, nil placeholder,
placeholder
	h := NewOpenAIGatewayHandler(gateway, service.NewConcurrencyService(cache), billingCache, &service.APIKeyService{placeholder, nil, nil, nil, nil, cfg)
	apiKey := &service.APIKey{
		ID: 902, GroupID: &groupID,
		User:  &service.User{ID: 903, Status: service.StatusActiveplaceholder,
		Group: &service.Group{ID: groupID, Platform: service.PlatformGrok, Status: service.StatusActive, AllowImageGeneration: trueplaceholder,
placeholder
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware.ContextKeyAPIKey), apiKey)
		c.Set(string(middleware.ContextKeyUser), middleware.AuthSubject{UserID: apiKey.User.ID, Concurrency: 1placeholder)
		c.Next()
placeholder)
	router.POST("/openai/v1/responses", h.Responses)
	router.GET("/openai/v1/responses", h.ResponsesWebSocket)
	router.POST("/openai/v1/messages", h.Messages)
	router.POST("/openai/v1/chat/completions", h.ChatCompletions)
	router.POST("/openai/v1/videos/generations", h.GrokVideoGeneration)
	router.GET("/openai/v1/videos/:request_id", h.GrokVideoStatus)
	handlerRefresherStarted.Store(router, refresher.started)
	cleanup := func() {
		handlerRefresherStarted.Delete(router)
		billingCache.Stop()
placeholder
	return h, repo, upstream, router, cleanup
placeholder
