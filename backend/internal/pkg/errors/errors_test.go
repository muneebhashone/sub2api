//go:build unit

package errors

import (
	stderrors "errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplicationError_Basics(t *testing.T) {
	tests := []struct {
		name    string
		err     *ApplicationError
		want    Status
		wantIs  bool
		target  error
		wrapped error
placeholder{
		{
			name: "new",
			err:  New(400, "BAD_REQUEST", "invalid input"),
			want: Status{
				Code:    400,
				Reason:  "BAD_REQUEST",
				Message: "invalid input",
		placeholder,
	placeholder,
		{
			name:   "is_matches_code_and_reason",
			err:    New(401, "UNAUTHORIZED", "nope"),
			want:   Status{Code: 401, Reason: "UNAUTHORIZED", Message: "nope"placeholder,
			target: New(401, "UNAUTHORIZED", "ignored message"),
			wantIs: true,
	placeholder,
		{
			name:   "is_does_not_match_reason",
			err:    New(401, "UNAUTHORIZED", "nope"),
			want:   Status{Code: 401, Reason: "UNAUTHORIZED", Message: "nope"placeholder,
			target: New(401, "DIFFERENT", "ignored message"),
			wantIs: false,
	placeholder,
		{
			name:    "from_error_unwraps_wrapped_application_error",
			err:     New(404, "NOT_FOUND", "missing"),
			wrapped: fmt.Errorf("wrap: %w", New(404, "NOT_FOUND", "missing")),
			want: Status{
				Code:    404,
				Reason:  "NOT_FOUND",
				Message: "missing",
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				require.Equal(t, tt.want, tt.err.Status)
		placeholder

			if tt.target != nil {
				require.Equal(t, tt.wantIs, stderrors.Is(tt.err, tt.target))
		placeholder

			if tt.wrapped != nil {
				got := FromError(tt.wrapped)
				require.Equal(t, tt.want, got.Status)
		placeholder
	placeholder)
placeholder
placeholder

func TestApplicationError_WithMetadataDeepCopy(t *testing.T) {
	tests := []struct {
		name string
		md   map[string]string
placeholder{
		{name: "non_nil", md: map[string]string{"a": "1"placeholderplaceholder,
		{name: "nil", md: nilplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appErr := BadRequest("BAD_REQUEST", "invalid input").WithMetadata(tt.md)

			if tt.md == nil {
				require.Nil(t, appErr.Metadata)
				return
		placeholder

			tt.md["a"] = "changed"
			require.Equal(t, "1", appErr.Metadata["a"])
	placeholder)
placeholder
placeholder

func TestFromError_Generic(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantCode   int32
		wantReason string
		wantMsg    string
placeholder{
		{
			name:       "plain_error",
			err:        stderrors.New("boom"),
			wantCode:   UnknownCode,
			wantReason: UnknownReason,
			wantMsg:    UnknownMessage,
	placeholder,
		{
			name:       "wrapped_plain_error",
			err:        fmt.Errorf("wrap: %w", io.EOF),
			wantCode:   UnknownCode,
			wantReason: UnknownReason,
			wantMsg:    UnknownMessage,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromError(tt.err)
			require.Equal(t, tt.wantCode, got.Code)
			require.Equal(t, tt.wantReason, got.Reason)
			require.Equal(t, tt.wantMsg, got.Message)
			require.Equal(t, tt.err, got.Unwrap())
	placeholder)
placeholder
placeholder

func TestToHTTP(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
		wantBody       Status
placeholder{
		{
			name:           "nil_error",
			err:            nil,
			wantStatusCode: http.StatusOK,
			wantBody:       Status{Code: int32(http.StatusOK)placeholder,
	placeholder,
		{
			name:           "application_error",
			err:            Forbidden("FORBIDDEN", "no access"),
			wantStatusCode: http.StatusForbidden,
			wantBody: Status{
				Code:    int32(http.StatusForbidden),
				Reason:  "FORBIDDEN",
				Message: "no access",
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body := ToHTTP(tt.err)
			require.Equal(t, tt.wantStatusCode, code)
			require.Equal(t, tt.wantBody, body)
	placeholder)
placeholder
placeholder

func TestToHTTP_MetadataDeepCopy(t *testing.T) {
	md := map[string]string{"k": "v"placeholder
	appErr := BadRequest("BAD_REQUEST", "invalid").WithMetadata(md)

	code, body := ToHTTP(appErr)
	require.Equal(t, http.StatusBadRequest, code)
	require.Equal(t, "v", body.Metadata["k"])

	md["k"] = "changed"
	require.Equal(t, "v", body.Metadata["k"])

	appErr.Metadata["k"] = "changed-again"
	require.Equal(t, "v", body.Metadata["k"])
placeholder
