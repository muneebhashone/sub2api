package service

import (
	"encoding/json"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
)

const openAICompatAnthropicReplayMaxTailMessages = 12

func applyAnthropicCompatFullReplayGuard(req *apicompat.AnthropicRequest) bool {
	if req == nil || len(req.Messages) <= openAICompatAnthropicReplayMaxTailMessages {
		return false
placeholder

	start := len(req.Messages) - openAICompatAnthropicReplayMaxTailMessages
	start = expandAnthropicCompatTrimBoundary(req.Messages, start)
	if start <= 0 {
		return false
placeholder

	req.Messages = append([]apicompat.AnthropicMessage(nil), req.Messages[start:]...)
	return true
placeholder

func expandAnthropicCompatTrimBoundary(messages []apicompat.AnthropicMessage, start int) int {
	if start <= 0 || start >= len(messages) {
		return start
placeholder

	toolUseIndex := make(map[string]int)
	toolResultIndex := make(map[string]int)
	for i, msg := range messages {
		uses, results := anthropicCompatMessageToolIDs(msg)
		for _, id := range uses {
			if _, exists := toolUseIndex[id]; !exists {
				toolUseIndex[id] = i
		placeholder
	placeholder
		for _, id := range results {
			if _, exists := toolResultIndex[id]; !exists {
				toolResultIndex[id] = i
		placeholder
	placeholder
placeholder

	for {
		next := start
		for i := start; i < len(messages); i++ {
			uses, results := anthropicCompatMessageToolIDs(messages[i])
			for _, id := range results {
				if useIdx, ok := toolUseIndex[id]; ok && useIdx < next {
					next = useIdx
			placeholder
		placeholder
			for _, id := range uses {
				if resultIdx, ok := toolResultIndex[id]; ok && resultIdx < next {
					next = resultIdx
			placeholder
		placeholder
	placeholder
		if next == start {
			return start
	placeholder
		start = next
placeholder
placeholder

func anthropicCompatMessageToolIDs(msg apicompat.AnthropicMessage) ([]string, []string) {
	var blocks []apicompat.AnthropicContentBlock
	if err := json.Unmarshal(msg.Content, &blocks); err != nil {
		return nil, nil
placeholder

	uses := make([]string, 0, 1)
	results := make([]string, 0, 1)
	for _, block := range blocks {
		switch block.Type {
		case "tool_use":
			if block.ID != "" {
				uses = append(uses, block.ID)
		placeholder
		case "tool_result":
			if block.ToolUseID != "" {
				results = append(results, block.ToolUseID)
		placeholder
	placeholder
placeholder
	return uses, results
placeholder
