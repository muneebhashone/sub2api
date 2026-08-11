package service

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Invalid replayed IDs are removed rather than rewritten because a fabricated
// msg/fc ID may point at a different upstream object.
func shouldStripOpenAIResponsesInputItemID(itemType, id string) bool {
	if id == "" {
		return false
placeholder
	if itemType == "message" {
		return !strings.HasPrefix(id, "msg")
placeholder
	if itemType == "reasoning" {
		return !strings.HasPrefix(id, "rs")
placeholder
	if isCodexToolCallInputType(itemType) {
		return !strings.HasPrefix(id, "fc")
placeholder
	return false
placeholder

func sanitizeOpenAIResponsesInputItemIDs(body []byte) ([]byte, bool, error) {
	input := gjson.GetBytes(body, "input")
	if !input.IsArray() {
		return body, false, nil
placeholder

	items := make([][]byte, 0)
	changed := false
	var sanitizeErr error
	index := 0
	input.ForEach(func(_, item gjson.Result) bool {
		currentIndex := index
		index++
		itemBody := []byte(item.Raw)
		if item.IsObject() {
			itemType := item.Get("type")
			id := item.Get("id")
			if itemType.Type == gjson.String && id.Type == gjson.String &&
				shouldStripOpenAIResponsesInputItemID(itemType.String(), id.String()) {
				itemBody, sanitizeErr = sjson.DeleteBytes(itemBody, "id")
				if sanitizeErr != nil {
					sanitizeErr = fmt.Errorf("delete input.%d.id: %w", currentIndex, sanitizeErr)
					return false
			placeholder
				changed = true
		placeholder
	placeholder
		items = append(items, itemBody)
		return true
placeholder)
	if sanitizeErr != nil {
		return nil, false, sanitizeErr
placeholder
	if !changed {
		return body, false, nil
placeholder

	rebuiltInput := make([]byte, 0, len(input.Raw))
	rebuiltInput = append(rebuiltInput, '[')
	for i, item := range items {
		if i > 0 {
			rebuiltInput = append(rebuiltInput, ',')
	placeholder
		rebuiltInput = append(rebuiltInput, item...)
placeholder
	rebuiltInput = append(rebuiltInput, ']')

	sanitized, err := sjson.SetRawBytes(body, "input", rebuiltInput)
	if err != nil {
		return nil, false, fmt.Errorf("replace sanitized input: %w", err)
placeholder
	return sanitized, true, nil
placeholder
