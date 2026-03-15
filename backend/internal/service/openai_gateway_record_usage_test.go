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

func newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo UsageLogRepository, billingRepo UsageBillingRepository, userRepo UserRepository, subRepo UserSubscriptionRepository, rateRepo UserGroupRateRepository) *OpenAIGatewayService {
	svc := newOpenAIRecordUsageServiceForTest(usageRepo, userRepo, subRepo, rateRepo)
	svc.usageBillingRepo = billingRepo
	return svc
placeholder

func expectedOpenAICost(t *testing.T, svc *OpenAIGatewayService, model string, usage OpenAIUsage, multiplier float64) *CostBreakdown {
placeholder

	cost, err := svc.billingService.CalculateCost(model, UsageTokens{
		InputTokens:         max(usage.InputTokens-usage.CacheReadInputTokens, 0),
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

func TestOpenAIGatewayServiceRecordUsage_BillingErrorSkipsUsageLogWrite(t *testing.T) {
	usageRepo := &openAIRecordUsageLogRepoStub{placeholder
	billingRepo := &openAIRecordUsageBillingRepoStub{err: errors.New("billing tx failed")placeholder
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

placeholder
	require.Equal(t, 1, billingRepo.calls)
	require.Equal(t, 0, usageRepo.calls)
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

func TestOpenAIGatewayServiceRecordUsage_Gpt54LongContextBillsWholeSession(t *testing.T) {
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
		APIKey:  &APIKey{ID: 1014placeholder,
		User:    &User{ID: 2014placeholder,
		Account: &Account{ID: 3014placeholder,
placeholder)

placeholder
	require.NotNil(t, usageRepo.lastLog)

	expectedInput := 300000 * 2.5e-6 * 2.0
	expectedOutput := 2000 * 15e-6 * 1.5
	require.InDelta(t, expectedInput, usageRepo.lastLog.InputCost, 1e-10)
	require.InDelta(t, expectedOutput, usageRepo.lastLog.OutputCost, 1e-10)
	require.InDelta(t, expectedInput+expectedOutput, usageRepo.lastLog.TotalCost, 1e-10)
	require.InDelta(t, (expectedInput+expectedOutput)*1.1, usageRepo.lastLog.ActualCost, 1e-10)
	require.Equal(t, 1, userRepo.deductCalls)
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

	t.Run("default ignored", func(t *testing.T) {
		require.Nil(t, normalizeOpenAIServiceTier("default"))
placeholder)

	t.Run("invalid ignored", func(t *testing.T) {
		require.Nil(t, normalizeOpenAIServiceTier("turbo"))
placeholder)
placeholder

func TestExtractOpenAIServiceTier(t *testing.T) {
	require.Equal(t, "priority", *extractOpenAIServiceTier(map[string]any{"service_tier": "fast"placeholder))
	require.Equal(t, "flex", *extractOpenAIServiceTier(map[string]any{"service_tier": "flex"placeholder))
	require.Nil(t, extractOpenAIServiceTier(map[string]any{"service_tier": 1placeholder))
	require.Nil(t, extractOpenAIServiceTier(nil))
placeholder

func TestExtractOpenAIServiceTierFromBody(t *testing.T) {
	require.Equal(t, "priority", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"fast"placeholder`)))
	require.Equal(t, "flex", *extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"flex"placeholder`)))
	require.Nil(t, extractOpenAIServiceTierFromBody([]byte(`{"service_tier":"default"placeholder`)))
	require.Nil(t, extractOpenAIServiceTierFromBody(nil))
placeholder

func TestOpenAIGatewayServiceRecordUsage_UsesBillingModelAndMetadataFields(t *testing.T) {
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
	require.Equal(t, "gpt-5.1-codex", usageRepo.lastLog.Model)
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
		APIKey:       &APIKey{ID: 100, GroupID: i64p(88), Group: &Group{ID: 88, SubscriptionType: SubscriptionTypeSubscriptionplaceholderplaceholder,
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
