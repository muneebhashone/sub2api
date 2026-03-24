package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

// PrivacyClientFactory creates an HTTP client for privacy API calls.
// Injected from repository layer to avoid import cycles.
type PrivacyClientFactory func(proxyURL string) (*req.Client, error)

const (
	openAISettingsURL = "https://chatgpt.com/backend-api/settings/account_user_setting"

	PrivacyModeTrainingOff = "training_off"
	PrivacyModeFailed      = "training_set_failed"
	PrivacyModeCFBlocked   = "training_set_cf_blocked"
)

// disableOpenAITraining calls ChatGPT settings API to turn off "Improve the model for everyone".
// Returns privacy_mode value: "training_off" on success, "cf_blocked" / "failed" on failure.
func disableOpenAITraining(ctx context.Context, clientFactory PrivacyClientFactory, accessToken, proxyURL string) string {
	if accessToken == "" || clientFactory == nil {
		return ""
placeholder

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client, err := clientFactory(proxyURL)
	if err != nil {
		slog.Warn("openai_privacy_client_error", "error", err.Error())
		return PrivacyModeFailed
placeholder

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Origin", "https://chatgpt.com").
		SetHeader("Referer", "https://chatgpt.com/").
		SetQueryParam("feature", "training_allowed").
		SetQueryParam("value", "false").
		Patch(openAISettingsURL)

	if err != nil {
		slog.Warn("openai_privacy_request_error", "error", err.Error())
		return PrivacyModeFailed
placeholder

	if resp.StatusCode == 403 || resp.StatusCode == 503 {
		body := resp.String()
		if strings.Contains(body, "cloudflare") || strings.Contains(body, "cf-") || strings.Contains(body, "Just a moment") {
			slog.Warn("openai_privacy_cf_blocked", "status", resp.StatusCode)
			return PrivacyModeCFBlocked
	placeholder
placeholder

	if !resp.IsSuccessState() {
		slog.Warn("openai_privacy_failed", "status", resp.StatusCode, "body", truncate(resp.String(), 200))
		return PrivacyModeFailed
placeholder

	slog.Info("openai_privacy_training_disabled")
	return PrivacyModeTrainingOff
placeholder

// ChatGPTAccountInfo 从 chatgpt.com/backend-api/accounts/check 获取的账号信息
type ChatGPTAccountInfo struct {
	PlanType string
	Email    string
placeholder

const chatGPTAccountsCheckURL = "https://chatgpt.com/backend-api/accounts/check/v4-2023-04-27"

// fetchChatGPTAccountInfo calls ChatGPT backend-api to get account info (plan_type, etc.).
// Used as fallback when id_token doesn't contain these fields (e.g., Mobile RT).
// orgID is used to match the correct account when multiple accounts exist (e.g., personal + team).
// Returns nil on any failure (best-effort, non-blocking).
func fetchChatGPTAccountInfo(ctx context.Context, clientFactory PrivacyClientFactory, accessToken, proxyURL, orgID string) *ChatGPTAccountInfo {
	if accessToken == "" || clientFactory == nil {
		return nil
placeholder

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	client, err := clientFactory(proxyURL)
	if err != nil {
		slog.Debug("chatgpt_account_check_client_error", "error", err.Error())
		return nil
placeholder

	var result map[string]any
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Origin", "https://chatgpt.com").
		SetHeader("Referer", "https://chatgpt.com/").
		SetHeader("Accept", "application/json").
		SetSuccessResult(&result).
		Get(chatGPTAccountsCheckURL)

	if err != nil {
		slog.Debug("chatgpt_account_check_request_error", "error", err.Error())
		return nil
placeholder

	if !resp.IsSuccessState() {
		slog.Debug("chatgpt_account_check_failed", "status", resp.StatusCode, "body", truncate(resp.String(), 200))
		return nil
placeholder

	info := &ChatGPTAccountInfo{placeholder

	accounts, ok := result["accounts"].(map[string]any)
	if !ok {
		slog.Debug("chatgpt_account_check_no_accounts", "body", truncate(resp.String(), 300))
		return nil
placeholder

	// 优先匹配 orgID 对应的账号（access_token JWT 中的 poid）
	if orgID != "" {
		if matched := extractPlanFromAccount(accounts, orgID); matched != "" {
			info.PlanType = matched
	placeholder
placeholder

	// 未匹配到时，遍历所有账号：优先 is_default，次选非 free
	if info.PlanType == "" {
		var defaultPlan, paidPlan, anyPlan string
		for _, acctRaw := range accounts {
			acct, ok := acctRaw.(map[string]any)
			if !ok {
				continue
		placeholder
			planType := extractPlanType(acct)
			if planType == "" {
				continue
		placeholder
			if anyPlan == "" {
				anyPlan = planType
		placeholder
			if account, ok := acct["account"].(map[string]any); ok {
				if isDefault, _ := account["is_default"].(bool); isDefault {
					defaultPlan = planType
			placeholder
		placeholder
			if !strings.EqualFold(planType, "free") && paidPlan == "" {
				paidPlan = planType
		placeholder
	placeholder
		// 优先级：default > 非 free > 任意
		switch {
		case defaultPlan != "":
			info.PlanType = defaultPlan
		case paidPlan != "":
			info.PlanType = paidPlan
		default:
			info.PlanType = anyPlan
	placeholder
placeholder

	if info.PlanType == "" {
		slog.Debug("chatgpt_account_check_no_plan_type", "body", truncate(resp.String(), 300))
		return nil
placeholder

	slog.Info("chatgpt_account_check_success", "plan_type", info.PlanType, "org_id", orgID)
	return info
placeholder

// extractPlanFromAccount 从 accounts map 中按 key（account_id）精确匹配并提取 plan_type
func extractPlanFromAccount(accounts map[string]any, accountKey string) string {
	acctRaw, ok := accounts[accountKey]
	if !ok {
		return ""
placeholder
	acct, ok := acctRaw.(map[string]any)
	if !ok {
		return ""
placeholder
	return extractPlanType(acct)
placeholder

// extractPlanType 从单个 account 对象中提取 plan_type
func extractPlanType(acct map[string]any) string {
	if account, ok := acct["account"].(map[string]any); ok {
		if planType, ok := account["plan_type"].(string); ok && planType != "" {
			return planType
	placeholder
placeholder
	if entitlement, ok := acct["entitlement"].(map[string]any); ok {
		if subPlan, ok := entitlement["subscription_plan"].(string); ok && subPlan != "" {
			return subPlan
	placeholder
placeholder
	return ""
placeholder

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
placeholder
	return s[:n] + fmt.Sprintf("...(%d more)", len(s)-n)
placeholder
