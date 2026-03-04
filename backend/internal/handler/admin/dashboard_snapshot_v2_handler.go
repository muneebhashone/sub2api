package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

var dashboardSnapshotV2Cache = newSnapshotCache(30 * time.Second)

type dashboardSnapshotV2Stats struct {
	usagestats.DashboardStats
	Uptime int64 `json:"uptime"`
placeholder

type dashboardSnapshotV2Response struct {
	GeneratedAt string `json:"generated_at"`

	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Granularity string `json:"granularity"`

	Stats      *dashboardSnapshotV2Stats        `json:"stats,omitempty"`
	Trend      []usagestats.TrendDataPoint      `json:"trend,omitempty"`
	Models     []usagestats.ModelStat           `json:"models,omitempty"`
	Groups     []usagestats.GroupStat           `json:"groups,omitempty"`
	UsersTrend []usagestats.UserUsageTrendPoint `json:"users_trend,omitempty"`
placeholder

type dashboardSnapshotV2Filters struct {
	UserID      int64
	APIKeyID    int64
	AccountID   int64
	GroupID     int64
	Model       string
	RequestType *int16
	Stream      *bool
	BillingType *int8
placeholder

type dashboardSnapshotV2CacheKey struct {
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	Granularity       string `json:"granularity"`
	UserID            int64  `json:"user_id"`
	APIKeyID          int64  `json:"api_key_id"`
	AccountID         int64  `json:"account_id"`
	GroupID           int64  `json:"group_id"`
	Model             string `json:"model"`
	RequestType       *int16 `json:"request_type"`
	Stream            *bool  `json:"stream"`
	BillingType       *int8  `json:"billing_type"`
	IncludeStats      bool   `json:"include_stats"`
	IncludeTrend      bool   `json:"include_trend"`
	IncludeModels     bool   `json:"include_models"`
	IncludeGroups     bool   `json:"include_groups"`
	IncludeUsersTrend bool   `json:"include_users_trend"`
	UsersTrendLimit   int    `json:"users_trend_limit"`
placeholder

func (h *DashboardHandler) GetSnapshotV2(c *gin.Context) {
	startTime, endTime := parseTimeRange(c)
	granularity := strings.TrimSpace(c.DefaultQuery("granularity", "day"))
	if granularity != "hour" {
		granularity = "day"
placeholder

	includeStats := parseBoolQueryWithDefault(c.Query("include_stats"), true)
	includeTrend := parseBoolQueryWithDefault(c.Query("include_trend"), true)
	includeModels := parseBoolQueryWithDefault(c.Query("include_model_stats"), true)
	includeGroups := parseBoolQueryWithDefault(c.Query("include_group_stats"), false)
	includeUsersTrend := parseBoolQueryWithDefault(c.Query("include_users_trend"), false)
	usersTrendLimit := 12
	if raw := strings.TrimSpace(c.Query("users_trend_limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 50 {
			usersTrendLimit = parsed
	placeholder
placeholder

	filters, err := parseDashboardSnapshotV2Filters(c)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	keyRaw, _ := json.Marshal(dashboardSnapshotV2CacheKey{
		StartTime:         startTime.UTC().Format(time.RFC3339),
		EndTime:           endTime.UTC().Format(time.RFC3339),
		Granularity:       granularity,
		UserID:            filters.UserID,
		APIKeyID:          filters.APIKeyID,
		AccountID:         filters.AccountID,
		GroupID:           filters.GroupID,
		Model:             filters.Model,
		RequestType:       filters.RequestType,
		Stream:            filters.Stream,
		BillingType:       filters.BillingType,
		IncludeStats:      includeStats,
		IncludeTrend:      includeTrend,
		IncludeModels:     includeModels,
		IncludeGroups:     includeGroups,
		IncludeUsersTrend: includeUsersTrend,
		UsersTrendLimit:   usersTrendLimit,
placeholder)
	cacheKey := string(keyRaw)

	if cached, ok := dashboardSnapshotV2Cache.Get(cacheKey); ok {
		if cached.ETag != "" {
			c.Header("ETag", cached.ETag)
			c.Header("Vary", "If-None-Match")
			if ifNoneMatchMatched(c.GetHeader("If-None-Match"), cached.ETag) {
				c.Status(http.StatusNotModified)
				return
		placeholder
	placeholder
		c.Header("X-Snapshot-Cache", "hit")
		response.Success(c, cached.Payload)
		return
placeholder

	resp := &dashboardSnapshotV2Response{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		StartDate:   startTime.Format("2006-01-02"),
		EndDate:     endTime.Add(-24 * time.Hour).Format("2006-01-02"),
		Granularity: granularity,
placeholder

	if includeStats {
		stats, err := h.dashboardService.GetDashboardStats(c.Request.Context())
		if err != nil {
			response.Error(c, 500, "Failed to get dashboard statistics")
			return
	placeholder
		resp.Stats = &dashboardSnapshotV2Stats{
			DashboardStats: *stats,
			Uptime:         int64(time.Since(h.startTime).Seconds()),
	placeholder
placeholder

	if includeTrend {
		trend, err := h.dashboardService.GetUsageTrendWithFilters(
			c.Request.Context(),
			startTime,
			endTime,
			granularity,
			filters.UserID,
			filters.APIKeyID,
			filters.AccountID,
			filters.GroupID,
			filters.Model,
			filters.RequestType,
			filters.Stream,
			filters.BillingType,
		)
		if err != nil {
			response.Error(c, 500, "Failed to get usage trend")
			return
	placeholder
		resp.Trend = trend
placeholder

	if includeModels {
		models, err := h.dashboardService.GetModelStatsWithFilters(
			c.Request.Context(),
			startTime,
			endTime,
			filters.UserID,
			filters.APIKeyID,
			filters.AccountID,
			filters.GroupID,
			filters.RequestType,
			filters.Stream,
			filters.BillingType,
		)
		if err != nil {
			response.Error(c, 500, "Failed to get model statistics")
			return
	placeholder
		resp.Models = models
placeholder

	if includeGroups {
		groups, err := h.dashboardService.GetGroupStatsWithFilters(
			c.Request.Context(),
			startTime,
			endTime,
			filters.UserID,
			filters.APIKeyID,
			filters.AccountID,
			filters.GroupID,
			filters.RequestType,
			filters.Stream,
			filters.BillingType,
		)
		if err != nil {
			response.Error(c, 500, "Failed to get group statistics")
			return
	placeholder
		resp.Groups = groups
placeholder

	if includeUsersTrend {
		usersTrend, err := h.dashboardService.GetUserUsageTrend(
			c.Request.Context(),
			startTime,
			endTime,
			granularity,
			usersTrendLimit,
		)
		if err != nil {
			response.Error(c, 500, "Failed to get user usage trend")
			return
	placeholder
		resp.UsersTrend = usersTrend
placeholder

	cached := dashboardSnapshotV2Cache.Set(cacheKey, resp)
	if cached.ETag != "" {
		c.Header("ETag", cached.ETag)
		c.Header("Vary", "If-None-Match")
placeholder
	c.Header("X-Snapshot-Cache", "miss")
	response.Success(c, resp)
placeholder

func parseDashboardSnapshotV2Filters(c *gin.Context) (*dashboardSnapshotV2Filters, error) {
	filters := &dashboardSnapshotV2Filters{
		Model: strings.TrimSpace(c.Query("model")),
placeholder

	if userIDStr := strings.TrimSpace(c.Query("user_id")); userIDStr != "" {
		id, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			return nil, err
	placeholder
		filters.UserID = id
placeholder
	if apiKeyIDStr := strings.TrimSpace(c.Query("api_key_id")); apiKeyIDStr != "" {
		id, err := strconv.ParseInt(apiKeyIDStr, 10, 64)
		if err != nil {
			return nil, err
	placeholder
		filters.APIKeyID = id
placeholder
	if accountIDStr := strings.TrimSpace(c.Query("account_id")); accountIDStr != "" {
		id, err := strconv.ParseInt(accountIDStr, 10, 64)
		if err != nil {
			return nil, err
	placeholder
		filters.AccountID = id
placeholder
	if groupIDStr := strings.TrimSpace(c.Query("group_id")); groupIDStr != "" {
		id, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			return nil, err
	placeholder
		filters.GroupID = id
placeholder

	if requestTypeStr := strings.TrimSpace(c.Query("request_type")); requestTypeStr != "" {
		parsed, err := service.ParseUsageRequestType(requestTypeStr)
		if err != nil {
			return nil, err
	placeholder
		value := int16(parsed)
		filters.RequestType = &value
placeholder else if streamStr := strings.TrimSpace(c.Query("stream")); streamStr != "" {
		streamVal, err := strconv.ParseBool(streamStr)
		if err != nil {
			return nil, err
	placeholder
		filters.Stream = &streamVal
placeholder

	if billingTypeStr := strings.TrimSpace(c.Query("billing_type")); billingTypeStr != "" {
		v, err := strconv.ParseInt(billingTypeStr, 10, 8)
		if err != nil {
			return nil, err
	placeholder
		bt := int8(v)
		filters.BillingType = &bt
placeholder

	return filters, nil
placeholder
