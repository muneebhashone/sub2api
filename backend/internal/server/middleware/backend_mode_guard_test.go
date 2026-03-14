//go:build unit

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

type bmSettingRepo struct {
	values map[string]string
placeholder

func (r *bmSettingRepo) Get(_ context.Context, _ string) (*service.Setting, error) {
	panic("unexpected Get call")
placeholder

func (r *bmSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	v, ok := r.values[key]
	if !ok {
		return "", service.ErrSettingNotFound
placeholder
	return v, nil
placeholder

func (r *bmSettingRepo) Set(_ context.Context, _, _ string) error {
	panic("unexpected Set call")
placeholder

func (r *bmSettingRepo) GetMultiple(_ context.Context, _ []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (r *bmSettingRepo) SetMultiple(_ context.Context, settings map[string]string) error {
	if r.values == nil {
		r.values = make(map[string]string, len(settings))
placeholder
	for key, value := range settings {
		r.values[key] = value
placeholder
	return nil
placeholder

func (r *bmSettingRepo) GetAll(_ context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (r *bmSettingRepo) Delete(_ context.Context, _ string) error {
	panic("unexpected Delete call")
placeholder

func newBackendModeSettingService(t *testing.T, enabled string) *service.SettingService {
placeholder

	repo := &bmSettingRepo{
		values: map[string]string{
			service.SettingKeyBackendModeEnabled: enabled,
	placeholder,
placeholder
	svc := service.NewSettingService(repo, &config.Config{placeholder)
	require.NoError(t, svc.UpdateSettings(context.Background(), &service.SystemSettings{
		BackendModeEnabled: enabled == "true",
placeholder))

	return svc
placeholder

func stringPtr(v string) *string {
	return &v
placeholder

func TestBackendModeUserGuard(t *testing.T) {
	tests := []struct {
		name       string
		nilService bool
		enabled    string
		role       *string
		wantStatus int
placeholder{
		{
			name:       "disabled_allows_all",
			enabled:    "false",
			role:       stringPtr("user"),
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "nil_service_allows_all",
			nilService: true,
			role:       stringPtr("user"),
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_admin_allowed",
			enabled:    "true",
			role:       stringPtr("admin"),
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_user_blocked",
			enabled:    "true",
			role:       stringPtr("user"),
			wantStatus: http.StatusForbidden,
	placeholder,
		{
			name:       "enabled_no_role_blocked",
			enabled:    "true",
			wantStatus: http.StatusForbidden,
	placeholder,
		{
			name:       "enabled_empty_role_blocked",
			enabled:    "true",
			role:       stringPtr(""),
			wantStatus: http.StatusForbidden,
	placeholder,
placeholder

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			r := gin.New()
			if tc.role != nil {
				role := *tc.role
				r.Use(func(c *gin.Context) {
					c.Set(string(ContextKeyUserRole), role)
					c.Next()
			placeholder)
		placeholder

			var svc *service.SettingService
			if !tc.nilService {
				svc = newBackendModeSettingService(t, tc.enabled)
		placeholder

			r.Use(BackendModeUserGuard(svc))
			r.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": trueplaceholder)
		placeholder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
	placeholder)
placeholder
placeholder

func TestBackendModeAuthGuard(t *testing.T) {
	tests := []struct {
		name       string
		nilService bool
		enabled    string
		path       string
		wantStatus int
placeholder{
		{
			name:       "disabled_allows_all",
			enabled:    "false",
			path:       "/api/v1/auth/register",
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "nil_service_allows_all",
			nilService: true,
			path:       "/api/v1/auth/register",
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_allows_login",
			enabled:    "true",
			path:       "/api/v1/auth/login",
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_allows_login_2fa",
			enabled:    "true",
			path:       "/api/v1/auth/login/2fa",
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_allows_logout",
			enabled:    "true",
			path:       "/api/v1/auth/logout",
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_allows_refresh",
			enabled:    "true",
			path:       "/api/v1/auth/refresh",
			wantStatus: http.StatusOK,
	placeholder,
		{
			name:       "enabled_blocks_register",
			enabled:    "true",
			path:       "/api/v1/auth/register",
			wantStatus: http.StatusForbidden,
	placeholder,
		{
			name:       "enabled_blocks_forgot_password",
			enabled:    "true",
			path:       "/api/v1/auth/forgot-password",
			wantStatus: http.StatusForbidden,
	placeholder,
placeholder

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			r := gin.New()

			var svc *service.SettingService
			if !tc.nilService {
				svc = newBackendModeSettingService(t, tc.enabled)
		placeholder

			r.Use(BackendModeAuthGuard(svc))
			r.Any("/*path", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": trueplaceholder)
		placeholder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
	placeholder)
placeholder
placeholder
