//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// --- billingModelForRestriction ---

func TestBillingModelForRestriction_Requested(t *testing.T) {
	t.Parallel()
	got := billingModelForRestriction(BillingModelSourceRequested, "claude-sonnet-4-5", "claude-sonnet-4-6")
	require.Equal(t, "claude-sonnet-4-5", got)
placeholder

func TestBillingModelForRestriction_ChannelMapped(t *testing.T) {
	t.Parallel()
	got := billingModelForRestriction(BillingModelSourceChannelMapped, "claude-sonnet-4-5", "claude-sonnet-4-6")
	require.Equal(t, "claude-sonnet-4-6", got)
placeholder

func TestBillingModelForRestriction_Upstream(t *testing.T) {
	t.Parallel()
	got := billingModelForRestriction(BillingModelSourceUpstream, "claude-sonnet-4-5", "claude-sonnet-4-6")
	require.Equal(t, "", got, "upstream should return empty (per-account check needed)")
placeholder

func TestBillingModelForRestriction_ResponseModelUsesMappedPrecheck(t *testing.T) {
	t.Parallel()
	got := billingModelForRestriction(BillingModelSourceResponse, "claude-fable-5", "claude-fable-5")
	require.Equal(t, "claude-fable-5", got)
placeholder

func TestBillingModelForRestriction_Empty(t *testing.T) {
	t.Parallel()
	got := billingModelForRestriction("", "claude-sonnet-4-5", "claude-sonnet-4-6")
	require.Equal(t, "claude-sonnet-4-6", got, "empty source defaults to channel_mapped")
placeholder

// --- resolveAccountUpstreamModel ---

func TestResolveAccountUpstreamModel_Antigravity(t *testing.T) {
	t.Parallel()
	account := &Account{
		Platform: PlatformAntigravity,
placeholder
	// Antigravity 平台使用 DefaultAntigravityModelMapping
	got := resolveAccountUpstreamModel(account, "claude-sonnet-4-6")
	require.Equal(t, "claude-sonnet-4-6", got)
placeholder

func TestResolveAccountUpstreamModel_Antigravity_Unsupported(t *testing.T) {
	t.Parallel()
	account := &Account{
		Platform: PlatformAntigravity,
placeholder
	got := resolveAccountUpstreamModel(account, "totally-unknown-model")
	require.Equal(t, "", got, "unsupported model should return empty")
placeholder

func TestResolveAccountUpstreamModel_NonAntigravity(t *testing.T) {
	t.Parallel()
	account := &Account{
		Platform: PlatformAnthropic,
placeholder
	got := resolveAccountUpstreamModel(account, "claude-sonnet-4-6")
	require.Equal(t, "claude-sonnet-4-6", got, "no mapping = passthrough")
placeholder

// --- checkChannelPricingRestriction ---

func TestCheckChannelPricingRestriction_NilGroupID(t *testing.T) {
	t.Parallel()
	svc := &GatewayService{channelService: &ChannelService{placeholderplaceholder
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), nil, "claude-sonnet-4"))
placeholder

func TestCheckChannelPricingRestriction_NilChannelService(t *testing.T) {
	t.Parallel()
	svc := &GatewayService{placeholder
	gid := int64(10)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "claude-sonnet-4"))
placeholder

func TestCheckChannelPricingRestriction_EmptyModel(t *testing.T) {
	t.Parallel()
	svc := &GatewayService{channelService: &ChannelService{placeholderplaceholder
	gid := int64(10)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, ""))
placeholder

func TestCheckChannelPricingRestriction_ChannelMapped_Restricted(t *testing.T) {
	t.Parallel()
	// 渠道映射 claude-sonnet-4-5 → claude-sonnet-4-6，但定价列表只有 claude-opus-4-6
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		RestrictModels:     true,
		BillingModelSource: BillingModelSourceChannelMapped,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4-6"placeholderplaceholder,
	placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {"claude-sonnet-4-5": "claude-sonnet-4-6"placeholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(10)
	require.True(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "claude-sonnet-4-5"),
		"mapped model claude-sonnet-4-6 is NOT in pricing → restricted")
placeholder

func TestCheckChannelPricingRestriction_ChannelMapped_Allowed(t *testing.T) {
	t.Parallel()
	// 渠道映射 claude-sonnet-4-5 → claude-sonnet-4-6，定价列表包含 claude-sonnet-4-6
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		RestrictModels:     true,
		BillingModelSource: BillingModelSourceChannelMapped,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4-6"placeholderplaceholder,
	placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {"claude-sonnet-4-5": "claude-sonnet-4-6"placeholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(10)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "claude-sonnet-4-5"),
		"mapped model claude-sonnet-4-6 IS in pricing → allowed")
placeholder

func TestCheckChannelPricingRestriction_Requested_Restricted(t *testing.T) {
	t.Parallel()
	// billing_model_source=requested，定价列表有 claude-sonnet-4-6 但请求的是 claude-sonnet-4-5
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		RestrictModels:     true,
		BillingModelSource: BillingModelSourceRequested,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4-6"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(10)
	require.True(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "claude-sonnet-4-5"),
		"requested model claude-sonnet-4-5 is NOT in pricing → restricted")
placeholder

func TestCheckChannelPricingRestriction_Requested_Allowed(t *testing.T) {
	t.Parallel()
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		RestrictModels:     true,
		BillingModelSource: BillingModelSourceRequested,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4-5"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(10)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "claude-sonnet-4-5"),
		"requested model IS in pricing → allowed")
placeholder

func TestCheckChannelPricingRestriction_Upstream_SkipsPreCheck(t *testing.T) {
	t.Parallel()
	// upstream 模式：预检查始终跳过（返回 false），需逐账号检查
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		RestrictModels:     true,
		BillingModelSource: BillingModelSourceUpstream,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4-6"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(10)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "unknown-model"),
		"upstream mode should skip pre-check (per-account check needed)")
placeholder

func TestCheckChannelPricingRestriction_RestrictModelsDisabled(t *testing.T) {
	t.Parallel()
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: false, // 未开启模型限制
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4-6"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(10)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "any-model"),
		"RestrictModels=false → always allowed")
placeholder

func TestCheckChannelPricingRestriction_NoChannel(t *testing.T) {
	t.Parallel()
	// 分组没有关联渠道
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) { return nil, nil placeholder,
placeholder
	channelSvc := newTestChannelService(repo)
	svc := &GatewayService{channelService: channelSvcplaceholder

	gid := int64(999)
	require.False(t, svc.checkChannelPricingRestriction(context.Background(), &gid, "any-model"),
		"no channel for group → allowed")
placeholder

// --- isUpstreamModelRestrictedByChannel ---

func TestIsUpstreamModelRestrictedByChannel_Restricted(t *testing.T) {
	t.Parallel()
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4-6"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	account := &Account{Platform: PlatformAntigravityplaceholder
	// claude-sonnet-4-6 在 DefaultAntigravityModelMapping 中，映射后仍为 claude-sonnet-4-6
	// 但定价列表只有 claude-opus-4-6
	require.True(t, svc.isUpstreamModelRestrictedByChannel(context.Background(), 10, account, "claude-sonnet-4-6"),
		"upstream model claude-sonnet-4-6 NOT in pricing → restricted")
placeholder

func TestIsUpstreamModelRestrictedByChannel_Allowed(t *testing.T) {
	t.Parallel()
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4-6"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	account := &Account{Platform: PlatformAntigravityplaceholder
	require.False(t, svc.isUpstreamModelRestrictedByChannel(context.Background(), 10, account, "claude-sonnet-4-6"),
		"upstream model claude-sonnet-4-6 IS in pricing → allowed")
placeholder

func TestIsUpstreamModelRestrictedByChannel_UnsupportedModel(t *testing.T) {
	t.Parallel()
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4-6"placeholderplaceholder,
	placeholder,
placeholder
	channelSvc := newTestChannelService(makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder))
	svc := &GatewayService{channelService: channelSvcplaceholder

	account := &Account{Platform: PlatformAntigravityplaceholder
	// totally-unknown-model 不在 DefaultAntigravityModelMapping 中 → 映射结果为空
	require.False(t, svc.isUpstreamModelRestrictedByChannel(context.Background(), 10, account, "totally-unknown-model"),
		"unmappable model → upstream model empty → not restricted (account filter handles this)")
placeholder
