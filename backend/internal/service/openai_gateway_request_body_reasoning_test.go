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
