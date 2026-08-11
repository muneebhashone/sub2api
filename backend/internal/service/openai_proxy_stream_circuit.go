package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultOpenAIProxyStreamFailureThreshold  = 2
	defaultOpenAIProxyStreamFailureWindow     = time.Minute
	defaultOpenAIProxyStreamQuarantineTTL     = 10 * time.Minute
	defaultOpenAIProxyStreamCircuitMaxEntries = 4096
	// defaultOpenAIProxyStreamFailureCollapse merges disconnects that land
	// within this interval into a single failure event. One proxy or HTTP/2
	// connection loss kills every stream multiplexed on it at the same
	// moment; counting each stream separately would trip the breaker from a
	// single upstream event (#5056).
	defaultOpenAIProxyStreamFailureCollapse = 3 * time.Second
	// openAIProxyStreamFailOpenLogInterval rate-limits the fail-open warning
	// so an outage-mode burst does not flood the log.
	openAIProxyStreamFailOpenLogInterval = 5 * time.Second
)

type openAIProxyStreamCircuitSettings struct {
	disabled         bool
	failureThreshold int
	failureWindow    time.Duration
	quarantineTTL    time.Duration
	collapseInterval time.Duration
	maxEntries       int
placeholder

type openAIProxyStreamCircuitEntry struct {
	failureCount  int
	windowStart   time.Time
	lastFailureAt time.Time
	blockedUntil  time.Time
	lastTouched   time.Time
placeholder

// openAIProxyStreamCircuit is an in-process, proxy-ID keyed circuit. It is
// intentionally bounded and ephemeral: a restart clears observations, while a
// tripped entry expires automatically after its TTL.
type openAIProxyStreamCircuit struct {
	mu       sync.Mutex
	settings openAIProxyStreamCircuitSettings
	entries  map[int64]openAIProxyStreamCircuitEntry
placeholder

func resolveOpenAIProxyStreamCircuitSettings(s *OpenAIGatewayService) openAIProxyStreamCircuitSettings {
	settings := openAIProxyStreamCircuitSettings{
		failureThreshold: defaultOpenAIProxyStreamFailureThreshold,
		failureWindow:    defaultOpenAIProxyStreamFailureWindow,
		quarantineTTL:    defaultOpenAIProxyStreamQuarantineTTL,
		collapseInterval: defaultOpenAIProxyStreamFailureCollapse,
		maxEntries:       defaultOpenAIProxyStreamCircuitMaxEntries,
placeholder
	if s == nil || s.cfg == nil {
		return settings
placeholder
	cfg := s.cfg.Gateway.OpenAIProxyStreamCircuit
	settings.disabled = cfg.Disabled
	if cfg.FailureThreshold > 0 {
		settings.failureThreshold = cfg.FailureThreshold
placeholder
	if cfg.WindowSeconds > 0 {
		settings.failureWindow = time.Duration(cfg.WindowSeconds) * time.Second
placeholder
	if cfg.TTLSeconds > 0 {
		settings.quarantineTTL = time.Duration(cfg.TTLSeconds) * time.Second
placeholder
	return settings
placeholder

func newOpenAIProxyStreamCircuit(settings openAIProxyStreamCircuitSettings) *openAIProxyStreamCircuit {
	if settings.failureThreshold <= 0 {
		settings.failureThreshold = defaultOpenAIProxyStreamFailureThreshold
placeholder
	if settings.failureWindow <= 0 {
		settings.failureWindow = defaultOpenAIProxyStreamFailureWindow
placeholder
	if settings.quarantineTTL <= 0 {
		settings.quarantineTTL = defaultOpenAIProxyStreamQuarantineTTL
placeholder
	if settings.maxEntries <= 0 {
		settings.maxEntries = defaultOpenAIProxyStreamCircuitMaxEntries
placeholder
	if settings.collapseInterval < 0 {
		settings.collapseInterval = 0
placeholder
	return &openAIProxyStreamCircuit{
		settings: settings,
		entries:  make(map[int64]openAIProxyStreamCircuitEntry),
placeholder
placeholder

func (s *OpenAIGatewayService) getOpenAIProxyStreamCircuit() *openAIProxyStreamCircuit {
	if s == nil {
		return nil
placeholder
	s.openaiProxyStreamCircuitOnce.Do(func() {
		if s.openaiProxyStreamCircuit == nil {
			s.openaiProxyStreamCircuit = newOpenAIProxyStreamCircuit(resolveOpenAIProxyStreamCircuitSettings(s))
	placeholder
placeholder)
	return s.openaiProxyStreamCircuit
placeholder

func (c *openAIProxyStreamCircuit) recordFailure(proxyID int64, now time.Time) (bool, time.Time) {
	if c == nil || c.settings.disabled || proxyID <= 0 {
		return false, time.Time{placeholder
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[proxyID]
	if exists && now.Before(entry.blockedUntil) {
		entry.lastTouched = now
		c.entries[proxyID] = entry
		return false, entry.blockedUntil
placeholder
	if !exists {
		c.ensureCapacityLocked(now)
placeholder
	if entry.windowStart.IsZero() || now.Before(entry.windowStart) || now.Sub(entry.windowStart) > c.settings.failureWindow {
		entry.failureCount = 0
		entry.windowStart = now
		entry.blockedUntil = time.Time{placeholder
placeholder
	// Collapse a burst of disconnects into one failure event: when a proxy or
	// a multiplexed HTTP/2 connection dies, every in-flight stream reports the
	// same underlying incident within moments of each other.
	if c.settings.collapseInterval > 0 && !entry.lastFailureAt.IsZero() &&
		now.Sub(entry.lastFailureAt) >= 0 && now.Sub(entry.lastFailureAt) < c.settings.collapseInterval {
		entry.lastTouched = now
		c.entries[proxyID] = entry
		return false, time.Time{placeholder
placeholder
	entry.failureCount++
	entry.lastFailureAt = now
	entry.lastTouched = now
	tripped := entry.failureCount >= c.settings.failureThreshold
	if tripped {
		entry.blockedUntil = now.Add(c.settings.quarantineTTL)
placeholder
	c.entries[proxyID] = entry
	return tripped, entry.blockedUntil
placeholder

func (c *openAIProxyStreamCircuit) recordSuccess(proxyID int64) bool {
	if c == nil || proxyID <= 0 {
		return false
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.entries[proxyID]; !ok {
		return false
placeholder
	delete(c.entries, proxyID)
	return true
placeholder

func (c *openAIProxyStreamCircuit) isBlocked(proxyID int64, now time.Time) bool {
	if c == nil || c.settings.disabled || proxyID <= 0 {
		return false
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[proxyID]
	if !ok || entry.blockedUntil.IsZero() {
		return false
placeholder
	if !now.Before(entry.blockedUntil) {
		delete(c.entries, proxyID)
		return false
placeholder
	return true
placeholder

// activeBlockCount reports how many proxies are currently quarantined. It
// gates the fail-open retry: a "no available accounts" selection result only
// warrants a second, quarantine-blind pass when the circuit is actually
// withholding capacity.
func (c *openAIProxyStreamCircuit) activeBlockCount(now time.Time) int {
	if c == nil || c.settings.disabled {
		return 0
placeholder
	c.mu.Lock()
	defer c.mu.Unlock()
	count := 0
	for _, entry := range c.entries {
		if !entry.blockedUntil.IsZero() && now.Before(entry.blockedUntil) {
			count++
	placeholder
placeholder
	return count
placeholder

func (c *openAIProxyStreamCircuit) ensureCapacityLocked(now time.Time) {
	if len(c.entries) < c.settings.maxEntries {
		return
placeholder
	for proxyID, entry := range c.entries {
		staleObservation := entry.blockedUntil.IsZero() && now.Sub(entry.lastTouched) > c.settings.failureWindow
		expiredQuarantine := !entry.blockedUntil.IsZero() && !now.Before(entry.blockedUntil)
		if staleObservation || expiredQuarantine {
			delete(c.entries, proxyID)
	placeholder
placeholder
	if len(c.entries) < c.settings.maxEntries {
		return
placeholder
	var oldestProxyID int64
	var oldest time.Time
	for proxyID, entry := range c.entries {
		if oldestProxyID == 0 || entry.lastTouched.Before(oldest) {
			oldestProxyID = proxyID
			oldest = entry.lastTouched
	placeholder
placeholder
	if oldestProxyID > 0 {
		delete(c.entries, oldestProxyID)
placeholder
placeholder

func openAIProxyStreamCircuitProxyID(account *Account) (int64, bool) {
	if account == nil || account.Platform != PlatformOpenAI || account.ProxyID == nil || *account.ProxyID <= 0 {
		return 0, false
placeholder
	return *account.ProxyID, true
placeholder

func (s *OpenAIGatewayService) recordOpenAIProxyStreamDisconnect(account *Account, streamErr error, upstreamRequestID string) {
	proxyID, ok := openAIProxyStreamCircuitProxyID(account)
	if !ok || streamErr == nil || errors.Is(streamErr, context.Canceled) || errors.Is(streamErr, context.DeadlineExceeded) {
		return
placeholder
	circuit := s.getOpenAIProxyStreamCircuit()
	tripped, until := circuit.recordFailure(proxyID, time.Now())
	if !tripped {
		return
placeholder
	logger.L().With(zap.String("component", "service.openai_gateway")).Warn(
		"openai.proxy_quarantined_stream_disconnect",
		zap.Int64("proxy_id", proxyID),
		zap.Int64("account_id", account.ID),
		zap.Time("until", until),
		zap.String("upstream_request_id", upstreamRequestID),
		zap.String("error", sanitizeUpstreamErrorMessage(streamErr.Error())),
	)
placeholder

func (s *OpenAIGatewayService) clearOpenAIProxyStreamDisconnect(account *Account) {
	proxyID, ok := openAIProxyStreamCircuitProxyID(account)
	if !ok {
		return
placeholder
	if circuit := s.getOpenAIProxyStreamCircuit(); circuit != nil {
		circuit.recordSuccess(proxyID)
placeholder
placeholder

// openAIProxyStreamQuarantineBypassKey marks a selection pass that must ignore
// proxy quarantine. It is set for the second, fail-open selection attempt when
// the first pass found no available account while the circuit was withholding
// proxies: a degraded proxy is strictly better than answering 502 (#5056).
type openAIProxyStreamQuarantineBypassKey struct{placeholder

func withOpenAIProxyStreamQuarantineBypass(ctx context.Context) context.Context {
	return context.WithValue(ctx, openAIProxyStreamQuarantineBypassKey{placeholder, true)
placeholder

func openAIProxyStreamQuarantineBypassed(ctx context.Context) bool {
	if ctx == nil {
		return false
placeholder
	bypassed, _ := ctx.Value(openAIProxyStreamQuarantineBypassKey{placeholder).(bool)
	return bypassed
placeholder

func (s *OpenAIGatewayService) isOpenAIProxyStreamQuarantined(ctx context.Context, account *Account) bool {
	proxyID, ok := openAIProxyStreamCircuitProxyID(account)
	if !ok {
		return false
placeholder
	if openAIProxyStreamQuarantineBypassed(ctx) {
		return false
placeholder
	circuit := s.getOpenAIProxyStreamCircuit()
	return circuit != nil && circuit.isBlocked(proxyID, time.Now())
placeholder

// logOpenAIProxyStreamQuarantineFailOpen emits a rate-limited warning when a
// selection pass had to re-admit quarantined proxies to serve at all.
func (s *OpenAIGatewayService) logOpenAIProxyStreamQuarantineFailOpen(requestedModel string, blockedProxies int) {
	now := time.Now().UnixNano()
	last := s.openaiProxyStreamFailOpenLogAt.Load()
	if now-last < int64(openAIProxyStreamFailOpenLogInterval) ||
		!s.openaiProxyStreamFailOpenLogAt.CompareAndSwap(last, now) {
		return
placeholder
	logger.L().With(zap.String("component", "service.openai_gateway")).Warn(
		"openai.proxy_stream_quarantine_fail_open",
		zap.Int("blocked_proxies", blockedProxies),
		zap.String("model", requestedModel),
	)
placeholder
