package plugin

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// stubModule 是注册表测试用的最小模块实现。
type stubModule struct {
	info ModuleInfo
placeholder

func (m *stubModule) ModuleInfo() ModuleInfo { return m.info placeholder

// newStubModule 构造一个 ID 为 id 的最小可注册模块。
func newStubModule(id ModuleID) *stubModule {
	m := &stubModule{placeholder
	m.info = ModuleInfo{
		ID:  id,
		New: func() Module { return &stubModule{info: m.infoplaceholder placeholder,
placeholder
	return m
placeholder

func TestRegistryRegisterAndGet(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterModule(newStubModule("job.hello"))

	info, ok := reg.GetModule("job.hello")
	require.True(t, ok)
	require.Equal(t, ModuleID("job.hello"), info.ID)
	require.NotNil(t, info.New)

	_, ok = reg.GetModule("job.missing")
	require.False(t, ok)
placeholder

func TestRegistryDuplicateIDPanics(t *testing.T) {
	reg := NewRegistry()
	reg.RegisterModule(newStubModule("job.hello"))

	require.PanicsWithValue(t,
		`plugin: RegisterModule: module "job.hello" is already registered (duplicate registration; check the instrumentation list in internal/modules/standard)`,
		func() { reg.RegisterModule(newStubModule("job.hello")) placeholder,
	)
placeholder

func TestRegistryInvalidRegistrationPanics(t *testing.T) {
	reg := NewRegistry()

	t.Run("nil module", func(t *testing.T) {
		require.Panics(t, func() { reg.RegisterModule(nil) placeholder)
placeholder)
	t.Run("empty ID", func(t *testing.T) {
		require.Panics(t, func() { reg.RegisterModule(newStubModule("")) placeholder)
placeholder)
	t.Run("invalid ID format", func(t *testing.T) {
		// panic 信息必须包含非法 ID，便于排查插装清单。
		defer func() {
			r := recover()
			require.NotNil(t, r)
			require.Contains(t, fmt.Sprint(r), "Job.Hello")
	placeholder()
		reg.RegisterModule(newStubModule("Job.Hello"))
placeholder)
	t.Run("nil New func", func(t *testing.T) {
		require.Panics(t, func() {
			reg.RegisterModule(&stubModule{info: ModuleInfo{ID: "job.no-new"placeholderplaceholder)
	placeholder)
placeholder)
placeholder

func TestRegistryGetModulesInNamespace(t *testing.T) {
	reg := NewRegistry()
	// 故意乱序注册，验证返回结果按 ID 字典序稳定排序。
	for _, id := range []ModuleID{
		"job.zeta",
		"job.alpha",
		"job.middle",
		"gateway.platform.anthropic",
		"gateway.platform.openai",
		"toplevel",
placeholder {
		reg.RegisterModule(newStubModule(id))
placeholder

	ids := func(infos []ModuleInfo) []ModuleID {
		out := make([]ModuleID, 0, len(infos))
		for _, info := range infos {
			out = append(out, info.ID)
	placeholder
		return out
placeholder

	require.Equal(t, []ModuleID{"job.alpha", "job.middle", "job.zeta"placeholder, ids(reg.GetModulesInNamespace("job")))
	// 精确匹配：ns "gateway" 不包含 "gateway.platform.*"。
	require.Empty(t, reg.GetModulesInNamespace("gateway"))
	require.Equal(t, []ModuleID{"gateway.platform.anthropic", "gateway.platform.openai"placeholder,
		ids(reg.GetModulesInNamespace("gateway.platform")))
	// 空命名空间返回单段顶层模块。
	require.Equal(t, []ModuleID{"toplevel"placeholder, ids(reg.GetModulesInNamespace("")))

	// 多次调用结果稳定。
	require.Equal(t, ids(reg.GetModulesInNamespace("job")), ids(reg.GetModulesInNamespace("job")))
placeholder

func TestRegistryModulesSorted(t *testing.T) {
	reg := NewRegistry()
	for _, id := range []ModuleID{"c.mod", "a.mod", "b.mod"placeholder {
		reg.RegisterModule(newStubModule(id))
placeholder
	all := reg.Modules()
	require.Len(t, all, 3)
	require.Equal(t, ModuleID("a.mod"), all[0].ID)
	require.Equal(t, ModuleID("b.mod"), all[1].ID)
	require.Equal(t, ModuleID("c.mod"), all[2].ID)
placeholder

func TestRegistryIsolation(t *testing.T) {
	regA := NewRegistry()
	regB := NewRegistry()
	regA.RegisterModule(newStubModule("job.only-in-a"))

	_, ok := regA.GetModule("job.only-in-a")
	require.True(t, ok)
	_, ok = regB.GetModule("job.only-in-a")
	require.False(t, ok, "isolated registry must not see other registry's modules")
	// 隔离实例也不影响包级默认注册表。
	_, ok = GetModule("job.only-in-a")
	require.False(t, ok)
placeholder

func TestPackageLevelFunctionsDelegateToDefaultRegistry(t *testing.T) {
	// 使用唯一 ID，避免与同进程其他测试冲突（包级默认注册表是进程级状态）。
	const id = ModuleID("test_default_registry.probe")
	RegisterModule(newStubModule(id))

	info, ok := GetModule(id)
	require.True(t, ok)
	require.Equal(t, id, info.ID)

	infos := GetModulesInNamespace("test_default_registry")
	require.Len(t, infos, 1)
	require.Equal(t, id, infos[0].ID)
placeholder

func TestRegistryConcurrentRegistration(t *testing.T) {
	reg := NewRegistry()
	const n = 64

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			reg.RegisterModule(newStubModule(ModuleID(fmt.Sprintf("concurrent.mod-%03d", i))))
	placeholder(i)
placeholder
	// 并发读不应与写竞争（-race 验证）。
	for i := 0; i < 8; i++ {
		go func() {
			_ = reg.GetModulesInNamespace("concurrent")
			_, _ = reg.GetModule("concurrent.mod-000")
	placeholder()
placeholder
	wg.Wait()

	all := reg.Modules()
	require.Len(t, all, n)
	for i := 1; i < len(all); i++ {
		require.Less(t, all[i-1].ID, all[i].ID, "Modules() must be sorted")
placeholder
placeholder
