package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var opsDashboardSnapshotV2Cache = newSnapshotCache(30 * time.Second)

type opsDashboardSnapshotV2Response struct {
	GeneratedAt string `json:"generated_at"`

	Overview        *service.OpsDashboardOverview       `json:"overview"`
	ThroughputTrend *service.OpsThroughputTrendResponse `json:"throughput_trend"`
	ErrorTrend      *service.OpsErrorTrendResponse      `json:"error_trend"`
placeholder

type opsDashboardSnapshotV2CacheKey struct {
	StartTime    string               `json:"start_time"`
	EndTime      string               `json:"end_time"`
	Platform     string               `json:"platform"`
	GroupID      *int64               `json:"group_id"`
	QueryMode    service.OpsQueryMode `json:"mode"`
	BucketSecond int                  `json:"bucket_second"`
placeholder

// GetDashboardSnapshotV2 returns ops dashboard core snapshot in one request.
// GET /api/v1/admin/ops/dashboard/snapshot-v2
func (h *OpsHandler) GetDashboardSnapshotV2(c *gin.Context) {
	if h.opsService == nil {
		response.Error(c, http.StatusServiceUnavailable, "Ops service not available")
		return
placeholder
	if err := h.opsService.RequireMonitoringEnabled(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	startTime, endTime, err := parseOpsTimeRange(c, "1h")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
placeholder

	filter := &service.OpsDashboardFilter{
		StartTime: startTime,
		EndTime:   endTime,
		Platform:  strings.TrimSpace(c.Query("platform")),
		QueryMode: parseOpsQueryMode(c),
placeholder
	if v := strings.TrimSpace(c.Query("group_id")); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			response.BadRequest(c, "Invalid group_id")
			return
	placeholder
		filter.GroupID = &id
placeholder
	bucketSeconds := pickThroughputBucketSeconds(endTime.Sub(startTime))

	keyRaw, _ := json.Marshal(opsDashboardSnapshotV2CacheKey{
		StartTime:    startTime.UTC().Format(time.RFC3339),
		EndTime:      endTime.UTC().Format(time.RFC3339),
		Platform:     filter.Platform,
		GroupID:      filter.GroupID,
		QueryMode:    filter.QueryMode,
		BucketSecond: bucketSeconds,
placeholder)
	cacheKey := string(keyRaw)

	if cached, ok := opsDashboardSnapshotV2Cache.Get(cacheKey); ok {
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

	var (
		overview *service.OpsDashboardOverview
		trend    *service.OpsThroughputTrendResponse
		errTrend *service.OpsErrorTrendResponse
	)
	g, gctx := errgroup.WithContext(c.Request.Context())
	g.Go(func() error {
		f := *filter
		result, err := h.opsService.GetDashboardOverview(gctx, &f)
		if err != nil {
			return err
	placeholder
		overview = result
		return nil
placeholder)
	g.Go(func() error {
		f := *filter
		result, err := h.opsService.GetThroughputTrend(gctx, &f, bucketSeconds)
		if err != nil {
			return err
	placeholder
		trend = result
		return nil
placeholder)
	g.Go(func() error {
		f := *filter
		result, err := h.opsService.GetErrorTrend(gctx, &f, bucketSeconds)
		if err != nil {
			return err
	placeholder
		errTrend = result
		return nil
placeholder)
	if err := g.Wait(); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	resp := &opsDashboardSnapshotV2Response{
		GeneratedAt:     time.Now().UTC().Format(time.RFC3339),
		Overview:        overview,
		ThroughputTrend: trend,
		ErrorTrend:      errTrend,
placeholder

	cached := opsDashboardSnapshotV2Cache.Set(cacheKey, resp)
	if cached.ETag != "" {
		c.Header("ETag", cached.ETag)
		c.Header("Vary", "If-None-Match")
placeholder
	c.Header("X-Snapshot-Cache", "miss")
	response.Success(c, resp)
placeholder
