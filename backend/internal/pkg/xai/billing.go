package xai

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// CLI client identity required by cli-chat-proxy billing endpoints.
	CLITokenAuthHeader     = "x-xai-token-auth"
	CLITokenAuthValue      = "xai-grok-cli"
	CLIClientVersionHeader = "x-grok-client-version"
	// CLIClientVersion is the one place the pinned Grok CLI version lives. The
	// repository and service layers build their own client identity from it, so
	// one bump here covers OAuth traffic and billing probes together.
	// Keep in sync with https://x.ai/cli/stable.
	CLIClientVersion = "0.2.114"
	CLIUserAgent     = "grok-pager/" + CLIClientVersion + " grok-shell/" + CLIClientVersion + " (macos; aarch64)"

	BillingWeeklyPath  = "/billing?format=credits"
	BillingMonthlyPath = "/billing"

	SuperGrokLimitCents      = 15_000  // $150.00
	SuperGrokHeavyLimitCents = 150_000 // $1,500.00
)

// BillingPeriod describes the current weekly/monthly window.
type BillingPeriod struct {
	Type  string `json:"type,omitempty"`
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
placeholder

// BillingProductUsage is per-product usage inside the weekly credits window.
type BillingProductUsage struct {
	Product      string   `json:"product,omitempty"`
	UsagePercent *float64 `json:"usagePercent,omitempty"`
placeholder

// BillingConfig is the nested config object from /v1/billing responses.
type BillingConfig struct {
	CurrentPeriod      *BillingPeriod        `json:"currentPeriod,omitempty"`
	CreditUsagePercent *float64              `json:"creditUsagePercent,omitempty"`
	ProductUsage       []BillingProductUsage `json:"productUsage,omitempty"`
	MonthlyLimit       json.RawMessage       `json:"monthlyLimit,omitempty"`
	Used               json.RawMessage       `json:"used,omitempty"`
	BillingPeriodStart string                `json:"billingPeriodStart,omitempty"`
	BillingPeriodEnd   string                `json:"billingPeriodEnd,omitempty"`
placeholder

// BillingPayload is the top-level body from /v1/billing.
type BillingPayload struct {
	Config *BillingConfig `json:"config,omitempty"`
placeholder

// BillingProductSummary is a normalized product usage row for UI.
type BillingProductSummary struct {
	Product      string   `json:"product"`
	UsagePercent *float64 `json:"usage_percent,omitempty"`
placeholder

// BillingSummary is the merged weekly + monthly billing view.
type BillingSummary struct {
	PeriodType         string                  `json:"period_type,omitempty"` // weekly | monthly | unknown
	UsagePercent       *float64                `json:"usage_percent,omitempty"`
	PeriodStart        string                  `json:"period_start,omitempty"`
	PeriodEnd          string                  `json:"period_end,omitempty"`
	ProductUsage       []BillingProductSummary `json:"product_usage,omitempty"`
	MonthlyLimitCents  *float64                `json:"monthly_limit_cents,omitempty"`
	UsedCents          *float64                `json:"used_cents,omitempty"`
	IncludedUsedCents  *float64                `json:"included_used_cents,omitempty"`
	BillingPeriodStart string                  `json:"billing_period_start,omitempty"`
	BillingPeriodEnd   string                  `json:"billing_period_end,omitempty"`
	UsedPercent        *float64                `json:"used_percent,omitempty"`
	Plan               string                  `json:"plan,omitempty"` // SuperGrok | SuperGrok Heavy | ""
	StatusCode         int                     `json:"status_code,omitempty"`
	WeeklyStatusCode   int                     `json:"weekly_status_code,omitempty"`
	MonthlyStatusCode  int                     `json:"monthly_status_code,omitempty"`
	Source             string                  `json:"source,omitempty"`
	FetchedAt          string                  `json:"fetched_at,omitempty"`
	UpdatedAt          string                  `json:"updated_at,omitempty"`
	WeeklyUpdatedAt    string                  `json:"weekly_updated_at,omitempty"`
	MonthlyUpdatedAt   string                  `json:"monthly_updated_at,omitempty"`
	Partial            bool                    `json:"partial,omitempty"`
	FailedWindows      []string                `json:"failed_windows,omitempty"`
placeholder

// BuildBillingURL builds weekly or monthly billing URL against the CLI chat proxy.
func BuildBillingURL(formatCredits bool) string {
	base := strings.TrimRight(DefaultCLIBaseURL, "/")
	if formatCredits {
		return base + BillingWeeklyPath
placeholder
	return base + BillingMonthlyPath
placeholder

// BuildBillingURLWithValidator builds the weekly or monthly billing URL against
// the caller-resolved base URL, applying the caller's outbound URL trust policy
// first. Accounts forwarding through a custom upstream keep their billing
// probes on the same upstream.
func BuildBillingURLWithValidator(baseURL string, formatCredits bool, validator BaseURLValidator) (string, error) {
	validatedBaseURL, err := validatedBaseURLWithValidator(baseURL, validator)
	if err != nil {
		return "", fmt.Errorf("invalid base url: %w", err)
placeholder
	if formatCredits {
		return validatedBaseURL + BillingWeeklyPath, nil
placeholder
	return validatedBaseURL + BillingMonthlyPath, nil
placeholder

// ApplyCLIBillingHeaders sets Authorization + CLI identity headers for billing GETs.
func ApplyCLIBillingHeaders(req *http.Request, accessToken string) {
	if req == nil {
		return
placeholder
	token := strings.TrimSpace(accessToken)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
placeholder
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(CLITokenAuthHeader, CLITokenAuthValue)
	req.Header.Set(CLIClientVersionHeader, CLIClientVersion)
	req.Header.Set("User-Agent", CLIUserAgent)
placeholder

// ParseBillingPayload unmarshals a billing API response body.
func ParseBillingPayload(body []byte) (*BillingPayload, error) {
	if len(body) == 0 {
		return nil, fmt.Errorf("empty billing body")
placeholder
	var payload BillingPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
placeholder
	return &payload, nil
placeholder

// BuildBillingSummary normalizes a billing config into a UI-friendly summary.
func BuildBillingSummary(config *BillingConfig) *BillingSummary {
	if config == nil {
		return nil
placeholder
	summary := &BillingSummary{placeholder
	period := config.CurrentPeriod
	periodType := resolvePeriodType(period)
	creditUsage := cloneFloat(config.CreditUsagePercent)

	periodStart := ""
	periodEnd := ""
	if period != nil {
		periodStart = strings.TrimSpace(period.Start)
		periodEnd = strings.TrimSpace(period.End)
placeholder
	if periodStart == "" {
		periodStart = strings.TrimSpace(config.BillingPeriodStart)
placeholder
	if periodEnd == "" {
		periodEnd = strings.TrimSpace(config.BillingPeriodEnd)
placeholder

	products := make([]BillingProductSummary, 0, len(config.ProductUsage))
	for _, item := range config.ProductUsage {
		product := strings.TrimSpace(item.Product)
		if product == "" {
			continue
	placeholder
		products = append(products, BillingProductSummary{
			Product:      product,
			UsagePercent: cloneFloat(item.UsagePercent),
	placeholder)
placeholder

	monthlyLimit := parseCentValue(config.MonthlyLimit)
	used := parseCentValue(config.Used)
	billingStart := strings.TrimSpace(config.BillingPeriodStart)
	billingEnd := strings.TrimSpace(config.BillingPeriodEnd)

	var includedUsed *float64
	if used != nil {
		if monthlyLimit != nil && *monthlyLimit > 0 {
			v := math.Min(*used, *monthlyLimit)
			includedUsed = &v
	placeholder else {
			includedUsed = cloneFloat(used)
	placeholder
placeholder

	var usedPercent *float64
	if monthlyLimit != nil && *monthlyLimit > 0 && includedUsed != nil {
		v := (*includedUsed / *monthlyLimit) * 100
		usedPercent = &v
placeholder

	hasWeekly := creditUsage != nil || periodType == "weekly" || len(products) > 0
	hasMonthly := monthlyLimit != nil || used != nil || (!hasWeekly && billingEnd != "")
	if !hasWeekly && !hasMonthly {
		return nil
placeholder

	if hasWeekly {
		if periodType == "unknown" {
			periodType = "weekly"
	placeholder
		summary.PeriodType = periodType
		summary.UsagePercent = creditUsage
		summary.PeriodStart = periodStart
		summary.PeriodEnd = periodEnd
placeholder else {
		// Monthly-only: do not put monthly % into UsagePercent (weekly bar field).
		// Frontend weekly bar only renders when PeriodType == weekly.
		summary.PeriodType = "monthly"
		summary.PeriodStart = billingStart
		summary.PeriodEnd = billingEnd
placeholder
	summary.ProductUsage = products
	summary.MonthlyLimitCents = monthlyLimit
	summary.UsedCents = used
	summary.IncludedUsedCents = includedUsed
	if hasMonthly {
		summary.BillingPeriodStart = billingStart
		summary.BillingPeriodEnd = billingEnd
placeholder
	summary.UsedPercent = usedPercent
	summary.Plan = resolvePlan(monthlyLimit)
	return summary
placeholder

// MergeBillingProbeResult updates successful billing domains while retaining
// the previous value for any domain that could not be refreshed.
func MergeBillingProbeResult(previous, weekly, monthly *BillingSummary, weeklyOK, monthlyOK bool) *BillingSummary {
	var out BillingSummary
	if previous != nil {
		out = *previous
		previousUpdatedAt := previous.UpdatedAt
		if previousUpdatedAt == "" {
			previousUpdatedAt = previous.FetchedAt
	placeholder
		if out.WeeklyUpdatedAt == "" && (out.UsagePercent != nil || len(out.ProductUsage) > 0) {
			out.WeeklyUpdatedAt = previousUpdatedAt
	placeholder
		if out.MonthlyUpdatedAt == "" && (out.MonthlyLimitCents != nil || out.UsedPercent != nil) {
			out.MonthlyUpdatedAt = previousUpdatedAt
	placeholder
placeholder
	now := time.Now().UTC().Format(time.RFC3339)

	if weeklyOK && weekly != nil {
		out.PeriodType = weekly.PeriodType
		out.UsagePercent = weekly.UsagePercent
		out.PeriodStart = weekly.PeriodStart
		out.PeriodEnd = weekly.PeriodEnd
		out.ProductUsage = weekly.ProductUsage
		out.WeeklyUpdatedAt = now
placeholder
	if monthlyOK && monthly != nil {
		if out.PeriodType == "" {
			out.PeriodType = "monthly"
	placeholder
		out.MonthlyLimitCents = monthly.MonthlyLimitCents
		out.UsedCents = monthly.UsedCents
		out.IncludedUsedCents = monthly.IncludedUsedCents
		out.BillingPeriodStart = monthly.BillingPeriodStart
		out.BillingPeriodEnd = monthly.BillingPeriodEnd
		out.UsedPercent = monthly.UsedPercent
		out.Plan = monthly.Plan
		out.MonthlyUpdatedAt = now
placeholder

	out.Partial = !weeklyOK || !monthlyOK
	out.FailedWindows = nil
	if !weeklyOK {
		out.FailedWindows = append(out.FailedWindows, "weekly")
placeholder
	if !monthlyOK {
		out.FailedWindows = append(out.FailedWindows, "monthly")
placeholder
	if !weeklyOK && !monthlyOK && previous == nil {
		return nil
placeholder
	return &out
placeholder

// StampBillingSummary sets fetch metadata.
func StampBillingSummary(summary *BillingSummary, statusCode int, source string) *BillingSummary {
	if summary == nil {
		return nil
placeholder
	now := time.Now().UTC().Format(time.RFC3339)
	summary.StatusCode = statusCode
	summary.Source = source
	summary.FetchedAt = now
	summary.UpdatedAt = now
	return summary
placeholder

func resolvePeriodType(period *BillingPeriod) string {
	if period == nil {
		return "unknown"
placeholder
	raw := strings.ToLower(strings.TrimSpace(period.Type))
	if strings.Contains(raw, "weekly") {
		return "weekly"
placeholder
	if strings.Contains(raw, "monthly") {
		return "monthly"
placeholder
	return "unknown"
placeholder

func resolvePlan(monthlyLimitCents *float64) string {
	if monthlyLimitCents == nil {
		return ""
placeholder
	// Allow small float noise.
	limit := math.Round(*monthlyLimitCents)
	switch limit {
	case SuperGrokLimitCents:
		return "SuperGrok"
	case SuperGrokHeavyLimitCents:
		return "SuperGrok Heavy"
	default:
		return ""
placeholder
placeholder

func parseCentValue(raw json.RawMessage) *float64 {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
placeholder
	// Object form: {"val": placeholder
	var obj struct {
		Val any `json:"val"`
placeholder
	if err := json.Unmarshal(raw, &obj); err == nil && obj.Val != nil {
		return anyToFloat(obj.Val)
placeholder
	// Bare number / string
	var n any
	if err := json.Unmarshal(raw, &n); err != nil {
		return nil
placeholder
	return anyToFloat(n)
placeholder

func anyToFloat(v any) *float64 {
	switch n := v.(type) {
	case float64:
		return &n
	case float32:
		f := float64(n)
		return &f
	case int:
		f := float64(n)
		return &f
	case int64:
		f := float64(n)
		return &f
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return nil
	placeholder
		return &f
	case string:
		s := strings.TrimSpace(n)
		if s == "" {
			return nil
	placeholder
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil
	placeholder
		return &f
	default:
		return nil
placeholder
placeholder

func cloneFloat(v *float64) *float64 {
	if v == nil {
		return nil
placeholder
	f := *v
	return &f
placeholder
