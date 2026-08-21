package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestTrimOpenAIEncryptedReasoningItems_ContentNull(t *testing.T) {
	reqBody := map[string]any{
		"model": "grok-4.5",
		"input": []any{
			map[string]any{"type": "message", "role": "user", "content": "hi"placeholder,
			map[string]any{
				"type":              "reasoning",
				"summary":           []any{map[string]any{"type": "summary_text", "text": "thinking..."placeholderplaceholder,
				"content":           nil,
				"encrypted_content": nil,
		placeholder,
			map[string]any{"type": "message", "role": "assistant", "content": "Hello!"placeholder,
	placeholder,
placeholder

	changed := trimOpenAIEncryptedReasoningItems(reqBody)
	require.True(t, changed)

	input, ok := reqBody["input"].([]any)
	require.True(t, ok)
	require.Len(t, input, 3)

	reasoning, ok := input[1].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "reasoning", reasoning["type"])
	assert.NotNil(t, reasoning["summary"])
	_, hasContent := reasoning["content"]
	assert.False(t, hasContent, "content: null should be stripped")
	_, hasEncrypted := reasoning["encrypted_content"]
	assert.False(t, hasEncrypted, "encrypted_content should be stripped")
placeholder

func TestTrimOpenAIEncryptedReasoningItems_ContentNullOnly(t *testing.T) {
	reqBody := map[string]any{
		"model": "grok-4.5",
		"input": []any{
			map[string]any{
				"type":    "reasoning",
				"summary": []any{map[string]any{"type": "summary_text", "text": "ok"placeholderplaceholder,
				"content": nil,
		placeholder,
	placeholder,
placeholder

	changed := trimOpenAIEncryptedReasoningItems(reqBody)
	require.True(t, changed)

	input, ok := reqBody["input"].([]any)
	require.True(t, ok)
	require.Len(t, input, 1)

	reasoning, ok := input[0].(map[string]any)
	require.True(t, ok)
	_, hasContent := reasoning["content"]
	assert.False(t, hasContent, "content: null should be stripped even without encrypted_content")
placeholder

func TestTrimOpenAIEncryptedReasoningItems_ContentNonNull(t *testing.T) {
	reqBody := map[string]any{
		"model": "grok-4.5",
		"input": []any{
			map[string]any{
				"type":    "reasoning",
				"summary": []any{map[string]any{"type": "summary_text", "text": "ok"placeholderplaceholder,
				"content": "some actual content",
		placeholder,
	placeholder,
placeholder

	changed := trimOpenAIEncryptedReasoningItems(reqBody)
	assert.False(t, changed, "non-null content should not be stripped")

	input, ok := reqBody["input"].([]any)
	require.True(t, ok)
	reasoning, ok := input[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "some actual content", reasoning["content"])
placeholder

func TestTrimOpenAIEncryptedReasoningItems_NoReasoningItems(t *testing.T) {
	reqBody := map[string]any{
		"model": "grok-4.5",
		"input": []any{
			map[string]any{"type": "message", "role": "user", "content": "hi"placeholder,
	placeholder,
placeholder

	changed := trimOpenAIEncryptedReasoningItems(reqBody)
	assert.False(t, changed)
placeholder

func TestTrimOpenAIEncryptedReasoningItems_Compaction(t *testing.T) {
	tests := []struct {
		name      string
		itemType  string
		encrypted bool
		changed   bool
placeholder{
		{name: "compaction", itemType: "compaction", encrypted: true, changed: trueplaceholder,
		{name: "compaction summary", itemType: "compaction_summary", encrypted: true, changed: trueplaceholder,
		{name: "unencrypted compaction", itemType: "compaction", encrypted: false, changed: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := map[string]any{"type": tt.itemType, "id": "cmp_stale"placeholder
			if tt.encrypted {
				item["encrypted_content"] = "gAAA"
		placeholder
			reqBody := map[string]any{"input": []any{
				item,
				map[string]any{"type": "message", "content": "hi"placeholder,
		placeholderplaceholder

			changed := trimOpenAIEncryptedReasoningItems(reqBody)
			assert.Equal(t, tt.changed, changed)
			input, ok := reqBody["input"].([]any)
			require.True(t, ok)
			require.NotEmpty(t, input)
			first, ok := input[0].(map[string]any)
			require.True(t, ok)
			if tt.changed {
				require.Len(t, input, 1)
				assert.Equal(t, "message", first["type"])
				return
		placeholder
			require.Len(t, input, 2)
			assert.Equal(t, tt.itemType, first["type"])
	placeholder)
placeholder
placeholder

func TestSanitizeOpenAICrossModeFailoverReasoning_DropsWholeEncryptedItem(t *testing.T) {
	body := []byte(`{"model":"gpt-5.1","input":[` +
		`{"type":"message","role":"user","content":"hi"placeholder,` +
		`{"type":"reasoning","id":"rs_kiro_1","encrypted_content":"ENC","summary":[{"type":"summary_text","text":"t"placeholder]placeholder,` +
		`{"type":"message","role":"assistant","content":"yo"placeholder` +
		`]placeholder`)

	sanitized, changed, err := SanitizeOpenAICrossModeFailoverReasoning(body)
placeholder
	require.True(t, changed)
	// The whole reasoning item is gone — id and summary go with encrypted_content,
	// unlike trimOpenAIEncryptedReasoningItems which keeps the skeleton.
	require.NotContains(t, string(sanitized), "reasoning")
	require.NotContains(t, string(sanitized), "rs_kiro_1")
	require.NotContains(t, string(sanitized), "summary_text")
	require.Equal(t, int64(2), gjson.GetBytes(sanitized, "input.#").Int())
placeholder

func TestSanitizeOpenAICrossModeFailoverReasoning_NoEncryptedIsNoop(t *testing.T) {
	body := []byte(`{"model":"gpt-5.1","input":[{"type":"reasoning","summary":[{"type":"summary_text","text":"t"placeholder]placeholder]placeholder`)
	sanitized, changed, err := SanitizeOpenAICrossModeFailoverReasoning(body)
placeholder
	require.False(t, changed, "reasoning without encrypted_content must be preserved")
	require.Equal(t, string(body), string(sanitized))
placeholder

func TestSanitizeOpenAICrossModeFailoverReasoning_NoInputIsNoop(t *testing.T) {
	body := []byte(`{"model":"gpt-5.1"placeholder`)
	sanitized, changed, err := SanitizeOpenAICrossModeFailoverReasoning(body)
placeholder
	require.False(t, changed)
	require.Equal(t, string(body), string(sanitized))
placeholder

func TestSanitizeOpenAICrossModeFailoverReasoning_PreservesLargeIntegers(t *testing.T) {
	body := []byte(`{"model":"gpt-5.1","input":[` +
		`{"type":"reasoning","id":"rs_kiro_1","encrypted_content":"ENC"placeholder,` +
		`{"type":"message","role":"user","content":"hi"placeholder` +
		`],"tools":[{"type":"function","function":{"name":"lookup","parameters":{"type":"object","properties":{"id":{"const":9007199254740993placeholderplaceholderplaceholderplaceholderplaceholder]placeholder`)

	sanitized, changed, err := SanitizeOpenAICrossModeFailoverReasoning(body)
placeholder
	require.True(t, changed)
	require.Contains(t, string(sanitized), `"const":9007199254740993`,
		"sanitization must not round JSON integers through float64")
placeholder

func TestTrimOpenAIEncryptedReasoningItems_ContentNullDropsBareSkeleton(t *testing.T) {
	reqBody := map[string]any{
		"input": []any{
			map[string]any{"type": "reasoning", "content": nilplaceholder,
	placeholder,
placeholder

	changed := trimOpenAIEncryptedReasoningItems(reqBody)
	require.True(t, changed)
	_, hasInput := reqBody["input"]
	assert.False(t, hasInput, "bare reasoning skeleton should be dropped, emptying input")
placeholder

func TestNormalizeOpenAIAPIKeyStoreFalseReasoningReplay(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","store":false,"input":[` +
		`{"type":"reasoning","id":"rs_encrypted","call_id":"remove","encrypted_content":"cipher","summary":null,"opaque":9007199254740993placeholder,` +
		`{"type":"reasoning","id":"rs_server_only","summary":[{"type":"summary_text","text":"drop"placeholder]placeholder,` +
		`{"type":"item_reference","id":"rs_server_only"placeholder,` +
		`{"type":"item_reference","id":"msg_keep"placeholder,` +
		`{"type":"message","id":"msg_keep","role":"user","content":"continue"placeholder` +
		`]placeholder`)

	normalized, changed, err := normalizeOpenAIAPIKeyStoreFalseReasoningReplay(body, false)
placeholder
	require.True(t, changed)
	require.Equal(t, int64(3), gjson.GetBytes(normalized, "input.#").Int())
	require.Equal(t, "reasoning", gjson.GetBytes(normalized, "input.0.type").String())
	require.False(t, gjson.GetBytes(normalized, "input.0.id").Exists())
	require.False(t, gjson.GetBytes(normalized, "input.0.call_id").Exists())
	require.Equal(t, "cipher", gjson.GetBytes(normalized, "input.0.encrypted_content").String())
	require.True(t, gjson.GetBytes(normalized, "input.0.summary").IsArray())
	require.Equal(t, "9007199254740993", gjson.GetBytes(normalized, "input.0.opaque").Raw)
	require.Equal(t, "msg_keep", gjson.GetBytes(normalized, "input.1.id").String())
	require.Equal(t, "message", gjson.GetBytes(normalized, "input.2.type").String())
placeholder

func TestNormalizeOpenAIAPIKeyStoreFalseReasoningReplayRequiresExplicitStoreFalse(t *testing.T) {
	for _, body := range []string{
		`{"input":[{"type":"reasoning","id":"rs_keep"placeholder]placeholder`,
		`{"store":true,"input":[{"type":"reasoning","id":"rs_keep"placeholder]placeholder`,
placeholder {
		normalized, changed, err := normalizeOpenAIAPIKeyStoreFalseReasoningReplay([]byte(body), false)
	placeholder
		require.False(t, changed)
		require.Equal(t, body, string(normalized))
placeholder
placeholder

func TestNormalizeOpenAIAPIKeyStoreFalseReasoningReplayKnownCompactMode(t *testing.T) {
	body := []byte(`{"input":[{"type":"reasoning","id":"rs_drop","summary":[]placeholder,{"type":"message","content":"continue"placeholder]placeholder`)

	normalized, changed, err := normalizeOpenAIAPIKeyStoreFalseReasoningReplay(body, true)

placeholder
	require.True(t, changed)
	require.Equal(t, int64(1), gjson.GetBytes(normalized, "input.#").Int())
	require.Equal(t, "message", gjson.GetBytes(normalized, "input.0.type").String())
placeholder

func TestNormalizeOpenAIAPIKeyStoreFalseReasoningReplayRejectsEmptyEncryptedContent(t *testing.T) {
	for _, encrypted := range []string{"null", `""`, `"   "`, "123"placeholder {
		body := []byte(`{"store":false,"input":[{"type":"reasoning","id":"rs_drop","encrypted_content":` + encrypted + `placeholder,{"type":"message","content":"continue"placeholder]placeholder`)
		normalized, changed, err := normalizeOpenAIAPIKeyStoreFalseReasoningReplay(body, false)
	placeholder
		require.True(t, changed)
		require.Equal(t, int64(1), gjson.GetBytes(normalized, "input.#").Int())
		require.Equal(t, "message", gjson.GetBytes(normalized, "input.0.type").String())
placeholder
placeholder

func TestNormalizeOpenAIParallelToolCallsWithoutTools(t *testing.T) {
	withTools := []byte(`{"tools":[{"type":"function","name":"lookup"placeholder],"parallel_tool_calls":falseplaceholder`)
	normalized, changed, err := normalizeOpenAIParallelToolCallsWithoutTools(withTools)
placeholder
	require.False(t, changed)
	require.Equal(t, string(withTools), string(normalized))

	withoutTools := []byte(`{"input":"hi","parallel_tool_calls":trueplaceholder`)
	normalized, changed, err = normalizeOpenAIParallelToolCallsWithoutTools(withoutTools)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(normalized, "parallel_tool_calls").Exists())
placeholder
