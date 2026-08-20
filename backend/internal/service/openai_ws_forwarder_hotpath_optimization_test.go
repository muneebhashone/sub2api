package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOpenAIWSEventEnvelope(t *testing.T) {
	eventType, responseID, response := parseOpenAIWSEventEnvelope([]byte(`{"type":"response.completed","response":{"id":"resp_1","model":"gpt-5.1"placeholderplaceholder`))
	require.Equal(t, "response.completed", eventType)
	require.Equal(t, "resp_1", responseID)
	require.True(t, response.Exists())
	require.Equal(t, `{"id":"resp_1","model":"gpt-5.1"placeholder`, response.Raw)

	eventType, responseID, response = parseOpenAIWSEventEnvelope([]byte(`{"type":"response.delta","id":"evt_1"placeholder`))
	require.Equal(t, "response.delta", eventType)
	require.Equal(t, "evt_1", responseID)
	require.False(t, response.Exists())
placeholder

func TestParseOpenAIWSResponseUsageFromCompletedEvent(t *testing.T) {
	usage := &OpenAIUsage{placeholder
	parseOpenAIWSResponseUsageFromCompletedEvent(
		[]byte(`{"type":"response.completed","response":{"usage":{"input_tokens":11,"output_tokens":7,"input_tokens_details":{"cached_tokens":3placeholderplaceholderplaceholderplaceholder`),
		usage,
	)
	require.Equal(t, 11, usage.InputTokens)
	require.Equal(t, 7, usage.OutputTokens)
	require.Equal(t, 3, usage.CacheReadInputTokens)

	parseOpenAIWSResponseUsageFromCompletedEvent(
		[]byte(`{"type":"response.completed","response":{"usage":{"prompt_tokens":19,"completion_tokens":5,"prompt_tokens_details":{"cached_tokens":4placeholderplaceholderplaceholderplaceholder`),
		usage,
	)
	require.Equal(t, 19, usage.InputTokens)
	require.Equal(t, 5, usage.OutputTokens)
	require.Equal(t, 4, usage.CacheReadInputTokens)

	parseOpenAIWSResponseUsageFromCompletedEvent(
		[]byte(`{"type":"response.completed","response":{"usage":{"input_tokens":0,"output_tokens":0,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholderplaceholder`),
		usage,
	)
	require.Equal(t, OpenAIUsage{InputTokens: 19, OutputTokens: 5, CacheReadInputTokens: 4placeholder, *usage)

	parseOpenAIWSResponseUsageFromCompletedEvent(
		[]byte(`{"type":"response.failed","response":{"usage":{"input_tokens":3,"output_tokens":0,"input_tokens_details":{"cached_tokens":0placeholderplaceholderplaceholderplaceholder`),
		usage,
	)
	require.Equal(t, OpenAIUsage{InputTokens: 3placeholder, *usage)
placeholder

func TestOpenAIWSEventShouldParseUsageTerminalEvents(t *testing.T) {
	t.Parallel()

	for _, eventType := range []string{
		"response.completed",
		"response.done",
		"response.failed",
		"response.incomplete",
		"response.cancelled",
		"response.canceled",
placeholder {
		require.True(t, openAIWSEventShouldParseUsage(eventType), eventType)
		require.True(t, openAIWSEventShouldParseUsage("  "+eventType+"  "), eventType)
placeholder
	require.False(t, openAIWSEventShouldParseUsage("response.output_text.delta"))
	require.True(t, openAIWSEventShouldParseUsage("response.output_text.done"))
	require.False(t, openAIWSEventShouldParseUsage(""))
	require.False(t, openAIWSMessageShouldParseUsage("response.in_progress", []byte(`{"type":"response.in_progress"placeholder`)))
	require.True(t, openAIWSMessageShouldParseUsage("response.in_progress", []byte(`{"type":"response.in_progress","usage":{placeholderplaceholder`)))
	require.False(t, openAIWSMessageShouldParseUsage("response.output_text.delta", []byte(`{"type":"response.output_text.delta","usage":{placeholderplaceholder`)))
placeholder

func TestOpenAIWSErrorEventHelpers_ConsistentWithWrapper(t *testing.T) {
	message := []byte(`{"type":"error","error":{"type":"invalid_request_error","code":"invalid_request","message":"invalid input"placeholderplaceholder`)
	codeRaw, errTypeRaw, errMsgRaw := parseOpenAIWSErrorEventFields(message)

	wrappedReason, wrappedRecoverable := classifyOpenAIWSErrorEvent(message)
	rawReason, rawRecoverable := classifyOpenAIWSErrorEventFromRaw(codeRaw, errTypeRaw, errMsgRaw)
	require.Equal(t, wrappedReason, rawReason)
	require.Equal(t, wrappedRecoverable, rawRecoverable)

	wrappedStatus := openAIWSErrorHTTPStatus(message)
	rawStatus := openAIWSErrorHTTPStatusFromRaw(codeRaw, errTypeRaw)
	require.Equal(t, wrappedStatus, rawStatus)
	require.Equal(t, http.StatusBadRequest, rawStatus)

	wrappedCode, wrappedType, wrappedMsg := summarizeOpenAIWSErrorEventFields(message)
	rawCode, rawType, rawMsg := summarizeOpenAIWSErrorEventFieldsFromRaw(codeRaw, errTypeRaw, errMsgRaw)
	require.Equal(t, wrappedCode, rawCode)
	require.Equal(t, wrappedType, rawType)
	require.Equal(t, wrappedMsg, rawMsg)
placeholder

func TestOpenAIWSMessageLikelyContainsToolCalls(t *testing.T) {
	require.False(t, openAIWSMessageLikelyContainsToolCalls([]byte(`{"type":"response.output_text.delta","delta":"hello"placeholder`)))
	require.True(t, openAIWSMessageLikelyContainsToolCalls([]byte(`{"type":"response.output_item.added","item":{"tool_calls":[{"id":"tc1"placeholder]placeholderplaceholder`)))
	require.True(t, openAIWSMessageLikelyContainsToolCalls([]byte(`{"type":"response.output_item.added","item":{"type":"function_call"placeholderplaceholder`)))
placeholder

func TestReplaceOpenAIWSMessageModel_OptimizedStillCorrect(t *testing.T) {
	noModel := []byte(`{"type":"response.output_text.delta","delta":"hello"placeholder`)
	require.Equal(t, string(noModel), string(replaceOpenAIWSMessageModel(noModel, "gpt-5.1", "custom-model")))

	rootOnly := []byte(`{"type":"response.created","model":"gpt-5.1"placeholder`)
	require.Equal(t, `{"type":"response.created","model":"custom-model"placeholder`, string(replaceOpenAIWSMessageModel(rootOnly, "gpt-5.1", "custom-model")))

	responseOnly := []byte(`{"type":"response.completed","response":{"model":"gpt-5.1"placeholderplaceholder`)
	require.Equal(t, `{"type":"response.completed","response":{"model":"custom-model"placeholderplaceholder`, string(replaceOpenAIWSMessageModel(responseOnly, "gpt-5.1", "custom-model")))

	both := []byte(`{"model":"gpt-5.1","response":{"model":"gpt-5.1"placeholderplaceholder`)
	require.Equal(t, `{"model":"custom-model","response":{"model":"custom-model"placeholderplaceholder`, string(replaceOpenAIWSMessageModel(both, "gpt-5.1", "custom-model")))
placeholder
