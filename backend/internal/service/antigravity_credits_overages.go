package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

const antigravityCreditsOveragesKey = "antigravity_credits_overages"

type antigravity429Category string

const (
	antigravity429Unknown        antigravity429Category = "unknown"
	antigravity429RateLimited    antigravity429Category = "rate_limited"
	antigravity429QuotaExhausted antigravity429Category = "quota_exhausted"
)

var (
	creditsExhaustedCache sync.Map

	antigravityQuotaExhaustedKeywords = []string{
		"quota_exhausted",
		"quota exhausted",
placeholder

	creditsExhaustedKeywords = []string{
		"google_one_ai",
		"insufficient credit",
		"insufficient credits",
		"not enough credit",
		"not enough credits",
		"credit exhausted",
		"credits exhausted",
		"credit balance",
		"minimumcreditamountforusage",
		"minimum credit amount for usage",
		"minimum credit",
placeholder
)

// isCreditsExhausted 检查账号的 AI Credits 是否已被标记为耗尽。
func isCreditsExhausted(accountID int64) bool {
	v, ok := creditsExhaustedCache.Load(accountID)
	if !ok {
		return false
placeholder
	until, ok := v.(time.Time)
	if !ok || time.Now().After(until) {
		creditsExhaustedCache.Delete(accountID)
		return false
placeholder
	return true
placeholder

// setCreditsExhausted 将账号标记为 AI Credits 已耗尽，直到指定时间。
func setCreditsExhausted(accountID int64, until time.Time) {
	creditsExhaustedCache.Store(accountID, until)
placeholder

// clearCreditsExhausted 清除账号的 AI Credits 耗尽标记。
func clearCreditsExhausted(accountID int64) {
	creditsExhaustedCache.Delete(accountID)
placeholder

// classifyAntigravity429 将 Antigravity 的 429 响应归类为配额耗尽、限流或未知。
func classifyAntigravity429(body []byte) antigravity429Category {
	if len(body) == 0 {
		return antigravity429Unknown
placeholder
	lowerBody := strings.ToLower(string(body))
	for _, keyword := range antigravityQuotaExhaustedKeywords {
		if strings.Contains(lowerBody, keyword) {
			return antigravity429QuotaExhausted
	placeholder
placeholder
	if info := parseAntigravitySmartRetryInfo(body); info != nil && !info.IsModelCapacityExhausted {
		return antigravity429RateLimited
placeholder
	return antigravity429Unknown
placeholder

// injectEnabledCreditTypes 在已序列化的 v1internal JSON body 中注入 AI Credits 类型。
func injectEnabledCreditTypes(body []byte) []byte {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil
placeholder
	payload["enabledCreditTypes"] = []string{"GOOGLE_ONE_AI"placeholder
	result, err := json.Marshal(payload)
	if err != nil {
		return nil
placeholder
	return result
placeholder

// resolveCreditsOveragesModelKey 解析当前请求对应的 overages 状态模型 key。
func resolveCreditsOveragesModelKey(ctx context.Context, account *Account, upstreamModelName, requestedModel string) string {
	modelKey := strings.TrimSpace(upstreamModelName)
	if modelKey != "" {
		return modelKey
placeholder
	if account == nil {
		return ""
placeholder
	modelKey = resolveFinalAntigravityModelKey(ctx, account, requestedModel)
	if strings.TrimSpace(modelKey) != "" {
		return modelKey
placeholder
	return resolveAntigravityModelKey(requestedModel)
placeholder

// canUseAntigravityCreditsOverages 判断当前请求是否应直接走已激活的 overages 链路。
func canUseAntigravityCreditsOverages(ctx context.Context, account *Account, requestedModel string) bool {
	if account == nil || account.Platform != PlatformAntigravity {
		return false
placeholder
	if !account.IsSchedulable() {
		return false
placeholder
	if !account.IsOveragesEnabled() || isCreditsExhausted(account.ID) {
		return false
placeholder
	return account.getAntigravityCreditsOveragesRemainingTimeWithContext(ctx, requestedModel) > 0
placeholder

func (a *Account) getAntigravityCreditsOveragesRemainingTimeWithContext(ctx context.Context, requestedModel string) time.Duration {
	if a == nil || a.Extra == nil {
		return 0
placeholder
	modelKey := resolveFinalAntigravityModelKey(ctx, a, requestedModel)
	modelKey = strings.TrimSpace(modelKey)
	if modelKey == "" {
		return 0
placeholder
	rawStates, ok := a.Extra[antigravityCreditsOveragesKey].(map[string]any)
	if !ok {
		return 0
placeholder
	rawState, ok := rawStates[modelKey].(map[string]any)
	if !ok {
		return 0
placeholder
	activeUntilRaw, ok := rawState["active_until"].(string)
	if !ok || strings.TrimSpace(activeUntilRaw) == "" {
		return 0
placeholder
	activeUntil, err := time.Parse(time.RFC3339, activeUntilRaw)
	if err != nil {
		return 0
placeholder
	remaining := time.Until(activeUntil)
	if remaining > 0 {
		return remaining
placeholder
	return 0
placeholder

func setAntigravityCreditsOveragesActive(ctx context.Context, repo AccountRepository, account *Account, modelKey string, activeUntil time.Time) {
	if repo == nil || account == nil || account.ID == 0 || strings.TrimSpace(modelKey) == "" {
		return
placeholder
	stateMap := copyAntigravityCreditsOveragesState(account)
	stateMap[modelKey] = map[string]any{
		"activated_at": time.Now().UTC().Format(time.RFC3339),
		"active_until": activeUntil.UTC().Format(time.RFC3339),
placeholder
	ensureAccountExtra(account)
	account.Extra[antigravityCreditsOveragesKey] = stateMap
	if err := repo.UpdateExtra(ctx, account.ID, map[string]any{antigravityCreditsOveragesKey: stateMapplaceholder); err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "set overages state failed: account=%d model=%s err=%v", account.ID, modelKey, err)
placeholder
placeholder

func clearAntigravityCreditsOveragesState(ctx context.Context, repo AccountRepository, accountID int64) error {
	if repo == nil || accountID == 0 {
		return nil
placeholder
	return repo.UpdateExtra(ctx, accountID, map[string]any{
		antigravityCreditsOveragesKey: map[string]any{placeholder,
placeholder)
placeholder

func clearAntigravityCreditsOveragesStateForModel(ctx context.Context, repo AccountRepository, account *Account, modelKey string) {
	if repo == nil || account == nil || account.ID == 0 {
		return
placeholder
	stateMap := copyAntigravityCreditsOveragesState(account)
	delete(stateMap, modelKey)
	ensureAccountExtra(account)
	account.Extra[antigravityCreditsOveragesKey] = stateMap
	if err := repo.UpdateExtra(ctx, account.ID, map[string]any{antigravityCreditsOveragesKey: stateMapplaceholder); err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "clear overages state failed: account=%d model=%s err=%v", account.ID, modelKey, err)
placeholder
placeholder

func copyAntigravityCreditsOveragesState(account *Account) map[string]any {
	result := make(map[string]any)
	if account == nil || account.Extra == nil {
		return result
placeholder
	rawState, ok := account.Extra[antigravityCreditsOveragesKey].(map[string]any)
	if !ok {
		return result
placeholder
	for key, value := range rawState {
		result[key] = value
placeholder
	return result
placeholder

func ensureAccountExtra(account *Account) {
	if account != nil && account.Extra == nil {
		account.Extra = make(map[string]any)
placeholder
placeholder

// shouldMarkCreditsExhausted 判断一次 credits 请求失败是否应标记为 credits 耗尽。
func shouldMarkCreditsExhausted(resp *http.Response, respBody []byte, reqErr error) bool {
	if reqErr != nil || resp == nil {
		return false
placeholder
	if resp.StatusCode >= 500 || resp.StatusCode == http.StatusRequestTimeout {
		return false
placeholder
	if isURLLevelRateLimit(respBody) {
		return false
placeholder
	if info := parseAntigravitySmartRetryInfo(respBody); info != nil {
		return false
placeholder
	bodyLower := strings.ToLower(string(respBody))
	for _, keyword := range creditsExhaustedKeywords {
		if strings.Contains(bodyLower, keyword) {
			return true
	placeholder
placeholder
	return false
placeholder

type creditsOveragesRetryResult struct {
	handled bool
	resp    *http.Response
placeholder

// attemptCreditsOveragesRetry 在确认免费配额耗尽后，尝试注入 AI Credits 继续请求。
func (s *AntigravityGatewayService) attemptCreditsOveragesRetry(
	p antigravityRetryLoopParams,
	baseURL string,
	modelName string,
	waitDuration time.Duration,
	originalStatusCode int,
	respBody []byte,
) *creditsOveragesRetryResult {
	creditsBody := injectEnabledCreditTypes(p.body)
	if creditsBody == nil {
		return &creditsOveragesRetryResult{handled: falseplaceholder
placeholder
	modelKey := resolveCreditsOveragesModelKey(p.ctx, p.account, modelName, p.requestedModel)
	logger.LegacyPrintf("service.antigravity_gateway", "%s status=429 credit_overages_retry model=%s account=%d (injecting enabledCreditTypes)",
		p.prefix, modelKey, p.account.ID)

	creditsReq, err := antigravity.NewAPIRequestWithURL(p.ctx, baseURL, p.action, p.accessToken, creditsBody)
	if err != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "%s credit_overages_failed model=%s account=%d build_request_err=%v",
			p.prefix, modelKey, p.account.ID, err)
		return &creditsOveragesRetryResult{handled: trueplaceholder
placeholder

	creditsResp, err := p.httpUpstream.Do(creditsReq, p.proxyURL, p.account.ID, p.account.Concurrency)
	if err == nil && creditsResp != nil && creditsResp.StatusCode < 400 {
		clearCreditsExhausted(p.account.ID)
		activeUntil := s.resolveCreditsOveragesActiveUntil(respBody, waitDuration)
		setAntigravityCreditsOveragesActive(p.ctx, p.accountRepo, p.account, modelKey, activeUntil)
		logger.LegacyPrintf("service.antigravity_gateway", "%s status=%d credit_overages_success model=%s account=%d active_until=%s",
			p.prefix, creditsResp.StatusCode, modelKey, p.account.ID, activeUntil.UTC().Format(time.RFC3339))
		return &creditsOveragesRetryResult{handled: true, resp: creditsRespplaceholder
placeholder

	s.handleCreditsRetryFailure(p.prefix, modelKey, p.account, waitDuration, creditsResp, err)
	return &creditsOveragesRetryResult{handled: trueplaceholder
placeholder

func (s *AntigravityGatewayService) handleCreditsRetryFailure(
	prefix string,
	modelKey string,
	account *Account,
	waitDuration time.Duration,
	creditsResp *http.Response,
	reqErr error,
) {
	var creditsRespBody []byte
	creditsStatusCode := 0
	if creditsResp != nil {
		creditsStatusCode = creditsResp.StatusCode
		if creditsResp.Body != nil {
			creditsRespBody, _ = io.ReadAll(io.LimitReader(creditsResp.Body, 64<<10))
			_ = creditsResp.Body.Close()
	placeholder
placeholder

	if shouldMarkCreditsExhausted(creditsResp, creditsRespBody, reqErr) && account != nil {
		exhaustedUntil := s.resolveCreditsOveragesActiveUntil(creditsRespBody, waitDuration)
		setCreditsExhausted(account.ID, exhaustedUntil)
		clearAntigravityCreditsOveragesStateForModel(context.Background(), s.accountRepo, account, modelKey)
		logger.LegacyPrintf("service.antigravity_gateway", "%s credit_overages_failed model=%s account=%d marked_exhausted=true status=%d exhausted_until=%v body=%s",
			prefix, modelKey, account.ID, creditsStatusCode, exhaustedUntil, truncateForLog(creditsRespBody, 200))
		return
placeholder
	if account != nil {
		logger.LegacyPrintf("service.antigravity_gateway", "%s credit_overages_failed model=%s account=%d marked_exhausted=false status=%d err=%v body=%s",
			prefix, modelKey, account.ID, creditsStatusCode, reqErr, truncateForLog(creditsRespBody, 200))
placeholder
placeholder

func (s *AntigravityGatewayService) resolveCreditsOveragesActiveUntil(respBody []byte, waitDuration time.Duration) time.Time {
	resetAt := ParseGeminiRateLimitResetTime(respBody)
	defaultDur := waitDuration
	if defaultDur <= 0 {
		defaultDur = s.getDefaultRateLimitDuration()
placeholder
	return s.resolveResetTime(resetAt, defaultDur)
placeholder
