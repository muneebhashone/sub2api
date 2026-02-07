//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ============ 临时限流单元测试 ============

// TestMatchTempUnschedKeyword 测试关键词匹配函数
func TestMatchTempUnschedKeyword(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		keywords []string
		want     string
placeholder{
		{
			name:     "match_first",
			body:     "server is overloaded",
			keywords: []string{"overloaded", "capacity"placeholder,
			want:     "overloaded",
	placeholder,
		{
			name:     "match_second",
			body:     "no capacity available",
			keywords: []string{"overloaded", "capacity"placeholder,
			want:     "capacity",
	placeholder,
		{
			name:     "no_match",
			body:     "internal error",
			keywords: []string{"overloaded", "capacity"placeholder,
			want:     "",
	placeholder,
		{
			name:     "empty_body",
			body:     "",
			keywords: []string{"overloaded"placeholder,
			want:     "",
	placeholder,
		{
			name:     "empty_keywords",
			body:     "server is overloaded",
			keywords: []string{placeholder,
			want:     "",
	placeholder,
		{
			name:     "whitespace_keyword",
			body:     "server is overloaded",
			keywords: []string{"  ", "overloaded"placeholder,
			want:     "overloaded",
	placeholder,
		{
			// matchTempUnschedKeyword 期望 body 已经是小写的
			// 所以要测试大小写不敏感匹配，需要传入小写的 body
			name:     "case_insensitive_body_lowered",
			body:     "server is overloaded", // body 已经是小写
			keywords: []string{"OVERLOADED"placeholder, // keyword 会被转为小写比较
			want:     "OVERLOADED",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchTempUnschedKeyword(tt.body, tt.keywords)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// TestAccountIsSchedulable_TempUnschedulable 测试临时限流账号不可调度
func TestAccountIsSchedulable_TempUnschedulable(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	past := time.Now().Add(-10 * time.Minute)

	tests := []struct {
		name    string
		account *Account
		want    bool
placeholder{
		{
			name: "temp_unschedulable_active",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: &future,
		placeholder,
			want: false,
	placeholder,
		{
			name: "temp_unschedulable_expired",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: &past,
		placeholder,
			want: true,
	placeholder,
		{
			name: "no_temp_unschedulable",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: nil,
		placeholder,
			want: true,
	placeholder,
		{
			name: "temp_unschedulable_with_rate_limit",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: &future,
				RateLimitResetAt:       &past, // 过期的限流不影响
		placeholder,
			want: false, // 临时限流生效
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.IsSchedulable()
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// TestAccount_IsTempUnschedulableEnabled 测试临时限流开关
func TestAccount_IsTempUnschedulableEnabled(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    bool
placeholder{
		{
			name: "enabled",
			account: &Account{
		placeholder
					"temp_unschedulable_enabled": true,
			placeholder,
		placeholder,
			want: true,
	placeholder,
		{
			name: "disabled",
			account: &Account{
		placeholder
					"temp_unschedulable_enabled": false,
			placeholder,
		placeholder,
			want: false,
	placeholder,
		{
			name: "not_set",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			want: false,
	placeholder,
		{
			name:    "nil_credentials",
			account: &Account{placeholder,
			want:    false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.IsTempUnschedulableEnabled()
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// TestAccount_GetTempUnschedulableRules 测试获取临时限流规则
func TestAccount_GetTempUnschedulableRules(t *testing.T) {
	tests := []struct {
		name      string
		account   *Account
		wantCount int
placeholder{
		{
			name: "has_rules",
			account: &Account{
		placeholder
					"temp_unschedulable_rules": []any{
						map[string]any{
							"error_code":       float64(503),
							"keywords":         []any{"overloaded"placeholder,
							"duration_minutes": float64(5),
					placeholder,
						map[string]any{
							"error_code":       float64(500),
							"keywords":         []any{"internal"placeholder,
							"duration_minutes": float64(10),
					placeholder,
				placeholder,
			placeholder,
		placeholder,
			wantCount: 2,
	placeholder,
		{
			name: "empty_rules",
			account: &Account{
		placeholder
					"temp_unschedulable_rules": []any{placeholder,
			placeholder,
		placeholder,
			wantCount: 0,
	placeholder,
		{
			name: "no_rules",
			account: &Account{
		placeholderplaceholder,
		placeholder,
			wantCount: 0,
	placeholder,
		{
			name:      "nil_credentials",
			account:   &Account{placeholder,
			wantCount: 0,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := tt.account.GetTempUnschedulableRules()
			require.Len(t, rules, tt.wantCount)
	placeholder)
placeholder
placeholder

// TestTempUnschedulableRule_Parse 测试规则解析
func TestTempUnschedulableRule_Parse(t *testing.T) {
	account := &Account{
placeholder
			"temp_unschedulable_rules": []any{
				map[string]any{
					"error_code":       float64(503),
					"keywords":         []any{"overloaded", "capacity"placeholder,
					"duration_minutes": float64(5),
			placeholder,
		placeholder,
	placeholder,
placeholder

	rules := account.GetTempUnschedulableRules()
	require.Len(t, rules, 1)

	rule := rules[0]
	require.Equal(t, 503, rule.ErrorCode)
	require.Equal(t, []string{"overloaded", "capacity"placeholder, rule.Keywords)
	require.Equal(t, 5, rule.DurationMinutes)
placeholder

// TestTruncateTempUnschedMessage 测试消息截断
func TestTruncateTempUnschedMessage(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		maxBytes int
		want     string
placeholder{
		{
			name:     "short_message",
			body:     []byte("short"),
			maxBytes: 100,
			want:     "short",
	placeholder,
		{
			// 截断后会 TrimSpace，所以末尾的空格会被移除
			name:     "truncate_long_message",
			body:     []byte("this is a very long message that needs to be truncated"),
			maxBytes: 20,
			want:     "this is a very long", // 截断后 TrimSpace
	placeholder,
		{
			name:     "empty_body",
			body:     []byte{placeholder,
			maxBytes: 100,
			want:     "",
	placeholder,
		{
			name:     "zero_max_bytes",
			body:     []byte("test"),
			maxBytes: 0,
			want:     "",
	placeholder,
		{
			name:     "whitespace_trimmed",
			body:     []byte("  test  "),
			maxBytes: 100,
			want:     "test",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateTempUnschedMessage(tt.body, tt.maxBytes)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// TestTempUnschedState 测试临时限流状态结构
func TestTempUnschedState(t *testing.T) {
	now := time.Now()
	until := now.Add(5 * time.Minute)

	state := &TempUnschedState{
		UntilUnix:       until.Unix(),
		TriggeredAtUnix: now.Unix(),
		StatusCode:      503,
		MatchedKeyword:  "overloaded",
		RuleIndex:       0,
		ErrorMessage:    "Server is overloaded",
placeholder

	require.Equal(t, 503, state.StatusCode)
	require.Equal(t, "overloaded", state.MatchedKeyword)
	require.Equal(t, 0, state.RuleIndex)

	// 验证时间戳
	require.Equal(t, until.Unix(), state.UntilUnix)
	require.Equal(t, now.Unix(), state.TriggeredAtUnix)
placeholder

// TestAccount_TempUnschedulableUntil 测试临时限流时间字段
func TestAccount_TempUnschedulableUntil(t *testing.T) {
	future := time.Now().Add(10 * time.Minute)
	past := time.Now().Add(-10 * time.Minute)

	tests := []struct {
		name        string
		account     *Account
		schedulable bool
placeholder{
		{
			name: "active_temp_unsched_not_schedulable",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: &future,
		placeholder,
			schedulable: false,
	placeholder,
		{
			name: "expired_temp_unsched_is_schedulable",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: &past,
		placeholder,
			schedulable: true,
	placeholder,
		{
			name: "nil_temp_unsched_is_schedulable",
			account: &Account{
				Status:                 StatusActive,
				Schedulable:            true,
				TempUnschedulableUntil: nil,
		placeholder,
			schedulable: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.account.IsSchedulable()
			require.Equal(t, tt.schedulable, got)
	placeholder)
placeholder
placeholder
