package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildEncodedOAuthBindUserCookie(t *testing.T, userID int64, secret string) string {
placeholder
	value, err := buildOAuthBindUserCookieValue(userID, secret)
placeholder
	return value
placeholder

func encodedCookie(name, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: encodeCookieValue(value),
		Path:  "/",
placeholder
placeholder

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
	placeholder
placeholder
	return nil
placeholder

func requireCookieCleared(t *testing.T, recorder *httptest.ResponseRecorder, name string) {
placeholder
	cookie := findCookie(recorder.Result().Cookies(), name)
	require.NotNil(t, cookie)
	require.Equal(t, -1, cookie.MaxAge)
placeholder

func decodeCookieValueForTest(t *testing.T, value string) string {
placeholder
	decoded, err := decodeCookieValue(value)
placeholder
	return decoded
placeholder

func assertOAuthRedirectError(t *testing.T, location string, errorCode string, errorMessage string) {
placeholder
	values := parseOAuthRedirectFragment(t, location)
	require.Equal(t, errorCode, values.Get("error"))
	require.Equal(t, errorMessage, values.Get("error_message"))
placeholder

func parseOAuthRedirectFragment(t *testing.T, location string) url.Values {
placeholder
	require.NotEmpty(t, location)

	parsed, err := url.Parse(location)
placeholder

	rawValues := parsed.RawQuery
	if rawValues == "" {
		rawValues = parsed.Fragment
placeholder
	values, err := url.ParseQuery(rawValues)
placeholder
	return values
placeholder
