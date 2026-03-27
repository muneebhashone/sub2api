//go:build unit

package service

import (
	"testing"
)

func applyAntigravitySubscriptionResult(account *Account, result AntigravitySubscriptionResult) (map[string]any, map[string]any) {
	credentials := make(map[string]any)
	for k, v := range account.Credentials {
		credentials[k] = v
placeholder
	credentials["plan_type"] = result.PlanType

	extra := make(map[string]any)
	for k, v := range account.Extra {
		extra[k] = v
placeholder
	if result.SubscriptionStatus != "" {
		extra["subscription_status"] = result.SubscriptionStatus
placeholder else {
		delete(extra, "subscription_status")
placeholder
	if result.SubscriptionError != "" {
		extra["subscription_error"] = result.SubscriptionError
placeholder else {
		delete(extra, "subscription_error")
placeholder
	return credentials, extra
placeholder

func TestApplyAntigravityPrivacyMode_SetsInMemoryExtra(t *testing.T) {
	account := &Account{placeholder

	applyAntigravityPrivacyMode(account, AntigravityPrivacySet)

	if account.Extra == nil {
		t.Fatal("expected account.Extra to be initialized")
placeholder
	if got := account.Extra["privacy_mode"]; got != AntigravityPrivacySet {
		t.Fatalf("expected privacy_mode %q, got %v", AntigravityPrivacySet, got)
placeholder
placeholder

func TestApplyAntigravityPrivacyMode_PreservedBySubscriptionResult(t *testing.T) {
	account := &Account{
placeholder
			"access_token": "token",
	placeholder,
		Extra: map[string]any{
			"existing": "value",
	placeholder,
placeholder
	applyAntigravityPrivacyMode(account, AntigravityPrivacySet)

	_, extra := applyAntigravitySubscriptionResult(account, AntigravitySubscriptionResult{
		PlanType: "Pro",
placeholder)

	if got := extra["privacy_mode"]; got != AntigravityPrivacySet {
		t.Fatalf("expected subscription writeback to keep privacy_mode %q, got %v", AntigravityPrivacySet, got)
placeholder
	if got := extra["existing"]; got != "value" {
		t.Fatalf("expected existing extra fields to be preserved, got %v", got)
placeholder
placeholder
