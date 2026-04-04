//go:build unit

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetModelPricing(t *testing.T) {
	ch := &Channel{
		ModelPricing: []ChannelModelPricing{
			{ID: 1, Models: []string{"claude-sonnet-4"placeholder, BillingMode: BillingModeToken, InputPrice: testPtrFloat64(3e-6)placeholder,
			{ID: 3, Models: []string{"gpt-5.1"placeholder, BillingMode: BillingModePerRequestplaceholder,
	placeholder,
placeholder

	tests := []struct {
		name    string
		model   string
		wantID  int64
		wantNil bool
placeholder{
		{"exact match", "claude-sonnet-4", 1, falseplaceholder,
		{"case insensitive", "Claude-Sonnet-4", 1, falseplaceholder,
		{"not found", "gemini-3.1-pro", 0, trueplaceholder,
		{"wildcard pattern not matched", "claude-opus-4-20250514", 0, trueplaceholder,
		{"per_request model", "gpt-5.1", 3, falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ch.GetModelPricing(tt.model)
			if tt.wantNil {
				require.Nil(t, result)
				return
		placeholder
			require.NotNil(t, result)
			require.Equal(t, tt.wantID, result.ID)
	placeholder)
placeholder
placeholder

func TestGetModelPricing_ReturnsCopy(t *testing.T) {
	ch := &Channel{
		ModelPricing: []ChannelModelPricing{
			{ID: 1, Models: []string{"claude-sonnet-4"placeholder, InputPrice: testPtrFloat64(3e-6)placeholder,
	placeholder,
placeholder

	result := ch.GetModelPricing("claude-sonnet-4")
	require.NotNil(t, result)

	// Modify the returned copy's slice — original should be unchanged
	result.Models = append(result.Models, "hacked")

	// Original should be unchanged
	require.Equal(t, 1, len(ch.ModelPricing[0].Models))
placeholder

func TestGetModelPricing_EmptyPricing(t *testing.T) {
	ch := &Channel{ModelPricing: nilplaceholder
	require.Nil(t, ch.GetModelPricing("any-model"))

	ch2 := &Channel{ModelPricing: []ChannelModelPricing{placeholderplaceholder
	require.Nil(t, ch2.GetModelPricing("any-model"))
placeholder

func TestGetIntervalForContext(t *testing.T) {
	p := &ChannelModelPricing{
		Intervals: []PricingInterval{
			{MinTokens: 0, MaxTokens: testPtrInt(128000), InputPrice: testPtrFloat64(1e-6)placeholder,
			{MinTokens: 128000, MaxTokens: nil, InputPrice: testPtrFloat64(2e-6)placeholder,
	placeholder,
placeholder

	tests := []struct {
		name      string
		tokens    int
		wantPrice *float64
		wantNil   bool
placeholder{
		{"first interval", 50000, testPtrFloat64(1e-6), falseplaceholder,
		// (min, max] — 128000 在第一个区间的 max，包含，所以匹配第一个
		{"boundary: max of first (inclusive)", 128000, testPtrFloat64(1e-6), falseplaceholder,
		// 128001 > 128000，匹配第二个区间
		{"boundary: just above first max", 128001, testPtrFloat64(2e-6), falseplaceholder,
		{"unbounded interval", 500000, testPtrFloat64(2e-6), falseplaceholder,
		// (0, max] — 0 不匹配任何区间（左开）
		{"zero tokens: no match", 0, nil, trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.GetIntervalForContext(tt.tokens)
			if tt.wantNil {
				require.Nil(t, result)
				return
		placeholder
			require.NotNil(t, result)
			require.InDelta(t, *tt.wantPrice, *result.InputPrice, 1e-12)
	placeholder)
placeholder
placeholder

func TestGetIntervalForContext_NoMatch(t *testing.T) {
	p := &ChannelModelPricing{
		Intervals: []PricingInterval{
			{MinTokens: 10000, MaxTokens: testPtrInt(50000)placeholder,
	placeholder,
placeholder
	require.Nil(t, p.GetIntervalForContext(5000))     // 5000 <= 10000, not > min
	require.Nil(t, p.GetIntervalForContext(10000))    // 10000 not > 10000 (left-open)
	require.NotNil(t, p.GetIntervalForContext(50000)) // 50000 <= 50000 (right-closed)
	require.Nil(t, p.GetIntervalForContext(50001))    // 50001 > 50000
placeholder

func TestGetIntervalForContext_Empty(t *testing.T) {
	p := &ChannelModelPricing{Intervals: nilplaceholder
	require.Nil(t, p.GetIntervalForContext(1000))
placeholder

func TestGetTierByLabel(t *testing.T) {
	p := &ChannelModelPricing{
		Intervals: []PricingInterval{
			{TierLabel: "1K", PerRequestPrice: testPtrFloat64(0.04)placeholder,
			{TierLabel: "2K", PerRequestPrice: testPtrFloat64(0.08)placeholder,
			{TierLabel: "HD", PerRequestPrice: testPtrFloat64(0.12)placeholder,
	placeholder,
placeholder

	tests := []struct {
		name    string
		label   string
		wantNil bool
		want    float64
placeholder{
		{"exact match", "1K", false, 0.04placeholder,
		{"case insensitive", "hd", false, 0.12placeholder,
		{"not found", "4K", true, 0placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.GetTierByLabel(tt.label)
			if tt.wantNil {
				require.Nil(t, result)
				return
		placeholder
			require.NotNil(t, result)
			require.InDelta(t, tt.want, *result.PerRequestPrice, 1e-12)
	placeholder)
placeholder
placeholder

func TestGetTierByLabel_Empty(t *testing.T) {
	p := &ChannelModelPricing{Intervals: nilplaceholder
	require.Nil(t, p.GetTierByLabel("1K"))
placeholder

func TestChannelClone(t *testing.T) {
	original := &Channel{
		ID:       1,
		Name:     "test",
		GroupIDs: []int64{10, 20placeholder,
		ModelPricing: []ChannelModelPricing{
			{
				ID:         100,
				Models:     []string{"model-a"placeholder,
				InputPrice: testPtrFloat64(5e-6),
		placeholder,
	placeholder,
placeholder

	cloned := original.Clone()
	require.NotNil(t, cloned)
	require.Equal(t, original.ID, cloned.ID)
	require.Equal(t, original.Name, cloned.Name)

	// Modify clone slices — original should not change
	cloned.GroupIDs[0] = 999
	require.Equal(t, int64(10), original.GroupIDs[0])

	cloned.ModelPricing[0].Models[0] = "hacked"
	require.Equal(t, "model-a", original.ModelPricing[0].Models[0])
placeholder

func TestChannelClone_Nil(t *testing.T) {
	var ch *Channel
	require.Nil(t, ch.Clone())
placeholder

func TestChannelModelPricingClone(t *testing.T) {
	original := ChannelModelPricing{
		Models: []string{"a", "b"placeholder,
		Intervals: []PricingInterval{
			{MinTokens: 0, TierLabel: "tier1"placeholder,
	placeholder,
placeholder

	cloned := original.Clone()

	// Modify clone slices — original unchanged
	cloned.Models[0] = "hacked"
	require.Equal(t, "a", original.Models[0])

	cloned.Intervals[0].TierLabel = "hacked"
	require.Equal(t, "tier1", original.Intervals[0].TierLabel)
placeholder

// --- BillingMode.IsValid ---

func TestBillingModeIsValid(t *testing.T) {
	tests := []struct {
		name string
		mode BillingMode
		want bool
placeholder{
		{"token", BillingModeToken, trueplaceholder,
		{"per_request", BillingModePerRequest, trueplaceholder,
		{"image", BillingModeImage, trueplaceholder,
		{"empty", BillingMode(""), trueplaceholder,
		{"unknown", BillingMode("unknown"), falseplaceholder,
		{"random", BillingMode("xyz"), falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.mode.IsValid())
	placeholder)
placeholder
placeholder

// --- Channel.IsActive ---

func TestChannelIsActive(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
placeholder{
		{"active", StatusActive, trueplaceholder,
		{"disabled", "disabled", falseplaceholder,
		{"empty", "", falseplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &Channel{Status: tt.statusplaceholder
			require.Equal(t, tt.want, ch.IsActive())
	placeholder)
placeholder
placeholder

// --- ChannelModelPricing.Clone edge cases ---

func TestChannelModelPricingClone_EdgeCases(t *testing.T) {
	t.Run("nil models", func(t *testing.T) {
		original := ChannelModelPricing{Models: nilplaceholder
		cloned := original.Clone()
		require.Nil(t, cloned.Models)
placeholder)

	t.Run("nil intervals", func(t *testing.T) {
		original := ChannelModelPricing{Intervals: nilplaceholder
		cloned := original.Clone()
		require.Nil(t, cloned.Intervals)
placeholder)

	t.Run("empty models", func(t *testing.T) {
		original := ChannelModelPricing{Models: []string{placeholderplaceholder
		cloned := original.Clone()
		require.NotNil(t, cloned.Models)
		require.Empty(t, cloned.Models)
placeholder)
placeholder

// --- Channel.Clone edge cases ---

func TestChannelClone_EdgeCases(t *testing.T) {
	t.Run("nil model mapping", func(t *testing.T) {
		original := &Channel{ID: 1, ModelMapping: nilplaceholder
		cloned := original.Clone()
		require.Nil(t, cloned.ModelMapping)
placeholder)

	t.Run("nil model pricing", func(t *testing.T) {
		original := &Channel{ID: 1, ModelPricing: nilplaceholder
		cloned := original.Clone()
		require.Nil(t, cloned.ModelPricing)
placeholder)

	t.Run("deep copy model mapping", func(t *testing.T) {
		original := &Channel{
			ID: 1,
			ModelMapping: map[string]map[string]string{
				"openai": {"gpt-4": "gpt-4-turbo"placeholder,
		placeholder,
	placeholder
		cloned := original.Clone()

		// Modify the cloned nested map
		cloned.ModelMapping["openai"]["gpt-4"] = "hacked"

		// Original must remain unchanged
		require.Equal(t, "gpt-4-turbo", original.ModelMapping["openai"]["gpt-4"])
placeholder)
placeholder

// --- ValidateIntervals ---

func TestValidateIntervals_Empty(t *testing.T) {
	require.NoError(t, ValidateIntervals(nil))
	require.NoError(t, ValidateIntervals([]PricingInterval{placeholder))
placeholder

func TestValidateIntervals_ValidIntervals(t *testing.T) {
	tests := []struct {
		name      string
		intervals []PricingInterval
placeholder{
		{
			name: "single bounded interval",
			intervals: []PricingInterval{
				{MinTokens: 0, MaxTokens: testPtrInt(128000), InputPrice: testPtrFloat64(1e-6)placeholder,
		placeholder,
	placeholder,
		{
			name: "two intervals with gap",
			intervals: []PricingInterval{
				{MinTokens: 0, MaxTokens: testPtrInt(100000), InputPrice: testPtrFloat64(1e-6)placeholder,
				{MinTokens: 128000, MaxTokens: nil, InputPrice: testPtrFloat64(2e-6)placeholder,
		placeholder,
	placeholder,
		{
			name: "two contiguous intervals",
			intervals: []PricingInterval{
				{MinTokens: 0, MaxTokens: testPtrInt(128000), InputPrice: testPtrFloat64(1e-6)placeholder,
				{MinTokens: 128000, MaxTokens: nil, InputPrice: testPtrFloat64(2e-6)placeholder,
		placeholder,
	placeholder,
		{
			name: "unsorted input (auto-sorted by validator)",
			intervals: []PricingInterval{
				{MinTokens: 128000, MaxTokens: nil, InputPrice: testPtrFloat64(2e-6)placeholder,
				{MinTokens: 0, MaxTokens: testPtrInt(128000), InputPrice: testPtrFloat64(1e-6)placeholder,
		placeholder,
	placeholder,
		{
			name: "single unbounded interval",
			intervals: []PricingInterval{
				{MinTokens: 0, MaxTokens: nil, InputPrice: testPtrFloat64(1e-6)placeholder,
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, ValidateIntervals(tt.intervals))
	placeholder)
placeholder
placeholder

func TestValidateIntervals_NegativeMinTokens(t *testing.T) {
	intervals := []PricingInterval{
		{MinTokens: -1, MaxTokens: testPtrInt(100), InputPrice: testPtrFloat64(1e-6)placeholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "min_tokens")
	require.Contains(t, err.Error(), ">= 0")
placeholder

func TestValidateIntervals_MaxTokensZero(t *testing.T) {
	intervals := []PricingInterval{
		{MinTokens: 0, MaxTokens: testPtrInt(0), InputPrice: testPtrFloat64(1e-6)placeholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "max_tokens")
	require.Contains(t, err.Error(), "> 0")
placeholder

func TestValidateIntervals_MaxLessThanMin(t *testing.T) {
	intervals := []PricingInterval{
		{MinTokens: 100, MaxTokens: testPtrInt(50), InputPrice: testPtrFloat64(1e-6)placeholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "max_tokens")
	require.Contains(t, err.Error(), "> min_tokens")
placeholder

func TestValidateIntervals_MaxEqualsMin(t *testing.T) {
	intervals := []PricingInterval{
		{MinTokens: 100, MaxTokens: testPtrInt(100), InputPrice: testPtrFloat64(1e-6)placeholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "max_tokens")
	require.Contains(t, err.Error(), "> min_tokens")
placeholder

func TestValidateIntervals_NegativePrice(t *testing.T) {
	negPrice := -0.01
	intervals := []PricingInterval{
		{MinTokens: 0, MaxTokens: testPtrInt(100), InputPrice: &negPriceplaceholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "input_price")
	require.Contains(t, err.Error(), ">= 0")
placeholder

func TestValidateIntervals_OverlappingIntervals(t *testing.T) {
	intervals := []PricingInterval{
		{MinTokens: 0, MaxTokens: testPtrInt(200), InputPrice: testPtrFloat64(1e-6)placeholder,
		{MinTokens: 100, MaxTokens: testPtrInt(300), InputPrice: testPtrFloat64(2e-6)placeholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "overlap")
placeholder

func TestValidateIntervals_UnboundedNotLast(t *testing.T) {
	intervals := []PricingInterval{
		{MinTokens: 0, MaxTokens: nil, InputPrice: testPtrFloat64(1e-6)placeholder,
		{MinTokens: 128000, MaxTokens: testPtrInt(256000), InputPrice: testPtrFloat64(2e-6)placeholder,
placeholder
	err := ValidateIntervals(intervals)
placeholder
	require.Contains(t, err.Error(), "unbounded")
	require.Contains(t, err.Error(), "last")
placeholder
