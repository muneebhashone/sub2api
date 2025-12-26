package service

import "time"

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
	if s.IsExpired() {
		return 0
placeholder
	return int(time.Until(s.ExpiresAt).Hours() / 24)
placeholder

func (s *UserSubscription) IsWindowActivated() bool {
	return s.DailyWindowStart != nil || s.WeeklyWindowStart != nil || s.MonthlyWindowStart != nil
placeholder

func (s *UserSubscription) NeedsDailyReset() bool {
	if s.DailyWindowStart == nil {
		return false
placeholder
	return time.Since(*s.DailyWindowStart) >= 24*time.Hour
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
