package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/DouDOU-start/go-sora2api/sora"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	openaioauth "github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/util/logredact"
	"github.com/tidwall/gjson"
)

// SoraSDKClient 基于 go-sora2api SDK 的 Sora 客户端实现。
// 它实现了 SoraClient 接口，用 SDK 替代原有的自建 HTTP/PoW/TLS 指纹逻辑。
type SoraSDKClient struct {
	cfg             *config.Config
	httpUpstream    HTTPUpstream
	tokenProvider   *OpenAITokenProvider
	accountRepo     AccountRepository
	soraAccountRepo SoraAccountRepository

	// 每个 proxyURL 对应一个 SDK 客户端实例
	sdkClients sync.Map // key: proxyURL (string), value: *sora.Client
placeholder

// NewSoraSDKClient 创建基于 SDK 的 Sora 客户端
func NewSoraSDKClient(cfg *config.Config, httpUpstream HTTPUpstream, tokenProvider *OpenAITokenProvider) *SoraSDKClient {
	return &SoraSDKClient{
		cfg:           cfg,
		httpUpstream:  httpUpstream,
		tokenProvider: tokenProvider,
placeholder
placeholder

// SetAccountRepositories 设置账号和 Sora 扩展仓库（用于 token 持久化）
func (c *SoraSDKClient) SetAccountRepositories(accountRepo AccountRepository, soraAccountRepo SoraAccountRepository) {
	if c == nil {
		return
placeholder
	c.accountRepo = accountRepo
	c.soraAccountRepo = soraAccountRepo
placeholder

// Enabled 判断是否启用 Sora
func (c *SoraSDKClient) Enabled() bool {
	if c == nil || c.cfg == nil {
		return false
placeholder
	return strings.TrimSpace(c.cfg.Sora.Client.BaseURL) != ""
placeholder

// PreflightCheck 在创建任务前执行账号能力预检。
// 当前仅对视频模型执行预检，用于提前识别额度耗尽或能力缺失。
func (c *SoraSDKClient) PreflightCheck(ctx context.Context, account *Account, requestedModel string, modelCfg SoraModelConfig) error {
	if modelCfg.Type != "video" {
		return nil
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return err
placeholder
	balance, err := sdkClient.GetCreditBalance(ctx, token)
	if err != nil {
		accountID := int64(0)
		if account != nil {
			accountID = account.ID
	placeholder
		logger.LegacyPrintf(
			"service.sora_sdk",
			"[PreflightCheckRawError] account_id=%d model=%s op=get_credit_balance raw_err=%s",
			accountID,
			requestedModel,
			logredact.RedactText(err.Error()),
		)
		return &SoraUpstreamError{
			StatusCode: http.StatusForbidden,
			Message:    "当前账号未开通 Sora2 能力或无可用配额",
	placeholder
placeholder
	if balance.RateLimitReached || balance.RemainingCount <= 0 {
		msg := "当前账号 Sora2 可用配额不足"
		if requestedModel != "" {
			msg = fmt.Sprintf("当前账号 %s 可用配额不足", requestedModel)
	placeholder
		return &SoraUpstreamError{
			StatusCode: http.StatusTooManyRequests,
			Message:    msg,
	placeholder
placeholder
	return nil
placeholder

func (c *SoraSDKClient) UploadImage(ctx context.Context, account *Account, data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty image data")
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	if filename == "" {
		filename = "image.png"
placeholder
	mediaID, err := sdkClient.UploadImage(ctx, token, data, filename)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return mediaID, nil
placeholder

func (c *SoraSDKClient) CreateImageTask(ctx context.Context, account *Account, req SoraImageRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	sentinel, err := sdkClient.GenerateSentinelToken(ctx, token)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	var taskID string
	if strings.TrimSpace(req.MediaID) != "" {
		taskID, err = sdkClient.CreateImageTaskWithImage(ctx, token, sentinel, req.Prompt, req.Width, req.Height, req.MediaID)
placeholder else {
		taskID, err = sdkClient.CreateImageTask(ctx, token, sentinel, req.Prompt, req.Width, req.Height)
placeholder
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return taskID, nil
placeholder

func (c *SoraSDKClient) CreateVideoTask(ctx context.Context, account *Account, req SoraVideoRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	sentinel, err := sdkClient.GenerateSentinelToken(ctx, token)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder

	orientation := req.Orientation
	if orientation == "" {
		orientation = "landscape"
placeholder
	nFrames := req.Frames
	if nFrames <= 0 {
		nFrames = 450
placeholder
	model := req.Model
	if model == "" {
		model = "sy_8"
placeholder
	size := req.Size
	if size == "" {
		size = "small"
placeholder
	videoCount := req.VideoCount
	if videoCount <= 0 {
		videoCount = 1
placeholder
	if videoCount > 3 {
		videoCount = 3
placeholder

	// Remix 模式
	if strings.TrimSpace(req.RemixTargetID) != "" {
		if videoCount > 1 {
			accountID := int64(0)
			if account != nil {
				accountID = account.ID
		placeholder
			c.debugLogf("video_count_ignored_for_remix account_id=%d count=%d", accountID, videoCount)
	placeholder
		styleID := "" // SDK ExtractStyle 可从 prompt 中提取
		taskID, err := sdkClient.RemixVideo(ctx, token, sentinel, req.RemixTargetID, req.Prompt, orientation, nFrames, styleID)
		if err != nil {
			return "", c.wrapSDKError(err, account)
	placeholder
		return taskID, nil
placeholder

	// 普通视频（文生视频或图生视频）
	var taskID string
	if videoCount <= 1 {
		taskID, err = sdkClient.CreateVideoTaskWithOptions(ctx, token, sentinel, req.Prompt, orientation, nFrames, model, size, req.MediaID, "")
placeholder else {
		taskID, err = c.createVideoTaskWithVariants(ctx, account, token, sentinel, req.Prompt, orientation, nFrames, model, size, req.MediaID, videoCount)
placeholder
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return taskID, nil
placeholder

func (c *SoraSDKClient) createVideoTaskWithVariants(
	ctx context.Context,
	account *Account,
	accessToken string,
	sentinelToken string,
	prompt string,
	orientation string,
	nFrames int,
	model string,
	size string,
	mediaID string,
	videoCount int,
) (string, error) {
	inpaintItems := make([]any, 0, 1)
	if strings.TrimSpace(mediaID) != "" {
		inpaintItems = append(inpaintItems, map[string]any{
			"kind":      "upload",
			"upload_id": mediaID,
	placeholder)
placeholder
	payload := map[string]any{
		"kind":          "video",
		"prompt":        prompt,
		"orientation":   orientation,
		"size":          size,
		"n_frames":      nFrames,
		"n_variants":    videoCount,
		"model":         model,
		"inpaint_items": inpaintItems,
		"style_id":      nil,
placeholder
	raw, err := c.doSoraBackendJSON(ctx, account, http.MethodPost, "/nf/create", accessToken, sentinelToken, payload)
	if err != nil {
		return "", err
placeholder
	taskID := strings.TrimSpace(gjson.GetBytes(raw, "id").String())
	if taskID == "" {
		return "", errors.New("create video task response missing id")
placeholder
	return taskID, nil
placeholder

func (c *SoraSDKClient) CreateStoryboardTask(ctx context.Context, account *Account, req SoraStoryboardRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	sentinel, err := sdkClient.GenerateSentinelToken(ctx, token)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder

	orientation := req.Orientation
	if orientation == "" {
		orientation = "landscape"
placeholder
	nFrames := req.Frames
	if nFrames <= 0 {
		nFrames = 450
placeholder

	taskID, err := sdkClient.CreateStoryboardTask(ctx, token, sentinel, req.Prompt, orientation, nFrames, req.MediaID, "")
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return taskID, nil
placeholder

func (c *SoraSDKClient) UploadCharacterVideo(ctx context.Context, account *Account, data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty video data")
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	cameoID, err := sdkClient.UploadCharacterVideo(ctx, token, data)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return cameoID, nil
placeholder

func (c *SoraSDKClient) GetCameoStatus(ctx context.Context, account *Account, cameoID string) (*SoraCameoStatus, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return nil, err
placeholder
	status, err := sdkClient.GetCameoStatus(ctx, token, cameoID)
	if err != nil {
		return nil, c.wrapSDKError(err, account)
placeholder
	return &SoraCameoStatus{
		Status:          status.Status,
		DisplayNameHint: status.DisplayNameHint,
		UsernameHint:    status.UsernameHint,
		ProfileAssetURL: status.ProfileAssetURL,
placeholder, nil
placeholder

func (c *SoraSDKClient) DownloadCharacterImage(ctx context.Context, account *Account, imageURL string) ([]byte, error) {
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return nil, err
placeholder
	data, err := sdkClient.DownloadCharacterImage(ctx, imageURL)
	if err != nil {
		return nil, c.wrapSDKError(err, account)
placeholder
	return data, nil
placeholder

func (c *SoraSDKClient) UploadCharacterImage(ctx context.Context, account *Account, data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty character image")
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	assetPointer, err := sdkClient.UploadCharacterImage(ctx, token, data)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return assetPointer, nil
placeholder

func (c *SoraSDKClient) FinalizeCharacter(ctx context.Context, account *Account, req SoraCharacterFinalizeRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	characterID, err := sdkClient.FinalizeCharacter(ctx, token, req.CameoID, req.Username, req.DisplayName, req.ProfileAssetPointer)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return characterID, nil
placeholder

func (c *SoraSDKClient) SetCharacterPublic(ctx context.Context, account *Account, cameoID string) error {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return err
placeholder
	if err := sdkClient.SetCharacterPublic(ctx, token, cameoID); err != nil {
		return c.wrapSDKError(err, account)
placeholder
	return nil
placeholder

func (c *SoraSDKClient) DeleteCharacter(ctx context.Context, account *Account, characterID string) error {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return err
placeholder
	if err := sdkClient.DeleteCharacter(ctx, token, characterID); err != nil {
		return c.wrapSDKError(err, account)
placeholder
	return nil
placeholder

func (c *SoraSDKClient) PostVideoForWatermarkFree(ctx context.Context, account *Account, generationID string) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	sentinel, err := sdkClient.GenerateSentinelToken(ctx, token)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	postID, err := sdkClient.PublishVideo(ctx, token, sentinel, generationID)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return postID, nil
placeholder

func (c *SoraSDKClient) DeletePost(ctx context.Context, account *Account, postID string) error {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return err
placeholder
	if err := sdkClient.DeletePost(ctx, token, postID); err != nil {
		return c.wrapSDKError(err, account)
placeholder
	return nil
placeholder

// GetWatermarkFreeURLCustom 使用自定义第三方解析服务获取去水印链接。
// SDK 不涉及此功能，保留自建实现。
func (c *SoraSDKClient) GetWatermarkFreeURLCustom(ctx context.Context, account *Account, parseURL, parseToken, postID string) (string, error) {
	parseURL = strings.TrimRight(strings.TrimSpace(parseURL), "/")
	if parseURL == "" {
		return "", errors.New("custom parse url is required")
placeholder
	if strings.TrimSpace(parseToken) == "" {
		return "", errors.New("custom parse token is required")
placeholder
	shareURL := "https://sora.chatgpt.com/p/" + strings.TrimSpace(postID)
	payload := map[string]any{
		"url":   shareURL,
		"token": strings.TrimSpace(parseToken),
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, parseURL+"/get-sora-link", bytes.NewReader(body))
	if err != nil {
		return "", err
placeholder
	req.Header.Set("Content-Type", "application/json")

	proxyURL := c.resolveProxyURL(account)
	accountID := int64(0)
	accountConcurrency := 0
	if account != nil {
		accountID = account.ID
		accountConcurrency = account.Concurrency
placeholder
	var resp *http.Response
	if c.httpUpstream != nil {
		resp, err = c.httpUpstream.Do(req, proxyURL, accountID, accountConcurrency)
placeholder else {
		resp, err = http.DefaultClient.Do(req)
placeholder
	if err != nil {
		return "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return "", err
placeholder
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("custom parse failed: %d %s", resp.StatusCode, truncateForLog(raw, 256))
placeholder
	downloadLink := strings.TrimSpace(gjson.GetBytes(raw, "download_link").String())
	if downloadLink == "" {
		return "", errors.New("custom parse response missing download_link")
placeholder
	return downloadLink, nil
placeholder

func (c *SoraSDKClient) EnhancePrompt(ctx context.Context, account *Account, prompt, expansionLevel string, durationS int) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder
	if strings.TrimSpace(expansionLevel) == "" {
		expansionLevel = "medium"
placeholder
	if durationS <= 0 {
		durationS = 10
placeholder
	enhanced, err := sdkClient.EnhancePrompt(ctx, token, prompt, expansionLevel, durationS)
	if err != nil {
		return "", c.wrapSDKError(err, account)
placeholder
	return enhanced, nil
placeholder

func (c *SoraSDKClient) GetImageTask(ctx context.Context, account *Account, taskID string) (*SoraImageTaskStatus, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return nil, err
placeholder
	result := sdkClient.QueryImageTaskOnce(ctx, token, taskID, time.Now().Add(-10*time.Second))
	if result.Err != nil {
		return &SoraImageTaskStatus{
			ID:       taskID,
			Status:   "failed",
			ErrorMsg: result.Err.Error(),
	placeholder, nil
placeholder
	if result.Done && result.ImageURL != "" {
		return &SoraImageTaskStatus{
			ID:     taskID,
			Status: "succeeded",
			URLs:   []string{result.ImageURLplaceholder,
	placeholder, nil
placeholder
	status := result.Progress.Status
	if status == "" {
		status = "processing"
placeholder
	return &SoraImageTaskStatus{
		ID:          taskID,
		Status:      status,
		ProgressPct: float64(result.Progress.Percent) / 100.0,
placeholder, nil
placeholder

func (c *SoraSDKClient) GetVideoTask(ctx context.Context, account *Account, taskID string) (*SoraVideoTaskStatus, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return nil, err
placeholder

	// 先查询 pending 列表
	result := sdkClient.QueryVideoTaskOnce(ctx, token, taskID, time.Now().Add(-10*time.Second), 0)
	if result.Err != nil {
		return &SoraVideoTaskStatus{
			ID:       taskID,
			Status:   "failed",
			ErrorMsg: result.Err.Error(),
	placeholder, nil
placeholder
	if !result.Done {
		return &SoraVideoTaskStatus{
			ID:          taskID,
			Status:      result.Progress.Status,
			ProgressPct: result.Progress.Percent,
	placeholder, nil
placeholder

	// 任务不在 pending 中，查询 drafts 获取下载链接
	downloadURLs, err := c.getVideoTaskDownloadURLs(ctx, account, token, taskID)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "内容违规") || strings.Contains(errMsg, "Content violates") {
			return &SoraVideoTaskStatus{
				ID:       taskID,
				Status:   "failed",
				ErrorMsg: errMsg,
		placeholder, nil
	placeholder
		// 可能还在处理中
		return &SoraVideoTaskStatus{
			ID:     taskID,
			Status: "processing",
	placeholder, nil
placeholder
	if len(downloadURLs) == 0 {
		return &SoraVideoTaskStatus{
			ID:     taskID,
			Status: "processing",
	placeholder, nil
placeholder
	return &SoraVideoTaskStatus{
		ID:     taskID,
		Status: "completed",
		URLs:   downloadURLs,
placeholder, nil
placeholder

func (c *SoraSDKClient) getVideoTaskDownloadURLs(ctx context.Context, account *Account, accessToken, taskID string) ([]string, error) {
	raw, err := c.doSoraBackendJSON(ctx, account, http.MethodGet, "/project_y/profile/drafts?limit=30", accessToken, "", nil)
	if err != nil {
		return nil, err
placeholder
	items := gjson.GetBytes(raw, "items")
	if !items.Exists() || !items.IsArray() {
		return nil, fmt.Errorf("drafts response missing items for task %s", taskID)
placeholder
	urlSet := make(map[string]struct{placeholder, 4)
	urls := make([]string, 0, 4)
	items.ForEach(func(_, item gjson.Result) bool {
		if strings.TrimSpace(item.Get("task_id").String()) != taskID {
			return true
	placeholder
		kind := strings.TrimSpace(item.Get("kind").String())
		reason := strings.TrimSpace(item.Get("reason_str").String())
		markdownReason := strings.TrimSpace(item.Get("markdown_reason_str").String())
		if kind == "sora_content_violation" || reason != "" || markdownReason != "" {
			if reason == "" {
				reason = markdownReason
		placeholder
			if reason == "" {
				reason = "内容违规"
		placeholder
			err = fmt.Errorf("内容违规: %s", reason)
			return false
	placeholder
		url := strings.TrimSpace(item.Get("downloadable_url").String())
		if url == "" {
			url = strings.TrimSpace(item.Get("url").String())
	placeholder
		if url == "" {
			return true
	placeholder
		if _, exists := urlSet[url]; exists {
			return true
	placeholder
		urlSet[url] = struct{placeholder{placeholder
		urls = append(urls, url)
		return true
placeholder)
	if err != nil {
		return nil, err
placeholder
	if len(urls) > 0 {
		return urls, nil
placeholder

	// 兼容旧 SDK 的兜底逻辑
	sdkClient, sdkErr := c.getSDKClient(account)
	if sdkErr != nil {
		return nil, sdkErr
placeholder
	downloadURL, sdkErr := sdkClient.GetDownloadURL(ctx, accessToken, taskID)
	if sdkErr != nil {
		return nil, sdkErr
placeholder
	if strings.TrimSpace(downloadURL) == "" {
		return nil, nil
placeholder
	return []string{downloadURLplaceholder, nil
placeholder

func (c *SoraSDKClient) doSoraBackendJSON(
	ctx context.Context,
	account *Account,
	method string,
	path string,
	accessToken string,
	sentinelToken string,
	payload map[string]any,
) ([]byte, error) {
	endpoint := "https://sora.chatgpt.com/backend" + path
	var body io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, err
	placeholder
		body = bytes.NewReader(raw)
placeholder

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
placeholder
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://sora.chatgpt.com")
	req.Header.Set("Referer", "https://sora.chatgpt.com/")
	req.Header.Set("User-Agent", "Sora/1.2026.007 (Android 15; 24122RKC7C; build 2600700)")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
placeholder
	if strings.TrimSpace(sentinelToken) != "" {
		req.Header.Set("openai-sentinel-token", sentinelToken)
placeholder

	proxyURL := c.resolveProxyURL(account)
	accountID := int64(0)
	accountConcurrency := 0
	if account != nil {
		accountID = account.ID
		accountConcurrency = account.Concurrency
placeholder

	var resp *http.Response
	if c.httpUpstream != nil {
		resp, err = c.httpUpstream.Do(req, proxyURL, accountID, accountConcurrency)
placeholder else {
		resp, err = http.DefaultClient.Do(req)
placeholder
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, err
placeholder
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncateForLog(raw, 256))
placeholder
	return raw, nil
placeholder

// --- 内部方法 ---

// getSDKClient 获取或创建指定代理的 SDK 客户端实例
func (c *SoraSDKClient) getSDKClient(account *Account) (*sora.Client, error) {
	proxyURL := c.resolveProxyURL(account)
	if v, ok := c.sdkClients.Load(proxyURL); ok {
		if cli, ok2 := v.(*sora.Client); ok2 {
			return cli, nil
	placeholder
placeholder
	client, err := sora.New(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("创建 Sora SDK 客户端失败: %w", err)
placeholder
	actual, _ := c.sdkClients.LoadOrStore(proxyURL, client)
	if cli, ok := actual.(*sora.Client); ok {
		return cli, nil
placeholder
	return client, nil
placeholder

func (c *SoraSDKClient) resolveProxyURL(account *Account) string {
	if account == nil || account.ProxyID == nil || account.Proxy == nil {
		return ""
placeholder
	return strings.TrimSpace(account.Proxy.URL())
placeholder

// getAccessToken 获取账号的 access_token，支持多种 token 来源和自动刷新。
// 此方法保留了原 SoraDirectClient 的 token 管理逻辑。
func (c *SoraSDKClient) getAccessToken(ctx context.Context, account *Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder

	// 优先尝试 OpenAI Token Provider
	allowProvider := c.allowOpenAITokenProvider(account)
	var providerErr error
	if allowProvider && c.tokenProvider != nil {
		token, err := c.tokenProvider.GetAccessToken(ctx, account)
		if err == nil && strings.TrimSpace(token) != "" {
			c.debugLogf("token_selected account_id=%d source=openai_token_provider", account.ID)
			return token, nil
	placeholder
		providerErr = err
		if err != nil && c.debugEnabled() {
			c.debugLogf("token_provider_failed account_id=%d err=%s", account.ID, logredact.RedactText(err.Error()))
	placeholder
placeholder

	// 尝试直接使用 credentials 中的 access_token
	token := strings.TrimSpace(account.GetCredential("access_token"))
	if token != "" {
		expiresAt := account.GetCredentialAsTime("expires_at")
		if expiresAt != nil && time.Until(*expiresAt) <= 2*time.Minute {
			refreshed, refreshErr := c.recoverAccessToken(ctx, account, "access_token_expiring")
			if refreshErr == nil && strings.TrimSpace(refreshed) != "" {
				return refreshed, nil
		placeholder
	placeholder
		return token, nil
placeholder

	// 尝试通过 session_token 或 refresh_token 恢复
	recovered, recoverErr := c.recoverAccessToken(ctx, account, "access_token_missing")
	if recoverErr == nil && strings.TrimSpace(recovered) != "" {
		return recovered, nil
placeholder
	if providerErr != nil {
		return "", providerErr
placeholder
	return "", errors.New("access_token not found")
placeholder

// recoverAccessToken 通过 session_token 或 refresh_token 恢复 access_token
func (c *SoraSDKClient) recoverAccessToken(ctx context.Context, account *Account, reason string) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder

	// 先尝试 session_token
	if sessionToken := strings.TrimSpace(account.GetCredential("session_token")); sessionToken != "" {
		accessToken, expiresAt, err := c.exchangeSessionToken(ctx, account, sessionToken)
		if err == nil && strings.TrimSpace(accessToken) != "" {
			c.applyRecoveredToken(ctx, account, accessToken, "", expiresAt, sessionToken)
			return accessToken, nil
	placeholder
placeholder

	// 再尝试 refresh_token
	refreshToken := strings.TrimSpace(account.GetCredential("refresh_token"))
	if refreshToken == "" {
		return "", errors.New("session_token/refresh_token not found")
placeholder

	sdkClient, err := c.getSDKClient(account)
	if err != nil {
		return "", err
placeholder

	// 尝试多个 client_id
	clientIDs := []string{
		strings.TrimSpace(account.GetCredential("client_id")),
		openaioauth.SoraClientID,
		openaioauth.ClientID,
placeholder
	tried := make(map[string]struct{placeholder, len(clientIDs))
	var lastErr error

	for _, clientID := range clientIDs {
		if clientID == "" {
			continue
	placeholder
		if _, ok := tried[clientID]; ok {
			continue
	placeholder
		tried[clientID] = struct{placeholder{placeholder

		newAccess, newRefresh, refreshErr := sdkClient.RefreshAccessToken(ctx, refreshToken, clientID)
		if refreshErr != nil {
			lastErr = refreshErr
			continue
	placeholder
		if strings.TrimSpace(newAccess) == "" {
			lastErr = errors.New("refreshed access_token is empty")
			continue
	placeholder
		c.applyRecoveredToken(ctx, account, newAccess, newRefresh, "", "")
		return newAccess, nil
placeholder

	if lastErr != nil {
		return "", lastErr
placeholder
	return "", errors.New("no available client_id for refresh_token exchange")
placeholder

// exchangeSessionToken 通过 session_token 换取 access_token
func (c *SoraSDKClient) exchangeSessionToken(ctx context.Context, account *Account, sessionToken string) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://sora.chatgpt.com/api/auth/session", nil)
	if err != nil {
		return "", "", err
placeholder
	req.Header.Set("Cookie", "__Secure-next-auth.session-token="+sessionToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Origin", "https://sora.chatgpt.com")
	req.Header.Set("Referer", "https://sora.chatgpt.com/")
	req.Header.Set("User-Agent", "Sora/1.2026.007 (Android 15; 24122RKC7C; build 2600700)")

	proxyURL := c.resolveProxyURL(account)
	accountID := int64(0)
	accountConcurrency := 0
	if account != nil {
		accountID = account.ID
		accountConcurrency = account.Concurrency
placeholder

	var resp *http.Response
	if c.httpUpstream != nil {
		resp, err = c.httpUpstream.Do(req, proxyURL, accountID, accountConcurrency)
placeholder else {
		resp, err = http.DefaultClient.Do(req)
placeholder
	if err != nil {
		return "", "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return "", "", err
placeholder
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("session exchange failed: %d", resp.StatusCode)
placeholder

	accessToken := strings.TrimSpace(gjson.GetBytes(body, "accessToken").String())
	if accessToken == "" {
		return "", "", errors.New("session exchange missing accessToken")
placeholder
	expiresAt := strings.TrimSpace(gjson.GetBytes(body, "expires").String())
	return accessToken, expiresAt, nil
placeholder

// applyRecoveredToken 将恢复的 token 写入账号内存和数据库
func (c *SoraSDKClient) applyRecoveredToken(ctx context.Context, account *Account, accessToken, refreshToken, expiresAt, sessionToken string) {
	if account == nil {
		return
placeholder
	if account.Credentials == nil {
		account.Credentials = make(map[string]any)
placeholder
	if strings.TrimSpace(accessToken) != "" {
		account.Credentials["access_token"] = accessToken
placeholder
	if strings.TrimSpace(refreshToken) != "" {
		account.Credentials["refresh_token"] = refreshToken
placeholder
	if strings.TrimSpace(expiresAt) != "" {
		account.Credentials["expires_at"] = expiresAt
placeholder
	if strings.TrimSpace(sessionToken) != "" {
		account.Credentials["session_token"] = sessionToken
placeholder

	if c.accountRepo != nil {
		if err := persistAccountCredentials(ctx, c.accountRepo, account, account.Credentials); err != nil && c.debugEnabled() {
			c.debugLogf("persist_recovered_token_failed account_id=%d err=%s", account.ID, logredact.RedactText(err.Error()))
	placeholder
placeholder
	c.updateSoraAccountExtension(ctx, account, accessToken, refreshToken, sessionToken)
placeholder

func (c *SoraSDKClient) updateSoraAccountExtension(ctx context.Context, account *Account, accessToken, refreshToken, sessionToken string) {
	if c == nil || c.soraAccountRepo == nil || account == nil || account.ID <= 0 {
		return
placeholder
	updates := make(map[string]any)
	if strings.TrimSpace(accessToken) != "" && strings.TrimSpace(refreshToken) != "" {
		updates["access_token"] = accessToken
		updates["refresh_token"] = refreshToken
placeholder
	if strings.TrimSpace(sessionToken) != "" {
		updates["session_token"] = sessionToken
placeholder
	if len(updates) == 0 {
		return
placeholder
	if err := c.soraAccountRepo.Upsert(ctx, account.ID, updates); err != nil && c.debugEnabled() {
		c.debugLogf("persist_sora_extension_failed account_id=%d err=%s", account.ID, logredact.RedactText(err.Error()))
placeholder
placeholder

func (c *SoraSDKClient) allowOpenAITokenProvider(account *Account) bool {
	if c == nil || c.tokenProvider == nil {
		return false
placeholder
	if account != nil && account.Platform == PlatformSora {
		return c.cfg != nil && c.cfg.Sora.Client.UseOpenAITokenProvider
placeholder
	return true
placeholder

// wrapSDKError 将 SDK 错误包装为 SoraUpstreamError
func (c *SoraSDKClient) wrapSDKError(err error, account *Account) error {
	if err == nil {
		return nil
placeholder
	msg := err.Error()
	statusCode := http.StatusBadGateway
	if strings.Contains(msg, "HTTP 401") || strings.Contains(msg, "HTTP 403") {
		statusCode = http.StatusUnauthorized
placeholder else if strings.Contains(msg, "HTTP 429") {
		statusCode = http.StatusTooManyRequests
placeholder else if strings.Contains(msg, "HTTP 404") {
		statusCode = http.StatusNotFound
placeholder
	accountID := int64(0)
	if account != nil {
		accountID = account.ID
placeholder
	logger.LegacyPrintf(
		"service.sora_sdk",
		"[WrapSDKError] account_id=%d mapped_status=%d raw_err=%s",
		accountID,
		statusCode,
		logredact.RedactText(msg),
	)
	return &SoraUpstreamError{
		StatusCode: statusCode,
		Message:    msg,
placeholder
placeholder

func (c *SoraSDKClient) debugEnabled() bool {
	return c != nil && c.cfg != nil && c.cfg.Sora.Client.Debug
placeholder

func (c *SoraSDKClient) debugLogf(format string, args ...any) {
	if c.debugEnabled() {
		log.Printf("[SoraSDK] "+format, args...)
placeholder
placeholder
