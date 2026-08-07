package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

func TestBuildGrokVoiceURL_UsesAPIDefaultForCLIProxyBase(t *testing.T) {
	account := &Account{
		Platform: PlatformGrok,
		Type:     AccountTypeOAuth,
placeholder
			"base_url": xai.DefaultCLIBaseURL,
	placeholder,
placeholder
	url, err := buildGrokVoiceURL(account, nil, "tts")
placeholder
	require.Equal(t, xai.DefaultBaseURL+"/tts", url)

	url, err = buildGrokVoiceURL(account, nil, "realtime")
placeholder
	require.Equal(t, xai.DefaultBaseURL+"/realtime", url)
placeholder

func TestBuildGrokVoiceURL_EmptyBaseFallsBackToAPI(t *testing.T) {
	account := &Account{
		Platform:    PlatformGrok,
		Type:        AccountTypeOAuth,
placeholderplaceholder,
placeholder
	url, err := buildGrokVoiceURL(account, nil, "stt")
placeholder
	require.Equal(t, xai.DefaultBaseURL+"/stt", url)
placeholder

func TestBuildGrokVoiceURL_RequiresEndpoint(t *testing.T) {
	account := &Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	_, err := buildGrokVoiceURL(account, nil, "  ")
placeholder
placeholder

func TestForwardGrokVoice_RejectsNonGrok(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	_, err := svc.ForwardGrokVoice(context.Background(), nil, &Account{Platform: PlatformOpenAIplaceholder, "tts", []byte(`{placeholder`), "application/json")
placeholder
	require.Contains(t, err.Error(), "not supported")
placeholder

func TestForwardGrokVoice_RejectsUnknownEndpoint(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	_, err := svc.ForwardGrokVoice(context.Background(), nil, &Account{Platform: PlatformGrokplaceholder, "unknown", []byte(`{placeholder`), "application/json")
placeholder
	require.Contains(t, err.Error(), "unsupported")
placeholder
