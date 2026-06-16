//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/xai"
	"github.com/stretchr/testify/require"
)

type grokOAuthClientStub struct {
	refreshResponse *xai.TokenResponse
placeholder

func (s *grokOAuthClientStub) ExchangeCode(context.Context, string, string, string, string, string) (*xai.TokenResponse, error) {
	return &xai.TokenResponse{placeholder, nil
placeholder

func (s *grokOAuthClientStub) RefreshToken(context.Context, string, string, string) (*xai.TokenResponse, error) {
	return s.refreshResponse, nil
placeholder

func TestGrokOAuthServiceRefreshTokenPreservesOriginalRefreshTokenWhenNotRotated(t *testing.T) {
	svc := NewGrokOAuthService(nil, &grokOAuthClientStub{
		refreshResponse: &xai.TokenResponse{
			AccessToken: "new-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
	placeholder,
placeholder)
	defer svc.Stop()

	info, err := svc.RefreshToken(context.Background(), "original-refresh-token", "", "client-id")
placeholder
	require.Equal(t, "new-access-token", info.AccessToken)
	require.Equal(t, "original-refresh-token", info.RefreshToken)
	require.Equal(t, "client-id", info.ClientID)
placeholder
