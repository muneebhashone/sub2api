package repository

import (
	"context"
	"database/sql"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

type sqlExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
placeholder

type sqlBeginner interface {
	sqlExecutor
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
placeholder

type groupRepository struct {
	client *dbent.Client
	sql    sqlExecutor
	begin  sqlBeginner
placeholder

func NewGroupRepository(client *dbent.Client, sqlDB *sql.DB) service.GroupRepository {
	return newGroupRepositoryWithSQL(client, sqlDB)
placeholder

func newGroupRepositoryWithSQL(client *dbent.Client, sqlq sqlExecutor) *groupRepository {
	var beginner sqlBeginner
	if b, ok := sqlq.(sqlBeginner); ok {
		beginner = b
placeholder
	return &groupRepository{client: client, sql: sqlq, begin: beginnerplaceholder
placeholder

func (r *groupRepository) Create(ctx context.Context, groupIn *service.Group) error {
	builder := r.client.Group.Create().
		SetName(groupIn.Name).
		SetDescription(groupIn.Description).
		SetPlatform(groupIn.Platform).
		SetRateMultiplier(groupIn.RateMultiplier).
		SetIsExclusive(groupIn.IsExclusive).
		SetStatus(groupIn.Status).
		SetSubscriptionType(groupIn.SubscriptionType).
		SetNillableDailyLimitUsd(groupIn.DailyLimitUSD).
		SetNillableWeeklyLimitUsd(groupIn.WeeklyLimitUSD).
		SetNillableMonthlyLimitUsd(groupIn.MonthlyLimitUSD)

	created, err := builder.Save(ctx)
	if err == nil {
		groupIn.ID = created.ID
		groupIn.CreatedAt = created.CreatedAt
		groupIn.UpdatedAt = created.UpdatedAt
placeholder
	return translatePersistenceError(err, nil, service.ErrGroupExists)
placeholder

func (r *groupRepository) GetByID(ctx context.Context, id int64) (*service.Group, error) {
	m, err := r.client.Group.Query().
		Where(group.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrGroupNotFound, nil)
placeholder

	out := groupEntityToService(m)
	count, _ := r.GetAccountCount(ctx, out.ID)
	out.AccountCount = count
	return out, nil
placeholder

func (r *groupRepository) Update(ctx context.Context, groupIn *service.Group) error {
	updated, err := r.client.Group.UpdateOneID(groupIn.ID).
		SetName(groupIn.Name).
		SetDescription(groupIn.Description).
		SetPlatform(groupIn.Platform).
		SetRateMultiplier(groupIn.RateMultiplier).
		SetIsExclusive(groupIn.IsExclusive).
		SetStatus(groupIn.Status).
		SetSubscriptionType(groupIn.SubscriptionType).
		SetNillableDailyLimitUsd(groupIn.DailyLimitUSD).
		SetNillableWeeklyLimitUsd(groupIn.WeeklyLimitUSD).
		SetNillableMonthlyLimitUsd(groupIn.MonthlyLimitUSD).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrGroupNotFound, service.ErrGroupExists)
placeholder
	groupIn.UpdatedAt = updated.UpdatedAt
	return nil
placeholder

func (r *groupRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.client.Group.Delete().Where(group.IDEQ(id)).Exec(ctx)
	return err
placeholder

func (r *groupRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.Group, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", nil)
placeholder

func (r *groupRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, status string, isExclusive *bool) ([]service.Group, *pagination.PaginationResult, error) {
	q := r.client.Group.Query()

	if platform != "" {
		q = q.Where(group.PlatformEQ(platform))
placeholder
	if status != "" {
		q = q.Where(group.StatusEQ(status))
placeholder
	if isExclusive != nil {
		q = q.Where(group.IsExclusiveEQ(*isExclusive))
placeholder

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	groups, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Asc(group.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	groupIDs := make([]int64, 0, len(groups))
	outGroups := make([]service.Group, 0, len(groups))
	for i := range groups {
		g := groupEntityToService(groups[i])
		outGroups = append(outGroups, *g)
		groupIDs = append(groupIDs, g.ID)
placeholder

	counts, err := r.loadAccountCounts(ctx, groupIDs)
	if err == nil {
		for i := range outGroups {
			outGroups[i].AccountCount = counts[outGroups[i].ID]
	placeholder
placeholder

	return outGroups, paginationResultFromTotal(int64(total), params), nil
placeholder

func (r *groupRepository) ListActive(ctx context.Context) ([]service.Group, error) {
	groups, err := r.client.Group.Query().
		Where(group.StatusEQ(service.StatusActive)).
		Order(dbent.Asc(group.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	groupIDs := make([]int64, 0, len(groups))
	outGroups := make([]service.Group, 0, len(groups))
	for i := range groups {
		g := groupEntityToService(groups[i])
		outGroups = append(outGroups, *g)
		groupIDs = append(groupIDs, g.ID)
placeholder

	counts, err := r.loadAccountCounts(ctx, groupIDs)
	if err == nil {
		for i := range outGroups {
			outGroups[i].AccountCount = counts[outGroups[i].ID]
	placeholder
placeholder

	return outGroups, nil
placeholder

func (r *groupRepository) ListActiveByPlatform(ctx context.Context, platform string) ([]service.Group, error) {
	groups, err := r.client.Group.Query().
		Where(group.StatusEQ(service.StatusActive), group.PlatformEQ(platform)).
		Order(dbent.Asc(group.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	groupIDs := make([]int64, 0, len(groups))
	outGroups := make([]service.Group, 0, len(groups))
	for i := range groups {
		g := groupEntityToService(groups[i])
		outGroups = append(outGroups, *g)
		groupIDs = append(groupIDs, g.ID)
placeholder

	counts, err := r.loadAccountCounts(ctx, groupIDs)
	if err == nil {
		for i := range outGroups {
			outGroups[i].AccountCount = counts[outGroups[i].ID]
	placeholder
placeholder

	return outGroups, nil
placeholder

func (r *groupRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.client.Group.Query().Where(group.NameEQ(name)).Exist(ctx)
placeholder

func (r *groupRepository) GetAccountCount(ctx context.Context, groupID int64) (int64, error) {
	var count int64
	if err := r.sql.QueryRowContext(ctx, "SELECT COUNT(*) FROM account_groups WHERE group_id = $1", groupID).Scan(&count); err != nil {
		return 0, err
placeholder
	return count, nil
placeholder

func (r *groupRepository) DeleteAccountGroupsByGroupID(ctx context.Context, groupID int64) (int64, error) {
	res, err := r.sql.ExecContext(ctx, "DELETE FROM account_groups WHERE group_id = $1", groupID)
	if err != nil {
		return 0, err
placeholder
	affected, _ := res.RowsAffected()
	return affected, nil
placeholder

func (r *groupRepository) DeleteCascade(ctx context.Context, id int64) ([]int64, error) {
	g, err := r.client.Group.Query().Where(group.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrGroupNotFound, nil)
placeholder
	groupSvc := groupEntityToService(g)

	exec := r.sql
	txClient := r.client
	var sqlTx *sql.Tx
	var txClientClose func() error

	if r.begin != nil {
		sqlTx, err = r.begin.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
	placeholder
		exec = sqlTx
		txClient = entClientFromSQLTx(sqlTx)
		txClientClose = txClient.Close
		defer func() { _ = sqlTx.Rollback() placeholder()
placeholder
	if txClientClose != nil {
		defer func() { _ = txClientClose() placeholder()
placeholder

	// Lock the group row to avoid concurrent writes while we cascade.
	var lockedID int64
	if err := exec.QueryRowContext(ctx, "SELECT id FROM groups WHERE id = $1 FOR UPDATE", id).Scan(&lockedID); err != nil {
		if errorsIsNoRows(err) {
			return nil, service.ErrGroupNotFound
	placeholder
		return nil, err
placeholder

	var affectedUserIDs []int64
	if groupSvc.IsSubscriptionType() {
		rows, err := exec.QueryContext(ctx, "SELECT user_id FROM user_subscriptions WHERE group_id = $1", id)
		if err != nil {
			return nil, err
	placeholder
		for rows.Next() {
			var userID int64
			if scanErr := rows.Scan(&userID); scanErr != nil {
				_ = rows.Close()
				return nil, scanErr
		placeholder
			affectedUserIDs = append(affectedUserIDs, userID)
	placeholder
		if err := rows.Close(); err != nil {
			return nil, err
	placeholder
		if err := rows.Err(); err != nil {
			return nil, err
	placeholder

		if _, err := exec.ExecContext(ctx, "DELETE FROM user_subscriptions WHERE group_id = $1", id); err != nil {
			return nil, err
	placeholder
placeholder

	// 2. Clear group_id for api keys bound to this group.
	if _, err := txClient.ApiKey.Update().
		Where(apikey.GroupIDEQ(id)).
		ClearGroupID().
		Save(ctx); err != nil {
		return nil, err
placeholder

	// 3. Remove the group id from users.allowed_groups array (legacy representation).
	// Phase 1 compatibility: also delete from user_allowed_groups join table when present.
	if _, err := exec.ExecContext(ctx, "DELETE FROM user_allowed_groups WHERE group_id = $1", id); err != nil {
		return nil, err
placeholder
	if _, err := exec.ExecContext(
		ctx,
		"UPDATE users SET allowed_groups = array_remove(allowed_groups, $1) WHERE $1 = ANY(allowed_groups)",
		id,
	); err != nil {
		return nil, err
placeholder

	// 4. Delete account_groups join rows.
	if _, err := exec.ExecContext(ctx, "DELETE FROM account_groups WHERE group_id = $1", id); err != nil {
		return nil, err
placeholder

	// 5. Soft-delete group itself.
	if _, err := txClient.Group.Delete().Where(group.IDEQ(id)).Exec(ctx); err != nil {
		return nil, err
placeholder

	if sqlTx != nil {
		if err := sqlTx.Commit(); err != nil {
			return nil, err
	placeholder
placeholder

	return affectedUserIDs, nil
placeholder

func (r *groupRepository) loadAccountCounts(ctx context.Context, groupIDs []int64) (map[int64]int64, error) {
	counts := make(map[int64]int64, len(groupIDs))
	if len(groupIDs) == 0 {
		return counts, nil
placeholder

	rows, err := r.sql.QueryContext(
		ctx,
		"SELECT group_id, COUNT(*) FROM account_groups WHERE group_id = ANY($1) GROUP BY group_id",
		pq.Array(groupIDs),
	)
	if err != nil {
		return nil, err
placeholder
	defer rows.Close()

	for rows.Next() {
		var groupID int64
		var count int64
		if err := rows.Scan(&groupID, &count); err != nil {
			return nil, err
	placeholder
		counts[groupID] = count
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder

	return counts, nil
placeholder

func entClientFromSQLTx(tx *sql.Tx) *dbent.Client {
	drv := entsql.NewDriver(dialect.Postgres, entsql.Conn{ExecQuerier: txplaceholder)
	return dbent.NewClient(dbent.Driver(drv))
placeholder

func errorsIsNoRows(err error) bool {
	return err == sql.ErrNoRows
placeholder
