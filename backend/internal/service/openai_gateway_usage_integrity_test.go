//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequiresBillableGrokChatUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		account *Account
		models  []string
		want    bool
placeholder{
		{
			name:    "grok platform",
			account: &Account{Platform: PlatformGrokplaceholder,
			models:  []string{"alias"placeholder,
			want:    true,
	placeholder,
		{
			name:    "OpenAI-compatible requested Grok model",
			account: &Account{Platform: PlatformOpenAIplaceholder,
			models:  []string{"grok-4.5"placeholder,
			want:    true,
	placeholder,
		{
			name:    "OpenAI-compatible mapped Grok model",
			account: &Account{Platform: PlatformOpenAIplaceholder,
			models:  []string{"alias", "grok-4.5"placeholder,
			want:    true,
	placeholder,
		{
			name:    "xAI-qualified Grok model",
			account: &Account{Platform: PlatformOpenAIplaceholder,
			models:  []string{"xai/grok-4.5"placeholder,
			want:    true,
	placeholder,
		{
			name:    "ordinary OpenAI model",
			account: &Account{Platform: PlatformOpenAIplaceholder,
			models:  []string{"gpt-5.4"placeholder,
			want:    false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, requiresBillableGrokChatUsage(tt.account, tt.models...))
	placeholder)
placeholder
placeholder
