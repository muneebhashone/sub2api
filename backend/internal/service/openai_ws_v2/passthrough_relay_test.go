package openai_ws_v2

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	coderws "github.com/coder/websocket"
	"github.com/stretchr/testify/require"
)

type passthroughTestFrame struct {
	msgType coderws.MessageType
	payload []byte
placeholder

type passthroughTestFrameConn struct {
	mu     sync.Mutex
	writes []passthroughTestFrame
	readCh chan passthroughTestFrame
	once   sync.Once
placeholder

type delayedReadFrameConn struct {
	base       FrameConn
	firstDelay time.Duration
	once       sync.Once
placeholder

type readStartSpyFrameConn struct {
	base      FrameConn
	started   chan struct{placeholder
	startOnce sync.Once
placeholder

type closeSpyFrameConn struct {
	closeCalls atomic.Int32
placeholder

func newPassthroughTestFrameConn(frames []passthroughTestFrame, autoClose bool) *passthroughTestFrameConn {
	c := &passthroughTestFrameConn{
		readCh: make(chan passthroughTestFrame, len(frames)+1),
placeholder
	for _, frame := range frames {
		copied := passthroughTestFrame{msgType: frame.msgType, payload: append([]byte(nil), frame.payload...)placeholder
		c.readCh <- copied
placeholder
	if autoClose {
		close(c.readCh)
placeholder
	return c
placeholder

func (c *passthroughTestFrameConn) ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error) {
	if ctx == nil {
		ctx = context.Background()
placeholder
	select {
	case <-ctx.Done():
		return coderws.MessageText, nil, ctx.Err()
	case frame, ok := <-c.readCh:
		if !ok {
			return coderws.MessageText, nil, io.EOF
	placeholder
		return frame.msgType, append([]byte(nil), frame.payload...), nil
placeholder
placeholder

func (c *passthroughTestFrameConn) WriteFrame(ctx context.Context, msgType coderws.MessageType, payload []byte) error {
	if ctx == nil {
		ctx = context.Background()
placeholder
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writes = append(c.writes, passthroughTestFrame{msgType: msgType, payload: append([]byte(nil), payload...)placeholder)
	return nil
placeholder

func (c *passthroughTestFrameConn) Close() error {
	c.once.Do(func() {
		defer func() { _ = recover() placeholder()
		close(c.readCh)
placeholder)
	return nil
placeholder

func (c *passthroughTestFrameConn) Writes() []passthroughTestFrame {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]passthroughTestFrame, len(c.writes))
	copy(out, c.writes)
	return out
placeholder

func (c *delayedReadFrameConn) ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error) {
	if c == nil || c.base == nil {
		return coderws.MessageText, nil, io.EOF
placeholder
	c.once.Do(func() {
		if c.firstDelay > 0 {
			timer := time.NewTimer(c.firstDelay)
			defer timer.Stop()
			select {
			case <-ctx.Done():
			case <-timer.C:
		placeholder
	placeholder
placeholder)
	return c.base.ReadFrame(ctx)
placeholder

func (c *delayedReadFrameConn) WriteFrame(ctx context.Context, msgType coderws.MessageType, payload []byte) error {
	if c == nil || c.base == nil {
		return io.EOF
placeholder
	return c.base.WriteFrame(ctx, msgType, payload)
placeholder

func (c *delayedReadFrameConn) Close() error {
	if c == nil || c.base == nil {
		return nil
placeholder
	return c.base.Close()
placeholder

func (c *readStartSpyFrameConn) ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error) {
	c.startOnce.Do(func() { close(c.started) placeholder)
	return c.base.ReadFrame(ctx)
placeholder

func (c *readStartSpyFrameConn) WriteFrame(ctx context.Context, msgType coderws.MessageType, payload []byte) error {
	return c.base.WriteFrame(ctx, msgType, payload)
placeholder

func (c *readStartSpyFrameConn) Close() error {
	return c.base.Close()
placeholder

func (c *closeSpyFrameConn) ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error) {
	if ctx == nil {
		ctx = context.Background()
placeholder
	<-ctx.Done()
	return coderws.MessageText, nil, ctx.Err()
placeholder

func (c *closeSpyFrameConn) WriteFrame(ctx context.Context, _ coderws.MessageType, _ []byte) error {
	if ctx == nil {
		ctx = context.Background()
placeholder
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
placeholder
placeholder

func (c *closeSpyFrameConn) Close() error {
	if c != nil {
		c.closeCalls.Add(1)
placeholder
	return nil
placeholder

func (c *closeSpyFrameConn) CloseCalls() int32 {
	if c == nil {
		return 0
placeholder
	return c.closeCalls.Load()
placeholder

func TestRelay_BasicRelayAndUsage(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_123","usage":{"input_tokens":7,"output_tokens":3,"input_tokens_details":{"cached_tokens":2placeholderplaceholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[{"type":"input_text","text":"hello"placeholder]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)
	require.Equal(t, "gpt-5.3-codex", result.RequestModel)
	require.Equal(t, "resp_123", result.RequestID)
	require.Equal(t, "response.completed", result.TerminalEventType)
	require.Equal(t, 7, result.Usage.InputTokens)
	require.Equal(t, 3, result.Usage.OutputTokens)
	require.Equal(t, 2, result.Usage.CacheReadInputTokens)
	require.Nil(t, result.FirstTokenMs)
	require.Equal(t, int64(1), result.ClientToUpstreamFrames)
	require.Equal(t, int64(1), result.UpstreamToClientFrames)
	require.Equal(t, int64(0), result.DroppedDownstreamFrames)

	upstreamWrites := upstreamConn.Writes()
	require.Len(t, upstreamWrites, 1)
	require.Equal(t, coderws.MessageText, upstreamWrites[0].msgType)
	require.JSONEq(t, string(firstPayload), string(upstreamWrites[0].payload))

	clientWrites := clientConn.Writes()
	require.Len(t, clientWrites, 1)
	require.Equal(t, coderws.MessageText, clientWrites[0].msgType)
	require.JSONEq(t, `{"type":"response.completed","response":{"id":"resp_123","usage":{"input_tokens":7,"output_tokens":3,"input_tokens_details":{"cached_tokens":2placeholderplaceholderplaceholderplaceholder`, string(clientWrites[0].payload))
placeholder

func TestRelay_FunctionCallOutputBytesPreserved(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_func","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[{"type":"function_call_output","call_id":"call_abc123","output":"{\"ok\":trueplaceholder"placeholder]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)

	upstreamWrites := upstreamConn.Writes()
	require.Len(t, upstreamWrites, 1)
	require.Equal(t, coderws.MessageText, upstreamWrites[0].msgType)
	require.Equal(t, firstPayload, upstreamWrites[0].payload)
placeholder

func TestRelay_UpstreamDisconnect(t *testing.T) {
	t.Parallel()

	// 上游立即关闭（EOF），客户端不发送额外帧
	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, true) // 立即 close -> EOF

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	// 上游 EOF 属于 disconnect，标记为 graceful
	require.Nil(t, relayExit, "上游 EOF 应被视为 graceful disconnect")
	require.Equal(t, "gpt-4o", result.RequestModel)
placeholder

func TestRelay_ClientDisconnect(t *testing.T) {
	t.Parallel()

	// 客户端立即关闭（EOF），上游阻塞读取直到 context 取消
	clientConn := newPassthroughTestFrameConn(nil, true) // 立即 close -> EOF
	upstreamConn := newPassthroughTestFrameConn(nil, false)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.NotNil(t, relayExit, "客户端 EOF 应返回可观测的中断状态")
	require.Equal(t, "client_disconnected", relayExit.Stage)
	require.Equal(t, "gpt-4o", result.RequestModel)
placeholder

func TestRelay_ClientDisconnect_DrainCapturesLateUsage(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, true)
	upstreamBase := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_drain","usage":{"input_tokens":6,"output_tokens":4,"input_tokens_details":{"cached_tokens":1placeholderplaceholderplaceholderplaceholder`),
	placeholder,
placeholder, true)
	upstreamConn := &delayedReadFrameConn{
		base:       upstreamBase,
		firstDelay: 80 * time.Millisecond,
placeholder

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		UpstreamDrainTimeout: 400 * time.Millisecond,
placeholder)
	require.NotNil(t, relayExit)
	require.Equal(t, "client_disconnected", relayExit.Stage)
	require.Equal(t, "resp_drain", result.RequestID)
	require.Equal(t, "response.completed", result.TerminalEventType)
	require.Equal(t, 6, result.Usage.InputTokens)
	require.Equal(t, 4, result.Usage.OutputTokens)
	require.Equal(t, 1, result.Usage.CacheReadInputTokens)
	require.Equal(t, int64(1), result.ClientToUpstreamFrames)
	require.Equal(t, int64(0), result.UpstreamToClientFrames)
	require.Equal(t, int64(1), result.DroppedDownstreamFrames)
placeholder

func TestRelay_IdleTimeout(t *testing.T) {
	t.Parallel()

	// 客户端和上游都不发送帧，idle timeout 应触发
	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, false)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用快进时间来加速 idle timeout
	now := time.Now()
	callCount := 0
	nowFn := func() time.Time {
		callCount++
		// 前几次调用返回正常时间（初始化阶段），之后快进
		if callCount <= 5 {
			return now
	placeholder
		return now.Add(time.Hour) // 快进到超时
placeholder

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		IdleTimeout: 2 * time.Second,
		Now:         nowFn,
placeholder)
	require.NotNil(t, relayExit, "应因 idle timeout 退出")
	require.Equal(t, "idle_timeout", relayExit.Stage)
	require.Equal(t, "gpt-4o", result.RequestModel)
placeholder

func TestRelay_IdleTimeoutDoesNotCloseClientOnError(t *testing.T) {
	t.Parallel()

	clientConn := &closeSpyFrameConn{placeholder
	upstreamConn := &closeSpyFrameConn{placeholder

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	callCount := 0
	nowFn := func() time.Time {
		callCount++
		if callCount <= 5 {
			return now
	placeholder
		return now.Add(time.Hour)
placeholder

	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		IdleTimeout: 2 * time.Second,
		Now:         nowFn,
placeholder)
	require.NotNil(t, relayExit, "应因 idle timeout 退出")
	require.Equal(t, "idle_timeout", relayExit.Stage)
	require.Zero(t, clientConn.CloseCalls(), "错误路径不应提前关闭客户端连接，交给上层决定 close code")
	require.GreaterOrEqual(t, upstreamConn.CloseCalls(), int32(1))
placeholder

func TestRelay_NilConnections(t *testing.T) {
	t.Parallel()

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx := context.Background()

	t.Run("nil client conn", func(t *testing.T) {
		upstreamConn := newPassthroughTestFrameConn(nil, true)
		_, relayExit := Relay(ctx, nil, upstreamConn, firstPayload, RelayOptions{placeholder)
		require.NotNil(t, relayExit)
		require.Equal(t, "relay_init", relayExit.Stage)
		require.Contains(t, relayExit.Err.Error(), "nil")
placeholder)

	t.Run("nil upstream conn", func(t *testing.T) {
		clientConn := newPassthroughTestFrameConn(nil, true)
		_, relayExit := Relay(ctx, clientConn, nil, firstPayload, RelayOptions{placeholder)
		require.NotNil(t, relayExit)
		require.Equal(t, "relay_init", relayExit.Stage)
		require.Contains(t, relayExit.Err.Error(), "nil")
placeholder)
placeholder

func TestRelay_MultipleUpstreamMessages(t *testing.T) {
	t.Parallel()

	// 上游发送多个事件（delta + completed），验证多帧中继和 usage 聚合
	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.output_text.delta","delta":"Hello"placeholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.output_text.delta","delta":" world"placeholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_multi","usage":{"input_tokens":10,"output_tokens":5,"input_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[{"type":"input_text","text":"hi"placeholder]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)
	require.Equal(t, "resp_multi", result.RequestID)
	require.Equal(t, "response.completed", result.TerminalEventType)
	require.Equal(t, 10, result.Usage.InputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
	require.Equal(t, 3, result.Usage.CacheReadInputTokens)
	require.NotNil(t, result.FirstTokenMs)

	// 验证所有 3 个上游帧都转发给了客户端
	clientWrites := clientConn.Writes()
	require.Len(t, clientWrites, 3)
placeholder

func TestRelay_OnTurnComplete_PerTerminalEvent(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_turn_1","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholderplaceholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.failed","response":{"id":"resp_turn_2","usage":{"input_tokens":3,"output_tokens":4placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	turns := make([]RelayTurnResult, 0, 2)
	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		OnTurnComplete: func(turn RelayTurnResult) {
			turns = append(turns, turn)
	placeholder,
placeholder)
	require.Nil(t, relayExit)
	require.Len(t, turns, 2)
	require.Equal(t, "resp_turn_1", turns[0].RequestID)
	require.Equal(t, "response.completed", turns[0].TerminalEventType)
	require.Equal(t, 2, turns[0].Usage.InputTokens)
	require.Equal(t, 1, turns[0].Usage.OutputTokens)
	require.Equal(t, "resp_turn_2", turns[1].RequestID)
	require.Equal(t, "response.failed", turns[1].TerminalEventType)
	require.Equal(t, 3, turns[1].Usage.InputTokens)
	require.Equal(t, 4, turns[1].Usage.OutputTokens)
	require.Equal(t, 5, result.Usage.InputTokens)
	require.Equal(t, 5, result.Usage.OutputTokens)
placeholder

func TestRelay_OnTurnComplete_UsesCurrentResponseCreateModel(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, false)
	firstPayload := []byte(`{"type":"response.create","model":"gpt-5.6-sol","input":[]placeholder`)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	turns := make(chan RelayTurnResult, 2)
	done := make(chan struct{placeholder)
	go func() {
		defer close(done)
		_, _ = Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
			OnTurnComplete: func(turn RelayTurnResult) {
				turns <- turn
		placeholder,
	placeholder)
placeholder()

	upstreamConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.completed","response":{"id":"resp_sol","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholderplaceholder`),
placeholder
	firstTurn := <-turns
	require.Equal(t, "gpt-5.6-sol", firstTurn.RequestModel)

	clientConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.create","model":"gpt-5.6-terra","input":[]placeholder`),
placeholder
	require.Eventually(t, func() bool {
		return len(upstreamConn.Writes()) == 2
placeholder, time.Second, 10*time.Millisecond)
	upstreamConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.completed","response":{"id":"resp_terra","usage":{"input_tokens":3,"output_tokens":1placeholderplaceholderplaceholder`),
placeholder
	secondTurn := <-turns
	require.Equal(t, "gpt-5.6-terra", secondTurn.RequestModel)

	_ = clientConn.Close()
	_ = upstreamConn.Close()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("relay did not stop after client close")
placeholder
placeholder

func TestRelay_OnTurnComplete_ProvidesTurnMetrics(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.output_text.delta","response_id":"resp_metric","delta":"hi"placeholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_metric","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	base := time.Unix(0, 0)
	var nowTick atomic.Int64
	nowFn := func() time.Time {
		step := nowTick.Add(1)
		return base.Add(time.Duration(step) * 5 * time.Millisecond)
placeholder

	var turn RelayTurnResult
	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		Now: nowFn,
		OnTurnComplete: func(current RelayTurnResult) {
			turn = current
	placeholder,
placeholder)
	require.Nil(t, relayExit)
	require.Equal(t, "resp_metric", turn.RequestID)
	require.Equal(t, "response.completed", turn.TerminalEventType)
	require.NotNil(t, turn.FirstTokenMs)
	require.GreaterOrEqual(t, *turn.FirstTokenMs, 0)
	require.Greater(t, turn.Duration.Milliseconds(), int64(0))
	require.NotNil(t, result.FirstTokenMs)
	require.Greater(t, result.Duration.Milliseconds(), int64(0))
placeholder

func TestRelay_OnTurnComplete_UsesResponseCreateTimeAcrossPricingBoundary(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_boundary","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	responseCreateAt := time.Date(2026, time.August, 17, 9, 59, 59, 0, time.UTC)
	upstreamResponseAt := responseCreateAt.Add(time.Second)
	var nowCalls atomic.Int64
	nowFn := func() time.Time {
		if nowCalls.Add(1) == 1 {
			return responseCreateAt
	placeholder
		return upstreamResponseAt
placeholder

	var turn RelayTurnResult
	_, relayExit := Relay(
		context.Background(),
		clientConn,
		upstreamConn,
		[]byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`),
		RelayOptions{
			Now:            nowFn,
			OnTurnComplete: func(current RelayTurnResult) { turn = current placeholder,
	placeholder,
	)

	require.Nil(t, relayExit)
	require.Equal(t, responseCreateAt, turn.StartedAt)
placeholder

func TestRelay_OnTurnComplete_UsesExplicitFirstTurnStartedAt(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_initial_boundary","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	responseCreateAt := time.Date(2026, time.August, 17, 9, 59, 59, 0, time.UTC)
	relayStartedAt := responseCreateAt.Add(time.Second)
	var turn RelayTurnResult
	_, relayExit := Relay(
		context.Background(),
		clientConn,
		upstreamConn,
		[]byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`),
		RelayOptions{
			FirstTurnStartedAt: responseCreateAt,
			Now:                func() time.Time { return relayStartedAt placeholder,
			OnTurnComplete:     func(current RelayTurnResult) { turn = current placeholder,
	placeholder,
	)

	require.Nil(t, relayExit)
	require.Equal(t, responseCreateAt, turn.StartedAt)
placeholder

func TestRelay_OnTurnComplete_UsesSubsequentResponseCreateTimeAcrossPricingBoundary(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, false)
	firstTurnAt := time.Date(2026, time.August, 17, 9, 0, 0, 0, time.UTC)
	secondTurnAt := time.Date(2026, time.August, 17, 9, 59, 59, 0, time.UTC)
	secondResponseAt := secondTurnAt.Add(time.Second)
	var clock atomic.Int64
	clock.Store(firstTurnAt.UnixNano())

	turns := make(chan RelayTurnResult, 2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan struct{placeholder)
	go func() {
		defer close(done)
		_, _ = Relay(
			ctx,
			clientConn,
			upstreamConn,
			[]byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`),
			RelayOptions{
				Now: func() time.Time { return time.Unix(0, clock.Load()).UTC() placeholder,
				OnTurnComplete: func(current RelayTurnResult) {
					turns <- current
			placeholder,
		placeholder,
		)
placeholder()

	require.Eventually(t, func() bool { return len(upstreamConn.Writes()) == 1 placeholder, time.Second, time.Millisecond)
	upstreamConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.completed","response":{"id":"resp_first","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
placeholder
	select {
	case <-turns:
	case <-time.After(time.Second):
		t.Fatal("first turn did not complete")
placeholder

	clock.Store(secondTurnAt.UnixNano())
	clientConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`),
placeholder
	require.Eventually(t, func() bool { return len(upstreamConn.Writes()) == 2 placeholder, time.Second, time.Millisecond)
	clock.Store(secondResponseAt.UnixNano())
	upstreamConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.completed","response":{"id":"resp_second","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
placeholder

	var secondTurn RelayTurnResult
	select {
	case secondTurn = <-turns:
	case <-time.After(time.Second):
		t.Fatal("second turn did not complete")
placeholder
	require.Equal(t, secondTurnAt, secondTurn.StartedAt)

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("relay did not stop after cancellation")
placeholder
placeholder

func TestRelay_BinaryFramePassthrough(t *testing.T) {
	t.Parallel()

	// 验证 binary frame 被透传但不进行 usage 解析
	binaryPayload := []byte{0x00, 0x01, 0x02, 0x03placeholder
	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageBinary,
			payload: binaryPayload,
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)
	// binary frame 不解析 usage
	require.Equal(t, 0, result.Usage.InputTokens)

	clientWrites := clientConn.Writes()
	require.Len(t, clientWrites, 1)
	require.Equal(t, coderws.MessageBinary, clientWrites[0].msgType)
	require.Equal(t, binaryPayload, clientWrites[0].payload)
placeholder

func TestRelay_BinaryJSONFrameSkipsObservation(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageBinary,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_binary","usage":{"input_tokens":7,"output_tokens":3placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)
	require.Equal(t, 0, result.Usage.InputTokens)
	require.Equal(t, "", result.RequestID)
	require.Equal(t, "", result.TerminalEventType)

	clientWrites := clientConn.Writes()
	require.Len(t, clientWrites, 1)
	require.Equal(t, coderws.MessageBinary, clientWrites[0].msgType)
placeholder

func TestRelay_UpstreamErrorEventPassthroughRaw(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	errorEvent := []byte(`{"type":"error","error":{"type":"invalid_request_error","message":"No tool call found"placeholderplaceholder`)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: errorEvent,
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)

	clientWrites := clientConn.Writes()
	require.Len(t, clientWrites, 1)
	require.Equal(t, coderws.MessageText, clientWrites[0].msgType)
	require.Equal(t, errorEvent, clientWrites[0].payload)
placeholder

func TestRelay_PreservesFirstMessageType(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		FirstMessageType: coderws.MessageBinary,
placeholder)
	require.Nil(t, relayExit)

	upstreamWrites := upstreamConn.Writes()
	require.Len(t, upstreamWrites, 1)
	require.Equal(t, coderws.MessageBinary, upstreamWrites[0].msgType)
	require.Equal(t, firstPayload, upstreamWrites[0].payload)
placeholder

func TestRelay_UsageParseFailureDoesNotBlockRelay(t *testing.T) {
	baseline := SnapshotMetrics().UsageParseFailureTotal

	// 上游发送无效 JSON（非 usage 格式），不应影响透传
	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_bad","usage":"not_an_object"placeholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	require.Nil(t, relayExit)
	// usage 解析失败，值为 0 但不影响透传
	require.Equal(t, 0, result.Usage.InputTokens)
	require.Equal(t, "response.completed", result.TerminalEventType)

	// 帧仍然被转发
	clientWrites := clientConn.Writes()
	require.Len(t, clientWrites, 1)
	require.GreaterOrEqual(t, SnapshotMetrics().UsageParseFailureTotal, baseline+1)
placeholder

func TestRelay_WriteUpstreamFirstMessageFails(t *testing.T) {
	t.Parallel()

	// 上游连接立即关闭，首包写入失败
	upstreamConn := newPassthroughTestFrameConn(nil, true)
	_ = upstreamConn.Close()

	// 覆盖 WriteFrame 使其返回错误
	errConn := &errorOnWriteFrameConn{placeholder
	clientConn := newPassthroughTestFrameConn(nil, false)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, relayExit := Relay(ctx, clientConn, errConn, firstPayload, RelayOptions{placeholder)
	require.NotNil(t, relayExit)
	require.Equal(t, "write_upstream", relayExit.Stage)
placeholder

func TestRelay_ContextCanceled(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, false)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)

	// 立即取消 context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{placeholder)
	// context 取消导致写首包失败
	require.NotNil(t, relayExit)
placeholder

func TestRelay_DownstreamPreambleStartsClientReader(t *testing.T) {
	clientBase := newPassthroughTestFrameConn(nil, false)
	clientConn := &readStartSpyFrameConn{base: clientBase, started: make(chan struct{placeholder)placeholder
	upstreamConn := newPassthroughTestFrameConn(nil, false)
	resultCh := make(chan *RelayExit, 1)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go func() {
		_, relayExit := Relay(
			ctx,
			clientConn,
			upstreamConn,
			[]byte(`{"type":"response.create","model":"gpt-5.1"placeholder`),
			RelayOptions{
				StartClientAfterFirstDownstream: true,
		placeholder,
		)
		resultCh <- relayExit
placeholder()

	upstreamConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.created","response":{"id":"resp_semantic_gate"placeholderplaceholder`),
placeholder
	require.Eventually(t, func() bool { return len(clientBase.Writes()) == 1 placeholder, time.Second, 10*time.Millisecond)
	select {
	case <-clientConn.started:
	case <-time.After(time.Second):
		t.Fatal("response.created did not start the client reader")
placeholder

	upstreamConn.readCh <- passthroughTestFrame{
		msgType: coderws.MessageText,
		payload: []byte(`{"type":"response.completed","response":{"id":"resp_semantic_gate","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
placeholder
	_ = upstreamConn.Close()
	select {
	case relayExit := <-resultCh:
		require.Nil(t, relayExit)
	case <-time.After(time.Second):
		t.Fatal("relay did not finish after terminal event and upstream close")
placeholder
placeholder

func TestRelay_TraceEvents_ContainsLifecycleStages(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_trace","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stages := make([]string, 0, 8)
	var stagesMu sync.Mutex
	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		OnTrace: func(event RelayTraceEvent) {
			stagesMu.Lock()
			stages = append(stages, event.Stage)
			stagesMu.Unlock()
	placeholder,
placeholder)
	require.Nil(t, relayExit)
	stagesMu.Lock()
	capturedStages := append([]string(nil), stages...)
	stagesMu.Unlock()
	require.Contains(t, capturedStages, "relay_start")
	require.Contains(t, capturedStages, "write_first_message_ok")
	require.Contains(t, capturedStages, "first_exit")
	require.Contains(t, capturedStages, "relay_complete")
placeholder

func TestRelay_TraceEvents_IdleTimeout(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn(nil, false)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-4o","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	callCount := 0
	nowFn := func() time.Time {
		callCount++
		if callCount <= 5 {
			return now
	placeholder
		return now.Add(time.Hour)
placeholder

	stages := make([]string, 0, 8)
	var stagesMu sync.Mutex
	_, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		IdleTimeout: 2 * time.Second,
		Now:         nowFn,
		OnTrace: func(event RelayTraceEvent) {
			stagesMu.Lock()
			stages = append(stages, event.Stage)
			stagesMu.Unlock()
	placeholder,
placeholder)
	require.NotNil(t, relayExit)
	require.Equal(t, "idle_timeout", relayExit.Stage)
	stagesMu.Lock()
	capturedStages := append([]string(nil), stages...)
	stagesMu.Unlock()
	require.Contains(t, capturedStages, "idle_timeout_triggered")
	require.Contains(t, capturedStages, "relay_exit")
placeholder

// errorOnWriteFrameConn 是一个写入总是失败的 FrameConn 实现，用于测试首包写入失败。
type errorOnWriteFrameConn struct{placeholder

func (c *errorOnWriteFrameConn) ReadFrame(ctx context.Context) (coderws.MessageType, []byte, error) {
	<-ctx.Done()
	return coderws.MessageText, nil, ctx.Err()
placeholder

func (c *errorOnWriteFrameConn) WriteFrame(_ context.Context, _ coderws.MessageType, _ []byte) error {
	return errors.New("write failed: connection refused")
placeholder

func (c *errorOnWriteFrameConn) Close() error {
	return nil
placeholder

func TestRelay_NoSemanticOutputTerminalSequence_FirstTokenMsNil(t *testing.T) {
	t.Parallel()

	for _, terminalEvent := range []string{"response.completed", "response.done"placeholder {
		terminalEvent := terminalEvent
		t.Run(terminalEvent, func(t *testing.T) {
			t.Parallel()

			clientConn := newPassthroughTestFrameConn(nil, false)
			upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
				{
					msgType: coderws.MessageText,
					payload: []byte(`{"type":"response.created","response":{"id":"resp_no_output"placeholderplaceholder`),
			placeholder,
				{
					msgType: coderws.MessageText,
					payload: []byte(`{"type":"response.in_progress","response":{"id":"resp_no_output"placeholderplaceholder`),
			placeholder,
				{
					msgType: coderws.MessageText,
					payload: []byte(`{"type":"response.content_part.done","response_id":"resp_no_output"placeholder`),
			placeholder,
				{
					msgType: coderws.MessageText,
					payload: []byte(`{"type":"response.output_item.done","response_id":"resp_no_output"placeholder`),
			placeholder,
				{
					msgType: coderws.MessageText,
					payload: []byte(`{"type":"` + terminalEvent + `","response":{"id":"resp_no_output","usage":{"input_tokens":2,"output_tokens":0placeholderplaceholderplaceholder`),
			placeholder,
		placeholder, true)

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			var turn RelayTurnResult
			result, relayExit := Relay(
				ctx,
				clientConn,
				upstreamConn,
				[]byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`),
				RelayOptions{OnTurnComplete: func(current RelayTurnResult) { turn = current placeholderplaceholder,
			)

			require.Nil(t, relayExit)
			require.Equal(t, terminalEvent, turn.TerminalEventType)
			require.Nil(t, turn.FirstTokenMs)
			require.Equal(t, terminalEvent, result.TerminalEventType)
			require.Nil(t, result.FirstTokenMs)
			require.Equal(t, int64(5), result.UpstreamToClientFrames)
	placeholder)
placeholder
placeholder

func TestRelay_NoDeltaOutputDoneEvent_RecordsFirstTokenBeforeTerminal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		donePayload string
placeholder{
		{
			name:        "output text done",
			donePayload: `{"type":"response.output_text.done","response_id":"resp_done","text":"hello"placeholder`,
	placeholder,
		{
			name:        "function call arguments done",
			donePayload: `{"type":"response.function_call_arguments.done","response_id":"resp_done","arguments":"{\"city\":\"Paris\"placeholder"placeholder`,
	placeholder,
placeholder
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			clientConn := newPassthroughTestFrameConn(nil, false)
			upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
				{msgType: coderws.MessageText, payload: []byte(`{"type":"response.created","response":{"id":"resp_done"placeholderplaceholder`)placeholder,
				{msgType: coderws.MessageText, payload: []byte(tt.donePayload)placeholder,
				{msgType: coderws.MessageText, payload: []byte(`{"type":"response.completed","response":{"id":"resp_done","usage":{"input_tokens":2,"output_tokens":1placeholderplaceholderplaceholder`)placeholder,
		placeholder, true)

			base := time.Unix(0, 0)
			var nowTick atomic.Int64
			nowFn := func() time.Time {
				return base.Add(time.Duration(nowTick.Add(1)) * 10 * time.Millisecond)
		placeholder
			var turn RelayTurnResult
			result, relayExit := Relay(
				context.Background(),
				clientConn,
				upstreamConn,
				[]byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`),
				RelayOptions{
					Now:            nowFn,
					OnTurnComplete: func(current RelayTurnResult) { turn = current placeholder,
			placeholder,
			)

			require.Nil(t, relayExit)
			require.NotNil(t, turn.FirstTokenMs)
			require.Less(t, int64(*turn.FirstTokenMs), turn.Duration.Milliseconds())
			require.NotNil(t, result.FirstTokenMs)
			require.Less(t, int64(*result.FirstTokenMs), result.Duration.Milliseconds())
	placeholder)
placeholder
placeholder

func TestRelay_OnTurnComplete_RealOpenAIStream_FirstTokenMs(t *testing.T) {
	t.Parallel()

	clientConn := newPassthroughTestFrameConn(nil, false)
	upstreamConn := newPassthroughTestFrameConn([]passthroughTestFrame{
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.created","response":{"id":"resp_real"placeholderplaceholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.output_text.delta","delta":"He"placeholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.output_text.delta","delta":"llo"placeholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.output_text.delta","delta":" world"placeholder`),
	placeholder,
		{
			msgType: coderws.MessageText,
			payload: []byte(`{"type":"response.completed","response":{"id":"resp_real","usage":{"input_tokens":2,"output_tokens":3placeholderplaceholderplaceholder`),
	placeholder,
placeholder, true)

	firstPayload := []byte(`{"type":"response.create","model":"gpt-5.3-codex","input":[]placeholder`)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	base := time.Unix(0, 0)
	var nowTick atomic.Int64
	nowFn := func() time.Time {
		step := nowTick.Add(1)
		return base.Add(time.Duration(step) * 10 * time.Millisecond)
placeholder

	var turn RelayTurnResult
	result, relayExit := Relay(ctx, clientConn, upstreamConn, firstPayload, RelayOptions{
		Now: nowFn,
		OnTurnComplete: func(current RelayTurnResult) {
			turn = current
	placeholder,
placeholder)
	require.Nil(t, relayExit)
	require.Equal(t, "resp_real", turn.RequestID)
	require.Equal(t, "response.completed", turn.TerminalEventType)

	require.NotNil(t, turn.FirstTokenMs, "per-turn FirstTokenMs must be captured for real OpenAI streams")
	require.Greater(t, turn.Duration.Milliseconds(), int64(0))

	require.Less(t,
		int64(*turn.FirstTokenMs),
		turn.Duration.Milliseconds(),
		"per-turn FirstTokenMs (%dms) should be strictly less than Duration (%dms); "+
			"equality indicates the bug where first_token is mistakenly stamped on the terminal event",
		*turn.FirstTokenMs, turn.Duration.Milliseconds(),
	)

	require.NotNil(t, result.FirstTokenMs)
	require.Greater(t, *result.FirstTokenMs, 0)
placeholder
