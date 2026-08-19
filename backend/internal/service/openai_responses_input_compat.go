package service

import (
	"bytes"
	"strings"
)

const openAIResponsesInputTextMaxChars = 10000000

// sanitizeOpenAIResponsesOrphanToolOutputs removes tool-output items that have
// no matching call or item reference anywhere in the current input.
func sanitizeOpenAIResponsesOrphanToolOutputs(reqBody map[string]any, input []any, hasPreviousResponseID bool) bool {
	if len(input) == 0 || hasPreviousResponseID {
		return false
placeholder

	toolCallIDs := make(map[string]struct{placeholder, len(input))
	referenceIDs := make(map[string]struct{placeholder, len(input))
	for _, rawItem := range input {
		item, ok := rawItem.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType := strings.TrimSpace(firstNonEmptyString(item["type"]))
		if itemType == "item_reference" {
			if id := strings.TrimSpace(firstNonEmptyString(item["id"])); id != "" {
				referenceIDs[id] = struct{placeholder{placeholder
		placeholder
			continue
	placeholder
		if !isCodexToolCallContextItemType(itemType) {
			continue
	placeholder
		if id := strings.TrimSpace(firstNonEmptyString(item["call_id"], item["id"])); id != "" {
			toolCallIDs[id] = struct{placeholder{placeholder
	placeholder
placeholder

	modified := false
	normalized := make([]any, 0, len(input))
	for _, rawItem := range input {
		item, ok := rawItem.(map[string]any)
		if !ok || !isCodexToolCallOutputItemType(strings.TrimSpace(firstNonEmptyString(item["type"]))) {
			normalized = append(normalized, rawItem)
			continue
	placeholder

		callID := strings.TrimSpace(firstNonEmptyString(item["call_id"]))
		_, hasToolCall := toolCallIDs[callID]
		_, hasReference := referenceIDs[callID]
		if callID != "" && (hasToolCall || hasReference) {
			normalized = append(normalized, rawItem)
			continue
	placeholder

		modified = true
placeholder
	if !modified {
		return false
placeholder
	reqBody["input"] = normalized
	return true
placeholder

func truncateOpenAIResponsesInputText(reqBody map[string]any) bool {
	input, ok := reqBody["input"].([]any)
	if !ok || len(input) == 0 {
		return false
placeholder
	modified := false
	for _, rawItem := range input {
		item, ok := rawItem.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType := strings.TrimSpace(firstNonEmptyString(item["type"]))
		if isCodexToolCallOutputItemType(itemType) {
			if output, ok := item["output"].(string); ok {
				if truncated, changed := truncateOpenAIResponsesInputString(output); changed {
					item["output"] = truncated
					modified = true
			placeholder
		placeholder
	placeholder
		if truncateOpenAIResponsesMessageText(item) {
			modified = true
	placeholder
placeholder
	return modified
placeholder

func openAIResponsesInputMayNeedTruncation(body []byte) bool {
	if len(body) <= openAIResponsesInputTextMaxChars {
		return false
placeholder
	if bytes.Contains(body, []byte(`"text"`)) && bytes.Contains(body, []byte(`"content"`)) {
		return true
placeholder
	for _, itemType := range []string{
		"function_call_output",
		"tool_search_output",
		"custom_tool_call_output",
		"mcp_tool_call_output",
placeholder {
		if bytes.Contains(body, []byte(itemType)) {
			return true
	placeholder
placeholder
	return false
placeholder

func truncateOpenAIResponsesMessageText(item map[string]any) bool {
	itemType := strings.TrimSpace(firstNonEmptyString(item["type"]))
	role := strings.TrimSpace(firstNonEmptyString(item["role"]))
	if itemType != "message" && role == "" {
		return false
placeholder
	parts, ok := item["content"].([]any)
	if !ok {
		return false
placeholder
	modified := false
	for _, rawPart := range parts {
		part, ok := rawPart.(map[string]any)
		if !ok {
			continue
	placeholder
		text, ok := part["text"].(string)
		if !ok {
			continue
	placeholder
		if truncated, changed := truncateOpenAIResponsesInputString(text); changed {
			part["text"] = truncated
			modified = true
	placeholder
placeholder
	return modified
placeholder

func truncateOpenAIResponsesInputString(value string) (string, bool) {
	if len(value) <= openAIResponsesInputTextMaxChars {
		return value, false
placeholder
	chars := 0
	for index := range value {
		if chars == openAIResponsesInputTextMaxChars {
			return value[:index], true
	placeholder
		chars++
placeholder
	return value, false
placeholder
