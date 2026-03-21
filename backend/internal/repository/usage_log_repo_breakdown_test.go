//go:build unit

package repository

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/stretchr/testify/require"
)

func TestResolveEndpointColumn(t *testing.T) {
	tests := []struct {
		endpointType string
		want         string
placeholder{
		{"inbound", "ul.inbound_endpoint"placeholder,
		{"upstream", "ul.upstream_endpoint"placeholder,
		{"path", "ul.inbound_endpoint || ' -> ' || ul.upstream_endpoint"placeholder,
		{"", "ul.inbound_endpoint"placeholder,        // default
		{"unknown", "ul.inbound_endpoint"placeholder, // fallback
placeholder

	for _, tc := range tests {
		t.Run(tc.endpointType, func(t *testing.T) {
			got := resolveEndpointColumn(tc.endpointType)
			require.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder

func TestResolveModelDimensionExpression(t *testing.T) {
	tests := []struct {
		modelType string
		want      string
placeholder{
		{usagestats.ModelSourceRequested, "COALESCE(NULLIF(TRIM(requested_model), ''), model)"placeholder,
		{usagestats.ModelSourceUpstream, "COALESCE(NULLIF(TRIM(upstream_model), ''), COALESCE(NULLIF(TRIM(requested_model), ''), model))"placeholder,
		{usagestats.ModelSourceMapping, "(COALESCE(NULLIF(TRIM(requested_model), ''), model) || ' -> ' || COALESCE(NULLIF(TRIM(upstream_model), ''), COALESCE(NULLIF(TRIM(requested_model), ''), model)))"placeholder,
		{"", "COALESCE(NULLIF(TRIM(requested_model), ''), model)"placeholder,
		{"invalid", "COALESCE(NULLIF(TRIM(requested_model), ''), model)"placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.modelType, func(t *testing.T) {
			got := resolveModelDimensionExpression(tc.modelType)
			require.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder
