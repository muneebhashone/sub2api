//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func mustCreateUser(t *testing.T, client *dbent.Client, u *service.User) *service.User {
placeholder
	ctx := context.Background()

	if u.Email == "" {
		u.Email = "user-" + time.Now().Format(time.RFC3339Nano) + "@example.com"
placeholder
	if u.PasswordHash == "" {
		u.PasswordHash = "test-password-hash"
placeholder
	if u.Role == "" {
		u.Role = service.RoleUser
placeholder
	if u.Status == "" {
		u.Status = service.StatusActive
placeholder
	if u.Concurrency == 0 {
		u.Concurrency = 5
placeholder

	create := client.User.Create().
		SetEmail(u.Email).
		SetPasswordHash(u.PasswordHash).
		SetRole(u.Role).
		SetStatus(u.Status).
		SetBalance(u.Balance).
		SetConcurrency(u.Concurrency).
		SetUsername(u.Username).
		SetWechat(u.Wechat).
		SetNotes(u.Notes)
	if !u.CreatedAt.IsZero() {
		create.SetCreatedAt(u.CreatedAt)
placeholder
	if !u.UpdatedAt.IsZero() {
		create.SetUpdatedAt(u.UpdatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create user")

	u.ID = created.ID
	u.CreatedAt = created.CreatedAt
	u.UpdatedAt = created.UpdatedAt

	if len(u.AllowedGroups) > 0 {
		for _, groupID := range u.AllowedGroups {
			_, err := client.UserAllowedGroup.Create().
				SetUserID(u.ID).
				SetGroupID(groupID).
				Save(ctx)
			require.NoError(t, err, "create user_allowed_groups row")
	placeholder
placeholder

	return u
placeholder

func mustCreateGroup(t *testing.T, client *dbent.Client, g *service.Group) *service.Group {
placeholder
	ctx := context.Background()

	if g.Platform == "" {
		g.Platform = service.PlatformAnthropic
placeholder
	if g.Status == "" {
		g.Status = service.StatusActive
placeholder
	if g.SubscriptionType == "" {
		g.SubscriptionType = service.SubscriptionTypeStandard
placeholder

	create := client.Group.Create().
		SetName(g.Name).
		SetPlatform(g.Platform).
		SetStatus(g.Status).
		SetSubscriptionType(g.SubscriptionType).
		SetRateMultiplier(g.RateMultiplier).
		SetIsExclusive(g.IsExclusive)
	if g.Description != "" {
		create.SetDescription(g.Description)
placeholder
	if g.DailyLimitUSD != nil {
		create.SetDailyLimitUsd(*g.DailyLimitUSD)
placeholder
	if g.WeeklyLimitUSD != nil {
		create.SetWeeklyLimitUsd(*g.WeeklyLimitUSD)
placeholder
	if g.MonthlyLimitUSD != nil {
		create.SetMonthlyLimitUsd(*g.MonthlyLimitUSD)
placeholder
	if !g.CreatedAt.IsZero() {
		create.SetCreatedAt(g.CreatedAt)
placeholder
	if !g.UpdatedAt.IsZero() {
		create.SetUpdatedAt(g.UpdatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create group")

	g.ID = created.ID
	g.CreatedAt = created.CreatedAt
	g.UpdatedAt = created.UpdatedAt
	return g
placeholder

func mustCreateProxy(t *testing.T, client *dbent.Client, p *service.Proxy) *service.Proxy {
placeholder
	ctx := context.Background()

	if p.Protocol == "" {
		p.Protocol = "http"
placeholder
	if p.Host == "" {
		p.Host = "127.0.0.1"
placeholder
	if p.Port == 0 {
		p.Port = 8080
placeholder
	if p.Status == "" {
		p.Status = service.StatusActive
placeholder

	create := client.Proxy.Create().
		SetName(p.Name).
		SetProtocol(p.Protocol).
		SetHost(p.Host).
		SetPort(p.Port).
		SetStatus(p.Status)
	if p.Username != "" {
		create.SetUsername(p.Username)
placeholder
	if p.Password != "" {
		create.SetPassword(p.Password)
placeholder
	if !p.CreatedAt.IsZero() {
		create.SetCreatedAt(p.CreatedAt)
placeholder
	if !p.UpdatedAt.IsZero() {
		create.SetUpdatedAt(p.UpdatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create proxy")

	p.ID = created.ID
	p.CreatedAt = created.CreatedAt
	p.UpdatedAt = created.UpdatedAt
	return p
placeholder

func mustCreateAccount(t *testing.T, client *dbent.Client, a *service.Account) *service.Account {
placeholder
	ctx := context.Background()

	if a.Platform == "" {
		a.Platform = service.PlatformAnthropic
placeholder
	if a.Type == "" {
		a.Type = service.AccountTypeOAuth
placeholder
	if a.Status == "" {
		a.Status = service.StatusActive
placeholder
	if a.Concurrency == 0 {
		a.Concurrency = 3
placeholder
	if a.Priority == 0 {
		a.Priority = 50
placeholder
	if !a.Schedulable {
		a.Schedulable = true
placeholder
	if a.Credentials == nil {
		a.Credentials = map[string]any{placeholder
placeholder
	if a.Extra == nil {
		a.Extra = map[string]any{placeholder
placeholder

	create := client.Account.Create().
		SetName(a.Name).
		SetPlatform(a.Platform).
		SetType(a.Type).
		SetCredentials(a.Credentials).
		SetExtra(a.Extra).
		SetConcurrency(a.Concurrency).
		SetPriority(a.Priority).
		SetStatus(a.Status).
		SetSchedulable(a.Schedulable).
		SetErrorMessage(a.ErrorMessage)

	if a.ProxyID != nil {
		create.SetProxyID(*a.ProxyID)
placeholder
	if a.LastUsedAt != nil {
		create.SetLastUsedAt(*a.LastUsedAt)
placeholder
	if a.RateLimitedAt != nil {
		create.SetRateLimitedAt(*a.RateLimitedAt)
placeholder
	if a.RateLimitResetAt != nil {
		create.SetRateLimitResetAt(*a.RateLimitResetAt)
placeholder
	if a.OverloadUntil != nil {
		create.SetOverloadUntil(*a.OverloadUntil)
placeholder
	if a.SessionWindowStart != nil {
		create.SetSessionWindowStart(*a.SessionWindowStart)
placeholder
	if a.SessionWindowEnd != nil {
		create.SetSessionWindowEnd(*a.SessionWindowEnd)
placeholder
	if a.SessionWindowStatus != "" {
		create.SetSessionWindowStatus(a.SessionWindowStatus)
placeholder
	if !a.CreatedAt.IsZero() {
		create.SetCreatedAt(a.CreatedAt)
placeholder
	if !a.UpdatedAt.IsZero() {
		create.SetUpdatedAt(a.UpdatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create account")

	a.ID = created.ID
	a.CreatedAt = created.CreatedAt
	a.UpdatedAt = created.UpdatedAt
	return a
placeholder

func mustCreateApiKey(t *testing.T, client *dbent.Client, k *service.ApiKey) *service.ApiKey {
placeholder
	ctx := context.Background()

	if k.Status == "" {
		k.Status = service.StatusActive
placeholder
	if k.Key == "" {
		k.Key = "sk-" + time.Now().Format("150405.000000")
placeholder
	if k.Name == "" {
		k.Name = "default"
placeholder

	create := client.ApiKey.Create().
		SetUserID(k.UserID).
		SetKey(k.Key).
		SetName(k.Name).
		SetStatus(k.Status)
	if k.GroupID != nil {
		create.SetGroupID(*k.GroupID)
placeholder
	if !k.CreatedAt.IsZero() {
		create.SetCreatedAt(k.CreatedAt)
placeholder
	if !k.UpdatedAt.IsZero() {
		create.SetUpdatedAt(k.UpdatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create api key")

	k.ID = created.ID
	k.CreatedAt = created.CreatedAt
	k.UpdatedAt = created.UpdatedAt
	return k
placeholder

func mustCreateRedeemCode(t *testing.T, client *dbent.Client, c *service.RedeemCode) *service.RedeemCode {
placeholder
	ctx := context.Background()

	if c.Status == "" {
		c.Status = service.StatusUnused
placeholder
	if c.Type == "" {
		c.Type = service.RedeemTypeBalance
placeholder
	if c.Code == "" {
		c.Code = "rc-" + time.Now().Format("150405.000000")
placeholder

	create := client.RedeemCode.Create().
		SetCode(c.Code).
		SetType(c.Type).
		SetValue(c.Value).
		SetStatus(c.Status).
		SetNotes(c.Notes).
		SetValidityDays(c.ValidityDays)
	if c.UsedBy != nil {
		create.SetUsedBy(*c.UsedBy)
placeholder
	if c.UsedAt != nil {
		create.SetUsedAt(*c.UsedAt)
placeholder
	if c.GroupID != nil {
		create.SetGroupID(*c.GroupID)
placeholder
	if !c.CreatedAt.IsZero() {
		create.SetCreatedAt(c.CreatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create redeem code")

	c.ID = created.ID
	c.CreatedAt = created.CreatedAt
	return c
placeholder

func mustCreateSubscription(t *testing.T, client *dbent.Client, s *service.UserSubscription) *service.UserSubscription {
placeholder
	ctx := context.Background()

	if s.Status == "" {
		s.Status = service.SubscriptionStatusActive
placeholder
	now := time.Now()
	if s.StartsAt.IsZero() {
		s.StartsAt = now.Add(-1 * time.Hour)
placeholder
	if s.ExpiresAt.IsZero() {
		s.ExpiresAt = now.Add(24 * time.Hour)
placeholder
	if s.AssignedAt.IsZero() {
		s.AssignedAt = now
placeholder
	if s.CreatedAt.IsZero() {
		s.CreatedAt = now
placeholder
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = now
placeholder

	create := client.UserSubscription.Create().
		SetUserID(s.UserID).
		SetGroupID(s.GroupID).
		SetStartsAt(s.StartsAt).
		SetExpiresAt(s.ExpiresAt).
		SetStatus(s.Status).
		SetAssignedAt(s.AssignedAt).
		SetNotes(s.Notes).
		SetDailyUsageUsd(s.DailyUsageUSD).
		SetWeeklyUsageUsd(s.WeeklyUsageUSD).
		SetMonthlyUsageUsd(s.MonthlyUsageUSD)

	if s.AssignedBy != nil {
		create.SetAssignedBy(*s.AssignedBy)
placeholder
	if !s.CreatedAt.IsZero() {
		create.SetCreatedAt(s.CreatedAt)
placeholder
	if !s.UpdatedAt.IsZero() {
		create.SetUpdatedAt(s.UpdatedAt)
placeholder

	created, err := create.Save(ctx)
	require.NoError(t, err, "create user subscription")

	s.ID = created.ID
	s.CreatedAt = created.CreatedAt
	s.UpdatedAt = created.UpdatedAt
	return s
placeholder

func mustBindAccountToGroup(t *testing.T, client *dbent.Client, accountID, groupID int64, priority int) {
placeholder
	ctx := context.Background()

	_, err := client.AccountGroup.Create().
		SetAccountID(accountID).
		SetGroupID(groupID).
		SetPriority(priority).
		Save(ctx)
	require.NoError(t, err, "create account_group")
placeholder
