package handler

import (
	"net/http"
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

func decodeCookieValueForTest(t *testing.T, value string) string {
placeholder
	decoded, err := decodeCookieValue(value)
placeholder
	return decoded
placeholder

func assertOAuthRedirectError(t *testing.T, location string, errorCode string, errorMessage string) {
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
	require.Equal(t, errorCode, values.Get("error"))
	require.Equal(t, errorMessage, values.Get("error_message"))
placeholder
