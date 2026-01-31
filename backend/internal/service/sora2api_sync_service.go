package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Sora2APISyncService 用于同步 Sora 账号到 sora2api token 池
type Sora2APISyncService struct {
	sora2api    *Sora2APIService
	accountRepo AccountRepository
	httpClient  *http.Client
placeholder

func NewSora2APISyncService(sora2api *Sora2APIService, accountRepo AccountRepository) *Sora2APISyncService {
	return &Sora2APISyncService{
		sora2api:    sora2api,
		accountRepo: accountRepo,
		httpClient:  &http.Client{Timeout: 10 * time.Secondplaceholder,
placeholder
placeholder

func (s *Sora2APISyncService) Enabled() bool {
	return s != nil && s.sora2api != nil && s.sora2api.AdminEnabled()
placeholder

// SyncAccount 将 Sora 账号同步到 sora2api（导入或更新）
func (s *Sora2APISyncService) SyncAccount(ctx context.Context, account *Account) error {
	if !s.Enabled() {
		return nil
placeholder
	if account == nil || account.Platform != PlatformSora {
		return nil
placeholder

	accessToken := strings.TrimSpace(account.GetCredential("access_token"))
	if accessToken == "" {
		return errors.New("sora 账号缺少 access_token")
placeholder

	email, updated := s.resolveAccountEmail(ctx, account)
	if email == "" {
		return errors.New("无法解析 Sora 账号邮箱")
placeholder
	if updated && s.accountRepo != nil {
		if err := s.accountRepo.Update(ctx, account); err != nil {
			log.Printf("[SoraSync] 更新账号邮箱失败: account_id=%d err=%v", account.ID, err)
	placeholder
placeholder

	item := Sora2APIImportTokenItem{
		Email:            email,
		AccessToken:      accessToken,
		SessionToken:     strings.TrimSpace(account.GetCredential("session_token")),
		RefreshToken:     strings.TrimSpace(account.GetCredential("refresh_token")),
		ClientID:         strings.TrimSpace(account.GetCredential("client_id")),
		Remark:           account.Name,
		IsActive:         account.IsActive() && account.Schedulable,
		ImageEnabled:     true,
		VideoEnabled:     true,
		ImageConcurrency: normalizeSoraConcurrency(account.Concurrency),
		VideoConcurrency: normalizeSoraConcurrency(account.Concurrency),
placeholder

	if err := s.sora2api.ImportTokens(ctx, []Sora2APIImportTokenItem{itemplaceholder); err != nil {
		return err
placeholder
	return nil
placeholder

// DisableAccount 禁用 sora2api 中的 token
func (s *Sora2APISyncService) DisableAccount(ctx context.Context, account *Account) error {
	if !s.Enabled() {
		return nil
placeholder
	if account == nil || account.Platform != PlatformSora {
		return nil
placeholder
	tokenID, err := s.resolveTokenID(ctx, account)
	if err != nil {
		return err
placeholder
	return s.sora2api.DisableToken(ctx, tokenID)
placeholder

// DeleteAccount 删除 sora2api 中的 token
func (s *Sora2APISyncService) DeleteAccount(ctx context.Context, account *Account) error {
	if !s.Enabled() {
		return nil
placeholder
	if account == nil || account.Platform != PlatformSora {
		return nil
placeholder
	tokenID, err := s.resolveTokenID(ctx, account)
	if err != nil {
		return err
placeholder
	return s.sora2api.DeleteToken(ctx, tokenID)
placeholder

func normalizeSoraConcurrency(value int) int {
	if value <= 0 {
		return -1
placeholder
	return value
placeholder

func (s *Sora2APISyncService) resolveAccountEmail(ctx context.Context, account *Account) (string, bool) {
	if account == nil {
		return "", false
placeholder
	if email := strings.TrimSpace(account.GetCredential("email")); email != "" {
		return email, false
placeholder
	if email := strings.TrimSpace(account.GetExtraString("email")); email != "" {
		if account.Credentials == nil {
			account.Credentials = map[string]any{placeholder
	placeholder
		account.Credentials["email"] = email
		return email, true
placeholder
	if email := strings.TrimSpace(account.GetExtraString("sora_email")); email != "" {
		if account.Credentials == nil {
			account.Credentials = map[string]any{placeholder
	placeholder
		account.Credentials["email"] = email
		return email, true
placeholder

	accessToken := strings.TrimSpace(account.GetCredential("access_token"))
	if accessToken != "" {
		if email := extractEmailFromAccessToken(accessToken); email != "" {
			if account.Credentials == nil {
				account.Credentials = map[string]any{placeholder
		placeholder
			account.Credentials["email"] = email
			return email, true
	placeholder
		if email := s.fetchEmailFromSora(ctx, accessToken); email != "" {
			if account.Credentials == nil {
				account.Credentials = map[string]any{placeholder
		placeholder
			account.Credentials["email"] = email
			return email, true
	placeholder
placeholder

	return "", false
placeholder

func (s *Sora2APISyncService) resolveTokenID(ctx context.Context, account *Account) (int64, error) {
	if account == nil {
		return 0, errors.New("account is nil")
placeholder

	if account.Extra != nil {
		if v, ok := account.Extra["sora2api_token_id"]; ok {
			if id, ok := v.(float64); ok && id > 0 {
				return int64(id), nil
		placeholder
			if id, ok := v.(int64); ok && id > 0 {
				return id, nil
		placeholder
			if id, ok := v.(int); ok && id > 0 {
				return int64(id), nil
		placeholder
	placeholder
placeholder

	email := strings.TrimSpace(account.GetCredential("email"))
	if email == "" {
		email, _ = s.resolveAccountEmail(ctx, account)
placeholder
	if email == "" {
		return 0, errors.New("sora2api token email missing")
placeholder

	tokenID, err := s.findTokenIDByEmail(ctx, email)
	if err != nil {
		return 0, err
placeholder
	return tokenID, nil
placeholder

func (s *Sora2APISyncService) findTokenIDByEmail(ctx context.Context, email string) (int64, error) {
	if !s.Enabled() {
		return 0, errors.New("sora2api admin not configured")
placeholder
	tokens, err := s.sora2api.ListTokens(ctx)
	if err != nil {
		return 0, err
placeholder
	for _, token := range tokens {
		if strings.EqualFold(strings.TrimSpace(token.Email), strings.TrimSpace(email)) {
			return token.ID, nil
	placeholder
placeholder
	return 0, fmt.Errorf("sora2api token not found for email: %s", email)
placeholder

func extractEmailFromAccessToken(accessToken string) string {
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	claims := jwt.MapClaims{placeholder
	_, _, err := parser.ParseUnverified(accessToken, claims)
	if err != nil {
		return ""
placeholder
	if email, ok := claims["email"].(string); ok && strings.TrimSpace(email) != "" {
		return email
placeholder
	if profile, ok := claims["https://api.openai.com/profile"].(map[string]any); ok {
		if email, ok := profile["email"].(string); ok && strings.TrimSpace(email) != "" {
			return email
	placeholder
placeholder
	return ""
placeholder

func (s *Sora2APISyncService) fetchEmailFromSora(ctx context.Context, accessToken string) string {
	if s.httpClient == nil {
		return ""
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, soraMeAPIURL, nil)
	if err != nil {
		return ""
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", "Sora/1.2026.007 (Android 15; 24122RKC7C; build 2600700)")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return ""
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		return ""
placeholder
	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return ""
placeholder
	if email, ok := payload["email"].(string); ok && strings.TrimSpace(email) != "" {
		return email
placeholder
	return ""
placeholder
