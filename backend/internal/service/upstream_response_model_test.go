package service

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestUpstreamResponseModelObserverTerminalWinsAndRecordsConflict(t *testing.T) {
	observer := &upstreamResponseModelObserver{placeholder

	observer.ObserveOpenAI([]byte(`{"type":"response.created","response":{"model":"gpt-5.5"placeholderplaceholder`), "response.created")
	observer.ObserveOpenAI([]byte(`{"type":"response.completed","response":{"model":"gpt-5.4"placeholderplaceholder`), "response.completed")

	require.Equal(t, "gpt-5.4", observer.Model())
	require.True(t, observer.Conflict())
placeholder

func TestUpstreamResponseModelObserverSupportsAnthropicAndGeminiShapes(t *testing.T) {
	t.Run("anthropic", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveAnthropic([]byte(`{"type":"message_start","message":{"model":"claude-sonnet-4-20250514"placeholderplaceholder`))
		require.Equal(t, "claude-sonnet-4-20250514", observer.Model())
placeholder)

	t.Run("gemini outer and nested", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveGemini([]byte(`{"response":{"modelVersion":"gemini-2.5-pro"placeholderplaceholder`))
		observer.ObserveGemini([]byte(`{"modelVersion":"gemini-2.5-pro-latest"placeholder`))
		require.Equal(t, "gemini-2.5-pro-latest", observer.Model())
		require.True(t, observer.Conflict())
placeholder)
placeholder

func TestUpstreamResponseModelObservationAttemptReset(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	first := beginUpstreamResponseModelObservation(c)
	first.Observe("failed-attempt-model", false)
	second := beginUpstreamResponseModelObservation(c)
	second.Observe("successful-attempt-model", false)

	require.Equal(t, "successful-attempt-model", observedUpstreamResponseModel(c))
	require.False(t, observedUpstreamResponseModelConflict(c))
placeholder

func TestUpstreamModelMismatchThreeStateAndCaseInsensitiveComparison(t *testing.T) {
	require.Nil(t, upstreamModelMismatch("gpt-5.5", ""))

	matched := upstreamModelMismatch("gpt-5.5", "GPT-5.5")
	require.NotNil(t, matched)
	require.False(t, *matched)

	mismatched := upstreamModelMismatch("gpt-5.5", "gpt-5.4")
	require.NotNil(t, mismatched)
	require.True(t, *mismatched)
placeholder

func TestObserveOpenAISSEBodyIgnoresMalformedPayload(t *testing.T) {
	observer := &upstreamResponseModelObserver{placeholder
	observeOpenAISSEBody(observer, "data: not-json\n\ndata: {\"type\":\"response.completed\",\"response\":{\"model\":\"gpt-5.4\"placeholderplaceholder\n\n")

	require.Equal(t, "gpt-5.4", observer.Model())
	require.False(t, observer.Conflict())
placeholder

func TestObserveAntigravityGeminiSSELineReadsWrapperModelWithoutUnwrap(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		want    string
placeholder{
		{
			name:    "top-level sibling",
			payload: `{"modelVersion":"gemini-3-pro","response":{"candidates":[]placeholderplaceholder`,
			want:    "gemini-3-pro",
	placeholder,
		{
			name:    "single wrapper",
			payload: `{"response":{"modelVersion":"gemini-3-pro","candidates":[]placeholderplaceholder`,
			want:    "gemini-3-pro",
	placeholder,
		{
			name:    "nested response after one wrapper",
			payload: `{"response":{"response":{"modelVersion":"gemini-3-pro","candidates":[]placeholderplaceholderplaceholder`,
			want:    "gemini-3-pro",
	placeholder,
		{
			name:    "outer declaration takes precedence",
			payload: `{"modelVersion":"gemini-outer","response":{"modelVersion":"gemini-inner","candidates":[]placeholderplaceholder`,
			want:    "gemini-outer",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			c, _ := gin.CreateTestContext(nil)
			beginUpstreamResponseModelObservation(c)

			svc := &AntigravityGatewayService{placeholder
			svc.observeAntigravityGeminiSSELine(c, "data: "+tt.payload)

			require.Equal(t, tt.want, observedUpstreamResponseModel(c))
			require.False(t, observedUpstreamResponseModelConflict(c))
	placeholder)
placeholder
placeholder

func TestUpstreamResponseModelObserverRejectsMalformedJSONWithModelField(t *testing.T) {
	observer := &upstreamResponseModelObserver{placeholder
	observer.ObserveOpenAI([]byte(`{"response":{"model":"gpt-5.4"placeholder`), "response.completed")

	require.Empty(t, observer.Model())
placeholder

func TestUpstreamResponseModelObserverBoundsUntrustedModelName(t *testing.T) {
	observer := &upstreamResponseModelObserver{placeholder
	observer.Observe("  "+strings.Repeat("模", upstreamResponseModelMaxLength+1)+"  ", false)

	require.Len(t, []rune(observer.Model()), upstreamResponseModelMaxLength)
placeholder
