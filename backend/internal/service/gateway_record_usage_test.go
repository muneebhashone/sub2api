//go:build unit

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

func newGatewayRecordUsageServiceForTest(usageRepo UsageLogRepository, userRepo UserRepository, subRepo UserSubscriptionRepository) *GatewayService {
	cfg := &config.Config{placeholder
	cfg.Default.RateMultiplier = 1.1
	return NewGatewayService(
		nil,
		nil,
		usageRepo,
		nil,
		userRepo,
		subRepo,
		nil,
		nil,
		cfg,
		nil,
		nil,
		NewBillingService(cfg, nil),
		nil,
		&BillingCacheService{placeholder,
		nil,
		nil,
		&DeferredService{placeholder,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil, // userPlatformQuotaRepo
	)
placeholder

func newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo UsageLogRepository, billingRepo UsageBillingRepository, userRepo UserRepository, subRepo UserSubscriptionRepository) *GatewayService {
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, subRepo)
	svc.usageBillingRepo = billingRepo
	return svc
placeholder

type openAIRecordUsageBestEffortLogRepoStub struct {
	UsageLogRepository

	bestEffortErr   error
	createErr       error
	bestEffortCalls int
	createCalls     int
	lastLog         *UsageLog
	lastCtxErr      error
placeholder

func (s *openAIRecordUsageBestEffortLogRepoStub) CreateBestEffort(ctx context.Context, log *UsageLog) error {
	s.bestEffortCalls++
	s.lastLog = log
	s.lastCtxErr = ctx.Err()
	return s.bestEffortErr
placeholder

func (s *openAIRecordUsageBestEffortLogRepoStub) Create(ctx context.Context, log *UsageLog) (bool, error) {
	s.createCalls++
	s.lastLog = log
	s.lastCtxErr = ctx.Err()
	return false, s.createErr
placeholder

func TestGatewayServiceRecordUsage_BillingUsesDetachedContext(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: false, err: context.DeadlineExceededplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, subRepo)

	reqCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := svc.RecordUsage(reqCtx, &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_detached_ctx",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    501,
			Quota: 100,
	placeholder,
		User:          &User{ID: 601placeholder,
		Account:       &Account{ID: 701placeholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 1, userRepo.deductCalls)
	require.NoError(t, userRepo.lastCtxErr)
	require.Equal(t, 1, quotaSvc.quotaCalls)
	require.NoError(t, quotaSvc.lastQuotaCtxErr)
placeholder

func TestGatewayServiceRecordUsage_BillingFingerprintIncludesRequestPayloadHash(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	payloadHash := HashUsageRequestPayload([]byte(`{"messages":[{"role":"user","content":"hello"placeholder]placeholder`))
	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_payload_hash",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
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

func TestGatewayServiceRecordUsage_BillingFingerprintFallsBackToContextRequestID(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	ctx := context.WithValue(context.Background(), ctxkey.RequestID, "req-local-123")
	err := svc.RecordUsage(ctx, &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_payload_fallback",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 501, Quota: 100placeholder,
		User:    &User{ID: 601placeholder,
		Account: &Account{ID: 701placeholder,
placeholder)
placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.Equal(t, "local:req-local-123", billingRepo.lastCmd.RequestPayloadHash)
placeholder

func TestGatewayServiceRecordUsage_PreservesRequestedAndUpstreamModels(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)
	mappedModel := "claude-sonnet-4-20250514"

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:     "gateway_models_split",
			Usage:         ClaudeUsage{InputTokens: 10, OutputTokens: 6placeholder,
			Model:         "claude-sonnet-4",
			UpstreamModel: mappedModel,
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 501, Quota: 100placeholder,
		User:    &User{ID: 601placeholder,
		Account: &Account{ID: 701placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "claude-sonnet-4", usageRepo.lastLog.Model)
	require.Equal(t, "claude-sonnet-4", usageRepo.lastLog.RequestedModel)
	require.NotNil(t, usageRepo.lastLog.UpstreamModel)
	require.Equal(t, mappedModel, *usageRepo.lastLog.UpstreamModel)
placeholder

func TestGatewayServiceRecordUsage_PreservesChannelMappedUpstreamModel(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:     "gateway_channel_mapping_models",
			Usage:         ClaudeUsage{InputTokens: 10, OutputTokens: 6placeholder,
			Model:         "gpt-5.6-terra",
			UpstreamModel: "gpt-5.6-terra",
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 501, Quota: 100placeholder,
		User:    &User{ID: 601placeholder,
		Account: &Account{ID: 701placeholder,
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

func TestGatewayServiceRecordUsage_PreservesLoopedChannelAndAccountUpstreamModel(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:     "gateway_looped_mapping_models",
			Usage:         ClaudeUsage{InputTokens: 10, OutputTokens: 6placeholder,
			Model:         "gpt-5.6-terra",
			UpstreamModel: "gpt-5.6-sol",
			Duration:      time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 501, Quota: 100placeholder,
		User:    &User{ID: 601placeholder,
		Account: &Account{ID: 701placeholder,
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

func TestGatewayServiceRecordUsage_EmptyImageSizeDefaultsBeforeBillingAndPersistence(t *testing.T) {
	imagePrice2K := 0.19
	groupID := int64(901)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:      "gateway_image_default_size",
			Model:          "gemini-image",
			ImageCount:     1,
			ImageInputSize: "auto",
			Duration:       time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      801,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:             groupID,
				RateMultiplier: 1.0,
				ImagePrice2K:   &imagePrice2K,
		placeholder,
	placeholder,
		User:    &User{ID: 601placeholder,
		Account: &Account{ID: 701placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 1, usageRepo.lastLog.ImageCount)
	require.NotNil(t, usageRepo.lastLog.ImageSize)
	require.Equal(t, ImageBillingSize2K, *usageRepo.lastLog.ImageSize)
	require.NotNil(t, usageRepo.lastLog.ImageInputSize)
	require.Equal(t, "auto", *usageRepo.lastLog.ImageInputSize)
	require.NotNil(t, usageRepo.lastLog.ImageSizeSource)
	require.Equal(t, ImageSizeSourceDefault, *usageRepo.lastLog.ImageSizeSource)
	require.InDelta(t, 0.19, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, 0.19, usageRepo.lastLog.ActualCost, 1e-12)
placeholder

func TestGatewayServiceRecordUsage_PeakRateAffectsTokenModeImageOutputTokens(t *testing.T) {
	groupID := int64(902)
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{placeholder)
	svc.resolver = newOpenAITokenImageChannelPricingResolverForTest(t, groupID, "gemini-image")

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:  "gateway_peak_image_tokens",
			Model:      "gemini-image",
			ImageCount: 1,
			Usage: ClaudeUsage{
				InputTokens:       1000,
				OutputTokens:      600,
				ImageOutputTokens: 100,
		placeholder,
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:      802,
			GroupID: i64p(groupID),
			Group: &Group{
				ID:                 groupID,
				RateMultiplier:     1.0,
				SubscriptionType:   SubscriptionTypeSubscription,
				PeakRateEnabled:    true,
				PeakStart:          "00:00",
				PeakEnd:            "23:59",
				PeakRateMultiplier: 3.0,
		placeholder,
	placeholder,
		User:    &User{ID: 602placeholder,
		Account: &Account{ID: 702placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.BillingMode)
	require.Equal(t, string(BillingModeToken), *usageRepo.lastLog.BillingMode)
	require.Equal(t, 3.0, usageRepo.lastLog.RateMultiplier)

	textInput := 1000 * 3e-6
	textOutput := 500 * 15e-6
	imageOutput := 100 * 15e-6
	expectedActual := (textInput + textOutput + imageOutput) * 3.0

	require.InDelta(t, textInput+textOutput+imageOutput, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, imageOutput, usageRepo.lastLog.ImageOutputCost, 1e-12)
	require.InDelta(t, expectedActual, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, expectedActual, userRepo.lastAmount, 1e-12)
placeholder

func TestGatewayServiceRecordUsage_TimePricingUsesPricingAt(t *testing.T) {
	groupID := int64(904)
	requestStart := time.Date(2024, time.January, 2, 2, 0, 0, 0, time.UTC) // 上海 10:00
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{placeholder)
	svc.resolver = newOpenAITokenImageChannelPricingResolverWithTimeForTest(t, groupID, "gpt-5.1", &ChannelTimePricing{
		Timezone: "Asia/Shanghai",
		Periods:  []ChannelTimePricingPeriod{{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholderplaceholder,
placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_time_pricing_request_start",
			Model:     "gpt-5.1",
			Usage:     ClaudeUsage{InputTokens: 1000, OutputTokens: 500placeholder,
	placeholder,
		APIKey: &APIKey{ID: 804, GroupID: i64p(groupID), Group: &Group{
			ID: groupID, RateMultiplier: 0.8, SubscriptionType: SubscriptionTypeSubscription,
placeholder
		User:      &User{ID: 604placeholder,
		Account:   &Account{ID: 704placeholder,
		PricingAt: requestStart,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	baseCost := 1000*3e-6 + 500*15e-6
	require.InDelta(t, baseCost*2, usageRepo.lastLog.TotalCost, 1e-12)
	require.InDelta(t, baseCost*2*0.8, usageRepo.lastLog.ActualCost, 1e-12)
	require.InDelta(t, 0.8, usageRepo.lastLog.RateMultiplier, 1e-12)
placeholder
func TestGatewayServiceRecordUsage_UsesExplicitPricingAtForPeakRate(t *testing.T) {
	for _, platform := range []string{PlatformAnthropic, PlatformGemini, PlatformGrok, PlatformAntigravityplaceholder {
		t.Run(platform, func(t *testing.T) {
			groupID := int64(903)
			usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
			userRepo := &openAIRecordUsageUserRepoStub{placeholder
			svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, &openAIRecordUsageSubRepoStub{placeholder)
			svc.resolver = newOpenAITokenImageChannelPricingResolverForTest(t, groupID, "gemini-image")

			pricingAt := time.Date(2026, time.January, 1, 0, 30, 0, 0, time.UTC)
			err := svc.RecordUsage(context.Background(), &RecordUsageInput{
				Result: &ForwardResult{
					RequestID:  "gateway_explicit_pricing_at_" + platform,
					Model:      "gemini-image",
					ImageCount: 1,
					Usage: ClaudeUsage{
						InputTokens:       1000,
						OutputTokens:      600,
						ImageOutputTokens: 100,
				placeholder,
			placeholder,
				APIKey: &APIKey{
					ID:      803,
					GroupID: i64p(groupID),
					Group: &Group{
						ID:                 groupID,
						Platform:           platform,
						RateMultiplier:     1.0,
						SubscriptionType:   SubscriptionTypeSubscription,
						PeakRateEnabled:    true,
						PeakStart:          "00:00",
						PeakEnd:            "01:00",
						PeakRateMultiplier: 3.0,
				placeholder,
			placeholder,
				User:      &User{ID: 603placeholder,
				Account:   &Account{ID: 703, Platform: platformplaceholder,
				PricingAt: pricingAt,
		placeholder)

		placeholder
			require.NotNil(t, usageRepo.lastLog)
			require.Equal(t, 3.0, usageRepo.lastLog.RateMultiplier)
	placeholder)
placeholder
placeholder

func TestGatewayServiceRecordUsage_UsageLogWriteErrorDoesNotSkipBilling(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: false, err: MarkUsageLogCreateNotPersisted(context.Canceled)placeholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, subRepo)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_not_persisted",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    503,
			Quota: 100,
	placeholder,
		User:          &User{ID: 603placeholder,
		Account:       &Account{ID: 703placeholder,
		APIKeyService: quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 1, userRepo.deductCalls)
	require.Equal(t, 1, quotaSvc.quotaCalls)
placeholder

func TestGatewayServiceRecordUsageWithLongContext_BillingUsesDetachedContext(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{inserted: false, err: context.DeadlineExceededplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	quotaSvc := &openAIRecordUsageAPIKeyQuotaStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, subRepo)

	reqCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := svc.RecordUsageWithLongContext(reqCtx, &RecordUsageLongContextInput{
		Result: &ForwardResult{
			RequestID: "gateway_long_context_detached_ctx",
			Usage: ClaudeUsage{
				InputTokens:  12,
				OutputTokens: 8,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey: &APIKey{
			ID:    502,
			Quota: 100,
	placeholder,
		User:                  &User{ID: 602placeholder,
		Account:               &Account{ID: 702placeholder,
		LongContextThreshold:  200000,
		LongContextMultiplier: 2,
		APIKeyService:         quotaSvc,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.calls)
	require.Equal(t, 1, userRepo.deductCalls)
	require.NoError(t, userRepo.lastCtxErr)
	require.Equal(t, 1, quotaSvc.quotaCalls)
	require.NoError(t, quotaSvc.lastQuotaCtxErr)
placeholder

func TestGatewayServiceRecordUsage_UsesFallbackRequestIDForUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, userRepo, subRepo)

	ctx := context.WithValue(context.Background(), ctxkey.RequestID, "gateway-local-fallback")
	err := svc.RecordUsage(ctx, &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 504placeholder,
		User:    &User{ID: 604placeholder,
		Account: &Account{ID: 704placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "local:gateway-local-fallback", usageRepo.lastLog.RequestID)
placeholder

func TestGatewayServiceRecordUsage_PrefersClientRequestIDOverUpstreamRequestID(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	ctx := context.WithValue(context.Background(), ctxkey.ClientRequestID, "client-stable-123")
	ctx = context.WithValue(ctx, ctxkey.RequestID, "req-local-ignored")
	err := svc.RecordUsage(ctx, &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "upstream-volatile-456",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 506placeholder,
		User:    &User{ID: 606placeholder,
		Account: &Account{ID: 706placeholder,
placeholder)

placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.Equal(t, "client:client-stable-123", billingRepo.lastCmd.RequestID)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "client:client-stable-123", usageRepo.lastLog.RequestID)
placeholder

func TestGatewayServiceRecordUsage_GeneratesRequestIDWhenAllSourcesMissing(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 507placeholder,
		User:    &User{ID: 607placeholder,
		Account: &Account{ID: 707placeholder,
placeholder)

placeholder
	require.NotNil(t, billingRepo.lastCmd)
	require.True(t, strings.HasPrefix(billingRepo.lastCmd.RequestID, "generated:"))
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, billingRepo.lastCmd.RequestID, usageRepo.lastLog.RequestID)
placeholder

func TestGatewayServiceRecordUsage_DroppedUsageLogFallsBackToSyncCreate(t *testing.T) {
	// 计费成功后 best-effort 写入被丢弃（队列超时）时必须同步兜底，
	// 否则出现“已扣费但无 usage_log”的对账缺口（issue #3656）。
	usageRepo := &openAIRecordUsageBestEffortLogRepoStub{
		bestEffortErr: MarkUsageLogCreateDropped(errors.New("usage log best-effort queue full")),
placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_drop_usage_log",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 508placeholder,
		User:    &User{ID: 608placeholder,
		Account: &Account{ID: 708placeholder,
placeholder)

placeholder
	require.Equal(t, 1, usageRepo.bestEffortCalls)
	require.Equal(t, 1, usageRepo.createCalls)
	// 兜底调用使用的 ctx 必须仍然存活，不能带着已死的 ctx 走过场。
	require.NoError(t, usageRepo.lastCtxErr)
placeholder

func TestGatewayServiceRecordUsage_BillingErrorWritesUnsettledUsageLog(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingErr := errors.New("billing tx failed")
	billingRepo := &openAIRecordUsageBillingRepoStub{err: billingErrplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newGatewayRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "gateway_billing_fail",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 6,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 505placeholder,
		User:    &User{ID: 605placeholder,
		Account: &Account{ID: 705placeholder,
placeholder)

	require.ErrorIs(t, err, billingErr)
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 1, usageRepo.calls)
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, 10, usageRepo.lastLog.InputTokens)
	require.Equal(t, 6, usageRepo.lastLog.OutputTokens)
	require.Greater(t, usageRepo.lastLog.InputCost, 0.0)
	require.Greater(t, usageRepo.lastLog.OutputCost, 0.0)
	require.Greater(t, usageRepo.lastLog.TotalCost, 0.0)
	require.Zero(t, usageRepo.lastLog.ActualCost)
placeholder

func TestGatewayServiceRecordUsage_ReasoningEffortPersisted(t *testing.T) {
	usageRepo := &openAIRecordUsageBestEffortLogRepoStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	effort := "max"
	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "effort_test",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 5,
		placeholder,
			Model:           "claude-opus-4-6",
			Duration:        time.Second,
			ReasoningEffort: &effort,
	placeholder,
		APIKey:  &APIKey{ID: 1placeholder,
		User:    &User{ID: 1placeholder,
		Account: &Account{ID: 1placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.ReasoningEffort)
	require.Equal(t, "max", *usageRepo.lastLog.ReasoningEffort)
placeholder

func TestGatewayServiceRecordUsage_ReasoningEffortNil(t *testing.T) {
	usageRepo := &openAIRecordUsageBestEffortLogRepoStub{placeholder
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)

	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID: "no_effort_test",
			Usage: ClaudeUsage{
				InputTokens:  10,
				OutputTokens: 5,
		placeholder,
			Model:    "claude-sonnet-4",
			Duration: time.Second,
	placeholder,
		APIKey:  &APIKey{ID: 1placeholder,
		User:    &User{ID: 1placeholder,
		Account: &Account{ID: 1placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Nil(t, usageRepo.lastLog.ReasoningEffort)
placeholder

// newGatewayRecordUsageServiceWithResolverForTest mirrors production wiring for
// token billing: a pricing resolver plus a grouped API key select the unified
// billing path, which is the only one that honours the service tier.
func newGatewayRecordUsageServiceWithResolverForTest(usageRepo UsageLogRepository) (*GatewayService, *APIKey) {
	svc := newGatewayRecordUsageServiceForTest(usageRepo, &openAIRecordUsageUserRepoStub{placeholder, &openAIRecordUsageSubRepoStub{placeholder)
	svc.resolver = NewModelPricingResolver(nil, svc.billingService)
	groupID := int64(7)
	return svc, &APIKey{ID: 1, GroupID: &groupID, Group: &Group{ID: groupID, RateMultiplier: 1.0placeholderplaceholder
placeholder

func TestGatewayServiceRecordUsage_FastSpeedDowngradedByUpstreamResponse(t *testing.T) {
	usageRepo := &openAIRecordUsageBestEffortLogRepoStub{placeholder
	svc, apiKey := newGatewayRecordUsageServiceWithResolverForTest(usageRepo)

	tier := "fast"
	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:                   "fast_downgraded_test",
			Usage:                       ClaudeUsage{InputTokens: 100, OutputTokens: 50placeholder,
			Model:                       "claude-opus-5",
			Duration:                    time.Second,
			ServiceTier:                 &tier,
			UpstreamResponseServiceTier: "standard",
	placeholder,
		APIKey:  apiKey,
		User:    &User{ID: 1placeholder,
		Account: &Account{ID: 1, Platform: PlatformAnthropicplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.NotNil(t, usageRepo.lastLog.ServiceTier)
	require.Equal(t, "standard", *usageRepo.lastLog.ServiceTier)

	tokens := UsageTokens{InputTokens: 100, OutputTokens: 50placeholder
	standardCost, err := svc.billingService.CalculateCost("claude-opus-5", tokens, 1.0)
placeholder
	fastCost, err := svc.billingService.CalculateCostWithServiceTier("claude-opus-5", tokens, 1.0, "fast")
placeholder
	require.Greater(t, fastCost.TotalCost, standardCost.TotalCost, "fast mode must carry a premium for the test to be meaningful")
	require.InDelta(t, standardCost.TotalCost, usageRepo.lastLog.TotalCost, 1e-10)
placeholder

func TestGatewayServiceRecordUsage_FastSpeedHonouredKeepsPremium(t *testing.T) {
	usageRepo := &openAIRecordUsageBestEffortLogRepoStub{placeholder
	svc, apiKey := newGatewayRecordUsageServiceWithResolverForTest(usageRepo)

	tier := "fast"
	err := svc.RecordUsage(context.Background(), &RecordUsageInput{
		Result: &ForwardResult{
			RequestID:                   "fast_honoured_test",
			Usage:                       ClaudeUsage{InputTokens: 100, OutputTokens: 50placeholder,
			Model:                       "claude-opus-5",
			Duration:                    time.Second,
			ServiceTier:                 &tier,
			UpstreamResponseServiceTier: "fast",
	placeholder,
		APIKey:  apiKey,
		User:    &User{ID: 1placeholder,
		Account: &Account{ID: 1, Platform: PlatformAnthropicplaceholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)
	require.Equal(t, "fast", *usageRepo.lastLog.ServiceTier)

	fastCost, err := svc.billingService.CalculateCostWithServiceTier("claude-opus-5", UsageTokens{InputTokens: 100, OutputTokens: 50placeholder, 1.0, "fast")
placeholder
	require.InDelta(t, fastCost.TotalCost, usageRepo.lastLog.TotalCost, 1e-10)
placeholder
