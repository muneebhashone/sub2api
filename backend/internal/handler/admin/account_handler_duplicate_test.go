//go:build unit

package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type duplicateAccountAdminServiceStub struct {
	service.AdminService
	account      *service.Account
	calls        int
	recoverCalls int
	accountID    int64
	actorScope   string
	operationKey string
	recoverScope string
	recoverKey   string
	recoverErr   error
	created      bool
placeholder

type blockingDuplicateAdminServiceStub struct {
	service.AdminService
	account      *service.Account
	started      chan struct{placeholder
	release      chan struct{placeholder
	calls        atomic.Int32
	recoverCalls atomic.Int32
	recoverErr   error
placeholder

type failOnceMarkSucceededRepo struct {
	*memoryIdempotencyRepoStub
	failNext bool
placeholder

func (r *failOnceMarkSucceededRepo) MarkSucceeded(ctx context.Context, id int64, responseStatus int, responseBody string, expiresAt time.Time) error {
	if r.failNext {
		r.failNext = false
		return errors.New("mark succeeded failed")
placeholder
	return r.memoryIdempotencyRepoStub.MarkSucceeded(ctx, id, responseStatus, responseBody, expiresAt)
placeholder

func (s *duplicateAccountAdminServiceStub) DuplicateAccount(_ context.Context, accountID int64, actorScope, operationKey string) (*service.Account, error) {
	s.calls++
	s.accountID = accountID
	s.actorScope = actorScope
	s.operationKey = operationKey
	s.created = true
	return s.account, nil
placeholder

func (s *duplicateAccountAdminServiceStub) RecoverDuplicateAccount(_ context.Context, _ int64, actorScope, operationKey string) (*service.Account, error) {
	s.recoverCalls++
	s.recoverScope = actorScope
	s.recoverKey = operationKey
	if s.recoverErr != nil {
		return nil, s.recoverErr
placeholder
	if !s.created {
		return nil, nil
placeholder
	return s.account, nil
placeholder

func (s *blockingDuplicateAdminServiceStub) DuplicateAccount(_ context.Context, _ int64, _, _ string) (*service.Account, error) {
	s.calls.Add(1)
	close(s.started)
	<-s.release
	return s.account, nil
placeholder

func (s *blockingDuplicateAdminServiceStub) RecoverDuplicateAccount(_ context.Context, _ int64, _, _ string) (*service.Account, error) {
	s.recoverCalls.Add(1)
	return nil, s.recoverErr
placeholder

func setupDuplicateAccountRouter(t *testing.T, svc service.AdminService) *gin.Engine {
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
	handler := NewAccountHandler(svc, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	router.POST("/api/v1/admin/accounts/:id/duplicate", handler.Duplicate)
	return router
placeholder

func TestDuplicateAccountHandlerRedactsCredentials(t *testing.T) {
	svc := &duplicateAccountAdminServiceStub{
		account: &service.Account{
			ID:          43,
			Name:        "primary (Copy)",
			Platform:    service.PlatformAnthropic,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: false,
	placeholder"api_key": "top-secret-key"placeholder,
	placeholder,
placeholder
	router := setupDuplicateAccountRouter(t, svc)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/42/duplicate", nil)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Equal(t, 1, svc.calls)
	require.Contains(t, recorder.Body.String(), `"name":"primary (Copy)"`)
	require.NotContains(t, recorder.Body.String(), "top-secret-key")
	var responseBody struct {
		Data struct {
			Credentials map[string]any `json:"credentials"`
			Schedulable bool           `json:"schedulable"`
	placeholder `json:"data"`
placeholder
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &responseBody))
	require.Empty(t, responseBody.Data.Credentials)
	require.False(t, responseBody.Data.Schedulable)
placeholder

func TestDuplicateAccountHandlerRejectsInvalidID(t *testing.T) {
	svc := &duplicateAccountAdminServiceStub{placeholder
	router := setupDuplicateAccountRouter(t, svc)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/not-a-number/duplicate", nil)

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Zero(t, svc.calls)
placeholder

func TestDuplicateAccountHandlerReplaysSameIdempotencyKey(t *testing.T) {
	svc := &duplicateAccountAdminServiceStub{
		account: &service.Account{
			ID:          43,
			Name:        "primary (Copy)",
			Platform:    service.PlatformAnthropic,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: false,
	placeholder,
placeholder
	router := setupDuplicateAccountRouter(t, svc)
	repo := newMemoryIdempotencyRepoStub()
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(repo, service.DefaultIdempotencyConfig()))

	call := func() *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/42/duplicate", nil)
		request.Header.Set("Idempotency-Key", "duplicate-account-42")
		router.ServeHTTP(recorder, request)
		return recorder
placeholder

	first := call()
	second := call()

	require.Equal(t, http.StatusOK, first.Code)
	require.Equal(t, http.StatusOK, second.Code)
	require.Equal(t, 1, svc.calls)
	require.Equal(t, int64(42), svc.accountID)
	require.Equal(t, "admin:77", svc.actorScope)
	require.Equal(t, "duplicate-account-42", svc.operationKey)
	require.Equal(t, "true", second.Header().Get("X-Idempotency-Replayed"))
placeholder

func TestDuplicateAccountHandlerRecoversAfterMarkSucceededFailure(t *testing.T) {
	svc := &duplicateAccountAdminServiceStub{
		account: &service.Account{
			ID:          43,
			Name:        "primary (Copy)",
			Platform:    service.PlatformAnthropic,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: false,
	placeholder,
placeholder
	router := setupDuplicateAccountRouter(t, svc)
	repo := &failOnceMarkSucceededRepo{memoryIdempotencyRepoStub: newMemoryIdempotencyRepoStub(), failNext: trueplaceholder
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(repo, service.DefaultIdempotencyConfig()))

	call := func() *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/42/duplicate", nil)
		request.Header.Set("Idempotency-Key", "duplicate-account-42-recovery")
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
	require.Equal(t, "duplicate-account-42-recovery", svc.recoverKey)
	require.Contains(t, second.Body.String(), `"id":43`)
placeholder

func TestDuplicateAccountHandlerPreservesIdempotencyErrorWhenRecoveryLookupFails(t *testing.T) {
	svc := &duplicateAccountAdminServiceStub{
		account: &service.Account{
			ID:          43,
			Name:        "primary (Copy)",
			Platform:    service.PlatformAnthropic,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: false,
	placeholder,
		recoverErr: errors.New("recovery database unavailable"),
placeholder
	router := setupDuplicateAccountRouter(t, svc)
	repo := &failOnceMarkSucceededRepo{memoryIdempotencyRepoStub: newMemoryIdempotencyRepoStub(), failNext: trueplaceholder
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(repo, service.DefaultIdempotencyConfig()))
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/42/duplicate", nil)
	request.Header.Set("Idempotency-Key", "duplicate-account-42-recovery-error")

	router.ServeHTTP(recorder, request)

	require.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	require.Contains(t, recorder.Body.String(), "IDEMPOTENCY_STORE_UNAVAILABLE")
	require.Equal(t, 1, svc.calls)
	require.Equal(t, 1, svc.recoverCalls)
	require.Equal(t, "admin:77", svc.recoverScope)
placeholder

func TestDuplicateAccountHandlerDoesNotReexecuteWhileOriginalIsProcessing(t *testing.T) {
	svc := &blockingDuplicateAdminServiceStub{
		account: &service.Account{
			ID:          43,
			Name:        "primary (Copy)",
			Platform:    service.PlatformAnthropic,
			Type:        service.AccountTypeAPIKey,
			Status:      service.StatusActive,
			Schedulable: false,
	placeholder,
		started:    make(chan struct{placeholder),
		release:    make(chan struct{placeholder),
		recoverErr: errors.New("recovery database unavailable"),
placeholder
	router := setupDuplicateAccountRouter(t, svc)
	service.SetDefaultIdempotencyCoordinator(service.NewIdempotencyCoordinator(newMemoryIdempotencyRepoStub(), service.DefaultIdempotencyConfig()))

	call := func() *httptest.ResponseRecorder {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/accounts/42/duplicate", nil)
		request.Header.Set("Idempotency-Key", "duplicate-account-42-active")
		router.ServeHTTP(recorder, request)
		return recorder
placeholder
	firstDone := make(chan *httptest.ResponseRecorder, 1)
	go func() { firstDone <- call() placeholder()
	<-svc.started

	second := call()
	require.Equal(t, http.StatusConflict, second.Code)
	require.Contains(t, second.Body.String(), "IDEMPOTENCY_IN_PROGRESS")
	require.Equal(t, int32(1), svc.calls.Load())
	require.Equal(t, int32(1), svc.recoverCalls.Load())

	close(svc.release)
	select {
	case first := <-firstDone:
		require.Equal(t, http.StatusOK, first.Code)
	case <-time.After(time.Second):
		t.Fatal("original duplicate request did not finish")
placeholder
placeholder
