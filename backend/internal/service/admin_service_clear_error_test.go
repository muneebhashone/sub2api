//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type accountRepoStubForClearAccountError struct {
	mockAccountRepoForGemini
	account                  *Account
	clearErrorCalls          int
	clearRateLimitCalls      int
	clearAntigravityCalls    int
	clearModelRateLimitCalls int
	clearTempUnschedCalls    int
placeholder

func (r *accountRepoStubForClearAccountError) GetByID(ctx context.Context, id int64) (*Account, error) {
	return r.account, nil
placeholder

func (r *accountRepoStubForClearAccountError) ClearError(ctx context.Context, id int64) error {
	r.clearErrorCalls++
	r.account.Status = StatusActive
	r.account.ErrorMessage = ""
	return nil
placeholder

func (r *accountRepoStubForClearAccountError) ClearRateLimit(ctx context.Context, id int64) error {
	r.clearRateLimitCalls++
	r.account.RateLimitedAt = nil
	r.account.RateLimitResetAt = nil
	return nil
placeholder

func (r *accountRepoStubForClearAccountError) ClearAntigravityQuotaScopes(ctx context.Context, id int64) error {
	r.clearAntigravityCalls++
	return nil
placeholder

func (r *accountRepoStubForClearAccountError) ClearModelRateLimits(ctx context.Context, id int64) error {
	r.clearModelRateLimitCalls++
	return nil
placeholder

func (r *accountRepoStubForClearAccountError) ClearTempUnschedulable(ctx context.Context, id int64) error {
	r.clearTempUnschedCalls++
	r.account.TempUnschedulableUntil = nil
	r.account.TempUnschedulableReason = ""
	return nil
placeholder

func TestAdminService_ClearAccountError_AlsoClearsRecoverableRuntimeState(t *testing.T) {
	until := time.Now().Add(10 * time.Minute)
	resetAt := time.Now().Add(5 * time.Minute)
	repo := &accountRepoStubForClearAccountError{
		account: &Account{
			ID:                      31,
			Platform:                PlatformOpenAI,
			Type:                    AccountTypeOAuth,
			Status:                  StatusError,
			ErrorMessage:            "refresh failed",
			RateLimitResetAt:        &resetAt,
			TempUnschedulableUntil:  &until,
			TempUnschedulableReason: "missing refresh token",
	placeholder,
placeholder
	svc := &adminServiceImpl{accountRepo: repoplaceholder

	updated, err := svc.ClearAccountError(context.Background(), 31)
placeholder
	require.NotNil(t, updated)
	require.Equal(t, 1, repo.clearErrorCalls)
	require.Equal(t, 1, repo.clearRateLimitCalls)
	require.Equal(t, 1, repo.clearAntigravityCalls)
	require.Equal(t, 1, repo.clearModelRateLimitCalls)
	require.Equal(t, 1, repo.clearTempUnschedCalls)
	require.Nil(t, updated.RateLimitResetAt)
	require.Nil(t, updated.TempUnschedulableUntil)
	require.Empty(t, updated.TempUnschedulableReason)
placeholder
