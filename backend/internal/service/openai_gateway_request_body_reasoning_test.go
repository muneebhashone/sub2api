package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrimOpenAIEncryptedReasoningItems_ContentNull(t *testing.T) {
	reqBody := map[string]any{
		"model": "grok-4.5",
		"input": []any{
			map[string]any{"type": "message", "role": "user", "content": "hi"placeholder,
			map[string]any{
				"type":              "reasoning",
				"summary":          []any{map[string]any{"type": "summary_text", "text": "thinking..."placeholderplaceholder,
				"content":          nil,
				"encrypted_content": nil,
		placeholder,
			map[string]any{"type": "message", "role": "assistant", "content": "Hello!"placeholder,
	placeholder,
placeholder

	changed := trimOpenAIEncryptedReasoningItems(reqBody)
	require.True(t, changed)

	input := reqBody["input"].([]any)
	require.Len(t, input, 3)

	reasoning := input[1].(map[string]any)
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

	input := reqBody["input"].([]any)
	require.Len(t, input, 1)

	reasoning := input[0].(map[string]any)
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

	input := reqBody["input"].([]any)
	reasoning := input[0].(map[string]any)
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
