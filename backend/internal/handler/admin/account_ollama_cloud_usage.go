package admin

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type ollamaCloudUsageSessionRequest struct {
	Session string `json:"session" binding:"required"`
placeholder

type ollamaCloudUsageAutoRefreshRequest struct {
	Enabled *bool `json:"enabled" binding:"required"`
placeholder

func (h *AccountHandler) GetOllamaCloudUsageSettings(c *gin.Context) {
	if h.ollamaCloudUsage == nil {
		response.ErrorFrom(c, service.ErrOllamaCloudUsageUnavailable)
		return
placeholder
	settings, err := h.ollamaCloudUsage.GetSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, settings)
placeholder

func (h *AccountHandler) UpdateOllamaCloudUsageSettings(c *gin.Context) {
	if h.ollamaCloudUsage == nil {
		response.ErrorFrom(c, service.ErrOllamaCloudUsageUnavailable)
		return
placeholder
	var req service.OllamaCloudUsageSettings
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	if err := h.ollamaCloudUsage.UpdateSettings(c.Request.Context(), &req); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	settings, err := h.ollamaCloudUsage.GetSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, settings)
placeholder

func (h *AccountHandler) GetOllamaCloudUsage(c *gin.Context) {
	if !h.requireOllamaCloudUsage(c) {
		return
placeholder
	accountID, ok := ollamaCloudUsageAccountID(c)
	if !ok {
		return
placeholder
	state, err := h.ollamaCloudUsage.GetState(c.Request.Context(), accountID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, state)
placeholder

func (h *AccountHandler) SaveOllamaCloudUsageSession(c *gin.Context) {
	if !h.requireOllamaCloudUsage(c) {
		return
placeholder
	accountID, ok := ollamaCloudUsageAccountID(c)
	if !ok {
		return
placeholder
	var req ollamaCloudUsageSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	state, err := h.ollamaCloudUsage.SaveSession(c.Request.Context(), accountID, req.Session)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, state)
placeholder

func (h *AccountHandler) DeleteOllamaCloudUsageSession(c *gin.Context) {
	if !h.requireOllamaCloudUsage(c) {
		return
placeholder
	accountID, ok := ollamaCloudUsageAccountID(c)
	if !ok {
		return
placeholder
	state, err := h.ollamaCloudUsage.DeleteSession(c.Request.Context(), accountID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, state)
placeholder

func (h *AccountHandler) SetOllamaCloudUsageAutoRefresh(c *gin.Context) {
	if !h.requireOllamaCloudUsage(c) {
		return
placeholder
	accountID, ok := ollamaCloudUsageAccountID(c)
	if !ok {
		return
placeholder
	var req ollamaCloudUsageAutoRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	state, err := h.ollamaCloudUsage.SetAutoRefresh(c.Request.Context(), accountID, *req.Enabled)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, state)
placeholder

func (h *AccountHandler) RefreshOllamaCloudUsage(c *gin.Context) {
	if !h.requireOllamaCloudUsage(c) {
		return
placeholder
	accountID, ok := ollamaCloudUsageAccountID(c)
	if !ok {
		return
placeholder
	state, err := h.ollamaCloudUsage.Refresh(c.Request.Context(), accountID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, state)
placeholder

func (h *AccountHandler) requireOllamaCloudUsage(c *gin.Context) bool {
	if h != nil && h.ollamaCloudUsage != nil {
		return true
placeholder
	response.ErrorFrom(c, service.ErrOllamaCloudUsageUnavailable)
	return false
placeholder

func ollamaCloudUsageAccountID(c *gin.Context) (int64, bool) {
	if c == nil {
		return 0, false
placeholder
	accountID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || accountID <= 0 {
		response.BadRequest(c, "Invalid account ID")
		return 0, false
placeholder
	return accountID, true
placeholder
