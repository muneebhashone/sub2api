package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/googleapi"
	"github.com/Wei-Shaw/sub2api/internal/securityaudit"
	"github.com/Wei-Shaw/sub2api/internal/service"
	coderws "github.com/coder/websocket"
	"github.com/gin-gonic/gin"
)

func (h *OpenAIGatewayHandler) openAISecurityAuditError(c *gin.Context, decision *securityaudit.Decision) {
	if decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		h.errorResponse(c, securityAuditStatus(decision), securityAuditErrorCode(decision), securityAuditMessage(decision))
		return
placeholder
	errType := "api_error"
	if decision.Kind == securityaudit.DecisionBlock {
		errType = "permission_error"
placeholder
	c.JSON(securityAuditStatus(decision), gin.H{"error": gin.H{
		"type": errType, "code": securityAuditErrorCode(decision), "message": securityAuditMessage(decision),
placeholderplaceholder)
placeholder

func (h *GatewayHandler) openAISecurityAuditError(c *gin.Context, decision *securityaudit.Decision) {
	if decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		h.chatCompletionsErrorResponse(c, securityAuditStatus(decision), securityAuditErrorCode(decision), securityAuditMessage(decision))
		return
placeholder
	errType := "api_error"
	if decision.Kind == securityaudit.DecisionBlock {
		errType = "permission_error"
placeholder
	c.JSON(securityAuditStatus(decision), gin.H{"error": gin.H{
		"type": errType, "code": securityAuditErrorCode(decision), "message": securityAuditMessage(decision),
placeholderplaceholder)
placeholder

func (h *GatewayHandler) responsesSecurityAuditError(c *gin.Context, decision *securityaudit.Decision) {
	if decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		h.responsesErrorResponse(c, securityAuditStatus(decision), securityAuditErrorCode(decision), securityAuditMessage(decision))
		return
placeholder
	c.JSON(securityAuditStatus(decision), gin.H{"error": gin.H{
		"type": "api_error", "code": securityAuditErrorCode(decision), "message": securityAuditMessage(decision),
placeholderplaceholder)
placeholder

func (h *GatewayHandler) anthropicSecurityAuditError(c *gin.Context, decision *securityaudit.Decision) {
	if decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		h.errorResponse(c, securityAuditStatus(decision), securityAuditErrorCode(decision), securityAuditMessage(decision))
		return
placeholder
	errType := "api_error"
	if decision.Kind == securityaudit.DecisionBlock {
		errType = "permission_error"
placeholder
	c.JSON(securityAuditStatus(decision), gin.H{"type": "error", "error": gin.H{
		"type": errType, "code": securityAuditErrorCode(decision), "message": securityAuditMessage(decision),
placeholderplaceholder)
placeholder

func (h *OpenAIGatewayHandler) anthropicSecurityAuditError(c *gin.Context, decision *securityaudit.Decision) {
	if decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		h.anthropicErrorResponse(c, securityAuditStatus(decision), securityAuditErrorCode(decision), securityAuditMessage(decision))
		return
placeholder
	errType := "api_error"
	if decision.Kind == securityaudit.DecisionBlock {
		errType = "permission_error"
placeholder
	c.JSON(securityAuditStatus(decision), gin.H{"type": "error", "error": gin.H{
		"type": errType, "code": securityAuditErrorCode(decision), "message": securityAuditMessage(decision),
placeholderplaceholder)
placeholder

func googleSecurityAuditError(c *gin.Context, decision *securityaudit.Decision) {
	if decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		googleError(c, securityAuditStatus(decision), securityAuditMessage(decision))
		return
placeholder
	status := securityAuditStatus(decision)
	googleStatus := googleapi.HTTPStatusToGoogleStatus(status)
	if status == http.StatusServiceUnavailable {
		googleStatus = "UNAVAILABLE"
placeholder
	requestID := ""
	if c != nil && c.Request != nil {
		requestID = contentModerationRequestID(c.Request.Context())
placeholder
	c.JSON(status, gin.H{"error": gin.H{
		"code": status, "message": securityAuditMessage(decision), "status": googleStatus,
		"details": []gin.H{{
			"@type":  "type.googleapis.com/google.rpc.ErrorInfo",
			"reason": securityAuditErrorCode(decision), "domain": "sub2api.securityaudit",
			"metadata": gin.H{"request_id": requestIDplaceholder,
placeholder
placeholderplaceholder)
placeholder

func writeSecurityAuditWSError(ctx context.Context, conn *coderws.Conn, decision *securityaudit.Decision) {
	if conn == nil || decision == nil {
		return
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		legacy := decision.Legacy
		writeContentModerationWSError(ctx, conn, (legacyContentModerationDecision{legacyplaceholder).toService())
		return
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder
	payload, err := json.Marshal(gin.H{
		"event_id": "evt_prompt_guard_rejected", "type": "error",
		"error": gin.H{"type": "invalid_request_error", "code": securityAuditErrorCode(decision), "message": securityAuditMessage(decision)placeholder,
placeholder)
	if err != nil {
		return
placeholder
	writeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_ = conn.Write(writeCtx, coderws.MessageText, payload)
placeholder

type legacyContentModerationDecision struct{ value *securityaudit.LegacyDecision placeholder

func (d legacyContentModerationDecision) toService() *service.ContentModerationDecision {
	if d.value == nil {
		return nil
placeholder
	return &service.ContentModerationDecision{Allowed: d.value.Allowed, Blocked: d.value.Blocked, Flagged: d.value.Flagged, Message: d.value.Message, StatusCode: d.value.StatusCode, Action: d.value.Actionplaceholder
placeholder

func securityAuditWSCloseStatus(decision *securityaudit.Decision) coderws.StatusCode {
	if decision == nil {
		return coderws.StatusInternalError
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		return coderws.StatusPolicyViolation
placeholder
	if decision.Kind == securityaudit.DecisionBlock {
		return coderws.StatusCode(4403)
placeholder
	return coderws.StatusTryAgainLater
placeholder

func securityAuditWSCloseReason(decision *securityaudit.Decision) string {
	if decision == nil {
		return securityaudit.ErrorCodeUnavailable
placeholder
	if decision.Legacy != nil && decision.Legacy.Blocked {
		message := strings.TrimSpace(decision.Legacy.Message)
		if message != "" {
			return message
	placeholder
		return "content_policy_violation"
placeholder
	code := securityAuditErrorCode(decision)
	if code == "" {
		return securityaudit.ErrorCodeUnavailable
placeholder
	return code
placeholder
