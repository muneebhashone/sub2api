package service

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

// In-memory team+model rate-limit overlay for Grok OAuth. When xAI rate-limits
// one account in a team for a model, sibling accounts sharing team_id skip the
// same model until the cooldown expires (mirrors grok2api teamModelRateLimit).
//
// Process-local only: multi-instance deployments each learn the block from their
// own 429s. Prefer short TTLs so drift self-heals.
type grokTeamModelRateLimit struct {
	Until time.Time
placeholder

type grokTeamModelRateLimitStore struct {
	mu    sync.Mutex
	items map[string]grokTeamModelRateLimit
placeholder

var globalGrokTeamModelRateLimits = &grokTeamModelRateLimitStore{
	items: make(map[string]grokTeamModelRateLimit),
placeholder

const (
	grokTeamRateLimitDefaultTTL = 10 * time.Minute
	grokTeamRateLimitMaxTTL     = time.Hour
	grokTeamRateLimitMinTTL     = 30 * time.Second
)

func grokTeamFingerprint(teamID string) string {
	teamID = strings.TrimSpace(teamID)
	if teamID == "" {
		return ""
placeholder
	sum := sha256.Sum256([]byte(strings.ToLower(teamID)))
	return hex.EncodeToString(sum[:8])
placeholder

func grokTeamModelRateLimitKey(teamFingerprint, model string) string {
	return teamFingerprint + "|" + strings.ToLower(strings.TrimSpace(model))
placeholder

func accountGrokTeamID(account *Account) string {
	if account == nil {
		return ""
placeholder
	return strings.TrimSpace(account.GetCredential("team_id"))
placeholder

// markGrokTeamModelRateLimit records that this team+model pair should be skipped
// until until. No-op when team_id or model is empty.
func markGrokTeamModelRateLimit(account *Account, model string, until time.Time) {
	if account == nil || !account.IsGrokOAuth() {
		return
placeholder
	fp := grokTeamFingerprint(accountGrokTeamID(account))
	model = strings.TrimSpace(model)
	if fp == "" || model == "" || until.IsZero() {
		return
placeholder
	now := time.Now()
	if !until.After(now) {
		until = now.Add(grokTeamRateLimitDefaultTTL)
placeholder
	maxUntil := now.Add(grokTeamRateLimitMaxTTL)
	if until.After(maxUntil) {
		until = maxUntil
placeholder
	key := grokTeamModelRateLimitKey(fp, model)
	globalGrokTeamModelRateLimits.mu.Lock()
	defer globalGrokTeamModelRateLimits.mu.Unlock()
	if cur, ok := globalGrokTeamModelRateLimits.items[key]; ok && cur.Until.After(until) {
		return
placeholder
	globalGrokTeamModelRateLimits.items[key] = grokTeamModelRateLimit{Until: untilplaceholder
	// Opportunistic prune of expired entries.
	for k, v := range globalGrokTeamModelRateLimits.items {
		if !v.Until.After(now) {
			delete(globalGrokTeamModelRateLimits.items, k)
	placeholder
placeholder
placeholder

// isGrokTeamModelRateLimited reports whether the account's team is currently
// blocked for the requested model.
func isGrokTeamModelRateLimited(account *Account, model string, now time.Time) bool {
	if account == nil || !account.IsGrokOAuth() {
		return false
placeholder
	fp := grokTeamFingerprint(accountGrokTeamID(account))
	model = strings.TrimSpace(model)
	if fp == "" || model == "" {
		return false
placeholder
	key := grokTeamModelRateLimitKey(fp, model)
	globalGrokTeamModelRateLimits.mu.Lock()
	defer globalGrokTeamModelRateLimits.mu.Unlock()
	cur, ok := globalGrokTeamModelRateLimits.items[key]
	if !ok {
		return false
placeholder
	if !cur.Until.After(now) {
		delete(globalGrokTeamModelRateLimits.items, key)
		return false
placeholder
	return true
placeholder

// filterGrokTeamModelRateLimitedAccounts drops candidates whose team is under a
// model-scoped rate-limit cool. Accounts without team_id pass through.
func filterGrokTeamModelRateLimitedAccounts(accounts []Account, model string, now time.Time) []Account {
	if len(accounts) == 0 || strings.TrimSpace(model) == "" {
		return accounts
placeholder
	out := accounts[:0]
	kept := false
	for i := range accounts {
		if isGrokTeamModelRateLimited(&accounts[i], model, now) {
			continue
	placeholder
		out = append(out, accounts[i])
		kept = true
placeholder
	if !kept && len(out) == 0 {
		// All filtered — return empty (caller treats as no capacity).
		return nil
placeholder
	return out
placeholder

// resolveGrokTeamRateLimitUntil derives a team cool window from an observed
// account rate-limit reset, with sane clamps.
func resolveGrokTeamRateLimitUntil(resetAt, now time.Time) time.Time {
	if resetAt.After(now.Add(grokTeamRateLimitMinTTL)) {
		maxUntil := now.Add(grokTeamRateLimitMaxTTL)
		if resetAt.After(maxUntil) {
			return maxUntil
	placeholder
		return resetAt
placeholder
	return now.Add(grokTeamRateLimitDefaultTTL)
placeholder
