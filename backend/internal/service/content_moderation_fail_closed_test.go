package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// newFailClosedTestService 起一个把审计 API 指向"已关闭端口"的服务：callModeration
// 必定失败（连接被拒），用来稳定复现"风控后端异常"。
func newFailClosedTestService(t *testing.T, mutate func(cfg *ContentModerationConfig)) (*ContentModerationService, *contentModerationTestRepo) {
placeholder

	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {placeholder))
	deadURL := dead.URL
	dead.Close() // 端口关闭 -> 后续调用连接被拒

	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.Mode = ContentModerationModePreBlock
	cfg.BaseURL = deadURL
	cfg.APIKeys = []string{"sk-test"placeholder
	cfg.RetryCount = 0
	cfg.TimeoutMS = 500
	if mutate != nil {
		mutate(cfg)
placeholder
	rawCfg, err := json.Marshal(cfg)
placeholder

	repo := &contentModerationTestRepo{placeholder
	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		repo, &contentModerationTestHashCache{placeholder, nil, nil, nil, nil, nil,
	)
	return svc, repo
placeholder

func failClosedTestInput() ContentModerationCheckInput {
	return ContentModerationCheckInput{
		UserID:   1,
		Endpoint: "/v1/messages",
		Provider: "anthropic",
		Protocol: ContentModerationProtocolAnthropicMessages,
		Body:     []byte(`{"messages":[{"role":"user","content":"hello"placeholder]placeholder`),
placeholder
placeholder

// 回归护栏（#5394）：配置加载失败时必须放行。此时快照尚未建立，我们无从得知
// risk_control_enabled 与 fail_closed_on_error，无法确认运维的阻断意图；若拦截，
// 一次设置读取抖动就会让全部网关流量（含从未开启风控的部署）返回 5xx。
func TestContentModerationCheck_ConfigLoadFailureAllows(t *testing.T) {
	svc := &ContentModerationService{
		settingRepo: &contentModerationTestSettingRepo{
			values: map[string]string{SettingKeyRiskControlEnabled: "false"placeholder,
			err:    fmt.Errorf("settings unavailable"),
	placeholder,
		repo: &contentModerationTestRepo{placeholder,
placeholder

	decision, err := svc.Check(context.Background(), failClosedTestInput())

placeholder
	require.True(t, decision.Allowed, "配置加载失败不得阻断流量：意图未知时只能放行")
	require.False(t, decision.Blocked)
	require.Equal(t, ContentModerationActionAllow, decision.Action)
placeholder

// 损坏的 content_moderation_config JSON 是确定性失败，快照永远建不起来。若在此
// fail-closed，会造成永久性全站中断，且管理员被锁在修复界面之外。
func TestContentModerationCheck_MalformedConfigJSONAllows(t *testing.T) {
	svc := &ContentModerationService{
		settingRepo: &contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "true",
			SettingKeyContentModerationConfig: `{"enabled": tru`,
placeholder
		repo: &contentModerationTestRepo{placeholder,
placeholder

	for i := 0; i < 3; i++ {
		decision, err := svc.Check(context.Background(), failClosedTestInput())
	placeholder
		require.True(t, decision.Allowed, "第 %d 次请求：损坏配置不得演变成永久中断", i+1)
placeholder
placeholder

// 默认行为不变：审计 API 挂掉照旧放行，升级不会给存量部署带来可用性变化。
func TestContentModerationCheck_AuditAPIFailureFailsOpenByDefault(t *testing.T) {
	svc, repo := newFailClosedTestService(t, nil) // FailClosedOnError 默认 false

	decision, err := svc.Check(context.Background(), failClosedTestInput())

placeholder
	require.True(t, decision.Allowed, "默认必须 fail-open")
	require.False(t, decision.Blocked)
	require.Empty(t, repo.snapshotLogs(), "RecordNonHits 关闭时不写日志")
placeholder

// #5388 的真正修复：开关打开后，审计 API 异常必须拒绝请求。
func TestContentModerationCheck_AuditAPIFailureFailsClosedWhenEnabled(t *testing.T) {
	svc, repo := newFailClosedTestService(t, func(cfg *ContentModerationConfig) {
		cfg.FailClosedOnError = true
placeholder)

	decision, err := svc.Check(context.Background(), failClosedTestInput())

placeholder
	require.False(t, decision.Allowed)
	require.True(t, decision.Blocked, "开关打开后审计后端异常必须阻断")
	require.Equal(t, http.StatusServiceUnavailable, decision.StatusCode, "系统不可用用 503，不用 403/500")
	require.Equal(t, ContentModerationErrorCodeUnavailable, decision.ErrorCode)
	require.Equal(t, ContentModerationUnavailableMessage, decision.Message)
	require.Equal(t, ContentModerationActionError, decision.Action)

	require.False(t, decision.Flagged, "系统故障不是内容违规，Flagged 必须为 false")
	for _, log := range repo.snapshotLogs() {
		require.False(t, log.Flagged, "fail-closed 不得写 flagged 日志，否则会喂给自动封号计数")
placeholder
placeholder

// RecordNonHits 打开时仍会落一条 error 日志，但必须是非 flagged 的，
// 不能进入 CountFlaggedByUserSince 的封号计数。
func TestContentModerationCheck_FailClosedLogIsNotFlagged(t *testing.T) {
	svc, repo := newFailClosedTestService(t, func(cfg *ContentModerationConfig) {
		cfg.FailClosedOnError = true
		cfg.RecordNonHits = true
placeholder)

	decision, err := svc.Check(context.Background(), failClosedTestInput())
placeholder
	require.True(t, decision.Blocked)

	logs := repo.snapshotLogs()
	require.Len(t, logs, 1)
	require.False(t, logs[0].Flagged, "error 日志必须非 flagged")
	require.Equal(t, ContentModerationActionError, logs[0].Action)
placeholder

// observe 模式没有阻断语义，开关打开也不得拦截。
func TestContentModerationCheck_FailClosedIgnoredInObserveMode(t *testing.T) {
	svc, _ := newFailClosedTestService(t, func(cfg *ContentModerationConfig) {
		cfg.FailClosedOnError = true
		cfg.Mode = ContentModerationModeObserve
placeholder)

	decision, err := svc.Check(context.Background(), failClosedTestInput())

placeholder
	require.True(t, decision.Allowed, "observe 模式永不阻断")
placeholder

// 一个审计 Key 都没配：默认放行，开关打开则拒绝。
func TestContentModerationCheck_NoAuditKeysRespectsFailClosedSwitch(t *testing.T) {
	t.Run("默认放行", func(t *testing.T) {
		svc, _ := newFailClosedTestService(t, func(cfg *ContentModerationConfig) {
			cfg.APIKeys = []string{placeholder
	placeholder)
		decision, err := svc.Check(context.Background(), failClosedTestInput())
	placeholder
		require.True(t, decision.Allowed)
placeholder)

	t.Run("开关打开则拒绝", func(t *testing.T) {
		svc, _ := newFailClosedTestService(t, func(cfg *ContentModerationConfig) {
			cfg.APIKeys = []string{placeholder
			cfg.FailClosedOnError = true
	placeholder)
		decision, err := svc.Check(context.Background(), failClosedTestInput())
	placeholder
		require.True(t, decision.Blocked)
		require.Equal(t, http.StatusServiceUnavailable, decision.StatusCode)
		require.Equal(t, ContentModerationErrorCodeUnavailable, decision.ErrorCode)
placeholder)
placeholder

// 审计总开关关闭时，无论 fail_closed_on_error 如何都不得阻断。
func TestContentModerationCheck_FailClosedIgnoredWhenModerationDisabled(t *testing.T) {
	svc, _ := newFailClosedTestService(t, func(cfg *ContentModerationConfig) {
		cfg.Enabled = false
		cfg.FailClosedOnError = true
placeholder)

	decision, err := svc.Check(context.Background(), failClosedTestInput())

placeholder
	require.True(t, decision.Allowed, "审计未启用的部署永远不受 fail-closed 影响")
placeholder

// 风控总开关（risk_control_enabled）关闭时同样不得阻断。
func TestContentModerationCheck_FailClosedIgnoredWhenRiskControlDisabled(t *testing.T) {
	cfg := defaultContentModerationConfig()
	cfg.Enabled = true
	cfg.Mode = ContentModerationModePreBlock
	cfg.FailClosedOnError = true
	cfg.APIKeys = []string{placeholder
	rawCfg, err := json.Marshal(cfg)
placeholder

	svc := NewContentModerationService(
		&contentModerationTestSettingRepo{values: map[string]string{
			SettingKeyRiskControlEnabled:      "false",
			SettingKeyContentModerationConfig: string(rawCfg),
placeholder
		&contentModerationTestRepo{placeholder, &contentModerationTestHashCache{placeholder, nil, nil, nil, nil, nil,
	)

	decision, err := svc.Check(context.Background(), failClosedTestInput())

placeholder
	require.True(t, decision.Allowed, "风控中心未启用的部署永远不受 fail-closed 影响")
placeholder

func TestContentModerationConfig_FailClosedOnErrorMatrix(t *testing.T) {
	build := func(enabled, failClosed bool, mode string) *ContentModerationConfig {
		cfg := defaultContentModerationConfig()
		cfg.Enabled = enabled
		cfg.FailClosedOnError = failClosed
		cfg.Mode = mode
		return cfg
placeholder

	require.True(t, build(true, true, ContentModerationModePreBlock).failClosedOnError())
	require.False(t, build(true, false, ContentModerationModePreBlock).failClosedOnError(), "开关关闭")
	require.False(t, build(false, true, ContentModerationModePreBlock).failClosedOnError(), "审计未启用")
	require.False(t, build(true, true, ContentModerationModeObserve).failClosedOnError(), "observe 模式")
	require.False(t, build(true, true, ContentModerationModeOff).failClosedOnError(), "off 模式")

	var nilCfg *ContentModerationConfig
	require.False(t, nilCfg.failClosedOnError(), "nil 配置不得 panic")
placeholder

func TestContentModerationUpdateConfig_FailClosedOnError(t *testing.T) {
	repo := &contentModerationTestSettingRepo{values: map[string]string{placeholderplaceholder
	svc := NewContentModerationService(repo, &contentModerationTestRepo{placeholder, nil, nil, nil, nil, nil, nil)
	ctx := context.Background()

	view, err := svc.GetConfig(ctx)
placeholder
	require.False(t, view.FailClosedOnError, "默认必须 fail-open")

	enable := true
	view, err = svc.UpdateConfig(ctx, UpdateContentModerationConfigInput{FailClosedOnError: &enableplaceholder)
placeholder
	require.True(t, view.FailClosedOnError)

	saved := &ContentModerationConfig{placeholder
	require.NoError(t, json.Unmarshal([]byte(repo.values[SettingKeyContentModerationConfig]), saved))
	require.True(t, saved.FailClosedOnError, "必须持久化")

	// 未携带该字段的后续更新不得把它冲掉。
	view, err = svc.UpdateConfig(ctx, UpdateContentModerationConfigInput{placeholder)
placeholder
	require.True(t, view.FailClosedOnError)

	disable := false
	view, err = svc.UpdateConfig(ctx, UpdateContentModerationConfigInput{FailClosedOnError: &disableplaceholder)
placeholder
	require.False(t, view.FailClosedOnError)
placeholder
