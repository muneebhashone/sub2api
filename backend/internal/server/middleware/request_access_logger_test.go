package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

type testLogSink struct {
	mu     sync.Mutex
	events []*logger.LogEvent
placeholder

func (s *testLogSink) WriteLogEvent(event *logger.LogEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
placeholder

func (s *testLogSink) list() []*logger.LogEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*logger.LogEvent, len(s.events))
	copy(out, s.events)
	return out
placeholder

func initMiddlewareTestLogger(t *testing.T) *testLogSink {
	return initMiddlewareTestLoggerWithLevel(t, "debug")
placeholder

func initMiddlewareTestLoggerWithLevel(t *testing.T, level string) *testLogSink {
placeholder
	level = strings.TrimSpace(level)
	if level == "" {
		level = "debug"
placeholder
	if err := logger.Init(logger.InitOptions{
		Level:       level,
		Format:      "json",
		ServiceName: "sub2api",
		Environment: "test",
		Output: logger.OutputOptions{
			ToStdout: false,
			ToFile:   false,
	placeholder,
placeholder); err != nil {
		t.Fatalf("init logger: %v", err)
placeholder
	sink := &testLogSink{placeholder
	logger.SetSink(sink)
	t.Cleanup(func() {
		logger.SetSink(nil)
placeholder)
	return sink
placeholder

func TestRequestLogger_GenerateAndPropagateRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLogger())
	r.GET("/t", func(c *gin.Context) {
		reqID, ok := c.Request.Context().Value(ctxkey.RequestID).(string)
		if !ok || reqID == "" {
			t.Fatalf("request_id missing in context")
	placeholder
		if got := c.Writer.Header().Get(requestIDHeader); got != reqID {
			t.Fatalf("response header request_id mismatch, header=%q ctx=%q", got, reqID)
	placeholder
		c.Status(http.StatusOK)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
placeholder
	if w.Header().Get(requestIDHeader) == "" {
		t.Fatalf("X-Request-ID should be set")
placeholder
placeholder

func TestRequestLogger_KeepIncomingRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLogger())
	r.GET("/t", func(c *gin.Context) {
		reqID, _ := c.Request.Context().Value(ctxkey.RequestID).(string)
		if reqID != "rid-fixed" {
			t.Fatalf("request_id=%q, want rid-fixed", reqID)
	placeholder
		c.Status(http.StatusOK)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	req.Header.Set(requestIDHeader, "rid-fixed")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
placeholder
	if got := w.Header().Get(requestIDHeader); got != "rid-fixed" {
		t.Fatalf("header=%q, want rid-fixed", got)
placeholder
placeholder

func TestRequestLoggerBoundsIncomingRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestLogger())
	r.GET("/t", func(c *gin.Context) {
		reqID, _ := c.Request.Context().Value(ctxkey.RequestID).(string)
		if len(reqID) != 36 {
			t.Fatalf("request_id length=%d", len(reqID))
	placeholder
		c.Status(http.StatusOK)
placeholder)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/t", nil)
	req.Header.Set(requestIDHeader, strings.Repeat("r", 1024))
	r.ServeHTTP(w, req)
	if got := len(w.Header().Get(requestIDHeader)); got != 36 {
		t.Fatalf("response request_id length=%d", got)
placeholder
placeholder

func TestLogger_AccessLogIncludesCoreFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sink := initMiddlewareTestLogger(t)

	r := gin.New()
	r.Use(Logger())
	r.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ctxkey.AccountID, int64(101))
		ctx = context.WithValue(ctx, ctxkey.Platform, "openai")
		ctx = context.WithValue(ctx, ctxkey.Model, "gpt-5")
		c.Request = c.Request.WithContext(ctx)
		c.Next()
placeholder)
	r.GET("/api/test", func(c *gin.Context) {
		c.Status(http.StatusCreated)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d", w.Code)
placeholder

	events := sink.list()
	if len(events) == 0 {
		t.Fatalf("expected at least one log event")
placeholder
	found := false
	for _, event := range events {
		if event == nil || event.Message != "http request completed" {
			continue
	placeholder
		found = true
		switch v := event.Fields["status_code"].(type) {
		case int:
			if v != http.StatusCreated {
				t.Fatalf("status_code field mismatch: %v", v)
		placeholder
		case int64:
			if v != int64(http.StatusCreated) {
				t.Fatalf("status_code field mismatch: %v", v)
		placeholder
		default:
			t.Fatalf("status_code type mismatch: %T", v)
	placeholder
		switch v := event.Fields["account_id"].(type) {
		case int64:
			if v != 101 {
				t.Fatalf("account_id field mismatch: %v", v)
		placeholder
		case int:
			if v != 101 {
				t.Fatalf("account_id field mismatch: %v", v)
		placeholder
		default:
			t.Fatalf("account_id type mismatch: %T", v)
	placeholder
		if event.Fields["platform"] != "openai" || event.Fields["model"] != "gpt-5" {
			t.Fatalf("platform/model mismatch: %+v", event.Fields)
	placeholder
placeholder
	if !found {
		t.Fatalf("access log event not found")
placeholder
placeholder

func TestLogger_IngressRejectRemainsInStandardAccessLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sink := initMiddlewareTestLogger(t)
	r := gin.New()
	r.Use(Logger())
	r.GET("/v1/messages", func(c *gin.Context) {
		MarkIngressRejected(c, IngressRejectInvalidAPIKey)
		c.Status(http.StatusUnauthorized)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/messages", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", w.Code)
placeholder
	events := sink.list()
	if len(events) != 1 {
		t.Fatalf("events=%d, want 1", len(events))
placeholder
	if got := events[0].Fields["ingress_reject_reason"]; got != string(IngressRejectInvalidAPIKey) {
		t.Fatalf("ingress_reject_reason=%v", got)
placeholder
	if got, _ := events[0].Fields[logger.OpsSystemLogSkipField].(bool); !got {
		t.Fatalf("%s must be true", logger.OpsSystemLogSkipField)
placeholder
placeholder

func TestLogger_AccessLogUsesForwardedClientIPFromTrustedProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sink := initMiddlewareTestLogger(t)

	r := gin.New()
	if err := r.SetTrustedProxies([]string{"104.23.251.120"placeholder); err != nil {
		t.Fatalf("set trusted proxies: %v", err)
placeholder
	r.Use(Logger())
	r.GET("/api/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.RemoteAddr = "104.23.251.120:443"
	req.Header.Set("X-Forwarded-For", "203.0.113.42")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
placeholder

	for _, event := range sink.list() {
		if event == nil || event.Message != "http request completed" {
			continue
	placeholder
		if got := event.Fields["client_ip"]; got != "203.0.113.42" {
			t.Fatalf("client_ip=%q, want real forwarded ip", got)
	placeholder
		return
placeholder
	t.Fatalf("access log event not found")
placeholder

func TestLogger_HealthPathSkipped(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sink := initMiddlewareTestLogger(t)

	r := gin.New()
	r.Use(Logger())
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d", w.Code)
placeholder
	if len(sink.list()) != 0 {
		t.Fatalf("health endpoint should not write access log")
placeholder
placeholder

func TestLogger_AccessLogDroppedWhenLevelWarn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sink := initMiddlewareTestLoggerWithLevel(t, "warn")

	r := gin.New()
	r.Use(RequestLogger())
	r.Use(Logger())
	r.GET("/api/test", func(c *gin.Context) {
		c.Status(http.StatusCreated)
placeholder)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d", w.Code)
placeholder

	events := sink.list()
	for _, event := range events {
		if event != nil && event.Message == "http request completed" {
			t.Fatalf("access log should not be indexed when level=warn: %+v", event)
	placeholder
placeholder
placeholder
