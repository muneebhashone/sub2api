package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	coderws "github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	defaultLiveMaxSessionDuration = time.Hour
	liveLeaseRefreshInterval      = 20 * time.Second
	liveRedisOperationTimeout     = 3 * time.Second
	liveClosedRecordTTL           = 24 * time.Hour
	liveObserverPollInterval      = 250 * time.Millisecond
	liveUpstreamBodyLimit         = 2 << 20
)

var (
	chatGPTLiveCallsURL        = "https://chatgpt.com/backend-api/codex/realtime/calls?intent=quicksilver&architecture=avas"
	chatGPTLiveSidebandBaseURL = "wss://chatgpt.com/backend-api/codex"
)

type liveFrameConn interface {
	ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error)
	WriteFrame(ctx context.Context, msgType coderws.MessageType, payload []byte) error
	Close() error
placeholder

func liveSidebandReadError(err error) error {
	if coderws.CloseStatus(err) == coderws.StatusNormalClosure {
		return ErrLiveCallNotFound
placeholder
	return err
placeholder

func hashLiveCallID(callID string) string {
	sum := sha256.Sum256([]byte(callID))
	return hex.EncodeToString(sum[:])
placeholder

func liveGroupID(groupID *int64) int64 {
	if groupID == nil {
		return 0
placeholder
	return *groupID
placeholder

func liveOptionalID(value int64) *int64 {
	if value <= 0 {
		return nil
placeholder
	result := value
	return &result
placeholder

func (s *OpenAIGatewayService) liveStore() (LiveCallStore, error) {
	if s == nil || s.cache == nil {
		return nil, ErrLiveUnavailable
placeholder
	store, ok := s.cache.(LiveCallStore)
	if !ok {
		return nil, ErrLiveUnavailable
placeholder
	return store, nil
placeholder

func (s *OpenAIGatewayService) liveConcurrencyCache() (LiveConcurrencyCache, error) {
	if s == nil || s.concurrencyService == nil || s.concurrencyService.cache == nil {
		return nil, ErrLiveUnavailable
placeholder
	cache, ok := s.concurrencyService.cache.(LiveConcurrencyCache)
	if !ok {
		return nil, ErrLiveUnavailable
placeholder
	return cache, nil
placeholder

func (s *OpenAIGatewayService) liveMaxSessionDuration() time.Duration {
	if s != nil && s.cfg != nil && s.cfg.Gateway.Live.MaxSessionDurationSeconds > 0 {
		return time.Duration(s.cfg.Gateway.Live.MaxSessionDurationSeconds) * time.Second
placeholder
	return defaultLiveMaxSessionDuration
placeholder

func ValidateLiveCallRequest(request *LiveCallRequest) error {
	if request == nil || strings.TrimSpace(request.SDP) == "" {
		return errors.New("sdp is required")
placeholder
	if len(request.Session) == 0 || !json.Valid(request.Session) {
		return errors.New("session must be valid JSON")
placeholder
	var sessionObject map[string]json.RawMessage
	if err := json.Unmarshal(request.Session, &sessionObject); err != nil {
		return errors.New("session must be a JSON object")
placeholder
	if sessionObject == nil {
		return errors.New("session must be a JSON object")
placeholder
	return nil
placeholder

// CreateLiveCall 创建 Frameless 会话。调用方须在调用期间持有普通用户槽位；
// 调度器持有的普通账号槽位会被同一个 Live 租约原子接替。
func (s *OpenAIGatewayService) CreateLiveCall(
	ctx context.Context,
	request *LiveCallRequest,
	identity LiveCallIdentity,
	userMaxConcurrency int,
) (*LiveCallCreated, error) {
	if err := ValidateLiveCallRequest(request); err != nil {
		return nil, err
placeholder
	store, err := s.liveStore()
	if err != nil {
		return nil, err
placeholder
	liveCache, err := s.liveConcurrencyCache()
	if err != nil {
		return nil, err
placeholder

	excluded := make(map[int64]struct{placeholder)
	var lastErr error
	for attempt := 0; attempt <= 3; attempt++ {
		selection, _, selectErr := s.SelectAccountWithSchedulerForCapability(
			ctx,
			identity.GroupID,
			"",
			uuid.NewString(),
			"",
			excluded,
			OpenAIUpstreamTransportHTTPSSE,
			OpenAIEndpointCapabilityLive,
			false,
			false,
			false,
		)
		if selectErr != nil {
			if lastErr != nil {
				return nil, lastErr
		placeholder
			return nil, selectErr
	placeholder
		if selection == nil || selection.Account == nil || !selection.Acquired {
			if selection != nil && selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
			return nil, ErrLiveConcurrencyFull
	placeholder

		account := selection.Account
		leaseID := generateRequestID()
		acquired, acquireErr := liveCache.AcquireLiveLease(
			ctx,
			account.ID,
			account.Concurrency,
			identity.UserID,
			userMaxConcurrency,
			identity.APIKeyID,
			leaseID,
			true,
		)
		if acquireErr != nil || !acquired {
			selection.ReleaseFunc()
			if acquireErr != nil {
				return nil, acquireErr
		placeholder
			return nil, ErrLiveConcurrencyFull
	placeholder

		created, createErr := s.createUpstreamLiveCall(ctx, account, request)
		selection.ReleaseFunc()
		if createErr != nil {
			s.releaseLiveLease(account.ID, identity.UserID, identity.APIKeyID, leaseID)
			if !s.shouldFailoverLiveCreateError(createErr) {
				return nil, createErr
		placeholder
			excluded[account.ID] = struct{placeholder{placeholder
			lastErr = createErr
			continue
	placeholder

		now := time.Now()
		model := strings.TrimSpace(gjson.GetBytes(request.Session, "model").String())
		if model == "" {
			model = "gpt-live"
	placeholder
		record := &LiveCallRecord{
			CallID:          created.CallID,
			CallHash:        hashLiveCallID(created.CallID),
			AccountID:       account.ID,
			APIKeyID:        identity.APIKeyID,
			UserID:          identity.UserID,
			GroupID:         liveGroupID(identity.GroupID),
			SubscriptionID:  liveGroupID(identity.SubscriptionID),
			LeaseID:         leaseID,
			Model:           model,
			CreatedAt:       now,
			ExpiresAt:       now.Add(s.liveMaxSessionDuration()),
			Controller:      LiveControllerPending,
			UserAgent:       identity.UserAgent,
			IPAddress:       identity.IPAddress,
			InboundEndpoint: identity.InboundEndpoint,
	placeholder
		mappingTTL := s.liveMaxSessionDuration() + 5*time.Minute
		if saveErr := store.SaveLiveCall(ctx, record, mappingTTL); saveErr != nil {
			s.releaseLiveLease(account.ID, identity.UserID, identity.APIKeyID, leaseID)
			return nil, fmt.Errorf("save live call mapping: %w", saveErr)
	placeholder
		created.Account = account
		go s.observeLiveCall(record.CallHash)
		return created, nil
placeholder
	if lastErr != nil {
		return nil, lastErr
placeholder
	return nil, ErrLiveUnavailable
placeholder

func (s *OpenAIGatewayService) shouldFailoverLiveCreateError(err error) bool {
	var upstreamErr *UpstreamFailoverError
	if !errors.As(err, &upstreamErr) {
		// 凭证读取和网络传输错误都可能只影响当前账号或代理。
		return true
placeholder
	return s.shouldFailoverOpenAIUpstreamResponse(
		upstreamErr.StatusCode,
		"",
		upstreamErr.ResponseBody,
	)
placeholder

func (s *OpenAIGatewayService) createUpstreamLiveCall(
	ctx context.Context,
	account *Account,
	request *LiveCallRequest,
) (*LiveCallCreated, error) {
	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		logLiveCreateStageFailure(ctx, account.ID, "access_token", err)
		return nil, err
placeholder
	body, err := json.Marshal(struct {
		SDP     string          `json:"sdp"`
		Session json.RawMessage `json:"session"`
placeholder{
		SDP:     request.SDP,
		Session: request.Session,
placeholder)
	if err != nil {
		return nil, err
placeholder
	reqCtx := WithHTTPUpstreamRedirectsDisabled(WithHTTPUpstreamProfile(ctx, HTTPUpstreamProfileOpenAI))
	upstreamReq, err := http.NewRequestWithContext(reqCtx, http.MethodPost, chatGPTLiveCallsURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder
	authHeaders, err := s.buildOpenAIAuthenticationHeaders(ctx, account, token)
	if err != nil {
		logLiveCreateStageFailure(ctx, account.ID, "authentication_headers", err)
		return nil, err
placeholder
	for key, values := range authHeaders {
		for _, value := range values {
			upstreamReq.Header.Add(key, value)
	placeholder
placeholder
	upstreamReq.Host = "chatgpt.com"
	if err := resolveAndSetOpenAIChatGPTAccountHeaders(ctx, s.accountRepo, upstreamReq.Header, account); err != nil {
		logLiveCreateStageFailure(ctx, account.ID, "account_headers", err)
		return nil, err
placeholder
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Accept", "application/sdp")
	applyLiveUpstreamIdentityHeaders(upstreamReq.Header)

	resp, err := s.httpUpstream.Do(upstreamReq, resolveAccountProxyURL(account), account.ID, account.Concurrency)
	if err != nil {
		logLiveCreateStageFailure(ctx, account.ID, "upstream_transport", err)
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	responseBody, readErr := io.ReadAll(io.LimitReader(resp.Body, liveUpstreamBodyLimit+1))
	if readErr != nil {
		return nil, readErr
placeholder
	if len(responseBody) > liveUpstreamBodyLimit {
		return nil, errors.New("live upstream response is too large")
placeholder
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logLiveUpstreamFailure(ctx, account.ID, resp.StatusCode, resp.Header, responseBody)
		return nil, &UpstreamFailoverError{
			StatusCode:      resp.StatusCode,
			ResponseBody:    responseBody,
			ResponseHeaders: resp.Header.Clone(),
	placeholder
placeholder
	callID, err := liveCallIDFromLocation(resp.Header.Get("Location"))
	if err != nil {
		return nil, err
placeholder
	return &LiveCallCreated{
		SDP:      responseBody,
		CallID:   callID,
		Location: resp.Header.Get("Location"),
placeholder, nil
placeholder

func logLiveCreateStageFailure(ctx context.Context, accountID int64, stage string, err error) {
	logger.FromContext(ctx).Warn(
		"OpenAI Live 创建阶段失败",
		zap.Int64("account_id", accountID),
		zap.String("stage", stage),
		zap.String("error_type", fmt.Sprintf("%T", err)),
	)
placeholder

func logLiveUpstreamFailure(
	ctx context.Context,
	accountID int64,
	statusCode int,
	headers http.Header,
	body []byte,
) {
	errorType := strings.TrimSpace(gjson.GetBytes(body, "error.type").String())
	errorCode := strings.TrimSpace(gjson.GetBytes(body, "error.code").String())
	errorMessage := strings.TrimSpace(gjson.GetBytes(body, "error.message").String())
	if errorType == "" {
		errorType = strings.TrimSpace(gjson.GetBytes(body, "type").String())
placeholder
	if errorCode == "" {
		errorCode = strings.TrimSpace(gjson.GetBytes(body, "code").String())
placeholder
	if errorMessage == "" {
		errorMessage = strings.TrimSpace(gjson.GetBytes(body, "message").String())
placeholder
	if errorMessage == "" {
		errorMessage = strings.TrimSpace(gjson.GetBytes(body, "detail").String())
placeholder

	logger.FromContext(ctx).Warn(
		"OpenAI Live 上游拒绝请求",
		zap.Int64("account_id", accountID),
		zap.Int("upstream_status_code", statusCode),
		zap.String("upstream_error_type", truncateOpenAIWSLogValue(errorType, 120)),
		zap.String("upstream_error_code", truncateOpenAIWSLogValue(errorCode, 120)),
		zap.String("upstream_error_message", truncateOpenAIWSLogValue(errorMessage, 300)),
		zap.String("upstream_content_type", truncateOpenAIWSLogValue(headers.Get("Content-Type"), 120)),
		zap.String("upstream_server", truncateOpenAIWSLogValue(headers.Get("Server"), 120)),
		zap.String("upstream_cf_mitigated", truncateOpenAIWSLogValue(headers.Get("Cf-Mitigated"), 120)),
		zap.String("upstream_cf_ray", truncateOpenAIWSLogValue(headers.Get("Cf-Ray"), 120)),
		zap.String("upstream_request_id", truncateOpenAIWSLogValue(headers.Get("X-Request-Id"), 120)),
	)
placeholder

func liveCallIDFromLocation(location string) (string, error) {
	location = strings.TrimSpace(location)
	if location == "" {
		return "", errors.New("live upstream response has no Location")
placeholder
	parsed, err := url.Parse(location)
	if err != nil {
		return "", fmt.Errorf("parse live Location: %w", err)
placeholder
	callID := strings.TrimSpace(path.Base(strings.TrimSuffix(parsed.Path, "/")))
	if callID == "" || callID == "." || callID == "codex" {
		return "", errors.New("live upstream Location has no call id")
placeholder
	return callID, nil
placeholder

func applyLiveUpstreamIdentityHeaders(headers http.Header) {
	headers.Set("OpenAI-Alpha", "quicksilver=v2")
	ensureCodexIdentityHeaders(headers)
	enforceCodexIdentityHeaders(headers)
	if strings.TrimSpace(headers.Get("session-id")) == "" {
		headers.Set("session-id", uuid.NewString())
placeholder
	if strings.TrimSpace(headers.Get("thread-id")) == "" {
		headers.Set("thread-id", uuid.NewString())
placeholder
	// Realtime/Live 不使用 Responses 的实验头。
	headers.Del("OpenAI-Beta")
placeholder

func (s *OpenAIGatewayService) liveSidebandHeaders(ctx context.Context, account *Account) (http.Header, error) {
	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	headers, err := s.buildOpenAIAuthenticationHeaders(ctx, account, token)
	if err != nil {
		return nil, err
placeholder
	if err := resolveAndSetOpenAIChatGPTAccountHeaders(ctx, s.accountRepo, headers, account); err != nil {
		return nil, err
placeholder
	applyLiveUpstreamIdentityHeaders(headers)
	return headers, nil
placeholder

func (s *OpenAIGatewayService) dialLiveSideband(ctx context.Context, record *LiveCallRecord) (liveFrameConn, error) {
	account, err := s.accountRepo.GetByID(ctx, record.AccountID)
	if err != nil {
		return nil, err
placeholder
	if account == nil || !account.SupportsOpenAIEndpointCapability(OpenAIEndpointCapabilityLive) {
		return nil, ErrLiveUnavailable
placeholder
	headers, err := s.liveSidebandHeaders(ctx, account)
	if err != nil {
		return nil, err
placeholder
	target := strings.TrimRight(chatGPTLiveSidebandBaseURL, "/") + "/" + url.PathEscape(record.CallID)
	conn, status, _, err := s.getOpenAIWSPassthroughDialer().Dial(ctx, target, headers, resolveAccountProxyURL(account))
	if err != nil {
		return nil, fmt.Errorf("dial live sideband (status %d): %w", status, err)
placeholder
	raw, ok := conn.(liveFrameConn)
	if !ok {
		_ = conn.Close()
		return nil, errors.New("live sideband transport does not support raw frames")
placeholder
	return raw, nil
placeholder

func (s *OpenAIGatewayService) GetLiveCallForIdentity(
	ctx context.Context,
	callID string,
	identity LiveCallIdentity,
) (*LiveCallRecord, error) {
	store, err := s.liveStore()
	if err != nil {
		return nil, err
placeholder
	record, err := store.GetLiveCall(ctx, hashLiveCallID(callID))
	if err != nil {
		return nil, err
placeholder
	if record.CallID != callID ||
		record.APIKeyID != identity.APIKeyID ||
		record.UserID != identity.UserID ||
		record.GroupID != liveGroupID(identity.GroupID) {
		return nil, ErrLiveIdentityMismatch
placeholder
	if record.Controller == LiveControllerClosed {
		return nil, ErrLiveCallNotFound
placeholder
	return record, nil
placeholder

// ProxyLiveSideband 让认证后的客户端接管控制连接；媒体始终不经过这里。
func (s *OpenAIGatewayService) ProxyLiveSideband(
	ctx context.Context,
	record *LiveCallRecord,
	downstream *coderws.Conn,
) error {
	if record == nil || downstream == nil {
		return ErrLiveCallNotFound
placeholder
	store, err := s.liveStore()
	if err != nil {
		return err
placeholder
	owner := uuid.NewString()
	claimed, err := store.ClaimLiveController(ctx, record.CallHash, LiveControllerProxy, owner)
	if err != nil {
		return err
placeholder
	if !claimed {
		return ErrLiveControllerChanged
placeholder

	// observer 轮询到接管状态后会关闭旧控制连接；同一个 call 可重新加入。
	time.Sleep(liveObserverPollInterval)
	upstream, err := s.dialLiveSideband(ctx, record)
	if err != nil {
		_, _ = store.ReleaseLiveController(context.Background(), record.CallHash, owner)
		go s.observeLiveCall(record.CallHash)
		return err
placeholder
	defer upstream.Close()
	downstream.SetReadLimit(openAIWSMessageReadLimitBytes)

	proxyCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	errCh := make(chan error, 2)
	go func() {
		for {
			messageType, payload, readErr := downstream.Read(proxyCtx)
			if readErr != nil {
				errCh <- readErr
				return
		placeholder
			if writeErr := upstream.WriteFrame(proxyCtx, messageType, payload); writeErr != nil {
				errCh <- writeErr
				return
		placeholder
	placeholder
placeholder()
	go func() {
		for {
			messageType, payload, readErr := upstream.ReadFrame(proxyCtx)
			if readErr != nil {
				errCh <- liveSidebandReadError(readErr)
				return
		placeholder
			if writeErr := downstream.Write(proxyCtx, messageType, payload); writeErr != nil {
				errCh <- writeErr
				return
		placeholder
			if messageType == coderws.MessageText {
				eventType := strings.TrimSpace(gjson.GetBytes(payload, "type").String())
				if eventType == "session.closed" || eventType == "session.ended" {
					errCh <- ErrLiveCallNotFound
					return
			placeholder
		placeholder
	placeholder
placeholder()

	runErr := s.runLiveController(proxyCtx, record, upstream, errCh)
	cancel()
	_, _ = store.ReleaseLiveController(context.Background(), record.CallHash, owner)
	if errors.Is(runErr, ErrLiveCallNotFound) {
		s.finalizeLiveCall(record)
		return runErr
placeholder
	if !errors.Is(runErr, context.DeadlineExceeded) && time.Now().Before(record.ExpiresAt) {
		go s.observeLiveCall(record.CallHash)
		return runErr
placeholder
	s.finalizeLiveCall(record)
	return runErr
placeholder

func (s *OpenAIGatewayService) runLiveController(
	ctx context.Context,
	record *LiveCallRecord,
	upstream liveFrameConn,
	errCh <-chan error,
) error {
	refreshTicker := time.NewTicker(liveLeaseRefreshInterval)
	defer refreshTicker.Stop()
	maxTimer := time.NewTimer(time.Until(record.ExpiresAt))
	defer maxTimer.Stop()
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case err := <-errCh:
			return err
		case <-maxTimer.C:
			closeCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_ = upstream.WriteFrame(closeCtx, coderws.MessageText, []byte(`{"type":"session.close"placeholder`))
			cancel()
			return context.DeadlineExceeded
		case <-refreshTicker.C:
			if !s.refreshLiveLease(record) {
				return ErrLiveUnavailable
		placeholder
	placeholder
placeholder
placeholder

func (s *OpenAIGatewayService) observeLiveCall(callHash string) {
	store, err := s.liveStore()
	if err != nil {
		return
placeholder
	owner := uuid.NewString()
	claimed, err := store.ClaimLiveController(context.Background(), callHash, LiveControllerObserver, owner)
	if err != nil || !claimed {
		return
placeholder
	for {
		record, getErr := store.GetLiveCall(context.Background(), callHash)
		if getErr != nil || record.Controller != LiveControllerObserver {
			return
	placeholder
		if !time.Now().Before(record.ExpiresAt) {
			s.finalizeLiveCall(record)
			return
	placeholder
		upstream, dialErr := s.dialLiveSideband(context.Background(), record)
		if dialErr != nil {
			if !s.waitForLiveObserverRetry(record) {
				return
		placeholder
			continue
	placeholder
		runErr := s.runLiveObserverConnection(record, upstream)
		_ = upstream.Close()
		if errors.Is(runErr, ErrLiveControllerChanged) {
			return
	placeholder
		if errors.Is(runErr, context.DeadlineExceeded) || errors.Is(runErr, ErrLiveCallNotFound) {
			s.finalizeLiveCall(record)
			return
	placeholder
		if !s.waitForLiveObserverRetry(record) {
			return
	placeholder
placeholder
placeholder

func (s *OpenAIGatewayService) runLiveObserverConnection(record *LiveCallRecord, upstream liveFrameConn) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	frameCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	go func() {
		for {
			messageType, payload, err := upstream.ReadFrame(ctx)
			if err != nil {
				select {
				case errCh <- liveSidebandReadError(err):
				case <-ctx.Done():
			placeholder
				return
		placeholder
			if messageType == coderws.MessageText {
				select {
				case frameCh <- payload:
				case <-ctx.Done():
					return
			placeholder
		placeholder
	placeholder
placeholder()
	refreshTicker := time.NewTicker(liveLeaseRefreshInterval)
	defer refreshTicker.Stop()
	controllerTicker := time.NewTicker(liveObserverPollInterval)
	defer controllerTicker.Stop()
	maxTimer := time.NewTimer(time.Until(record.ExpiresAt))
	defer maxTimer.Stop()
	store, _ := s.liveStore()
	for {
		select {
		case payload := <-frameCh:
			eventType := strings.TrimSpace(gjson.GetBytes(payload, "type").String())
			if eventType == "session.closed" || eventType == "session.ended" {
				return ErrLiveCallNotFound
		placeholder
		case err := <-errCh:
			return err
		case <-controllerTicker.C:
			controller, err := store.GetLiveController(context.Background(), record.CallHash)
			if err != nil {
				return err
		placeholder
			if controller != LiveControllerObserver {
				return ErrLiveControllerChanged
		placeholder
		case <-refreshTicker.C:
			if !s.refreshLiveLease(record) {
				return ErrLiveUnavailable
		placeholder
		case <-maxTimer.C:
			closeCtx, closeCancel := context.WithTimeout(context.Background(), 2*time.Second)
			_ = upstream.WriteFrame(closeCtx, coderws.MessageText, []byte(`{"type":"session.close"placeholder`))
			closeCancel()
			return context.DeadlineExceeded
	placeholder
placeholder
placeholder

func (s *OpenAIGatewayService) waitForLiveObserverRetry(record *LiveCallRecord) bool {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
	store, err := s.liveStore()
	if err != nil {
		return false
placeholder
	controller, err := store.GetLiveController(context.Background(), record.CallHash)
	return err == nil && controller == LiveControllerObserver && time.Now().Before(record.ExpiresAt)
placeholder

func (s *OpenAIGatewayService) refreshLiveLease(record *LiveCallRecord) bool {
	cache, err := s.liveConcurrencyCache()
	if err != nil {
		return false
placeholder
	ctx, cancel := context.WithTimeout(context.Background(), liveRedisOperationTimeout)
	defer cancel()
	refreshed, err := cache.RefreshLiveLease(ctx, record.AccountID, record.UserID, record.APIKeyID, record.LeaseID)
	return err == nil && refreshed
placeholder

func (s *OpenAIGatewayService) releaseLiveLease(accountID, userID, apiKeyID int64, leaseID string) {
	cache, err := s.liveConcurrencyCache()
	if err != nil {
		return
placeholder
	ctx, cancel := context.WithTimeout(context.Background(), liveRedisOperationTimeout)
	defer cancel()
	_ = cache.ReleaseLiveLease(ctx, accountID, userID, apiKeyID, leaseID)
placeholder

func (s *OpenAIGatewayService) finalizeLiveCall(record *LiveCallRecord) {
	if record == nil {
		return
placeholder
	store, err := s.liveStore()
	if err != nil {
		return
placeholder
	ctx, cancel := context.WithTimeout(context.Background(), liveRedisOperationTimeout)
	first, err := store.MarkLiveCallClosed(ctx, record.CallHash, liveClosedRecordTTL)
	cancel()
	if err != nil || !first {
		return
placeholder
	s.releaseLiveLease(record.AccountID, record.UserID, record.APIKeyID, record.LeaseID)
	if s.usageLogRepo == nil {
		return
placeholder
	duration := int(time.Since(record.CreatedAt).Milliseconds())
	if duration < 0 {
		duration = 0
placeholder
	inboundEndpoint := record.InboundEndpoint
	upstreamEndpoint := "/backend-api/codex/realtime/calls"
	userAgent := record.UserAgent
	ipAddress := record.IPAddress
	billingType := int8(BillingTypeBalance)
	if record.SubscriptionID > 0 {
		billingType = BillingTypeSubscription
placeholder
	_, _ = s.usageLogRepo.Create(context.Background(), &UsageLog{
		UserID:           record.UserID,
		APIKeyID:         record.APIKeyID,
		AccountID:        record.AccountID,
		RequestID:        record.CallHash,
		Model:            record.Model,
		RequestedModel:   record.Model,
		GroupID:          liveOptionalID(record.GroupID),
		SubscriptionID:   liveOptionalID(record.SubscriptionID),
		RateMultiplier:   1,
		BillingType:      billingType,
		RequestType:      RequestTypeLive,
		DurationMs:       &duration,
		UserAgent:        &userAgent,
		IPAddress:        &ipAddress,
		InboundEndpoint:  &inboundEndpoint,
		UpstreamEndpoint: &upstreamEndpoint,
		CreatedAt:        record.CreatedAt,
placeholder)
placeholder
