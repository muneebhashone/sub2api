// Package service provides business logic and domain services for the application.
package service

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	ID                 int64
	Name               string
	Notes              *string
	Platform           string
	Type               string
	Credentials        map[string]any
	Extra              map[string]any
	ProxyID            *int64
	Concurrency        int
	Priority           int
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
	if a.Credentials == nil {
		return nil
placeholder
	raw, ok := a.Credentials["model_mapping"]
	if !ok || raw == nil {
		return nil
placeholder
	if m, ok := raw.(map[string]any); ok {
		result := make(map[string]string)
		for k, v := range m {
			if s, ok := v.(string); ok {
				result[k] = s
		placeholder
	placeholder
		if len(result) > 0 {
			return result
	placeholder
placeholder
	return nil
placeholder

func (a *Account) IsModelSupported(requestedModel string) bool {
	mapping := a.GetModelMapping()
	if len(mapping) == 0 {
		return true
placeholder
	_, exists := mapping[requestedModel]
	return exists
placeholder

func (a *Account) GetMappedModel(requestedModel string) string {
	mapping := a.GetModelMapping()
	if len(mapping) == 0 {
		return requestedModel
placeholder
	if mappedModel, exists := mapping[requestedModel]; exists {
		return mappedModel
placeholder
	return requestedModel
placeholder

func (a *Account) GetBaseURL() string {
	if a.Type != AccountTypeAPIKey {
		return ""
placeholder
	baseURL := a.GetCredential("base_url")
	if baseURL == "" {
		return "https://api.anthropic.com"
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
