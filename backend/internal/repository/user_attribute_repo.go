package repository

import (
	"context"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/userattributedefinition"
	"github.com/Wei-Shaw/sub2api/ent/userattributevalue"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

// UserAttributeDefinitionRepository implementation
type userAttributeDefinitionRepository struct {
	client *dbent.Client
placeholder

// NewUserAttributeDefinitionRepository creates a new repository instance
func NewUserAttributeDefinitionRepository(client *dbent.Client) service.UserAttributeDefinitionRepository {
	return &userAttributeDefinitionRepository{client: clientplaceholder
placeholder

func (r *userAttributeDefinitionRepository) Create(ctx context.Context, def *service.UserAttributeDefinition) error {
	client := clientFromContext(ctx, r.client)

	created, err := client.UserAttributeDefinition.Create().
		SetKey(def.Key).
		SetName(def.Name).
		SetDescription(def.Description).
		SetType(string(def.Type)).
		SetOptions(toEntOptions(def.Options)).
		SetRequired(def.Required).
		SetValidation(toEntValidation(def.Validation)).
		SetPlaceholder(def.Placeholder).
		SetEnabled(def.Enabled).
		Save(ctx)

	if err != nil {
		return translatePersistenceError(err, nil, service.ErrAttributeKeyExists)
placeholder

	def.ID = created.ID
	def.DisplayOrder = created.DisplayOrder
	def.CreatedAt = created.CreatedAt
	def.UpdatedAt = created.UpdatedAt
	return nil
placeholder

func (r *userAttributeDefinitionRepository) GetByID(ctx context.Context, id int64) (*service.UserAttributeDefinition, error) {
	client := clientFromContext(ctx, r.client)

	e, err := client.UserAttributeDefinition.Query().
		Where(userattributedefinition.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrAttributeDefinitionNotFound, nil)
placeholder
	return defEntityToService(e), nil
placeholder

func (r *userAttributeDefinitionRepository) GetByKey(ctx context.Context, key string) (*service.UserAttributeDefinition, error) {
	client := clientFromContext(ctx, r.client)

	e, err := client.UserAttributeDefinition.Query().
		Where(userattributedefinition.KeyEQ(key)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrAttributeDefinitionNotFound, nil)
placeholder
	return defEntityToService(e), nil
placeholder

func (r *userAttributeDefinitionRepository) Update(ctx context.Context, def *service.UserAttributeDefinition) error {
	client := clientFromContext(ctx, r.client)

	updated, err := client.UserAttributeDefinition.UpdateOneID(def.ID).
		SetName(def.Name).
		SetDescription(def.Description).
		SetType(string(def.Type)).
		SetOptions(toEntOptions(def.Options)).
		SetRequired(def.Required).
		SetValidation(toEntValidation(def.Validation)).
		SetPlaceholder(def.Placeholder).
		SetDisplayOrder(def.DisplayOrder).
		SetEnabled(def.Enabled).
		Save(ctx)

	if err != nil {
		return translatePersistenceError(err, service.ErrAttributeDefinitionNotFound, service.ErrAttributeKeyExists)
placeholder

	def.UpdatedAt = updated.UpdatedAt
	return nil
placeholder

func (r *userAttributeDefinitionRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)

	_, err := client.UserAttributeDefinition.Delete().
		Where(userattributedefinition.IDEQ(id)).
		Exec(ctx)
	return translatePersistenceError(err, service.ErrAttributeDefinitionNotFound, nil)
placeholder

func (r *userAttributeDefinitionRepository) List(ctx context.Context, enabledOnly bool) ([]service.UserAttributeDefinition, error) {
	client := clientFromContext(ctx, r.client)

	q := client.UserAttributeDefinition.Query()
	if enabledOnly {
		q = q.Where(userattributedefinition.EnabledEQ(true))
placeholder

	entities, err := q.Order(dbent.Asc(userattributedefinition.FieldDisplayOrder)).All(ctx)
	if err != nil {
		return nil, err
placeholder

	result := make([]service.UserAttributeDefinition, 0, len(entities))
	for _, e := range entities {
		result = append(result, *defEntityToService(e))
placeholder
	return result, nil
placeholder

func (r *userAttributeDefinitionRepository) UpdateDisplayOrders(ctx context.Context, orders map[int64]int) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
placeholder
	defer func() { _ = tx.Rollback() placeholder()

	for id, order := range orders {
		if _, err := tx.UserAttributeDefinition.UpdateOneID(id).
			SetDisplayOrder(order).
			Save(ctx); err != nil {
			return translatePersistenceError(err, service.ErrAttributeDefinitionNotFound, nil)
	placeholder
placeholder

	return tx.Commit()
placeholder

func (r *userAttributeDefinitionRepository) ExistsByKey(ctx context.Context, key string) (bool, error) {
	client := clientFromContext(ctx, r.client)
	return client.UserAttributeDefinition.Query().
		Where(userattributedefinition.KeyEQ(key)).
		Exist(ctx)
placeholder

// UserAttributeValueRepository implementation
type userAttributeValueRepository struct {
	client *dbent.Client
placeholder

// NewUserAttributeValueRepository creates a new repository instance
func NewUserAttributeValueRepository(client *dbent.Client) service.UserAttributeValueRepository {
	return &userAttributeValueRepository{client: clientplaceholder
placeholder

func (r *userAttributeValueRepository) GetByUserID(ctx context.Context, userID int64) ([]service.UserAttributeValue, error) {
	client := clientFromContext(ctx, r.client)

	entities, err := client.UserAttributeValue.Query().
		Where(userattributevalue.UserIDEQ(userID)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	result := make([]service.UserAttributeValue, 0, len(entities))
	for _, e := range entities {
		result = append(result, service.UserAttributeValue{
			ID:          e.ID,
			UserID:      e.UserID,
			AttributeID: e.AttributeID,
			Value:       e.Value,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
	placeholder)
placeholder
	return result, nil
placeholder

func (r *userAttributeValueRepository) GetByUserIDs(ctx context.Context, userIDs []int64) ([]service.UserAttributeValue, error) {
	if len(userIDs) == 0 {
		return []service.UserAttributeValue{placeholder, nil
placeholder

	client := clientFromContext(ctx, r.client)

	entities, err := client.UserAttributeValue.Query().
		Where(userattributevalue.UserIDIn(userIDs...)).
		All(ctx)
	if err != nil {
		return nil, err
placeholder

	result := make([]service.UserAttributeValue, 0, len(entities))
	for _, e := range entities {
		result = append(result, service.UserAttributeValue{
			ID:          e.ID,
			UserID:      e.UserID,
			AttributeID: e.AttributeID,
			Value:       e.Value,
			CreatedAt:   e.CreatedAt,
			UpdatedAt:   e.UpdatedAt,
	placeholder)
placeholder
	return result, nil
placeholder

func (r *userAttributeValueRepository) UpsertBatch(ctx context.Context, userID int64, inputs []service.UpdateUserAttributeInput) error {
	if len(inputs) == 0 {
		return nil
placeholder

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
placeholder
	defer func() { _ = tx.Rollback() placeholder()

	for _, input := range inputs {
		// Use upsert (ON CONFLICT DO UPDATE)
		err := tx.UserAttributeValue.Create().
			SetUserID(userID).
			SetAttributeID(input.AttributeID).
			SetValue(input.Value).
			OnConflictColumns(userattributevalue.FieldUserID, userattributevalue.FieldAttributeID).
			UpdateValue().
			UpdateUpdatedAt().
			Exec(ctx)
		if err != nil {
			return err
	placeholder
placeholder

	return tx.Commit()
placeholder

func (r *userAttributeValueRepository) DeleteByAttributeID(ctx context.Context, attributeID int64) error {
	client := clientFromContext(ctx, r.client)

	_, err := client.UserAttributeValue.Delete().
		Where(userattributevalue.AttributeIDEQ(attributeID)).
		Exec(ctx)
	return err
placeholder

func (r *userAttributeValueRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	client := clientFromContext(ctx, r.client)

	_, err := client.UserAttributeValue.Delete().
		Where(userattributevalue.UserIDEQ(userID)).
		Exec(ctx)
	return err
placeholder

// Helper functions for entity to service conversion
func defEntityToService(e *dbent.UserAttributeDefinition) *service.UserAttributeDefinition {
	if e == nil {
		return nil
placeholder
	return &service.UserAttributeDefinition{
		ID:           e.ID,
		Key:          e.Key,
		Name:         e.Name,
		Description:  e.Description,
		Type:         service.UserAttributeType(e.Type),
		Options:      toServiceOptions(e.Options),
		Required:     e.Required,
		Validation:   toServiceValidation(e.Validation),
		Placeholder:  e.Placeholder,
		DisplayOrder: e.DisplayOrder,
		Enabled:      e.Enabled,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
placeholder
placeholder

// Type conversion helpers (map types <-> service types)
func toEntOptions(opts []service.UserAttributeOption) []map[string]any {
	if opts == nil {
		return []map[string]any{placeholder
placeholder
	result := make([]map[string]any, len(opts))
	for i, o := range opts {
		result[i] = map[string]any{"value": o.Value, "label": o.Labelplaceholder
placeholder
	return result
placeholder

func toServiceOptions(opts []map[string]any) []service.UserAttributeOption {
	if opts == nil {
		return []service.UserAttributeOption{placeholder
placeholder
	result := make([]service.UserAttributeOption, len(opts))
	for i, o := range opts {
		result[i] = service.UserAttributeOption{
			Value: getString(o, "value"),
			Label: getString(o, "label"),
	placeholder
placeholder
	return result
placeholder

func toEntValidation(v service.UserAttributeValidation) map[string]any {
	result := map[string]any{placeholder
	if v.MinLength != nil {
		result["min_length"] = *v.MinLength
placeholder
	if v.MaxLength != nil {
		result["max_length"] = *v.MaxLength
placeholder
	if v.Min != nil {
		result["min"] = *v.Min
placeholder
	if v.Max != nil {
		result["max"] = *v.Max
placeholder
	if v.Pattern != nil {
		result["pattern"] = *v.Pattern
placeholder
	if v.Message != nil {
		result["message"] = *v.Message
placeholder
	return result
placeholder

func toServiceValidation(v map[string]any) service.UserAttributeValidation {
	result := service.UserAttributeValidation{placeholder
	if val := getInt(v, "min_length"); val != nil {
		result.MinLength = val
placeholder
	if val := getInt(v, "max_length"); val != nil {
		result.MaxLength = val
placeholder
	if val := getInt(v, "min"); val != nil {
		result.Min = val
placeholder
	if val := getInt(v, "max"); val != nil {
		result.Max = val
placeholder
	if val := getStringPtr(v, "pattern"); val != nil {
		result.Pattern = val
placeholder
	if val := getStringPtr(v, "message"); val != nil {
		result.Message = val
placeholder
	return result
placeholder

// Helper functions for type conversion
func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
	placeholder
placeholder
	return ""
placeholder

func getStringPtr(m map[string]any, key string) *string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return &s
	placeholder
placeholder
	return nil
placeholder

func getInt(m map[string]any, key string) *int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case int:
			return &n
		case int64:
			i := int(n)
			return &i
		case float64:
			i := int(n)
			return &i
	placeholder
placeholder
	return nil
placeholder
