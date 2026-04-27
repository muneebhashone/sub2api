package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type openAIFastPolicyRepoStub struct {
	values map[string]string
placeholder

func (s *openAIFastPolicyRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *openAIFastPolicyRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if v, ok := s.values[key]; ok {
		return v, nil
placeholder
	return "", ErrSettingNotFound
placeholder

func (s *openAIFastPolicyRepoStub) Set(ctx context.Context, key, value string) error {
	if s.values == nil {
		s.values = map[string]string{placeholder
placeholder
	s.values[key] = value
	return nil
placeholder

func (s *openAIFastPolicyRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *openAIFastPolicyRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *openAIFastPolicyRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *openAIFastPolicyRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

func newOpenAIGatewayServiceWithSettings(t *testing.T, settings *OpenAIFastPolicySettings) *OpenAIGatewayService {
placeholder
	repo := &openAIFastPolicyRepoStub{values: map[string]string{placeholderplaceholder
	if settings != nil {
		raw, err := json.Marshal(settings)
	placeholder
		repo.values[SettingKeyOpenAIFastPolicySettings] = string(raw)
placeholder
	return &OpenAIGatewayService{
		settingService: NewSettingService(repo, &config.Config{placeholder),
placeholder
placeholder

func TestEvaluateOpenAIFastPolicy_DefaultFiltersAllModelsPriority(t *testing.T) {
	svc := newOpenAIGatewayServiceWithSettings(t, DefaultOpenAIFastPolicySettings())
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	// 默认策略对所有模型生效（whitelist 为空），因为 codex 的 service_tier=fast
	// 是用户级开关，与 model 正交。
	// gpt-5.5 + priority → filter
	action, _ := svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-5.5", OpenAIFastTierPriority)
	require.Equal(t, BetaPolicyActionFilter, action)

	// gpt-5.5-turbo → filter
	action, _ = svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-5.5-turbo", OpenAIFastTierPriority)
	require.Equal(t, BetaPolicyActionFilter, action)

	// gpt-4 + priority → filter（默认策略覆盖所有模型）
	action, _ = svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-4", OpenAIFastTierPriority)
	require.Equal(t, BetaPolicyActionFilter, action)

	// gpt-5.5 + flex → pass (tier doesn't match)
	action, _ = svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-5.5", OpenAIFastTierFlex)
	require.Equal(t, BetaPolicyActionPass, action)

	// empty tier → pass
	action, _ = svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-5.5", "")
	require.Equal(t, BetaPolicyActionPass, action)
placeholder

func TestEvaluateOpenAIFastPolicy_BlockRuleCarriesMessage(t *testing.T) {
	settings := &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier:    OpenAIFastTierPriority,
			Action:         BetaPolicyActionBlock,
			Scope:          BetaPolicyScopeAll,
			ErrorMessage:   "fast mode is not allowed",
			ModelWhitelist: []string{"gpt-5.5"placeholder,
			FallbackAction: BetaPolicyActionPass,
placeholder
placeholder
	svc := newOpenAIGatewayServiceWithSettings(t, settings)
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	action, msg := svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-5.5", OpenAIFastTierPriority)
	require.Equal(t, BetaPolicyActionBlock, action)
	require.Equal(t, "fast mode is not allowed", msg)
placeholder

func TestEvaluateOpenAIFastPolicy_ScopeFiltersOAuth(t *testing.T) {
	settings := &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier: OpenAIFastTierAny,
			Action:      BetaPolicyActionFilter,
			Scope:       BetaPolicyScopeOAuth,
placeholder
placeholder
	svc := newOpenAIGatewayServiceWithSettings(t, settings)

	// OAuth account → rule matches
	oauthAccount := &Account{Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	action, _ := svc.evaluateOpenAIFastPolicy(context.Background(), oauthAccount, "gpt-4", OpenAIFastTierPriority)
	require.Equal(t, BetaPolicyActionFilter, action)

	// API Key account → rule skipped → pass
	apiKeyAccount := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder
	action, _ = svc.evaluateOpenAIFastPolicy(context.Background(), apiKeyAccount, "gpt-4", OpenAIFastTierPriority)
	require.Equal(t, BetaPolicyActionPass, action)
placeholder

func TestApplyOpenAIFastPolicyToBody_FilterRemovesField(t *testing.T) {
	svc := newOpenAIGatewayServiceWithSettings(t, DefaultOpenAIFastPolicySettings())
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	// gpt-5.5 fast → service_tier stripped
	body := []byte(`{"model":"gpt-5.5","service_tier":"priority","messages":[]placeholder`)
	updated, err := svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
placeholder
	require.NotContains(t, string(updated), `"service_tier"`)

	// Client sending "fast" (alias for priority) also filtered
	body = []byte(`{"model":"gpt-5.5","service_tier":"fast"placeholder`)
	updated, err = svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
placeholder
	require.NotContains(t, string(updated), `"service_tier"`)

	// gpt-4 priority → 默认策略对所有模型 filter，service_tier 被移除
	body = []byte(`{"model":"gpt-4","service_tier":"priority"placeholder`)
	updated, err = svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-4", body)
placeholder
	require.NotContains(t, string(updated), `"service_tier"`)

	// No service_tier → no-op
	body = []byte(`{"model":"gpt-5.5"placeholder`)
	updated, err = svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
placeholder
	require.Equal(t, string(body), string(updated))
placeholder

// TestApplyOpenAIFastPolicyToBody_OfficialTiersBypassDefaultRule 验证扩展白名单后
// 客户端显式发送的 OpenAI 官方合法 tier（auto/default/scale）能透传到上游而不被
// 静默剥离。默认策略只针对 priority，所以这些 tier 落在 fall-through pass 分支。
func TestApplyOpenAIFastPolicyToBody_OfficialTiersBypassDefaultRule(t *testing.T) {
	svc := newOpenAIGatewayServiceWithSettings(t, DefaultOpenAIFastPolicySettings())
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	for _, tier := range []string{"auto", "default", "scale"placeholder {
		body := []byte(`{"model":"gpt-5.5","service_tier":"` + tier + `"placeholder`)
		updated, err := svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
		require.NoError(t, err, "tier %q should pass without error", tier)
		require.Contains(t, string(updated), `"service_tier":"`+tier+`"`,
			"tier %q should be preserved in body under default rule", tier)
placeholder

	// evaluate 层也应判定为 pass（默认规则 ServiceTier=priority 与 auto/default/scale 不匹配）
	for _, tier := range []string{"auto", "default", "scale"placeholder {
		action, _ := svc.evaluateOpenAIFastPolicy(context.Background(), account, "gpt-5.5", tier)
		require.Equal(t, BetaPolicyActionPass, action, "tier %q should evaluate to pass", tier)
placeholder
placeholder

// TestApplyOpenAIFastPolicyToBody_AllRuleStripsOfficialTiers 验证管理员显式配置
// ServiceTier=all + Action=filter 规则后，auto/default/scale 等官方 tier 也会
// 被剥离。这是符合预期的——首条匹配 short-circuit，"all" 覆盖任意已识别 tier。
func TestApplyOpenAIFastPolicyToBody_AllRuleStripsOfficialTiers(t *testing.T) {
	settings := &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier: OpenAIFastTierAny,
			Action:      BetaPolicyActionFilter,
			Scope:       BetaPolicyScopeAll,
placeholder
placeholder
	svc := newOpenAIGatewayServiceWithSettings(t, settings)
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	for _, tier := range []string{"auto", "default", "scale", "priority", "flex"placeholder {
		body := []byte(`{"model":"gpt-5.5","service_tier":"` + tier + `"placeholder`)
		updated, err := svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
	placeholder
		require.NotContains(t, string(updated), `"service_tier"`,
			"tier %q should be stripped under ServiceTier=all + filter rule", tier)
placeholder
placeholder

// TestApplyOpenAIFastPolicyToBody_UnknownTierStripped 验证真未知 tier 仍被剥离
// （normalize 返回 nil → normalizeResponsesBodyServiceTier 删除字段；
// applyOpenAIFastPolicyToBody 在 normTier 为空时直接 no-op，因为字段已不可能存在
// 于经过前置归一化的请求里。这里直接调 apply 验证它对未识别值不会异常）。
func TestApplyOpenAIFastPolicyToBody_UnknownTierStripped(t *testing.T) {
	svc := newOpenAIGatewayServiceWithSettings(t, DefaultOpenAIFastPolicySettings())
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	// normalize 阶段会将未知值剥离
	require.Nil(t, normalizeOpenAIServiceTier("xxx"))

	// applyOpenAIFastPolicyToBody 收到未识别 tier 时不报错，body 透传不变
	// （不属于本函数职责——上层 normalizeResponsesBodyServiceTier 已剥离）
	body := []byte(`{"model":"gpt-5.5","service_tier":"xxx"placeholder`)
	updated, err := svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
placeholder
	require.Equal(t, string(body), string(updated))
placeholder

func TestApplyOpenAIFastPolicyToBody_BlockReturnsTypedError(t *testing.T) {
	settings := &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier:    OpenAIFastTierPriority,
			Action:         BetaPolicyActionBlock,
			Scope:          BetaPolicyScopeAll,
			ErrorMessage:   "fast mode is blocked for gpt-5.5",
			ModelWhitelist: []string{"gpt-5.5"placeholder,
			FallbackAction: BetaPolicyActionPass,
placeholder
placeholder
	svc := newOpenAIGatewayServiceWithSettings(t, settings)
	account := &Account{Platform: PlatformOpenAI, Type: AccountTypeAPIKeyplaceholder

	body := []byte(`{"model":"gpt-5.5","service_tier":"priority"placeholder`)
	updated, err := svc.applyOpenAIFastPolicyToBody(context.Background(), account, "gpt-5.5", body)
placeholder
	var blocked *OpenAIFastBlockedError
	require.True(t, errors.As(err, &blocked))
	require.Contains(t, blocked.Message, "fast mode is blocked")
	require.Equal(t, string(body), string(updated)) // body not mutated on block
placeholder

func TestSetOpenAIFastPolicySettings_Validation(t *testing.T) {
	repo := &openAIFastPolicyRepoStub{values: map[string]string{placeholderplaceholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	// Invalid action rejected
	err := svc.SetOpenAIFastPolicySettings(context.Background(), &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier: OpenAIFastTierPriority,
			Action:      "bogus",
			Scope:       BetaPolicyScopeAll,
placeholder
placeholder)
placeholder

	// Invalid service_tier rejected
	err = svc.SetOpenAIFastPolicySettings(context.Background(), &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier: "turbo",
			Action:      BetaPolicyActionPass,
			Scope:       BetaPolicyScopeAll,
placeholder
placeholder)
placeholder

	// Valid settings persisted
	err = svc.SetOpenAIFastPolicySettings(context.Background(), &OpenAIFastPolicySettings{
		Rules: []OpenAIFastPolicyRule{{
			ServiceTier: OpenAIFastTierPriority,
			Action:      BetaPolicyActionFilter,
			Scope:       BetaPolicyScopeAll,
placeholder
placeholder)
placeholder

	got, err := svc.GetOpenAIFastPolicySettings(context.Background())
placeholder
	require.Len(t, got.Rules, 1)
	require.Equal(t, OpenAIFastTierPriority, got.Rules[0].ServiceTier)
placeholder
