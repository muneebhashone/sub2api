package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	mathrand "math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	antigravityStickySessionTTL = time.Hour
	antigravityMaxRetries       = 3
	antigravityRetryBaseDelay   = 1 * time.Second
	antigravityRetryMaxDelay    = 16 * time.Second
)

// getSessionID 从 gin.Context 获取 session_id（用于日志追踪）
func getSessionID(c *gin.Context) string {
	if c == nil {
		return ""
placeholder
	return c.GetHeader("session_id")
placeholder

// logPrefix 生成统一的日志前缀
func logPrefix(sessionID, accountName string) string {
	if sessionID != "" {
		return fmt.Sprintf("[antigravity-Forward] session=%s account=%s", sessionID, accountName)
placeholder
	return fmt.Sprintf("[antigravity-Forward] account=%s", accountName)
placeholder

// Antigravity 直接支持的模型（精确匹配透传）
var antigravitySupportedModels = map[string]bool{
	"claude-opus-4-5-thinking":   true,
	"claude-sonnet-4-5":          true,
	"claude-sonnet-4-5-thinking": true,
	"gemini-2.5-flash":           true,
	"gemini-2.5-flash-lite":      true,
	"gemini-2.5-flash-thinking":  true,
	"gemini-3-flash":             true,
	"gemini-3-pro-low":           true,
	"gemini-3-pro-high":          true,
	"gemini-3-pro-image":         true,
placeholder

// Antigravity 前缀映射表（按前缀长度降序排列，确保最长匹配优先）
// 用于处理模型版本号变化（如 -20251111, -thinking, -preview 等后缀）
var antigravityPrefixMapping = []struct {
	prefix string
	target string
placeholder{
	// 长前缀优先
	{"gemini-2.5-flash-image", "gemini-3-pro-image"placeholder, // gemini-2.5-flash-image → 3-pro-image
	{"gemini-3-pro-image", "gemini-3-pro-image"placeholder,     // gemini-3-pro-image-preview 等
	{"claude-3-5-sonnet", "claude-sonnet-4-5"placeholder,       // 旧版 claude-3-5-sonnet-xxx
	{"claude-sonnet-4-5", "claude-sonnet-4-5"placeholder,       // claude-sonnet-4-5-xxx
	{"claude-haiku-4-5", "claude-sonnet-4-5"placeholder,        // claude-haiku-4-5-xxx → sonnet
	{"claude-opus-4-5", "claude-opus-4-5-thinking"placeholder,
	{"claude-3-haiku", "claude-sonnet-4-5"placeholder, // 旧版 claude-3-haiku-xxx → sonnet
	{"claude-sonnet-4", "claude-sonnet-4-5"placeholder,
	{"claude-haiku-4", "claude-sonnet-4-5"placeholder, // → sonnet
	{"claude-opus-4", "claude-opus-4-5-thinking"placeholder,
	{"gemini-3-pro", "gemini-3-pro-high"placeholder, // gemini-3-pro, gemini-3-pro-preview 等
placeholder

// AntigravityGatewayService 处理 Antigravity 平台的 API 转发
type AntigravityGatewayService struct {
	accountRepo      AccountRepository
	tokenProvider    *AntigravityTokenProvider
	rateLimitService *RateLimitService
	httpUpstream     HTTPUpstream
	settingService   *SettingService
placeholder

func NewAntigravityGatewayService(
	accountRepo AccountRepository,
	_ GatewayCache,
	tokenProvider *AntigravityTokenProvider,
	rateLimitService *RateLimitService,
	httpUpstream HTTPUpstream,
	settingService *SettingService,
) *AntigravityGatewayService {
	return &AntigravityGatewayService{
		accountRepo:      accountRepo,
		tokenProvider:    tokenProvider,
		rateLimitService: rateLimitService,
		httpUpstream:     httpUpstream,
		settingService:   settingService,
placeholder
placeholder

// GetTokenProvider 返回 token provider
func (s *AntigravityGatewayService) GetTokenProvider() *AntigravityTokenProvider {
	return s.tokenProvider
placeholder

// getMappedModel 获取映射后的模型名
// 逻辑：账户映射 → 直接支持透传 → 前缀映射 → gemini透传 → 默认值
func (s *AntigravityGatewayService) getMappedModel(account *Account, requestedModel string) string {
	// 1. 账户级映射（用户自定义优先）
	if mapped := account.GetMappedModel(requestedModel); mapped != requestedModel {
		return mapped
placeholder

	// 2. 直接支持的模型透传
	if antigravitySupportedModels[requestedModel] {
		return requestedModel
placeholder

	// 3. 前缀映射（处理版本号变化，如 -20251111, -thinking, -preview）
	for _, pm := range antigravityPrefixMapping {
		if strings.HasPrefix(requestedModel, pm.prefix) {
			return pm.target
	placeholder
placeholder

	// 4. Gemini 模型透传（未匹配到前缀的 gemini 模型）
	if strings.HasPrefix(requestedModel, "gemini-") {
		return requestedModel
placeholder

	// 5. 默认值
	return "claude-sonnet-4-5"
placeholder

// IsModelSupported 检查模型是否被支持
// 所有 claude- 和 gemini- 前缀的模型都能通过映射或透传支持
func (s *AntigravityGatewayService) IsModelSupported(requestedModel string) bool {
	return strings.HasPrefix(requestedModel, "claude-") ||
		strings.HasPrefix(requestedModel, "gemini-")
placeholder

// TestConnectionResult 测试连接结果
type TestConnectionResult struct {
	Text        string // 响应文本
	MappedModel string // 实际使用的模型
placeholder

// TestConnection 测试 Antigravity 账号连接（非流式，无重试、无计费）
// 支持 Claude 和 Gemini 两种协议，根据 modelID 前缀自动选择
func (s *AntigravityGatewayService) TestConnection(ctx context.Context, account *Account, modelID string) (*TestConnectionResult, error) {
	// 获取 token
	if s.tokenProvider == nil {
		return nil, errors.New("antigravity token provider not configured")
placeholder
	accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("获取 access_token 失败: %w", err)
placeholder

	// 获取 project_id（部分账户类型可能没有）
	projectID := strings.TrimSpace(account.GetCredential("project_id"))

	// 模型映射
	mappedModel := s.getMappedModel(account, modelID)

	// 构建请求体
	var requestBody []byte
	if strings.HasPrefix(modelID, "gemini-") {
		// Gemini 模型：直接使用 Gemini 格式
		requestBody, err = s.buildGeminiTestRequest(projectID, mappedModel)
placeholder else {
		// Claude 模型：使用协议转换
		requestBody, err = s.buildClaudeTestRequest(projectID, mappedModel)
placeholder
	if err != nil {
		return nil, fmt.Errorf("构建请求失败: %w", err)
placeholder

	// 构建 HTTP 请求（总是使用流式 endpoint，与官方客户端一致）
	req, err := antigravity.NewAPIRequest(ctx, "streamGenerateContent", accessToken, requestBody)
	if err != nil {
		return nil, err
placeholder

	// 调试日志：Test 请求信息
	log.Printf("[antigravity-Test] account=%s request_size=%d url=%s", account.Name, len(requestBody), req.URL.String())

	// 代理 URL
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	// 发送请求
	resp, err := s.httpUpstream.Do(req, proxyURL, account.ID, account.Concurrency)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	// 读取响应
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
placeholder

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API 返回 %d: %s", resp.StatusCode, string(respBody))
placeholder

	// 解析流式响应，提取文本
	text := extractTextFromSSEResponse(respBody)

	return &TestConnectionResult{
		Text:        text,
		MappedModel: mappedModel,
placeholder, nil
placeholder

// buildGeminiTestRequest 构建 Gemini 格式测试请求
func (s *AntigravityGatewayService) buildGeminiTestRequest(projectID, model string) ([]byte, error) {
	payload := map[string]any{
		"contents": []map[string]any{
			{
				"role": "user",
				"parts": []map[string]any{
					{"text": "hi"placeholder,
			placeholder,
		placeholder,
	placeholder,
		// Antigravity 上游要求必须包含身份提示词
		"systemInstruction": map[string]any{
			"parts": []map[string]any{
				{"text": antigravity.GetDefaultIdentityPatch()placeholder,
		placeholder,
	placeholder,
placeholder
	payloadBytes, _ := json.Marshal(payload)
	return s.wrapV1InternalRequest(projectID, model, payloadBytes)
placeholder

// buildClaudeTestRequest 构建 Claude 格式测试请求并转换为 Gemini 格式
func (s *AntigravityGatewayService) buildClaudeTestRequest(projectID, mappedModel string) ([]byte, error) {
	claudeReq := &antigravity.ClaudeRequest{
		Model: mappedModel,
		Messages: []antigravity.ClaudeMessage{
			{
				Role:    "user",
				Content: json.RawMessage(`"hi"`),
		placeholder,
	placeholder,
		MaxTokens: 1024,
		Stream:    false,
placeholder
	return antigravity.TransformClaudeToGemini(claudeReq, projectID, mappedModel)
placeholder

func (s *AntigravityGatewayService) getClaudeTransformOptions(ctx context.Context) antigravity.TransformOptions {
	opts := antigravity.DefaultTransformOptions()
	if s.settingService == nil {
		return opts
placeholder
	opts.EnableIdentityPatch = s.settingService.IsIdentityPatchEnabled(ctx)
	opts.IdentityPatch = s.settingService.GetIdentityPatchPrompt(ctx)
	return opts
placeholder

// extractTextFromSSEResponse 从 SSE 流式响应中提取文本
func extractTextFromSSEResponse(respBody []byte) string {
	var texts []string
	lines := bytes.Split(respBody, []byte("\n"))

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
	placeholder

		// 跳过 SSE 前缀
		if bytes.HasPrefix(line, []byte("data:")) {
			line = bytes.TrimPrefix(line, []byte("data:"))
			line = bytes.TrimSpace(line)
	placeholder

		// 跳过非 JSON 行
		if len(line) == 0 || line[0] != '{' {
			continue
	placeholder

		// 解析 JSON
		var data map[string]any
		if err := json.Unmarshal(line, &data); err != nil {
			continue
	placeholder

		// 尝试从 response.candidates[0].content.parts[].text 提取
		response, ok := data["response"].(map[string]any)
		if !ok {
			// 尝试直接从 candidates 提取（某些响应格式）
			response = data
	placeholder

		candidates, ok := response["candidates"].([]any)
		if !ok || len(candidates) == 0 {
			continue
	placeholder

		candidate, ok := candidates[0].(map[string]any)
		if !ok {
			continue
	placeholder

		content, ok := candidate["content"].(map[string]any)
		if !ok {
			continue
	placeholder

		parts, ok := content["parts"].([]any)
		if !ok {
			continue
	placeholder

		for _, part := range parts {
			if partMap, ok := part.(map[string]any); ok {
				if text, ok := partMap["text"].(string); ok && text != "" {
					texts = append(texts, text)
			placeholder
		placeholder
	placeholder
placeholder

	return strings.Join(texts, "")
placeholder

// injectIdentityPatchToGeminiRequest 为 Gemini 格式请求注入身份提示词
// 如果请求中已包含 "You are Antigravity" 则不重复注入
func injectIdentityPatchToGeminiRequest(body []byte) ([]byte, error) {
	var request map[string]any
	if err := json.Unmarshal(body, &request); err != nil {
		return nil, fmt.Errorf("解析 Gemini 请求失败: %w", err)
placeholder

	// 检查现有 systemInstruction 是否已包含身份提示词
	if sysInst, ok := request["systemInstruction"].(map[string]any); ok {
		if parts, ok := sysInst["parts"].([]any); ok {
			for _, part := range parts {
				if partMap, ok := part.(map[string]any); ok {
					if text, ok := partMap["text"].(string); ok {
						if strings.Contains(text, "You are Antigravity") {
							// 已包含身份提示词，直接返回原始请求
							return body, nil
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder
placeholder

	// 获取默认身份提示词
	identityPatch := antigravity.GetDefaultIdentityPatch()

	// 构建新的 systemInstruction
	newPart := map[string]any{"text": identityPatchplaceholder

	if existing, ok := request["systemInstruction"].(map[string]any); ok {
		// 已有 systemInstruction，在开头插入身份提示词
		if parts, ok := existing["parts"].([]any); ok {
			existing["parts"] = append([]any{newPartplaceholder, parts...)
	placeholder else {
			existing["parts"] = []any{newPartplaceholder
	placeholder
placeholder else {
		// 没有 systemInstruction，创建新的
		request["systemInstruction"] = map[string]any{
			"parts": []any{newPartplaceholder,
	placeholder
placeholder

	return json.Marshal(request)
placeholder

// wrapV1InternalRequest 包装请求为 v1internal 格式
func (s *AntigravityGatewayService) wrapV1InternalRequest(projectID, model string, originalBody []byte) ([]byte, error) {
	var request any
	if err := json.Unmarshal(originalBody, &request); err != nil {
		return nil, fmt.Errorf("解析请求体失败: %w", err)
placeholder

	wrapped := map[string]any{
		"project":     projectID,
		"requestId":   "agent-" + uuid.New().String(),
		"userAgent":   "antigravity", // 固定值，与官方客户端一致
		"requestType": "agent",
		"model":       model,
		"request":     request,
placeholder

	return json.Marshal(wrapped)
placeholder

// unwrapV1InternalResponse 解包 v1internal 响应
func (s *AntigravityGatewayService) unwrapV1InternalResponse(body []byte) ([]byte, error) {
	var outer map[string]any
	if err := json.Unmarshal(body, &outer); err != nil {
		return nil, err
placeholder

	if resp, ok := outer["response"]; ok {
		return json.Marshal(resp)
placeholder

	return body, nil
placeholder

// isModelNotFoundError 检测是否为模型不存在的 404 错误
func isModelNotFoundError(statusCode int, body []byte) bool {
	if statusCode != 404 {
		return false
placeholder

	bodyStr := strings.ToLower(string(body))
	keywords := []string{"model not found", "unknown model", "not found"placeholder
	for _, keyword := range keywords {
		if strings.Contains(bodyStr, keyword) {
			return true
	placeholder
placeholder
	return true // 404 without specific message also treated as model not found
placeholder

// Forward 转发 Claude 协议请求（Claude → Gemini 转换）
func (s *AntigravityGatewayService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte) (*ForwardResult, error) {
	startTime := time.Now()
	sessionID := getSessionID(c)
	prefix := logPrefix(sessionID, account.Name)

	// 解析 Claude 请求
	var claudeReq antigravity.ClaudeRequest
	if err := json.Unmarshal(body, &claudeReq); err != nil {
		return nil, fmt.Errorf("parse claude request: %w", err)
placeholder
	if strings.TrimSpace(claudeReq.Model) == "" {
		return nil, fmt.Errorf("missing model")
placeholder

	originalModel := claudeReq.Model
	mappedModel := s.getMappedModel(account, claudeReq.Model)

	// 获取 access_token
	if s.tokenProvider == nil {
		return nil, errors.New("antigravity token provider not configured")
placeholder
	accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("获取 access_token 失败: %w", err)
placeholder

	// 获取 project_id（部分账户类型可能没有）
	projectID := strings.TrimSpace(account.GetCredential("project_id"))

	// 代理 URL
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	// 获取转换选项
	// Antigravity 上游要求必须包含身份提示词，否则会返回 429
	transformOpts := s.getClaudeTransformOptions(ctx)
	transformOpts.EnableIdentityPatch = true // 强制启用，Antigravity 上游必需

	// 转换 Claude 请求为 Gemini 格式
	geminiBody, err := antigravity.TransformClaudeToGeminiWithOptions(&claudeReq, projectID, mappedModel, transformOpts)
	if err != nil {
		return nil, fmt.Errorf("transform request: %w", err)
placeholder

	// Antigravity 上游只支持流式请求，统一使用 streamGenerateContent
	// 如果客户端请求非流式，在响应处理阶段会收集完整流式响应后转换返回
	action := "streamGenerateContent"

	// 重试循环
	var resp *http.Response
	for attempt := 1; attempt <= antigravityMaxRetries; attempt++ {
		// 检查 context 是否已取消（客户端断开连接）
		select {
		case <-ctx.Done():
			log.Printf("%s status=context_canceled error=%v", prefix, ctx.Err())
			return nil, ctx.Err()
		default:
	placeholder

		upstreamReq, err := antigravity.NewAPIRequest(ctx, action, accessToken, geminiBody)
		if err != nil {
			return nil, err
	placeholder

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		if err != nil {
			if attempt < antigravityMaxRetries {
				log.Printf("%s status=request_failed retry=%d/%d error=%v", prefix, attempt, antigravityMaxRetries, err)
				if !sleepAntigravityBackoffWithContext(ctx, attempt) {
					log.Printf("%s status=context_canceled_during_backoff", prefix)
					return nil, ctx.Err()
			placeholder
				continue
		placeholder
			log.Printf("%s status=request_failed retries_exhausted error=%v", prefix, err)
			return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Upstream request failed after retries")
	placeholder

		if resp.StatusCode >= 400 && s.shouldRetryUpstreamError(resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()

			if attempt < antigravityMaxRetries {
				log.Printf("%s status=%d retry=%d/%d body=%s", prefix, resp.StatusCode, attempt, antigravityMaxRetries, truncateForLog(respBody, 500))
				if !sleepAntigravityBackoffWithContext(ctx, attempt) {
					log.Printf("%s status=context_canceled_during_backoff", prefix)
					return nil, ctx.Err()
			placeholder
				continue
		placeholder
			// 所有重试都失败，标记限流状态
			if resp.StatusCode == 429 {
				s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody)
		placeholder
			// 最后一次尝试也失败
			resp = &http.Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header.Clone(),
				Body:       io.NopCloser(bytes.NewReader(respBody)),
		placeholder
			break
	placeholder

		break
placeholder
	defer func() { _ = resp.Body.Close() placeholder()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))

		// 优先检测 thinking block 的 signature 相关错误（400）并重试一次：
		// Antigravity /v1internal 链路在部分场景会对 thought/thinking signature 做严格校验，
		// 当历史消息携带的 signature 不合法时会直接 400；去除 thinking 后可继续完成请求。
		if resp.StatusCode == http.StatusBadRequest && isSignatureRelatedError(respBody) {
			// Conservative two-stage fallback:
			// 1) Disable top-level thinking + thinking->text
			// 2) Only if still signature-related 400: also downgrade tool_use/tool_result to text.

			retryStages := []struct {
				name  string
				strip func(*antigravity.ClaudeRequest) (bool, error)
		placeholder{
				{name: "thinking-only", strip: stripThinkingFromClaudeRequestplaceholder,
				{name: "thinking+tools", strip: stripSignatureSensitiveBlocksFromClaudeRequestplaceholder,
		placeholder

			for _, stage := range retryStages {
				retryClaudeReq := claudeReq
				retryClaudeReq.Messages = append([]antigravity.ClaudeMessage(nil), claudeReq.Messages...)

				stripped, stripErr := stage.strip(&retryClaudeReq)
				if stripErr != nil || !stripped {
					continue
			placeholder

				log.Printf("Antigravity account %d: detected signature-related 400, retrying once (%s)", account.ID, stage.name)

				retryGeminiBody, txErr := antigravity.TransformClaudeToGeminiWithOptions(&retryClaudeReq, projectID, mappedModel, s.getClaudeTransformOptions(ctx))
				if txErr != nil {
					continue
			placeholder
				retryReq, buildErr := antigravity.NewAPIRequest(ctx, action, accessToken, retryGeminiBody)
				if buildErr != nil {
					continue
			placeholder
				retryResp, retryErr := s.httpUpstream.Do(retryReq, proxyURL, account.ID, account.Concurrency)
				if retryErr != nil {
					log.Printf("Antigravity account %d: signature retry request failed (%s): %v", account.ID, stage.name, retryErr)
					continue
			placeholder

				if retryResp.StatusCode < 400 {
					_ = resp.Body.Close()
					resp = retryResp
					respBody = nil
					break
			placeholder

				retryBody, _ := io.ReadAll(io.LimitReader(retryResp.Body, 2<<20))
				_ = retryResp.Body.Close()

				// If this stage fixed the signature issue, we stop; otherwise we may try the next stage.
				if retryResp.StatusCode != http.StatusBadRequest || !isSignatureRelatedError(retryBody) {
					respBody = retryBody
					resp = &http.Response{
						StatusCode: retryResp.StatusCode,
						Header:     retryResp.Header.Clone(),
						Body:       io.NopCloser(bytes.NewReader(retryBody)),
				placeholder
					break
			placeholder

				// Still signature-related; capture context and allow next stage.
				respBody = retryBody
				resp = &http.Response{
					StatusCode: retryResp.StatusCode,
					Header:     retryResp.Header.Clone(),
					Body:       io.NopCloser(bytes.NewReader(retryBody)),
			placeholder
		placeholder
	placeholder

		// 处理错误响应（重试后仍失败或不触发重试）
		if resp.StatusCode >= 400 {
			s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody)

			if s.shouldFailoverUpstreamError(resp.StatusCode) {
				return nil, &UpstreamFailoverError{StatusCode: resp.StatusCodeplaceholder
		placeholder

			return nil, s.writeMappedClaudeError(c, resp.StatusCode, respBody)
	placeholder
placeholder

	requestID := resp.Header.Get("x-request-id")
	if requestID != "" {
		c.Header("x-request-id", requestID)
placeholder

	var usage *ClaudeUsage
	var firstTokenMs *int
	if claudeReq.Stream {
		// 客户端要求流式，直接透传转换
		streamRes, err := s.handleClaudeStreamingResponse(c, resp, startTime, originalModel)
		if err != nil {
			log.Printf("%s status=stream_error error=%v", prefix, err)
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder else {
		// 客户端要求非流式，收集流式响应后转换返回
		streamRes, err := s.handleClaudeStreamToNonStreaming(c, resp, startTime, originalModel)
		if err != nil {
			log.Printf("%s status=stream_collect_error error=%v", prefix, err)
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel, // 使用原始模型用于计费和日志
		Stream:       claudeReq.Stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
placeholder, nil
placeholder

func isSignatureRelatedError(respBody []byte) bool {
	msg := strings.ToLower(strings.TrimSpace(extractAntigravityErrorMessage(respBody)))
	if msg == "" {
		// Fallback: best-effort scan of the raw payload.
		msg = strings.ToLower(string(respBody))
placeholder

	// Keep this intentionally broad: different upstreams may use "signature" or "thought_signature".
	if strings.Contains(msg, "thought_signature") || strings.Contains(msg, "signature") {
		return true
placeholder

	// Also detect thinking block structural errors:
	// "Expected `thinking` or `redacted_thinking`, but found `text`"
	if strings.Contains(msg, "expected") && (strings.Contains(msg, "thinking") || strings.Contains(msg, "redacted_thinking")) {
		return true
placeholder

	return false
placeholder

func extractAntigravityErrorMessage(body []byte) string {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
placeholder

	// Google-style: {"error": {"message": "..."placeholderplaceholder
	if errObj, ok := payload["error"].(map[string]any); ok {
		if msg, ok := errObj["message"].(string); ok && strings.TrimSpace(msg) != "" {
			return msg
	placeholder
placeholder

	// Fallback: top-level message
	if msg, ok := payload["message"].(string); ok && strings.TrimSpace(msg) != "" {
		return msg
placeholder

	return ""
placeholder

// stripThinkingFromClaudeRequest converts thinking blocks to text blocks in a Claude Messages request.
// This preserves the thinking content while avoiding signature validation errors.
// Note: redacted_thinking blocks are removed because they cannot be converted to text.
// It also disables top-level `thinking` to avoid upstream structural constraints for thinking mode.
func stripThinkingFromClaudeRequest(req *antigravity.ClaudeRequest) (bool, error) {
	if req == nil {
		return false, nil
placeholder

	changed := false
	if req.Thinking != nil {
		req.Thinking = nil
		changed = true
placeholder

	for i := range req.Messages {
		raw := req.Messages[i].Content
		if len(raw) == 0 {
			continue
	placeholder

		// If content is a string, nothing to strip.
		var str string
		if json.Unmarshal(raw, &str) == nil {
			continue
	placeholder

		// Otherwise treat as an array of blocks and convert thinking blocks to text.
		var blocks []map[string]any
		if err := json.Unmarshal(raw, &blocks); err != nil {
			continue
	placeholder

		filtered := make([]map[string]any, 0, len(blocks))
		modifiedAny := false
		for _, block := range blocks {
			t, _ := block["type"].(string)
			switch t {
			case "thinking":
				thinkingText, _ := block["thinking"].(string)
				if thinkingText != "" {
					filtered = append(filtered, map[string]any{
						"type": "text",
						"text": thinkingText,
				placeholder)
			placeholder
				modifiedAny = true
			case "redacted_thinking":
				modifiedAny = true
			case "":
				if thinkingText, hasThinking := block["thinking"].(string); hasThinking {
					if thinkingText != "" {
						filtered = append(filtered, map[string]any{
							"type": "text",
							"text": thinkingText,
					placeholder)
				placeholder
					modifiedAny = true
			placeholder else {
					filtered = append(filtered, block)
			placeholder
			default:
				filtered = append(filtered, block)
		placeholder
	placeholder

		if !modifiedAny {
			continue
	placeholder

		if len(filtered) == 0 {
			filtered = append(filtered, map[string]any{
				"type": "text",
				"text": "(content removed)",
		placeholder)
	placeholder

		newRaw, err := json.Marshal(filtered)
		if err != nil {
			return changed, err
	placeholder
		req.Messages[i].Content = newRaw
		changed = true
placeholder

	return changed, nil
placeholder

// stripSignatureSensitiveBlocksFromClaudeRequest is a stronger retry degradation that additionally converts
// tool blocks to plain text. Use this only after a thinking-only retry still fails with signature errors.
func stripSignatureSensitiveBlocksFromClaudeRequest(req *antigravity.ClaudeRequest) (bool, error) {
	if req == nil {
		return false, nil
placeholder

	changed := false
	if req.Thinking != nil {
		req.Thinking = nil
		changed = true
placeholder

	for i := range req.Messages {
		raw := req.Messages[i].Content
		if len(raw) == 0 {
			continue
	placeholder

		// If content is a string, nothing to strip.
		var str string
		if json.Unmarshal(raw, &str) == nil {
			continue
	placeholder

		// Otherwise treat as an array of blocks and convert signature-sensitive blocks to text.
		var blocks []map[string]any
		if err := json.Unmarshal(raw, &blocks); err != nil {
			continue
	placeholder

		filtered := make([]map[string]any, 0, len(blocks))
		modifiedAny := false
		for _, block := range blocks {
			t, _ := block["type"].(string)
			switch t {
			case "thinking":
				// Convert thinking to text, skip if empty
				thinkingText, _ := block["thinking"].(string)
				if thinkingText != "" {
					filtered = append(filtered, map[string]any{
						"type": "text",
						"text": thinkingText,
				placeholder)
			placeholder
				modifiedAny = true
			case "redacted_thinking":
				// Remove redacted_thinking (cannot convert encrypted content)
				modifiedAny = true
			case "tool_use":
				// Convert tool_use to text to avoid upstream signature/thought_signature validation errors.
				// This is a retry-only degradation path, so we prioritise request validity over tool semantics.
				name, _ := block["name"].(string)
				id, _ := block["id"].(string)
				input := block["input"]
				inputJSON, _ := json.Marshal(input)
				text := "(tool_use)"
				if name != "" {
					text += " name=" + name
			placeholder
				if id != "" {
					text += " id=" + id
			placeholder
				if len(inputJSON) > 0 && string(inputJSON) != "null" {
					text += " input=" + string(inputJSON)
			placeholder
				filtered = append(filtered, map[string]any{
					"type": "text",
					"text": text,
			placeholder)
				modifiedAny = true
			case "tool_result":
				// Convert tool_result to text so it stays consistent when tool_use is downgraded.
				toolUseID, _ := block["tool_use_id"].(string)
				isError, _ := block["is_error"].(bool)
				content := block["content"]
				contentJSON, _ := json.Marshal(content)
				text := "(tool_result)"
				if toolUseID != "" {
					text += " tool_use_id=" + toolUseID
			placeholder
				if isError {
					text += " is_error=true"
			placeholder
				if len(contentJSON) > 0 && string(contentJSON) != "null" {
					text += "\n" + string(contentJSON)
			placeholder
				filtered = append(filtered, map[string]any{
					"type": "text",
					"text": text,
			placeholder)
				modifiedAny = true
			case "":
				// Handle untyped block with "thinking" field
				if thinkingText, hasThinking := block["thinking"].(string); hasThinking {
					if thinkingText != "" {
						filtered = append(filtered, map[string]any{
							"type": "text",
							"text": thinkingText,
					placeholder)
				placeholder
					modifiedAny = true
			placeholder else {
					filtered = append(filtered, block)
			placeholder
			default:
				filtered = append(filtered, block)
		placeholder
	placeholder

		if !modifiedAny {
			continue
	placeholder

		if len(filtered) == 0 {
			// Keep request valid: upstream rejects empty content arrays.
			filtered = append(filtered, map[string]any{
				"type": "text",
				"text": "(content removed)",
		placeholder)
	placeholder

		newRaw, err := json.Marshal(filtered)
		if err != nil {
			return changed, err
	placeholder
		req.Messages[i].Content = newRaw
		changed = true
placeholder

	return changed, nil
placeholder

// ForwardGemini 转发 Gemini 协议请求
func (s *AntigravityGatewayService) ForwardGemini(ctx context.Context, c *gin.Context, account *Account, originalModel string, action string, stream bool, body []byte) (*ForwardResult, error) {
	startTime := time.Now()
	sessionID := getSessionID(c)
	prefix := logPrefix(sessionID, account.Name)

	if strings.TrimSpace(originalModel) == "" {
		return nil, s.writeGoogleError(c, http.StatusBadRequest, "Missing model in URL")
placeholder
	if strings.TrimSpace(action) == "" {
		return nil, s.writeGoogleError(c, http.StatusBadRequest, "Missing action in URL")
placeholder
	if len(body) == 0 {
		return nil, s.writeGoogleError(c, http.StatusBadRequest, "Request body is empty")
placeholder

	// 解析请求以获取 image_size（用于图片计费）
	imageSize := s.extractImageSize(body)

	switch action {
	case "generateContent", "streamGenerateContent":
		// ok
	case "countTokens":
		// 直接返回空值，不透传上游
		c.JSON(http.StatusOK, map[string]any{"totalTokens": 0placeholder)
		return &ForwardResult{
			RequestID:    "",
			Usage:        ClaudeUsage{placeholder,
			Model:        originalModel,
			Stream:       false,
			Duration:     time.Since(time.Now()),
			FirstTokenMs: nil,
	placeholder, nil
	default:
		return nil, s.writeGoogleError(c, http.StatusNotFound, "Unsupported action: "+action)
placeholder

	mappedModel := s.getMappedModel(account, originalModel)

	// 获取 access_token
	if s.tokenProvider == nil {
		return nil, errors.New("antigravity token provider not configured")
placeholder
	accessToken, err := s.tokenProvider.GetAccessToken(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("获取 access_token 失败: %w", err)
placeholder

	// 获取 project_id（部分账户类型可能没有）
	projectID := strings.TrimSpace(account.GetCredential("project_id"))

	// 代理 URL
	proxyURL := ""
	if account.ProxyID != nil && account.Proxy != nil {
		proxyURL = account.Proxy.URL()
placeholder

	// Antigravity 上游要求必须包含身份提示词，注入到请求中
	injectedBody, err := injectIdentityPatchToGeminiRequest(body)
	if err != nil {
		return nil, err
placeholder

	// 包装请求
	wrappedBody, err := s.wrapV1InternalRequest(projectID, mappedModel, injectedBody)
	if err != nil {
		return nil, err
placeholder

	// Antigravity 上游只支持流式请求，统一使用 streamGenerateContent
	// 如果客户端请求非流式，在响应处理阶段会收集完整流式响应后返回
	upstreamAction := "streamGenerateContent"

	// 重试循环
	var resp *http.Response
	for attempt := 1; attempt <= antigravityMaxRetries; attempt++ {
		// 检查 context 是否已取消（客户端断开连接）
		select {
		case <-ctx.Done():
			log.Printf("%s status=context_canceled error=%v", prefix, ctx.Err())
			return nil, ctx.Err()
		default:
	placeholder

		upstreamReq, err := antigravity.NewAPIRequest(ctx, upstreamAction, accessToken, wrappedBody)
		if err != nil {
			return nil, err
	placeholder

		resp, err = s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		if err != nil {
			if attempt < antigravityMaxRetries {
				log.Printf("%s status=request_failed retry=%d/%d error=%v", prefix, attempt, antigravityMaxRetries, err)
				if !sleepAntigravityBackoffWithContext(ctx, attempt) {
					log.Printf("%s status=context_canceled_during_backoff", prefix)
					return nil, ctx.Err()
			placeholder
				continue
		placeholder
			log.Printf("%s status=request_failed retries_exhausted error=%v", prefix, err)
			return nil, s.writeGoogleError(c, http.StatusBadGateway, "Upstream request failed after retries")
	placeholder

		if resp.StatusCode >= 400 && s.shouldRetryUpstreamError(resp.StatusCode) {
			respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
			_ = resp.Body.Close()

			if attempt < antigravityMaxRetries {
				log.Printf("%s status=%d retry=%d/%d", prefix, resp.StatusCode, attempt, antigravityMaxRetries)
				if !sleepAntigravityBackoffWithContext(ctx, attempt) {
					log.Printf("%s status=context_canceled_during_backoff", prefix)
					return nil, ctx.Err()
			placeholder
				continue
		placeholder
			// 所有重试都失败，标记限流状态
			if resp.StatusCode == 429 {
				s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody)
		placeholder
			resp = &http.Response{
				StatusCode: resp.StatusCode,
				Header:     resp.Header.Clone(),
				Body:       io.NopCloser(bytes.NewReader(respBody)),
		placeholder
			break
	placeholder

		break
placeholder
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
	placeholder
placeholder()

	// 处理错误响应
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
		// 尽早关闭原始响应体，释放连接；后续逻辑仍可能需要读取 body，因此用内存副本重新包装。
		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewReader(respBody))

		// 模型兜底：模型不存在且开启 fallback 时，自动用 fallback 模型重试一次
		if s.settingService != nil && s.settingService.IsModelFallbackEnabled(ctx) &&
			isModelNotFoundError(resp.StatusCode, respBody) {
			fallbackModel := s.settingService.GetFallbackModel(ctx, PlatformAntigravity)
			if fallbackModel != "" && fallbackModel != mappedModel {
				log.Printf("[Antigravity] Model not found (%s), retrying with fallback model %s (account: %s)", mappedModel, fallbackModel, account.Name)

				fallbackWrapped, err := s.wrapV1InternalRequest(projectID, fallbackModel, injectedBody)
				if err == nil {
					fallbackReq, err := antigravity.NewAPIRequest(ctx, upstreamAction, accessToken, fallbackWrapped)
					if err == nil {
						fallbackResp, err := s.httpUpstream.Do(fallbackReq, proxyURL, account.ID, account.Concurrency)
						if err == nil && fallbackResp.StatusCode < 400 {
							_ = resp.Body.Close()
							resp = fallbackResp
					placeholder else if fallbackResp != nil {
							_ = fallbackResp.Body.Close()
					placeholder
				placeholder
			placeholder
		placeholder
	placeholder

		// fallback 成功：继续按正常响应处理
		if resp.StatusCode < 400 {
			goto handleSuccess
	placeholder

		s.handleUpstreamError(ctx, prefix, account, resp.StatusCode, resp.Header, respBody)

		if s.shouldFailoverUpstreamError(resp.StatusCode) {
			return nil, &UpstreamFailoverError{StatusCode: resp.StatusCodeplaceholder
	placeholder

		// 解包并返回错误
		requestID := resp.Header.Get("x-request-id")
		if requestID != "" {
			c.Header("x-request-id", requestID)
	placeholder
		unwrapped, _ := s.unwrapV1InternalResponse(respBody)
		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
	placeholder
		c.Data(resp.StatusCode, contentType, unwrapped)
		return nil, fmt.Errorf("antigravity upstream error: %d", resp.StatusCode)
placeholder

handleSuccess:
	requestID := resp.Header.Get("x-request-id")
	if requestID != "" {
		c.Header("x-request-id", requestID)
placeholder

	var usage *ClaudeUsage
	var firstTokenMs *int

	if stream {
		// 客户端要求流式，直接透传
		streamRes, err := s.handleGeminiStreamingResponse(c, resp, startTime)
		if err != nil {
			log.Printf("%s status=stream_error error=%v", prefix, err)
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder else {
		// 客户端要求非流式，收集流式响应后返回
		streamRes, err := s.handleGeminiStreamToNonStreaming(c, resp, startTime)
		if err != nil {
			log.Printf("%s status=stream_collect_error error=%v", prefix, err)
			return nil, err
	placeholder
		usage = streamRes.usage
		firstTokenMs = streamRes.firstTokenMs
placeholder

	if usage == nil {
		usage = &ClaudeUsage{placeholder
placeholder

	// 判断是否为图片生成模型
	imageCount := 0
	if isImageGenerationModel(mappedModel) {
		// Gemini 图片生成 API 每次请求只生成一张图片（API 限制）
		imageCount = 1
placeholder

	return &ForwardResult{
		RequestID:    requestID,
		Usage:        *usage,
		Model:        originalModel,
		Stream:       stream,
		Duration:     time.Since(startTime),
		FirstTokenMs: firstTokenMs,
		ImageCount:   imageCount,
		ImageSize:    imageSize,
placeholder, nil
placeholder

func (s *AntigravityGatewayService) shouldRetryUpstreamError(statusCode int) bool {
	switch statusCode {
	case 429, 500, 502, 503, 504, 529:
		return true
	default:
		return false
placeholder
placeholder

func (s *AntigravityGatewayService) shouldFailoverUpstreamError(statusCode int) bool {
	switch statusCode {
	case 401, 403, 429, 529:
		return true
	default:
		return statusCode >= 500
placeholder
placeholder

// sleepAntigravityBackoffWithContext 带 context 取消检查的退避等待
// 返回 true 表示正常完成等待，false 表示 context 已取消
func sleepAntigravityBackoffWithContext(ctx context.Context, attempt int) bool {
	delay := antigravityRetryBaseDelay * time.Duration(1<<uint(attempt-1))
	if delay > antigravityRetryMaxDelay {
		delay = antigravityRetryMaxDelay
placeholder

	// +/- 20% jitter
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))
	jitter := time.Duration(float64(delay) * 0.2 * (r.Float64()*2 - 1))
	sleepFor := delay + jitter
	if sleepFor < 0 {
		sleepFor = 0
placeholder

	select {
	case <-ctx.Done():
		return false
	case <-time.After(sleepFor):
		return true
placeholder
placeholder

func (s *AntigravityGatewayService) handleUpstreamError(ctx context.Context, prefix string, account *Account, statusCode int, headers http.Header, body []byte) {
	// 429 使用 Gemini 格式解析（从 body 解析重置时间）
	if statusCode == 429 {
		resetAt := ParseGeminiRateLimitResetTime(body)
		if resetAt == nil {
			// 解析失败：Gemini 有重试时间用 5 分钟，Claude 没有用 1 分钟
			defaultDur := 1 * time.Minute
			if bytes.Contains(body, []byte("Please retry in")) || bytes.Contains(body, []byte("retryDelay")) {
				defaultDur = 5 * time.Minute
		placeholder
			ra := time.Now().Add(defaultDur)
			log.Printf("%s status=429 rate_limited reset_in=%v (fallback)", prefix, defaultDur)
			_ = s.accountRepo.SetRateLimited(ctx, account.ID, ra)
			return
	placeholder
		resetTime := time.Unix(*resetAt, 0)
		log.Printf("%s status=429 rate_limited reset_at=%v reset_in=%v", prefix, resetTime.Format("15:04:05"), time.Until(resetTime).Truncate(time.Second))
		_ = s.accountRepo.SetRateLimited(ctx, account.ID, resetTime)
		return
placeholder
	// 其他错误码继续使用 rateLimitService
	if s.rateLimitService == nil {
		return
placeholder
	shouldDisable := s.rateLimitService.HandleUpstreamError(ctx, account, statusCode, headers, body)
	if shouldDisable {
		log.Printf("%s status=%d marked_error", prefix, statusCode)
placeholder
placeholder

type antigravityStreamResult struct {
	usage        *ClaudeUsage
	firstTokenMs *int
placeholder

func (s *AntigravityGatewayService) handleGeminiStreamingResponse(c *gin.Context, resp *http.Response, startTime time.Time) (*antigravityStreamResult, error) {
	c.Status(resp.StatusCode)
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "text/event-stream; charset=utf-8"
placeholder
	c.Header("Content-Type", contentType)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder

	// 使用 Scanner 并限制单行大小，避免 ReadString 无上限导致 OOM
	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.settingService.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)
	usage := &ClaudeUsage{placeholder
	var firstTokenMs *int

	type scanEvent struct {
		line string
		err  error
placeholder
	// 独立 goroutine 读取上游，避免读取阻塞影响超时处理
	events := make(chan scanEvent, 16)
	done := make(chan struct{placeholder)
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
	placeholder
placeholder
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()placeholder) {
				return
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: errplaceholder)
	placeholder
placeholder()
	defer close(done)

	// 上游数据间隔超时保护（防止上游挂起长期占用连接）
	streamInterval := time.Duration(0)
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.settingService.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	// 仅发送一次错误事件，避免多次写入导致协议混乱
	errorEventSent := false
	sendErrorEvent := func(reason string) {
		if errorEventSent {
			return
	placeholder
		errorEventSent = true
		_, _ = fmt.Fprintf(c.Writer, "event: error\ndata: {\"error\":\"%s\"placeholder\n\n", reason)
		flusher.Flush()
placeholder

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, nil
		placeholder
			if ev.err != nil {
				if errors.Is(ev.err, bufio.ErrTooLong) {
					log.Printf("SSE line too long (antigravity): max_size=%d error=%v", maxLineSize, ev.err)
					sendErrorEvent("response_too_large")
					return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, ev.err
			placeholder
				sendErrorEvent("stream_read_error")
				return nil, ev.err
		placeholder

			line := ev.line
			trimmed := strings.TrimRight(line, "\r\n")
			if strings.HasPrefix(trimmed, "data:") {
				payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "data:"))
				if payload == "" || payload == "[DONE]" {
					if _, err := fmt.Fprintf(c.Writer, "%s\n", line); err != nil {
						sendErrorEvent("write_failed")
						return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, err
				placeholder
					flusher.Flush()
					continue
			placeholder

				// 解包 v1internal 响应
				inner, parseErr := s.unwrapV1InternalResponse([]byte(payload))
				if parseErr == nil && inner != nil {
					payload = string(inner)
			placeholder

				// 解析 usage
				var parsed map[string]any
				if json.Unmarshal(inner, &parsed) == nil {
					if u := extractGeminiUsage(parsed); u != nil {
						usage = u
				placeholder
			placeholder

				if firstTokenMs == nil {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
			placeholder

				if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", payload); err != nil {
					sendErrorEvent("write_failed")
					return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, err
			placeholder
				flusher.Flush()
				continue
		placeholder

			if _, err := fmt.Fprintf(c.Writer, "%s\n", line); err != nil {
				sendErrorEvent("write_failed")
				return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, err
		placeholder
			flusher.Flush()

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			log.Printf("Stream data interval timeout (antigravity)")
			sendErrorEvent("stream_timeout")
			return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, fmt.Errorf("stream data interval timeout")
	placeholder
placeholder
placeholder

// handleGeminiStreamToNonStreaming 读取上游流式响应，合并为非流式响应返回给客户端
// Gemini 流式响应中每个 chunk 都包含累积的完整文本，只需保留最后一个有效响应
func (s *AntigravityGatewayService) handleGeminiStreamToNonStreaming(c *gin.Context, resp *http.Response, startTime time.Time) (*antigravityStreamResult, error) {
	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.settingService.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	usage := &ClaudeUsage{placeholder
	var firstTokenMs *int
	var last map[string]any
	var lastWithParts map[string]any

	type scanEvent struct {
		line string
		err  error
placeholder

	// 独立 goroutine 读取上游，避免读取阻塞影响超时处理
	events := make(chan scanEvent, 16)
	done := make(chan struct{placeholder)
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
	placeholder
placeholder

	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()placeholder) {
				return
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: errplaceholder)
	placeholder
placeholder()
	defer close(done)

	// 上游数据间隔超时保护（防止上游挂起长期占用连接）
	streamInterval := time.Duration(0)
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.settingService.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				// 流结束，返回收集的响应
				goto returnResponse
		placeholder
			if ev.err != nil {
				if errors.Is(ev.err, bufio.ErrTooLong) {
					log.Printf("SSE line too long (antigravity non-stream): max_size=%d error=%v", maxLineSize, ev.err)
			placeholder
				return nil, ev.err
		placeholder

			line := ev.line
			trimmed := strings.TrimRight(line, "\r\n")

			if !strings.HasPrefix(trimmed, "data:") {
				continue
		placeholder

			payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "data:"))
			if payload == "" || payload == "[DONE]" {
				continue
		placeholder

			// 解包 v1internal 响应
			inner, parseErr := s.unwrapV1InternalResponse([]byte(payload))
			if parseErr != nil {
				continue
		placeholder

			var parsed map[string]any
			if err := json.Unmarshal(inner, &parsed); err != nil {
				continue
		placeholder

			// 记录首 token 时间
			if firstTokenMs == nil {
				ms := int(time.Since(startTime).Milliseconds())
				firstTokenMs = &ms
		placeholder

			last = parsed

			// 提取 usage
			if u := extractGeminiUsage(parsed); u != nil {
				usage = u
		placeholder

			// 保留最后一个有 parts 的响应
			if parts := extractGeminiParts(parsed); len(parts) > 0 {
				lastWithParts = parsed
		placeholder

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			log.Printf("Stream data interval timeout (antigravity non-stream)")
			return nil, fmt.Errorf("stream data interval timeout")
	placeholder
placeholder

returnResponse:
	// 选择最后一个有效响应
	finalResponse := pickGeminiCollectResult(last, lastWithParts)

	// 处理空响应情况
	if last == nil && lastWithParts == nil {
		log.Printf("[antigravity-Forward] warning: empty stream response, no valid chunks received")
placeholder

	respBody, err := json.Marshal(finalResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
placeholder
	c.Data(http.StatusOK, "application/json", respBody)

	return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, nil
placeholder

func (s *AntigravityGatewayService) writeClaudeError(c *gin.Context, status int, errType, message string) error {
	c.JSON(status, gin.H{
		"type":  "error",
		"error": gin.H{"type": errType, "message": messageplaceholder,
placeholder)
	return fmt.Errorf("%s", message)
placeholder

func (s *AntigravityGatewayService) writeMappedClaudeError(c *gin.Context, upstreamStatus int, body []byte) error {
	// 记录上游错误详情便于调试
	log.Printf("[antigravity-Forward] upstream_error status=%d body=%s", upstreamStatus, string(body))

	var statusCode int
	var errType, errMsg string

	switch upstreamStatus {
	case 400:
		statusCode = http.StatusBadRequest
		errType = "invalid_request_error"
		errMsg = "Invalid request"
	case 401:
		statusCode = http.StatusBadGateway
		errType = "authentication_error"
		errMsg = "Upstream authentication failed"
	case 403:
		statusCode = http.StatusBadGateway
		errType = "permission_error"
		errMsg = "Upstream access forbidden"
	case 429:
		statusCode = http.StatusTooManyRequests
		errType = "rate_limit_error"
		errMsg = "Upstream rate limit exceeded"
	case 529:
		statusCode = http.StatusServiceUnavailable
		errType = "overloaded_error"
		errMsg = "Upstream service overloaded"
	default:
		statusCode = http.StatusBadGateway
		errType = "upstream_error"
		errMsg = "Upstream request failed"
placeholder

	c.JSON(statusCode, gin.H{
		"type":  "error",
		"error": gin.H{"type": errType, "message": errMsgplaceholder,
placeholder)
	return fmt.Errorf("upstream error: %d", upstreamStatus)
placeholder

func (s *AntigravityGatewayService) writeGoogleError(c *gin.Context, status int, message string) error {
	statusStr := "UNKNOWN"
	switch status {
	case 400:
		statusStr = "INVALID_ARGUMENT"
	case 404:
		statusStr = "NOT_FOUND"
	case 429:
		statusStr = "RESOURCE_EXHAUSTED"
	case 500:
		statusStr = "INTERNAL"
	case 502, 503:
		statusStr = "UNAVAILABLE"
placeholder

	c.JSON(status, gin.H{
		"error": gin.H{
			"code":    status,
			"message": message,
			"status":  statusStr,
	placeholder,
placeholder)
	return fmt.Errorf("%s", message)
placeholder

// handleClaudeStreamToNonStreaming 收集上游流式响应，转换为 Claude 非流式格式返回
// 用于处理客户端非流式请求但上游只支持流式的情况
func (s *AntigravityGatewayService) handleClaudeStreamToNonStreaming(c *gin.Context, resp *http.Response, startTime time.Time, originalModel string) (*antigravityStreamResult, error) {
	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.settingService.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	var firstTokenMs *int
	var last map[string]any
	var lastWithParts map[string]any

	type scanEvent struct {
		line string
		err  error
placeholder

	// 独立 goroutine 读取上游，避免读取阻塞影响超时处理
	events := make(chan scanEvent, 16)
	done := make(chan struct{placeholder)
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
	placeholder
placeholder

	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()placeholder) {
				return
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: errplaceholder)
	placeholder
placeholder()
	defer close(done)

	// 上游数据间隔超时保护（防止上游挂起长期占用连接）
	streamInterval := time.Duration(0)
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.settingService.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				// 流结束，转换并返回响应
				goto returnResponse
		placeholder
			if ev.err != nil {
				if errors.Is(ev.err, bufio.ErrTooLong) {
					log.Printf("SSE line too long (antigravity claude non-stream): max_size=%d error=%v", maxLineSize, ev.err)
			placeholder
				return nil, ev.err
		placeholder

			line := ev.line
			trimmed := strings.TrimRight(line, "\r\n")

			if !strings.HasPrefix(trimmed, "data:") {
				continue
		placeholder

			payload := strings.TrimSpace(strings.TrimPrefix(trimmed, "data:"))
			if payload == "" || payload == "[DONE]" {
				continue
		placeholder

			// 解包 v1internal 响应
			inner, parseErr := s.unwrapV1InternalResponse([]byte(payload))
			if parseErr != nil {
				continue
		placeholder

			var parsed map[string]any
			if err := json.Unmarshal(inner, &parsed); err != nil {
				continue
		placeholder

			// 记录首 token 时间
			if firstTokenMs == nil {
				ms := int(time.Since(startTime).Milliseconds())
				firstTokenMs = &ms
		placeholder

			last = parsed

			// 保留最后一个有 parts 的响应
			if parts := extractGeminiParts(parsed); len(parts) > 0 {
				lastWithParts = parsed
		placeholder

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			log.Printf("Stream data interval timeout (antigravity claude non-stream)")
			return nil, fmt.Errorf("stream data interval timeout")
	placeholder
placeholder

returnResponse:
	// 选择最后一个有效响应
	finalResponse := pickGeminiCollectResult(last, lastWithParts)

	// 处理空响应情况
	if last == nil && lastWithParts == nil {
		log.Printf("[antigravity-Forward] warning: empty stream response, no valid chunks received")
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Empty response from upstream")
placeholder

	// 序列化为 JSON（Gemini 格式）
	geminiBody, err := json.Marshal(finalResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gemini response: %w", err)
placeholder

	// 转换 Gemini 响应为 Claude 格式
	claudeResp, agUsage, err := antigravity.TransformGeminiToClaude(geminiBody, originalModel)
	if err != nil {
		log.Printf("[antigravity-Forward] transform_error error=%v body=%s", err, string(geminiBody))
		return nil, s.writeClaudeError(c, http.StatusBadGateway, "upstream_error", "Failed to parse upstream response")
placeholder

	c.Data(http.StatusOK, "application/json", claudeResp)

	// 转换为 service.ClaudeUsage
	usage := &ClaudeUsage{
		InputTokens:              agUsage.InputTokens,
		OutputTokens:             agUsage.OutputTokens,
		CacheCreationInputTokens: agUsage.CacheCreationInputTokens,
		CacheReadInputTokens:     agUsage.CacheReadInputTokens,
placeholder

	return &antigravityStreamResult{usage: usage, firstTokenMs: firstTokenMsplaceholder, nil
placeholder

// handleClaudeStreamingResponse 处理 Claude 流式响应（Gemini SSE → Claude SSE 转换）
func (s *AntigravityGatewayService) handleClaudeStreamingResponse(c *gin.Context, resp *http.Response, startTime time.Time, originalModel string) (*antigravityStreamResult, error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming not supported")
placeholder

	processor := antigravity.NewStreamingProcessor(originalModel)
	var firstTokenMs *int
	// 使用 Scanner 并限制单行大小，避免 ReadString 无上限导致 OOM
	scanner := bufio.NewScanner(resp.Body)
	maxLineSize := defaultMaxLineSize
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.MaxLineSize > 0 {
		maxLineSize = s.settingService.cfg.Gateway.MaxLineSize
placeholder
	scanner.Buffer(make([]byte, 64*1024), maxLineSize)

	// 辅助函数：转换 antigravity.ClaudeUsage 到 service.ClaudeUsage
	convertUsage := func(agUsage *antigravity.ClaudeUsage) *ClaudeUsage {
		if agUsage == nil {
			return &ClaudeUsage{placeholder
	placeholder
		return &ClaudeUsage{
			InputTokens:              agUsage.InputTokens,
			OutputTokens:             agUsage.OutputTokens,
			CacheCreationInputTokens: agUsage.CacheCreationInputTokens,
			CacheReadInputTokens:     agUsage.CacheReadInputTokens,
	placeholder
placeholder

	type scanEvent struct {
		line string
		err  error
placeholder
	// 独立 goroutine 读取上游，避免读取阻塞影响超时处理
	events := make(chan scanEvent, 16)
	done := make(chan struct{placeholder)
	sendEvent := func(ev scanEvent) bool {
		select {
		case events <- ev:
			return true
		case <-done:
			return false
	placeholder
placeholder
	var lastReadAt int64
	atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
	go func() {
		defer close(events)
		for scanner.Scan() {
			atomic.StoreInt64(&lastReadAt, time.Now().UnixNano())
			if !sendEvent(scanEvent{line: scanner.Text()placeholder) {
				return
		placeholder
	placeholder
		if err := scanner.Err(); err != nil {
			_ = sendEvent(scanEvent{err: errplaceholder)
	placeholder
placeholder()
	defer close(done)

	streamInterval := time.Duration(0)
	if s.settingService.cfg != nil && s.settingService.cfg.Gateway.StreamDataIntervalTimeout > 0 {
		streamInterval = time.Duration(s.settingService.cfg.Gateway.StreamDataIntervalTimeout) * time.Second
placeholder
	var intervalTicker *time.Ticker
	if streamInterval > 0 {
		intervalTicker = time.NewTicker(streamInterval)
		defer intervalTicker.Stop()
placeholder
	var intervalCh <-chan time.Time
	if intervalTicker != nil {
		intervalCh = intervalTicker.C
placeholder

	// 仅发送一次错误事件，避免多次写入导致协议混乱
	errorEventSent := false
	sendErrorEvent := func(reason string) {
		if errorEventSent {
			return
	placeholder
		errorEventSent = true
		_, _ = fmt.Fprintf(c.Writer, "event: error\ndata: {\"error\":\"%s\"placeholder\n\n", reason)
		flusher.Flush()
placeholder

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				// 发送结束事件
				finalEvents, agUsage := processor.Finish()
				if len(finalEvents) > 0 {
					_, _ = c.Writer.Write(finalEvents)
					flusher.Flush()
			placeholder
				return &antigravityStreamResult{usage: convertUsage(agUsage), firstTokenMs: firstTokenMsplaceholder, nil
		placeholder
			if ev.err != nil {
				if errors.Is(ev.err, bufio.ErrTooLong) {
					log.Printf("SSE line too long (antigravity): max_size=%d error=%v", maxLineSize, ev.err)
					sendErrorEvent("response_too_large")
					return &antigravityStreamResult{usage: convertUsage(nil), firstTokenMs: firstTokenMsplaceholder, ev.err
			placeholder
				sendErrorEvent("stream_read_error")
				return nil, fmt.Errorf("stream read error: %w", ev.err)
		placeholder

			line := ev.line
			// 处理 SSE 行，转换为 Claude 格式
			claudeEvents := processor.ProcessLine(strings.TrimRight(line, "\r\n"))

			if len(claudeEvents) > 0 {
				if firstTokenMs == nil {
					ms := int(time.Since(startTime).Milliseconds())
					firstTokenMs = &ms
			placeholder

				if _, writeErr := c.Writer.Write(claudeEvents); writeErr != nil {
					finalEvents, agUsage := processor.Finish()
					if len(finalEvents) > 0 {
						_, _ = c.Writer.Write(finalEvents)
				placeholder
					sendErrorEvent("write_failed")
					return &antigravityStreamResult{usage: convertUsage(agUsage), firstTokenMs: firstTokenMsplaceholder, writeErr
			placeholder
				flusher.Flush()
		placeholder

		case <-intervalCh:
			lastRead := time.Unix(0, atomic.LoadInt64(&lastReadAt))
			if time.Since(lastRead) < streamInterval {
				continue
		placeholder
			log.Printf("Stream data interval timeout (antigravity)")
			sendErrorEvent("stream_timeout")
			return &antigravityStreamResult{usage: convertUsage(nil), firstTokenMs: firstTokenMsplaceholder, fmt.Errorf("stream data interval timeout")
	placeholder
placeholder

placeholder

// extractImageSize 从 Gemini 请求中提取 image_size 参数
func (s *AntigravityGatewayService) extractImageSize(body []byte) string {
	var req antigravity.GeminiRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return "2K" // 默认 2K
placeholder

	if req.GenerationConfig != nil && req.GenerationConfig.ImageConfig != nil {
		size := strings.ToUpper(strings.TrimSpace(req.GenerationConfig.ImageConfig.ImageSize))
		if size == "1K" || size == "2K" || size == "4K" {
			return size
	placeholder
placeholder

	return "2K" // 默认 2K
placeholder

// isImageGenerationModel 判断模型是否为图片生成模型
// 支持的模型：gemini-3-pro-image, gemini-3-pro-image-preview, gemini-2.5-flash-image 等
func isImageGenerationModel(model string) bool {
	modelLower := strings.ToLower(model)
	// 移除 models/ 前缀
	modelLower = strings.TrimPrefix(modelLower, "models/")

	// 精确匹配或前缀匹配
	return modelLower == "gemini-3-pro-image" ||
		modelLower == "gemini-3-pro-image-preview" ||
		strings.HasPrefix(modelLower, "gemini-3-pro-image-") ||
		modelLower == "gemini-2.5-flash-image" ||
		modelLower == "gemini-2.5-flash-image-preview" ||
		strings.HasPrefix(modelLower, "gemini-2.5-flash-image-")
placeholder
