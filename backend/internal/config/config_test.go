package config

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func resetViperWithJWTSecret(t *testing.T) {
placeholder
	viper.Reset()
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
placeholder

func TestLoadServerTimingConfig(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		resetViperWithJWTSecret(t)
		cfg, err := Load()
	placeholder
		require.False(t, cfg.Server.EnableServerTiming)
placeholder)

	t.Run("enabled by exact environment variable", func(t *testing.T) {
		resetViperWithJWTSecret(t)
		t.Setenv("ENABLE_SERVER_TIMING", "true")
		cfg, err := Load()
	placeholder
		require.True(t, cfg.Server.EnableServerTiming)
placeholder)
placeholder

func TestLoadHTTPIngressSafetyDefaults(t *testing.T) {
	resetViperWithJWTSecret(t)
	cfg, err := Load()
placeholder
	require.Equal(t, 10, cfg.Server.ReadHeaderTimeout)
	require.Equal(t, 64*1024, cfg.Server.MaxHeaderBytes)
	require.Equal(t, []string{
		"127.0.0.0/8",
		"::1/128",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
placeholder, cfg.Server.TrustedProxies)
	require.Equal(t, int64(32*1024*1024), cfg.Gateway.TextMaxBodySize)
	require.True(t, cfg.APIKeyAuth.InvalidAbuse.Enabled)
	require.Equal(t, 120, cfg.APIKeyAuth.InvalidAbuse.Threshold)
	require.Equal(t, 16384, cfg.APIKeyAuth.InvalidAbuse.Capacity)
placeholder

func TestLoadExplicitEmptyTrustedProxiesKeepsLegacyDefault(t *testing.T) {
	resetViperWithJWTSecret(t)
	viper.Set("server.trusted_proxies", []string{placeholder)

	cfg, err := Load()
placeholder
	require.Empty(t, cfg.Server.TrustedProxies)
placeholder

func TestLoadForBootstrapAllowsMissingJWTSecret(t *testing.T) {
	viper.Reset()
	t.Setenv("JWT_SECRET", "")

	cfg, err := LoadForBootstrap()
	if err != nil {
		t.Fatalf("LoadForBootstrap() error: %v", err)
placeholder
	if cfg.JWT.Secret != "" {
		t.Fatalf("LoadForBootstrap() should keep empty jwt.secret during bootstrap")
placeholder
placeholder

func TestNormalizeRunMode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
placeholder{
		{"simple", "simple"placeholder,
		{"SIMPLE", "simple"placeholder,
		{"standard", "standard"placeholder,
		{"invalid", "standard"placeholder,
		{"", "standard"placeholder,
placeholder

	for _, tt := range tests {
		result := NormalizeRunMode(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeRunMode(%q) = %q, want %q", tt.input, result, tt.expected)
	placeholder
placeholder
placeholder

func TestLoadDefaultSchedulingConfig(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.Gateway.Scheduling.StickySessionMaxWaiting != 3 {
		t.Fatalf("StickySessionMaxWaiting = %d, want 3", cfg.Gateway.Scheduling.StickySessionMaxWaiting)
placeholder
	if cfg.Gateway.Scheduling.StickySessionWaitTimeout != 120*time.Second {
		t.Fatalf("StickySessionWaitTimeout = %v, want 120s", cfg.Gateway.Scheduling.StickySessionWaitTimeout)
placeholder
	if cfg.Gateway.Scheduling.FallbackWaitTimeout != 30*time.Second {
		t.Fatalf("FallbackWaitTimeout = %v, want 30s", cfg.Gateway.Scheduling.FallbackWaitTimeout)
placeholder
	if cfg.Gateway.Scheduling.FallbackMaxWaiting != 100 {
		t.Fatalf("FallbackMaxWaiting = %d, want 100", cfg.Gateway.Scheduling.FallbackMaxWaiting)
placeholder
	if !cfg.Gateway.Scheduling.LoadBatchEnabled {
		t.Fatalf("LoadBatchEnabled = false, want true")
placeholder
	if cfg.Gateway.Scheduling.LoadBatchCacheTTLMS != 200 {
		t.Fatalf("LoadBatchCacheTTLMS = %d, want 200", cfg.Gateway.Scheduling.LoadBatchCacheTTLMS)
placeholder
	if cfg.Gateway.Scheduling.SlotCleanupInterval != 30*time.Second {
		t.Fatalf("SlotCleanupInterval = %v, want 30s", cfg.Gateway.Scheduling.SlotCleanupInterval)
placeholder
placeholder

func TestLoadDefaultOpenAIFirstOutputTimeoutsDisabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
placeholder
	require.Zero(t, cfg.Gateway.OpenAIFirstOutputTimeoutSeconds)
	require.Zero(t, cfg.Gateway.OpenAIHighEffortFirstOutputTimeoutSeconds)
placeholder

func TestLoadOpenAIFirstOutputTimeoutsFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_OPENAI_FIRST_OUTPUT_TIMEOUT_SECONDS", "90")
	t.Setenv("GATEWAY_OPENAI_HIGH_EFFORT_FIRST_OUTPUT_TIMEOUT_SECONDS", "240")

	cfg, err := Load()
placeholder
	require.Equal(t, 90, cfg.Gateway.OpenAIFirstOutputTimeoutSeconds)
	require.Equal(t, 240, cfg.Gateway.OpenAIHighEffortFirstOutputTimeoutSeconds)
placeholder

func TestValidateOpenAIFirstOutputTimeoutMinimum(t *testing.T) {
	resetViperWithJWTSecret(t)
	cfg, err := Load()
placeholder
	cfg.Gateway.OpenAIFirstOutputTimeoutSeconds = 30
	require.NoError(t, cfg.Validate())
placeholder

func TestLoadDefaultOpenAIWSConfig(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if !cfg.Gateway.OpenAIWS.Enabled {
		t.Fatalf("Gateway.OpenAIWS.Enabled = false, want true")
placeholder
	if !cfg.Gateway.OpenAIWS.ResponsesWebsocketsV2 {
		t.Fatalf("Gateway.OpenAIWS.ResponsesWebsocketsV2 = false, want true")
placeholder
	if cfg.Gateway.OpenAIWS.ResponsesWebsockets {
		t.Fatalf("Gateway.OpenAIWS.ResponsesWebsockets = true, want false")
placeholder
	if !cfg.Gateway.OpenAIWS.DynamicMaxConnsByAccountConcurrencyEnabled {
		t.Fatalf("Gateway.OpenAIWS.DynamicMaxConnsByAccountConcurrencyEnabled = false, want true")
placeholder
	if cfg.Gateway.OpenAIWS.OAuthMaxConnsFactor != 1.0 {
		t.Fatalf("Gateway.OpenAIWS.OAuthMaxConnsFactor = %v, want 1.0", cfg.Gateway.OpenAIWS.OAuthMaxConnsFactor)
placeholder
	if cfg.Gateway.OpenAIWS.APIKeyMaxConnsFactor != 1.0 {
		t.Fatalf("Gateway.OpenAIWS.APIKeyMaxConnsFactor = %v, want 1.0", cfg.Gateway.OpenAIWS.APIKeyMaxConnsFactor)
placeholder
	if cfg.Gateway.OpenAIWS.StickySessionTTLSeconds != 3600 {
		t.Fatalf("Gateway.OpenAIWS.StickySessionTTLSeconds = %d, want 3600", cfg.Gateway.OpenAIWS.StickySessionTTLSeconds)
placeholder
	if !cfg.Gateway.OpenAIScheduler.StickyEscapeEnabled {
		t.Fatalf("Gateway.OpenAIScheduler.StickyEscapeEnabled = false, want true")
placeholder
	if cfg.Gateway.OpenAIScheduler.StickyEscapeTTFTMs != 15000 {
		t.Fatalf("Gateway.OpenAIScheduler.StickyEscapeTTFTMs = %d, want 15000", cfg.Gateway.OpenAIScheduler.StickyEscapeTTFTMs)
placeholder
	if cfg.Gateway.OpenAIScheduler.StickyEscapeErrorRate != 0.5 {
		t.Fatalf("Gateway.OpenAIScheduler.StickyEscapeErrorRate = %v, want 0.5", cfg.Gateway.OpenAIScheduler.StickyEscapeErrorRate)
placeholder
	if !cfg.Gateway.OpenAIWS.SessionHashReadOldFallback {
		t.Fatalf("Gateway.OpenAIWS.SessionHashReadOldFallback = false, want true")
placeholder
	if !cfg.Gateway.OpenAIWS.SessionHashDualWriteOld {
		t.Fatalf("Gateway.OpenAIWS.SessionHashDualWriteOld = false, want true")
placeholder
	if !cfg.Gateway.OpenAIWS.MetadataBridgeEnabled {
		t.Fatalf("Gateway.OpenAIWS.MetadataBridgeEnabled = false, want true")
placeholder
	if cfg.Gateway.OpenAIWS.StickyResponseIDTTLSeconds != 3600 {
		t.Fatalf("Gateway.OpenAIWS.StickyResponseIDTTLSeconds = %d, want 3600", cfg.Gateway.OpenAIWS.StickyResponseIDTTLSeconds)
placeholder
	if cfg.Gateway.OpenAIWS.FallbackCooldownSeconds != 30 {
		t.Fatalf("Gateway.OpenAIWS.FallbackCooldownSeconds = %d, want 30", cfg.Gateway.OpenAIWS.FallbackCooldownSeconds)
placeholder
	if cfg.Gateway.OpenAIWS.EventFlushBatchSize != 1 {
		t.Fatalf("Gateway.OpenAIWS.EventFlushBatchSize = %d, want 1", cfg.Gateway.OpenAIWS.EventFlushBatchSize)
placeholder
	if cfg.Gateway.OpenAIWS.EventFlushIntervalMS != 10 {
		t.Fatalf("Gateway.OpenAIWS.EventFlushIntervalMS = %d, want 10", cfg.Gateway.OpenAIWS.EventFlushIntervalMS)
placeholder
	if cfg.Gateway.OpenAIWS.PrewarmCooldownMS != 300 {
		t.Fatalf("Gateway.OpenAIWS.PrewarmCooldownMS = %d, want 300", cfg.Gateway.OpenAIWS.PrewarmCooldownMS)
placeholder
	if cfg.Gateway.OpenAIWS.ClientReadLimitBytes != 64*1024*1024 {
		t.Fatalf("Gateway.OpenAIWS.ClientReadLimitBytes = %d, want %d", cfg.Gateway.OpenAIWS.ClientReadLimitBytes, 64*1024*1024)
placeholder
	if !cfg.Gateway.OpenAIWS.HTTPBridgeEnabled {
		t.Fatalf("Gateway.OpenAIWS.HTTPBridgeEnabled = false, want true")
placeholder
	if cfg.Gateway.OpenAIWS.HTTPBridgeThresholdBytes != 15*1024*1024 {
		t.Fatalf("Gateway.OpenAIWS.HTTPBridgeThresholdBytes = %d, want %d", cfg.Gateway.OpenAIWS.HTTPBridgeThresholdBytes, 15*1024*1024)
placeholder
	if cfg.Gateway.OpenAIWS.RetryBackoffInitialMS != 120 {
		t.Fatalf("Gateway.OpenAIWS.RetryBackoffInitialMS = %d, want 120", cfg.Gateway.OpenAIWS.RetryBackoffInitialMS)
placeholder
	if cfg.Gateway.OpenAIWS.RetryBackoffMaxMS != 2000 {
		t.Fatalf("Gateway.OpenAIWS.RetryBackoffMaxMS = %d, want 2000", cfg.Gateway.OpenAIWS.RetryBackoffMaxMS)
placeholder
	if cfg.Gateway.OpenAIWS.RetryJitterRatio != 0.2 {
		t.Fatalf("Gateway.OpenAIWS.RetryJitterRatio = %v, want 0.2", cfg.Gateway.OpenAIWS.RetryJitterRatio)
placeholder
	if cfg.Gateway.OpenAIWS.RetryTotalBudgetMS != 5000 {
		t.Fatalf("Gateway.OpenAIWS.RetryTotalBudgetMS = %d, want 5000", cfg.Gateway.OpenAIWS.RetryTotalBudgetMS)
placeholder
	if cfg.Gateway.OpenAIWS.PayloadLogSampleRate != 0.2 {
		t.Fatalf("Gateway.OpenAIWS.PayloadLogSampleRate = %v, want 0.2", cfg.Gateway.OpenAIWS.PayloadLogSampleRate)
placeholder
	if cfg.Gateway.OpenAIWS.SchedulerScoreWeights.QuotaHeadroom != 0 {
		t.Fatalf("Gateway.OpenAIWS.SchedulerScoreWeights.QuotaHeadroom = %v, want 0", cfg.Gateway.OpenAIWS.SchedulerScoreWeights.QuotaHeadroom)
placeholder
	if cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost != 0 {
		t.Fatalf("Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = %v, want 0", cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost)
placeholder
	if !cfg.Gateway.OpenAIWS.StoreDisabledForceNewConn {
		t.Fatalf("Gateway.OpenAIWS.StoreDisabledForceNewConn = false, want true")
placeholder
	if cfg.Gateway.OpenAIWS.StoreDisabledConnMode != "strict" {
		t.Fatalf("Gateway.OpenAIWS.StoreDisabledConnMode = %q, want %q", cfg.Gateway.OpenAIWS.StoreDisabledConnMode, "strict")
placeholder
	if cfg.Gateway.OpenAIWS.ModeRouterV2Enabled {
		t.Fatalf("Gateway.OpenAIWS.ModeRouterV2Enabled = true, want false")
placeholder
	if cfg.Gateway.OpenAIWS.IngressModeDefault != "ctx_pool" {
		t.Fatalf("Gateway.OpenAIWS.IngressModeDefault = %q, want %q", cfg.Gateway.OpenAIWS.IngressModeDefault, "ctx_pool")
placeholder
	if cfg.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds != DefaultOpenAIWSClientFirstMessageTimeoutSeconds {
		t.Fatalf(
			"Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds = %d, want %d",
			cfg.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds,
			DefaultOpenAIWSClientFirstMessageTimeoutSeconds,
		)
placeholder
	if cfg.Gateway.OpenAIWS.IngressInterTurnIdleTimeoutSeconds != 300 {
		t.Fatalf("Gateway.OpenAIWS.IngressInterTurnIdleTimeoutSeconds = %d, want 300", cfg.Gateway.OpenAIWS.IngressInterTurnIdleTimeoutSeconds)
placeholder
	if cfg.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey != 64 {
		t.Fatalf("Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey = %d, want 64", cfg.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey)
placeholder
placeholder

func TestLoadOpenAIWSClientFirstMessageTimeoutFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_OPENAI_WS_CLIENT_FIRST_MESSAGE_TIMEOUT_SECONDS", "120")

	cfg, err := Load()
placeholder
	require.Equal(t, 120, cfg.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds)
placeholder

func TestLoadDefaultOpenAICompactModel(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
placeholder
	require.Equal(t, "gpt-5.4", cfg.Gateway.OpenAICompactModel)
placeholder

func TestLoadOpenAICompactModelFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_OPENAI_COMPACT_MODEL", "gpt-5.3-codex")

	cfg, err := Load()
placeholder
	require.Equal(t, "gpt-5.3-codex", cfg.Gateway.OpenAICompactModel)
placeholder

func TestLoadDefaultOpenAIHTTP2Enabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
placeholder
	require.True(t, cfg.Gateway.OpenAIHTTP2.Enabled)
	require.True(t, cfg.Gateway.OpenAIHTTP2.AllowProxyFallbackToHTTP1)
placeholder

func TestLoadOpenAIHTTP2DisabledFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_OPENAI_HTTP2_ENABLED", "false")

	cfg, err := Load()
placeholder
	require.False(t, cfg.Gateway.OpenAIHTTP2.Enabled)
placeholder

func TestLoadDefaultOpenAIResponseHeaderTimeoutUnlimited(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
placeholder
	require.Equal(t, 0, cfg.Gateway.OpenAIResponseHeaderTimeout)
placeholder

func TestLoadOpenAIResponseHeaderTimeoutFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_OPENAI_RESPONSE_HEADER_TIMEOUT", "1800")

	cfg, err := Load()
placeholder
	require.Equal(t, 1800, cfg.Gateway.OpenAIResponseHeaderTimeout)
placeholder

func TestLoadImageNonstreamKeepaliveFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_IMAGE_NONSTREAM_KEEPALIVE_INTERVAL", "15")

	cfg, err := Load()
placeholder
	require.Equal(t, 15, cfg.Gateway.ImageNonstreamKeepaliveInterval)
placeholder

func TestLoadOpenAIWSStickyTTLCompatibility(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_OPENAI_WS_STICKY_RESPONSE_ID_TTL_SECONDS", "0")
	t.Setenv("GATEWAY_OPENAI_WS_STICKY_PREVIOUS_RESPONSE_TTL_SECONDS", "7200")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.Gateway.OpenAIWS.StickyResponseIDTTLSeconds != 7200 {
		t.Fatalf("StickyResponseIDTTLSeconds = %d, want 7200", cfg.Gateway.OpenAIWS.StickyResponseIDTTLSeconds)
placeholder
placeholder

func TestLoadDefaultIdempotencyConfig(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if !cfg.Idempotency.ObserveOnly {
		t.Fatalf("Idempotency.ObserveOnly = false, want true")
placeholder
	if cfg.Idempotency.DefaultTTLSeconds != 86400 {
		t.Fatalf("Idempotency.DefaultTTLSeconds = %d, want 86400", cfg.Idempotency.DefaultTTLSeconds)
placeholder
	if cfg.Idempotency.SystemOperationTTLSeconds != 3600 {
		t.Fatalf("Idempotency.SystemOperationTTLSeconds = %d, want 3600", cfg.Idempotency.SystemOperationTTLSeconds)
placeholder
placeholder

func TestLoadDefaultBatchImageQueueDisabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
placeholder
	require.False(t, cfg.BatchImage.QueueEnabled)
placeholder

func TestLoadIdempotencyConfigFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("IDEMPOTENCY_OBSERVE_ONLY", "false")
	t.Setenv("IDEMPOTENCY_DEFAULT_TTL_SECONDS", "600")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder
	if cfg.Idempotency.ObserveOnly {
		t.Fatalf("Idempotency.ObserveOnly = true, want false")
placeholder
	if cfg.Idempotency.DefaultTTLSeconds != 600 {
		t.Fatalf("Idempotency.DefaultTTLSeconds = %d, want 600", cfg.Idempotency.DefaultTTLSeconds)
placeholder
placeholder

func TestLoadSchedulingConfigFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("GATEWAY_SCHEDULING_STICKY_SESSION_MAX_WAITING", "5")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.Gateway.Scheduling.StickySessionMaxWaiting != 5 {
		t.Fatalf("StickySessionMaxWaiting = %d, want 5", cfg.Gateway.Scheduling.StickySessionMaxWaiting)
placeholder
placeholder

func TestLoadWeChatConnectConfigFromLegacyEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("WECHAT_OAUTH_OPEN_APP_ID", "wx-open-app")
	t.Setenv("WECHAT_OAUTH_OPEN_APP_SECRET", "wx-open-secret")
	t.Setenv("WECHAT_OAUTH_MP_APP_ID", "wx-mp-app")
	t.Setenv("WECHAT_OAUTH_MP_APP_SECRET", "wx-mp-secret")
	t.Setenv("WECHAT_OAUTH_FRONTEND_REDIRECT_URL", "/auth/wechat/legacy-callback")

	cfg, err := Load()
placeholder
	require.True(t, cfg.WeChat.Enabled)
	require.True(t, cfg.WeChat.OpenEnabled)
	require.True(t, cfg.WeChat.MPEnabled)
	require.False(t, cfg.WeChat.MobileEnabled)
	require.Equal(t, "open", cfg.WeChat.Mode)
	require.Equal(t, "wx-open-app", cfg.WeChat.OpenAppID)
	require.Equal(t, "wx-open-secret", cfg.WeChat.OpenAppSecret)
	require.Equal(t, "wx-mp-app", cfg.WeChat.MPAppID)
	require.Equal(t, "wx-mp-secret", cfg.WeChat.MPAppSecret)
	require.Equal(t, "/auth/wechat/legacy-callback", cfg.WeChat.FrontendRedirectURL)
placeholder

func TestLoadDefaultOIDCSecurityDefaults(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
placeholder
	require.True(t, cfg.OIDC.UsePKCE)
	require.True(t, cfg.OIDC.ValidateIDToken)
	require.False(t, cfg.OIDC.UsePKCEExplicit)
	require.False(t, cfg.OIDC.ValidateIDTokenExplicit)
placeholder

func TestLoadExplicitOIDCSecurityDefaultsFromEnvMarksFlagsExplicit(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("OIDC_CONNECT_USE_PKCE", "false")
	t.Setenv("OIDC_CONNECT_VALIDATE_ID_TOKEN", "false")

	cfg, err := Load()
placeholder
	require.False(t, cfg.OIDC.UsePKCE)
	require.False(t, cfg.OIDC.ValidateIDToken)
	require.True(t, cfg.OIDC.UsePKCEExplicit)
	require.True(t, cfg.OIDC.ValidateIDTokenExplicit)
placeholder

func TestLoadForcedCodexInstructionsTemplate(t *testing.T) {
	resetViperWithJWTSecret(t)

	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "codex-instructions.md.tmpl")
	configPath := filepath.Join(tempDir, "config.yaml")

	require.NoError(t, os.WriteFile(templatePath, []byte("server-prefix\n\n{{ .ExistingInstructions placeholderplaceholder"), 0o644))
	yamlSafePath := filepath.ToSlash(templatePath)
	require.NoError(t, os.WriteFile(configPath, []byte("gateway:\n  forced_codex_instructions_template_file: \""+yamlSafePath+"\"\n"), 0o644))
	t.Setenv("DATA_DIR", tempDir)

	cfg, err := Load()
placeholder
	require.Equal(t, yamlSafePath, cfg.Gateway.ForcedCodexInstructionsTemplateFile)
	require.Equal(t, "server-prefix\n\n{{ .ExistingInstructions placeholderplaceholder", cfg.Gateway.ForcedCodexInstructionsTemplate)
placeholder

func TestLoadDefaultSecurityToggles(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.Security.URLAllowlist.Enabled {
		t.Fatalf("URLAllowlist.Enabled = true, want false")
placeholder
	if !cfg.Security.URLAllowlist.AllowInsecureHTTP {
		t.Fatalf("URLAllowlist.AllowInsecureHTTP = false, want true")
placeholder
	if !cfg.Security.URLAllowlist.AllowPrivateHosts {
		t.Fatalf("URLAllowlist.AllowPrivateHosts = false, want true")
placeholder
	if !cfg.Security.ResponseHeaders.Enabled {
		t.Fatalf("ResponseHeaders.Enabled = false, want true")
placeholder
placeholder

func TestLoadDefaultServerMode(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.Server.Mode != "release" {
		t.Fatalf("Server.Mode = %q, want %q", cfg.Server.Mode, "release")
placeholder
placeholder

func TestLoadDefaultJWTAccessTokenExpireMinutes(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.JWT.ExpireHour != 24 {
		t.Fatalf("JWT.ExpireHour = %d, want 24", cfg.JWT.ExpireHour)
placeholder
	if cfg.JWT.AccessTokenExpireMinutes != 0 {
		t.Fatalf("JWT.AccessTokenExpireMinutes = %d, want 0", cfg.JWT.AccessTokenExpireMinutes)
placeholder
placeholder

func TestLoadJWTAccessTokenExpireMinutesFromEnv(t *testing.T) {
	resetViperWithJWTSecret(t)
	t.Setenv("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", "90")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.JWT.AccessTokenExpireMinutes != 90 {
		t.Fatalf("JWT.AccessTokenExpireMinutes = %d, want 90", cfg.JWT.AccessTokenExpireMinutes)
placeholder
placeholder

func TestLoadDefaultDatabaseSSLMode(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if cfg.Database.SSLMode != "prefer" {
		t.Fatalf("Database.SSLMode = %q, want %q", cfg.Database.SSLMode, "prefer")
placeholder
placeholder

func TestValidateLinuxDoFrontendRedirectURL(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.LinuxDo.Enabled = true
	cfg.LinuxDo.ClientID = "test-client"
	cfg.LinuxDo.ClientSecret = "test-secret"
	cfg.LinuxDo.RedirectURL = "https://example.com/api/v1/auth/oauth/linuxdo/callback"
	cfg.LinuxDo.TokenAuthMethod = "client_secret_post"
	cfg.LinuxDo.UsePKCE = true

	cfg.LinuxDo.FrontendRedirectURL = "javascript:alert(1)"
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for javascript scheme, got nil")
placeholder
	if !strings.Contains(err.Error(), "linuxdo_connect.frontend_redirect_url") {
		t.Fatalf("Validate() expected frontend_redirect_url error, got: %v", err)
placeholder
placeholder

func TestValidateLinuxDoAllowsDisablingPKCEForCompatibility(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.LinuxDo.Enabled = true
	cfg.LinuxDo.ClientID = "test-client"
	cfg.LinuxDo.ClientSecret = ""
	cfg.LinuxDo.RedirectURL = "https://example.com/api/v1/auth/oauth/linuxdo/callback"
	cfg.LinuxDo.FrontendRedirectURL = "/auth/linuxdo/callback"
	cfg.LinuxDo.TokenAuthMethod = "none"
	cfg.LinuxDo.UsePKCE = false

	err = cfg.Validate()
	if err != nil {
		t.Fatalf("Validate() expected LinuxDo config without PKCE to pass for compatibility, got: %v", err)
placeholder
placeholder

func TestValidateOIDCScopesMustContainOpenID(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.OIDC.Enabled = true
	cfg.OIDC.ClientID = "oidc-client"
	cfg.OIDC.ClientSecret = "oidc-secret"
	cfg.OIDC.IssuerURL = "https://issuer.example.com"
	cfg.OIDC.AuthorizeURL = "https://issuer.example.com/auth"
	cfg.OIDC.TokenURL = "https://issuer.example.com/token"
	cfg.OIDC.JWKSURL = "https://issuer.example.com/jwks"
	cfg.OIDC.RedirectURL = "https://example.com/api/v1/auth/oauth/oidc/callback"
	cfg.OIDC.FrontendRedirectURL = "/auth/oidc/callback"
	cfg.OIDC.Scopes = "profile email"
	cfg.OIDC.UsePKCE = true

	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error when scopes do not include openid, got nil")
placeholder
	if !strings.Contains(err.Error(), "oidc_connect.scopes") {
		t.Fatalf("Validate() expected oidc_connect.scopes error, got: %v", err)
placeholder
placeholder

func TestValidateOIDCAllowsIssuerOnlyEndpointsWithDiscoveryFallback(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.OIDC.Enabled = true
	cfg.OIDC.ClientID = "oidc-client"
	cfg.OIDC.ClientSecret = "oidc-secret"
	cfg.OIDC.IssuerURL = "https://issuer.example.com"
	cfg.OIDC.AuthorizeURL = ""
	cfg.OIDC.TokenURL = ""
	cfg.OIDC.JWKSURL = ""
	cfg.OIDC.RedirectURL = "https://example.com/api/v1/auth/oauth/oidc/callback"
	cfg.OIDC.FrontendRedirectURL = "/auth/oidc/callback"
	cfg.OIDC.Scopes = "openid email profile"
	cfg.OIDC.ValidateIDToken = true
	cfg.OIDC.UsePKCE = true

	err = cfg.Validate()
	if err != nil {
		t.Fatalf("Validate() expected issuer-only OIDC config to pass with discovery fallback, got: %v", err)
placeholder
placeholder

func TestValidateOIDCAllowsExplicitCompatibilityOverridesForPKCEAndIDTokenValidation(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.OIDC.Enabled = true
	cfg.OIDC.ClientID = "oidc-client"
	cfg.OIDC.ClientSecret = "oidc-secret"
	cfg.OIDC.IssuerURL = "https://issuer.example.com"
	cfg.OIDC.AuthorizeURL = "https://issuer.example.com/auth"
	cfg.OIDC.TokenURL = "https://issuer.example.com/token"
	cfg.OIDC.UserInfoURL = "https://issuer.example.com/userinfo"
	cfg.OIDC.RedirectURL = "https://example.com/api/v1/auth/oauth/oidc/callback"
	cfg.OIDC.FrontendRedirectURL = "/auth/oidc/callback"
	cfg.OIDC.Scopes = "openid email profile"
	cfg.OIDC.UsePKCE = false
	cfg.OIDC.ValidateIDToken = false
	cfg.OIDC.JWKSURL = ""
	cfg.OIDC.AllowedSigningAlgs = ""

	err = cfg.Validate()
	if err != nil {
		t.Fatalf("Validate() expected OIDC config without PKCE/id_token validation to pass for compatibility, got: %v", err)
placeholder
placeholder

func TestLoadDefaultDashboardCacheConfig(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if !cfg.Dashboard.Enabled {
		t.Fatalf("Dashboard.Enabled = false, want true")
placeholder
	if cfg.Dashboard.KeyPrefix != "sub2api:" {
		t.Fatalf("Dashboard.KeyPrefix = %q, want %q", cfg.Dashboard.KeyPrefix, "sub2api:")
placeholder
	if cfg.Dashboard.StatsFreshTTLSeconds != 15 {
		t.Fatalf("Dashboard.StatsFreshTTLSeconds = %d, want 15", cfg.Dashboard.StatsFreshTTLSeconds)
placeholder
	if cfg.Dashboard.StatsTTLSeconds != 30 {
		t.Fatalf("Dashboard.StatsTTLSeconds = %d, want 30", cfg.Dashboard.StatsTTLSeconds)
placeholder
	if cfg.Dashboard.StatsRefreshTimeoutSeconds != 30 {
		t.Fatalf("Dashboard.StatsRefreshTimeoutSeconds = %d, want 30", cfg.Dashboard.StatsRefreshTimeoutSeconds)
placeholder
placeholder

func TestValidateDashboardCacheConfigEnabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Dashboard.Enabled = true
	cfg.Dashboard.StatsFreshTTLSeconds = 10
	cfg.Dashboard.StatsTTLSeconds = 5
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for stats_fresh_ttl_seconds > stats_ttl_seconds, got nil")
placeholder
	if !strings.Contains(err.Error(), "dashboard_cache.stats_fresh_ttl_seconds") {
		t.Fatalf("Validate() expected stats_fresh_ttl_seconds error, got: %v", err)
placeholder
placeholder

func TestValidateDashboardCacheConfigDisabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Dashboard.Enabled = false
	cfg.Dashboard.StatsTTLSeconds = -1
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for negative stats_ttl_seconds, got nil")
placeholder
	if !strings.Contains(err.Error(), "dashboard_cache.stats_ttl_seconds") {
		t.Fatalf("Validate() expected stats_ttl_seconds error, got: %v", err)
placeholder
placeholder

func TestLoadDefaultDashboardAggregationConfig(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if !cfg.DashboardAgg.Enabled {
		t.Fatalf("DashboardAgg.Enabled = false, want true")
placeholder
	if cfg.DashboardAgg.IntervalSeconds != 60 {
		t.Fatalf("DashboardAgg.IntervalSeconds = %d, want 60", cfg.DashboardAgg.IntervalSeconds)
placeholder
	if cfg.DashboardAgg.LookbackSeconds != 120 {
		t.Fatalf("DashboardAgg.LookbackSeconds = %d, want 120", cfg.DashboardAgg.LookbackSeconds)
placeholder
	if cfg.DashboardAgg.BackfillEnabled {
		t.Fatalf("DashboardAgg.BackfillEnabled = true, want false")
placeholder
	if cfg.DashboardAgg.BackfillMaxDays != 31 {
		t.Fatalf("DashboardAgg.BackfillMaxDays = %d, want 31", cfg.DashboardAgg.BackfillMaxDays)
placeholder
	if cfg.DashboardAgg.Retention.UsageLogsDays != 90 {
		t.Fatalf("DashboardAgg.Retention.UsageLogsDays = %d, want 90", cfg.DashboardAgg.Retention.UsageLogsDays)
placeholder
	if cfg.DashboardAgg.Retention.UsageBillingDedupDays != 365 {
		t.Fatalf("DashboardAgg.Retention.UsageBillingDedupDays = %d, want 365", cfg.DashboardAgg.Retention.UsageBillingDedupDays)
placeholder
	if cfg.DashboardAgg.Retention.HourlyDays != 180 {
		t.Fatalf("DashboardAgg.Retention.HourlyDays = %d, want 180", cfg.DashboardAgg.Retention.HourlyDays)
placeholder
	if cfg.DashboardAgg.Retention.DailyDays != 730 {
		t.Fatalf("DashboardAgg.Retention.DailyDays = %d, want 730", cfg.DashboardAgg.Retention.DailyDays)
placeholder
	if cfg.DashboardAgg.RecomputeDays != 2 {
		t.Fatalf("DashboardAgg.RecomputeDays = %d, want 2", cfg.DashboardAgg.RecomputeDays)
placeholder
placeholder

func TestValidateDashboardAggregationConfigDisabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.DashboardAgg.Enabled = false
	cfg.DashboardAgg.IntervalSeconds = -1
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for negative dashboard_aggregation.interval_seconds, got nil")
placeholder
	if !strings.Contains(err.Error(), "dashboard_aggregation.interval_seconds") {
		t.Fatalf("Validate() expected interval_seconds error, got: %v", err)
placeholder
placeholder

func TestValidateDashboardAggregationBackfillMaxDays(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.DashboardAgg.BackfillEnabled = true
	cfg.DashboardAgg.BackfillMaxDays = 0
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for dashboard_aggregation.backfill_max_days, got nil")
placeholder
	if !strings.Contains(err.Error(), "dashboard_aggregation.backfill_max_days") {
		t.Fatalf("Validate() expected backfill_max_days error, got: %v", err)
placeholder
placeholder

func TestLoadDefaultUsageCleanupConfig(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if !cfg.UsageCleanup.Enabled {
		t.Fatalf("UsageCleanup.Enabled = false, want true")
placeholder
	if cfg.UsageCleanup.MaxRangeDays != 31 {
		t.Fatalf("UsageCleanup.MaxRangeDays = %d, want 31", cfg.UsageCleanup.MaxRangeDays)
placeholder
	if cfg.UsageCleanup.BatchSize != 5000 {
		t.Fatalf("UsageCleanup.BatchSize = %d, want 5000", cfg.UsageCleanup.BatchSize)
placeholder
	if cfg.UsageCleanup.WorkerIntervalSeconds != 10 {
		t.Fatalf("UsageCleanup.WorkerIntervalSeconds = %d, want 10", cfg.UsageCleanup.WorkerIntervalSeconds)
placeholder
	if cfg.UsageCleanup.TaskTimeoutSeconds != 1800 {
		t.Fatalf("UsageCleanup.TaskTimeoutSeconds = %d, want 1800", cfg.UsageCleanup.TaskTimeoutSeconds)
placeholder
placeholder

func TestValidateUsageCleanupConfigEnabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.UsageCleanup.Enabled = true
	cfg.UsageCleanup.MaxRangeDays = 0
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for usage_cleanup.max_range_days, got nil")
placeholder
	if !strings.Contains(err.Error(), "usage_cleanup.max_range_days") {
		t.Fatalf("Validate() expected max_range_days error, got: %v", err)
placeholder
placeholder

func TestValidateUsageCleanupConfigDisabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.UsageCleanup.Enabled = false
	cfg.UsageCleanup.BatchSize = -1
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for usage_cleanup.batch_size, got nil")
placeholder
	if !strings.Contains(err.Error(), "usage_cleanup.batch_size") {
		t.Fatalf("Validate() expected batch_size error, got: %v", err)
placeholder
placeholder

func TestConfigAddressHelpers(t *testing.T) {
	server := ServerConfig{Host: "127.0.0.1", Port: 9000placeholder
	if server.Address() != "127.0.0.1:9000" {
		t.Fatalf("ServerConfig.Address() = %q", server.Address())
placeholder

	dbCfg := DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DBName:   "sub2api",
		SSLMode:  "disable",
placeholder
	if !strings.Contains(dbCfg.DSN(), "password=") {
placeholder else {
		t.Fatalf("DatabaseConfig.DSN() should not include password when empty")
placeholder

	dbCfg.Password = "secret"
	if !strings.Contains(dbCfg.DSN(), "password=secret") {
		t.Fatalf("DatabaseConfig.DSN() missing password")
placeholder

	dbCfg.Password = ""
	if strings.Contains(dbCfg.DSNWithTimezone("UTC"), "password=") {
		t.Fatalf("DatabaseConfig.DSNWithTimezone() should omit password when empty")
placeholder

	if !strings.Contains(dbCfg.DSNWithTimezone(""), "TimeZone=Asia/Shanghai") {
		t.Fatalf("DatabaseConfig.DSNWithTimezone() should use default timezone")
placeholder
	if !strings.Contains(dbCfg.DSNWithTimezone("UTC"), "TimeZone=UTC") {
		t.Fatalf("DatabaseConfig.DSNWithTimezone() should use provided timezone")
placeholder

	redis := RedisConfig{Host: "redis", Port: 6379placeholder
	if redis.Address() != "redis:6379" {
		t.Fatalf("RedisConfig.Address() = %q", redis.Address())
placeholder
placeholder

func TestNormalizeStringSlice(t *testing.T) {
	values := normalizeStringSlice([]string{" a ", "", "b", "   ", "c"placeholder)
	if len(values) != 3 || values[0] != "a" || values[1] != "b" || values[2] != "c" {
		t.Fatalf("normalizeStringSlice() unexpected result: %#v", values)
placeholder
	if normalizeStringSlice(nil) != nil {
		t.Fatalf("normalizeStringSlice(nil) expected nil slice")
placeholder
placeholder

func TestGetServerAddressFromEnv(t *testing.T) {
	t.Setenv("SERVER_HOST", "127.0.0.1")
	t.Setenv("SERVER_PORT", "9090")

	address := GetServerAddress()
	if address != "127.0.0.1:9090" {
		t.Fatalf("GetServerAddress() = %q", address)
placeholder
placeholder

func TestValidateAbsoluteHTTPURL(t *testing.T) {
	if err := ValidateAbsoluteHTTPURL("https://example.com/path"); err != nil {
		t.Fatalf("ValidateAbsoluteHTTPURL valid url error: %v", err)
placeholder
	if err := ValidateAbsoluteHTTPURL(""); err == nil {
		t.Fatalf("ValidateAbsoluteHTTPURL should reject empty url")
placeholder
	if err := ValidateAbsoluteHTTPURL("/relative"); err == nil {
		t.Fatalf("ValidateAbsoluteHTTPURL should reject relative url")
placeholder
	if err := ValidateAbsoluteHTTPURL("ftp://example.com"); err == nil {
		t.Fatalf("ValidateAbsoluteHTTPURL should reject ftp scheme")
placeholder
	if err := ValidateAbsoluteHTTPURL("https://example.com/#frag"); err == nil {
		t.Fatalf("ValidateAbsoluteHTTPURL should reject fragment")
placeholder
placeholder

func TestValidateServerFrontendURL(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Server.FrontendURL = "https://example.com"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() frontend_url valid error: %v", err)
placeholder

	cfg.Server.FrontendURL = "https://example.com/path"
	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() frontend_url with path valid error: %v", err)
placeholder

	cfg.Server.FrontendURL = "https://example.com?utm=1"
	if err := cfg.Validate(); err == nil {
		t.Fatalf("Validate() should reject server.frontend_url with query")
placeholder

	cfg.Server.FrontendURL = "https://user:pass@example.com"
	if err := cfg.Validate(); err == nil {
		t.Fatalf("Validate() should reject server.frontend_url with userinfo")
placeholder

	cfg.Server.FrontendURL = "/relative"
	if err := cfg.Validate(); err == nil {
		t.Fatalf("Validate() should reject relative server.frontend_url")
placeholder
placeholder

func TestValidateFrontendRedirectURL(t *testing.T) {
	if err := ValidateFrontendRedirectURL("/auth/callback"); err != nil {
		t.Fatalf("ValidateFrontendRedirectURL relative error: %v", err)
placeholder
	if err := ValidateFrontendRedirectURL("https://example.com/auth"); err != nil {
		t.Fatalf("ValidateFrontendRedirectURL absolute error: %v", err)
placeholder
	if err := ValidateFrontendRedirectURL("example.com/path"); err == nil {
		t.Fatalf("ValidateFrontendRedirectURL should reject non-absolute url")
placeholder
	if err := ValidateFrontendRedirectURL("//evil.com"); err == nil {
		t.Fatalf("ValidateFrontendRedirectURL should reject // prefix")
placeholder
	if err := ValidateFrontendRedirectURL("javascript:alert(1)"); err == nil {
		t.Fatalf("ValidateFrontendRedirectURL should reject javascript scheme")
placeholder
placeholder

func TestWarnIfInsecureURL(t *testing.T) {
	warnIfInsecureURL("test", "http://example.com")
	warnIfInsecureURL("test", "bad://url")
	warnIfInsecureURL("test", "://invalid")
placeholder

func TestGenerateJWTSecretDefaultLength(t *testing.T) {
	secret, err := generateJWTSecret(0)
	if err != nil {
		t.Fatalf("generateJWTSecret error: %v", err)
placeholder
	if len(secret) == 0 {
		t.Fatalf("generateJWTSecret returned empty string")
placeholder
placeholder

func TestValidateOpsCleanupScheduleRequired(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder
	cfg.Ops.Cleanup.Enabled = true
	cfg.Ops.Cleanup.Schedule = ""
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for ops.cleanup.schedule")
placeholder
	if !strings.Contains(err.Error(), "ops.cleanup.schedule") {
		t.Fatalf("Validate() expected ops.cleanup.schedule error, got: %v", err)
placeholder
placeholder

func TestValidateConcurrencyPingInterval(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder
	cfg.Concurrency.PingInterval = 3
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for concurrency.ping_interval")
placeholder
	if !strings.Contains(err.Error(), "concurrency.ping_interval") {
		t.Fatalf("Validate() expected concurrency.ping_interval error, got: %v", err)
placeholder
placeholder

func TestProvideConfig(t *testing.T) {
	resetViperWithJWTSecret(t)
	if _, err := ProvideConfig(); err != nil {
		t.Fatalf("ProvideConfig() error: %v", err)
placeholder
placeholder

func TestValidateConfigWithLinuxDoEnabled(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Security.CSP.Enabled = true
	cfg.Security.CSP.Policy = "default-src 'self'"

	cfg.LinuxDo.Enabled = true
	cfg.LinuxDo.ClientID = "client"
	cfg.LinuxDo.ClientSecret = "secret"
	cfg.LinuxDo.AuthorizeURL = "https://example.com/oauth2/authorize"
	cfg.LinuxDo.TokenURL = "https://example.com/oauth2/token"
	cfg.LinuxDo.UserInfoURL = "https://example.com/oauth2/userinfo"
	cfg.LinuxDo.RedirectURL = "https://example.com/api/v1/auth/oauth/linuxdo/callback"
	cfg.LinuxDo.FrontendRedirectURL = "/auth/linuxdo/callback"
	cfg.LinuxDo.TokenAuthMethod = "client_secret_post"
	cfg.LinuxDo.UsePKCE = true

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() unexpected error: %v", err)
placeholder
placeholder

func TestValidateJWTSecretStrength(t *testing.T) {
	if !isWeakJWTSecret("change-me-in-production") {
		t.Fatalf("isWeakJWTSecret should detect weak secret")
placeholder
	if isWeakJWTSecret("StrongSecretValue") {
		t.Fatalf("isWeakJWTSecret should accept strong secret")
placeholder
placeholder

func TestGenerateJWTSecretWithLength(t *testing.T) {
	secret, err := generateJWTSecret(16)
	if err != nil {
		t.Fatalf("generateJWTSecret error: %v", err)
placeholder
	if len(secret) == 0 {
		t.Fatalf("generateJWTSecret returned empty string")
placeholder
placeholder

func TestDatabaseDSNWithTimezone_WithPassword(t *testing.T) {
	d := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "u",
		Password: "p",
		DBName:   "db",
		SSLMode:  "prefer",
placeholder
	got := d.DSNWithTimezone("UTC")
	if !strings.Contains(got, "password=p") {
		t.Fatalf("DSNWithTimezone should include password: %q", got)
placeholder
	if !strings.Contains(got, "TimeZone=UTC") {
		t.Fatalf("DSNWithTimezone should include TimeZone=UTC: %q", got)
placeholder
placeholder

func TestValidateAbsoluteHTTPURLMissingHost(t *testing.T) {
	if err := ValidateAbsoluteHTTPURL("https://"); err == nil {
		t.Fatalf("ValidateAbsoluteHTTPURL should reject missing host")
placeholder
placeholder

func TestValidateFrontendRedirectURLInvalidChars(t *testing.T) {
	if err := ValidateFrontendRedirectURL("/auth/\ncallback"); err == nil {
		t.Fatalf("ValidateFrontendRedirectURL should reject invalid chars")
placeholder
	if err := ValidateFrontendRedirectURL("http://"); err == nil {
		t.Fatalf("ValidateFrontendRedirectURL should reject missing host")
placeholder
	if err := ValidateFrontendRedirectURL("mailto:user@example.com"); err == nil {
		t.Fatalf("ValidateFrontendRedirectURL should reject mailto")
placeholder
placeholder

func TestWarnIfInsecureURLHTTPS(t *testing.T) {
	warnIfInsecureURL("secure", "https://example.com")
placeholder

func TestValidateJWTSecret_UTF8Bytes(t *testing.T) {
	resetViperWithJWTSecret(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	// 31 bytes (< 32) even though it's 31 characters.
	cfg.JWT.Secret = strings.Repeat("a", 31)
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() should reject 31-byte secret")
placeholder
	if !strings.Contains(err.Error(), "at least 32 bytes") {
		t.Fatalf("Validate() error = %v", err)
placeholder

	// 32 bytes OK.
	cfg.JWT.Secret = strings.Repeat("a", 32)
	err = cfg.Validate()
	if err != nil {
		t.Fatalf("Validate() should accept 32-byte secret: %v", err)
placeholder
placeholder

func TestValidateConfigErrors(t *testing.T) {
	buildValid := func(t *testing.T) *Config {
	placeholder
		resetViperWithJWTSecret(t)
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error: %v", err)
	placeholder
		return cfg
placeholder

	cases := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
placeholder{
		{
			name:    "server read header timeout",
			mutate:  func(c *Config) { c.Server.ReadHeaderTimeout = 0 placeholder,
			wantErr: "server.read_header_timeout",
	placeholder,
		{
			name:    "server max header bytes too small",
			mutate:  func(c *Config) { c.Server.MaxHeaderBytes = 4096 placeholder,
			wantErr: "server.max_header_bytes",
	placeholder,
		{
			name:    "server max request body size",
			mutate:  func(c *Config) { c.Server.MaxRequestBodySize = -1 placeholder,
			wantErr: "server.max_request_body_size",
	placeholder,
		{
			name: "h2c zero concurrent streams",
			mutate: func(c *Config) {
				c.Server.H2C.Enabled = true
				c.Server.H2C.MaxConcurrentStreams = 0
		placeholder,
			wantErr: "server.h2c.max_concurrent_streams",
	placeholder,
		{
			name: "h2c oversized read frame",
			mutate: func(c *Config) {
				c.Server.H2C.Enabled = true
				c.Server.H2C.MaxReadFrameSize = 16 * 1024 * 1024
		placeholder,
			wantErr: "server.h2c.max_read_frame_size",
	placeholder,
		{
			name: "invalid auth abuse threshold too small",
			mutate: func(c *Config) {
				c.APIKeyAuth.InvalidAbuse.Threshold = 9
		placeholder,
			wantErr: "api_key_auth_cache.invalid_abuse.threshold",
	placeholder,
		{
			name: "invalid auth abuse capacity too small",
			mutate: func(c *Config) {
				c.APIKeyAuth.InvalidAbuse.Capacity = 255
		placeholder,
			wantErr: "api_key_auth_cache.invalid_abuse.capacity",
	placeholder,
		{
			name:    "jwt secret required",
			mutate:  func(c *Config) { c.JWT.Secret = "" placeholder,
			wantErr: "jwt.secret is required",
	placeholder,
		{
			name:    "jwt secret min bytes",
			mutate:  func(c *Config) { c.JWT.Secret = strings.Repeat("a", 31) placeholder,
			wantErr: "jwt.secret must be at least 32 bytes",
	placeholder,
		{
			name:    "subscription maintenance worker_count non-negative",
			mutate:  func(c *Config) { c.SubscriptionMaintenance.WorkerCount = -1 placeholder,
			wantErr: "subscription_maintenance.worker_count",
	placeholder,
		{
			name:    "subscription maintenance queue_size non-negative",
			mutate:  func(c *Config) { c.SubscriptionMaintenance.QueueSize = -1 placeholder,
			wantErr: "subscription_maintenance.queue_size",
	placeholder,
		{
			name:    "jwt expire hour positive",
			mutate:  func(c *Config) { c.JWT.ExpireHour = 0 placeholder,
			wantErr: "jwt.expire_hour must be positive",
	placeholder,
		{
			name:    "jwt expire hour max",
			mutate:  func(c *Config) { c.JWT.ExpireHour = 200 placeholder,
			wantErr: "jwt.expire_hour must be <= 168",
	placeholder,
		{
			name:    "jwt access token expire minutes non-negative",
			mutate:  func(c *Config) { c.JWT.AccessTokenExpireMinutes = -1 placeholder,
			wantErr: "jwt.access_token_expire_minutes must be non-negative",
	placeholder,
		{
			name:    "csp policy required",
			mutate:  func(c *Config) { c.Security.CSP.Enabled = true; c.Security.CSP.Policy = "" placeholder,
			wantErr: "security.csp.policy",
	placeholder,
		{
			name: "linuxdo client id required",
			mutate: func(c *Config) {
				c.LinuxDo.Enabled = true
				c.LinuxDo.UsePKCE = true
				c.LinuxDo.ClientID = ""
		placeholder,
			wantErr: "linuxdo_connect.client_id",
	placeholder,
		{
			name: "linuxdo token auth method",
			mutate: func(c *Config) {
				c.LinuxDo.Enabled = true
				c.LinuxDo.UsePKCE = true
				c.LinuxDo.ClientID = "client"
				c.LinuxDo.ClientSecret = "secret"
				c.LinuxDo.AuthorizeURL = "https://example.com/authorize"
				c.LinuxDo.TokenURL = "https://example.com/token"
				c.LinuxDo.UserInfoURL = "https://example.com/userinfo"
				c.LinuxDo.RedirectURL = "https://example.com/callback"
				c.LinuxDo.FrontendRedirectURL = "/auth/callback"
				c.LinuxDo.TokenAuthMethod = "invalid"
		placeholder,
			wantErr: "linuxdo_connect.token_auth_method",
	placeholder,
		{
			name:    "billing circuit breaker threshold",
			mutate:  func(c *Config) { c.Billing.CircuitBreaker.FailureThreshold = 0 placeholder,
			wantErr: "billing.circuit_breaker.failure_threshold",
	placeholder,
		{
			name:    "billing circuit breaker reset",
			mutate:  func(c *Config) { c.Billing.CircuitBreaker.ResetTimeoutSeconds = 0 placeholder,
			wantErr: "billing.circuit_breaker.reset_timeout_seconds",
	placeholder,
		{
			name:    "billing circuit breaker half open",
			mutate:  func(c *Config) { c.Billing.CircuitBreaker.HalfOpenRequests = 0 placeholder,
			wantErr: "billing.circuit_breaker.half_open_requests",
	placeholder,
		{
			name:    "billing minimum balance reserve",
			mutate:  func(c *Config) { c.Billing.MinimumBalanceReserve = -0.01 placeholder,
			wantErr: "billing.minimum_balance_reserve",
	placeholder,
		{
			name:    "database max open conns",
			mutate:  func(c *Config) { c.Database.MaxOpenConns = 0 placeholder,
			wantErr: "database.max_open_conns",
	placeholder,
		{
			name:    "database max lifetime",
			mutate:  func(c *Config) { c.Database.ConnMaxLifetimeMinutes = -1 placeholder,
			wantErr: "database.conn_max_lifetime_minutes",
	placeholder,
		{
			name:    "database idle exceeds open",
			mutate:  func(c *Config) { c.Database.MaxIdleConns = c.Database.MaxOpenConns + 1 placeholder,
			wantErr: "database.max_idle_conns cannot exceed",
	placeholder,
		{
			name:    "redis dial timeout",
			mutate:  func(c *Config) { c.Redis.DialTimeoutSeconds = 0 placeholder,
			wantErr: "redis.dial_timeout_seconds",
	placeholder,
		{
			name:    "redis read timeout",
			mutate:  func(c *Config) { c.Redis.ReadTimeoutSeconds = 0 placeholder,
			wantErr: "redis.read_timeout_seconds",
	placeholder,
		{
			name:    "redis write timeout",
			mutate:  func(c *Config) { c.Redis.WriteTimeoutSeconds = 0 placeholder,
			wantErr: "redis.write_timeout_seconds",
	placeholder,
		{
			name:    "redis pool size",
			mutate:  func(c *Config) { c.Redis.PoolSize = 0 placeholder,
			wantErr: "redis.pool_size",
	placeholder,
		{
			name:    "redis idle exceeds pool",
			mutate:  func(c *Config) { c.Redis.MinIdleConns = c.Redis.PoolSize + 1 placeholder,
			wantErr: "redis.min_idle_conns cannot exceed",
	placeholder,
		{
			name:    "dashboard cache disabled negative",
			mutate:  func(c *Config) { c.Dashboard.Enabled = false; c.Dashboard.StatsTTLSeconds = -1 placeholder,
			wantErr: "dashboard_cache.stats_ttl_seconds",
	placeholder,
		{
			name:    "dashboard cache fresh ttl positive",
			mutate:  func(c *Config) { c.Dashboard.Enabled = true; c.Dashboard.StatsFreshTTLSeconds = 0 placeholder,
			wantErr: "dashboard_cache.stats_fresh_ttl_seconds",
	placeholder,
		{
			name:    "dashboard aggregation enabled interval",
			mutate:  func(c *Config) { c.DashboardAgg.Enabled = true; c.DashboardAgg.IntervalSeconds = 0 placeholder,
			wantErr: "dashboard_aggregation.interval_seconds",
	placeholder,
		{
			name: "dashboard aggregation backfill positive",
			mutate: func(c *Config) {
				c.DashboardAgg.Enabled = true
				c.DashboardAgg.BackfillEnabled = true
				c.DashboardAgg.BackfillMaxDays = 0
		placeholder,
			wantErr: "dashboard_aggregation.backfill_max_days",
	placeholder,
		{
			name:    "dashboard aggregation retention",
			mutate:  func(c *Config) { c.DashboardAgg.Enabled = true; c.DashboardAgg.Retention.UsageLogsDays = 0 placeholder,
			wantErr: "dashboard_aggregation.retention.usage_logs_days",
	placeholder,
		{
			name: "dashboard aggregation dedup retention",
			mutate: func(c *Config) {
				c.DashboardAgg.Enabled = true
				c.DashboardAgg.Retention.UsageBillingDedupDays = 0
		placeholder,
			wantErr: "dashboard_aggregation.retention.usage_billing_dedup_days",
	placeholder,
		{
			name: "dashboard aggregation dedup retention smaller than usage logs",
			mutate: func(c *Config) {
				c.DashboardAgg.Enabled = true
				c.DashboardAgg.Retention.UsageLogsDays = 30
				c.DashboardAgg.Retention.UsageBillingDedupDays = 29
		placeholder,
			wantErr: "dashboard_aggregation.retention.usage_billing_dedup_days",
	placeholder,
		{
			name:    "dashboard aggregation disabled interval",
			mutate:  func(c *Config) { c.DashboardAgg.Enabled = false; c.DashboardAgg.IntervalSeconds = -1 placeholder,
			wantErr: "dashboard_aggregation.interval_seconds",
	placeholder,
		{
			name:    "usage cleanup max range",
			mutate:  func(c *Config) { c.UsageCleanup.Enabled = true; c.UsageCleanup.MaxRangeDays = 0 placeholder,
			wantErr: "usage_cleanup.max_range_days",
	placeholder,
		{
			name:    "usage cleanup worker interval",
			mutate:  func(c *Config) { c.UsageCleanup.Enabled = true; c.UsageCleanup.WorkerIntervalSeconds = 0 placeholder,
			wantErr: "usage_cleanup.worker_interval_seconds",
	placeholder,
		{
			name:    "usage cleanup batch size",
			mutate:  func(c *Config) { c.UsageCleanup.Enabled = true; c.UsageCleanup.BatchSize = 0 placeholder,
			wantErr: "usage_cleanup.batch_size",
	placeholder,
		{
			name:    "usage cleanup disabled negative",
			mutate:  func(c *Config) { c.UsageCleanup.Enabled = false; c.UsageCleanup.BatchSize = -1 placeholder,
			wantErr: "usage_cleanup.batch_size",
	placeholder,
		{
			name:    "gateway max body size",
			mutate:  func(c *Config) { c.Gateway.MaxBodySize = 0 placeholder,
			wantErr: "gateway.max_body_size",
	placeholder,
		{
			name:    "gateway text body exceeds media body",
			mutate:  func(c *Config) { c.Gateway.TextMaxBodySize = c.Gateway.MaxBodySize + 1 placeholder,
			wantErr: "gateway.text_max_body_size",
	placeholder,
		{
			name:    "gateway response header timeout",
			mutate:  func(c *Config) { c.Gateway.ResponseHeaderTimeout = -1 placeholder,
			wantErr: "gateway.response_header_timeout",
	placeholder,
		{
			name:    "gateway openai response header timeout",
			mutate:  func(c *Config) { c.Gateway.OpenAIResponseHeaderTimeout = -1 placeholder,
			wantErr: "gateway.openai_response_header_timeout",
	placeholder,
		{
			name:    "gateway openai first output timeout below minimum",
			mutate:  func(c *Config) { c.Gateway.OpenAIFirstOutputTimeoutSeconds = 29 placeholder,
			wantErr: "gateway.openai_first_output_timeout_seconds",
	placeholder,
		{
			name:    "gateway openai high effort first output timeout too large",
			mutate:  func(c *Config) { c.Gateway.OpenAIHighEffortFirstOutputTimeoutSeconds = 1801 placeholder,
			wantErr: "gateway.openai_high_effort_first_output_timeout_seconds",
	placeholder,
		{
			name:    "gateway max idle conns",
			mutate:  func(c *Config) { c.Gateway.MaxIdleConns = 0 placeholder,
			wantErr: "gateway.max_idle_conns",
	placeholder,
		{
			name:    "gateway max idle conns per host",
			mutate:  func(c *Config) { c.Gateway.MaxIdleConnsPerHost = 0 placeholder,
			wantErr: "gateway.max_idle_conns_per_host",
	placeholder,
		{
			name:    "gateway idle timeout",
			mutate:  func(c *Config) { c.Gateway.IdleConnTimeoutSeconds = 0 placeholder,
			wantErr: "gateway.idle_conn_timeout_seconds",
	placeholder,
		{
			name:    "gateway max upstream clients",
			mutate:  func(c *Config) { c.Gateway.MaxUpstreamClients = 0 placeholder,
			wantErr: "gateway.max_upstream_clients",
	placeholder,
		{
			name:    "gateway client idle ttl",
			mutate:  func(c *Config) { c.Gateway.ClientIdleTTLSeconds = 0 placeholder,
			wantErr: "gateway.client_idle_ttl_seconds",
	placeholder,
		{
			name:    "gateway concurrency slot ttl",
			mutate:  func(c *Config) { c.Gateway.ConcurrencySlotTTLMinutes = 0 placeholder,
			wantErr: "gateway.concurrency_slot_ttl_minutes",
	placeholder,
		{
			name:    "gateway max conns per host",
			mutate:  func(c *Config) { c.Gateway.MaxConnsPerHost = -1 placeholder,
			wantErr: "gateway.max_conns_per_host",
	placeholder,
		{
			name:    "gateway connection isolation",
			mutate:  func(c *Config) { c.Gateway.ConnectionPoolIsolation = "invalid" placeholder,
			wantErr: "gateway.connection_pool_isolation",
	placeholder,
		{
			name:    "gateway stream keepalive range",
			mutate:  func(c *Config) { c.Gateway.StreamKeepaliveInterval = 4 placeholder,
			wantErr: "gateway.stream_keepalive_interval",
	placeholder,
		{
			name:    "gateway openai ws oauth max conns factor",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.OAuthMaxConnsFactor = 0 placeholder,
			wantErr: "gateway.openai_ws.oauth_max_conns_factor",
	placeholder,
		{
			name:    "gateway openai ws apikey max conns factor",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.APIKeyMaxConnsFactor = 0 placeholder,
			wantErr: "gateway.openai_ws.apikey_max_conns_factor",
	placeholder,
		{
			name:    "gateway openai http2 fallback threshold",
			mutate:  func(c *Config) { c.Gateway.OpenAIHTTP2.FallbackErrorThreshold = -1 placeholder,
			wantErr: "gateway.openai_http2.fallback_error_threshold",
	placeholder,
		{
			name:    "gateway openai http2 fallback window",
			mutate:  func(c *Config) { c.Gateway.OpenAIHTTP2.FallbackWindowSeconds = -1 placeholder,
			wantErr: "gateway.openai_http2.fallback_window_seconds",
	placeholder,
		{
			name:    "gateway openai http2 fallback ttl",
			mutate:  func(c *Config) { c.Gateway.OpenAIHTTP2.FallbackTTLSeconds = -1 placeholder,
			wantErr: "gateway.openai_http2.fallback_ttl_seconds",
	placeholder,
		{
			name:    "gateway stream data interval range",
			mutate:  func(c *Config) { c.Gateway.StreamDataIntervalTimeout = 5 placeholder,
			wantErr: "gateway.stream_data_interval_timeout",
	placeholder,
		{
			name:    "gateway stream data interval negative",
			mutate:  func(c *Config) { c.Gateway.StreamDataIntervalTimeout = -1 placeholder,
			wantErr: "gateway.stream_data_interval_timeout must be non-negative",
	placeholder,
		{
			name:    "gateway image stream keepalive range",
			mutate:  func(c *Config) { c.Gateway.ImageStreamKeepaliveInterval = 4 placeholder,
			wantErr: "gateway.image_stream_keepalive_interval",
	placeholder,
		{
			name:    "gateway image stream keepalive negative",
			mutate:  func(c *Config) { c.Gateway.ImageStreamKeepaliveInterval = -1 placeholder,
			wantErr: "gateway.image_stream_keepalive_interval must be non-negative",
	placeholder,
		{
			name:    "gateway image nonstream keepalive range",
			mutate:  func(c *Config) { c.Gateway.ImageNonstreamKeepaliveInterval = 4 placeholder,
			wantErr: "gateway.image_nonstream_keepalive_interval",
	placeholder,
		{
			name:    "gateway image nonstream keepalive negative",
			mutate:  func(c *Config) { c.Gateway.ImageNonstreamKeepaliveInterval = -1 placeholder,
			wantErr: "gateway.image_nonstream_keepalive_interval must be non-negative",
	placeholder,
		{
			name:    "gateway image stream data interval range",
			mutate:  func(c *Config) { c.Gateway.ImageStreamDataIntervalTimeout = 30 placeholder,
			wantErr: "gateway.image_stream_data_interval_timeout",
	placeholder,
		{
			name:    "gateway image stream data interval negative",
			mutate:  func(c *Config) { c.Gateway.ImageStreamDataIntervalTimeout = -1 placeholder,
			wantErr: "gateway.image_stream_data_interval_timeout must be non-negative",
	placeholder,
		{
			name:    "gateway image concurrency max negative",
			mutate:  func(c *Config) { c.Gateway.ImageConcurrency.MaxConcurrentRequests = -1 placeholder,
			wantErr: "gateway.image_concurrency.max_concurrent_requests must be non-negative",
	placeholder,
		{
			name:    "gateway image concurrency overflow mode invalid",
			mutate:  func(c *Config) { c.Gateway.ImageConcurrency.OverflowMode = "queue" placeholder,
			wantErr: "gateway.image_concurrency.overflow_mode",
	placeholder,
		{
			name:    "gateway image concurrency wait timeout negative",
			mutate:  func(c *Config) { c.Gateway.ImageConcurrency.WaitTimeoutSeconds = -1 placeholder,
			wantErr: "gateway.image_concurrency.wait_timeout_seconds must be non-negative",
	placeholder,
		{
			name:    "gateway image concurrency max waiting negative",
			mutate:  func(c *Config) { c.Gateway.ImageConcurrency.MaxWaitingRequests = -1 placeholder,
			wantErr: "gateway.image_concurrency.max_waiting_requests must be non-negative",
	placeholder,
		{
			name:    "gateway max line size",
			mutate:  func(c *Config) { c.Gateway.MaxLineSize = 1024 placeholder,
			wantErr: "gateway.max_line_size must be at least",
	placeholder,
		{
			name:    "gateway max line size negative",
			mutate:  func(c *Config) { c.Gateway.MaxLineSize = -1 placeholder,
			wantErr: "gateway.max_line_size must be non-negative",
	placeholder,
		{
			name:    "gateway usage record worker count",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.WorkerCount = 0 placeholder,
			wantErr: "gateway.usage_record.worker_count",
	placeholder,
		{
			name:    "gateway usage record queue size",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.QueueSize = 0 placeholder,
			wantErr: "gateway.usage_record.queue_size",
	placeholder,
		{
			name:    "gateway usage record timeout",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.TaskTimeoutSeconds = 0 placeholder,
			wantErr: "gateway.usage_record.task_timeout_seconds",
	placeholder,
		{
			name:    "gateway usage record overflow policy",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.OverflowPolicy = "invalid" placeholder,
			wantErr: "gateway.usage_record.overflow_policy",
	placeholder,
		{
			name:    "gateway usage record sample percent range",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.OverflowSamplePercent = 101 placeholder,
			wantErr: "gateway.usage_record.overflow_sample_percent",
	placeholder,
		{
			name: "gateway usage record sample percent required for sample policy",
			mutate: func(c *Config) {
				c.Gateway.UsageRecord.OverflowPolicy = UsageRecordOverflowPolicySample
				c.Gateway.UsageRecord.OverflowSamplePercent = 0
		placeholder,
			wantErr: "gateway.usage_record.overflow_sample_percent must be positive",
	placeholder,
		{
			name: "gateway usage record auto scale max gte min",
			mutate: func(c *Config) {
				c.Gateway.UsageRecord.AutoScaleMinWorkers = 256
				c.Gateway.UsageRecord.AutoScaleMaxWorkers = 128
		placeholder,
			wantErr: "gateway.usage_record.auto_scale_max_workers",
	placeholder,
		{
			name: "gateway usage record worker in auto scale range",
			mutate: func(c *Config) {
				c.Gateway.UsageRecord.AutoScaleMinWorkers = 200
				c.Gateway.UsageRecord.AutoScaleMaxWorkers = 300
				c.Gateway.UsageRecord.WorkerCount = 128
		placeholder,
			wantErr: "gateway.usage_record.worker_count must be between auto_scale_min_workers and auto_scale_max_workers",
	placeholder,
		{
			name: "gateway usage record auto scale queue thresholds order",
			mutate: func(c *Config) {
				c.Gateway.UsageRecord.AutoScaleUpQueuePercent = 50
				c.Gateway.UsageRecord.AutoScaleDownQueuePercent = 50
		placeholder,
			wantErr: "gateway.usage_record.auto_scale_down_queue_percent must be less",
	placeholder,
		{
			name:    "gateway usage record auto scale up step",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.AutoScaleUpStep = 0 placeholder,
			wantErr: "gateway.usage_record.auto_scale_up_step",
	placeholder,
		{
			name:    "gateway usage record auto scale interval",
			mutate:  func(c *Config) { c.Gateway.UsageRecord.AutoScaleCheckIntervalSeconds = 0 placeholder,
			wantErr: "gateway.usage_record.auto_scale_check_interval_seconds",
	placeholder,
		{
			name:    "gateway user group rate cache ttl",
			mutate:  func(c *Config) { c.Gateway.UserGroupRateCacheTTLSeconds = 0 placeholder,
			wantErr: "gateway.user_group_rate_cache_ttl_seconds",
	placeholder,
		{
			name:    "gateway models list cache ttl range",
			mutate:  func(c *Config) { c.Gateway.ModelsListCacheTTLSeconds = 31 placeholder,
			wantErr: "gateway.models_list_cache_ttl_seconds",
	placeholder,
		{
			name:    "gateway scheduling sticky waiting",
			mutate:  func(c *Config) { c.Gateway.Scheduling.StickySessionMaxWaiting = 0 placeholder,
			wantErr: "gateway.scheduling.sticky_session_max_waiting",
	placeholder,
		{
			name:    "gateway scheduling load batch cache ttl",
			mutate:  func(c *Config) { c.Gateway.Scheduling.LoadBatchCacheTTLMS = -1 placeholder,
			wantErr: "gateway.scheduling.load_batch_cache_ttl_ms",
	placeholder,
		{
			name:    "gateway scheduling outbox poll",
			mutate:  func(c *Config) { c.Gateway.Scheduling.OutboxPollIntervalSeconds = 0 placeholder,
			wantErr: "gateway.scheduling.outbox_poll_interval_seconds",
	placeholder,
		{
			name:    "gateway scheduling outbox failures",
			mutate:  func(c *Config) { c.Gateway.Scheduling.OutboxLagRebuildFailures = 0 placeholder,
			wantErr: "gateway.scheduling.outbox_lag_rebuild_failures",
	placeholder,
		{
			name: "gateway outbox lag rebuild",
			mutate: func(c *Config) {
				c.Gateway.Scheduling.OutboxLagWarnSeconds = 10
				c.Gateway.Scheduling.OutboxLagRebuildSeconds = 5
		placeholder,
			wantErr: "gateway.scheduling.outbox_lag_rebuild_seconds",
	placeholder,
		{
			name:    "log level invalid",
			mutate:  func(c *Config) { c.Log.Level = "trace" placeholder,
			wantErr: "log.level",
	placeholder,
		{
			name:    "log format invalid",
			mutate:  func(c *Config) { c.Log.Format = "plain" placeholder,
			wantErr: "log.format",
	placeholder,
		{
			name: "log output disabled",
			mutate: func(c *Config) {
				c.Log.Output.ToStdout = false
				c.Log.Output.ToFile = false
		placeholder,
			wantErr: "log.output.to_stdout and log.output.to_file cannot both be false",
	placeholder,
		{
			name:    "log rotation size",
			mutate:  func(c *Config) { c.Log.Rotation.MaxSizeMB = 0 placeholder,
			wantErr: "log.rotation.max_size_mb",
	placeholder,
		{
			name: "log sampling enabled invalid",
			mutate: func(c *Config) {
				c.Log.Sampling.Enabled = true
				c.Log.Sampling.Initial = 0
		placeholder,
			wantErr: "log.sampling.initial",
	placeholder,
		{
			name:    "ops metrics collector ttl",
			mutate:  func(c *Config) { c.Ops.MetricsCollectorCache.TTL = -1 placeholder,
			wantErr: "ops.metrics_collector_cache.ttl",
	placeholder,
		{
			name:    "ops cleanup retention",
			mutate:  func(c *Config) { c.Ops.Cleanup.ErrorLogRetentionDays = -1 placeholder,
			wantErr: "ops.cleanup.error_log_retention_days",
	placeholder,
		{
			name:    "ops cleanup minute retention",
			mutate:  func(c *Config) { c.Ops.Cleanup.MinuteMetricsRetentionDays = -1 placeholder,
			wantErr: "ops.cleanup.minute_metrics_retention_days",
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := buildValid(t)
			tt.mutate(cfg)
			err := cfg.Validate()
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %q", err, tt.wantErr)
		placeholder
	placeholder)
placeholder
placeholder

func TestValidateConfig_OpenAIWSRules(t *testing.T) {
	buildValid := func(t *testing.T) *Config {
	placeholder
		resetViperWithJWTSecret(t)
		cfg, err := Load()
	placeholder
		return cfg
placeholder

	t.Run("sticky response id ttl 兼容旧键回填", func(t *testing.T) {
		cfg := buildValid(t)
		cfg.Gateway.OpenAIWS.StickyResponseIDTTLSeconds = 0
		cfg.Gateway.OpenAIWS.StickyPreviousResponseTTLSeconds = 7200

		require.NoError(t, cfg.Validate())
		require.Equal(t, 7200, cfg.Gateway.OpenAIWS.StickyResponseIDTTLSeconds)
placeholder)

	cases := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
placeholder{
		{
			name:    "max_conns_per_account 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.MaxConnsPerAccount = 0 placeholder,
			wantErr: "gateway.openai_ws.max_conns_per_account",
	placeholder,
		{
			name:    "client_first_message_timeout_seconds 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds = 0 placeholder,
			wantErr: "gateway.openai_ws.client_first_message_timeout_seconds",
	placeholder,
		{
			name:    "client_first_message_timeout_seconds 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.ClientFirstMessageTimeoutSeconds = -1 placeholder,
			wantErr: "gateway.openai_ws.client_first_message_timeout_seconds",
	placeholder,
		{
			name:    "ingress_inter_turn_idle_timeout_seconds 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.IngressInterTurnIdleTimeoutSeconds = -1 placeholder,
			wantErr: "gateway.openai_ws.ingress_inter_turn_idle_timeout_seconds",
	placeholder,
		{
			name:    "max_ingress_connections_per_api_key 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.MaxIngressConnectionsPerAPIKey = -1 placeholder,
			wantErr: "gateway.openai_ws.max_ingress_connections_per_api_key",
	placeholder,
		{
			name:    "min_idle_per_account 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.MinIdlePerAccount = -1 placeholder,
			wantErr: "gateway.openai_ws.min_idle_per_account",
	placeholder,
		{
			name:    "max_idle_per_account 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.MaxIdlePerAccount = -1 placeholder,
			wantErr: "gateway.openai_ws.max_idle_per_account",
	placeholder,
		{
			name: "min_idle_per_account 不能大于 max_idle_per_account",
			mutate: func(c *Config) {
				c.Gateway.OpenAIWS.MinIdlePerAccount = 3
				c.Gateway.OpenAIWS.MaxIdlePerAccount = 2
		placeholder,
			wantErr: "gateway.openai_ws.min_idle_per_account must be <= max_idle_per_account",
	placeholder,
		{
			name: "max_idle_per_account 不能大于 max_conns_per_account",
			mutate: func(c *Config) {
				c.Gateway.OpenAIWS.MaxConnsPerAccount = 2
				c.Gateway.OpenAIWS.MinIdlePerAccount = 1
				c.Gateway.OpenAIWS.MaxIdlePerAccount = 3
		placeholder,
			wantErr: "gateway.openai_ws.max_idle_per_account must be <= max_conns_per_account",
	placeholder,
		{
			name:    "dial_timeout_seconds 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.DialTimeoutSeconds = 0 placeholder,
			wantErr: "gateway.openai_ws.dial_timeout_seconds",
	placeholder,
		{
			name:    "read_timeout_seconds 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.ReadTimeoutSeconds = 0 placeholder,
			wantErr: "gateway.openai_ws.read_timeout_seconds",
	placeholder,
		{
			name:    "write_timeout_seconds 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.WriteTimeoutSeconds = 0 placeholder,
			wantErr: "gateway.openai_ws.write_timeout_seconds",
	placeholder,
		{
			name:    "pool_target_utilization 必须在 (0,1]",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.PoolTargetUtilization = 0 placeholder,
			wantErr: "gateway.openai_ws.pool_target_utilization",
	placeholder,
		{
			name:    "queue_limit_per_conn 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.QueueLimitPerConn = 0 placeholder,
			wantErr: "gateway.openai_ws.queue_limit_per_conn",
	placeholder,
		{
			name:    "fallback_cooldown_seconds 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.FallbackCooldownSeconds = -1 placeholder,
			wantErr: "gateway.openai_ws.fallback_cooldown_seconds",
	placeholder,
		{
			name:    "store_disabled_conn_mode 必须为 strict|adaptive|off",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.StoreDisabledConnMode = "invalid" placeholder,
			wantErr: "gateway.openai_ws.store_disabled_conn_mode",
	placeholder,
		{
			name:    "ingress_mode_default 必须为 off|ctx_pool|passthrough|http_bridge",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.IngressModeDefault = "invalid" placeholder,
			wantErr: "gateway.openai_ws.ingress_mode_default",
	placeholder,
		{
			name:    "payload_log_sample_rate 必须在 [0,1] 范围内",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.PayloadLogSampleRate = 1.2 placeholder,
			wantErr: "gateway.openai_ws.payload_log_sample_rate",
	placeholder,
		{
			name:    "retry_total_budget_ms 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.RetryTotalBudgetMS = -1 placeholder,
			wantErr: "gateway.openai_ws.retry_total_budget_ms",
	placeholder,
		{
			name:    "lb_top_k 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.LBTopK = 0 placeholder,
			wantErr: "gateway.openai_ws.lb_top_k",
	placeholder,
		{
			name:    "sticky_session_ttl_seconds 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.StickySessionTTLSeconds = 0 placeholder,
			wantErr: "gateway.openai_ws.sticky_session_ttl_seconds",
	placeholder,
		{
			name: "sticky_response_id_ttl_seconds 必须为正数",
			mutate: func(c *Config) {
				c.Gateway.OpenAIWS.StickyResponseIDTTLSeconds = 0
				c.Gateway.OpenAIWS.StickyPreviousResponseTTLSeconds = 0
		placeholder,
			wantErr: "gateway.openai_ws.sticky_response_id_ttl_seconds",
	placeholder,
		{
			name:    "sticky_previous_response_ttl_seconds 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.StickyPreviousResponseTTLSeconds = -1 placeholder,
			wantErr: "gateway.openai_ws.sticky_previous_response_ttl_seconds",
	placeholder,
		{
			name:    "scheduler_score_weights 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.SchedulerScoreWeights.Queue = -0.1 placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights.* must be non-negative",
	placeholder,
		{
			name:    "scheduler_score_weights quota_headroom 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.SchedulerScoreWeights.QuotaHeadroom = -0.1 placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights.* must be non-negative",
	placeholder,
		{
			name:    "scheduler_score_weights upstream_cost 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = -0.1 placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights.* must be non-negative",
	placeholder,
		{
			name:    "scheduler_score_weights reset 不能为负数",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.SchedulerScoreWeights.Reset = -0.1 placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights.* must be non-negative",
	placeholder,
		{
			name:    "scheduler_score_weights 不能为 NaN",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.SchedulerScoreWeights.PreviousResponse = math.NaN() placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights.* must be non-negative and finite",
	placeholder,
		{
			name:    "scheduler_score_weights 不能为 Inf",
			mutate:  func(c *Config) { c.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = math.Inf(1) placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights.* must be non-negative and finite",
	placeholder,
		{
			name: "scheduler_score_weights 总和不能溢出",
			mutate: func(c *Config) {
				c.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = math.MaxFloat64
				c.Gateway.OpenAIWS.SchedulerScoreWeights.Load = math.MaxFloat64
		placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights base-weight sum must be finite",
	placeholder,
		{
			name: "scheduler_score_weights 含 sticky 总和不能溢出",
			mutate: func(c *Config) {
				c.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = math.MaxFloat64
				c.Gateway.OpenAIWS.SchedulerScoreWeights.PreviousResponse = math.MaxFloat64
		placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights total-weight sum must be finite",
	placeholder,
		{
			name: "scheduler_score_weights 不能全为 0",
			mutate: func(c *Config) {
				c.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = 0
				c.Gateway.OpenAIWS.SchedulerScoreWeights.Load = 0
				c.Gateway.OpenAIWS.SchedulerScoreWeights.Queue = 0
				c.Gateway.OpenAIWS.SchedulerScoreWeights.ErrorRate = 0
				c.Gateway.OpenAIWS.SchedulerScoreWeights.TTFT = 0
		placeholder,
			wantErr: "gateway.openai_ws.scheduler_score_weights must not all be zero",
	placeholder,
		{
			name:    "sticky_escape_ttft_ms 必须为正数",
			mutate:  func(c *Config) { c.Gateway.OpenAIScheduler.StickyEscapeTTFTMs = 0 placeholder,
			wantErr: "gateway.openai_scheduler.sticky_escape_ttft_ms",
	placeholder,
		{
			name:    "sticky_escape_error_rate 不能小于 0",
			mutate:  func(c *Config) { c.Gateway.OpenAIScheduler.StickyEscapeErrorRate = -0.1 placeholder,
			wantErr: "gateway.openai_scheduler.sticky_escape_error_rate",
	placeholder,
		{
			name:    "sticky_escape_error_rate 不能大于 1",
			mutate:  func(c *Config) { c.Gateway.OpenAIScheduler.StickyEscapeErrorRate = 1.1 placeholder,
			wantErr: "gateway.openai_scheduler.sticky_escape_error_rate",
	placeholder,
placeholder

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cfg := buildValid(t)
			tc.mutate(cfg)

			err := cfg.Validate()
		placeholder
			require.Contains(t, err.Error(), tc.wantErr)
	placeholder)
placeholder

	t.Run("quota_headroom 可作为唯一有效调度权重", func(t *testing.T) {
		cfg := buildValid(t)
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Load = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Queue = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.ErrorRate = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.TTFT = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.QuotaHeadroom = 0.1

		require.NoError(t, cfg.Validate())
placeholder)

	t.Run("upstream_cost 可作为唯一有效调度权重", func(t *testing.T) {
		cfg := buildValid(t)
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Load = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Queue = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.ErrorRate = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.TTFT = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost = 0.1

		require.NoError(t, cfg.Validate())
placeholder)

	t.Run("reset 可作为唯一有效调度权重", func(t *testing.T) {
		cfg := buildValid(t)
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Priority = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Load = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Queue = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.ErrorRate = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.TTFT = 0
		cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Reset = 0.1

		require.NoError(t, cfg.Validate())
placeholder)
placeholder

func TestValidateConfig_AutoScaleDisabledIgnoreAutoScaleFields(t *testing.T) {
	resetViperWithJWTSecret(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Gateway.UsageRecord.AutoScaleEnabled = false
	cfg.Gateway.UsageRecord.WorkerCount = 64

	// 自动扩缩容关闭时，这些字段应被忽略，不应导致校验失败。
	cfg.Gateway.UsageRecord.AutoScaleMinWorkers = 0
	cfg.Gateway.UsageRecord.AutoScaleMaxWorkers = 0
	cfg.Gateway.UsageRecord.AutoScaleUpQueuePercent = 0
	cfg.Gateway.UsageRecord.AutoScaleDownQueuePercent = 100
	cfg.Gateway.UsageRecord.AutoScaleUpStep = 0
	cfg.Gateway.UsageRecord.AutoScaleDownStep = 0
	cfg.Gateway.UsageRecord.AutoScaleCheckIntervalSeconds = 0
	cfg.Gateway.UsageRecord.AutoScaleCooldownSeconds = -1

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate() should ignore auto scale fields when disabled: %v", err)
placeholder
placeholder

func TestValidateConfig_LogRequiredAndRotationBounds(t *testing.T) {
	resetViperWithJWTSecret(t)

	cases := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
placeholder{
		{
			name: "log level required",
			mutate: func(c *Config) {
				c.Log.Level = ""
		placeholder,
			wantErr: "log.level is required",
	placeholder,
		{
			name: "log format required",
			mutate: func(c *Config) {
				c.Log.Format = ""
		placeholder,
			wantErr: "log.format is required",
	placeholder,
		{
			name: "log stacktrace required",
			mutate: func(c *Config) {
				c.Log.StacktraceLevel = ""
		placeholder,
			wantErr: "log.stacktrace_level is required",
	placeholder,
		{
			name: "log max backups non-negative",
			mutate: func(c *Config) {
				c.Log.Rotation.MaxBackups = -1
		placeholder,
			wantErr: "log.rotation.max_backups must be non-negative",
	placeholder,
		{
			name: "log max age non-negative",
			mutate: func(c *Config) {
				c.Log.Rotation.MaxAgeDays = -1
		placeholder,
			wantErr: "log.rotation.max_age_days must be non-negative",
	placeholder,
		{
			name: "sampling thereafter non-negative when disabled",
			mutate: func(c *Config) {
				c.Log.Sampling.Enabled = false
				c.Log.Sampling.Thereafter = -1
		placeholder,
			wantErr: "log.sampling.thereafter must be non-negative",
	placeholder,
placeholder

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() error: %v", err)
		placeholder
			tt.mutate(cfg)
			err = cfg.Validate()
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %q", err, tt.wantErr)
		placeholder
	placeholder)
placeholder
placeholder

func TestLoad_DefaultGatewayUsageRecordConfig(t *testing.T) {
	resetViperWithJWTSecret(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder
	if cfg.Gateway.UsageRecord.WorkerCount != 128 {
		t.Fatalf("worker_count = %d, want 128", cfg.Gateway.UsageRecord.WorkerCount)
placeholder
	if cfg.Gateway.UsageRecord.QueueSize != 16384 {
		t.Fatalf("queue_size = %d, want 16384", cfg.Gateway.UsageRecord.QueueSize)
placeholder
	if cfg.Gateway.UsageRecord.TaskTimeoutSeconds != 5 {
		t.Fatalf("task_timeout_seconds = %d, want 5", cfg.Gateway.UsageRecord.TaskTimeoutSeconds)
placeholder
	if cfg.Gateway.UsageRecord.OverflowPolicy != UsageRecordOverflowPolicySync {
		t.Fatalf("overflow_policy = %s, want %s", cfg.Gateway.UsageRecord.OverflowPolicy, UsageRecordOverflowPolicySync)
placeholder
	if cfg.Gateway.UsageRecord.OverflowSamplePercent != 10 {
		t.Fatalf("overflow_sample_percent = %d, want 10", cfg.Gateway.UsageRecord.OverflowSamplePercent)
placeholder
	if !cfg.Gateway.UsageRecord.AutoScaleEnabled {
		t.Fatalf("auto_scale_enabled = false, want true")
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleMinWorkers != 128 {
		t.Fatalf("auto_scale_min_workers = %d, want 128", cfg.Gateway.UsageRecord.AutoScaleMinWorkers)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleMaxWorkers != 512 {
		t.Fatalf("auto_scale_max_workers = %d, want 512", cfg.Gateway.UsageRecord.AutoScaleMaxWorkers)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleUpQueuePercent != 70 {
		t.Fatalf("auto_scale_up_queue_percent = %d, want 70", cfg.Gateway.UsageRecord.AutoScaleUpQueuePercent)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleDownQueuePercent != 15 {
		t.Fatalf("auto_scale_down_queue_percent = %d, want 15", cfg.Gateway.UsageRecord.AutoScaleDownQueuePercent)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleUpStep != 32 {
		t.Fatalf("auto_scale_up_step = %d, want 32", cfg.Gateway.UsageRecord.AutoScaleUpStep)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleDownStep != 16 {
		t.Fatalf("auto_scale_down_step = %d, want 16", cfg.Gateway.UsageRecord.AutoScaleDownStep)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleCheckIntervalSeconds != 3 {
		t.Fatalf("auto_scale_check_interval_seconds = %d, want 3", cfg.Gateway.UsageRecord.AutoScaleCheckIntervalSeconds)
placeholder
	if cfg.Gateway.UsageRecord.AutoScaleCooldownSeconds != 10 {
		t.Fatalf("auto_scale_cooldown_seconds = %d, want 10", cfg.Gateway.UsageRecord.AutoScaleCooldownSeconds)
placeholder
placeholder

func TestLoad_DefaultGatewayImageStreamConfig(t *testing.T) {
	resetViperWithJWTSecret(t)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder
	if cfg.Gateway.StreamDataIntervalTimeout != 180 {
		t.Fatalf("stream_data_interval_timeout = %d, want 180", cfg.Gateway.StreamDataIntervalTimeout)
placeholder
	if cfg.Gateway.StreamKeepaliveInterval != 10 {
		t.Fatalf("stream_keepalive_interval = %d, want 10", cfg.Gateway.StreamKeepaliveInterval)
placeholder
	if cfg.Gateway.ImageStreamDataIntervalTimeout != 900 {
		t.Fatalf("image_stream_data_interval_timeout = %d, want 900", cfg.Gateway.ImageStreamDataIntervalTimeout)
placeholder
	if cfg.Gateway.ImageStreamKeepaliveInterval != 10 {
		t.Fatalf("image_stream_keepalive_interval = %d, want 10", cfg.Gateway.ImageStreamKeepaliveInterval)
placeholder
	if cfg.Gateway.ImageNonstreamKeepaliveInterval != 0 {
		t.Fatalf("image_nonstream_keepalive_interval = %d, want 0", cfg.Gateway.ImageNonstreamKeepaliveInterval)
placeholder
	if cfg.Gateway.ImageConcurrency.Enabled {
		t.Fatalf("image_concurrency.enabled = true, want false")
placeholder
	if cfg.Gateway.ImageConcurrency.MaxConcurrentRequests != 0 {
		t.Fatalf("image_concurrency.max_concurrent_requests = %d, want 0", cfg.Gateway.ImageConcurrency.MaxConcurrentRequests)
placeholder
	if cfg.Gateway.ImageConcurrency.OverflowMode != ImageConcurrencyOverflowModeReject {
		t.Fatalf("image_concurrency.overflow_mode = %q, want %q", cfg.Gateway.ImageConcurrency.OverflowMode, ImageConcurrencyOverflowModeReject)
placeholder
	if cfg.Gateway.ImageConcurrency.WaitTimeoutSeconds != 30 {
		t.Fatalf("image_concurrency.wait_timeout_seconds = %d, want 30", cfg.Gateway.ImageConcurrency.WaitTimeoutSeconds)
placeholder
	if cfg.Gateway.ImageConcurrency.MaxWaitingRequests != 100 {
		t.Fatalf("image_concurrency.max_waiting_requests = %d, want 100", cfg.Gateway.ImageConcurrency.MaxWaitingRequests)
placeholder
	if cfg.Gateway.ImageStreamDataIntervalTimeout <= cfg.Gateway.StreamDataIntervalTimeout {
		t.Fatalf("image stream timeout = %d, want greater than ordinary stream timeout %d", cfg.Gateway.ImageStreamDataIntervalTimeout, cfg.Gateway.StreamDataIntervalTimeout)
placeholder
placeholder
