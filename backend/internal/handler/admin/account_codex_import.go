package admin

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const codexImportClockSkewSeconds int64 = 120

type CodexSessionImportRequest struct {
	Content                 string         `json:"content"`
	Contents                []string       `json:"contents"`
	Name                    string         `json:"name"`
	Notes                   *string        `json:"notes"`
	GroupIDs                []int64        `json:"group_ids"`
	ProxyID                 *int64         `json:"proxy_id"`
	Concurrency             *int           `json:"concurrency"`
	Priority                *int           `json:"priority"`
	RateMultiplier          *float64       `json:"rate_multiplier"`
	LoadFactor              *int           `json:"load_factor"`
	ExpiresAt               *int64         `json:"expires_at"`
	AutoPauseOnExpired      *bool          `json:"auto_pause_on_expired"`
	CredentialExtras        map[string]any `json:"credential_extras"`
	Extra                   map[string]any `json:"extra"`
	UpdateExisting          *bool          `json:"update_existing"`
	SkipDefaultGroupBind    *bool          `json:"skip_default_group_bind"`
	ConfirmMixedChannelRisk *bool          `json:"confirm_mixed_channel_risk"`
placeholder

type CodexSessionImportResult struct {
	Total    int                         `json:"total"`
	Created  int                         `json:"created"`
	Updated  int                         `json:"updated"`
	Skipped  int                         `json:"skipped"`
	Failed   int                         `json:"failed"`
	Items    []CodexSessionImportItem    `json:"items,omitempty"`
	Warnings []CodexSessionImportMessage `json:"warnings,omitempty"`
	Errors   []CodexSessionImportMessage `json:"errors,omitempty"`
placeholder

type CodexSessionImportItem struct {
	Index     int    `json:"index"`
	Name      string `json:"name,omitempty"`
	Action    string `json:"action"`
	AccountID int64  `json:"account_id,omitempty"`
	Message   string `json:"message,omitempty"`
placeholder

type CodexSessionImportMessage struct {
	Index   int    `json:"index"`
	Name    string `json:"name,omitempty"`
	Message string `json:"message"`
placeholder

type codexImportEntry struct {
	Index int
	Value any
placeholder

type codexImportAccount struct {
	Name            string
	AccessToken     string
	RefreshToken    string
	IDToken         string
	Email           string
	AccountID       string
	UserID          string
	PlanType        string
	Organization    string
	AgentRuntimeID  string
	AgentPrivateKey string
	AgentTaskID     string
	AgentFedRAMP    bool
	IsAgentIdentity bool
	Credentials     map[string]any
	Extra           map[string]any
	TokenExpiresAt  *time.Time
	IdentityKeys    []string
	WarningTexts    []string
placeholder

type codexJWTClaims struct {
	Sub        string                `json:"sub"`
	Email      string                `json:"email"`
	Exp        int64                 `json:"exp"`
	Iat        int64                 `json:"iat"`
	OpenAIAuth *codexJWTOpenAIClaims `json:"https://api.openai.com/auth,omitempty"`
placeholder

type codexJWTOpenAIClaims struct {
	ChatGPTAccountID string                     `json:"chatgpt_account_id"`
	ChatGPTUserID    string                     `json:"chatgpt_user_id"`
	ChatGPTPlanType  string                     `json:"chatgpt_plan_type"`
	UserID           string                     `json:"user_id"`
	POID             string                     `json:"poid"`
	Organizations    []openai.OrganizationClaim `json:"organizations"`
placeholder

type codexAccountIndex struct {
	accountsByKey map[string][]service.Account
placeholder

func (h *AccountHandler) ImportCodexSession(c *gin.Context) {
	var req CodexSessionImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	if err := service.ValidateOpenAILongContextBillingExtra(service.PlatformOpenAI, req.Extra); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if req.Concurrency != nil && *req.Concurrency < 0 {
		response.BadRequest(c, "concurrency must be >= 0")
		return
placeholder
	if req.Priority != nil && *req.Priority < 0 {
		response.BadRequest(c, "priority must be >= 0")
		return
placeholder
	if req.RateMultiplier != nil && *req.RateMultiplier < 0 {
		response.BadRequest(c, "rate_multiplier must be >= 0")
		return
placeholder
	if req.LoadFactor != nil && *req.LoadFactor > 10000 {
		response.BadRequest(c, "load_factor must be <= 10000")
		return
placeholder

	entries, err := parseCodexSessionImportEntries(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder
	if len(entries) == 0 {
		response.BadRequest(c, "请输入 accessToken 或 Codex session JSON")
		return
placeholder

	executeAdminIdempotentJSON(c, "admin.accounts.import_codex_session", req, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		return h.importCodexSessions(ctx, req, entries)
placeholder)
placeholder

func (h *AccountHandler) importCodexSessions(ctx context.Context, req CodexSessionImportRequest, entries []codexImportEntry) (CodexSessionImportResult, error) {
	result := CodexSessionImportResult{
		Total: len(entries),
		Items: make([]CodexSessionImportItem, 0, len(entries)),
placeholder

	existingAccounts, err := h.listAccountsFiltered(ctx, service.PlatformOpenAI, service.AccountTypeOAuth, "", "", 0, "", "created_at", "desc")
	if err != nil {
		return result, err
placeholder
	index := buildCodexAccountIndex(existingAccounts)

	updateExisting := true
	if req.UpdateExisting != nil {
		updateExisting = *req.UpdateExisting
placeholder
	concurrency := 3
	if req.Concurrency != nil {
		concurrency = *req.Concurrency
placeholder
	priority := 50
	if req.Priority != nil {
		priority = *req.Priority
placeholder
	credentialExtras := sanitizeCodexImportCredentialExtras(req.CredentialExtras)
	skipDefaultGroupBind := false
	if req.SkipDefaultGroupBind != nil {
		skipDefaultGroupBind = *req.SkipDefaultGroupBind
placeholder
	skipMixedChannelCheck := req.ConfirmMixedChannelRisk != nil && *req.ConfirmMixedChannelRisk

	seenIdentity := map[string]codexSeenIdentity{placeholder
	for _, entry := range entries {
		item, err := normalizeCodexImportEntry(entry)
		if err != nil {
			result.Failed++
			result.Items = append(result.Items, CodexSessionImportItem{
				Index:   entry.Index,
				Action:  "failed",
				Message: err.Error(),
		placeholder)
			result.Errors = append(result.Errors, CodexSessionImportMessage{
				Index:   entry.Index,
				Message: err.Error(),
		placeholder)
			continue
	placeholder
		accountName := buildCodexCreateAccountName(req.Name, item, entry.Index, len(entries))
		effectiveExpiresAt, credentialExpiresAt, autoPauseOnExpired, expiryWarnings, expiryErr := resolveCodexImportExpiry(req, item)
		if expiryErr != nil {
			result.Failed++
			result.Items = append(result.Items, CodexSessionImportItem{
				Index:   entry.Index,
				Name:    accountName,
				Action:  "failed",
				Message: expiryErr.Error(),
		placeholder)
			result.Errors = append(result.Errors, CodexSessionImportMessage{
				Index:   entry.Index,
				Name:    accountName,
				Message: expiryErr.Error(),
		placeholder)
			continue
	placeholder
		item.WarningTexts = append(item.WarningTexts, expiryWarnings...)
		if credentialExpiresAt != nil {
			item.Credentials["expires_at"] = credentialExpiresAt.Format(time.RFC3339)
	placeholder
		credentials := mergeCodexImportMap(item.Credentials, credentialExtras)
		extra := mergeCodexImportMap(req.Extra, item.Extra)
		for _, warning := range item.WarningTexts {
			result.Warnings = append(result.Warnings, CodexSessionImportMessage{
				Index:   entry.Index,
				Name:    accountName,
				Message: warning,
		placeholder)
	placeholder

		if duplicateIndex, ok := firstSeenCodexIdentity(seenIdentity, item.IdentityKeys, item.UserID); ok {
			message := fmt.Sprintf("与第 %d 条导入项重复，已跳过", duplicateIndex)
			result.Skipped++
			result.Items = append(result.Items, CodexSessionImportItem{
				Index:   entry.Index,
				Name:    accountName,
				Action:  "skipped",
				Message: message,
		placeholder)
			result.Warnings = append(result.Warnings, CodexSessionImportMessage{
				Index:   entry.Index,
				Name:    accountName,
				Message: message,
		placeholder)
			continue
	placeholder
		markCodexIdentitySeen(seenIdentity, item.IdentityKeys, entry.Index, item.UserID)

		existing, matchedKey := index.Find(item.IdentityKeys, item.UserID)
		if existing != nil && updateExisting {
			if strings.HasPrefix(matchedKey, "account:") && item.UserID != "" &&
				codexCredentialString(existing.Credentials, "chatgpt_user_id") == "" {
				result.Warnings = append(result.Warnings, CodexSessionImportMessage{
					Index:   entry.Index,
					Name:    accountName,
					Message: "已有账号未记录 chatgpt_user_id，已按共享的 chatgpt_account_id 匹配并回填，请确认两者属于同一用户",
			placeholder)
		placeholder
			preserveExistingRefresh := item.RefreshToken == "" &&
				codexCredentialString(existing.Credentials, "refresh_token") != ""
			if preserveExistingRefresh {
				result.Warnings = append(result.Warnings, CodexSessionImportMessage{
					Index:   entry.Index,
					Name:    accountName,
					Message: "已有账号包含 refresh_token，本次 accessToken-only 导入已保留自动续期凭据",
			placeholder)
				effectiveExpiresAt = nil
				autoPauseOnExpired = nil
		placeholder
			mergedCredentials := mergeCodexImportCredentials(existing.Credentials, credentials, item)
			mergedExtra := mergeCodexImportMap(existing.Extra, extra)
			updateInput := &service.UpdateAccountInput{
				Credentials:        mergedCredentials,
				Extra:              mergedExtra,
				Concurrency:        req.Concurrency,
				Priority:           req.Priority,
				RateMultiplier:     req.RateMultiplier,
				LoadFactor:         req.LoadFactor,
				ExpiresAt:          effectiveExpiresAt,
				AutoPauseOnExpired: autoPauseOnExpired,
		placeholder
			if req.ProxyID != nil {
				updateInput.ProxyID = req.ProxyID
		placeholder
			if len(req.GroupIDs) > 0 {
				groupIDs := append([]int64(nil), req.GroupIDs...)
				updateInput.GroupIDs = &groupIDs
				updateInput.SkipMixedChannelCheck = skipMixedChannelCheck
		placeholder
			updated, updateErr := h.adminService.UpdateAccount(ctx, existing.ID, updateInput)
			if updateErr != nil {
				result.Failed++
				result.Items = append(result.Items, CodexSessionImportItem{
					Index:   entry.Index,
					Name:    accountName,
					Action:  "failed",
					Message: updateErr.Error(),
			placeholder)
				result.Errors = append(result.Errors, CodexSessionImportMessage{
					Index:   entry.Index,
					Name:    accountName,
					Message: updateErr.Error(),
			placeholder)
				continue
		placeholder
			if h.tokenCacheInvalidator != nil && updated != nil {
				_ = h.tokenCacheInvalidator.InvalidateToken(ctx, updated)
		placeholder
			result.Updated++
			accountID := existing.ID
			if updated != nil {
				accountID = updated.ID
				index.Add(*updated)
		placeholder
			result.Items = append(result.Items, CodexSessionImportItem{
				Index:     entry.Index,
				Name:      accountName,
				Action:    "updated",
				AccountID: accountID,
		placeholder)
			continue
	placeholder

		account, createErr := h.adminService.CreateAccount(ctx, &service.CreateAccountInput{
			Name:                  accountName,
			Notes:                 req.Notes,
			Platform:              service.PlatformOpenAI,
			Type:                  service.AccountTypeOAuth,
			Credentials:           credentials,
			Extra:                 extra,
			ProxyID:               req.ProxyID,
			Concurrency:           concurrency,
			Priority:              priority,
			RateMultiplier:        req.RateMultiplier,
			LoadFactor:            req.LoadFactor,
			GroupIDs:              req.GroupIDs,
			ExpiresAt:             effectiveExpiresAt,
			AutoPauseOnExpired:    autoPauseOnExpired,
			SkipDefaultGroupBind:  skipDefaultGroupBind,
			SkipMixedChannelCheck: skipMixedChannelCheck,
	placeholder)
		if createErr != nil {
			result.Failed++
			result.Items = append(result.Items, CodexSessionImportItem{
				Index:   entry.Index,
				Name:    accountName,
				Action:  "failed",
				Message: createErr.Error(),
		placeholder)
			result.Errors = append(result.Errors, CodexSessionImportMessage{
				Index:   entry.Index,
				Name:    accountName,
				Message: createErr.Error(),
		placeholder)
			continue
	placeholder
		if account != nil {
			index.Add(*account)
	placeholder
		result.Created++
		accountID := int64(0)
		if account != nil {
			accountID = account.ID
	placeholder
		result.Items = append(result.Items, CodexSessionImportItem{
			Index:     entry.Index,
			Name:      accountName,
			Action:    "created",
			AccountID: accountID,
	placeholder)
placeholder

	return result, nil
placeholder

func parseCodexSessionImportEntries(req CodexSessionImportRequest) ([]codexImportEntry, error) {
	contents := make([]string, 0, 1+len(req.Contents))
	if strings.TrimSpace(req.Content) != "" {
		contents = append(contents, req.Content)
placeholder
	for _, content := range req.Contents {
		if strings.TrimSpace(content) != "" {
			contents = append(contents, content)
	placeholder
placeholder

	var entries []codexImportEntry
	for _, content := range contents {
		values, err := parseCodexSessionImportContent(content)
		if err != nil {
			return nil, err
	placeholder
		for _, value := range values {
			entries = append(entries, codexImportEntry{
				Index: len(entries) + 1,
				Value: value,
		placeholder)
	placeholder
placeholder
	return entries, nil
placeholder

func parseCodexSessionImportContent(content string) ([]any, error) {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, nil
placeholder

	if looksLikeJSON(trimmed) {
		values, err := decodeCodexJSONStream(trimmed)
		if err != nil {
			if strings.Contains(trimmed, "\n") {
				if lineValues, lineErr := parseCodexSessionImportLines(trimmed); lineErr == nil {
					return lineValues, nil
			placeholder
		placeholder
			return nil, fmt.Errorf("JSON 解析失败: %w", err)
	placeholder
		return flattenCodexImportValues(values), nil
placeholder

	return parseCodexSessionImportLines(trimmed)
placeholder

func parseCodexSessionImportLines(content string) ([]any, error) {
	values := make([]any, 0)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
	placeholder
		if looksLikeJSON(line) {
			lineValues, err := decodeCodexJSONStream(line)
			if err != nil {
				return nil, fmt.Errorf("第 %d 行 JSON 解析失败: %w", len(values)+1, err)
		placeholder
			values = append(values, flattenCodexImportValues(lineValues)...)
			continue
	placeholder
		values = append(values, line)
placeholder
	return values, nil
placeholder

func decodeCodexJSONStream(content string) ([]any, error) {
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()
	values := make([]any, 0, 1)
	for {
		var value any
		err := decoder.Decode(&value)
		if errors.Is(err, io.EOF) {
			break
	placeholder
		if err != nil {
			return nil, err
	placeholder
		values = append(values, value)
placeholder
	if len(values) == 0 {
		return nil, errors.New("空 JSON 内容")
placeholder
	return values, nil
placeholder

func flattenCodexImportValues(values []any) []any {
	out := make([]any, 0, len(values))
	var appendValue func(any)
	appendValue = func(value any) {
		if arr, ok := value.([]any); ok {
			for _, item := range arr {
				appendValue(item)
		placeholder
			return
	placeholder
		out = append(out, value)
placeholder
	for _, value := range values {
		appendValue(value)
placeholder
	return out
placeholder

func normalizeCodexImportEntry(entry codexImportEntry) (*codexImportAccount, error) {
	now := time.Now().UTC()
	item := &codexImportAccount{
placeholderplaceholder,
		Extra: map[string]any{
			"import_source": "codex_session",
			"imported_at":   now.Format(time.RFC3339),
	placeholder,
placeholder

	switch raw := entry.Value.(type) {
	case string:
		item.AccessToken = strings.TrimSpace(raw)
	case map[string]any:
		if agentIdentity, ok := firstCodexMap(raw, []string{"agent_identity"placeholder, []string{"agentIdentity"placeholder); ok || strings.EqualFold(firstCodexString(raw, []string{"auth_mode"placeholder, []string{"authMode"placeholder), service.OpenAIAuthModeAgentIdentity) {
			if !ok {
				agentIdentity = raw
		placeholder
			item.IsAgentIdentity = true
			item.AgentRuntimeID = firstCodexString(agentIdentity, []string{"agent_runtime_id"placeholder, []string{"agentRuntimeId"placeholder)
			item.AgentPrivateKey = firstCodexString(agentIdentity, []string{"agent_private_key"placeholder, []string{"agentPrivateKey"placeholder)
			item.AgentTaskID = firstCodexString(agentIdentity, []string{"task_id"placeholder, []string{"taskId"placeholder)
			item.AccountID = firstCodexString(agentIdentity, []string{"account_id"placeholder, []string{"accountId"placeholder)
			item.UserID = firstCodexString(agentIdentity, []string{"chatgpt_user_id"placeholder, []string{"chatgptUserId"placeholder)
			item.Email = firstCodexString(agentIdentity, []string{"email"placeholder)
			item.PlanType = firstCodexString(agentIdentity, []string{"plan_type"placeholder, []string{"planType"placeholder)
			item.AgentFedRAMP = firstCodexBool(agentIdentity, []string{"chatgpt_account_is_fedramp"placeholder, []string{"chatgptAccountIsFedramp"placeholder)
			if item.AgentRuntimeID == "" || item.AgentPrivateKey == "" || item.AccountID == "" || item.UserID == "" {
				return nil, errors.New("agent identity 缺少必要字段")
		placeholder
			if err := service.ValidateOpenAIAgentIdentityPrivateKey(item.AgentPrivateKey); err != nil {
				return nil, errors.New("agent identity private key 格式无效")
		placeholder
			item.Credentials["auth_mode"] = service.OpenAIAuthModeAgentIdentity
			item.Credentials["agent_runtime_id"] = item.AgentRuntimeID
			item.Credentials["agent_private_key"] = item.AgentPrivateKey
			item.Credentials["chatgpt_account_id"] = item.AccountID
			item.Credentials["chatgpt_user_id"] = item.UserID
			item.Credentials["chatgpt_account_is_fedramp"] = item.AgentFedRAMP
			setCodexCredentialIfNotEmpty(item.Credentials, "task_id", item.AgentTaskID)
			setCodexCredentialIfNotEmpty(item.Credentials, "email", item.Email)
			setCodexCredentialIfNotEmpty(item.Credentials, "plan_type", item.PlanType)
			if item.AgentTaskID == "" {
				item.WarningTexts = append(item.WarningTexts, "未包含 task_id，首次请求会使用现有 runtime 注册新 task")
		placeholder
			item.IdentityKeys = buildCodexAgentIdentityKeys(item.AccountID, item.UserID, item.Email, item.AgentRuntimeID)
			item.Name = buildCodexImportAccountName(item, entry.Index)
			return item, nil
	placeholder
		item.AccessToken = firstCodexString(raw,
			[]string{"tokens", "access_token"placeholder,
			[]string{"tokens", "accessToken"placeholder,
			[]string{"access_token"placeholder,
			[]string{"accessToken"placeholder,
			[]string{"token"placeholder,
		)
		item.RefreshToken = firstCodexString(raw,
			[]string{"tokens", "refresh_token"placeholder,
			[]string{"tokens", "refreshToken"placeholder,
			[]string{"refresh_token"placeholder,
			[]string{"refreshToken"placeholder,
		)
		item.IDToken = firstCodexString(raw,
			[]string{"tokens", "id_token"placeholder,
			[]string{"tokens", "idToken"placeholder,
			[]string{"id_token"placeholder,
			[]string{"idToken"placeholder,
		)
		item.Email = firstCodexString(raw, []string{"email"placeholder, []string{"user", "email"placeholder)
		item.AccountID = firstCodexString(raw,
			[]string{"chatgpt_account_id"placeholder,
			[]string{"chatgptAccountId"placeholder,
			[]string{"account_id"placeholder,
			[]string{"accountId"placeholder,
			[]string{"account", "id"placeholder,
			[]string{"account", "account_id"placeholder,
			[]string{"account", "chatgpt_account_id"placeholder,
		)
		item.UserID = firstCodexString(raw,
			[]string{"chatgpt_user_id"placeholder,
			[]string{"chatgptUserId"placeholder,
			[]string{"user_id"placeholder,
			[]string{"userId"placeholder,
			[]string{"user", "id"placeholder,
		)
		item.PlanType = firstCodexString(raw,
			[]string{"plan_type"placeholder,
			[]string{"planType"placeholder,
			[]string{"account", "plan_type"placeholder,
			[]string{"account", "planType"placeholder,
		)
		item.Organization = firstCodexString(raw,
			[]string{"organization_id"placeholder,
			[]string{"organizationId"placeholder,
			[]string{"org_id"placeholder,
			[]string{"orgId"placeholder,
		)
		item.Name = firstCodexString(raw, []string{"name"placeholder, []string{"user", "name"placeholder)
		authProvider := firstCodexString(raw, []string{"auth_provider"placeholder, []string{"authProvider"placeholder)
		if authProvider != "" {
			item.Extra["auth_provider"] = authProvider
	placeholder
		if sessionToken := firstCodexString(raw, []string{"session_token"placeholder, []string{"sessionToken"placeholder); sessionToken != "" {
			item.Extra["session_token_present"] = true
			item.WarningTexts = append(item.WarningTexts, "sessionToken 已忽略，不会作为 OAuth refresh_token 存储")
	placeholder
		if sessionExpiresAt, ok := codexTimeAt(raw, []string{"expires"placeholder); ok {
			item.Extra["session_expires_at"] = sessionExpiresAt.Format(time.RFC3339)
	placeholder
		if tokenExpiresAt, ok := firstCodexTime(raw,
			[]string{"tokens", "expires_at"placeholder,
			[]string{"tokens", "expiresAt"placeholder,
			[]string{"expires_at"placeholder,
			[]string{"expiresAt"placeholder,
		); ok {
			if tokenExpiresAt.Unix() <= now.Unix()-codexImportClockSkewSeconds {
				return nil, fmt.Errorf("access_token 已过期: %s", tokenExpiresAt.Format(time.RFC3339))
		placeholder
			item.TokenExpiresAt = &tokenExpiresAt
			item.Credentials["expires_at"] = tokenExpiresAt.Format(time.RFC3339)
	placeholder
		copyCodexExtraString(raw, item.Extra, "user_image", []string{"user", "image"placeholder)
		copyCodexExtraString(raw, item.Extra, "user_picture", []string{"user", "picture"placeholder)
		copyCodexExtraString(raw, item.Extra, "account_structure", []string{"account", "structure"placeholder)
		copyCodexExtraString(raw, item.Extra, "account_residency_region", []string{"account", "residencyRegion"placeholder)
		copyCodexExtraString(raw, item.Extra, "compute_residency", []string{"account", "computeResidency"placeholder)
	default:
		return nil, fmt.Errorf("第 %d 条格式不支持", entry.Index)
placeholder

	if item.IsAgentIdentity {
		return item, nil
placeholder
	if item.AccessToken == "" {
		return nil, errors.New("缺少 accessToken/access_token")
placeholder
	item.Credentials["access_token"] = item.AccessToken
	if item.RefreshToken != "" {
		item.Credentials["refresh_token"] = item.RefreshToken
		item.Credentials["client_id"] = openai.ClientID
placeholder
	if item.IDToken != "" {
		item.Credentials["id_token"] = item.IDToken
		_ = enrichCodexImportAccountFromJWT(item, item.IDToken, false, now)
placeholder
	if err := enrichCodexImportAccountFromJWT(item, item.AccessToken, true, now); err != nil {
		return nil, err
placeholder
	if _, ok := item.Credentials["expires_at"]; !ok {
		item.WarningTexts = append(item.WarningTexts, "无法从 accessToken 解析过期时间，导入后需自行确认令牌有效性")
placeholder
	if item.RefreshToken == "" {
		item.WarningTexts = append(item.WarningTexts, "未包含 refresh_token，accessToken 过期后无法自动续期")
placeholder

	setCodexCredentialIfNotEmpty(item.Credentials, "email", item.Email)
	setCodexCredentialIfNotEmpty(item.Credentials, "chatgpt_account_id", item.AccountID)
	setCodexCredentialIfNotEmpty(item.Credentials, "chatgpt_user_id", item.UserID)
	setCodexCredentialIfNotEmpty(item.Credentials, "organization_id", item.Organization)
	setCodexCredentialIfNotEmpty(item.Credentials, "plan_type", item.PlanType)

	fingerprint := codexTokenFingerprint(item.AccessToken)
	item.Extra["access_token_sha256"] = fingerprint
	item.IdentityKeys = buildCodexImportIdentityKeys(item.AccountID, item.UserID, item.Email, item.AccessToken, item.RefreshToken)
	item.Name = buildCodexImportAccountName(item, entry.Index)

	return item, nil
placeholder

func enrichCodexImportAccountFromJWT(item *codexImportAccount, token string, validateExpiry bool, now time.Time) error {
	claims, err := decodeCodexJWTClaims(token)
	if err != nil {
		if validateExpiry {
			item.WarningTexts = append(item.WarningTexts, "accessToken 不是可解析 JWT，无法校验过期时间和账号身份")
	placeholder
		return nil
placeholder
	if validateExpiry && claims.Exp > 0 {
		if now.Unix() > claims.Exp+codexImportClockSkewSeconds {
			return fmt.Errorf("access_token 已过期: %s", time.Unix(claims.Exp, 0).UTC().Format(time.RFC3339))
	placeholder
		expiresAt := time.Unix(claims.Exp, 0).UTC()
		item.TokenExpiresAt = &expiresAt
		item.Credentials["expires_at"] = expiresAt.Format(time.RFC3339)
placeholder
	if item.Email == "" {
		item.Email = strings.TrimSpace(claims.Email)
placeholder
	if claims.OpenAIAuth == nil {
		if item.UserID == "" {
			item.UserID = strings.TrimSpace(claims.Sub)
	placeholder
		return nil
placeholder
	if item.AccountID == "" {
		item.AccountID = strings.TrimSpace(claims.OpenAIAuth.ChatGPTAccountID)
placeholder
	if item.UserID == "" {
		item.UserID = strings.TrimSpace(claims.OpenAIAuth.ChatGPTUserID)
placeholder
	if item.UserID == "" {
		item.UserID = strings.TrimSpace(claims.OpenAIAuth.UserID)
placeholder
	if item.PlanType == "" {
		item.PlanType = strings.TrimSpace(claims.OpenAIAuth.ChatGPTPlanType)
placeholder
	if item.Organization == "" {
		item.Organization = strings.TrimSpace(claims.OpenAIAuth.POID)
placeholder
	if item.Organization == "" {
		for _, org := range claims.OpenAIAuth.Organizations {
			if org.IsDefault {
				item.Organization = org.ID
				break
		placeholder
	placeholder
placeholder
	if item.Organization == "" && len(claims.OpenAIAuth.Organizations) > 0 {
		item.Organization = claims.OpenAIAuth.Organizations[0].ID
placeholder
	if item.UserID == "" {
		item.UserID = strings.TrimSpace(claims.Sub)
placeholder
	return nil
placeholder

func decodeCodexJWTClaims(token string) (*codexJWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
placeholder
	payload, err := decodeCodexJWTSegment(parts[1])
	if err != nil {
		return nil, err
placeholder
	var claims codexJWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
placeholder
	return &claims, nil
placeholder

func decodeCodexJWTSegment(segment string) ([]byte, error) {
	if decoded, err := base64.RawURLEncoding.DecodeString(segment); err == nil {
		return decoded, nil
placeholder
	if decoded, err := base64.RawStdEncoding.DecodeString(segment); err == nil {
		return decoded, nil
placeholder
	padded := segment
	if rem := len(padded) % 4; rem > 0 {
		padded += strings.Repeat("=", 4-rem)
placeholder
	if decoded, err := base64.URLEncoding.DecodeString(padded); err == nil {
		return decoded, nil
placeholder
	return base64.StdEncoding.DecodeString(padded)
placeholder

func buildCodexImportAccountName(item *codexImportAccount, index int) string {
	for _, candidate := range []string{item.Name, item.Email, item.AccountID, item.UserIDplaceholder {
		candidate = strings.TrimSpace(candidate)
		if candidate != "" {
			return candidate
	placeholder
placeholder
	return fmt.Sprintf("Codex 导入账号 %d", index)
placeholder

func buildCodexCreateAccountName(base string, item *codexImportAccount, index, total int) string {
	base = strings.TrimSpace(base)
	if base == "" {
		if item == nil {
			return fmt.Sprintf("Codex 导入账号 %d", index)
	placeholder
		return item.Name
placeholder
	if total > 1 {
		return fmt.Sprintf("%s #%d", base, index)
placeholder
	return base
placeholder

func resolveCodexImportExpiry(req CodexSessionImportRequest, item *codexImportAccount) (*int64, *time.Time, *bool, []string, error) {
	if item == nil {
		return nil, nil, nil, nil, errors.New("导入项为空")
placeholder
	// Agent Identity has no OAuth access-token lifetime. Its runtime/task
	// lifecycle is handled by the upstream task recovery path, so it must not
	// be rejected or auto-paused by the OAuth import expiry policy.
	if item.IsAgentIdentity {
		return nil, nil, nil, nil, nil
placeholder

	var requestExpiresAt *time.Time
	if req.ExpiresAt != nil && *req.ExpiresAt > 0 {
		t := time.Unix(*req.ExpiresAt, 0).UTC()
		requestExpiresAt = &t
placeholder

	var accountExpiresAt *time.Time
	var credentialExpiresAt *time.Time
	warnings := make([]string, 0, 2)
	if item.RefreshToken == "" {
		if item.TokenExpiresAt != nil {
			tokenExpiresAt := item.TokenExpiresAt.UTC()
			accountExpiresAt = &tokenExpiresAt
			credentialExpiresAt = &tokenExpiresAt
	placeholder
		if requestExpiresAt != nil {
			accountExpiresAt = earlierCodexTime(accountExpiresAt, requestExpiresAt)
			credentialExpiresAt = earlierCodexTime(credentialExpiresAt, requestExpiresAt)
	placeholder
		if accountExpiresAt == nil {
			return nil, nil, nil, nil, errors.New("未包含 refresh_token，且无法解析 accessToken 过期时间；请在第一步设置过期时间后再导入")
	placeholder
		if accountExpiresAt.Unix() <= time.Now().UTC().Unix()-codexImportClockSkewSeconds {
			return nil, nil, nil, nil, fmt.Errorf("过期时间已过期: %s", accountExpiresAt.Format(time.RFC3339))
	placeholder
		warnings = append(warnings, "未包含 refresh_token，已按 accessToken/账号过期时间设置自动停止调度")
		if req.AutoPauseOnExpired != nil && !*req.AutoPauseOnExpired {
			warnings = append(warnings, "未包含 refresh_token，已强制开启过期自动暂停")
	placeholder
		autoPause := true
		expiresAtUnix := accountExpiresAt.Unix()
		return &expiresAtUnix, credentialExpiresAt, &autoPause, warnings, nil
placeholder

	if requestExpiresAt != nil {
		accountExpiresAt = requestExpiresAt
placeholder
	if item.TokenExpiresAt != nil {
		tokenExpiresAt := item.TokenExpiresAt.UTC()
		credentialExpiresAt = &tokenExpiresAt
placeholder
	var expiresAtUnix *int64
	if accountExpiresAt != nil {
		v := accountExpiresAt.Unix()
		expiresAtUnix = &v
placeholder
	return expiresAtUnix, credentialExpiresAt, req.AutoPauseOnExpired, warnings, nil
placeholder

func earlierCodexTime(current, candidate *time.Time) *time.Time {
	if candidate == nil {
		return current
placeholder
	if current == nil || candidate.Before(*current) {
		t := candidate.UTC()
		return &t
placeholder
	t := current.UTC()
	return &t
placeholder

func sanitizeCodexImportCredentialExtras(input map[string]any) map[string]any {
	if len(input) == 0 {
		return nil
placeholder
	protected := map[string]struct{placeholder{
		"access_token":               {placeholder,
		"refresh_token":              {placeholder,
		"id_token":                   {placeholder,
		"expires_at":                 {placeholder,
		"email":                      {placeholder,
		"chatgpt_account_id":         {placeholder,
		"chatgpt_user_id":            {placeholder,
		"organization_id":            {placeholder,
		"plan_type":                  {placeholder,
		"client_id":                  {placeholder,
		"auth_mode":                  {placeholder,
		"openai_auth_mode":           {placeholder,
		"token_type":                 {placeholder,
		"chatgpt_account_is_fedramp": {placeholder,
		"agent_runtime_id":           {placeholder,
		"agent_private_key":          {placeholder,
		"task_id":                    {placeholder,
placeholder
	out := make(map[string]any, len(input))
	for key, value := range input {
		normalizedKey := strings.TrimSpace(key)
		if normalizedKey == "" {
			continue
	placeholder
		if _, ok := protected[strings.ToLower(normalizedKey)]; ok {
			continue
	placeholder
		out[normalizedKey] = value
placeholder
	if len(out) == 0 {
		return nil
placeholder
	return out
placeholder

// buildCodexImportIdentityKeys 生成导入条目的匹配键。refresh_token 缺失时
// Codex session 只能作为 accessToken-only 凭据使用，此时以 access token
// 指纹作为唯一稳定身份，避免同 workspace 下共享的 account/user 标识误合并。
func buildCodexImportIdentityKeys(accountID, userID, email, accessToken, refreshToken string) []string {
	accessToken = strings.TrimSpace(accessToken)
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" && accessToken != "" {
		return []string{"access:" + codexTokenFingerprint(accessToken)placeholder
placeholder
	return buildCodexStoredIdentityKeys(accountID, userID, email, accessToken)
placeholder

func buildCodexAgentIdentityKeys(accountID, userID, email, runtimeID string) []string {
	keys := buildCodexStoredIdentityKeys(accountID, userID, email, "")
	if runtimeID = strings.TrimSpace(runtimeID); runtimeID != "" {
		keys = append([]string{"agent:" + runtimeIDplaceholder, keys...)
placeholder
	return keys
placeholder

// buildCodexStoredIdentityKeys 生成存量账号索引键，保留 user/account 维度，
// 让 accessToken-only 账号后续升级为完整 OAuth 时仍能命中并更新原账号。
func buildCodexStoredIdentityKeys(accountID, userID, email, accessToken string) []string {
	keys := make([]string, 0, 3)
	accountID = strings.TrimSpace(accountID)
	userID = strings.TrimSpace(userID)
	accessToken = strings.TrimSpace(accessToken)
	if userID != "" {
		keys = append(keys, "user:"+userID)
placeholder
	if accountID == "" && userID == "" {
		if email = strings.ToLower(strings.TrimSpace(email)); email != "" {
			keys = append(keys, "email:"+email)
	placeholder
placeholder
	if accessToken != "" {
		keys = append(keys, "access:"+codexTokenFingerprint(accessToken))
placeholder
	if accountID != "" {
		keys = append(keys, "account:"+accountID)
placeholder
	return keys
placeholder

func buildCodexAccountIndex(accounts []service.Account) *codexAccountIndex {
	index := &codexAccountIndex{accountsByKey: map[string][]service.Account{placeholderplaceholder
	for _, account := range accounts {
		index.Add(account)
placeholder
	return index
placeholder

func (i *codexAccountIndex) Add(account service.Account) {
	if i == nil {
		return
placeholder
	if i.accountsByKey == nil {
		i.accountsByKey = map[string][]service.Account{placeholder
placeholder
	i.remove(account.ID)
	keys := buildCodexStoredIdentityKeys(
		codexCredentialString(account.Credentials, "chatgpt_account_id"),
		codexCredentialString(account.Credentials, "chatgpt_user_id"),
		codexCredentialString(account.Credentials, "email"),
		codexCredentialString(account.Credentials, "access_token"),
	)
	for _, key := range keys {
		i.accountsByKey[key] = upsertCodexAccount(i.accountsByKey[key], account)
placeholder
	if runtimeID := codexCredentialString(account.Credentials, "agent_runtime_id"); runtimeID != "" {
		key := "agent:" + runtimeID
		i.accountsByKey[key] = upsertCodexAccount(i.accountsByKey[key], account)
placeholder
placeholder

func (i *codexAccountIndex) remove(accountID int64) {
	for key, accounts := range i.accountsByKey {
		kept := accounts[:0]
		for _, account := range accounts {
			if account.ID != accountID {
				kept = append(kept, account)
		placeholder
	placeholder
		if len(kept) == 0 {
			delete(i.accountsByKey, key)
			continue
	placeholder
		i.accountsByKey[key] = kept
placeholder
placeholder

// upsertCodexAccount 保留同一键下的全部候选账号（共享的 account: 键可对应
// 团队内多个账号），同一账号重复 Add 时原位替换为最新状态。
func upsertCodexAccount(accounts []service.Account, account service.Account) []service.Account {
	for idx := range accounts {
		if accounts[idx].ID == account.ID {
			accounts[idx] = account
			return accounts
	placeholder
placeholder
	return append(accounts, account)
placeholder

// Find 返回第一个通过跨用户校验的候选账号及其命中的匹配键。
func (i *codexAccountIndex) Find(keys []string, userID string) (*service.Account, string) {
	if i == nil {
		return nil, ""
placeholder
	for _, key := range keys {
		for _, account := range i.accountsByKey[key] {
			if codexIdentityConflicts(key, userID, codexCredentialString(account.Credentials, "chatgpt_user_id")) {
				continue
		placeholder
			return &account, key
	placeholder
placeholder
	return nil, ""
placeholder

// codexIdentityConflicts 判断 account: 键的命中是否把同一 ChatGPT 团队的两个
// 不同成员误连到一起：双方都携带 user id 且不相等时视为冲突。存量索引侧
// 仍保留 account 键，任一侧缺少 user id 时允许匹配，使含 refresh_token
// 的常规导入和 accessToken-only 账号升级为完整 OAuth 时仍能更新原账号。
func codexIdentityConflicts(key, userID, storedUserID string) bool {
	if !strings.HasPrefix(key, "account:") {
		return false
placeholder
	userID = strings.TrimSpace(userID)
	storedUserID = strings.TrimSpace(storedUserID)
	return userID != "" && storedUserID != "" && userID != storedUserID
placeholder

type codexSeenIdentity struct {
	index  int
	userID string
placeholder

func firstSeenCodexIdentity(seen map[string]codexSeenIdentity, keys []string, userID string) (int, bool) {
	for _, key := range keys {
		entry, ok := seen[key]
		if !ok {
			continue
	placeholder
		if codexIdentityConflicts(key, userID, entry.userID) {
			continue
	placeholder
		return entry.index, true
placeholder
	return 0, false
placeholder

func markCodexIdentitySeen(seen map[string]codexSeenIdentity, keys []string, index int, userID string) {
	for _, key := range keys {
		seen[key] = codexSeenIdentity{index: index, userID: userIDplaceholder
placeholder
placeholder

func mergeCodexImportMap(existing, incoming map[string]any) map[string]any {
	out := make(map[string]any, len(existing)+len(incoming))
	for k, v := range existing {
		out[k] = v
placeholder
	for k, v := range incoming {
		out[k] = v
placeholder
	return out
placeholder

func mergeCodexImportCredentials(existing, incoming map[string]any, item *codexImportAccount) map[string]any {
	out := mergeCodexImportMap(existing, incoming)
	if item == nil {
		return out
placeholder
	if strings.TrimSpace(item.RefreshToken) == "" {
		if codexCredentialString(existing, "refresh_token") == "" {
			delete(out, "refresh_token")
			delete(out, "client_id")
	placeholder else {
			out["refresh_token"] = existing["refresh_token"]
			if clientID, ok := existing["client_id"]; ok {
				out["client_id"] = clientID
		placeholder
	placeholder
placeholder
	if strings.TrimSpace(item.IDToken) == "" {
		delete(out, "id_token")
placeholder
	return out
placeholder

func codexCredentialString(credentials map[string]any, key string) string {
	if credentials == nil {
		return ""
placeholder
	return codexStringValue(credentials[key])
placeholder

func codexTokenFingerprint(token string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(token)))
	return hex.EncodeToString(sum[:])
placeholder

func looksLikeJSON(content string) bool {
	if content == "" {
		return false
placeholder
	switch content[0] {
	case '{', '[':
		return true
	default:
		return false
placeholder
placeholder

func firstCodexString(obj map[string]any, paths ...[]string) string {
	for _, path := range paths {
		if value, ok := codexPathValue(obj, path); ok {
			if str := codexStringValue(value); str != "" {
				return str
		placeholder
	placeholder
placeholder
	return ""
placeholder

func firstCodexMap(obj map[string]any, paths ...[]string) (map[string]any, bool) {
	for _, path := range paths {
		value, ok := codexPathValue(obj, path)
		if !ok || value == nil {
			continue
	placeholder
		if mapped, ok := value.(map[string]any); ok {
			return mapped, true
	placeholder
placeholder
	return nil, false
placeholder

func firstCodexBool(obj map[string]any, paths ...[]string) bool {
	for _, path := range paths {
		value, ok := codexPathValue(obj, path)
		if !ok {
			continue
	placeholder
		switch value := value.(type) {
		case bool:
			return value
		case string:
			parsed, err := strconv.ParseBool(strings.TrimSpace(value))
			if err == nil {
				return parsed
		placeholder
	placeholder
placeholder
	return false
placeholder

func copyCodexExtraString(obj map[string]any, extra map[string]any, key string, path []string) {
	value := firstCodexString(obj, path)
	if value != "" {
		extra[key] = value
placeholder
placeholder

func firstCodexTime(obj map[string]any, paths ...[]string) (time.Time, bool) {
	for _, path := range paths {
		if value, ok := codexTimeAt(obj, path); ok {
			return value, true
	placeholder
placeholder
	return time.Time{placeholder, false
placeholder

func codexTimeAt(obj map[string]any, path []string) (time.Time, bool) {
	value, ok := codexPathValue(obj, path)
	if !ok {
		return time.Time{placeholder, false
placeholder
	return parseCodexTimeValue(value)
placeholder

func codexPathValue(obj map[string]any, path []string) (any, bool) {
	var current any = obj
	for _, key := range path {
		currentObj, ok := current.(map[string]any)
		if !ok {
			return nil, false
	placeholder
		value, ok := currentObj[key]
		if !ok {
			return nil, false
	placeholder
		current = value
placeholder
	return current, true
placeholder

func codexStringValue(value any) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	case json.Number:
		return strings.TrimSpace(v.String())
	case float64:
		return strings.TrimSpace(strconv.FormatFloat(v, 'f', -1, 64))
	case float32:
		return strings.TrimSpace(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
placeholder
placeholder

func setCodexCredentialIfNotEmpty(credentials map[string]any, key, value string) {
	value = strings.TrimSpace(value)
	if value != "" {
		credentials[key] = value
placeholder
placeholder

func parseCodexTimeValue(value any) (time.Time, bool) {
	switch v := value.(type) {
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return time.Time{placeholder, false
	placeholder
		if parsed, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return parsed.UTC(), true
	placeholder
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return codexUnixTime(n), true
	placeholder
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return codexUnixTime(n), true
	placeholder
		if f, err := v.Float64(); err == nil {
			return codexUnixTime(int64(f)), true
	placeholder
	case float64:
		return codexUnixTime(int64(v)), true
	case int:
		return codexUnixTime(int64(v)), true
	case int64:
		return codexUnixTime(v), true
placeholder
	return time.Time{placeholder, false
placeholder

func codexUnixTime(value int64) time.Time {
	if value > 1_000_000_000_000 {
		return time.UnixMilli(value).UTC()
placeholder
	return time.Unix(value, 0).UTC()
placeholder
