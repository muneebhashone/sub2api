package service

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// ConvertChatCompletionsToResponses converts an OpenAI Chat Completions request to a Responses request.
func ConvertChatCompletionsToResponses(req map[string]any) (map[string]any, error) {
	if req == nil {
		return nil, errors.New("request is nil")
placeholder

	model := strings.TrimSpace(getString(req["model"]))
	if model == "" {
		return nil, errors.New("model is required")
placeholder

	messagesRaw, ok := req["messages"]
	if !ok {
		return nil, errors.New("messages is required")
placeholder
	messages, ok := messagesRaw.([]any)
	if !ok {
		return nil, errors.New("messages must be an array")
placeholder

	input, err := convertChatMessagesToResponsesInput(messages)
	if err != nil {
		return nil, err
placeholder

	out := make(map[string]any, len(req)+1)
	for key, value := range req {
		switch key {
		case "messages", "max_tokens", "max_completion_tokens", "stream_options", "functions", "function_call":
			continue
		default:
			out[key] = value
	placeholder
placeholder

	out["model"] = model
	out["input"] = input

	if _, ok := out["max_output_tokens"]; !ok {
		if v, ok := req["max_tokens"]; ok {
			out["max_output_tokens"] = v
	placeholder else if v, ok := req["max_completion_tokens"]; ok {
			out["max_output_tokens"] = v
	placeholder
placeholder

	if _, ok := out["tools"]; !ok {
		if functions, ok := req["functions"].([]any); ok && len(functions) > 0 {
			tools := make([]any, 0, len(functions))
			for _, fn := range functions {
				if fnMap, ok := fn.(map[string]any); ok {
					tools = append(tools, map[string]any{
						"type":     "function",
						"function": fnMap,
				placeholder)
			placeholder
		placeholder
			if len(tools) > 0 {
				out["tools"] = tools
		placeholder
	placeholder
placeholder

	if _, ok := out["tool_choice"]; !ok {
		if functionCall, ok := req["function_call"]; ok {
			out["tool_choice"] = functionCall
	placeholder
placeholder

	return out, nil
placeholder

// ConvertResponsesToChatCompletion converts an OpenAI Responses response body to Chat Completions format.
func ConvertResponsesToChatCompletion(body []byte) ([]byte, error) {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
placeholder

	id := strings.TrimSpace(getString(resp["id"]))
	if id == "" {
		id = "chatcmpl-" + safeRandomHex(12)
placeholder
	model := strings.TrimSpace(getString(resp["model"]))

	created := getInt64(resp["created_at"])
	if created == 0 {
		created = getInt64(resp["created"])
placeholder
	if created == 0 {
		created = time.Now().Unix()
placeholder

	text, toolCalls := extractResponseTextAndToolCalls(resp)
	finishReason := "stop"
	if len(toolCalls) > 0 {
		finishReason = "tool_calls"
placeholder

	message := map[string]any{
		"role":    "assistant",
		"content": text,
placeholder
	if len(toolCalls) > 0 {
		message["tool_calls"] = toolCalls
placeholder

	chatResp := map[string]any{
		"id":      id,
		"object":  "chat.completion",
		"created": created,
		"model":   model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"message":       message,
				"finish_reason": finishReason,
		placeholder,
	placeholder,
placeholder

	if usage := extractResponseUsage(resp); usage != nil {
		chatResp["usage"] = usage
placeholder
	if fingerprint := strings.TrimSpace(getString(resp["system_fingerprint"])); fingerprint != "" {
		chatResp["system_fingerprint"] = fingerprint
placeholder

	return json.Marshal(chatResp)
placeholder

func convertChatMessagesToResponsesInput(messages []any) ([]any, error) {
	input := make([]any, 0, len(messages))
	for _, msg := range messages {
		msgMap, ok := msg.(map[string]any)
		if !ok {
			return nil, errors.New("message must be an object")
	placeholder
		role := strings.TrimSpace(getString(msgMap["role"]))
		if role == "" {
			return nil, errors.New("message role is required")
	placeholder

		switch role {
		case "tool":
			callID := strings.TrimSpace(getString(msgMap["tool_call_id"]))
			if callID == "" {
				callID = strings.TrimSpace(getString(msgMap["id"]))
		placeholder
			output := extractMessageContentText(msgMap["content"])
			input = append(input, map[string]any{
				"type":    "function_call_output",
				"call_id": callID,
				"output":  output,
		placeholder)
		case "function":
			callID := strings.TrimSpace(getString(msgMap["name"]))
			output := extractMessageContentText(msgMap["content"])
			input = append(input, map[string]any{
				"type":    "function_call_output",
				"call_id": callID,
				"output":  output,
		placeholder)
		default:
			convertedContent := convertChatContent(msgMap["content"])
			toolCalls := []any(nil)
			if role == "assistant" {
				toolCalls = extractToolCallsFromMessage(msgMap)
		placeholder
			skipAssistantMessage := role == "assistant" && len(toolCalls) > 0 && isEmptyContent(convertedContent)
			if !skipAssistantMessage {
				msgItem := map[string]any{
					"role":    role,
					"content": convertedContent,
			placeholder
				if name := strings.TrimSpace(getString(msgMap["name"])); name != "" {
					msgItem["name"] = name
			placeholder
				input = append(input, msgItem)
		placeholder
			if role == "assistant" && len(toolCalls) > 0 {
				input = append(input, toolCalls...)
		placeholder
	placeholder
placeholder
	return input, nil
placeholder

func convertChatContent(content any) any {
	switch v := content.(type) {
	case nil:
		return ""
	case string:
		return v
	case []any:
		converted := make([]any, 0, len(v))
		for _, part := range v {
			partMap, ok := part.(map[string]any)
			if !ok {
				converted = append(converted, part)
				continue
		placeholder
			partType := strings.TrimSpace(getString(partMap["type"]))
			switch partType {
			case "text":
				text := getString(partMap["text"])
				if text != "" {
					converted = append(converted, map[string]any{
						"type": "input_text",
						"text": text,
				placeholder)
					continue
			placeholder
			case "image_url":
				imageURL := ""
				if imageObj, ok := partMap["image_url"].(map[string]any); ok {
					imageURL = getString(imageObj["url"])
			placeholder else {
					imageURL = getString(partMap["image_url"])
			placeholder
				if imageURL != "" {
					converted = append(converted, map[string]any{
						"type":      "input_image",
						"image_url": imageURL,
				placeholder)
					continue
			placeholder
			case "input_text", "input_image":
				converted = append(converted, partMap)
				continue
		placeholder
			converted = append(converted, partMap)
	placeholder
		return converted
	default:
		return v
placeholder
placeholder

func extractToolCallsFromMessage(msg map[string]any) []any {
	var out []any
	if toolCalls, ok := msg["tool_calls"].([]any); ok {
		for _, call := range toolCalls {
			callMap, ok := call.(map[string]any)
			if !ok {
				continue
		placeholder
			callID := strings.TrimSpace(getString(callMap["id"]))
			if callID == "" {
				callID = strings.TrimSpace(getString(callMap["call_id"]))
		placeholder
			name := ""
			args := ""
			if fn, ok := callMap["function"].(map[string]any); ok {
				name = strings.TrimSpace(getString(fn["name"]))
				args = getString(fn["arguments"])
		placeholder
			if name == "" && args == "" {
				continue
		placeholder
			item := map[string]any{
				"type": "tool_call",
		placeholder
			if callID != "" {
				item["call_id"] = callID
		placeholder
			if name != "" {
				item["name"] = name
		placeholder
			if args != "" {
				item["arguments"] = args
		placeholder
			out = append(out, item)
	placeholder
placeholder

	if fnCall, ok := msg["function_call"].(map[string]any); ok {
		name := strings.TrimSpace(getString(fnCall["name"]))
		args := getString(fnCall["arguments"])
		if name != "" || args != "" {
			callID := strings.TrimSpace(getString(msg["tool_call_id"]))
			if callID == "" {
				callID = name
		placeholder
			item := map[string]any{
				"type": "function_call",
		placeholder
			if callID != "" {
				item["call_id"] = callID
		placeholder
			if name != "" {
				item["name"] = name
		placeholder
			if args != "" {
				item["arguments"] = args
		placeholder
			out = append(out, item)
	placeholder
placeholder

	return out
placeholder

func extractMessageContentText(content any) string {
	switch v := content.(type) {
	case string:
		return v
	case []any:
		parts := make([]string, 0, len(v))
		for _, part := range v {
			partMap, ok := part.(map[string]any)
			if !ok {
				continue
		placeholder
			partType := strings.TrimSpace(getString(partMap["type"]))
			if partType == "" || partType == "text" || partType == "output_text" || partType == "input_text" {
				text := getString(partMap["text"])
				if text != "" {
					parts = append(parts, text)
			placeholder
		placeholder
	placeholder
		return strings.Join(parts, "")
	default:
		return ""
placeholder
placeholder

func isEmptyContent(content any) bool {
	switch v := content.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(v) == ""
	case []any:
		return len(v) == 0
	default:
		return false
placeholder
placeholder

func extractResponseTextAndToolCalls(resp map[string]any) (string, []any) {
	output, ok := resp["output"].([]any)
	if !ok {
		if text, ok := resp["output_text"].(string); ok {
			return text, nil
	placeholder
		return "", nil
placeholder

	textParts := make([]string, 0)
	toolCalls := make([]any, 0)

	for _, item := range output {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType := strings.TrimSpace(getString(itemMap["type"]))

		if itemType == "tool_call" || itemType == "function_call" {
			if tc := responseItemToChatToolCall(itemMap); tc != nil {
				toolCalls = append(toolCalls, tc)
		placeholder
			continue
	placeholder

		content := itemMap["content"]
		switch v := content.(type) {
		case string:
			if v != "" {
				textParts = append(textParts, v)
		placeholder
		case []any:
			for _, part := range v {
				partMap, ok := part.(map[string]any)
				if !ok {
					continue
			placeholder
				partType := strings.TrimSpace(getString(partMap["type"]))
				switch partType {
				case "output_text", "text", "input_text":
					text := getString(partMap["text"])
					if text != "" {
						textParts = append(textParts, text)
				placeholder
				case "tool_call", "function_call":
					if tc := responseItemToChatToolCall(partMap); tc != nil {
						toolCalls = append(toolCalls, tc)
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	return strings.Join(textParts, ""), toolCalls
placeholder

func responseItemToChatToolCall(item map[string]any) map[string]any {
	callID := strings.TrimSpace(getString(item["call_id"]))
	if callID == "" {
		callID = strings.TrimSpace(getString(item["id"]))
placeholder
	name := strings.TrimSpace(getString(item["name"]))
	arguments := getString(item["arguments"])
	if fn, ok := item["function"].(map[string]any); ok {
		if name == "" {
			name = strings.TrimSpace(getString(fn["name"]))
	placeholder
		if arguments == "" {
			arguments = getString(fn["arguments"])
	placeholder
placeholder

	if name == "" && arguments == "" && callID == "" {
		return nil
placeholder

	if callID == "" {
		callID = "call_" + safeRandomHex(6)
placeholder

	return map[string]any{
		"id":   callID,
		"type": "function",
		"function": map[string]any{
			"name":      name,
			"arguments": arguments,
	placeholder,
placeholder
placeholder

func extractResponseUsage(resp map[string]any) map[string]any {
	usage, ok := resp["usage"].(map[string]any)
	if !ok {
		return nil
placeholder
	promptTokens := int(getNumber(usage["input_tokens"]))
	completionTokens := int(getNumber(usage["output_tokens"]))
	if promptTokens == 0 && completionTokens == 0 {
		return nil
placeholder

	return map[string]any{
		"prompt_tokens":     promptTokens,
		"completion_tokens": completionTokens,
		"total_tokens":      promptTokens + completionTokens,
placeholder
placeholder

func getString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case json.Number:
		return v.String()
	default:
		return ""
placeholder
placeholder

func getNumber(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		f, _ := v.Float64()
		return f
	default:
		return 0
placeholder
placeholder

func getInt64(value any) int64 {
	switch v := value.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	case json.Number:
		i, _ := v.Int64()
		return i
	default:
		return 0
placeholder
placeholder

func safeRandomHex(byteLength int) string {
	value, err := randomHexString(byteLength)
	if err != nil || value == "" {
		return "000000"
placeholder
	return value
placeholder
