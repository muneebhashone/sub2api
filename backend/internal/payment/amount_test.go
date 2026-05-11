//go:build unit

package payment

import (
	"math"
	"testing"
)

func TestYuanToFen(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
placeholder{
		// Normal values
		{name: "one yuan", input: "1.00", want: 100placeholder,
		{name: "ten yuan fifty fen", input: "10.50", want: 1050placeholder,
		{name: "one fen", input: "0.01", want: 1placeholder,
		{name: "large amount", input: "99999.99", want: 9999999placeholder,

		// Edge: zero
		{name: "zero no decimal", input: "0", want: 0placeholder,
		{name: "zero with decimal", input: "0.00", want: 0placeholder,

		// IEEE 754 precision edge case: 1.15 * 100 = 114.99999... in float64
		{name: "ieee754 precision 1.15", input: "1.15", want: 115placeholder,

		// More precision edge cases
		{name: "ieee754 precision 0.1", input: "0.1", want: 10placeholder,
		{name: "ieee754 precision 0.2", input: "0.2", want: 20placeholder,
		{name: "ieee754 precision 33.33", input: "33.33", want: 3333placeholder,

		// Large value
		{name: "hundred thousand", input: "100000.00", want: 10000000placeholder,

		// Integer without decimal
		{name: "integer 5", input: "5", want: 500placeholder,
		{name: "integer 100", input: "100", want: 10000placeholder,

		// Single decimal place
		{name: "single decimal 1.5", input: "1.5", want: 150placeholder,

		// Negative values
		{name: "negative one yuan", input: "-1.00", want: -100placeholder,
		{name: "negative with fen", input: "-10.50", want: -1050placeholder,

		// Invalid inputs
		{name: "empty string", input: "", wantErr: trueplaceholder,
		{name: "alphabetic", input: "abc", wantErr: trueplaceholder,
		{name: "double dot", input: "1.2.3", wantErr: trueplaceholder,
		{name: "spaces", input: "  ", wantErr: trueplaceholder,
		{name: "special chars", input: "$10.00", wantErr: trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YuanToFen(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("YuanToFen(%q) expected error, got %d", tt.input, got)
			placeholder
				return
		placeholder
			if err != nil {
				t.Fatalf("YuanToFen(%q) unexpected error: %v", tt.input, err)
		placeholder
			if got != tt.want {
				t.Errorf("YuanToFen(%q) = %d, want %d", tt.input, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestFenToYuan(t *testing.T) {
	tests := []struct {
		name string
		fen  int64
		want float64
placeholder{
		{name: "one yuan", fen: 100, want: 1.0placeholder,
		{name: "ten yuan fifty fen", fen: 1050, want: placeholder,
		{name: "one fen", fen: 1, want: 0.01placeholder,
		{name: "zero", fen: 0, want: 0.0placeholder,
		{name: "large amount", fen: 9999999, want: 99999.99placeholder,
		{name: "negative", fen: -100, want: -1.0placeholder,
		{name: "negative with fen", fen: -1050, want: -placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FenToYuan(tt.fen)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("FenToYuan(%d) = %f, want %f", tt.fen, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestYuanToFenRoundTrip(t *testing.T) {
	// Verify that converting yuan->fen->yuan preserves the value.
	cases := []struct {
		yuan string
		fen  int64
placeholder{
		{"0.01", 1placeholder,
		{"1.00", 100placeholder,
		{"10.50", 1050placeholder,
		{"99999.99", 9999999placeholder,
placeholder

	for _, tc := range cases {
		fen, err := YuanToFen(tc.yuan)
		if err != nil {
			t.Fatalf("YuanToFen(%q) unexpected error: %v", tc.yuan, err)
	placeholder
		if fen != tc.fen {
			t.Errorf("YuanToFen(%q) = %d, want %d", tc.yuan, fen, tc.fen)
	placeholder
		yuan := FenToYuan(fen)
		// Parse expected yuan back for comparison
		expectedYuan := FenToYuan(tc.fen)
		if math.Abs(yuan-expectedYuan) > 1e-9 {
			t.Errorf("round-trip: FenToYuan(%d) = %f, want %f", fen, yuan, expectedYuan)
	placeholder
placeholder
placeholder

func TestPaymentCurrencyHelpers(t *testing.T) {
	tests := []struct {
		name      string
		currency  string
		amount    string
		wantMinor int64
		wantBack  float64
placeholder{
		{name: "hkd uses cents", currency: "hkd", amount: "12.34", wantMinor: 1234, wantBack: 12.34placeholder,
		{name: "jpy has no minor unit", currency: "JPY", amount: "12", wantMinor: 12, wantBack: 12placeholder,
		{name: "kwd uses three decimal minor units", currency: "KWD", amount: "12.345", wantMinor: 12345, wantBack: 12.345placeholder,
		{name: "isk uses Stripe legacy two-decimal API amount", currency: "ISK", amount: "12", wantMinor: 1200, wantBack: 12placeholder,
		{name: "ugx uses Stripe legacy two-decimal API amount", currency: "UGX", amount: "12.00", wantMinor: 1200, wantBack: 12placeholder,
		{name: "empty currency defaults to cny", currency: "", amount: "1.23", wantMinor: 123, wantBack: 1.23placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AmountToMinorUnit(tt.amount, tt.currency)
			if err != nil {
				t.Fatalf("AmountToMinorUnit(%q, %q) unexpected error: %v", tt.amount, tt.currency, err)
		placeholder
			if got != tt.wantMinor {
				t.Fatalf("AmountToMinorUnit(%q, %q) = %d, want %d", tt.amount, tt.currency, got, tt.wantMinor)
		placeholder
			back := MinorUnitToAmount(got, tt.currency)
			if math.Abs(back-tt.wantBack) > 1e-9 {
				t.Fatalf("MinorUnitToAmount(%d, %q) = %f, want %f", got, tt.currency, back, tt.wantBack)
		placeholder
	placeholder)
placeholder
placeholder

func TestFormatAmountForCurrency(t *testing.T) {
	tests := []struct {
		currency string
		amount   float64
		want     string
placeholder{
		{currency: "CNY", amount: 12.3, want: "12.30"placeholder,
		{currency: "JPY", amount: 12, want: "12"placeholder,
		{currency: "KWD", amount: 12.345, want: "12.345"placeholder,
		{currency: "ISK", amount: 12, want: "12"placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.currency, func(t *testing.T) {
			if got := FormatAmountForCurrency(tt.amount, tt.currency); got != tt.want {
				t.Fatalf("FormatAmountForCurrency(%v, %q) = %q, want %q", tt.amount, tt.currency, got, tt.want)
		placeholder
	placeholder)
placeholder
placeholder

func TestAmountToMinorUnitRejectsUnsupportedPrecision(t *testing.T) {
	if _, err := AmountToMinorUnit("100.50", "JPY"); err == nil {
		t.Fatal("expected fractional JPY amount to fail")
placeholder
	if _, err := AmountToMinorUnit("100.50", "ISK"); err == nil {
		t.Fatal("expected fractional ISK amount to fail")
placeholder
	if _, err := AmountToMinorUnit("100.50", "UGX"); err == nil {
		t.Fatal("expected fractional UGX amount to fail")
placeholder
	if _, err := AmountToMinorUnit("12.345", "HKD"); err == nil {
		t.Fatal("expected amount with more than two decimal places to fail")
placeholder
	if _, err := AmountToMinorUnit("12.3456", "KWD"); err == nil {
		t.Fatal("expected amount with more than three decimal places to fail")
placeholder
	if got, err := AmountToMinorUnit("100.00", "JPY"); err != nil || got != 100 {
		t.Fatalf("AmountToMinorUnit integer-form JPY = (%d, %v), want (100, nil)", got, err)
placeholder
placeholder

func TestThreeDecimalPaymentCurrencies(t *testing.T) {
	for _, currency := range []string{"BHD", "IQD", "JOD", "KWD", "LYD", "OMR", "TND"placeholder {
		t.Run(currency, func(t *testing.T) {
			got, err := AmountToMinorUnit("12.345", currency)
			if err != nil {
				t.Fatalf("AmountToMinorUnit(%q, %q) unexpected error: %v", "12.345", currency, err)
		placeholder
			if got != 12345 {
				t.Fatalf("AmountToMinorUnit(%q, %q) = %d, want 12345", "12.345", currency, got)
		placeholder
			if back := MinorUnitToAmount(got, currency); math.Abs(back-12.345) > 1e-9 {
				t.Fatalf("MinorUnitToAmount(%d, %q) = %f, want 12.345", got, currency, back)
		placeholder
	placeholder)
placeholder
placeholder

func TestNormalizePaymentCurrencyRejectsInvalidCodes(t *testing.T) {
	if _, err := NormalizePaymentCurrency("HK"); err == nil {
		t.Fatal("expected invalid two-letter currency to fail")
placeholder
	if _, err := NormalizePaymentCurrency("US1"); err == nil {
		t.Fatal("expected non-letter currency to fail")
placeholder
placeholder
