//go:build integration

package repository

import (
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func mustCreateUser(t *testing.T, db *gorm.DB, u *model.User) *model.User {
placeholder
	if u.PasswordHash == "" {
		u.PasswordHash = "test-password-hash"
placeholder
	if u.Role == "" {
		u.Role = model.RoleUser
placeholder
	if u.Status == "" {
		u.Status = model.StatusActive
placeholder
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
placeholder
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
placeholder
	require.NoError(t, db.Create(u).Error, "create user")
	return u
placeholder

func mustCreateGroup(t *testing.T, db *gorm.DB, g *model.Group) *model.Group {
placeholder
	if g.Platform == "" {
		g.Platform = model.PlatformAnthropic
placeholder
	if g.Status == "" {
		g.Status = model.StatusActive
placeholder
	if g.SubscriptionType == "" {
		g.SubscriptionType = model.SubscriptionTypeStandard
placeholder
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now()
placeholder
	if g.UpdatedAt.IsZero() {
		g.UpdatedAt = g.CreatedAt
placeholder
	require.NoError(t, db.Create(g).Error, "create group")
	return g
placeholder

func mustCreateProxy(t *testing.T, db *gorm.DB, p *model.Proxy) *model.Proxy {
placeholder
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
		p.Status = model.StatusActive
placeholder
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
placeholder
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = p.CreatedAt
placeholder
	require.NoError(t, db.Create(p).Error, "create proxy")
	return p
placeholder

func mustCreateAccount(t *testing.T, db *gorm.DB, a *model.Account) *model.Account {
placeholder
	if a.Platform == "" {
		a.Platform = model.PlatformAnthropic
placeholder
	if a.Type == "" {
		a.Type = model.AccountTypeOAuth
placeholder
	if a.Status == "" {
		a.Status = model.StatusActive
placeholder
	if !a.Schedulable {
		a.Schedulable = true
placeholder
	if a.Credentials == nil {
		a.Credentials = model.JSONB{placeholder
placeholder
	if a.Extra == nil {
		a.Extra = model.JSONB{placeholder
placeholder
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
placeholder
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = a.CreatedAt
placeholder
	require.NoError(t, db.Create(a).Error, "create account")
	return a
placeholder

func mustCreateApiKey(t *testing.T, db *gorm.DB, k *model.ApiKey) *model.ApiKey {
placeholder
	if k.Status == "" {
		k.Status = model.StatusActive
placeholder
	if k.CreatedAt.IsZero() {
		k.CreatedAt = time.Now()
placeholder
	if k.UpdatedAt.IsZero() {
		k.UpdatedAt = k.CreatedAt
placeholder
	require.NoError(t, db.Create(k).Error, "create api key")
	return k
placeholder

func mustCreateRedeemCode(t *testing.T, db *gorm.DB, c *model.RedeemCode) *model.RedeemCode {
placeholder
	if c.Status == "" {
		c.Status = model.StatusUnused
placeholder
	if c.Type == "" {
		c.Type = model.RedeemTypeBalance
placeholder
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
placeholder
	require.NoError(t, db.Create(c).Error, "create redeem code")
	return c
placeholder

func mustCreateSubscription(t *testing.T, db *gorm.DB, s *model.UserSubscription) *model.UserSubscription {
placeholder
	if s.Status == "" {
		s.Status = model.SubscriptionStatusActive
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
	require.NoError(t, db.Create(s).Error, "create user subscription")
	return s
placeholder

func mustBindAccountToGroup(t *testing.T, db *gorm.DB, accountID, groupID int64, priority int) {
placeholder
	require.NoError(t, db.Create(&model.AccountGroup{
		AccountID: accountID,
		GroupID:   groupID,
		Priority:  priority,
placeholder).Error, "create account_group")
placeholder
