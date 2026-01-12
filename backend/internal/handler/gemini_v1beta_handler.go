package handler

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/Wei-Shaw/sub2api/internal/pkg/gemini"
	"github.com/Wei-Shaw/sub2api/internal/pkg/googleapi"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// GeminiV1BetaListModels proxies:
// GET /v1beta/models
func (h *GatewayHandler) GeminiV1BetaListModels(c *gin.Context) {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		googleError(c, http.StatusUnauthorized, "Invalid API key")
		return
placeholder
	// 检查平台：优先使用强制平台（/antigravity 路由），否则要求 gemini 分组
	forcePlatform, hasForcePlatform := middleware.GetForcePlatformFromContext(c)
	if !hasForcePlatform && (apiKey.Group == nil || apiKey.Group.Platform != service.PlatformGemini) {
		googleError(c, http.StatusBadRequest, "API key group platform is not gemini")
		return
placeholder

	// 强制 antigravity 模式：返回 antigravity 支持的模型列表
	if forcePlatform == service.PlatformAntigravity {
		c.JSON(http.StatusOK, antigravity.FallbackGeminiModelsList())
		return
placeholder

	account, err := h.geminiCompatService.SelectAccountForAIStudioEndpoints(c.Request.Context(), apiKey.GroupID)
	if err != nil {
		// 没有 gemini 账户，检查是否有 antigravity 账户可用
		hasAntigravity, _ := h.geminiCompatService.HasAntigravityAccounts(c.Request.Context(), apiKey.GroupID)
		if hasAntigravity {
			// antigravity 账户使用静态模型列表
			c.JSON(http.StatusOK, gemini.FallbackModelsList())
			return
	placeholder
		googleError(c, http.StatusServiceUnavailable, "No available Gemini accounts: "+err.Error())
		return
placeholder

	res, err := h.geminiCompatService.ForwardAIStudioGET(c.Request.Context(), account, "/v1beta/models")
	if err != nil {
		googleError(c, http.StatusBadGateway, err.Error())
		return
placeholder
	if shouldFallbackGeminiModels(res) {
		c.JSON(http.StatusOK, gemini.FallbackModelsList())
		return
placeholder
	writeUpstreamResponse(c, res)
placeholder

// GeminiV1BetaGetModel proxies:
// GET /v1beta/models/{modelplaceholder
func (h *GatewayHandler) GeminiV1BetaGetModel(c *gin.Context) {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		googleError(c, http.StatusUnauthorized, "Invalid API key")
		return
placeholder
	// 检查平台：优先使用强制平台（/antigravity 路由），否则要求 gemini 分组
	forcePlatform, hasForcePlatform := middleware.GetForcePlatformFromContext(c)
	if !hasForcePlatform && (apiKey.Group == nil || apiKey.Group.Platform != service.PlatformGemini) {
		googleError(c, http.StatusBadRequest, "API key group platform is not gemini")
		return
placeholder

	modelName := strings.TrimSpace(c.Param("model"))
	if modelName == "" {
		googleError(c, http.StatusBadRequest, "Missing model in URL")
		return
placeholder

	// 强制 antigravity 模式：返回 antigravity 模型信息
	if forcePlatform == service.PlatformAntigravity {
		c.JSON(http.StatusOK, antigravity.FallbackGeminiModel(modelName))
		return
placeholder

	account, err := h.geminiCompatService.SelectAccountForAIStudioEndpoints(c.Request.Context(), apiKey.GroupID)
	if err != nil {
		// 没有 gemini 账户，检查是否有 antigravity 账户可用
		hasAntigravity, _ := h.geminiCompatService.HasAntigravityAccounts(c.Request.Context(), apiKey.GroupID)
		if hasAntigravity {
			// antigravity 账户使用静态模型信息
			c.JSON(http.StatusOK, gemini.FallbackModel(modelName))
			return
	placeholder
		googleError(c, http.StatusServiceUnavailable, "No available Gemini accounts: "+err.Error())
		return
placeholder

	res, err := h.geminiCompatService.ForwardAIStudioGET(c.Request.Context(), account, "/v1beta/models/"+modelName)
	if err != nil {
		googleError(c, http.StatusBadGateway, err.Error())
		return
placeholder
	if shouldFallbackGeminiModels(res) {
		c.JSON(http.StatusOK, gemini.FallbackModel(modelName))
		return
placeholder
	writeUpstreamResponse(c, res)
placeholder

// GeminiV1BetaModels proxies Gemini native REST endpoints like:
// POST /v1beta/models/{modelplaceholder:generateContent
// POST /v1beta/models/{modelplaceholder:streamGenerateContent?alt=sse
func (h *GatewayHandler) GeminiV1BetaModels(c *gin.Context) {
	apiKey, ok := middleware.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		googleError(c, http.StatusUnauthorized, "Invalid API key")
		return
placeholder
	authSubject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		googleError(c, http.StatusInternalServerError, "User context not found")
		return
placeholder

	// 检查平台：优先使用强制平台（/antigravity 路由，中间件已设置 request.Context），否则要求 gemini 分组
	if !middleware.HasForcePlatform(c) {
		if apiKey.Group == nil || apiKey.Group.Platform != service.PlatformGemini {
			googleError(c, http.StatusBadRequest, "API key group platform is not gemini")
			return
	placeholder
placeholder

	modelName, action, err := parseGeminiModelAction(strings.TrimPrefix(c.Param("modelAction"), "/"))
	if err != nil {
		googleError(c, http.StatusNotFound, err.Error())
		return
placeholder

	stream := action == "streamGenerateContent"

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if maxErr, ok := extractMaxBytesError(err); ok {
			googleError(c, http.StatusRequestEntityTooLarge, buildBodyTooLargeMessage(maxErr.Limit))
			return
	placeholder
		googleError(c, http.StatusBadRequest, "Failed to read request body")
		return
placeholder
	if len(body) == 0 {
		googleError(c, http.StatusBadRequest, "Request body is empty")
		return
placeholder

	setOpsRequestContext(c, modelName, stream, body)

	// Get subscription (may be nil)
	subscription, _ := middleware.GetSubscriptionFromContext(c)

	// For Gemini native API, do not send Claude-style ping frames.
	geminiConcurrency := NewConcurrencyHelper(h.concurrencyHelper.concurrencyService, SSEPingFormatNone, 0)

	// 0) wait queue check
	maxWait := service.CalculateMaxWait(authSubject.Concurrency)
	canWait, err := geminiConcurrency.IncrementWaitCount(c.Request.Context(), authSubject.UserID, maxWait)
	waitCounted := false
	if err != nil {
		log.Printf("Increment wait count failed: %v", err)
placeholder else if !canWait {
		googleError(c, http.StatusTooManyRequests, "Too many pending requests, please retry later")
		return
placeholder
	if err == nil && canWait {
		waitCounted = true
placeholder
	defer func() {
		if waitCounted {
			geminiConcurrency.DecrementWaitCount(c.Request.Context(), authSubject.UserID)
	placeholder
placeholder()

	// 1) user concurrency slot
	streamStarted := false
	userReleaseFunc, err := geminiConcurrency.AcquireUserSlotWithWait(c, authSubject.UserID, authSubject.Concurrency, stream, &streamStarted)
	if err != nil {
		googleError(c, http.StatusTooManyRequests, err.Error())
		return
placeholder
	if waitCounted {
		geminiConcurrency.DecrementWaitCount(c.Request.Context(), authSubject.UserID)
		waitCounted = false
placeholder
	// 确保请求取消时也会释放槽位，避免长连接被动中断造成泄漏
	userReleaseFunc = wrapReleaseOnDone(c.Request.Context(), userReleaseFunc)
	if userReleaseFunc != nil {
		defer userReleaseFunc()
placeholder

	// 2) billing eligibility check (after wait)
	if err := h.billingCacheService.CheckBillingEligibility(c.Request.Context(), apiKey.User, apiKey, apiKey.Group, subscription); err != nil {
		status, _, message := billingErrorDetails(err)
		googleError(c, status, message)
		return
placeholder

	// 3) select account (sticky session based on request body)
	parsedReq, _ := service.ParseGatewayRequest(body)
	sessionHash := h.gatewayService.GenerateSessionHash(parsedReq)
	sessionKey := sessionHash
	if sessionHash != "" {
		sessionKey = "gemini:" + sessionHash
placeholder
	const maxAccountSwitches = 3
	switchCount := 0
	failedAccountIDs := make(map[int64]struct{placeholder)
	lastFailoverStatus := 0

	for {
		selection, err := h.gatewayService.SelectAccountWithLoadAwareness(c.Request.Context(), apiKey.GroupID, sessionKey, modelName, failedAccountIDs)
		if err != nil {
			if len(failedAccountIDs) == 0 {
				googleError(c, http.StatusServiceUnavailable, "No available Gemini accounts: "+err.Error())
				return
		placeholder
			handleGeminiFailoverExhausted(c, lastFailoverStatus)
			return
	placeholder
		account := selection.Account
		setOpsSelectedAccount(c, account.ID)

		// 4) account concurrency slot
		accountReleaseFunc := selection.ReleaseFunc
		if !selection.Acquired {
			if selection.WaitPlan == nil {
				googleError(c, http.StatusServiceUnavailable, "No available Gemini accounts")
				return
		placeholder
			accountWaitCounted := false
			canWait, err := geminiConcurrency.IncrementAccountWaitCount(c.Request.Context(), account.ID, selection.WaitPlan.MaxWaiting)
			if err != nil {
				log.Printf("Increment account wait count failed: %v", err)
		placeholder else if !canWait {
				log.Printf("Account wait queue full: account=%d", account.ID)
				googleError(c, http.StatusTooManyRequests, "Too many pending requests, please retry later")
				return
		placeholder
			if err == nil && canWait {
				accountWaitCounted = true
		placeholder
			defer func() {
				if accountWaitCounted {
					geminiConcurrency.DecrementAccountWaitCount(c.Request.Context(), account.ID)
			placeholder
		placeholder()

			accountReleaseFunc, err = geminiConcurrency.AcquireAccountSlotWithWaitTimeout(
				c,
				account.ID,
				selection.WaitPlan.MaxConcurrency,
				selection.WaitPlan.Timeout,
				stream,
				&streamStarted,
			)
			if err != nil {
				googleError(c, http.StatusTooManyRequests, err.Error())
				return
		placeholder
			if accountWaitCounted {
				geminiConcurrency.DecrementAccountWaitCount(c.Request.Context(), account.ID)
				accountWaitCounted = false
		placeholder
			if err := h.gatewayService.BindStickySession(c.Request.Context(), apiKey.GroupID, sessionKey, account.ID); err != nil {
				log.Printf("Bind sticky session failed: %v", err)
		placeholder
	placeholder
		// 账号槽位/等待计数需要在超时或断开时安全回收
		accountReleaseFunc = wrapReleaseOnDone(c.Request.Context(), accountReleaseFunc)

		// 5) forward (根据平台分流)
		var result *service.ForwardResult
		if account.Platform == service.PlatformAntigravity {
			result, err = h.antigravityGatewayService.ForwardGemini(c.Request.Context(), c, account, modelName, action, stream, body)
	placeholder else {
			result, err = h.geminiCompatService.ForwardNative(c.Request.Context(), c, account, modelName, action, stream, body)
	placeholder
		if accountReleaseFunc != nil {
			accountReleaseFunc()
	placeholder
		if err != nil {
			var failoverErr *service.UpstreamFailoverError
			if errors.As(err, &failoverErr) {
				failedAccountIDs[account.ID] = struct{placeholder{placeholder
				if switchCount >= maxAccountSwitches {
					lastFailoverStatus = failoverErr.StatusCode
					handleGeminiFailoverExhausted(c, lastFailoverStatus)
					return
			placeholder
				lastFailoverStatus = failoverErr.StatusCode
				switchCount++
				log.Printf("Gemini account %d: upstream error %d, switching account %d/%d", account.ID, failoverErr.StatusCode, switchCount, maxAccountSwitches)
				continue
		placeholder
			// ForwardNative already wrote the response
			log.Printf("Gemini native forward failed: %v", err)
			return
	placeholder

		// 捕获请求信息（用于异步记录，避免在 goroutine 中访问 gin.Context）
		userAgent := c.GetHeader("User-Agent")
		clientIP := c.ClientIP()

		// 6) record usage async
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

func parseGeminiModelAction(rest string) (model string, action string, err error) {
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return "", "", &pathParseError{"missing path"placeholder
placeholder

	// Standard: {modelplaceholder:{actionplaceholder
	if i := strings.Index(rest, ":"); i > 0 && i < len(rest)-1 {
		return rest[:i], rest[i+1:], nil
placeholder

	// Fallback: {modelplaceholder/{actionplaceholder
	if i := strings.Index(rest, "/"); i > 0 && i < len(rest)-1 {
		return rest[:i], rest[i+1:], nil
placeholder

	return "", "", &pathParseError{"invalid model action path"placeholder
placeholder

func handleGeminiFailoverExhausted(c *gin.Context, statusCode int) {
	status, message := mapGeminiUpstreamError(statusCode)
	googleError(c, status, message)
placeholder

func mapGeminiUpstreamError(statusCode int) (int, string) {
	switch statusCode {
	case 401:
		return http.StatusBadGateway, "Upstream authentication failed, please contact administrator"
	case 403:
		return http.StatusBadGateway, "Upstream access forbidden, please contact administrator"
	case 429:
		return http.StatusTooManyRequests, "Upstream rate limit exceeded, please retry later"
	case 529:
		return http.StatusServiceUnavailable, "Upstream service overloaded, please retry later"
	case 500, 502, 503, 504:
		return http.StatusBadGateway, "Upstream service temporarily unavailable"
	default:
		return http.StatusBadGateway, "Upstream request failed"
placeholder
placeholder

type pathParseError struct{ msg string placeholder

func (e *pathParseError) Error() string { return e.msg placeholder

func googleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    status,
			"message": message,
			"status":  googleapi.HTTPStatusToGoogleStatus(status),
	placeholder,
placeholder)
placeholder

func writeUpstreamResponse(c *gin.Context, res *service.UpstreamHTTPResult) {
	if res == nil {
		googleError(c, http.StatusBadGateway, "Empty upstream response")
		return
placeholder
	for k, vv := range res.Headers {
		// Avoid overriding content-length and hop-by-hop headers.
		if strings.EqualFold(k, "Content-Length") || strings.EqualFold(k, "Transfer-Encoding") || strings.EqualFold(k, "Connection") {
			continue
	placeholder
		for _, v := range vv {
			c.Writer.Header().Add(k, v)
	placeholder
placeholder
	contentType := res.Headers.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
placeholder
	c.Data(res.StatusCode, contentType, res.Body)
placeholder

func shouldFallbackGeminiModels(res *service.UpstreamHTTPResult) bool {
	if res == nil {
		return true
placeholder
	if res.StatusCode != http.StatusUnauthorized && res.StatusCode != http.StatusForbidden {
		return false
placeholder
	if strings.Contains(strings.ToLower(res.Headers.Get("Www-Authenticate")), "insufficient_scope") {
		return true
placeholder
	if strings.Contains(strings.ToLower(string(res.Body)), "insufficient authentication scopes") {
		return true
placeholder
	if strings.Contains(strings.ToLower(string(res.Body)), "access_token_scope_insufficient") {
		return true
placeholder
	return false
placeholder
