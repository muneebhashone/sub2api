package service

import (
	"strings"
	"testing"
)

// TestOpenAIGatewayService_ToolCorrection 测试 OpenAIGatewayService 中的工具修正集成
func TestOpenAIGatewayService_ToolCorrection(t *testing.T) {
	// 创建一个简单的 service 实例来测试工具修正
	service := &OpenAIGatewayService{
		toolCorrector: NewCodexToolCorrector(),
placeholder

	tests := []struct {
		name     string
		input    []byte
		expected string
		changed  bool
placeholder{
		{
			name: "correct apply_patch in response body",
			input: []byte(`{
				"choices": [{
					"message": {
						"tool_calls": [{
							"function": {"name": "apply_patch"placeholder
					placeholder]
				placeholder
			placeholder]
		placeholder`),
			expected: "edit",
			changed:  true,
	placeholder,
		{
			name: "correct update_plan in response body",
			input: []byte(`{
				"tool_calls": [{
					"function": {"name": "update_plan"placeholder
			placeholder]
		placeholder`),
			expected: "todowrite",
			changed:  true,
	placeholder,
		{
			name: "no change for correct tool name",
			input: []byte(`{
				"tool_calls": [{
					"function": {"name": "edit"placeholder
			placeholder]
		placeholder`),
			expected: "edit",
			changed:  false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.correctToolCallsInResponseBody(tt.input)
			resultStr := string(result)

			// 检查是否包含期望的工具名称
			if !strings.Contains(resultStr, tt.expected) {
				t.Errorf("expected result to contain %q, got %q", tt.expected, resultStr)
		placeholder

			// 对于预期有变化的情况，验证结果与输入不同
			if tt.changed && string(result) == string(tt.input) {
				t.Error("expected result to be different from input, but they are the same")
		placeholder

			// 对于预期无变化的情况，验证结果与输入相同
			if !tt.changed && string(result) != string(tt.input) {
				t.Error("expected result to be same as input, but they are different")
		placeholder
	placeholder)
placeholder
placeholder

// TestOpenAIGatewayService_ToolCorrectorInitialization 测试工具修正器是否正确初始化
func TestOpenAIGatewayService_ToolCorrectorInitialization(t *testing.T) {
	service := &OpenAIGatewayService{
		toolCorrector: NewCodexToolCorrector(),
placeholder

	if service.toolCorrector == nil {
		t.Fatal("toolCorrector should not be nil")
placeholder

	// 测试修正器可以正常工作
	data := `{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholder`
	corrected, changed := service.toolCorrector.CorrectToolCallsInSSEData(data)

	if !changed {
		t.Error("expected tool call to be corrected")
placeholder

	if !strings.Contains(corrected, "edit") {
		t.Errorf("expected corrected data to contain 'edit', got %q", corrected)
placeholder
placeholder

// TestToolCorrectionStats 测试工具修正统计功能
func TestToolCorrectionStats(t *testing.T) {
	service := &OpenAIGatewayService{
		toolCorrector: NewCodexToolCorrector(),
placeholder

	// 执行几次修正
	testData := []string{
		`{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholder`,
		`{"tool_calls":[{"function":{"name":"update_plan"placeholderplaceholder]placeholder`,
		`{"tool_calls":[{"function":{"name":"apply_patch"placeholderplaceholder]placeholder`,
placeholder

	for _, data := range testData {
		service.toolCorrector.CorrectToolCallsInSSEData(data)
placeholder

	stats := service.toolCorrector.GetStats()

	if stats.TotalCorrected != 3 {
		t.Errorf("expected 3 corrections, got %d", stats.TotalCorrected)
placeholder

	if stats.CorrectionsByTool["apply_patch->edit"] != 2 {
		t.Errorf("expected 2 apply_patch->edit corrections, got %d", stats.CorrectionsByTool["apply_patch->edit"])
placeholder

	if stats.CorrectionsByTool["update_plan->todowrite"] != 1 {
		t.Errorf("expected 1 update_plan->todowrite correction, got %d", stats.CorrectionsByTool["update_plan->todowrite"])
placeholder
placeholder
