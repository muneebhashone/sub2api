package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBedrockSignerFromAccount_DefaultRegion(t *testing.T) {
	account := &Account{
		Platform: PlatformAnthropic,
		Type:     AccountTypeBedrock,
placeholder
			"aws_access_key_id":     "test-akid",
			"aws_secret_access_key": "test-secret",
	placeholder,
placeholder

	signer, err := NewBedrockSignerFromAccount(account)
placeholder
	require.NotNil(t, signer)
	assert.Equal(t, defaultBedrockRegion, signer.region)
placeholder

func TestFilterBetaTokens(t *testing.T) {
	tokens := []string{"placeholder", "tool-search-tool-2025-10-19"placeholder
	filterSet := map[string]struct{placeholder{
		"tool-search-tool-2025-10-19": {placeholder,
placeholder

	assert.Equal(t, []string{"placeholder"placeholder, filterBetaTokens(tokens, filterSet))
	assert.Equal(t, tokens, filterBetaTokens(tokens, nil))
	assert.Nil(t, filterBetaTokens(nil, filterSet))
placeholder
