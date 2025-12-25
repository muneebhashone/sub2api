package errors

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	UnknownCode   = http.StatusInternalServerError
	UnknownReason = ""
)

type Status struct {
	Code     int32             `json:"code"`
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata,omitempty"`
placeholder

// ApplicationError is the standard error type used to control HTTP responses.
//
// Code is expected to be an HTTP status code (e.g. 400/401/403/404/409/500).
type ApplicationError struct {
	Status
	cause error
placeholder

// Error is kept for backwards compatibility within this package.
type Error = ApplicationError

func (e *ApplicationError) Error() string {
	if e == nil {
		return "<nil>"
placeholder
	if e.cause == nil {
		return fmt.Sprintf("error: code=%d reason=%q message=%q metadata=%v", e.Code, e.Reason, e.Message, e.Metadata)
placeholder
	return fmt.Sprintf("error: code=%d reason=%q message=%q metadata=%v cause=%v", e.Code, e.Reason, e.Message, e.Metadata, e.cause)
placeholder

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *ApplicationError) Unwrap() error { return e.cause placeholder

// Is matches each error in the chain with the target value.
func (e *ApplicationError) Is(err error) bool {
	if se := new(ApplicationError); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
placeholder
	return false
placeholder

// WithCause attaches the underlying cause of the error.
func (e *ApplicationError) WithCause(cause error) *ApplicationError {
	err := Clone(e)
	err.cause = cause
	return err
placeholder

// WithMetadata deep-copies the given metadata map.
func (e *ApplicationError) WithMetadata(md map[string]string) *ApplicationError {
	err := Clone(e)
	if md == nil {
		err.Metadata = nil
		return err
placeholder
	err.Metadata = make(map[string]string, len(md))
	for k, v := range md {
		err.Metadata[k] = v
placeholder
	return err
placeholder

// New returns an error object for the code, message.
func New(code int, reason, message string) *ApplicationError {
	return &ApplicationError{
		Status: Status{
			Code:    int32(code),
			Message: message,
			Reason:  reason,
	placeholder,
placeholder
placeholder

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code int, reason, format string, a ...any) *ApplicationError {
	return New(code, reason, fmt.Sprintf(format, a...))
placeholder

// Errorf returns an error object for the code, message and error info.
func Errorf(code int, reason, format string, a ...any) error {
	return New(code, reason, fmt.Sprintf(format, a...))
placeholder

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return http.StatusOK
placeholder
	return int(FromError(err).Code)
placeholder

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
placeholder
	return FromError(err).Reason
placeholder

// Message returns the message for a particular error.
// It supports wrapped errors.
func Message(err error) string {
	if err == nil {
		return ""
placeholder
	return FromError(err).Message
placeholder

// Clone deep clone error to a new error.
func Clone(err *ApplicationError) *ApplicationError {
	if err == nil {
		return nil
placeholder
	var metadata map[string]string
	if err.Metadata != nil {
		metadata = make(map[string]string, len(err.Metadata))
		for k, v := range err.Metadata {
			metadata[k] = v
	placeholder
placeholder
	return &ApplicationError{
		cause: err.cause,
		Status: Status{
			Code:     err.Code,
			Reason:   err.Reason,
			Message:  err.Message,
			Metadata: metadata,
	placeholder,
placeholder
placeholder

// FromError tries to convert an error to *ApplicationError.
// It supports wrapped errors.
func FromError(err error) *ApplicationError {
	if err == nil {
		return nil
placeholder
	if se := new(ApplicationError); errors.As(err, &se) {
		return se
placeholder

	// Fall back to a generic internal error.
	return New(UnknownCode, UnknownReason, err.Error()).WithCause(err)
placeholder
