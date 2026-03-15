//go:build unit

package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/stretchr/testify/require"
)

func TestClassifyAntigravity429(t *testing.T) {
	t.Run("明确配额耗尽", func(t *testing.T) {
		body := []byte(`{"error":{"status":"RESOURCE_EXHAUSTED","message":"QUOTA_EXHAUSTED"placeholderplaceholder`)
		require.Equal(t, antigravity429QuotaExhausted, classifyAntigravity429(body))
placeholder)

	t.Run("结构化限流", func(t *testing.T) {
		body := []byte(`{
			"error": {
				"status": "RESOURCE_EXHAUSTED",
				"details": [
					{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "metadata": {"model": "claude-sonnet-4-5"placeholder, "reason": "RATE_LIMIT_EXCEEDED"placeholder,
					{"@type": "type.googleapis.com/google.rpc.RetryInfo", "retryDelay": "0.5s"placeholder
				]
		placeholder
	placeholder`)
		require.Equal(t, antigravity429RateLimited, classifyAntigravity429(body))
placeholder)

	t.Run("未知429", func(t *testing.T) {
		body := []byte(`{"error":{"message":"too many requests"placeholderplaceholder`)
		require.Equal(t, antigravity429Unknown, classifyAntigravity429(body))
placeholder)
placeholder

func TestIsCreditsExhausted_UsesAICreditsKey(t *testing.T) {
	t.Run("无 AICredits key 则积分可用", func(t *testing.T) {
		account := &Account{
			ID:       1,
			Platform: PlatformAntigravity,
			Extra: map[string]any{
				"allow_overages": true,
		placeholder,
	placeholder
		require.False(t, account.isCreditsExhausted())
placeholder)

	t.Run("AICredits key 生效则积分耗尽", func(t *testing.T) {
		account := &Account{
			ID:       2,
			Platform: PlatformAntigravity,
			Extra: map[string]any{
				"allow_overages": true,
				modelRateLimitsKey: map[string]any{
					creditsExhaustedKey: map[string]any{
						"rate_limited_at":     time.Now().UTC().Format(time.RFC3339),
						"rate_limit_reset_at": time.Now().Add(5 * time.Hour).UTC().Format(time.RFC3339),
				placeholder,
			placeholder,
		placeholder,
	placeholder
		require.True(t, account.isCreditsExhausted())
placeholder)

	t.Run("AICredits key 过期则积分可用", func(t *testing.T) {
		account := &Account{
			ID:       3,
			Platform: PlatformAntigravity,
			Extra: map[string]any{
				"allow_overages": true,
				modelRateLimitsKey: map[string]any{
					creditsExhaustedKey: map[string]any{
						"rate_limited_at":     time.Now().Add(-6 * time.Hour).UTC().Format(time.RFC3339),
						"rate_limit_reset_at": time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339),
				placeholder,
			placeholder,
		placeholder,
	placeholder
		require.False(t, account.isCreditsExhausted())
placeholder)
placeholder

func TestHandleSmartRetry_QuotaExhausted_UsesCreditsAndStoresIndependentState(t *testing.T) {
	successResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(`{"ok":trueplaceholder`)),
placeholder
	upstream := &mockSmartRetryUpstream{
		responses: []*http.Response{successRespplaceholder,
		errors:    []error{nilplaceholder,
placeholder
	repo := &stubAntigravityAccountRepo{placeholder
	account := &Account{
		ID:       101,
		Name:     "acc-101",
		Type:     AccountTypeOAuth,
		Platform: PlatformAntigravity,
		Extra: map[string]any{
			"allow_overages": true,
	placeholder,
placeholder
			"model_mapping": map[string]any{
				"claude-opus-4-6": "claude-sonnet-4-5",
		placeholder,
	placeholder,
placeholder

	respBody := []byte(`{"error":{"status":"RESOURCE_EXHAUSTED","message":"QUOTA_EXHAUSTED"placeholderplaceholder`)
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
placeholder
	params := antigravityRetryLoopParams{
		ctx:            context.Background(),
		prefix:         "[test]",
		account:        account,
		accessToken:    "token",
		action:         "generateContent",
		body:           []byte(`{"model":"claude-opus-4-6","request":{placeholderplaceholder`),
		httpUpstream:   upstream,
		accountRepo:    repo,
		requestedModel: "claude-opus-4-6",
		handleError: func(ctx context.Context, prefix string, account *Account, statusCode int, headers http.Header, body []byte, requestedModel string, groupID int64, sessionHash string, isStickySession bool) *handleModelRateLimitResult {
			return nil
	placeholder,
placeholder

	svc := &AntigravityGatewayService{placeholder
	result := svc.handleSmartRetry(params, resp, respBody, "https://ag-1.test", 0, []string{"https://ag-1.test"placeholder)

	require.NotNil(t, result)
	require.Equal(t, smartRetryActionBreakWithResp, result.action)
	require.NotNil(t, result.resp)
	require.Nil(t, result.switchError)
	require.Len(t, upstream.requestBodies, 1)
	require.Contains(t, string(upstream.requestBodies[0]), "enabledCreditTypes")
	require.Empty(t, repo.modelRateLimitCalls, "overages 成功后不应写入普通 model_rate_limits")
placeholder

func TestHandleSmartRetry_RateLimited_DoesNotUseCredits(t *testing.T) {
	successResp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(strings.NewReader(`{"ok":trueplaceholder`)),
placeholder
	upstream := &mockSmartRetryUpstream{
		responses: []*http.Response{successRespplaceholder,
		errors:    []error{nilplaceholder,
placeholder
	repo := &stubAntigravityAccountRepo{placeholder
	account := &Account{
		ID:       102,
		Name:     "acc-102",
		Type:     AccountTypeOAuth,
		Platform: PlatformAntigravity,
		Extra: map[string]any{
			"allow_overages": true,
	placeholder,
placeholder

	respBody := []byte(`{
		"error": {
			"status": "RESOURCE_EXHAUSTED",
			"details": [
				{"@type": "type.googleapis.com/google.rpc.ErrorInfo", "metadata": {"model": "claude-sonnet-4-5"placeholder, "reason": "RATE_LIMIT_EXCEEDED"placeholder,
				{"@type": "type.googleapis.com/google.rpc.RetryInfo", "retryDelay": "0.1s"placeholder
			]
	placeholder
placeholder`)
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{placeholder,
		Body:       io.NopCloser(bytes.NewReader(respBody)),
placeholder
	params := antigravityRetryLoopParams{
		ctx:          context.Background(),
		prefix:       "[test]",
		account:      account,
		accessToken:  "token",
		action:       "generateContent",
		body:         []byte(`{"model":"claude-sonnet-4-5","request":{placeholderplaceholder`),
		httpUpstream: upstream,
		accountRepo:  repo,
		handleError: func(ctx context.Context, prefix string, account *Account, statusCode int, headers http.Header, body []byte, requestedModel string, groupID int64, sessionHash string, isStickySession bool) *handleModelRateLimitResult {
			return nil
	placeholder,
placeholder

	svc := &AntigravityGatewayService{placeholder
	result := svc.handleSmartRetry(params, resp, respBody, "https://ag-1.test", 0, []string{"https://ag-1.test"placeholder)

	require.NotNil(t, result)
	require.Equal(t, smartRetryActionBreakWithResp, result.action)
	require.NotNil(t, result.resp)
	require.Len(t, upstream.requestBodies, 1)
	require.NotContains(t, string(upstream.requestBodies[0]), "enabledCreditTypes")
	require.Empty(t, repo.extraUpdateCalls)
	require.Empty(t, repo.modelRateLimitCalls)
placeholder

func TestAntigravityRetryLoop_ModelRateLimited_InjectsCredits(t *testing.T) {
	oldBaseURLs := append([]string(nil), antigravity.BaseURLs...)
	oldAvailability := antigravity.DefaultURLAvailability
	defer func() {
		antigravity.BaseURLs = oldBaseURLs
		antigravity.DefaultURLAvailability = oldAvailability
placeholder()

	antigravity.BaseURLs = []string{"https://ag-1.test"placeholder
	antigravity.DefaultURLAvailability = antigravity.NewURLAvailability(time.Minute)

	upstream := &queuedHTTPUpstreamStub{
		responses: []*http.Response{
			{
				StatusCode: http.StatusOK,
				Header:     http.Header{placeholder,
				Body:       io.NopCloser(strings.NewReader(`{"ok":trueplaceholder`)),
		placeholder,
	placeholder,
		errors: []error{nilplaceholder,
placeholder
	// 模型已限流 + overages 启用 + 无 AICredits key → 应直接注入积分
	account := &Account{
		ID:          103,
		Name:        "acc-103",
		Type:        AccountTypeOAuth,
		Platform:    PlatformAntigravity,
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"allow_overages": true,
			modelRateLimitsKey: map[string]any{
				"claude-sonnet-4-5": map[string]any{
					"rate_limited_at":     time.Now().UTC().Format(time.RFC3339),
					"rate_limit_reset_at": time.Now().Add(30 * time.Minute).UTC().Format(time.RFC3339),
			placeholder,
		placeholder,
	placeholder,
placeholder

	svc := &AntigravityGatewayService{placeholder
	result, err := svc.antigravityRetryLoop(antigravityRetryLoopParams{
		ctx:            context.Background(),
		prefix:         "[test]",
		account:        account,
		accessToken:    "token",
		action:         "generateContent",
		body:           []byte(`{"model":"claude-sonnet-4-5","request":{placeholderplaceholder`),
		httpUpstream:   upstream,
		requestedModel: "claude-sonnet-4-5",
		handleError: func(ctx context.Context, prefix string, account *Account, statusCode int, headers http.Header, body []byte, requestedModel string, groupID int64, sessionHash string, isStickySession bool) *handleModelRateLimitResult {
			return nil
	placeholder,
placeholder)

placeholder
	require.NotNil(t, result)
	require.Len(t, upstream.requestBodies, 1)
	require.Contains(t, string(upstream.requestBodies[0]), "enabledCreditTypes")
placeholder

func TestAntigravityRetryLoop_CreditsExhausted_DoesNotInject(t *testing.T) {
	oldBaseURLs := append([]string(nil), antigravity.BaseURLs...)
	oldAvailability := antigravity.DefaultURLAvailability
	defer func() {
		antigravity.BaseURLs = oldBaseURLs
		antigravity.DefaultURLAvailability = oldAvailability
placeholder()

	antigravity.BaseURLs = []string{"https://ag-1.test"placeholder
	antigravity.DefaultURLAvailability = antigravity.NewURLAvailability(time.Minute)

	// 模型限流 + overages 启用 + AICredits key 生效 → 不应注入积分，应切号
	account := &Account{
		ID:          104,
		Name:        "acc-104",
		Type:        AccountTypeOAuth,
		Platform:    PlatformAntigravity,
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"allow_overages": true,
			modelRateLimitsKey: map[string]any{
				"claude-sonnet-4-5": map[string]any{
					"rate_limited_at":     time.Now().UTC().Format(time.RFC3339),
					"rate_limit_reset_at": time.Now().Add(30 * time.Minute).UTC().Format(time.RFC3339),
			placeholder,
				creditsExhaustedKey: map[string]any{
					"rate_limited_at":     time.Now().UTC().Format(time.RFC3339),
					"rate_limit_reset_at": time.Now().Add(5 * time.Hour).UTC().Format(time.RFC3339),
			placeholder,
		placeholder,
	placeholder,
placeholder

	svc := &AntigravityGatewayService{placeholder
	_, err := svc.antigravityRetryLoop(antigravityRetryLoopParams{
		ctx:            context.Background(),
		prefix:         "[test]",
		account:        account,
		accessToken:    "token",
		action:         "generateContent",
		body:           []byte(`{"model":"claude-sonnet-4-5","request":{placeholderplaceholder`),
		requestedModel: "claude-sonnet-4-5",
		handleError: func(ctx context.Context, prefix string, account *Account, statusCode int, headers http.Header, body []byte, requestedModel string, groupID int64, sessionHash string, isStickySession bool) *handleModelRateLimitResult {
			return nil
	placeholder,
placeholder)

	// 模型限流 + 积分耗尽 → 应触发切号错误
placeholder
	var switchErr *AntigravityAccountSwitchError
	require.ErrorAs(t, err, &switchErr)
placeholder

func TestAntigravityRetryLoop_CreditErrorMarksExhausted(t *testing.T) {
	oldBaseURLs := append([]string(nil), antigravity.BaseURLs...)
	oldAvailability := antigravity.DefaultURLAvailability
	defer func() {
		antigravity.BaseURLs = oldBaseURLs
		antigravity.DefaultURLAvailability = oldAvailability
placeholder()

	antigravity.BaseURLs = []string{"https://ag-1.test"placeholder
	antigravity.DefaultURLAvailability = antigravity.NewURLAvailability(time.Minute)

	repo := &stubAntigravityAccountRepo{placeholder
	upstream := &queuedHTTPUpstreamStub{
		responses: []*http.Response{
			{
				StatusCode: http.StatusForbidden,
				Header:     http.Header{placeholder,
				Body:       io.NopCloser(strings.NewReader(`{"error":{"message":"Insufficient GOOGLE_ONE_AI credits"placeholderplaceholder`)),
		placeholder,
	placeholder,
		errors: []error{nilplaceholder,
placeholder
	// 模型限流 + overages 启用 + 积分可用 → 注入积分但上游返回积分不足
	account := &Account{
		ID:          105,
		Name:        "acc-105",
		Type:        AccountTypeOAuth,
		Platform:    PlatformAntigravity,
		Status:      StatusActive,
		Schedulable: true,
		Extra: map[string]any{
			"allow_overages": true,
			modelRateLimitsKey: map[string]any{
				"claude-sonnet-4-5": map[string]any{
					"rate_limited_at":     time.Now().UTC().Format(time.RFC3339),
					"rate_limit_reset_at": time.Now().Add(30 * time.Minute).UTC().Format(time.RFC3339),
			placeholder,
		placeholder,
	placeholder,
placeholder

	svc := &AntigravityGatewayService{accountRepo: repoplaceholder
	result, err := svc.antigravityRetryLoop(antigravityRetryLoopParams{
		ctx:            context.Background(),
		prefix:         "[test]",
		account:        account,
		accessToken:    "token",
		action:         "generateContent",
		body:           []byte(`{"model":"claude-sonnet-4-5","request":{placeholderplaceholder`),
		httpUpstream:   upstream,
		accountRepo:    repo,
		requestedModel: "claude-sonnet-4-5",
		handleError: func(ctx context.Context, prefix string, account *Account, statusCode int, headers http.Header, body []byte, requestedModel string, groupID int64, sessionHash string, isStickySession bool) *handleModelRateLimitResult {
			return nil
	placeholder,
placeholder)

placeholder
	require.NotNil(t, result)
	// 验证 AICredits key 已通过 SetModelRateLimit 写入数据库
	require.Len(t, repo.modelRateLimitCalls, 1, "应通过 SetModelRateLimit 写入 AICredits key")
	require.Equal(t, creditsExhaustedKey, repo.modelRateLimitCalls[0].modelKey)
placeholder

func TestShouldMarkCreditsExhausted(t *testing.T) {
	t.Run("reqErr 不为 nil 时不标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusForbiddenplaceholder
		require.False(t, shouldMarkCreditsExhausted(resp, []byte(`{"error":"Insufficient credits"placeholder`), io.ErrUnexpectedEOF))
placeholder)

	t.Run("resp 为 nil 时不标记", func(t *testing.T) {
		require.False(t, shouldMarkCreditsExhausted(nil, []byte(`{"error":"Insufficient credits"placeholder`), nil))
placeholder)

	t.Run("5xx 响应不标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusInternalServerErrorplaceholder
		require.False(t, shouldMarkCreditsExhausted(resp, []byte(`{"error":"Insufficient credits"placeholder`), nil))
placeholder)

	t.Run("408 RequestTimeout 不标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusRequestTimeoutplaceholder
		require.False(t, shouldMarkCreditsExhausted(resp, []byte(`{"error":"Insufficient credits"placeholder`), nil))
placeholder)

	t.Run("URL 级限流不标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusTooManyRequestsplaceholder
		body := []byte(`{"error":{"message":"Resource has been exhausted"placeholderplaceholder`)
		require.False(t, shouldMarkCreditsExhausted(resp, body, nil))
placeholder)

	t.Run("结构化限流不标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusTooManyRequestsplaceholder
		body := []byte(`{"error":{"status":"RESOURCE_EXHAUSTED","details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"RATE_LIMIT_EXCEEDED"placeholder,{"@type":"type.googleapis.com/google.rpc.RetryInfo","retryDelay":"0.5s"placeholder]placeholderplaceholder`)
		require.False(t, shouldMarkCreditsExhausted(resp, body, nil))
placeholder)

	t.Run("含 credits 关键词时标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusForbiddenplaceholder
		for _, keyword := range []string{
			"Insufficient GOOGLE_ONE_AI credits",
			"insufficient credit balance",
			"not enough credits for this request",
			"Credits exhausted",
			"minimumCreditAmountForUsage requirement not met",
	placeholder {
			body := []byte(`{"error":{"message":"` + keyword + `"placeholderplaceholder`)
			require.True(t, shouldMarkCreditsExhausted(resp, body, nil), "should mark for keyword: %s", keyword)
	placeholder
placeholder)

	t.Run("无 credits 关键词时不标记", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusForbiddenplaceholder
		body := []byte(`{"error":{"message":"permission denied"placeholderplaceholder`)
		require.False(t, shouldMarkCreditsExhausted(resp, body, nil))
placeholder)
placeholder

func TestInjectEnabledCreditTypes(t *testing.T) {
	t.Run("正常 JSON 注入成功", func(t *testing.T) {
		body := []byte(`{"model":"claude-sonnet-4-5","request":{placeholderplaceholder`)
		result := injectEnabledCreditTypes(body)
		require.NotNil(t, result)
		require.Contains(t, string(result), `"enabledCreditTypes"`)
		require.Contains(t, string(result), `GOOGLE_ONE_AI`)
placeholder)

	t.Run("非法 JSON 返回 nil", func(t *testing.T) {
		require.Nil(t, injectEnabledCreditTypes([]byte(`not json`)))
placeholder)

	t.Run("空 body 返回 nil", func(t *testing.T) {
		require.Nil(t, injectEnabledCreditTypes([]byte{placeholder))
placeholder)

	t.Run("已有 enabledCreditTypes 会被覆盖", func(t *testing.T) {
		body := []byte(`{"enabledCreditTypes":["OLD"],"model":"test"placeholder`)
		result := injectEnabledCreditTypes(body)
		require.NotNil(t, result)
		require.Contains(t, string(result), `GOOGLE_ONE_AI`)
		require.NotContains(t, string(result), `OLD`)
placeholder)
placeholder

func TestClearCreditsExhausted(t *testing.T) {
	t.Run("account 为 nil 不操作", func(t *testing.T) {
		repo := &stubAntigravityAccountRepo{placeholder
		svc := &AntigravityGatewayService{accountRepo: repoplaceholder
		svc.clearCreditsExhausted(context.Background(), nil)
		require.Empty(t, repo.extraUpdateCalls)
placeholder)

	t.Run("Extra 为 nil 不操作", func(t *testing.T) {
		repo := &stubAntigravityAccountRepo{placeholder
		svc := &AntigravityGatewayService{accountRepo: repoplaceholder
		svc.clearCreditsExhausted(context.Background(), &Account{ID: 1placeholder)
		require.Empty(t, repo.extraUpdateCalls)
placeholder)

	t.Run("无 modelRateLimitsKey 不操作", func(t *testing.T) {
		repo := &stubAntigravityAccountRepo{placeholder
		svc := &AntigravityGatewayService{accountRepo: repoplaceholder
		svc.clearCreditsExhausted(context.Background(), &Account{
			ID:    1,
			Extra: map[string]any{"some_key": "value"placeholder,
	placeholder)
		require.Empty(t, repo.extraUpdateCalls)
placeholder)

	t.Run("无 AICredits key 不操作", func(t *testing.T) {
		repo := &stubAntigravityAccountRepo{placeholder
		svc := &AntigravityGatewayService{accountRepo: repoplaceholder
		svc.clearCreditsExhausted(context.Background(), &Account{
			ID: 1,
			Extra: map[string]any{
				modelRateLimitsKey: map[string]any{
					"claude-sonnet-4-5": map[string]any{
						"rate_limited_at":     "2026-03-15T00:00:00Z",
						"rate_limit_reset_at": "2099-03-15T00:00:00Z",
				placeholder,
			placeholder,
		placeholder,
	placeholder)
		require.Empty(t, repo.extraUpdateCalls)
placeholder)

	t.Run("有 AICredits key 时删除并调用 UpdateExtra", func(t *testing.T) {
		repo := &stubAntigravityAccountRepo{placeholder
		svc := &AntigravityGatewayService{accountRepo: repoplaceholder
		account := &Account{
			ID: 1,
			Extra: map[string]any{
				modelRateLimitsKey: map[string]any{
					"claude-sonnet-4-5": map[string]any{
						"rate_limited_at":     "2026-03-15T00:00:00Z",
						"rate_limit_reset_at": "2099-03-15T00:00:00Z",
				placeholder,
					creditsExhaustedKey: map[string]any{
						"rate_limited_at":     "2026-03-15T00:00:00Z",
						"rate_limit_reset_at": time.Now().Add(5 * time.Hour).UTC().Format(time.RFC3339),
				placeholder,
			placeholder,
		placeholder,
	placeholder
		svc.clearCreditsExhausted(context.Background(), account)
		require.Len(t, repo.extraUpdateCalls, 1)
		// AICredits key 应被删除
		rawLimits := account.Extra[modelRateLimitsKey].(map[string]any)
		_, exists := rawLimits[creditsExhaustedKey]
		require.False(t, exists, "AICredits key 应被删除")
		// 普通模型限流应保留
		_, exists = rawLimits["claude-sonnet-4-5"]
		require.True(t, exists, "普通模型限流应保留")
placeholder)
placeholder
