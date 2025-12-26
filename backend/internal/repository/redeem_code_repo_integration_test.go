//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RedeemCodeRepoSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	repo *redeemCodeRepository
placeholder

func (s *RedeemCodeRepoSuite) SetupTest() {
	s.ctx = context.Background()
	s.db = testTx(s.T())
	s.repo = NewRedeemCodeRepository(s.db).(*redeemCodeRepository)
placeholder

func TestRedeemCodeRepoSuite(t *testing.T) {
	suite.Run(t, new(RedeemCodeRepoSuite))
placeholder

// --- Create / CreateBatch / GetByID / GetByCode ---

func (s *RedeemCodeRepoSuite) TestCreate() {
	code := &service.RedeemCode{
		Code:   "TEST-CREATE",
		Type:   service.RedeemTypeBalance,
		Value:  100,
		Status: service.StatusUnused,
placeholder

	err := s.repo.Create(s.ctx, code)
	s.Require().NoError(err, "Create")
	s.Require().NotZero(code.ID, "expected ID to be set")

	got, err := s.repo.GetByID(s.ctx, code.ID)
	s.Require().NoError(err, "GetByID")
	s.Require().Equal("TEST-CREATE", got.Code)
placeholder

func (s *RedeemCodeRepoSuite) TestCreateBatch() {
	codes := []service.RedeemCode{
		{Code: "BATCH-1", Type: service.RedeemTypeBalance, Value: 10, Status: service.StatusUnusedplaceholder,
		{Code: "BATCH-2", Type: service.RedeemTypeBalance, Value: 20, Status: service.StatusUnusedplaceholder,
placeholder

	err := s.repo.CreateBatch(s.ctx, codes)
	s.Require().NoError(err, "CreateBatch")

	got1, err := s.repo.GetByCode(s.ctx, "BATCH-1")
	s.Require().NoError(err)
	s.Require().Equal(float64(10), got1.Value)

	got2, err := s.repo.GetByCode(s.ctx, "BATCH-2")
	s.Require().NoError(err)
	s.Require().Equal(float64(20), got2.Value)
placeholder

func (s *RedeemCodeRepoSuite) TestGetByID_NotFound() {
	_, err := s.repo.GetByID(s.ctx, 999999)
	s.Require().Error(err, "expected error for non-existent ID")
placeholder

func (s *RedeemCodeRepoSuite) TestGetByCode() {
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "GET-BY-CODE", Type: service.RedeemTypeBalanceplaceholder)

	got, err := s.repo.GetByCode(s.ctx, "GET-BY-CODE")
	s.Require().NoError(err, "GetByCode")
	s.Require().Equal("GET-BY-CODE", got.Code)
placeholder

func (s *RedeemCodeRepoSuite) TestGetByCode_NotFound() {
	_, err := s.repo.GetByCode(s.ctx, "NON-EXISTENT")
	s.Require().Error(err, "expected error for non-existent code")
placeholder

// --- Delete ---

func (s *RedeemCodeRepoSuite) TestDelete() {
	code := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "TO-DELETE", Type: service.RedeemTypeBalanceplaceholder)

	err := s.repo.Delete(s.ctx, code.ID)
	s.Require().NoError(err, "Delete")

	_, err = s.repo.GetByID(s.ctx, code.ID)
	s.Require().Error(err, "expected error after delete")
placeholder

// --- List / ListWithFilters ---

func (s *RedeemCodeRepoSuite) TestList() {
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "LIST-1", Type: service.RedeemTypeBalanceplaceholder)
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "LIST-2", Type: service.RedeemTypeBalanceplaceholder)

	codes, page, err := s.repo.List(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder)
	s.Require().NoError(err, "List")
	s.Require().Len(codes, 2)
	s.Require().Equal(int64(2), page.Total)
placeholder

func (s *RedeemCodeRepoSuite) TestListWithFilters_Type() {
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "TYPE-BAL", Type: service.RedeemTypeBalanceplaceholder)
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "TYPE-SUB", Type: service.RedeemTypeSubscriptionplaceholder)

	codes, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, service.RedeemTypeSubscription, "", "")
	s.Require().NoError(err)
	s.Require().Len(codes, 1)
	s.Require().Equal(service.RedeemTypeSubscription, codes[0].Type)
placeholder

func (s *RedeemCodeRepoSuite) TestListWithFilters_Status() {
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "STAT-UNUSED", Type: service.RedeemTypeBalance, Status: service.StatusUnusedplaceholder)
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "STAT-USED", Type: service.RedeemTypeBalance, Status: service.StatusUsedplaceholder)

	codes, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", service.StatusUsed, "")
	s.Require().NoError(err)
	s.Require().Len(codes, 1)
	s.Require().Equal(service.StatusUsed, codes[0].Status)
placeholder

func (s *RedeemCodeRepoSuite) TestListWithFilters_Search() {
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "ALPHA-CODE", Type: service.RedeemTypeBalanceplaceholder)
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "BETA-CODE", Type: service.RedeemTypeBalanceplaceholder)

	codes, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "alpha")
	s.Require().NoError(err)
	s.Require().Len(codes, 1)
	s.Require().Contains(codes[0].Code, "ALPHA")
placeholder

func (s *RedeemCodeRepoSuite) TestListWithFilters_GroupPreload() {
	group := mustCreateGroup(s.T(), s.db, &groupModel{Name: "g-preload"placeholder)
	mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{
		Code:    "WITH-GROUP",
		Type:    service.RedeemTypeSubscription,
		GroupID: &group.ID,
placeholder)

	codes, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "")
	s.Require().NoError(err)
	s.Require().Len(codes, 1)
	s.Require().NotNil(codes[0].Group, "expected Group preload")
	s.Require().Equal(group.ID, codes[0].Group.ID)
placeholder

// --- Update ---

func (s *RedeemCodeRepoSuite) TestUpdate() {
	code := redeemCodeModelToService(mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "UPDATE-ME", Type: service.RedeemTypeBalance, Value: 10placeholder))

	code.Value = 50
	err := s.repo.Update(s.ctx, code)
	s.Require().NoError(err, "Update")

	got, err := s.repo.GetByID(s.ctx, code.ID)
	s.Require().NoError(err)
	s.Require().Equal(float64(50), got.Value)
placeholder

// --- Use ---

func (s *RedeemCodeRepoSuite) TestUse() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "use@test.com"placeholder)
	code := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "USE-ME", Type: service.RedeemTypeBalance, Status: service.StatusUnusedplaceholder)

	err := s.repo.Use(s.ctx, code.ID, user.ID)
	s.Require().NoError(err, "Use")

	got, err := s.repo.GetByID(s.ctx, code.ID)
	s.Require().NoError(err)
	s.Require().Equal(service.StatusUsed, got.Status)
	s.Require().NotNil(got.UsedBy)
	s.Require().Equal(user.ID, *got.UsedBy)
	s.Require().NotNil(got.UsedAt)
placeholder

func (s *RedeemCodeRepoSuite) TestUse_Idempotency() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "idem@test.com"placeholder)
	code := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "IDEM-CODE", Type: service.RedeemTypeBalance, Status: service.StatusUnusedplaceholder)

	err := s.repo.Use(s.ctx, code.ID, user.ID)
	s.Require().NoError(err, "Use first time")

	// Second use should fail
	err = s.repo.Use(s.ctx, code.ID, user.ID)
	s.Require().Error(err, "Use expected error on second call")
	s.Require().ErrorIs(err, service.ErrRedeemCodeUsed)
placeholder

func (s *RedeemCodeRepoSuite) TestUse_AlreadyUsed() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "already@test.com"placeholder)
	code := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{Code: "ALREADY-USED", Type: service.RedeemTypeBalance, Status: service.StatusUsedplaceholder)

	err := s.repo.Use(s.ctx, code.ID, user.ID)
	s.Require().Error(err, "expected error for already used code")
	s.Require().ErrorIs(err, service.ErrRedeemCodeUsed)
placeholder

// --- ListByUser ---

func (s *RedeemCodeRepoSuite) TestListByUser() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "listby@test.com"placeholder)
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	// Create codes with explicit used_at for ordering
	c1 := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{
		Code:   "USER-1",
		Type:   service.RedeemTypeBalance,
		Status: service.StatusUsed,
		UsedBy: &user.ID,
placeholder)
	s.db.Model(c1).Update("used_at", base)

	c2 := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{
		Code:   "USER-2",
		Type:   service.RedeemTypeBalance,
		Status: service.StatusUsed,
		UsedBy: &user.ID,
placeholder)
	s.db.Model(c2).Update("used_at", base.Add(1*time.Hour))

	codes, err := s.repo.ListByUser(s.ctx, user.ID, 10)
	s.Require().NoError(err, "ListByUser")
	s.Require().Len(codes, 2)
	// Ordered by used_at DESC, so USER-2 first
	s.Require().Equal("USER-2", codes[0].Code)
	s.Require().Equal("USER-1", codes[1].Code)
placeholder

func (s *RedeemCodeRepoSuite) TestListByUser_WithGroupPreload() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "grp@test.com"placeholder)
	group := mustCreateGroup(s.T(), s.db, &groupModel{Name: "g-listby"placeholder)

	c := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{
		Code:    "WITH-GRP",
		Type:    service.RedeemTypeSubscription,
		Status:  service.StatusUsed,
		UsedBy:  &user.ID,
		GroupID: &group.ID,
placeholder)
	s.db.Model(c).Update("used_at", time.Now())

	codes, err := s.repo.ListByUser(s.ctx, user.ID, 10)
	s.Require().NoError(err)
	s.Require().Len(codes, 1)
	s.Require().NotNil(codes[0].Group)
	s.Require().Equal(group.ID, codes[0].Group.ID)
placeholder

func (s *RedeemCodeRepoSuite) TestListByUser_DefaultLimit() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "deflimit@test.com"placeholder)
	c := mustCreateRedeemCode(s.T(), s.db, &redeemCodeModel{
		Code:   "DEF-LIM",
		Type:   service.RedeemTypeBalance,
		Status: service.StatusUsed,
		UsedBy: &user.ID,
placeholder)
	s.db.Model(c).Update("used_at", time.Now())

	// limit <= 0 should default to 10
	codes, err := s.repo.ListByUser(s.ctx, user.ID, 0)
	s.Require().NoError(err)
	s.Require().Len(codes, 1)
placeholder

// --- Combined original test ---

func (s *RedeemCodeRepoSuite) TestCreateBatch_Filters_Use_Idempotency_ListByUser() {
	user := mustCreateUser(s.T(), s.db, &userModel{Email: "rc@example.com"placeholder)
	group := mustCreateGroup(s.T(), s.db, &groupModel{Name: "g-rc"placeholder)

	codes := []service.RedeemCode{
		{Code: "CODEA", Type: service.RedeemTypeBalance, Value: 1, Status: service.StatusUnused, CreatedAt: time.Now()placeholder,
		{Code: "CODEB", Type: service.RedeemTypeSubscription, Value: 0, Status: service.StatusUnused, GroupID: &group.ID, ValidityDays: 7, CreatedAt: time.Now()placeholder,
placeholder
	s.Require().NoError(s.repo.CreateBatch(s.ctx, codes), "CreateBatch")

	list, page, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, service.RedeemTypeSubscription, service.StatusUnused, "code")
	s.Require().NoError(err, "ListWithFilters")
	s.Require().Equal(int64(1), page.Total)
	s.Require().Len(list, 1)
	s.Require().NotNil(list[0].Group, "expected Group preload")
	s.Require().Equal(group.ID, list[0].Group.ID)

	codeB, err := s.repo.GetByCode(s.ctx, "CODEB")
	s.Require().NoError(err, "GetByCode")
	s.Require().NoError(s.repo.Use(s.ctx, codeB.ID, user.ID), "Use")
	err = s.repo.Use(s.ctx, codeB.ID, user.ID)
	s.Require().Error(err, "Use expected error on second call")
	s.Require().ErrorIs(err, service.ErrRedeemCodeUsed)

	codeA, err := s.repo.GetByCode(s.ctx, "CODEA")
	s.Require().NoError(err, "GetByCode")

	// Use fixed time instead of time.Sleep for deterministic ordering
	s.db.Model(&redeemCodeModel{placeholder).Where("id = ?", codeB.ID).Update("used_at", time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC))
	s.Require().NoError(s.repo.Use(s.ctx, codeA.ID, user.ID), "Use codeA")
	s.db.Model(&redeemCodeModel{placeholder).Where("id = ?", codeA.ID).Update("used_at", time.Date(2025, 1, 1, 13, 0, 0, 0, time.UTC))

	used, err := s.repo.ListByUser(s.ctx, user.ID, 10)
	s.Require().NoError(err, "ListByUser")
	s.Require().Len(used, 2, "expected 2 used codes")
	s.Require().Equal("CODEA", used[0].Code, "expected newest used code first")
placeholder
