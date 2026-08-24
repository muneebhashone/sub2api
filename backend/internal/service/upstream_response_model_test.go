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

func TestUpstreamModelMismatchTreatsGrokBuildRuntimeIDsAsAliases(t *testing.T) {
	tests := []struct {
		name          string
		sentModel     string
		responseModel string
placeholder{
		{
			name:          "issue 5634 grok 4.6",
			sentModel:     "grok-4.6",
			responseModel: "grok-4.6-build",
	placeholder,
		{
			name:          "grok 4.6 latest",
			sentModel:     "grok-4.6-latest",
			responseModel: "grok-4.6-build",
	placeholder,
		{
			name:          "issue 5647 grok 4.5 latest",
			sentModel:     "grok-4.5-latest",
			responseModel: "grok-4.5-build",
	placeholder,
		{
			name:          "grok 4.5 canonical",
			sentModel:     "grok-4.5",
			responseModel: "GROK-4.5-BUILD",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mismatch := upstreamModelMismatch(tt.sentModel, tt.responseModel)

			require.NotNil(t, mismatch)
			require.False(t, *mismatch)
	placeholder)
placeholder
placeholder

func TestUpstreamModelMismatchDoesNotCollapseDifferentModels(t *testing.T) {
	tests := []struct {
		name          string
		sentModel     string
		responseModel string
placeholder{
		{
			name:          "different grok versions",
			sentModel:     "grok-4.5",
			responseModel: "grok-4.6-build",
	placeholder,
		{
			name:          "unrelated build suffix",
			sentModel:     "gpt-5.5",
			responseModel: "gpt-5.5-build",
	placeholder,
		{
			name:          "different grok runtime",
			sentModel:     "grok-build-0.1",
			responseModel: "grok-4.5-build",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mismatch := upstreamModelMismatch(tt.sentModel, tt.responseModel)

			require.NotNil(t, mismatch)
			require.True(t, *mismatch)
	placeholder)
placeholder
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

func TestUpstreamResponseModelObserverServiceTierTerminalEventWins(t *testing.T) {
	observer := &upstreamResponseModelObserver{placeholder

	// response.created echoes the requested tier and must not count as a declaration.
	observer.ObserveOpenAI([]byte(`{"type":"response.created","response":{"id":"resp_1","model":"gpt-5.6-sol","service_tier":"priority"placeholderplaceholder`), "response.created")
	require.Equal(t, "", observer.ServiceTier())

	observer.ObserveOpenAI([]byte(`{"type":"response.completed","response":{"id":"resp_1","model":"gpt-5.6-sol","service_tier":"default","usage":{"input_tokens":1,"output_tokens":2placeholderplaceholderplaceholder`), "response.completed")
	require.Equal(t, "default", observer.ServiceTier())
	require.False(t, observer.Conflict(), "tier echo must not be reported as a model conflict")
placeholder

func TestUpstreamResponseModelObserverServiceTierUntypedPayloads(t *testing.T) {
	t.Run("chat completions chunks agree", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveOpenAI([]byte(`{"id":"chatcmpl-1","object":"chat.completion.chunk","model":"gpt-5.4","service_tier":"default","choices":[]placeholder`), "")
		observer.ObserveOpenAI([]byte(`{"id":"chatcmpl-1","object":"chat.completion.chunk","model":"gpt-5.4","service_tier":"default","choices":[]placeholder`), "")
		require.Equal(t, "default", observer.ServiceTier())
placeholder)

	t.Run("disagreeing chunks are not trusted", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveOpenAI([]byte(`{"model":"gpt-5.4","service_tier":"priority"placeholder`), "")
		observer.ObserveOpenAI([]byte(`{"model":"gpt-5.4","service_tier":"default"placeholder`), "")
		require.Equal(t, "", observer.ServiceTier())
placeholder)

	t.Run("non-stream body normalises fast and drops auto", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveOpenAI([]byte(`{"id":"resp_1","object":"response","model":"gpt-5.6-sol","service_tier":"fast"placeholder`), "")
		require.Equal(t, "priority", observer.ServiceTier())

		auto := &upstreamResponseModelObserver{placeholder
		auto.ObserveOpenAI([]byte(`{"id":"resp_2","object":"response","model":"gpt-5.6-sol","service_tier":"auto"placeholder`), "")
		require.Equal(t, "", auto.ServiceTier())
placeholder)

	t.Run("model-free deltas never declare a tier", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveOpenAI([]byte(`{"type":"response.output_text.delta","delta":"hi","service_tier":"default"placeholder`), "response.output_text.delta")
		require.Equal(t, "", observer.ServiceTier())
placeholder)
placeholder

func TestUpstreamResponseModelObserverServiceTierAnthropicSpeed(t *testing.T) {
	t.Run("message_start reports standard", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveAnthropic([]byte(`{"type":"message_start","message":{"id":"msg_1","model":"claude-opus-5","usage":{"input_tokens":8,"output_tokens":1,"speed":"standard"placeholderplaceholderplaceholder`))
		observer.ObserveAnthropic([]byte(`{"type":"message_delta","delta":{"stop_reason":"end_turn"placeholder,"usage":{"output_tokens":12placeholderplaceholder`))
		require.Equal(t, "standard", observer.ServiceTier())
		require.Equal(t, "claude-opus-5", observer.Model())
placeholder)

	t.Run("non-stream body reports fast", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveAnthropic([]byte(`{"id":"msg_1","type":"message","model":"claude-opus-5","usage":{"input_tokens":8,"output_tokens":12,"speed":"fast"placeholderplaceholder`))
		require.Equal(t, "fast", observer.ServiceTier())
placeholder)

	t.Run("missing speed declares nothing", func(t *testing.T) {
		observer := &upstreamResponseModelObserver{placeholder
		observer.ObserveAnthropic([]byte(`{"id":"msg_1","type":"message","model":"claude-opus-5","usage":{"input_tokens":8,"output_tokens":12placeholderplaceholder`))
		require.Equal(t, "", observer.ServiceTier())
placeholder)
placeholder

func TestObservedUpstreamResponseServiceTierFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	require.Equal(t, "", observedUpstreamResponseServiceTier(c))

	observer := beginUpstreamResponseModelObservation(c)
	observer.ObserveOpenAI([]byte(`{"type":"response.completed","response":{"id":"resp_1","model":"gpt-5.6-sol","service_tier":"default"placeholderplaceholder`), "response.completed")
	require.Equal(t, "default", observedUpstreamResponseServiceTier(c))
placeholder
