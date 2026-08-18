//go:build unit

package service

import "testing"

func TestApplyMonitorUpdate_ProviderOnlySwitchToGrokUsesDefaultModel(t *testing.T) {
	grok := MonitorProviderGrok
	existing := &ChannelMonitor{
		Provider:        MonitorProviderOpenAI,
		APIMode:         MonitorAPIModeResponses,
		Endpoint:        "https://api.openai.com/v1",
		PrimaryModel:    "gpt-5",
		IntervalSeconds: 60,
placeholder

	err := applyMonitorUpdate(existing, ChannelMonitorUpdateParams{Provider: &grokplaceholder)
	if err != nil {
		t.Fatalf("provider-only switch to Grok failed: %v", err)
placeholder
	if existing.PrimaryModel != MonitorDefaultGrokModel {
		t.Fatalf("expected Grok default model %q, got %q", MonitorDefaultGrokModel, existing.PrimaryModel)
placeholder
	if existing.APIMode != MonitorAPIModeChatCompletions {
		t.Fatalf("expected Grok API mode %q, got %q", MonitorAPIModeChatCompletions, existing.APIMode)
placeholder
placeholder

func TestApplyMonitorUpdate_SwitchToGrokPreservesExplicitModel(t *testing.T) {
	grok := MonitorProviderGrok
	explicitModel := "grok-4.3"
	existing := &ChannelMonitor{
		Provider:        MonitorProviderOpenAI,
		APIMode:         MonitorAPIModeChatCompletions,
		Endpoint:        "https://api.openai.com/v1",
		PrimaryModel:    "gpt-5",
		IntervalSeconds: 60,
placeholder

	err := applyMonitorUpdate(existing, ChannelMonitorUpdateParams{
		Provider:     &grok,
		PrimaryModel: &explicitModel,
placeholder)
	if err != nil {
		t.Fatalf("switch to Grok with explicit model failed: %v", err)
placeholder
	if existing.PrimaryModel != explicitModel {
		t.Fatalf("expected explicit model %q, got %q", explicitModel, existing.PrimaryModel)
placeholder
placeholder

func TestApplyMonitorUpdate_SameGrokProviderDoesNotResetExistingModel(t *testing.T) {
	grok := MonitorProviderGrok
	existing := &ChannelMonitor{
		Provider:        MonitorProviderGrok,
		APIMode:         MonitorAPIModeChatCompletions,
		Endpoint:        "https://api.x.ai",
		PrimaryModel:    "grok-4.3",
		IntervalSeconds: 60,
placeholder

	err := applyMonitorUpdate(existing, ChannelMonitorUpdateParams{Provider: &grokplaceholder)
	if err != nil {
		t.Fatalf("same-provider Grok update failed: %v", err)
placeholder
	if existing.PrimaryModel != "grok-4.3" {
		t.Fatalf("same-provider update reset existing model to %q", existing.PrimaryModel)
placeholder
placeholder

func TestApplyMonitorUpdate_SwitchToGrokRejectsResponsesMode(t *testing.T) {
	grok := MonitorProviderGrok
	responses := MonitorAPIModeResponses
	existing := &ChannelMonitor{
		Provider:        MonitorProviderOpenAI,
		APIMode:         MonitorAPIModeChatCompletions,
		PrimaryModel:    "gpt-5",
		IntervalSeconds: 60,
placeholder

	err := applyMonitorUpdate(existing, ChannelMonitorUpdateParams{
		Provider: &grok,
		APIMode:  &responses,
placeholder)
	if err == nil {
		t.Fatal("Grok responses mode should remain unsupported")
placeholder
placeholder
