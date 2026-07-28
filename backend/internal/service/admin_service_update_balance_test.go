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
	adjustErr error
	// changes 记录每次原子余额变更，顺序与调用顺序一致。
	changes []BalanceChange
placeholder

func (s *balanceUserRepoStub) AdjustBalance(ctx context.Context, id int64, delta float64) (BalanceChange, error) {
	return s.apply(func(current float64) float64 { return current + delta placeholder)
placeholder

func (s *balanceUserRepoStub) SetBalance(ctx context.Context, id int64, value float64) (BalanceChange, error) {
	return s.apply(func(float64) float64 { return value placeholder)
placeholder

func (s *balanceUserRepoStub) apply(next func(current float64) float64) (BalanceChange, error) {
	if s.adjustErr != nil {
		return BalanceChange{placeholder, s.adjustErr
placeholder
	if s.userRepoStub == nil || s.userRepoStub.user == nil {
		return BalanceChange{placeholder, ErrUserNotFound
placeholder
	change := BalanceChange{Old: s.userRepoStub.user.Balanceplaceholder
	change.New = next(change.Old)
	if change.New < 0 {
		return change, ErrBalanceNegative
placeholder
	s.userRepoStub.user.Balance = change.New
	s.changes = append(s.changes, change)
	return change, nil
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

// 管理员调账必须走原子的 AdjustBalance/SetBalance，而不是"读余额→算新值→整行写回"，
// 后者会把并发的计费扣款覆盖掉。userRepoStub.Update 对未预期的调用会 panic，
// 因此这里同时证明它没被走到。
func TestAdminService_UpdateUserBalance_UsesAtomicPrimitives(t *testing.T) {
	tests := []struct {
		name      string
		operation string
		amount    float64
		want      BalanceChange
placeholder{
		{name: "add", operation: "add", amount: 5, want: BalanceChange{Old: 10, New: 15placeholderplaceholder,
		{name: "subtract", operation: "subtract", amount: 4, want: BalanceChange{Old: 10, New: 6placeholderplaceholder,
		{name: "set", operation: "set", amount: 2, want: BalanceChange{Old: 10, New: 2placeholderplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &balanceUserRepoStub{userRepoStub: &userRepoStub{user: &User{ID: 7, Balance: 10placeholderplaceholderplaceholder
			svc := &adminServiceImpl{
				userRepo:       repo,
				redeemCodeRepo: &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder,
		placeholder

			user, err := svc.UpdateUserBalance(context.Background(), 7, tt.amount, tt.operation, "")
		placeholder
			require.Equal(t, []BalanceChange{tt.wantplaceholder, repo.changes)
			require.Equal(t, tt.want.New, user.Balance)
	placeholder)
placeholder
placeholder

func TestAdminService_UpdateUserBalance_RejectsNegativeResult(t *testing.T) {
	repo := &balanceUserRepoStub{userRepoStub: &userRepoStub{user: &User{ID: 7, Balance: 3placeholderplaceholderplaceholder
	svc := &adminServiceImpl{
		userRepo:       repo,
		redeemCodeRepo: &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder,
placeholder

	_, err := svc.UpdateUserBalance(context.Background(), 7, 4, "subtract", "")
placeholder
	require.Contains(t, err.Error(), "balance cannot be negative")
	require.Empty(t, repo.changes, "refused adjustment must not be applied")
	require.Equal(t, 3.0, repo.userRepoStub.user.Balance)
placeholder

func TestAdminService_UpdateUserBalance_RejectsUnknownOperation(t *testing.T) {
	repo := &balanceUserRepoStub{userRepoStub: &userRepoStub{user: &User{ID: 7, Balance: 10placeholderplaceholderplaceholder
	svc := &adminServiceImpl{
		userRepo:       repo,
		redeemCodeRepo: &balanceRedeemRepoStub{redeemRepoStub: &redeemRepoStub{placeholderplaceholder,
placeholder

	_, err := svc.UpdateUserBalance(context.Background(), 7, 1, "multiply", "")
placeholder
	require.Empty(t, repo.changes)
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
