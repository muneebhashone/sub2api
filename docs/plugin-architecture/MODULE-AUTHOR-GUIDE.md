# Sub2API 模块作者指南（Module Author Guide）

> 本指南面向需要编写或迁移插件模块的开发者，内容以 `backend/internal/plugin/` 已落地代码为准。
> 配套文档：[插件化改造 ROADMAP](../../.claude/plugin-refactor/ROADMAP.md)（命名空间规划与全程铁律）、
> [回归不变量清单 INVARIANTS](../../.claude/plugin-refactor/INVARIANTS.md)（每 PR 必须保持的外部行为）。

## 1. 概述与设计理念

Sub2API 的插件系统是 **Caddy 式进程内模块系统**：模块与宿主编译进同一个二进制，没有动态加载。内核由四部分组成——**命名空间化模块注册表**（模块包在 `init()` 期自注册，编译期插装清单决定哪些模块被编译进来）、**`modules:` 配置子树**（每模块的 enabled 开关 + 私有配置）、**宿主能力面 Host**（模块可使用的全部宿主能力，只暴露 ports 级接口）、**生命周期驱动器 Runtime**（按稳定序驱动 Provision → Validate → Start → Stop）。注册表只是"目录"，模块的实例化与生命周期统一由 Runtime 驱动。引入新模块在缺省配置下必须是零行为变更。

内核代码全部位于 `backend/internal/plugin/`：

| 文件 | 职责 |
|---|---|
| `module.go` | `ModuleID` / `ModuleInfo` / `Module` 与 4 个可选生命周期接口 |
| `registry.go` | 注册表：`RegisterModule` / `GetModule` / `GetModulesInNamespace`，可实例化的 `Registry` |
| `host.go` | 宿主能力面 `Host`（Logger / ConfigOf / DB / Redis） |
| `config.go` | `modules:` 配置子树解析（`ParseConfig` / `Config.Of`） |
| `runtime.go` | 生命周期驱动器 `Runtime`（Build / Start / Stop / Snapshot） |
| `wire.go` | Wire 装配（ProvideModuleConfig / ProvideHost / ProvideModuleRuntime） |

## 2. 模块结构与命名空间

### 2.1 ModuleID 规则

`ModuleID` 是模块的全局唯一标识，采用点分层级命名，形如 `job.hello`、`gateway.platform.anthropic`：

- 最后一段为模块名（`Name()`），之前的部分为命名空间（`Namespace()`）；
- 每一段只允许 `[a-z0-9_-]` 字符且不能为空；
- 单段 ID（如 `hello`）的命名空间为空字符串；
- `GetModulesInNamespace(ns)` 为**精确匹配**（非前缀匹配）：`ns` 为 `gateway.platform` 时返回 `gateway.platform.anthropic` 等直接子模块，不包含 `gateway.platform.x.y`。

ID 格式非法（空段、非法字符）在注册时直接 panic，在配置解析时返回错误，均为 fail-fast。

### 2.2 命名空间规划

命名空间在 [ROADMAP](../../.claude/plugin-refactor/ROADMAP.md) 中一次定义、分阶段填充。新模块必须落入已规划的命名空间；需要新命名空间时先在 ROADMAP 登记：

| 命名空间 | 宿主接口 | 迁移阶段 | 状态 |
|---|---|---|---|
| `job.*` | Starter/Stopper | Phase 1 示例 / Phase 2 可选 | **已实现**（示例模块 `job.hello`） |
| `payment.provider.*` | 现有 `payment.Provider` | Phase 2 | 规划中 |
| `gateway.hook.*` | pre/post-flight 管道钩子 | Phase 2 | 规划中 |
| `gateway.platform.*` | PlatformProvider | Phase 3 | 规划中 |
| `channel.*` / `notify.*` | 数据契约 / Notifier | Phase 4 | 规划中 |
| `plugin.bridge.*` | gRPC 进程外桥 | Phase 5（可选） | 规划中 |

### 2.3 模块的最低要求

模块是一个实现 `plugin.Module` 接口的 struct，放在 `backend/internal/modules/<模块名>/` 包内：

```go
// Module 是所有插件模块必须实现的最小接口。
type Module interface {
	// ModuleInfo 返回模块的注册信息。该方法必须无副作用且可在零值实例上调用。
	ModuleInfo() ModuleInfo
placeholder
```

`ModuleInfo` 是注册表中的一条记录：

```go
type ModuleInfo struct {
	// ID 是模块的全局唯一标识。
	ID ModuleID

	// New 返回该模块的一个新实例。
	// 每次 Runtime.Build 都会通过 New 创建全新实例，模块不应依赖包级可变状态。
	New func() Module

	// EnabledByDefault 声明该模块在 `modules:` 配置未显式设置 enabled 时的默认启停状态。
	// 零值 false 意味着新模块默认不启用，保证引入新模块不会在未配置时改变现有行为。
	EnabledByDefault bool
placeholder
```

要点：

- **零值即可用**：`ModuleInfo()` 必须可在零值实例上调用且无副作用；依赖在 `Provision` 中装配；
- **默认 disabled**：除非有明确决策，新模块的 `EnabledByDefault` 应为 `false`（零值），未配置时模块不实例化、无任何副作用；
- **注册在 `init()` 期完成**，重复 ID、nil New、非法 ID 一律 panic（注册错误属于编译期插装错误，必须在进程启动最早期暴露）：

```go
const ID plugin.ModuleID = "job.hello"

func init() {
	plugin.RegisterModule(&Module{placeholder)
placeholder
```

建议在模块包中加编译期断言，确保实现的生命周期接口与预期一致：

```go
var (
	_ plugin.Module      = (*Module)(nil)
	_ plugin.Provisioner = (*Module)(nil)
	_ plugin.Validator   = (*Module)(nil)
	_ plugin.Starter     = (*Module)(nil)
	_ plugin.Stopper     = (*Module)(nil)
)
```

## 3. 生命周期契约

除 `Module` 外，模块可按需实现四个**可选**生命周期接口：

```go
type Provisioner interface{ Provision(host *Host) error placeholder
type Validator   interface{ Validate() error placeholder
type Starter     interface{ Start(ctx context.Context) error placeholder
type Stopper     interface{ Stop(ctx context.Context) error placeholder
```

Runtime 按**注册表稳定序（ID 字典序）**驱动生命周期，宿主侧用法（`cmd/server/main.go` 已接好，模块作者无需关心）：

```go
rt := plugin.NewRuntime(host, cfg)
if err := rt.Build(); err != nil { ... placeholder   // 启动中止
if err := rt.Start(ctx); err != nil { ... placeholder // 已自动逆序回滚
defer rt.Stop(shutdownCtx)                  // 接入既有 cleanup 序列
```

### 3.1 各阶段语义

| 阶段 | 调用时机 | 失败后果 |
|---|---|---|
| `New()` | `Runtime.Build`，仅实例化 enabled 模块 | 返回 nil 实例 → Build 报错，启动中止 |
| `Provision(host)` | Build 阶段，实例化后。模块在此从 Host 获取所需能力（Logger / ConfigOf / DB / Redis）并完成自身初始化 | Build 报错，**启动中止**（fail-fast） |
| `Validate()` | Build 阶段，Provision 之后。校验模块配置与状态 | Build 报错，**启动中止** |
| `Start(ctx)` | `Runtime.Start`，HTTP server 启动前，按稳定序。启动后台工作（goroutine、监听等）；未实现 Starter 的 enabled 模块直接视为 running | 见 3.2 逆序回滚 |
| `Stop(ctx)` | `Runtime.Stop`，进程优雅关闭时，按**稳定序的逆序** | 见 3.3 |

`Build` 任一模块失败立即返回包含模块 ID 的错误，不会继续处理后续模块；已成功 Provision 的模块保持 provisioned 状态，由进程启动失败路径整体退出（`main` 中 `log.Fatalf`）。Runtime 为单次使用：Build/Start 各只能调用一次，不支持 Stop 后重新 Start。

### 3.2 Start 半途失败：逆序回滚

任一模块 `Start` 失败时，Runtime 对**已进入 running 的模块按逆序执行 Stop 回滚**（回滚沿用传入的 ctx；回滚失败仅记日志），然后返回包含失败模块 ID 的错误，进程启动中止。因此：

- `Start` 失败的模块应自行清理 `Start` 内部已分配的资源（Runtime 不会对失败模块本身调用 Stop）；
- `Stop` 必须容忍"Start 从未成功过"的状态（见 9.3 hello 的 no-op 处理）。

### 3.3 Stop 逆序与 ctx deadline

- `Stop` 按稳定序的**逆序**执行：后启动的模块先停止；
- 单个模块 Stop 失败**只记日志并继续**关闭其余模块，最终返回 `errors.Join` 聚合的失败集合（与既有 cleanup 序列的容错纪律一致）；
- ctx 原样传递给各模块 Stop，**模块必须自行尊重其 deadline**——宿主的 cleanup 整体超时为 10 秒，模块 Stop 与其他业务服务关闭并行执行，且先于 Redis/Ent 等基础设施关闭（模块停止时仍可写 DB/Redis）。

### 3.4 状态机

```
registered ──Build──▶ provisioned ──Start──▶ running ──Stop──▶ stopped
     │（disabled 模块停留在 registered）            任一阶段失败 ──▶ errored
```

- `registered`：已注册未实例化（disabled 模块始终处于该状态）；
- `provisioned`：已实例化并通过 Provision/Validate，等待启动；
- `running`：已启动（未实现 Starter 的 enabled 模块在 Start 阶段也进入该状态）；
- `stopped`：已正常停止；
- `errored`：在某个生命周期阶段失败，错误见 `ModuleStatus.Err`。

`Runtime.Snapshot()` 返回全部已注册模块（含 disabled）的状态快照，按 ID 字典序稳定排序，供可观测使用。

## 4. Host 能力面

`Provision` 注入的 `*plugin.Host` 是模块可使用的**全部**宿主能力，当前仅四项：

```go
type Host struct {
	// Logger 是宿主提供的结构化日志器。模块应通过它输出日志，而不是自建全局 logger。
	Logger *zap.Logger

	// ConfigOf 将 `modules:` 配置子树中 id 对应的 raw 私有配置解码到 out
	// （out 必须是指向模块自定义配置 struct 的指针）。
	ConfigOf func(id ModuleID, out any) error

	// DB 是 Ent ORM 客户端（数据库访问能力）。
	DB *ent.Client

	// Redis 是 go-redis 客户端（缓存/队列能力）。
	Redis redis.UniversalClient
placeholder
```

**铁律（host.go 边界纪律）：**

1. Host 只暴露 ports 级基础能力，**禁止暴露任何具体 Service**——模块依赖具体 Service 等于绕过插件边界，回到单体耦合；
2. **后续每新增一项 Host 能力，必须在[插件化 ROADMAP](../../.claude/plugin-refactor/ROADMAP.md) 决策记录中登记**，防止 Host 膨胀为上帝对象；
3. 模块内部依赖一律经 Host 获取，**不进入 Wire graph**（Wire 只负责装配 Host 与 Runtime，模块目录由注册表负责）。

防御性约定：`Provision` 应容忍 `host`、`host.Logger`、`host.ConfigOf` 为 nil（测试中常以部分填充的 Host 驱动模块），参考 hello 模块的处理（见第 9 节）。

## 5. 配置子树（`modules:`）

### 5.1 用法

模块配置位于全局配置（`config.yaml`）的 `modules:` 子树，**顶层 key 是完整的带点模块 ID 字面量**（一个 key，不是嵌套层级）：

```yaml
modules:
  job.hello:
    enabled: true       # 内核保留键：启停开关
    interval: 10s       # 以下为模块私有配置，内核不解释
    greeting: "hello from sub2api"
```

每个模块项中只有 `enabled` 键由内核解释（接受 bool 或可解析为 bool 的字符串），其余键原样保留为模块私有配置，由 `Config.Of` 解码到模块自定义的配置 struct。

### 5.2 enabled 三态语义

| 配置情况 | 最终启停 |
|---|---|
| `enabled: true` | 启用 |
| `enabled: false` | 禁用 |
| 未显式配置（含整个模块项缺失、`modules:` 子树不存在、或写了空的 `enabled:` 即 YAML null） | 由注册时声明的 `ModuleInfo.EnabledByDefault` 决定 |

缺省（`modules:` 子树不存在）时为空配置，所有模块按各自默认 enabled 值运行——对默认 disabled 的模块即"完全不存在"，行为与未引入插件内核前一致。

### 5.3 私有配置解码

`Config.Of`（经 `host.ConfigOf` 暴露给模块）的解码行为与项目全局配置（Viper 默认解码器）保持一致：

- `mapstructure` 标签、弱类型转换、字符串到 `time.Duration` / 逗号分隔切片的钩子；
- **模块未配置或私有配置为空时不修改 out**——模块应先填默认值再调用 `ConfigOf`（参考 hello 的 `defaultConfig()`）；
- 类型不匹配等解码失败时返回包含模块 ID 的错误，模块应将其从 `Provision` 原样返回（启动中止）。

### 5.4 未知 ID 忽略 / 非法 ID fail-fast

`ParseConfig` 的语义（在应用启动装配阶段执行，见 `plugin/wire.go` 的 `ProvideModuleConfig`）：

- **ID 格式非法**（如 `modules:` 下出现 `Job.Hello`、空段）或 **enabled 类型错误** → 返回错误，**应用启动 fail-fast**，配置笔误在启动时即暴露；
- **格式合法但未注册**的模块 ID（例如不同编译变体的配置共用）→ 允许存在，**被 Runtime 忽略**，不影响其他模块。

### 5.5 viper 含点 key 的注意事项

`modules:` 子树**不走 `viper.Unmarshal`**，而是在配置加载时通过 `viper.Get("modules")` 手工提取（见 `internal/config/config.go` 的 `Config.Modules` 字段与 `normalizeModulesSubtree`）——Unmarshal 基于 AllSettings 重建嵌套结构，会把含点的模块 ID key（`job.hello`）错误拆分为多级嵌套（`job: { hello: ... placeholder`）。模块作者无需做任何事，但需注意：

- 在 YAML 中模块 ID 必须写成**单个 key**（`job.hello:`），不要写成嵌套层级；
- 空模块项（如 `job.hello:` 后无内容）规整为空 map，语义是"提及但无私有配置"（enabled 仍走默认值）；
- 全局配置的环境变量覆盖（`.` → `_` 替换）不适用于 `modules:` 子树下含点的模块 ID key，模块配置请写在配置文件中。

## 6. 插装清单（Instrumentation List）

`backend/internal/modules/standard/imports.go` 是**唯一插装点**：所有需要编译进二进制的内置模块必须且只能在该文件中匿名 import。模块包在 `init()` 期向包级默认注册表自注册，`cmd/server/main.go` 匿名 import 本包即完成全部内置模块的插装。

```go
// Package standard 是插件模块的唯一插装清单（instrumentation list）。
package standard

import (
	// job.hello：示例模块（默认 disabled），验证插件链路连通性。
	_ "github.com/Wei-Shaw/sub2api/internal/modules/hello"
)
```

**新增一个模块 = 在此文件加一行 import**（带一行用途注释），使插装变更集中、可 review。不要在其他任何位置（包括测试外的业务代码）匿名 import 模块包。

## 7. 测试要求

### 7.1 用隔离 Registry，不污染包级默认注册表

涉及 Runtime 生命周期的测试**必须**在隔离注册表中驱动，避免依赖（或污染）包级默认注册表。默认入口是 `internal/plugin/plugintest` 夹具包——`RunLifecycle` / `BuildOnly` / `BuildExpectingError` 内部已用 `plugin.NewRegistry()` 构造隔离环境，一行驱动完整生命周期（API 速览与典型用法见 8.2）。

plugintest 未覆盖的复杂编排（Start 半途失败回滚、多模块顺序等）直接使用完全公开的内核 API：

```go
reg := plugin.NewRegistry()
reg.RegisterModule(&Module{placeholder)
rt := plugin.NewRuntimeWithRegistry(&plugin.Host{Logger: loggerplaceholder, cfg, reg)

require.NoError(t, rt.Build())
require.NoError(t, rt.Start(context.Background()))
require.NoError(t, rt.Stop(context.Background()))
```

测内核语义本身（注册 panic、命名空间查询、生命周期顺序等）时，用 fake 模块（实现部分生命周期接口的最小 struct）注册进隔离 Registry，参考 `internal/plugin/registry_test.go`、`runtime_test.go` 中的写法。

### 7.2 模块必测清单

参照 `internal/modules/hello/hello_test.go`，新模块至少覆盖：

1. **已注册且默认值正确**：`plugin.GetModule(ID)` 可查到，`EnabledByDefault` 符合预期；
2. **Provision 默认配置**：未提供私有配置时使用默认值并通过 Validate；nil host 不 panic；
3. **私有配置解码**：经 `host.ConfigOf` 解码（含 duration 字符串等弱类型转换）；解码类型错误时 Provision 返回含模块 ID 的错误；
4. **Validate 拒绝非法配置**（表驱动）；
5. **Start/Stop 行为**：启动后副作用可观测（如用 `zaptest/observer` 断言日志）、Stop 后 goroutine 优雅退出、未 Start 时 Stop 为 no-op；
6. **缺省配置零副作用**：`plugin.Config{placeholder` 下模块保持 disabled，Build/Start/Stop 全程 no-op、Snapshot 状态为 `registered`、无任何模块日志。

### 7.3 回归 gate（每 PR 硬性要求）

插件化改造全程零回归，模块相关 PR 必须通过 Phase-0 回归安全网（在 `backend/` 目录执行）：

```bash
make test-invariants            # 不变量/特征化测试（详见 .claude/plugin-refactor/INVARIANTS.md）
./scripts/bench-baseline.sh compare   # 热路径基准 vs 基线（allocs/op 严格、ns/op 容差 15%）
```

不变量清单见 [INVARIANTS.md](../../.claude/plugin-refactor/INVARIANTS.md)，基准脚本说明见 [bench-baseline.sh](../../backend/scripts/bench-baseline.sh) 文件头注释。**默认配置（不启用任何模块）下，全部不变量必须绿、基准无回归。** 两条命令的预期输出、阈值语义与 Wire 重生成注意事项见 8.5。

## 8. 开发与调试工作流

本章把"从零开发并调试一个模块"串成可照做的闭环，每一步都给出实测过的命令与预期输出（实测环境：2026-06 主线）。除特别说明外，命令均在 `backend/` 目录下执行。

```
make new-module（8.1 脚手架）→ plugintest 单测驱动开发（8.2）
→ imports.go 插装 + config.yaml 启用 + go run 本地运行（8.3）
→ 日志 / admin API / 前端页面观测（8.4）
→ make test-invariants + bench compare 提交前自检（8.5）
```

### 8.1 脚手架：生成模块骨架

```bash
make new-module ID=job.foo        # 等价：go run ./tools/newmodule -id job.foo
```

预期输出（实测，全文）：

```
module job.foo generated:
  internal/modules/foo/foo.go
  internal/modules/foo/foo_test.go

Next steps (the generator never edits imports.go -- explicit instrumentation keeps it reviewable):

  1. 编译期插装：在 internal/modules/standard/imports.go 的 import 块中加入：

	// job.foo：TODO 一句话描述模块用途（默认 disabled）。
	_ "github.com/Wei-Shaw/sub2api/internal/modules/foo"

  2. 配置启用（config.yaml；新模块默认 disabled，未显式 enabled 时零行为变更）：

	modules:
	  job.foo:
	    enabled: true
	    interval: 30s

  3. 填充代码中的 TODO 后运行测试：

	go test -count=1 ./internal/modules/foo/

生命周期契约与完整示例见 docs/plugin-architecture/MODULE-AUTHOR-GUIDE.md（蓝本实现：internal/modules/hello/）。
```

产物说明（模板从 hello 蓝本提炼，渲染时经 gofmt 校验，ID/包名非法直接报错、目标目录已存在拒绝覆盖）：

- `foo.go`：全部四个可选生命周期接口的骨架 + 编译期断言 + `EnabledByDefault: false` + 示例私有配置（`interval`）与校验，`TODO` 注释标出需要填充/删改的位置；
- `foo_test.go`：基于 plugintest 的 5 个用例（注册校验、默认配置、非法配置、缺省 no-op、启用生命周期），**开箱即过**：

```bash
$ go test -count=1 ./internal/modules/foo/
ok      github.com/Wei-Shaw/sub2api/internal/modules/foo        0.017s
```

生成器**不会**自动修改 `standard/imports.go`——插装行按 next steps 提示手工加入（设计原则与唯一插装点纪律见第 6 节）。

### 8.2 单测驱动开发：plugintest 夹具

`internal/plugin/plugintest` 是模块测试的默认入口（**仅供 `_test.go` import**，进入生产依赖面在评审中一票否决）。API 速览：

| 入口 | 用途 |
|---|---|
| `NewHost(t, ...Option)` | 一行构造 Host，默认零依赖可用：nop logger + 空 ConfigOf + nil DB/Redis |
| `WithConfig(raw)` | 把与 `modules:` 子树同形状的 raw 绑定为 ConfigOf，供直接调 `Provision` 的测试 |
| `WithObservedLogger()` | 返回 `(Option, *observer.ObservedLogs)`，debug 级别，用于断言日志 |
| `WithRedis(t)` | 附带 miniredis 支撑的 go-redis 客户端，生命周期挂 `t.Cleanup` 自动清理 |
| `RunLifecycle(t, m, host, raw)` | 隔离注册 → 解析 raw → Build → Start；`t.Cleanup` 自动 Stop（10s 上限）并断言无错 |
| `BuildOnly(...)` | 同 RunLifecycle 但停在 Build 之后，只验证 Provision/Validate |
| `BuildExpectingError(...)` | 期望 Build 失败，返回错误供断言内容（错误信息含模块 ID） |

典型用法（均摘自 `internal/modules/hello/hello_test.go`）。

直接调 Provision 验证私有配置解码——`NewHost` + `WithConfig`：

```go
host := plugintest.NewHost(t, plugintest.WithConfig(map[string]map[string]any{
	"job.hello": {"enabled": true, "interval": "15ms", "greeting": "hi"placeholder,
placeholder))

m := new(Module)
require.NoError(t, m.Provision(host))
require.Equal(t, 15*time.Millisecond, m.cfg.Interval)
require.Equal(t, "hi", m.cfg.Greeting)
```

完整生命周期 + 日志断言——`RunLifecycle` + `WithObservedLogger`：

```go
logOpt, logs := plugintest.WithObservedLogger()
rt := plugintest.RunLifecycle(t, &Module{placeholder, plugintest.NewHost(t, logOpt), map[string]map[string]any{
	"job.hello": {"enabled": true, "interval": "5ms"placeholder,
placeholder)

require.Eventually(t, func() bool {
	return logs.FilterMessage(defaultConfig().Greeting).Len() >= 1
placeholder, 2*time.Second, 5*time.Millisecond, "periodic debug log should appear after Start")

require.NoError(t, rt.Stop(context.Background()))
snap := rt.Snapshot()
require.True(t, snap[0].Enabled)
require.Equal(t, plugin.StateStopped, snap[0].State)
```

非法配置应使 Build 失败（fail-fast）——`BuildExpectingError`：

```go
err := plugintest.BuildExpectingError(t, &Module{placeholder, plugintest.NewHost(t), map[string]map[string]any{
	"job.hello": {"enabled": true, "interval": "-1s"placeholder,
placeholder)
require.Contains(t, err.Error(), `validate module "job.hello"`)
require.Contains(t, err.Error(), "interval must be positive")
```

语义提醒：

- `RunLifecycle` / `BuildOnly` 的 `raw` 是模块配置的**唯一来源**（会覆盖 host 上 `WithConfig` 设置的 ConfigOf），两者不要混用；
- raw 未写 `enabled: true` 且模块默认 disabled 时，Build/Start 全程 no-op，夹具会 `t.Logf` 提示（防"忘写 enabled"空跑：`module "job.foo" resolved to disabled under the given config ...`）；
- Host 是纯数据结构，夹具未覆盖的少见场景（如注入 enttest 客户端）直接对字段赋值即可。

模块必测清单见 7.2；plugintest 未覆盖的复杂编排用内核 API（见 7.1）。

### 8.3 本地运行：在真实进程中启用模块

**依赖准备（PostgreSQL 15+ / Redis 7+）。** 调试模块推荐在宿主机直跑后端：

- 本机已有 PG/Redis：`config.yaml` 的 `database:` / `redis:` 指向 `127.0.0.1` 即可（本地环境约定参照仓库根 `DEV_GUIDE.md`）；
- 用 Docker 起依赖：`deploy/docker-compose.dev.yml` 内有 `postgres` / `redis` 服务，但**默认不映射宿主端口**（仅供同 compose 网络内的 sub2api 容器使用）。整栈在容器内运行（改模块代码需重新 build）：

  ```bash
  cd deploy && POSTGRES_PASSWORD=devpass docker compose -f docker-compose.dev.yml up -d --build
  ```

  若坚持宿主跑后端 + Docker 跑依赖，需自行给 `postgres` / `redis` 服务加 `ports:` 映射（`5432:5432` / `6379:6379`）后再 `docker compose -f docker-compose.dev.yml up -d postgres redis`。

**配置。** 后端从工作目录查找 `config.yaml`（还会查 `./config`、`/etc/sub2api` 等；完整模板见 `deploy/config.example.yaml`），在 `backend/` 下运行即读取 `backend/config.yaml`。在其中追加 `modules:` 子树（顶层任意位置，语义见第 5 节）：

```yaml
modules:
  job.hello:
    enabled: true
    interval: 10s
```

**启动与日志。**

```bash
go run ./cmd/server
```

启动序列中、HTTP server 监听之前，可看到模块生命周期日志（实测摘录；`interval` 字段为毫秒数）：

```
2026-06-11T20:01:31.765+0800  INFO  hello/hello.go:111     hello module started   {"service": "sub2api", "env": "production", "module": "job.hello", "interval": 10000placeholder
2026-06-11T20:01:31.765+0800  INFO  plugin/runtime.go:209  plugin module started  {"service": "sub2api", "env": "production", "module": "job.hello"placeholder
```

随后出现 `Server started on <host>:<port>`。注意 `plugin module provisioned` 与 hello 的周期 greeting 是 **Debug** 级别，默认 `log.level: info` 下不可见——要观察周期行为，在 `config.yaml` 中设 `log: { level: debug placeholder`。模块配置非法时启动 fail-fast：进程以 `Failed to build plugin modules: ...`（错误含模块 ID）退出，符合 5.4 的语义。

**优雅关闭。** Ctrl+C（SIGINT/SIGTERM）后模块按逆序停止，与既有 cleanup 步骤交织输出（实测摘录）：

```
2026-06-11T20:03:50.449+0800  INFO  hello/hello.go:149     hello module stopped   {"service": "sub2api", "env": "production", "module": "job.hello"placeholder
2026-06-11T20:03:50.449+0800  INFO  plugin/runtime.go:252  plugin module stopped  {"service": "sub2api", "env": "production", "module": "job.hello"placeholder
2026-06-11T20:03:50.449+0800  INFO  stdlog  runtime/asm_amd64.s:1771  [Cleanup] PluginModuleRuntime succeeded  {"service": "sub2api", "env": "production", "legacy_stdlog": trueplaceholder
```

### 8.4 三种观测方式

**1) 结构化日志。** 模块经 `Host.Logger` 输出的日志自带全局字段（`service` / `env`），按模块自身附加的 `module` 字段过滤即可。Runtime 侧固定输出：`plugin module started` / `plugin module stopped`（Info）、`plugin module provisioned`（Debug）、`plugin module stop failed`（Error），均含 `"module": "<ID>"` 字段。

**2) admin API（只读）。** 支持管理员 JWT 或系统设置中配置的 Admin API Key 两种鉴权（端口以 `server.port` 为准）：

```bash
curl -s http://127.0.0.1:8080/api/v1/admin/modules -H "x-api-key: <admin-api-key>"
# 或：curl -s http://127.0.0.1:8080/api/v1/admin/modules -H "Authorization: Bearer <admin-jwt>"
```

实测响应（job.hello 启用且运行中）：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "modules": [
      {
        "id": "job.hello",
        "namespace": "job",
        "name": "hello",
        "enabled": true,
        "state": "running",
        "error": ""
      placeholder
    ]
  placeholder
placeholder
```

要点：返回**全部已注册模块**（disabled 模块也在列，`state` 为 `registered`），按 ID 字典序；`error` 文本已统一脱敏；全新部署需先完成管理员合规确认，否则该接口返回 423 `ADMIN_COMPLIANCE_ACK_REQUIRED`（用前端登录一次按引导确认即可）。

**3) 前端「插件模块」页。**

```bash
cd frontend && pnpm install && pnpm run dev    # http://127.0.0.1:3000
```

dev server 默认把 `/api` 代理到 `http://localhost:8080`；后端端口不是 8080 时在 `frontend/.env` 设 `VITE_DEV_PROXY_TARGET=http://127.0.0.1:<port>`。管理员登录后，经侧边栏「插件模块」进入 `/admin/modules`：只读列表展示模块 ID / 命名空间 / 名称 / 启用 / 状态 / 错误信息，状态以语义徽章着色（running=绿、errored=红、registered=灰、provisioned/stopped=主色），错误长文本截断 + 悬浮查看完整内容，右上角手动刷新。页面消费的就是上面的 admin API（实现：`frontend/src/views/admin/ModulesView.vue`）。

### 8.5 提交前自检

```bash
make test-invariants                  # 回归不变量/特征化测试（实测约 16s）
./scripts/bench-baseline.sh compare   # 热路径基准 vs 基线（BENCH_COUNT=6 实测约 4 分钟）
```

`test-invariants` 只筛选 `Characterization|Invariant` 用例，多数包显示 `[no tests to run]` 属正常，关注每行是否全 `ok`、退出码为 0。

`compare` 先跑当前基准（有 `benchstat` 时附其报告），再对每个基准输出对比行并按阈值判定（实测样例）：

```
BenchmarkGatewayForward_AnthropicNonStreamPassthrough  ns/op  10164 -> 9882 (-2.8%)  allocs/op  93.0 -> 93.0  ok
```

阈值语义（详见脚本头注释）：

- **allocs/op 严格**：均值增加超过 `ALLOC_TOLERANCE`（默认 0）即 FAIL；
- **ns/op 宽松**：劣化超过 `TIME_TOLERANCE_PCT`%（默认 15，吸收 runner 抖动）即 FAIL；
- **基线中的基准在本次运行缺失**（被删除/改名）→ FAIL（回归网已失效，须修复或重新 collect）；新增基准只 warn、不参与对比；
- 基线文件不存在 → 退出码 2，提示先 `./scripts/bench-baseline.sh collect`（换机器/换 Go 版本后也要重新采集）。

**Wire 重生成注意。** 普通模块不进 Wire graph，无需此步；只有改了 Wire 装配（如 `internal/plugin/wire.go`、`cmd/server/wire.go`）才需要重生成。当前 `make generate` / `go generate ./cmd/server` 中无版本号的 wire 调用因 go.sum 缺 `github.com/google/subcommands` 条目而失败（预存在问题，实测报错 `missing go.sum entry for module providing package github.com/google/subcommands`），实测可用的调用方式：

```bash
go run github.com/google/wire/cmd/wire@v0.7.0 ./cmd/server/
```

**前端有改动时**，在仓库根执行 `make test-frontend`（lint + typecheck + 关键 vitest 集）。

### 8.6 常见坑

| 坑 | 现象 | 处置 |
|---|---|---|
| 忘加 imports.go | 单包 `go test` 全绿（同包测试必然触发 `init()` 注册），但真实进程里模块"不存在"：无生命周期日志，Snapshot / admin API / 前端页面的列表里**根本没有这个模块**（配置中格式合法但未注册的 ID 被 Runtime 静默忽略，见 5.4） | 在 `standard/imports.go` 加匿名 import（见第 6 节） |
| `enabled:` 留空 | YAML null = 未显式配置 = 走 `EnabledByDefault`（默认 false），模块不会运行 | 启用必须写 `enabled: true`（三态语义见 5.2） |
| 含点 ID 的 viper 行为 | 模块 ID 在 YAML 里写成嵌套层级（`job:` 下挂 `hello:`）不被识别；环境变量覆盖对 `modules:` 子树不生效 | ID 写成单个 key、模块配置写在配置文件中（见 5.5）。改动 modules 配置解析本身时，注意该子树不走 `viper.Unmarshal`，而是 `internal/config` 的 `normalizeModulesSubtree` 手工提取 |
| `EnabledByDefault: true` | 模块未配置即运行——"引入新模块在缺省配置下零行为变更"的铁律被打破 | 保持 false；确需默认启用须先在 ROADMAP 登记决策（见 2.3） |
| Stop 不尊重 ctx deadline | 宿主 cleanup 整体超时 10 秒、模块 Stop 与其他服务关闭并行，悬挂的 Stop 拖死优雅关闭；plugintest 兜底 Stop 同样 10s 上限，悬挂表现为测试失败 | 等待 worker 退出必须同时 `select` 上 `ctx.Done()`（照抄 9.3 的 Stop 写法） |
| Start 里跑长任务 | Runtime 在持锁状态下调用 Start/Stop，阻塞的 Start 会卡住 Snapshot（admin API 与前端页面跟着挂起）和后续模块的启动 | Start 必须快速返回，长任务自起 goroutine；goroutine 不要挂在 Start 的 ctx 上，寿命由 Stop 终结（见 9.3） |
| mapstructure 标签含大写 | 全局配置经 Viper 加载时 key 被静默小写化，含大写的标签永远解码不到值 | `mapstructure` 标签全小写（脚手架模板注释已注明） |
| 前端单文件测试跑成全量 | `pnpm run test:run -- <file>` 的 `--` 透传会使 vitest 文件过滤失效 | 用 `pnpm exec vitest run <file>`（如 `pnpm exec vitest run src/views/admin/__tests__/ModulesView.spec.ts`） |

## 9. 完整示例：job.hello

`backend/internal/modules/hello/hello.go` 是最小可运行的参考实现：实现全部四个可选生命周期接口，启动后按配置间隔周期性输出一条 debug 日志，默认 disabled。以下代码均摘自真实文件（有裁剪）。

### 9.1 注册

```go
// ID 是 hello 模块的模块 ID。
const ID plugin.ModuleID = "job.hello"

func init() {
	plugin.RegisterModule(&Module{placeholder)
placeholder

// Module 是 job.hello 模块实例。零值即可用，依赖在 Provision 中装配；
// 每次 Runtime.Build 都会通过 ModuleInfo.New 创建全新实例。
type Module struct {
	cfg    Config
	logger *zap.Logger

	cancel context.CancelFunc
	done   chan struct{placeholder
placeholder

// ModuleInfo 返回模块注册信息。EnabledByDefault 为 false：
// 未在 `modules:` 中显式 enabled 时模块不实例化、无任何副作用。
func (*Module) ModuleInfo() plugin.ModuleInfo {
	return plugin.ModuleInfo{
		ID:               ID,
		New:              func() plugin.Module { return new(Module) placeholder,
		EnabledByDefault: false,
placeholder
placeholder
```

插装：`internal/modules/standard/imports.go` 中已有 `_ "github.com/Wei-Shaw/sub2api/internal/modules/hello"` 一行。

### 9.2 配置（默认值 + 解码 + 校验）

```go
// Config 是 hello 模块的私有配置
// （`modules:` 子树中 "job.hello" 项除 enabled 外的部分）。
type Config struct {
	// Interval 是周期日志的输出间隔，必须大于 0。
	Interval time.Duration `mapstructure:"interval"`

	// Greeting 是周期日志输出的内容，不能为空。
	Greeting string `mapstructure:"greeting"`
placeholder

// defaultConfig 返回 hello 模块的默认配置（用户未配置的字段保持默认值）。
func defaultConfig() Config {
	return Config{
		Interval: 30 * time.Second,
		Greeting: "hello from job.hello",
placeholder
placeholder

// Provision 从宿主获取 logger，并在默认配置之上解码模块私有配置。
func (m *Module) Provision(host *plugin.Host) error {
	m.cfg = defaultConfig()
	if host != nil {
		m.logger = host.Logger
		if host.ConfigOf != nil {
			if err := host.ConfigOf(ID, &m.cfg); err != nil {
				return err
		placeholder
	placeholder
placeholder
	if m.logger == nil {
		m.logger = zap.NewNop()
placeholder
	return nil
placeholder

// Validate 校验模块配置：interval 必须大于 0，greeting 不能为空。
func (m *Module) Validate() error {
	if m.cfg.Interval <= 0 {
		return fmt.Errorf("hello: interval must be positive, got %s", m.cfg.Interval)
placeholder
	if m.cfg.Greeting == "" {
		return errors.New("hello: greeting must not be empty")
placeholder
	return nil
placeholder
```

注意先填默认值再解码：`ConfigOf` 对未配置的模块不修改 out，默认值由模块自带。

### 9.3 启动与停止

```go
// Start 启动周期日志 goroutine。
//
// 模块持有独立于启动 ctx 的生命周期上下文，由 Stop 负责取消，
// 因此模块的运行寿命不受启动流程上下文的影响。
func (m *Module) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	m.done = make(chan struct{placeholder)
	go m.run(ctx)
	m.logger.Info("hello module started",
		zap.String("module", string(ID)),
		zap.Duration("interval", m.cfg.Interval))
	return nil
placeholder

// Stop 取消周期 goroutine 并等待其退出，尊重 ctx 的 deadline。
// 未 Start 过（或 Start 之前失败）时调用为 no-op。
func (m *Module) Stop(ctx context.Context) error {
	if m.cancel == nil {
		return nil
placeholder
	m.cancel()
	select {
	case <-m.done:
		m.logger.Info("hello module stopped", zap.String("module", string(ID)))
		return nil
	case <-ctx.Done():
		return fmt.Errorf("hello: stop interrupted while waiting for worker exit: %w", ctx.Err())
placeholder
placeholder
```

两个值得照抄的细节：

- **后台 goroutine 不挂在 Start 的 ctx 上**——Start 的 ctx 只是启动流程上下文，模块寿命由 Stop 显式终结；
- **Stop 等待 worker 真正退出，且尊重 ctx deadline**——超时返回包装了 `ctx.Err()` 的错误，Runtime 会记日志并继续关闭其余模块。

### 9.4 启用与观测

启用（`config.yaml`）：

```yaml
modules:
  job.hello:
    enabled: true
    interval: 10s
```

启动/关闭日志样例、`GET /api/v1/admin/modules`（实现见 `internal/handler/admin/module_handler.go`）与前端「插件模块」页等观测方式，统一见第 8 章（8.3 本地运行、8.4 三种观测方式）；程序内观测用 `Runtime.Snapshot()`（语义见 3.4）。

### 9.5 新模块检查清单

- [ ] 模块包位于 `internal/modules/<模块名>/`，ID 落入 ROADMAP 已规划命名空间；
- [ ] `init()` 中 `plugin.RegisterModule`，零值可用、`ModuleInfo()` 无副作用；
- [ ] `EnabledByDefault: false`（除非有登记过的决策）；
- [ ] 依赖只经 Host 获取，不依赖任何具体 Service、不进 Wire graph；
- [ ] 私有配置：默认值自带 + `ConfigOf` 解码 + `Validate` 校验；
- [ ] Stop 容忍未 Start、尊重 ctx deadline；
- [ ] `standard/imports.go` 加一行匿名 import；
- [ ] 测试覆盖第 7.2 节清单，使用隔离 Registry（默认经 plugintest，见 7.1 / 8.2）；
- [ ] `make test-invariants` 与 `./scripts/bench-baseline.sh compare` 全绿（运行细节见 8.5）。

---

后续阶段的迁移计划与远期能力（试点模块、平台 Provider 模块化、模块自带前端页面、进程外桥等）以 [插件化 ROADMAP](../../.claude/plugin-refactor/ROADMAP.md) 为准，本指南只描述已落地的内核行为。
