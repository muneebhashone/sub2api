package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// UserAttributeService handles attribute management
type UserAttributeService struct {
	defRepo   UserAttributeDefinitionRepository
	valueRepo UserAttributeValueRepository
placeholder

// NewUserAttributeService creates a new service instance
func NewUserAttributeService(
	defRepo UserAttributeDefinitionRepository,
	valueRepo UserAttributeValueRepository,
) *UserAttributeService {
	return &UserAttributeService{
		defRepo:   defRepo,
		valueRepo: valueRepo,
placeholder
placeholder

// CreateDefinition creates a new attribute definition
func (s *UserAttributeService) CreateDefinition(ctx context.Context, input CreateAttributeDefinitionInput) (*UserAttributeDefinition, error) {
	// Validate type
	if !isValidAttributeType(input.Type) {
		return nil, ErrInvalidAttributeType
placeholder

	// Check if key exists
	exists, err := s.defRepo.ExistsByKey(ctx, input.Key)
	if err != nil {
		return nil, fmt.Errorf("check key exists: %w", err)
placeholder
	if exists {
		return nil, ErrAttributeKeyExists
placeholder

	def := &UserAttributeDefinition{
		Key:         input.Key,
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Options:     input.Options,
		Required:    input.Required,
		Validation:  input.Validation,
		Placeholder: input.Placeholder,
		Enabled:     input.Enabled,
placeholder

	if err := validateDefinitionPattern(def); err != nil {
		return nil, err
placeholder

	if err := s.defRepo.Create(ctx, def); err != nil {
		return nil, fmt.Errorf("create definition: %w", err)
placeholder

	return def, nil
placeholder

// GetDefinition retrieves a definition by ID
func (s *UserAttributeService) GetDefinition(ctx context.Context, id int64) (*UserAttributeDefinition, error) {
	return s.defRepo.GetByID(ctx, id)
placeholder

// ListDefinitions lists all definitions
func (s *UserAttributeService) ListDefinitions(ctx context.Context, enabledOnly bool) ([]UserAttributeDefinition, error) {
	return s.defRepo.List(ctx, enabledOnly)
placeholder

// UpdateDefinition updates an existing definition
func (s *UserAttributeService) UpdateDefinition(ctx context.Context, id int64, input UpdateAttributeDefinitionInput) (*UserAttributeDefinition, error) {
	def, err := s.defRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
placeholder

	if input.Name != nil {
		def.Name = *input.Name
placeholder
	if input.Description != nil {
		def.Description = *input.Description
placeholder
	if input.Type != nil {
		if !isValidAttributeType(*input.Type) {
			return nil, ErrInvalidAttributeType
	placeholder
		def.Type = *input.Type
placeholder
	if input.Options != nil {
		def.Options = *input.Options
placeholder
	if input.Required != nil {
		def.Required = *input.Required
placeholder
	if input.Validation != nil {
		def.Validation = *input.Validation
placeholder
	if input.Placeholder != nil {
		def.Placeholder = *input.Placeholder
placeholder
	if input.Enabled != nil {
		def.Enabled = *input.Enabled
placeholder

	if err := validateDefinitionPattern(def); err != nil {
		return nil, err
placeholder

	if err := s.defRepo.Update(ctx, def); err != nil {
		return nil, fmt.Errorf("update definition: %w", err)
placeholder

	return def, nil
placeholder

// DeleteDefinition soft-deletes a definition and hard-deletes associated values
func (s *UserAttributeService) DeleteDefinition(ctx context.Context, id int64) error {
	// Check if definition exists
	_, err := s.defRepo.GetByID(ctx, id)
	if err != nil {
		return err
placeholder

	// First delete all values (hard delete)
	if err := s.valueRepo.DeleteByAttributeID(ctx, id); err != nil {
		return fmt.Errorf("delete values: %w", err)
placeholder

	// Then soft-delete the definition
	if err := s.defRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete definition: %w", err)
placeholder

	return nil
placeholder

// ReorderDefinitions updates display order for multiple definitions
func (s *UserAttributeService) ReorderDefinitions(ctx context.Context, orders map[int64]int) error {
	return s.defRepo.UpdateDisplayOrders(ctx, orders)
placeholder

// GetUserAttributes retrieves all attribute values for a user
func (s *UserAttributeService) GetUserAttributes(ctx context.Context, userID int64) ([]UserAttributeValue, error) {
	return s.valueRepo.GetByUserID(ctx, userID)
placeholder

// GetBatchUserAttributes retrieves attribute values for multiple users
// Returns a map of userID -> map of attributeID -> value
func (s *UserAttributeService) GetBatchUserAttributes(ctx context.Context, userIDs []int64) (map[int64]map[int64]string, error) {
	values, err := s.valueRepo.GetByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
placeholder

	result := make(map[int64]map[int64]string)
	for _, v := range values {
		if result[v.UserID] == nil {
			result[v.UserID] = make(map[int64]string)
	placeholder
		result[v.UserID][v.AttributeID] = v.Value
placeholder

	return result, nil
placeholder

// UpdateUserAttributes batch updates attribute values for a user
func (s *UserAttributeService) UpdateUserAttributes(ctx context.Context, userID int64, inputs []UpdateUserAttributeInput) error {
	// Validate all values before updating
	defs, err := s.defRepo.List(ctx, true)
	if err != nil {
		return fmt.Errorf("list definitions: %w", err)
placeholder

	defMap := make(map[int64]*UserAttributeDefinition, len(defs))
	for i := range defs {
		defMap[defs[i].ID] = &defs[i]
placeholder

	for _, input := range inputs {
		def, ok := defMap[input.AttributeID]
		if !ok {
			return ErrAttributeDefinitionNotFound
	placeholder

		if err := s.validateValue(def, input.Value); err != nil {
			return err
	placeholder
placeholder

	return s.valueRepo.UpsertBatch(ctx, userID, inputs)
placeholder

// validateValue validates a value against its definition
func (s *UserAttributeService) validateValue(def *UserAttributeDefinition, value string) error {
	// Skip validation for empty non-required fields
	if value == "" && !def.Required {
		return nil
placeholder

	// Required check
	if def.Required && value == "" {
		return validationError(fmt.Sprintf("%s is required", def.Name))
placeholder

	v := def.Validation

	// String length validation
	if v.MinLength != nil && len(value) < *v.MinLength {
		return validationError(fmt.Sprintf("%s must be at least %d characters", def.Name, *v.MinLength))
placeholder
	if v.MaxLength != nil && len(value) > *v.MaxLength {
		return validationError(fmt.Sprintf("%s must be at most %d characters", def.Name, *v.MaxLength))
placeholder

	// Number validation
	if def.Type == AttributeTypeNumber && value != "" {
		num, err := strconv.Atoi(value)
		if err != nil {
			return validationError(fmt.Sprintf("%s must be a number", def.Name))
	placeholder
		if v.Min != nil && num < *v.Min {
			return validationError(fmt.Sprintf("%s must be at least %d", def.Name, *v.Min))
	placeholder
		if v.Max != nil && num > *v.Max {
			return validationError(fmt.Sprintf("%s must be at most %d", def.Name, *v.Max))
	placeholder
placeholder

	// Pattern validation
	if v.Pattern != nil && *v.Pattern != "" && value != "" {
		re, err := regexp.Compile(*v.Pattern)
		if err != nil {
			return validationError(def.Name + " has an invalid pattern")
	placeholder
		if !re.MatchString(value) {
			msg := def.Name + " format is invalid"
			if v.Message != nil && *v.Message != "" {
				msg = *v.Message
		placeholder
			return validationError(msg)
	placeholder
placeholder

	// Select validation
	if def.Type == AttributeTypeSelect && value != "" {
		found := false
		for _, opt := range def.Options {
			if opt.Value == value {
				found = true
				break
		placeholder
	placeholder
		if !found {
			return validationError(fmt.Sprintf("%s: invalid option", def.Name))
	placeholder
placeholder

	// Multi-select validation (stored as JSON array)
	if def.Type == AttributeTypeMultiSelect && value != "" {
		var values []string
		if err := json.Unmarshal([]byte(value), &values); err != nil {
			// Try comma-separated fallback
			values = strings.Split(value, ",")
	placeholder
		for _, val := range values {
			val = strings.TrimSpace(val)
			found := false
			for _, opt := range def.Options {
				if opt.Value == val {
					found = true
					break
			placeholder
		placeholder
			if !found {
				return validationError(fmt.Sprintf("%s: invalid option %s", def.Name, val))
		placeholder
	placeholder
placeholder

	return nil
placeholder

// validationError creates a validation error with a custom message
func validationError(msg string) error {
	return infraerrors.BadRequest("ATTRIBUTE_VALIDATION_FAILED", msg)
placeholder

func isValidAttributeType(t UserAttributeType) bool {
	switch t {
	case AttributeTypeText, AttributeTypeTextarea, AttributeTypeNumber,
		AttributeTypeEmail, AttributeTypeURL, AttributeTypeDate,
		AttributeTypeSelect, AttributeTypeMultiSelect:
		return true
placeholder
	return false
placeholder

func validateDefinitionPattern(def *UserAttributeDefinition) error {
	if def == nil {
		return nil
placeholder
	if def.Validation.Pattern == nil {
		return nil
placeholder
	pattern := strings.TrimSpace(*def.Validation.Pattern)
	if pattern == "" {
		return nil
placeholder
	if _, err := regexp.Compile(pattern); err != nil {
		return infraerrors.BadRequest("INVALID_ATTRIBUTE_PATTERN", fmt.Sprintf("invalid pattern for %s: %v", def.Name, err))
placeholder
	return nil
placeholder
