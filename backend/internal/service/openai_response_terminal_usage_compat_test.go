//go:build unit

package service

import (
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/stretchr/testify/require"
)

func TestEffectiveOpenAISSEEventTypePrefersPayload(t *testing.T) {
	t.Parallel()

	require.Equal(t, "response.failed", effectiveOpenAISSEEventType([]byte(`{"type":"response.failed"placeholder`), "error"))
	require.Equal(t, "error", effectiveOpenAISSEEventType([]byte(`{"error":{"message":"failed"placeholderplaceholder`), " error "))
	require.JSONEq(t, `{"type":"error","error":{"message":"failed"placeholderplaceholder`, openAICompatPayloadWithEventType(`{"type":"","error":{"message":"failed"placeholderplaceholder`, "error"))
placeholder

func TestExtractOpenAISSETerminalEventUsesEventField(t *testing.T) {
	t.Parallel()

	body := "event: error\n" +
		"data: {\"error\":{\"message\":\"provider failed\"placeholderplaceholder\n\n" +
		"data: [DONE]\n\n"
	eventType, payload, ok := extractOpenAISSETerminalEvent(body)
	require.True(t, ok)
	require.Equal(t, "error", eventType)
	require.Equal(t, "provider failed", extractOpenAISSEErrorMessage(payload))
placeholder

func TestParseSSEUsageEffectiveTerminalRules(t *testing.T) {
	t.Parallel()

	svc := &OpenAIGatewayService{placeholder
	usage := &OpenAIUsage{placeholder
	svc.parseSSEUsageBytesWithType([]byte(`{"usage":{"input_tokens":17,"output_tokens":5,"input_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholder`), "response.in_progress", usage)
	svc.parseSSEUsageBytesWithType([]byte(`{"response":{"id":"resp_1"placeholderplaceholder`), "response.completed", usage)
	require.Equal(t, OpenAIUsage{InputTokens: 17, OutputTokens: 5, CacheReadInputTokens: 3placeholder, *usage)

	svc.parseSSEUsageBytesWithType([]byte(`{"response":{"usage":{"input_tokens":0,"output_tokens":0,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholderplaceholder`), "response.completed", usage)
	require.Equal(t, OpenAIUsage{InputTokens: 17, OutputTokens: 5, CacheReadInputTokens: 3placeholder, *usage)

	svc.parseSSEUsageBytesWithType([]byte(`{"response":{"usage":{"input_tokens":2,"output_tokens":0,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholderplaceholder`), "response.completed", usage)
	require.Equal(t, OpenAIUsage{InputTokens: 2placeholder, *usage)
placeholder

func TestOpenAICompatTerminalResponseSynthesizesBareError(t *testing.T) {
	t.Parallel()

	payload := []byte(`{"type":"error","code":"upstream_error","error":{"message":"provider failed"placeholderplaceholder`)
	event := &apicompat.ResponsesStreamEvent{Type: "error", Code: "upstream_error"placeholder
	response := openAICompatTerminalResponse(event, payload)
	require.NotNil(t, response)
	require.Equal(t, "failed", response.Status)
	require.Equal(t, "upstream_error", response.Error.Code)
	require.Equal(t, "provider failed", response.Error.Message)
placeholder

func TestForEachOpenAISSEFrameDataTypeOverridesEventField(t *testing.T) {
	t.Parallel()

	var types []string
	forEachOpenAISSEFrame(strings.Join([]string{
		"event: response.in_progress",
		`data: {"type":"response.completed","response":{"id":"resp_1"placeholderplaceholder`,
		"",
placeholder, "\n"), func(eventType string, _ []byte) {
		types = append(types, eventType)
placeholder)
	require.Equal(t, []string{"response.completed"placeholder, types)
placeholder

func TestSampleOpenAIMissingUsageLogRateLimitsAndCounts(t *testing.T) {
	var sampler openAIMissingUsageLogSampler
	now := time.Unix(100, 0)
	logNow, total, suppressed := sampler.sample(now)
	require.True(t, logNow)
	require.Equal(t, uint64(1), total)
	require.Zero(t, suppressed)

	logNow, total, _ = sampler.sample(now.Add(time.Second))
	require.False(t, logNow)
	require.Equal(t, uint64(2), total)

	logNow, total, suppressed = sampler.sample(now.Add(openAIMissingUsageLogInterval))
	require.True(t, logNow)
	require.Equal(t, uint64(3), total)
	require.Equal(t, uint64(1), suppressed)
placeholder
