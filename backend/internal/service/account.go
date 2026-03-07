// Package service provides business logic and domain services for the application.
package service

import (
	"encoding/json"
	"hash/fnv"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/domain"
)

type Account struct {
	ID          int64
	Name        string
	Notes       *string
	Platform    string
	Type        string
	Credentials map[string]any
	Extra       map[string]any
	ProxyID     *int64
	Concurrency int
	Priority    int
	// RateMultiplier 账号计费倍率（>=0，允许 0 表示该账号计费为 0）。
	// 使用指针用于兼容旧版本调度缓存（Redis）中缺字段的情况：nil 表示按 1.0 处理。
	RateMultiplier     *float64
	LoadFactor         *int // 调度负载因子；nil 表示使用 Concurrency
	Status             string
	ErrorMessage       string
	LastUsedAt         *time.Time
	ExpiresAt          *time.Time
	AutoPauseOnExpired bool
	CreatedAt          time.Time
	UpdatedAt          time.Time

	Schedulable bool

	RateLimitedAt    *time.Time
	RateLimitResetAt *time.Time
	OverloadUntil    *time.Time

	TempUnschedulableUntil  *time.Time
	TempUnschedulableReason string

	SessionWindowStart  *time.Time
	SessionWindowEnd    *time.Time
	SessionWindowStatus string

	Proxy         *Proxy
	AccountGroups []AccountGroup
	GroupIDs      []int64
	Groups        []*Group

	// model_mapping 热路径缓存（非持久化字段）
	modelMappingCache               map[string]string
	modelMappingCacheReady          bool
	modelMappingCacheCredentialsPtr uintptr
	modelMappingCacheRawPtr         uintptr
	modelMappingCacheRawLen         int
	modelMappingCacheRawSig         uint64
placeholder

type TempUnschedulableRule struct {
	ErrorCode       int      `json:"error_code"`
	Keywords        []string `json:"keywords"`
	DurationMinutes int      `json:"duration_minutes"`
	Description     string   `json:"description"`
placeholder

func (a *Account) IsActive() bool {
	return a.Status == StatusActive
placeholder

// BillingRateMultiplier 返回账号计费倍率。
// - nil 表示未配置/旧缓存缺字段，按 1.0 处理
// - 允许 0，表示该账号计费为 0
// - 负数属于非法数据，出于安全考虑按 1.0 处理
func (a *Account) BillingRateMultiplier() float64 {
	if a == nil || a.RateMultiplier == nil {
		return 1.0
placeholder
	if *a.RateMultiplier < 0 {
		return 1.0
placeholder
	return *a.RateMultiplier
placeholder

func (a *Account) EffectiveLoadFactor() int {
	if a == nil {
		return 1
placeholder
	if a.LoadFactor != nil && *a.LoadFactor > 0 {
		return *a.LoadFactor
placeholder
	if a.Concurrency > 0 {
		return a.Concurrency
placeholder
	return 1
placeholder

func (a *Account) IsSchedulable() bool {
	if !a.IsActive() || !a.Schedulable {
		return false
placeholder
	now := time.Now()
	if a.AutoPauseOnExpired && a.ExpiresAt != nil && !now.Before(*a.ExpiresAt) {
		return false
placeholder
	if a.OverloadUntil != nil && now.Before(*a.OverloadUntil) {
		return false
placeholder
	if a.RateLimitResetAt != nil && now.Before(*a.RateLimitResetAt) {
		return false
placeholder
	if a.TempUnschedulableUntil != nil && now.Before(*a.TempUnschedulableUntil) {
		return false
placeholder
	return true
placeholder

func (a *Account) IsRateLimited() bool {
	if a.RateLimitResetAt == nil {
		return false
placeholder
	return time.Now().Before(*a.RateLimitResetAt)
placeholder

func (a *Account) IsOverloaded() bool {
	if a.OverloadUntil == nil {
		return false
placeholder
	return time.Now().Before(*a.OverloadUntil)
placeholder

func (a *Account) IsOAuth() bool {
	return a.Type == AccountTypeOAuth || a.Type == AccountTypeSetupToken
placeholder

func (a *Account) IsGemini() bool {
	return a.Platform == PlatformGemini
placeholder

func (a *Account) GeminiOAuthType() string {
	if a.Platform != PlatformGemini || a.Type != AccountTypeOAuth {
		return ""
placeholder
	oauthType := strings.TrimSpace(a.GetCredential("oauth_type"))
	if oauthType == "" && strings.TrimSpace(a.GetCredential("project_id")) != "" {
		return "code_assist"
placeholder
	return oauthType
placeholder

func (a *Account) GeminiTierID() string {
	tierID := strings.TrimSpace(a.GetCredential("tier_id"))
	return tierID
placeholder

func (a *Account) IsGeminiCodeAssist() bool {
	if a.Platform != PlatformGemini || a.Type != AccountTypeOAuth {
		return false
placeholder
	oauthType := a.GeminiOAuthType()
	if oauthType == "" {
		return strings.TrimSpace(a.GetCredential("project_id")) != ""
placeholder
	return oauthType == "code_assist"
placeholder

func (a *Account) CanGetUsage() bool {
	return a.Type == AccountTypeOAuth
placeholder

func (a *Account) GetCredential(key string) string {
	if a.Credentials == nil {
		return ""
placeholder
	v, ok := a.Credentials[key]
	if !ok || v == nil {
		return ""
placeholder

	// 支持多种类型（兼容历史数据中 expires_at 等字段可能是数字或字符串）
	switch val := v.(type) {
	case string:
		return val
	case json.Number:
		// GORM datatypes.JSONMap 使用 UseNumber() 解析，数字类型为 json.Number
		return val.String()
	case float64:
		// JSON 解析后数字默认为 float64
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case int:
		return strconv.Itoa(val)
	default:
		return ""
placeholder
placeholder

// GetCredentialAsTime 解析凭证中的时间戳字段，支持多种格式
// 兼容以下格式：
//   - RFC3339 字符串: "2025-01-01T00:00:00Z"
//   - Unix 时间戳字符串: "1735689600"
//   - Unix 时间戳数字: 1735689600 (float64/int64/json.Number)
func (a *Account) GetCredentialAsTime(key string) *time.Time {
	s := a.GetCredential(key)
	if s == "" {
		return nil
placeholder
	// 尝试 RFC3339 格式
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t
placeholder
	// 尝试 Unix 时间戳（纯数字字符串）
	if ts, err := strconv.ParseInt(s, 10, 64); err == nil {
		t := time.Unix(ts, 0)
		return &t
placeholder
	return nil
placeholder

// GetCredentialAsInt64 解析凭证中的 int64 字段
// 用于读取 _token_version 等内部字段
func (a *Account) GetCredentialAsInt64(key string) int64 {
	if a == nil || a.Credentials == nil {
		return 0
placeholder
	val, ok := a.Credentials[key]
	if !ok || val == nil {
		return 0
placeholder
	switch v := val.(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i
	placeholder
	case string:
		if i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return i
	placeholder
placeholder
	return 0
placeholder

func (a *Account) IsTempUnschedulableEnabled() bool {
	if a.Credentials == nil {
		return false
placeholder
	raw, ok := a.Credentials["temp_unschedulable_enabled"]
	if !ok || raw == nil {
		return false
placeholder
	enabled, ok := raw.(bool)
	return ok && enabled
placeholder

func (a *Account) GetTempUnschedulableRules() []TempUnschedulableRule {
	if a.Credentials == nil {
		return nil
placeholder
	raw, ok := a.Credentials["temp_unschedulable_rules"]
	if !ok || raw == nil {
		return nil
placeholder

	arr, ok := raw.([]any)
	if !ok {
		return nil
placeholder

	rules := make([]TempUnschedulableRule, 0, len(arr))
	for _, item := range arr {
		entry, ok := item.(map[string]any)
		if !ok || entry == nil {
			continue
	placeholder

		rule := TempUnschedulableRule{
			ErrorCode:       parseTempUnschedInt(entry["error_code"]),
			Keywords:        parseTempUnschedStrings(entry["keywords"]),
			DurationMinutes: parseTempUnschedInt(entry["duration_minutes"]),
			Description:     parseTempUnschedString(entry["description"]),
	placeholder

		if rule.ErrorCode <= 0 || rule.DurationMinutes <= 0 || len(rule.Keywords) == 0 {
			continue
	placeholder

		rules = append(rules, rule)
placeholder

	return rules
placeholder

func parseTempUnschedString(value any) string {
	s, ok := value.(string)
	if !ok {
		return ""
placeholder
	return strings.TrimSpace(s)
placeholder

func parseTempUnschedStrings(value any) []string {
	if value == nil {
		return nil
placeholder

	var raw []string
	switch v := value.(type) {
	case []string:
		raw = v
	case []any:
		raw = make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				raw = append(raw, s)
		placeholder
	placeholder
	default:
		return nil
placeholder

	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s := strings.TrimSpace(item)
		if s != "" {
			out = append(out, s)
	placeholder
placeholder
	return out
placeholder

func normalizeAccountNotes(value *string) *string {
	if value == nil {
		return nil
placeholder
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
placeholder
	return &trimmed
placeholder

func parseTempUnschedInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return int(i)
	placeholder
	case string:
		if i, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return i
	placeholder
placeholder
	return 0
placeholder

func (a *Account) GetModelMapping() map[string]string {
	credentialsPtr := mapPtr(a.Credentials)
	rawMapping, _ := a.Credentials["model_mapping"].(map[string]any)
	rawPtr := mapPtr(rawMapping)
	rawLen := len(rawMapping)
	rawSig := uint64(0)
	rawSigReady := false

	if a.modelMappingCacheReady &&
		a.modelMappingCacheCredentialsPtr == credentialsPtr &&
		a.modelMappingCacheRawPtr == rawPtr &&
		a.modelMappingCacheRawLen == rawLen {
		rawSig = modelMappingSignature(rawMapping)
		rawSigReady = true
		if a.modelMappingCacheRawSig == rawSig {
			return a.modelMappingCache
	placeholder
placeholder

	mapping := a.resolveModelMapping(rawMapping)
	if !rawSigReady {
		rawSig = modelMappingSignature(rawMapping)
placeholder

	a.modelMappingCache = mapping
	a.modelMappingCacheReady = true
	a.modelMappingCacheCredentialsPtr = credentialsPtr
	a.modelMappingCacheRawPtr = rawPtr
	a.modelMappingCacheRawLen = rawLen
	a.modelMappingCacheRawSig = rawSig
	return mapping
placeholder

func (a *Account) resolveModelMapping(rawMapping map[string]any) map[string]string {
	if a.Credentials == nil {
		// Antigravity 平台使用默认映射
		if a.Platform == domain.PlatformAntigravity {
			return domain.DefaultAntigravityModelMapping
	placeholder
		return nil
placeholder
	if len(rawMapping) == 0 {
		// Antigravity 平台使用默认映射
		if a.Platform == domain.PlatformAntigravity {
			return domain.DefaultAntigravityModelMapping
	placeholder
		return nil
placeholder

	result := make(map[string]string)
	for k, v := range rawMapping {
		if s, ok := v.(string); ok {
			result[k] = s
	placeholder
placeholder
	if len(result) > 0 {
		if a.Platform == domain.PlatformAntigravity {
			ensureAntigravityDefaultPassthroughs(result, []string{
				"gemini-3-flash",
				"gemini-3.1-pro-high",
				"gemini-3.1-pro-low",
		placeholder)
	placeholder
		return result
placeholder

	// Antigravity 平台使用默认映射
	if a.Platform == domain.PlatformAntigravity {
		return domain.DefaultAntigravityModelMapping
placeholder
	return nil
placeholder

func mapPtr(m map[string]any) uintptr {
	if m == nil {
		return 0
placeholder
	return reflect.ValueOf(m).Pointer()
placeholder

func modelMappingSignature(rawMapping map[string]any) uint64 {
	if len(rawMapping) == 0 {
		return 0
placeholder
	keys := make([]string, 0, len(rawMapping))
	for k := range rawMapping {
		keys = append(keys, k)
placeholder
	sort.Strings(keys)

	h := fnv.New64a()
	for _, k := range keys {
		_, _ = h.Write([]byte(k))
		_, _ = h.Write([]byte{0placeholder)
		if v, ok := rawMapping[k].(string); ok {
			_, _ = h.Write([]byte(v))
	placeholder else {
			_, _ = h.Write([]byte{1placeholder)
	placeholder
		_, _ = h.Write([]byte{0xffplaceholder)
placeholder
	return h.Sum64()
placeholder

func ensureAntigravityDefaultPassthrough(mapping map[string]string, model string) {
	if mapping == nil || model == "" {
		return
placeholder
	if _, exists := mapping[model]; exists {
		return
placeholder
	for pattern := range mapping {
		if matchWildcard(pattern, model) {
			return
	placeholder
placeholder
	mapping[model] = model
placeholder

func ensureAntigravityDefaultPassthroughs(mapping map[string]string, models []string) {
	for _, model := range models {
		ensureAntigravityDefaultPassthrough(mapping, model)
placeholder
placeholder

// IsModelSupported 检查模型是否在 model_mapping 中（支持通配符）
// 如果未配置 mapping，返回 true（允许所有模型）
func (a *Account) IsModelSupported(requestedModel string) bool {
	mapping := a.GetModelMapping()
	if len(mapping) == 0 {
		return true // 无映射 = 允许所有
placeholder
	// 精确匹配
	if _, exists := mapping[requestedModel]; exists {
		return true
placeholder
	// 通配符匹配
	for pattern := range mapping {
		if matchWildcard(pattern, requestedModel) {
			return true
	placeholder
placeholder
	return false
placeholder

// GetMappedModel 获取映射后的模型名（支持通配符，最长优先匹配）
// 如果未配置 mapping，返回原始模型名
func (a *Account) GetMappedModel(requestedModel string) string {
	mapping := a.GetModelMapping()
	if len(mapping) == 0 {
		return requestedModel
placeholder
	// 精确匹配优先
	if mappedModel, exists := mapping[requestedModel]; exists {
		return mappedModel
placeholder
	// 通配符匹配（最长优先）
	return matchWildcardMapping(mapping, requestedModel)
placeholder

func (a *Account) GetBaseURL() string {
	if a.Type != AccountTypeAPIKey {
		return ""
placeholder
	baseURL := a.GetCredential("base_url")
	if baseURL == "" {
		return "https://api.anthropic.com"
placeholder
	if a.Platform == PlatformAntigravity {
		return strings.TrimRight(baseURL, "/") + "/antigravity"
placeholder
	return baseURL
placeholder

// GetGeminiBaseURL 返回 Gemini 兼容端点的 base URL。
// Antigravity 平台的 APIKey 账号自动拼接 /antigravity。
func (a *Account) GetGeminiBaseURL(defaultBaseURL string) string {
	baseURL := strings.TrimSpace(a.GetCredential("base_url"))
	if baseURL == "" {
		return defaultBaseURL
placeholder
	if a.Platform == PlatformAntigravity && a.Type == AccountTypeAPIKey {
		return strings.TrimRight(baseURL, "/") + "/antigravity"
placeholder
	return baseURL
placeholder

func (a *Account) GetExtraString(key string) string {
	if a.Extra == nil {
		return ""
placeholder
	if v, ok := a.Extra[key]; ok {
		if s, ok := v.(string); ok {
			return s
	placeholder
placeholder
	return ""
placeholder

func (a *Account) GetClaudeUserID() string {
	if v := strings.TrimSpace(a.GetExtraString("claude_user_id")); v != "" {
		return v
placeholder
	if v := strings.TrimSpace(a.GetExtraString("anthropic_user_id")); v != "" {
		return v
placeholder
	if v := strings.TrimSpace(a.GetCredential("claude_user_id")); v != "" {
		return v
placeholder
	if v := strings.TrimSpace(a.GetCredential("anthropic_user_id")); v != "" {
		return v
placeholder
	return ""
placeholder

// matchAntigravityWildcard 通配符匹配（仅支持末尾 *）
// 用于 model_mapping 的通配符匹配
func matchAntigravityWildcard(pattern, str string) bool {
	if strings.HasSuffix(pattern, "*") {
		prefix := pattern[:len(pattern)-1]
		return strings.HasPrefix(str, prefix)
placeholder
	return pattern == str
placeholder

// matchWildcard 通用通配符匹配（仅支持末尾 *）
// 复用 Antigravity 的通配符逻辑，供其他平台使用
func matchWildcard(pattern, str string) bool {
	return matchAntigravityWildcard(pattern, str)
placeholder

// matchWildcardMapping 通配符映射匹配（最长优先）
// 如果没有匹配，返回原始字符串
func matchWildcardMapping(mapping map[string]string, requestedModel string) string {
	// 收集所有匹配的 pattern，按长度降序排序（最长优先）
	type patternMatch struct {
		pattern string
		target  string
placeholder
	var matches []patternMatch

	for pattern, target := range mapping {
		if matchWildcard(pattern, requestedModel) {
			matches = append(matches, patternMatch{pattern, targetplaceholder)
	placeholder
placeholder

	if len(matches) == 0 {
		return requestedModel // 无匹配，返回原始模型名
placeholder

	// 按 pattern 长度降序排序
	sort.Slice(matches, func(i, j int) bool {
		if len(matches[i].pattern) != len(matches[j].pattern) {
			return len(matches[i].pattern) > len(matches[j].pattern)
	placeholder
		return matches[i].pattern < matches[j].pattern
placeholder)

	return matches[0].target
placeholder

func (a *Account) IsCustomErrorCodesEnabled() bool {
	if a.Type != AccountTypeAPIKey || a.Credentials == nil {
		return false
placeholder
	if v, ok := a.Credentials["custom_error_codes_enabled"]; ok {
		if enabled, ok := v.(bool); ok {
			return enabled
	placeholder
placeholder
	return false
placeholder

func (a *Account) GetCustomErrorCodes() []int {
	if a.Credentials == nil {
		return nil
placeholder
	raw, ok := a.Credentials["custom_error_codes"]
	if !ok || raw == nil {
		return nil
placeholder
	if arr, ok := raw.([]any); ok {
		result := make([]int, 0, len(arr))
		for _, v := range arr {
			if f, ok := v.(float64); ok {
				result = append(result, int(f))
		placeholder
	placeholder
		return result
placeholder
	return nil
placeholder

func (a *Account) ShouldHandleErrorCode(statusCode int) bool {
	if !a.IsCustomErrorCodesEnabled() {
		return true
placeholder
	codes := a.GetCustomErrorCodes()
	if len(codes) == 0 {
		return true
placeholder
	for _, code := range codes {
		if code == statusCode {
			return true
	placeholder
placeholder
	return false
placeholder

func (a *Account) IsInterceptWarmupEnabled() bool {
	if a.Credentials == nil {
		return false
placeholder
	if v, ok := a.Credentials["intercept_warmup_requests"]; ok {
		if enabled, ok := v.(bool); ok {
			return enabled
	placeholder
placeholder
	return false
placeholder

func (a *Account) IsOpenAI() bool {
	return a.Platform == PlatformOpenAI
placeholder

func (a *Account) IsAnthropic() bool {
	return a.Platform == PlatformAnthropic
placeholder

func (a *Account) IsOpenAIOAuth() bool {
	return a.IsOpenAI() && a.Type == AccountTypeOAuth
placeholder

func (a *Account) IsOpenAIApiKey() bool {
	return a.IsOpenAI() && a.Type == AccountTypeAPIKey
placeholder

func (a *Account) GetOpenAIBaseURL() string {
	if !a.IsOpenAI() {
		return ""
placeholder
	if a.Type == AccountTypeAPIKey {
		baseURL := a.GetCredential("base_url")
		if baseURL != "" {
			return baseURL
	placeholder
placeholder
	return "https://api.openai.com"
placeholder

func (a *Account) GetOpenAIAccessToken() string {
	if !a.IsOpenAI() {
		return ""
placeholder
	return a.GetCredential("access_token")
placeholder

func (a *Account) GetOpenAIRefreshToken() string {
	if !a.IsOpenAIOAuth() {
		return ""
placeholder
	return a.GetCredential("refresh_token")
placeholder

func (a *Account) GetOpenAIIDToken() string {
	if !a.IsOpenAIOAuth() {
		return ""
placeholder
	return a.GetCredential("id_token")
placeholder

func (a *Account) GetOpenAIApiKey() string {
	if !a.IsOpenAIApiKey() {
		return ""
placeholder
	return a.GetCredential("api_key")
placeholder

func (a *Account) GetOpenAIUserAgent() string {
	if !a.IsOpenAI() {
		return ""
placeholder
	return a.GetCredential("user_agent")
placeholder

func (a *Account) GetChatGPTAccountID() string {
	if !a.IsOpenAIOAuth() {
		return ""
placeholder
	return a.GetCredential("chatgpt_account_id")
placeholder

func (a *Account) GetChatGPTUserID() string {
	if !a.IsOpenAIOAuth() {
		return ""
placeholder
	return a.GetCredential("chatgpt_user_id")
placeholder

func (a *Account) GetOpenAIOrganizationID() string {
	if !a.IsOpenAIOAuth() {
		return ""
placeholder
	return a.GetCredential("organization_id")
placeholder

func (a *Account) GetOpenAITokenExpiresAt() *time.Time {
	if !a.IsOpenAIOAuth() {
		return nil
placeholder
	return a.GetCredentialAsTime("expires_at")
placeholder

func (a *Account) IsOpenAITokenExpired() bool {
	expiresAt := a.GetOpenAITokenExpiresAt()
	if expiresAt == nil {
		return false
placeholder
	return time.Now().Add(60 * time.Second).After(*expiresAt)
placeholder

// IsMixedSchedulingEnabled 检查 antigravity 账户是否启用混合调度
// 启用后可参与 anthropic/gemini 分组的账户调度
func (a *Account) IsMixedSchedulingEnabled() bool {
	if a.Platform != PlatformAntigravity {
		return false
placeholder
	if a.Extra == nil {
		return false
placeholder
	if v, ok := a.Extra["mixed_scheduling"]; ok {
		if enabled, ok := v.(bool); ok {
			return enabled
	placeholder
placeholder
	return false
placeholder

// IsOpenAIPassthroughEnabled 返回 OpenAI 账号是否启用“自动透传（仅替换认证）”。
//
// 新字段：accounts.extra.openai_passthrough。
// 兼容字段：accounts.extra.openai_oauth_passthrough（历史 OAuth 开关）。
// 字段缺失或类型不正确时，按 false（关闭）处理。
func (a *Account) IsOpenAIPassthroughEnabled() bool {
	if a == nil || !a.IsOpenAI() || a.Extra == nil {
		return false
placeholder
	if enabled, ok := a.Extra["openai_passthrough"].(bool); ok {
		return enabled
placeholder
	if enabled, ok := a.Extra["openai_oauth_passthrough"].(bool); ok {
		return enabled
placeholder
	return false
placeholder

// IsOpenAIResponsesWebSocketV2Enabled 返回 OpenAI 账号是否开启 Responses WebSocket v2。
//
// 分类型新字段：
// - OAuth 账号：accounts.extra.openai_oauth_responses_websockets_v2_enabled
// - API Key 账号：accounts.extra.openai_apikey_responses_websockets_v2_enabled
//
// 兼容字段：
// - accounts.extra.responses_websockets_v2_enabled
// - accounts.extra.openai_ws_enabled（历史开关）
//
// 优先级：
// 1. 按账号类型读取分类型字段
// 2. 分类型字段缺失时，回退兼容字段
func (a *Account) IsOpenAIResponsesWebSocketV2Enabled() bool {
	if a == nil || !a.IsOpenAI() || a.Extra == nil {
		return false
placeholder
	if a.IsOpenAIOAuth() {
		if enabled, ok := a.Extra["openai_oauth_responses_websockets_v2_enabled"].(bool); ok {
			return enabled
	placeholder
placeholder
	if a.IsOpenAIApiKey() {
		if enabled, ok := a.Extra["openai_apikey_responses_websockets_v2_enabled"].(bool); ok {
			return enabled
	placeholder
placeholder
	if enabled, ok := a.Extra["responses_websockets_v2_enabled"].(bool); ok {
		return enabled
placeholder
	if enabled, ok := a.Extra["openai_ws_enabled"].(bool); ok {
		return enabled
placeholder
	return false
placeholder

const (
	OpenAIWSIngressModeOff         = "off"
	OpenAIWSIngressModeShared      = "shared"
	OpenAIWSIngressModeDedicated   = "dedicated"
	OpenAIWSIngressModeCtxPool     = "ctx_pool"
	OpenAIWSIngressModePassthrough = "passthrough"
)

func normalizeOpenAIWSIngressMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case OpenAIWSIngressModeOff:
		return OpenAIWSIngressModeOff
	case OpenAIWSIngressModeCtxPool:
		return OpenAIWSIngressModeCtxPool
	case OpenAIWSIngressModePassthrough:
		return OpenAIWSIngressModePassthrough
	case OpenAIWSIngressModeShared:
		return OpenAIWSIngressModeShared
	case OpenAIWSIngressModeDedicated:
		return OpenAIWSIngressModeDedicated
	default:
		return ""
placeholder
placeholder

func normalizeOpenAIWSIngressDefaultMode(mode string) string {
	if normalized := normalizeOpenAIWSIngressMode(mode); normalized != "" {
		if normalized == OpenAIWSIngressModeShared || normalized == OpenAIWSIngressModeDedicated {
			return OpenAIWSIngressModeCtxPool
	placeholder
		return normalized
placeholder
	return OpenAIWSIngressModeCtxPool
placeholder

// ResolveOpenAIResponsesWebSocketV2Mode 返回账号在 WSv2 ingress 下的有效模式（off/ctx_pool/passthrough）。
//
// 优先级：
// 1. 分类型 mode 新字段（string）
// 2. 分类型 enabled 旧字段（bool）
// 3. 兼容 enabled 旧字段（bool）
// 4. defaultMode（非法时回退 ctx_pool）
func (a *Account) ResolveOpenAIResponsesWebSocketV2Mode(defaultMode string) string {
	resolvedDefault := normalizeOpenAIWSIngressDefaultMode(defaultMode)
	if a == nil || !a.IsOpenAI() {
		return OpenAIWSIngressModeOff
placeholder
	if a.Extra == nil {
		return resolvedDefault
placeholder

	resolveModeString := func(key string) (string, bool) {
		raw, ok := a.Extra[key]
		if !ok {
			return "", false
	placeholder
		mode, ok := raw.(string)
		if !ok {
			return "", false
	placeholder
		normalized := normalizeOpenAIWSIngressMode(mode)
		if normalized == "" {
			return "", false
	placeholder
		return normalized, true
placeholder
	resolveBoolMode := func(key string) (string, bool) {
		raw, ok := a.Extra[key]
		if !ok {
			return "", false
	placeholder
		enabled, ok := raw.(bool)
		if !ok {
			return "", false
	placeholder
		if enabled {
			return OpenAIWSIngressModeCtxPool, true
	placeholder
		return OpenAIWSIngressModeOff, true
placeholder

	if a.IsOpenAIOAuth() {
		if mode, ok := resolveModeString("openai_oauth_responses_websockets_v2_mode"); ok {
			return mode
	placeholder
		if mode, ok := resolveBoolMode("openai_oauth_responses_websockets_v2_enabled"); ok {
			return mode
	placeholder
placeholder
	if a.IsOpenAIApiKey() {
		if mode, ok := resolveModeString("openai_apikey_responses_websockets_v2_mode"); ok {
			return mode
	placeholder
		if mode, ok := resolveBoolMode("openai_apikey_responses_websockets_v2_enabled"); ok {
			return mode
	placeholder
placeholder
	if mode, ok := resolveBoolMode("responses_websockets_v2_enabled"); ok {
		return mode
placeholder
	if mode, ok := resolveBoolMode("openai_ws_enabled"); ok {
		return mode
placeholder
	// 兼容旧值：shared/dedicated 语义都归并到 ctx_pool。
	if resolvedDefault == OpenAIWSIngressModeShared || resolvedDefault == OpenAIWSIngressModeDedicated {
		return OpenAIWSIngressModeCtxPool
placeholder
	return resolvedDefault
placeholder

// IsOpenAIWSForceHTTPEnabled 返回账号级“强制 HTTP”开关。
// 字段：accounts.extra.openai_ws_force_http。
func (a *Account) IsOpenAIWSForceHTTPEnabled() bool {
	if a == nil || !a.IsOpenAI() || a.Extra == nil {
		return false
placeholder
	enabled, ok := a.Extra["openai_ws_force_http"].(bool)
	return ok && enabled
placeholder

// IsOpenAIWSAllowStoreRecoveryEnabled 返回账号级 store 恢复开关。
// 字段：accounts.extra.openai_ws_allow_store_recovery。
func (a *Account) IsOpenAIWSAllowStoreRecoveryEnabled() bool {
	if a == nil || !a.IsOpenAI() || a.Extra == nil {
		return false
placeholder
	enabled, ok := a.Extra["openai_ws_allow_store_recovery"].(bool)
	return ok && enabled
placeholder

// IsOpenAIOAuthPassthroughEnabled 兼容旧接口，等价于 OAuth 账号的 IsOpenAIPassthroughEnabled。
func (a *Account) IsOpenAIOAuthPassthroughEnabled() bool {
	return a != nil && a.IsOpenAIOAuth() && a.IsOpenAIPassthroughEnabled()
placeholder

// IsAnthropicAPIKeyPassthroughEnabled 返回 Anthropic API Key 账号是否启用“自动透传（仅替换认证）”。
// 字段：accounts.extra.anthropic_passthrough。
// 字段缺失或类型不正确时，按 false（关闭）处理。
func (a *Account) IsAnthropicAPIKeyPassthroughEnabled() bool {
	if a == nil || a.Platform != PlatformAnthropic || a.Type != AccountTypeAPIKey || a.Extra == nil {
		return false
placeholder
	enabled, ok := a.Extra["anthropic_passthrough"].(bool)
	return ok && enabled
placeholder

// IsCodexCLIOnlyEnabled 返回 OpenAI OAuth 账号是否启用“仅允许 Codex 官方客户端”。
// 字段：accounts.extra.codex_cli_only。
// 字段缺失或类型不正确时，按 false（关闭）处理。
func (a *Account) IsCodexCLIOnlyEnabled() bool {
	if a == nil || !a.IsOpenAIOAuth() || a.Extra == nil {
		return false
placeholder
	enabled, ok := a.Extra["codex_cli_only"].(bool)
	return ok && enabled
placeholder

// WindowCostSchedulability 窗口费用调度状态
type WindowCostSchedulability int

const (
	// WindowCostSchedulable 可正常调度
	WindowCostSchedulable WindowCostSchedulability = iota
	// WindowCostStickyOnly 仅允许粘性会话
	WindowCostStickyOnly
	// WindowCostNotSchedulable 完全不可调度
	WindowCostNotSchedulable
)

// IsAnthropicOAuthOrSetupToken 判断是否为 Anthropic OAuth 或 SetupToken 类型账号
// 仅这两类账号支持 5h 窗口额度控制和会话数量控制
func (a *Account) IsAnthropicOAuthOrSetupToken() bool {
	return a.Platform == PlatformAnthropic && (a.Type == AccountTypeOAuth || a.Type == AccountTypeSetupToken)
placeholder

// IsTLSFingerprintEnabled 检查是否启用 TLS 指纹伪装
// 仅适用于 Anthropic OAuth/SetupToken 类型账号
// 启用后将模拟 Claude Code (Node.js) 客户端的 TLS 握手特征
func (a *Account) IsTLSFingerprintEnabled() bool {
	// 仅支持 Anthropic OAuth/SetupToken 账号
	if !a.IsAnthropicOAuthOrSetupToken() {
		return false
placeholder
	if a.Extra == nil {
		return false
placeholder
	if v, ok := a.Extra["enable_tls_fingerprint"]; ok {
		if enabled, ok := v.(bool); ok {
			return enabled
	placeholder
placeholder
	return false
placeholder

// GetUserMsgQueueMode 获取用户消息队列模式
// "serialize" = 串行队列, "throttle" = 软性限速, "" = 未设置（使用全局配置）
func (a *Account) GetUserMsgQueueMode() string {
	if a.Extra == nil {
		return ""
placeholder
	// 优先读取新字段 user_msg_queue_mode（白名单校验，非法值视为未设置）
	if mode, ok := a.Extra["user_msg_queue_mode"].(string); ok && mode != "" {
		if mode == config.UMQModeSerialize || mode == config.UMQModeThrottle {
			return mode
	placeholder
		return "" // 非法值 fallback 到全局配置
placeholder
	// 向后兼容: user_msg_queue_enabled: true → "serialize"
	if enabled, ok := a.Extra["user_msg_queue_enabled"].(bool); ok && enabled {
		return config.UMQModeSerialize
placeholder
	return ""
placeholder

// IsSessionIDMaskingEnabled 检查是否启用会话ID伪装
// 仅适用于 Anthropic OAuth/SetupToken 类型账号
// 启用后将在一段时间内（15分钟）固定 metadata.user_id 中的 session ID，
// 使上游认为请求来自同一个会话
func (a *Account) IsSessionIDMaskingEnabled() bool {
	if !a.IsAnthropicOAuthOrSetupToken() {
		return false
placeholder
	if a.Extra == nil {
		return false
placeholder
	if v, ok := a.Extra["session_id_masking_enabled"]; ok {
		if enabled, ok := v.(bool); ok {
			return enabled
	placeholder
placeholder
	return false
placeholder

// IsCacheTTLOverrideEnabled 检查是否启用缓存 TTL 强制替换
// 仅适用于 Anthropic OAuth/SetupToken 类型账号
// 启用后将所有 cache creation tokens 归入指定的 TTL 类型（5m 或 1h）
func (a *Account) IsCacheTTLOverrideEnabled() bool {
	if !a.IsAnthropicOAuthOrSetupToken() {
		return false
placeholder
	if a.Extra == nil {
		return false
placeholder
	if v, ok := a.Extra["cache_ttl_override_enabled"]; ok {
		if enabled, ok := v.(bool); ok {
			return enabled
	placeholder
placeholder
	return false
placeholder

// GetCacheTTLOverrideTarget 获取缓存 TTL 强制替换的目标类型
// 返回 "5m" 或 "1h"，默认 "5m"
func (a *Account) GetCacheTTLOverrideTarget() string {
	if a.Extra == nil {
		return "5m"
placeholder
	if v, ok := a.Extra["cache_ttl_override_target"]; ok {
		if target, ok := v.(string); ok && (target == "5m" || target == "1h") {
			return target
	placeholder
placeholder
	return "5m"
placeholder

// GetQuotaLimit 获取 API Key 账号的配额限制（美元）
// 返回 0 表示未启用
func (a *Account) GetQuotaLimit() float64 {
	return a.getExtraFloat64("quota_limit")
placeholder

// GetQuotaUsed 获取 API Key 账号的已用配额（美元）
func (a *Account) GetQuotaUsed() float64 {
	return a.getExtraFloat64("quota_used")
placeholder

// GetQuotaDailyLimit 获取日额度限制（美元），0 表示未启用
func (a *Account) GetQuotaDailyLimit() float64 {
	return a.getExtraFloat64("quota_daily_limit")
placeholder

// GetQuotaDailyUsed 获取当日已用额度（美元）
func (a *Account) GetQuotaDailyUsed() float64 {
	return a.getExtraFloat64("quota_daily_used")
placeholder

// GetQuotaWeeklyLimit 获取周额度限制（美元），0 表示未启用
func (a *Account) GetQuotaWeeklyLimit() float64 {
	return a.getExtraFloat64("quota_weekly_limit")
placeholder

// GetQuotaWeeklyUsed 获取本周已用额度（美元）
func (a *Account) GetQuotaWeeklyUsed() float64 {
	return a.getExtraFloat64("quota_weekly_used")
placeholder

// getExtraFloat64 从 Extra 中读取指定 key 的 float64 值
func (a *Account) getExtraFloat64(key string) float64 {
	if a.Extra == nil {
		return 0
placeholder
	if v, ok := a.Extra[key]; ok {
		return parseExtraFloat64(v)
placeholder
	return 0
placeholder

// getExtraTime 从 Extra 中读取 RFC3339 时间戳
func (a *Account) getExtraTime(key string) time.Time {
	if a.Extra == nil {
		return time.Time{placeholder
placeholder
	if v, ok := a.Extra[key]; ok {
		if s, ok := v.(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
				return t
		placeholder
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t
		placeholder
	placeholder
placeholder
	return time.Time{placeholder
placeholder

// HasAnyQuotaLimit 检查是否配置了任一维度的配额限制
func (a *Account) HasAnyQuotaLimit() bool {
	return a.GetQuotaLimit() > 0 || a.GetQuotaDailyLimit() > 0 || a.GetQuotaWeeklyLimit() > 0
placeholder

// isPeriodExpired 检查指定周期（自 periodStart 起经过 dur）是否已过期
func isPeriodExpired(periodStart time.Time, dur time.Duration) bool {
	if periodStart.IsZero() {
		return true // 从未使用过，视为过期（下次 increment 会初始化）
placeholder
	return time.Since(periodStart) >= dur
placeholder

// IsQuotaExceeded 检查 API Key 账号配额是否已超限（任一维度超限即返回 true）
func (a *Account) IsQuotaExceeded() bool {
	// 总额度
	if limit := a.GetQuotaLimit(); limit > 0 && a.GetQuotaUsed() >= limit {
		return true
placeholder
	// 日额度（周期过期视为未超限，下次 increment 会重置）
	if limit := a.GetQuotaDailyLimit(); limit > 0 {
		start := a.getExtraTime("quota_daily_start")
		if !isPeriodExpired(start, 24*time.Hour) && a.GetQuotaDailyUsed() >= limit {
			return true
	placeholder
placeholder
	// 周额度
	if limit := a.GetQuotaWeeklyLimit(); limit > 0 {
		start := a.getExtraTime("quota_weekly_start")
		if !isPeriodExpired(start, 7*24*time.Hour) && a.GetQuotaWeeklyUsed() >= limit {
			return true
	placeholder
placeholder
	return false
placeholder

// GetWindowCostLimit 获取 5h 窗口费用阈值（美元）
// 返回 0 表示未启用
func (a *Account) GetWindowCostLimit() float64 {
	if a.Extra == nil {
		return 0
placeholder
	if v, ok := a.Extra["window_cost_limit"]; ok {
		return parseExtraFloat64(v)
placeholder
	return 0
placeholder

// GetWindowCostStickyReserve 获取粘性会话预留额度（美元）
// 默认值为 10
func (a *Account) GetWindowCostStickyReserve() float64 {
	if a.Extra == nil {
		return 10.0
placeholder
	if v, ok := a.Extra["window_cost_sticky_reserve"]; ok {
		val := parseExtraFloat64(v)
		if val > 0 {
			return val
	placeholder
placeholder
	return 10.0
placeholder

// GetMaxSessions 获取最大并发会话数
// 返回 0 表示未启用
func (a *Account) GetMaxSessions() int {
	if a.Extra == nil {
		return 0
placeholder
	if v, ok := a.Extra["max_sessions"]; ok {
		return parseExtraInt(v)
placeholder
	return 0
placeholder

// GetSessionIdleTimeoutMinutes 获取会话空闲超时分钟数
// 默认值为 5 分钟
func (a *Account) GetSessionIdleTimeoutMinutes() int {
	if a.Extra == nil {
		return 5
placeholder
	if v, ok := a.Extra["session_idle_timeout_minutes"]; ok {
		val := parseExtraInt(v)
		if val > 0 {
			return val
	placeholder
placeholder
	return 5
placeholder

// GetBaseRPM 获取基础 RPM 限制
// 返回 0 表示未启用（负数视为无效配置，按 0 处理）
func (a *Account) GetBaseRPM() int {
	if a.Extra == nil {
		return 0
placeholder
	if v, ok := a.Extra["base_rpm"]; ok {
		val := parseExtraInt(v)
		if val > 0 {
			return val
	placeholder
placeholder
	return 0
placeholder

// GetRPMStrategy 获取 RPM 策略
// "tiered" = 三区模型（默认）, "sticky_exempt" = 粘性豁免
func (a *Account) GetRPMStrategy() string {
	if a.Extra == nil {
		return "tiered"
placeholder
	if v, ok := a.Extra["rpm_strategy"]; ok {
		if s, ok := v.(string); ok && s == "sticky_exempt" {
			return "sticky_exempt"
	placeholder
placeholder
	return "tiered"
placeholder

// GetRPMStickyBuffer 获取 RPM 粘性缓冲数量
// tiered 模式下的黄区大小，默认为 base_rpm 的 20%（至少 1）
func (a *Account) GetRPMStickyBuffer() int {
	if a.Extra == nil {
		return 0
placeholder
	if v, ok := a.Extra["rpm_sticky_buffer"]; ok {
		val := parseExtraInt(v)
		if val > 0 {
			return val
	placeholder
placeholder
	base := a.GetBaseRPM()
	buffer := base / 5
	if buffer < 1 && base > 0 {
		buffer = 1
placeholder
	return buffer
placeholder

// CheckRPMSchedulability 根据当前 RPM 计数检查调度状态
// 复用 WindowCostSchedulability 三态：Schedulable / StickyOnly / NotSchedulable
func (a *Account) CheckRPMSchedulability(currentRPM int) WindowCostSchedulability {
	baseRPM := a.GetBaseRPM()
	if baseRPM <= 0 {
		return WindowCostSchedulable
placeholder

	if currentRPM < baseRPM {
		return WindowCostSchedulable
placeholder

	strategy := a.GetRPMStrategy()
	if strategy == "sticky_exempt" {
		return WindowCostStickyOnly // 粘性豁免无红区
placeholder

	// tiered: 黄区 + 红区
	buffer := a.GetRPMStickyBuffer()
	if currentRPM < baseRPM+buffer {
		return WindowCostStickyOnly
placeholder
	return WindowCostNotSchedulable
placeholder

// CheckWindowCostSchedulability 根据当前窗口费用检查调度状态
// - 费用 < 阈值: WindowCostSchedulable（可正常调度）
// - 费用 >= 阈值 且 < 阈值+预留: WindowCostStickyOnly（仅粘性会话）
// - 费用 >= 阈值+预留: WindowCostNotSchedulable（不可调度）
func (a *Account) CheckWindowCostSchedulability(currentWindowCost float64) WindowCostSchedulability {
	limit := a.GetWindowCostLimit()
	if limit <= 0 {
		return WindowCostSchedulable
placeholder

	if currentWindowCost < limit {
		return WindowCostSchedulable
placeholder

	stickyReserve := a.GetWindowCostStickyReserve()
	if currentWindowCost < limit+stickyReserve {
		return WindowCostStickyOnly
placeholder

	return WindowCostNotSchedulable
placeholder

// GetCurrentWindowStartTime 获取当前有效的窗口开始时间
// 逻辑：
// 1. 如果窗口未过期（SessionWindowEnd 存在且在当前时间之后），使用记录的 SessionWindowStart
// 2. 否则（窗口过期或未设置），使用新的预测窗口开始时间（从当前整点开始）
func (a *Account) GetCurrentWindowStartTime() time.Time {
	now := time.Now()

	// 窗口未过期，使用记录的窗口开始时间
	if a.SessionWindowStart != nil && a.SessionWindowEnd != nil && now.Before(*a.SessionWindowEnd) {
		return *a.SessionWindowStart
placeholder

	// 窗口已过期或未设置，预测新的窗口开始时间（从当前整点开始）
	// 与 ratelimit_service.go 中 UpdateSessionWindow 的预测逻辑保持一致
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
placeholder

// parseExtraFloat64 从 extra 字段解析 float64 值
func parseExtraFloat64(value any) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f
	placeholder
	case string:
		if f, err := strconv.ParseFloat(strings.TrimSpace(v), 64); err == nil {
			return f
	placeholder
placeholder
	return 0
placeholder

// parseExtraInt 从 extra 字段解析 int 值
// ParseExtraInt 从 extra 字段的 any 值解析为 int。
// 支持 int, int64, float64, json.Number, string 类型，无法解析时返回 0。
func ParseExtraInt(value any) int {
	return parseExtraInt(value)
placeholder

func parseExtraInt(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return int(i)
	placeholder
	case string:
		if i, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return i
	placeholder
placeholder
	return 0
placeholder
