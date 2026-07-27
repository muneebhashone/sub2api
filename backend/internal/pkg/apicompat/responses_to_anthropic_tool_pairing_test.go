package apicompat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// assertAnthropicPairing enforces the Anthropic Messages tool-pairing invariants
// that, when violated, surface as upstream 400s.
func assertAnthropicPairing(t *testing.T, messages []AnthropicMessage) {
placeholder
	for i, m := range messages {
		blocks := parseContentBlocks(m.Content)

		// No two consecutive same-role messages.
		if i > 0 {
			require.NotEqualf(t, messages[i-1].Role, m.Role, "consecutive %s messages at %d", m.Role, i)
	placeholder

		for _, b := range blocks {
			switch b.Type {
			case "tool_result":
				// Must have a matching tool_use in the immediately previous message.
				require.Positivef(t, i, "tool_result %s has no previous message", b.ToolUseID)
				prev := parseContentBlocks(messages[i-1].Content)
				require.Truef(t, hasToolUse(prev, b.ToolUseID),
					"tool_result %s has no corresponding tool_use in previous message", b.ToolUseID)
			case "tool_use":
				// Must be answered by a tool_result in the immediately next message.
				require.Lessf(t, i+1, len(messages), "tool_use %s has no following message", b.ID)
				next := parseContentBlocks(messages[i+1].Content)
				require.Truef(t, hasToolResult(next, b.ID),
					"tool_use %s is not answered in the next message", b.ID)
		placeholder
	placeholder
placeholder
placeholder

func hasToolUse(blocks []AnthropicContentBlock, id string) bool {
	for _, b := range blocks {
		if b.Type == "tool_use" && b.ID == id {
			return true
	placeholder
placeholder
	return false
placeholder

func hasToolResult(blocks []AnthropicContentBlock, toolUseID string) bool {
	for _, b := range blocks {
		if b.Type == "tool_result" && b.ToolUseID == toolUseID {
			return true
	placeholder
placeholder
	return false
placeholder

func convertAnthropic(t *testing.T, input string) []AnthropicMessage {
placeholder
	_, messages, err := convertResponsesInputToAnthropic("", json.RawMessage(input))
placeholder
	assertAnthropicPairing(t, messages)
	return messages
placeholder

// Tests use call_-prefixed ids because fromResponsesCallIDToAnthropic passes
// those through unchanged (matching codex's real call_00_... ids); bare ids
// would be rewritten to toolu_<id>.

// A developer/approval message injected between a function_call and its output
// must be moved out of the tool_use→tool_result adjacency. This is the shape
// that produced the production 400 "tool_result ... must have a corresponding
// tool_use block in the previous message".
func TestAnthropicPairing_DeveloperMessageBetween(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"do it"placeholder]placeholder,
		{"type":"function_call","call_id":"call_A","name":"exec","arguments":"{placeholder"placeholder,
		{"type":"message","role":"developer","content":[{"type":"input_text","text":"Approved command prefix saved"placeholder]placeholder,
		{"type":"function_call_output","call_id":"call_A","output":"ok"placeholder
	]`)
	// The assistant tool_use message is immediately followed by its tool_result.
	for i, m := range msgs {
		if hasToolUse(parseContentBlocks(m.Content), "call_A") {
			require.Equal(t, "user", msgs[i+1].Role)
			require.True(t, hasToolResult(parseContentBlocks(msgs[i+1].Content), "call_A"))
	placeholder
placeholder
placeholder

// Parallel tool calls where both outputs arrive stay grouped: one assistant
// message with both tool_use blocks, the next user message with both results.
func TestAnthropicPairing_ParallelBothAnswered(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"features?"placeholder]placeholder,
		{"type":"function_call","call_id":"call_c0","name":"exec","arguments":"{placeholder"placeholder,
		{"type":"function_call","call_id":"call_c1","name":"exec","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_c0","output":"log"placeholder,
		{"type":"function_call_output","call_id":"call_c1","output":"tags"placeholder
	]`)
	var sawGrouped bool
	for _, m := range msgs {
		blocks := parseContentBlocks(m.Content)
		if hasToolUse(blocks, "call_c0") && hasToolUse(blocks, "call_c1") {
			sawGrouped = true
	placeholder
placeholder
	require.True(t, sawGrouped, "parallel tool_use blocks should share one assistant message")
placeholder

// A parallel call whose sibling output never arrived must be dropped so every
// remaining tool_use is answered.
func TestAnthropicPairing_ParallelOneUnanswered(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"q"placeholder]placeholder,
		{"type":"function_call","call_id":"call_A","name":"exec","arguments":"{placeholder"placeholder,
		{"type":"function_call","call_id":"call_B","name":"exec","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_A","output":"oa"placeholder
	]`)
	for _, m := range msgs {
		require.Falsef(t, hasToolUse(parseContentBlocks(m.Content), "call_B"),
			"unanswered tool_use call_B should have been dropped")
placeholder
placeholder

// An orphan tool_result whose tool_use was never announced must be dropped.
func TestAnthropicPairing_OrphanToolResultDropped(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"q"placeholder]placeholder,
		{"type":"function_call_output","call_id":"call_ghost","output":"orphan"placeholder
	]`)
	for _, m := range msgs {
		require.Falsef(t, hasToolResult(parseContentBlocks(m.Content), "call_ghost"),
			"orphan tool_result should have been dropped")
placeholder
placeholder

// A dangling tool_call at the end of the history (no output yet) drops the
// assistant message holding only that call, leaving no tool_use behind.
func TestAnthropicPairing_DanglingCallDropped(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"q"placeholder]placeholder,
		{"type":"function_call","call_id":"call_A","name":"exec","arguments":"{placeholder"placeholder
	]`)
	for _, m := range msgs {
		require.Falsef(t, hasToolUse(parseContentBlocks(m.Content), "call_A"),
			"dangling tool_use call_A should have been dropped")
placeholder
placeholder

// Baseline: a single answered call pairs correctly and preserves the surrounding
// turns.
func TestAnthropicPairing_SingleCall(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"message","role":"user","content":[{"type":"input_text","text":"latest sha?"placeholder]placeholder,
		{"type":"function_call","call_id":"call_A","name":"exec","arguments":"{\"cmd\":\"git rev-parse HEAD\"placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_A","output":"deadbeef"placeholder,
		{"type":"message","role":"assistant","content":[{"type":"output_text","text":"It is deadbeef."placeholder]placeholder
	]`)
	// user, assistant(tool_use), user(tool_result), assistant(text)
	require.GreaterOrEqual(t, len(msgs), 4)
	require.Equal(t, "user", msgs[0].Role)
	require.True(t, hasToolUse(parseContentBlocks(msgs[1].Content), "call_A"))
	require.True(t, hasToolResult(parseContentBlocks(msgs[2].Content), "call_A"))
placeholder

func TestResponsesToAnthropic_FunctionOutputContentArray(t *testing.T) {
	msgs := convertAnthropic(t, `[
		{"type":"function_call","call_id":"call_A","name":"view_image","arguments":"{placeholder"placeholder,
		{"type":"function_call_output","call_id":"call_A","output":[
			{"type":"input_text","text":"image loaded"placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,YQ=="placeholder
		]placeholder
	]`)

	require.Len(t, msgs, 2)
	resultBlocks := parseContentBlocks(msgs[1].Content)
	require.Len(t, resultBlocks, 1)
	require.Equal(t, "tool_result", resultBlocks[0].Type)

	var content []AnthropicContentBlock
	require.NoError(t, json.Unmarshal(resultBlocks[0].Content, &content))
	require.Len(t, content, 2)
	require.Equal(t, "text", content[0].Type)
	require.Equal(t, "image loaded", content[0].Text)
	require.Equal(t, "image", content[1].Type)
	require.NotNil(t, content[1].Source)
	require.Equal(t, "image/png", content[1].Source.MediaType)
	require.Equal(t, "YQ==", content[1].Source.Data)
placeholder
