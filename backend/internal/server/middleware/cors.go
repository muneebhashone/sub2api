package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

var corsWarningOnce sync.Once

// CORS 跨域中间件
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	allowedOrigins := normalizeOrigins(cfg.AllowedOrigins)
	allowAll := false
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAll = true
			break
	placeholder
placeholder
	wildcardWithSpecific := allowAll && len(allowedOrigins) > 1
	if wildcardWithSpecific {
		allowedOrigins = []string{"*"placeholder
placeholder
	allowCredentials := cfg.AllowCredentials

	corsWarningOnce.Do(func() {
		if len(allowedOrigins) == 0 {
			log.Println("Warning: CORS allowed_origins not configured; cross-origin requests will be rejected.")
	placeholder
		if wildcardWithSpecific {
			log.Println("Warning: CORS allowed_origins includes '*'; wildcard will take precedence over explicit origins.")
	placeholder
		if allowAll && allowCredentials {
			log.Println("Warning: CORS allowed_origins set to '*', disabling allow_credentials.")
	placeholder
placeholder)
	if allowAll && allowCredentials {
		allowCredentials = false
placeholder

	allowedSet := make(map[string]struct{placeholder, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		if origin == "" || origin == "*" {
			continue
	placeholder
		allowedSet[origin] = struct{placeholder{placeholder
placeholder

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		originAllowed := allowAll
		if origin != "" && !allowAll {
			_, originAllowed = allowedSet[origin]
	placeholder

		if originAllowed {
			if allowAll {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		placeholder else if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Add("Vary", "Origin")
		placeholder
			if allowCredentials {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		placeholder
	placeholder

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// 处理预检请求
		if c.Request.Method == http.MethodOptions {
			if originAllowed {
				c.AbortWithStatus(http.StatusNoContent)
		placeholder else {
				c.AbortWithStatus(http.StatusForbidden)
		placeholder
			return
	placeholder

		c.Next()
placeholder
placeholder

func normalizeOrigins(values []string) []string {
	if len(values) == 0 {
		return nil
placeholder
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
	placeholder
		normalized = append(normalized, trimmed)
placeholder
	return normalized
placeholder
