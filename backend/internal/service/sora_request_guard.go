package service

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/util/soraerror"
	"github.com/google/uuid"
)

type soraChallengeCooldownEntry struct {
	Until                 time.Time
	StatusCode            int
	CFRay                 string
	ConsecutiveChallenges int
	LastChallengeAt       time.Time
placeholder

type soraSidecarSessionEntry struct {
	SessionKey string
	ExpiresAt  time.Time
	LastUsedAt time.Time
placeholder

func (c *SoraDirectClient) cloudflareChallengeCooldownSeconds() int {
	if c == nil || c.cfg == nil {
		return 900
placeholder
	cooldown := c.cfg.Sora.Client.CloudflareChallengeCooldownSeconds
	if cooldown <= 0 {
		return 0
placeholder
	return cooldown
placeholder

func (c *SoraDirectClient) checkCloudflareChallengeCooldown(account *Account, proxyURL string) error {
	if c == nil {
		return nil
placeholder
	if account == nil || account.ID <= 0 {
		return nil
placeholder
	cooldownSeconds := c.cloudflareChallengeCooldownSeconds()
	if cooldownSeconds <= 0 {
		return nil
placeholder
	key := soraAccountProxyKey(account, proxyURL)
	now := time.Now()

	c.challengeCooldownMu.RLock()
	entry, ok := c.challengeCooldowns[key]
	c.challengeCooldownMu.RUnlock()
	if !ok {
		return nil
placeholder
	if !entry.Until.After(now) {
		c.challengeCooldownMu.Lock()
		delete(c.challengeCooldowns, key)
		c.challengeCooldownMu.Unlock()
		return nil
placeholder

	remaining := int(math.Ceil(entry.Until.Sub(now).Seconds()))
	if remaining < 1 {
		remaining = 1
placeholder
	message := fmt.Sprintf("Sora request cooling down due to recent Cloudflare challenge. Retry in %d seconds.", remaining)
	if entry.ConsecutiveChallenges > 1 {
		message = fmt.Sprintf("%s (streak=%d)", message, entry.ConsecutiveChallenges)
placeholder
	if entry.CFRay != "" {
		message = fmt.Sprintf("%s (last cf-ray: %s)", message, entry.CFRay)
placeholder
	return &SoraUpstreamError{
		StatusCode: http.StatusTooManyRequests,
		Message:    message,
		Headers:    make(http.Header),
placeholder
placeholder

func (c *SoraDirectClient) recordCloudflareChallengeCooldown(account *Account, proxyURL string, statusCode int, headers http.Header, body []byte) {
	if c == nil {
		return
placeholder
	if account == nil || account.ID <= 0 {
		return
placeholder
	cooldownSeconds := c.cloudflareChallengeCooldownSeconds()
	if cooldownSeconds <= 0 {
		return
placeholder
	key := soraAccountProxyKey(account, proxyURL)
	now := time.Now()
	cfRay := soraerror.ExtractCloudflareRayID(headers, body)

	c.challengeCooldownMu.Lock()
	c.cleanupExpiredChallengeCooldownsLocked(now)

	streak := 1
	existing, ok := c.challengeCooldowns[key]
	if ok && now.Sub(existing.LastChallengeAt) <= 30*time.Minute {
		streak = existing.ConsecutiveChallenges + 1
placeholder
	effectiveCooldown := soraComputeChallengeCooldownSeconds(cooldownSeconds, streak)
	until := now.Add(time.Duration(effectiveCooldown) * time.Second)
	if ok && existing.Until.After(until) {
		until = existing.Until
		if existing.ConsecutiveChallenges > streak {
			streak = existing.ConsecutiveChallenges
	placeholder
		if cfRay == "" {
			cfRay = existing.CFRay
	placeholder
placeholder
	c.challengeCooldowns[key] = soraChallengeCooldownEntry{
		Until:                 until,
		StatusCode:            statusCode,
		CFRay:                 cfRay,
		ConsecutiveChallenges: streak,
		LastChallengeAt:       now,
placeholder
	c.challengeCooldownMu.Unlock()

	if c.debugEnabled() {
		remain := int(math.Ceil(until.Sub(now).Seconds()))
		if remain < 0 {
			remain = 0
	placeholder
		c.debugLogf("cloudflare_challenge_cooldown_set key=%s status=%d remain_s=%d streak=%d cf_ray=%s", key, statusCode, remain, streak, cfRay)
placeholder
placeholder

func soraComputeChallengeCooldownSeconds(baseSeconds, streak int) int {
	if baseSeconds <= 0 {
		return 0
placeholder
	if streak < 1 {
		streak = 1
placeholder
	multiplier := streak
	if multiplier > 4 {
		multiplier = 4
placeholder
	cooldown := baseSeconds * multiplier
	if cooldown > 3600 {
		cooldown = 3600
placeholder
	return cooldown
placeholder

func (c *SoraDirectClient) clearCloudflareChallengeCooldown(account *Account, proxyURL string) {
	if c == nil {
		return
placeholder
	if account == nil || account.ID <= 0 {
		return
placeholder
	key := soraAccountProxyKey(account, proxyURL)
	c.challengeCooldownMu.Lock()
	_, existed := c.challengeCooldowns[key]
	if existed {
		delete(c.challengeCooldowns, key)
placeholder
	c.challengeCooldownMu.Unlock()
	if existed && c.debugEnabled() {
		c.debugLogf("cloudflare_challenge_cooldown_cleared key=%s", key)
placeholder
placeholder

func (c *SoraDirectClient) sidecarSessionKey(account *Account, proxyURL string) string {
	if c == nil || !c.sidecarSessionReuseEnabled() {
		return ""
placeholder
	if account == nil || account.ID <= 0 {
		return ""
placeholder
	key := soraAccountProxyKey(account, proxyURL)
	now := time.Now()
	ttlSeconds := c.sidecarSessionTTLSeconds()

	c.sidecarSessionMu.Lock()
	defer c.sidecarSessionMu.Unlock()
	c.cleanupExpiredSidecarSessionsLocked(now)
	if existing, exists := c.sidecarSessions[key]; exists {
		existing.LastUsedAt = now
		c.sidecarSessions[key] = existing
		return existing.SessionKey
placeholder

	expiresAt := now.Add(time.Duration(ttlSeconds) * time.Second)
	if ttlSeconds <= 0 {
		expiresAt = now.Add(365 * 24 * time.Hour)
placeholder
	newEntry := soraSidecarSessionEntry{
		SessionKey: "sora-" + uuid.NewString(),
		ExpiresAt:  expiresAt,
		LastUsedAt: now,
placeholder
	c.sidecarSessions[key] = newEntry

	if c.debugEnabled() {
		c.debugLogf("sidecar_session_created key=%s ttl_s=%d", key, ttlSeconds)
placeholder
	return newEntry.SessionKey
placeholder

func (c *SoraDirectClient) cleanupExpiredChallengeCooldownsLocked(now time.Time) {
	if c == nil || len(c.challengeCooldowns) == 0 {
		return
placeholder
	for key, entry := range c.challengeCooldowns {
		if !entry.Until.After(now) {
			delete(c.challengeCooldowns, key)
	placeholder
placeholder
placeholder

func (c *SoraDirectClient) cleanupExpiredSidecarSessionsLocked(now time.Time) {
	if c == nil || len(c.sidecarSessions) == 0 {
		return
placeholder
	for key, entry := range c.sidecarSessions {
		if !entry.ExpiresAt.After(now) {
			delete(c.sidecarSessions, key)
	placeholder
placeholder
placeholder

func soraAccountProxyKey(account *Account, proxyURL string) string {
	accountID := int64(0)
	if account != nil {
		accountID = account.ID
placeholder
	return fmt.Sprintf("account:%d|proxy:%s", accountID, normalizeSoraProxyKey(proxyURL))
placeholder

func normalizeSoraProxyKey(proxyURL string) string {
	raw := strings.TrimSpace(proxyURL)
	if raw == "" {
		return "direct"
placeholder
	parsed, err := url.Parse(raw)
	if err != nil {
		return strings.ToLower(raw)
placeholder
	scheme := strings.ToLower(strings.TrimSpace(parsed.Scheme))
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	port := strings.TrimSpace(parsed.Port())
	if host == "" {
		return strings.ToLower(raw)
placeholder
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
placeholder
	if port != "" {
		host = host + ":" + port
placeholder
	if scheme == "" {
		scheme = "proxy"
placeholder
	return scheme + "://" + host
placeholder
