package service

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

var (
	benchmarkToolContinuationBoolSink bool
	benchmarkWSParseStringSink        string
	benchmarkWSParseMapSink           map[string]any
	benchmarkUsageSink                OpenAIUsage
)

func BenchmarkToolContinuationValidationLegacy(b *testing.B) {
	reqBody := benchmarkToolContinuationRequestBody()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkToolContinuationBoolSink = legacyValidateFunctionCallOutputContext(reqBody)
placeholder
placeholder

func BenchmarkToolContinuationValidationOptimized(b *testing.B) {
	reqBody := benchmarkToolContinuationRequestBody()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkToolContinuationBoolSink = optimizedValidateFunctionCallOutputContext(reqBody)
placeholder
placeholder

func BenchmarkWSIngressPayloadParseLegacy(b *testing.B) {
	raw := benchmarkWSIngressPayloadBytes()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventType, model, promptCacheKey, previousResponseID, payload, err := legacyParseWSIngressPayload(raw)
		if err == nil {
			benchmarkWSParseStringSink = eventType + model + promptCacheKey + previousResponseID
			benchmarkWSParseMapSink = payload
	placeholder
placeholder
placeholder

func BenchmarkWSIngressPayloadParseOptimized(b *testing.B) {
	raw := benchmarkWSIngressPayloadBytes()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eventType, model, promptCacheKey, previousResponseID, payload, err := optimizedParseWSIngressPayload(raw)
		if err == nil {
			benchmarkWSParseStringSink = eventType + model + promptCacheKey + previousResponseID
			benchmarkWSParseMapSink = payload
	placeholder
placeholder
placeholder

func BenchmarkOpenAIUsageExtractLegacy(b *testing.B) {
	body := benchmarkOpenAIUsageJSONBytes()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usage, ok := legacyExtractOpenAIUsageFromJSONBytes(body)
		if ok {
			benchmarkUsageSink = usage
	placeholder
placeholder
placeholder

func BenchmarkOpenAIUsageExtractOptimized(b *testing.B) {
	body := benchmarkOpenAIUsageJSONBytes()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		usage, ok := extractOpenAIUsageFromJSONBytes(body)
		if ok {
			benchmarkUsageSink = usage
	placeholder
placeholder
placeholder

func benchmarkToolContinuationRequestBody() map[string]any {
	input := make([]any, 0, 64)
	for i := 0; i < 24; i++ {
		input = append(input, map[string]any{
			"type": "text",
			"text": "benchmark text",
	placeholder)
placeholder
	for i := 0; i < 10; i++ {
		callID := "call_" + strconv.Itoa(i)
		input = append(input, map[string]any{
			"type":    "tool_call",
			"call_id": callID,
	placeholder)
		input = append(input, map[string]any{
			"type":    "function_call_output",
			"call_id": callID,
	placeholder)
		input = append(input, map[string]any{
			"type": "item_reference",
			"id":   callID,
	placeholder)
placeholder
	return map[string]any{
		"model": "gpt-5.3-codex",
		"input": input,
placeholder
placeholder

func benchmarkWSIngressPayloadBytes() []byte {
	return []byte(`{"type":"response.create","model":"gpt-5.3-codex","prompt_cache_key":"cache_bench","previous_response_id":"resp_prev_bench","input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"placeholder]placeholder]placeholder`)
placeholder

func benchmarkOpenAIUsageJSONBytes() []byte {
	return []byte(`{"id":"resp_bench","object":"response","model":"gpt-5.3-codex","usage":{"input_tokens":3210,"output_tokens":987,"input_tokens_details":{"cached_tokens":456placeholderplaceholderplaceholder`)
placeholder

func legacyValidateFunctionCallOutputContext(reqBody map[string]any) bool {
	if !legacyHasFunctionCallOutput(reqBody) {
		return true
placeholder
	previousResponseID, _ := reqBody["previous_response_id"].(string)
	if strings.TrimSpace(previousResponseID) != "" {
		return true
placeholder
	if legacyHasToolCallContext(reqBody) {
		return true
placeholder
	if legacyHasFunctionCallOutputMissingCallID(reqBody) {
		return false
placeholder
	callIDs := legacyFunctionCallOutputCallIDs(reqBody)
	return legacyHasItemReferenceForCallIDs(reqBody, callIDs)
placeholder

func optimizedValidateFunctionCallOutputContext(reqBody map[string]any) bool {
	validation := ValidateFunctionCallOutputContext(reqBody)
	if !validation.HasFunctionCallOutput {
		return true
placeholder
	previousResponseID, _ := reqBody["previous_response_id"].(string)
	if strings.TrimSpace(previousResponseID) != "" {
		return true
placeholder
	if validation.HasToolCallContext {
		return true
placeholder
	if validation.HasFunctionCallOutputMissingCallID {
		return false
placeholder
	return validation.HasItemReferenceForAllCallIDs
placeholder

func legacyHasFunctionCallOutput(reqBody map[string]any) bool {
	if reqBody == nil {
		return false
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return false
placeholder
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if itemType == "function_call_output" {
			return true
	placeholder
placeholder
	return false
placeholder

func legacyHasToolCallContext(reqBody map[string]any) bool {
	if reqBody == nil {
		return false
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return false
placeholder
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if itemType != "tool_call" && itemType != "function_call" {
			continue
	placeholder
		if callID, ok := itemMap["call_id"].(string); ok && strings.TrimSpace(callID) != "" {
			return true
	placeholder
placeholder
	return false
placeholder

func legacyFunctionCallOutputCallIDs(reqBody map[string]any) []string {
	if reqBody == nil {
		return nil
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return nil
placeholder
	ids := make(map[string]struct{placeholder)
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if itemType != "function_call_output" {
			continue
	placeholder
		if callID, ok := itemMap["call_id"].(string); ok && strings.TrimSpace(callID) != "" {
			ids[callID] = struct{placeholder{placeholder
	placeholder
placeholder
	if len(ids) == 0 {
		return nil
placeholder
	callIDs := make([]string, 0, len(ids))
	for id := range ids {
		callIDs = append(callIDs, id)
placeholder
	return callIDs
placeholder

func legacyHasFunctionCallOutputMissingCallID(reqBody map[string]any) bool {
	if reqBody == nil {
		return false
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return false
placeholder
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if itemType != "function_call_output" {
			continue
	placeholder
		callID, _ := itemMap["call_id"].(string)
		if strings.TrimSpace(callID) == "" {
			return true
	placeholder
placeholder
	return false
placeholder

func legacyHasItemReferenceForCallIDs(reqBody map[string]any, callIDs []string) bool {
	if reqBody == nil || len(callIDs) == 0 {
		return false
placeholder
	input, ok := reqBody["input"].([]any)
	if !ok {
		return false
placeholder
	referenceIDs := make(map[string]struct{placeholder)
	for _, item := range input {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		itemType, _ := itemMap["type"].(string)
		if itemType != "item_reference" {
			continue
	placeholder
		idValue, _ := itemMap["id"].(string)
		idValue = strings.TrimSpace(idValue)
		if idValue == "" {
			continue
	placeholder
		referenceIDs[idValue] = struct{placeholder{placeholder
placeholder
	if len(referenceIDs) == 0 {
		return false
placeholder
	for _, callID := range callIDs {
		if _, ok := referenceIDs[callID]; !ok {
			return false
	placeholder
placeholder
	return true
placeholder

func legacyParseWSIngressPayload(raw []byte) (eventType, model, promptCacheKey, previousResponseID string, payload map[string]any, err error) {
	values := gjson.GetManyBytes(raw, "type", "model", "prompt_cache_key", "previous_response_id")
	eventType = strings.TrimSpace(values[0].String())
	if eventType == "" {
		eventType = "response.create"
placeholder
	model = strings.TrimSpace(values[1].String())
	promptCacheKey = strings.TrimSpace(values[2].String())
	previousResponseID = strings.TrimSpace(values[3].String())
	payload = make(map[string]any)
	if err = json.Unmarshal(raw, &payload); err != nil {
		return "", "", "", "", nil, err
placeholder
	if _, exists := payload["type"]; !exists {
		payload["type"] = "response.create"
placeholder
	return eventType, model, promptCacheKey, previousResponseID, payload, nil
placeholder

func optimizedParseWSIngressPayload(raw []byte) (eventType, model, promptCacheKey, previousResponseID string, payload map[string]any, err error) {
	payload = make(map[string]any)
	if err = json.Unmarshal(raw, &payload); err != nil {
		return "", "", "", "", nil, err
placeholder
	eventType = openAIWSPayloadString(payload, "type")
	if eventType == "" {
		eventType = "response.create"
		payload["type"] = eventType
placeholder
	model = openAIWSPayloadString(payload, "model")
	promptCacheKey = openAIWSPayloadString(payload, "prompt_cache_key")
	previousResponseID = openAIWSPayloadString(payload, "previous_response_id")
	return eventType, model, promptCacheKey, previousResponseID, payload, nil
placeholder

func legacyExtractOpenAIUsageFromJSONBytes(body []byte) (OpenAIUsage, bool) {
	var response struct {
		Usage struct {
			InputTokens       int `json:"input_tokens"`
			OutputTokens      int `json:"output_tokens"`
			InputTokenDetails struct {
				CachedTokens int `json:"cached_tokens"`
		placeholder `json:"input_tokens_details"`
	placeholder `json:"usage"`
placeholder
	if err := json.Unmarshal(body, &response); err != nil {
		return OpenAIUsage{placeholder, false
placeholder
	return OpenAIUsage{
		InputTokens:          response.Usage.InputTokens,
		OutputTokens:         response.Usage.OutputTokens,
		CacheReadInputTokens: response.Usage.InputTokenDetails.CachedTokens,
placeholder, true
placeholder
