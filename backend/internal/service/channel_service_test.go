//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mock: ChannelRepository
// ---------------------------------------------------------------------------

type mockChannelRepository struct {
	listAllFn                  func(ctx context.Context) ([]Channel, error)
	getGroupPlatformsFn        func(ctx context.Context, groupIDs []int64) (map[int64]string, error)
	createFn                   func(ctx context.Context, channel *Channel) error
	getByIDFn                  func(ctx context.Context, id int64) (*Channel, error)
	updateFn                   func(ctx context.Context, channel *Channel) error
	deleteFn                   func(ctx context.Context, id int64) error
	listFn                     func(ctx context.Context, params pagination.PaginationParams, status, search string) ([]Channel, *pagination.PaginationResult, error)
	existsByNameFn             func(ctx context.Context, name string) (bool, error)
	existsByNameExcludingFn    func(ctx context.Context, name string, excludeID int64) (bool, error)
	getGroupIDsFn              func(ctx context.Context, channelID int64) ([]int64, error)
	setGroupIDsFn              func(ctx context.Context, channelID int64, groupIDs []int64) error
	getChannelIDByGroupIDFn    func(ctx context.Context, groupID int64) (int64, error)
	getGroupsInOtherChannelsFn func(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error)
	listModelPricingFn         func(ctx context.Context, channelID int64) ([]ChannelModelPricing, error)
	createModelPricingFn       func(ctx context.Context, pricing *ChannelModelPricing) error
	updateModelPricingFn       func(ctx context.Context, pricing *ChannelModelPricing) error
	deleteModelPricingFn       func(ctx context.Context, id int64) error
	replaceModelPricingFn      func(ctx context.Context, channelID int64, pricingList []ChannelModelPricing) error
placeholder

func (m *mockChannelRepository) Create(ctx context.Context, channel *Channel) error {
	if m.createFn != nil {
		return m.createFn(ctx, channel)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) GetByID(ctx context.Context, id int64) (*Channel, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
placeholder
	return nil, ErrChannelNotFound
placeholder

func (m *mockChannelRepository) Update(ctx context.Context, channel *Channel) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, channel)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) List(ctx context.Context, params pagination.PaginationParams, status, search string) ([]Channel, *pagination.PaginationResult, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params, status, search)
placeholder
	return nil, nil, nil
placeholder

func (m *mockChannelRepository) ListAll(ctx context.Context) ([]Channel, error) {
	if m.listAllFn != nil {
		return m.listAllFn(ctx)
placeholder
	return nil, nil
placeholder

func (m *mockChannelRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if m.existsByNameFn != nil {
		return m.existsByNameFn(ctx, name)
placeholder
	return false, nil
placeholder

func (m *mockChannelRepository) ExistsByNameExcluding(ctx context.Context, name string, excludeID int64) (bool, error) {
	if m.existsByNameExcludingFn != nil {
		return m.existsByNameExcludingFn(ctx, name, excludeID)
placeholder
	return false, nil
placeholder

func (m *mockChannelRepository) GetGroupIDs(ctx context.Context, channelID int64) ([]int64, error) {
	if m.getGroupIDsFn != nil {
		return m.getGroupIDsFn(ctx, channelID)
placeholder
	return nil, nil
placeholder

func (m *mockChannelRepository) SetGroupIDs(ctx context.Context, channelID int64, groupIDs []int64) error {
	if m.setGroupIDsFn != nil {
		return m.setGroupIDsFn(ctx, channelID, groupIDs)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) GetChannelIDByGroupID(ctx context.Context, groupID int64) (int64, error) {
	if m.getChannelIDByGroupIDFn != nil {
		return m.getChannelIDByGroupIDFn(ctx, groupID)
placeholder
	return 0, nil
placeholder

func (m *mockChannelRepository) GetGroupsInOtherChannels(ctx context.Context, channelID int64, groupIDs []int64) ([]int64, error) {
	if m.getGroupsInOtherChannelsFn != nil {
		return m.getGroupsInOtherChannelsFn(ctx, channelID, groupIDs)
placeholder
	return nil, nil
placeholder

func (m *mockChannelRepository) GetGroupPlatforms(ctx context.Context, groupIDs []int64) (map[int64]string, error) {
	if m.getGroupPlatformsFn != nil {
		return m.getGroupPlatformsFn(ctx, groupIDs)
placeholder
	return nil, nil
placeholder

func (m *mockChannelRepository) ListModelPricing(ctx context.Context, channelID int64) ([]ChannelModelPricing, error) {
	if m.listModelPricingFn != nil {
		return m.listModelPricingFn(ctx, channelID)
placeholder
	return nil, nil
placeholder

func (m *mockChannelRepository) CreateModelPricing(ctx context.Context, pricing *ChannelModelPricing) error {
	if m.createModelPricingFn != nil {
		return m.createModelPricingFn(ctx, pricing)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) UpdateModelPricing(ctx context.Context, pricing *ChannelModelPricing) error {
	if m.updateModelPricingFn != nil {
		return m.updateModelPricingFn(ctx, pricing)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) DeleteModelPricing(ctx context.Context, id int64) error {
	if m.deleteModelPricingFn != nil {
		return m.deleteModelPricingFn(ctx, id)
placeholder
	return nil
placeholder

func (m *mockChannelRepository) ReplaceModelPricing(ctx context.Context, channelID int64, pricingList []ChannelModelPricing) error {
	if m.replaceModelPricingFn != nil {
		return m.replaceModelPricingFn(ctx, channelID, pricingList)
placeholder
	return nil
placeholder

// ---------------------------------------------------------------------------
// Mock: APIKeyAuthCacheInvalidator
// ---------------------------------------------------------------------------

type mockChannelAuthCacheInvalidator struct {
	invalidatedGroupIDs []int64
	invalidatedKeys     []string
	invalidatedUserIDs  []int64
placeholder

func (m *mockChannelAuthCacheInvalidator) InvalidateAuthCacheByKey(_ context.Context, key string) {
	m.invalidatedKeys = append(m.invalidatedKeys, key)
placeholder

func (m *mockChannelAuthCacheInvalidator) InvalidateAuthCacheByUserID(_ context.Context, userID int64) {
	m.invalidatedUserIDs = append(m.invalidatedUserIDs, userID)
placeholder

func (m *mockChannelAuthCacheInvalidator) InvalidateAuthCacheByGroupID(_ context.Context, groupID int64) {
	m.invalidatedGroupIDs = append(m.invalidatedGroupIDs, groupID)
placeholder

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestChannelService(repo *mockChannelRepository) *ChannelService {
	return NewChannelService(repo, nil)
placeholder

func newTestChannelServiceWithAuth(repo *mockChannelRepository, auth *mockChannelAuthCacheInvalidator) *ChannelService {
	return NewChannelService(repo, auth)
placeholder

// makeStandardRepo returns a repo that serves one active channel with anthropic pricing
// for group 1, with the given model pricing and model mapping.
func makeStandardRepo(ch Channel, groupPlatforms map[int64]string) *mockChannelRepository {
	return &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return []Channel{chplaceholder, nil
	placeholder,
		getGroupPlatformsFn: func(_ context.Context, _ []int64) (map[int64]string, error) {
			return groupPlatforms, nil
	placeholder,
placeholder
placeholder

// ===========================================================================
// 1. BuildModelMappingChain
// ===========================================================================

func TestBuildModelMappingChain(t *testing.T) {
	tests := []struct {
		name          string
		result        ChannelMappingResult
		requestModel  string
		upstreamModel string
		want          string
placeholder{
		{
			name:          "no mapping, no upstream diff",
			result:        ChannelMappingResult{Mapped: false, MappedModel: "claude-sonnet-4"placeholder,
			requestModel:  "claude-sonnet-4",
			upstreamModel: "claude-sonnet-4",
			want:          "",
	placeholder,
		{
			name:          "no mapping, upstream differs",
			result:        ChannelMappingResult{Mapped: false, MappedModel: "claude-sonnet-4"placeholder,
			requestModel:  "claude-sonnet-4",
			upstreamModel: "claude-sonnet-4-20250514",
			want:          "claude-sonnet-4\u2192claude-sonnet-4-20250514",
	placeholder,
		{
			name:          "mapped, upstream differs",
			result:        ChannelMappingResult{Mapped: true, MappedModel: "claude-sonnet-4-20250514"placeholder,
			requestModel:  "my-model",
			upstreamModel: "actual-upstream",
			want:          "my-model\u2192claude-sonnet-4-20250514\u2192actual-upstream",
	placeholder,
		{
			name:          "mapped, upstream same as mapped",
			result:        ChannelMappingResult{Mapped: true, MappedModel: "claude-sonnet-4-20250514"placeholder,
			requestModel:  "claude-sonnet-4",
			upstreamModel: "claude-sonnet-4-20250514",
			want:          "claude-sonnet-4\u2192claude-sonnet-4-20250514",
	placeholder,
		{
			name:          "mapped, upstream empty",
			result:        ChannelMappingResult{Mapped: true, MappedModel: "target-model"placeholder,
			requestModel:  "my-model",
			upstreamModel: "",
			want:          "my-model\u2192target-model",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.BuildModelMappingChain(tt.requestModel, tt.upstreamModel)
			require.Equal(t, tt.want, got)
	placeholder)
placeholder
placeholder

// ===========================================================================
// 2. ReplaceModelInBody
// ===========================================================================

func TestReplaceModelInBody(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		newModel string
		check    func(t *testing.T, result []byte)
placeholder{
		{
			name:     "empty body",
			body:     []byte{placeholder,
			newModel: "new-model",
			check: func(t *testing.T, result []byte) {
				require.Equal(t, []byte{placeholder, result)
		placeholder,
	placeholder,
		{
			name:     "model already equal",
			body:     []byte(`{"model":"claude-sonnet-4","temperature":0.7placeholder`),
			newModel: "claude-sonnet-4",
			check: func(t *testing.T, result []byte) {
				require.Equal(t, []byte(`{"model":"claude-sonnet-4","temperature":0.7placeholder`), result)
		placeholder,
	placeholder,
		{
			name:     "model different",
			body:     []byte(`{"model":"claude-sonnet-4","temperature":0.7placeholder`),
			newModel: "claude-opus-4",
			check: func(t *testing.T, result []byte) {
				require.Contains(t, string(result), `"model":"claude-opus-4"`)
				require.Contains(t, string(result), `"temperature"`)
		placeholder,
	placeholder,
		{
			name:     "no model field",
			body:     []byte(`{"temperature":0.7placeholder`),
			newModel: "claude-opus-4",
			check: func(t *testing.T, result []byte) {
				require.Contains(t, string(result), `"model":"claude-opus-4"`)
				require.Contains(t, string(result), `"temperature"`)
		placeholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceModelInBody(tt.body, tt.newModel)
			tt.check(t, result)
	placeholder)
placeholder
placeholder

// ===========================================================================
// 3. validateNoConflictingModels + validateNoConflictingMappings
// ===========================================================================

func TestValidateNoConflictingModels(t *testing.T) {
	tests := []struct {
		name        string
		pricingList []ChannelModelPricing
		wantErr     bool
		errContains string
placeholder{
		{
			name: "no duplicates",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"claude-sonnet-4", "claude-opus-4"placeholderplaceholder,
				{Platform: "openai", Models: []string{"gpt-5.1"placeholderplaceholder,
		placeholder,
			wantErr: false,
	placeholder,
		{
			name: "same platform duplicate",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"claude-sonnet-4"placeholderplaceholder,
				{Platform: "anthropic", Models: []string{"claude-sonnet-4"placeholderplaceholder,
		placeholder,
			wantErr:     true,
			errContains: "claude-sonnet-4",
	placeholder,
		{
			name: "same model different platform",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"model-a"placeholderplaceholder,
				{Platform: "openai", Models: []string{"model-a"placeholderplaceholder,
		placeholder,
			wantErr: false,
	placeholder,
		{
			name: "case insensitive",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"Claude"placeholderplaceholder,
				{Platform: "anthropic", Models: []string{"claude"placeholderplaceholder,
		placeholder,
			wantErr: true,
	placeholder,
		{
			name:        "empty list (nil)",
			pricingList: nil,
			wantErr:     false,
	placeholder,
		{
			name: "wildcard_vs_wildcard_conflict",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"claude-*"placeholderplaceholder,
				{Platform: "anthropic", Models: []string{"claude-opus-*"placeholderplaceholder,
		placeholder,
			wantErr:     true,
			errContains: "conflict",
	placeholder,
		{
			name: "wildcard_vs_exact_conflict",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"claude-*"placeholderplaceholder,
				{Platform: "anthropic", Models: []string{"claude-opus-4-6"placeholderplaceholder,
		placeholder,
			wantErr:     true,
			errContains: "conflict",
	placeholder,
		{
			name: "no_conflict_different_platform",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"claude-opus-*"placeholderplaceholder,
				{Platform: "openai", Models: []string{"claude-*"placeholderplaceholder,
		placeholder,
			wantErr: false,
	placeholder,
		{
			name: "no_conflict_same_platform_different_prefix",
			pricingList: []ChannelModelPricing{
				{Platform: "anthropic", Models: []string{"claude-opus-*"placeholderplaceholder,
				{Platform: "anthropic", Models: []string{"gpt-*"placeholderplaceholder,
		placeholder,
			wantErr: false,
	placeholder,
		{
			name: "catch_all_wildcard_conflicts_with_everything",
			pricingList: []ChannelModelPricing{
				{Platform: "openai", Models: []string{"*"placeholderplaceholder,
				{Platform: "openai", Models: []string{"gpt-5"placeholderplaceholder,
		placeholder,
			wantErr:     true,
			errContains: "conflict",
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNoConflictingModels(tt.pricingList)
			if tt.wantErr {
			placeholder
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
			placeholder
		placeholder else {
			placeholder
		placeholder
	placeholder)
placeholder

	// Additional sub-case: explicit empty slice
	t.Run("empty list (empty slice)", func(t *testing.T) {
		err := validateNoConflictingModels([]ChannelModelPricing{placeholder)
	placeholder
placeholder)
placeholder

func TestValidateNoConflictingMappings(t *testing.T) {
	tests := []struct {
		name        string
		mapping     map[string]map[string]string
		wantErr     bool
		errContains string
placeholder{
		{
			name:    "nil mapping",
			mapping: nil,
			wantErr: false,
	placeholder,
		{
			name:    "empty mapping",
			mapping: map[string]map[string]string{placeholder,
			wantErr: false,
	placeholder,
		{
			name: "no conflict",
			mapping: map[string]map[string]string{
				"anthropic": {"claude-opus-*": "opus", "gpt-*": "gpt"placeholder,
		placeholder,
			wantErr: false,
	placeholder,
		{
			name: "wildcard vs wildcard conflict",
			mapping: map[string]map[string]string{
				"anthropic": {"claude-*": "a", "claude-opus-*": "b"placeholder,
		placeholder,
			wantErr:     true,
			errContains: "conflict",
	placeholder,
		{
			name: "wildcard vs exact conflict",
			mapping: map[string]map[string]string{
				"openai": {"gpt-*": "a", "gpt-4o": "b"placeholder,
		placeholder,
			wantErr:     true,
			errContains: "conflict",
	placeholder,
		{
			name: "exact duplicate conflict",
			mapping: map[string]map[string]string{
				"anthropic": {"claude-opus-4": "a"placeholder,
				"openai":    {"claude-opus-4": "b"placeholder,
		placeholder,
			wantErr: false, // different platforms
	placeholder,
		{
			name: "different platforms no conflict",
			mapping: map[string]map[string]string{
				"anthropic": {"claude-*": "a"placeholder,
				"openai":    {"claude-*": "b"placeholder,
		placeholder,
			wantErr: false,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNoConflictingMappings(tt.mapping)
			if tt.wantErr {
			placeholder
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
			placeholder
		placeholder else {
			placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestConflictsBetween(t *testing.T) {
	tests := []struct {
		name string
		a, b modelEntry
		want bool
placeholder{
		{
			name: "exact same",
			a:    modelEntry{prefix: "claude-opus-4", wildcard: falseplaceholder,
			b:    modelEntry{prefix: "claude-opus-4", wildcard: falseplaceholder,
			want: true,
	placeholder,
		{
			name: "exact different",
			a:    modelEntry{prefix: "claude-opus-4", wildcard: falseplaceholder,
			b:    modelEntry{prefix: "gpt-4o", wildcard: falseplaceholder,
			want: false,
	placeholder,
		{
			name: "wildcard matches exact",
			a:    modelEntry{prefix: "claude-", wildcard: trueplaceholder,
			b:    modelEntry{prefix: "claude-opus-4", wildcard: falseplaceholder,
			want: true,
	placeholder,
		{
			name: "exact does not match unrelated wildcard",
			a:    modelEntry{prefix: "gpt-4o", wildcard: falseplaceholder,
			b:    modelEntry{prefix: "claude-", wildcard: trueplaceholder,
			want: false,
	placeholder,
		{
			name: "wildcard prefix overlap",
			a:    modelEntry{prefix: "claude-", wildcard: trueplaceholder,
			b:    modelEntry{prefix: "claude-opus-", wildcard: trueplaceholder,
			want: true,
	placeholder,
		{
			name: "wildcards no overlap",
			a:    modelEntry{prefix: "claude-", wildcard: trueplaceholder,
			b:    modelEntry{prefix: "gpt-", wildcard: trueplaceholder,
			want: false,
	placeholder,
		{
			name: "catch-all wildcard vs any",
			a:    modelEntry{prefix: "", wildcard: trueplaceholder,
			b:    modelEntry{prefix: "anything", wildcard: falseplaceholder,
			want: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, conflictsBetween(tt.a, tt.b))
	placeholder)
placeholder
placeholder

// ===========================================================================
// 4. Cache Building + Hot Path Methods
// ===========================================================================

// --- 4.1 GetChannelForGroup ---

func TestGetChannelForGroup_Success(t *testing.T) {
	ch := Channel{
		ID:       1,
		Name:     "test-channel",
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.ID)
	require.Equal(t, "test-channel", result.Name)

	// returned value should be a clone
	result.Name = "mutated"
	result2, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.Equal(t, "test-channel", result2.Name)
placeholder

func TestGetChannelForGroup_InactiveChannel(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusDisabled,
		GroupIDs: []int64{10placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.Nil(t, result)
placeholder

func TestGetChannelForGroup_NoChannel(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result, err := svc.GetChannelForGroup(context.Background(), 999)
placeholder
	require.Nil(t, result)
placeholder

func TestGetChannelForGroup_CacheError(t *testing.T) {
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, errors.New("db connection failed")
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	result, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.Nil(t, result)
	require.Contains(t, err.Error(), "db connection failed")
placeholder

// --- 4.2 GetChannelModelPricing ---

func TestGetChannelModelPricing_ExactMatch(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result)
	require.Equal(t, int64(100), result.ID)
	require.InDelta(t, 15e-6, *result.InputPrice, 1e-12)
placeholder

func TestGetChannelModelPricing_CaseInsensitive(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "Claude-Opus-4")
	require.NotNil(t, result)
	require.Equal(t, int64(100), result.ID)
placeholder

func TestGetChannelModelPricing_WildcardMatch(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 200, Platform: "anthropic", Models: []string{"claude-*"placeholder, InputPrice: testPtrFloat64(10e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-sonnet-4")
	require.NotNil(t, result)
	require.Equal(t, int64(200), result.ID)
placeholder

func TestGetChannelModelPricing_WildcardFirstMatch(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 200, Platform: "anthropic", Models: []string{"claude-*"placeholder, InputPrice: testPtrFloat64(10e-6)placeholder,
			{ID: 300, Platform: "anthropic", Models: []string{"claude-sonnet-*"placeholder, InputPrice: testPtrFloat64(5e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-sonnet-4-20250514")
	require.NotNil(t, result)
	// "claude-*" is defined first, so it matches first regardless of prefix length
	require.Equal(t, int64(200), result.ID)
	require.InDelta(t, 10e-6, *result.InputPrice, 1e-12)
placeholder

func TestGetChannelModelPricing_NoMatch(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "gpt-5.1")
	require.Nil(t, result)
placeholder

func TestGetChannelModelPricing_InactiveChannel(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusDisabled,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.Nil(t, result)
placeholder

func TestGetChannelModelPricing_PlatformFiltering(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10, 20placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "openai", Models: []string{"gpt-5.1"placeholder, InputPrice: testPtrFloat64(5e-6)placeholder,
			{ID: 200, Platform: "anthropic", Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic", 20: "openai"placeholder)
	svc := newTestChannelService(repo)

	// Group 10 (anthropic) should NOT see openai pricing
	result := svc.GetChannelModelPricing(context.Background(), 10, "gpt-5.1")
	require.Nil(t, result)

	// Group 10 (anthropic) should see anthropic pricing
	result = svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result)
	require.Equal(t, int64(200), result.ID)

	// Group 20 (openai) should see openai pricing
	result = svc.GetChannelModelPricing(context.Background(), 20, "gpt-5.1")
	require.NotNil(t, result)
	require.Equal(t, int64(100), result.ID)

	// Group 20 (openai) should NOT see anthropic pricing
	result = svc.GetChannelModelPricing(context.Background(), 20, "claude-opus-4")
	require.Nil(t, result)
placeholder

func TestGetChannelModelPricing_ReturnsCopy(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result)

	// Mutate the returned pricing's slice fields — original cache should not be affected
	// (Clone copies slices independently, pointer fields are shared per design)
	result.Models = append(result.Models, "hacked")
	result.ID = 999

	// Original cache should not be affected (slice independence + struct copy)
	result2 := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result2)
	require.Equal(t, 1, len(result2.Models))
	require.Equal(t, int64(100), result2.ID)
placeholder

// --- 4.3 ResolveChannelMapping ---

func TestResolveChannelMapping_NoChannel(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	// Group 999 is not in any channel
	result := svc.ResolveChannelMapping(context.Background(), 999, "claude-opus-4")
	require.Equal(t, "claude-opus-4", result.MappedModel)
	require.False(t, result.Mapped)
	require.Equal(t, int64(0), result.ChannelID)
placeholder

func TestResolveChannelMapping_ExactMapping(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"claude-sonnet-4": "claude-sonnet-4-20250514",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-sonnet-4")
	require.True(t, result.Mapped)
	require.Equal(t, "claude-sonnet-4-20250514", result.MappedModel)
	require.Equal(t, int64(1), result.ChannelID)
placeholder

func TestResolveChannelMapping_WildcardMapping(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"*": "gpt-5.4",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "any-model-name")
	require.True(t, result.Mapped)
	require.Equal(t, "gpt-5.4", result.MappedModel)
placeholder

func TestResolveChannelMapping_WildcardFirstMatch(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"claude-*":        "target2",
				"claude-sonnet-*": "target1",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-sonnet-4")
	require.True(t, result.Mapped)
	// map iteration order is non-deterministic, so the first-match depends on
	// insertion order which Go maps don't guarantee; verify that one of the
	// wildcard targets matched
	require.Contains(t, []string{"target1", "target2"placeholder, result.MappedModel)
placeholder

func TestResolveChannelMapping_NoMapping(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"claude-sonnet-4": "mapped",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-opus-4")
	require.False(t, result.Mapped)
	require.Equal(t, "claude-opus-4", result.MappedModel)
	require.Equal(t, int64(1), result.ChannelID)
placeholder

func TestResolveChannelMapping_DefaultBillingModelSource(t *testing.T) {
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		BillingModelSource: "", // empty
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-opus-4")
	require.Equal(t, BillingModelSourceChannelMapped, result.BillingModelSource)
placeholder

func TestResolveChannelMapping_UpstreamBillingModelSource(t *testing.T) {
	ch := Channel{
		ID:                 1,
		Status:             StatusActive,
		GroupIDs:           []int64{10placeholder,
		BillingModelSource: BillingModelSourceUpstream,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-opus-4")
	require.Equal(t, BillingModelSourceUpstream, result.BillingModelSource)
placeholder

func TestResolveChannelMapping_InactiveChannel(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusDisabled,
		GroupIDs: []int64{10placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"claude-sonnet-4": "mapped",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-sonnet-4")
	require.False(t, result.Mapped)
	require.Equal(t, "claude-sonnet-4", result.MappedModel)
	require.Equal(t, int64(0), result.ChannelID) // no channel
placeholder

// --- 4.4 IsModelRestricted ---

func TestIsModelRestricted_NoChannel(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	// Group 999 is not in any channel
	restricted := svc.IsModelRestricted(context.Background(), 999, "claude-opus-4")
	require.False(t, restricted)
placeholder

func TestIsModelRestricted_RestrictDisabled(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: false,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	// Even though model is not in pricing, RestrictModels=false
	restricted := svc.IsModelRestricted(context.Background(), 10, "nonexistent-model")
	require.False(t, restricted)
placeholder

func TestIsModelRestricted_InactiveChannel(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusDisabled,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	restricted := svc.IsModelRestricted(context.Background(), 10, "any-model")
	require.False(t, restricted)
placeholder

func TestIsModelRestricted_ModelInPricing(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4", "claude-sonnet-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	restricted := svc.IsModelRestricted(context.Background(), 10, "claude-opus-4")
	require.False(t, restricted)
placeholder

func TestIsModelRestricted_ModelInWildcard(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-*"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	restricted := svc.IsModelRestricted(context.Background(), 10, "claude-sonnet-4")
	require.False(t, restricted)
placeholder

func TestIsModelRestricted_ModelNotFound(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	restricted := svc.IsModelRestricted(context.Background(), 10, "gpt-5.1")
	require.True(t, restricted)
placeholder

func TestIsModelRestricted_CaseInsensitive(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	restricted := svc.IsModelRestricted(context.Background(), 10, "Claude-Opus-4")
	require.False(t, restricted)
placeholder

// --- 4.5 ResolveChannelMappingAndRestrict ---

func TestResolveChannelMappingAndRestrict_NilGroupID(t *testing.T) {
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	mapping, restricted := svc.ResolveChannelMappingAndRestrict(context.Background(), nil, "claude-opus-4")
	require.False(t, restricted)
	require.False(t, mapping.Mapped)
	require.Equal(t, "claude-opus-4", mapping.MappedModel)
placeholder

func TestResolveChannelMappingAndRestrict_ModelInPricing_WithMapping(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4"placeholderplaceholder,
	placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"claude-sonnet-4": "claude-sonnet-4-20250514",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	gid := int64(10)
	mapping, restricted := svc.ResolveChannelMappingAndRestrict(context.Background(), &gid, "claude-sonnet-4")
	require.False(t, restricted) // model IS in pricing
	require.True(t, mapping.Mapped)
	require.Equal(t, "claude-sonnet-4-20250514", mapping.MappedModel)
placeholder

func TestResolveChannelMappingAndRestrict_ModelNotInPricing_WithMapping(t *testing.T) {
	// CRITICAL: this test verifies that restriction checks the ORIGINAL model
	// against pricing BEFORE applying mapping. The model "unknown-model" is NOT
	// in pricing, so even though the wildcard mapping "*" matches it, it should
	// still be restricted.
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4"placeholderplaceholder,
	placeholder,
		ModelMapping: map[string]map[string]string{
			"anthropic": {
				"*": "catch-all-target",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	gid := int64(10)
	mapping, restricted := svc.ResolveChannelMappingAndRestrict(context.Background(), &gid, "unknown-model")
	require.True(t, restricted) // model NOT in pricing, even though mapping exists
	require.True(t, mapping.Mapped)
	require.Equal(t, "catch-all-target", mapping.MappedModel)
placeholder

func TestResolveChannelMappingAndRestrict_ModelNotInPricing_NoMapping(t *testing.T) {
	ch := Channel{
		ID:             1,
		Status:         StatusActive,
		GroupIDs:       []int64{10placeholder,
		RestrictModels: true,
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-sonnet-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	gid := int64(10)
	mapping, restricted := svc.ResolveChannelMappingAndRestrict(context.Background(), &gid, "unknown-model")
	require.True(t, restricted) // model NOT in pricing
	require.False(t, mapping.Mapped)
	require.Equal(t, "unknown-model", mapping.MappedModel)
placeholder

// --- 4.6 Cache Building Specifics ---

func TestBuildCache_DBError(t *testing.T) {
	callCount := 0
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			callCount++
			return nil, errors.New("database down")
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	// First call should fail
	_, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.Contains(t, err.Error(), "database down")
	require.Equal(t, 1, callCount)

	// Second call within error-TTL should use error cache, but still return error
	// Because buildCache stores error-TTL cache and returns error, the cached value
	// is still within TTL and loadCache returns it (which is an empty cache).
	// Actually, re-reading the code: buildCache returns nil, err, and the error cache
	// only serves as a "don't retry immediately" mechanism. The singleflight.Do
	// returns the error. On next call within error-TTL, the cache has an empty but
	// valid entry, so loadCache returns it (with empty maps). GetChannelForGroup
	// will find nothing and return nil, nil.
	result, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.Nil(t, result)
	// Should NOT have hit DB again (error-TTL cache is active)
	require.Equal(t, 1, callCount)
placeholder

func TestBuildCache_GroupPlatformError(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return []Channel{chplaceholder, nil
	placeholder,
		getGroupPlatformsFn: func(_ context.Context, _ []int64) (map[int64]string, error) {
			return nil, errors.New("group platforms failed")
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	// Should degrade gracefully: channel is found, but without platform info
	// pricing won't match because platform will be "" and pricing platform is "anthropic"
	result, err := svc.GetChannelForGroup(context.Background(), 10)
placeholder
	require.NotNil(t, result) // channel still found
	require.Equal(t, int64(1), result.ID)
placeholder

func TestBuildCache_MultipleGroupsSameChannel(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10, 20, 30placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{
		10: "anthropic",
		20: "anthropic",
		30: "anthropic",
placeholder)
	svc := newTestChannelService(repo)

	for _, gid := range []int64{10, 20, 30placeholder {
		result := svc.GetChannelModelPricing(context.Background(), gid, "claude-opus-4")
		require.NotNil(t, result, "group %d should have pricing", gid)
		require.Equal(t, int64(100), result.ID)
placeholder
placeholder

func TestBuildCache_PlatformFiltering(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10, 20placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
			{ID: 200, Platform: "openai", Models: []string{"gpt-5.1"placeholderplaceholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{
		10: "anthropic",
		20: "openai",
placeholder)
	svc := newTestChannelService(repo)

	// anthropic group sees only anthropic models
	require.NotNil(t, svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4"))
	require.Nil(t, svc.GetChannelModelPricing(context.Background(), 10, "gpt-5.1"))

	// openai group sees only openai models
	require.NotNil(t, svc.GetChannelModelPricing(context.Background(), 20, "gpt-5.1"))
	require.Nil(t, svc.GetChannelModelPricing(context.Background(), 20, "claude-opus-4"))
placeholder

func TestBuildCache_WildcardPreservesConfigOrder(t *testing.T) {
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			// Configuration order: shortest prefix first
			{ID: 100, Platform: "anthropic", Models: []string{"c-*"placeholder, InputPrice: testPtrFloat64(1e-6)placeholder,
			{ID: 200, Platform: "anthropic", Models: []string{"c-son-*"placeholder, InputPrice: testPtrFloat64(2e-6)placeholder,
			{ID: 300, Platform: "anthropic", Models: []string{"c-son-4-*"placeholder, InputPrice: testPtrFloat64(3e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: "anthropic"placeholder)
	svc := newTestChannelService(repo)

	// "c-son-4-xxx" matches all three wildcards, but "c-*" (ID=100) is first in config
	result := svc.GetChannelModelPricing(context.Background(), 10, "c-son-4-xxx")
	require.NotNil(t, result)
	require.Equal(t, int64(100), result.ID)

	// "c-son-yyy" matches "c-*" and "c-son-*", but "c-*" (ID=100) is first
	result = svc.GetChannelModelPricing(context.Background(), 10, "c-son-yyy")
	require.NotNil(t, result)
	require.Equal(t, int64(100), result.ID)

	// "c-other" only matches "c-*" (ID=100)
	result = svc.GetChannelModelPricing(context.Background(), 10, "c-other")
	require.NotNil(t, result)
	require.Equal(t, int64(100), result.ID)
placeholder

// --- 4.7 invalidateCache ---

func TestInvalidateCache(t *testing.T) {
	callCount := 0
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			callCount++
			return []Channel{chplaceholder, nil
	placeholder,
		getGroupPlatformsFn: func(_ context.Context, _ []int64) (map[int64]string, error) {
			return map[int64]string{10: "anthropic"placeholder, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	// First load
	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result)
	require.Equal(t, 1, callCount)

	// Second call should use cache
	result = svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result)
	require.Equal(t, 1, callCount) // no new DB call

	// Invalidate
	svc.invalidateCache()

	// Next call should rebuild from DB
	result = svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.NotNil(t, result)
	require.Equal(t, 2, callCount) // rebuilt
placeholder

// ===========================================================================
// 5. CRUD Methods
// ===========================================================================

// --- 5.1 Create ---

func TestCreate_Success(t *testing.T) {
	createdID := int64(42)
	repo := &mockChannelRepository{
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
	placeholder,
		getGroupsInOtherChannelsFn: func(_ context.Context, _ int64, _ []int64) ([]int64, error) {
			return nil, nil
	placeholder,
		createFn: func(_ context.Context, ch *Channel) error {
			ch.ID = createdID
			return nil
	placeholder,
		getByIDFn: func(_ context.Context, id int64) (*Channel, error) {
			return &Channel{ID: id, Name: "new-channel", Status: StatusActiveplaceholder, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	result, err := svc.Create(context.Background(), &CreateChannelInput{
		Name:     "new-channel",
		GroupIDs: []int64{10placeholder,
placeholder)
placeholder
	require.NotNil(t, result)
	require.Equal(t, createdID, result.ID)
placeholder

func TestCreate_NameExists(t *testing.T) {
	repo := &mockChannelRepository{
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return true, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	_, err := svc.Create(context.Background(), &CreateChannelInput{
		Name: "existing-channel",
placeholder)
placeholder
	require.ErrorIs(t, err, ErrChannelExists)
placeholder

func TestCreate_GroupConflict(t *testing.T) {
	repo := &mockChannelRepository{
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
	placeholder,
		getGroupsInOtherChannelsFn: func(_ context.Context, _ int64, _ []int64) ([]int64, error) {
			return []int64{10placeholder, nil // group 10 already in another channel
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	_, err := svc.Create(context.Background(), &CreateChannelInput{
		Name:     "new-channel",
		GroupIDs: []int64{10, 20placeholder,
placeholder)
placeholder
	require.ErrorIs(t, err, ErrGroupAlreadyInChannel)
placeholder

func TestCreate_DuplicateModel(t *testing.T) {
	repo := &mockChannelRepository{
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	_, err := svc.Create(context.Background(), &CreateChannelInput{
		Name: "new-channel",
		ModelPricing: []ChannelModelPricing{
			{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
			{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder, // duplicate
	placeholder,
placeholder)
placeholder
	require.Contains(t, err.Error(), "claude-opus-4")
placeholder

func TestCreate_DefaultBillingModelSource(t *testing.T) {
	var capturedChannel *Channel
	repo := &mockChannelRepository{
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
	placeholder,
		createFn: func(_ context.Context, ch *Channel) error {
			capturedChannel = ch
			ch.ID = 1
			return nil
	placeholder,
		getByIDFn: func(_ context.Context, id int64) (*Channel, error) {
			return capturedChannel, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	result, err := svc.Create(context.Background(), &CreateChannelInput{
		Name:               "new-channel",
		BillingModelSource: "", // empty, should default to "channel_mapped"
placeholder)
placeholder
	require.NotNil(t, result)
	require.Equal(t, BillingModelSourceChannelMapped, result.BillingModelSource)
placeholder

func TestCreate_InvalidatesCache(t *testing.T) {
	loadCount := 0
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
	placeholder,
placeholder
	repo := &mockChannelRepository{
		listAllFn: func(_ context.Context) ([]Channel, error) {
			loadCount++
			return []Channel{chplaceholder, nil
	placeholder,
		getGroupPlatformsFn: func(_ context.Context, _ []int64) (map[int64]string, error) {
			return map[int64]string{10: "anthropic"placeholder, nil
	placeholder,
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
	placeholder,
		createFn: func(_ context.Context, c *Channel) error {
			c.ID = 2
			return nil
	placeholder,
		getByIDFn: func(_ context.Context, id int64) (*Channel, error) {
			return &Channel{ID: id, Name: "new", Status: StatusActiveplaceholder, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	// Load cache
	_ = svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.Equal(t, 1, loadCount)

	// Create triggers cache invalidation
	_, err := svc.Create(context.Background(), &CreateChannelInput{Name: "new"placeholder)
placeholder

	// Next cache access should rebuild
	_ = svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4")
	require.Equal(t, 2, loadCount)
placeholder

// --- 5.2 Update ---

func TestUpdate_Success(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "original",
		Status: StatusActive,
placeholder
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, id int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
		updateFn: func(_ context.Context, _ *Channel) error {
			return nil
	placeholder,
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return nil, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	result, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		Name:        "updated-name",
		Description: testPtrString("new desc"),
placeholder)
placeholder
	require.NotNil(t, result)
placeholder

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, _ int64) (*Channel, error) {
			return nil, ErrChannelNotFound
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	_, err := svc.Update(context.Background(), 999, &UpdateChannelInput{
		Name: "whatever",
placeholder)
placeholder
	require.Contains(t, err.Error(), "channel")
placeholder

func TestUpdate_NameConflict(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "original",
		Status: StatusActive,
placeholder
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, _ int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
		existsByNameExcludingFn: func(_ context.Context, _ string, _ int64) (bool, error) {
			return true, nil // name conflicts with another channel
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	_, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		Name: "conflicting-name",
placeholder)
placeholder
	require.ErrorIs(t, err, ErrChannelExists)
placeholder

func TestUpdate_GroupConflict(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "original",
		Status: StatusActive,
placeholder
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, _ int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
		getGroupsInOtherChannelsFn: func(_ context.Context, _ int64, _ []int64) ([]int64, error) {
			return []int64{20placeholder, nil // group 20 in another channel
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	newGroupIDs := []int64{10, 20placeholder
	_, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		GroupIDs: &newGroupIDs,
placeholder)
placeholder
	require.ErrorIs(t, err, ErrGroupAlreadyInChannel)
placeholder

func TestUpdate_DuplicateModel(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "original",
		Status: StatusActive,
placeholder
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, _ int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	dupPricing := []ChannelModelPricing{
		{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
		{Platform: "anthropic", Models: []string{"claude-opus-4"placeholderplaceholder,
placeholder
	_, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		ModelPricing: &dupPricing,
placeholder)
placeholder
	require.Contains(t, err.Error(), "claude-opus-4")
placeholder

func TestUpdate_InvalidatesChannelCache(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "original",
		Status: StatusActive,
placeholder
	loadCount := 0
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, _ int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
		updateFn: func(_ context.Context, _ *Channel) error {
			return nil
	placeholder,
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return []int64{10, 20placeholder, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			loadCount++
			return []Channel{*existingplaceholder, nil
	placeholder,
		getGroupPlatformsFn: func(_ context.Context, _ []int64) (map[int64]string, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	// Load cache first
	_, _ = svc.GetChannelForGroup(context.Background(), 10)
	require.Equal(t, 1, loadCount)

	result, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		Description: testPtrString("updated"),
placeholder)
placeholder
	require.NotNil(t, result)

	// Channel cache should be invalidated (next access rebuilds)
	_, _ = svc.GetChannelForGroup(context.Background(), 10)
	require.Equal(t, 2, loadCount)
placeholder

func TestUpdate_InvalidatesAuthCache(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "original",
		Status: StatusActive,
placeholder
	auth := &mockChannelAuthCacheInvalidator{placeholder
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, _ int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
		updateFn: func(_ context.Context, _ *Channel) error {
			return nil
	placeholder,
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return []int64{10, 20placeholder, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelServiceWithAuth(repo, auth)

	result, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		Description: testPtrString("updated"),
placeholder)
placeholder
	require.NotNil(t, result)

	// Auth cache should be invalidated for both groups
	require.ElementsMatch(t, []int64{10, 20placeholder, auth.invalidatedGroupIDs)
placeholder

// --- 5.3 Delete ---

func TestChannelDelete_Success(t *testing.T) {
	deleted := false
	repo := &mockChannelRepository{
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return nil, nil
	placeholder,
		deleteFn: func(_ context.Context, _ int64) error {
			deleted = true
			return nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	err := svc.Delete(context.Background(), 1)
placeholder
	require.True(t, deleted)
placeholder

func TestChannelDelete_InvalidatesCaches(t *testing.T) {
	auth := &mockChannelAuthCacheInvalidator{placeholder
	loadCount := 0
	repo := &mockChannelRepository{
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return []int64{10, 20placeholder, nil
	placeholder,
		deleteFn: func(_ context.Context, _ int64) error {
			return nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			loadCount++
			return []Channel{{ID: 1, Status: StatusActive, GroupIDs: []int64{10, 20placeholderplaceholderplaceholder, nil
	placeholder,
		getGroupPlatformsFn: func(_ context.Context, _ []int64) (map[int64]string, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelServiceWithAuth(repo, auth)

	// Load cache first
	_, _ = svc.GetChannelForGroup(context.Background(), 10)
	require.Equal(t, 1, loadCount)

	err := svc.Delete(context.Background(), 1)
placeholder

	// Auth cache invalidated for both groups
	require.ElementsMatch(t, []int64{10, 20placeholder, auth.invalidatedGroupIDs)

	// Channel cache invalidated
	_, _ = svc.GetChannelForGroup(context.Background(), 10)
	require.Equal(t, 2, loadCount)
placeholder

func TestChannelDelete_NotFound(t *testing.T) {
	repo := &mockChannelRepository{
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return nil, nil
	placeholder,
		deleteFn: func(_ context.Context, _ int64) error {
			return errors.New("record not found")
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	err := svc.Delete(context.Background(), 999)
placeholder
	require.Contains(t, err.Error(), "not found")
placeholder

// ===========================================================================
// 6. Edge Case Tests
// ===========================================================================

// --- 6.1 Create with empty GroupIDs ---

func TestCreate_NoGroups(t *testing.T) {
	createdID := int64(55)
	getGroupsInOtherChannelsCalled := false
	repo := &mockChannelRepository{
		existsByNameFn: func(_ context.Context, _ string) (bool, error) {
			return false, nil
	placeholder,
		getGroupsInOtherChannelsFn: func(_ context.Context, _ int64, _ []int64) ([]int64, error) {
			getGroupsInOtherChannelsCalled = true
			return nil, nil
	placeholder,
		createFn: func(_ context.Context, ch *Channel) error {
			ch.ID = createdID
			return nil
	placeholder,
		getByIDFn: func(_ context.Context, id int64) (*Channel, error) {
			return &Channel{ID: id, Name: "no-groups-channel", Status: StatusActiveplaceholder, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	result, err := svc.Create(context.Background(), &CreateChannelInput{
		Name:     "no-groups-channel",
		GroupIDs: []int64{placeholder, // empty slice
placeholder)
placeholder
	require.NotNil(t, result)
	require.Equal(t, createdID, result.ID)
	// GetGroupsInOtherChannels should NOT have been called (skipped by len(input.GroupIDs) > 0)
	require.False(t, getGroupsInOtherChannelsCalled)
placeholder

// --- 6.2 Update only Status ---

func TestUpdate_StatusOnly(t *testing.T) {
	existing := &Channel{
		ID:     1,
		Name:   "test-channel",
		Status: StatusActive,
placeholder
	var capturedChannel *Channel
	repo := &mockChannelRepository{
		getByIDFn: func(_ context.Context, id int64) (*Channel, error) {
			return existing.Clone(), nil
	placeholder,
		updateFn: func(_ context.Context, ch *Channel) error {
			capturedChannel = ch
			return nil
	placeholder,
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return nil, nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	result, err := svc.Update(context.Background(), 1, &UpdateChannelInput{
		Status: StatusDisabled,
placeholder)
placeholder
	require.NotNil(t, result)
	// Verify that the channel passed to repo.Update has the new status
	require.NotNil(t, capturedChannel)
	require.Equal(t, StatusDisabled, capturedChannel.Status)
	// Name should remain unchanged
	require.Equal(t, "test-channel", capturedChannel.Name)
placeholder

// --- 6.3 Delete when GetGroupIDs fails ---

func TestChannelDelete_GetGroupIDsError(t *testing.T) {
	deleted := false
	repo := &mockChannelRepository{
		getGroupIDsFn: func(_ context.Context, _ int64) ([]int64, error) {
			return nil, errors.New("group IDs lookup failed")
	placeholder,
		deleteFn: func(_ context.Context, _ int64) error {
			deleted = true
			return nil
	placeholder,
		listAllFn: func(_ context.Context) ([]Channel, error) {
			return nil, nil
	placeholder,
placeholder
	svc := newTestChannelService(repo)

	// Delete should still succeed even though GetGroupIDs returned error (degradation path L588-591)
	err := svc.Delete(context.Background(), 1)
placeholder
	require.True(t, deleted)
placeholder

// --- 6.4 ReplaceModelInBody with invalid JSON ---

func TestReplaceModelInBody_InvalidJSON(t *testing.T) {
	// Case 1: broken JSON object — gjson won't find "model", sjson does best-effort set
	// (no panic, no error from sjson, but result is mutated garbage)
	brokenBody := []byte("{broken")
	result := ReplaceModelInBody(brokenBody, "new-model")
	require.NotNil(t, result)
	// sjson does not error on this input, so result differs from original — just verify no panic

	// Case 2: JSON array — sjson.SetBytes returns error on non-object,
	// triggering the L447 error fallback path that returns original body.
	arrayBody := []byte("[]")
	result2 := ReplaceModelInBody(arrayBody, "new-model")
	require.Equal(t, arrayBody, result2)
placeholder

// ===========================================================================
// 7. isPlatformPricingMatch
// ===========================================================================

func TestIsPlatformPricingMatch(t *testing.T) {
	tests := []struct {
		name            string
		groupPlatform   string
		pricingPlatform string
		want            bool
placeholder{
		{"antigravity matches anthropic", PlatformAntigravity, PlatformAnthropic, trueplaceholder,
		{"antigravity matches gemini", PlatformAntigravity, PlatformGemini, trueplaceholder,
		{"antigravity matches antigravity", PlatformAntigravity, PlatformAntigravity, trueplaceholder,
		{"antigravity does NOT match openai", PlatformAntigravity, PlatformOpenAI, falseplaceholder,
		{"anthropic matches anthropic", PlatformAnthropic, PlatformAnthropic, trueplaceholder,
		{"anthropic does NOT match antigravity", PlatformAnthropic, PlatformAntigravity, falseplaceholder,
		{"anthropic does NOT match gemini", PlatformAnthropic, PlatformGemini, falseplaceholder,
		{"gemini matches gemini", PlatformGemini, PlatformGemini, trueplaceholder,
		{"gemini does NOT match antigravity", PlatformGemini, PlatformAntigravity, falseplaceholder,
		{"gemini does NOT match anthropic", PlatformGemini, PlatformAnthropic, falseplaceholder,
		{"empty string matches nothing", "", PlatformAnthropic, falseplaceholder,
		{"empty string matches empty", "", "", trueplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, isPlatformPricingMatch(tt.groupPlatform, tt.pricingPlatform))
	placeholder)
placeholder
placeholder

// ===========================================================================
// 8. matchingPlatforms
// ===========================================================================

func TestMatchingPlatforms(t *testing.T) {
	tests := []struct {
		name          string
		groupPlatform string
		want          []string
placeholder{
		{"antigravity returns all three", PlatformAntigravity, []string{PlatformAntigravity, PlatformAnthropic, PlatformGeminiplaceholderplaceholder,
		{"anthropic returns itself", PlatformAnthropic, []string{PlatformAnthropicplaceholderplaceholder,
		{"gemini returns itself", PlatformGemini, []string{PlatformGeminiplaceholderplaceholder,
		{"openai returns itself", PlatformOpenAI, []string{PlatformOpenAIplaceholderplaceholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchingPlatforms(tt.groupPlatform)
			require.Equal(t, tt.want, result)
	placeholder)
placeholder
placeholder

// ===========================================================================
// 9. Antigravity cross-platform channel pricing
// ===========================================================================

func TestGetChannelModelPricing_AntigravityCrossPlatform(t *testing.T) {
	// Channel has anthropic pricing for claude-opus-4-6.
	// Group 10 is antigravity — should see the anthropic pricing.
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: PlatformAnthropic, Models: []string{"claude-opus-4-6"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: PlatformAntigravityplaceholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4-6")
	require.NotNil(t, result, "antigravity group should see anthropic pricing")
	require.Equal(t, int64(100), result.ID)
	require.InDelta(t, 15e-6, *result.InputPrice, 1e-12)
placeholder

func TestGetChannelModelPricing_AnthropicCannotSeeAntigravityPricing(t *testing.T) {
	// Channel has antigravity-platform pricing for claude-opus-4-6.
	// Group 10 is anthropic — should NOT see antigravity pricing (no cross-platform leakage).
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelPricing: []ChannelModelPricing{
			{ID: 100, Platform: PlatformAntigravity, Models: []string{"claude-opus-4-6"placeholder, InputPrice: testPtrFloat64(15e-6)placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: PlatformAnthropicplaceholder)
	svc := newTestChannelService(repo)

	result := svc.GetChannelModelPricing(context.Background(), 10, "claude-opus-4-6")
	require.Nil(t, result, "anthropic group should NOT see antigravity-platform pricing")
placeholder

// ===========================================================================
// 10. Antigravity cross-platform model mapping
// ===========================================================================

func TestResolveChannelMapping_AntigravityCrossPlatform(t *testing.T) {
	// Channel has anthropic model mapping: claude-opus-4-5 → claude-opus-4-6.
	// Group 10 is antigravity — should apply the anthropic mapping.
	ch := Channel{
		ID:       1,
		Status:   StatusActive,
		GroupIDs: []int64{10placeholder,
		ModelMapping: map[string]map[string]string{
			PlatformAnthropic: {
				"claude-opus-4-5": "claude-opus-4-6",
		placeholder,
	placeholder,
placeholder
	repo := makeStandardRepo(ch, map[int64]string{10: PlatformAntigravityplaceholder)
	svc := newTestChannelService(repo)

	result := svc.ResolveChannelMapping(context.Background(), 10, "claude-opus-4-5")
	require.True(t, result.Mapped, "antigravity group should apply anthropic mapping")
	require.Equal(t, "claude-opus-4-6", result.MappedModel)
	require.Equal(t, int64(1), result.ChannelID)
placeholder
