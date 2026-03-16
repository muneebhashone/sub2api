package admin

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	backupService *service.BackupService
	userService   *service.UserService
placeholder

func NewBackupHandler(backupService *service.BackupService, userService *service.UserService) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
		userService:   userService,
placeholder
placeholder

// ─── S3 配置 ───

func (h *BackupHandler) GetS3Config(c *gin.Context) {
	cfg, err := h.backupService.GetS3Config(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, cfg)
placeholder

func (h *BackupHandler) UpdateS3Config(c *gin.Context) {
	var req service.BackupS3Config
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	cfg, err := h.backupService.UpdateS3Config(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, cfg)
placeholder

func (h *BackupHandler) TestS3Connection(c *gin.Context) {
	var req service.BackupS3Config
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	err := h.backupService.TestS3Connection(c.Request.Context(), req)
	if err != nil {
		response.Success(c, gin.H{"ok": false, "message": err.Error()placeholder)
		return
placeholder
	response.Success(c, gin.H{"ok": true, "message": "connection successful"placeholder)
placeholder

// ─── 定时备份 ───

func (h *BackupHandler) GetSchedule(c *gin.Context) {
	cfg, err := h.backupService.GetSchedule(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, cfg)
placeholder

func (h *BackupHandler) UpdateSchedule(c *gin.Context) {
	var req service.BackupScheduleConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
placeholder
	cfg, err := h.backupService.UpdateSchedule(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, cfg)
placeholder

// ─── 备份操作 ───

type CreateBackupRequest struct {
	ExpireDays *int `json:"expire_days"` // nil=使用默认值14，0=永不过期
placeholder

func (h *BackupHandler) CreateBackup(c *gin.Context) {
	var req CreateBackupRequest
	_ = c.ShouldBindJSON(&req) // 允许空 body

	expireDays := 14 // 默认14天过期
	if req.ExpireDays != nil {
		expireDays = *req.ExpireDays
placeholder

	record, err := h.backupService.StartBackup(c.Request.Context(), "manual", expireDays)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Accepted(c, record)
placeholder

func (h *BackupHandler) ListBackups(c *gin.Context) {
	records, err := h.backupService.ListBackups(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if records == nil {
		records = []service.BackupRecord{placeholder
placeholder
	response.Success(c, gin.H{"items": recordsplaceholder)
placeholder

func (h *BackupHandler) GetBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		response.BadRequest(c, "backup ID is required")
		return
placeholder
	record, err := h.backupService.GetBackupRecord(c.Request.Context(), backupID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, record)
placeholder

func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		response.BadRequest(c, "backup ID is required")
		return
placeholder
	if err := h.backupService.DeleteBackup(c.Request.Context(), backupID); err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, gin.H{"deleted": trueplaceholder)
placeholder

func (h *BackupHandler) GetDownloadURL(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		response.BadRequest(c, "backup ID is required")
		return
placeholder
	url, err := h.backupService.GetBackupDownloadURL(c.Request.Context(), backupID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Success(c, gin.H{"url": urlplaceholder)
placeholder

// ─── 恢复操作（需要重新输入管理员密码） ───

type RestoreBackupRequest struct {
	Password string `json:"password" binding:"required"`
placeholder

func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	backupID := c.Param("id")
	if backupID == "" {
		response.BadRequest(c, "backup ID is required")
		return
placeholder

	var req RestoreBackupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "password is required for restore operation")
		return
placeholder

	// 从上下文获取当前管理员用户 ID
	sub, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "unauthorized")
		return
placeholder

	// 获取管理员用户并验证密码
	user, err := h.userService.GetByID(c.Request.Context(), sub.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	if !user.CheckPassword(req.Password) {
		response.BadRequest(c, "incorrect admin password")
		return
placeholder

	record, err := h.backupService.StartRestore(c.Request.Context(), backupID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
placeholder
	response.Accepted(c, record)
placeholder
