package repository

import (
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestBuildOpsErrorLogsWhere_UserScopedFilters(t *testing.T) {
	uid := int64(42)
	kid := int64(7)
	filter := &service.OpsErrorLogFilter{
		UserID:             &uid,
		APIKeyID:           &kid,
		Model:              "claude-sonnet-4-5",
		ExcludeCountTokens: true,
		ErrorPhasesAny:     []string{"auth"placeholder,
		ErrorTypesAny:      []string{"rate_limit_error"placeholder,
		View:               "all",
placeholder
	where, args := buildOpsErrorLogsWhere(filter)

	for _, want := range []string{
		"e.user_id = $",
		"e.api_key_id = $",
		"COALESCE(e.requested_model, e.model, '') = $",
		"COALESCE(e.is_count_tokens, false) = false",
		"e.error_phase = ANY($",
		"e.error_type = ANY($",
placeholder {
		if !strings.Contains(where, want) {
			t.Fatalf("where missing %q\nfull: %s", want, where)
	placeholder
placeholder
	if len(args) != 5 {
		t.Fatalf("expected 5 args, got %d", len(args))
placeholder
placeholder

func TestBuildOpsErrorLogsWhere_ModelFuzzy(t *testing.T) {
	// 默认（ModelFuzzy=false）保持精确匹配
	exact := &service.OpsErrorLogFilter{Model: "claude"placeholder
	whereExact, _ := buildOpsErrorLogsWhere(exact)
	if !strings.Contains(whereExact, "COALESCE(e.requested_model, e.model, '') = $") {
		t.Fatalf("default should be exact match, got: %s", whereExact)
placeholder

	// ModelFuzzy=true → ILIKE
	fuzzy := &service.OpsErrorLogFilter{Model: "claude", ModelFuzzy: trueplaceholder
	whereFuzzy, args := buildOpsErrorLogsWhere(fuzzy)
	if !strings.Contains(whereFuzzy, "COALESCE(e.requested_model, e.model, '') ILIKE $") {
		t.Fatalf("ModelFuzzy should use ILIKE, got: %s", whereFuzzy)
placeholder
	if len(args) != 1 || args[0] != "%claude%" {
		t.Fatalf("expected arg \"%%claude%%\", got %v", args)
placeholder

	// 通配符转义：输入含 % 应被转义为字面量
	esc := &service.OpsErrorLogFilter{Model: "50%off", ModelFuzzy: trueplaceholder
	_, escArgs := buildOpsErrorLogsWhere(esc)
	if len(escArgs) != 1 || escArgs[0] != `%50\%off%` {
		t.Fatalf("expected escaped arg, got %v", escArgs)
placeholder

	esc2 := &service.OpsErrorLogFilter{Model: "gpt_4o", ModelFuzzy: trueplaceholder
	_, escArgs2 := buildOpsErrorLogsWhere(esc2)
	if len(escArgs2) != 1 || escArgs2[0] != `%gpt\_4o%` {
		t.Fatalf("underscore should be escaped, got %v", escArgs2)
placeholder
placeholder

func TestBuildOpsErrorLogsWhere_MatchDeletedKeyOwner(t *testing.T) {
	uid := int64(42)

	// 开关开启 → 归属放宽为 OR(user_id 或 deleted_key_owner_user_id),且共用同一占位符
	on := &service.OpsErrorLogFilter{UserID: &uid, MatchDeletedKeyOwner: trueplaceholder
	whereOn, argsOn := buildOpsErrorLogsWhere(on)
	if !strings.Contains(whereOn, "(e.user_id = $1 OR e.deleted_key_owner_user_id = $1)") {
		t.Fatalf("MatchDeletedKeyOwner=true should widen to OR, got: %s", whereOn)
placeholder
	if len(argsOn) != 1 || argsOn[0] != uid {
		t.Fatalf("expected single reused arg %d, got %v", uid, argsOn)
placeholder

	// 开关关闭(默认)→ 仅精确 user_id,绝不出现 deleted_key_owner_user_id(admin 回归)
	off := &service.OpsErrorLogFilter{UserID: &uidplaceholder
	whereOff, _ := buildOpsErrorLogsWhere(off)
	if !strings.Contains(whereOff, "e.user_id = $1") {
		t.Fatalf("default should match user_id exactly, got: %s", whereOff)
placeholder
	if strings.Contains(whereOff, "deleted_key_owner_user_id") {
		t.Fatalf("default must NOT include deleted_key_owner_user_id, got: %s", whereOff)
placeholder
placeholder
