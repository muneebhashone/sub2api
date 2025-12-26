//go:build integration

package repository

import (
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func mustCreateUser(t *testing.T, db *gorm.DB, u *userModel) *userModel {
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
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
placeholder
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = u.CreatedAt
placeholder
	require.NoError(t, db.Create(u).Error, "create user")
	return u
placeholder

func mustCreateGroup(t *testing.T, db *gorm.DB, g *groupModel) *groupModel {
placeholder
	if g.Platform == "" {
		g.Platform = service.PlatformAnthropic
placeholder
	if g.Status == "" {
		g.Status = service.StatusActive
placeholder
	if g.SubscriptionType == "" {
		g.SubscriptionType = service.SubscriptionTypeStandard
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

func mustCreateProxy(t *testing.T, db *gorm.DB, p *proxyModel) *proxyModel {
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
		p.Status = service.StatusActive
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

func mustCreateAccount(t *testing.T, db *gorm.DB, a *accountModel) *accountModel {
placeholder
	if a.Platform == "" {
		a.Platform = service.PlatformAnthropic
placeholder
	if a.Type == "" {
		a.Type = service.AccountTypeOAuth
placeholder
	if a.Status == "" {
		a.Status = service.StatusActive
placeholder
	if !a.Schedulable {
		a.Schedulable = true
placeholder
	if a.Credentials == nil {
		a.Credentials = datatypes.JSONMap{placeholder
placeholder
	if a.Extra == nil {
		a.Extra = datatypes.JSONMap{placeholder
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

func mustCreateApiKey(t *testing.T, db *gorm.DB, k *apiKeyModel) *apiKeyModel {
placeholder
	if k.Status == "" {
		k.Status = service.StatusActive
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

func mustCreateRedeemCode(t *testing.T, db *gorm.DB, c *redeemCodeModel) *redeemCodeModel {
placeholder
	if c.Status == "" {
		c.Status = service.StatusUnused
placeholder
	if c.Type == "" {
		c.Type = service.RedeemTypeBalance
placeholder
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now()
placeholder
	require.NoError(t, db.Create(c).Error, "create redeem code")
	return c
placeholder

func mustCreateSubscription(t *testing.T, db *gorm.DB, s *userSubscriptionModel) *userSubscriptionModel {
placeholder
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
	require.NoError(t, db.Create(s).Error, "create user subscription")
	return s
placeholder

func mustBindAccountToGroup(t *testing.T, db *gorm.DB, accountID, groupID int64, priority int) {
placeholder
	require.NoError(t, db.Create(&accountGroupModel{
		AccountID: accountID,
		GroupID:   groupID,
		Priority:  priority,
		CreatedAt: time.Now(),
placeholder).Error, "create account_group")
placeholder
