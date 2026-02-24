package geminicli

import (
	"encoding/hex"
	"strings"
	"sync"
	"testing"
	"time"
)

// ---------------------------------------------------------------------------
// SessionStore 测试
// ---------------------------------------------------------------------------

func TestSessionStore_SetAndGet(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:     "test-state",
		OAuthType: "code_assist",
		CreatedAt: time.Now(),
placeholder
	store.Set("sid-1", session)

	got, ok := store.Get("sid-1")
	if !ok {
		t.Fatal("期望 Get 返回 ok=true，实际返回 false")
placeholder
	if got.State != "test-state" {
		t.Errorf("期望 State=%q，实际=%q", "test-state", got.State)
placeholder
placeholder

func TestSessionStore_GetNotFound(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	_, ok := store.Get("不存在的ID")
	if ok {
		t.Error("期望不存在的 sessionID 返回 ok=false")
placeholder
placeholder

func TestSessionStore_GetExpired(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	// 创建一个已过期的 session（CreatedAt 设置为 SessionTTL+1 分钟之前）
	session := &OAuthSession{
		State:     "expired-state",
		OAuthType: "code_assist",
		CreatedAt: time.Now().Add(-(SessionTTL + 1*time.Minute)),
placeholder
	store.Set("expired-sid", session)

	_, ok := store.Get("expired-sid")
	if ok {
		t.Error("期望过期的 session 返回 ok=false")
placeholder
placeholder

func TestSessionStore_Delete(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	session := &OAuthSession{
		State:     "to-delete",
		OAuthType: "code_assist",
		CreatedAt: time.Now(),
placeholder
	store.Set("del-sid", session)

	// 先确认存在
	if _, ok := store.Get("del-sid"); !ok {
		t.Fatal("删除前 session 应该存在")
placeholder

	store.Delete("del-sid")

	if _, ok := store.Get("del-sid"); ok {
		t.Error("删除后 session 不应该存在")
placeholder
placeholder

func TestSessionStore_Stop_Idempotent(t *testing.T) {
	store := NewSessionStore()

	// 多次调用 Stop 不应 panic
	store.Stop()
	store.Stop()
	store.Stop()
placeholder

func TestSessionStore_ConcurrentAccess(t *testing.T) {
	store := NewSessionStore()
	defer store.Stop()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	// 并发写入
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			sid := "concurrent-" + string(rune('A'+idx%26))
			store.Set(sid, &OAuthSession{
				State:     sid,
				OAuthType: "code_assist",
				CreatedAt: time.Now(),
		placeholder)
	placeholder(i)
placeholder

	// 并发读取
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			sid := "concurrent-" + string(rune('A'+idx%26))
			store.Get(sid) // 可能找到也可能没找到，关键是不 panic
	placeholder(i)
placeholder

	// 并发删除
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			sid := "concurrent-" + string(rune('A'+idx%26))
			store.Delete(sid)
	placeholder(i)
placeholder

	wg.Wait()
placeholder

// ---------------------------------------------------------------------------
// GenerateRandomBytes 测试
// ---------------------------------------------------------------------------

func TestGenerateRandomBytes(t *testing.T) {
	tests := []int{0, 1, 16, 32, 64placeholder
	for _, n := range tests {
		b, err := GenerateRandomBytes(n)
		if err != nil {
			t.Errorf("GenerateRandomBytes(%d) 出错: %v", n, err)
			continue
	placeholder
		if len(b) != n {
			t.Errorf("GenerateRandomBytes(%d) 返回长度=%d，期望=%d", n, len(b), n)
	placeholder
placeholder
placeholder

func TestGenerateRandomBytes_Uniqueness(t *testing.T) {
	// 两次调用应该返回不同的结果（极小概率相同，32字节足够）
	a, _ := GenerateRandomBytes(32)
	b, _ := GenerateRandomBytes(32)
	if string(a) == string(b) {
		t.Error("两次 GenerateRandomBytes(32) 返回了相同结果，随机性可能有问题")
placeholder
placeholder

// ---------------------------------------------------------------------------
// GenerateState 测试
// ---------------------------------------------------------------------------

func TestGenerateState(t *testing.T) {
	state, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState() 出错: %v", err)
placeholder
	if state == "" {
		t.Error("GenerateState() 返回空字符串")
placeholder
	// base64url 编码不应包含 padding '='
	if strings.Contains(state, "=") {
		t.Errorf("GenerateState() 结果包含 '=' padding: %s", state)
placeholder
	// base64url 不应包含 '+' 或 '/'
	if strings.ContainsAny(state, "+/") {
		t.Errorf("GenerateState() 结果包含非 base64url 字符: %s", state)
placeholder
placeholder

// ---------------------------------------------------------------------------
// GenerateSessionID 测试
// ---------------------------------------------------------------------------

func TestGenerateSessionID(t *testing.T) {
	sid, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("GenerateSessionID() 出错: %v", err)
placeholder
	// 16 字节 -> 32 个 hex 字符
	if len(sid) != 32 {
		t.Errorf("GenerateSessionID() 长度=%d，期望=32", len(sid))
placeholder
	// 必须是合法的 hex 字符串
	if _, err := hex.DecodeString(sid); err != nil {
		t.Errorf("GenerateSessionID() 不是合法的 hex 字符串: %s, err=%v", sid, err)
placeholder
placeholder

func TestGenerateSessionID_Uniqueness(t *testing.T) {
	a, _ := GenerateSessionID()
	b, _ := GenerateSessionID()
	if a == b {
		t.Error("两次 GenerateSessionID() 返回了相同结果")
placeholder
placeholder

// ---------------------------------------------------------------------------
// GenerateCodeVerifier 测试
// ---------------------------------------------------------------------------

func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("GenerateCodeVerifier() 出错: %v", err)
placeholder
	if verifier == "" {
		t.Error("GenerateCodeVerifier() 返回空字符串")
placeholder
	// RFC 7636 要求 code_verifier 至少 43 个字符
	if len(verifier) < 43 {
		t.Errorf("GenerateCodeVerifier() 长度=%d，RFC 7636 要求至少 43 字符", len(verifier))
placeholder
	// base64url 编码不应包含 padding 和非 URL 安全字符
	if strings.Contains(verifier, "=") {
		t.Errorf("GenerateCodeVerifier() 包含 '=' padding: %s", verifier)
placeholder
	if strings.ContainsAny(verifier, "+/") {
		t.Errorf("GenerateCodeVerifier() 包含非 base64url 字符: %s", verifier)
placeholder
placeholder

// ---------------------------------------------------------------------------
// GenerateCodeChallenge 测试
// ---------------------------------------------------------------------------

func TestGenerateCodeChallenge(t *testing.T) {
	// 使用已知输入验证输出
	// RFC 7636 附录 B 示例: verifier = "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	// 预期 challenge = "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expected := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	challenge := GenerateCodeChallenge(verifier)
	if challenge != expected {
		t.Errorf("GenerateCodeChallenge(%q) = %q，期望 %q", verifier, challenge, expected)
placeholder
placeholder

func TestGenerateCodeChallenge_NoPadding(t *testing.T) {
	challenge := GenerateCodeChallenge("test-verifier-string")
	if strings.Contains(challenge, "=") {
		t.Errorf("GenerateCodeChallenge() 结果包含 '=' padding: %s", challenge)
placeholder
placeholder

// ---------------------------------------------------------------------------
// base64URLEncode 测试
// ---------------------------------------------------------------------------

func TestBase64URLEncode(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
placeholder{
		{"空字节", []byte{placeholderplaceholder,
		{"单字节", []byte{0xffplaceholderplaceholder,
		{"多字节", []byte{0x01, 0x02, 0x03, 0x04, 0x05placeholderplaceholder,
		{"全零", []byte{0x00, 0x00, 0x00placeholderplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := base64URLEncode(tt.input)
			// 不应包含 '=' padding
			if strings.Contains(result, "=") {
				t.Errorf("base64URLEncode(%v) 包含 '=' padding: %s", tt.input, result)
		placeholder
			// 不应包含标准 base64 的 '+' 或 '/'
			if strings.ContainsAny(result, "+/") {
				t.Errorf("base64URLEncode(%v) 包含非 URL 安全字符: %s", tt.input, result)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// hasRestrictedScope 测试
// ---------------------------------------------------------------------------

func TestHasRestrictedScope(t *testing.T) {
	tests := []struct {
		scope    string
		expected bool
placeholder{
		// 受限 scope
		{"https://www.googleapis.com/auth/generative-language", trueplaceholder,
		{"https://www.googleapis.com/auth/generative-language.retriever", trueplaceholder,
		{"https://www.googleapis.com/auth/generative-language.tuning", trueplaceholder,
		{"https://www.googleapis.com/auth/drive", trueplaceholder,
		{"https://www.googleapis.com/auth/drive.readonly", trueplaceholder,
		{"https://www.googleapis.com/auth/drive.file", trueplaceholder,
		// 非受限 scope
		{"https://www.googleapis.com/auth/cloud-platform", falseplaceholder,
		{"https://www.googleapis.com/auth/userinfo.email", falseplaceholder,
		{"https://www.googleapis.com/auth/userinfo.profile", falseplaceholder,
		// 边界情况
		{"", falseplaceholder,
		{"random-scope", falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.scope, func(t *testing.T) {
			got := hasRestrictedScope(tt.scope)
			if got != tt.expected {
				t.Errorf("hasRestrictedScope(%q) = %v，期望 %v", tt.scope, got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// BuildAuthorizationURL 测试
// ---------------------------------------------------------------------------

func TestBuildAuthorizationURL(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-secret")

	authURL, err := BuildAuthorizationURL(
		OAuthConfig{placeholder,
		"test-state",
		"test-challenge",
		"https://example.com/callback",
		"",
		"code_assist",
	)
	if err != nil {
		t.Fatalf("BuildAuthorizationURL() 出错: %v", err)
placeholder

	// 检查返回的 URL 包含期望的参数
	checks := []string{
		"response_type=code",
		"client_id=" + GeminiCLIOAuthClientID,
		"redirect_uri=",
		"state=test-state",
		"code_challenge=test-challenge",
		"code_challenge_method=S256",
		"access_type=offline",
		"prompt=consent",
		"include_granted_scopes=true",
placeholder
	for _, check := range checks {
		if !strings.Contains(authURL, check) {
			t.Errorf("BuildAuthorizationURL() URL 缺少参数 %q\nURL: %s", check, authURL)
	placeholder
placeholder

	// 不应包含 project_id（因为传的是空字符串）
	if strings.Contains(authURL, "project_id=") {
		t.Errorf("BuildAuthorizationURL() 空 projectID 时不应包含 project_id 参数")
placeholder

	// URL 应该以正确的授权端点开头
	if !strings.HasPrefix(authURL, AuthorizeURL+"?") {
		t.Errorf("BuildAuthorizationURL() URL 应以 %s? 开头，实际: %s", AuthorizeURL, authURL)
placeholder
placeholder

func TestBuildAuthorizationURL_EmptyRedirectURI(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-secret")

	_, err := BuildAuthorizationURL(
		OAuthConfig{placeholder,
		"test-state",
		"test-challenge",
		"", // 空 redirectURI
		"",
		"code_assist",
	)
	if err == nil {
		t.Error("BuildAuthorizationURL() 空 redirectURI 应该报错")
placeholder
	if !strings.Contains(err.Error(), "redirect_uri") {
		t.Errorf("错误消息应包含 'redirect_uri'，实际: %v", err)
placeholder
placeholder

func TestBuildAuthorizationURL_WithProjectID(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-secret")

	authURL, err := BuildAuthorizationURL(
		OAuthConfig{placeholder,
		"test-state",
		"test-challenge",
		"https://example.com/callback",
		"my-project-123",
		"code_assist",
	)
	if err != nil {
		t.Fatalf("BuildAuthorizationURL() 出错: %v", err)
placeholder
	if !strings.Contains(authURL, "project_id=my-project-123") {
		t.Errorf("BuildAuthorizationURL() 带 projectID 时应包含 project_id 参数\nURL: %s", authURL)
placeholder
placeholder

func TestBuildAuthorizationURL_OAuthConfigError(t *testing.T) {
	// 不设置环境变量，也不提供 client 凭据，EffectiveOAuthConfig 应该报错
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "")

	_, err := BuildAuthorizationURL(
		OAuthConfig{placeholder,
		"test-state",
		"test-challenge",
		"https://example.com/callback",
		"",
		"code_assist",
	)
	if err == nil {
		t.Error("当 EffectiveOAuthConfig 失败时，BuildAuthorizationURL 应该返回错误")
placeholder
placeholder

// ---------------------------------------------------------------------------
// EffectiveOAuthConfig 测试 - 原有测试
// ---------------------------------------------------------------------------

func TestEffectiveOAuthConfig_GoogleOne(t *testing.T) {
	// 内置的 Gemini CLI client secret 不嵌入在此仓库中。
	// 测试通过环境变量设置一个假的 secret 来模拟运维配置。
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	tests := []struct {
		name         string
		input        OAuthConfig
		oauthType    string
		wantClientID string
		wantScopes   string
		wantErr      bool
placeholder{
		{
			name:         "Google One 使用内置客户端（空配置）",
			input:        OAuthConfig{placeholder,
			oauthType:    "google_one",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
	placeholder,
		{
			name: "Google One 使用自定义客户端（传入自定义凭据时使用自定义）",
			input: OAuthConfig{
				ClientID:     "custom-client-id",
				ClientSecret: "custom-client-secret",
		placeholder,
			oauthType:    "google_one",
			wantClientID: "custom-client-id",
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
	placeholder,
		{
			name: "Google One 内置客户端 + 自定义 scopes（应过滤受限 scopes）",
			input: OAuthConfig{
				Scopes: "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly",
		placeholder,
			oauthType:    "google_one",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   "https://www.googleapis.com/auth/cloud-platform",
			wantErr:      false,
	placeholder,
		{
			name: "Google One 内置客户端 + 仅受限 scopes（应回退到默认）",
			input: OAuthConfig{
				Scopes: "https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly",
		placeholder,
			oauthType:    "google_one",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
	placeholder,
		{
			name:         "Code Assist 使用内置客户端",
			input:        OAuthConfig{placeholder,
			oauthType:    "code_assist",
			wantClientID: GeminiCLIOAuthClientID,
			wantScopes:   DefaultCodeAssistScopes,
			wantErr:      false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EffectiveOAuthConfig(tt.input, tt.oauthType)
			if (err != nil) != tt.wantErr {
				t.Errorf("EffectiveOAuthConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
		placeholder
			if err != nil {
				return
		placeholder
			if got.ClientID != tt.wantClientID {
				t.Errorf("EffectiveOAuthConfig() ClientID = %v, want %v", got.ClientID, tt.wantClientID)
		placeholder
			if got.Scopes != tt.wantScopes {
				t.Errorf("EffectiveOAuthConfig() Scopes = %v, want %v", got.Scopes, tt.wantScopes)
		placeholder
	placeholder)
placeholder
placeholder

func TestEffectiveOAuthConfig_ScopeFiltering(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 测试 Google One + 内置客户端过滤受限 scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		Scopes: "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly https://www.googleapis.com/auth/userinfo.profile",
placeholder, "google_one")

	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder

	// 应仅包含 cloud-platform、userinfo.email 和 userinfo.profile
	// 不应包含 generative-language 或 drive scopes
	if strings.Contains(cfg.Scopes, "generative-language") {
		t.Errorf("使用内置客户端时 Scopes 不应包含 generative-language，实际: %v", cfg.Scopes)
placeholder
	if strings.Contains(cfg.Scopes, "drive") {
		t.Errorf("使用内置客户端时 Scopes 不应包含 drive，实际: %v", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "cloud-platform") {
		t.Errorf("Scopes 应包含 cloud-platform，实际: %v", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "userinfo.email") {
		t.Errorf("Scopes 应包含 userinfo.email，实际: %v", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "userinfo.profile") {
		t.Errorf("Scopes 应包含 userinfo.profile，实际: %v", cfg.Scopes)
placeholder
placeholder

// ---------------------------------------------------------------------------
// EffectiveOAuthConfig 测试 - 新增分支覆盖
// ---------------------------------------------------------------------------

func TestEffectiveOAuthConfig_OnlyClientID_NoSecret(t *testing.T) {
	// 只提供 clientID 不提供 secret 应报错
	_, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID: "some-client-id",
placeholder, "code_assist")
	if err == nil {
		t.Error("只提供 ClientID 不提供 ClientSecret 应该报错")
placeholder
	if !strings.Contains(err.Error(), "client_id") || !strings.Contains(err.Error(), "client_secret") {
		t.Errorf("错误消息应提及 client_id 和 client_secret，实际: %v", err)
placeholder
placeholder

func TestEffectiveOAuthConfig_OnlyClientSecret_NoID(t *testing.T) {
	// 只提供 secret 不提供 clientID 应报错
	_, err := EffectiveOAuthConfig(OAuthConfig{
		ClientSecret: "some-client-secret",
placeholder, "code_assist")
	if err == nil {
		t.Error("只提供 ClientSecret 不提供 ClientID 应该报错")
placeholder
	if !strings.Contains(err.Error(), "client_id") || !strings.Contains(err.Error(), "client_secret") {
		t.Errorf("错误消息应提及 client_id 和 client_secret，实际: %v", err)
placeholder
placeholder

func TestEffectiveOAuthConfig_AIStudio_DefaultScopes_BuiltinClient(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// ai_studio 类型，使用内置客户端，scopes 为空 -> 应使用 DefaultCodeAssistScopes（因为内置客户端不能请求 generative-language scope）
	cfg, err := EffectiveOAuthConfig(OAuthConfig{placeholder, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	if cfg.Scopes != DefaultCodeAssistScopes {
		t.Errorf("ai_studio + 内置客户端应使用 DefaultCodeAssistScopes，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_AIStudio_DefaultScopes_CustomClient(t *testing.T) {
	// ai_studio 类型，使用自定义客户端，scopes 为空 -> 应使用 DefaultAIStudioScopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
placeholder, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	if cfg.Scopes != DefaultAIStudioScopes {
		t.Errorf("ai_studio + 自定义客户端应使用 DefaultAIStudioScopes，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_AIStudio_ScopeNormalization(t *testing.T) {
	// ai_studio 类型，旧的 generative-language scope 应被归一化为 generative-language.retriever
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/generative-language https://www.googleapis.com/auth/cloud-platform",
placeholder, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	if strings.Contains(cfg.Scopes, "auth/generative-language ") || strings.HasSuffix(cfg.Scopes, "auth/generative-language") {
		// 确保不包含未归一化的旧 scope（仅 generative-language 而非 generative-language.retriever）
		parts := strings.Fields(cfg.Scopes)
		for _, p := range parts {
			if p == "https://www.googleapis.com/auth/generative-language" {
				t.Errorf("ai_studio 应将 generative-language 归一化为 generative-language.retriever，实际 scopes: %q", cfg.Scopes)
		placeholder
	placeholder
placeholder
	if !strings.Contains(cfg.Scopes, "generative-language.retriever") {
		t.Errorf("ai_studio 归一化后应包含 generative-language.retriever，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_CommaSeparatedScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 逗号分隔的 scopes 应被归一化为空格分隔
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/userinfo.email",
placeholder, "code_assist")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	// 应该用空格分隔，而非逗号
	if strings.Contains(cfg.Scopes, ",") {
		t.Errorf("逗号分隔的 scopes 应被归一化为空格分隔，实际: %q", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "cloud-platform") {
		t.Errorf("归一化后应包含 cloud-platform，实际: %q", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "userinfo.email") {
		t.Errorf("归一化后应包含 userinfo.email，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_MixedCommaAndSpaceScopes(t *testing.T) {
	// 混合逗号和空格分隔的 scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/cloud-platform, https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile",
placeholder, "code_assist")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	parts := strings.Fields(cfg.Scopes)
	if len(parts) != 3 {
		t.Errorf("归一化后应有 3 个 scope，实际: %d，scopes: %q", len(parts), cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_WhitespaceTriming(t *testing.T) {
	// 输入中的前后空白应被清理
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "  custom-id  ",
		ClientSecret: "  custom-secret  ",
		Scopes:       "  https://www.googleapis.com/auth/cloud-platform  ",
placeholder, "code_assist")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	if cfg.ClientID != "custom-id" {
		t.Errorf("ClientID 应去除前后空白，实际: %q", cfg.ClientID)
placeholder
	if cfg.ClientSecret != "custom-secret" {
		t.Errorf("ClientSecret 应去除前后空白，实际: %q", cfg.ClientSecret)
placeholder
	if cfg.Scopes != "https://www.googleapis.com/auth/cloud-platform" {
		t.Errorf("Scopes 应去除前后空白，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_NoEnvSecret(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "")

	cfg, err := EffectiveOAuthConfig(OAuthConfig{placeholder, "code_assist")
	if err != nil {
		t.Fatalf("不设置环境变量时应回退到内置 secret，实际报错: %v", err)
placeholder
	if strings.TrimSpace(cfg.ClientSecret) == "" {
		t.Error("ClientSecret 不应为空")
placeholder
	if cfg.ClientID != GeminiCLIOAuthClientID {
		t.Errorf("ClientID 应回退为内置客户端 ID，实际: %q", cfg.ClientID)
placeholder
placeholder

func TestEffectiveOAuthConfig_AIStudio_BuiltinClient_CustomScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// ai_studio + 内置客户端 + 自定义 scopes -> 应过滤受限 scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		Scopes: "https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/generative-language.retriever",
placeholder, "ai_studio")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	// 内置客户端应过滤 generative-language.retriever
	if strings.Contains(cfg.Scopes, "generative-language") {
		t.Errorf("ai_studio + 内置客户端应过滤受限 scopes，实际: %q", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "cloud-platform") {
		t.Errorf("应保留 cloud-platform scope，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_UnknownOAuthType_DefaultScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 未知的 oauthType 应回退到默认的 code_assist scopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{placeholder, "unknown_type")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	if cfg.Scopes != DefaultCodeAssistScopes {
		t.Errorf("未知 oauthType 应使用 DefaultCodeAssistScopes，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_EmptyOAuthType_DefaultScopes(t *testing.T) {
	t.Setenv(GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	// 空的 oauthType 应走 default 分支，使用 DefaultCodeAssistScopes
	cfg, err := EffectiveOAuthConfig(OAuthConfig{placeholder, "")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	if cfg.Scopes != DefaultCodeAssistScopes {
		t.Errorf("空 oauthType 应使用 DefaultCodeAssistScopes，实际: %q", cfg.Scopes)
placeholder
placeholder

func TestEffectiveOAuthConfig_CustomClient_NoScopeFiltering(t *testing.T) {
	// 自定义客户端 + google_one + 包含受限 scopes -> 不应被过滤（因为不是内置客户端）
	cfg, err := EffectiveOAuthConfig(OAuthConfig{
		ClientID:     "custom-id",
		ClientSecret: "custom-secret",
		Scopes:       "https://www.googleapis.com/auth/generative-language.retriever https://www.googleapis.com/auth/drive.readonly",
placeholder, "google_one")
	if err != nil {
		t.Fatalf("EffectiveOAuthConfig() error = %v", err)
placeholder
	// 自定义客户端不应过滤任何 scope
	if !strings.Contains(cfg.Scopes, "generative-language.retriever") {
		t.Errorf("自定义客户端不应过滤 generative-language.retriever，实际: %q", cfg.Scopes)
placeholder
	if !strings.Contains(cfg.Scopes, "drive.readonly") {
		t.Errorf("自定义客户端不应过滤 drive.readonly，实际: %q", cfg.Scopes)
placeholder
placeholder
