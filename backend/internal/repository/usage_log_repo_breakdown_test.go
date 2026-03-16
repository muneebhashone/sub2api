//go:build unit

package repository

import (
	"testing"

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
		{"", "ul.inbound_endpoint"placeholder,           // default
		{"unknown", "ul.inbound_endpoint"placeholder,     // fallback
placeholder

	for _, tc := range tests {
		t.Run(tc.endpointType, func(t *testing.T) {
			got := resolveEndpointColumn(tc.endpointType)
			require.Equal(t, tc.want, got)
	placeholder)
placeholder
placeholder
