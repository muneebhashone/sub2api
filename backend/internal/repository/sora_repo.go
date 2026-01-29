package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Wei-Shaw/sub2api/ent"
	dbsoraaccount "github.com/Wei-Shaw/sub2api/ent/soraaccount"
	dbsoracachefile "github.com/Wei-Shaw/sub2api/ent/soracachefile"
	dbsoratask "github.com/Wei-Shaw/sub2api/ent/soratask"
	dbsorausagestat "github.com/Wei-Shaw/sub2api/ent/sorausagestat"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	entsql "entgo.io/ent/dialect/sql"
)

// SoraAccount

type soraAccountRepository struct {
	client *ent.Client
placeholder

func NewSoraAccountRepository(client *ent.Client) service.SoraAccountRepository {
	return &soraAccountRepository{client: clientplaceholder
placeholder

func (r *soraAccountRepository) GetByAccountID(ctx context.Context, accountID int64) (*service.SoraAccount, error) {
	if accountID <= 0 {
		return nil, nil
placeholder
	acc, err := r.client.SoraAccount.Query().Where(dbsoraaccount.AccountIDEQ(accountID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, err
placeholder
	return mapSoraAccount(acc), nil
placeholder

func (r *soraAccountRepository) GetByAccountIDs(ctx context.Context, accountIDs []int64) (map[int64]*service.SoraAccount, error) {
	if len(accountIDs) == 0 {
		return map[int64]*service.SoraAccount{placeholder, nil
placeholder
	records, err := r.client.SoraAccount.Query().Where(dbsoraaccount.AccountIDIn(accountIDs...)).All(ctx)
	if err != nil {
		return nil, err
placeholder
	result := make(map[int64]*service.SoraAccount, len(records))
	for _, record := range records {
		if record == nil {
			continue
	placeholder
		result[record.AccountID] = mapSoraAccount(record)
placeholder
	return result, nil
placeholder

func (r *soraAccountRepository) Upsert(ctx context.Context, accountID int64, updates map[string]any) error {
	if accountID <= 0 {
		return errors.New("invalid account_id")
placeholder
	acc, err := r.client.SoraAccount.Query().Where(dbsoraaccount.AccountIDEQ(accountID)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return err
placeholder
	if acc == nil {
		builder := r.client.SoraAccount.Create().SetAccountID(accountID)
		applySoraAccountUpdates(builder.Mutation(), updates)
		return builder.Exec(ctx)
placeholder
	updater := r.client.SoraAccount.UpdateOneID(acc.ID)
	applySoraAccountUpdates(updater.Mutation(), updates)
	return updater.Exec(ctx)
placeholder

func applySoraAccountUpdates(m *ent.SoraAccountMutation, updates map[string]any) {
	if updates == nil {
		return
placeholder
	for key, val := range updates {
		switch key {
		case "access_token":
			if v, ok := val.(string); ok {
				m.SetAccessToken(v)
		placeholder
		case "session_token":
			if v, ok := val.(string); ok {
				m.SetSessionToken(v)
		placeholder
		case "refresh_token":
			if v, ok := val.(string); ok {
				m.SetRefreshToken(v)
		placeholder
		case "client_id":
			if v, ok := val.(string); ok {
				m.SetClientID(v)
		placeholder
		case "email":
			if v, ok := val.(string); ok {
				m.SetEmail(v)
		placeholder
		case "username":
			if v, ok := val.(string); ok {
				m.SetUsername(v)
		placeholder
		case "remark":
			if v, ok := val.(string); ok {
				m.SetRemark(v)
		placeholder
		case "plan_type":
			if v, ok := val.(string); ok {
				m.SetPlanType(v)
		placeholder
		case "plan_title":
			if v, ok := val.(string); ok {
				m.SetPlanTitle(v)
		placeholder
		case "subscription_end":
			if v, ok := val.(time.Time); ok {
				m.SetSubscriptionEnd(v)
		placeholder
			if v, ok := val.(*time.Time); ok && v != nil {
				m.SetSubscriptionEnd(*v)
		placeholder
		case "sora_supported":
			if v, ok := val.(bool); ok {
				m.SetSoraSupported(v)
		placeholder
		case "sora_invite_code":
			if v, ok := val.(string); ok {
				m.SetSoraInviteCode(v)
		placeholder
		case "sora_redeemed_count":
			if v, ok := val.(int); ok {
				m.SetSoraRedeemedCount(v)
		placeholder
		case "sora_remaining_count":
			if v, ok := val.(int); ok {
				m.SetSoraRemainingCount(v)
		placeholder
		case "sora_total_count":
			if v, ok := val.(int); ok {
				m.SetSoraTotalCount(v)
		placeholder
		case "sora_cooldown_until":
			if v, ok := val.(time.Time); ok {
				m.SetSoraCooldownUntil(v)
		placeholder
			if v, ok := val.(*time.Time); ok && v != nil {
				m.SetSoraCooldownUntil(*v)
		placeholder
		case "cooled_until":
			if v, ok := val.(time.Time); ok {
				m.SetCooledUntil(v)
		placeholder
			if v, ok := val.(*time.Time); ok && v != nil {
				m.SetCooledUntil(*v)
		placeholder
		case "image_enabled":
			if v, ok := val.(bool); ok {
				m.SetImageEnabled(v)
		placeholder
		case "video_enabled":
			if v, ok := val.(bool); ok {
				m.SetVideoEnabled(v)
		placeholder
		case "image_concurrency":
			if v, ok := val.(int); ok {
				m.SetImageConcurrency(v)
		placeholder
		case "video_concurrency":
			if v, ok := val.(int); ok {
				m.SetVideoConcurrency(v)
		placeholder
		case "is_expired":
			if v, ok := val.(bool); ok {
				m.SetIsExpired(v)
		placeholder
	placeholder
placeholder
placeholder

func mapSoraAccount(acc *ent.SoraAccount) *service.SoraAccount {
	if acc == nil {
		return nil
placeholder
	return &service.SoraAccount{
		AccountID:          acc.AccountID,
		AccessToken:        derefString(acc.AccessToken),
		SessionToken:       derefString(acc.SessionToken),
		RefreshToken:       derefString(acc.RefreshToken),
		ClientID:           derefString(acc.ClientID),
		Email:              derefString(acc.Email),
		Username:           derefString(acc.Username),
		Remark:             derefString(acc.Remark),
		UseCount:           acc.UseCount,
		PlanType:           derefString(acc.PlanType),
		PlanTitle:          derefString(acc.PlanTitle),
		SubscriptionEnd:    acc.SubscriptionEnd,
		SoraSupported:      acc.SoraSupported,
		SoraInviteCode:     derefString(acc.SoraInviteCode),
		SoraRedeemedCount:  acc.SoraRedeemedCount,
		SoraRemainingCount: acc.SoraRemainingCount,
		SoraTotalCount:     acc.SoraTotalCount,
		SoraCooldownUntil:  acc.SoraCooldownUntil,
		CooledUntil:        acc.CooledUntil,
		ImageEnabled:       acc.ImageEnabled,
		VideoEnabled:       acc.VideoEnabled,
		ImageConcurrency:   acc.ImageConcurrency,
		VideoConcurrency:   acc.VideoConcurrency,
		IsExpired:          acc.IsExpired,
		CreatedAt:          acc.CreatedAt,
		UpdatedAt:          acc.UpdatedAt,
placeholder
placeholder

func mapSoraUsageStat(stat *ent.SoraUsageStat) *service.SoraUsageStat {
	if stat == nil {
		return nil
placeholder
	return &service.SoraUsageStat{
		AccountID:             stat.AccountID,
		ImageCount:            stat.ImageCount,
		VideoCount:            stat.VideoCount,
		ErrorCount:            stat.ErrorCount,
		LastErrorAt:           stat.LastErrorAt,
		TodayImageCount:       stat.TodayImageCount,
		TodayVideoCount:       stat.TodayVideoCount,
		TodayErrorCount:       stat.TodayErrorCount,
		TodayDate:             stat.TodayDate,
		ConsecutiveErrorCount: stat.ConsecutiveErrorCount,
		CreatedAt:             stat.CreatedAt,
		UpdatedAt:             stat.UpdatedAt,
placeholder
placeholder

func mapSoraCacheFile(file *ent.SoraCacheFile) *service.SoraCacheFile {
	if file == nil {
		return nil
placeholder
	return &service.SoraCacheFile{
		ID:          int64(file.ID),
		TaskID:      derefString(file.TaskID),
		AccountID:   file.AccountID,
		UserID:      file.UserID,
		MediaType:   file.MediaType,
		OriginalURL: file.OriginalURL,
		CachePath:   file.CachePath,
		CacheURL:    file.CacheURL,
		SizeBytes:   file.SizeBytes,
		CreatedAt:   file.CreatedAt,
placeholder
placeholder

// SoraUsageStat

type soraUsageStatRepository struct {
	client *ent.Client
	sql    sqlExecutor
placeholder

func NewSoraUsageStatRepository(client *ent.Client, sqlDB *sql.DB) service.SoraUsageStatRepository {
	return &soraUsageStatRepository{client: client, sql: sqlDBplaceholder
placeholder

func (r *soraUsageStatRepository) RecordSuccess(ctx context.Context, accountID int64, isVideo bool) error {
	if accountID <= 0 {
		return nil
placeholder
	field := "image_count"
	todayField := "today_image_count"
	if isVideo {
		field = "video_count"
		todayField = "today_video_count"
placeholder
	today := time.Now().UTC().Truncate(24 * time.Hour)
	query := "INSERT INTO sora_usage_stats (account_id, " + field + ", " + todayField + ", today_date, consecutive_error_count, created_at, updated_at) " +
		"VALUES ($1, 1, 1, $2, 0, NOW(), NOW()) " +
		"ON CONFLICT (account_id) DO UPDATE SET " +
		field + " = sora_usage_stats." + field + " + 1, " +
		todayField + " = CASE WHEN sora_usage_stats.today_date = $2 THEN sora_usage_stats." + todayField + " + 1 ELSE 1 END, " +
		"today_date = $2, consecutive_error_count = 0, updated_at = NOW()"
	_, err := r.sql.ExecContext(ctx, query, accountID, today)
	return err
placeholder

func (r *soraUsageStatRepository) RecordError(ctx context.Context, accountID int64) (int, error) {
	if accountID <= 0 {
		return 0, nil
placeholder
	today := time.Now().UTC().Truncate(24 * time.Hour)
	query := "INSERT INTO sora_usage_stats (account_id, error_count, today_error_count, today_date, consecutive_error_count, last_error_at, created_at, updated_at) " +
		"VALUES ($1, 1, 1, $2, 1, NOW(), NOW(), NOW()) " +
		"ON CONFLICT (account_id) DO UPDATE SET " +
		"error_count = sora_usage_stats.error_count + 1, " +
		"today_error_count = CASE WHEN sora_usage_stats.today_date = $2 THEN sora_usage_stats.today_error_count + 1 ELSE 1 END, " +
		"today_date = $2, consecutive_error_count = sora_usage_stats.consecutive_error_count + 1, last_error_at = NOW(), updated_at = NOW() " +
		"RETURNING consecutive_error_count"
	var consecutive int
	err := scanSingleRow(ctx, r.sql, query, []any{accountID, todayplaceholder, &consecutive)
	if err != nil {
		return 0, err
placeholder
	return consecutive, nil
placeholder

func (r *soraUsageStatRepository) ResetConsecutiveErrors(ctx context.Context, accountID int64) error {
	if accountID <= 0 {
		return nil
placeholder
	err := r.client.SoraUsageStat.Update().Where(dbsorausagestat.AccountIDEQ(accountID)).
		SetConsecutiveErrorCount(0).
		Exec(ctx)
	if ent.IsNotFound(err) {
		return nil
placeholder
	return err
placeholder

func (r *soraUsageStatRepository) GetByAccountID(ctx context.Context, accountID int64) (*service.SoraUsageStat, error) {
	if accountID <= 0 {
		return nil, nil
placeholder
	stat, err := r.client.SoraUsageStat.Query().Where(dbsorausagestat.AccountIDEQ(accountID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, err
placeholder
	return mapSoraUsageStat(stat), nil
placeholder

func (r *soraUsageStatRepository) GetByAccountIDs(ctx context.Context, accountIDs []int64) (map[int64]*service.SoraUsageStat, error) {
	if len(accountIDs) == 0 {
		return map[int64]*service.SoraUsageStat{placeholder, nil
placeholder
	stats, err := r.client.SoraUsageStat.Query().Where(dbsorausagestat.AccountIDIn(accountIDs...)).All(ctx)
	if err != nil {
		return nil, err
placeholder
	result := make(map[int64]*service.SoraUsageStat, len(stats))
	for _, stat := range stats {
		if stat == nil {
			continue
	placeholder
		result[stat.AccountID] = mapSoraUsageStat(stat)
placeholder
	return result, nil
placeholder

func (r *soraUsageStatRepository) List(ctx context.Context, params pagination.PaginationParams) ([]*service.SoraUsageStat, *pagination.PaginationResult, error) {
	query := r.client.SoraUsageStat.Query()
	total, err := query.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder
	stats, err := query.Order(ent.Desc(dbsorausagestat.FieldUpdatedAt)).
		Limit(params.Limit()).
		Offset(params.Offset()).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder
	result := make([]*service.SoraUsageStat, 0, len(stats))
	for _, stat := range stats {
		result = append(result, mapSoraUsageStat(stat))
placeholder
	return result, paginationResultFromTotal(int64(total), params), nil
placeholder

// SoraTask

type soraTaskRepository struct {
	client *ent.Client
placeholder

func NewSoraTaskRepository(client *ent.Client) service.SoraTaskRepository {
	return &soraTaskRepository{client: clientplaceholder
placeholder

func (r *soraTaskRepository) Create(ctx context.Context, task *service.SoraTask) error {
	if task == nil {
		return nil
placeholder
	builder := r.client.SoraTask.Create().
		SetTaskID(task.TaskID).
		SetAccountID(task.AccountID).
		SetModel(task.Model).
		SetPrompt(task.Prompt).
		SetStatus(task.Status).
		SetProgress(task.Progress).
		SetRetryCount(task.RetryCount)
	if task.ResultURLs != "" {
		builder.SetResultUrls(task.ResultURLs)
placeholder
	if task.ErrorMessage != "" {
		builder.SetErrorMessage(task.ErrorMessage)
placeholder
	if task.CreatedAt.IsZero() {
		builder.SetCreatedAt(time.Now())
placeholder else {
		builder.SetCreatedAt(task.CreatedAt)
placeholder
	if task.CompletedAt != nil {
		builder.SetCompletedAt(*task.CompletedAt)
placeholder
	return builder.Exec(ctx)
placeholder

func (r *soraTaskRepository) UpdateStatus(ctx context.Context, taskID string, status string, progress float64, resultURLs string, errorMessage string, completedAt *time.Time) error {
	if taskID == "" {
		return nil
placeholder
	builder := r.client.SoraTask.Update().Where(dbsoratask.TaskIDEQ(taskID)).
		SetStatus(status).
		SetProgress(progress)
	if resultURLs != "" {
		builder.SetResultUrls(resultURLs)
placeholder
	if errorMessage != "" {
		builder.SetErrorMessage(errorMessage)
placeholder
	if completedAt != nil {
		builder.SetCompletedAt(*completedAt)
placeholder
	_, err := builder.Save(ctx)
	if ent.IsNotFound(err) {
		return nil
placeholder
	return err
placeholder

// SoraCacheFile

type soraCacheFileRepository struct {
	client *ent.Client
placeholder

func NewSoraCacheFileRepository(client *ent.Client) service.SoraCacheFileRepository {
	return &soraCacheFileRepository{client: clientplaceholder
placeholder

func (r *soraCacheFileRepository) Create(ctx context.Context, file *service.SoraCacheFile) error {
	if file == nil {
		return nil
placeholder
	builder := r.client.SoraCacheFile.Create().
		SetAccountID(file.AccountID).
		SetUserID(file.UserID).
		SetMediaType(file.MediaType).
		SetOriginalURL(file.OriginalURL).
		SetCachePath(file.CachePath).
		SetCacheURL(file.CacheURL).
		SetSizeBytes(file.SizeBytes)
	if file.TaskID != "" {
		builder.SetTaskID(file.TaskID)
placeholder
	if file.CreatedAt.IsZero() {
		builder.SetCreatedAt(time.Now())
placeholder else {
		builder.SetCreatedAt(file.CreatedAt)
placeholder
	return builder.Exec(ctx)
placeholder

func (r *soraCacheFileRepository) ListOldest(ctx context.Context, limit int) ([]*service.SoraCacheFile, error) {
	if limit <= 0 {
		return []*service.SoraCacheFile{placeholder, nil
placeholder
	records, err := r.client.SoraCacheFile.Query().
		Order(dbsoracachefile.ByCreatedAt(entsql.OrderAsc())).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	result := make([]*service.SoraCacheFile, 0, len(records))
	for _, record := range records {
		if record == nil {
			continue
	placeholder
		result = append(result, mapSoraCacheFile(record))
placeholder
	return result, nil
placeholder

func (r *soraCacheFileRepository) DeleteByIDs(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
placeholder
	_, err := r.client.SoraCacheFile.Delete().Where(dbsoracachefile.IDIn(ids...)).Exec(ctx)
	return err
placeholder
