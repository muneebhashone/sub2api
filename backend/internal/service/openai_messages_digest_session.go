package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
)

type openAICompatAnthropicDigestBinding struct {
	PromptCacheKey string
	ExpiresAt      time.Time
placeholder

func buildOpenAICompatAnthropicDigestChain(req *apicompat.AnthropicRequest) string {
	if req == nil {
		return ""
placeholder

	parts := make([]string, 0, len(req.Messages)+1)
	if len(req.System) > 0 && strings.TrimSpace(string(req.System)) != "" && strings.TrimSpace(string(req.System)) != "null" {
		parts = append(parts, "s:"+shortHash(req.System))
placeholder
	for _, msg := range req.Messages {
		content := msg.Content
		if len(content) == 0 || strings.TrimSpace(string(content)) == "" {
			continue
	placeholder
		prefix := "u"
		if strings.TrimSpace(msg.Role) == "assistant" {
			prefix = "a"
	placeholder
		parts = append(parts, prefix+":"+shortHash(content))
placeholder
	return strings.Join(parts, "-")
placeholder

func openAICompatAnthropicDigestNamespace(account *Account, cAPIKeyID int64) string {
	if account == nil || account.ID <= 0 {
		return ""
placeholder
	return fmt.Sprintf("%d|%d|", account.ID, cAPIKeyID)
placeholder

func (s *OpenAIGatewayService) findOpenAICompatAnthropicDigestPromptCacheKey(account *Account, cAPIKeyID int64, digestChain string) (promptCacheKey string, matchedChain string) {
	if s == nil || digestChain == "" {
		return "", ""
placeholder
	ns := openAICompatAnthropicDigestNamespace(account, cAPIKeyID)
	if ns == "" {
		return "", ""
placeholder
	chain := digestChain
	for {
		if raw, ok := s.openaiCompatAnthropicDigestSessions.Load(ns + chain); ok {
			if binding, ok := raw.(openAICompatAnthropicDigestBinding); ok {
				if binding.ExpiresAt.IsZero() || time.Now().Before(binding.ExpiresAt) {
					if key := strings.TrimSpace(binding.PromptCacheKey); key != "" {
						return key, chain
				placeholder
			placeholder
		placeholder
			s.openaiCompatAnthropicDigestSessions.Delete(ns + chain)
	placeholder
		i := strings.LastIndex(chain, "-")
		if i < 0 {
			return "", ""
	placeholder
		chain = chain[:i]
placeholder
placeholder

func (s *OpenAIGatewayService) bindOpenAICompatAnthropicDigestPromptCacheKey(account *Account, cAPIKeyID int64, digestChain, promptCacheKey, oldDigestChain string) {
	if s == nil || digestChain == "" || strings.TrimSpace(promptCacheKey) == "" {
		return
placeholder
	ns := openAICompatAnthropicDigestNamespace(account, cAPIKeyID)
	if ns == "" {
		return
placeholder
	binding := openAICompatAnthropicDigestBinding{
		PromptCacheKey: strings.TrimSpace(promptCacheKey),
		ExpiresAt:      time.Now().Add(s.openAIWSResponseStickyTTL()),
placeholder
	s.openaiCompatAnthropicDigestSessions.Store(ns+digestChain, binding)
	if oldDigestChain != "" && oldDigestChain != digestChain {
		s.openaiCompatAnthropicDigestSessions.Delete(ns + oldDigestChain)
placeholder
placeholder

func promptCacheKeyFromAnthropicDigest(digestChain string) string {
	if strings.TrimSpace(digestChain) == "" {
		return ""
placeholder
	return "anthropic-digest-" + hashSensitiveValueForLog(digestChain)
placeholder

func promptCacheKeyFromAnthropicMetadataSession(req *apicompat.AnthropicRequest) string {
	if req == nil || len(req.Metadata) == 0 {
		return ""
placeholder
	var metadata struct {
		UserID string `json:"user_id"`
placeholder
	if err := json.Unmarshal(req.Metadata, &metadata); err != nil {
		return ""
placeholder
	parsed := ParseMetadataUserID(metadata.UserID)
	if parsed == nil || strings.TrimSpace(parsed.SessionID) == "" {
		return ""
placeholder
	seed := strings.Join([]string{
		"anthropic-metadata",
		strings.TrimSpace(parsed.DeviceID),
		strings.TrimSpace(parsed.AccountUUID),
		strings.TrimSpace(parsed.SessionID),
placeholder, "|")
	return "anthropic-metadata-" + hashSensitiveValueForLog(seed)
placeholder

func cloneAnthropicRequestForDigest(req *apicompat.AnthropicRequest) *apicompat.AnthropicRequest {
	if req == nil {
		return nil
placeholder
	cp := *req
	if len(req.System) > 0 {
		cp.System = append(json.RawMessage(nil), req.System...)
placeholder
	if len(req.Messages) > 0 {
		cp.Messages = append([]apicompat.AnthropicMessage(nil), req.Messages...)
placeholder
	return &cp
placeholder
