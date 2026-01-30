package repository

import (
	"context"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/announcement"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

type announcementRepository struct {
	client *dbent.Client
placeholder

func NewAnnouncementRepository(client *dbent.Client) service.AnnouncementRepository {
	return &announcementRepository{client: clientplaceholder
placeholder

func (r *announcementRepository) Create(ctx context.Context, a *service.Announcement) error {
	client := clientFromContext(ctx, r.client)
	builder := client.Announcement.Create().
		SetTitle(a.Title).
		SetContent(a.Content).
		SetStatus(a.Status).
		SetTargeting(a.Targeting)

	if a.StartsAt != nil {
		builder.SetStartsAt(*a.StartsAt)
placeholder
	if a.EndsAt != nil {
		builder.SetEndsAt(*a.EndsAt)
placeholder
	if a.CreatedBy != nil {
		builder.SetCreatedBy(*a.CreatedBy)
placeholder
	if a.UpdatedBy != nil {
		builder.SetUpdatedBy(*a.UpdatedBy)
placeholder

	created, err := builder.Save(ctx)
	if err != nil {
		return err
placeholder

	applyAnnouncementEntityToService(a, created)
	return nil
placeholder

func (r *announcementRepository) GetByID(ctx context.Context, id int64) (*service.Announcement, error) {
	m, err := r.client.Announcement.Query().
		Where(announcement.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrAnnouncementNotFound, nil)
placeholder
	return announcementEntityToService(m), nil
placeholder

func (r *announcementRepository) Update(ctx context.Context, a *service.Announcement) error {
	client := clientFromContext(ctx, r.client)
	builder := client.Announcement.UpdateOneID(a.ID).
		SetTitle(a.Title).
		SetContent(a.Content).
		SetStatus(a.Status).
		SetTargeting(a.Targeting)

	if a.StartsAt != nil {
		builder.SetStartsAt(*a.StartsAt)
placeholder else {
		builder.ClearStartsAt()
placeholder
	if a.EndsAt != nil {
		builder.SetEndsAt(*a.EndsAt)
placeholder else {
		builder.ClearEndsAt()
placeholder
	if a.CreatedBy != nil {
		builder.SetCreatedBy(*a.CreatedBy)
placeholder else {
		builder.ClearCreatedBy()
placeholder
	if a.UpdatedBy != nil {
		builder.SetUpdatedBy(*a.UpdatedBy)
placeholder else {
		builder.ClearUpdatedBy()
placeholder

	updated, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, service.ErrAnnouncementNotFound, nil)
placeholder

	a.UpdatedAt = updated.UpdatedAt
	return nil
placeholder

func (r *announcementRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.Announcement.Delete().Where(announcement.IDEQ(id)).Exec(ctx)
	return err
placeholder

func (r *announcementRepository) List(
	ctx context.Context,
	params pagination.PaginationParams,
	filters service.AnnouncementListFilters,
) ([]service.Announcement, *pagination.PaginationResult, error) {
	q := r.client.Announcement.Query()

	if filters.Status != "" {
		q = q.Where(announcement.StatusEQ(filters.Status))
placeholder
	if filters.Search != "" {
		q = q.Where(
			announcement.Or(
				announcement.TitleContainsFold(filters.Search),
				announcement.ContentContainsFold(filters.Search),
			),
		)
placeholder

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	items, err := q.
		Offset(params.Offset()).
		Limit(params.Limit()).
		Order(dbent.Desc(announcement.FieldID)).
		All(ctx)
	if err != nil {
		return nil, nil, err
placeholder

	out := announcementEntitiesToService(items)
	return out, paginationResultFromTotal(int64(total), params), nil
placeholder

func (r *announcementRepository) ListActive(ctx context.Context, now time.Time) ([]service.Announcement, error) {
	q := r.client.Announcement.Query().
		Where(
			announcement.StatusEQ(service.AnnouncementStatusActive),
			announcement.Or(announcement.StartsAtIsNil(), announcement.StartsAtLTE(now)),
			announcement.Or(announcement.EndsAtIsNil(), announcement.EndsAtGT(now)),
		).
		Order(dbent.Desc(announcement.FieldID))

	items, err := q.All(ctx)
	if err != nil {
		return nil, err
placeholder
	return announcementEntitiesToService(items), nil
placeholder

func applyAnnouncementEntityToService(dst *service.Announcement, src *dbent.Announcement) {
	if dst == nil || src == nil {
		return
placeholder
	dst.ID = src.ID
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
placeholder

func announcementEntityToService(m *dbent.Announcement) *service.Announcement {
	if m == nil {
		return nil
placeholder
	return &service.Announcement{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		Status:    m.Status,
		Targeting: m.Targeting,
		StartsAt:  m.StartsAt,
		EndsAt:    m.EndsAt,
		CreatedBy: m.CreatedBy,
		UpdatedBy: m.UpdatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
placeholder
placeholder

func announcementEntitiesToService(models []*dbent.Announcement) []service.Announcement {
	out := make([]service.Announcement, 0, len(models))
	for i := range models {
		if s := announcementEntityToService(models[i]); s != nil {
			out = append(out, *s)
	placeholder
placeholder
	return out
placeholder

