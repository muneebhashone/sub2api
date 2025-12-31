// Package repository 实现数据访问层（Repository Pattern）。
//
// 该包提供了与数据库交互的所有操作，包括 CRUD、复杂查询和批量操作。
// 采用 Repository 模式将数据访问逻辑与业务逻辑分离，便于测试和维护。
//
// 主要特性：
//   - 使用 Ent ORM 进行类型安全的数据库操作
//   - 对于复杂查询（如批量更新、聚合统计）使用原生 SQL
//   - 提供统一的错误翻译机制，将数据库错误转换为业务错误
//   - 支持软删除，所有查询自动过滤已删除记录
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	dbaccount "github.com/Wei-Shaw/sub2api/ent/account"
	dbaccountgroup "github.com/Wei-Shaw/sub2api/ent/accountgroup"
	dbgroup "github.com/Wei-Shaw/sub2api/ent/group"
	dbpredicate "github.com/Wei-Shaw/sub2api/ent/predicate"
	dbproxy "github.com/Wei-Shaw/sub2api/ent/proxy"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"

	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
)

// accountRepository 实现 service.AccountRepository 接口。
// 提供 AI API 账户的完整数据访问功能。
//
// 设计说明：
//   - client: Ent 客户端，用于类型安全的 ORM 操作
//   - sql: 原生 SQL 执行器，用于复杂查询和批量操作
type accountRepository struct {
	client *dbent.Client // Ent ORM 客户端
	sql    sqlExecutor   // 原生 SQL 执行接口
placeholder

// NewAccountRepository 创建账户仓储实例。
// 这是对外暴露的构造函数，返回接口类型以便于依赖注入。
func NewAccountRepository(client *dbent.Client, sqlDB *sql.DB) service.AccountRepository {
	return newAccountRepositoryWithSQL(client, sqlDB)
placeholder

// newAccountRepositoryWithSQL 是内部构造函数，支持依赖注入 SQL 执行器。
// 这种设计便于单元测试时注入 mock 对象。
func newAccountRepositoryWithSQL(client *dbent.Client, sqlq sqlExecutor) *accountRepository {
	return &accountRepository{client: client, sql: sqlqplaceholder
placeholder

func (r *accountRepository) Create(ctx context.Context, account *service.Account) error {
	if account == nil {
		return service.ErrAccountNilInput
placeholder

	builder := r.client.Account.Create().
		SetName(account.Name).
		SetPlatform(account.Platform).
		SetType(account.Type).
		SetCredentials(normalizeJSONMap(account.Credentials)).
		SetExtra(normalizeJSONMap(account.Extra)).
		SetConcurrency(account.Concurrency).
		SetPriority(account.Priority).
		SetStatus(account.Status).
		SetErrorMessage(account.ErrorMessage).
		SetSchedulable(account.Schedulable)

	if account.ProxyID != nil {
		builder.SetProxyID(*account.ProxyID)
placeholder
	if account.LastUsedAt != nil {
		builder.SetLastUsedAt(*account.LastUsedAt)
placeholder
	if account.RateLimitedAt != nil {
		builder.SetRateLimitedAt(*account.RateLimitedAt)
placeholder
	if account.RateLimitResetAt != nil {
		builder.SetRateLimitResetAt(*account.RateLimitResetAt)
placeholder
	if account.OverloadUntil != nil {
		builder.SetOverloadUntil(*account.OverloadUntil)
placeholder
	if account.SessionWindowStart != nil {
		builder.SetSessionWindowStart(*account.SessionWindowStart)
placeholder
	if account.SessionWindowEnd != nil {
		builder.SetSessionWindowEnd(*account.SessionWindowEnd)
placeholder
	if account.SessionWindowStatus != "" {
		builder.SetSessionWindowStatus(account.SessionWindowStatus)
placeholder

	created, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrAccountNotFound, nil)
placeholder

	account.ID = created.ID
	account.CreatedAt = created.CreatedAt
	account.UpdatedAt = created.UpdatedAt
	return nil
placeholder

func (r *accountRepository) GetByID(ctx context.Context, id int64) (*service.Account, error) {
	m, err := r.client.Account.Query().Where(dbaccount.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrAccountNotFound, nil)
placeholder

	accounts, err := r.accountsToService(ctx, []*dbent.Account{mplaceholder)
	if err != nil {
		return nil, err
placeholder
	if len(accounts) == 0 {
		return nil, service.ErrAccountNotFound
placeholder
	return &accounts[0], nil
placeholder

// ExistsByID 检查指定 ID 的账号是否存在。
// 相比 GetByID，此方法性能更优，因为：
//   - 使用 Exist() 方法生成 SELECT EXISTS 查询，只返回布尔值
//   - 不加载完整的账号实体及其关联数据（Groups、Proxy 等）
//   - 适用于删除前的存在性检查等只需判断有无的场景
func (r *accountRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	exists, err := r.client.Account.Query().Where(dbaccount.IDEQ(id)).Exist(ctx)
	if err != nil {
		return false, err
placeholder
	return exists, nil
placeholder

func (r *accountRepository) GetByCRSAccountID(ctx context.Context, crsAccountID string) (*service.Account, error) {
	if crsAccountID == "" {
		return nil, nil
placeholder

	// 使用 sqljson.ValueEQ 生成 JSON 路径过滤，避免手写 SQL 片段导致语法兼容问题。
	m, err := r.client.Account.Query().
		Where(func(s *entsql.Selector) {
			s.Where(sqljson.ValueEQ(dbaccount.FieldExtra, crsAccountID, sqljson.Path("crs_account_id")))
	placeholder).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, err
placeholder

	accounts, err := r.accountsToService(ctx, []*dbent.Account{mplaceholder)
	if err != nil {
		return nil, err
placeholder
	if len(accounts) == 0 {
		return nil, nil
placeholder
	return &accounts[0], nil
placeholder

func (r *accountRepository) Update(ctx context.Context, account *service.Account) error {
	if account == nil {
		return nil
placeholder

	builder := r.client.Account.UpdateOneID(account.ID).
		SetName(account.Name).
		SetPlatform(account.Platform).
		SetType(account.Type).
		SetCredentials(normalizeJSONMap(account.Credentials)).
		SetExtra(normalizeJSONMap(account.Extra)).
		SetConcurrency(account.Concurrency).
		SetPriority(account.Priority).
		SetStatus(account.Status).
		SetErrorMessage(account.ErrorMessage).
		SetSchedulable(account.Schedulable)

	if account.ProxyID != nil {
		builder.SetProxyID(*account.ProxyID)
placeholder else {
		builder.ClearProxyID()
placeholder
	if account.LastUsedAt != nil {
		builder.SetLastUsedAt(*account.LastUsedAt)
placeholder else {
		builder.ClearLastUsedAt()
placeholder
	if account.RateLimitedAt != nil {
		builder.SetRateLimitedAt(*account.RateLimitedAt)
placeholder else {
		builder.ClearRateLimitedAt()
placeholder
	if account.RateLimitResetAt != nil {
		builder.SetRateLimitResetAt(*account.RateLimitResetAt)
placeholder else {
		builder.ClearRateLimitResetAt()
placeholder
	if account.OverloadUntil != nil {
		builder.SetOverloadUntil(*account.OverloadUntil)
placeholder else {
		builder.ClearOverloadUntil()
placeholder
	if account.SessionWindowStart != nil {
		builder.SetSessionWindowStart(*account.SessionWindowStart)
placeholder else {
		builder.ClearSessionWindowStart()
placeholder
	if account.SessionWindowEnd != nil {
		builder.SetSessionWindowEnd(*account.SessionWindowEnd)
placeholder else {
		builder.ClearSessionWindowEnd()
placeholder
	if account.SessionWindowStatus != "" {
		builder.SetSessionWindowStatus(account.SessionWindowStatus)
placeholder else {
		builder.ClearSessionWindowStatus()
placeholder

	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrAccountNotFound, nil)
placeholder
	account.UpdatedAt = updated.UpdatedAt
	return nil
placeholder

func (r *accountRepository) Delete(ctx context.Context, id int64) error {
	// 使用事务保证账号与关联分组的删除原子性
	tx, err := r.client.Tx(ctx)
	if err != nil && !errors.Is(err, dbent.ErrTxStarted) {
		return err
placeholder

	var txClient *dbent.Client
	if err == nil {
		defer func() { _ = tx.Rollback() placeholder()
		txClient = tx.Client()
placeholder else {
		// 已处于外部事务中（ErrTxStarted），复用当前 client
		txClient = r.client
placeholder

	if _, err := txClient.AccountGroup.Delete().Where(dbaccountgroup.AccountIDEQ(id)).Exec(ctx); err != nil {
		return err
placeholder
	if _, err := txClient.Account.Delete().Where(dbaccount.IDEQ(id)).Exec(ctx); err != nil {
		return err
placeholder

	if tx != nil {
		return tx.Commit()
placeholder
	return nil
placeholder

func (r *accountRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.Account, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "", "", "")
placeholder

func (r *accountRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, platform, accountType, status, search string) ([]service.Account, *pagination.PaginationResult, error) {
	q := r.client.Account.Query()

	if platform != "" {
		q = q.Where(dbaccount.PlatformEQ(platform))
placeholder
	if accountType != "" {
		q = q.Where(dbaccount.TypeEQ(accountType))
placeholder
	if status != "" {
		q = q.Where(dbaccount.StatusEQ(status))
placeholder
	if search != "" {
		q = q.Where(dbaccount.NameContainsFold(search))
placeholder

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	accounts, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(dbaccount.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	outAccounts, err := r.accountsToService(ctx, accounts)
	if err != nil {
		return nil, nil, err
placeholder
	return outAccounts, paginationResultFromTotal(int64(total), params), nil
placeholder

func (r *accountRepository) ListByGroup(ctx context.Context, groupID int64) ([]service.Account, error) {
	accounts, err := r.queryAccountsByGroup(ctx, groupID, accountGroupQueryOptions{
		status: service.StatusActive,
placeholder)
	if err != nil {
		return nil, err
placeholder
	return accounts, nil
placeholder

func (r *accountRepository) ListActive(ctx context.Context) ([]service.Account, error) {
	accounts, err := r.client.Account.Query().
		Where(dbaccount.StatusEQ(service.StatusActive)).
		Order(dbent.Asc(dbaccount.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.accountsToService(ctx, accounts)
placeholder

func (r *accountRepository) ListByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	accounts, err := r.client.Account.Query().
		Where(
			dbaccount.PlatformEQ(platform),
			dbaccount.StatusEQ(service.StatusActive),
		).
		Order(dbent.Asc(dbaccount.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.accountsToService(ctx, accounts)
placeholder

func (r *accountRepository) UpdateLastUsed(ctx context.Context, id int64) error {
	now := time.Now()
	_, err := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		SetLastUsedAt(now).
		Save(ctx)
	return err
placeholder

func (r *accountRepository) BatchUpdateLastUsed(ctx context.Context, updates map[int64]time.Time) error {
	if len(updates) == 0 {
		return nil
placeholder

	ids := make([]int64, 0, len(updates))
	args := make([]any, 0, len(updates)*2+1)
	caseSQL := "UPDATE accounts SET last_used_at = CASE id"

	idx := 1
	for id, ts := range updates {
		caseSQL += " WHEN $" + itoa(idx) + " THEN $" + itoa(idx+1) + "::timestamptz"
		args = append(args, id, ts)
		ids = append(ids, id)
		idx += 2
placeholder

	caseSQL += " END, updated_at = NOW() WHERE id = ANY($" + itoa(idx) + ") AND deleted_at IS NULL"
	args = append(args, pq.Array(ids))

	_, err := r.sql.ExecContext(ctx, caseSQL, args...)
	return err
placeholder

func (r *accountRepository) SetError(ctx context.Context, id int64, errorMsg string) error {
	_, err := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		SetStatus(service.StatusError).
		SetErrorMessage(errorMsg).
		Save(ctx)
	return err
placeholder

func (r *accountRepository) AddToGroup(ctx context.Context, accountID, groupID int64, priority int) error {
	_, err := r.client.AccountGroup.Create().
		SetAccountID(accountID).
		SetGroupID(groupID).
		SetPriority(priority).
		Save(ctx)
	return err
placeholder

func (r *accountRepository) RemoveFromGroup(ctx context.Context, accountID, groupID int64) error {
	_, err := r.client.AccountGroup.Delete().
		Where(
			dbaccountgroup.AccountIDEQ(accountID),
			dbaccountgroup.GroupIDEQ(groupID),
		).
		Exec(ctx)
	return err
placeholder

func (r *accountRepository) GetGroups(ctx context.Context, accountID int64) ([]service.Group, error) {
	groups, err := r.client.Group.Query().
		Where(
			dbgroup.HasAccountsWith(dbaccount.IDEQ(accountID)),
		).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	outGroups := make([]service.Group, 0, len(groups))
	for i := range groups {
		outGroups = append(outGroups, *groupEntityToService(groups[i]))
placeholder
	return outGroups, nil
placeholder

func (r *accountRepository) BindGroups(ctx context.Context, accountID int64, groupIDs []int64) error {
	// 使用事务保证删除旧绑定与创建新绑定的原子性
	tx, err := r.client.Tx(ctx)
	if err != nil && !errors.Is(err, dbent.ErrTxStarted) {
		return err
placeholder

	var txClient *dbent.Client
	if err == nil {
		defer func() { _ = tx.Rollback() placeholder()
		txClient = tx.Client()
placeholder else {
		// 已处于外部事务中（ErrTxStarted），复用当前 client
		txClient = r.client
placeholder

	if _, err := txClient.AccountGroup.Delete().Where(dbaccountgroup.AccountIDEQ(accountID)).Exec(ctx); err != nil {
		return err
placeholder

	if len(groupIDs) == 0 {
		if tx != nil {
			return tx.Commit()
	placeholder
		return nil
placeholder

	builders := make([]*dbent.AccountGroupCreate, 0, len(groupIDs))
	for i, groupID := range groupIDs {
		builders = append(builders, txClient.AccountGroup.Create().
			SetAccountID(accountID).
			SetGroupID(groupID).
			SetPriority(i+1),
		)
placeholder

	if _, err := txClient.AccountGroup.CreateBulk(builders...).Save(ctx); err != nil {
		return err
placeholder

	if tx != nil {
		return tx.Commit()
placeholder
	return nil
placeholder

func (r *accountRepository) ListSchedulable(ctx context.Context) ([]service.Account, error) {
	now := time.Now()
	accounts, err := r.client.Account.Query().
		Where(
			dbaccount.StatusEQ(service.StatusActive),
			dbaccount.SchedulableEQ(true),
			dbaccount.Or(dbaccount.OverloadUntilIsNil(), dbaccount.OverloadUntilLTE(now)),
			dbaccount.Or(dbaccount.RateLimitResetAtIsNil(), dbaccount.RateLimitResetAtLTE(now)),
		).
		Order(dbent.Asc(dbaccount.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.accountsToService(ctx, accounts)
placeholder

func (r *accountRepository) ListSchedulableByGroupID(ctx context.Context, groupID int64) ([]service.Account, error) {
	return r.queryAccountsByGroup(ctx, groupID, accountGroupQueryOptions{
		status:      service.StatusActive,
		schedulable: true,
placeholder)
placeholder

func (r *accountRepository) ListSchedulableByPlatform(ctx context.Context, platform string) ([]service.Account, error) {
	now := time.Now()
	accounts, err := r.client.Account.Query().
		Where(
			dbaccount.PlatformEQ(platform),
			dbaccount.StatusEQ(service.StatusActive),
			dbaccount.SchedulableEQ(true),
			dbaccount.Or(dbaccount.OverloadUntilIsNil(), dbaccount.OverloadUntilLTE(now)),
			dbaccount.Or(dbaccount.RateLimitResetAtIsNil(), dbaccount.RateLimitResetAtLTE(now)),
		).
		Order(dbent.Asc(dbaccount.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.accountsToService(ctx, accounts)
placeholder

func (r *accountRepository) ListSchedulableByGroupIDAndPlatform(ctx context.Context, groupID int64, platform string) ([]service.Account, error) {
	// 单平台查询复用多平台逻辑，保持过滤条件与排序策略一致。
	return r.queryAccountsByGroup(ctx, groupID, accountGroupQueryOptions{
		status:      service.StatusActive,
		schedulable: true,
		platforms:   []string{platformplaceholder,
placeholder)
placeholder

func (r *accountRepository) ListSchedulableByPlatforms(ctx context.Context, platforms []string) ([]service.Account, error) {
	if len(platforms) == 0 {
		return nil, nil
placeholder
	// 仅返回可调度的活跃账号，并过滤处于过载/限流窗口的账号。
	// 代理与分组信息统一在 accountsToService 中批量加载，避免 N+1 查询。
	now := time.Now()
	accounts, err := r.client.Account.Query().
		Where(
			dbaccount.PlatformIn(platforms...),
			dbaccount.StatusEQ(service.StatusActive),
			dbaccount.SchedulableEQ(true),
			dbaccount.Or(dbaccount.OverloadUntilIsNil(), dbaccount.OverloadUntilLTE(now)),
			dbaccount.Or(dbaccount.RateLimitResetAtIsNil(), dbaccount.RateLimitResetAtLTE(now)),
		).
		Order(dbent.Asc(dbaccount.FieldPriority)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	return r.accountsToService(ctx, accounts)
placeholder

func (r *accountRepository) ListSchedulableByGroupIDAndPlatforms(ctx context.Context, groupID int64, platforms []string) ([]service.Account, error) {
	if len(platforms) == 0 {
		return nil, nil
placeholder
	// 复用按分组查询逻辑，保证分组优先级 + 账号优先级的排序与筛选一致。
	return r.queryAccountsByGroup(ctx, groupID, accountGroupQueryOptions{
		status:      service.StatusActive,
		schedulable: true,
		platforms:   platforms,
placeholder)
placeholder

func (r *accountRepository) SetRateLimited(ctx context.Context, id int64, resetAt time.Time) error {
	now := time.Now()
	_, err := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		SetRateLimitedAt(now).
		SetRateLimitResetAt(resetAt).
		Save(ctx)
	return err
placeholder

func (r *accountRepository) SetOverloaded(ctx context.Context, id int64, until time.Time) error {
	_, err := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		SetOverloadUntil(until).
		Save(ctx)
	return err
placeholder

func (r *accountRepository) ClearRateLimit(ctx context.Context, id int64) error {
	_, err := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		ClearRateLimitedAt().
		ClearRateLimitResetAt().
		ClearOverloadUntil().
		Save(ctx)
	return err
placeholder

func (r *accountRepository) UpdateSessionWindow(ctx context.Context, id int64, start, end *time.Time, status string) error {
	builder := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		SetSessionWindowStatus(status)
	if start != nil {
		builder.SetSessionWindowStart(*start)
placeholder
	if end != nil {
		builder.SetSessionWindowEnd(*end)
placeholder
	_, err := builder.Save(ctx)
	return err
placeholder

func (r *accountRepository) SetSchedulable(ctx context.Context, id int64, schedulable bool) error {
	_, err := r.client.Account.Update().
		Where(dbaccount.IDEQ(id)).
		SetSchedulable(schedulable).
		Save(ctx)
	return err
placeholder

func (r *accountRepository) UpdateExtra(ctx context.Context, id int64, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
placeholder

	// 使用 JSONB 合并操作实现原子更新，避免读-改-写的并发丢失更新问题
	payload, err := json.Marshal(updates)
	if err != nil {
		return err
placeholder

	client := clientFromContext(ctx, r.client)
	result, err := client.ExecContext(
		ctx,
		"UPDATE accounts SET extra = COALESCE(extra, '{placeholder'::jsonb) || $1::jsonb, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL",
		payload, id,
	)
	if err != nil {
		return err
placeholder

	affected, err := result.RowsAffected()
	if err != nil {
		return err
placeholder
	if affected == 0 {
		return service.ErrAccountNotFound
placeholder
	return nil
placeholder

func (r *accountRepository) BulkUpdate(ctx context.Context, ids []int64, updates service.AccountBulkUpdate) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
placeholder

	setClauses := make([]string, 0, 8)
	args := make([]any, 0, 8)

	idx := 1
	if updates.Name != nil {
		setClauses = append(setClauses, "name = $"+itoa(idx))
		args = append(args, *updates.Name)
		idx++
placeholder
	if updates.ProxyID != nil {
		setClauses = append(setClauses, "proxy_id = $"+itoa(idx))
		args = append(args, *updates.ProxyID)
		idx++
placeholder
	if updates.Concurrency != nil {
		setClauses = append(setClauses, "concurrency = $"+itoa(idx))
		args = append(args, *updates.Concurrency)
		idx++
placeholder
	if updates.Priority != nil {
		setClauses = append(setClauses, "priority = $"+itoa(idx))
		args = append(args, *updates.Priority)
		idx++
placeholder
	if updates.Status != nil {
		setClauses = append(setClauses, "status = $"+itoa(idx))
		args = append(args, *updates.Status)
		idx++
placeholder
	// JSONB 需要合并而非覆盖，使用 raw SQL 保持旧行为。
	if len(updates.Credentials) > 0 {
		payload, err := json.Marshal(updates.Credentials)
		if err != nil {
			return 0, err
	placeholder
		setClauses = append(setClauses, "credentials = COALESCE(credentials, '{placeholder'::jsonb) || $"+itoa(idx)+"::jsonb")
		args = append(args, payload)
		idx++
placeholder
	if len(updates.Extra) > 0 {
		payload, err := json.Marshal(updates.Extra)
		if err != nil {
			return 0, err
	placeholder
		setClauses = append(setClauses, "extra = COALESCE(extra, '{placeholder'::jsonb) || $"+itoa(idx)+"::jsonb")
		args = append(args, payload)
		idx++
placeholder

	if len(setClauses) == 0 {
		return 0, nil
placeholder

	setClauses = append(setClauses, "updated_at = NOW()")

	query := "UPDATE accounts SET " + joinClauses(setClauses, ", ") + " WHERE id = ANY($" + itoa(idx) + ") AND deleted_at IS NULL"
	args = append(args, pq.Array(ids))

	result, err := r.sql.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
placeholder
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
placeholder
	return rows, nil
placeholder

type accountGroupQueryOptions struct {
	status      string
	schedulable bool
	platforms   []string // 允许的多个平台，空切片表示不进行平台过滤
placeholder

func (r *accountRepository) queryAccountsByGroup(ctx context.Context, groupID int64, opts accountGroupQueryOptions) ([]service.Account, error) {
	q := r.client.AccountGroup.Query().
		Where(dbaccountgroup.GroupIDEQ(groupID))

	// 通过 account_groups 中间表查询账号，并按需叠加状态/平台/调度能力过滤。
	preds := make([]dbpredicate.Account, 0, 6)
	preds = append(preds, dbaccount.DeletedAtIsNil())
	if opts.status != "" {
		preds = append(preds, dbaccount.StatusEQ(opts.status))
placeholder
	if len(opts.platforms) > 0 {
		preds = append(preds, dbaccount.PlatformIn(opts.platforms...))
placeholder
	if opts.schedulable {
		now := time.Now()
		preds = append(preds,
			dbaccount.SchedulableEQ(true),
			dbaccount.Or(dbaccount.OverloadUntilIsNil(), dbaccount.OverloadUntilLTE(now)),
			dbaccount.Or(dbaccount.RateLimitResetAtIsNil(), dbaccount.RateLimitResetAtLTE(now)),
		)
placeholder

	if len(preds) > 0 {
		q = q.Where(dbaccountgroup.HasAccountWith(preds...))
placeholder

	groups, err := q.
		Order(
			dbaccountgroup.ByPriority(),
			dbaccountgroup.ByAccountField(dbaccount.FieldPriority),
		).
		WithAccount().
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	orderedIDs := make([]int64, 0, len(groups))
	accountMap := make(map[int64]*dbent.Account, len(groups))
	for _, ag := range groups {
		if ag.Edges.Account == nil {
			continue
	placeholder
		if _, exists := accountMap[ag.AccountID]; exists {
			continue
	placeholder
		accountMap[ag.AccountID] = ag.Edges.Account
		orderedIDs = append(orderedIDs, ag.AccountID)
placeholder

	accounts := make([]*dbent.Account, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		if acc, ok := accountMap[id]; ok {
			accounts = append(accounts, acc)
	placeholder
placeholder

	return r.accountsToService(ctx, accounts)
placeholder

func (r *accountRepository) accountsToService(ctx context.Context, accounts []*dbent.Account) ([]service.Account, error) {
	if len(accounts) == 0 {
		return []service.Account{placeholder, nil
placeholder

	accountIDs := make([]int64, 0, len(accounts))
	proxyIDs := make([]int64, 0, len(accounts))
	for _, acc := range accounts {
		accountIDs = append(accountIDs, acc.ID)
		if acc.ProxyID != nil {
			proxyIDs = append(proxyIDs, *acc.ProxyID)
	placeholder
placeholder

	proxyMap, err := r.loadProxies(ctx, proxyIDs)
	if err != nil {
		return nil, err
placeholder
	groupsByAccount, groupIDsByAccount, accountGroupsByAccount, err := r.loadAccountGroups(ctx, accountIDs)
	if err != nil {
		return nil, err
placeholder

	outAccounts := make([]service.Account, 0, len(accounts))
	for _, acc := range accounts {
		out := accountEntityToService(acc)
		if out == nil {
			continue
	placeholder
		if acc.ProxyID != nil {
			if proxy, ok := proxyMap[*acc.ProxyID]; ok {
				out.Proxy = proxy
		placeholder
	placeholder
		if groups, ok := groupsByAccount[acc.ID]; ok {
			out.Groups = groups
	placeholder
		if groupIDs, ok := groupIDsByAccount[acc.ID]; ok {
			out.GroupIDs = groupIDs
	placeholder
		if ags, ok := accountGroupsByAccount[acc.ID]; ok {
			out.AccountGroups = ags
	placeholder
		outAccounts = append(outAccounts, *out)
placeholder

	return outAccounts, nil
placeholder

func (r *accountRepository) loadProxies(ctx context.Context, proxyIDs []int64) (map[int64]*service.Proxy, error) {
	proxyMap := make(map[int64]*service.Proxy)
	if len(proxyIDs) == 0 {
		return proxyMap, nil
placeholder

	proxies, err := r.client.Proxy.Query().Where(dbproxy.IDIn(proxyIDs...)).All(ctx)
	if err != nil {
		return nil, err
placeholder

	for _, p := range proxies {
		proxyMap[p.ID] = proxyEntityToService(p)
placeholder
	return proxyMap, nil
placeholder

func (r *accountRepository) loadAccountGroups(ctx context.Context, accountIDs []int64) (map[int64][]*service.Group, map[int64][]int64, map[int64][]service.AccountGroup, error) {
	groupsByAccount := make(map[int64][]*service.Group)
	groupIDsByAccount := make(map[int64][]int64)
	accountGroupsByAccount := make(map[int64][]service.AccountGroup)

	if len(accountIDs) == 0 {
		return groupsByAccount, groupIDsByAccount, accountGroupsByAccount, nil
placeholder

	entries, err := r.client.AccountGroup.Query().
		Where(dbaccountgroup.AccountIDIn(accountIDs...)).
		WithGroup().
		Order(dbaccountgroup.ByAccountID(), dbaccountgroup.ByPriority()).
		All(ctx)
	if err != nil {
		return nil, nil, nil, err
placeholder

	for _, ag := range entries {
		groupSvc := groupEntityToService(ag.Edges.Group)
		agSvc := service.AccountGroup{
			AccountID: ag.AccountID,
			GroupID:   ag.GroupID,
			Priority:  ag.Priority,
			CreatedAt: ag.CreatedAt,
			Group:     groupSvc,
	placeholder
		accountGroupsByAccount[ag.AccountID] = append(accountGroupsByAccount[ag.AccountID], agSvc)
		groupIDsByAccount[ag.AccountID] = append(groupIDsByAccount[ag.AccountID], ag.GroupID)
		if groupSvc != nil {
			groupsByAccount[ag.AccountID] = append(groupsByAccount[ag.AccountID], groupSvc)
	placeholder
placeholder

	return groupsByAccount, groupIDsByAccount, accountGroupsByAccount, nil
placeholder

func accountEntityToService(m *dbent.Account) *service.Account {
	if m == nil {
		return nil
placeholder

	return &service.Account{
		ID:                  m.ID,
		Name:                m.Name,
		Platform:            m.Platform,
		Type:                m.Type,
		Credentials:         copyJSONMap(m.Credentials),
		Extra:               copyJSONMap(m.Extra),
		ProxyID:             m.ProxyID,
		Concurrency:         m.Concurrency,
		Priority:            m.Priority,
		Status:              m.Status,
		ErrorMessage:        derefString(m.ErrorMessage),
		LastUsedAt:          m.LastUsedAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		Schedulable:         m.Schedulable,
		RateLimitedAt:       m.RateLimitedAt,
		RateLimitResetAt:    m.RateLimitResetAt,
		OverloadUntil:       m.OverloadUntil,
		SessionWindowStart:  m.SessionWindowStart,
		SessionWindowEnd:    m.SessionWindowEnd,
		SessionWindowStatus: derefString(m.SessionWindowStatus),
placeholder
placeholder

func normalizeJSONMap(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{placeholder
placeholder
	return in
placeholder

func copyJSONMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
placeholder
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
placeholder
	return out
placeholder

func joinClauses(clauses []string, sep string) string {
	if len(clauses) == 0 {
		return ""
placeholder
	out := clauses[0]
	for i := 1; i < len(clauses); i++ {
		out += sep + clauses[i]
placeholder
	return out
placeholder

func itoa(v int) string {
	return strconv.Itoa(v)
placeholder
