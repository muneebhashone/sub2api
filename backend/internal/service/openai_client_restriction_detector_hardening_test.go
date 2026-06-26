package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func hdrCtx(h map[string]string) *gin.Context {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	for k, v := range h {
		c.Request.Header.Set(k, v)
placeholder
	return c
placeholder

func codexOnlyAccount() *Account {
placeholderPlatform: PlatformOpenAI, Type: AccountTypeOAuth, Extra: map[string]any{"codex_cli_only": trueplaceholderplaceholder
placeholder

func TestDetect_N1_StrictOfficialUA(t *testing.T) {
	det := NewOpenAICodexClientRestrictionDetector(nil)
	acc := codexOnlyAccount()
	// 构造「中段 codex token」伪装：首段带可解析版本(绕过版本门)，codex_app 在中段。
	// 旧 lax(Contains) 会判为官方 UA → 放行(空策略指纹门开放)；N1 收紧后应判非官方 → NotMatchedUA。
	ua := "x/0.141.0 codex_app/0.141.0"
	r := det.Detect(hdrCtx(map[string]string{"User-Agent": uaplaceholder), acc, CodexRestrictionPolicy{placeholder, nil)
	require.False(t, r.Matched)
	require.Equal(t, CodexClientRestrictionReasonNotMatchedUA, r.Reason)
placeholder

func TestDetectCodexClientRestriction_NilSettingServiceFailsClosed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// settingService 缺失（仅测试/误配可达）：账号已开 codex_cli_only、官方 UA、但无 x-codex- 指纹头。
	// 零值 policy 不得让指纹门失败开放——gateway 应回退默认种子指纹信号并拒（MissingEngineFingerprint）。
	s := &OpenAIGatewayService{placeholder
	r := s.detectCodexClientRestriction(hdrCtx(map[string]string{"User-Agent": "codex_cli_rs/0.141.0 (x)"placeholder), codexOnlyAccount(), nil)
	require.True(t, r.Enabled)
	require.False(t, r.Matched)
	require.Equal(t, CodexClientRestrictionReasonMissingEngineFingerprint, r.Reason)
placeholder

func TestDetect_Hardening(t *testing.T) {
	det := NewOpenAICodexClientRestrictionDetector(nil)
	acc := codexOnlyAccount()
	fp := map[string]string{"x-codex-installation-id": "i1"placeholder // 引擎指纹

	t.Run("黑名单优先于官方身份", func(t *testing.T) {
		pol := CodexRestrictionPolicy{Blacklist: []openai.AllowedClientEntry{{Originator: "codex_cli_rs"placeholderplaceholderplaceholder
		h := map[string]string{"User-Agent": "codex_cli_rs/0.141.0 (x)", "originator": "codex_cli_rs", "x-codex-installation-id": "i1"placeholder
		r := det.Detect(hdrCtx(h), acc, pol, nil)
		require.False(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonBlacklisted, r.Reason)
placeholder)

	t.Run("strict 指纹缺失→拒(即便官方 UA)", func(t *testing.T) {
		r := det.Detect(hdrCtx(map[string]string{"User-Agent": "codex_cli_rs/0.141.0 (x)"placeholder), acc, CodexRestrictionPolicy{EngineFingerprintSignals: openai.DefaultEngineFingerprintSignalsplaceholder, nil)
		require.False(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonMissingEngineFingerprint, r.Reason)
placeholder)

	t.Run("strict 带指纹→放行", func(t *testing.T) {
		h := map[string]string{"User-Agent": "codex_cli_rs/0.141.0 (x)"placeholder
		for k, v := range fp {
			h[k] = v
	placeholder
		r := det.Detect(hdrCtx(h), acc, CodexRestrictionPolicy{EngineFingerprintSignals: openai.DefaultEngineFingerprintSignalsplaceholder, nil)
		require.True(t, r.Matched)
placeholder)

	t.Run("版本过低→拒", func(t *testing.T) {
		h := map[string]string{"User-Agent": "codex_cli_rs/0.130.0 (x)"placeholder
		for k, v := range fp {
			h[k] = v
	placeholder
		r := det.Detect(hdrCtx(h), acc, CodexRestrictionPolicy{MinCodexVersion: "0.141.0"placeholder, nil)
		require.False(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonVersionTooLow, r.Reason)
placeholder)

	t.Run("白名单 app-server 新 client + 指纹→放行", func(t *testing.T) {
		h := map[string]string{"User-Agent": "opencode/0.141.0 (x)", "originator": "opencode"placeholder
		for k, v := range fp {
			h[k] = v
	placeholder
		pol := CodexRestrictionPolicy{
			Whitelist:                []openai.AllowedClientEntry{{Originator: "opencode", UAContains: []string{"opencode/"placeholderplaceholderplaceholder,
			EngineFingerprintSignals: openai.DefaultEngineFingerprintSignals,
	placeholder
		r := det.Detect(hdrCtx(h), acc, pol, nil)
		require.True(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonMatchedWhitelistClient, r.Reason)
placeholder)

	t.Run("非 codex 不命中→拒", func(t *testing.T) {
		r := det.Detect(hdrCtx(map[string]string{"User-Agent": "curl/8"placeholder), acc, CodexRestrictionPolicy{placeholder, nil)
		require.False(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonNotMatchedUA, r.Reason)
placeholder)

	t.Run("版本不可识别→拒(originator 命中但 UA 无版本)", func(t *testing.T) {
		h := map[string]string{"User-Agent": "curl/8.0", "originator": "codex_chatgpt_desktop"placeholder
		for k, v := range fp {
			h[k] = v
	placeholder
		r := det.Detect(hdrCtx(h), acc, CodexRestrictionPolicy{placeholder, nil)
		require.False(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonVersionUndetectable, r.Reason)
placeholder)

	t.Run("版本过高→拒", func(t *testing.T) {
		h := map[string]string{"User-Agent": "codex_cli_rs/0.200.0 (x)"placeholder
		for k, v := range fp {
			h[k] = v
	placeholder
		r := det.Detect(hdrCtx(h), acc, CodexRestrictionPolicy{MaxCodexVersion: "0.141.0"placeholder, nil)
		require.False(t, r.Matched)
		require.Equal(t, CodexClientRestrictionReasonVersionTooHigh, r.Reason)
placeholder)

	t.Run("max 边界(==max)放行", func(t *testing.T) {
		h := map[string]string{"User-Agent": "codex_cli_rs/0.141.0 (x)"placeholder
		for k, v := range fp {
			h[k] = v
	placeholder
		r := det.Detect(hdrCtx(h), acc, CodexRestrictionPolicy{MaxCodexVersion: "0.141.0"placeholder, nil)
		require.True(t, r.Matched)
placeholder)
placeholder
