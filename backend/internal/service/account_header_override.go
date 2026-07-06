package service

import (
	"net/http"
	"sort"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"

	"golang.org/x/net/http/httpguts"
)

// 请求头覆写（header override）：仅对 Anthropic / OpenAI 平台的 api_key 账号生效。
// 管理员在账号上配置一组 header name -> value，转发到上游前用配置值覆盖同名请求头
// （匹配不区分大小写）；value 为空的条目视为"未填写"，不参与覆盖。
const (
	credKeyHeaderOverrideEnabled = "header_override_enabled"
	credKeyHeaderOverrides       = "header_overrides"

	maxHeaderOverrideEntries     = 64
	maxHeaderOverrideNameLength  = 200
	maxHeaderOverrideValueLength = 8192
)

// headerOverrideBlockedNames 禁止覆写的请求头（小写）。
//   - 连接控制/逐跳头：由 HTTP 栈管理，覆写会破坏请求传输；
//   - host/content-length：由 Go 的 Request.Host / ContentLength 字段管理，header 覆写不生效或产生冲突；
//   - authorization/x-api-key：上游认证头由账号凭据统一注入，禁止通过覆写篡改；
//   - accept-encoding：强制压缩会破坏网关对上游流式响应（SSE/usage）的解析；
//   - sec-websocket-*：WebSocket 握手头由拨号器管理（OpenAI WS 模式）；
//   - session_id/conversation_id 等：逐请求会话隔离头，固定值会造成会话串扰。
var headerOverrideBlockedNames = map[string]struct{placeholder{
	"host":                     {placeholder,
	"content-length":           {placeholder,
	"transfer-encoding":        {placeholder,
	"connection":               {placeholder,
	"keep-alive":               {placeholder,
	"proxy-authenticate":       {placeholder,
	"proxy-authorization":      {placeholder,
	"proxy-connection":         {placeholder,
	"te":                       {placeholder,
	"trailer":                  {placeholder,
	"upgrade":                  {placeholder,
	"authorization":            {placeholder,
	"x-api-key":                {placeholder,
	"accept-encoding":          {placeholder,
	"sec-websocket-key":        {placeholder,
	"sec-websocket-version":    {placeholder,
	"sec-websocket-extensions": {placeholder,
	"sec-websocket-protocol":   {placeholder,
	"sec-websocket-accept":     {placeholder,
	"session_id":               {placeholder,
	"conversation_id":          {placeholder,
	"x-codex-turn-state":       {placeholder,
	"x-codex-turn-metadata":    {placeholder,
	"chatgpt-account-id":       {placeholder,
placeholder

func isHeaderOverrideBlockedName(lowerName string) bool {
	_, blocked := headerOverrideBlockedNames[lowerName]
	return blocked
placeholder

// IsHeaderOverrideEligible 报告账号类型是否支持请求头覆写。
// 目前仅开放 Anthropic / OpenAI 两个平台的 api_key 账号。
func (a *Account) IsHeaderOverrideEligible() bool {
	if a == nil || a.Type != AccountTypeAPIKey {
		return false
placeholder
	return a.Platform == PlatformAnthropic || a.Platform == PlatformOpenAI
placeholder

// IsHeaderOverrideEnabled 报告账号是否启用了请求头覆写。
func (a *Account) IsHeaderOverrideEnabled() bool {
	if !a.IsHeaderOverrideEligible() || a.Credentials == nil {
		return false
placeholder
	enabled, ok := a.Credentials[credKeyHeaderOverrideEnabled].(bool)
	return ok && enabled
placeholder

// GetHeaderOverrides 返回生效的请求头覆写表（key 统一小写）。
// 未启用、不符合平台/类型条件或配置为空时返回 nil。
// 空 value 的条目（模板占位）与非法/禁止的 header 名会被跳过。
func (a *Account) GetHeaderOverrides() map[string]string {
	if !a.IsHeaderOverrideEnabled() {
		return nil
placeholder
	raw := stringMappingFromRaw(a.Credentials[credKeyHeaderOverrides])
	if len(raw) == 0 {
		return nil
placeholder
	result := make(map[string]string, len(raw))
	for name, value := range raw {
		lowerName := strings.ToLower(strings.TrimSpace(name))
		value = strings.TrimSpace(value)
		if lowerName == "" || value == "" {
			continue
	placeholder
		// 防御性过滤：保存路径已做校验，这里兜底未经 Normalize 落库的数据
		if len(lowerName) > maxHeaderOverrideNameLength || len(value) > maxHeaderOverrideValueLength {
			continue
	placeholder
		if isHeaderOverrideBlockedName(lowerName) {
			continue
	placeholder
		if !httpguts.ValidHeaderFieldName(lowerName) || !httpguts.ValidHeaderFieldValue(value) {
			continue
	placeholder
		result[lowerName] = value
placeholder
	if len(result) == 0 {
		return nil
placeholder
	return result
placeholder

// ApplyHeaderOverrides 将账号配置的请求头覆写应用到出站请求头。
// 对每个覆写条目：先删除所有大小写变体（转发链路会以 wire casing 直接写入 map，
// 可能存在非 canonical key），再按已知 wire casing 写入，避免产生重复头。
// 账号未启用或不符合条件时为 no-op，可安全地在 OAuth/api_key 共用的构建器中调用。
func (a *Account) ApplyHeaderOverrides(h http.Header) {
	if h == nil {
		return
placeholder
	overrides := a.GetHeaderOverrides()
	if len(overrides) == 0 {
		return
placeholder
	names := make([]string, 0, len(overrides))
	for name := range overrides {
		names = append(names, name)
placeholder
	sort.Strings(names)
	for _, name := range names {
		for existing := range h {
			if strings.EqualFold(existing, name) {
				delete(h, existing)
		placeholder
	placeholder
		h[resolveWireCasing(name)] = []string{overrides[name]placeholder
placeholder
placeholder

// NormalizeHeaderOverrideCredentials 校验并原地规范化 credentials 中的请求头覆写字段。
// 供账号创建/更新/批量更新的保存路径调用；credentials 未携带相关字段时为 no-op。
// 规范化内容：header 名转小写并去除首尾空白，value 去除首尾空白，丢弃名和值均为空的条目。
func NormalizeHeaderOverrideCredentials(credentials map[string]any) error {
	if credentials == nil {
		return nil
placeholder
	if raw, ok := credentials[credKeyHeaderOverrideEnabled]; ok && raw != nil {
		if _, isBool := raw.(bool); !isBool {
			return infraerrors.New(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header_override_enabled must be a boolean")
	placeholder
placeholder
	raw, ok := credentials[credKeyHeaderOverrides]
	if !ok || raw == nil {
		return nil
placeholder

	var entries map[string]any
	switch m := raw.(type) {
	case map[string]any:
		entries = m
	case map[string]string:
		entries = make(map[string]any, len(m))
		for k, v := range m {
			entries[k] = v
	placeholder
	default:
		return infraerrors.New(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
			"header_overrides must be an object of header name to string value")
placeholder

	if len(entries) > maxHeaderOverrideEntries {
		return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
			"header_overrides supports at most %d entries", maxHeaderOverrideEntries)
placeholder

	normalized := make(map[string]any, len(entries))
	for name, rawValue := range entries {
		value, isString := rawValue.(string)
		if !isString {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header %q value must be a string", name)
	placeholder
		lowerName := strings.ToLower(strings.TrimSpace(name))
		value = strings.TrimSpace(value)
		if lowerName == "" {
			if value == "" {
				continue // 丢弃完全为空的占位行
		placeholder
			return infraerrors.New(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header name must not be empty")
	placeholder
		if len(lowerName) > maxHeaderOverrideNameLength {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header name %q exceeds %d characters", lowerName, maxHeaderOverrideNameLength)
	placeholder
		if !httpguts.ValidHeaderFieldName(lowerName) {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"invalid header name %q", lowerName)
	placeholder
		if isHeaderOverrideBlockedName(lowerName) {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header %q is not allowed to be overridden", lowerName)
	placeholder
		if len(value) > maxHeaderOverrideValueLength {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header %q value exceeds %d characters", lowerName, maxHeaderOverrideValueLength)
	placeholder
		if !httpguts.ValidHeaderFieldValue(value) {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"header %q has an invalid value", lowerName)
	placeholder
		if _, dup := normalized[lowerName]; dup {
			return infraerrors.Newf(http.StatusBadRequest, "INVALID_HEADER_OVERRIDE",
				"duplicate header name %q (matching is case-insensitive)", lowerName)
	placeholder
		normalized[lowerName] = value
placeholder
	credentials[credKeyHeaderOverrides] = normalized
	return nil
placeholder
