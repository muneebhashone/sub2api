package apicompat

import "encoding/json"

func chatResponseFormatToResponsesTextFormat(raw json.RawMessage) json.RawMessage {
	raw = normalizedRawJSON(raw)
	if len(raw) == 0 {
		return nil
placeholder

	obj, ok := rawJSONObject(raw)
	if !ok || rawString(obj["type"]) != "json_schema" {
		return raw
placeholder

	schemaRaw := normalizedRawJSON(obj["json_schema"])
	if len(schemaRaw) == 0 {
		return raw
placeholder

	var schema map[string]json.RawMessage
	if err := json.Unmarshal(schemaRaw, &schema); err != nil {
		return raw
placeholder
	schema["type"] = rawJSONString("json_schema")

	out, err := json.Marshal(schema)
	if err != nil {
		return raw
placeholder
	return out
placeholder

func responsesTextFormatToChatResponseFormat(raw json.RawMessage) json.RawMessage {
	raw = normalizedRawJSON(raw)
	if len(raw) == 0 {
		return nil
placeholder

	obj, ok := rawJSONObject(raw)
	if !ok || rawString(obj["type"]) != "json_schema" {
		return raw
placeholder
	if _, alreadyChatShape := obj["json_schema"]; alreadyChatShape {
		return raw
placeholder

	schema := make(map[string]json.RawMessage, len(obj))
	for key, value := range obj {
		if key == "type" {
			continue
	placeholder
		schema[key] = value
placeholder
	if len(schema) == 0 {
		return raw
placeholder

	schemaRaw, err := json.Marshal(schema)
	if err != nil {
		return raw
placeholder
	out, err := json.Marshal(map[string]json.RawMessage{
		"type":        rawJSONString("json_schema"),
		"json_schema": schemaRaw,
placeholder)
	if err != nil {
		return raw
placeholder
	return out
placeholder

func normalizedRawJSON(raw json.RawMessage) json.RawMessage {
	raw = bytesTrimSpace(raw)
	if len(raw) == 0 || string(raw) == "null" {
		return nil
placeholder
	return append(json.RawMessage(nil), raw...)
placeholder

func rawJSONObject(raw json.RawMessage) (map[string]json.RawMessage, bool) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, false
placeholder
	return obj, true
placeholder

func rawJSONString(value string) json.RawMessage {
	data, _ := json.Marshal(value)
	return data
placeholder
