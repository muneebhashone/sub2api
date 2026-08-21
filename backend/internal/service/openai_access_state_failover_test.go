package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type openAIStream403AccountRepo struct {
	AccountRepository
	setErrorCalls int
placeholder

func (r *openAIStream403AccountRepo) SetError(context.Context, int64, string) error {
	r.setErrorCalls++
	return nil
placeholder

type openAIAuthPolicyAccountRepo struct {
	AccountRepository
	tempCalls     int
	setErrorCalls int
placeholder

func (r *openAIAuthPolicyAccountRepo) SetTempUnschedulable(context.Context, int64, time.Time, string) error {
	r.tempCalls++
	return nil
placeholder

func (r *openAIAuthPolicyAccountRepo) SetError(context.Context, int64, string) error {
	r.setErrorCalls++
	return nil
placeholder

type openAIAuthPolicy403Counter struct {
	counts []int64
placeholder

func (s *openAIAuthPolicy403Counter) IncrementOpenAI403Count(context.Context, int64, int) (int64, error) {
	if len(s.counts) == 0 {
		return 1, nil
placeholder
	count := s.counts[0]
	s.counts = s.counts[1:]
	return count, nil
placeholder

func (*openAIAuthPolicy403Counter) ResetOpenAI403Count(context.Context, int64) error {
	return nil
placeholder

func TestOpenAIUpstreamAccessStateClassification(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
placeholder{
		{"workspace_code", `{"detail":{"code":"deactivated_workspace"placeholderplaceholder`, trueplaceholder,
		{"disabled_account_message", `{"error":{"message":"Your account is disabled"placeholderplaceholder`, falseplaceholder,
		{"suspended_workspace_message", `{"response":{"error":{"message":"This workspace has been suspended"placeholderplaceholderplaceholder`, falseplaceholder,
		{"deactivated_organization_message", `{"detail":{"message":"The organization is deactivated"placeholderplaceholder`, falseplaceholder,
		{"scalar_detail", `{"detail":"This workspace has been disabled"placeholder`, falseplaceholder,
		{"suspended_org_code", `{"error":{"code":"org_suspended"placeholderplaceholder`, trueplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(tt.body)
			require.Equal(t, tt.want, isOpenAIUpstreamAccessStateError("", body))
			if !tt.want {
				return
		placeholder
			require.True(t, (&OpenAIGatewayService{placeholder).shouldFailoverOpenAIUpstreamResponse(http.StatusForbidden, "", body))
			require.True(t, shouldFailoverOpenAIPassthroughResponse(&Account{Type: AccountTypeOAuthplaceholder, http.StatusForbidden, body))

			err := newOpenAIUpstreamFailoverError(http.StatusForbidden, nil, body, "", true)
			require.True(t, err.IsCredentialFailure())
			require.Equal(t, GatewayFailureScopeAccount, err.Scope)
			require.Equal(t, OpenAIUpstreamAccessStateReason, err.Reason)
			require.Equal(t, NextAccountRetry, err.NextAccountAction)
			require.False(t, err.RetryableOnSameAccount)
			require.False(t, err.RequestScopedTransient)
			require.Equal(t, http.StatusBadGateway, err.ClientStatusCode)
			require.Equal(t, openAIUpstreamAccessUnavailableClientMessage, err.ClientMessage)
	placeholder)
placeholder
placeholder

func TestOpenAIUpstreamAccessStateDoesNotScanEchoedJSON(t *testing.T) {
	body := []byte(`{"error":{"code":"invalid_request_error","message":"Invalid input"placeholder,"echo":{"prompt":"my account is disabled"placeholderplaceholder`)
	require.False(t, isOpenAIUpstreamAccessStateError("", body))
	require.False(t, shouldFailoverOpenAIPassthroughResponse(&Account{Type: AccountTypeOAuthplaceholder, http.StatusBadRequest, body))
placeholder

func TestOpenAIHTTPAccessStateDoesNotTrustBadRequestMessage(t *testing.T) {
	body := []byte(`{"error":{"type":"invalid_request_error","code":"unknown_parameter","message":"Unknown parameter: account disabled"placeholderplaceholder`)
	svc := &OpenAIGatewayService{placeholder

	require.False(t, isOpenAIUpstreamAccessStateError("", body), "free-form stream messages are not durable account evidence")
	require.False(t, isOpenAIHTTPUpstreamAccessStateError(http.StatusBadRequest, "", body))
	require.False(t, svc.shouldFailoverOpenAIUpstreamResponse(http.StatusBadRequest, "", body))
	require.False(t, shouldFailoverOpenAIPassthroughResponse(&Account{Type: AccountTypeOAuthplaceholder, http.StatusBadRequest, body))

	err := newOpenAIUpstreamFailoverError(http.StatusBadRequest, nil, body, "", false)
	require.False(t, err.IsCredentialFailure())
placeholder

func TestOpenAIHTTPAccessStateBadRequestDoesNotDisableAccount(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	svc := &OpenAIGatewayService{rateLimitService: &RateLimitService{accountRepo: repoplaceholderplaceholder
	account := &Account{ID: 925, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
	body := []byte(`{"error":{"code":"unknown_parameter","message":"Unknown parameter: account disabled"placeholderplaceholder`)

	disabled := svc.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusBadRequest, nil, body)

	require.False(t, disabled)
	require.Zero(t, repo.setErrorCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIStreamEchoedAccessStateMessageDoesNotDisableOrFailover(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	svc := &OpenAIGatewayService{rateLimitService: &RateLimitService{accountRepo: repoplaceholderplaceholder
	account := &Account{ID: 926, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
	payload := []byte(`{"type":"response.failed","response":{"error":{"type":"invalid_request_error","code":"unknown_parameter","message":"Unknown parameter: account disabled"placeholderplaceholderplaceholder`)
	message := extractOpenAISSEErrorMessage(payload)

	require.False(t, isOpenAIUpstreamAccessStateError(message, payload))
	require.False(t, openAIStreamFailedEventShouldFailover(payload, message))
	status, disabled := svc.handleOpenAIStreamTerminalAccountSideEffects(nil, account, payload, message, nil)
	require.Equal(t, http.StatusBadGateway, status)
	require.False(t, disabled)
	require.Zero(t, repo.setErrorCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIHTTPAccessStateTrustsStructuredCode(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	svc := &OpenAIGatewayService{rateLimitService: &RateLimitService{accountRepo: repoplaceholderplaceholder
	account := &Account{ID: 930, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
	body := []byte(`{"error":{"code":"organization_deactivated","message":"request rejected"placeholderplaceholder`)

	require.True(t, isOpenAIHTTPUpstreamAccessStateError(http.StatusBadRequest, "", body))
	require.True(t, svc.shouldFailoverOpenAIUpstreamResponse(http.StatusBadRequest, "", body))
	require.True(t, svc.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusBadRequest, nil, body))
	require.Equal(t, 1, repo.setErrorCalls)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIHTTPAuthMessagesUseExistingStatusPolicies(t *testing.T) {
	t.Run("oauth 401 remains recoverable", func(t *testing.T) {
		repo := &openAIAuthPolicyAccountRepo{placeholder
		rateLimits := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		svc := &OpenAIGatewayService{rateLimitService: rateLimitsplaceholder
		account := &Account{ID: 931, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: true,
	placeholder"refresh_token": "refreshable"placeholderplaceholder
		body := []byte(`{"error":{"message":"account is disabled"placeholderplaceholder`)

		require.False(t, isOpenAIHTTPUpstreamAccessStateError(http.StatusUnauthorized, "", body))
		require.True(t, svc.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusUnauthorized, nil, body))
		require.Zero(t, repo.setErrorCalls)
		require.Equal(t, 1, repo.tempCalls)
placeholder)

	t.Run("403 uses counter cooldown", func(t *testing.T) {
		repo := &openAIAuthPolicyAccountRepo{placeholder
		counter := &openAIAuthPolicy403Counter{counts: []int64{1placeholderplaceholder
		rateLimits := NewRateLimitService(repo, nil, &config.Config{placeholder, nil, nil)
		rateLimits.openAI403CounterCache = counter
		svc := &OpenAIGatewayService{rateLimitService: rateLimitsplaceholder
		rateLimits.SetAccountRuntimeBlocker(svc)
		account := &Account{ID: 932, Platform: PlatformOpenAI, Type: AccountTypeOAuth, Status: StatusActive, Schedulable: trueplaceholder
		body := []byte(`{"error":{"message":"workspace has been suspended"placeholderplaceholder`)

		require.False(t, isOpenAIHTTPUpstreamAccessStateError(http.StatusForbidden, "", body))
		require.True(t, svc.handleOpenAIAccountUpstreamError(context.Background(), account, http.StatusForbidden, nil, body))
		require.Zero(t, repo.setErrorCalls)
		require.Equal(t, 1, repo.tempCalls)
placeholder)
placeholder

func TestOpenAICyberPolicyWrapped5xxNeverFailsOver(t *testing.T) {
	body := []byte(`{"error":{"code":"cyber_policy","message":"blocked"placeholderplaceholder`)
	svc := &OpenAIGatewayService{placeholder
	require.False(t, svc.shouldFailoverOpenAIUpstreamResponse(http.StatusBadGateway, "wrapped upstream failure", body))
	require.False(t, shouldFailoverOpenAIPassthroughResponse(&Account{Type: AccountTypeOAuthplaceholder, http.StatusBadGateway, body))
placeholder

func TestOpenAICapacityFailoverCarriesSafeTerminalResponse(t *testing.T) {
	message := "Our servers are currently overloaded. Please try again later."
	body := []byte(`{"error":{"code":"server_is_overloaded","message":"` + message + `"placeholderplaceholder`)
	err := newOpenAIUpstreamFailoverError(http.StatusBadRequest, nil, body, message, false)

	require.True(t, err.IsOpenAICapacityShed())
	require.Equal(t, http.StatusServiceUnavailable, err.ClientStatusCode)
	require.Equal(t, message, err.ClientMessage)
	require.NotContains(t, err.ClientMessage, "server_is_overloaded")
placeholder

func TestOpenAIStreamSemanticStatusesPreservedAcrossTerminalShapes(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		status       int
		wantFailover bool
placeholder{
		{"unauthorized", `{"type":"error","error":{"type":"authentication_error","code":"invalid_api_key","message":"unauthorized"placeholderplaceholder`, http.StatusUnauthorized, trueplaceholder,
		{"forbidden", `{"type":"response.failed","response":{"error":{"type":"permission_error","message":"forbidden"placeholderplaceholderplaceholder`, http.StatusForbidden, falseplaceholder,
		{"rate_limit", `{"error":{"type":"rate_limit_error","code":"rate_limit_exceeded","message":"slow down"placeholderplaceholder`, http.StatusTooManyRequests, trueplaceholder,
		{"overload_529", `{"type":"error","error":{"status_code":529,"code":"overloaded","message":"overloaded"placeholderplaceholder`, 529, trueplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := []byte(tt.body)
			message := extractOpenAISSEErrorMessage(payload)
			require.Equal(t, tt.status, openAIStreamFailureStatus(payload, message))
			require.Equal(t, tt.wantFailover, openAIStreamErrorEventShouldFailover(payload, message))
	placeholder)
placeholder
placeholder

func TestOpenAIStreamBareErrorUsesSemanticFailover(t *testing.T) {
	payload := []byte(`{"error":{"type":"rate_limit_error","code":"rate_limit_exceeded","message":"slow down"placeholderplaceholder`)
	require.True(t, openAIStreamErrorEventShouldFailover(payload, "slow down"))
placeholder

func TestOpenAIStream403FailoverRequiresStructuredAccountCredentialSignal(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    bool
placeholder{
		{
			name:    "ordinary permission error",
			payload: `{"type":"response.failed","response":{"error":{"type":"permission_error","code":"forbidden","message":"access denied for this request"placeholderplaceholderplaceholder`,
	placeholder,
		{
			name:    "explicit 403 request status",
			payload: `{"type":"error","error":{"type":"permission_error","code":"forbidden","status_code":403,"message":"forbidden content"placeholderplaceholder`,
	placeholder,
		{
			name:    "structured access state",
			payload: `{"type":"response.failed","response":{"error":{"code":"workspace_suspended","message":"workspace is suspended"placeholderplaceholderplaceholder`,
			want:    true,
	placeholder,
		{
			name:    "explicit credential auth code",
			payload: `{"type":"error","error":{"type":"permission_error","code":"invalid_api_key","status_code":403,"message":"credential rejected"placeholderplaceholder`,
			want:    true,
	placeholder,
		{
			name:    "explicit authentication type",
			payload: `{"type":"error","error":{"type":"authentication_error","status_code":403,"message":"credential rejected"placeholderplaceholder`,
			want:    true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := []byte(tt.payload)
			message := extractOpenAISSEErrorMessage(payload)
			require.Equal(t, http.StatusForbidden, openAIStreamFailureStatus(payload, message))
			require.Equal(t, tt.want, openAIStreamFailedEventShouldFailover(payload, message))
			require.Equal(t, tt.want, openAIStreamErrorEventShouldFailover(payload, message))
	placeholder)
placeholder
placeholder

func TestOpenAIStream403PostOutputAccountSideEffectsIgnoreRequestPermissionErrors(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	rateLimits := &RateLimitService{accountRepo: repoplaceholder
	svc := &OpenAIGatewayService{rateLimitService: rateLimitsplaceholder
	account := &Account{ID: 918, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	payload := []byte(`{"type":"error","error":{"type":"permission_error","code":"forbidden","status_code":403,"message":"access denied for this request"placeholderplaceholder`)

	status, disabled := svc.handleOpenAIStreamTerminalAccountSideEffects(nil, account, payload, "access denied for this request", nil)

	require.Equal(t, http.StatusForbidden, status)
	require.False(t, disabled)
	require.Zero(t, repo.setErrorCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIStream403ExplicitCredentialAuthAppliesAccountSideEffects(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	rateLimits := &RateLimitService{accountRepo: repoplaceholder
	svc := &OpenAIGatewayService{rateLimitService: rateLimitsplaceholder
	account := &Account{ID: 917, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	payload := []byte(`{"type":"error","error":{"type":"permission_error","code":"invalid_api_key","status_code":403,"message":"credential rejected"placeholderplaceholder`)

	status, disabled := svc.handleOpenAIStreamTerminalAccountSideEffects(nil, account, payload, "credential rejected", nil)

	require.Equal(t, http.StatusForbidden, status)
	require.True(t, disabled)
	require.Equal(t, 1, repo.setErrorCalls)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIWSStandaloneFailedStructured403AppliesAccountSideEffectsOnce(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	svc := &OpenAIGatewayService{rateLimitService: &RateLimitService{accountRepo: repoplaceholderplaceholder
	account := &Account{ID: 923, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	failed := []byte(`{"type":"response.failed","response":{"error":{"type":"permission_error","code":"invalid_api_key","status_code":403,"message":"credential rejected"placeholderplaceholderplaceholder`)

	require.True(t, svc.handleOpenAIWSFailureAccountSideEffects(context.Background(), account, "gpt-5", nil, failed))
	require.Equal(t, 1, repo.setErrorCalls)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIWSPairedStructured403SideEffectsCanBeDeduplicated(t *testing.T) {
	repo := &openAIStream403AccountRepo{placeholder
	svc := &OpenAIGatewayService{rateLimitService: &RateLimitService{accountRepo: repoplaceholderplaceholder
	account := &Account{ID: 924, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	errorEvent := []byte(`{"type":"error","error":{"code":"workspace_suspended","status_code":403,"message":"workspace is suspended"placeholderplaceholder`)
	failedEvent := []byte(`{"type":"response.failed","response":{"error":{"code":"workspace_suspended","status_code":403,"message":"workspace is suspended"placeholderplaceholderplaceholder`)

	applied := svc.handleOpenAIWSFailureAccountSideEffects(context.Background(), account, "gpt-5", nil, errorEvent)
	if !applied {
		applied = svc.handleOpenAIWSFailureAccountSideEffects(context.Background(), account, "gpt-5", nil, failedEvent)
placeholder

	require.True(t, applied)
	require.Equal(t, 1, repo.setErrorCalls)
placeholder

func TestOpenAIStreamAccessStateAppliesAccountHealthBeforeFailover(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	account := &Account{ID: 919, Platform: PlatformOpenAI, Type: AccountTypeSetupTokenplaceholder
	payload := []byte(`{"type":"response.failed","response":{"error":{"code":"workspace_suspended","message":"workspace is suspended"placeholderplaceholderplaceholder`)

	status, disabled := svc.handleOpenAIStreamTerminalAccountSideEffects(nil, account, payload, "workspace is suspended", nil)

	require.Equal(t, http.StatusForbidden, status)
	require.True(t, disabled)
	require.True(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestOpenAIStreamPairedFailureAppliesAccountSideEffectsOnce(t *testing.T) {
	const upstream = "data: {\"type\":\"response.output_text.delta\",\"delta\":\"partial\"placeholder\n\n" +
		"data: {\"type\":\"error\",\"error\":{\"status_code\":403,\"code\":\"workspace_suspended\",\"message\":\"workspace is suspended\"placeholderplaceholder\n\n" +
		"data: {\"type\":\"response.failed\",\"response\":{\"id\":\"resp_failed\",\"error\":{\"status_code\":403,\"code\":\"workspace_suspended\",\"message\":\"workspace is suspended\"placeholderplaceholderplaceholder\n\n"

	t.Run("native", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		repo := &openAIStream403AccountRepo{placeholder
		svc := &OpenAIGatewayService{
			cfg:              &config.Config{placeholder,
			toolCorrector:    NewCodexToolCorrector(),
			rateLimitService: &RateLimitService{accountRepo: repoplaceholder,
	placeholder
		account := &Account{ID: 921, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
		recorder := newOpenAIResponseFlushRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstream)),
	placeholder

		result, err := svc.handleStreamingResponse(context.Background(), resp, c, account, time.Now(), "gpt-5", "gpt-5")

	placeholder
		require.NotNil(t, result)
		require.Equal(t, 1, repo.setErrorCalls)
placeholder)

	t.Run("passthrough", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		repo := &openAIStream403AccountRepo{placeholder
		svc := &OpenAIGatewayService{
			cfg:              &config.Config{placeholder,
			rateLimitService: &RateLimitService{accountRepo: repoplaceholder,
	placeholder
		account := &Account{ID: 922, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
		writer := &passthroughFlushTestWriter{
			ResponseWriter:  c.Writer,
			recorder:        recorder,
			failAfterWrites: -1,
	placeholder
		c.Writer = writer
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
			Body:       io.NopCloser(strings.NewReader(upstream)),
	placeholder

		result, err := svc.handleStreamingResponsePassthrough(
			context.Background(), resp, c, account, time.Now(), "gpt-5", "gpt-5",
		)

	placeholder
		require.NotNil(t, result)
		require.Equal(t, 1, repo.setErrorCalls)
placeholder)
placeholder

func TestOpenAIStreamOAuthLike429GetsDeadlineWithoutImmediateRuntimeBlock(t *testing.T) {
	for _, accountType := range []string{AccountTypeOAuth, AccountTypeSetupTokenplaceholder {
		t.Run(accountType, func(t *testing.T) {
			svc := &OpenAIGatewayService{placeholder
			account := &Account{ID: 920, Platform: PlatformOpenAI, Type: accountTypeplaceholder
			payload := []byte(`{"type":"error","error":{"type":"rate_limit_error","code":"rate_limit_exceeded","message":"slow down"placeholderplaceholder`)
			status, disabled := svc.handleOpenAIStreamTerminalAccountSideEffects(nil, account, payload, "slow down", nil)
			err := svc.newOpenAIAccountFailoverError(account, status, nil, payload, "slow down", disabled, false)

			require.Equal(t, http.StatusTooManyRequests, status)
			require.False(t, disabled)
			require.True(t, err.RetryableOnSameAccount)
			require.False(t, err.SameAccountRetryDeadline.IsZero())
			require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
	placeholder)
placeholder
placeholder
