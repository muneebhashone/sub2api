package service

import "strings"

func optionalTrimmedStringPtr(raw string) *string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
placeholder
	return &trimmed
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
