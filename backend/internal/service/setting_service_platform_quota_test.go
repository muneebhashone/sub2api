//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMergePlatformQuotaDefaults_PatchSemantics(t *testing.T) {
	five := 5.0
	base := DefaultPlatformQuotaSetting{
		DailyLimitUSD:  &five,
		WeeklyLimitUSD: &five,
placeholder
	ten := 10.0
	patch := DefaultPlatformQuotaSetting{DailyLimitUSD: &tenplaceholder

	mergePlatformQuotaDefaults(&base, &patch)
	if base.DailyLimitUSD == nil || *base.DailyLimitUSD != 10.0 {
		t.Errorf("daily not patched: %+v", base.DailyLimitUSD)
placeholder
	if base.WeeklyLimitUSD == nil || *base.WeeklyLimitUSD != 5.0 {
		t.Errorf("weekly should remain 5.0: %+v", base.WeeklyLimitUSD)
placeholder
placeholder

func TestMergePlatformQuotaDefaults_ZeroIsExplicitDisable(t *testing.T) {
	five := 5.0
	base := DefaultPlatformQuotaSetting{DailyLimitUSD: &fiveplaceholder
	zero := 0.0
	patch := DefaultPlatformQuotaSetting{DailyLimitUSD: &zeroplaceholder

	mergePlatformQuotaDefaults(&base, &patch)
	if base.DailyLimitUSD == nil || *base.DailyLimitUSD != 0 {
		t.Errorf("explicit 0 should patch base, got %+v", base.DailyLimitUSD)
placeholder
placeholder

func TestMergePlatformQuotaDefaults_NilSrcIsNoop(t *testing.T) {
	five := 5.0
	base := DefaultPlatformQuotaSetting{DailyLimitUSD: &fiveplaceholder
	mergePlatformQuotaDefaults(&base, nil)
	if base.DailyLimitUSD == nil || *base.DailyLimitUSD != 5.0 {
		t.Errorf("nil src should be no-op: %+v", base.DailyLimitUSD)
placeholder
placeholder

func floatPtrPQ(v float64) *float64 { return &v placeholder

func newSettingServiceForPlatformQuotaTest(seed map[string]string) *SettingService {
	repo := newMockSettingRepo()
	for k, v := range seed {
		repo.data[k] = v
placeholder
	return NewSettingService(repo, &config.Config{placeholder)
placeholder

func TestGetDefaultPlatformQuotas_ReturnsAllowedPlatforms(t *testing.T) {
	zero := 0.0
	svc := newSettingServiceForPlatformQuotaTest(map[string]string{
		// 新 JSON 格式：anthropic daily=10.5, openai monthly=0, 其他平台无配置
		SettingKeyDefaultPlatformQuotas: `{"anthropic":{"daily":placeholder,"openai":{"monthly":0placeholderplaceholder`,
placeholder)
	got, err := svc.GetDefaultPlatformQuotas(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
placeholder
	// 必须包含全部允许 platform key（补齐契约）
	for _, platform := range AllowedQuotaPlatforms {
		if _, ok := got[platform]; !ok {
			t.Errorf("missing platform key: %q", platform)
	placeholder
placeholder
	// anthropic daily = 10.5
	if v := got["anthropic"].DailyLimitUSD; v == nil || *v != 10.5 {
		t.Errorf("anthropic daily want 10.5, got %v", v)
placeholder
	// openai monthly = 0（显式禁用）
	if v := got["openai"].MonthlyLimitUSD; v == nil || *v != zero {
		t.Errorf("openai monthly want 0 (explicit disable), got %v", v)
placeholder
	// gemini 无配置 → weekly = nil
	if v := got["gemini"].WeeklyLimitUSD; v != nil {
		t.Errorf("gemini weekly want nil (not configured), got %v", *v)
placeholder
	// antigravity 无配置 → daily = nil
	if v := got["antigravity"].DailyLimitUSD; v != nil {
		t.Errorf("antigravity daily want nil (not configured), got %v", *v)
placeholder
placeholder

func TestGetAuthSourcePlatformQuotas_OnlyConfiguredReturned(t *testing.T) {
	source := "email"
	// 新 JSON 格式：anthropic daily=5, monthly=100；openai weekly=0；gemini/antigravity 无配置
	svc := newSettingServiceForPlatformQuotaTest(map[string]string{
		SettingKeyAuthSourcePlatformQuotas(source): `{"anthropic":{"daily":5,"monthly":100placeholder,"openai":{"weekly":0placeholderplaceholder`,
placeholder)
	got := svc.GetAuthSourcePlatformQuotas(context.Background(), source)

	// anthropic 有配置 → 在结果中
	anthro, ok := got["anthropic"]
	if !ok {
		t.Fatal("expected anthropic to be present")
placeholder
	if anthro.DailyLimitUSD == nil || *anthro.DailyLimitUSD != 5.0 {
		t.Errorf("anthropic daily want 5.0, got %v", anthro.DailyLimitUSD)
placeholder
	if anthro.MonthlyLimitUSD == nil || *anthro.MonthlyLimitUSD != 100.0 {
		t.Errorf("anthropic monthly want 100.0, got %v", anthro.MonthlyLimitUSD)
placeholder
	if anthro.WeeklyLimitUSD != nil {
		t.Errorf("anthropic weekly not configured, want nil, got %v", *anthro.WeeklyLimitUSD)
placeholder

	// openai weekly=0 → 在结果中
	oai, ok := got["openai"]
	if !ok {
		t.Fatal("expected openai to be present")
placeholder
	if oai.WeeklyLimitUSD == nil || *oai.WeeklyLimitUSD != 0 {
		t.Errorf("openai weekly want 0, got %v", oai.WeeklyLimitUSD)
placeholder

	// gemini / antigravity 无配置 → 不在结果中（override 语义）
	if _, ok := got["gemini"]; ok {
		t.Error("gemini not configured, should be absent from result")
placeholder
	if _, ok := got["antigravity"]; ok {
		t.Error("antigravity not configured, should be absent from result")
placeholder
placeholder

func TestGetAuthSourcePlatformQuotas_AllNegativeOrEmpty_NoEntry(t *testing.T) {
	source := "linuxdo"
	// 新 JSON 格式：未配置任何平台（空 JSON key）→ 返回空 map
	svc := newSettingServiceForPlatformQuotaTest(map[string]string{
		SettingKeyAuthSourcePlatformQuotas(source): `{placeholder`,
placeholder)
	got := svc.GetAuthSourcePlatformQuotas(context.Background(), source)
	// 空 map → override 语义，无 openai 条目
	if _, ok := got["openai"]; ok {
		t.Error("empty JSON object should result in no openai entry")
placeholder
	if len(got) != 0 {
		t.Errorf("expected empty result map, got %v", got)
placeholder
placeholder

// TestSystemPlatformQuotas_WriteReadRoundTrip 验证系统层 platform quota 经 buildSystemSettingsUpdates（写）
// 再由 GetDefaultPlatformQuotas（读）正确往返——覆盖真实 write→read 路径，锁住 4-key 补齐契约。
func TestSystemPlatformQuotas_WriteReadRoundTrip(t *testing.T) {
	svc := newSettingServiceForPlatformQuotaTest(nil)
	ctx := context.Background()

	ten := 10.0
	ss := &SystemSettings{
		DefaultPlatformQuotas: map[string]*DefaultPlatformQuotaSetting{
			"anthropic": {DailyLimitUSD: &ten, WeeklyLimitUSD: nil, MonthlyLimitUSD: nilplaceholder,
	placeholder,
placeholder
	if err := svc.UpdateSettings(ctx, ss); err != nil {
		t.Fatalf("UpdateSettings: %v", err)
placeholder

	got, err := svc.GetDefaultPlatformQuotas(ctx)
	if err != nil {
		t.Fatal(err)
placeholder
	// 4-key 补齐契约：无论写了几个 platform，读回必须含全部 4 个
	for _, p := range []string{"anthropic", "openai", "gemini", "antigravity"placeholder {
		if _, ok := got[p]; !ok {
			t.Errorf("4-key contract violated: missing platform %q", p)
	placeholder
placeholder
	// 写入值正确往返
	if v := got["anthropic"].DailyLimitUSD; v == nil || *v != ten {
		t.Fatalf("anthropic daily round-trip failed: got %v, want 10", v)
placeholder
	// 未写入的平台字段为 nil
	if got["openai"].DailyLimitUSD != nil {
		t.Errorf("openai daily should be nil (not written), got %v", got["openai"].DailyLimitUSD)
placeholder
placeholder

// TestSystemPlatformQuotas_EmptyMapClearsAll 验证空 map 的整体替换语义：
// 写入 DefaultPlatformQuotas={placeholder 后，GetDefaultPlatformQuotas 返回全部允许平台、所有字段均为 nil，
// 明确文档化"空 map = 清空全部配额"是有意为之的 whole-replace 语义。
func TestSystemPlatformQuotas_EmptyMapClearsAll(t *testing.T) {
	svc := newSettingServiceForPlatformQuotaTest(nil)
	ctx := context.Background()

	// 先写入有值的配置
	ten := 10.0
	if err := svc.UpdateSettings(ctx, &SystemSettings{
		DefaultPlatformQuotas: map[string]*DefaultPlatformQuotaSetting{
			"anthropic": {DailyLimitUSD: &tenplaceholder,
	placeholder,
placeholder); err != nil {
		t.Fatalf("initial write: %v", err)
placeholder

	// 再写入空 map（整体替换语义：清空全部）
	if err := svc.UpdateSettings(ctx, &SystemSettings{
		DefaultPlatformQuotas: map[string]*DefaultPlatformQuotaSetting{placeholder,
placeholder); err != nil {
		t.Fatalf("empty map write: %v", err)
placeholder

	got, err := svc.GetDefaultPlatformQuotas(ctx)
	if err != nil {
		t.Fatal(err)
placeholder
	// 全部允许平台 key 仍然存在（补齐契约）
	for _, p := range AllowedQuotaPlatforms {
		if _, ok := got[p]; !ok {
			t.Errorf("allowed-platform contract violated after empty write: missing %q", p)
	placeholder
placeholder
	// 所有字段 nil（全部已清空）
	for _, p := range AllowedQuotaPlatforms {
		pq := got[p]
		if pq == nil {
			continue
	placeholder
		if pq.DailyLimitUSD != nil || pq.WeeklyLimitUSD != nil || pq.MonthlyLimitUSD != nil {
			t.Errorf("platform %q should have all-nil limits after empty-map write, got %+v", p, pq)
	placeholder
placeholder
placeholder

// TestUpdateSettingsWithAuthSourceDefaults_PlatformQuotaRoundTrip 验证 round-4 fix：
// PUT /admin/settings 携带的 auth source × platform × window 限额能完整写入并被 GetAuthSourcePlatformQuotas 读回。
// Round-4 之前 writeProviderDefaultGrantUpdates 完全没写 PQ key，前端配置静默丢失。
func TestUpdateSettingsWithAuthSourceDefaults_PlatformQuotaRoundTrip(t *testing.T) {
	svc := newSettingServiceForPlatformQuotaTest(nil)
	systemSettings := &SystemSettings{placeholder
	authDefaults := &AuthSourceDefaultSettings{
		Email: ProviderDefaultGrantSettings{
			PlatformQuotas: map[string]*DefaultPlatformQuotaSetting{
				"anthropic": {
					DailyLimitUSD:   floatPtrPQ(5.0),
					WeeklyLimitUSD:  nil, // 无限额
					MonthlyLimitUSD: floatPtrPQ(100.0),
			placeholder,
				"openai": {
					DailyLimitUSD: floatPtrPQ(0), // 显式禁用
			placeholder,
		placeholder,
	placeholder,
placeholder
	if err := svc.UpdateSettingsWithAuthSourceDefaults(context.Background(), systemSettings, authDefaults); err != nil {
		t.Fatalf("UpdateSettingsWithAuthSourceDefaults: %v", err)
placeholder
	got := svc.GetAuthSourcePlatformQuotas(context.Background(), "email")
	anthro := got["anthropic"]
	if anthro == nil || anthro.DailyLimitUSD == nil || *anthro.DailyLimitUSD != 5.0 {
		t.Errorf("anthropic daily round-trip failed: %+v", anthro)
placeholder
	if anthro != nil && anthro.WeeklyLimitUSD != nil {
		t.Errorf("anthropic weekly want nil (无限额), got %v", *anthro.WeeklyLimitUSD)
placeholder
	if anthro == nil || anthro.MonthlyLimitUSD == nil || *anthro.MonthlyLimitUSD != 100.0 {
		t.Errorf("anthropic monthly round-trip failed: %+v", anthro)
placeholder
	oai := got["openai"]
	if oai == nil || oai.DailyLimitUSD == nil || *oai.DailyLimitUSD != 0 {
		t.Errorf("openai daily=0 (禁用) round-trip failed: %+v", oai)
placeholder
	// 其他 source 不应有 quota（authDefaults 只填了 Email）
	if linux := svc.GetAuthSourcePlatformQuotas(context.Background(), "linuxdo"); len(linux) != 0 {
		t.Errorf("linuxdo should be empty, got %+v", linux)
placeholder
placeholder

// TestUpdateSettingsWithAuthSourceDefaults_NilPlatformQuotaPreservesExisting 验证 #2 防御：
// 请求未携带某 auth source 的 platform quota（nil）时跳过写入、保留既有配置，
// 而非整体替换为空 map 清空（与系统层 nil 守卫一致）。
func TestUpdateSettingsWithAuthSourceDefaults_NilPlatformQuotaPreservesExisting(t *testing.T) {
	svc := newSettingServiceForPlatformQuotaTest(map[string]string{
		SettingKeyAuthSourcePlatformQuotas("email"): `{"anthropic":{"daily":5,"weekly":null,"monthly":nullplaceholderplaceholder`,
placeholder)
	// authDefaults 不携带 Email 的 PlatformQuotas（nil）——应保留既有配置
	authDefaults := &AuthSourceDefaultSettings{
		Email: ProviderDefaultGrantSettings{PlatformQuotas: nilplaceholder,
placeholder
	if err := svc.UpdateSettingsWithAuthSourceDefaults(context.Background(), &SystemSettings{placeholder, authDefaults); err != nil {
		t.Fatalf("UpdateSettingsWithAuthSourceDefaults: %v", err)
placeholder
	anthro := svc.GetAuthSourcePlatformQuotas(context.Background(), "email")["anthropic"]
	if anthro == nil || anthro.DailyLimitUSD == nil || *anthro.DailyLimitUSD != 5.0 {
		t.Errorf("nil PlatformQuotas 应保留既有 anthropic daily=5，got %+v", anthro)
placeholder
placeholder

// TestGetAuthSourcePlatformQuotas_JSON 验证新 JSON key 读写语义：
// 写入 JSON，断言已配置平台在结果中、未配置平台不在结果中（override 语义）。
func TestGetAuthSourcePlatformQuotas_JSON(t *testing.T) {
	svc := newSettingServiceForPlatformQuotaTest(map[string]string{
		SettingKeyAuthSourcePlatformQuotas("email"): `{"openai":{"daily":null,"weekly":null,"monthly":20placeholderplaceholder`,
placeholder)
	got := svc.GetAuthSourcePlatformQuotas(context.Background(), "email")

	// openai monthly = 20
	oai, ok := got["openai"]
	if !ok {
		t.Fatal("expected openai to be present")
placeholder
	if oai.MonthlyLimitUSD == nil || *oai.MonthlyLimitUSD != 20 {
		t.Errorf("openai monthly want 20, got %v", oai.MonthlyLimitUSD)
placeholder
	if oai.DailyLimitUSD != nil {
		t.Errorf("openai daily want nil, got %v", *oai.DailyLimitUSD)
placeholder
	if oai.WeeklyLimitUSD != nil {
		t.Errorf("openai weekly want nil, got %v", *oai.WeeklyLimitUSD)
placeholder

	// anthropic 未配置 → 不在结果中（override 语义）
	if _, ok := got["anthropic"]; ok {
		t.Error("anthropic not configured, should be absent from result")
placeholder
placeholder

// TestUpdateSettingsWithAuthSourceDefaults_NegativeQuotaRejected 验证改动 C：
// auth-source platform quota 含负数时，UpdateSettingsWithAuthSourceDefaults 返回 BadRequest 错误。
func TestUpdateSettingsWithAuthSourceDefaults_NegativeQuotaRejected(t *testing.T) {
	svc := newSettingServiceForPlatformQuotaTest(nil)
	neg := -1.0
	authDefaults := &AuthSourceDefaultSettings{
		Email: ProviderDefaultGrantSettings{
			PlatformQuotas: map[string]*DefaultPlatformQuotaSetting{
				"anthropic": {DailyLimitUSD: &negplaceholder,
		placeholder,
	placeholder,
placeholder
	err := svc.UpdateSettingsWithAuthSourceDefaults(context.Background(), &SystemSettings{placeholder, authDefaults)
	require.Error(t, err, "expected error for negative quota")
	require.Equal(t, "INVALID_DEFAULT_PLATFORM_QUOTA", infraerrors.Reason(err))
placeholder
