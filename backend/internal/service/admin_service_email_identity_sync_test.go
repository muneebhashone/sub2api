//go:build unit

package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ensureEmailCall struct {
	userID int64
	email  string
placeholder

type replaceEmailCall struct {
	userID   int64
	oldEmail string
	newEmail string
placeholder

type emailSyncUserRepoStub struct {
	*userRepoStub
	ensureCalls  []ensureEmailCall
	replaceCalls []replaceEmailCall
placeholder

func (s *emailSyncUserRepoStub) EnsureEmailAuthIdentity(_ context.Context, userID int64, email string) error {
	s.ensureCalls = append(s.ensureCalls, ensureEmailCall{userID: userID, email: emailplaceholder)
	return nil
placeholder

func (s *emailSyncUserRepoStub) ReplaceEmailAuthIdentity(_ context.Context, userID int64, oldEmail, newEmail string) error {
	s.replaceCalls = append(s.replaceCalls, replaceEmailCall{
		userID:   userID,
		oldEmail: oldEmail,
		newEmail: newEmail,
placeholder)
	return nil
placeholder

func (s *emailSyncUserRepoStub) GetLatestUsedAtByUserIDs(context.Context, []int64) (map[int64]*time.Time, error) {
	return map[int64]*time.Time{placeholder, nil
placeholder

func (s *emailSyncUserRepoStub) GetLatestUsedAtByUserID(context.Context, int64) (*time.Time, error) {
	return nil, nil
placeholder

func TestAdminService_CreateUser_EnsuresEmailAuthIdentity(t *testing.T) {
	repo := &emailSyncUserRepoStub{userRepoStub: &userRepoStub{nextID: 55placeholderplaceholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	user, err := svc.CreateUser(context.Background(), &CreateUserInput{
		Email:    "admin-created@example.com",
		Password: "strong-pass",
placeholder)
placeholder
	require.NotNil(t, user)
	require.Equal(t, []ensureEmailCall{{
		userID: 55,
		email:  "admin-created@example.com",
placeholderplaceholder, repo.ensureCalls)
	require.Empty(t, repo.replaceCalls)
placeholder

func TestAdminService_UpdateUser_ReplacesEmailAuthIdentity(t *testing.T) {
	repo := &emailSyncUserRepoStub{
		userRepoStub: &userRepoStub{
			user: &User{
				ID:          91,
				Email:       "before@example.com",
				Role:        RoleUser,
				Status:      StatusActive,
				Concurrency: 3,
		placeholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	updated, err := svc.UpdateUser(context.Background(), 91, &UpdateUserInput{
		Email: "after@example.com",
placeholder)
placeholder
	require.NotNil(t, updated)
	require.Equal(t, "after@example.com", updated.Email)
	require.Equal(t, []replaceEmailCall{{
		userID:   91,
		oldEmail: "before@example.com",
		newEmail: "after@example.com",
placeholderplaceholder, repo.replaceCalls)
	require.Empty(t, repo.ensureCalls)
placeholder
