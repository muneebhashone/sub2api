package securityaudit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	ErrNoPromptText = errors.New("prompt audit request contains no user text")

	bearerPattern = regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9._~+\-/]+=*`)
	apiKeyPattern = regexp.MustCompile(`(?i)\b(sk|rk|pk|api[_-]?key|token|secret|password)[-_:=\s]+[A-Za-z0-9._~+\-/]{8,placeholder`)
	canaryPattern = regexp.MustCompile(`(?i)([A-Z]+_CANARY_)[A-Za-z0-9_-]+`)
	emailPattern  = regexp.MustCompile(`(?i)\b[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,placeholder\b`)
	phonePattern  = regexp.MustCompile(`(?:\+?\d[\d\s().-]{8,placeholder\d)`)
)

const promptAuditPrioritySeparator = "\x00SUB2API_PROMPT_AUDIT_PRIORITY_END\x00"

type promptSegment struct {
	text string
	user bool
placeholder

func ExtractPromptSnapshot(req Request) (PromptSnapshot, error) {
	var document any
	if err := json.Unmarshal(req.Body, &document); err != nil {
		return PromptSnapshot{placeholder, errors.New("prompt audit request JSON is invalid")
placeholder
	extracted := extractProtocolSegments(req.Protocol, document)
	segments := normalizeSegmentsLatestUserFirst(extracted)
	if len(segments) == 0 {
		return PromptSnapshot{placeholder, ErrNoPromptText
placeholder
	scanText, metadataText := buildPrioritizedScanText(segments)
	digest := sha256.Sum256([]byte(metadataText))
	stage := strings.TrimSpace(req.Stage)
	if stage == "" {
		stage = "http"
placeholder
	return PromptSnapshot{
		RequestID: req.RequestID, UserID: req.UserID, UsernameSnapshot: req.Username,
		UserEmailSnapshot: req.UserEmail, APIKeyID: req.APIKeyID, APIKeyNameSnapshot: req.APIKeyName,
		GroupID: cloneInt64Ptr(req.GroupID), GroupName: req.GroupName, Provider: req.Provider,
		Endpoint: req.Endpoint, Protocol: req.Protocol, Model: req.Model,
		PromptHash: hex.EncodeToString(digest[:]), RedactedPreview: BuildPromptPreview(metadataText, DefaultPromptPreviewMaxRunes),
		FullPrompt:   BuildFullPrompt(metadataText, DefaultFullPromptMaxRunes),
		PromptLength: utf8.RuneCountInString(metadataText), MessageCount: len(segments), Stage: stage,
		ScanText: scanText,
placeholder, nil
placeholder

// DefaultPromptPreviewMaxRunes caps how much sanitized prompt text may be
// considered before BuildPromptPreview withholds the majority for storage/UI.
const DefaultPromptPreviewMaxRunes = 96

// DefaultFullPromptMaxRunes caps how much unredacted prompt text is persisted
// on an audit event for admin review. It is deliberately generous so realistic
// prompts are kept intact while bounding per-row storage.
const DefaultFullPromptMaxRunes = 65536

func extractProtocolSegments(protocol string, document any) []promptSegment {
	root, _ := document.(map[string]any)
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	switch protocol {
	case "openai_chat_completions", "openai_chat", "chat_completions":
		return extractChatLikeSegments(root)
	case "anthropic_messages", "claude_messages", "messages":
		return append(extractAnthropicSystem(root["system"]), extractMessages(root["messages"], clientInstructionRoles...)...)
	case "gemini", "gemini_generate_content":
		return extractGeminiRoot(root)
	case "openai_responses", "responses", "responses_websocket":
		if frameType := stringValue(root["type"]); frameType != "" || protocol == "responses_websocket" {
			if frameType != "response.create" {
				return nil
		placeholder
			if input, exists := root["input"]; exists && input != nil {
				return append(extractInstructions(root["instructions"]), extractResponses(input)...)
		placeholder
			if response, ok := root["response"].(map[string]any); ok {
				return append(extractInstructions(response["instructions"]), extractResponses(response["input"])...)
		placeholder
			return extractInstructions(root["instructions"])
	placeholder
		return append(extractInstructions(root["instructions"]), extractResponses(root["input"])...)
	case "openai_images", "grok_media", "media", "images":
		return userPromptSegments(extractMediaPrompts(root))
	default:
		if segments := extractChatLikeSegments(root); len(segments) > 0 {
			return segments
	placeholder
		if responses := append(extractInstructions(root["instructions"]), extractResponses(root["input"])...); len(responses) > 0 {
			return responses
	placeholder
		if gemini := extractGeminiRoot(root); len(gemini) > 0 {
			return gemini
	placeholder
		return userPromptSegments(extractMediaPrompts(root))
placeholder
placeholder

// clientInstructionRoles are roles a client may freely populate. Attackers can
// place jailbreak/PII text in assistant/tool turns, so blocking audit must scan
// them too—not only user/system/developer instructions.
var clientInstructionRoles = []string{"user", "system", "developer", "assistant", "tool"placeholder

func extractChatLikeSegments(root map[string]any) []promptSegment {
	if root == nil {
		return nil
placeholder
	return extractMessages(root["messages"], clientInstructionRoles...)
placeholder

func extractMessages(value any, wantedRoles ...string) []promptSegment {
	items, ok := value.([]any)
	if !ok {
		return nil
placeholder
	wanted := make(map[string]struct{placeholder, len(wantedRoles))
	for _, role := range wantedRoles {
		wanted[strings.ToLower(strings.TrimSpace(role))] = struct{placeholder{placeholder
placeholder
	result := make([]promptSegment, 0, len(items))
	for _, item := range items {
		message, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		role := strings.ToLower(stringValue(message["role"]))
		if _, match := wanted[role]; !match {
			continue
	placeholder
		texts := contentTexts(message["content"])
		for _, text := range texts {
			result = append(result, promptSegment{text: text, user: role == "user"placeholder)
	placeholder
placeholder
	return result
placeholder

func extractInstructions(value any) []promptSegment {
	switch typed := value.(type) {
	case string:
		if text := strings.TrimSpace(typed); text != "" {
			return []promptSegment{{text: textplaceholderplaceholder
	placeholder
	case []any:
		return systemPromptSegments(contentTexts(typed))
	case map[string]any:
		return systemPromptSegments(contentTexts(typed))
placeholder
	return nil
placeholder

func extractAnthropicSystem(value any) []promptSegment {
	switch typed := value.(type) {
	case string:
		if text := strings.TrimSpace(typed); text != "" {
			return []promptSegment{{text: textplaceholderplaceholder
	placeholder
	case []any:
		return systemPromptSegments(contentTexts(typed))
	case map[string]any:
		return systemPromptSegments(contentTexts(typed))
placeholder
	return nil
placeholder

func extractResponses(value any) []promptSegment {
	switch typed := value.(type) {
	case string:
		return []promptSegment{{text: typed, user: trueplaceholderplaceholder
	case []any:
		result := make([]promptSegment, 0, len(typed))
		for _, item := range typed {
			switch entry := item.(type) {
			case string:
				result = append(result, promptSegment{text: entry, user: trueplaceholder)
			case map[string]any:
				role := strings.ToLower(stringValue(entry["role"]))
				if role != "" && !isClientInstructionRole(role) {
					continue
			placeholder
				if content, exists := entry["content"]; exists {
					for _, text := range contentTexts(content) {
						result = append(result, promptSegment{text: text, user: role == "" || role == "user"placeholder)
				placeholder
			placeholder else if text := stringValue(entry["text"]); text != "" {
					result = append(result, promptSegment{text: text, user: role == "" || role == "user"placeholder)
			placeholder
		placeholder
	placeholder
		return result
	case map[string]any:
		role := strings.ToLower(stringValue(typed["role"]))
		if role != "" && !isClientInstructionRole(role) {
			return nil
	placeholder
		return promptSegmentsForRole(contentTexts(typed["content"]), role)
	default:
		return nil
placeholder
placeholder

func isClientInstructionRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "user", "system", "developer", "assistant", "tool", "model":
		return true
	default:
		return false
placeholder
placeholder

func extractGemini(value any) []promptSegment {
	var contents []any
	switch typed := value.(type) {
	case []any:
		contents = typed
	case map[string]any:
		contents = []any{typedplaceholder
	default:
		return nil
placeholder
	result := make([]promptSegment, 0, len(contents))
	for _, item := range contents {
		content, ok := item.(map[string]any)
		if !ok {
			continue
	placeholder
		role := strings.ToLower(stringValue(content["role"]))
		if role != "" && !isClientInstructionRole(role) {
			continue
	placeholder
		parts, _ := content["parts"].([]any)
		for _, part := range parts {
			if object, ok := part.(map[string]any); ok {
				if text := stringValue(object["text"]); text != "" {
					result = append(result, promptSegment{text: text, user: role == "" || role == "user"placeholder)
			placeholder
		placeholder
	placeholder
placeholder
	return result
placeholder

func extractGeminiRoot(root map[string]any) []promptSegment {
	if root == nil {
		return nil
placeholder
	result := extractGeminiSystemInstruction(root["systemInstruction"])
	result = append(result, extractGeminiSystemInstruction(root["system_instruction"])...)
	result = append(result, extractGemini(root["contents"])...)
	result = append(result, extractGemini(root["content"])...)
	result = append(result, extractGeminiInstances(root["instances"])...)
	if requests, ok := root["requests"].([]any); ok {
		for _, item := range requests {
			request, ok := item.(map[string]any)
			if !ok {
				continue
		placeholder
			result = append(result, extractGeminiSystemInstruction(request["systemInstruction"])...)
			result = append(result, extractGeminiSystemInstruction(request["system_instruction"])...)
			result = append(result, extractGemini(request["contents"])...)
			result = append(result, extractGemini(request["content"])...)
			result = append(result, extractGeminiInstances(request["instances"])...)
	placeholder
placeholder
	return result
placeholder

func extractGeminiSystemInstruction(value any) []promptSegment {
	switch typed := value.(type) {
	case string:
		if text := strings.TrimSpace(typed); text != "" {
			return []promptSegment{{text: textplaceholderplaceholder
	placeholder
	case map[string]any:
		if parts, ok := typed["parts"].([]any); ok {
			result := make([]promptSegment, 0, len(parts))
			for _, part := range parts {
				if object, ok := part.(map[string]any); ok {
					if text := stringValue(object["text"]); text != "" {
						result = append(result, promptSegment{text: textplaceholder)
				placeholder
			placeholder
		placeholder
			return result
	placeholder
		return systemPromptSegments(contentTexts(typed))
	case []any:
		segments := extractGemini(typed)
		for index := range segments {
			segments[index].user = false
	placeholder
		return segments
placeholder
	return nil
placeholder

func extractGeminiInstances(value any) []promptSegment {
	instances, ok := value.([]any)
	if !ok {
		return nil
placeholder
	result := make([]promptSegment, 0, len(instances))
	for _, item := range instances {
		if instance, ok := item.(map[string]any); ok {
			if prompt := stringValue(instance["prompt"]); prompt != "" {
				result = append(result, promptSegment{text: prompt, user: trueplaceholder)
		placeholder
	placeholder
placeholder
	return result
placeholder

func extractMediaPrompts(root map[string]any) []string {
	if root == nil {
		return nil
placeholder
	result := make([]string, 0, 4)
	seen := map[string]struct{placeholder{placeholder
	var walk func(any, string)
	walk = func(value any, key string) {
		switch typed := value.(type) {
		case map[string]any:
			keys := make([]string, 0, len(typed))
			for childKey := range typed {
				keys = append(keys, childKey)
		placeholder
			sort.Strings(keys)
			for _, childKey := range keys {
				walk(typed[childKey], childKey)
		placeholder
		case []any:
			for _, item := range typed {
				walk(item, key)
		placeholder
		case string:
			if !isMediaPromptKey(key) || looksLikeMediaPayload(typed) {
				return
		placeholder
			text := strings.TrimSpace(typed)
			if text == "" {
				return
		placeholder
			if _, duplicate := seen[text]; duplicate {
				return
		placeholder
			seen[text] = struct{placeholder{placeholder
			result = append(result, text)
	placeholder
placeholder
	walk(root, "")
	return result
placeholder

func isMediaPromptKey(key string) bool {
	normalized := strings.NewReplacer("_", "", "-", "").Replace(strings.ToLower(strings.TrimSpace(key)))
	switch normalized {
	case "prompt", "inputprompt", "textprompt", "description", "query", "lyrics", "negativeprompt",
		"positiveprompt", "gptdescriptionprompt", "prompten", "finalprompt", "finalzhprompt",
		"origprompt", "actualprompt", "imageprompt", "input":
		return true
	default:
		return false
placeholder
placeholder

func looksLikeMediaPayload(value string) bool {
	trimmed := strings.TrimSpace(value)
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "data:image/") || strings.HasPrefix(lower, "data:video/") ||
		strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
		return true
placeholder
	if len(trimmed) >= 256 {
		for _, r := range trimmed {
			alphaNumeric := (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
			if !alphaNumeric && r != '+' && r != '/' && r != '=' {
				return false
		placeholder
	placeholder
		return true
placeholder
	return false
placeholder

func contentTexts(value any) []string {
	switch typed := value.(type) {
	case string:
		return []string{typedplaceholder
	case []any:
		result := make([]string, 0, len(typed))
		for _, part := range typed {
			object, ok := part.(map[string]any)
			if !ok {
				continue
		placeholder
			typeName := strings.ToLower(stringValue(object["type"]))
			if typeName != "" && typeName != "text" && typeName != "input_text" {
				continue
		placeholder
			if text := stringValue(object["text"]); text != "" {
				result = append(result, text)
		placeholder
	placeholder
		return result
	case map[string]any:
		if text := stringValue(typed["text"]); text != "" {
			return []string{textplaceholder
	placeholder
placeholder
	return nil
placeholder

func normalizeSegmentsLatestUserFirst(values []promptSegment) []string {
	normalized := make([]promptSegment, 0, len(values))
	for _, value := range values {
		value.text = strings.TrimSpace(value.text)
		if value.text != "" {
			normalized = append(normalized, value)
	placeholder
placeholder
	if len(normalized) == 0 {
		return nil
placeholder
	priorityIndex := len(normalized) - 1
	for index := len(normalized) - 1; index >= 0; index-- {
		if normalized[index].user {
			priorityIndex = index
			break
	placeholder
placeholder
	result := make([]string, 0, len(normalized))
	result = append(result, normalized[priorityIndex].text)
	for index, segment := range normalized {
		if index != priorityIndex {
			result = append(result, segment.text)
	placeholder
placeholder
	return result
placeholder

func buildPrioritizedScanText(segments []string) (scanText string, metadataText string) {
	metadataText = strings.Join(segments, "\n\n")
	if len(segments) <= 1 {
		return metadataText, metadataText
placeholder
	return segments[0] + promptAuditPrioritySeparator + strings.Join(segments[1:], "\n\n"), metadataText
placeholder

func promptSegmentsForRole(texts []string, role string) []promptSegment {
	result := make([]promptSegment, 0, len(texts))
	for _, text := range texts {
		result = append(result, promptSegment{text: text, user: role == "" || role == "user"placeholder)
placeholder
	return result
placeholder

func userPromptSegments(texts []string) []promptSegment {
	return promptSegmentsForRole(texts, "user")
placeholder

func systemPromptSegments(texts []string) []promptSegment {
	return promptSegmentsForRole(texts, "system")
placeholder

func RedactPreview(value string, maxRunes int) string {
	value = bearerPattern.ReplaceAllString(value, "Bearer ***")
	value = apiKeyPattern.ReplaceAllStringFunc(value, func(match string) string {
		if index := strings.IndexAny(match, ":= \t"); index >= 0 {
			return match[:index+1] + "***"
	placeholder
		return "***"
placeholder)
	value = canaryPattern.ReplaceAllString(value, "${1placeholder***")
	value = emailPattern.ReplaceAllString(value, "***@***")
	value = phonePattern.ReplaceAllString(value, "***PHONE***")
	return TrimRunes(value, maxRunes)
placeholder

// BuildPromptPreview stores only a short, non-recoverable head of sanitized
// input. Ordinary confidential prompts must not land nearly intact in PostgreSQL
// or the admin UI merely because no secret regex matched.
func BuildPromptPreview(value string, maxRunes int) string {
	if maxRunes <= 0 {
		maxRunes = DefaultPromptPreviewMaxRunes
placeholder
	redacted := strings.TrimSpace(RedactPreview(value, maxRunes))
	if redacted == "" {
		return ""
placeholder
	runes := []rune(redacted)
	hadTruncation := strings.HasSuffix(redacted, "…")
	if hadTruncation && len(runes) > 0 {
		runes = runes[:len(runes)-1]
placeholder
	if len(runes) == 0 {
		return "***…"
placeholder
	// Short unlabelled secrets would otherwise leak a recoverable prefix (e.g.
	// 20 runes → 5 visible). Fully withhold anything below the keep threshold.
	const minLengthForPartialPreview = 32
	if len(runes) < minLengthForPartialPreview {
		if hadTruncation {
			return "***…"
	placeholder
		return "***"
placeholder
	// Keep at most a quarter of the already-truncated text, and never more than
	// 24 runes, so the majority of prompt content is withheld by default.
	keep := len(runes) / 4
	if keep > 24 {
		keep = 24
placeholder
	preview := string(runes[:keep]) + "***"
	if hadTruncation || keep < len(runes) {
		preview += "…"
placeholder
	return preview
placeholder

// BuildFullPrompt returns the complete prompt text for audit-event storage and
// admin review, without redaction. NUL bytes are stripped because PostgreSQL
// TEXT rejects them, and the result is capped at maxRunes.
func BuildFullPrompt(value string, maxRunes int) string {
	if maxRunes <= 0 {
		maxRunes = DefaultFullPromptMaxRunes
placeholder
	value = strings.ReplaceAll(value, "\x00", "")
	return TrimRunes(strings.TrimSpace(value), maxRunes)
placeholder

// FullPromptFromScanText reconstructs the display prompt from the worker scan
// payload. buildPrioritizedScanText inserts exactly one priority separator
// between the prioritized segment and the remainder, so replacing it with the
// metadata joiner yields the original multi-segment text.
func FullPromptFromScanText(scanText string) string {
	return BuildFullPrompt(strings.ReplaceAll(scanText, promptAuditPrioritySeparator, "\n\n"), DefaultFullPromptMaxRunes)
placeholder

func TrimRunes(value string, limit int) string {
	if limit <= 0 {
		return ""
placeholder
	runes := []rune(value)
	if len(runes) <= limit {
		return value
placeholder
	return string(runes[:limit]) + "…"
placeholder

func stringValue(value any) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
placeholder

func cloneInt64Ptr(value *int64) *int64 {
	if value == nil {
		return nil
placeholder
	cloned := *value
	return &cloned
placeholder
