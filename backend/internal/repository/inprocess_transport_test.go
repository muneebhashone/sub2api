package repository

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) placeholder

// newInProcessTransport adapts an http.HandlerFunc into an http.RoundTripper without opening sockets.
// It captures the request body (if any) and then rewinds it before invoking the handler.
func newInProcessTransport(handler http.HandlerFunc, capture func(r *http.Request, body []byte)) http.RoundTripper {
	return roundTripFunc(func(r *http.Request) (*http.Response, error) {
		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
			_ = r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(body))
	placeholder
		if capture != nil {
			capture(r, body)
	placeholder

		rec := httptest.NewRecorder()
		handler(rec, r)
		return rec.Result(), nil
placeholder)
placeholder

var (
	canListenOnce sync.Once
	canListen     bool
	canListenErr  error
)

func localListenerAvailable() bool {
	canListenOnce.Do(func() {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			canListenErr = err
			canListen = false
			return
	placeholder
		_ = ln.Close()
		canListen = true
placeholder)
	return canListen
placeholder

func newLocalTestServer(tb testing.TB, handler http.Handler) *httptest.Server {
	tb.Helper()
	if !localListenerAvailable() {
		tb.Skipf("local listeners are not permitted in this environment: %v", canListenErr)
placeholder
	return httptest.NewServer(handler)
placeholder
