//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func TestAdminService_CreateUser_WithAdminRole(t *testing.T) {
	repo := &userRepoStub{nextID: 30placeholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	user, err := svc.CreateUser(context.Background(), &CreateUserInput{
		Email:    "admin@test.com",
		Password: "strong-pass",
		Role:     RoleAdmin,
placeholder)
placeholder
	require.Equal(t, RoleAdmin, user.Role)
placeholder

func TestAdminService_CreateUser_DefaultsToUserRole(t *testing.T) {
	repo := &userRepoStub{nextID: 31placeholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	user, err := svc.CreateUser(context.Background(), &CreateUserInput{
		Email:    "plain@test.com",
		Password: "strong-pass",
placeholder)
placeholder
	require.Equal(t, RoleUser, user.Role)
placeholder

func TestAdminService_CreateUser_InvalidRoleRejected(t *testing.T) {
	repo := &userRepoStub{nextID: 32placeholder
	svc := &adminServiceImpl{userRepo: repoplaceholder

	_, err := svc.CreateUser(context.Background(), &CreateUserInput{
		Email:    "bad@test.com",
		Password: "strong-pass",
		Role:     "superuser",
placeholder)
placeholder
	require.Empty(t, repo.created, "非法角色不应写入用户")
placeholder

func TestAdminService_UpdateUser_PromoteToAdmin(t *testing.T) {
	base := &userRepoStub{user: &User{ID: 42, Email: "u@example.com", Role: RoleUserplaceholderplaceholder
	repo := &rpmUserRepoStub{userRepoStub: baseplaceholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		userRepo:             repo,
		redeemCodeRepo:       &redeemRepoStub{placeholder,
		authCacheInvalidator: invalidator,
placeholder

	updated, err := svc.UpdateUser(context.Background(), 42, &UpdateUserInput{Role: RoleAdminplaceholder)
placeholder
	require.Equal(t, RoleAdmin, updated.Role)
	require.Equal(t, []int64{42placeholder, invalidator.userIDs, "角色变更应失效认证缓存")
placeholder

func TestAdminService_UpdateUser_RoleOmittedKeepsExisting(t *testing.T) {
	base := &userRepoStub{user: &User{ID: 42, Email: "u@example.com", Role: RoleAdminplaceholderplaceholder
	repo := &rpmUserRepoStub{userRepoStub: baseplaceholder
	svc := &adminServiceImpl{userRepo: repo, redeemCodeRepo: &redeemRepoStub{placeholderplaceholder

	newName := "renamed"
	updated, err := svc.UpdateUser(context.Background(), 42, &UpdateUserInput{Username: &newNameplaceholder)
placeholder
	require.Equal(t, RoleAdmin, updated.Role, "未提供 role 时不应改变现有角色")
placeholder

func TestAdminService_UpdateUser_InvalidRoleRejected(t *testing.T) {
	base := &userRepoStub{user: &User{ID: 42, Email: "u@example.com", Role: RoleUserplaceholderplaceholder
	repo := &rpmUserRepoStub{userRepoStub: baseplaceholder
	svc := &adminServiceImpl{userRepo: repo, redeemCodeRepo: &redeemRepoStub{placeholderplaceholder

	_, err := svc.UpdateUser(context.Background(), 42, &UpdateUserInput{Role: "root"placeholder)
placeholder
	require.Nil(t, repo.lastUpdated, "非法角色不应触发持久化")
placeholder

// roleGuardUserRepoStub 在 rpmUserRepoStub 之上提供可控的管理员计数，
// 用于测试"最后一个管理员不可降级"守卫。
type roleGuardUserRepoStub struct {
	*rpmUserRepoStub
	adminTotal int64
	listCalls  int
placeholder

func (s *roleGuardUserRepoStub) ListWithFilters(_ context.Context, _ pagination.PaginationParams, _ UserListFilters) ([]User, *pagination.PaginationResult, error) {
	s.listCalls++
	return nil, &pagination.PaginationResult{Total: s.adminTotalplaceholder, nil
placeholder

func TestAdminService_UpdateUser_DemoteLastAdminRejected(t *testing.T) {
	base := &userRepoStub{user: &User{ID: 42, Email: "a@example.com", Role: RoleAdminplaceholderplaceholder
	repo := &roleGuardUserRepoStub{rpmUserRepoStub: &rpmUserRepoStub{userRepoStub: baseplaceholder, adminTotal: 1placeholder
	svc := &adminServiceImpl{userRepo: repo, redeemCodeRepo: &redeemRepoStub{placeholderplaceholder

	_, err := svc.UpdateUser(context.Background(), 42, &UpdateUserInput{Role: RoleUserplaceholder)
placeholder
	require.Contains(t, err.Error(), "last admin")
	require.Nil(t, repo.lastUpdated, "最后一个管理员不应被降级持久化")
	require.Equal(t, 1, repo.listCalls, "降级路径应触发管理员计数")
placeholder

func TestAdminService_UpdateUser_DemoteAdminAllowedWhenOthersExist(t *testing.T) {
	base := &userRepoStub{user: &User{ID: 42, Email: "a@example.com", Role: RoleAdminplaceholderplaceholder
	repo := &roleGuardUserRepoStub{rpmUserRepoStub: &rpmUserRepoStub{userRepoStub: baseplaceholder, adminTotal: 2placeholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		userRepo:             repo,
		redeemCodeRepo:       &redeemRepoStub{placeholder,
		authCacheInvalidator: invalidator,
placeholder

	updated, err := svc.UpdateUser(context.Background(), 42, &UpdateUserInput{Role: RoleUserplaceholder)
placeholder
	require.Equal(t, RoleUser, updated.Role)
	require.NotNil(t, repo.lastUpdated)
	require.Equal(t, RoleUser, repo.lastUpdated.Role, "存在其他管理员时允许降级")
placeholder

func TestAdminService_UpdateUser_PromoteDoesNotCountAdmins(t *testing.T) {
	base := &userRepoStub{user: &User{ID: 42, Email: "u@example.com", Role: RoleUserplaceholderplaceholder
	repo := &roleGuardUserRepoStub{rpmUserRepoStub: &rpmUserRepoStub{userRepoStub: baseplaceholder, adminTotal: 1placeholder
	svc := &adminServiceImpl{
		userRepo:             repo,
		redeemCodeRepo:       &redeemRepoStub{placeholder,
		authCacheInvalidator: &authCacheInvalidatorStub{placeholder,
placeholder

	updated, err := svc.UpdateUser(context.Background(), 42, &UpdateUserInput{Role: RoleAdminplaceholder)
placeholder
	require.Equal(t, RoleAdmin, updated.Role)
	require.Equal(t, 0, repo.listCalls, "升级路径不应触发管理员计数")
placeholder
