package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncBillingHeaderVersion(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		userAgent string
		wantSub   string // substring expected in result
		unchanged bool   // expect body to remain the same
placeholder{
		{
			name:      "replaces cc_version preserving message-derived suffix",
			body:      `{"system":[{"type":"text","text":"x-anthropic-billing-header: cc_version=2.1.81.df2; cc_entrypoint=cli; cch=00000;"placeholder,{"type":"text","text":"You are Claude Code.","cache_control":{"type":"ephemeral"placeholderplaceholder],"messages":[]placeholder`,
			userAgent: "claude-cli/2.1.22 (external, cli)",
			wantSub:   "cc_version=2.1.22.df2",
	placeholder,
		{
			name:      "no billing header in system",
			body:      `{"system":[{"type":"text","text":"You are Claude Code."placeholder],"messages":[]placeholder`,
			userAgent: "claude-cli/2.1.22",
			unchanged: true,
	placeholder,
		{
			name:      "no system field",
			body:      `{"messages":[]placeholder`,
			userAgent: "claude-cli/2.1.22",
			unchanged: true,
	placeholder,
		{
			name:      "user-agent without version",
			body:      `{"system":[{"type":"text","text":"x-anthropic-billing-header: cc_version=2.1.81; cc_entrypoint=cli; cch=00000;"placeholder],"messages":[]placeholder`,
			userAgent: "Mozilla/5.0",
			unchanged: true,
	placeholder,
		{
			name:      "empty user-agent",
			body:      `{"system":[{"type":"text","text":"x-anthropic-billing-header: cc_version=2.1.81; cc_entrypoint=cli; cch=00000;"placeholder],"messages":[]placeholder`,
			userAgent: "",
			unchanged: true,
	placeholder,
		{
			name:      "version already matches",
			body:      `{"system":[{"type":"text","text":"x-anthropic-billing-header: cc_version=2.1.22; cc_entrypoint=cli; cch=00000;"placeholder],"messages":[]placeholder`,
			userAgent: "claude-cli/2.1.22",
			unchanged: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := syncBillingHeaderVersion([]byte(tt.body), tt.userAgent)
			if tt.unchanged {
				assert.Equal(t, tt.body, string(result), "body should remain unchanged")
		placeholder else {
				assert.Contains(t, string(result), tt.wantSub)
				// Ensure old semver is gone
				assert.NotContains(t, string(result), "cc_version=2.1.81")
		placeholder
	placeholder)
placeholder
placeholder
