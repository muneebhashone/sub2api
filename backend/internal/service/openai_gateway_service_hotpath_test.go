package service

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestOpenAIRequestView_ExtractsRawScalars(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":" gpt-5 ","stream":true,"prompt_cache_key":" ses-1 ","previous_response_id":" resp-1 ","service_tier":" fast ","reasoning":{"effort":" medium "placeholderplaceholder`))

	require.Equal(t, "gpt-5", view.Model)
	require.True(t, view.Stream)
	require.Equal(t, "ses-1", view.PromptCacheKey)
	require.Equal(t, "resp-1", view.PreviousResponseID)
	require.Equal(t, "fast", view.ServiceTier)
	require.Equal(t, "medium", view.ReasoningEffort)
placeholder

func TestOpenAIRequestView_DecodeKeepsFullMapBehavior(t *testing.T) {
	view := newOpenAIRequestView([]byte(`{"model":"gpt-5","stream":true,"input":[{"type":"message","content":"hi"placeholder]placeholder`))

	reqBody, err := view.Decode(nil)
placeholder
	require.Equal(t, "gpt-5", reqBody["model"])
	require.IsType(t, []any{placeholder, reqBody["input"])
placeholder

func TestExtractOpenAIRequestMetaFromBody(t *testing.T) {
	tests := []struct {
		name          string
		body          []byte
		wantModel     string
		wantStream    bool
		wantPromptKey string
placeholder{
		{
			name:          "完整字段",
			body:          []byte(`{"model":"gpt-5","stream":true,"prompt_cache_key":" ses-1 "placeholder`),
			wantModel:     "gpt-5",
			wantStream:    true,
			wantPromptKey: "ses-1",
	placeholder,
		{
			name:          "缺失可选字段",
			body:          []byte(`{"model":"gpt-4"placeholder`),
			wantModel:     "gpt-4",
			wantStream:    false,
			wantPromptKey: "",
	placeholder,
		{
			name:          "空请求体",
			body:          nil,
			wantModel:     "",
			wantStream:    false,
			wantPromptKey: "",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, stream, promptKey := extractOpenAIRequestMetaFromBody(tt.body)
			require.Equal(t, tt.wantModel, model)
			require.Equal(t, tt.wantStream, stream)
			require.Equal(t, tt.wantPromptKey, promptKey)
	placeholder)
placeholder
placeholder

func TestExtractOpenAIReasoningEffortFromBody(t *testing.T) {
	tests := []struct {
		name      string
		body      []byte
		model     string
		wantNil   bool
		wantValue string
placeholder{
		{
			name:      "优先读取 reasoning.effort",
			body:      []byte(`{"reasoning":{"effort":"medium"placeholderplaceholder`),
			model:     "gpt-5-high",
			wantNil:   false,
			wantValue: "medium",
	placeholder,
		{
			name:      "兼容 reasoning_effort",
			body:      []byte(`{"reasoning_effort":"x-high"placeholder`),
			model:     "",
			wantNil:   false,
			wantValue: "xhigh",
	placeholder,
		{
			name:    "minimal 归一化为空",
			body:    []byte(`{"reasoning":{"effort":"minimal"placeholderplaceholder`),
			model:   "gpt-5-high",
			wantNil: true,
	placeholder,
		{
			name:      "缺失字段时从模型后缀推导",
			body:      []byte(`{"input":"hi"placeholder`),
			model:     "gpt-5-high",
			wantNil:   false,
			wantValue: "high",
	placeholder,
		{
			name:    "未知后缀不返回",
			body:    []byte(`{"input":"hi"placeholder`),
			model:   "gpt-5-unknown",
			wantNil: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractOpenAIReasoningEffortFromBody(tt.body, tt.model)
			if tt.wantNil {
				require.Nil(t, got)
				return
		placeholder
			require.NotNil(t, got)
			require.Equal(t, tt.wantValue, *got)
	placeholder)
placeholder
placeholder

func TestGetOpenAIRequestBodyMap_ParseError(t *testing.T) {
	_, err := getOpenAIRequestBodyMap(nil, []byte(`{invalid-json`))
placeholder
	require.Contains(t, err.Error(), "parse request")
placeholder

func TestGetOpenAIRequestBodyMap_DoesNotWriteContextCache(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	got, err := getOpenAIRequestBodyMap(c, []byte(`{"model":"gpt-5","stream":trueplaceholder`))
placeholder
	require.Equal(t, "gpt-5", got["model"])
	require.Empty(t, c.Keys)
placeholder

func TestSanitizeEmptyBase64InputImagesInOpenAIRequestBodyMap(t *testing.T) {
	var reqBody map[string]any
	require.NoError(t, json.Unmarshal([]byte(`{
		"model":"gpt-5.4",
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,   "placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,abc123"placeholder
			]placeholder,
			{"role":"user","content":[
				{"type":"input_image","image_url":"data:image/png;base64,"placeholder
			]placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,"placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,top-level-valid"placeholder
		]
placeholder`), &reqBody))

	require.True(t, sanitizeEmptyBase64InputImagesInOpenAIRequestBodyMap(reqBody))

	normalized, err := json.Marshal(reqBody)
placeholder
	require.JSONEq(t, `{
		"model":"gpt-5.4",
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,abc123"placeholder
			]placeholder,
			{"type":"input_image","image_url":"data:image/png;base64,top-level-valid"placeholder
		]
placeholder`, string(normalized))
placeholder

func TestSanitizeEmptyBase64InputImagesInOpenAIBody(t *testing.T) {
	body, changed, err := sanitizeEmptyBase64InputImagesInOpenAIBody([]byte(`{
		"model":"gpt-5.4",
		"stream":true,
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder,
				{"type":"input_image","image_url":"data:image/png;base64,"placeholder
			]placeholder
		]
placeholder`))
placeholder
	require.True(t, changed)
	require.JSONEq(t, `{
		"model":"gpt-5.4",
		"stream":true,
		"input":[
			{"role":"user","content":[
				{"type":"input_text","text":"Describe this"placeholder
			]placeholder
		]
placeholder`, string(body))
placeholder
