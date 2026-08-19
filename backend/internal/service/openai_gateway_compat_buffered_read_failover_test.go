package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type openAICompatBufferedReadErrorCloser struct {
	err error
placeholder

func (r *openAICompatBufferedReadErrorCloser) Read([]byte) (int, error) { return 0, r.err placeholder
func (r *openAICompatBufferedReadErrorCloser) Close() error             { return nil placeholder

func TestChatCompletionsBufferedResponsesReadErrorReturnsFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	readErrors := []struct {
		name         string
		err          error
		expectedCode string
placeholder{
		{name: "unexpected_eof", err: io.ErrUnexpectedEOF, expectedCode: OpenAIUpstreamStreamReadErrorCodeplaceholder,
		{name: "http2_reset", err: errors.New("stream error: stream ID 7; INTERNAL_ERROR; received from peer"), expectedCode: OpenAIUpstreamHTTP2StreamErrorCodeplaceholder,
placeholder

	for _, readError := range readErrors {
		t.Run(readError.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholder, "X-Request-Id": []string{"upstream-rid"placeholderplaceholder,
				Body:       &openAICompatBufferedReadErrorCloser{err: readError.errplaceholder,
		placeholder
			account := &Account{ID: 40, Name: "openai-oauth", Platform: PlatformOpenAIplaceholder

			result, err := (&OpenAIGatewayService{placeholder).handleChatBufferedStreamingResponse(
				resp, c, account, "gpt-5.6-sol", "gpt-5.6-sol", "gpt-5.6-sol", time.Now(),
			)

		placeholder
			require.Nil(t, result)
			var failoverErr *UpstreamFailoverError
			require.ErrorAs(t, err, &failoverErr)
			require.Equal(t, http.StatusBadGateway, failoverErr.StatusCode)
			require.Equal(t, "upstream-rid", failoverErr.ResponseHeaders.Get("x-request-id"))
			require.Equal(t, readError.expectedCode, gjson.GetBytes(failoverErr.ResponseBody, "error.code").String())
			require.Empty(t, rec.Body.String())
			require.False(t, c.Writer.Written())
	placeholder)
placeholder
placeholder

func TestChatCompletionsBufferedResponsesReadErrorDoesNotFailoverAfterClientCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	requestContext, cancel := context.WithCancel(context.Background())
	cancel()
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil).WithContext(requestContext)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       &openAICompatBufferedReadErrorCloser{err: io.ErrUnexpectedEOFplaceholder,
placeholder

	result, err := (&OpenAIGatewayService{placeholder).handleChatBufferedStreamingResponse(
		resp,
		c,
		&Account{ID: 40, Name: "openai-oauth", Platform: PlatformOpenAIplaceholder,
		"gpt-5.6-sol",
		"gpt-5.6-sol",
		"gpt-5.6-sol",
		time.Now(),
	)

placeholder
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.NotErrorAs(t, err, &failoverErr)
	require.Empty(t, rec.Body.String())
placeholder

func TestChatCompletionsBufferedResponsesOversizedLineDoesNotFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       &openAICompatBufferedReadErrorCloser{err: bufio.ErrTooLongplaceholder,
placeholder

	result, err := (&OpenAIGatewayService{placeholder).handleChatBufferedStreamingResponse(
		resp,
		c,
		&Account{ID: 40, Name: "openai-oauth", Platform: PlatformOpenAIplaceholder,
		"gpt-5.6-sol",
		"gpt-5.6-sol",
		"gpt-5.6-sol",
		time.Now(),
	)

	require.ErrorIs(t, err, bufio.ErrTooLong)
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.NotErrorAs(t, err, &failoverErr)
placeholder

func TestAnthropicBufferedResponsesReadErrorKeepsExistingBehavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"text/event-stream"placeholderplaceholder,
		Body:       &openAICompatBufferedReadErrorCloser{err: io.ErrUnexpectedEOFplaceholder,
placeholder

	result, err := (&OpenAIGatewayService{placeholder).handleAnthropicBufferedStreamingResponse(
		resp,
		c,
		&Account{ID: 40, Name: "openai-oauth", Platform: PlatformOpenAIplaceholder,
		"gpt-5.6-sol",
		"gpt-5.6-sol",
		"gpt-5.6-sol",
		time.Now(),
	)

	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
	require.Equal(t, io.ErrUnexpectedEOF, err, "Messages 路径必须保持原始读取错误，不引入 Chat failover 包装")
	require.Nil(t, result)
	var failoverErr *UpstreamFailoverError
	require.NotErrorAs(t, err, &failoverErr)
placeholder
