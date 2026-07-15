package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
	"github.com/gin-gonic/gin"
)

const (
	snapshotCacheHeader = "X-Snapshot-Cache"
	usageCacheHeader    = "X-Usage-Stats-Cache"
)

type serverTimingResponseWriter struct {
	gin.ResponseWriter
	context *gin.Context
	once    sync.Once
placeholder

func (w *serverTimingResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
placeholder

// ServerTiming collects timing for Admin and User web UI requests when enabled.
func ServerTiming(enabled bool) gin.HandlerFunc {
	if !enabled {
		return func(c *gin.Context) {
			c.Next()
	placeholder
placeholder
	return func(c *gin.Context) {
		if !shouldCollectServerTiming(c) || c.Request == nil {
			c.Next()
			return
	placeholder

		collector := servertiming.New(time.Now())
		c.Request = c.Request.WithContext(servertiming.WithCollector(c.Request.Context(), collector))
		writer := &serverTimingResponseWriter{
			ResponseWriter: c.Writer,
			context:        c,
	placeholder
		c.Writer = writer
		c.Next()
		writer.finalize()
placeholder
placeholder

func (w *serverTimingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
placeholder

func (w *serverTimingResponseWriter) WriteHeaderNow() {
	w.finalize()
	w.ResponseWriter.WriteHeaderNow()
placeholder

func (w *serverTimingResponseWriter) Write(data []byte) (int, error) {
	w.finalize()
	return w.ResponseWriter.Write(data)
placeholder

func (w *serverTimingResponseWriter) WriteString(data string) (int, error) {
	w.finalize()
	return w.ResponseWriter.WriteString(data)
placeholder

func (w *serverTimingResponseWriter) Flush() {
	w.finalize()
	w.ResponseWriter.Flush()
placeholder

func (w *serverTimingResponseWriter) finalize() {
	if w == nil {
		return
placeholder
	w.once.Do(func() {
		if value := ServerTimingHeaderValue(w.context); value != "" {
			w.ResponseWriter.Header().Set(servertiming.HeaderName, value)
	placeholder
placeholder)
placeholder

// ServerTimingHeaderValue returns a timing value only for authorized UI scopes.
// Admins may receive timing for any collected Admin/User UI request. Non-admin
// authenticated users may receive timing only on allowlisted user-facing paths.
// X-User-UI-Request is a scope signal and is never used as authorization.
func ServerTimingHeaderValue(c *gin.Context) string {
	if c == nil || c.Request == nil {
		return ""
placeholder
	role, ok := GetUserRoleFromContext(c)
	if !ok || role == "" {
		return ""
placeholder
	if role != "admin" && !isUserTimingPath(c.Request.URL.Path) {
		return ""
placeholder
	return servertiming.HeaderValue(c.Request.Context(), time.Now(), responseCacheStatus(c.Writer.Header()))
placeholder

// ServerTimingResponseHeader builds the extra header map required by WebSocket upgrades.
func ServerTimingResponseHeader(c *gin.Context) http.Header {
	value := ServerTimingHeaderValue(c)
	if value == "" {
		return nil
placeholder
	return http.Header{servertiming.HeaderName: []string{valueplaceholderplaceholder
placeholder

func shouldCollectServerTiming(c *gin.Context) bool {
	return isAdminUIRequest(c) || isUserUIRequest(c)
placeholder

func isAdminUIRequest(c *gin.Context) bool {
	if c == nil || c.Request == nil || c.Request.URL == nil {
		return false
placeholder
	if strings.TrimSpace(c.GetHeader(servertiming.AdminUIHeader)) == "1" {
		return true
placeholder
	path := strings.TrimSpace(c.Request.URL.Path)
	return path == "/api/v1/admin" || strings.HasPrefix(path, "/api/v1/admin/")
placeholder

func isUserUIRequest(c *gin.Context) bool {
	if c == nil || c.Request == nil || c.Request.URL == nil {
		return false
placeholder
	if strings.TrimSpace(c.GetHeader(servertiming.UserUIHeader)) == "1" {
		return true
placeholder
	return isUserTimingPath(c.Request.URL.Path)
placeholder

// isUserTimingPath reports whether the path is a user-facing web API that may
// emit Server-Timing for authenticated callers (excluding public payment routes).
func isUserTimingPath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
placeholder
	const prefix = "/api/v1"
	if !strings.HasPrefix(path, prefix) {
		return false
placeholder
	rest := strings.TrimPrefix(path, prefix)
	if rest == "" {
		return false
placeholder
	if !strings.HasPrefix(rest, "/") {
		rest = "/" + rest
placeholder

	switch {
	case rest == "/auth/me",
		rest == "/auth/revoke-all-sessions",
		rest == "/auth/oauth/bind-token":
		return true
	case rest == "/user", strings.HasPrefix(rest, "/user/"):
		return true
	case rest == "/keys", strings.HasPrefix(rest, "/keys/"):
		return true
	case rest == "/groups/available", rest == "/groups/rates":
		return true
	case rest == "/channels/available":
		return true
	case rest == "/usage", strings.HasPrefix(rest, "/usage/"):
		return true
	case rest == "/announcements", strings.HasPrefix(rest, "/announcements/"):
		return true
	case rest == "/redeem", strings.HasPrefix(rest, "/redeem/"):
		return true
	case rest == "/subscriptions", strings.HasPrefix(rest, "/subscriptions/"):
		return true
	case rest == "/channel-monitors", strings.HasPrefix(rest, "/channel-monitors/"):
		return true
	case strings.HasPrefix(rest, "/payment/"):
		// Exclude public and webhook payment surfaces.
		if strings.HasPrefix(rest, "/payment/public") || strings.HasPrefix(rest, "/payment/webhook") {
			return false
	placeholder
		return true
	default:
		return false
placeholder
placeholder

func responseCacheStatus(header http.Header) string {
	for _, name := range []string{snapshotCacheHeader, usageCacheHeaderplaceholder {
		switch strings.ToLower(strings.TrimSpace(header.Get(name))) {
		case "hit":
			return "hit"
		case "miss":
			return "miss"
		case "bypass":
			return "bypass"
	placeholder
placeholder
	return "bypass"
placeholder
