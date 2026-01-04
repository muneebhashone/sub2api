package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
)

type geminiModelClass string

const (
	geminiModelPro   geminiModelClass = "pro"
	geminiModelFlash geminiModelClass = "flash"
)

type GeminiQuota struct {
	// SharedRPD is a shared requests-per-day pool across models.
	// When SharedRPD > 0, callers should treat ProRPD/FlashRPD as not applicable for daily quota checks.
	SharedRPD int64 `json:"shared_rpd,omitempty"`
	// SharedRPM is a shared requests-per-minute pool across models.
	// When SharedRPM > 0, callers should treat ProRPM/FlashRPM as not applicable for minute quota checks.
	SharedRPM int64 `json:"shared_rpm,omitempty"`

	// Per-model quotas (AI Studio / API key).
	// A value of -1 means "unlimited" (pay-as-you-go).
	ProRPD   int64 `json:"pro_rpd,omitempty"`
	ProRPM   int64 `json:"pro_rpm,omitempty"`
	FlashRPD int64 `json:"flash_rpd,omitempty"`
	FlashRPM int64 `json:"flash_rpm,omitempty"`
placeholder

type GeminiTierPolicy struct {
	Quota    GeminiQuota
	Cooldown time.Duration
placeholder

type GeminiQuotaPolicy struct {
	tiers map[string]GeminiTierPolicy
placeholder

type GeminiUsageTotals struct {
	ProRequests   int64
	FlashRequests int64
	ProTokens     int64
	FlashTokens   int64
	ProCost       float64
	FlashCost     float64
placeholder

const geminiQuotaCacheTTL = time.Minute

type geminiQuotaOverridesV1 struct {
	Tiers map[string]config.GeminiTierQuotaConfig `json:"tiers"`
placeholder

type geminiQuotaOverridesV2 struct {
	QuotaRules map[string]geminiQuotaRuleOverride `json:"quota_rules"`
placeholder

type geminiQuotaRuleOverride struct {
	SharedRPD   *int64                    `json:"shared_rpd,omitempty"`
	SharedRPM   *int64                    `json:"rpm,omitempty"`
	GeminiPro   *geminiModelQuotaOverride `json:"gemini_pro,omitempty"`
	GeminiFlash *geminiModelQuotaOverride `json:"gemini_flash,omitempty"`
	Desc        *string                   `json:"desc,omitempty"`
placeholder

type geminiModelQuotaOverride struct {
	RPD *int64 `json:"rpd,omitempty"`
	RPM *int64 `json:"rpm,omitempty"`
placeholder

type GeminiQuotaService struct {
	cfg         *config.Config
	settingRepo SettingRepository
	mu          sync.Mutex
	cachedAt    time.Time
	policy      *GeminiQuotaPolicy
placeholder

func NewGeminiQuotaService(cfg *config.Config, settingRepo SettingRepository) *GeminiQuotaService {
	return &GeminiQuotaService{
		cfg:         cfg,
		settingRepo: settingRepo,
placeholder
placeholder

func (s *GeminiQuotaService) Policy(ctx context.Context) *GeminiQuotaPolicy {
	if s == nil {
		return newGeminiQuotaPolicy()
placeholder

	now := time.Now()
	s.mu.Lock()
	if s.policy != nil && now.Sub(s.cachedAt) < geminiQuotaCacheTTL {
		policy := s.policy
		s.mu.Unlock()
		return policy
placeholder
	s.mu.Unlock()

	policy := newGeminiQuotaPolicy()
	if s.cfg != nil {
		policy.ApplyOverrides(s.cfg.Gemini.Quota.Tiers)
		if strings.TrimSpace(s.cfg.Gemini.Quota.Policy) != "" {
			raw := []byte(s.cfg.Gemini.Quota.Policy)
			var overridesV2 geminiQuotaOverridesV2
			if err := json.Unmarshal(raw, &overridesV2); err == nil && len(overridesV2.QuotaRules) > 0 {
				policy.ApplyQuotaRulesOverrides(overridesV2.QuotaRules)
		placeholder else {
				var overridesV1 geminiQuotaOverridesV1
				if err := json.Unmarshal(raw, &overridesV1); err != nil {
					log.Printf("gemini quota: parse config policy failed: %v", err)
			placeholder else {
					policy.ApplyOverrides(overridesV1.Tiers)
			placeholder
		placeholder
	placeholder
placeholder

	if s.settingRepo != nil {
		value, err := s.settingRepo.GetValue(ctx, SettingKeyGeminiQuotaPolicy)
		if err != nil && !errors.Is(err, ErrSettingNotFound) {
			log.Printf("gemini quota: load setting failed: %v", err)
	placeholder else if strings.TrimSpace(value) != "" {
			raw := []byte(value)
			var overridesV2 geminiQuotaOverridesV2
			if err := json.Unmarshal(raw, &overridesV2); err == nil && len(overridesV2.QuotaRules) > 0 {
				policy.ApplyQuotaRulesOverrides(overridesV2.QuotaRules)
		placeholder else {
				var overridesV1 geminiQuotaOverridesV1
				if err := json.Unmarshal(raw, &overridesV1); err != nil {
					log.Printf("gemini quota: parse setting failed: %v", err)
			placeholder else {
					policy.ApplyOverrides(overridesV1.Tiers)
			placeholder
		placeholder
	placeholder
placeholder

	s.mu.Lock()
	s.policy = policy
	s.cachedAt = now
	s.mu.Unlock()

	return policy
placeholder

func (s *GeminiQuotaService) QuotaForAccount(ctx context.Context, account *Account) (GeminiQuota, bool) {
	if account == nil || account.Platform != PlatformGemini {
		return GeminiQuota{placeholder, false
placeholder

	// Map (oauth_type + tier_id) to a canonical policy tier key.
	// This keeps the policy table stable even if upstream tier_id strings vary.
	tierKey := geminiQuotaTierKeyForAccount(account)
	if tierKey == "" {
		return GeminiQuota{placeholder, false
placeholder

	policy := s.Policy(ctx)
	return policy.QuotaForTier(tierKey)
placeholder

func (s *GeminiQuotaService) CooldownForTier(ctx context.Context, tierID string) time.Duration {
	policy := s.Policy(ctx)
	return policy.CooldownForTier(tierID)
placeholder

func (s *GeminiQuotaService) CooldownForAccount(ctx context.Context, account *Account) time.Duration {
	if s == nil || account == nil || account.Platform != PlatformGemini {
		return 5 * time.Minute
placeholder
	tierKey := geminiQuotaTierKeyForAccount(account)
	if strings.TrimSpace(tierKey) == "" {
		return 5 * time.Minute
placeholder
	return s.CooldownForTier(ctx, tierKey)
placeholder

func newGeminiQuotaPolicy() *GeminiQuotaPolicy {
	return &GeminiQuotaPolicy{
		tiers: map[string]GeminiTierPolicy{
			// --- AI Studio / API Key (per-model) ---
			// aistudio_free:
			//   - gemini_pro:   50 RPD / 2 RPM
			//   - gemini_flash: 1500 RPD / 15 RPM
			GeminiTierAIStudioFree: {Quota: GeminiQuota{ProRPD: 50, ProRPM: 2, FlashRPD: 1500, FlashRPM: 15placeholder, Cooldown: 30 * time.Minuteplaceholder,
			// aistudio_paid: -1 means "unlimited/pay-as-you-go" for RPD.
			GeminiTierAIStudioPaid: {Quota: GeminiQuota{ProRPD: -1, ProRPM: 1000, FlashRPD: -1, FlashRPM: 2000placeholder, Cooldown: 5 * time.Minuteplaceholder,

			// --- Google One (shared pool) ---
			GeminiTierGoogleOneFree: {Quota: GeminiQuota{SharedRPD: 1000, SharedRPM: 60placeholder, Cooldown: 30 * time.Minuteplaceholder,
			GeminiTierGoogleAIPro:   {Quota: GeminiQuota{SharedRPD: 1500, SharedRPM: 120placeholder, Cooldown: 5 * time.Minuteplaceholder,
			GeminiTierGoogleAIUltra: {Quota: GeminiQuota{SharedRPD: 2000, SharedRPM: 120placeholder, Cooldown: 5 * time.Minuteplaceholder,

			// --- GCP Code Assist (shared pool) ---
			GeminiTierGCPStandard:   {Quota: GeminiQuota{SharedRPD: 1500, SharedRPM: 120placeholder, Cooldown: 5 * time.Minuteplaceholder,
			GeminiTierGCPEnterprise: {Quota: GeminiQuota{SharedRPD: 2000, SharedRPM: 120placeholder, Cooldown: 5 * time.Minuteplaceholder,
	placeholder,
placeholder
placeholder

func (p *GeminiQuotaPolicy) ApplyOverrides(tiers map[string]config.GeminiTierQuotaConfig) {
	if p == nil || len(tiers) == 0 {
		return
placeholder
	for rawID, override := range tiers {
		tierID := normalizeGeminiTierID(rawID)
		if tierID == "" {
			continue
	placeholder
		policy, ok := p.tiers[tierID]
		if !ok {
			policy = GeminiTierPolicy{Cooldown: 5 * time.Minuteplaceholder
	placeholder
		// Backward-compatible overrides:
		// - If the tier uses shared quota, interpret pro_rpd as shared_rpd.
		// - Otherwise apply per-model overrides.
		if override.ProRPD != nil {
			if policy.Quota.SharedRPD > 0 {
				policy.Quota.SharedRPD = clampGeminiQuotaInt64WithUnlimited(*override.ProRPD)
		placeholder else {
				policy.Quota.ProRPD = clampGeminiQuotaInt64WithUnlimited(*override.ProRPD)
		placeholder
	placeholder
		if override.FlashRPD != nil {
			if policy.Quota.SharedRPD > 0 {
				// No separate flash RPD for shared tiers.
		placeholder else {
				policy.Quota.FlashRPD = clampGeminiQuotaInt64WithUnlimited(*override.FlashRPD)
		placeholder
	placeholder
		if override.CooldownMinutes != nil {
			minutes := clampGeminiQuotaInt(*override.CooldownMinutes)
			policy.Cooldown = time.Duration(minutes) * time.Minute
	placeholder
		p.tiers[tierID] = policy
placeholder
placeholder

func (p *GeminiQuotaPolicy) ApplyQuotaRulesOverrides(rules map[string]geminiQuotaRuleOverride) {
	if p == nil || len(rules) == 0 {
		return
placeholder
	for rawID, override := range rules {
		tierID := normalizeGeminiTierID(rawID)
		if tierID == "" {
			continue
	placeholder
		policy, ok := p.tiers[tierID]
		if !ok {
			policy = GeminiTierPolicy{Cooldown: 5 * time.Minuteplaceholder
	placeholder

		if override.SharedRPD != nil {
			policy.Quota.SharedRPD = clampGeminiQuotaInt64WithUnlimited(*override.SharedRPD)
	placeholder
		if override.SharedRPM != nil {
			policy.Quota.SharedRPM = clampGeminiQuotaRPM(*override.SharedRPM)
	placeholder
		if override.GeminiPro != nil {
			if override.GeminiPro.RPD != nil {
				policy.Quota.ProRPD = clampGeminiQuotaInt64WithUnlimited(*override.GeminiPro.RPD)
		placeholder
			if override.GeminiPro.RPM != nil {
				policy.Quota.ProRPM = clampGeminiQuotaRPM(*override.GeminiPro.RPM)
		placeholder
	placeholder
		if override.GeminiFlash != nil {
			if override.GeminiFlash.RPD != nil {
				policy.Quota.FlashRPD = clampGeminiQuotaInt64WithUnlimited(*override.GeminiFlash.RPD)
		placeholder
			if override.GeminiFlash.RPM != nil {
				policy.Quota.FlashRPM = clampGeminiQuotaRPM(*override.GeminiFlash.RPM)
		placeholder
	placeholder

		p.tiers[tierID] = policy
placeholder
placeholder

func (p *GeminiQuotaPolicy) QuotaForTier(tierID string) (GeminiQuota, bool) {
	policy, ok := p.policyForTier(tierID)
	if !ok {
		return GeminiQuota{placeholder, false
placeholder
	return policy.Quota, true
placeholder

func (p *GeminiQuotaPolicy) CooldownForTier(tierID string) time.Duration {
	policy, ok := p.policyForTier(tierID)
	if ok && policy.Cooldown > 0 {
		return policy.Cooldown
placeholder
	return 5 * time.Minute
placeholder

func (p *GeminiQuotaPolicy) policyForTier(tierID string) (GeminiTierPolicy, bool) {
	if p == nil {
		return GeminiTierPolicy{placeholder, false
placeholder
	normalized := normalizeGeminiTierID(tierID)
	if policy, ok := p.tiers[normalized]; ok {
		return policy, true
placeholder
	return GeminiTierPolicy{placeholder, false
placeholder

func normalizeGeminiTierID(tierID string) string {
	tierID = strings.TrimSpace(tierID)
	if tierID == "" {
		return ""
placeholder
	// Prefer canonical mapping (handles legacy tier strings).
	if canonical := canonicalGeminiTierID(tierID); canonical != "" {
		return canonical
placeholder
	// Accept older policy keys that used uppercase names.
	switch strings.ToUpper(tierID) {
	case "AISTUDIO_FREE":
		return GeminiTierAIStudioFree
	case "AISTUDIO_PAID":
		return GeminiTierAIStudioPaid
	case "GOOGLE_ONE_FREE":
		return GeminiTierGoogleOneFree
	case "GOOGLE_AI_PRO":
		return GeminiTierGoogleAIPro
	case "GOOGLE_AI_ULTRA":
		return GeminiTierGoogleAIUltra
	case "GCP_STANDARD":
		return GeminiTierGCPStandard
	case "GCP_ENTERPRISE":
		return GeminiTierGCPEnterprise
placeholder
	return strings.ToLower(tierID)
placeholder

func clampGeminiQuotaInt64WithUnlimited(value int64) int64 {
	if value < -1 {
		return 0
placeholder
	return value
placeholder

func clampGeminiQuotaInt(value int) int {
	if value < 0 {
		return 0
placeholder
	return value
placeholder

func clampGeminiQuotaRPM(value int64) int64 {
	if value < 0 {
		return 0
placeholder
	return value
placeholder

func geminiCooldownForTier(tierID string) time.Duration {
	policy := newGeminiQuotaPolicy()
	return policy.CooldownForTier(tierID)
placeholder

func geminiQuotaTierKeyForAccount(account *Account) string {
	if account == nil || account.Platform != PlatformGemini {
		return ""
placeholder

	// Note: GeminiOAuthType() already defaults legacy (project_id present) to code_assist.
	oauthType := strings.ToLower(strings.TrimSpace(account.GeminiOAuthType()))
	rawTier := strings.TrimSpace(account.GeminiTierID())

	// Prefer the canonical tier stored in credentials.
	if tierID := canonicalGeminiTierIDForOAuthType(oauthType, rawTier); tierID != "" && tierID != GeminiTierGoogleOneUnknown {
		return tierID
placeholder

	// Fallback defaults when tier_id is missing or unknown.
	switch oauthType {
	case "google_one":
		return GeminiTierGoogleOneFree
	case "code_assist":
		return GeminiTierGCPStandard
	case "ai_studio":
		return GeminiTierAIStudioFree
	default:
		// API Key accounts (type=apikey) have empty oauth_type and are treated as AI Studio.
		return GeminiTierAIStudioFree
placeholder
placeholder

func geminiModelClassFromName(model string) geminiModelClass {
	name := strings.ToLower(strings.TrimSpace(model))
	if strings.Contains(name, "flash") || strings.Contains(name, "lite") {
		return geminiModelFlash
placeholder
	return geminiModelPro
placeholder

func geminiAggregateUsage(stats []usagestats.ModelStat) GeminiUsageTotals {
	var totals GeminiUsageTotals
	for _, stat := range stats {
		switch geminiModelClassFromName(stat.Model) {
		case geminiModelFlash:
			totals.FlashRequests += stat.Requests
			totals.FlashTokens += stat.TotalTokens
			totals.FlashCost += stat.ActualCost
		default:
			totals.ProRequests += stat.Requests
			totals.ProTokens += stat.TotalTokens
			totals.ProCost += stat.ActualCost
	placeholder
placeholder
	return totals
placeholder

func geminiQuotaLocation() *time.Location {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return time.FixedZone("PST", -8*3600)
placeholder
	return loc
placeholder

func geminiDailyWindowStart(now time.Time) time.Time {
	loc := geminiQuotaLocation()
	localNow := now.In(loc)
	return time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
placeholder

func geminiDailyResetTime(now time.Time) time.Time {
	loc := geminiQuotaLocation()
	localNow := now.In(loc)
	start := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	reset := start.Add(24 * time.Hour)
	if !reset.After(localNow) {
		reset = reset.Add(24 * time.Hour)
placeholder
	return reset
placeholder
