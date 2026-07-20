package admin

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/pkg/sysutil"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles system-related operations
type SystemHandler struct {
	updateSvc systemUpdateService
	lockSvc   *service.SystemOperationLockService
placeholder

// systemUpdateTimeout bounds a full in-place update or rollback: the release
// manifest fetch plus a large binary download over slow links. It must stay
// above the GitHub download client timeout (10 minutes) so the download owns
// its own deadline.
const systemUpdateTimeout = 15 * time.Minute

// systemUpdateContext detaches a long-running update/rollback from the HTTP
// request lifetime. Browsers and reverse proxies commonly abort idle requests
// after 30-60s (axios default, nginx proxy_read_timeout), which canceled
// c.Request.Context() mid-download and killed the update with
// "download failed: context canceled" (#4504). The swap keeps running after a
// client disconnect; a later retry then hits the system operation lock or
// reports "Already up to date".
func systemUpdateContext(ctx context.Context) (context.Context, context.CancelFunc) {
	base := context.Background()
	if ctx != nil {
		base = context.WithoutCancel(ctx)
placeholder
	return context.WithTimeout(base, systemUpdateTimeout)
placeholder

type systemUpdateService interface {
	CheckUpdate(ctx context.Context, force bool) (*service.UpdateInfo, error)
	PerformUpdate(ctx context.Context) error
	Rollback() error
	ListRollbackVersions(ctx context.Context) ([]service.RollbackVersion, error)
	RollbackToVersion(ctx context.Context, version string) error
placeholder

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(updateSvc systemUpdateService, lockSvc *service.SystemOperationLockService) *SystemHandler {
	return &SystemHandler{
		updateSvc: updateSvc,
		lockSvc:   lockSvc,
placeholder
placeholder

// GetVersion returns the current version
// GET /api/v1/admin/system/version
func (h *SystemHandler) GetVersion(c *gin.Context) {
	info, _ := h.updateSvc.CheckUpdate(c.Request.Context(), false)
	response.Success(c, gin.H{
		"version": info.CurrentVersion,
placeholder)
placeholder

// CheckUpdates checks for available updates
// GET /api/v1/admin/system/check-updates
func (h *SystemHandler) CheckUpdates(c *gin.Context) {
	force := c.Query("force") == "true"
	info, err := h.updateSvc.CheckUpdate(c.Request.Context(), force)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
placeholder
	response.Success(c, info)
placeholder

// PerformUpdate downloads and applies the update
// POST /api/v1/admin/system/update
func (h *SystemHandler) PerformUpdate(c *gin.Context) {
	operationID := buildSystemOperationID(c, "update")
	payload := gin.H{"operation_id": operationIDplaceholder
	executeAdminIdempotentJSON(c, "admin.system.update", payload, service.DefaultSystemOperationIdempotencyTTL(), func(ctx context.Context) (any, error) {
		lock, release, err := h.acquireSystemLock(ctx, operationID)
		if err != nil {
			return nil, err
	placeholder
		var releaseReason string
		succeeded := false
		defer func() {
			release(releaseReason, succeeded)
	placeholder()

		updateCtx, cancel := systemUpdateContext(ctx)
		defer cancel()

		if err := h.updateSvc.PerformUpdate(updateCtx); err != nil {
			if errors.Is(err, service.ErrNoUpdateAvailable) {
				info, checkErr := h.updateSvc.CheckUpdate(updateCtx, false)
				if checkErr != nil {
					releaseReason = "SYSTEM_UPDATE_FAILED"
					return nil, checkErr
			placeholder
				succeeded = true
				return gin.H{
					"message":            "Already up to date",
					"already_up_to_date": true,
					"current_version":    info.CurrentVersion,
					"latest_version":     info.LatestVersion,
					"operation_id":       lock.OperationID(),
			placeholder, nil
		placeholder
			releaseReason = "SYSTEM_UPDATE_FAILED"
			return nil, err
	placeholder
		succeeded = true

		return gin.H{
			"message":      "Update completed. Please restart the service.",
			"need_restart": true,
			"operation_id": lock.OperationID(),
	placeholder, nil
placeholder)
placeholder

// GetRollbackVersions lists versions available for rollback
// GET /api/v1/admin/system/rollback-versions
func (h *SystemHandler) GetRollbackVersions(c *gin.Context) {
	versions, err := h.updateSvc.ListRollbackVersions(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
placeholder
	response.Success(c, gin.H{
		"versions": versions,
placeholder)
placeholder

// Rollback restores a previous version.
// Without a body (or with an empty version) it restores the local .backup binary
// left by the last in-place update. With {"version": "x.y.z"placeholder it downloads and
// installs that specific release (must be one of the recent rollback versions).
// POST /api/v1/admin/system/rollback
func (h *SystemHandler) Rollback(c *gin.Context) {
	var req struct {
		Version string `json:"version"`
placeholder
	if c.Request.Body != nil && c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "invalid request body")
			return
	placeholder
placeholder
	targetVersion := strings.TrimSpace(req.Version)

	operation := "rollback"
	if targetVersion != "" {
		operation = "rollback:" + targetVersion
placeholder
	operationID := buildSystemOperationID(c, operation)
	payload := gin.H{"operation_id": operationID, "version": targetVersionplaceholder
	executeAdminIdempotentJSON(c, "admin.system.rollback", payload, service.DefaultSystemOperationIdempotencyTTL(), func(ctx context.Context) (any, error) {
		lock, release, err := h.acquireSystemLock(ctx, operationID)
		if err != nil {
			return nil, err
	placeholder
		var releaseReason string
		succeeded := false
		defer func() {
			release(releaseReason, succeeded)
	placeholder()

		if targetVersion != "" {
			// 指定版本回退同样要下载完整二进制，与更新一样和请求生命周期解耦。
			rollbackCtx, cancel := systemUpdateContext(ctx)
			defer cancel()
			err = h.updateSvc.RollbackToVersion(rollbackCtx, targetVersion)
	placeholder else {
			err = h.updateSvc.Rollback()
	placeholder
		if err != nil {
			releaseReason = "SYSTEM_ROLLBACK_FAILED"
			return nil, err
	placeholder
		succeeded = true

		return gin.H{
			"message":      "Rollback completed. Please restart the service.",
			"need_restart": true,
			"version":      targetVersion,
			"operation_id": lock.OperationID(),
	placeholder, nil
placeholder)
placeholder

// RestartService restarts the systemd service
// POST /api/v1/admin/system/restart
func (h *SystemHandler) RestartService(c *gin.Context) {
	operationID := buildSystemOperationID(c, "restart")
	payload := gin.H{"operation_id": operationIDplaceholder
	executeAdminIdempotentJSON(c, "admin.system.restart", payload, service.DefaultSystemOperationIdempotencyTTL(), func(ctx context.Context) (any, error) {
		lock, release, err := h.acquireSystemLock(ctx, operationID)
		if err != nil {
			return nil, err
	placeholder
		succeeded := false
		defer func() {
			release("", succeeded)
	placeholder()

		// Schedule service restart in background after sending response
		// This ensures the client receives the success response before the service restarts
		go func() {
			// Wait a moment to ensure the response is sent
			time.Sleep(500 * time.Millisecond)
			sysutil.RestartServiceAsync()
	placeholder()
		succeeded = true
		return gin.H{
			"message":      "Service restart initiated",
			"operation_id": lock.OperationID(),
	placeholder, nil
placeholder)
placeholder

func (h *SystemHandler) acquireSystemLock(
	ctx context.Context,
	operationID string,
) (*service.SystemOperationLock, func(string, bool), error) {
	if h.lockSvc == nil {
		return nil, nil, service.ErrIdempotencyStoreUnavail
placeholder
	lock, err := h.lockSvc.Acquire(ctx, operationID)
	if err != nil {
		return nil, nil, err
placeholder
	release := func(reason string, succeeded bool) {
		releaseCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = h.lockSvc.Release(releaseCtx, lock, succeeded, reason)
placeholder
	return lock, release, nil
placeholder

func buildSystemOperationID(c *gin.Context, operation string) string {
	key := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if key == "" {
		return "sysop-" + operation + "-" + strconv.FormatInt(time.Now().UnixNano(), 36)
placeholder
	actorScope := "admin:0"
	if subject, ok := middleware2.GetAuthSubjectFromContext(c); ok {
		actorScope = "admin:" + strconv.FormatInt(subject.UserID, 10)
placeholder
	seed := operation + "|" + actorScope + "|" + c.FullPath() + "|" + key
	hash := service.HashIdempotencyKey(seed)
	if len(hash) > 24 {
		hash = hash[:24]
placeholder
	return "sysop-" + hash
placeholder
