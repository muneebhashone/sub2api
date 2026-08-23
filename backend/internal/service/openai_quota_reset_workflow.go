package service

import (
	"context"
	"log/slog"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	OpenAIQuotaResetWarningCacheRefreshFailed    = "reset_credit_cache_refresh_failed"
	OpenAIQuotaResetWarningAccountRecoveryFailed = "account_state_recovery_failed"
	OpenAIQuotaResetWarningAccountRefreshFailed  = "account_state_refresh_failed"
)

type openAIQuotaResetWorkflowQuota interface {
	QueryUsage(ctx context.Context, accountID int64) (*OpenAIQuotaUsage, error)
	CacheResetCreditsSnapshot(ctx context.Context, accountID int64, credits *OpenAIRateLimitResetCredits) error
placeholder

type openAIQuotaResetWorkflowRecoverer interface {
	RecoverAccountState(ctx context.Context, accountID int64, options AccountRecoveryOptions) (*SuccessfulTestRecoveryResult, error)
placeholder

// OpenAIQuotaResetPostProcessResult 汇总手动和自动用卡后的共享恢复结果。
// WarningCode 允许上游消费成功但本地恢复部分失败时仍然准确呈现状态。
type OpenAIQuotaResetPostProcessResult struct {
	Quota                 *OpenAIQuotaUsage
	Account               *Account
	CacheRefreshed        bool
	AccountStateRecovered bool
	WarningCode           string
placeholder

// RunOpenAIQuotaResetPostProcess 按“解除限流、刷新额度缓存、刷新账号行”的固定顺序
// 执行消费后恢复，避免手动和自动入口在失败语义上逐渐分叉。
func RunOpenAIQuotaResetPostProcess(
	ctx context.Context,
	accountID int64,
	quota openAIQuotaResetWorkflowQuota,
	recoverer openAIQuotaResetWorkflowRecoverer,
	loadAccount func(context.Context, int64) (*Account, error),
) OpenAIQuotaResetPostProcessResult {
	result := OpenAIQuotaResetPostProcessResult{placeholder
	if recoverer == nil {
		result.WarningCode = OpenAIQuotaResetWarningAccountRecoveryFailed
		return result
placeholder
	if _, err := recoverer.RecoverAccountState(ctx, accountID, AccountRecoveryOptions{InvalidateToken: trueplaceholder); err != nil {
		slog.Warn("openai_quota_reset_account_recovery_failed", "account_id", accountID, "error_code", infraerrors.Reason(err))
		result.WarningCode = OpenAIQuotaResetWarningAccountRecoveryFailed
		return result
placeholder
	result.AccountStateRecovered = true

	if quota != nil {
		usage, usageErr := quota.QueryUsage(ctx, accountID)
		switch {
		case usageErr != nil || usage == nil:
			slog.Warn("openai_quota_reset_cache_refresh_failed", "account_id", accountID, "error_code", infraerrors.Reason(usageErr))
			result.WarningCode = OpenAIQuotaResetWarningCacheRefreshFailed
		default:
			if err := quota.CacheResetCreditsSnapshot(ctx, accountID, usage.RateLimitResetCredits); err != nil {
				slog.Warn("openai_quota_reset_cache_refresh_failed", "account_id", accountID, "error_code", infraerrors.Reason(err))
				result.WarningCode = OpenAIQuotaResetWarningCacheRefreshFailed
		placeholder else {
				result.Quota = usage
				result.CacheRefreshed = true
		placeholder
	placeholder
placeholder

	if loadAccount == nil {
		return result
placeholder
	account, err := loadAccount(ctx, accountID)
	if err != nil {
		slog.Warn("openai_quota_reset_account_refresh_failed", "account_id", accountID, "error_code", infraerrors.Reason(err))
		if result.WarningCode == "" {
			result.WarningCode = OpenAIQuotaResetWarningAccountRefreshFailed
	placeholder
		return result
placeholder
	result.Account = account
	return result
placeholder
