package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestOpenAICompatibleHandlersRejectInvalidStreamFieldType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		path string
		body string
		run  func(*gin.Context)
placeholder{
		{
			name: "gateway_responses_string_stream",
			path: "/v1/responses",
			body: `{"model":"gpt-5","stream":"true","input":"hello"placeholder`,
			run: func(c *gin.Context) {
				(&GatewayHandler{placeholder).Responses(c)
		placeholder,
	placeholder,
		{
			name: "gateway_responses_number_stream",
			path: "/v1/responses",
			body: `{"model":"gpt-5","stream":1,"input":"hello"placeholder`,
			run: func(c *gin.Context) {
				(&GatewayHandler{placeholder).Responses(c)
		placeholder,
	placeholder,
		{
			name: "gateway_chat_completions_string_stream",
			path: "/v1/chat/completions",
			body: `{"model":"gpt-5","stream":"true","messages":[{"role":"user","content":"hello"placeholder]placeholder`,
			run: func(c *gin.Context) {
				(&GatewayHandler{placeholder).ChatCompletions(c)
		placeholder,
	placeholder,
		{
			name: "gateway_chat_completions_number_stream",
			path: "/v1/chat/completions",
			body: `{"model":"gpt-5","stream":1,"messages":[{"role":"user","content":"hello"placeholder]placeholder`,
			run: func(c *gin.Context) {
				(&GatewayHandler{placeholder).ChatCompletions(c)
		placeholder,
	placeholder,
		{
			name: "openai_chat_completions_string_stream",
			path: "/openai/v1/chat/completions",
			body: `{"model":"gpt-5","stream":"true","messages":[{"role":"user","content":"hello"placeholder]placeholder`,
			run: func(c *gin.Context) {
				newOpenAIHandlerForPreviousResponseIDValidation(t, nil).ChatCompletions(c)
		placeholder,
	placeholder,
		{
			name: "openai_chat_completions_number_stream",
			path: "/openai/v1/chat/completions",
			body: `{"model":"gpt-5","stream":1,"messages":[{"role":"user","content":"hello"placeholder]placeholder`,
			run: func(c *gin.Context) {
				newOpenAIHandlerForPreviousResponseIDValidation(t, nil).ChatCompletions(c)
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newOpenAICompatibleStreamValidationContext(tt.path, tt.body, false)

			tt.run(c)

			require.Equal(t, http.StatusBadRequest, rec.Code)
			require.Equal(t, invalidStreamFieldTypeMessage, gjson.GetBytes(rec.Body.Bytes(), "error.message").String())
			require.Contains(t, rec.Body.String(), "invalid_request_error")
	placeholder)
placeholder
placeholder

func TestGatewayOpenAICompatibleHandlersAllowBooleanStreamToContinue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		path string
		body string
		run  func(*gin.Context)
placeholder{
		{
			name: "responses_false",
			path: "/v1/responses",
			body: `{"model":"gpt-5","stream":false,"input":"hello"placeholder`,
			run: func(c *gin.Context) {
				(&GatewayHandler{gatewayService: &service.GatewayService{placeholderplaceholder).Responses(c)
		placeholder,
	placeholder,
		{
			name: "chat_completions_true",
			path: "/v1/chat/completions",
			body: `{"model":"gpt-5","stream":true,"messages":[{"role":"user","content":"hello"placeholder]placeholder`,
			run: func(c *gin.Context) {
				(&GatewayHandler{gatewayService: &service.GatewayService{placeholderplaceholder).ChatCompletions(c)
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := newOpenAICompatibleStreamValidationContext(tt.path, tt.body, true)

			tt.run(c)

			require.Equal(t, http.StatusForbidden, rec.Code)
			require.Contains(t, rec.Body.String(), "This group is restricted to Claude Code clients")
	placeholder)
placeholder
placeholder

func newOpenAICompatibleStreamValidationContext(path, body string, claudeCodeOnly bool) (*gin.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	groupID := int64(7)
	c.Set(string(middleware2.ContextKeyAPIKey), &service.APIKey{
		ID:      11,
		GroupID: &groupID,
		Group:   &service.Group{ID: groupID, ClaudeCodeOnly: claudeCodeOnlyplaceholder,
		User:    &service.User{ID: 13placeholder,
placeholder)
	c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 13, Concurrency: 1placeholder)

	return c, rec
placeholder
