package admin

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var validOpsAlertMetricTypes = []string{
	"success_rate",
	"error_rate",
	"upstream_error_rate",
	"p95_latency_ms",
	"p99_latency_ms",
	"cpu_usage_percent",
	"memory_usage_percent",
	"concurrency_queue_depth",
placeholder

var validOpsAlertMetricTypeSet = func() map[string]struct{placeholder {
	set := make(map[string]struct{placeholder, len(validOpsAlertMetricTypes))
	for _, v := range validOpsAlertMetricTypes {
		set[v] = struct{placeholder{placeholder
placeholder
	return set
placeholder()

var validOpsAlertOperators = []string{">", "<", ">=", "<=", "==", "!="placeholder

var validOpsAlertOperatorSet = func() map[string]struct{placeholder {
	set := make(map[string]struct{placeholder, len(validOpsAlertOperators))
	for _, v := range validOpsAlertOperators {
		set[v] = struct{placeholder{placeholder
placeholder
	return set
placeholder()

var validOpsAlertSeverities = []string{"P0", "P1", "P2", "P3"placeholder

var validOpsAlertSeveritySet = func() map[string]struct{placeholder {
	set := make(map[string]struct{placeholder, len(validOpsAlertSeverities))
	for _, v := range validOpsAlertSeverities {
		set[v] = struct{placeholder{placeholder
placeholder
	return set
placeholder()

type opsAlertRuleValidatedInput struct {
	Name       string
	MetricType string
	Operator   string
	Threshold  float64

	Severity string

	WindowMinutes    int
	SustainedMinutes int
	CooldownMinutes  int

	Enabled     bool
	NotifyEmail bool

	WindowProvided    bool
	SustainedProvided bool
	CooldownProvided  bool
	SeverityProvided  bool
	EnabledProvided   bool
	NotifyProvided    bool
placeholder

func isPercentOrRateMetric(metricType string) bool {
	switch metricType {
	case "success_rate",
		"error_rate",
		"upstream_error_rate",
		"cpu_usage_percent",
		"memory_usage_percent":
		return true
	default:
		return false
placeholder
placeholder

func validateOpsAlertRulePayload(raw map[string]json.RawMessage) (*opsAlertRuleValidatedInput, error) {
	if raw == nil {
		return nil, fmt.Errorf("invalid request body")
placeholder

	requiredFields := []string{"name", "metric_type", "operator", "threshold"placeholder
	for _, field := range requiredFields {
		if _, ok := raw[field]; !ok {
			return nil, fmt.Errorf("%s is required", field)
	placeholder
placeholder

	var name string
	if err := json.Unmarshal(raw["name"], &name); err != nil || strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name is required")
placeholder
	name = strings.TrimSpace(name)

	var metricType string
	if err := json.Unmarshal(raw["metric_type"], &metricType); err != nil || strings.TrimSpace(metricType) == "" {
		return nil, fmt.Errorf("metric_type is required")
placeholder
	metricType = strings.TrimSpace(metricType)
	if _, ok := validOpsAlertMetricTypeSet[metricType]; !ok {
		return nil, fmt.Errorf("metric_type must be one of: %s", strings.Join(validOpsAlertMetricTypes, ", "))
placeholder

	var operator string
	if err := json.Unmarshal(raw["operator"], &operator); err != nil || strings.TrimSpace(operator) == "" {
		return nil, fmt.Errorf("operator is required")
placeholder
	operator = strings.TrimSpace(operator)
	if _, ok := validOpsAlertOperatorSet[operator]; !ok {
		return nil, fmt.Errorf("operator must be one of: %s", strings.Join(validOpsAlertOperators, ", "))
placeholder

	var threshold float64
	if err := json.Unmarshal(raw["threshold"], &threshold); err != nil {
		return nil, fmt.Errorf("threshold must be a number")
placeholder
	if math.IsNaN(threshold) || math.IsInf(threshold, 0) {
		return nil, fmt.Errorf("threshold must be a finite number")
placeholder
	if isPercentOrRateMetric(metricType) {
		if threshold < 0 || threshold > 100 {
			return nil, fmt.Errorf("threshold must be between 0 and 100 for metric_type %s", metricType)
	placeholder
placeholder else if threshold < 0 {
		return nil, fmt.Errorf("threshold must be >= 0")
placeholder

	validated := &opsAlertRuleValidatedInput{
		Name:       name,
		MetricType: metricType,
		Operator:   operator,
		Threshold:  threshold,
placeholder

	if v, ok := raw["severity"]; ok {
		validated.SeverityProvided = true
		var sev string
		if err := json.Unmarshal(v, &sev); err != nil {
			return nil, fmt.Errorf("severity must be a string")
	placeholder
		sev = strings.ToUpper(strings.TrimSpace(sev))
		if sev != "" {
			if _, ok := validOpsAlertSeveritySet[sev]; !ok {
				return nil, fmt.Errorf("severity must be one of: %s", strings.Join(validOpsAlertSeverities, ", "))
		placeholder
			validated.Severity = sev
	placeholder
placeholder
	if validated.Severity == "" {
		validated.Severity = "P2"
placeholder

	if v, ok := raw["enabled"]; ok {
		validated.EnabledProvided = true
		if err := json.Unmarshal(v, &validated.Enabled); err != nil {
			return nil, fmt.Errorf("enabled must be a boolean")
	placeholder
placeholder else {
		validated.Enabled = true
placeholder

	if v, ok := raw["notify_email"]; ok {
		validated.NotifyProvided = true
		if err := json.Unmarshal(v, &validated.NotifyEmail); err != nil {
			return nil, fmt.Errorf("notify_email must be a boolean")
	placeholder
placeholder else {
		validated.NotifyEmail = true
placeholder

	if v, ok := raw["window_minutes"]; ok {
		validated.WindowProvided = true
		if err := json.Unmarshal(v, &validated.WindowMinutes); err != nil {
			return nil, fmt.Errorf("window_minutes must be an integer")
	placeholder
		switch validated.WindowMinutes {
		case 1, 5, 60:
		default:
			return nil, fmt.Errorf("window_minutes must be one of: 1, 5, 60")
	placeholder
placeholder else {
		validated.WindowMinutes = 1
placeholder

	if v, ok := raw["sustained_minutes"]; ok {
		validated.SustainedProvided = true
		if err := json.Unmarshal(v, &validated.SustainedMinutes); err != nil {
			return nil, fmt.Errorf("sustained_minutes must be an integer")
	placeholder
		if validated.SustainedMinutes < 1 || validated.SustainedMinutes > 1440 {
			return nil, fmt.Errorf("sustained_minutes must be between 1 and 1440")
	placeholder
placeholder else {
		validated.SustainedMinutes = 1
placeholder

	if v, ok := raw["cooldown_minutes"]; ok {
		validated.CooldownProvided = true
		if err := json.Unmarshal(v, &validated.CooldownMinutes); err != nil {
			return nil, fmt.Errorf("cooldown_minutes must be an integer")
	placeholder
		if validated.CooldownMinutes < 0 || validated.CooldownMinutes > 1440 {
			return nil, fmt.Errorf("cooldown_minutes must be between 0 and 1440")
	placeholder
placeholder else {
		validated.CooldownMinutes = 0
placeholder

	return validated, nil
placeholder

// ListAlertRules returns all ops alert rules.
// GET /api/v1/admin/ops/alert-rules
func (h *OpsHandler) ListAlertRules(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	rules, err := h.opsService.ListAlertRules(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, rules)
placeholder

// CreateAlertRule creates an ops alert rule.
// POST /api/v1/admin/ops/alert-rules
func (h *OpsHandler) CreateAlertRule(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	var raw map[string]json.RawMessage
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
placeholder
	validated, err := validateOpsAlertRulePayload(raw)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	var rule service.OpsAlertRule
	if err := c.ShouldBindBodyWith(&rule, binding.JSON); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
placeholder

	rule.Name = validated.Name
	rule.MetricType = validated.MetricType
	rule.Operator = validated.Operator
	rule.Threshold = validated.Threshold
	rule.WindowMinutes = validated.WindowMinutes
	rule.SustainedMinutes = validated.SustainedMinutes
	rule.CooldownMinutes = validated.CooldownMinutes
	rule.Severity = validated.Severity
	rule.Enabled = validated.Enabled
	rule.NotifyEmail = validated.NotifyEmail

	created, err := h.opsService.CreateAlertRule(c.Request.Context(), &rule)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, created)
placeholder

// UpdateAlertRule updates an existing ops alert rule.
// PUT /api/v1/admin/ops/alert-rules/:id
func (h *OpsHandler) UpdateAlertRule(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid rule ID")
		return
placeholder

	var raw map[string]json.RawMessage
	if err := c.ShouldBindBodyWith(&raw, binding.JSON); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
placeholder
	validated, err := validateOpsAlertRulePayload(raw)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	var rule service.OpsAlertRule
	if err := c.ShouldBindBodyWith(&rule, binding.JSON); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
placeholder

	rule.ID = id
	rule.Name = validated.Name
	rule.MetricType = validated.MetricType
	rule.Operator = validated.Operator
	rule.Threshold = validated.Threshold
	rule.WindowMinutes = validated.WindowMinutes
	rule.SustainedMinutes = validated.SustainedMinutes
	rule.CooldownMinutes = validated.CooldownMinutes
	rule.Severity = validated.Severity
	rule.Enabled = validated.Enabled
	rule.NotifyEmail = validated.NotifyEmail

	updated, err := h.opsService.UpdateAlertRule(c.Request.Context(), &rule)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, updated)
placeholder

// DeleteAlertRule deletes an ops alert rule.
// DELETE /api/v1/admin/ops/alert-rules/:id
func (h *OpsHandler) DeleteAlertRule(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid rule ID")
		return
placeholder

	if err := h.opsService.DeleteAlertRule(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, gin.H{"deleted": trueplaceholder)
placeholder

// ListAlertEvents lists recent ops alert events.
// GET /api/v1/admin/ops/alert-events
func (h *OpsHandler) ListAlertEvents(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	limit := 100
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			response.BadRequest(c, "Invalid limit")
			return
	placeholder
		limit = n
placeholder

	filter := &service.OpsAlertEventFilter{
		Limit:    limit,
		Status:   strings.TrimSpace(c.Query("status")),
		Severity: strings.TrimSpace(c.Query("severity")),
placeholder

	// Optional global filter support (platform/group/time range).
	if platform := strings.TrimSpace(c.Query("platform")); platform != "" {
		filter.Platform = platform
placeholder
	if v := strings.TrimSpace(c.Query("group_id")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			response.BadRequest(c, "Invalid group_id")
			return
	placeholder
		filter.GroupID = &id
placeholder
	if startTime, endTime, err := parseOpsTimeRange(c, "24h"); err == nil {
		// Only apply when explicitly provided to avoid surprising default narrowing.
		if strings.TrimSpace(c.Query("start_time")) != "" || strings.TrimSpace(c.Query("end_time")) != "" || strings.TrimSpace(c.Query("time_range")) != "" {
			filter.StartTime = &startTime
			filter.EndTime = &endTime
	placeholder
placeholder else {
		response.BadRequest(c, err.Error())
		return
placeholder

	events, err := h.opsService.ListAlertEvents(c.Request.Context(), filter)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, events)
placeholder
