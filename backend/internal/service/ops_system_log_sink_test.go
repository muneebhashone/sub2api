package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

func TestOpsSystemLogSink_ShouldIndex(t *testing.T) {
	sink := &OpsSystemLogSink{placeholder

	cases := []struct {
		name  string
		event *logger.LogEvent
		want  bool
placeholder{
		{
			name:  "warn level",
			event: &logger.LogEvent{Level: "warn", Component: "app"placeholder,
			want:  true,
	placeholder,
		{
			name:  "error level",
			event: &logger.LogEvent{Level: "error", Component: "app"placeholder,
			want:  true,
	placeholder,
		{
			name:  "access component",
			event: &logger.LogEvent{Level: "info", Component: "http.access"placeholder,
			want:  true,
	placeholder,
		{
			name: "rejected access excluded from database sink",
			event: &logger.LogEvent{
				Level:     "info",
				Component: "http.access",
				Fields:    map[string]any{logger.OpsSystemLogSkipField: trueplaceholder,
		placeholder,
			want: false,
	placeholder,
		{
			name: "access component from fields (real zap path)",
			event: &logger.LogEvent{
				Level:     "info",
				Component: "",
				Fields:    map[string]any{"component": "http.access"placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name:  "audit component",
			event: &logger.LogEvent{Level: "info", Component: "audit.log_config_change"placeholder,
			want:  true,
	placeholder,
		{
			name: "audit component from fields (real zap path)",
			event: &logger.LogEvent{
				Level:     "info",
				Component: "",
				Fields:    map[string]any{"component": "audit.log_config_change"placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name:  "plain info",
			event: &logger.LogEvent{Level: "info", Component: "app"placeholder,
			want:  false,
	placeholder,
placeholder

	for _, tc := range cases {
		if got := sink.shouldIndex(tc.event); got != tc.want {
			t.Fatalf("%s: shouldIndex()=%v, want %v", tc.name, got, tc.want)
	placeholder
placeholder
placeholder

func TestOpsSystemLogSink_WriteLogEvent_ShouldDropWhenQueueFull(t *testing.T) {
	sink := &OpsSystemLogSink{
		queue: make(chan *logger.LogEvent, 1),
placeholder

	sink.WriteLogEvent(&logger.LogEvent{Level: "warn", Component: "app"placeholder)
	sink.WriteLogEvent(&logger.LogEvent{Level: "warn", Component: "app"placeholder)

	if got := len(sink.queue); got != 1 {
		t.Fatalf("queue len = %d, want 1", got)
placeholder
	if dropped := atomic.LoadUint64(&sink.droppedCount); dropped != 1 {
		t.Fatalf("droppedCount = %d, want 1", dropped)
placeholder
placeholder

func TestOpsSystemLogSink_Health(t *testing.T) {
	sink := &OpsSystemLogSink{
		queue: make(chan *logger.LogEvent, 10),
placeholder
	sink.lastError.Store("db timeout")
	atomic.StoreUint64(&sink.droppedCount, 3)
	atomic.StoreUint64(&sink.writeFailed, 2)
	atomic.StoreUint64(&sink.writtenCount, 5)
	atomic.StoreUint64(&sink.totalDelayNs, uint64(5000000)) // 5ms total -> avg 1ms
	sink.queue <- &logger.LogEvent{Level: "warn", Component: "app"placeholder
	sink.queue <- &logger.LogEvent{Level: "warn", Component: "app"placeholder

	health := sink.Health()
	if health.QueueDepth != 2 {
		t.Fatalf("queue depth = %d, want 2", health.QueueDepth)
placeholder
	if health.QueueCapacity != 10 {
		t.Fatalf("queue capacity = %d, want 10", health.QueueCapacity)
placeholder
	if health.DroppedCount != 3 {
		t.Fatalf("dropped = %d, want 3", health.DroppedCount)
placeholder
	if health.WriteFailed != 2 {
		t.Fatalf("write failed = %d, want 2", health.WriteFailed)
placeholder
	if health.WrittenCount != 5 {
		t.Fatalf("written = %d, want 5", health.WrittenCount)
placeholder
	if health.AvgWriteDelayMs != 1 {
		t.Fatalf("avg delay ms = %d, want 1", health.AvgWriteDelayMs)
placeholder
	if health.LastError != "db timeout" {
		t.Fatalf("last error = %q, want db timeout", health.LastError)
placeholder
placeholder

func TestOpsSystemLogSink_StartStopAndFlushSuccess(t *testing.T) {
	done := make(chan struct{placeholder, 1)
	var captured []*OpsInsertSystemLogInput
	repo := &opsRepoMock{
		BatchInsertSystemLogsFn: func(_ context.Context, inputs []*OpsInsertSystemLogInput) (int64, error) {
			captured = append(captured, inputs...)
			select {
			case done <- struct{placeholder{placeholder:
			default:
		placeholder
			return int64(len(inputs)), nil
	placeholder,
placeholder

	sink := NewOpsSystemLogSink(repo)
	sink.host = "api-node-1"
	sink.batchSize = 1
	sink.flushInterval = 10 * time.Millisecond
	sink.Start()
	defer sink.Stop()

	sink.WriteLogEvent(&logger.LogEvent{
		Time:      time.Now().UTC(),
		Level:     "warn",
		Component: "http.access",
		Message:   `authorization="Bearer sk-test-123"`,
		Fields: map[string]any{
			"component":         "http.access",
			"request_id":        "req-1",
			"client_request_id": "creq-1",
			"user_id":           "12",
			"api_key_id":        int64(56),
			"account_id":        json.Number("34"),
			"platform":          "openai",
			"model":             "gpt-5",
	placeholder,
placeholder)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for sink flush")
placeholder

	if len(captured) != 1 {
		t.Fatalf("captured len = %d, want 1", len(captured))
placeholder
	item := captured[0]
	if item.Host != "api-node-1" {
		t.Fatalf("host = %q, want api-node-1", item.Host)
placeholder
	if item.RequestID != "req-1" || item.ClientRequestID != "creq-1" {
		t.Fatalf("unexpected request ids: %+v", item)
placeholder
	if item.UserID == nil || *item.UserID != 12 {
		t.Fatalf("unexpected user_id: %+v", item.UserID)
placeholder
	if item.APIKeyID == nil || *item.APIKeyID != 56 {
		t.Fatalf("unexpected api_key_id: %+v", item.APIKeyID)
placeholder
	if item.AccountID == nil || *item.AccountID != 34 {
		t.Fatalf("unexpected account_id: %+v", item.AccountID)
placeholder
	if strings.TrimSpace(item.Message) == "" {
		t.Fatalf("message should not be empty")
placeholder
	// writtenCount is incremented after BatchInsertSystemLogsFn returns,
	// so poll briefly to avoid a race between the done signal and the atomic add.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if sink.Health().WrittenCount > 0 {
			break
	placeholder
		time.Sleep(time.Millisecond)
placeholder
	health := sink.Health()
	if health.WrittenCount == 0 {
		t.Fatalf("written_count should be >0")
placeholder
placeholder

func TestOpsSystemLogSink_FlushFailureUpdatesHealth(t *testing.T) {
	repo := &opsRepoMock{
		BatchInsertSystemLogsFn: func(_ context.Context, inputs []*OpsInsertSystemLogInput) (int64, error) {
			return 0, errors.New("db unavailable")
	placeholder,
placeholder
	sink := NewOpsSystemLogSink(repo)
	sink.batchSize = 1
	sink.flushInterval = 10 * time.Millisecond
	sink.Start()
	defer sink.Stop()

	sink.WriteLogEvent(&logger.LogEvent{
		Time:      time.Now().UTC(),
		Level:     "warn",
		Component: "app",
		Message:   "boom",
		Fields:    map[string]any{placeholder,
placeholder)

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		health := sink.Health()
		if health.WriteFailed > 0 {
			if !strings.Contains(health.LastError, "db unavailable") {
				t.Fatalf("unexpected last error: %s", health.LastError)
		placeholder
			return
	placeholder
		time.Sleep(20 * time.Millisecond)
placeholder
	t.Fatalf("write_failed_count not updated")
placeholder

func TestOpsSystemLogSink_StopFlushUsesActiveContextAndDrainsQueue(t *testing.T) {
	var inserted int64
	var canceledCtxCalls int64
	repo := &opsRepoMock{
		BatchInsertSystemLogsFn: func(ctx context.Context, inputs []*OpsInsertSystemLogInput) (int64, error) {
			if err := ctx.Err(); err != nil {
				atomic.AddInt64(&canceledCtxCalls, 1)
				return 0, err
		placeholder
			atomic.AddInt64(&inserted, int64(len(inputs)))
			return int64(len(inputs)), nil
	placeholder,
placeholder

	sink := NewOpsSystemLogSink(repo)
	sink.batchSize = 200
	sink.flushInterval = time.Hour
	sink.Start()

	sink.WriteLogEvent(&logger.LogEvent{
		Time:      time.Now().UTC(),
		Level:     "warn",
		Component: "app",
		Message:   "pending-on-shutdown",
		Fields:    map[string]any{"component": "http.access"placeholder,
placeholder)

	sink.Stop()

	if got := atomic.LoadInt64(&inserted); got != 1 {
		t.Fatalf("inserted = %d, want 1", got)
placeholder
	if got := atomic.LoadInt64(&canceledCtxCalls); got != 0 {
		t.Fatalf("canceled ctx calls = %d, want 0", got)
placeholder
	health := sink.Health()
	if health.WrittenCount != 1 {
		t.Fatalf("written_count = %d, want 1", health.WrittenCount)
placeholder
placeholder

type stringerValue string

func (s stringerValue) String() string { return string(s) placeholder

func TestOpsSystemLogSink_HelperFunctions(t *testing.T) {
	src := map[string]any{"a": 1placeholder
	cloned := copyMap(src)
	src["a"] = 2
	v, ok := cloned["a"].(int)
	if !ok || v != 1 {
		t.Fatalf("copyMap should create copy")
placeholder
	if got := asString(stringerValue(" hello ")); got != "hello" {
		t.Fatalf("asString stringer = %q", got)
placeholder
	if got := asString(fmt.Errorf("x")); got != "" {
		t.Fatalf("asString error should be empty, got %q", got)
placeholder
	if got := asString(123); got != "" {
		t.Fatalf("asString non-string should be empty, got %q", got)
placeholder

	cases := []struct {
		in   any
		want int64
		ok   bool
placeholder{
		{in: 5, want: 5, ok: trueplaceholder,
		{in: int64(6), want: 6, ok: trueplaceholder,
		{in: float64(7), want: 7, ok: trueplaceholder,
		{in: json.Number("8"), want: 8, ok: trueplaceholder,
		{in: "9", want: 9, ok: trueplaceholder,
		{in: "0", ok: falseplaceholder,
		{in: -1, ok: falseplaceholder,
		{in: "abc", ok: falseplaceholder,
placeholder
	for _, tc := range cases {
		got := asInt64Ptr(tc.in)
		if tc.ok {
			if got == nil || *got != tc.want {
				t.Fatalf("asInt64Ptr(%v) = %+v, want %d", tc.in, got, tc.want)
		placeholder
	placeholder else if got != nil {
			t.Fatalf("asInt64Ptr(%v) should be nil, got %d", tc.in, *got)
	placeholder
placeholder
placeholder

func TestNormalizeSystemLogHost(t *testing.T) {
	if got := normalizeSystemLogHost(" api-node-1 ", nil); got != "api-node-1" {
		t.Fatalf("trimmed host = %q, want api-node-1", got)
placeholder
	if got := normalizeSystemLogHost("", nil); got != "unknown" {
		t.Fatalf("empty host = %q, want unknown", got)
placeholder
	if got := normalizeSystemLogHost("api-node-1", errors.New("hostname unavailable")); got != "unknown" {
		t.Fatalf("errored host = %q, want unknown", got)
placeholder
	longHost := strings.Repeat("节", maxSystemLogHostLength+1)
	got := normalizeSystemLogHost(longHost, nil)
	if runeCount := len([]rune(got)); runeCount != maxSystemLogHostLength {
		t.Fatalf("truncated host rune count = %d, want %d", runeCount, maxSystemLogHostLength)
placeholder
placeholder
