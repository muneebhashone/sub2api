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

// Sora2APIModel represents a model entry returned by sora2api.
type Sora2APIModel struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	OwnedBy     string `json:"owned_by,omitempty"`
	Description string `json:"description,omitempty"`
placeholder

// Sora2APIModelList represents /v1/models response.
type Sora2APIModelList struct {
	Object string          `json:"object"`
	Data   []Sora2APIModel `json:"data"`
placeholder

// Sora2APIImportTokenItem mirrors sora2api ImportTokenItem.
type Sora2APIImportTokenItem struct {
	Email            string `json:"email"`
	AccessToken      string `json:"access_token,omitempty"`
	SessionToken     string `json:"session_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	ClientID         string `json:"client_id,omitempty"`
	ProxyURL         string `json:"proxy_url,omitempty"`
	Remark           string `json:"remark,omitempty"`
	IsActive         bool   `json:"is_active"`
	ImageEnabled     bool   `json:"image_enabled"`
	VideoEnabled     bool   `json:"video_enabled"`
	ImageConcurrency int    `json:"image_concurrency"`
	VideoConcurrency int    `json:"video_concurrency"`
placeholder

// Sora2APIToken represents minimal fields for admin list.
type Sora2APIToken struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
placeholder

// Sora2APIService provides access to sora2api endpoints.
type Sora2APIService struct {
	cfg *config.Config

	baseURL         string
	apiKey          string
	adminUsername   string
	adminPassword   string
	adminTokenTTL   time.Duration
	tokenImportMode string

	client      *http.Client
	adminClient *http.Client

	adminToken   string
	adminTokenAt time.Time
	adminMu      sync.Mutex

	modelCache []Sora2APIModel
	modelMu    sync.RWMutex
placeholder

func NewSora2APIService(cfg *config.Config) *Sora2APIService {
	if cfg == nil {
		return &Sora2APIService{placeholder
placeholder
	adminTTL := time.Duration(cfg.Sora2API.AdminTokenTTLSeconds) * time.Second
	if adminTTL <= 0 {
		adminTTL = 15 * time.Minute
placeholder
	adminTimeout := time.Duration(cfg.Sora2API.AdminTimeoutSeconds) * time.Second
	if adminTimeout <= 0 {
		adminTimeout = 10 * time.Second
placeholder
	return &Sora2APIService{
		cfg:             cfg,
		baseURL:         strings.TrimRight(strings.TrimSpace(cfg.Sora2API.BaseURL), "/"),
		apiKey:          strings.TrimSpace(cfg.Sora2API.APIKey),
		adminUsername:   strings.TrimSpace(cfg.Sora2API.AdminUsername),
		adminPassword:   strings.TrimSpace(cfg.Sora2API.AdminPassword),
		adminTokenTTL:   adminTTL,
		tokenImportMode: strings.ToLower(strings.TrimSpace(cfg.Sora2API.TokenImportMode)),
		client:          &http.Client{placeholder,
		adminClient:     &http.Client{Timeout: adminTimeoutplaceholder,
placeholder
placeholder

func (s *Sora2APIService) Enabled() bool {
	return s != nil && s.baseURL != "" && s.apiKey != ""
placeholder

func (s *Sora2APIService) AdminEnabled() bool {
	return s != nil && s.baseURL != "" && s.adminUsername != "" && s.adminPassword != ""
placeholder

func (s *Sora2APIService) buildURL(path string) string {
	if s.baseURL == "" {
		return path
placeholder
	if strings.HasPrefix(path, "/") {
		return s.baseURL + path
placeholder
	return s.baseURL + "/" + path
placeholder

// BuildURL 返回完整的 sora2api URL（用于代理媒体）
func (s *Sora2APIService) BuildURL(path string) string {
	return s.buildURL(path)
placeholder

func (s *Sora2APIService) NewAPIRequest(ctx context.Context, method string, path string, body []byte) (*http.Request, error) {
	if !s.Enabled() {
		return nil, errors.New("sora2api not configured")
placeholder
	req, err := http.NewRequestWithContext(ctx, method, s.buildURL(path), bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
placeholder

func (s *Sora2APIService) ListModels(ctx context.Context) ([]Sora2APIModel, error) {
	if !s.Enabled() {
		return nil, errors.New("sora2api not configured")
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.buildURL("/v1/models"), nil)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	resp, err := s.client.Do(req)
	if err != nil {
		return s.cachedModelsOnError(err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode != http.StatusOK {
		return s.cachedModelsOnError(fmt.Errorf("sora2api models status: %d", resp.StatusCode))
placeholder

	var payload Sora2APIModelList
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return s.cachedModelsOnError(err)
placeholder
	models := payload.Data
	if s.cfg != nil && s.cfg.Gateway.SoraModelFilters.HidePromptEnhance {
		filtered := make([]Sora2APIModel, 0, len(models))
		for _, m := range models {
			if strings.HasPrefix(strings.ToLower(m.ID), "prompt-enhance") {
				continue
		placeholder
			filtered = append(filtered, m)
	placeholder
		models = filtered
placeholder

	s.modelMu.Lock()
	s.modelCache = models
	s.modelMu.Unlock()

	return models, nil
placeholder

func (s *Sora2APIService) cachedModelsOnError(err error) ([]Sora2APIModel, error) {
	s.modelMu.RLock()
	cached := append([]Sora2APIModel(nil), s.modelCache...)
	s.modelMu.RUnlock()
	if len(cached) > 0 {
		log.Printf("[Sora2API] 模型列表拉取失败，回退缓存: %v", err)
		return cached, nil
placeholder
	return nil, err
placeholder

func (s *Sora2APIService) ImportTokens(ctx context.Context, items []Sora2APIImportTokenItem) error {
	if !s.AdminEnabled() {
		return errors.New("sora2api admin not configured")
placeholder
	mode := s.tokenImportMode
	if mode == "" {
		mode = "at"
placeholder
	payload := map[string]any{
		"tokens": items,
		"mode":   mode,
placeholder
	_, err := s.doAdminRequest(ctx, http.MethodPost, "/api/tokens/import", payload, nil)
	return err
placeholder

func (s *Sora2APIService) ListTokens(ctx context.Context) ([]Sora2APIToken, error) {
	if !s.AdminEnabled() {
		return nil, errors.New("sora2api admin not configured")
placeholder
	var tokens []Sora2APIToken
	_, err := s.doAdminRequest(ctx, http.MethodGet, "/api/tokens", nil, &tokens)
	return tokens, err
placeholder

func (s *Sora2APIService) DisableToken(ctx context.Context, tokenID int64) error {
	if !s.AdminEnabled() {
		return errors.New("sora2api admin not configured")
placeholder
	path := fmt.Sprintf("/api/tokens/%d/disable", tokenID)
	_, err := s.doAdminRequest(ctx, http.MethodPost, path, nil, nil)
	return err
placeholder

func (s *Sora2APIService) DeleteToken(ctx context.Context, tokenID int64) error {
	if !s.AdminEnabled() {
		return errors.New("sora2api admin not configured")
placeholder
	path := fmt.Sprintf("/api/tokens/%d", tokenID)
	_, err := s.doAdminRequest(ctx, http.MethodDelete, path, nil, nil)
	return err
placeholder

func (s *Sora2APIService) doAdminRequest(ctx context.Context, method string, path string, body any, out any) (*http.Response, error) {
	if !s.AdminEnabled() {
		return nil, errors.New("sora2api admin not configured")
placeholder
	token, err := s.getAdminToken(ctx)
	if err != nil {
		return nil, err
placeholder
	resp, err := s.doAdminRequestWithToken(ctx, method, path, token, body, out)
	if err == nil && resp != nil && resp.StatusCode != http.StatusUnauthorized {
		return resp, nil
placeholder
	if resp != nil && resp.StatusCode == http.StatusUnauthorized {
		s.invalidateAdminToken()
		token, err = s.getAdminToken(ctx)
		if err != nil {
			return resp, err
	placeholder
		return s.doAdminRequestWithToken(ctx, method, path, token, body, out)
placeholder
	return resp, err
placeholder

func (s *Sora2APIService) doAdminRequestWithToken(ctx context.Context, method string, path string, token string, body any, out any) (*http.Response, error) {
	var reader *bytes.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, err
	placeholder
		reader = bytes.NewReader(buf)
placeholder else {
		reader = bytes.NewReader(nil)
placeholder
	req, err := http.NewRequestWithContext(ctx, method, s.buildURL(path), reader)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
placeholder
	resp, err := s.adminClient.Do(req)
	if err != nil {
		return resp, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, fmt.Errorf("sora2api admin status: %d", resp.StatusCode)
placeholder
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return resp, err
	placeholder
placeholder
	return resp, nil
placeholder

func (s *Sora2APIService) getAdminToken(ctx context.Context) (string, error) {
	s.adminMu.Lock()
	defer s.adminMu.Unlock()

	if s.adminToken != "" && time.Since(s.adminTokenAt) < s.adminTokenTTL {
		return s.adminToken, nil
placeholder

	if !s.AdminEnabled() {
		return "", errors.New("sora2api admin not configured")
placeholder

	payload := map[string]string{
		"username": s.adminUsername,
		"password": s.adminPassword,
placeholder
	buf, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.buildURL("/api/login"), bytes.NewReader(buf))
	if err != nil {
		return "", err
placeholder
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.adminClient.Do(req)
	if err != nil {
		return "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("sora2api login failed: %d", resp.StatusCode)
placeholder
	var result struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
		Message string `json:"message"`
placeholder
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
placeholder
	if !result.Success || result.Token == "" {
		if result.Message == "" {
			result.Message = "sora2api login failed"
	placeholder
		return "", errors.New(result.Message)
placeholder
	s.adminToken = result.Token
	s.adminTokenAt = time.Now()
	return result.Token, nil
placeholder

func (s *Sora2APIService) invalidateAdminToken() {
	s.adminMu.Lock()
	defer s.adminMu.Unlock()
	s.adminToken = ""
	s.adminTokenAt = time.Time{placeholder
placeholder
