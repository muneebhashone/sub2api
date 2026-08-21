package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai_compat"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// Forward forwards request to OpenAI API
func (s *OpenAIGatewayService) Forward(ctx context.Context, c *gin.Context, account *Account, body []byte) (*OpenAIForwardResult, error) {
	beginUpstreamResponseModelObservation(c)
	clearGrokResponsesClientToolMapping(c)
	clearOpenAIResponsesClientToolMapping(c)
	clearOpenAIResponsesNamespaceNames(c)
	setCodexToolNameReverse(c, nil)
	startTime := time.Now()
	// 固定渠道映射后的请求级 canonical body；账号 normalize/strip 不得改写跨 failover hint。
	canonicalImageIntentBody := body

	restrictionResult := s.detectCodexClientRestriction(c, account, body)
	apiKeyID := getAPIKeyIDFromContext(c)
	logCodexCLIOnlyDetection(ctx, c, account, apiKeyID, restrictionResult, body)
	if restrictionResult.Enabled && !restrictionResult.Matched {
		MarkOpsClientBusinessLimited(c, OpsClientBusinessLimitedReasonLocalPolicyDenied)
		c.JSON(http.StatusForbidden, gin.H{
			"error": gin.H{
				"type":    "forbidden_error",
				"message": CodexClientRestrictionMessage(restrictionResult),
		placeholder,
	placeholder)
		return nil, errors.New("codex_cli_only restriction: only codex official clients are allowed")
placeholder

	normalizedBody, normalized, err := normalizeOpenAICodexCompactReasoningEffortForAccount(c, account, body)
	if err != nil {
		return nil, err
placeholder
	if normalized {
		body = normalizedBody
placeholder
	legacyIngressBody, legacyIngressChanged, legacyIngressErr := normalizeOpenAIResponsesLegacyIngress(body)
	if legacyIngressErr != nil {
		return nil, legacyIngressErr
placeholder
	if legacyIngressChanged {
		body = legacyIngressBody
placeholder
	// 在分流到 passthrough / Codex transform / 原生 ChatCompletions 之前统一修正
	// 显式为 null 的工具 Schema type，否则 upstream 的 400 会被归一成可重试的 502，
	// 同一份坏定义在账号池里反复重放。
	if sanitizedToolBody, toolSchemaSanitized, toolSchemaErr := sanitizeOpenAIResponsesToolSchemasForPlatform(body, account.Platform); toolSchemaErr != nil {
		return nil, toolSchemaErr
placeholder else if toolSchemaSanitized {
		body = sanitizedToolBody
placeholder
	if account.IsOpenAIOAuthLike() {
		reasoningBody, reasoningChanged, reasoningErr := normalizeOpenAIResponsesReasoningMode(body)
		if reasoningErr != nil {
			return nil, fmt.Errorf("normalize OpenAI Responses reasoning.mode: %w", reasoningErr)
	placeholder
		if reasoningChanged {
			body = reasoningBody
	placeholder
placeholder
	if account.IsOpenAIOAuthLike() && isOpenAIResponsesLiteHeader(c.GetHeader(responsesLiteHeader)) {
		liteBody, changed, liteErr := normalizeOpenAIResponsesLiteToolsPayload(body)
		if liteErr != nil {
			setOpsUpstreamError(c, http.StatusBadRequest, liteErr.Error(), "")
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
				"type": "invalid_request_error", "message": liteErr.Error(), "param": "tools",
		placeholderplaceholder)
			return nil, liteErr
	placeholder
		if changed {
			body = liteBody
	placeholder
placeholder
	wsDecision := s.getOpenAIWSProtocolResolver().Resolve(account)
	// 仅允许 WS 入站请求走 WS 上游，避免出现 HTTP -> WS 协议混用。
	wsDecision = resolveOpenAIWSDecisionByClientTransport(wsDecision, GetOpenAIClientTransport(c))
	passthroughEnabled := account.IsOpenAIPassthroughEnabled()
	compactPath := isOpenAIResponsesCompactPath(c)
	if shouldFlattenOpenAIResponsesNamespaces(account, wsDecision.Transport, passthroughEnabled, compactPath) {
		body, err = flattenOpenAIResponsesNamespaces(c, body)
		if err != nil {
			setOpsUpstreamError(c, http.StatusBadRequest, err.Error(), "")
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
				"type": "invalid_request_error", "message": err.Error(), "param": "tools",
		placeholderplaceholder)
			return nil, err
	placeholder
placeholder
	if shouldStripOpenAIResponsesInputNamespaces(account, wsDecision.Transport, passthroughEnabled) {
		keepToolCallNamespaces := shouldKeepOpenAIResponsesToolCallNamespaces(
			account, wsDecision.Transport, passthroughEnabled, compactPath,
		)
		body, err = stripOpenAIResponsesInputNamespaces(body, keepToolCallNamespaces)
		if err != nil {
			setOpsUpstreamError(c, http.StatusBadRequest, err.Error(), "")
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
				"type": "invalid_request_error", "message": err.Error(), "param": "input",
		placeholderplaceholder)
			return nil, err
	placeholder
placeholder

	nativeDeepSeekResponses := account.Platform == PlatformDeepseek &&
		(account.GetAPIProtocol() == APIProtocolResponses || account.IsAdaptiveAPIProtocol())
	if nativeDeepSeekResponses && account.Type == AccountTypeAPIKey && !compactPath &&
		needsOpenAIResponsesClientToolAdaptation(body) {
		adaptedBody, mapping, adaptErr := adaptOpenAIResponsesClientTools(body)
		if adaptErr != nil {
			return nil, fmt.Errorf("adapt DeepSeek Responses client tools: %w", adaptErr)
	placeholder
		body = adaptedBody
		setOpenAIResponsesClientToolMapping(c, mapping)
placeholder

	originalBody := body
	requestView := newOpenAIRequestView(body)
	reqModel, reqStream, promptCacheKey := requestView.Model, requestView.Stream, requestView.PromptCacheKey
	originalModel := reqModel

	if account.Platform == PlatformGrok {
		return s.forwardGrokResponses(ctx, c, account, body, originalModel, reqStream, startTime)
placeholder

	// CN 供应商 anthropic 协议账号：/v1/responses 入站是交叉协议组合
	// （Responses 客户端 × Anthropic 上游），转成 Anthropic 请求走原生端点。
	// 不能落到下面的 raw-CC 分支——其 URL 构造会把 anthropic base 当 CC base 用。
	if account.IsAnthropicProtocol() {
		return s.forwardResponsesViaNativeAnthropic(ctx, c, account, body, reqModel)
placeholder
	if account.IsOpenAIApiKey() {
		if normalized, changed, normalizeErr := normalizeOpenAIParallelToolCallsWithoutTools(body); normalizeErr != nil {
			return nil, normalizeErr
	placeholder else if changed {
			body = normalized
			originalBody = normalized
	placeholder
		if normalized, changed, normalizeErr := normalizeOpenAIAPIKeyStoreFalseReasoningReplay(body, isOpenAIResponsesCompactPath(c)); normalizeErr != nil {
			return nil, normalizeErr
	placeholder else if changed {
			body = normalized
			originalBody = normalized
	placeholder
		requestView = newOpenAIRequestView(body)
		reqModel, reqStream, promptCacheKey = requestView.Model, requestView.Stream, requestView.PromptCacheKey
		originalModel = reqModel
placeholder

	if shouldForwardOpenAIResponsesViaRawChatCompletions(account) {
		return s.forwardResponsesViaRawChatCompletions(ctx, c, account, body)
placeholder
	if account.IsOpenAI() && (account.IsOpenAIApiKey() || account.IsOpenAIOAuthLike()) {
		sanitizedBody, changed, sanitizeErr := sanitizeOpenAIResponsesInputItemIDs(body)
		if sanitizeErr != nil {
			return nil, fmt.Errorf("sanitize OpenAI Responses input item IDs: %w", sanitizeErr)
	placeholder
		if changed {
			body = sanitizedBody
			originalBody = sanitizedBody
			requestView = newOpenAIRequestView(sanitizedBody)
			reqModel, reqStream, promptCacheKey = requestView.Model, requestView.Stream, requestView.PromptCacheKey
			originalModel = reqModel
	placeholder
placeholder

	compatMessagesBridge := isOpenAICompatMessagesBridgeBody(body)
	setOpenAICompatMessagesBridgeContext(c, compatMessagesBridge)

	isCodexCLI := openai.IsCodexOfficialClientByHeaders(c.GetHeader("User-Agent"), c.GetHeader("originator")) || (s.cfg != nil && s.cfg.Gateway.ForceCodexCLI)
	codexImageGenerationExplicitToolPolicy := codexImageGenerationExplicitToolPolicyAllow
	if isCodexCLI {
		codexImageGenerationExplicitToolPolicy = account.CodexImageGenerationExplicitToolPolicy()
placeholder
	if c != nil {
		c.Set("openai_ws_transport_decision", string(wsDecision.Transport))
		c.Set("openai_ws_transport_reason", wsDecision.Reason)
placeholder
	if wsDecision.Transport == OpenAIUpstreamTransportResponsesWebsocketV2 {
		logOpenAIWSModeDebug(
			"selected account_id=%d account_type=%s transport=%s reason=%s model=%s stream=%v",
			account.ID,
			account.Type,
			normalizeOpenAIWSLogValue(string(wsDecision.Transport)),
			normalizeOpenAIWSLogValue(wsDecision.Reason),
			reqModel,
			reqStream,
		)
placeholder
	// 当前仅支持 WSv2；WSv1 命中时直接返回错误，避免出现“配置可开但行为不确定”。
	if wsDecision.Transport == OpenAIUpstreamTransportResponsesWebsocket {
		if c != nil {
			MarkOpsClientBusinessLimited(c, OpsClientBusinessLimitedReasonLocalFeatureGate)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"type":    "invalid_request_error",
					"message": "OpenAI WSv1 is temporarily unsupported. Please enable responses_websockets_v2.",
			placeholder,
		placeholder)
	placeholder
		return nil, errors.New("openai ws v1 is temporarily unsupported; use ws v2")
placeholder
	if passthroughEnabled {
		attemptImageIntentInvalidated := false
		if isCodexCLI && codexImageGenerationExplicitToolPolicy == codexImageGenerationExplicitToolPolicyStrip {
			strippedBody, changed, stripErr := stripOpenAIImageGenerationToolsFromRawPayload(body)
			if stripErr != nil {
				return nil, stripErr
		placeholder
			if changed {
				body = strippedBody
				originalBody = strippedBody
				attemptImageIntentInvalidated = true
				logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Stripped /responses image_generation tool for Codex client by account policy")
		placeholder
	placeholder
		// 透传分支只需要轻量提取字段，避免热路径全量 Unmarshal。
		mappedModel := account.GetMappedModel(reqModel)
		reasoningEffort := extractOpenAIReasoningEffortFromBody(body, mappedModel)
		// 国产模型默认 effort 补充：也要用 mappedModel 判定是否是 passback-required 上游。
		reasoningEffort = ApplyThinkingEnabledFallback(reasoningEffort, body, mappedModel)
		return s.forwardOpenAIPassthrough(
			ctx,
			c,
			account,
			originalBody,
			canonicalImageIntentBody,
			reqModel,
			attemptImageIntentInvalidated,
			reasoningEffort,
			reqStream,
			startTime,
		)
placeholder

	bodyModified := false
	var reqBody map[string]any
	ensureReqBody := func() (map[string]any, error) {
		if requestView.HasPatches() {
			patchedBody, patchErr := requestView.ApplyPatches()
			if patchErr != nil {
				return nil, patchErr
		placeholder
			body = patchedBody
			requestView = newOpenAIRequestView(body)
			reqBody = nil
			bodyModified = false
	placeholder
		if reqBody != nil {
			return reqBody, nil
	placeholder
		decoded, decodeErr := requestView.Decode(c)
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		reqBody = decoded
		return reqBody, nil
placeholder
	markPatchSet := func(path string, value any) {
		bodyModified = true
		if requestView.patchesDisabled {
			if reqBody != nil {
				setOpenAIRequestMapPath(reqBody, path, value)
		placeholder
			return
	placeholder
		requestView.MarkPatchSet(path, value)
placeholder
	markPatchDelete := func(path string) {
		bodyModified = true
		if requestView.patchesDisabled {
			if reqBody != nil {
				deleteOpenAIRequestMapPath(reqBody, path)
		placeholder
			return
	placeholder
		requestView.MarkPatchDelete(path)
placeholder
	disablePatch := func() {
		requestView.DisablePatches()
placeholder
	markDecodedModified := func() {
		bodyModified = true
		disablePatch()
placeholder

	apiKey := getAPIKeyFromContext(c)
	imageGenerationAllowed := GroupAllowsImageGeneration(nil)
	if apiKey != nil {
		imageGenerationAllowed = GroupAllowsImageGeneration(apiKey.Group)
placeholder
	codexImageGenerationBridgeEnabled := isCodexCLI &&
		!isOpenAIResponsesLiteHeader(c.GetHeader(responsesLiteHeader)) &&
		imageGenerationAllowed &&
		codexImageGenerationExplicitToolPolicy != codexImageGenerationExplicitToolPolicyStrip &&
		s.isCodexImageGenerationBridgeEnabled(ctx, account, apiKey)
	var imageIntent bool
	canonicalImageIntent := resolveOpenAIImageIntentHint(c, reqModel, canonicalImageIntentBody, IsImageGenerationIntent)
	if isCodexCLI && codexImageGenerationExplicitToolPolicy == codexImageGenerationExplicitToolPolicyStrip {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if stripOpenAIImageGenerationTools(decoded) {
			markDecodedModified()
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Stripped /responses image_generation tool for Codex client by account policy")
	placeholder
		imageIntent = IsImageGenerationIntentMap(openAIResponsesEndpoint, reqModel, decoded)
placeholder else {
		imageIntent = canonicalImageIntent
placeholder
	if imageIntent && !imageGenerationAllowed {
		MarkOpsClientBusinessLimited(c, OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"type": "permission_error", "message": ImageGenerationPermissionMessage()placeholderplaceholder)
		return nil, errors.New("image generation disabled for group")
placeholder

	instructions := gjson.GetBytes(body, "instructions")
	instructionsEmpty := !instructions.Exists() || instructions.Type != gjson.String || strings.TrimSpace(instructions.String()) == ""
	if instructionsEmpty && !compatMessagesBridge && !nativeDeepSeekResponses {
		markPatchSet("instructions", defaultCodexSynthInstructions(reqModel))
placeholder

	isCompactRequest := compactPath
	requestedModel := reqModel
	billingModel, upstreamModel := resolveOpenAIForwardMappedModels(account, requestedModel, isCompactRequest)
	if isCompactRequest {
		if compactModel := s.resolveOpenAICompactFallbackModel(account, requestedModel); compactModel != "" {
			upstreamModel = compactModel
	placeholder
placeholder
	if billingModel != requestedModel {
		logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Model mapping applied: %s -> %s (account: %s, isCodexCLI: %v)", requestedModel, billingModel, account.Name, isCodexCLI)
placeholder
	reqModel = billingModel
	if upstreamModel != requestedModel {
		markPatchSet("model", upstreamModel)
placeholder
	if upstreamModel != billingModel {
		if isCompactRequest {
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Compact model mapping applied: %s -> %s (account: %s, isCodexCLI: %v)", requestedModel, upstreamModel, account.Name, isCodexCLI)
	placeholder else {
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Upstream model resolved: %s -> %s (account: %s, type: %s, isCodexCLI: %v)", billingModel, upstreamModel, account.Name, account.Type, isCodexCLI)
	placeholder
placeholder
	if strings.TrimSpace(gjson.GetBytes(body, "reasoning.effort").String()) == "minimal" {
		markPatchSet("reasoning.effort", "none")
		logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Normalized reasoning.effort: minimal -> none (account: %s)", account.Name)
placeholder
	if strings.TrimSpace(gjson.GetBytes(body, "text.format.type").String()) == "json_schema" ||
		strings.TrimSpace(gjson.GetBytes(body, "response_format.type").String()) == "json_schema" {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if normalizeOpenAIResponseFormatSchemas(decoded) {
			markDecodedModified()
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Normalized Responses JSON schema compatibility")
	placeholder
placeholder

	imageIntent = imageIntent || IsImageGenerationIntent(openAIResponsesEndpoint, reqModel, nil) || isOpenAIImageGenerationModel(upstreamModel)
	if imageIntent && !imageGenerationAllowed {
		MarkOpsClientBusinessLimited(c, OpsClientBusinessLimitedReasonLocalFeatureGate)
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"type": "permission_error", "message": ImageGenerationPermissionMessage()placeholderplaceholder)
		return nil, errors.New("image generation disabled for group")
placeholder

	// /responses/compact 是会话压缩请求：上游不接受 tool_choice（400 unknown_parameter），
	// 注入 image_generation 工具也没有意义，整块豁免。
	if imageGenerationAllowed && !isCompactRequest && (codexImageGenerationBridgeEnabled || isOpenAIImageGenerationModel(requestView.Model) || openAIRequestBodyImageGenerationToolNeedsNormalization(body) || isOpenAIImageGenerationModel(upstreamModel)) {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if codexImageGenerationBridgeEnabled && ensureOpenAIResponsesImageGenerationTool(decoded) {
			markDecodedModified()
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Injected /responses image_generation tool for Codex client")
	placeholder
		if codexImageGenerationBridgeEnabled && ensureOpenAIResponsesImageGenerationToolChoiceAuto(decoded) {
			markDecodedModified()
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Set /responses image_generation tool_choice=auto for Codex client")
	placeholder
		if normalizeOpenAIResponsesImageGenerationTools(decoded) {
			markDecodedModified()
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Normalized /responses image_generation tool payload")
	placeholder
		if normalizeOpenAIResponsesImageOnlyModel(decoded) {
			markDecodedModified()
			if model, ok := decoded["model"].(string); ok {
				upstreamModel = strings.TrimSpace(model)
		placeholder
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Normalized /responses image-only model request inbound_model=%s image_model=%s upstream_model=%s", requestView.Model, billingModel, upstreamModel)
	placeholder
		if err := validateOpenAIResponsesImageModel(decoded, upstreamModel); err != nil {
			setOpsUpstreamError(c, http.StatusBadRequest, err.Error(), "")
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"type": "invalid_request_error", "message": err.Error(), "param": "model"placeholderplaceholder)
			return nil, err
	placeholder
		if hasOpenAIImageGenerationTool(decoded) {
			imageIntent = true
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] /responses image_generation request inbound_model=%s mapped_model=%s account_type=%s", requestView.Model, upstreamModel, account.Type)
	placeholder
		if codexImageGenerationBridgeEnabled && applyCodexImageGenerationBridgeInstructions(decoded) {
			markDecodedModified()
			logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Added Codex image_generation bridge instructions")
	placeholder
placeholder else if imageGenerationAllowed && imageIntent && openAIRequestBodyHasImageGenerationDeclaration(body) {
		// 完整 image_generation tool 只做 raw 计费读取，校验/桥接/旧字段迁移命中时才展开大 input map。
		logger.LegacyPrintf("service.openai_gateway", "[OpenAI] /responses image_generation request inbound_model=%s mapped_model=%s account_type=%s", requestView.Model, upstreamModel, account.Type)
placeholder

	if isCodexSparkModel(upstreamModel) && openAIRequestBodyMayContainImageInput(body) {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if err := validateCodexSparkInput(decoded, upstreamModel); err != nil {
			setOpsUpstreamError(c, http.StatusBadRequest, err.Error(), "")
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"type": "invalid_request_error", "message": err.Error(), "param": "input"placeholderplaceholder)
			return nil, err
	placeholder
placeholder

	// gpt-5.3-codex-spark also rejects the image_generation tool (HTTP 400,
	// param=tools). Strip it here so both APIKey and OAuth /responses paths are
	// covered regardless of the image-generation feature gate.
	if isCodexSparkModel(upstreamModel) && openAIRequestBodyHasImageGenerationDeclaration(body) {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if stripCodexSparkImageGenerationTools(decoded) {
			markDecodedModified()
	placeholder
placeholder

	if account.UsesOpenAICodexProtocol() {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		codexResult := codexTransformResult{placeholder
		if compatMessagesBridge {
			codexResult = applyCodexOAuthTransformWithOptions(decoded, codexOAuthTransformOptions{IsCodexCLI: isCodexCLI, IsCompact: isCompactRequest, SkipDefaultInstructions: true, PreserveToolCallIDs: trueplaceholder)
			ensureCodexOAuthInstructionsField(decoded)
			markDecodedModified()
	placeholder else {
			codexResult = applyCodexOAuthTransform(decoded, isCodexCLI, isCompactRequest)
	placeholder
		if codexResult.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"type": "invalid_request_error", "message": codexResult.Error.Error()placeholderplaceholder)
			return nil, codexResult.Error
	placeholder
		setCodexToolNameReverse(c, codexResult.ToolNameReverse)
		if codexResult.Modified {
			markDecodedModified()
	placeholder
		// 带真实 device_id 时补齐 client_metadata 安装标识，与真实 Codex 对齐（compact 形态不同，跳过）。
		if !isCompactRequest && applyCodexClientMetadata(decoded, account) {
			markDecodedModified()
	placeholder
		stageCodexFingerprintIDs(c, nil)
		// 指纹收敛：一次性解析收敛 ID，请求体和出站头共享同一份 IDs（保证 turn_id 等随机字段一致）。
		// fingerprintIDs 在此处解析，后续 buildUpstreamRequest 中使用同一份。
		if !isCompactRequest {
			var clientHeaders http.Header
			if c != nil && c.Request != nil {
				clientHeaders = c.Request.Header
		placeholder
			fpIDs := resolveCodexFingerprintIDsFromRequest(account, clientHeaders)
			if fpIDs != nil {
				if applyCodexFingerprintClientMetadata(decoded, fpIDs) {
					markDecodedModified()
			placeholder
		placeholder
			// 将 fpIDs 存入 gin context，供 buildUpstreamRequest 中头改写使用。
			// 无条件覆写（含 nil）：failover 从收敛账号切到 off 账号时，上一
			// 账号的 IDs 不得残留（stageCodexFingerprintIDs 注释）。
			stageCodexFingerprintIDs(c, fpIDs)
	placeholder
		if codexResult.NormalizedModel != "" {
			upstreamModel = codexResult.NormalizedModel
	placeholder
		if currentPromptCacheKey, ok := decoded["prompt_cache_key"].(string); ok && currentPromptCacheKey != "" {
			promptCacheKey = currentPromptCacheKey
	placeholder else if codexResult.PromptCacheKey != "" {
			promptCacheKey = codexResult.PromptCacheKey
	placeholder
placeholder

	if !SupportsVerbosity(upstreamModel) && gjson.GetBytes(body, "text.verbosity").Exists() {
		markPatchDelete("text.verbosity")
placeholder

	if !isCodexCLI {
		maxOutputTokens := gjson.GetBytes(body, "max_output_tokens")
		if maxOutputTokens.Exists() {
			switch account.Platform {
			case PlatformOpenAI, PlatformDeepseek:
				// Preserve Responses-native output limits unless the selected upstream
				// explicitly rejects the field in the bounded HTTP retry loop below.
			case PlatformAnthropic:
				decoded, decodeErr := ensureReqBody()
				if decodeErr != nil {
					return nil, decodeErr
			placeholder
				delete(decoded, "max_output_tokens")
				if _, hasMaxTokens := decoded["max_tokens"]; !hasMaxTokens {
					decoded["max_tokens"] = maxOutputTokens.Value()
			placeholder
				markDecodedModified()
			case PlatformGemini:
				markPatchDelete("max_output_tokens")
			default:
				markPatchDelete("max_output_tokens")
		placeholder
	placeholder
		// /v1/responses 的规范输出上限字段是 max_output_tokens；部分客户端仍按
		// Chat Completions 习惯发送 max_tokens，兼容 Responses 上游会拒绝该字段（#4417）。
		// 仅对 OpenAI 平台归一化：Anthropic 合法使用 max_tokens，其 max_output_tokens
		// 反向转换已在上方 switch 中处理。
		if account.Platform == PlatformOpenAI {
			if maxTokens := gjson.GetBytes(body, "max_tokens"); maxTokens.Exists() {
				if !gjson.GetBytes(body, "max_output_tokens").Exists() {
					markPatchSet("max_output_tokens", maxTokens.Value())
			placeholder
				markPatchDelete("max_tokens")
		placeholder
	placeholder
		if gjson.GetBytes(body, "max_completion_tokens").Exists() && (account.Type == AccountTypeAPIKey || account.Platform != PlatformOpenAI) {
			markPatchDelete("max_completion_tokens")
	placeholder
		for _, unsupportedField := range []string{"prompt_cache_retention", "safety_identifier", "prompt_cache_options"placeholder {
			if gjson.GetBytes(body, unsupportedField).Exists() {
				markPatchDelete(unsupportedField)
		placeholder
	placeholder
placeholder
	if wsDecision.Transport != OpenAIUpstreamTransportResponsesWebsocketV2 &&
		!account.IsOpenAIApiKey() && gjson.GetBytes(body, "previous_response_id").Exists() {
		markPatchDelete("previous_response_id")
placeholder
	if openAIRequestBodyMayContainEmptyBase64InputImage(body) {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if sanitizeEmptyBase64InputImagesInOpenAIRequestBodyMap(decoded) {
			markDecodedModified()
	placeholder
placeholder

	if rawTier := requestView.ServiceTier; rawTier != "" {
		if normTier := normalizedOpenAIServiceTierValue(rawTier); normTier != "" {
			action, errMsg := s.evaluateOpenAIFastPolicy(ctx, account, upstreamModel, normTier)
			switch action {
			case BetaPolicyActionBlock:
				msg := errMsg
				if msg == "" {
					msg = fmt.Sprintf("openai service_tier=%s is not allowed for model %s", normTier, upstreamModel)
			placeholder
				blocked := &OpenAIFastBlockedError{Message: msgplaceholder
				writeOpenAIFastPolicyBlockedResponse(c, blocked)
				return nil, blocked
			case BetaPolicyActionFilter:
				markPatchDelete("service_tier")
			case OpenAIFastPolicyActionForcePriority:
				if rawTier != OpenAIFastTierPriority {
					markPatchSet("service_tier", OpenAIFastTierPriority)
			placeholder
			default:
				if normTier != rawTier {
					markPatchSet("service_tier", normTier)
			placeholder
		placeholder
	placeholder
placeholder

	if account.UsesOpenAICodexProtocol() {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if input, ok := decoded["input"].([]any); ok && sanitizeOpenAIResponsesOrphanToolOutputs(
			decoded,
			input,
			strings.TrimSpace(firstNonEmptyString(decoded["previous_response_id"])) != "",
		) {
			markDecodedModified()
	placeholder
placeholder
	if reqBody != nil || openAIResponsesInputMayNeedTruncation(body) {
		decoded, decodeErr := ensureReqBody()
		if decodeErr != nil {
			return nil, decodeErr
	placeholder
		if truncateOpenAIResponsesInputText(decoded) {
			markDecodedModified()
	placeholder
placeholder

	if bodyModified {
		if requestView.HasPatches() {
			if patchedBody, patchErr := requestView.ApplyPatches(); patchErr == nil {
				body = patchedBody
				requestView = newOpenAIRequestView(body)
				reqBody = nil
				bodyModified = false
		placeholder
	placeholder
		if bodyModified {
			decoded, decodeErr := ensureReqBody()
			if decodeErr != nil {
				return nil, decodeErr
		placeholder
			var marshalErr error
			body, marshalErr = marshalOpenAIUpstreamJSON(decoded)
			if marshalErr != nil {
				return nil, fmt.Errorf("serialize request body: %w", marshalErr)
		placeholder
			requestView = newOpenAIRequestView(body)
	placeholder
placeholder
	// Run after orphan-output filtering and all request-map rebuilds so a
	// compaction trigger cannot remain ahead of surviving history items.
	if normalizedBody, changed, normalizeErr := NormalizeCompactionTriggerInputOrder(body); normalizeErr != nil {
		return nil, fmt.Errorf("normalize compaction trigger order: %w", normalizeErr)
placeholder else if changed {
		body = normalizedBody
		requestView = newOpenAIRequestView(body)
		reqBody = nil
placeholder
	imageBillingModel := ""
	imageSizeTier := ""
	imageInputSize := ""
	if imageIntent {
		var imageCfg OpenAIResponsesImageBillingConfig
		var imageCfgErr error
		if reqBody != nil {
			imageCfg, imageCfgErr = resolveOpenAIResponsesImageBillingConfigDetailed(reqBody, billingModel)
	placeholder else {
			imageCfg, imageCfgErr = resolveOpenAIResponsesImageBillingConfigDetailedFromBody(body, billingModel)
	placeholder
		if imageCfgErr != nil {
			setOpsUpstreamError(c, http.StatusBadRequest, imageCfgErr.Error(), "")
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"type": "invalid_request_error", "message": imageCfgErr.Error(), "param": "size"placeholderplaceholder)
			return nil, imageCfgErr
	placeholder
		imageBillingModel = imageCfg.Model
		imageSizeTier = imageCfg.SizeTier
		imageInputSize = imageCfg.InputSize
placeholder
	// Get access token
	token, _, err := s.GetAccessToken(ctx, account)
	if err != nil {
		return nil, err
placeholder
	SetOpsUpstreamModel(c, upstreamModel)

	// 命中 WS 时仅走 WebSocket Mode；不再自动回退 HTTP。
	if wsDecision.Transport == OpenAIUpstreamTransportResponsesWebsocketV2 {
		// WS 分支需要结构化 payload 与重连恢复，命中后再触发 full-map decode。
		wsReqBody, err := ensureReqBody()
		if err != nil {
			return nil, err
	placeholder
		_, hasPreviousResponseID := wsReqBody["previous_response_id"]
		logOpenAIWSModeDebug(
			"forward_start account_id=%d account_type=%s model=%s stream=%v has_previous_response_id=%v",
			account.ID,
			account.Type,
			upstreamModel,
			reqStream,
			hasPreviousResponseID,
		)
		maxAttempts := openAIWSReconnectRetryLimit + 1
		wsAttempts := 0
		var wsResult *OpenAIForwardResult
		var wsErr error
		wsLastFailureReason := ""
		agentTaskRecoveryTried := false
		wsPrevResponseRecoveryTried := false
		wsInvalidEncryptedContentRecoveryTried := false
		recoverPrevResponseNotFound := func(attempt int) bool {
			if wsPrevResponseRecoveryTried {
				return false
		placeholder
			previousResponseID := openAIWSPayloadString(wsReqBody, "previous_response_id")
			if previousResponseID == "" {
				logOpenAIWSModeInfo(
					"reconnect_prev_response_recovery_skip account_id=%d attempt=%d reason=missing_previous_response_id previous_response_id_present=false",
					account.ID,
					attempt,
				)
				return false
		placeholder
			if HasFunctionCallOutput(wsReqBody) {
				logOpenAIWSModeInfo(
					"reconnect_prev_response_recovery_skip account_id=%d attempt=%d reason=has_function_call_output previous_response_id_present=true",
					account.ID,
					attempt,
				)
				return false
		placeholder
			delete(wsReqBody, "previous_response_id")
			wsPrevResponseRecoveryTried = true
			logOpenAIWSModeInfo(
				"reconnect_prev_response_recovery account_id=%d attempt=%d action=drop_previous_response_id retry=1 previous_response_id=%s previous_response_id_kind=%s",
				account.ID,
				attempt,
				truncateOpenAIWSLogValue(previousResponseID, openAIWSIDValueMaxLen),
				normalizeOpenAIWSLogValue(ClassifyOpenAIPreviousResponseIDKind(previousResponseID)),
			)
			return true
	placeholder
		recoverInvalidEncryptedContent := func(attempt int) bool {
			if wsInvalidEncryptedContentRecoveryTried {
				return false
		placeholder
			removedReasoningItems := trimOpenAIEncryptedReasoningItems(wsReqBody)
			if !removedReasoningItems {
				logOpenAIWSModeInfo(
					"reconnect_invalid_encrypted_content_recovery_skip account_id=%d attempt=%d reason=missing_encrypted_reasoning_items",
					account.ID,
					attempt,
				)
				return false
		placeholder
			previousResponseID := openAIWSPayloadString(wsReqBody, "previous_response_id")
			hasFunctionCallOutput := HasFunctionCallOutput(wsReqBody)
			if previousResponseID != "" && !hasFunctionCallOutput {
				delete(wsReqBody, "previous_response_id")
		placeholder
			wsInvalidEncryptedContentRecoveryTried = true
			logOpenAIWSModeInfo(
				"reconnect_invalid_encrypted_content_recovery account_id=%d attempt=%d action=drop_encrypted_reasoning_items retry=1 previous_response_id_present=%v previous_response_id=%s previous_response_id_kind=%s has_function_call_output=%v dropped_previous_response_id=%v",
				account.ID,
				attempt,
				previousResponseID != "",
				truncateOpenAIWSLogValue(previousResponseID, openAIWSIDValueMaxLen),
				normalizeOpenAIWSLogValue(ClassifyOpenAIPreviousResponseIDKind(previousResponseID)),
				hasFunctionCallOutput,
				previousResponseID != "" && !hasFunctionCallOutput,
			)
			return true
	placeholder
		retryBudget := s.openAIWSRetryTotalBudget()
		retryStartedAt := time.Now()
	wsRetryLoop:
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			wsAttempts = attempt
			wsResult, wsErr = s.forwardOpenAIWSV2(
				ctx,
				c,
				account,
				wsReqBody,
				token,
				wsDecision,
				isCodexCLI,
				reqStream,
				originalModel,
				upstreamModel,
				startTime,
				attempt,
				wsLastFailureReason,
				&agentTaskRecoveryTried,
			)
			if wsErr == nil {
				break
		placeholder
			if c != nil && c.Writer != nil && c.Writer.Written() {
				break
		placeholder
			var taskRecoveredErr *agentIdentityTaskRecoveredError
			if errors.As(wsErr, &taskRecoveredErr) {
				continue
		placeholder

			reason, retryable := classifyOpenAIWSReconnectReason(wsErr)
			if reason != "" {
				wsLastFailureReason = reason
		placeholder
			// previous_response_not_found 说明续链锚点不可用：
			// 对非 function_call_output 场景，允许一次“去掉 previous_response_id 后重放”。
			if reason == "previous_response_not_found" && recoverPrevResponseNotFound(attempt) {
				continue
		placeholder
			if reason == "invalid_encrypted_content" && recoverInvalidEncryptedContent(attempt) {
				continue
		placeholder
			if retryable && attempt < maxAttempts {
				backoff := s.openAIWSRetryBackoff(attempt)
				if retryBudget > 0 && time.Since(retryStartedAt)+backoff > retryBudget {
					s.recordOpenAIWSRetryExhausted()
					logOpenAIWSModeInfo(
						"reconnect_budget_exhausted account_id=%d attempts=%d max_retries=%d reason=%s elapsed_ms=%d budget_ms=%d",
						account.ID,
						attempt,
						openAIWSReconnectRetryLimit,
						normalizeOpenAIWSLogValue(reason),
						time.Since(retryStartedAt).Milliseconds(),
						retryBudget.Milliseconds(),
					)
					break
			placeholder
				s.recordOpenAIWSRetryAttempt(backoff)
				logOpenAIWSModeInfo(
					"reconnect_retry account_id=%d retry=%d max_retries=%d reason=%s backoff_ms=%d",
					account.ID,
					attempt,
					openAIWSReconnectRetryLimit,
					normalizeOpenAIWSLogValue(reason),
					backoff.Milliseconds(),
				)
				if backoff > 0 {
					timer := time.NewTimer(backoff)
					select {
					case <-ctx.Done():
						if !timer.Stop() {
							<-timer.C
					placeholder
						wsErr = wrapOpenAIWSFallback("retry_backoff_canceled", ctx.Err())
						break wsRetryLoop
					case <-timer.C:
				placeholder
			placeholder
				continue
		placeholder
			if retryable {
				s.recordOpenAIWSRetryExhausted()
				logOpenAIWSModeInfo(
					"reconnect_exhausted account_id=%d attempts=%d max_retries=%d reason=%s",
					account.ID,
					attempt,
					openAIWSReconnectRetryLimit,
					normalizeOpenAIWSLogValue(reason),
				)
		placeholder else if reason != "" {
				s.recordOpenAIWSNonRetryableFastFallback()
				logOpenAIWSModeInfo(
					"reconnect_stop account_id=%d attempt=%d reason=%s",
					account.ID,
					attempt,
					normalizeOpenAIWSLogValue(reason),
				)
		placeholder
			break
	placeholder
		if wsErr == nil {
			firstTokenMs := int64(0)
			hasFirstTokenMs := wsResult != nil && wsResult.FirstTokenMs != nil
			if hasFirstTokenMs {
				firstTokenMs = int64(*wsResult.FirstTokenMs)
		placeholder
			requestID := ""
			if wsResult != nil {
				requestID = strings.TrimSpace(wsResult.RequestID)
		placeholder
			logOpenAIWSModeDebug(
				"forward_succeeded account_id=%d request_id=%s stream=%v has_first_token_ms=%v first_token_ms=%d ws_attempts=%d",
				account.ID,
				requestID,
				reqStream,
				hasFirstTokenMs,
				firstTokenMs,
				wsAttempts,
			)
			wsResult.UpstreamModel = upstreamModel
			if wsResult.BillingModel == "" {
				wsResult.BillingModel = billingModel
		placeholder
			if wsResult.ImageCount > 0 {
				wsResult.ImageSize = imageSizeTier
				wsResult.ImageInputSize = imageInputSize
				wsResult.BillingModel = imageBillingModel
		placeholder
			return wsResult, nil
	placeholder
		s.writeOpenAIWSFallbackErrorResponse(c, account, wsErr)
		return nil, wsErr
placeholder

	reasoningEffort := extractOpenAIReasoningEffortFromBody(body, upstreamModel, billingModel, originalModel)
	// 国产模型默认 effort 补充：此处 reqModel 已被 mapping 重写为 billingModel。
	reasoningEffort = ApplyThinkingEnabledFallback(reasoningEffort, body, reqModel)
	reasoningEffortValue := ""
	if reasoningEffort != nil {
		reasoningEffortValue = *reasoningEffort
placeholder
	firstOutputTimeout := time.Duration(0)
	if reqStream && account.Platform == PlatformOpenAI {
		firstOutputTimeout = s.openAIFirstOutputTimeout(reasoningEffortValue)
placeholder

	httpInvalidEncryptedContentRetryTried := false
	compactModelFallbackRetried := false
	agentTaskRecoveryTried := false
	rejectedFieldRetryState := openAIResponsesRejectedFieldRetryStateForRequest(c, body)
	for {
		// Build upstream request
		upstreamCtx, releaseUpstreamCtx := detachUpstreamContext(ctx)
		var headerGuard *openAIFirstOutputHeaderGuard
		if firstOutputTimeout > 0 {
			upstreamCtx, headerGuard = newOpenAIFirstOutputHeaderGuard(
				upstreamCtx, releaseUpstreamCtx, startTime.Add(firstOutputTimeout),
			)
	placeholder
		upstreamReq, err := s.buildUpstreamRequest(upstreamCtx, c, account, body, token, reqStream, promptCacheKey, isCodexCLI)
		if headerGuard == nil {
			releaseUpstreamCtx()
	placeholder
		if err != nil {
			if headerGuard != nil {
				headerGuard.close()
		placeholder
			return nil, err
	placeholder

		// Get proxy URL
		proxyURL := ""
		if account.ProxyID != nil && account.Proxy != nil {
			proxyURL = account.Proxy.URL()
	placeholder

		// Send request
		upstreamStart := time.Now()
		resp, err := s.httpUpstream.Do(upstreamReq, proxyURL, account.ID, account.Concurrency)
		SetOpsLatencyMs(c, OpsUpstreamLatencyMsKey, time.Since(upstreamStart).Milliseconds())
		if headerGuard != nil && headerGuard.stopHeaderWait() {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
		placeholder
			headerGuard.close()
			return nil, s.newOpenAIFirstOutputTimeoutError(
				ctx, c, account, startTime, originalModel, reasoningEffortValue,
				firstOutputTimeout, "response_headers", nil,
			)
	placeholder
		if err != nil {
			if resp != nil && resp.Body != nil {
				_ = resp.Body.Close()
		placeholder
			if headerGuard != nil {
				headerGuard.close()
		placeholder
			// Transport-level failure (proxy/DNS/TCP/TLS — no HTTP response). Convert to
			// a failover so the handler switches to a healthy account, and temporarily
			// unschedule the account on durable faults (e.g. rejected proxy credentials).
			return nil, s.handleOpenAIUpstreamTransportError(ctx, c, account, err, false)
	placeholder
		if headerGuard != nil {
			resp.Body = &openAIRequestContextReadCloser{ReadCloser: resp.Body, cleanup: headerGuard.closeplaceholder
	placeholder

		// Handle error response
		if resp.StatusCode >= 400 {
			respBody := s.readUpstreamErrorBody(resp)
			_ = resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(respBody))

			upstreamMsg := strings.TrimSpace(extractUpstreamErrorMessage(respBody))
			upstreamMsg = sanitizeUpstreamErrorMessage(upstreamMsg)
			upstreamCode := extractUpstreamErrorCode(respBody)
			if !agentTaskRecoveryTried && s.isAgentIdentityAccount(ctx, account) && isAgentIdentityTaskInvalidHTTPResponse(resp.StatusCode, respBody) {
				agentTaskRecoveryTried = true
				expectedTaskID := account.GetCredential("task_id")
				if err := s.recoverAgentIdentityTask(ctx, account, expectedTaskID); err != nil {
					return nil, fmt.Errorf("agent identity task recovery failed: %w", err)
			placeholder
				continue
		placeholder
			respBody = s.redactAgentIdentitySensitiveBody(ctx, account, respBody)
			resp.Body = io.NopCloser(bytes.NewReader(respBody))
			if !httpInvalidEncryptedContentRetryTried && resp.StatusCode == http.StatusBadRequest && upstreamCode == "invalid_encrypted_content" {
				decoded, decodeErr := ensureReqBody()
				if decodeErr != nil {
					return nil, decodeErr
			placeholder
				if trimOpenAIEncryptedReasoningItems(decoded) {
					body, err = marshalOpenAIUpstreamJSON(decoded)
					if err != nil {
						return nil, fmt.Errorf("serialize invalid_encrypted_content retry body: %w", err)
				placeholder
					httpInvalidEncryptedContentRetryTried = true
					rejectedFieldRetryState.remember(body)
					logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Retrying non-WSv2 request once after invalid_encrypted_content (account: %s)", account.Name)
					continue
			placeholder
				logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Skip non-WSv2 invalid_encrypted_content retry because encrypted reasoning items are missing (account: %s)", account.Name)
		placeholder
			if retryBody, reason, changed, retryErr := normalizeOpenAIResponsesRejectedFieldRetryBody(resp.StatusCode, body, respBody); retryErr != nil {
				return nil, fmt.Errorf("normalize rejected Responses field retry body: %w", retryErr)
		placeholder else if changed && rejectedFieldRetryState.Allow(retryBody) {
				body = retryBody
				requestView = newOpenAIRequestView(body)
				reqBody = nil
				logger.LegacyPrintf("service.openai_gateway", "[OpenAI] Retrying non-WSv2 request after %s (account: %s)", reason, account.Name)
				continue
		placeholder
			if retryBody, fallbackModel, retry := s.prepareOpenAICompactFallbackRetry(
				c, account, requestedModel, body, resp.StatusCode, upstreamMsg, respBody, compactModelFallbackRetried,
			); retry {
				s.appendOpenAICompactFallbackRetryOps(c, account, resp, respBody, upstreamMsg, false)
				fromModel := strings.TrimSpace(gjson.GetBytes(body, "model").String())
				body = retryBody
				requestView = newOpenAIRequestView(body)
				reqBody = nil
				upstreamModel = fallbackModel
				compactModelFallbackRetried = true
				SetOpsUpstreamModel(c, fallbackModel)
				logger.LegacyPrintf(
					"service.openai_gateway",
					"[OpenAI] Retrying explicit compact request once with fallback model (account: %s, from: %s, to: %s, upstream_code: %s)",
					account.Name, fromModel, fallbackModel, upstreamCode,
				)
				continue
		placeholder
			if s.shouldFailoverOpenAIUpstreamResponse(resp.StatusCode, upstreamMsg, respBody) {
				upstreamDetail := ""
				if s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
					maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
					if maxBytes <= 0 {
						maxBytes = 2048
				placeholder
					upstreamDetail = truncateString(string(respBody), maxBytes)
			placeholder
				appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
					Platform:           account.Platform,
					AccountID:          account.ID,
					AccountName:        account.Name,
					UpstreamStatusCode: resp.StatusCode,
					UpstreamRequestID:  resp.Header.Get("x-request-id"),
					Kind:               "failover",
					Message:            upstreamMsg,
					Detail:             upstreamDetail,
			placeholder)

				shouldDisable := s.handleFailoverSideEffects(ctx, resp, account, respBody, upstreamModel)
				return nil, s.newOpenAIAccountFailoverError(
					account,
					resp.StatusCode,
					resp.Header,
					respBody,
					upstreamMsg,
					shouldDisable,
					!shouldDisable && account.IsPoolMode() && (account.IsPoolModeRetryableStatus(resp.StatusCode) || isOpenAITransientProcessingError(resp.StatusCode, upstreamMsg, respBody)),
				)
		placeholder
			return s.handleErrorResponse(ctx, resp, c, account, body, resolveOpenAIErrorSchedulingModel(billingModel, upstreamModel))
	placeholder
		defer func() { _ = resp.Body.Close() placeholder()

		if mapping, ok := openAIResponsesClientToolMapping(c); ok && isEventStreamResponse(resp.Header) {
			maxLineSize := defaultMaxLineSize
			if s.cfg != nil && s.cfg.Gateway.MaxLineSize > 0 {
				maxLineSize = s.cfg.Gateway.MaxLineSize
		placeholder
			resp.Body = newResponsesClientToolStreamBody(resp.Body, mapping, maxLineSize)
	placeholder

		serviceTier := extractOpenAIServiceTierFromBody(body)
		// 上游接受后只保留计费需要的标量，避免响应处理期间继续保活完整 input/tools map。
		reqBody = nil

		// Handle normal response
		var usage *OpenAIUsage
		var firstTokenMs *int
		responseID := ""
		imageCount := 0
		searchCount := 0
		var imageOutputSizes []string
		if reqStream {
			streamResult, err := s.handleStreamingResponseWithReasoning(ctx, resp, c, account, startTime, originalModel, upstreamModel, reasoningEffortValue)
			if err != nil {
				if signal, ok := asOpenAICompactFallbackSignal(err); ok {
					if retryBody, fallbackModel, retry := s.prepareOpenAICompactFallbackRetry(
						c, account, requestedModel, body, http.StatusBadRequest, signal.message, signal.payload, compactModelFallbackRetried,
					); retry {
						s.appendOpenAICompactFallbackRetryOps(c, account, resp, signal.payload, signal.message, false)
						body = retryBody
						requestView = newOpenAIRequestView(body)
						upstreamModel = fallbackModel
						compactModelFallbackRetried = true
						SetOpsUpstreamModel(c, fallbackModel)
						continue
				placeholder
					if resp.Body != nil {
						_ = resp.Body.Close()
				placeholder
					compactResp, compactBody := openAICompactFallbackErrorResponse(resp, signal)
					if s.shouldFailoverOpenAIUpstreamResponse(compactResp.StatusCode, signal.message, compactBody) {
						appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
							Platform:           account.Platform,
							AccountID:          account.ID,
							AccountName:        account.Name,
							UpstreamStatusCode: compactResp.StatusCode,
							UpstreamRequestID:  compactResp.Header.Get("x-request-id"),
							Kind:               "failover",
							Message:            signal.message,
					placeholder)
						shouldDisable := s.handleFailoverSideEffects(ctx, compactResp, account, compactBody, upstreamModel)
						return nil, s.newOpenAIAccountFailoverError(
							account, compactResp.StatusCode, compactResp.Header, compactBody, signal.message, shouldDisable,
							!shouldDisable && account.IsPoolMode() && (account.IsPoolModeRetryableStatus(compactResp.StatusCode) || isOpenAITransientProcessingError(compactResp.StatusCode, signal.message, compactBody)),
						)
				placeholder
					return s.handleErrorResponse(ctx, compactResp, c, account, body, resolveOpenAIErrorSchedulingModel(billingModel, upstreamModel))
			placeholder
				return nil, err
		placeholder
			usage = streamResult.usage
			firstTokenMs = streamResult.firstTokenMs
			responseID = strings.TrimSpace(streamResult.responseID)
			imageCount = streamResult.imageCount
			imageOutputSizes = streamResult.imageOutputSizes
			searchCount = streamResult.searchCount
	placeholder else {
			nonStreamResult, err := s.handleNonStreamingResponse(ctx, resp, c, account, originalModel, upstreamModel)
			if err != nil {
				if signal, ok := asOpenAICompactFallbackSignal(err); ok {
					if retryBody, fallbackModel, retry := s.prepareOpenAICompactFallbackRetry(
						c, account, requestedModel, body, http.StatusBadRequest, signal.message, signal.payload, compactModelFallbackRetried,
					); retry {
						body = retryBody
						requestView = newOpenAIRequestView(body)
						upstreamModel = fallbackModel
						compactModelFallbackRetried = true
						SetOpsUpstreamModel(c, fallbackModel)
						continue
				placeholder
			placeholder
				return nil, err
		placeholder
			usage = nonStreamResult.usage
			responseID = strings.TrimSpace(nonStreamResult.responseID)
			imageCount = nonStreamResult.imageCount
			imageOutputSizes = nonStreamResult.imageOutputSizes
			searchCount = nonStreamResult.searchCount
	placeholder
		s.bindHTTPResponseAccount(ctx, c, account, responseID)

		// Extract and save Codex usage snapshot from response headers (for OAuth accounts).
		// 排除 spark 影子:其 codex_* 仅由 QueryUsage(/wham/usage bengalfox)更新(外审第7轮 P1)。
		if account.UsesOpenAICodexProtocol() && !account.IsShadow() {
			if snapshot := ParseCodexRateLimitHeaders(resp.Header); snapshot != nil {
				s.updateCodexUsageSnapshot(ctx, account.ID, snapshot)
		placeholder
	placeholder

		if usage == nil {
			usage = &OpenAIUsage{placeholder
	placeholder

		forwardResult := &OpenAIForwardResult{
			RequestID:                     resp.Header.Get("x-request-id"),
			ResponseID:                    responseID,
			Usage:                         *usage,
			Model:                         originalModel,
			BillingModel:                  billingModel,
			UpstreamModel:                 upstreamModel,
			UpstreamResponseModel:         observedUpstreamResponseModel(c),
			UpstreamResponseModelConflict: observedUpstreamResponseModelConflict(c),
			ServiceTier:                   serviceTier,
			ReasoningEffort:               reasoningEffort,
			Stream:                        reqStream,
			OpenAIWSMode:                  false,
			Duration:                      time.Since(startTime),
			FirstTokenMs:                  firstTokenMs,
	placeholder
		if imageCount > 0 {
			forwardResult.ImageCount = imageCount
			forwardResult.ImageSize = imageSizeTier
			forwardResult.ImageInputSize = imageInputSize
			forwardResult.ImageOutputSizes = imageOutputSizes
			forwardResult.BillingModel = imageBillingModel
	placeholder
		// Grok-native web_search / x_search / tool_search tool invocations (per-1k pricing).
		// Token cost still applies separately when usage is present; search is additive only
		// when search_price_per_1k is configured (nil price → $0 from CalculateSearchCost).
		if searchCount > 0 && account != nil && account.IsGrok() {
			forwardResult.SearchCount = searchCount
	placeholder
		return forwardResult, nil
placeholder
placeholder

func shouldForwardOpenAIResponsesViaRawChatCompletions(account *Account) bool {
	if account == nil || account.Type != AccountTypeAPIKey {
		return false
placeholder
	if account.IsCNProvider() {
		// CN 的显式协议配置优先于异步探针 Extra；adaptive 仅 DeepSeek 有原生
		// Responses，Kimi/GLM 回退 Chat Completions。
		switch account.GetAPIProtocol() {
		case APIProtocolChatCompletions:
			return true
		case APIProtocolAdaptive:
			return account.Platform != PlatformDeepseek
		default:
			return false
	placeholder
placeholder
	return !openai_compat.ShouldUseResponsesAPI(account.Extra)
placeholder

func (s *OpenAIGatewayService) buildUpstreamRequest(ctx context.Context, c *gin.Context, account *Account, body []byte, token string, isStream bool, promptCacheKey string, isCodexCLI bool) (*http.Request, error) {
	// Determine target URL based on account type
	var targetURL string
	switch account.Type {
	case AccountTypeOAuth:
		// OAuth accounts use ChatGPT internal API
		targetURL = chatgptCodexURL
	case AccountTypeSetupToken:
		if account.IsOpenAIOAuthLike() {
			targetURL = chatgptCodexURL
	placeholder else {
			targetURL = openaiPlatformAPIURL
	placeholder
	case AccountTypeAPIKey:
		// API Key accounts use Platform API or custom base URL
		baseURL := account.GetOpenAIBaseURL()
		if account.Platform == PlatformDeepseek && account.IsAdaptiveAPIProtocol() {
			baseURL = account.GetCNProtocolBaseURL(APIProtocolResponses)
	placeholder
		if baseURL == "" {
			targetURL = openaiPlatformAPIURL
	placeholder else {
			validatedURL, err := s.validateUpstreamBaseURL(baseURL)
			if err != nil {
				return nil, err
		placeholder
			targetURL = buildOpenAIResponsesURLForPlatform(account.Platform, validatedURL)
	placeholder
	default:
		targetURL = openaiPlatformAPIURL
placeholder
	targetURL = appendOpenAIResponsesRequestPathSuffix(targetURL, openAIResponsesRequestPathSuffix(c))

	// DeepSeek 原生 Responses 端点为无状态实现：强制 store=false、清除
	// previous_response_id，避免携带状态字段被上游拒绝。
	body = normalizeDeepSeekResponsesRequestBody(account, body)

	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
placeholder
	req = req.WithContext(WithHTTPUpstreamProfile(req.Context(), HTTPUpstreamProfileOpenAI))

	// Build authentication for this request. Agent Identity signs a fresh
	// assertion here; OAuth/PAT/API-key keep their existing Bearer behavior.
	authHeaders, err := s.buildOpenAIAuthenticationHeaders(ctx, account, token)
	if err != nil {
		return nil, fmt.Errorf("build openai authentication headers: %w", err)
placeholder
	for key, values := range authHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
	placeholder
placeholder

	// Set headers specific to OAuth accounts (ChatGPT internal API)
	if account.UsesOpenAICodexProtocol() {
		// Required: set Host for ChatGPT API (must use req.Host, not Header.Set)
		req.Host = "chatgpt.com"
		if err := resolveAndSetOpenAIChatGPTAccountHeaders(ctx, s.accountRepo, req.Header, account); err != nil {
			return nil, fmt.Errorf("resolve chatgpt account headers: %w", err)
	placeholder
placeholder

	// Whitelist passthrough headers
	for key, values := range c.Request.Header {
		lowerKey := strings.ToLower(key)
		if openaiAllowedHeaders[lowerKey] {
			for _, v := range values {
				req.Header.Add(key, v)
		placeholder
	placeholder
placeholder
	// 客户端回带的 x-codex-turn-state 若已知由其他账号铸造（failover 换号），
	// 剥离后再出站——异账号 blob 与本账号的（指纹收敛后）出站身份自相矛盾。
	s.guardOpenAICodexTurnStateEcho(c, account, req.Header)
	if account.UsesOpenAICodexProtocol() {
		compatMessagesBridge := isOpenAICompatMessagesBridgeContext(c) || isOpenAICompatMessagesBridgeBody(body)
		// 清除客户端透传的 session 头，后续用隔离后的值重新设置，防止跨用户会话碰撞。
		clientConversationID := strings.TrimSpace(req.Header.Get("conversation_id"))
		req.Header.Del("conversation_id")
		req.Header.Del("session_id")

		if compatMessagesBridge {
			req.Header.Del("OpenAI-Beta")
			req.Header.Del("originator")
	placeholder else {
			req.Header.Set("originator", resolveOpenAIUpstreamOriginator(c, isCodexCLI))
	placeholder
		apiKeyID := getAPIKeyIDFromContext(c)
		if isOpenAIResponsesCompactPath(c) {
			req.Header.Set("accept", "application/json")
			if req.Header.Get("version") == "" {
				req.Header.Set("version", CodexCanonicalClientVersion())
		placeholder
			compactSession := resolveOpenAICompactSessionID(c)
			req.Header.Set("session_id", isolateOpenAISessionID(apiKeyID, compactSession))
	placeholder else {
			req.Header.Set("accept", "text/event-stream")
	placeholder
		if promptCacheKey != "" {
			isolated := isolateOpenAISessionID(apiKeyID, promptCacheKey)
			req.Header.Set("session_id", isolated)
			if !compatMessagesBridge || clientConversationID != "" {
				req.Header.Set("conversation_id", isolated)
		placeholder
	placeholder
placeholder else if isOpenAIResponsesCompactPath(c) {
		// compact 上游是 unary JSON 协议：API-key 账号也显式声明 Accept，
		// 避免 OpenAI 兼容网关按 SSE 返回（#3777 期望行为 4）。
		req.Header.Set("accept", "application/json")
placeholder

	// Apply custom User-Agent if configured
	customUA := account.GetOpenAIUserAgent()
	if customUA != "" {
		req.Header.Set("user-agent", customUA)
placeholder

	// 若开启 ForceCodexCLI，则强制将上游 User-Agent 伪装为规范 Codex 身份。
	// 用于网关未透传/改写 User-Agent 时，仍能命中 Codex 侧识别逻辑。
	if s.cfg != nil && s.cfg.Gateway.ForceCodexCLI {
		req.Header.Set("user-agent", CodexCanonicalUserAgent())
placeholder

	// 指纹收敛：使用 Forward() 中预计算的收敛 ID 改写出站头，与请求体使用同一份 IDs。
	applyStagedCodexFingerprintHeaders(c, account, req.Header)

	// 终态收口：强制统一 OAuth 出站身份（User-Agent / originator / version 同源自洽）。
	// 客户端自报身份不参与构造，浏览器型 UA 也因此不会再到达上游（原浏览器 UA 兜底已被吸收）。
	if account.UsesOpenAICodexProtocol() {
		enforceCodexIdentityHeadersWithUA(req.Header, s.codexIdentityOverrideUA(account))
placeholder

	// Ensure required headers exist
	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/json")
placeholder

	// 账号级请求头覆写（仅 openai api_key 账号启用时生效；OAuth 路径 no-op）
	account.ApplyHeaderOverrides(req.Header)
	// x-codex-beta-features：按真实 Codex 的会话级行为补注（在账号级覆写之后，
	// 保证不被覆盖丢失）。
	applyOpenAICodexBetaFeatures(c, account, req.Header)
	setOpenAICodexRoutingHintFromBody(req.Header, account, body)
	logOpenAIRoutingDiagnosticsFromBody(ctx, account, "http", req.Header, body, "not_applicable")

	return req, nil
placeholder

// codexIdentityOverrideUA 返回账号级显式配置的出站 User-Agent，供强制统一身份时作为覆写来源。
// ForceCodexCLI 语义是「强制使用 Codex CLI 身份」，等价于使用网关规范身份，故返回空串；
// 该优先级与历史行为一致（ForceCodexCLI 在账号自定义 UA 之后生效）。
func (s *OpenAIGatewayService) codexIdentityOverrideUA(account *Account) string {
	if s != nil && s.cfg != nil && s.cfg.Gateway.ForceCodexCLI {
		return ""
placeholder
	return account.GetOpenAIUserAgent()
placeholder
