package service

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"log/slog"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"golang.org/x/sync/singleflight"
)

const (
	openAIAccountScheduleLayerPreviousResponse = "previous_response_id"
	openAIAccountScheduleLayerSessionSticky    = "session_hash"
	openAIAccountScheduleLayerLoadBalance      = "load_balance"
	openAIAdvancedSchedulerSettingKey          = "openai_advanced_scheduler_enabled"
)

const (
	openAIAdvancedSchedulerSettingCacheTTL  = 5 * time.Second
	openAIAdvancedSchedulerSettingDBTimeout = 2 * time.Second
	// ponytail: cap probes added when cost ordering expands configured Top-K;
	// use bulk acquisition if a measured workload needs a higher ceiling.
	openAIAccountSelectionProbeLimit = 64
)

const (
	openAIQuotaHeadroomNeutralFactor           = 0.5
	openAIQuotaHeadroomSecondaryLowRemain      = 0.10
	openAIQuotaHeadroomSnapshotStaleAfter      = 8 * time.Hour
	openAIUpstreamCostNeutralFactor            = 0.5
	defaultOpenAIOAuthSchedulingRateMultiplier = 1.0
)

type cachedOpenAIAdvancedSchedulerSetting struct {
	lowUpstreamRatePriorityEnabled bool
	oauthSchedulingRateMultiplier  float64
	enabled                        bool
	stickyWeightedEnabled          bool
	subscriptionPriorityEnabled    bool
	lbTopKOverride                 int
	weightOverrides                map[string]float64
	expiresAt                      int64
placeholder

type openAIAdvancedSchedulerRuntimeSettings struct {
	lowUpstreamRatePriorityEnabled bool
	oauthSchedulingRateMultiplier  float64
	enabled                        bool
	stickyWeightedEnabled          bool
	subscriptionPriorityEnabled    bool
	lbTopKOverride                 int
	weightOverrides                map[string]float64
placeholder

var openAIAdvancedSchedulerSettingCache atomic.Value // *cachedOpenAIAdvancedSchedulerSetting
var openAIAdvancedSchedulerSettingSF singleflight.Group

type OpenAIAccountScheduleRequest struct {
	GroupID                 *int64
	Platform                string
	SessionHash             string
	StickyAccountID         int64
	StickyPreviousAccountID int64
	StickyWeighted          bool
	SubscriptionPriority    bool
	PreserveStickyBinding   bool
	PreviousResponseID      string
	PreviousResponseCanMove bool
	UseUpstreamTokenCost    bool
	RequestedModel          string
	RequiredTransport       OpenAIUpstreamTransport
	RequiredCapability      OpenAIEndpointCapability
	RequiredImageCapability OpenAIImagesCapability
	RequireCompact          bool
	ExcludedIDs             map[int64]struct{placeholder
placeholder

type OpenAIAccountScheduleDecision struct {
	Layer               string
	StickyPreviousHit   bool
	StickySessionHit    bool
	CandidateCount      int
	TopK                int
	LatencyMs           int64
	LoadSkew            float64
	SelectedAccountID   int64
	SelectedAccountType string
placeholder

type OpenAIAccountSchedulerMetricsSnapshot struct {
	SelectTotal              int64
	StickyPreviousHitTotal   int64
	StickySessionHitTotal    int64
	LoadBalanceSelectTotal   int64
	AccountSwitchTotal       int64
	SchedulerLatencyMsTotal  int64
	SchedulerLatencyMsAvg    float64
	StickyHitRatio           float64
	AccountSwitchRate        float64
	LoadSkewAvg              float64
	RuntimeStatsAccountCount int
placeholder

type OpenAIAccountScheduler interface {
	Select(ctx context.Context, req OpenAIAccountScheduleRequest) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error)
	ReportResult(accountID int64, success bool, firstTokenMs *int)
	ReportSwitch()
	SnapshotMetrics() OpenAIAccountSchedulerMetricsSnapshot
placeholder

type openAIAccountSchedulerMetrics struct {
	selectTotal            atomic.Int64
	stickyPreviousHitTotal atomic.Int64
	stickySessionHitTotal  atomic.Int64
	loadBalanceSelectTotal atomic.Int64
	accountSwitchTotal     atomic.Int64
	latencyMsTotal         atomic.Int64
	loadSkewMilliTotal     atomic.Int64
placeholder

type openAIAccountLoadPlan struct {
	allCandidates             []openAIAccountCandidateScore
	candidates                []openAIAccountCandidateScore
	staleSnapshotCompactRetry []openAIAccountCandidateScore
	selectionOrder            []openAIAccountCandidateScore
	candidateCount            int
	topK                      int
	loadSkew                  float64
	includeOverflowFallback   bool
placeholder

type openAIAccountLoadSelectionAttempt struct {
	result              *AccountSelectionResult
	selectionOrder      []openAIAccountCandidateScore
	candidateCount      int
	topK                int
	loadSkew            float64
	compactBlocked      bool
	noCompactCandidates bool
	err                 error
placeholder

func (m *openAIAccountSchedulerMetrics) recordSelect(decision OpenAIAccountScheduleDecision) {
	if m == nil {
		return
placeholder
	m.selectTotal.Add(1)
	m.latencyMsTotal.Add(decision.LatencyMs)
	m.loadSkewMilliTotal.Add(int64(math.Round(decision.LoadSkew * 1000)))
	if decision.StickyPreviousHit {
		m.stickyPreviousHitTotal.Add(1)
placeholder
	if decision.StickySessionHit {
		m.stickySessionHitTotal.Add(1)
placeholder
	if decision.Layer == openAIAccountScheduleLayerLoadBalance {
		m.loadBalanceSelectTotal.Add(1)
placeholder
placeholder

func (m *openAIAccountSchedulerMetrics) recordSwitch() {
	if m == nil {
		return
placeholder
	m.accountSwitchTotal.Add(1)
placeholder

type openAIAccountRuntimeStats struct {
	accounts     sync.Map
	accountCount atomic.Int64
placeholder

type openAIAccountRuntimeStat struct {
	errorRateEWMABits atomic.Uint64
	ttftEWMABits      atomic.Uint64
placeholder

func newOpenAIAccountRuntimeStats() *openAIAccountRuntimeStats {
	return &openAIAccountRuntimeStats{placeholder
placeholder

func (s *openAIAccountRuntimeStats) loadOrCreate(accountID int64) *openAIAccountRuntimeStat {
	if value, ok := s.accounts.Load(accountID); ok {
		stat, _ := value.(*openAIAccountRuntimeStat)
		if stat != nil {
			return stat
	placeholder
placeholder

	stat := &openAIAccountRuntimeStat{placeholder
	stat.ttftEWMABits.Store(math.Float64bits(math.NaN()))
	actual, loaded := s.accounts.LoadOrStore(accountID, stat)
	if !loaded {
		s.accountCount.Add(1)
		return stat
placeholder
	existing, _ := actual.(*openAIAccountRuntimeStat)
	if existing != nil {
		return existing
placeholder
	return stat
placeholder

func updateEWMAAtomic(target *atomic.Uint64, sample float64, alpha float64) {
	for {
		oldBits := target.Load()
		oldValue := math.Float64frombits(oldBits)
		newValue := alpha*sample + (1-alpha)*oldValue
		if target.CompareAndSwap(oldBits, math.Float64bits(newValue)) {
			return
	placeholder
placeholder
placeholder

func (s *openAIAccountRuntimeStats) report(accountID int64, success bool, firstTokenMs *int) {
	if s == nil || accountID <= 0 {
		return
placeholder
	const alpha = 0.2
	stat := s.loadOrCreate(accountID)

	errorSample := 1.0
	if success {
		errorSample = 0.0
placeholder
	updateEWMAAtomic(&stat.errorRateEWMABits, errorSample, alpha)

	if firstTokenMs != nil && *firstTokenMs > 0 {
		ttft := float64(*firstTokenMs)
		ttftBits := math.Float64bits(ttft)
		for {
			oldBits := stat.ttftEWMABits.Load()
			oldValue := math.Float64frombits(oldBits)
			if math.IsNaN(oldValue) {
				if stat.ttftEWMABits.CompareAndSwap(oldBits, ttftBits) {
					break
			placeholder
				continue
		placeholder
			newValue := alpha*ttft + (1-alpha)*oldValue
			if stat.ttftEWMABits.CompareAndSwap(oldBits, math.Float64bits(newValue)) {
				break
		placeholder
	placeholder
placeholder
placeholder

func (s *openAIAccountRuntimeStats) snapshot(accountID int64) (errorRate float64, ttft float64, hasTTFT bool) {
	if s == nil || accountID <= 0 {
		return 0, 0, false
placeholder
	value, ok := s.accounts.Load(accountID)
	if !ok {
		return 0, 0, false
placeholder
	stat, _ := value.(*openAIAccountRuntimeStat)
	if stat == nil {
		return 0, 0, false
placeholder
	errorRate = clamp01(math.Float64frombits(stat.errorRateEWMABits.Load()))
	ttftValue := math.Float64frombits(stat.ttftEWMABits.Load())
	if math.IsNaN(ttftValue) {
		return errorRate, 0, false
placeholder
	return errorRate, ttftValue, true
placeholder

func (s *openAIAccountRuntimeStats) size() int {
	if s == nil {
		return 0
placeholder
	return int(s.accountCount.Load())
placeholder

type defaultOpenAIAccountScheduler struct {
	service                *OpenAIGatewayService
	metrics                openAIAccountSchedulerMetrics
	stats                  *openAIAccountRuntimeStats
	grokFreeQuotaGateCache sync.Map // key: int64(accountID), value: grokFreeQuotaGateCacheEntry
placeholder

type openAISelectionProbeBudget struct {
	acquires  int
	rechecks  int
	attempted map[int64]struct{placeholder
	limited   bool
placeholder

func newOpenAISelectionProbeBudget() *openAISelectionProbeBudget {
	return &openAISelectionProbeBudget{attempted: make(map[int64]struct{placeholder)placeholder
placeholder

func (b *openAISelectionProbeBudget) enableLimit() {
	if b != nil {
		b.limited = true
placeholder
placeholder

func (b *openAISelectionProbeBudget) recordAcquire(accountID int64) bool {
	if b == nil {
		return false
placeholder
	if !b.limited {
		return true
placeholder
	if b.acquires >= openAIAccountSelectionProbeLimit {
		return false
placeholder
	if b.attempted == nil {
		b.attempted = make(map[int64]struct{placeholder)
placeholder
	b.acquires++
	b.attempted[accountID] = struct{placeholder{placeholder
	return true
placeholder

func (b *openAISelectionProbeBudget) recordRecheck() bool {
	if b == nil {
		return false
placeholder
	if !b.limited {
		return true
placeholder
	if b.rechecks >= openAIAccountSelectionProbeLimit {
		return false
placeholder
	b.rechecks++
	return true
placeholder

func (b *openAISelectionProbeBudget) acquireExhausted() bool {
	return b != nil && b.limited && b.acquires >= openAIAccountSelectionProbeLimit
placeholder

func (b *openAISelectionProbeBudget) wasAttempted(accountID int64) bool {
	if b == nil {
		return false
placeholder
	_, ok := b.attempted[accountID]
	return ok
placeholder

type openAIStickyEscapeConfig struct {
	enabled   bool
	ttftMs    float64
	errorRate float64
placeholder

func newDefaultOpenAIAccountScheduler(service *OpenAIGatewayService, stats *openAIAccountRuntimeStats) OpenAIAccountScheduler {
	if stats == nil {
		stats = newOpenAIAccountRuntimeStats()
placeholder
	return &defaultOpenAIAccountScheduler{
		service: service,
		stats:   stats,
placeholder
placeholder

func (s *defaultOpenAIAccountScheduler) Select(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error) {
	decision := OpenAIAccountScheduleDecision{placeholder
	start := time.Now()
	defer func() {
		decision.LatencyMs = time.Since(start).Milliseconds()
		s.metrics.recordSelect(decision)
placeholder()

	previousResponseID := strings.TrimSpace(req.PreviousResponseID)
	if previousResponseID != "" && normalizeOpenAICompatiblePlatform(req.Platform) == PlatformOpenAI &&
		(!req.StickyWeighted || !req.PreviousResponseCanMove) {
		selection, err := s.service.selectAccountByPreviousResponseIDForCapability(
			ctx,
			req.GroupID,
			previousResponseID,
			req.RequestedModel,
			req.ExcludedIDs,
			req.RequiredCapability,
			req.RequireCompact,
		)
		if err != nil {
			return nil, decision, err
	placeholder
		if selection != nil && selection.Account != nil {
			if !s.isAccountTransportCompatible(selection.Account, req.RequiredTransport) {
				if selection.ReleaseFunc != nil {
					selection.ReleaseFunc()
			placeholder
				selection = nil
		placeholder
	placeholder
		if selection != nil && selection.Account != nil {
			decision.Layer = openAIAccountScheduleLayerPreviousResponse
			decision.StickyPreviousHit = true
			decision.SelectedAccountID = selection.Account.ID
			decision.SelectedAccountType = selection.Account.Type
			if req.SessionHash != "" {
				_ = s.service.bindOpenAIStickySessionDuringSelection(ctx, req.GroupID, req.SessionHash, selection.Account.ID)
		placeholder
			return selection, decision, nil
	placeholder
placeholder

	if !req.StickyWeighted {
		selection, escapedSticky, err := s.selectBySessionHash(ctx, req)
		if err != nil {
			return nil, decision, err
	placeholder
		if selection != nil && selection.Account != nil {
			decision.Layer = openAIAccountScheduleLayerSessionSticky
			decision.StickySessionHit = true
			decision.SelectedAccountID = selection.Account.ID
			decision.SelectedAccountType = selection.Account.Type
			return selection, decision, nil
	placeholder
		if escapedSticky {
			req.PreserveStickyBinding = true
	placeholder
placeholder

	selection, candidateCount, topK, loadSkew, err := s.selectByLoadBalance(ctx, req)
	decision.Layer = openAIAccountScheduleLayerLoadBalance
	decision.CandidateCount = candidateCount
	decision.TopK = topK
	decision.LoadSkew = loadSkew
	if err != nil {
		return nil, decision, err
placeholder
	if selection != nil && selection.Account != nil {
		decision.SelectedAccountID = selection.Account.ID
		decision.SelectedAccountType = selection.Account.Type
		if req.StickyWeighted {
			if req.StickyPreviousAccountID > 0 && selection.Account.ID == req.StickyPreviousAccountID {
				decision.StickyPreviousHit = true
		placeholder
			if req.StickyAccountID > 0 && selection.Account.ID == req.StickyAccountID {
				decision.StickySessionHit = true
		placeholder
	placeholder
placeholder
	return selection, decision, nil
placeholder

func (s *defaultOpenAIAccountScheduler) selectBySessionHash(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
) (*AccountSelectionResult, bool, error) {
	sessionHash := strings.TrimSpace(req.SessionHash)
	if sessionHash == "" || s == nil || s.service == nil || s.service.cache == nil {
		return nil, false, nil
placeholder

	accountID := req.StickyAccountID
	if accountID <= 0 {
		var err error
		accountID, err = s.service.getStickySessionAccountID(ctx, req.GroupID, sessionHash)
		if err != nil || accountID <= 0 {
			return nil, false, nil
	placeholder
placeholder
	if accountID <= 0 {
		return nil, false, nil
placeholder
	if req.ExcludedIDs != nil {
		if _, excluded := req.ExcludedIDs[accountID]; excluded {
			return nil, false, nil
	placeholder
placeholder

	account, err := s.service.getSchedulableAccount(ctx, accountID)
	if err != nil || account == nil {
		_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, sessionHash)
		return nil, false, nil
placeholder
	if shouldClearStickySession(account, req.RequestedModel) || account.Platform != normalizeOpenAICompatiblePlatform(req.Platform) || !account.IsOpenAICompatible() || !account.IsSchedulable() {
		_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, sessionHash)
		return nil, false, nil
placeholder
	if !s.isAccountRequestCompatible(ctx, account, req) {
		return nil, false, nil
placeholder
	if !s.isAccountTransportCompatible(account, req.RequiredTransport) {
		_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, sessionHash)
		return nil, false, nil
placeholder
	account = s.service.recheckSelectedOpenAIAccountFromDB(ctx, account, req.GroupID, req.Platform, req.RequestedModel, req.RequireCompact, req.RequiredCapability)
	if account == nil || !s.service.openAIAccountMatchesSchedulingGroup(account, req.GroupID) || !s.isAccountTransportCompatible(account, req.RequiredTransport) {
		_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, sessionHash)
		return nil, false, nil
placeholder
	// Free-tier soft gate: sticky session must not pin an over-quota free OAuth account.
	// Admin QueryQuota / import probes do not use this path.
	if account != nil && len(s.filterGrokFreeQuotaAccounts(ctx, []Account{*accountplaceholder)) == 0 {
		_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, sessionHash)
		return nil, false, nil
placeholder
	escapeCfg := s.service.openAIStickyEscapeConfig()
	if reason, errorRate, ttft, shouldEscape := s.shouldEscapeStickyAccount(accountID, escapeCfg); shouldEscape {
		slog.Info("sticky_escape_triggered",
			"account_id", accountID,
			"reason", reason,
			"error_rate", errorRate,
			"ttft", ttft,
		)
		return nil, true, nil
placeholder
	result, acquireErr := s.service.tryAcquireAccountSlot(ctx, accountID, account.Concurrency)
	if acquireErr == nil && result != nil && result.Acquired {
		_ = s.service.refreshStickySessionTTL(ctx, req.GroupID, sessionHash, s.service.openAIWSSessionStickyTTL())
		return attachSelectionProfitGate(ctx, &AccountSelectionResult{
			Account:     account,
			Acquired:    true,
			ReleaseFunc: result.ReleaseFunc,
	placeholder), false, nil
placeholder

	cfg := s.service.schedulingConfig()
	// WaitPlan.MaxConcurrency 使用 Concurrency（非 EffectiveLoadFactor），因为 WaitPlan 控制的是 Redis 实际并发槽位等待。
	if s.service.concurrencyService != nil {
		if escapeCfg.enabled && acquireErr == nil && result != nil && !result.Acquired {
			errorRate, ttft, _ := s.stats.snapshot(accountID)
			slog.Info("sticky_escape_triggered",
				"account_id", accountID,
				"reason", "concurrency_full",
				"error_rate", errorRate,
				"ttft", ttft,
			)
			return nil, true, nil
	placeholder
		return attachSelectionProfitGate(ctx, &AccountSelectionResult{
			Account: account,
			WaitPlan: &AccountWaitPlan{
				AccountID:      accountID,
				MaxConcurrency: account.Concurrency,
				Timeout:        cfg.StickySessionWaitTimeout,
				MaxWaiting:     cfg.StickySessionMaxWaiting,
		placeholder,
	placeholder), false, nil
placeholder
	return nil, false, nil
placeholder

func openAIStickyAccountMatchesGroup(account *Account, groupID *int64) bool {
	if account == nil {
		return false
placeholder
	if groupID == nil {
		return len(account.AccountGroups) == 0 && len(account.GroupIDs) == 0
placeholder
	for _, accountGroupID := range account.GroupIDs {
		if accountGroupID == *groupID {
			return true
	placeholder
placeholder
	for _, accountGroup := range account.AccountGroups {
		if accountGroup.GroupID == *groupID {
			return true
	placeholder
placeholder
	return false
placeholder

func openAIAccountSchedulingPriority(account *Account) int {
	if account == nil {
		return 0
placeholder
	return account.Priority
placeholder

func (s *defaultOpenAIAccountScheduler) shouldEscapeStickyAccount(accountID int64, cfg openAIStickyEscapeConfig) (reason string, errorRate float64, ttft float64, shouldEscape bool) {
	if !cfg.enabled || s == nil || s.stats == nil || accountID <= 0 {
		return "", 0, 0, false
placeholder
	errorRate, ttft, hasTTFT := s.stats.snapshot(accountID)
	if hasTTFT && ttft > cfg.ttftMs {
		return "ttft", errorRate, ttft, true
placeholder
	if errorRate > cfg.errorRate {
		return "error_rate", errorRate, ttft, true
placeholder
	return "", errorRate, ttft, false
placeholder

type openAIAccountCandidateScore struct {
	account   *Account
	loadInfo  *AccountLoadInfo
	loadKnown bool
	score     float64
	priority  int
	errorRate float64
	ttft      float64
	hasTTFT   bool
placeholder

type openAIAccountCandidateHeap []openAIAccountCandidateScore

func (h openAIAccountCandidateHeap) Len() int {
	return len(h)
placeholder

func (h openAIAccountCandidateHeap) Less(i, j int) bool {
	// 最小堆根节点保存“最差”候选，便于 O(log k) 维护 topK。
	return isOpenAIAccountCandidateBetter(h[j], h[i])
placeholder

func (h openAIAccountCandidateHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
placeholder

func (h *openAIAccountCandidateHeap) Push(x any) {
	candidate, ok := x.(openAIAccountCandidateScore)
	if !ok {
		panic("openAIAccountCandidateHeap: invalid element type")
placeholder
	*h = append(*h, candidate)
placeholder

func (h *openAIAccountCandidateHeap) Pop() any {
	old := *h
	n := len(old)
	last := old[n-1]
	*h = old[:n-1]
	return last
placeholder

func isOpenAIAccountCandidateBetter(left openAIAccountCandidateScore, right openAIAccountCandidateScore) bool {
	if left.score != right.score {
		return left.score > right.score
placeholder
	if left.account.Priority != right.account.Priority {
		return left.account.Priority < right.account.Priority
placeholder
	if left.loadInfo.LoadRate != right.loadInfo.LoadRate {
		return left.loadInfo.LoadRate < right.loadInfo.LoadRate
placeholder
	if left.loadInfo.WaitingCount != right.loadInfo.WaitingCount {
		return left.loadInfo.WaitingCount < right.loadInfo.WaitingCount
placeholder
	return left.account.ID < right.account.ID
placeholder

func selectTopKOpenAICandidates(candidates []openAIAccountCandidateScore, topK int) []openAIAccountCandidateScore {
	if len(candidates) == 0 {
		return nil
placeholder
	if topK <= 0 {
		topK = 1
placeholder
	if topK >= len(candidates) {
		ranked := append([]openAIAccountCandidateScore(nil), candidates...)
		sort.Slice(ranked, func(i, j int) bool {
			return isOpenAIAccountCandidateBetter(ranked[i], ranked[j])
	placeholder)
		return ranked
placeholder

	best := make(openAIAccountCandidateHeap, 0, topK)
	for _, candidate := range candidates {
		if len(best) < topK {
			heap.Push(&best, candidate)
			continue
	placeholder
		if isOpenAIAccountCandidateBetter(candidate, best[0]) {
			best[0] = candidate
			heap.Fix(&best, 0)
	placeholder
placeholder

	ranked := make([]openAIAccountCandidateScore, len(best))
	copy(ranked, best)
	sort.Slice(ranked, func(i, j int) bool {
		return isOpenAIAccountCandidateBetter(ranked[i], ranked[j])
placeholder)
	return ranked
placeholder

type openAISelectionRNG struct {
	state uint64
placeholder

func newOpenAISelectionRNG(seed uint64) openAISelectionRNG {
	if seed == 0 {
		seed = 0x9e3779b97f4a7c15
placeholder
	return openAISelectionRNG{state: seedplaceholder
placeholder

func (r *openAISelectionRNG) nextUint64() uint64 {
	// xorshift64*
	x := r.state
	x ^= x >> 12
	x ^= x << 25
	x ^= x >> 27
	r.state = x
	return x * 2685821657736338717
placeholder

func (r *openAISelectionRNG) nextFloat64() float64 {
	// [0,1)
	return float64(r.nextUint64()>>11) / (1 << 53)
placeholder

func deriveOpenAISelectionSeed(req OpenAIAccountScheduleRequest) uint64 {
	hasher := fnv.New64a()
	writeValue := func(value string) {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return
	placeholder
		_, _ = hasher.Write([]byte(trimmed))
		_, _ = hasher.Write([]byte{0placeholder)
placeholder

	writeValue(req.SessionHash)
	writeValue(req.PreviousResponseID)
	writeValue(req.RequestedModel)
	if req.GroupID != nil {
		_, _ = hasher.Write([]byte(strconv.FormatInt(*req.GroupID, 10)))
placeholder

	seed := hasher.Sum64()
	// 对“无会话锚点”的纯负载均衡请求引入时间熵，避免固定命中同一账号。
	if strings.TrimSpace(req.SessionHash) == "" && strings.TrimSpace(req.PreviousResponseID) == "" {
		seed ^= uint64(time.Now().UnixNano())
placeholder
	if seed == 0 {
		seed = uint64(time.Now().UnixNano()) ^ 0x9e3779b97f4a7c15
placeholder
	return seed
placeholder

func buildOpenAIWeightedSelectionOrder(
	candidates []openAIAccountCandidateScore,
	req OpenAIAccountScheduleRequest,
) []openAIAccountCandidateScore {
	if len(candidates) <= 1 {
		return append([]openAIAccountCandidateScore(nil), candidates...)
placeholder

	pool := append([]openAIAccountCandidateScore(nil), candidates...)
	weights := make([]float64, len(pool))
	minScore := pool[0].score
	for i := 1; i < len(pool); i++ {
		if pool[i].score < minScore {
			minScore = pool[i].score
	placeholder
placeholder
	for i := range pool {
		// 将 top-K 分值平移到正区间，避免“单一最高分账号”长期垄断。
		weight := (pool[i].score - minScore) + 1.0
		if math.IsNaN(weight) || math.IsInf(weight, 0) || weight <= 0 {
			weight = 1.0
	placeholder
		weights[i] = weight
placeholder

	order := make([]openAIAccountCandidateScore, 0, len(pool))
	rng := newOpenAISelectionRNG(deriveOpenAISelectionSeed(req))
	for len(pool) > 0 {
		total := 0.0
		for _, w := range weights {
			total += w
	placeholder

		selectedIdx := 0
		if total > 0 {
			r := rng.nextFloat64() * total
			acc := 0.0
			for i, w := range weights {
				acc += w
				if r <= acc {
					selectedIdx = i
					break
			placeholder
		placeholder
	placeholder else {
			selectedIdx = int(rng.nextUint64() % uint64(len(pool)))
	placeholder

		order = append(order, pool[selectedIdx])
		pool = append(pool[:selectedIdx], pool[selectedIdx+1:]...)
		weights = append(weights[:selectedIdx], weights[selectedIdx+1:]...)
placeholder
	return order
placeholder

func (s *defaultOpenAIAccountScheduler) buildOpenAIAccountLoadPlan(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
	filtered []*Account,
	loadMap map[int64]*AccountLoadInfo,
) openAIAccountLoadPlan {
	allCandidates := make([]openAIAccountCandidateScore, 0, len(filtered))
	for _, account := range filtered {
		loadInfo, loadKnown := loadMap[account.ID]
		if !loadKnown || loadInfo == nil {
			loadInfo = &AccountLoadInfo{AccountID: account.IDplaceholder
			loadKnown = false
	placeholder
		errorRate, ttft, hasTTFT := 0.0, 0.0, false
		if s.stats != nil {
			errorRate, ttft, hasTTFT = s.stats.snapshot(account.ID)
	placeholder
		allCandidates = append(allCandidates, openAIAccountCandidateScore{
			account:   account,
			loadInfo:  loadInfo,
			loadKnown: loadKnown,
			errorRate: errorRate,
			ttft:      ttft,
			hasTTFT:   hasTTFT,
	placeholder)
placeholder

	candidates := allCandidates
	staleSnapshotCompactRetry := make([]openAIAccountCandidateScore, 0, len(allCandidates))
	if req.RequireCompact {
		candidates = make([]openAIAccountCandidateScore, 0, len(allCandidates))
		for _, candidate := range allCandidates {
			if openAICompactSupportTier(candidate.account) == 0 {
				staleSnapshotCompactRetry = append(staleSnapshotCompactRetry, candidate)
				continue
		placeholder
			candidates = append(candidates, candidate)
	placeholder
placeholder

	plan := openAIAccountLoadPlan{
		allCandidates:             allCandidates,
		candidates:                candidates,
		staleSnapshotCompactRetry: staleSnapshotCompactRetry,
		candidateCount:            len(candidates),
placeholder
	if len(candidates) == 0 {
		plan.selectionOrder = s.buildOpenAISelectionOrder(req, plan)
		return plan
placeholder

	minPriority, maxPriority := openAIAccountSchedulingPriority(candidates[0].account), openAIAccountSchedulingPriority(candidates[0].account)
	maxWaiting := 1
	loadRateSum := 0.0
	loadRateSumSquares := 0.0
	minTTFT, maxTTFT := 0.0, 0.0
	hasTTFTSample := false
	for i := range candidates {
		candidate := &candidates[i]
		candidate.priority = openAIAccountSchedulingPriority(candidate.account)
		if candidate.priority < minPriority {
			minPriority = candidate.priority
	placeholder
		if candidate.priority > maxPriority {
			maxPriority = candidate.priority
	placeholder
		if candidate.loadInfo.WaitingCount > maxWaiting {
			maxWaiting = candidate.loadInfo.WaitingCount
	placeholder
		if candidate.hasTTFT && candidate.ttft > 0 {
			if !hasTTFTSample {
				minTTFT, maxTTFT = candidate.ttft, candidate.ttft
				hasTTFTSample = true
		placeholder else {
				if candidate.ttft < minTTFT {
					minTTFT = candidate.ttft
			placeholder
				if candidate.ttft > maxTTFT {
					maxTTFT = candidate.ttft
			placeholder
		placeholder
	placeholder
		loadRate := float64(candidate.loadInfo.LoadRate)
		loadRateSum += loadRate
		loadRateSumSquares += loadRate * loadRate
placeholder
	plan.loadSkew = calcLoadSkewByMoments(loadRateSum, loadRateSumSquares, len(candidates))

	weights := s.service.openAIWSSchedulerWeightsForRequest(ctx)
	now := time.Now()
	upstreamCostFactors := map[int64]float64(nil)
	if req.UseUpstreamTokenCost && weights.UpstreamCost > 0 {
		accounts := make([]*Account, 0, len(candidates))
		for _, candidate := range candidates {
			accounts = append(accounts, candidate.account)
	placeholder
		upstreamCostFactors = openAIUpstreamCostFactors(accounts, now, s.service.openAIOAuthSchedulingRateMultiplier(ctx))
		for _, factor := range upstreamCostFactors {
			if factor != openAIUpstreamCostNeutralFactor {
				plan.includeOverflowFallback = true
				break
		placeholder
	placeholder
placeholder

	// Reset 因子（use-it-or-lose-it）：在拥有「未来会话窗口结束时间」的账号中，
	// 剩余时间越短 → 因子越接近 1（越早重置越优先用尽）。无活跃窗口的账号因子为 0。
	// 仅在 weights.Reset > 0 时计算，默认关闭不影响原有行为。
	minResetRemaining, maxResetRemaining := 0.0, 0.0
	hasResetSample := false
	if weights.Reset > 0 {
		for _, candidate := range candidates {
			end := candidate.account.SessionWindowEnd
			if end == nil || !now.Before(*end) {
				continue
		placeholder
			remaining := end.Sub(now).Seconds()
			if !hasResetSample {
				minResetRemaining, maxResetRemaining = remaining, remaining
				hasResetSample = true
				continue
		placeholder
			if remaining < minResetRemaining {
				minResetRemaining = remaining
		placeholder
			if remaining > maxResetRemaining {
				maxResetRemaining = remaining
		placeholder
	placeholder
placeholder

	for i := range candidates {
		item := &candidates[i]
		priorityFactor := 1.0
		if maxPriority > minPriority {
			priorityFactor = 1 - float64(item.priority-minPriority)/float64(maxPriority-minPriority)
	placeholder
		loadFactor := 1 - clamp01(float64(item.loadInfo.LoadRate)/100.0)
		queueFactor := 1 - clamp01(float64(item.loadInfo.WaitingCount)/float64(maxWaiting))
		errorFactor := 1 - clamp01(item.errorRate)
		ttftFactor := 0.5
		if item.hasTTFT && hasTTFTSample && maxTTFT > minTTFT {
			ttftFactor = 1 - clamp01((item.ttft-minTTFT)/(maxTTFT-minTTFT))
	placeholder
		resetFactor := 0.0
		if weights.Reset > 0 && hasResetSample {
			if end := item.account.SessionWindowEnd; end != nil && now.Before(*end) {
				if maxResetRemaining > minResetRemaining {
					resetFactor = 1 - clamp01((end.Sub(now).Seconds()-minResetRemaining)/(maxResetRemaining-minResetRemaining))
			placeholder else {
					// 所有有窗口的账号剩余时间相同：一律给满分，让其优于无窗口账号。
					resetFactor = 1
			placeholder
		placeholder
	placeholder
		quotaHeadroomFactor := 0.0
		if weights.QuotaHeadroom > 0 {
			quotaHeadroomFactor = openAIQuotaHeadroomFactor(item.account, now)
	placeholder
		upstreamCostFactor := openAIUpstreamCostNeutralFactor
		if factor, ok := upstreamCostFactors[item.account.ID]; ok {
			upstreamCostFactor = factor
	placeholder

		item.score = weights.Priority*priorityFactor +
			weights.Load*loadFactor +
			weights.Queue*queueFactor +
			weights.ErrorRate*errorFactor +
			weights.TTFT*ttftFactor +
			weights.Reset*resetFactor +
			weights.QuotaHeadroom*quotaHeadroomFactor +
			weights.UpstreamCost*(upstreamCostFactor-openAIUpstreamCostNeutralFactor)
		if req.StickyWeighted {
			if req.PreviousResponseCanMove && req.StickyPreviousAccountID > 0 && item.account.ID == req.StickyPreviousAccountID {
				item.score += weights.Previous
		placeholder
			if req.StickyAccountID > 0 && item.account.ID == req.StickyAccountID {
				item.score += weights.SessionSticky
		placeholder
	placeholder
placeholder
	plan.candidates = candidates

	plan.topK = s.service.openAIWSLBTopKForRequest(ctx)
	if plan.topK > len(candidates) {
		plan.topK = len(candidates)
placeholder
	if plan.topK <= 0 {
		plan.topK = 1
placeholder

	plan.selectionOrder = s.buildOpenAISelectionOrder(req, plan)
	return plan
placeholder

func (s *defaultOpenAIAccountScheduler) buildOpenAISelectionOrder(
	req OpenAIAccountScheduleRequest,
	plan openAIAccountLoadPlan,
) []openAIAccountCandidateScore {
	buildSelectionOrder := func(pool []openAIAccountCandidateScore) []openAIAccountCandidateScore {
		if len(pool) == 0 || plan.topK <= 0 {
			return nil
	placeholder
		groupTopK := plan.topK
		if groupTopK > len(pool) {
			groupTopK = len(pool)
	placeholder
		ranked := selectTopKOpenAICandidates(pool, groupTopK)
		var primary []openAIAccountCandidateScore
		if req.StickyWeighted {
			for _, stickyID := range []int64{req.StickyPreviousAccountID, req.StickyAccountIDplaceholder {
				if stickyID <= 0 {
					continue
			placeholder
				for i, candidate := range ranked {
					if candidate.account != nil && candidate.account.ID == stickyID {
						primary = append([]openAIAccountCandidateScore{candidateplaceholder, ranked[:i]...)
						primary = append(primary, ranked[i+1:]...)
						break
				placeholder
			placeholder
				if len(primary) > 0 {
					break
			placeholder
		placeholder
	placeholder
		if len(primary) == 0 {
			primary = buildOpenAIWeightedSelectionOrder(ranked, req)
	placeholder
		if !plan.includeOverflowFallback || groupTopK >= len(pool) {
			return primary
	placeholder

		selected := make(map[int64]struct{placeholder, len(primary))
		for _, candidate := range primary {
			selected[candidate.account.ID] = struct{placeholder{placeholder
	placeholder
		overflow := make([]openAIAccountCandidateScore, 0, len(pool)-len(primary))
		for _, candidate := range pool {
			if _, ok := selected[candidate.account.ID]; !ok {
				overflow = append(overflow, candidate)
		placeholder
	placeholder
		sort.Slice(overflow, func(i, j int) bool {
			return isOpenAIAccountCandidateBetter(overflow[i], overflow[j])
	placeholder)
		return append(primary, overflow...)
placeholder

	if req.RequireCompact {
		supported := make([]openAIAccountCandidateScore, 0, len(plan.candidates))
		unknown := make([]openAIAccountCandidateScore, 0, len(plan.candidates))
		for _, candidate := range plan.candidates {
			switch openAICompactSupportTier(candidate.account) {
			case 2:
				supported = append(supported, candidate)
			case 1:
				unknown = append(unknown, candidate)
		placeholder
	placeholder
		selectionOrder := make([]openAIAccountCandidateScore, 0, len(plan.allCandidates))
		selectionOrder = append(selectionOrder, buildSelectionOrder(supported)...)
		selectionOrder = append(selectionOrder, buildSelectionOrder(unknown)...)
		if len(plan.staleSnapshotCompactRetry) > 0 && s.service.schedulerSnapshot != nil {
			selectionOrder = append(selectionOrder, sortOpenAICompactRetryCandidates(plan.staleSnapshotCompactRetry)...)
	placeholder
		return selectionOrder
placeholder

	return buildSelectionOrder(plan.candidates)
placeholder

func sortOpenAICompactRetryCandidates(pool []openAIAccountCandidateScore) []openAIAccountCandidateScore {
	if len(pool) == 0 {
		return nil
placeholder
	ordered := append([]openAIAccountCandidateScore(nil), pool...)
	sort.SliceStable(ordered, func(i, j int) bool {
		a, b := ordered[i], ordered[j]
		if a.account.Priority != b.account.Priority {
			return a.account.Priority < b.account.Priority
	placeholder
		if a.loadInfo.LoadRate != b.loadInfo.LoadRate {
			return a.loadInfo.LoadRate < b.loadInfo.LoadRate
	placeholder
		if a.loadInfo.WaitingCount != b.loadInfo.WaitingCount {
			return a.loadInfo.WaitingCount < b.loadInfo.WaitingCount
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
	return ordered
placeholder

func (s *defaultOpenAIAccountScheduler) tryAcquireOpenAISelectionOrder(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
	selectionOrder []openAIAccountCandidateScore,
) (*AccountSelectionResult, bool, error) {
	budget := newOpenAISelectionProbeBudget()
	budget.enableLimit()
	return s.tryAcquireOpenAISelectionOrderWithBudget(ctx, req, selectionOrder, budget)
placeholder

func (s *defaultOpenAIAccountScheduler) tryAcquireOpenAISelectionOrderWithBudget(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
	selectionOrder []openAIAccountCandidateScore,
	budget *openAISelectionProbeBudget,
) (*AccountSelectionResult, bool, error) {
	compactBlocked := false
	release := func(result *AcquireResult) {
		if result != nil && result.ReleaseFunc != nil {
			result.ReleaseFunc()
	placeholder
placeholder
	for i := 0; i < len(selectionOrder); i++ {
		candidate := selectionOrder[i]
		if candidate.account == nil {
			continue
	placeholder
		if candidate.loadKnown && candidate.account.Concurrency > 0 &&
			candidate.loadInfo.CurrentConcurrency >= candidate.account.Concurrency {
			continue
	placeholder

		result, attempted, acquireErr := s.tryAcquireOpenAIAccountSlot(ctx, candidate.account.ID, candidate.account.Concurrency, budget)
		if !attempted {
			break
	placeholder
		if acquireErr != nil {
			return nil, compactBlocked, acquireErr
	placeholder
		if result == nil || !result.Acquired {
			continue
	placeholder

		fresh := s.service.resolveFreshSchedulableOpenAIAccount(ctx, candidate.account, req.Platform, req.RequestedModel, false, req.RequiredCapability)
		if fresh == nil || !s.isAccountTransportCompatible(fresh, req.RequiredTransport) || !s.isAccountRequestCompatible(ctx, fresh, req) {
			release(result)
			continue
	placeholder
		if !s.consumeOpenAISelectionDBRecheck(budget) {
			release(result)
			break
	placeholder
		fresh = s.service.recheckSelectedOpenAIAccountFromDB(ctx, fresh, req.GroupID, req.Platform, req.RequestedModel, false, req.RequiredCapability)
		if fresh == nil || !s.isAccountTransportCompatible(fresh, req.RequiredTransport) || !s.isAccountRequestCompatible(ctx, fresh, req) {
			release(result)
			continue
	placeholder
		if req.RequireCompact && openAICompactSupportTier(fresh) == 0 {
			compactBlocked = true
			release(result)
			continue
	placeholder

		if fresh.Concurrency != candidate.account.Concurrency {
			release(result)
			result, attempted, acquireErr = s.tryAcquireOpenAIAccountSlot(ctx, fresh.ID, fresh.Concurrency, budget)
			if !attempted {
				continue
		placeholder
			if acquireErr != nil {
				return nil, compactBlocked, acquireErr
		placeholder
			if result == nil || !result.Acquired {
				continue
		placeholder
	placeholder
		if req.SessionHash != "" && !req.PreserveStickyBinding {
			_ = s.service.bindOpenAIStickySessionDuringSelection(ctx, req.GroupID, req.SessionHash, fresh.ID)
	placeholder
		return attachSelectionProfitGate(ctx, &AccountSelectionResult{
			Account:     fresh,
			Acquired:    true,
			ReleaseFunc: result.ReleaseFunc,
	placeholder), compactBlocked, nil
placeholder
	return nil, compactBlocked, nil
placeholder

func (s *defaultOpenAIAccountScheduler) tryAcquireOpenAIAccountSlot(
	ctx context.Context,
	accountID int64,
	maxConcurrency int,
	budget *openAISelectionProbeBudget,
) (*AcquireResult, bool, error) {
	if s.service.concurrencyService != nil && maxConcurrency > 0 && !budget.recordAcquire(accountID) {
		return nil, false, nil
placeholder
	result, err := s.service.tryAcquireAccountSlot(ctx, accountID, maxConcurrency)
	return result, true, err
placeholder

func (s *defaultOpenAIAccountScheduler) consumeOpenAISelectionDBRecheck(budget *openAISelectionProbeBudget) bool {
	if s.service.schedulerSnapshot == nil || s.service.accountRepo == nil {
		return true
placeholder
	return budget.recordRecheck()
placeholder

func (s *defaultOpenAIAccountScheduler) tryFallbackToWeightedSticky(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
) (*AccountSelectionResult, error) {
	if !req.StickyWeighted {
		return nil, nil
placeholder
	for _, accountID := range []int64{req.StickyPreviousAccountID, req.StickyAccountIDplaceholder {
		if accountID <= 0 {
			continue
	placeholder
		if req.ExcludedIDs != nil {
			if _, excluded := req.ExcludedIDs[accountID]; excluded {
				continue
		placeholder
	placeholder
		account, err := s.service.getSchedulableAccount(ctx, accountID)
		if err != nil || account == nil {
			continue
	placeholder
		if !s.isAccountRequestCompatible(ctx, account, req) || !s.isAccountTransportCompatible(account, req.RequiredTransport) {
			continue
	placeholder
		account = s.service.recheckSelectedOpenAIAccountFromDB(ctx, account, req.GroupID, req.Platform, req.RequestedModel, req.RequireCompact, req.RequiredCapability)
		if account == nil {
			if accountID == req.StickyAccountID && strings.TrimSpace(req.SessionHash) != "" {
				_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, req.SessionHash)
		placeholder
			continue
	placeholder
		if !s.service.openAIAccountMatchesSchedulingGroup(account, req.GroupID) {
			if accountID == req.StickyAccountID && strings.TrimSpace(req.SessionHash) != "" {
				_ = s.service.deleteStickySessionAccountID(ctx, req.GroupID, req.SessionHash)
		placeholder
			continue
	placeholder
		if !s.isAccountRequestCompatible(ctx, account, req) || !s.isAccountTransportCompatible(account, req.RequiredTransport) {
			continue
	placeholder
		if req.RequireCompact && openAICompactSupportTier(account) == 0 {
			continue
	placeholder
		result, acquireErr := s.service.tryAcquireAccountSlot(ctx, account.ID, account.Concurrency)
		if acquireErr != nil {
			return nil, acquireErr
	placeholder
		if result != nil && result.Acquired {
			if req.SessionHash != "" && !req.PreserveStickyBinding {
				_ = s.service.bindOpenAIStickySessionDuringSelection(ctx, req.GroupID, req.SessionHash, account.ID)
		placeholder
			return attachSelectionProfitGate(ctx, &AccountSelectionResult{
				Account:     account,
				Acquired:    true,
				ReleaseFunc: result.ReleaseFunc,
		placeholder), nil
	placeholder
		if s.service.concurrencyService != nil {
			cfg := s.service.schedulingConfig()
			return attachSelectionProfitGate(ctx, &AccountSelectionResult{
				Account: account,
				WaitPlan: &AccountWaitPlan{
					AccountID:      account.ID,
					MaxConcurrency: account.Concurrency,
					Timeout:        cfg.StickySessionWaitTimeout,
					MaxWaiting:     cfg.StickySessionMaxWaiting,
			placeholder,
		placeholder), nil
	placeholder
placeholder
	return nil, nil
placeholder

// openAISelectionFilterStats counts why candidates were dropped by the
// selectByLoadBalance initial filter. Historically these exclusions were
// silent (debug logs at best), so a "no available accounts" failure with
// excluded_account_count=0 was undiagnosable from the error alone (#4599).
// The reasons map is lazily allocated: on the happy path (nothing filtered
// out, or an account is eventually selected) no extra allocation happens.
type openAISelectionFilterStats struct {
	pool    int
	reasons map[string]int
placeholder

func (s *openAISelectionFilterStats) exclude(reason string) {
	if s.reasons == nil {
		s.reasons = make(map[string]int, 4)
placeholder
	s.reasons[reason]++
placeholder

// summary renders deterministic exclusion statistics for scheduling error
// messages, e.g. "pool=3, filtered: model_not_supported=2 quota_auto_pause_7d=1".
// Reasons are sorted lexicographically so the output is stable for tests and
// log aggregation. extra, when non-empty, is appended as a trailing marker.
func (s openAISelectionFilterStats) summary(extra string) string {
	var b strings.Builder
	_, _ = b.WriteString("pool=")
	_, _ = b.WriteString(strconv.Itoa(s.pool))
	if len(s.reasons) > 0 {
		reasons := make([]string, 0, len(s.reasons))
		for reason := range s.reasons {
			reasons = append(reasons, reason)
	placeholder
		sort.Strings(reasons)
		_, _ = b.WriteString(", filtered:")
		for _, reason := range reasons {
			_, _ = b.WriteString(" ")
			_, _ = b.WriteString(reason)
			_, _ = b.WriteString("=")
			_, _ = b.WriteString(strconv.Itoa(s.reasons[reason]))
	placeholder
placeholder
	if extra != "" {
		_, _ = b.WriteString(", ")
		_, _ = b.WriteString(extra)
placeholder
	return b.String()
placeholder

func (s *defaultOpenAIAccountScheduler) selectByLoadBalance(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
) (*AccountSelectionResult, int, int, float64, error) {
	budget := newOpenAISelectionProbeBudget()
	accounts, err := s.service.listSchedulableAccounts(ctx, req.GroupID, req.Platform)
	if err != nil {
		return nil, 0, 0, 0, err
placeholder
	if len(accounts) == 0 {
		return nil, 0, 0, 0, noAvailableOpenAISelectionError(req.RequestedModel, false, openAISelectionFilterStats{placeholder.summary(""))
placeholder
	// Local free-tier soft gate on the Grok scheduling path only (not admin probe).
	accounts = s.filterGrokFreeQuotaAccounts(ctx, accounts)
	if len(accounts) == 0 {
		return nil, 0, 0, 0, noAvailableOpenAISelectionError(req.RequestedModel, false, openAISelectionFilterStats{placeholder.summary("grok_free_quota_soft_gate"))
placeholder

	// require_privacy_set: 获取分组信息
	var schedGroup *Group
	if req.GroupID != nil && s.service.schedulerSnapshot != nil {
		schedGroup, _ = s.service.schedulerSnapshot.GetGroupByID(ctx, *req.GroupID)
placeholder

	filterStats := openAISelectionFilterStats{pool: len(accounts)placeholder
	filtered := make([]*Account, 0, len(accounts))
	loadReq := make([]AccountWithConcurrency, 0, len(accounts))
	for i := range accounts {
		account := &accounts[i]
		if req.ExcludedIDs != nil {
			if _, excluded := req.ExcludedIDs[account.ID]; excluded {
				filterStats.exclude("excluded")
				continue
		placeholder
	placeholder
		if !account.IsSchedulable() {
			filterStats.exclude("not_schedulable")
			continue
	placeholder
		if account.Platform != normalizeOpenAICompatiblePlatform(req.Platform) || !account.IsOpenAICompatible() {
			filterStats.exclude("platform_mismatch")
			continue
	placeholder
		if s.service.isOpenAIAccountRequestRuntimeBlocked(account, req.RequestedModel) {
			filterStats.exclude("runtime_blocked")
			continue
	placeholder
		// require_privacy_set: 跳过 privacy 未设置的账号并标记异常
		if schedGroup != nil && schedGroup.RequirePrivacySet && !account.IsPrivacySet() {
			s.service.BlockAccountScheduling(account, time.Time{placeholder, "privacy_not_set")
			_ = s.service.accountRepo.SetError(ctx, account.ID,
				fmt.Sprintf("Privacy not set, required by group [%s]", schedGroup.Name))
			filterStats.exclude("privacy_not_set")
			continue
	placeholder
		if compatible, reason := s.isAccountRequestCompatibleReason(ctx, account, req); !compatible {
			filterStats.exclude(reason)
			continue
	placeholder
		if !s.isAccountTransportCompatible(account, req.RequiredTransport) {
			filterStats.exclude("transport_incompatible")
			continue
	placeholder
		filtered = append(filtered, account)
		loadReq = append(loadReq, AccountWithConcurrency{
			ID:             account.ID,
			MaxConcurrency: account.EffectiveLoadFactor(),
	placeholder)
placeholder
	if len(filtered) == 0 {
		return nil, 0, 0, 0, noAvailableOpenAISelectionError(req.RequestedModel, false, filterStats.summary(""))
placeholder

	loadMap := map[int64]*AccountLoadInfo{placeholder
	if s.service.concurrencyService != nil {
		if batchLoad, loadErr := s.service.concurrencyService.GetAccountsLoadBatch(ctx, loadReq); loadErr == nil {
			loadMap = batchLoad
	placeholder
placeholder

	if req.SubscriptionPriority {
		subscriptionAccounts, regularAccounts := partitionOpenAIChatGPTSubscriptionAccounts(filtered)
		if len(subscriptionAccounts) > 0 {
			attempt := s.trySelectByLoadBalancePool(ctx, req, subscriptionAccounts, loadMap, budget)
			if attempt.err != nil && (!attempt.noCompactCandidates || len(regularAccounts) <= 0) {
				return nil, attempt.candidateCount, attempt.topK, attempt.loadSkew, attempt.err
		placeholder
			if attempt.result != nil {
				return attempt.result, attempt.candidateCount, attempt.topK, attempt.loadSkew, nil
		placeholder
			if len(regularAccounts) > 0 {
				regularAttempt := s.trySelectByLoadBalancePool(ctx, req, regularAccounts, loadMap, budget)
				if regularAttempt.err != nil && !regularAttempt.noCompactCandidates {
					return nil, regularAttempt.candidateCount, regularAttempt.topK, regularAttempt.loadSkew, regularAttempt.err
			placeholder
				if regularAttempt.result != nil {
					return regularAttempt.result, regularAttempt.candidateCount, regularAttempt.topK, regularAttempt.loadSkew, nil
			placeholder
				var result *AccountSelectionResult
				candidateCount, topK, loadSkew := regularAttempt.candidateCount, regularAttempt.topK, regularAttempt.loadSkew
				fallbackErr := regularAttempt.err
				if regularAttempt.err == nil {
					result, candidateCount, topK, loadSkew, fallbackErr = s.finishLoadBalanceSelectionFallback(ctx, req, regularAttempt, budget, filterStats)
					if fallbackErr == nil && result != nil {
						return result, candidateCount, topK, loadSkew, nil
				placeholder
			placeholder
				// 常规池既无法获取也无法排队（含仅剩不支持 compact 的候选）时，
				// 回退到订阅池的等待计划：busy-but-waitable 的订阅账号不应因常规池存在
				// 而被丢弃，否则开启订阅优先反而让本可排队成功的请求硬失败。
				subResult, subCandidateCount, subTopK, subLoadSkew, subErr := s.finishLoadBalanceSelectionFallback(ctx, req, attempt, budget, filterStats)
				if subErr == nil && subResult != nil {
					return subResult, subCandidateCount, subTopK, subLoadSkew, nil
			placeholder
				return result, candidateCount, topK, loadSkew, fallbackErr
		placeholder
			return s.finishLoadBalanceSelectionFallback(ctx, req, attempt, budget, filterStats)
	placeholder
placeholder

	attempt := s.trySelectByLoadBalancePool(ctx, req, filtered, loadMap, budget)
	if attempt.err != nil {
		return nil, attempt.candidateCount, attempt.topK, attempt.loadSkew, attempt.err
placeholder
	if attempt.result != nil {
		return attempt.result, attempt.candidateCount, attempt.topK, attempt.loadSkew, nil
placeholder
	return s.finishLoadBalanceSelectionFallback(ctx, req, attempt, budget, filterStats)
placeholder

func partitionOpenAIChatGPTSubscriptionAccounts(accounts []*Account) ([]*Account, []*Account) {
	subscriptionAccounts := make([]*Account, 0, len(accounts))
	regularAccounts := make([]*Account, 0, len(accounts))
	for _, account := range accounts {
		if account != nil && account.IsOpenAIChatGPTSubscription() {
			subscriptionAccounts = append(subscriptionAccounts, account)
			continue
	placeholder
		regularAccounts = append(regularAccounts, account)
placeholder
	return subscriptionAccounts, regularAccounts
placeholder

func (s *defaultOpenAIAccountScheduler) trySelectByLoadBalancePool(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
	filtered []*Account,
	loadMap map[int64]*AccountLoadInfo,
	budget *openAISelectionProbeBudget,
) openAIAccountLoadSelectionAttempt {
	plan := s.buildOpenAIAccountLoadPlan(ctx, req, filtered, loadMap)
	if openAICostOverflowExpanded(req, plan) {
		budget.enableLimit()
placeholder
	attempt := openAIAccountLoadSelectionAttempt{
		selectionOrder: plan.selectionOrder,
		candidateCount: plan.candidateCount,
		topK:           plan.topK,
		loadSkew:       plan.loadSkew,
placeholder
	if req.RequireCompact && len(plan.candidates) == 0 && len(plan.staleSnapshotCompactRetry) == 0 {
		attempt.noCompactCandidates = true
		attempt.err = ErrNoAvailableCompactAccounts
		return attempt
placeholder
	if req.RequireCompact && len(attempt.selectionOrder) == 0 && s.service.schedulerSnapshot == nil {
		attempt.noCompactCandidates = true
		attempt.err = ErrNoAvailableCompactAccounts
		return attempt
placeholder
	if len(attempt.selectionOrder) == 0 {
		attempt.compactBlocked = req.RequireCompact && len(plan.allCandidates) > 0
		return attempt
placeholder

	result, compactBlocked, acquireErr := s.tryAcquireOpenAISelectionOrderWithBudget(ctx, req, attempt.selectionOrder, budget)
	attempt.compactBlocked = compactBlocked
	if acquireErr != nil {
		attempt.err = acquireErr
		return attempt
placeholder
	if result != nil {
		attempt.result = result
		return attempt
placeholder

	if s.service.concurrencyService != nil && !budget.acquireExhausted() {
		loadReq := buildOpenAIAccountLoadRequest(filtered)
		if freshLoadMap, loadErr := s.service.concurrencyService.GetAccountsLoadBatchFresh(ctx, loadReq); loadErr == nil {
			freshPlan := s.buildOpenAIAccountLoadPlan(ctx, req, filtered, freshLoadMap)
			if openAICostOverflowExpanded(req, freshPlan) {
				budget.enableLimit()
		placeholder
			if len(freshPlan.selectionOrder) > 0 {
				freshResult, freshCompactBlocked, freshAcquireErr := s.tryAcquireOpenAISelectionOrderWithBudget(ctx, req, freshPlan.selectionOrder, budget)
				if freshAcquireErr != nil {
					attempt.err = freshAcquireErr
					return attempt
			placeholder
				if freshResult != nil {
					attempt.result = freshResult
					attempt.selectionOrder = freshPlan.selectionOrder
					attempt.candidateCount = freshPlan.candidateCount
					attempt.topK = freshPlan.topK
					attempt.loadSkew = freshPlan.loadSkew
					return attempt
			placeholder
				attempt.compactBlocked = attempt.compactBlocked || freshCompactBlocked
				attempt.selectionOrder = freshPlan.selectionOrder
				attempt.candidateCount = freshPlan.candidateCount
				attempt.topK = freshPlan.topK
				attempt.loadSkew = freshPlan.loadSkew
		placeholder
	placeholder
placeholder

	return attempt
placeholder

func openAICostOverflowExpanded(req OpenAIAccountScheduleRequest, plan openAIAccountLoadPlan) bool {
	if !plan.includeOverflowFallback || plan.topK <= 0 {
		return false
placeholder
	if !req.RequireCompact {
		return len(plan.candidates) > plan.topK
placeholder
	supported, unknown := 0, 0
	for _, candidate := range plan.candidates {
		switch openAICompactSupportTier(candidate.account) {
		case 2:
			supported++
		case 1:
			unknown++
	placeholder
placeholder
	return supported > plan.topK || unknown > plan.topK
placeholder

func buildOpenAIAccountLoadRequest(accounts []*Account) []AccountWithConcurrency {
	loadReq := make([]AccountWithConcurrency, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
	placeholder
		loadReq = append(loadReq, AccountWithConcurrency{
			ID:             account.ID,
			MaxConcurrency: account.EffectiveLoadFactor(),
	placeholder)
placeholder
	return loadReq
placeholder

func (s *defaultOpenAIAccountScheduler) finishLoadBalanceSelectionFallback(
	ctx context.Context,
	req OpenAIAccountScheduleRequest,
	attempt openAIAccountLoadSelectionAttempt,
	budget *openAISelectionProbeBudget,
	filterStats openAISelectionFilterStats,
) (*AccountSelectionResult, int, int, float64, error) {
	candidateCount := attempt.candidateCount
	topK := attempt.topK
	loadSkew := attempt.loadSkew

	if len(attempt.selectionOrder) == 0 {
		return nil, candidateCount, topK, loadSkew, noAvailableOpenAISelectionError(req.RequestedModel, attempt.compactBlocked, filterStats.summary("selection_order_empty"))
placeholder

	if stickyFallback, stickyErr := s.tryFallbackToWeightedSticky(ctx, req); stickyErr != nil {
		return nil, candidateCount, topK, loadSkew, stickyErr
placeholder else if stickyFallback != nil {
		return stickyFallback, candidateCount, topK, loadSkew, nil
placeholder

	cfg := s.service.schedulingConfig()
	compactBlocked := attempt.compactBlocked
	// WaitPlan.MaxConcurrency 使用 Concurrency（非 EffectiveLoadFactor），因为 WaitPlan 控制的是 Redis 实际并发槽位等待。
	passes := 1
	if budget != nil && budget.limited {
		passes = 4
placeholder
	for pass := 0; pass < passes; pass++ {
		wantAttempted := pass == 1 || pass == 3
		wantKnownFull := pass >= 2
		for _, candidate := range attempt.selectionOrder {
			if candidate.account == nil {
				continue
		placeholder
			if budget != nil && budget.limited {
				knownFull := candidate.loadKnown && candidate.account.Concurrency > 0 &&
					candidate.loadInfo.CurrentConcurrency >= candidate.account.Concurrency
				if budget.wasAttempted(candidate.account.ID) != wantAttempted || knownFull != wantKnownFull {
					continue
			placeholder
		placeholder
			fresh := s.service.resolveFreshSchedulableOpenAIAccount(ctx, candidate.account, req.Platform, req.RequestedModel, false, req.RequiredCapability)
			if fresh == nil || !s.isAccountTransportCompatible(fresh, req.RequiredTransport) || !s.isAccountRequestCompatible(ctx, fresh, req) {
				continue
		placeholder
			if !s.consumeOpenAISelectionDBRecheck(budget) {
				return nil, candidateCount, topK, loadSkew, noAvailableOpenAISelectionError(req.RequestedModel, compactBlocked, filterStats.summary("selection_order_exhausted"))
		placeholder
			fresh = s.service.recheckSelectedOpenAIAccountFromDB(ctx, fresh, req.GroupID, req.Platform, req.RequestedModel, false, req.RequiredCapability)
			if fresh == nil || !s.isAccountTransportCompatible(fresh, req.RequiredTransport) || !s.isAccountRequestCompatible(ctx, fresh, req) {
				continue
		placeholder
			if req.RequireCompact && openAICompactSupportTier(fresh) == 0 {
				compactBlocked = true
				continue
		placeholder
			return attachSelectionProfitGate(ctx, &AccountSelectionResult{
				Account: fresh,
				WaitPlan: &AccountWaitPlan{
					AccountID:      fresh.ID,
					MaxConcurrency: fresh.Concurrency,
					Timeout:        cfg.FallbackWaitTimeout,
					MaxWaiting:     cfg.FallbackMaxWaiting,
			placeholder,
		placeholder), candidateCount, topK, loadSkew, nil
	placeholder
placeholder

	return nil, candidateCount, topK, loadSkew, noAvailableOpenAISelectionError(req.RequestedModel, compactBlocked, filterStats.summary("selection_order_exhausted"))
placeholder

func (s *defaultOpenAIAccountScheduler) isAccountTransportCompatible(account *Account, requiredTransport OpenAIUpstreamTransport) bool {
	if requiredTransport == OpenAIUpstreamTransportAny || requiredTransport == OpenAIUpstreamTransportHTTPSSE {
		return true
placeholder
	if s == nil || s.service == nil {
		return false
placeholder
	return s.service.isOpenAIAccountTransportCompatible(account, requiredTransport)
placeholder

func (s *defaultOpenAIAccountScheduler) lookupShadowParentAccount(ctx context.Context, id int64) *Account {
	if s == nil || s.service == nil {
		return nil
placeholder
	if s.service.schedulerSnapshot != nil {
		if account, err := s.service.schedulerSnapshot.GetAccount(ctx, id); err == nil && account != nil {
			return account
	placeholder
placeholder
	if s.service.accountRepo == nil {
		return nil
placeholder
	account, _ := s.service.accountRepo.GetByID(ctx, id)
	return account
placeholder

func (s *defaultOpenAIAccountScheduler) isAccountRequestCompatible(ctx context.Context, account *Account, req OpenAIAccountScheduleRequest) bool {
	compatible, _ := s.isAccountRequestCompatibleReason(ctx, account, req)
	return compatible
placeholder

// isAccountRequestCompatibleReason reports whether the account can serve the
// request, and when it cannot, names the veto point. The reason feeds
// openAISelectionFilterStats so that "no available accounts" errors state why
// each candidate was dropped instead of failing silently (#4599).
func (s *defaultOpenAIAccountScheduler) isAccountRequestCompatibleReason(ctx context.Context, account *Account, req OpenAIAccountScheduleRequest) (bool, string) {
	if account == nil {
		return false, "account_nil"
placeholder
	if s != nil && s.service != nil && s.service.isOpenAIAccountRequestRuntimeBlocked(account, req.RequestedModel) {
		return false, "runtime_blocked"
placeholder
	if s != nil && s.service != nil && s.service.isOpenAIProxyStreamQuarantined(ctx, account) {
		return false, "proxy_stream_quarantined"
placeholder
	// Quota auto-pause must be evaluated during the initial filter too. Without it the
	// TopK candidate pool can be filled with paused accounts and the later fresh/DB
	// rechecks won't reach healthy accounts that fell outside TopK — manifesting as
	// "no available accounts" even though healthy ones exist.
	if paused, decision := shouldAutoPauseOpenAIAccountByQuota(ctx, account); paused {
		reason := "quota_auto_pause"
		if decision.window != "" {
			reason += "_" + decision.window
	placeholder
		return false, reason
placeholder
	// 母账号健康联动：影子账号的凭据来自母账号，母账号不可调度时影子也不应被选中。
	// Parent-health gate: shadow borrows the parent's credentials; an unschedulable
	// parent must block the shadow across all scheduler paths.
	if !parentHealthyForShadow(account, func(id int64) *Account {
		return s.lookupShadowParentAccount(ctx, id)
placeholder) {
		return false, "shadow_parent_unhealthy"
placeholder
	if req.RequestedModel != "" && !account.IsModelSupported(req.RequestedModel) {
		return false, "model_not_supported"
placeholder
	if req.GroupID != nil && s != nil && s.service != nil &&
		s.service.needsUpstreamChannelRestrictionCheck(ctx, req.GroupID) &&
		s.service.isUpstreamModelRestrictedByChannel(ctx, *req.GroupID, account, req.RequestedModel, req.RequireCompact) {
		return false, "channel_upstream_restricted"
placeholder
	if !accountSupportsOpenAICapabilities(account, req.RequiredCapability, req.RequiredImageCapability) {
		return false, "capability_mismatch"
placeholder
	// 分组利润控制：不合格账号在候选过滤与抢槽后终检阶段即被排除，
	// 排序/评分/粘性/熔断只在合格账号之间工作；named reason 进入 filter stats。
	if vetoed, reason := openAIProfitControlVetoReason(ctx, account); vetoed {
		return false, reason
placeholder
	return true, ""
placeholder

func (s *defaultOpenAIAccountScheduler) ReportResult(accountID int64, success bool, firstTokenMs *int) {
	if s == nil || s.stats == nil {
		return
placeholder
	s.stats.report(accountID, success, firstTokenMs)
placeholder

func (s *defaultOpenAIAccountScheduler) ReportSwitch() {
	if s == nil {
		return
placeholder
	s.metrics.recordSwitch()
placeholder

func (s *defaultOpenAIAccountScheduler) SnapshotMetrics() OpenAIAccountSchedulerMetricsSnapshot {
	if s == nil {
		return OpenAIAccountSchedulerMetricsSnapshot{placeholder
placeholder

	selectTotal := s.metrics.selectTotal.Load()
	prevHit := s.metrics.stickyPreviousHitTotal.Load()
	sessionHit := s.metrics.stickySessionHitTotal.Load()
	switchTotal := s.metrics.accountSwitchTotal.Load()
	latencyTotal := s.metrics.latencyMsTotal.Load()
	loadSkewTotal := s.metrics.loadSkewMilliTotal.Load()

	snapshot := OpenAIAccountSchedulerMetricsSnapshot{
		SelectTotal:              selectTotal,
		StickyPreviousHitTotal:   prevHit,
		StickySessionHitTotal:    sessionHit,
		LoadBalanceSelectTotal:   s.metrics.loadBalanceSelectTotal.Load(),
		AccountSwitchTotal:       switchTotal,
		SchedulerLatencyMsTotal:  latencyTotal,
		RuntimeStatsAccountCount: s.stats.size(),
placeholder
	if selectTotal > 0 {
		snapshot.SchedulerLatencyMsAvg = float64(latencyTotal) / float64(selectTotal)
		snapshot.StickyHitRatio = float64(prevHit+sessionHit) / float64(selectTotal)
		snapshot.AccountSwitchRate = float64(switchTotal) / float64(selectTotal)
		snapshot.LoadSkewAvg = float64(loadSkewTotal) / 1000 / float64(selectTotal)
placeholder
	return snapshot
placeholder

func (s *OpenAIGatewayService) openAIAdvancedSchedulerSettingRepo() SettingRepository {
	if s == nil || s.rateLimitService == nil || s.rateLimitService.settingService == nil {
		return nil
placeholder
	return s.rateLimitService.settingService.settingRepo
placeholder

func (s *OpenAIGatewayService) openAIAdvancedSchedulerRuntimeSettings(ctx context.Context) openAIAdvancedSchedulerRuntimeSettings {
	if cached, ok := openAIAdvancedSchedulerSettingCache.Load().(*cachedOpenAIAdvancedSchedulerSetting); ok && cached != nil {
		if time.Now().UnixNano() < cached.expiresAt {
			return openAIAdvancedSchedulerRuntimeSettings{
				lowUpstreamRatePriorityEnabled: cached.lowUpstreamRatePriorityEnabled,
				oauthSchedulingRateMultiplier:  cached.oauthSchedulingRateMultiplier,
				enabled:                        cached.enabled,
				stickyWeightedEnabled:          cached.stickyWeightedEnabled,
				subscriptionPriorityEnabled:    cached.subscriptionPriorityEnabled,
				lbTopKOverride:                 cached.lbTopKOverride,
				weightOverrides:                cloneOpenAIAdvancedSchedulerWeightOverrides(cached.weightOverrides),
		placeholder
	placeholder
placeholder

	result, _, _ := openAIAdvancedSchedulerSettingSF.Do(openAIAdvancedSchedulerSettingKey, func() (any, error) {
		if cached, ok := openAIAdvancedSchedulerSettingCache.Load().(*cachedOpenAIAdvancedSchedulerSetting); ok && cached != nil {
			if time.Now().UnixNano() < cached.expiresAt {
				return openAIAdvancedSchedulerRuntimeSettings{
					lowUpstreamRatePriorityEnabled: cached.lowUpstreamRatePriorityEnabled,
					oauthSchedulingRateMultiplier:  cached.oauthSchedulingRateMultiplier,
					enabled:                        cached.enabled,
					stickyWeightedEnabled:          cached.stickyWeightedEnabled,
					subscriptionPriorityEnabled:    cached.subscriptionPriorityEnabled,
					lbTopKOverride:                 cached.lbTopKOverride,
					weightOverrides:                cloneOpenAIAdvancedSchedulerWeightOverrides(cached.weightOverrides),
			placeholder, nil
		placeholder
	placeholder

		lowUpstreamRatePriorityEnabled := false
		oauthSchedulingRateMultiplier := defaultOpenAIOAuthSchedulingRateMultiplier
		enabled := false
		stickyWeightedEnabled := false
		subscriptionPriorityEnabled := false
		lbTopKOverride := 0
		weightOverrides := map[string]float64{placeholder
		if repo := s.openAIAdvancedSchedulerSettingRepo(); repo != nil {
			dbCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), openAIAdvancedSchedulerSettingDBTimeout)
			defer cancel()

			if values, err := repo.GetMultiple(dbCtx, openAIAdvancedSchedulerRuntimeSettingKeys()); err == nil {
				lowUpstreamRatePriorityEnabled = strings.EqualFold(strings.TrimSpace(values[SettingKeyOpenAILowUpstreamRatePriorityEnabled]), "true")
				oauthSchedulingRateMultiplier = parseOpenAIOAuthSchedulingRateMultiplier(values[SettingKeyOpenAIOAuthSchedulingRateMultiplier])
				enabled = strings.EqualFold(strings.TrimSpace(values[openAIAdvancedSchedulerSettingKey]), "true")
				stickyWeightedEnabled = strings.EqualFold(strings.TrimSpace(values[SettingKeyOpenAIAdvancedSchedulerStickyWeightedEnabled]), "true")
				subscriptionPriorityEnabled = strings.EqualFold(strings.TrimSpace(values[SettingKeyOpenAIAdvancedSchedulerSubscriptionPriorityEnabled]), "true")
				lbTopKOverride = parsePositiveIntOverride(values[SettingKeyOpenAIAdvancedSchedulerLBTopK])
				weightOverrides = parseOpenAIAdvancedSchedulerWeightOverrides(values)
		placeholder else {
				// 批量读取失败时逐键降级，覆盖全部键（含 TopK/权重），避免只加载布尔开关
				// 而静默丢弃管理员配置的覆盖值；降级状态会被缓存一个 TTL，必须留痕。
				slog.Warn("openai_advanced_scheduler_settings_batch_load_failed", "error", err)
				fallbackValues := make(map[string]string)
				for _, key := range openAIAdvancedSchedulerRuntimeSettingKeys() {
					if value, valueErr := repo.GetValue(dbCtx, key); valueErr == nil {
						fallbackValues[key] = value
				placeholder
			placeholder
				lowUpstreamRatePriorityEnabled = strings.EqualFold(strings.TrimSpace(fallbackValues[SettingKeyOpenAILowUpstreamRatePriorityEnabled]), "true")
				oauthSchedulingRateMultiplier = parseOpenAIOAuthSchedulingRateMultiplier(fallbackValues[SettingKeyOpenAIOAuthSchedulingRateMultiplier])
				enabled = strings.EqualFold(strings.TrimSpace(fallbackValues[openAIAdvancedSchedulerSettingKey]), "true")
				stickyWeightedEnabled = strings.EqualFold(strings.TrimSpace(fallbackValues[SettingKeyOpenAIAdvancedSchedulerStickyWeightedEnabled]), "true")
				subscriptionPriorityEnabled = strings.EqualFold(strings.TrimSpace(fallbackValues[SettingKeyOpenAIAdvancedSchedulerSubscriptionPriorityEnabled]), "true")
				lbTopKOverride = parsePositiveIntOverride(fallbackValues[SettingKeyOpenAIAdvancedSchedulerLBTopK])
				weightOverrides = parseOpenAIAdvancedSchedulerWeightOverrides(fallbackValues)
		placeholder
	placeholder

		openAIAdvancedSchedulerSettingCache.Store(&cachedOpenAIAdvancedSchedulerSetting{
			lowUpstreamRatePriorityEnabled: lowUpstreamRatePriorityEnabled,
			oauthSchedulingRateMultiplier:  oauthSchedulingRateMultiplier,
			enabled:                        enabled,
			stickyWeightedEnabled:          stickyWeightedEnabled,
			subscriptionPriorityEnabled:    subscriptionPriorityEnabled,
			lbTopKOverride:                 lbTopKOverride,
			weightOverrides:                cloneOpenAIAdvancedSchedulerWeightOverrides(weightOverrides),
			expiresAt:                      time.Now().Add(openAIAdvancedSchedulerSettingCacheTTL).UnixNano(),
	placeholder)
		return openAIAdvancedSchedulerRuntimeSettings{
			lowUpstreamRatePriorityEnabled: lowUpstreamRatePriorityEnabled,
			oauthSchedulingRateMultiplier:  oauthSchedulingRateMultiplier,
			enabled:                        enabled,
			stickyWeightedEnabled:          stickyWeightedEnabled,
			subscriptionPriorityEnabled:    subscriptionPriorityEnabled,
			lbTopKOverride:                 lbTopKOverride,
			weightOverrides:                weightOverrides,
	placeholder, nil
placeholder)

	settings, _ := result.(openAIAdvancedSchedulerRuntimeSettings)
	return settings
placeholder

func (s *OpenAIGatewayService) isOpenAIAdvancedSchedulerEnabled(ctx context.Context) bool {
	return s.openAIAdvancedSchedulerRuntimeSettings(ctx).enabled
placeholder

func (s *OpenAIGatewayService) isOpenAILowUpstreamRatePriorityEnabled(ctx context.Context) bool {
	settings := s.openAIAdvancedSchedulerRuntimeSettings(ctx)
	return !settings.enabled && settings.lowUpstreamRatePriorityEnabled
placeholder

func (s *OpenAIGatewayService) openAIOAuthSchedulingRateMultiplier(ctx context.Context) float64 {
	return s.openAIAdvancedSchedulerRuntimeSettings(ctx).oauthSchedulingRateMultiplier
placeholder

func (s *OpenAIGatewayService) isOpenAIAdvancedSchedulerStickyWeightedEnabled(ctx context.Context) bool {
	settings := s.openAIAdvancedSchedulerRuntimeSettings(ctx)
	return settings.enabled && settings.stickyWeightedEnabled
placeholder

func (s *OpenAIGatewayService) isOpenAIAdvancedSchedulerSubscriptionPriorityEnabled(ctx context.Context) bool {
	settings := s.openAIAdvancedSchedulerRuntimeSettings(ctx)
	return settings.enabled && settings.subscriptionPriorityEnabled
placeholder

func openAIAdvancedSchedulerRuntimeSettingKeys() []string {
	keys := []string{
		SettingKeyOpenAILowUpstreamRatePriorityEnabled,
		SettingKeyOpenAIOAuthSchedulingRateMultiplier,
		openAIAdvancedSchedulerSettingKey,
		SettingKeyOpenAIAdvancedSchedulerStickyWeightedEnabled,
		SettingKeyOpenAIAdvancedSchedulerSubscriptionPriorityEnabled,
		SettingKeyOpenAIAdvancedSchedulerLBTopK,
placeholder
	for _, spec := range openAIAdvancedSchedulerWeightOverrideSpecs() {
		keys = append(keys, spec.key)
placeholder
	return keys
placeholder

type openAIAdvancedSchedulerWeightOverrideSpec struct {
	key  string
	name string
placeholder

func openAIAdvancedSchedulerWeightOverrideSpecs() []openAIAdvancedSchedulerWeightOverrideSpec {
	return []openAIAdvancedSchedulerWeightOverrideSpec{
		{key: SettingKeyOpenAIAdvancedSchedulerWeightPriority, name: "priority"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightLoad, name: "load"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightQueue, name: "queue"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightErrorRate, name: "error_rate"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightTTFT, name: "ttft"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightReset, name: "reset"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightQuotaHeadroom, name: "quota_headroom"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightUpstreamCost, name: "upstream_cost"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightPreviousResponse, name: "previous_response"placeholder,
		{key: SettingKeyOpenAIAdvancedSchedulerWeightSessionSticky, name: "session_sticky"placeholder,
placeholder
placeholder

func parsePositiveIntOverride(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
placeholder
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0
placeholder
	return value
placeholder

func parseOpenAIAdvancedSchedulerWeightOverrides(values map[string]string) map[string]float64 {
	overrides := map[string]float64{placeholder
	for _, spec := range openAIAdvancedSchedulerWeightOverrideSpecs() {
		raw := strings.TrimSpace(values[spec.key])
		if raw == "" {
			continue
	placeholder
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil || value < 0 || math.IsNaN(value) || math.IsInf(value, 0) {
			continue
	placeholder
		overrides[spec.name] = value
placeholder
	return overrides
placeholder

func cloneOpenAIAdvancedSchedulerWeightOverrides(in map[string]float64) map[string]float64 {
	if len(in) == 0 {
		return nil
placeholder
	out := make(map[string]float64, len(in))
	for key, value := range in {
		out[key] = value
placeholder
	return out
placeholder

func (s *OpenAIGatewayService) getOpenAIAccountScheduler(ctx context.Context) OpenAIAccountScheduler {
	if s == nil {
		return nil
placeholder
	if !s.isOpenAIAdvancedSchedulerEnabled(ctx) {
		return nil
placeholder
	s.openaiSchedulerOnce.Do(func() {
		if s.openaiAccountStats == nil {
			s.openaiAccountStats = newOpenAIAccountRuntimeStats()
	placeholder
		if s.openaiScheduler == nil {
			s.openaiScheduler = newDefaultOpenAIAccountScheduler(s, s.openaiAccountStats)
	placeholder
placeholder)
	return s.openaiScheduler
placeholder

func resetOpenAIAdvancedSchedulerSettingCacheForTest() {
	openAIAdvancedSchedulerSettingCache = atomic.Value{placeholder
	openAIAdvancedSchedulerSettingSF = singleflight.Group{placeholder
placeholder

func (s *OpenAIGatewayService) SelectAccountWithScheduler(
	ctx context.Context,
	groupID *int64,
	previousResponseID string,
	sessionHash string,
	requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	requiredTransport OpenAIUpstreamTransport,
	requireCompact bool,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error) {
	return s.selectAccountWithScheduler(ctx, groupID, previousResponseID, sessionHash, requestedModel, excludedIDs, requiredTransport, "", "", requireCompact, PlatformOpenAI, false, true)
placeholder

// SelectAccountWithSchedulerForCapability 按能力要求调度账号。
// previousResponseCanMove 表示首包 input 可自行重建工具续链，previous_response_id 允许跨账号迁移
// （粘性加权模式下改为加权偏好而非硬粘连）。
func (s *OpenAIGatewayService) SelectAccountWithSchedulerForCapability(
	ctx context.Context,
	groupID *int64,
	previousResponseID string,
	sessionHash string,
	requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	requiredTransport OpenAIUpstreamTransport,
	requiredCapability OpenAIEndpointCapability,
	requireCompact bool,
	previousResponseCanMove bool,
	useUpstreamTokenCost bool,
	platformOverride ...string,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error) {
	platform := PlatformOpenAI
	if len(platformOverride) > 0 {
		platform = platformOverride[0]
placeholder
	return s.selectAccountWithScheduler(ctx, groupID, previousResponseID, sessionHash, requestedModel, excludedIDs, requiredTransport, requiredCapability, "", requireCompact, platform, previousResponseCanMove, useUpstreamTokenCost)
placeholder

func (s *OpenAIGatewayService) SelectAccountWithSchedulerForImages(
	ctx context.Context,
	groupID *int64,
	sessionHash string,
	requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	requiredCapability OpenAIImagesCapability,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error) {
	selection, decision, err := s.selectAccountWithScheduler(ctx, groupID, "", sessionHash, requestedModel, excludedIDs, OpenAIUpstreamTransportHTTPSSE, "", requiredCapability, false, PlatformOpenAI, false, false)
	if err == nil && selection != nil && selection.Account != nil {
		return selection, decision, nil
placeholder
	// 如果要求 native 能力（如指定了模型）但没有可用的 APIKey 账号，回退到 basic（OAuth 账号）
	if requiredCapability == OpenAIImagesCapabilityNative {
		return s.selectAccountWithScheduler(ctx, groupID, "", sessionHash, requestedModel, excludedIDs, OpenAIUpstreamTransportHTTPSSE, "", OpenAIImagesCapabilityBasic, false, PlatformOpenAI, false, false)
placeholder
	return selection, decision, err
placeholder

// selectAccountWithScheduler wraps selectAccountWithSchedulerOnce with a
// fail-open second pass for the proxy stream circuit (#5056): when the only
// reason no account is available is that every candidate sits behind a
// quarantined proxy, the quarantine must degrade to a preference instead of
// zeroing out capacity. The retry re-runs the exact same selection with the
// quarantine checks bypassed, so healthy proxies always win the first pass
// and quarantined ones only serve when nothing else can.
func (s *OpenAIGatewayService) selectAccountWithScheduler(
	ctx context.Context,
	groupID *int64,
	previousResponseID string,
	sessionHash string,
	requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	requiredTransport OpenAIUpstreamTransport,
	requiredCapability OpenAIEndpointCapability,
	requiredImageCapability OpenAIImagesCapability,
	requireCompact bool,
	platform string,
	previousResponseCanMove bool,
	useUpstreamTokenCost bool,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error) {
	selection, decision, err := s.selectAccountWithSchedulerOnce(ctx, groupID, previousResponseID, sessionHash, requestedModel, excludedIDs, requiredTransport, requiredCapability, requiredImageCapability, requireCompact, platform, previousResponseCanMove, useUpstreamTokenCost)
	if err == nil || openAIProxyStreamQuarantineBypassed(ctx) {
		return selection, decision, err
placeholder
	if !errors.Is(err, ErrNoAvailableAccounts) && !errors.Is(err, ErrNoAvailableCompactAccounts) {
		return selection, decision, err
placeholder
	// The circuit only ever quarantines PlatformOpenAI accounts.
	if normalizeOpenAICompatiblePlatform(platform) != PlatformOpenAI {
		return selection, decision, err
placeholder
	blocked := s.getOpenAIProxyStreamCircuit().activeBlockCount(time.Now())
	if blocked == 0 {
		return selection, decision, err
placeholder
	s.logOpenAIProxyStreamQuarantineFailOpen(requestedModel, blocked)
	return s.selectAccountWithSchedulerOnce(withOpenAIProxyStreamQuarantineBypass(ctx), groupID, previousResponseID, sessionHash, requestedModel, excludedIDs, requiredTransport, requiredCapability, requiredImageCapability, requireCompact, platform, previousResponseCanMove, useUpstreamTokenCost)
placeholder

func (s *OpenAIGatewayService) selectAccountWithSchedulerOnce(
	ctx context.Context,
	groupID *int64,
	previousResponseID string,
	sessionHash string,
	requestedModel string,
	excludedIDs map[int64]struct{placeholder,
	requiredTransport OpenAIUpstreamTransport,
	requiredCapability OpenAIEndpointCapability,
	requiredImageCapability OpenAIImagesCapability,
	requireCompact bool,
	platform string,
	previousResponseCanMove bool,
	useUpstreamTokenCost bool,
) (*AccountSelectionResult, OpenAIAccountScheduleDecision, error) {
	ctx = s.withOpenAIQuotaAutoPauseContext(ctx)
	// 分组利润控制：唯一文本调度入口的防御性装门。handler 文本
	// 入口已在请求开始经 WithOpenAIRequestPricingContext 装门并固定 pricingAt，
	// 此处对同分组门直接复用（failover 重入阈值稳定），仅为不经 handler 装配的
	// 内部调用兜底。图片/视频调度不在利润门范围：requiredImageCapability 非空的
	// Images 调度不装门；requiredCapability == OpenAIEndpointCapabilityResponses
	// 当前仅显式生图意图的 /v1/responses 设置（HTTP openAIResponsesRequiredCapability
	// 与 WS 桥同款判定），同样不装门——若未来把该 capability 用于非生图流量，
	// 需要同步收窄本条件（有测试钉死该映射）。
	if requiredImageCapability == "" && requiredCapability != OpenAIEndpointCapabilityResponses {
		ctx = s.withOpenAIProfitControlGate(ctx, groupID)
placeholder
	platform = normalizeOpenAICompatiblePlatform(platform)
	decision := OpenAIAccountScheduleDecision{placeholder
	scheduler := s.getOpenAIAccountScheduler(ctx)
	if scheduler == nil {
		decision.Layer = openAIAccountScheduleLayerLoadBalance
		if requiredTransport == OpenAIUpstreamTransportAny || requiredTransport == OpenAIUpstreamTransportHTTPSSE {
			effectiveExcludedIDs := cloneExcludedAccountIDs(excludedIDs)
			for {
				selection, err := s.selectAccountWithLoadAwareness(ctx, groupID, platform, sessionHash, requestedModel, effectiveExcludedIDs, requireCompact, requiredCapability, useUpstreamTokenCost)
				if err != nil {
					return nil, decision, err
			placeholder
				if selection == nil || selection.Account == nil {
					return selection, decision, nil
			placeholder
				if accountSupportsOpenAICapabilities(selection.Account, requiredCapability, requiredImageCapability) {
					return selection, decision, nil
			placeholder
				if selection.ReleaseFunc != nil {
					selection.ReleaseFunc()
			placeholder
				if effectiveExcludedIDs == nil {
					effectiveExcludedIDs = make(map[int64]struct{placeholder)
			placeholder
				if _, exists := effectiveExcludedIDs[selection.Account.ID]; exists {
					return nil, decision, ErrNoAvailableAccounts
			placeholder
				effectiveExcludedIDs[selection.Account.ID] = struct{placeholder{placeholder
		placeholder
	placeholder

		effectiveExcludedIDs := cloneExcludedAccountIDs(excludedIDs)
		for {
			selection, err := s.selectAccountWithLoadAwareness(ctx, groupID, platform, sessionHash, requestedModel, effectiveExcludedIDs, requireCompact, requiredCapability, useUpstreamTokenCost)
			if err != nil {
				return nil, decision, err
		placeholder
			if selection == nil || selection.Account == nil {
				return selection, decision, nil
		placeholder
			if s.isOpenAIAccountTransportCompatible(selection.Account, requiredTransport) &&
				accountSupportsOpenAICapabilities(selection.Account, requiredCapability, requiredImageCapability) {
				return selection, decision, nil
		placeholder
			if selection.ReleaseFunc != nil {
				selection.ReleaseFunc()
		placeholder
			if effectiveExcludedIDs == nil {
				effectiveExcludedIDs = make(map[int64]struct{placeholder)
		placeholder
			if _, exists := effectiveExcludedIDs[selection.Account.ID]; exists {
				return nil, decision, ErrNoAvailableAccounts
		placeholder
			effectiveExcludedIDs[selection.Account.ID] = struct{placeholder{placeholder
	placeholder
placeholder

	if s.checkChannelPricingRestriction(ctx, groupID, requestedModel) {
		slog.Warn("channel pricing restriction blocked request",
			"group_id", derefGroupID(groupID),
			"model", requestedModel)
		return nil, decision, fmt.Errorf("%w supporting model: %s (channel pricing restriction)", ErrNoAvailableAccounts, requestedModel)
placeholder

	var stickyAccountID int64
	if sessionHash != "" && s.cache != nil {
		if accountID, err := s.getStickySessionAccountID(ctx, groupID, sessionHash); err == nil && accountID > 0 {
			stickyAccountID = accountID
	placeholder
placeholder
	stickyWeighted := s.isOpenAIAdvancedSchedulerStickyWeightedEnabled(ctx)
	subscriptionPriority := s.isOpenAIAdvancedSchedulerSubscriptionPriorityEnabled(ctx)
	stickyPreviousAccountID := int64(0)
	if stickyWeighted && previousResponseCanMove && strings.TrimSpace(previousResponseID) != "" && platform == PlatformOpenAI {
		stickyPreviousAccountID = s.ResolveAccountIDByPreviousResponseIDForScheduler(ctx, groupID, previousResponseID, requestedModel, excludedIDs, requiredCapability, requireCompact)
placeholder

	return scheduler.Select(ctx, OpenAIAccountScheduleRequest{
		GroupID:                 groupID,
		Platform:                platform,
		SessionHash:             sessionHash,
		StickyAccountID:         stickyAccountID,
		StickyPreviousAccountID: stickyPreviousAccountID,
		StickyWeighted:          stickyWeighted,
		SubscriptionPriority:    subscriptionPriority,
		PreviousResponseID:      previousResponseID,
		PreviousResponseCanMove: previousResponseCanMove,
		UseUpstreamTokenCost:    useUpstreamTokenCost,
		RequestedModel:          requestedModel,
		RequiredTransport:       requiredTransport,
		RequiredCapability:      requiredCapability,
		RequiredImageCapability: requiredImageCapability,
		RequireCompact:          requireCompact,
		ExcludedIDs:             excludedIDs,
placeholder)
placeholder

func accountSupportsOpenAICapabilities(account *Account, requiredCapability OpenAIEndpointCapability, requiredImageCapability OpenAIImagesCapability) bool {
	if account == nil {
		return false
placeholder
	return account.SupportsOpenAIEndpointCapability(requiredCapability) &&
		account.SupportsOpenAIImageCapability(requiredImageCapability)
placeholder

func cloneExcludedAccountIDs(excludedIDs map[int64]struct{placeholder) map[int64]struct{placeholder {
	if len(excludedIDs) == 0 {
		return nil
placeholder
	cloned := make(map[int64]struct{placeholder, len(excludedIDs))
	for id := range excludedIDs {
		cloned[id] = struct{placeholder{placeholder
placeholder
	return cloned
placeholder

func (s *OpenAIGatewayService) isOpenAIAccountTransportCompatible(account *Account, requiredTransport OpenAIUpstreamTransport) bool {
	if requiredTransport == OpenAIUpstreamTransportAny || requiredTransport == OpenAIUpstreamTransportHTTPSSE {
		return true
placeholder
	if s == nil || account == nil {
		return false
placeholder
	if requiredTransport == OpenAIUpstreamTransportResponsesWebsocketV2Ingress {
		if s.cfg == nil || !s.cfg.Gateway.OpenAIWS.ModeRouterV2Enabled {
			return s.getOpenAIWSProtocolResolver().Resolve(account).Transport == OpenAIUpstreamTransportResponsesWebsocketV2
	placeholder
		mode := account.ResolveOpenAIResponsesWebSocketV2Mode(s.cfg.Gateway.OpenAIWS.IngressModeDefault)
		switch mode {
		case OpenAIWSIngressModeCtxPool, OpenAIWSIngressModePassthrough, OpenAIWSIngressModeHTTPBridge, OpenAIWSIngressModeShared, OpenAIWSIngressModeDedicated:
			return true
		default:
			return false
	placeholder
placeholder
	return s.getOpenAIWSProtocolResolver().Resolve(account).Transport == requiredTransport
placeholder

func (s *OpenAIGatewayService) ReportOpenAIAccountScheduleResult(accountID int64, model string, success bool, firstTokenMs *int) {
	if success {
		s.clearOpenAIAccountModelTransientState(accountID, normalizeOpenAIAccountModelTransientModel(model))
placeholder
	scheduler := s.getOpenAIAccountScheduler(context.Background())
	if scheduler == nil {
		return
placeholder
	scheduler.ReportResult(accountID, success, firstTokenMs)
placeholder

func (s *OpenAIGatewayService) RecordOpenAIAccountSwitch() {
	scheduler := s.getOpenAIAccountScheduler(context.Background())
	if scheduler == nil {
		return
placeholder
	scheduler.ReportSwitch()
placeholder

func (s *OpenAIGatewayService) SnapshotOpenAIAccountSchedulerMetrics() OpenAIAccountSchedulerMetricsSnapshot {
	scheduler := s.getOpenAIAccountScheduler(context.Background())
	if scheduler == nil {
		return OpenAIAccountSchedulerMetricsSnapshot{placeholder
placeholder
	return scheduler.SnapshotMetrics()
placeholder

func (s *OpenAIGatewayService) openAIWSSessionStickyTTL() time.Duration {
	if s != nil && s.cfg != nil && s.cfg.Gateway.OpenAIWS.StickySessionTTLSeconds > 0 {
		return time.Duration(s.cfg.Gateway.OpenAIWS.StickySessionTTLSeconds) * time.Second
placeholder
	return openaiStickySessionTTL
placeholder

func (s *OpenAIGatewayService) openAIWSLBTopK() int {
	if s != nil && s.cfg != nil && s.cfg.Gateway.OpenAIWS.LBTopK > 0 {
		return s.cfg.Gateway.OpenAIWS.LBTopK
placeholder
	return 7
placeholder

func (s *OpenAIGatewayService) openAIWSLBTopKForRequest(ctx context.Context) int {
	base := s.openAIWSLBTopK()
	settings := s.openAIAdvancedSchedulerRuntimeSettings(ctx)
	// DB 覆盖值与 stickyWeighted/subscriptionPriority 一样受总开关门控：
	// 关闭高级调度器后所有调用方（含管理页分数快照）都应回到配置/默认行为。
	if !settings.enabled {
		return base
placeholder
	if settings.lbTopKOverride > 0 {
		return settings.lbTopKOverride
placeholder
	return base
placeholder

func (s *OpenAIGatewayService) openAIStickyEscapeConfig() openAIStickyEscapeConfig {
	if s != nil && s.cfg != nil {
		cfg := s.cfg.Gateway.OpenAIScheduler
		enabled := cfg.StickyEscapeEnabled
		if !enabled && cfg.StickyEscapeTTFTMs == 0 && cfg.StickyEscapeErrorRate == 0 {
			enabled = true
	placeholder
		ttftMs := float64(cfg.StickyEscapeTTFTMs)
		if ttftMs <= 0 {
			ttftMs = 15000
	placeholder
		errorRate := cfg.StickyEscapeErrorRate
		if errorRate < 0 || errorRate > 1 {
			errorRate = 0.5
	placeholder
		if errorRate == 0 && cfg.StickyEscapeTTFTMs == 0 && cfg.StickyEscapeErrorRate == 0 {
			errorRate = 0.5
	placeholder
		return openAIStickyEscapeConfig{
			enabled:   enabled,
			ttftMs:    ttftMs,
			errorRate: errorRate,
	placeholder
placeholder
	return openAIStickyEscapeConfig{
		enabled:   true,
		ttftMs:    15000,
		errorRate: 0.5,
placeholder
placeholder

func (s *OpenAIGatewayService) openAIWSSchedulerWeights() GatewayOpenAIWSSchedulerScoreWeightsView {
	if s != nil && s.cfg != nil {
		return GatewayOpenAIWSSchedulerScoreWeightsView{
			Priority:      s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Priority,
			Load:          s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Load,
			Queue:         s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Queue,
			ErrorRate:     s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.ErrorRate,
			TTFT:          s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.TTFT,
			Reset:         s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.Reset,
			QuotaHeadroom: s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.QuotaHeadroom,
			UpstreamCost:  s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.UpstreamCost,
			Previous:      s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.PreviousResponse,
			SessionSticky: s.cfg.Gateway.OpenAIWS.SchedulerScoreWeights.SessionSticky,
	placeholder
placeholder
	return GatewayOpenAIWSSchedulerScoreWeightsView{
		Priority:      1.0,
		Load:          1.0,
		Queue:         0.7,
		ErrorRate:     0.8,
		TTFT:          0.5,
		Reset:         0.0,
		QuotaHeadroom: 0.0,
		UpstreamCost:  0.0,
		Previous:      5.0,
		SessionSticky: 3.0,
placeholder
placeholder

func (s *OpenAIGatewayService) openAIWSSchedulerWeightsForRequest(ctx context.Context) GatewayOpenAIWSSchedulerScoreWeightsView {
	weights := s.openAIWSSchedulerWeights()
	settings := s.openAIAdvancedSchedulerRuntimeSettings(ctx)
	// 同 openAIWSLBTopKForRequest：总开关关闭时不应用 DB 覆盖值。
	if !settings.enabled {
		return weights
placeholder
	overridden := applyOpenAIAdvancedSchedulerWeightOverrides(weights, settings.weightOverrides)
	if !overridden.configWeights().IsValid() {
		return weights
placeholder
	return overridden
placeholder

func applyOpenAIAdvancedSchedulerWeightOverrides(
	weights GatewayOpenAIWSSchedulerScoreWeightsView,
	overrides map[string]float64,
) GatewayOpenAIWSSchedulerScoreWeightsView {
	for key, value := range overrides {
		switch key {
		case "priority":
			weights.Priority = value
		case "load":
			weights.Load = value
		case "queue":
			weights.Queue = value
		case "error_rate":
			weights.ErrorRate = value
		case "ttft":
			weights.TTFT = value
		case "reset":
			weights.Reset = value
		case "quota_headroom":
			weights.QuotaHeadroom = value
		case "upstream_cost":
			weights.UpstreamCost = value
		case "previous_response":
			weights.Previous = value
		case "session_sticky":
			weights.SessionSticky = value
	placeholder
placeholder
	return weights
placeholder

type GatewayOpenAIWSSchedulerScoreWeightsView struct {
	Priority  float64
	Load      float64
	Queue     float64
	ErrorRate float64
	TTFT      float64
	// Reset 倾向「会话窗口最早重置」的账号；0 表示关闭（默认）。
	Reset         float64
	QuotaHeadroom float64
	UpstreamCost  float64
	Previous      float64
	SessionSticky float64
placeholder

func (w GatewayOpenAIWSSchedulerScoreWeightsView) configWeights() config.GatewayOpenAIWSSchedulerScoreWeights {
	return config.GatewayOpenAIWSSchedulerScoreWeights{
		Priority:         w.Priority,
		Load:             w.Load,
		Queue:            w.Queue,
		ErrorRate:        w.ErrorRate,
		TTFT:             w.TTFT,
		Reset:            w.Reset,
		QuotaHeadroom:    w.QuotaHeadroom,
		UpstreamCost:     w.UpstreamCost,
		PreviousResponse: w.Previous,
		SessionSticky:    w.SessionSticky,
placeholder
placeholder

type OpenAIAccountSchedulerScoreSnapshot struct {
	BaseScore             float64
	StickyScore           float64
	StickyScoreInfinity   bool
	StickyWeightedEnabled bool
placeholder

func (s *RateLimitService) BuildOpenAIAccountSchedulerScoreSnapshot(
	ctx context.Context,
	accounts []*Account,
	loadMap map[int64]*AccountLoadInfo,
) map[int64]OpenAIAccountSchedulerScoreSnapshot {
	gateway := &OpenAIGatewayService{cfg: nil, rateLimitService: splaceholder
	if s != nil {
		gateway.cfg = s.cfg
placeholder
	return buildOpenAIAccountSchedulerScoreSnapshot(
		accounts,
		loadMap,
		gateway.openAIWSSchedulerWeightsForRequest(ctx),
		gateway.isOpenAIAdvancedSchedulerStickyWeightedEnabled(ctx),
		gateway.openAIOAuthSchedulingRateMultiplier(ctx),
	)
placeholder

func BuildOpenAIAccountSchedulerScoreSnapshot(
	accounts []*Account,
	loadMap map[int64]*AccountLoadInfo,
) map[int64]OpenAIAccountSchedulerScoreSnapshot {
	gateway := &OpenAIGatewayService{placeholder
	return buildOpenAIAccountSchedulerScoreSnapshot(accounts, loadMap, gateway.openAIWSSchedulerWeights(), false, defaultOpenAIOAuthSchedulingRateMultiplier)
placeholder

func buildOpenAIAccountSchedulerScoreSnapshot(
	accounts []*Account,
	loadMap map[int64]*AccountLoadInfo,
	weights GatewayOpenAIWSSchedulerScoreWeightsView,
	stickyWeightedEnabled bool,
	oauthSchedulingRateMultiplier float64,
) map[int64]OpenAIAccountSchedulerScoreSnapshot {
	if len(accounts) == 0 {
		return nil
placeholder
	candidates := make([]openAIAccountCandidateScore, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
	placeholder
		loadInfo := loadMap[account.ID]
		if loadInfo == nil {
			loadInfo = &AccountLoadInfo{AccountID: account.IDplaceholder
	placeholder
		candidates = append(candidates, openAIAccountCandidateScore{
			account:   account,
			loadInfo:  loadInfo,
			errorRate: 0,
			ttft:      0,
			hasTTFT:   false,
	placeholder)
placeholder
	if len(candidates) == 0 {
		return nil
placeholder

	minPriority, maxPriority := openAIAccountSchedulingPriority(candidates[0].account), openAIAccountSchedulingPriority(candidates[0].account)
	maxWaiting := 1
	for i := range candidates {
		candidate := &candidates[i]
		candidate.priority = openAIAccountSchedulingPriority(candidate.account)
		if candidate.priority < minPriority {
			minPriority = candidate.priority
	placeholder
		if candidate.priority > maxPriority {
			maxPriority = candidate.priority
	placeholder
		if candidate.loadInfo.WaitingCount > maxWaiting {
			maxWaiting = candidate.loadInfo.WaitingCount
	placeholder
placeholder

	minResetRemaining, maxResetRemaining := 0.0, 0.0
	hasResetSample := false
	now := time.Now()
	upstreamCostFactors := map[int64]float64(nil)
	if weights.UpstreamCost > 0 {
		accounts := make([]*Account, 0, len(candidates))
		for _, candidate := range candidates {
			accounts = append(accounts, candidate.account)
	placeholder
		upstreamCostFactors = openAIUpstreamCostFactors(accounts, now, oauthSchedulingRateMultiplier)
placeholder
	if weights.Reset > 0 {
		for _, candidate := range candidates {
			end := candidate.account.SessionWindowEnd
			if end == nil || !now.Before(*end) {
				continue
		placeholder
			remaining := end.Sub(now).Seconds()
			if !hasResetSample {
				minResetRemaining, maxResetRemaining = remaining, remaining
				hasResetSample = true
				continue
		placeholder
			if remaining < minResetRemaining {
				minResetRemaining = remaining
		placeholder
			if remaining > maxResetRemaining {
				maxResetRemaining = remaining
		placeholder
	placeholder
placeholder

	result := make(map[int64]OpenAIAccountSchedulerScoreSnapshot, len(candidates))
	for _, candidate := range candidates {
		priorityFactor := 1.0
		if maxPriority > minPriority {
			priorityFactor = 1 - float64(candidate.priority-minPriority)/float64(maxPriority-minPriority)
	placeholder
		loadFactor := 1 - clamp01(float64(candidate.loadInfo.LoadRate)/100.0)
		queueFactor := 1 - clamp01(float64(candidate.loadInfo.WaitingCount)/float64(maxWaiting))
		errorFactor := 1.0
		ttftFactor := 0.5
		resetFactor := 0.0
		if weights.Reset > 0 && hasResetSample {
			if end := candidate.account.SessionWindowEnd; end != nil && now.Before(*end) {
				if maxResetRemaining > minResetRemaining {
					resetFactor = 1 - clamp01((end.Sub(now).Seconds()-minResetRemaining)/(maxResetRemaining-minResetRemaining))
			placeholder else {
					resetFactor = 1
			placeholder
		placeholder
	placeholder
		quotaHeadroomFactor := 0.0
		if weights.QuotaHeadroom > 0 {
			quotaHeadroomFactor = openAIQuotaHeadroomFactor(candidate.account, now)
	placeholder
		upstreamCostFactor := openAIUpstreamCostNeutralFactor
		if factor, ok := upstreamCostFactors[candidate.account.ID]; ok {
			upstreamCostFactor = factor
	placeholder
		baseScore := weights.Priority*priorityFactor +
			weights.Load*loadFactor +
			weights.Queue*queueFactor +
			weights.ErrorRate*errorFactor +
			weights.TTFT*ttftFactor +
			weights.Reset*resetFactor +
			weights.QuotaHeadroom*quotaHeadroomFactor +
			weights.UpstreamCost*(upstreamCostFactor-openAIUpstreamCostNeutralFactor)
		score := OpenAIAccountSchedulerScoreSnapshot{
			BaseScore:             baseScore,
			StickyWeightedEnabled: stickyWeightedEnabled,
			StickyScoreInfinity:   !stickyWeightedEnabled,
	placeholder
		if stickyWeightedEnabled {
			score.StickyScore = baseScore + weights.Previous + weights.SessionSticky
	placeholder
		result[candidate.account.ID] = score
placeholder
	return result
placeholder

func openAIUpstreamCostFactors(accounts []*Account, now time.Time, oauthSchedulingRateMultiplier float64) map[int64]float64 {
	type rateSample struct {
		accountID int64
		rate      float64
placeholder

	factors := make(map[int64]float64, len(accounts))
	samples := make([]rateSample, 0, len(accounts))
	eligibleCount := 0
	for _, account := range accounts {
		if account == nil {
			continue
	placeholder
		factors[account.ID] = openAIUpstreamCostNeutralFactor
		if !account.IsOpenAIApiKey() && !account.IsOpenAIOAuth() {
			continue
	placeholder
		eligibleCount++
		if rate, ok := openAISchedulingRate(account, now, oauthSchedulingRateMultiplier); ok {
			samples = append(samples, rateSample{accountID: account.ID, rate: rateplaceholder)
	placeholder
placeholder
	if len(samples) < 2 || eligibleCount == 0 {
		return factors
placeholder

	allEqual := true
	positiveLogs := make([]float64, 0, len(samples))
	for i, sample := range samples {
		if i > 0 && sample.rate != samples[0].rate {
			allEqual = false
	placeholder
		if sample.rate > 0 {
			positiveLogs = append(positiveLogs, math.Log(sample.rate))
	placeholder
placeholder
	if allEqual || len(positiveLogs) == 0 {
		return factors
placeholder

	sort.Float64s(positiveLogs)
	middle := len(positiveLogs) / 2
	medianLog := positiveLogs[middle]
	if len(positiveLogs)%2 == 0 {
		medianLog = (positiveLogs[middle-1] + positiveLogs[middle]) / 2
placeholder
	center := math.Exp(medianLog)
	if center <= 0 || math.IsNaN(center) || math.IsInf(center, 0) {
		return factors
placeholder

	coverage := float64(len(samples)) / float64(eligibleCount)
	for _, sample := range samples {
		rawFactor := 1.0
		if sample.rate > 0 {
			rawFactor = 1 / (1 + sample.rate/center)
	placeholder
		factors[sample.accountID] = clamp01(openAIUpstreamCostNeutralFactor + coverage*(rawFactor-openAIUpstreamCostNeutralFactor))
placeholder
	return factors
placeholder

type openAILegacyUpstreamRateOrder struct {
	enabled bool
	rates   map[int64]float64
placeholder

func newOpenAILegacyUpstreamRateOrder(accounts []*Account, now time.Time, oauthSchedulingRateMultiplier float64) openAILegacyUpstreamRateOrder {
	rates := make(map[int64]float64, len(accounts))
	var first float64
	distinct := false
	for _, account := range accounts {
		if account == nil {
			continue
	placeholder
		// 与 openAIUpstreamCostFactors 使用同一道平台门控：只有 OpenAI 平台账号
		// 的倍率参与 legacy 低倍率优先排序。上游自报倍率来自中转方，不能让它对
		// 其他平台的调度产生影响——否则自报低价即可吸走流量，而实际结算走本地倍率。
		if !account.IsOpenAIApiKey() && !account.IsOpenAIOAuth() {
			continue
	placeholder
		rate, ok := openAISchedulingRate(account, now, oauthSchedulingRateMultiplier)
		if !ok {
			continue
	placeholder
		if len(rates) == 0 {
			first = rate
	placeholder else if rate != first {
			distinct = true
	placeholder
		rates[account.ID] = rate
placeholder
	return openAILegacyUpstreamRateOrder{enabled: len(rates) >= 2 && distinct, rates: ratesplaceholder
placeholder

func openAISchedulingRate(account *Account, now time.Time, oauthSchedulingRateMultiplier float64) (float64, bool) {
	if account != nil && account.IsOpenAIOAuth() {
		return oauthSchedulingRateMultiplier, true
placeholder
	return openAIFreshUpstreamBillingRate(account, now)
placeholder

// compare returns -1 when a should be selected before b, 1 when b should be
// selected first, and 0 when the rate signal does not distinguish them.
func (o openAILegacyUpstreamRateOrder) compare(a, b *Account) int {
	if !o.enabled || a == nil || b == nil {
		return 0
placeholder
	aRate, aKnown := o.rates[a.ID]
	bRate, bKnown := o.rates[b.ID]
	if aKnown != bKnown {
		if aKnown {
			return -1
	placeholder
		return 1
placeholder
	if !aKnown || aRate == bRate {
		return 0
placeholder
	if aRate < bRate {
		return -1
placeholder
	return 1
placeholder

func openAIFreshUpstreamBillingRate(account *Account, now time.Time) (float64, bool) {
	if !isUpstreamBillingProbeAccount(account) {
		return 0, false
placeholder
	snapshot := decodeUpstreamBillingProbeSnapshot(account.Extra)
	if snapshot == nil || (snapshot.Status != UpstreamBillingProbeStatusOK && snapshot.Status != UpstreamBillingProbeStatusFailed) ||
		snapshot.ReceivedAt == nil || snapshot.ReceivedAt.IsZero() {
		return 0, false
placeholder
	receivedAt := *snapshot.ReceivedAt
	freshUntil := snapshot.FreshUntil
	if freshUntil == nil && snapshot.Status == UpstreamBillingProbeStatusOK {
		interval := snapshot.NextProbeAt.Sub(receivedAt)
		if interval > 0 {
			freshUntil = probeTimePtr(receivedAt.Add(2 * interval))
	placeholder
placeholder
	if freshUntil == nil || !freshUntil.After(receivedAt) || now.Before(receivedAt) || now.After(*freshUntil) {
		return 0, false
placeholder
	return upstreamBillingRateAt(snapshot.Data, now)
placeholder

func openAIQuotaHeadroomFactor(account *Account, now time.Time) float64 {
	if account == nil || len(account.Extra) == 0 || openAIQuotaHeadroomSnapshotStale(account.Extra, now) {
		return openAIQuotaHeadroomNeutralFactor
placeholder
	primaryUsedPercent, ok := resolveAccountExtraNumber(account.Extra, "codex_primary_used_percent", "codex_7d_used_percent")
	if !ok || openAIQuotaWindowResetAny(account.Extra, now, "primary", "7d") {
		return openAIQuotaHeadroomNeutralFactor
placeholder

	factor := 1 - clamp01(primaryUsedPercent/100)
	if secondaryUsedPercent, ok := resolveAccountExtraNumber(account.Extra, "codex_secondary_used_percent", "codex_5h_used_percent"); ok &&
		!openAIQuotaWindowResetAny(account.Extra, now, "secondary", "5h") {
		secondaryRemaining := 1 - clamp01(secondaryUsedPercent/100)
		if secondaryRemaining < openAIQuotaHeadroomSecondaryLowRemain {
			factor *= openAIQuotaHeadroomNeutralFactor
	placeholder
placeholder
	return factor
placeholder

func openAIQuotaHeadroomSnapshotStale(extra map[string]any, now time.Time) bool {
	updatedRaw, ok := extra["codex_usage_updated_at"]
	if !ok {
		return true
placeholder
	updatedAt, err := parseTime(fmt.Sprint(updatedRaw))
	if err != nil {
		return true
placeholder
	return now.Sub(updatedAt) >= openAIQuotaHeadroomSnapshotStaleAfter
placeholder

func openAIQuotaWindowResetAny(extra map[string]any, now time.Time, windows ...string) bool {
	for _, window := range windows {
		if openAIQuotaWindowReset(extra, window, now) {
			return true
	placeholder
placeholder
	return false
placeholder

func clamp01(value float64) float64 {
	switch {
	case value < 0:
		return 0
	case value > 1:
		return 1
	default:
		return value
placeholder
placeholder

func calcLoadSkewByMoments(sum float64, sumSquares float64, count int) float64 {
	if count <= 1 {
		return 0
placeholder
	mean := sum / float64(count)
	variance := sumSquares/float64(count) - mean*mean
	if variance < 0 {
		variance = 0
placeholder
	return math.Sqrt(variance)
placeholder
