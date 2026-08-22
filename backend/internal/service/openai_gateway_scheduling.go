package service

// 本文件由 openai_gateway_service.go 纯移动拆分而来：粘性会话哈希、账号选择与
// 负载感知调度、配额自动暂停判定、并发槽位获取。仅做代码搬迁，无任何行为变更。

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	openCodeSessionAffinityHeader = "X-Session-Affinity"
	openCodeSessionIDHeader       = "X-Session-Id"
	openCodeNativeSessionHeader   = "X-OpenCode-Session"
	codeBuddyConversationHeader   = "X-Conversation-ID"
)

var explicitOpenAIHeaderSessionNames = []string{
	"session_id",
	"conversation_id",
	openCodeSessionAffinityHeader,
	openCodeSessionIDHeader,
	openCodeNativeSessionHeader,
	codeBuddyConversationHeader,
placeholder

// explicitOpenAIHeaderSessionID resolves stable conversation identifiers sent
// by OpenAI-compatible clients. Keep this list limited to session-scoped
// fields: request/message IDs rotate every turn and would defeat sticky routing
// and upstream prompt caching.
func explicitOpenAIHeaderSessionID(c *gin.Context) string {
	if c == nil {
		return ""
placeholder

	for _, header := range explicitOpenAIHeaderSessionNames {
		if sessionID := strings.TrimSpace(c.GetHeader(header)); sessionID != "" {
			return sessionID
	placeholder
placeholder
	return ""
placeholder

// ExtractSessionID extracts the raw session ID from headers or body without hashing.
// Used by ForwardAsAnthropic to pass as prompt_cache_key for upstream cache.
func (s *OpenAIGatewayService) ExtractSessionID(c *gin.Context, body []byte) string {
	return explicitOpenAIRequestSessionID(c, body)
placeholder

func explicitOpenAISessionID(c *gin.Context, body []byte) string {
	if c == nil {
		return ""
placeholder

	sessionID := explicitOpenAIHeaderSessionID(c)
	if sessionID == "" && len(body) > 0 {
		sessionID = strings.TrimSpace(openAIRequestPayloadView(body).Get("prompt_cache_key").String())
placeholder
	return sessionID
placeholder

// openAIRequestPayloadView unwraps Responses WebSocket event envelopes while
// leaving ordinary HTTP objects untouched even when they contain a response
// field for another purpose.
func openAIRequestPayloadView(body []byte) gjson.Result {
	root := parseRawJSONView(body)
	eventType := strings.ToLower(strings.TrimSpace(root.Get("type").String()))
	if strings.HasPrefix(eventType, "response.") {
		if response := root.Get("response"); response.Exists() && response.IsObject() {
			return response
	placeholder
placeholder
	return root
placeholder

// explicitOpenAIRequestSessionID extends the common OpenAI session signals
// with Grok's native conversation header only for requests authenticated to a
// Grok group. This keeps an unrelated x-grok-conv-id header from changing
// scheduling or upstream session behavior for non-Grok groups.
//
// For Grok groups only, previous_response_id is a last-resort sticky seed so
// multi-turn Responses chains stay on the same OAuth account when no explicit
// session/conversation/prompt_cache_key is present. Non-Grok groups omit this
// so HTTP OpenAI paths that delete previous_response_id before upstream are
// unchanged.
func explicitOpenAIRequestSessionID(c *gin.Context, body []byte) string {
	if c == nil {
		return ""
placeholder

	sessionID := explicitOpenAIHeaderSessionID(c)
	if sessionID == "" && isGrokRequestContext(c) {
		sessionID = strings.TrimSpace(c.GetHeader(grokConversationIDHeader))
placeholder
	if sessionID == "" && len(body) > 0 {
		sessionID = strings.TrimSpace(openAIRequestPayloadView(body).Get("prompt_cache_key").String())
placeholder
	if sessionID == "" && isGrokRequestContext(c) && len(body) > 0 {
		sessionID = grokPreviousResponseSessionSeed(body)
placeholder
	return sessionID
placeholder

// grokPreviousResponseSessionSeed returns a stable sticky seed from a Responses
// previous_response_id. Only resp_* response ids are accepted; message ids and
// unknown shapes must not pin sticky routing or prompt-cache identity.
func grokPreviousResponseSessionSeed(body []byte) string {
	id := strings.TrimSpace(gjson.GetBytes(body, "previous_response_id").String())
	if id == "" {
		return ""
placeholder
	if ClassifyOpenAIPreviousResponseIDKind(id) != OpenAIPreviousResponseIDKindResponseID {
		return ""
placeholder
	// Namespace so content-derived seeds never collide with response ids.
	return "grok-prev-resp:" + id
placeholder

// GenerateExplicitSessionHash generates a sticky-session hash only from explicit
// client session signals. It intentionally skips content-derived fallback and is
// used by stateless endpoints such as /v1/images.
func (s *OpenAIGatewayService) GenerateExplicitSessionHash(c *gin.Context, body []byte) string {
	sessionID := explicitOpenAIRequestSessionID(c, body)
	if sessionID == "" {
		return ""
placeholder

	currentHash, legacyHash := deriveOpenAISessionHashes(sessionID)
	attachOpenAILegacySessionHashToGin(c, legacyHash)
	return currentHash
placeholder

// GenerateSessionHash generates a sticky-session hash for OpenAI requests.
//
// Priority:
//  1. Header: session_id
//  2. Header: conversation_id
//  3. Header: x-session-affinity / x-session-id / x-opencode-session (OpenCode)
//  4. Header: x-conversation-id (CodeBuddy)
//  5. Header: x-grok-conv-id (Grok groups only)
//  6. Body:   prompt_cache_key
//  7. Body:   content-based fallback (model + system + tools + first user message)
//
// Grok sticky affinity is intentionally separate from the upstream
// prompt_cache_key identity (resolveGrokCacheIdentity): sticky pins an OAuth
// account for multi-turn routing, while the cache identity is tenant+model
// isolated for xAI server-side prompt cache. For Grok groups we scope the
// sticky seed with the client-requested model so switching models does not
// inherit a stale account binding (grok2api affinityKey pattern).
func (s *OpenAIGatewayService) GenerateSessionHash(c *gin.Context, body []byte) string {
	if c == nil {
		return ""
placeholder

	sessionID := explicitOpenAIRequestSessionID(c, body)
	if sessionID == "" && len(body) > 0 {
		sessionID = deriveOpenAIContentSessionSeed(body)
placeholder
	if sessionID == "" {
		return ""
placeholder

	if isGrokRequestContext(c) {
		sessionID = grokStickyAffinitySeed(sessionID, body)
placeholder

	currentHash, legacyHash := deriveOpenAISessionHashes(sessionID)
	attachOpenAILegacySessionHashToGin(c, legacyHash)
	return currentHash
placeholder

// grokStickyAffinitySeed scopes sticky routing by model without changing the
// upstream prompt_cache_key written by applyGrokResponsesCacheIdentity.
func grokStickyAffinitySeed(sessionID string, body []byte) string {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return ""
placeholder
	model := ""
	if len(body) > 0 {
		model = strings.ToLower(strings.TrimSpace(gjson.GetBytes(body, "model").String()))
placeholder
	if model == "" {
		return "grok-affinity:v1:" + sessionID
placeholder
	return "grok-affinity:v1:" + model + ":" + sessionID
placeholder

// GenerateSessionHashWithFallback 先按常规信号生成会话哈希；
// 当未携带 session_id/conversation_id/prompt_cache_key 时，使用 fallbackSeed 生成稳定哈希。
// 该方法用于 WS ingress，避免会话信号缺失时发生跨账号漂移。
func (s *OpenAIGatewayService) GenerateSessionHashWithFallback(c *gin.Context, body []byte, fallbackSeed string) string {
	sessionHash := s.GenerateSessionHash(c, body)
	if sessionHash != "" {
		return sessionHash
placeholder

	seed := strings.TrimSpace(fallbackSeed)
	if seed == "" {
		return ""
placeholder

	currentHash, legacyHash := deriveOpenAISessionHashes(seed)
	attachOpenAILegacySessionHashToGin(c, legacyHash)
	return currentHash
placeholder

func resolveOpenAIUpstreamOriginator(c *gin.Context, isOfficialClient bool) string {
	if c != nil {
		if originator := strings.TrimSpace(c.GetHeader("originator")); originator != "" {
			return originator
	placeholder
placeholder
	if isOfficialClient {
		return openai.CodexDefaultOriginator
placeholder
	return "opencode"
placeholder

// BindStickySession sets session -> account binding with standard TTL.
func (s *OpenAIGatewayService) BindStickySession(ctx context.Context, groupID *int64, sessionHash string, accountID int64) error {
	if sessionHash == "" || accountID <= 0 {
		return nil
placeholder
	ttl := openaiStickySessionTTL
	if s != nil && s.cfg != nil && s.cfg.Gateway.OpenAIWS.StickySessionTTLSeconds > 0 {
		ttl = time.Duration(s.cfg.Gateway.OpenAIWS.StickySessionTTLSeconds) * time.Second
placeholder
	return s.setStickySessionAccountID(ctx, groupID, sessionHash, accountID, ttl)
placeholder

// SelectAccount selects an OpenAI account with sticky session support
func (s *OpenAIGatewayService) SelectAccount(ctx context.Context, groupID *int64, sessionHash string) (*Account, error) {
	return s.SelectAccountForModel(ctx, groupID, sessionHash, "")
placeholder

// SelectAccountForModel selects an account supporting the requested model
func (s *OpenAIGatewayService) SelectAccountForModel(ctx context.Context, groupID *int64, sessionHash string, requestedModel string) (*Account, error) {
	return s.SelectAccountForModelWithExclusions(ctx, groupID, sessionHash, requestedModel, nil)
placeholder

// SelectAccountForModelWithExclusions selects an account supporting the requested model while excluding specified accounts.
// SelectAccountForModelWithExclusions 选择支持指定模型的账号，同时排除指定的账号。
func (s *OpenAIGatewayService) SelectAccountForModelWithExclusions(ctx context.Context, groupID *int64, sessionHash string, requestedModel string, excludedIDs map[int64]struct{placeholder) (*Account, error) {
	return s.selectAccountForModelWithExclusions(s.withOpenAIQuotaAutoPauseContext(ctx), groupID, PlatformOpenAI, sessionHash, requestedModel, excludedIDs, false, 0, "", false)
placeholder

// SelectAccountForTokenCount selects an account for a non-billable token-count
// request. It applies the normal platform, model, capability, and runtime
// eligibility checks without acquiring or waiting for a generation slot.
func (s *OpenAIGatewayService) SelectAccountForTokenCount(
	ctx context.Context,
	groupID *int64,
	sessionHash string,
	requestedModel string,
	requiredCapability OpenAIEndpointCapability,
	platform string,
) (*Account, error) {
	ctx = WithOpenAIProfitControlSuppressed(ctx)
	ctx = s.withOpenAIQuotaAutoPauseContext(ctx)
	return s.selectAccountForModelWithExclusions(
		ctx,
		groupID,
		platform,
		sessionHash,
		requestedModel,
		nil,
		false,
		0,
		requiredCapability,
		false,
	)
placeholder

// NormalizeOpenAICompatiblePlatform 保留 grok 与国产 OpenAI 兼容供应商（kimi/zhipu/
// deepseek）的原值，其他值一律归一为 openai。调度器据此对账号与请求做精确平台匹配：
// kimi 分组请求只命中 kimi 账号，语义与 openai/grok 一致。
// （upstream 曾将本函数改为未导出 normalizeOpenAICompatiblePlatform，本分支的
// handler 调度入口仍需导出，保持导出名。）
func NormalizeOpenAICompatiblePlatform(platform string) string {
	switch platform {
	case PlatformGrok, PlatformKimi, PlatformZhipu, PlatformDeepseek:
		return platform
	default:
		return PlatformOpenAI
placeholder
placeholder

// noAvailableOpenAISelectionError builds the standard "no account available" error
// while preserving the legacy /responses/compact error when applicable.
// details carries an optional machine-parseable exclusion summary (e.g.
// "pool=2, filtered: quota_auto_pause_7d=1 runtime_blocked=1") appended in
// parentheses. It is for server-side logs / ops diagnostics only: handlers
// never forward this error text to OpenAI-platform clients (they respond with
// the generic classification message). Callers that must preserve the legacy
// message pass "".
func noAvailableOpenAISelectionError(requestedModel string, compactBlocked bool, details string) error {
	if compactBlocked {
		return ErrNoAvailableCompactAccounts
placeholder
	message := "no available OpenAI accounts"
	if requestedModel != "" {
		message = fmt.Sprintf("no available OpenAI accounts supporting model: %s", requestedModel)
placeholder
	if details != "" {
		message += " (" + details + ")"
placeholder
	return openAINoAvailableSelectionError{message: messageplaceholder
placeholder

type openAINoAvailableSelectionError struct {
	message string
placeholder

func (e openAINoAvailableSelectionError) Error() string {
	return e.message
placeholder

func (e openAINoAvailableSelectionError) Unwrap() error {
	return ErrNoAvailableAccounts
placeholder

// openAICompactSupportTier classifies an OpenAI-compatible account by compact capability.
// 0 = explicitly unsupported, 1 = unknown / not yet probed, 2 = explicitly supported.
func openAICompactSupportTier(account *Account) int {
	if account == nil {
		return 0
placeholder
	if account.IsGrok() {
		return 2
placeholder
	if !account.IsOpenAI() {
		return 0
placeholder
	supported, known := account.OpenAICompactSupportKnown()
	if !known {
		return 1
placeholder
	if supported {
		return 2
placeholder
	return 0
placeholder

// isOpenAICompatibleAccountEligibleForRequest 判断 OpenAI 兼容账号是否满足本次请求的调度条件。
// 检查内容包括：平台匹配、账号可用性、quota 自动暂停、spark 路由限制、模型支持及端点能力。
//
// 注意：对 spark 影子账号，调用方还须额外调用 parentHealthyForShadow(account, lookup)
// 检查母账号凭据可用性；该检查未内置于本函数，以避免注入 DB 依赖。
func isOpenAICompatibleAccountEligibleForRequest(ctx context.Context, account *Account, platform string, requestedModel string, requireCompact bool, requiredCapability OpenAIEndpointCapability) bool {
	if !isOpenAICompatibleAccountEligibleForRequestBeforeProfit(ctx, account, platform, requestedModel, requireCompact, requiredCapability) {
		return false
placeholder
	// 分组利润控制：legacy 引擎的粘性/候选循环与 DB recheck 共用
	// 本判定，任何 fallback 都不能把利润不合格账号重新放回候选。
	if vetoed, _ := openAIProfitControlVetoReason(ctx, account); vetoed {
		return false
placeholder
	return true
placeholder

// isOpenAICompatibleAccountEligibleForRequestBeforeProfit applies every
// ordinary scheduling gate. Legacy selection uses it before classifying the
// profit veto so earlier failures retain their actual reason.
func isOpenAICompatibleAccountEligibleForRequestBeforeProfit(ctx context.Context, account *Account, platform string, requestedModel string, requireCompact bool, requiredCapability OpenAIEndpointCapability) bool {
	platform = NormalizeOpenAICompatiblePlatform(platform)
	if account == nil || account.Platform != platform || !account.IsOpenAICompatible() || !account.IsSchedulableForModelWithContext(ctx, requestedModel) {
		return false
placeholder
	if account.IsOpenAI() {
		if paused, reason := shouldAutoPauseOpenAIAccountByQuota(ctx, account); paused {
			// Debug level: this fires per-candidate on the scheduling hot path, so Info
			// would amplify into log spam once several accounts cross the threshold.
			slog.Debug("account_auto_paused_by_quota",
				"account_id", account.ID,
				"window", reason.window,
				"threshold", reason.threshold,
				"utilization", reason.utilization,
			)
			return false
	placeholder
placeholder
	if account.IsGrok() {
		if paused, reason := shouldAutoPauseGrokAccountByQuota(account); paused {
			slog.Debug("grok_account_auto_paused_by_quota",
				"account_id", account.ID,
				"window", reason.window,
				"threshold", reason.threshold,
				"utilization", reason.utilization,
			)
			return false
	placeholder
placeholder
	if requestedModel != "" && !account.IsModelSupported(requestedModel) {
		return false
placeholder
	if !account.SupportsOpenAIEndpointCapability(requiredCapability) {
		if account.IsGrok() && requiredCapability == OpenAIEndpointCapabilityGrokMediaGeneration {
			_, reason := account.GrokMediaGenerationEligibility()
			slog.Debug("grok_media_account_ineligible", "account_id", account.ID, "reason", reason)
	placeholder
		return false
placeholder
	if requireCompact && openAICompactSupportTier(account) == 0 {
		return false
placeholder
	return true
placeholder

type openAIQuotaAutoPauseDecision struct {
	window      string
	threshold   float64
	utilization float64
placeholder

func shouldAutoPauseGrokAccountByQuota(account *Account) (bool, openAIQuotaAutoPauseDecision) {
	if account == nil || !account.IsGrok() || account.Type != AccountTypeOAuth {
		return false, openAIQuotaAutoPauseDecision{placeholder
placeholder
	snapshot, err := grokQuotaSnapshotFromExtra(account.Extra)
	if err != nil || snapshot == nil {
		return false, openAIQuotaAutoPauseDecision{placeholder
placeholder
	now := time.Now()
	if grokQuotaSnapshotStaleForPause(snapshot, now) {
		return false, openAIQuotaAutoPauseDecision{placeholder
placeholder
	if grokQuotaRetryAfterActive(snapshot, now) {
		return true, openAIQuotaAutoPauseDecision{window: "retry_after", threshold: 1, utilization: 1placeholder
placeholder
	if paused, decision := shouldAutoPauseGrokQuotaWindow("requests", snapshot.Requests, now); paused {
		return true, decision
placeholder
	if paused, decision := shouldAutoPauseGrokQuotaWindow("tokens", snapshot.Tokens, now); paused {
		return true, decision
placeholder
	return false, openAIQuotaAutoPauseDecision{placeholder
placeholder

func grokQuotaRetryAfterActive(snapshot *xai.QuotaSnapshot, now time.Time) bool {
	if snapshot == nil || snapshot.RetryAfterSeconds == nil || *snapshot.RetryAfterSeconds <= 0 {
		return false
placeholder
	if strings.TrimSpace(snapshot.UpdatedAt) == "" {
		return true
placeholder
	updatedAt, err := parseTime(snapshot.UpdatedAt)
	if err != nil {
		return true
placeholder
	retryAfterUntil := updatedAt.Add(time.Duration(*snapshot.RetryAfterSeconds) * time.Second)
	return now.Before(retryAfterUntil)
placeholder

func shouldAutoPauseGrokQuotaWindow(name string, window *xai.QuotaWindow, now time.Time) (bool, openAIQuotaAutoPauseDecision) {
	if window == nil || window.Limit == nil || window.Remaining == nil || *window.Limit <= 0 {
		return false, openAIQuotaAutoPauseDecision{placeholder
placeholder
	if window.ResetUnix != nil && *window.ResetUnix > 0 && !now.Before(time.Unix(*window.ResetUnix, 0)) {
		return false, openAIQuotaAutoPauseDecision{placeholder
placeholder
	utilization := float64(*window.Limit-*window.Remaining) / float64(*window.Limit)
	if *window.Remaining <= 0 || utilization >= 1 {
		return true, openAIQuotaAutoPauseDecision{window: name, threshold: 1, utilization: utilizationplaceholder
placeholder
	return false, openAIQuotaAutoPauseDecision{placeholder
placeholder

func grokQuotaSnapshotStaleForPause(snapshot *xai.QuotaSnapshot, now time.Time) bool {
	if snapshot == nil || strings.TrimSpace(snapshot.UpdatedAt) == "" {
		return false
placeholder
	updatedAt, err := parseTime(snapshot.UpdatedAt)
	if err != nil {
		return false
placeholder
	return now.Sub(updatedAt) >= openAICodexAutoPauseStaleAfter
placeholder

func shouldAutoPauseOpenAIAccountByQuota(ctx context.Context, account *Account) (bool, openAIQuotaAutoPauseDecision) {
	if account == nil || !account.IsOpenAI() {
		return false, openAIQuotaAutoPauseDecision{placeholder
placeholder
	// Per-account explicit-disable flags must take precedence over the global default.
	// Without these, leaving the account threshold blank means "use global default",
	// so an admin has no way to exempt a single account from auto-pause once a global
	// default exists. The disable flag is per-window so an account can opt out of
	// only 5h or only 7d auto-pause.
	disabled5h := resolveAccountExtraBool(account.Extra, "auto_pause_5h_disabled")
	disabled7d := resolveAccountExtraBool(account.Extra, "auto_pause_7d_disabled")
	threshold5h, threshold7d := resolveOpenAIQuotaAutoPauseThresholds(ctx, account)
	now := time.Now()
	if !disabled5h && threshold5h > 0 {
		if utilization, ok := resolveOpenAIQuotaUtilization(account.Extra, "5h", now); ok && utilization >= threshold5h {
			return true, openAIQuotaAutoPauseDecision{window: "5h", threshold: threshold5h, utilization: utilizationplaceholder
	placeholder
placeholder
	if !disabled7d && threshold7d > 0 {
		if utilization, ok := resolveOpenAIQuotaUtilization(account.Extra, "7d", now); ok && utilization >= threshold7d {
			return true, openAIQuotaAutoPauseDecision{window: "7d", threshold: threshold7d, utilization: utilizationplaceholder
	placeholder
placeholder
	return false, openAIQuotaAutoPauseDecision{placeholder
placeholder

// resolveAccountExtraBool reads a bool-like value from account extra, tolerating
// the few shapes JSON unmarshalling may produce (real bool, "true"/"false"
// strings, 0/1 numbers).
func resolveAccountExtraBool(extra map[string]any, key string) bool {
	if len(extra) == 0 {
		return false
placeholder
	value, ok := extra[key]
	if !ok || value == nil {
		return false
placeholder
	switch v := value.(type) {
	case bool:
		return v
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(v))
		return err == nil && parsed
	case float64:
		return v != 0
	case float32:
		return v != 0
	case int:
		return v != 0
	case int64:
		return v != 0
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i != 0
	placeholder
placeholder
	return false
placeholder

func resolveOpenAIQuotaAutoPauseThresholds(ctx context.Context, account *Account) (float64, float64) {
	threshold5h, _ := resolveAccountExtraNumber(account.Extra, "auto_pause_5h_threshold")
	threshold7d, _ := resolveAccountExtraNumber(account.Extra, "auto_pause_7d_threshold")
	threshold5h = clamp01(threshold5h)
	threshold7d = clamp01(threshold7d)
	if threshold5h > 0 && threshold7d > 0 {
		return threshold5h, threshold7d
placeholder
	settings := openAIQuotaAutoPauseSettingsFromContext(ctx)
	if threshold5h <= 0 {
		threshold5h = clamp01(settings.DefaultThreshold5h)
placeholder
	if threshold7d <= 0 {
		threshold7d = clamp01(settings.DefaultThreshold7d)
placeholder
	return threshold5h, threshold7d
placeholder

func resolveAccountExtraNumber(extra map[string]any, keys ...string) (float64, bool) {
	if len(extra) == 0 {
		return 0, false
placeholder
	for _, key := range keys {
		value, ok := extra[key]
		if !ok || value == nil {
			continue
	placeholder
		switch v := value.(type) {
		case float64:
			return v, true
		case float32:
			return float64(v), true
		case int:
			return float64(v), true
		case int64:
			return float64(v), true
		case json.Number:
			parsed, err := v.Float64()
			if err == nil {
				return parsed, true
		placeholder
		case string:
			parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err == nil {
				return parsed, true
		placeholder
	placeholder
placeholder
	return 0, false
placeholder

// resolveOpenAIQuotaUtilization returns the current utilization ratio (0..1) for the
// given Codex usage window. ok=false means there is no usable signal to pause on:
// either no snapshot exists, or the window has already rolled over so the cached
// percentage is stale. The stale guard matters because a paused account stops
// receiving requests, so its snapshot is never refreshed from upstream headers —
// without this check an old used_percent would keep the account paused forever even
// after the real window reset.
func resolveOpenAIQuotaUtilization(extra map[string]any, window string, now time.Time) (float64, bool) {
	usedPercent := readOpenAIQuotaUsedPercent(extra, window)
	if usedPercent <= 0 {
		return 0, false
placeholder
	if openAIQuotaWindowReset(extra, window, now) {
		return 0, false
placeholder
	// 快照过于陈旧（账号长期未收到流量刷新）时，不再据此暂停。放行后下一次响应头
	// 会刷新快照实现自愈，避免账号在错误/过期的 used% 上被永久跳过（issue #2994）。
	if openAICodexSnapshotStaleForPause(extra, now) {
		return 0, false
placeholder
	return usedPercent / 100, true
placeholder

// openAICodexSnapshotStaleForPause reports whether the Codex usage snapshot is stale
// enough that it should no longer keep an account auto-paused. It anchors on
// codex_usage_updated_at (always written by buildCodexUsageExtraUpdates). A missing or
// unparseable timestamp returns false (treated as fresh, so the account stays paused) —
// this is deliberate: it prevents any snapshot without a write time from silently escaping
// auto-pause, and a genuinely-exhausted account that is actively served refreshes the
// timestamp on every response so it never crosses the staleness bound.
func openAICodexSnapshotStaleForPause(extra map[string]any, now time.Time) bool {
	if len(extra) == 0 {
		return false
placeholder
	updatedRaw, ok := extra["codex_usage_updated_at"]
	if !ok {
		return false
placeholder
	updatedAt, err := parseTime(fmt.Sprint(updatedRaw))
	if err != nil {
		return false
placeholder
	return now.Sub(updatedAt) >= openAICodexAutoPauseStaleAfter
placeholder

// openAIQuotaWindowReset reports whether the Codex usage window's reset time has
// already passed relative to now. It prefers the absolute codex_<window>_reset_at
// timestamp and falls back to codex_<window>_reset_after_seconds anchored at
// codex_usage_updated_at, mirroring AccountUsageService's window-progress logic.
func openAIQuotaWindowReset(extra map[string]any, window string, now time.Time) bool {
	if len(extra) == 0 {
		return false
placeholder
	if resetAtRaw, ok := extra["codex_"+window+"_reset_at"]; ok {
		if resetAt, err := parseTime(fmt.Sprint(resetAtRaw)); err == nil {
			return !now.Before(resetAt)
	placeholder
placeholder
	resetAfter := parseExtraInt(extra["codex_"+window+"_reset_after_seconds"])
	if resetAfter <= 0 {
		return false
placeholder
	base := now
	if updatedRaw, ok := extra["codex_usage_updated_at"]; ok {
		if updatedAt, err := parseTime(fmt.Sprint(updatedRaw)); err == nil {
			base = updatedAt
	placeholder
placeholder
	resetAt := base.Add(time.Duration(resetAfter) * time.Second)
	return !now.Before(resetAt)
placeholder

func readOpenAIQuotaUsedPercent(extra map[string]any, window string) float64 {
	if len(extra) == 0 {
		return 0
placeholder
	if value, ok := resolveAccountExtraNumber(extra, "codex_"+window+"_used_percent"); ok {
		return value
placeholder
	return 0
placeholder

type openAIQuotaAutoPauseCtxKey struct{placeholder

func withOpenAIQuotaAutoPauseSettings(ctx context.Context, settings OpsOpenAIAccountQuotaAutoPauseSettings) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	return context.WithValue(ctx, openAIQuotaAutoPauseCtxKey{placeholder, settings)
placeholder

func openAIQuotaAutoPauseSettingsFromContext(ctx context.Context) OpsOpenAIAccountQuotaAutoPauseSettings {
	if ctx == nil {
		return OpsOpenAIAccountQuotaAutoPauseSettings{placeholder
placeholder
	settings, _ := ctx.Value(openAIQuotaAutoPauseCtxKey{placeholder).(OpsOpenAIAccountQuotaAutoPauseSettings)
	return settings
placeholder

func (s *OpenAIGatewayService) withOpenAIQuotaAutoPauseContext(ctx context.Context) context.Context {
	if s == nil || s.settingService == nil {
		return ctx
placeholder
	return withOpenAIQuotaAutoPauseSettings(ctx, s.settingService.GetOpenAIQuotaAutoPauseSettings(ctx))
placeholder

// prioritizeOpenAICompactAccounts re-orders a slice so that accounts with known
// compact support are tried first, followed by unknown, then explicitly unsupported.
// The relative order within each tier is preserved.
func prioritizeOpenAICompactAccounts(accounts []*Account) []*Account {
	if len(accounts) == 0 {
		return nil
placeholder
	supported := make([]*Account, 0, len(accounts))
	unknown := make([]*Account, 0, len(accounts))
	unsupported := make([]*Account, 0, len(accounts))
	for _, account := range accounts {
		switch openAICompactSupportTier(account) {
		case 2:
			supported = append(supported, account)
		case 1:
			unknown = append(unknown, account)
		default:
			unsupported = append(unsupported, account)
	placeholder
placeholder
	out := make([]*Account, 0, len(accounts))
	out = append(out, supported...)
	out = append(out, unknown...)
	out = append(out, unsupported...)
	return out
placeholder

// resolveOpenAIAccountUpstreamModelForRequest resolves the upstream model that
// would be sent for a given request, honoring the legacy compact-only mapping
// when the caller is on the /responses/compact path.
func resolveOpenAIAccountUpstreamModelForRequest(account *Account, requestedModel string, requireCompact bool) string {
	// Forward checks the raw Chat Completions fallback before passthrough.
	// These API-key accounts therefore apply normal account model_mapping and
	// upstream normalization, but never compact_model_mapping.
	if shouldForwardOpenAIResponsesViaRawChatCompletions(account) {
		upstreamModel := resolveOpenAIForwardModel(account, requestedModel, "")
		return normalizeOpenAIModelForUpstream(account, upstreamModel)
placeholder

	// Passthrough accounts only replace authentication. Their Forward path
	// keeps the channel-mapped model in the request body and does not apply the
	// account's normal model_mapping. Legacy /responses/compact is the one
	// exception: forwardOpenAIPassthrough applies compact_model_mapping
	// directly to that channel-mapped model.
	if account != nil && account.IsOpenAIPassthroughEnabled() {
		upstreamModel := strings.TrimSpace(requestedModel)
		if upstreamModel == "" {
			return ""
	placeholder
		if requireCompact {
			return resolveOpenAICompactForwardModel(account, upstreamModel)
	placeholder
		return upstreamModel
placeholder

	// Compact mappings are keyed by the client-visible model. Prefer an exact
	// compact rule before ordinary account mapping; otherwise a normal alias can
	// hide the compact-specific rule and make scheduling disagree with Forward.
	if requireCompact && account != nil {
		if compactModel, matched := account.ResolveCompactMappedModel(strings.TrimSpace(requestedModel)); matched {
			if compactModel = strings.TrimSpace(compactModel); compactModel != "" {
				return compactModel
		placeholder
	placeholder
placeholder

	upstreamModel := resolveOpenAIForwardModel(account, requestedModel, "")
	if upstreamModel == "" {
		return ""
placeholder
	if requireCompact {
		compactModel := resolveOpenAICompactForwardModel(account, upstreamModel)
		if compactModel != upstreamModel {
			return compactModel
	placeholder
placeholder
	return normalizeOpenAIModelForUpstream(account, upstreamModel)
placeholder

// ResolveOpenAIAccountUpstreamModelForRequest exposes the scheduler's exact
// account mapping chain to handler-side outcome reporting.
func ResolveOpenAIAccountUpstreamModelForRequest(account *Account, requestedModel string, requireCompact bool) string {
	return resolveOpenAIAccountUpstreamModelForRequest(account, requestedModel, requireCompact)
placeholder

// resolveOpenAIForwardMappedModels is the shared account mapping chain for
// Forward callers. billingModel retains the ordinary mapping used for usage
// accounting, while upstreamModel is the model the scheduler has admitted.
func resolveOpenAIForwardMappedModels(account *Account, requestedModel string, requireCompact bool) (billingModel, upstreamModel string) {
	requestedModel = strings.TrimSpace(requestedModel)
	if account != nil && account.IsOpenAIPassthroughEnabled() {
		billingModel = requestedModel
placeholder else if account != nil {
		billingModel = strings.TrimSpace(account.GetMappedModel(requestedModel))
placeholder
	if billingModel == "" {
		billingModel = requestedModel
placeholder
	upstreamModel = resolveOpenAIAccountUpstreamModelForRequest(account, requestedModel, requireCompact)
	if strings.TrimSpace(upstreamModel) == "" {
		upstreamModel = billingModel
placeholder
	return billingModel, upstreamModel
placeholder

func resolveOpenAIErrorSchedulingModel(billingModel, upstreamModel string) string {
	if upstreamModel = strings.TrimSpace(upstreamModel); upstreamModel != "" {
		return upstreamModel
placeholder
	return strings.TrimSpace(billingModel)
placeholder

func (s *OpenAIGatewayService) selectAccountForModelWithExclusions(ctx context.Context, groupID *int64, platform string, sessionHash string, requestedModel string, excludedIDs map[int64]struct{placeholder, requireCompact bool, stickyAccountID int64, requiredCapability OpenAIEndpointCapability, preferLowUpstreamRate bool) (*Account, error) {
	platform = NormalizeOpenAICompatiblePlatform(platform)
	if s.checkChannelPricingRestriction(ctx, groupID, requestedModel) {
		slog.Warn("channel pricing restriction blocked request",
			"group_id", derefGroupID(groupID),
			"model", requestedModel)
		return nil, fmt.Errorf("%w supporting model: %s (channel pricing restriction)", ErrNoAvailableAccounts, requestedModel)
placeholder

	// 1. 尝试粘性会话命中
	// Try sticky session hit
	if account := s.tryStickySessionHit(ctx, groupID, platform, sessionHash, requestedModel, excludedIDs, requireCompact, stickyAccountID, requiredCapability); account != nil {
		return account, nil
placeholder

	// 2. 获取可调度的 OpenAI 账号
	// Get schedulable OpenAI accounts
	accounts, err := s.listSchedulableAccounts(ctx, groupID, platform)
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder

	// 3. 按优先级 + LRU 选择最佳账号
	// Select by priority + LRU
	selected, compactBlocked, filterStats := s.selectBestAccount(ctx, groupID, platform, accounts, requestedModel, excludedIDs, requireCompact, requiredCapability, preferLowUpstreamRate)

	if selected == nil {
		return nil, noAvailableOpenAISelectionError(requestedModel, compactBlocked, filterStats.summary(""))
placeholder

	hydrated, err := s.hydrateSelectedAccount(ctx, selected)
	if err != nil {
		return nil, err
placeholder

	// 4. 设置粘性会话绑定（利润门下推迟到 handler 终检通过后再绑定，
	// 终检否决的账号不得成为新的粘性目标；无门保持既有 eager 绑定与 TTL）
	// Set sticky session binding (deferred until terminal admission under a profit gate)
	if sessionHash != "" && !gatewayProfitControlGateActive(ctx) {
		_ = s.setStickySessionAccountID(ctx, groupID, sessionHash, selected.ID, openaiStickySessionTTL)
placeholder

	return hydrated, nil
placeholder

// tryStickySessionHit 尝试从粘性会话获取账号。
// 如果命中且账号可用则返回账号；如果账号不可用则清理会话并返回 nil。
//
// tryStickySessionHit attempts to get account from sticky session.
// Returns account if hit and usable; clears session and returns nil if account is unavailable.
func (s *OpenAIGatewayService) tryStickySessionHit(ctx context.Context, groupID *int64, platform string, sessionHash, requestedModel string, excludedIDs map[int64]struct{placeholder, requireCompact bool, stickyAccountID int64, requiredCapability OpenAIEndpointCapability) *Account {
	if sessionHash == "" {
		return nil
placeholder
	platform = NormalizeOpenAICompatiblePlatform(platform)

	accountID := stickyAccountID
	if accountID <= 0 {
		var err error
		accountID, err = s.getStickySessionAccountID(ctx, groupID, sessionHash)
		if err != nil || accountID <= 0 {
			return nil
	placeholder
placeholder

	if _, excluded := excludedIDs[accountID]; excluded {
		return nil
placeholder

	account, err := s.getSchedulableAccount(ctx, accountID)
	if err != nil {
		return nil
placeholder

	// 检查账号是否需要清理粘性会话
	// Check if sticky session should be cleared
	if shouldClearStickySession(account, requestedModel) {
		_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
		return nil
placeholder

	// 验证账号是否可用于当前请求
	// Verify account is usable for current request
	if !isOpenAICompatibleAccountEligibleForRequest(ctx, account, platform, requestedModel, false, requiredCapability) {
		return nil
placeholder
	if !parentHealthyForShadow(account, s.parentAccountLookup(ctx)) {
		_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
		return nil
placeholder
	if s.isOpenAIAccountRequestRuntimeBlocked(account, requestedModel) {
		_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
		return nil
placeholder
	account = s.recheckSelectedOpenAIAccountFromDB(ctx, account, groupID, platform, requestedModel, requireCompact, requiredCapability)
	if account == nil || !s.openAIAccountMatchesSchedulingGroup(account, groupID) {
		_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
		return nil
placeholder
	if groupID != nil && s.needsUpstreamChannelRestrictionCheck(ctx, groupID) &&
		s.isUpstreamModelRestrictedByChannel(ctx, *groupID, account, requestedModel, requireCompact) {
		_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
		return nil
placeholder

	// 刷新会话 TTL 并返回账号
	// Refresh session TTL and return account
	_ = s.refreshStickySessionTTL(ctx, groupID, sessionHash, openaiStickySessionTTL)
	return account
placeholder

// selectBestAccount 从候选账号中选择最佳账号（优先级 + LRU）。
// 返回 nil 表示无可用账号。
//
// selectBestAccount selects the best account from candidates (priority + LRU).
// Returns nil if no available account. The second return reports whether at
// least one candidate was filtered out solely because it lacks compact support
// (only meaningful when the legacy /responses/compact requireCompact flag is
// true); the third contains deterministic
// exclusion diagnostics for the evaluated snapshot.
func (s *OpenAIGatewayService) selectBestAccount(ctx context.Context, groupID *int64, platform string, accounts []Account, requestedModel string, excludedIDs map[int64]struct{placeholder, requireCompact bool, requiredCapability OpenAIEndpointCapability, preferLowUpstreamRate bool) (*Account, bool, openAISelectionFilterStats) {
	platform = NormalizeOpenAICompatiblePlatform(platform)
	compactBlocked := false
	filterStats := openAISelectionFilterStats{pool: len(accounts)placeholder
	needsUpstreamCheck := s.needsUpstreamChannelRestrictionCheck(ctx, groupID)
	eligible := make([]*Account, 0, len(accounts))
	compactTiers := make(map[int64]int, len(accounts))

	for i := range accounts {
		acc := &accounts[i]

		// 跳过被排除的账号
		// Skip excluded accounts
		if _, excluded := excludedIDs[acc.ID]; excluded {
			filterStats.exclude("excluded")
			continue
	placeholder

		fresh := s.resolveFreshSchedulableOpenAIAccountBeforeProfit(ctx, acc, platform, requestedModel, false, requiredCapability)
		if fresh == nil {
			filterStats.exclude("ineligible")
			continue
	placeholder
		fresh = s.recheckSelectedOpenAIAccountFromDBBeforeProfit(ctx, fresh, groupID, platform, requestedModel, false, requiredCapability)
		if fresh == nil {
			filterStats.exclude("ineligible")
			continue
	placeholder
		if needsUpstreamCheck && s.isUpstreamModelRestrictedByChannel(ctx, *groupID, fresh, requestedModel, requireCompact) {
			filterStats.exclude("channel_restricted")
			continue
	placeholder
		if vetoed, reason := openAIProfitControlVetoReason(ctx, fresh); vetoed {
			filterStats.exclude(reason)
			continue
	placeholder
		compactTier := 0
		if requireCompact {
			compactTier = openAICompactSupportTier(fresh)
			if compactTier == 0 {
				compactBlocked = true
				filterStats.exclude("compact_unsupported")
				continue
		placeholder
	placeholder

		eligible = append(eligible, fresh)
		compactTiers[fresh.ID] = compactTier
placeholder

	if len(eligible) == 0 {
		return nil, compactBlocked, filterStats
placeholder
	rateOrder := openAILegacyUpstreamRateOrder{placeholder
	if preferLowUpstreamRate {
		rateOrder = newOpenAILegacyUpstreamRateOrder(eligible, time.Now(), s.openAIOAuthSchedulingRateMultiplier(ctx))
placeholder
	sort.SliceStable(eligible, func(i, j int) bool {
		a, b := eligible[i], eligible[j]
		if requireCompact && compactTiers[a.ID] != compactTiers[b.ID] {
			return compactTiers[a.ID] > compactTiers[b.ID]
	placeholder
		if rateCmp := rateOrder.compare(a, b); rateCmp != 0 {
			return rateCmp < 0
	placeholder
		return s.isBetterAccount(a, b)
placeholder)
	return eligible[0], compactBlocked, filterStats
placeholder

// isBetterAccount 判断 candidate 是否比 current 更优。
// 规则：优先级更高（数值更小）优先；同优先级时，未使用过的优先，其次是最久未使用的。
//
// isBetterAccount checks if candidate is better than current.
// Rules: higher priority (lower value) wins; same priority: never used > least recently used.
func (s *OpenAIGatewayService) isBetterAccount(candidate, current *Account) bool {
	// 优先级更高（数值更小）
	// Higher priority (lower value)
	if candidate.Priority < current.Priority {
		return true
placeholder
	if candidate.Priority > current.Priority {
		return false
placeholder

	// 同优先级，比较最后使用时间
	// Same priority, compare last used time
	switch {
	case candidate.LastUsedAt == nil && current.LastUsedAt != nil:
		// candidate 从未使用，优先
		return true
	case candidate.LastUsedAt != nil && current.LastUsedAt == nil:
		// current 从未使用，保持
		return false
	case candidate.LastUsedAt == nil && current.LastUsedAt == nil:
		// 都未使用，保持
		return false
	default:
		// 都使用过，选择最久未使用的
		return candidate.LastUsedAt.Before(*current.LastUsedAt)
placeholder
placeholder

// SelectAccountWithLoadAwareness selects an account with load-awareness and wait plan.
func (s *OpenAIGatewayService) SelectAccountWithLoadAwareness(ctx context.Context, groupID *int64, sessionHash string, requestedModel string, excludedIDs map[int64]struct{placeholder) (*AccountSelectionResult, error) {
	ctx = s.withOpenAIQuotaAutoPauseContext(ctx)
	ctx = s.withOpenAIGroupPrivacyRequirement(ctx, groupID)
	// 分组利润控制：legacy 公共入口同样装门，保证不经
	// selectAccountWithScheduler 的调用方也无法绕过利润准入。
	ctx = s.withOpenAIProfitControlGate(ctx, groupID)
	return s.selectAccountWithLoadAwareness(ctx, groupID, PlatformOpenAI, sessionHash, requestedModel, excludedIDs, false, "", true)
placeholder

func (s *OpenAIGatewayService) selectAccountWithLoadAwareness(ctx context.Context, groupID *int64, platform string, sessionHash string, requestedModel string, excludedIDs map[int64]struct{placeholder, requireCompact bool, requiredCapability OpenAIEndpointCapability, useUpstreamTokenCost bool) (*AccountSelectionResult, error) {
	platform = NormalizeOpenAICompatiblePlatform(platform)
	if s.checkChannelPricingRestriction(ctx, groupID, requestedModel) {
		slog.Warn("channel pricing restriction blocked request",
			"group_id", derefGroupID(groupID),
			"model", requestedModel)
		return nil, fmt.Errorf("%w supporting model: %s (channel pricing restriction)", ErrNoAvailableAccounts, requestedModel)
placeholder

	cfg := s.schedulingConfig()
	preferLowUpstreamRate := useUpstreamTokenCost && s.isOpenAILowUpstreamRatePriorityEnabled(ctx)
	needsUpstreamCheck := s.needsUpstreamChannelRestrictionCheck(ctx, groupID)
	var stickyAccountID int64
	if sessionHash != "" && s.cache != nil {
		if accountID, err := s.getStickySessionAccountID(ctx, groupID, sessionHash); err == nil {
			stickyAccountID = accountID
	placeholder
placeholder
	if s.concurrencyService == nil || !cfg.LoadBatchEnabled {
		account, err := s.selectAccountForModelWithExclusions(ctx, groupID, platform, sessionHash, requestedModel, excludedIDs, requireCompact, stickyAccountID, requiredCapability, preferLowUpstreamRate)
		if err != nil {
			return nil, err
	placeholder
		result, err := s.tryAcquireAccountSlot(ctx, account.ID, account.Concurrency)
		if err == nil && result != nil && result.Acquired {
			return s.newAcquiredSelectionResult(ctx, account, result.ReleaseFunc)
	placeholder
		if stickyAccountID > 0 && stickyAccountID == account.ID && s.concurrencyService != nil {
			waitingCount, _ := s.concurrencyService.GetAccountWaitingCount(ctx, account.ID)
			if waitingCount < cfg.StickySessionMaxWaiting {
				return s.newSelectionResult(ctx, account, false, nil, &AccountWaitPlan{
					AccountID:      account.ID,
					MaxConcurrency: account.Concurrency,
					Timeout:        cfg.StickySessionWaitTimeout,
					MaxWaiting:     cfg.StickySessionMaxWaiting,
			placeholder)
		placeholder
	placeholder
		return s.newSelectionResult(ctx, account, false, nil, &AccountWaitPlan{
			AccountID:      account.ID,
			MaxConcurrency: account.Concurrency,
			Timeout:        cfg.FallbackWaitTimeout,
			MaxWaiting:     cfg.FallbackMaxWaiting,
	placeholder)
placeholder

	accounts, err := s.listSchedulableAccounts(ctx, groupID, platform)
	if err != nil {
		return nil, err
placeholder
	if len(accounts) == 0 {
		return nil, ErrNoAvailableAccounts
placeholder

	isExcluded := func(accountID int64) bool {
		if excludedIDs == nil {
			return false
	placeholder
		_, excluded := excludedIDs[accountID]
		return excluded
placeholder

	// ============ Layer 1: Sticky session ============
	if sessionHash != "" {
		accountID := stickyAccountID
		if accountID > 0 && !isExcluded(accountID) {
			account, err := s.getSchedulableAccount(ctx, accountID)
			if err == nil {
				clearSticky := shouldClearStickySession(account, requestedModel)
				if clearSticky {
					_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
			placeholder
				if !clearSticky && isOpenAICompatibleAccountEligibleForRequest(ctx, account, platform, requestedModel, false, requiredCapability) {
					account = s.recheckSelectedOpenAIAccountFromDB(ctx, account, groupID, platform, requestedModel, requireCompact, requiredCapability)
					if account == nil {
						_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
				placeholder else if !s.openAIAccountMatchesSchedulingGroup(account, groupID) {
						_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
				placeholder else if s.isOpenAIAccountRequestRuntimeBlocked(account, requestedModel) {
						_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
				placeholder else if needsUpstreamCheck && s.isUpstreamModelRestrictedByChannel(ctx, *groupID, account, requestedModel, requireCompact) {
						_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
				placeholder else if !parentHealthyForShadow(account, s.parentAccountLookup(ctx)) {
						_ = s.deleteStickySessionAccountID(ctx, groupID, sessionHash)
				placeholder else {
						result, err := s.tryAcquireAccountSlot(ctx, accountID, account.Concurrency)
						if err == nil && result != nil && result.Acquired {
							selection, selectErr := s.newAcquiredSelectionResult(ctx, account, result.ReleaseFunc)
							if selectErr != nil {
								return nil, selectErr
						placeholder
							_ = s.refreshStickySessionTTL(ctx, groupID, sessionHash, openaiStickySessionTTL)
							return selection, nil
					placeholder

						waitingCount, _ := s.concurrencyService.GetAccountWaitingCount(ctx, accountID)
						if waitingCount < cfg.StickySessionMaxWaiting {
							return s.newSelectionResult(ctx, account, false, nil, &AccountWaitPlan{
								AccountID:      accountID,
								MaxConcurrency: account.Concurrency,
								Timeout:        cfg.StickySessionWaitTimeout,
								MaxWaiting:     cfg.StickySessionMaxWaiting,
						placeholder)
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// ============ Layer 2: Load-aware selection ============
	// Per-pass parent-health cache to avoid repeated DB calls when multiple shadow
	// accounts share the same parent.
	parentCacheL2 := make(map[int64]*Account)
	parentLookupL2 := func(id int64) *Account {
		if a, ok := parentCacheL2[id]; ok {
			return a
	placeholder
		if s.accountRepo == nil {
			return nil
	placeholder
		a, _ := s.accountRepo.GetByID(ctx, id)
		parentCacheL2[id] = a
		return a
placeholder
	baseCandidateCount := 0
	candidates := make([]*Account, 0, len(accounts))
	for i := range accounts {
		acc := &accounts[i]
		if isExcluded(acc.ID) {
			continue
	placeholder
		// Scheduler snapshots can be temporarily stale (bucket rebuild is throttled);
		// re-check schedulability here so recently rate-limited/overloaded accounts
		// are not selected again before the bucket is rebuilt.
		if !isOpenAICompatibleAccountEligibleForRequest(ctx, acc, platform, requestedModel, false, requiredCapability) {
			continue
	placeholder
		if !parentHealthyForShadow(acc, parentLookupL2) {
			continue
	placeholder
		if s.isOpenAIAccountRequestRuntimeBlocked(acc, requestedModel) {
			continue
	placeholder
		if needsUpstreamCheck && s.isUpstreamModelRestrictedByChannel(ctx, *groupID, acc, requestedModel, requireCompact) {
			continue
	placeholder
		baseCandidateCount++
		candidates = append(candidates, acc)
placeholder

	if len(candidates) == 0 {
		return nil, ErrNoAvailableAccounts
placeholder
	rateOrder := openAILegacyUpstreamRateOrder{placeholder
	if preferLowUpstreamRate {
		rateOrder = newOpenAILegacyUpstreamRateOrder(candidates, time.Now(), s.openAIOAuthSchedulingRateMultiplier(ctx))
placeholder

	accountLoads := make([]AccountWithConcurrency, 0, len(candidates))
	for _, acc := range candidates {
		accountLoads = append(accountLoads, AccountWithConcurrency{
			ID:             acc.ID,
			MaxConcurrency: acc.EffectiveLoadFactor(),
	placeholder)
placeholder

	tryAcquireFromLoadMap := func(loadMap map[int64]*AccountLoadInfo) (*AccountSelectionResult, bool, error) {
		var available []accountWithLoad
		for _, acc := range candidates {
			loadInfo := loadMap[acc.ID]
			if loadInfo == nil {
				loadInfo = &AccountLoadInfo{AccountID: acc.IDplaceholder
		placeholder
			if loadInfo.LoadRate < 100 {
				available = append(available, accountWithLoad{
					account:  acc,
					loadInfo: loadInfo,
			placeholder)
		placeholder
	placeholder

		if len(available) == 0 {
			return nil, false, nil
	placeholder

		sort.SliceStable(available, func(i, j int) bool {
			a, b := available[i], available[j]
			if a.account.Priority != b.account.Priority {
				return a.account.Priority < b.account.Priority
		placeholder
			if a.loadInfo.LoadRate != b.loadInfo.LoadRate {
				return a.loadInfo.LoadRate < b.loadInfo.LoadRate
		placeholder
			switch {
			case a.account.LastUsedAt == nil && b.account.LastUsedAt != nil:
				return true
			case a.account.LastUsedAt != nil && b.account.LastUsedAt == nil:
				return false
			case a.account.LastUsedAt == nil && b.account.LastUsedAt == nil:
				return false
			default:
				return a.account.LastUsedAt.Before(*b.account.LastUsedAt)
		placeholder
	placeholder)
		shuffleWithinSortGroups(available)
		if rateOrder.enabled {
			sort.SliceStable(available, func(i, j int) bool {
				return rateOrder.compare(available[i].account, available[j].account) < 0
		placeholder)
	placeholder

		selectionOrder := make([]accountWithLoad, 0, len(available))
		if requireCompact {
			appendTier := func(out []accountWithLoad, tier int) []accountWithLoad {
				for _, item := range available {
					if openAICompactSupportTier(item.account) == tier {
						out = append(out, item)
				placeholder
			placeholder
				return out
		placeholder
			selectionOrder = appendTier(selectionOrder, 2)
			selectionOrder = appendTier(selectionOrder, 1)
			// tier 0 候选作为兜底追加：DB recheck 时若发现 cache tier 0 实际
			// 已升级为 1/2（探测刚跑完，cache 尚未刷新），仍可正常命中。
			selectionOrder = appendTier(selectionOrder, 0)
	placeholder else {
			selectionOrder = append(selectionOrder, available...)
	placeholder

		for _, item := range selectionOrder {
			fresh := s.resolveFreshSchedulableOpenAIAccount(ctx, item.account, platform, requestedModel, false, requiredCapability)
			if fresh == nil {
				continue
		placeholder
			fresh = s.recheckSelectedOpenAIAccountFromDB(ctx, fresh, groupID, platform, requestedModel, requireCompact, requiredCapability)
			if fresh == nil {
				continue
		placeholder
			if needsUpstreamCheck && s.isUpstreamModelRestrictedByChannel(ctx, *groupID, fresh, requestedModel, requireCompact) {
				continue
		placeholder
			result, err := s.tryAcquireAccountSlot(ctx, fresh.ID, fresh.Concurrency)
			if err == nil && result != nil && result.Acquired {
				selection, selectErr := s.newAcquiredSelectionResult(ctx, fresh, result.ReleaseFunc)
				if selectErr != nil {
					return nil, true, selectErr
			placeholder
				if sessionHash != "" && !gatewayProfitControlGateActive(ctx) {
					_ = s.setStickySessionAccountID(ctx, groupID, sessionHash, fresh.ID, openaiStickySessionTTL)
			placeholder
				return selection, true, nil
		placeholder
	placeholder
		return nil, true, nil
placeholder

	loadMap, err := s.concurrencyService.GetAccountsLoadBatch(ctx, accountLoads)
	if err != nil {
		ordered := append([]*Account(nil), candidates...)
		sortAccountsByPriorityAndLastUsed(ordered, false)
		if rateOrder.enabled {
			sort.SliceStable(ordered, func(i, j int) bool {
				return rateOrder.compare(ordered[i], ordered[j]) < 0
		placeholder)
	placeholder
		if requireCompact {
			ordered = prioritizeOpenAICompactAccounts(ordered)
	placeholder
		for _, acc := range ordered {
			fresh := s.resolveFreshSchedulableOpenAIAccount(ctx, acc, platform, requestedModel, false, requiredCapability)
			if fresh == nil {
				continue
		placeholder
			fresh = s.recheckSelectedOpenAIAccountFromDB(ctx, fresh, groupID, platform, requestedModel, requireCompact, requiredCapability)
			if fresh == nil {
				continue
		placeholder
			if needsUpstreamCheck && s.isUpstreamModelRestrictedByChannel(ctx, *groupID, fresh, requestedModel, requireCompact) {
				continue
		placeholder
			result, err := s.tryAcquireAccountSlot(ctx, fresh.ID, fresh.Concurrency)
			if err == nil && result != nil && result.Acquired {
				selection, selectErr := s.newAcquiredSelectionResult(ctx, fresh, result.ReleaseFunc)
				if selectErr != nil {
					return nil, selectErr
			placeholder
				if sessionHash != "" && !gatewayProfitControlGateActive(ctx) {
					_ = s.setStickySessionAccountID(ctx, groupID, sessionHash, fresh.ID, openaiStickySessionTTL)
			placeholder
				return selection, nil
		placeholder
	placeholder
placeholder else {
		if selection, attempted, selectErr := tryAcquireFromLoadMap(loadMap); selectErr != nil {
			return nil, selectErr
	placeholder else if selection != nil {
			return selection, nil
	placeholder else if attempted {
			if freshLoadMap, loadErr := s.concurrencyService.GetAccountsLoadBatchFresh(ctx, accountLoads); loadErr == nil {
				if selection, _, selectErr := tryAcquireFromLoadMap(freshLoadMap); selectErr != nil {
					return nil, selectErr
			placeholder else if selection != nil {
					return selection, nil
			placeholder
		placeholder
	placeholder
placeholder

	// ============ Layer 3: Fallback wait ============
	sortAccountsByPriorityAndLastUsed(candidates, false)
	if rateOrder.enabled {
		sort.SliceStable(candidates, func(i, j int) bool {
			return rateOrder.compare(candidates[i], candidates[j]) < 0
	placeholder)
placeholder
	if requireCompact {
		candidates = prioritizeOpenAICompactAccounts(candidates)
placeholder
	for _, acc := range candidates {
		fresh := s.resolveFreshSchedulableOpenAIAccount(ctx, acc, platform, requestedModel, false, requiredCapability)
		if fresh == nil {
			continue
	placeholder
		fresh = s.recheckSelectedOpenAIAccountFromDB(ctx, fresh, groupID, platform, requestedModel, requireCompact, requiredCapability)
		if fresh == nil {
			continue
	placeholder
		if needsUpstreamCheck && s.isUpstreamModelRestrictedByChannel(ctx, *groupID, fresh, requestedModel, requireCompact) {
			continue
	placeholder
		return s.newSelectionResult(ctx, fresh, false, nil, &AccountWaitPlan{
			AccountID:      fresh.ID,
			MaxConcurrency: fresh.Concurrency,
			Timeout:        cfg.FallbackWaitTimeout,
			MaxWaiting:     cfg.FallbackMaxWaiting,
	placeholder)
placeholder

	if requireCompact && baseCandidateCount > 0 {
		return nil, ErrNoAvailableCompactAccounts
placeholder
	return nil, ErrNoAvailableAccounts
placeholder

func (s *OpenAIGatewayService) listSchedulableAccounts(ctx context.Context, groupID *int64, platform string) ([]Account, error) {
	platform = NormalizeOpenAICompatiblePlatform(platform)
	if s.schedulerSnapshot != nil {
		accounts, _, err := s.schedulerSnapshot.ListSchedulableAccounts(ctx, groupID, platform, false)
		if err != nil {
			return accounts, err
	placeholder
		accounts = s.filterOpenAIAccountsBySchedulingThreshold(ctx, accounts)
		if platform == PlatformGrok {
			accounts = s.filterGrokFreeQuotaAccountsForOpenAI(ctx, accounts)
	placeholder
		return accounts, nil
placeholder
	var accounts []Account
	var err error
	if s.cfg != nil && s.cfg.RunMode == config.RunModeSimple {
		accounts, err = s.accountRepo.ListSchedulableByPlatform(ctx, platform)
placeholder else if groupID != nil {
		accounts, err = s.accountRepo.ListSchedulableByGroupIDAndPlatform(ctx, *groupID, platform)
placeholder else {
		accounts, err = s.accountRepo.ListSchedulableUngroupedByPlatform(ctx, platform)
placeholder
	if err != nil {
		return nil, fmt.Errorf("query accounts failed: %w", err)
placeholder
	accounts = s.filterOpenAIAccountsBySchedulingThreshold(ctx, accounts)
	if platform == PlatformGrok {
		accounts = s.filterGrokFreeQuotaAccountsForOpenAI(ctx, accounts)
placeholder
	return accounts, nil
placeholder

func (s *OpenAIGatewayService) tryAcquireAccountSlot(ctx context.Context, accountID int64, maxConcurrency int) (*AcquireResult, error) {
	if s.concurrencyService == nil {
		return &AcquireResult{Acquired: true, ReleaseFunc: func() {placeholderplaceholder, nil
placeholder
	return s.concurrencyService.AcquireAccountSlot(ctx, accountID, maxConcurrency)
placeholder

func (s *OpenAIGatewayService) resolveFreshSchedulableOpenAIAccount(ctx context.Context, account *Account, platform string, requestedModel string, requireCompact bool, requiredCapability OpenAIEndpointCapability) *Account {
	fresh := s.resolveFreshSchedulableOpenAIAccountBeforeProfit(ctx, account, platform, requestedModel, requireCompact, requiredCapability)
	if fresh == nil {
		return nil
placeholder
	if vetoed, _ := openAIProfitControlVetoReason(ctx, fresh); vetoed {
		return nil
placeholder
	return fresh
placeholder

func (s *OpenAIGatewayService) resolveFreshSchedulableOpenAIAccountBeforeProfit(ctx context.Context, account *Account, platform string, requestedModel string, requireCompact bool, requiredCapability OpenAIEndpointCapability) *Account {
	if account == nil {
		return nil
placeholder
	platform = NormalizeOpenAICompatiblePlatform(platform)

	fresh := account
	if s.schedulerSnapshot != nil {
		current, err := s.getSchedulableAccount(ctx, account.ID)
		if err != nil || current == nil {
			return nil
	placeholder
		fresh = current
placeholder

	if !isOpenAICompatibleAccountEligibleForRequestBeforeProfit(ctx, fresh, platform, requestedModel, requireCompact, requiredCapability) {
		return nil
placeholder
	if !parentHealthyForShadow(fresh, s.parentAccountLookup(ctx)) {
		return nil
placeholder
	if s.isOpenAIAccountRequestRuntimeBlocked(fresh, requestedModel) {
		return nil
placeholder
	if s.isOpenAIAccountBlockedBySchedulingThreshold(ctx, fresh) {
		return nil
placeholder
	if s.isOpenAIProxyStreamQuarantined(ctx, fresh) {
		return nil
placeholder
	return fresh
placeholder

// parentAccountLookup 返回供 parentHealthyForShadow 使用的母账号解析闭包:经 accountRepo
// 按 ID 取当前 Account(repo 为空时 fail-closed 返回 nil)。统一调度/粘连各路径的母账号解析,
// 取代各调用点重复内联的同一闭包(历史上 recheck 等路径还漏写过 accountRepo==nil 守卫)。
// L2 候选循环改用带 per-pass 缓存的 parentLookupL2,不走此方法。
func (s *OpenAIGatewayService) parentAccountLookup(ctx context.Context) func(int64) *Account {
	return func(id int64) *Account {
		if s.accountRepo == nil {
			return nil
	placeholder
		a, _ := s.accountRepo.GetByID(ctx, id)
		return a
placeholder
placeholder

func (s *OpenAIGatewayService) recheckSelectedOpenAIAccountFromDB(ctx context.Context, account *Account, groupID *int64, platform string, requestedModel string, requireCompact bool, requiredCapability OpenAIEndpointCapability) *Account {
	latest := s.recheckSelectedOpenAIAccountFromDBBeforeProfit(ctx, account, groupID, platform, requestedModel, requireCompact, requiredCapability)
	if latest == nil {
		return nil
placeholder
	if vetoed, _ := openAIProfitControlVetoReason(ctx, latest); vetoed {
		return nil
placeholder
	return latest
placeholder

func (s *OpenAIGatewayService) recheckSelectedOpenAIAccountFromDBBeforeProfit(ctx context.Context, account *Account, groupID *int64, platform string, requestedModel string, requireCompact bool, requiredCapability OpenAIEndpointCapability) *Account {
	if account == nil {
		return nil
placeholder
	platform = NormalizeOpenAICompatiblePlatform(platform)
	if s.schedulerSnapshot == nil || s.accountRepo == nil {
		if s.openAIGroupRequiresPrivacySet(ctx, groupID) && !account.IsPrivacySet() {
			return nil
	placeholder
		if !isOpenAICompatibleAccountEligibleForRequestBeforeProfit(ctx, account, platform, requestedModel, requireCompact, requiredCapability) {
			return nil
	placeholder
		if s.isOpenAIAccountBlockedBySchedulingThreshold(ctx, account) {
			return nil
	placeholder
		if !parentHealthyForShadow(account, s.parentAccountLookup(ctx)) {
			return nil
	placeholder
		if s.isOpenAIProxyStreamQuarantined(ctx, account) {
			return nil
	placeholder
		return account
placeholder

	latest, err := s.accountRepo.GetByID(ctx, account.ID)
	if err != nil || latest == nil {
		return nil
placeholder
	if !s.openAIAccountMatchesSchedulingGroup(latest, groupID) {
		return nil
placeholder
	if s.openAIGroupRequiresPrivacySet(ctx, groupID) && !latest.IsPrivacySet() {
		return nil
placeholder
	if !isOpenAICompatibleAccountEligibleForRequestBeforeProfit(ctx, latest, platform, requestedModel, requireCompact, requiredCapability) {
		return nil
placeholder
	if !parentHealthyForShadow(latest, s.parentAccountLookup(ctx)) {
		return nil
placeholder
	if s.isOpenAIAccountRequestRuntimeBlocked(latest, requestedModel) {
		return nil
placeholder
	if s.isOpenAIAccountBlockedBySchedulingThreshold(ctx, latest) {
		return nil
placeholder
	if s.isOpenAIProxyStreamQuarantined(ctx, latest) {
		return nil
placeholder
	return latest
placeholder

func (s *OpenAIGatewayService) openAIAccountMatchesSchedulingGroup(account *Account, groupID *int64) bool {
	if s != nil && s.cfg != nil && s.cfg.RunMode == config.RunModeSimple {
		return account != nil
placeholder
	return openAIStickyAccountMatchesGroup(account, groupID)
placeholder

func (s *OpenAIGatewayService) getSchedulableAccount(ctx context.Context, accountID int64) (*Account, error) {
	var (
		account *Account
		err     error
	)
	if s.schedulerSnapshot != nil {
		account, err = s.schedulerSnapshot.GetAccount(ctx, accountID)
placeholder else {
		account, err = s.accountRepo.GetByID(ctx, accountID)
placeholder
	if err != nil || account == nil {
		return account, err
placeholder
	if s.isOpenAIAccountBlockedBySchedulingThreshold(ctx, account) {
		return nil, nil
placeholder
	// Legacy sticky (advanced scheduler off) must still free-gate Grok OAuth.
	if account.IsGrok() {
		if gated := s.filterGrokFreeQuotaAccountsForOpenAI(ctx, []Account{*accountplaceholder); len(gated) == 0 {
			return nil, nil
	placeholder
placeholder
	return account, nil
placeholder

// filterGrokFreeQuotaAccountsForOpenAI applies the same local free soft-gate as
// GatewayService / advanced scheduler, for OpenAI-compatible legacy selection.
func (s *OpenAIGatewayService) filterGrokFreeQuotaAccountsForOpenAI(ctx context.Context, accounts []Account) []Account {
	if s == nil {
		return accounts
placeholder
	return filterGrokFreeQuotaAccountsCore(ctx, s.cfg, s.usageLogRepo, &openaiGrokFreeQuotaGateCache, accounts)
placeholder

func (s *OpenAIGatewayService) filterOpenAIAccountsBySchedulingThreshold(ctx context.Context, accounts []Account) []Account {
	if len(accounts) == 0 {
		return accounts
placeholder

	filtered := make([]Account, 0, len(accounts))
	for i := range accounts {
		if s.isOpenAIAccountBlockedBySchedulingThreshold(ctx, &accounts[i]) {
			continue
	placeholder
		filtered = append(filtered, accounts[i])
placeholder
	return filtered
placeholder

func (s *OpenAIGatewayService) isOpenAIAccountBlockedBySchedulingThreshold(ctx context.Context, account *Account) bool {
	if s == nil || s.rateLimitService == nil || account == nil {
		return false
placeholder
	return s.rateLimitService.ApplyAccountSchedulingThreshold(ctx, account)
placeholder

func (s *OpenAIGatewayService) hydrateSelectedAccount(ctx context.Context, account *Account) (*Account, error) {
	if account == nil || s.schedulerSnapshot == nil {
		return account, nil
placeholder
	hydrated, err := s.schedulerSnapshot.GetAccount(ctx, account.ID)
	if err != nil {
		return nil, err
placeholder
	if hydrated == nil {
		return nil, fmt.Errorf("selected openai account %d not found during hydration", account.ID)
placeholder
	return hydrated, nil
placeholder

func (s *OpenAIGatewayService) newSelectionResult(ctx context.Context, account *Account, acquired bool, release func(), waitPlan *AccountWaitPlan) (*AccountSelectionResult, error) {
	hydrated, err := s.hydrateSelectedAccount(ctx, account)
	if err != nil {
		return nil, err
placeholder
	return attachSelectionProfitGate(ctx, &AccountSelectionResult{
		Account:     hydrated,
		Acquired:    acquired,
		ReleaseFunc: release,
		WaitPlan:    waitPlan,
placeholder), nil
placeholder

func (s *OpenAIGatewayService) newAcquiredSelectionResult(ctx context.Context, account *Account, release func()) (*AccountSelectionResult, error) {
	selection, err := s.newSelectionResult(ctx, account, true, release, nil)
	if err != nil && release != nil {
		release()
placeholder
	return selection, err
placeholder

func (s *OpenAIGatewayService) schedulingConfig() config.GatewaySchedulingConfig {
	if s.cfg != nil {
		return s.cfg.Gateway.Scheduling
placeholder
	return config.GatewaySchedulingConfig{
		StickySessionMaxWaiting:  3,
		StickySessionWaitTimeout: 45 * time.Second,
		FallbackWaitTimeout:      30 * time.Second,
		FallbackMaxWaiting:       100,
		LoadBatchEnabled:         true,
		SlotCleanupInterval:      30 * time.Second,
placeholder
placeholder
