package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) placeholder

// ──────────────────────────────────────────────────────────
// NormalizeInboundEndpoint
// ──────────────────────────────────────────────────────────

func TestNormalizeInboundEndpoint(t *testing.T) {
	tests := []struct {
		path string
		want string
placeholder{
		// Direct canonical paths.
		{"/v1/messages", EndpointMessagesplaceholder,
		{"/v1/chat/completions", EndpointChatCompletionsplaceholder,
		{"/v1/embeddings", EndpointEmbeddingsplaceholder,
		{"/v1/responses", EndpointResponsesplaceholder,
		{"/v1/responses/compact", EndpointResponsesCompactplaceholder,
		{"/v1/responses/compact/detail", EndpointResponsesCompactplaceholder,
		{"/v1/images/generations", EndpointImagesGenerationsplaceholder,
		{"/v1/images/edits", EndpointImagesEditsplaceholder,
		{"/v1/videos/generations", EndpointVideosGenerationsplaceholder,
		{"/v1/videos/req_123", EndpointVideosplaceholder,
		{"/v1beta/models", EndpointGeminiModelsplaceholder,

		// Prefixed paths (antigravity, openai) — root Responses.
		{"/antigravity/v1/messages", EndpointMessagesplaceholder,
		{"/openai/v1/responses", EndpointResponsesplaceholder,
		{"/openai/v1/images/generations", EndpointImagesGenerationsplaceholder,
		{"/openai/v1/images/edits", EndpointImagesEditsplaceholder,
		{"/antigravity/v1beta/models/gemini:generateContent", EndpointGeminiModelsplaceholder,

		// Prefixed paths — "/responses/compact" is its OWN distinct
		// inbound endpoint, not folded into the root Responses endpoint.
		{"/openai/v1/responses/compact", EndpointResponsesCompactplaceholder,
		{"/openai/v1/responses/compact/detail", EndpointResponsesCompactplaceholder,

		// Bare top-level alias route "/responses" — root vs. compact.
		{"/responses", EndpointResponsesplaceholder,
		{"/responses/compact", EndpointResponsesCompactplaceholder,
		{"/responses/compact/detail", EndpointResponsesCompactplaceholder,

		// Bare Codex direct alias route — root vs. compact.
		{"/backend-api/codex/responses", EndpointResponsesplaceholder,
		{"/backend-api/codex/responses/compact", EndpointResponsesCompactplaceholder,
		{"/backend-api/codex/responses/compact/detail", EndpointResponsesCompactplaceholder,

		// Must NOT generalize to arbitrary paths merely ending in
		// "/responses" (or "/responses/compact") that are unrelated to
		// the two known bare alias roots, unless they already carry a
		// supported "/v1/responses..." prefix form.
		{"/foo/responses", "/foo/responses"placeholder,
		{"/foo/responses/compact", "/foo/responses/compact"placeholder,

		// Gin route patterns with wildcards. The literal wildcard token
		// ("*subpath") is not the "compact" segment itself, so these
		// generic FullPath patterns normalize to the root Responses
		// endpoint; only a concrete "compact" path segment (tested above)
		// resolves to EndpointResponsesCompact.
		{"/v1beta/models/*modelAction", EndpointGeminiModelsplaceholder,
		{"/v1/responses/*subpath", EndpointResponsesplaceholder,
		{"/responses/*subpath", EndpointResponsesplaceholder,
		{"/backend-api/codex/responses/*subpath", EndpointResponsesplaceholder,

		// Unknown path is returned as-is.
		{"/v1/embeddings", "/v1/embeddings"placeholder,
		{"", ""placeholder,
		{"  /v1/messages  ", EndpointMessagesplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			require.Equal(t, tt.want, NormalizeInboundEndpoint(tt.path))
	placeholder)
placeholder
placeholder

// ──────────────────────────────────────────────────────────
// DeriveUpstreamEndpoint
// ──────────────────────────────────────────────────────────

func TestDeriveUpstreamEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		inbound  string
		rawPath  string
		platform string
		want     string
placeholder{
		// Anthropic.
		{"anthropic messages", EndpointMessages, "/v1/messages", service.PlatformAnthropic, EndpointMessagesplaceholder,

		// Gemini.
		{"gemini models", EndpointGeminiModels, "/v1beta/models/gemini:gen", service.PlatformGemini, EndpointGeminiModelsplaceholder,

		// OpenAI — root Responses.
		{"openai responses root", EndpointResponses, "/v1/responses", service.PlatformOpenAI, EndpointResponsesplaceholder,

		// OpenAI — compact, raw path carries the derivable "/compact"
		// (or nested) suffix, which must be preserved on the upstream
		// endpoint.
		{"openai responses compact", EndpointResponsesCompact, "/openai/v1/responses/compact", service.PlatformOpenAI, "/v1/responses/compact"placeholder,
		{"openai responses nested", EndpointResponsesCompact, "/openai/v1/responses/compact/detail", service.PlatformOpenAI, "/v1/responses/compact/detail"placeholder,
		{"openai bare responses compact", EndpointResponsesCompact, "/responses/compact", service.PlatformOpenAI, "/v1/responses/compact"placeholder,
		{"openai bare responses compact detail", EndpointResponsesCompact, "/responses/compact/detail", service.PlatformOpenAI, "/v1/responses/compact/detail"placeholder,
		{"openai codex direct responses compact", EndpointResponsesCompact, "/backend-api/codex/responses/compact", service.PlatformOpenAI, "/v1/responses/compact"placeholder,
		{"openai codex direct responses compact detail", EndpointResponsesCompact, "/backend-api/codex/responses/compact/detail", service.PlatformOpenAI, "/v1/responses/compact/detail"placeholder,

		// OpenAI — bare root alias routes normalize to root Responses.
		{"openai bare responses", EndpointResponses, "/responses", service.PlatformOpenAI, EndpointResponsesplaceholder,
		{"openai codex direct responses", EndpointResponses, "/backend-api/codex/responses", service.PlatformOpenAI, EndpointResponsesplaceholder,

		// OpenAI — inbound is already the canonical compact endpoint but
		// the raw path carries no derivable "/responses..." suffix (e.g.
		// it was already normalized upstream). Must not silently fall
		// back to the root Responses endpoint.
		{"openai responses compact inbound only, unrelated raw path", EndpointResponsesCompact, "/v1/messages", service.PlatformOpenAI, EndpointResponsesCompactplaceholder,

		{"openai from messages", EndpointMessages, "/v1/messages", service.PlatformOpenAI, EndpointResponsesplaceholder,
		{"openai from completions", EndpointChatCompletions, "/v1/chat/completions", service.PlatformOpenAI, EndpointResponsesplaceholder,
		{"openai embeddings", EndpointEmbeddings, "/v1/embeddings", service.PlatformOpenAI, EndpointEmbeddingsplaceholder,
		{"openai image generations", EndpointImagesGenerations, "/v1/images/generations", service.PlatformOpenAI, EndpointImagesGenerationsplaceholder,
		{"openai image edits", EndpointImagesEdits, "/openai/v1/images/edits", service.PlatformOpenAI, EndpointImagesEditsplaceholder,
		{"grok video generations", EndpointVideosGenerations, "/v1/videos/generations", service.PlatformGrok, EndpointVideosGenerationsplaceholder,
		{"grok video status", EndpointVideos, "/videos/req_123", service.PlatformGrok, EndpointVideosplaceholder,

		// Antigravity — uses inbound to pick Claude vs Gemini upstream.
		{"antigravity claude", EndpointMessages, "/antigravity/v1/messages", service.PlatformAntigravity, EndpointMessagesplaceholder,
		{"antigravity gemini", EndpointGeminiModels, "/antigravity/v1beta/models", service.PlatformAntigravity, EndpointGeminiModelsplaceholder,

		// Unknown platform — passthrough.
		{"unknown platform", "/v1/embeddings", "/v1/embeddings", "unknown", "/v1/embeddings"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, DeriveUpstreamEndpoint(tt.inbound, tt.rawPath, tt.platform))
	placeholder)
placeholder
placeholder

// ──────────────────────────────────────────────────────────
// responsesSubpathSuffix
// ──────────────────────────────────────────────────────────

func TestResponsesSubpathSuffix(t *testing.T) {
	tests := []struct {
		raw  string
		want string
placeholder{
		{"/v1/responses", ""placeholder,
		{"/v1/responses/", ""placeholder,
		{"/v1/responses/compact", "/compact"placeholder,
		{"/openai/v1/responses/compact/detail", "/compact/detail"placeholder,
		{"/responses", ""placeholder,
		{"/responses/compact", "/compact"placeholder,
		{"/responses/compact/detail", "/compact/detail"placeholder,
		{"/backend-api/codex/responses", ""placeholder,
		{"/backend-api/codex/responses/compact", "/compact"placeholder,
		{"/backend-api/codex/responses/compact/detail", "/compact/detail"placeholder,
		{"/v1/messages", ""placeholder,
		{"", ""placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			require.Equal(t, tt.want, responsesSubpathSuffix(tt.raw))
	placeholder)
placeholder
placeholder

// ──────────────────────────────────────────────────────────
// InboundEndpointMiddleware + context helpers
// ──────────────────────────────────────────────────────────

func TestInboundEndpointMiddleware(t *testing.T) {
	router := gin.New()
	router.Use(InboundEndpointMiddleware())

	var captured string
	router.POST("/v1/messages", func(c *gin.Context) {
		captured = GetInboundEndpoint(c)
		c.Status(http.StatusOK)
placeholder)

	req := httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, EndpointMessages, captured)
placeholder

func TestGetInboundEndpoint_FallbackWithoutMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/antigravity/v1/messages", nil)

	// Middleware did not run — fallback to normalizing c.Request.URL.Path.
	got := GetInboundEndpoint(c)
	require.Equal(t, EndpointMessages, got)
placeholder

func TestGetUpstreamEndpoint_FullFlow(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses/compact", nil)

	// Simulate middleware.
	c.Set(ctxKeyInboundEndpoint, NormalizeInboundEndpoint(c.Request.URL.Path))

	got := GetUpstreamEndpoint(c, service.PlatformOpenAI)
	require.Equal(t, "/v1/responses/compact", got)
placeholder
