//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type balanceUserRepoStub struct {
	*userRepoStub
	updateErr error
	updated   []*User
placeholder

func (s *balanceUserRepoStub) Update(ctx context.Context, user *User) error {
	if s.updateErr != nil {
		return s.updateErr
placeholder
	if user == nil {
		return nil
placeholder
	clone := *user
	s.updated = append(s.updated, &clone)
	if s.userRepoStub != nil {
		s.userRepoStub.user = &clone
placeholder
	return nil
placeholder

type balanceRedeemRepoStub struct {
	*redeemRepoStub
	created []*RedeemCode
placeholder

func (s *balanceRedeemRepoStub) Create(ctx context.Context, code *RedeemCode) error {
	if code == nil {
		return nil
placeholder
	clone := *code
	s.created = append(s.created, &clone)
	return nil
placeholder

type authCacheInvalidatorStub struct {
	userIDs  []int64
	groupIDs []int64
	keys     []string
placeholder

type adminRechargeAffiliateAccruerStub struct {
	calls  []adminRechargeAffiliateAccrual
	rebate float64
	err    error
placeholder

type adminRechargeAffiliateAccrual struct {
	userID int64
	amount float64
placeholder

func (s *adminRechargeAffiliateAccruerStub) AccrueInviteRebate(_ context.Context, userID int64, amount float64) (float64, error) {
	s.calls = append(s.calls, adminRechargeAffiliateAccrual{userID: userID, amount: amountplaceholder)
	return s.rebate, s.err
placeholder

func adminRechargeSettingService(enabled bool) *SettingService {
	values := map[string]string{placeholder
	if enabled {
		values[SettingKeyAffiliateAdminRechargeEnabled] = "true"
placeholder
	return NewSettingService(&settingRepoStub{values: valuesplaceholder, nil)
placeholder

func (s *authCacheInvalidatorStub) InvalidateAuthCacheByKey(ctx context.Context, key string) {
	s.keys = append(s.keys, key)
placeholder

func (s *authCacheInvalidatorStub) InvalidateAuthCacheByUserID(ctx context.Context, userID int64) {
	s.userIDs = append(s.userIDs, userID)
placeholder

func (s *authCacheInvalidatorStub) InvalidateAuthCacheByGroupID(ctx context.Context, groupID int64) {
	s.groupIDs = append(s.groupIDs, groupID)
placeholder

func TestAdminService_UpdateUserBalance_InvalidatesAuthCache(t *testing.T) {
	baseRepo := &userRepoStub{user: &User{ID: 7, Balance: 10placeholderplaceholder
	repo := &balanceUserRepoStub{userRepoStub: baseRepoplaceholder
	redeemRepo := &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		userRepo:             repo,
		redeemCodeRepo:       redeemRepo,
		authCacheInvalidator: invalidator,
placeholder

	_, err := svc.UpdateUserBalance(context.Background(), 7, 5, "add", "")
placeholder
	require.Equal(t, []int64{7placeholder, invalidator.userIDs)
	require.Len(t, redeemRepo.created, 1)
placeholder

func TestAdminService_UpdateUserBalance_NoChangeNoInvalidate(t *testing.T) {
	baseRepo := &userRepoStub{user: &User{ID: 7, Balance: 10placeholderplaceholder
	repo := &balanceUserRepoStub{userRepoStub: baseRepoplaceholder
	redeemRepo := &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		userRepo:             repo,
		redeemCodeRepo:       redeemRepo,
		authCacheInvalidator: invalidator,
placeholder

	_, err := svc.UpdateUserBalance(context.Background(), 7, 10, "set", "")
placeholder
	require.Empty(t, invalidator.userIDs)
	require.Empty(t, redeemRepo.created)
placeholder

func TestAdminService_UpdateUserBalance_AdminRechargeAffiliateRebate(t *testing.T) {
	tests := []struct {
		name      string
		enabled   bool
		operation string
		amount    float64
		wantCalls []adminRechargeAffiliateAccrual
placeholder{
		{
			name:      "disabled by default",
			operation: "add",
			amount:    5,
	placeholder,
		{
			name:      "enabled add",
			enabled:   true,
			operation: "add",
			amount:    0.1,
			wantCalls: []adminRechargeAffiliateAccrual{{userID: 7, amount: 0.1placeholderplaceholder,
	placeholder,
		{
			name:      "enabled set increase",
			enabled:   true,
			operation: "set",
			amount:    15,
	placeholder,
		{
			name:      "enabled subtract",
			enabled:   true,
			operation: "subtract",
			amount:    5,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseRepo := &userRepoStub{user: &User{ID: 7, Balance: 10placeholderplaceholder
			repo := &balanceUserRepoStub{userRepoStub: baseRepoplaceholder
			redeemRepo := &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder
			affiliate := &adminRechargeAffiliateAccruerStub{placeholder
			svc := &adminServiceImpl{
				userRepo:         repo,
				redeemCodeRepo:   redeemRepo,
				settingService:   adminRechargeSettingService(tt.enabled),
				affiliateService: affiliate,
		placeholder

			_, err := svc.UpdateUserBalance(context.Background(), 7, tt.amount, tt.operation, "")
		placeholder
			require.Equal(t, tt.wantCalls, affiliate.calls)
	placeholder)
placeholder
placeholder

func TestAdminService_UpdateUserBalance_AffiliateFailureDoesNotRollbackRecharge(t *testing.T) {
	baseRepo := &userRepoStub{user: &User{ID: 7, Balance: 10placeholderplaceholder
	repo := &balanceUserRepoStub{userRepoStub: baseRepoplaceholder
	redeemRepo := &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder
	affiliate := &adminRechargeAffiliateAccruerStub{err: errors.New("affiliate unavailable")placeholder
	svc := &adminServiceImpl{
		userRepo:         repo,
		redeemCodeRepo:   redeemRepo,
		settingService:   adminRechargeSettingService(true),
		affiliateService: affiliate,
placeholder

	user, err := svc.UpdateUserBalance(context.Background(), 7, 5, "add", "")
placeholder
	require.Equal(t, 15.0, user.Balance)
	require.Equal(t, []adminRechargeAffiliateAccrual{{userID: 7, amount: 5placeholderplaceholder, affiliate.calls)
	require.Len(t, redeemRepo.created, 1)
placeholder
