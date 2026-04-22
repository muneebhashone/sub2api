package repository

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/authidentity"
	"github.com/Wei-Shaw/sub2api/ent/authidentitychannel"
	"github.com/Wei-Shaw/sub2api/ent/identityadoptiondecision"
	dbpredicate "github.com/Wei-Shaw/sub2api/ent/predicate"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

var (
	ErrAuthIdentityOwnershipConflict = infraerrors.Conflict(
		"AUTH_IDENTITY_OWNERSHIP_CONFLICT",
		"auth identity already belongs to another user",
	)
	ErrAuthIdentityChannelOwnershipConflict = infraerrors.Conflict(
		"AUTH_IDENTITY_CHANNEL_OWNERSHIP_CONFLICT",
		"auth identity channel already belongs to another user",
	)
	ErrAuthIdentityChannelProviderMismatch = infraerrors.BadRequest(
		"AUTH_IDENTITY_CHANNEL_PROVIDER_MISMATCH",
		"auth identity channel provider must match canonical identity",
	)
)

type ProviderGrantReason string

const (
	ProviderGrantReasonSignup    ProviderGrantReason = "signup"
	ProviderGrantReasonFirstBind ProviderGrantReason = "first_bind"
)

type AuthIdentityKey struct {
	ProviderType    string
	ProviderKey     string
	ProviderSubject string
placeholder

type AuthIdentityChannelKey struct {
	ProviderType   string
	ProviderKey    string
	Channel        string
	ChannelAppID   string
	ChannelSubject string
placeholder

type CreateAuthIdentityInput struct {
	UserID          int64
	Canonical       AuthIdentityKey
	Channel         *AuthIdentityChannelKey
	Issuer          *string
	VerifiedAt      *time.Time
	Metadata        map[string]any
	ChannelMetadata map[string]any
placeholder

type BindAuthIdentityInput = CreateAuthIdentityInput

type CreateAuthIdentityResult struct {
	Identity *dbent.AuthIdentity
	Channel  *dbent.AuthIdentityChannel
placeholder

func (r *CreateAuthIdentityResult) IdentityRef() AuthIdentityKey {
	if r == nil || r.Identity == nil {
		return AuthIdentityKey{placeholder
placeholder
	return AuthIdentityKey{
		ProviderType:    r.Identity.ProviderType,
		ProviderKey:     r.Identity.ProviderKey,
		ProviderSubject: r.Identity.ProviderSubject,
placeholder
placeholder

func (r *CreateAuthIdentityResult) ChannelRef() *AuthIdentityChannelKey {
	if r == nil || r.Channel == nil {
		return nil
placeholder
	return &AuthIdentityChannelKey{
		ProviderType:   r.Channel.ProviderType,
		ProviderKey:    r.Channel.ProviderKey,
		Channel:        r.Channel.Channel,
		ChannelAppID:   r.Channel.ChannelAppID,
		ChannelSubject: r.Channel.ChannelSubject,
placeholder
placeholder

type UserAuthIdentityLookup struct {
	User     *dbent.User
	Identity *dbent.AuthIdentity
	Channel  *dbent.AuthIdentityChannel
placeholder

type ProviderGrantRecordInput struct {
	UserID       int64
	ProviderType string
	GrantReason  ProviderGrantReason
placeholder

type IdentityAdoptionDecisionInput struct {
	PendingAuthSessionID int64
	IdentityID           *int64
	AdoptDisplayName     bool
	AdoptAvatar          bool
placeholder

type sqlQueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
placeholder

var repositoryScopedKeyLocks = newScopedKeyLockRegistry()

type scopedKeyLockRegistry struct {
	mu    sync.Mutex
	locks map[string]*scopedKeyLockEntry
placeholder

type scopedKeyLockEntry struct {
	mu   sync.Mutex
	refs int
placeholder

func newScopedKeyLockRegistry() *scopedKeyLockRegistry {
	return &scopedKeyLockRegistry{
		locks: make(map[string]*scopedKeyLockEntry),
placeholder
placeholder

func (r *scopedKeyLockRegistry) lock(keys ...string) func() {
	normalized := normalizeLockKeys(keys...)
	if len(normalized) == 0 {
		return func() {placeholder
placeholder

	entries := make([]*scopedKeyLockEntry, 0, len(normalized))
	r.mu.Lock()
	for _, key := range normalized {
		entry := r.locks[key]
		if entry == nil {
			entry = &scopedKeyLockEntry{placeholder
			r.locks[key] = entry
	placeholder
		entry.refs++
		entries = append(entries, entry)
placeholder
	r.mu.Unlock()

	for _, entry := range entries {
		entry.mu.Lock()
placeholder

	return func() {
		for i := len(entries) - 1; i >= 0; i-- {
			entries[i].mu.Unlock()
	placeholder

		r.mu.Lock()
		defer r.mu.Unlock()
		for idx, key := range normalized {
			entry := entries[idx]
			entry.refs--
			if entry.refs == 0 {
				delete(r.locks, key)
		placeholder
	placeholder
placeholder
placeholder

func normalizeLockKeys(keys ...string) []string {
	if len(keys) == 0 {
		return nil
placeholder

	deduped := make(map[string]struct{placeholder, len(keys))
	for _, key := range keys {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
	placeholder
		deduped[trimmed] = struct{placeholder{placeholder
placeholder
	if len(deduped) == 0 {
		return nil
placeholder

	normalized := make([]string, 0, len(deduped))
	for key := range deduped {
		normalized = append(normalized, key)
placeholder
	sort.Strings(normalized)
	return normalized
placeholder

func advisoryLockHash(key string) int64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(key))
	return int64(hasher.Sum64())
placeholder

func lockRepositoryScopedKeys(ctx context.Context, client *dbent.Client, exec sqlQueryExecutor, keys ...string) (func(), error) {
	release := repositoryScopedKeyLocks.lock(keys...)
	normalized := normalizeLockKeys(keys...)
	if len(normalized) == 0 || client == nil || exec == nil || client.Driver().Dialect() != dialect.Postgres {
		return release, nil
placeholder

	for _, key := range normalized {
		rows, err := exec.QueryContext(ctx, "SELECT pg_advisory_xact_lock($1)", advisoryLockHash(key))
		if err != nil {
			release()
			return nil, err
	placeholder
		_ = rows.Close()
placeholder
	return release, nil
placeholder

func (r *userRepository) WithUserProfileIdentityTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	if dbent.TxFromContext(ctx) != nil {
		return fn(ctx)
placeholder

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
placeholder
	defer func() { _ = tx.Rollback() placeholder()

	txCtx := dbent.NewTxContext(ctx, tx)
	if err := fn(txCtx); err != nil {
		return err
placeholder
	return tx.Commit()
placeholder

func (r *userRepository) CreateAuthIdentity(ctx context.Context, input CreateAuthIdentityInput) (*CreateAuthIdentityResult, error) {
	if err := validateAuthIdentityChannelProviderMatch(input.Canonical, input.Channel); err != nil {
		return nil, err
placeholder

	client := clientFromContext(ctx, r.client)

	create := client.AuthIdentity.Create().
		SetUserID(input.UserID).
		SetProviderType(strings.TrimSpace(input.Canonical.ProviderType)).
		SetProviderKey(strings.TrimSpace(input.Canonical.ProviderKey)).
		SetProviderSubject(strings.TrimSpace(input.Canonical.ProviderSubject)).
		SetMetadata(copyMetadata(input.Metadata)).
		SetNillableIssuer(input.Issuer).
		SetNillableVerifiedAt(input.VerifiedAt)

	identity, err := create.Save(ctx)
	if err != nil {
		return nil, err
placeholder

	var channel *dbent.AuthIdentityChannel
	if input.Channel != nil {
		channel, err = client.AuthIdentityChannel.Create().
			SetIdentityID(identity.ID).
			SetProviderType(strings.TrimSpace(input.Channel.ProviderType)).
			SetProviderKey(strings.TrimSpace(input.Channel.ProviderKey)).
			SetChannel(strings.TrimSpace(input.Channel.Channel)).
			SetChannelAppID(strings.TrimSpace(input.Channel.ChannelAppID)).
			SetChannelSubject(strings.TrimSpace(input.Channel.ChannelSubject)).
			SetMetadata(copyMetadata(input.ChannelMetadata)).
			Save(ctx)
		if err != nil {
			return nil, err
	placeholder
placeholder

	return &CreateAuthIdentityResult{Identity: identity, Channel: channelplaceholder, nil
placeholder

func (r *userRepository) GetUserByCanonicalIdentity(ctx context.Context, key AuthIdentityKey) (*UserAuthIdentityLookup, error) {
	identity, err := clientFromContext(ctx, r.client).AuthIdentity.Query().
		Where(
			authidentity.ProviderTypeEQ(strings.TrimSpace(key.ProviderType)),
			authidentity.ProviderKeyEQ(strings.TrimSpace(key.ProviderKey)),
			authidentity.ProviderSubjectEQ(strings.TrimSpace(key.ProviderSubject)),
		).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
placeholder

	return &UserAuthIdentityLookup{
		User:     identity.Edges.User,
		Identity: identity,
placeholder, nil
placeholder

func (r *userRepository) GetUserByChannelIdentity(ctx context.Context, key AuthIdentityChannelKey) (*UserAuthIdentityLookup, error) {
	channel, err := clientFromContext(ctx, r.client).AuthIdentityChannel.Query().
		Where(
			authidentitychannel.ProviderTypeEQ(strings.TrimSpace(key.ProviderType)),
			authidentitychannel.ProviderKeyEQ(strings.TrimSpace(key.ProviderKey)),
			authidentitychannel.ChannelEQ(strings.TrimSpace(key.Channel)),
			authidentitychannel.ChannelAppIDEQ(strings.TrimSpace(key.ChannelAppID)),
			authidentitychannel.ChannelSubjectEQ(strings.TrimSpace(key.ChannelSubject)),
		).
		WithIdentity(func(q *dbent.AuthIdentityQuery) {
			q.WithUser()
	placeholder).
		Only(ctx)
	if err != nil {
		return nil, err
placeholder

	return &UserAuthIdentityLookup{
		User:     channel.Edges.Identity.Edges.User,
		Identity: channel.Edges.Identity,
		Channel:  channel,
placeholder, nil
placeholder

func (r *userRepository) ListUserAuthIdentities(ctx context.Context, userID int64) ([]service.UserAuthIdentityRecord, error) {
	identities, err := clientFromContext(ctx, r.client).AuthIdentity.Query().
		Where(authidentity.UserIDEQ(userID)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	records := make([]service.UserAuthIdentityRecord, 0, len(identities))
	for _, identity := range identities {
		if identity == nil {
			continue
	placeholder
		records = append(records, service.UserAuthIdentityRecord{
			ProviderType:    strings.TrimSpace(identity.ProviderType),
			ProviderKey:     strings.TrimSpace(identity.ProviderKey),
			ProviderSubject: strings.TrimSpace(identity.ProviderSubject),
			VerifiedAt:      identity.VerifiedAt,
			Issuer:          identity.Issuer,
			Metadata:        copyMetadata(identity.Metadata),
			CreatedAt:       identity.CreatedAt,
			UpdatedAt:       identity.UpdatedAt,
	placeholder)
placeholder

	return records, nil
placeholder

func (r *userRepository) UnbindUserAuthProvider(ctx context.Context, userID int64, provider string) error {
	provider = strings.ToLower(strings.TrimSpace(provider))
	if provider == "" || provider == "email" {
		return service.ErrIdentityProviderInvalid
placeholder

	return r.WithUserProfileIdentityTx(ctx, func(txCtx context.Context) error {
		client := clientFromContext(txCtx, r.client)
		identityIDs, err := client.AuthIdentity.Query().
			Where(
				authidentity.UserIDEQ(userID),
				authidentity.ProviderTypeEQ(provider),
			).
			IDs(txCtx)
		if err != nil {
			return err
	placeholder
		if len(identityIDs) == 0 {
			return nil
	placeholder

		if _, err := client.IdentityAdoptionDecision.Update().
			Where(identityadoptiondecision.IdentityIDIn(identityIDs...)).
			ClearIdentityID().
			Save(txCtx); err != nil {
			return err
	placeholder
		if _, err := client.AuthIdentityChannel.Delete().
			Where(authidentitychannel.IdentityIDIn(identityIDs...)).
			Exec(txCtx); err != nil {
			return err
	placeholder
		_, err = client.AuthIdentity.Delete().
			Where(
				authidentity.UserIDEQ(userID),
				authidentity.ProviderTypeEQ(provider),
			).
			Exec(txCtx)
		return err
placeholder)
placeholder

func (r *userRepository) BindAuthIdentityToUser(ctx context.Context, input BindAuthIdentityInput) (*CreateAuthIdentityResult, error) {
	if err := validateAuthIdentityChannelProviderMatch(input.Canonical, input.Channel); err != nil {
		return nil, err
placeholder

	var result *CreateAuthIdentityResult
	err := r.WithUserProfileIdentityTx(ctx, func(txCtx context.Context) error {
		client := clientFromContext(txCtx, r.client)
		canonical := input.Canonical

		identityRecords, err := client.AuthIdentity.Query().
			Where(
				authidentity.ProviderTypeEQ(strings.TrimSpace(canonical.ProviderType)),
				authidentity.ProviderKeyIn(compatibleIdentityProviderKeys(canonical.ProviderType, canonical.ProviderKey)...),
				authidentity.ProviderSubjectEQ(strings.TrimSpace(canonical.ProviderSubject)),
			).
			All(txCtx)
		if err != nil {
			return err
	placeholder
		identity := selectOwnedCompatibleIdentity(identityRecords, input.UserID)
		if identity == nil && hasCompatibleIdentityConflict(identityRecords, input.UserID) {
			return ErrAuthIdentityOwnershipConflict
	placeholder
		if identity == nil {
			identity, err = client.AuthIdentity.Create().
				SetUserID(input.UserID).
				SetProviderType(strings.TrimSpace(canonical.ProviderType)).
				SetProviderKey(strings.TrimSpace(canonical.ProviderKey)).
				SetProviderSubject(strings.TrimSpace(canonical.ProviderSubject)).
				SetMetadata(copyMetadata(input.Metadata)).
				SetNillableIssuer(input.Issuer).
				SetNillableVerifiedAt(input.VerifiedAt).
				Save(txCtx)
			if err != nil {
				return err
		placeholder
	placeholder else {
			targetProviderKey := canonicalizeCompatibleIdentityProviderKey(canonical.ProviderType, identity.ProviderKey, canonical.ProviderKey)
			update := client.AuthIdentity.UpdateOneID(identity.ID)
			if targetProviderKey != "" && !strings.EqualFold(targetProviderKey, identity.ProviderKey) {
				update = update.SetProviderKey(targetProviderKey)
		placeholder
			if input.Metadata != nil {
				update = update.SetMetadata(copyMetadata(input.Metadata))
		placeholder
			if input.Issuer != nil {
				update = update.SetIssuer(strings.TrimSpace(*input.Issuer))
		placeholder
			if input.VerifiedAt != nil {
				update = update.SetVerifiedAt(*input.VerifiedAt)
		placeholder
			identity, err = update.Save(txCtx)
			if err != nil {
				return err
		placeholder
	placeholder

		var channel *dbent.AuthIdentityChannel
		if input.Channel != nil {
			channelRecords, err := client.AuthIdentityChannel.Query().
				Where(
					authidentitychannel.ProviderTypeEQ(strings.TrimSpace(input.Channel.ProviderType)),
					authidentitychannel.ProviderKeyIn(compatibleIdentityProviderKeys(input.Channel.ProviderType, input.Channel.ProviderKey)...),
					authidentitychannel.ChannelEQ(strings.TrimSpace(input.Channel.Channel)),
					authidentitychannel.ChannelAppIDEQ(strings.TrimSpace(input.Channel.ChannelAppID)),
					authidentitychannel.ChannelSubjectEQ(strings.TrimSpace(input.Channel.ChannelSubject)),
				).
				WithIdentity().
				All(txCtx)
			if err != nil {
				return err
		placeholder
			channel = selectOwnedCompatibleChannel(channelRecords, input.UserID)
			if channel == nil && hasCompatibleChannelConflict(channelRecords, input.UserID) {
				return ErrAuthIdentityChannelOwnershipConflict
		placeholder
			if channel == nil {
				channel, err = client.AuthIdentityChannel.Create().
					SetIdentityID(identity.ID).
					SetProviderType(strings.TrimSpace(input.Channel.ProviderType)).
					SetProviderKey(strings.TrimSpace(input.Channel.ProviderKey)).
					SetChannel(strings.TrimSpace(input.Channel.Channel)).
					SetChannelAppID(strings.TrimSpace(input.Channel.ChannelAppID)).
					SetChannelSubject(strings.TrimSpace(input.Channel.ChannelSubject)).
					SetMetadata(copyMetadata(input.ChannelMetadata)).
					Save(txCtx)
				if err != nil {
					return err
			placeholder
		placeholder else {
				targetProviderKey := canonicalizeCompatibleIdentityProviderKey(input.Channel.ProviderType, channel.ProviderKey, input.Channel.ProviderKey)
				update := client.AuthIdentityChannel.UpdateOneID(channel.ID).
					SetIdentityID(identity.ID)
				if targetProviderKey != "" && !strings.EqualFold(targetProviderKey, channel.ProviderKey) {
					update = update.SetProviderKey(targetProviderKey)
			placeholder
				if input.ChannelMetadata != nil {
					update = update.SetMetadata(copyMetadata(input.ChannelMetadata))
			placeholder
				channel, err = update.Save(txCtx)
				if err != nil {
					return err
			placeholder
		placeholder
	placeholder

		result = &CreateAuthIdentityResult{Identity: identity, Channel: channelplaceholder
		return nil
placeholder)
	if err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

func compatibleIdentityProviderKeys(providerType, providerKey string) []string {
	providerType = strings.TrimSpace(strings.ToLower(providerType))
	providerKey = strings.TrimSpace(providerKey)
	if providerKey == "" {
		return []string{providerKeyplaceholder
placeholder
	if providerType != "wechat" {
		return []string{providerKeyplaceholder
placeholder
	keys := []string{providerKeyplaceholder
	if !strings.EqualFold(providerKey, "wechat-main") {
		keys = append(keys, "wechat-main")
placeholder
	if !strings.EqualFold(providerKey, "wechat") {
		keys = append(keys, "wechat")
placeholder
	return keys
placeholder

func canonicalizeCompatibleIdentityProviderKey(providerType, existingKey, requestedKey string) string {
	providerType = strings.TrimSpace(strings.ToLower(providerType))
	existingKey = strings.TrimSpace(existingKey)
	requestedKey = strings.TrimSpace(requestedKey)
	if providerType != "wechat" {
		if requestedKey != "" {
			return requestedKey
	placeholder
		return existingKey
placeholder
	if strings.EqualFold(existingKey, "wechat") || strings.EqualFold(existingKey, "wechat-main") || strings.EqualFold(requestedKey, "wechat-main") {
		return "wechat-main"
placeholder
	if requestedKey != "" {
		return requestedKey
placeholder
	return existingKey
placeholder

func compatibleIdentityProviderKeyRank(providerType, providerKey string) int {
	providerType = strings.TrimSpace(strings.ToLower(providerType))
	providerKey = strings.TrimSpace(providerKey)
	if providerType != "wechat" {
		return 0
placeholder
	switch {
	case strings.EqualFold(providerKey, "wechat-main"):
		return 0
	case strings.EqualFold(providerKey, "wechat"):
		return 2
	default:
		return 1
placeholder
placeholder

func selectOwnedCompatibleIdentity(records []*dbent.AuthIdentity, userID int64) *dbent.AuthIdentity {
	var selected *dbent.AuthIdentity
	for _, record := range records {
		if record.UserID != userID {
			continue
	placeholder
		if selected == nil || compatibleIdentityProviderKeyRank(record.ProviderType, record.ProviderKey) < compatibleIdentityProviderKeyRank(selected.ProviderType, selected.ProviderKey) {
			selected = record
	placeholder
placeholder
	return selected
placeholder

func hasCompatibleIdentityConflict(records []*dbent.AuthIdentity, userID int64) bool {
	for _, record := range records {
		if record.UserID != userID {
			return true
	placeholder
placeholder
	return false
placeholder

func selectOwnedCompatibleChannel(records []*dbent.AuthIdentityChannel, userID int64) *dbent.AuthIdentityChannel {
	var selected *dbent.AuthIdentityChannel
	for _, record := range records {
		if record.Edges.Identity == nil || record.Edges.Identity.UserID != userID {
			continue
	placeholder
		if selected == nil || compatibleIdentityProviderKeyRank(record.ProviderType, record.ProviderKey) < compatibleIdentityProviderKeyRank(selected.ProviderType, selected.ProviderKey) {
			selected = record
	placeholder
placeholder
	return selected
placeholder

func hasCompatibleChannelConflict(records []*dbent.AuthIdentityChannel, userID int64) bool {
	for _, record := range records {
		if record.Edges.Identity != nil && record.Edges.Identity.UserID != userID {
			return true
	placeholder
placeholder
	return false
placeholder

func (r *userRepository) RecordProviderGrant(ctx context.Context, input ProviderGrantRecordInput) (bool, error) {
	exec := txAwareSQLExecutor(ctx, r.sql, r.client)
	if exec == nil {
		return false, fmt.Errorf("sql executor is not configured")
placeholder

	result, err := exec.ExecContext(ctx, `
INSERT INTO user_provider_default_grants (user_id, provider_type, grant_reason)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, provider_type, grant_reason) DO NOTHING`,
		input.UserID,
		strings.TrimSpace(input.ProviderType),
		string(input.GrantReason),
	)
	if err != nil {
		return false, err
placeholder
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
placeholder
	return affected > 0, nil
placeholder

func (r *userRepository) UpsertIdentityAdoptionDecision(ctx context.Context, input IdentityAdoptionDecisionInput) (*dbent.IdentityAdoptionDecision, error) {
	var result *dbent.IdentityAdoptionDecision
	err := r.WithUserProfileIdentityTx(ctx, func(txCtx context.Context) error {
		client := clientFromContext(txCtx, r.client)
		releaseLocks, err := lockRepositoryScopedKeys(
			txCtx,
			client,
			txAwareSQLExecutor(txCtx, r.sql, r.client),
			identityAdoptionDecisionLockKeys(input.PendingAuthSessionID, input.IdentityID)...,
		)
		if err != nil {
			return err
	placeholder
		defer releaseLocks()

		if input.IdentityID != nil && *input.IdentityID > 0 {
			if _, err := client.IdentityAdoptionDecision.Update().
				Where(
					identityadoptiondecision.IdentityIDEQ(*input.IdentityID),
					dbpredicate.IdentityAdoptionDecision(func(s *entsql.Selector) {
						col := s.C(identityadoptiondecision.FieldPendingAuthSessionID)
						s.Where(entsql.Or(
							entsql.IsNull(col),
							entsql.NEQ(col, input.PendingAuthSessionID),
						))
				placeholder),
				).
				ClearIdentityID().
				Save(txCtx); err != nil {
				return err
		placeholder
	placeholder

		create := client.IdentityAdoptionDecision.Create().
			SetPendingAuthSessionID(input.PendingAuthSessionID).
			SetAdoptDisplayName(input.AdoptDisplayName).
			SetAdoptAvatar(input.AdoptAvatar).
			SetDecidedAt(time.Now().UTC())
		if input.IdentityID != nil && *input.IdentityID > 0 {
			create = create.SetIdentityID(*input.IdentityID)
	placeholder

		decisionID, err := create.
			OnConflictColumns(identityadoptiondecision.FieldPendingAuthSessionID).
			UpdateNewValues().
			ID(txCtx)
		if err != nil {
			return err
	placeholder

		result, err = client.IdentityAdoptionDecision.Get(txCtx, decisionID)
		return err
placeholder)
	if err != nil {
		return nil, err
placeholder
	return result, nil
placeholder

func identityAdoptionDecisionLockKeys(pendingAuthSessionID int64, identityID *int64) []string {
	keys := []string{fmt.Sprintf("identity-adoption:pending:%d", pendingAuthSessionID)placeholder
	if identityID != nil && *identityID > 0 {
		keys = append(keys, fmt.Sprintf("identity-adoption:identity:%d", *identityID))
placeholder
	return keys
placeholder

func (r *userRepository) GetIdentityAdoptionDecisionByPendingAuthSessionID(ctx context.Context, pendingAuthSessionID int64) (*dbent.IdentityAdoptionDecision, error) {
	return clientFromContext(ctx, r.client).IdentityAdoptionDecision.Query().
		Where(identityadoptiondecision.PendingAuthSessionIDEQ(pendingAuthSessionID)).
		Only(ctx)
placeholder

func (r *userRepository) UpdateUserLastLoginAt(ctx context.Context, userID int64, loginAt time.Time) error {
	_, err := clientFromContext(ctx, r.client).User.UpdateOneID(userID).
		SetLastLoginAt(loginAt).
		Save(ctx)
	return err
placeholder

func (r *userRepository) UpdateUserLastActiveAt(ctx context.Context, userID int64, activeAt time.Time) error {
	_, err := clientFromContext(ctx, r.client).User.UpdateOneID(userID).
		SetLastActiveAt(activeAt).
		Save(ctx)
	return err
placeholder

func (r *userRepository) GetUserAvatar(ctx context.Context, userID int64) (*service.UserAvatar, error) {
	exec, err := r.userProfileIdentitySQL(ctx)
	if err != nil {
		return nil, err
placeholder

	rows, err := exec.QueryContext(ctx, `
SELECT storage_provider, storage_key, url, content_type, byte_size, sha256
FROM user_avatars
WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	if !rows.Next() {
		return nil, rows.Err()
placeholder

	var avatar service.UserAvatar
	if err := rows.Scan(
		&avatar.StorageProvider,
		&avatar.StorageKey,
		&avatar.URL,
		&avatar.ContentType,
		&avatar.ByteSize,
		&avatar.SHA256,
	); err != nil {
		return nil, err
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return &avatar, nil
placeholder

func (r *userRepository) UpsertUserAvatar(ctx context.Context, userID int64, input service.UpsertUserAvatarInput) (*service.UserAvatar, error) {
	exec, err := r.userProfileIdentitySQL(ctx)
	if err != nil {
		return nil, err
placeholder

	_, err = exec.ExecContext(ctx, `
INSERT INTO user_avatars (user_id, storage_provider, storage_key, url, content_type, byte_size, sha256, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
ON CONFLICT (user_id) DO UPDATE SET
	storage_provider = EXCLUDED.storage_provider,
	storage_key = EXCLUDED.storage_key,
	url = EXCLUDED.url,
	content_type = EXCLUDED.content_type,
	byte_size = EXCLUDED.byte_size,
	sha256 = EXCLUDED.sha256,
	updated_at = NOW()`,
		userID,
		strings.TrimSpace(input.StorageProvider),
		strings.TrimSpace(input.StorageKey),
		strings.TrimSpace(input.URL),
		strings.TrimSpace(input.ContentType),
		input.ByteSize,
		strings.TrimSpace(input.SHA256),
	)
	if err != nil {
		return nil, err
placeholder

	return &service.UserAvatar{
		StorageProvider: strings.TrimSpace(input.StorageProvider),
		StorageKey:      strings.TrimSpace(input.StorageKey),
		URL:             strings.TrimSpace(input.URL),
		ContentType:     strings.TrimSpace(input.ContentType),
		ByteSize:        input.ByteSize,
		SHA256:          strings.TrimSpace(input.SHA256),
placeholder, nil
placeholder

func (r *userRepository) DeleteUserAvatar(ctx context.Context, userID int64) error {
	exec, err := r.userProfileIdentitySQL(ctx)
	if err != nil {
		return err
placeholder
	_, err = exec.ExecContext(ctx, `DELETE FROM user_avatars WHERE user_id = $1`, userID)
	return err
placeholder

func copyMetadata(in map[string]any) map[string]any {
	if len(in) == 0 {
		return map[string]any{placeholder
placeholder
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
placeholder
	return out
placeholder

func validateAuthIdentityChannelProviderMatch(canonical AuthIdentityKey, channel *AuthIdentityChannelKey) error {
	if channel == nil {
		return nil
placeholder

	canonicalProviderType := strings.TrimSpace(canonical.ProviderType)
	canonicalProviderKey := strings.TrimSpace(canonical.ProviderKey)
	channelProviderType := strings.TrimSpace(channel.ProviderType)
	channelProviderKey := strings.TrimSpace(channel.ProviderKey)

	if canonicalProviderType != channelProviderType || canonicalProviderKey != channelProviderKey {
		return ErrAuthIdentityChannelProviderMismatch
placeholder

	return nil
placeholder

func txAwareSQLExecutor(ctx context.Context, fallback sqlExecutor, client *dbent.Client) sqlQueryExecutor {
	if tx := dbent.TxFromContext(ctx); tx != nil {
		if exec := sqlExecutorFromEntClient(tx.Client()); exec != nil {
			return exec
	placeholder
placeholder
	if fallback != nil {
		return fallback
placeholder
	return sqlExecutorFromEntClient(client)
placeholder

func (r *userRepository) userProfileIdentitySQL(ctx context.Context) (sqlQueryExecutor, error) {
	exec := txAwareSQLExecutor(ctx, r.sql, r.client)
	if exec == nil {
		return nil, fmt.Errorf("sql executor is not configured")
placeholder
	return exec, nil
placeholder

func sqlExecutorFromEntClient(client *dbent.Client) sqlQueryExecutor {
	if client == nil {
		return nil
placeholder

	clientValue := reflect.ValueOf(client).Elem()
	configValue := clientValue.FieldByName("config")
	driverValue := configValue.FieldByName("driver")
	if !driverValue.IsValid() {
		return nil
placeholder

	driver := reflect.NewAt(driverValue.Type(), unsafe.Pointer(driverValue.UnsafeAddr())).Elem().Interface()
	exec, ok := driver.(sqlQueryExecutor)
	if !ok {
		return nil
placeholder
	return exec
placeholder
