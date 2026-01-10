// Package dto provides data transfer objects for HTTP handlers.
package dto

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func UserFromServiceShallow(u *service.User) *User {
	if u == nil {
		return nil
placeholder
	return &User{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		Notes:         u.Notes,
		Role:          u.Role,
		Balance:       u.Balance,
		Concurrency:   u.Concurrency,
		Status:        u.Status,
		AllowedGroups: u.AllowedGroups,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
placeholder
placeholder

func UserFromService(u *service.User) *User {
	if u == nil {
		return nil
placeholder
	out := UserFromServiceShallow(u)
	if len(u.APIKeys) > 0 {
		out.APIKeys = make([]APIKey, 0, len(u.APIKeys))
		for i := range u.APIKeys {
			k := u.APIKeys[i]
			out.APIKeys = append(out.APIKeys, *APIKeyFromService(&k))
	placeholder
placeholder
	if len(u.Subscriptions) > 0 {
		out.Subscriptions = make([]UserSubscription, 0, len(u.Subscriptions))
		for i := range u.Subscriptions {
			s := u.Subscriptions[i]
			out.Subscriptions = append(out.Subscriptions, *UserSubscriptionFromService(&s))
	placeholder
placeholder
	return out
placeholder

func APIKeyFromService(k *service.APIKey) *APIKey {
	if k == nil {
		return nil
placeholder
	return &APIKey{
		ID:          k.ID,
		UserID:      k.UserID,
		Key:         k.Key,
		Name:        k.Name,
		GroupID:     k.GroupID,
		Status:      k.Status,
		IPWhitelist: k.IPWhitelist,
		IPBlacklist: k.IPBlacklist,
		CreatedAt:   k.CreatedAt,
		UpdatedAt:   k.UpdatedAt,
		User:        UserFromServiceShallow(k.User),
		Group:       GroupFromServiceShallow(k.Group),
placeholder
placeholder

func GroupFromServiceShallow(g *service.Group) *Group {
	if g == nil {
		return nil
placeholder
	return &Group{
		ID:               g.ID,
		Name:             g.Name,
		Description:      g.Description,
		Platform:         g.Platform,
		RateMultiplier:   g.RateMultiplier,
		IsExclusive:      g.IsExclusive,
		Status:           g.Status,
		SubscriptionType: g.SubscriptionType,
		DailyLimitUSD:    g.DailyLimitUSD,
		WeeklyLimitUSD:   g.WeeklyLimitUSD,
		MonthlyLimitUSD:  g.MonthlyLimitUSD,
		ImagePrice1K:     g.ImagePrice1K,
		ImagePrice2K:     g.ImagePrice2K,
		ImagePrice4K:     g.ImagePrice4K,
		ClaudeCodeOnly:   g.ClaudeCodeOnly,
		FallbackGroupID:  g.FallbackGroupID,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		AccountCount:     g.AccountCount,
placeholder
placeholder

func GroupFromService(g *service.Group) *Group {
	if g == nil {
		return nil
placeholder
	out := GroupFromServiceShallow(g)
	if len(g.AccountGroups) > 0 {
		out.AccountGroups = make([]AccountGroup, 0, len(g.AccountGroups))
		for i := range g.AccountGroups {
			ag := g.AccountGroups[i]
			out.AccountGroups = append(out.AccountGroups, *AccountGroupFromService(&ag))
	placeholder
placeholder
	return out
placeholder

func AccountFromServiceShallow(a *service.Account) *Account {
	if a == nil {
		return nil
placeholder
placeholder
		ID:                      a.ID,
		Name:                    a.Name,
		Notes:                   a.Notes,
		Platform:                a.Platform,
		Type:                    a.Type,
		Credentials:             a.Credentials,
		Extra:                   a.Extra,
		ProxyID:                 a.ProxyID,
		Concurrency:             a.Concurrency,
		Priority:                a.Priority,
		Status:                  a.Status,
		ErrorMessage:            a.ErrorMessage,
		LastUsedAt:              a.LastUsedAt,
		ExpiresAt:               timeToUnixSeconds(a.ExpiresAt),
		AutoPauseOnExpired:      a.AutoPauseOnExpired,
		CreatedAt:               a.CreatedAt,
		UpdatedAt:               a.UpdatedAt,
		Schedulable:             a.Schedulable,
		RateLimitedAt:           a.RateLimitedAt,
		RateLimitResetAt:        a.RateLimitResetAt,
		OverloadUntil:           a.OverloadUntil,
		TempUnschedulableUntil:  a.TempUnschedulableUntil,
		TempUnschedulableReason: a.TempUnschedulableReason,
		SessionWindowStart:      a.SessionWindowStart,
		SessionWindowEnd:        a.SessionWindowEnd,
		SessionWindowStatus:     a.SessionWindowStatus,
		GroupIDs:                a.GroupIDs,
placeholder
placeholder

func AccountFromService(a *service.Account) *Account {
	if a == nil {
		return nil
placeholder
	out := AccountFromServiceShallow(a)
	out.Proxy = ProxyFromService(a.Proxy)
	if len(a.AccountGroups) > 0 {
		out.AccountGroups = make([]AccountGroup, 0, len(a.AccountGroups))
		for i := range a.AccountGroups {
			ag := a.AccountGroups[i]
			out.AccountGroups = append(out.AccountGroups, *AccountGroupFromService(&ag))
	placeholder
placeholder
	if len(a.Groups) > 0 {
		out.Groups = make([]*Group, 0, len(a.Groups))
		for _, g := range a.Groups {
			out.Groups = append(out.Groups, GroupFromServiceShallow(g))
	placeholder
placeholder
	return out
placeholder

func timeToUnixSeconds(value *time.Time) *int64 {
	if value == nil {
		return nil
placeholder
	ts := value.Unix()
	return &ts
placeholder

func AccountGroupFromService(ag *service.AccountGroup) *AccountGroup {
	if ag == nil {
		return nil
placeholder
	return &AccountGroup{
		AccountID: ag.AccountID,
		GroupID:   ag.GroupID,
		Priority:  ag.Priority,
		CreatedAt: ag.CreatedAt,
		Account:   AccountFromServiceShallow(ag.Account),
		Group:     GroupFromServiceShallow(ag.Group),
placeholder
placeholder

func ProxyFromService(p *service.Proxy) *Proxy {
	if p == nil {
		return nil
placeholder
	return &Proxy{
		ID:        p.ID,
		Name:      p.Name,
		Protocol:  p.Protocol,
		Host:      p.Host,
		Port:      p.Port,
		Username:  p.Username,
		Password:  p.Password,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
placeholder
placeholder

func ProxyWithAccountCountFromService(p *service.ProxyWithAccountCount) *ProxyWithAccountCount {
	if p == nil {
		return nil
placeholder
	return &ProxyWithAccountCount{
		Proxy:        *ProxyFromService(&p.Proxy),
		AccountCount: p.AccountCount,
placeholder
placeholder

func RedeemCodeFromService(rc *service.RedeemCode) *RedeemCode {
	if rc == nil {
		return nil
placeholder
	return &RedeemCode{
		ID:           rc.ID,
		Code:         rc.Code,
		Type:         rc.Type,
		Value:        rc.Value,
		Status:       rc.Status,
		UsedBy:       rc.UsedBy,
		UsedAt:       rc.UsedAt,
		Notes:        rc.Notes,
		CreatedAt:    rc.CreatedAt,
		GroupID:      rc.GroupID,
		ValidityDays: rc.ValidityDays,
		User:         UserFromServiceShallow(rc.User),
		Group:        GroupFromServiceShallow(rc.Group),
placeholder
placeholder

// AccountSummaryFromService returns a minimal AccountSummary for usage log display.
// Only includes ID and Name - no sensitive fields like Credentials, Proxy, etc.
func AccountSummaryFromService(a *service.Account) *AccountSummary {
	if a == nil {
		return nil
placeholder
	return &AccountSummary{
		ID:   a.ID,
		Name: a.Name,
placeholder
placeholder

// usageLogFromServiceBase is a helper that converts service UsageLog to DTO.
// The account parameter allows caller to control what Account info is included.
// The includeIPAddress parameter controls whether to include the IP address (admin-only).
func usageLogFromServiceBase(l *service.UsageLog, account *AccountSummary, includeIPAddress bool) *UsageLog {
	if l == nil {
		return nil
placeholder
	result := &UsageLog{
		ID:                    l.ID,
		UserID:                l.UserID,
		APIKeyID:              l.APIKeyID,
		AccountID:             l.AccountID,
		RequestID:             l.RequestID,
		Model:                 l.Model,
		GroupID:               l.GroupID,
		SubscriptionID:        l.SubscriptionID,
		InputTokens:           l.InputTokens,
		OutputTokens:          l.OutputTokens,
		CacheCreationTokens:   l.CacheCreationTokens,
		CacheReadTokens:       l.CacheReadTokens,
		CacheCreation5mTokens: l.CacheCreation5mTokens,
		CacheCreation1hTokens: l.CacheCreation1hTokens,
		InputCost:             l.InputCost,
		OutputCost:            l.OutputCost,
		CacheCreationCost:     l.CacheCreationCost,
		CacheReadCost:         l.CacheReadCost,
		TotalCost:             l.TotalCost,
		ActualCost:            l.ActualCost,
		RateMultiplier:        l.RateMultiplier,
		BillingType:           l.BillingType,
		Stream:                l.Stream,
		DurationMs:            l.DurationMs,
		FirstTokenMs:          l.FirstTokenMs,
		ImageCount:            l.ImageCount,
		ImageSize:             l.ImageSize,
		UserAgent:             l.UserAgent,
		CreatedAt:             l.CreatedAt,
		User:                  UserFromServiceShallow(l.User),
		APIKey:                APIKeyFromService(l.APIKey),
		Account:               account,
		Group:                 GroupFromServiceShallow(l.Group),
		Subscription:          UserSubscriptionFromService(l.Subscription),
placeholder
	// IP 地址仅对管理员可见
	if includeIPAddress {
		result.IPAddress = l.IPAddress
placeholder
	return result
placeholder

// UsageLogFromService converts a service UsageLog to DTO for regular users.
// It excludes Account details and IP address - users should not see these.
func UsageLogFromService(l *service.UsageLog) *UsageLog {
	return usageLogFromServiceBase(l, nil, false)
placeholder

// UsageLogFromServiceAdmin converts a service UsageLog to DTO for admin users.
// It includes minimal Account info (ID, Name only) and IP address.
func UsageLogFromServiceAdmin(l *service.UsageLog) *UsageLog {
	if l == nil {
		return nil
placeholder
	return usageLogFromServiceBase(l, AccountSummaryFromService(l.Account), true)
placeholder

func SettingFromService(s *service.Setting) *Setting {
	if s == nil {
		return nil
placeholder
	return &Setting{
		ID:        s.ID,
		Key:       s.Key,
		Value:     s.Value,
		UpdatedAt: s.UpdatedAt,
placeholder
placeholder

func UserSubscriptionFromService(sub *service.UserSubscription) *UserSubscription {
	if sub == nil {
		return nil
placeholder
	return &UserSubscription{
		ID:                 sub.ID,
		UserID:             sub.UserID,
		GroupID:            sub.GroupID,
		StartsAt:           sub.StartsAt,
		ExpiresAt:          sub.ExpiresAt,
		Status:             sub.Status,
		DailyWindowStart:   sub.DailyWindowStart,
		WeeklyWindowStart:  sub.WeeklyWindowStart,
		MonthlyWindowStart: sub.MonthlyWindowStart,
		DailyUsageUSD:      sub.DailyUsageUSD,
		WeeklyUsageUSD:     sub.WeeklyUsageUSD,
		MonthlyUsageUSD:    sub.MonthlyUsageUSD,
		AssignedBy:         sub.AssignedBy,
		AssignedAt:         sub.AssignedAt,
		Notes:              sub.Notes,
		CreatedAt:          sub.CreatedAt,
		UpdatedAt:          sub.UpdatedAt,
		User:               UserFromServiceShallow(sub.User),
		Group:              GroupFromServiceShallow(sub.Group),
		AssignedByUser:     UserFromServiceShallow(sub.AssignedByUser),
placeholder
placeholder

func BulkAssignResultFromService(r *service.BulkAssignResult) *BulkAssignResult {
	if r == nil {
		return nil
placeholder
	subs := make([]UserSubscription, 0, len(r.Subscriptions))
	for i := range r.Subscriptions {
		subs = append(subs, *UserSubscriptionFromService(&r.Subscriptions[i]))
placeholder
	return &BulkAssignResult{
		SuccessCount:  r.SuccessCount,
		FailedCount:   r.FailedCount,
		Subscriptions: subs,
		Errors:        r.Errors,
placeholder
placeholder
