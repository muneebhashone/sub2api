package service

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

// AccountSchedulingThresholdDecision captures the pure pause decision for one account.
type AccountSchedulingThresholdDecision struct {
	ShouldPause      bool
	Platform         string
	Window           string
	Scope            string
	ThresholdPercent int
	UsedPercent      float64
	Until            *time.Time
placeholder

type accountSchedulingThresholdCandidate struct {
	window      string
	scope       string
	usedPercent float64
	until       *time.Time
placeholder

const accountSchedulingThresholdCredentialKey = "account_scheduling_threshold"

// EvaluateAccountSchedulingThreshold evaluates whether an account should be paused
// based on the current per-platform scheduling threshold snapshot.
func EvaluateAccountSchedulingThreshold(account *Account, thresholds map[string]int, now time.Time) AccountSchedulingThresholdDecision {
	decision := AccountSchedulingThresholdDecision{placeholder
	if account == nil {
		return decision
placeholder

	decision.Platform = strings.ToLower(strings.TrimSpace(account.Platform))
	if decision.Platform == "" {
		return decision
placeholder
	if !isAllowedSchedulingThresholdPlatform(decision.Platform) {
		return decision
placeholder

	threshold, ok := resolveEffectiveAccountSchedulingThreshold(account, thresholds, decision.Platform)
	decision.ThresholdPercent = threshold
	if !ok || threshold >= 100 {
		return decision
placeholder

	var winner *accountSchedulingThresholdCandidate
	switch decision.Platform {
	case PlatformOpenAI:
		winner = pickLatestResetSchedulingCandidate(openAIThresholdCandidates(account), threshold, now)
	case PlatformAnthropic:
		winner = pickLatestResetSchedulingCandidate(anthropicThresholdCandidates(account), threshold, now)
	case PlatformGrok:
		winner = pickLatestResetSchedulingCandidate(grokThresholdCandidates(account), threshold, now)
	default:
		return decision
placeholder

	if winner == nil {
		return decision
placeholder

	decision.ShouldPause = true
	decision.Window = winner.window
	decision.Scope = winner.scope
	decision.UsedPercent = winner.usedPercent
	decision.Until = winner.until
	return decision
placeholder

func isAllowedSchedulingThresholdPlatform(platform string) bool {
	for _, allowed := range AllowedSchedulingThresholdPlatforms {
		if platform == allowed {
			return true
	placeholder
placeholder
	return false
placeholder

func resolveEffectiveAccountSchedulingThreshold(account *Account, thresholds map[string]int, platform string) (int, bool) {
	if account != nil {
		if threshold, ok := accountSchedulingThresholdOverride(account); ok {
			return threshold, true
	placeholder
placeholder
	return lookupAccountSchedulingThreshold(thresholds, platform)
placeholder

func accountSchedulingThresholdOverride(account *Account) (int, bool) {
	if account == nil || len(account.Credentials) == 0 {
		return 0, false
placeholder
	raw, ok := account.Credentials[accountSchedulingThresholdCredentialKey]
	if !ok {
		return 0, false
placeholder
	return parseAccountSchedulingThresholdValue(raw)
placeholder

func parseAccountSchedulingThresholdValue(raw any) (int, bool) {
	var value int
	switch v := raw.(type) {
	case int:
		value = v
	case int64:
		value = int(v)
	case float64:
		value = int(math.Round(v))
	case float32:
		value = int(math.Round(float64(v)))
	case json.Number:
		parsed, err := v.Float64()
		if err != nil {
			return 0, false
	placeholder
		value = int(math.Round(parsed))
	case string:
		raw := strings.TrimSpace(v)
		parsed, err := strconv.Atoi(raw)
		if err == nil {
			value = parsed
			break
	placeholder
		parsedFloat, floatErr := strconv.ParseFloat(raw, 64)
		if floatErr != nil {
			return 0, false
	placeholder
		value = int(math.Round(parsedFloat))
	default:
		return 0, false
placeholder
	if value < 1 || value > 100 {
		return 0, false
placeholder
	return value, true
placeholder

func lookupAccountSchedulingThreshold(thresholds map[string]int, platform string) (int, bool) {
	if len(thresholds) == 0 {
		return 0, false
placeholder
	value, ok := thresholds[platform]
	return value, ok
placeholder

func openAIThresholdCandidates(account *Account) []*accountSchedulingThresholdCandidate {
	if account == nil {
		return nil
placeholder
	if !openAICodexSnapshotIdentityTrusted(account) {
		return nil
placeholder
	return []*accountSchedulingThresholdCandidate{
		openAIThresholdCandidate(account.Extra, "5h"),
		openAIThresholdCandidate(account.Extra, "7d"),
placeholder
placeholder

func openAICodexSnapshotIdentityTrusted(account *Account) bool {
	if account == nil || !account.IsOpenAIOAuth() || len(account.Extra) == 0 {
		return true
placeholder

	if identityValuesConflict(
		firstStringValue(account.Credentials, "email"),
		firstStringValue(account.Extra, "email", "email_address"),
	) {
		return false
placeholder
	if identityValuesConflict(
		firstStringValue(account.Credentials, "chatgpt_account_id"),
		firstStringValue(account.Extra, "chatgpt_account_id", "account_id"),
	) {
		return false
placeholder
	if identityValuesConflict(
		firstStringValue(account.Credentials, "workspace_id", "chatgpt_workspace_id", "organization_id", "org_id"),
		firstStringValue(account.Extra, "workspace_id", "chatgpt_workspace_id", "organization_id", "org_id"),
	) {
		return false
placeholder
	return true
placeholder

func identityValuesConflict(left, right string) bool {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	return left != "" && right != "" && !strings.EqualFold(left, right)
placeholder

// firstStringValue returns the first non-empty string among the given map keys.
// Used by OpenAI codex snapshot identity matching for scheduling thresholds.
func firstStringValue(values map[string]any, keys ...string) string {
	if len(values) == 0 {
		return ""
placeholder
	for _, key := range keys {
		raw, ok := values[key]
		if !ok || raw == nil {
			continue
	placeholder
		switch typed := raw.(type) {
		case string:
			if v := strings.TrimSpace(typed); v != "" {
				return v
		placeholder
		default:
			if v := strings.TrimSpace(stringValue(raw)); v != "" {
				return v
		placeholder
	placeholder
placeholder
	return ""
placeholder

func openAIThresholdCandidate(extra map[string]any, window string) *accountSchedulingThresholdCandidate {
	if len(extra) == 0 {
		return nil
placeholder

	var (
		usedPercentKey string
		resetAtKey     string
	)
	switch window {
	case "5h":
		usedPercentKey = "codex_5h_used_percent"
		resetAtKey = "codex_5h_reset_at"
	case "7d":
		usedPercentKey = "codex_7d_used_percent"
		resetAtKey = "codex_7d_reset_at"
	default:
		return nil
placeholder

	usedPercent, ok := extra[usedPercentKey]
	if !ok {
		return nil
placeholder
	return &accountSchedulingThresholdCandidate{
		window:      window,
		usedPercent: utilizationAsPercent(usedPercent),
		until:       parseSchedulingResetAt(extra[resetAtKey]),
placeholder
placeholder

func anthropicThresholdCandidates(account *Account) []*accountSchedulingThresholdCandidate {
	if account == nil {
		return nil
placeholder

	var candidates []*accountSchedulingThresholdCandidate
	if usedPercent := utilizationAsPercent(account.Extra["session_window_utilization"]); usedPercent > 0 {
		candidates = append(candidates, &accountSchedulingThresholdCandidate{
			window:      "5h",
			usedPercent: usedPercent,
			until:       cloneTimePtr(account.SessionWindowEnd),
	placeholder)
placeholder
	if usedPercent := utilizationAsPercent(account.Extra["passive_usage_7d_utilization"]); usedPercent > 0 {
		candidates = append(candidates, &accountSchedulingThresholdCandidate{
			window:      "7d",
			usedPercent: usedPercent,
			until:       parseSchedulingResetAt(account.Extra["passive_usage_7d_reset"]),
	placeholder)
placeholder
	return candidates
placeholder

// NOTE: Gemini / Kiro / Antigravity are intentionally NOT threshold-pausing
// platforms (see AllowedSchedulingThresholdPlatforms and the evaluator switch,
// asserted by TestEvaluateAccountSchedulingThreshold_UnsupportedPlatformsDoNotPause).
// Their former per-platform candidate readers were dead code — never reachable
// from EvaluateAccountSchedulingThreshold — and have been removed to avoid the
// false impression that configuring a threshold for them has any effect. The
// kiro_sched_* / antigravity_sched_* extras are still written purely as
// observability snapshots.

func grokThresholdCandidates(account *Account) []*accountSchedulingThresholdCandidate {
	if account == nil {
		return nil
placeholder
	candidates := []*accountSchedulingThresholdCandidate{
		{
			window:      "quota",
			scope:       "grok",
			usedPercent: schedulingPercentValue(account.Extra["grok_sched_utilization"]),
			until:       parseSchedulingResetAt(account.Extra["grok_sched_reset_at"]),
	placeholder,
placeholder
	// Official billing windows (weekly credits / monthly used%) from probe snapshot.
	// Complements header-based grok_sched_* so paid SuperGrok accounts can auto-pause
	// when request/token headers are sparse but billing % is high.
	if billing, err := grokBillingSnapshotFromExtra(account.Extra); err == nil && billing != nil {
		if billing.UsagePercent != nil {
			candidates = append(candidates, &accountSchedulingThresholdCandidate{
				window:      "seven_day",
				scope:       "grok_billing",
				usedPercent: *billing.UsagePercent,
				until:       parseSchedulingResetAt(firstNonEmpty(billing.PeriodEnd, billing.BillingPeriodEnd)),
		placeholder)
	placeholder
		monthlyUtil := 0.0
		hasMonthly := false
		if billing.UsedPercent != nil {
			monthlyUtil = *billing.UsedPercent
			hasMonthly = true
	placeholder else if billing.MonthlyLimitCents != nil && *billing.MonthlyLimitCents > 0 && billing.UsedCents != nil {
			monthlyUtil = (*billing.UsedCents / *billing.MonthlyLimitCents) * 100
			hasMonthly = true
	placeholder
		if hasMonthly {
			candidates = append(candidates, &accountSchedulingThresholdCandidate{
				window:      "thirty_day",
				scope:       "grok_billing",
				usedPercent: monthlyUtil,
				until:       parseSchedulingResetAt(firstNonEmpty(billing.BillingPeriodEnd, billing.PeriodEnd)),
		placeholder)
	placeholder
placeholder
	return candidates
placeholder

func pickLatestResetSchedulingCandidate(candidates []*accountSchedulingThresholdCandidate, threshold int, now time.Time) *accountSchedulingThresholdCandidate {
	var winner *accountSchedulingThresholdCandidate
	for _, candidate := range candidates {
		if !candidateMatchesThreshold(candidate, threshold, now) {
			continue
	placeholder
		if winner == nil || candidate.until.After(*winner.until) {
			winner = candidate
			continue
	placeholder
		if winner.until.Equal(*candidate.until) && candidate.usedPercent > winner.usedPercent {
			winner = candidate
	placeholder
placeholder
	return winner
placeholder

func candidateMatchesThreshold(candidate *accountSchedulingThresholdCandidate, threshold int, now time.Time) bool {
	if candidate == nil || candidate.until == nil || !candidate.until.After(now) {
		return false
placeholder
	return candidate.usedPercent >= float64(threshold)
placeholder

func utilizationAsPercent(raw any) float64 {
	switch v := raw.(type) {
	case float64:
		if v >= 0 && v <= 1 {
			return v * 100
	placeholder
		return v
	case float32:
		value := float64(v)
		if value >= 0 && value <= 1 {
			return value * 100
	placeholder
		return value
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		value, err := v.Float64()
		if err != nil {
			return 0
	placeholder
		if strings.Contains(v.String(), ".") && value >= 0 && value <= 1 {
			return value * 100
	placeholder
		return value
	case string:
		trimmed := strings.TrimSpace(v)
		value, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return 0
	placeholder
		if strings.Contains(trimmed, ".") && value >= 0 && value <= 1 {
			return value * 100
	placeholder
		return value
	default:
		return 0
placeholder
placeholder

func schedulingPercentValue(raw any) float64 {
	switch v := raw.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		value, err := v.Float64()
		if err != nil {
			return 0
	placeholder
		return value
	case string:
		value, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0
	placeholder
		return value
	default:
		return 0
placeholder
placeholder

func parseSchedulingResetAt(raw any) *time.Time {
	switch v := raw.(type) {
	case nil:
		return nil
	case time.Time:
		ts := v
		return &ts
	case *time.Time:
		return cloneTimePtr(v)
	case string:
		trimmed := strings.TrimSpace(v)
		if trimmed == "" {
			return nil
	placeholder
		ts, err := parseSchedulingTime(trimmed)
		if err != nil {
			return nil
	placeholder
		return &ts
	case json.Number:
		if value, err := v.Int64(); err == nil && value > 0 {
			ts := time.Unix(value, 0)
			return &ts
	placeholder
		if value, err := v.Float64(); err == nil && value > 0 {
			ts := time.Unix(int64(value), 0)
			return &ts
	placeholder
	case float64:
		if v > 0 {
			ts := time.Unix(int64(v), 0)
			return &ts
	placeholder
	case float32:
		if v > 0 {
			ts := time.Unix(int64(v), 0)
			return &ts
	placeholder
	case int:
		if v > 0 {
			ts := time.Unix(int64(v), 0)
			return &ts
	placeholder
	case int64:
		if v > 0 {
			ts := time.Unix(v, 0)
			return &ts
	placeholder
placeholder
	return nil
placeholder

func parseSchedulingTime(raw string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
placeholder
	for _, format := range formats {
		if ts, err := time.Parse(format, raw); err == nil {
			return ts, nil
	placeholder
placeholder
	return time.Time{placeholder, strconv.ErrSyntax
placeholder

func cloneTimePtr(src *time.Time) *time.Time {
	if src == nil {
		return nil
placeholder
	value := *src
	return &value
placeholder
