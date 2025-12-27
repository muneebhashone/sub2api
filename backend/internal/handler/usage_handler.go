package handler

import (
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// UsageHandler handles usage-related requests
type UsageHandler struct {
	usageService  *service.UsageService
	apiKeyService *service.ApiKeyService
placeholder

// NewUsageHandler creates a new UsageHandler
func NewUsageHandler(usageService *service.UsageService, apiKeyService *service.ApiKeyService) *UsageHandler {
	return &UsageHandler{
		usageService:  usageService,
		apiKeyService: apiKeyService,
placeholder
placeholder

// List handles listing usage records with pagination
// GET /api/v1/usage
func (h *UsageHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	page, pageSize := response.ParsePagination(c)

	var apiKeyID int64
	if apiKeyIDStr := c.Query("api_key_id"); apiKeyIDStr != "" {
		id, err := strconv.ParseInt(apiKeyIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid api_key_id")
			return
	placeholder

		// [Security Fix] Verify API Key ownership to prevent horizontal privilege escalation
		apiKey, err := h.apiKeyService.GetByID(c.Request.Context(), id)
		if err != nil {
			response.ErrorFrom(c, err)
			return
	placeholder
		if apiKey.UserID != subject.UserID {
			response.Forbidden(c, "Not authorized to access this API key's usage records")
			return
	placeholder

		apiKeyID = id
placeholder

	// Parse additional filters
	model := c.Query("model")

	var stream *bool
	if streamStr := c.Query("stream"); streamStr != "" {
		val, err := strconv.ParseBool(streamStr)
		if err != nil {
			response.BadRequest(c, "Invalid stream value, use true or false")
			return
	placeholder
		stream = &val
placeholder

	var billingType *int8
	if billingTypeStr := c.Query("billing_type"); billingTypeStr != "" {
		val, err := strconv.ParseInt(billingTypeStr, 10, 8)
		if err != nil {
			response.BadRequest(c, "Invalid billing_type")
			return
	placeholder
		bt := int8(val)
		billingType = &bt
placeholder

	// Parse date range
	var startTime, endTime *time.Time
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		t, err := timezone.ParseInLocation("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD")
			return
	placeholder
		startTime = &t
placeholder

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		t, err := timezone.ParseInLocation("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD")
			return
	placeholder
		// Set end time to end of day
		t = t.Add(24*time.Hour - time.Nanosecond)
		endTime = &t
placeholder

	params := pagination.PaginationParams{Page: page, PageSize: pageSizeplaceholder
	filters := usagestats.UsageLogFilters{
		UserID:      subject.UserID, // Always filter by current user for security
		ApiKeyID:    apiKeyID,
		Model:       model,
		Stream:      stream,
		BillingType: billingType,
		StartTime:   startTime,
		EndTime:     endTime,
placeholder

	records, result, err := h.usageService.ListWithFilters(c.Request.Context(), params, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	out := make([]dto.UsageLog, 0, len(records))
	for i := range records {
		out = append(out, *dto.UsageLogFromService(&records[i]))
placeholder
	response.Paginated(c, out, result.Total, page, pageSize)
placeholder

// GetByID handles getting a single usage record
// GET /api/v1/usage/:id
func (h *UsageHandler) GetByID(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	usageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid usage ID")
		return
placeholder

	record, err := h.usageService.GetByID(c.Request.Context(), usageID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	// 验证所有权
	if record.UserID != subject.UserID {
		response.Forbidden(c, "Not authorized to access this record")
		return
placeholder

	response.Success(c, dto.UsageLogFromService(record))
placeholder

// Stats handles getting usage statistics
// GET /api/v1/usage/stats
func (h *UsageHandler) Stats(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	var apiKeyID int64
	if apiKeyIDStr := c.Query("api_key_id"); apiKeyIDStr != "" {
		id, err := strconv.ParseInt(apiKeyIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "Invalid api_key_id")
			return
	placeholder

		// [Security Fix] Verify API Key ownership to prevent horizontal privilege escalation
		apiKey, err := h.apiKeyService.GetByID(c.Request.Context(), id)
		if err != nil {
			response.NotFound(c, "API key not found")
			return
	placeholder
		if apiKey.UserID != subject.UserID {
			response.Forbidden(c, "Not authorized to access this API key's statistics")
			return
	placeholder

		apiKeyID = id
placeholder

	// 获取时间范围参数
	now := timezone.Now()
	var startTime, endTime time.Time

	// 优先使用 start_date 和 end_date 参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr != "" && endDateStr != "" {
		// 使用自定义日期范围
		var err error
		startTime, err = timezone.ParseInLocation("2006-01-02", startDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid start_date format, use YYYY-MM-DD")
			return
	placeholder
		endTime, err = timezone.ParseInLocation("2006-01-02", endDateStr)
		if err != nil {
			response.BadRequest(c, "Invalid end_date format, use YYYY-MM-DD")
			return
	placeholder
		// 设置结束时间为当天结束
		endTime = endTime.Add(24*time.Hour - time.Nanosecond)
placeholder else {
		// 使用 period 参数
		period := c.DefaultQuery("period", "today")
		switch period {
		case "today":
			startTime = timezone.StartOfDay(now)
		case "week":
			startTime = now.AddDate(0, 0, -7)
		case "month":
			startTime = now.AddDate(0, -1, 0)
		default:
			startTime = timezone.StartOfDay(now)
	placeholder
		endTime = now
placeholder

	var stats *service.UsageStats
	var err error
	if apiKeyID > 0 {
		stats, err = h.usageService.GetStatsByApiKey(c.Request.Context(), apiKeyID, startTime, endTime)
placeholder else {
		stats, err = h.usageService.GetStatsByUser(c.Request.Context(), subject.UserID, startTime, endTime)
placeholder
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, stats)
placeholder

// parseUserTimeRange parses start_date, end_date query parameters for user dashboard
func parseUserTimeRange(c *gin.Context) (time.Time, time.Time) {
	now := timezone.Now()
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var startTime, endTime time.Time

	if startDate != "" {
		if t, err := timezone.ParseInLocation("2006-01-02", startDate); err == nil {
			startTime = t
	placeholder else {
			startTime = timezone.StartOfDay(now.AddDate(0, 0, -7))
	placeholder
placeholder else {
		startTime = timezone.StartOfDay(now.AddDate(0, 0, -7))
placeholder

	if endDate != "" {
		if t, err := timezone.ParseInLocation("2006-01-02", endDate); err == nil {
			endTime = t.Add(24 * time.Hour) // Include the end date
	placeholder else {
			endTime = timezone.StartOfDay(now.AddDate(0, 0, 1))
	placeholder
placeholder else {
		endTime = timezone.StartOfDay(now.AddDate(0, 0, 1))
placeholder

	return startTime, endTime
placeholder

// DashboardStats handles getting user dashboard statistics
// GET /api/v1/usage/dashboard/stats
func (h *UsageHandler) DashboardStats(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	stats, err := h.usageService.GetUserDashboardStats(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, stats)
placeholder

// DashboardTrend handles getting user usage trend data
// GET /api/v1/usage/dashboard/trend
func (h *UsageHandler) DashboardTrend(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	startTime, endTime := parseUserTimeRange(c)
	granularity := c.DefaultQuery("granularity", "day")

	trend, err := h.usageService.GetUserUsageTrendByUserID(c.Request.Context(), subject.UserID, startTime, endTime, granularity)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{
		"trend":       trend,
		"start_date":  startTime.Format("2006-01-02"),
		"end_date":    endTime.Add(-24 * time.Hour).Format("2006-01-02"),
		"granularity": granularity,
placeholder)
placeholder

// DashboardModels handles getting user model usage statistics
// GET /api/v1/usage/dashboard/models
func (h *UsageHandler) DashboardModels(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	startTime, endTime := parseUserTimeRange(c)

	stats, err := h.usageService.GetUserModelStats(c.Request.Context(), subject.UserID, startTime, endTime)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{
		"models":     stats,
		"start_date": startTime.Format("2006-01-02"),
		"end_date":   endTime.Add(-24 * time.Hour).Format("2006-01-02"),
placeholder)
placeholder

// BatchApiKeysUsageRequest represents the request for batch API keys usage
type BatchApiKeysUsageRequest struct {
	ApiKeyIDs []int64 `json:"api_key_ids" binding:"required"`
placeholder

// DashboardApiKeysUsage handles getting usage stats for user's own API keys
// POST /api/v1/usage/dashboard/api-keys-usage
func (h *UsageHandler) DashboardApiKeysUsage(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	var req BatchApiKeysUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	if len(req.ApiKeyIDs) == 0 {
		response.Success(c, gin.H{"stats": map[string]any{placeholderplaceholder)
		return
placeholder

	// Limit the number of API key IDs to prevent SQL parameter overflow
	if len(req.ApiKeyIDs) > 100 {
		response.BadRequest(c, "Too many API key IDs (maximum 100 allowed)")
		return
placeholder

	validApiKeyIDs, err := h.apiKeyService.VerifyOwnership(c.Request.Context(), subject.UserID, req.ApiKeyIDs)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	if len(validApiKeyIDs) == 0 {
		response.Success(c, gin.H{"stats": map[string]any{placeholderplaceholder)
		return
placeholder

	stats, err := h.usageService.GetBatchApiKeyUsageStats(c.Request.Context(), validApiKeyIDs)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"stats": statsplaceholder)
placeholder
