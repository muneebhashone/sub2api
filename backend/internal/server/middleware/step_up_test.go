package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type stubStepUpGrantChecker struct {
	granted bool
	err     error
placeholder

func (s stubStepUpGrantChecker) HasStepUpGrant(ctx context.Context, userID int64, sessionKey string) (bool, error) {
	return s.granted, s.err
placeholder

type stubStepUpUserReader struct {
	user *service.User
	err  error
placeholder

func (s stubStepUpUserReader) GetByID(ctx context.Context, id int64) (*service.User, error) {
	return s.user, s.err
placeholder

type stubStepUpSettingReader struct {
	enabled bool
placeholder

func (s stubStepUpSettingReader) IsStepUpEnabled(ctx context.Context) bool {
	return s.enabled
placeholder

// stepUpEnabled 功能开关开启的设置桩，供既有门控分支测试使用。
var stepUpEnabled = stubStepUpSettingReader{enabled: trueplaceholder

func newStepUpTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
placeholder
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/sensitive", nil)
	return c, rec
placeholder

func TestEnforceStepUpRejectsAdminAPIKey(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set("auth_method", service.AuditAuthMethodAdminAPIKey)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{TotpEnabled: trueplaceholderplaceholder, stepUpEnabled)

	require.False(t, ok)
	require.True(t, c.IsAborted())
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_ADMIN_API_KEY_FORBIDDEN")
placeholder

func TestEnforceStepUpRequiresAuthSubject(t *testing.T) {
	c, rec := newStepUpTestContext(t)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{TotpEnabled: trueplaceholderplaceholder, stepUpEnabled)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
placeholder

func TestEnforceStepUpRequiresTotpEnabled(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: falseplaceholderplaceholder, stepUpEnabled)

	require.False(t, ok)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_TOTP_NOT_ENABLED")
placeholder

func TestEnforceStepUpFailsClosedOnGrantError(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{err: errors.New("redis down")placeholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder, stepUpEnabled)

	require.False(t, ok)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_UNAVAILABLE")
placeholder

func TestEnforceStepUpRequiresGrant(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: falseplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder, stepUpEnabled)

	require.False(t, ok)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_REQUIRED")
placeholder

func TestEnforceStepUpPassesWithGrant(t *testing.T) {
	c, _ := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder, stepUpEnabled)

	require.True(t, ok)
	require.False(t, c.IsAborted())
placeholder

// 功能开关关闭时：不论 TOTP/grant/凭证类型，一律放行（恢复门控引入前行为）。
func TestEnforceStepUpDisabledSkipsAllChecks(t *testing.T) {
	disabled := stubStepUpSettingReader{enabled: falseplaceholder

	t.Run("no totp, no grant", func(t *testing.T) {
		c, _ := newStepUpTestContext(t)
		c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

		ok := enforceStepUp(c, stubStepUpGrantChecker{granted: falseplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: falseplaceholderplaceholder, disabled)

		require.True(t, ok)
		require.False(t, c.IsAborted())
placeholder)

	t.Run("admin api key", func(t *testing.T) {
		c, _ := newStepUpTestContext(t)
		c.Set("auth_method", service.AuditAuthMethodAdminAPIKey)

		ok := enforceStepUp(c, stubStepUpGrantChecker{granted: falseplaceholder, stubStepUpUserReader{user: nil, err: errors.New("should not be called")placeholder, disabled)

		require.True(t, ok)
		require.False(t, c.IsAborted())
placeholder)
placeholder

// settings 为 nil 时保持门控（fail-closed），避免装配缺陷静默关闭安全控制。
func TestEnforceStepUpNilSettingsFailsClosed(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: falseplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder, nil)

	require.False(t, ok)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_REQUIRED")
placeholder

// EnforceStepUp 收到 nil *service.SettingService 时不得因 typed-nil 装箱绕过门控：
// 未认证请求仍应被拦截（401），而不是当作"开关关闭"放行。
func TestEnforceStepUpTypedNilSettingServiceFailsClosed(t *testing.T) {
	require.Nil(t, stepUpSettingsOrNil(nil))

	c, rec := newStepUpTestContext(t)

	ok := EnforceStepUp(c, nil, nil, nil)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
placeholder
