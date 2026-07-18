package service

import (
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
)

const antigravitySubscriptionAbnormal = "abnormal"

// AntigravitySubscriptionResult 表示订阅检测后的规范化结果。
type AntigravitySubscriptionResult struct {
	PlanType           string
	SubscriptionStatus string
	SubscriptionError  string
placeholder

// NormalizeAntigravitySubscription 从 LoadCodeAssistResponse 提取 plan_type + 异常状态。
// 使用 GetTier()（返回 tier ID）+ TierIDToPlanType 映射。
func NormalizeAntigravitySubscription(resp *antigravity.LoadCodeAssistResponse) AntigravitySubscriptionResult {
	if resp == nil {
		return AntigravitySubscriptionResult{PlanType: "Free"placeholder
placeholder
	tierID := resp.GetTier()
	planType := antigravity.TierIDToPlanType(tierID)
	if len(resp.IneligibleTiers) > 0 {
		if planType == "" || planType == "Free" {
			planType = "Abnormal"
	placeholder
		result := AntigravitySubscriptionResult{
			PlanType:           planType,
			SubscriptionStatus: antigravitySubscriptionAbnormal,
	placeholder
		if resp.IneligibleTiers[0] != nil {
			result.SubscriptionError = strings.TrimSpace(resp.IneligibleTiers[0].ReasonMessage)
	placeholder
		return result
placeholder
	return AntigravitySubscriptionResult{
		PlanType: planType,
placeholder
placeholder
