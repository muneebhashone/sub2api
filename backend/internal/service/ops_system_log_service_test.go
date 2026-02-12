package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

func TestOpsServiceListSystemLogs_DefaultClampAndSuccess(t *testing.T) {
	var gotFilter *OpsSystemLogFilter
	repo := &opsRepoMock{
		ListSystemLogsFn: func(ctx context.Context, filter *OpsSystemLogFilter) (*OpsSystemLogList, error) {
			gotFilter = filter
			return &OpsSystemLogList{
				Logs:     []*OpsSystemLog{{ID: 1, Level: "warn", Message: "x"placeholderplaceholder,
				Total:    1,
				Page:     filter.Page,
				PageSize: filter.PageSize,
		placeholder, nil
	placeholder,
placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

	out, err := svc.ListSystemLogs(context.Background(), &OpsSystemLogFilter{
		Page:     0,
		PageSize: 999,
placeholder)
	if err != nil {
		t.Fatalf("ListSystemLogs() error: %v", err)
placeholder
	if gotFilter == nil {
		t.Fatalf("expected repository to receive filter")
placeholder
	if gotFilter.Page != 1 || gotFilter.PageSize != 200 {
		t.Fatalf("filter normalized unexpectedly: page=%d pageSize=%d", gotFilter.Page, gotFilter.PageSize)
placeholder
	if out.Total != 1 || len(out.Logs) != 1 {
		t.Fatalf("unexpected result: %+v", out)
placeholder
placeholder

func TestOpsServiceListSystemLogs_MonitoringDisabled(t *testing.T) {
	svc := NewOpsService(
		&opsRepoMock{placeholder,
		nil,
		&config.Config{Ops: config.OpsConfig{Enabled: falseplaceholderplaceholder,
		nil, nil, nil, nil, nil, nil, nil, nil,
	)
	_, err := svc.ListSystemLogs(context.Background(), &OpsSystemLogFilter{placeholder)
	if err == nil {
		t.Fatalf("expected disabled error")
placeholder
placeholder

func TestOpsServiceListSystemLogs_NilRepoReturnsEmpty(t *testing.T) {
	svc := NewOpsService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	out, err := svc.ListSystemLogs(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListSystemLogs() error: %v", err)
placeholder
	if out == nil || out.Page != 1 || out.PageSize != 50 || out.Total != 0 || len(out.Logs) != 0 {
		t.Fatalf("unexpected nil-repo result: %+v", out)
placeholder
placeholder

func TestOpsServiceListSystemLogs_RepoErrorMapped(t *testing.T) {
	repo := &opsRepoMock{
		ListSystemLogsFn: func(ctx context.Context, filter *OpsSystemLogFilter) (*OpsSystemLogList, error) {
			return nil, errors.New("db down")
	placeholder,
placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	_, err := svc.ListSystemLogs(context.Background(), &OpsSystemLogFilter{placeholder)
	if err == nil {
		t.Fatalf("expected mapped internal error")
placeholder
	if !strings.Contains(err.Error(), "OPS_SYSTEM_LOG_LIST_FAILED") {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

func TestOpsServiceCleanupSystemLogs_SuccessAndAudit(t *testing.T) {
	var audit *OpsSystemLogCleanupAudit
	repo := &opsRepoMock{
		DeleteSystemLogsFn: func(ctx context.Context, filter *OpsSystemLogCleanupFilter) (int64, error) {
			return 3, nil
	placeholder,
		InsertSystemLogCleanupAuditFn: func(ctx context.Context, input *OpsSystemLogCleanupAudit) error {
			audit = input
			return nil
	placeholder,
placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	userID := int64(7)
	now := time.Now().UTC()
	filter := &OpsSystemLogCleanupFilter{
		StartTime:       &now,
		Level:           "warn",
		RequestID:       "req-1",
		ClientRequestID: "creq-1",
		UserID:          &userID,
		Query:           "timeout",
placeholder

	deleted, err := svc.CleanupSystemLogs(context.Background(), filter, 99)
	if err != nil {
		t.Fatalf("CleanupSystemLogs() error: %v", err)
placeholder
	if deleted != 3 {
		t.Fatalf("deleted=%d, want 3", deleted)
placeholder
	if audit == nil {
		t.Fatalf("expected cleanup audit")
placeholder
	if !strings.Contains(audit.Conditions, `"client_request_id":"creq-1"`) {
		t.Fatalf("audit conditions should include client_request_id: %s", audit.Conditions)
placeholder
	if !strings.Contains(audit.Conditions, `"user_id":7`) {
		t.Fatalf("audit conditions should include user_id: %s", audit.Conditions)
placeholder
placeholder

func TestOpsServiceCleanupSystemLogs_RepoUnavailableAndInvalidOperator(t *testing.T) {
	svc := NewOpsService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if _, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{RequestID: "r"placeholder, 1); err == nil {
		t.Fatalf("expected repo unavailable error")
placeholder

	svc = NewOpsService(&opsRepoMock{placeholder, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	if _, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{RequestID: "r"placeholder, 0); err == nil {
		t.Fatalf("expected invalid operator error")
placeholder
placeholder

func TestOpsServiceCleanupSystemLogs_FilterRequired(t *testing.T) {
	repo := &opsRepoMock{
		DeleteSystemLogsFn: func(ctx context.Context, filter *OpsSystemLogCleanupFilter) (int64, error) {
			return 0, errors.New("cleanup requires at least one filter condition")
	placeholder,
placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	_, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{placeholder, 1)
	if err == nil {
		t.Fatalf("expected filter required error")
placeholder
	if !strings.Contains(strings.ToLower(err.Error()), "filter") {
		t.Fatalf("unexpected error: %v", err)
placeholder
placeholder

func TestOpsServiceCleanupSystemLogs_InvalidRange(t *testing.T) {
	repo := &opsRepoMock{placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	start := time.Now().UTC()
	end := start.Add(-time.Hour)
	_, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{
		StartTime: &start,
		EndTime:   &end,
placeholder, 1)
	if err == nil {
		t.Fatalf("expected invalid range error")
placeholder
placeholder

func TestOpsServiceCleanupSystemLogs_NoRowsAndInternalError(t *testing.T) {
	repo := &opsRepoMock{
		DeleteSystemLogsFn: func(ctx context.Context, filter *OpsSystemLogCleanupFilter) (int64, error) {
			return 0, sql.ErrNoRows
	placeholder,
placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	deleted, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{
		RequestID: "req-1",
placeholder, 1)
	if err != nil || deleted != 0 {
		t.Fatalf("expected no rows shortcut, deleted=%d err=%v", deleted, err)
placeholder

	repo.DeleteSystemLogsFn = func(ctx context.Context, filter *OpsSystemLogCleanupFilter) (int64, error) {
		return 0, errors.New("boom")
placeholder
	if _, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{
		RequestID: "req-1",
placeholder, 1); err == nil {
		t.Fatalf("expected internal cleanup error")
placeholder
placeholder

func TestOpsServiceCleanupSystemLogs_AuditFailureIgnored(t *testing.T) {
	repo := &opsRepoMock{
		DeleteSystemLogsFn: func(ctx context.Context, filter *OpsSystemLogCleanupFilter) (int64, error) {
			return 5, nil
	placeholder,
		InsertSystemLogCleanupAuditFn: func(ctx context.Context, input *OpsSystemLogCleanupAudit) error {
			return errors.New("audit down")
	placeholder,
placeholder
	svc := NewOpsService(repo, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	deleted, err := svc.CleanupSystemLogs(context.Background(), &OpsSystemLogCleanupFilter{
		RequestID: "r1",
placeholder, 1)
	if err != nil || deleted != 5 {
		t.Fatalf("audit failure should not break cleanup, deleted=%d err=%v", deleted, err)
placeholder
placeholder

func TestMarshalSystemLogCleanupConditions_NilAndMarshalError(t *testing.T) {
	if got := marshalSystemLogCleanupConditions(nil); got != "{placeholder" {
		t.Fatalf("nil filter should return {placeholder, got %s", got)
placeholder

	now := time.Now().UTC()
	userID := int64(1)
	filter := &OpsSystemLogCleanupFilter{
		StartTime: &now,
		EndTime:   &now,
		UserID:    &userID,
placeholder
	got := marshalSystemLogCleanupConditions(filter)
	if !strings.Contains(got, `"start_time"`) || !strings.Contains(got, `"user_id":1`) {
		t.Fatalf("unexpected marshal payload: %s", got)
placeholder
placeholder

func TestOpsServiceGetSystemLogSinkHealth(t *testing.T) {
	svc := NewOpsService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	health := svc.GetSystemLogSinkHealth()
	if health.QueueCapacity != 0 || health.QueueDepth != 0 {
		t.Fatalf("unexpected health for nil sink: %+v", health)
placeholder

	sink := NewOpsSystemLogSink(&opsRepoMock{placeholder)
	svc = NewOpsService(&opsRepoMock{placeholder, nil, nil, nil, nil, nil, nil, nil, nil, nil, sink)
	health = svc.GetSystemLogSinkHealth()
	if health.QueueCapacity <= 0 {
		t.Fatalf("expected non-zero queue capacity: %+v", health)
placeholder
placeholder
