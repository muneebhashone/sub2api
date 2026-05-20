package apicompat

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ResponsesToChatCompletionsRequest converts a Responses API request into a
// Chat Completions request for upstreams that only implement
// /v1/chat/completions.
func ResponsesToChatCompletionsRequest(req *ResponsesRequest) (*ChatCompletionsRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("responses request is nil")
placeholder

	messages, err := responsesInputToChatMessages(req.Instructions, req.Input)
	if err != nil {
		return nil, err
placeholder

	out := &ChatCompletionsRequest{
		Model:               req.Model,
		Messages:            messages,
		MaxCompletionTokens: req.MaxOutputTokens,
		Temperature:         req.Temperature,
		TopP:                req.TopP,
		Stream:              req.Stream,
		ServiceTier:         req.ServiceTier,
placeholder
	if req.Reasoning != nil {
		out.ReasoningEffort = req.Reasoning.Effort
placeholder
	if len(req.Tools) > 0 {
		out.Tools = responsesToolsToChatTools(req.Tools)
placeholder
	if len(req.ToolChoice) > 0 {
		out.ToolChoice = responsesToolChoiceToChatToolChoice(req.ToolChoice)
placeholder

	return out, nil
placeholder

func responsesInputToChatMessages(instructions string, inputRaw json.RawMessage) ([]ChatMessage, error) {
	var messages []ChatMessage
	if strings.TrimSpace(instructions) != "" {
		content, _ := json.Marshal(instructions)
		messages = append(messages, ChatMessage{
			Role:    "system",
			Content: content,
	placeholder)
placeholder

	inputRaw = bytesTrimSpace(inputRaw)
	if len(inputRaw) == 0 || string(inputRaw) == "null" {
		return messages, nil
placeholder

	var inputText string
	if err := json.Unmarshal(inputRaw, &inputText); err == nil {
		content, _ := json.Marshal(inputText)
		messages = append(messages, ChatMessage{
			Role:    "user",
			Content: content,
	placeholder)
		return messages, nil
placeholder

	var rawItems []json.RawMessage
	if err := json.Unmarshal(inputRaw, &rawItems); err != nil {
		return nil, fmt.Errorf("parse responses input: %w", err)
placeholder

	for _, raw := range rawItems {
		raw = bytesTrimSpace(raw)
		if len(raw) == 0 || string(raw) == "null" {
			continue
	placeholder

		var item map[string]json.RawMessage
		if err := json.Unmarshal(raw, &item); err != nil {
			var text string
			if textErr := json.Unmarshal(raw, &text); textErr == nil {
				content, _ := json.Marshal(text)
				messages = append(messages, ChatMessage{Role: "user", Content: contentplaceholder)
				continue
		placeholder
			return nil, fmt.Errorf("parse responses input item: %w", err)
	placeholder

		role := rawString(item["role"])
		itemType := rawString(item["type"])
		switch itemType {
		case "function_call":
			arguments := rawString(item["arguments"])
			if strings.TrimSpace(arguments) == "" {
				arguments = "{placeholder"
		placeholder
			messages = append(messages, ChatMessage{
				Role: "assistant",
				ToolCalls: []ChatToolCall{{
					ID:   rawString(item["call_id"]),
					Type: "function",
					Function: ChatFunctionCall{
						Name:      rawString(item["name"]),
						Arguments: arguments,
				placeholder,
		placeholder
		placeholder)
			continue
		case "function_call_output":
			content, _ := json.Marshal(rawString(item["output"]))
			messages = append(messages, ChatMessage{
				Role:       "tool",
				ToolCallID: rawString(item["call_id"]),
				Content:    content,
		placeholder)
			continue
		case "input_text", "text":
			content, _ := json.Marshal(rawString(item["text"]))
			messages = append(messages, ChatMessage{Role: "user", Content: contentplaceholder)
			continue
		case "input_image":
			content, err := chatContentFromSingleResponsesPart(itemType, item)
			if err != nil {
				return nil, err
		placeholder
			messages = append(messages, ChatMessage{Role: "user", Content: contentplaceholder)
			continue
	placeholder

		if role == "" {
			role = "user"
	placeholder
		content := item["content"]
		if len(bytesTrimSpace(content)) == 0 {
			if text := rawString(item["text"]); text != "" {
				content, _ = json.Marshal(text)
		placeholder
	placeholder
		chatContent, err := responsesContentToChatContent(content, role)
		if err != nil {
			return nil, err
	placeholder
		messages = append(messages, ChatMessage{
			Role:    role,
			Content: chatContent,
	placeholder)
placeholder

	return messages, nil
placeholder

func responsesContentToChatContent(raw json.RawMessage, role string) (json.RawMessage, error) {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		empty, _ := json.Marshal("")
		return empty, nil
placeholder

	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return raw, nil
placeholder

	var rawParts []json.RawMessage
	if err := json.Unmarshal(raw, &rawParts); err == nil {
		return responsesContentPartsToChatContent(rawParts, role)
placeholder

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err == nil {
		return chatContentFromSingleResponsesPart(rawString(obj["type"]), obj)
placeholder

	return raw, nil
placeholder

func responsesContentPartsToChatContent(rawParts []json.RawMessage, role string) (json.RawMessage, error) {
	var textParts []string
	var chatParts []ChatContentPart
	hasNonText := false

	for _, rawPart := range rawParts {
		var part map[string]json.RawMessage
		if err := json.Unmarshal(rawPart, &part); err != nil {
			continue
	placeholder
		partType := rawString(part["type"])
		switch partType {
		case "input_text", "output_text", "text", "":
			text := rawString(part["text"])
			if text == "" {
				continue
		placeholder
			textParts = append(textParts, text)
			chatParts = append(chatParts, ChatContentPart{Type: "text", Text: textplaceholder)
		case "input_image", "image_url":
			imageURL := rawString(part["image_url"])
			if imageURL == "" {
				imageURL = rawNestedString(part["image_url"], "url")
		placeholder
			if imageURL == "" {
				continue
		placeholder
			hasNonText = true
			chatParts = append(chatParts, ChatContentPart{
				Type:     "image_url",
				ImageURL: &ChatImageURL{URL: imageURLplaceholder,
		placeholder)
	placeholder
placeholder

	if !hasNonText {
		joined, _ := json.Marshal(strings.Join(textParts, "\n\n"))
		return joined, nil
placeholder
	if role != "user" {
		joined, _ := json.Marshal(strings.Join(textParts, "\n\n"))
		return joined, nil
placeholder
	if len(chatParts) == 0 {
		empty, _ := json.Marshal("")
		return empty, nil
placeholder
	return json.Marshal(chatParts)
placeholder

func chatContentFromSingleResponsesPart(partType string, part map[string]json.RawMessage) (json.RawMessage, error) {
	switch partType {
	case "input_image", "image_url":
		imageURL := rawString(part["image_url"])
		if imageURL == "" {
			imageURL = rawNestedString(part["image_url"], "url")
	placeholder
		return json.Marshal([]ChatContentPart{{
			Type:     "image_url",
			ImageURL: &ChatImageURL{URL: imageURLplaceholder,
	placeholderplaceholder)
	default:
		return json.Marshal(rawString(part["text"]))
placeholder
placeholder

func responsesToolsToChatTools(tools []ResponsesTool) []ChatTool {
	out := make([]ChatTool, 0, len(tools))
	for _, tool := range tools {
		if tool.Type != "function" {
			continue
	placeholder
		out = append(out, ChatTool{
			Type: "function",
			Function: &ChatFunction{
				Name:        tool.Name,
				Description: tool.Description,
				Parameters:  tool.Parameters,
				Strict:      tool.Strict,
		placeholder,
	placeholder)
placeholder
	return out
placeholder

func responsesToolChoiceToChatToolChoice(raw json.RawMessage) json.RawMessage {
	var choice map[string]json.RawMessage
	if err := json.Unmarshal(raw, &choice); err != nil {
		return raw
placeholder
	if rawString(choice["type"]) != "function" {
		return raw
placeholder
	name := rawString(choice["name"])
	if name == "" {
		name = rawNestedString(choice["function"], "name")
placeholder
	if name == "" {
		return raw
placeholder
	out, err := json.Marshal(map[string]any{
		"type": "function",
		"function": map[string]string{
			"name": name,
	placeholder,
placeholder)
	if err != nil {
		return raw
placeholder
	return out
placeholder

// ChatCompletionsResponseToResponses converts a non-streaming Chat Completions
// response into a Responses API response.
func ChatCompletionsResponseToResponses(resp *ChatCompletionsResponse, model string) *ResponsesResponse {
	id := ""
	if resp != nil {
		id = resp.ID
placeholder
	if id == "" {
		id = generateResponsesID()
placeholder

	out := &ResponsesResponse{
		ID:     id,
		Object: "response",
		Model:  model,
		Status: "completed",
placeholder
	if resp == nil {
		out.Output = []ResponsesOutput{emptyResponsesMessageOutput()placeholder
		return out
placeholder
	if out.Model == "" {
		out.Model = resp.Model
placeholder

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		out.Output = chatMessageToResponsesOutput(choice.Message)
		if choice.FinishReason == "length" {
			out.Status = "incomplete"
			out.IncompleteDetails = &ResponsesIncompleteDetails{Reason: "max_output_tokens"placeholder
	placeholder
placeholder
	if len(out.Output) == 0 {
		out.Output = []ResponsesOutput{emptyResponsesMessageOutput()placeholder
placeholder
	if resp.Usage != nil {
		out.Usage = ChatUsageToResponsesUsage(resp.Usage)
placeholder
	return out
placeholder

func chatMessageToResponsesOutput(message ChatMessage) []ResponsesOutput {
	var outputs []ResponsesOutput
	if message.ReasoningContent != "" {
		outputs = append(outputs, ResponsesOutput{
			Type: "reasoning",
			ID:   generateItemID(),
			Summary: []ResponsesSummary{{
				Type: "summary_text",
				Text: message.ReasoningContent,
	placeholder
	placeholder)
placeholder

	text := chatMessageContentText(message.Content)
	if text != "" || len(message.ToolCalls) == 0 {
		outputs = append(outputs, ResponsesOutput{
			Type: "message",
			ID:   generateItemID(),
			Role: "assistant",
			Content: []ResponsesContentPart{{
				Type: "output_text",
				Text: text,
	placeholder
			Status: "completed",
	placeholder)
placeholder

	for _, toolCall := range message.ToolCalls {
		arguments := toolCall.Function.Arguments
		if strings.TrimSpace(arguments) == "" {
			arguments = "{placeholder"
	placeholder
		outputs = append(outputs, ResponsesOutput{
			Type:      "function_call",
			ID:        generateItemID(),
			CallID:    toolCall.ID,
			Name:      toolCall.Function.Name,
			Arguments: arguments,
			Status:    "completed",
	placeholder)
placeholder

	return outputs
placeholder

func emptyResponsesMessageOutput() ResponsesOutput {
	return ResponsesOutput{
		Type:    "message",
		ID:      generateItemID(),
		Role:    "assistant",
		Content: []ResponsesContentPart{{Type: "output_text", Text: ""placeholderplaceholder,
		Status:  "completed",
placeholder
placeholder

func chatMessageContentText(raw json.RawMessage) string {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		return ""
placeholder
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
placeholder
	var parts []ChatContentPart
	if err := json.Unmarshal(raw, &parts); err == nil {
		var texts []string
		for _, part := range parts {
			if part.Type == "text" && part.Text != "" {
				texts = append(texts, part.Text)
		placeholder
	placeholder
		return strings.Join(texts, "\n\n")
placeholder
	return ""
placeholder

// ChatUsageToResponsesUsage converts Chat Completions token usage to Responses
// usage shape.
func ChatUsageToResponsesUsage(usage *ChatUsage) *ResponsesUsage {
	if usage == nil {
		return nil
placeholder
	out := &ResponsesUsage{
		InputTokens:  usage.PromptTokens,
		OutputTokens: usage.CompletionTokens,
		TotalTokens:  usage.TotalTokens,
placeholder
	if out.TotalTokens == 0 {
		out.TotalTokens = out.InputTokens + out.OutputTokens
placeholder
	if usage.PromptTokensDetails != nil && usage.PromptTokensDetails.CachedTokens > 0 {
		out.InputTokensDetails = &ResponsesInputTokensDetails{
			CachedTokens: usage.PromptTokensDetails.CachedTokens,
	placeholder
placeholder
	return out
placeholder

// ChatCompletionsToResponsesStreamState tracks state while converting Chat
// Completions SSE chunks into Responses SSE events.
type ChatCompletionsToResponsesStreamState struct {
	ResponseID     string
	Model          string
	Created        int64
	SequenceNumber int
	CreatedSent    bool
	CompletedSent  bool

	MessageItemID string
	Text          strings.Builder
	Reasoning     strings.Builder
	ToolCalls     map[int]*ChatToolCall

	FinishReason string
	Usage        *ResponsesUsage
placeholder

// NewChatCompletionsToResponsesStreamState returns an initialized stream state.
func NewChatCompletionsToResponsesStreamState(model string) *ChatCompletionsToResponsesStreamState {
	return &ChatCompletionsToResponsesStreamState{
		ResponseID: generateResponsesID(),
		Model:      model,
		Created:    time.Now().Unix(),
		ToolCalls:  make(map[int]*ChatToolCall),
placeholder
placeholder

// ChatCompletionsChunkToResponsesEvents converts one Chat Completions stream
// chunk into zero or more Responses stream events.
func ChatCompletionsChunkToResponsesEvents(
	chunk *ChatCompletionsChunk,
	state *ChatCompletionsToResponsesStreamState,
) []ResponsesStreamEvent {
	if chunk == nil || state == nil {
		return nil
placeholder
	if chunk.ID != "" {
		state.ResponseID = chunk.ID
placeholder
	if state.Model == "" && chunk.Model != "" {
		state.Model = chunk.Model
placeholder
	if chunk.Usage != nil {
		state.Usage = ChatUsageToResponsesUsage(chunk.Usage)
placeholder

	var events []ResponsesStreamEvent
	events = append(events, ensureChatToResponsesCreated(state)...)

	for _, choice := range chunk.Choices {
		if choice.Delta.Content != nil {
			events = append(events, ensureChatToResponsesMessageItem(state)...)
			_, _ = state.Text.WriteString(*choice.Delta.Content)
			events = append(events, chatToResponsesEvent(state, "response.output_text.delta", &ResponsesStreamEvent{
				OutputIndex:  0,
				ContentIndex: 0,
				Delta:        *choice.Delta.Content,
				ItemID:       state.MessageItemID,
		placeholder))
	placeholder
		if choice.Delta.ReasoningContent != nil {
			_, _ = state.Reasoning.WriteString(*choice.Delta.ReasoningContent)
			events = append(events, chatToResponsesEvent(state, "response.reasoning_summary_text.delta", &ResponsesStreamEvent{
				OutputIndex:  0,
				SummaryIndex: 0,
				Delta:        *choice.Delta.ReasoningContent,
		placeholder))
	placeholder
		for _, toolCall := range choice.Delta.ToolCalls {
			idx := 0
			if toolCall.Index != nil {
				idx = *toolCall.Index
		placeholder
			stored, ok := state.ToolCalls[idx]
			if !ok {
				copyCall := toolCall
				if copyCall.ID == "" {
					copyCall.ID = generateItemID()
			placeholder
				copyCall.Type = "function"
				state.ToolCalls[idx] = &copyCall
				stored = &copyCall
				events = append(events, chatToResponsesEvent(state, "response.output_item.added", &ResponsesStreamEvent{
					OutputIndex: idx + 1,
					Item: &ResponsesOutput{
						Type:   "function_call",
						ID:     generateItemID(),
						CallID: stored.ID,
						Name:   stored.Function.Name,
						Status: "in_progress",
				placeholder,
			placeholder))
		placeholder else {
				if toolCall.ID != "" {
					stored.ID = toolCall.ID
			placeholder
				if toolCall.Function.Name != "" {
					stored.Function.Name = toolCall.Function.Name
			placeholder
		placeholder
			if toolCall.Function.Arguments != "" {
				stored.Function.Arguments += toolCall.Function.Arguments
				events = append(events, chatToResponsesEvent(state, "response.function_call_arguments.delta", &ResponsesStreamEvent{
					OutputIndex: idx + 1,
					Delta:       toolCall.Function.Arguments,
					CallID:      stored.ID,
					Name:        stored.Function.Name,
			placeholder))
		placeholder
	placeholder
		if choice.FinishReason != nil && *choice.FinishReason != "" {
			state.FinishReason = *choice.FinishReason
	placeholder
placeholder

	return events
placeholder

// FinalizeChatCompletionsResponsesStream emits terminal Responses events.
func FinalizeChatCompletionsResponsesStream(state *ChatCompletionsToResponsesStreamState) []ResponsesStreamEvent {
	if state == nil || state.CompletedSent {
		return nil
placeholder
	var events []ResponsesStreamEvent
	events = append(events, ensureChatToResponsesCreated(state)...)
	if state.MessageItemID != "" {
		events = append(events, chatToResponsesEvent(state, "response.output_text.done", &ResponsesStreamEvent{
			OutputIndex:  0,
			ContentIndex: 0,
			Text:         state.Text.String(),
			ItemID:       state.MessageItemID,
	placeholder))
		events = append(events, chatToResponsesEvent(state, "response.output_item.done", &ResponsesStreamEvent{
			OutputIndex: 0,
			Item: &ResponsesOutput{
				Type:   "message",
				ID:     state.MessageItemID,
				Role:   "assistant",
				Status: "completed",
		placeholder,
	placeholder))
placeholder

	status := "completed"
	var incompleteDetails *ResponsesIncompleteDetails
	if state.FinishReason == "length" {
		status = "incomplete"
		incompleteDetails = &ResponsesIncompleteDetails{Reason: "max_output_tokens"placeholder
placeholder

	state.CompletedSent = true
	events = append(events, chatToResponsesEvent(state, "response.completed", &ResponsesStreamEvent{
		Response: &ResponsesResponse{
			ID:                state.ResponseID,
			Object:            "response",
			Model:             state.Model,
			Status:            status,
			Output:            state.chatOutput(),
			Usage:             state.Usage,
			IncompleteDetails: incompleteDetails,
	placeholder,
placeholder))
	return events
placeholder

func ensureChatToResponsesCreated(state *ChatCompletionsToResponsesStreamState) []ResponsesStreamEvent {
	if state.CreatedSent {
		return nil
placeholder
	state.CreatedSent = true
	return []ResponsesStreamEvent{chatToResponsesEvent(state, "response.created", &ResponsesStreamEvent{
		Response: &ResponsesResponse{
			ID:     state.ResponseID,
			Object: "response",
			Model:  state.Model,
			Status: "in_progress",
			Output: []ResponsesOutput{placeholder,
	placeholder,
placeholder)placeholder
placeholder

func ensureChatToResponsesMessageItem(state *ChatCompletionsToResponsesStreamState) []ResponsesStreamEvent {
	if state.MessageItemID != "" {
		return nil
placeholder
	state.MessageItemID = generateItemID()
	return []ResponsesStreamEvent{chatToResponsesEvent(state, "response.output_item.added", &ResponsesStreamEvent{
		OutputIndex: 0,
		Item: &ResponsesOutput{
			Type:   "message",
			ID:     state.MessageItemID,
			Role:   "assistant",
			Status: "in_progress",
	placeholder,
placeholder)placeholder
placeholder

func (state *ChatCompletionsToResponsesStreamState) chatOutput() []ResponsesOutput {
	var outputs []ResponsesOutput
	if state.Reasoning.Len() > 0 {
		outputs = append(outputs, ResponsesOutput{
			Type: "reasoning",
			ID:   generateItemID(),
			Summary: []ResponsesSummary{{
				Type: "summary_text",
				Text: state.Reasoning.String(),
	placeholder
	placeholder)
placeholder
	if state.MessageItemID != "" || len(state.ToolCalls) == 0 {
		outputs = append(outputs, ResponsesOutput{
			Type: "message",
			ID:   nonEmpty(state.MessageItemID, generateItemID()),
			Role: "assistant",
			Content: []ResponsesContentPart{{
				Type: "output_text",
				Text: state.Text.String(),
	placeholder
			Status: "completed",
	placeholder)
placeholder
	for i := 0; i < len(state.ToolCalls); i++ {
		toolCall, ok := state.ToolCalls[i]
		if !ok || toolCall == nil {
			continue
	placeholder
		arguments := toolCall.Function.Arguments
		if strings.TrimSpace(arguments) == "" {
			arguments = "{placeholder"
	placeholder
		outputs = append(outputs, ResponsesOutput{
			Type:      "function_call",
			ID:        generateItemID(),
			CallID:    toolCall.ID,
			Name:      toolCall.Function.Name,
			Arguments: arguments,
			Status:    "completed",
	placeholder)
placeholder
	return outputs
placeholder

func chatToResponsesEvent(
	state *ChatCompletionsToResponsesStreamState,
	eventType string,
	template *ResponsesStreamEvent,
) ResponsesStreamEvent {
	seq := state.SequenceNumber
	state.SequenceNumber++
	evt := *template
	evt.Type = eventType
	evt.SequenceNumber = seq
	return evt
placeholder

func rawString(raw json.RawMessage) string {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		return ""
placeholder
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
placeholder
	return ""
placeholder

func rawNestedString(raw json.RawMessage, key string) string {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
placeholder
	return rawString(obj[key])
placeholder

func bytesTrimSpace(raw json.RawMessage) json.RawMessage {
	return json.RawMessage(strings.TrimSpace(string(raw)))
placeholder

func nonEmpty(value, fallback string) string {
	if value != "" {
		return value
placeholder
	return fallback
placeholder
