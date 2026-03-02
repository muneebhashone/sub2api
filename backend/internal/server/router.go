package server

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/server/routes"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/Wei-Shaw/sub2api/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// extractOrigin returns the scheme+host origin from rawURL, or "" on error.
// Only http and https schemes are accepted; other values (e.g. "//host/path") return "".
func extractOrigin(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
placeholder
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		return ""
placeholder
	if u.Scheme != "http" && u.Scheme != "https" {
		return ""
placeholder
	return u.Scheme + "://" + u.Host
placeholder

const paymentOriginFetchTimeout = 5 * time.Second

// SetupRouter 配置路由器中间件和路由
func SetupRouter(
	r *gin.Engine,
	handlers *handler.Handlers,
	jwtAuth middleware2.JWTAuthMiddleware,
	adminAuth middleware2.AdminAuthMiddleware,
	apiKeyAuth middleware2.APIKeyAuthMiddleware,
	apiKeyService *service.APIKeyService,
	subscriptionService *service.SubscriptionService,
	opsService *service.OpsService,
	settingService *service.SettingService,
	cfg *config.Config,
	redisClient *redis.Client,
) *gin.Engine {
	// 缓存 iframe 页面的 origin 列表，用于动态注入 CSP frame-src
	// 包含 purchase_subscription_url 和所有 custom_menu_items 的 origin（去重）
	var cachedFrameOrigins atomic.Pointer[[]string]
	emptyOrigins := []string{placeholder
	cachedFrameOrigins.Store(&emptyOrigins)

	refreshFrameOrigins := func() {
		ctx, cancel := context.WithTimeout(context.Background(), paymentOriginFetchTimeout)
		defer cancel()
		settings, err := settingService.GetPublicSettings(ctx)
		if err != nil {
			// 获取失败时保留已有缓存，避免 frame-src 被意外清空
			return
	placeholder

		seen := make(map[string]struct{placeholder)
		var origins []string

		// purchase subscription URL
		if settings.PurchaseSubscriptionEnabled {
			if origin := extractOrigin(settings.PurchaseSubscriptionURL); origin != "" {
				if _, ok := seen[origin]; !ok {
					seen[origin] = struct{placeholder{placeholder
					origins = append(origins, origin)
			placeholder
		placeholder
	placeholder

		// custom menu items
		if raw := strings.TrimSpace(settings.CustomMenuItems); raw != "" && raw != "[]" {
			var items []struct {
				URL string `json:"url"`
		placeholder
			if err := json.Unmarshal([]byte(raw), &items); err == nil {
				for _, item := range items {
					if origin := extractOrigin(item.URL); origin != "" {
						if _, ok := seen[origin]; !ok {
							seen[origin] = struct{placeholder{placeholder
							origins = append(origins, origin)
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder

		cachedFrameOrigins.Store(&origins)
placeholder
	refreshFrameOrigins() // 启动时初始化

	// 应用中间件
	r.Use(middleware2.RequestLogger())
	r.Use(middleware2.Logger())
	r.Use(middleware2.CORS(cfg.CORS))
	r.Use(middleware2.SecurityHeaders(cfg.Security.CSP, func() []string {
		if p := cachedFrameOrigins.Load(); p != nil {
			return *p
	placeholder
		return nil
placeholder))

	// Serve embedded frontend with settings injection if available
	if web.HasEmbeddedFrontend() {
		frontendServer, err := web.NewFrontendServer(settingService)
		if err != nil {
			log.Printf("Warning: Failed to create frontend server with settings injection: %v, using legacy mode", err)
			r.Use(web.ServeEmbeddedFrontend())
			settingService.SetOnUpdateCallback(refreshFrameOrigins)
	placeholder else {
			// Register combined callback: invalidate HTML cache + refresh frame origins
			settingService.SetOnUpdateCallback(func() {
				frontendServer.InvalidateCache()
				refreshFrameOrigins()
		placeholder)
			r.Use(frontendServer.Middleware())
	placeholder
placeholder else {
		settingService.SetOnUpdateCallback(refreshFrameOrigins)
placeholder

	// 注册路由
	registerRoutes(r, handlers, jwtAuth, adminAuth, apiKeyAuth, apiKeyService, subscriptionService, opsService, cfg, redisClient)

	return r
placeholder

// registerRoutes 注册所有 HTTP 路由
func registerRoutes(
	r *gin.Engine,
	h *handler.Handlers,
	jwtAuth middleware2.JWTAuthMiddleware,
	adminAuth middleware2.AdminAuthMiddleware,
	apiKeyAuth middleware2.APIKeyAuthMiddleware,
	apiKeyService *service.APIKeyService,
	subscriptionService *service.SubscriptionService,
	opsService *service.OpsService,
	cfg *config.Config,
	redisClient *redis.Client,
) {
	// 通用路由（健康检查、状态等）
	routes.RegisterCommonRoutes(r)

	// API v1
	v1 := r.Group("/api/v1")

	// 注册各模块路由
	routes.RegisterAuthRoutes(v1, h, jwtAuth, redisClient)
	routes.RegisterUserRoutes(v1, h, jwtAuth)
	routes.RegisterSoraClientRoutes(v1, h, jwtAuth)
	routes.RegisterAdminRoutes(v1, h, adminAuth)
	routes.RegisterGatewayRoutes(r, h, apiKeyAuth, apiKeyService, subscriptionService, opsService, cfg)
placeholder
