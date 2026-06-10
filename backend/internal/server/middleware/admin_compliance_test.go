package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type complianceGuardRepoStub struct {
	values map[string]string
placeholder

func (r *complianceGuardRepoStub) Get(ctx context.Context, key string) (*service.Setting, error) {
	if value, ok := r.values[key]; ok {
		return &service.Setting{Key: key, Value: valueplaceholder, nil
placeholder
	return nil, service.ErrSettingNotFound
placeholder

func (r *complianceGuardRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	setting, err := r.Get(ctx, key)
	if err != nil {
		return "", err
placeholder
	return setting.Value, nil
placeholder

func (r *complianceGuardRepoStub) Set(ctx context.Context, key, value string) error { return nil placeholder
func (r *complianceGuardRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder
func (r *complianceGuardRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	return nil
placeholder
func (r *complianceGuardRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder
func (r *complianceGuardRepoStub) Delete(ctx context.Context, key string) error { return nil placeholder

func TestAdminComplianceGuardBlocksAdminRouteWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := service.NewSettingService(&complianceGuardRepoStub{placeholder, &config.Config{placeholder)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)
		c.Next()
placeholder)
	router.Use(AdminComplianceGuard(svc))
	router.GET("/api/v1/admin/users", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
placeholder)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusLocked, w.Code)
	require.Contains(t, w.Body.String(), "ADMIN_COMPLIANCE_ACK_REQUIRED")
placeholder

func TestAdminComplianceGuardBypassesComplianceEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := service.NewSettingService(&complianceGuardRepoStub{placeholder, &config.Config{placeholder)
	router := gin.New()
	router.Use(AdminComplianceGuard(svc))
	router.GET("/api/v1/admin/compliance", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
placeholder)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/compliance", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "ok", w.Body.String())
placeholder
