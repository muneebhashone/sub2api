package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/stretchr/testify/require"
)

func TestGetCodexRestrictionPolicy(t *testing.T) {
	svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{
		SettingKeyMinCodexVersion:       "0.141.0",
		SettingKeyMaxCodexVersion:       "0.200.0",
		SettingKeyCodexCLIOnlyWhitelist: `[{"originator":"opencode","ua_contains":["opencode/"]placeholder]`,
		SettingKeyCodexCLIOnlyBlacklist: `[{"originator":"evil"placeholder]`,
placeholderplaceholder, &config.Config{placeholder)

	pol := svc.GetCodexRestrictionPolicy(context.Background())
	require.Equal(t, "0.141.0", pol.MinCodexVersion)
	require.Equal(t, "0.200.0", pol.MaxCodexVersion)
	require.Len(t, pol.Whitelist, 1)
	require.Equal(t, "opencode", pol.Whitelist[0].Originator)
	require.Equal(t, []string{"opencode/"placeholder, pol.Whitelist[0].UAContains)
	require.Len(t, pol.Blacklist, 1)
	require.Equal(t, "evil", pol.Blacklist[0].Originator)
placeholder

func TestGetCodexRestrictionPolicy_DefaultsSafe(t *testing.T) {
	svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{placeholderplaceholder, &config.Config{placeholder)

	pol := svc.GetCodexRestrictionPolicy(context.Background())
	require.Empty(t, pol.MinCodexVersion)
	require.Empty(t, pol.Whitelist)
	require.Empty(t, pol.Blacklist)
placeholder

func TestGetCodexRestrictionPolicy_InvalidJSONSafe(t *testing.T) {
	svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{
		SettingKeyCodexCLIOnlyWhitelist: "not-json",
		SettingKeyCodexCLIOnlyBlacklist: "{bad",
placeholderplaceholder, &config.Config{placeholder)

	pol := svc.GetCodexRestrictionPolicy(context.Background())
	require.Empty(t, pol.Whitelist, "非法 JSON → 安全空名单")
	require.Empty(t, pol.Blacklist, "非法 JSON → 安全空名单")
placeholder

type codexPolicyMigrationRepoStub struct {
	values map[string]string
	sets   map[string]string
placeholder

func (s *codexPolicyMigrationRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unused")
placeholder
func (s *codexPolicyMigrationRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
placeholder
	return "", ErrSettingNotFound
placeholder
func (s *codexPolicyMigrationRepoStub) Set(ctx context.Context, key, value string) error {
	if s.sets == nil {
		s.sets = map[string]string{placeholder
placeholder
	s.sets[key] = value
	s.values[key] = value
	return nil
placeholder
func (s *codexPolicyMigrationRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unused")
placeholder
func (s *codexPolicyMigrationRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unused")
placeholder
func (s *codexPolicyMigrationRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unused")
placeholder
func (s *codexPolicyMigrationRepoStub) Delete(ctx context.Context, key string) error {
	panic("unused")
placeholder

func TestMigrateOpenAIAllowClaudeCodeCodexPluginSetting(t *testing.T) {
	t.Run("legacy true appends Claude Code entry to whitelist", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyOpenAIAllowClaudeCodeCodexPlugin: "true",
			SettingKeyCodexCLIOnlyWhitelist:            `[{"originator":"opencode","ua_contains":["opencode/"]placeholder]`,
	placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)

		require.NoError(t, svc.MigrateOpenAIAllowClaudeCodeCodexPluginSetting(context.Background()))

		raw := repo.sets[SettingKeyCodexCLIOnlyWhitelist]
		require.NotEmpty(t, raw)
		var entries []struct {
			Originator string   `json:"originator"`
			UAContains []string `json:"ua_contains"`
	placeholder
		require.NoError(t, json.Unmarshal([]byte(raw), &entries))
		require.Len(t, entries, 2)
		require.Equal(t, "opencode", entries[0].Originator)
		require.Equal(t, "Claude Code", entries[1].Originator)
		require.Equal(t, []string{"Claude Code/"placeholder, entries[1].UAContains)
placeholder)

	t.Run("legacy true does not duplicate existing Claude Code entry", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyOpenAIAllowClaudeCodeCodexPlugin: "true",
			SettingKeyCodexCLIOnlyWhitelist:            `[{"originator":"Claude Code","ua_contains":["Claude Code/"]placeholder]`,
	placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)

		require.NoError(t, svc.MigrateOpenAIAllowClaudeCodeCodexPluginSetting(context.Background()))

		_, wrote := repo.sets[SettingKeyCodexCLIOnlyWhitelist]
		require.False(t, wrote)
placeholder)
placeholder

func TestGetCodexRestrictionPolicy_AllowAppServerClients(t *testing.T) {
	t.Run("显式 true 开启", func(t *testing.T) {
		svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyCodexCLIOnlyAllowAppServerClients: "true",
placeholder &config.Config{placeholder)
		require.True(t, svc.GetCodexRestrictionPolicy(context.Background()).AllowAppServerClients)
placeholder)
	t.Run("缺失默认 false", func(t *testing.T) {
		svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{placeholderplaceholder, &config.Config{placeholder)
		require.False(t, svc.GetCodexRestrictionPolicy(context.Background()).AllowAppServerClients)
placeholder)
	t.Run("非 true 值视为 false", func(t *testing.T) {
		svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyCodexCLIOnlyAllowAppServerClients: "1",
placeholder &config.Config{placeholder)
		require.False(t, svc.GetCodexRestrictionPolicy(context.Background()).AllowAppServerClients)
placeholder)
placeholder

func TestGetCodexRestrictionPolicy_EngineFingerprintSignals(t *testing.T) {
	t.Run("键缺失 → 默认种子(只勾x-codex-)", func(t *testing.T) {
		svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{placeholderplaceholder, &config.Config{placeholder)
		pol := svc.GetCodexRestrictionPolicy(context.Background())
		require.True(t, len(pol.EngineFingerprintSignals) > 0)
		require.True(t, openaiEngineSignalsEqual(pol.EngineFingerprintSignals, openai.DefaultEngineFingerprintSignals))
placeholder)
	t.Run("显式配置 → 原样采用", func(t *testing.T) {
		raw := `[{"type":"header_exact","match":["session-id"],"required":trueplaceholder]`
		svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyCodexCLIOnlyEngineFingerprintSignals: raw,
placeholder &config.Config{placeholder)
		pol := svc.GetCodexRestrictionPolicy(context.Background())
		require.Len(t, pol.EngineFingerprintSignals, 1)
		require.Equal(t, "session-id", pol.EngineFingerprintSignals[0].Match[0])
placeholder)
	t.Run("非法JSON → 回落默认种子", func(t *testing.T) {
		svc := NewSettingService(&codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyCodexCLIOnlyEngineFingerprintSignals: "not json",
placeholder &config.Config{placeholder)
		pol := svc.GetCodexRestrictionPolicy(context.Background())
		require.True(t, openaiEngineSignalsEqual(pol.EngineFingerprintSignals, openai.DefaultEngineFingerprintSignals))
placeholder)
placeholder

func openaiEngineSignalsEqual(a, b []openai.EngineFingerprintSignal) bool {
	return reflect.DeepEqual(a, b)
placeholder

func TestMigrateCodexBodyFingerprintToSignals(t *testing.T) {
	t.Run("信号键已存在 → 不动", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyCodexCLIOnlyEngineFingerprintSignals:   `[{"type":"header_exact","match":["session-id"],"required":trueplaceholder]`,
			SettingKeyCodexCLIOnlyAllowBodyEngineFingerprint: "true",
	placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)
		require.NoError(t, svc.MigrateCodexBodyFingerprintToSignals(context.Background()))
		require.Equal(t, `[{"type":"header_exact","match":["session-id"],"required":trueplaceholder]`, repo.values[SettingKeyCodexCLIOnlyEngineFingerprintSignals])
placeholder)
	t.Run("信号键缺失 + 旧body=true → 写种子且body行Required=true", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyCodexCLIOnlyAllowBodyEngineFingerprint: "true",
	placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)
		require.NoError(t, svc.MigrateCodexBodyFingerprintToSignals(context.Background()))
		sigs, ok := openai.ParseEngineFingerprintSignals(repo.values[SettingKeyCodexCLIOnlyEngineFingerprintSignals])
		require.True(t, ok)
		var bodyReq bool
		for _, s := range sigs {
			if s.Type == openai.FingerprintSignalBodyPath {
				bodyReq = s.Required
		placeholder
	placeholder
		require.True(t, bodyReq)
placeholder)
	t.Run("信号键缺失 + 旧body=false/缺 → 写种子(body不勾)", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)
		require.NoError(t, svc.MigrateCodexBodyFingerprintToSignals(context.Background()))
		require.Equal(t, openai.DefaultEngineFingerprintSignalsJSON(), repo.values[SettingKeyCodexCLIOnlyEngineFingerprintSignals])
placeholder)
placeholder

func TestMigrateGrokDefaultTextModel(t *testing.T) {
	t.Run("upgrades legacy built-in default", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyGrokDefaultTextModel: "grok-4.5",
	placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)
		require.NoError(t, svc.MigrateGrokDefaultTextModel(context.Background()))
		require.Equal(t, "grok-4.6", repo.values[SettingKeyGrokDefaultTextModel])
		require.Equal(t, "grok-4.6", repo.sets[SettingKeyGrokDefaultTextModel])
placeholder)

	t.Run("does not overwrite an explicit model", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{
			SettingKeyGrokDefaultTextModel: "grok-4.3",
	placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)
		require.NoError(t, svc.MigrateGrokDefaultTextModel(context.Background()))
		require.Equal(t, "grok-4.3", repo.values[SettingKeyGrokDefaultTextModel])
		_, wrote := repo.sets[SettingKeyGrokDefaultTextModel]
		require.False(t, wrote)
placeholder)

	t.Run("missing setting is left for normal defaults", func(t *testing.T) {
		repo := &codexPolicyMigrationRepoStub{values: map[string]string{placeholderplaceholder
		svc := NewSettingService(repo, &config.Config{placeholder)
		require.NoError(t, svc.MigrateGrokDefaultTextModel(context.Background()))
		_, wrote := repo.sets[SettingKeyGrokDefaultTextModel]
		require.False(t, wrote)
placeholder)
placeholder
