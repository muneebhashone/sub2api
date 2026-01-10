//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsageService_InvalidateUsageCaches(t *testing.T) {
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &UsageService{authCacheInvalidator: invalidatorplaceholder

	svc.invalidateUsageCaches(context.Background(), 7, false)
	require.Empty(t, invalidator.userIDs)

	svc.invalidateUsageCaches(context.Background(), 7, true)
	require.Equal(t, []int64{7placeholder, invalidator.userIDs)
placeholder

func TestRedeemService_InvalidateRedeemCaches_AuthCache(t *testing.T) {
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &RedeemService{authCacheInvalidator: invalidatorplaceholder

	svc.invalidateRedeemCaches(context.Background(), 11, &RedeemCode{Type: RedeemTypeBalanceplaceholder)
	svc.invalidateRedeemCaches(context.Background(), 11, &RedeemCode{Type: RedeemTypeConcurrencyplaceholder)

	require.Equal(t, []int64{11, 11placeholder, invalidator.userIDs)
placeholder
