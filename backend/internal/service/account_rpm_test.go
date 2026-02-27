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
