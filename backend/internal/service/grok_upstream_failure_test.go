//go:build unit

package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClassifyGrokUpstreamFailure_FreeUsage(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
placeholder{
		{
			name:   "code free-usage-exhausted",
			status: http.StatusTooManyRequests,
			body:   `{"error":{"code":"subscription:free-usage-exhausted","message":"You've used all the included free usage for model grok-4.5. Usage resets over a rolling 24-hour window."placeholderplaceholder`,
	placeholder,
		{
			name:   "chinese body without 429",
			status: http.StatusBadRequest,
			body:   `{"error":{"message":"模型额度用完，请稍后再试"placeholderplaceholder`,
	placeholder,
		{
			name:   "token pair with free marker",
			status: http.StatusOK,
			body:   `{"error":{"message":"free usage tokens (actual / limit): 2000000 / 2000000 for model grok-4.5"placeholderplaceholder`,
	placeholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := classifyGrokUpstreamFailure(tc.status, []byte(tc.body), "grok-4.5")
			require.Equal(t, GrokFailureFreeUsage, d.Class)
			require.True(t, d.ShouldCooldown)
			require.True(t, d.ShouldFailover)
			require.False(t, d.BlockModel, "free-usage must not soft-block models")
			require.GreaterOrEqual(t, d.Cooldown, 20*time.Minute)
	placeholder)
placeholder
placeholder

func TestClassifyGrokUpstreamFailure_EmptyUpstream(t *testing.T) {
	d := classifyGrokUpstreamFailure(http.StatusBadGateway, []byte(`empty model output: no content/tool_calls`), "grok-4.5")
	require.Equal(t, GrokFailureEmptyUpstream, d.Class)
	require.True(t, d.ShouldCooldown)
	require.True(t, d.ShouldFailover)
	require.True(t, d.BlockModel)
	require.Equal(t, 4*time.Minute, d.Cooldown)
placeholder

func TestClassifyGrokUpstreamFailure_Billing(t *testing.T) {
	d := classifyGrokUpstreamFailure(http.StatusForbidden, []byte(`{"code":"personal-team-blocked:spending-limit","error":"spending limit reached"placeholder`), "")
	require.Equal(t, GrokFailureBilling, d.Class)
	require.True(t, d.ShouldCooldown)
	require.True(t, d.ShouldFailover)
placeholder

func TestClassifyGrokUpstreamFailure_ValidationNoCool(t *testing.T) {
	d := classifyGrokUpstreamFailure(http.StatusBadRequest, []byte(`{"error":{"message":"invalid tool schema"placeholderplaceholder`), "")
	require.Equal(t, GrokFailureNone, d.Class)
	require.False(t, d.ShouldCooldown)
	require.False(t, d.ShouldFailover)
placeholder

func TestClassifyGrokUpstreamFailure_FreeUsageWinsOver5xx(t *testing.T) {
	// Proxy may rewrite free-usage into synthetic 502; body must win.
	d := classifyGrokUpstreamFailure(http.StatusBadGateway, []byte(`subscription:free-usage-exhausted for model grok-4.3`), "grok-4.3")
	require.Equal(t, GrokFailureFreeUsage, d.Class)
	require.NotEqual(t, GrokFailureServer, d.Class)
placeholder

func TestShouldFailoverGrokUpstreamError_FreeUsageBody(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	body := []byte(`{"error":{"code":"subscription:free-usage-exhausted","message":"free usage exhausted"placeholderplaceholder`)
	require.True(t, svc.shouldFailoverGrokUpstreamError(http.StatusBadRequest, body))
placeholder

func TestShouldFailoverGrokUpstreamError_ContentPolicyStillNoFailover(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	body := []byte(`{"error":{"code":"new_sensitive","message":"text is sensitive"placeholderplaceholder`)
	require.False(t, svc.shouldFailoverGrokUpstreamError(http.StatusForbidden, body))
placeholder

func TestHandleGrokAccountUpstreamError_FreeUsageBodyCoolsAccount(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 9101, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	before := time.Now()
	body := []byte(`{"error":{"code":"subscription:free-usage-exhausted","message":"You've used all the included free usage. Usage resets over a rolling 24-hour window."placeholderplaceholder`)

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusBadRequest, nil, body)

	require.Equal(t, 1, repo.tempUnschedCalls)
	require.Equal(t, "grok free usage exhausted", repo.lastTempUnschedReason)
	// 24h rolling → 2h cool
	require.Greater(t, repo.lastTempUnschedUntil, before.Add(119*time.Minute))
	require.Less(t, repo.lastTempUnschedUntil, before.Add(121*time.Minute))
placeholder

func TestHandleGrokAccountUpstreamError_EmptyOutputCoolsAccount(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 9102, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	before := time.Now()

	svc.handleGrokAccountUpstreamError(
		context.Background(), account, http.StatusBadGateway, nil,
		[]byte(`empty model output: no content/tool_calls`),
	)

	require.Equal(t, 1, repo.tempUnschedCalls)
	require.Equal(t, "grok empty model output", repo.lastTempUnschedReason)
	require.WithinDuration(t, before.Add(4*time.Minute), repo.lastTempUnschedUntil, time.Second)
placeholder

func TestHandleGrokAccountUpstreamError_FreeUsageDoesNotCoolPoolMode(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{
		ID:       9103,
		Platform: PlatformGrok,
		Type:     AccountTypeAPIKey,
placeholder
			"pool_mode": true,
	placeholder,
placeholder
	body := []byte(`{"error":{"code":"subscription:free-usage-exhausted","message":"free usage exhausted"placeholderplaceholder`)

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusBadRequest, nil, body)

	require.Zero(t, repo.tempUnschedCalls)
	require.False(t, svc.isOpenAIAccountRuntimeBlocked(account))
placeholder

func TestHandleGrokAccountUpstreamError_ContentPolicyStillNoMutation(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 9104, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	body := []byte(`{"error":{"code":"new_sensitive","message":"text is sensitive"placeholderplaceholder`)

	svc.handleGrokAccountUpstreamError(context.Background(), account, http.StatusForbidden, nil, body)

	require.Zero(t, repo.tempUnschedCalls)
placeholder

func TestHandleGrokAccountUpstreamError_Entitlement403Unchanged(t *testing.T) {
	repo := &grokQuotaAccountRepo{placeholder
	svc := &OpenAIGatewayService{accountRepo: repoplaceholder
	account := &Account{ID: 9105, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	before := time.Now()

	svc.handleGrokAccountUpstreamError(
		context.Background(), account, http.StatusForbidden, nil,
		[]byte(`{"error":{"message":"subscription required"placeholderplaceholder`),
	)

	require.Equal(t, 1, repo.tempUnschedCalls)
	require.Equal(t, "grok access or entitlement denied", repo.lastTempUnschedReason)
	require.Greater(t, repo.lastTempUnschedUntil, before.Add(29*time.Minute))
	require.Less(t, repo.lastTempUnschedUntil, before.Add(31*time.Minute))
placeholder
