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

// TestBuildOpsErrorLogsWhere_CyberPolicyStatusExemption verifies that streaming
// cyber_policy hits (status_code=200) remain visible in admin + user error-request
// lists.  The repository filter must emit an OR exemption for error_type='cyber_policy'
// so that stream-path cyber rows (upstream delivers 200 with a failed SSE event) are
// not silently excluded by the COALESCE(status_code,0) >= 400 guard.
func TestBuildOpsErrorLogsWhere_CyberPolicyStatusExemption(t *testing.T) {
	// Default filter (no phase) must include the cyber_policy exemption.
	where, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{placeholder)
	if !strings.Contains(where, "e.error_type = 'cyber_policy'") {
		t.Fatalf("default filter must exempt cyber_policy from status >= 400 guard\nfull: %s", where)
placeholder
	if !strings.Contains(where, "COALESCE(e.status_code, 0) >= 400") {
		t.Fatalf("default filter must still include the status >= 400 guard for non-cyber rows\nfull: %s", where)
placeholder

	// phase=upstream WITHOUT the recovered-upstream opt-in keeps the status guard:
	// request-error list endpoints filter by phase=upstream as a plain condition.
	whereUpstream, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{Phase: "upstream"placeholder)
	if !strings.Contains(whereUpstream, "COALESCE(e.status_code, 0) >= 400") {
		t.Fatalf("upstream phase without IncludeRecoveredUpstream must keep the status guard\nfull: %s", whereUpstream)
placeholder
	if !strings.Contains(whereUpstream, "e.error_phase = $") {
		t.Fatalf("upstream phase filter must emit the error_phase condition\nfull: %s", whereUpstream)
placeholder

	// phase=upstream WITH IncludeRecoveredUpstream (ops 上游列表) skips the guard,
	// exposing recovered (<400) upstream rows.
	whereRecovered, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{Phase: "upstream", IncludeRecoveredUpstream: trueplaceholder)
	if strings.Contains(whereRecovered, "status_code") {
		t.Fatalf("upstream phase with IncludeRecoveredUpstream must not add any status_code clause\nfull: %s", whereRecovered)
placeholder

	// account_auth uses the same explicit provider-health opt-in but remains a
	// distinct phase from inference upstream errors.
	whereAccountAuth, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{Phase: "account_auth", IncludeRecoveredUpstream: trueplaceholder)
	if strings.Contains(whereAccountAuth, "status_code") {
		t.Fatalf("account_auth phase with IncludeRecoveredUpstream must expose recovered rows\nfull: %s", whereAccountAuth)
placeholder
	if !strings.Contains(whereAccountAuth, "e.error_phase = $") {
		t.Fatalf("account_auth recovered filter must retain its explicit phase\nfull: %s", whereAccountAuth)
placeholder

	whereProviderHealth, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{
		ErrorPhasesAny:           []string{"upstream", "account_auth"placeholder,
		IncludeRecoveredUpstream: true,
placeholder)
	if strings.Contains(whereProviderHealth, "status_code") {
		t.Fatalf("provider-health ANY filter must expose recovered inference and credential rows\nfull: %s", whereProviderHealth)
placeholder
	if !strings.Contains(whereProviderHealth, "e.error_phase = ANY($") {
		t.Fatalf("provider-health filter must preserve distinct phase values\nfull: %s", whereProviderHealth)
placeholder

	whereUserAccountAuth, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{ErrorPhasesAny: []string{"account_auth"placeholderplaceholder)
	if !strings.Contains(whereUserAccountAuth, "COALESCE(e.status_code, 0) >= 400") {
		t.Fatalf("request-error account_auth filters must exclude recovered successes\nfull: %s", whereUserAccountAuth)
placeholder

	whereMixed, _ := buildOpsErrorLogsWhere(&service.OpsErrorLogFilter{
		ErrorPhasesAny:           []string{"account_auth", "request"placeholder,
		IncludeRecoveredUpstream: true,
placeholder)
	if !strings.Contains(whereMixed, "COALESCE(e.status_code, 0) >= 400") {
		t.Fatalf("recovered opt-in must not bypass the guard for non-provider phases\nfull: %s", whereMixed)
placeholder
placeholder

func TestBuildOpsErrorLogsWhere_UserOwnershipIsDirectOnly(t *testing.T) {
	uid := int64(42)
	filter := &service.OpsErrorLogFilter{UserID: &uidplaceholder
	where, args := buildOpsErrorLogsWhere(filter)
	if !strings.Contains(where, "e.user_id = $1") {
		t.Fatalf("user scope should match user_id exactly, got: %s", where)
placeholder
	if len(args) != 1 || args[0] != uid {
		t.Fatalf("expected user id arg %d, got %v", uid, args)
placeholder
	if strings.Contains(where, "deleted_key_owner_user_id") {
		t.Fatalf("user ownership must not depend on deleted-key attribution: %s", where)
placeholder
placeholder
