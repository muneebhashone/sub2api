package service

import (
	"encoding/json"
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
		{name: "tool_search_output", body: map[string]any{"input": []any{map[string]any{"type": "tool_search_output"placeholderplaceholderplaceholder, want: trueplaceholder,
		{name: "custom_tool_call_output", body: map[string]any{"input": []any{map[string]any{"type": "custom_tool_call_output"placeholderplaceholderplaceholder, want: trueplaceholder,
		{name: "mcp_tool_call_output", body: map[string]any{"input": []any{map[string]any{"type": "mcp_tool_call_output"placeholderplaceholderplaceholder, want: trueplaceholder,
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
	// 所有 Codex 工具输出都应视为续链输出，避免 WS 续链时丢失 previous_response_id。
	require.False(t, HasFunctionCallOutput(nil))
	for _, typ := range []string{
		"function_call_output",
		"tool_search_output",
		"custom_tool_call_output",
		"mcp_tool_call_output",
placeholder {
		require.True(t, HasFunctionCallOutput(map[string]any{
			"input": []any{map[string]any{"type": typplaceholderplaceholder,
	placeholder), typ)
placeholder
	require.False(t, HasFunctionCallOutput(map[string]any{
		"input": "text",
placeholder))
placeholder

func TestHasToolCallContext(t *testing.T) {
	// 工具调用上下文必须包含 call_id，才能作为可关联上下文。
	require.False(t, HasToolCallContext(nil))
	for _, typ := range []string{
		"tool_call",
		"function_call",
		"local_shell_call",
		"tool_search_call",
		"custom_tool_call",
		"mcp_tool_call",
placeholder {
		require.True(t, HasToolCallContext(map[string]any{
			"input": []any{map[string]any{"type": typ, "call_id": "call_1"placeholderplaceholder,
	placeholder), typ)
placeholder
	require.False(t, HasToolCallContext(map[string]any{
		"input": []any{map[string]any{"type": "tool_call"placeholderplaceholder,
placeholder))
placeholder

func TestFunctionCallOutputCallIDs(t *testing.T) {
	// 仅提取工具输出的非空 call_id，去重后返回。
	require.Empty(t, FunctionCallOutputCallIDs(nil))
	callIDs := FunctionCallOutputCallIDs(map[string]any{
		"input": []any{
			map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
			map[string]any{"type": "tool_search_output", "call_id": "call_search"placeholder,
			map[string]any{"type": "custom_tool_call_output", "call_id": "call_custom"placeholder,
			map[string]any{"type": "mcp_tool_call_output", "call_id": "call_mcp"placeholder,
			map[string]any{"type": "function_call_output", "call_id": ""placeholder,
			map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
	placeholder,
placeholder)
	require.ElementsMatch(t, []string{"call_1", "call_search", "call_custom", "call_mcp"placeholder, callIDs)
placeholder

func TestHasFunctionCallOutputMissingCallID(t *testing.T) {
	require.False(t, HasFunctionCallOutputMissingCallID(nil))
	require.True(t, HasFunctionCallOutputMissingCallID(map[string]any{
		"input": []any{map[string]any{"type": "function_call_output"placeholderplaceholder,
placeholder))
	require.True(t, HasFunctionCallOutputMissingCallID(map[string]any{
		"input": []any{map[string]any{"type": "tool_search_output"placeholderplaceholder,
placeholder))
	require.False(t, HasFunctionCallOutputMissingCallID(map[string]any{
		"input": []any{map[string]any{"type": "tool_search_output", "call_id": "call_1"placeholderplaceholder,
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

func TestValidateFunctionCallOutputContextBytesMatchesMapValidation(t *testing.T) {
	// handler 预校验走 raw JSON 扫描，语义必须与 service 内部 map 校验保持一致。
	cases := []struct {
		name string
		body map[string]any
placeholder{
		{
			name: "no_input",
			body: map[string]any{"model": "gpt-5.4"placeholder,
	placeholder,
		{
			name: "missing_call_id",
			body: map[string]any{"input": []any{map[string]any{"type": "function_call_output"placeholderplaceholderplaceholder,
	placeholder,
		{
			name: "call_id_without_reference",
			body: map[string]any{"input": []any{map[string]any{"type": "function_call_output", "call_id": "call_1"placeholderplaceholderplaceholder,
	placeholder,
		{
			name: "matching_reference",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
				map[string]any{"type": "item_reference", "id": "call_1"placeholder,
	placeholder
	placeholder,
		{
			name: "partial_reference",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
				map[string]any{"type": "tool_search_output", "call_id": "call_2"placeholder,
				map[string]any{"type": "item_reference", "id": "call_1"placeholder,
	placeholder
	placeholder,
		{
			name: "tool_context",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call_output", "call_id": "call_1"placeholder,
				map[string]any{"type": "function_call", "call_id": "call_1"placeholder,
	placeholder
	placeholder,
		{
			name: "all_codex_tool_outputs",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call_output", "call_id": "call_function"placeholder,
				map[string]any{"type": "tool_search_output", "call_id": "call_search"placeholder,
				map[string]any{"type": "custom_tool_call_output", "call_id": "call_custom"placeholder,
				map[string]any{"type": "mcp_tool_call_output", "call_id": "call_mcp"placeholder,
				map[string]any{"type": "item_reference", "id": "call_function"placeholder,
				map[string]any{"type": "item_reference", "id": "call_search"placeholder,
				map[string]any{"type": "item_reference", "id": "call_custom"placeholder,
				map[string]any{"type": "item_reference", "id": "call_mcp"placeholder,
	placeholder
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.body)
		placeholder

			require.Equal(t, ValidateFunctionCallOutputContext(tt.body), ValidateFunctionCallOutputContextBytes(bodyBytes))
	placeholder)
placeholder
placeholder

func TestAnalyzeToolCallOutputContextCoverageBytes(t *testing.T) {
	cases := []struct {
		name         string
		body         map[string]any
		hasOutput    bool
		coversAllIDs bool
placeholder{
		{
			name:         "no_input",
			body:         map[string]any{"model": "gpt-5.1"placeholder,
			hasOutput:    false,
			coversAllIDs: false,
	placeholder,
		{
			name: "no_tool_output",
			body: map[string]any{"input": []any{
				map[string]any{"type": "message", "content": "hi"placeholder,
	placeholder
			hasOutput:    false,
			coversAllIDs: false,
	placeholder,
		{
			name: "object_tool_output_requires_context_replay",
			body: map[string]any{"input": map[string]any{
				"type": "custom_tool_call_output", "call_id": "call_a",
	placeholder
			hasOutput:    true,
			coversAllIDs: false,
	placeholder,
		{
			name: "all_outputs_covered_by_context",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call", "call_id": "call_a"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_a"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: true,
	placeholder,
		{
			name: "all_outputs_covered_by_item_reference",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call_output", "call_id": "call_a"placeholder,
				map[string]any{"type": "item_reference", "id": "call_a"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: true,
	placeholder,
		{
			// 关键回归用例：input 内存在某一个上下文项，但另一个输出的 call_id
			// 只能由上游会话链（previous_response_id）解析——不可剥离。
			name: "partial_coverage_not_movable",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call", "call_id": "call_a"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_a"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_b"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: false,
	placeholder,
		{
			name: "unrelated_context_does_not_cover",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call", "call_id": "call_x"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_b"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: false,
	placeholder,
		{
			name: "output_missing_call_id_not_movable",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call", "call_id": "call_a"placeholder,
				map[string]any{"type": "function_call_output"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_a"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: false,
	placeholder,
		{
			name: "mixed_context_and_reference_cover_all",
			body: map[string]any{"input": []any{
				map[string]any{"type": "function_call", "call_id": "call_a"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_a"placeholder,
				map[string]any{"type": "function_call_output", "call_id": "call_b"placeholder,
				map[string]any{"type": "item_reference", "id": "call_b"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: true,
	placeholder,
		{
			name: "all_codex_output_types_covered",
			body: map[string]any{"input": []any{
				map[string]any{"type": "tool_search_output", "call_id": "call_s"placeholder,
				map[string]any{"type": "tool_search_call", "call_id": "call_s"placeholder,
				map[string]any{"type": "mcp_tool_call_output", "call_id": "call_m"placeholder,
				map[string]any{"type": "mcp_tool_call", "call_id": "call_m"placeholder,
	placeholder
			hasOutput:    true,
			coversAllIDs: true,
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(tt.body)
		placeholder

			coverage := AnalyzeToolCallOutputContextCoverageBytes(bodyBytes)
			require.Equal(t, tt.hasOutput, coverage.HasFunctionCallOutput, "HasFunctionCallOutput")
			require.Equal(t, tt.coversAllIDs, coverage.ContextCoversAllCallIDs, "ContextCoversAllCallIDs")
	placeholder)
placeholder
placeholder
