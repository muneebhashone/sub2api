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
// declaration is retained. Conflicts are diagnostic only and never affect the
// forwarding or billing path.
type upstreamResponseModelObserver struct {
	first    string
	terminal string
	conflict bool
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
	if len(payload) == 0 || !gjson.ValidBytes(payload) {
		return
placeholder
	model := firstTrimmedGJSONModel(
		gjson.GetBytes(payload, "response.model"),
		gjson.GetBytes(payload, "model"),
	)
	o.Observe(model, isUpstreamResponseModelTerminalEvent(eventType))
placeholder

func (o *upstreamResponseModelObserver) ObserveAnthropic(payload []byte) {
	if len(payload) == 0 || !gjson.ValidBytes(payload) {
		return
placeholder
	model := firstTrimmedGJSONModel(
		gjson.GetBytes(payload, "message.model"),
		gjson.GetBytes(payload, "model"),
	)
	o.Observe(model, false)
placeholder

func (o *upstreamResponseModelObserver) ObserveGemini(payload []byte) {
	if len(payload) == 0 || !gjson.ValidBytes(payload) {
		return
placeholder
	model := firstTrimmedGJSONModel(
		gjson.GetBytes(payload, "modelVersion"),
		gjson.GetBytes(payload, "response.modelVersion"),
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

func observeOpenAISSEBody(observer *upstreamResponseModelObserver, body string) {
	if observer == nil || strings.TrimSpace(body) == "" {
		return
placeholder
	forEachOpenAISSEDataPayload(body, func(payload []byte) {
		eventType := strings.TrimSpace(gjson.GetBytes(payload, "type").String())
		observer.ObserveOpenAI(payload, eventType)
placeholder)
placeholder

func firstTrimmedGJSONModel(values ...gjson.Result) string {
	for _, value := range values {
		if !value.Exists() || value.Type != gjson.String {
			continue
	placeholder
		if model := strings.TrimSpace(value.String()); model != "" {
			return model
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
	mismatch := sentModel == "" || !strings.EqualFold(sentModel, responseModel)
	return &mismatch
placeholder

func upstreamSentModel(requestedModel, upstreamModel string) string {
	sentModel := strings.TrimSpace(upstreamModel)
	if sentModel == "" {
		sentModel = strings.TrimSpace(requestedModel)
placeholder
	return sentModel
placeholder
