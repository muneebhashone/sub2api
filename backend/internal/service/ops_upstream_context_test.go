package service

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSafeUpstreamURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{"strips query", "https://api.anthropic.com/v1/messages?beta=true", "https://api.anthropic.com/v1/messages"placeholder,
		{"strips fragment", "https://api.openai.com/v1/responses#frag", "https://api.openai.com/v1/responses"placeholder,
		{"strips both", "https://host/path?token=secret#x", "https://host/path"placeholder,
		{"no query or fragment", "https://host/path", "https://host/path"placeholder,
		{"empty string", "", ""placeholder,
		{"whitespace only", "  ", ""placeholder,
		{"query before fragment", "https://h/p?a=1#f", "https://h/p"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, safeUpstreamURL(tt.input))
	placeholder)
placeholder
placeholder

func TestAppendOpsUpstreamError_UsesRequestBodyBytesFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	setOpsUpstreamRequestBody(c, []byte(`{"model":"gpt-5"placeholder`))
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Kind:    "http_error",
		Message: "upstream failed",
placeholder)

	v, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := v.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.Len(t, events, 1)
	require.Equal(t, `{"model":"gpt-5"placeholder`, events[0].UpstreamRequestBody)
placeholder

func TestAppendOpsUpstreamError_UsesRequestBodyStringFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	c.Set(OpsUpstreamRequestBodyKey, `{"model":"gpt-4"placeholder`)
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Kind:    "request_error",
		Message: "dial timeout",
placeholder)

	v, ok := c.Get(OpsUpstreamErrorsKey)
	require.True(t, ok)
	events, ok := v.([]*OpsUpstreamErrorEvent)
	require.True(t, ok)
	require.Len(t, events, 1)
	require.Equal(t, `{"model":"gpt-4"placeholder`, events[0].UpstreamRequestBody)
placeholder
