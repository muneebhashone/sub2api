package service

import "strings"

func optionalTrimmedStringPtr(raw string) *string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
placeholder
	return &trimmed
placeholder

// optionalNonEqualStringPtr returns a pointer to value if it is non-empty and
// differs from compare; otherwise nil. Usage logging passes the requested
// model as compare so a channel mapping still records its effective upstream.
func optionalNonEqualStringPtr(value, compare string) *string {
	if value == "" || value == compare {
		return nil
placeholder
	return &value
placeholder

func forwardResultBillingModel(requestedModel, upstreamModel string) string {
	if trimmed := strings.TrimSpace(requestedModel); trimmed != "" {
		return trimmed
placeholder
	return strings.TrimSpace(upstreamModel)
placeholder

func optionalInt64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
placeholder
	return &v
placeholder
