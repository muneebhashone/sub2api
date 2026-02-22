package config

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func resetViperWithJWTSecret(t *testing.T) {
placeholder
	viper.Reset()
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
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
	if cfg.Gateway.Scheduling.SlotCleanupInterval != 30*time.Second {
		t.Fatalf("SlotCleanupInterval = %v, want 30s", cfg.Gateway.Scheduling.SlotCleanupInterval)
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
	cfg.LinuxDo.UsePKCE = false

	cfg.LinuxDo.FrontendRedirectURL = "javascript:alert(1)"
	err = cfg.Validate()
	if err == nil {
		t.Fatalf("Validate() expected error for javascript scheme, got nil")
placeholder
	if !strings.Contains(err.Error(), "linuxdo_connect.frontend_redirect_url") {
		t.Fatalf("Validate() expected frontend_redirect_url error, got: %v", err)
placeholder
placeholder

func TestValidateLinuxDoPKCERequiredForPublicClient(t *testing.T) {
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
	if err == nil {
		t.Fatalf("Validate() expected error when token_auth_method=none and use_pkce=false, got nil")
placeholder
	if !strings.Contains(err.Error(), "linuxdo_connect.use_pkce") {
		t.Fatalf("Validate() expected use_pkce error, got: %v", err)
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
			name:    "csp policy required",
			mutate:  func(c *Config) { c.Security.CSP.Enabled = true; c.Security.CSP.Policy = "" placeholder,
			wantErr: "security.csp.policy",
	placeholder,
		{
			name: "linuxdo client id required",
			mutate: func(c *Config) {
				c.LinuxDo.Enabled = true
				c.LinuxDo.ClientID = ""
		placeholder,
			wantErr: "linuxdo_connect.client_id",
	placeholder,
		{
			name: "linuxdo token auth method",
			mutate: func(c *Config) {
				c.LinuxDo.Enabled = true
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

func TestSoraCurlCFFISidecarDefaults(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	if !cfg.Sora.Client.CurlCFFISidecar.Enabled {
		t.Fatalf("Sora curl_cffi sidecar should be enabled by default")
placeholder
	if cfg.Sora.Client.CloudflareChallengeCooldownSeconds <= 0 {
		t.Fatalf("Sora cloudflare challenge cooldown should be positive by default")
placeholder
	if cfg.Sora.Client.CurlCFFISidecar.BaseURL == "" {
		t.Fatalf("Sora curl_cffi sidecar base_url should not be empty by default")
placeholder
	if cfg.Sora.Client.CurlCFFISidecar.Impersonate == "" {
		t.Fatalf("Sora curl_cffi sidecar impersonate should not be empty by default")
placeholder
	if !cfg.Sora.Client.CurlCFFISidecar.SessionReuseEnabled {
		t.Fatalf("Sora curl_cffi sidecar session reuse should be enabled by default")
placeholder
	if cfg.Sora.Client.CurlCFFISidecar.SessionTTLSeconds <= 0 {
		t.Fatalf("Sora curl_cffi sidecar session ttl should be positive by default")
placeholder
placeholder

func TestValidateSoraCurlCFFISidecarRequired(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Sora.Client.CurlCFFISidecar.Enabled = false
	err = cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "sora.client.curl_cffi_sidecar.enabled must be true") {
		t.Fatalf("Validate() error = %v, want sidecar enabled error", err)
placeholder
placeholder

func TestValidateSoraCurlCFFISidecarBaseURLRequired(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Sora.Client.CurlCFFISidecar.BaseURL = "   "
	err = cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "sora.client.curl_cffi_sidecar.base_url is required") {
		t.Fatalf("Validate() error = %v, want sidecar base_url required error", err)
placeholder
placeholder

func TestValidateSoraCurlCFFISidecarSessionTTLNonNegative(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Sora.Client.CurlCFFISidecar.SessionTTLSeconds = -1
	err = cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "sora.client.curl_cffi_sidecar.session_ttl_seconds must be non-negative") {
		t.Fatalf("Validate() error = %v, want sidecar session ttl error", err)
placeholder
placeholder

func TestValidateSoraCloudflareChallengeCooldownNonNegative(t *testing.T) {
	resetViperWithJWTSecret(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
placeholder

	cfg.Sora.Client.CloudflareChallengeCooldownSeconds = -1
	err = cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "sora.client.cloudflare_challenge_cooldown_seconds must be non-negative") {
		t.Fatalf("Validate() error = %v, want cloudflare cooldown error", err)
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
	if cfg.Gateway.UsageRecord.OverflowPolicy != UsageRecordOverflowPolicySample {
		t.Fatalf("overflow_policy = %s, want %s", cfg.Gateway.UsageRecord.OverflowPolicy, UsageRecordOverflowPolicySample)
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
