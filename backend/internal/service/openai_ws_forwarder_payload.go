package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func validateOpenAIWSBearerToken(account *Account, token string) error {
	if account == nil {
		return errors.New("account is nil")
placeholder
	if strings.TrimSpace(token) == "" && !account.IsOpenAIAgentIdentity() {
		return errors.New("token is empty")
placeholder
	return nil
placeholder

func (s *OpenAIGatewayService) buildOpenAIResponsesWSURL(account *Account) (string, error) {
	if account == nil {
		return "", errors.New("account is nil")
placeholder
	var targetURL string
	switch account.Type {
	case AccountTypeOAuth:
		targetURL = chatgptCodexURL
	case AccountTypeAPIKey:
		baseURL := account.GetOpenAIBaseURL()
		if baseURL == "" {
			targetURL = openaiPlatformAPIURL
	placeholder else {
			validatedURL, err := s.validateUpstreamBaseURL(baseURL)
			if err != nil {
				return "", err
		placeholder
			targetURL = buildOpenAIResponsesURL(validatedURL)
	placeholder
	default:
		targetURL = openaiPlatformAPIURL
placeholder

	parsed, err := url.Parse(strings.TrimSpace(targetURL))
	if err != nil {
		return "", fmt.Errorf("invalid target url: %w", err)
placeholder
	switch strings.ToLower(parsed.Scheme) {
	case "https":
		parsed.Scheme = "wss"
	case "http":
		parsed.Scheme = "ws"
	case "wss", "ws":
		// 保持不变
	default:
		return "", fmt.Errorf("unsupported scheme for ws: %s", parsed.Scheme)
placeholder
	return parsed.String(), nil
placeholder

func (s *OpenAIGatewayService) buildOpenAIWSHeaders(
	ctx context.Context,
	c *gin.Context,
	account *Account,
	token string,
	decision OpenAIWSProtocolDecision,
	isCodexCLI bool,
	turnState string,
	turnMetadata string,
	promptCacheKey string,
	routingModel string,
	routingServiceTier string,
) (http.Header, openAIWSSessionHeaderResolution, error) {
	headers := make(http.Header)
	if account == nil || !account.IsOpenAIAgentIdentity() {
		headers.Set("authorization", "Bearer "+token)
placeholder

	sessionResolution := resolveOpenAIWSSessionHeaders(c, promptCacheKey)
	if c != nil && c.Request != nil {
		if v := strings.TrimSpace(c.Request.Header.Get("accept-language")); v != "" {
			headers.Set("accept-language", v)
	placeholder
		for _, value := range c.Request.Header.Values("x-codex-beta-features") {
			if value = strings.TrimSpace(value); value != "" {
				headers.Add("x-codex-beta-features", value)
		placeholder
	placeholder
		for _, name := range [...]string{"x-codex-window-id", "x-codex-installation-id"placeholder {
			if value := c.Request.Header.Get(name); strings.TrimSpace(value) != "" {
				headers.Set(name, value)
		placeholder
	placeholder
placeholder
	// OAuth 账号：将 apiKeyID 混入 session 标识符，防止跨用户会话碰撞。
	if account != nil && account.Type == AccountTypeOAuth {
		apiKeyID := getAPIKeyIDFromContext(c)
		if sessionResolution.SessionID != "" {
			headers.Set("session_id", isolateOpenAISessionID(apiKeyID, sessionResolution.SessionID))
	placeholder
		if sessionResolution.ConversationID != "" {
			headers.Set("conversation_id", isolateOpenAISessionID(apiKeyID, sessionResolution.ConversationID))
	placeholder
placeholder else {
		if sessionResolution.SessionID != "" {
			headers.Set("session_id", sessionResolution.SessionID)
	placeholder
		if sessionResolution.ConversationID != "" {
			headers.Set("conversation_id", sessionResolution.ConversationID)
	placeholder
placeholder
	if state := strings.TrimSpace(turnState); state != "" {
		headers.Set(openAIWSTurnStateHeader, state)
placeholder
	if metadata := strings.TrimSpace(turnMetadata); metadata != "" {
		headers.Set(openAIWSTurnMetadataHeader, metadata)
placeholder

	if account != nil && account.Type == AccountTypeOAuth {
		if err := resolveAndSetOpenAIChatGPTAccountHeaders(ctx, s.accountRepo, headers, account); err != nil {
			return nil, sessionResolution, fmt.Errorf("resolve chatgpt account headers: %w", err)
	placeholder
		headers.Set("originator", resolveOpenAIUpstreamOriginator(c, isCodexCLI))
placeholder

	betaValue := openAIWSBetaV2Value
	if decision.Transport == OpenAIUpstreamTransportResponsesWebsocket {
		betaValue = openAIWSBetaV1Value
placeholder
	headers.Set("OpenAI-Beta", betaValue)

	customUA := ""
	if account != nil {
		customUA = account.GetOpenAIUserAgent()
placeholder
	if strings.TrimSpace(customUA) != "" {
		headers.Set("user-agent", customUA)
placeholder else if c != nil {
		if ua := strings.TrimSpace(c.GetHeader("User-Agent")); ua != "" {
			headers.Set("user-agent", ua)
	placeholder
placeholder
	if s != nil && s.cfg != nil && s.cfg.Gateway.ForceCodexCLI {
		headers.Set("user-agent", codexCLIUserAgent)
placeholder
	// 终态收口：WS 握手与 HTTP 出站共用同一套身份语义，账号级自定义 UA 同样作为
	// 管理员显式配置传入（上面写进 headers 的值只在强制统一被关闭时才参与配对）。
	if account != nil && account.Type == AccountTypeOAuth {
		enforceCodexIdentityHeadersWithUA(headers, s.codexIdentityOverrideUA(account))
placeholder

	// 账号级请求头覆写（仅 openai api_key 账号启用时生效；OAuth 路径 no-op）。
	// 覆盖所有 WS 模式（ctx_pool/dedicated/passthrough）的握手头。
	account.ApplyHeaderOverrides(headers)
	setOpenAICodexRoutingHint(headers, account, routingModel, routingServiceTier)

	return headers, sessionResolution, nil
placeholder

func (s *OpenAIGatewayService) buildOpenAIWSCreatePayload(reqBody map[string]any, account *Account) map[string]any {
	// OpenAI WS Mode 协议：response.create 字段与 HTTP /responses 基本一致。
	// 保留 stream 字段（与 Codex CLI 一致），仅移除 background。
	payload := make(map[string]any, len(reqBody)+1)
	for k, v := range reqBody {
		payload[k] = v
placeholder

	delete(payload, "background")
	if _, exists := payload["stream"]; !exists {
		payload["stream"] = true
placeholder
	payload["type"] = "response.create"

	// OAuth 默认保持 store=false，避免误依赖服务端历史。
	if account != nil && account.Type == AccountTypeOAuth && !s.isOpenAIWSStoreRecoveryAllowed(account) {
		payload["store"] = false
placeholder
	return payload
placeholder

func setOpenAIWSTurnMetadata(payload map[string]any, turnMetadata string) {
	if len(payload) == 0 {
		return
placeholder
	metadata := strings.TrimSpace(turnMetadata)
	if metadata == "" {
		return
placeholder

	switch existing := payload["client_metadata"].(type) {
	case map[string]any:
		existing[openAIWSTurnMetadataHeader] = metadata
		payload["client_metadata"] = existing
	case map[string]string:
		next := make(map[string]any, len(existing)+1)
		for k, v := range existing {
			next[k] = v
	placeholder
		next[openAIWSTurnMetadataHeader] = metadata
		payload["client_metadata"] = next
	default:
		payload["client_metadata"] = map[string]any{
			openAIWSTurnMetadataHeader: metadata,
	placeholder
placeholder
placeholder

func (s *OpenAIGatewayService) isOpenAIWSStoreRecoveryAllowed(account *Account) bool {
	if account != nil && account.IsOpenAIWSAllowStoreRecoveryEnabled() {
		return true
placeholder
	if s != nil && s.cfg != nil && s.cfg.Gateway.OpenAIWS.AllowStoreRecovery {
		return true
placeholder
	return false
placeholder

func (s *OpenAIGatewayService) isOpenAIWSStoreDisabledInRequest(reqBody map[string]any, account *Account) bool {
	if account != nil && account.Type == AccountTypeOAuth && !s.isOpenAIWSStoreRecoveryAllowed(account) {
		return true
placeholder
	if len(reqBody) == 0 {
		return false
placeholder
	rawStore, ok := reqBody["store"]
	if !ok {
		return false
placeholder
	storeEnabled, ok := rawStore.(bool)
	if !ok {
		return false
placeholder
	return !storeEnabled
placeholder

func (s *OpenAIGatewayService) isOpenAIWSStoreDisabledInRequestRaw(reqBody []byte, account *Account) bool {
	if account != nil && account.Type == AccountTypeOAuth && !s.isOpenAIWSStoreRecoveryAllowed(account) {
		return true
placeholder
	if len(reqBody) == 0 {
		return false
placeholder
	storeValue := gjson.GetBytes(reqBody, "store")
	if !storeValue.Exists() {
		return false
placeholder
	if storeValue.Type != gjson.True && storeValue.Type != gjson.False {
		return false
placeholder
	return !storeValue.Bool()
placeholder

func (s *OpenAIGatewayService) openAIWSStoreDisabledConnMode() string {
	if s == nil || s.cfg == nil {
		return openAIWSStoreDisabledConnModeStrict
placeholder
	mode := strings.ToLower(strings.TrimSpace(s.cfg.Gateway.OpenAIWS.StoreDisabledConnMode))
	switch mode {
	case openAIWSStoreDisabledConnModeStrict, openAIWSStoreDisabledConnModeAdaptive, openAIWSStoreDisabledConnModeOff:
		return mode
	case "":
		// 兼容旧配置：仅配置了布尔开关时按旧语义推导。
		if s.cfg.Gateway.OpenAIWS.StoreDisabledForceNewConn {
			return openAIWSStoreDisabledConnModeStrict
	placeholder
		return openAIWSStoreDisabledConnModeOff
	default:
		return openAIWSStoreDisabledConnModeStrict
placeholder
placeholder

func shouldForceNewConnOnStoreDisabled(mode, lastFailureReason string) bool {
	switch mode {
	case openAIWSStoreDisabledConnModeOff:
		return false
	case openAIWSStoreDisabledConnModeAdaptive:
		reason := strings.TrimPrefix(strings.TrimSpace(lastFailureReason), "prewarm_")
		switch reason {
		case "policy_violation", "message_too_big", "auth_failed", "write_request", "write":
			return true
		default:
			return false
	placeholder
	default:
		return true
placeholder
placeholder

func dropPreviousResponseIDFromRawPayload(payload []byte) ([]byte, bool, error) {
	return dropPreviousResponseIDFromRawPayloadWithDeleteFn(payload, sjson.DeleteBytes)
placeholder

func dropPreviousResponseIDFromRawPayloadWithDeleteFn(
	payload []byte,
	deleteFn func([]byte, string) ([]byte, error),
) ([]byte, bool, error) {
	if len(payload) == 0 {
		return payload, false, nil
placeholder
	if !gjson.GetBytes(payload, "previous_response_id").Exists() {
		return payload, false, nil
placeholder
	if deleteFn == nil {
		deleteFn = sjson.DeleteBytes
placeholder

	updated := payload
	for i := 0; i < openAIWSMaxPrevResponseIDDeletePasses &&
		gjson.GetBytes(updated, "previous_response_id").Exists(); i++ {
		next, err := deleteFn(updated, "previous_response_id")
		if err != nil {
			return payload, false, err
	placeholder
		updated = next
placeholder
	return updated, !gjson.GetBytes(updated, "previous_response_id").Exists(), nil
placeholder

func setPreviousResponseIDToRawPayload(payload []byte, previousResponseID string) ([]byte, error) {
	normalizedPrevID := strings.TrimSpace(previousResponseID)
	if len(payload) == 0 || normalizedPrevID == "" {
		return payload, nil
placeholder
	updated, err := sjson.SetBytes(payload, "previous_response_id", normalizedPrevID)
	if err == nil {
		return updated, nil
placeholder

	var reqBody map[string]any
	if unmarshalErr := json.Unmarshal(payload, &reqBody); unmarshalErr != nil {
		return nil, err
placeholder
	reqBody["previous_response_id"] = normalizedPrevID
	rebuilt, marshalErr := json.Marshal(reqBody)
	if marshalErr != nil {
		return nil, marshalErr
placeholder
	return rebuilt, nil
placeholder

func shouldInferIngressFunctionCallOutputPreviousResponseID(
	storeDisabled bool,
	turn int,
	signals ToolContinuationSignals,
	currentPreviousResponseID string,
	expectedPreviousResponseID string,
) bool {
	if !storeDisabled || turn <= 1 || !signals.HasFunctionCallOutput {
		return false
placeholder
	if strings.TrimSpace(currentPreviousResponseID) != "" {
		return false
placeholder
	if signals.HasFunctionCallOutputMissingCallID {
		return false
placeholder
	// If the client already sent the actual tool-call context, treat this as
	// a full replay / self-contained continuation payload rather than
	// downgrading it into an inferred delta continuation. item_reference alone
	// is not enough on the store=false WS path: it still needs a valid prior
	// response anchor so upstream can resolve the referenced function_call.
	if signals.HasToolCallContext {
		return false
placeholder
	return strings.TrimSpace(expectedPreviousResponseID) != ""
placeholder

func alignStoreDisabledPreviousResponseID(
	payload []byte,
	expectedPreviousResponseID string,
) ([]byte, bool, error) {
	if len(payload) == 0 {
		return payload, false, nil
placeholder
	expected := strings.TrimSpace(expectedPreviousResponseID)
	if expected == "" {
		return payload, false, nil
placeholder
	current := openAIWSPayloadStringFromRaw(payload, "previous_response_id")
	if current == "" || current == expected {
		return payload, false, nil
placeholder

	withoutPrev, removed, dropErr := dropPreviousResponseIDFromRawPayload(payload)
	if dropErr != nil {
		return payload, false, dropErr
placeholder
	if !removed {
		return payload, false, nil
placeholder
	updated, setErr := setPreviousResponseIDToRawPayload(withoutPrev, expected)
	if setErr != nil {
		return payload, false, setErr
placeholder
	return updated, true, nil
placeholder

func cloneOpenAIWSPayloadBytes(payload []byte) []byte {
	if len(payload) == 0 {
		return nil
placeholder
	cloned := make([]byte, len(payload))
	copy(cloned, payload)
	return cloned
placeholder

func cloneOpenAIWSRawMessages(items []json.RawMessage) []json.RawMessage {
	if items == nil {
		return nil
placeholder
	cloned := make([]json.RawMessage, 0, len(items))
	for idx := range items {
		cloned = append(cloned, json.RawMessage(cloneOpenAIWSPayloadBytes(items[idx])))
placeholder
	return cloned
placeholder

func normalizeOpenAIWSJSONForCompare(raw []byte) ([]byte, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, errors.New("json is empty")
placeholder
	var decoded any
	if err := json.Unmarshal(trimmed, &decoded); err != nil {
		return nil, err
placeholder
	return json.Marshal(decoded)
placeholder

func normalizeOpenAIWSJSONForCompareOrRaw(raw []byte) []byte {
	normalized, err := normalizeOpenAIWSJSONForCompare(raw)
	if err != nil {
		return bytes.TrimSpace(raw)
placeholder
	return normalized
placeholder

func normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(payload []byte) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errors.New("payload is empty")
placeholder
	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		return nil, err
placeholder
	delete(decoded, "input")
	delete(decoded, "previous_response_id")
	// Codex changes transport-only metadata for every response.create. These fields
	// do not alter the context referenced by previous_response_id and are excluded
	// from Codex's own websocket reuse comparison.
	delete(decoded, "client_metadata")
	delete(decoded, "stream_options")
	// Official Codex prewarms a connection with generate=false, then omits the
	// field on the business request that continues from the prewarm response.
	// Only normalize false so a meaningful generate=true change remains visible.
	if generate, ok := decoded["generate"].(bool); ok && !generate {
		delete(decoded, "generate")
placeholder
	return json.Marshal(decoded)
placeholder

func openAIWSExtractNormalizedInputSequence(payload []byte) ([]json.RawMessage, bool, error) {
	if len(payload) == 0 {
		return nil, false, nil
placeholder
	inputValue := gjson.GetBytes(payload, "input")
	if !inputValue.Exists() {
		return nil, false, nil
placeholder
	if inputValue.Type == gjson.JSON {
		raw := strings.TrimSpace(inputValue.Raw)
		if strings.HasPrefix(raw, "[") {
			var items []json.RawMessage
			if err := json.Unmarshal([]byte(raw), &items); err != nil {
				return nil, true, err
		placeholder
			return items, true, nil
	placeholder
		return []json.RawMessage{json.RawMessage(raw)placeholder, true, nil
placeholder
	if inputValue.Type == gjson.String {
		encoded, _ := json.Marshal(inputValue.String())
		return []json.RawMessage{encodedplaceholder, true, nil
placeholder
	return []json.RawMessage{json.RawMessage(inputValue.Raw)placeholder, true, nil
placeholder

func openAIWSInputIsPrefixExtended(previousPayload, currentPayload []byte) (bool, error) {
	previousItems, previousExists, prevErr := openAIWSExtractNormalizedInputSequence(previousPayload)
	if prevErr != nil {
		return false, prevErr
placeholder
	currentItems, currentExists, currentErr := openAIWSExtractNormalizedInputSequence(currentPayload)
	if currentErr != nil {
		return false, currentErr
placeholder
	if !previousExists && !currentExists {
		return true, nil
placeholder
	if !previousExists {
		return len(currentItems) == 0, nil
placeholder
	if !currentExists {
		return len(previousItems) == 0, nil
placeholder
	if len(currentItems) < len(previousItems) {
		return false, nil
placeholder

	for idx := range previousItems {
		previousNormalized := normalizeOpenAIWSJSONForCompareOrRaw(previousItems[idx])
		currentNormalized := normalizeOpenAIWSJSONForCompareOrRaw(currentItems[idx])
		if !bytes.Equal(previousNormalized, currentNormalized) {
			return false, nil
	placeholder
placeholder
	return true, nil
placeholder

func openAIWSRawItemsHasPrefix(items []json.RawMessage, prefix []json.RawMessage) bool {
	if len(prefix) == 0 {
		return true
placeholder
	if len(items) < len(prefix) {
		return false
placeholder
	for idx := range prefix {
		previousNormalized := normalizeOpenAIWSJSONForCompareOrRaw(prefix[idx])
		currentNormalized := normalizeOpenAIWSJSONForCompareOrRaw(items[idx])
		if !bytes.Equal(previousNormalized, currentNormalized) {
			return false
	placeholder
placeholder
	return true
placeholder

func openAIWSRawItemsHasFunctionCallOutput(items []json.RawMessage) bool {
	for _, item := range items {
		if isCodexToolCallOutputItemType(gjson.GetBytes(item, "type").String()) {
			return true
	placeholder
placeholder
	return false
placeholder

func openAIWSRawItemsHaveToolCallContextForOutputs(items []json.RawMessage) bool {
	if len(items) == 0 {
		return false
placeholder
	contextCallIDs := make(map[string]struct{placeholder)
	outputCallIDs := make(map[string]struct{placeholder)
	for _, item := range items {
		itemType := gjson.GetBytes(item, "type").String()
		callID := strings.TrimSpace(gjson.GetBytes(item, "call_id").String())
		switch {
		case isCodexToolCallContextItemType(itemType):
			if callID != "" {
				contextCallIDs[callID] = struct{placeholder{placeholder
		placeholder
		case isCodexToolCallOutputItemType(itemType):
			if callID == "" {
				return false
		placeholder
			outputCallIDs[callID] = struct{placeholder{placeholder
	placeholder
placeholder
	if len(outputCallIDs) == 0 || len(contextCallIDs) == 0 {
		return false
placeholder
	for callID := range outputCallIDs {
		if _, ok := contextCallIDs[callID]; !ok {
			return false
	placeholder
placeholder
	return true
placeholder

func openAIWSRawPayloadHasToolCallOutput(payload []byte) bool {
	if len(payload) == 0 {
		return false
placeholder
	input := gjson.GetBytes(payload, "input")
	if !input.Exists() {
		return false
placeholder
	if input.IsArray() {
		for _, item := range input.Array() {
			if isCodexToolCallOutputItemType(item.Get("type").String()) {
				return true
		placeholder
	placeholder
		return false
placeholder
	if input.Type == gjson.JSON {
		return isCodexToolCallOutputItemType(input.Get("type").String())
placeholder
	return false
placeholder

func buildOpenAIWSReplayInputSequence(
	previousFullInput []json.RawMessage,
	previousFullInputExists bool,
	currentPayload []byte,
	hasPreviousResponseID bool,
) ([]json.RawMessage, bool, error) {
	currentItems, currentExists, currentErr := openAIWSExtractNormalizedInputSequence(currentPayload)
	if currentErr != nil {
		return nil, false, currentErr
placeholder
	if !hasPreviousResponseID {
		return cloneOpenAIWSRawMessages(currentItems), currentExists, nil
placeholder
	if !previousFullInputExists {
		return cloneOpenAIWSRawMessages(currentItems), currentExists, nil
placeholder
	if !currentExists || len(currentItems) == 0 {
		return cloneOpenAIWSRawMessages(previousFullInput), true, nil
placeholder
	if openAIWSRawItemsHasPrefix(currentItems, previousFullInput) {
		return cloneOpenAIWSRawMessages(currentItems), true, nil
placeholder
	merged := make([]json.RawMessage, 0, len(previousFullInput)+len(currentItems))
	merged = append(merged, cloneOpenAIWSRawMessages(previousFullInput)...)
	merged = append(merged, cloneOpenAIWSRawMessages(currentItems)...)
	return merged, true, nil
placeholder

func setOpenAIWSPayloadInputSequence(
	payload []byte,
	fullInput []json.RawMessage,
	fullInputExists bool,
) ([]byte, error) {
	if !fullInputExists {
		return payload, nil
placeholder
	// Preserve [] vs null semantics when input exists but is empty.
	inputForMarshal := fullInput
	if inputForMarshal == nil {
		inputForMarshal = []json.RawMessage{placeholder
placeholder
	inputRaw, marshalErr := json.Marshal(inputForMarshal)
	if marshalErr != nil {
		return nil, marshalErr
placeholder
	return sjson.SetRawBytes(payload, "input", inputRaw)
placeholder

func shouldKeepIngressPreviousResponseID(
	previousPayload []byte,
	currentPayload []byte,
	lastTurnResponseID string,
	hasFunctionCallOutput bool,
) (bool, string, error) {
	if hasFunctionCallOutput {
		return true, "has_function_call_output", nil
placeholder
	currentPreviousResponseID := strings.TrimSpace(openAIWSPayloadStringFromRaw(currentPayload, "previous_response_id"))
	if currentPreviousResponseID == "" {
		return false, "missing_previous_response_id", nil
placeholder
	expectedPreviousResponseID := strings.TrimSpace(lastTurnResponseID)
	if expectedPreviousResponseID == "" {
		return false, "missing_last_turn_response_id", nil
placeholder
	if currentPreviousResponseID != expectedPreviousResponseID {
		return false, "previous_response_id_mismatch", nil
placeholder
	if len(previousPayload) == 0 {
		return false, "missing_previous_turn_payload", nil
placeholder

	previousComparable, previousComparableErr := normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(previousPayload)
	if previousComparableErr != nil {
		return false, "non_input_compare_error", previousComparableErr
placeholder
	currentComparable, currentComparableErr := normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(currentPayload)
	if currentComparableErr != nil {
		return false, "non_input_compare_error", currentComparableErr
placeholder
	if !bytes.Equal(previousComparable, currentComparable) {
		return false, "non_input_changed", nil
placeholder
	return true, "strict_incremental_ok", nil
placeholder

type openAIWSIngressPreviousTurnStrictState struct {
	nonInputComparable []byte
placeholder

func buildOpenAIWSIngressPreviousTurnStrictState(payload []byte) (*openAIWSIngressPreviousTurnStrictState, error) {
	if len(payload) == 0 {
		return nil, nil
placeholder
	nonInputComparable, nonInputErr := normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(payload)
	if nonInputErr != nil {
		return nil, nonInputErr
placeholder
	return &openAIWSIngressPreviousTurnStrictState{
		nonInputComparable: nonInputComparable,
placeholder, nil
placeholder

func shouldKeepIngressPreviousResponseIDWithStrictState(
	previousState *openAIWSIngressPreviousTurnStrictState,
	currentPayload []byte,
	lastTurnResponseID string,
	hasFunctionCallOutput bool,
) (bool, string, error) {
	if hasFunctionCallOutput {
		return true, "has_function_call_output", nil
placeholder
	currentPreviousResponseID := strings.TrimSpace(openAIWSPayloadStringFromRaw(currentPayload, "previous_response_id"))
	if currentPreviousResponseID == "" {
		return false, "missing_previous_response_id", nil
placeholder
	expectedPreviousResponseID := strings.TrimSpace(lastTurnResponseID)
	if expectedPreviousResponseID == "" {
		return false, "missing_last_turn_response_id", nil
placeholder
	if currentPreviousResponseID != expectedPreviousResponseID {
		return false, "previous_response_id_mismatch", nil
placeholder
	if previousState == nil {
		return false, "missing_previous_turn_payload", nil
placeholder

	currentComparable, currentComparableErr := normalizeOpenAIWSPayloadWithoutInputAndPreviousResponseID(currentPayload)
	if currentComparableErr != nil {
		return false, "non_input_compare_error", currentComparableErr
placeholder
	if !bytes.Equal(previousState.nonInputComparable, currentComparable) {
		return false, "non_input_changed", nil
placeholder
	return true, "strict_incremental_ok", nil
placeholder
