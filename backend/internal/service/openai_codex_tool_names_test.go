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
	require.Equal(t, codexPythonToolAlias, reqBody["tools"].([]any)[0].(map[string]any)["name"])
	require.Equal(t, codexPythonToolAlias, reqBody["tool_choice"].(map[string]any)["name"])
	require.Equal(t, codexPythonToolAlias, reqBody["input"].([]any)[0].(map[string]any)["name"])
	nested := reqBody["input"].([]any)[1].(map[string]any)["tools"].([]any)[0].(map[string]any)["function"].(map[string]any)
	require.Equal(t, codexPythonToolAlias, nested["name"])
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
	require.Equal(t, codexPythonToolAlias, reqBody["tools"].([]any)[0].(map[string]any)["name"])

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
