package repository

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PricingServiceSuite struct {
	suite.Suite
	ctx    context.Context
	srv    *httptest.Server
	client *pricingRemoteClient
placeholder

func (s *PricingServiceSuite) SetupTest() {
	s.ctx = context.Background()
	client, ok := NewPricingRemoteClient(&config.Config{
		Security: config.SecurityConfig{
			URLAllowlist: config.URLAllowlistConfig{
				AllowPrivateHosts: true,
		placeholder,
	placeholder,
placeholder).(*pricingRemoteClient)
	require.True(s.T(), ok, "type assertion failed")
	s.client = client
placeholder

func (s *PricingServiceSuite) TearDownTest() {
	if s.srv != nil {
		s.srv.Close()
		s.srv = nil
placeholder
placeholder

func (s *PricingServiceSuite) setupServer(handler http.HandlerFunc) {
	s.srv = httptest.NewServer(handler)
placeholder

func (s *PricingServiceSuite) TestFetchPricingJSON_Success() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"ok":trueplaceholder`))
			return
	placeholder
		w.WriteHeader(http.StatusInternalServerError)
placeholder))

	body, err := s.client.FetchPricingJSON(s.ctx, s.srv.URL+"/ok")
	require.NoError(s.T(), err, "FetchPricingJSON")
	require.Equal(s.T(), `{"ok":trueplaceholder`, string(body), "body mismatch")
placeholder

func (s *PricingServiceSuite) TestFetchPricingJSON_NonOKStatus() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
placeholder))

	_, err := s.client.FetchPricingJSON(s.ctx, s.srv.URL+"/err")
	require.Error(s.T(), err, "expected error for non-200 status")
placeholder

func (s *PricingServiceSuite) TestFetchHashText_ParsesFields() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/hashfile":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("abc123  model_prices.json\n"))
		case "/hashonly":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("def456\n"))
		default:
			w.WriteHeader(http.StatusNotFound)
	placeholder
placeholder))

	hash, err := s.client.FetchHashText(s.ctx, s.srv.URL+"/hashfile")
	require.NoError(s.T(), err, "FetchHashText")
	require.Equal(s.T(), "abc123", hash, "hash mismatch")

	hash2, err := s.client.FetchHashText(s.ctx, s.srv.URL+"/hashonly")
	require.NoError(s.T(), err, "FetchHashText")
	require.Equal(s.T(), "def456", hash2, "hash mismatch")
placeholder

func (s *PricingServiceSuite) TestFetchHashText_NonOKStatus() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
placeholder))

	_, err := s.client.FetchHashText(s.ctx, s.srv.URL+"/nope")
	require.Error(s.T(), err, "expected error for non-200 status")
placeholder

func (s *PricingServiceSuite) TestFetchPricingJSON_InvalidURL() {
	_, err := s.client.FetchPricingJSON(s.ctx, "://invalid-url")
	require.Error(s.T(), err, "expected error for invalid URL")
placeholder

func (s *PricingServiceSuite) TestFetchHashText_EmptyBody() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// empty body
placeholder))

	hash, err := s.client.FetchHashText(s.ctx, s.srv.URL+"/empty")
	require.NoError(s.T(), err, "FetchHashText empty body should not error")
	require.Equal(s.T(), "", hash, "expected empty hash")
placeholder

func (s *PricingServiceSuite) TestFetchHashText_WhitespaceOnly() {
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("   \n"))
placeholder))

	hash, err := s.client.FetchHashText(s.ctx, s.srv.URL+"/ws")
	require.NoError(s.T(), err, "FetchHashText whitespace body should not error")
	require.Equal(s.T(), "", hash, "expected empty hash after trimming")
placeholder

func (s *PricingServiceSuite) TestFetchPricingJSON_ContextCancel() {
	started := make(chan struct{placeholder)
	s.setupServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
placeholder))

	ctx, cancel := context.WithCancel(s.ctx)

	done := make(chan error, 1)
	go func() {
		_, err := s.client.FetchPricingJSON(ctx, s.srv.URL+"/block")
		done <- err
placeholder()

	<-started
	cancel()

	err := <-done
	require.Error(s.T(), err)
placeholder

func TestPricingServiceSuite(t *testing.T) {
	suite.Run(t, new(PricingServiceSuite))
placeholder
