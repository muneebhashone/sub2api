//go:build unit

package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newSessionHeaderContext(t *testing.T, headers map[string]string) *gin.Context {
placeholder
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
placeholder
	c.Request = req
	return c
placeholder

func TestSanitizeSessionID(t *testing.T) {
	longRunes := strings.Repeat("a", maxPersistedSessionIDLength+50)
	multibyte := strings.Repeat("好", maxPersistedSessionIDLength+10)

	tests := []struct {
		name string
		in   string
		want string
placeholder{
		{"empty", "", ""placeholder,
		{"whitespace only", "   \t  ", ""placeholder,
		{"trims surrounding whitespace", "  sess-123  ", "sess-123"placeholder,
		{"plain value", "conv_abc-123.XYZ", "conv_abc-123.XYZ"placeholder,
		{"uuid", "550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440000"placeholder,
		{"reject CR", "sess\r123", ""placeholder,
		{"reject LF", "sess\n123", ""placeholder,
		{"reject CRLF injection", "sess-1\r\nSet-Cookie: x=y", ""placeholder,
		{"reject tab inside", "sess\t123", ""placeholder,
		{"reject NUL", "sess\x00123", ""placeholder,
		{"reject DEL", "sess\x7f123", ""placeholder,
		{"reject invalid UTF-8", string([]byte{'s', 'e', 's', 's', '-', 0xffplaceholder), ""placeholder,
		{"accepts value at column bound", strings.Repeat("b", maxPersistedSessionIDLength), strings.Repeat("b", maxPersistedSessionIDLength)placeholder,
		{"rejects overlong value", longRunes, ""placeholder,
		{"rejects overlong multibyte value", multibyte, ""placeholder,
placeholder
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sanitizeSessionID(tc.in)
			require.Equal(t, tc.want, got, "sanitizeSessionID(%q)", tc.in)
			// Sanitized output must never exceed the DB column bound (rune-counted).
			require.LessOrEqual(t, len([]rune(got)), maxPersistedSessionIDLength)
	placeholder)
placeholder
placeholder

func TestExtractClientSessionID_NilContext(t *testing.T) {
	require.Equal(t, "", ExtractClientSessionID(nil))
placeholder

func TestExtractClientSessionID_NilRequest(t *testing.T) {
	require.Equal(t, "", ExtractClientSessionID(&gin.Context{placeholder))
placeholder

func TestExtractClientSessionID_AbsentReturnsEmpty(t *testing.T) {
	c := newSessionHeaderContext(t, nil)
	require.Equal(t, "", ExtractClientSessionID(c))
placeholder

func TestExtractClientSessionID_SupportedHeaders(t *testing.T) {
	tests := []struct {
		name   string
		header string
		value  string
placeholder{
		{"session_id", "session_id", "sess-A"placeholder,
		{"conversation_id", "conversation_id", "conv-B"placeholder,
		{"X-Session-Affinity", openCodeSessionAffinityHeader, "aff-C"placeholder,
		{"X-Session-Id", openCodeSessionIDHeader, "sid-D"placeholder,
		{"X-OpenCode-Session", openCodeNativeSessionHeader, "oc-E"placeholder,
		{"X-Conversation-ID", codeBuddyConversationHeader, "cb-F"placeholder,
		{"X-Claude-Code-Session-Id", claudeCodeSessionHeader, "cc-G"placeholder,
placeholder
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := newSessionHeaderContext(t, map[string]string{tc.header: tc.valueplaceholder)
			require.Equal(t, tc.value, ExtractClientSessionID(c))
	placeholder)
placeholder
placeholder

func TestExtractClientSessionID_HeaderPrecedence(t *testing.T) {
	// session_id ranks ahead of conversation_id and the X-* variants.
	c := newSessionHeaderContext(t, map[string]string{
		"session_id":                "primary",
		"conversation_id":           "secondary",
		openCodeSessionIDHeader:     "tertiary",
		codeBuddyConversationHeader: "quaternary",
placeholder)
	require.Equal(t, "primary", ExtractClientSessionID(c))
placeholder

func TestExtractClientSessionID_Sanitizes(t *testing.T) {
	c := newSessionHeaderContext(t, map[string]string{openCodeSessionIDHeader: "  clean-123  "placeholder)
	require.Equal(t, "clean-123", ExtractClientSessionID(c))
placeholder

func TestExtractClientSessionID_IgnoresNonSessionHeaders(t *testing.T) {
	// prompt_cache_key, request/message ids, and a Grok conversation header on a
	// non-Grok request are NOT persisted as session_id.
	c := newSessionHeaderContext(t, map[string]string{
		"prompt_cache_key": "cache-key-should-not-persist",
		"X-Request-Id":     "req-should-not-persist",
		"x-grok-conv-id":   "grok-conv-should-not-persist",
placeholder)
	require.Equal(t, "", ExtractClientSessionID(c))
placeholder

func TestExtractClientSessionID_GrokConversationHeader(t *testing.T) {
	c := newSessionHeaderContext(t, map[string]string{
		grokConversationIDHeader: "grok-native-session",
placeholder)
	c.Set("api_key", &APIKey{
		ID:    42,
		Group: &Group{Platform: PlatformGrokplaceholder,
placeholder)

	require.Equal(t, "grok-native-session", ExtractClientSessionID(c))
placeholder

func TestExtractClientSessionID_GrokConversationHeaderForCompositeRoute(t *testing.T) {
	c := newSessionHeaderContext(t, map[string]string{
		grokConversationIDHeader: "grok-composite-session",
placeholder)
	c.Set("api_key", &APIKey{
		ID:    43,
		Group: &Group{Platform: PlatformCompositeplaceholder,
placeholder)
	c.Request = c.Request.WithContext(WithResolvedTargetPlatform(context.Background(), PlatformGrok))

	require.Equal(t, "grok-composite-session", ExtractClientSessionID(c))
placeholder

func TestExtractClientSessionID_InjectionHeaderDropped(t *testing.T) {
	// A supported header carrying a CRLF payload is rejected, not persisted mangled.
	c := newSessionHeaderContext(t, map[string]string{"session_id": "abc"placeholder)
	c.Request.Header.Set("session_id", "abc\r\nX-Injected: 1")
	require.Equal(t, "", ExtractClientSessionID(c))
placeholder
