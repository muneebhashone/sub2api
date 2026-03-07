package apicompat

import (
	"encoding/json"
	"fmt"
	"time"
)

// ---------------------------------------------------------------------------
// Non-streaming: ResponsesResponse → AnthropicResponse
// ---------------------------------------------------------------------------

// ResponsesToAnthropic converts a Responses API response directly into an
// Anthropic Messages response. Reasoning output items are mapped to thinking
// blocks; function_call items become tool_use blocks.
func ResponsesToAnthropic(resp *ResponsesResponse, model string) *AnthropicResponse {
	out := &AnthropicResponse{
		ID:    resp.ID,
		Type:  "message",
		Role:  "assistant",
		Model: model,
placeholder

	var blocks []AnthropicContentBlock

	for _, item := range resp.Output {
		switch item.Type {
		case "reasoning":
			summaryText := ""
			for _, s := range item.Summary {
				if s.Type == "summary_text" && s.Text != "" {
					summaryText += s.Text
			placeholder
		placeholder
			if summaryText != "" {
				blocks = append(blocks, AnthropicContentBlock{
					Type:     "thinking",
					Thinking: summaryText,
			placeholder)
		placeholder
		case "message":
			for _, part := range item.Content {
				if part.Type == "output_text" && part.Text != "" {
					blocks = append(blocks, AnthropicContentBlock{
						Type: "text",
						Text: part.Text,
				placeholder)
			placeholder
		placeholder
		case "function_call":
			blocks = append(blocks, AnthropicContentBlock{
				Type:  "tool_use",
				ID:    fromResponsesCallID(item.CallID),
				Name:  item.Name,
				Input: json.RawMessage(item.Arguments),
		placeholder)
		case "web_search_call":
			toolUseID := "srvtoolu_" + item.ID
			query := ""
			if item.Action != nil {
				query = item.Action.Query
		placeholder
			inputJSON, _ := json.Marshal(map[string]string{"query": queryplaceholder)
			blocks = append(blocks, AnthropicContentBlock{
				Type:  "server_tool_use",
				ID:    toolUseID,
				Name:  "web_search",
				Input: inputJSON,
		placeholder)
			emptyResults, _ := json.Marshal([]struct{placeholder{placeholder)
			blocks = append(blocks, AnthropicContentBlock{
				Type:      "web_search_tool_result",
				ToolUseID: toolUseID,
				Content:   emptyResults,
		placeholder)
	placeholder
placeholder

	if len(blocks) == 0 {
		blocks = append(blocks, AnthropicContentBlock{Type: "text", Text: ""placeholder)
placeholder
	out.Content = blocks

	out.StopReason = responsesStatusToAnthropicStopReason(resp.Status, resp.IncompleteDetails, blocks)

	if resp.Usage != nil {
		out.Usage = AnthropicUsage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
	placeholder
		if resp.Usage.InputTokensDetails != nil {
			out.Usage.CacheReadInputTokens = resp.Usage.InputTokensDetails.CachedTokens
	placeholder
placeholder

	return out
placeholder

func responsesStatusToAnthropicStopReason(status string, details *ResponsesIncompleteDetails, blocks []AnthropicContentBlock) string {
	switch status {
	case "incomplete":
		if details != nil && details.Reason == "max_output_tokens" {
			return "max_tokens"
	placeholder
		return "end_turn"
	case "completed":
		if len(blocks) > 0 && blocks[len(blocks)-1].Type == "tool_use" {
			return "tool_use"
	placeholder
		return "end_turn"
	default:
		return "end_turn"
placeholder
placeholder

// ---------------------------------------------------------------------------
// Streaming: ResponsesStreamEvent → []AnthropicStreamEvent (stateful converter)
// ---------------------------------------------------------------------------

// ResponsesEventToAnthropicState tracks state for converting a sequence of
// Responses SSE events directly into Anthropic SSE events.
type ResponsesEventToAnthropicState struct {
	MessageStartSent bool
	MessageStopSent  bool

	ContentBlockIndex int
	ContentBlockOpen  bool
	CurrentBlockType  string // "text" | "thinking" | "tool_use"

	// OutputIndexToBlockIdx maps Responses output_index → Anthropic content block index.
	OutputIndexToBlockIdx map[int]int

	InputTokens          int
	OutputTokens         int
	CacheReadInputTokens int

	ResponseID string
	Model      string
	Created    int64
placeholder

// NewResponsesEventToAnthropicState returns an initialised stream state.
func NewResponsesEventToAnthropicState() *ResponsesEventToAnthropicState {
	return &ResponsesEventToAnthropicState{
		OutputIndexToBlockIdx: make(map[int]int),
		Created:               time.Now().Unix(),
placeholder
placeholder

// ResponsesEventToAnthropicEvents converts a single Responses SSE event into
// zero or more Anthropic SSE events, updating state as it goes.
func ResponsesEventToAnthropicEvents(
	evt *ResponsesStreamEvent,
	state *ResponsesEventToAnthropicState,
) []AnthropicStreamEvent {
	switch evt.Type {
	case "response.created":
		return resToAnthHandleCreated(evt, state)
	case "response.output_item.added":
		return resToAnthHandleOutputItemAdded(evt, state)
	case "response.output_text.delta":
		return resToAnthHandleTextDelta(evt, state)
	case "response.output_text.done":
		return resToAnthHandleBlockDone(state)
	case "response.function_call_arguments.delta":
		return resToAnthHandleFuncArgsDelta(evt, state)
	case "response.function_call_arguments.done":
		return resToAnthHandleBlockDone(state)
	case "response.output_item.done":
		return resToAnthHandleOutputItemDone(evt, state)
	case "response.reasoning_summary_text.delta":
		return resToAnthHandleReasoningDelta(evt, state)
	case "response.reasoning_summary_text.done":
		return resToAnthHandleBlockDone(state)
	case "response.completed", "response.incomplete", "response.failed":
		return resToAnthHandleCompleted(evt, state)
	default:
		return nil
placeholder
placeholder

// FinalizeResponsesAnthropicStream emits synthetic termination events if the
// stream ended without a proper completion event.
func FinalizeResponsesAnthropicStream(state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if !state.MessageStartSent || state.MessageStopSent {
		return nil
placeholder

	var events []AnthropicStreamEvent
	events = append(events, closeCurrentBlock(state)...)

	events = append(events,
		AnthropicStreamEvent{
			Type: "message_delta",
			Delta: &AnthropicDelta{
				StopReason: "end_turn",
		placeholder,
			Usage: &AnthropicUsage{
				InputTokens:          state.InputTokens,
				OutputTokens:         state.OutputTokens,
				CacheReadInputTokens: state.CacheReadInputTokens,
		placeholder,
	placeholder,
		AnthropicStreamEvent{Type: "message_stop"placeholder,
	)
	state.MessageStopSent = true
	return events
placeholder

// ResponsesAnthropicEventToSSE formats an AnthropicStreamEvent as an SSE line pair.
func ResponsesAnthropicEventToSSE(evt AnthropicStreamEvent) (string, error) {
	data, err := json.Marshal(evt)
	if err != nil {
		return "", err
placeholder
	return fmt.Sprintf("event: %s\ndata: %s\n\n", evt.Type, data), nil
placeholder

// --- internal handlers ---

func resToAnthHandleCreated(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if evt.Response != nil {
		state.ResponseID = evt.Response.ID
		// Only use upstream model if no override was set (e.g. originalModel)
		if state.Model == "" {
			state.Model = evt.Response.Model
	placeholder
placeholder

	if state.MessageStartSent {
		return nil
placeholder
	state.MessageStartSent = true

	return []AnthropicStreamEvent{{
		Type: "message_start",
		Message: &AnthropicResponse{
			ID:      state.ResponseID,
			Type:    "message",
			Role:    "assistant",
			Content: []AnthropicContentBlock{placeholder,
			Model:   state.Model,
			Usage: AnthropicUsage{
				InputTokens:  0,
				OutputTokens: 0,
		placeholder,
	placeholder,
placeholderplaceholder
placeholder

func resToAnthHandleOutputItemAdded(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if evt.Item == nil {
		return nil
placeholder

	switch evt.Item.Type {
	case "function_call":
		var events []AnthropicStreamEvent
		events = append(events, closeCurrentBlock(state)...)

		idx := state.ContentBlockIndex
		state.OutputIndexToBlockIdx[evt.OutputIndex] = idx
		state.ContentBlockOpen = true
		state.CurrentBlockType = "tool_use"

		events = append(events, AnthropicStreamEvent{
			Type:  "content_block_start",
			Index: &idx,
			ContentBlock: &AnthropicContentBlock{
				Type:  "tool_use",
				ID:    fromResponsesCallID(evt.Item.CallID),
				Name:  evt.Item.Name,
				Input: json.RawMessage("{placeholder"),
		placeholder,
	placeholder)
		return events

	case "reasoning":
		var events []AnthropicStreamEvent
		events = append(events, closeCurrentBlock(state)...)

		idx := state.ContentBlockIndex
		state.OutputIndexToBlockIdx[evt.OutputIndex] = idx
		state.ContentBlockOpen = true
		state.CurrentBlockType = "thinking"

		events = append(events, AnthropicStreamEvent{
			Type:  "content_block_start",
			Index: &idx,
			ContentBlock: &AnthropicContentBlock{
				Type:     "thinking",
				Thinking: "",
		placeholder,
	placeholder)
		return events

	case "message":
		return nil
placeholder

	return nil
placeholder

func resToAnthHandleTextDelta(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if evt.Delta == "" {
		return nil
placeholder

	var events []AnthropicStreamEvent

	if !state.ContentBlockOpen || state.CurrentBlockType != "text" {
		events = append(events, closeCurrentBlock(state)...)

		idx := state.ContentBlockIndex
		state.ContentBlockOpen = true
		state.CurrentBlockType = "text"

		events = append(events, AnthropicStreamEvent{
			Type:  "content_block_start",
			Index: &idx,
			ContentBlock: &AnthropicContentBlock{
				Type: "text",
				Text: "",
		placeholder,
	placeholder)
placeholder

	idx := state.ContentBlockIndex
	events = append(events, AnthropicStreamEvent{
		Type:  "content_block_delta",
		Index: &idx,
		Delta: &AnthropicDelta{
			Type: "text_delta",
			Text: evt.Delta,
	placeholder,
placeholder)
	return events
placeholder

func resToAnthHandleFuncArgsDelta(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if evt.Delta == "" {
		return nil
placeholder

	blockIdx, ok := state.OutputIndexToBlockIdx[evt.OutputIndex]
	if !ok {
		return nil
placeholder

	return []AnthropicStreamEvent{{
		Type:  "content_block_delta",
		Index: &blockIdx,
		Delta: &AnthropicDelta{
			Type:        "input_json_delta",
			PartialJSON: evt.Delta,
	placeholder,
placeholderplaceholder
placeholder

func resToAnthHandleReasoningDelta(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if evt.Delta == "" {
		return nil
placeholder

	blockIdx, ok := state.OutputIndexToBlockIdx[evt.OutputIndex]
	if !ok {
		return nil
placeholder

	return []AnthropicStreamEvent{{
		Type:  "content_block_delta",
		Index: &blockIdx,
		Delta: &AnthropicDelta{
			Type:     "thinking_delta",
			Thinking: evt.Delta,
	placeholder,
placeholderplaceholder
placeholder

func resToAnthHandleBlockDone(state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if !state.ContentBlockOpen {
		return nil
placeholder
	return closeCurrentBlock(state)
placeholder

func resToAnthHandleOutputItemDone(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if evt.Item == nil {
		return nil
placeholder

	// Handle web_search_call → synthesize server_tool_use + web_search_tool_result blocks.
	if evt.Item.Type == "web_search_call" && evt.Item.Status == "completed" {
		return resToAnthHandleWebSearchDone(evt, state)
placeholder

	if state.ContentBlockOpen {
		return closeCurrentBlock(state)
placeholder
	return nil
placeholder

// resToAnthHandleWebSearchDone converts an OpenAI web_search_call output item
// into Anthropic server_tool_use + web_search_tool_result content block pairs.
// This allows Claude Code to count the searches performed.
func resToAnthHandleWebSearchDone(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	var events []AnthropicStreamEvent
	events = append(events, closeCurrentBlock(state)...)

	toolUseID := "srvtoolu_" + evt.Item.ID
	query := ""
	if evt.Item.Action != nil {
		query = evt.Item.Action.Query
placeholder
	inputJSON, _ := json.Marshal(map[string]string{"query": queryplaceholder)

	// Emit server_tool_use block (start + stop).
	idx1 := state.ContentBlockIndex
	events = append(events, AnthropicStreamEvent{
		Type:  "content_block_start",
		Index: &idx1,
		ContentBlock: &AnthropicContentBlock{
			Type:  "server_tool_use",
			ID:    toolUseID,
			Name:  "web_search",
			Input: inputJSON,
	placeholder,
placeholder)
	events = append(events, AnthropicStreamEvent{
		Type:  "content_block_stop",
		Index: &idx1,
placeholder)
	state.ContentBlockIndex++

	// Emit web_search_tool_result block (start + stop).
	// Content is empty because OpenAI does not expose individual search results;
	// the model consumes them internally and produces text output.
	emptyResults, _ := json.Marshal([]struct{placeholder{placeholder)
	idx2 := state.ContentBlockIndex
	events = append(events, AnthropicStreamEvent{
		Type:  "content_block_start",
		Index: &idx2,
		ContentBlock: &AnthropicContentBlock{
			Type:      "web_search_tool_result",
			ToolUseID: toolUseID,
			Content:   emptyResults,
	placeholder,
placeholder)
	events = append(events, AnthropicStreamEvent{
		Type:  "content_block_stop",
		Index: &idx2,
placeholder)
	state.ContentBlockIndex++

	return events
placeholder

func resToAnthHandleCompleted(evt *ResponsesStreamEvent, state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if state.MessageStopSent {
		return nil
placeholder

	var events []AnthropicStreamEvent
	events = append(events, closeCurrentBlock(state)...)

	stopReason := "end_turn"
	if evt.Response != nil {
		if evt.Response.Usage != nil {
			state.InputTokens = evt.Response.Usage.InputTokens
			state.OutputTokens = evt.Response.Usage.OutputTokens
			if evt.Response.Usage.InputTokensDetails != nil {
				state.CacheReadInputTokens = evt.Response.Usage.InputTokensDetails.CachedTokens
		placeholder
	placeholder
		switch evt.Response.Status {
		case "incomplete":
			if evt.Response.IncompleteDetails != nil && evt.Response.IncompleteDetails.Reason == "max_output_tokens" {
				stopReason = "max_tokens"
		placeholder
		case "completed":
			if state.ContentBlockIndex > 0 && state.CurrentBlockType == "tool_use" {
				stopReason = "tool_use"
		placeholder
	placeholder
placeholder

	events = append(events,
		AnthropicStreamEvent{
			Type: "message_delta",
			Delta: &AnthropicDelta{
				StopReason: stopReason,
		placeholder,
			Usage: &AnthropicUsage{
				InputTokens:          state.InputTokens,
				OutputTokens:         state.OutputTokens,
				CacheReadInputTokens: state.CacheReadInputTokens,
		placeholder,
	placeholder,
		AnthropicStreamEvent{Type: "message_stop"placeholder,
	)
	state.MessageStopSent = true
	return events
placeholder

func closeCurrentBlock(state *ResponsesEventToAnthropicState) []AnthropicStreamEvent {
	if !state.ContentBlockOpen {
		return nil
placeholder
	idx := state.ContentBlockIndex
	state.ContentBlockOpen = false
	state.ContentBlockIndex++
	return []AnthropicStreamEvent{{
		Type:  "content_block_stop",
		Index: &idx,
placeholderplaceholder
placeholder
