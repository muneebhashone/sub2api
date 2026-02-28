package service

import (
	"encoding/json"
	"testing"
)

func TestMayContainToolCallPayload(t *testing.T) {
	if mayContainToolCallPayload([]byte(`{"type":"response.output_text.delta","delta":"hello"placeholder`)) {
		t.Fatalf("plain text event should not trigger tool-call parsing")
placeholder
	if !mayContainToolCallPayload([]byte(`{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholder`)) {
		t.Fatalf("tool_calls event should trigger tool-call parsing")
placeholder
placeholder

func TestCorrectToolCallsInSSEData(t *testing.T) {
	corrector := NewCodexToolCorrector()

	tests := []struct {
		name            string
		input           string
		expectCorrected bool
		checkFunc       func(t *testing.T, result string)
placeholder{
		{
			name:            "empty string",
			input:           "",
			expectCorrected: false,
	placeholder,
		{
			name:            "newline only",
			input:           "\n",
			expectCorrected: false,
	placeholder,
		{
			name:            "invalid json",
			input:           "not a json",
			expectCorrected: false,
	placeholder,
		{
			name:            "correct apply_patch in tool_calls",
			input:           `{"tool_calls":[{"function":{"name":"apply_patch","arguments":"{placeholder"placeholderplaceholder]placeholder`,
			expectCorrected: true,
			checkFunc: func(t *testing.T, result string) {
				var payload map[string]any
				if err := json.Unmarshal([]byte(result), &payload); err != nil {
					t.Fatalf("Failed to parse result: %v", err)
			placeholder
				toolCalls, ok := payload["tool_calls"].([]any)
				if !ok || len(toolCalls) == 0 {
					t.Fatal("No tool_calls found in result")
			placeholder
				toolCall, ok := toolCalls[0].(map[string]any)
				if !ok {
					t.Fatal("Invalid tool_call format")
			placeholder
				functionCall, ok := toolCall["function"].(map[string]any)
				if !ok {
					t.Fatal("Invalid function format")
			placeholder
				if functionCall["name"] != "edit" {
					t.Errorf("Expected tool name 'edit', got '%v'", functionCall["name"])
			placeholder
		placeholder,
	placeholder,
		{
			name:            "correct update_plan in function_call",
			input:           `{"function_call":{"name":"update_plan","arguments":"{placeholder"placeholderplaceholder`,
			expectCorrected: true,
			checkFunc: func(t *testing.T, result string) {
				var payload map[string]any
				if err := json.Unmarshal([]byte(result), &payload); err != nil {
					t.Fatalf("Failed to parse result: %v", err)
			placeholder
				functionCall, ok := payload["function_call"].(map[string]any)
				if !ok {
					t.Fatal("Invalid function_call format")
			placeholder
				if functionCall["name"] != "todowrite" {
					t.Errorf("Expected tool name 'todowrite', got '%v'", functionCall["name"])
			placeholder
		placeholder,
	placeholder,
		{
			name:            "correct search_files in delta.tool_calls",
			input:           `{"delta":{"tool_calls":[{"function":{"name":"search_files"placeholderplaceholder]placeholderplaceholder`,
			expectCorrected: true,
			checkFunc: func(t *testing.T, result string) {
				var payload map[string]any
				if err := json.Unmarshal([]byte(result), &payload); err != nil {
					t.Fatalf("Failed to parse result: %v", err)
			placeholder
				delta, ok := payload["delta"].(map[string]any)
				if !ok {
					t.Fatal("Invalid delta format")
			placeholder
				toolCalls, ok := delta["tool_calls"].([]any)
				if !ok || len(toolCalls) == 0 {
					t.Fatal("No tool_calls found in delta")
			placeholder
				toolCall, ok := toolCalls[0].(map[string]any)
				if !ok {
					t.Fatal("Invalid tool_call format")
			placeholder
				functionCall, ok := toolCall["function"].(map[string]any)
				if !ok {
					t.Fatal("Invalid function format")
			placeholder
				if functionCall["name"] != "grep" {
					t.Errorf("Expected tool name 'grep', got '%v'", functionCall["name"])
			placeholder
		placeholder,
	placeholder,
		{
			name:            "correct list_files in choices.message.tool_calls",
			input:           `{"choices":[{"message":{"tool_calls":[{"function":{"name":"list_files"placeholderplaceholder]placeholderplaceholder]placeholder`,
			expectCorrected: true,
			checkFunc: func(t *testing.T, result string) {
				var payload map[string]any
				if err := json.Unmarshal([]byte(result), &payload); err != nil {
					t.Fatalf("Failed to parse result: %v", err)
			placeholder
				choices, ok := payload["choices"].([]any)
				if !ok || len(choices) == 0 {
					t.Fatal("No choices found in result")
			placeholder
				choice, ok := choices[0].(map[string]any)
				if !ok {
					t.Fatal("Invalid choice format")
			placeholder
				message, ok := choice["message"].(map[string]any)
				if !ok {
					t.Fatal("Invalid message format")
			placeholder
				toolCalls, ok := message["tool_calls"].([]any)
				if !ok || len(toolCalls) == 0 {
					t.Fatal("No tool_calls found in message")
			placeholder
				toolCall, ok := toolCalls[0].(map[string]any)
				if !ok {
					t.Fatal("Invalid tool_call format")
			placeholder
				functionCall, ok := toolCall["function"].(map[string]any)
				if !ok {
					t.Fatal("Invalid function format")
			placeholder
				if functionCall["name"] != "glob" {
					t.Errorf("Expected tool name 'glob', got '%v'", functionCall["name"])
			placeholder
		placeholder,
	placeholder,
		{
			name:            "no correction needed",
			input:           `{"tool_calls":[{"function":{"name":"read","arguments":"{placeholder"placeholderplaceholder]placeholder`,
			expectCorrected: false,
	placeholder,
		{
			name:            "correct multiple tool calls",
			input:           `{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder,{"function":{"name":"read_file"placeholderplaceholder]placeholder`,
			expectCorrected: true,
			checkFunc: func(t *testing.T, result string) {
				var payload map[string]any
				if err := json.Unmarshal([]byte(result), &payload); err != nil {
					t.Fatalf("Failed to parse result: %v", err)
			placeholder
				toolCalls, ok := payload["tool_calls"].([]any)
				if !ok || len(toolCalls) < 2 {
					t.Fatal("Expected at least 2 tool_calls")
			placeholder

				toolCall1, ok := toolCalls[0].(map[string]any)
				if !ok {
					t.Fatal("Invalid first tool_call format")
			placeholder
				func1, ok := toolCall1["function"].(map[string]any)
				if !ok {
					t.Fatal("Invalid first function format")
			placeholder
				if func1["name"] != "edit" {
					t.Errorf("Expected first tool name 'edit', got '%v'", func1["name"])
			placeholder

				toolCall2, ok := toolCalls[1].(map[string]any)
				if !ok {
					t.Fatal("Invalid second tool_call format")
			placeholder
				func2, ok := toolCall2["function"].(map[string]any)
				if !ok {
					t.Fatal("Invalid second function format")
			placeholder
				if func2["name"] != "read" {
					t.Errorf("Expected second tool name 'read', got '%v'", func2["name"])
			placeholder
		placeholder,
	placeholder,
		{
			name:            "camelCase format - applyPatch",
			input:           `{"tool_calls":[{"function":{"name":"applyPatch"placeholderplaceholder]placeholder`,
			expectCorrected: true,
			checkFunc: func(t *testing.T, result string) {
				var payload map[string]any
				if err := json.Unmarshal([]byte(result), &payload); err != nil {
					t.Fatalf("Failed to parse result: %v", err)
			placeholder
				toolCalls, ok := payload["tool_calls"].([]any)
				if !ok || len(toolCalls) == 0 {
					t.Fatal("No tool_calls found in result")
			placeholder
				toolCall, ok := toolCalls[0].(map[string]any)
				if !ok {
					t.Fatal("Invalid tool_call format")
			placeholder
				functionCall, ok := toolCall["function"].(map[string]any)
				if !ok {
					t.Fatal("Invalid function format")
			placeholder
				if functionCall["name"] != "edit" {
					t.Errorf("Expected tool name 'edit', got '%v'", functionCall["name"])
			placeholder
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, corrected := corrector.CorrectToolCallsInSSEData(tt.input)

			if corrected != tt.expectCorrected {
				t.Errorf("Expected corrected=%v, got %v", tt.expectCorrected, corrected)
		placeholder

			if !corrected && result != tt.input {
				t.Errorf("Expected unchanged result when not corrected")
		placeholder

			if tt.checkFunc != nil {
				tt.checkFunc(t, result)
		placeholder
	placeholder)
placeholder
placeholder

func TestCorrectToolName(t *testing.T) {
	tests := []struct {
		input     string
		expected  string
		corrected bool
placeholder{
		{"apply_patch", "edit", trueplaceholder,
		{"applyPatch", "edit", trueplaceholder,
		{"update_plan", "todowrite", trueplaceholder,
		{"updatePlan", "todowrite", trueplaceholder,
		{"read_plan", "todoread", trueplaceholder,
		{"readPlan", "todoread", trueplaceholder,
		{"search_files", "grep", trueplaceholder,
		{"searchFiles", "grep", trueplaceholder,
		{"list_files", "glob", trueplaceholder,
		{"listFiles", "glob", trueplaceholder,
		{"read_file", "read", trueplaceholder,
		{"readFile", "read", trueplaceholder,
		{"write_file", "write", trueplaceholder,
		{"writeFile", "write", trueplaceholder,
		{"execute_bash", "bash", trueplaceholder,
		{"executeBash", "bash", trueplaceholder,
		{"exec_bash", "bash", trueplaceholder,
		{"execBash", "bash", trueplaceholder,
		{"unknown_tool", "unknown_tool", falseplaceholder,
		{"read", "read", falseplaceholder,
		{"edit", "edit", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, corrected := CorrectToolName(tt.input)

			if corrected != tt.corrected {
				t.Errorf("Expected corrected=%v, got %v", tt.corrected, corrected)
		placeholder

			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
		placeholder
	placeholder)
placeholder
placeholder

func TestGetToolNameMapping(t *testing.T) {
	mapping := GetToolNameMapping()

	expectedMappings := map[string]string{
		"apply_patch":  "edit",
		"update_plan":  "todowrite",
		"read_plan":    "todoread",
		"search_files": "grep",
		"list_files":   "glob",
placeholder

	for from, to := range expectedMappings {
		if mapping[from] != to {
			t.Errorf("Expected mapping[%s] = %s, got %s", from, to, mapping[from])
	placeholder
placeholder

	mapping["test_tool"] = "test_value"
	newMapping := GetToolNameMapping()
	if _, exists := newMapping["test_tool"]; exists {
		t.Error("Modifications to returned mapping should not affect original")
placeholder
placeholder

func TestCorrectorStats(t *testing.T) {
	corrector := NewCodexToolCorrector()

	stats := corrector.GetStats()
	if stats.TotalCorrected != 0 {
		t.Errorf("Expected TotalCorrected=0, got %d", stats.TotalCorrected)
placeholder
	if len(stats.CorrectionsByTool) != 0 {
		t.Errorf("Expected empty CorrectionsByTool, got length %d", len(stats.CorrectionsByTool))
placeholder

	corrector.CorrectToolCallsInSSEData(`{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholder`)
	corrector.CorrectToolCallsInSSEData(`{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholder`)
	corrector.CorrectToolCallsInSSEData(`{"tool_calls":[{"function":{"name":"update_plan"placeholderplaceholder]placeholder`)

	stats = corrector.GetStats()
	if stats.TotalCorrected != 3 {
		t.Errorf("Expected TotalCorrected=3, got %d", stats.TotalCorrected)
placeholder

	if stats.CorrectionsByTool["apply_patch->edit"] != 2 {
		t.Errorf("Expected apply_patch->edit count=2, got %d", stats.CorrectionsByTool["apply_patch->edit"])
placeholder

	if stats.CorrectionsByTool["update_plan->todowrite"] != 1 {
		t.Errorf("Expected update_plan->todowrite count=1, got %d", stats.CorrectionsByTool["update_plan->todowrite"])
placeholder

	corrector.ResetStats()
	stats = corrector.GetStats()
	if stats.TotalCorrected != 0 {
		t.Errorf("Expected TotalCorrected=0 after reset, got %d", stats.TotalCorrected)
placeholder
	if len(stats.CorrectionsByTool) != 0 {
		t.Errorf("Expected empty CorrectionsByTool after reset, got length %d", len(stats.CorrectionsByTool))
placeholder
placeholder

func TestComplexSSEData(t *testing.T) {
	corrector := NewCodexToolCorrector()

	input := `{
		"id": "chatcmpl-123",
		"object": "chat.completion.chunk",
		"created": 1234567890,
		"model": "gpt-5.1-codex",
		"choices": [
			{
				"index": 0,
				"delta": {
					"tool_calls": [
						{
							"index": 0,
							"function": {
								"name": "apply_patch",
								"arguments": "{\"file\":\"test.go\"placeholder"
						placeholder
					placeholder
					]
			placeholder,
				"finish_reason": null
		placeholder
		]
placeholder`

	result, corrected := corrector.CorrectToolCallsInSSEData(input)

	if !corrected {
		t.Error("Expected data to be corrected")
placeholder

	var payload map[string]any
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("Failed to parse result: %v", err)
placeholder

	choices, ok := payload["choices"].([]any)
	if !ok || len(choices) == 0 {
		t.Fatal("No choices found in result")
placeholder
	choice, ok := choices[0].(map[string]any)
	if !ok {
		t.Fatal("Invalid choice format")
placeholder
	delta, ok := choice["delta"].(map[string]any)
	if !ok {
		t.Fatal("Invalid delta format")
placeholder
	toolCalls, ok := delta["tool_calls"].([]any)
	if !ok || len(toolCalls) == 0 {
		t.Fatal("No tool_calls found in delta")
placeholder
	toolCall, ok := toolCalls[0].(map[string]any)
	if !ok {
		t.Fatal("Invalid tool_call format")
placeholder
	function, ok := toolCall["function"].(map[string]any)
	if !ok {
		t.Fatal("Invalid function format")
placeholder

	if function["name"] != "edit" {
		t.Errorf("Expected tool name 'edit', got '%v'", function["name"])
placeholder
placeholder

// TestCorrectToolParameters 测试工具参数修正
func TestCorrectToolParameters(t *testing.T) {
	corrector := NewCodexToolCorrector()

	tests := []struct {
		name     string
		input    string
		expected map[string]bool // key: 期待存在的参数, value: true表示应该存在
placeholder{
		{
			name: "rename work_dir to workdir in bash tool",
			input: `{
				"tool_calls": [{
					"function": {
						"name": "bash",
						"arguments": "{\"command\":\"ls\",\"work_dir\":\"/tmp\"placeholder"
				placeholder
			placeholder]
		placeholder`,
			expected: map[string]bool{
				"command":  true,
				"workdir":  true,
				"work_dir": false,
		placeholder,
	placeholder,
		{
			name: "rename snake_case edit params to camelCase",
			input: `{
				"tool_calls": [{
					"function": {
						"name": "apply_patch",
						"arguments": "{\"path\":\"/foo/bar.go\",\"old_string\":\"old\",\"new_string\":\"new\"placeholder"
				placeholder
			placeholder]
		placeholder`,
			expected: map[string]bool{
				"filePath":   true,
				"path":       false,
				"oldString":  true,
				"old_string": false,
				"newString":  true,
				"new_string": false,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corrected, changed := corrector.CorrectToolCallsInSSEData(tt.input)
			if !changed {
				t.Error("expected data to be corrected")
		placeholder

			// 解析修正后的数据
			var result map[string]any
			if err := json.Unmarshal([]byte(corrected), &result); err != nil {
				t.Fatalf("failed to parse corrected data: %v", err)
		placeholder

			// 检查工具调用
			toolCalls, ok := result["tool_calls"].([]any)
			if !ok || len(toolCalls) == 0 {
				t.Fatal("no tool_calls found in corrected data")
		placeholder

			toolCall, ok := toolCalls[0].(map[string]any)
			if !ok {
				t.Fatal("invalid tool_call structure")
		placeholder

			function, ok := toolCall["function"].(map[string]any)
			if !ok {
				t.Fatal("no function found in tool_call")
		placeholder

			argumentsStr, ok := function["arguments"].(string)
			if !ok {
				t.Fatal("arguments is not a string")
		placeholder

			var args map[string]any
			if err := json.Unmarshal([]byte(argumentsStr), &args); err != nil {
				t.Fatalf("failed to parse arguments: %v", err)
		placeholder

			// 验证期望的参数
			for param, shouldExist := range tt.expected {
				_, exists := args[param]
				if shouldExist && !exists {
					t.Errorf("expected parameter %q to exist, but it doesn't", param)
			placeholder
				if !shouldExist && exists {
					t.Errorf("expected parameter %q to not exist, but it does", param)
			placeholder
		placeholder
	placeholder)
placeholder
placeholder
