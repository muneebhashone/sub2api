package service

import (
	"context"
	"encoding/json"
	"testing"
)

func TestGetOpsAdvancedSettings_DefaultHidesOpenAITokenStats(t *testing.T) {
	repo := newRuntimeSettingRepoStub()
	svc := &OpsService{settingRepo: repoplaceholder

	cfg, err := svc.GetOpsAdvancedSettings(context.Background())
	if err != nil {
		t.Fatalf("GetOpsAdvancedSettings() error = %v", err)
placeholder
	if cfg.DisplayOpenAITokenStats {
		t.Fatalf("DisplayOpenAITokenStats = true, want false by default")
placeholder
	if !cfg.DisplayAlertEvents {
		t.Fatalf("DisplayAlertEvents = false, want true by default")
placeholder
	if repo.setCalls != 1 {
		t.Fatalf("expected defaults to be persisted once, got %d", repo.setCalls)
placeholder
placeholder

func TestUpdateOpsAdvancedSettings_PersistsOpenAITokenStatsVisibility(t *testing.T) {
	repo := newRuntimeSettingRepoStub()
	svc := &OpsService{settingRepo: repoplaceholder

	cfg := defaultOpsAdvancedSettings()
	cfg.DisplayOpenAITokenStats = true
	cfg.DisplayAlertEvents = false

	updated, err := svc.UpdateOpsAdvancedSettings(context.Background(), cfg)
	if err != nil {
		t.Fatalf("UpdateOpsAdvancedSettings() error = %v", err)
placeholder
	if !updated.DisplayOpenAITokenStats {
		t.Fatalf("DisplayOpenAITokenStats = false, want true")
placeholder
	if updated.DisplayAlertEvents {
		t.Fatalf("DisplayAlertEvents = true, want false")
placeholder

	reloaded, err := svc.GetOpsAdvancedSettings(context.Background())
	if err != nil {
		t.Fatalf("GetOpsAdvancedSettings() after update error = %v", err)
placeholder
	if !reloaded.DisplayOpenAITokenStats {
		t.Fatalf("reloaded DisplayOpenAITokenStats = false, want true")
placeholder
	if reloaded.DisplayAlertEvents {
		t.Fatalf("reloaded DisplayAlertEvents = true, want false")
placeholder
placeholder

func TestGetOpsAdvancedSettings_BackfillsNewDisplayFlagsFromDefaults(t *testing.T) {
	repo := newRuntimeSettingRepoStub()
	svc := &OpsService{settingRepo: repoplaceholder

	legacyCfg := map[string]any{
		"data_retention": map[string]any{
			"cleanup_enabled":               false,
			"cleanup_schedule":              "0 2 * * *",
			"error_log_retention_days":      30,
			"minute_metrics_retention_days": 30,
			"hourly_metrics_retention_days": 30,
	placeholder,
		"aggregation": map[string]any{
			"aggregation_enabled": false,
	placeholder,
		"ignore_count_tokens_errors":    true,
		"ignore_context_canceled":       true,
		"ignore_no_available_accounts":  false,
		"ignore_invalid_api_key_errors": false,
		"auto_refresh_enabled":          false,
		"auto_refresh_interval_seconds": 30,
placeholder
	raw, err := json.Marshal(legacyCfg)
	if err != nil {
		t.Fatalf("marshal legacy config: %v", err)
placeholder
	repo.values[SettingKeyOpsAdvancedSettings] = string(raw)

	cfg, err := svc.GetOpsAdvancedSettings(context.Background())
	if err != nil {
		t.Fatalf("GetOpsAdvancedSettings() error = %v", err)
placeholder
	if cfg.DisplayOpenAITokenStats {
		t.Fatalf("DisplayOpenAITokenStats = true, want false default backfill")
placeholder
	if !cfg.DisplayAlertEvents {
		t.Fatalf("DisplayAlertEvents = false, want true default backfill")
placeholder
placeholder
