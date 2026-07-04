package service

import (
	"context"
	"time"
)

// GroupCapacitySummary holds aggregated capacity for a single group.
type GroupCapacitySummary struct {
	GroupID         int64 `json:"group_id"`
	ConcurrencyUsed int   `json:"concurrency_used"`
	ConcurrencyMax  int   `json:"concurrency_max"`
	SessionsUsed    int   `json:"sessions_used"`
	SessionsMax     int   `json:"sessions_max"`
	RPMUsed         int   `json:"rpm_used"`
	RPMMax          int   `json:"rpm_max"`
placeholder

// GroupAccountCapacityRow is the lightweight account projection needed for
// capacity summary aggregation.
type GroupAccountCapacityRow struct {
	GroupID             int64
	AccountID           int64
	Concurrency         int
	Extra               map[string]any
	SessionWindowStart  *time.Time
	SessionWindowEnd    *time.Time
	SessionWindowStatus string
placeholder

type groupCapacityActiveGroupIDLister interface {
	ListActiveIDs(ctx context.Context) ([]int64, error)
placeholder

type groupCapacityAccountLister interface {
	ListSchedulableCapacityByGroupIDs(ctx context.Context, groupIDs []int64) ([]GroupAccountCapacityRow, error)
placeholder

// GroupCapacityService aggregates per-group capacity from runtime data.
type GroupCapacityService struct {
	accountRepo        AccountRepository
	groupRepo          GroupRepository
	concurrencyService *ConcurrencyService
	sessionLimitCache  SessionLimitCache
	rpmCache           RPMCache
placeholder

// NewGroupCapacityService creates a new GroupCapacityService.
func NewGroupCapacityService(
	accountRepo AccountRepository,
	groupRepo GroupRepository,
	concurrencyService *ConcurrencyService,
	sessionLimitCache SessionLimitCache,
	rpmCache RPMCache,
) *GroupCapacityService {
	return &GroupCapacityService{
		accountRepo:        accountRepo,
		groupRepo:          groupRepo,
		concurrencyService: concurrencyService,
		sessionLimitCache:  sessionLimitCache,
		rpmCache:           rpmCache,
placeholder
placeholder

// GetAllGroupCapacity returns capacity summary for all active groups.
func (s *GroupCapacityService) GetAllGroupCapacity(ctx context.Context) ([]GroupCapacitySummary, error) {
	groupIDs, err := s.listActiveGroupIDs(ctx)
	if err != nil {
		return nil, err
placeholder

	if lister, ok := s.accountRepo.(groupCapacityAccountLister); ok {
		return s.getGroupCapacitiesBatch(ctx, groupIDs, lister)
placeholder

	return s.getGroupCapacitiesSequential(ctx, groupIDs), nil
placeholder

func (s *GroupCapacityService) listActiveGroupIDs(ctx context.Context) ([]int64, error) {
	if lister, ok := s.groupRepo.(groupCapacityActiveGroupIDLister); ok {
		return lister.ListActiveIDs(ctx)
placeholder

	groups, err := s.groupRepo.ListActive(ctx)
	if err != nil {
		return nil, err
placeholder
	groupIDs := make([]int64, 0, len(groups))
	for i := range groups {
		groupIDs = append(groupIDs, groups[i].ID)
placeholder
	return groupIDs, nil
placeholder

func (s *GroupCapacityService) getGroupCapacitiesSequential(ctx context.Context, groupIDs []int64) []GroupCapacitySummary {
	results := make([]GroupCapacitySummary, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		cap, err := s.getGroupCapacity(ctx, groupID)
		if err != nil {
			// Skip groups with errors, return partial results
			continue
	placeholder
		cap.GroupID = groupID
		results = append(results, cap)
placeholder
	return results
placeholder

type groupCapacityAccountRef struct {
	groupID   int64
	accountID int64
placeholder

func (s *GroupCapacityService) getGroupCapacitiesBatch(ctx context.Context, groupIDs []int64, lister groupCapacityAccountLister) ([]GroupCapacitySummary, error) {
	results := make([]GroupCapacitySummary, len(groupIDs))
	groupIndex := make(map[int64]int, len(groupIDs))
	for i, groupID := range groupIDs {
		results[i].GroupID = groupID
		groupIndex[groupID] = i
placeholder
	if len(groupIDs) == 0 {
		return results, nil
placeholder

	rows, err := lister.ListSchedulableCapacityByGroupIDs(ctx, groupIDs)
	if err != nil {
		return nil, err
placeholder
	if len(rows) == 0 {
		return results, nil
placeholder

	refs := make([]groupCapacityAccountRef, 0, len(rows))
	seenGroupAccount := make(map[groupCapacityAccountRef]struct{placeholder, len(rows))
	accountIDSet := make(map[int64]struct{placeholder, len(rows))
	accountIDs := make([]int64, 0, len(rows))
	sessionTimeouts := make(map[int64]time.Duration)

	for _, row := range rows {
		idx, ok := groupIndex[row.GroupID]
		if !ok || row.AccountID <= 0 {
			continue
	placeholder

		ref := groupCapacityAccountRef{groupID: row.GroupID, accountID: row.AccountIDplaceholder
		if _, ok := seenGroupAccount[ref]; ok {
			continue
	placeholder
		seenGroupAccount[ref] = struct{placeholder{placeholder
		refs = append(refs, ref)

		if _, ok := accountIDSet[row.AccountID]; !ok {
			accountIDSet[row.AccountID] = struct{placeholder{placeholder
			accountIDs = append(accountIDs, row.AccountID)
	placeholder

		acc := Account{
			ID:                  row.AccountID,
			Concurrency:         row.Concurrency,
			Extra:               row.Extra,
			SessionWindowStart:  row.SessionWindowStart,
			SessionWindowEnd:    row.SessionWindowEnd,
			SessionWindowStatus: row.SessionWindowStatus,
	placeholder

		results[idx].ConcurrencyMax += acc.Concurrency

		if maxSessions := acc.GetMaxSessions(); maxSessions > 0 {
			results[idx].SessionsMax += maxSessions
			timeout := time.Duration(acc.GetSessionIdleTimeoutMinutes()) * time.Minute
			if timeout <= 0 {
				timeout = 5 * time.Minute
		placeholder
			sessionTimeouts[acc.ID] = timeout
	placeholder

		if rpm := acc.GetBaseRPM(); rpm > 0 {
			results[idx].RPMMax += rpm
	placeholder
placeholder

	if len(accountIDs) == 0 {
		return results, nil
placeholder

	concurrencyMap := map[int64]int{placeholder
	if s.concurrencyService != nil {
		concurrencyMap, _ = s.concurrencyService.GetAccountConcurrencyBatch(ctx, accountIDs)
placeholder

	sessionAccountIDs := accountIDsForGroupsWithLimit(refs, groupIndex, results, func(summary GroupCapacitySummary) bool {
		return summary.SessionsMax > 0
placeholder)
	var sessionsMap map[int64]int
	if len(sessionAccountIDs) > 0 && s.sessionLimitCache != nil {
		sessionsMap, _ = s.sessionLimitCache.GetActiveSessionCountBatch(ctx, sessionAccountIDs, sessionTimeouts)
placeholder

	rpmAccountIDs := accountIDsForGroupsWithLimit(refs, groupIndex, results, func(summary GroupCapacitySummary) bool {
		return summary.RPMMax > 0
placeholder)
	var rpmMap map[int64]int
	if len(rpmAccountIDs) > 0 && s.rpmCache != nil {
		rpmMap, _ = s.rpmCache.GetRPMBatch(ctx, rpmAccountIDs)
placeholder

	for _, ref := range refs {
		idx := groupIndex[ref.groupID]
		results[idx].ConcurrencyUsed += concurrencyMap[ref.accountID]
		if sessionsMap != nil && results[idx].SessionsMax > 0 {
			results[idx].SessionsUsed += sessionsMap[ref.accountID]
	placeholder
		if rpmMap != nil && results[idx].RPMMax > 0 {
			results[idx].RPMUsed += rpmMap[ref.accountID]
	placeholder
placeholder
	return results, nil
placeholder

func accountIDsForGroupsWithLimit(refs []groupCapacityAccountRef, groupIndex map[int64]int, summaries []GroupCapacitySummary, include func(GroupCapacitySummary) bool) []int64 {
	seen := make(map[int64]struct{placeholder)
	accountIDs := make([]int64, 0)
	for _, ref := range refs {
		idx, ok := groupIndex[ref.groupID]
		if !ok || !include(summaries[idx]) {
			continue
	placeholder
		if _, ok := seen[ref.accountID]; ok {
			continue
	placeholder
		seen[ref.accountID] = struct{placeholder{placeholder
		accountIDs = append(accountIDs, ref.accountID)
placeholder
	return accountIDs
placeholder

func (s *GroupCapacityService) getGroupCapacity(ctx context.Context, groupID int64) (GroupCapacitySummary, error) {
	accounts, err := s.accountRepo.ListSchedulableByGroupID(ctx, groupID)
	if err != nil {
		return GroupCapacitySummary{placeholder, err
placeholder
	if len(accounts) == 0 {
		return GroupCapacitySummary{placeholder, nil
placeholder

	// Collect account IDs and config values
	accountIDs := make([]int64, 0, len(accounts))
	sessionTimeouts := make(map[int64]time.Duration)
	var concurrencyMax, sessionsMax, rpmMax int

	for i := range accounts {
		acc := &accounts[i]
		accountIDs = append(accountIDs, acc.ID)
		concurrencyMax += acc.Concurrency

		if ms := acc.GetMaxSessions(); ms > 0 {
			sessionsMax += ms
			timeout := time.Duration(acc.GetSessionIdleTimeoutMinutes()) * time.Minute
			if timeout <= 0 {
				timeout = 5 * time.Minute
		placeholder
			sessionTimeouts[acc.ID] = timeout
	placeholder

		if rpm := acc.GetBaseRPM(); rpm > 0 {
			rpmMax += rpm
	placeholder
placeholder

	// Batch query runtime data from Redis
	concurrencyMap, _ := s.concurrencyService.GetAccountConcurrencyBatch(ctx, accountIDs)

	var sessionsMap map[int64]int
	if sessionsMax > 0 && s.sessionLimitCache != nil {
		sessionsMap, _ = s.sessionLimitCache.GetActiveSessionCountBatch(ctx, accountIDs, sessionTimeouts)
placeholder

	var rpmMap map[int64]int
	if rpmMax > 0 && s.rpmCache != nil {
		rpmMap, _ = s.rpmCache.GetRPMBatch(ctx, accountIDs)
placeholder

	// Aggregate
	var concurrencyUsed, sessionsUsed, rpmUsed int
	for _, id := range accountIDs {
		concurrencyUsed += concurrencyMap[id]
		if sessionsMap != nil {
			sessionsUsed += sessionsMap[id]
	placeholder
		if rpmMap != nil {
			rpmUsed += rpmMap[id]
	placeholder
placeholder

	return GroupCapacitySummary{
		ConcurrencyUsed: concurrencyUsed,
		ConcurrencyMax:  concurrencyMax,
		SessionsUsed:    sessionsUsed,
		SessionsMax:     sessionsMax,
		RPMUsed:         rpmUsed,
		RPMMax:          rpmMax,
placeholder, nil
placeholder
