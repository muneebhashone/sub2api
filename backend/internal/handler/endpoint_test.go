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
		{"/v1/responses", EndpointResponsesplaceholder,
		{"/v1beta/models", EndpointGeminiModelsplaceholder,

		// Prefixed paths (antigravity, openai, sora).
		{"/antigravity/v1/messages", EndpointMessagesplaceholder,
		{"/openai/v1/responses", EndpointResponsesplaceholder,
		{"/openai/v1/responses/compact", EndpointResponsesplaceholder,
		{"/sora/v1/chat/completions", EndpointChatCompletionsplaceholder,
		{"/antigravity/v1beta/models/gemini:generateContent", EndpointGeminiModelsplaceholder,

		// Gin route patterns with wildcards.
		{"/v1beta/models/*modelAction", EndpointGeminiModelsplaceholder,
		{"/v1/responses/*subpath", EndpointResponsesplaceholder,

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

		// Sora.
		{"sora completions", EndpointChatCompletions, "/sora/v1/chat/completions", service.PlatformSora, EndpointChatCompletionsplaceholder,

		// OpenAI — always /v1/responses.
		{"openai responses root", EndpointResponses, "/v1/responses", service.PlatformOpenAI, EndpointResponsesplaceholder,
		{"openai responses compact", EndpointResponses, "/openai/v1/responses/compact", service.PlatformOpenAI, "/v1/responses/compact"placeholder,
		{"openai responses nested", EndpointResponses, "/openai/v1/responses/compact/detail", service.PlatformOpenAI, "/v1/responses/compact/detail"placeholder,
		{"openai from messages", EndpointMessages, "/v1/messages", service.PlatformOpenAI, EndpointResponsesplaceholder,
		{"openai from completions", EndpointChatCompletions, "/v1/chat/completions", service.PlatformOpenAI, EndpointResponsesplaceholder,

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
