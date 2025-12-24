package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RedeemHandler handles redeem code-related requests
type RedeemHandler struct {
	redeemService *service.RedeemService
placeholder

// NewRedeemHandler creates a new RedeemHandler
func NewRedeemHandler(redeemService *service.RedeemService) *RedeemHandler {
	return &RedeemHandler{
		redeemService: redeemService,
placeholder
placeholder

// RedeemRequest represents the redeem code request payload
type RedeemRequest struct {
	Code string `json:"code" binding:"required"`
placeholder

// RedeemResponse represents the redeem response
type RedeemResponse struct {
	Message        string   `json:"message"`
	Type           string   `json:"type"`
	Value          float64  `json:"value"`
	NewBalance     *float64 `json:"new_balance,omitempty"`
	NewConcurrency *int     `json:"new_concurrency,omitempty"`
placeholder

// Redeem handles redeeming a code
// POST /api/v1/redeem
func (h *RedeemHandler) Redeem(c *gin.Context) {
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

	var req RedeemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	result, err := h.redeemService.Redeem(c.Request.Context(), user.ID, req.Code)
	if err != nil {
		response.BadRequest(c, "Failed to redeem code: "+err.Error())
		return
placeholder

	response.Success(c, result)
placeholder

// GetHistory returns the user's redemption history
// GET /api/v1/redeem/history
func (h *RedeemHandler) GetHistory(c *gin.Context) {
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

	// Default limit is 25
	limit := 25

	codes, err := h.redeemService.GetUserHistory(c.Request.Context(), user.ID, limit)
	if err != nil {
		response.InternalError(c, "Failed to get history: "+err.Error())
		return
placeholder

	response.Success(c, codes)
placeholder
