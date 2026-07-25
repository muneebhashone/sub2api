package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/apikey"
	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/authidentitychannel"
	dbgroup "github.com/Wei-Shaw/sub2api/ent/group"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	"github.com/Wei-Shaw/sub2api/ent/predicate"
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"
	dbuser "github.com/Wei-Shaw/sub2api/ent/user"
	"github.com/Wei-Shaw/sub2api/ent/userallowedgroup"
	"github.com/Wei-Shaw/sub2api/ent/usersubscription"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"

	entsql "entgo.io/ent/dialect/sql"
)

type userRepository struct {
	client *dbent.Client
	sql    sqlExecutor
placeholder

var _ service.RedeemUserAdjustmentRepository = (*userRepository)(nil)

func NewUserRepository(client *dbent.Client, sqlDB *sql.DB) service.UserRepository {
	return newUserRepositoryWithSQL(client, sqlDB)
placeholder

func newUserRepositoryWithSQL(client *dbent.Client, sqlq sqlExecutor) *userRepository {
	return &userRepository{client: client, sql: sqlqplaceholder
placeholder

func (r *userRepository) Create(ctx context.Context, userIn *service.User) error {
	return r.create(ctx, userIn, false)
placeholder

// CreateWithEmailAliasGuard 见 service.UserRepository：在邮箱唯一性锁内复查收件箱身份，
// 供注册路径使用。
func (r *userRepository) CreateWithEmailAliasGuard(ctx context.Context, userIn *service.User) error {
	return r.create(ctx, userIn, true)
placeholder

func (r *userRepository) create(ctx context.Context, userIn *service.User, guardEmailAlias bool) error {
	if userIn == nil {
		return nil
placeholder

	// 统一使用 ent 的事务：保证用户与允许分组的更新原子化，
	// 并避免基于 *sql.Tx 手动构造 ent client 导致的 ExecQuerier 断言错误。
	tx, err := r.client.Tx(ctx)
	if err != nil && !errors.Is(err, dbent.ErrTxStarted) {
		return err
placeholder

	var txClient *dbent.Client
	txCtx := ctx
	if err == nil {
		defer func() { _ = tx.Rollback() placeholder()
		txClient = tx.Client()
		txCtx = dbent.NewTxContext(ctx, tx)
placeholder else {
		// 已处于外部事务中（ErrTxStarted），复用当前事务 client 并由调用方负责提交/回滚。
		if existingTx := dbent.TxFromContext(ctx); existingTx != nil {
			txClient = existingTx.Client()
	placeholder else {
			txClient = r.client
	placeholder
placeholder

	lockKeys := []string{normalizedEmailUniquenessLockKey(userIn.Email)placeholder
	if guardEmailAlias {
		// 别名变体的字面量不同，唯一索引无法兜底；用收件箱身份锁把同一收件箱的并发注册串行化。
		lockKeys = append(lockKeys, emailAliasUniquenessLockKey(userIn.Email))
placeholder
	releaseEmailLock, err := lockRepositoryScopedKeys(
		txCtx,
		txClient,
		txAwareSQLExecutor(txCtx, r.sql, r.client),
		lockKeys...,
	)
	if err != nil {
		return err
placeholder
	defer releaseEmailLock()

	if err := ensureNormalizedEmailAvailableWithClient(txCtx, txClient, 0, userIn.Email); err != nil {
		return err
placeholder

	if guardEmailAlias {
		aliasExists, err := existsByEmailAliasWithClient(txCtx, txClient, userIn.Email)
		if err != nil {
			return err
	placeholder
		if aliasExists {
			return service.ErrEmailExists
	placeholder
placeholder

	created, err := txClient.User.Create().
		SetEmail(userIn.Email).
		SetUsername(userIn.Username).
		SetNotes(userIn.Notes).
		SetPasswordHash(userIn.PasswordHash).
		SetRole(userIn.Role).
		SetBalance(userIn.Balance).
		SetConcurrency(userIn.Concurrency).
		SetStatus(userIn.Status).
		SetSignupSource(userSignupSourceOrDefault(userIn.SignupSource)).
		SetNillableLastLoginAt(userIn.LastLoginAt).
		SetNillableLastActiveAt(userIn.LastActiveAt).
		SetRpmLimit(userIn.RPMLimit).
		Save(txCtx)
	if err != nil {
		return translatePersistenceError(err, nil, service.ErrEmailExists)
placeholder

	if err := r.syncUserAllowedGroupsWithClient(txCtx, txClient, created.ID, userIn.AllowedGroups); err != nil {
		return err
placeholder
	if err := ensureEmailAuthIdentityWithClient(txCtx, txClient, created.ID, created.Email, "user_repo_create"); err != nil {
		return err
placeholder

	if tx != nil {
		if err := tx.Commit(); err != nil {
			return err
	placeholder
placeholder

	applyUserEntityToService(userIn, created)
	return nil
placeholder

func (r *userRepository) GetByID(ctx context.Context, id int64) (*service.User, error) {
	m, err := r.client.User.Query().Where(dbuser.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder

	out := userEntityToService(m)
	groups, err := r.loadAllowedGroups(ctx, []int64{idplaceholder)
	if err != nil {
		return nil, err
placeholder
	if v, ok := groups[id]; ok {
		out.AllowedGroups = v
placeholder
	return out, nil
placeholder

func (r *userRepository) GetByIDIncludeDeleted(ctx context.Context, id int64) (*service.User, error) {
	ctx = mixins.SkipSoftDelete(ctx)
	m, err := r.client.User.Query().Where(dbuser.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	out := userEntityToService(m)
	groups, err := r.loadAllowedGroups(ctx, []int64{idplaceholder)
	if err != nil {
		return nil, err
placeholder
	if v, ok := groups[id]; ok {
		out.AllowedGroups = v
placeholder
	return out, nil
placeholder

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	matches, err := r.client.User.Query().
		Where(userEmailLookupPredicate(email)).
		Order(dbent.Asc(dbuser.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder
	if len(matches) == 0 {
		return nil, service.ErrUserNotFound
placeholder
	if len(matches) > 1 {
		return nil, fmt.Errorf("normalized email lookup matched multiple users for %q", strings.TrimSpace(email))
placeholder
	m := matches[0]

	out := userEntityToService(m)
	groups, err := r.loadAllowedGroups(ctx, []int64{m.IDplaceholder)
	if err != nil {
		return nil, err
placeholder
	if v, ok := groups[m.ID]; ok {
		out.AllowedGroups = v
placeholder
	return out, nil
placeholder

func (r *userRepository) Update(ctx context.Context, userIn *service.User) error {
	if userIn == nil {
		return nil
placeholder

	// 使用 ent 事务包裹用户更新与 allowed_groups 同步，避免跨层事务不一致。
	tx, err := r.client.Tx(ctx)
	if err != nil && !errors.Is(err, dbent.ErrTxStarted) {
		return err
placeholder

	var txClient *dbent.Client
	txCtx := ctx
	if err == nil {
		defer func() { _ = tx.Rollback() placeholder()
		txClient = tx.Client()
		txCtx = dbent.NewTxContext(ctx, tx)
placeholder else {
		// 已处于外部事务中（ErrTxStarted），复用当前事务 client 并由调用方负责提交/回滚。
		if existingTx := dbent.TxFromContext(ctx); existingTx != nil {
			txClient = existingTx.Client()
	placeholder else {
			txClient = r.client
	placeholder
placeholder

	releaseEmailLock, err := lockRepositoryScopedKeys(
		txCtx,
		txClient,
		txAwareSQLExecutor(txCtx, r.sql, r.client),
		normalizedEmailUniquenessLockKey(userIn.Email),
	)
	if err != nil {
		return err
placeholder
	defer releaseEmailLock()

	if err := ensureNormalizedEmailAvailableWithClient(txCtx, txClient, userIn.ID, userIn.Email); err != nil {
		return err
placeholder

	existing, err := clientFromContext(txCtx, txClient).User.Get(txCtx, userIn.ID)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	oldEmail := existing.Email

	updateOp := txClient.User.UpdateOneID(userIn.ID).
		SetEmail(userIn.Email).
		SetUsername(userIn.Username).
		SetNotes(userIn.Notes).
		SetPasswordHash(userIn.PasswordHash).
		SetRole(userIn.Role).
		SetBalance(userIn.Balance).
		SetConcurrency(userIn.Concurrency).
		SetStatus(userIn.Status).
		SetBalanceNotifyEnabled(userIn.BalanceNotifyEnabled).
		SetBalanceNotifyThresholdType(userIn.BalanceNotifyThresholdType).
		SetNillableBalanceNotifyThreshold(userIn.BalanceNotifyThreshold).
		SetBalanceNotifyExtraEmails(marshalExtraEmails(userIn.BalanceNotifyExtraEmails)).
		SetTotalRecharged(userIn.TotalRecharged).
		SetRpmLimit(userIn.RPMLimit)
	if userIn.SignupSource != "" {
		updateOp = updateOp.SetSignupSource(userIn.SignupSource)
placeholder
	if userIn.LastLoginAt != nil {
		updateOp = updateOp.SetLastLoginAt(*userIn.LastLoginAt)
placeholder
	if userIn.LastActiveAt != nil {
		updateOp = updateOp.SetLastActiveAt(*userIn.LastActiveAt)
placeholder
	if userIn.BalanceNotifyThreshold == nil {
		updateOp = updateOp.ClearBalanceNotifyThreshold()
placeholder
	updated, err := updateOp.Save(txCtx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, service.ErrEmailExists)
placeholder

	if err := r.syncUserAllowedGroupsWithClient(txCtx, txClient, updated.ID, userIn.AllowedGroups); err != nil {
		return err
placeholder
	if err := replaceEmailAuthIdentityWithClient(txCtx, txClient, updated.ID, oldEmail, updated.Email, "user_repo_update"); err != nil {
		return err
placeholder

	if tx != nil {
		if err := tx.Commit(); err != nil {
			return err
	placeholder
placeholder

	userIn.UpdatedAt = updated.UpdatedAt
	return nil
placeholder

func ensureEmailAuthIdentityWithClient(ctx context.Context, client *dbent.Client, userID int64, email string, source string) error {
	client = clientFromContext(ctx, client)
	if client == nil || userID <= 0 {
		return nil
placeholder

	subject := normalizeEmailAuthIdentitySubject(email)
	if subject == "" {
		return nil
placeholder

	if err := client.AuthIdentity.Create().
		SetUserID(userID).
		SetProviderType("email").
		SetProviderKey("email").
		SetProviderSubject(subject).
		SetVerifiedAt(time.Now().UTC()).
		SetMetadata(map[string]any{"source": sourceplaceholder).
		OnConflictColumns(
			authidentity.FieldProviderType,
			authidentity.FieldProviderKey,
			authidentity.FieldProviderSubject,
		).
		DoNothing().
		Exec(ctx); err != nil {
		if !isSQLNoRowsError(err) {
			return err
	placeholder
placeholder

	identity, err := client.AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ(subject),
		).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil
	placeholder
		return err
placeholder
	if identity.UserID != userID {
		return ErrAuthIdentityOwnershipConflict
placeholder
	return nil
placeholder

func replaceEmailAuthIdentityWithClient(ctx context.Context, client *dbent.Client, userID int64, oldEmail, newEmail string, source string) error {
	newSubject := normalizeEmailAuthIdentitySubject(newEmail)
	if err := ensureEmailAuthIdentityWithClient(ctx, client, userID, newEmail, source); err != nil {
		return err
placeholder

	oldSubject := normalizeEmailAuthIdentitySubject(oldEmail)
	if oldSubject == "" || oldSubject == newSubject {
		return nil
placeholder

	_, err := clientFromContext(ctx, client).AuthIdentity.Delete().
		Where(
			authidentity.UserIDEQ(userID),
			authidentity.ProviderTypeEQ("email"),
			authidentity.ProviderKeyEQ("email"),
			authidentity.ProviderSubjectEQ(oldSubject),
		).
		Exec(ctx)
	return err
placeholder

func normalizeEmailAuthIdentitySubject(email string) string {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return ""
placeholder
	if strings.HasSuffix(normalized, service.LinuxDoConnectSyntheticEmailDomain) ||
		strings.HasSuffix(normalized, service.OIDCConnectSyntheticEmailDomain) ||
		strings.HasSuffix(normalized, service.WeChatConnectSyntheticEmailDomain) ||
		strings.HasSuffix(normalized, service.DingTalkConnectSyntheticEmailDomain) {
		return ""
placeholder
	return normalized
placeholder

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	// 复用 context 中已存在的事务（如 AdminService.DeleteUser 把删 Key 与删 User 包在同一事务中），
	// 由调用方负责提交/回滚，保证两者的原子性。
	if existingTx := dbent.TxFromContext(ctx); existingTx != nil {
		return r.deleteUser(ctx, existingTx.Client(), id)
placeholder

	tx, err := r.client.Tx(ctx)
	if err != nil && !errors.Is(err, dbent.ErrTxStarted) {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	exec := r.client
	if err == nil {
		defer func() { _ = tx.Rollback() placeholder()
		exec = tx.Client()
placeholder
	// err == dbent.ErrTxStarted 时复用当前事务（exec = r.client）。

	if err := r.deleteUser(ctx, exec, id); err != nil {
		return err
placeholder

	if tx != nil {
		if err := tx.Commit(); err != nil {
			return translatePersistenceError(err, service.ErrUserNotFound, nil)
	placeholder
placeholder
	return nil
placeholder

// deleteUser 在给定 client（可能是外部事务 client）上删除用户及其身份关联记录，自身不开启/提交事务。
func (r *userRepository) deleteUser(ctx context.Context, exec *dbent.Client, id int64) error {
	identityIDs, err := exec.AuthIdentity.Query().
		Where(authidentity.UserIDEQ(id)).
		IDs(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	if len(identityIDs) > 0 {
		if _, err := exec.IdentityAdoptionDecision.Update().
			Where(identityadoptiondecision.IdentityIDIn(identityIDs...)).
			ClearIdentityID().
			Save(ctx); err != nil {
			return translatePersistenceError(err, service.ErrUserNotFound, nil)
	placeholder
		if _, err := exec.AuthIdentityChannel.Delete().
			Where(authidentitychannel.IdentityIDIn(identityIDs...)).
			Exec(ctx); err != nil {
			return translatePersistenceError(err, service.ErrUserNotFound, nil)
	placeholder
		if _, err := exec.AuthIdentity.Delete().
			Where(authidentity.UserIDEQ(id)).
			Exec(ctx); err != nil {
			return translatePersistenceError(err, service.ErrUserNotFound, nil)
	placeholder
placeholder

	affected, err := exec.User.Delete().Where(dbuser.IDEQ(id)).Exec(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	if affected == 0 {
		return service.ErrUserNotFound
placeholder
	return nil
placeholder

func (r *userRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.User, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, service.UserListFilters{placeholder)
placeholder

func (r *userRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, filters service.UserListFilters) ([]service.User, *pagination.PaginationResult, error) {
	// SkipSoftDelete 仅作用于 User 身份解析（下方 Count/All）；订阅、分组等关联实体沿用原始 ctx，避免穿透到这些同样带软删除的实体而带出已删除行。
	userCtx := ctx
	if filters.IncludeDeleted {
		userCtx = mixins.SkipSoftDelete(ctx)
placeholder

	q := r.client.User.Query()

	if filters.Status != "" {
		q = q.Where(dbuser.StatusEQ(filters.Status))
placeholder
	if filters.Role != "" {
		q = q.Where(dbuser.RoleEQ(filters.Role))
placeholder
	if filters.Search != "" {
		q = q.Where(
			dbuser.Or(
				dbuser.EmailContainsFold(filters.Search),
				dbuser.UsernameContainsFold(filters.Search),
				dbuser.NotesContainsFold(filters.Search),
				dbuser.HasAPIKeysWith(apikey.KeyContainsFold(filters.Search)),
			),
		)
placeholder

	if filters.GroupName != "" {
		q = q.Where(dbuser.HasAllowedGroupsWith(
			dbgroup.NameContainsFold(filters.GroupName),
		))
placeholder

	if filters.APIKeyGroupID > 0 {
		// 按"API Key 实际绑定的分组"过滤：用户只要有任意一个未软删除的 API Key
		// 绑定到该分组即命中（EXISTS 语义）。
		// 注意：SoftDeleteMixin 的拦截器不会自动下沉到 HasAPIKeysWith 子查询，
		// 必须显式加 apikey.DeletedAtIsNil()，否则已软删除的 key 会污染过滤结果。
		q = q.Where(dbuser.HasAPIKeysWith(
			apikey.GroupIDEQ(filters.APIKeyGroupID),
			apikey.DeletedAtIsNil(),
		))
placeholder

	// If attribute filters are specified, we need to filter by user IDs first
	var allowedUserIDs []int64
	if len(filters.Attributes) > 0 {
		var attrErr error
		allowedUserIDs, attrErr = r.filterUsersByAttributes(ctx, filters.Attributes)
		if attrErr != nil {
			return nil, nil, attrErr
	placeholder
		if len(allowedUserIDs) == 0 {
			// No users match the attribute filters
			return []service.User{placeholder, paginationResultFromTotal(0, params), nil
	placeholder
		q = q.Where(dbuser.IDIn(allowedUserIDs...))
placeholder

	total, err := q.Clone().Count(userCtx)
	if err != nil {
		return nil, nil, err
placeholder

	usersQuery := q.
		Offset(params.Offset()).
		Limit(params.Limit())
	for _, order := range userListOrder(params) {
		usersQuery = usersQuery.Order(order)
placeholder

	users, err := usersQuery.All(userCtx)
	if err != nil {
		return nil, nil, err
placeholder

	outUsers := make([]service.User, 0, len(users))
	if len(users) == 0 {
		return outUsers, paginationResultFromTotal(int64(total), params), nil
placeholder

	userIDs := make([]int64, 0, len(users))
	userMap := make(map[int64]*service.User, len(users))
	for i := range users {
		userIDs = append(userIDs, users[i].ID)
		u := userEntityToService(users[i])
		outUsers = append(outUsers, *u)
		userMap[u.ID] = &outUsers[len(outUsers)-1]
placeholder

	shouldLoadSubscriptions := filters.IncludeSubscriptions == nil || *filters.IncludeSubscriptions
	if shouldLoadSubscriptions {
		// Batch load active subscriptions with groups to avoid N+1.
		subs, err := r.client.UserSubscription.Query().
			Where(
				usersubscription.UserIDIn(userIDs...),
				usersubscription.StatusEQ(service.SubscriptionStatusActive),
			).
			WithGroup().
			All(ctx)
		if err != nil {
			return nil, nil, err
	placeholder

		for i := range subs {
			if u, ok := userMap[subs[i].UserID]; ok {
				u.Subscriptions = append(u.Subscriptions, *userSubscriptionEntityToService(subs[i]))
		placeholder
	placeholder
placeholder

	allowedGroupsByUser, err := r.loadAllowedGroups(ctx, userIDs)
	if err != nil {
		return nil, nil, err
placeholder
	for id, u := range userMap {
		if groups, ok := allowedGroupsByUser[id]; ok {
			u.AllowedGroups = groups
	placeholder
placeholder

	return outUsers, paginationResultFromTotal(int64(total), params), nil
placeholder

func userListOrder(params pagination.PaginationParams) []func(*entsql.Selector) {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := params.NormalizedSortOrder(pagination.SortOrderDesc)

	if sortBy == "last_used_at" {
		return userLastUsedAtOrder(sortOrder)
placeholder

	var field string
	defaultField := true
	nullsLastField := false
	switch sortBy {
	case "email":
		field = dbuser.FieldEmail
		defaultField = false
	case "username":
		field = dbuser.FieldUsername
		defaultField = false
	case "role":
		field = dbuser.FieldRole
		defaultField = false
	case "balance":
		field = dbuser.FieldBalance
		defaultField = false
	case "concurrency":
		field = dbuser.FieldConcurrency
		defaultField = false
	case "status":
		field = dbuser.FieldStatus
		defaultField = false
	case "created_at":
		field = dbuser.FieldCreatedAt
		defaultField = false
	case "last_active_at":
		field = dbuser.FieldLastActiveAt
		defaultField = false
		nullsLastField = true
	default:
		field = dbuser.FieldID
placeholder

	if sortOrder == pagination.SortOrderAsc {
		if defaultField && field == dbuser.FieldID {
			return []func(*entsql.Selector){dbent.Asc(dbuser.FieldID)placeholder
	placeholder
		if nullsLastField {
			return []func(*entsql.Selector){
				entsql.OrderByField(field, entsql.OrderNullsLast()).ToFunc(),
				dbent.Asc(dbuser.FieldID),
		placeholder
	placeholder
		return []func(*entsql.Selector){dbent.Asc(field), dbent.Asc(dbuser.FieldID)placeholder
placeholder
	if defaultField && field == dbuser.FieldID {
		return []func(*entsql.Selector){dbent.Desc(dbuser.FieldID)placeholder
placeholder
	if nullsLastField {
		return []func(*entsql.Selector){
			entsql.OrderByField(field, entsql.OrderDesc(), entsql.OrderNullsLast()).ToFunc(),
			dbent.Desc(dbuser.FieldID),
	placeholder
placeholder
	return []func(*entsql.Selector){dbent.Desc(field), dbent.Desc(dbuser.FieldID)placeholder
placeholder

func (r *userRepository) GetLatestUsedAtByUserIDs(ctx context.Context, userIDs []int64) (map[int64]*time.Time, error) {
	result := make(map[int64]*time.Time, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
placeholder
	if r.sql == nil {
		return nil, fmt.Errorf("sql executor is not configured")
placeholder

	const query = `
		SELECT user_id, MAX(created_at) AS last_used_at
		FROM usage_logs
		WHERE user_id = ANY($1)
		GROUP BY user_id
	`

	rows, err := r.sql.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	for rows.Next() {
		var (
			userID     int64
			lastUsedAt time.Time
		)
		if scanErr := rows.Scan(&userID, &lastUsedAt); scanErr != nil {
			return nil, scanErr
	placeholder
		ts := lastUsedAt.UTC()
		result[userID] = &ts
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

func (r *userRepository) GetLatestUsedAtByUserID(ctx context.Context, userID int64) (*time.Time, error) {
	latestByUserID, err := r.GetLatestUsedAtByUserIDs(ctx, []int64{userIDplaceholder)
	if err != nil {
		return nil, err
placeholder
	return latestByUserID[userID], nil
placeholder

func userLastUsedAtOrder(sortOrder string) []func(*entsql.Selector) {
	orderExpr := func(direction, nulls string, tieOrder func(string) string) func(*entsql.Selector) {
		return func(s *entsql.Selector) {
			subquery := fmt.Sprintf("(SELECT MAX(created_at) FROM usage_logs WHERE user_id = %s)", s.C(dbuser.FieldID))
			s.OrderExpr(entsql.Expr(subquery + " " + direction + " NULLS " + nulls))
			s.OrderBy(tieOrder(s.C(dbuser.FieldID)))
	placeholder
placeholder

	if sortOrder == pagination.SortOrderAsc {
		return []func(*entsql.Selector){
			orderExpr("ASC", "FIRST", entsql.Asc),
	placeholder
placeholder
	return []func(*entsql.Selector){
		orderExpr("DESC", "LAST", entsql.Desc),
placeholder
placeholder

// filterUsersByAttributes returns user IDs that match ALL the given attribute filters
func (r *userRepository) filterUsersByAttributes(ctx context.Context, attrs map[int64]string) ([]int64, error) {
	if len(attrs) == 0 {
		return nil, nil
placeholder

	if r.sql == nil {
		return nil, fmt.Errorf("sql executor is not configured")
placeholder

	clauses := make([]string, 0, len(attrs))
	args := make([]any, 0, len(attrs)*2+1)
	argIndex := 1
	for attrID, value := range attrs {
		clauses = append(clauses, fmt.Sprintf("(attribute_id = $%d AND value ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, attrID, "%"+value+"%")
		argIndex += 2
placeholder

	query := fmt.Sprintf(
		`SELECT user_id
		 FROM user_attribute_values
		 WHERE %s
		 GROUP BY user_id
		 HAVING COUNT(DISTINCT attribute_id) = $%d`,
		strings.Join(clauses, " OR "),
		argIndex,
	)
	args = append(args, len(attrs))

	rows, err := r.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	result := make([]int64, 0)
	for rows.Next() {
		var userID int64
		if scanErr := rows.Scan(&userID); scanErr != nil {
			return nil, scanErr
	placeholder
		result = append(result, userID)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

func (r *userRepository) UpdateBalance(ctx context.Context, id int64, amount float64) error {
	client := clientFromContext(ctx, r.client)
	update := client.User.Update().Where(dbuser.IDEQ(id)).AddBalance(amount)
	// Track cumulative recharge amount for percentage-based notifications
	if amount > 0 {
		update = update.AddTotalRecharged(amount)
placeholder
	n, err := update.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	if n == 0 {
		return service.ErrUserNotFound
placeholder
	return nil
placeholder

func (r *userRepository) ApplyRedeemBalanceAdjustment(ctx context.Context, id int64, delta float64) error {
	const updateSQL = `
		UPDATE users
		SET balance = GREATEST(balance + $1, 0), updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	client := clientFromContext(ctx, r.client)
	result, err := client.ExecContext(ctx, updateSQL, delta, id)
	if err != nil {
		return err
placeholder
	affected, err := result.RowsAffected()
	if err != nil {
		return err
placeholder
	if affected == 0 {
		return service.ErrUserNotFound
placeholder
	return nil
placeholder

// DeductBalance 扣除用户余额
// 透支策略：允许余额变为负数，确保当前请求能够完成
// 中间件会阻止余额 <= 0 的用户发起后续请求
func (r *userRepository) DeductBalance(ctx context.Context, id int64, amount float64) error {
	client := clientFromContext(ctx, r.client)
	n, err := client.User.Update().
		Where(dbuser.IDEQ(id), dbuser.BalanceGTE(amount)).
		AddBalance(-amount).
		Save(ctx)
	if err != nil {
		return err
placeholder
	if n > 0 {
		return nil
placeholder

	n, err = client.User.Update().
		Where(dbuser.IDEQ(id)).
		AddBalance(-amount).
		Save(ctx)
	if err != nil {
		return err
placeholder
	if n == 0 {
		return service.ErrUserNotFound
placeholder
	return nil
placeholder

func (r *userRepository) UpdateConcurrency(ctx context.Context, id int64, amount int) error {
	client := clientFromContext(ctx, r.client)
	n, err := client.User.Update().Where(dbuser.IDEQ(id)).AddConcurrency(amount).Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	if n == 0 {
		return service.ErrUserNotFound
placeholder
	return nil
placeholder

func (r *userRepository) ApplyRedeemConcurrencyAdjustment(ctx context.Context, id int64, delta int) error {
	const updateSQL = `
		UPDATE users
		SET concurrency = GREATEST(concurrency + $1, 0), updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`
	client := clientFromContext(ctx, r.client)
	result, err := client.ExecContext(ctx, updateSQL, delta, id)
	if err != nil {
		return err
placeholder
	affected, err := result.RowsAffected()
	if err != nil {
		return err
placeholder
	if affected == 0 {
		return service.ErrUserNotFound
placeholder
	return nil
placeholder

func (r *userRepository) BatchSetConcurrency(ctx context.Context, userIDs []int64, value int) (int, error) {
	if len(userIDs) == 0 {
		return 0, nil
placeholder
	if value < 0 {
		value = 0
placeholder
	res, err := r.sql.ExecContext(ctx,
		"UPDATE users SET concurrency = $1, updated_at = NOW() WHERE id = ANY($2) AND deleted_at IS NULL",
		value, pq.Array(userIDs))
	if err != nil {
		return 0, fmt.Errorf("batch set concurrency: %w", err)
placeholder
	affected, _ := res.RowsAffected()
	return int(affected), nil
placeholder

func (r *userRepository) BatchAddConcurrency(ctx context.Context, userIDs []int64, delta int) (int, error) {
	if len(userIDs) == 0 {
		return 0, nil
placeholder
	res, err := r.sql.ExecContext(ctx,
		"UPDATE users SET concurrency = GREATEST(concurrency + $1, 0), updated_at = NOW() WHERE id = ANY($2) AND deleted_at IS NULL",
		delta, pq.Array(userIDs))
	if err != nil {
		return 0, fmt.Errorf("batch add concurrency: %w", err)
placeholder
	affected, _ := res.RowsAffected()
	return int(affected), nil
placeholder

func (r *userRepository) BatchUpdateLimits(ctx context.Context, userIDs []int64, concurrency, rpmLimit *int) (int, error) {
	if len(userIDs) == 0 || (concurrency == nil && rpmLimit == nil) {
		return 0, nil
placeholder

	setClauses := make([]string, 0, 3)
	args := make([]any, 0, 3)
	if concurrency != nil {
		value := max(*concurrency, 0)
		args = append(args, value)
		setClauses = append(setClauses, fmt.Sprintf("concurrency = $%d", len(args)))
placeholder
	if rpmLimit != nil {
		value := max(*rpmLimit, 0)
		args = append(args, value)
		setClauses = append(setClauses, fmt.Sprintf("rpm_limit = $%d", len(args)))
placeholder
	setClauses = append(setClauses, "updated_at = NOW()")
	args = append(args, pq.Array(userIDs))

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = ANY($%d) AND deleted_at IS NULL",
		strings.Join(setClauses, ", "),
		len(args),
	)
	res, err := r.sql.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("batch update user limits: %w", err)
placeholder
	affected, _ := res.RowsAffected()
	return int(affected), nil
placeholder

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.client.User.Query().Where(userEmailLookupPredicate(email)).Exist(ctx)
placeholder

// emailAliasCandidateLimit 限制一次别名查重最多取回的候选行数。探针都以去点后的
// 本地部分为前缀锚定（见 dotStrippedEmailExpr），正常收件箱的变体只有个位数；
// 上限只是兜底，避免公开未鉴权的注册/发码端点把大表整张读进内存。
const emailAliasCandidateLimit = 50

// ExistsByEmailAlias 见 service.UserRepository。软删除过滤沿用 ExistsByEmail 的默认行为。
func (r *userRepository) ExistsByEmailAlias(ctx context.Context, email string) (bool, error) {
	return existsByEmailAliasWithClient(ctx, clientFromContext(ctx, r.client), email)
placeholder

func existsByEmailAliasWithClient(ctx context.Context, client *dbent.Client, email string) (bool, error) {
	if client == nil {
		return false, nil
placeholder
	probes := service.EmailAliasDedupProbes(email)
	if len(probes) == 0 {
		return false, nil
placeholder

	preds := make([]predicate.User, 0, 2*len(probes))
	for _, probe := range probes {
		preds = append(preds,
			dotStrippedEmailEQ(probe.Local+"@"+probe.Domain),
			// "+后缀"的内容未知，只能按前缀匹配。
			dotStrippedEmailLike(escapeLikeWildcards(probe.Local)+"+%@"+escapeLikeWildcards(probe.Domain)),
		)
placeholder
	candidates, err := client.User.Query().
		Where(dbuser.Or(preds...)).
		Limit(emailAliasCandidateLimit).
		Select(dbuser.FieldEmail).
		Strings(ctx)
	if err != nil {
		return false, err
placeholder

	// 探针会有过度匹配（点号只在 Gmail 家族无意义），最终判定必须回到完整归一化规则。
	identity := service.NormalizeEmailForAliasDedup(email)
	for _, candidate := range candidates {
		if service.NormalizeEmailForAliasDedup(candidate) == identity {
			return true, nil
	placeholder
placeholder
	return false, nil
placeholder

// dotStrippedEmailExpr 渲染下面的表达式：去掉存量邮箱的大小写、首尾空白（与
// userEmailLookupPredicate 的精确匹配口径一致，历史数据存在带空白的行）以及全部点号。
//
//	REPLACE(LOWER(TRIM(email)), '.', '')
//
// 两侧都去点，因此一个域名探针即可同时覆盖 Gmail 点号变体与 FQDN 根点（user@gmail.com.）。
// migrations/190 为同一表达式建了索引。
func dotStrippedEmailExpr(b *entsql.Builder, s *entsql.Selector) *entsql.Builder {
	return b.WriteString("REPLACE(LOWER(TRIM(").
		Ident(s.C(dbuser.FieldEmail)).
		WriteString(")), '.', '')")
placeholder

func dotStrippedEmailEQ(value string) predicate.User {
	return predicate.User(func(s *entsql.Selector) {
		s.Where(entsql.P(func(b *entsql.Builder) {
			dotStrippedEmailExpr(b, s).WriteString(" = ").Arg(value)
	placeholder))
placeholder)
placeholder

func dotStrippedEmailLike(pattern string) predicate.User {
	return predicate.User(func(s *entsql.Selector) {
		s.Where(entsql.P(func(b *entsql.Builder) {
			dotStrippedEmailExpr(b, s).WriteString(" LIKE ").Arg(pattern).WriteString(` ESCAPE '\'`)
	placeholder))
placeholder)
placeholder

// escapeLikeWildcards 转义 LIKE 元字符：本地部分合法可含 % 与 _，不转义会扩大匹配面。
var likeWildcardEscaper = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)

func escapeLikeWildcards(value string) string {
	return likeWildcardEscaper.Replace(value)
placeholder

func ensureNormalizedEmailAvailableWithClient(ctx context.Context, client *dbent.Client, userID int64, email string) error {
	client = clientFromContext(ctx, client)
	if client == nil {
		return nil
placeholder

	matches, err := client.User.Query().
		Where(userEmailLookupPredicate(email)).
		All(ctx)
	if err != nil {
		return err
placeholder
	for _, match := range matches {
		if match.ID != userID {
			return service.ErrEmailExists
	placeholder
placeholder
	return nil
placeholder

func userEmailLookupPredicate(email string) predicate.User {
	normalized := normalizeEmailLookupValue(email)
	if normalized == "" {
		return dbuser.EmailEQ(email)
placeholder
	return predicate.User(func(s *entsql.Selector) {
		s.Where(entsql.P(func(b *entsql.Builder) {
			b.WriteString("LOWER(TRIM(").
				Ident(s.C(dbuser.FieldEmail)).
				WriteString(")) = ").
				Arg(normalized)
	placeholder))
placeholder)
placeholder

func normalizeEmailLookupValue(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
placeholder

func normalizedEmailUniquenessLockKey(email string) string {
	normalized := normalizeEmailLookupValue(email)
	if normalized == "" {
		return ""
placeholder
	return "users:normalized-email:" + normalized
placeholder

// emailAliasUniquenessLockKey 按收件箱身份（而非邮箱字面量）加锁，使同一收件箱的不同
// 别名变体在注册时互斥。
func emailAliasUniquenessLockKey(email string) string {
	identity := service.NormalizeEmailForAliasDedup(email)
	if identity == "" {
		return ""
placeholder
	return "users:email-alias-identity:" + identity
placeholder

func (r *userRepository) AddGroupToAllowedGroups(ctx context.Context, userID int64, groupID int64) error {
	client := clientFromContext(ctx, r.client)
	err := client.UserAllowedGroup.Create().
		SetUserID(userID).
		SetGroupID(groupID).
		OnConflictColumns(userallowedgroup.FieldUserID, userallowedgroup.FieldGroupID).
		DoNothing().
		Exec(ctx)
	if isSQLNoRowsError(err) {
		return nil
placeholder
	return err
placeholder

func (r *userRepository) RemoveGroupFromAllowedGroups(ctx context.Context, groupID int64) (int64, error) {
	// 仅操作 user_allowed_groups 联接表，legacy users.allowed_groups 列已弃用。
	affected, err := r.client.UserAllowedGroup.Delete().
		Where(userallowedgroup.GroupIDEQ(groupID)).
		Exec(ctx)
	if err != nil {
		return 0, err
placeholder
	return int64(affected), nil
placeholder

// RemoveGroupFromUserAllowedGroups 移除单个用户的指定分组权限
func (r *userRepository) RemoveGroupFromUserAllowedGroups(ctx context.Context, userID int64, groupID int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.UserAllowedGroup.Delete().
		Where(userallowedgroup.UserIDEQ(userID), userallowedgroup.GroupIDEQ(groupID)).
		Exec(ctx)
	return err
placeholder

func (r *userRepository) GetFirstAdmin(ctx context.Context) (*service.User, error) {
	m, err := r.client.User.Query().
		Where(
			dbuser.RoleEQ(service.RoleAdmin),
			dbuser.StatusEQ(service.StatusActive),
		).
		Order(dbent.Asc(dbuser.FieldID)).
		First(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder

	out := userEntityToService(m)
	groups, err := r.loadAllowedGroups(ctx, []int64{m.IDplaceholder)
	if err != nil {
		return nil, err
placeholder
	if v, ok := groups[m.ID]; ok {
		out.AllowedGroups = v
placeholder
	return out, nil
placeholder

func (r *userRepository) loadAllowedGroups(ctx context.Context, userIDs []int64) (map[int64][]int64, error) {
	out := make(map[int64][]int64, len(userIDs))
	if len(userIDs) == 0 {
		return out, nil
placeholder

	rows, err := r.client.UserAllowedGroup.Query().
		Where(userallowedgroup.UserIDIn(userIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	for i := range rows {
		out[rows[i].UserID] = append(out[rows[i].UserID], rows[i].GroupID)
placeholder

	for userID := range out {
		sort.Slice(out[userID], func(i, j int) bool { return out[userID][i] < out[userID][j] placeholder)
placeholder

	return out, nil
placeholder

// syncUserAllowedGroupsWithClient 在 ent client/事务内同步用户允许分组：
// 仅操作 user_allowed_groups 联接表，legacy users.allowed_groups 列已弃用。
func (r *userRepository) syncUserAllowedGroupsWithClient(ctx context.Context, client *dbent.Client, userID int64, groupIDs []int64) error {
	if client == nil {
		return nil
placeholder

	existingRows, err := client.UserAllowedGroup.Query().
		Where(userallowedgroup.UserIDEQ(userID)).
		All(ctx)
	if err != nil {
		return err
placeholder

	desired := make(map[int64]struct{placeholder, len(groupIDs))
	for _, id := range groupIDs {
		if id <= 0 {
			continue
	placeholder
		desired[id] = struct{placeholder{placeholder
placeholder

	existing := make(map[int64]struct{placeholder, len(existingRows))
	removed := make([]int64, 0)
	for _, row := range existingRows {
		existing[row.GroupID] = struct{placeholder{placeholder
		if _, keep := desired[row.GroupID]; !keep {
			removed = append(removed, row.GroupID)
	placeholder
placeholder
	if len(removed) > 0 {
		if _, err := client.UserAllowedGroup.Delete().
			Where(userallowedgroup.UserIDEQ(userID), userallowedgroup.GroupIDIn(removed...)).
			Exec(ctx); err != nil {
			return err
	placeholder
placeholder

	creates := make([]*dbent.UserAllowedGroupCreate, 0, len(desired))
	for groupID := range desired {
		if _, present := existing[groupID]; !present {
			creates = append(creates, client.UserAllowedGroup.Create().SetUserID(userID).SetGroupID(groupID))
	placeholder
placeholder
	if len(creates) > 0 {
		if err := client.UserAllowedGroup.
			CreateBulk(creates...).
			OnConflictColumns(userallowedgroup.FieldUserID, userallowedgroup.FieldGroupID).
			DoNothing().
			Exec(ctx); err != nil {
			if isSQLNoRowsError(err) {
				return nil
		placeholder
			return err
	placeholder
placeholder

	return nil
placeholder

func applyUserEntityToService(dst *service.User, src *dbent.User) {
	if dst == nil || src == nil {
		return
placeholder
	dst.ID = src.ID
	dst.SignupSource = src.SignupSource
	dst.LastLoginAt = src.LastLoginAt
	dst.LastActiveAt = src.LastActiveAt
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
placeholder

func userSignupSourceOrDefault(signupSource string) string {
	switch strings.TrimSpace(strings.ToLower(signupSource)) {
	case "", "email":
		return "email"
	case "linuxdo", "wechat", "oidc", "dingtalk":
		return strings.TrimSpace(strings.ToLower(signupSource))
	default:
		return "email"
placeholder
placeholder

// marshalExtraEmails serializes notify email entries to JSON for storage.
func marshalExtraEmails(entries []service.NotifyEmailEntry) string {
	return service.MarshalNotifyEmails(entries)
placeholder

// UpdateTotpSecret 更新用户的 TOTP 加密密钥
func (r *userRepository) UpdateTotpSecret(ctx context.Context, userID int64, encryptedSecret *string) error {
	client := clientFromContext(ctx, r.client)
	update := client.User.UpdateOneID(userID)
	if encryptedSecret == nil {
		update = update.ClearTotpSecretEncrypted()
placeholder else {
		update = update.SetTotpSecretEncrypted(*encryptedSecret)
placeholder
	_, err := update.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	return nil
placeholder

// EnableTotp 启用用户的 TOTP 双因素认证
func (r *userRepository) EnableTotp(ctx context.Context, userID int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.User.UpdateOneID(userID).
		SetTotpEnabled(true).
		SetTotpEnabledAt(time.Now()).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	return nil
placeholder

// DisableTotp 禁用用户的 TOTP 双因素认证
func (r *userRepository) DisableTotp(ctx context.Context, userID int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.User.UpdateOneID(userID).
		SetTotpEnabled(false).
		ClearTotpEnabledAt().
		ClearTotpSecretEncrypted().
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrUserNotFound, nil)
placeholder
	return nil
placeholder
