package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNeedsToolContinuationSignals(t *testing.T) {
	// 覆盖所有触发续链的信号来源，确保判定逻辑完整。
	cases := []struct {
		name string
		body map[string]any
		want bool
placeholder{
		{name: "nil", body: nil, want: falseplaceholder,
		{name: "previous_response_id", body: map[string]any{"previous_response_id": "resp_1"placeholder, want: trueplaceholder,
		{name: "previous_response_id_blank", body: map[string]any{"previous_response_id": "  "placeholder, want: falseplaceholder,
		{name: "function_call_output", body: map[string]any{"input": []any{map[string]any{"type": "function_call_output"placeholderplaceholderplaceholder, want: trueplaceholder,
		{name: "item_reference", body: map[string]any{"input": []any{map[string]any{"type": "item_reference"placeholderplaceholderplaceholder, want: trueplaceholder,
		{name: "tools", body: map[string]any{"tools": []any{map[string]any{"type": "function"placeholderplaceholderplaceholder, want: trueplaceholder,
		{name: "tools_empty", body: map[string]any{"tools": []any{placeholderplaceholder, want: falseplaceholder,
		{name: "tools_invalid", body: map[string]any{"tools": "bad"placeholder, want: falseplaceholder,
		{name: "tool_choice", body: map[string]any{"tool_choice": "auto"placeholder, want: trueplaceholder,
		{name: "tool_choice_object", body: map[string]any{"tool_choice": map[string]any{"type": "function"placeholderplaceholder, want: trueplaceholder,
		{name: "tool_choice_empty_object", body: map[string]any{"tool_choice": map[string]any{placeholderplaceholder, want: falseplaceholder,
		{name: "none", body: map[string]any{"input": []any{map[string]any{"type": "text", "text": "hi"placeholderplaceholderplaceholder, want: falseplaceholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NeedsToolContinuation(tt.body))
	placeholder)
placeholder
placeholder

func TestHasFunctionCallOutput(t *testing.T) {
	// 仅当 input 中存在 function_call_output 才视为续链输出。
	require.False(t, HasFunctionCallOutput(nil))
	require.True(t, HasFunctionCallOutput(map[string]any{
		"input": []any{map[string]any{"type": "function_call_output"placeholderplaceholder,
placeholder))
	require.False(t, HasFunctionCallOutput(map[string]any{
		"input": "text",
placeholder))
placeholder

func TestHasToolCallContext(t *testing.T) {
	// tool_call/function_call 必须包含 call_id，才能作为可关联上下文。
	require.False(t, HasToolCallContext(nil))
	require.True(t, HasToolCallContext(map[string]any{
		"input": []any{map[string]any{"type": "tool_call", "call_id": "call_1"placeholderplaceholder,
placeholder))
	require.True(t, HasToolCallContext(map[string]any{
		"input": []any{map[string]any{"type": "function_call", "call_id": "call_2"placeholderplaceholder,
placeholder))
	require.False(t, HasToolCallContext(map[string]any{
		"input": []any{map[string]any{"type": "tool_call"placeholderplaceholder,
placeholder))
placeholder

func TestFunctionCallOutputCallIDs(t *testing.T) {
	// 仅提取非空 call_id，去重后返回。
	require.Empty(t, FunctionCallOutputCallIDs(nil))
	callIDs := FunctionCallOutputCallIDs(map[string]any{
		"input": []any{
			map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
			map[string]any{"type": "function_call_output", "call_id": ""placeholder,
			map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
	placeholder,
placeholder)
	require.ElementsMatch(t, []string{"call_1"placeholder, callIDs)
placeholder

func TestHasFunctionCallOutputMissingCallID(t *testing.T) {
	require.False(t, HasFunctionCallOutputMissingCallID(nil))
	require.True(t, HasFunctionCallOutputMissingCallID(map[string]any{
		"input": []any{map[string]any{"type": "function_call_output"placeholderplaceholder,
placeholder))
	require.False(t, HasFunctionCallOutputMissingCallID(map[string]any{
		"input": []any{map[string]any{"type": "function_call_output", "call_id": "call_1"placeholderplaceholder,
placeholder))
placeholder

func TestHasItemReferenceForCallIDs(t *testing.T) {
	// item_reference 需要覆盖所有 call_id 才视为可关联上下文。
	require.False(t, HasItemReferenceForCallIDs(nil, []string{"call_1"placeholder))
	require.False(t, HasItemReferenceForCallIDs(map[string]any{placeholder, []string{"call_1"placeholder))
	req := map[string]any{
		"input": []any{
			map[string]any{"type": "item_reference", "id": "call_1"placeholder,
			map[string]any{"type": "item_reference", "id": "call_2"placeholder,
	placeholder,
placeholder
	require.True(t, HasItemReferenceForCallIDs(req, []string{"call_1"placeholder))
	require.True(t, HasItemReferenceForCallIDs(req, []string{"call_1", "call_2"placeholder))
	require.False(t, HasItemReferenceForCallIDs(req, []string{"call_1", "call_3"placeholder))
placeholder
