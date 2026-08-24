package service

import (
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

var (
	upstreamResponseModelBenchmarkSink string

	upstreamResponseModelBenchmarkTerminal = []byte(`{"type":"response.completed","response":{"id":"resp_123","model":"gpt-5.5-2026-04-23"placeholder,"usage":{"input_tokens":128,"output_tokens":64placeholderplaceholder`)
	upstreamResponseModelBenchmarkDelta    = []byte(`{"type":"response.output_text.delta","delta":"hello"placeholder`)
	upstreamResponseModelBenchmarkWrapper  = []byte(`{"response":{"response":{"modelVersion":"gemini-3-pro","candidates":[]placeholderplaceholderplaceholder`)
)

func BenchmarkUpstreamResponseModelOpenAI(b *testing.B) {
	tests := []struct {
		name    string
		payload []byte
placeholder{
		{name: "terminal_with_model", payload: upstreamResponseModelBenchmarkTerminalplaceholder,
		{name: "delta_without_model", payload: upstreamResponseModelBenchmarkDeltaplaceholder,
placeholder

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.Run("legacy", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					upstreamResponseModelBenchmarkSink = benchmarkLegacyOpenAIModel(tt.payload)
			placeholder
		placeholder)
			b.Run("optimized", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					upstreamResponseModelBenchmarkSink = firstValidTrimmedGJSONString(tt.payload, "response.model", "model")
			placeholder
		placeholder)
	placeholder)
placeholder
placeholder

func BenchmarkUpstreamResponseModelAntigravityWrapper(b *testing.B) {
	b.Run("legacy", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			upstreamResponseModelBenchmarkSink = benchmarkLegacyWrappedGeminiModel(upstreamResponseModelBenchmarkWrapper)
	placeholder
placeholder)
	b.Run("optimized", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			upstreamResponseModelBenchmarkSink = firstValidTrimmedGJSONString(
				upstreamResponseModelBenchmarkWrapper,
				"modelVersion",
				"response.modelVersion",
				"response.response.modelVersion",
			)
	placeholder
placeholder)
placeholder

func benchmarkLegacyOpenAIModel(payload []byte) string {
	if len(payload) == 0 || !gjson.ValidBytes(payload) {
		return ""
placeholder
	return benchmarkFirstTrimmedGJSONModel(
		gjson.GetBytes(payload, "response.model"),
		gjson.GetBytes(payload, "model"),
	)
placeholder

func benchmarkLegacyWrappedGeminiModel(payload []byte) string {
	if inner := gjson.GetBytes(payload, "response"); inner.Exists() {
		payload = []byte(inner.Raw)
placeholder
	if len(payload) == 0 || !gjson.ValidBytes(payload) {
		return ""
placeholder
	return benchmarkFirstTrimmedGJSONModel(
		gjson.GetBytes(payload, "modelVersion"),
		gjson.GetBytes(payload, "response.modelVersion"),
	)
placeholder

func benchmarkFirstTrimmedGJSONModel(values ...gjson.Result) string {
	for _, value := range values {
		if !value.Exists() || value.Type != gjson.String {
			continue
	placeholder
		if model := strings.TrimSpace(value.String()); model != "" {
			return model
	placeholder
placeholder
	return ""
placeholder
