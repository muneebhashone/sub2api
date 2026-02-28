package service

import "testing"

func TestGetBaseRPM(t *testing.T) {
	tests := []struct {
		name     string
		extra    map[string]any
		expected int
placeholder{
		{"nil extra", nil, 0placeholder,
		{"no key", map[string]any{placeholder, 0placeholder,
		{"zero", map[string]any{"base_rpm": 0placeholder, 0placeholder,
		{"int value", map[string]any{"base_rpm": 15placeholder, 15placeholder,
		{"float value", map[string]any{"base_rpm": 15.0placeholder, 15placeholder,
		{"string value", map[string]any{"base_rpm": "15"placeholder, 15placeholder,
		{"negative value", map[string]any{"base_rpm": -5placeholder, 0placeholder,
		{"int64 value", map[string]any{"base_rpm": int64(20)placeholder, 20placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{Extra: tt.extraplaceholder
			if got := a.GetBaseRPM(); got != tt.expected {
				t.Errorf("GetBaseRPM() = %d, want %d", got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestGetRPMStrategy(t *testing.T) {
	tests := []struct {
		name     string
		extra    map[string]any
		expected string
placeholder{
		{"nil extra", nil, "tiered"placeholder,
		{"no key", map[string]any{placeholder, "tiered"placeholder,
		{"tiered", map[string]any{"rpm_strategy": "tiered"placeholder, "tiered"placeholder,
		{"sticky_exempt", map[string]any{"rpm_strategy": "sticky_exempt"placeholder, "sticky_exempt"placeholder,
		{"invalid", map[string]any{"rpm_strategy": "foobar"placeholder, "tiered"placeholder,
		{"empty string fallback", map[string]any{"rpm_strategy": ""placeholder, "tiered"placeholder,
		{"numeric value fallback", map[string]any{"rpm_strategy": placeholder, "tiered"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{Extra: tt.extraplaceholder
			if got := a.GetRPMStrategy(); got != tt.expected {
				t.Errorf("GetRPMStrategy() = %q, want %q", got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestCheckRPMSchedulability(t *testing.T) {
	tests := []struct {
		name       string
		extra      map[string]any
		currentRPM int
		expected   WindowCostSchedulability
placeholder{
		{"disabled", map[string]any{placeholder, 100, WindowCostSchedulableplaceholder,
		{"green zone", map[string]any{"base_rpm": 15placeholder, 10, WindowCostSchedulableplaceholder,
		{"yellow zone tiered", map[string]any{"base_rpm": 15placeholder, 15, WindowCostStickyOnlyplaceholder,
		{"red zone tiered", map[string]any{"base_rpm": 15placeholder, 18, WindowCostNotSchedulableplaceholder,
		{"sticky_exempt at limit", map[string]any{"base_rpm": 15, "rpm_strategy": "sticky_exempt"placeholder, 15, WindowCostStickyOnlyplaceholder,
		{"sticky_exempt over limit", map[string]any{"base_rpm": 15, "rpm_strategy": "sticky_exempt"placeholder, 100, WindowCostStickyOnlyplaceholder,
		{"custom buffer", map[string]any{"base_rpm": 10, "rpm_sticky_buffer": 5placeholder, 14, WindowCostStickyOnlyplaceholder,
		{"custom buffer red", map[string]any{"base_rpm": 10, "rpm_sticky_buffer": 5placeholder, 15, WindowCostNotSchedulableplaceholder,
		{"base_rpm=1 green", map[string]any{"base_rpm": 1placeholder, 0, WindowCostSchedulableplaceholder,
		{"base_rpm=1 yellow (at limit)", map[string]any{"base_rpm": 1placeholder, 1, WindowCostStickyOnlyplaceholder,
		{"base_rpm=1 red (at limit+buffer)", map[string]any{"base_rpm": 1placeholder, 2, WindowCostNotSchedulableplaceholder,
		{"negative currentRPM", map[string]any{"base_rpm": 15placeholder, -1, WindowCostSchedulableplaceholder,
		{"base_rpm negative disabled", map[string]any{"base_rpm": -5placeholder, 10, WindowCostSchedulableplaceholder,
		{"very high currentRPM", map[string]any{"base_rpm": 10placeholder, 9999, WindowCostNotSchedulableplaceholder,
		{"sticky_exempt very high currentRPM", map[string]any{"base_rpm": 10, "rpm_strategy": "sticky_exempt"placeholder, 9999, WindowCostStickyOnlyplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{Extra: tt.extraplaceholder
			if got := a.CheckRPMSchedulability(tt.currentRPM); got != tt.expected {
				t.Errorf("CheckRPMSchedulability(%d) = %d, want %d", tt.currentRPM, got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder

func TestGetRPMStickyBuffer(t *testing.T) {
	tests := []struct {
		name     string
		extra    map[string]any
		expected int
placeholder{
		{"nil extra", nil, 0placeholder,
		{"no keys", map[string]any{placeholder, 0placeholder,
		{"base_rpm=0", map[string]any{"base_rpm": 0placeholder, 0placeholder,
		{"base_rpm=1 min buffer 1", map[string]any{"base_rpm": 1placeholder, 1placeholder,
		{"base_rpm=4 min buffer 1", map[string]any{"base_rpm": 4placeholder, 1placeholder,
		{"base_rpm=5 buffer 1", map[string]any{"base_rpm": 5placeholder, 1placeholder,
		{"base_rpm=10 buffer 2", map[string]any{"base_rpm": 10placeholder, 2placeholder,
		{"base_rpm=15 buffer 3", map[string]any{"base_rpm": 15placeholder, 3placeholder,
		{"base_rpm=100 buffer 20", map[string]any{"base_rpm": 100placeholder, 20placeholder,
		{"custom buffer=5", map[string]any{"base_rpm": 10, "rpm_sticky_buffer": 5placeholder, 5placeholder,
		{"custom buffer=0 fallback to default", map[string]any{"base_rpm": 10, "rpm_sticky_buffer": 0placeholder, 2placeholder,
		{"custom buffer negative fallback", map[string]any{"base_rpm": 10, "rpm_sticky_buffer": -1placeholder, 2placeholder,
		{"custom buffer with float", map[string]any{"base_rpm": 10, "rpm_sticky_buffer": float64(7)placeholder, 7placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{Extra: tt.extraplaceholder
			if got := a.GetRPMStickyBuffer(); got != tt.expected {
				t.Errorf("GetRPMStickyBuffer() = %d, want %d", got, tt.expected)
		placeholder
	placeholder)
placeholder
placeholder
