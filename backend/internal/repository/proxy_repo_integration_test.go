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

type ProxyRepoSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	repo *proxyRepository
placeholder

func (s *ProxyRepoSuite) SetupTest() {
	s.ctx = context.Background()
	s.db = testTx(s.T())
	s.repo = NewProxyRepository(s.db).(*proxyRepository)
placeholder

func TestProxyRepoSuite(t *testing.T) {
	suite.Run(t, new(ProxyRepoSuite))
placeholder

// --- Create / GetByID / Update / Delete ---

func (s *ProxyRepoSuite) TestCreate() {
	proxy := &service.Proxy{
		Name:     "test-create",
		Protocol: "http",
		Host:     "127.0.0.1",
		Port:     8080,
		Status:   service.StatusActive,
placeholder

	err := s.repo.Create(s.ctx, proxy)
	s.Require().NoError(err, "Create")
	s.Require().NotZero(proxy.ID, "expected ID to be set")

	got, err := s.repo.GetByID(s.ctx, proxy.ID)
	s.Require().NoError(err, "GetByID")
	s.Require().Equal("test-create", got.Name)
placeholder

func (s *ProxyRepoSuite) TestGetByID_NotFound() {
	_, err := s.repo.GetByID(s.ctx, 999999)
	s.Require().Error(err, "expected error for non-existent ID")
placeholder

func (s *ProxyRepoSuite) TestUpdate() {
	proxy := proxyModelToService(mustCreateProxy(s.T(), s.db, &proxyModel{Name: "original"placeholder))

	proxy.Name = "updated"
	err := s.repo.Update(s.ctx, proxy)
	s.Require().NoError(err, "Update")

	got, err := s.repo.GetByID(s.ctx, proxy.ID)
	s.Require().NoError(err, "GetByID after update")
	s.Require().Equal("updated", got.Name)
placeholder

func (s *ProxyRepoSuite) TestDelete() {
	proxy := mustCreateProxy(s.T(), s.db, &proxyModel{Name: "to-delete"placeholder)

	err := s.repo.Delete(s.ctx, proxy.ID)
	s.Require().NoError(err, "Delete")

	_, err = s.repo.GetByID(s.ctx, proxy.ID)
	s.Require().Error(err, "expected error after delete")
placeholder

// --- List / ListWithFilters ---

func (s *ProxyRepoSuite) TestList() {
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p1"placeholder)
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p2"placeholder)

	proxies, page, err := s.repo.List(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder)
	s.Require().NoError(err, "List")
	s.Require().Len(proxies, 2)
	s.Require().Equal(int64(2), page.Total)
placeholder

func (s *ProxyRepoSuite) TestListWithFilters_Protocol() {
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p1", Protocol: "http"placeholder)
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p2", Protocol: "socks5"placeholder)

	proxies, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "socks5", "", "")
	s.Require().NoError(err)
	s.Require().Len(proxies, 1)
	s.Require().Equal("socks5", proxies[0].Protocol)
placeholder

func (s *ProxyRepoSuite) TestListWithFilters_Status() {
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p1", Status: service.StatusActiveplaceholder)
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p2", Status: service.StatusDisabledplaceholder)

	proxies, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", service.StatusDisabled, "")
	s.Require().NoError(err)
	s.Require().Len(proxies, 1)
	s.Require().Equal(service.StatusDisabled, proxies[0].Status)
placeholder

func (s *ProxyRepoSuite) TestListWithFilters_Search() {
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "production-proxy"placeholder)
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "dev-proxy"placeholder)

	proxies, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 10placeholder, "", "", "prod")
	s.Require().NoError(err)
	s.Require().Len(proxies, 1)
	s.Require().Contains(proxies[0].Name, "production")
placeholder

// --- ListActive ---

func (s *ProxyRepoSuite) TestListActive() {
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "active1", Status: service.StatusActiveplaceholder)
	mustCreateProxy(s.T(), s.db, &proxyModel{Name: "inactive1", Status: service.StatusDisabledplaceholder)

	proxies, err := s.repo.ListActive(s.ctx)
	s.Require().NoError(err, "ListActive")
	s.Require().Len(proxies, 1)
	s.Require().Equal("active1", proxies[0].Name)
placeholder

// --- ExistsByHostPortAuth ---

func (s *ProxyRepoSuite) TestExistsByHostPortAuth() {
	mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:     "p1",
		Protocol: "http",
		Host:     "1.2.3.4",
		Port:     8080,
		Username: "user",
		Password: "pass",
placeholder)

	exists, err := s.repo.ExistsByHostPortAuth(s.ctx, "1.2.3.4", 8080, "user", "pass")
	s.Require().NoError(err, "ExistsByHostPortAuth")
	s.Require().True(exists)

	notExists, err := s.repo.ExistsByHostPortAuth(s.ctx, "1.2.3.4", 8080, "wrong", "creds")
	s.Require().NoError(err)
	s.Require().False(notExists)
placeholder

func (s *ProxyRepoSuite) TestExistsByHostPortAuth_NoAuth() {
	mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:     "p-noauth",
		Protocol: "http",
		Host:     "5.6.7.8",
		Port:     8081,
		Username: "",
		Password: "",
placeholder)

	exists, err := s.repo.ExistsByHostPortAuth(s.ctx, "5.6.7.8", 8081, "", "")
	s.Require().NoError(err)
	s.Require().True(exists)
placeholder

// --- CountAccountsByProxyID ---

func (s *ProxyRepoSuite) TestCountAccountsByProxyID() {
	proxy := mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p-count"placeholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a1", ProxyID: &proxy.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a2", ProxyID: &proxy.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a3"placeholder) // no proxy

	count, err := s.repo.CountAccountsByProxyID(s.ctx, proxy.ID)
	s.Require().NoError(err, "CountAccountsByProxyID")
	s.Require().Equal(int64(2), count)
placeholder

func (s *ProxyRepoSuite) TestCountAccountsByProxyID_Zero() {
	proxy := mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p-zero"placeholder)

	count, err := s.repo.CountAccountsByProxyID(s.ctx, proxy.ID)
	s.Require().NoError(err)
	s.Require().Zero(count)
placeholder

// --- GetAccountCountsForProxies ---

func (s *ProxyRepoSuite) TestGetAccountCountsForProxies() {
	p1 := mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p1"placeholder)
	p2 := mustCreateProxy(s.T(), s.db, &proxyModel{Name: "p2"placeholder)

	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a1", ProxyID: &p1.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a2", ProxyID: &p1.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a3", ProxyID: &p2.IDplaceholder)

	counts, err := s.repo.GetAccountCountsForProxies(s.ctx)
	s.Require().NoError(err, "GetAccountCountsForProxies")
	s.Require().Equal(int64(2), counts[p1.ID])
	s.Require().Equal(int64(1), counts[p2.ID])
placeholder

func (s *ProxyRepoSuite) TestGetAccountCountsForProxies_Empty() {
	counts, err := s.repo.GetAccountCountsForProxies(s.ctx)
	s.Require().NoError(err)
	s.Require().Empty(counts)
placeholder

// --- ListActiveWithAccountCount ---

func (s *ProxyRepoSuite) TestListActiveWithAccountCount() {
	base := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	p1 := mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:      "p1",
		Status:    service.StatusActive,
		CreatedAt: base.Add(-1 * time.Hour),
placeholder)
	p2 := mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:      "p2",
		Status:    service.StatusActive,
		CreatedAt: base,
placeholder)
	mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:   "p3-inactive",
		Status: service.StatusDisabled,
placeholder)

	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a1", ProxyID: &p1.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a2", ProxyID: &p1.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a3", ProxyID: &p2.IDplaceholder)

	withCounts, err := s.repo.ListActiveWithAccountCount(s.ctx)
	s.Require().NoError(err, "ListActiveWithAccountCount")
	s.Require().Len(withCounts, 2, "expected 2 active proxies")

	// Sorted by created_at DESC, so p2 first
	s.Require().Equal(p2.ID, withCounts[0].ID)
	s.Require().Equal(int64(1), withCounts[0].AccountCount)
	s.Require().Equal(p1.ID, withCounts[1].ID)
	s.Require().Equal(int64(2), withCounts[1].AccountCount)
placeholder

// --- Combined original test ---

func (s *ProxyRepoSuite) TestExistsByHostPortAuth_And_AccountCountAggregates() {
	p1 := mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:      "p1",
		Protocol:  "http",
		Host:      "1.2.3.4",
		Port:      8080,
		Username:  "u",
		Password:  "p",
		CreatedAt: time.Now().Add(-1 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
placeholder)
	p2 := mustCreateProxy(s.T(), s.db, &proxyModel{
		Name:      "p2",
		Protocol:  "http",
		Host:      "5.6.7.8",
		Port:      8081,
		Username:  "",
		Password:  "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
placeholder)

	exists, err := s.repo.ExistsByHostPortAuth(s.ctx, "1.2.3.4", 8080, "u", "p")
	s.Require().NoError(err, "ExistsByHostPortAuth")
	s.Require().True(exists, "expected proxy to exist")

	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a1", ProxyID: &p1.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a2", ProxyID: &p1.IDplaceholder)
	mustCreateAccount(s.T(), s.db, &accountModel{Name: "a3", ProxyID: &p2.IDplaceholder)

	count1, err := s.repo.CountAccountsByProxyID(s.ctx, p1.ID)
	s.Require().NoError(err, "CountAccountsByProxyID")
	s.Require().Equal(int64(2), count1, "expected 2 accounts for p1")

	counts, err := s.repo.GetAccountCountsForProxies(s.ctx)
	s.Require().NoError(err, "GetAccountCountsForProxies")
	s.Require().Equal(int64(2), counts[p1.ID])
	s.Require().Equal(int64(1), counts[p2.ID])

	withCounts, err := s.repo.ListActiveWithAccountCount(s.ctx)
	s.Require().NoError(err, "ListActiveWithAccountCount")
	s.Require().Len(withCounts, 2, "expected 2 proxies")
	for _, pc := range withCounts {
		switch pc.ID {
		case p1.ID:
			s.Require().Equal(int64(2), pc.AccountCount, "p1 count mismatch")
		case p2.ID:
			s.Require().Equal(int64(1), pc.AccountCount, "p2 count mismatch")
		default:
			s.Require().Fail("unexpected proxy id", pc.ID)
	placeholder
placeholder
placeholder
