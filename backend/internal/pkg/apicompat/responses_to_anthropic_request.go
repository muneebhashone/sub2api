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
	system, messages, err := convertResponsesInputToAnthropic(req.Instructions, req.Input)
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
// a Responses API instructions + input array. Returns the system as raw JSON
// (for Anthropic's polymorphic system field) and a list of Anthropic messages.
func convertResponsesInputToAnthropic(instructions string, inputRaw json.RawMessage) (json.RawMessage, []AnthropicMessage, error) {
	var systemParts []string
	if strings.TrimSpace(instructions) != "" {
		systemParts = append(systemParts, strings.TrimSpace(instructions))
placeholder

	// Try as plain string input.
	var inputStr string
	if err := json.Unmarshal(inputRaw, &inputStr); err == nil {
		content, _ := json.Marshal(inputStr)
		var system json.RawMessage
		if len(systemParts) > 0 {
			system, _ = json.Marshal(strings.Join(systemParts, "\n\n"))
	placeholder
		return system, []AnthropicMessage{{Role: "user", Content: contentplaceholderplaceholder, nil
placeholder

	var items []ResponsesInputItem
	if err := json.Unmarshal(inputRaw, &items); err != nil {
		return nil, nil, fmt.Errorf("parse responses input: %w", err)
placeholder

	var messages []AnthropicMessage

	for _, item := range items {
		switch {
		case item.Role == "system" || item.Role == "developer":
			text := extractTextFromContent(item.Content)
			if text != "" {
				systemParts = append(systemParts, text)
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
			contentJSON := responsesFunctionOutputToAnthropicContent(item)
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

		case item.Type == "reasoning":
			// Anthropic 无法摄入 OpenAI 的 reasoning：encrypted_content 是不透明的，
			// 而 thinking 块的重放需要 Anthropic 自己签发的 signature，无法伪造。
			// Codex 常见形态（只带 summary + encrypted_content）本来就会被丢弃，
			// 这里让带 content 数组的形态保持同样行为——否则 reasoning_text 块会被
			// 原样塞进 Anthropic 请求体，上游直接回 400。

		case item.Role == "user":
			content, err := convertResponsesUserToAnthropicContent(item.Content)
			if err != nil {
				return nil, nil, err
		placeholder
			// 内容里只有网关不认识的分片时，sanitize 会得到空串。Anthropic 拒收
			// 空内容消息（"all messages must have non-empty content"），整条丢掉
			// 比发一条必然 400 的消息更可用。
			if anthropicContentIsEmpty(content) {
				continue
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
			// 同上：分片全不认识时会退化成单个空 text 块，而 Anthropic 拒收
			// 空文本块（"text content blocks must contain non-whitespace text"）。
			if anthropicContentIsEmpty(content) || anthropicContentIsOnlyBlankText(content) {
				continue
		placeholder
			messages = append(messages, AnthropicMessage{
				Role:    "assistant",
				Content: content,
		placeholder)

		default:
			// 未知 role/type —— 尽量当作 user 消息保留其中的文本/图片。
			// 必须走与真实 user 消息同一套白名单转换：直接透传 item.Content 会把
			// Responses 专有的分片类型（reasoning_text、web_search_call 的载荷等）
			// 原样发给 Anthropic，上游只会回 400 把整轮打挂。
			if item.Content == nil {
				continue
		placeholder
			content, err := convertResponsesUserToAnthropicContent(item.Content)
			if err != nil {
				return nil, nil, err
		placeholder
			if anthropicContentIsEmpty(content) {
				continue
		placeholder
			messages = append(messages, AnthropicMessage{
				Role:    "user",
				Content: content,
		placeholder)
	placeholder
placeholder

	// Repair tool_use/tool_result pairing, then merge consecutive same-role
	// messages (Anthropic requires alternating roles). The first merge groups
	// parallel calls (and their results) so the pairing pass sees them together;
	// the pairing pass may re-split a user turn (e.g. when an injected message
	// sat between a call and its output), so a second merge restores alternation.
	messages = mergeConsecutiveMessages(messages)
	messages = normalizeAnthropicToolPairing(messages)
	messages = mergeConsecutiveMessages(messages)

	var system json.RawMessage
	if len(systemParts) > 0 {
		system, _ = json.Marshal(strings.Join(systemParts, "\n\n"))
placeholder

	return system, messages, nil
placeholder

func responsesFunctionOutputToAnthropicContent(item ResponsesInputItem) json.RawMessage {
	if len(item.outputRaw) == 0 {
		output := item.Output
		if output == "" {
			output = "(empty)"
	placeholder
		content, _ := json.Marshal(output)
		return content
placeholder

	var parts []ResponsesContentPart
	if err := json.Unmarshal(item.outputRaw, &parts); err == nil {
		blocks := make([]AnthropicContentBlock, 0, len(parts))
		for _, part := range parts {
			switch part.Type {
			case "input_text", "output_text", "text":
				if part.Text != "" {
					blocks = append(blocks, AnthropicContentBlock{Type: "text", Text: part.Textplaceholder)
			placeholder
			case "input_image":
				if source := dataURIToAnthropicImageSource(part.ImageURL); source != nil {
					blocks = append(blocks, AnthropicContentBlock{Type: "image", Source: sourceplaceholder)
			placeholder
		placeholder
	placeholder
		if len(blocks) > 0 {
			content, _ := json.Marshal(blocks)
			return content
	placeholder
		if len(parts) == 0 {
			content, _ := json.Marshal("(empty)")
			return content
	placeholder
placeholder

	content, _ := json.Marshal(item.Output)
	return content
placeholder

// normalizeAnthropicToolPairing rebuilds the message sequence so it satisfies
// Anthropic's tool_use/tool_result invariants, which the naive item-by-item
// conversion violates whenever the Responses history interleaves anything
// between a function_call and its function_call_output:
//
//   - every tool_result block must have a matching tool_use in the immediately
//     preceding assistant message ("tool_result ... must have a corresponding
//     tool_use block in the previous message");
//   - every tool_use block must be answered by a tool_result in the immediately
//     following user message (Anthropic rejects unanswered tool_use ids);
//   - user/assistant turns must alternate.
//
// codex (Responses, store:false) re-sends the whole history each turn and
// frequently injects items between a call and its output — a developer/approval
// notice, or a sibling parallel call whose output never arrived. The unrepaired
// converter emits each function_call as its own assistant message and each
// output as its own user message, so any such interleaving breaks
// tool_use↔tool_result adjacency and yields an upstream 400.
//
// The repair indexes every tool_result by its tool_use id, then for each
// assistant message carrying tool_use blocks keeps only the answered ones
// (dropping unanswered/dangling calls — and the assistant message entirely if it
// has no other content) and emits the matching tool_result blocks, in call
// order, as the very next user message. Standalone tool_result blocks are
// dropped from their original position (re-emitted adjacent to their call);
// orphan tool_results with no announcing tool_use are dropped. Non-tool content
// passes through in place. This mirrors normalizeChatMessages on the
// Responses→Chat path.
func normalizeAnthropicToolPairing(messages []AnthropicMessage) []AnthropicMessage {
	// Index every tool_result block by its tool_use id (last wins on dup).
	results := make(map[string]AnthropicContentBlock)
	for _, m := range messages {
		if m.Role != "user" {
			continue
	placeholder
		for _, b := range parseContentBlocks(m.Content) {
			if b.Type == "tool_result" && b.ToolUseID != "" {
				results[b.ToolUseID] = b
		placeholder
	placeholder
placeholder

	out := make([]AnthropicMessage, 0, len(messages))
	for _, m := range messages {
		blocks := parseContentBlocks(m.Content)
		switch m.Role {
		case "assistant":
			var toolUses, others []AnthropicContentBlock
			for _, b := range blocks {
				if b.Type == "tool_use" {
					toolUses = append(toolUses, b)
			placeholder else {
					others = append(others, b)
			placeholder
		placeholder
			if len(toolUses) == 0 {
				out = append(out, m)
				continue
		placeholder
			kept := make([]AnthropicContentBlock, 0, len(toolUses))
			for _, tu := range toolUses {
				if _, ok := results[tu.ID]; ok {
					kept = append(kept, tu)
			placeholder
		placeholder
			if len(kept) == 0 {
				// No answered calls: keep any non-tool content, else drop.
				if len(others) > 0 {
					out = append(out, anthropicMessageFromBlocks("assistant", others))
			placeholder
				continue
		placeholder
			asstBlocks := make([]AnthropicContentBlock, 0, len(others)+len(kept))
			asstBlocks = append(asstBlocks, others...)
			asstBlocks = append(asstBlocks, kept...)
			out = append(out, anthropicMessageFromBlocks("assistant", asstBlocks))

			resBlocks := make([]AnthropicContentBlock, 0, len(kept))
			for _, tu := range kept {
				resBlocks = append(resBlocks, results[tu.ID])
		placeholder
			out = append(out, anthropicMessageFromBlocks("user", resBlocks))

		case "user":
			var nonResult []AnthropicContentBlock
			hasResult := false
			for _, b := range blocks {
				if b.Type == "tool_result" {
					hasResult = true
					continue
			placeholder
				nonResult = append(nonResult, b)
		placeholder
			if !hasResult {
				out = append(out, m)
				continue
		placeholder
			// The tool_result blocks are re-emitted next to their call; keep any
			// other content of this user turn in place, drop it if there is none.
			if len(nonResult) > 0 {
				out = append(out, anthropicMessageFromBlocks("user", nonResult))
		placeholder

		default:
			out = append(out, m)
	placeholder
placeholder
	return out
placeholder

// anthropicMessageFromBlocks builds an AnthropicMessage whose content is the
// marshaled block array.
func anthropicMessageFromBlocks(role string, blocks []AnthropicContentBlock) AnthropicMessage {
	content, _ := json.Marshal(blocks)
	return AnthropicMessage{Role: role, Content: contentplaceholder
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
// anthropicContentIsEmpty 判断转换结果是否为"空内容"。
// convertResponsesUserToAnthropicContent 在没有任何可识别分片时返回 JSON 空串，
// 而 Anthropic 拒收空内容消息。
func anthropicContentIsEmpty(content json.RawMessage) bool {
	trimmed := strings.TrimSpace(string(content))
	switch trimmed {
	case "", "null", `""`, "[]":
		return true
placeholder
	return false
placeholder

// anthropicContentIsOnlyBlankText 判断内容是否只由空白 text 块组成。
func anthropicContentIsOnlyBlankText(content json.RawMessage) bool {
	blocks := parseContentBlocks(content)
	if len(blocks) == 0 {
		return false
placeholder
	for _, b := range blocks {
		if b.Type != "text" || strings.TrimSpace(b.Text) != "" {
			return false
	placeholder
placeholder
	return true
placeholder

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
		case "web_search", "google_search", "web_search_20250305":
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
		case "custom":
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
				InputSchema: normalizeAnthropicInputSchema(t.Parameters),
		placeholder)
	placeholder
placeholder
	return out
placeholder

// normalizeAnthropicInputSchema ensures input_schema is a valid object schema.
func normalizeAnthropicInputSchema(schema json.RawMessage) json.RawMessage {
	const emptyObjectSchema = `{"type":"object","properties":{placeholderplaceholder`

	trimmed := strings.TrimSpace(string(schema))
	if trimmed == "" || trimmed == "null" {
		return json.RawMessage(emptyObjectSchema)
placeholder

	var m map[string]json.RawMessage
	if err := json.Unmarshal(schema, &m); err != nil {
		return json.RawMessage(`{"type":"object","properties":{placeholderplaceholder`)
placeholder

	typeRaw, ok := m["type"]
	if !ok || strings.TrimSpace(string(typeRaw)) == "" || string(typeRaw) == "null" {
		m["type"] = json.RawMessage(`"object"`)
placeholder else {
		var typ string
		if err := json.Unmarshal(typeRaw, &typ); err != nil || typ != "object" {
			return json.RawMessage(emptyObjectSchema)
	placeholder
placeholder

	if _, ok := m["properties"]; !ok {
		m["properties"] = json.RawMessage(`{placeholder`)
placeholder

	out, err := json.Marshal(m)
	if err != nil {
		return json.RawMessage(emptyObjectSchema)
placeholder
	return out
placeholder

// convertResponsesToAnthropicToolChoice maps Responses tool_choice to Anthropic format.
// Reverse of convertAnthropicToolChoiceToResponses.
//
//	"auto"                                     → {"type":"auto"placeholder
//	"required"                                 → {"type":"any"placeholder
//	"none"                                     → {"type":"none"placeholder
//	{"type":"function","name":"X"placeholder                 → {"type":"tool","name":"X"placeholder
//	{"type":"function","function":{"name":"X"placeholderplaceholder     → {"type":"tool","name":"X"placeholder // legacy
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
		Name     string `json:"name"`
		Function struct {
			Name string `json:"name"`
	placeholder `json:"function"`
placeholder
	if err := json.Unmarshal(raw, &tc); err == nil && tc.Type == "function" {
		name := strings.TrimSpace(tc.Name)
		if name == "" {
			name = strings.TrimSpace(tc.Function.Name)
	placeholder
		if name == "" {
			return raw, nil
	placeholder
		return json.Marshal(map[string]string{
			"type": "tool",
			"name": name,
	placeholder)
placeholder

	// Pass through unknown
	return raw, nil
placeholder
