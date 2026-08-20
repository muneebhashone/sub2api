package service

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestAliasOpenAIOAuthReservedToolNames_RewritesDeclarationsAndReferences(t *testing.T) {
	reqBody := map[string]any{
		"tools": []any{
			map[string]any{"type": "function", "name": "python"placeholder,
			map[string]any{"type": "namespace", "name": "code", "tools": []any{
				map[string]any{"type": "function", "name": "shell"placeholder,
	placeholder
	placeholder,
		"tool_choice": map[string]any{"type": "function", "name": "python"placeholder,
		"input": []any{
			map[string]any{"type": "function_call", "name": "python", "call_id": "fc_1"placeholder,
			map[string]any{"type": "additional_tools", "tools": []any{
				map[string]any{"type": "function", "function": map[string]any{"name": "python"placeholderplaceholder,
	placeholder
	placeholder,
placeholder

	reverse, changed, err := aliasOpenAIOAuthReservedToolNames(reqBody)
placeholder
	require.True(t, changed)
	require.Equal(t, "python", reverse[codexPythonToolAlias])
	tools, ok := reqBody["tools"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, tools)
	firstTool, ok := tools[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, firstTool["name"])
	toolChoice, ok := reqBody["tool_choice"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, toolChoice["name"])
	input, ok := reqBody["input"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, input)
	firstInput, ok := input[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, firstInput["name"])
	secondInput, ok := input[1].(map[string]any)
	require.True(t, ok)
	nestedTools, ok := secondInput["tools"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, nestedTools)
	nestedTool, ok := nestedTools[0].(map[string]any)
	require.True(t, ok)
	nestedFn, ok := nestedTool["function"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, nestedFn["name"])
placeholder

func TestAliasOpenAIOAuthReservedToolNames_CollisionDoesNotMutate(t *testing.T) {
	reqBody := map[string]any{"tools": []any{
		map[string]any{"type": "function", "name": "python"placeholder,
		map[string]any{"type": "function", "name": codexPythonToolAliasplaceholder,
placeholderplaceholder
	before, err := json.Marshal(reqBody)
placeholder

	reverse, changed, err := aliasOpenAIOAuthReservedToolNames(reqBody)
	require.ErrorContains(t, err, `both normalize to "python__sub2api"`)
	require.False(t, changed)
	require.Nil(t, reverse)
	after, marshalErr := json.Marshal(reqBody)
	require.NoError(t, marshalErr)
	require.JSONEq(t, string(before), string(after))
placeholder

func TestApplyCodexOAuthTransform_ReservedPythonNameIsOAuthOnly(t *testing.T) {
	reqBody := map[string]any{
		"model": "gpt-5.5",
		"tools": []any{map[string]any{"type": "function", "name": "PYTHON"placeholderplaceholder,
placeholder

	result := applyCodexOAuthTransform(reqBody, true, false)
	require.NoError(t, result.Error)
	require.Equal(t, "PYTHON", result.ToolNameReverse[codexPythonToolAlias])
	tools, ok := reqBody["tools"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, tools)
	tool, ok := tools[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, tool["name"])

	apiKeyBody := []byte(`{"type":"response.create","tools":[{"type":"function","name":"python"placeholder]placeholder`)
	normalized, changed, err := normalizeOpenAIResponsesWebSocketCompatibilityBody(apiKeyBody, &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder)
placeholder
	require.False(t, changed)
	require.JSONEq(t, string(apiKeyBody), string(normalized))
placeholder

func TestRestoreCodexToolNamesFromContext_HTTPAndWSPayloadShapes(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	setCodexToolNameReverse(c, map[string]string{codexPythonToolAlias: "python"placeholder)

	streamEvent := restoreCodexToolNamesFromContext(c, []byte(
		`{"type":"response.output_item.done","item":{"type":"function_call","name":"python__sub2api"placeholder,"note":"python__sub2api"placeholder`,
	))
	require.Equal(t, "python", gjson.GetBytes(streamEvent, "item.name").String())
	require.Equal(t, "python__sub2api", gjson.GetBytes(streamEvent, "note").String())

	nonStreaming := restoreCodexToolNamesFromContext(c, []byte(
		`{"id":"resp_1","output":[{"type":"function_call","name":"python__sub2api"placeholder]placeholder`,
	))
	require.Equal(t, "python", gjson.GetBytes(nonStreaming, "output.0.name").String())

	setCodexToolNameReverse(c, nil)
	require.JSONEq(t,
		`{"type":"response.output_item.added","item":{"name":"python__sub2api"placeholderplaceholder`,
		string(restoreCodexToolNamesFromContext(c, []byte(`{"type":"response.output_item.added","item":{"name":"python__sub2api"placeholderplaceholder`))),
	)
placeholder

func TestAliasOpenAIOAuthReservedToolNames_SessionUpdateOnlyTouchesFunctionProtocolNodes(t *testing.T) {
	body := []byte(`{"type":"session.update","session":{"tools":[{"type":"function","name":"python"placeholder,{"type":"image_generation","name":"python"placeholder,{"type":"namespace","name":"python","tools":[{"type":"function","name":"shell"placeholder]placeholder]placeholder,"metadata":{"name":"python"placeholder,"sequence":900719925474099312345placeholder`)

	aliased, reverse, changed, err := aliasOpenAIOAuthReservedToolNamesBody(body)
placeholder
	require.True(t, changed)
	require.Equal(t, "python", reverse[codexPythonToolAlias])
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(aliased, "session.tools.0.name").String())
	require.Equal(t, "python", gjson.GetBytes(aliased, "session.tools.1.name").String())
	require.Equal(t, "python", gjson.GetBytes(aliased, "session.tools.2.name").String())
	require.Equal(t, "shell", gjson.GetBytes(aliased, "session.tools.2.tools.0.name").String())
	require.Equal(t, "python", gjson.GetBytes(aliased, "metadata.name").String())
	require.Equal(t, "900719925474099312345", gjson.GetBytes(aliased, "sequence").Raw)
placeholder

func TestRestoreCodexToolNamesInJSON_OnlyTouchesResponseToolCallNodesAndPreservesNumbers(t *testing.T) {
	reverse := map[string]string{codexPythonToolAlias: "python"placeholder
	body := []byte(`{"type":"response.completed","response":{"output":[{"type":"function_call","name":"python__sub2api"placeholder,{"type":"message","name":"python__sub2api","content":[]placeholder]placeholder,"metadata":{"name":"python__sub2api"placeholder,"sequence":900719925474099312345placeholder`)

	restored := restoreCodexToolNamesInJSON(body, reverse)
	require.Equal(t, "python", gjson.GetBytes(restored, "response.output.0.name").String())
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(restored, "response.output.1.name").String())
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(restored, "metadata.name").String())
	require.Equal(t, "900719925474099312345", gjson.GetBytes(restored, "sequence").Raw)
placeholder

func TestRestoreCodexToolNamesInJSON_ExplicitHTTPAndSSEToolCallProtocols(t *testing.T) {
	reverse := map[string]string{codexPythonToolAlias: "python"placeholder
	tests := []struct {
		name string
		body string
		path string
placeholder{
		{
			name: "chat http",
			body: `{"choices":[{"message":{"tool_calls":[{"type":"function","function":{"name":"python__sub2api"placeholderplaceholder]placeholderplaceholder],"metadata":{"name":"python__sub2api"placeholderplaceholder`,
			path: "choices.0.message.tool_calls.0.function.name",
	placeholder,
		{
			name: "chat sse",
			body: `{"choices":[{"delta":{"tool_calls":[{"type":"function","function":{"name":"python__sub2api"placeholderplaceholder]placeholderplaceholder],"metadata":{"name":"python__sub2api"placeholderplaceholder`,
			path: "choices.0.delta.tool_calls.0.function.name",
	placeholder,
		{
			name: "messages tool use",
			body: `{"type":"content_block_start","content":[{"type":"tool_use","name":"python__sub2api"placeholder],"metadata":{"name":"python__sub2api"placeholderplaceholder`,
			path: "content.0.name",
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restored := restoreCodexToolNamesInJSON([]byte(tt.body), reverse)
			require.Equal(t, "python", gjson.GetBytes(restored, tt.path).String())
			require.Equal(t, codexPythonToolAlias, gjson.GetBytes(restored, "metadata.name").String())
	placeholder)
placeholder
placeholder

func TestAliasOpenAIOAuthReservedToolNames_PromptCompatibilityRunsFirst(t *testing.T) {
	body := []byte(`{"model":"gpt-5.5","prompt":[{"type":"function_call","name":"python","call_id":"fc_1"placeholder],"functions":[{"name":"python"placeholder],"function_call":{"name":"python"placeholder,"sequence":900719925474099312345placeholder`)
	reqBody, err := getOpenAIRequestBodyMap(nil, body)
placeholder
	result := applyCodexOAuthTransform(reqBody, true, false)
	require.NoError(t, result.Error)
	input, ok := reqBody["input"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, input)
	firstInput, ok := input[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, firstInput["name"])
	tools, ok := reqBody["tools"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, tools)
	tool, ok := tools[0].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, tool["name"])
	toolChoice, ok := reqBody["tool_choice"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, codexPythonToolAlias, toolChoice["name"])
	encoded, err := json.Marshal(reqBody)
placeholder
	require.Equal(t, "900719925474099312345", gjson.GetBytes(encoded, "sequence").Raw)
placeholder

func TestCodexToolNameReverse_WSSessionReplacementDoesNotChangeActiveTurn(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	setCodexToolNameReverse(c, nil)
	first := []byte(`{"type":"response.create","tools":[{"type":"function","name":"python"placeholder]placeholder`)
	updateCodexToolNameReverseForWSFrame(c, first, map[string]string{codexPythonToolAlias: "python"placeholder)

	update := []byte(`{"type":"session.update","session":{"tools":[{"type":"function","name":"python__sub2api"placeholder]placeholderplaceholder`)
	updateCodexToolNameReverseForWSFrame(c, update, nil)
	currentOutput := restoreCodexToolNamesFromContext(c, []byte(`{"type":"response.output_item.done","item":{"type":"function_call","name":"python__sub2api"placeholderplaceholder`))
	require.Equal(t, "python", gjson.GetBytes(currentOutput, "item.name").String())
	sessionEcho := restoreCodexToolNamesFromContext(c, []byte(`{"type":"session.updated","session":{"tools":[{"type":"function","name":"python__sub2api"placeholder]placeholderplaceholder`))
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(sessionEcho, "session.tools.0.name").String())

	next := []byte(`{"type":"response.create","input":"next"placeholder`)
	updateCodexToolNameReverseForWSFrame(c, next, nil)
	nextOutput := restoreCodexToolNamesFromContext(c, []byte(`{"type":"response.output_item.done","item":{"type":"function_call","name":"python__sub2api"placeholderplaceholder`))
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(nextOutput, "item.name").String())

	sessionPython := []byte(`{"type":"session.update","session":{"tools":[{"type":"function","name":"python"placeholder]placeholderplaceholder`)
	updateCodexToolNameReverseForWSFrame(c, sessionPython, map[string]string{codexPythonToolAlias: "python"placeholder)
	explicitLiteral := []byte(`{"type":"response.create","input":[{"type":"additional_tools","tools":[{"type":"function","name":"python__sub2api"placeholder]placeholder]placeholder`)
	updateCodexToolNameReverseForWSFrame(c, explicitLiteral, nil)
	literalOutput := restoreCodexToolNamesFromContext(c, []byte(`{"type":"response.output_item.done","item":{"type":"function_call","name":"python__sub2api"placeholderplaceholder`))
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(literalOutput, "item.name").String())
	updateCodexToolNameReverseForWSFrame(c, next, nil)
	inheritedOutput := restoreCodexToolNamesFromContext(c, []byte(`{"type":"response.output_item.done","item":{"type":"function_call","name":"python__sub2api"placeholderplaceholder`))
	require.Equal(t, "python", gjson.GetBytes(inheritedOutput, "item.name").String())
placeholder

func TestDecodeOpenAIJSONUseNumberRejectsTrailingDocument(t *testing.T) {
	var decoded map[string]any
	require.Error(t, decodeOpenAIJSONUseNumber([]byte(`{"name":"python"placeholder{"extra":trueplaceholder`), &decoded))
placeholder

func TestRestoreCodexToolNamesFromSSEContextUsesEventLineTypeWithoutAddingType(t *testing.T) {
	c, _ := gin.CreateTestContext(nil)
	setCodexToolNameReverse(c, map[string]string{codexPythonToolAlias: codexReservedPythonToolNameplaceholder)
	payload := []byte(`{"item":{"type":"function_call","name":"python__sub2api"placeholder,"metadata":{"name":"python__sub2api"placeholderplaceholder`)

	restored := restoreCodexToolNamesFromSSEContext(c, payload, "response.output_item.done")

	require.Equal(t, codexReservedPythonToolName, gjson.GetBytes(restored, "item.name").String())
	require.Equal(t, codexPythonToolAlias, gjson.GetBytes(restored, "metadata.name").String())
	require.False(t, gjson.GetBytes(restored, "type").Exists())
placeholder
