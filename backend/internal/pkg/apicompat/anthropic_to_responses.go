package apicompat

import (
	"encoding/json"
	"strings"
)

// AnthropicToResponses converts an Anthropic Messages request directly into
// a Responses API request. This preserves fields that would be lost in a
// Chat Completions intermediary round-trip (e.g. thinking, cache_control,
// structured system prompts).
func AnthropicToResponses(req *AnthropicRequest) (*ResponsesRequest, error) {
	input, err := convertAnthropicToResponsesInput(req.System, req.Messages)
	if err != nil {
		return nil, err
placeholder

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, err
placeholder

	out := &ResponsesRequest{
		Model:       req.Model,
		Input:       inputJSON,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
		Include:     []string{"reasoning.encrypted_content"placeholder,
placeholder

	storeFalse := false
	out.Store = &storeFalse

	if req.MaxTokens > 0 {
		v := req.MaxTokens
		if v < minMaxOutputTokens {
			v = minMaxOutputTokens
	placeholder
		out.MaxOutputTokens = &v
placeholder

	if len(req.Tools) > 0 {
		out.Tools = convertAnthropicToolsToResponses(req.Tools)
placeholder

	return out, nil
placeholder

// convertAnthropicToResponsesInput builds the Responses API input items array
// from the Anthropic system field and message list.
func convertAnthropicToResponsesInput(system json.RawMessage, msgs []AnthropicMessage) ([]ResponsesInputItem, error) {
	var out []ResponsesInputItem

	// System prompt → system role input item.
	if len(system) > 0 {
		sysText, err := parseAnthropicSystemPrompt(system)
		if err != nil {
			return nil, err
	placeholder
		if sysText != "" {
			content, _ := json.Marshal(sysText)
			out = append(out, ResponsesInputItem{
				Role:    "system",
				Content: content,
		placeholder)
	placeholder
placeholder

	for _, m := range msgs {
		items, err := anthropicMsgToResponsesItems(m)
		if err != nil {
			return nil, err
	placeholder
		out = append(out, items...)
placeholder
	return out, nil
placeholder

// parseAnthropicSystemPrompt handles the Anthropic system field which can be
// a plain string or an array of text blocks.
func parseAnthropicSystemPrompt(raw json.RawMessage) (string, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s, nil
placeholder
	var blocks []AnthropicContentBlock
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return "", err
placeholder
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" && b.Text != "" {
			parts = append(parts, b.Text)
	placeholder
placeholder
	return strings.Join(parts, "\n\n"), nil
placeholder

// anthropicMsgToResponsesItems converts a single Anthropic message into one
// or more Responses API input items.
func anthropicMsgToResponsesItems(m AnthropicMessage) ([]ResponsesInputItem, error) {
	switch m.Role {
	case "user":
		return anthropicUserToResponses(m.Content)
	case "assistant":
		return anthropicAssistantToResponses(m.Content)
	default:
		return anthropicUserToResponses(m.Content)
placeholder
placeholder

// anthropicUserToResponses handles an Anthropic user message. Content can be a
// plain string or an array of blocks. tool_result blocks are extracted into
// function_call_output items.
func anthropicUserToResponses(raw json.RawMessage) ([]ResponsesInputItem, error) {
	// Try plain string.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		content, _ := json.Marshal(s)
		return []ResponsesInputItem{{Role: "user", Content: contentplaceholderplaceholder, nil
placeholder

	var blocks []AnthropicContentBlock
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return nil, err
placeholder

	var out []ResponsesInputItem

	// Extract tool_result blocks → function_call_output items.
	for _, b := range blocks {
		if b.Type != "tool_result" {
			continue
	placeholder
		text := extractAnthropicToolResultText(b)
		if text == "" {
			// OpenAI Responses API requires "output" field; use placeholder for empty results.
			text = "(empty)"
	placeholder
		out = append(out, ResponsesInputItem{
			Type:   "function_call_output",
			CallID: toResponsesCallID(b.ToolUseID),
			Output: text,
	placeholder)
placeholder

	// Remaining text blocks → user message.
	text := extractAnthropicTextFromBlocks(blocks)
	if text != "" {
		content, _ := json.Marshal(text)
		out = append(out, ResponsesInputItem{Role: "user", Content: contentplaceholder)
placeholder

	return out, nil
placeholder

// anthropicAssistantToResponses handles an Anthropic assistant message.
// Text content → assistant message with output_text parts.
// tool_use blocks → function_call items.
// thinking blocks → ignored (OpenAI doesn't accept them as input).
func anthropicAssistantToResponses(raw json.RawMessage) ([]ResponsesInputItem, error) {
	// Try plain string.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		parts := []ResponsesContentPart{{Type: "output_text", Text: splaceholderplaceholder
		partsJSON, err := json.Marshal(parts)
		if err != nil {
			return nil, err
	placeholder
		return []ResponsesInputItem{{Role: "assistant", Content: partsJSONplaceholderplaceholder, nil
placeholder

	var blocks []AnthropicContentBlock
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return nil, err
placeholder

	var items []ResponsesInputItem

	// Text content → assistant message with output_text content parts.
	text := extractAnthropicTextFromBlocks(blocks)
	if text != "" {
		parts := []ResponsesContentPart{{Type: "output_text", Text: textplaceholderplaceholder
		partsJSON, err := json.Marshal(parts)
		if err != nil {
			return nil, err
	placeholder
		items = append(items, ResponsesInputItem{Role: "assistant", Content: partsJSONplaceholder)
placeholder

	// tool_use → function_call items.
	for _, b := range blocks {
		if b.Type != "tool_use" {
			continue
	placeholder
		args := "{placeholder"
		if len(b.Input) > 0 {
			args = string(b.Input)
	placeholder
		fcID := toResponsesCallID(b.ID)
		items = append(items, ResponsesInputItem{
			Type:      "function_call",
			CallID:    fcID,
			Name:      b.Name,
			Arguments: args,
			ID:        fcID,
	placeholder)
placeholder

	return items, nil
placeholder

// toResponsesCallID converts an Anthropic tool ID (toolu_xxx / call_xxx) to a
// Responses API function_call ID that starts with "fc_".
func toResponsesCallID(id string) string {
	if strings.HasPrefix(id, "fc_") {
		return id
placeholder
	return "fc_" + id
placeholder

// fromResponsesCallID reverses toResponsesCallID, stripping the "fc_" prefix
// that was added during request conversion.
func fromResponsesCallID(id string) string {
	if after, ok := strings.CutPrefix(id, "fc_"); ok {
		// Only strip if the remainder doesn't look like it was already "fc_" prefixed.
		// E.g. "fc_toolu_xxx" → "toolu_xxx", "fc_call_xxx" → "call_xxx"
		if strings.HasPrefix(after, "toolu_") || strings.HasPrefix(after, "call_") {
			return after
	placeholder
placeholder
	return id
placeholder

// extractAnthropicToolResultText gets the text content from a tool_result block.
func extractAnthropicToolResultText(b AnthropicContentBlock) string {
	if len(b.Content) == 0 {
		return ""
placeholder
	var s string
	if err := json.Unmarshal(b.Content, &s); err == nil {
		return s
placeholder
	var inner []AnthropicContentBlock
	if err := json.Unmarshal(b.Content, &inner); err == nil {
		var parts []string
		for _, ib := range inner {
			if ib.Type == "text" && ib.Text != "" {
				parts = append(parts, ib.Text)
		placeholder
	placeholder
		return strings.Join(parts, "\n\n")
placeholder
	return ""
placeholder

// extractAnthropicTextFromBlocks joins all text blocks, ignoring thinking/
// tool_use/tool_result blocks.
func extractAnthropicTextFromBlocks(blocks []AnthropicContentBlock) string {
	var parts []string
	for _, b := range blocks {
		if b.Type == "text" && b.Text != "" {
			parts = append(parts, b.Text)
	placeholder
placeholder
	return strings.Join(parts, "\n\n")
placeholder

// convertAnthropicToolsToResponses maps Anthropic tool definitions to
// Responses API function tools (input_schema → parameters).
func convertAnthropicToolsToResponses(tools []AnthropicTool) []ResponsesTool {
	out := make([]ResponsesTool, len(tools))
	for i, t := range tools {
		out[i] = ResponsesTool{
			Type:        "function",
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.InputSchema,
	placeholder
placeholder
	return out
placeholder
