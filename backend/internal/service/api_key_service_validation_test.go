//go:build unit

package service

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCreateAPIKeyRequestNumericLimits(t *testing.T) {
	positiveExpiry := 1
	require.NoError(t, validateCreateAPIKeyRequest(CreateAPIKeyRequest{
		Quota: 1e100, RateLimit5h: 1e100, ExpiresInDays: &positiveExpiry,
placeholder))
	require.NoError(t, validateCreateAPIKeyRequest(CreateAPIKeyRequest{placeholder))

	invalidExpiry := 0
	tests := []CreateAPIKeyRequest{
		{Quota: -1placeholder,
		{Quota: math.NaN()placeholder,
		{Quota: math.Inf(1)placeholder,
		{RateLimit5h: -1placeholder,
		{RateLimit1d: math.NaN()placeholder,
		{RateLimit7d: math.Inf(-1)placeholder,
		{ExpiresInDays: &invalidExpiryplaceholder,
placeholder
	for _, req := range tests {
		require.Error(t, validateCreateAPIKeyRequest(req))
placeholder
placeholder

func TestValidateUpdateAPIKeyRequestNumericLimits(t *testing.T) {
	zero, large, negative, nan, inf := 0.0, 1e100, -1.0, math.NaN(), math.Inf(1)
	require.NoError(t, validateUpdateAPIKeyRequest(UpdateAPIKeyRequest{Quota: &zero, RateLimit7d: &largeplaceholder))

	for _, req := range []UpdateAPIKeyRequest{
		{Quota: &negativeplaceholder,
		{RateLimit5h: &nanplaceholder,
		{RateLimit1d: &infplaceholder,
		{RateLimit7d: &negativeplaceholder,
placeholder {
		require.Error(t, validateUpdateAPIKeyRequest(req))
placeholder
placeholder
