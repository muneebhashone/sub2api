package service

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestBuildDynamicToolMap_BelowThreshold(t *testing.T) {
	// Parrot 行为：tools 数量 ≤ 5 时不做动态映射。
	names := []string{"bash", "edit", "read", "write", "search"placeholder
	require.Nil(t, buildDynamicToolMap(names))
placeholder

func TestBuildDynamicToolMap_AboveThresholdIsStable(t *testing.T) {
	// Parrot 不变量：同一组 tool_names 在同进程内映射稳定（保证 cache 命中）。
	names := []string{"alpha", "beta", "gamma", "delta", "epsilon", "zeta"placeholder
	a := buildDynamicToolMap(names)
	b := buildDynamicToolMap(names)
	require.NotNil(t, a)
	require.Equal(t, a, b, "same input tool_names must yield identical mapping")
	require.Len(t, a, 6)
	for _, name := range names {
		require.Contains(t, a, name)
		require.NotEqual(t, name, a[name])
placeholder
placeholder

func TestSanitizeToolName_StaticPrefix(t *testing.T) {
	require.Equal(t, "cc_sess_list", sanitizeToolName("sessions_list", nil))
	require.Equal(t, "cc_ses_get", sanitizeToolName("session_get", nil))
	require.Equal(t, "bash", sanitizeToolName("bash", nil))
placeholder

func TestSanitizeToolName_DynamicTakesPrecedence(t *testing.T) {
	dyn := map[string]string{"sessions_list": "analyze_ses00"placeholder
	got := sanitizeToolName("sessions_list", dyn)
	require.Equal(t, "analyze_ses00", got, "dynamic mapping wins over static prefix")
placeholder

func TestRestoreToolNamesInBytes_LongestFirst(t *testing.T) {
	// 当假名 "abc_12" 是另一个更长假名的子串（真实场景极少但算法必须防御）时，
	// 长的必须先替换。本测试用显式构造的映射来验证排序不变量。
	rw := &ToolNameRewrite{
		Forward: map[string]string{"foo": "abc_12", "bar": "abc_12_ext"placeholder,
		Reverse: map[string]string{"abc_12": "foo", "abc_12_ext": "bar"placeholder,
placeholder
	// 手工构造 ReverseOrdered：长的在前
	rw.ReverseOrdered = [][2]string{
		{"abc_12_ext", "bar"placeholder,
		{"abc_12", "foo"placeholder,
placeholder
	data := []byte(`{"tool":"abc_12_ext","other":"abc_12"placeholder`)
	restored := string(restoreToolNamesInBytes(data, rw))
	require.Equal(t, `{"tool":"bar","other":"foo"placeholder`, restored)
placeholder

func TestRestoreToolNamesInBytes_StaticPrefixRollback(t *testing.T) {
	data := []byte(`{"name":"sessions_list","id":"cc_ses_xyz"placeholder`)
	got := string(restoreToolNamesInBytes(data, nil))
	require.Equal(t, `{"name":"sessions_list","id":"session_xyz"placeholder`, got)
placeholder

func TestApplyToolNameRewriteToBody_RenamesToolsAndToolChoice(t *testing.T) {
	body := []byte(`{"tools":[{"name":"sessions_list","input_schema":{placeholderplaceholder,{"name":"session_get","input_schema":{placeholderplaceholder,{"name":"web_search","type":"web_search_20250305"placeholder],"tool_choice":{"type":"tool","name":"sessions_list"placeholderplaceholder`)
	rw := buildToolNameRewriteFromBody(body)
	require.NotNil(t, rw)
	require.Contains(t, rw.Forward, "sessions_list")
	require.Contains(t, rw.Forward, "session_get")
	// web_search 是 server tool，不参与工具名改写
	require.NotContains(t, rw.Forward, "web_search")

	out := applyToolNameRewriteToBody(body, rw)

	// tools[0].name 和 tools[1].name 被改写，tools[2].name 保持不变
	require.Equal(t, "cc_sess_list", gjson.GetBytes(out, "tools.0.name").String())
	require.Equal(t, "cc_ses_get", gjson.GetBytes(out, "tools.1.name").String())
	require.Equal(t, "web_search", gjson.GetBytes(out, "tools.2.name").String())

	// tool_choice.name 被同步改写
	require.Equal(t, "cc_sess_list", gjson.GetBytes(out, "tool_choice.name").String())
	require.Equal(t, "tool", gjson.GetBytes(out, "tool_choice.type").String())
placeholder

func TestApplyToolNameRewriteToBody_RenamesToolUseInMessages(t *testing.T) {
	// sessions_list 通过静态前缀规则改写为 cc_sess_list
	// web_search 是 server tool（type != ""），不参与工具名改写
	// messages 中的 tool_use.name 必须同步改写，才能和 tools[] 保持一致
	body := []byte(`{"tools":[{"name":"sessions_list","input_schema":{placeholderplaceholder,{"name":"web_search","type":"web_search_20250305"placeholder],"messages":[{"role":"user","content":[{"type":"text","text":"hi"placeholder]placeholder,{"role":"assistant","content":[{"type":"tool_use","id":"tu_01","name":"sessions_list","input":{placeholderplaceholder,{"type":"text","text":"thinking"placeholder]placeholder,{"role":"user","content":[{"type":"tool_result","tool_use_id":"tu_01","content":"ok"placeholder]placeholder]placeholder`)
	rw := buildToolNameRewriteFromBody(body)
	require.NotNil(t, rw)
	require.Equal(t, "cc_sess_list", rw.Forward["sessions_list"])

	out := applyToolNameRewriteToBody(body, rw)

	// tools[0].name 被改写
	require.Equal(t, "cc_sess_list", gjson.GetBytes(out, "tools.0.name").String())
	// tools[1].name 是 server tool，保持不变
	require.Equal(t, "web_search", gjson.GetBytes(out, "tools.1.name").String())
	// messages[1].content[0].name 是 tool_use，必须同步改写以匹配 tools[]
	require.Equal(t, "cc_sess_list", gjson.GetBytes(out, "messages.1.content.0.name").String())
	// messages[1].content[1] 是 text，保持不变
	require.Equal(t, "thinking", gjson.GetBytes(out, "messages.1.content.1.text").String())
	// messages[2].content[0] 是 tool_result，不包含 name 字段，保持不变
	require.Equal(t, "ok", gjson.GetBytes(out, "messages.2.content.0.content").String())
placeholder

func TestApplyToolNameRewriteToBody_RenamesToolUseWithDynamicMapping(t *testing.T) {
	body := []byte(`{"tools":[{"name":"alpha_search","input_schema":{placeholderplaceholder,{"name":"beta_lookup","input_schema":{placeholderplaceholder,{"name":"gamma_fetch","input_schema":{placeholderplaceholder,{"name":"delta_update","input_schema":{placeholderplaceholder,{"name":"epsilon_parse","input_schema":{placeholderplaceholder,{"name":"zeta_render","input_schema":{placeholderplaceholder,{"name":"web_search","type":"web_search_20250305"placeholder],"tool_choice":{"type":"tool","name":"gamma_fetch"placeholder,"messages":[{"role":"assistant","content":[{"type":"tool_use","id":"tu_dyn","name":"gamma_fetch","input":{placeholderplaceholder,{"type":"tool_use","id":"tu_srv","name":"web_search","input":{placeholderplaceholder,{"type":"text","text":"done"placeholder]placeholder,{"role":"user","content":[{"type":"tool_result","tool_use_id":"tu_dyn","content":"ok"placeholder]placeholder]placeholder`)
	rw := buildToolNameRewriteFromBody(body)
	require.NotNil(t, rw)
	require.Len(t, rw.Forward, 6)

	fakeGamma := rw.Forward["gamma_fetch"]
	require.NotEmpty(t, fakeGamma)
	require.NotEqual(t, "gamma_fetch", fakeGamma)
	require.NotContains(t, rw.Forward, "web_search")

	out := applyToolNameRewriteToBody(body, rw)

	// 动态映射会改写 tools[]、tool_choice 和历史 tool_use 中的同一个工具名
	require.Equal(t, fakeGamma, gjson.GetBytes(out, "tools.2.name").String())
	require.Equal(t, fakeGamma, gjson.GetBytes(out, "tool_choice.name").String())
	require.Equal(t, fakeGamma, gjson.GetBytes(out, "messages.0.content.0.name").String())
	// server tool 不参与动态映射，历史 tool_use 中同名引用也保持不变
	require.Equal(t, "web_search", gjson.GetBytes(out, "tools.6.name").String())
	require.Equal(t, "web_search", gjson.GetBytes(out, "messages.0.content.1.name").String())
	// tool_result 依靠 tool_use_id 关联，不需要 name 字段
	require.Equal(t, "ok", gjson.GetBytes(out, "messages.1.content.0.content").String())
placeholder

func TestApplyToolsLastCacheBreakpoint_InjectsDefault(t *testing.T) {
	body := []byte(`{"tools":[{"name":"a","input_schema":{placeholderplaceholder,{"name":"b","input_schema":{placeholderplaceholder]placeholder`)
	out := applyToolsLastCacheBreakpoint(body)
	require.Equal(t, "ephemeral", gjson.GetBytes(out, "tools.1.cache_control.type").String())
	require.Equal(t, "5m", gjson.GetBytes(out, "tools.1.cache_control.ttl").String())
	// First tool untouched
	require.False(t, gjson.GetBytes(out, "tools.0.cache_control").Exists())
placeholder

func TestApplyToolsLastCacheBreakpoint_PassesThroughClientTTL(t *testing.T) {
	body := []byte(`{"tools":[{"name":"a","input_schema":{placeholder,"cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder]placeholder`)
	out := applyToolsLastCacheBreakpoint(body)
	// User-provided ttl must be preserved.
	require.Equal(t, "1h", gjson.GetBytes(out, "tools.0.cache_control.ttl").String())
placeholder

func TestApplyToolsLastCacheBreakpoint_StripsDeferredToolCacheControl(t *testing.T) {
	body := []byte(`{"tools":[{"name":"a","custom":{"defer_loading":trueplaceholder,"cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder,{"name":"b","custom":{"defer_loading":trueplaceholderplaceholder]placeholder`)
	out := applyToolsLastCacheBreakpoint(body)

	require.False(t, gjson.GetBytes(out, "tools.0.cache_control").Exists())
	require.False(t, gjson.GetBytes(out, "tools.1.cache_control").Exists())
placeholder

func TestApplyToolsLastCacheBreakpoint_SkipsDeferredFinalTool(t *testing.T) {
	body := []byte(`{"tools":[{"name":"a","input_schema":{placeholderplaceholder,{"name":"b","defer_loading":trueplaceholder]placeholder`)
	out := applyToolsLastCacheBreakpoint(body)

	require.Equal(t, "ephemeral", gjson.GetBytes(out, "tools.0.cache_control.type").String())
	require.Equal(t, "5m", gjson.GetBytes(out, "tools.0.cache_control.ttl").String())
	require.False(t, gjson.GetBytes(out, "tools.1.cache_control").Exists())
placeholder

func TestApplyToolsLastCacheBreakpoint_OnlyLiteralTrueIsDeferred(t *testing.T) {
	body := []byte(`{"tools":[{"name":"custom-true","custom":{"defer_loading":trueplaceholder,"cache_control":{"type":"ephemeral"placeholderplaceholder,{"name":"top-level-true","defer_loading":true,"cache_control":{"type":"ephemeral"placeholderplaceholder,{"name":"false","defer_loading":false,"cache_control":{"type":"ephemeral"placeholderplaceholder,{"name":"string","defer_loading":"true","cache_control":{"type":"ephemeral"placeholderplaceholder,{"name":"number","defer_loading":1,"cache_control":{"type":"ephemeral"placeholderplaceholder,{"name":"object","defer_loading":{placeholder,"cache_control":{"type":"ephemeral"placeholderplaceholder]placeholder`)
	out := stripDeferredToolCacheControl(body)

	require.False(t, gjson.GetBytes(out, "tools.0.cache_control").Exists())
	require.False(t, gjson.GetBytes(out, "tools.1.cache_control").Exists())
	for idx := 2; idx < 6; idx++ {
		require.Equal(t, "ephemeral", gjson.GetBytes(out, fmt.Sprintf("tools.%d.cache_control.type", idx)).String())
placeholder
placeholder

func TestStripMessageCacheControl(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":[{"type":"text","text":"hi","cache_control":{"type":"ephemeral"placeholderplaceholder]placeholder]placeholder`)
	out := stripMessageCacheControl(body)
	require.False(t, gjson.GetBytes(out, "messages.0.content.0.cache_control").Exists())
placeholder

func TestAddMessageCacheBreakpoints_LastMessageOnly(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":[{"type":"text","text":"hello"placeholder]placeholder]placeholder`)
	out := addMessageCacheBreakpoints(body)
	require.Equal(t, "ephemeral", gjson.GetBytes(out, "messages.0.content.0.cache_control.type").String())
	require.Equal(t, "5m", gjson.GetBytes(out, "messages.0.content.0.cache_control.ttl").String())
placeholder

func TestAddMessageCacheBreakpoints_SecondToLastUserTurn(t *testing.T) {
	// Parrot 不变量：messages ≥ 4 时才打第二个断点，且位置是"倒数第二个 user turn"。
	body := []byte(`{"messages":[
        {"role":"user","content":[{"type":"text","text":"q1"placeholder]placeholder,
        {"role":"assistant","content":[{"type":"text","text":"a1"placeholder]placeholder,
        {"role":"user","content":[{"type":"text","text":"q2"placeholder]placeholder,
        {"role":"assistant","content":[{"type":"text","text":"a2"placeholder]placeholder
    ]placeholder`)
	out := addMessageCacheBreakpoints(body)
	// 最后一条 assistant 被打断点
	require.Equal(t, "ephemeral", gjson.GetBytes(out, "messages.3.content.0.cache_control.type").String())
	// 倒数第二个 user turn = index 0（唯一另一个 user）
	require.Equal(t, "ephemeral", gjson.GetBytes(out, "messages.0.content.0.cache_control.type").String())
	// 其他不打断点
	require.False(t, gjson.GetBytes(out, "messages.1.content.0.cache_control").Exists())
	require.False(t, gjson.GetBytes(out, "messages.2.content.0.cache_control").Exists())
placeholder

func TestAddMessageCacheBreakpoints_StringContentPromoted(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`)
	out := addMessageCacheBreakpoints(body)
	// content 升级成数组
	require.True(t, gjson.GetBytes(out, "messages.0.content").IsArray())
	require.Equal(t, "text", gjson.GetBytes(out, "messages.0.content.0.type").String())
	require.Equal(t, "hi", gjson.GetBytes(out, "messages.0.content.0.text").String())
	require.Equal(t, "5m", gjson.GetBytes(out, "messages.0.content.0.cache_control.ttl").String())
placeholder

func TestRewriteMessageCacheControlIfEnabled_DefaultKeepsClientAnchors(t *testing.T) {
	body := []byte(`{"messages":[
		{"role":"user","content":[{"type":"text","text":"stable","cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder]placeholder,
		{"role":"assistant","content":[{"type":"text","text":"ok"placeholder]placeholder,
		{"role":"user","content":[{"type":"text","text":"latest","cache_control":{"type":"ephemeral","ttl":"5m"placeholderplaceholder]placeholder
	]placeholder`)

	out := (&GatewayService{placeholder).rewriteMessageCacheControlIfEnabled(context.Background(), body)

	require.JSONEq(t, string(body), string(out))
	require.Equal(t, "1h", gjson.GetBytes(out, "messages.0.content.0.cache_control.ttl").String())
	require.Equal(t, "5m", gjson.GetBytes(out, "messages.2.content.0.cache_control.ttl").String())
placeholder

func TestRewriteMessageCacheControlIfEnabled_OptInPreservesLegacyRewrite(t *testing.T) {
	body := []byte(`{"messages":[
		{"role":"user","content":[{"type":"text","text":"stable","cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder]placeholder,
		{"role":"assistant","content":[{"type":"text","text":"ok"placeholder]placeholder,
		{"role":"user","content":[{"type":"text","text":"latest","cache_control":{"type":"ephemeral","ttl":"1h"placeholderplaceholder]placeholder,
		{"role":"assistant","content":[{"type":"text","text":"done"placeholder]placeholder
	]placeholder`)
	repo := &gatewayTTLSettingRepo{data: map[string]string{
		SettingKeyRewriteMessageCacheControl: "true",
placeholderplaceholder
	gatewayForwardingCache.Store(&cachedGatewayForwardingSettings{placeholder)
	svc := &GatewayService{settingService: NewSettingService(repo, &config.Config{placeholder)placeholder

	out := svc.rewriteMessageCacheControlIfEnabled(context.Background(), body)

	require.Equal(t, "5m", gjson.GetBytes(out, "messages.0.content.0.cache_control.ttl").String())
	require.False(t, gjson.GetBytes(out, "messages.2.content.0.cache_control").Exists())
	require.Equal(t, "5m", gjson.GetBytes(out, "messages.3.content.0.cache_control.ttl").String())
placeholder

func TestBuildToolNameRewriteFromBody_ReverseOrderedByLengthDesc(t *testing.T) {
	// 超过阈值触发动态映射，验证 ReverseOrdered 按假名长度倒序排列
	body := []byte(`{"tools":[
        {"name":"t1","input_schema":{placeholderplaceholder,
        {"name":"t2","input_schema":{placeholderplaceholder,
        {"name":"t3","input_schema":{placeholderplaceholder,
        {"name":"t4","input_schema":{placeholderplaceholder,
        {"name":"t5","input_schema":{placeholderplaceholder,
        {"name":"t6","input_schema":{placeholderplaceholder
    ]placeholder`)
	rw := buildToolNameRewriteFromBody(body)
	require.NotNil(t, rw)
	require.NotEmpty(t, rw.ReverseOrdered)
	for i := 1; i < len(rw.ReverseOrdered); i++ {
		require.GreaterOrEqual(t, len(rw.ReverseOrdered[i-1][0]), len(rw.ReverseOrdered[i][0]),
			"ReverseOrdered must be sorted by fake-name length descending")
placeholder
placeholder

func TestRestoreToolNamesInBytes_NoMapping_NoStaticMatch_IsNoop(t *testing.T) {
	data := []byte("plain text without any tool names")
	require.Equal(t, string(data), string(restoreToolNamesInBytes(data, nil)))
placeholder

// Ensure the fake name format follows Parrot's "{prefixplaceholder{name[:3]placeholder{i:02dplaceholder".
func TestBuildDynamicToolMap_FakeNameShape(t *testing.T) {
	names := []string{"alphabet", "bravo", "charlie", "delta", "echo", "foxtrot"placeholder
	m := buildDynamicToolMap(names)
	require.NotNil(t, m)
	for _, name := range names {
		fake, ok := m[name]
		require.True(t, ok)
		// fake = prefix + head3 + "%02d"
		// ends with two decimal digits
		require.Regexp(t, `^[a-z]+_[a-z0-9]{1,3placeholder\d{2placeholder$`, fake)
		head := name
		if len(head) > 3 {
			head = head[:3]
	placeholder
		require.True(t, strings.Contains(fake, head), "fake %q should contain head3 %q of %q", fake, head, name)
placeholder
placeholder
