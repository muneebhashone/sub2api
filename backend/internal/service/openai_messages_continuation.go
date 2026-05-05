package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

type openAICompatSessionResponseBinding struct {
	ResponseID           string
	TurnState            string
	ContinuationDisabled bool
	ExpiresAt            time.Time
placeholder

func openAICompatContinuationEnabled(account *Account, model string) bool {
	if account == nil || account.Type != AccountTypeAPIKey {
		return false
placeholder
	return shouldAutoInjectPromptCacheKeyForCompat(model)
placeholder

func trimAnthropicCompatResponsesInputToLatestTurn(req *apicompat.ResponsesRequest) {
	if req == nil || len(req.Input) == 0 {
		return
placeholder

	var items []apicompat.ResponsesInputItem
	if err := json.Unmarshal(req.Input, &items); err != nil || len(items) == 0 {
		return
placeholder

	start := len(items) - 1
	for start > 0 && items[start].Type == "function_call_output" {
		start--
placeholder
	trimmed := append([]apicompat.ResponsesInputItem(nil), items[start:]...)
	if len(trimmed) == len(items) {
		return
placeholder
	if input, err := json.Marshal(trimmed); err == nil {
		req.Input = input
placeholder
placeholder

func isOpenAICompatPreviousResponseNotFound(statusCode int, upstreamMsg string, upstreamBody []byte) bool {
	if statusCode != http.StatusBadRequest && statusCode != http.StatusNotFound {
		return false
placeholder
	check := func(s string) bool {
		lower := strings.ToLower(strings.TrimSpace(s))
		return strings.Contains(lower, "previous_response_not_found") ||
			(strings.Contains(lower, "previous response") && strings.Contains(lower, "not found")) ||
			(strings.Contains(lower, "unsupported parameter") && strings.Contains(lower, "previous_response_id"))
placeholder
	if check(upstreamMsg) || check(string(upstreamBody)) {
		return true
placeholder
	return check(gjson.GetBytes(upstreamBody, "error.code").String()) ||
		check(gjson.GetBytes(upstreamBody, "error.message").String())
placeholder

func isOpenAICompatPreviousResponseUnsupported(statusCode int, upstreamMsg string, upstreamBody []byte) bool {
	if statusCode != http.StatusBadRequest {
		return false
placeholder
	check := func(s string) bool {
		lower := strings.ToLower(strings.TrimSpace(s))
		if !strings.Contains(lower, "previous_response_id") {
			return false
	placeholder
		return strings.Contains(lower, "unsupported parameter") ||
			strings.Contains(lower, "only supported on responses websocket") ||
			strings.Contains(lower, "not supported")
placeholder
	if check(upstreamMsg) || check(string(upstreamBody)) {
		return true
placeholder
	return check(gjson.GetBytes(upstreamBody, "error.code").String()) ||
		check(gjson.GetBytes(upstreamBody, "error.message").String())
placeholder

func openAICompatSessionResponseKey(c *gin.Context, account *Account, promptCacheKey string) string {
	key := strings.TrimSpace(promptCacheKey)
	if account == nil || key == "" {
		return ""
placeholder
	apiKeyID := int64(0)
	if c != nil {
		apiKeyID = getAPIKeyIDFromContext(c)
placeholder
	return strings.Join([]string{
		strconv.FormatInt(account.ID, 10),
		strconv.FormatInt(apiKeyID, 10),
		key,
placeholder, "\x00")
placeholder

func (s *OpenAIGatewayService) getOpenAICompatSessionResponseID(_ context.Context, c *gin.Context, account *Account, promptCacheKey string) string {
	if s == nil {
		return ""
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	if key == "" {
		return ""
placeholder
	raw, ok := s.openaiCompatSessionResponses.Load(key)
	if !ok {
		return ""
placeholder
	binding, ok := raw.(openAICompatSessionResponseBinding)
	if !ok {
		s.openaiCompatSessionResponses.Delete(key)
		return ""
placeholder
	if !binding.ExpiresAt.IsZero() && time.Now().After(binding.ExpiresAt) {
		s.openaiCompatSessionResponses.Delete(key)
		return ""
placeholder
	if binding.ContinuationDisabled {
		return ""
placeholder
	if strings.TrimSpace(binding.ResponseID) == "" {
		s.openaiCompatSessionResponses.Delete(key)
		return ""
placeholder
	return strings.TrimSpace(binding.ResponseID)
placeholder

func (s *OpenAIGatewayService) bindOpenAICompatSessionResponseID(_ context.Context, c *gin.Context, account *Account, promptCacheKey, responseID string) {
	if s == nil {
		return
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	id := strings.TrimSpace(responseID)
	if key == "" || id == "" {
		return
placeholder
	binding := openAICompatSessionResponseBinding{
		ResponseID: id,
		ExpiresAt:  time.Now().Add(s.openAIWSResponseStickyTTL()),
placeholder
	if raw, ok := s.openaiCompatSessionResponses.Load(key); ok {
		if existing, ok := raw.(openAICompatSessionResponseBinding); ok {
			if existing.ContinuationDisabled {
				existing.ResponseID = ""
				existing.ExpiresAt = time.Now().Add(s.openAIWSResponseStickyTTL())
				s.openaiCompatSessionResponses.Store(key, existing)
				return
		placeholder
			binding.TurnState = existing.TurnState
	placeholder
placeholder
	s.openaiCompatSessionResponses.Store(key, binding)
placeholder

func (s *OpenAIGatewayService) deleteOpenAICompatSessionResponseID(_ context.Context, c *gin.Context, account *Account, promptCacheKey string) {
	if s == nil {
		return
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	if key == "" {
		return
placeholder
	raw, ok := s.openaiCompatSessionResponses.Load(key)
	if !ok {
		return
placeholder
	binding, ok := raw.(openAICompatSessionResponseBinding)
	if !ok {
		s.openaiCompatSessionResponses.Delete(key)
		return
placeholder
	binding.ResponseID = ""
	if strings.TrimSpace(binding.TurnState) == "" && !binding.ContinuationDisabled {
		s.openaiCompatSessionResponses.Delete(key)
		return
placeholder
	binding.ExpiresAt = time.Now().Add(s.openAIWSResponseStickyTTL())
	s.openaiCompatSessionResponses.Store(key, binding)
placeholder

func (s *OpenAIGatewayService) disableOpenAICompatSessionContinuation(_ context.Context, c *gin.Context, account *Account, promptCacheKey string) {
	if s == nil {
		return
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	if key == "" {
		return
placeholder
	binding := openAICompatSessionResponseBinding{
		ContinuationDisabled: true,
		ExpiresAt:            time.Now().Add(s.openAIWSResponseStickyTTL()),
placeholder
	if raw, ok := s.openaiCompatSessionResponses.Load(key); ok {
		if existing, ok := raw.(openAICompatSessionResponseBinding); ok {
			binding.TurnState = existing.TurnState
	placeholder
placeholder
	s.openaiCompatSessionResponses.Store(key, binding)
placeholder

func (s *OpenAIGatewayService) isOpenAICompatSessionContinuationDisabled(_ context.Context, c *gin.Context, account *Account, promptCacheKey string) bool {
	if s == nil {
		return false
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	if key == "" {
		return false
placeholder
	raw, ok := s.openaiCompatSessionResponses.Load(key)
	if !ok {
		return false
placeholder
	binding, ok := raw.(openAICompatSessionResponseBinding)
	if !ok {
		s.openaiCompatSessionResponses.Delete(key)
		return false
placeholder
	if !binding.ExpiresAt.IsZero() && time.Now().After(binding.ExpiresAt) {
		s.openaiCompatSessionResponses.Delete(key)
		return false
placeholder
	return binding.ContinuationDisabled
placeholder

func (s *OpenAIGatewayService) getOpenAICompatSessionTurnState(_ context.Context, c *gin.Context, account *Account, promptCacheKey string) string {
	if s == nil {
		return ""
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	if key == "" {
		return ""
placeholder
	raw, ok := s.openaiCompatSessionResponses.Load(key)
	if !ok {
		return ""
placeholder
	binding, ok := raw.(openAICompatSessionResponseBinding)
	if !ok || strings.TrimSpace(binding.TurnState) == "" {
		return ""
placeholder
	if !binding.ExpiresAt.IsZero() && time.Now().After(binding.ExpiresAt) {
		s.openaiCompatSessionResponses.Delete(key)
		return ""
placeholder
	return strings.TrimSpace(binding.TurnState)
placeholder

func (s *OpenAIGatewayService) bindOpenAICompatSessionTurnState(_ context.Context, c *gin.Context, account *Account, promptCacheKey, turnState string) {
	if s == nil {
		return
placeholder
	key := openAICompatSessionResponseKey(c, account, promptCacheKey)
	state := strings.TrimSpace(turnState)
	if key == "" || state == "" {
		return
placeholder
	binding := openAICompatSessionResponseBinding{
		TurnState: state,
		ExpiresAt: time.Now().Add(s.openAIWSResponseStickyTTL()),
placeholder
	if raw, ok := s.openaiCompatSessionResponses.Load(key); ok {
		if existing, ok := raw.(openAICompatSessionResponseBinding); ok {
			binding.ResponseID = existing.ResponseID
			binding.ContinuationDisabled = existing.ContinuationDisabled
	placeholder
placeholder
	s.openaiCompatSessionResponses.Store(key, binding)
placeholder
