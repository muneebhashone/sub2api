package service

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/sjson"
)

type claudeMaxResponseRewriteContext struct {
	Parsed *ParsedRequest
	Group  *Group
placeholder

type claudeMaxResponseRewriteContextKeyType struct{placeholder

var claudeMaxResponseRewriteContextKey = claudeMaxResponseRewriteContextKeyType{placeholder

func withClaudeMaxResponseRewriteContext(ctx context.Context, c *gin.Context, parsed *ParsedRequest) context.Context {
	if ctx == nil {
		ctx = context.Background()
placeholder
	value := claudeMaxResponseRewriteContext{
		Parsed: parsed,
		Group:  claudeMaxGroupFromGinContext(c),
placeholder
	return context.WithValue(ctx, claudeMaxResponseRewriteContextKey, value)
placeholder

func claudeMaxResponseRewriteContextFromContext(ctx context.Context) claudeMaxResponseRewriteContext {
	if ctx == nil {
		return claudeMaxResponseRewriteContext{placeholder
placeholder
	value, _ := ctx.Value(claudeMaxResponseRewriteContextKey).(claudeMaxResponseRewriteContext)
	return value
placeholder

func claudeMaxGroupFromGinContext(c *gin.Context) *Group {
	if c == nil {
		return nil
placeholder
	raw, exists := c.Get("api_key")
	if !exists {
		return nil
placeholder
	apiKey, ok := raw.(*APIKey)
	if !ok || apiKey == nil {
		return nil
placeholder
	return apiKey.Group
placeholder

func applyClaudeMaxSimulationToUsage(ctx context.Context, usage *ClaudeUsage, model string, accountID int64) claudeMaxCacheBillingOutcome {
	var out claudeMaxCacheBillingOutcome
	if usage == nil {
		return out
placeholder
	rewriteCtx := claudeMaxResponseRewriteContextFromContext(ctx)
	return applyClaudeMaxCacheBillingPolicyToUsage(usage, rewriteCtx.Parsed, rewriteCtx.Group, model, accountID)
placeholder

func applyClaudeMaxSimulationToUsageJSONMap(ctx context.Context, usageObj map[string]any, model string, accountID int64) claudeMaxCacheBillingOutcome {
	var out claudeMaxCacheBillingOutcome
	if usageObj == nil {
		return out
placeholder
	usage := claudeUsageFromJSONMap(usageObj)
	out = applyClaudeMaxSimulationToUsage(ctx, &usage, model, accountID)
	if out.Simulated {
		rewriteClaudeUsageJSONMap(usageObj, usage)
placeholder
	return out
placeholder

func rewriteClaudeUsageJSONBytes(body []byte, usage ClaudeUsage) []byte {
	updated := body
	var err error

	updated, err = sjson.SetBytes(updated, "usage.input_tokens", usage.InputTokens)
	if err != nil {
		return body
placeholder
	updated, err = sjson.SetBytes(updated, "usage.cache_creation_input_tokens", usage.CacheCreationInputTokens)
	if err != nil {
		return body
placeholder
	updated, err = sjson.SetBytes(updated, "usage.cache_creation.ephemeral_5m_input_tokens", usage.CacheCreation5mTokens)
	if err != nil {
		return body
placeholder
	updated, err = sjson.SetBytes(updated, "usage.cache_creation.ephemeral_1h_input_tokens", usage.CacheCreation1hTokens)
	if err != nil {
		return body
placeholder
	return updated
placeholder

func claudeUsageFromJSONMap(usageObj map[string]any) ClaudeUsage {
	var usage ClaudeUsage
	if usageObj == nil {
		return usage
placeholder

	usage.InputTokens = usageIntFromAny(usageObj["input_tokens"])
	usage.OutputTokens = usageIntFromAny(usageObj["output_tokens"])
	usage.CacheCreationInputTokens = usageIntFromAny(usageObj["cache_creation_input_tokens"])
	usage.CacheReadInputTokens = usageIntFromAny(usageObj["cache_read_input_tokens"])

	if ccObj, ok := usageObj["cache_creation"].(map[string]any); ok {
		usage.CacheCreation5mTokens = usageIntFromAny(ccObj["ephemeral_5m_input_tokens"])
		usage.CacheCreation1hTokens = usageIntFromAny(ccObj["ephemeral_1h_input_tokens"])
placeholder
	return usage
placeholder

func rewriteClaudeUsageJSONMap(usageObj map[string]any, usage ClaudeUsage) {
	if usageObj == nil {
		return
placeholder
	usageObj["input_tokens"] = usage.InputTokens
	usageObj["cache_creation_input_tokens"] = usage.CacheCreationInputTokens

	ccObj, _ := usageObj["cache_creation"].(map[string]any)
	if ccObj == nil {
		ccObj = make(map[string]any, 2)
		usageObj["cache_creation"] = ccObj
placeholder
	ccObj["ephemeral_5m_input_tokens"] = usage.CacheCreation5mTokens
	ccObj["ephemeral_1h_input_tokens"] = usage.CacheCreation1hTokens
placeholder

func usageIntFromAny(v any) int {
	switch value := v.(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	case json.Number:
		if n, err := value.Int64(); err == nil {
			return int(n)
	placeholder
placeholder
	return 0
placeholder
