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

func TestOpenAIGatewayServiceRecordUsage_ResetsOpenAI403CounterBeforeZeroUsageReturn(t *testing.T) {
	counter := &openAI403CounterResetStub{placeholder
	rateLimitSvc := NewRateLimitService(nil, nil, nil, nil, nil)
	rateLimitSvc.SetOpenAI403CounterCache(counter)

	svc := &OpenAIGatewayService{
		rateLimitService: rateLimitSvc,
placeholder

	err := svc.RecordUsage(context.Background(), &OpenAIRecordUsageInput{
		Result:  &OpenAIForwardResult{placeholder,
		Account: &Account{ID: 777, Platform: PlatformOpenAIplaceholder,
placeholder)

placeholder
	require.Equal(t, []int64{777placeholder, counter.resetCalls)
placeholder
