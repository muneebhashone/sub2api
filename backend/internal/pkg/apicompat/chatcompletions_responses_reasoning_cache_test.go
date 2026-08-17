package apicompat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// Encrypted-only reasoning items (empty summary + opaque encrypted_content,
// e.g. after codex remote compaction) carry no plaintext the bridge can map to
// reasoning_content. The gateway-side cache keyed by reasoning item id restores
// it; without the restore, DeepSeek thinking mode rejects the history with 400
// "The `reasoning_content` in the thinking mode must be passed back to the API".
func TestResponsesToChat_ReasoningCacheLookup_RestoresEncryptedOnlyItem(t *testing.T) {
	req := &ResponsesRequest{
		Model: "deepseek-reasoner",
		Input: json.RawMessage(`[
			{"type":"reasoning","id":"item_enc1","summary":[],"encrypted_content":"opaque"placeholder,
			{"type":"function_call","call_id":"call_1","name":"get_value","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_1","output":"ok"placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"go on"placeholder]placeholder
		]`),
placeholder

	out, err := ResponsesToChatCompletionsRequestWithOptions(req, &ResponsesToChatOptions{
		ReasoningContentByID: func(itemID string) string {
			if itemID == "item_enc1" {
				return "cached thinking"
		placeholder
			return ""
	placeholder,
placeholder)
placeholder
	require.Len(t, out.Messages, 3)
	require.Equal(t, "assistant", out.Messages[0].Role)
	require.Equal(t, "cached thinking", out.Messages[0].ReasoningContent)
	require.Len(t, out.Messages[0].ToolCalls, 1)
	require.Equal(t, "call_1", out.Messages[0].ToolCalls[0].ID)
	require.Equal(t, "tool", out.Messages[1].Role)
	require.Equal(t, "user", out.Messages[2].Role)
placeholder

// A cache miss keeps the original behavior: no reasoning_content, no error.
func TestResponsesToChat_ReasoningCacheLookup_MissKeepsOriginalBehavior(t *testing.T) {
	req := &ResponsesRequest{
		Model: "deepseek-reasoner",
		Input: json.RawMessage(`[
			{"type":"reasoning","id":"item_unknown","summary":[],"encrypted_content":"opaque"placeholder,
			{"type":"function_call","call_id":"call_1","name":"get_value","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_1","output":"ok"placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"go on"placeholder]placeholder
		]`),
placeholder

	out, err := ResponsesToChatCompletionsRequestWithOptions(req, &ResponsesToChatOptions{
		ReasoningContentByID: func(string) string { return "" placeholder,
placeholder)
placeholder
	require.Len(t, out.Messages, 3)
	require.Empty(t, out.Messages[0].ReasoningContent)

	// Nil options (legacy path) behaves identically.
	legacy, err := ResponsesToChatCompletionsRequest(req)
placeholder
	require.Equal(t, out.Messages, legacy.Messages)
placeholder

// Plaintext summary wins and the cache lookup is not consulted.
func TestResponsesToChat_ReasoningCacheLookup_PlaintextPreferred(t *testing.T) {
	req := &ResponsesRequest{
		Model: "deepseek-reasoner",
		Input: json.RawMessage(`[
			{"type":"reasoning","id":"item_plain","summary":[{"type":"summary_text","text":"plain thinking"placeholder]placeholder,
			{"type":"function_call","call_id":"call_1","name":"get_value","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_1","output":"ok"placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"go on"placeholder]placeholder
		]`),
placeholder

	lookupCalled := false
	out, err := ResponsesToChatCompletionsRequestWithOptions(req, &ResponsesToChatOptions{
		ReasoningContentByID: func(string) string {
			lookupCalled = true
			return "cached thinking"
	placeholder,
placeholder)
placeholder
	require.Len(t, out.Messages, 3)
	require.Equal(t, "plain thinking", out.Messages[0].ReasoningContent)
	require.False(t, lookupCalled, "plaintext summary present → cache lookup must not run")
placeholder

// DeepSeek emits reasoning only once per turn; chained tool calls
// (reasoning → call A → output A → call B) have no reasoning item before call
// B. The turn's reasoning must be replayed on B's assistant message, otherwise
// DeepSeek thinking mode 400s the history ("reasoning_content ... must be
// passed back"). Reproduced from a real codex 0.147.0 resume history.
func TestResponsesToChat_ChainedToolCallsReplayTurnReasoning(t *testing.T) {
	req := &ResponsesRequest{
		Model: "deepseek-reasoner",
		Input: json.RawMessage(`[
			{"type":"reasoning","id":"item_r1","summary":[{"type":"summary_text","text":"turn thinking"placeholder]placeholder,
			{"type":"message","role":"assistant","content":[{"type":"output_text","text":"\n\n"placeholder]placeholder,
			{"type":"function_call","call_id":"call_a","name":"exec_command","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_a","output":"ok"placeholder,
			{"type":"function_call","call_id":"call_b","name":"exec_command","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_b","output":"ok"placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"next"placeholder]placeholder,
			{"type":"reasoning","id":"item_r2","summary":[{"type":"summary_text","text":"second turn"placeholder]placeholder,
			{"type":"function_call","call_id":"call_c","name":"exec_command","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_c","output":"ok"placeholder
		]`),
placeholder

	out, err := ResponsesToChatCompletionsRequest(req)
placeholder

	byCallID := map[string]ChatMessage{placeholder
	for _, m := range out.Messages {
		for _, tc := range m.ToolCalls {
			byCallID[tc.ID] = m
	placeholder
placeholder
	require.Len(t, byCallID, 3)
	require.Equal(t, "turn thinking", byCallID["call_a"].ReasoningContent)
	require.Equal(t, "turn thinking", byCallID["call_b"].ReasoningContent,
		"链式第二个工具调用必须回放本轮 reasoning")
	require.Equal(t, "second turn", byCallID["call_c"].ReasoningContent,
		"user 消息后开启新轮次，不得沿用上一轮 reasoning")

	// 每一条 assistant 消息都必须带 reasoning_content（DeepSeek 契约）。
	for i, m := range out.Messages {
		if m.Role == "assistant" {
			require.NotEmpty(t, m.ReasoningContent, "messages[%d] 缺 reasoning_content", i)
	placeholder
placeholder
placeholder

func TestExtractResponsesReasoningItem(t *testing.T) {
	id, text, ok := ExtractResponsesReasoningItem(json.RawMessage(
		`{"type":"reasoning","id":"item_a","summary":[{"type":"summary_text","text":"think"placeholder]placeholder`))
	require.True(t, ok)
	require.Equal(t, "item_a", id)
	require.Equal(t, "think", text)

	// Encrypted-only item: ok with id but empty text.
	id, text, ok = ExtractResponsesReasoningItem(json.RawMessage(
		`{"type":"reasoning","id":"item_b","summary":[],"encrypted_content":"opaque"placeholder`))
	require.True(t, ok)
	require.Equal(t, "item_b", id)
	require.Empty(t, text)

	// Non-reasoning items are skipped.
	_, _, ok = ExtractResponsesReasoningItem(json.RawMessage(
		`{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"placeholder]placeholder`))
	require.False(t, ok)

	_, _, ok = ExtractResponsesReasoningItem(json.RawMessage(`"bare string"`))
	require.False(t, ok)
placeholder
