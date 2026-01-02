package logredact

import (
	"encoding/json"
	"strings"
)

// maxRedactDepth 限制递归深度以防止栈溢出
const maxRedactDepth = 32

var defaultSensitiveKeys = map[string]struct{placeholder{
	"authorization_code": {placeholder,
	"code":               {placeholder,
	"code_verifier":      {placeholder,
	"access_token":       {placeholder,
	"refresh_token":      {placeholder,
	"id_token":           {placeholder,
	"client_secret":      {placeholder,
	"password":           {placeholder,
placeholder

func RedactMap(input map[string]any, extraKeys ...string) map[string]any {
	if input == nil {
		return map[string]any{placeholder
placeholder
	keys := buildKeySet(extraKeys)
	redacted, ok := redactValueWithDepth(input, keys, 0).(map[string]any)
	if !ok {
		return map[string]any{placeholder
placeholder
	return redacted
placeholder

func RedactJSON(raw []byte, extraKeys ...string) string {
	if len(raw) == 0 {
		return ""
placeholder
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return "<non-json payload redacted>"
placeholder
	keys := buildKeySet(extraKeys)
	redacted := redactValueWithDepth(value, keys, 0)
	encoded, err := json.Marshal(redacted)
	if err != nil {
		return "<redacted>"
placeholder
	return string(encoded)
placeholder

func buildKeySet(extraKeys []string) map[string]struct{placeholder {
	keys := make(map[string]struct{placeholder, len(defaultSensitiveKeys)+len(extraKeys))
	for k := range defaultSensitiveKeys {
		keys[k] = struct{placeholder{placeholder
placeholder
	for _, key := range extraKeys {
		normalized := normalizeKey(key)
		if normalized == "" {
			continue
	placeholder
		keys[normalized] = struct{placeholder{placeholder
placeholder
	return keys
placeholder

func redactValueWithDepth(value any, keys map[string]struct{placeholder, depth int) any {
	if depth > maxRedactDepth {
		return "<depth limit exceeded>"
placeholder

	switch v := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))
		for k, val := range v {
			if isSensitiveKey(k, keys) {
				out[k] = "***"
				continue
		placeholder
			out[k] = redactValueWithDepth(val, keys, depth+1)
	placeholder
		return out
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = redactValueWithDepth(item, keys, depth+1)
	placeholder
		return out
	default:
		return value
placeholder
placeholder

func isSensitiveKey(key string, keys map[string]struct{placeholder) bool {
	_, ok := keys[normalizeKey(key)]
	return ok
placeholder

func normalizeKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
placeholder
