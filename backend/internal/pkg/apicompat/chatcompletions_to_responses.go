package apicompat

import (
	"encoding/json"
	"fmt"
)

// ChatCompletionsToResponses converts a Chat Completions request into a
// Responses API request. The upstream always streams, so Stream is forced to
// true. store is always false and reasoning.encrypted_content is always
// included so that the response translator has full context.
func ChatCompletionsToResponses(req *ChatCompletionsRequest) (*ResponsesRequest, error) {
	input, err := convertChatMessagesToResponsesInput(req.Messages)
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
		Stream:      true, // upstream always streams
		Include:     []string{"reasoning.encrypted_content"placeholder,
		ServiceTier: req.ServiceTier,
placeholder

	storeFalse := false
	out.Store = &storeFalse

	// max_tokens / max_completion_tokens → max_output_tokens, prefer max_completion_tokens
	maxTokens := 0
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
placeholder
	if req.MaxCompletionTokens != nil {
		maxTokens = *req.MaxCompletionTokens
placeholder
	if maxTokens > 0 {
		v := maxTokens
		if v < minMaxOutputTokens {
			v = minMaxOutputTokens
	placeholder
		out.MaxOutputTokens = &v
placeholder

	// reasoning_effort → reasoning.effort + reasoning.summary="auto"
	if req.ReasoningEffort != "" {
		out.Reasoning = &ResponsesReasoning{
			Effort:  req.ReasoningEffort,
			Summary: "auto",
	placeholder
placeholder

	// tools[] and legacy functions[] → ResponsesTool[]
	if len(req.Tools) > 0 || len(req.Functions) > 0 {
		out.Tools = convertChatToolsToResponses(req.Tools, req.Functions)
placeholder

	// tool_choice: already compatible format — pass through directly.
	// Legacy function_call needs mapping.
	if len(req.ToolChoice) > 0 {
		out.ToolChoice = req.ToolChoice
placeholder else if len(req.FunctionCall) > 0 {
		tc, err := convertChatFunctionCallToToolChoice(req.FunctionCall)
		if err != nil {
			return nil, fmt.Errorf("convert function_call: %w", err)
	placeholder
		out.ToolChoice = tc
placeholder

	return out, nil
placeholder

// convertChatMessagesToResponsesInput converts the Chat Completions messages
// array into a Responses API input items array.
func convertChatMessagesToResponsesInput(msgs []ChatMessage) ([]ResponsesInputItem, error) {
	var out []ResponsesInputItem
	for _, m := range msgs {
		items, err := chatMessageToResponsesItems(m)
		if err != nil {
			return nil, err
	placeholder
		out = append(out, items...)
placeholder
	return out, nil
placeholder

// chatMessageToResponsesItems converts a single ChatMessage into one or more
// ResponsesInputItem values.
func chatMessageToResponsesItems(m ChatMessage) ([]ResponsesInputItem, error) {
	switch m.Role {
	case "system":
		return chatSystemToResponses(m)
	case "user":
		return chatUserToResponses(m)
	case "assistant":
		return chatAssistantToResponses(m)
	case "tool":
		return chatToolToResponses(m)
	case "function":
		return chatFunctionToResponses(m)
	default:
		return chatUserToResponses(m)
placeholder
placeholder

// chatSystemToResponses converts a system message.
func chatSystemToResponses(m ChatMessage) ([]ResponsesInputItem, error) {
	text, err := parseChatContent(m.Content)
	if err != nil {
		return nil, err
placeholder
	content, err := json.Marshal(text)
	if err != nil {
		return nil, err
placeholder
	return []ResponsesInputItem{{Role: "system", Content: contentplaceholderplaceholder, nil
placeholder

// chatUserToResponses converts a user message, handling both plain strings and
// multi-modal content arrays.
func chatUserToResponses(m ChatMessage) ([]ResponsesInputItem, error) {
	// Try plain string first.
	var s string
	if err := json.Unmarshal(m.Content, &s); err == nil {
		content, _ := json.Marshal(s)
		return []ResponsesInputItem{{Role: "user", Content: contentplaceholderplaceholder, nil
placeholder

	var parts []ChatContentPart
	if err := json.Unmarshal(m.Content, &parts); err != nil {
		return nil, fmt.Errorf("parse user content: %w", err)
placeholder

	var responseParts []ResponsesContentPart
	for _, p := range parts {
		switch p.Type {
		case "text":
			if p.Text != "" {
				responseParts = append(responseParts, ResponsesContentPart{
					Type: "input_text",
					Text: p.Text,
			placeholder)
		placeholder
		case "image_url":
			if p.ImageURL != nil && p.ImageURL.URL != "" {
				responseParts = append(responseParts, ResponsesContentPart{
					Type:     "input_image",
					ImageURL: p.ImageURL.URL,
			placeholder)
		placeholder
	placeholder
placeholder

	content, err := json.Marshal(responseParts)
	if err != nil {
		return nil, err
placeholder
	return []ResponsesInputItem{{Role: "user", Content: contentplaceholderplaceholder, nil
placeholder

// chatAssistantToResponses converts an assistant message. If there is both
// text content and tool_calls, the text is emitted as an assistant message
// first, then each tool_call becomes a function_call item. If the content is
// empty/nil and there are tool_calls, only function_call items are emitted.
func chatAssistantToResponses(m ChatMessage) ([]ResponsesInputItem, error) {
	var items []ResponsesInputItem

	// Emit assistant message with output_text if content is non-empty.
	if len(m.Content) > 0 {
		var s string
		if err := json.Unmarshal(m.Content, &s); err == nil && s != "" {
			parts := []ResponsesContentPart{{Type: "output_text", Text: splaceholderplaceholder
			partsJSON, err := json.Marshal(parts)
			if err != nil {
				return nil, err
		placeholder
			items = append(items, ResponsesInputItem{Role: "assistant", Content: partsJSONplaceholder)
	placeholder
placeholder

	// Emit one function_call item per tool_call.
	for _, tc := range m.ToolCalls {
		args := tc.Function.Arguments
		if args == "" {
			args = "{placeholder"
	placeholder
		items = append(items, ResponsesInputItem{
			Type:      "function_call",
			CallID:    tc.ID,
			Name:      tc.Function.Name,
			Arguments: args,
			ID:        tc.ID,
	placeholder)
placeholder

	return items, nil
placeholder

// chatToolToResponses converts a tool result message (role=tool) into a
// function_call_output item.
func chatToolToResponses(m ChatMessage) ([]ResponsesInputItem, error) {
	output, err := parseChatContent(m.Content)
	if err != nil {
		return nil, err
placeholder
	if output == "" {
		output = "(empty)"
placeholder
	return []ResponsesInputItem{{
		Type:   "function_call_output",
		CallID: m.ToolCallID,
		Output: output,
placeholderplaceholder, nil
placeholder

// chatFunctionToResponses converts a legacy function result message
// (role=function) into a function_call_output item. The Name field is used as
// call_id since legacy function calls do not carry a separate call_id.
func chatFunctionToResponses(m ChatMessage) ([]ResponsesInputItem, error) {
	output, err := parseChatContent(m.Content)
	if err != nil {
		return nil, err
placeholder
	if output == "" {
		output = "(empty)"
placeholder
	return []ResponsesInputItem{{
		Type:   "function_call_output",
		CallID: m.Name,
		Output: output,
placeholderplaceholder, nil
placeholder

// parseChatContent returns the string value of a ChatMessage Content field.
// Content must be a JSON string. Returns "" if content is null or empty.
func parseChatContent(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
placeholder
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", fmt.Errorf("parse content as string: %w", err)
placeholder
	return s, nil
placeholder

// convertChatToolsToResponses maps Chat Completions tool definitions and legacy
// function definitions to Responses API tool definitions.
func convertChatToolsToResponses(tools []ChatTool, functions []ChatFunction) []ResponsesTool {
	var out []ResponsesTool

	for _, t := range tools {
		if t.Type != "function" || t.Function == nil {
			continue
	placeholder
		rt := ResponsesTool{
			Type:        "function",
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
			Strict:      t.Function.Strict,
	placeholder
		out = append(out, rt)
placeholder

	// Legacy functions[] are treated as function-type tools.
	for _, f := range functions {
		rt := ResponsesTool{
			Type:        "function",
			Name:        f.Name,
			Description: f.Description,
			Parameters:  f.Parameters,
			Strict:      f.Strict,
	placeholder
		out = append(out, rt)
placeholder

	return out
placeholder

// convertChatFunctionCallToToolChoice maps the legacy function_call field to a
// Responses API tool_choice value.
//
//	"auto" → "auto"
//	"none" → "none"
//	{"name":"X"placeholder → {"type":"function","function":{"name":"X"placeholderplaceholder
func convertChatFunctionCallToToolChoice(raw json.RawMessage) (json.RawMessage, error) {
	// Try string first ("auto", "none", etc.) — pass through as-is.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return json.Marshal(s)
placeholder

	// Object form: {"name":"X"placeholder
	var obj struct {
		Name string `json:"name"`
placeholder
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
placeholder
	return json.Marshal(map[string]any{
		"type":     "function",
		"function": map[string]string{"name": obj.Nameplaceholder,
placeholder)
placeholder
