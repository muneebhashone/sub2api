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

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/pkg/googleapi"
	"github.com/Wei-Shaw/sub2api/internal/util/responseheaders"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"

	"github.com/gin-gonic/gin"
)

const geminiStickySessionTTL = time.Hour

const (
	geminiMaxRetries     = 5
	geminiRetryBaseDelay = 1 * time.Second
	geminiRetryMaxDelay  = 16 * time.Second
)

// Gemini tool calling now requires `thoughtSignature` in parts that include `functionCall`.
// Many clients don't send it; we inject a known dummy signature to satisfy the validator.
// Ref: https://ai.google.dev/gemini-api/docs/thought-signatures
const geminiDummyThoughtSignature = "skip_thought_signature_validator"

type GeminiMessagesCompatService struct {
	accountRepo               AccountRepository
	groupRepo                 GroupRepository
	cache                     GatewayCache
	schedulerSnapshot         *SchedulerSnapshotService
	tokenProvider             *GeminiTokenProvider
	rateLimitService          *RateLimitService
	httpUpstream              HTTPUpstream
	antigravityGatewayService *AntigravityGatewayService
	cfg                       *config.Config
placeholder

func NewGeminiMessagesCompatService(
	accountRepo AccountRepository,
	groupRepo GroupRepository,
	cache GatewayCache,
	schedulerSnapshot *SchedulerSnapshotService,
	tokenProvider *GeminiTokenProvider,
	rateLimitService *RateLimitService,
	httpUpstream HTTPUpstream,
	antigravityGatewayService *AntigravityGatewayService,
	cfg *config.Config,
) *GeminiMessagesCompatService {
	return &GeminiMessagesCompatService{
		accountRepo:               accountRepo,
		groupRepo:                 groupRepo,
		cache:                     cache,
		schedulerSnapshot:         schedulerSnapshot,
		tokenProvider:             tokenProvider,
		rateLimitService:          rateLimitService,
		httpUpstream:              httpUpstream,
		antigravityGatewayService: antigravityGatewayService,
		cfg:                       cfg,
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
	// 1. 确定目标平台和调度模式
	// Determine target platform and scheduling mode
	platform, useMixedScheduling, hasForcePlatform, err := s.resolvePlatformAndSchedulingMode(ctx, groupID)
	if err != nil {
		return nil, err
placeholder

	cacheKey := "gemini:" + sessionHash

	// 2. 尝试粘性会话命中
	// Try sticky session hit
	if account := s.tryStickySessionHit(ctx, groupID, sessionHash, cacheKey, requestedModel, excludedIDs, platform, useMixedScheduling); account != nil {
		return account, nil
placeholder

	// 3. 查询可调度账户（强制平台模式：优先按分组查找，找不到再查全部）
	// Query schedulable accounts (force platform mode: try group first, fallback to all)
	accounts, err := s.listSchedulableAccountsOnce(ctx, groupID, platform, hasForcePlatform)
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder
	// 强制平台模式下，分组中找不到账户时回退查询全部
	if len(accounts) == 0 && groupID != nil && hasForcePlatform {
		accounts, err = s.listSchedulableAccountsOnce(ctx, nil, platform, hasForcePlatform)
		if err != nil {
			return nil, fmt.Errorf("query accounts failed: %w", err)
	placeholder
placeholder

	// 4. 按优先级 + LRU 选择最佳账号
	// Select best account by priority + LRU
	selected := s.selectBestGeminiAccount(ctx, accounts, requestedModel, excludedIDs, platform, useMixedScheduling)

	if selected == nil {
		if requestedModel != "" {
			return nil, fmt.Errorf("no available Gemini accounts supporting model: %s", requestedModel)
	placeholder
		return nil, errors.New("no available Gemini accounts")
placeholder

	// 5. 设置粘性会话绑定
	// Set sticky session binding
	if sessionHash != "" {
		_ = s.cache.SetSessionAccountID(ctx, derefGroupID(groupID), cacheKey, selected.ID, geminiStickySessionTTL)
placeholder

	return selected, nil
placeholder

// resolvePlatformAndSchedulingMode 解析目标平台和调度模式。
// 返回：平台名称、是否使用混合调度、是否强制平台、错误。
//
// resolvePlatformAndSchedulingMode resolves target platform and scheduling mode.
// Returns: platform name, whether to use mixed scheduling, whether force platform, error.
func (s *GeminiMessagesCompatService) resolvePlatformAndSchedulingMode(ctx context.Context, groupID *int64) (platform string, useMixedScheduling bool, hasForcePlatform bool, err error) {
	// 优先检查 context 中的强制平台（/antigravity 路由）
	forcePlatform, hasForcePlatform := ctx.Value(ctxkey.ForcePlatform).(string)
	if hasForcePlatform && forcePlatform != "" {
		return forcePlatform, false, true, nil
placeholder

	if groupID != nil {
		// 根据分组 platform 决定查询哪种账号
		var group *Group
		if ctxGroup, ok := ctx.Value(ctxkey.Group).(*Group); ok && IsGroupContextValid(ctxGroup) && ctxGroup.ID == *groupID {
			group = ctxGroup
	placeholder else {
			group, err = s.groupRepo.GetByIDLite(ctx, *groupID)
			if err != nil {
				return "", false, false, fmt.Errorf("get group failed: %w", err)
		placeholder
	placeholder
		// gemini 分组支持混合调度（包含启用了 mixed_scheduling 的 antigravity 账户）
		return group.Platform, group.Platform == PlatformGemini, false, nil
placeholder

	// 无分组时只使用原生 gemini 平台
	return PlatformGemini, true, false, nil
placeholder

// tryStickySessionHit 尝试从粘性会话获取账号。
// 如果命中且账号可用则返回账号；如果账号不可用则清理会话并返回 nil。
//
// tryStickySessionHit attempts to get account from sticky session.
// Returns account if hit and usable; clears session and returns nil if account unavailable.
func (s *GeminiMessagesCompatService) tryStickySessionHit(
	ctx context.Context,
	groupID *int64,
	sessionHash, cacheKey, requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	platform string,
	useMixedScheduling bool,
) *Account {
	if sessionHash == "" {
		return nil
placeholder

	accountID, err := s.cache.GetSessionAccountID(ctx, derefGroupID(groupID), cacheKey)
	if err != nil || accountID <= 0 {
		return nil
placeholder

	if _, excluded := excludedIDs[accountID]; excluded {
		return nil
placeholder

	account, err := s.getSchedulableAccount(ctx, accountID)
	if err != nil {
		return nil
placeholder

	// 检查账号是否需要清理粘性会话
	// Check if sticky session should be cleared
	if shouldClearStickySession(account) {
		_ = s.cache.DeleteSessionAccountID(ctx, derefGroupID(groupID), cacheKey)
		return nil
placeholder

	// 验证账号是否可用于当前请求
	// Verify account is usable for current request
	if !s.isAccountUsableForRequest(ctx, account, requestedModel, platform, useMixedScheduling) {
		return nil
placeholder

	// 刷新会话 TTL 并返回账号
	// Refresh session TTL and return account
	_ = s.cache.RefreshSessionTTL(ctx, derefGroupID(groupID), cacheKey, geminiStickySessionTTL)
	return account
placeholder

// isAccountUsableForRequest 检查账号是否可用于当前请求。
// 验证：模型调度、模型支持、平台匹配、速率限制预检。
//
// isAccountUsableForRequest checks if account is usable for current request.
// Validates: model scheduling, model support, platform matching, rate limit precheck.
func (s *GeminiMessagesCompatService) isAccountUsableForRequest(
	ctx context.Context,
	account *Account,
	requestedModel, platform string,
	useMixedScheduling bool,
) bool {
	// 检查模型调度能力
	// Check model scheduling capability
	if !account.IsSchedulableForModel(requestedModel) {
		return false
placeholder

	// 检查模型支持
	// Check model support
	if requestedModel != "" && !s.isModelSupportedByAccount(account, requestedModel) {
		return false
placeholder

	// 检查平台匹配
	// Check platform matching
	if !s.isAccountValidForPlatform(account, platform, useMixedScheduling) {
		return false
placeholder

	// 速率限制预检
	// Rate limit precheck
	if !s.passesRateLimitPreCheck(ctx, account, requestedModel) {
		return false
placeholder

	return true
placeholder

// isAccountValidForPlatform 检查账号是否匹配目标平台。
// 原生平台直接匹配；混合调度模式下 antigravity 需要启用 mixed_scheduling。
//
// isAccountValidForPlatform checks if account matches target platform.
// Native platform matches directly; mixed scheduling mode requires antigravity to enable mixed_scheduling.
func (s *GeminiMessagesCompatService) isAccountValidForPlatform(account *Account, platform string, useMixedScheduling bool) bool {
	if account.Platform == platform {
		return true
placeholder
	if useMixedScheduling && account.Platform == PlatformAntigravity && account.IsMixedSchedulingEnabled() {
		return true
placeholder
	return false
placeholder

// passesRateLimitPreCheck 执行速率限制预检。
// 返回 true 表示通过预检或无需预检。
//
// passesRateLimitPreCheck performs rate limit precheck.
// Returns true if passed or precheck not required.
func (s *GeminiMessagesCompatService) passesRateLimitPreCheck(ctx context.Context, account *Account, requestedModel string) bool {
	if s.rateLimitService == nil || requestedModel == "" {
		return true
placeholder
	ok, err := s.rateLimitService.PreCheckUsage(ctx, account, requestedModel)
	if err != nil {
		log.Printf("[Gemini PreCheck] Account %d precheck error: %v", account.ID, err)
placeholder
	return ok
placeholder

// selectBestGeminiAccount 从候选账号中选择最佳账号（优先级 + LRU + OAuth 优先）。
// 返回 nil 表示无可用账号。
//
// selectBestGeminiAccount selects best account from candidates (priority + LRU + OAuth preferred).
// Returns nil if no available account.
func (s *GeminiMessagesCompatService) selectBestGeminiAccount(
	ctx context.Context,
	accounts []Account,
	requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	platform string,
	useMixedScheduling bool,
) *Account {
	var selected *Account

	for i := range accounts {
		acc := &accounts[i]

		// 跳过被排除的账号
		if _, excluded := excludedIDs[acc.ID]; excluded {
			continue
	placeholder

		// 检查账号是否可用于当前请求
		if !s.isAccountUsableForRequest(ctx, acc, requestedModel, platform, useMixedScheduling) {
			continue
	placeholder

		// 选择最佳账号
		if selected == nil {
			selected = acc
			continue
	placeholder

		if s.isBetterGeminiAccount(acc, selected) {
			selected = acc
	placeholder
placeholder

	return selected
placeholder

// isBetterGeminiAccount 判断 candidate 是否比 current 更优。
// 规则：优先级更高（数值更小）优先；同优先级时，未使用过的优先（OAuth > 非 OAuth），其次是最久未使用的。
//
// isBetterGeminiAccount checks if candidate is better than current.
// Rules: higher priority (lower value) wins; same priority: never used (OAuth > non-OAuth) > least recently used.
func (s *GeminiMessagesCompatService) isBetterGeminiAccount(candidate, current *Account) bool {
	// 优先级更高（数值更小）
	if candidate.Priority < current.Priority {
		return true
placeholder
	if candidate.Priority > current.Priority {
		return false
placeholder

	// 同优先级，比较最后使用时间
	switch {
	case candidate.LastUsedAt == nil && current.LastUsedAt != nil:
		// candidate 从未使用，优先
		return true
	case candidate.LastUsedAt != nil && current.LastUsedAt == nil:
		// current 从未使用，保持
		return false
	case candidate.LastUsedAt == nil && current.LastUsedAt == nil:
		// 都未使用，优先选择 OAuth 账号（更兼容 Code Assist 流程）
		return candidate.Type == AccountTypeOAuth && current.Type != AccountTypeOAuth
	default:
		// 都使用过，选择最久未使用的
		return candidate.LastUsedAt.Before(*current.LastUsedAt)
placeholder
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

func (s *GeminiMessagesCompatService) getSchedulableAccount(ctx context.Context, accountID int64) (*Account, error) {
	if s.schedulerSnapshot != nil {
		return s.schedulerSnapshot.GetAccount(ctx, accountID)
placeholder
	return s.accountRepo.GetByID(ctx, accountID)
placeholder

func (s *GeminiMessagesCompatService) listSchedulableAccountsOnce(ctx context.Context, groupID *int64, platform string, hasForcePlatform bool) ([]Account, error) {
	if s.schedulerSnapshot != nil {
		accounts, _, err := s.schedulerSnapshot.ListSchedulableAccounts(ctx, groupID, platform, hasForcePlatform)
		return accounts, err
placeholder

	useMixedScheduling := platform == PlatformGemini && !hasForcePlatform
	queryPlatforms := []string{platformplaceholder
	if useMixedScheduling {
		queryPlatforms = []string{platform, PlatformAntigravityplaceholder
placeholder

	if groupID != nil {
		return s.accountRepo.ListSchedulableByGroupIDAndPlatforms(ctx, *groupID, queryPlatforms)
placeholder
	return s.accountRepo.ListSchedulableByPlatforms(ctx, queryPlatforms)
placeholder

func (s *GeminiMessagesCompatService) validateUpstreamBaseURL(raw string) (string, error) {
	if s.cfg != nil && !s.cfg.Security.URLAllowlist.Enabled {
		normalized, err := urlvalidator.ValidateURLFormat(raw, s.cfg.Security.URLAllowlist.AllowInsecureHTTP)
		if err != nil {
			return "", fmt.Errorf("invalid base_url: %w", err)
	placeholder
		return normalized, nil
placeholder
	normalized, err := urlvalidator.ValidateHTTPSURL(raw, urlvalidator.ValidationOptions{
		AllowedHosts:     s.cfg.Security.URLAllowlist.UpstreamHosts,
		RequireAllowlist: true,
		AllowPrivate:     s.cfg.Security.URLAllowlist.AllowPrivateHosts,
placeholder)
	if err != nil {
		return "", fmt.Errorf("invalid base_url: %w", err)
placeholder
	return normalized, nil
placeholder

// HasAntigravityAccounts 检查是否有可用的 antigravity 账户
func (s *GeminiMessagesCompatService) HasAntigravityAccounts(ctx context.Context, groupID *int64) (bool, error) {
	accounts, err := s.listSchedulableAccountsOnce(ctx, groupID, PlatformAntigravity, false)
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
	accounts, err := s.listSchedulableAccountsOnce(ctx, groupID, PlatformGemini, true)
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
		case AccountTypeAPIKey:
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
	if account.Type == AccountTypeAPIKey {
		mappedModel = account.GetMappedModel(req.Model)
placeholder

	geminiReq, err := convertClaudeMessagesToGeminiGenerateContent(body)
	if err != nil {
		return nil, s.writeClaudeError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
placeholder
	geminiReq = ensureGeminiFunctionCallThoughtSignatures(geminiReq)
	originalClaudeBody := body

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
	case AccountTypeAPIKey:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			apiKey := account.GetCredential("api_key")
			if strings.TrimSpace(apiKey) == "" {
				return nil, "", errors.New("gemini api_key not configured")
		placeholder

			baseURL := strings.TrimSpace(account.GetCredential("base_url"))
			if baseURL == "" {
				baseURL = geminicli.AIStudioBaseURL
		placeholder
			normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
			if err != nil {
				return nil, "", err
		placeholder

			action := "generateContent"
			if req.Stream {
				action = "streamGenerateContent"
		placeholder
			fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(normalizedBaseURL, "/"), mappedModel, action)
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
				baseURL, err := s.validateUpstreamBaseURL(geminicli.GeminiCliBaseURL)
				if err != nil {
					return nil, "", err
			placeholder
				fullURL := fmt.Sprintf("%s/v1internal:%s", strings.TrimRight(baseURL, "/"), action)
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
				baseURL := strings.TrimSpace(account.GetCredential("base_url"))
				if baseURL == "" {
					baseURL = geminicli.AIStudioBaseURL
			placeholder
				normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
				if err != nil {
					return nil, "", err
			placeholder

				fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(normalizedBaseURL, "/"), mappedModel, action)
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
	signatureRetryStage := 0
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

		// Capture upstream request body for ops retry of this attempt.
		if c != nil {
			// In this code path `body` is already the JSON sent to upstream.
			c.Set(OpsUpstreamRequestBodyKey, string(body))
	placeholder

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		if err != nil {
			safeErr := sanitizeUpstreamErrorMessage(err.Error())
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: 0,
				Kind:               "request_error",
				Message:            safeErr,
		placeholder)
			if attempt < geminiMaxRetries {
				log.Printf("Gemini account %d: upstream request failed, retry %d/%d: %v", account.ID, attempt, geminiMaxRetries, err)
				sleepGeminiBackoff(attempt)
				continue
		placeholder
			setOpsUpstreamError(c, 0, safeErr, "")
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed after retries: "+safeErr)
	placeholder

		// Special-case: signature/thought_signature validation errors are not transient, but may be fixed by
		// downgrading Claude thinking/tool history to plain text (conservative two-stage retry).
		if resp.StatusCode == http.StatusBadRequest && signatureRetryStage < 2 {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()

			if isGeminiSignatureRelatedError(respBody) {
				upstreamReqID := resp.Header.Get(requestIDHeader)
				if upstreamReqID == "" {
					upstreamReqID = resp.Header.Get("x-goog-request-id")
			placeholder
				upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
				upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
				upstreamDetail := ""
				if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
					maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
					if maxBytes <= 0 {
						maxBytes = 2048
				placeholder
					upstreamDetail = truncateString(string(respBody), maxBytes)
			placeholder
				appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
					Platform:           account.Platform,
					AccountID:          account.ID,
					AccountName:        account.Name,
					UpstreamStatusCode: resp.StatusCode,
					UpstreamRequestID:  upstreamReqID,
					Kind:               "signature_error",
					Message:            upstreamMsg,
					Detail:             upstreamDetail,
			placeholder)

				var strippedClaudeBody []byte
				stageName := ""
				switch signatureRetryStage {
				case 0:
					// Stage 1: disable thinking + thinking->text
					strippedClaudeBody = FilterThinkingBlocksForRetry(originalClaudeBody)
					stageName = "thinking-only"
					signatureRetryStage = 1
				default:
					// Stage 2: additionally downgrade tool_use/tool_result blocks to text
					strippedClaudeBody = FilterSignatureSensitiveBlocksForRetry(originalClaudeBody)
					stageName = "thinking+tools"
					signatureRetryStage = 2
			placeholder
				retryGeminiReq, txErr := convertClaudeMessagesToGeminiGenerateContent(strippedClaudeBody)
				if txErr == nil {
					log.Printf("Gemini account %d: detected signature-related 400, retrying with downgraded Claude blocks (%s)", account.ID, stageName)
					geminiReq = retryGeminiReq
					// Consume one retry budget attempt and continue with the updated request payload.
					sleepGeminiBackoff(1)
					continue
			placeholder
		placeholder

			// Restore body for downstream error handling.
			resp = &http.Response{
				StatusCode: http.StatusBadRequest,
				Header:     resp.Header.Clone(),
				Body:       io.NopCloser(bytes.NewReader(respBody)),
		placeholder
			break
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
				upstreamReqID := resp.Header.Get(requestIDHeader)
				if upstreamReqID == "" {
					upstreamReqID = resp.Header.Get("x-goog-request-id")
			placeholder
				upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
				upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
				upstreamDetail := ""
				if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
					maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
					if maxBytes <= 0 {
						maxBytes = 2048
				placeholder
					upstreamDetail = truncateString(string(respBody), maxBytes)
			placeholder
				appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
					Platform:           account.Platform,
					AccountID:          account.ID,
					AccountName:        account.Name,
					UpstreamStatusCode: resp.StatusCode,
					UpstreamRequestID:  upstreamReqID,
					Kind:               "retry",
					Message:            upstreamMsg,
					Detail:             upstreamDetail,
			placeholder)

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
		tempMatched := false
		if s.rateLimitService != nil {
			tempMatched = s.rateLimitService.HandleTempUnschedulable(ctx, account, resp.StatusCode, respBody)
	placeholder
		s.handleGeminiUpstreamError(ctx, account, resp.StatusCode, resp.Header, respBody)
		if tempMatched {
			upstreamReqID := resp.Header.Get(requestIDHeader)
			if upstreamReqID == "" {
				upstreamReqID = resp.Header.Get("x-goog-request-id")
		placeholder
			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			upstreamDetail := ""
			if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
				maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
				if maxBytes <= 0 {
					maxBytes = 2048
			placeholder
				upstreamDetail = truncateString(string(respBody), maxBytes)
		placeholder
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  upstreamReqID,
				Kind:               "failover",
				Message:            upstreamMsg,
				Detail:             upstreamDetail,
		placeholder)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode, ResponseBody: respBodyplaceholder
	placeholder
		if s.shouldFailoverGeminiUpstreamError(resp.StatusCode) {
			upstreamReqID := resp.Header.Get(requestIDHeader)
			if upstreamReqID == "" {
				upstreamReqID = resp.Header.Get("x-goog-request-id")
		placeholder
			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			upstreamDetail := ""
			if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
				maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
				if maxBytes <= 0 {
					maxBytes = 2048
			placeholder
				upstreamDetail = truncateString(string(respBody), maxBytes)
		placeholder
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  upstreamReqID,
				Kind:               "failover",
				Message:            upstreamMsg,
				Detail:             upstreamDetail,
		placeholder)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode, ResponseBody: respBodyplaceholder
	placeholder
		upstreamReqID := resp.Header.Get(requestIDHeader)
		if upstreamReqID == "" {
			upstreamReqID = resp.Header.Get("x-goog-request-id")
	placeholder
		return nil, s.writeGeminiMappedError(c, account, resp.StatusCode, upstreamReqID, respBody)
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

	// 图片生成计费
	imageCount := 0
	imageSize := s.extractImageSize(body)
	if isImageGenerationModel(originalModel) {
		imageCount = 1
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel,
		Stream:       req.Stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		ImageCount:   imageCount,
		ImageSize:    imageSize,
placeholder, nil
placeholder

func isGeminiSignatureRelatedError(respBody []byte) bool {
	msg := strings.ToLower(strings.TrimSpace(extractAntigravityErrorMessage(respBody)))
	if msg == "" {
		msg = strings.ToLower(string(respBody))
placeholder
	return strings.Contains(msg, "thought_signature") || strings.Contains(msg, "signature")
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

	// 过滤掉 parts 为空的消息（Gemini API 不接受空 parts）
	if filteredBody, err := filterEmptyPartsFromGeminiRequest(body); err == nil {
		body = filteredBody
placeholder

	switch action {
	case "generateContent", "streamGenerateContent", "countTokens":
		// ok
	default:
		return nil, s.writeGoogleError(c, http.StatusNotFound, "Unsupported action: "+action)
placeholder

	// Some Gemini upstreams validate tool call parts strictly; ensure any `functionCall` part includes a
	// `thoughtSignature` to avoid frequent INVALID_ARGUMENT 400s.
	body = ensureGeminiFunctionCallThoughtSignatures(body)

	mappedModel := originalModel
	if account.Type == AccountTypeAPIKey {
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
	case AccountTypeAPIKey:
		buildReq = func(ctx context.Context) (*http.Request, string, error) {
			apiKey := account.GetCredential("api_key")
			if strings.TrimSpace(apiKey) == "" {
				return nil, "", errors.New("gemini api_key not configured")
		placeholder

			baseURL := strings.TrimSpace(account.GetCredential("base_url"))
			if baseURL == "" {
				baseURL = geminicli.AIStudioBaseURL
		placeholder
			normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
			if err != nil {
				return nil, "", err
		placeholder

			fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(normalizedBaseURL, "/"), mappedModel, upstreamAction)
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
				baseURL, err := s.validateUpstreamBaseURL(geminicli.GeminiCliBaseURL)
				if err != nil {
					return nil, "", err
			placeholder
				fullURL := fmt.Sprintf("%s/v1internal:%s", strings.TrimRight(baseURL, "/"), upstreamAction)
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
				baseURL := strings.TrimSpace(account.GetCredential("base_url"))
				if baseURL == "" {
					baseURL = geminicli.AIStudioBaseURL
			placeholder
				normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
				if err != nil {
					return nil, "", err
			placeholder

				fullURL := fmt.Sprintf("%s/v1beta/models/%s:%s", strings.TrimRight(normalizedBaseURL, "/"), mappedModel, upstreamAction)
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

		// Capture upstream request body for ops retry of this attempt.
		if c != nil {
			// In this code path `body` is already the JSON sent to upstream.
			c.Set(OpsUpstreamRequestBodyKey, string(body))
	placeholder

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		if err != nil {
			safeErr := sanitizeUpstreamErrorMessage(err.Error())
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: 0,
				Kind:               "request_error",
				Message:            safeErr,
		placeholder)
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
			setOpsUpstreamError(c, 0, safeErr, "")
			return nil, s.writeGoogleError(c, http.StatusBadGateway, "Upstream request failed after retries: "+safeErr)
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
				upstreamReqID := resp.Header.Get(requestIDHeader)
				if upstreamReqID == "" {
					upstreamReqID = resp.Header.Get("x-goog-request-id")
			placeholder
				upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
				upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
				upstreamDetail := ""
				if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
					maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
					if maxBytes <= 0 {
						maxBytes = 2048
				placeholder
					upstreamDetail = truncateString(string(respBody), maxBytes)
			placeholder
				appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
					Platform:           account.Platform,
					AccountID:          account.ID,
					AccountName:        account.Name,
					UpstreamStatusCode: resp.StatusCode,
					UpstreamRequestID:  upstreamReqID,
					Kind:               "retry",
					Message:            upstreamMsg,
					Detail:             upstreamDetail,
			placeholder)

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
		tempMatched := false
		if s.rateLimitService != nil {
			tempMatched = s.rateLimitService.HandleTempUnschedulable(ctx, account, resp.StatusCode, respBody)
	placeholder
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

		if tempMatched {
			evBody := unwrapIfNeeded(isOAuth, respBody)
			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(evBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			upstreamDetail := ""
			if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
				maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
				if maxBytes <= 0 {
					maxBytes = 2048
			placeholder
				upstreamDetail = truncateString(string(evBody), maxBytes)
		placeholder
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  requestID,
				Kind:               "failover",
				Message:            upstreamMsg,
				Detail:             upstreamDetail,
		placeholder)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode, ResponseBody: respBodyplaceholder
	placeholder
		if s.shouldFailoverGeminiUpstreamError(resp.StatusCode) {
			evBody := unwrapIfNeeded(isOAuth, respBody)
			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(evBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			upstreamDetail := ""
			if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
				maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
				if maxBytes <= 0 {
					maxBytes = 2048
			placeholder
				upstreamDetail = truncateString(string(evBody), maxBytes)
		placeholder
			appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
				Platform:           account.Platform,
				AccountID:          account.ID,
				AccountName:        account.Name,
				UpstreamStatusCode: resp.StatusCode,
				UpstreamRequestID:  requestID,
				Kind:               "failover",
				Message:            upstreamMsg,
				Detail:             upstreamDetail,
		placeholder)
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCode, ResponseBody: evBodyplaceholder
	placeholder

		respBody = unwrapIfNeeded(isOAuth, respBody)
		upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
		upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
		upstreamDetail := ""
		if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
			maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
			if maxBytes <= 0 {
				maxBytes = 2048
		placeholder
			upstreamDetail = truncateString(string(respBody), maxBytes)
			log.Printf("[Gemini] native upstream error %d: %s", resp.StatusCode, truncateForLog(respBody, s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes))
	placeholder
		setOpsUpstreamError(c, resp.StatusCode, upstreamMsg, upstreamDetail)
		appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
			Platform:           account.Platform,
			AccountID:          account.ID,
			AccountName:        account.Name,
			UpstreamStatusCode: resp.StatusCode,
			UpstreamRequestID:  requestID,
			Kind:               "http_error",
			Message:            upstreamMsg,
			Detail:             upstreamDetail,
	placeholder)

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
	placeholder
		c.Data(resp.StatusCode, contentType, respBody)
		if upstreamMsg == "" {
			return nil, fmt.Errorf("gemini upstream error: %d", resp.StatusCode)
	placeholder
		return nil, fmt.Errorf("gemini upstream error: %d message=%s", resp.StatusCode, upstreamMsg)
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

	// 图片生成计费
	imageCount := 0
	imageSize := s.extractImageSize(body)
	if isImageGenerationModel(originalModel) {
		imageCount = 1
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel,
		Stream:       stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		ImageCount:   imageCount,
		ImageSize:    imageSize,
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

func (s *GeminiMessagesCompatService) writeGeminiMappedError(c *gin.Context, account *Account, upstreamStatus int, upstreamRequestID string, body []byte) error {
	upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(body))
	upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
	upstreamDetail := ""
	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
		if maxBytes <= 0 {
			maxBytes = 2048
	placeholder
		upstreamDetail = truncateString(string(body), maxBytes)
placeholder
	setOpsUpstreamError(c, upstreamStatus, upstreamMsg, upstreamDetail)
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:           account.Platform,
		AccountID:          account.ID,
		AccountName:        account.Name,
		UpstreamStatusCode: upstreamStatus,
		UpstreamRequestID:  upstreamRequestID,
		Kind:               "http_error",
		Message:            upstreamMsg,
		Detail:             upstreamDetail,
placeholder)

	if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		log.Printf("[Gemini] upstream error %d: %s", upstreamStatus, truncateForLog(body, s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes))
placeholder

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
	if upstreamMsg == "" {
		return fmt.Errorf("upstream error: %d", upstreamStatus)
placeholder
	return fmt.Errorf("upstream error: %d message=%s", upstreamStatus, upstreamMsg)
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
	var collectedTextParts []string // Collect all text parts for aggregation
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
						return mergeCollectedTextParts(pickGeminiCollectResult(last, lastWithParts), collectedTextParts), usage, nil
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
							// Collect text from each part for aggregation
							for _, part := range parts {
								if text, ok := part["text"].(string); ok && text != "" {
									collectedTextParts = append(collectedTextParts, text)
							placeholder
						placeholder
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

	return mergeCollectedTextParts(pickGeminiCollectResult(last, lastWithParts), collectedTextParts), usage, nil
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

// mergeCollectedTextParts merges all collected text chunks into the final response.
// This fixes the issue where non-streaming responses only returned the last chunk
// instead of the complete aggregated text.
func mergeCollectedTextParts(response map[string]any, textParts []string) map[string]any {
	if len(textParts) == 0 {
		return response
placeholder

	// Join all text parts
	mergedText := strings.Join(textParts, "")

	// Deep copy response
	result := make(map[string]any)
	for k, v := range response {
		result[k] = v
placeholder

	// Get or create candidates
	candidates, ok := result["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		candidates = []any{map[string]any{placeholderplaceholder
placeholder

	// Get first candidate
	candidate, ok := candidates[0].(map[string]any)
	if !ok {
		candidate = make(map[string]any)
		candidates[0] = candidate
placeholder

	// Get or create content
	content, ok := candidate["content"].(map[string]any)
	if !ok {
		content = map[string]any{"role": "model"placeholder
		candidate["content"] = content
placeholder

	// Get existing parts
	existingParts, ok := content["parts"].([]any)
	if !ok {
		existingParts = []any{placeholder
placeholder

	// Find and update first text part, or create new one
	newParts := make([]any, 0, len(existingParts)+1)
	textUpdated := false

	for _, p := range existingParts {
		pm, ok := p.(map[string]any)
		if !ok {
			newParts = append(newParts, p)
			continue
	placeholder
		if _, hasText := pm["text"]; hasText && !textUpdated {
			// Replace with merged text
			newPart := make(map[string]any)
			for k, v := range pm {
				newPart[k] = v
		placeholder
			newPart["text"] = mergedText
			newParts = append(newParts, newPart)
			textUpdated = true
	placeholder else {
			newParts = append(newParts, pm)
	placeholder
placeholder

	if !textUpdated {
		newParts = append([]any{map[string]any{"text": mergedTextplaceholderplaceholder, newParts...)
placeholder

	content["parts"] = newParts
	result["candidates"] = candidates

	return result
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
	// Log response headers for debugging
	log.Printf("[GeminiAPI] ========== Response Headers ==========")
	for key, values := range resp.Header {
		if strings.HasPrefix(strings.ToLower(key), "x-ratelimit") {
			log.Printf("[GeminiAPI] %s: %v", key, values)
	placeholder
placeholder
	log.Printf("[GeminiAPI] ========================================")

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

	responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)

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
	// Log response headers for debugging
	log.Printf("[GeminiAPI] ========== Streaming Response Headers ==========")
	for key, values := range resp.Header {
		if strings.HasPrefix(strings.ToLower(key), "x-ratelimit") {
			log.Printf("[GeminiAPI] %s: %v", key, values)
	placeholder
placeholder
	log.Printf("[GeminiAPI] ====================================================")

	if s.cfg != nil {
		responseheaders.WriteFilteredHeaders(c.Writer.Header(), resp.Header, s.cfg.Security.ResponseHeaders)
placeholder

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

	baseURL := strings.TrimSpace(account.GetCredential("base_url"))
	if baseURL == "" {
		baseURL = geminicli.AIStudioBaseURL
placeholder
	normalizedBaseURL, err := s.validateUpstreamBaseURL(baseURL)
	if err != nil {
		return nil, err
placeholder
	fullURL := strings.TrimRight(normalizedBaseURL, "/") + path

	var proxyURL string
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
placeholder

	switch account.Type {
	case AccountTypeAPIKey:
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
	wwwAuthenticate := resp.Header.Get("Www-Authenticate")
	filteredHeaders := responseheaders.FilterHeaders(resp.Header, s.cfg.Security.ResponseHeaders)
	if wwwAuthenticate != "" {
		filteredHeaders.Set("Www-Authenticate", wwwAuthenticate)
placeholder
	return &UpstreamHTTPResult{
		StatusCode: resp.StatusCode,
		Headers:    filteredHeaders,
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
	cached, _ := asInt(usageMeta["cachedContentTokenCount"])
	// 注意：Gemini 的 promptTokenCount 包含 cachedContentTokenCount，
	// 但 Claude 的 input_tokens 不包含 cache_read_input_tokens，需要减去
	return &ClaudeUsage{
		InputTokens:          prompt - cached,
		OutputTokens:         cand,
		CacheReadInputTokens: cached,
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

func ensureGeminiFunctionCallThoughtSignatures(body []byte) []byte {
	// Fast path: only run when functionCall is present.
	if !bytes.Contains(body, []byte(`"functionCall"`)) {
		return body
placeholder

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return body
placeholder

	contentsAny, ok := payload["contents"].([]any)
	if !ok || len(contentsAny) == 0 {
		return body
placeholder

	modified := false
	for _, c := range contentsAny {
		cm, ok := c.(map[string]any)
		if !ok {
			continue
	placeholder
		partsAny, ok := cm["parts"].([]any)
		if !ok || len(partsAny) == 0 {
			continue
	placeholder
		for _, p := range partsAny {
			pm, ok := p.(map[string]any)
			if !ok || pm == nil {
				continue
		placeholder
			if fc, ok := pm["functionCall"].(map[string]any); !ok || fc == nil {
				continue
		placeholder
			ts, _ := pm["thoughtSignature"].(string)
			if strings.TrimSpace(ts) == "" {
				pm["thoughtSignature"] = geminiDummyThoughtSignature
				modified = true
		placeholder
	placeholder
placeholder

	if !modified {
		return body
placeholder
	b, err := json.Marshal(payload)
	if err != nil {
		return body
placeholder
	return b
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
			// 字符串形式的 content，保留所有内容（包括空白）
			parts = append(parts, map[string]any{"text": contentplaceholder)
		case []any:
			// 如果只有一个 block，不过滤空白（让上游 API 报错）
			singleBlock := len(content) == 1

			for _, block := range content {
				bm, ok := block.(map[string]any)
				if !ok {
					continue
			placeholder
				bt, _ := bm["type"].(string)
				switch bt {
				case "text":
					if text, ok := bm["text"].(string); ok {
						// 单个 block 时保留所有内容（包括空白）
						// 多个 blocks 时过滤掉空白
						if singleBlock || strings.TrimSpace(text) != "" {
							parts = append(parts, map[string]any{"text": textplaceholder)
					placeholder
				placeholder
				case "tool_use":
					id, _ := bm["id"].(string)
					name, _ := bm["name"].(string)
					if strings.TrimSpace(id) != "" && strings.TrimSpace(name) != "" {
						toolUseIDToName[id] = name
				placeholder
					signature, _ := bm["signature"].(string)
					signature = strings.TrimSpace(signature)
					if signature == "" {
						signature = geminiDummyThoughtSignature
				placeholder
					parts = append(parts, map[string]any{
						"thoughtSignature": signature,
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

// extractImageSize 从 Gemini 请求中提取 image_size 参数
func (s *GeminiMessagesCompatService) extractImageSize(body []byte) string {
	var req struct {
		GenerationConfig *struct {
			ImageConfig *struct {
				ImageSize string `json:"imageSize"`
		placeholder `json:"imageConfig"`
	placeholder `json:"generationConfig"`
placeholder
	if err := json.Unmarshal(body, &req); err != nil {
		return "2K"
placeholder

	if req.GenerationConfig != nil && req.GenerationConfig.ImageConfig != nil {
		size := strings.ToUpper(strings.TrimSpace(req.GenerationConfig.ImageConfig.ImageSize))
		if size == "1K" || size == "2K" || size == "4K" {
			return size
	placeholder
placeholder

	return "2K"
placeholder
