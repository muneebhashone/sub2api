package openai

import "testing"

func TestIsCodexCLIRequest(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want bool
placeholder{
		{name: "codex_cli_rs 前缀", ua: "codex_cli_rs/0.1.0", want: trueplaceholder,
		{name: "codex_vscode 前缀", ua: "codex_vscode/1.2.3", want: trueplaceholder,
		{name: "大小写混合", ua: "Codex_CLI_Rs/0.1.0", want: trueplaceholder,
		{name: "复合 UA 包含 codex", ua: "Mozilla/5.0 codex_cli_rs/0.1.0", want: trueplaceholder,
		{name: "空白包裹", ua: "  codex_vscode/1.2.3  ", want: trueplaceholder,
		{name: "非 codex", ua: "curl/8.0.1", want: falseplaceholder,
		{name: "空字符串", ua: "", want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCodexCLIRequest(tt.ua)
			if got != tt.want {
				t.Fatalf("IsCodexCLIRequest(%q) = %v, want %v", tt.ua, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestIsCodexOfficialClientRequest(t *testing.T) {
	tests := []struct {
		name string
		ua   string
		want bool
placeholder{
		{name: "codex_cli_rs 前缀", ua: "codex_cli_rs/0.98.0", want: trueplaceholder,
		{name: "codex_vscode 前缀", ua: "codex_vscode/1.0.0", want: trueplaceholder,
		{name: "codex_app 前缀", ua: "codex_app/0.1.0", want: trueplaceholder,
		{name: "codex_chatgpt_desktop 前缀", ua: "codex_chatgpt_desktop/1.0.0", want: trueplaceholder,
		{name: "codex_atlas 前缀", ua: "codex_atlas/1.0.0", want: trueplaceholder,
		{name: "codex_exec 前缀", ua: "codex_exec/0.1.0", want: trueplaceholder,
		{name: "codex_sdk_ts 前缀", ua: "codex_sdk_ts/0.1.0", want: trueplaceholder,
		{name: "Codex 桌面 UA", ua: "Codex Desktop/1.2.3", want: trueplaceholder,
		{name: "复合 UA 包含 codex_app", ua: "Mozilla/5.0 codex_app/0.1.0", want: trueplaceholder,
		{name: "大小写混合", ua: "Codex_VSCode/1.2.3", want: trueplaceholder,
		{name: "非 codex", ua: "curl/8.0.1", want: falseplaceholder,
		{name: "空字符串", ua: "", want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCodexOfficialClientRequest(tt.ua)
			if got != tt.want {
				t.Fatalf("IsCodexOfficialClientRequest(%q) = %v, want %v", tt.ua, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestIsCodexOfficialClientOriginator(t *testing.T) {
	tests := []struct {
		name       string
		originator string
		want       bool
placeholder{
		{name: "codex_cli_rs", originator: "codex_cli_rs", want: trueplaceholder,
		{name: "codex_vscode", originator: "codex_vscode", want: trueplaceholder,
		{name: "codex_app", originator: "codex_app", want: trueplaceholder,
		{name: "codex_chatgpt_desktop", originator: "codex_chatgpt_desktop", want: trueplaceholder,
		{name: "codex_atlas", originator: "codex_atlas", want: trueplaceholder,
		{name: "codex_exec", originator: "codex_exec", want: trueplaceholder,
		{name: "codex_sdk_ts", originator: "codex_sdk_ts", want: trueplaceholder,
		{name: "Codex 前缀", originator: "Codex Desktop", want: trueplaceholder,
		{name: "空白包裹", originator: "  codex_vscode  ", want: trueplaceholder,
		{name: "非 codex", originator: "my_client", want: falseplaceholder,
		{name: "空字符串", originator: "", want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCodexOfficialClientOriginator(tt.originator)
			if got != tt.want {
				t.Fatalf("IsCodexOfficialClientOriginator(%q) = %v, want %v", tt.originator, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestIsCodexOfficialClientByHeaders(t *testing.T) {
	tests := []struct {
		name       string
		ua         string
		originator string
		want       bool
placeholder{
		{name: "仅 originator 命中 desktop", originator: "Codex Desktop", want: trueplaceholder,
		{name: "仅 originator 命中 vscode", originator: "codex_vscode", want: trueplaceholder,
		{name: "仅 ua 命中 desktop", ua: "Codex Desktop/1.2.3", want: trueplaceholder,
		{name: "ua 与 originator 都未命中", ua: "curl/8.0.1", originator: "my_client", want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCodexOfficialClientByHeaders(tt.ua, tt.originator)
			if got != tt.want {
				t.Fatalf("IsCodexOfficialClientByHeaders(%q, %q) = %v, want %v", tt.ua, tt.originator, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
