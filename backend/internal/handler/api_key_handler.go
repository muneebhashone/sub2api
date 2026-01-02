package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// APIKeyHandler handles API key-related requests
type APIKeyHandler struct {
	apiKeyService *service.APIKeyService
placeholder

// NewAPIKeyHandler creates a new APIKeyHandler
func NewAPIKeyHandler(apiKeyService *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
placeholder
placeholder

// CreateAPIKeyRequest represents the create API key request payload
type CreateAPIKeyRequest struct {
	Name      string  `json:"name" binding:"required"`
	GroupID   *int64  `json:"group_id"`   // nullable
	CustomKey *string `json:"custom_key"` // 可选的自定义key
placeholder

// UpdateAPIKeyRequest represents the update API key request payload
type UpdateAPIKeyRequest struct {
	Name    string `json:"name"`
	GroupID *int64 `json:"group_id"`
	Status  string `json:"status" binding:"omitempty,oneof=active inactive"`
placeholder

// List handles listing user's API keys with pagination
// GET /api/v1/api-keys
func (h *APIKeyHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{Page: page, PageSize: pageSizeplaceholder

	keys, result, err := h.apiKeyService.List(c.Request.Context(), subject.UserID, params)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	out := make([]dto.APIKey, 0, len(keys))
	for i := range keys {
		out = append(out, *dto.APIKeyFromService(&keys[i]))
placeholder
	response.Paginated(c, out, result.Total, page, pageSize)
placeholder

// GetByID handles getting a single API key
// GET /api/v1/api-keys/:id
func (h *APIKeyHandler) GetByID(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid key ID")
		return
placeholder

	key, err := h.apiKeyService.GetByID(c.Request.Context(), keyID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	// 验证所有权
	if key.UserID != subject.UserID {
		response.Forbidden(c, "Not authorized to access this key")
		return
placeholder

	response.Success(c, dto.APIKeyFromService(key))
placeholder

// Create handles creating a new API key
// POST /api/v1/api-keys
func (h *APIKeyHandler) Create(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	svcReq := service.CreateAPIKeyRequest{
		Name:      req.Name,
		GroupID:   req.GroupID,
		CustomKey: req.CustomKey,
placeholder
	key, err := h.apiKeyService.Create(c.Request.Context(), subject.UserID, svcReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.APIKeyFromService(key))
placeholder

// Update handles updating an API key
// PUT /api/v1/api-keys/:id
func (h *APIKeyHandler) Update(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid key ID")
		return
placeholder

	var req UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	svcReq := service.UpdateAPIKeyRequest{placeholder
	if req.Name != "" {
		svcReq.Name = &req.Name
placeholder
	svcReq.GroupID = req.GroupID
	if req.Status != "" {
		svcReq.Status = &req.Status
placeholder

	key, err := h.apiKeyService.Update(c.Request.Context(), keyID, subject.UserID, svcReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.APIKeyFromService(key))
placeholder

// Delete handles deleting an API key
// DELETE /api/v1/api-keys/:id
func (h *APIKeyHandler) Delete(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid key ID")
		return
placeholder

	err = h.apiKeyService.Delete(c.Request.Context(), keyID, subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "API key deleted successfully"placeholder)
placeholder

// GetAvailableGroups 获取用户可以绑定的分组列表
// GET /api/v1/groups/available
func (h *APIKeyHandler) GetAvailableGroups(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	groups, err := h.apiKeyService.GetAvailableGroups(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	out := make([]dto.Group, 0, len(groups))
	for i := range groups {
		out = append(out, *dto.GroupFromService(&groups[i]))
placeholder
	response.Success(c, out)
placeholder
