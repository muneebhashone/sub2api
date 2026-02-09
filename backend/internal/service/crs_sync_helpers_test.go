package service

import (
	"testing"
)

func TestBuildSelectedSet(t *testing.T) {
	tests := []struct {
		name     string
		ids      []string
		wantNil  bool
		wantSize int
placeholder{
		{
			name:    "nil input returns nil (backward compatible: create all)",
			ids:     nil,
			wantNil: true,
	placeholder,
		{
			name:     "empty slice returns empty map (create none)",
			ids:      []string{placeholder,
			wantNil:  false,
			wantSize: 0,
	placeholder,
		{
			name:     "single ID",
			ids:      []string{"abc-123"placeholder,
			wantNil:  false,
			wantSize: 1,
	placeholder,
		{
			name:     "multiple IDs",
			ids:      []string{"a", "b", "c"placeholder,
			wantNil:  false,
			wantSize: 3,
	placeholder,
		{
			name:     "duplicate IDs are deduplicated",
			ids:      []string{"a", "a", "b"placeholder,
			wantNil:  false,
			wantSize: 2,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSelectedSet(tt.ids)
			if tt.wantNil {
				if got != nil {
					t.Errorf("buildSelectedSet(%v) = %v, want nil", tt.ids, got)
			placeholder
				return
		placeholder
			if got == nil {
				t.Fatalf("buildSelectedSet(%v) = nil, want non-nil map", tt.ids)
		placeholder
			if len(got) != tt.wantSize {
				t.Errorf("buildSelectedSet(%v) has %d entries, want %d", tt.ids, len(got), tt.wantSize)
		placeholder
			// Verify all unique IDs are present
			for _, id := range tt.ids {
				if _, ok := got[id]; !ok {
					t.Errorf("buildSelectedSet(%v) missing key %q", tt.ids, id)
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestShouldCreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		crsID       string
		selectedSet map[string]struct{placeholder
		want        bool
placeholder{
		{
			name:        "nil set allows all (backward compatible)",
			crsID:       "any-id",
			selectedSet: nil,
			want:        true,
	placeholder,
		{
			name:        "empty set blocks all",
			crsID:       "any-id",
			selectedSet: map[string]struct{placeholder{placeholder,
			want:        false,
	placeholder,
		{
			name:        "ID in set is allowed",
			crsID:       "abc-123",
			selectedSet: map[string]struct{placeholder{"abc-123": {placeholder, "def-456": {placeholderplaceholder,
			want:        true,
	placeholder,
		{
			name:        "ID not in set is blocked",
			crsID:       "xyz-789",
			selectedSet: map[string]struct{placeholder{"abc-123": {placeholder, "def-456": {placeholderplaceholder,
			want:        false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldCreateAccount(tt.crsID, tt.selectedSet)
			if got != tt.want {
				t.Errorf("shouldCreateAccount(%q, %v) = %v, want %v",
					tt.crsID, tt.selectedSet, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder
