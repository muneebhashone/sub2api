package service

import (
	"context"
	"io"
	"sync/atomic"
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

func TestBuildGrokVoiceURL_EncodesCustomVoicePathSegments(t *testing.T) {
	account := &Account{Platform: PlatformGrok, Type: AccountTypeOAuthplaceholder
	got, err := buildGrokVoiceURL(account, nil, "custom-voices/nlbqfwie/audio")
placeholder
	require.Equal(t, xai.DefaultBaseURL+"/custom-voices/nlbqfwie/audio", got)

	_, err = buildGrokVoiceURL(account, nil, "custom-voices/../audio")
placeholder
placeholder

func TestForwardGrokVoice_RejectsNonGrok(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	_, err := svc.ForwardGrokVoice(context.Background(), nil, &Account{Platform: PlatformOpenAIplaceholder, "tts", []byte(`{placeholder`), "application/json")
placeholder
	require.Contains(t, err.Error(), "not supported")
placeholder

func TestAwaitGrokRealtimeAudioObservedReadsFlagAfterRelayExits(t *testing.T) {
	errCh := make(chan error, 1)
	var observed atomic.Bool
	go func() {
		observed.Store(true)
		errCh <- io.EOF
placeholder()
	got, err := awaitGrokRealtimeAudioObserved(errCh, &observed)
	require.ErrorIs(t, err, io.EOF)
	require.True(t, got, "audioObserved must be read after the relay returns, not before <-errCh")
placeholder

func TestGrokRealtimeEventHasAudio(t *testing.T) {
	require.False(t, grokRealtimeEventHasAudio([]byte(`{"type":"session.created"placeholder`)))
	require.False(t, grokRealtimeEventHasAudio([]byte(`{"type":"response.audio_transcript.delta","delta":"hi"placeholder`)))
	require.False(t, grokRealtimeEventHasAudio([]byte(`{"type":"response.audio.delta","delta":""placeholder`)))
	require.True(t, grokRealtimeEventHasAudio([]byte(`{"type":"response.audio.delta","delta":"abc"placeholder`)))
	require.True(t, grokRealtimeEventHasAudio([]byte(`{"type":"response.output_audio.delta","audio":"abc"placeholder`)))
placeholder

func TestForwardGrokVoice_RejectsUnknownEndpoint(t *testing.T) {
	svc := &OpenAIGatewayService{placeholder
	_, err := svc.ForwardGrokVoice(context.Background(), nil, &Account{Platform: PlatformGrokplaceholder, "unknown", []byte(`{placeholder`), "application/json")
placeholder
	require.Contains(t, err.Error(), "unsupported")
placeholder
