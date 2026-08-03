package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type lockingRenewalRepo struct {
	userSubRepoNoop
	mu        sync.Mutex
	stale     UserSubscription
	current   UserSubscription
	lockReads int
placeholder

func (r *lockingRenewalRepo) ExistsByUserIDAndGroupID(context.Context, int64, int64) (bool, error) {
	return true, nil
placeholder

func (r *lockingRenewalRepo) GetByUserIDAndGroupID(context.Context, int64, int64) (*UserSubscription, error) {
	copy := r.stale
	return &copy, nil
placeholder

func (r *lockingRenewalRepo) GetByID(_ context.Context, _ int64) (*UserSubscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := r.current
	return &copy, nil
placeholder

func (r *lockingRenewalRepo) GetByIDForUpdate(_ context.Context, _ int64) (*UserSubscription, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lockReads++
	copy := r.current
	return &copy, nil
placeholder

func (r *lockingRenewalRepo) ExtendExpiry(_ context.Context, _ int64, expiresAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.current.ExpiresAt = expiresAt
	return nil
placeholder

func (r *lockingRenewalRepo) UpdateStatus(_ context.Context, _ int64, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.current.Status = status
	return nil
placeholder

func (r *lockingRenewalRepo) UpdateNotes(_ context.Context, _ int64, notes string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.current.Notes = notes
	return nil
placeholder

func (r *lockingRenewalRepo) Update(_ context.Context, sub *UserSubscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.current = *sub
	return nil
placeholder

func TestAssignOrExtendSubscriptionUsesLockedCurrentRow(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	lockedExpiry := now.AddDate(0, 0, 20)
	windowStart := now.Add(-24 * time.Hour)
	repo := &lockingRenewalRepo{
		stale: UserSubscription{ID: 7, UserID: 11, GroupID: 13, ExpiresAt: now.Add(-time.Hour), Status: SubscriptionStatusExpired, Notes: "stale"placeholder,
		current: UserSubscription{
			ID: 7, UserID: 11, GroupID: 13, StartsAt: now.AddDate(0, 0, -10), ExpiresAt: lockedExpiry,
			Status: SubscriptionStatusSuspended, Notes: "current", DailyWindowStart: &windowStart, DailyUsageUSD: 4,
	placeholder,
placeholder
	svc := NewSubscriptionService(&subscriptionGroupRepoStub{group: &Group{ID: 13, SubscriptionType: SubscriptionTypeSubscriptionplaceholderplaceholder, repo, nil, nil, nil)
	svc.now = func() time.Time { return now placeholder

	sub, extended, err := svc.AssignOrExtendSubscription(context.Background(), &AssignSubscriptionInput{
		UserID: 11, GroupID: 13, ValidityDays: 5, Notes: "renewed",
placeholder)

placeholder
	require.True(t, extended)
	require.Equal(t, 1, repo.lockReads)
	require.Equal(t, lockedExpiry.AddDate(0, 0, 5), sub.ExpiresAt)
	require.Equal(t, SubscriptionStatusActive, sub.Status)
	require.Equal(t, "current\nrenewed", sub.Notes)
	require.Equal(t, windowStart, *sub.DailyWindowStart)
	require.Equal(t, float64(4), sub.DailyUsageUSD)
placeholder

func TestAssignOrExtendSubscriptionSerializedRenewalsAccumulateDays(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	initialExpiry := now.AddDate(0, 0, 10)
	stale := UserSubscription{ID: 17, UserID: 21, GroupID: 23, StartsAt: now, ExpiresAt: initialExpiry, Status: SubscriptionStatusActiveplaceholder
	repo := &lockingRenewalRepo{stale: stale, current: staleplaceholder
	svc := NewSubscriptionService(&subscriptionGroupRepoStub{group: &Group{ID: 23, SubscriptionType: SubscriptionTypeSubscriptionplaceholderplaceholder, repo, nil, nil, nil)
	svc.now = func() time.Time { return now placeholder
	input := &AssignSubscriptionInput{UserID: 21, GroupID: 23, ValidityDays: 7placeholder

	_, _, err := svc.AssignOrExtendSubscription(context.Background(), input)
placeholder
	second, _, err := svc.AssignOrExtendSubscription(context.Background(), input)
placeholder

	require.Equal(t, 2, repo.lockReads)
	require.Equal(t, initialExpiry.AddDate(0, 0, 14), second.ExpiresAt)
placeholder

func TestAssignSubscriptionDoesNotReactivateRowSuspendedAfterStaleRead(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	windowStart := now.Add(-24 * time.Hour)
	current := UserSubscription{
		ID: 27, UserID: 31, GroupID: 33, StartsAt: now.AddDate(0, 0, -10), ExpiresAt: now.Add(-time.Hour),
		Status: SubscriptionStatusSuspended, Notes: "suspended", DailyWindowStart: &windowStart, DailyUsageUSD: 4,
placeholder
	repo := &lockingRenewalRepo{
		stale:   UserSubscription{ID: 27, UserID: 31, GroupID: 33, ExpiresAt: now.Add(-time.Hour), Status: SubscriptionStatusExpiredplaceholder,
		current: current,
placeholder
	svc := NewSubscriptionService(&subscriptionGroupRepoStub{group: &Group{ID: 33, SubscriptionType: SubscriptionTypeSubscriptionplaceholderplaceholder, repo, nil, nil, nil)
	svc.now = func() time.Time { return now placeholder

	sub, reused, err := svc.assignSubscriptionWithReuse(context.Background(), &AssignSubscriptionInput{
		UserID: 31, GroupID: 33, ValidityDays: 5, Notes: "renewed",
placeholder)

placeholder
	require.True(t, reused)
	require.Equal(t, 1, repo.lockReads)
	require.Equal(t, current, repo.current)
	require.Equal(t, SubscriptionStatusSuspended, sub.Status)
	require.Equal(t, current.ExpiresAt, sub.ExpiresAt)
	require.Equal(t, current.Notes, sub.Notes)
placeholder
