//go:build unit

package service

import (
	"math"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func timeConfig(periods ...ChannelTimePricingPeriod) *ChannelTimePricing {
	return &ChannelTimePricing{Timezone: "Asia/Shanghai", Periods: periodsplaceholder
placeholder

func onePeriod() []ChannelTimePricingPeriod {
	return []ChannelTimePricingPeriod{{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholderplaceholder
placeholder

func TestValidateChannelTimePricing(t *testing.T) {
	tests := []struct {
		name    string
		config  *ChannelTimePricing
		wantErr string
placeholder{
		{name: "nil disabled", config: nilplaceholder,
		{name: "empty disabled", config: &ChannelTimePricing{Timezone: "Asia/Shanghai"placeholderplaceholder,
		{name: "adjacent", config: timeConfig(
			ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholder,
			ChannelTimePricingPeriod{StartTime: "12:00", EndTime: "14:00", Multiplier: placeholder)placeholder,
		{name: "midnight split", config: timeConfig(
			ChannelTimePricingPeriod{StartTime: "22:00", EndTime: "00:00", Multiplier: 2placeholder,
			ChannelTimePricingPeriod{StartTime: "00:00", EndTime: "02:00", Multiplier: 2placeholder)placeholder,
		{name: "second precision", config: timeConfig(
			ChannelTimePricingPeriod{StartTime: "09:00:00", EndTime: "12:00:00", Multiplier: 2placeholder,
			ChannelTimePricingPeriod{StartTime: "14:00:00", EndTime: "18:00:00", Multiplier: 2placeholder)placeholder,
		{name: "second precision overlap", config: timeConfig(
			ChannelTimePricingPeriod{StartTime: "09:00:00", EndTime: "12:00:00", Multiplier: 2placeholder,
			ChannelTimePricingPeriod{StartTime: "11:59:59", EndTime: "14:00:00", Multiplier: 2placeholder), wantErr: "overlap"placeholder,
		{name: "empty timezone", config: &ChannelTimePricing{Periods: onePeriod()placeholder, wantErr: "timezone"placeholder,
		{name: "whitespace timezone", config: &ChannelTimePricing{Timezone: "  ", Periods: onePeriod()placeholder, wantErr: "timezone"placeholder,
		{name: "timezone", config: &ChannelTimePricing{Timezone: "UTC+8", Periods: onePeriod()placeholder, wantErr: "timezone"placeholder,
		{name: "format", config: timeConfig(ChannelTimePricingPeriod{StartTime: "9:00", EndTime: "12:00", Multiplier: 2placeholder), wantErr: "HH:mm"placeholder,
		{name: "equal midnight", config: timeConfig(ChannelTimePricingPeriod{StartTime: "00:00", EndTime: "00:00", Multiplier: 2placeholder), wantErr: "before"placeholder,
		{name: "cross midnight", config: timeConfig(ChannelTimePricingPeriod{StartTime: "22:00", EndTime: "02:00", Multiplier: 2placeholder), wantErr: "before"placeholder,
		{name: "overlap", config: timeConfig(
			ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholder,
			ChannelTimePricingPeriod{StartTime: "11:59", EndTime: "14:00", Multiplier: 2placeholder), wantErr: "overlap"placeholder,
		{name: "zero", config: timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 0placeholder), wantErr: "greater than 0"placeholder,
		{name: "minimum positive", config: timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 0.01placeholder)placeholder,
		{name: "tiny positive", config: timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 1e-12placeholder), wantErr: "at least 0.01"placeholder,
		{name: "below minimum", config: timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 0.001placeholder), wantErr: "at least 0.01"placeholder,
		{name: "three decimals", config: timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 1.001placeholder), wantErr: "decimal"placeholder,
		{name: "scaled overflow", config: timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: math.MaxFloat64placeholder), wantErr: "finite"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChannelTimePricing(tt.config)
			if tt.wantErr == "" {
			placeholder
				return
		placeholder
		placeholder
			require.True(t, strings.Contains(err.Error(), tt.wantErr), "error %q does not contain %q", err, tt.wantErr)
	placeholder)
placeholder
placeholder

func TestChannelTimePricingMultiplierAt(t *testing.T) {
	config := timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholder)
	tests := []struct {
		name string
		at   time.Time
		want float64
placeholder{
		{name: "Shanghai 08:59", at: time.Date(2026, 6, 29, 0, 59, 0, 0, time.UTC), want: 1placeholder,
		{name: "Shanghai 09:00", at: time.Date(2026, 6, 29, 1, 0, 0, 0, time.UTC), want: 2placeholder,
		{name: "Shanghai 11:59", at: time.Date(2026, 6, 29, 3, 59, 0, 0, time.UTC), want: 2placeholder,
		{name: "Shanghai 12:00", at: time.Date(2026, 6, 29, 4, 0, 0, 0, time.UTC), want: 1placeholder,
		{name: "Shanghai 14:00", at: time.Date(2026, 6, 29, 6, 0, 0, 0, time.UTC), want: 1placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, config.MultiplierAt(tt.at))
	placeholder)
placeholder

	newYork := &ChannelTimePricing{Timezone: "America/New_York", Periods: onePeriod()placeholder
	at := time.Date(2026, 6, 29, 14, 0, 0, 0, time.UTC)
	require.Equal(t, 1.0, config.MultiplierAt(at))
	require.Equal(t, 2.0, newYork.MultiplierAt(at))
placeholder

func TestChannelTimePricingMultiplierAtWeekdaysOnly(t *testing.T) {
	config := &ChannelTimePricing{
		Timezone:     "Asia/Shanghai",
		WeekdaysOnly: true,
		Periods:      onePeriod(),
placeholder

	tests := []struct {
		name string
		at   time.Time
		want float64
placeholder{
		{name: "Monday in configured timezone", at: time.Date(2026, 6, 29, 1, 0, 0, 0, time.UTC), want: 2placeholder,
		{name: "Saturday in configured timezone", at: time.Date(2026, 7, 4, 1, 0, 0, 0, time.UTC), want: 1placeholder,
		{name: "Sunday in configured timezone", at: time.Date(2026, 7, 5, 1, 0, 0, 0, time.UTC), want: 1placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, config.MultiplierAt(tt.at))
	placeholder)
placeholder
placeholder

func TestChannelTimePricingMultiplierAtSecondPrecision(t *testing.T) {
	config := timeConfig(ChannelTimePricingPeriod{StartTime: "09:00:30", EndTime: "09:00:45", Multiplier: 2placeholder)
	shanghai, err := time.LoadLocation("Asia/Shanghai")
placeholder

	tests := []struct {
		name string
		at   time.Time
		want float64
placeholder{
		{name: "before", at: time.Date(2026, 6, 29, 9, 0, 29, 0, shanghai), want: 1placeholder,
		{name: "start", at: time.Date(2026, 6, 29, 9, 0, 30, 0, shanghai), want: 2placeholder,
		{name: "last matching second", at: time.Date(2026, 6, 29, 9, 0, 44, 999_999_999, shanghai), want: 2placeholder,
		{name: "end", at: time.Date(2026, 6, 29, 9, 0, 45, 0, shanghai), want: 1placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, config.MultiplierAt(tt.at))
	placeholder)
placeholder
placeholder

func TestChannelTimePricingMultiplierAtMidnightSplit(t *testing.T) {
	config := timeConfig(
		ChannelTimePricingPeriod{StartTime: "22:00", EndTime: "00:00", Multiplier: 2placeholder,
		ChannelTimePricingPeriod{StartTime: "00:00", EndTime: "02:00", Multiplier: 3placeholder,
	)
	shanghai, err := time.LoadLocation("Asia/Shanghai")
placeholder
	tests := []struct {
		name string
		at   time.Time
		want float64
placeholder{
		{name: "23:59", at: time.Date(2026, 6, 29, 23, 59, 0, 0, shanghai), want: 2placeholder,
		{name: "next day 00:00", at: time.Date(2026, 6, 30, 0, 0, 0, 0, shanghai), want: 3placeholder,
		{name: "02:00", at: time.Date(2026, 6, 30, 2, 0, 0, 0, shanghai), want: 1placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, config.MultiplierAt(tt.at))
	placeholder)
placeholder
placeholder

func TestChannelTimePricingMultiplierAtDegradesForInvalidConfigurations(t *testing.T) {
	var nilConfig *ChannelTimePricing
	zeroTime := time.Time{placeholder
	validAt := time.Date(2026, 6, 29, 1, 0, 0, 0, time.UTC)

	require.Equal(t, 1.0, nilConfig.MultiplierAt(validAt))
	require.Equal(t, 1.0, timeConfig().MultiplierAt(validAt))
	require.Equal(t, 1.0, timeConfig(ChannelTimePricingPeriod{StartTime: "09:00", EndTime: "12:00", Multiplier: 2placeholder).MultiplierAt(zeroTime))
	require.Equal(t, 1.0, (&ChannelTimePricing{Periods: onePeriod()placeholder).MultiplierAt(validAt))
	require.Equal(t, 1.0, (&ChannelTimePricing{Timezone: "  ", Periods: onePeriod()placeholder).MultiplierAt(validAt))
	require.Equal(t, 1.0, (&ChannelTimePricing{Timezone: "UTC+8", Periods: onePeriod()placeholder).MultiplierAt(validAt))
	require.Equal(t, 1.0, timeConfig(ChannelTimePricingPeriod{StartTime: "22:00", EndTime: "02:00", Multiplier: 2placeholder).MultiplierAt(validAt))
placeholder

func TestChannelTimePricingRejectsLocalTimezone(t *testing.T) {
	config := &ChannelTimePricing{Timezone: "Local", Periods: onePeriod()placeholder

	err := validateChannelTimePricing(config)
placeholder
	require.Contains(t, err.Error(), "timezone")
	require.Equal(t, 1.0, config.MultiplierAt(time.Date(2026, 6, 29, 1, 0, 0, 0, time.UTC)))
placeholder
