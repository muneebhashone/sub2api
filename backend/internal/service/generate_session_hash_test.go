//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// ============ 基础优先级测试 ============

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

	parsed := &ParsedRequest{
		MetadataUserID: "user_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2_account__session_123e4567-e89b-12d3-a456-426614174000",
		System:         "You are a helpful assistant.",
		HasSystem:      true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder

	hash := svc.GenerateSessionHash(parsed)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", hash, "metadata session_id should have highest priority")
placeholder

// ============ System + Messages 基础测试 ============

func TestGenerateSessionHash_SystemPlusMessages(t *testing.T) {
	svc := &GatewayService{placeholder

	withSystem := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder
	withoutSystem := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(withSystem)
	h2 := svc.GenerateSessionHash(withoutSystem)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEqual(t, h1, h2, "system prompt should be part of digest, producing different hash")
placeholder

func TestGenerateSessionHash_SystemOnlyProducesHash(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
placeholder
	hash := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, hash, "system prompt alone should produce a hash as part of full digest")
placeholder

func TestGenerateSessionHash_DifferentSystemsSameMessages(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed1 := &ParsedRequest{
		System:    "You are assistant A.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder
	parsed2 := &ParsedRequest{
		System:    "You are assistant B.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "different system prompts with same messages should produce different hashes")
placeholder

func TestGenerateSessionHash_SameSystemSameMessages(t *testing.T) {
	svc := &GatewayService{placeholder

	mk := func() *ParsedRequest {
		return &ParsedRequest{
			System:    "You are a helpful assistant.",
			HasSystem: true,
			Messages: []any{
				map[string]any{"role": "user", "content": "hello"placeholder,
				map[string]any{"role": "assistant", "content": "hi"placeholder,
		placeholder,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same system + same messages should produce identical hash")
placeholder

func TestGenerateSessionHash_DifferentMessagesProduceDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed1 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "help me with Go"placeholder,
	placeholder,
placeholder
	parsed2 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "help me with Python"placeholder,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "same system but different messages should produce different hashes")
placeholder

// ============ SessionContext 核心测试 ============

func TestGenerateSessionHash_DifferentSessionContextProducesDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder

	// 相同消息 + 不同 SessionContext → 不同 hash（解决碰撞问题的核心场景）
	parsed1 := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "192.168.1.1",
			UserAgent: "Mozilla/5.0",
			APIKeyID:  100,
	placeholder,
placeholder
	parsed2 := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "10.0.0.1",
			UserAgent: "curl/7.0",
			APIKeyID:  200,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEqual(t, h1, h2, "same messages but different SessionContext should produce different hashes")
placeholder

func TestGenerateSessionHash_SameSessionContextProducesSameHash(t *testing.T) {
	svc := &GatewayService{placeholder

	mk := func() *ParsedRequest {
		return &ParsedRequest{
			Messages: []any{
				map[string]any{"role": "user", "content": "hello"placeholder,
		placeholder,
			SessionContext: &SessionContext{
				ClientIP:  "192.168.1.1",
				UserAgent: "Mozilla/5.0",
				APIKeyID:  100,
		placeholder,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same messages + same SessionContext should produce identical hash")
placeholder

func TestGenerateSessionHash_MetadataOverridesSessionContext(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed := &ParsedRequest{
		MetadataUserID: "user_a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2_account__session_123e4567-e89b-12d3-a456-426614174000",
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "192.168.1.1",
			UserAgent: "Mozilla/5.0",
			APIKeyID:  100,
	placeholder,
placeholder

	hash := svc.GenerateSessionHash(parsed)
	require.Equal(t, "123e4567-e89b-12d3-a456-426614174000", hash,
		"metadata session_id should take priority over SessionContext")
placeholder

func TestGenerateSessionHash_MetadataJSON_HasHighestPriority(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed := &ParsedRequest{
		MetadataUserID: `{"device_id":"a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2","account_uuid":"","session_id":"c72554f2-1234-5678-abcd-123456789abc"placeholder`,
		System:         "You are a helpful assistant.",
		HasSystem:      true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder

	hash := svc.GenerateSessionHash(parsed)
	require.Equal(t, "c72554f2-1234-5678-abcd-123456789abc", hash, "JSON format metadata session_id should have highest priority")
placeholder

func TestGenerateSessionHash_NilSessionContextBackwardCompatible(t *testing.T) {
	svc := &GatewayService{placeholder

	withCtx := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: nil,
placeholder
	withoutCtx := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(withCtx)
	h2 := svc.GenerateSessionHash(withoutCtx)
	require.Equal(t, h1, h2, "nil SessionContext should produce same hash as no SessionContext")
placeholder

// ============ 多轮连续会话测试 ============

func TestGenerateSessionHash_ContinuousConversation_HashChangesWithMessages(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 模拟连续会话：每增加一轮对话，hash 应该不同（内容累积变化）
	round1 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	round2 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
			map[string]any{"role": "assistant", "content": "Hi there!"placeholder,
			map[string]any{"role": "user", "content": "How are you?"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	round3 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
			map[string]any{"role": "assistant", "content": "Hi there!"placeholder,
			map[string]any{"role": "user", "content": "How are you?"placeholder,
			map[string]any{"role": "assistant", "content": "I'm doing well!"placeholder,
			map[string]any{"role": "user", "content": "Tell me a joke"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

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

	// 同一轮对话重复请求（如重试）应产生相同 hash
	mk := func() *ParsedRequest {
		return &ParsedRequest{
			System:    "You are a helpful assistant.",
			HasSystem: true,
			Messages: []any{
				map[string]any{"role": "user", "content": "hello"placeholder,
				map[string]any{"role": "assistant", "content": "Hi there!"placeholder,
				map[string]any{"role": "user", "content": "How are you?"placeholder,
		placeholder,
			SessionContext: ctx,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same conversation state should produce identical hash on retry")
placeholder

// ============ 消息回退测试 ============

func TestGenerateSessionHash_MessageRollback(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 模拟消息回退：用户删掉最后一轮再重发
	original := &ParsedRequest{
		System:    "System prompt",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "msg1"placeholder,
			map[string]any{"role": "assistant", "content": "reply1"placeholder,
			map[string]any{"role": "user", "content": "msg2"placeholder,
			map[string]any{"role": "assistant", "content": "reply2"placeholder,
			map[string]any{"role": "user", "content": "msg3"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	// 回退到 msg2 后，用新的 msg3 替代
	rollback := &ParsedRequest{
		System:    "System prompt",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "msg1"placeholder,
			map[string]any{"role": "assistant", "content": "reply1"placeholder,
			map[string]any{"role": "user", "content": "msg2"placeholder,
			map[string]any{"role": "assistant", "content": "reply2"placeholder,
			map[string]any{"role": "user", "content": "different msg3"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	hOrig := svc.GenerateSessionHash(original)
	hRollback := svc.GenerateSessionHash(rollback)
	require.NotEqual(t, hOrig, hRollback, "rollback with different last message should produce different hash")
placeholder

func TestGenerateSessionHash_MessageRollbackSameContent(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 回退后重新发送相同内容 → 相同 hash（合理的粘性恢复）
	mk := func() *ParsedRequest {
		return &ParsedRequest{
			System:    "System prompt",
			HasSystem: true,
			Messages: []any{
				map[string]any{"role": "user", "content": "msg1"placeholder,
				map[string]any{"role": "assistant", "content": "reply1"placeholder,
				map[string]any{"role": "user", "content": "msg2"placeholder,
		placeholder,
			SessionContext: ctx,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "rollback and resend same content should produce same hash")
placeholder

// ============ 相同 System、不同用户消息 ============

func TestGenerateSessionHash_SameSystemDifferentUsers(t *testing.T) {
	svc := &GatewayService{placeholder

	// 两个不同用户使用相同 system prompt 但发送不同消息
	user1 := &ParsedRequest{
		System:    "You are a code reviewer.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "Review this Go code"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "1.1.1.1",
			UserAgent: "vscode",
			APIKeyID:  1,
	placeholder,
placeholder
	user2 := &ParsedRequest{
		System:    "You are a code reviewer.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "Review this Python code"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "2.2.2.2",
			UserAgent: "vscode",
			APIKeyID:  2,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(user1)
	h2 := svc.GenerateSessionHash(user2)
	require.NotEqual(t, h1, h2, "different users with different messages should get different hashes")
placeholder

func TestGenerateSessionHash_SameSystemSameMessageDifferentContext(t *testing.T) {
	svc := &GatewayService{placeholder

	// 这是修复的核心场景：两个不同用户发送完全相同的 system + messages（如 "hello"）
	// 有了 SessionContext 后应该产生不同 hash
	user1 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "1.1.1.1",
			UserAgent: "Mozilla/5.0",
			APIKeyID:  10,
	placeholder,
placeholder
	user2 := &ParsedRequest{
		System:    "You are a helpful assistant.",
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "2.2.2.2",
			UserAgent: "Mozilla/5.0",
			APIKeyID:  20,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(user1)
	h2 := svc.GenerateSessionHash(user2)
	require.NotEqual(t, h1, h2, "CRITICAL: same system+messages but different users should get different hashes")
placeholder

// ============ SessionContext 各字段独立影响测试 ============

func TestGenerateSessionHash_SessionContext_IPDifference(t *testing.T) {
	svc := &GatewayService{placeholder

	base := func(ip string) *ParsedRequest {
		return &ParsedRequest{
			Messages: []any{
				map[string]any{"role": "user", "content": "test"placeholder,
		placeholder,
			SessionContext: &SessionContext{
				ClientIP:  ip,
				UserAgent: "same-ua",
				APIKeyID:  1,
		placeholder,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(base("1.1.1.1"))
	h2 := svc.GenerateSessionHash(base("2.2.2.2"))
	require.NotEqual(t, h1, h2, "different IP should produce different hash")
placeholder

func TestGenerateSessionHash_SessionContext_UADifference(t *testing.T) {
	svc := &GatewayService{placeholder

	base := func(ua string) *ParsedRequest {
		return &ParsedRequest{
			Messages: []any{
				map[string]any{"role": "user", "content": "test"placeholder,
		placeholder,
			SessionContext: &SessionContext{
				ClientIP:  "1.1.1.1",
				UserAgent: ua,
				APIKeyID:  1,
		placeholder,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(base("Mozilla/5.0"))
	h2 := svc.GenerateSessionHash(base("curl/7.0"))
	require.NotEqual(t, h1, h2, "different User-Agent should produce different hash")
placeholder

func TestGenerateSessionHash_SessionContext_APIKeyIDDifference(t *testing.T) {
	svc := &GatewayService{placeholder

	base := func(keyID int64) *ParsedRequest {
		return &ParsedRequest{
			Messages: []any{
				map[string]any{"role": "user", "content": "test"placeholder,
		placeholder,
			SessionContext: &SessionContext{
				ClientIP:  "1.1.1.1",
				UserAgent: "same-ua",
				APIKeyID:  keyID,
		placeholder,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(base(1))
	h2 := svc.GenerateSessionHash(base(2))
	require.NotEqual(t, h1, h2, "different APIKeyID should produce different hash")
placeholder

// ============ 多用户并发相同消息场景 ============

func TestGenerateSessionHash_MultipleUsersSameFirstMessage(t *testing.T) {
	svc := &GatewayService{placeholder

	// 模拟 5 个不同用户同时发送 "hello" → 应该产生 5 个不同的 hash
	hashes := make(map[string]bool)
	for i := 0; i < 5; i++ {
		parsed := &ParsedRequest{
			Messages: []any{
				map[string]any{"role": "user", "content": "hello"placeholder,
		placeholder,
			SessionContext: &SessionContext{
				ClientIP:  "192.168.1." + string(rune('1'+i)),
				UserAgent: "client-" + string(rune('A'+i)),
				APIKeyID:  int64(i + 1),
		placeholder,
	placeholder
		h := svc.GenerateSessionHash(parsed)
		require.NotEmpty(t, h)
		require.False(t, hashes[h], "hash collision detected for user %d", i)
		hashes[h] = true
placeholder
	require.Len(t, hashes, 5, "5 different users should produce 5 unique hashes")
placeholder

// ============ 连续会话粘性：多轮对话同一用户 ============

func TestGenerateSessionHash_SameUserGrowingConversation(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "browser", APIKeyID: 42placeholder

	// 模拟同一用户的连续会话，每轮 hash 不同但同用户重试保持一致
	messages := []map[string]any{
		{"role": "user", "content": "msg1"placeholder,
		{"role": "assistant", "content": "reply1"placeholder,
		{"role": "user", "content": "msg2"placeholder,
		{"role": "assistant", "content": "reply2"placeholder,
		{"role": "user", "content": "msg3"placeholder,
		{"role": "assistant", "content": "reply3"placeholder,
		{"role": "user", "content": "msg4"placeholder,
placeholder

	prevHash := ""
	for round := 1; round <= len(messages); round += 2 {
		// 构建前 round 条消息
		msgs := make([]any, round)
		for j := 0; j < round; j++ {
			msgs[j] = messages[j]
	placeholder
		parsed := &ParsedRequest{
			System:         "System",
			HasSystem:      true,
			Messages:       msgs,
			SessionContext: ctx,
	placeholder
		h := svc.GenerateSessionHash(parsed)
		require.NotEmpty(t, h, "round %d hash should not be empty", round)

		if prevHash != "" {
			require.NotEqual(t, prevHash, h, "round %d hash should differ from previous round", round)
	placeholder
		prevHash = h

		// 同一轮重试应该相同
		h2 := svc.GenerateSessionHash(parsed)
		require.Equal(t, h, h2, "retry of round %d should produce same hash", round)
placeholder
placeholder

// ============ 多轮消息内容结构化测试 ============

func TestGenerateSessionHash_MultipleUserMessages(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 5 条用户消息（无 assistant 回复）
	parsed := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "first"placeholder,
			map[string]any{"role": "user", "content": "second"placeholder,
			map[string]any{"role": "user", "content": "third"placeholder,
			map[string]any{"role": "user", "content": "fourth"placeholder,
			map[string]any{"role": "user", "content": "fifth"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h)

	// 修改中间一条消息应该改变 hash
	parsed2 := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "first"placeholder,
			map[string]any{"role": "user", "content": "CHANGED"placeholder,
			map[string]any{"role": "user", "content": "third"placeholder,
			map[string]any{"role": "user", "content": "fourth"placeholder,
			map[string]any{"role": "user", "content": "fifth"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h, h2, "changing any message should change the hash")
placeholder

func TestGenerateSessionHash_MessageOrderMatters(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	parsed1 := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "alpha"placeholder,
			map[string]any{"role": "user", "content": "beta"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder
	parsed2 := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "beta"placeholder,
			map[string]any{"role": "user", "content": "alpha"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "message order should affect the hash")
placeholder

// ============ 复杂内容格式测试 ============

func TestGenerateSessionHash_StructuredContent(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 结构化 content（数组形式）
	parsed := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{"type": "text", "text": "Look at this"placeholder,
					map[string]any{"type": "text", "text": "And this too"placeholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "structured content should produce a hash")
placeholder

func TestGenerateSessionHash_ArraySystemPrompt(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 数组格式的 system prompt
	parsed := &ParsedRequest{
		System: []any{
			map[string]any{"type": "text", "text": "You are a helpful assistant."placeholder,
			map[string]any{"type": "text", "text": "Be concise."placeholder,
	placeholder,
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "array system prompt should produce a hash")
placeholder

// ============ SessionContext 与 cache_control 优先级 ============

func TestGenerateSessionHash_CacheControlOverridesSessionContext(t *testing.T) {
	svc := &GatewayService{placeholder

	// 当有 cache_control: ephemeral 时，使用第 2 级优先级
	// SessionContext 不应影响结果
	parsed1 := &ParsedRequest{
		System: []any{
			map[string]any{
				"type":          "text",
				"text":          "You are a tool-specific assistant.",
				"cache_control": map[string]any{"type": "ephemeral"placeholder,
		placeholder,
	placeholder,
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "1.1.1.1",
			UserAgent: "ua1",
			APIKeyID:  100,
	placeholder,
placeholder
	parsed2 := &ParsedRequest{
		System: []any{
			map[string]any{
				"type":          "text",
				"text":          "You are a tool-specific assistant.",
				"cache_control": map[string]any{"type": "ephemeral"placeholder,
		placeholder,
	placeholder,
		HasSystem: true,
		Messages: []any{
			map[string]any{"role": "user", "content": "hello"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "2.2.2.2",
			UserAgent: "ua2",
			APIKeyID:  200,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.Equal(t, h1, h2, "cache_control ephemeral has higher priority, SessionContext should not affect result")
placeholder

// ============ 边界情况 ============

func TestGenerateSessionHash_EmptyMessages(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed := &ParsedRequest{
		Messages: []any{placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "1.1.1.1",
			UserAgent: "test",
			APIKeyID:  1,
	placeholder,
placeholder

	// 空 messages + 只有 SessionContext 时，combined.Len() > 0 因为有 context 写入
	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "empty messages with SessionContext should still produce a hash from context")
placeholder

func TestGenerateSessionHash_EmptyMessagesNoContext(t *testing.T) {
	svc := &GatewayService{placeholder

	parsed := &ParsedRequest{
		Messages: []any{placeholder,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.Empty(t, h, "empty messages without SessionContext should produce empty hash")
placeholder

func TestGenerateSessionHash_SessionContextWithEmptyFields(t *testing.T) {
	svc := &GatewayService{placeholder

	// SessionContext 字段为空字符串和零值时仍应影响 hash
	withEmptyCtx := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "test"placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "",
			UserAgent: "",
			APIKeyID:  0,
	placeholder,
placeholder
	withoutCtx := &ParsedRequest{
		Messages: []any{
			map[string]any{"role": "user", "content": "test"placeholder,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(withEmptyCtx)
	h2 := svc.GenerateSessionHash(withoutCtx)
	// 有 SessionContext（即使字段为空）仍然会写入分隔符 "::" 等
	require.NotEqual(t, h1, h2, "empty-field SessionContext should still differ from nil SessionContext")
placeholder

// ============ 长对话历史测试 ============

func TestGenerateSessionHash_LongConversation(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "test", APIKeyID: 1placeholder

	// 构建 20 轮对话
	messages := make([]any, 0, 40)
	for i := 0; i < 20; i++ {
		messages = append(messages, map[string]any{
			"role":    "user",
			"content": "user message " + string(rune('A'+i)),
	placeholder)
		messages = append(messages, map[string]any{
			"role":    "assistant",
			"content": "assistant reply " + string(rune('A'+i)),
	placeholder)
placeholder

	parsed := &ParsedRequest{
		System:         "System prompt",
		HasSystem:      true,
		Messages:       messages,
		SessionContext: ctx,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h)

	// 再加一轮应该不同
	moreMessages := make([]any, len(messages)+2)
	copy(moreMessages, messages)
	moreMessages[len(messages)] = map[string]any{"role": "user", "content": "one more"placeholder
	moreMessages[len(messages)+1] = map[string]any{"role": "assistant", "content": "ok"placeholder

	parsed2 := &ParsedRequest{
		System:         "System prompt",
		HasSystem:      true,
		Messages:       moreMessages,
		SessionContext: ctx,
placeholder

	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h, h2, "adding more messages to long conversation should change hash")
placeholder

// ============ Gemini 原生格式 session hash 测试 ============

func TestGenerateSessionHash_GeminiContentsProducesHash(t *testing.T) {
	svc := &GatewayService{placeholder

	// Gemini 格式: contents[].parts[].text
	parsed := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{"text": "Hello from Gemini"placeholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "1.2.3.4",
			UserAgent: "gemini-cli",
			APIKeyID:  1,
	placeholder,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "Gemini contents with parts should produce a non-empty hash")
placeholder

func TestGenerateSessionHash_GeminiDifferentContentsDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder

	parsed1 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{"text": "Hello"placeholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder
	parsed2 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{"text": "Goodbye"placeholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h1 := svc.GenerateSessionHash(parsed1)
	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h1, h2, "different Gemini contents should produce different hashes")
placeholder

func TestGenerateSessionHash_GeminiSameContentsSameHash(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder

	mk := func() *ParsedRequest {
		return &ParsedRequest{
			Messages: []any{
				map[string]any{
					"role": "user",
					"parts": []any{
						map[string]any{"text": "Hello"placeholder,
				placeholder,
			placeholder,
				map[string]any{
					"role": "model",
					"parts": []any{
						map[string]any{"text": "Hi there!"placeholder,
				placeholder,
			placeholder,
		placeholder,
			SessionContext: ctx,
	placeholder
placeholder

	h1 := svc.GenerateSessionHash(mk())
	h2 := svc.GenerateSessionHash(mk())
	require.Equal(t, h1, h2, "same Gemini contents should produce identical hash")
placeholder

func TestGenerateSessionHash_GeminiMultiTurnHashChanges(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder

	round1 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "hello"placeholderplaceholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	round2 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "hello"placeholderplaceholder,
		placeholder,
			map[string]any{
				"role":  "model",
				"parts": []any{map[string]any{"text": "Hi!"placeholderplaceholder,
		placeholder,
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "How are you?"placeholderplaceholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h1 := svc.GenerateSessionHash(round1)
	h2 := svc.GenerateSessionHash(round2)
	require.NotEmpty(t, h1)
	require.NotEmpty(t, h2)
	require.NotEqual(t, h1, h2, "Gemini multi-turn should produce different hashes per round")
placeholder

func TestGenerateSessionHash_GeminiDifferentUsersSameContentDifferentHash(t *testing.T) {
	svc := &GatewayService{placeholder

	// 核心场景：两个不同用户发送相同 Gemini 格式消息应得到不同 hash
	user1 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "hello"placeholderplaceholder,
		placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "1.1.1.1",
			UserAgent: "gemini-cli",
			APIKeyID:  10,
	placeholder,
placeholder
	user2 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "hello"placeholderplaceholder,
		placeholder,
	placeholder,
		SessionContext: &SessionContext{
			ClientIP:  "2.2.2.2",
			UserAgent: "gemini-cli",
			APIKeyID:  20,
	placeholder,
placeholder

	h1 := svc.GenerateSessionHash(user1)
	h2 := svc.GenerateSessionHash(user2)
	require.NotEqual(t, h1, h2, "CRITICAL: different Gemini users with same content must get different hashes")
placeholder

func TestGenerateSessionHash_GeminiSystemInstructionAffectsHash(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder

	// systemInstruction 经 ParseGatewayRequest 解析后存入 parsed.System
	withSys := &ParsedRequest{
		System: []any{
			map[string]any{"text": "You are a coding assistant."placeholder,
	placeholder,
		Messages: []any{
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "hello"placeholderplaceholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder
	withoutSys := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role":  "user",
				"parts": []any{map[string]any{"text": "hello"placeholderplaceholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h1 := svc.GenerateSessionHash(withSys)
	h2 := svc.GenerateSessionHash(withoutSys)
	require.NotEqual(t, h1, h2, "systemInstruction should affect the hash")
placeholder

func TestGenerateSessionHash_GeminiMultiPartMessage(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder

	// 多 parts 的消息
	parsed := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{"text": "Part 1"placeholder,
					map[string]any{"text": "Part 2"placeholder,
					map[string]any{"text": "Part 3"placeholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "multi-part Gemini message should produce a hash")

	// 不同内容的多 parts
	parsed2 := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{"text": "Part 1"placeholder,
					map[string]any{"text": "CHANGED"placeholder,
					map[string]any{"text": "Part 3"placeholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h2 := svc.GenerateSessionHash(parsed2)
	require.NotEqual(t, h, h2, "changing a part should change the hash")
placeholder

func TestGenerateSessionHash_GeminiNonTextPartsIgnored(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "1.2.3.4", UserAgent: "gemini-cli", APIKeyID: 1placeholder

	// 含非 text 类型 parts（如 inline_data），应被跳过但不报错
	parsed := &ParsedRequest{
		Messages: []any{
			map[string]any{
				"role": "user",
				"parts": []any{
					map[string]any{"text": "Describe this image"placeholder,
					map[string]any{"inline_data": map[string]any{"mime_type": "image/png", "data": "base64..."placeholderplaceholder,
			placeholder,
		placeholder,
	placeholder,
		SessionContext: ctx,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "Gemini message with mixed parts should still produce a hash from text parts")
placeholder

func TestGenerateSessionHash_GeminiMultiTurnHashNotSticky(t *testing.T) {
	svc := &GatewayService{placeholder

	ctx := &SessionContext{ClientIP: "10.0.0.1", UserAgent: "gemini-cli", APIKeyID: 42placeholder

	// 模拟同一 Gemini 会话的三轮请求，每轮 contents 累积增长。
	// 验证预期行为：每轮 hash 都不同，即 GenerateSessionHash 不具备跨轮粘性。
	// 这是 by-design 的——Gemini 的跨轮粘性由 Digest Fallback（BuildGeminiDigestChain）负责。
	round1Body := []byte(`{
		"systemInstruction": {"parts": [{"text": "You are a coding assistant."placeholder]placeholder,
		"contents": [
			{"role": "user", "parts": [{"text": "Write a Go function"placeholder]placeholder
		]
placeholder`)
	round2Body := []byte(`{
		"systemInstruction": {"parts": [{"text": "You are a coding assistant."placeholder]placeholder,
		"contents": [
			{"role": "user", "parts": [{"text": "Write a Go function"placeholder]placeholder,
			{"role": "model", "parts": [{"text": "func hello() {placeholder"placeholder]placeholder,
			{"role": "user", "parts": [{"text": "Add error handling"placeholder]placeholder
		]
placeholder`)
	round3Body := []byte(`{
		"systemInstruction": {"parts": [{"text": "You are a coding assistant."placeholder]placeholder,
		"contents": [
			{"role": "user", "parts": [{"text": "Write a Go function"placeholder]placeholder,
			{"role": "model", "parts": [{"text": "func hello() {placeholder"placeholder]placeholder,
			{"role": "user", "parts": [{"text": "Add error handling"placeholder]placeholder,
			{"role": "model", "parts": [{"text": "func hello() error { return nil placeholder"placeholder]placeholder,
			{"role": "user", "parts": [{"text": "Now add tests"placeholder]placeholder
		]
placeholder`)

	hashes := make([]string, 3)
	for i, body := range [][]byte{round1Body, round2Body, round3Bodyplaceholder {
		parsed, err := ParseGatewayRequest(body, "gemini")
	placeholder
		parsed.SessionContext = ctx
		hashes[i] = svc.GenerateSessionHash(parsed)
		require.NotEmpty(t, hashes[i], "round %d hash should not be empty", i+1)
placeholder

	// 每轮 hash 都不同——这是预期行为
	require.NotEqual(t, hashes[0], hashes[1], "round 1 vs 2 hash should differ (contents grow)")
	require.NotEqual(t, hashes[1], hashes[2], "round 2 vs 3 hash should differ (contents grow)")
	require.NotEqual(t, hashes[0], hashes[2], "round 1 vs 3 hash should differ")

	// 同一轮重试应产生相同 hash
	parsed1Again, err := ParseGatewayRequest(round2Body, "gemini")
placeholder
	parsed1Again.SessionContext = ctx
	h2Again := svc.GenerateSessionHash(parsed1Again)
	require.Equal(t, hashes[1], h2Again, "retry of same round should produce same hash")
placeholder

func TestGenerateSessionHash_GeminiEndToEnd(t *testing.T) {
	svc := &GatewayService{placeholder

	// 端到端测试：模拟 ParseGatewayRequest + GenerateSessionHash 完整流程
	body := []byte(`{
		"model": "gemini-2.5-pro",
		"systemInstruction": {
			"parts": [{"text": "You are a coding assistant."placeholder]
	placeholder,
		"contents": [
			{"role": "user", "parts": [{"text": "Write a Go function"placeholder]placeholder,
			{"role": "model", "parts": [{"text": "Here is a function..."placeholder]placeholder,
			{"role": "user", "parts": [{"text": "Now add error handling"placeholder]placeholder
		]
placeholder`)

	parsed, err := ParseGatewayRequest(body, "gemini")
placeholder
	parsed.SessionContext = &SessionContext{
		ClientIP:  "10.0.0.1",
		UserAgent: "gemini-cli/1.0",
		APIKeyID:  42,
placeholder

	h := svc.GenerateSessionHash(parsed)
	require.NotEmpty(t, h, "end-to-end Gemini flow should produce a hash")

	// 同一请求再次解析应产生相同 hash
	parsed2, err := ParseGatewayRequest(body, "gemini")
placeholder
	parsed2.SessionContext = &SessionContext{
		ClientIP:  "10.0.0.1",
		UserAgent: "gemini-cli/1.0",
		APIKeyID:  42,
placeholder

	h2 := svc.GenerateSessionHash(parsed2)
	require.Equal(t, h, h2, "same request should produce same hash")

	// 不同用户发送相同请求应产生不同 hash
	parsed3, err := ParseGatewayRequest(body, "gemini")
placeholder
	parsed3.SessionContext = &SessionContext{
		ClientIP:  "10.0.0.2",
		UserAgent: "gemini-cli/1.0",
		APIKeyID:  99,
placeholder

	h3 := svc.GenerateSessionHash(parsed3)
	require.NotEqual(t, h, h3, "different user with same Gemini request should get different hash")
placeholder
