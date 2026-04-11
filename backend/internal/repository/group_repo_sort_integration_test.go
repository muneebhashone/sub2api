//go:build integration

package repository

import (
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func (s *GroupRepoSuite) TestList_DefaultSortBySortOrderAsc() {
	g1 := &service.Group{Name: "g1", Platform: service.PlatformAnthropic, RateMultiplier: 1, Status: service.StatusActive, SubscriptionType: service.SubscriptionTypeStandard, SortOrder: 20placeholder
	g2 := &service.Group{Name: "g2", Platform: service.PlatformAnthropic, RateMultiplier: 1, Status: service.StatusActive, SubscriptionType: service.SubscriptionTypeStandard, SortOrder: 10placeholder
	s.Require().NoError(s.repo.Create(s.ctx, g1))
	s.Require().NoError(s.repo.Create(s.ctx, g2))

	groups, _, err := s.repo.List(s.ctx, pagination.PaginationParams{Page: 1, PageSize: 100placeholder)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(groups), 2)
	indexByID := make(map[int64]int, len(groups))
	for i, g := range groups {
		indexByID[g.ID] = i
placeholder
	s.Require().Contains(indexByID, g1.ID)
	s.Require().Contains(indexByID, g2.ID)
	// g2 has SortOrder=10, g1 has SortOrder=20; ascending means g2 comes first
	s.Require().Less(indexByID[g2.ID], indexByID[g1.ID])
placeholder

func (s *GroupRepoSuite) TestList_SortBySortOrderDesc() {
	g1 := &service.Group{Name: "g1", Platform: service.PlatformAnthropic, RateMultiplier: 1, Status: service.StatusActive, SubscriptionType: service.SubscriptionTypeStandard, SortOrder: 40placeholder
	g2 := &service.Group{Name: "g2", Platform: service.PlatformAnthropic, RateMultiplier: 1, Status: service.StatusActive, SubscriptionType: service.SubscriptionTypeStandard, SortOrder: 50placeholder
	s.Require().NoError(s.repo.Create(s.ctx, g1))
	s.Require().NoError(s.repo.Create(s.ctx, g2))

	groups, _, err := s.repo.List(s.ctx, pagination.PaginationParams{
		Page:      1,
		PageSize:  10,
		SortBy:    "sort_order",
		SortOrder: "desc",
placeholder)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(groups), 2)
	indexByID := make(map[int64]int, len(groups))
	for i, group := range groups {
		indexByID[group.ID] = i
placeholder
	s.Require().Contains(indexByID, g1.ID)
	s.Require().Contains(indexByID, g2.ID)
	s.Require().Less(indexByID[g2.ID], indexByID[g1.ID])
placeholder
