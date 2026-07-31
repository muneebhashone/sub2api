package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestNormalizeOpenAIPassthroughOAuthBody_RemovesUnsupportedUser(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":"hello","user":"user_123","metadata":{"user_id":"user_123"placeholder,"prompt_cache_retention":"24h","safety_identifier":"sid","stream_options":{"include_usage":trueplaceholderplaceholder`)

	normalized, changed, err := normalizeOpenAIPassthroughOAuthBody(body, false)
placeholder
	require.True(t, changed)
	for _, field := range openAIChatGPTInternalUnsupportedFields {
		require.False(t, gjson.GetBytes(normalized, field).Exists(), "%s should be stripped", field)
placeholder
	require.True(t, gjson.GetBytes(normalized, "stream").Bool())
	require.False(t, gjson.GetBytes(normalized, "store").Bool())
placeholder

func TestNormalizeOpenAIPassthroughOAuthBody_CompactRemovesUnsupportedUser(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":"hello","user":"user_123","metadata":{"user_id":"user_123"placeholder,"stream":true,"store":trueplaceholder`)

	normalized, changed, err := normalizeOpenAIPassthroughOAuthBody(body, true)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(normalized, "user").Exists())
	require.False(t, gjson.GetBytes(normalized, "metadata").Exists())
	require.False(t, gjson.GetBytes(normalized, "stream").Exists())
	require.False(t, gjson.GetBytes(normalized, "store").Exists())
placeholder

func TestNormalizeOpenAIPassthroughOAuthBody_StringInputWrappedAsArray(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":"hello world"placeholder`)

	normalized, changed, err := normalizeOpenAIPassthroughOAuthBody(body, false)
placeholder
	require.True(t, changed)

	input := gjson.GetBytes(normalized, "input")
	require.True(t, input.IsArray(), "string input should be converted to array")
	items := input.Array()
	require.Len(t, items, 1)
	require.Equal(t, "message", items[0].Get("type").String())
	require.Equal(t, "user", items[0].Get("role").String())
	require.Equal(t, "hello world", items[0].Get("content").String())
placeholder

func TestNormalizeOpenAIPassthroughOAuthBody_EmptyStringInputWrappedAsEmptyArray(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":"  "placeholder`)

	normalized, changed, err := normalizeOpenAIPassthroughOAuthBody(body, false)
placeholder
	require.True(t, changed)

	input := gjson.GetBytes(normalized, "input")
	require.True(t, input.IsArray())
	require.Len(t, input.Array(), 0, "whitespace-only input should become empty array")
placeholder

func TestNormalizeOpenAIPassthroughOAuthBody_ObjectInputWrappedAsArray(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":{"type":"message","role":"user","content":"hi"placeholderplaceholder`)

	normalized, changed, err := normalizeOpenAIPassthroughOAuthBody(body, false)
placeholder
	require.True(t, changed)

	input := gjson.GetBytes(normalized, "input")
	require.True(t, input.IsArray(), "object input should be wrapped in array")
	items := input.Array()
	require.Len(t, items, 1)
	require.Equal(t, "message", items[0].Get("type").String())
placeholder

func TestNormalizeOpenAIPassthroughOAuthBody_ArrayInputUnchanged(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":[{"type":"message","role":"user","content":"hi"placeholder]placeholder`)

	normalized, _, err := normalizeOpenAIPassthroughOAuthBody(body, false)
placeholder

	input := gjson.GetBytes(normalized, "input")
	require.True(t, input.IsArray())
	require.Len(t, input.Array(), 1)
	require.Equal(t, "message", input.Array()[0].Get("type").String())
placeholder

func TestDetectOpenAIPassthroughInstructionsRejectReason(t *testing.T) {
	for _, tt := range []struct {
		name string
		body string
		want string
placeholder{
		{name: "missing is optional", body: `{"model":"gpt-5.1-codex"placeholder`, want: ""placeholder,
		{name: "non string remains rejected", body: `{"instructions":{"text":"invalid"placeholderplaceholder`, want: "instructions_not_string"placeholder,
		{name: "empty remains rejected", body: `{"instructions":"  "placeholder`, want: "instructions_empty"placeholder,
		{name: "non empty remains accepted", body: `{"instructions":"client guidance"placeholder`, want: ""placeholder,
placeholder {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, detectOpenAIPassthroughInstructionsRejectReason("gpt-5.1-codex", []byte(tt.body)))
	placeholder)
placeholder
placeholder
