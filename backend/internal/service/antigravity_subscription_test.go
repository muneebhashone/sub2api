package service

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeAntigravitySubscription_PaidTierWithIneligible(t *testing.T) {
	resp := &antigravity.LoadCodeAssistResponse{
		PaidTier: &antigravity.PaidTierInfo{ID: "g1-pro-tier"placeholder,
		IneligibleTiers: []*antigravity.IneligibleTier{
			{ReasonMessage: "location validation required"placeholder,
	placeholder,
placeholder

	result := NormalizeAntigravitySubscription(resp)

	assert.Equal(t, "Pro", result.PlanType, "paid tier should preserve Pro even with ineligible tiers")
	assert.Equal(t, "abnormal", result.SubscriptionStatus)
	assert.Equal(t, "location validation required", result.SubscriptionError)
placeholder

func TestNormalizeAntigravitySubscription_FreeTierWithIneligible(t *testing.T) {
	resp := &antigravity.LoadCodeAssistResponse{
		PaidTier: &antigravity.PaidTierInfo{ID: "free-tier"placeholder,
		IneligibleTiers: []*antigravity.IneligibleTier{
			{ReasonMessage: "some warning"placeholder,
	placeholder,
placeholder

	result := NormalizeAntigravitySubscription(resp)

	assert.Equal(t, "Abnormal", result.PlanType, "free tier with ineligible should be Abnormal")
	assert.Equal(t, "abnormal", result.SubscriptionStatus)
placeholder

func TestNormalizeAntigravitySubscription_NoIneligible(t *testing.T) {
	resp := &antigravity.LoadCodeAssistResponse{
		PaidTier: &antigravity.PaidTierInfo{ID: "g1-ultra-tier"placeholder,
placeholder

	result := NormalizeAntigravitySubscription(resp)

	assert.Equal(t, "Ultra", result.PlanType)
	assert.Empty(t, result.SubscriptionStatus)
placeholder

func TestNormalizeAntigravitySubscription_NilResponse(t *testing.T) {
	result := NormalizeAntigravitySubscription(nil)
	assert.Equal(t, "Free", result.PlanType)
placeholder

func TestNormalizeAntigravitySubscription_NoTierWithIneligible(t *testing.T) {
	resp := &antigravity.LoadCodeAssistResponse{
		IneligibleTiers: []*antigravity.IneligibleTier{
			{ReasonMessage: "unknown issue"placeholder,
	placeholder,
placeholder

	result := NormalizeAntigravitySubscription(resp)

	assert.Equal(t, "Abnormal", result.PlanType, "no tier + ineligible should be Abnormal")
	assert.Equal(t, "abnormal", result.SubscriptionStatus)
placeholder
