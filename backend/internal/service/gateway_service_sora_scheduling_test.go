package service

import (
	"context"
	"testing"
	"time"
)

func TestGatewayServiceIsAccountSchedulableForSelectionSoraIgnoresGenericWindows(t *testing.T) {
	svc := &GatewayService{placeholder
	now := time.Now()
	past := now.Add(-1 * time.Minute)
	future := now.Add(5 * time.Minute)

	acc := &Account{
		Platform:           PlatformSora,
		Status:             StatusActive,
		Schedulable:        true,
		AutoPauseOnExpired: true,
		ExpiresAt:          &past,
		OverloadUntil:      &future,
		RateLimitResetAt:   &future,
placeholder

	if !svc.isAccountSchedulableForSelection(acc) {
		t.Fatalf("expected sora account to ignore generic expiry/overload/rate-limit windows")
placeholder
placeholder

func TestGatewayServiceIsAccountSchedulableForSelectionNonSoraKeepsGenericLogic(t *testing.T) {
	svc := &GatewayService{placeholder
	future := time.Now().Add(5 * time.Minute)

	acc := &Account{
		Platform:         PlatformAnthropic,
		Status:           StatusActive,
		Schedulable:      true,
		RateLimitResetAt: &future,
placeholder

	if svc.isAccountSchedulableForSelection(acc) {
		t.Fatalf("expected non-sora account to keep generic schedulable checks")
placeholder
placeholder

func TestGatewayServiceIsAccountSchedulableForModelSelectionSoraChecksModelScopeOnly(t *testing.T) {
	svc := &GatewayService{placeholder
	model := "sora2-landscape-10s"
	resetAt := time.Now().Add(2 * time.Minute).UTC().Format(time.RFC3339)
	globalResetAt := time.Now().Add(2 * time.Minute)

	acc := &Account{
		Platform:         PlatformSora,
		Status:           StatusActive,
		Schedulable:      true,
		RateLimitResetAt: &globalResetAt,
		Extra: map[string]any{
			"model_rate_limits": map[string]any{
				model: map[string]any{
					"rate_limit_reset_at": resetAt,
			placeholder,
		placeholder,
	placeholder,
placeholder

	if svc.isAccountSchedulableForModelSelection(context.Background(), acc, model) {
		t.Fatalf("expected sora account to be blocked by model scope rate limit")
placeholder
placeholder

func TestCollectSelectionFailureStatsSoraIgnoresGenericUnschedulableWindows(t *testing.T) {
	svc := &GatewayService{placeholder
	future := time.Now().Add(3 * time.Minute)

	accounts := []Account{
		{
			ID:               1,
			Platform:         PlatformSora,
			Status:           StatusActive,
			Schedulable:      true,
			RateLimitResetAt: &future,
	placeholder,
placeholder

	stats := svc.collectSelectionFailureStats(context.Background(), accounts, "sora2-landscape-10s", PlatformSora, map[int64]struct{placeholder{placeholder, false)
	if stats.Unschedulable != 0 || stats.Eligible != 1 {
		t.Fatalf("unexpected stats: unschedulable=%d eligible=%d", stats.Unschedulable, stats.Eligible)
placeholder
placeholder
