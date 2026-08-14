package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

func newCompactBodySignalTestContext(t *testing.T, path string, body []byte) *gin.Context {
placeholder
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_RemoteV2StaysOnResponses(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	tests := []struct {
		name       string
		betaHeader string
		userAgent  string
placeholder{
		{name: "headerless"placeholder,
		{name: "unrelated_header", betaHeader: "responses_websockets_v2"placeholder,
		{name: "wrong_case_header", betaHeader: "REMOTE_COMPACTION_V2"placeholder,
		{name: "declared_header", betaHeader: "remote_compaction_v2"placeholder,
		{name: "codex_cli_user_agent", userAgent: "codex_cli_rs/0.144.1 (Ubuntu 22.4.0; x86_64) xterm-256color"placeholder,
		{name: "codex_desktop_user_agent", userAgent: "Codex Desktop/0.139.0 (Mac OS X 14; arm64) unknown"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(`{
				"model":"gpt-5.6-sol",
				"stream":true,
				"store":true,
				"prompt_cache_key":"pck-signal-1",
				"reasoning":{"effort":"max","context":"all_turns"placeholder,
				"input":[
					{"type":"message","role":"user","content":"hello"placeholder,
					{"type":"compaction_trigger"placeholder
				]
		placeholder`)
			c := newCompactBodySignalTestContext(t, "/v1/responses", body)
			if tt.betaHeader != "" {
				c.Request.Header.Set("x-codex-beta-features", tt.betaHeader)
		placeholder
			if tt.userAgent != "" {
				c.Request.Header.Set("User-Agent", tt.userAgent)
		placeholder

			normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
			require.True(t, ok)

			require.Equal(t, "/v1/responses", c.Request.URL.Path)
			require.False(t, isOpenAIRemoteCompactPath(c))
			require.Equal(t, body, normalized)
			require.True(t, gjson.GetBytes(normalized, "stream").Bool())
			require.True(t, gjson.GetBytes(normalized, "store").Bool())
			require.Equal(t, "pck-signal-1", gjson.GetBytes(normalized, "prompt_cache_key").String())
			require.Equal(t, "max", gjson.GetBytes(normalized, "reasoning.effort").String())
			require.Equal(t, "all_turns", gjson.GetBytes(normalized, "reasoning.context").String())
			require.True(t, requiresOpenAICompactAccount(c, normalized))

			reqStream, streamOK := parseOpenAICompatibleStream(normalized)
			require.True(t, streamOK)
			require.True(t, reqStream)

			_, seedExists := c.Get(service.OpenAICompactSessionSeedKeyForTest())
			require.False(t, seedExists)
			_, streamMarkerExists := c.Get(service.OpenAICompactClientStreamKeyForTest())
			require.False(t, streamMarkerExists)
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_RemoteV2PathAliasesStayOnResponses(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"compaction_trigger"placeholder]placeholder`)
	for _, path := range []string{"/v1/responses/", "/backend-api/codex/responses"placeholder {
		t.Run(path, func(t *testing.T) {
			c := newCompactBodySignalTestContext(t, path, body)

			normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
			require.True(t, ok)
			require.Equal(t, path, c.Request.URL.Path)
			require.Equal(t, body, normalized)
	placeholder)
placeholder
placeholder

func TestRequiresOpenAICompactAccount_RemoteV2UsesNativeRequestSignal(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	tests := []struct {
		name       string
		body       []byte
		betaHeader string
		wantBefore bool
		wantAfter  bool
		path       string
placeholder{
		{
			name:       "native_v2_headerless",
			body:       []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
			wantBefore: true,
			wantAfter:  true,
			path:       "/v1/responses",
	placeholder,
		{
			name:       "native_v2_without_trigger",
			body:       []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"message","role":"user","content":"hello"placeholder]placeholder`),
			betaHeader: "remote_compaction_v2",
			wantBefore: false,
			wantAfter:  false,
			path:       "/v1/responses",
	placeholder,
		{
			name:       "native_v2_unrelated_header",
			body:       []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
			betaHeader: "responses_websockets_v2",
			wantBefore: true,
			wantAfter:  true,
			path:       "/v1/responses",
	placeholder,
		{
			name:       "native_v2_wrong_case_header",
			body:       []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
			betaHeader: "REMOTE_COMPACTION_V2",
			wantBefore: true,
			wantAfter:  true,
			path:       "/v1/responses",
	placeholder,
		{
			name:       "responses_subpath_with_native_signal",
			body:       []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
			betaHeader: "remote_compaction_v2",
			wantBefore: false,
			wantAfter:  false,
			path:       "/v1/responses/resp_123/responses",
	placeholder,
		{
			name:       "stream_false_promotes",
			body:       []byte(`{"model":"gpt-5.6-sol","stream":false,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
			wantBefore: false,
			wantAfter:  true,
			path:       "/v1/responses",
	placeholder,
		{
			name:       "stream_absent_promotes",
			body:       []byte(`{"model":"gpt-5.6-sol","input":[{"type":"compaction_trigger"placeholder]placeholder`),
			betaHeader: "remote_compaction_v2",
			wantBefore: false,
			wantAfter:  true,
			path:       "/v1/responses",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCompactBodySignalTestContext(t, tt.path, tt.body)
			if tt.betaHeader != "" {
				c.Request.Header.Set("x-codex-beta-features", tt.betaHeader)
		placeholder

			require.Equal(t, tt.wantBefore, requiresOpenAICompactAccount(c, tt.body))
			normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), tt.body)
			require.True(t, ok)
			require.Equal(t, tt.wantAfter, requiresOpenAICompactAccount(c, normalized))
	placeholder)
placeholder
placeholder

func TestRequiresOpenAICompactAccount_RemoteV2SupportsResponsesRootAliases(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.6-sol","stream":true,"input":[{"type":"compaction_trigger"placeholder]placeholder`)

	for _, path := range []string{
		"/v1/responses",
		"/v1/responses/",
		"/openai/v1/responses",
		"/responses",
		"/backend-api/codex/responses",
placeholder {
		t.Run(path, func(t *testing.T) {
			c := newCompactBodySignalTestContext(t, path, body)

			normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
			require.True(t, ok)
			require.Equal(t, path, c.Request.URL.Path)
			require.True(t, requiresOpenAICompactAccount(c, normalized))
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_BodySignalTrailingSlashPromoted(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`)
	c := newCompactBodySignalTestContext(t, "/v1/responses/", body)

	_, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	require.Equal(t, "/v1/responses/compact", c.Request.URL.Path)
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_CodexDirectAliasPromoted(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`)
	c := newCompactBodySignalTestContext(t, "/backend-api/codex/responses", body)

	_, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	require.Equal(t, "/backend-api/codex/responses/compact", c.Request.URL.Path)
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_NonRemoteV2BodySignalPromoted(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	tests := []struct {
		name       string
		body       []byte
		betaHeader string
		wantMarked bool
placeholder{
		{
			name: "stream_false_headerless",
			body: []byte(`{"model":"gpt-5.5","stream":false,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
	placeholder,
		{
			name: "stream_absent_headerless",
			body: []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`),
	placeholder,
		{
			name:       "stream_false_declared_header",
			body:       []byte(`{"model":"gpt-5.5","stream":false,"input":[{"type":"compaction_trigger"placeholder]placeholder`),
			betaHeader: "remote_compaction_v2",
	placeholder,
		{
			name:       "stream_absent_wrong_case_header",
			body:       []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`),
			betaHeader: "REMOTE_COMPACTION_V2",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCompactBodySignalTestContext(t, "/v1/responses", tt.body)
			if tt.betaHeader != "" {
				c.Request.Header.Set("x-codex-beta-features", tt.betaHeader)
		placeholder

			normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), tt.body)
			require.True(t, ok)
			require.Equal(t, "/v1/responses/compact", c.Request.URL.Path)
			require.False(t, gjson.GetBytes(normalized, "stream").Exists())

			marked, exists := c.Get(service.OpenAICompactClientStreamKeyForTest())
			require.Equal(t, tt.wantMarked, exists)
			if tt.wantMarked {
				require.Equal(t, true, marked)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_NoTriggerUntouched(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.5","stream":true,"input":[{"type":"message","role":"user","content":"hello"placeholder]placeholder`)
	c := newCompactBodySignalTestContext(t, "/v1/responses", body)

	normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	require.Equal(t, "/v1/responses", c.Request.URL.Path)
	require.False(t, isOpenAIRemoteCompactPath(c))
	require.Equal(t, body, normalized)
	require.True(t, gjson.GetBytes(normalized, "stream").Bool())
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_PathBasedNoDoubleSuffix(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.5","stream":true,"store":true,"input":[{"type":"message","role":"user","content":"hello"placeholder]placeholder`)
	c := newCompactBodySignalTestContext(t, "/v1/responses/compact", body)

	normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	require.Equal(t, "/v1/responses/compact", c.Request.URL.Path)
	require.False(t, gjson.GetBytes(normalized, "stream").Exists())
	require.False(t, gjson.GetBytes(normalized, "store").Exists())
placeholder

func TestNormalizeOpenAIResponsesCompactRequest_SubpathNotPromoted(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.5","input":[{"type":"compaction_trigger"placeholder]placeholder`)
	c := newCompactBodySignalTestContext(t, "/v1/responses/resp_123/cancel", body)

	normalized, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	require.Equal(t, "/v1/responses/resp_123/cancel", c.Request.URL.Path)
	require.Equal(t, body, normalized)
placeholder

// path-based compact（Codex v1 unary 协议）即使 body 带 stream:true 也不标记，
// 保持 JSON 写回行为不变。
func TestNormalizeOpenAIResponsesCompactRequest_PathBasedStreamTrueNotMarked(t *testing.T) {
	h := &OpenAIGatewayHandler{placeholder
	body := []byte(`{"model":"gpt-5.5","stream":true,"input":[{"type":"message","role":"user","content":"hello"placeholder]placeholder`)
	c := newCompactBodySignalTestContext(t, "/v1/responses/compact", body)

	_, ok := h.normalizeOpenAIResponsesCompactRequest(c, zap.NewNop(), body)
	require.True(t, ok)
	_, exists := c.Get(service.OpenAICompactClientStreamKeyForTest())
	require.False(t, exists)
placeholder
