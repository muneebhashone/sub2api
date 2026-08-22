package apicompat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestAdaptResponsesClientTools_LowersDeclarationsHistoryChoiceAndNamespaces(t *testing.T) {
	req := map[string]any{
		"tools": []any{
			map[string]any{"type": "custom", "name": "exec", "format": map[string]any{"type": "grammar"placeholderplaceholder,
			map[string]any{"type": "tool_search"placeholder,
			map[string]any{"type": "namespace", "name": "team", "tools": []any{map[string]any{"type": "function", "name": "send"placeholderplaceholderplaceholder,
	placeholder,
		"tool_choice": map[string]any{"type": "custom", "name": "exec"placeholder,
		"input": []any{
			map[string]any{"type": "custom_tool_call", "id": "ctc_client", "call_id": "c1", "name": "exec", "input": "dir"placeholder,
			map[string]any{"type": "custom_tool_call_output", "id": "ctco_client", "call_id": "c1", "output": "ok"placeholder,
			map[string]any{"type": "tool_search_call", "id": "tsc_client", "call_id": "s1", "arguments": map[string]any{"query": "git"placeholderplaceholder,
			map[string]any{"type": "tool_search_output", "id": "tso_client", "call_id": "s1", "output": map[string]any{"groups": []string{"git"placeholderplaceholderplaceholder,
			map[string]any{"type": "function_call", "call_id": "n1", "namespace": "team", "name": "send", "arguments": "{placeholder"placeholder,
	placeholder,
placeholder

	mapping, changed, err := AdaptResponsesClientTools(req)
placeholder
	require.True(t, changed)
	require.True(t, mapping.CustomTools["exec"])
	require.True(t, mapping.ToolSearch)
	require.Equal(t, ResponsesNamespaceName{Namespace: "team", Name: "send"placeholder, mapping.NamespaceTools["team__send"])

	tools := requireResponsesClientToolValue[[]any](t, req["tools"])
	require.Len(t, tools, 3)
	exec := requireResponsesClientToolValue[map[string]any](t, tools[0])
	require.Equal(t, "function", exec["type"])
	parameters := requireResponsesClientToolValue[json.RawMessage](t, exec["parameters"])
	require.JSONEq(t, customToolInputSchema, string(parameters))
	search := requireResponsesClientToolValue[map[string]any](t, tools[1])
	require.Equal(t, toolSearchProxyName, search["name"])
	namespaceTool := requireResponsesClientToolValue[map[string]any](t, tools[2])
	require.Equal(t, "team__send", namespaceTool["name"])

	choice := requireResponsesClientToolValue[map[string]any](t, req["tool_choice"])
	require.Equal(t, "function", choice["type"])
	input := requireResponsesClientToolValue[[]any](t, req["input"])
	customCall := requireResponsesClientToolValue[map[string]any](t, input[0])
	require.Equal(t, "function_call", customCall["type"])
	require.NotContains(t, customCall, "id")
	require.JSONEq(t, `{"input":"dir"placeholder`, requireResponsesClientToolValue[string](t, customCall["arguments"]))
	customOutput := requireResponsesClientToolValue[map[string]any](t, input[1])
	require.Equal(t, "function_call_output", customOutput["type"])
	require.NotContains(t, customOutput, "id")
	searchCall := requireResponsesClientToolValue[map[string]any](t, input[2])
	require.Equal(t, "function_call", searchCall["type"])
	require.NotContains(t, searchCall, "id")
	require.Equal(t, toolSearchProxyName, searchCall["name"])
	require.JSONEq(t, `{"query":"git"placeholder`, requireResponsesClientToolValue[string](t, searchCall["arguments"]))
	searchOutput := requireResponsesClientToolValue[map[string]any](t, input[3])
	require.Equal(t, "function_call_output", searchOutput["type"])
	require.NotContains(t, searchOutput, "id")
	require.JSONEq(t, `{"groups":["git"]placeholder`, requireResponsesClientToolValue[string](t, searchOutput["output"]))
	namespaceCall := requireResponsesClientToolValue[map[string]any](t, input[4])
	require.Equal(t, "team__send", namespaceCall["name"])
placeholder

func TestAdaptResponsesClientTools_RemovesDeferredFlagsWhenToolSearchIsLowered(t *testing.T) {
	req := map[string]any{
		"tools": []any{
			map[string]any{"type": "tool_search"placeholder,
			map[string]any{"type": "function", "name": "shell", "defer_loading": trueplaceholder,
			map[string]any{"type": "function", "name": "apply_patch"placeholder,
	placeholder,
placeholder

	_, changed, err := AdaptResponsesClientTools(req)
placeholder
	require.True(t, changed)
	tools := requireResponsesClientToolValue[[]any](t, req["tools"])
	require.Equal(t, toolSearchProxyName, requireResponsesClientToolValue[map[string]any](t, tools[0])["name"])
	require.NotContains(t, requireResponsesClientToolValue[map[string]any](t, tools[1]), "defer_loading")
placeholder

func TestStripResponsesDeferredToolFlags_PreservesFlagsWithBuiltInToolSearch(t *testing.T) {
	tools := []any{
		map[string]any{"type": "tool_search"placeholder,
		map[string]any{"type": "function", "name": "shell", "defer_loading": trueplaceholder,
placeholder

	require.False(t, stripResponsesDeferredToolFlags(tools))
	require.Equal(t, true, requireResponsesClientToolValue[map[string]any](t, tools[1])["defer_loading"])
placeholder

func TestAdaptResponsesClientTools_LowersDiscoveredToolSearchOutput(t *testing.T) {
	requestJSON := `{
		"tools":[{"type":"tool_search"placeholder],
		"input":[
			{"type":"tool_search_call","id":"tsc_client","call_id":"call_search","arguments":{"query":"codex app"placeholder,"execution":"client","status":"completed"placeholder,
			{"type":"tool_search_output","id":"tso_client","call_id":"call_search","execution":"client","status":"completed","tools":[
				{"type":"namespace","name":"codex_app","tools":[{"type":"function","name":"load_workspace_dependencies","description":"Load workspace dependencies","parameters":{"type":"object","properties":{placeholder,"additionalProperties":falseplaceholderplaceholder]placeholder,
				{"type":"namespace","name":"multi_agent_v1","tools":[
					{"type":"function","name":"spawn_agent","description":"Spawn an agent","parameters":{"type":"object","properties":{"message":{"type":"string"placeholderplaceholder,"required":["message"],"additionalProperties":falseplaceholderplaceholder,
					{"type":"function","name":"wait_agent","description":"Wait for agents","parameters":{"type":"object","properties":{placeholder,"additionalProperties":falseplaceholderplaceholder
				]placeholder
			]placeholder
		]
placeholder`

	type adaptedRequest struct {
		req     map[string]any
		mapping ResponsesClientToolMapping
placeholder
	adapt := func() adaptedRequest {
		var req map[string]any
		require.NoError(t, json.Unmarshal([]byte(requestJSON), &req))
		mapping, changed, err := AdaptResponsesClientTools(req)
	placeholder
		require.True(t, changed)
		return adaptedRequest{req: req, mapping: mappingplaceholder
placeholder

	first := adapt()
	second := adapt()
	firstInput := requireResponsesClientToolValue[[]any](t, first.req["input"])
	secondInput := requireResponsesClientToolValue[[]any](t, second.req["input"])

	tools := requireResponsesClientToolValue[[]any](t, first.req["tools"])
	require.Len(t, tools, 4)
	require.Equal(t, []string{
		"tool_search",
		"codex_app__load_workspace_dependencies",
		"multi_agent_v1__spawn_agent",
		"multi_agent_v1__wait_agent",
placeholder, responsesClientToolNames(t, tools))
	require.Equal(t, ResponsesNamespaceName{Namespace: "multi_agent_v1", Name: "spawn_agent"placeholder, first.mapping.NamespaceTools["multi_agent_v1__spawn_agent"])
	require.Equal(t, ResponsesNamespaceName{Namespace: "multi_agent_v1", Name: "wait_agent"placeholder, first.mapping.NamespaceTools["multi_agent_v1__wait_agent"])

	call := requireResponsesClientToolValue[map[string]any](t, firstInput[0])
	require.Equal(t, "function_call", call["type"])
	require.Equal(t, toolSearchProxyName, call["name"])
	require.JSONEq(t, `{"query":"codex app"placeholder`, requireResponsesClientToolValue[string](t, call["arguments"]))
	require.NotContains(t, call, "execution")

	output := requireResponsesClientToolValue[map[string]any](t, firstInput[1])
	require.Equal(t, map[string]any{
		"type":    "function_call_output",
		"call_id": "call_search",
		"output":  output["output"],
placeholder, output)
	outputText := requireResponsesClientToolValue[string](t, output["output"])
	require.JSONEq(t, `[
		{"type":"namespace","name":"codex_app","tools":[{"type":"function","name":"load_workspace_dependencies","description":"Load workspace dependencies","parameters":{"type":"object","properties":{placeholder,"additionalProperties":falseplaceholderplaceholder]placeholder,
		{"type":"namespace","name":"multi_agent_v1","tools":[
			{"type":"function","name":"spawn_agent","description":"Spawn an agent","parameters":{"type":"object","properties":{"message":{"type":"string"placeholderplaceholder,"required":["message"],"additionalProperties":falseplaceholderplaceholder,
			{"type":"function","name":"wait_agent","description":"Wait for agents","parameters":{"type":"object","properties":{placeholder,"additionalProperties":falseplaceholderplaceholder
		]placeholder
	]`, outputText)
	secondOutput := requireResponsesClientToolValue[map[string]any](t, secondInput[1])
	require.Equal(t, outputText, secondOutput["output"], "tool discovery output encoding must be deterministic")

	restored, changed, err := RestoreResponsesClientToolPayload(
		[]byte(`{"output":[{"type":"function_call","name":"multi_agent_v1__spawn_agent","call_id":"call_spawn","arguments":"{\"message\":\"work\"placeholder"placeholder]placeholder`),
		first.mapping,
	)
placeholder
	require.True(t, changed)
	require.JSONEq(t, `{"output":[{"type":"function_call","name":"spawn_agent","namespace":"multi_agent_v1","call_id":"call_spawn","arguments":"{\"message\":\"work\"placeholder"placeholder]placeholder`, string(restored))
placeholder

func TestAdaptResponsesClientTools_PromotesDirectDiscoveryAndDeduplicatesIdenticalDeclarations(t *testing.T) {
	direct := map[string]any{
		"type": "function", "name": "inspect_result", "description": "Inspect a result",
		"parameters": map[string]any{"type": "object", "properties": map[string]any{placeholderplaceholder,
placeholder
	custom := map[string]any{
		"type": "custom", "name": "run_script", "description": "Run a script",
		"format": map[string]any{"type": "grammar"placeholder,
placeholder
	namespace := map[string]any{
		"type": "namespace", "name": "multi_agent_v1", "tools": []any{map[string]any{
			"type": "function", "name": "spawn_agent", "parameters": map[string]any{"type": "object"placeholder,
placeholder
placeholder
	req := map[string]any{
		"tools": []any{
			map[string]any{"type": "function", "name": "static_first", "parameters": map[string]any{"type": "object"placeholderplaceholder,
			map[string]any{"type": "tool_search"placeholder,
	placeholder,
		"input": []any{
			map[string]any{"type": "tool_search_output", "status": "completed", "call_id": "search_1", "tools": []any{direct, custom, namespaceplaceholderplaceholder,
			map[string]any{"type": "tool_search_output", "status": "completed", "call_id": "search_2", "tools": []any{copyClientTool(direct), copyClientTool(custom), copyClientTool(namespace)placeholderplaceholder,
	placeholder,
placeholder

	mapping, changed, err := AdaptResponsesClientTools(req)
placeholder
	require.True(t, changed)
	require.True(t, mapping.CustomTools["run_script"])
	require.Equal(t, ResponsesNamespaceName{Namespace: "multi_agent_v1", Name: "spawn_agent"placeholder, mapping.NamespaceTools["multi_agent_v1__spawn_agent"])
	tools := requireResponsesClientToolValue[[]any](t, req["tools"])
	require.Equal(t, []string{"static_first", "tool_search", "inspect_result", "run_script", "multi_agent_v1__spawn_agent"placeholder, responsesClientToolNames(t, tools))
	customTool := requireResponsesClientToolValue[map[string]any](t, tools[3])
	require.Equal(t, "function", customTool["type"])
	require.NotContains(t, customTool, "format")
	for _, raw := range requireResponsesClientToolValue[[]any](t, req["input"]) {
		item := requireResponsesClientToolValue[map[string]any](t, raw)
		require.Equal(t, "function_call_output", item["type"])
		require.NotContains(t, item, "tools")
		require.NotContains(t, item, "status")
placeholder
placeholder

func TestAdaptResponsesClientTools_RejectsDiscoveredSchemaAndNamespaceCollisions(t *testing.T) {
	tests := []struct {
		name        string
		staticTools []any
		discovered  []any
placeholder{
		{
			name: "direct schema collision",
			staticTools: []any{map[string]any{
				"type": "function", "name": "inspect", "parameters": map[string]any{"type": "object"placeholder,
	placeholder
			discovered: []any{map[string]any{
				"type": "function", "name": "inspect", "parameters": map[string]any{"type": "string"placeholder,
	placeholder
	placeholder,
		{
			name: "namespace schema collision",
			staticTools: []any{map[string]any{
				"type": "namespace", "name": "multi_agent_v1", "tools": []any{map[string]any{
					"type": "function", "name": "spawn_agent", "parameters": map[string]any{"type": "object"placeholder,
		placeholder
	placeholder
			discovered: []any{map[string]any{
				"type": "namespace", "name": "multi_agent_v1", "tools": []any{map[string]any{
					"type": "function", "name": "spawn_agent", "parameters": map[string]any{"type": "string"placeholder,
		placeholder
	placeholder
	placeholder,
		{
			name:        "flattened namespace collision",
			staticTools: []any{map[string]any{"type": "function", "name": "multi_agent_v1__spawn_agent"placeholderplaceholder,
			discovered: []any{map[string]any{
				"type": "namespace", "name": "multi_agent_v1", "tools": []any{map[string]any{"type": "function", "name": "spawn_agent"placeholderplaceholder,
	placeholder
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := map[string]any{
				"tools": append(tt.staticTools, map[string]any{"type": "tool_search"placeholder),
				"input": []any{map[string]any{
					"type": "tool_search_output", "status": "completed", "tools": tt.discovered,
		placeholder
		placeholder
			_, _, err := AdaptResponsesClientTools(req)
			require.ErrorContains(t, err, "conflicts")
	placeholder)
placeholder
placeholder

func TestAdaptResponsesClientTools_DoesNotPromoteUnusableDiscoveries(t *testing.T) {
	req := map[string]any{
		"tools": []any{map[string]any{"type": "tool_search"placeholderplaceholder,
		"input": []any{
			map[string]any{"type": "tool_search_output", "call_id": "search_in_progress", "status": "in_progress", "tools": []any{map[string]any{"type": "function", "name": "not_ready"placeholderplaceholderplaceholder,
			map[string]any{"type": "tool_search_output", "call_id": "search_malformed", "status": "completed", "tools": []any{map[string]any{"type": "function"placeholderplaceholderplaceholder,
	placeholder,
placeholder

	_, changed, err := AdaptResponsesClientTools(req)
placeholder
	require.True(t, changed, "the static tool_search declaration is still lowered")
	tools := requireResponsesClientToolValue[[]any](t, req["tools"])
	require.Equal(t, []string{"tool_search"placeholder, responsesClientToolNames(t, tools))
placeholder

func responsesClientToolNames(t *testing.T, tools []any) []string {
placeholder
	names := make([]string, 0, len(tools))
	for _, raw := range tools {
		tool := requireResponsesClientToolValue[map[string]any](t, raw)
		names = append(names, requireResponsesClientToolValue[string](t, tool["name"]))
placeholder
	return names
placeholder

func TestAdaptResponsesClientTools_ToolSearchOutputEdgeCases(t *testing.T) {
	unencodableOutput := make(chan struct{placeholder)
	tests := []struct {
		name             string
		item             map[string]any
		wantOutput       any
		wantOutputExists bool
		wantPrivateKeys  []string
		wantExactOutput  bool
		wantErr          bool
placeholder{
		{
			name:             "absent tools and output is rejected",
			item:             map[string]any{"type": "tool_search_output", "call_id": "call_empty", "status": "completed"placeholder,
			wantOutputExists: false,
			wantErr:          true,
	placeholder,
		{
			name: "preexisting string output wins",
			item: map[string]any{
				"type": "tool_search_output", "call_id": "call_legacy", "output": "legacy",
				"tools": []any{map[string]any{"type": "function", "name": "ignored"placeholderplaceholder, "execution": "client",
		placeholder,
			wantOutput:       "legacy",
			wantOutputExists: true,
			wantExactOutput:  true,
	placeholder,
		{
			name: "preexisting object output remains legacy representation",
			item: map[string]any{
				"type": "tool_search_output", "call_id": "call_object", "output": map[string]any{"groups": []any{"github"placeholderplaceholder,
				"tools": []any{map[string]any{"type": "function", "name": "ignored"placeholderplaceholder,
		placeholder,
			wantOutput:       `{"groups":["github"]placeholder`,
			wantOutputExists: true,
			wantExactOutput:  true,
	placeholder,
		{
			name: "unencodable preexisting output is rejected",
			item: map[string]any{
				"type": "tool_search_output", "call_id": "call_bad_output", "output": unencodableOutput,
				"tools": []any{map[string]any{"type": "function", "name": "retained"placeholderplaceholder, "status": "completed", "execution": "client",
		placeholder,
			wantOutput: unencodableOutput,
			wantErr:    true,
	placeholder,
		{
			name: "empty tools array is a valid empty output",
			item: map[string]any{
				"type": "tool_search_output", "call_id": "call_empty_tools",
				"tools": []any{placeholder, "status": "completed", "execution": "client",
		placeholder,
			wantOutput:       `[]`,
			wantOutputExists: true,
			wantExactOutput:  true,
	placeholder,
		{
			name: "non-array tools value is serialized directly",
			item: map[string]any{
				"type": "tool_search_output", "call_id": "call_malformed",
				"tools": map[string]any{"unexpected": trueplaceholder, "status": "completed", "execution": "client",
		placeholder,
			wantOutput:       `{"unexpected":trueplaceholder`,
			wantOutputExists: true,
			wantExactOutput:  true,
	placeholder,
		{
			name: "unencodable tools is rejected",
			item: map[string]any{
				"type": "tool_search_output", "call_id": "call_unencodable", "tools": make(chan struct{placeholder), "status": "completed",
		placeholder,
			wantErr: true,
	placeholder,
		{
			name: "missing call id is rejected",
			item: map[string]any{
				"type": "tool_search_output", "tools": []any{placeholder, "status": "completed",
		placeholder,
			wantErr: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := map[string]any{
				"tools": []any{map[string]any{"type": "tool_search"placeholderplaceholder,
				"input": []any{tt.itemplaceholder,
		placeholder
			_, changed, err := AdaptResponsesClientTools(req)
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
			require.True(t, changed)
			input := requireResponsesClientToolValue[[]any](t, req["input"])
			output := requireResponsesClientToolValue[map[string]any](t, input[0])
			require.Equal(t, "function_call_output", output["type"])
			actualOutput, outputExists := output["output"]
			require.Equal(t, tt.wantOutputExists, outputExists)
			if tt.wantOutputExists {
				require.Equal(t, tt.wantOutput, actualOutput)
		placeholder
			if tt.wantExactOutput {
				require.Equal(t, map[string]any{
					"type":    "function_call_output",
					"call_id": output["call_id"],
					"output":  tt.wantOutput,
			placeholder, output)
		placeholder
			if len(tt.wantPrivateKeys) > 0 {
				for _, key := range tt.wantPrivateKeys {
					require.Contains(t, output, key)
			placeholder
		placeholder else {
				require.NotContains(t, output, "tools")
				require.NotContains(t, output, "status")
				require.NotContains(t, output, "execution")
		placeholder
	placeholder)
placeholder
placeholder

func requireResponsesClientToolValue[T any](t *testing.T, value any) T {
placeholder
	typed, ok := value.(T)
	require.True(t, ok, "unexpected value type %T", value)
	return typed
placeholder

func TestAdaptResponsesClientTools_RejectsAmbiguousNames(t *testing.T) {
	cases := []map[string]any{
		{"tools": []any{map[string]any{"type": "custom", "name": "same"placeholder, map[string]any{"type": "function", "name": "same"placeholderplaceholderplaceholder,
		{"tools": []any{map[string]any{"type": "tool_search"placeholder, map[string]any{"type": "function", "name": "tool_search"placeholderplaceholderplaceholder,
		{"tools": []any{map[string]any{"type": "function", "name": "team__send"placeholder, map[string]any{"type": "namespace", "name": "team", "tools": []any{map[string]any{"type": "function", "name": "send"placeholderplaceholderplaceholderplaceholderplaceholder,
placeholder
	for _, req := range cases {
		_, _, err := AdaptResponsesClientTools(req)
	placeholder
placeholder
placeholder

func TestAdaptResponsesClientToolsWithInheritedMapping_LowersFollowupHistoryWithoutTools(t *testing.T) {
	req := map[string]any{
		"input": []any{
			map[string]any{
				"type": "custom_tool_call", "name": "exec",
				"call_id": "call_1", "input": "pwd",
		placeholder,
			map[string]any{
				"type": "custom_tool_call_output", "call_id": "call_1",
				"id":     "ctco_client_output_1",
				"output": []any{map[string]any{"type": "input_text", "text": "ok"placeholderplaceholder,
		placeholder,
	placeholder,
placeholder
	inherited := ResponsesClientToolMapping{CustomTools: map[string]bool{"exec": trueplaceholderplaceholder

	mapping, changed, err := AdaptResponsesClientToolsWithInheritedMapping(req, inherited)

placeholder
	require.True(t, changed)
	require.Equal(t, inherited, mapping)
	items := requireResponsesClientToolValue[[]any](t, req["input"])
	call := requireResponsesClientToolValue[map[string]any](t, items[0])
	require.Equal(t, "function_call", call["type"])
	require.JSONEq(t, `{"input":"pwd"placeholder`, requireResponsesClientToolValue[string](t, call["arguments"]))
	require.NotContains(t, call, "input")
	output := requireResponsesClientToolValue[map[string]any](t, items[1])
	require.Equal(t, "function_call_output", output["type"])
	require.NotContains(t, output, "id")
	require.JSONEq(t, `[{"text":"ok","type":"input_text"placeholder]`, requireResponsesClientToolValue[string](t, output["output"]))
placeholder

func TestAdaptResponsesClientToolsWithInheritedMapping_PromotesOmittedToolsDiscoveryIntoEffectiveDeclarations(t *testing.T) {
	req := map[string]any{
		"input": []any{map[string]any{
			"type": "tool_search_output", "call_id": "call_search", "status": "completed", "execution": "client",
			"tools": []any{map[string]any{
				"type": "namespace", "name": "multi_agent_v1", "tools": []any{map[string]any{
					"type": "function", "name": "spawn_agent", "parameters": map[string]any{"type": "object"placeholder,
		placeholder
	placeholder
placeholder
placeholder
	inherited := ResponsesClientToolMapping{
		ToolSearch: true,
		NamespaceTools: map[string]ResponsesNamespaceName{
			"codex_app__read_resource": {Namespace: "codex_app", Name: "read_resource"placeholder,
	placeholder,
placeholder
	lowered := []any{
		map[string]any{"type": "function", "name": "static_first", "parameters": map[string]any{"type": "object"placeholderplaceholder,
		map[string]any{"type": "function", "name": "tool_search", "parameters": json.RawMessage(toolSearchProxySchema)placeholder,
		map[string]any{"type": "function", "name": "codex_app__read_resource", "parameters": map[string]any{"type": "object"placeholderplaceholder,
placeholder

	mapping, changed, err := AdaptResponsesClientToolsWithInheritedMapping(req, inherited, lowered)
placeholder
	require.True(t, changed)
	require.True(t, mapping.ToolSearch)
	require.Equal(t, ResponsesNamespaceName{Namespace: "codex_app", Name: "read_resource"placeholder, mapping.NamespaceTools["codex_app__read_resource"])
	require.Equal(t, ResponsesNamespaceName{Namespace: "multi_agent_v1", Name: "spawn_agent"placeholder, mapping.NamespaceTools["multi_agent_v1__spawn_agent"])
	tools := requireResponsesClientToolValue[[]any](t, req["tools"])
	require.Equal(t, []string{
		"static_first", "tool_search", "codex_app__read_resource", "multi_agent_v1__spawn_agent",
placeholder, responsesClientToolNames(t, tools))
	output := requireResponsesClientToolValue[map[string]any](t, requireResponsesClientToolValue[[]any](t, req["input"])[0])
	require.Equal(t, "function_call_output", output["type"])
	require.IsType(t, "", output["output"])
	require.NotContains(t, output, "tools")
	require.NotContains(t, output, "status")
	require.NotContains(t, output, "execution")
placeholder

func TestAdaptResponsesClientToolsWithInheritedMapping_ExplicitToolsReplaceInheritedMapping(t *testing.T) {
	req := map[string]any{
		"tools": []any{placeholder,
		"input": []any{map[string]any{
			"type": "custom_tool_call", "name": "exec", "input": "pwd",
placeholder
placeholder

	mapping, changed, err := AdaptResponsesClientToolsWithInheritedMapping(
		req,
		ResponsesClientToolMapping{CustomTools: map[string]bool{"exec": trueplaceholderplaceholder,
	)

placeholder
	require.False(t, changed)
	require.Empty(t, mapping)
	items := requireResponsesClientToolValue[[]any](t, req["input"])
	call := requireResponsesClientToolValue[map[string]any](t, items[0])
	require.Equal(t, "custom_tool_call", call["type"])
placeholder

func TestAdaptResponsesClientToolsWithInheritedMapping_ExplicitToolResetDoesNotPromoteDiscovery(t *testing.T) {
	for _, reset := range []any{nil, []any{placeholderplaceholder {
		req := map[string]any{
			"tools": reset,
			"input": []any{map[string]any{
				"type": "tool_search_output", "call_id": "call_reset", "status": "completed",
				"tools": []any{map[string]any{"type": "function", "name": "must_not_promote"placeholderplaceholder,
	placeholder
	placeholder
		mapping, changed, err := AdaptResponsesClientToolsWithInheritedMapping(
			req,
			ResponsesClientToolMapping{ToolSearch: trueplaceholder,
			[]any{map[string]any{"type": "function", "name": "tool_search"placeholderplaceholder,
		)
	placeholder
		require.False(t, changed)
		require.Empty(t, mapping)
		item := requireResponsesClientToolValue[map[string]any](t, requireResponsesClientToolValue[[]any](t, req["input"])[0])
		require.Equal(t, "tool_search_output", item["type"])
placeholder
placeholder

func TestRestoreResponsesClientToolPayload_RestoresClientAndNamespaceCalls(t *testing.T) {
	mapping := ResponsesClientToolMapping{
		CustomTools: map[string]bool{"exec": trueplaceholder, ToolSearch: true,
		NamespaceTools: map[string]ResponsesNamespaceName{"team__send": {Namespace: "team", Name: "send"placeholderplaceholder,
placeholder
	payload := []byte(`{"id":"resp","output":[{"type":"function_call","id":"i1","call_id":"c1","name":"exec","arguments":"{\"input\":\"dir\"placeholder","namespace":"ignore"placeholder,{"type":"function_call","id":"i2","call_id":"s1","name":"tool_search","arguments":"{\"query\":\"git\"placeholder"placeholder,{"type":"function_call","id":"i3","call_id":"n1","name":"team__send","arguments":"{placeholder"placeholder]placeholder`)

	restored, changed, err := RestoreResponsesClientToolPayload(payload, mapping)
placeholder
	require.True(t, changed)
	require.JSONEq(t, `{"id":"resp","output":[{"type":"custom_tool_call","id":"i1","call_id":"c1","name":"exec","input":"dir"placeholder,{"type":"tool_search_call","id":"i2","call_id":"s1","execution":"client","arguments":{"query":"git"placeholderplaceholder,{"type":"function_call","id":"i3","call_id":"n1","name":"send","namespace":"team","arguments":"{placeholder"placeholder]placeholder`, string(restored))
placeholder

func TestResponsesClientToolStreamRestorer_CustomToolBuffersWrapperAndSequences(t *testing.T) {
	restorer := NewResponsesClientToolStreamRestorer(ResponsesClientToolMapping{CustomTools: map[string]bool{"exec": trueplaceholderplaceholder)
	added := restorer.Restore(ResponsesStreamEvent{Type: "response.output_item.added", SequenceNumber: 7, OutputIndex: 0, Item: &ResponsesOutput{Type: "function_call", ID: "i1", CallID: "c1", Name: "exec", Status: "in_progress"placeholderplaceholder)
	require.Len(t, added, 1)
	require.Equal(t, 7, added[0].SequenceNumber)
	require.Equal(t, "custom_tool_call", added[0].Item.Type)
	require.Empty(t, restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.delta", SequenceNumber: 8, ItemID: "i1", Delta: `{"input":"di`placeholder))
	done := restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.done", SequenceNumber: 9, ItemID: "i1", CallID: "c1", Name: "exec", Arguments: `{"input":"dir"placeholder`placeholder)
	require.Len(t, done, 2)
	require.Equal(t, 8, done[0].SequenceNumber)
	require.Equal(t, "response.custom_tool_call_input.delta", done[0].Type)
	require.Equal(t, "dir", done[0].Delta)
	require.Equal(t, 9, done[1].SequenceNumber)
	require.Equal(t, "response.custom_tool_call_input.done", done[1].Type)
	require.Equal(t, "dir", done[1].Input)
	closed := restorer.Restore(ResponsesStreamEvent{Type: "response.output_item.done", SequenceNumber: 10, OutputIndex: 0, Item: &ResponsesOutput{Type: "function_call", ID: "i1", CallID: "c1", Name: "exec", Arguments: `{"input":"dir"placeholder`, Status: "completed"placeholderplaceholder)
	require.Equal(t, 10, closed[0].SequenceNumber)
	require.Equal(t, "custom_tool_call", closed[0].Item.Type)
	require.Equal(t, "dir", closed[0].Item.Input)
placeholder

func TestResponsesClientToolStreamRestorer_ToolSearchAndFunction(t *testing.T) {
	restorer := NewResponsesClientToolStreamRestorer(ResponsesClientToolMapping{ToolSearch: trueplaceholder)
	search := restorer.Restore(ResponsesStreamEvent{Type: "response.output_item.added", SequenceNumber: 0, OutputIndex: 0, Item: &ResponsesOutput{Type: "function_call", ID: "s1", CallID: "c1", Name: "tool_search", Status: "in_progress"placeholderplaceholder)
	require.Equal(t, "tool_search_call", search[0].Item.Type)
	require.Empty(t, restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.delta", SequenceNumber: 1, ItemID: "s1", Delta: `{"query":"git"placeholder`placeholder))
	require.Empty(t, restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.done", SequenceNumber: 2, ItemID: "s1", Arguments: `{"query":"git"placeholder`placeholder))
	closed := restorer.Restore(ResponsesStreamEvent{Type: "response.output_item.done", SequenceNumber: 3, OutputIndex: 0, Item: &ResponsesOutput{Type: "function_call", ID: "s1", CallID: "c1", Name: "tool_search", Status: "completed"placeholderplaceholder)
	require.Equal(t, 1, closed[0].SequenceNumber)
	require.Equal(t, "tool_search_call", closed[0].Item.Type)
	require.JSONEq(t, `{"query":"git"placeholder`, string(toolSearchCallArgumentsJSON(closed[0].Item.Arguments)))

	function := restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.done", SequenceNumber: 4, ItemID: "plain", Name: "plain", Arguments: "{placeholder"placeholder)
	require.Len(t, function, 1)
	require.Equal(t, "response.function_call_arguments.done", function[0].Type)
	require.Equal(t, 2, function[0].SequenceNumber)
placeholder

func TestResponsesClientToolStreamRestorer_RestoresNamespaceLifecycle(t *testing.T) {
	restorer := NewResponsesClientToolStreamRestorer(ResponsesClientToolMapping{
		NamespaceTools: map[string]ResponsesNamespaceName{
			"browser__open": {Namespace: "browser", Name: "open"placeholder,
	placeholder,
placeholder)

	added, changed, err := restorer.RestoreEvent([]byte(`{"type":"response.output_item.added","sequence_number":4,"output_index":0,"item":{"type":"function_call","id":"i1","call_id":"c1","name":"browser__open","arguments":"","status":"in_progress"placeholderplaceholder`))
placeholder
	require.True(t, changed)
	require.Len(t, added, 1)
	require.Equal(t, "open", gjson.GetBytes(added[0], "item.name").String())
	require.Equal(t, "browser", gjson.GetBytes(added[0], "item.namespace").String())

	delta, changed, err := restorer.RestoreEvent([]byte(`{"type":"response.function_call_arguments.delta","sequence_number":5,"output_index":0,"item_id":"i1","name":"browser__open","delta":"{\"url\":"placeholder`))
placeholder
	require.True(t, changed)
	require.Len(t, delta, 1)
	require.Equal(t, "open", gjson.GetBytes(delta[0], "name").String())

	done, changed, err := restorer.RestoreEvent([]byte(`{"type":"response.function_call_arguments.done","sequence_number":6,"output_index":0,"item_id":"i1","name":"browser__open","arguments":"{placeholder"placeholder`))
placeholder
	require.True(t, changed)
	require.Len(t, done, 1)
	require.Equal(t, "open", gjson.GetBytes(done[0], "name").String())
placeholder

func TestResponsesClientToolStreamRestorer_RawEventsPreserveUnknownFieldsAndOutputFallback(t *testing.T) {
	restorer := NewResponsesClientToolStreamRestorer(ResponsesClientToolMapping{CustomTools: map[string]bool{"exec": trueplaceholderplaceholder)
	passthrough, changed, err := restorer.RestoreEvent([]byte(`{"type":"response.created","sequence_number":4,"response":{"id":"r"placeholder,"upstream_extension":{"keep":trueplaceholderplaceholder`))
placeholder
	require.False(t, changed)
	require.Len(t, passthrough, 1)
	require.Contains(t, string(passthrough[0]), `"upstream_extension":{"keep":trueplaceholder`)

	restorer.Restore(ResponsesStreamEvent{Type: "response.output_item.added", SequenceNumber: 5, OutputIndex: 9, Item: &ResponsesOutput{Type: "function_call", ID: "item", CallID: "call", Name: "exec"placeholderplaceholder)
	// Some upstreams omit every tool identity field on later argument chunks.
	require.Empty(t, restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.delta", SequenceNumber: 6, OutputIndex: 9, Delta: `{"input":"pwd"placeholder`placeholder))
	done := restorer.Restore(ResponsesStreamEvent{Type: "response.function_call_arguments.done", SequenceNumber: 7, OutputIndex: 9placeholder)
	require.Len(t, done, 2)
	require.Equal(t, "pwd", done[1].Input)
placeholder

func TestResponsesClientToolStreamRestorer_RestoresAllTerminalEvents(t *testing.T) {
	for _, eventType := range []string{
		"response.completed",
		"response.done",
		"response.incomplete",
		"response.failed",
		"response.cancelled",
		"response.canceled",
placeholder {
		t.Run(eventType, func(t *testing.T) {
			restorer := NewResponsesClientToolStreamRestorer(ResponsesClientToolMapping{CustomTools: map[string]bool{"exec": trueplaceholderplaceholder)
			payload := []byte(`{"type":"` + eventType + `","sequence_number":7,"response":{"id":"resp_tools","output":[{"type":"function_call","id":"item_exec","call_id":"call_exec","name":"exec","arguments":"{\"input\":\"pwd\"placeholder"placeholder]placeholderplaceholder`)

			restored, changed, err := restorer.RestoreEvent(payload)

		placeholder
			require.True(t, changed)
			require.Len(t, restored, 1)
			require.Equal(t, eventType, gjson.GetBytes(restored[0], "type").String())
			require.Equal(t, int64(7), gjson.GetBytes(restored[0], "sequence_number").Int())
			require.Equal(t, "custom_tool_call", gjson.GetBytes(restored[0], "response.output.0.type").String())
			require.Equal(t, "pwd", gjson.GetBytes(restored[0], "response.output.0.input").String())
			require.False(t, gjson.GetBytes(restored[0], "response.output.0.arguments").Exists())
	placeholder)
placeholder
placeholder
