package repository

import (
	"context"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/promocode"
	"github.com/Wei-Shaw/sub2api/ent/promocodeusage"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"

	entsql "entgo.io/ent/dialect/sql"
)

type promoCodeRepository struct {
	client *dbent.Client
placeholder

func NewPromoCodeRepository(client *dbent.Client) service.PromoCodeRepository {
	return &promoCodeRepository{client: clientplaceholder
placeholder

func (r *promoCodeRepository) Create(ctx context.Context, code *service.PromoCode) error {
	client := clientFromContext(ctx, r.client)
	builder := client.PromoCode.Create().
		SetCode(code.Code).
		SetBonusAmount(code.BonusAmount).
		SetMaxUses(code.MaxUses).
		SetUsedCount(code.UsedCount).
		SetStatus(code.Status).
		SetNotes(code.Notes)

	if code.ExpiresAt != nil {
		builder.SetExpiresAt(*code.ExpiresAt)
placeholder

	created, err := builder.Save(ctx)
	if err != nil {
		return err
placeholder

	code.ID = created.ID
	code.CreatedAt = created.CreatedAt
	code.UpdatedAt = created.UpdatedAt
	return nil
placeholder

func (r *promoCodeRepository) GetByID(ctx context.Context, id int64) (*service.PromoCode, error) {
	m, err := r.client.PromoCode.Query().
		Where(promocode.IDEQ(id)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrPromoCodeNotFound
	placeholder
		return nil, err
placeholder
	return promoCodeEntityToService(m), nil
placeholder

func (r *promoCodeRepository) GetByCode(ctx context.Context, code string) (*service.PromoCode, error) {
	m, err := r.client.PromoCode.Query().
		Where(promocode.CodeEqualFold(code)).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrPromoCodeNotFound
	placeholder
		return nil, err
placeholder
	return promoCodeEntityToService(m), nil
placeholder

func (r *promoCodeRepository) GetByCodeForUpdate(ctx context.Context, code string) (*service.PromoCode, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.PromoCode.Query().
		Where(promocode.CodeEqualFold(code)).
		ForUpdate().
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, service.ErrPromoCodeNotFound
	placeholder
		return nil, err
placeholder
	return promoCodeEntityToService(m), nil
placeholder

// Update 写入管理员可编辑的字段。
//
// 这里刻意不写 used_count：它由兑换路径的 IncrementUsedCount 原子递增，
// 而 used_count >= max_uses 正是"优惠码用完了"的判定依据。若管理员编辑
// （改有效期、改额度……）时按快照把 used_count 回写，并发的兑换计数就会被抹掉，
// 兑换次数统计随之失真。PromoService.Update 也从不修改该字段。
func (r *promoCodeRepository) Update(ctx context.Context, code *service.PromoCode) error {
	client := clientFromContext(ctx, r.client)
	builder := client.PromoCode.UpdateOneID(code.ID).
		SetCode(code.Code).
		SetBonusAmount(code.BonusAmount).
		SetMaxUses(code.MaxUses).
		SetStatus(code.Status).
		SetNotes(code.Notes)

	if code.ExpiresAt != nil {
		builder.SetExpiresAt(*code.ExpiresAt)
placeholder else {
		builder.ClearExpiresAt()
placeholder

	updated, err := builder.Save(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return service.ErrPromoCodeNotFound
	placeholder
		return err
placeholder

	code.UpdatedAt = updated.UpdatedAt
	return nil
placeholder

func (r *promoCodeRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.PromoCode.Delete().Where(promocode.IDEQ(id)).Exec(ctx)
	return err
placeholder

func (r *promoCodeRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.PromoCode, *pagination.PaginationResult, error) {
	return r.ListWithFilters(ctx, params, "", "")
placeholder

func (r *promoCodeRepository) ListWithFilters(ctx context.Context, params pagination.PaginationParams, status, search string) ([]service.PromoCode, *pagination.PaginationResult, error) {
	q := r.client.PromoCode.Query()

	if status != "" {
		q = q.Where(promocode.StatusEQ(status))
placeholder
	if search != "" {
		q = q.Where(promocode.CodeContainsFold(search))
placeholder

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	codesQuery := q.
		Offset(params.Offset()).
		Limit(params.Limit())
	for _, order := range promoCodeListOrder(params) {
		codesQuery = codesQuery.Order(order)
placeholder

	codes, err := codesQuery.All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	outCodes := promoCodeEntitiesToService(codes)

	return outCodes, paginationResultFromTotal(int64(total), params), nil
placeholder

func promoCodeListOrder(params pagination.PaginationParams) []func(*entsql.Selector) {
	sortBy := strings.ToLower(strings.TrimSpace(params.SortBy))
	sortOrder := params.NormalizedSortOrder(pagination.SortOrderDesc)

	var field string
	switch sortBy {
	case "bonus_amount":
		field = promocode.FieldBonusAmount
	case "status":
		field = promocode.FieldStatus
	case "expires_at":
		field = promocode.FieldExpiresAt
	case "created_at":
		field = promocode.FieldCreatedAt
	case "code":
		field = promocode.FieldCode
	default:
		field = promocode.FieldID
placeholder

	if sortOrder == pagination.SortOrderAsc {
		return []func(*entsql.Selector){dbent.Asc(field), dbent.Asc(promocode.FieldID)placeholder
placeholder
	return []func(*entsql.Selector){dbent.Desc(field), dbent.Desc(promocode.FieldID)placeholder
placeholder

func (r *promoCodeRepository) CreateUsage(ctx context.Context, usage *service.PromoCodeUsage) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.PromoCodeUsage.Create().
		SetPromoCodeID(usage.PromoCodeID).
		SetUserID(usage.UserID).
		SetBonusAmount(usage.BonusAmount).
		SetUsedAt(usage.UsedAt).
		Save(ctx)
	if err != nil {
		return err
placeholder

	usage.ID = created.ID
	return nil
placeholder

func (r *promoCodeRepository) GetUsageByPromoCodeAndUser(ctx context.Context, promoCodeID, userID int64) (*service.PromoCodeUsage, error) {
	m, err := r.client.PromoCodeUsage.Query().
		Where(
			promocodeusage.PromoCodeIDEQ(promoCodeID),
			promocodeusage.UserIDEQ(userID),
		).
		Only(ctx)
	if err != nil {
		if dbent.IsNotFound(err) {
			return nil, nil
	placeholder
		return nil, err
placeholder
	return promoCodeUsageEntityToService(m), nil
placeholder

func (r *promoCodeRepository) ListUsagesByPromoCode(ctx context.Context, promoCodeID int64, params pagination.PaginationParams) ([]service.PromoCodeUsage, *pagination.PaginationResult, error) {
	q := r.client.PromoCodeUsage.Query().
		Where(promocodeusage.PromoCodeIDEQ(promoCodeID))

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	usages, err := q.
		WithUser().
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(promocodeusage.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	outUsages := promoCodeUsageEntitiesToService(usages)

	return outUsages, paginationResultFromTotal(int64(total), params), nil
placeholder

func (r *promoCodeRepository) IncrementUsedCount(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.PromoCode.UpdateOneID(id).
		AddUsedCount(1).
		Save(ctx)
	return err
placeholder

// Entity to Service conversions

func promoCodeEntityToService(m *dbent.PromoCode) *service.PromoCode {
	if m == nil {
		return nil
placeholder
	return &service.PromoCode{
		ID:          m.ID,
		Code:        m.Code,
		BonusAmount: m.BonusAmount,
		MaxUses:     m.MaxUses,
		UsedCount:   m.UsedCount,
		Status:      m.Status,
		ExpiresAt:   m.ExpiresAt,
		Notes:       derefString(m.Notes),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
placeholder
placeholder

func promoCodeEntitiesToService(models []*dbent.PromoCode) []service.PromoCode {
	out := make([]service.PromoCode, 0, len(models))
	for i := range models {
		if s := promoCodeEntityToService(models[i]); s != nil {
			out = append(out, *s)
	placeholder
placeholder
	return out
placeholder

func promoCodeUsageEntityToService(m *dbent.PromoCodeUsage) *service.PromoCodeUsage {
	if m == nil {
		return nil
placeholder
	out := &service.PromoCodeUsage{
		ID:          m.ID,
		PromoCodeID: m.PromoCodeID,
		UserID:      m.UserID,
		BonusAmount: m.BonusAmount,
		UsedAt:      m.UsedAt,
placeholder
	if m.Edges.User != nil {
		out.User = userEntityToService(m.Edges.User)
placeholder
	return out
placeholder

func promoCodeUsageEntitiesToService(models []*dbent.PromoCodeUsage) []service.PromoCodeUsage {
	out := make([]service.PromoCodeUsage, 0, len(models))
	for i := range models {
		if s := promoCodeUsageEntityToService(models[i]); s != nil {
			out = append(out, *s)
	placeholder
placeholder
	return out
placeholder
