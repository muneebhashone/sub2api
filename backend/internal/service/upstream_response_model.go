package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

const (
	upstreamResponseModelObserverContextKey = "upstream_response_model_observer"
	upstreamResponseModelMaxLength          = 200
)

// upstreamResponseModelObserver tracks one forwarding attempt (or one WS turn).
// A terminal declaration wins over an earlier declaration; otherwise the first
// declaration is retained. Observation never affects the forwarding path.
//
// Billing normally ignores the observed model as well; the only exception is a
// channel explicitly configured with billing_model_source = response_model,
// where a conflict flag makes billing fall back to the baseline model
// (see responseModelBillingDeclaration).
//
// The same observer also records the service tier the upstream reports having
// used (OpenAI service_tier, Anthropic usage.speed). The billable tier is
// resolved by resolvedOpenAIUpstreamServiceTierFromObserver (upstream echo
// first, outbound body tier as fallback); the upstream ResolveBillingServiceTier
// only-lowers path additionally audits downgrades at usage-record time.
type upstreamResponseModelObserver struct {
	first    string
	terminal string
	conflict bool

	// firstTier holds the first non-terminal tier declaration; it is discarded
	// when later non-terminal declarations disagree. terminalTier comes from a
	// terminal event and always wins.
	firstTier         string
	firstTierConflict bool
	terminalTier      string
placeholder

func (o *upstreamResponseModelObserver) Observe(model string, terminal bool) {
	model = normalizeObservedUpstreamResponseModel(model)
	if model == "" {
		return
placeholder
	current := o.Model()
	if current != "" && !strings.EqualFold(current, model) {
		o.conflict = true
placeholder
	if terminal {
		o.terminal = model
		return
placeholder
	if o.first == "" {
		o.first = model
placeholder
placeholder

func normalizeObservedUpstreamResponseModel(model string) string {
	model = strings.TrimSpace(model)
	if model == "" {
		return ""
placeholder
	runes := []rune(model)
	if len(runes) > upstreamResponseModelMaxLength {
		model = string(runes[:upstreamResponseModelMaxLength])
placeholder
	return model
placeholder

func (o *upstreamResponseModelObserver) ObserveOpenAI(payload []byte, eventType string) {
	model := firstValidTrimmedGJSONString(payload, "response.model", "model")
	terminal := isUpstreamResponseModelTerminalEvent(eventType)
	o.Observe(model, terminal)
	// Every payload that declares a service tier also declares a model, so
	// model-free delta frames skip the extra lookups entirely.
	if model == "" {
		return
placeholder
	// Non-terminal Responses API events echo the requested tier rather than the
	// tier actually used. Only terminal events and untyped payloads (chat
	// completions chunks, non-streaming bodies) report the processing tier.
	if !terminal && strings.TrimSpace(eventType) != "" {
		return
placeholder
	tier := normalizeObservedOpenAIServiceTier(firstValidTrimmedGJSONString(payload, "response.service_tier", "service_tier"))
	o.ObserveServiceTier(tier, terminal)
placeholder

func (o *upstreamResponseModelObserver) ObserveAnthropic(payload []byte) {
	model := firstValidTrimmedGJSONString(payload, "message.model", "model")
	o.Observe(model, false)
	// usage.speed travels with the message object (message_start in streams,
	// the top-level body otherwise), i.e. only in payloads that declare a model.
	if model == "" {
		return
placeholder
	tier := normalizeObservedAnthropicSpeed(firstValidTrimmedGJSONString(payload, "message.usage.speed", "usage.speed"))
	o.ObserveServiceTier(tier, false)
placeholder

// ObserveServiceTier records a tier declared by the upstream response. A
// terminal declaration always wins; non-terminal declarations are only trusted
// when they agree with each other.
func (o *upstreamResponseModelObserver) ObserveServiceTier(tier string, terminal bool) {
	if o == nil || tier == "" {
		return
placeholder
	if terminal {
		o.terminalTier = tier
		return
placeholder
	if o.firstTier == "" {
		o.firstTier = tier
		return
placeholder
	if o.firstTier != tier {
		o.firstTierConflict = true
placeholder
placeholder

// ServiceTier returns the tier the upstream reports having used, or "" when the
// response never declared one unambiguously.
func (o *upstreamResponseModelObserver) ServiceTier() string {
	if o == nil {
		return ""
placeholder
	if o.terminalTier != "" {
		return o.terminalTier
placeholder
	if o.firstTierConflict {
		return ""
placeholder
	return o.firstTier
placeholder

// normalizeObservedOpenAIServiceTier maps a tier reported by an OpenAI response
// onto the billing vocabulary. "auto" never describes a processing tier and
// unknown values are ignored rather than guessed at.
func normalizeObservedOpenAIServiceTier(raw string) string {
	switch value := strings.ToLower(strings.TrimSpace(raw)); value {
	case "priority", "fast":
		return OpenAIFastTierPriority
	case "default", "flex", "scale":
		return value
	default:
		return ""
placeholder
placeholder

// normalizeObservedAnthropicSpeed maps Anthropic usage.speed onto the billing
// vocabulary: "fast" is the billable fast-mode tier, "standard" the base rate.
func normalizeObservedAnthropicSpeed(raw string) string {
	switch value := strings.ToLower(strings.TrimSpace(raw)); value {
	case "fast", "standard":
		return value
	default:
		return ""
placeholder
placeholder

func (o *upstreamResponseModelObserver) ObserveGemini(payload []byte) {
	model := firstValidTrimmedGJSONString(
		payload,
		"modelVersion",
		"response.modelVersion",
		"response.response.modelVersion",
	)
	// Gemini streaming has no universal terminal event carrying modelVersion;
	// treating each declaration as terminal retains the latest chunk.
	o.Observe(model, true)
placeholder

func (o *upstreamResponseModelObserver) Model() string {
	if o == nil {
		return ""
placeholder
	if o.terminal != "" {
		return o.terminal
placeholder
	return o.first
placeholder

func (o *upstreamResponseModelObserver) Conflict() bool {
	return o != nil && o.conflict
placeholder

func beginUpstreamResponseModelObservation(c *gin.Context) *upstreamResponseModelObserver {
	observer := &upstreamResponseModelObserver{placeholder
	if c != nil {
		c.Set(upstreamResponseModelObserverContextKey, observer)
placeholder
	return observer
placeholder

func upstreamResponseModelObserverFromContext(c *gin.Context) *upstreamResponseModelObserver {
	if c == nil {
		return nil
placeholder
	value, ok := c.Get(upstreamResponseModelObserverContextKey)
	if !ok {
		return nil
placeholder
	observer, _ := value.(*upstreamResponseModelObserver)
	return observer
placeholder

func observedUpstreamResponseModel(c *gin.Context) string {
	return upstreamResponseModelObserverFromContext(c).Model()
placeholder

func observedUpstreamResponseModelConflict(c *gin.Context) bool {
	return upstreamResponseModelObserverFromContext(c).Conflict()
placeholder

func observedUpstreamResponseServiceTier(c *gin.Context) string {
	return upstreamResponseModelObserverFromContext(c).ServiceTier()
placeholder

// resolvedOpenAIUpstreamServiceTierFromObserver 返回计费/用量日志实际使用的
// service tier：
//
//  1. observer 记录到的上游真实回显优先——只有上游实际给了 priority/fast 才按
//     Fast 计费；上游回显 default/flex/auto 等则如实采用并据此计费；
//  2. 上游未回显时，回退到「最终出站 body」里的 tier（经过 fast policy
//     filter/force 之后），保证 policy filter 删掉字段后不再按原请求 Fast 计费。
//
// HTTP→WS 等使用局部 observer 的路径必须把该 observer 传进来，不能只读
// Gin context——局部 observer 不会自动写入 context。
func resolvedOpenAIUpstreamServiceTierFromObserver(observer *upstreamResponseModelObserver, outboundBodyTier *string) *string {
	if observer != nil {
		if tier := strings.TrimSpace(observer.ServiceTier()); tier != "" {
			return normalizeOpenAIServiceTier(tier)
	placeholder
placeholder
	return outboundBodyTier
placeholder

// resolvedOpenAIUpstreamServiceTier 读取 Gin context 上的 observer 后委托
// resolvedOpenAIUpstreamServiceTierFromObserver。标准 HTTP 转发路径通过
// beginUpstreamResponseModelObservation 把 observer 挂到 context；局部
// observer 路径应直接调用 FromObserver。
func resolvedOpenAIUpstreamServiceTier(c *gin.Context, outboundBodyTier *string) *string {
	return resolvedOpenAIUpstreamServiceTierFromObserver(upstreamResponseModelObserverFromContext(c), outboundBodyTier)
placeholder

func observeOpenAISSEBody(observer *upstreamResponseModelObserver, body string) {
	if observer == nil || strings.TrimSpace(body) == "" {
		return
placeholder
	forEachOpenAISSEFrame(body, func(eventType string, payload []byte) {
		observer.ObserveOpenAI(payload, eventType)
placeholder)
placeholder

func firstValidTrimmedGJSONString(payload []byte, paths ...string) string {
	if len(payload) == 0 {
		return ""
placeholder
	for _, path := range paths {
		value := gjson.GetBytes(payload, path)
		if !value.Exists() || value.Type != gjson.String {
			continue
	placeholder
		if text := strings.TrimSpace(value.String()); text != "" {
			// Validate only after finding a candidate. This avoids a full validation
			// pass on the common model-free delta path while still rejecting malformed
			// payloads that appear to declare a value.
			if !gjson.ValidBytes(payload) {
				return ""
		placeholder
			return text
	placeholder
placeholder
	return ""
placeholder

func isUpstreamResponseModelTerminalEvent(eventType string) bool {
	switch strings.TrimSpace(eventType) {
	case "response.completed", "response.done", "response.failed", "response.incomplete", "response.cancelled", "response.canceled":
		return true
	default:
		return false
placeholder
placeholder

func upstreamModelMismatch(sentModel, responseModel string) *bool {
	responseModel = strings.TrimSpace(responseModel)
	if responseModel == "" {
		return nil
placeholder
	sentModel = strings.TrimSpace(sentModel)
	mismatch := sentModel == "" || !upstreamModelsMatchForAudit(sentModel, responseModel)
	return &mismatch
placeholder

func upstreamModelsMatchForAudit(sentModel, responseModel string) bool {
	if strings.EqualFold(sentModel, responseModel) {
		return true
placeholder

	// xAI reports the runtime build ID for these supported public aliases.
	// Canonicalize only for mismatch auditing; keep the raw response model for
	// observability and for the separate response-model billing safeguards.
	sentGrokModel := canonicalGrokBuildRuntimeModel(sentModel)
	return sentGrokModel != "" && sentGrokModel == canonicalGrokBuildRuntimeModel(responseModel)
placeholder

func canonicalGrokBuildRuntimeModel(model string) string {
	switch strings.ToLower(strings.TrimSpace(model)) {
	case "grok-4.5", "grok-4.5-latest", "grok-4.5-build":
		return "grok-4.5-build"
	case "grok-4.6", "grok-4.6-latest", "grok-4.6-build":
		return "grok-4.6-build"
	default:
		return ""
placeholder
placeholder

func upstreamSentModel(requestedModel, upstreamModel string) string {
	sentModel := strings.TrimSpace(upstreamModel)
	if sentModel == "" {
		sentModel = strings.TrimSpace(requestedModel)
placeholder
	return sentModel
placeholder
