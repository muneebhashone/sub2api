package repository

import (
	"context"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/setting"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

const (
	simpleModeAdminConcurrencyUpgradeKey = "simple_mode_admin_concurrency_upgraded_30"
	simpleModeLegacyAdminConcurrency     = 5
	simpleModeTargetAdminConcurrency     = 30
)

func ensureSimpleModeAdminConcurrency(ctx context.Context, client *dbent.Client) error {
	if client == nil {
		return fmt.Errorf("nil ent client")
placeholder

	upgraded, err := client.Setting.Query().Where(setting.KeyEQ(simpleModeAdminConcurrencyUpgradeKey)).Exist(ctx)
	if err != nil {
		return fmt.Errorf("check admin concurrency upgrade marker: %w", err)
placeholder
	if upgraded {
		return nil
placeholder

	if _, err := client.User.Update().
		Where(
			dbuser.RoleEQ(service.RoleAdmin),
			dbuser.ConcurrencyEQ(simpleModeLegacyAdminConcurrency),
		).
		SetConcurrency(simpleModeTargetAdminConcurrency).
		Save(ctx); err != nil {
		return fmt.Errorf("upgrade simple mode admin concurrency: %w", err)
placeholder

	now := time.Now()
	if err := client.Setting.Create().
		SetKey(simpleModeAdminConcurrencyUpgradeKey).
		SetValue(now.Format(time.RFC3339)).
		SetUpdatedAt(now).
		OnConflictColumns(setting.FieldKey).
		UpdateNewValues().
		Exec(ctx); err != nil {
		return fmt.Errorf("persist admin concurrency upgrade marker: %w", err)
placeholder

	return nil
placeholder
