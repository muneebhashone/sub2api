package usagestats

import "testing"

func TestIsValidModelSource(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
placeholder{
		{name: "requested", source: ModelSourceRequested, want: trueplaceholder,
		{name: "upstream", source: ModelSourceUpstream, want: trueplaceholder,
		{name: "mapping", source: ModelSourceMapping, want: trueplaceholder,
		{name: "invalid", source: "foobar", want: falseplaceholder,
		{name: "empty", source: "", want: falseplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidModelSource(tc.source); got != tc.want {
				t.Fatalf("IsValidModelSource(%q)=%v want %v", tc.source, got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizeModelSource(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
placeholder{
		{name: "requested", source: ModelSourceRequested, want: ModelSourceRequestedplaceholder,
		{name: "upstream", source: ModelSourceUpstream, want: ModelSourceUpstreamplaceholder,
		{name: "mapping", source: ModelSourceMapping, want: ModelSourceMappingplaceholder,
		{name: "invalid falls back", source: "foobar", want: ModelSourceRequestedplaceholder,
		{name: "empty falls back", source: "", want: ModelSourceRequestedplaceholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := NormalizeModelSource(tc.source); got != tc.want {
				t.Fatalf("NormalizeModelSource(%q)=%q want %q", tc.source, got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder
