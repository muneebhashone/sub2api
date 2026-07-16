package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type batchLimitsAdminServiceStub struct {
	*stubAdminService
	calls []batchLimitsAdminServiceCall
placeholder

type batchLimitsAdminServiceCall struct {
	userIDs     []int64
	concurrency *int
	rpmLimit    *int
placeholder

func cloneIntPointer(value *int) *int {
	if value == nil {
		return nil
placeholder
	cloned := *value
	return &cloned
placeholder

func (s *batchLimitsAdminServiceStub) BatchUpdateLimits(_ context.Context, userIDs []int64, concurrency, rpmLimit *int) (int, error) {
	s.calls = append(s.calls, batchLimitsAdminServiceCall{
		userIDs:     append([]int64(nil), userIDs...),
		concurrency: cloneIntPointer(concurrency),
		rpmLimit:    cloneIntPointer(rpmLimit),
placeholder)
	return len(userIDs), nil
placeholder

func setupBatchLimitsRouter(serviceStub service.AdminService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewUserHandler(serviceStub, nil, nil, nil)
	router.POST("/api/v1/admin/users/batch-limits", handler.BatchUpdateLimits)
	return router
placeholder

func postBatchLimits(t *testing.T, router *gin.Engine, body []byte) *httptest.ResponseRecorder {
placeholder
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/admin/users/batch-limits",
		bytes.NewReader(body),
	)
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	return recorder
placeholder

func TestUserHandlerBatchUpdateLimitsAcceptsPartialAndZeroValues(t *testing.T) {
	tests := []struct {
		name                string
		body                string
		expectedConcurrency *int
		expectedRPMLimit    *int
placeholder{
		{name: "concurrency only", body: `{"user_ids":[1,2],"concurrency":10placeholder`, expectedConcurrency: pointerTo(10)placeholder,
		{name: "both limits", body: `{"user_ids":[1,2],"concurrency":8,"rpm_limit":60placeholder`, expectedConcurrency: pointerTo(8), expectedRPMLimit: pointerTo(60)placeholder,
		{name: "explicit zero", body: `{"user_ids":[1,2],"concurrency":0,"rpm_limit":0placeholder`, expectedConcurrency: pointerTo(0), expectedRPMLimit: pointerTo(0)placeholder,
placeholder

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceStub := &batchLimitsAdminServiceStub{stubAdminService: newStubAdminService()placeholder
			recorder := postBatchLimits(t, setupBatchLimitsRouter(serviceStub), []byte(test.body))

			require.Equal(t, http.StatusOK, recorder.Code)
			require.Len(t, serviceStub.calls, 1)
			require.Equal(t, []int64{1, 2placeholder, serviceStub.calls[0].userIDs)
			require.Equal(t, test.expectedConcurrency, serviceStub.calls[0].concurrency)
			require.Equal(t, test.expectedRPMLimit, serviceStub.calls[0].rpmLimit)

			var response struct {
				Data struct {
					Affected int `json:"affected"`
			placeholder `json:"data"`
		placeholder
			require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			require.Equal(t, 2, response.Data.Affected)
	placeholder)
placeholder
placeholder

func TestUserHandlerBatchUpdateLimitsRejectsInvalidRequests(t *testing.T) {
	tooManyIDs := make([]int64, 501)
	for index := range tooManyIDs {
		tooManyIDs[index] = int64(index + 1)
placeholder
	tooManyBody, err := json.Marshal(map[string]any{"user_ids": tooManyIDs, "rpm_limit": 10placeholder)
placeholder

	tests := []struct {
		name string
		body []byte
placeholder{
		{name: "no limits", body: []byte(`{"user_ids":[1]placeholder`)placeholder,
		{name: "invalid json", body: []byte(`{"user_ids":`)placeholder,
		{name: "missing user ids", body: []byte(`{"rpm_limit":10placeholder`)placeholder,
		{name: "more than 500 ids", body: tooManyBodyplaceholder,
placeholder

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			serviceStub := &batchLimitsAdminServiceStub{stubAdminService: newStubAdminService()placeholder
			recorder := postBatchLimits(t, setupBatchLimitsRouter(serviceStub), test.body)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			require.Empty(t, serviceStub.calls)
	placeholder)
placeholder
placeholder

func TestUserHandlerBatchUpdateLimitsAllUsesEveryListedUser(t *testing.T) {
	base := newStubAdminService()
	base.users = []service.User{{ID: 11placeholder, {ID: 12placeholder, {ID: 13placeholderplaceholder
	serviceStub := &batchLimitsAdminServiceStub{stubAdminService: baseplaceholder
	recorder := postBatchLimits(
		t,
		setupBatchLimitsRouter(serviceStub),
		[]byte(`{"all":true,"user_ids":[999],"rpm_limit":0placeholder`),
	)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Len(t, serviceStub.calls, 1)
	require.Equal(t, []int64{11, 12, 13placeholder, serviceStub.calls[0].userIDs)
	require.Equal(t, 1, base.lastListUsers.calls)
placeholder

func pointerTo(value int) *int {
	return &value
placeholder
