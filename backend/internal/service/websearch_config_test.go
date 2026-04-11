package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// --- validateWebSearchConfig ---

func TestValidateWebSearchConfig_Nil(t *testing.T) {
	require.NoError(t, validateWebSearchConfig(nil))
placeholder

func TestValidateWebSearchConfig_Valid(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Enabled: true,
		Providers: []WebSearchProviderConfig{
			{Type: "brave", Priority: 1, QuotaLimit: 1000, QuotaRefreshInterval: "monthly"placeholder,
			{Type: "tavily", Priority: 2, QuotaLimit: 500, QuotaRefreshInterval: "daily"placeholder,
	placeholder,
placeholder
	require.NoError(t, validateWebSearchConfig(cfg))
placeholder

func TestValidateWebSearchConfig_TooManyProviders(t *testing.T) {
	cfg := &WebSearchEmulationConfig{Providers: make([]WebSearchProviderConfig, 11)placeholder
	for i := range cfg.Providers {
		cfg.Providers[i] = WebSearchProviderConfig{Type: "brave"placeholder
placeholder
	err := validateWebSearchConfig(cfg)
	require.ErrorContains(t, err, "too many providers")
placeholder

func TestValidateWebSearchConfig_InvalidType(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "bing"placeholderplaceholder,
placeholder
	require.ErrorContains(t, validateWebSearchConfig(cfg), "invalid type")
placeholder

func TestValidateWebSearchConfig_InvalidQuotaInterval(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "brave", QuotaRefreshInterval: "hourly"placeholderplaceholder,
placeholder
	require.ErrorContains(t, validateWebSearchConfig(cfg), "invalid quota_refresh_interval")
placeholder

func TestValidateWebSearchConfig_NegativeQuotaLimit(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "brave", QuotaLimit: -1placeholderplaceholder,
placeholder
	require.ErrorContains(t, validateWebSearchConfig(cfg), "quota_limit must be >= 0")
placeholder

func TestValidateWebSearchConfig_DuplicateType(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{
			{Type: "brave", Priority: 1placeholder,
			{Type: "brave", Priority: 2placeholder,
	placeholder,
placeholder
	require.ErrorContains(t, validateWebSearchConfig(cfg), "duplicate type")
placeholder

func TestValidateWebSearchConfig_EmptyQuotaInterval(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "brave", QuotaRefreshInterval: ""placeholderplaceholder,
placeholder
	require.NoError(t, validateWebSearchConfig(cfg))
placeholder

func TestValidateWebSearchConfig_ZeroQuotaLimit(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "brave", QuotaLimit: 0placeholderplaceholder,
placeholder
	require.NoError(t, validateWebSearchConfig(cfg))
placeholder

// --- parseWebSearchConfigJSON ---

func TestParseWebSearchConfigJSON_ValidJSON(t *testing.T) {
	raw := `{"enabled":true,"providers":[{"type":"brave","api_key":"sk-xxx"placeholder]placeholder`
	cfg := parseWebSearchConfigJSON(raw)
	require.True(t, cfg.Enabled)
	require.Len(t, cfg.Providers, 1)
	require.Equal(t, "brave", cfg.Providers[0].Type)
placeholder

func TestParseWebSearchConfigJSON_EmptyString(t *testing.T) {
	cfg := parseWebSearchConfigJSON("")
	require.False(t, cfg.Enabled)
	require.Empty(t, cfg.Providers)
placeholder

func TestParseWebSearchConfigJSON_InvalidJSON(t *testing.T) {
	cfg := parseWebSearchConfigJSON("not{json")
	require.False(t, cfg.Enabled)
	require.Empty(t, cfg.Providers)
placeholder

// --- SanitizeWebSearchConfig ---

func TestSanitizeWebSearchConfig_MaskAPIKey(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Enabled: true,
		Providers: []WebSearchProviderConfig{
			{Type: "brave", APIKey: "sk-secret-xxx"placeholder,
	placeholder,
placeholder
	out := SanitizeWebSearchConfig(cfg)
	require.Equal(t, "", out.Providers[0].APIKey)
	require.True(t, out.Providers[0].APIKeyConfigured)
placeholder

func TestSanitizeWebSearchConfig_NoAPIKey(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "brave", APIKey: ""placeholderplaceholder,
placeholder
	out := SanitizeWebSearchConfig(cfg)
	require.Equal(t, "", out.Providers[0].APIKey)
	require.False(t, out.Providers[0].APIKeyConfigured)
placeholder

func TestSanitizeWebSearchConfig_Nil(t *testing.T) {
	require.Nil(t, SanitizeWebSearchConfig(nil))
placeholder

func TestSanitizeWebSearchConfig_PreservesOtherFields(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Enabled: true,
		Providers: []WebSearchProviderConfig{
			{Type: "brave", APIKey: "secret", Priority: 10, QuotaLimit: 1000placeholder,
	placeholder,
placeholder
	out := SanitizeWebSearchConfig(cfg)
	require.True(t, out.Enabled)
	require.Equal(t, 10, out.Providers[0].Priority)
	require.Equal(t, int64(1000), out.Providers[0].QuotaLimit)
placeholder

func TestSanitizeWebSearchConfig_DoesNotMutateOriginal(t *testing.T) {
	cfg := &WebSearchEmulationConfig{
		Providers: []WebSearchProviderConfig{{Type: "brave", APIKey: "secret"placeholderplaceholder,
placeholder
	_ = SanitizeWebSearchConfig(cfg)
	require.Equal(t, "secret", cfg.Providers[0].APIKey)
placeholder
