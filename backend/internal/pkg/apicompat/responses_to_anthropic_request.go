package apicompat

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResponsesToAnthropicRequest converts a Responses API request into an
// Anthropic Messages request. This is the reverse of AnthropicToResponses and
// enables Anthropic platform groups to accept OpenAI Responses API requests
// by converting them to the native /v1/messages format before forwarding upstream.
func ResponsesToAnthropicRequest(req *ResponsesRequest) (*AnthropicRequest, error) {
	system, messages, err := convertResponsesInputToAnthropic(req.Input)
	if err != nil {
		return nil, err
placeholder

	out := &AnthropicRequest{
		Model:       req.Model,
		Messages:    messages,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      req.Stream,
placeholder

	if len(system) > 0 {
		out.System = system
placeholder

	// max_output_tokens → max_tokens
	if req.MaxOutputTokens != nil && *req.MaxOutputTokens > 0 {
		out.MaxTokens = *req.MaxOutputTokens
placeholder
	if out.MaxTokens == 0 {
		// Anthropic requires max_tokens; default to a sensible value.
		out.MaxTokens = 8192
placeholder

	// Convert tools
	if len(req.Tools) > 0 {
		out.Tools = convertResponsesToAnthropicTools(req.Tools)
placeholder

	// Convert tool_choice (reverse of convertAnthropicToolChoiceToResponses)
	if len(req.ToolChoice) > 0 {
		tc, err := convertResponsesToAnthropicToolChoice(req.ToolChoice)
		if err != nil {
			return nil, fmt.Errorf("convert tool_choice: %w", err)
	placeholder
		out.ToolChoice = tc
placeholder

	// reasoning.effort → output_config.effort + thinking
	if req.Reasoning != nil && req.Reasoning.Effort != "" {
		effort := mapResponsesEffortToAnthropic(req.Reasoning.Effort)
		out.OutputConfig = &AnthropicOutputConfig{Effort: effortplaceholder
		// Enable thinking for non-low efforts
		if effort != "low" {
			out.Thinking = &AnthropicThinking{
				Type:         "enabled",
				BudgetTokens: defaultThinkingBudget(effort),
		placeholder
	placeholder
placeholder

	return out, nil
placeholder

// defaultThinkingBudget returns a sensible thinking budget based on effort level.
func defaultThinkingBudget(effort string) int {
	switch effort {
	case "low":
		return 1024
	case "medium":
		return 4096
	case "high":
		return 10240
	case "max":
		return 32768
	default:
		return 10240
placeholder
placeholder

// mapResponsesEffortToAnthropic converts OpenAI Responses reasoning effort to
// Anthropic effort levels. Reverse of mapAnthropicEffortToResponses.
//
//	low    → low
//	medium → medium
//	high   → high
//	xhigh  → max
func mapResponsesEffortToAnthropic(effort string) string {
	if effort == "xhigh" {
		return "max"
placeholder
	return effort // low→low, medium→medium, high→high, unknown→passthrough
placeholder

// convertResponsesInputToAnthropic extracts system prompt and messages from
// a Responses API input array. Returns the system as raw JSON (for Anthropic's
// polymorphic system field) and a list of Anthropic messages.
func convertResponsesInputToAnthropic(inputRaw json.RawMessage) (json.RawMessage, []AnthropicMessage, error) {
	// Try as plain string input.
	var inputStr string
	if err := json.Unmarshal(inputRaw, &inputStr); err == nil {
		content, _ := json.Marshal(inputStr)
		return nil, []AnthropicMessage{{Role: "user", Content: contentplaceholderplaceholder, nil
placeholder

	var items []ResponsesInputItem
	if err := json.Unmarshal(inputRaw, &items); err != nil {
		return nil, nil, fmt.Errorf("parse responses input: %w", err)
placeholder

	var system json.RawMessage
	var messages []AnthropicMessage

	for _, item := range items {
		switch {
		case item.Role == "system":
			// System prompt → Anthropic system field
			text := extractTextFromContent(item.Content)
			if text != "" {
				system, _ = json.Marshal(text)
		placeholder

		case item.Type == "function_call":
			// function_call → assistant message with tool_use block
			input := json.RawMessage("{placeholder")
			if item.Arguments != "" {
				input = json.RawMessage(item.Arguments)
		placeholder
			block := AnthropicContentBlock{
				Type:  "tool_use",
				ID:    fromResponsesCallIDToAnthropic(item.CallID),
				Name:  item.Name,
				Input: input,
		placeholder
			blockJSON, _ := json.Marshal([]AnthropicContentBlock{blockplaceholder)
			messages = append(messages, AnthropicMessage{
				Role:    "assistant",
				Content: blockJSON,
		placeholder)

		case item.Type == "function_call_output":
			// function_call_output → user message with tool_result block
			outputContent := item.Output
			if outputContent == "" {
				outputContent = "(empty)"
		placeholder
			contentJSON, _ := json.Marshal(outputContent)
			block := AnthropicContentBlock{
				Type:      "tool_result",
				ToolUseID: fromResponsesCallIDToAnthropic(item.CallID),
				Content:   contentJSON,
		placeholder
			blockJSON, _ := json.Marshal([]AnthropicContentBlock{blockplaceholder)
			messages = append(messages, AnthropicMessage{
				Role:    "user",
				Content: blockJSON,
		placeholder)

		case item.Role == "user":
			content, err := convertResponsesUserToAnthropicContent(item.Content)
			if err != nil {
				return nil, nil, err
		placeholder
			messages = append(messages, AnthropicMessage{
				Role:    "user",
				Content: content,
		placeholder)

		case item.Role == "assistant":
			content, err := convertResponsesAssistantToAnthropicContent(item.Content)
			if err != nil {
				return nil, nil, err
		placeholder
			messages = append(messages, AnthropicMessage{
				Role:    "assistant",
				Content: content,
		placeholder)

		default:
			// Unknown role/type — attempt as user message
			if item.Content != nil {
				messages = append(messages, AnthropicMessage{
					Role:    "user",
					Content: item.Content,
			placeholder)
		placeholder
	placeholder
placeholder

	// Merge consecutive same-role messages (Anthropic requires alternating roles)
	messages = mergeConsecutiveMessages(messages)

	return system, messages, nil
placeholder

// extractTextFromContent extracts text from a content field that may be a
// plain string or an array of content parts.
func extractTextFromContent(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
placeholder
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
placeholder
	var parts []ResponsesContentPart
	if err := json.Unmarshal(raw, &parts); err == nil {
		var texts []string
		for _, p := range parts {
			if (p.Type == "input_text" || p.Type == "output_text" || p.Type == "text") && p.Text != "" {
				texts = append(texts, p.Text)
		placeholder
	placeholder
		return strings.Join(texts, "\n\n")
placeholder
	return ""
placeholder

// convertResponsesUserToAnthropicContent converts a Responses user message
// content field into Anthropic content blocks JSON.
func convertResponsesUserToAnthropicContent(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return json.Marshal("") // empty string content
placeholder

	// Try plain string.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return json.Marshal(s)
placeholder

	// Array of content parts → Anthropic content blocks.
	var parts []ResponsesContentPart
	if err := json.Unmarshal(raw, &parts); err != nil {
		// Pass through as-is if we can't parse
		return raw, nil
placeholder

	var blocks []AnthropicContentBlock
	for _, p := range parts {
		switch p.Type {
		case "input_text", "text":
			if p.Text != "" {
				blocks = append(blocks, AnthropicContentBlock{
					Type: "text",
					Text: p.Text,
			placeholder)
		placeholder
		case "input_image":
			src := dataURIToAnthropicImageSource(p.ImageURL)
			if src != nil {
				blocks = append(blocks, AnthropicContentBlock{
					Type:   "image",
					Source: src,
			placeholder)
		placeholder
	placeholder
placeholder

	if len(blocks) == 0 {
		return json.Marshal("")
placeholder
	return json.Marshal(blocks)
placeholder

// convertResponsesAssistantToAnthropicContent converts a Responses assistant
// message content field into Anthropic content blocks JSON.
func convertResponsesAssistantToAnthropicContent(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return json.Marshal([]AnthropicContentBlock{{Type: "text", Text: ""placeholderplaceholder)
placeholder

	// Try plain string.
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return json.Marshal([]AnthropicContentBlock{{Type: "text", Text: splaceholderplaceholder)
placeholder

	// Array of content parts → Anthropic content blocks.
	var parts []ResponsesContentPart
	if err := json.Unmarshal(raw, &parts); err != nil {
		return raw, nil
placeholder

	var blocks []AnthropicContentBlock
	for _, p := range parts {
		switch p.Type {
		case "output_text", "text":
			if p.Text != "" {
				blocks = append(blocks, AnthropicContentBlock{
					Type: "text",
					Text: p.Text,
			placeholder)
		placeholder
	placeholder
placeholder

	if len(blocks) == 0 {
		blocks = append(blocks, AnthropicContentBlock{Type: "text", Text: ""placeholder)
placeholder
	return json.Marshal(blocks)
placeholder

// fromResponsesCallIDToAnthropic converts an OpenAI function call ID back to
// Anthropic format. Reverses toResponsesCallID.
func fromResponsesCallIDToAnthropic(id string) string {
	// If it has our "fc_" prefix wrapping a known Anthropic prefix, strip it
	if after, ok := strings.CutPrefix(id, "fc_"); ok {
		if strings.HasPrefix(after, "toolu_") || strings.HasPrefix(after, "call_") {
			return after
	placeholder
placeholder
	// Generate a synthetic Anthropic tool ID
	if !strings.HasPrefix(id, "toolu_") && !strings.HasPrefix(id, "call_") {
		return "toolu_" + id
placeholder
	return id
placeholder

// dataURIToAnthropicImageSource parses a data URI into an AnthropicImageSource.
func dataURIToAnthropicImageSource(dataURI string) *AnthropicImageSource {
	if !strings.HasPrefix(dataURI, "data:") {
		return nil
placeholder
	// Format: data:<media_type>;base64,<data>
	rest := strings.TrimPrefix(dataURI, "data:")
	semicolonIdx := strings.Index(rest, ";")
	if semicolonIdx < 0 {
		return nil
placeholder
	mediaType := rest[:semicolonIdx]
	rest = rest[semicolonIdx+1:]
	if !strings.HasPrefix(rest, "base64,") {
		return nil
placeholder
	data := strings.TrimPrefix(rest, "base64,")
	return &AnthropicImageSource{
		Type:      "base64",
		MediaType: mediaType,
		Data:      data,
placeholder
placeholder

// mergeConsecutiveMessages merges consecutive messages with the same role
// because Anthropic requires alternating user/assistant turns.
func mergeConsecutiveMessages(messages []AnthropicMessage) []AnthropicMessage {
	if len(messages) <= 1 {
		return messages
placeholder

	var merged []AnthropicMessage
	for _, msg := range messages {
		if len(merged) == 0 || merged[len(merged)-1].Role != msg.Role {
			merged = append(merged, msg)
			continue
	placeholder

		// Same role — merge content arrays
		last := &merged[len(merged)-1]
		lastBlocks := parseContentBlocks(last.Content)
		newBlocks := parseContentBlocks(msg.Content)
		combined := append(lastBlocks, newBlocks...)
		last.Content, _ = json.Marshal(combined)
placeholder
	return merged
placeholder

// parseContentBlocks attempts to parse content as []AnthropicContentBlock.
// If it's a string, wraps it in a text block.
func parseContentBlocks(raw json.RawMessage) []AnthropicContentBlock {
	var blocks []AnthropicContentBlock
	if err := json.Unmarshal(raw, &blocks); err == nil {
		return blocks
placeholder
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return []AnthropicContentBlock{{Type: "text", Text: splaceholderplaceholder
placeholder
	return nil
placeholder

// convertResponsesToAnthropicTools maps Responses API tools to Anthropic format.
// Reverse of convertAnthropicToolsToResponses.
func convertResponsesToAnthropicTools(tools []ResponsesTool) []AnthropicTool {
	var out []AnthropicTool
	for _, t := range tools {
		switch t.Type {
		case "web_search":
			out = append(out, AnthropicTool{
				Type: "web_search_20250305",
				Name: "web_search",
		placeholder)
		case "function":
			out = append(out, AnthropicTool{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: normalizeAnthropicInputSchema(t.Parameters),
		placeholder)
		default:
			// Pass through unknown tool types
			out = append(out, AnthropicTool{
				Type:        t.Type,
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.Parameters,
		placeholder)
	placeholder
placeholder
	return out
placeholder

// normalizeAnthropicInputSchema ensures the input_schema has a "type" field.
func normalizeAnthropicInputSchema(schema json.RawMessage) json.RawMessage {
	if len(schema) == 0 || string(schema) == "null" {
		return json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)
placeholder
	return schema
placeholder

// convertResponsesToAnthropicToolChoice maps Responses tool_choice to Anthropic format.
// Reverse of convertAnthropicToolChoiceToResponses.
//
//	"auto"                                     → {"type":"auto"placeholder
//	"required"                                 → {"type":"any"placeholder
//	"none"                                     → {"type":"none"placeholder
//	{"type":"function","function":{"name":"X"placeholderplaceholder → {"type":"tool","name":"X"placeholder
func convertResponsesToAnthropicToolChoice(raw json.RawMessage) (json.RawMessage, error) {
	// Try as string first
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		switch s {
		case "auto":
			return json.Marshal(map[string]string{"type": "auto"placeholder)
		case "required":
			return json.Marshal(map[string]string{"type": "any"placeholder)
		case "none":
			return json.Marshal(map[string]string{"type": "none"placeholder)
		default:
			return raw, nil
	placeholder
placeholder

	// Try as object with type=function
	var tc struct {
		Type     string `json:"type"`
		Function struct {
			Name string `json:"name"`
	placeholder `json:"function"`
placeholder
	if err := json.Unmarshal(raw, &tc); err == nil && tc.Type == "function" && tc.Function.Name != "" {
		return json.Marshal(map[string]string{
			"type": "tool",
			"name": tc.Function.Name,
	placeholder)
placeholder

	// Pass through unknown
	return raw, nil
placeholder
