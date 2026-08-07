//go:build unit

package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestForwardGrokResponses_PropagatesSearchCountFromJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"grok","input":"search something","tools":[{"type":"web_search"placeholder],"stream":falseplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

	account := healthyGrokOAuthGatewayTestAccount(9901, "access-token")
	repo := &mockAccountRepoForPlatform{accountsByID: map[int64]*Account{account.ID: accountplaceholderplaceholder
	upstreamBody := `{
		"id":"resp_search_bill",
		"object":"response",
		"model":"grok-4.5",
		"status":"completed",
		"output":[
			{"type":"web_search_call","id":"ws1","call_id":"c1","status":"completed"placeholder,
			{"type":"x_search_call","id":"xs1","call_id":"c2"placeholder,
			{"type":"message","role":"assistant","content":[{"type":"output_text","text":"ok"placeholder]placeholder
		],
		"usage":{"input_tokens":10,"output_tokens":5placeholder
placeholder`
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader([]byte(upstreamBody))),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", false, time.Now())
placeholder
	require.NotNil(t, result)
	require.Equal(t, 2, result.SearchCount, "Grok Responses must surface search tool calls for surcharge billing")
	require.Equal(t, 10, result.Usage.InputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
placeholder

func TestForwardGrokResponses_PropagatesSearchCountFromSSE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"model":"grok","input":"search","tools":[{"type":"web_search"placeholder],"stream":trueplaceholder`)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", bytes.NewReader(body))

	account := healthyGrokOAuthGatewayTestAccount(9902, "access-token")
	repo := &mockAccountRepoForPlatform{accountsByID: map[int64]*Account{account.ID: accountplaceholderplaceholder
	// item.done + response.completed for same call_id must count once after wire-up.
	sse := "data: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"web_search_call\",\"id\":\"ws1\",\"call_id\":\"c1\"placeholderplaceholder\n\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_s\",\"status\":\"completed\",\"output\":[{\"type\":\"web_search_call\",\"id\":\"ws1\",\"call_id\":\"c1\"placeholder],\"usage\":{\"input_tokens\":3,\"output_tokens\":1placeholderplaceholderplaceholder\n\n"
	upstream := &httpUpstreamRecorder{resp: &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       io.NopCloser(bytes.NewReader([]byte(sse))),
placeholderplaceholder
	svc := &OpenAIGatewayService{
		httpUpstream:      upstream,
		grokTokenProvider: NewGrokTokenProvider(repo, nil),
		accountRepo:       repo,
placeholder

	result, err := svc.forwardGrokResponses(context.Background(), c, account, body, "grok", true, time.Now())
placeholder
	require.NotNil(t, result)
	require.Equal(t, 1, result.SearchCount, "stream SearchCount must be wired and deduped")
placeholder

func TestGetSchedulableAccount_AppliesGrokFreeSoftGate(t *testing.T) {
	// Sticky/non-list path must not return over-gate free OAuth accounts.
	cfg := &config.Config{placeholder
	cfg.Gateway.Grok.FreeQuotaSoftGateEnabled = true
	cfg.Gateway.Grok.FreeQuotaTokenLimit = 1_000_000
	cfg.Gateway.Grok.FreeQuotaSoftGatePercent = 95
	cfg.Gateway.Grok.FreeQuotaWindowHours = 24
	cfg.Gateway.Grok.FreeQuotaStatsCacheSeconds = 0

	account := healthyGrokOAuthGatewayTestAccount(8801, "tok")
	account.Credentials["subscription_tier"] = "free"
	account.Status = StatusActive
	account.Schedulable = true

	repo := &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholder
	usageRepo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		account.ID: {Tokens: 999_000placeholder,
placeholderplaceholder
	// Clear shared gateway free-gate cache so this test is deterministic.
	gatewayGrokFreeQuotaGateCache.Range(func(key, _ any) bool {
		gatewayGrokFreeQuotaGateCache.Delete(key)
		return true
placeholder)
	svc := &GatewayService{
		cfg:          cfg,
		accountRepo:  repo,
		usageLogRepo: usageRepo,
placeholder

	got, err := svc.getSchedulableAccount(context.Background(), account.ID)
placeholder
	require.Nil(t, got, "over free soft-gate sticky hit must miss")
placeholder

func TestOpenAIGetSchedulableAccount_AppliesGrokFreeSoftGate(t *testing.T) {
	// Legacy OpenAI-compatible sticky (advanced scheduler off) must free-gate Grok.
	cfg := &config.Config{placeholder
	cfg.Gateway.Grok.FreeQuotaSoftGateEnabled = true
	cfg.Gateway.Grok.FreeQuotaTokenLimit = 1_000_000
	cfg.Gateway.Grok.FreeQuotaSoftGatePercent = 95
	cfg.Gateway.Grok.FreeQuotaWindowHours = 24
	cfg.Gateway.Grok.FreeQuotaStatsCacheSeconds = 0

	account := healthyGrokOAuthGatewayTestAccount(8802, "tok")
	account.Credentials["subscription_tier"] = "free"
	account.Status = StatusActive
	account.Schedulable = true

	repo := &mockAccountRepoForPlatform{
		accountsByID: map[int64]*Account{account.ID: accountplaceholder,
placeholder
	usageRepo := &grokFreeQuotaUsageRepoStub{stats: map[int64]*usagestats.AccountStats{
		account.ID: {Tokens: 999_000placeholder,
placeholderplaceholder
	openaiGrokFreeQuotaGateCache.Range(func(key, _ any) bool {
		openaiGrokFreeQuotaGateCache.Delete(key)
		return true
placeholder)
	svc := &OpenAIGatewayService{
		cfg:          cfg,
		accountRepo:  repo,
		usageLogRepo: usageRepo,
placeholder

	got, err := svc.getSchedulableAccount(context.Background(), account.ID)
placeholder
	require.Nil(t, got, "OpenAI legacy sticky must apply free soft-gate")
placeholder

func TestCountGrokNativeSearchCallsFromJSON_MessagesStyleBody(t *testing.T) {
	// Proves the same counter used by Anthropic-buffered Grok /v1/messages path.
	body := []byte(`{"id":"r1","output":[{"type":"web_search_call","id":"ws1"placeholder,{"type":"message","role":"assistant"placeholder]placeholder`)
	require.Equal(t, 1, countGrokNativeSearchCallsFromJSONBytes(body))
placeholder
