package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	openaioauth "github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/util/logredact"
	"github.com/google/uuid"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/sha3"
)

const (
	soraChatGPTBaseURL   = "https://chatgpt.com"
	soraSentinelFlow     = "sora_2_create_task"
	soraDefaultUserAgent = "Sora/1.2026.007 (Android 15; 24122RKC7C; build 2600700)"
)

var (
	soraSessionAuthURL = "https://sora.chatgpt.com/api/auth/session"
	soraOAuthTokenURL  = "https://auth.openai.com/oauth/token"
)

const (
	soraPowMaxIteration = 500000
)

var soraPowCores = []int{8, 16, 24, 32placeholder

var soraPowScripts = []string{
	"https://cdn.oaistatic.com/_next/static/cXh69klOLzS0Gy2joLDRS/_ssgManifest.js?dpl=453ebaec0d44c2decab71692e1bfe39be35a24b3",
placeholder

var soraPowDPL = []string{
	"prod-f501fe933b3edf57aea882da888e1a544df99840",
placeholder

var soraPowNavigatorKeys = []string{
	"registerProtocolHandler−function registerProtocolHandler() { [native code] placeholder",
	"storage−[object StorageManager]",
	"locks−[object LockManager]",
	"appCodeName−Mozilla",
	"permissions−[object Permissions]",
	"webdriver−false",
	"vendor−Google Inc.",
	"mediaDevices−[object MediaDevices]",
	"cookieEnabled−true",
	"product−Gecko",
	"productSub−20030107",
	"hardwareConcurrency−32",
	"onLine−true",
placeholder

var soraPowDocumentKeys = []string{
	"_reactListeningo743lnnpvdg",
	"location",
placeholder

var soraPowWindowKeys = []string{
	"0", "window", "self", "document", "name", "location",
	"navigator", "screen", "innerWidth", "innerHeight",
	"localStorage", "sessionStorage", "crypto", "performance",
	"fetch", "setTimeout", "setInterval", "console",
placeholder

var soraDesktopUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 11.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
placeholder

var soraMobileUserAgents = []string{
	"Sora/1.2026.007 (Android 15; 24122RKC7C; build 2600700)",
	"Sora/1.2026.007 (Android 14; SM-G998B; build 2600700)",
	"Sora/1.2026.007 (Android 15; Pixel 8 Pro; build 2600700)",
	"Sora/1.2026.007 (Android 14; Pixel 7; build 2600700)",
	"Sora/1.2026.007 (Android 15; 2211133C; build 2600700)",
	"Sora/1.2026.007 (Android 14; SM-S918B; build 2600700)",
	"Sora/1.2026.007 (Android 15; OnePlus 12; build 2600700)",
placeholder

var soraRand = rand.New(rand.NewSource(time.Now().UnixNano()))
var soraRandMu sync.Mutex
var soraPerfStart = time.Now()
var soraPowTokenGenerator = soraGetPowToken

// SoraClient 定义直连 Sora 的任务操作接口。
type SoraClient interface {
	Enabled() bool
	UploadImage(ctx context.Context, account *Account, data []byte, filename string) (string, error)
	CreateImageTask(ctx context.Context, account *Account, req SoraImageRequest) (string, error)
	CreateVideoTask(ctx context.Context, account *Account, req SoraVideoRequest) (string, error)
	CreateStoryboardTask(ctx context.Context, account *Account, req SoraStoryboardRequest) (string, error)
	UploadCharacterVideo(ctx context.Context, account *Account, data []byte) (string, error)
	GetCameoStatus(ctx context.Context, account *Account, cameoID string) (*SoraCameoStatus, error)
	DownloadCharacterImage(ctx context.Context, account *Account, imageURL string) ([]byte, error)
	UploadCharacterImage(ctx context.Context, account *Account, data []byte) (string, error)
	FinalizeCharacter(ctx context.Context, account *Account, req SoraCharacterFinalizeRequest) (string, error)
	SetCharacterPublic(ctx context.Context, account *Account, cameoID string) error
	DeleteCharacter(ctx context.Context, account *Account, characterID string) error
	PostVideoForWatermarkFree(ctx context.Context, account *Account, generationID string) (string, error)
	DeletePost(ctx context.Context, account *Account, postID string) error
	GetWatermarkFreeURLCustom(ctx context.Context, account *Account, parseURL, parseToken, postID string) (string, error)
	EnhancePrompt(ctx context.Context, account *Account, prompt, expansionLevel string, durationS int) (string, error)
	GetImageTask(ctx context.Context, account *Account, taskID string) (*SoraImageTaskStatus, error)
	GetVideoTask(ctx context.Context, account *Account, taskID string) (*SoraVideoTaskStatus, error)
placeholder

// SoraImageRequest 图片生成请求参数
type SoraImageRequest struct {
	Prompt  string
	Width   int
	Height  int
	MediaID string
placeholder

// SoraVideoRequest 视频生成请求参数
type SoraVideoRequest struct {
	Prompt        string
	Orientation   string
	Frames        int
	Model         string
	Size          string
	MediaID       string
	RemixTargetID string
	CameoIDs      []string
placeholder

// SoraStoryboardRequest 分镜视频生成请求参数
type SoraStoryboardRequest struct {
	Prompt      string
	Orientation string
	Frames      int
	Model       string
	Size        string
	MediaID     string
placeholder

// SoraImageTaskStatus 图片任务状态
type SoraImageTaskStatus struct {
	ID          string
	Status      string
	ProgressPct float64
	URLs        []string
	ErrorMsg    string
placeholder

// SoraVideoTaskStatus 视频任务状态
type SoraVideoTaskStatus struct {
	ID           string
	Status       string
	ProgressPct  int
	URLs         []string
	GenerationID string
	ErrorMsg     string
placeholder

// SoraCameoStatus 角色处理中间态
type SoraCameoStatus struct {
	Status             string
	StatusMessage      string
	DisplayNameHint    string
	UsernameHint       string
	ProfileAssetURL    string
	InstructionSetHint any
	InstructionSet     any
placeholder

// SoraCharacterFinalizeRequest 角色定稿请求参数
type SoraCharacterFinalizeRequest struct {
	CameoID             string
	Username            string
	DisplayName         string
	ProfileAssetPointer string
	InstructionSet      any
placeholder

// SoraUpstreamError 上游错误
type SoraUpstreamError struct {
	StatusCode int
	Message    string
	Headers    http.Header
	Body       []byte
placeholder

func (e *SoraUpstreamError) Error() string {
	if e == nil {
		return "sora upstream error"
placeholder
	if e.Message != "" {
		return fmt.Sprintf("sora upstream error: %d %s", e.StatusCode, e.Message)
placeholder
	return fmt.Sprintf("sora upstream error: %d", e.StatusCode)
placeholder

// SoraDirectClient 直连 Sora 实现
type SoraDirectClient struct {
	cfg             *config.Config
	httpUpstream    HTTPUpstream
	tokenProvider   *OpenAITokenProvider
	accountRepo     AccountRepository
	soraAccountRepo SoraAccountRepository
	baseURL         string
placeholder

// NewSoraDirectClient 创建 Sora 直连客户端
func NewSoraDirectClient(cfg *config.Config, httpUpstream HTTPUpstream, tokenProvider *OpenAITokenProvider) *SoraDirectClient {
	baseURL := ""
	if cfg != nil {
		rawBaseURL := strings.TrimRight(strings.TrimSpace(cfg.Sora.Client.BaseURL), "/")
		baseURL = normalizeSoraBaseURL(rawBaseURL)
		if rawBaseURL != "" && baseURL != rawBaseURL {
			log.Printf("[SoraClient] normalized base_url from %s to %s", sanitizeSoraLogURL(rawBaseURL), sanitizeSoraLogURL(baseURL))
	placeholder
placeholder
	return &SoraDirectClient{
		cfg:           cfg,
		httpUpstream:  httpUpstream,
		tokenProvider: tokenProvider,
		baseURL:       baseURL,
placeholder
placeholder

func (c *SoraDirectClient) SetAccountRepositories(accountRepo AccountRepository, soraAccountRepo SoraAccountRepository) {
	if c == nil {
		return
placeholder
	c.accountRepo = accountRepo
	c.soraAccountRepo = soraAccountRepo
placeholder

// Enabled 判断是否启用 Sora 直连
func (c *SoraDirectClient) Enabled() bool {
	if c == nil {
		return false
placeholder
	if strings.TrimSpace(c.baseURL) != "" {
		return true
placeholder
	if c.cfg == nil {
		return false
placeholder
	return strings.TrimSpace(normalizeSoraBaseURL(c.cfg.Sora.Client.BaseURL)) != ""
placeholder

// PreflightCheck 在创建任务前执行账号能力预检。
// 当前仅对视频模型执行 /nf/check 预检，用于提前识别额度耗尽或能力缺失。
func (c *SoraDirectClient) PreflightCheck(ctx context.Context, account *Account, requestedModel string, modelCfg SoraModelConfig) error {
	if modelCfg.Type != "video" {
		return nil
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Accept", "application/json")
	body, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodGet, c.buildURL("/nf/check"), headers, nil, false)
	if err != nil {
		var upstreamErr *SoraUpstreamError
		if errors.As(err, &upstreamErr) && upstreamErr.StatusCode == http.StatusNotFound {
			return &SoraUpstreamError{
				StatusCode: http.StatusForbidden,
				Message:    "当前账号未开通 Sora2 能力或无可用配额",
				Headers:    upstreamErr.Headers,
				Body:       upstreamErr.Body,
		placeholder
	placeholder
		return err
placeholder

	rateLimitReached := gjson.GetBytes(body, "rate_limit_and_credit_balance.rate_limit_reached").Bool()
	remaining := gjson.GetBytes(body, "rate_limit_and_credit_balance.estimated_num_videos_remaining")
	if rateLimitReached || (remaining.Exists() && remaining.Int() <= 0) {
		msg := "当前账号 Sora2 可用配额不足"
		if requestedModel != "" {
			msg = fmt.Sprintf("当前账号 %s 可用配额不足", requestedModel)
	placeholder
		return &SoraUpstreamError{
			StatusCode: http.StatusTooManyRequests,
			Message:    msg,
			Headers:    http.Header{placeholder,
	placeholder
placeholder
	return nil
placeholder

func (c *SoraDirectClient) UploadImage(ctx context.Context, account *Account, data []byte, filename string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty image data")
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	if filename == "" {
		filename = "image.png"
placeholder
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	contentType := mime.TypeByExtension(path.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
placeholder
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	partHeader.Set("Content-Type", contentType)
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
placeholder
	if _, err := part.Write(data); err != nil {
		return "", err
placeholder
	if err := writer.WriteField("file_name", filename); err != nil {
		return "", err
placeholder
	if err := writer.Close(); err != nil {
		return "", err
placeholder

	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", writer.FormDataContentType())

	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/uploads"), headers, &body, false)
	if err != nil {
		return "", err
placeholder
	id := strings.TrimSpace(gjson.GetBytes(respBody, "id").String())
	if id == "" {
		return "", errors.New("upload response missing id")
placeholder
	return id, nil
placeholder

func (c *SoraDirectClient) CreateImageTask(ctx context.Context, account *Account, req SoraImageRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	operation := "simple_compose"
	inpaintItems := []map[string]any{placeholder
	if strings.TrimSpace(req.MediaID) != "" {
		operation = "remix"
		inpaintItems = append(inpaintItems, map[string]any{
			"type":            "image",
			"frame_index":     0,
			"upload_media_id": req.MediaID,
	placeholder)
placeholder
	payload := map[string]any{
		"type":          "image_gen",
		"operation":     operation,
		"prompt":        req.Prompt,
		"width":         req.Width,
		"height":        req.Height,
		"n_variants":    1,
		"n_frames":      1,
		"inpaint_items": inpaintItems,
placeholder
	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	sentinel, err := c.generateSentinelToken(ctx, account, token, userAgent, proxyURL)
	if err != nil {
		return "", err
placeholder
	headers.Set("openai-sentinel-token", sentinel)

	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/video_gen"), headers, bytes.NewReader(body), true)
	if err != nil {
		return "", err
placeholder
	taskID := strings.TrimSpace(gjson.GetBytes(respBody, "id").String())
	if taskID == "" {
		return "", errors.New("image task response missing id")
placeholder
	return taskID, nil
placeholder

func (c *SoraDirectClient) CreateVideoTask(ctx context.Context, account *Account, req SoraVideoRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
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

	inpaintItems := []map[string]any{placeholder
	if strings.TrimSpace(req.MediaID) != "" {
		inpaintItems = append(inpaintItems, map[string]any{
			"kind":      "upload",
			"upload_id": req.MediaID,
	placeholder)
placeholder
	payload := map[string]any{
		"kind":          "video",
		"prompt":        req.Prompt,
		"orientation":   orientation,
		"size":          size,
		"n_frames":      nFrames,
		"model":         model,
		"inpaint_items": inpaintItems,
placeholder
	if strings.TrimSpace(req.RemixTargetID) != "" {
		payload["remix_target_id"] = req.RemixTargetID
		payload["cameo_ids"] = []string{placeholder
		payload["cameo_replacements"] = map[string]any{placeholder
placeholder else if len(req.CameoIDs) > 0 {
		payload["cameo_ids"] = req.CameoIDs
		payload["cameo_replacements"] = map[string]any{placeholder
placeholder

	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	sentinel, err := c.generateSentinelToken(ctx, account, token, userAgent, proxyURL)
	if err != nil {
		return "", err
placeholder
	headers.Set("openai-sentinel-token", sentinel)

	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/nf/create"), headers, bytes.NewReader(body), true)
	if err != nil {
		return "", err
placeholder
	taskID := strings.TrimSpace(gjson.GetBytes(respBody, "id").String())
	if taskID == "" {
		return "", errors.New("video task response missing id")
placeholder
	return taskID, nil
placeholder

func (c *SoraDirectClient) CreateStoryboardTask(ctx context.Context, account *Account, req SoraStoryboardRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
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

	inpaintItems := []map[string]any{placeholder
	if strings.TrimSpace(req.MediaID) != "" {
		inpaintItems = append(inpaintItems, map[string]any{
			"kind":      "upload",
			"upload_id": req.MediaID,
	placeholder)
placeholder
	payload := map[string]any{
		"kind":               "video",
		"prompt":             req.Prompt,
		"title":              "Draft your video",
		"orientation":        orientation,
		"size":               size,
		"n_frames":           nFrames,
		"storyboard_id":      nil,
		"inpaint_items":      inpaintItems,
		"remix_target_id":    nil,
		"model":              model,
		"metadata":           nil,
		"style_id":           nil,
		"cameo_ids":          nil,
		"cameo_replacements": nil,
		"audio_caption":      nil,
		"audio_transcript":   nil,
		"video_caption":      nil,
placeholder

	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	sentinel, err := c.generateSentinelToken(ctx, account, token, userAgent, proxyURL)
	if err != nil {
		return "", err
placeholder
	headers.Set("openai-sentinel-token", sentinel)

	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/nf/create/storyboard"), headers, bytes.NewReader(body), true)
	if err != nil {
		return "", err
placeholder
	taskID := strings.TrimSpace(gjson.GetBytes(respBody, "id").String())
	if taskID == "" {
		return "", errors.New("storyboard task response missing id")
placeholder
	return taskID, nil
placeholder

func (c *SoraDirectClient) UploadCharacterVideo(ctx context.Context, account *Account, data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty video data")
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="file"; filename="video.mp4"`)
	partHeader.Set("Content-Type", "video/mp4")
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
placeholder
	if _, err := part.Write(data); err != nil {
		return "", err
placeholder
	if err := writer.WriteField("timestamps", "0,3"); err != nil {
		return "", err
placeholder
	if err := writer.Close(); err != nil {
		return "", err
placeholder

	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", writer.FormDataContentType())
	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/characters/upload"), headers, &body, false)
	if err != nil {
		return "", err
placeholder
	cameoID := strings.TrimSpace(gjson.GetBytes(respBody, "id").String())
	if cameoID == "" {
		return "", errors.New("character upload response missing id")
placeholder
	return cameoID, nil
placeholder

func (c *SoraDirectClient) GetCameoStatus(ctx context.Context, account *Account, cameoID string) (*SoraCameoStatus, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)
	respBody, _, err := c.doRequestWithProxy(
		ctx,
		account,
		proxyURL,
		http.MethodGet,
		c.buildURL("/project_y/cameos/in_progress/"+strings.TrimSpace(cameoID)),
		headers,
		nil,
		false,
	)
	if err != nil {
		return nil, err
placeholder
	return &SoraCameoStatus{
		Status:             strings.TrimSpace(gjson.GetBytes(respBody, "status").String()),
		StatusMessage:      strings.TrimSpace(gjson.GetBytes(respBody, "status_message").String()),
		DisplayNameHint:    strings.TrimSpace(gjson.GetBytes(respBody, "display_name_hint").String()),
		UsernameHint:       strings.TrimSpace(gjson.GetBytes(respBody, "username_hint").String()),
		ProfileAssetURL:    strings.TrimSpace(gjson.GetBytes(respBody, "profile_asset_url").String()),
		InstructionSetHint: gjson.GetBytes(respBody, "instruction_set_hint").Value(),
		InstructionSet:     gjson.GetBytes(respBody, "instruction_set").Value(),
placeholder, nil
placeholder

func (c *SoraDirectClient) DownloadCharacterImage(ctx context.Context, account *Account, imageURL string) ([]byte, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Accept", "image/*,*/*;q=0.8")

	respBody, _, err := c.doRequestWithProxy(
		ctx,
		account,
		proxyURL,
		http.MethodGet,
		strings.TrimSpace(imageURL),
		headers,
		nil,
		false,
	)
	if err != nil {
		return nil, err
placeholder
	return respBody, nil
placeholder

func (c *SoraDirectClient) UploadCharacterImage(ctx context.Context, account *Account, data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty character image")
placeholder
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="file"; filename="profile.webp"`)
	partHeader.Set("Content-Type", "image/webp")
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
placeholder
	if _, err := part.Write(data); err != nil {
		return "", err
placeholder
	if err := writer.WriteField("use_case", "profile"); err != nil {
		return "", err
placeholder
	if err := writer.Close(); err != nil {
		return "", err
placeholder

	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", writer.FormDataContentType())
	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/project_y/file/upload"), headers, &body, false)
	if err != nil {
		return "", err
placeholder
	assetPointer := strings.TrimSpace(gjson.GetBytes(respBody, "asset_pointer").String())
	if assetPointer == "" {
		return "", errors.New("character image upload response missing asset_pointer")
placeholder
	return assetPointer, nil
placeholder

func (c *SoraDirectClient) FinalizeCharacter(ctx context.Context, account *Account, req SoraCharacterFinalizeRequest) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	payload := map[string]any{
		"cameo_id":               req.CameoID,
		"username":               req.Username,
		"display_name":           req.DisplayName,
		"profile_asset_pointer":  req.ProfileAssetPointer,
		"instruction_set":        nil,
		"safety_instruction_set": nil,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/characters/finalize"), headers, bytes.NewReader(body), false)
	if err != nil {
		return "", err
placeholder
	characterID := strings.TrimSpace(gjson.GetBytes(respBody, "character.character_id").String())
	if characterID == "" {
		return "", errors.New("character finalize response missing character_id")
placeholder
	return characterID, nil
placeholder

func (c *SoraDirectClient) SetCharacterPublic(ctx context.Context, account *Account, cameoID string) error {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	payload := map[string]any{"visibility": "public"placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return err
placeholder
	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	_, _, err = c.doRequestWithProxy(
		ctx,
		account,
		proxyURL,
		http.MethodPost,
		c.buildURL("/project_y/cameos/by_id/"+strings.TrimSpace(cameoID)+"/update_v2"),
		headers,
		bytes.NewReader(body),
		false,
	)
	return err
placeholder

func (c *SoraDirectClient) DeleteCharacter(ctx context.Context, account *Account, characterID string) error {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)
	_, _, err = c.doRequestWithProxy(
		ctx,
		account,
		proxyURL,
		http.MethodDelete,
		c.buildURL("/project_y/characters/"+strings.TrimSpace(characterID)),
		headers,
		nil,
		false,
	)
	return err
placeholder

func (c *SoraDirectClient) PostVideoForWatermarkFree(ctx context.Context, account *Account, generationID string) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	payload := map[string]any{
		"attachments_to_create": []map[string]any{
			{
				"generation_id": generationID,
				"kind":          "sora",
		placeholder,
	placeholder,
		"post_text": "",
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	sentinel, err := c.generateSentinelToken(ctx, account, token, userAgent, proxyURL)
	if err != nil {
		return "", err
placeholder
	headers.Set("openai-sentinel-token", sentinel)
	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/project_y/post"), headers, bytes.NewReader(body), true)
	if err != nil {
		return "", err
placeholder
	postID := strings.TrimSpace(gjson.GetBytes(respBody, "post.id").String())
	if postID == "" {
		return "", errors.New("watermark-free publish response missing post.id")
placeholder
	return postID, nil
placeholder

func (c *SoraDirectClient) DeletePost(ctx context.Context, account *Account, postID string) error {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)
	_, _, err = c.doRequestWithProxy(
		ctx,
		account,
		proxyURL,
		http.MethodDelete,
		c.buildURL("/project_y/post/"+strings.TrimSpace(postID)),
		headers,
		nil,
		false,
	)
	return err
placeholder

func (c *SoraDirectClient) GetWatermarkFreeURLCustom(ctx context.Context, account *Account, parseURL, parseToken, postID string) (string, error) {
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

func (c *SoraDirectClient) EnhancePrompt(ctx context.Context, account *Account, prompt, expansionLevel string, durationS int) (string, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return "", err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	if strings.TrimSpace(expansionLevel) == "" {
		expansionLevel = "medium"
placeholder
	if durationS <= 0 {
		durationS = 10
placeholder

	payload := map[string]any{
		"prompt":          prompt,
		"expansion_level": expansionLevel,
		"duration_s":      durationS,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder

	headers := c.buildBaseHeaders(token, userAgent)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")

	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, c.buildURL("/editor/enhance_prompt"), headers, bytes.NewReader(body), false)
	if err != nil {
		return "", err
placeholder
	enhancedPrompt := strings.TrimSpace(gjson.GetBytes(respBody, "enhanced_prompt").String())
	if enhancedPrompt == "" {
		return "", errors.New("enhance_prompt response missing enhanced_prompt")
placeholder
	return enhancedPrompt, nil
placeholder

func (c *SoraDirectClient) GetImageTask(ctx context.Context, account *Account, taskID string) (*SoraImageTaskStatus, error) {
	status, found, err := c.fetchRecentImageTask(ctx, account, taskID, c.recentTaskLimit())
	if err != nil {
		return nil, err
placeholder
	if found {
		return status, nil
placeholder
	maxLimit := c.recentTaskLimitMax()
	if maxLimit > 0 && maxLimit != c.recentTaskLimit() {
		status, found, err = c.fetchRecentImageTask(ctx, account, taskID, maxLimit)
		if err != nil {
			return nil, err
	placeholder
		if found {
			return status, nil
	placeholder
placeholder
	return &SoraImageTaskStatus{ID: taskID, Status: "processing"placeholder, nil
placeholder

func (c *SoraDirectClient) fetchRecentImageTask(ctx context.Context, account *Account, taskID string, limit int) (*SoraImageTaskStatus, bool, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, false, err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)
	if limit <= 0 {
		limit = 20
placeholder
	endpoint := fmt.Sprintf("/v2/recent_tasks?limit=%d", limit)
	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodGet, c.buildURL(endpoint), headers, nil, false)
	if err != nil {
		return nil, false, err
placeholder
	var found *SoraImageTaskStatus
	gjson.GetBytes(respBody, "task_responses").ForEach(func(_, item gjson.Result) bool {
		if item.Get("id").String() != taskID {
			return true // continue
	placeholder
		status := strings.TrimSpace(item.Get("status").String())
		progress := item.Get("progress_pct").Float()
		var urls []string
		item.Get("generations").ForEach(func(_, gen gjson.Result) bool {
			if u := strings.TrimSpace(gen.Get("url").String()); u != "" {
				urls = append(urls, u)
		placeholder
			return true
	placeholder)
		found = &SoraImageTaskStatus{
			ID:          taskID,
			Status:      status,
			ProgressPct: progress,
			URLs:        urls,
	placeholder
		return false // break
placeholder)
	if found != nil {
		return found, true, nil
placeholder
	return &SoraImageTaskStatus{ID: taskID, Status: "processing"placeholder, false, nil
placeholder

func (c *SoraDirectClient) recentTaskLimit() int {
	if c == nil || c.cfg == nil {
		return 20
placeholder
	if c.cfg.Sora.Client.RecentTaskLimit > 0 {
		return c.cfg.Sora.Client.RecentTaskLimit
placeholder
	return 20
placeholder

func (c *SoraDirectClient) recentTaskLimitMax() int {
	if c == nil || c.cfg == nil {
		return 0
placeholder
	if c.cfg.Sora.Client.RecentTaskLimitMax > 0 {
		return c.cfg.Sora.Client.RecentTaskLimitMax
placeholder
	return 0
placeholder

func (c *SoraDirectClient) GetVideoTask(ctx context.Context, account *Account, taskID string) (*SoraVideoTaskStatus, error) {
	token, err := c.getAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	userAgent := c.taskUserAgent()
	proxyURL := c.resolveProxyURL(account)
	headers := c.buildBaseHeaders(token, userAgent)

	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodGet, c.buildURL("/nf/pending/v2"), headers, nil, false)
	if err != nil {
		return nil, err
placeholder
	// 搜索 pending 列表（JSON 数组）
	pendingResult := gjson.ParseBytes(respBody)
	if pendingResult.IsArray() {
		var pendingFound *SoraVideoTaskStatus
		pendingResult.ForEach(func(_, task gjson.Result) bool {
			if task.Get("id").String() != taskID {
				return true
		placeholder
			progress := 0
			if v := task.Get("progress_pct"); v.Exists() {
				progress = int(v.Float() * 100)
		placeholder
			status := strings.TrimSpace(task.Get("status").String())
			pendingFound = &SoraVideoTaskStatus{
				ID:          taskID,
				Status:      status,
				ProgressPct: progress,
		placeholder
			return false
	placeholder)
		if pendingFound != nil {
			return pendingFound, nil
	placeholder
placeholder

	respBody, _, err = c.doRequestWithProxy(ctx, account, proxyURL, http.MethodGet, c.buildURL("/project_y/profile/drafts?limit=15"), headers, nil, false)
	if err != nil {
		return nil, err
placeholder
	var draftFound *SoraVideoTaskStatus
	gjson.GetBytes(respBody, "items").ForEach(func(_, draft gjson.Result) bool {
		if draft.Get("task_id").String() != taskID {
			return true
	placeholder
		generationID := strings.TrimSpace(draft.Get("id").String())
		kind := strings.TrimSpace(draft.Get("kind").String())
		reason := strings.TrimSpace(draft.Get("reason_str").String())
		if reason == "" {
			reason = strings.TrimSpace(draft.Get("markdown_reason_str").String())
	placeholder
		urlStr := strings.TrimSpace(draft.Get("downloadable_url").String())
		if urlStr == "" {
			urlStr = strings.TrimSpace(draft.Get("url").String())
	placeholder

		if kind == "sora_content_violation" || reason != "" || urlStr == "" {
			msg := reason
			if msg == "" {
				msg = "Content violates guardrails"
		placeholder
			draftFound = &SoraVideoTaskStatus{
				ID:           taskID,
				Status:       "failed",
				GenerationID: generationID,
				ErrorMsg:     msg,
		placeholder
	placeholder else {
			draftFound = &SoraVideoTaskStatus{
				ID:           taskID,
				Status:       "completed",
				GenerationID: generationID,
				URLs:         []string{urlStrplaceholder,
		placeholder
	placeholder
		return false
placeholder)
	if draftFound != nil {
		return draftFound, nil
placeholder

	return &SoraVideoTaskStatus{ID: taskID, Status: "processing"placeholder, nil
placeholder

func (c *SoraDirectClient) buildURL(endpoint string) string {
	base := strings.TrimRight(strings.TrimSpace(c.baseURL), "/")
	if base == "" && c != nil && c.cfg != nil {
		base = normalizeSoraBaseURL(c.cfg.Sora.Client.BaseURL)
		c.baseURL = base
placeholder
	if base == "" {
		return endpoint
placeholder
	if strings.HasPrefix(endpoint, "/") {
		return base + endpoint
placeholder
	return base + "/" + endpoint
placeholder

func (c *SoraDirectClient) defaultUserAgent() string {
	if c == nil || c.cfg == nil {
		return soraDefaultUserAgent
placeholder
	ua := strings.TrimSpace(c.cfg.Sora.Client.UserAgent)
	if ua == "" {
		return soraDefaultUserAgent
placeholder
	return ua
placeholder

func (c *SoraDirectClient) taskUserAgent() string {
	if c != nil && c.cfg != nil {
		if ua := strings.TrimSpace(c.cfg.Sora.Client.UserAgent); ua != "" {
			return ua
	placeholder
placeholder
	if len(soraMobileUserAgents) > 0 {
		return soraMobileUserAgents[soraRandInt(len(soraMobileUserAgents))]
placeholder
	if len(soraDesktopUserAgents) > 0 {
		return soraDesktopUserAgents[soraRandInt(len(soraDesktopUserAgents))]
placeholder
	return soraDefaultUserAgent
placeholder

func (c *SoraDirectClient) resolveProxyURL(account *Account) string {
	if account == nil || account.ProxyID == nil || account.Proxy == nil {
		return ""
placeholder
	return strings.TrimSpace(account.Proxy.URL())
placeholder

func (c *SoraDirectClient) getAccessToken(ctx context.Context, account *Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder

	allowProvider := c.allowOpenAITokenProvider(account)
	var providerErr error
	if allowProvider && c.tokenProvider != nil {
		token, err := c.tokenProvider.GetAccessToken(ctx, account)
		if err == nil && strings.TrimSpace(token) != "" {
			c.logTokenSource(account, "openai_token_provider")
			return token, nil
	placeholder
		providerErr = err
		if err != nil && c.debugEnabled() {
			c.debugLogf(
				"token_provider_failed account_id=%d platform=%s err=%s",
				account.ID,
				account.Platform,
				logredact.RedactText(err.Error()),
			)
	placeholder
placeholder
	token := strings.TrimSpace(account.GetCredential("access_token"))
	if token != "" {
		expiresAt := account.GetCredentialAsTime("expires_at")
		if expiresAt != nil && time.Until(*expiresAt) <= 2*time.Minute {
			refreshed, refreshErr := c.recoverAccessToken(ctx, account, "access_token_expiring")
			if refreshErr == nil && strings.TrimSpace(refreshed) != "" {
				c.logTokenSource(account, "refresh_token_recovered")
				return refreshed, nil
		placeholder
			if refreshErr != nil && c.debugEnabled() {
				c.debugLogf("token_refresh_before_use_failed account_id=%d err=%s", account.ID, logredact.RedactText(refreshErr.Error()))
		placeholder
	placeholder
		c.logTokenSource(account, "account_credentials")
		return token, nil
placeholder

	recovered, recoverErr := c.recoverAccessToken(ctx, account, "access_token_missing")
	if recoverErr == nil && strings.TrimSpace(recovered) != "" {
		c.logTokenSource(account, "session_or_refresh_recovered")
		return recovered, nil
placeholder
	if recoverErr != nil && c.debugEnabled() {
		c.debugLogf("token_recover_failed account_id=%d platform=%s err=%s", account.ID, account.Platform, logredact.RedactText(recoverErr.Error()))
placeholder
	if providerErr != nil {
		return "", providerErr
placeholder
	if c.tokenProvider != nil && !allowProvider {
		c.logTokenSource(account, "account_credentials(provider_disabled)")
placeholder
	return "", errors.New("access_token not found")
placeholder

func (c *SoraDirectClient) recoverAccessToken(ctx context.Context, account *Account, reason string) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder

	if sessionToken := strings.TrimSpace(account.GetCredential("session_token")); sessionToken != "" {
		accessToken, expiresAt, err := c.exchangeSessionToken(ctx, account, sessionToken)
		if err == nil && strings.TrimSpace(accessToken) != "" {
			c.applyRecoveredToken(ctx, account, accessToken, "", expiresAt, sessionToken)
			c.logTokenRecover(account, "session_token", reason, true, nil)
			return accessToken, nil
	placeholder
		c.logTokenRecover(account, "session_token", reason, false, err)
placeholder

	refreshToken := strings.TrimSpace(account.GetCredential("refresh_token"))
	if refreshToken == "" {
		return "", errors.New("session_token/refresh_token not found")
placeholder
	accessToken, newRefreshToken, expiresAt, err := c.exchangeRefreshToken(ctx, account, refreshToken)
	if err != nil {
		c.logTokenRecover(account, "refresh_token", reason, false, err)
		return "", err
placeholder
	if strings.TrimSpace(accessToken) == "" {
		return "", errors.New("refreshed access_token is empty")
placeholder
	c.applyRecoveredToken(ctx, account, accessToken, newRefreshToken, expiresAt, "")
	c.logTokenRecover(account, "refresh_token", reason, true, nil)
	return accessToken, nil
placeholder

func (c *SoraDirectClient) exchangeSessionToken(ctx context.Context, account *Account, sessionToken string) (string, string, error) {
	headers := http.Header{placeholder
	headers.Set("Cookie", "__Secure-next-auth.session-token="+sessionToken)
	headers.Set("Accept", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	headers.Set("User-Agent", c.defaultUserAgent())
	body, _, err := c.doRequest(ctx, account, http.MethodGet, soraSessionAuthURL, headers, nil, false)
	if err != nil {
		return "", "", err
placeholder
	accessToken := strings.TrimSpace(gjson.GetBytes(body, "accessToken").String())
	if accessToken == "" {
		return "", "", errors.New("session exchange missing accessToken")
placeholder
	expiresAt := strings.TrimSpace(gjson.GetBytes(body, "expires").String())
	return accessToken, expiresAt, nil
placeholder

func (c *SoraDirectClient) exchangeRefreshToken(ctx context.Context, account *Account, refreshToken string) (string, string, string, error) {
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

		formData := url.Values{placeholder
		formData.Set("client_id", clientID)
		formData.Set("grant_type", "refresh_token")
		formData.Set("refresh_token", refreshToken)
		formData.Set("redirect_uri", "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback")
		headers := http.Header{placeholder
		headers.Set("Accept", "application/json")
		headers.Set("Content-Type", "application/x-www-form-urlencoded")
		headers.Set("User-Agent", c.defaultUserAgent())

		respBody, _, err := c.doRequest(ctx, account, http.MethodPost, soraOAuthTokenURL, headers, strings.NewReader(formData.Encode()), false)
		if err != nil {
			lastErr = err
			if c.debugEnabled() {
				c.debugLogf("refresh_token_exchange_failed account_id=%d client_id=%s err=%s", account.ID, clientID, logredact.RedactText(err.Error()))
		placeholder
			continue
	placeholder
		accessToken := strings.TrimSpace(gjson.GetBytes(respBody, "access_token").String())
		if accessToken == "" {
			lastErr = errors.New("oauth refresh response missing access_token")
			continue
	placeholder
		newRefreshToken := strings.TrimSpace(gjson.GetBytes(respBody, "refresh_token").String())
		expiresIn := gjson.GetBytes(respBody, "expires_in").Int()
		expiresAt := ""
		if expiresIn > 0 {
			expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second).Format(time.RFC3339)
	placeholder
		return accessToken, newRefreshToken, expiresAt, nil
placeholder

	if lastErr != nil {
		return "", "", "", lastErr
placeholder
	return "", "", "", errors.New("no available client_id for refresh_token exchange")
placeholder

func (c *SoraDirectClient) applyRecoveredToken(ctx context.Context, account *Account, accessToken, refreshToken, expiresAt, sessionToken string) {
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
		if err := c.accountRepo.Update(ctx, account); err != nil {
			if c.debugEnabled() {
				c.debugLogf("persist_recovered_token_failed account_id=%d err=%s", account.ID, logredact.RedactText(err.Error()))
		placeholder
	placeholder
placeholder
	c.updateSoraAccountExtension(ctx, account, accessToken, refreshToken, sessionToken)
placeholder

func (c *SoraDirectClient) updateSoraAccountExtension(ctx context.Context, account *Account, accessToken, refreshToken, sessionToken string) {
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

func (c *SoraDirectClient) logTokenRecover(account *Account, source, reason string, success bool, err error) {
	if !c.debugEnabled() || account == nil {
		return
placeholder
	if success {
		c.debugLogf("token_recover_success account_id=%d platform=%s source=%s reason=%s", account.ID, account.Platform, source, reason)
		return
placeholder
	if err == nil {
		c.debugLogf("token_recover_failed account_id=%d platform=%s source=%s reason=%s", account.ID, account.Platform, source, reason)
		return
placeholder
	c.debugLogf("token_recover_failed account_id=%d platform=%s source=%s reason=%s err=%s", account.ID, account.Platform, source, reason, logredact.RedactText(err.Error()))
placeholder

func (c *SoraDirectClient) allowOpenAITokenProvider(account *Account) bool {
	if c == nil || c.tokenProvider == nil {
		return false
placeholder
	if account != nil && account.Platform == PlatformSora {
		return c.cfg != nil && c.cfg.Sora.Client.UseOpenAITokenProvider
placeholder
	return true
placeholder

func (c *SoraDirectClient) logTokenSource(account *Account, source string) {
	if !c.debugEnabled() || account == nil {
		return
placeholder
	c.debugLogf(
		"token_selected account_id=%d platform=%s account_type=%s source=%s",
		account.ID,
		account.Platform,
		account.Type,
		source,
	)
placeholder

func (c *SoraDirectClient) buildBaseHeaders(token, userAgent string) http.Header {
	headers := http.Header{placeholder
	if token != "" {
		headers.Set("Authorization", "Bearer "+token)
placeholder
	if userAgent != "" {
		headers.Set("User-Agent", userAgent)
placeholder
	if c != nil && c.cfg != nil {
		for key, value := range c.cfg.Sora.Client.Headers {
			if strings.EqualFold(key, "authorization") || strings.EqualFold(key, "openai-sentinel-token") {
				continue
		placeholder
			headers.Set(key, value)
	placeholder
placeholder
	return headers
placeholder

func (c *SoraDirectClient) doRequest(ctx context.Context, account *Account, method, urlStr string, headers http.Header, body io.Reader, allowRetry bool) ([]byte, http.Header, error) {
	return c.doRequestWithProxy(ctx, account, c.resolveProxyURL(account), method, urlStr, headers, body, allowRetry)
placeholder

func (c *SoraDirectClient) doRequestWithProxy(
	ctx context.Context,
	account *Account,
	proxyURL string,
	method,
	urlStr string,
	headers http.Header,
	body io.Reader,
	allowRetry bool,
) ([]byte, http.Header, error) {
	if strings.TrimSpace(urlStr) == "" {
		return nil, nil, errors.New("empty upstream url")
placeholder
	proxyURL = strings.TrimSpace(proxyURL)
	if proxyURL == "" {
		proxyURL = c.resolveProxyURL(account)
placeholder
	timeout := 0
	if c != nil && c.cfg != nil {
		timeout = c.cfg.Sora.Client.TimeoutSeconds
placeholder
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
placeholder
	maxRetries := 0
	if allowRetry && c != nil && c.cfg != nil {
		maxRetries = c.cfg.Sora.Client.MaxRetries
placeholder
	if maxRetries < 0 {
		maxRetries = 0
placeholder

	var bodyBytes []byte
	if body != nil {
		b, err := io.ReadAll(body)
		if err != nil {
			return nil, nil, err
	placeholder
		bodyBytes = b
placeholder

	attempts := maxRetries + 1
	authRecovered := false
	authRecoverExtraAttemptGranted := false
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if c.debugEnabled() {
			c.debugLogf(
				"request_start method=%s url=%s attempt=%d/%d timeout_s=%d body_bytes=%d proxy_bound=%t headers=%s",
				method,
				sanitizeSoraLogURL(urlStr),
				attempt,
				attempts,
				timeout,
				len(bodyBytes),
				proxyURL != "",
				formatSoraHeaders(headers),
			)
	placeholder

		var reader io.Reader
		if bodyBytes != nil {
			reader = bytes.NewReader(bodyBytes)
	placeholder
		req, err := http.NewRequestWithContext(ctx, method, urlStr, reader)
		if err != nil {
			return nil, nil, err
	placeholder
		req.Header = headers.Clone()
		start := time.Now()

		resp, err := c.doHTTP(req, proxyURL, account)
		if err != nil {
			lastErr = err
			if c.debugEnabled() {
				c.debugLogf(
					"request_transport_error method=%s url=%s attempt=%d/%d err=%s",
					method,
					sanitizeSoraLogURL(urlStr),
					attempt,
					attempts,
					logredact.RedactText(err.Error()),
				)
		placeholder
			if attempt < attempts && allowRetry {
				if c.debugEnabled() {
					c.debugLogf("request_retry_scheduled method=%s url=%s reason=transport_error next_attempt=%d/%d", method, sanitizeSoraLogURL(urlStr), attempt+1, attempts)
			placeholder
				c.sleepRetry(attempt)
				continue
		placeholder
			return nil, nil, err
	placeholder

		respBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, resp.Header, readErr
	placeholder

		if c.cfg != nil && c.cfg.Sora.Client.Debug {
			c.debugLogf(
				"response_received method=%s url=%s attempt=%d/%d status=%d cost=%s resp_bytes=%d resp_headers=%s",
				method,
				sanitizeSoraLogURL(urlStr),
				attempt,
				attempts,
				resp.StatusCode,
				time.Since(start),
				len(respBody),
				formatSoraHeaders(resp.Header),
			)
	placeholder

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			if !authRecovered && shouldAttemptSoraTokenRecover(resp.StatusCode, urlStr) && account != nil {
				if recovered, recoverErr := c.recoverAccessToken(ctx, account, fmt.Sprintf("upstream_status_%d", resp.StatusCode)); recoverErr == nil && strings.TrimSpace(recovered) != "" {
					headers.Set("Authorization", "Bearer "+recovered)
					authRecovered = true
					if attempt == attempts && !authRecoverExtraAttemptGranted {
						attempts++
						authRecoverExtraAttemptGranted = true
				placeholder
					if c.debugEnabled() {
						c.debugLogf("request_retry_with_recovered_token method=%s url=%s status=%d", method, sanitizeSoraLogURL(urlStr), resp.StatusCode)
				placeholder
					continue
			placeholder else if recoverErr != nil && c.debugEnabled() {
					c.debugLogf("request_recover_token_failed method=%s url=%s status=%d err=%s", method, sanitizeSoraLogURL(urlStr), resp.StatusCode, logredact.RedactText(recoverErr.Error()))
			placeholder
		placeholder
			if c.debugEnabled() {
				c.debugLogf(
					"response_non_success method=%s url=%s attempt=%d/%d status=%d body=%s",
					method,
					sanitizeSoraLogURL(urlStr),
					attempt,
					attempts,
					resp.StatusCode,
					summarizeSoraResponseBody(respBody, 512),
				)
		placeholder
			upstreamErr := c.buildUpstreamError(resp.StatusCode, resp.Header, respBody, urlStr)
			lastErr = upstreamErr
			if allowRetry && attempt < attempts && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500) {
				if c.debugEnabled() {
					c.debugLogf("request_retry_scheduled method=%s url=%s reason=status_%d next_attempt=%d/%d", method, sanitizeSoraLogURL(urlStr), resp.StatusCode, attempt+1, attempts)
			placeholder
				c.sleepRetry(attempt)
				continue
		placeholder
			return nil, resp.Header, upstreamErr
	placeholder
		return respBody, resp.Header, nil
placeholder
	if lastErr != nil {
		return nil, nil, lastErr
placeholder
	return nil, nil, errors.New("upstream retries exhausted")
placeholder

func shouldAttemptSoraTokenRecover(statusCode int, rawURL string) bool {
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		parsed, err := url.Parse(strings.TrimSpace(rawURL))
		if err != nil {
			return false
	placeholder
		host := strings.ToLower(parsed.Hostname())
		if host != "sora.chatgpt.com" && host != "chatgpt.com" {
			return false
	placeholder
		// 避免在 ST->AT 转换接口上递归触发 token 恢复导致死循环。
		path := strings.ToLower(strings.TrimSpace(parsed.Path))
		if path == "/api/auth/session" {
			return false
	placeholder
		return true
	default:
		return false
placeholder
placeholder

func (c *SoraDirectClient) doHTTP(req *http.Request, proxyURL string, account *Account) (*http.Response, error) {
	enableTLS := c == nil || c.cfg == nil || !c.cfg.Sora.Client.DisableTLSFingerprint
	if c.httpUpstream != nil {
		accountID := int64(0)
		accountConcurrency := 0
		if account != nil {
			accountID = account.ID
			accountConcurrency = account.Concurrency
	placeholder
		return c.httpUpstream.DoWithTLS(req, proxyURL, accountID, accountConcurrency, enableTLS)
placeholder
	return http.DefaultClient.Do(req)
placeholder

func (c *SoraDirectClient) sleepRetry(attempt int) {
	backoff := time.Duration(attempt*attempt) * time.Second
	if backoff > 10*time.Second {
		backoff = 10 * time.Second
placeholder
	time.Sleep(backoff)
placeholder

func (c *SoraDirectClient) buildUpstreamError(status int, headers http.Header, body []byte, requestURL string) error {
	msg := strings.TrimSpace(extractUpstreamErrorMessage(body))
	msg = sanitizeUpstreamErrorMessage(msg)
	if status == http.StatusNotFound && strings.Contains(strings.ToLower(msg), "not found") {
		if hint := soraBaseURLNotFoundHint(requestURL); hint != "" {
			msg = strings.TrimSpace(msg + " " + hint)
	placeholder
placeholder
	if msg == "" {
		msg = truncateForLog(body, 256)
placeholder
	return &SoraUpstreamError{
		StatusCode: status,
		Message:    msg,
		Headers:    headers,
		Body:       body,
placeholder
placeholder

func normalizeSoraBaseURL(raw string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(raw), "/")
	if trimmed == "" {
		return ""
placeholder
	parsed, err := url.Parse(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return trimmed
placeholder
	host := strings.ToLower(parsed.Hostname())
	if host != "sora.chatgpt.com" && host != "chatgpt.com" {
		return trimmed
placeholder
	pathVal := strings.TrimRight(strings.TrimSpace(parsed.Path), "/")
	switch pathVal {
	case "", "/":
		parsed.Path = "/backend"
	case "/backend-api":
		parsed.Path = "/backend"
placeholder
	return strings.TrimRight(parsed.String(), "/")
placeholder

func soraBaseURLNotFoundHint(requestURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(requestURL))
	if err != nil || parsed.Host == "" {
		return ""
placeholder
	host := strings.ToLower(parsed.Hostname())
	if host != "sora.chatgpt.com" && host != "chatgpt.com" {
		return ""
placeholder
	pathVal := strings.TrimSpace(parsed.Path)
	if strings.HasPrefix(pathVal, "/backend/") || pathVal == "/backend" {
		return ""
placeholder
	return "(请检查 sora.client.base_url，建议配置为 https://sora.chatgpt.com/backend)"
placeholder

func (c *SoraDirectClient) generateSentinelToken(ctx context.Context, account *Account, accessToken, userAgent, proxyURL string) (string, error) {
	reqID := uuid.NewString()
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		userAgent = c.taskUserAgent()
placeholder
	powToken := soraPowTokenGenerator(userAgent)
	payload := map[string]any{
		"p":    powToken,
		"flow": soraSentinelFlow,
		"id":   reqID,
placeholder
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
placeholder
	headers := http.Header{placeholder
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("Content-Type", "application/json")
	headers.Set("Origin", "https://sora.chatgpt.com")
	headers.Set("Referer", "https://sora.chatgpt.com/")
	headers.Set("User-Agent", userAgent)
	if accessToken != "" {
		headers.Set("Authorization", "Bearer "+accessToken)
placeholder

	urlStr := soraChatGPTBaseURL + "/backend-api/sentinel/req"
	respBody, _, err := c.doRequestWithProxy(ctx, account, proxyURL, http.MethodPost, urlStr, headers, bytes.NewReader(body), true)
	if err != nil {
		return "", err
placeholder
	var resp map[string]any
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return "", err
placeholder

	sentinel := soraBuildSentinelToken(soraSentinelFlow, reqID, powToken, resp, userAgent)
	if sentinel == "" {
		return "", errors.New("failed to build sentinel token")
placeholder
	return sentinel, nil
placeholder

func soraGetPowToken(userAgent string) string {
	configList := soraBuildPowConfig(userAgent)
	seed := strconv.FormatFloat(soraRandFloat(), 'f', -1, 64)
	difficulty := "0fffff"
	solution, _ := soraSolvePow(seed, difficulty, configList)
	return "gAAAAAC" + solution
placeholder

func soraRandFloat() float64 {
	soraRandMu.Lock()
	defer soraRandMu.Unlock()
	return soraRand.Float64()
placeholder

func soraRandInt(max int) int {
	if max <= 1 {
		return 0
placeholder
	soraRandMu.Lock()
	defer soraRandMu.Unlock()
	return soraRand.Intn(max)
placeholder

func soraBuildPowConfig(userAgent string) []any {
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" && len(soraDesktopUserAgents) > 0 {
		userAgent = soraDesktopUserAgents[0]
placeholder
	screenVal := soraStableChoiceInt([]int{
		1920 + 1080,
		2560 + 1440,
		1920 + 1200,
		2560 + 1600,
placeholder, userAgent+"|screen")
	perfMs := float64(time.Since(soraPerfStart).Milliseconds())
	wallMs := float64(time.Now().UnixNano()) / 1e6
	diff := wallMs - perfMs
	return []any{
		screenVal,
		soraPowParseTime(),
		4294705152,
		0,
		userAgent,
		soraStableChoice(soraPowScripts, userAgent+"|script"),
		soraStableChoice(soraPowDPL, userAgent+"|dpl"),
		"en-US",
		"en-US,es-US,en,es",
		0,
		soraStableChoice(soraPowNavigatorKeys, userAgent+"|navigator"),
		soraStableChoice(soraPowDocumentKeys, userAgent+"|document"),
		soraStableChoice(soraPowWindowKeys, userAgent+"|window"),
		perfMs,
		uuid.NewString(),
		"",
		soraStableChoiceInt(soraPowCores, userAgent+"|cores"),
		diff,
placeholder
placeholder

func soraStableChoice(items []string, seed string) string {
	if len(items) == 0 {
		return ""
placeholder
	idx := soraStableIndex(seed, len(items))
	return items[idx]
placeholder

func soraStableChoiceInt(items []int, seed string) int {
	if len(items) == 0 {
		return 0
placeholder
	idx := soraStableIndex(seed, len(items))
	return items[idx]
placeholder

func soraStableIndex(seed string, size int) int {
	if size <= 0 {
		return 0
placeholder
	h := fnv.New32a()
	_, _ = h.Write([]byte(seed))
	return int(h.Sum32() % uint32(size))
placeholder

func soraPowParseTime() string {
	loc := time.FixedZone("EST", -5*3600)
	return time.Now().In(loc).Format("Mon Jan 02 2006 15:04:05 GMT-0700 (Eastern Standard Time)")
placeholder

func soraSolvePow(seed, difficulty string, configList []any) (string, bool) {
	diffLen := len(difficulty) / 2
	target, err := hexDecodeString(difficulty)
	if err != nil {
		return "", false
placeholder
	seedBytes := []byte(seed)

	part1 := mustMarshalJSON(configList[:3])
	part2 := mustMarshalJSON(configList[4:9])
	part3 := mustMarshalJSON(configList[10:])

	staticPart1 := append(part1[:len(part1)-1], ',')
	staticPart2 := append([]byte(","), append(part2[1:len(part2)-1], ',')...)
	staticPart3 := append([]byte(","), part3[1:]...)

	for i := 0; i < soraPowMaxIteration; i++ {
		dynamicI := []byte(strconv.Itoa(i))
		dynamicJ := []byte(strconv.Itoa(i >> 1))
		finalJSON := make([]byte, 0, len(staticPart1)+len(dynamicI)+len(staticPart2)+len(dynamicJ)+len(staticPart3))
		finalJSON = append(finalJSON, staticPart1...)
		finalJSON = append(finalJSON, dynamicI...)
		finalJSON = append(finalJSON, staticPart2...)
		finalJSON = append(finalJSON, dynamicJ...)
		finalJSON = append(finalJSON, staticPart3...)

		b64 := base64.StdEncoding.EncodeToString(finalJSON)
		hash := sha3.Sum512(append(seedBytes, []byte(b64)...))
		if bytes.Compare(hash[:diffLen], target[:diffLen]) <= 0 {
			return b64, true
	placeholder
placeholder

	errorToken := "wQ8Lk5FbGpA2NcR9dShT6gYjU7VxZ4D" + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("\"%s\"", seed)))
	return errorToken, false
placeholder

func soraBuildSentinelToken(flow, reqID, powToken string, resp map[string]any, userAgent string) string {
	finalPow := powToken
	proof, _ := resp["proofofwork"].(map[string]any)
	if required, _ := proof["required"].(bool); required {
		seed, _ := proof["seed"].(string)
		difficulty, _ := proof["difficulty"].(string)
		if seed != "" && difficulty != "" {
			configList := soraBuildPowConfig(userAgent)
			solution, _ := soraSolvePow(seed, difficulty, configList)
			finalPow = "gAAAAAB" + solution
	placeholder
placeholder
	if !strings.HasSuffix(finalPow, "~S") {
		finalPow += "~S"
placeholder
	turnstile, _ := resp["turnstile"].(map[string]any)
	tokenPayload := map[string]any{
		"p":    finalPow,
		"t":    safeMapString(turnstile, "dx"),
		"c":    safeString(resp["token"]),
		"id":   reqID,
		"flow": flow,
placeholder
	encoded, _ := json.Marshal(tokenPayload)
	return string(encoded)
placeholder

func safeMapString(m map[string]any, key string) string {
	if m == nil {
		return ""
placeholder
	if v, ok := m[key]; ok {
		return safeString(v)
placeholder
	return ""
placeholder

func safeString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
placeholder
placeholder

func mustMarshalJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
placeholder

func hexDecodeString(s string) ([]byte, error) {
	dst := make([]byte, len(s)/2)
	_, err := hex.Decode(dst, []byte(s))
	return dst, err
placeholder

func sanitizeSoraLogURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
placeholder
	q := parsed.Query()
	q.Del("sig")
	q.Del("expires")
	parsed.RawQuery = q.Encode()
	return parsed.String()
placeholder

func (c *SoraDirectClient) debugEnabled() bool {
	return c != nil && c.cfg != nil && c.cfg.Sora.Client.Debug
placeholder

func (c *SoraDirectClient) debugLogf(format string, args ...any) {
	if !c.debugEnabled() {
		return
placeholder
	log.Printf("[SoraClient] "+format, args...)
placeholder

func formatSoraHeaders(headers http.Header) string {
	if len(headers) == 0 {
		return "{placeholder"
placeholder
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
placeholder
	sort.Strings(keys)
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		values := headers.Values(key)
		if len(values) == 0 {
			continue
	placeholder
		val := strings.Join(values, ",")
		if isSensitiveHeader(key) {
			out[key] = "***"
			continue
	placeholder
		out[key] = truncateForLog([]byte(logredact.RedactText(val)), 160)
placeholder
	encoded, err := json.Marshal(out)
	if err != nil {
		return "{placeholder"
placeholder
	return string(encoded)
placeholder

func isSensitiveHeader(key string) bool {
	k := strings.ToLower(strings.TrimSpace(key))
	switch k {
	case "authorization", "openai-sentinel-token", "cookie", "set-cookie", "x-api-key":
		return true
	default:
		return false
placeholder
placeholder

func summarizeSoraResponseBody(body []byte, maxLen int) string {
	if len(body) == 0 {
		return ""
placeholder
	var text string
	if json.Valid(body) {
		text = logredact.RedactJSON(body)
placeholder else {
		text = logredact.RedactText(string(body))
placeholder
	text = strings.TrimSpace(text)
	if maxLen <= 0 || len(text) <= maxLen {
		return text
placeholder
	return text[:maxLen] + "...(truncated)"
placeholder
