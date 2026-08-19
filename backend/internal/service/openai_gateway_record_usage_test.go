package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/ctxkey"
	"github.com/stretchr/testify/require"
)

type openAIRecordUsageLogRepoStub struct {
	UsageLogRepository

	inserted   bool
	err        error
	calls      int
	lastLog    *UsageLog
	lastCtxErr error
placeholder

func (s *openAIRecordUsageLogRepoStub) Create(ctx context.Context, log *UsageLog) (bool, error) {
	s.calls++
	s.lastLog = log
	s.lastCtxErr = ctx.Err()
	return s.inserted, s.err
placeholder

type openAIRecordUsageBillingRepoStub struct {
	UsageBillingRepository

	result     *UsageBillingApplyResult
	err        error
	calls      int
	lastCmd    *UsageBillingCommand
	lastCtxErr error
placeholder

type openAIRecordUsageAccountRepoStub struct {
	AccountRepository
	account *Account
	calls   int
placeholder

func (s *openAIRecordUsageAccountRepoStub) GetByID(_ context.Context, _ int64) (*Account, error) {
	s.calls++
	return s.account, nil
placeholder

func (s *openAIRecordUsageBillingRepoStub) Apply(ctx context.Context, cmd *UsageBillingCommand) (*UsageBillingApplyResult, error) {
	s.calls++
	s.lastCmd = cmd
	s.lastCtxErr = ctx.Err()
	if s.err != nil {
		return nil, s.err
placeholder
	if s.result != nil {
		return s.result, nil
placeholder
	return &UsageBillingApplyResult{Applied: trueplaceholder, nil
placeholder

func TestOpenAIGatewayServiceRecordUsage_RejectsNilInput(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	require.Error(t, svc.RecordUsage(context.Background(), nil))
	require.Error(t, svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{placeholder))
placeholder

func TestRecordCyberPolicyUsageLog_BillsRealUpstreamTokens(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	usage := OpenAIUsage{InputTokens: 1200, OutputTokens: 300placeholder

	// 流式 cyber：上游 response.failed 报告了真实 token，须按真实 token 计费并扣费，
	// 与 WS cyber / 正常请求口径一致（不再是 tokens=0 免费行）。
	svc.RecordCyberPolicyUsageLog(context.Background(), CyberPolicyUsageInput{
		APIKey:       &APIKey{ID: 2, User: &User{ID: 1placeholderplaceholder,
		Account:      &Account{ID: 3placeholder,
		RequestID:    "rid-cyber-stream",
		Model:        "gpt-5.1",
		Stream:       true,
		InputTokens:  1200,
		OutputTokens: 300,
placeholder)

	require.Equal(t, 1, usageRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "gpt-5.1", usageRepo.lastLog.Model)
	require.Equal(t, 1200, usageRepo.lastLog.InputTokens)
	require.Equal(t, 300, usageRepo.lastLog.OutputTokens)
	require.Equal(t, RequestTypeCyberBlocked, usageRepo.lastLog.RequestType, "cyber 行须标 request_type=cyber")
	require.True(t, usageRepo.lastLog.Stream, "cyber 不覆盖真实 stream 字段")

	expected := expectedOpenAICost(t, svc, "gpt-5.1", usage, 1.1)
	require.Greater(t, usageRepo.lastLog.ActualCost, 0.0, "流式 cyber 有真实 token，须计费")
	require.InDelta(t, expected.ActualCost, usageRepo.lastLog.ActualCost, 1e-12)
	require.Equal(t, 1, userRepo.deductCalls, "按真实 token 扣费，与 WS/正常请求一致")
	require.InDelta(t, expected.ActualCost, userRepo.lastAmount, 1e-12)
placeholder

func TestRecordCyberPolicyUsageLog_NonStreamZeroTokensZeroCost(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	// 非流式直接拒：上游未报 token，mark token 为 0 → cost 自然为 0，仍写一条 cyber 行（可见）。
	svc.RecordCyberPolicyUsageLog(context.Background(), CyberPolicyUsageInput{
		APIKey:    &APIKey{ID: 2, User: &User{ID: 1placeholderplaceholder,
		Account:   &Account{ID: 3placeholder,
		RequestID: "rid-cyber-400",
		Model:     "gpt-5.1",
		Stream:    false,
placeholder)

	require.Equal(t, 1, usageRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 0, usageRepo.lastLog.InputTokens)
	require.Equal(t, 0, usageRepo.lastLog.OutputTokens)
	require.Zero(t, usageRepo.lastLog.TotalCost)
	require.Equal(t, RequestTypeCyberBlocked, usageRepo.lastLog.RequestType)
placeholder

func TestRecordCyberPolicyUsageLog_SkipsWhenIncomplete(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	acct := &Account{ID: 3placeholder
	svc.RecordCyberPolicyUsageLog(context.Background(), CyberPolicyUsageInput{Account: acct, Model: "gpt-5"placeholder)                              // APIKey nil
	svc.RecordCyberPolicyUsageLog(context.Background(), CyberPolicyUsageInput{APIKey: &APIKey{ID: 2placeholder, Account: acct, Model: "gpt-5"placeholder)      // User nil
	svc.RecordCyberPolicyUsageLog(context.Background(), CyberPolicyUsageInput{APIKey: &APIKey{ID: 2, User: &User{ID: 1placeholderplaceholder, Model: "gpt-5"placeholder) // Account nil
	svc.RecordCyberPolicyUsageLog(context.Background(), CyberPolicyUsageInput{APIKey: &APIKey{ID: 2, User: &User{ID: 1placeholderplaceholder, Account: acctplaceholder)  // Model 空
	require.Equal(t, 0, usageRepo.calls, "APIKey/User/Account 缺失或 Model 空时跳过，不记不扣费")
placeholder

type openAIRecordUsageUserRepoStub struct {
	UserRepository

	deductCalls int
	deductErr   error
	lastAmount  float64
	lastCtxErr  error
placeholder

func (s *openAIRecordUsageUserRepoStub) DeductBalance(ctx context.Context, id int64, amount float64) error {
	s.deductCalls++
	s.lastAmount = amount
	s.lastCtxErr = ctx.Err()
	return s.deductErr
placeholder

func (s *openAIRecordUsageUserRepoStub) AdjustBalance(ctx context.Context, id int64, delta float64) (BalanceChange, error) {
	panic("unexpected AdjustBalance call")
placeholder

func (s *openAIRecordUsageUserRepoStub) SetBalance(ctx context.Context, id int64, value float64) (BalanceChange, error) {
	panic("unexpected SetBalance call")
placeholder

type openAIRecordUsageSubRepoStub struct {
	UserSubscriptionRepository

	incrementCalls int
	incrementErr   error
	lastCtxErr     error
placeholder

func (s *openAIRecordUsageSubRepoStub) IncrementUsage(ctx context.Context, id int64, costUSD float64) error {
	s.incrementCalls++
	s.lastCtxErr = ctx.Err()
	return s.incrementErr
placeholder

type openAIRecordUsageAPIKeyQuotaStub struct {
	quotaCalls          int
	rateLimitCalls      int
	err                 error
	lastAmount          float64
	lastQuotaCtxErr     error
	lastRateLimitCtxErr error
placeholder

func (s *openAIRecordUsageAPIKeyQuotaStub) UpdateQuotaUsed(ctx context.Context, apiKeyID int64, cost float64) error {
	s.quotaCalls++
	s.lastAmount = cost
	s.lastQuotaCtxErr = ctx.Err()
	return s.err
placeholder

func (s *openAIRecordUsageAPIKeyQuotaStub) UpdateRateLimitUsage(ctx context.Context, apiKeyID int64, cost float64) error {
	s.rateLimitCalls++
	s.lastAmount = cost
	s.lastRateLimitCtxErr = ctx.Err()
	return s.err
placeholder

type openAIUserGroupRateRepoStub struct {
	UserGroupRateRepository

	rate  *float64
	err   error
	calls int
placeholder

func (s *openAIUserGroupRateRepoStub) GetByUserAndGroup(ctx context.Context, userID, groupID int64) (*float64, error) {
	s.calls++
	if s.err != nil {
		return nil, s.err
placeholder
	return s.rate, nil
placeholder

func i64p(v int64) *int64 {
	return &v
placeholder

func newOpenAIRecordUsageServiceForTest(usageRepo UsageLogRepository, userRepo UserRepository, subRepo UserSubscriptionRepository, rateRepo UserGroupRateRepository) *OpenAIGatewayService {
	cfg := &config.Config{placeholder
	cfg.Default.RateMultiplier = 1.1
	svc := NewOpenAIGatewayService(
		nil,
		usageRepo,
		nil,
		userRepo,
		subRepo,
		rateRepo,
		nil,
		cfg,
		nil,
		nil,
		NewBillingService(cfg, nil),
		nil,
		&BillingCacheService{placeholder,
		nil,
		&DeferredService{placeholder,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil, // userPlatformQuotaRepo
	)
	svc.userGroupRateResolver = newUserGroupRateResolver(
		rateRepo,
		nil,
		resolveUserGroupRateCacheTTL(cfg),
		nil,
		"service.openai_gateway.test",
	)
	return svc
placeholder

func openAIRecordUsageAPIKeyWithGroup(svc *OpenAIGatewayService, id int64, groupLongContext bool) *APIKey {
	svc.resolver = NewModelPricingResolver(nil, svc.billingService)
	return &APIKey{
		ID: id,
		Group: &Group{
			ID:                        1,
			LongContextPricingEnabled: groupLongContext,
	placeholder,
placeholder
placeholder

func newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo UsageLogRepository, billingRepo UsageBillingRepository, userRepo UserRepository, subRepo UserSubscriptionRepository, rateRepo UserGroupRateRepository) *OpenAIGatewayService {
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, rateRepo)
	svc.usageBillingRepo = billingRepo
	return svc
placeholder

func expectedOpenAICost(t *testing.T, svc *OpenAIGatewayService, model string, usage OpenAIUsage, multiplier float64) *CostBreakdown {
placeholder

	cost, err := svc.billingService.CalculateCost(model, UsageTokens{
		InputTokens:         max(usage.InputTokens-usage.CacheReadInputTokens-usage.CacheCreationInputTokens, 0),
		OutputTokens:        usage.OutputTokens,
		CacheCreationTokens: usage.CacheCreationInputTokens,
		CacheReadTokens:     usage.CacheReadInputTokens,
placeholder, multiplier)
placeholder
	return cost
placeholder

func max(a, b int) int {
	if a > b {
		return a
placeholder
	return b
placeholder

func TestOpenAIGatewayServiceRecordUsage_ZeroUsageStillWritesUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_zero_usage",
			Usage:     OpenAIUsage{placeholder,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey:        &APIKey{ID: 1000, Quota: 100, Group: &Group{RateMultiplier: 1placeholderplaceholder,
		User:          &User{ID: 2000placeholder,
		Account:       &Account{ID: 3000, Type: AccountTypeAPIKeyplaceholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 0, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
	require.Equal(t, 0, quotaSvc.quotaCalls)
	require.Equal(t, 0, quotaSvc.rateLimitCalls)

	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "resp_zero_usage", usageRepo.lastLog.RequestID)
	require.Zero(t, usageRepo.lastLog.InputTokens)
	require.Zero(t, usageRepo.lastLog.OutputTokens)
	require.Zero(t, usageRepo.lastLog.CacheCreationTokens)
	require.Zero(t, usageRepo.lastLog.CacheReadTokens)
	require.Zero(t, usageRepo.lastLog.ImageOutputTokens)
	require.Zero(t, usageRepo.lastLog.ImageCount)
	require.Zero(t, usageRepo.lastLog.InputCost)
	require.Zero(t, usageRepo.lastLog.OutputCost)
	require.Zero(t, usageRepo.lastLog.TotalCost)
	require.Zero(t, usageRepo.lastLog.ActualCost)

	require.NotNil(t, billingRepo.lastCmd)
	require.Zero(t, billingRepo.lastCmd.BalanceCost)
	require.Zero(t, billingRepo.lastCmd.SubscriptionCost)
	require.Zero(t, billingRepo.lastCmd.APIKeyQuotaCost)
	require.Zero(t, billingRepo.lastCmd.APIKeyRateLimitCost)
	require.Zero(t, billingRepo.lastCmd.AccountQuotaCost)
placeholder

func TestOpenAIGatewayServiceRecordUsage_MissingPricingRecordsZeroCostUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_missing_pricing",
			Usage: OpenAIUsage{
				InputTokens:  1200,
				OutputTokens: 300,
		placeholder,
			Model:    "pricing-missing-test-model",
			Duration: time.Second,
	placeholder,
		APIKey:        &APIKey{ID: 1002, Quota: 100, Group: &Group{RateMultiplier: 1placeholderplaceholder,
		User:          &User{ID: 2002placeholder,
		Account:       &Account{ID: 3002, Type: AccountTypeAPIKeyplaceholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 0, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
	require.Equal(t, 0, quotaSvc.quotaCalls)
	require.Equal(t, 0, quotaSvc.rateLimitCalls)

	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "resp_missing_pricing", usageRepo.lastLog.RequestID)
	require.Equal(t, "pricing-missing-test-model", usageRepo.lastLog.Model)
	require.Equal(t, "pricing-missing-test-model", usageRepo.lastLog.RequestedModel)
	require.Equal(t, 1200, usageRepo.lastLog.InputTokens)
	require.Equal(t, 300, usageRepo.lastLog.OutputTokens)
	require.Zero(t, usageRepo.lastLog.TotalCost)
	require.Zero(t, usageRepo.lastLog.ActualCost)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeToken), *usageRepo.lastLog.BillingMode)

	require.NotNil(t, billingRepo.lastCmd)
	require.Zero(t, billingRepo.lastCmd.BalanceCost)
	require.Zero(t, billingRepo.lastCmd.SubscriptionCost)
	require.Zero(t, billingRepo.lastCmd.APIKeyQuotaCost)
	require.Zero(t, billingRepo.lastCmd.APIKeyRateLimitCost)
	require.Zero(t, billingRepo.lastCmd.AccountQuotaCost)
placeholder

func TestOpenAIGatewayServiceRecordUsage_UsesUserSpecificGroupRate(t *testing.T) {
	groupID := int64(11)
	groupRate := 1.4
	userRate := 1.8
	usage := OpenAIUsage{InputTokens: 15, OutputTokens: 4, CacheReadInputTokens: 3placeholder

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	rateRepo := &openAIUserGroupRateRepoStub{rate: &userRateplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, rateRepo)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_user_group_rate",
			Usage:     usage,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      1001,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: groupRate,
		placeholder,
	placeholder,
		User:    &User{ID: 2001placeholder,
		Account: &Account{ID: 3001placeholder,
placeholder)

placeholder
	require.Equal(t, 1, rateRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, userRate, usageRepo.lastLog.RateMultiplier)
	require.Equal(t, 12, usageRepo.lastLog.InputTokens)
	require.Equal(t, 3, usageRepo.lastLog.CacheReadTokens)

	expected := expectedOpenAICost(t, svc, "gpt-5.1", usage, userRate)
	require.InDelta(t, expected.ActualCost, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, expected.ActualCost, userRepo.lastAmount, 1e-12)
	require.Equal(t, 1, userRepo.deductCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_PeakRateAffectsTokenModeImageOutputTokens(t *testing.T) {
	groupID := int64(14)
	groupRate := 1.0
	usage := OpenAIUsage{
		InputTokens:       1000,
		OutputTokens:      600,
		ImageOutputTokens: 100,
placeholder

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	svc.resolver = newOpenAITokenImageChannelPricingResolverForTest(t, groupID, "gpt-5.1")

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_peak_image_tokens",
			Usage:      usage,
			Model:      "gpt-5.1",
			Duration:   time.Second,
			ImageCount: 1,
	placeholder,
		APIKey: &APIKey{
			ID:      1004,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                 groupID,
				RateMultiplier:     groupRate,
				SubscriptionType:   "subscription",
				PeakRateEnabled:    true,
				PeakStart:          "00:00",
				PeakEnd:            "23:59",
				PeakRateMultiplier: 3.0,
		placeholder,
	placeholder,
		User:    &User{ID: 2004placeholder,
		Account: &Account{ID: 3004placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 3.0, usageRepo.lastLog.RateMultiplier)
	require.Equal(t, usage.ImageOutputTokens, usageRepo.lastLog.ImageOutputTokens)

	expected, err := svc.billingService.CalculateCostUnified(CostInput{
		Ctx:     context.Background(),
		Model:   "gpt-5.1",
		GroupID: i64p(groupID),
		Tokens: UsageTokens{
			InputTokens:       usage.InputTokens,
			OutputTokens:      usage.OutputTokens,
			ImageOutputTokens: usage.ImageOutputTokens,
	placeholder,
		RateMultiplier: 1.0,
		Resolver:       svc.resolver,
placeholder)
placeholder
	expectedActual := expected.TotalCost * 3.0

	require.InDelta(t, expected.TotalCost, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, expected.ImageOutputCost, usageRepo.lastLog.ImageOutputCost, 1e-12)
	require.InDelta(t, expectedActual, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, expectedActual, userRepo.lastAmount, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_TimePricingUsesPricingAt(t *testing.T) {
	groupID := int64(16)
	requestStart := time.Date(2024, time.January, 2, 2, 0, 0, 0, time.UTC) // 上海 10:00
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAITokenImageChannelPricingResolverWithTimeForTest(t, groupID, "gpt-5.1", &ChannelTimePricing{
		Timezone: "Asia/Shanghai",
		Periods:  []ChannelTimePricingPeriod{{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholderplaceholder,
placeholder)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "openai_time_pricing_request_start",
			Model:     "gpt-5.1",
			Usage:     OpenAIUsage{InputTokens: 1000, OutputTokens: 500placeholder,
	placeholder,
		APIKey: &APIKey{ID: 1006, GroupID: i64p(groupID), Group: &Group{
			ID: groupID, RateMultiplier: 0.8, SubscriptionType: SubscriptionTypeSubscription,
placeholder
		User:      &User{ID: 2006placeholder,
		Account:   &Account{ID: 3006placeholder,
		PricingAt: requestStart,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	baseCost := 1000*3e-6 + 500*15e-6
	require.InDelta(t, baseCost*2, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, baseCost*2*0.8, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.8, usageRepo.lastLog.RateMultiplier, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_TimePricingUsesExplicitPricingAt(t *testing.T) {
	groupID := int64(17)
	pricingAt := time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC) // 上海 08:00
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAITokenImageChannelPricingResolverWithTimeForTest(t, groupID, "gpt-5.1", &ChannelTimePricing{
		Timezone: "Asia/Shanghai",
		Periods:  []ChannelTimePricingPeriod{{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholderplaceholder,
placeholder)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "openai_time_pricing_explicit",
			Model:     "gpt-5.1",
			Usage:     OpenAIUsage{InputTokens: 1000, OutputTokens: 500placeholder,
	placeholder,
		APIKey: &APIKey{ID: 1007, GroupID: i64p(groupID), Group: &Group{
			ID: groupID, RateMultiplier: 0.8, SubscriptionType: SubscriptionTypeSubscription,
placeholder
		User:      &User{ID: 2007placeholder,
		Account:   &Account{ID: 3007placeholder,
		PricingAt: pricingAt,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	baseCost := 1000*3e-6 + 500*15e-6
	require.InDelta(t, baseCost, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, baseCost*0.8, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.8, usageRepo.lastLog.RateMultiplier, 1e-12)
placeholder
func TestOpenAIGatewayServiceRecordUsage_IncludesEndpointMetadata(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	rateRepo := &openAIUserGroupRateRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, rateRepo)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_endpoint_metadata",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 2,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    1002,
			Group: &Group{RateMultiplier: 1placeholder,
	placeholder,
		User:             &User{ID: 2002placeholder,
		Account:          &Account{ID: 3002placeholder,
		InboundEndpoint:  " /v1/chat/completions ",
		UpstreamEndpoint: " /v1/responses ",
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.InboundEndpoint)
	require.Equal(t, "/v1/chat/completions", *usageRepo.lastLog.InboundEndpoint)
	require.NotNil(t, usageRepo.lastLog.UpstreamEndpoint)
	require.Equal(t, "/v1/responses", *usageRepo.lastLog.UpstreamEndpoint)
placeholder

func TestOpenAIGatewayServiceRecordUsage_FallsBackToGroupDefaultRateOnResolverError(t *testing.T) {
	groupID := int64(12)
	groupRate := 1.6
	usage := OpenAIUsage{InputTokens: 10, OutputTokens: 5, CacheReadInputTokens: 2placeholder

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	rateRepo := &openAIUserGroupRateRepoStub{err: errors.New("db unavailable")placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, rateRepo)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_group_default_on_error",
			Usage:     usage,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      1002,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: groupRate,
		placeholder,
	placeholder,
		User:    &User{ID: 2002placeholder,
		Account: &Account{ID: 3002placeholder,
placeholder)

placeholder
	require.Equal(t, 1, rateRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, groupRate, usageRepo.lastLog.RateMultiplier)

	expected := expectedOpenAICost(t, svc, "gpt-5.1", usage, groupRate)
	require.InDelta(t, expected.ActualCost, userRepo.lastAmount, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_FallsBackToGroupDefaultRateWhenResolverMissing(t *testing.T) {
	groupID := int64(13)
	groupRate := 1.25
	usage := OpenAIUsage{InputTokens: 9, OutputTokens: 4, CacheReadInputTokens: 1placeholder

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	svc.userGroupRateResolver = nil

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_group_default_nil_resolver",
			Usage:     usage,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      1003,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: groupRate,
		placeholder,
	placeholder,
		User:    &User{ID: 2003placeholder,
		Account: &Account{ID: 3003placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, groupRate, usageRepo.lastLog.RateMultiplier)
placeholder

func TestOpenAIGatewayServiceRecordUsage_DuplicateUsageLogSkipsBilling(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: falseplaceholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: falseplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_duplicate",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1004placeholder,
		User:    &User{ID: 2004placeholder,
		Account: &Account{ID: 3004placeholder,
placeholder)

placeholder
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 0, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_DuplicateBillingKeySkipsBillingWithRepo(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: falseplaceholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: falseplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_duplicate_billing_key",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    10045,
			Quota: 100,
	placeholder,
		User:          &User{ID: 20045placeholder,
		Account:       &Account{ID: 30045placeholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 0, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
	require.Equal(t, 0, quotaSvc.quotaCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillsWhenUsageLogCreateReturnsError(t *testing.T) {
	usage := OpenAIUsage{InputTokens: 8, OutputTokens: 4placeholder
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: false, err: errors.New("usage log batch state uncertain")placeholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_usage_log_error",
			Usage:     usage,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10041placeholder,
		User:    &User{ID: 20041placeholder,
		Account: &Account{ID: 30041placeholder,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 1, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_UsageLogWriteErrorDoesNotSkipBilling(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: false, err: MarkUsageLogCreateNotPersisted(context.Canceled)placeholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_not_persisted",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    10043,
			Quota: 100,
	placeholder,
		User:          &User{ID: 20043placeholder,
		Account:       &Account{ID: 30043placeholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 1, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
	require.Equal(t, 1, quotaSvc.quotaCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillingUsesDetachedContext(t *testing.T) {
	usage := OpenAIUsage{InputTokens: 10, OutputTokens: 6, CacheReadInputTokens: 2placeholder
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: false, err: context.DeadlineExceededplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	reqCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := svc.RecordUsage(reqCtx, &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_detached_billing_ctx",
			Usage:     usage,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    10042,
			Quota: 100,
	placeholder,
		User:          &User{ID: 20042placeholder,
		Account:       &Account{ID: 30042placeholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, userRepo.deductCalls)
	require.NoError(t, userRepo.lastCtxErr)
	require.Equal(t, 1, quotaSvc.quotaCalls)
	require.NoError(t, quotaSvc.lastQuotaCtxErr)
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillingRepoUsesDetachedContext(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	reqCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := svc.RecordUsage(reqCtx, &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_detached_billing_repo_ctx",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10046placeholder,
		User:    &User{ID: 20046placeholder,
		Account: &Account{ID: 30046placeholder,
placeholder)

placeholder
	require.Equal(t, 1, billingRepo.calls)
	require.NoError(t, billingRepo.lastCtxErr)
	require.Equal(t, 1, usageRepo.calls)
	require.NoError(t, usageRepo.lastCtxErr)
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillingFingerprintIncludesRequestPayloadHash(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	payloadHash := HashUsageRequestPayload([]byte(`{"model":"gpt-5","input":"hello"placeholder`))
	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "openai_payload_hash",
			Usage: OpenAIUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "gpt-5",
			Duration: time.Second,
	placeholder,
		APIKey:             &APIKey{ID: 501, Quota: 100placeholder,
		User:               &User{ID: 601placeholder,
		Account:            &Account{ID: 701placeholder,
		RequestPayloadHash: payloadHash,
placeholder)
placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.Equal(t, payloadHash, billingRepo.lastCmd.RequestPayloadHash)
placeholder

func TestOpenAIGatewayServiceRecordUsage_UsesFallbackRequestIDForBillingAndUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	ctx := context.WithValue(context.Background(), ctxkey.RequestID, "req-local-fallback")
	err := svc.RecordUsage(ctx, &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10047placeholder,
		User:    &User{ID: 20047placeholder,
		Account: &Account{ID: 30047placeholder,
placeholder)

placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.Equal(t, "local:req-local-fallback", billingRepo.lastCmd.RequestID)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "local:req-local-fallback", usageRepo.lastLog.RequestID)
placeholder

func TestOpenAIGatewayServiceRecordUsage_PrefersClientRequestIDOverUpstreamRequestID(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	ctx := context.WithValue(context.Background(), ctxkey.ClientRequestID, "openai-client-stable-123")
	err := svc.RecordUsage(ctx, &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "upstream-openai-volatile-456",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10049placeholder,
		User:    &User{ID: 20049placeholder,
		Account: &Account{ID: 30049placeholder,
placeholder)

placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.Equal(t, "client:openai-client-stable-123", billingRepo.lastCmd.RequestID)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "client:openai-client-stable-123", usageRepo.lastLog.RequestID)
placeholder

func TestOpenAIGatewayServiceRecordUsage_WSModePrefersUpstreamRequestIDOverClientRequestID(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	ctx := context.WithValue(context.Background(), ctxkey.ClientRequestID, "openai-ws-connection-123")
	err := svc.RecordUsage(ctx, &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:    "resp_openai_ws_turn_456",
			OpenAIWSMode: true,
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10050placeholder,
		User:    &User{ID: 20050placeholder,
		Account: &Account{ID: 30050placeholder,
placeholder)

placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.Equal(t, "resp_openai_ws_turn_456", billingRepo.lastCmd.RequestID)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "resp_openai_ws_turn_456", usageRepo.lastLog.RequestID)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GeneratesRequestIDWhenAllSourcesMissing(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10050placeholder,
		User:    &User{ID: 20050placeholder,
		Account: &Account{ID: 30050placeholder,
placeholder)

placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.True(t, strings.HasPrefix(billingRepo.lastCmd.RequestID, "generated:"))
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, billingRepo.lastCmd.RequestID, usageRepo.lastLog.RequestID)
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillingErrorWritesUnsettledUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingErr := errors.New("billing tx failed")
	billingRepo := &openAIRecordUsageBillingRepoStub{err: billingErrplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_billing_fail",
			Usage: OpenAIUsage{
				InputTokens:  8,
				OutputTokens: 4,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10048placeholder,
		User:    &User{ID: 20048placeholder,
		Account: &Account{ID: 30048placeholder,
placeholder)

	require.ErrorIs(t, err, billingErr)
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 1, usageRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 8, usageRepo.lastLog.InputTokens)
	require.Equal(t, 4, usageRepo.lastLog.OutputTokens)
	require.Greater(t, usageRepo.lastLog.InputCost, 0.0)
	require.Greater(t, usageRepo.lastLog.OutputCost, 0.0)
	require.Greater(t, usageRepo.lastLog.TotalCost, 0.0)
	require.Zero(t, usageRepo.lastLog.ActualCost)
placeholder

func TestOpenAIGatewayServiceRecordUsage_UpdatesAPIKeyQuotaWhenConfigured(t *testing.T) {
	usage := OpenAIUsage{InputTokens: 10, OutputTokens: 6, CacheReadInputTokens: 2placeholder
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_quota_update",
			Usage:     usage,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    1005,
			Quota: 100,
	placeholder,
		User:          &User{ID: 2005placeholder,
		Account:       &Account{ID: 3005placeholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, quotaSvc.quotaCalls)
	require.Equal(t, 0, quotaSvc.rateLimitCalls)
	expected := expectedOpenAICost(t, svc, "gpt-5.1", usage, 1.1)
	require.InDelta(t, expected.ActualCost, quotaSvc.lastAmount, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ClampsActualInputTokensToZero(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_clamp_actual_input",
			Usage: OpenAIUsage{
				InputTokens:          2,
				OutputTokens:         1,
				CacheReadInputTokens: 5,
		placeholder,
			Model:    "gpt-5.1",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1006placeholder,
		User:    &User{ID: 2006placeholder,
		Account: &Account{ID: 3006placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 0, usageRepo.lastLog.InputTokens)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GPT56SeparatesCacheWriteForBillingAndStats(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	svc.billingService = NewBillingService(svc.cfg, &PricingService{pricingData: map[string]*LiteLLMModelPricing{
		"gpt-5.6-sol": {
			InputCostPerToken:       5e-6,
			OutputCostPerToken:      30e-6,
			CacheReadInputTokenCost: 0.5e-6,
	placeholder,
placeholderplaceholder)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_gpt56_cache_write",
			Usage: OpenAIUsage{
				InputTokens:              1000,
				OutputTokens:             50,
				CacheCreationInputTokens: 200,
				CacheReadInputTokens:     100,
		placeholder,
			Model:    "gpt-5.6-sol",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1056placeholder,
		User:    &User{ID: 2056placeholder,
		Account: &Account{ID: 3056placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 700, usageRepo.lastLog.InputTokens)
	require.Equal(t, 200, usageRepo.lastLog.CacheCreationTokens)
	require.Equal(t, 100, usageRepo.lastLog.CacheReadTokens)
	require.Equal(t, 1050, usageRepo.lastLog.TotalTokens())
	require.InDelta(t, 700*5e-6, usageRepo.lastLog.InputCost, 1e-12)
	require.InDelta(t, 200*6.25e-6, usageRepo.lastLog.CacheCreationCost, 1e-12)
	require.InDelta(t, 100*0.5e-6, usageRepo.lastLog.CacheReadCost, 1e-12)
	require.InDelta(t, 50*30e-6, usageRepo.lastLog.OutputCost, 1e-12)
	require.InDelta(t, usageRepo.lastLog.TotalCost*1.1, usageRepo.lastLog.ActualCost, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_Gpt54LongContextBillingDisabledByDefault(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_gpt54_long_context",
			Usage: OpenAIUsage{
				InputTokens:  300000,
				OutputTokens: 2000,
		placeholder,
			Model:    "gpt-5.4-2026-03-05",
			Duration: time.Second,
	placeholder,
		APIKey:  openAIRecordUsageAPIKeyWithGroup(svc, 1014, false),
		User:    &User{ID: 2014placeholder,
		Account: &Account{ID: 3014, Platform: PlatformOpenAIplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)

	expectedInput := 300000 * 2.5e-6
	expectedOutput := 2000 * 15e-6
	require.InDelta(t, expectedInput, usageRepo.lastLog.InputCost, 1e-10)
	require.InDelta(t, expectedOutput, usageRepo.lastLog.OutputCost, 1e-10)
	require.InDelta(t, expectedInput+expectedOutput, usageRepo.lastLog.TotalCost, 1e-10)
	require.InDelta(t, (expectedInput+expectedOutput)*1.1, usageRepo.lastLog.ActualCost, 1e-10)
	require.False(t, usageRepo.lastLog.LongContextBillingApplied)
	require.Equal(t, 1, userRepo.deductCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_Gpt54LongContextBillingEnabledPerAccount(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_gpt54_long_context_disabled",
			Usage: OpenAIUsage{
				InputTokens:  300000,
				OutputTokens: 2000,
		placeholder,
			Model:    "gpt-5.4-2026-03-05",
			Duration: time.Second,
	placeholder,
		APIKey: openAIRecordUsageAPIKeyWithGroup(svc, 1015, true),
		User:   &User{ID: 2015placeholder,
		Account: &Account{
			ID:       3015,
			Platform: PlatformOpenAI,
			Extra:    map[string]any{"openai_long_context_billing_enabled": trueplaceholder,
	placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)

	expectedInput := 300000 * 2.5e-6 * 2.0
	expectedOutput := 2000 * 15e-6 * 1.5
	require.InDelta(t, expectedInput, usageRepo.lastLog.InputCost, 1e-10)
	require.InDelta(t, expectedOutput, usageRepo.lastLog.OutputCost, 1e-10)
	require.InDelta(t, expectedInput+expectedOutput, usageRepo.lastLog.TotalCost, 1e-10)
	require.InDelta(t, (expectedInput+expectedOutput)*1.1, usageRepo.lastLog.ActualCost, 1e-10)
	require.True(t, usageRepo.lastLog.LongContextBillingApplied)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GroupOrAccountLongContextAllows(t *testing.T) {
	tokens := OpenAIUsage{InputTokens: 300000, OutputTokens: 2000placeholder
	baseInput := 300000 * 2.5e-6
	baseOutput := 2000 * 15e-6

	t.Run("group on account off", func(t *testing.T) {
		usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
		svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
		err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
			Result:  &OpenAIForwardResult{RequestID: "resp_and_off", Usage: tokens, Model: "gpt-5.4-2026-03-05", Duration: time.Secondplaceholder,
			APIKey:  openAIRecordUsageAPIKeyWithGroup(svc, 1020, true),
			User:    &User{ID: 2020placeholder,
			Account: &Account{ID: 3020, Platform: PlatformOpenAIplaceholder,
	placeholder)
	placeholder
		require.True(t, usageRepo.lastLog.LongContextBillingApplied)
		require.InDelta(t, baseInput*2, usageRepo.lastLog.InputCost, 1e-10)
		require.InDelta(t, baseOutput*1.5, usageRepo.lastLog.OutputCost, 1e-10)
placeholder)

	t.Run("group off account on", func(t *testing.T) {
		usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
		svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
		err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
			Result: &OpenAIForwardResult{RequestID: "resp_and_group_off", Usage: tokens, Model: "gpt-5.4-2026-03-05", Duration: time.Secondplaceholder,
			APIKey: openAIRecordUsageAPIKeyWithGroup(svc, 1021, false),
			User:   &User{ID: 2021placeholder,
			Account: &Account{
				ID: 3021, Platform: PlatformOpenAI,
				Extra: map[string]any{"openai_long_context_billing_enabled": trueplaceholder,
		placeholder,
	placeholder)
	placeholder
		require.True(t, usageRepo.lastLog.LongContextBillingApplied)
		require.InDelta(t, baseInput*2, usageRepo.lastLog.InputCost, 1e-10)
		require.InDelta(t, baseOutput*1.5, usageRepo.lastLog.OutputCost, 1e-10)
placeholder)

	t.Run("group on account on", func(t *testing.T) {
		usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
		svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
		err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
			Result: &OpenAIForwardResult{RequestID: "resp_and_on", Usage: tokens, Model: "gpt-5.4-2026-03-05", Duration: time.Secondplaceholder,
			APIKey: openAIRecordUsageAPIKeyWithGroup(svc, 1022, true),
			User:   &User{ID: 2022placeholder,
			Account: &Account{
				ID: 3022, Platform: PlatformOpenAI,
				Extra: map[string]any{"openai_long_context_billing_enabled": trueplaceholder,
		placeholder,
	placeholder)
	placeholder
		require.True(t, usageRepo.lastLog.LongContextBillingApplied)
		require.InDelta(t, baseInput*2, usageRepo.lastLog.InputCost, 1e-10)
		require.InDelta(t, baseOutput*1.5, usageRepo.lastLog.OutputCost, 1e-10)
placeholder)
placeholder

// openai_long_context_billing_enabled is an OpenAI-only account setting, so it
// must not veto the official Grok >=200k ladder: a Grok account has no way to
// ever set that flag, which would make the group toggle unreachable.
func TestOpenAIGatewayServiceRecordUsage_GrokLongContextFollowsGroupToggleOnly(t *testing.T) {
	baseInput := 250000 * 2e-6
	baseOutput := 1000 * 6e-6

	grokAccount := func(id int64) *Account {
	placeholderID: id, Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
placeholder

	t.Run("group on applies the official ladder", func(t *testing.T) {
		usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
		svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
		err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
			Result: &OpenAIForwardResult{
				RequestID: "resp_grok_longctx_on",
				Usage:     OpenAIUsage{InputTokens: 250000, OutputTokens: 1000placeholder,
				Model:     "grok-4.5",
				Duration:  time.Second,
		placeholder,
			APIKey:  openAIRecordUsageAPIKeyWithGroup(svc, 1030, true),
			User:    &User{ID: 2030placeholder,
			Account: grokAccount(3030),
	placeholder)
	placeholder
		require.True(t, usageRepo.lastLog.LongContextBillingApplied)
		require.InDelta(t, baseInput*2, usageRepo.lastLog.InputCost, 1e-10)
		require.InDelta(t, baseOutput*2, usageRepo.lastLog.OutputCost, 1e-10)
placeholder)

	t.Run("group off keeps the base card", func(t *testing.T) {
		usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
		svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
		err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
			Result: &OpenAIForwardResult{
				RequestID: "resp_grok_longctx_off",
				Usage:     OpenAIUsage{InputTokens: 250000, OutputTokens: 1000placeholder,
				Model:     "grok-4.5",
				Duration:  time.Second,
		placeholder,
			APIKey:  openAIRecordUsageAPIKeyWithGroup(svc, 1031, false),
			User:    &User{ID: 2031placeholder,
			Account: grokAccount(3031),
	placeholder)
	placeholder
		require.False(t, usageRepo.lastLog.LongContextBillingApplied)
		require.InDelta(t, baseInput, usageRepo.lastLog.InputCost, 1e-10)
		require.InDelta(t, baseOutput, usageRepo.lastLog.OutputCost, 1e-10)
placeholder)
placeholder

func TestOpenAIGatewayServiceRecordUsage_SparkShadowUsesCurrentParentBillingSetting(t *testing.T) {
	tests := []struct {
		name          string
		parentEnabled bool
placeholder{
		{name: "parent opt out overrides stale enabled shadow", parentEnabled: falseplaceholder,
		{name: "parent opt in overrides stale disabled shadow", parentEnabled: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
			accountRepo := &openAIRecordUsageAccountRepoStub{account: &Account{
				ID:       4016,
				Platform: PlatformOpenAI,
				Type:     AccountTypeOAuth,
				Extra:    map[string]any{openAILongContextBillingEnabledKey: tt.parentEnabledplaceholder,
		placeholderplaceholder
			svc := newOpenAIRecordUsageServiceForTest(
				usageRepo,
				&openAIRecordUsageUserRepoStub{placeholder,
				&openAIRecordUsageSubRepoStub{placeholder,
				nil,
			)
			svc.accountRepo = accountRepo
			parentID := int64(4016)

			err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
				Result: &OpenAIForwardResult{
					RequestID: "resp_gpt54_shadow_parent_setting",
					Usage:     OpenAIUsage{InputTokens: 300000, OutputTokens: 2000placeholder,
					Model:     "gpt-5.4-2026-03-05",
					Duration:  time.Second,
			placeholder,
				APIKey: openAIRecordUsageAPIKeyWithGroup(svc, 1016, false),
				User:   &User{ID: 2016placeholder,
				Account: &Account{
					ID:              3016,
					Platform:        PlatformOpenAI,
					Type:            AccountTypeOAuth,
					ParentAccountID: &parentID,
					QuotaDimension:  QuotaDimensionSpark,
					Extra: map[string]any{
						openAILongContextBillingEnabledKey: !tt.parentEnabled,
				placeholder,
			placeholder,
		placeholder)

		placeholder
			require.Equal(t, 1, accountRepo.calls)
			require.Equal(t, tt.parentEnabled, usageRepo.lastLog.LongContextBillingApplied)
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceRecordUsage_ServiceTierPriorityUsesFastPricing(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	serviceTier := "priority"
	usage := OpenAIUsage{InputTokens: 100, OutputTokens: 50placeholder

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:   "resp_service_tier_priority",
			ServiceTier: &serviceTier,
			Usage:       usage,
			Model:       "gpt-5.4",
			Duration:    time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1015placeholder,
		User:    &User{ID: 2015placeholder,
		Account: &Account{ID: 3015placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.ServiceTier)
	require.Equal(t, serviceTier, *usageRepo.lastLog.ServiceTier)

	baseCost, calcErr := svc.billingService.CalculateCost("gpt-5.4", UsageTokens{InputTokens: 100, OutputTokens: 50placeholder, 1.0)
	require.NoError(t, calcErr)
	require.InDelta(t, baseCost.TotalCost*2, usageRepo.lastLog.TotalCost, 1e-10)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ServiceTierFlexHalvesCost(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	serviceTier := "flex"
	usage := OpenAIUsage{InputTokens: 100, OutputTokens: 50, CacheReadInputTokens: 20placeholder

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:   "resp_service_tier_flex",
			ServiceTier: &serviceTier,
			Usage:       usage,
			Model:       "gpt-5.4",
			Duration:    time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1016placeholder,
		User:    &User{ID: 2016placeholder,
		Account: &Account{ID: 3016placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)

	baseCost, calcErr := svc.billingService.CalculateCost("gpt-5.4", UsageTokens{InputTokens: 80, OutputTokens: 50, CacheReadTokens: 20placeholder, 1.0)
	require.NoError(t, calcErr)
	require.InDelta(t, baseCost.TotalCost*0.5, usageRepo.lastLog.TotalCost, 1e-10)
placeholder

func TestNormalizeOpenAIServiceTier(t *testing.T) {
	t.Run("fast maps to priority", func(t *testing.T) {
		got := normalizeOpenAIServiceTier(" fast ")
		require.NotNil(t, got)
		require.Equal(t, "priority", *got)
placeholder)

	t.Run("openai official tiers preserved", func(t *testing.T) {
		// OpenAI 官方文档定义的合法 tier 值都应被透传保留，避免因白名单过窄
		// 静默剥离客户端显式发送的合法字段。Codex 客户端只发 priority/flex，
		// 所以扩大白名单对 Codex 流量零影响（见 codex-rs/core/src/client.rs）。
		for _, tier := range []string{"priority", "flex", "auto", "default", "scale"placeholder {
			got := normalizeOpenAIServiceTier(tier)
			require.NotNil(t, got, "tier %q should not be normalized to nil", tier)
			require.Equal(t, tier, *got)
	placeholder
placeholder)

	t.Run("invalid ignored", func(t *testing.T) {
		require.Nil(t, normalizeOpenAIServiceTier("turbo"))
		require.Nil(t, normalizeOpenAIServiceTier("xxx"))
placeholder)
placeholder

func TestExtractOpenAIServiceTier(t *testing.T) {
	require.Equal(t, "priority", *extractOpenAIServiceTier(map[string]any{"service_tier": "fast"placeholder))
	require.Equal(t, "flex", *extractOpenAIServiceTier(map[string]any{"service_tier": "flex"placeholder))
	require.Equal(t, "auto", *extractOpenAIServiceTier(map[string]any{"service_tier": "auto"placeholder))
	require.Equal(t, "default", *extractOpenAIServiceTier(map[string]any{"service_tier": "default"placeholder))
	require.Equal(t, "scale", *extractOpenAIServiceTier(map[string]any{"service_tier": "scale"placeholder))
	require.Nil(t, extractOpenAIServiceTier(map[string]any{"service_tier": 1placeholder))
	require.Nil(t, extractOpenAIServiceTier(nil))
placeholder

func TestExtractOpenAIServiceTierFromBody(t *testing.T) {
	require.Equal(t, "priority", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"fast"placeholder`)))
	require.Equal(t, "flex", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"flex"placeholder`)))
	require.Equal(t, "auto", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"auto"placeholder`)))
	require.Equal(t, "default", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"default"placeholder`)))
	require.Equal(t, "scale", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"scale"placeholder`)))
	require.Nil(t, extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"turbo"placeholder`)))
	require.Nil(t, extractOpenAIServiceTierFromBody(nil))
placeholder

func TestOpenAIGatewayServiceRecordUsage_UsesRequestedModelAndUpstreamModelMetadataFields(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	serviceTier := "priority"
	reasoning := "high"

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:       "resp_billing_model_override",
			BillingModel:    "gpt-5.1-codex",
			Model:           "gpt-5.1",
			UpstreamModel:   "gpt-5.1-codex",
			ServiceTier:     &serviceTier,
			ReasoningEffort: &reasoning,
			Usage: OpenAIUsage{
				InputTokens:  20,
				OutputTokens: 10,
		placeholder,
			Duration:     2 * time.Second,
			FirstTokenMs: func() *int { v := 120; return &v placeholder(),
	placeholder,
		APIKey:    &APIKey{ID: 10, GroupID: i64p(11), Group: &Group{ID: 11, RateMultiplier: 1.2placeholderplaceholder,
		User:      &User{ID: 20placeholder,
		Account:   &Account{ID: 30placeholder,
		UserAgent: "codex-cli/1.0",
		IPAddress: "127.0.0.1",
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "gpt-5.1", usageRepo.lastLog.Model)
	require.Equal(t, "gpt-5.1", usageRepo.lastLog.RequestedModel)
	require.NotNil(t, usageRepo.lastLog.UpstreamModel)
	require.Equal(t, "gpt-5.1-codex", *usageRepo.lastLog.UpstreamModel)
	require.NotNil(t, usageRepo.lastLog.ServiceTier)
	require.Equal(t, serviceTier, *usageRepo.lastLog.ServiceTier)
	require.NotNil(t, usageRepo.lastLog.ReasoningEffort)
	require.Equal(t, reasoning, *usageRepo.lastLog.ReasoningEffort)
	require.NotNil(t, usageRepo.lastLog.UserAgent)
	require.Equal(t, "codex-cli/1.0", *usageRepo.lastLog.UserAgent)
	require.NotNil(t, usageRepo.lastLog.IPAddress)
	require.Equal(t, "127.0.0.1", *usageRepo.lastLog.IPAddress)
	require.NotNil(t, usageRepo.lastLog.GroupID)
	require.Equal(t, int64(11), *usageRepo.lastLog.GroupID)
	require.Equal(t, 1, userRepo.deductCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_PreservesChannelMappedUpstreamModel(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "openai_channel_mapping_models",
			Model:         "gpt-5.6-terra",
			UpstreamModel: "gpt-5.6-terra",
			Usage: OpenAIUsage{
				InputTokens:  20,
				OutputTokens: 10,
		placeholder,
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
		ChannelUsageFields: ChannelUsageFields{
			OriginalModel:      "gpt-5.6-sol",
			ChannelMappedModel: "gpt-5.6-terra",
	placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "gpt-5.6-sol", usageRepo.lastLog.RequestedModel)
	require.Equal(t, "gpt-5.6-terra", usageRepo.lastLog.Model)
	require.NotNil(t, usageRepo.lastLog.UpstreamModel)
	require.Equal(t, "gpt-5.6-terra", *usageRepo.lastLog.UpstreamModel)
placeholder

func TestOpenAIGatewayServiceRecordUsage_PreservesLoopedChannelAndAccountUpstreamModel(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "openai_looped_mapping_models",
			Model:         "gpt-5.6-terra",
			UpstreamModel: "gpt-5.6-sol",
			Usage:         OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
		ChannelUsageFields: ChannelUsageFields{
			OriginalModel:      "gpt-5.6-sol",
			ChannelMappedModel: "gpt-5.6-terra",
	placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "gpt-5.6-sol", usageRepo.lastLog.RequestedModel)
	require.Equal(t, "gpt-5.6-terra", usageRepo.lastLog.Model)
	require.NotNil(t, usageRepo.lastLog.UpstreamModel)
	require.Equal(t, "gpt-5.6-sol", *usageRepo.lastLog.UpstreamModel)
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillsMappedRequestsUsingRequestedModel(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	usage := OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder

	// Billing should use the requested model ("gpt-5.1"), not the upstream mapped model ("gpt-5.1-codex").
	// This ensures pricing is always based on the model the user requested.
	expectedCost, err := svc.billingService.CalculateCost("gpt-5.1", UsageTokens{
		InputTokens:  20,
		OutputTokens: 10,
placeholder, 1.1)
placeholder

	err = svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "resp_upstream_model_billing_fallback",
			Model:         "gpt-5.1",
			UpstreamModel: "gpt-5.1-codex",
			Usage:         usage,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "gpt-5.1", usageRepo.lastLog.Model)
	require.Equal(t, expectedCost.ActualCost, usageRepo.lastLog.ActualCost)
	require.Equal(t, expectedCost.TotalCost, usageRepo.lastLog.TotalCost)
	require.Equal(t, expectedCost.ActualCost, userRepo.lastAmount)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ChannelMappedDoesNotOverrideBillingModelWhenUnmapped(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	usage := OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder

	// 渠道未发生模型映射时，应使用 result.BillingModel 中记录的实际上游计费模型，
	// 而不是未映射的原始请求模型。
	expectedCost, err := svc.billingService.CalculateCost("gpt-5.1", UsageTokens{
		InputTokens:  20,
		OutputTokens: 10,
placeholder, 1.1)
placeholder

	err = svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "resp_channel_unmapped_billing",
			Model:         "glm",
			BillingModel:  "gpt-5.1",
			UpstreamModel: "gpt-5.1",
			Usage:         usage,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
		ChannelUsageFields: ChannelUsageFields{
			ChannelID:          1,
			OriginalModel:      "glm",
			ChannelMappedModel: "glm", // channel did NOT map
			BillingModelSource: BillingModelSourceChannelMapped,
	placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, expectedCost.ActualCost, usageRepo.lastLog.ActualCost)
	require.True(t, usageRepo.lastLog.ActualCost > 0, "cost must not be zero")
placeholder

func TestOpenAIGatewayServiceRecordUsage_ChannelMappedOverridesBillingModelWhenMapped(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	usage := OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder

	// When channel DID map the model (ChannelMappedModel != OriginalModel),
	// billing should use the channel-mapped model, honoring admin intent.
	expectedCost, err := svc.billingService.CalculateCost("gpt-5.1", UsageTokens{
		InputTokens:  20,
		OutputTokens: 10,
placeholder, 1.1)
placeholder

	err = svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "resp_channel_mapped_billing",
			Model:         "glm",
			BillingModel:  "gpt-5.1-codex",
			UpstreamModel: "gpt-5.1-codex",
			Usage:         usage,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
		ChannelUsageFields: ChannelUsageFields{
			ChannelID:          1,
			OriginalModel:      "glm",
			ChannelMappedModel: "gpt-5.1", // channel mapped glm → gpt-5.1
			BillingModelSource: BillingModelSourceChannelMapped,
	placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, expectedCost.ActualCost, usageRepo.lastLog.ActualCost)
	require.True(t, usageRepo.lastLog.ActualCost > 0, "cost must not be zero")
placeholder

func TestOpenAIGatewayServiceRecordUsage_ResponsesMappedBillingModelHonorsBillingModelSource(t *testing.T) {
	usage := OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder
	tokens := UsageTokens{InputTokens: 20, OutputTokens: 10placeholder

	tests := []struct {
		name               string
		billingModelSource string
		wantBillingModel   string
placeholder{
		{
			name:               "upstream uses mapped billing model",
			billingModelSource: BillingModelSourceUpstream,
			wantBillingModel:   "gpt-5.5",
	placeholder,
		{
			name:               "requested overrides mapped billing model",
			billingModelSource: BillingModelSourceRequested,
			wantBillingModel:   "gpt-5.4",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
			userRepo := &openAIRecordUsageUserRepoStub{placeholder
			subRepo := &openAIRecordUsageSubRepoStub{placeholder
			svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

			expectedCost, err := svc.billingService.CalculateCost(tt.wantBillingModel, tokens, 1.1)
		placeholder

			err = svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
				Result: &OpenAIForwardResult{
					RequestID:     "resp_mapped_billing_model_source",
					Model:         "gpt-5.4",
					BillingModel:  "gpt-5.5",
					UpstreamModel: "gpt-5.5",
					Usage:         usage,
					Duration:      time.Second,
			placeholder,
				APIKey:  &APIKey{ID: 10placeholder,
				User:    &User{ID: 20placeholder,
				Account: &Account{ID: 30placeholder,
				ChannelUsageFields: ChannelUsageFields{
					OriginalModel:      "gpt-5.4",
					ChannelMappedModel: "gpt-5.4",
					BillingModelSource: tt.billingModelSource,
			placeholder,
		placeholder)

		placeholder
			require.NotNil(t, usageRepo.lastLog)
			require.Equal(t, "gpt-5.4", usageRepo.lastLog.Model)
			require.InDelta(t, expectedCost.ActualCost, usageRepo.lastLog.ActualCost, 1e-12)
			require.InDelta(t, expectedCost.ActualCost, userRepo.lastAmount, 1e-12)
			require.True(t, usageRepo.lastLog.ActualCost > 0, "cost must not be zero")
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServiceRecordUsage_BillsCompactOpenAIModelAlias(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	usage := OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder

	expectedCost, err := svc.billingService.CalculateCost("gpt-5.5", UsageTokens{
		InputTokens:  20,
		OutputTokens: 10,
placeholder, 1.1)
placeholder

	err = svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "resp_compact_openai_alias",
			Model:         "gpt5.5",
			UpstreamModel: "gpt-5.4",
			Usage:         usage,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "gpt5.5", usageRepo.lastLog.Model)
	require.NotNil(t, usageRepo.lastLog.UpstreamModel)
	require.Equal(t, "gpt-5.4", *usageRepo.lastLog.UpstreamModel)
	require.InDelta(t, expectedCost.ActualCost, usageRepo.lastLog.ActualCost, 1e-12)
	require.True(t, usageRepo.lastLog.ActualCost > 0, "cost must not be zero")
	require.InDelta(t, expectedCost.ActualCost, userRepo.lastAmount, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_FallsBackToUpstreamModelWhenPrimaryUnpriceable(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	usage := OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder

	expectedCost, err := svc.billingService.CalculateCost("gpt-5.4", UsageTokens{
		InputTokens:  20,
		OutputTokens: 10,
placeholder, 1.1)
placeholder

	err = svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:     "resp_unpriceable_primary_upstream_fallback",
			Model:         "not-priceable-alias",
			BillingModel:  "not-priceable-alias",
			UpstreamModel: "gpt-5.4",
			Usage:         usage,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, expectedCost.ActualCost, usageRepo.lastLog.ActualCost, 1e-12)
	require.True(t, usageRepo.lastLog.ActualCost > 0, "cost must not be zero")
	require.InDelta(t, expectedCost.ActualCost, userRepo.lastAmount, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_UnpricedTokenModelFallsBackToZeroCostUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_unpriceable_without_upstream",
			Model:     "not-priceable-alias",
			Usage:     OpenAIUsage{InputTokens: 20, OutputTokens: 10placeholder,
			Duration:  time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 10placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "not-priceable-alias", usageRepo.lastLog.Model)
	require.Equal(t, 20, usageRepo.lastLog.InputTokens)
	require.Equal(t, 10, usageRepo.lastLog.OutputTokens)
	require.Zero(t, usageRepo.lastLog.TotalCost)
	require.Zero(t, usageRepo.lastLog.ActualCost)
	require.Equal(t, 0, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_SubscriptionBillingSetsSubscriptionFields(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	subscription := &UserSubscription{ID: 99placeholder

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_subscription_billing",
			Usage:     OpenAIUsage{InputTokens: 10, OutputTokens: 5placeholder,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey:       &APIKey{ID: 100, GroupID: i64p(88), Group: &Group{ID: 88, SubscriptionType: SubscriptionTypeSubscription, RateMultiplier: 1.0placeholderplaceholder,
		User:         &User{ID: 200placeholder,
		Account:      &Account{ID: 300placeholder,
		Subscription: subscription,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, BillingTypeSubscription, usageRepo.lastLog.BillingType)
	require.NotNil(t, usageRepo.lastLog.SubscriptionID)
	require.Equal(t, subscription.ID, *usageRepo.lastLog.SubscriptionID)
	require.Equal(t, 1, subRepo.incrementCalls)
	require.Equal(t, 0, userRepo.deductCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_SimpleModeSkipsBillingAfterPersist(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)
	svc.cfg.RunMode = config.RunModeSimple

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_simple_mode",
			Usage:     OpenAIUsage{InputTokens: 10, OutputTokens: 5placeholder,
			Model:     "gpt-5.1",
			Duration:  time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1000placeholder,
		User:    &User{ID: 2000placeholder,
		Account: &Account{ID: 3000placeholder,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 0, userRepo.deductCalls)
	require.Equal(t, 0, subRepo.incrementCalls)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ImageOnlyUsageStillPersists(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_only_usage",
			Model:      "gpt-image-2",
			ImageCount: 2,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1007placeholder,
		User:    &User{ID: 2007placeholder,
		Account: &Account{ID: 3007placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 2, usageRepo.lastLog.ImageCount)
	require.NotNil(t, usageRepo.lastLog.ImageSize)
	require.Equal(t, "1K", *usageRepo.lastLog.ImageSize)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_EmptyImageSizeDefaultsBeforeBillingAndPersistence(t *testing.T) {
	imagePrice2K := 0.31
	groupID := int64(1201)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_default_size",
			Model:      "gpt-image-2",
			ImageCount: 2,
			ImageSize:  "",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      11201,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: 1.0,
				ImagePrice2K:   &imagePrice2K,
		placeholder,
	placeholder,
		User:    &User{ID: 21201placeholder,
		Account: &Account{ID: 31201placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 2, usageRepo.lastLog.ImageCount)
	require.NotNil(t, usageRepo.lastLog.ImageSize)
	require.Equal(t, ImageBillingSize2K, *usageRepo.lastLog.ImageSize)
	require.NotNil(t, usageRepo.lastLog.ImageSizeSource)
	require.Equal(t, ImageSizeSourceDefault, *usageRepo.lastLog.ImageSizeSource)
	require.Nil(t, usageRepo.lastLog.ImageInputSize)
	require.Nil(t, usageRepo.lastLog.ImageOutputSize)
	require.InDelta(t, 0.62, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.62, usageRepo.lastLog.ActualCost, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_OutputImageSizeWinsBeforeBillingAndPersistence(t *testing.T) {
	imagePrice1K := 0.11
	imagePrice4K := 0.44
	groupID := int64(1202)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:        "resp_image_output_size",
			Model:            "gpt-image-2",
			ImageCount:       1,
			ImageInputSize:   "1024x1024",
			ImageOutputSizes: []string{"3840x2160"placeholder,
			Duration:         time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      11202,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: 1.0,
				ImagePrice1K:   &imagePrice1K,
				ImagePrice4K:   &imagePrice4K,
		placeholder,
	placeholder,
		User:    &User{ID: 21202placeholder,
		Account: &Account{ID: 31202placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.ImageSize)
	require.Equal(t, ImageBillingSize4K, *usageRepo.lastLog.ImageSize)
	require.NotNil(t, usageRepo.lastLog.ImageInputSize)
	require.Equal(t, "1024x1024", *usageRepo.lastLog.ImageInputSize)
	require.NotNil(t, usageRepo.lastLog.ImageOutputSize)
	require.Equal(t, "3840x2160", *usageRepo.lastLog.ImageOutputSize)
	require.NotNil(t, usageRepo.lastLog.ImageSizeSource)
	require.Equal(t, ImageSizeSourceOutput, *usageRepo.lastLog.ImageSizeSource)
	require.Equal(t, map[string]int{ImageBillingSize4K: 1placeholder, usageRepo.lastLog.ImageSizeBreakdown)
	require.InDelta(t, 0.44, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.44, usageRepo.lastLog.ActualCost, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ImageUsesPerImageBillingEvenWithUsageTokens(t *testing.T) {
	imagePrice := 0.02
	groupID := int64(12)

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_image_per_request",
			Model:     "gpt-image-2",
			Usage: OpenAIUsage{
				InputTokens:       1110,
				OutputTokens:      1756,
				ImageOutputTokens: 1756,
		placeholder,
			ImageCount: 2,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      1008,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: 1.0,
				ImagePrice1K:   &imagePrice,
		placeholder,
	placeholder,
		User:    &User{ID: 2008placeholder,
		Account: &Account{ID: 3008placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
	require.Equal(t, 2, usageRepo.lastLog.ImageCount)
	require.InDelta(t, 0.04, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.04, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.0, usageRepo.lastLog.InputCost, 1e-12)
	require.InDelta(t, 0.0, usageRepo.lastLog.OutputCost, 1e-12)
	require.InDelta(t, 0.0, usageRepo.lastLog.ImageOutputCost, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ImageSharedMultiplierPreservesExistingBehavior(t *testing.T) {
	imagePrice := 0.2
	groupID := int64(121)

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_shared_multiplier",
			Model:      "gpt-image-2",
			ImageCount: 1,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10121,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				RateMultiplier:       0.15,
				ImageRateIndependent: false,
				ImageRateMultiplier:  1,
				ImagePrice1K:         &imagePrice,
		placeholder,
	placeholder,
		User:    &User{ID: 20121placeholder,
		Account: &Account{ID: 30121placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, 0.2, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.03, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.15, usageRepo.lastLog.RateMultiplier, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ImageSharedMultiplierUsesUserGroupOverride(t *testing.T) {
	imagePrice := 0.5
	userRate := 0.2
	groupID := int64(125)

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(
		usageRepo,
		&openAIRecordUsageUserRepoStub{placeholder,
		&openAIRecordUsageSubRepoStub{placeholder,
		&openAIUserGroupRateRepoStub{rate: &userRateplaceholder,
	)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_user_group_override",
			Model:      "gpt-image-2",
			ImageCount: 1,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10125,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				RateMultiplier:       0.15,
				ImageRateIndependent: false,
				ImageRateMultiplier:  1,
				ImagePrice1K:         &imagePrice,
		placeholder,
	placeholder,
		User:    &User{ID: 20125placeholder,
		Account: &Account{ID: 30125placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, 0.5, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.1, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.2, usageRepo.lastLog.RateMultiplier, 1e-12)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ImageIndependentMultiplierUsesImageRate(t *testing.T) {
	imagePrice := 0.2
	groupID := int64(122)

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_independent_multiplier",
			Model:      "gpt-image-2",
			ImageCount: 1,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10122,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				RateMultiplier:       0.15,
				ImageRateIndependent: true,
				ImageRateMultiplier:  1,
				ImagePrice1K:         &imagePrice,
		placeholder,
	placeholder,
		User:    &User{ID: 20122placeholder,
		Account: &Account{ID: 30122placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, 0.2, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.2, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 1.0, usageRepo.lastLog.RateMultiplier, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestGrokVideoBillingUsesSeparateVideoRateMultiplier(t *testing.T) {
	imagePrice2K := 0.4
	videoPrice480P := 0.08
	groupID := int64(126)

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:    "video-request-123",
			ResponseID:   "video-request-123",
			Model:        "grok-imagine-video-1.5",
			BillingModel: "grok-imagine-video-1.5",
			// Pure video completion clears ImageCount (handler contract).
			ImageCount:           0,
			VideoCount:           1,
			VideoResolution:      VideoBillingResolution480P,
			VideoDurationSeconds: 1,
			Duration:             time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10126,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				Platform:             PlatformGrok,
				RateMultiplier:       0.15,
				ImageRateIndependent: true,
				ImageRateMultiplier:  0.5,
				ImagePrice2K:         &imagePrice2K,
				VideoRateIndependent: true,
				VideoRateMultiplier:  0.25,
				VideoPrice480P:       &videoPrice480P,
		placeholder,
	placeholder,
		User:    &User{ID: 20126placeholder,
		Account: &Account{ID: 30126, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "grok-imagine-video-1.5", usageRepo.lastLog.Model)
	require.Equal(t, 0, usageRepo.lastLog.ImageCount)
	require.Nil(t, usageRepo.lastLog.ImageSize)
	require.InDelta(t, 0.08, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.02, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.25, usageRepo.lastLog.RateMultiplier, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeVideo), *usageRepo.lastLog.BillingMode)
	require.Equal(t, 1, usageRepo.lastLog.VideoCount)
	require.NotNil(t, usageRepo.lastLog.VideoResolution)
	require.Equal(t, VideoBillingResolution480P, *usageRepo.lastLog.VideoResolution)
	require.NotNil(t, usageRepo.lastLog.VideoDurationSeconds)
	require.Equal(t, 1, *usageRepo.lastLog.VideoDurationSeconds)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GrokVideoUsesDefaultRateCard(t *testing.T) {
	groupID := int64(1261)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:       "video-default-rate-card",
			ResponseID:      "video-default-rate-card",
			Model:           "grok-imagine-video-1.5",
			BillingModel:    "grok-imagine-video-1.5",
			ImageCount:      0,
			VideoCount:      1,
			VideoResolution: VideoBillingResolution720P,
			Duration:        time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      101261,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				Platform:       PlatformGrok,
				RateMultiplier: 1,
		placeholder,
	placeholder,
		User:    &User{ID: 201261placeholder,
		Account: &Account{ID: 301261, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Nil(t, usageRepo.lastLog.ImageSize)
	// 结果未携带 duration 时按上游默认 8 秒计费：0.14 USD/s × 8s。
	require.InDelta(t, 0.14*8, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.14*8, usageRepo.lastLog.ActualCost, 1e-12)
	require.Equal(t, 0, usageRepo.lastLog.ImageCount)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeVideo), *usageRepo.lastLog.BillingMode)
	require.Equal(t, 1, usageRepo.lastLog.VideoCount)
	require.NotNil(t, usageRepo.lastLog.VideoDurationSeconds)
	require.Equal(t, VideoBillingDefaultDurationSeconds, *usageRepo.lastLog.VideoDurationSeconds)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GroupImagePriceOverridesChannelImagePrice(t *testing.T) {
	groupID := int64(127)
	channelPrice := 0.201
	groupImagePrice2K := 0.021
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAIImageChannelPricingResolverForTest(t, groupID, "grok-imagine-image-quality", channelPrice)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:    "resp_grok_image_group_price",
			Model:        "grok-imagine-image-quality",
			BillingModel: "grok-imagine-image-quality",
			ImageCount:   1,
			ImageSize:    ImageBillingSize2K,
			Duration:     time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10127,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				Platform:             PlatformGrok,
				RateMultiplier:       1,
				ImageRateIndependent: true,
				ImageRateMultiplier:  1,
				ImagePrice2K:         &groupImagePrice2K,
		placeholder,
	placeholder,
		User:    &User{ID: 20127placeholder,
		Account: &Account{ID: 30127, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 1, usageRepo.lastLog.ImageCount)
	require.Equal(t, ImageBillingSize2K, *usageRepo.lastLog.ImageSize)
	require.InDelta(t, 0.021, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.021, usageRepo.lastLog.ActualCost, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GroupVideoPriceOverridesChannelImagePrice(t *testing.T) {
	groupID := int64(128)
	channelPrice := 0.201
	groupVideoPrice720P := 0.037
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAIImageChannelPricingResolverForTest(t, groupID, "grok-imagine-video", channelPrice)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:            "resp_grok_video_group_price",
			Model:                "grok-imagine-video",
			BillingModel:         "grok-imagine-video",
			ImageCount:           0,
			VideoCount:           1,
			VideoResolution:      VideoBillingResolution720P,
			VideoDurationSeconds: 1,
			Duration:             time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10128,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				Platform:             PlatformGrok,
				RateMultiplier:       1,
				VideoRateIndependent: true,
				VideoRateMultiplier:  1,
				VideoPrice720P:       &groupVideoPrice720P,
		placeholder,
	placeholder,
		User:    &User{ID: 20128placeholder,
		Account: &Account{ID: 30128, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 0, usageRepo.lastLog.ImageCount)
	require.Nil(t, usageRepo.lastLog.ImageSize)
	require.InDelta(t, 0.037, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.037, usageRepo.lastLog.ActualCost, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeVideo), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_GroupVideoModelPriceOverridesFlatAndChannelPrice(t *testing.T) {
	groupID := int64(129)
	channelPrice := 0.201
	flatVideoPrice720P := 0.037
	modelVideoPrice720P := 0.123
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAIImageChannelPricingResolverForTest(t, groupID, "grok-imagine-video-1.5-preview", channelPrice)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:            "resp_grok_video_model_price",
			Model:                "grok-imagine-video-1.5-preview",
			BillingModel:         "grok-imagine-video-1.5-preview",
			ImageCount:           0,
			VideoCount:           1,
			VideoResolution:      VideoBillingResolution720P,
			VideoDurationSeconds: 2,
			Duration:             time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10129,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				Platform:             PlatformGrok,
				RateMultiplier:       1,
				VideoRateIndependent: true,
				VideoRateMultiplier:  1,
				VideoPrice720P:       &flatVideoPrice720P,
				VideoModelPrices: map[string]map[string]float64{
					VideoPriceFamilyGrokImagineVideo15: {VideoBillingResolution720P: modelVideoPrice720Pplaceholder,
			placeholder,
		placeholder,
	placeholder,
		User:    &User{ID: 20129placeholder,
		Account: &Account{ID: 30129, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, modelVideoPrice720P*2, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, modelVideoPrice720P*2, usageRepo.lastLog.ActualCost, 1e-12)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeVideo), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_HydratesGroupImagePriceWhenAuthSnapshotOmitsIt(t *testing.T) {
	groupID := int64(130)
	groupImagePrice2K := 0.021
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	channelService := &ChannelService{groupRepo: &openAIMediaPriceGroupRepoStub{group: &Group{
		ID:             groupID,
		Platform:       PlatformGrok,
		RateMultiplier: 1,
		ImagePrice2K:   &groupImagePrice2K,
placeholderplaceholderplaceholder
	channelCache := newEmptyChannelCache()
	channelCache.loadedAt = time.Now()
	channelService.cache.Store(channelCache)
	svc.channelService = channelService
	refreshed := svc.apiKeyWithFreshGroupMediaPricing(context.Background(), &APIKey{GroupID: i64p(groupID), Group: &Group{ID: groupIDplaceholderplaceholder)
	require.NotNil(t, refreshed.Group.ImagePrice2K)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:    "resp_grok_image_hydrated_price",
			Model:        "grok-imagine-image-quality",
			BillingModel: "grok-imagine-image-quality",
			ImageCount:   1,
			ImageSize:    ImageBillingSize2K,
			Duration:     time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10130,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				Platform:       PlatformGrok,
				RateMultiplier: 1,
		placeholder,
	placeholder,
		User:    &User{ID: 20130placeholder,
		Account: &Account{ID: 30130, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, 0.021, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.021, usageRepo.lastLog.ActualCost, 1e-12)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_HydratesGroupVideoPriceWhenAuthSnapshotOmitsIt(t *testing.T) {
	groupID := int64(131)
	groupVideoPrice720P := 0.037
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	channelService := &ChannelService{groupRepo: &openAIMediaPriceGroupRepoStub{group: &Group{
		ID:             groupID,
		Platform:       PlatformGrok,
		RateMultiplier: 1,
		VideoPrice720P: &groupVideoPrice720P,
placeholderplaceholderplaceholder
	channelCache := newEmptyChannelCache()
	channelCache.loadedAt = time.Now()
	channelService.cache.Store(channelCache)
	svc.channelService = channelService
	refreshed := svc.apiKeyWithFreshGroupMediaPricing(context.Background(), &APIKey{GroupID: i64p(groupID), Group: &Group{ID: groupIDplaceholderplaceholder)
	require.NotNil(t, refreshed.Group.VideoPrice720P)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:            "resp_grok_video_hydrated_price",
			Model:                "grok-imagine-video",
			BillingModel:         "grok-imagine-video",
			ImageCount:           0,
			VideoCount:           1,
			VideoResolution:      VideoBillingResolution720P,
			VideoDurationSeconds: 1,
			Duration:             time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10131,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				Platform:       PlatformGrok,
				RateMultiplier: 1,
		placeholder,
	placeholder,
		User:    &User{ID: 20131placeholder,
		Account: &Account{ID: 30131, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Nil(t, usageRepo.lastLog.ImageSize)
	require.InDelta(t, 0.037, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.037, usageRepo.lastLog.ActualCost, 1e-12)
	require.Equal(t, string(BillingModeVideo), *usageRepo.lastLog.BillingMode)
placeholder

// 视频请求命中渠道 token 计费时走 token 路径；此时行是 billing_mode='token'、image_count=1、
// image_size=NULL，必须携带 video_count>0 才能通过 usage_logs 的 image_size check 约束
// （迁移 172），否则整个计费事务会因约束违反而丢失。
func TestOpenAIGatewayServiceRecordUsage_GrokVideoWithTokenChannelPricingKeepsVideoMetadata(t *testing.T) {
	groupID := int64(132)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAITokenImageChannelPricingResolverForTest(t, groupID, "grok-imagine-video")

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:            "resp_grok_video_token_channel",
			Model:                "grok-imagine-video",
			BillingModel:         "grok-imagine-video",
			ImageCount:           0,
			VideoCount:           1,
			VideoResolution:      VideoBillingResolution720P,
			VideoDurationSeconds: 5,
			Usage:                OpenAIUsage{InputTokens: 100, OutputTokens: 200placeholder,
			Duration:             time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10132,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				Platform:       PlatformGrok,
				RateMultiplier: 1,
		placeholder,
	placeholder,
		User:    &User{ID: 20132placeholder,
		Account: &Account{ID: 30132, Platform: PlatformGrokplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeToken), *usageRepo.lastLog.BillingMode)
	require.Nil(t, usageRepo.lastLog.ImageSize)
	require.Equal(t, 0, usageRepo.lastLog.ImageCount)
	require.Equal(t, 1, usageRepo.lastLog.VideoCount)
	require.NotNil(t, usageRepo.lastLog.VideoResolution)
	require.Equal(t, VideoBillingResolution720P, *usageRepo.lastLog.VideoResolution)
	require.NotNil(t, usageRepo.lastLog.VideoDurationSeconds)
	require.Equal(t, 5, *usageRepo.lastLog.VideoDurationSeconds)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ChannelImageBillingUsesImageCountAndSharedMultiplier(t *testing.T) {
	groupID := int64(123)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAIImageChannelPricingResolverForTest(t, groupID, "gpt-image-2", 0.25)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_channel_shared",
			Model:      "gpt-image-2",
			ImageCount: 3,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10123,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				RateMultiplier:       0.15,
				ImageRateIndependent: false,
				ImageRateMultiplier:  1,
		placeholder,
	placeholder,
		User:    &User{ID: 20placeholder,
		Account: &Account{ID: 30placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, 0.75, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.1125, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.15, usageRepo.lastLog.RateMultiplier, 1e-12)
	require.Equal(t, 3, usageRepo.lastLog.ImageCount)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func TestOpenAIGatewayServiceRecordUsage_ChannelImageBillingUsesImageCountAndIndependentMultiplier(t *testing.T) {
	groupID := int64(124)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder, nil)
	svc.resolver = newOpenAIImageChannelPricingResolverForTest(t, groupID, "gpt-image-2", 0.25)

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID:  "resp_image_channel_independent",
			Model:      "gpt-image-2",
			ImageCount: 3,
			ImageSize:  "1K",
			Duration:   time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      10124,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                   groupID,
				RateMultiplier:       0.15,
				ImageRateIndependent: true,
				ImageRateMultiplier:  1,
		placeholder,
	placeholder,
		User:    &User{ID: 20124placeholder,
		Account: &Account{ID: 30124placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.InDelta(t, 0.75, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.75, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 1.0, usageRepo.lastLog.RateMultiplier, 1e-12)
	require.Equal(t, 3, usageRepo.lastLog.ImageCount)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeImage), *usageRepo.lastLog.BillingMode)
placeholder

func newOpenAIImageChannelPricingResolverForTest(t *testing.T, groupID int64, model string, price float64) *ModelPricingResolver {
placeholder
	cache := newEmptyChannelCache()
	cache.pricingByGroupModel[channelModelKey{groupID: groupID, model: modelplaceholder] = &ChannelModelPricing{
		BillingMode:     BillingModeImage,
		PerRequestPrice: &price,
placeholder
	cache.channelByGroupID[groupID] = &Channel{ID: groupID, Status: StatusActiveplaceholder
	cache.groupPlatform[groupID] = ""
	cache.loadedAt = time.Now()
	cs := &ChannelService{placeholder
	cs.cache.Store(cache)
	return NewModelPricingResolver(cs, NewBillingService(&config.Config{placeholder, nil))
placeholder

func newOpenAITokenImageChannelPricingResolverForTest(t *testing.T, groupID int64, model string) *ModelPricingResolver {
placeholder
	inputPrice := 3e-6
	outputPrice := 15e-6
	imageOutputPrice := 15e-6
	cache := newEmptyChannelCache()
	cache.pricingByGroupModel[channelModelKey{groupID: groupID, model: modelplaceholder] = &ChannelModelPricing{
		BillingMode:      BillingModeToken,
		InputPrice:       &inputPrice,
		OutputPrice:      &outputPrice,
		ImageOutputPrice: &imageOutputPrice,
placeholder
	cache.channelByGroupID[groupID] = &Channel{ID: groupID, Status: StatusActiveplaceholder
	cache.groupPlatform[groupID] = ""
	cache.loadedAt = time.Now()
	cs := &ChannelService{placeholder
	cs.cache.Store(cache)
	return NewModelPricingResolver(cs, NewBillingService(&config.Config{placeholder, nil))
placeholder

func newOpenAITokenImageChannelPricingResolverWithTimeForTest(
	t *testing.T,
	groupID int64,
	model string,
	timePricing *ChannelTimePricing,
) *ModelPricingResolver {
placeholder
	resolver := newOpenAITokenImageChannelPricingResolverForTest(t, groupID, model)
	cached, ok := resolver.channelService.cache.Load().(*channelCache)
	require.True(t, ok)
	cached.pricingByGroupModel[channelModelKey{groupID: groupID, model: modelplaceholder].TimePricing = timePricing
	return resolver
placeholder

type openAIMediaPriceGroupRepoStub struct {
	GroupRepository
	group *Group
	err   error
placeholder

func (s *openAIMediaPriceGroupRepoStub) GetByIDLite(context.Context, int64) (*Group, error) {
	if s.err != nil {
		return nil, s.err
placeholder
	return s.group, nil
placeholder

func TestGatewayServiceCalculateRecordUsageCost_ChannelImageBillingUsesImageCount(t *testing.T) {
	groupID := int64(126)
	billingService := NewBillingService(&config.Config{placeholder, nil)
	svc := &GatewayService{
		billingService: billingService,
		resolver:       newOpenAIImageChannelPricingResolverForTest(t, groupID, "gemini-image", 0.25),
placeholder

	cost := svc.calculateRecordUsageCost(
		context.Background(),
		&ForwardResult{Model: "gemini-image", ImageCount: 2, ImageSize: "1K"placeholder,
		&APIKey{GroupID: i64p(groupID), Group: &Group{ID: groupIDplaceholderplaceholder,
		"gemini-image",
		0.15,
		1.0,
		time.Time{placeholder,
		nil,
	)

	require.NotNil(t, cost)
	require.Equal(t, string(BillingModeImage), cost.BillingMode)
	require.InDelta(t, 0.5, cost.TotalCost, 1e-12)
	require.InDelta(t, 0.5, cost.ActualCost, 1e-12)
placeholder

func TestGatewayServiceCalculateRecordUsageCost_ChannelImageBillingUsesSizeTier(t *testing.T) {
	groupID := int64(127)
	defaultPrice := 0.10
	price4K := 0.40
	cache := newEmptyChannelCache()
	cache.pricingByGroupModel[channelModelKey{groupID: groupID, model: "gemini-image"placeholder] = &ChannelModelPricing{
		BillingMode:     BillingModeImage,
		PerRequestPrice: &defaultPrice,
		Intervals: []PricingInterval{{
			TierLabel:       "4K",
			PerRequestPrice: &price4K,
placeholder
placeholder
	cache.channelByGroupID[groupID] = &Channel{ID: groupID, Status: StatusActiveplaceholder
	cache.loadedAt = time.Now()
	channelService := &ChannelService{placeholder
	channelService.cache.Store(cache)

	svc := &GatewayService{
		billingService: NewBillingService(&config.Config{placeholder, nil),
		resolver:       NewModelPricingResolver(channelService, NewBillingService(&config.Config{placeholder, nil)),
placeholder

	cost := svc.calculateRecordUsageCost(
		context.Background(),
		&ForwardResult{Model: "gemini-image", ImageCount: 2, ImageSize: "4K"placeholder,
		&APIKey{GroupID: i64p(groupID), Group: &Group{ID: groupIDplaceholderplaceholder,
		"gemini-image",
		1.0,
		1.0,
		time.Time{placeholder,
		nil,
	)

	require.NotNil(t, cost)
	require.Equal(t, string(BillingModeImage), cost.BillingMode)
	require.InDelta(t, 0.80, cost.TotalCost, 1e-12)
	require.InDelta(t, 0.80, cost.ActualCost, 1e-12)
placeholder

func TestGatewayServiceCalculateRecordUsageCost_GroupImagePriceOverridesChannelImagePrice(t *testing.T) {
	groupID := int64(129)
	channelPrice := 0.25
	groupImagePrice2K := 0.021

	svc := &GatewayService{
		billingService: NewBillingService(&config.Config{placeholder, nil),
		resolver:       newOpenAIImageChannelPricingResolverForTest(t, groupID, "gemini-image", channelPrice),
placeholder

	cost := svc.calculateRecordUsageCost(
		context.Background(),
		&ForwardResult{Model: "gemini-image", ImageCount: 2, ImageSize: ImageBillingSize2Kplaceholder,
		&APIKey{
			GroupID: i64p(groupID),
			Group: &Group{
				ID:           groupID,
				ImagePrice2K: &groupImagePrice2K,
		placeholder,
	placeholder,
		"gemini-image",
		1.0,
		1.0,
		time.Time{placeholder,
		nil,
	)

	require.NotNil(t, cost)
	require.Equal(t, string(BillingModeImage), cost.BillingMode)
	require.InDelta(t, 0.042, cost.TotalCost, 1e-12)
	require.InDelta(t, 0.042, cost.ActualCost, 1e-12)
placeholder

func TestRecordUsageMarksCyberRequestType(t *testing.T) {
	logStub := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userStub := &openAIRecordUsageUserRepoStub{placeholder
	subStub := &openAIRecordUsageSubRepoStub{placeholder
	rateStub := &openAIUserGroupRateRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceForTest(logStub, userStub, subStub, rateStub)

	in := &OpenAIRecordUsageInput{
		CyberBlocked: true,
		Result: &OpenAIForwardResult{
			Model:    "gpt-5",
			Duration: time.Second,
			Usage:    OpenAIUsage{InputTokens: 100, OutputTokens: 0placeholder,
	placeholder,
		APIKey:  &APIKey{ID: 2, Group: &Group{RateMultiplier: 1placeholderplaceholder,
		User:    &User{ID: 1placeholder,
		Account: &Account{ID: 3placeholder,
placeholder
	require.NoError(t, svc.RecordUsage(context.Background(), in))
	require.NotNil(t, logStub.lastLog)
	require.Equal(t, RequestTypeCyberBlocked, logStub.lastLog.RequestType)
	require.Equal(t, 100, logStub.lastLog.InputTokens, "计费 token 不变(正常计费)")
placeholder

func TestGatewayServiceCalculateRecordUsageCost_ChannelImageBillingNormalizesMissingSizeTier(t *testing.T) {
	groupID := int64(128)
	defaultPrice := 0.10
	price2K := 0.22
	cache := newEmptyChannelCache()
	cache.pricingByGroupModel[channelModelKey{groupID: groupID, model: "gemini-image"placeholder] = &ChannelModelPricing{
		BillingMode:     BillingModeImage,
		PerRequestPrice: &defaultPrice,
		Intervals: []PricingInterval{{
			TierLabel:       "2K",
			PerRequestPrice: &price2K,
placeholder
placeholder
	cache.channelByGroupID[groupID] = &Channel{ID: groupID, Status: StatusActiveplaceholder
	cache.loadedAt = time.Now()
	channelService := &ChannelService{placeholder
	channelService.cache.Store(cache)

	svc := &GatewayService{
		billingService: NewBillingService(&config.Config{placeholder, nil),
		resolver:       NewModelPricingResolver(channelService, NewBillingService(&config.Config{placeholder, nil)),
placeholder

	cost := svc.calculateRecordUsageCost(
		context.Background(),
		&ForwardResult{Model: "gemini-image", ImageCount: 2, ImageSize: ""placeholder,
		&APIKey{GroupID: i64p(groupID), Group: &Group{ID: groupIDplaceholderplaceholder,
		"gemini-image",
		1.0,
		1.0,
		time.Time{placeholder,
		nil,
	)

	require.NotNil(t, cost)
	require.Equal(t, string(BillingModeImage), cost.BillingMode)
	require.InDelta(t, 0.44, cost.TotalCost, 1e-12)
	require.InDelta(t, 0.44, cost.ActualCost, 1e-12)
placeholder
