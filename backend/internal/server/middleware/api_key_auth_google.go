package middleware

import (
	"errors"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/googleapi"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// ApiKeyAuthGoogle is a Google-style error wrapper for API key auth.
func ApiKeyAuthGoogle(apiKeyService *service.ApiKeyService, cfg *config.Config) gin.HandlerFunc {
	return ApiKeyAuthWithSubscriptionGoogle(apiKeyService, nil, cfg)
placeholder

// ApiKeyAuthWithSubscriptionGoogle behaves like ApiKeyAuthWithSubscription but returns Google-style errors:
// {"error":{"code":401,"message":"...","status":"UNAUTHENTICATED"placeholderplaceholder
//
// It is intended for Gemini native endpoints (/v1beta) to match Gemini SDK expectations.
func ApiKeyAuthWithSubscriptionGoogle(apiKeyService *service.ApiKeyService, subscriptionService *service.SubscriptionService, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if v := strings.TrimSpace(c.Query("api_key")); v != "" {
			abortWithGoogleError(c, 400, "Query parameter api_key is deprecated. Use Authorization header or key instead.")
			return
	placeholder
		apiKeyString := extractAPIKeyFromRequest(c)
		if apiKeyString == "" {
			abortWithGoogleError(c, 401, "API key is required")
			return
	placeholder

		apiKey, err := apiKeyService.GetByKey(c.Request.Context(), apiKeyString)
		if err != nil {
			if errors.Is(err, service.ErrApiKeyNotFound) {
				abortWithGoogleError(c, 401, "Invalid API key")
				return
		placeholder
			abortWithGoogleError(c, 500, "Failed to validate API key")
			return
	placeholder

		if !apiKey.IsActive() {
			abortWithGoogleError(c, 401, "API key is disabled")
			return
	placeholder
		if apiKey.User == nil {
			abortWithGoogleError(c, 401, "User associated with API key not found")
			return
	placeholder
		if !apiKey.User.IsActive() {
			abortWithGoogleError(c, 401, "User account is not active")
			return
	placeholder

		// 简易模式：跳过余额和订阅检查
		if cfg.RunMode == config.RunModeSimple {
			c.Set(string(ContextKeyApiKey), apiKey)
			c.Set(string(ContextKeyUser), AuthSubject{
				UserID:      apiKey.User.ID,
				Concurrency: apiKey.User.Concurrency,
		placeholder)
			c.Set(string(ContextKeyUserRole), apiKey.User.Role)
			c.Next()
			return
	placeholder

		isSubscriptionType := apiKey.Group != nil && apiKey.Group.IsSubscriptionType()
		if isSubscriptionType && subscriptionService != nil {
			subscription, err := subscriptionService.GetActiveSubscription(
				c.Request.Context(),
				apiKey.User.ID,
				apiKey.Group.ID,
			)
			if err != nil {
				abortWithGoogleError(c, 403, "No active subscription found for this group")
				return
		placeholder
			if err := subscriptionService.ValidateSubscription(c.Request.Context(), subscription); err != nil {
				abortWithGoogleError(c, 403, err.Error())
				return
		placeholder
			_ = subscriptionService.CheckAndActivateWindow(c.Request.Context(), subscription)
			_ = subscriptionService.CheckAndResetWindows(c.Request.Context(), subscription)
			if err := subscriptionService.CheckUsageLimits(c.Request.Context(), subscription, apiKey.Group, 0); err != nil {
				abortWithGoogleError(c, 429, err.Error())
				return
		placeholder
			c.Set(string(ContextKeySubscription), subscription)
	placeholder else {
			if apiKey.User.Balance <= 0 {
				abortWithGoogleError(c, 403, "Insufficient account balance")
				return
		placeholder
	placeholder

		c.Set(string(ContextKeyApiKey), apiKey)
		c.Set(string(ContextKeyUser), AuthSubject{
			UserID:      apiKey.User.ID,
			Concurrency: apiKey.User.Concurrency,
	placeholder)
		c.Set(string(ContextKeyUserRole), apiKey.User.Role)
		c.Next()
placeholder
placeholder

func extractAPIKeyFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" && strings.TrimSpace(parts[1]) != "" {
			return strings.TrimSpace(parts[1])
	placeholder
placeholder
	if v := strings.TrimSpace(c.GetHeader("x-api-key")); v != "" {
		return v
placeholder
	if v := strings.TrimSpace(c.GetHeader("x-goog-api-key")); v != "" {
		return v
placeholder
	if allowGoogleQueryKey(c.Request.URL.Path) {
		if v := strings.TrimSpace(c.Query("key")); v != "" {
			return v
	placeholder
placeholder
	return ""
placeholder

func allowGoogleQueryKey(path string) bool {
	return strings.HasPrefix(path, "/v1beta") || strings.HasPrefix(path, "/antigravity/v1beta")
placeholder

func abortWithGoogleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    status,
			"message": message,
			"status":  googleapi.HTTPStatusToGoogleStatus(status),
	placeholder,
placeholder)
	c.Abort()
placeholder
