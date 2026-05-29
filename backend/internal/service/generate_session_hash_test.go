//go:build unit

package service

import (
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/stretchr/testify/require"
)

func mustParseSessionHashRequest(t *testing.T, body string, ctx *SessionContext) *ParsedRequest {
placeholder
	parsed, err := ParseGatewayRequest(NewRequestBodyRef([]byte(body)), domain.PlatformAnthropic)
placeholder
	parsed.SessionContext = ctx
	return parsed
placeholder

func mustParseGeminiSessionHashRequest(t *testing.T, body string, ctx *SessionContext) *ParsedRequest {
placeholder
	parsed, err := ParseGatewayRequest(NewRequestBodyRef([]byte(body)), domain.PlatformGemini)
placeholder
	parsed.SessionContext = ctx
	return parsed
placeholder

func anthropicSessionBody(system any, messages []any, metadataUserID string) string {
	body := map[string]any{placeholder
	if system != nil {
		body["system"] = system
placeholder
	if messages != nil {
		body["messages"] = messages
placeholder
	if metadataUserID != "" {
		body["metadata"] = map[string]any{"user_id": metadataUserIDplaceholder
placeholder
	data, _ := json.Marshal(body)
	return string(data)
placeholder

func geminiSessionBody(systemParts []any, contents []any) string {
	body := map[string]any{placeholder
	if systemParts != nil {
		body["systemInstruction"] = map[string]any{"parts": systemPartsplaceholder
placeholder
	if contents != nil {
		body["contents"] = contents
placeholder
	data, _ := json.Marshal(body)
	return string(data)
placeholder

func msg(role string, content any) map[string]any {
	return map[string]any{"role": role, "content": contentplaceholder
placeholder

func geminiMsg(role string, texts ...string) map[string]any {
	parts := make([]any, 0, len(texts))
	for _, text := range texts {
		parts = append(parts, map[string]any{"text": textplaceholder)
placeholder
	return map[string]any{"role": role, "parts": partsplaceholder
placeholder

func TestGenerateSessionHash_NilParsedRequest(t *testing.T) {
	svc := &GatewayService{placeholder
	require.Empty(t, svc.GenerateSessionHash(nil))
placeholder

func TestGenerateSessionHash_EmptyRequest(t *testing.T) {
	svc := &GatewayService{placeholder
	require.Empty(t, svc.GenerateSessionHash(&ParsedRequest{placeholder))
placeholder

func TestGenerateSessionHash_MetadataHasHighestPriority(t *testing.T) {
	svc := &GatewayService{placeholder
	metadata := "user_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2_account__session_123e4567-e89b-12d3-a456-426614174000"
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello")placeholder, metadata), nil)

	hash := svc.GenerateSessionHash(parsed)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", hash, "metadata session_id should have highest priority")
placeholder

func TestGenerateSessionHash_SystemPlusMessages(t *testing.T) {
	svc := &GatewayService{placeholder
	withSystem := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello")placeholder, ""), nil)
	withoutSystem := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", "hello")placeholder, ""), nil)

	h1 := svc.GenerateSessionHash(withSystem)
	h2 := svc.GenerateSessionHash(withoutSystem)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEqual(t, h1, h2, "system prompt should be part of digest, producing different hash")
placeholder

func TestGenerateSessionHash_SystemOnlyProducesHash(t *testing.T) {
	svc := &GatewayService{placeholder
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", nil, ""), nil)

	hash := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, hash, "system prompt alone should produce a hash as part of full digest")
placeholder

func TestGenerateSessionHash_DifferentSystemsSameMessages(t *testing.T) {
	svc := &GatewayService{placeholder
	parsed1 := mustParseSessionHashRequest(t, anthropicSessionBody("You are assistant A.", []any{msg("user", "hello")placeholder, ""), nil)
	parsed2 := mustParseSessionHashRequest(t, anthropicSessionBody("You are assistant B.", []any{msg("user", "hello")placeholder, ""), nil)

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "different system prompts with same messages should produce different hashes")
placeholder

func TestGenerateSessionHash_SameSystemSameMessages(t *testing.T) {
	svc := &GatewayService{placeholder
	mk := func() *ParsedRequest {
		return mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello"), msg("assistant", "hi")placeholder, ""), nil)
placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same system + same messages should produce identical hash")
placeholder

func TestGenerateSessionHash_DifferentMessagesProduceDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder
	parsed1 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "help me with Go")placeholder, ""), nil)
	parsed2 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "help me with Python")placeholder, ""), nil)

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "same system but different messages should produce different hashes")
placeholder

func TestGenerateSessionHash_DifferentSessionContextProducesDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "hello")placeholder, "")
	parsed1 := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "192.168.1.1", UserAgent: "Mozilla/5.0", APIKeyID: 100placeholder)
	parsed2 := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "10.0.0.1", UserAgent: "curl/7.0", APIKeyID: 200placeholder)

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEqual(t, h1, h2, "same messages but different SessionContext should produce different hashes")
placeholder

func TestGenerateSessionHash_SameSessionContextProducesSameHash(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "192.168.1.1", UserAgent: "Mozilla/5.0", APIKeyID: 100placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "hello")placeholder, "")
	mk := func() *ParsedRequest { return mustParseSessionHashRequest(t, body, ctx) placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same messages + same SessionContext should produce identical hash")
placeholder

func TestGenerateSessionHash_MetadataOverridesSessionContext(t *testing.T) {
	svc := &GatewayService{placeholder
	metadata := "user_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2_account__session_123e4567-e89b-12d3-a456-426614174000"
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", "hello")placeholder, metadata), &SessionContext{ClientIP: "192.168.1.1", UserAgent: "Mozilla/5.0", APIKeyID: 100placeholder)

	hash := svc.GenerateSessionHash(parsed)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", hash, "metadata session_id should take priority over SessionContext")
placeholder

func TestGenerateSessionHash_MetadataJSON_HasHighestPriority(t *testing.T) {
	svc := &GatewayService{placeholder
	metadata := `{"device_id":"a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2","account_uuid":"","session_id":"c72554f2-1234-5678-abcd-123456789abc"placeholder`
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello")placeholder, metadata), nil)

	hash := svc.GenerateSessionHash(parsed)
	require.Equal(t, "c72554f2-1234-5678-abcd-123456789abc", hash, "JSON format metadata session_id should have highest priority")
placeholder

func TestGenerateSessionHash_NilSessionContextBackwardCompatible(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "hello")placeholder, "")
	withCtx := mustParseSessionHashRequest(t, body, nil)
	withoutCtx := mustParseSessionHashRequest(t, body, nil)

	h1 := svc.GenerateSessionHash(withCtx)
	h2 := svc.GenerateSessionHash(withoutCtx)
	require.Equal(t, h1, h2, "nil SessionContext should produce same hash as no SessionContext")
placeholder

func TestGenerateSessionHash_ContinuousConversation_HashChangesWithMessages(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	round1 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello")placeholder, ""), ctx)
	round2 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello"), msg("assistant", "Hi there!"), msg("user", "How are you?")placeholder, ""), ctx)
	round3 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello"), msg("assistant", "Hi there!"), msg("user", "How are you?"), msg("assistant", "I'm doing well!"), msg("user", "Tell me a joke")placeholder, ""), ctx)

	h1 := svc.GenerateSessionHash(round1)
	h2 := svc.GenerateSessionHash(round2)
	h3 := svc.GenerateSessionHash(round3)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEmpty(t, h3)
	require.NotEqual(t, h1, h2, "different conversation rounds should produce different hashes")
	require.NotEqual(t, h2, h3, "each new round should produce a different hash")
	require.NotEqual(t, h1, h3, "round 1 and round 3 should differ")
placeholder

func TestGenerateSessionHash_ContinuousConversation_SameRoundSameHash(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	body := anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello"), msg("assistant", "Hi there!"), msg("user", "How are you?")placeholder, "")
	mk := func() *ParsedRequest { return mustParseSessionHashRequest(t, body, ctx) placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same conversation state should produce identical hash on retry")
placeholder

func TestGenerateSessionHash_MessageRollback(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	original := mustParseSessionHashRequest(t, anthropicSessionBody("System prompt", []any{msg("user", "msg1"), msg("assistant", "reply1"), msg("user", "msg2"), msg("assistant", "reply2"), msg("user", "msg3")placeholder, ""), ctx)
	rollback := mustParseSessionHashRequest(t, anthropicSessionBody("System prompt", []any{msg("user", "msg1"), msg("assistant", "reply1"), msg("user", "msg2"), msg("assistant", "reply2"), msg("user", "different msg3")placeholder, ""), ctx)

	hOrig := svc.GenerateSessionHash(original)
	hRollback := svc.GenerateSessionHash(rollback)
	require.NotEqual(t, hOrig, hRollback, "rollback with different last message should produce different hash")
placeholder

func TestGenerateSessionHash_MessageRollbackSameContent(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	body := anthropicSessionBody("System prompt", []any{msg("user", "msg1"), msg("assistant", "reply1"), msg("user", "msg2")placeholder, "")
	mk := func() *ParsedRequest { return mustParseSessionHashRequest(t, body, ctx) placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "rollback and resend same content should produce same hash")
placeholder

func TestGenerateSessionHash_SameSystemDifferentUsers(t *testing.T) {
	svc := &GatewayService{placeholder
	user1 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a code reviewer.", []any{msg("user", "Review this Go code")placeholder, ""), &SessionContext{ClientIP: "1.1.1.1", UserAgent: "vscode", APIKeyID: 1placeholder)
	user2 := mustParseSessionHashRequest(t, anthropicSessionBody("You are a code reviewer.", []any{msg("user", "Review this Python code")placeholder, ""), &SessionContext{ClientIP: "2.2.2.2", UserAgent: "vscode", APIKeyID: 2placeholder)

	h1 := svc.GenerateSessionHash(user1)
	h2 := svc.GenerateSessionHash(user2)
	require.NotEqual(t, h1, h2, "different users with different messages should get different hashes")
placeholder

func TestGenerateSessionHash_SameSystemSameMessageDifferentContext(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody("You are a helpful assistant.", []any{msg("user", "hello")placeholder, "")
	user1 := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: "Mozilla/5.0", APIKeyID: 10placeholder)
	user2 := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "2.2.2.2", UserAgent: "Mozilla/5.0", APIKeyID: 20placeholder)

	h1 := svc.GenerateSessionHash(user1)
	h2 := svc.GenerateSessionHash(user2)
	require.NotEqual(t, h1, h2, "CRITICAL: same system+messages but different users should get different hashes")
placeholder

func TestGenerateSessionHash_SessionContext_IPDifference(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "test")placeholder, "")
	base := func(ip string) *ParsedRequest {
		return mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: ip, UserAgent: "same-ua", APIKeyID: 1placeholder)
placeholder

	h1 := svc.GenerateSessionHash(base("1.1.1.1"))
	h2 := svc.GenerateSessionHash(base("2.2.2.2"))
	require.NotEqual(t, h1, h2, "different IP should produce different hash")
placeholder

func TestGenerateSessionHash_SessionContext_UADifference(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "test")placeholder, "")
	base := func(ua string) *ParsedRequest {
		return mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: ua, APIKeyID: 1placeholder)
placeholder

	h1 := svc.GenerateSessionHash(base("Mozilla/5.0"))
	h2 := svc.GenerateSessionHash(base("curl/7.0"))
	require.NotEqual(t, h1, h2, "different User-Agent should produce different hash")
placeholder

func TestGenerateSessionHash_SessionContext_UAVersionNoiseIgnored(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "test")placeholder, "")
	base := func(ua string) *ParsedRequest {
		return mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: ua, APIKeyID: 1placeholder)
placeholder

	h1 := svc.GenerateSessionHash(base("Mozilla/5.0 codex_cli_rs/0.1.0"))
	h2 := svc.GenerateSessionHash(base("Mozilla/5.0 codex_cli_rs/0.1.1"))
	require.Equal(t, h1, h2, "version-only User-Agent changes should not perturb the sticky session hash")
placeholder

func TestGenerateSessionHash_SessionContext_FreeformUAVersionNoiseIgnored(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "test")placeholder, "")
	base := func(ua string) *ParsedRequest {
		return mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: ua, APIKeyID: 1placeholder)
placeholder

	h1 := svc.GenerateSessionHash(base("Codex CLI 0.1.0"))
	h2 := svc.GenerateSessionHash(base("Codex CLI 0.1.1"))
	require.Equal(t, h1, h2, "free-form version-only User-Agent changes should not perturb the sticky session hash")
placeholder

func TestGenerateSessionHash_SessionContext_APIKeyIDDifference(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "test")placeholder, "")
	base := func(keyID int64) *ParsedRequest {
		return mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: "same-ua", APIKeyID: keyIDplaceholder)
placeholder

	h1 := svc.GenerateSessionHash(base(1))
	h2 := svc.GenerateSessionHash(base(2))
	require.NotEqual(t, h1, h2, "different APIKeyID should produce different hash")
placeholder

func TestGenerateSessionHash_MultipleUsersSameFirstMessage(t *testing.T) {
	svc := &GatewayService{placeholder
	hashes := make(map[string]bool)
	body := anthropicSessionBody(nil, []any{msg("user", "hello")placeholder, "")
	for i := 0; i < 5; i++ {
		parsed := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "192.168.1." + string(rune('1'+i)), UserAgent: "client-" + string(rune('A'+i)), APIKeyID: int64(i + 1)placeholder)
		h := svc.GenerateSessionHash(parsed)
		require.NotEmpty(t, h)
		require.False(t, hashes[h], "hash collision detected for user %d", i)
		hashes[h] = true
placeholder
	require.Len(t, hashes, 5, "5 different users should produce 5 unique hashes")
placeholder

func TestGenerateSessionHash_SameUserGrowingConversation(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "browser", APIKeyID: 42placeholder
	messages := []any{
		msg("user", "msg1"), msg("assistant", "reply1"), msg("user", "msg2"), msg("assistant", "reply2"),
		msg("user", "msg3"), msg("assistant", "reply3"), msg("user", "msg4"),
placeholder

	prevHash := ""
	for round := 1; round <= len(messages); round += 2 {
		parsed := mustParseSessionHashRequest(t, anthropicSessionBody("System", messages[:round], ""), ctx)
		h := svc.GenerateSessionHash(parsed)
		require.NotEmpty(t, h, "round %d hash should not be empty", round)
		if prevHash != "" {
			require.NotEqual(t, prevHash, h, "round %d hash should differ from previous round", round)
	placeholder
		prevHash = h
		h2 := svc.GenerateSessionHash(parsed)
		require.Equal(t, h, h2, "retry of round %d should produce same hash", round)
placeholder
placeholder

func TestGenerateSessionHash_MultipleUserMessages(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", "first"), msg("user", "second"), msg("user", "third"), msg("user", "fourth"), msg("user", "fifth")placeholder, ""), ctx)
	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h)

	parsed2 := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", "first"), msg("user", "CHANGED"), msg("user", "third"), msg("user", "fourth"), msg("user", "fifth")placeholder, ""), ctx)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h, h2, "changing any message should change the hash")
placeholder

func TestGenerateSessionHash_MessageOrderMatters(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	parsed1 := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", "alpha"), msg("user", "beta")placeholder, ""), ctx)
	parsed2 := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", "beta"), msg("user", "alpha")placeholder, ""), ctx)

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "message order should affect the hash")
placeholder

func TestGenerateSessionHash_StructuredContent(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	content := []any{map[string]any{"type": "text", "text": "Look at this"placeholder, map[string]any{"type": "text", "text": "And this too"placeholderplaceholder
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{msg("user", content)placeholder, ""), ctx)

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "structured content should produce a hash")
placeholder

func TestGenerateSessionHash_ArraySystemPrompt(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	system := []any{map[string]any{"type": "text", "text": "You are a helpful assistant."placeholder, map[string]any{"type": "text", "text": "Be concise."placeholderplaceholder
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody(system, []any{msg("user", "hello")placeholder, ""), ctx)

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "array system prompt should produce a hash")
placeholder

func TestGenerateSessionHash_CacheControlOverridesSessionContext(t *testing.T) {
	svc := &GatewayService{placeholder
	system := []any{map[string]any{"type": "text", "text": "You are a tool-specific assistant.", "cache_control": map[string]any{"type": "ephemeral"placeholderplaceholderplaceholder
	body := anthropicSessionBody(system, []any{msg("user", "hello")placeholder, "")
	parsed1 := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: "ua1", APIKeyID: 100placeholder)
	parsed2 := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "2.2.2.2", UserAgent: "ua2", APIKeyID: 200placeholder)

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.Equal(t, h1, h2, "cache_control ephemeral has higher priority, SessionContext should not affect result")
placeholder

func TestGenerateSessionHash_EmptyMessages(t *testing.T) {
	svc := &GatewayService{placeholder
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{placeholder, ""), &SessionContext{ClientIP: "1.1.1.1", UserAgent: "test", APIKeyID: 1placeholder)

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "empty messages with SessionContext should still produce a hash from context")
placeholder

func TestGenerateSessionHash_EmptyMessagesNoContext(t *testing.T) {
	svc := &GatewayService{placeholder
	parsed := mustParseSessionHashRequest(t, anthropicSessionBody(nil, []any{placeholder, ""), nil)

	h := svc.GenerateSessionHash(parsed)
	require.Empty(t, h, "empty messages without SessionContext should produce empty hash")
placeholder

func TestGenerateSessionHash_SessionContextWithEmptyFields(t *testing.T) {
	svc := &GatewayService{placeholder
	body := anthropicSessionBody(nil, []any{msg("user", "test")placeholder, "")
	withEmptyCtx := mustParseSessionHashRequest(t, body, &SessionContext{ClientIP: "", UserAgent: "", APIKeyID: 0placeholder)
	withoutCtx := mustParseSessionHashRequest(t, body, nil)

	h1 := svc.GenerateSessionHash(withEmptyCtx)
	h2 := svc.GenerateSessionHash(withoutCtx)
	require.NotEqual(t, h1, h2, "empty-field SessionContext should still differ from nil SessionContext")
placeholder

func TestGenerateSessionHash_LongConversation(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder
	messages := make([]any, 0, 40)
	for i := 0; i < 20; i++ {
		messages = append(messages, msg("user", "user message "+string(rune('A'+i))))
		messages = append(messages, msg("assistant", "assistant reply "+string(rune('A'+i))))
placeholder

	parsed := mustParseSessionHashRequest(t, anthropicSessionBody("System prompt", messages, ""), ctx)
	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h)

	moreMessages := append(append([]any{placeholder, messages...), msg("user", "one more"), msg("assistant", "ok"))
	parsed2 := mustParseSessionHashRequest(t, anthropicSessionBody("System prompt", moreMessages, ""), ctx)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h, h2, "adding more messages to long conversation should change hash")
placeholder

func TestGenerateSessionHash_GeminiContentsProducesHash(t *testing.T) {
	svc := &GatewayService{placeholder
	parsed := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "Hello from Gemini")placeholder), &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder)

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "Gemini contents with parts should produce a non-empty hash")
placeholder

func TestGenerateSessionHash_GeminiDifferentContentsDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder
	parsed1 := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "Hello")placeholder), ctx)
	parsed2 := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "Goodbye")placeholder), ctx)

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "different Gemini contents should produce different hashes")
placeholder

func TestGenerateSessionHash_GeminiSameContentsSameHash(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder
	body := geminiSessionBody(nil, []any{geminiMsg("user", "Hello"), geminiMsg("model", "Hi there!")placeholder)
	mk := func() *ParsedRequest { return mustParseGeminiSessionHashRequest(t, body, ctx) placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same Gemini contents should produce identical hash")
placeholder

func TestGenerateSessionHash_GeminiMultiTurnHashChanges(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder
	round1 := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "hello")placeholder), ctx)
	round2 := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "hello"), geminiMsg("model", "Hi!"), geminiMsg("user", "How are you?")placeholder), ctx)

	h1 := svc.GenerateSessionHash(round1)
	h2 := svc.GenerateSessionHash(round2)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEqual(t, h1, h2, "Gemini multi-turn should produce different hashes per round")
placeholder

func TestGenerateSessionHash_GeminiDifferentUsersSameContentDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder
	body := geminiSessionBody(nil, []any{geminiMsg("user", "hello")placeholder)
	user1 := mustParseGeminiSessionHashRequest(t, body, &SessionContext{ClientIP: "1.1.1.1", UserAgent: "gemini-cli", APIKeyID: 10placeholder)
	user2 := mustParseGeminiSessionHashRequest(t, body, &SessionContext{ClientIP: "2.2.2.2", UserAgent: "gemini-cli", APIKeyID: 20placeholder)

	h1 := svc.GenerateSessionHash(user1)
	h2 := svc.GenerateSessionHash(user2)
	require.NotEqual(t, h1, h2, "CRITICAL: different Gemini users with same content must get different hashes")
placeholder

func TestGenerateSessionHash_GeminiSystemInstructionAffectsHash(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder
	withSys := mustParseGeminiSessionHashRequest(t, geminiSessionBody([]any{map[string]any{"text": "You are a coding assistant."placeholderplaceholder, []any{geminiMsg("user", "hello")placeholder), ctx)
	withoutSys := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "hello")placeholder), ctx)

	h1 := svc.GenerateSessionHash(withSys)
	h2 := svc.GenerateSessionHash(withoutSys)
	require.NotEqual(t, h1, h2, "systemInstruction should affect the hash")
placeholder

func TestGenerateSessionHash_GeminiMultiPartMessage(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder
	parsed := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "Part 1", "Part 2", "Part 3")placeholder), ctx)
	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "multi-part Gemini message should produce a hash")

	parsed2 := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, []any{geminiMsg("user", "Part 1", "CHANGED", "Part 3")placeholder), ctx)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h, h2, "changing a part should change the hash")
placeholder

func TestGenerateSessionHash_GeminiNonTextPartsIgnored(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder
	content := []any{map[string]any{"role": "user", "parts": []any{map[string]any{"text": "Describe this image"placeholder, map[string]any{"inline_data": map[string]any{"mime_type": "image/png", "data": "base64..."placeholderplaceholderplaceholderplaceholderplaceholder
	parsed := mustParseGeminiSessionHashRequest(t, geminiSessionBody(nil, content), ctx)

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "Gemini message with mixed parts should still produce a hash from text parts")
placeholder

func TestGenerateSessionHash_GeminiMultiTurnHashNotSticky(t *testing.T) {
	svc := &GatewayService{placeholder
	ctx := &SessionContext{ClientIP: "10.0.0.1", UserAgent: "gemini-cli", APIKeyID: 42placeholder
	rounds := []string{
		geminiSessionBody([]any{map[string]any{"text": "You are a coding assistant."placeholderplaceholder, []any{geminiMsg("user", "Write a Go function")placeholder),
		geminiSessionBody([]any{map[string]any{"text": "You are a coding assistant."placeholderplaceholder, []any{geminiMsg("user", "Write a Go function"), geminiMsg("model", "func hello() {placeholder"), geminiMsg("user", "Add error handling")placeholder),
		geminiSessionBody([]any{map[string]any{"text": "You are a coding assistant."placeholderplaceholder, []any{geminiMsg("user", "Write a Go function"), geminiMsg("model", "func hello() {placeholder"), geminiMsg("user", "Add error handling"), geminiMsg("model", "func hello() error { return nil placeholder"), geminiMsg("user", "Now add tests")placeholder),
placeholder

	hashes := make([]string, len(rounds))
	for i, body := range rounds {
		parsed := mustParseGeminiSessionHashRequest(t, body, ctx)
		hashes[i] = svc.GenerateSessionHash(parsed)
		require.NotEmpty(t, hashes[i], "round %d hash should not be empty", i+1)
placeholder
	require.NotEqual(t, hashes[0], hashes[1], "round 1 vs 2 hash should differ (contents grow)")
	require.NotEqual(t, hashes[1], hashes[2], "round 2 vs 3 hash should differ (contents grow)")
	require.NotEqual(t, hashes[0], hashes[2], "round 1 vs 3 hash should differ")

	parsedAgain := mustParseGeminiSessionHashRequest(t, rounds[1], ctx)
	h2Again := svc.GenerateSessionHash(parsedAgain)
	require.Equal(t, hashes[1], h2Again, "retry of same round should produce same hash")
placeholder

func TestGenerateSessionHash_GeminiEndToEnd(t *testing.T) {
	svc := &GatewayService{placeholder
	body := geminiSessionBody([]any{map[string]any{"text": "You are a coding assistant."placeholderplaceholder, []any{geminiMsg("user", "Write a Go function"), geminiMsg("model", "Here is a function..."), geminiMsg("user", "Now add error handling")placeholder)
	parsed := mustParseGeminiSessionHashRequest(t, body, &SessionContext{ClientIP: "10.0.0.1", UserAgent: "gemini-cli/1.0", APIKeyID: 42placeholder)

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "end-to-end Gemini flow should produce a hash")

	parsed2 := mustParseGeminiSessionHashRequest(t, body, &SessionContext{ClientIP: "10.0.0.1", UserAgent: "gemini-cli/1.0", APIKeyID: 42placeholder)
	h2 := svc.GenerateSessionHash(parsed2)
	require.Equal(t, h, h2, "same request should produce same hash")

	parsed3 := mustParseGeminiSessionHashRequest(t, body, &SessionContext{ClientIP: "10.0.0.2", UserAgent: "gemini-cli/1.0", APIKeyID: 99placeholder)
	h3 := svc.GenerateSessionHash(parsed3)
	require.NotEqual(t, h, h3, "different user with same Gemini request should get different hash")
placeholder
