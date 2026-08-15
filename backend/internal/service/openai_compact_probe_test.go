package service

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestNormalizeAccountTestMode(t *testing.T) {
	tests := []struct {
		input string
		want  string
placeholder{
		{input: "", want: AccountTestModeDefaultplaceholder,
		{input: "default", want: AccountTestModeDefaultplaceholder,
		{input: " compact ", want: AccountTestModeCompactplaceholder,
		{input: "COMPACT", want: AccountTestModeCompactplaceholder,
		{input: "unknown", want: AccountTestModeDefaultplaceholder,
placeholder

	for _, tt := range tests {
		if got := normalizeAccountTestMode(tt.input); got != tt.want {
			t.Fatalf("normalizeAccountTestMode(%q) = %q, want %q", tt.input, got, tt.want)
	placeholder
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_SuccessMarksSupported(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	updates := buildOpenAICompactProbeExtraUpdates(&http.Response{StatusCode: http.StatusOKplaceholder, []byte(`{"id":"cmp_1"placeholder`), nil, true, now)

	if got := updates["openai_compact_supported"]; got != true {
		t.Fatalf("openai_compact_supported = %v, want true", got)
placeholder
	if got := updates["openai_compact_last_status"]; got != http.StatusOK {
		t.Fatalf("openai_compact_last_status = %v, want %d", got, http.StatusOK)
placeholder
	if got := updates["openai_compact_last_error"]; got != "" {
		t.Fatalf("openai_compact_last_error = %v, want empty string", got)
placeholder
	if got := updates["openai_compact_checked_at"]; got != now.Format(time.RFC3339) {
		t.Fatalf("openai_compact_checked_at = %v, want %s", got, now.Format(time.RFC3339))
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_404MarksUnsupported(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	body := []byte(`404 page not found`)
	updates := buildOpenAICompactProbeExtraUpdates(&http.Response{StatusCode: http.StatusNotFoundplaceholder, body, nil, false, now)

	if got := updates["openai_compact_supported"]; got != false {
		t.Fatalf("openai_compact_supported = %v, want false", got)
placeholder
	if got := updates["openai_compact_last_status"]; got != http.StatusNotFound {
		t.Fatalf("openai_compact_last_status = %v, want %d", got, http.StatusNotFound)
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_502DoesNotMarkUnsupported(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	updates := buildOpenAICompactProbeExtraUpdates(&http.Response{StatusCode: http.StatusBadGatewayplaceholder, []byte(`Upstream request failed`), nil, false, now)

	if _, exists := updates["openai_compact_supported"]; exists {
		t.Fatalf("did not expect openai_compact_supported for 502 response")
placeholder
	if got := updates["openai_compact_last_status"]; got != http.StatusBadGateway {
		t.Fatalf("openai_compact_last_status = %v, want %d", got, http.StatusBadGateway)
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_RequestErrorDoesNotMarkUnsupported(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	updates := buildOpenAICompactProbeExtraUpdates(nil, nil, errors.New("dial tcp timeout"), false, now)

	if _, exists := updates["openai_compact_supported"]; exists {
		t.Fatalf("did not expect openai_compact_supported for request error")
placeholder
	if got, exists := updates["openai_compact_last_status"]; !exists || got != nil {
		t.Fatalf("openai_compact_last_status = %v, want nil key", got)
placeholder
	if got := updates["openai_compact_last_error"]; got == "" {
		t.Fatalf("expected openai_compact_last_error to be populated")
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_NoResponseClearsLastStatus(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	updates := buildOpenAICompactProbeExtraUpdates(nil, nil, nil, false, now)

	if got, exists := updates["openai_compact_last_status"]; !exists || got != nil {
		t.Fatalf("openai_compact_last_status = %v, want nil key", got)
placeholder
	if got := updates["openai_compact_last_error"]; got != "compact probe failed" {
		t.Fatalf("openai_compact_last_error = %v, want compact probe failed", got)
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_UnknownModelDoesNotMarkUnsupported(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	body := []byte(`{"error":{"message":"unknown model gpt-5.4-openai-compact"placeholderplaceholder`)
	updates := buildOpenAICompactProbeExtraUpdates(&http.Response{StatusCode: http.StatusBadRequestplaceholder, body, nil, false, now)

	if _, exists := updates["openai_compact_supported"]; exists {
		t.Fatalf("did not expect openai_compact_supported for unknown-model diagnostics")
placeholder
	if got := updates["openai_compact_last_status"]; got != http.StatusBadRequest {
		t.Fatalf("openai_compact_last_status = %v, want %d", got, http.StatusBadRequest)
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_EmptyFailureBodyFallsBackToHTTPStatus(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	updates := buildOpenAICompactProbeExtraUpdates(&http.Response{StatusCode: http.StatusServiceUnavailableplaceholder, nil, nil, false, now)

	if got := updates["openai_compact_last_status"]; got != http.StatusServiceUnavailable {
		t.Fatalf("openai_compact_last_status = %v, want %d", got, http.StatusServiceUnavailable)
placeholder
	if got := updates["openai_compact_last_error"]; got != "HTTP 503" {
		t.Fatalf("openai_compact_last_error = %v, want HTTP 503", got)
placeholder
placeholder

func TestBuildOpenAICompactProbeExtraUpdates_2xxWithoutCompactionItemMarksUnsupported(t *testing.T) {
	now := time.Date(2026, 4, 10, 10, 0, 0, 0, time.UTC)
	updates := buildOpenAICompactProbeExtraUpdates(&http.Response{StatusCode: http.StatusOKplaceholder, []byte(`{"id":"resp_1","output":[]placeholder`), nil, false, now)

	if got := updates["openai_compact_supported"]; got != false {
		t.Fatalf("openai_compact_supported = %v, want false（2xx 无 compaction item = v2 不可用）", got)
placeholder
	if got := updates["openai_compact_last_error"]; got == "" {
		t.Fatalf("expected openai_compact_last_error to explain the missing compaction item")
placeholder
placeholder

func TestCreateOpenAICompactProbePayload_NativeV2Shape(t *testing.T) {
	payload := createOpenAICompactProbePayload("gpt-5.6-sol", true)
	if payload["stream"] != true {
		t.Fatalf("v2 probe payload must be streaming")
placeholder
	if payload["store"] != false {
		t.Fatalf("OAuth probe payload must carry store:false")
placeholder
	input, ok := payload["input"].([]any)
	if !ok || len(input) != 2 {
		t.Fatalf("expected 2 input items, got %v", payload["input"])
placeholder
	last, ok := input[len(input)-1].(map[string]any)
	if !ok || last["type"] != "compaction_trigger" {
		t.Fatalf("last input item must be compaction_trigger, got %v", input[len(input)-1])
placeholder

	apiKeyPayload := createOpenAICompactProbePayload("gpt-5.6-sol", false)
	if _, has := apiKeyPayload["store"]; has {
		t.Fatalf("API-key probe payload must not force store")
placeholder
placeholder

func TestOpenAICompactProbeFoundCompactionItem(t *testing.T) {
	sseWithItem := []byte("event: response.output_item.done\ndata: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"compaction\",\"id\":\"cmp_1\",\"encrypted_content\":\"blob\"placeholderplaceholder\n\ndata: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_1\",\"output\":[]placeholderplaceholder\n\n")
	if !openAICompactProbeFoundCompactionItem(sseWithItem) {
		t.Fatalf("SSE output_item.done 携带 compaction item 应判定为支持")
placeholder

	sseAlias := []byte("data: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"compaction_summary\",\"id\":\"cmp_2\"placeholderplaceholder\n\n")
	if !openAICompactProbeFoundCompactionItem(sseAlias) {
		t.Fatalf("compaction_summary 别名应判定为支持")
placeholder

	sseNoItem := []byte("data: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"message\",\"id\":\"msg_1\"placeholderplaceholder\n\ndata: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_2\",\"output\":[]placeholderplaceholder\n\n")
	if openAICompactProbeFoundCompactionItem(sseNoItem) {
		t.Fatalf("无 compaction item 的流不应判定为支持")
placeholder

	jsonWithItem := []byte(`{"id":"resp_3","output":[{"type":"compaction","id":"cmp_3"placeholder]placeholder`)
	if !openAICompactProbeFoundCompactionItem(jsonWithItem) {
		t.Fatalf("JSON output[] 携带 compaction item 应判定为支持（降级链兜底）")
placeholder

	if openAICompactProbeFoundCompactionItem(nil) {
		t.Fatalf("空响应不应判定为支持")
placeholder
placeholder

func TestOpenAICompactProbeFoundCompactionItem_TerminalResponseOutput(t *testing.T) {
	// 部分上游只在终态 response.completed 的 output[] 给出 compaction item，
	// 事件流里没有 output_item.done——探针同样必须判定为支持。
	sseTerminalOnly := []byte("data: {\"type\":\"response.output_item.done\",\"item\":{\"type\":\"message\",\"id\":\"msg_1\"placeholderplaceholder\n\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_t\",\"output\":[{\"type\":\"compaction\",\"id\":\"cmp_t\"placeholder]placeholderplaceholder\n\n")
	if !openAICompactProbeFoundCompactionItem(sseTerminalOnly) {
		t.Fatalf("终态 response.output 中的 compaction item 应判定为支持")
placeholder

	sseTerminalEmpty := []byte("data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_e\",\"output\":[]placeholderplaceholder\n\n")
	if openAICompactProbeFoundCompactionItem(sseTerminalEmpty) {
		t.Fatalf("终态 output 为空且无 item 事件时不应判定为支持")
placeholder
placeholder
