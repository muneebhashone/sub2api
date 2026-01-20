package repository

import (
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/imroc/req/v3"
	"github.com/stretchr/testify/require"
)

func forceHTTPVersion(t *testing.T, client *req.Client) string {
placeholder
	transport := client.GetTransport()
	field := reflect.ValueOf(transport).Elem().FieldByName("forceHttpVersion")
	require.True(t, field.IsValid(), "forceHttpVersion field not found")
	require.True(t, field.CanAddr(), "forceHttpVersion field not addressable")
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().String()
placeholder

func TestGetSharedReqClient_ForceHTTP2SeparatesCache(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	base := reqClientOptions{
		ProxyURL: "http://proxy.local:8080",
		Timeout:  time.Second,
placeholder
	clientDefault := getSharedReqClient(base)

	force := base
	force.ForceHTTP2 = true
	clientForce := getSharedReqClient(force)

	require.NotSame(t, clientDefault, clientForce)
	require.NotEqual(t, buildReqClientKey(base), buildReqClientKey(force))
placeholder

func TestGetSharedReqClient_ReuseCachedClient(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	opts := reqClientOptions{
		ProxyURL: "http://proxy.local:8080",
		Timeout:  2 * time.Second,
placeholder
	first := getSharedReqClient(opts)
	second := getSharedReqClient(opts)
	require.Same(t, first, second)
placeholder

func TestGetSharedReqClient_IgnoresNonClientCache(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	opts := reqClientOptions{
		ProxyURL: " http://proxy.local:8080 ",
		Timeout:  3 * time.Second,
placeholder
	key := buildReqClientKey(opts)
	sharedReqClients.Store(key, "invalid")

	client := getSharedReqClient(opts)

	require.NotNil(t, client)
	loaded, ok := sharedReqClients.Load(key)
	require.True(t, ok)
	require.IsType(t, "invalid", loaded)
placeholder

func TestGetSharedReqClient_ImpersonateAndProxy(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	opts := reqClientOptions{
		ProxyURL:    "  http://proxy.local:8080  ",
		Timeout:     4 * time.Second,
		Impersonate: true,
placeholder
	client := getSharedReqClient(opts)

	require.NotNil(t, client)
	require.Equal(t, "http://proxy.local:8080|4s|true|false", buildReqClientKey(opts))
placeholder

func TestCreateOpenAIReqClient_ForceHTTP2Enabled(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	client := createOpenAIReqClient("https://auth.openai.com/oauth/token", "http://proxy.local:8080")
	require.Equal(t, "2", forceHTTPVersion(t, client))
placeholder

func TestCreateOpenAIReqClient_ForceHTTP2DisabledForHTTP(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	client := createOpenAIReqClient("http://localhost/oauth/token", "http://proxy.local:8080")
	require.Equal(t, "", forceHTTPVersion(t, client))
placeholder

func TestCreateOpenAIReqClient_Timeout120Seconds(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	client := createOpenAIReqClient("https://auth.openai.com/oauth/token", "http://proxy.local:8080")
	require.Equal(t, 120*time.Second, client.GetClient().Timeout)
placeholder

func TestCreateGeminiReqClient_ForceHTTP2Disabled(t *testing.T) {
	sharedReqClients = sync.Map{placeholder
	client := createGeminiReqClient("http://proxy.local:8080")
	require.Equal(t, "", forceHTTPVersion(t, client))
placeholder
