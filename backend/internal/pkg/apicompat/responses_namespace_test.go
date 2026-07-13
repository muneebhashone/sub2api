package apicompat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlattenResponsesNamespaces_RewritesDeclarationHistoryAndChoice(t *testing.T) {
	req := map[string]any{
		"model": "gpt-5.5",
		"tools": []any{
			map[string]any{"type": "function", "name": "plain", "description": "keep"placeholder,
			map[string]any{
				"type": "namespace",
				"name": "collaboration",
				"tools": []any{
					map[string]any{"type": "function", "name": "spawn_agent", "description": "spawn", "parameters": map[string]any{"type": "object"placeholderplaceholder,
			placeholder,
		placeholder,
	placeholder,
		"tool_choice": map[string]any{"type": "function", "name": "spawn_agent", "namespace": "collaboration"placeholder,
		"input": []any{
			map[string]any{"type": "function_call", "call_id": "call_1", "name": "spawn_agent", "namespace": "collaboration", "arguments": "{placeholder"placeholder,
			map[string]any{"type": "message", "role": "user", "content": "hi", "name": "spawn_agent", "namespace": "collaboration"placeholder,
	placeholder,
placeholder

	names, changed, err := FlattenResponsesNamespaces(req)
placeholder
	require.True(t, changed)
	require.Equal(t, ResponsesNamespaceName{Namespace: "collaboration", Name: "spawn_agent"placeholder, names["collaboration__spawn_agent"])

	tools := req["tools"].([]any)
	require.Len(t, tools, 2)
	require.Equal(t, "plain", tools[0].(map[string]any)["name"])
	require.Equal(t, "collaboration__spawn_agent", tools[1].(map[string]any)["name"])
	require.Equal(t, "spawn", tools[1].(map[string]any)["description"])

	choice := req["tool_choice"].(map[string]any)
	require.Equal(t, "collaboration__spawn_agent", choice["name"])
	require.NotContains(t, choice, "namespace")

	call := req["input"].([]any)[0].(map[string]any)
	require.Equal(t, "collaboration__spawn_agent", call["name"])
	require.NotContains(t, call, "namespace")
	message := req["input"].([]any)[1].(map[string]any)
	require.Equal(t, "spawn_agent", message["name"])
	require.Equal(t, "collaboration", message["namespace"])
	require.Equal(t, "gpt-5.5", req["model"])
placeholder

func TestFlattenResponsesNamespaces_RejectsFlatNameCollision(t *testing.T) {
	req := map[string]any{"tools": []any{
		map[string]any{"type": "function", "name": "collaboration__spawn_agent"placeholder,
		map[string]any{"type": "namespace", "name": "collaboration", "tools": []any{
			map[string]any{"type": "function", "name": "spawn_agent"placeholder,
placeholder
placeholderplaceholder

	_, _, err := FlattenResponsesNamespaces(req)
	require.ErrorContains(t, err, "conflicts with a top-level tool")
placeholder

func TestFlattenResponsesNamespaces_NamespaceGroupChoiceFallsBackToAuto(t *testing.T) {
	req := map[string]any{
		"tools": []any{map[string]any{
			"type": "namespace", "name": "collaboration", "tools": []any{
				map[string]any{"type": "function", "name": "spawn_agent"placeholder,
				map[string]any{"type": "function", "name": "send_message"placeholder,
		placeholder,
placeholder
		"tool_choice": map[string]any{"type": "namespace", "name": "collaboration"placeholder,
placeholder

	_, changed, err := FlattenResponsesNamespaces(req)
placeholder
	require.True(t, changed)
	require.Equal(t, "auto", req["tool_choice"])
placeholder

func TestFlattenResponsesNamespacesExcept_PreservesBuiltInNamespaceAndChoice(t *testing.T) {
	req := map[string]any{
		"tools": []any{
			map[string]any{"type": "namespace", "name": "image_gen", "tools": []any{
				map[string]any{"type": "function", "name": "imagegen"placeholder,
	placeholder
			map[string]any{"type": "namespace", "name": "collaboration", "tools": []any{
				map[string]any{"type": "function", "name": "spawn_agent"placeholder,
	placeholder
	placeholder,
		"tool_choice": map[string]any{"type": "namespace", "name": "image_gen"placeholder,
placeholder

	names, changed, err := FlattenResponsesNamespacesExcept(req, map[string]bool{"image_gen": trueplaceholder)
placeholder
	require.True(t, changed)
	require.Contains(t, names, "collaboration__spawn_agent")
	tools := req["tools"].([]any)
	require.Equal(t, "namespace", tools[0].(map[string]any)["type"])
	require.Equal(t, "image_gen", tools[0].(map[string]any)["name"])
	require.Equal(t, "function", tools[1].(map[string]any)["type"])
	require.Equal(t, "collaboration__spawn_agent", tools[1].(map[string]any)["name"])
	require.Equal(t, map[string]any{"type": "namespace", "name": "image_gen"placeholder, req["tool_choice"])
placeholder

func TestFlattenResponsesNamespaces_RejectsNamespaceCollision(t *testing.T) {
	req := map[string]any{"tools": []any{
		map[string]any{"type": "namespace", "name": "a", "tools": []any{
			map[string]any{"type": "function", "name": "b__c"placeholder,
placeholder
		map[string]any{"type": "namespace", "name": "a__b", "tools": []any{
			map[string]any{"type": "function", "name": "c"placeholder,
placeholder
placeholderplaceholder

	_, _, err := FlattenResponsesNamespaces(req)
	require.ErrorContains(t, err, "both flatten")
placeholder

func TestRestoreResponsesNamespaceCalls_RewritesOnlyFunctionCalls(t *testing.T) {
	payload := []byte(`{"type":"response.completed","response":{"output":[{"type":"function_call","name":"collaboration__spawn_agent","call_id":"call_1","arguments":"{placeholder","extra":"keep"placeholder,{"type":"function_call","name":"plain","arguments":"{placeholder"placeholder,{"type":"message","name":"collaboration__spawn_agent","content":"<tag>&value</tag>"placeholder]placeholderplaceholder`)
	names := map[string]ResponsesNamespaceName{
		"collaboration__spawn_agent": {Namespace: "collaboration", Name: "spawn_agent"placeholder,
placeholder

	got, changed, err := RestoreResponsesNamespaceCalls(payload, names)
placeholder
	require.True(t, changed)
	require.JSONEq(t, `{"type":"response.completed","response":{"output":[{"type":"function_call","name":"spawn_agent","namespace":"collaboration","call_id":"call_1","arguments":"{placeholder","extra":"keep"placeholder,{"type":"function_call","name":"plain","arguments":"{placeholder"placeholder,{"type":"message","name":"collaboration__spawn_agent","content":"<tag>&value</tag>"placeholder]placeholderplaceholder`, string(got))
	require.Contains(t, string(got), "<tag>&value</tag>")
	require.NotContains(t, string(got), `\u003c`)
placeholder

func TestRestoreResponsesNamespaceCalls_RewritesLifecycleItems(t *testing.T) {
	for _, eventType := range []string{"response.output_item.added", "response.output_item.done"placeholder {
		t.Run(eventType, func(t *testing.T) {
			payload := []byte(`{"type":"` + eventType + `","item":{"type":"function_call","name":"collaboration__spawn_agent","arguments":"{placeholder"placeholderplaceholder`)
			got, changed, err := RestoreResponsesNamespaceCalls(payload, map[string]ResponsesNamespaceName{
				"collaboration__spawn_agent": {Namespace: "collaboration", Name: "spawn_agent"placeholder,
		placeholder)
		placeholder
			require.True(t, changed)
			require.JSONEq(t, `{"type":"`+eventType+`","item":{"type":"function_call","name":"spawn_agent","namespace":"collaboration","arguments":"{placeholder"placeholderplaceholder`, string(got))
	placeholder)
placeholder
placeholder
