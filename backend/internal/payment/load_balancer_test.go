//go:build unit

package payment

import (
	"encoding/json"
	"testing"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
)

func TestInstanceSupportsType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		supportedTypes string
		target         PaymentType
		expected       bool
placeholder{
		{
			name:           "exact match single type",
			supportedTypes: "alipay",
			target:         "alipay",
			expected:       true,
	placeholder,
		{
			name:           "no match single type",
			supportedTypes: "wxpay",
			target:         "alipay",
			expected:       false,
	placeholder,
		{
			name:           "match in comma-separated list",
			supportedTypes: "alipay,wxpay,stripe",
			target:         "wxpay",
			expected:       true,
	placeholder,
		{
			name:           "first in comma-separated list",
			supportedTypes: "alipay,wxpay",
			target:         "alipay",
			expected:       true,
	placeholder,
		{
			name:           "last in comma-separated list",
			supportedTypes: "alipay,wxpay,stripe",
			target:         "stripe",
			expected:       true,
	placeholder,
		{
			name:           "no match in comma-separated list",
			supportedTypes: "alipay,wxpay",
			target:         "stripe",
			expected:       false,
	placeholder,
		{
			name:           "empty target",
			supportedTypes: "alipay,wxpay",
			target:         "",
			expected:       false,
	placeholder,
		{
			name:           "types with spaces are trimmed",
			supportedTypes: " alipay , wxpay ",
			target:         "alipay",
			expected:       true,
	placeholder,
		{
			name:           "partial match should not succeed",
			supportedTypes: "alipay_direct",
			target:         "alipay",
			expected:       false,
	placeholder,
		{
			name:           "empty supported types means all supported",
			supportedTypes: "",
			target:         "alipay",
			expected:       true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := InstanceSupportsType(tt.supportedTypes, tt.target)
			if got != tt.expected {
				t.Fatalf("InstanceSupportsType(%q, %q) = %v, want %v", tt.supportedTypes, tt.target, got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// Helper to build test PaymentProviderInstance values
// ---------------------------------------------------------------------------

func testInstance(id int64, providerKey, limits string) *dbent.PaymentProviderInstance {
	return &dbent.PaymentProviderInstance{
		ID:          id,
		ProviderKey: providerKey,
		Limits:      limits,
		Enabled:     true,
placeholder
placeholder

// makeLimitsJSON builds a limits JSON string for a single payment type.
func makeLimitsJSON(paymentType string, cl ChannelLimits) string {
	m := map[string]ChannelLimits{paymentType: clplaceholder
	b, _ := json.Marshal(m)
	return string(b)
placeholder

// ---------------------------------------------------------------------------
// filterByLimits
// ---------------------------------------------------------------------------

func TestFilterByLimits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		candidates  []instanceCandidate
		paymentType PaymentType
		orderAmount float64
		wantIDs     []int64 // expected surviving instance IDs
placeholder{
		{
			name: "order below SingleMin is filtered out",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMin: 10placeholder)), dailyUsed: 0placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 5,
			wantIDs:     nil,
	placeholder,
		{
			name: "order at exact SingleMin boundary passes",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMin: 10placeholder)), dailyUsed: 0placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 10,
			wantIDs:     []int64{1placeholder,
	placeholder,
		{
			name: "order above SingleMax is filtered out",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMax: 100placeholder)), dailyUsed: 0placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 150,
			wantIDs:     nil,
	placeholder,
		{
			name: "order at exact SingleMax boundary passes",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMax: 100placeholder)), dailyUsed: 0placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 100,
			wantIDs:     []int64{1placeholder,
	placeholder,
		{
			name: "daily used + orderAmount exceeding dailyLimit is filtered out",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{DailyLimit: 500placeholder)), dailyUsed: 480placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 30,
			wantIDs:     nil, // 480+30=510 > 500
	placeholder,
		{
			name: "daily used + orderAmount equal to dailyLimit passes (strict greater-than)",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{DailyLimit: 500placeholder)), dailyUsed: 480placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 20,
			wantIDs:     []int64{1placeholder, // 480+20=500, 500 > 500 is false → passes
	placeholder,
		{
			name: "daily used + orderAmount below dailyLimit passes",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{DailyLimit: 500placeholder)), dailyUsed: 400placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 50,
			wantIDs:     []int64{1placeholder,
	placeholder,
		{
			name: "no limits configured passes through",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", ""), dailyUsed: 99999placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 100,
			wantIDs:     []int64{1placeholder,
	placeholder,
		{
			name: "multiple candidates with partial filtering",
			candidates: []instanceCandidate{
				// singleMax=50, order=80 → filtered out
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMax: 50placeholder)), dailyUsed: 0placeholder,
				// no limits → passes
				{inst: testInstance(2, "easypay", ""), dailyUsed: 0placeholder,
				// singleMin=100, order=80 → filtered out
				{inst: testInstance(3, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMin: 100placeholder)), dailyUsed: 0placeholder,
				// daily limit ok → passes (500+80=580 < 1000)
				{inst: testInstance(4, "easypay", makeLimitsJSON("alipay", ChannelLimits{DailyLimit: 1000placeholder)), dailyUsed: 500placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 80,
			wantIDs:     []int64{2, 4placeholder,
	placeholder,
		{
			name: "zero SingleMin and SingleMax means no single-transaction limit",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMin: 0, SingleMax: 0, DailyLimit: 0placeholder)), dailyUsed: 0placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 99999,
			wantIDs:     []int64{1placeholder,
	placeholder,
		{
			name: "all limits combined - order passes all checks",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMin: 10, SingleMax: 200, DailyLimit: 1000placeholder)), dailyUsed: 500placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 50,
			wantIDs:     []int64{1placeholder,
	placeholder,
		{
			name: "all limits combined - order fails SingleMin",
			candidates: []instanceCandidate{
				{inst: testInstance(1, "easypay", makeLimitsJSON("alipay", ChannelLimits{SingleMin: 10, SingleMax: 200, DailyLimit: 1000placeholder)), dailyUsed: 500placeholder,
		placeholder,
			paymentType: "alipay",
			orderAmount: 5,
			wantIDs:     nil,
	placeholder,
		{
			name:        "empty candidates returns empty",
			candidates:  nil,
			paymentType: "alipay",
			orderAmount: 10,
			wantIDs:     nil,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := filterByLimits(tt.candidates, tt.paymentType, tt.orderAmount)
			gotIDs := make([]int64, len(got))
			for i, c := range got {
				gotIDs[i] = c.inst.ID
		placeholder
			if !int64SliceEqual(gotIDs, tt.wantIDs) {
				t.Fatalf("filterByLimits() returned IDs %v, want %v", gotIDs, tt.wantIDs)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// pickLeastAmount
// ---------------------------------------------------------------------------

func TestPickLeastAmount(t *testing.T) {
	t.Parallel()

	t.Run("picks candidate with lowest dailyUsed", func(t *testing.T) {
		t.Parallel()
		candidates := []instanceCandidate{
			{inst: testInstance(1, "easypay", ""), dailyUsed: 300placeholder,
			{inst: testInstance(2, "easypay", ""), dailyUsed: 100placeholder,
			{inst: testInstance(3, "easypay", ""), dailyUsed: 200placeholder,
	placeholder
		got := pickLeastAmount(candidates)
		if got.inst.ID != 2 {
			t.Fatalf("pickLeastAmount() picked instance %d, want 2", got.inst.ID)
	placeholder
placeholder)

	t.Run("with equal dailyUsed picks the first one", func(t *testing.T) {
		t.Parallel()
		candidates := []instanceCandidate{
			{inst: testInstance(1, "easypay", ""), dailyUsed: 100placeholder,
			{inst: testInstance(2, "easypay", ""), dailyUsed: 100placeholder,
			{inst: testInstance(3, "easypay", ""), dailyUsed: 200placeholder,
	placeholder
		got := pickLeastAmount(candidates)
		if got.inst.ID != 1 {
			t.Fatalf("pickLeastAmount() picked instance %d, want 1 (first with lowest)", got.inst.ID)
	placeholder
placeholder)

	t.Run("single candidate returns that candidate", func(t *testing.T) {
		t.Parallel()
		candidates := []instanceCandidate{
			{inst: testInstance(42, "easypay", ""), dailyUsed: 999placeholder,
	placeholder
		got := pickLeastAmount(candidates)
		if got.inst.ID != 42 {
			t.Fatalf("pickLeastAmount() picked instance %d, want 42", got.inst.ID)
	placeholder
placeholder)

	t.Run("zero usage among non-zero picks zero", func(t *testing.T) {
		t.Parallel()
		candidates := []instanceCandidate{
			{inst: testInstance(1, "easypay", ""), dailyUsed: 500placeholder,
			{inst: testInstance(2, "easypay", ""), dailyUsed: 0placeholder,
			{inst: testInstance(3, "easypay", ""), dailyUsed: 300placeholder,
	placeholder
		got := pickLeastAmount(candidates)
		if got.inst.ID != 2 {
			t.Fatalf("pickLeastAmount() picked instance %d, want 2", got.inst.ID)
	placeholder
placeholder)
placeholder

// ---------------------------------------------------------------------------
// getInstanceChannelLimits
// ---------------------------------------------------------------------------

func TestGetInstanceChannelLimits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inst        *dbent.PaymentProviderInstance
		paymentType PaymentType
		want        ChannelLimits
placeholder{
		{
			name:        "empty limits string returns zero ChannelLimits",
			inst:        testInstance(1, "easypay", ""),
			paymentType: "alipay",
			want:        ChannelLimits{placeholder,
	placeholder,
		{
			name:        "invalid JSON returns zero ChannelLimits",
			inst:        testInstance(1, "easypay", "not-json{"),
			paymentType: "alipay",
			want:        ChannelLimits{placeholder,
	placeholder,
		{
			name: "valid JSON with matching payment type",
			inst: testInstance(1, "easypay",
				`{"alipay":{"singleMin":5,"singleMax":200,"dailyLimit":1000placeholderplaceholder`),
			paymentType: "alipay",
			want:        ChannelLimits{SingleMin: 5, SingleMax: 200, DailyLimit: 1000placeholder,
	placeholder,
		{
			name: "payment type not in limits returns zero ChannelLimits",
			inst: testInstance(1, "easypay",
				`{"alipay":{"singleMin":5,"singleMax":200placeholderplaceholder`),
			paymentType: "wxpay",
			want:        ChannelLimits{placeholder,
	placeholder,
		{
			name: "stripe provider uses stripe lookup key regardless of payment type",
			inst: testInstance(1, "stripe",
				`{"stripe":{"singleMin":10,"singleMax":500,"dailyLimit":5000placeholderplaceholder`),
			paymentType: "alipay",
			want:        ChannelLimits{SingleMin: 10, SingleMax: 500, DailyLimit: 5000placeholder,
	placeholder,
		{
			name: "stripe provider ignores payment type key even if present",
			inst: testInstance(1, "stripe",
				`{"stripe":{"singleMin":10,"singleMax":500placeholder,"alipay":{"singleMin":1,"singleMax":100placeholderplaceholder`),
			paymentType: "alipay",
			want:        ChannelLimits{SingleMin: 10, SingleMax: 500placeholder,
	placeholder,
		{
			name: "non-stripe provider uses payment type as lookup key",
			inst: testInstance(1, "easypay",
				`{"alipay":{"singleMin":5placeholder,"wxpay":{"singleMin":10placeholderplaceholder`),
			paymentType: "wxpay",
			want:        ChannelLimits{SingleMin: 10placeholder,
	placeholder,
		{
			name: "valid JSON with partial limits (only dailyLimit)",
			inst: testInstance(1, "easypay",
				`{"alipay":{"dailyLimit":800placeholderplaceholder`),
			paymentType: "alipay",
			want:        ChannelLimits{DailyLimit: 800placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getInstanceChannelLimits(tt.inst, tt.paymentType)
			if got != tt.want {
				t.Fatalf("getInstanceChannelLimits() = %+v, want %+v", got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// startOfDay
// ---------------------------------------------------------------------------

func TestStartOfDay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   time.Time
		want time.Time
placeholder{
		{
			name: "midday returns midnight of same day",
			in:   time.Date(2025, 6, 15, 14, 30, 45, 123456789, time.UTC),
			want: time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC),
	placeholder,
		{
			name: "midnight returns same time",
			in:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			want: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	placeholder,
		{
			name: "last second of day returns midnight of same day",
			in:   time.Date(2025, 12, 31, 23, 59, 59, 999999999, time.UTC),
			want: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
	placeholder,
		{
			name: "preserves timezone location",
			in:   time.Date(2025, 3, 10, 15, 0, 0, 0, time.FixedZone("CST", 8*3600)),
			want: time.Date(2025, 3, 10, 0, 0, 0, 0, time.FixedZone("CST", 8*3600)),
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := startOfDay(tt.in)
			if !got.Equal(tt.want) {
				t.Fatalf("startOfDay(%v) = %v, want %v", tt.in, got, tt.want)
		placeholder
			// Also verify location is preserved.
			if got.Location().String() != tt.want.Location().String() {
				t.Fatalf("startOfDay() location = %v, want %v", got.Location(), tt.want.Location())
		placeholder
	placeholder)
placeholder
placeholder

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// int64SliceEqual compares two int64 slices for equality.
// Both nil and empty slices are treated as equal.
func int64SliceEqual(a, b []int64) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
placeholder
	if len(a) != len(b) {
		return false
placeholder
	for i := range a {
		if a[i] != b[i] {
			return false
	placeholder
placeholder
	return true
placeholder
