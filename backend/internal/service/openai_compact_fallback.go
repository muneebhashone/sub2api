package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type openAICompactFallbackSignal struct {
	payload []byte
	message string
placeholder

func (e *openAICompactFallbackSignal) Error() string {
	if e == nil || strings.TrimSpace(e.message) == "" {
		return "upstream compact request failed"
placeholder
	return e.message
placeholder

func asOpenAICompactFallbackSignal(err error) (*openAICompactFallbackSignal, bool) {
	var signal *openAICompactFallbackSignal
	return signal, errors.As(err, &signal) && signal != nil
placeholder

func isExplicitOpenAICompactContext(c *gin.Context) bool {
	return isOpenAIResponsesCompactPath(c) || isOpenAINativeCompactionV2(c)
placeholder

func newOpenAICompactFallbackSignal(c *gin.Context, payload []byte, message string) error {
	if !isExplicitOpenAICompactContext(c) ||
		!isOpenAICompactModelFailure(http.StatusBadRequest, message, payload) {
		return nil
placeholder
	return &openAICompactFallbackSignal{
		payload: append([]byte(nil), payload...),
		message: sanitizeUpstreamErrorMessage(strings.TrimSpace(message)),
placeholder
placeholder

func isExplicitOpenAICompactRequest(c *gin.Context, body []byte) bool {
	return isOpenAIResponsesCompactPath(c) || HasCompactionTriggerInInput(body)
placeholder

// resolveOpenAICompactFallbackModel prefers the account's compact-only rule
// for the client-visible model. The process-wide fallback is used only when
// that account has no matching compact rule.
func (s *OpenAIGatewayService) resolveOpenAICompactFallbackModel(account *Account, requestedModel string) string {
	requestedModel = strings.TrimSpace(requestedModel)
	if account != nil {
		if mapped, matched := account.ResolveCompactMappedModel(requestedModel); matched {
			if mapped = strings.TrimSpace(mapped); mapped != "" {
				return mapped
		placeholder
	placeholder
placeholder
	if s == nil || s.cfg == nil {
		return ""
placeholder
	fallback := strings.TrimSpace(s.cfg.Gateway.OpenAICompactModel)
	if fallback == "" {
		return ""
placeholder
	return strings.TrimSpace(resolveOpenAIAccountUpstreamModelForRequest(account, fallback, false))
placeholder

func isOpenAICompactModelFailure(statusCode int, upstreamMsg string, upstreamBody []byte) bool {
	if isOpenAIContextWindowError(upstreamMsg, upstreamBody) {
		return true
placeholder
	if statusCode != http.StatusBadRequest && statusCode != http.StatusNotFound {
		return false
placeholder

	values := []string{
		extractUpstreamErrorCode(upstreamBody),
		upstreamMsg,
		gjson.GetBytes(upstreamBody, "error.type").String(),
		gjson.GetBytes(upstreamBody, "response.error.code").String(),
		gjson.GetBytes(upstreamBody, "response.error.type").String(),
placeholder
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		switch value {
		case "model_not_found", "model_not_available", "unsupported_model", "invalid_model":
			return true
	placeholder
		if isExplicitOpenAIModelAvailabilityMessage(value) {
			return true
	placeholder
placeholder
	// Some compact providers return only a failed response shell. It is safe to
	// retry that shape for an explicit compact request, but a populated error is
	// left untouched so business and policy failures keep their original wire.
	if strings.EqualFold(strings.TrimSpace(gjson.GetBytes(upstreamBody, "response.status").String()), "failed") ||
		strings.EqualFold(strings.TrimSpace(gjson.GetBytes(upstreamBody, "status").String()), "failed") {
		for _, path := range []string{
			"error.message", "error.code", "error.type",
			"response.error.message", "response.error.code", "response.error.type",
	placeholder {
			if strings.TrimSpace(gjson.GetBytes(upstreamBody, path).String()) != "" {
				return false
		placeholder
	placeholder
		return strings.TrimSpace(upstreamMsg) == ""
placeholder
	return false
placeholder

func isExplicitOpenAIModelAvailabilityMessage(value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return false
placeholder
	for _, phrase := range []string{
		"model not found",
		"model does not exist",
		"model is unavailable",
		"model is not available",
		"model is unsupported",
		"model is not supported",
		"unsupported model",
placeholder {
		if strings.Contains(value, phrase) {
			return true
	placeholder
placeholder
	// OpenAI commonly identifies the missing model between the word "model"
	// and the terminal availability phrase, for example: "The model `x` does
	// not exist". Requiring the message to start with the model subject avoids
	// treating unrelated feature errors such as "model output is not supported"
	// as a signal to change models.
	if strings.HasPrefix(value, "the model ") || strings.HasPrefix(value, "model ") {
		return strings.Contains(value, " does not exist") ||
			strings.Contains(value, " was not found") ||
			strings.Contains(value, " is unavailable") ||
			strings.Contains(value, " is not available")
placeholder
	return false
placeholder

func openAICompactFallbackErrorResponse(resp *http.Response, signal *openAICompactFallbackSignal) (*http.Response, []byte) {
	headers := make(http.Header)
	if resp != nil {
		headers = resp.Header.Clone()
placeholder
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/json")
placeholder
	payload := normalizeOpenAICompactFallbackHTTPErrorPayload(signal)
	return &http.Response{
		StatusCode: http.StatusBadRequest,
		Header:     headers,
		Body:       io.NopCloser(bytes.NewReader(payload)),
placeholder, payload
placeholder

func normalizeOpenAICompactFallbackHTTPErrorPayload(signal *openAICompactFallbackSignal) []byte {
	if signal == nil {
		return nil
placeholder
	payload := append([]byte(nil), signal.payload...)
	var terminal struct {
		Error    json.RawMessage `json:"error"`
		Response struct {
			Error json.RawMessage `json:"error"`
	placeholder `json:"response"`
placeholder
	if json.Unmarshal(payload, &terminal) != nil || len(bytes.TrimSpace(terminal.Response.Error)) == 0 ||
		bytes.Equal(bytes.TrimSpace(terminal.Response.Error), []byte("null")) {
		return payload
placeholder
	// Standard HTTP error handlers consume error.message/type/code. A streamed
	// response.failed terminal nests the same object under response.error, so
	// normalize only that envelope at the stream-to-HTTP boundary.
	normalized, err := json.Marshal(struct {
		Error json.RawMessage `json:"error"`
placeholder{Error: terminal.Response.Errorplaceholder)
	if err != nil {
		return payload
placeholder
	return normalized
placeholder

func (s *OpenAIGatewayService) appendOpenAICompactFallbackRetryOps(
	c *gin.Context,
	account *Account,
	resp *http.Response,
	payload []byte,
	message string,
	passthrough bool,
) {
	if account == nil {
		return
placeholder
	statusCode := http.StatusBadRequest
	requestID := ""
	if resp != nil {
		statusCode = resp.StatusCode
		requestID = resp.Header.Get("x-request-id")
placeholder
	detail := ""
	if s != nil && s.cfg != nil && s.cfg.Gateway.LogUpstreamErrorBody {
		maxBytes := s.cfg.Gateway.LogUpstreamErrorBodyMaxBytes
		if maxBytes <= 0 {
			maxBytes = 2048
	placeholder
		detail = truncateString(string(payload), maxBytes)
placeholder
	appendOpsUpstreamError(c, OpsUpstreamErrorEvent{
		Platform:             account.Platform,
		AccountID:            account.ID,
		AccountName:          account.Name,
		UpstreamStatusCode:   statusCode,
		UpstreamRequestID:    requestID,
		Passthrough:          passthrough,
		Kind:                 "retry",
		Reason:               "compact_model_fallback",
		Message:              sanitizeUpstreamErrorMessage(strings.TrimSpace(message)),
		Detail:               detail,
		UpstreamResponseBody: detail,
placeholder)
placeholder

// prepareOpenAICompactFallbackRetry returns a body for one safe, same-account
// retry. Callers invoke it only before any downstream response has been
// written; it changes the model and deliberately leaves path, trigger, and
// native-v2 context state untouched.
func (s *OpenAIGatewayService) prepareOpenAICompactFallbackRetry(
	c *gin.Context,
	account *Account,
	requestedModel string,
	currentBody []byte,
	statusCode int,
	upstreamMsg string,
	upstreamBody []byte,
	alreadyRetried bool,
) ([]byte, string, bool) {
	if alreadyRetried || !isExplicitOpenAICompactRequest(c, currentBody) ||
		!isOpenAICompactModelFailure(statusCode, upstreamMsg, upstreamBody) {
		return currentBody, "", false
placeholder
	fallbackModel := s.resolveOpenAICompactFallbackModel(account, requestedModel)
	currentModel := strings.TrimSpace(gjson.GetBytes(currentBody, "model").String())
	if fallbackModel == "" || strings.EqualFold(fallbackModel, currentModel) {
		return currentBody, "", false
placeholder
	retryBody := ReplaceModelInBody(currentBody, fallbackModel)
	if strings.EqualFold(strings.TrimSpace(gjson.GetBytes(retryBody, "model").String()), currentModel) {
		return currentBody, "", false
placeholder
	return retryBody, fallbackModel, true
placeholder

func (s *OpenAIGatewayService) applyOpenAIPassthroughCompactFallbackFromSignal(
	c *gin.Context,
	account *Account,
	requestedModel string,
	body []byte,
	err error,
	alreadyRetried bool,
	resp *http.Response,
) ([]byte, string, bool) {
	signal, ok := asOpenAICompactFallbackSignal(err)
	if !ok {
		return body, "", false
placeholder
	retryBody, fallbackModel, retry := s.prepareOpenAICompactFallbackRetry(
		c, account, requestedModel, body, http.StatusBadRequest, signal.message, signal.payload, alreadyRetried,
	)
	if !retry {
		return body, "", false
placeholder
	s.appendOpenAICompactFallbackRetryOps(c, account, resp, signal.payload, signal.message, true)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
placeholder
	fromModel := strings.TrimSpace(gjson.GetBytes(body, "model").String())
	accountName := ""
	if account != nil {
		accountName = account.Name
placeholder
	SetOpsUpstreamModel(c, fallbackModel)
	logger.LegacyPrintf(
		"service.openai_gateway",
		"[OpenAI passthrough] Retrying explicit compact request once with fallback model (account: %s, from: %s, to: %s, upstream_code: %s)",
		accountName, fromModel, fallbackModel, extractUpstreamErrorCode(signal.payload),
	)
	return retryBody, fallbackModel, true
placeholder
