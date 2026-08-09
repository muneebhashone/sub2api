//go:build unit

package xai

import (
	"context"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ssoDeviceFakeClient struct {
	t             *testing.T
	tokenCalls    int
	cookieHeaders []string
placeholder

func (c *ssoDeviceFakeClient) Do(req *http.Request) (*http.Response, error) {
	c.cookieHeaders = append(c.cookieHeaders, req.Header.Get("Cookie"))
	switch req.URL.String() {
	case SSOAccountsURL:
		require.Equal(c.t, http.MethodGet, req.Method)
		return ssoDeviceResponse(http.StatusOK, http.Header{"Set-Cookie": {"session=web-session; Domain=x.ai; Path=/"placeholderplaceholder, `{placeholder`), nil
	case SSODeviceURL:
		require.Equal(c.t, http.MethodPost, req.Method)
		values := readSSODeviceForm(c.t, req)
		require.Equal(c.t, DefaultClientID, values.Get("client_id"))
		require.Equal(c.t, SSOBuildScope, values.Get("scope"))
		return ssoDeviceResponse(http.StatusOK, http.Header{"Set-Cookie": {"csrf=csrf-token; Path=/"placeholderplaceholder, `{"device_code":"device-1","user_code":"USER-1","verification_uri_complete":"https://auth.x.ai/oauth2/device/complete","interval":1,"expires_in":60placeholder`), nil
	case "https://auth.x.ai/oauth2/device/complete":
		require.Equal(c.t, http.MethodGet, req.Method)
		return ssoDeviceResponse(http.StatusOK, nil, `<html>ok</html>`), nil
	case SSOVerifyURL:
		require.Equal(c.t, http.MethodPost, req.Method)
		values := readSSODeviceForm(c.t, req)
		require.Equal(c.t, "USER-1", values.Get("user_code"))
		return ssoDeviceResponse(http.StatusFound, http.Header{"Location": {"/oauth2/device/consent"placeholderplaceholder, ``), nil
	case "https://auth.x.ai/oauth2/device/consent":
		require.Equal(c.t, http.MethodGet, req.Method)
		return ssoDeviceResponse(http.StatusOK, nil, `<html>consent</html>`), nil
	case SSOApproveURL:
		require.Equal(c.t, http.MethodPost, req.Method)
		values := readSSODeviceForm(c.t, req)
		require.Equal(c.t, "USER-1", values.Get("user_code"))
		require.Equal(c.t, "allow", values.Get("action"))
		require.Equal(c.t, "User", values.Get("principal_type"))
		return ssoDeviceResponse(http.StatusSeeOther, http.Header{"Location": {"/oauth2/device/done"placeholderplaceholder, ``), nil
	case "https://auth.x.ai/oauth2/device/done":
		require.Equal(c.t, http.MethodGet, req.Method)
		return ssoDeviceResponse(http.StatusOK, nil, `<html>done</html>`), nil
	case SSOTokenURL:
		require.Equal(c.t, http.MethodPost, req.Method)
		c.tokenCalls++
		values := readSSODeviceForm(c.t, req)
		require.Equal(c.t, "urn:ietf:params:oauth:grant-type:device_code", values.Get("grant_type"))
		require.Equal(c.t, "device-1", values.Get("device_code"))
		return ssoDeviceResponse(http.StatusOK, nil, `{"access_token":"access-token","refresh_token":"refresh-token","id_token":"id-token","token_type":"Bearer","expires_in":3600,"scope":"`+SSOBuildScope+`"placeholder`), nil
	default:
		c.t.Fatalf("unexpected request: %s %s", req.Method, req.URL.String())
		return nil, nil
placeholder
placeholder

func TestConvertSSOToBuildCompletesDeviceFlow(t *testing.T) {
	t.Setenv(EnvClientID, "")
	client := &ssoDeviceFakeClient{t: tplaceholder
	token, err := ConvertSSOToBuild(context.Background(), "sso=sso-token; ignored=1", &SSODeviceOptions{
		HTTPClient: client,
		Sleep: func(context.Context, time.Duration) error {
			return nil
	placeholder,
placeholder)

placeholder
	require.Equal(t, "access-token", token.AccessToken)
	require.Equal(t, "refresh-token", token.RefreshToken)
	require.Equal(t, "id-token", token.IDToken)
	require.Equal(t, SSOBuildScope, token.Scope)
	require.Equal(t, 1, client.tokenCalls)
	require.Contains(t, client.cookieHeaders[0], "sso=sso-token")
	require.Contains(t, client.cookieHeaders[0], "sso-rw=sso-token")
	require.Contains(t, client.cookieHeaders[len(client.cookieHeaders)-1], "session=web-session")
	require.Contains(t, client.cookieHeaders[len(client.cookieHeaders)-1], "csrf=csrf-token")
placeholder

func TestNormalizeSSOTokenAcceptsCookieHeader(t *testing.T) {
	require.Equal(t, "token-1", NormalizeSSOToken("Cookie: foo=bar; sso=token-1; sso-rw=token-2"))
	require.Equal(t, "token-2", NormalizeSSOToken("sso-rw=token-2; foo=bar"))
	require.Equal(t, "raw-token", NormalizeSSOToken(" raw-token ; ignored=1"))
	require.Empty(t, NormalizeSSOToken(strings.Repeat("x", ssoMaxTokenLength+1)))
placeholder

func TestSSODeviceCookieJarHonorsDomainAndPath(t *testing.T) {
	jar, err := cookiejar.New(nil)
placeholder
	flow := &ssoDeviceFlow{cookieJar: jarplaceholder
	accountsURL, err := url.Parse("https://accounts.x.ai/")
placeholder
	authURL, err := url.Parse("https://auth.x.ai/oauth2/device/verify")
placeholder

	flow.captureCookies(accountsURL, ssoDeviceResponse(http.StatusOK, http.Header{"Set-Cookie": {
		"host-only=accounts; Path=/",
		"shared=all-xai; Domain=x.ai; Path=/",
		"narrow=oauth-only; Domain=x.ai; Path=/oauth2",
placeholderplaceholder, ""))

	authCookies := flow.cookieHeader(authURL)
	require.NotContains(t, authCookies, "host-only=accounts")
	require.Contains(t, authCookies, "shared=all-xai")
	require.Contains(t, authCookies, "narrow=oauth-only")

	accountsCookies := flow.cookieHeader(accountsURL)
	require.Contains(t, accountsCookies, "host-only=accounts")
	require.Contains(t, accountsCookies, "shared=all-xai")
	require.NotContains(t, accountsCookies, "narrow=oauth-only")
placeholder

func ssoDeviceResponse(status int, header http.Header, body string) *http.Response {
	if header == nil {
		header = http.Header{placeholder
placeholder
	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
placeholder
placeholder

func readSSODeviceForm(t *testing.T, req *http.Request) url.Values {
placeholder
	data, err := io.ReadAll(req.Body)
placeholder
	values, err := url.ParseQuery(string(data))
placeholder
	return values
placeholder
