package service

import (
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

func truncateOpenAIResponsesInputText(_ map[string]any) bool {
	// Do not silently rewrite client or tool output. If an upstream enforces a
	// text limit, forwarding the original value preserves its explicit error for
	// the client and the normal Ops error pipeline. This compatibility shim is
	// retained until the two callers can remove the old mutation hook together.
	return false
placeholder

func openAIResponsesInputMayNeedTruncation(_ []byte) bool {
	// See truncateOpenAIResponsesInputText. Returning false also avoids decoding
	// very large bodies solely for a mutation that must not happen.
	return false
placeholder
