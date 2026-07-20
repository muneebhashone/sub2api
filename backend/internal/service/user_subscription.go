package service

import "time"

const subscriptionDayDuration = 24 * time.Hour

type UserSubscription struct {
	ID      int64
	UserID  int64
	GroupID int64

	StartsAt  time.Time
	ExpiresAt time.Time
	Status    string

	DailyWindowStart   *time.Time
	WeeklyWindowStart  *time.Time
	MonthlyWindowStart *time.Time

	DailyUsageUSD   float64
	WeeklyUsageUSD  float64
	MonthlyUsageUSD float64

	AssignedBy *int64
	AssignedAt time.Time
	Notes      string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	User           *User
	Group          *Group
	AssignedByUser *User
placeholder

func (s *UserSubscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.ExpiresAt)
placeholder

func (s *UserSubscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
placeholder

func (s *UserSubscription) DaysRemaining() int {
	return s.daysRemainingAt(time.Now())
placeholder

func (s *UserSubscription) daysRemainingAt(now time.Time) int {
	remaining := s.ExpiresAt.Sub(now)
	if remaining <= 0 {
		return 0
placeholder

	days := int(remaining / subscriptionDayDuration)
	if remaining%subscriptionDayDuration != 0 {
		days++
placeholder
	return days
placeholder

func (s *UserSubscription) IsWindowActivated() bool {
	return s.DailyWindowStart != nil || s.WeeklyWindowStart != nil || s.MonthlyWindowStart != nil
placeholder

func (s *UserSubscription) HasOneTimeDailyQuota() bool {
	if s == nil || s.StartsAt.IsZero() || s.ExpiresAt.IsZero() {
		return false
placeholder
	return !s.ExpiresAt.After(s.StartsAt.AddDate(0, 0, 1))
placeholder

func (s *UserSubscription) NeedsDailyReset() bool {
	return s.NeedsDailyResetAt(time.Now())
placeholder

func (s *UserSubscription) NeedsDailyResetAt(now time.Time) bool {
	if s.DailyWindowStart == nil {
		return false
placeholder
	if s.HasOneTimeDailyQuota() {
		return false
placeholder
	return !now.Before(s.DailyWindowStart.Add(24 * time.Hour))
placeholder

func (s *UserSubscription) NeedsWeeklyReset() bool {
	if s.WeeklyWindowStart == nil {
		return false
placeholder
	return time.Since(*s.WeeklyWindowStart) >= 7*24*time.Hour
placeholder

func (s *UserSubscription) NeedsMonthlyReset() bool {
	if s.MonthlyWindowStart == nil {
		return false
placeholder
	return time.Since(*s.MonthlyWindowStart) >= 30*24*time.Hour
placeholder

func (s *UserSubscription) DailyResetTime() *time.Time {
	if s.DailyWindowStart == nil {
		return nil
placeholder
	if s.HasOneTimeDailyQuota() {
		t := s.ExpiresAt
		return &t
placeholder
	t := s.DailyWindowStart.Add(24 * time.Hour)
	return &t
placeholder

func (s *UserSubscription) WeeklyResetTime() *time.Time {
	if s.WeeklyWindowStart == nil {
		return nil
placeholder
	t := s.WeeklyWindowStart.Add(7 * 24 * time.Hour)
	return &t
placeholder

func (s *UserSubscription) MonthlyResetTime() *time.Time {
	if s.MonthlyWindowStart == nil {
		return nil
placeholder
	t := s.MonthlyWindowStart.Add(30 * 24 * time.Hour)
	return &t
placeholder

func (s *UserSubscription) CheckDailyLimit(group *Group, additionalCost float64) bool {
	if !group.HasDailyLimit() {
		return true
placeholder
	return s.DailyUsageUSD+additionalCost <= *group.DailyLimitUSD
placeholder

func (s *UserSubscription) CheckWeeklyLimit(group *Group, additionalCost float64) bool {
	if !group.HasWeeklyLimit() {
		return true
placeholder
	return s.WeeklyUsageUSD+additionalCost <= *group.WeeklyLimitUSD
placeholder

func (s *UserSubscription) CheckMonthlyLimit(group *Group, additionalCost float64) bool {
	if !group.HasMonthlyLimit() {
		return true
placeholder
	return s.MonthlyUsageUSD+additionalCost <= *group.MonthlyLimitUSD
placeholder

func (s *UserSubscription) CheckAllLimits(group *Group, additionalCost float64) (daily, weekly, monthly bool) {
	daily = s.CheckDailyLimit(group, additionalCost)
	weekly = s.CheckWeeklyLimit(group, additionalCost)
	monthly = s.CheckMonthlyLimit(group, additionalCost)
	return
placeholder
