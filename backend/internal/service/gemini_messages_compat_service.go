package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	mathrand "math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/pkg/googleapi"

	"github.com/gin-gonic/gin"
)

const geminiStickySessionTTL = time.Hour

const (
	geminiMaxRetries     = 5
	geminiRetryBaseDelay = 1 * time.Second
	geminiRetryMaxDelay  = 16 * time.Second
)

type GeminiMessagesCompatService struct {
	accountRepo               AccountRepository
	groupRepo                 GroupRepository
	cache                     GatewayCache
	tokenProvider             *GeminiTokenProvider
	rateLimitService          *RateLimitService
	httpUpstream              HTTPUpstream
	antigravityGatewayService *AntigravityGatewayService
placeholder

func NewGeminiMessagesCompatService(
	accountRepo AccountRepository,
	groupRepo GroupRepository,
	cache GatewayCache,
	tokenProvider *GeminiTokenProvider,
	rateLimitService *RateLimitService,
	httpUpstream HTTPUpstream,
	antigravityGatewayService *AntigravityGatewayService,
) *GeminiMessagesCompatService {
	return &GeminiMessagesCompatService{
		accountRepo:               accountRepo,
		groupRepo:                 groupRepo,
		cache:                     cache,
		tokenProvider:             tokenProvider,
		rateLimitService:          rateLimitService,
		httpUpstream:              httpUpstream,
		antigravityGatewayService: antigravityGatewayService,
placeholder
placeholder

// GetTokenProvider returns the token provider for OAuth accounts
func (s *GeminiMessagesCompatService) GetTokenProvider() *GeminiTokenProvider {
	return s.tokenProvider
placeholder

func (s *GeminiMessagesCompatService) SelectAccountForModel(ctx context.Context, groupID *int64, sessionHash string, requestedModel string) (*Account, error) {
	return s.SelectAccountForModelWithExclusions(ctx, groupID, sessionHash, requestedModel, nil)
placeholder

func (s *GeminiMessagesCompatService) SelectAccountForModelWithExclusions(ctx context.Context, groupID *int64, sessionHash string, requestedModel string, excludedIDs map[int64]struct{placeholder) (*Account, error) {
	// 优先检查 context 中的强制平台（/antigravity 路由）
	var platform string
	forcePlatform, hasForcePlatform := ctx.Value(ctxkey.ForcePlatform).(string)
	if hasForcePlatform && forcePlatform != "" {
		platform = forcePlatform
placeholder else if groupID != nil {
		// 根据分组 platform 决定查询哪种账号
		group, err := s.groupRepo.GetByID(ctx, *groupID)
		if err != nil {
			return nil, fmt.Errorf("get group failed: %w", err)
	placeholder
		platform = group.Platform
placeholder else {
		// 无分组时只使用原生 gemini 平台
		platform = PlatformGemini
placeholder

	// gemini 分组支持混合调度（包含启用了 mixed_scheduling 的 antigravity 账户）
	// 注意：强制平台模式不走混合调度
	useMixedScheduling := platform == PlatformGemini && !hasForcePlatform
	var queryPlatforms []string
	if useMixedScheduling {
		queryPlatforms = []string{PlatformGemini, PlatformAntigravityplaceholder
placeholder else {
		queryPlatforms = []string{platformplaceholder
placeholder

	cacheKey := "gemini:" + sessionHash

	if sessionHash != "" {
		accountID, err := s.cache.GetSessionAccountID(ctx, cacheKey)
		if err == nil && accountID > 0 {
			if _, excluded := excludedIDs[accountID]; !excluded {
				account, err := s.accountRepo.GetByID(ctx, accountID)
				// 检查账号是否有效：原生平台直接匹配，antigravity 需要启用混合调度
				if err == nil && account.IsSchedulable() && (requestedModel == "" || s.isModelSupportedByAccount(account, requestedModel)) {
					valid := false
					if account.Platform == platform {
						valid = true
				placeholder else if useMixedScheduling && account.Platform == PlatformAntigravity && account.IsMixedSchedulingEnabled() {
						valid = true
				placeholder
					if valid {
						usable := true
						if s.rateLimitService != nil && requestedModel != "" {
							ok, err := s.rateLimitService.PreCheckUsage(ctx, account, requestedModel)
							if err != nil {
								log.Printf("[Gemini PreCheck] Account %d precheck error: %v", account.ID, err)
						placeholder
							if !ok {
								usable = false
						placeholder
					placeholder
						if usable {
							_ = s.cache.RefreshSessionTTL(ctx, cacheKey, geminiStickySessionTTL)
							return account, nil
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// 查询可调度账户（强制平台模式：优先按分组查找，找不到再查全部）
	var accounts []Account
	var err error
	if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatforms(ctx, *groupID, queryPlatforms)
		if err != nil {
			return nil, fmt.Errorf("query accounts failed: %w", err)
	placeholder
		// 强制平台模式下，分组中找不到账户时回退查询全部
		if len(accounts) == 0 && hasForcePlatform {
			accounts, err = s.accountRepo.ListSchedulableByPlatforms(ctx, queryPlatforms)
	placeholder
placeholder else {
		accounts, err = s.accountRepo.ListSchedulableByPlatforms(ctx, queryPlatforms)
placeholder
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder

	var selected *Account
	for i := range accounts {
		acc := &accounts[i]
		if _, excluded := excludedIDs[acc.ID]; excluded {
			continue
	placeholder
		// 混合调度模式下：原生平台直接通过，antigravity 需要启用 mixed_scheduling
		// 非混合调度模式（antigravity 分组）：不需要过滤
		if useMixedScheduling && acc.Platform == PlatformAntigravity && !acc.IsMixedSchedulingEnabled() {
			continue
	placeholder
		if requestedModel != "" && !s.isModelSupportedByAccount(acc, requestedModel) {
			continue
	placeholder
		if s.rateLimitService != nil && requestedModel != "" {
			ok, err := s.rateLimitService.PreCheckUsage(ctx, acc, requestedModel)
			if err != nil {
				log.Printf("[Gemini PreCheck] Account %d precheck error: %v", acc.ID, err)
		placeholder
			if !ok {
				continue
		placeholder
	placeholder
		if selected == nil {
			selected = acc
			continue
	placeholder
		if acc.Priority < selected.Priority {
			selected = acc
	placeholder else if acc.Priority == selected.Priority {
			switch {
			case acc.LastUsedAt == nil && selected.LastUsedAt != nil:
				selected = acc
			case acc.LastUsedAt != nil && selected.LastUsedAt == nil:
				// keep selected (never used is preferred)
			case acc.LastUsedAt == nil && selected.LastUsedAt == nil:
				// Prefer OAuth accounts when both are unused (more compatible for Code Assist flows).
				if acc.Type == AccountTypeOAuth && selected.Type != AccountTypeOAuth {
					selected = acc
			placeholder
			default:
				if acc.LastUsedAt.Before(*selected.LastUsedAt) {
					selected = acc
			placeholder
		placeholder
	placeholder
placeholder

	if selected == nil {
		if requestedModel != "" {
			return nil, fmt.Errorf("no available Gemini accounts supporting model: %s", requestedModel)
	placeholder
		return nil, errors.New("no available Gemini accounts")
placeholder

	if sessionHash != "" {
		_ = s.cache.SetSessionAccountID(ctx, cacheKey, selected.ID, geminiStickySessionTTL)
placeholder

	return selected, nil
placeholder

// isModelSupportedByAccount 根据账户平台检查模型支持
func (s *GeminiMessagesCompatService) isModelSupportedByAccount(account *Account, requestedModel string) bool {
	if account.Platform == PlatformAntigravity {
		return IsAntigravityModelSupported(requestedModel)
placeholder
	return account.IsModelSupported(requestedModel)
placeholder

// GetAntigravityGatewayService 返回 AntigravityGatewayService
func (s *GeminiMessagesCompatService) GetAntigravityGatewayService() *AntigravityGatewayService {
	return s.antigravityGatewayService
placeholder

// HasAntigravityAccounts 检查是否有可用的 antigravity 账户
func (s *GeminiMessagesCompatService) HasAntigravityAccounts(ctx context.Context, groupID *int64) (bool, error) {
	var accounts []Account
	var err error
	if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, PlatformAntigravity)
placeholder else {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, PlatformAntigravity)
placeholder
	if err != nil {
		return false, err
placeholder
	return len(accounts) > 0, nil
placeholder

// SelectAccountForAIStudioEndpoints selects an account that is likely to succeed against
// generativelanguage.googleapis.com (e.g. GET /v1beta/models).
//
// Preference order:
// 1) API key accounts (AI Studio)
// 2) OAuth accounts without project_id (AI Studio OAuth)
// 3) OAuth accounts explicitly marked as ai_studio
// 4) Any remaining Gemini accounts (fallback)
func (s *GeminiMessagesCompatService) SelectAccountForAIStudioEndpoints(ctx context.Context, groupID *int64) (*Account, error) {
	var accounts []Account
	var err error
	if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, PlatformGemini)
placeholder else {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, PlatformGemini)
placeholder
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder
	if len(accounts) == 0 {
		return nil, errors.New("no available Gemini accounts")
placeholder

	rank := func(a *Account) int {
		if a == nil {
			return 999
	placeholder
		switch a.Type {
		case AccountTypeApiKey:
			if strings.TrimSpace(a.GetCredential("api_key")) != "" {
				return 0
		placeholder
			return 9
		case AccountTypeOAuth:
			if strings.TrimSpace(a.GetCredential("project_id")) == "" {
				return 1
		placeholder
			if strings.TrimSpace(a.GetCredential("oauth_type")) == "ai_studio" {
				return 2
		placeholder
			// Code Assist OAuth tokens often lack AI Studio scopes for models listing.
			return 3
		default:
			return 10
	placeholder
placeholder

	var selected *Account
	for i := range accounts {
		acc := &accounts[i]
		if selected == nil {
			selected = acc
			continue
	placeholder

		r1, r2 := rank(acc), rank(selected)
		if r1 < r2 {
			selected = acc
			continue
	placeholder
		if r1 > r2 {
			continue
	placeholder

		if acc.Priority < selected.Priority {
			selected = acc
	placeholder else if acc.Priority == selected.Priority {
			switch {
			case acc.LastUsedAt == nil && selected.LastUsedAt != nil:
				selected = acc
			case acc.LastUsedAt != nil && selected.LastUsedAt == nil:
				// keep selected
			case acc.LastUsedAt == nil && selected.LastUsedAt == nil:
				if acc.Type == AccountTypeOAuth && selected.Type != AccountTypeOAuth {
					selected = acc
			placeholder
			default:
				if acc.LastUsedAt.Before(*selected.LastUsedAt) {
					selected = acc
			placeholder
		placeholder
	placeholder
placeholder

	if selected == nil {
		return nil, errors.New("no available Gemini accounts")
placeholder
	return selected, nil
placeholder

func (s *GeminiMessagesCompatService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
	startTime := time.Now()

	var req struct {
		Model  string `json:"model"`
		Stream bool   `json:"stream"`
placeholder
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
placeholder
	if strings.TrimSpace(req.Model) == "" {
		return nil, fmt.Errorf("missing model")
placeholder

	originalModel := req.Model
	mappedModel := req.Model
	if account.Type == AccountTypeApiKey {
		mappedModel = account.GetMappedModel(req.Model)
placeholder

	geminiReq, err := convertClaudeMessagesToGeminiGenerateContent(body)
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	var requestIDHeader string
	var buildReq func(ctx context.Context) (*http.Request, string, error)
	useUpstreamStream := req.Stream
	if account.Type == AccountTypeOAuth && !req.Stream && strings.TrimSpace(account.GetCredential("project_id")) != "" {
		// Code Assist's non-streaming generateContent may return no content; use streaming upstream and aggregate.
		useUpstreamStream = true
placeholder

	switch account.Type {
	case AccountTypeApiKey:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			apiKey := account.GetCredential("api_key")
			if strings.TrimSpace(apiKey) == "" {
				return nil, "", errors.New("gemini api_key not configured")
		placeholder

			baseURL := strings.TrimRight(account.GetCredential("base_url"), "/")
			if baseURL == "" {
				baseURL = geminicli.AIStudioBaseURL
		placeholder

			action := "generateContent"
			if req.Stream {
				action = "streamGenerateContent"
		placeholder
			fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(baseURL, "/"), mappedModel, action)
			if req.Stream {
				fullURL += "?alt=sse"
		placeholder

			upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(geminiReq))
			if err != nil {
				return nil, "", err
		placeholder
			upstreamReq.Header.Set("Content-Type", "application/json")
			upstreamReq.Header.Set("x-goog-api-key", apiKey)
			return upstreamReq, "x-request-id", nil
	placeholder
		requestIDHeader = "x-request-id"

	case AccountTypeOAuth:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			if s.tokenProvider == nil {
				return nil, "", errors.New("gemini token provider not configured")
		placeholder
			accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
			if err != nil {
				return nil, "", err
		placeholder

			projectID := strings.TrimSpace(account.GetCredential("project_id"))

			action := "generateContent"
			if useUpstreamStream {
				action = "streamGenerateContent"
		placeholder

			// Two modes for OAuth:
			// 1. With project_id -> Code Assist API (wrapped request)
			// 2. Without project_id -> AI Studio API (direct OAuth, like API key but with Bearer token)
			if projectID != "" {
				// Mode 1: Code Assist API
				fullURL := fmt.Sprintf("%s/v1internal:%s", geminicli.GeminiCliBaseURL, action)
				if useUpstreamStream {
					fullURL += "?alt=sse"
			placeholder

				wrapped := map[string]any{
					"model":   mappedModel,
					"project": projectID,
			placeholder
				var inner any
				if err := json.Unmarshal(geminiReq, &inner); err != nil {
					return nil, "", fmt.Errorf("failed to parse gemini request: %w", err)
			placeholder
				wrapped["request"] = inner
				wrappedBytes, _ := json.Marshal(wrapped)

				upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(wrappedBytes))
				if err != nil {
					return nil, "", err
			placeholder
				upstreamReq.Header.Set("Content-Type", "application/json")
				upstreamReq.Header.Set("Authorization", "Bearer "+accessToken)
				upstreamReq.Header.Set("User-Agent", geminicli.GeminiCLIUserAgent)
				return upstreamReq, "x-request-id", nil
		placeholder else {
				// Mode 2: AI Studio API with OAuth (like API key mode, but using Bearer token)
				baseURL := strings.TrimRight(account.GetCredential("base_url"), "/")
				if baseURL == "" {
					baseURL = geminicli.AIStudioBaseURL
			placeholder

				fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", baseURL, mappedModel, action)
				if useUpstreamStream {
					fullURL += "?alt=sse"
			placeholder

				upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(geminiReq))
				if err != nil {
					return nil, "", err
			placeholder
				upstreamReq.Header.Set("Content-Type", "application/json")
				upstreamReq.Header.Set("Authorization", "Bearer "+accessToken)
				return upstreamReq, "x-request-id", nil
		placeholder
	placeholder
		requestIDHeader = "x-request-id"

	default:
		return nil, fmt.Errorf("unsupported account type: %s", account.Type)
placeholder

	var resp *http.Response
	for attempt := 1; attempt <= geminiMaxRetries; attempt++ {
		upstreamReq, idHeader, err := buildReq(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
		placeholder
			// Local build error: don't retry.
			if strings.Contains(err.Error(), "missing project_id") {
				return nil, s.writeClaudeError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		placeholder
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", err.Error())
	placeholder
		requestIDHeader = idHeader

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		if err != nil {
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream request failed, retry %d/%d: %v", account.ID, attempt, geminiMaxRetries, err)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed after retries: "+sanitizeUpstreamErrorMessage(err.Error()))
	placeholder

		if resp.StatusCode >= 400 && s.shouldRetryGeminiUpstreamError(account, resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
			// Don't treat insufficient-scope as transient.
			if resp.StatusCode == 403 && isGeminiInsufficientScope(resp.Header, respBody) {
				resp = &http.Response{
					StatusCode: resp.StatusCode,
					Header:     resp.Header.Clone(),
					Body:       io.NopCloser(bytes.NewReader(respBody)),
			placeholder
				break
		placeholder
			if resp.StatusCode == 429 {
				// Mark as rate-limited early so concurrent requests avoid this account.
				s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		placeholder
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream status %d, retry %d/%d", account.ID, resp.StatusCode, attempt, geminiMaxRetries)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			// Final attempt: surface the upstream error body (mapped below) instead of a generic retry error.
			resp = &http.Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header.Clone(),
				Body:       io.NopCloser(bytes.NewReader(respBody)),
		placeholder
			break
	placeholder

		break
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		if s.shouldFailoverGeminiUpstreamError(resp.StatusCode) {
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCodeplaceholder
	placeholder
		return nil, s.writeGeminiMappedError(c, resp.StatusCode, respBody)
placeholder

	requestID := resp.Header.Get(requestIDHeader)
	if requestID == "" {
		requestID = resp.Header.Get("x-goog-request-id")
placeholder
	if requestID != "" {
		c.Header("x-request-id", requestID)
placeholder

	var usage *ClaudeUsage
	var firstTokenMs *int
	if req.Stream {
		streamRes, err := s.handleStreamingResponse(c, resp, startTime, originalModel)
		if err != nil {
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder else {
		if useUpstreamStream {
			collected, usageObj, err := collectGeminiSSE(resp.Body, true)
			if err != nil {
				return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Failed to read upstream stream")
		placeholder
			claudeResp, usageObj2 := convertGeminiToClaudeMessage(collected, originalModel)
			c.JSON(http.StatusOK, claudeResp)
			usage = usageObj2
			if usageObj != nil && (usageObj.InputTokens > 0 || usageObj.OutputTokens > 0) {
				usage = usageObj
		placeholder
	placeholder else {
			usage, err = s.handleNonStreamingResponse(c, resp, originalModel)
			if err != nil {
				return nil, err
		placeholder
	placeholder
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel,
		Stream:       req.Stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
placeholder, nil
placeholder

func (s *GeminiMessagesCompatService) ForwardNative(ctx context.Context, c *gin.Context, account *Account, originalModel string, action string, stream bool, body []byte) (*ForwardResult, error) {
	startTime := time.Now()

	if strings.TrimSpace(originalModel) == "" {
		return nil, s.writeGoogleError(c, http.StatusBadRequest, "Missing model in URL")
placeholder
	if strings.TrimSpace(action) == "" {
		return nil, s.writeGoogleError(c, http.StatusBadRequest, "Missing action in URL")
placeholder
	if len(body) == 0 {
		return nil, s.writeGoogleError(c, http.StatusBadRequest, "Request body is empty")
placeholder

	switch action {
	case "generateContent", "streamGenerateContent", "countTokens":
		// ok
	default:
		return nil, s.writeGoogleError(c, http.StatusNotFound, "Unsupported action: "+action)
placeholder

	mappedModel := originalModel
	if account.Type == AccountTypeApiKey {
		mappedModel = account.GetMappedModel(originalModel)
placeholder

	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	useUpstreamStream := stream
	upstreamAction := action
	if account.Type == AccountTypeOAuth && !stream && action == "generateContent" && strings.TrimSpace(account.GetCredential("project_id")) != "" {
		// Code Assist's non-streaming generateContent may return no content; use streaming upstream and aggregate.
		useUpstreamStream = true
		upstreamAction = "streamGenerateContent"
placeholder
	forceAIStudio := action == "countTokens"

	var requestIDHeader string
	var buildReq func(ctx context.Context) (*http.Request, string, error)

	switch account.Type {
	case AccountTypeApiKey:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			apiKey := account.GetCredential("api_key")
			if strings.TrimSpace(apiKey) == "" {
				return nil, "", errors.New("gemini api_key not configured")
		placeholder

			baseURL := strings.TrimRight(account.GetCredential("base_url"), "/")
			if baseURL == "" {
				baseURL = geminicli.AIStudioBaseURL
		placeholder

			fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(baseURL, "/"), mappedModel, upstreamAction)
			if useUpstreamStream {
				fullURL += "?alt=sse"
		placeholder

			upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(body))
			if err != nil {
				return nil, "", err
		placeholder
			upstreamReq.Header.Set("Content-Type", "application/json")
			upstreamReq.Header.Set("x-goog-api-key", apiKey)
			return upstreamReq, "x-request-id", nil
	placeholder
		requestIDHeader = "x-request-id"

	case AccountTypeOAuth:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			if s.tokenProvider == nil {
				return nil, "", errors.New("gemini token provider not configured")
		placeholder
			accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
			if err != nil {
				return nil, "", err
		placeholder

			projectID := strings.TrimSpace(account.GetCredential("project_id"))

			// Two modes for OAuth:
			// 1. With project_id -> Code Assist API (wrapped request)
			// 2. Without project_id -> AI Studio API (direct OAuth, like API key but with Bearer token)
			if projectID != "" && !forceAIStudio {
				// Mode 1: Code Assist API
				fullURL := fmt.Sprintf("%s/v1internal:%s", geminicli.GeminiCliBaseURL, upstreamAction)
				if useUpstreamStream {
					fullURL += "?alt=sse"
			placeholder

				wrapped := map[string]any{
					"model":   mappedModel,
					"project": projectID,
			placeholder
				var inner any
				if err := json.Unmarshal(body, &inner); err != nil {
					return nil, "", fmt.Errorf("failed to parse gemini request: %w", err)
			placeholder
				wrapped["request"] = inner
				wrappedBytes, _ := json.Marshal(wrapped)

				upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(wrappedBytes))
				if err != nil {
					return nil, "", err
			placeholder
				upstreamReq.Header.Set("Content-Type", "application/json")
				upstreamReq.Header.Set("Authorization", "Bearer "+accessToken)
				upstreamReq.Header.Set("User-Agent", geminicli.GeminiCLIUserAgent)
				return upstreamReq, "x-request-id", nil
		placeholder else {
				// Mode 2: AI Studio API with OAuth (like API key mode, but using Bearer token)
				baseURL := strings.TrimRight(account.GetCredential("base_url"), "/")
				if baseURL == "" {
					baseURL = geminicli.AIStudioBaseURL
			placeholder

				fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", baseURL, mappedModel, upstreamAction)
				if useUpstreamStream {
					fullURL += "?alt=sse"
			placeholder

				upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(body))
				if err != nil {
					return nil, "", err
			placeholder
				upstreamReq.Header.Set("Content-Type", "application/json")
				upstreamReq.Header.Set("Authorization", "Bearer "+accessToken)
				return upstreamReq, "x-request-id", nil
		placeholder
	placeholder
		requestIDHeader = "x-request-id"

	default:
		return nil, s.writeGoogleError(c, http.StatusBadGateway, "Unsupported account type: "+account.Type)
placeholder

	var resp *http.Response
	for attempt := 1; attempt <= geminiMaxRetries; attempt++ {
		upstreamReq, idHeader, err := buildReq(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
		placeholder
			// Local build error: don't retry.
			if strings.Contains(err.Error(), "missing project_id") {
				return nil, s.writeGoogleError(c, http.StatusBadRequest, err.Error())
		placeholder
			return nil, s.writeGoogleError(c, http.StatusBadGateway, err.Error())
	placeholder
		requestIDHeader = idHeader

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		if err != nil {
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream request failed, retry %d/%d: %v", account.ID, attempt, geminiMaxRetries, err)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			if action == "countTokens" {
				estimated := estimateGeminiCountTokens(body)
				c.JSON(http.StatusOK, map[string]any{"totalTokens": estimatedplaceholder)
				return &ForwardResult{
					RequestID:    "",
					Usage:        ClaudeUsage{placeholder,
					Model:        originalModel,
					Stream:       false,
					Duration:     time.Since(startTime),
					FirstTokenMs: nil,
			placeholder, nil
		placeholder
			return nil, s.writeGoogleError(c, http.StatusBadGateway, "Upstream request failed after retries: "+sanitizeUpstreamErrorMessage(err.Error()))
	placeholder

		if resp.StatusCode >= 400 && s.shouldRetryGeminiUpstreamError(account, resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()
			// Don't treat insufficient-scope as transient.
			if resp.StatusCode == 403 && isGeminiInsufficientScope(resp.Header, respBody) {
				resp = &http.Response{
					StatusCode: resp.StatusCode,
					Header:     resp.Header.Clone(),
					Body:       io.NopCloser(bytes.NewReader(respBody)),
			placeholder
				break
		placeholder
			if resp.StatusCode == 429 {
				s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		placeholder
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream status %d, retry %d/%d", account.ID, resp.StatusCode, attempt, geminiMaxRetries)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			if action == "countTokens" {
				estimated := estimateGeminiCountTokens(body)
				c.JSON(http.StatusOK, map[string]any{"totalTokens": estimatedplaceholder)
				return &ForwardResult{
					RequestID:    "",
					Usage:        ClaudeUsage{placeholder,
					Model:        originalModel,
					Stream:       false,
					Duration:     time.Since(startTime),
					FirstTokenMs: nil,
			placeholder, nil
		placeholder
			// Final attempt: surface the upstream error body (passed through below) instead of a generic retry error.
			resp = &http.Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header.Clone(),
				Body:       io.NopCloser(bytes.NewReader(respBody)),
		placeholder
			break
	placeholder

		break
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	requestID := resp.Header.Get(requestIDHeader)
	if requestID == "" {
		requestID = resp.Header.Get("x-goog-request-id")
placeholder
	if requestID != "" {
		c.Header("x-request-id", requestID)
placeholder

	isOAuth := account.Type == AccountTypeOAuth

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)

		// Best-effort fallback for OAuth tokens missing AI Studio scopes when calling countTokens.
		// This avoids Gemini SDKs failing hard during preflight token counting.
		if action == "countTokens" && isOAuth && isGeminiInsufficientScope(resp.Header, respBody) {
			estimated := estimateGeminiCountTokens(body)
			c.JSON(http.StatusOK, map[string]any{"totalTokens": estimatedplaceholder)
			return &ForwardResult{
				RequestID:    requestID,
				Usage:        ClaudeUsage{placeholder,
				Model:        originalModel,
				Stream:       false,
				Duration:     time.Since(startTime),
				FirstTokenMs: nil,
		placeholder, nil
	placeholder

		if s.shouldFailoverGeminiUpstreamError(resp.StatusCode) {
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCodeplaceholder
	placeholder

		respBody = unwrapIfNeeded(isOAuth, respBody)
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
	placeholder
		c.Data(resp.StatusCode, contentType, respBody)
		return nil, fmt.Errorf("gemini upstream error: %d", resp.StatusCode)
placeholder

	var usage *ClaudeUsage
	var firstTokenMs *int

	if stream {
		streamRes, err := s.handleNativeStreamingResponse(c, resp, startTime, isOAuth)
		if err != nil {
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder else {
		if useUpstreamStream {
			collected, usageObj, err := collectGeminiSSE(resp.Body, isOAuth)
			if err != nil {
				return nil, s.writeGoogleError(c, http.StatusBadGateway, "Failed to read upstream stream")
		placeholder
			b, _ := json.Marshal(collected)
			c.Data(http.StatusOK, "application/json", b)
			usage = usageObj
	placeholder else {
			usageResp, err := s.handleNativeNonStreamingResponse(c, resp, isOAuth)
			if err != nil {
				return nil, err
		placeholder
			usage = usageResp
	placeholder
placeholder

	if usage == nil {
		usage = &ClaudeUsage{placeholder
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel,
		Stream:       stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
placeholder, nil
placeholder

func (s *GeminiMessagesCompatService) shouldRetryGeminiUpstreamError(account *Account, statusCode int) bool {
	switch statusCode {
	case 429, 500, 502, 503, 504, 529:
		return true
	case 403:
		// GeminiCli OAuth occasionally returns 403 transiently (activation/quota propagation); allow retry.
		if account == nil || account.Type != AccountTypeOAuth {
			return false
	placeholder
		oauthType := strings.ToLower(strings.TrimSpace(account.GetCredential("oauth_type")))
		if oauthType == "" && strings.TrimSpace(account.GetCredential("project_id")) != "" {
			// Legacy/implicit Code Assist OAuth accounts.
			oauthType = "code_assist"
	placeholder
		return oauthType == "code_assist"
	default:
		return false
placeholder
placeholder

func (s *GeminiMessagesCompatService) shouldFailoverGeminiUpstreamError(statusCode int) bool {
	switch statusCode {
	case 401, 403, 429, 529:
		return true
	default:
		return statusCode >= 500
placeholder
placeholder

func sleepGeminiBackoff(attempt int) {
	delay := geminiRetryBaseDelay * time.Duration(1<<uint(attempt-1))
	if delay > geminiRetryMaxDelay {
		delay = geminiRetryMaxDelay
placeholder

	// +/- 20% jitter
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	jitter := time.Duration(float64(delay) * 0.2 * (r.Float64()*2 - 1))
	sleepFor := delay + jitter
	if sleepFor < 0 {
		sleepFor = 0
placeholder
	time.Sleep(sleepFor)
placeholder

var (
	sensitiveQueryParamRegex = regexp.MustCompile(`(?i)([?&](?:key|client_secret|access_token|refresh_token)=)[^&"\s]+`)
	retryInRegex             = regexp.MustCompile(`Please retry in ([0-9.]+)s`)
)

func sanitizeUpstreamErrorMessage(msg string) string {
	if msg == "" {
		return msg
placeholder
	return sensitiveQueryParamRegex.ReplaceAllString(msg, `$1***`)
placeholder

func (s *GeminiMessagesCompatService) writeGeminiMappedError(c *gin.Context, upstreamStatus int, body []byte) error {
	var statusCode int
	var errType, errMsg string

	if mapped := mapGeminiErrorBodyToClaudeError(body); mapped != nil {
		errType = mapped.Type
		if mapped.Message != "" {
			errMsg = mapped.Message
	placeholder
		if mapped.StatusCode > 0 {
			statusCode = mapped.StatusCode
	placeholder
placeholder

	switch upstreamStatus {
	case 400:
		if statusCode == 0 {
			statusCode = http.StatusBadRequest
	placeholder
		if errType == "" {
			errType = "invalid_request_error"
	placeholder
		if errMsg == "" {
			errMsg = "Invalid request"
	placeholder
	case 401:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			errType = "authentication_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream authentication failed, please contact administrator"
	placeholder
	case 403:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			errType = "permission_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream access forbidden, please contact administrator"
	placeholder
	case 404:
		if statusCode == 0 {
			statusCode = http.StatusNotFound
	placeholder
		if errType == "" {
			errType = "not_found_error"
	placeholder
		if errMsg == "" {
			errMsg = "Resource not found"
	placeholder
	case 429:
		if statusCode == 0 {
			statusCode = http.StatusTooManyRequests
	placeholder
		if errType == "" {
			errType = "rate_limit_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream rate limit exceeded, please retry later"
	placeholder
	case 529:
		if statusCode == 0 {
			statusCode = http.StatusServiceUnavailable
	placeholder
		if errType == "" {
			errType = "overloaded_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream service overloaded, please retry later"
	placeholder
	case 500, 502, 503, 504:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			switch upstreamStatus {
			case 504:
				errType = "timeout_error"
			case 503:
				errType = "overloaded_error"
			default:
				errType = "api_error"
		placeholder
	placeholder
		if errMsg == "" {
			errMsg = "Upstream service temporarily unavailable"
	placeholder
	default:
		if statusCode == 0 {
			statusCode = http.StatusBadGateway
	placeholder
		if errType == "" {
			errType = "upstream_error"
	placeholder
		if errMsg == "" {
			errMsg = "Upstream request failed"
	placeholder
placeholder

	c.JSON(statusCode, gin.H{
		"type":  "error",
		"error": gin.H{"type": errType, "message": errMsgplaceholder,
placeholder)
	return fmt.Errorf("upstream error: %d", upstreamStatus)
placeholder

type claudeErrorMapping struct {
	Type       string
	Message    string
	StatusCode int
placeholder

func mapGeminiErrorBodyToClaudeError(body []byte) *claudeErrorMapping {
	if len(body) == 0 {
		return nil
placeholder

	var parsed struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
	placeholder `json:"error"`
placeholder
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil
placeholder
	if strings.TrimSpace(parsed.Error.Status) == "" && parsed.Error.Code == 0 && strings.TrimSpace(parsed.Error.Message) == "" {
		return nil
placeholder

	mapped := &claudeErrorMapping{
		Type:    mapGeminiStatusToClaudeErrorType(parsed.Error.Status),
		Message: "",
placeholder
	if mapped.Type == "" {
		mapped.Type = "upstream_error"
placeholder

	switch strings.ToUpper(strings.TrimSpace(parsed.Error.Status)) {
	case "INVALID_ARGUMENT":
		mapped.StatusCode = http.StatusBadRequest
	case "NOT_FOUND":
		mapped.StatusCode = http.StatusNotFound
	case "RESOURCE_EXHAUSTED":
		mapped.StatusCode = http.StatusTooManyRequests
	default:
		// Keep StatusCode unset and let HTTP status mapping decide.
placeholder

	// Keep messages generic by default; upstream error message can be long or include sensitive fragments.
	return mapped
placeholder

func mapGeminiStatusToClaudeErrorType(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "INVALID_ARGUMENT":
		return "invalid_request_error"
	case "PERMISSION_DENIED":
		return "permission_error"
	case "NOT_FOUND":
		return "not_found_error"
	case "RESOURCE_EXHAUSTED":
		return "rate_limit_error"
	case "UNAUTHENTICATED":
		return "authentication_error"
	case "UNAVAILABLE":
		return "overloaded_error"
	case "INTERNAL":
		return "api_error"
	case "DEADLINE_EXCEEDED":
		return "timeout_error"
	default:
		return ""
placeholder
placeholder

type geminiStreamResult struct {
	usage        *ClaudeUsage
	firstTokenMs *int
placeholder

func (s *GeminiMessagesCompatService) handleNonStreamingResponse(c *gin.Context, resp *http.Response, originalModel string) (*ClaudeUsage, error) {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Failed to read upstream response")
placeholder

	geminiResp, err := unwrapGeminiResponse(body)
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Failed to parse upstream response")
placeholder

	claudeResp, usage := convertGeminiToClaudeMessage(geminiResp, originalModel)
	c.JSON(http.StatusOK, claudeResp)

	return usage, nil
placeholder

func (s *GeminiMessagesCompatService) handleStreamingResponse(c *gin.Context, resp *http.Response, startTime time.Time, originalModel string) (*geminiStreamResult, error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder

	messageID := "msg_" + randomHex(12)
	messageStart := map[string]any{
		"type": "message_start",
		"message": map[string]any{
			"id":            messageID,
			"type":          "message",
			"role":          "assistant",
			"model":         originalModel,
			"content":       []any{placeholder,
			"stop_reason":   nil,
			"stop_sequence": nil,
			"usage": map[string]any{
				"input_tokens":  0,
				"output_tokens": 0,
		placeholder,
	placeholder,
placeholder
	writeSSE(c.Writer, "message_start", messageStart)
	flusher.Flush()

	var firstTokenMs *int
	var usage ClaudeUsage
	finishReason := ""
	sawToolUse := false

	nextBlockIndex := 0
	openBlockIndex := -1
	openBlockType := ""
	seenText := ""
	openToolIndex := -1
	openToolID := ""
	openToolName := ""
	seenToolJSON := ""

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("stream read error: %w", err)
	placeholder

		if !strings.HasPrefix(line, "data:") {
			if errors.Is(err, io.EOF) {
				break
		placeholder
			continue
	placeholder
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "" || payload == "[DONE]" {
			if errors.Is(err, io.EOF) {
				break
		placeholder
			continue
	placeholder

		geminiResp, err := unwrapGeminiResponse([]byte(payload))
		if err != nil {
			continue
	placeholder

		if fr := extractGeminiFinishReason(geminiResp); fr != "" {
			finishReason = fr
	placeholder

		parts := extractGeminiParts(geminiResp)
		for _, part := range parts {
			if text, ok := part["text"].(string); ok && text != "" {
				delta, newSeen := computeGeminiTextDelta(seenText, text)
				seenText = newSeen
				if delta == "" {
					continue
			placeholder

				if openBlockType != "text" {
					if openBlockIndex >= 0 {
						writeSSE(c.Writer, "content_block_stop", map[string]any{
							"type":  "content_block_stop",
							"index": openBlockIndex,
					placeholder)
				placeholder
					openBlockType = "text"
					openBlockIndex = nextBlockIndex
					nextBlockIndex++
					writeSSE(c.Writer, "content_block_start", map[string]any{
						"type":  "content_block_start",
						"index": openBlockIndex,
						"content_block": map[string]any{
							"type": "text",
							"text": "",
					placeholder,
				placeholder)
			placeholder

				if firstTokenMs == nil {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
			placeholder
				writeSSE(c.Writer, "content_block_delta", map[string]any{
					"type":  "content_block_delta",
					"index": openBlockIndex,
					"delta": map[string]any{
						"type": "text_delta",
						"text": delta,
				placeholder,
			placeholder)
				flusher.Flush()
				continue
		placeholder

			if fc, ok := part["functionCall"].(map[string]any); ok && fc != nil {
				name, _ := fc["name"].(string)
				args := fc["args"]
				if strings.TrimSpace(name) == "" {
					name = "tool"
			placeholder

				// Close any open text block before tool_use.
				if openBlockIndex >= 0 {
					writeSSE(c.Writer, "content_block_stop", map[string]any{
						"type":  "content_block_stop",
						"index": openBlockIndex,
				placeholder)
					openBlockIndex = -1
					openBlockType = ""
			placeholder

				// If we receive streamed tool args in pieces, keep a single tool block open and emit deltas.
				if openToolIndex >= 0 && openToolName != name {
					writeSSE(c.Writer, "content_block_stop", map[string]any{
						"type":  "content_block_stop",
						"index": openToolIndex,
				placeholder)
					openToolIndex = -1
					openToolName = ""
					seenToolJSON = ""
			placeholder

				if openToolIndex < 0 {
					openToolID = "toolu_" + randomHex(8)
					openToolIndex = nextBlockIndex
					openToolName = name
					nextBlockIndex++
					sawToolUse = true

					writeSSE(c.Writer, "content_block_start", map[string]any{
						"type":  "content_block_start",
						"index": openToolIndex,
						"content_block": map[string]any{
							"type":  "tool_use",
							"id":    openToolID,
							"name":  name,
							"input": map[string]any{placeholder,
					placeholder,
				placeholder)
			placeholder

				argsJSONText := "{placeholder"
				switch v := args.(type) {
				case nil:
					// keep default "{placeholder"
				case string:
					if strings.TrimSpace(v) != "" {
						argsJSONText = v
				placeholder
				default:
					if b, err := json.Marshal(args); err == nil && len(b) > 0 {
						argsJSONText = string(b)
				placeholder
			placeholder

				delta, newSeen := computeGeminiTextDelta(seenToolJSON, argsJSONText)
				seenToolJSON = newSeen
				if delta != "" {
					writeSSE(c.Writer, "content_block_delta", map[string]any{
						"type":  "content_block_delta",
						"index": openToolIndex,
						"delta": map[string]any{
							"type":         "input_json_delta",
							"partial_json": delta,
					placeholder,
				placeholder)
			placeholder
				flusher.Flush()
		placeholder
	placeholder

		if u := extractGeminiUsage(geminiResp); u != nil {
			usage = *u
	placeholder

		// Process the final unterminated line at EOF as well.
		if errors.Is(err, io.EOF) {
			break
	placeholder
placeholder

	if openBlockIndex >= 0 {
		writeSSE(c.Writer, "content_block_stop", map[string]any{
			"type":  "content_block_stop",
			"index": openBlockIndex,
	placeholder)
placeholder
	if openToolIndex >= 0 {
		writeSSE(c.Writer, "content_block_stop", map[string]any{
			"type":  "content_block_stop",
			"index": openToolIndex,
	placeholder)
placeholder

	stopReason := mapGeminiFinishReasonToClaudeStopReason(finishReason)
	if sawToolUse {
		stopReason = "tool_use"
placeholder

	usageObj := map[string]any{
		"output_tokens": usage.OutputTokens,
placeholder
	if usage.InputTokens > 0 {
		usageObj["input_tokens"] = usage.InputTokens
placeholder
	writeSSE(c.Writer, "message_delta", map[string]any{
		"type": "message_delta",
		"delta": map[string]any{
			"stop_reason":   stopReason,
			"stop_sequence": nil,
	placeholder,
		"usage": usageObj,
placeholder)
	writeSSE(c.Writer, "message_stop", map[string]any{
		"type": "message_stop",
placeholder)
	flusher.Flush()

	return &geminiStreamResult{usage: &usage, firstTokenMs: firstTokenMsplaceholder, nil
placeholder

func writeSSE(w io.Writer, event string, data any) {
	if event != "" {
		_, _ = fmt.Fprintf(w, "event: %s\n", event)
placeholder
	b, _ := json.Marshal(data)
	_, _ = fmt.Fprintf(w, "data: %s\n\n", string(b))
placeholder

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
placeholder

func (s *GeminiMessagesCompatService) writeClaudeError(c *gin.Context, status int, errType, message string) error {
	c.JSON(status, gin.H{
		"type":  "error",
		"error": gin.H{"type": errType, "message": messageplaceholder,
placeholder)
	return fmt.Errorf("%s", message)
placeholder

func (s *GeminiMessagesCompatService) writeGoogleError(c *gin.Context, status int, message string) error {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    status,
			"message": message,
			"status":  googleapi.HTTPStatusToGoogleStatus(status),
	placeholder,
placeholder)
	return fmt.Errorf("%s", message)
placeholder

func unwrapIfNeeded(isOAuth bool, raw []byte) []byte {
	if !isOAuth {
		return raw
placeholder
	inner, err := unwrapGeminiResponse(raw)
	if err != nil {
		return raw
placeholder
	b, err := json.Marshal(inner)
	if err != nil {
		return raw
placeholder
	return b
placeholder

func collectGeminiSSE(body io.Reader, isOAuth bool) (map[string]any, *ClaudeUsage, error) {
	reader := bufio.NewReader(body)

	var last map[string]any
	var lastWithParts map[string]any
	usage := &ClaudeUsage{placeholder

	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			trimmed := strings.TrimRight(line, "\r\n")
			if strings.HasPrefix(trimmed, "data:") {
				payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "data:"))
				switch payload {
				case "", "[DONE]":
					if payload == "[DONE]" {
						return pickGeminiCollectResult(last, lastWithParts), usage, nil
				placeholder
				default:
					var parsed map[string]any
					if isOAuth {
						inner, err := unwrapGeminiResponse([]byte(payload))
						if err == nil && inner != nil {
							parsed = inner
					placeholder
				placeholder else {
						_ = json.Unmarshal([]byte(payload), &parsed)
				placeholder
					if parsed != nil {
						last = parsed
						if u := extractGeminiUsage(parsed); u != nil {
							usage = u
					placeholder
						if parts := extractGeminiParts(parsed); len(parts) > 0 {
							lastWithParts = parsed
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder

		if errors.Is(err, io.EOF) {
			break
	placeholder
		if err != nil {
			return nil, nil, err
	placeholder
placeholder

	return pickGeminiCollectResult(last, lastWithParts), usage, nil
placeholder

func pickGeminiCollectResult(last map[string]any, lastWithParts map[string]any) map[string]any {
	if lastWithParts != nil {
		return lastWithParts
placeholder
	if last != nil {
		return last
placeholder
	return map[string]any{placeholder
placeholder

type geminiNativeStreamResult struct {
	usage        *ClaudeUsage
	firstTokenMs *int
placeholder

func isGeminiInsufficientScope(headers http.Header, body []byte) bool {
	if strings.Contains(strings.ToLower(headers.Get("Www-Authenticate")), "insufficient_scope") {
		return true
placeholder
	lower := strings.ToLower(string(body))
	return strings.Contains(lower, "insufficient authentication scopes") || strings.Contains(lower, "access_token_scope_insufficient")
placeholder

func estimateGeminiCountTokens(reqBody []byte) int {
	var obj map[string]any
	if err := json.Unmarshal(reqBody, &obj); err != nil {
		return 0
placeholder

	var texts []string

	// systemInstruction.parts[].text
	if si, ok := obj["systemInstruction"].(map[string]any); ok {
		if parts, ok := si["parts"].([]any); ok {
			for _, p := range parts {
				if pm, ok := p.(map[string]any); ok {
					if t, ok := pm["text"].(string); ok && strings.TrimSpace(t) != "" {
						texts = append(texts, t)
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// contents[].parts[].text
	if contents, ok := obj["contents"].([]any); ok {
		for _, c := range contents {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
		placeholder
			parts, ok := cm["parts"].([]any)
			if !ok {
				continue
		placeholder
			for _, p := range parts {
				pm, ok := p.(map[string]any)
				if !ok {
					continue
			placeholder
				if t, ok := pm["text"].(string); ok && strings.TrimSpace(t) != "" {
					texts = append(texts, t)
			placeholder
		placeholder
	placeholder
placeholder

	total := 0
	for _, t := range texts {
		total += estimateTokensForText(t)
placeholder
	if total < 0 {
		return 0
placeholder
	return total
placeholder

func estimateTokensForText(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
placeholder
	runes := []rune(s)
	if len(runes) == 0 {
		return 0
placeholder
	ascii := 0
	for _, r := range runes {
		if r <= 0x7f {
			ascii++
	placeholder
placeholder
	asciiRatio := float64(ascii) / float64(len(runes))
	if asciiRatio >= 0.8 {
		// Roughly 4 chars per token for English-like text.
		return (len(runes) + 3) / 4
placeholder
	// For CJK-heavy text, approximate 1 rune per token.
	return len(runes)
placeholder

type UpstreamHTTPResult struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
placeholder

func (s *GeminiMessagesCompatService) handleNativeNonStreamingResponse(c *gin.Context, resp *http.Response, isOAuth bool) (*ClaudeUsage, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
placeholder

	var parsed map[string]any
	if isOAuth {
		parsed, err = unwrapGeminiResponse(respBody)
		if err == nil && parsed != nil {
			respBody, _ = json.Marshal(parsed)
	placeholder
placeholder else {
		_ = json.Unmarshal(respBody, &parsed)
placeholder

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
placeholder
	c.Data(resp.StatusCode, contentType, respBody)

	if parsed != nil {
		if u := extractGeminiUsage(parsed); u != nil {
			return u, nil
	placeholder
placeholder
	return &ClaudeUsage{placeholder, nil
placeholder

func (s *GeminiMessagesCompatService) handleNativeStreamingResponse(c *gin.Context, resp *http.Response, startTime time.Time, isOAuth bool) (*geminiNativeStreamResult, error) {
	c.Status(resp.StatusCode)
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/event-stream; charset=utf-8"
placeholder
	c.Header("Content-Type", contentType)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder

	reader := bufio.NewReader(resp.Body)
	usage := &ClaudeUsage{placeholder
	var firstTokenMs *int

	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			trimmed := strings.TrimRight(line, "\r\n")
			if strings.HasPrefix(trimmed, "data:") {
				payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "data:"))
				// Keepalive / done markers
				if payload == "" || payload == "[DONE]" {
					_, _ = io.WriteString(c.Writer, line)
					flusher.Flush()
			placeholder else {
					var rawToWrite string
					rawToWrite = payload

					var parsed map[string]any
					if isOAuth {
						inner, err := unwrapGeminiResponse([]byte(payload))
						if err == nil && inner != nil {
							parsed = inner
							if b, err := json.Marshal(inner); err == nil {
								rawToWrite = string(b)
						placeholder
					placeholder
				placeholder else {
						_ = json.Unmarshal([]byte(payload), &parsed)
				placeholder

					if parsed != nil {
						if u := extractGeminiUsage(parsed); u != nil {
							usage = u
					placeholder
				placeholder

					if firstTokenMs == nil {
						ms := int(time.Since(startTime).Milliseconds())
						firstTokenMs = &ms
				placeholder

					if isOAuth {
						// SSE format requires double newline (\n\n) to separate events
						_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", rawToWrite)
				placeholder else {
						// Pass-through for AI Studio responses.
						_, _ = io.WriteString(c.Writer, line)
				placeholder
					flusher.Flush()
			placeholder
		placeholder else {
				_, _ = io.WriteString(c.Writer, line)
				flusher.Flush()
		placeholder
	placeholder

		if errors.Is(err, io.EOF) {
			break
	placeholder
		if err != nil {
			return nil, err
	placeholder
placeholder

	return &geminiNativeStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, nil
placeholder

// ForwardAIStudioGET forwards a GET request to AI Studio (generativelanguage.googleapis.com) for
// endpoints like /v1beta/models and /v1beta/models/{modelplaceholder.
//
// This is used to support Gemini SDKs that call models listing endpoints before generation.
func (s *GeminiMessagesCompatService) ForwardAIStudioGET(ctx context.Context, account *Account, path string) (*UpstreamHTTPResult, error) {
	if account == nil {
		return nil, errors.New("account is nil")
placeholder
	path = strings.TrimSpace(path)
	if path == "" || !strings.HasPrefix(path, "/") {
		return nil, errors.New("invalid path")
placeholder

	baseURL := strings.TrimRight(account.GetCredential("base_url"), "/")
	if baseURL == "" {
		baseURL = geminicli.AIStudioBaseURL
placeholder
	fullURL := strings.TrimRight(baseURL, "/") + path

	var proxyURL string
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
placeholder

	switch account.Type {
	case AccountTypeApiKey:
		apiKey := strings.TrimSpace(account.GetCredential("api_key"))
		if apiKey == "" {
			return nil, errors.New("gemini api_key not configured")
	placeholder
		req.Header.Set("x-goog-api-key", apiKey)
	case AccountTypeOAuth:
		if s.tokenProvider == nil {
			return nil, errors.New("gemini token provider not configured")
	placeholder
		accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
		if err != nil {
			return nil, err
	placeholder
		req.Header.Set("Authorization", "Bearer "+accessToken)
	default:
		return nil, fmt.Errorf("unsupported account type: %s", account.Type)
placeholder

	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	return &UpstreamHTTPResult{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
placeholder, nil
placeholder

func unwrapGeminiResponse(raw []byte) (map[string]any, error) {
	var outer map[string]any
	if err := json.Unmarshal(raw, &outer); err != nil {
		return nil, err
placeholder
	if resp, ok := outer["response"].(map[string]any); ok && resp != nil {
		return resp, nil
placeholder
	return outer, nil
placeholder

func convertGeminiToClaudeMessage(geminiResp map[string]any, originalModel string) (map[string]any, *ClaudeUsage) {
	usage := extractGeminiUsage(geminiResp)
	if usage == nil {
		usage = &ClaudeUsage{placeholder
placeholder

	contentBlocks := make([]any, 0)
	sawToolUse := false
	if candidates, ok := geminiResp["candidates"].([]any); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]any); ok {
			if content, ok := cand["content"].(map[string]any); ok {
				if parts, ok := content["parts"].([]any); ok {
					for _, part := range parts {
						pm, ok := part.(map[string]any)
						if !ok {
							continue
					placeholder
						if text, ok := pm["text"].(string); ok && text != "" {
							contentBlocks = append(contentBlocks, map[string]any{
								"type": "text",
								"text": text,
						placeholder)
					placeholder
						if fc, ok := pm["functionCall"].(map[string]any); ok {
							name, _ := fc["name"].(string)
							if strings.TrimSpace(name) == "" {
								name = "tool"
						placeholder
							args := fc["args"]
							sawToolUse = true
							contentBlocks = append(contentBlocks, map[string]any{
								"type":  "tool_use",
								"id":    "toolu_" + randomHex(8),
								"name":  name,
								"input": args,
						placeholder)
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	stopReason := mapGeminiFinishReasonToClaudeStopReason(extractGeminiFinishReason(geminiResp))
	if sawToolUse {
		stopReason = "tool_use"
placeholder

	resp := map[string]any{
		"id":            "msg_" + randomHex(12),
		"type":          "message",
		"role":          "assistant",
		"model":         originalModel,
		"content":       contentBlocks,
		"stop_reason":   stopReason,
		"stop_sequence": nil,
		"usage": map[string]any{
			"input_tokens":  usage.InputTokens,
			"output_tokens": usage.OutputTokens,
	placeholder,
placeholder

	return resp, usage
placeholder

func extractGeminiUsage(geminiResp map[string]any) *ClaudeUsage {
	usageMeta, ok := geminiResp["usageMetadata"].(map[string]any)
	if !ok || usageMeta == nil {
		return nil
placeholder
	prompt, _ := asInt(usageMeta["promptTokenCount"])
	cand, _ := asInt(usageMeta["candidatesTokenCount"])
	return &ClaudeUsage{
		InputTokens:  prompt,
		OutputTokens: cand,
placeholder
placeholder

func asInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case int:
		return t, true
	case int64:
		return int(t), true
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			return 0, false
	placeholder
		return int(i), true
	default:
		return 0, false
placeholder
placeholder

func (s *GeminiMessagesCompatService) handleGeminiUpstreamError(ctx context.Context, account *Account, statusCode int, headers http.Header, body []byte) {
	if s.rateLimitService != nil && (statusCode == 401 || statusCode == 403 || statusCode == 529) {
		s.rateLimitService.HandleUpstreamError(ctx, account, statusCode, headers, body)
		return
placeholder
	if statusCode != 429 {
		return
placeholder

	oauthType := account.GeminiOAuthType()
	tierID := account.GeminiTierID()
	projectID := strings.TrimSpace(account.GetCredential("project_id"))
	isCodeAssist := account.IsGeminiCodeAssist()

	resetAt := ParseGeminiRateLimitResetTime(body)
	if resetAt == nil {
		// 根据账号类型使用不同的默认重置时间
		var ra time.Time
		if isCodeAssist {
			// Code Assist: fallback cooldown by tier
			cooldown := geminiCooldownForTier(tierID)
			if s.rateLimitService != nil {
				cooldown = s.rateLimitService.GeminiCooldown(ctx, account)
		placeholder
			ra = time.Now().Add(cooldown)
			log.Printf("[Gemini 429] Account %d (Code Assist, tier=%s, project=%s) rate limited, cooldown=%v", account.ID, tierID, projectID, time.Until(ra).Truncate(time.Second))
	placeholder else {
			// API Key / AI Studio OAuth: PST 午夜
			if ts := nextGeminiDailyResetUnix(); ts != nil {
				ra = time.Unix(*ts, 0)
				log.Printf("[Gemini 429] Account %d (API Key/AI Studio, type=%s) rate limited, reset at PST midnight (%v)", account.ID, account.Type, ra)
		placeholder else {
				// 兜底：5 分钟
				ra = time.Now().Add(5 * time.Minute)
				log.Printf("[Gemini 429] Account %d rate limited, fallback to 5min", account.ID)
		placeholder
	placeholder
		_ = s.accountRepo.SetRateLimited(ctx, account.ID, ra)
		return
placeholder

	// 使用解析到的重置时间
	resetTime := time.Unix(*resetAt, 0)
	_ = s.accountRepo.SetRateLimited(ctx, account.ID, resetTime)
	log.Printf("[Gemini 429] Account %d rate limited until %v (oauth_type=%s, tier=%s)",
		account.ID, resetTime, oauthType, tierID)
placeholder

// ParseGeminiRateLimitResetTime 解析 Gemini 格式的 429 响应，返回重置时间的 Unix 时间戳
func ParseGeminiRateLimitResetTime(body []byte) *int64 {
	// Try to parse metadata.quotaResetDelay like "12.345s"
	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err == nil {
		if errObj, ok := parsed["error"].(map[string]any); ok {
			if msg, ok := errObj["message"].(string); ok {
				if looksLikeGeminiDailyQuota(msg) {
					if ts := nextGeminiDailyResetUnix(); ts != nil {
						return ts
				placeholder
			placeholder
		placeholder
			if details, ok := errObj["details"].([]any); ok {
				for _, d := range details {
					dm, ok := d.(map[string]any)
					if !ok {
						continue
				placeholder
					if meta, ok := dm["metadata"].(map[string]any); ok {
						if v, ok := meta["quotaResetDelay"].(string); ok {
							if dur, err := time.ParseDuration(v); err == nil {
								ts := time.Now().Unix() + int64(dur.Seconds())
								return &ts
						placeholder
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// Match "Please retry in Xs"
	matches := retryInRegex.FindStringSubmatch(string(body))
	if len(matches) == 2 {
		if dur, err := time.ParseDuration(matches[1] + "s"); err == nil {
			ts := time.Now().Unix() + int64(math.Ceil(dur.Seconds()))
			return &ts
	placeholder
placeholder

	return nil
placeholder

func looksLikeGeminiDailyQuota(message string) bool {
	m := strings.ToLower(message)
	if strings.Contains(m, "per day") || strings.Contains(m, "requests per day") || strings.Contains(m, "quota") && strings.Contains(m, "per day") {
		return true
placeholder
	return false
placeholder

func nextGeminiDailyResetUnix() *int64 {
	reset := geminiDailyResetTime(time.Now())
	ts := reset.Unix()
	return &ts
placeholder

func extractGeminiFinishReason(geminiResp map[string]any) string {
	if candidates, ok := geminiResp["candidates"].([]any); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]any); ok {
			if fr, ok := cand["finishReason"].(string); ok {
				return fr
		placeholder
	placeholder
placeholder
	return ""
placeholder

func extractGeminiParts(geminiResp map[string]any) []map[string]any {
	if candidates, ok := geminiResp["candidates"].([]any); ok && len(candidates) > 0 {
		if cand, ok := candidates[0].(map[string]any); ok {
			if content, ok := cand["content"].(map[string]any); ok {
				if partsAny, ok := content["parts"].([]any); ok && len(partsAny) > 0 {
					out := make([]map[string]any, 0, len(partsAny))
					for _, p := range partsAny {
						pm, ok := p.(map[string]any)
						if !ok {
							continue
					placeholder
						out = append(out, pm)
				placeholder
					return out
			placeholder
		placeholder
	placeholder
placeholder
	return nil
placeholder

func computeGeminiTextDelta(seen, incoming string) (delta, newSeen string) {
	incoming = strings.TrimSuffix(incoming, "\u0000")
	if incoming == "" {
		return "", seen
placeholder

	// Cumulative mode: incoming contains full text so far.
	if strings.HasPrefix(incoming, seen) {
		return strings.TrimPrefix(incoming, seen), incoming
placeholder
	// Duplicate/rewind: ignore.
	if strings.HasPrefix(seen, incoming) {
		return "", seen
placeholder
	// Delta mode: treat incoming as incremental chunk.
	return incoming, seen + incoming
placeholder

func mapGeminiFinishReasonToClaudeStopReason(finishReason string) string {
	switch strings.ToUpper(strings.TrimSpace(finishReason)) {
	case "MAX_TOKENS":
		return "max_tokens"
	case "STOP":
		return "end_turn"
	default:
		return "end_turn"
placeholder
placeholder

func convertClaudeMessagesToGeminiGenerateContent(body []byte) ([]byte, error) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
placeholder

	toolUseIDToName := make(map[string]string)

	systemText := extractClaudeSystemText(req["system"])
	contents, err := convertClaudeMessagesToGeminiContents(req["messages"], toolUseIDToName)
	if err != nil {
		return nil, err
placeholder

	out := make(map[string]any)
	if systemText != "" {
		out["systemInstruction"] = map[string]any{
			"parts": []any{map[string]any{"text": systemTextplaceholderplaceholder,
	placeholder
placeholder
	out["contents"] = contents

	if tools := convertClaudeToolsToGeminiTools(req["tools"]); tools != nil {
		out["tools"] = tools
placeholder

	generationConfig := convertClaudeGenerationConfig(req)
	if generationConfig != nil {
		out["generationConfig"] = generationConfig
placeholder

	stripGeminiFunctionIDs(out)
	return json.Marshal(out)
placeholder

func stripGeminiFunctionIDs(req map[string]any) {
	// Defensive cleanup: some upstreams reject unexpected `id` fields in functionCall/functionResponse.
	contents, ok := req["contents"].([]any)
	if !ok {
		return
placeholder
	for _, c := range contents {
		cm, ok := c.(map[string]any)
		if !ok {
			continue
	placeholder
		contentParts, ok := cm["parts"].([]any)
		if !ok {
			continue
	placeholder
		for _, p := range contentParts {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
		placeholder
			if fc, ok := pm["functionCall"].(map[string]any); ok && fc != nil {
				delete(fc, "id")
		placeholder
			if fr, ok := pm["functionResponse"].(map[string]any); ok && fr != nil {
				delete(fr, "id")
		placeholder
	placeholder
placeholder
placeholder

func extractClaudeSystemText(system any) string {
	switch v := system.(type) {
	case string:
		return strings.TrimSpace(v)
	case []any:
		var parts []string
		for _, p := range v {
			pm, ok := p.(map[string]any)
			if !ok {
				continue
		placeholder
			if t, _ := pm["type"].(string); t != "text" {
				continue
		placeholder
			if text, ok := pm["text"].(string); ok && strings.TrimSpace(text) != "" {
				parts = append(parts, text)
		placeholder
	placeholder
		return strings.TrimSpace(strings.Join(parts, "\n"))
	default:
		return ""
placeholder
placeholder

func convertClaudeMessagesToGeminiContents(messages any, toolUseIDToName map[string]string) ([]any, error) {
	arr, ok := messages.([]any)
	if !ok {
		return nil, errors.New("messages must be an array")
placeholder

	out := make([]any, 0, len(arr))
	for _, m := range arr {
		mm, ok := m.(map[string]any)
		if !ok {
			continue
	placeholder
		role, _ := mm["role"].(string)
		role = strings.ToLower(strings.TrimSpace(role))
		gRole := "user"
		if role == "assistant" {
			gRole = "model"
	placeholder

		parts := make([]any, 0)
		switch content := mm["content"].(type) {
		case string:
			if strings.TrimSpace(content) != "" {
				parts = append(parts, map[string]any{"text": contentplaceholder)
		placeholder
		case []any:
			for _, block := range content {
				bm, ok := block.(map[string]any)
				if !ok {
					continue
			placeholder
				bt, _ := bm["type"].(string)
				switch bt {
				case "text":
					if text, ok := bm["text"].(string); ok && strings.TrimSpace(text) != "" {
						parts = append(parts, map[string]any{"text": textplaceholder)
				placeholder
				case "tool_use":
					id, _ := bm["id"].(string)
					name, _ := bm["name"].(string)
					if strings.TrimSpace(id) != "" && strings.TrimSpace(name) != "" {
						toolUseIDToName[id] = name
				placeholder
					parts = append(parts, map[string]any{
						"functionCall": map[string]any{
							"name": name,
							"args": bm["input"],
					placeholder,
				placeholder)
				case "tool_result":
					toolUseID, _ := bm["tool_use_id"].(string)
					name := toolUseIDToName[toolUseID]
					if name == "" {
						name = "tool"
				placeholder
					parts = append(parts, map[string]any{
						"functionResponse": map[string]any{
							"name": name,
							"response": map[string]any{
								"content": extractClaudeContentText(bm["content"]),
						placeholder,
					placeholder,
				placeholder)
				case "image":
					if src, ok := bm["source"].(map[string]any); ok {
						if srcType, _ := src["type"].(string); srcType == "base64" {
							mediaType, _ := src["media_type"].(string)
							data, _ := src["data"].(string)
							if mediaType != "" && data != "" {
								parts = append(parts, map[string]any{
									"inlineData": map[string]any{
										"mimeType": mediaType,
										"data":     data,
								placeholder,
							placeholder)
						placeholder
					placeholder
				placeholder
				default:
					// best-effort: preserve unknown blocks as text
					if b, err := json.Marshal(bm); err == nil {
						parts = append(parts, map[string]any{"text": string(b)placeholder)
				placeholder
			placeholder
		placeholder
		default:
			// ignore
	placeholder

		out = append(out, map[string]any{
			"role":  gRole,
			"parts": parts,
	placeholder)
placeholder
	return out, nil
placeholder

func extractClaudeContentText(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []any:
		var sb strings.Builder
		for _, part := range t {
			pm, ok := part.(map[string]any)
			if !ok {
				continue
		placeholder
			if pm["type"] == "text" {
				if text, ok := pm["text"].(string); ok {
					_, _ = sb.WriteString(text)
			placeholder
		placeholder
	placeholder
		return sb.String()
	default:
		b, _ := json.Marshal(t)
		return string(b)
placeholder
placeholder

func convertClaudeToolsToGeminiTools(tools any) []any {
	arr, ok := tools.([]any)
	if !ok || len(arr) == 0 {
		return nil
placeholder

	funcDecls := make([]any, 0, len(arr))
	for _, t := range arr {
		tm, ok := t.(map[string]any)
		if !ok {
			continue
	placeholder

		var name, desc string
		var params any

		// 检查是否为 custom 类型工具 (MCP)
		toolType, _ := tm["type"].(string)
		if toolType == "custom" {
			// Custom 格式: 从 custom 字段获取 description 和 input_schema
			custom, ok := tm["custom"].(map[string]any)
			if !ok {
				continue
		placeholder
			name, _ = tm["name"].(string)
			desc, _ = custom["description"].(string)
			params = custom["input_schema"]
	placeholder else {
			// 标准格式: 从顶层字段获取
			name, _ = tm["name"].(string)
			desc, _ = tm["description"].(string)
			params = tm["input_schema"]
	placeholder

		if name == "" {
			continue
	placeholder

		// 为 nil params 提供默认值
		if params == nil {
			params = map[string]any{
				"type":       "object",
				"properties": map[string]any{placeholder,
		placeholder
	placeholder
		// 清理 JSON Schema
		cleanedParams := cleanToolSchema(params)

		funcDecls = append(funcDecls, map[string]any{
			"name":        name,
			"description": desc,
			"parameters":  cleanedParams,
	placeholder)
placeholder

	if len(funcDecls) == 0 {
		return nil
placeholder
	return []any{
		map[string]any{
			"functionDeclarations": funcDecls,
	placeholder,
placeholder
placeholder

// cleanToolSchema 清理工具的 JSON Schema，移除 Gemini 不支持的字段
func cleanToolSchema(schema any) any {
	if schema == nil {
		return nil
placeholder

	switch v := schema.(type) {
	case map[string]any:
		cleaned := make(map[string]any)
		for key, value := range v {
			// 跳过不支持的字段
			if key == "$schema" || key == "$id" || key == "$ref" ||
				key == "additionalProperties" || key == "minLength" ||
				key == "maxLength" || key == "minItems" || key == "maxItems" {
				continue
		placeholder
			// 递归清理嵌套对象
			cleaned[key] = cleanToolSchema(value)
	placeholder
		// 规范化 type 字段为大写
		if typeVal, ok := cleaned["type"].(string); ok {
			cleaned["type"] = strings.ToUpper(typeVal)
	placeholder
		return cleaned
	case []any:
		cleaned := make([]any, len(v))
		for i, item := range v {
			cleaned[i] = cleanToolSchema(item)
	placeholder
		return cleaned
	default:
		return v
placeholder
placeholder

func convertClaudeGenerationConfig(req map[string]any) map[string]any {
	out := make(map[string]any)
	if mt, ok := asInt(req["max_tokens"]); ok && mt > 0 {
		out["maxOutputTokens"] = mt
placeholder
	if temp, ok := req["temperature"].(float64); ok {
		out["temperature"] = temp
placeholder
	if topP, ok := req["top_p"].(float64); ok {
		out["topP"] = topP
placeholder
	if stopSeq, ok := req["stop_sequences"].([]any); ok && len(stopSeq) > 0 {
		out["stopSequences"] = stopSeq
placeholder
	if len(out) == 0 {
		return nil
placeholder
	return out
placeholder
