//go:build unit

package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// --- marshalModelMapping ---

func TestMarshalModelMapping(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]map[string]string
		wantJSON string // expected JSON output (exact match)
placeholder{
		{
			name:     "empty map",
			input:    map[string]map[string]string{placeholder,
			wantJSON: "{placeholder",
	placeholder,
		{
			name:     "nil map",
			input:    nil,
			wantJSON: "{placeholder",
	placeholder,
		{
			name: "populated map",
			input: map[string]map[string]string{
				"openai": {"gpt-4": "gpt-4-turbo"placeholder,
		placeholder,
	placeholder,
		{
			name: "nested values",
			input: map[string]map[string]string{
				"openai":    {"*": "gpt-5.4"placeholder,
				"anthropic": {"claude-old": "claude-new"placeholder,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := marshalModelMapping(tt.input)
		placeholder

			if tt.wantJSON != "" {
				require.Equal(t, []byte(tt.wantJSON), result)
		placeholder else {
				// round-trip: unmarshal and compare with input
				var parsed map[string]map[string]string
				require.NoError(t, json.Unmarshal(result, &parsed))
				require.Equal(t, tt.input, parsed)
		placeholder
	placeholder)
placeholder
placeholder

// --- unmarshalModelMapping ---

func TestUnmarshalModelMapping(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantNil bool
		want    map[string]map[string]string
placeholder{
		{
			name:    "nil data",
			input:   nil,
			wantNil: true,
	placeholder,
		{
			name:    "empty data",
			input:   []byte{placeholder,
			wantNil: true,
	placeholder,
		{
			name:    "invalid JSON",
			input:   []byte("not-json"),
			wantNil: true,
	placeholder,
		{
			name:    "type error - number",
			input:   []byte("42"),
			wantNil: true,
	placeholder,
		{
			name:    "type error - array",
			input:   []byte("[1,2,3]"),
			wantNil: true,
	placeholder,
		{
			name:  "valid JSON",
			input: []byte(`{"openai":{"gpt-4":"gpt-4-turbo"placeholder,"anthropic":{"old":"new"placeholderplaceholder`),
			want: map[string]map[string]string{
				"openai":    {"gpt-4": "gpt-4-turbo"placeholder,
				"anthropic": {"old": "new"placeholder,
		placeholder,
	placeholder,
		{
			name:  "empty object",
			input: []byte("{placeholder"),
			want:  map[string]map[string]string{placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unmarshalModelMapping(tt.input)
			if tt.wantNil {
				require.Nil(t, result)
		placeholder else {
				require.NotNil(t, result)
				require.Equal(t, tt.want, result)
		placeholder
	placeholder)
placeholder
placeholder

// --- escapeLike ---

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
placeholder{
		{
			name:  "no special chars",
			input: "hello",
			want:  "hello",
	placeholder,
		{
			name:  "backslash",
			input: `a\b`,
			want:  `a\\b`,
	placeholder,
		{
			name:  "percent",
			input: "50%",
			want:  `50\%`,
	placeholder,
		{
			name:  "underscore",
			input: "a_b",
			want:  `a\_b`,
	placeholder,
		{
			name:  "all special chars",
			input: `a\b%c_d`,
			want:  `a\\b\%c\_d`,
	placeholder,
		{
			name:  "empty string",
			input: "",
			want:  "",
	placeholder,
		{
			name:  "consecutive special chars",
			input: "%_%",
			want:  `\%\_\%`,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, escapeLike(tt.input))
	placeholder)
placeholder
placeholder

// --- isUniqueViolation ---

func TestIsUniqueViolation(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
placeholder{
		{
			name: "unique violation code 23505",
			err:  &pq.Error{Code: "23505"placeholder,
			want: true,
	placeholder,
		{
			name: "different pq error code",
			err:  &pq.Error{Code: "23503"placeholder,
			want: false,
	placeholder,
		{
			name: "non-pq error",
			err:  errors.New("some generic error"),
			want: false,
	placeholder,
		{
			name: "typed nil pq.Error",
			err: func() error {
				var pqErr *pq.Error
				return pqErr
		placeholder(),
			want: false,
	placeholder,
		{
			name: "bare nil",
			err:  nil,
			want: false,
	placeholder,
		{
			name: "wrapped pq error with 23505",
			err:  fmt.Errorf("wrapped: %w", &pq.Error{Code: "23505"placeholder),
			want: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isUniqueViolation(tt.err))
	placeholder)
placeholder
placeholder
