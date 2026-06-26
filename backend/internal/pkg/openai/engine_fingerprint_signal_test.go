package openai

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// hdr 构造一个 http.Header(键值对)。
func hdr(kv ...string) http.Header {
	h := http.Header{placeholder
	for i := 0; i+1 < len(kv); i += 2 {
		h.Set(kv[i], kv[i+1])
placeholder
	return h
placeholder

func TestEvaluateEngineFingerprint_DefaultSeed(t *testing.T) {
	sigs := DefaultEngineFingerprintSignals // 仅 x-codex- 前缀 Required
	cases := []struct {
		name string
		h    http.Header
		body string
		want bool
placeholder{
		{"R1 真CLI 带x-codex-window-id", hdr("x-codex-window-id", "a1", "session-id", "u1"), ``, trueplaceholder,
		{"R2 纯伪装 无指纹", hdr("user-agent", "codex/1"), ``, falseplaceholder,
		{"R3 仅body有", hdr(), `{"client_metadata":{"x-codex-window-id":"c3"placeholderplaceholder`, falseplaceholder,
		{"R4 旧版 仅session_id无x-codex-", hdr("session_id", "u4"), ``, falseplaceholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, EvaluateEngineFingerprint(tc.h, []byte(tc.body), sigs))
	placeholder)
placeholder
placeholder

func TestEvaluateEngineFingerprint_Rules(t *testing.T) {
	exactSession := EngineFingerprintSignal{Type: FingerprintSignalHeaderExact, Match: []string{"session-id", "session_id"placeholder, Required: trueplaceholder
	prefixCodex := EngineFingerprintSignal{Type: FingerprintSignalHeaderPrefix, Match: []string{"x-codex-"placeholder, Required: trueplaceholder
	bodyWin := EngineFingerprintSignal{Type: FingerprintSignalBodyPath, Match: []string{"client_metadata.x-codex-window-id"placeholder, Required: trueplaceholder

	t.Run("行内变体OR: 配置session-id 命中下划线session_id", func(t *testing.T) {
		require.True(t, EvaluateEngineFingerprint(hdr("session_id", "x"), nil, []EngineFingerprintSignal{exactSessionplaceholder))
placeholder)
	t.Run("跨条AND: 勾x-codex-与session 缺一即拒", func(t *testing.T) {
		both := []EngineFingerprintSignal{prefixCodex, exactSessionplaceholder
		require.True(t, EvaluateEngineFingerprint(hdr("x-codex-window-id", "a", "session-id", "b"), nil, both))
		require.False(t, EvaluateEngineFingerprint(hdr("session-id", "b"), nil, both)) // 缺 x-codex-
placeholder)
	t.Run("body_path 命中/ body空", func(t *testing.T) {
		require.True(t, EvaluateEngineFingerprint(hdr(), []byte(`{"client_metadata":{"x-codex-window-id":"1"placeholderplaceholder`), []EngineFingerprintSignal{bodyWinplaceholder))
		require.False(t, EvaluateEngineFingerprint(hdr(), nil, []EngineFingerprintSignal{bodyWinplaceholder))
placeholder)
	t.Run("无任何Required → true", func(t *testing.T) {
		none := []EngineFingerprintSignal{{Type: FingerprintSignalHeaderPrefix, Match: []string{"x-codex-"placeholder, Required: falseplaceholderplaceholder
		require.True(t, EvaluateEngineFingerprint(hdr(), nil, none))
		require.True(t, EvaluateEngineFingerprint(hdr(), nil, nil))
placeholder)
placeholder

func TestParseAndValidateEngineFingerprintSignals(t *testing.T) {
	t.Run("空串=合法空", func(t *testing.T) {
		sigs, ok := ParseEngineFingerprintSignals("")
		require.True(t, ok)
		require.Nil(t, sigs)
		require.NoError(t, ValidateEngineFingerprintSignalsJSON(""))
placeholder)
	t.Run("合法数组", func(t *testing.T) {
		raw := `[{"type":"header_prefix","match":["x-codex-"],"required":trueplaceholder]`
		sigs, ok := ParseEngineFingerprintSignals(raw)
		require.True(t, ok)
		require.Len(t, sigs, 1)
		require.NoError(t, ValidateEngineFingerprintSignalsJSON(raw))
placeholder)
	t.Run("非法JSON", func(t *testing.T) {
		_, ok := ParseEngineFingerprintSignals("not json")
		require.False(t, ok)
		require.Error(t, ValidateEngineFingerprintSignalsJSON("not json"))
placeholder)
	t.Run("非法type 被校验拒绝", func(t *testing.T) {
		require.Error(t, ValidateEngineFingerprintSignalsJSON(`[{"type":"bogus","match":["x"]placeholder]`))
placeholder)
	t.Run("match全空 被校验拒绝", func(t *testing.T) {
		require.Error(t, ValidateEngineFingerprintSignalsJSON(`[{"type":"header_exact","match":["",""]placeholder]`))
placeholder)
	t.Run("默认种子JSON 可解析且只勾x-codex-", func(t *testing.T) {
		sigs, ok := ParseEngineFingerprintSignals(DefaultEngineFingerprintSignalsJSON())
		require.True(t, ok)
		requiredTypes := []string{placeholder
		for _, s := range sigs {
			if s.Required {
				requiredTypes = append(requiredTypes, s.Type+":"+s.Match[0])
		placeholder
	placeholder
		require.Equal(t, []string{"header_prefix:x-codex-"placeholder, requiredTypes)
placeholder)
placeholder
