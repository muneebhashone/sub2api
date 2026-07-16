//go:build unit

package admin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type duplicateGroupAdminServiceStub struct {
	service.AdminService
	group        *service.Group
	calls        int
	recoverCalls int
	groupID      int64
	actorScope   string
	operationKey string
	recoverScope string
	recoverKey   string
	created      bool
placeholder

func (s *duplicateGroupAdminServiceStub) DuplicateGroup(_ context.Context, groupID int64, actorScope, operationKey string) (*service.Group, error) {
	s.calls++
	s.groupID = groupID
	s.actorScope = actorScope
	s.operationKey = operationKey
	s.created = true
	return s.group, nil
placeholder

func (s *duplicateGroupAdminServiceStub) RecoverDuplicateGroup(_ context.Context, _ int64, actorScope, operationKey string) (*service.Group, error) {
	s.recoverCalls++
	s.recoverScope = actorScope
	s.recoverKey = operationKey
	if !s.created {
		return nil, nil
placeholder
	return s.group, nil
placeholder

func setupDuplicateGroupRouter(t *testing.T, svc service.AdminService) *gin.Engine {
placeholder
	previousCoordinator := service.DefaultIdempotencyCoordinator()
	service.SetDefaultIdempotencyCoordinator(nil)
	t.Cleanup(func() { service.SetDefaultIdempotencyCoordinator(previousCoordinator) placeholder)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 77placeholder)
		c.Next()
placeholder)
	handler := NewGroupHandler(svc, nil, nil)
	router.POST("/api/v1/admin/groups/:id/duplicate", handler.Duplicate)
	return router
placeholder

func duplicateGroupHandlerFixture() *service.Group {
	return &service.Group{
		ID:                   43,
		Name:                 "primary (Copy)",
		Platform:             service.PlatformAnthropic,
		Status:               "inactive",
		RateMultiplier:       1,
		AccountCount:         3,
		ActiveAccountCount:   2,
		DuplicateOperationID: "internal-operation-must-not-leak",
		ModelRouting:         map[string][]int64{"claude-*": {7placeholderplaceholder,
placeholder
placeholder

func TestDuplicateGroupHandlerReturnsAdminDTOWithoutOperationMetadata(t *testing.T) {
	svc := &duplicateGroupAdminServiceStub{group: duplicateGroupHandlerFixture()placeholder
	router := setupDuplicateGroupRouter(t, svc)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/groups/42/duplicate", nil)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, svc.calls)
	require.Contains(t, recorder.Body.String(), `"name":"primary (Copy)"`)
	require.Contains(t, recorder.Body.String(), `"status":"inactive"`)
	require.Contains(t, recorder.Body.String(), `"account_count":3`)
	require.NotContains(t, recorder.Body.String(), "duplicate_operation_id")
	require.NotContains(t, recorder.Body.String(), "internal-operation-must-not-leak")
placeholder

func TestDuplicateGroupHandlerRejectsInvalidID(t *testing.T) {
	for _, id := range []string{"not-a-number", "0", "-1"placeholder {
		t.Run(id, func(t *testing.T) {
			svc := &duplicateGroupAdminServiceStub{placeholder
			router := setupDuplicateGroupRouter(t, svc)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/groups/"+id+"/duplicate", nil)

			router.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Zero(t, svc.calls)
	placeholder)
placeholder
placeholder

func TestDuplicateGroupHandlerReplaysSameIdempotencyKey(t *testing.T) {
	svc := &duplicateGroupAdminServiceStub{group: duplicateGroupHandlerFixture()placeholder
	router := setupDuplicateGroupRouter(t, svc)
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(newMemoryIdempotencyRepoStub(), service.DefaultIdempotencyConfig()))

	call := func() *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/groups/42/duplicate", nil)
		request.Header.Set("Idempotency-Key", "duplicate-group-42")
		router.ServeHTTP(recorder, request)
		return recorder
placeholder

	first := call()
	second := call()

	require.Equal(t, http.StatusOK, first.Code)
	require.Equal(t, http.StatusOK, second.Code)
	require.Equal(t, 1, svc.calls)
	require.Equal(t, int64(42), svc.groupID)
	require.Equal(t, "admin:77", svc.actorScope)
	require.Equal(t, "duplicate-group-42", svc.operationKey)
	require.Equal(t, "true", second.Header().Get("X-Idempotency-Replayed"))
placeholder

func TestDuplicateGroupHandlerRecoversAfterMarkSucceededFailure(t *testing.T) {
	svc := &duplicateGroupAdminServiceStub{group: duplicateGroupHandlerFixture()placeholder
	router := setupDuplicateGroupRouter(t, svc)
	repo := &failOnceMarkSucceededRepo{memoryIdempotencyRepoStub: newMemoryIdempotencyRepoStub(), failNext: trueplaceholder
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(repo, service.DefaultIdempotencyConfig()))

	call := func() *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/groups/42/duplicate", nil)
		request.Header.Set("Idempotency-Key", "duplicate-group-42-recovery")
		router.ServeHTTP(recorder, request)
		return recorder
placeholder

	first := call()
	second := call()

	require.Equal(t, http.StatusOK, first.Code)
	require.Equal(t, http.StatusOK, second.Code)
	require.Equal(t, "true", first.Header().Get("X-Idempotency-Recovered"))
	require.Equal(t, "true", second.Header().Get("X-Idempotency-Recovered"))
	require.Equal(t, 1, svc.calls, "ambiguous retries must not repeat the create side effect")
	require.Equal(t, 2, svc.recoverCalls)
	require.Equal(t, "admin:77", svc.recoverScope)
	require.Equal(t, "duplicate-group-42-recovery", svc.recoverKey)
placeholder
