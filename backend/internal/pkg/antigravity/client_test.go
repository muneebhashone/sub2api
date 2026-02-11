package antigravity

import (
	"testing"
)

func TestExtractProjectIDFromOnboardResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		resp map[string]any
		want string
placeholder{
		{
			name: "nil response",
			resp: nil,
			want: "",
	placeholder,
		{
			name: "empty response",
			resp: map[string]any{placeholder,
			want: "",
	placeholder,
		{
			name: "project as string",
			resp: map[string]any{
				"cloudaicompanionProject": "my-project-123",
		placeholder,
			want: "my-project-123",
	placeholder,
		{
			name: "project as string with spaces",
			resp: map[string]any{
				"cloudaicompanionProject": "  my-project-123  ",
		placeholder,
			want: "my-project-123",
	placeholder,
		{
			name: "project as map with id",
			resp: map[string]any{
				"cloudaicompanionProject": map[string]any{
					"id": "proj-from-map",
			placeholder,
		placeholder,
			want: "proj-from-map",
	placeholder,
		{
			name: "project as map without id",
			resp: map[string]any{
				"cloudaicompanionProject": map[string]any{
					"name": "some-name",
			placeholder,
		placeholder,
			want: "",
	placeholder,
		{
			name: "missing cloudaicompanionProject key",
			resp: map[string]any{
				"otherField": "value",
		placeholder,
			want: "",
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := extractProjectIDFromOnboardResponse(tc.resp)
			if got != tc.want {
				t.Fatalf("extractProjectIDFromOnboardResponse() = %q, want %q", got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder
