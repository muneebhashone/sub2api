package service

import (
	"context"
	"database/sql"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

type stubOpsRepoForUserErr struct {
	OpsRepository // 嵌入接口，未实现的方法 panic，仅覆盖 ListErrorLogs
	gotFilter     *OpsErrorLogFilter

	// GetErrorLogByID 控制字段
	detailToReturn    *OpsErrorLogDetail
	detailErrToReturn error
placeholder

func (s *stubOpsRepoForUserErr) ListErrorLogs(ctx context.Context, f *OpsErrorLogFilter) (*OpsErrorLogList, error) {
	s.gotFilter = f
	return &OpsErrorLogList{
		Errors: []*OpsErrorLog{{
			Phase: "request", Type: "rate_limit_error",
			Model: "m", RequestedModel: "rm", StatusCode: 429,
			Message: "secret", UserEmail: "a@b.c",
placeholder
		Total: 1, Page: 1, PageSize: 20,
placeholder, nil
placeholder

func (s *stubOpsRepoForUserErr) GetErrorLogByID(ctx context.Context, id int64) (*OpsErrorLogDetail, error) {
	if s.detailErrToReturn != nil {
		return nil, s.detailErrToReturn
placeholder
	return s.detailToReturn, nil
placeholder

func TestListUserErrorRequests_ForcesScopeAndRedacts(t *testing.T) {
	stub := &stubOpsRepoForUserErr{placeholder
	svc := &OpsService{opsRepo: stubplaceholder
	uid := int64(42)
	kid := int64(7)
	in := &OpsErrorLogFilter{UserID: nil, View: "errors", Phase: "upstream", APIKeyID: &kidplaceholder
	out, err := svc.ListUserErrorRequests(context.Background(), uid, in)
	if err != nil {
		t.Fatal(err)
placeholder
	// 强制按用户
	if stub.gotFilter.UserID == nil || *stub.gotFilter.UserID != uid {
		t.Fatalf("UserID not forced: %+v", stub.gotFilter.UserID)
placeholder
	// 强制 View=all（含业务限流/余额）
	if stub.gotFilter.View != "all" {
		t.Fatalf("View not forced to all: %q", stub.gotFilter.View)
placeholder
	// 强制排除 count_tokens
	if !stub.gotFilter.ExcludeCountTokens {
		t.Fatal("ExcludeCountTokens not forced")
placeholder
	// 强制清空 Phase（防止 "upstream" 绕过 status>=400 子句 + 与 ErrorPhasesAny 双重约束）
	if stub.gotFilter.Phase != "" {
		t.Fatalf("Phase not cleared: %q", stub.gotFilter.Phase)
placeholder
	// APIKeyID 透传保留（用户可按自己 key 过滤；越权由 user_id AND api_key_id 双重防护）
	if stub.gotFilter.APIKeyID == nil || *stub.gotFilter.APIKeyID != kid {
		t.Fatalf("APIKeyID should be preserved, got %v", stub.gotFilter.APIKeyID)
placeholder
	// 调用方传入的 filter 不应被原地篡改（验证 shallow copy 隔离生效）
	if in.View != "errors" || in.UserID != nil || in.Phase != "upstream" {
		t.Fatalf("caller filter was mutated: View=%q UserID=%v Phase=%q", in.View, in.UserID, in.Phase)
placeholder
	// 脱敏：返回条目含 message 字段
	if len(out.Items) != 1 || out.Items[0].Category != "rate_limit" || out.Items[0].Model != "rm" {
		t.Fatalf("bad item: %+v", out.Items)
placeholder
placeholder

func TestGetUserErrorRequestDetail_OwnershipEnforced(t *testing.T) {
	ownerUID := int64(999)
	callerUID := int64(1)
	upstreamStatus := 503

	detail := &OpsErrorLogDetail{
		OpsErrorLog: OpsErrorLog{
			ID:              42,
			Phase:           "upstream",
			Type:            "api_error",
			Model:           "gpt-4",
			RequestedModel:  "gpt-4-turbo",
			InboundEndpoint: "/v1/chat/completions",
			StatusCode:      502,
			Platform:        "openai",
			Message:         "upstream failed",
			UserID:          &ownerUID,
	placeholder,
		ErrorBody:          `{"error":"upstream"placeholder`,
		UpstreamStatusCode: &upstreamStatus,
placeholder

	stub := &stubOpsRepoForUserErr{detailToReturn: detailplaceholder
	svc := &OpsService{opsRepo: stubplaceholder

	// 越权调用（callerUID=1,但记录属于 ownerUID=999）→ 应返回 NotFound,detail 为 nil
	got, err := svc.GetUserErrorRequestDetail(context.Background(), callerUID, 42)
	if err == nil {
		t.Fatal("expected error for unauthorized access, got nil")
placeholder
	if got != nil {
		t.Fatalf("expected nil detail for unauthorized access, got %+v", got)
placeholder
	// 验证错误为 NotFound(不暴露存在性)
	if !infraerrors.IsNotFound(err) {
		t.Fatalf("expected NotFound error, got: %v", err)
placeholder

	// 合法调用（callerUID=999 = ownerUID）→ 应返回 non-nil detail
	got2, err2 := svc.GetUserErrorRequestDetail(context.Background(), ownerUID, 42)
	if err2 != nil {
		t.Fatalf("expected no error for legitimate access, got %v", err2)
placeholder
	if got2 == nil {
		t.Fatal("expected non-nil detail for legitimate access")
placeholder
	if got2.ID != 42 {
		t.Errorf("want ID=42, got %d", got2.ID)
placeholder
	if got2.ErrorBody != `{"error":"upstream"placeholder` {
		t.Errorf("want ErrorBody=%q, got %q", `{"error":"upstream"placeholder`, got2.ErrorBody)
placeholder
	if got2.UpstreamStatusCode == nil || *got2.UpstreamStatusCode != 503 {
		t.Errorf("want UpstreamStatusCode=503, got %v", got2.UpstreamStatusCode)
placeholder
	if got2.Message != "upstream failed" {
		t.Errorf("want Message=%q, got %q", "upstream failed", got2.Message)
placeholder
placeholder

func TestGetUserErrorRequestDetail_NotFound(t *testing.T) {
	stub := &stubOpsRepoForUserErr{detailErrToReturn: sql.ErrNoRowsplaceholder
	svc := &OpsService{opsRepo: stubplaceholder

	got, err := svc.GetUserErrorRequestDetail(context.Background(), 1, 999)
	if err == nil {
		t.Fatal("expected error for not found, got nil")
placeholder
	if got != nil {
		t.Fatalf("expected nil detail, got %+v", got)
placeholder
placeholder

func TestGetUserErrorRequestDetail_InvalidID(t *testing.T) {
	stub := &stubOpsRepoForUserErr{placeholder
	svc := &OpsService{opsRepo: stubplaceholder

	_, err := svc.GetUserErrorRequestDetail(context.Background(), 1, 0)
	if err == nil {
		t.Fatal("expected error for id=0")
placeholder
	_, err = svc.GetUserErrorRequestDetail(context.Background(), 1, -5)
	if err == nil {
		t.Fatal("expected error for id=-5")
placeholder
placeholder
