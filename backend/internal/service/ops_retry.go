package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

const (
	OpsRetryModeClient   = "client"
	OpsRetryModeUpstream = "upstream"
)

const (
	opsRetryStatusRunning   = "running"
	opsRetryStatusSucceeded = "succeeded"
	opsRetryStatusFailed    = "failed"
)

const (
	opsRetryTimeout             = 60 * time.Second
	opsRetryCaptureBytesLimit   = 64 * 1024
	opsRetryResponsePreviewMax  = 8 * 1024
	opsRetryMinIntervalPerError = 10 * time.Second
	opsRetryMaxAccountSwitches  = 3
)

var opsRetryRequestHeaderAllowlist = map[string]bool{
	"anthropic-beta":    true,
	"anthropic-version": true,
placeholder

type opsRetryRequestType string

const (
	opsRetryTypeMessages  opsRetryRequestType = "messages"
	opsRetryTypeOpenAI    opsRetryRequestType = "openai_responses"
	opsRetryTypeGeminiV1B opsRetryRequestType = "gemini_v1beta"
)

type limitedResponseWriter struct {
	header      http.Header
	wroteHeader bool

	limit        int
	totalWritten int64
	buf          bytes.Buffer
placeholder

func newLimitedResponseWriter(limit int) *limitedResponseWriter {
	if limit <= 0 {
		limit = 1
placeholder
	return &limitedResponseWriter{
		header: make(http.Header),
		limit:  limit,
placeholder
placeholder

func (w *limitedResponseWriter) Header() http.Header {
	return w.header
placeholder

func (w *limitedResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
placeholder
	w.wroteHeader = true
placeholder

func (w *limitedResponseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
placeholder
	w.totalWritten += int64(len(p))

	if w.buf.Len() < w.limit {
		remaining := w.limit - w.buf.Len()
		if len(p) > remaining {
			_, _ = w.buf.Write(p[:remaining])
	placeholder else {
			_, _ = w.buf.Write(p)
	placeholder
placeholder

	// Pretend we wrote everything to avoid upstream/client code treating it as an error.
	return len(p), nil
placeholder

func (w *limitedResponseWriter) Flush() {placeholder

func (w *limitedResponseWriter) bodyBytes() []byte {
	return w.buf.Bytes()
placeholder

func (w *limitedResponseWriter) truncated() bool {
	return w.totalWritten > int64(w.limit)
placeholder

func (s *OpsService) RetryError(ctx context.Context, requestedByUserID int64, errorID int64, mode string, pinnedAccountID *int64) (*OpsRetryResult, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder

	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case OpsRetryModeClient, OpsRetryModeUpstream:
	default:
		return nil, infraerrors.BadRequest("OPS_RETRY_INVALID_MODE", "mode must be client or upstream")
placeholder

	latest, err := s.opsRepo.GetLatestRetryAttemptForError(ctx, errorID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, infraerrors.InternalServer("OPS_RETRY_LOAD_LATEST_FAILED", "Failed to check retry status").WithCause(err)
placeholder
	if latest != nil {
		if strings.EqualFold(latest.Status, opsRetryStatusRunning) || strings.EqualFold(latest.Status, "queued") {
			return nil, infraerrors.Conflict("OPS_RETRY_IN_PROGRESS", "A retry is already in progress for this error")
	placeholder

		lastAttemptAt := latest.CreatedAt
		if latest.FinishedAt != nil && !latest.FinishedAt.IsZero() {
			lastAttemptAt = *latest.FinishedAt
	placeholder else if latest.StartedAt != nil && !latest.StartedAt.IsZero() {
			lastAttemptAt = *latest.StartedAt
	placeholder

		if time.Since(lastAttemptAt) < opsRetryMinIntervalPerError {
			return nil, infraerrors.Conflict("OPS_RETRY_TOO_FREQUENT", "Please wait before retrying this error again")
	placeholder
placeholder

	errorLog, err := s.GetErrorLogByID(ctx, errorID)
	if err != nil {
		return nil, err
placeholder
	if strings.TrimSpace(errorLog.RequestBody) == "" {
		return nil, infraerrors.BadRequest("OPS_RETRY_NO_REQUEST_BODY", "No request body found to retry")
placeholder

	var pinned *int64
	if mode == OpsRetryModeUpstream {
		if pinnedAccountID != nil && *pinnedAccountID > 0 {
			pinned = pinnedAccountID
	placeholder else if errorLog.AccountID != nil && *errorLog.AccountID > 0 {
			pinned = errorLog.AccountID
	placeholder else {
			return nil, infraerrors.BadRequest("OPS_RETRY_PINNED_ACCOUNT_REQUIRED", "pinned_account_id is required for upstream retry")
	placeholder
placeholder

	startedAt := time.Now()
	attemptID, err := s.opsRepo.InsertRetryAttempt(ctx, &OpsInsertRetryAttemptInput{
		RequestedByUserID: requestedByUserID,
		SourceErrorID:     errorID,
		Mode:              mode,
		PinnedAccountID:   pinned,
		Status:            opsRetryStatusRunning,
		StartedAt:         startedAt,
placeholder)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
			return nil, infraerrors.Conflict("OPS_RETRY_IN_PROGRESS", "A retry is already in progress for this error")
	placeholder
		return nil, infraerrors.InternalServer("OPS_RETRY_CREATE_ATTEMPT_FAILED", "Failed to create retry attempt").WithCause(err)
placeholder

	result := &OpsRetryResult{
		AttemptID:         attemptID,
		Mode:              mode,
		Status:            opsRetryStatusFailed,
		PinnedAccountID:   pinned,
		HTTPStatusCode:    0,
		UpstreamRequestID: "",
		ResponsePreview:   "",
		ResponseTruncated: false,
		ErrorMessage:      "",
		StartedAt:         startedAt,
placeholder

	execCtx, cancel := context.WithTimeout(ctx, opsRetryTimeout)
	defer cancel()

	execRes := s.executeRetry(execCtx, errorLog, mode, pinned)

	finishedAt := time.Now()
	result.FinishedAt = finishedAt
	result.DurationMs = finishedAt.Sub(startedAt).Milliseconds()

	if execRes != nil {
		result.Status = execRes.status
		result.UsedAccountID = execRes.usedAccountID
		result.HTTPStatusCode = execRes.httpStatusCode
		result.UpstreamRequestID = execRes.upstreamRequestID
		result.ResponsePreview = execRes.responsePreview
		result.ResponseTruncated = execRes.responseTruncated
		result.ErrorMessage = execRes.errorMessage
placeholder

	updateCtx, updateCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer updateCancel()

	var updateErrMsg *string
	if strings.TrimSpace(result.ErrorMessage) != "" {
		msg := result.ErrorMessage
		updateErrMsg = &msg
placeholder
	var resultRequestID *string
	if strings.TrimSpace(result.UpstreamRequestID) != "" {
		v := result.UpstreamRequestID
		resultRequestID = &v
placeholder

	finalStatus := result.Status
	if strings.TrimSpace(finalStatus) == "" {
		finalStatus = opsRetryStatusFailed
placeholder

	if err := s.opsRepo.UpdateRetryAttempt(updateCtx, &OpsUpdateRetryAttemptInput{
		ID:              attemptID,
		Status:          finalStatus,
		FinishedAt:      finishedAt,
		DurationMs:      result.DurationMs,
		ResultRequestID: resultRequestID,
		ErrorMessage:    updateErrMsg,
placeholder); err != nil {
		// Best-effort: retry itself already executed; do not fail the API response.
		log.Printf("[Ops] UpdateRetryAttempt failed: %v", err)
placeholder

	return result, nil
placeholder

type opsRetryExecution struct {
	status string

	usedAccountID     *int64
	httpStatusCode    int
	upstreamRequestID string

	responsePreview   string
	responseTruncated bool

	errorMessage string
placeholder

func (s *OpsService) executeRetry(ctx context.Context, errorLog *OpsErrorLogDetail, mode string, pinnedAccountID *int64) *opsRetryExecution {
	if errorLog == nil {
		return &opsRetryExecution{
			status:       opsRetryStatusFailed,
			errorMessage: "missing error log",
	placeholder
placeholder

	reqType := detectOpsRetryType(errorLog.RequestPath)
	bodyBytes := []byte(errorLog.RequestBody)

	switch reqType {
	case opsRetryTypeMessages:
		bodyBytes = FilterThinkingBlocksForRetry(bodyBytes)
	case opsRetryTypeOpenAI, opsRetryTypeGeminiV1B:
		// No-op
placeholder

	switch strings.ToLower(strings.TrimSpace(mode)) {
	case OpsRetryModeUpstream:
		if pinnedAccountID == nil || *pinnedAccountID <= 0 {
			return &opsRetryExecution{
				status:       opsRetryStatusFailed,
				errorMessage: "pinned_account_id required for upstream retry",
		placeholder
	placeholder
		return s.executePinnedRetry(ctx, reqType, errorLog, bodyBytes, *pinnedAccountID)
	case OpsRetryModeClient:
		return s.executeClientRetry(ctx, reqType, errorLog, bodyBytes)
	default:
		return &opsRetryExecution{
			status:       opsRetryStatusFailed,
			errorMessage: "invalid retry mode",
	placeholder
placeholder
placeholder

func detectOpsRetryType(path string) opsRetryRequestType {
	p := strings.ToLower(strings.TrimSpace(path))
	switch {
	case strings.Contains(p, "/responses"):
		return opsRetryTypeOpenAI
	case strings.Contains(p, "/v1beta/"):
		return opsRetryTypeGeminiV1B
	default:
		return opsRetryTypeMessages
placeholder
placeholder

func (s *OpsService) executePinnedRetry(ctx context.Context, reqType opsRetryRequestType, errorLog *OpsErrorLogDetail, body []byte, pinnedAccountID int64) *opsRetryExecution {
	if s.accountRepo == nil {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "account repository not available"placeholder
placeholder

	account, err := s.accountRepo.GetByID(ctx, pinnedAccountID)
	if err != nil {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: fmt.Sprintf("account not found: %v", err)placeholder
placeholder
	if account == nil {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "account not found"placeholder
placeholder
	if !account.IsSchedulable() {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "account is not schedulable"placeholder
placeholder
	if errorLog.GroupID != nil && *errorLog.GroupID > 0 {
		if !containsInt64(account.GroupIDs, *errorLog.GroupID) {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "pinned account is not in the same group as the original request"placeholder
	placeholder
placeholder

	var release func()
	if s.concurrencyService != nil {
		acq, err := s.concurrencyService.AcquireAccountSlot(ctx, account.ID, account.Concurrency)
		if err != nil {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: fmt.Sprintf("acquire account slot failed: %v", err)placeholder
	placeholder
		if acq == nil || !acq.Acquired {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "account concurrency limit reached"placeholder
	placeholder
		release = acq.ReleaseFunc
placeholder
	if release != nil {
		defer release()
placeholder

	usedID := account.ID
	exec := s.executeWithAccount(ctx, reqType, errorLog, body, account)
	exec.usedAccountID = &usedID
	if exec.status == "" {
		exec.status = opsRetryStatusFailed
placeholder
	return exec
placeholder

func (s *OpsService) executeClientRetry(ctx context.Context, reqType opsRetryRequestType, errorLog *OpsErrorLogDetail, body []byte) *opsRetryExecution {
	groupID := errorLog.GroupID
	if groupID == nil || *groupID <= 0 {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "group_id missing; cannot reselect account"placeholder
placeholder

	model, stream, parsedErr := extractRetryModelAndStream(reqType, errorLog, body)
	if parsedErr != nil {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: parsedErr.Error()placeholder
placeholder
	_ = stream

	excluded := make(map[int64]struct{placeholder)
	switches := 0

	for {
		if switches >= opsRetryMaxAccountSwitches {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "retry failed after exhausting account failovers"placeholder
	placeholder

		selection, selErr := s.selectAccountForRetry(ctx, reqType, groupID, model, excluded)
		if selErr != nil {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: selErr.Error()placeholder
	placeholder
		if selection == nil || selection.Account == nil {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "no available accounts"placeholder
	placeholder

		account := selection.Account
		if !selection.Acquired || selection.ReleaseFunc == nil {
			excluded[account.ID] = struct{placeholder{placeholder
			switches++
			continue
	placeholder

		exec := func() *opsRetryExecution {
			defer selection.ReleaseFunc()
			return s.executeWithAccount(ctx, reqType, errorLog, body, account)
	placeholder()

		if exec != nil {
			if exec.status == opsRetryStatusSucceeded {
				usedID := account.ID
				exec.usedAccountID = &usedID
				return exec
		placeholder
			// If the gateway services ask for failover, try another account.
			if s.isFailoverError(exec.errorMessage) {
				excluded[account.ID] = struct{placeholder{placeholder
				switches++
				continue
		placeholder
			usedID := account.ID
			exec.usedAccountID = &usedID
			return exec
	placeholder

		excluded[account.ID] = struct{placeholder{placeholder
		switches++
placeholder
placeholder

func (s *OpsService) selectAccountForRetry(ctx context.Context, reqType opsRetryRequestType, groupID *int64, model string, excludedIDs map[int64]struct{placeholder) (*AccountSelectionResult, error) {
	switch reqType {
	case opsRetryTypeOpenAI:
		if s.openAIGatewayService == nil {
			return nil, fmt.Errorf("openai gateway service not available")
	placeholder
		return s.openAIGatewayService.SelectAccountWithLoadAwareness(ctx, groupID, "", model, excludedIDs)
	case opsRetryTypeGeminiV1B, opsRetryTypeMessages:
		if s.gatewayService == nil {
			return nil, fmt.Errorf("gateway service not available")
	placeholder
		return s.gatewayService.SelectAccountWithLoadAwareness(ctx, groupID, "", model, excludedIDs)
	default:
		return nil, fmt.Errorf("unsupported retry type: %s", reqType)
placeholder
placeholder

func extractRetryModelAndStream(reqType opsRetryRequestType, errorLog *OpsErrorLogDetail, body []byte) (model string, stream bool, err error) {
	switch reqType {
	case opsRetryTypeMessages:
		parsed, parseErr := ParseGatewayRequest(body)
		if parseErr != nil {
			return "", false, fmt.Errorf("failed to parse messages request body: %w", parseErr)
	placeholder
		return parsed.Model, parsed.Stream, nil
	case opsRetryTypeOpenAI:
		var v struct {
			Model  string `json:"model"`
			Stream bool   `json:"stream"`
	placeholder
		if err := json.Unmarshal(body, &v); err != nil {
			return "", false, fmt.Errorf("failed to parse openai request body: %w", err)
	placeholder
		return strings.TrimSpace(v.Model), v.Stream, nil
	case opsRetryTypeGeminiV1B:
		if strings.TrimSpace(errorLog.Model) == "" {
			return "", false, fmt.Errorf("missing model for gemini v1beta retry")
	placeholder
		return strings.TrimSpace(errorLog.Model), errorLog.Stream, nil
	default:
		return "", false, fmt.Errorf("unsupported retry type: %s", reqType)
placeholder
placeholder

func (s *OpsService) executeWithAccount(ctx context.Context, reqType opsRetryRequestType, errorLog *OpsErrorLogDetail, body []byte, account *Account) *opsRetryExecution {
	if account == nil {
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "missing account"placeholder
placeholder

	c, w := newOpsRetryContext(ctx, errorLog)

	var err error
	switch reqType {
	case opsRetryTypeOpenAI:
		if s.openAIGatewayService == nil {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "openai gateway service not available"placeholder
	placeholder
		_, err = s.openAIGatewayService.Forward(ctx, c, account, body)
	case opsRetryTypeGeminiV1B:
		if s.geminiCompatService == nil || s.antigravityGatewayService == nil {
			return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "gemini services not available"placeholder
	placeholder
		modelName := strings.TrimSpace(errorLog.Model)
		action := "generateContent"
		if errorLog.Stream {
			action = "streamGenerateContent"
	placeholder
		if account.Platform == PlatformAntigravity {
			_, err = s.antigravityGatewayService.ForwardGemini(ctx, c, account, modelName, action, errorLog.Stream, body)
	placeholder else {
			_, err = s.geminiCompatService.ForwardNative(ctx, c, account, modelName, action, errorLog.Stream, body)
	placeholder
	case opsRetryTypeMessages:
		switch account.Platform {
		case PlatformAntigravity:
			if s.antigravityGatewayService == nil {
				return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "antigravity gateway service not available"placeholder
		placeholder
			_, err = s.antigravityGatewayService.Forward(ctx, c, account, body)
		case PlatformGemini:
			if s.geminiCompatService == nil {
				return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "gemini gateway service not available"placeholder
		placeholder
			_, err = s.geminiCompatService.Forward(ctx, c, account, body)
		default:
			if s.gatewayService == nil {
				return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "gateway service not available"placeholder
		placeholder
			parsedReq, parseErr := ParseGatewayRequest(body)
			if parseErr != nil {
				return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "failed to parse request body"placeholder
		placeholder
			_, err = s.gatewayService.Forward(ctx, c, account, parsedReq)
	placeholder
	default:
		return &opsRetryExecution{status: opsRetryStatusFailed, errorMessage: "unsupported retry type"placeholder
placeholder

	statusCode := http.StatusOK
	if c != nil && c.Writer != nil {
		statusCode = c.Writer.Status()
placeholder

	upstreamReqID := extractUpstreamRequestID(c)
	preview, truncated := extractResponsePreview(w)

	exec := &opsRetryExecution{
		status:            opsRetryStatusFailed,
		httpStatusCode:    statusCode,
		upstreamRequestID: upstreamReqID,
		responsePreview:   preview,
		responseTruncated: truncated,
		errorMessage:      "",
placeholder

	if err == nil && statusCode < 400 {
		exec.status = opsRetryStatusSucceeded
		return exec
placeholder

	if err != nil {
		exec.errorMessage = err.Error()
placeholder else {
		exec.errorMessage = fmt.Sprintf("upstream returned status %d", statusCode)
placeholder

	return exec
placeholder

func newOpsRetryContext(ctx context.Context, errorLog *OpsErrorLogDetail) (*gin.Context, *limitedResponseWriter) {
	w := newLimitedResponseWriter(opsRetryCaptureBytesLimit)
	c, _ := gin.CreateTestContext(w)

	path := "/"
	if errorLog != nil && strings.TrimSpace(errorLog.RequestPath) != "" {
		path = errorLog.RequestPath
placeholder

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost"+path, bytes.NewReader(nil))
	req.Header.Set("content-type", "application/json")
	if errorLog != nil && strings.TrimSpace(errorLog.UserAgent) != "" {
		req.Header.Set("user-agent", errorLog.UserAgent)
placeholder
	// Restore a minimal, whitelisted subset of request headers to improve retry fidelity
	// (e.g. anthropic-beta / anthropic-version). Never replay auth credentials.
	if errorLog != nil && strings.TrimSpace(errorLog.RequestHeaders) != "" {
		var stored map[string]string
		if err := json.Unmarshal([]byte(errorLog.RequestHeaders), &stored); err == nil {
			for k, v := range stored {
				key := strings.TrimSpace(k)
				if key == "" {
					continue
			placeholder
				if !opsRetryRequestHeaderAllowlist[strings.ToLower(key)] {
					continue
			placeholder
				val := strings.TrimSpace(v)
				if val == "" {
					continue
			placeholder
				req.Header.Set(key, val)
		placeholder
	placeholder
placeholder

	c.Request = req
	return c, w
placeholder

func extractUpstreamRequestID(c *gin.Context) string {
	if c == nil || c.Writer == nil {
		return ""
placeholder
	h := c.Writer.Header()
	if h == nil {
		return ""
placeholder
	for _, key := range []string{"x-request-id", "X-Request-Id", "X-Request-ID"placeholder {
		if v := strings.TrimSpace(h.Get(key)); v != "" {
			return v
	placeholder
placeholder
	return ""
placeholder

func extractResponsePreview(w *limitedResponseWriter) (preview string, truncated bool) {
	if w == nil {
		return "", false
placeholder
	b := bytes.TrimSpace(w.bodyBytes())
	if len(b) == 0 {
		return "", w.truncated()
placeholder
	if len(b) > opsRetryResponsePreviewMax {
		return string(b[:opsRetryResponsePreviewMax]), true
placeholder
	return string(b), w.truncated()
placeholder

func containsInt64(items []int64, needle int64) bool {
	for _, v := range items {
		if v == needle {
			return true
	placeholder
placeholder
	return false
placeholder

func (s *OpsService) isFailoverError(message string) bool {
	msg := strings.ToLower(strings.TrimSpace(message))
	if msg == "" {
		return false
placeholder
	return strings.Contains(msg, "upstream error:") && strings.Contains(msg, "failover")
placeholder
