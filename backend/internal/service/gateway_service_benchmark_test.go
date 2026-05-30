package service

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/domain"
)

var (
	benchmarkStringSink string
	benchmarkIntSink    int
)

// BenchmarkGenerateSessionHash_Metadata 关注 JSON 解析与正则匹配开销。
func BenchmarkGenerateSessionHash_Metadata(b *testing.B) {
	svc := &GatewayService{placeholder
	body := []byte(`{"metadata":{"user_id":"session_123e4567-e89b-12d3-a456-426614174000"placeholder,"messages":[{"content":"hello"placeholder]placeholder`)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), "")
		if err != nil {
			b.Fatalf("解析请求失败: %v", err)
	placeholder
		benchmarkStringSink = svc.GenerateSessionHash(parsed)
placeholder
placeholder

func BenchmarkParseGatewayRequest_LargeAnthropicMessages(b *testing.B) {
	for _, size := range benchmarkBodySizes() {
		b.Run(size.name, func(b *testing.B) {
			body := buildLargeAnthropicMessagesBody(size.bytes, false)

			b.SetBytes(int64(len(body)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), domain.PlatformAnthropic)
				if err != nil {
					b.Fatalf("解析 Anthropic 请求失败: %v", err)
			placeholder
				benchmarkIntSink = len(parsed.MessagesRaw())
		placeholder
	placeholder)
placeholder
placeholder

func BenchmarkParseGatewayRequest_LargeGeminiContents(b *testing.B) {
	for _, size := range benchmarkBodySizes() {
		b.Run(size.name, func(b *testing.B) {
			body := buildLargeGeminiContentsBody(size.bytes)

			b.SetBytes(int64(len(body)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), domain.PlatformGemini)
				if err != nil {
					b.Fatalf("解析 Gemini 请求失败: %v", err)
			placeholder
				benchmarkIntSink = len(parsed.MessagesRaw())
		placeholder
	placeholder)
placeholder
placeholder

func BenchmarkGenerateSessionHash_LargeAnthropicMessages(b *testing.B) {
	svc := &GatewayService{placeholder
	for _, size := range benchmarkBodySizes() {
		b.Run(size.name, func(b *testing.B) {
			body := buildLargeAnthropicMessagesBody(size.bytes, true)
			parsed, err := ParseGatewayRequest(NewRequestBodyRef(body), domain.PlatformAnthropic)
			if err != nil {
				b.Fatalf("解析请求失败: %v", err)
		placeholder

			b.SetBytes(int64(len(body)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkStringSink = svc.GenerateSessionHash(parsed)
		placeholder
	placeholder)
placeholder
placeholder

func BenchmarkOpenAIResponses_LargeInputMeta(b *testing.B) {
	for _, size := range benchmarkBodySizes() {
		b.Run(size.name, func(b *testing.B) {
			body := buildLargeOpenAIResponsesBody(size.bytes)

			b.SetBytes(int64(len(body)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				model, stream, promptCacheKey := extractOpenAIRequestMetaFromBody(body)
				benchmarkStringSink = model + promptCacheKey
				if stream {
					benchmarkIntSink++
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func BenchmarkOpenAIResponses_LargeInputDecodeMap(b *testing.B) {
	for _, size := range benchmarkBodySizes() {
		b.Run(size.name, func(b *testing.B) {
			body := buildLargeOpenAIResponsesBody(size.bytes)

			b.SetBytes(int64(len(body)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reqBody, err := getOpenAIRequestBodyMap(nil, body)
				if err != nil {
					b.Fatalf("解析 OpenAI 请求失败: %v", err)
			placeholder
				benchmarkIntSink = len(reqBody)
		placeholder
	placeholder)
placeholder
placeholder

func BenchmarkOpenAIResponses_LargeInputFunctionCallValidation(b *testing.B) {
	for _, size := range benchmarkBodySizes() {
		b.Run(size.name, func(b *testing.B) {
			body := buildLargeOpenAIResponsesToolContinuationBody(size.bytes)

			b.SetBytes(int64(len(body)))
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				validation := ValidateFunctionCallOutputContextBytes(body)
				if !validation.HasFunctionCallOutput || !validation.HasItemReferenceForAllCallIDs {
					b.Fatalf("工具续链校验结果异常: %+v", validation)
			placeholder
				benchmarkIntSink++
		placeholder
	placeholder)
placeholder
placeholder

// BenchmarkExtractCacheableContent_System 关注字符串拼接路径的性能。
func BenchmarkExtractCacheableContent_System(b *testing.B) {
	svc := &GatewayService{placeholder
	req := buildSystemCacheableRequest(12)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchmarkStringSink = svc.extractCacheableContent(req)
placeholder
placeholder

func benchmarkBodySizes() []struct {
	name  string
	bytes int
placeholder {
	return []struct {
		name  string
		bytes int
placeholder{
		{name: "4MB", bytes: 4 << 20placeholder,
		{name: "8MB", bytes: 8 << 20placeholder,
		{name: "16MB", bytes: 16 << 20placeholder,
		{name: "32MB", bytes: 32 << 20placeholder,
placeholder
placeholder

func buildSystemCacheableRequest(parts int) *ParsedRequest {
	var builder strings.Builder
	builder.WriteString(`{"system":[`)
	for i := 0; i < parts; i++ {
		if i > 0 {
			builder.WriteByte(',')
	placeholder
		builder.WriteString(`{"text":"system_part_`)
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(`","cache_control":{"type":"ephemeral"placeholderplaceholder`)
placeholder
	builder.WriteString(`]placeholder`)
	parsed, err := ParseGatewayRequest(NewRequestBodyRef([]byte(builder.String())), "")
	if err != nil {
		panic(err)
placeholder
	return parsed
placeholder

func buildLargeAnthropicMessagesBody(targetBytes int, includeCacheControl bool) []byte {
	var builder strings.Builder
	builder.Grow(targetBytes + 1024)
	builder.WriteString(`{"model":"claude-sonnet-4-5","stream":true,"system":[{"type":"text","text":"system seed"placeholder],"messages":[`)
	for i := 0; builder.Len() < targetBytes; i++ {
		if i > 0 {
			builder.WriteByte(',')
	placeholder
		builder.WriteString(`{"role":"user","content":[{"type":"text","text":"`)
		builder.WriteString(strings.Repeat("anthropic payload ", 64))
		builder.WriteString(strconv.Itoa(i))
		builder.WriteByte('"')
		if includeCacheControl && i%32 == 0 {
			builder.WriteString(`,"cache_control":{"type":"ephemeral"placeholder`)
	placeholder
		builder.WriteString(`placeholder]placeholder`)
placeholder
	builder.WriteString(`]placeholder`)
	return []byte(builder.String())
placeholder

func buildLargeGeminiContentsBody(targetBytes int) []byte {
	var builder strings.Builder
	builder.Grow(targetBytes + 1024)
	builder.WriteString(`{"model":"gemini-2.5-pro","systemInstruction":{"parts":[{"text":"system seed"placeholder]placeholder,"contents":[`)
	for i := 0; builder.Len() < targetBytes; i++ {
		if i > 0 {
			builder.WriteByte(',')
	placeholder
		builder.WriteString(`{"role":"user","parts":[{"text":"`)
		builder.WriteString(strings.Repeat("gemini payload ", 64))
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(`"placeholder]placeholder`)
placeholder
	builder.WriteString(`]placeholder`)
	return []byte(builder.String())
placeholder

func buildLargeOpenAIResponsesBody(targetBytes int) []byte {
	var builder strings.Builder
	builder.Grow(targetBytes + 1024)
	builder.WriteString(`{"model":"gpt-5.4","stream":true,"prompt_cache_key":"session-benchmark","input":[`)
	for i := 0; builder.Len() < targetBytes; i++ {
		if i > 0 {
			builder.WriteByte(',')
	placeholder
		builder.WriteString(`{"type":"message","role":"user","content":[{"type":"input_text","text":"`)
		builder.WriteString(strings.Repeat("openai responses payload ", 48))
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(`"placeholder]placeholder`)
placeholder
	builder.WriteString(`],"tools":[{"type":"function","name":"lookup","parameters":{"type":"object","properties":{"query":{"type":"string"placeholderplaceholderplaceholderplaceholder]placeholder`)
	return []byte(builder.String())
placeholder

func buildLargeOpenAIResponsesToolContinuationBody(targetBytes int) []byte {
	var builder strings.Builder
	builder.Grow(targetBytes + 1024)
	builder.WriteString(`{"model":"gpt-5.4","stream":true,"previous_response_id":"resp_benchmark","input":[`)
	for i := 0; builder.Len() < targetBytes; i++ {
		if i > 0 {
			builder.WriteByte(',')
	placeholder
		callID := "call_" + strconv.Itoa(i)
		builder.WriteString(`{"type":"item_reference","id":"`)
		builder.WriteString(callID)
		builder.WriteString(`"placeholder,{"type":"function_call_output","call_id":"`)
		builder.WriteString(callID)
		builder.WriteString(`","output":"`)
		builder.WriteString(strings.Repeat("tool output payload ", 48))
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString(`"placeholder`)
placeholder
	builder.WriteString(`]placeholder`)
	return []byte(builder.String())
placeholder
