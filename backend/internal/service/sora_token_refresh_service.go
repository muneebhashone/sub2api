package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

const defaultSoraClientID = "app_LlGpXReQgckcGGUo2JrYvtJK"

// SoraTokenRefreshService handles Sora access token refresh.
type SoraTokenRefreshService struct {
	accountRepo     AccountRepository
	soraAccountRepo SoraAccountRepository
	settingService  *SettingService
	httpUpstream    HTTPUpstream
	cfg             *config.Config
	stopCh          chan struct{placeholder
	stopOnce        sync.Once
placeholder

func NewSoraTokenRefreshService(
	accountRepo AccountRepository,
	soraAccountRepo SoraAccountRepository,
	settingService *SettingService,
	httpUpstream HTTPUpstream,
	cfg *config.Config,
) *SoraTokenRefreshService {
	return &SoraTokenRefreshService{
		accountRepo:     accountRepo,
		soraAccountRepo: soraAccountRepo,
		settingService:  settingService,
		httpUpstream:    httpUpstream,
		cfg:             cfg,
		stopCh:          make(chan struct{placeholder),
placeholder
placeholder

func (s *SoraTokenRefreshService) Start() {
	if s == nil {
		return
placeholder
	go s.refreshLoop()
placeholder

func (s *SoraTokenRefreshService) Stop() {
	if s == nil {
		return
placeholder
	s.stopOnce.Do(func() {
		close(s.stopCh)
placeholder)
placeholder

func (s *SoraTokenRefreshService) refreshLoop() {
	for {
		wait := s.nextRunDelay()
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
			s.refreshOnce()
		case <-s.stopCh:
			timer.Stop()
			return
	placeholder
placeholder
placeholder

func (s *SoraTokenRefreshService) refreshOnce() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	if !s.isEnabled(ctx) {
		log.Println("[SoraTokenRefresh] disabled by settings")
		return
placeholder
	if s.accountRepo == nil || s.soraAccountRepo == nil {
		log.Println("[SoraTokenRefresh] repository not configured")
		return
placeholder

	accounts, err := s.accountRepo.ListByPlatform(ctx, PlatformSora)
	if err != nil {
		log.Printf("[SoraTokenRefresh] list accounts failed: %v", err)
		return
placeholder
	if len(accounts) == 0 {
		log.Println("[SoraTokenRefresh] no sora accounts")
		return
placeholder
	ids := make([]int64, 0, len(accounts))
	accountMap := make(map[int64]*Account, len(accounts))
	for i := range accounts {
		acc := accounts[i]
		ids = append(ids, acc.ID)
		accountMap[acc.ID] = &acc
placeholder
	accountExtras, err := s.soraAccountRepo.GetByAccountIDs(ctx, ids)
	if err != nil {
		log.Printf("[SoraTokenRefresh] load sora accounts failed: %v", err)
		return
placeholder

	success := 0
	failed := 0
	skipped := 0
	for accountID, account := range accountMap {
		extra := accountExtras[accountID]
		if extra == nil {
			skipped++
			continue
	placeholder
		result, err := s.refreshForAccount(ctx, account, extra)
		if err != nil {
			failed++
			log.Printf("[SoraTokenRefresh] account %d refresh failed: %v", accountID, err)
			continue
	placeholder
		if result == nil {
			skipped++
			continue
	placeholder

		updates := map[string]any{
			"access_token": result.AccessToken,
	placeholder
		if result.RefreshToken != "" {
			updates["refresh_token"] = result.RefreshToken
	placeholder
		if result.Email != "" {
			updates["email"] = result.Email
	placeholder
		if err := s.soraAccountRepo.Upsert(ctx, accountID, updates); err != nil {
			failed++
			log.Printf("[SoraTokenRefresh] account %d update failed: %v", accountID, err)
			continue
	placeholder
		success++
placeholder
	log.Printf("[SoraTokenRefresh] done: success=%d failed=%d skipped=%d", success, failed, skipped)
placeholder

func (s *SoraTokenRefreshService) refreshForAccount(ctx context.Context, account *Account, extra *SoraAccount) (*soraRefreshResult, error) {
	if extra == nil {
		return nil, nil
placeholder
	if strings.TrimSpace(extra.SessionToken) == "" && strings.TrimSpace(extra.RefreshToken) == "" {
		return nil, nil
placeholder

	if extra.SessionToken != "" {
		result, err := s.refreshWithSessionToken(ctx, account, extra.SessionToken)
		if err == nil && result != nil && result.AccessToken != "" {
			return result, nil
	placeholder
		if strings.TrimSpace(extra.RefreshToken) == "" {
			return nil, err
	placeholder
placeholder

	clientID := strings.TrimSpace(extra.ClientID)
	if clientID == "" {
		clientID = defaultSoraClientID
placeholder
	return s.refreshWithRefreshToken(ctx, account, extra.RefreshToken, clientID)
placeholder

type soraRefreshResult struct {
	AccessToken  string
	RefreshToken string
	Email        string
placeholder

type soraSessionResponse struct {
	AccessToken string `json:"accessToken"`
	User        struct {
		Email string `json:"email"`
placeholder `json:"user"`
placeholder

type soraRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
placeholder

func (s *SoraTokenRefreshService) refreshWithSessionToken(ctx context.Context, account *Account, sessionToken string) (*soraRefreshResult, error) {
	if s.httpUpstream == nil {
		return nil, fmt.Errorf("upstream not configured")
placeholder
	req, err := http.NewRequestWithContext(ctx, "GET", "https://sora.chatgpt.com/api/auth/session", nil)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Cookie", "__Secure-next-auth.session-token="+sessionToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://sora.chatgpt.com")
	req.Header.Set("Referer", "https://sora.chatgpt.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	enableTLS := false
	if s.cfg != nil {
		enableTLS = s.cfg.Gateway.TLSFingerprint.Enabled
placeholder
	proxyURL := ""
	accountConcurrency := 0
	accountID := int64(0)
	if account != nil {
		accountID = account.ID
		accountConcurrency = account.Concurrency
		if account.Proxy != nil {
			proxyURL = account.Proxy.URL()
	placeholder
placeholder
	resp, err := s.httpUpstream.DoWithTLS(req, proxyURL, accountID, accountConcurrency, enableTLS)
	if err != nil {
		return nil, err
placeholder
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("session refresh failed: %d", resp.StatusCode)
placeholder
	var payload soraSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
placeholder
	if payload.AccessToken == "" {
		return nil, errors.New("session refresh missing access token")
placeholder
	return &soraRefreshResult{AccessToken: payload.AccessToken, Email: payload.User.Emailplaceholder, nil
placeholder

func (s *SoraTokenRefreshService) refreshWithRefreshToken(ctx context.Context, account *Account, refreshToken, clientID string) (*soraRefreshResult, error) {
	if s.httpUpstream == nil {
		return nil, fmt.Errorf("upstream not configured")
placeholder
	payload := map[string]any{
		"client_id":     clientID,
		"grant_type":    "refresh_token",
		"redirect_uri":  "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
		"refresh_token": refreshToken,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
placeholder
	req, err := http.NewRequestWithContext(ctx, "POST", "https://auth.openai.com/oauth/token", bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	enableTLS := false
	if s.cfg != nil {
		enableTLS = s.cfg.Gateway.TLSFingerprint.Enabled
placeholder
	proxyURL := ""
	accountConcurrency := 0
	accountID := int64(0)
	if account != nil {
		accountID = account.ID
		accountConcurrency = account.Concurrency
		if account.Proxy != nil {
			proxyURL = account.Proxy.URL()
	placeholder
placeholder
	resp, err := s.httpUpstream.DoWithTLS(req, proxyURL, accountID, accountConcurrency, enableTLS)
	if err != nil {
		return nil, err
placeholder
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("refresh token failed: %d", resp.StatusCode)
placeholder
	var payloadResp soraRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&payloadResp); err != nil {
		return nil, err
placeholder
	if payloadResp.AccessToken == "" {
		return nil, errors.New("refresh token missing access token")
placeholder
	return &soraRefreshResult{AccessToken: payloadResp.AccessToken, RefreshToken: payloadResp.RefreshTokenplaceholder, nil
placeholder

func (s *SoraTokenRefreshService) nextRunDelay() time.Duration {
	location := time.Local
	if s.cfg != nil && strings.TrimSpace(s.cfg.Timezone) != "" {
		if tz, err := time.LoadLocation(strings.TrimSpace(s.cfg.Timezone)); err == nil {
			location = tz
	placeholder
placeholder
	now := time.Now().In(location)
	next := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).Add(24 * time.Hour)
	return time.Until(next)
placeholder

func (s *SoraTokenRefreshService) isEnabled(ctx context.Context) bool {
	if s.settingService == nil {
		return s.cfg != nil && s.cfg.Sora.TokenRefresh.Enabled
placeholder
	cfg := s.settingService.GetSoraConfig(ctx)
	return cfg.TokenRefresh.Enabled
placeholder
