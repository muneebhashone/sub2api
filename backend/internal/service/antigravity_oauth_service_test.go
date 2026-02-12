package service

import (
	"testing"
)

func TestResolveDefaultTierID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		loadRaw map[string]any
		want    string
placeholder{
		{
			name:    "nil loadRaw",
			loadRaw: nil,
			want:    "",
	placeholder,
		{
			name: "missing allowedTiers",
			loadRaw: map[string]any{
				"paidTier": map[string]any{"id": "g1-pro-tier"placeholder,
		placeholder,
			want: "",
	placeholder,
		{
			name:    "empty allowedTiers",
			loadRaw: map[string]any{"allowedTiers": []any{placeholderplaceholder,
			want:    "",
	placeholder,
		{
			name: "tier missing id field",
			loadRaw: map[string]any{
				"allowedTiers": []any{
					map[string]any{"isDefault": trueplaceholder,
			placeholder,
		placeholder,
			want: "",
	placeholder,
		{
			name: "allowedTiers but no default",
			loadRaw: map[string]any{
				"allowedTiers": []any{
					map[string]any{"id": "free-tier", "isDefault": falseplaceholder,
					map[string]any{"id": "standard-tier", "isDefault": falseplaceholder,
			placeholder,
		placeholder,
			want: "",
	placeholder,
		{
			name: "default tier found",
			loadRaw: map[string]any{
				"allowedTiers": []any{
					map[string]any{"id": "free-tier", "isDefault": trueplaceholder,
					map[string]any{"id": "standard-tier", "isDefault": falseplaceholder,
			placeholder,
		placeholder,
			want: "free-tier",
	placeholder,
		{
			name: "default tier id with spaces",
			loadRaw: map[string]any{
				"allowedTiers": []any{
					map[string]any{"id": "  standard-tier  ", "isDefault": trueplaceholder,
			placeholder,
		placeholder,
			want: "standard-tier",
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := resolveDefaultTierID(tc.loadRaw)
			if got != tc.want {
				t.Fatalf("resolveDefaultTierID() = %q, want %q", got, tc.want)
		placeholder
	placeholder)
placeholder
placeholder
