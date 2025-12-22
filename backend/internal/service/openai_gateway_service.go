package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sub2api/internal/config"
	"sub2api/internal/model"
	"sub2api/internal/service/ports"

	"github.com/gin-gonic/gin"
)

const (
	// ChatGPT internal API for OAuth accounts
	chatgptCodexURL = "https://chatgpt.com/backend-api/codex/responses"
	// OpenAI Platform API for API Key accounts (fallback)
	openaiPlatformAPIURL   = "https://api.openai.com/v1/responses"
	openaiStickySessionTTL = time.Hour // 粘性会话TTL
)

// OpenAI allowed headers whitelist (for non-OAuth accounts)
var openaiAllowedHeaders = map[string]bool{
	"accept-language": true,
	"content-type":    true,
	"user-agent":      true,
	"originator":      true,
	"session_id":      true,
placeholder

// OpenAIUsage represents OpenAI API response usage
type OpenAIUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
placeholder

// OpenAIForwardResult represents the result of forwarding
type OpenAIForwardResult struct {
	RequestID    string
	Usage        OpenAIUsage
	Model        string
	Stream       bool
	Duration     time.Duration
	FirstTokenMs *int
placeholder

// OpenAIGatewayService handles OpenAI API gateway operations
type OpenAIGatewayService struct {
	accountRepo         ports.AccountRepository
	usageLogRepo        ports.UsageLogRepository
	userRepo            ports.UserRepository
	userSubRepo         ports.UserSubscriptionRepository
	cache               ports.GatewayCache
	cfg                 *config.Config
	billingService      *BillingService
	rateLimitService    *RateLimitService
	billingCacheService *BillingCacheService
	httpUpstream        ports.HTTPUpstream
placeholder

// NewOpenAIGatewayService creates a new OpenAIGatewayService
func NewOpenAIGatewayService(
	accountRepo ports.AccountRepository,
	usageLogRepo ports.UsageLogRepository,
	userRepo ports.UserRepository,
	userSubRepo ports.UserSubscriptionRepository,
	cache ports.GatewayCache,
	cfg *config.Config,
	billingService *BillingService,
	rateLimitService *RateLimitService,
	billingCacheService *BillingCacheService,
	httpUpstream ports.HTTPUpstream,
) *OpenAIGatewayService {
	return &OpenAIGatewayService{
		accountRepo:         accountRepo,
		usageLogRepo:        usageLogRepo,
		userRepo:            userRepo,
		userSubRepo:         userSubRepo,
		cache:               cache,
		cfg:                 cfg,
		billingService:      billingService,
		rateLimitService:    rateLimitService,
		billingCacheService: billingCacheService,
		httpUpstream:        httpUpstream,
placeholder
placeholder

// GenerateSessionHash generates session hash from header (OpenAI uses session_id header)
func (s *OpenAIGatewayService) GenerateSessionHash(c *gin.Context) string {
	sessionID := c.GetHeader("session_id")
	if sessionID == "" {
		return ""
placeholder
	hash := sha256.Sum256([]byte(sessionID))
	return hex.EncodeToString(hash[:])
placeholder

// SelectAccount selects an OpenAI account with sticky session support
func (s *OpenAIGatewayService) SelectAccount(ctx context.Context, groupID *int64, sessionHash string) (*model.Account, error) {
	return s.SelectAccountForModel(ctx, groupID, sessionHash, "")
placeholder

// SelectAccountForModel selects an account supporting the requested model
func (s *OpenAIGatewayService) SelectAccountForModel(ctx context.Context, groupID *int64, sessionHash string, requestedModel string) (*model.Account, error) {
	// 1. Check sticky session
	if sessionHash != "" {
		accountID, err := s.cache.GetSessionAccountID(ctx, "openai:"+sessionHash)
		if err == nil && accountID > 0 {
			account, err := s.accountRepo.GetByID(ctx, accountID)
			if err == nil && account.IsSchedulable() && account.IsOpenAI() && (requestedModel == "" || account.IsModelSupported(requestedModel)) {
				// Refresh sticky session TTL
				_ = s.cache.RefreshSessionTTL(ctx, "openai:"+sessionHash, openaiStickySessionTTL)
				return account, nil
		placeholder
	placeholder
placeholder

	// 2. Get schedulable OpenAI accounts
	var accounts []model.Account
	var err error
	if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, model.PlatformOpenAI)
placeholder else {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, model.PlatformOpenAI)
placeholder
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder

	// 3. Select by priority + LRU
	var selected *model.Account
	for i := range accounts {
		acc := &accounts[i]
		// Check model support
		if requestedModel != "" && !acc.IsModelSupported(requestedModel) {
			continue
	placeholder
		if selected == nil {
			selected = acc
			continue
	placeholder
		// Lower priority value means higher priority
		if acc.Priority < selected.Priority {
			selected = acc
	placeholder else if acc.Priority == selected.Priority {
			// Same priority, select least recently used
			if acc.LastUsedAt == nil || (selected.LastUsedAt != nil && acc.LastUsedAt.Before(*selected.LastUsedAt)) {
				selected = acc
		placeholder
	placeholder
placeholder

	if selected == nil {
		if requestedModel != "" {
			return nil, fmt.Errorf("no available OpenAI accounts supporting model: %s", requestedModel)
	placeholder
		return nil, errors.New("no available OpenAI accounts")
placeholder

	// 4. Set sticky session
	if sessionHash != "" {
		_ = s.cache.SetSessionAccountID(ctx, "openai:"+sessionHash, selected.ID, openaiStickySessionTTL)
placeholder

	return selected, nil
placeholder

// GetAccessToken gets the access token for an OpenAI account
func (s *OpenAIGatewayService) GetAccessToken(ctx context.Context, account *model.Account) (string, string, error) {
	if account.Type == model.AccountTypeOAuth {
		accessToken := account.GetOpenAIAccessToken()
		if accessToken == "" {
			return "", "", errors.New("access_token not found in credentials")
	placeholder
		return accessToken, "oauth", nil
placeholder else if account.Type == model.AccountTypeApiKey {
		apiKey := account.GetOpenAIApiKey()
		if apiKey == "" {
			return "", "", errors.New("api_key not found in credentials")
	placeholder
		return apiKey, "apikey", nil
placeholder
	return "", "", fmt.Errorf("unsupported account type: %s", account.Type)
placeholder

// Forward forwards request to OpenAI API
func (s *OpenAIGatewayService) Forward(ctx context.Context, c *gin.Context, account *model.Account, body []byte) (*OpenAIForwardResult, error) {
	startTime := time.Now()

	// Parse request body once (avoid multiple parse/serialize cycles)
	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		return nil, fmt.Errorf("parse request: %w", err)
placeholder

	// Extract model and stream from parsed body
	reqModel, _ := reqBody["model"].(string)
	reqStream, _ := reqBody["stream"].(bool)

	// Track if body needs re-serialization
	bodyModified := false
	originalModel := reqModel

	// Apply model mapping
	mappedModel := account.GetMappedModel(reqModel)
	if mappedModel != reqModel {
		reqBody["model"] = mappedModel
		bodyModified = true
placeholder

	// For OAuth accounts using ChatGPT internal API, add store: false
	if account.Type == model.AccountTypeOAuth {
		reqBody["store"] = false
		bodyModified = true
placeholder

	// Re-serialize body only if modified
	if bodyModified {
		var err error
		body, err = json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("serialize request body: %w", err)
	placeholder
placeholder

	// Get access token
	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder

	// Build upstream request
	upstreamReq, err := s.buildUpstreamRequest(ctx, c, account, body, token, reqStream)
	if err != nil {
		return nil, err
placeholder

	// Get proxy URL
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	// Send request
	resp, err := s.httpUpstream.Do(upstreamReq, proxyURL)
	if err != nil {
		return nil, fmt.Errorf("upstream request failed: %w", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	// Handle error response
	if resp.StatusCode >= 400 {
		return s.handleErrorResponse(ctx, resp, c, account)
placeholder

	// Handle normal response
	var usage *OpenAIUsage
	var firstTokenMs *int
	if reqStream {
		streamResult, err := s.handleStreamingResponse(ctx, resp, c, account, startTime, originalModel, mappedModel)
		if err != nil {
			return nil, err
	placeholder
		usage = streamResult.usage
		firstTokenMs = streamResult.firstTokenMs
placeholder else {
		usage, err = s.handleNonStreamingResponse(ctx, resp, c, account, originalModel, mappedModel)
		if err != nil {
			return nil, err
	placeholder
placeholder

	return &OpenAIForwardResult{
		RequestID:    resp.Header.Get("x-request-id"),
		Usage:        *usage,
		Model:        originalModel,
		Stream:       reqStream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
placeholder, nil
placeholder

func (s *OpenAIGatewayService) buildUpstreamRequest(ctx context.Context, c *gin.Context, account *model.Account, body []byte, token string, isStream bool) (*http.Request, error) {
	// Determine target URL based on account type
	var targetURL string
	if account.Type == model.AccountTypeOAuth {
		// OAuth accounts use ChatGPT internal API
		targetURL = chatgptCodexURL
placeholder else if account.Type == model.AccountTypeApiKey {
		// API Key accounts use Platform API or custom base URL
		baseURL := account.GetOpenAIBaseURL()
		if baseURL != "" {
			targetURL = baseURL + "/v1/responses"
	placeholder else {
			targetURL = openaiPlatformAPIURL
	placeholder
placeholder else {
		targetURL = openaiPlatformAPIURL
placeholder

	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder

	// Set authentication header
	req.Header.Set("authorization", "Bearer "+token)

	// Set headers specific to OAuth accounts (ChatGPT internal API)
	if account.Type == model.AccountTypeOAuth {
		// Required: set Host for ChatGPT API (must use req.Host, not Header.Set)
		req.Host = "chatgpt.com"
		// Required: set chatgpt-account-id header
		chatgptAccountID := account.GetChatGPTAccountID()
		if chatgptAccountID != "" {
			req.Header.Set("chatgpt-account-id", chatgptAccountID)
	placeholder
		// Set accept header based on stream mode
		if isStream {
			req.Header.Set("accept", "text/event-stream")
	placeholder else {
			req.Header.Set("accept", "application/json")
	placeholder
placeholder

	// Whitelist passthrough headers
	for key, values := range c.Request.Header {
		lowerKey := strings.ToLower(key)
		if openaiAllowedHeaders[lowerKey] {
			for _, v := range values {
				req.Header.Add(key, v)
		placeholder
	placeholder
placeholder

	// Apply custom User-Agent if configured
	customUA := account.GetOpenAIUserAgent()
	if customUA != "" {
		req.Header.Set("user-agent", customUA)
placeholder

	// Ensure required headers exist
	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/json")
placeholder

	return req, nil
placeholder

func (s *OpenAIGatewayService) handleErrorResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *model.Account) (*OpenAIForwardResult, error) {
	body, _ := io.ReadAll(resp.Body)

	// Check custom error codes
	if !account.ShouldHandleErrorCode(resp.StatusCode) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"type":    "upstream_error",
				"message": "Upstream gateway error",
		placeholder,
	placeholder)
		return nil, fmt.Errorf("upstream error: %d (not in custom error codes)", resp.StatusCode)
placeholder

	// Handle upstream error (mark account status)
	s.rateLimitService.HandleUpstreamError(ctx, account, resp.StatusCode, resp.Header, body)

	// Return appropriate error response
	var errType, errMsg string
	var statusCode int

	switch resp.StatusCode {
	case 401:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream authentication failed, please contact administrator"
	case 403:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream access forbidden, please contact administrator"
	case 429:
		statusCode = http.StatusTooManyRequests
		errType = "rate_limit_error"
		errMsg = "Upstream rate limit exceeded, please retry later"
	default:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream request failed"
placeholder

	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": errMsg,
	placeholder,
placeholder)

	return nil, fmt.Errorf("upstream error: %d", resp.StatusCode)
placeholder

// openaiStreamingResult streaming response result
type openaiStreamingResult struct {
	usage        *OpenAIUsage
	firstTokenMs *int
placeholder

func (s *OpenAIGatewayService) handleStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *model.Account, startTime time.Time, originalModel, mappedModel string) (*openaiStreamingResult, error) {
	// Set SSE response headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// Pass through other headers
	if v := resp.Header.Get("x-request-id"); v != "" {
		c.Header("x-request-id", v)
placeholder

	w := c.Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder

	usage := &OpenAIUsage{placeholder
	var firstTokenMs *int
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	needModelReplace := originalModel != mappedModel

	for scanner.Scan() {
		line := scanner.Text()

		// Replace model in response if needed
		if needModelReplace && strings.HasPrefix(line, "data: ") {
			line = s.replaceModelInSSELine(line, mappedModel, originalModel)
	placeholder

		// Forward line
		if _, err := fmt.Fprintf(w, "%s\n", line); err != nil {
			return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, err
	placeholder
		flusher.Flush()

		// Parse usage data
		if strings.HasPrefix(line, "data: ") {
			data := line[6:]
			// Record first token time
			if firstTokenMs == nil && data != "" && data != "[DONE]" {
				ms := int(time.Since(startTime).Milliseconds())
				firstTokenMs = &ms
		placeholder
			s.parseSSEUsage(data, usage)
	placeholder
placeholder

	if err := scanner.Err(); err != nil {
		return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, fmt.Errorf("stream read error: %w", err)
placeholder

	return &openaiStreamingResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, nil
placeholder

func (s *OpenAIGatewayService) replaceModelInSSELine(line, fromModel, toModel string) string {
	data := line[6:]
	if data == "" || data == "[DONE]" {
		return line
placeholder

	var event map[string]any
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		return line
placeholder

	// Replace model in response
	if m, ok := event["model"].(string); ok && m == fromModel {
		event["model"] = toModel
		newData, err := json.Marshal(event)
		if err != nil {
			return line
	placeholder
		return "data: " + string(newData)
placeholder

	// Check nested response
	if response, ok := event["response"].(map[string]any); ok {
		if m, ok := response["model"].(string); ok && m == fromModel {
			response["model"] = toModel
			newData, err := json.Marshal(event)
			if err != nil {
				return line
		placeholder
			return "data: " + string(newData)
	placeholder
placeholder

	return line
placeholder

func (s *OpenAIGatewayService) parseSSEUsage(data string, usage *OpenAIUsage) {
	// Parse response.completed event for usage (OpenAI Responses format)
	var event struct {
		Type     string `json:"type"`
		Response struct {
			Usage struct {
				InputTokens       int `json:"input_tokens"`
				OutputTokens      int `json:"output_tokens"`
				InputTokenDetails struct {
					CachedTokens int `json:"cached_tokens"`
			placeholder `json:"input_tokens_details"`
		placeholder `json:"usage"`
	placeholder `json:"response"`
placeholder

	if json.Unmarshal([]byte(data), &event) == nil && event.Type == "response.completed" {
		usage.InputTokens = event.Response.Usage.InputTokens
		usage.OutputTokens = event.Response.Usage.OutputTokens
		usage.CacheReadInputTokens = event.Response.Usage.InputTokenDetails.CachedTokens
placeholder
placeholder

func (s *OpenAIGatewayService) handleNonStreamingResponse(ctx context.Context, resp *http.Response, c *gin.Context, account *model.Account, originalModel, mappedModel string) (*OpenAIUsage, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
placeholder

	// Parse usage
	var response struct {
		Usage struct {
			InputTokens       int `json:"input_tokens"`
			OutputTokens      int `json:"output_tokens"`
			InputTokenDetails struct {
				CachedTokens int `json:"cached_tokens"`
		placeholder `json:"input_tokens_details"`
	placeholder `json:"usage"`
placeholder
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
placeholder

	usage := &OpenAIUsage{
		InputTokens:          response.Usage.InputTokens,
		OutputTokens:         response.Usage.OutputTokens,
		CacheReadInputTokens: response.Usage.InputTokenDetails.CachedTokens,
placeholder

	// Replace model in response if needed
	if originalModel != mappedModel {
		body = s.replaceModelInResponseBody(body, mappedModel, originalModel)
placeholder

	// Pass through headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
	placeholder
placeholder

	c.Data(resp.StatusCode, "application/json", body)

	return usage, nil
placeholder

func (s *OpenAIGatewayService) replaceModelInResponseBody(body []byte, fromModel, toModel string) []byte {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return body
placeholder

	model, ok := resp["model"].(string)
	if !ok || model != fromModel {
		return body
placeholder

	resp["model"] = toModel
	newBody, err := json.Marshal(resp)
	if err != nil {
		return body
placeholder

	return newBody
placeholder

// OpenAIRecordUsageInput input for recording usage
type OpenAIRecordUsageInput struct {
	Result       *OpenAIForwardResult
	ApiKey       *model.ApiKey
	User         *model.User
	Account      *model.Account
	Subscription *model.UserSubscription
placeholder

// RecordUsage records usage and deducts balance
func (s *OpenAIGatewayService) RecordUsage(ctx context.Context, input *OpenAIRecordUsageInput) error {
	result := input.Result
	apiKey := input.ApiKey
	user := input.User
	account := input.Account
	subscription := input.Subscription

	// Calculate cost
	tokens := UsageTokens{
		InputTokens:         result.Usage.InputTokens,
		OutputTokens:        result.Usage.OutputTokens,
		CacheCreationTokens: result.Usage.CacheCreationInputTokens,
		CacheReadTokens:     result.Usage.CacheReadInputTokens,
placeholder

	// Get rate multiplier
	multiplier := s.cfg.Default.RateMultiplier
	if apiKey.GroupID != nil && apiKey.Group != nil {
		multiplier = apiKey.Group.RateMultiplier
placeholder

	cost, err := s.billingService.CalculateCost(result.Model, tokens, multiplier)
	if err != nil {
		cost = &CostBreakdown{ActualCost: 0placeholder
placeholder

	// Determine billing type
	isSubscriptionBilling := subscription != nil && apiKey.Group != nil && apiKey.Group.IsSubscriptionType()
	billingType := model.BillingTypeBalance
	if isSubscriptionBilling {
		billingType = model.BillingTypeSubscription
placeholder

	// Create usage log
	durationMs := int(result.Duration.Milliseconds())
	usageLog := &model.UsageLog{
		UserID:              user.ID,
		ApiKeyID:            apiKey.ID,
		AccountID:           account.ID,
		RequestID:           result.RequestID,
		Model:               result.Model,
		InputTokens:         result.Usage.InputTokens,
		OutputTokens:        result.Usage.OutputTokens,
		CacheCreationTokens: result.Usage.CacheCreationInputTokens,
		CacheReadTokens:     result.Usage.CacheReadInputTokens,
		InputCost:           cost.InputCost,
		OutputCost:          cost.OutputCost,
		CacheCreationCost:   cost.CacheCreationCost,
		CacheReadCost:       cost.CacheReadCost,
		TotalCost:           cost.TotalCost,
		ActualCost:          cost.ActualCost,
		RateMultiplier:      multiplier,
		BillingType:         billingType,
		Stream:              result.Stream,
		DurationMs:          &durationMs,
		FirstTokenMs:        result.FirstTokenMs,
		CreatedAt:           time.Now(),
placeholder

	if apiKey.GroupID != nil {
		usageLog.GroupID = apiKey.GroupID
placeholder
	if subscription != nil {
		usageLog.SubscriptionID = &subscription.ID
placeholder

	_ = s.usageLogRepo.Create(ctx, usageLog)

	// Deduct based on billing type
	if isSubscriptionBilling {
		if cost.TotalCost > 0 {
			_ = s.userSubRepo.IncrementUsage(ctx, subscription.ID, cost.TotalCost)
			go func() {
				cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = s.billingCacheService.UpdateSubscriptionUsage(cacheCtx, user.ID, *apiKey.GroupID, cost.TotalCost)
		placeholder()
	placeholder
placeholder else {
		if cost.ActualCost > 0 {
			_ = s.userRepo.DeductBalance(ctx, user.ID, cost.ActualCost)
			go func() {
				cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = s.billingCacheService.DeductBalanceCache(cacheCtx, user.ID, cost.ActualCost)
		placeholder()
	placeholder
placeholder

	// Update account last used
	_ = s.accountRepo.UpdateLastUsed(ctx, account.ID)

	return nil
placeholder
