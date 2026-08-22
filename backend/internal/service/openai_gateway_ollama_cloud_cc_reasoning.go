package service

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Ollama Cloud 的 OpenAI 兼容 /v1/chat/completions 把思维放在 reasoning / thinking，
// 而 DeepSeek/OpenAI 客户端只认 reasoning_content。仅在 raw CC 直转路径上做 wire JSON
// 双向补齐，不改 CC↔Responses / Anthropic / Grok 桥。

func isOllamaCloudRawChatCompletionsAccount(account *Account) bool {
	if account == nil || account.Platform != PlatformOpenAI || account.Type != AccountTypeAPIKey {
		return false
placeholder
	mode, _ := account.Extra[openai_compat.ExtraKeyResponsesMode].(string)
	if openai_compat.NormalizeResponsesSupportMode(mode) != openai_compat.ResponsesSupportModeForceChatCompletions {
		return false
placeholder
	if accountHasOllamaCloudUsageExtra(account) {
		return true
placeholder
	if account.Credentials == nil {
		return false
placeholder
	baseURL, _ := account.Credentials["base_url"].(string)
	return isOllamaCloudBaseURL(baseURL)
placeholder

func accountHasOllamaCloudUsageExtra(account *Account) bool {
	if account == nil || account.Extra == nil {
		return false
placeholder
	for _, key := range []string{
		OllamaCloudUsageSessionExtraKey,
		OllamaCloudUsageAutoRefreshExtraKey,
		OllamaCloudUsageSnapshotExtraKey,
placeholder {
		if _, ok := account.Extra[key]; ok {
			return true
	placeholder
placeholder
	return false
placeholder

func applyOllamaCloudRawChatCompletionsRequest(account *Account, body []byte) []byte {
	if !isOllamaCloudRawChatCompletionsAccount(account) || len(body) == 0 {
		return body
placeholder
	return normalizeOllamaCloudChatCompletionsRequest(body)
placeholder

func applyOllamaCloudRawChatCompletionsResponse(account *Account, body []byte) []byte {
	if !isOllamaCloudRawChatCompletionsAccount(account) || len(body) == 0 {
		return body
placeholder
	return normalizeOllamaCloudChatCompletionsResponseJSON(body)
placeholder

func applyOllamaCloudRawChatCompletionsSSELine(account *Account, line string) string {
	if !isOllamaCloudRawChatCompletionsAccount(account) || line == "" {
		return line
placeholder
	return normalizeOllamaCloudChatCompletionsSSELine(line)
placeholder

func normalizeOllamaCloudChatCompletionsRequest(body []byte) []byte {
	if !gjson.ValidBytes(body) {
		return body
placeholder
	messages := gjson.GetBytes(body, "messages")
	if !messages.IsArray() {
		return body
placeholder
	updated := body
	changed := false
	for i, msg := range messages.Array() {
		if msg.Get("role").String() != "assistant" {
			continue
	placeholder
		reasoningContent, ok := jsonNonEmptyString(msg.Get("reasoning_content"))
		if !ok {
			continue
	placeholder
		if _, has := jsonNonEmptyString(msg.Get("reasoning")); has {
			continue
	placeholder
		if _, has := jsonNonEmptyString(msg.Get("thinking")); has {
			continue
	placeholder
		next, err := sjson.SetBytes(updated, "messages."+strconv.Itoa(i)+".reasoning", reasoningContent)
		if err != nil {
			return body
	placeholder
		updated = next
		changed = true
placeholder
	if !changed {
		return body
placeholder
	return updated
placeholder

func normalizeOllamaCloudChatCompletionsResponseJSON(body []byte) []byte {
	if !gjson.ValidBytes(body) {
		return body
placeholder
	choices := gjson.GetBytes(body, "choices")
	if !choices.IsArray() {
		return body
placeholder
	updated := body
	changed := false
	for i, choice := range choices.Array() {
		for _, container := range []string{"message", "delta"placeholder {
			obj := choice.Get(container)
			if !obj.Exists() || !obj.IsObject() {
				continue
		placeholder
			if obj.Get("reasoning_content").Exists() {
				continue
		placeholder
			src, ok := jsonNonEmptyString(obj.Get("reasoning"))
			if !ok {
				src, ok = jsonNonEmptyString(obj.Get("thinking"))
		placeholder
			if !ok {
				continue
		placeholder
			next, err := sjson.SetBytes(updated, "choices."+strconv.Itoa(i)+"."+container+".reasoning_content", src)
			if err != nil {
				return body
		placeholder
			updated = next
			changed = true
	placeholder
placeholder
	if !changed {
		return body
placeholder
	return updated
placeholder

func normalizeOllamaCloudChatCompletionsSSELine(line string) string {
	payload, ok := extractOpenAISSEDataLine(line)
	if !ok {
		return line
placeholder
	trimmed := strings.TrimSpace(payload)
	if trimmed == "" || trimmed == "[DONE]" {
		return line
placeholder
	rewritten := normalizeOllamaCloudChatCompletionsResponseJSON([]byte(payload))
	if string(rewritten) == payload {
		return line
placeholder
	prefixLen := len(line) - len(payload)
	if prefixLen < 0 {
		return line
placeholder
	return line[:prefixLen] + string(rewritten)
placeholder

func jsonNonEmptyString(v gjson.Result) (string, bool) {
	if v.Type != gjson.String || v.Str == "" {
		return "", false
placeholder
	return v.Str, true
placeholder
