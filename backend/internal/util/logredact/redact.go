package logredact

import (
	"encoding/json"
	"regexp"
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

var defaultSensitiveKeyList = []string{
	"authorization_code",
	"code",
	"code_verifier",
	"access_token",
	"refresh_token",
	"id_token",
	"client_secret",
	"password",
placeholder

var (
	reGOCSPX = regexp.MustCompile(`GOCSPX-[0-9A-Za-z_-]{24,placeholder`)
	reAIza   = regexp.MustCompile(`AIza[0-9A-Za-z_-]{35placeholder`)
)

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

// RedactText 对非结构化文本做轻量脱敏。
//
// 规则：
// - 如果文本本身是 JSON，则按 RedactJSON 处理。
// - 否则尝试对常见 key=value / key:"value" 片段做脱敏。
//
// 注意：该函数用于日志/错误信息兜底，不保证覆盖所有格式。
func RedactText(input string, extraKeys ...string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
placeholder

	raw := []byte(input)
	if json.Valid(raw) {
		return RedactJSON(raw, extraKeys...)
placeholder

	keyAlt := buildKeyAlternation(extraKeys)
	// JSON-like: "access_token":"..."
	reJSONLike := regexp.MustCompile(`(?i)("(?:` + keyAlt + `)"\s*:\s*")([^"]*)(")`)
	// Query-like: access_token=...
	reQueryLike := regexp.MustCompile(`(?i)\b((?:` + keyAlt + `))=([^&\s]+)`)
	// Plain: access_token: ... / access_token = ...
	rePlain := regexp.MustCompile(`(?i)\b((?:` + keyAlt + `))\b(\s*[:=]\s*)([^,\s]+)`)

	out := input
	out = reGOCSPX.ReplaceAllString(out, "GOCSPX-***")
	out = reAIza.ReplaceAllString(out, "AIza***")
	out = reJSONLike.ReplaceAllString(out, `$1***$3`)
	out = reQueryLike.ReplaceAllString(out, `$1=***`)
	out = rePlain.ReplaceAllString(out, `$1$2***`)
	return out
placeholder

func buildKeyAlternation(extraKeys []string) string {
	seen := make(map[string]struct{placeholder, len(defaultSensitiveKeyList)+len(extraKeys))
	keys := make([]string, 0, len(defaultSensitiveKeyList)+len(extraKeys))
	for _, k := range defaultSensitiveKeyList {
		seen[k] = struct{placeholder{placeholder
		keys = append(keys, regexp.QuoteMeta(k))
placeholder
	for _, k := range extraKeys {
		n := normalizeKey(k)
		if n == "" {
			continue
	placeholder
		if _, ok := seen[n]; ok {
			continue
	placeholder
		seen[n] = struct{placeholder{placeholder
		keys = append(keys, regexp.QuoteMeta(n))
placeholder
	return strings.Join(keys, "|")
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
