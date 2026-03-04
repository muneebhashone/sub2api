//go:build unit

package service

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/geminicli"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

// =====================
// 保留原有测试
// =====================

func TestGeminiOAuthService_GenerateAuthURL_RedirectURIStrategy(t *testing.T) {
	// NOTE: This test sets process env; it must not run in parallel.
	// The built-in Gemini CLI client secret is not embedded in this repository.
	// Tests set a dummy secret via env to simulate operator-provided configuration.
	t.Setenv(geminicli.GeminiCLIOAuthClientSecretEnv, "test-built-in-secret")

	type testCase struct {
		name          string
		cfg           *config.Config
		oauthType     string
		projectID     string
		wantClientID  string
		wantRedirect  string
		wantScope     string
		wantProjectID string
		wantErrSubstr string
placeholder

	tests := []testCase{
		{
			name: "google_one uses built-in client when not configured and redirects to upstream",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{placeholder,
			placeholder,
		placeholder,
			oauthType:     "google_one",
			wantClientID:  geminicli.GeminiCLIOAuthClientID,
			wantRedirect:  geminicli.GeminiCLIRedirectURI,
			wantScope:     geminicli.DefaultCodeAssistScopes,
			wantProjectID: "",
	placeholder,
		{
			name: "google_one always forces built-in client even when custom client configured",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "custom-client-id",
						ClientSecret: "custom-client-secret",
				placeholder,
			placeholder,
		placeholder,
			oauthType:     "google_one",
			wantClientID:  geminicli.GeminiCLIOAuthClientID,
			wantRedirect:  geminicli.GeminiCLIRedirectURI,
			wantScope:     geminicli.DefaultCodeAssistScopes,
			wantProjectID: "",
	placeholder,
		{
			name: "code_assist always forces built-in client even when custom client configured",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "custom-client-id",
						ClientSecret: "custom-client-secret",
				placeholder,
			placeholder,
		placeholder,
			oauthType:     "code_assist",
			projectID:     "my-gcp-project",
			wantClientID:  geminicli.GeminiCLIOAuthClientID,
			wantRedirect:  geminicli.GeminiCLIRedirectURI,
			wantScope:     geminicli.DefaultCodeAssistScopes,
			wantProjectID: "my-gcp-project",
	placeholder,
		{
			name: "ai_studio requires custom client",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{placeholder,
			placeholder,
		placeholder,
			oauthType:     "ai_studio",
			wantErrSubstr: "AI Studio OAuth requires a custom OAuth Client",
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewGeminiOAuthService(nil, nil, nil, nil, tt.cfg)
			got, err := svc.GenerateAuthURL(context.Background(), nil, "https://example.com/auth/callback", tt.projectID, tt.oauthType, "")
			if tt.wantErrSubstr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrSubstr)
			placeholder
				if !strings.Contains(err.Error(), tt.wantErrSubstr) {
					t.Fatalf("expected error containing %q, got: %v", tt.wantErrSubstr, err)
			placeholder
				return
		placeholder
			if err != nil {
				t.Fatalf("GenerateAuthURL returned error: %v", err)
		placeholder

			parsed, err := url.Parse(got.AuthURL)
			if err != nil {
				t.Fatalf("failed to parse auth_url: %v", err)
		placeholder
			q := parsed.Query()

			if gotState := q.Get("state"); gotState != got.State {
				t.Fatalf("state mismatch: query=%q result=%q", gotState, got.State)
		placeholder
			if gotClientID := q.Get("client_id"); gotClientID != tt.wantClientID {
				t.Fatalf("client_id mismatch: got=%q want=%q", gotClientID, tt.wantClientID)
		placeholder
			if gotRedirect := q.Get("redirect_uri"); gotRedirect != tt.wantRedirect {
				t.Fatalf("redirect_uri mismatch: got=%q want=%q", gotRedirect, tt.wantRedirect)
		placeholder
			if gotScope := q.Get("scope"); gotScope != tt.wantScope {
				t.Fatalf("scope mismatch: got=%q want=%q", gotScope, tt.wantScope)
		placeholder
			if gotProjectID := q.Get("project_id"); gotProjectID != tt.wantProjectID {
				t.Fatalf("project_id mismatch: got=%q want=%q", gotProjectID, tt.wantProjectID)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：validateTierID
// =====================

func TestValidateTierID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tierID  string
		wantErr bool
placeholder{
		{name: "空字符串合法", tierID: "", wantErr: falseplaceholder,
		{name: "正常 tier_id", tierID: "google_one_free", wantErr: falseplaceholder,
		{name: "包含斜杠", tierID: "tier/sub", wantErr: falseplaceholder,
		{name: "包含连字符", tierID: "gcp-standard", wantErr: falseplaceholder,
		{name: "纯数字", tierID: "12345", wantErr: falseplaceholder,
		{name: "超长字符串（65个字符）", tierID: strings.Repeat("a", 65), wantErr: trueplaceholder,
		{name: "刚好64个字符", tierID: strings.Repeat("b", 64), wantErr: falseplaceholder,
		{name: "非法字符_空格", tierID: "tier id", wantErr: trueplaceholder,
		{name: "非法字符_中文", tierID: "tier_中文", wantErr: trueplaceholder,
		{name: "非法字符_特殊符号", tierID: "tier@id", wantErr: trueplaceholder,
		{name: "非法字符_感叹号", tierID: "tier!id", wantErr: trueplaceholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateTierID(tt.tierID)
			if tt.wantErr && err == nil {
				t.Fatalf("期望返回错误，但返回 nil")
		placeholder
			if !tt.wantErr && err != nil {
				t.Fatalf("不期望返回错误，但返回: %v", err)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：canonicalGeminiTierID
// =====================

func TestCanonicalGeminiTierID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
placeholder{
		// 空值
		{name: "空字符串", raw: "", want: ""placeholder,
		{name: "纯空白", raw: "   ", want: ""placeholder,

		// 已规范化的值（直接返回）
		{name: "google_one_free", raw: "google_one_free", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "google_ai_pro", raw: "google_ai_pro", want: GeminiTierGoogleAIProplaceholder,
		{name: "google_ai_ultra", raw: "google_ai_ultra", want: GeminiTierGoogleAIUltraplaceholder,
		{name: "gcp_standard", raw: "gcp_standard", want: GeminiTierGCPStandardplaceholder,
		{name: "gcp_enterprise", raw: "gcp_enterprise", want: GeminiTierGCPEnterpriseplaceholder,
		{name: "aistudio_free", raw: "aistudio_free", want: GeminiTierAIStudioFreeplaceholder,
		{name: "aistudio_paid", raw: "aistudio_paid", want: GeminiTierAIStudioPaidplaceholder,
		{name: "google_one_unknown", raw: "google_one_unknown", want: GeminiTierGoogleOneUnknownplaceholder,

		// 大小写不敏感
		{name: "Google_One_Free 大写", raw: "Google_One_Free", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "GCP_STANDARD 全大写", raw: "GCP_STANDARD", want: GeminiTierGCPStandardplaceholder,

		// legacy 映射: Google One
		{name: "AI_PREMIUM -> google_ai_pro", raw: "AI_PREMIUM", want: GeminiTierGoogleAIProplaceholder,
		{name: "FREE -> google_one_free", raw: "FREE", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "GOOGLE_ONE_BASIC -> google_one_free", raw: "GOOGLE_ONE_BASIC", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "GOOGLE_ONE_STANDARD -> google_one_free", raw: "GOOGLE_ONE_STANDARD", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "GOOGLE_ONE_UNLIMITED -> google_ai_ultra", raw: "GOOGLE_ONE_UNLIMITED", want: GeminiTierGoogleAIUltraplaceholder,
		{name: "GOOGLE_ONE_UNKNOWN -> google_one_unknown", raw: "GOOGLE_ONE_UNKNOWN", want: GeminiTierGoogleOneUnknownplaceholder,

		// legacy 映射: Code Assist
		{name: "STANDARD -> gcp_standard", raw: "STANDARD", want: GeminiTierGCPStandardplaceholder,
		{name: "PRO -> gcp_standard", raw: "PRO", want: GeminiTierGCPStandardplaceholder,
		{name: "LEGACY -> gcp_standard", raw: "LEGACY", want: GeminiTierGCPStandardplaceholder,
		{name: "ENTERPRISE -> gcp_enterprise", raw: "ENTERPRISE", want: GeminiTierGCPEnterpriseplaceholder,
		{name: "ULTRA -> gcp_enterprise", raw: "ULTRA", want: GeminiTierGCPEnterpriseplaceholder,

		// kebab-case
		{name: "standard-tier -> gcp_standard", raw: "standard-tier", want: GeminiTierGCPStandardplaceholder,
		{name: "pro-tier -> gcp_standard", raw: "pro-tier", want: GeminiTierGCPStandardplaceholder,
		{name: "ultra-tier -> gcp_enterprise", raw: "ultra-tier", want: GeminiTierGCPEnterpriseplaceholder,

		// 未知值
		{name: "unknown_value -> 空", raw: "unknown_value", want: ""placeholder,
		{name: "random-text -> 空", raw: "random-text", want: ""placeholder,

		// 带空白
		{name: "带前后空白", raw: "  google_one_free  ", want: GeminiTierGoogleOneFreeplaceholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := canonicalGeminiTierID(tt.raw)
			if got != tt.want {
				t.Fatalf("canonicalGeminiTierID(%q) = %q, want %q", tt.raw, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：canonicalGeminiTierIDForOAuthType
// =====================

func TestCanonicalGeminiTierIDForOAuthType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		oauthType string
		tierID    string
		want      string
placeholder{
		// google_one 类型过滤
		{name: "google_one + google_one_free", oauthType: "google_one", tierID: "google_one_free", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "google_one + google_ai_pro", oauthType: "google_one", tierID: "google_ai_pro", want: GeminiTierGoogleAIProplaceholder,
		{name: "google_one + google_ai_ultra", oauthType: "google_one", tierID: "google_ai_ultra", want: GeminiTierGoogleAIUltraplaceholder,
		{name: "google_one + gcp_standard 被过滤", oauthType: "google_one", tierID: "gcp_standard", want: ""placeholder,
		{name: "google_one + aistudio_free 被过滤", oauthType: "google_one", tierID: "aistudio_free", want: ""placeholder,
		{name: "google_one + AI_PREMIUM 遗留映射", oauthType: "google_one", tierID: "AI_PREMIUM", want: GeminiTierGoogleAIProplaceholder,

		// code_assist 类型过滤
		{name: "code_assist + gcp_standard", oauthType: "code_assist", tierID: "gcp_standard", want: GeminiTierGCPStandardplaceholder,
		{name: "code_assist + gcp_enterprise", oauthType: "code_assist", tierID: "gcp_enterprise", want: GeminiTierGCPEnterpriseplaceholder,
		{name: "code_assist + google_one_free 被过滤", oauthType: "code_assist", tierID: "google_one_free", want: ""placeholder,
		{name: "code_assist + aistudio_free 被过滤", oauthType: "code_assist", tierID: "aistudio_free", want: ""placeholder,
		{name: "code_assist + STANDARD 遗留映射", oauthType: "code_assist", tierID: "STANDARD", want: GeminiTierGCPStandardplaceholder,
		{name: "code_assist + standard-tier kebab", oauthType: "code_assist", tierID: "standard-tier", want: GeminiTierGCPStandardplaceholder,

		// ai_studio 类型过滤
		{name: "ai_studio + aistudio_free", oauthType: "ai_studio", tierID: "aistudio_free", want: GeminiTierAIStudioFreeplaceholder,
		{name: "ai_studio + aistudio_paid", oauthType: "ai_studio", tierID: "aistudio_paid", want: GeminiTierAIStudioPaidplaceholder,
		{name: "ai_studio + gcp_standard 被过滤", oauthType: "ai_studio", tierID: "gcp_standard", want: ""placeholder,
		{name: "ai_studio + google_one_free 被过滤", oauthType: "ai_studio", tierID: "google_one_free", want: ""placeholder,

		// 空值
		{name: "空 tierID", oauthType: "google_one", tierID: "", want: ""placeholder,
		{name: "空 oauthType + 有效 tierID", oauthType: "", tierID: "gcp_standard", want: GeminiTierGCPStandardplaceholder,
		{name: "未知 oauthType 接受规范化值", oauthType: "unknown_type", tierID: "gcp_standard", want: GeminiTierGCPStandardplaceholder,

		// oauthType 大小写和空白
		{name: "GOOGLE_ONE 大写", oauthType: "GOOGLE_ONE", tierID: "google_one_free", want: GeminiTierGoogleOneFreeplaceholder,
		{name: "oauthType 带空白", oauthType: "  code_assist  ", tierID: "gcp_standard", want: GeminiTierGCPStandardplaceholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := canonicalGeminiTierIDForOAuthType(tt.oauthType, tt.tierID)
			if got != tt.want {
				t.Fatalf("canonicalGeminiTierIDForOAuthType(%q, %q) = %q, want %q", tt.oauthType, tt.tierID, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：extractTierIDFromAllowedTiers
// =====================

func TestExtractTierIDFromAllowedTiers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		allowedTiers []geminicli.AllowedTier
		want         string
placeholder{
		{
			name:         "nil 列表返回 LEGACY",
			allowedTiers: nil,
			want:         "LEGACY",
	placeholder,
		{
			name:         "空列表返回 LEGACY",
			allowedTiers: []geminicli.AllowedTier{placeholder,
			want:         "LEGACY",
	placeholder,
		{
			name: "有 IsDefault 的 tier",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "STANDARD", IsDefault: falseplaceholder,
				{ID: "PRO", IsDefault: trueplaceholder,
				{ID: "ENTERPRISE", IsDefault: falseplaceholder,
		placeholder,
			want: "PRO",
	placeholder,
		{
			name: "没有 IsDefault 取第一个非空",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "STANDARD", IsDefault: falseplaceholder,
				{ID: "ENTERPRISE", IsDefault: falseplaceholder,
		placeholder,
			want: "STANDARD",
	placeholder,
		{
			name: "IsDefault 的 ID 为空，取第一个非空",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "", IsDefault: trueplaceholder,
				{ID: "PRO", IsDefault: falseplaceholder,
		placeholder,
			want: "PRO",
	placeholder,
		{
			name: "所有 ID 都为空返回 LEGACY",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "", IsDefault: falseplaceholder,
				{ID: "   ", IsDefault: falseplaceholder,
		placeholder,
			want: "LEGACY",
	placeholder,
		{
			name: "ID 带空白会被 trim",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "  STANDARD  ", IsDefault: trueplaceholder,
		placeholder,
			want: "STANDARD",
	placeholder,
		{
			name: "单个 tier 且 IsDefault",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "ENTERPRISE", IsDefault: trueplaceholder,
		placeholder,
			want: "ENTERPRISE",
	placeholder,
		{
			name: "单个 tier 非 IsDefault",
			allowedTiers: []geminicli.AllowedTier{
				{ID: "ENTERPRISE", IsDefault: falseplaceholder,
		placeholder,
			want: "ENTERPRISE",
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := extractTierIDFromAllowedTiers(tt.allowedTiers)
			if got != tt.want {
				t.Fatalf("extractTierIDFromAllowedTiers() = %q, want %q", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：inferGoogleOneTier
// =====================

func TestInferGoogleOneTier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		storageBytes int64
		want         string
placeholder{
		// 边界：<= 0
		{name: "0 bytes -> unknown", storageBytes: 0, want: GeminiTierGoogleOneUnknownplaceholder,
		{name: "负数 -> unknown", storageBytes: -1, want: GeminiTierGoogleOneUnknownplaceholder,

		// > 100TB -> ultra
		{name: "> 100TB -> ultra", storageBytes: int64(StorageTierUnlimited) + 1, want: GeminiTierGoogleAIUltraplaceholder,
		{name: "200TB -> ultra", storageBytes: 200 * int64(TB), want: GeminiTierGoogleAIUltraplaceholder,

		// >= 2TB -> pro (但 <= 100TB)
		{name: "正好 2TB -> pro", storageBytes: int64(StorageTierAIPremium), want: GeminiTierGoogleAIProplaceholder,
		{name: "5TB -> pro", storageBytes: 5 * int64(TB), want: GeminiTierGoogleAIProplaceholder,
		{name: "100TB 正好 -> pro (不是 > 100TB)", storageBytes: int64(StorageTierUnlimited), want: GeminiTierGoogleAIProplaceholder,

		// >= 15GB -> free (但 < 2TB)
		{name: "正好 15GB -> free", storageBytes: int64(StorageTierFree), want: GeminiTierGoogleOneFreeplaceholder,
		{name: "100GB -> free", storageBytes: 100 * int64(GB), want: GeminiTierGoogleOneFreeplaceholder,
		{name: "略低于 2TB -> free", storageBytes: int64(StorageTierAIPremium) - 1, want: GeminiTierGoogleOneFreeplaceholder,

		// < 15GB -> unknown
		{name: "1GB -> unknown", storageBytes: int64(GB), want: GeminiTierGoogleOneUnknownplaceholder,
		{name: "略低于 15GB -> unknown", storageBytes: int64(StorageTierFree) - 1, want: GeminiTierGoogleOneUnknownplaceholder,
		{name: "1 byte -> unknown", storageBytes: 1, want: GeminiTierGoogleOneUnknownplaceholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := inferGoogleOneTier(tt.storageBytes)
			if got != tt.want {
				t.Fatalf("inferGoogleOneTier(%d) = %q, want %q", tt.storageBytes, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：isNonRetryableGeminiOAuthError
// =====================

func TestIsNonRetryableGeminiOAuthError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
placeholder{
		{name: "invalid_grant", err: fmt.Errorf("error: invalid_grant"), want: trueplaceholder,
		{name: "invalid_client", err: fmt.Errorf("oauth error: invalid_client"), want: trueplaceholder,
		{name: "unauthorized_client", err: fmt.Errorf("unauthorized_client: mismatch"), want: trueplaceholder,
		{name: "access_denied", err: fmt.Errorf("access_denied by user"), want: trueplaceholder,
		{name: "普通网络错误", err: fmt.Errorf("connection timeout"), want: falseplaceholder,
		{name: "HTTP 500 错误", err: fmt.Errorf("server error 500"), want: falseplaceholder,
		{name: "空错误信息", err: fmt.Errorf(""), want: falseplaceholder,
		{name: "包含 invalid 但不是完整匹配", err: fmt.Errorf("invalid request"), want: falseplaceholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := isNonRetryableGeminiOAuthError(tt.err)
			if got != tt.want {
				t.Fatalf("isNonRetryableGeminiOAuthError(%v) = %v, want %v", tt.err, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：BuildAccountCredentials
// =====================

func TestGeminiOAuthService_BuildAccountCredentials(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	t.Run("完整字段", func(t *testing.T) {
		t.Parallel()
		tokenInfo := &GeminiTokenInfo{
			AccessToken:  "access-123",
			RefreshToken: "refresh-456",
			ExpiresIn:    3600,
			ExpiresAt:    1700000000,
			TokenType:    "Bearer",
			Scope:        "openid email",
			ProjectID:    "my-project",
			TierID:       "gcp_standard",
			OAuthType:    "code_assist",
			Extra: map[string]any{
				"drive_storage_limit": int64(2199023255552),
		placeholder,
	placeholder

		creds := svc.BuildAccountCredentials(tokenInfo)

		assertCredStr(t, creds, "access_token", "access-123")
		assertCredStr(t, creds, "refresh_token", "refresh-456")
		assertCredStr(t, creds, "token_type", "Bearer")
		assertCredStr(t, creds, "scope", "openid email")
		assertCredStr(t, creds, "project_id", "my-project")
		assertCredStr(t, creds, "tier_id", "gcp_standard")
		assertCredStr(t, creds, "oauth_type", "code_assist")
		assertCredStr(t, creds, "expires_at", "1700000000")

		if _, ok := creds["drive_storage_limit"]; !ok {
			t.Fatal("extra 字段 drive_storage_limit 未包含在 creds 中")
	placeholder
placeholder)

	t.Run("最小字段（仅 access_token 和 expires_at）", func(t *testing.T) {
		t.Parallel()
		tokenInfo := &GeminiTokenInfo{
			AccessToken: "token-only",
			ExpiresAt:   1700000000,
	placeholder

		creds := svc.BuildAccountCredentials(tokenInfo)

		assertCredStr(t, creds, "access_token", "token-only")
		assertCredStr(t, creds, "expires_at", "1700000000")

		// 可选字段不应存在
		for _, key := range []string{"refresh_token", "token_type", "scope", "project_id", "tier_id", "oauth_type"placeholder {
			if _, ok := creds[key]; ok {
				t.Fatalf("不应包含空字段 %q", key)
		placeholder
	placeholder
placeholder)

	t.Run("无效 tier_id 被静默跳过", func(t *testing.T) {
		t.Parallel()
		tokenInfo := &GeminiTokenInfo{
			AccessToken: "token",
			ExpiresAt:   1700000000,
			TierID:      "tier with spaces",
	placeholder

		creds := svc.BuildAccountCredentials(tokenInfo)

		if _, ok := creds["tier_id"]; ok {
			t.Fatal("无效 tier_id 不应被存入 creds")
	placeholder
placeholder)

	t.Run("超长 tier_id 被静默跳过", func(t *testing.T) {
		t.Parallel()
		tokenInfo := &GeminiTokenInfo{
			AccessToken: "token",
			ExpiresAt:   1700000000,
			TierID:      strings.Repeat("x", 65),
	placeholder

		creds := svc.BuildAccountCredentials(tokenInfo)

		if _, ok := creds["tier_id"]; ok {
			t.Fatal("超长 tier_id 不应被存入 creds")
	placeholder
placeholder)

	t.Run("无 extra 字段", func(t *testing.T) {
		t.Parallel()
		tokenInfo := &GeminiTokenInfo{
			AccessToken:  "token",
			ExpiresAt:    1700000000,
			RefreshToken: "rt",
	placeholder

		creds := svc.BuildAccountCredentials(tokenInfo)

		// 仅包含基础字段
		if len(creds) != 3 { // access_token, expires_at, refresh_token
			t.Fatalf("creds 字段数量不匹配: got=%d want=3, keys=%v", len(creds), credKeys(creds))
	placeholder
placeholder)
placeholder

// =====================
// 新增测试：GetOAuthConfig
// =====================

func TestGeminiOAuthService_GetOAuthConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		cfg         *config.Config
		wantEnabled bool
placeholder{
		{
			name: "无自定义 OAuth 客户端",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{placeholder,
			placeholder,
		placeholder,
			wantEnabled: false,
	placeholder,
		{
			name: "仅 ClientID 无 ClientSecret",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID: "custom-id",
				placeholder,
			placeholder,
		placeholder,
			wantEnabled: false,
	placeholder,
		{
			name: "仅 ClientSecret 无 ClientID",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientSecret: "custom-secret",
				placeholder,
			placeholder,
		placeholder,
			wantEnabled: false,
	placeholder,
		{
			name: "使用内置 Gemini CLI ClientID（不算自定义）",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     geminicli.GeminiCLIOAuthClientID,
						ClientSecret: "some-secret",
				placeholder,
			placeholder,
		placeholder,
			wantEnabled: false,
	placeholder,
		{
			name: "自定义 OAuth 客户端（非内置 ID）",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "my-custom-client-id",
						ClientSecret: "my-custom-client-secret",
				placeholder,
			placeholder,
		placeholder,
			wantEnabled: true,
	placeholder,
		{
			name: "带空白的自定义客户端",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "  my-custom-client-id  ",
						ClientSecret: "  my-custom-client-secret  ",
				placeholder,
			placeholder,
		placeholder,
			wantEnabled: true,
	placeholder,
		{
			name: "纯空白字符串不算配置",
			cfg: &config.Config{
				Gemini: config.GeminiConfig{
					OAuth: config.GeminiOAuthConfig{
						ClientID:     "   ",
						ClientSecret: "   ",
				placeholder,
			placeholder,
		placeholder,
			wantEnabled: false,
	placeholder,
placeholder

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := NewGeminiOAuthService(nil, nil, nil, nil, tt.cfg)
			defer svc.Stop()

			result := svc.GetOAuthConfig()
			if result.AIStudioOAuthEnabled != tt.wantEnabled {
				t.Fatalf("AIStudioOAuthEnabled = %v, want %v", result.AIStudioOAuthEnabled, tt.wantEnabled)
		placeholder
			// RequiredRedirectURIs 始终包含 AI Studio redirect URI
			if len(result.RequiredRedirectURIs) != 1 || result.RequiredRedirectURIs[0] != geminicli.AIStudioOAuthRedirectURI {
				t.Fatalf("RequiredRedirectURIs 不匹配: got=%v", result.RequiredRedirectURIs)
		placeholder
	placeholder)
placeholder
placeholder

// =====================
// 新增测试：GeminiOAuthService.Stop
// =====================

func TestGeminiOAuthService_Stop_NoPanic(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)

	// 调用 Stop 不应 panic
	svc.Stop()
	// 多次调用也不应 panic
	svc.Stop()
placeholder

// =====================
// mock: GeminiOAuthClient
// =====================

type mockGeminiOAuthClient struct {
	exchangeCodeFunc func(ctx context.Context, oauthType, code, codeVerifier, redirectURI, proxyURL string) (*geminicli.TokenResponse, error)
	refreshTokenFunc func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error)
placeholder

func (m *mockGeminiOAuthClient) ExchangeCode(ctx context.Context, oauthType, code, codeVerifier, redirectURI, proxyURL string) (*geminicli.TokenResponse, error) {
	if m.exchangeCodeFunc != nil {
		return m.exchangeCodeFunc(ctx, oauthType, code, codeVerifier, redirectURI, proxyURL)
placeholder
	panic("ExchangeCode not implemented")
placeholder

func (m *mockGeminiOAuthClient) RefreshToken(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
	if m.refreshTokenFunc != nil {
		return m.refreshTokenFunc(ctx, oauthType, refreshToken, proxyURL)
placeholder
	panic("RefreshToken not implemented")
placeholder

// =====================
// mock: GeminiCliCodeAssistClient
// =====================

type mockGeminiCodeAssistClient struct {
	loadCodeAssistFunc func(ctx context.Context, accessToken, proxyURL string, req *geminicli.LoadCodeAssistRequest) (*geminicli.LoadCodeAssistResponse, error)
	onboardUserFunc    func(ctx context.Context, accessToken, proxyURL string, req *geminicli.OnboardUserRequest) (*geminicli.OnboardUserResponse, error)
placeholder

func (m *mockGeminiCodeAssistClient) LoadCodeAssist(ctx context.Context, accessToken, proxyURL string, req *geminicli.LoadCodeAssistRequest) (*geminicli.LoadCodeAssistResponse, error) {
	if m.loadCodeAssistFunc != nil {
		return m.loadCodeAssistFunc(ctx, accessToken, proxyURL, req)
placeholder
	panic("LoadCodeAssist not implemented")
placeholder

func (m *mockGeminiCodeAssistClient) OnboardUser(ctx context.Context, accessToken, proxyURL string, req *geminicli.OnboardUserRequest) (*geminicli.OnboardUserResponse, error) {
	if m.onboardUserFunc != nil {
		return m.onboardUserFunc(ctx, accessToken, proxyURL, req)
placeholder
	panic("OnboardUser not implemented")
placeholder

// =====================
// mock: ProxyRepository (最小实现)
// =====================

type mockGeminiProxyRepo struct {
	getByIDFunc func(ctx context.Context, id int64) (*Proxy, error)
placeholder

func (m *mockGeminiProxyRepo) Create(ctx context.Context, proxy *Proxy) error { panic("not impl") placeholder
func (m *mockGeminiProxyRepo) GetByID(ctx context.Context, id int64) (*Proxy, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
placeholder
	return nil, fmt.Errorf("proxy not found")
placeholder
func (m *mockGeminiProxyRepo) ListByIDs(ctx context.Context, ids []int64) ([]Proxy, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) Update(ctx context.Context, proxy *Proxy) error { panic("not impl") placeholder
func (m *mockGeminiProxyRepo) Delete(ctx context.Context, id int64) error     { panic("not impl") placeholder
func (m *mockGeminiProxyRepo) List(ctx context.Context, params pagination.PaginationParams) ([]Proxy, *pagination.PaginationResult, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) ListWithFilters(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]Proxy, *pagination.PaginationResult, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) ListWithFiltersAndAccountCount(ctx context.Context, params pagination.PaginationParams, protocol, status, search string) ([]ProxyWithAccountCount, *pagination.PaginationResult, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) ListActive(ctx context.Context) ([]Proxy, error) { panic("not impl") placeholder
func (m *mockGeminiProxyRepo) ListActiveWithAccountCount(ctx context.Context) ([]ProxyWithAccountCount, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) ExistsByHostPortAuth(ctx context.Context, host string, port int, username, password string) (bool, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) CountAccountsByProxyID(ctx context.Context, proxyID int64) (int64, error) {
	panic("not impl")
placeholder
func (m *mockGeminiProxyRepo) ListAccountSummariesByProxyID(ctx context.Context, proxyID int64) ([]ProxyAccountSummary, error) {
	panic("not impl")
placeholder

// mockDriveClient implements geminicli.DriveClient for tests.
type mockDriveClient struct {
	getStorageQuotaFunc func(ctx context.Context, accessToken, proxyURL string) (*geminicli.DriveStorageInfo, error)
placeholder

func (m *mockDriveClient) GetStorageQuota(ctx context.Context, accessToken, proxyURL string) (*geminicli.DriveStorageInfo, error) {
	if m.getStorageQuotaFunc != nil {
		return m.getStorageQuotaFunc(ctx, accessToken, proxyURL)
placeholder
	return nil, fmt.Errorf("drive API not available in test")
placeholder

// =====================
// 新增测试：GeminiOAuthService.RefreshToken（含重试逻辑）
// =====================

func TestGeminiOAuthService_RefreshToken_Success(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken:  "new-access",
				RefreshToken: "new-refresh",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
				Scope:        "openid",
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(nil, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	info, err := svc.RefreshToken(context.Background(), "code_assist", "old-refresh", "")
	if err != nil {
		t.Fatalf("RefreshToken 返回错误: %v", err)
placeholder
	if info.AccessToken != "new-access" {
		t.Fatalf("AccessToken 不匹配: got=%q", info.AccessToken)
placeholder
	if info.RefreshToken != "new-refresh" {
		t.Fatalf("RefreshToken 不匹配: got=%q", info.RefreshToken)
placeholder
	if info.ExpiresAt == 0 {
		t.Fatal("ExpiresAt 不应为 0")
placeholder
placeholder

func TestGeminiOAuthService_RefreshToken_NonRetryableError(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return nil, fmt.Errorf("invalid_grant: token revoked")
	placeholder,
placeholder

	svc := NewGeminiOAuthService(nil, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	_, err := svc.RefreshToken(context.Background(), "code_assist", "revoked-token", "")
	if err == nil {
		t.Fatal("RefreshToken 应返回错误（不可重试的 invalid_grant）")
placeholder
	if !strings.Contains(err.Error(), "invalid_grant") {
		t.Fatalf("错误应包含 invalid_grant: got=%q", err.Error())
placeholder
placeholder

func TestGeminiOAuthService_RefreshToken_RetryableError(t *testing.T) {
	t.Parallel()

	callCount := 0
	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			callCount++
			if callCount <= 2 {
				return nil, fmt.Errorf("temporary network error")
		placeholder
			return &geminicli.TokenResponse{
				AccessToken: "recovered",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(nil, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	info, err := svc.RefreshToken(context.Background(), "code_assist", "rt", "")
	if err != nil {
		t.Fatalf("RefreshToken 应在重试后成功: %v", err)
placeholder
	if info.AccessToken != "recovered" {
		t.Fatalf("AccessToken 不匹配: got=%q", info.AccessToken)
placeholder
	if callCount < 3 {
		t.Fatalf("应至少调用 3 次（2 次失败 + 1 次成功）: got=%d", callCount)
placeholder
placeholder

// =====================
// 新增测试：GeminiOAuthService.RefreshAccountToken
// =====================

func TestGeminiOAuthService_RefreshAccountToken_NotGeminiOAuth(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeOAuth,
placeholder

	_, err := svc.RefreshAccountToken(context.Background(), account)
	if err == nil {
		t.Fatal("应返回错误（非 Gemini OAuth 账号）")
placeholder
	if !strings.Contains(err.Error(), "not a Gemini OAuth account") {
		t.Fatalf("错误信息不匹配: got=%q", err.Error())
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_NoRefreshToken(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"access_token": "at",
			"oauth_type":   "code_assist",
	placeholder,
placeholder

	_, err := svc.RefreshAccountToken(context.Background(), account)
	if err == nil {
		t.Fatal("应返回错误（无 refresh_token）")
placeholder
	if !strings.Contains(err.Error(), "no refresh token") {
		t.Fatalf("错误信息不匹配: got=%q", err.Error())
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_AIStudio(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken:  "refreshed-at",
				RefreshToken: "refreshed-rt",
				ExpiresIn:    3600,
				TokenType:    "Bearer",
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"access_token":  "old-at",
			"refresh_token": "old-rt",
			"oauth_type":    "ai_studio",
			"tier_id":       "aistudio_free",
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
	if info.AccessToken != "refreshed-at" {
		t.Fatalf("AccessToken 不匹配: got=%q", info.AccessToken)
placeholder
	if info.OAuthType != "ai_studio" {
		t.Fatalf("OAuthType 不匹配: got=%q", info.OAuthType)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_CodeAssist_WithProjectID(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken:  "refreshed",
				RefreshToken: "new-rt",
				ExpiresIn:    3600,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"access_token":  "old-at",
			"refresh_token": "old-rt",
			"oauth_type":    "code_assist",
			"project_id":    "my-project",
			"tier_id":       "gcp_standard",
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
	if info.ProjectID != "my-project" {
		t.Fatalf("ProjectID 应保留: got=%q", info.ProjectID)
placeholder
	if info.TierID != GeminiTierGCPStandard {
		t.Fatalf("TierID 不匹配: got=%q want=%q", info.TierID, GeminiTierGCPStandard)
placeholder
	if info.OAuthType != "code_assist" {
		t.Fatalf("OAuthType 不匹配: got=%q", info.OAuthType)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_DefaultOAuthType(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			if oauthType != "code_assist" {
				t.Errorf("默认 oauthType 应为 code_assist: got=%q", oauthType)
		placeholder
			return &geminicli.TokenResponse{
				AccessToken: "refreshed",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	// 无 oauth_type 凭据的旧账号
	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "old-rt",
			"project_id":    "proj",
			"tier_id":       "STANDARD",
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
	if info.OAuthType != "code_assist" {
		t.Fatalf("OAuthType 应默认为 code_assist: got=%q", info.OAuthType)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_WithProxy(t *testing.T) {
	t.Parallel()

	proxyRepo := &mockGeminiProxyRepo{
		getByIDFunc: func(ctx context.Context, id int64) (*Proxy, error) {
			return &Proxy{
				Protocol: "http",
				Host:     "proxy.test",
				Port:     3128,
		placeholder, nil
	placeholder,
placeholder

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			if proxyURL != "http://proxy.test:3128" {
				t.Errorf("proxyURL 不匹配: got=%q", proxyURL)
		placeholder
			return &geminicli.TokenResponse{
				AccessToken: "refreshed",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(proxyRepo, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	proxyID := int64(5)
	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
		ProxyID:  &proxyID,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "code_assist",
			"project_id":    "proj",
	placeholder,
placeholder

	_, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_CodeAssist_NoProjectID_AutoDetect(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken: "at",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	codeAssist := &mockGeminiCodeAssistClient{
		loadCodeAssistFunc: func(ctx context.Context, accessToken, proxyURL string, req *geminicli.LoadCodeAssistRequest) (*geminicli.LoadCodeAssistResponse, error) {
			return &geminicli.LoadCodeAssistResponse{
				CloudAICompanionProject: "auto-project-123",
				CurrentTier:             &geminicli.TierInfo{ID: "STANDARD"placeholder,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, codeAssist, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "code_assist",
			// 无 project_id，触发 fetchProjectID
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
	if info.ProjectID != "auto-project-123" {
		t.Fatalf("ProjectID 应为自动检测值: got=%q", info.ProjectID)
placeholder
	if info.TierID != GeminiTierGCPStandard {
		t.Fatalf("TierID 不匹配: got=%q", info.TierID)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_CodeAssist_NoProjectID_FailsEmpty(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken: "at",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	// 返回有 currentTier 但无 cloudaicompanionProject 的响应，
	// 使 fetchProjectID 走"已注册用户"路径（尝试 Cloud Resource Manager -> 失败 -> 返回错误），
	// 避免走 onboardUser 路径（5 次重试 x 2 秒 = 10 秒超时）
	codeAssist := &mockGeminiCodeAssistClient{
		loadCodeAssistFunc: func(ctx context.Context, accessToken, proxyURL string, req *geminicli.LoadCodeAssistRequest) (*geminicli.LoadCodeAssistResponse, error) {
			return &geminicli.LoadCodeAssistResponse{
				CurrentTier: &geminicli.TierInfo{ID: "STANDARD"placeholder,
				// 无 CloudAICompanionProject
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, codeAssist, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "code_assist",
	placeholder,
placeholder

	_, err := svc.RefreshAccountToken(context.Background(), account)
	if err == nil {
		t.Fatal("应返回错误（无法检测 project_id）")
placeholder
	if !strings.Contains(err.Error(), "project_id") {
		t.Fatalf("错误信息应包含 project_id: got=%q", err.Error())
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_GoogleOne_FreshCache(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken: "at",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "google_one",
			"project_id":    "proj",
			"tier_id":       "google_ai_pro",
	placeholder,
		Extra: map[string]any{
			// 缓存刷新时间在 24 小时内
			"drive_tier_updated_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
	// 缓存新鲜，应使用已有的 tier_id
	if info.TierID != GeminiTierGoogleAIPro {
		t.Fatalf("TierID 应使用缓存值: got=%q want=%q", info.TierID, GeminiTierGoogleAIPro)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_GoogleOne_NoTierID_DefaultsFree(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return &geminicli.TokenResponse{
				AccessToken: "at",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, &mockDriveClient{placeholder, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "google_one",
			"project_id":    "proj",
			// 无 tier_id
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 返回错误: %v", err)
placeholder
	// FetchGoogleOneTier 会被调用但 oauthClient（此处 mock）不实现 Drive API，
	// svc.FetchGoogleOneTier 使用真实 DriveClient 会失败，最终回退到默认值。
	// 由于没有 tier_id 且 FetchGoogleOneTier 失败，应默认为 google_one_free
	if info.TierID != GeminiTierGoogleOneFree {
		t.Fatalf("TierID 应为默认 free: got=%q", info.TierID)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_UnauthorizedClient_Fallback(t *testing.T) {
	t.Parallel()

	callCount := 0
	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			callCount++
			if oauthType == "code_assist" {
				return nil, fmt.Errorf("unauthorized_client: client mismatch")
		placeholder
			// ai_studio 路径成功
			return &geminicli.TokenResponse{
				AccessToken: "recovered",
				ExpiresIn:   3600,
		placeholder, nil
	placeholder,
placeholder

	// 启用自定义 OAuth 客户端以触发 fallback 路径
	cfg := &config.Config{
		Gemini: config.GeminiConfig{
			OAuth: config.GeminiOAuthConfig{
				ClientID:     "custom-id",
				ClientSecret: "custom-secret",
		placeholder,
	placeholder,
placeholder

	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, nil, cfg)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "code_assist",
			"project_id":    "proj",
			"tier_id":       "gcp_standard",
	placeholder,
placeholder

	info, err := svc.RefreshAccountToken(context.Background(), account)
	if err != nil {
		t.Fatalf("RefreshAccountToken 应在 fallback 后成功: %v", err)
placeholder
	if info.AccessToken != "recovered" {
		t.Fatalf("AccessToken 不匹配: got=%q", info.AccessToken)
placeholder
placeholder

func TestGeminiOAuthService_RefreshAccountToken_UnauthorizedClient_NoFallback(t *testing.T) {
	t.Parallel()

	client := &mockGeminiOAuthClient{
		refreshTokenFunc: func(ctx context.Context, oauthType, refreshToken, proxyURL string) (*geminicli.TokenResponse, error) {
			return nil, fmt.Errorf("unauthorized_client: client mismatch")
	placeholder,
placeholder

	// 无自定义 OAuth 客户端，无法 fallback
	svc := NewGeminiOAuthService(&mockGeminiProxyRepo{placeholder, client, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	account := &Account{
placeholder
		Type:     AccountTypeOAuth,
placeholder
			"refresh_token": "rt",
			"oauth_type":    "code_assist",
			"project_id":    "proj",
	placeholder,
placeholder

	_, err := svc.RefreshAccountToken(context.Background(), account)
	if err == nil {
		t.Fatal("应返回错误（无 fallback）")
placeholder
	if !strings.Contains(err.Error(), "OAuth client mismatch") {
		t.Fatalf("错误应包含 OAuth client mismatch: got=%q", err.Error())
placeholder
placeholder

// =====================
// 新增测试：GeminiOAuthService.ExchangeCode
// =====================

func TestGeminiOAuthService_ExchangeCode_SessionNotFound(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	_, err := svc.ExchangeCode(context.Background(), &GeminiExchangeCodeInput{
		SessionID: "nonexistent",
		State:     "some-state",
		Code:      "some-code",
placeholder)
	if err == nil {
		t.Fatal("应返回错误（session 不存在）")
placeholder
	if !strings.Contains(err.Error(), "session not found") {
		t.Fatalf("错误信息不匹配: got=%q", err.Error())
placeholder
placeholder

func TestGeminiOAuthService_ExchangeCode_InvalidState(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	// 手动创建 session（必须设置 CreatedAt，否则会因 TTL 过期被拒绝）
	svc.sessionStore.Set("test-session", &geminicli.OAuthSession{
		State:        "correct-state",
		CodeVerifier: "verifier",
		OAuthType:    "ai_studio",
		CreatedAt:    time.Now(),
placeholder)

	_, err := svc.ExchangeCode(context.Background(), &GeminiExchangeCodeInput{
		SessionID: "test-session",
		State:     "wrong-state",
		Code:      "code",
placeholder)
	if err == nil {
		t.Fatal("应返回错误（state 不匹配）")
placeholder
	if !strings.Contains(err.Error(), "invalid state") {
		t.Fatalf("错误信息不匹配: got=%q", err.Error())
placeholder
placeholder

func TestGeminiOAuthService_ExchangeCode_EmptyState(t *testing.T) {
	t.Parallel()

	svc := NewGeminiOAuthService(nil, nil, nil, nil, &config.Config{placeholder)
	defer svc.Stop()

	svc.sessionStore.Set("test-session", &geminicli.OAuthSession{
		State:        "correct-state",
		CodeVerifier: "verifier",
		CreatedAt:    time.Now(),
placeholder)

	_, err := svc.ExchangeCode(context.Background(), &GeminiExchangeCodeInput{
		SessionID: "test-session",
		State:     "",
		Code:      "code",
placeholder)
	if err == nil {
		t.Fatal("应返回错误（空 state）")
placeholder
placeholder

// =====================
// 辅助函数
// =====================

func assertCredStr(t *testing.T, creds map[string]any, key, want string) {
placeholder
	raw, ok := creds[key]
	if !ok {
		t.Fatalf("creds 缺少 key=%q", key)
placeholder
	got, ok := raw.(string)
	if !ok {
		t.Fatalf("creds[%q] 不是 string: %T", key, raw)
placeholder
	if got != want {
		t.Fatalf("creds[%q] = %q, want %q", key, got, want)
placeholder
placeholder

func credKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
placeholder
	return keys
placeholder
