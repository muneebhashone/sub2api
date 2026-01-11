package repository

import (
	"fmt"
	"strings"
)

type latencyHistogramBucket struct {
	upperMs int
	label   string
placeholder

var latencyHistogramBuckets = []latencyHistogramBucket{
	{upperMs: 100, label: "0-100ms"placeholder,
	{upperMs: 200, label: "100-200ms"placeholder,
	{upperMs: 500, label: "200-500ms"placeholder,
	{upperMs: 1000, label: "500-1000ms"placeholder,
	{upperMs: 2000, label: "1000-2000ms"placeholder,
	{upperMs: 0, label: "2000ms+"placeholder, // default bucket
placeholder

var latencyHistogramOrderedRanges = func() []string {
	out := make([]string, 0, len(latencyHistogramBuckets))
	for _, b := range latencyHistogramBuckets {
		out = append(out, b.label)
placeholder
	return out
placeholder()

func latencyHistogramRangeCaseExpr(column string) string {
	var sb strings.Builder
	_ = sb.WriteString("CASE\n")

	for _, b := range latencyHistogramBuckets {
		if b.upperMs <= 0 {
			continue
	placeholder
		_ = sb.WriteString(fmt.Sprintf("\tWHEN %s < %d THEN '%s'\n", column, b.upperMs, b.label))
placeholder

	// Default bucket.
	last := latencyHistogramBuckets[len(latencyHistogramBuckets)-1]
	_ = sb.WriteString(fmt.Sprintf("\tELSE '%s'\n", last.label))
	sb.WriteString("END")
	return sb.String()
placeholder

func latencyHistogramRangeOrderCaseExpr(column string) string {
	var sb strings.Builder
	sb.WriteString("CASE\n")

	order := 1
	for _, b := range latencyHistogramBuckets {
		if b.upperMs <= 0 {
			continue
	placeholder
		sb.WriteString(fmt.Sprintf("\tWHEN %s < %d THEN %d\n", column, b.upperMs, order))
		order++
placeholder

	sb.WriteString(fmt.Sprintf("\tELSE %d\n", order))
	sb.WriteString("END")
	return sb.String()
placeholder
