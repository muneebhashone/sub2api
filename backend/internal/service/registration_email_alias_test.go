//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeEmailForAliasDedup(t *testing.T) {
	cases := []struct {
		name  string
		email string
		want  string
placeholder{
		{"plain", "user@example.com", "user@example.com"placeholder,
		{"uppercase and spaces", "  User@Example.COM ", "user@example.com"placeholder,
		{"plus alias stripped", "user+tag@example.com", "user@example.com"placeholder,
		{"gmail plus alias", "someone+bulk294@gmail.com", "someone@gmail.com"placeholder,
		{"gmail dots removed", "some.one@gmail.com", "someone@gmail.com"placeholder,
		{"gmail dots and plus", "s.o.m.e+x@gmail.com", "some@gmail.com"placeholder,
		{"googlemail folded to gmail", "user@googlemail.com", "user@gmail.com"placeholder,
		{"non-gmail keeps dots", "first.last@qq.com", "first.last@qq.com"placeholder,
		{"invalid keeps lowered raw", "not-an-email", "not-an-email"placeholder,
placeholder
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, NormalizeEmailForAliasDedup(tc.email))
	placeholder)
placeholder
placeholder

func TestAliasDedupCandidateDomains(t *testing.T) {
	require.ElementsMatch(t, []string{"gmail.com", "googlemail.com"placeholder, aliasDedupCandidateDomains("user@gmail.com"))
	require.ElementsMatch(t, []string{"gmail.com", "googlemail.com"placeholder, aliasDedupCandidateDomains("user@googlemail.com"))
	require.Equal(t, []string{"qq.com"placeholder, aliasDedupCandidateDomains("user@qq.com"))
	require.Nil(t, aliasDedupCandidateDomains("not-an-email"))
placeholder

// aliasDedupRepoStub implements only the methods alias dedup uses; other
// UserRepository methods come from the embedded nil interface (a wrong call
// would panic, failing the test).
type aliasDedupRepoStub struct {
	UserRepository
	exists    bool
	existsErr error
	emails    []string
	listErr   error
	scanned   [][]string
placeholder

func (s *aliasDedupRepoStub) ExistsByEmail(context.Context, string) (bool, error) {
	return s.exists, s.existsErr
placeholder

func (s *aliasDedupRepoStub) ListEmailsByDomains(_ context.Context, domains []string) ([]string, error) {
	s.scanned = append(s.scanned, domains)
	return s.emails, s.listErr
placeholder

// exactOnlyRepoStub only supports the exact check (no alias-lookup capability).
type exactOnlyRepoStub struct {
	UserRepository
	exists bool
placeholder

func (s *exactOnlyRepoStub) ExistsByEmail(context.Context, string) (bool, error) {
	return s.exists, nil
placeholder

func TestExistsByEmailOrAlias(t *testing.T) {
	ctx := context.Background()

	t.Run("exact duplicate short-circuits", func(t *testing.T) {
		repo := &aliasDedupRepoStub{exists: trueplaceholder
		svc := &AuthService{userRepo: repoplaceholder
		got, err := svc.existsByEmailOrAlias(ctx, "user@gmail.com")
	placeholder
		require.True(t, got)
		require.Empty(t, repo.scanned, "no alias scan expected after exact hit")
placeholder)

	t.Run("plus alias variant detected", func(t *testing.T) {
		repo := &aliasDedupRepoStub{emails: []string{"someone+bulk294@gmail.com"placeholderplaceholder
		svc := &AuthService{userRepo: repoplaceholder
		got, err := svc.existsByEmailOrAlias(ctx, "Someone@gmail.com")
	placeholder
		require.True(t, got)
placeholder)

	t.Run("gmail dot variant detected", func(t *testing.T) {
		repo := &aliasDedupRepoStub{emails: []string{"some.one@gmail.com"placeholderplaceholder
		svc := &AuthService{userRepo: repoplaceholder
		got, err := svc.existsByEmailOrAlias(ctx, "someone@gmail.com")
	placeholder
		require.True(t, got)
placeholder)

	t.Run("gmail scans both gmail-family domains", func(t *testing.T) {
		repo := &aliasDedupRepoStub{placeholder
		svc := &AuthService{userRepo: repoplaceholder
		_, err := svc.existsByEmailOrAlias(ctx, "user@googlemail.com")
	placeholder
		require.Len(t, repo.scanned, 1)
		require.ElementsMatch(t, []string{"gmail.com", "googlemail.com"placeholder, repo.scanned[0])
placeholder)

	t.Run("different inbox allowed", func(t *testing.T) {
		repo := &aliasDedupRepoStub{emails: []string{"other@gmail.com"placeholderplaceholder
		svc := &AuthService{userRepo: repoplaceholder
		got, err := svc.existsByEmailOrAlias(ctx, "user@gmail.com")
	placeholder
		require.False(t, got)
placeholder)

	t.Run("list error fails closed", func(t *testing.T) {
		repo := &aliasDedupRepoStub{listErr: errors.New("db down")placeholder
		svc := &AuthService{userRepo: repoplaceholder
		_, err := svc.existsByEmailOrAlias(ctx, "user@gmail.com")
	placeholder
placeholder)

	t.Run("exact check error propagates", func(t *testing.T) {
		repo := &aliasDedupRepoStub{existsErr: errors.New("db down")placeholder
		svc := &AuthService{userRepo: repoplaceholder
		_, err := svc.existsByEmailOrAlias(ctx, "user@gmail.com")
	placeholder
placeholder)

	t.Run("repo without capability falls back to exact check", func(t *testing.T) {
		svc := &AuthService{userRepo: &exactOnlyRepoStub{exists: falseplaceholderplaceholder
		got, err := svc.existsByEmailOrAlias(ctx, "user@gmail.com")
	placeholder
		require.False(t, got)
placeholder)
placeholder
