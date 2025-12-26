//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserRepoSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	repo *userRepository
placeholder

func (s *UserRepoSuite) SetupTest() {
	s.ctx = context.Background()
	s.db = testTx(s.T())
	s.repo = NewUserRepository(s.db).(*userRepository)
placeholder

func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuite))
placeholder

// --- Create / GetByID / GetByEmail / Update / Delete ---

func (s *UserRepoSuite) TestCreate() {
	user := &service.User{
		Email:        "create@test.com",
		Username:     "testuser",
		PasswordHash: "test-password-hash",
		Role:         service.RoleUser,
		Status:       service.StatusActive,
placeholder

	err := s.repo.Create(s.ctx, user)
	s.Require().NoError(err, "Create")
	s.Require().NotZero(user.ID, "expected ID to be set")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err, "GetByID")
	s.Require().Equal("create@test.com", got.Email)
placeholder

func (s *UserRepoSuite) TestGetByID_NotFound() {
	_, err := s.repo.GetByID(s.ctx, 999999)
	s.Require().Error(err, "expected error for non-existent ID")
placeholder

func (s *UserRepoSuite) TestGetByEmail() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "byemail@test.com"placeholder)

	got, err := s.repo.GetByEmail(s.ctx, user.Email)
	s.Require().NoError(err, "GetByEmail")
	s.Require().Equal(user.ID, got.ID)
placeholder

func (s *UserRepoSuite) TestGetByEmail_NotFound() {
	_, err := s.repo.GetByEmail(s.ctx, "nonexistent@test.com")
	s.Require().Error(err, "expected error for non-existent email")
placeholder

func (s *UserRepoSuite) TestUpdate() {
	user := userModelToService(mustCreateUser(s.T(), s.db, &userModel{Email: "update@test.com", Username: "original"placeholder))

	user.Username = "updated"
	err := s.repo.Update(s.ctx, user)
	s.Require().NoError(err, "Update")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err, "GetByID after update")
	s.Require().Equal("updated", got.Username)
placeholder

func (s *UserRepoSuite) TestDelete() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "delete@test.com"placeholder)

	err := s.repo.Delete(s.ctx, user.ID)
	s.Require().NoError(err, "Delete")

	_, err = s.repo.GetByID(s.ctx, user.ID)
	s.Require().Error(err, "expected error after delete")
placeholder

// --- List / ListWithFilters ---

func (s *UserRepoSuite) TestList() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "list1@test.com"placeholder)
	mustCreateUser(s.T(), s.db, &userModel{Email: "list2@test.com"placeholder)

	users, page, err := s.repo.List(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder)
	s.Require().NoError(err, "List")
	s.Require().Len(users, 2)
	s.Require().Equal(int64(2), page.Total)
placeholder

func (s *UserRepoSuite) TestListWithFilters_Status() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "active@test.com", Status: service.StatusActiveplaceholder)
	mustCreateUser(s.T(), s.db, &userModel{Email: "disabled@test.com", Status: service.StatusDisabledplaceholder)

	users, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, service.StatusActive, "", "")
	s.Require().NoError(err)
	s.Require().Len(users, 1)
	s.Require().Equal(service.StatusActive, users[0].Status)
placeholder

func (s *UserRepoSuite) TestListWithFilters_Role() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "user@test.com", Role: service.RoleUserplaceholder)
	mustCreateUser(s.T(), s.db, &userModel{Email: "admin@test.com", Role: service.RoleAdminplaceholder)

	users, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", service.RoleAdmin, "")
	s.Require().NoError(err)
	s.Require().Len(users, 1)
	s.Require().Equal(service.RoleAdmin, users[0].Role)
placeholder

func (s *UserRepoSuite) TestListWithFilters_Search() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "alice@test.com", Username: "Alice"placeholder)
	mustCreateUser(s.T(), s.db, &userModel{Email: "bob@test.com", Username: "Bob"placeholder)

	users, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "alice")
	s.Require().NoError(err)
	s.Require().Len(users, 1)
	s.Require().Contains(users[0].Email, "alice")
placeholder

func (s *UserRepoSuite) TestListWithFilters_SearchByUsername() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "u1@test.com", Username: "JohnDoe"placeholder)
	mustCreateUser(s.T(), s.db, &userModel{Email: "u2@test.com", Username: "JaneSmith"placeholder)

	users, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "john")
	s.Require().NoError(err)
	s.Require().Len(users, 1)
	s.Require().Equal("JohnDoe", users[0].Username)
placeholder

func (s *UserRepoSuite) TestListWithFilters_SearchByWechat() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "w1@test.com", Wechat: "wx_hello"placeholder)
	mustCreateUser(s.T(), s.db, &userModel{Email: "w2@test.com", Wechat: "wx_world"placeholder)

	users, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "wx_hello")
	s.Require().NoError(err)
	s.Require().Len(users, 1)
	s.Require().Equal("wx_hello", users[0].Wechat)
placeholder

func (s *UserRepoSuite) TestListWithFilters_LoadsActiveSubscriptions() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "sub@test.com", Status: service.StatusActiveplaceholder)
	group := mustCreateGroup(s.T(), s.db, &groupModel{Name: "g-sub"placeholder)

	_ = mustCreateSubscription(s.T(), s.db, &userSubscriptionModel{
		UserID:    user.ID,
		GroupID:   group.ID,
		Status:    service.SubscriptionStatusActive,
		ExpiresAt: time.Now().Add(1 * time.Hour),
placeholder)
	_ = mustCreateSubscription(s.T(), s.db, &userSubscriptionModel{
		UserID:    user.ID,
		GroupID:   group.ID,
		Status:    service.SubscriptionStatusExpired,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
placeholder)

	users, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "sub@")
	s.Require().NoError(err, "ListWithFilters")
	s.Require().Len(users, 1, "expected 1 user")
	s.Require().Len(users[0].Subscriptions, 1, "expected 1 active subscription")
	s.Require().NotNil(users[0].Subscriptions[0].Group, "expected subscription group preload")
	s.Require().Equal(group.ID, users[0].Subscriptions[0].Group.ID, "group ID mismatch")
placeholder

func (s *UserRepoSuite) TestListWithFilters_CombinedFilters() {
	mustCreateUser(s.T(), s.db, &userModel{
		Email:    "a@example.com",
		Username: "Alice",
		Wechat:   "wx_a",
		Role:     service.RoleUser,
		Status:   service.StatusActive,
		Balance:  10,
placeholder)
	target := mustCreateUser(s.T(), s.db, &userModel{
		Email:    "b@example.com",
		Username: "Bob",
		Wechat:   "wx_b",
		Role:     service.RoleAdmin,
		Status:   service.StatusActive,
		Balance:  1,
placeholder)
	mustCreateUser(s.T(), s.db, &userModel{
		Email:  "c@example.com",
		Role:   service.RoleAdmin,
		Status: service.StatusDisabled,
placeholder)

	users, page, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, service.StatusActive, service.RoleAdmin, "b@")
	s.Require().NoError(err, "ListWithFilters")
	s.Require().Equal(int64(1), page.Total, "ListWithFilters total mismatch")
	s.Require().Len(users, 1, "ListWithFilters len mismatch")
	s.Require().Equal(target.ID, users[0].ID, "ListWithFilters result mismatch")
placeholder

// --- Balance operations ---

func (s *UserRepoSuite) TestUpdateBalance() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "bal@test.com", Balance: 10placeholder)

	err := s.repo.UpdateBalance(s.ctx, user.ID, 2.5)
	s.Require().NoError(err, "UpdateBalance")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Equal(12.5, got.Balance)
placeholder

func (s *UserRepoSuite) TestUpdateBalance_Negative() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "balneg@test.com", Balance: 10placeholder)

	err := s.repo.UpdateBalance(s.ctx, user.ID, -3)
	s.Require().NoError(err, "UpdateBalance with negative")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Equal(7.0, got.Balance)
placeholder

func (s *UserRepoSuite) TestDeductBalance() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "deduct@test.com", Balance: 10placeholder)

	err := s.repo.DeductBalance(s.ctx, user.ID, 5)
	s.Require().NoError(err, "DeductBalance")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Equal(5.0, got.Balance)
placeholder

func (s *UserRepoSuite) TestDeductBalance_InsufficientFunds() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "insuf@test.com", Balance: 5placeholder)

	err := s.repo.DeductBalance(s.ctx, user.ID, 999)
	s.Require().Error(err, "expected error for insufficient balance")
	s.Require().ErrorIs(err, service.ErrInsufficientBalance)
placeholder

func (s *UserRepoSuite) TestDeductBalance_ExactAmount() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "exact@test.com", Balance: 10placeholder)

	err := s.repo.DeductBalance(s.ctx, user.ID, 10)
	s.Require().NoError(err, "DeductBalance exact amount")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Zero(got.Balance)
placeholder

// --- Concurrency ---

func (s *UserRepoSuite) TestUpdateConcurrency() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "conc@test.com", Concurrency: 5placeholder)

	err := s.repo.UpdateConcurrency(s.ctx, user.ID, 3)
	s.Require().NoError(err, "UpdateConcurrency")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Equal(8, got.Concurrency)
placeholder

func (s *UserRepoSuite) TestUpdateConcurrency_Negative() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "concneg@test.com", Concurrency: 5placeholder)

	err := s.repo.UpdateConcurrency(s.ctx, user.ID, -2)
	s.Require().NoError(err, "UpdateConcurrency negative")

	got, err := s.repo.GetByID(s.ctx, user.ID)
	s.Require().NoError(err)
	s.Require().Equal(3, got.Concurrency)
placeholder

// --- ExistsByEmail ---

func (s *UserRepoSuite) TestExistsByEmail() {
	mustCreateUser(s.T(), s.db, &userModel{Email: "exists@test.com"placeholder)

	exists, err := s.repo.ExistsByEmail(s.ctx, "exists@test.com")
	s.Require().NoError(err, "ExistsByEmail")
	s.Require().True(exists)

	notExists, err := s.repo.ExistsByEmail(s.ctx, "notexists@test.com")
	s.Require().NoError(err)
	s.Require().False(notExists)
placeholder

// --- RemoveGroupFromAllowedGroups ---

func (s *UserRepoSuite) TestRemoveGroupFromAllowedGroups() {
	groupID := int64(42)
	userA := mustCreateUser(s.T(), s.db, &userModel{
		Email:         "a1@example.com",
		AllowedGroups: pq.Int64Array{groupID, 7placeholder,
placeholder)
	mustCreateUser(s.T(), s.db, &userModel{
		Email:         "a2@example.com",
		AllowedGroups: pq.Int64Array{7placeholder,
placeholder)

	affected, err := s.repo.RemoveGroupFromAllowedGroups(s.ctx, groupID)
	s.Require().NoError(err, "RemoveGroupFromAllowedGroups")
	s.Require().Equal(int64(1), affected, "expected 1 affected row")

	got, err := s.repo.GetByID(s.ctx, userA.ID)
	s.Require().NoError(err, "GetByID")
	for _, id := range got.AllowedGroups {
		s.Require().NotEqual(groupID, id, "expected groupID to be removed from allowed_groups")
placeholder
placeholder

func (s *UserRepoSuite) TestRemoveGroupFromAllowedGroups_NoMatch() {
	mustCreateUser(s.T(), s.db, &userModel{
		Email:         "nomatch@test.com",
		AllowedGroups: pq.Int64Array{1, 2, 3placeholder,
placeholder)

	affected, err := s.repo.RemoveGroupFromAllowedGroups(s.ctx, 999)
	s.Require().NoError(err)
	s.Require().Zero(affected, "expected no affected rows")
placeholder

// --- GetFirstAdmin ---

func (s *UserRepoSuite) TestGetFirstAdmin() {
	admin1 := mustCreateUser(s.T(), s.db, &userModel{
		Email:  "admin1@example.com",
		Role:   service.RoleAdmin,
		Status: service.StatusActive,
placeholder)
	mustCreateUser(s.T(), s.db, &userModel{
		Email:  "admin2@example.com",
		Role:   service.RoleAdmin,
		Status: service.StatusActive,
placeholder)

	got, err := s.repo.GetFirstAdmin(s.ctx)
	s.Require().NoError(err, "GetFirstAdmin")
	s.Require().Equal(admin1.ID, got.ID, "GetFirstAdmin mismatch")
placeholder

func (s *UserRepoSuite) TestGetFirstAdmin_NoAdmin() {
	mustCreateUser(s.T(), s.db, &userModel{
		Email:  "user@example.com",
		Role:   service.RoleUser,
		Status: service.StatusActive,
placeholder)

	_, err := s.repo.GetFirstAdmin(s.ctx)
	s.Require().Error(err, "expected error when no admin exists")
placeholder

func (s *UserRepoSuite) TestGetFirstAdmin_DisabledAdminIgnored() {
	mustCreateUser(s.T(), s.db, &userModel{
		Email:  "disabled@example.com",
		Role:   service.RoleAdmin,
		Status: service.StatusDisabled,
placeholder)
	activeAdmin := mustCreateUser(s.T(), s.db, &userModel{
		Email:  "active@example.com",
		Role:   service.RoleAdmin,
		Status: service.StatusActive,
placeholder)

	got, err := s.repo.GetFirstAdmin(s.ctx)
	s.Require().NoError(err, "GetFirstAdmin")
	s.Require().Equal(activeAdmin.ID, got.ID, "should return only active admin")
placeholder

// --- Combined original test ---

func (s *UserRepoSuite) TestCRUD_And_Filters_And_AtomicUpdates() {
	user1 := mustCreateUser(s.T(), s.db, &userModel{
		Email:    "a@example.com",
		Username: "Alice",
		Wechat:   "wx_a",
		Role:     service.RoleUser,
		Status:   service.StatusActive,
		Balance:  10,
placeholder)
	user2 := mustCreateUser(s.T(), s.db, &userModel{
		Email:    "b@example.com",
		Username: "Bob",
		Wechat:   "wx_b",
		Role:     service.RoleAdmin,
		Status:   service.StatusActive,
		Balance:  1,
placeholder)
	_ = mustCreateUser(s.T(), s.db, &userModel{
		Email:  "c@example.com",
		Role:   service.RoleAdmin,
		Status: service.StatusDisabled,
placeholder)

	got, err := s.repo.GetByID(s.ctx, user1.ID)
	s.Require().NoError(err, "GetByID")
	s.Require().Equal(user1.Email, got.Email, "GetByID email mismatch")

	gotByEmail, err := s.repo.GetByEmail(s.ctx, user2.Email)
	s.Require().NoError(err, "GetByEmail")
	s.Require().Equal(user2.ID, gotByEmail.ID, "GetByEmail ID mismatch")

	got.Username = "Alice2"
	s.Require().NoError(s.repo.Update(s.ctx, got), "Update")
	got2, err := s.repo.GetByID(s.ctx, user1.ID)
	s.Require().NoError(err, "GetByID after update")
	s.Require().Equal("Alice2", got2.Username, "Update did not persist")

	s.Require().NoError(s.repo.UpdateBalance(s.ctx, user1.ID, 2.5), "UpdateBalance")
	got3, err := s.repo.GetByID(s.ctx, user1.ID)
	s.Require().NoError(err, "GetByID after UpdateBalance")
	s.Require().Equal(12.5, got3.Balance, "UpdateBalance mismatch")

	s.Require().NoError(s.repo.DeductBalance(s.ctx, user1.ID, 5), "DeductBalance")
	got4, err := s.repo.GetByID(s.ctx, user1.ID)
	s.Require().NoError(err, "GetByID after DeductBalance")
	s.Require().Equal(7.5, got4.Balance, "DeductBalance mismatch")

	err = s.repo.DeductBalance(s.ctx, user1.ID, 999)
	s.Require().Error(err, "DeductBalance expected error for insufficient balance")
	s.Require().ErrorIs(err, service.ErrInsufficientBalance, "DeductBalance unexpected error")

	s.Require().NoError(s.repo.UpdateConcurrency(s.ctx, user1.ID, 3), "UpdateConcurrency")
	got5, err := s.repo.GetByID(s.ctx, user1.ID)
	s.Require().NoError(err, "GetByID after UpdateConcurrency")
	s.Require().Equal(user1.Concurrency+3, got5.Concurrency, "UpdateConcurrency mismatch")

	params := pagination.PaginationParams{Page: 1, PageSize: 10placeholder
	users, page, err := s.repo.ListWithFilters(s.ctx, params, service.StatusActive, service.RoleAdmin, "b@")
	s.Require().NoError(err, "ListWithFilters")
	s.Require().Equal(int64(1), page.Total, "ListWithFilters total mismatch")
	s.Require().Len(users, 1, "ListWithFilters len mismatch")
	s.Require().Equal(user2.ID, users[0].ID, "ListWithFilters result mismatch")
placeholder
