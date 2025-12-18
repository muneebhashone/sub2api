package handler

import (
	"strconv"
	"time"

	"sub2api/internal/model"
	"sub2api/internal/pkg/response"
	"sub2api/internal/pkg/timezone"
	"sub2api/internal/repository"
	"sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// UsageHandler handles usage-related requests
type UsageHandler struct {
	usageService  *service.UsageService
	usageRepo     *repository.UsageLogRepository
	apiKeyService *service.ApiKeyService
placeholder

// NewUsageHandler creates a new UsageHandler
func NewUsageHandler(usageService *service.UsageService, usageRepo *repository.UsageLogRepository, apiKeyService *service.ApiKeyService) *UsageHandler {
	return &UsageHandler{
		usageService:  usageService,
		usageRepo:     usageRepo,
		apiKeyService: apiKeyService,
placeholder
placeholder

// List handles listing usage records with pagination
// GET /api/v1/usage
func (h *UsageHandler) List(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
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
			response.NotFound(c, "API key not found")
			return
	placeholder
		if apiKey.UserID != user.ID {
			response.Forbidden(c, "Not authorized to access this API key's usage records")
			return
	placeholder

		apiKeyID = id
placeholder

	params := repository.PaginationParams{Page: page, PageSize: pageSizeplaceholder
	var records []model.UsageLog
	var result *repository.PaginationResult
	var err error

	if apiKeyID > 0 {
		records, result, err = h.usageService.ListByApiKey(c.Request.Context(), apiKeyID, params)
placeholder else {
		records, result, err = h.usageService.ListByUser(c.Request.Context(), user.ID, params)
placeholder
	if err != nil {
		response.InternalError(c, "Failed to list usage records: "+err.Error())
		return
placeholder

	response.Paginated(c, records, result.Total, page, pageSize)
placeholder

// GetByID handles getting a single usage record
// GET /api/v1/usage/:id
func (h *UsageHandler) GetByID(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
		return
placeholder

	usageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid usage ID")
		return
placeholder

	record, err := h.usageService.GetByID(c.Request.Context(), usageID)
	if err != nil {
		response.NotFound(c, "Usage record not found")
		return
placeholder

	// 验证所有权
	if record.UserID != user.ID {
		response.Forbidden(c, "Not authorized to access this record")
		return
placeholder

	response.Success(c, record)
placeholder

// Stats handles getting usage statistics
// GET /api/v1/usage/stats
func (h *UsageHandler) Stats(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
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
		if apiKey.UserID != user.ID {
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
		stats, err = h.usageService.GetStatsByUser(c.Request.Context(), user.ID, startTime, endTime)
placeholder
	if err != nil {
		response.InternalError(c, "Failed to get usage statistics: "+err.Error())
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
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
		return
placeholder

	stats, err := h.usageRepo.GetUserDashboardStats(c.Request.Context(), user.ID)
	if err != nil {
		response.InternalError(c, "Failed to get dashboard statistics")
		return
placeholder

	response.Success(c, stats)
placeholder

// DashboardTrend handles getting user usage trend data
// GET /api/v1/usage/dashboard/trend
func (h *UsageHandler) DashboardTrend(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
		return
placeholder

	startTime, endTime := parseUserTimeRange(c)
	granularity := c.DefaultQuery("granularity", "day")

	trend, err := h.usageRepo.GetUserUsageTrendByUserID(c.Request.Context(), user.ID, startTime, endTime, granularity)
	if err != nil {
		response.InternalError(c, "Failed to get usage trend")
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
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
		return
placeholder

	startTime, endTime := parseUserTimeRange(c)

	stats, err := h.usageRepo.GetUserModelStats(c.Request.Context(), user.ID, startTime, endTime)
	if err != nil {
		response.InternalError(c, "Failed to get model statistics")
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
	userValue, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	user, ok := userValue.(*model.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
		return
placeholder

	var req BatchApiKeysUsageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	if len(req.ApiKeyIDs) == 0 {
		response.Success(c, gin.H{"stats": map[string]interface{placeholder{placeholderplaceholder)
		return
placeholder

	// Verify ownership of all requested API keys
	userApiKeys, _, err := h.apiKeyService.List(c.Request.Context(), user.ID, repository.PaginationParams{Page: 1, PageSize: 1000placeholder)
	if err != nil {
		response.InternalError(c, "Failed to verify API key ownership")
		return
placeholder

	userApiKeyIDs := make(map[int64]bool)
	for _, key := range userApiKeys {
		userApiKeyIDs[key.ID] = true
placeholder

	// Filter to only include user's own API keys
	validApiKeyIDs := make([]int64, 0)
	for _, id := range req.ApiKeyIDs {
		if userApiKeyIDs[id] {
			validApiKeyIDs = append(validApiKeyIDs, id)
	placeholder
placeholder

	if len(validApiKeyIDs) == 0 {
		response.Success(c, gin.H{"stats": map[string]interface{placeholder{placeholderplaceholder)
		return
placeholder

	stats, err := h.usageRepo.GetBatchApiKeyUsageStats(c.Request.Context(), validApiKeyIDs)
	if err != nil {
		response.InternalError(c, "Failed to get API key usage stats")
		return
placeholder

	response.Success(c, gin.H{"stats": statsplaceholder)
placeholder
