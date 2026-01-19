package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userService *service.UserService
placeholder

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
placeholder
placeholder

// ChangePasswordRequest represents the change password request payload
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
placeholder

// UpdateProfileRequest represents the update profile request payload
type UpdateProfileRequest struct {
	Username *string `json:"username"`
placeholder

// GetProfile handles getting user profile
// GET /api/v1/users/me
func (h *UserHandler) GetProfile(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	userData, err := h.userService.GetByID(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.UserFromService(userData))
placeholder

// ChangePassword handles changing user password
// POST /api/v1/users/me/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	svcReq := service.ChangePasswordRequest{
		CurrentPassword: req.OldPassword,
		NewPassword:     req.NewPassword,
placeholder
	err := h.userService.ChangePassword(c.Request.Context(), subject.UserID, svcReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, gin.H{"message": "Password changed successfully"placeholder)
placeholder

// UpdateProfile handles updating user profile
// PUT /api/v1/users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	svcReq := service.UpdateProfileRequest{
		Username: req.Username,
placeholder
	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), subject.UserID, svcReq)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder

	response.Success(c, dto.UserFromService(updatedUser))
placeholder
