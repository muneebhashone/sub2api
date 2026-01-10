//go:build unit

package service

import (
	"context"
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
