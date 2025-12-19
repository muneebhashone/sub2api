package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// JSONB 用于存储JSONB数据
type JSONB map[string]interface{placeholder

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
placeholder
	return json.Marshal(j)
placeholder

func (j *JSONB) Scan(value interface{placeholder) error {
	if value == nil {
		*j = nil
		return nil
placeholder
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
placeholder
	return json.Unmarshal(bytes, j)
placeholder

type Account struct {
	ID           int64          `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null" json:"name"`
	Platform     string         `gorm:"size:50;not null" json:"platform"`           // anthropic/openai/gemini
	Type         string         `gorm:"size:20;not null" json:"type"`               // oauth/apikey
	Credentials  JSONB          `gorm:"type:jsonb;default:'{placeholder'" json:"credentials"` // 凭证(加密存储)
	Extra        JSONB          `gorm:"type:jsonb;default:'{placeholder'" json:"extra"`       // 扩展信息
	ProxyID      *int64         `gorm:"index" json:"proxy_id"`
	Concurrency  int            `gorm:"default:3;not null" json:"concurrency"`
	Priority     int            `gorm:"default:50;not null" json:"priority"`            // 1-100，越小越高
	Status       string         `gorm:"size:20;default:active;not null" json:"status"`  // active/disabled/error
	ErrorMessage string         `gorm:"type:text" json:"error_message"`
	LastUsedAt   *time.Time     `gorm:"index" json:"last_used_at"`
	CreatedAt    time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 调度控制
	Schedulable bool `gorm:"default:true;not null" json:"schedulable"`

	// 限流状态 (429)
	RateLimitedAt    *time.Time `gorm:"index" json:"rate_limited_at"`
	RateLimitResetAt *time.Time `gorm:"index" json:"rate_limit_reset_at"`

	// 过载状态 (529)
	OverloadUntil *time.Time `gorm:"index" json:"overload_until"`

	// 5小时时间窗口
	SessionWindowStart  *time.Time `json:"session_window_start"`
	SessionWindowEnd    *time.Time `json:"session_window_end"`
	SessionWindowStatus string     `gorm:"size:20" json:"session_window_status"` // allowed/allowed_warning/rejected

	// 关联
	Proxy         *Proxy         `gorm:"foreignKey:ProxyID" json:"proxy,omitempty"`
	AccountGroups []AccountGroup `gorm:"foreignKey:AccountID" json:"account_groups,omitempty"`

	// 虚拟字段 (不存储到数据库)
	GroupIDs []int64 `gorm:"-" json:"group_ids,omitempty"`
placeholder

func (Account) TableName() string {
	return "accounts"
placeholder

// IsActive 检查是否激活
func (a *Account) IsActive() bool {
	return a.Status == "active"
placeholder

// IsSchedulable 检查账号是否可调度
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

// IsRateLimited 检查是否处于限流状态
func (a *Account) IsRateLimited() bool {
	if a.RateLimitResetAt == nil {
		return false
placeholder
	return time.Now().Before(*a.RateLimitResetAt)
placeholder

// IsOverloaded 检查是否处于过载状态
func (a *Account) IsOverloaded() bool {
	if a.OverloadUntil == nil {
		return false
placeholder
	return time.Now().Before(*a.OverloadUntil)
placeholder

// IsOAuth 检查是否为OAuth类型账号（包括oauth和setup-token）
func (a *Account) IsOAuth() bool {
	return a.Type == AccountTypeOAuth || a.Type == AccountTypeSetupToken
placeholder

// CanGetUsage 检查账号是否可以获取usage信息（只有oauth类型可以，setup-token没有profile权限）
func (a *Account) CanGetUsage() bool {
	return a.Type == AccountTypeOAuth
placeholder

// GetCredential 获取凭证字段
func (a *Account) GetCredential(key string) string {
	if a.Credentials == nil {
		return ""
placeholder
	if v, ok := a.Credentials[key]; ok {
		if s, ok := v.(string); ok {
			return s
	placeholder
placeholder
	return ""
placeholder

// GetModelMapping 获取模型映射配置
// 返回格式: map[请求模型名]实际模型名
func (a *Account) GetModelMapping() map[string]string {
	if a.Credentials == nil {
		return nil
placeholder
	raw, ok := a.Credentials["model_mapping"]
	if !ok || raw == nil {
		return nil
placeholder
	// 处理map[string]interface{placeholder类型
	if m, ok := raw.(map[string]interface{placeholder); ok {
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

// IsModelSupported 检查请求的模型是否被该账号支持
// 如果没有设置模型映射，则支持所有模型
func (a *Account) IsModelSupported(requestedModel string) bool {
	mapping := a.GetModelMapping()
	if mapping == nil || len(mapping) == 0 {
		return true // 没有映射配置，支持所有模型
placeholder
	_, exists := mapping[requestedModel]
	return exists
placeholder

// GetMappedModel 获取映射后的实际模型名
// 如果没有映射，返回原始模型名
func (a *Account) GetMappedModel(requestedModel string) string {
	mapping := a.GetModelMapping()
	if mapping == nil || len(mapping) == 0 {
		return requestedModel
placeholder
	if mappedModel, exists := mapping[requestedModel]; exists {
		return mappedModel
placeholder
	return requestedModel
placeholder

// GetBaseURL 获取API基础URL（用于apikey类型账号）
func (a *Account) GetBaseURL() string {
	if a.Type != AccountTypeApiKey {
		return ""
placeholder
	baseURL := a.GetCredential("base_url")
	if baseURL == "" {
		return "https://api.anthropic.com" // 默认URL
placeholder
	return baseURL
placeholder

// GetExtraString 从Extra字段获取字符串值
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

// IsCustomErrorCodesEnabled 检查是否启用自定义错误码功能（仅适用于 apikey 类型）
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

// GetCustomErrorCodes 获取自定义错误码列表
func (a *Account) GetCustomErrorCodes() []int {
	if a.Credentials == nil {
		return nil
placeholder
	raw, ok := a.Credentials["custom_error_codes"]
	if !ok || raw == nil {
		return nil
placeholder
	// 处理 []interface{placeholder 类型（JSON反序列化后的格式）
	if arr, ok := raw.([]interface{placeholder); ok {
		result := make([]int, 0, len(arr))
		for _, v := range arr {
			// JSON 数字默认解析为 float64
			if f, ok := v.(float64); ok {
				result = append(result, int(f))
		placeholder
	placeholder
		return result
placeholder
	return nil
placeholder

// ShouldHandleErrorCode 检查指定错误码是否应该被处理（停止调度/标记限流等）
// 如果未启用自定义错误码或列表为空，返回 true（使用默认策略）
// 如果启用且列表非空，只有在列表中的错误码才返回 true
func (a *Account) ShouldHandleErrorCode(statusCode int) bool {
	if !a.IsCustomErrorCodesEnabled() {
		return true // 未启用，使用默认策略
placeholder
	codes := a.GetCustomErrorCodes()
	if len(codes) == 0 {
		return true // 启用但列表为空，fallback到默认策略
placeholder
	// 检查是否在自定义列表中
	for _, code := range codes {
		if code == statusCode {
			return true
	placeholder
placeholder
	return false
placeholder

// IsInterceptWarmupEnabled 检查是否启用预热请求拦截
// 启用后，标题生成、Warmup等预热请求将返回mock响应，不消耗上游token
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
