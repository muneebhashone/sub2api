//go:build unit

package handler

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateAPIKeyCreateRequest(t *testing.T) {
	zero, large, negative, nan, inf := 0.0, 1e100, -1.0, math.NaN(), math.Inf(1)
	positiveDays, zeroDays, negativeDays := 1, 0, -1
	require.NoError(t, validateAPIKeyCreateRequest(CreateAPIKeyRequest{placeholder))
	require.NoError(t, validateAPIKeyCreateRequest(CreateAPIKeyRequest{Quota: &zero, RateLimit5h: &large, ExpiresInDays: &positiveDaysplaceholder))

	for _, req := range []CreateAPIKeyRequest{
		{Quota: &negativeplaceholder,
		{Quota: &nanplaceholder,
		{RateLimit5h: &infplaceholder,
		{RateLimit1d: &negativeplaceholder,
		{RateLimit7d: &negativeplaceholder,
		{ExpiresInDays: &zeroDaysplaceholder,
		{ExpiresInDays: &negativeDaysplaceholder,
placeholder {
		require.Error(t, validateAPIKeyCreateRequest(req))
placeholder
placeholder

func TestValidateAPIKeyUpdateRequest(t *testing.T) {
	zero, large, negative, nan, inf := 0.0, 1e100, -1.0, math.NaN(), math.Inf(-1)
	require.NoError(t, validateAPIKeyUpdateRequest(UpdateAPIKeyRequest{Quota: &zero, RateLimit7d: &largeplaceholder))

	for _, req := range []UpdateAPIKeyRequest{
		{Quota: &negativeplaceholder,
		{RateLimit5h: &nanplaceholder,
		{RateLimit1d: &infplaceholder,
		{RateLimit7d: &negativeplaceholder,
placeholder {
		require.Error(t, validateAPIKeyUpdateRequest(req))
placeholder
placeholder
