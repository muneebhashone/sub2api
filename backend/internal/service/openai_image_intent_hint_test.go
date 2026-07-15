package service

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func newOpenAIImageIntentHintTestContext(transport OpenAIClientTransport) *gin.Context {
	c := &gin.Context{placeholder
	SetOpenAIClientTransport(c, transport)
	return c
placeholder

func countingOpenAIImageIntentClassifier(calls *atomic.Int64) openAIImageIntentClassifier {
	return func(endpoint string, requestedModel string, body []byte) bool {
		calls.Add(1)
		return IsImageGenerationIntent(endpoint, requestedModel, body)
placeholder
placeholder

func TestResolveOpenAIImageIntentHintCachesTrueAndFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		body []byte
		want bool
placeholder{
		{name: "true", body: []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation"placeholder]placeholder`), want: trueplaceholder,
		{name: "false is known", body: []byte(`{"model":"gpt-5.4","input":"write code"placeholder`), want: falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
			var calls atomic.Int64
			classify := countingOpenAIImageIntentClassifier(&calls)

			require.Equal(t, tt.want, resolveOpenAIImageIntentHint(c, "gpt-5.4", tt.body, classify))
			require.Equal(t, tt.want, resolveOpenAIImageIntentHint(c, "gpt-5.4", tt.body, classify))
			require.Equal(t, int64(1), calls.Load())
			cached, known := getOpenAIImageIntentHint(c)
			require.True(t, known)
			require.Equal(t, tt.want, cached)
	placeholder)
placeholder
placeholder

func TestResolveOpenAIImageIntentHintUsesHandlerSeed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, seeded := range []bool{false, trueplaceholder {
		c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
		SetOpenAIImageIntentHint(c, seeded)
		var calls atomic.Int64

		got := resolveOpenAIImageIntentHint(c, "gpt-5.4", []byte(`{"model":"gpt-5.4"placeholder`), countingOpenAIImageIntentClassifier(&calls))

		require.Equal(t, seeded, got)
		require.Zero(t, calls.Load())
placeholder
placeholder

func TestResolveOpenAIPassthroughImageIntentReusesCanonicalAcrossFailover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
	body := []byte(`{"model":"gpt-5.4","input":"write code"placeholder`)
	var calls atomic.Int64
	classify := countingOpenAIImageIntentClassifier(&calls)

	for range 3 {
		require.False(t, resolveOpenAIPassthroughImageIntent(c, "gpt-5.4", body, "gpt-5.4", body, false, classify))
placeholder
	require.Equal(t, int64(1), calls.Load())
placeholder

func TestResolveOpenAIPassthroughImageIntentKeepsCompactMappingAttemptLocal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("text to image", func(t *testing.T) {
		c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
		body := []byte(`{"model":"draw-alias","input":"draw"placeholder`)
		compactBody := []byte(`{"model":"gpt-image-2","input":"draw"placeholder`)
		var calls atomic.Int64
		classify := countingOpenAIImageIntentClassifier(&calls)

		require.True(t, resolveOpenAIPassthroughImageIntent(c, "draw-alias", body, "gpt-image-2", compactBody, true, classify))
		cached, known := getOpenAIImageIntentHint(c)
		require.True(t, known)
		require.False(t, cached)

		require.False(t, resolveOpenAIPassthroughImageIntent(c, "draw-alias", body, "draw-alias", body, false, classify))
		require.Equal(t, int64(2), calls.Load())
		cached, known = getOpenAIImageIntentHint(c)
		require.True(t, known)
		require.False(t, cached)
placeholder)

	t.Run("image to text", func(t *testing.T) {
		c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
		body := []byte(`{"model":"gpt-image-2","input":"draw"placeholder`)
		compactBody := []byte(`{"model":"gpt-5.4","input":"draw"placeholder`)
		var calls atomic.Int64
		classify := countingOpenAIImageIntentClassifier(&calls)

		require.False(t, resolveOpenAIPassthroughImageIntent(c, "gpt-image-2", body, "gpt-5.4", compactBody, true, classify))
		cached, known := getOpenAIImageIntentHint(c)
		require.True(t, known)
		require.True(t, cached)

		require.True(t, resolveOpenAIPassthroughImageIntent(c, "gpt-image-2", body, "gpt-image-2", body, false, classify))
		require.Equal(t, int64(2), calls.Load())
placeholder)
placeholder

func TestResolveOpenAIPassthroughImageIntentInvalidationDoesNotPolluteCanonical(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
	canonicalBody := []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation"placeholder]placeholder`)
	strippedBody := []byte(`{"model":"gpt-5.4","tools":[]placeholder`)
	var calls atomic.Int64
	classify := countingOpenAIImageIntentClassifier(&calls)

	require.False(t, resolveOpenAIPassthroughImageIntent(c, "gpt-5.4", canonicalBody, "gpt-5.4", strippedBody, true, classify))
	require.Equal(t, int64(2), calls.Load(), "unknown canonical and invalidated attempt are classified independently")
	cached, known := getOpenAIImageIntentHint(c)
	require.True(t, known)
	require.True(t, cached)

	require.True(t, resolveOpenAIPassthroughImageIntent(c, "gpt-5.4", canonicalBody, "gpt-5.4", canonicalBody, false, classify))
	require.Equal(t, int64(2), calls.Load())
placeholder

func TestResolveOpenAIPassthroughImageIntentMappedBodyStartsUnknownThenSeedsCanonical(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
	canonicalBody := []byte(`{"model":"gpt-image-2","input":"draw"placeholder`)
	strippedAttemptBody := []byte(`{"model":"gpt-5.4","input":"draw"placeholder`)
	_, known := getOpenAIImageIntentHint(c)
	require.False(t, known)
	var calls atomic.Int64
	classify := countingOpenAIImageIntentClassifier(&calls)

	require.False(t, resolveOpenAIPassthroughImageIntent(c, "gpt-image-2", canonicalBody, "gpt-5.4", strippedAttemptBody, true, classify))
	require.Equal(t, int64(2), calls.Load())
	cached, known := getOpenAIImageIntentHint(c)
	require.True(t, known)
	require.True(t, cached)
placeholder

func TestResolveOpenAIPassthroughImageIntentReusesAcrossInvariantMutations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name          string
		canonicalBody []byte
		attemptBody   []byte
		want          bool
placeholder{
		{
			name:          "oauth sanitize fast policy and reasoning",
			canonicalBody: []byte(`{"model":"gpt-5.4","input":[{"type":"input_image","image_url":"data:image/png;base64,"placeholder],"service_tier":"fast","reasoning":{"effort":"minimal"placeholderplaceholder`),
			attemptBody:   []byte(`{"model":"gpt-5.4","input":[],"service_tier":"priority","reasoning":{"effort":"none"placeholder,"store":false,"stream":trueplaceholder`),
			want:          false,
	placeholder,
		{
			name:          "namespace flatten",
			canonicalBody: []byte(`{"model":"gpt-5.4","tools":[{"type":"namespace","name":"code_tools"placeholder]placeholder`),
			attemptBody:   []byte(`{"model":"gpt-5.4","tools":[{"type":"function","name":"code_tools.run"placeholder]placeholder`),
			want:          false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
			var calls atomic.Int64
			classify := countingOpenAIImageIntentClassifier(&calls)

			require.Equal(t, tt.want, resolveOpenAIPassthroughImageIntent(c, "gpt-5.4", tt.canonicalBody, "gpt-5.4", tt.attemptBody, false, classify))
			require.Equal(t, tt.want, resolveOpenAIPassthroughImageIntent(c, "gpt-5.4", tt.canonicalBody, "gpt-5.4", tt.attemptBody, false, classify))
			require.Equal(t, int64(1), calls.Load())
	placeholder)
placeholder
placeholder

func TestOpenAIGatewayServicePassthroughCompactImageIntentIsAttemptLocal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		canonicalModel string
		compactModel   string
		wantRejected   bool
		wantCanonical  bool
placeholder{
		{
			name:           "text to image rejects",
			canonicalModel: "gpt-5.4",
			compactModel:   "gpt-image-2",
			wantRejected:   true,
	placeholder,
		{
			name:           "image to text reaches upstream",
			canonicalModel: "gpt-image-2",
			compactModel:   "gpt-5.4",
			wantCanonical:  true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := &httpUpstreamRecorder{resp: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"placeholderplaceholder,
				Body:       io.NopCloser(strings.NewReader(`{"id":"resp_compact","model":"` + tt.compactModel + `","usage":{"input_tokens":1,"output_tokens":1placeholderplaceholder`)),
		placeholderplaceholder
			svc := newOpenAIImageGenerationControlTestService(upstream)
			c, recorder := newOpenAIImageGenerationControlTestContext(false, "unit-test-agent/1.0")
			c.Request = httptest.NewRequest(http.MethodPost, "/openai/v1/responses/compact", nil)
			SetOpenAIClientTransport(c, OpenAIClientTransportHTTP)
			account := newOpenAIImageGenerationControlTestAccount()
			account.Extra = map[string]any{"openai_passthrough": trueplaceholder
			account.Credentials = map[string]any{
				"api_key": "sk-test",
				"compact_model_mapping": map[string]any{
					tt.canonicalModel: tt.compactModel,
			placeholder,
		placeholder
			body := []byte(`{"model":"` + tt.canonicalModel + `","stream":false,"input":"draw"placeholder`)

			result, err := svc.Forward(context.Background(), c, account, body)

			cached, known := getOpenAIImageIntentHint(c)
			require.True(t, known)
			require.Equal(t, tt.wantCanonical, cached)
			if tt.wantRejected {
			placeholder
				require.Nil(t, result)
				require.Equal(t, http.StatusForbidden, recorder.Code)
				require.Nil(t, upstream.lastReq)
				return
		placeholder
		placeholder
			require.NotNil(t, result)
			require.NotNil(t, upstream.lastReq)
			require.Equal(t, tt.compactModel, gjson.GetBytes(upstream.lastBody, "model").String())
	placeholder)
placeholder
placeholder

func TestResolveOpenAIImageIntentHintExcludesWebSocketAndUnknownTransport(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, transport := range []OpenAIClientTransport{OpenAIClientTransportWS, OpenAIClientTransportUnknownplaceholder {
		c := newOpenAIImageIntentHintTestContext(transport)
		var calls atomic.Int64
		classify := countingOpenAIImageIntentClassifier(&calls)
		body := []byte(`{"model":"gpt-5.4","input":"write code"placeholder`)

		require.False(t, resolveOpenAIImageIntentHint(c, "gpt-5.4", body, classify))
		require.False(t, resolveOpenAIImageIntentHint(c, "gpt-5.4", body, classify))
		require.Equal(t, int64(2), calls.Load())
		_, known := getOpenAIImageIntentHint(c)
		require.False(t, known)
placeholder
placeholder

func TestResolveOpenAIImageIntentHintConcurrentRequestsAreIsolated(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const requests = 32
	var calls atomic.Int64
	classify := countingOpenAIImageIntentClassifier(&calls)
	var wg sync.WaitGroup
	results := make([][2]bool, requests)

	for i := range requests {
		wg.Add(1)
		go func(index int, image bool) {
			defer wg.Done()
			c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
			body := []byte(`{"model":"gpt-5.4","input":"write code"placeholder`)
			if image {
				body = []byte(`{"model":"gpt-5.4","tools":[{"type":"image_generation"placeholder]placeholder`)
		placeholder
			results[index][0] = resolveOpenAIImageIntentHint(c, "gpt-5.4", body, classify)
			results[index][1] = resolveOpenAIImageIntentHint(c, "gpt-5.4", body, classify)
	placeholder(i, i%2 == 0)
placeholder
	wg.Wait()
	for i, result := range results {
		require.Equal(t, i%2 == 0, result[0])
		require.Equal(t, result[0], result[1])
placeholder
	require.Equal(t, int64(requests), calls.Load())
placeholder

var openAIImageIntentHintBenchmarkSink bool

func BenchmarkOpenAIPassthroughImageIntentHintLargeBody(b *testing.B) {
	body := []byte(`{"model":"gpt-5.4","input":"` + strings.Repeat("x", 4<<20) + `"placeholder`)
	const attempts = 4

	b.Run("scan_each_attempt", func(b *testing.B) {
		c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
		b.ReportAllocs()
		calls := 0
		for range b.N {
			c.Set(openAIImageIntentHintContextKey, struct{placeholder{placeholder)
			for range attempts {
				calls++
				openAIImageIntentHintBenchmarkSink = IsImageGenerationIntent(openAIResponsesEndpoint, "gpt-5.4", body)
		placeholder
	placeholder
		b.ReportMetric(float64(calls)/float64(b.N), "classifier_calls/op")
placeholder)

	b.Run("request_scoped_hint", func(b *testing.B) {
		c := newOpenAIImageIntentHintTestContext(OpenAIClientTransportHTTP)
		b.ReportAllocs()
		calls := 0
		classify := func(endpoint string, requestedModel string, candidate []byte) bool {
			calls++
			return IsImageGenerationIntent(endpoint, requestedModel, candidate)
	placeholder
		for range b.N {
			c.Set(openAIImageIntentHintContextKey, struct{placeholder{placeholder)
			for range attempts {
				openAIImageIntentHintBenchmarkSink = resolveOpenAIPassthroughImageIntent(c, "gpt-5.4", body, "gpt-5.4", body, false, classify)
		placeholder
	placeholder
		b.ReportMetric(float64(calls)/float64(b.N), "classifier_calls/op")
placeholder)
placeholder
