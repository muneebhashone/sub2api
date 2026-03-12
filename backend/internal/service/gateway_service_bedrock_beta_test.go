package service

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

type betaPolicySettingRepoStub struct {
	values map[string]string
placeholder

func (s *betaPolicySettingRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *betaPolicySettingRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
placeholder
	return "", ErrSettingNotFound
placeholder

func (s *betaPolicySettingRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *betaPolicySettingRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *betaPolicySettingRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *betaPolicySettingRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *betaPolicySettingRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

func TestResolveBedrockBetaTokensForRequest_BlocksOnOriginalAnthropicToken(t *testing.T) {
	settings := &BetaPolicySettings{
		Rules: []BetaPolicyRule{
			{
				BetaToken:    "advanced-tool-use-2025-11-20",
				Action:       BetaPolicyActionBlock,
				Scope:        BetaPolicyScopeAll,
				ErrorMessage: "advanced tool use is blocked",
		placeholder,
	placeholder,
placeholder
	raw, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("marshal settings: %v", err)
placeholder

	svc := &GatewayService{
		settingService: NewSettingService(
			&betaPolicySettingRepoStub{values: map[string]string{
				SettingKeyBetaPolicySettings: string(raw),
	placeholder
			&config.Config{placeholder,
		),
placeholder
	account := &Account{Platform: PlatformAnthropic, Type: AccountTypeBedrockplaceholder

	_, err = svc.resolveBedrockBetaTokensForRequest(
		context.Background(),
		account,
		"advanced-tool-use-2025-11-20",
		[]byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
		"us.anthropic.claude-opus-4-6-v1",
	)
	if err == nil {
		t.Fatal("expected raw advanced-tool-use token to be blocked before Bedrock transform")
placeholder
	if err.Error() != "advanced tool use is blocked" {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

func TestResolveBedrockBetaTokensForRequest_FiltersAfterBedrockTransform(t *testing.T) {
	settings := &BetaPolicySettings{
		Rules: []BetaPolicyRule{
			{
				BetaToken: "tool-search-tool-2025-10-19",
				Action:    BetaPolicyActionFilter,
				Scope:     BetaPolicyScopeAll,
		placeholder,
	placeholder,
placeholder
	raw, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("marshal settings: %v", err)
placeholder

	svc := &GatewayService{
		settingService: NewSettingService(
			&betaPolicySettingRepoStub{values: map[string]string{
				SettingKeyBetaPolicySettings: string(raw),
	placeholder
			&config.Config{placeholder,
		),
placeholder
	account := &Account{Platform: PlatformAnthropic, Type: AccountTypeBedrockplaceholder

	betaTokens, err := svc.resolveBedrockBetaTokensForRequest(
		context.Background(),
		account,
		"advanced-tool-use-2025-11-20",
		[]byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
		"us.anthropic.claude-opus-4-6-v1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	for _, token := range betaTokens {
		if token == "tool-search-tool-2025-10-19" {
			t.Fatalf("expected transformed Bedrock token to be filtered")
	placeholder
placeholder
placeholder

// TestResolveBedrockBetaTokensForRequest_BlocksBodyAutoInjectedThinking 验证：
// 管理员 block 了 interleaved-thinking，客户端不在 header 中带该 token，
// 但请求体包含 thinking 字段 → 自动注入后应被 block。
func TestResolveBedrockBetaTokensForRequest_BlocksBodyAutoInjectedThinking(t *testing.T) {
	settings := &BetaPolicySettings{
		Rules: []BetaPolicyRule{
			{
				BetaToken:    "placeholder",
				Action:       BetaPolicyActionBlock,
				Scope:        BetaPolicyScopeAll,
				ErrorMessage: "thinking is blocked",
		placeholder,
	placeholder,
placeholder
	raw, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("marshal settings: %v", err)
placeholder

	svc := &GatewayService{
		settingService: NewSettingService(
			&betaPolicySettingRepoStub{values: map[string]string{
				SettingKeyBetaPolicySettings: string(raw),
	placeholder
			&config.Config{placeholder,
		),
placeholder
	account := &Account{Platform: PlatformAnthropic, Type: AccountTypeBedrockplaceholder

	// header 中不带 beta token，但 body 中有 thinking 字段
	_, err = svc.resolveBedrockBetaTokensForRequest(
		context.Background(),
		account,
		"", // 空 header
		[]byte(`{"thinking":{"type":"enabled","budget_tokens":10000placeholder,"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
		"us.anthropic.claude-opus-4-6-v1",
	)
	if err == nil {
		t.Fatal("expected body-injected interleaved-thinking to be blocked")
placeholder
	if err.Error() != "thinking is blocked" {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

// TestResolveBedrockBetaTokensForRequest_BlocksBodyAutoInjectedToolSearch 验证：
// 管理员 block 了 tool-search-tool，客户端不在 header 中带 beta token，
// 但请求体包含 tool search 工具 → 自动注入后应被 block。
func TestResolveBedrockBetaTokensForRequest_BlocksBodyAutoInjectedToolSearch(t *testing.T) {
	settings := &BetaPolicySettings{
		Rules: []BetaPolicyRule{
			{
				BetaToken:    "tool-search-tool-2025-10-19",
				Action:       BetaPolicyActionBlock,
				Scope:        BetaPolicyScopeAll,
				ErrorMessage: "tool search is blocked",
		placeholder,
	placeholder,
placeholder
	raw, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("marshal settings: %v", err)
placeholder

	svc := &GatewayService{
		settingService: NewSettingService(
			&betaPolicySettingRepoStub{values: map[string]string{
				SettingKeyBetaPolicySettings: string(raw),
	placeholder
			&config.Config{placeholder,
		),
placeholder
	account := &Account{Platform: PlatformAnthropic, Type: AccountTypeBedrockplaceholder

	// header 中不带 beta token，但 body 中有 tool_search_tool 工具
	_, err = svc.resolveBedrockBetaTokensForRequest(
		context.Background(),
		account,
		"",
		[]byte(`{"tools":[{"type":"tool_search_tool_regex_20251119","name":"search"placeholder],"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
		"us.anthropic.claude-sonnet-4-6",
	)
	if err == nil {
		t.Fatal("expected body-injected tool-search-tool to be blocked")
placeholder
	if err.Error() != "tool search is blocked" {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

// TestResolveBedrockBetaTokensForRequest_PassesWhenNoBlockRuleMatches 验证：
// body 自动注入的 token 如果没有对应的 block 规则，应正常通过。
func TestResolveBedrockBetaTokensForRequest_PassesWhenNoBlockRuleMatches(t *testing.T) {
	settings := &BetaPolicySettings{
		Rules: []BetaPolicyRule{
			{
				BetaToken:    "computer-use-2025-11-24",
				Action:       BetaPolicyActionBlock,
				Scope:        BetaPolicyScopeAll,
				ErrorMessage: "computer use is blocked",
		placeholder,
	placeholder,
placeholder
	raw, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("marshal settings: %v", err)
placeholder

	svc := &GatewayService{
		settingService: NewSettingService(
			&betaPolicySettingRepoStub{values: map[string]string{
				SettingKeyBetaPolicySettings: string(raw),
	placeholder
			&config.Config{placeholder,
		),
placeholder
	account := &Account{Platform: PlatformAnthropic, Type: AccountTypeBedrockplaceholder

	// body 中有 thinking（会注入 interleaved-thinking），但 block 规则只针对 computer-use
	tokens, err := svc.resolveBedrockBetaTokensForRequest(
		context.Background(),
		account,
		"",
		[]byte(`{"thinking":{"type":"enabled","budget_tokens":10000placeholder,"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
		"us.anthropic.claude-opus-4-6-v1",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	found := false
	for _, token := range tokens {
		if token == "placeholder" {
			found = true
	placeholder
placeholder
	if !found {
		t.Fatal("expected interleaved-thinking token to be present")
placeholder
placeholder
