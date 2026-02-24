package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
)

const soraImageInputMaxBytes = 20 << 20
const soraImageInputMaxRedirects = 3
const soraImageInputTimeout = 20 * time.Second
const soraVideoInputMaxBytes = 200 << 20
const soraVideoInputMaxRedirects = 3
const soraVideoInputTimeout = 60 * time.Second

var soraImageSizeMap = map[string]string{
	"gpt-image":           "360",
	"gpt-image-landscape": "540",
	"gpt-image-portrait":  "540",
placeholder

var soraBlockedHostnames = map[string]struct{placeholder{
	"localhost":                 {placeholder,
	"localhost.localdomain":     {placeholder,
	"metadata.google.internal":  {placeholder,
	"metadata.google.internal.": {placeholder,
placeholder

var soraBlockedCIDRs = mustParseCIDRs([]string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"::/128",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
placeholder)

// SoraGatewayService handles forwarding requests to Sora upstream.
type SoraGatewayService struct {
	soraClient       SoraClient
	mediaStorage     *SoraMediaStorage
	rateLimitService *RateLimitService
	cfg              *config.Config
placeholder

type soraWatermarkOptions struct {
	Enabled           bool
	ParseMethod       string
	ParseURL          string
	ParseToken        string
	FallbackOnFailure bool
	DeletePost        bool
placeholder

type soraCharacterOptions struct {
	SetPublic           bool
	DeleteAfterGenerate bool
placeholder

type soraCharacterFlowResult struct {
	CameoID     string
	CharacterID string
	Username    string
	DisplayName string
placeholder

var soraStoryboardPattern = regexp.MustCompile(`\[\d+(?:\.\d+)?s\]`)
var soraStoryboardShotPattern = regexp.MustCompile(`\[(\d+(?:\.\d+)?)s\]\s*([^\[]+)`)
var soraRemixTargetPattern = regexp.MustCompile(`s_[a-f0-9]{32placeholder`)
var soraRemixTargetInURLPattern = regexp.MustCompile(`https://sora\.chatgpt\.com/p/s_[a-f0-9]{32placeholder`)

type soraPreflightChecker interface {
	PreflightCheck(ctx context.Context, account *Account, requestedModel string, modelCfg SoraModelConfig) error
placeholder

func NewSoraGatewayService(
	soraClient SoraClient,
	mediaStorage *SoraMediaStorage,
	rateLimitService *RateLimitService,
	cfg *config.Config,
) *SoraGatewayService {
	return &SoraGatewayService{
		soraClient:       soraClient,
		mediaStorage:     mediaStorage,
		rateLimitService: rateLimitService,
		cfg:              cfg,
placeholder
placeholder

func (s *SoraGatewayService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte, clientStream bool) (*ForwardResult, error) {
	startTime := time.Now()

	if s.soraClient == nil || !s.soraClient.Enabled() {
		if c != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": gin.H{
					"type":    "api_error",
					"message": "Sora 上游未配置",
			placeholder,
		placeholder)
	placeholder
		return nil, errors.New("sora upstream not configured")
placeholder

	var reqBody map[string]any
	if err := json.Unmarshal(body, &reqBody); err != nil {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "Failed to parse request body", clientStream)
		return nil, fmt.Errorf("parse request: %w", err)
placeholder
	reqModel, _ := reqBody["model"].(string)
	reqStream, _ := reqBody["stream"].(bool)
	if strings.TrimSpace(reqModel) == "" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "model is required", clientStream)
		return nil, errors.New("model is required")
placeholder

	mappedModel := account.GetMappedModel(reqModel)
	if mappedModel != "" && mappedModel != reqModel {
		reqModel = mappedModel
placeholder

	modelCfg, ok := GetSoraModelConfig(reqModel)
	if !ok {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "Unsupported Sora model", clientStream)
		return nil, fmt.Errorf("unsupported model: %s", reqModel)
placeholder
	prompt, imageInput, videoInput, remixTargetID := extractSoraInput(reqBody)
	prompt = strings.TrimSpace(prompt)
	imageInput = strings.TrimSpace(imageInput)
	videoInput = strings.TrimSpace(videoInput)
	remixTargetID = strings.TrimSpace(remixTargetID)

	if videoInput != "" && modelCfg.Type != "video" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "video input only supports video models", clientStream)
		return nil, errors.New("video input only supports video models")
placeholder
	if videoInput != "" && imageInput != "" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "image input and video input cannot be used together", clientStream)
		return nil, errors.New("image input and video input cannot be used together")
placeholder
	characterOnly := videoInput != "" && prompt == ""
	if modelCfg.Type == "prompt_enhance" && prompt == "" {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "prompt is required", clientStream)
		return nil, errors.New("prompt is required")
placeholder
	if modelCfg.Type != "prompt_enhance" && prompt == "" && !characterOnly {
		s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", "prompt is required", clientStream)
		return nil, errors.New("prompt is required")
placeholder

	reqCtx, cancel := s.withSoraTimeout(ctx, reqStream)
	if cancel != nil {
		defer cancel()
placeholder
	if checker, ok := s.soraClient.(soraPreflightChecker); ok && !characterOnly {
		if err := checker.PreflightCheck(reqCtx, account, reqModel, modelCfg); err != nil {
			return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
	placeholder
placeholder

	if modelCfg.Type == "prompt_enhance" {
		enhancedPrompt, err := s.soraClient.EnhancePrompt(reqCtx, account, prompt, modelCfg.ExpansionLevel, modelCfg.DurationS)
		if err != nil {
			return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
	placeholder
		content := strings.TrimSpace(enhancedPrompt)
		if content == "" {
			content = prompt
	placeholder
		var firstTokenMs *int
		if clientStream {
			ms, streamErr := s.writeSoraStream(c, reqModel, content, startTime)
			if streamErr != nil {
				return nil, streamErr
		placeholder
			firstTokenMs = ms
	placeholder else if c != nil {
			c.JSON(http.StatusOK, buildSoraNonStreamResponse(content, reqModel))
	placeholder
		return &ForwardResult{
			RequestID:    "",
			Model:        reqModel,
			Stream:       clientStream,
			Duration:     time.Since(startTime),
			FirstTokenMs: firstTokenMs,
			Usage:        ClaudeUsage{placeholder,
			MediaType:    "prompt",
	placeholder, nil
placeholder

	characterOpts := parseSoraCharacterOptions(reqBody)
	watermarkOpts := parseSoraWatermarkOptions(reqBody)
	var characterResult *soraCharacterFlowResult
	if videoInput != "" {
		videoData, videoErr := decodeSoraVideoInput(reqCtx, videoInput)
		if videoErr != nil {
			s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", videoErr.Error(), clientStream)
			return nil, videoErr
	placeholder
		characterResult, videoErr = s.createCharacterFromVideo(reqCtx, account, videoData, characterOpts)
		if videoErr != nil {
			return nil, s.handleSoraRequestError(ctx, account, videoErr, reqModel, c, clientStream)
	placeholder
		if characterResult != nil && characterOpts.DeleteAfterGenerate && strings.TrimSpace(characterResult.CharacterID) != "" && !characterOnly {
			characterID := strings.TrimSpace(characterResult.CharacterID)
			defer func() {
				cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancelCleanup()
				if err := s.soraClient.DeleteCharacter(cleanupCtx, account, characterID); err != nil {
					log.Printf("[Sora] cleanup character failed, character_id=%s err=%v", characterID, err)
			placeholder
		placeholder()
	placeholder
		if characterOnly {
			content := "角色创建成功"
			if characterResult != nil && strings.TrimSpace(characterResult.Username) != "" {
				content = fmt.Sprintf("角色创建成功，角色名@%s", strings.TrimSpace(characterResult.Username))
		placeholder
			var firstTokenMs *int
			if clientStream {
				ms, streamErr := s.writeSoraStream(c, reqModel, content, startTime)
				if streamErr != nil {
					return nil, streamErr
			placeholder
				firstTokenMs = ms
		placeholder else if c != nil {
				resp := buildSoraNonStreamResponse(content, reqModel)
				if characterResult != nil {
					resp["character_id"] = characterResult.CharacterID
					resp["cameo_id"] = characterResult.CameoID
					resp["character_username"] = characterResult.Username
					resp["character_display_name"] = characterResult.DisplayName
			placeholder
				c.JSON(http.StatusOK, resp)
		placeholder
			return &ForwardResult{
				RequestID:    "",
				Model:        reqModel,
				Stream:       clientStream,
				Duration:     time.Since(startTime),
				FirstTokenMs: firstTokenMs,
				Usage:        ClaudeUsage{placeholder,
				MediaType:    "prompt",
		placeholder, nil
	placeholder
		if characterResult != nil && strings.TrimSpace(characterResult.Username) != "" {
			prompt = fmt.Sprintf("@%s %s", characterResult.Username, prompt)
	placeholder
placeholder

	var imageData []byte
	imageFilename := ""
	if imageInput != "" {
		decoded, filename, err := decodeSoraImageInput(reqCtx, imageInput)
		if err != nil {
			s.writeSoraError(c, http.StatusBadRequest, "invalid_request_error", err.Error(), clientStream)
			return nil, err
	placeholder
		imageData = decoded
		imageFilename = filename
placeholder

	mediaID := ""
	if len(imageData) > 0 {
		uploadID, err := s.soraClient.UploadImage(reqCtx, account, imageData, imageFilename)
		if err != nil {
			return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
	placeholder
		mediaID = uploadID
placeholder

	taskID := ""
	var err error
	switch modelCfg.Type {
	case "image":
		taskID, err = s.soraClient.CreateImageTask(reqCtx, account, SoraImageRequest{
			Prompt:  prompt,
			Width:   modelCfg.Width,
			Height:  modelCfg.Height,
			MediaID: mediaID,
	placeholder)
	case "video":
		if remixTargetID == "" && isSoraStoryboardPrompt(prompt) {
			taskID, err = s.soraClient.CreateStoryboardTask(reqCtx, account, SoraStoryboardRequest{
				Prompt:      formatSoraStoryboardPrompt(prompt),
				Orientation: modelCfg.Orientation,
				Frames:      modelCfg.Frames,
				Model:       modelCfg.Model,
				Size:        modelCfg.Size,
				MediaID:     mediaID,
		placeholder)
	placeholder else {
			taskID, err = s.soraClient.CreateVideoTask(reqCtx, account, SoraVideoRequest{
				Prompt:        prompt,
				Orientation:   modelCfg.Orientation,
				Frames:        modelCfg.Frames,
				Model:         modelCfg.Model,
				Size:          modelCfg.Size,
				MediaID:       mediaID,
				RemixTargetID: remixTargetID,
				CameoIDs:      extractSoraCameoIDs(reqBody),
		placeholder)
	placeholder
	default:
		err = fmt.Errorf("unsupported model type: %s", modelCfg.Type)
placeholder
	if err != nil {
		return nil, s.handleSoraRequestError(ctx, account, err, reqModel, c, clientStream)
placeholder

	if clientStream && c != nil {
		s.prepareSoraStream(c, taskID)
placeholder

	var mediaURLs []string
	videoGenerationID := ""
	mediaType := modelCfg.Type
	imageCount := 0
	imageSize := ""
	switch modelCfg.Type {
	case "image":
		urls, pollErr := s.pollImageTask(reqCtx, c, account, taskID, clientStream)
		if pollErr != nil {
			return nil, s.handleSoraRequestError(ctx, account, pollErr, reqModel, c, clientStream)
	placeholder
		mediaURLs = urls
		imageCount = len(urls)
		imageSize = soraImageSizeFromModel(reqModel)
	case "video":
		videoStatus, pollErr := s.pollVideoTaskDetailed(reqCtx, c, account, taskID, clientStream)
		if pollErr != nil {
			return nil, s.handleSoraRequestError(ctx, account, pollErr, reqModel, c, clientStream)
	placeholder
		if videoStatus != nil {
			mediaURLs = videoStatus.URLs
			videoGenerationID = strings.TrimSpace(videoStatus.GenerationID)
	placeholder
	default:
		mediaType = "prompt"
placeholder

	watermarkPostID := ""
	if modelCfg.Type == "video" && watermarkOpts.Enabled {
		watermarkURL, postID, watermarkErr := s.resolveWatermarkFreeURL(reqCtx, account, videoGenerationID, watermarkOpts)
		if watermarkErr != nil {
			if !watermarkOpts.FallbackOnFailure {
				return nil, s.handleSoraRequestError(ctx, account, watermarkErr, reqModel, c, clientStream)
		placeholder
			log.Printf("[Sora] watermark-free fallback to original URL, task_id=%s err=%v", taskID, watermarkErr)
	placeholder else if strings.TrimSpace(watermarkURL) != "" {
			mediaURLs = []string{strings.TrimSpace(watermarkURL)placeholder
			watermarkPostID = strings.TrimSpace(postID)
	placeholder
placeholder

	finalURLs := s.normalizeSoraMediaURLs(mediaURLs)
	if len(mediaURLs) > 0 && s.mediaStorage != nil && s.mediaStorage.Enabled() {
		stored, storeErr := s.mediaStorage.StoreFromURLs(reqCtx, mediaType, mediaURLs)
		if storeErr != nil {
			// 存储失败时降级使用原始 URL，不中断用户请求
			log.Printf("[Sora] StoreFromURLs failed, falling back to original URLs: %v", storeErr)
	placeholder else {
			finalURLs = s.normalizeSoraMediaURLs(stored)
	placeholder
placeholder
	if watermarkPostID != "" && watermarkOpts.DeletePost {
		if deleteErr := s.soraClient.DeletePost(reqCtx, account, watermarkPostID); deleteErr != nil {
			log.Printf("[Sora] delete post failed, post_id=%s err=%v", watermarkPostID, deleteErr)
	placeholder
placeholder

	content := buildSoraContent(mediaType, finalURLs)
	var firstTokenMs *int
	if clientStream {
		ms, streamErr := s.writeSoraStream(c, reqModel, content, startTime)
		if streamErr != nil {
			return nil, streamErr
	placeholder
		firstTokenMs = ms
placeholder else if c != nil {
		response := buildSoraNonStreamResponse(content, reqModel)
		if len(finalURLs) > 0 {
			response["media_url"] = finalURLs[0]
			if len(finalURLs) > 1 {
				response["media_urls"] = finalURLs
		placeholder
	placeholder
		c.JSON(http.StatusOK, response)
placeholder

	return &ForwardResult{
		RequestID:    taskID,
		Model:        reqModel,
		Stream:       clientStream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		Usage:        ClaudeUsage{placeholder,
		MediaType:    mediaType,
		MediaURL:     firstMediaURL(finalURLs),
		ImageCount:   imageCount,
		ImageSize:    imageSize,
placeholder, nil
placeholder

func (s *SoraGatewayService) withSoraTimeout(ctx context.Context, stream bool) (context.Context, context.CancelFunc) {
	if s == nil || s.cfg == nil {
		return ctx, nil
placeholder
	timeoutSeconds := s.cfg.Gateway.SoraRequestTimeoutSeconds
	if stream {
		timeoutSeconds = s.cfg.Gateway.SoraStreamTimeoutSeconds
placeholder
	if timeoutSeconds <= 0 {
		return ctx, nil
placeholder
	return context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
placeholder

func parseSoraWatermarkOptions(body map[string]any) soraWatermarkOptions {
	opts := soraWatermarkOptions{
		Enabled:           parseBoolWithDefault(body, "watermark_free", false),
		ParseMethod:       strings.ToLower(strings.TrimSpace(parseStringWithDefault(body, "watermark_parse_method", "third_party"))),
		ParseURL:          strings.TrimSpace(parseStringWithDefault(body, "watermark_parse_url", "")),
		ParseToken:        strings.TrimSpace(parseStringWithDefault(body, "watermark_parse_token", "")),
		FallbackOnFailure: parseBoolWithDefault(body, "watermark_fallback_on_failure", true),
		DeletePost:        parseBoolWithDefault(body, "watermark_delete_post", false),
placeholder
	if opts.ParseMethod == "" {
		opts.ParseMethod = "third_party"
placeholder
	return opts
placeholder

func parseSoraCharacterOptions(body map[string]any) soraCharacterOptions {
	return soraCharacterOptions{
		SetPublic:           parseBoolWithDefault(body, "character_set_public", true),
		DeleteAfterGenerate: parseBoolWithDefault(body, "character_delete_after_generate", true),
placeholder
placeholder

func parseBoolWithDefault(body map[string]any, key string, def bool) bool {
	if body == nil {
		return def
placeholder
	val, ok := body[key]
	if !ok {
		return def
placeholder
	switch typed := val.(type) {
	case bool:
		return typed
	case int:
		return typed != 0
	case int32:
		return typed != 0
	case int64:
		return typed != 0
	case float64:
		return typed != 0
	case string:
		typed = strings.ToLower(strings.TrimSpace(typed))
		if typed == "true" || typed == "1" || typed == "yes" {
			return true
	placeholder
		if typed == "false" || typed == "0" || typed == "no" {
			return false
	placeholder
placeholder
	return def
placeholder

func parseStringWithDefault(body map[string]any, key, def string) string {
	if body == nil {
		return def
placeholder
	val, ok := body[key]
	if !ok {
		return def
placeholder
	if str, ok := val.(string); ok {
		return str
placeholder
	return def
placeholder

func extractSoraCameoIDs(body map[string]any) []string {
	if body == nil {
		return nil
placeholder
	raw, ok := body["cameo_ids"]
	if !ok {
		return nil
placeholder
	switch typed := raw.(type) {
	case []string:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			item = strings.TrimSpace(item)
			if item != "" {
				out = append(out, item)
		placeholder
	placeholder
		return out
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				continue
		placeholder
			str = strings.TrimSpace(str)
			if str != "" {
				out = append(out, str)
		placeholder
	placeholder
		return out
	default:
		return nil
placeholder
placeholder

func (s *SoraGatewayService) createCharacterFromVideo(ctx context.Context, account *Account, videoData []byte, opts soraCharacterOptions) (*soraCharacterFlowResult, error) {
	cameoID, err := s.soraClient.UploadCharacterVideo(ctx, account, videoData)
	if err != nil {
		return nil, err
placeholder

	cameoStatus, err := s.pollCameoStatus(ctx, account, cameoID)
	if err != nil {
		return nil, err
placeholder
	username := processSoraCharacterUsername(cameoStatus.UsernameHint)
	displayName := strings.TrimSpace(cameoStatus.DisplayNameHint)
	if displayName == "" {
		displayName = "Character"
placeholder
	profileAssetURL := strings.TrimSpace(cameoStatus.ProfileAssetURL)
	if profileAssetURL == "" {
		return nil, errors.New("profile asset url not found in cameo status")
placeholder

	avatarData, err := s.soraClient.DownloadCharacterImage(ctx, account, profileAssetURL)
	if err != nil {
		return nil, err
placeholder
	assetPointer, err := s.soraClient.UploadCharacterImage(ctx, account, avatarData)
	if err != nil {
		return nil, err
placeholder
	instructionSet := cameoStatus.InstructionSetHint
	if instructionSet == nil {
		instructionSet = cameoStatus.InstructionSet
placeholder

	characterID, err := s.soraClient.FinalizeCharacter(ctx, account, SoraCharacterFinalizeRequest{
		CameoID:             strings.TrimSpace(cameoID),
		Username:            username,
		DisplayName:         displayName,
		ProfileAssetPointer: assetPointer,
		InstructionSet:      instructionSet,
placeholder)
	if err != nil {
		return nil, err
placeholder

	if opts.SetPublic {
		if err := s.soraClient.SetCharacterPublic(ctx, account, cameoID); err != nil {
			return nil, err
	placeholder
placeholder

	return &soraCharacterFlowResult{
		CameoID:     strings.TrimSpace(cameoID),
		CharacterID: strings.TrimSpace(characterID),
		Username:    strings.TrimSpace(username),
		DisplayName: displayName,
placeholder, nil
placeholder

func (s *SoraGatewayService) pollCameoStatus(ctx context.Context, account *Account, cameoID string) (*SoraCameoStatus, error) {
	timeout := 10 * time.Minute
	interval := 5 * time.Second
	maxAttempts := int(math.Ceil(timeout.Seconds() / interval.Seconds()))
	if maxAttempts < 1 {
		maxAttempts = 1
placeholder

	var lastErr error
	consecutiveErrors := 0
	for attempt := 0; attempt < maxAttempts; attempt++ {
		status, err := s.soraClient.GetCameoStatus(ctx, account, cameoID)
		if err != nil {
			lastErr = err
			consecutiveErrors++
			if consecutiveErrors >= 3 {
				break
		placeholder
			if attempt < maxAttempts-1 {
				if sleepErr := sleepWithContext(ctx, interval); sleepErr != nil {
					return nil, sleepErr
			placeholder
		placeholder
			continue
	placeholder
		consecutiveErrors = 0
		if status == nil {
			if attempt < maxAttempts-1 {
				if sleepErr := sleepWithContext(ctx, interval); sleepErr != nil {
					return nil, sleepErr
			placeholder
		placeholder
			continue
	placeholder
		currentStatus := strings.ToLower(strings.TrimSpace(status.Status))
		statusMessage := strings.TrimSpace(status.StatusMessage)
		if currentStatus == "failed" {
			if statusMessage == "" {
				statusMessage = "character creation failed"
		placeholder
			return nil, errors.New(statusMessage)
	placeholder
		if strings.EqualFold(statusMessage, "Completed") || currentStatus == "finalized" {
			return status, nil
	placeholder
		if attempt < maxAttempts-1 {
			if sleepErr := sleepWithContext(ctx, interval); sleepErr != nil {
				return nil, sleepErr
		placeholder
	placeholder
placeholder
	if lastErr != nil {
		return nil, fmt.Errorf("poll cameo status failed: %w", lastErr)
placeholder
	return nil, errors.New("cameo processing timeout")
placeholder

func processSoraCharacterUsername(usernameHint string) string {
	usernameHint = strings.TrimSpace(usernameHint)
	if usernameHint == "" {
		usernameHint = "character"
placeholder
	if strings.Contains(usernameHint, ".") {
		parts := strings.Split(usernameHint, ".")
		usernameHint = strings.TrimSpace(parts[len(parts)-1])
placeholder
	if usernameHint == "" {
		usernameHint = "character"
placeholder
	return fmt.Sprintf("%s%d", usernameHint, soraRandInt(900)+100)
placeholder

func (s *SoraGatewayService) resolveWatermarkFreeURL(ctx context.Context, account *Account, generationID string, opts soraWatermarkOptions) (string, string, error) {
	generationID = strings.TrimSpace(generationID)
	if generationID == "" {
		return "", "", errors.New("generation id is required for watermark-free mode")
placeholder
	postID, err := s.soraClient.PostVideoForWatermarkFree(ctx, account, generationID)
	if err != nil {
		return "", "", err
placeholder
	postID = strings.TrimSpace(postID)
	if postID == "" {
		return "", "", errors.New("watermark-free publish returned empty post id")
placeholder

	switch opts.ParseMethod {
	case "custom":
		urlVal, parseErr := s.soraClient.GetWatermarkFreeURLCustom(ctx, account, opts.ParseURL, opts.ParseToken, postID)
		if parseErr != nil {
			return "", postID, parseErr
	placeholder
		return strings.TrimSpace(urlVal), postID, nil
	case "", "third_party":
		return fmt.Sprintf("https://oscdn2.dyysy.com/MP4/%s.mp4", postID), postID, nil
	default:
		return "", postID, fmt.Errorf("unsupported watermark parse method: %s", opts.ParseMethod)
placeholder
placeholder

func (s *SoraGatewayService) shouldFailoverUpstreamError(statusCode int) bool {
	switch statusCode {
	case 401, 402, 403, 404, 429, 529:
		return true
	default:
		return statusCode >= 500
placeholder
placeholder

func buildSoraNonStreamResponse(content, model string) map[string]any {
	return map[string]any{
		"id":      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		"object":  "chat.completion",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": content,
			placeholder,
				"finish_reason": "stop",
		placeholder,
	placeholder,
placeholder
placeholder

func soraImageSizeFromModel(model string) string {
	modelLower := strings.ToLower(model)
	if size, ok := soraImageSizeMap[modelLower]; ok {
		return size
placeholder
	if strings.Contains(modelLower, "landscape") || strings.Contains(modelLower, "portrait") {
		return "540"
placeholder
	return "360"
placeholder

func soraProErrorMessage(model, upstreamMsg string) string {
	modelLower := strings.ToLower(model)
	if strings.Contains(modelLower, "sora2pro-hd") {
		return "当前账号无法使用 Sora Pro-HD 模型，请更换模型或账号"
placeholder
	if strings.Contains(modelLower, "sora2pro") {
		return "当前账号无法使用 Sora Pro 模型，请更换模型或账号"
placeholder
	return ""
placeholder

func firstMediaURL(urls []string) string {
	if len(urls) == 0 {
		return ""
placeholder
	return urls[0]
placeholder

func (s *SoraGatewayService) buildSoraMediaURL(path string, rawQuery string) string {
	if path == "" {
		return path
placeholder
	prefix := "/sora/media"
	values := url.Values{placeholder
	if rawQuery != "" {
		if parsed, err := url.ParseQuery(rawQuery); err == nil {
			values = parsed
	placeholder
placeholder

	signKey := ""
	ttlSeconds := 0
	if s != nil && s.cfg != nil {
		signKey = strings.TrimSpace(s.cfg.Gateway.SoraMediaSigningKey)
		ttlSeconds = s.cfg.Gateway.SoraMediaSignedURLTTLSeconds
placeholder
	values.Del("sig")
	values.Del("expires")
	signingQuery := values.Encode()
	if signKey != "" && ttlSeconds > 0 {
		expires := time.Now().Add(time.Duration(ttlSeconds) * time.Second).Unix()
		signature := SignSoraMediaURL(path, signingQuery, expires, signKey)
		if signature != "" {
			values.Set("expires", strconv.FormatInt(expires, 10))
			values.Set("sig", signature)
			prefix = "/sora/media-signed"
	placeholder
placeholder

	encoded := values.Encode()
	if encoded == "" {
		return prefix + path
placeholder
	return prefix + path + "?" + encoded
placeholder

func (s *SoraGatewayService) prepareSoraStream(c *gin.Context, requestID string) {
	if c == nil {
		return
placeholder
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	if strings.TrimSpace(requestID) != "" {
		c.Header("x-request-id", requestID)
placeholder
placeholder

func (s *SoraGatewayService) writeSoraStream(c *gin.Context, model, content string, startTime time.Time) (*int, error) {
	if c == nil {
		return nil, nil
placeholder
	writer := c.Writer
	flusher, _ := writer.(http.Flusher)

	chunk := map[string]any{
		"id":      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index": 0,
				"delta": map[string]any{
					"content": content,
			placeholder,
		placeholder,
	placeholder,
placeholder
	encoded, _ := json.Marshal(chunk)
	if _, err := fmt.Fprintf(writer, "data: %s\n\n", encoded); err != nil {
		return nil, err
placeholder
	if flusher != nil {
		flusher.Flush()
placeholder
	ms := int(time.Since(startTime).Milliseconds())
	finalChunk := map[string]any{
		"id":      chunk["id"],
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []any{
			map[string]any{
				"index":         0,
				"delta":         map[string]any{placeholder,
				"finish_reason": "stop",
		placeholder,
	placeholder,
placeholder
	finalEncoded, _ := json.Marshal(finalChunk)
	if _, err := fmt.Fprintf(writer, "data: %s\n\n", finalEncoded); err != nil {
		return &ms, err
placeholder
	if _, err := fmt.Fprint(writer, "data: [DONE]\n\n"); err != nil {
		return &ms, err
placeholder
	if flusher != nil {
		flusher.Flush()
placeholder
	return &ms, nil
placeholder

func (s *SoraGatewayService) writeSoraError(c *gin.Context, status int, errType, message string, stream bool) {
	if c == nil {
		return
placeholder
	if stream {
		flusher, _ := c.Writer.(http.Flusher)
		errorData := map[string]any{
			"error": map[string]string{
				"type":    errType,
				"message": message,
		placeholder,
	placeholder
		jsonBytes, err := json.Marshal(errorData)
		if err != nil {
			_ = c.Error(err)
			return
	placeholder
		errorEvent := fmt.Sprintf("event: error\ndata: %s\n\n", string(jsonBytes))
		_, _ = fmt.Fprint(c.Writer, errorEvent)
		_, _ = fmt.Fprint(c.Writer, "data: [DONE]\n\n")
		if flusher != nil {
			flusher.Flush()
	placeholder
		return
placeholder
	c.JSON(status, gin.H{
		"error": gin.H{
			"type":    errType,
			"message": message,
	placeholder,
placeholder)
placeholder

func (s *SoraGatewayService) handleSoraRequestError(ctx context.Context, account *Account, err error, model string, c *gin.Context, stream bool) error {
	if err == nil {
		return nil
placeholder
	var upstreamErr *SoraUpstreamError
	if errors.As(err, &upstreamErr) {
		if s.rateLimitService != nil && account != nil {
			s.rateLimitService.HandleUpstreamError(ctx, account, upstreamErr.StatusCode, upstreamErr.Headers, upstreamErr.Body)
	placeholder
		if s.shouldFailoverUpstreamError(upstreamErr.StatusCode) {
			var responseHeaders http.Header
			if upstreamErr.Headers != nil {
				responseHeaders = upstreamErr.Headers.Clone()
		placeholder
			return &UpstreamFailoverError{
				StatusCode:      upstreamErr.StatusCode,
				ResponseBody:    upstreamErr.Body,
				ResponseHeaders: responseHeaders,
		placeholder
	placeholder
		msg := upstreamErr.Message
		if override := soraProErrorMessage(model, msg); override != "" {
			msg = override
	placeholder
		s.writeSoraError(c, upstreamErr.StatusCode, "upstream_error", msg, stream)
		return err
placeholder
	if errors.Is(err, context.DeadlineExceeded) {
		s.writeSoraError(c, http.StatusGatewayTimeout, "timeout_error", "Sora generation timeout", stream)
		return err
placeholder
	s.writeSoraError(c, http.StatusBadGateway, "api_error", err.Error(), stream)
	return err
placeholder

func (s *SoraGatewayService) pollImageTask(ctx context.Context, c *gin.Context, account *Account, taskID string, stream bool) ([]string, error) {
	interval := s.pollInterval()
	maxAttempts := s.pollMaxAttempts()
	lastPing := time.Now()
	for attempt := 0; attempt < maxAttempts; attempt++ {
		status, err := s.soraClient.GetImageTask(ctx, account, taskID)
		if err != nil {
			return nil, err
	placeholder
		switch strings.ToLower(status.Status) {
		case "succeeded", "completed":
			return status.URLs, nil
		case "failed":
			if status.ErrorMsg != "" {
				return nil, errors.New(status.ErrorMsg)
		placeholder
			return nil, errors.New("sora image generation failed")
	placeholder
		if stream {
			s.maybeSendPing(c, &lastPing)
	placeholder
		if err := sleepWithContext(ctx, interval); err != nil {
			return nil, err
	placeholder
placeholder
	return nil, errors.New("sora image generation timeout")
placeholder

func (s *SoraGatewayService) pollVideoTaskDetailed(ctx context.Context, c *gin.Context, account *Account, taskID string, stream bool) (*SoraVideoTaskStatus, error) {
	interval := s.pollInterval()
	maxAttempts := s.pollMaxAttempts()
	lastPing := time.Now()
	for attempt := 0; attempt < maxAttempts; attempt++ {
		status, err := s.soraClient.GetVideoTask(ctx, account, taskID)
		if err != nil {
			return nil, err
	placeholder
		switch strings.ToLower(status.Status) {
		case "completed", "succeeded":
			return status, nil
		case "failed":
			if status.ErrorMsg != "" {
				return nil, errors.New(status.ErrorMsg)
		placeholder
			return nil, errors.New("sora video generation failed")
	placeholder
		if stream {
			s.maybeSendPing(c, &lastPing)
	placeholder
		if err := sleepWithContext(ctx, interval); err != nil {
			return nil, err
	placeholder
placeholder
	return nil, errors.New("sora video generation timeout")
placeholder

func (s *SoraGatewayService) pollInterval() time.Duration {
	if s == nil || s.cfg == nil {
		return 2 * time.Second
placeholder
	interval := s.cfg.Sora.Client.PollIntervalSeconds
	if interval <= 0 {
		interval = 2
placeholder
	return time.Duration(interval) * time.Second
placeholder

func (s *SoraGatewayService) pollMaxAttempts() int {
	if s == nil || s.cfg == nil {
		return 600
placeholder
	maxAttempts := s.cfg.Sora.Client.MaxPollAttempts
	if maxAttempts <= 0 {
		maxAttempts = 600
placeholder
	return maxAttempts
placeholder

func (s *SoraGatewayService) maybeSendPing(c *gin.Context, lastPing *time.Time) {
	if c == nil {
		return
placeholder
	interval := 10 * time.Second
	if s != nil && s.cfg != nil && s.cfg.Concurrency.PingInterval > 0 {
		interval = time.Duration(s.cfg.Concurrency.PingInterval) * time.Second
placeholder
	if time.Since(*lastPing) < interval {
		return
placeholder
	if _, err := fmt.Fprint(c.Writer, ":\n\n"); err == nil {
		if flusher, ok := c.Writer.(http.Flusher); ok {
			flusher.Flush()
	placeholder
		*lastPing = time.Now()
placeholder
placeholder

func (s *SoraGatewayService) normalizeSoraMediaURLs(urls []string) []string {
	if len(urls) == 0 {
		return urls
placeholder
	output := make([]string, 0, len(urls))
	for _, raw := range urls {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
	placeholder
		if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
			output = append(output, raw)
			continue
	placeholder
		pathVal := raw
		if !strings.HasPrefix(pathVal, "/") {
			pathVal = "/" + pathVal
	placeholder
		output = append(output, s.buildSoraMediaURL(pathVal, ""))
placeholder
	return output
placeholder

func buildSoraContent(mediaType string, urls []string) string {
	switch mediaType {
	case "image":
		parts := make([]string, 0, len(urls))
		for _, u := range urls {
			parts = append(parts, fmt.Sprintf("![image](%s)", u))
	placeholder
		return strings.Join(parts, "\n")
	case "video":
		if len(urls) == 0 {
			return ""
	placeholder
		return fmt.Sprintf("```html\n<video src='%s' controls></video>\n```", urls[0])
	default:
		return ""
placeholder
placeholder

func extractSoraInput(body map[string]any) (prompt, imageInput, videoInput, remixTargetID string) {
	if body == nil {
		return "", "", "", ""
placeholder
	if v, ok := body["remix_target_id"].(string); ok {
		remixTargetID = strings.TrimSpace(v)
placeholder
	if v, ok := body["image"].(string); ok {
		imageInput = v
placeholder
	if v, ok := body["video"].(string); ok {
		videoInput = v
placeholder
	if v, ok := body["prompt"].(string); ok && strings.TrimSpace(v) != "" {
		prompt = v
placeholder
	if messages, ok := body["messages"].([]any); ok {
		builder := strings.Builder{placeholder
		for _, raw := range messages {
			msg, ok := raw.(map[string]any)
			if !ok {
				continue
		placeholder
			role, _ := msg["role"].(string)
			if role != "" && role != "user" {
				continue
		placeholder
			content := msg["content"]
			text, img, vid := parseSoraMessageContent(content)
			if text != "" {
				if builder.Len() > 0 {
					_, _ = builder.WriteString("\n")
			placeholder
				_, _ = builder.WriteString(text)
		placeholder
			if imageInput == "" && img != "" {
				imageInput = img
		placeholder
			if videoInput == "" && vid != "" {
				videoInput = vid
		placeholder
	placeholder
		if prompt == "" {
			prompt = builder.String()
	placeholder
placeholder
	if remixTargetID == "" {
		remixTargetID = extractRemixTargetIDFromPrompt(prompt)
placeholder
	prompt = cleanRemixLinkFromPrompt(prompt)
	return prompt, imageInput, videoInput, remixTargetID
placeholder

func parseSoraMessageContent(content any) (text, imageInput, videoInput string) {
	switch val := content.(type) {
	case string:
		return val, "", ""
	case []any:
		builder := strings.Builder{placeholder
		for _, item := range val {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
		placeholder
			t, _ := itemMap["type"].(string)
			switch t {
			case "text":
				if txt, ok := itemMap["text"].(string); ok && strings.TrimSpace(txt) != "" {
					if builder.Len() > 0 {
						_, _ = builder.WriteString("\n")
				placeholder
					_, _ = builder.WriteString(txt)
			placeholder
			case "image_url":
				if imageInput == "" {
					if urlVal, ok := itemMap["image_url"].(map[string]any); ok {
						imageInput = fmt.Sprintf("%v", urlVal["url"])
				placeholder else if urlStr, ok := itemMap["image_url"].(string); ok {
						imageInput = urlStr
				placeholder
			placeholder
			case "video_url":
				if videoInput == "" {
					if urlVal, ok := itemMap["video_url"].(map[string]any); ok {
						videoInput = fmt.Sprintf("%v", urlVal["url"])
				placeholder else if urlStr, ok := itemMap["video_url"].(string); ok {
						videoInput = urlStr
				placeholder
			placeholder
		placeholder
	placeholder
		return builder.String(), imageInput, videoInput
	default:
		return "", "", ""
placeholder
placeholder

func isSoraStoryboardPrompt(prompt string) bool {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return false
placeholder
	return len(soraStoryboardPattern.FindAllString(prompt, -1)) >= 1
placeholder

func formatSoraStoryboardPrompt(prompt string) string {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return ""
placeholder
	matches := soraStoryboardShotPattern.FindAllStringSubmatch(prompt, -1)
	if len(matches) == 0 {
		return prompt
placeholder
	firstBracketPos := strings.Index(prompt, "[")
	instructions := ""
	if firstBracketPos > 0 {
		instructions = strings.TrimSpace(prompt[:firstBracketPos])
placeholder
	shots := make([]string, 0, len(matches))
	for i, match := range matches {
		if len(match) < 3 {
			continue
	placeholder
		duration := strings.TrimSpace(match[1])
		scene := strings.TrimSpace(match[2])
		if scene == "" {
			continue
	placeholder
		shots = append(shots, fmt.Sprintf("Shot %d:\nduration: %ssec\nScene: %s", i+1, duration, scene))
placeholder
	if len(shots) == 0 {
		return prompt
placeholder
	timeline := strings.Join(shots, "\n\n")
	if instructions == "" {
		return timeline
placeholder
	return fmt.Sprintf("current timeline:\n%s\n\ninstructions:\n%s", timeline, instructions)
placeholder

func extractRemixTargetIDFromPrompt(prompt string) string {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return ""
placeholder
	return strings.TrimSpace(soraRemixTargetPattern.FindString(prompt))
placeholder

func cleanRemixLinkFromPrompt(prompt string) string {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return prompt
placeholder
	cleaned := soraRemixTargetInURLPattern.ReplaceAllString(prompt, "")
	cleaned = soraRemixTargetPattern.ReplaceAllString(cleaned, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return strings.TrimSpace(cleaned)
placeholder

func decodeSoraImageInput(ctx context.Context, input string) ([]byte, string, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, "", errors.New("empty image input")
placeholder
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, "", errors.New("invalid data url")
	placeholder
		meta := parts[0]
		payload := parts[1]
		decoded, err := decodeBase64WithLimit(payload, soraImageInputMaxBytes)
		if err != nil {
			return nil, "", err
	placeholder
		ext := ""
		if strings.HasPrefix(meta, "data:") {
			metaParts := strings.SplitN(meta[5:], ";", 2)
			if len(metaParts) > 0 {
				if exts, err := mime.ExtensionsByType(metaParts[0]); err == nil && len(exts) > 0 {
					ext = exts[0]
			placeholder
		placeholder
	placeholder
		filename := "image" + ext
		return decoded, filename, nil
placeholder
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return downloadSoraImageInput(ctx, raw)
placeholder
	decoded, err := decodeBase64WithLimit(raw, soraImageInputMaxBytes)
	if err != nil {
		return nil, "", errors.New("invalid base64 image")
placeholder
	return decoded, "image.png", nil
placeholder

func decodeSoraVideoInput(ctx context.Context, input string) ([]byte, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return nil, errors.New("empty video input")
placeholder
	if strings.HasPrefix(raw, "data:") {
		parts := strings.SplitN(raw, ",", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid video data url")
	placeholder
		decoded, err := decodeBase64WithLimit(parts[1], soraVideoInputMaxBytes)
		if err != nil {
			return nil, errors.New("invalid base64 video")
	placeholder
		if len(decoded) == 0 {
			return nil, errors.New("empty video data")
	placeholder
		return decoded, nil
placeholder
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return downloadSoraVideoInput(ctx, raw)
placeholder
	decoded, err := decodeBase64WithLimit(raw, soraVideoInputMaxBytes)
	if err != nil {
		return nil, errors.New("invalid base64 video")
placeholder
	if len(decoded) == 0 {
		return nil, errors.New("empty video data")
placeholder
	return decoded, nil
placeholder

func downloadSoraImageInput(ctx context.Context, rawURL string) ([]byte, string, error) {
	parsed, err := validateSoraRemoteURL(rawURL)
	if err != nil {
		return nil, "", err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, "", err
placeholder
	client := &http.Client{
		Timeout: soraImageInputTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= soraImageInputMaxRedirects {
				return errors.New("too many redirects")
		placeholder
			return validateSoraRemoteURLValue(req.URL)
	placeholder,
placeholder
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download image failed: %d", resp.StatusCode)
placeholder
	data, err := io.ReadAll(io.LimitReader(resp.Body, soraImageInputMaxBytes))
	if err != nil {
		return nil, "", err
placeholder
	ext := fileExtFromURL(parsed.String())
	if ext == "" {
		ext = fileExtFromContentType(resp.Header.Get("Content-Type"))
placeholder
	filename := "image" + ext
	return data, filename, nil
placeholder

func downloadSoraVideoInput(ctx context.Context, rawURL string) ([]byte, error) {
	parsed, err := validateSoraRemoteURL(rawURL)
	if err != nil {
		return nil, err
placeholder
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, err
placeholder
	client := &http.Client{
		Timeout: soraVideoInputTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= soraVideoInputMaxRedirects {
				return errors.New("too many redirects")
		placeholder
			return validateSoraRemoteURLValue(req.URL)
	placeholder,
placeholder
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = resp.Body.Close() placeholder()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download video failed: %d", resp.StatusCode)
placeholder
	data, err := io.ReadAll(io.LimitReader(resp.Body, soraVideoInputMaxBytes))
	if err != nil {
		return nil, err
placeholder
	if len(data) == 0 {
		return nil, errors.New("empty video content")
placeholder
	return data, nil
placeholder

func decodeBase64WithLimit(encoded string, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, errors.New("invalid max bytes limit")
placeholder
	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(encoded))
	limited := io.LimitReader(decoder, maxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
placeholder
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("input exceeds %d bytes limit", maxBytes)
placeholder
	return data, nil
placeholder

func validateSoraRemoteURL(raw string) (*url.URL, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, errors.New("empty remote url")
placeholder
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid remote url: %w", err)
placeholder
	if err := validateSoraRemoteURLValue(parsed); err != nil {
		return nil, err
placeholder
	return parsed, nil
placeholder

func validateSoraRemoteURLValue(parsed *url.URL) error {
	if parsed == nil {
		return errors.New("invalid remote url")
placeholder
	scheme := strings.ToLower(strings.TrimSpace(parsed.Scheme))
	if scheme != "http" && scheme != "https" {
		return errors.New("only http/https remote url is allowed")
placeholder
	if parsed.User != nil {
		return errors.New("remote url cannot contain userinfo")
placeholder
	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" {
		return errors.New("remote url missing host")
placeholder
	if _, blocked := soraBlockedHostnames[host]; blocked {
		return errors.New("remote url is not allowed")
placeholder
	if ip := net.ParseIP(host); ip != nil {
		if isSoraBlockedIP(ip) {
			return errors.New("remote url is not allowed")
	placeholder
		return nil
placeholder
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("resolve remote url failed: %w", err)
placeholder
	for _, ip := range ips {
		if isSoraBlockedIP(ip) {
			return errors.New("remote url is not allowed")
	placeholder
placeholder
	return nil
placeholder

func isSoraBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
placeholder
	for _, cidr := range soraBlockedCIDRs {
		if cidr.Contains(ip) {
			return true
	placeholder
placeholder
	return false
placeholder

func mustParseCIDRs(values []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(values))
	for _, val := range values {
		_, cidr, err := net.ParseCIDR(val)
		if err != nil {
			continue
	placeholder
		out = append(out, cidr)
placeholder
	return out
placeholder
