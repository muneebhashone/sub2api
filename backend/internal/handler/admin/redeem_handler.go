package admin

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RedeemHandler handles admin redeem code management
type RedeemHandler struct {
	adminService  service.AdminService
	redeemService *service.RedeemService
placeholder

// NewRedeemHandler creates a new admin redeem handler
func NewRedeemHandler(adminService service.AdminService, redeemService *service.RedeemService) *RedeemHandler {
	return &RedeemHandler{
		adminService:  adminService,
		redeemService: redeemService,
placeholder
placeholder

// GenerateRedeemCodesRequest represents generate redeem codes request
type GenerateRedeemCodesRequest struct {
	Count         int        `json:"count" binding:"required,min=1,max=100"`
	Type          string     `json:"type" binding:"required,oneof=balance concurrency subscription invitation"`
	Value         float64    `json:"value"`
	GroupID       *int64     `json:"group_id"`      // 订阅类型必填
	ValidityDays  int        `json:"validity_days"` // 订阅类型使用，正数增加/负数退款扣减
	ExpiresAt     *time.Time `json:"expires_at"`
	ExpiresInDays *int       `json:"expires_in_days" binding:"omitempty,min=1,max=3650"`
placeholder

// CreateAndRedeemCodeRequest represents creating a fixed code and redeeming it for a target user.
// Type 为 omitempty 而非 required 是为了向后兼容旧版调用方（不传 type 时默认 balance）。
type CreateAndRedeemCodeRequest struct {
	Code          string     `json:"code" binding:"required,min=3,max=128"`
	Type          string     `json:"type" binding:"omitempty,oneof=balance concurrency subscription invitation"` // 不传时默认 balance（向后兼容）
	Value         float64    `json:"value" binding:"required"`
	UserID        int64      `json:"user_id" binding:"required,gt=0"`
	GroupID       *int64     `json:"group_id"`      // subscription 类型必填
	ValidityDays  int        `json:"validity_days"` // subscription 类型：正数增加，负数退款扣减
	Notes         string     `json:"notes"`
	ExpiresAt     *time.Time `json:"expires_at"`
	ExpiresInDays *int       `json:"expires_in_days" binding:"omitempty,min=1,max=3650"`
placeholder

func resolveRedeemCodeExpiresAt(expiresAt *time.Time, expiresInDays *int) (*time.Time, error) {
	if expiresAt != nil && expiresInDays != nil {
		return nil, infraerrors.BadRequest("REDEEM_CODE_EXPIRY_CONFLICT", "expires_at and expires_in_days cannot both be set")
placeholder

	now := time.Now().UTC()
	if expiresInDays != nil {
		if *expiresInDays <= 0 {
			return nil, infraerrors.BadRequest("REDEEM_CODE_EXPIRES_IN_DAYS_INVALID", "expires_in_days must be greater than zero")
	placeholder
		expires := now.AddDate(0, 0, *expiresInDays)
		return &expires, nil
placeholder
	if expiresAt == nil {
		return nil, nil
placeholder

	expires := expiresAt.UTC()
	if !expires.After(now) {
		return nil, infraerrors.BadRequest("REDEEM_CODE_EXPIRES_AT_INVALID", "expires_at must be in the future")
placeholder
	return &expires, nil
placeholder

// List handles listing all redeem codes with pagination
// GET /api/v1/admin/redeem-codes
func (h *RedeemHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	codeType := c.Query("type")
	status := c.Query("status")
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	// 标准化和验证 search 参数
	search = strings.TrimSpace(search)
	if len(search) > 100 {
		search = search[:100]
placeholder

	codes, total, err := h.adminService.ListRedeemCodes(c.Request.Context(), page, pageSize, codeType, status, search, sortBy, sortOrder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	out := make([]dto.AdminRedeemCode, 0, len(codes))
	for i := range codes {
		out = append(out, *dto.RedeemCodeFromServiceAdmin(&codes[i]))
placeholder
	response.Paginated(c, out, total, page, pageSize)
placeholder

// GetByID handles getting a redeem code by ID
// GET /api/v1/admin/redeem-codes/:id
func (h *RedeemHandler) GetByID(c *gin.Context) {
	codeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid redeem code ID")
		return
placeholder

	code, err := h.adminService.GetRedeemCode(c.Request.Context(), codeID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.RedeemCodeFromServiceAdmin(code))
placeholder

// Generate handles generating new redeem codes
// POST /api/v1/admin/redeem-codes/generate
func (h *RedeemHandler) Generate(c *gin.Context) {
	var req GenerateRedeemCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	expiresAt, err := resolveRedeemCodeExpiresAt(req.ExpiresAt, req.ExpiresInDays)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	executeAdminIdempotentJSON(c, "admin.redeem_codes.generate", req, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		codes, execErr := h.adminService.GenerateRedeemCodes(ctx, &service.GenerateRedeemCodesInput{
			Count:        req.Count,
			Type:         req.Type,
			Value:        req.Value,
			GroupID:      req.GroupID,
			ValidityDays: req.ValidityDays,
			ExpiresAt:    expiresAt,
	placeholder)
		if execErr != nil {
			return nil, execErr
	placeholder

		out := make([]dto.AdminRedeemCode, 0, len(codes))
		for i := range codes {
			out = append(out, *dto.RedeemCodeFromServiceAdmin(&codes[i]))
	placeholder
		return out, nil
placeholder)
placeholder

// CreateAndRedeem creates a fixed redeem code and redeems it for a target user in one step.
// POST /api/v1/admin/redeem-codes/create-and-redeem
func (h *RedeemHandler) CreateAndRedeem(c *gin.Context) {
	if h.redeemService == nil {
		response.InternalError(c, "redeem service not configured")
		return
placeholder

	var req CreateAndRedeemCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	req.Code = strings.TrimSpace(req.Code)
	// 向后兼容：旧版调用方（如 Sub2ApiPay）不传 type 字段，默认当作 balance 充值处理。
	// 请勿删除此默认值逻辑，否则会导致旧版调用方 400 报错。
	if req.Type == "" {
		req.Type = "balance"
placeholder

	if req.Type == "subscription" {
		if req.GroupID == nil {
			response.BadRequest(c, "group_id is required for subscription type")
			return
	placeholder
		if req.ValidityDays == 0 {
			response.BadRequest(c, "validity_days must not be zero for subscription type")
			return
	placeholder
placeholder

	expiresAt, err := resolveRedeemCodeExpiresAt(req.ExpiresAt, req.ExpiresInDays)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	executeAdminIdempotentJSON(c, "admin.redeem_codes.create_and_redeem", req, service.DefaultWriteIdempotencyTTL(), func(ctx context.Context) (any, error) {
		existing, err := h.redeemService.GetByCode(ctx, req.Code)
		if err == nil {
			return h.resolveCreateAndRedeemExisting(ctx, existing, req.UserID)
	placeholder
		if !errors.Is(err, service.ErrRedeemCodeNotFound) {
			return nil, err
	placeholder

		createErr := h.redeemService.CreateCode(ctx, &service.RedeemCode{
			Code:         req.Code,
			Type:         req.Type,
			Value:        req.Value,
			Status:       service.StatusUnused,
			Notes:        req.Notes,
			GroupID:      req.GroupID,
			ValidityDays: req.ValidityDays,
			ExpiresAt:    expiresAt,
	placeholder)
		if createErr != nil {
			// Unique code race: if code now exists, use idempotent semantics by used_by.
			existingAfterCreateErr, getErr := h.redeemService.GetByCode(ctx, req.Code)
			if getErr == nil {
				return h.resolveCreateAndRedeemExisting(ctx, existingAfterCreateErr, req.UserID)
		placeholder
			return nil, createErr
	placeholder

		redeemed, redeemErr := h.redeemService.Redeem(ctx, req.UserID, req.Code)
		if redeemErr != nil {
			return nil, redeemErr
	placeholder
		return gin.H{"redeem_code": dto.RedeemCodeFromServiceAdmin(redeemed)placeholder, nil
placeholder)
placeholder

func (h *RedeemHandler) resolveCreateAndRedeemExisting(ctx context.Context, existing *service.RedeemCode, userID int64) (any, error) {
	if existing == nil {
		return nil, infraerrors.Conflict("REDEEM_CODE_CONFLICT", "redeem code conflict")
placeholder

	// If previous run created the code but crashed before redeem, redeem it now.
	if existing.IsExpired() {
		return nil, service.ErrRedeemCodeExpired
placeholder
	if existing.CanUse() {
		redeemed, err := h.redeemService.Redeem(ctx, userID, existing.Code)
		if err == nil {
			return gin.H{"redeem_code": dto.RedeemCodeFromServiceAdmin(redeemed)placeholder, nil
	placeholder
		if !errors.Is(err, service.ErrRedeemCodeUsed) {
			return nil, err
	placeholder
		latest, getErr := h.redeemService.GetByCode(ctx, existing.Code)
		if getErr == nil {
			existing = latest
	placeholder
placeholder

	if existing.UsedBy != nil && *existing.UsedBy == userID {
		return gin.H{"redeem_code": dto.RedeemCodeFromServiceAdmin(existing)placeholder, nil
placeholder

	return nil, infraerrors.Conflict("REDEEM_CODE_CONFLICT", "redeem code already used by another user")
placeholder

// Delete handles deleting a redeem code
// DELETE /api/v1/admin/redeem-codes/:id
func (h *RedeemHandler) Delete(c *gin.Context) {
	codeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid redeem code ID")
		return
placeholder

	err = h.adminService.DeleteRedeemCode(c.Request.Context(), codeID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "Redeem code deleted successfully"placeholder)
placeholder

// BatchDelete handles batch deleting redeem codes
// POST /api/v1/admin/redeem-codes/batch-delete
func (h *RedeemHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids" binding:"required,min=1"`
placeholder
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	deleted, err := h.adminService.BatchDeleteRedeemCodes(c.Request.Context(), req.IDs)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{
		"deleted": deleted,
		"message": "Redeem codes deleted successfully",
placeholder)
placeholder

// BatchUpdate handles batch updating redeem codes
// POST /api/v1/admin/redeem-codes/batch-update
func (h *RedeemHandler) BatchUpdate(c *gin.Context) {
	if h.redeemService == nil {
		response.InternalError(c, "redeem service not configured")
		return
placeholder

	var req dto.BatchUpdateRedeemCodesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	result, err := h.redeemService.BatchUpdate(c.Request.Context(), &service.RedeemCodeBatchUpdateInput{
		IDs:    req.IDs,
		Fields: redeemBatchUpdateFieldsFromDTO(req.Fields),
placeholder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{
		"updated": result.Updated,
		"message": "Redeem codes updated successfully",
placeholder)
placeholder

func redeemBatchUpdateFieldsFromDTO(in dto.BatchUpdateRedeemCodeFields) service.RedeemCodeBatchUpdateFields {
	out := service.RedeemCodeBatchUpdateFields{
		Status: in.Status,
		Notes:  in.Notes,
		Type:   in.Type,
		Value:  in.Value,
placeholder
	if in.ExpiresAt.Set {
		out.ExpiresAt = service.NullableTimeUpdate{Set: true, Value: in.ExpiresAt.Valueplaceholder
placeholder
	if in.GroupID.Set {
		out.GroupID = service.NullableInt64Update{Set: true, Value: in.GroupID.Valueplaceholder
placeholder
	return out
placeholder

// Expire handles expiring a redeem code
// POST /api/v1/admin/redeem-codes/:id/expire
func (h *RedeemHandler) Expire(c *gin.Context) {
	codeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Invalid redeem code ID")
		return
placeholder

	code, err := h.adminService.ExpireRedeemCode(c.Request.Context(), codeID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.RedeemCodeFromServiceAdmin(code))
placeholder

// GetStats handles getting redeem code statistics
// GET /api/v1/admin/redeem-codes/stats
func (h *RedeemHandler) GetStats(c *gin.Context) {
	// Return mock data for now
	response.Success(c, gin.H{
		"total_codes":             0,
		"active_codes":            0,
		"used_codes":              0,
		"expired_codes":           0,
		"total_value_distributed": 0.0,
		"by_type": gin.H{
			"balance":     0,
			"concurrency": 0,
			"trial":       0,
	placeholder,
placeholder)
placeholder

// Export handles exporting redeem codes to CSV
// GET /api/v1/admin/redeem-codes/export
func (h *RedeemHandler) Export(c *gin.Context) {
	codeType := c.Query("type")
	status := c.Query("status")
	search := strings.TrimSpace(c.Query("search"))
	sortBy := c.DefaultQuery("sort_by", "id")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	if len(search) > 100 {
		search = search[:100]
placeholder

	// Get all codes without pagination (use large page size)
	codes, _, err := h.adminService.ListRedeemCodes(c.Request.Context(), 1, 10000, codeType, status, search, sortBy, sortOrder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	// Create CSV buffer
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	if err := writer.Write([]string{"id", "code", "type", "value", "status", "used_by", "used_by_email", "used_at", "expires_at", "created_at"placeholder); err != nil {
		response.InternalError(c, "Failed to export redeem codes: "+err.Error())
		return
placeholder

	// Write data rows
	for _, code := range codes {
		usedBy := ""
		if code.UsedBy != nil {
			usedBy = fmt.Sprintf("%d", *code.UsedBy)
	placeholder
		usedByEmail := ""
		if code.User != nil {
			usedByEmail = code.User.Email
	placeholder
		usedAt := ""
		if code.UsedAt != nil {
			usedAt = code.UsedAt.Format("2006-01-02 15:04:05")
	placeholder
		expiresAt := ""
		if code.ExpiresAt != nil {
			expiresAt = code.ExpiresAt.Format("2006-01-02 15:04:05")
	placeholder
		if err := writer.Write([]string{
			fmt.Sprintf("%d", code.ID),
			code.Code,
			code.Type,
			fmt.Sprintf("%.2f", code.Value),
			code.Status,
			usedBy,
			usedByEmail,
			usedAt,
			expiresAt,
			code.CreatedAt.Format("2006-01-02 15:04:05"),
	placeholder); err != nil {
			response.InternalError(c, "Failed to export redeem codes: "+err.Error())
			return
	placeholder
placeholder

	writer.Flush()
	if err := writer.Error(); err != nil {
		response.InternalError(c, "Failed to export redeem codes: "+err.Error())
		return
placeholder

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=redeem_codes.csv")
	c.Data(200, "text/csv", buf.Bytes())
placeholder
