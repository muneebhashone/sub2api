package service

import (
	"encoding/json"
	"strconv"
	"time"
)

type Account struct {
	ID           int64
	Name         string
	Platform     string
	Type         string
	Credentials  map[string]any
	Extra        map[string]any
	ProxyID      *int64
	Concurrency  int
	Priority     int
	Status       string
	ErrorMessage string
	LastUsedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Schedulable bool

	RateLimitedAt    *time.Time
	RateLimitResetAt *time.Time
	OverloadUntil    *time.Time

	SessionWindowStart  *time.Time
	SessionWindowEnd    *time.Time
	SessionWindowStatus string

	Proxy         *Proxy
	AccountGroups []AccountGroup
	GroupIDs      []int64
	Groups        []*Group
placeholder

func (a *Account) IsActive() bool {
	return a.Status == StatusActive
placeholder

func (a *Account) IsSchedulable() bool {
	if !a.IsActive() || !a.Schedulable {
		return false
placeholder
	now := time.Now()
	if a.OverloadUntil != nil && now.Before(*a.OverloadUntil) {
		return false
placeholder
	if a.RateLimitResetAt != nil && now.Before(*a.RateLimitResetAt) {
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
	if a.Type != AccountTypeApiKey {
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
	if a.Type != AccountTypeApiKey || a.Credentials == nil {
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
	return a.IsOpenAI() && a.Type == AccountTypeApiKey
placeholder

func (a *Account) GetOpenAIBaseURL() string {
	if !a.IsOpenAI() {
		return ""
placeholder
	if a.Type == AccountTypeApiKey {
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
	expiresAtStr := a.GetCredential("expires_at")
	if expiresAtStr == "" {
		return nil
placeholder
	t, err := time.Parse(time.RFC3339, expiresAtStr)
	if err != nil {
		if v, ok := a.Credentials["expires_at"].(float64); ok {
			tt := time.Unix(int64(v), 0)
			return &tt
	placeholder
		return nil
placeholder
	return &t
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
