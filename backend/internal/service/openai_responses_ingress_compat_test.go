package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestNormalizeOpenAIResponsesLegacyIngressMessagesOnly(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","stream":false,"messages":[{"role":"system","content":"repo policy"placeholder,{"role":"assistant","tool_calls":[{"id":"call_1","type":"function","function":{"name":"lookup","arguments":"{\"q\":\"hello\"placeholder"placeholderplaceholder]placeholder,{"role":"tool","tool_call_id":"call_1","content":"ok"placeholder],"previous_response_id":"resp_stale"placeholder`)

	normalized, changed, err := normalizeOpenAIResponsesLegacyIngress(body)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(normalized, "messages").Exists())
	require.False(t, gjson.GetBytes(normalized, "previous_response_id").Exists())
	require.Equal(t, "system", gjson.GetBytes(normalized, "input.0.role").String())
	require.Equal(t, "repo policy", gjson.GetBytes(normalized, "input.0.content").String())
	require.Equal(t, "function_call", gjson.GetBytes(normalized, "input.1.type").String())
	require.Equal(t, "call_1", gjson.GetBytes(normalized, "input.1.call_id").String())
	require.Equal(t, "function_call_output", gjson.GetBytes(normalized, "input.2.type").String())
	require.Equal(t, "call_1", gjson.GetBytes(normalized, "input.2.call_id").String())
	require.Equal(t, gjson.False, gjson.GetBytes(normalized, "stream").Type)
placeholder

func TestNormalizeOpenAIResponsesLegacyIngressConvertsChatTopLevelFields(t *testing.T) {
	body := []byte(`{"model":"gpt-4.1","messages":[{"role":"user","content":"weather?"placeholder],"tools":[{"type":"function","function":{"name":"lookup","description":"Lookup weather","parameters":{"type":"object"placeholderplaceholderplaceholder],"tool_choice":{"type":"function","function":{"name":"lookup"placeholderplaceholder,"max_tokens":64,"reasoning_effort":"high","service_tier":"flex","temperature":0.2,"top_p":0.7placeholder`)

	normalized, changed, err := normalizeOpenAIResponsesLegacyIngress(body)
placeholder
	require.True(t, changed)
	require.Equal(t, "user", gjson.GetBytes(normalized, "input.0.role").String())
	require.Equal(t, "function", gjson.GetBytes(normalized, "tools.0.type").String())
	require.Equal(t, "lookup", gjson.GetBytes(normalized, "tools.0.name").String())
	require.False(t, gjson.GetBytes(normalized, "tools.0.function").Exists())
	require.Equal(t, "function", gjson.GetBytes(normalized, "tool_choice.type").String())
	require.Equal(t, "lookup", gjson.GetBytes(normalized, "tool_choice.name").String())
	require.False(t, gjson.GetBytes(normalized, "tool_choice.function").Exists())
	require.Equal(t, int64(128), gjson.GetBytes(normalized, "max_output_tokens").Int())
	require.Equal(t, "high", gjson.GetBytes(normalized, "reasoning.effort").String())
	require.Equal(t, "auto", gjson.GetBytes(normalized, "reasoning.summary").String())
	require.Equal(t, "flex", gjson.GetBytes(normalized, "service_tier").String())
	require.False(t, gjson.GetBytes(normalized, "max_tokens").Exists())
	require.False(t, gjson.GetBytes(normalized, "reasoning_effort").Exists())
placeholder

func TestNormalizeOpenAIResponsesLegacyIngressNativeInputRemainsAuthoritative(t *testing.T) {
	body := []byte(`{"model":"gpt-5.4","input":[{"id":"msg_native","type":"message","role":"user","content":[{"type":"input_text","text":"keep me"placeholder]placeholder],"messages":[{"role":"user","content":"do not merge"placeholder],"tools":[{"type":"function","function":{"name":"native_shape"placeholderplaceholder],"tool_choice":{"type":"function","function":{"name":"native_shape"placeholderplaceholder,"previous_response_id":"resp_keep"placeholder`)

	normalized, changed, err := normalizeOpenAIResponsesLegacyIngress(body)
placeholder
	require.True(t, changed)
	require.False(t, gjson.GetBytes(normalized, "messages").Exists())
	require.Equal(t, "msg_native", gjson.GetBytes(normalized, "input.0.id").String())
	require.Equal(t, int64(1), gjson.GetBytes(normalized, "input.#").Int())
	require.Equal(t, "native_shape", gjson.GetBytes(normalized, "tools.0.function.name").String())
	require.False(t, gjson.GetBytes(normalized, "tools.0.name").Exists())
	require.Equal(t, "native_shape", gjson.GetBytes(normalized, "tool_choice.function.name").String())
	require.Equal(t, "resp_keep", gjson.GetBytes(normalized, "previous_response_id").String())
placeholder

func TestNormalizeOpenAIResponsesLegacyIngressKeepsPromptAliasAndDropsCommands(t *testing.T) {
	body := []byte(`{"model":"gpt-4.1","prompt":"hello","commands":[{"name":"legacy"placeholder]placeholder`)

	normalized, changed, err := normalizeOpenAIResponsesLegacyIngress(body)
placeholder
	require.True(t, changed)
	require.Equal(t, "hello", gjson.GetBytes(normalized, "input").String())
	require.False(t, gjson.GetBytes(normalized, "prompt").Exists())
	require.False(t, gjson.GetBytes(normalized, "commands").Exists())
placeholder
