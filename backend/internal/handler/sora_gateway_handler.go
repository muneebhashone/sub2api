package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SoraGatewayHandler handles Sora chat completions requests
type SoraGatewayHandler struct {
	gatewayService      *service.GatewayService
	soraGatewayService  *service.SoraGatewayService
	billingCacheService *service.BillingCacheService
	concurrencyHelper   *ConcurrencyHelper
	maxAccountSwitches  int
	streamMode          string
	sora2apiBaseURL     string
	soraMediaSigningKey string
	mediaClient         *http.Client
placeholder

// NewSoraGatewayHandler creates a new SoraGatewayHandler
func NewSoraGatewayHandler(
	gatewayService *service.GatewayService,
	soraGatewayService *service.SoraGatewayService,
	concurrencyService *service.ConcurrencyService,
	billingCacheService *service.BillingCacheService,
	cfg *config.Config,
) *SoraGatewayHandler {
	pingInterval := time.Duration(0)
	maxAccountSwitches := 3
	streamMode := "force"
	signKey := ""
	if cfg != nil {
		pingInterval = time.Duration(cfg.Concurrency.PingInterval) * time.Second
		if cfg.Gateway.MaxAccountSwitches > 0 {
			maxAccountSwitches = cfg.Gateway.MaxAccountSwitches
	placeholder
		if mode := strings.TrimSpace(cfg.Gateway.SoraStreamMode); mode != "" {
			streamMode = mode
	placeholder
		signKey = strings.TrimSpace(cfg.Gateway.SoraMediaSigningKey)
placeholder
	baseURL := ""
	if cfg != nil {
		baseURL = strings.TrimRight(strings.TrimSpace(cfg.Sora2API.BaseURL), "/")
placeholder
	mediaTimeout := 180 * time.Second
	if cfg != nil && cfg.Gateway.SoraRequestTimeoutSeconds > 0 {
		mediaTimeout = time.Duration(cfg.Gateway.SoraRequestTimeoutSeconds) * time.Second
placeholder
	return &SoraGatewayHandler{
		gatewayService:      gatewayService,
		soraGatewayService:  soraGatewayService,
		billingCacheService: billingCacheService,
		concurrencyHelper:   NewConcurrencyHelper(concurrencyService, SSEPingFormatComment, pingInterval),
		maxAccountSwitches:  maxAccountSwitches,
		streamMode:          strings.ToLower(streamMode),
		sora2apiBaseURL:     baseURL,
		soraMediaSigningKey: signKey,
		mediaClient:         &http.Client{Timeout: mediaTimeoutplaceholder,
placeholder
placeholder

// ChatCompletions handles Sora /v1/chat/completions endpoint
func (h *SoraGatewayHandler) ChatCompletions(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
placeholder

	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		h.errorResponse(c, http.StatusInternalServerError, "api_error", "User context not found")
		return
placeholder

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			h.errorResponse(c, http.StatusRequestEntityTooLarge, "invalid_request_error", buildBodyTooLargeMessage(maxErr.Limit))
			return
	placeholder
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to read request body")
		return
placeholder
	if len(body) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Request body is empty")
		return
placeholder

	setOpsRequestContext(c, "", false, body)

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body")
		return
placeholder

	reqModel, _ := reqBody["model"].(string)
	if reqModel == "" {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "model is required")
		return
placeholder
	reqMessages, _ := reqBody["messages"].([]any)
	if len(reqMessages) == 0 {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "messages is required")
		return
placeholder

	clientStream, _ := reqBody["stream"].(bool)
	if !clientStream {
		if h.streamMode == "error" {
			h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "Sora requires stream=true")
			return
	placeholder
		reqBody["stream"] = true
		updated, err := json.Marshal(reqBody)
		if err != nil {
			h.errorResponse(c, http.StatusInternalServerError, "api_error", "Failed to process request")
			return
	placeholder
		body = updated
placeholder

	setOpsRequestContext(c, reqModel, clientStream, body)

	platform := ""
	if forced, ok := middleware2.GetForcePlatformFromContext(c); ok {
		platform = forced
placeholder else if apiKey.Group != nil {
		platform = apiKey.Group.Platform
placeholder
	if platform != service.PlatformSora {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "This endpoint only supports Sora platform")
		return
placeholder

	streamStarted := false
	subscription, _ := middleware2.GetSubscriptionFromContext(c)

	maxWait := service.CalculateMaxWait(subject.Concurrency)
	canWait, err := h.concurrencyHelper.IncrementWaitCount(c.Request.Context(), subject.UserID, maxWait)
	waitCounted := false
	if err != nil {
		log.Printf("Increment wait count failed: %v", err)
placeholder else if !canWait {
		h.errorResponse(c, http.StatusTooManyRequests, "rate_limit_error", "Too many pending requests, please retry later")
		return
placeholder
	if err == nil && canWait {
		waitCounted = true
placeholder
	defer func() {
		if waitCounted {
			h.concurrencyHelper.DecrementWaitCount(c.Request.Context(), subject.UserID)
	placeholder
placeholder()

	userReleaseFunc, err := h.concurrencyHelper.AcquireUserSlotWithWait(c, subject.UserID, subject.Concurrency, clientStream, &streamStarted)
	if err != nil {
		log.Printf("User concurrency acquire failed: %v", err)
		h.handleConcurrencyError(c, err, "user", streamStarted)
		return
placeholder
	if waitCounted {
		h.concurrencyHelper.DecrementWaitCount(c.Request.Context(), subject.UserID)
		waitCounted = false
placeholder
	userReleaseFunc = wrapReleaseOnDone(c.Request.Context(), userReleaseFunc)
	if userReleaseFunc != nil {
		defer userReleaseFunc()
placeholder

	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription); err != nil {
		log.Printf("Billing eligibility check failed after wait: %v", err)
		status, code, message := billingErrorDetails(err)
		h.handleStreamingAwareError(c, status, code, message, streamStarted)
		return
placeholder

	sessionHash := generateOpenAISessionHash(c, reqBody)

	maxAccountSwitches := h.maxAccountSwitches
	switchCount := 0
	failedAccountIDs := make(map[int64]struct{placeholder)
	lastFailoverStatus := 0

	for {
		selection, err := h.gatewayService.SelectAccountWithLoadAwareness(c.Request.Context(), apiKey.GroupID, sessionHash, reqModel, failedAccountIDs, "")
		if err != nil {
			log.Printf("[Sora Handler] SelectAccount failed: %v", err)
			if len(failedAccountIDs) == 0 {
				h.handleStreamingAwareError(c, http.StatusServiceUnavailable, "api_error", "No available accounts: "+err.Error(), streamStarted)
				return
		placeholder
			h.handleFailoverExhausted(c, lastFailoverStatus, streamStarted)
			return
	placeholder
		account := selection.Account
		setOpsSelectedAccount(c, account.ID)

		accountReleaseFunc := selection.ReleaseFunc
		if !selection.Acquired {
			if selection.WaitPlan == nil {
				h.handleStreamingAwareError(c, http.StatusServiceUnavailable, "api_error", "No available accounts", streamStarted)
				return
		placeholder
			accountWaitCounted := false
			canWait, err := h.concurrencyHelper.IncrementAccountWaitCount(c.Request.Context(), account.ID, selection.WaitPlan.MaxWaiting)
			if err != nil {
				log.Printf("Increment account wait count failed: %v", err)
		placeholder else if !canWait {
				log.Printf("Account wait queue full: account=%d", account.ID)
				h.handleStreamingAwareError(c, http.StatusTooManyRequests, "rate_limit_error", "Too many pending requests, please retry later", streamStarted)
				return
		placeholder
			if err == nil && canWait {
				accountWaitCounted = true
		placeholder
			defer func() {
				if accountWaitCounted {
					h.concurrencyHelper.DecrementAccountWaitCount(c.Request.Context(), account.ID)
			placeholder
		placeholder()

			accountReleaseFunc, err = h.concurrencyHelper.AcquireAccountSlotWithWaitTimeout(
				c,
				account.ID,
				selection.WaitPlan.MaxConcurrency,
				selection.WaitPlan.Timeout,
				clientStream,
				&streamStarted,
			)
			if err != nil {
				log.Printf("Account concurrency acquire failed: %v", err)
				h.handleConcurrencyError(c, err, "account", streamStarted)
				return
		placeholder
			if accountWaitCounted {
				h.concurrencyHelper.DecrementAccountWaitCount(c.Request.Context(), account.ID)
				accountWaitCounted = false
		placeholder
	placeholder
		accountReleaseFunc = wrapReleaseOnDone(c.Request.Context(), accountReleaseFunc)

		result, err := h.soraGatewayService.Forward(c.Request.Context(), c, account, body, clientStream)
		if accountReleaseFunc != nil {
			accountReleaseFunc()
	placeholder
		if err != nil {
			var failoverErr *service.UpstreamFailoverError
			if errors.As(err, &failoverErr) {
				failedAccountIDs[account.ID] = struct{placeholder{placeholder
				if switchCount >= maxAccountSwitches {
					lastFailoverStatus = failoverErr.StatusCode
					h.handleFailoverExhausted(c, lastFailoverStatus, streamStarted)
					return
			placeholder
				lastFailoverStatus = failoverErr.StatusCode
				switchCount++
				log.Printf("Account %d: upstream error %d, switching account %d/%d", account.ID, failoverErr.StatusCode, switchCount, maxAccountSwitches)
				continue
		placeholder
			log.Printf("Account %d: Forward request failed: %v", account.ID, err)
			return
	placeholder

		userAgent := c.GetHeader("User-Agent")
		clientIP := ip.GetClientIP(c)

		go func(result *service.ForwardResult, usedAccount *service.Account, ua, ip string) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := h.gatewayService.RecordUsage(ctx, &service.RecordUsageInput{
				Result:       result,
				APIKey:       apiKey,
				User:         apiKey.User,
				Account:      usedAccount,
				Subscription: subscription,
				UserAgent:    ua,
				IPAddress:    ip,
		placeholder); err != nil {
				log.Printf("Record usage failed: %v", err)
		placeholder
	placeholder(result, account, userAgent, clientIP)
		return
placeholder
placeholder

func generateOpenAISessionHash(c *gin.Context, reqBody map[string]any) string {
	if c == nil {
		return ""
placeholder
	sessionID := strings.TrimSpace(c.GetHeader("session_id"))
	if sessionID == "" {
		sessionID = strings.TrimSpace(c.GetHeader("conversation_id"))
placeholder
	if sessionID == "" && reqBody != nil {
		if v, ok := reqBody["prompt_cache_key"].(string); ok {
			sessionID = strings.TrimSpace(v)
	placeholder
placeholder
	if sessionID == "" {
		return ""
placeholder
	hash := sha256.Sum256([]byte(sessionID))
	return hex.EncodeToString(hash[:])
placeholder

func (h *SoraGatewayHandler) handleConcurrencyError(c *gin.Context, err error, slotType string, streamStarted bool) {
	h.handleStreamingAwareError(c, http.StatusTooManyRequests, "rate_limit_error",
		fmt.Sprintf("Concurrency limit exceeded for %s, please retry later", slotType), streamStarted)
placeholder

func (h *SoraGatewayHandler) handleFailoverExhausted(c *gin.Context, statusCode int, streamStarted bool) {
	status, errType, errMsg := h.mapUpstreamError(statusCode)
	h.handleStreamingAwareError(c, status, errType, errMsg, streamStarted)
placeholder

func (h *SoraGatewayHandler) mapUpstreamError(statusCode int) (int, string, string) {
	switch statusCode {
	case 401:
		return http.StatusBadGateway, "upstream_error", "Upstream authentication failed, please contact administrator"
	case 403:
		return http.StatusBadGateway, "upstream_error", "Upstream access forbidden, please contact administrator"
	case 429:
		return http.StatusTooManyRequests, "rate_limit_error", "Upstream rate limit exceeded, please retry later"
	case 529:
		return http.StatusServiceUnavailable, "upstream_error", "Upstream service overloaded, please retry later"
	case 500, 502, 503, 504:
		return http.StatusBadGateway, "upstream_error", "Upstream service temporarily unavailable"
	default:
		return http.StatusBadGateway, "upstream_error", "Upstream request failed"
placeholder
placeholder

func (h *SoraGatewayHandler) handleStreamingAwareError(c *gin.Context, status int, errType, message string, streamStarted bool) {
	if streamStarted {
		flusher, ok := c.Writer.(http.Flusher)
		if ok {
			errorEvent := fmt.Sprintf(`event: error`+"\n"+`data: {"error": {"type": "%s", "message": "%s"placeholderplaceholder`+"\n\n", errType, message)
			if _, err := fmt.Fprint(c.Writer, errorEvent); err != nil {
				_ = c.Error(err)
		placeholder
			flusher.Flush()
	placeholder
		return
placeholder
	h.errorResponse(c, status, errType, message)
placeholder

func (h *SoraGatewayHandler) errorResponse(c *gin.Context, status int, errType, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": message,
	placeholder,
placeholder)
placeholder

// MediaProxy proxies /tmp or /static media files from sora2api
func (h *SoraGatewayHandler) MediaProxy(c *gin.Context) {
	h.proxySoraMedia(c, false)
placeholder

// MediaProxySigned proxies /tmp or /static media files with signature verification
func (h *SoraGatewayHandler) MediaProxySigned(c *gin.Context) {
	h.proxySoraMedia(c, true)
placeholder

func (h *SoraGatewayHandler) proxySoraMedia(c *gin.Context, requireSignature bool) {
	if h.sora2apiBaseURL == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"type":    "api_error",
				"message": "sora2api 未配置",
		placeholder,
	placeholder)
		return
placeholder

	rawPath := c.Param("filepath")
	if rawPath == "" {
		c.Status(http.StatusNotFound)
		return
placeholder
	cleaned := path.Clean(rawPath)
	if !strings.HasPrefix(cleaned, "/tmp/") && !strings.HasPrefix(cleaned, "/static/") {
		c.Status(http.StatusNotFound)
		return
placeholder

	query := c.Request.URL.Query()
	if requireSignature {
		if h.soraMediaSigningKey == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "api_error",
					"message": "Sora 媒体签名未配置",
			placeholder,
		placeholder)
			return
	placeholder
		expiresStr := strings.TrimSpace(query.Get("expires"))
		signature := strings.TrimSpace(query.Get("sig"))
		expires, err := strconv.ParseInt(expiresStr, 10, 64)
		if err != nil || expires <= time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"type":    "authentication_error",
					"message": "Sora 媒体签名已过期",
			placeholder,
		placeholder)
			return
	placeholder
		query.Del("sig")
		query.Del("expires")
		signingQuery := query.Encode()
		if !service.VerifySoraMediaURL(cleaned, signingQuery, expires, signature, h.soraMediaSigningKey) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"type":    "authentication_error",
					"message": "Sora 媒体签名无效",
			placeholder,
		placeholder)
			return
	placeholder
placeholder

	targetURL := h.sora2apiBaseURL + cleaned
	if rawQuery := query.Encode(); rawQuery != "" {
		targetURL += "?" + rawQuery
placeholder

	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, nil)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
placeholder
	copyHeaders := []string{"Range", "If-Range", "If-Modified-Since", "If-None-Match", "Accept", "User-Agent"placeholder
	for _, key := range copyHeaders {
		if val := c.GetHeader(key); val != "" {
			req.Header.Set(key, val)
	placeholder
placeholder

	client := h.mediaClient
	if client == nil {
		client = http.DefaultClient
placeholder
	resp, err := client.Do(req)
	if err != nil {
		c.Status(http.StatusBadGateway)
		return
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	for _, key := range []string{"Content-Type", "Content-Length", "Accept-Ranges", "Content-Range", "Cache-Control", "Last-Modified", "ETag"placeholder {
		if val := resp.Header.Get(key); val != "" {
			c.Header(key, val)
	placeholder
placeholder
	c.Status(resp.StatusCode)
	_, _ = io.Copy(c.Writer, resp.Body)
placeholder
