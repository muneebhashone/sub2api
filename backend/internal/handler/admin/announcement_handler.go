package admin

import (
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// AnnouncementHandler handles admin announcement management
type AnnouncementHandler struct {
	announcementService *service.AnnouncementService
placeholder

// NewAnnouncementHandler creates a new admin announcement handler
func NewAnnouncementHandler(announcementService *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementService: announcementService,
placeholder
placeholder

type CreateAnnouncementRequest struct {
	Title     string                       `json:"title" binding:"required"`
	Content   string                       `json:"content" binding:"required"`
	Status    string                       `json:"status" binding:"omitempty,oneof=draft active archived"`
	Targeting service.AnnouncementTargeting `json:"targeting"`
	StartsAt  *int64                       `json:"starts_at"` // Unix seconds, 0/empty = immediate
	EndsAt    *int64                       `json:"ends_at"`   // Unix seconds, 0/empty = never
placeholder

type UpdateAnnouncementRequest struct {
	Title     *string                        `json:"title"`
	Content   *string                        `json:"content"`
	Status    *string                        `json:"status" binding:"omitempty,oneof=draft active archived"`
	Targeting *service.AnnouncementTargeting `json:"targeting"`
	StartsAt  *int64                         `json:"starts_at"` // Unix seconds, 0 = clear
	EndsAt    *int64                         `json:"ends_at"`   // Unix seconds, 0 = clear
placeholder

// List handles listing announcements with filters
// GET /api/v1/admin/announcements
func (h *AnnouncementHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	status := strings.TrimSpace(c.Query("status"))
	search := strings.TrimSpace(c.Query("search"))
	if len(search) > 200 {
		search = search[:200]
placeholder

	params := pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
placeholder

	items, paginationResult, err := h.announcementService.List(
		c.Request.Context(),
		params,
		service.AnnouncementListFilters{Status: status, Search: searchplaceholder,
	)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	out := make([]dto.Announcement, 0, len(items))
	for i := range items {
		out = append(out, *dto.AnnouncementFromService(&items[i]))
placeholder
	response.Paginated(c, out, paginationResult.Total, page, pageSize)
placeholder

// GetByID handles getting an announcement by ID
// GET /api/v1/admin/announcements/:id
func (h *AnnouncementHandler) GetByID(c *gin.Context) {
	announcementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || announcementID <= 0 {
		response.BadRequest(c, "Invalid announcement ID")
		return
placeholder

	item, err := h.announcementService.GetByID(c.Request.Context(), announcementID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.AnnouncementFromService(item))
placeholder

// Create handles creating a new announcement
// POST /api/v1/admin/announcements
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
placeholder

	input := &service.CreateAnnouncementInput{
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
		Targeting: req.Targeting,
		ActorID:   &subject.UserID,
placeholder

	if req.StartsAt != nil && *req.StartsAt > 0 {
		t := time.Unix(*req.StartsAt, 0)
		input.StartsAt = &t
placeholder
	if req.EndsAt != nil && *req.EndsAt > 0 {
		t := time.Unix(*req.EndsAt, 0)
		input.EndsAt = &t
placeholder

	created, err := h.announcementService.Create(c.Request.Context(), input)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.AnnouncementFromService(created))
placeholder

// Update handles updating an announcement
// PUT /api/v1/admin/announcements/:id
func (h *AnnouncementHandler) Update(c *gin.Context) {
	announcementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || announcementID <= 0 {
		response.BadRequest(c, "Invalid announcement ID")
		return
placeholder

	var req UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
placeholder

	input := &service.UpdateAnnouncementInput{
		Title:     req.Title,
		Content:   req.Content,
		Status:    req.Status,
		Targeting: req.Targeting,
		ActorID:   &subject.UserID,
placeholder

	if req.StartsAt != nil {
		if *req.StartsAt == 0 {
			var cleared *time.Time = nil
			input.StartsAt = &cleared
	placeholder else {
			t := time.Unix(*req.StartsAt, 0)
			ptr := &t
			input.StartsAt = &ptr
	placeholder
placeholder

	if req.EndsAt != nil {
		if *req.EndsAt == 0 {
			var cleared *time.Time = nil
			input.EndsAt = &cleared
	placeholder else {
			t := time.Unix(*req.EndsAt, 0)
			ptr := &t
			input.EndsAt = &ptr
	placeholder
placeholder

	updated, err := h.announcementService.Update(c.Request.Context(), announcementID, input)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.AnnouncementFromService(updated))
placeholder

// Delete handles deleting an announcement
// DELETE /api/v1/admin/announcements/:id
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	announcementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || announcementID <= 0 {
		response.BadRequest(c, "Invalid announcement ID")
		return
placeholder

	if err := h.announcementService.Delete(c.Request.Context(), announcementID); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "Announcement deleted successfully"placeholder)
placeholder

// ListReadStatus handles listing users read status for an announcement
// GET /api/v1/admin/announcements/:id/read-status
func (h *AnnouncementHandler) ListReadStatus(c *gin.Context) {
	announcementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || announcementID <= 0 {
		response.BadRequest(c, "Invalid announcement ID")
		return
placeholder

	page, pageSize := response.ParsePagination(c)
	params := pagination.PaginationParams{
		Page:     page,
		PageSize: pageSize,
placeholder
	search := strings.TrimSpace(c.Query("search"))
	if len(search) > 200 {
		search = search[:200]
placeholder

	items, paginationResult, err := h.announcementService.ListUserReadStatus(
		c.Request.Context(),
		announcementID,
		params,
		search,
	)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Paginated(c, items, paginationResult.Total, page, pageSize)
placeholder

