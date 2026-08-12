package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOpenAIVisibleOutputClassification(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		eventType string
		want      bool
placeholder{
		{name: "keepalive", data: `{"type":"keepalive"placeholder`, want: falseplaceholder,
		{name: "created", data: `{"type":"response.created"placeholder`, want: falseplaceholder,
		{name: "empty output item", data: `{"type":"response.output_item.added","item":{"id":"item_test","type":"reasoning","summary":[]placeholderplaceholder`, want: falseplaceholder,
		{name: "empty delta", data: `{"type":"response.output_text.delta","delta":""placeholder`, want: falseplaceholder,
		{name: "text delta", data: `{"type":"response.output_text.delta","delta":"test output"placeholder`, want: trueplaceholder,
		{name: "tool arguments", data: `{"type":"response.function_call_arguments.delta","delta":"{placeholder"placeholder`, want: trueplaceholder,
		{name: "partial image", data: `{"type":"response.image_generation_call.partial_image","partial_image_b64":"dGVzdA=="placeholder`, want: trueplaceholder,
		{name: "completed image item", data: `{"type":"response.output_item.done","item":{"id":"item_test","type":"image_generation_call","result":"dGVzdA=="placeholderplaceholder`, want: trueplaceholder,
		{name: "empty completed", data: `{"type":"response.completed","response":{"id":"resp_test","output":[]placeholderplaceholder`, want: falseplaceholder,
		{name: "completed with output usage only", data: `{"type":"response.completed","response":{"id":"resp_test","usage":{"input_tokens":1,"output_tokens":2placeholderplaceholderplaceholder`, want: falseplaceholder,
		{name: "completed with text", data: `{"type":"response.completed","response":{"id":"resp_test","output":[{"type":"message","content":[{"type":"output_text","text":"test output"placeholder]placeholder]placeholderplaceholder`, want: trueplaceholder,
		{name: "done marker", data: `[DONE]`, want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, openAIStreamDataStartsVisibleOutput(tt.data, tt.eventType))
	placeholder)
placeholder
placeholder

func TestOpenAIResponsesTTFTStartsAtVisibleOutput(t *testing.T) {
	for _, passthrough := range []bool{false, trueplaceholder {
		name := "native"
		if passthrough {
			name = "passthrough"
	placeholder
		t.Run(name, func(t *testing.T) {
			result := runSyntheticVisibleTTFTStream(t, passthrough, 120*time.Millisecond, 0,
				`{"type":"response.output_text.delta","delta":"test output"placeholder`)
			require.NotNil(t, result.firstTokenMs)
			require.GreaterOrEqual(t, *result.firstTokenMs, 100)
	placeholder)
placeholder
placeholder

func TestOpenAIResponsesTTFTStartsAtCompletedImage(t *testing.T) {
	for _, passthrough := range []bool{false, trueplaceholder {
		name := "native"
		if passthrough {
			name = "passthrough"
	placeholder
		t.Run(name, func(t *testing.T) {
			result := runSyntheticVisibleTTFTStream(t, passthrough, 120*time.Millisecond, 0,
				`{"type":"response.output_item.done","item":{"id":"item_test","type":"image_generation_call","result":"dGVzdA=="placeholderplaceholder`)
			require.NotNil(t, result.firstTokenMs)
			require.GreaterOrEqual(t, *result.firstTokenMs, 100)
	placeholder)
placeholder
placeholder

func TestOpenAINativeProgressDisarmsTimeoutWithoutStartingTTFT(t *testing.T) {
	result := runSyntheticVisibleTTFTStream(t, false, 1200*time.Millisecond, 1,
		`{"type":"response.output_text.delta","delta":"test output"placeholder`)
	require.NotNil(t, result.firstTokenMs)
	require.GreaterOrEqual(t, *result.firstTokenMs, 1100)
placeholder

func runSyntheticVisibleTTFTStream(t *testing.T, passthrough bool, visibleDelay time.Duration, timeoutSeconds int, visibleEvent string) *openaiStreamingResult {
placeholder
	gin.SetMode(gin.TestMode)
	svc := &OpenAIGatewayService{cfg: &config.Config{Gateway: config.GatewayConfig{
		MaxLineSize:                     defaultMaxLineSize,
		OpenAIFirstOutputTimeoutSeconds: timeoutSeconds,
placeholderplaceholderplaceholder
	reader, writer := io.Pipe()
	writerDone := make(chan struct{placeholder)
	go func() {
		defer close(writerDone)
		defer func() { _ = writer.Close() placeholder()
		_, _ = io.WriteString(writer, "data: {\"type\":\"response.created\",\"response\":{\"id\":\"resp_test\"placeholderplaceholder\n\n")
		_, _ = io.WriteString(writer, "data: {\"type\":\"response.output_item.added\",\"item\":{\"id\":\"item_test\",\"type\":\"reasoning\",\"summary\":[]placeholderplaceholder\n\n")
		time.Sleep(visibleDelay)
		_, _ = io.WriteString(writer, "data: "+visibleEvent+"\n\n")
		_, _ = io.WriteString(writer, "data: {\"type\":\"response.completed\",\"response\":{\"id\":\"resp_test\",\"usage\":{\"input_tokens\":1,\"output_tokens\":1placeholderplaceholderplaceholder\n\n")
placeholder()

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/responses", nil)
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{placeholder, Body: readerplaceholder
	account := &Account{ID: 1, Name: "account_test", Platform: PlatformOpenAIplaceholder
	started := time.Now()

	var result *openaiStreamingResult
	var err error
	if passthrough {
		var passthroughResult *openaiStreamingResultPassthrough
		passthroughResult, err = svc.handleStreamingResponsePassthrough(context.Background(), resp, c, account, started, "test-model", "test-model")
		if passthroughResult != nil {
			result = &openaiStreamingResult{firstTokenMs: passthroughResult.firstTokenMsplaceholder
	placeholder
placeholder else {
		result, err = svc.handleStreamingResponse(context.Background(), resp, c, account, started, "test-model", "test-model")
placeholder
placeholder
	require.NotNil(t, result)
	require.Contains(t, recorder.Body.String(), `"type":"response.output_item.added"`)
	require.Contains(t, recorder.Body.String(), visibleEvent)
	select {
	case <-writerDone:
	case <-time.After(time.Second):
		t.Fatal("synthetic upstream writer did not exit")
placeholder
	return result
placeholder
