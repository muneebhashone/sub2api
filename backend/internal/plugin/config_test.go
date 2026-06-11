package plugin

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseConfigEnabledTriState(t *testing.T) {
	cfg, err := ParseConfig(map[string]map[string]any{
		"job.on":      {"enabled": trueplaceholder,
		"job.off":     {"enabled": falseplaceholder,
		"job.unset":   {"greeting": "hi"placeholder,
		"job.strbool": {"enabled": "true"placeholder, // 弱类型：字符串 bool（如来自环境变量）
placeholder)
placeholder

	require.NotNil(t, cfg["job.on"].Enabled)
	require.True(t, *cfg["job.on"].Enabled)

	require.NotNil(t, cfg["job.off"].Enabled)
	require.False(t, *cfg["job.off"].Enabled)

	require.Nil(t, cfg["job.unset"].Enabled, "未显式配置 enabled 时必须为 nil（用模块默认值）")

	require.NotNil(t, cfg["job.strbool"].Enabled)
	require.True(t, *cfg["job.strbool"].Enabled)
placeholder

func TestParseConfigSeparatesEnabledFromRaw(t *testing.T) {
	cfg, err := ParseConfig(map[string]map[string]any{
		"job.hello": {"enabled": true, "greeting": "hi", "interval": "5s"placeholder,
placeholder)
placeholder

	mc := cfg["job.hello"]
	require.Equal(t, map[string]any{"greeting": "hi", "interval": "5s"placeholder, mc.Raw,
		"enabled 键属于内核，不得混入模块私有配置")
placeholder

func TestParseConfigErrors(t *testing.T) {
	t.Run("invalid module ID", func(t *testing.T) {
		_, err := ParseConfig(map[string]map[string]any{"Job.Hello": {placeholderplaceholder)
	placeholder
		require.Contains(t, err.Error(), "Job.Hello")
placeholder)
	t.Run("invalid enabled type", func(t *testing.T) {
		_, err := ParseConfig(map[string]map[string]any{"job.hello": {"enabled": placeholderplaceholder)
	placeholder
		require.Contains(t, err.Error(), "job.hello")
placeholder)
	t.Run("invalid enabled string", func(t *testing.T) {
		_, err := ParseConfig(map[string]map[string]any{"job.hello": {"enabled": "not-a-bool"placeholderplaceholder)
	placeholder
placeholder)
	t.Run("unregistered but valid ID is allowed", func(t *testing.T) {
		cfg, err := ParseConfig(map[string]map[string]any{"future.module": {"enabled": trueplaceholderplaceholder)
	placeholder
		require.Contains(t, cfg, ModuleID("future.module"))
placeholder)
placeholder

type helloModuleConfig struct {
	Greeting string        `mapstructure:"greeting"`
	Interval time.Duration `mapstructure:"interval"`
	Targets  []string      `mapstructure:"targets"`
	Limit    int           `mapstructure:"limit"`
placeholder

func TestConfigOfDecodesRawIntoStruct(t *testing.T) {
	cfg, err := ParseConfig(map[string]map[string]any{
		"job.hello": {
			"enabled":  true,
			"greeting": "hi",
			"interval": "5s",    // 字符串 → time.Duration（与 viper 默认钩子一致）
			"targets":  "a,b,c", // 逗号分隔 → 切片（与 viper 默认钩子一致）
			"limit":    "42",    // 弱类型：字符串 → int
	placeholder,
placeholder)
placeholder

	var out helloModuleConfig
	require.NoError(t, cfg.Of("job.hello", &out))
	require.Equal(t, "hi", out.Greeting)
	require.Equal(t, 5*time.Second, out.Interval)
	require.Equal(t, []string{"a", "b", "c"placeholder, out.Targets)
	require.Equal(t, 42, out.Limit)
placeholder

func TestConfigOfMissingModuleLeavesOutUntouched(t *testing.T) {
	cfg := Config{placeholder
	out := helloModuleConfig{Greeting: "default-value"placeholder
	require.NoError(t, cfg.Of("job.absent", &out))
	require.Equal(t, "default-value", out.Greeting, "模块未配置时不得修改 out（保留模块默认值）")
placeholder

func TestConfigOfTypeErrorContainsModuleID(t *testing.T) {
	cfg, err := ParseConfig(map[string]map[string]any{
		"job.hello": {"limit": "not-an-int"placeholder,
placeholder)
placeholder

	var out helloModuleConfig
	err = cfg.Of("job.hello", &out)
placeholder
	require.Contains(t, err.Error(), "job.hello", "解码错误必须包含模块 ID")
placeholder

func TestConfigEnabledFor(t *testing.T) {
	on, off := true, false
	cfg := Config{
		"job.explicit-on":  {Enabled: &onplaceholder,
		"job.explicit-off": {Enabled: &offplaceholder,
		"job.mentioned":    {placeholder,
placeholder

	tests := []struct {
		name string
		info ModuleInfo
		want bool
placeholder{
		{"explicit on overrides default off", ModuleInfo{ID: "job.explicit-on", EnabledByDefault: falseplaceholder, trueplaceholder,
		{"explicit off overrides default on", ModuleInfo{ID: "job.explicit-off", EnabledByDefault: trueplaceholder, falseplaceholder,
		{"mentioned without enabled uses default true", ModuleInfo{ID: "job.mentioned", EnabledByDefault: trueplaceholder, trueplaceholder,
		{"mentioned without enabled uses default false", ModuleInfo{ID: "job.mentioned", EnabledByDefault: falseplaceholder, falseplaceholder,
		{"absent uses default true", ModuleInfo{ID: "job.absent", EnabledByDefault: trueplaceholder, trueplaceholder,
		{"absent uses default false", ModuleInfo{ID: "job.absent", EnabledByDefault: falseplaceholder, falseplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, cfg.enabledFor(tt.info))
	placeholder)
placeholder
placeholder
