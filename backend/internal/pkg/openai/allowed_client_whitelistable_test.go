package openai

import "testing"

func TestAllowedClientEntry_IsWhitelistable(t *testing.T) {
	cases := []struct {
		name  string
		entry AllowedClientEntry
		want  bool
placeholder{
		{name: "完整条目可白名单", entry: AllowedClientEntry{Originator: "opencode", UAContains: []string{"opencode/"placeholderplaceholder, want: trueplaceholder,
		{name: "多个有效 marker 可白名单", entry: AllowedClientEntry{Originator: "x", UAContains: []string{"a/", "b/"placeholderplaceholder, want: trueplaceholder,
		{name: "缺 originator → 不可(静默失效)", entry: AllowedClientEntry{UAContains: []string{"opencode/"placeholderplaceholder, want: falseplaceholder,
		{name: "originator 全空白 → 不可", entry: AllowedClientEntry{Originator: "   ", UAContains: []string{"opencode/"placeholderplaceholder, want: falseplaceholder,
		{name: "缺 ua_contains → 不可(静默失效)", entry: AllowedClientEntry{Originator: "opencode"placeholder, want: falseplaceholder,
		{name: "ua_contains 全空白 → 不可", entry: AllowedClientEntry{Originator: "opencode", UAContains: []string{"", "  "placeholderplaceholder, want: falseplaceholder,
		{name: "含一个空白 marker → 不可(空白会让整条 IsAllowedClientMatch 永不命中)", entry: AllowedClientEntry{Originator: "opencode", UAContains: []string{"opencode/", ""placeholderplaceholder, want: falseplaceholder,
placeholder
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.IsWhitelistable(); got != tt.want {
				t.Fatalf("IsWhitelistable(%+v) = %v, want %v", tt.entry, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
