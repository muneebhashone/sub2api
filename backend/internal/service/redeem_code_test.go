package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedeemCodeExpiry(t *testing.T) {
	now := time.Now().UTC()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tests := []struct {
		name        string
		code        RedeemCode
		wantExpired bool
		wantCanUse  bool
placeholder{
		{
			name:        "unused without expiry can be used",
			code:        RedeemCode{Status: StatusUnusedplaceholder,
			wantExpired: false,
			wantCanUse:  true,
	placeholder,
		{
			name:        "unused before expiry can be used",
			code:        RedeemCode{Status: StatusUnused, ExpiresAt: &futureplaceholder,
			wantExpired: false,
			wantCanUse:  true,
	placeholder,
		{
			name:        "unused after expiry cannot be used",
			code:        RedeemCode{Status: StatusUnused, ExpiresAt: &pastplaceholder,
			wantExpired: true,
			wantCanUse:  false,
	placeholder,
		{
			name:        "explicit expired status is expired",
			code:        RedeemCode{Status: StatusExpiredplaceholder,
			wantExpired: true,
			wantCanUse:  false,
	placeholder,
		{
			name:        "used code remains used even after expiry time",
			code:        RedeemCode{Status: StatusUsed, ExpiresAt: &pastplaceholder,
			wantExpired: false,
			wantCanUse:  false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantExpired, tt.code.IsExpiredAt(now))
			require.Equal(t, tt.wantCanUse, tt.code.CanUse())
	placeholder)
placeholder
placeholder
