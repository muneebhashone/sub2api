package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	coderws "github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

type liveTestFrame struct {
	messageType coderws.MessageType
	payload     []byte
	err         error
placeholder

type liveTestFrameConn struct {
	reads     chan liveTestFrame
	writes    chan liveTestFrame
	closed    chan struct{placeholder
	closeOnce sync.Once
placeholder

func newLiveTestFrameConn() *liveTestFrameConn {
	return &liveTestFrameConn{
		reads:  make(chan liveTestFrame, 8),
		writes: make(chan liveTestFrame, 8),
		closed: make(chan struct{placeholder),
placeholder
placeholder

func (c *liveTestFrameConn) ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error) {
	select {
	case frame := <-c.reads:
		return frame.messageType, frame.payload, frame.err
	case <-c.closed:
		return coderws.MessageText, nil, coderws.CloseError{Code: coderws.StatusNormalClosureplaceholder
	case <-ctx.Done():
		return coderws.MessageText, nil, context.Cause(ctx)
placeholder
placeholder

func (c *liveTestFrameConn) WriteFrame(ctx context.Context, messageType coderws.MessageType, payload []byte) error {
	frame := liveTestFrame{messageType: messageType, payload: append([]byte(nil), payload...)placeholder
	select {
	case c.writes <- frame:
		return nil
	case <-c.closed:
		return errors.New("connection closed")
	case <-ctx.Done():
		return context.Cause(ctx)
placeholder
placeholder

func (c *liveTestFrameConn) WriteJSON(ctx context.Context, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
placeholder
	return c.WriteFrame(ctx, coderws.MessageText, payload)
placeholder

func (c *liveTestFrameConn) ReadMessage(ctx context.Context) ([]byte, error) {
	_, payload, err := c.ReadFrame(ctx)
	return payload, err
placeholder

func (c *liveTestFrameConn) Ping(context.Context) error { return nil placeholder

func (c *liveTestFrameConn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) placeholder)
	return nil
placeholder

type liveTestDialer struct {
	conn    *liveTestFrameConn
	url     string
	headers http.Header
placeholder

func (d *liveTestDialer) Dial(
	_ context.Context,
	wsURL string,
	headers http.Header,
	_ string,
) (openAIWSClientConn, int, http.Header, error) {
	d.url = wsURL
	d.headers = headers.Clone()
	return d.conn, http.StatusSwitchingProtocols, nil, nil
placeholder

type liveTestAccountRepo struct {
	AccountRepository
	account *Account
placeholder

func (r *liveTestAccountRepo) GetByID(context.Context, int64) (*Account, error) {
	return r.account, nil
placeholder

type liveTestStore struct {
	GatewayCache
	mu     sync.Mutex
	record *LiveCallRecord
placeholder

func (s *liveTestStore) SaveLiveCall(_ context.Context, record *LiveCallRecord, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := *record
	s.record = &copy
	return nil
placeholder

func (s *liveTestStore) GetLiveCall(_ context.Context, callHash string) (*LiveCallRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.record == nil || s.record.CallHash != callHash {
		return nil, ErrLiveCallNotFound
placeholder
	copy := *s.record
	return &copy, nil
placeholder

func (s *liveTestStore) ClaimLiveController(_ context.Context, callHash, controller, owner string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.record == nil || s.record.CallHash != callHash || s.record.Controller == LiveControllerClosed {
		return false, nil
placeholder
	if controller == LiveControllerObserver && s.record.Controller != LiveControllerPending {
		return false, nil
placeholder
	if controller == LiveControllerProxy && s.record.Controller != LiveControllerPending && s.record.Controller != LiveControllerObserver {
		return false, nil
placeholder
	s.record.Controller = controller
	s.record.ControllerOwner = owner
	return true, nil
placeholder

func (s *liveTestStore) ReleaseLiveController(_ context.Context, callHash, owner string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.record == nil || s.record.CallHash != callHash || s.record.ControllerOwner != owner {
		return false, nil
placeholder
	s.record.Controller = LiveControllerPending
	s.record.ControllerOwner = ""
	return true, nil
placeholder

func (s *liveTestStore) GetLiveController(_ context.Context, callHash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.record == nil || s.record.CallHash != callHash {
		return "", ErrLiveCallNotFound
placeholder
	return s.record.Controller, nil
placeholder

func (s *liveTestStore) MarkLiveCallClosed(_ context.Context, callHash string, _ time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.record == nil || s.record.CallHash != callHash || s.record.Controller == LiveControllerClosed {
		return false, nil
placeholder
	s.record.Controller = LiveControllerClosed
	s.record.ControllerOwner = ""
	return true, nil
placeholder

type liveTestConcurrencyCache struct {
	ConcurrencyCache
	mu       sync.Mutex
	releases int
placeholder

func (c *liveTestConcurrencyCache) AcquireLiveLease(
	context.Context,
	int64,
	int,
	int64,
	int,
	int64,
	string,
	bool,
) (bool, error) {
	return true, nil
placeholder

func (c *liveTestConcurrencyCache) RefreshLiveLease(
	context.Context,
	int64,
	int64,
	int64,
	string,
) (bool, error) {
	return true, nil
placeholder

func (c *liveTestConcurrencyCache) ReleaseLiveLease(
	context.Context,
	int64,
	int64,
	int64,
	string,
) error {
	c.mu.Lock()
	c.releases++
	c.mu.Unlock()
	return nil
placeholder

type liveTestUsageRepo struct {
	UsageLogRepository
	mu   sync.Mutex
	logs []*UsageLog
placeholder

func (r *liveTestUsageRepo) Create(_ context.Context, log *UsageLog) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *log
	r.logs = append(r.logs, &copy)
	return true, nil
placeholder

func TestRunLiveControllerClosesExpiredSession(t *testing.T) {
	upstream := newLiveTestFrameConn()
	record := &LiveCallRecord{ExpiresAt: time.Now().Add(20 * time.Millisecond)placeholder
	service := &OpenAIGatewayService{placeholder

	err := service.runLiveController(context.Background(), record, upstream, make(chan error))
	require.ErrorIs(t, err, context.DeadlineExceeded)

	select {
	case frame := <-upstream.writes:
		require.Equal(t, coderws.MessageText, frame.messageType)
		require.JSONEq(t, `{"type":"session.close"placeholder`, string(frame.payload))
	case <-time.After(time.Second):
		t.Fatal("没有向上游发送 session.close")
placeholder
placeholder

func TestFinalizeLiveCallIsIdempotentAndWritesZeroUsage(t *testing.T) {
	record := &LiveCallRecord{
		CallID:          "call_secret",
		CallHash:        hashLiveCallID("call_secret"),
		AccountID:       11,
		APIKeyID:        22,
		UserID:          33,
		GroupID:         44,
		LeaseID:         "lease-1",
		Model:           "gpt-live-test",
		CreatedAt:       time.Now().Add(-time.Second),
		ExpiresAt:       time.Now().Add(time.Hour),
		Controller:      LiveControllerPending,
		InboundEndpoint: "/v1/live",
placeholder
	store := &liveTestStore{placeholder
	require.NoError(t, store.SaveLiveCall(context.Background(), record, time.Hour))
	concurrencyCache := &liveTestConcurrencyCache{placeholder
	usageRepo := &liveTestUsageRepo{placeholder
	service := &OpenAIGatewayService{
		cache:              store,
		concurrencyService: NewConcurrencyService(concurrencyCache),
		usageLogRepo:       usageRepo,
placeholder

	service.finalizeLiveCall(record)
	service.finalizeLiveCall(record)

	concurrencyCache.mu.Lock()
	require.Equal(t, 1, concurrencyCache.releases)
	concurrencyCache.mu.Unlock()
	usageRepo.mu.Lock()
	require.Len(t, usageRepo.logs, 1)
	log := usageRepo.logs[0]
	usageRepo.mu.Unlock()
	require.Equal(t, RequestTypeLive, log.RequestType)
	require.Equal(t, record.CallHash, log.RequestID)
	require.NotEqual(t, record.CallID, log.RequestID)
	require.NotNil(t, log.DurationMs)
	require.Zero(t, log.InputTokens)
	require.Zero(t, log.OutputTokens)
	require.Zero(t, log.TotalCost)
	require.Zero(t, log.ActualCost)
placeholder

func TestGetLiveCallForIdentityRejectsMismatchedCaller(t *testing.T) {
	groupID := int64(44)
	record := &LiveCallRecord{
		CallID:     "call_identity",
		CallHash:   hashLiveCallID("call_identity"),
		APIKeyID:   22,
		UserID:     33,
		GroupID:    groupID,
		Controller: LiveControllerPending,
placeholder
	store := &liveTestStore{placeholder
	require.NoError(t, store.SaveLiveCall(context.Background(), record, time.Hour))
	service := &OpenAIGatewayService{cache: storeplaceholder

	_, err := service.GetLiveCallForIdentity(context.Background(), record.CallID, LiveCallIdentity{
		APIKeyID: 99,
		UserID:   record.UserID,
		GroupID:  &groupID,
placeholder)
	require.ErrorIs(t, err, ErrLiveIdentityMismatch)

	loaded, err := service.GetLiveCallForIdentity(context.Background(), record.CallID, LiveCallIdentity{
		APIKeyID: record.APIKeyID,
		UserID:   record.UserID,
		GroupID:  &groupID,
placeholder)
placeholder
	require.Equal(t, record.AccountID, loaded.AccountID)
placeholder

func TestProxyLiveSidebandForwardsTextAndBinary(t *testing.T) {
	account := &Account{
		ID:          11,
		Platform:    PlatformOpenAI,
		Type:        AccountTypeOAuth,
		Concurrency: 2,
placeholder
			"access_token":       "test-access-token",
			"chatgpt_account_id": "acct_test",
	placeholder,
placeholder
	record := &LiveCallRecord{
		CallID:     "call_proxy",
		CallHash:   hashLiveCallID("call_proxy"),
		AccountID:  account.ID,
		APIKeyID:   22,
		UserID:     33,
		LeaseID:    "lease-1",
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(time.Minute),
		Controller: LiveControllerPending,
placeholder
	store := &liveTestStore{placeholder
	require.NoError(t, store.SaveLiveCall(context.Background(), record, time.Hour))
	upstream := newLiveTestFrameConn()
	dialer := &liveTestDialer{conn: upstreamplaceholder
	service := &OpenAIGatewayService{
		accountRepo:               &liveTestAccountRepo{account: accountplaceholder,
		cache:                     store,
		openaiWSPassthroughDialer: dialer,
placeholder
	proxyResult := make(chan error, 1)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		downstream, err := coderws.Accept(writer, request, nil)
		if err != nil {
			proxyResult <- err
			return
	placeholder
		defer downstream.CloseNow()
		proxyResult <- service.ProxyLiveSideband(request.Context(), record, downstream)
placeholder))
	defer server.Close()

	client, _, err := coderws.Dial(
		context.Background(),
		"ws"+strings.TrimPrefix(server.URL, "http"),
		nil,
	)
placeholder
	defer client.CloseNow()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	require.NoError(t, client.Write(ctx, coderws.MessageText, []byte(`{"type":"client.text"placeholder`)))
	clientText := <-upstream.writes
	require.Equal(t, coderws.MessageText, clientText.messageType)
	require.JSONEq(t, `{"type":"client.text"placeholder`, string(clientText.payload))

	require.NoError(t, client.Write(ctx, coderws.MessageBinary, []byte{1, 2, 3placeholder))
	clientBinary := <-upstream.writes
	require.Equal(t, coderws.MessageBinary, clientBinary.messageType)
	require.Equal(t, []byte{1, 2, 3placeholder, clientBinary.payload)

	upstream.reads <- liveTestFrame{messageType: coderws.MessageText, payload: []byte(`{"type":"server.text"placeholder`)placeholder
	messageType, payload, err := client.Read(ctx)
placeholder
	require.Equal(t, coderws.MessageText, messageType)
	require.JSONEq(t, `{"type":"server.text"placeholder`, string(payload))

	upstream.reads <- liveTestFrame{messageType: coderws.MessageBinary, payload: []byte{4, 5, 6placeholderplaceholder
	messageType, payload, err = client.Read(ctx)
placeholder
	require.Equal(t, coderws.MessageBinary, messageType)
	require.Equal(t, []byte{4, 5, 6placeholder, payload)

	require.Equal(t, "wss://chatgpt.com/backend-api/codex/call_proxy", dialer.url)
	require.Equal(t, "Bearer test-access-token", dialer.headers.Get("Authorization"))
	require.Equal(t, "acct_test", dialer.headers.Get("Chatgpt-Account-Id"))
	upstream.reads <- liveTestFrame{err: coderws.CloseError{Code: coderws.StatusNormalClosureplaceholderplaceholder
	require.ErrorIs(t, <-proxyResult, ErrLiveCallNotFound)
placeholder
