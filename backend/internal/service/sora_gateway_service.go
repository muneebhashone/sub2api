package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/sora"
	"github.com/Wei-Shaw/sub2api/internal/util/urlvalidator"
)

const (
	soraErrorDisableThreshold = 5
	maxImageDownloadSize      = 20 * 1024 * 1024  // 20MB
	maxVideoDownloadSize      = 200 * 1024 * 1024 // 200MB
)

var (
	ErrSoraAccountMissingToken = errors.New("sora account missing access token")
	ErrSoraAccountNotEligible  = errors.New("sora account not eligible")
)

// SoraGenerationRequest 表示 Sora 生成请求。
type SoraGenerationRequest struct {
	Model         string
	Prompt        string
	Image         string
	Video         string
	RemixTargetID string
	Stream        bool
	UserID        int64
placeholder

// SoraGenerationResult 表示 Sora 生成结果。
type SoraGenerationResult struct {
	Content    string
	MediaType  string
	ResultURLs []string
	TaskID     string
placeholder

// SoraGatewayService 处理 Sora 生成流程。
type SoraGatewayService struct {
	accountRepo     AccountRepository
	soraAccountRepo SoraAccountRepository
	usageRepo       SoraUsageStatRepository
	taskRepo        SoraTaskRepository
	cacheService    *SoraCacheService
	settingService  *SettingService
	concurrency     *ConcurrencyService
	cfg             *config.Config
	httpUpstream    HTTPUpstream
placeholder

// NewSoraGatewayService 创建 SoraGatewayService。
func NewSoraGatewayService(
	accountRepo AccountRepository,
	soraAccountRepo SoraAccountRepository,
	usageRepo SoraUsageStatRepository,
	taskRepo SoraTaskRepository,
	cacheService *SoraCacheService,
	settingService *SettingService,
	concurrencyService *ConcurrencyService,
	cfg *config.Config,
	httpUpstream HTTPUpstream,
) *SoraGatewayService {
	return &SoraGatewayService{
		accountRepo:     accountRepo,
		soraAccountRepo: soraAccountRepo,
		usageRepo:       usageRepo,
		taskRepo:        taskRepo,
		cacheService:    cacheService,
		settingService:  settingService,
		concurrency:     concurrencyService,
		cfg:             cfg,
		httpUpstream:    httpUpstream,
placeholder
placeholder

// ListModels 返回 Sora 模型列表。
func (s *SoraGatewayService) ListModels() []sora.ModelListItem {
	return sora.ListModels()
placeholder

// Generate 执行 Sora 生成流程。
func (s *SoraGatewayService) Generate(ctx context.Context, account *Account, req SoraGenerationRequest) (*SoraGenerationResult, error) {
	client, cfg := s.getClient(ctx)
	if client == nil {
		return nil, errors.New("sora client is not configured")
placeholder
	modelCfg, ok := sora.ModelConfigs[req.Model]
	if !ok {
		return nil, fmt.Errorf("unsupported model: %s", req.Model)
placeholder
	accessToken, soraAcc, err := s.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	if soraAcc != nil && soraAcc.SoraCooldownUntil != nil && time.Now().Before(*soraAcc.SoraCooldownUntil) {
		return nil, ErrSoraAccountNotEligible
placeholder
	if modelCfg.RequirePro && !isSoraProAccount(soraAcc) {
		return nil, ErrSoraAccountNotEligible
placeholder
	if modelCfg.Type == "video" && soraAcc != nil {
		if !soraAcc.VideoEnabled || !soraAcc.SoraSupported || soraAcc.IsExpired {
			return nil, ErrSoraAccountNotEligible
	placeholder
placeholder
	if modelCfg.Type == "image" && soraAcc != nil {
		if !soraAcc.ImageEnabled || soraAcc.IsExpired {
			return nil, ErrSoraAccountNotEligible
	placeholder
placeholder

	opts := sora.RequestOptions{
		AccountID:          account.ID,
		AccountConcurrency: account.Concurrency,
		AccessToken:        accessToken,
placeholder
	if account.Proxy != nil {
		opts.ProxyURL = account.Proxy.URL()
placeholder

	releaseFunc, err := s.acquireSoraSlots(ctx, account, soraAcc, modelCfg.Type == "video")
	if err != nil {
		return nil, err
placeholder
	if releaseFunc != nil {
		defer releaseFunc()
placeholder

	if modelCfg.Type == "prompt_enhance" {
		content, err := client.EnhancePrompt(ctx, opts, req.Prompt, modelCfg.ExpansionLevel, modelCfg.DurationS)
		if err != nil {
			return nil, err
	placeholder
		return &SoraGenerationResult{Content: content, MediaType: "text"placeholder, nil
placeholder

	var mediaID string
	if req.Image != "" {
		data, err := s.loadImageBytes(ctx, opts, req.Image)
		if err != nil {
			return nil, err
	placeholder
		mediaID, err = client.UploadImage(ctx, opts, data, "image.png")
		if err != nil {
			return nil, err
	placeholder
placeholder
	if req.Video != "" && modelCfg.Type != "video" {
		return nil, errors.New("视频输入仅支持视频模型")
placeholder
	if req.Video != "" && req.Image != "" {
		return nil, errors.New("不能同时传入 image 与 video")
placeholder

	var cleanupCharacter func()
	if req.Video != "" && req.RemixTargetID == "" {
		username, characterID, err := s.createCharacter(ctx, client, opts, req.Video)
		if err != nil {
			return nil, err
	placeholder
		if strings.TrimSpace(req.Prompt) == "" {
			return &SoraGenerationResult{
				Content:   fmt.Sprintf("角色创建成功，角色名@%s", username),
				MediaType: "text",
		placeholder, nil
	placeholder
		if username != "" {
			req.Prompt = fmt.Sprintf("@%s %s", username, strings.TrimSpace(req.Prompt))
	placeholder
		if characterID != "" {
			cleanupCharacter = func() {
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				_ = client.DeleteCharacter(ctx, opts, characterID)
		placeholder
	placeholder
placeholder
	if cleanupCharacter != nil {
		defer cleanupCharacter()
placeholder

	var taskID string
	if modelCfg.Type == "image" {
		taskID, err = client.GenerateImage(ctx, opts, req.Prompt, modelCfg.Width, modelCfg.Height, mediaID)
placeholder else {
		orientation := modelCfg.Orientation
		if orientation == "" {
			orientation = "landscape"
	placeholder
		modelName := modelCfg.Model
		if modelName == "" {
			modelName = "sy_8"
	placeholder
		size := modelCfg.Size
		if size == "" {
			size = "small"
	placeholder
		if req.RemixTargetID != "" {
			taskID, err = client.RemixVideo(ctx, opts, req.RemixTargetID, req.Prompt, orientation, modelCfg.NFrames, "")
	placeholder else if sora.IsStoryboardPrompt(req.Prompt) {
			formatted := sora.FormatStoryboardPrompt(req.Prompt)
			taskID, err = client.GenerateStoryboard(ctx, opts, formatted, orientation, modelCfg.NFrames, mediaID, "")
	placeholder else {
			taskID, err = client.GenerateVideo(ctx, opts, req.Prompt, orientation, modelCfg.NFrames, mediaID, "", modelName, size)
	placeholder
placeholder
	if err != nil {
		return nil, err
placeholder

	if s.taskRepo != nil {
		_ = s.taskRepo.Create(ctx, &SoraTask{
			TaskID:    taskID,
			AccountID: account.ID,
			Model:     req.Model,
			Prompt:    req.Prompt,
			Status:    "processing",
			Progress:  0,
			CreatedAt: time.Now(),
	placeholder)
placeholder

	result, err := s.pollResult(ctx, client, cfg, opts, taskID, modelCfg.Type == "video", req)
	if err != nil {
		if s.taskRepo != nil {
			_ = s.taskRepo.UpdateStatus(ctx, taskID, "failed", 0, "", err.Error(), timePtr(time.Now()))
	placeholder
		consecutive := 0
		if s.usageRepo != nil {
			consecutive, _ = s.usageRepo.RecordError(ctx, account.ID)
	placeholder
		if consecutive >= soraErrorDisableThreshold {
			_ = s.accountRepo.SetError(ctx, account.ID, "Sora 连续错误次数过多，已自动禁用")
	placeholder
		return nil, err
placeholder

	if s.taskRepo != nil {
		payload, _ := json.Marshal(result.ResultURLs)
		_ = s.taskRepo.UpdateStatus(ctx, taskID, "completed", 100, string(payload), "", timePtr(time.Now()))
placeholder
	if s.usageRepo != nil {
		_ = s.usageRepo.RecordSuccess(ctx, account.ID, modelCfg.Type == "video")
placeholder
	return result, nil
placeholder

func (s *SoraGatewayService) pollResult(ctx context.Context, client *sora.Client, cfg config.SoraConfig, opts sora.RequestOptions, taskID string, isVideo bool, req SoraGenerationRequest) (*SoraGenerationResult, error) {
	if taskID == "" {
		return nil, errors.New("missing task id")
placeholder
	pollInterval := 2 * time.Second
	if cfg.PollInterval > 0 {
		pollInterval = time.Duration(cfg.PollInterval*1000) * time.Millisecond
placeholder
	timeout := 300 * time.Second
	if cfg.Timeout > 0 {
		timeout = time.Duration(cfg.Timeout) * time.Second
placeholder
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
	placeholder
		if isVideo {
			pending, err := client.GetPendingTasks(ctx, opts)
			if err == nil {
				for _, task := range pending {
					if stringFromMap(task, "id") == taskID {
						continue
				placeholder
			placeholder
		placeholder
			drafts, err := client.GetVideoDrafts(ctx, opts)
			if err != nil {
				return nil, err
		placeholder
			items, _ := drafts["items"].([]any)
			for _, item := range items {
				entry, ok := item.(map[string]any)
				if !ok {
					continue
			placeholder
				if stringFromMap(entry, "task_id") != taskID {
					continue
			placeholder
				url := firstNonEmpty(stringFromMap(entry, "downloadable_url"), stringFromMap(entry, "url"))
				reason := stringFromMap(entry, "reason_str")
				if url == "" {
					if reason == "" {
						reason = "视频生成失败"
				placeholder
					return nil, errors.New(reason)
			placeholder
				finalURL, err := s.handleWatermark(ctx, client, cfg, opts, url, entry, req, opts.AccountID, taskID)
				if err != nil {
					return nil, err
			placeholder
				return &SoraGenerationResult{
					Content:    buildVideoMarkdown(finalURL),
					MediaType:  "video",
					ResultURLs: []string{finalURLplaceholder,
					TaskID:     taskID,
			placeholder, nil
		placeholder
	placeholder else {
			resp, err := client.GetImageTasks(ctx, opts)
			if err != nil {
				return nil, err
		placeholder
			tasks, _ := resp["task_responses"].([]any)
			for _, item := range tasks {
				entry, ok := item.(map[string]any)
				if !ok {
					continue
			placeholder
				if stringFromMap(entry, "id") != taskID {
					continue
			placeholder
				status := stringFromMap(entry, "status")
				switch status {
				case "succeeded":
					urls := extractImageURLs(entry)
					if len(urls) == 0 {
						return nil, errors.New("image urls empty")
				placeholder
					content := buildImageMarkdown(urls)
					return &SoraGenerationResult{
						Content:    content,
						MediaType:  "image",
						ResultURLs: urls,
						TaskID:     taskID,
				placeholder, nil
				case "failed":
					message := stringFromMap(entry, "error_message")
					if message == "" {
						message = "image generation failed"
				placeholder
					return nil, errors.New(message)
			placeholder
		placeholder
	placeholder

		time.Sleep(pollInterval)
placeholder
	return nil, errors.New("generation timeout")
placeholder

func (s *SoraGatewayService) handleWatermark(ctx context.Context, client *sora.Client, cfg config.SoraConfig, opts sora.RequestOptions, url string, entry map[string]any, req SoraGenerationRequest, accountID int64, taskID string) (string, error) {
	if !cfg.WatermarkFree.Enabled {
		return s.cacheVideo(ctx, url, req, accountID, taskID), nil
placeholder
	generationID := stringFromMap(entry, "id")
	if generationID == "" {
		return s.cacheVideo(ctx, url, req, accountID, taskID), nil
placeholder
	postID, err := client.PostVideoForWatermarkFree(ctx, opts, generationID)
	if err != nil {
		if cfg.WatermarkFree.FallbackOnFailure {
			return s.cacheVideo(ctx, url, req, accountID, taskID), nil
	placeholder
		return "", err
placeholder
	if postID == "" {
		if cfg.WatermarkFree.FallbackOnFailure {
			return s.cacheVideo(ctx, url, req, accountID, taskID), nil
	placeholder
		return "", errors.New("watermark-free post id empty")
placeholder
	var parsedURL string
	if cfg.WatermarkFree.ParseMethod == "custom" {
		if cfg.WatermarkFree.CustomParseURL == "" || cfg.WatermarkFree.CustomParseToken == "" {
			return "", errors.New("custom parse 未配置")
	placeholder
		parsedURL, err = s.fetchCustomWatermarkURL(ctx, cfg.WatermarkFree.CustomParseURL, cfg.WatermarkFree.CustomParseToken, postID)
		if err != nil {
			if cfg.WatermarkFree.FallbackOnFailure {
				return s.cacheVideo(ctx, url, req, accountID, taskID), nil
		placeholder
			return "", err
	placeholder
placeholder else {
		parsedURL = fmt.Sprintf("https://oscdn2.dyysy.com/MP4/%s.mp4", postID)
placeholder
	cached := s.cacheVideo(ctx, parsedURL, req, accountID, taskID)
	_ = client.DeletePost(ctx, opts, postID)
	return cached, nil
placeholder

func (s *SoraGatewayService) cacheVideo(ctx context.Context, url string, req SoraGenerationRequest, accountID int64, taskID string) string {
	if s.cacheService == nil {
		return url
placeholder
	file, err := s.cacheService.CacheVideo(ctx, accountID, req.UserID, taskID, url)
	if err != nil || file == nil {
		return url
placeholder
	return file.CacheURL
placeholder

func (s *SoraGatewayService) getAccessToken(ctx context.Context, account *Account) (string, *SoraAccount, error) {
	if account == nil {
		return "", nil, errors.New("account is nil")
placeholder
	var soraAcc *SoraAccount
	if s.soraAccountRepo != nil {
		soraAcc, _ = s.soraAccountRepo.GetByAccountID(ctx, account.ID)
placeholder
	if soraAcc != nil && soraAcc.AccessToken != "" {
		return soraAcc.AccessToken, soraAcc, nil
placeholder
	if account.Credentials != nil {
		if v, ok := account.Credentials["access_token"].(string); ok && v != "" {
			return v, soraAcc, nil
	placeholder
		if v, ok := account.Credentials["token"].(string); ok && v != "" {
			return v, soraAcc, nil
	placeholder
placeholder
	return "", soraAcc, ErrSoraAccountMissingToken
placeholder

func (s *SoraGatewayService) getClient(ctx context.Context) (*sora.Client, config.SoraConfig) {
	cfg := s.getSoraConfig(ctx)
	if s.httpUpstream == nil {
		return nil, cfg
placeholder
	baseURL := strings.TrimSpace(cfg.BaseURL)
	if baseURL == "" {
		return nil, cfg
placeholder
	timeout := time.Duration(cfg.Timeout) * time.Second
	if cfg.Timeout <= 0 {
		timeout = 120 * time.Second
placeholder
	enableTLS := false
	if s.cfg != nil {
		enableTLS = s.cfg.Gateway.TLSFingerprint.Enabled
placeholder
	return sora.NewClient(baseURL, timeout, s.httpUpstream, enableTLS), cfg
placeholder

func decodeBase64(raw string) ([]byte, error) {
	data := raw
	if idx := strings.Index(raw, "base64,"); idx != -1 {
		data = raw[idx+7:]
placeholder
	return base64.StdEncoding.DecodeString(data)
placeholder

func extractImageURLs(entry map[string]any) []string {
	generations, _ := entry["generations"].([]any)
	urls := make([]string, 0, len(generations))
	for _, gen := range generations {
		m, ok := gen.(map[string]any)
		if !ok {
			continue
	placeholder
		if url, ok := m["url"].(string); ok && url != "" {
			urls = append(urls, url)
	placeholder
placeholder
	return urls
placeholder

func buildImageMarkdown(urls []string) string {
	parts := make([]string, 0, len(urls))
	for _, u := range urls {
		parts = append(parts, fmt.Sprintf("![Generated Image](%s)", u))
placeholder
	return strings.Join(parts, "\n")
placeholder

func buildVideoMarkdown(url string) string {
	return fmt.Sprintf("```html\n<video src='%s' controls></video>\n```", url)
placeholder

func stringFromMap(m map[string]any, key string) string {
	if m == nil {
		return ""
placeholder
	if v, ok := m[key].(string); ok {
		return v
placeholder
	return ""
placeholder

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
	placeholder
placeholder
	return ""
placeholder

func isSoraProAccount(acc *SoraAccount) bool {
	if acc == nil {
		return false
placeholder
	return strings.EqualFold(acc.PlanType, "chatgpt_pro")
placeholder

func timePtr(t time.Time) *time.Time {
	return &t
placeholder

// fetchCustomWatermarkURL 使用自定义解析服务获取无水印视频 URL
func (s *SoraGatewayService) fetchCustomWatermarkURL(ctx context.Context, parseURL, parseToken, postID string) (string, error) {
	// 使用项目的 URL 校验器验证 parseURL 格式，防止 SSRF 攻击
	if _, err := urlvalidator.ValidateHTTPSURL(parseURL, urlvalidator.ValidationOptions{placeholder); err != nil {
		return "", fmt.Errorf("无效的解析服务地址: %w", err)
placeholder

	payload := map[string]any{
		"url":   fmt.Sprintf("https://sora.chatgpt.com/p/%s", postID),
		"token": parseToken,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	req, err := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(parseURL, "/")+"/get-sora-link", strings.NewReader(string(body)))
	if err != nil {
		return "", err
placeholder
	req.Header.Set("Content-Type", "application/json")

	// 复用 httpUpstream，遵守代理和 TLS 配置
	enableTLS := false
	if s.cfg != nil {
		enableTLS = s.cfg.Gateway.TLSFingerprint.Enabled
placeholder
	resp, err := s.httpUpstream.DoWithTLS(req, "", 0, 1, enableTLS)
	if err != nil {
		return "", err
placeholder
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("custom parse failed: %d", resp.StatusCode)
placeholder
	var parsed map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
placeholder
	if errMsg, ok := parsed["error"].(string); ok && errMsg != "" {
		return "", errors.New(errMsg)
placeholder
	if link, ok := parsed["download_link"].(string); ok {
		return link, nil
placeholder
	return "", errors.New("custom parse response missing download_link")
placeholder

const (
	soraSlotImageLock   int64 = 1
	soraSlotImageLimit  int64 = 2
	soraSlotVideoLimit  int64 = 3
	soraDefaultUsername       = "character"
)

func (s *SoraGatewayService) CallLogicMode(ctx context.Context) string {
	return strings.TrimSpace(s.getSoraConfig(ctx).CallLogicMode)
placeholder

func (s *SoraGatewayService) getSoraConfig(ctx context.Context) config.SoraConfig {
	if s.settingService != nil {
		return s.settingService.GetSoraConfig(ctx)
placeholder
	if s.cfg != nil {
		return s.cfg.Sora
placeholder
	return config.SoraConfig{placeholder
placeholder

func (s *SoraGatewayService) acquireSoraSlots(ctx context.Context, account *Account, soraAcc *SoraAccount, isVideo bool) (func(), error) {
	if s.concurrency == nil || account == nil || soraAcc == nil {
		return nil, nil
placeholder
	releases := make([]func(), 0, 2)
	appendRelease := func(release func()) {
		if release != nil {
			releases = append(releases, release)
	placeholder
placeholder
	// 错误时释放所有已获取的槽位
	releaseAll := func() {
		for _, r := range releases {
			r()
	placeholder
placeholder

	if isVideo {
		if soraAcc.VideoConcurrency > 0 {
			release, err := s.acquireSoraSlot(ctx, account.ID, soraAcc.VideoConcurrency, soraSlotVideoLimit)
			if err != nil {
				releaseAll()
				return nil, err
		placeholder
			appendRelease(release)
	placeholder
placeholder else {
		release, err := s.acquireSoraSlot(ctx, account.ID, 1, soraSlotImageLock)
		if err != nil {
			releaseAll()
			return nil, err
	placeholder
		appendRelease(release)
		if soraAcc.ImageConcurrency > 0 {
			release, err := s.acquireSoraSlot(ctx, account.ID, soraAcc.ImageConcurrency, soraSlotImageLimit)
			if err != nil {
				releaseAll() // 释放已获取的 soraSlotImageLock
				return nil, err
		placeholder
			appendRelease(release)
	placeholder
placeholder

	if len(releases) == 0 {
		return nil, nil
placeholder
	return func() {
		for _, release := range releases {
			release()
	placeholder
placeholder, nil
placeholder

func (s *SoraGatewayService) acquireSoraSlot(ctx context.Context, accountID int64, maxConcurrency int, slotType int64) (func(), error) {
	if s.concurrency == nil || maxConcurrency <= 0 {
		return nil, nil
placeholder
	derivedID := soraConcurrencyAccountID(accountID, slotType)
	result, err := s.concurrency.AcquireAccountSlot(ctx, derivedID, maxConcurrency)
	if err != nil {
		return nil, err
placeholder
	if !result.Acquired {
		return nil, ErrSoraAccountNotEligible
placeholder
	return result.ReleaseFunc, nil
placeholder

func soraConcurrencyAccountID(accountID int64, slotType int64) int64 {
	if accountID < 0 {
		accountID = -accountID
placeholder
	return -(accountID*10 + slotType)
placeholder

func (s *SoraGatewayService) createCharacter(ctx context.Context, client *sora.Client, opts sora.RequestOptions, rawVideo string) (string, string, error) {
	videoBytes, err := s.loadVideoBytes(ctx, opts, rawVideo)
	if err != nil {
		return "", "", err
placeholder
	cameoID, err := client.UploadCharacterVideo(ctx, opts, videoBytes)
	if err != nil {
		return "", "", err
placeholder
	status, err := s.pollCameoStatus(ctx, client, opts, cameoID)
	if err != nil {
		return "", "", err
placeholder
	username := processCharacterUsername(stringFromMap(status, "username_hint"))
	if username == "" {
		username = soraDefaultUsername
placeholder
	displayName := stringFromMap(status, "display_name_hint")
	if displayName == "" {
		displayName = "Character"
placeholder
	profileURL := stringFromMap(status, "profile_asset_url")
	if profileURL == "" {
		return "", "", errors.New("profile asset url missing")
placeholder
	avatarData, err := client.DownloadCharacterImage(ctx, opts, profileURL)
	if err != nil {
		return "", "", err
placeholder
	assetPointer, err := client.UploadCharacterImage(ctx, opts, avatarData)
	if err != nil {
		return "", "", err
placeholder
	characterID, err := client.FinalizeCharacter(ctx, opts, cameoID, username, displayName, assetPointer)
	if err != nil {
		return "", "", err
placeholder
	if err := client.SetCharacterPublic(ctx, opts, cameoID); err != nil {
		return "", "", err
placeholder
	return username, characterID, nil
placeholder

func (s *SoraGatewayService) pollCameoStatus(ctx context.Context, client *sora.Client, opts sora.RequestOptions, cameoID string) (map[string]any, error) {
	if cameoID == "" {
		return nil, errors.New("cameo id empty")
placeholder
	timeout := 600 * time.Second
	pollInterval := 5 * time.Second
	deadline := time.Now().Add(timeout)
	consecutiveErrors := 0
	maxConsecutiveErrors := 3

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
	placeholder
		time.Sleep(pollInterval)
		status, err := client.GetCameoStatus(ctx, opts, cameoID)
		if err != nil {
			consecutiveErrors++
			if consecutiveErrors >= maxConsecutiveErrors {
				return nil, err
		placeholder
			continue
	placeholder
		consecutiveErrors = 0
		statusValue := stringFromMap(status, "status")
		statusMessage := stringFromMap(status, "status_message")
		if statusValue == "failed" {
			if statusMessage == "" {
				statusMessage = "角色创建失败"
		placeholder
			return nil, fmt.Errorf("角色创建失败: %s", statusMessage)
	placeholder
		if strings.EqualFold(statusMessage, "Completed") || strings.EqualFold(statusValue, "finalized") {
			return status, nil
	placeholder
placeholder
	return nil, errors.New("角色创建超时")
placeholder

func (s *SoraGatewayService) loadVideoBytes(ctx context.Context, opts sora.RequestOptions, rawVideo string) ([]byte, error) {
	trimmed := strings.TrimSpace(rawVideo)
	if trimmed == "" {
		return nil, errors.New("video data is empty")
placeholder
	if looksLikeURL(trimmed) {
		if err := s.validateMediaURL(trimmed); err != nil {
			return nil, err
	placeholder
		return s.downloadMedia(ctx, opts, trimmed, maxVideoDownloadSize)
placeholder
	return decodeBase64(trimmed)
placeholder

func (s *SoraGatewayService) loadImageBytes(ctx context.Context, opts sora.RequestOptions, rawImage string) ([]byte, error) {
	trimmed := strings.TrimSpace(rawImage)
	if trimmed == "" {
		return nil, errors.New("image data is empty")
placeholder
	if looksLikeURL(trimmed) {
		if err := s.validateMediaURL(trimmed); err != nil {
			return nil, err
	placeholder
		return s.downloadMedia(ctx, opts, trimmed, maxImageDownloadSize)
placeholder
	return decodeBase64(trimmed)
placeholder

func (s *SoraGatewayService) validateMediaURL(rawURL string) error {
	cfg := s.cfg
	if cfg == nil {
		return nil
placeholder
	if cfg.Security.URLAllowlist.Enabled {
		_, err := urlvalidator.ValidateHTTPSURL(rawURL, urlvalidator.ValidationOptions{
			AllowedHosts:     cfg.Security.URLAllowlist.UpstreamHosts,
			RequireAllowlist: true,
			AllowPrivate:     cfg.Security.URLAllowlist.AllowPrivateHosts,
	placeholder)
		if err != nil {
			return fmt.Errorf("媒体地址不合法: %w", err)
	placeholder
		return nil
placeholder
	if _, err := urlvalidator.ValidateURLFormat(rawURL, cfg.Security.URLAllowlist.AllowInsecureHTTP); err != nil {
		return fmt.Errorf("媒体地址不合法: %w", err)
placeholder
	return nil
placeholder

func (s *SoraGatewayService) downloadMedia(ctx context.Context, opts sora.RequestOptions, mediaURL string, maxSize int64) ([]byte, error) {
	if s.httpUpstream == nil {
		return nil, errors.New("upstream is nil")
placeholder
	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	enableTLS := false
	if s.cfg != nil {
		enableTLS = s.cfg.Gateway.TLSFingerprint.Enabled
placeholder
	resp, err := s.httpUpstream.DoWithTLS(req, opts.ProxyURL, opts.AccountID, opts.AccountConcurrency, enableTLS)
	if err != nil {
		return nil, err
placeholder
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("下载失败: %d", resp.StatusCode)
placeholder

	// 使用 LimitReader 限制最大读取大小，防止 DoS 攻击
	limitedReader := io.LimitReader(resp.Body, maxSize+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
placeholder

	// 检查是否超过大小限制
	if int64(len(data)) > maxSize {
		return nil, fmt.Errorf("媒体文件过大 (最大 %d 字节, 实际 %d 字节)", maxSize, len(data))
placeholder

	return data, nil
placeholder

func processCharacterUsername(usernameHint string) string {
	trimmed := strings.TrimSpace(usernameHint)
	if trimmed == "" {
		return ""
placeholder
	base := trimmed
	if idx := strings.LastIndex(trimmed, "."); idx != -1 && idx+1 < len(trimmed) {
		base = trimmed[idx+1:]
placeholder
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s%d", base, rng.Intn(900)+100)
placeholder

func looksLikeURL(value string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://")
placeholder
