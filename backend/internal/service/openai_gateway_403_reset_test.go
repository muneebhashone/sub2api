package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type openAI403CounterResetStub struct {
	resetCalls []int64
placeholder

func (s *openAI403CounterResetStub) IncrementOpenAI403Count(context.Context, int64, int) (int64, error) {
	return 0, nil
placeholder

func (s *openAI403CounterResetStub) ResetOpenAI403Count(_ context.Context, accountID int64) error {
	s.resetCalls = append(s.resetCalls, accountID)
	return nil
placeholder

func TestOpenAIGatewayServiceRecordUsage_ResetsOpenAI403CounterForZeroUsage(t *testing.T) {
	counter := &openAI403CounterResetStub{placeholder
	rateLimitSvc := NewRateLimitService(nil, nil, nil, nil, nil)
	rateLimitSvc.SetOpenAI403CounterCache(counter)

	usageRepo := &openAIRecordUsageLogRepoStub{inserted: trueplaceholder
	billingRepo := &openAIRecordUsageBillingRepoStub{result: &UsageBillingApplyResult{Applied: trueplaceholderplaceholder
	userRepo := &openAIRecordUsageUserRepoStub{placeholder
	subRepo := &openAIRecordUsageSubRepoStub{placeholder
	svc := newOpenAIRecordUsageServiceWithBillingRepoForTest(usageRepo, billingRepo, userRepo, subRepo, nil)
	svc.rateLimitService = rateLimitSvc

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result: &OpenAIForwardResult{
			RequestID: "resp_zero_usage_reset_403",
			Model:     "gpt-5.1",
	placeholder,
		APIKey:  &APIKey{ID: 1001, Group: &Group{RateMultiplier: 1placeholderplaceholder,
		User:    &User{ID: 2001placeholder,
		Account: &Account{ID: 777, Platform: PlatformOpenAIplaceholder,
placeholder)

placeholder
	require.Equal(t, []int64{777placeholder, counter.resetCalls)
	require.Equal(t, 1, usageRepo.calls)
placeholder
