package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// 当数组末尾不是用户消息时（典型场景：Agent 工具循环结束于 tool/assistant），
// 应直接跳过审计——不再回溯查找历史中的某条用户消息。

func TestExtractContentModerationInput_AnthropicAgentToolLoopSkipsAudit(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"user","content":"调用一下天气工具"placeholder,
			{"role":"assistant","content":[{"type":"tool_use","id":"tool_1","name":"weather","input":{placeholderplaceholder]placeholder,
			{"role":"user","content":[{"type":"tool_result","tool_use_id":"tool_1","content":"晴 25 度"placeholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolAnthropicMessages, body)

	require.Empty(t, input.Text)
	require.Empty(t, input.Images)
placeholder

func TestExtractContentModerationInput_AnthropicFirstTurnExtractsUser(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"user","content":"Q1"placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolAnthropicMessages, body)

	require.Equal(t, "Q1", input.Text)
placeholder

func TestExtractContentModerationInput_AnthropicMultiTurnExtractsLatestUser(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"user","content":"Q1"placeholder,
			{"role":"assistant","content":"A1"placeholder,
			{"role":"user","content":"Q2"placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolAnthropicMessages, body)

	require.Equal(t, "Q2", input.Text)
placeholder

func TestExtractContentModerationInput_AnthropicStreamResendExtractsResend(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"user","content":"原问题"placeholder,
			{"role":"assistant","content":"部分回答……"placeholder,
			{"role":"user","content":"重发"placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolAnthropicMessages, body)

	require.Equal(t, "重发", input.Text)
placeholder

func TestExtractContentModerationInput_OpenAIChatAgentToolLoopSkipsAudit(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"system","content":"sys"placeholder,
			{"role":"user","content":"列出我的订单"placeholder,
			{"role":"assistant","content":null,"tool_calls":[{"id":"call_1","type":"function","function":{"name":"orders","arguments":"{placeholder"placeholderplaceholder]placeholder,
			{"role":"tool","tool_call_id":"call_1","content":"[]"placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIChat, body)

	require.Empty(t, input.Text)
	require.Empty(t, input.Images)
placeholder

func TestExtractContentModerationInput_OpenAIChatMultiTurnExtractsLatestUser(t *testing.T) {
	body := []byte(`{
		"messages": [
			{"role":"user","content":"Q1"placeholder,
			{"role":"assistant","content":"A1"placeholder,
			{"role":"user","content":"Q2"placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIChat, body)

	require.Equal(t, "Q2", input.Text)
placeholder

func TestExtractContentModerationInput_GeminiAgentToolLoopSkipsAudit(t *testing.T) {
	body := []byte(`{
		"contents": [
			{"role":"user","parts":[{"text":"查询天气"placeholder]placeholder,
			{"role":"model","parts":[{"functionCall":{"name":"weather","args":{placeholderplaceholderplaceholder]placeholder,
			{"role":"user","parts":[{"functionResponse":{"name":"weather","response":{"temp":25placeholderplaceholderplaceholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolGemini, body)

	require.Empty(t, input.Text)
	require.Empty(t, input.Images)
placeholder

func TestExtractContentModerationInput_GeminiFirstTurnExtractsUser(t *testing.T) {
	body := []byte(`{
		"contents": [
			{"role":"user","parts":[{"text":"你好"placeholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolGemini, body)

	require.Equal(t, "你好", input.Text)
placeholder

func TestExtractContentModerationInput_GeminiMultiTurnExtractsLatestUser(t *testing.T) {
	body := []byte(`{
		"contents": [
			{"role":"user","parts":[{"text":"Q1"placeholder]placeholder,
			{"role":"model","parts":[{"text":"A1"placeholder]placeholder,
			{"role":"user","parts":[{"text":"Q2"placeholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolGemini, body)

	require.Equal(t, "Q2", input.Text)
placeholder

func TestExtractContentModerationInput_ResponsesAgentToolLoopSkipsAudit(t *testing.T) {
	body := []byte(`{
		"input":[
			{"type":"message","role":"user","content":[{"type":"input_text","text":"运行测试"placeholder]placeholder,
			{"type":"function_call","call_id":"call_1","name":"run_tests","arguments":"{placeholder"placeholder,
			{"type":"function_call_output","call_id":"call_1","output":"all passed"placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIResponses, body)

	require.Empty(t, input.Text)
	require.Empty(t, input.Images)
placeholder

func TestExtractContentModerationInput_ResponsesLastUserMessageExtracted(t *testing.T) {
	body := []byte(`{
		"input":[
			{"type":"message","role":"user","content":[{"type":"input_text","text":"first"placeholder]placeholder,
			{"type":"message","role":"assistant","content":[{"type":"output_text","text":"answer"placeholder]placeholder,
			{"type":"message","role":"user","content":[{"type":"input_text","text":"latest"placeholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIResponses, body)

	require.Equal(t, "latest", input.Text)
placeholder

func TestExtractContentModerationInput_ResponsesLastIsAssistantSkipped(t *testing.T) {
	body := []byte(`{
		"input":[
			{"type":"message","role":"user","content":[{"type":"input_text","text":"q1"placeholder]placeholder,
			{"type":"message","role":"assistant","content":[{"type":"output_text","text":"a1"placeholder]placeholder
		]
placeholder`)

	input := ExtractContentModerationInput(ContentModerationProtocolOpenAIResponses, body)

	require.Empty(t, input.Text)
	require.Empty(t, input.Images)
placeholder
