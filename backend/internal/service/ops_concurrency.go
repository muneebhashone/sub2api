package service

import (
	"context"
	"log"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
)

const (
	opsAccountsPageSize          = 100
	opsConcurrencyBatchChunkSize = 200
)

func (s *OpsService) listAllAccountsForOps(ctx context.Context, platformFilter string) ([]Account, error) {
	if s == nil || s.accountRepo == nil {
		return []Account{placeholder, nil
placeholder

	out := make([]Account, 0, 128)
	page := 1
	for {
		accounts, pageInfo, err := s.accountRepo.ListWithFilters(ctx, pagination.PaginationParams{
			Page:     page,
			PageSize: opsAccountsPageSize,
	placeholder, platformFilter, "", "", "", 0)
		if err != nil {
			return nil, err
	placeholder
		if len(accounts) == 0 {
			break
	placeholder

		out = append(out, accounts...)
		if pageInfo != nil && int64(len(out)) >= pageInfo.Total {
			break
	placeholder
		if len(accounts) < opsAccountsPageSize {
			break
	placeholder

		page++
		if page > 10_000 {
			log.Printf("[Ops] listAllAccountsForOps: aborting after too many pages (platform=%q)", platformFilter)
			break
	placeholder
placeholder

	return out, nil
placeholder

func (s *OpsService) getAccountsLoadMapBestEffort(ctx context.Context, accounts []Account) map[int64]*AccountLoadInfo {
	if s == nil || s.concurrencyService == nil {
		return map[int64]*AccountLoadInfo{placeholder
placeholder
	if len(accounts) == 0 {
		return map[int64]*AccountLoadInfo{placeholder
placeholder

	// De-duplicate IDs (and keep the max concurrency to avoid under-reporting).
	unique := make(map[int64]int, len(accounts))
	for _, acc := range accounts {
		if acc.ID <= 0 {
			continue
	placeholder
		c := acc.Concurrency
		if c <= 0 {
			c = 1
	placeholder
		if prev, ok := unique[acc.ID]; !ok || c > prev {
			unique[acc.ID] = c
	placeholder
placeholder

	batch := make([]AccountWithConcurrency, 0, len(unique))
	for id, maxConc := range unique {
		batch = append(batch, AccountWithConcurrency{
			ID:             id,
			MaxConcurrency: maxConc,
	placeholder)
placeholder

	out := make(map[int64]*AccountLoadInfo, len(batch))
	for i := 0; i < len(batch); i += opsConcurrencyBatchChunkSize {
		end := i + opsConcurrencyBatchChunkSize
		if end > len(batch) {
			end = len(batch)
	placeholder
		part, err := s.concurrencyService.GetAccountsLoadBatch(ctx, batch[i:end])
		if err != nil {
			// Best-effort: return zeros rather than failing the ops UI.
			log.Printf("[Ops] GetAccountsLoadBatch failed: %v", err)
			continue
	placeholder
		for k, v := range part {
			out[k] = v
	placeholder
placeholder

	return out
placeholder

// GetConcurrencyStats returns real-time concurrency usage aggregated by platform/group/account.
//
// Optional filters:
// - platformFilter: only include accounts in that platform (best-effort reduces DB load)
// - groupIDFilter: only include accounts that belong to that group
func (s *OpsService) GetConcurrencyStats(
	ctx context.Context,
	platformFilter string,
	groupIDFilter *int64,
) (map[string]*PlatformConcurrencyInfo, map[int64]*GroupConcurrencyInfo, map[int64]*AccountConcurrencyInfo, *time.Time, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, nil, nil, nil, err
placeholder

	accounts, err := s.listAllAccountsForOps(ctx, platformFilter)
	if err != nil {
		return nil, nil, nil, nil, err
placeholder

	collectedAt := time.Now()
	loadMap := s.getAccountsLoadMapBestEffort(ctx, accounts)

	platform := make(map[string]*PlatformConcurrencyInfo)
	group := make(map[int64]*GroupConcurrencyInfo)
	account := make(map[int64]*AccountConcurrencyInfo)

	for _, acc := range accounts {
		if acc.ID <= 0 {
			continue
	placeholder

		var matchedGroup *Group
		if groupIDFilter != nil && *groupIDFilter > 0 {
			for _, grp := range acc.Groups {
				if grp == nil || grp.ID <= 0 {
					continue
			placeholder
				if grp.ID == *groupIDFilter {
					matchedGroup = grp
					break
			placeholder
		placeholder
			// Group filter provided: skip accounts not in that group.
			if matchedGroup == nil {
				continue
		placeholder
	placeholder

		load := loadMap[acc.ID]
		currentInUse := int64(0)
		waiting := int64(0)
		if load != nil {
			currentInUse = int64(load.CurrentConcurrency)
			waiting = int64(load.WaitingCount)
	placeholder

		// Account-level view picks one display group (the first group).
		displayGroupID := int64(0)
		displayGroupName := ""
		if matchedGroup != nil {
			displayGroupID = matchedGroup.ID
			displayGroupName = matchedGroup.Name
	placeholder else if len(acc.Groups) > 0 && acc.Groups[0] != nil {
			displayGroupID = acc.Groups[0].ID
			displayGroupName = acc.Groups[0].Name
	placeholder

		if _, ok := account[acc.ID]; !ok {
			info := &AccountConcurrencyInfo{
				AccountID:      acc.ID,
				AccountName:    acc.Name,
				Platform:       acc.Platform,
				GroupID:        displayGroupID,
				GroupName:      displayGroupName,
				CurrentInUse:   currentInUse,
				MaxCapacity:    int64(acc.Concurrency),
				WaitingInQueue: waiting,
		placeholder
			if info.MaxCapacity > 0 {
				info.LoadPercentage = float64(info.CurrentInUse) / float64(info.MaxCapacity) * 100
		placeholder
			account[acc.ID] = info
	placeholder

		// Platform aggregation.
		if acc.Platform != "" {
			if _, ok := platform[acc.Platform]; !ok {
				platform[acc.Platform] = &PlatformConcurrencyInfo{
					Platform: acc.Platform,
			placeholder
		placeholder
			p := platform[acc.Platform]
			p.MaxCapacity += int64(acc.Concurrency)
			p.CurrentInUse += currentInUse
			p.WaitingInQueue += waiting
	placeholder

		// Group aggregation (one account may contribute to multiple groups).
		if matchedGroup != nil {
			grp := matchedGroup
			if _, ok := group[grp.ID]; !ok {
				group[grp.ID] = &GroupConcurrencyInfo{
					GroupID:   grp.ID,
					GroupName: grp.Name,
					Platform:  grp.Platform,
			placeholder
		placeholder
			g := group[grp.ID]
			if g.GroupName == "" && grp.Name != "" {
				g.GroupName = grp.Name
		placeholder
			if g.Platform != "" && grp.Platform != "" && g.Platform != grp.Platform {
				// Groups are expected to be platform-scoped. If mismatch is observed, avoid misleading labels.
				g.Platform = ""
		placeholder
			g.MaxCapacity += int64(acc.Concurrency)
			g.CurrentInUse += currentInUse
			g.WaitingInQueue += waiting
	placeholder else {
			for _, grp := range acc.Groups {
				if grp == nil || grp.ID <= 0 {
					continue
			placeholder
				if _, ok := group[grp.ID]; !ok {
					group[grp.ID] = &GroupConcurrencyInfo{
						GroupID:   grp.ID,
						GroupName: grp.Name,
						Platform:  grp.Platform,
				placeholder
			placeholder
				g := group[grp.ID]
				if g.GroupName == "" && grp.Name != "" {
					g.GroupName = grp.Name
			placeholder
				if g.Platform != "" && grp.Platform != "" && g.Platform != grp.Platform {
					// Groups are expected to be platform-scoped. If mismatch is observed, avoid misleading labels.
					g.Platform = ""
			placeholder
				g.MaxCapacity += int64(acc.Concurrency)
				g.CurrentInUse += currentInUse
				g.WaitingInQueue += waiting
		placeholder
	placeholder
placeholder

	for _, info := range platform {
		if info.MaxCapacity > 0 {
			info.LoadPercentage = float64(info.CurrentInUse) / float64(info.MaxCapacity) * 100
	placeholder
placeholder
	for _, info := range group {
		if info.MaxCapacity > 0 {
			info.LoadPercentage = float64(info.CurrentInUse) / float64(info.MaxCapacity) * 100
	placeholder
placeholder

	return platform, group, account, &collectedAt, nil
placeholder

// listAllActiveUsersForOps returns all active users with their concurrency settings.
func (s *OpsService) listAllActiveUsersForOps(ctx context.Context) ([]User, error) {
	if s == nil || s.userRepo == nil {
		return []User{placeholder, nil
placeholder

	out := make([]User, 0, 128)
	page := 1
	for {
		users, pageInfo, err := s.userRepo.ListWithFilters(ctx, pagination.PaginationParams{
			Page:     page,
			PageSize: opsAccountsPageSize,
	placeholder, UserListFilters{
			Status: StatusActive,
	placeholder)
		if err != nil {
			return nil, err
	placeholder
		if len(users) == 0 {
			break
	placeholder

		out = append(out, users...)
		if pageInfo != nil && int64(len(out)) >= pageInfo.Total {
			break
	placeholder
		if len(users) < opsAccountsPageSize {
			break
	placeholder

		page++
		if page > 10_000 {
			log.Printf("[Ops] listAllActiveUsersForOps: aborting after too many pages")
			break
	placeholder
placeholder

	return out, nil
placeholder

// getUsersLoadMapBestEffort returns user load info for the given users.
func (s *OpsService) getUsersLoadMapBestEffort(ctx context.Context, users []User) map[int64]*UserLoadInfo {
	if s == nil || s.concurrencyService == nil {
		return map[int64]*UserLoadInfo{placeholder
placeholder
	if len(users) == 0 {
		return map[int64]*UserLoadInfo{placeholder
placeholder

	// De-duplicate IDs (and keep the max concurrency to avoid under-reporting).
	unique := make(map[int64]int, len(users))
	for _, u := range users {
		if u.ID <= 0 {
			continue
	placeholder
		if prev, ok := unique[u.ID]; !ok || u.Concurrency > prev {
			unique[u.ID] = u.Concurrency
	placeholder
placeholder

	batch := make([]UserWithConcurrency, 0, len(unique))
	for id, maxConc := range unique {
		batch = append(batch, UserWithConcurrency{
			ID:             id,
			MaxConcurrency: maxConc,
	placeholder)
placeholder

	out := make(map[int64]*UserLoadInfo, len(batch))
	for i := 0; i < len(batch); i += opsConcurrencyBatchChunkSize {
		end := i + opsConcurrencyBatchChunkSize
		if end > len(batch) {
			end = len(batch)
	placeholder
		part, err := s.concurrencyService.GetUsersLoadBatch(ctx, batch[i:end])
		if err != nil {
			// Best-effort: return zeros rather than failing the ops UI.
			log.Printf("[Ops] GetUsersLoadBatch failed: %v", err)
			continue
	placeholder
		for k, v := range part {
			out[k] = v
	placeholder
placeholder

	return out
placeholder

// GetUserConcurrencyStats returns real-time concurrency usage for all active users.
func (s *OpsService) GetUserConcurrencyStats(ctx context.Context) (map[int64]*UserConcurrencyInfo, *time.Time, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, nil, err
placeholder

	users, err := s.listAllActiveUsersForOps(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	collectedAt := time.Now()
	loadMap := s.getUsersLoadMapBestEffort(ctx, users)

	result := make(map[int64]*UserConcurrencyInfo)

	for _, u := range users {
		if u.ID <= 0 {
			continue
	placeholder

		load := loadMap[u.ID]
		currentInUse := int64(0)
		waiting := int64(0)
		if load != nil {
			currentInUse = int64(load.CurrentConcurrency)
			waiting = int64(load.WaitingCount)
	placeholder

		// Skip users with no concurrency activity
		if currentInUse == 0 && waiting == 0 {
			continue
	placeholder

		info := &UserConcurrencyInfo{
			UserID:         u.ID,
			UserEmail:      u.Email,
			Username:       u.Username,
			CurrentInUse:   currentInUse,
			MaxCapacity:    int64(u.Concurrency),
			WaitingInQueue: waiting,
	placeholder
		if info.MaxCapacity > 0 {
			info.LoadPercentage = float64(info.CurrentInUse) / float64(info.MaxCapacity) * 100
	placeholder
		result[u.ID] = info
placeholder

	return result, &collectedAt, nil
placeholder
