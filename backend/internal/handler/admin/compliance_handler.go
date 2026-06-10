package admin

import (
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/ip"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

type ComplianceHandler struct {
	settingService *service.SettingService
placeholder

func NewComplianceHandler(settingService *service.SettingService) *ComplianceHandler {
	return &ComplianceHandler{settingService: settingServiceplaceholder
placeholder

type AcceptAdminComplianceRequest struct {
	Phrase   string `json:"phrase" binding:"required"`
	Language string `json:"language"`
placeholder

func (h *ComplianceHandler) GetStatus(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	status, err := h.settingService.GetAdminComplianceStatus(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, status)
placeholder

func (h *ComplianceHandler) Accept(c *gin.Context) {
	var req AcceptAdminComplianceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder

	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not authenticated")
		return
placeholder

	status, err := h.settingService.AcceptAdminCompliance(c.Request.Context(), service.AdminComplianceAcceptInput{
		AdminUserID: subject.UserID,
		Phrase:      req.Phrase,
		Language:    req.Language,
		IPAddress:   ip.GetClientIP(c),
		UserAgent:   strings.TrimSpace(c.GetHeader("User-Agent")),
placeholder)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, status)
placeholder
