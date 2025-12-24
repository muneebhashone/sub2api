package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/service/ports"
)

type CRSSyncService struct {
	accountRepo ports.AccountRepository
	proxyRepo   ports.ProxyRepository
placeholder

func NewCRSSyncService(accountRepo ports.AccountRepository, proxyRepo ports.ProxyRepository) *CRSSyncService {
	return &CRSSyncService{
		accountRepo: accountRepo,
		proxyRepo:   proxyRepo,
placeholder
placeholder

type SyncFromCRSInput struct {
	BaseURL     string
	Username    string
	Password    string
	SyncProxies bool
placeholder

type SyncFromCRSItemResult struct {
	CRSAccountID string `json:"crs_account_id"`
	Kind         string `json:"kind"`
	Name         string `json:"name"`
	Action       string `json:"action"` // created/updated/failed/skipped
	Error        string `json:"error,omitempty"`
placeholder

type SyncFromCRSResult struct {
	Created int                     `json:"created"`
	Updated int                     `json:"updated"`
	Skipped int                     `json:"skipped"`
	Failed  int                     `json:"failed"`
	Items   []SyncFromCRSItemResult `json:"items"`
placeholder

type crsLoginResponse struct {
	Success  bool   `json:"success"`
	Token    string `json:"token"`
	Message  string `json:"message"`
	Error    string `json:"error"`
	Username string `json:"username"`
placeholder

type crsExportResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
	Data    struct {
		ExportedAt              string                      `json:"exportedAt"`
		ClaudeAccounts          []crsClaudeAccount          `json:"claudeAccounts"`
		ClaudeConsoleAccounts   []crsConsoleAccount         `json:"claudeConsoleAccounts"`
		OpenAIOAuthAccounts     []crsOpenAIOAuthAccount     `json:"openaiOAuthAccounts"`
		OpenAIResponsesAccounts []crsOpenAIResponsesAccount `json:"openaiResponsesAccounts"`
placeholder `json:"data"`
placeholder

type crsProxy struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
placeholder

type crsClaudeAccount struct {
	Kind        string         `json:"kind"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Platform    string         `json:"platform"`
	AuthType    string         `json:"authType"` // oauth/setup-token
	IsActive    bool           `json:"isActive"`
	Schedulable bool           `json:"schedulable"`
	Priority    int            `json:"priority"`
	Status      string         `json:"status"`
	Proxy       *crsProxy      `json:"proxy"`
	Credentials map[string]any `json:"credentials"`
placeholder

type crsConsoleAccount struct {
	Kind               string         `json:"kind"`
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	Platform           string         `json:"platform"`
	IsActive           bool           `json:"isActive"`
	Schedulable        bool           `json:"schedulable"`
	Priority           int            `json:"priority"`
	Status             string         `json:"status"`
	MaxConcurrentTasks int            `json:"maxConcurrentTasks"`
	Proxy              *crsProxy      `json:"proxy"`
	Credentials        map[string]any `json:"credentials"`
placeholder

type crsOpenAIResponsesAccount struct {
	Kind        string         `json:"kind"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Platform    string         `json:"platform"`
	IsActive    bool           `json:"isActive"`
	Schedulable bool           `json:"schedulable"`
	Priority    int            `json:"priority"`
	Status      string         `json:"status"`
	Proxy       *crsProxy      `json:"proxy"`
	Credentials map[string]any `json:"credentials"`
placeholder

type crsOpenAIOAuthAccount struct {
	Kind        string         `json:"kind"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Platform    string         `json:"platform"`
	AuthType    string         `json:"authType"` // oauth
	IsActive    bool           `json:"isActive"`
	Schedulable bool           `json:"schedulable"`
	Priority    int            `json:"priority"`
	Status      string         `json:"status"`
	Proxy       *crsProxy      `json:"proxy"`
	Credentials map[string]any `json:"credentials"`
placeholder

func (s *CRSSyncService) SyncFromCRS(ctx context.Context, input SyncFromCRSInput) (*SyncFromCRSResult, error) {
	baseURL, err := normalizeBaseURL(input.BaseURL)
	if err != nil {
		return nil, err
placeholder
	if strings.TrimSpace(input.Username) == "" || strings.TrimSpace(input.Password) == "" {
		return nil, errors.New("username and password are required")
placeholder

	client := &http.Client{Timeout: 20 * time.Secondplaceholder

	adminToken, err := crsLogin(ctx, client, baseURL, input.Username, input.Password)
	if err != nil {
		return nil, err
placeholder

	exported, err := crsExportAccounts(ctx, client, baseURL, adminToken)
	if err != nil {
		return nil, err
placeholder

	now := time.Now().UTC().Format(time.RFC3339)

	result := &SyncFromCRSResult{
		Items: make(
			[]SyncFromCRSItemResult,
			0,
			len(exported.Data.ClaudeAccounts)+len(exported.Data.ClaudeConsoleAccounts)+len(exported.Data.OpenAIOAuthAccounts)+len(exported.Data.OpenAIResponsesAccounts),
		),
placeholder

	var proxies []model.Proxy
	if input.SyncProxies {
		proxies, _ = s.proxyRepo.ListActive(ctx)
placeholder

	// Claude OAuth / Setup Token -> sub2api anthropic oauth/setup-token
	for _, src := range exported.Data.ClaudeAccounts {
		item := SyncFromCRSItemResult{
			CRSAccountID: src.ID,
			Kind:         src.Kind,
			Name:         src.Name,
	placeholder

		targetType := strings.TrimSpace(src.AuthType)
		if targetType == "" {
			targetType = "oauth"
	placeholder
		if targetType != model.AccountTypeOAuth && targetType != model.AccountTypeSetupToken {
			item.Action = "skipped"
			item.Error = "unsupported authType: " + targetType
			result.Skipped++
			result.Items = append(result.Items, item)
			continue
	placeholder

		accessToken, _ := src.Credentials["access_token"].(string)
		if strings.TrimSpace(accessToken) == "" {
			item.Action = "failed"
			item.Error = "missing access_token"
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		proxyID, err := s.mapOrCreateProxy(ctx, input.SyncProxies, &proxies, src.Proxy, fmt.Sprintf("crs-%s", src.Name))
		if err != nil {
			item.Action = "failed"
			item.Error = "proxy sync failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		credentials := sanitizeCredentialsMap(src.Credentials)
		priority := clampPriority(src.Priority)
		concurrency := 3
		status := mapCRSStatus(src.IsActive, src.Status)

		extra := map[string]any{
			"crs_account_id": src.ID,
			"crs_kind":       src.Kind,
			"crs_synced_at":  now,
	placeholder

		existing, err := s.accountRepo.GetByCRSAccountID(ctx, src.ID)
		if err != nil {
			item.Action = "failed"
			item.Error = "db lookup failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing == nil {
			account := &model.Account{
				Name:        defaultName(src.Name, src.ID),
				Platform:    model.PlatformAnthropic,
				Type:        targetType,
				Credentials: model.JSONB(credentials),
				Extra:       model.JSONB(extra),
				ProxyID:     proxyID,
				Concurrency: concurrency,
				Priority:    priority,
				Status:      status,
				Schedulable: src.Schedulable,
		placeholder
			if err := s.accountRepo.Create(ctx, account); err != nil {
				item.Action = "failed"
				item.Error = "create failed: " + err.Error()
				result.Failed++
				result.Items = append(result.Items, item)
				continue
		placeholder
			item.Action = "created"
			result.Created++
			result.Items = append(result.Items, item)
			continue
	placeholder

		// Update existing
		if existing.Extra == nil {
			existing.Extra = make(model.JSONB)
	placeholder
		for k, v := range extra {
			existing.Extra[k] = v
	placeholder
		existing.Name = defaultName(src.Name, src.ID)
		existing.Platform = model.PlatformAnthropic
		existing.Type = targetType
		existing.Credentials = model.JSONB(credentials)
		existing.ProxyID = proxyID
		existing.Concurrency = concurrency
		existing.Priority = priority
		existing.Status = status
		existing.Schedulable = src.Schedulable

		if err := s.accountRepo.Update(ctx, existing); err != nil {
			item.Action = "failed"
			item.Error = "update failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		item.Action = "updated"
		result.Updated++
		result.Items = append(result.Items, item)
placeholder

	// Claude Console API Key -> sub2api anthropic apikey
	for _, src := range exported.Data.ClaudeConsoleAccounts {
		item := SyncFromCRSItemResult{
			CRSAccountID: src.ID,
			Kind:         src.Kind,
			Name:         src.Name,
	placeholder

		apiKey, _ := src.Credentials["api_key"].(string)
		if strings.TrimSpace(apiKey) == "" {
			item.Action = "failed"
			item.Error = "missing api_key"
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		proxyID, err := s.mapOrCreateProxy(ctx, input.SyncProxies, &proxies, src.Proxy, fmt.Sprintf("crs-%s", src.Name))
		if err != nil {
			item.Action = "failed"
			item.Error = "proxy sync failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		credentials := sanitizeCredentialsMap(src.Credentials)
		priority := clampPriority(src.Priority)
		concurrency := 3
		if src.MaxConcurrentTasks > 0 {
			concurrency = src.MaxConcurrentTasks
	placeholder
		status := mapCRSStatus(src.IsActive, src.Status)

		extra := map[string]any{
			"crs_account_id": src.ID,
			"crs_kind":       src.Kind,
			"crs_synced_at":  now,
	placeholder

		existing, err := s.accountRepo.GetByCRSAccountID(ctx, src.ID)
		if err != nil {
			item.Action = "failed"
			item.Error = "db lookup failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing == nil {
			account := &model.Account{
				Name:        defaultName(src.Name, src.ID),
				Platform:    model.PlatformAnthropic,
				Type:        model.AccountTypeApiKey,
				Credentials: model.JSONB(credentials),
				Extra:       model.JSONB(extra),
				ProxyID:     proxyID,
				Concurrency: concurrency,
				Priority:    priority,
				Status:      status,
				Schedulable: src.Schedulable,
		placeholder
			if err := s.accountRepo.Create(ctx, account); err != nil {
				item.Action = "failed"
				item.Error = "create failed: " + err.Error()
				result.Failed++
				result.Items = append(result.Items, item)
				continue
		placeholder
			item.Action = "created"
			result.Created++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing.Extra == nil {
			existing.Extra = make(model.JSONB)
	placeholder
		for k, v := range extra {
			existing.Extra[k] = v
	placeholder
		existing.Name = defaultName(src.Name, src.ID)
		existing.Platform = model.PlatformAnthropic
		existing.Type = model.AccountTypeApiKey
		existing.Credentials = model.JSONB(credentials)
		existing.ProxyID = proxyID
		existing.Concurrency = concurrency
		existing.Priority = priority
		existing.Status = status
		existing.Schedulable = src.Schedulable

		if err := s.accountRepo.Update(ctx, existing); err != nil {
			item.Action = "failed"
			item.Error = "update failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		item.Action = "updated"
		result.Updated++
		result.Items = append(result.Items, item)
placeholder

	// OpenAI OAuth -> sub2api openai oauth
	for _, src := range exported.Data.OpenAIOAuthAccounts {
		item := SyncFromCRSItemResult{
			CRSAccountID: src.ID,
			Kind:         src.Kind,
			Name:         src.Name,
	placeholder

		accessToken, _ := src.Credentials["access_token"].(string)
		if strings.TrimSpace(accessToken) == "" {
			item.Action = "failed"
			item.Error = "missing access_token"
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		proxyID, err := s.mapOrCreateProxy(
			ctx,
			input.SyncProxies,
			&proxies,
			src.Proxy,
			fmt.Sprintf("crs-%s", src.Name),
		)
		if err != nil {
			item.Action = "failed"
			item.Error = "proxy sync failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		credentials := sanitizeCredentialsMap(src.Credentials)
		// Normalize token_type
		if v, ok := credentials["token_type"].(string); !ok || strings.TrimSpace(v) == "" {
			credentials["token_type"] = "Bearer"
	placeholder
		priority := clampPriority(src.Priority)
		concurrency := 3
		status := mapCRSStatus(src.IsActive, src.Status)

		extra := map[string]any{
			"crs_account_id": src.ID,
			"crs_kind":       src.Kind,
			"crs_synced_at":  now,
	placeholder

		existing, err := s.accountRepo.GetByCRSAccountID(ctx, src.ID)
		if err != nil {
			item.Action = "failed"
			item.Error = "db lookup failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing == nil {
			account := &model.Account{
				Name:        defaultName(src.Name, src.ID),
				Platform:    model.PlatformOpenAI,
				Type:        model.AccountTypeOAuth,
				Credentials: model.JSONB(credentials),
				Extra:       model.JSONB(extra),
				ProxyID:     proxyID,
				Concurrency: concurrency,
				Priority:    priority,
				Status:      status,
				Schedulable: src.Schedulable,
		placeholder
			if err := s.accountRepo.Create(ctx, account); err != nil {
				item.Action = "failed"
				item.Error = "create failed: " + err.Error()
				result.Failed++
				result.Items = append(result.Items, item)
				continue
		placeholder
			item.Action = "created"
			result.Created++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing.Extra == nil {
			existing.Extra = make(model.JSONB)
	placeholder
		for k, v := range extra {
			existing.Extra[k] = v
	placeholder
		existing.Name = defaultName(src.Name, src.ID)
		existing.Platform = model.PlatformOpenAI
		existing.Type = model.AccountTypeOAuth
		existing.Credentials = model.JSONB(credentials)
		existing.ProxyID = proxyID
		existing.Concurrency = concurrency
		existing.Priority = priority
		existing.Status = status
		existing.Schedulable = src.Schedulable

		if err := s.accountRepo.Update(ctx, existing); err != nil {
			item.Action = "failed"
			item.Error = "update failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		item.Action = "updated"
		result.Updated++
		result.Items = append(result.Items, item)
placeholder

	// OpenAI Responses API Key -> sub2api openai apikey
	for _, src := range exported.Data.OpenAIResponsesAccounts {
		item := SyncFromCRSItemResult{
			CRSAccountID: src.ID,
			Kind:         src.Kind,
			Name:         src.Name,
	placeholder

		apiKey, _ := src.Credentials["api_key"].(string)
		if strings.TrimSpace(apiKey) == "" {
			item.Action = "failed"
			item.Error = "missing api_key"
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if baseURL, ok := src.Credentials["base_url"].(string); !ok || strings.TrimSpace(baseURL) == "" {
			src.Credentials["base_url"] = "https://api.openai.com"
	placeholder

		proxyID, err := s.mapOrCreateProxy(
			ctx,
			input.SyncProxies,
			&proxies,
			src.Proxy,
			fmt.Sprintf("crs-%s", src.Name),
		)
		if err != nil {
			item.Action = "failed"
			item.Error = "proxy sync failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		credentials := sanitizeCredentialsMap(src.Credentials)
		priority := clampPriority(src.Priority)
		concurrency := 3
		status := mapCRSStatus(src.IsActive, src.Status)

		extra := map[string]any{
			"crs_account_id": src.ID,
			"crs_kind":       src.Kind,
			"crs_synced_at":  now,
	placeholder

		existing, err := s.accountRepo.GetByCRSAccountID(ctx, src.ID)
		if err != nil {
			item.Action = "failed"
			item.Error = "db lookup failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing == nil {
			account := &model.Account{
				Name:        defaultName(src.Name, src.ID),
				Platform:    model.PlatformOpenAI,
				Type:        model.AccountTypeApiKey,
				Credentials: model.JSONB(credentials),
				Extra:       model.JSONB(extra),
				ProxyID:     proxyID,
				Concurrency: concurrency,
				Priority:    priority,
				Status:      status,
				Schedulable: src.Schedulable,
		placeholder
			if err := s.accountRepo.Create(ctx, account); err != nil {
				item.Action = "failed"
				item.Error = "create failed: " + err.Error()
				result.Failed++
				result.Items = append(result.Items, item)
				continue
		placeholder
			item.Action = "created"
			result.Created++
			result.Items = append(result.Items, item)
			continue
	placeholder

		if existing.Extra == nil {
			existing.Extra = make(model.JSONB)
	placeholder
		for k, v := range extra {
			existing.Extra[k] = v
	placeholder
		existing.Name = defaultName(src.Name, src.ID)
		existing.Platform = model.PlatformOpenAI
		existing.Type = model.AccountTypeApiKey
		existing.Credentials = model.JSONB(credentials)
		existing.ProxyID = proxyID
		existing.Concurrency = concurrency
		existing.Priority = priority
		existing.Status = status
		existing.Schedulable = src.Schedulable

		if err := s.accountRepo.Update(ctx, existing); err != nil {
			item.Action = "failed"
			item.Error = "update failed: " + err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
	placeholder

		item.Action = "updated"
		result.Updated++
		result.Items = append(result.Items, item)
placeholder

	return result, nil
placeholder

func (s *CRSSyncService) mapOrCreateProxy(ctx context.Context, enabled bool, cached *[]model.Proxy, src *crsProxy, defaultName string) (*int64, error) {
	if !enabled || src == nil {
		return nil, nil
placeholder
	protocol := strings.ToLower(strings.TrimSpace(src.Protocol))
	switch protocol {
	case "socks":
		protocol = "socks5"
	case "socks5h":
		protocol = "socks5"
placeholder
	host := strings.TrimSpace(src.Host)
	port := src.Port
	username := strings.TrimSpace(src.Username)
	password := strings.TrimSpace(src.Password)

	if protocol == "" || host == "" || port <= 0 {
		return nil, nil
placeholder
	if protocol != "http" && protocol != "https" && protocol != "socks5" {
		return nil, nil
placeholder

	// Find existing proxy (active only).
	for _, p := range *cached {
		if strings.EqualFold(p.Protocol, protocol) &&
			p.Host == host &&
			p.Port == port &&
			p.Username == username &&
			p.Password == password {
			id := p.ID
			return &id, nil
	placeholder
placeholder

	// Create new proxy
	proxy := &model.Proxy{
		Name:     defaultProxyName(defaultName, protocol, host, port),
		Protocol: protocol,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Status:   model.StatusActive,
placeholder
	if err := s.proxyRepo.Create(ctx, proxy); err != nil {
		return nil, err
placeholder

	*cached = append(*cached, *proxy)
	id := proxy.ID
	return &id, nil
placeholder

func defaultProxyName(base, protocol, host string, port int) string {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "crs"
placeholder
	return fmt.Sprintf("%s (%s://%s:%d)", base, protocol, host, port)
placeholder

func defaultName(name, id string) string {
	if strings.TrimSpace(name) != "" {
		return strings.TrimSpace(name)
placeholder
	return "CRS " + id
placeholder

func clampPriority(priority int) int {
	if priority < 1 || priority > 100 {
		return 50
placeholder
	return priority
placeholder

func sanitizeCredentialsMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{placeholder
placeholder
	out := make(map[string]any, len(input))
	for k, v := range input {
		// Avoid nil values to keep JSONB cleaner
		if v != nil {
			out[k] = v
	placeholder
placeholder
	return out
placeholder

func mapCRSStatus(isActive bool, status string) string {
	if !isActive {
		return "inactive"
placeholder
	if strings.EqualFold(strings.TrimSpace(status), "error") {
		return "error"
placeholder
	return "active"
placeholder

func normalizeBaseURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", errors.New("base_url is required")
placeholder
	u, err := url.Parse(trimmed)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid base_url: %s", trimmed)
placeholder
	u.Path = strings.TrimRight(u.Path, "/")
	return strings.TrimRight(u.String(), "/"), nil
placeholder

func crsLogin(ctx context.Context, client *http.Client, baseURL, username, password string) (string, error) {
	payload := map[string]any{
		"username": username,
		"password": password,
placeholder
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/web/auth/login", bytes.NewReader(body))
	if err != nil {
		return "", err
placeholder
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
placeholder
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("crs login failed: status=%d body=%s", resp.StatusCode, string(raw))
placeholder

	var parsed crsLoginResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", fmt.Errorf("crs login parse failed: %w", err)
placeholder
	if !parsed.Success || strings.TrimSpace(parsed.Token) == "" {
		msg := parsed.Message
		if msg == "" {
			msg = parsed.Error
	placeholder
		if msg == "" {
			msg = "unknown error"
	placeholder
		return "", errors.New("crs login failed: " + msg)
placeholder
	return parsed.Token, nil
placeholder

func crsExportAccounts(ctx context.Context, client *http.Client, baseURL, adminToken string) (*crsExportResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/admin/sync/export-accounts?include_secrets=true", nil)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer resp.Body.Close()

	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 5<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("crs export failed: status=%d body=%s", resp.StatusCode, string(raw))
placeholder

	var parsed crsExportResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("crs export parse failed: %w", err)
placeholder
	if !parsed.Success {
		msg := parsed.Message
		if msg == "" {
			msg = parsed.Error
	placeholder
		if msg == "" {
			msg = "unknown error"
	placeholder
		return nil, errors.New("crs export failed: " + msg)
placeholder
	return &parsed, nil
placeholder
