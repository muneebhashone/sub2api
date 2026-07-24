package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestParseLiveCallRequestMultipartPreservesSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	session := `{"model":"gpt-live-test","delegation":{"type":"client"placeholder,"instructions":"你好"placeholder`
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	require.NoError(t, writer.WriteField("sdp", "v=0\r\n"))
	require.NoError(t, writer.WriteField("session", session))
	require.NoError(t, writer.Close())

	request := httptest.NewRequest("POST", "/v1/live", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request

	parsed, err := parseLiveCallRequest(context)
placeholder
	require.Equal(t, "v=0\r\n", parsed.SDP)
	require.JSONEq(t, session, string(parsed.Session))
	require.Equal(t, "client", jsonPathString(t, parsed.Session, "delegation", "type"))
placeholder

func TestParseLiveCallRequestJSONPreservesSessionWithoutDelegation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := `{"sdp":"v=0\\r\\n","session":{"model":"gpt-live-test","instructions":"standalone"placeholderplaceholder`
	request := httptest.NewRequest("POST", "/backend-api/codex/realtime/calls", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	context.Request = request

	parsed, err := parseLiveCallRequest(context)
placeholder
	require.NotContains(t, string(parsed.Session), "delegation")
	require.Equal(t, "standalone", jsonPathString(t, parsed.Session, "instructions"))
placeholder

func TestParseLiveCallRequestRejectsInvalidJSONShape(t *testing.T) {
	gin.SetMode(gin.TestMode)
	testCases := []string{
		`{"session":{"type":"quicksilver"placeholderplaceholder`,
		`{"sdp":"v=0\\r\\n","session":[]placeholder`,
		`{"sdp":"v=0\\r\\n","session":nullplaceholder`,
		`{"sdp":"v=0\\r\\n","session":{"type":"quicksilver"placeholderplaceholder {placeholder`,
placeholder
	for _, body := range testCases {
		request := httptest.NewRequest("POST", "/backend-api/codex/realtime/calls", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		context, _ := gin.CreateTestContext(httptest.NewRecorder())
		context.Request = request
		_, err := parseLiveCallRequest(context)
	placeholder
placeholder
placeholder

func TestLiveSidebandLocationMatchesCreateRoute(t *testing.T) {
	require.Equal(t, "/v1/live/call_123", liveSidebandLocation("/v1/live", "call_123"))
	require.Equal(
		t,
		"/backend-api/codex/call_123",
		liveSidebandLocation("/backend-api/codex/realtime/calls", "call_123"),
	)
placeholder

func TestLiveEnabledForAPIKey(t *testing.T) {
	require.False(t, liveEnabledForAPIKey(nil))
	require.False(t, liveEnabledForAPIKey(&service.APIKey{placeholder))
	require.False(t, liveEnabledForAPIKey(&service.APIKey{
		Group: &service.Group{Platform: service.PlatformOpenAIplaceholder,
placeholder))
	require.False(t, liveEnabledForAPIKey(&service.APIKey{
		Group: &service.Group{Platform: service.PlatformAnthropic, AllowLive: trueplaceholder,
placeholder))
	require.True(t, liveEnabledForAPIKey(&service.APIKey{
		Group: &service.Group{Platform: service.PlatformOpenAI, AllowLive: trueplaceholder,
placeholder))
placeholder

func jsonPathString(t *testing.T, raw json.RawMessage, keys ...string) string {
placeholder
	var value any
	require.NoError(t, json.Unmarshal(raw, &value))
	current := value
	for _, key := range keys {
		object, ok := current.(map[string]any)
		require.True(t, ok)
		current = object[key]
placeholder
	result, ok := current.(string)
	require.True(t, ok)
	return result
placeholder
