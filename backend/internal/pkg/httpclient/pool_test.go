package httpclient

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
placeholder

func TestValidatedTransport_CacheHostValidation(t *testing.T) {
	originalValidate := validateResolvedIP
	defer func() { validateResolvedIP = originalValidate placeholder()

	var validateCalls int32
	validateResolvedIP = func(host string) error {
		atomic.AddInt32(&validateCalls, 1)
		require.Equal(t, "api.openai.com", host)
		return nil
placeholder

	var baseCalls int32
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		atomic.AddInt32(&baseCalls, 1)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{placeholder`)),
			Header:     make(http.Header),
	placeholder, nil
placeholder)

	now := time.Unix(1730000000, 0)
	transport := newValidatedTransport(base)
	transport.now = func() time.Time { return now placeholder

	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/responses", nil)
placeholder

	_, err = transport.RoundTrip(req)
placeholder
	_, err = transport.RoundTrip(req)
placeholder

	require.Equal(t, int32(1), atomic.LoadInt32(&validateCalls))
	require.Equal(t, int32(2), atomic.LoadInt32(&baseCalls))
placeholder

func TestValidatedTransport_ExpiredCacheTriggersRevalidation(t *testing.T) {
	originalValidate := validateResolvedIP
	defer func() { validateResolvedIP = originalValidate placeholder()

	var validateCalls int32
	validateResolvedIP = func(_ string) error {
		atomic.AddInt32(&validateCalls, 1)
		return nil
placeholder

	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{placeholder`)),
			Header:     make(http.Header),
	placeholder, nil
placeholder)

	now := time.Unix(1730001000, 0)
	transport := newValidatedTransport(base)
	transport.now = func() time.Time { return now placeholder

	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/responses", nil)
placeholder

	_, err = transport.RoundTrip(req)
placeholder

	now = now.Add(validatedHostTTL + time.Second)
	_, err = transport.RoundTrip(req)
placeholder

	require.Equal(t, int32(2), atomic.LoadInt32(&validateCalls))
placeholder

func TestValidatedTransport_ValidationErrorStopsRoundTrip(t *testing.T) {
	originalValidate := validateResolvedIP
	defer func() { validateResolvedIP = originalValidate placeholder()

	expectedErr := errors.New("dns rebinding rejected")
	validateResolvedIP = func(_ string) error {
		return expectedErr
placeholder

	var baseCalls int32
	base := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		atomic.AddInt32(&baseCalls, 1)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{placeholder`))placeholder, nil
placeholder)

	transport := newValidatedTransport(base)
	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/responses", nil)
placeholder

	_, err = transport.RoundTrip(req)
	require.ErrorIs(t, err, expectedErr)
	require.Equal(t, int32(0), atomic.LoadInt32(&baseCalls))
placeholder
