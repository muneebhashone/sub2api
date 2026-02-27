package service

import (
	"encoding/json"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/tidwall/gjson"
)

type claudeMaxCacheBillingOutcome struct {
	Simulated     bool
	ForcedCache1H bool
placeholder

func applyClaudeMaxCacheBillingPolicy(input *RecordUsageInput) claudeMaxCacheBillingOutcome {
	var out claudeMaxCacheBillingOutcome
	if !shouldApplyClaudeMaxBillingRules(input) {
		return out
placeholder

	if input == nil || input.Result == nil {
		return out
placeholder
	result := input.Result
	usage := &result.Usage
	accountID := int64(0)
	if input.Account != nil {
		accountID = input.Account.ID
placeholder

	if hasCacheCreationTokens(*usage) {
		before5m := usage.CacheCreation5mTokens
		before1h := usage.CacheCreation1hTokens
		out.ForcedCache1H = safelyForceCacheCreationTo1H(usage)
		if out.ForcedCache1H {
			logger.LegacyPrintf("service.gateway", "force_claude_max_cache_1h: model=%s account=%d cache_creation_5m:%d->%d cache_creation_1h:%d->%d",
				result.Model,
				accountID,
				before5m,
				usage.CacheCreation5mTokens,
				before1h,
				usage.CacheCreation1hTokens,
			)
	placeholder
		return out
placeholder

	if !shouldSimulateClaudeMaxUsage(input) {
		return out
placeholder
	beforeInputTokens := usage.InputTokens
	out.Simulated = safelyApplyClaudeMaxUsageSimulation(result, input.ParsedRequest)
	if out.Simulated {
		logger.LegacyPrintf("service.gateway", "simulate_claude_max_usage: model=%s account=%d input_tokens:%d->%d cache_creation_1h=%d",
			result.Model,
			accountID,
			beforeInputTokens,
			usage.InputTokens,
			usage.CacheCreation1hTokens,
		)
placeholder
	return out
placeholder

func isClaudeFamilyModel(model string) bool {
	normalized := strings.ToLower(strings.TrimSpace(claude.NormalizeModelID(model)))
	if normalized == "" {
		return false
placeholder
	return strings.Contains(normalized, "claude-")
placeholder

func shouldApplyClaudeMaxBillingRules(input *RecordUsageInput) bool {
	if input == nil || input.Result == nil || input.APIKey == nil || input.APIKey.Group == nil {
		return false
placeholder
	group := input.APIKey.Group
	if !group.SimulateClaudeMaxEnabled || group.Platform != PlatformAnthropic {
		return false
placeholder

	model := input.Result.Model
	if model == "" && input.ParsedRequest != nil {
		model = input.ParsedRequest.Model
placeholder
	if !isClaudeFamilyModel(model) {
		return false
placeholder
	return true
placeholder

func hasCacheCreationTokens(usage ClaudeUsage) bool {
	return usage.CacheCreationInputTokens > 0 || usage.CacheCreation5mTokens > 0 || usage.CacheCreation1hTokens > 0
placeholder

func shouldSimulateClaudeMaxUsage(input *RecordUsageInput) bool {
	if !shouldApplyClaudeMaxBillingRules(input) {
		return false
placeholder
	if !hasClaudeCacheSignals(input.ParsedRequest) {
		return false
placeholder
	usage := input.Result.Usage
	if usage.InputTokens <= 0 {
		return false
placeholder
	if hasCacheCreationTokens(usage) {
		return false
placeholder
	return true
placeholder

func forceCacheCreationTo1H(usage *ClaudeUsage) bool {
	if usage == nil || !hasCacheCreationTokens(*usage) {
		return false
placeholder

	before5m := usage.CacheCreation5mTokens
	before1h := usage.CacheCreation1hTokens
	beforeAgg := usage.CacheCreationInputTokens

	_ = applyCacheTTLOverride(usage, "1h")
	total := usage.CacheCreation5mTokens + usage.CacheCreation1hTokens
	if total <= 0 {
		total = usage.CacheCreationInputTokens
placeholder
	if total <= 0 {
		return false
placeholder

	usage.CacheCreation5mTokens = 0
	usage.CacheCreation1hTokens = total
	usage.CacheCreationInputTokens = total

	return before5m != usage.CacheCreation5mTokens ||
		before1h != usage.CacheCreation1hTokens ||
		beforeAgg != usage.CacheCreationInputTokens
placeholder

func safelyApplyClaudeMaxUsageSimulation(result *ForwardResult, parsed *ParsedRequest) (changed bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.LegacyPrintf("service.gateway", "simulate_claude_max_usage skipped: panic=%v", r)
			changed = false
	placeholder
placeholder()
	return applyClaudeMaxUsageSimulation(result, parsed)
placeholder

func safelyForceCacheCreationTo1H(usage *ClaudeUsage) (changed bool) {
	defer func() {
		if r := recover(); r != nil {
			logger.LegacyPrintf("service.gateway", "force_cache_creation_1h skipped: panic=%v", r)
			changed = false
	placeholder
placeholder()
	return forceCacheCreationTo1H(usage)
placeholder

func applyClaudeMaxUsageSimulation(result *ForwardResult, parsed *ParsedRequest) bool {
	if result == nil {
		return false
placeholder
	return projectUsageToClaudeMax1H(&result.Usage, parsed)
placeholder

func projectUsageToClaudeMax1H(usage *ClaudeUsage, parsed *ParsedRequest) bool {
	if usage == nil {
		return false
placeholder
	totalWindowTokens := usage.InputTokens + usage.CacheCreation5mTokens + usage.CacheCreation1hTokens
	if totalWindowTokens <= 1 {
		return false
placeholder

	simulatedInputTokens := computeClaudeMaxProjectedInputTokens(totalWindowTokens, parsed)
	if simulatedInputTokens <= 0 {
		simulatedInputTokens = 1
placeholder
	if simulatedInputTokens >= totalWindowTokens {
		simulatedInputTokens = totalWindowTokens - 1
placeholder

	cacheCreation1hTokens := totalWindowTokens - simulatedInputTokens
	if usage.InputTokens == simulatedInputTokens &&
		usage.CacheCreation5mTokens == 0 &&
		usage.CacheCreation1hTokens == cacheCreation1hTokens &&
		usage.CacheCreationInputTokens == cacheCreation1hTokens {
		return false
placeholder

	usage.InputTokens = simulatedInputTokens
	usage.CacheCreation5mTokens = 0
	usage.CacheCreation1hTokens = cacheCreation1hTokens
	usage.CacheCreationInputTokens = cacheCreation1hTokens
	return true
placeholder

type claudeCacheProjection struct {
	HasBreakpoint        bool
	BreakpointCount      int
	TotalEstimatedTokens int
	TailEstimatedTokens  int
placeholder

func computeClaudeMaxProjectedInputTokens(totalWindowTokens int, parsed *ParsedRequest) int {
	if totalWindowTokens <= 1 {
		return totalWindowTokens
placeholder

	projection := analyzeClaudeCacheProjection(parsed)
	if !projection.HasBreakpoint || projection.TotalEstimatedTokens <= 0 || projection.TailEstimatedTokens <= 0 {
		return totalWindowTokens
placeholder

	totalEstimate := int64(projection.TotalEstimatedTokens)
	tailEstimate := int64(projection.TailEstimatedTokens)
	if tailEstimate > totalEstimate {
		tailEstimate = totalEstimate
placeholder

	scaled := (int64(totalWindowTokens)*tailEstimate + totalEstimate/2) / totalEstimate
	if scaled <= 0 {
		scaled = 1
placeholder
	if scaled >= int64(totalWindowTokens) {
		scaled = int64(totalWindowTokens - 1)
placeholder
	return int(scaled)
placeholder

func hasClaudeCacheSignals(parsed *ParsedRequest) bool {
	if parsed == nil {
		return false
placeholder
	if hasTopLevelEphemeralCacheControl(parsed) {
		return true
placeholder
	return countExplicitCacheBreakpoints(parsed) > 0
placeholder

func hasTopLevelEphemeralCacheControl(parsed *ParsedRequest) bool {
	if parsed == nil || len(parsed.Body) == 0 {
		return false
placeholder
	cacheType := strings.TrimSpace(gjson.GetBytes(parsed.Body, "cache_control.type").String())
	return strings.EqualFold(cacheType, "ephemeral")
placeholder

func analyzeClaudeCacheProjection(parsed *ParsedRequest) claudeCacheProjection {
	var projection claudeCacheProjection
	if parsed == nil {
		return projection
placeholder

	total := 0
	lastBreakpointAt := -1

	switch system := parsed.System.(type) {
	case string:
		total += claudeMaxMessageOverheadTokens + estimateClaudeTextTokens(system)
	case []any:
		for _, raw := range system {
			block, ok := raw.(map[string]any)
			if !ok {
				total += claudeMaxUnknownContentTokens
				continue
		placeholder
			total += estimateClaudeBlockTokens(block)
			if hasEphemeralCacheControl(block) {
				lastBreakpointAt = total
				projection.BreakpointCount++
				projection.HasBreakpoint = true
		placeholder
	placeholder
placeholder

	for _, rawMsg := range parsed.Messages {
		total += claudeMaxMessageOverheadTokens
		msg, ok := rawMsg.(map[string]any)
		if !ok {
			total += claudeMaxUnknownContentTokens
			continue
	placeholder
		content, exists := msg["content"]
		if !exists {
			continue
	placeholder
		msgTokens, msgLastBreak, msgBreakCount := estimateClaudeContentTokens(content)
		total += msgTokens
		if msgBreakCount > 0 {
			lastBreakpointAt = total - msgTokens + msgLastBreak
			projection.BreakpointCount += msgBreakCount
			projection.HasBreakpoint = true
	placeholder
placeholder

	if total <= 0 {
		total = 1
placeholder
	projection.TotalEstimatedTokens = total

	if projection.HasBreakpoint && lastBreakpointAt >= 0 {
		tail := total - lastBreakpointAt
		if tail <= 0 {
			tail = 1
	placeholder
		projection.TailEstimatedTokens = tail
		return projection
placeholder

	if hasTopLevelEphemeralCacheControl(parsed) {
		tail := estimateLastUserMessageTokens(parsed)
		if tail <= 0 {
			tail = 1
	placeholder
		projection.HasBreakpoint = true
		projection.BreakpointCount = 1
		projection.TailEstimatedTokens = tail
placeholder
	return projection
placeholder

func countExplicitCacheBreakpoints(parsed *ParsedRequest) int {
	if parsed == nil {
		return 0
placeholder
	total := 0
	if system, ok := parsed.System.([]any); ok {
		for _, raw := range system {
			if block, ok := raw.(map[string]any); ok && hasEphemeralCacheControl(block) {
				total++
		placeholder
	placeholder
placeholder
	for _, rawMsg := range parsed.Messages {
		msg, ok := rawMsg.(map[string]any)
		if !ok {
			continue
	placeholder
		content, ok := msg["content"].([]any)
		if !ok {
			continue
	placeholder
		for _, raw := range content {
			if block, ok := raw.(map[string]any); ok && hasEphemeralCacheControl(block) {
				total++
		placeholder
	placeholder
placeholder
	return total
placeholder

func hasEphemeralCacheControl(block map[string]any) bool {
	if block == nil {
		return false
placeholder
	raw, ok := block["cache_control"]
	if !ok || raw == nil {
		return false
placeholder
	switch cc := raw.(type) {
	case map[string]any:
		cacheType, _ := cc["type"].(string)
		return strings.EqualFold(strings.TrimSpace(cacheType), "ephemeral")
	case map[string]string:
		return strings.EqualFold(strings.TrimSpace(cc["type"]), "ephemeral")
	default:
		return false
placeholder
placeholder

func estimateClaudeContentTokens(content any) (tokens int, lastBreakAt int, breakpointCount int) {
	switch value := content.(type) {
	case string:
		return estimateClaudeTextTokens(value), -1, 0
	case []any:
		total := 0
		lastBreak := -1
		breaks := 0
		for _, raw := range value {
			block, ok := raw.(map[string]any)
			if !ok {
				total += claudeMaxUnknownContentTokens
				continue
		placeholder
			total += estimateClaudeBlockTokens(block)
			if hasEphemeralCacheControl(block) {
				lastBreak = total
				breaks++
		placeholder
	placeholder
		return total, lastBreak, breaks
	default:
		return estimateStructuredTokens(value), -1, 0
placeholder
placeholder

func estimateClaudeBlockTokens(block map[string]any) int {
	if block == nil {
		return claudeMaxUnknownContentTokens
placeholder
	tokens := claudeMaxBlockOverheadTokens
	blockType, _ := block["type"].(string)
	switch blockType {
	case "text":
		if text, ok := block["text"].(string); ok {
			tokens += estimateClaudeTextTokens(text)
	placeholder
	case "tool_result":
		if content, ok := block["content"]; ok {
			nested, _, _ := estimateClaudeContentTokens(content)
			tokens += nested
	placeholder
	case "tool_use":
		if name, ok := block["name"].(string); ok {
			tokens += estimateClaudeTextTokens(name)
	placeholder
		if input, ok := block["input"]; ok {
			tokens += estimateStructuredTokens(input)
	placeholder
	default:
		if text, ok := block["text"].(string); ok {
			tokens += estimateClaudeTextTokens(text)
	placeholder else if content, ok := block["content"]; ok {
			nested, _, _ := estimateClaudeContentTokens(content)
			tokens += nested
	placeholder
placeholder
	if tokens <= claudeMaxBlockOverheadTokens {
		tokens += claudeMaxUnknownContentTokens
placeholder
	return tokens
placeholder

func estimateLastUserMessageTokens(parsed *ParsedRequest) int {
	if parsed == nil || len(parsed.Messages) == 0 {
		return 0
placeholder
	for i := len(parsed.Messages) - 1; i >= 0; i-- {
		msg, ok := parsed.Messages[i].(map[string]any)
		if !ok {
			continue
	placeholder
		role, _ := msg["role"].(string)
		if !strings.EqualFold(strings.TrimSpace(role), "user") {
			continue
	placeholder
		tokens, _, _ := estimateClaudeContentTokens(msg["content"])
		return claudeMaxMessageOverheadTokens + tokens
placeholder
	return 0
placeholder

func estimateStructuredTokens(v any) int {
	if v == nil {
		return 0
placeholder
	raw, err := json.Marshal(v)
	if err != nil {
		return claudeMaxUnknownContentTokens
placeholder
	return estimateClaudeTextTokens(string(raw))
placeholder

func estimateClaudeTextTokens(text string) int {
	if tokens, ok := estimateTokensByThirdPartyTokenizer(text); ok {
		return tokens
placeholder
	return estimateClaudeTextTokensHeuristic(text)
placeholder

func estimateClaudeTextTokensHeuristic(text string) int {
	normalized := strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
	if normalized == "" {
		return 0
placeholder
	asciiChars := 0
	nonASCIIChars := 0
	for _, r := range normalized {
		if r <= 127 {
			asciiChars++
	placeholder else {
			nonASCIIChars++
	placeholder
placeholder
	tokens := nonASCIIChars
	if asciiChars > 0 {
		tokens += (asciiChars + 3) / 4
placeholder
	if words := len(strings.Fields(normalized)); words > tokens {
		tokens = words
placeholder
	if tokens <= 0 {
		return 1
placeholder
	return tokens
placeholder
