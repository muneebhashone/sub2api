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

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{TotpEnabled: trueplaceholderplaceholder)

	require.False(t, ok)
	require.True(t, c.IsAborted())
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_ADMIN_API_KEY_FORBIDDEN")
placeholder

func TestEnforceStepUpRequiresAuthSubject(t *testing.T) {
	c, rec := newStepUpTestContext(t)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{TotpEnabled: trueplaceholderplaceholder)

	require.False(t, ok)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
placeholder

func TestEnforceStepUpRequiresTotpEnabled(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: falseplaceholderplaceholder)

	require.False(t, ok)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_TOTP_NOT_ENABLED")
placeholder

func TestEnforceStepUpFailsClosedOnGrantError(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{err: errors.New("redis down")placeholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder)

	require.False(t, ok)
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_UNAVAILABLE")
placeholder

func TestEnforceStepUpRequiresGrant(t *testing.T) {
	c, rec := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: falseplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder)

	require.False(t, ok)
	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "STEP_UP_REQUIRED")
placeholder

func TestEnforceStepUpPassesWithGrant(t *testing.T) {
	c, _ := newStepUpTestContext(t)
	c.Set(string(ContextKeyUser), AuthSubject{UserID: 1placeholder)

	ok := enforceStepUp(c, stubStepUpGrantChecker{granted: trueplaceholder, stubStepUpUserReader{user: &service.User{ID: 1, TotpEnabled: trueplaceholderplaceholder)

	require.True(t, ok)
	require.False(t, c.IsAborted())
placeholder
