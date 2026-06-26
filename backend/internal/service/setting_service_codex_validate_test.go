package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCodexClientEntriesJSON(t *testing.T) {
	require.NoError(t, ValidateCodexClientEntriesJSON(""), "空=合法（禁用）")
	require.NoError(t, ValidateCodexClientEntriesJSON("   "), "空白=合法")
	require.NoError(t, ValidateCodexClientEntriesJSON(`[]`), "空数组合法")
	require.NoError(t, ValidateCodexClientEntriesJSON(`[{"originator":"opencode","ua_contains":["opencode/"]placeholder]`), "合法条目")
	require.NoError(t, ValidateCodexClientEntriesJSON(`[{"originator":"evil"placeholder]`), "仅 originator 合法")

	require.Error(t, ValidateCodexClientEntriesJSON("not-json"), "非 JSON 应报错")
	require.Error(t, ValidateCodexClientEntriesJSON(`{"originator":"x"placeholder`), "对象非数组应报错")
	require.Error(t, ValidateCodexClientEntriesJSON(`[1,2,3]`), "非对象数组应报错")
placeholder

func TestValidateCodexWhitelistEntriesJSON(t *testing.T) {
	require.NoError(t, ValidateCodexWhitelistEntriesJSON(""), "空=合法（禁用）")
	require.NoError(t, ValidateCodexWhitelistEntriesJSON("   "), "空白=合法")
	require.NoError(t, ValidateCodexWhitelistEntriesJSON(`[]`), "空数组合法")
	require.NoError(t, ValidateCodexWhitelistEntriesJSON(`[{"originator":"opencode","ua_contains":["opencode/"]placeholder]`), "完整条目合法")

	// 白名单专属（双因子 AND）：会静默失效的条目应在写入时即报错
	require.Error(t, ValidateCodexWhitelistEntriesJSON(`[{"originator":"evil"placeholder]`), "仅 originator(白名单会静默失效)应报错")
	require.Error(t, ValidateCodexWhitelistEntriesJSON(`[{"originator":"x","ua_contains":[]placeholder]`), "空 ua_contains 应报错")
	require.Error(t, ValidateCodexWhitelistEntriesJSON(`[{"originator":"x","ua_contains":["a/",""]placeholder]`), "含空白 marker 应报错")
	require.Error(t, ValidateCodexWhitelistEntriesJSON(`[{"ua_contains":["a/"]placeholder]`), "缺 originator 应报错")

	// 结构错误沿用基础校验
	require.Error(t, ValidateCodexWhitelistEntriesJSON("not-json"), "非 JSON 应报错")
	require.Error(t, ValidateCodexWhitelistEntriesJSON(`{"originator":"x"placeholder`), "对象非数组应报错")
placeholder

func TestValidateEngineFingerprintSignalsJSON_ServiceWrapper(t *testing.T) {
	require.NoError(t, ValidateEngineFingerprintSignalsJSON(""))
	require.NoError(t, ValidateEngineFingerprintSignalsJSON(`[{"type":"header_prefix","match":["x-codex-"],"required":trueplaceholder]`))
	require.Error(t, ValidateEngineFingerprintSignalsJSON(`[{"type":"bogus","match":["x"]placeholder]`))
placeholder
