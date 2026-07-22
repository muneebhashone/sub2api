package apicompat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// ResponsesClientToolMapping records the reversible lowering applied before a
// native Responses request is sent to an upstream that only understands
// function tools.
type ResponsesClientToolMapping struct {
	CustomTools    map[string]bool
	ToolSearch     bool
	NamespaceTools map[string]ResponsesNamespaceName
placeholder

// AdaptResponsesClientTools lowers Codex client-only tools in req to
// ordinary function tools. It mutates req and returns the mapping required to
// restore the upstream response.
func AdaptResponsesClientTools(req map[string]any) (ResponsesClientToolMapping, bool, error) {
	if req == nil {
		return ResponsesClientToolMapping{placeholder, false, nil
placeholder
	tools, ok := req["tools"].([]any)
	if !ok || len(tools) == 0 {
		return ResponsesClientToolMapping{placeholder, false, nil
placeholder

	adapter := ResponsesClientToolMapping{CustomTools: make(map[string]bool)placeholder
	functionNames := make(map[string]bool)
	customNames := make(map[string]bool)
	for _, raw := range tools {
		tool, ok := raw.(map[string]any)
		if !ok {
			continue
	placeholder
		name := strings.TrimSpace(stringValue(tool["name"]))
		switch strings.TrimSpace(stringValue(tool["type"])) {
		case "function":
			if name != "" {
				functionNames[name] = true
		placeholder
		case "custom":
			if name != "" {
				customNames[name] = true
		placeholder
		case "tool_search":
			adapter.ToolSearch = true
	placeholder
placeholder
	for name := range customNames {
		if functionNames[name] {
			return ResponsesClientToolMapping{placeholder, false, fmt.Errorf("custom tool %q conflicts with a function tool of the same name; this upstream cannot disambiguate them, rename one of the tools", name)
	placeholder
placeholder
	if adapter.ToolSearch && (functionNames[toolSearchProxyName] || customNames[toolSearchProxyName]) {
		return ResponsesClientToolMapping{placeholder, false, fmt.Errorf("built-in tool_search conflicts with a declared tool named %q; this upstream cannot disambiguate them, rename the tool", toolSearchProxyName)
placeholder

	// Namespace flattening also rewrites namespace-qualified history and choice.
	names, flattened, err := FlattenResponsesNamespaces(req)
	if err != nil {
		return ResponsesClientToolMapping{placeholder, false, err
placeholder
	adapter.NamespaceTools = names
	if adapter.ToolSearch {
		if _, exists := names[toolSearchProxyName]; exists {
			return ResponsesClientToolMapping{placeholder, false, fmt.Errorf("built-in tool_search conflicts with namespace tool flattened as %q; this upstream cannot disambiguate them, rename the tool", toolSearchProxyName)
	placeholder
placeholder

	tools, _ = req["tools"].([]any)
	lowered := make([]any, 0, len(tools))
	changed := flattened
	seenSearch := false
	for _, raw := range tools {
		tool, ok := raw.(map[string]any)
		if !ok {
			lowered = append(lowered, raw)
			continue
	placeholder
		typ := strings.TrimSpace(stringValue(tool["type"]))
		name := strings.TrimSpace(stringValue(tool["name"]))
		switch typ {
		case "custom":
			if name == "" {
				lowered = append(lowered, raw)
				continue
		placeholder
			copy := copyClientTool(tool)
			copy["type"] = "function"
			copy["parameters"] = json.RawMessage(customToolInputSchema)
			delete(copy, "format")
			adapter.CustomTools[name] = true
			lowered = append(lowered, copy)
			changed = true
		case "tool_search":
			if seenSearch {
				changed = true
				continue
		placeholder
			seenSearch = true
			lowered = append(lowered, map[string]any{
				"type": "function", "name": toolSearchProxyName,
				"description": "Search and load Codex tools, plugins, connectors, and MCP namespaces for the current task.",
				"parameters":  json.RawMessage(toolSearchProxySchema),
		placeholder)
			changed = true
		default:
			lowered = append(lowered, raw)
	placeholder
placeholder
	if changed {
		req["tools"] = lowered
placeholder
	if rewriteClientToolHistory(req["input"], &adapter) {
		changed = true
placeholder
	if rewriteClientToolChoice(req, &adapter) {
		changed = true
placeholder
	if len(adapter.CustomTools) == 0 {
		adapter.CustomTools = nil
placeholder
	if len(adapter.NamespaceTools) == 0 {
		adapter.NamespaceTools = nil
placeholder
	return adapter, changed, nil
placeholder

func copyClientTool(tool map[string]any) map[string]any {
	copy := make(map[string]any, len(tool))
	for key, value := range tool {
		copy[key] = value
placeholder
	return copy
placeholder

func rewriteClientToolHistory(value any, adapter *ResponsesClientToolMapping) bool {
	changed := false
	var visit func(any)
	visit = func(value any) {
		switch typed := value.(type) {
		case []any:
			for _, item := range typed {
				visit(item)
		placeholder
		case map[string]any:
			typ := strings.TrimSpace(stringValue(typed["type"]))
			switch typ {
			case "custom_tool_call":
				if adapter.CustomTools[strings.TrimSpace(stringValue(typed["name"]))] {
					typed["type"] = "function_call"
					typed["arguments"] = customToolCallArguments(stringValue(typed["input"]))
					delete(typed, "input")
					changed = true
			placeholder
			case "custom_tool_call_output":
				typed["type"] = "function_call_output"
				normalizeClientToolOutput(typed)
				changed = true
			case "tool_search_call":
				if adapter.ToolSearch {
					typed["type"] = "function_call"
					typed["name"] = toolSearchProxyName
					typed["arguments"] = rawObjectString(typed["arguments"])
					delete(typed, "execution")
					changed = true
			placeholder
			case "tool_search_output":
				if adapter.ToolSearch {
					typed["type"] = "function_call_output"
					normalizeClientToolOutput(typed)
					changed = true
			placeholder
		placeholder
			for _, child := range typed {
				visit(child)
		placeholder
	placeholder
placeholder
	visit(value)
	return changed
placeholder

func normalizeClientToolOutput(item map[string]any) {
	output, exists := item["output"]
	if !exists {
		return
placeholder
	if _, ok := output.(string); ok {
		return
placeholder
	if output == nil {
		item["output"] = ""
		return
placeholder
	encoded, err := json.Marshal(output)
	if err != nil {
		item["output"] = ""
		return
placeholder
	item["output"] = string(encoded)
placeholder

func rewriteClientToolChoice(req map[string]any, adapter *ResponsesClientToolMapping) bool {
	choice, ok := req["tool_choice"].(map[string]any)
	if !ok {
		return false
placeholder
	typ := strings.TrimSpace(stringValue(choice["type"]))
	name := strings.TrimSpace(stringValue(choice["name"]))
	if typ == "custom" && adapter.CustomTools[name] {
		choice["type"] = "function"
		return true
placeholder
	if typ == "tool_search" && adapter.ToolSearch {
		req["tool_choice"] = map[string]any{"type": "function", "name": toolSearchProxyNameplaceholder
		return true
placeholder
	return false
placeholder

func customToolCallArguments(input string) string {
	encoded, _ := json.Marshal(map[string]string{"input": inputplaceholder)
	return string(encoded)
placeholder

func rawObjectString(value any) string {
	if text, ok := value.(string); ok {
		return text
placeholder
	encoded, err := json.Marshal(value)
	if err != nil {
		return "{placeholder"
placeholder
	return string(encoded)
placeholder

// RestoreResponsesClientToolPayload restores client tool calls in a non-stream
// native Responses JSON payload.
func RestoreResponsesClientToolPayload(payload []byte, mapping ResponsesClientToolMapping) ([]byte, bool, error) {
	if len(payload) == 0 {
		return payload, false, nil
placeholder
	var value any
	if err := json.Unmarshal(payload, &value); err != nil {
		return payload, false, err
placeholder
	changed := restoreClientToolValue(value, &mapping)
	if !changed {
		if len(mapping.NamespaceTools) == 0 {
			return payload, false, nil
	placeholder
		return RestoreResponsesNamespaceCalls(payload, mapping.NamespaceTools)
placeholder
	var rebuilt bytes.Buffer
	encoder := json.NewEncoder(&rebuilt)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return payload, false, err
placeholder
	rebuiltPayload := bytes.TrimSuffix(rebuilt.Bytes(), []byte("\n"))
	if len(mapping.NamespaceTools) == 0 {
		return rebuiltPayload, true, nil
placeholder
	restored, _, err := RestoreResponsesNamespaceCalls(rebuiltPayload, mapping.NamespaceTools)
	if err != nil {
		return payload, false, err
placeholder
	return restored, true, nil
placeholder

func restoreClientToolValue(value any, adapter *ResponsesClientToolMapping) bool {
	changed := false
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			changed = restoreClientToolValue(item, adapter) || changed
	placeholder
	case map[string]any:
		if strings.TrimSpace(stringValue(typed["type"])) == "function_call" {
			name := strings.TrimSpace(stringValue(typed["name"]))
			if adapter.CustomTools[name] {
				typed["type"] = "custom_tool_call"
				typed["input"] = extractCustomToolCallInput(rawObjectString(typed["arguments"]))
				delete(typed, "arguments")
				delete(typed, "namespace")
				changed = true
		placeholder else if adapter.ToolSearch && name == toolSearchProxyName {
				typed["type"] = "tool_search_call"
				typed["execution"] = "client"
				typed["arguments"] = json.RawMessage(toolSearchCallArgumentsJSON(rawObjectString(typed["arguments"])))
				delete(typed, "name")
				delete(typed, "namespace")
				changed = true
		placeholder
	placeholder
		for _, child := range typed {
			changed = restoreClientToolValue(child, adapter) || changed
	placeholder
placeholder
	return changed
placeholder

// ResponsesClientToolStreamRestorer restores client tool stream lifecycles.
// It is intentionally stateful because custom tools need their function
// arguments buffered until the upstream signals the call is complete.
type ResponsesClientToolStreamRestorer struct {
	adapter  ResponsesClientToolMapping
	nextSeq  int
	seenSeq  bool
	calls    map[string]*responsesClientToolStreamCall
	byOutput map[int]*responsesClientToolStreamCall
placeholder

type responsesClientToolStreamCall struct {
	kind      string
	name      string
	callID    string
	itemID    string
	outputIdx int
	arguments strings.Builder
placeholder

func NewResponsesClientToolStreamRestorer(mapping ResponsesClientToolMapping) *ResponsesClientToolStreamRestorer {
	return &ResponsesClientToolStreamRestorer{adapter: mapping, calls: make(map[string]*responsesClientToolStreamCall), byOutput: make(map[int]*responsesClientToolStreamCall)placeholder
placeholder

// Restore transforms one upstream SSE event into zero or more client events.
// Returned sequence numbers are continuous even when function argument events
// are suppressed or a custom completion expands into two events.
func (r *ResponsesClientToolStreamRestorer) Restore(event ResponsesStreamEvent) []ResponsesStreamEvent {
	if r == nil {
		return []ResponsesStreamEvent{eventplaceholder
placeholder
	if !r.seenSeq {
		r.nextSeq = event.SequenceNumber
		r.seenSeq = true
placeholder
	var out []ResponsesStreamEvent
	emit := func(event ResponsesStreamEvent) {
		event.SequenceNumber = r.nextSeq
		r.nextSeq++
		out = append(out, event)
placeholder

	switch event.Type {
	case "response.output_item.added":
		if call := r.recordItem(event); call != nil {
			if call.kind == "custom" {
				event.Item.Type = "custom_tool_call"
				event.Item.Input = ""
				event.Item.Arguments = ""
				event.Item.Namespace = ""
		placeholder else {
				event.Item.Type = "tool_search_call"
				event.Item.Name = ""
				event.Item.Arguments = "{placeholder"
				event.Item.Namespace = ""
		placeholder
	placeholder
		emit(r.restoreNamespaceEvent(event))
	case "response.function_call_arguments.delta":
		if call := r.callFor(event); call != nil {
			_, _ = call.arguments.WriteString(event.Delta)
			return nil
	placeholder
		emit(r.restoreNamespaceEvent(event))
	case "response.function_call_arguments.done":
		if call := r.callFor(event); call != nil {
			if event.Arguments != "" {
				call.arguments.Reset()
				_, _ = call.arguments.WriteString(event.Arguments)
		placeholder
			if call.kind == "custom" {
				input := extractCustomToolCallInput(call.arguments.String())
				if input != "" {
					emit(ResponsesStreamEvent{Type: "response.custom_tool_call_input.delta", OutputIndex: call.outputIdx, ItemID: call.itemID, Delta: inputplaceholder)
			placeholder
				emit(ResponsesStreamEvent{Type: "response.custom_tool_call_input.done", OutputIndex: call.outputIdx, ItemID: call.itemID, CallID: call.callID, Name: call.name, Input: inputplaceholder)
		placeholder
			return out
	placeholder
		emit(r.restoreNamespaceEvent(event))
	case "response.output_item.done":
		if call := r.recordItem(event); call != nil {
			if call.kind == "custom" {
				event.Item.Type = "custom_tool_call"
				event.Item.Input = extractCustomToolCallInput(call.arguments.String())
				event.Item.Arguments = ""
				event.Item.Namespace = ""
		placeholder else {
				event.Item.Type = "tool_search_call"
				event.Item.Name = ""
				event.Item.Arguments = call.arguments.String()
				if strings.TrimSpace(event.Item.Arguments) == "" {
					event.Item.Arguments = "{placeholder"
			placeholder
				event.Item.Namespace = ""
		placeholder
			delete(r.calls, call.itemID)
			delete(r.calls, call.callID)
			delete(r.byOutput, call.outputIdx)
	placeholder
		emit(r.restoreNamespaceEvent(event))
	default:
		// response.completed carries the non-stream representation.
		if event.Response != nil {
			restoreResponsesOutputClientTools(event.Response.Output, &r.adapter)
	placeholder
		emit(r.restoreNamespaceEvent(event))
placeholder
	return out
placeholder

// RestoreEvent restores one Responses SSE JSON data payload. Custom tool
// completions can expand to multiple payloads and proxy argument deltas can be
// intentionally dropped, hence the slice return value.
func (r *ResponsesClientToolStreamRestorer) RestoreEvent(payload []byte) ([][]byte, bool, error) {
	if len(payload) == 0 {
		return nil, false, nil
placeholder
	var wire struct {
		Type     string `json:"type"`
		Sequence int    `json:"sequence_number"`
placeholder
	if err := json.Unmarshal(payload, &wire); err != nil {
		return nil, false, err
placeholder
	if wire.Type == "response.completed" || wire.Type == "response.incomplete" || wire.Type == "response.failed" {
		restored, changed, err := RestoreResponsesClientToolPayload(payload, r.adapter)
		if err != nil {
			return nil, false, err
	placeholder
		return r.resequenceRaw(restored, wire.Sequence, changed)
placeholder
	if !clientToolLifecycleEvent(wire.Type) {
		return r.resequenceRaw(payload, wire.Sequence, false)
placeholder
	if !r.clientToolEventPayload(payload) {
		return r.resequenceRaw(payload, wire.Sequence, false)
placeholder
	var event ResponsesStreamEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, false, err
placeholder
	events := r.Restore(event)
	if len(events) == 1 {
		unchanged, err := json.Marshal(events[0])
		if err == nil && bytes.Equal(bytes.TrimSpace(unchanged), bytes.TrimSpace(payload)) {
			return [][]byte{payloadplaceholder, false, nil
	placeholder
placeholder
	result := make([][]byte, 0, len(events))
	for _, restored := range events {
		encoded, err := json.Marshal(restored)
		if err != nil {
			return nil, false, err
	placeholder
		result = append(result, encoded)
placeholder
	return result, true, nil
placeholder

func (r *ResponsesClientToolStreamRestorer) clientToolEventPayload(payload []byte) bool {
	var raw struct {
		ItemID      string `json:"item_id"`
		CallID      string `json:"call_id"`
		Name        string `json:"name"`
		OutputIndex int    `json:"output_index"`
		Item        *struct {
			Type   string `json:"type"`
			ID     string `json:"id"`
			CallID string `json:"call_id"`
			Name   string `json:"name"`
	placeholder `json:"item"`
placeholder
	if err := json.Unmarshal(payload, &raw); err != nil {
		return false
placeholder
	if raw.Item != nil {
		if raw.Item.Type != "function_call" {
			return false
	placeholder
		_, namespaceTool := r.adapter.NamespaceTools[raw.Item.Name]
		return r.adapter.CustomTools[raw.Item.Name] || (r.adapter.ToolSearch && raw.Item.Name == toolSearchProxyName) || namespaceTool || r.calls[raw.Item.ID] != nil || r.calls[raw.Item.CallID] != nil
placeholder
	if _, namespaceTool := r.adapter.NamespaceTools[raw.Name]; namespaceTool {
		return true
placeholder
	if r.calls[raw.ItemID] != nil || r.calls[raw.CallID] != nil || r.byOutput[raw.OutputIndex] != nil {
		return true
placeholder
	return false
placeholder

func clientToolLifecycleEvent(typ string) bool {
	switch typ {
	case "response.output_item.added", "response.output_item.done", "response.function_call_arguments.delta", "response.function_call_arguments.done":
		return true
	default:
		return false
placeholder
placeholder

// resequenceRaw deliberately keeps opaque upstream event fields untouched.
func (r *ResponsesClientToolStreamRestorer) resequenceRaw(payload []byte, sequence int, changed bool) ([][]byte, bool, error) {
	if !r.seenSeq {
		r.nextSeq, r.seenSeq = sequence, true
placeholder
	if r.nextSeq == sequence && !changed {
		r.nextSeq++
		return [][]byte{payloadplaceholder, false, nil
placeholder
	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, false, err
placeholder
	raw["sequence_number"] = r.nextSeq
	r.nextSeq++
	encoded, err := json.Marshal(raw)
	if err != nil {
		return nil, false, err
placeholder
	return [][]byte{encodedplaceholder, true, nil
placeholder

func (r *ResponsesClientToolStreamRestorer) recordItem(event ResponsesStreamEvent) *responsesClientToolStreamCall {
	if event.Item == nil || event.Item.Type != "function_call" {
		return nil
placeholder
	name := event.Item.Name
	kind := ""
	if r.adapter.CustomTools[name] {
		kind = "custom"
placeholder else if r.adapter.ToolSearch && name == toolSearchProxyName {
		kind = "tool_search"
placeholder
	if kind == "" {
		return nil
placeholder
	key := event.Item.ID
	if key == "" {
		key = event.Item.CallID
placeholder
	call := r.calls[key]
	if call == nil {
		call = &responsesClientToolStreamCall{kind: kind, name: name, callID: event.Item.CallID, itemID: event.Item.ID, outputIdx: event.OutputIndexplaceholder
		r.calls[key] = call
		if call.callID != "" {
			r.calls[call.callID] = call
	placeholder
		r.byOutput[call.outputIdx] = call
placeholder
	if event.Item.Arguments != "" {
		call.arguments.Reset()
		_, _ = call.arguments.WriteString(event.Item.Arguments)
placeholder
	return call
placeholder

func (r *ResponsesClientToolStreamRestorer) callFor(event ResponsesStreamEvent) *responsesClientToolStreamCall {
	if call := r.calls[event.ItemID]; call != nil {
		return call
placeholder
	if call := r.byOutput[event.OutputIndex]; call != nil {
		return call
placeholder
	for _, call := range r.calls {
		if (event.CallID != "" && call.callID == event.CallID) || (event.ItemID == "" && event.Name != "" && call.name == event.Name) {
			return call
	placeholder
placeholder
	return nil
placeholder

func (r *ResponsesClientToolStreamRestorer) restoreNamespaceEvent(event ResponsesStreamEvent) ResponsesStreamEvent {
	if len(r.adapter.NamespaceTools) == 0 {
		return event
placeholder
	if event.Item != nil && event.Item.Type == "function_call" {
		if name, ok := r.adapter.NamespaceTools[event.Item.Name]; ok {
			event.Item.Name, event.Item.Namespace = name.Name, name.Namespace
	placeholder
placeholder
	if event.Type == "response.function_call_arguments.done" {
		if name, ok := r.adapter.NamespaceTools[event.Name]; ok {
			event.Name = name.Name
	placeholder
placeholder
	return event
placeholder

func restoreResponsesOutputClientTools(outputs []ResponsesOutput, adapter *ResponsesClientToolMapping) {
	for index := range outputs {
		output := &outputs[index]
		if output.Type != "function_call" {
			continue
	placeholder
		if adapter.CustomTools[output.Name] {
			output.Type = "custom_tool_call"
			output.Input = extractCustomToolCallInput(output.Arguments)
			output.Arguments = ""
			output.Namespace = ""
	placeholder else if adapter.ToolSearch && output.Name == toolSearchProxyName {
			output.Type = "tool_search_call"
			output.Name = ""
			output.Namespace = ""
	placeholder
		if name, ok := adapter.NamespaceTools[output.Name]; ok && output.Type == "function_call" {
			output.Name, output.Namespace = name.Name, name.Namespace
	placeholder
placeholder
placeholder
