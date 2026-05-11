package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type gatewayTTLSettingRepo struct {
	data map[string]string
placeholder

func (r *gatewayTTLSettingRepo) Get(context.Context, string) (*Setting, error) {
	return nil, ErrSettingNotFound
placeholder

func (r *gatewayTTLSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	if r == nil {
		return "", ErrSettingNotFound
placeholder
	v, ok := r.data[key]
	if !ok {
		return "", ErrSettingNotFound
placeholder
	return v, nil
placeholder

func (r *gatewayTTLSettingRepo) Set(_ context.Context, key, value string) error {
	if r == nil {
		return errors.New("setting repo is nil")
placeholder
	if r.data == nil {
		r.data = map[string]string{placeholder
placeholder
	r.data[key] = value
	return nil
placeholder

func (r *gatewayTTLSettingRepo) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	result := make(map[string]string)
	if r == nil {
		return result, nil
placeholder
	for _, key := range keys {
		if v, ok := r.data[key]; ok {
			result[key] = v
	placeholder
placeholder
	return result, nil
placeholder

func (r *gatewayTTLSettingRepo) SetMultiple(_ context.Context, settings map[string]string) error {
	if r == nil {
		return errors.New("setting repo is nil")
placeholder
	if r.data == nil {
		r.data = map[string]string{placeholder
placeholder
	for key, value := range settings {
		r.data[key] = value
placeholder
	return nil
placeholder

func (r *gatewayTTLSettingRepo) GetAll(context.Context) (map[string]string, error) {
	result := make(map[string]string)
	if r == nil {
		return result, nil
placeholder
	for key, value := range r.data {
		result[key] = value
placeholder
	return result, nil
placeholder

func (r *gatewayTTLSettingRepo) Delete(_ context.Context, key string) error {
	if r != nil {
		delete(r.data, key)
placeholder
	return nil
placeholder

func assertJSONTokenOrder(t *testing.T, body string, tokens ...string) {
placeholder

	last := -1
	for _, token := range tokens {
		pos := strings.Index(body, token)
		require.NotEqualf(t, -1, pos, "missing token %s in body %s", token, body)
		require.Greaterf(t, pos, last, "token %s should appear after previous tokens in body %s", token, body)
		last = pos
placeholder
placeholder

func TestReplaceModelInBody_PreservesTopLevelFieldOrder(t *testing.T) {
	svc := &GatewayService{placeholder
	body := []byte(`{"alpha":1,"model":"claude-3-5-sonnet-latest","messages":[],"omega":2placeholder`)

	result := svc.replaceModelInBody(body, "claude-3-5-sonnet-20241022")
	resultStr := string(result)

	assertJSONTokenOrder(t, resultStr, `"alpha"`, `"model"`, `"messages"`, `"omega"`)
	require.Contains(t, resultStr, `"model":"claude-3-5-sonnet-20241022"`)
placeholder

func TestNormalizeClaudeOAuthRequestBody_PreservesTopLevelFieldOrder(t *testing.T) {
	body := []byte(`{"alpha":1,"model":"claude-3-5-sonnet-latest","temperature":0.2,"system":"You are OpenCode, the best coding agent on the planet.","messages":[],"tool_choice":{"type":"auto"placeholder,"omega":2placeholder`)

	result, modelID := normalizeClaudeOAuthRequestBody(body, "claude-3-5-sonnet-latest", claudeOAuthNormalizeOptions{
		injectMetadata: true,
		metadataUserID: "user-1",
placeholder)
	resultStr := string(result)

	require.Equal(t, claude.NormalizeModelID("claude-3-5-sonnet-latest"), modelID)
	assertJSONTokenOrder(t, resultStr, `"alpha"`, `"model"`, `"temperature"`, `"system"`, `"messages"`, `"omega"`, `"tools"`, `"metadata"`, `"max_tokens"`)
	require.Contains(t, resultStr, `"temperature":0.2`)
	require.NotContains(t, resultStr, `"tool_choice"`)
	require.Contains(t, resultStr, `"system":"`+claudeCodeSystemPrompt+`"`)
	require.Contains(t, resultStr, `"tools":[]`)
	require.Contains(t, resultStr, `"metadata":{"user_id":"user-1"placeholder`)
	require.Contains(t, resultStr, `"max_tokens":128000`)
placeholder

func TestInjectClaudeCodePrompt_PreservesFieldOrder(t *testing.T) {
	body := []byte(`{"alpha":1,"system":[{"id":"block-1","type":"text","text":"Custom"placeholder],"messages":[],"omega":2placeholder`)

	result := injectClaudeCodePrompt(body, []any{
		map[string]any{"id": "block-1", "type": "text", "text": "Custom"placeholder,
placeholder)
	resultStr := string(result)

	assertJSONTokenOrder(t, resultStr, `"alpha"`, `"system"`, `"messages"`, `"omega"`)
	require.Contains(t, resultStr, `{"id":"block-1","type":"text","text":"`+claudeCodeSystemPrompt+`\n\nCustom"placeholder`)
placeholder

func TestEnforceCacheControlLimit_PreservesTopLevelFieldOrder(t *testing.T) {
	body := []byte(`{"alpha":1,"system":[{"type":"text","text":"s1","cache_control":{"type":"ephemeral"placeholderplaceholder,{"type":"text","text":"s2","cache_control":{"type":"ephemeral"placeholderplaceholder],"messages":[{"role":"user","content":[{"type":"text","text":"m1","cache_control":{"type":"ephemeral"placeholderplaceholder,{"type":"text","text":"m2","cache_control":{"type":"ephemeral"placeholderplaceholder,{"type":"text","text":"m3","cache_control":{"type":"ephemeral"placeholderplaceholder]placeholder],"omega":2placeholder`)

	result := enforceCacheControlLimit(body)
	resultStr := string(result)

	assertJSONTokenOrder(t, resultStr, `"alpha"`, `"system"`, `"messages"`, `"omega"`)
	require.Equal(t, 4, strings.Count(resultStr, `"cache_control"`))
placeholder

func TestEnforceCacheControlLimit_CountsToolsAndPreservesMessageAnchorsFirst(t *testing.T) {
	body := []byte(`{"alpha":1,"system":[{"type":"text","text":"sys","cache_control":{"type":"ephemeral"placeholderplaceholder],"messages":[{"role":"user","content":[{"type":"text","text":"m1","cache_control":{"type":"ephemeral"placeholderplaceholder,{"type":"text","text":"m2","cache_control":{"type":"ephemeral"placeholderplaceholder,{"type":"text","text":"m3","cache_control":{"type":"ephemeral"placeholderplaceholder]placeholder],"tools":[{"name":"a","input_schema":{placeholder,"cache_control":{"type":"ephemeral"placeholderplaceholder],"omega":2placeholder`)

	result := enforceCacheControlLimit(body)
	resultStr := string(result)

	assertJSONTokenOrder(t, resultStr, `"alpha"`, `"system"`, `"messages"`, `"tools"`, `"omega"`)
	require.Equal(t, 4, strings.Count(resultStr, `"cache_control"`))
	require.True(t, gjson.GetBytes(result, "system.0.cache_control").Exists())
	require.True(t, gjson.GetBytes(result, "messages.0.content.0.cache_control").Exists())
	require.True(t, gjson.GetBytes(result, "messages.0.content.1.cache_control").Exists())
	require.True(t, gjson.GetBytes(result, "messages.0.content.2.cache_control").Exists())
	require.False(t, gjson.GetBytes(result, "tools.0.cache_control").Exists())
placeholder

func TestInjectAnthropicCacheControlTTL1h_OnlyUpdatesExistingEphemeralCacheControl(t *testing.T) {
	body := []byte(`{"alpha":1,"cache_control":{"type":"ephemeral"placeholder,"system":[{"type":"text","text":"sys","cache_control":{"type":"ephemeral","ttl":"5m"placeholderplaceholder,{"type":"text","text":"plain"placeholder],"messages":[{"role":"user","content":[{"type":"text","text":"hi","cache_control":{"type":"ephemeral"placeholderplaceholder,{"type":"text","text":"non","cache_control":{"type":"persistent","ttl":"5m"placeholderplaceholder]placeholder],"tools":[{"name":"a","input_schema":{placeholder,"cache_control":{"type":"ephemeral"placeholderplaceholder],"omega":2placeholder`)

	result := injectAnthropicCacheControlTTL1h(body)
	resultStr := string(result)

	assertJSONTokenOrder(t, resultStr, `"alpha"`, `"cache_control"`, `"system"`, `"messages"`, `"tools"`, `"omega"`)
	require.Equal(t, "1h", gjson.GetBytes(result, "cache_control.ttl").String())
	require.Equal(t, "1h", gjson.GetBytes(result, "system.0.cache_control.ttl").String())
	require.False(t, gjson.GetBytes(result, "system.1.cache_control").Exists())
	require.Equal(t, "1h", gjson.GetBytes(result, "messages.0.content.0.cache_control.ttl").String())
	require.Equal(t, "5m", gjson.GetBytes(result, "messages.0.content.1.cache_control.ttl").String())
	require.Equal(t, "1h", gjson.GetBytes(result, "tools.0.cache_control.ttl").String())
placeholder

func TestGatewayCacheTTLGlobalSetting_TargetResolution(t *testing.T) {
	repo := &gatewayTTLSettingRepo{data: map[string]string{
		SettingKeyEnableAnthropicCacheTTL1hInjection: "true",
placeholderplaceholder
	gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{placeholder)
	svc := &GatewayService{
		settingService: NewSettingService(repo, &config.Config{placeholder),
placeholder
	account := &Account{Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder

	target, ok := svc.resolveCacheTTLUsageOverrideTarget(context.Background(), account)
	require.True(t, ok)
	require.Equal(t, cacheTTLTarget5m, target)

	account.Extra = map[string]any{
		"cache_ttl_override_enabled": true,
		"cache_ttl_override_target":  "1h",
placeholder
	target, ok = svc.resolveCacheTTLUsageOverrideTarget(context.Background(), account)
	require.True(t, ok)
	require.Equal(t, cacheTTLTarget1h, target)
placeholder

func TestGatewayCacheTTLGlobalSetting_RequestInjectionScope(t *testing.T) {
	repo := &gatewayTTLSettingRepo{data: map[string]string{
		SettingKeyEnableAnthropicCacheTTL1hInjection: "true",
placeholderplaceholder
	gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{placeholder)
	svc := &GatewayService{
		settingService: NewSettingService(repo, &config.Config{placeholder),
placeholder

	require.True(t, svc.shouldInjectAnthropicCacheTTL1h(context.Background(), &Account{Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder))
	require.True(t, svc.shouldInjectAnthropicCacheTTL1h(context.Background(), &Account{Platform: PlatformAnthropic, Type: AccountTypeSetupTokenplaceholder))
	require.False(t, svc.shouldInjectAnthropicCacheTTL1h(context.Background(), &Account{Platform: PlatformAnthropic, Type: AccountTypeAPIKeyplaceholder))
	require.False(t, svc.shouldInjectAnthropicCacheTTL1h(context.Background(), &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder))

	repo.data[SettingKeyEnableAnthropicCacheTTL1hInjection] = "false"
	gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{placeholder)
	require.False(t, svc.shouldInjectAnthropicCacheTTL1h(context.Background(), &Account{Platform: PlatformAnthropic, Type: AccountTypeOAuthplaceholder))
placeholder
