package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	opsModelKey       = "ops_model"
	opsStreamKey      = "ops_stream"
	opsRequestBodyKey = "ops_request_body"
	opsAccountIDKey   = "ops_account_id"
)

const (
	opsErrorLogTimeout      = 5 * time.Second
	opsErrorLogDrainTimeout = 10 * time.Second

	opsErrorLogMinWorkerCount = 4
	opsErrorLogMaxWorkerCount = 32

	opsErrorLogQueueSizePerWorker = 128
	opsErrorLogMinQueueSize       = 256
	opsErrorLogMaxQueueSize       = 8192
)

type opsErrorLogJob struct {
	ops         *service.OpsService
	entry       *service.OpsInsertErrorLogInput
	requestBody []byte
placeholder

var (
	opsErrorLogOnce  sync.Once
	opsErrorLogQueue chan opsErrorLogJob

	opsErrorLogStopOnce  sync.Once
	opsErrorLogWorkersWg sync.WaitGroup
	opsErrorLogMu        sync.RWMutex
	opsErrorLogStopping  bool
	opsErrorLogQueueLen  atomic.Int64
	opsErrorLogEnqueued  atomic.Int64
	opsErrorLogDropped   atomic.Int64
	opsErrorLogProcessed atomic.Int64

	opsErrorLogLastDropLogAt atomic.Int64

	opsErrorLogShutdownCh   = make(chan struct{placeholder)
	opsErrorLogShutdownOnce sync.Once
	opsErrorLogDrained      atomic.Bool
)

func startOpsErrorLogWorkers() {
	opsErrorLogMu.Lock()
	defer opsErrorLogMu.Unlock()

	if opsErrorLogStopping {
		return
placeholder

	workerCount, queueSize := opsErrorLogConfig()
	opsErrorLogQueue = make(chan opsErrorLogJob, queueSize)
	opsErrorLogQueueLen.Store(0)

	opsErrorLogWorkersWg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer opsErrorLogWorkersWg.Done()
			for job := range opsErrorLogQueue {
				opsErrorLogQueueLen.Add(-1)
				if job.ops == nil || job.entry == nil {
					continue
			placeholder
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[OpsErrorLogger] worker panic: %v\n%s", r, debug.Stack())
					placeholder
				placeholder()
					ctx, cancel := context.WithTimeout(context.Background(), opsErrorLogTimeout)
					_ = job.ops.RecordError(ctx, job.entry, job.requestBody)
					cancel()
					opsErrorLogProcessed.Add(1)
			placeholder()
		placeholder
	placeholder()
placeholder
placeholder

func enqueueOpsErrorLog(ops *service.OpsService, entry *service.OpsInsertErrorLogInput, requestBody []byte) {
	if ops == nil || entry == nil {
		return
placeholder
	select {
	case <-opsErrorLogShutdownCh:
		return
	default:
placeholder

	opsErrorLogMu.RLock()
	stopping := opsErrorLogStopping
	opsErrorLogMu.RUnlock()
	if stopping {
		return
placeholder

	opsErrorLogOnce.Do(startOpsErrorLogWorkers)

	opsErrorLogMu.RLock()
	defer opsErrorLogMu.RUnlock()
	if opsErrorLogStopping || opsErrorLogQueue == nil {
		return
placeholder

	select {
	case opsErrorLogQueue <- opsErrorLogJob{ops: ops, entry: entry, requestBody: requestBodyplaceholder:
		opsErrorLogQueueLen.Add(1)
		opsErrorLogEnqueued.Add(1)
	default:
		// Queue is full; drop to avoid blocking request handling.
		opsErrorLogDropped.Add(1)
		maybeLogOpsErrorLogDrop()
placeholder
placeholder

func StopOpsErrorLogWorkers() bool {
	opsErrorLogStopOnce.Do(func() {
		opsErrorLogShutdownOnce.Do(func() {
			close(opsErrorLogShutdownCh)
	placeholder)
		opsErrorLogDrained.Store(stopOpsErrorLogWorkers())
placeholder)
	return opsErrorLogDrained.Load()
placeholder

func stopOpsErrorLogWorkers() bool {
	opsErrorLogMu.Lock()
	opsErrorLogStopping = true
	ch := opsErrorLogQueue
	if ch != nil {
		close(ch)
placeholder
	opsErrorLogQueue = nil
	opsErrorLogMu.Unlock()

	if ch == nil {
		opsErrorLogQueueLen.Store(0)
		return true
placeholder

	done := make(chan struct{placeholder)
	go func() {
		opsErrorLogWorkersWg.Wait()
		close(done)
placeholder()

	select {
	case <-done:
		opsErrorLogQueueLen.Store(0)
		return true
	case <-time.After(opsErrorLogDrainTimeout):
		return false
placeholder
placeholder

func OpsErrorLogQueueLength() int64 {
	return opsErrorLogQueueLen.Load()
placeholder

func OpsErrorLogQueueCapacity() int {
	opsErrorLogMu.RLock()
	ch := opsErrorLogQueue
	opsErrorLogMu.RUnlock()
	if ch == nil {
		return 0
placeholder
	return cap(ch)
placeholder

func OpsErrorLogDroppedTotal() int64 {
	return opsErrorLogDropped.Load()
placeholder

func OpsErrorLogEnqueuedTotal() int64 {
	return opsErrorLogEnqueued.Load()
placeholder

func OpsErrorLogProcessedTotal() int64 {
	return opsErrorLogProcessed.Load()
placeholder

func maybeLogOpsErrorLogDrop() {
	now := time.Now().Unix()

	for {
		last := opsErrorLogLastDropLogAt.Load()
		if last != 0 && now-last < 60 {
			return
	placeholder
		if opsErrorLogLastDropLogAt.CompareAndSwap(last, now) {
			break
	placeholder
placeholder

	queued := opsErrorLogQueueLen.Load()
	queueCap := OpsErrorLogQueueCapacity()

	log.Printf(
		"[OpsErrorLogger] queue is full; dropping logs (queued=%d cap=%d enqueued_total=%d dropped_total=%d processed_total=%d)",
		queued,
		queueCap,
		opsErrorLogEnqueued.Load(),
		opsErrorLogDropped.Load(),
		opsErrorLogProcessed.Load(),
	)
placeholder

func opsErrorLogConfig() (workerCount int, queueSize int) {
	workerCount = runtime.GOMAXPROCS(0) * 2
	if workerCount < opsErrorLogMinWorkerCount {
		workerCount = opsErrorLogMinWorkerCount
placeholder
	if workerCount > opsErrorLogMaxWorkerCount {
		workerCount = opsErrorLogMaxWorkerCount
placeholder

	queueSize = workerCount * opsErrorLogQueueSizePerWorker
	if queueSize < opsErrorLogMinQueueSize {
		queueSize = opsErrorLogMinQueueSize
placeholder
	if queueSize > opsErrorLogMaxQueueSize {
		queueSize = opsErrorLogMaxQueueSize
placeholder

	return workerCount, queueSize
placeholder

func setOpsRequestContext(c *gin.Context, model string, stream bool, requestBody []byte) {
	if c == nil {
		return
placeholder
	c.Set(opsModelKey, model)
	c.Set(opsStreamKey, stream)
	if len(requestBody) > 0 {
		c.Set(opsRequestBodyKey, requestBody)
placeholder
placeholder

func setOpsSelectedAccount(c *gin.Context, accountID int64) {
	if c == nil || accountID <= 0 {
		return
placeholder
	c.Set(opsAccountIDKey, accountID)
placeholder

type opsCaptureWriter struct {
	gin.ResponseWriter
	limit int
	buf   bytes.Buffer
placeholder

func (w *opsCaptureWriter) Write(b []byte) (int, error) {
	if w.Status() >= 400 && w.limit > 0 && w.buf.Len() < w.limit {
		remaining := w.limit - w.buf.Len()
		if len(b) > remaining {
			_, _ = w.buf.Write(b[:remaining])
	placeholder else {
			_, _ = w.buf.Write(b)
	placeholder
placeholder
	return w.ResponseWriter.Write(b)
placeholder

func (w *opsCaptureWriter) WriteString(s string) (int, error) {
	if w.Status() >= 400 && w.limit > 0 && w.buf.Len() < w.limit {
		remaining := w.limit - w.buf.Len()
		if len(s) > remaining {
			_, _ = w.buf.WriteString(s[:remaining])
	placeholder else {
			_, _ = w.buf.WriteString(s)
	placeholder
placeholder
	return w.ResponseWriter.WriteString(s)
placeholder

// OpsErrorLoggerMiddleware records error responses (status >= 400) into ops_error_logs.
//
// Notes:
// - It buffers response bodies only when status >= 400 to avoid overhead for successful traffic.
// - Streaming errors after the response has started (SSE) may still need explicit logging.
func OpsErrorLoggerMiddleware(ops *service.OpsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &opsCaptureWriter{ResponseWriter: c.Writer, limit: 64 * placeholder
		c.Writer = w
		c.Next()

		status := c.Writer.Status()
		if status < 400 {
			return
	placeholder
		if ops == nil {
			return
	placeholder
		if !ops.IsMonitoringEnabled(c.Request.Context()) {
			return
	placeholder

		body := w.buf.Bytes()
		parsed := parseOpsErrorResponse(body)

		apiKey, _ := middleware2.GetAPIKeyFromContext(c)

		clientRequestID, _ := c.Request.Context().Value(ctxkey.ClientRequestID).(string)

		model, _ := c.Get(opsModelKey)
		streamV, _ := c.Get(opsStreamKey)
		accountIDV, _ := c.Get(opsAccountIDKey)

		var modelName string
		if s, ok := model.(string); ok {
			modelName = s
	placeholder
		stream := false
		if b, ok := streamV.(bool); ok {
			stream = b
	placeholder
		var accountID *int64
		if v, ok := accountIDV.(int64); ok && v > 0 {
			accountID = &v
	placeholder

		fallbackPlatform := guessPlatformFromPath(c.Request.URL.Path)
		platform := resolveOpsPlatform(apiKey, fallbackPlatform)

		requestID := c.Writer.Header().Get("X-Request-Id")
		if requestID == "" {
			requestID = c.Writer.Header().Get("x-request-id")
	placeholder

		phase := classifyOpsPhase(parsed.ErrorType, parsed.Message, parsed.Code)
		isBusinessLimited := classifyOpsIsBusinessLimited(parsed.ErrorType, phase, parsed.Code, status, parsed.Message)

		errorOwner := classifyOpsErrorOwner(phase, parsed.Message)
		errorSource := classifyOpsErrorSource(phase, parsed.Message)

		entry := &service.OpsInsertErrorLogInput{
			RequestID:       requestID,
			ClientRequestID: clientRequestID,

			AccountID: accountID,
			Platform:  platform,
			Model:     modelName,
			RequestPath: func() string {
				if c.Request != nil && c.Request.URL != nil {
					return c.Request.URL.Path
			placeholder
				return ""
		placeholder(),
			Stream:    stream,
			UserAgent: c.GetHeader("User-Agent"),

			ErrorPhase:        phase,
			ErrorType:         normalizeOpsErrorType(parsed.ErrorType, parsed.Code),
			Severity:          classifyOpsSeverity(parsed.ErrorType, status),
			StatusCode:        status,
			IsBusinessLimited: isBusinessLimited,

			ErrorMessage: parsed.Message,
			// Keep the full captured error body (capture is already capped at 64KB) so the
			// service layer can sanitize JSON before truncating for storage.
			ErrorBody:   string(body),
			ErrorSource: errorSource,
			ErrorOwner:  errorOwner,

			IsRetryable: classifyOpsIsRetryable(parsed.ErrorType, status),
			RetryCount:  0,
			CreatedAt:   time.Now(),
	placeholder

		// Capture upstream error context set by gateway services (if present).
		// This does NOT affect the client response; it enriches Ops troubleshooting data.
		{
			if v, ok := c.Get(service.OpsUpstreamStatusCodeKey); ok {
				switch t := v.(type) {
				case int:
					if t > 0 {
						code := t
						entry.UpstreamStatusCode = &code
				placeholder
				case int64:
					if t > 0 {
						code := int(t)
						entry.UpstreamStatusCode = &code
				placeholder
			placeholder
		placeholder
			if v, ok := c.Get(service.OpsUpstreamErrorMessageKey); ok {
				if s, ok := v.(string); ok {
					if msg := strings.TrimSpace(s); msg != "" {
						entry.UpstreamErrorMessage = &msg
				placeholder
			placeholder
		placeholder
			if v, ok := c.Get(service.OpsUpstreamErrorDetailKey); ok {
				if s, ok := v.(string); ok {
					if detail := strings.TrimSpace(s); detail != "" {
						entry.UpstreamErrorDetail = &detail
				placeholder
			placeholder
		placeholder
	placeholder

		if apiKey != nil {
			entry.APIKeyID = &apiKey.ID
			if apiKey.User != nil {
				entry.UserID = &apiKey.User.ID
		placeholder
			if apiKey.GroupID != nil {
				entry.GroupID = apiKey.GroupID
		placeholder
			// Prefer group platform if present (more stable than inferring from path).
			if apiKey.Group != nil && apiKey.Group.Platform != "" {
				entry.Platform = apiKey.Group.Platform
		placeholder
	placeholder

		var clientIP string
		if ip := strings.TrimSpace(c.ClientIP()); ip != "" {
			clientIP = ip
			entry.ClientIP = &clientIP
	placeholder

		var requestBody []byte
		if v, ok := c.Get(opsRequestBodyKey); ok {
			if b, ok := v.([]byte); ok && len(b) > 0 {
				requestBody = b
		placeholder
	placeholder
		// Persist only a minimal, whitelisted set of request headers to improve retry fidelity.
		// Do NOT store Authorization/Cookie/etc.
		entry.RequestHeadersJSON = extractOpsRetryRequestHeaders(c)

		enqueueOpsErrorLog(ops, entry, requestBody)
placeholder
placeholder

var opsRetryRequestHeaderAllowlist = []string{
	"anthropic-beta",
	"anthropic-version",
placeholder

func extractOpsRetryRequestHeaders(c *gin.Context) *string {
	if c == nil || c.Request == nil {
		return nil
placeholder

	headers := make(map[string]string, 4)
	for _, key := range opsRetryRequestHeaderAllowlist {
		v := strings.TrimSpace(c.GetHeader(key))
		if v == "" {
			continue
	placeholder
		// Keep headers small even if a client sends something unexpected.
		headers[key] = truncateString(v, 512)
placeholder
	if len(headers) == 0 {
		return nil
placeholder

	raw, err := json.Marshal(headers)
	if err != nil {
		return nil
placeholder
	s := string(raw)
	return &s
placeholder

type parsedOpsError struct {
	ErrorType string
	Message   string
	Code      string
placeholder

func parseOpsErrorResponse(body []byte) parsedOpsError {
	if len(body) == 0 {
		return parsedOpsError{placeholder
placeholder

	// Fast path: attempt to decode into a generic map.
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return parsedOpsError{Message: truncateString(string(body), 1024)placeholder
placeholder

	// Claude/OpenAI-style gateway error: { type:"error", error:{ type, message placeholder placeholder
	if errObj, ok := m["error"].(map[string]any); ok {
		t, _ := errObj["type"].(string)
		msg, _ := errObj["message"].(string)
		// Gemini googleError also uses "error": { code, message, status placeholder
		if msg == "" {
			if v, ok := errObj["message"]; ok {
				msg, _ = v.(string)
		placeholder
	placeholder
		if t == "" {
			// Gemini error does not have "type" field.
			t = "api_error"
	placeholder
		// For gemini error, capture numeric code as string for business-limited mapping if needed.
		var code string
		if v, ok := errObj["code"]; ok {
			switch n := v.(type) {
			case float64:
				code = strconvItoa(int(n))
			case int:
				code = strconvItoa(n)
		placeholder
	placeholder
		return parsedOpsError{ErrorType: t, Message: msg, Code: codeplaceholder
placeholder

	// APIKeyAuth-style: { code:"INSUFFICIENT_BALANCE", message:"..." placeholder
	code, _ := m["code"].(string)
	msg, _ := m["message"].(string)
	if code != "" || msg != "" {
		return parsedOpsError{ErrorType: "api_error", Message: msg, Code: codeplaceholder
placeholder

	return parsedOpsError{Message: truncateString(string(body), 1024)placeholder
placeholder

func resolveOpsPlatform(apiKey *service.APIKey, fallback string) string {
	if apiKey != nil && apiKey.Group != nil && apiKey.Group.Platform != "" {
		return apiKey.Group.Platform
placeholder
	return fallback
placeholder

func guessPlatformFromPath(path string) string {
	p := strings.ToLower(path)
	switch {
	case strings.HasPrefix(p, "/antigravity/"):
		return service.PlatformAntigravity
	case strings.HasPrefix(p, "/v1beta/"):
		return service.PlatformGemini
	case strings.Contains(p, "/responses"):
		return service.PlatformOpenAI
	default:
		return ""
placeholder
placeholder

func normalizeOpsErrorType(errType string, code string) string {
	if errType != "" {
		return errType
placeholder
	switch strings.TrimSpace(code) {
	case "INSUFFICIENT_BALANCE":
		return "billing_error"
	case "USAGE_LIMIT_EXCEEDED", "SUBSCRIPTION_NOT_FOUND", "SUBSCRIPTION_INVALID":
		return "subscription_error"
	default:
		return "api_error"
placeholder
placeholder

func classifyOpsPhase(errType, message, code string) string {
	msg := strings.ToLower(message)
	switch strings.TrimSpace(code) {
	case "INSUFFICIENT_BALANCE", "USAGE_LIMIT_EXCEEDED", "SUBSCRIPTION_NOT_FOUND", "SUBSCRIPTION_INVALID":
		return "billing"
placeholder

	switch errType {
	case "authentication_error":
		return "auth"
	case "billing_error", "subscription_error":
		return "billing"
	case "rate_limit_error":
		if strings.Contains(msg, "concurrency") || strings.Contains(msg, "pending") || strings.Contains(msg, "queue") {
			return "concurrency"
	placeholder
		return "upstream"
	case "invalid_request_error":
		return "response"
	case "upstream_error", "overloaded_error":
		return "upstream"
	case "api_error":
		if strings.Contains(msg, "no available accounts") {
			return "scheduling"
	placeholder
		return "internal"
	default:
		return "internal"
placeholder
placeholder

func classifyOpsSeverity(errType string, status int) string {
	switch errType {
	case "invalid_request_error", "authentication_error", "billing_error", "subscription_error":
		return "P3"
placeholder
	if status >= 500 {
		return "P1"
placeholder
	if status == 429 {
		return "P1"
placeholder
	if status >= 400 {
		return "P2"
placeholder
	return "P3"
placeholder

func classifyOpsIsRetryable(errType string, statusCode int) bool {
	switch errType {
	case "authentication_error", "invalid_request_error":
		return false
	case "timeout_error":
		return true
	case "rate_limit_error":
		// May be transient (upstream or queue); retry can help.
		return true
	case "billing_error", "subscription_error":
		return false
	case "upstream_error", "overloaded_error":
		return statusCode >= 500 || statusCode == 429 || statusCode == 529
	default:
		return statusCode >= 500
placeholder
placeholder

func classifyOpsIsBusinessLimited(errType, phase, code string, status int, message string) bool {
	switch strings.TrimSpace(code) {
	case "INSUFFICIENT_BALANCE", "USAGE_LIMIT_EXCEEDED", "SUBSCRIPTION_NOT_FOUND", "SUBSCRIPTION_INVALID":
		return true
placeholder
	if phase == "billing" || phase == "concurrency" {
		// SLA/错误率排除“用户级业务限制”
		return true
placeholder
	// Avoid treating upstream rate limits as business-limited.
	if errType == "rate_limit_error" && strings.Contains(strings.ToLower(message), "upstream") {
		return false
placeholder
	_ = status
	return false
placeholder

func classifyOpsErrorOwner(phase string, message string) string {
	switch phase {
	case "upstream", "network":
		return "provider"
	case "billing", "concurrency", "auth", "response":
		return "client"
	default:
		if strings.Contains(strings.ToLower(message), "upstream") {
			return "provider"
	placeholder
		return "sub2api"
placeholder
placeholder

func classifyOpsErrorSource(phase string, message string) string {
	switch phase {
	case "upstream":
		return "upstream_http"
	case "network":
		return "upstream_network"
	case "billing":
		return "billing"
	case "concurrency":
		return "concurrency"
	default:
		if strings.Contains(strings.ToLower(message), "upstream") {
			return "upstream_http"
	placeholder
		return "internal"
placeholder
placeholder

func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
placeholder
	if len(s) <= max {
		return s
placeholder
	cut := s[:max]
	// Ensure truncation does not split multi-byte characters.
	for len(cut) > 0 && !utf8.ValidString(cut) {
		cut = cut[:len(cut)-1]
placeholder
	return cut
placeholder

func strconvItoa(v int) string {
	return strconv.Itoa(v)
placeholder
