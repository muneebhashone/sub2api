//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

func ptrString[T ~string](v T) *string {
	s := string(v)
	return &s
placeholder

// groupRepoStubForAdmin 用于测试 AdminService 的 GroupRepository Stub
type groupRepoStubForAdmin struct {
	created  *Group // 记录 Create 调用的参数
	updated  *Group // 记录 Update 调用的参数
	getByID  *Group // GetByID 返回值
	getErr   error  // GetByID 返回的错误
	createID int64

	getByIDByID map[int64]*Group

	deleteAccountGroupsByGroupIDFn func(groupID int64) (int64, error)
	bindAccountsToGroupFn          func(groupID int64, accountIDs []int64) error
	getAccountIDsByGroupIDsFn      func(groupIDs []int64) ([]int64, error)

	listWithFiltersCalls       int
	listWithFiltersParams      pagination.PaginationParams
	listWithFiltersPlatform    string
	listWithFiltersStatus      string
	listWithFiltersSearch      string
	listWithFiltersIsExclusive *bool
	listWithFiltersGroups      []Group
	listWithFiltersResult      *pagination.PaginationResult
	listWithFiltersErr         error
placeholder

func (s *groupRepoStubForAdmin) Create(_ context.Context, g *Group) error {
	if s.createID > 0 {
		g.ID = s.createID
placeholder
	s.created = g
	return nil
placeholder

func (s *groupRepoStubForAdmin) Update(_ context.Context, g *Group) error {
	s.updated = g
	return nil
placeholder

func (s *groupRepoStubForAdmin) GetByID(_ context.Context, id int64) (*Group, error) {
	if s.getErr != nil {
		return nil, s.getErr
placeholder
	if s.getByIDByID != nil {
		if group, ok := s.getByIDByID[id]; ok {
			return group, nil
	placeholder
		return nil, ErrGroupNotFound
placeholder
	return s.getByID, nil
placeholder

func (s *groupRepoStubForAdmin) GetByIDLite(_ context.Context, id int64) (*Group, error) {
	if s.getErr != nil {
		return nil, s.getErr
placeholder
	if s.getByIDByID != nil {
		if group, ok := s.getByIDByID[id]; ok {
			return group, nil
	placeholder
		return nil, ErrGroupNotFound
placeholder
	return s.getByID, nil
placeholder

func (s *groupRepoStubForAdmin) Delete(_ context.Context, _ int64) error {
	panic("unexpected Delete call")
placeholder

func (s *groupRepoStubForAdmin) DeleteCascade(_ context.Context, _ int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
placeholder

func (s *groupRepoStubForAdmin) List(_ context.Context, _ pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *groupRepoStubForAdmin) ListWithFilters(_ context.Context, params pagination.PaginationParams, platform, status, search string, isExclusive *bool) ([]Group, *pagination.PaginationResult, error) {
	s.listWithFiltersCalls++
	s.listWithFiltersParams = params
	s.listWithFiltersPlatform = platform
	s.listWithFiltersStatus = status
	s.listWithFiltersSearch = search
	s.listWithFiltersIsExclusive = isExclusive

	if s.listWithFiltersErr != nil {
		return nil, nil, s.listWithFiltersErr
placeholder

	result := s.listWithFiltersResult
	if result == nil {
		result = &pagination.PaginationResult{
			Total:    int64(len(s.listWithFiltersGroups)),
			Page:     params.Page,
			PageSize: params.PageSize,
	placeholder
placeholder

	return s.listWithFiltersGroups, result, nil
placeholder

func (s *groupRepoStubForAdmin) ListActive(_ context.Context) ([]Group, error) {
	panic("unexpected ListActive call")
placeholder

func (s *groupRepoStubForAdmin) ListActiveByPlatform(_ context.Context, _ string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
placeholder

func (s *groupRepoStubForAdmin) ExistsByName(_ context.Context, _ string) (bool, error) {
	panic("unexpected ExistsByName call")
placeholder

func (s *groupRepoStubForAdmin) GetAccountCount(_ context.Context, _ int64) (int64, int64, error) {
	panic("unexpected GetAccountCount call")
placeholder

func (s *groupRepoStubForAdmin) DeleteAccountGroupsByGroupID(_ context.Context, groupID int64) (int64, error) {
	if s.deleteAccountGroupsByGroupIDFn != nil {
		return s.deleteAccountGroupsByGroupIDFn(groupID)
placeholder
	panic("unexpected DeleteAccountGroupsByGroupID call")
placeholder

func (s *groupRepoStubForAdmin) BindAccountsToGroup(_ context.Context, groupID int64, accountIDs []int64) error {
	if s.bindAccountsToGroupFn != nil {
		return s.bindAccountsToGroupFn(groupID, accountIDs)
placeholder
	panic("unexpected BindAccountsToGroup call")
placeholder

func (s *groupRepoStubForAdmin) GetAccountIDsByGroupIDs(_ context.Context, groupIDs []int64) ([]int64, error) {
	if s.getAccountIDsByGroupIDsFn != nil {
		return s.getAccountIDsByGroupIDsFn(groupIDs)
placeholder
	panic("unexpected GetAccountIDsByGroupIDs call")
placeholder

func (s *groupRepoStubForAdmin) UpdateSortOrders(_ context.Context, _ []GroupSortOrderUpdate) error {
	return nil
placeholder

type compositeRouteRepoStubForAdmin struct {
	routes    []CompositeModelRoute
	created   *CompositeModelRoute
	updated   *CompositeModelRoute
	deleted   []int64
	nextID    int64
	listErr   error
	createErr error
	updateErr error
	deleteErr error
placeholder

func (s *compositeRouteRepoStubForAdmin) ListByGroup(_ context.Context, groupID int64, includeDisabled bool) ([]CompositeModelRoute, error) {
	if s.listErr != nil {
		return nil, s.listErr
placeholder
	routes := make([]CompositeModelRoute, 0, len(s.routes))
	for _, route := range s.routes {
		if route.GroupID != groupID {
			continue
	placeholder
		if !includeDisabled && !route.Enabled {
			continue
	placeholder
		routes = append(routes, route)
placeholder
	return routes, nil
placeholder

func (s *compositeRouteRepoStubForAdmin) Create(_ context.Context, route *CompositeModelRoute) error {
	if s.createErr != nil {
		return s.createErr
placeholder
	if s.nextID > 0 {
		route.ID = s.nextID
placeholder
	cloned := *route
	s.created = &cloned
	s.routes = append(s.routes, cloned)
	return nil
placeholder

func (s *compositeRouteRepoStubForAdmin) Update(_ context.Context, route *CompositeModelRoute) error {
	if s.updateErr != nil {
		return s.updateErr
placeholder
	cloned := *route
	s.updated = &cloned
	for i := range s.routes {
		if s.routes[i].ID == route.ID {
			s.routes[i] = cloned
			return nil
	placeholder
placeholder
	s.routes = append(s.routes, cloned)
	return nil
placeholder

func (s *compositeRouteRepoStubForAdmin) Delete(_ context.Context, id int64) error {
	if s.deleteErr != nil {
		return s.deleteErr
placeholder
	s.deleted = append(s.deleted, id)
	return nil
placeholder

func (s *compositeRouteRepoStubForAdmin) DeleteByGroup(_ context.Context, groupID int64) error {
	next := s.routes[:0]
	for _, route := range s.routes {
		if route.GroupID != groupID {
			next = append(next, route)
	placeholder
placeholder
	s.routes = next
	return nil
placeholder

func TestAdminService_ListGroups_PassesSortParams(t *testing.T) {
	repo := &groupRepoStubForAdmin{
		listWithFiltersGroups: []Group{{ID: 1, Name: "g1"placeholderplaceholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, _, err := svc.ListGroups(context.Background(), 3, 25, PlatformOpenAI, StatusActive, "needle", nil, "account_count", "ASC")
placeholder
	require.Equal(t, pagination.PaginationParams{
		Page:      3,
		PageSize:  25,
		SortBy:    "account_count",
		SortOrder: "ASC",
placeholder, repo.listWithFiltersParams)
placeholder

// TestAdminService_CreateGroup_WithImagePricing 测试创建分组时 ImagePrice 字段正确传递
func TestAdminService_CreateGroup_WithImagePricing(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	price1K := 0.10
	price2K := 0.15
	price4K := 0.30

	input := &CreateGroupInput{
		Name:           "test-group",
		Description:    "Test group",
		Platform:       PlatformAntigravity,
		RateMultiplier: 1.0,
		ImagePrice1K:   &price1K,
		ImagePrice2K:   &price2K,
		ImagePrice4K:   &price4K,
placeholder

	group, err := svc.CreateGroup(context.Background(), input)
placeholder
	require.NotNil(t, group)

	// 验证 repo 收到了正确的字段
	require.NotNil(t, repo.created)
	require.NotNil(t, repo.created.ImagePrice1K)
	require.NotNil(t, repo.created.ImagePrice2K)
	require.NotNil(t, repo.created.ImagePrice4K)
	require.InDelta(t, 0.10, *repo.created.ImagePrice1K, 0.0001)
	require.InDelta(t, 0.15, *repo.created.ImagePrice2K, 0.0001)
	require.InDelta(t, 0.30, *repo.created.ImagePrice4K, 0.0001)
placeholder

func TestAdminService_CreateGroup_WithVideoPricing(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	price480P := 0.08
	price720P := 0.12
	price1080P := 0.18
	videoMultiplier := 0.75

	input := &CreateGroupInput{
		Name:                 "grok-video",
		Description:          "Grok video group",
		Platform:             PlatformGrok,
		RateMultiplier:       1.0,
		VideoRateIndependent: true,
		VideoRateMultiplier:  &videoMultiplier,
		VideoPrice480P:       &price480P,
		VideoPrice720P:       &price720P,
		VideoPrice1080P:      &price1080P,
placeholder

	group, err := svc.CreateGroup(context.Background(), input)
placeholder
	require.NotNil(t, group)

	require.NotNil(t, repo.created)
	require.True(t, repo.created.VideoRateIndependent)
	require.InDelta(t, 0.75, repo.created.VideoRateMultiplier, 1e-12)
	require.NotNil(t, repo.created.VideoPrice480P)
	require.NotNil(t, repo.created.VideoPrice720P)
	require.NotNil(t, repo.created.VideoPrice1080P)
	require.InDelta(t, 0.08, *repo.created.VideoPrice480P, 0.0001)
	require.InDelta(t, 0.12, *repo.created.VideoPrice720P, 0.0001)
	require.InDelta(t, 0.18, *repo.created.VideoPrice1080P, 0.0001)
placeholder

// TestAdminService_CreateGroup_NilImagePricing 测试 ImagePrice 为 nil 时正常创建
func TestAdminService_CreateGroup_NilImagePricing(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	input := &CreateGroupInput{
		Name:           "test-group",
		Description:    "Test group",
		Platform:       PlatformAntigravity,
		RateMultiplier: 1.0,
		// ImagePrice 字段全部为 nil
placeholder

	group, err := svc.CreateGroup(context.Background(), input)
placeholder
	require.NotNil(t, group)

	// 验证 ImagePrice 字段为 nil
	require.NotNil(t, repo.created)
	require.Nil(t, repo.created.ImagePrice1K)
	require.Nil(t, repo.created.ImagePrice2K)
	require.Nil(t, repo.created.ImagePrice4K)
placeholder

func TestAdminService_CreateGroup_DefaultsGrokMediaGenerationEnabled(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:           "grok-media",
		Description:    "Grok media group",
		Platform:       PlatformGrok,
		RateMultiplier: 1.0,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.True(t, repo.created.AllowImageGeneration)
	require.True(t, group.AllowImageGeneration)
placeholder

func TestAdminService_CreateGroup_PreservesNonGrokImageGenerationDisabled(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:           "anthropic-text",
		Description:    "Anthropic text group",
		Platform:       PlatformAnthropic,
		RateMultiplier: 1.0,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.False(t, repo.created.AllowImageGeneration)
	require.False(t, group.AllowImageGeneration)
placeholder

func TestAdminService_CreateGroup_DisablesBatchImageWhenImageGenerationDisabled(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                      "gemini-no-image",
		Description:               "Gemini group without image generation",
		Platform:                  PlatformGemini,
		RateMultiplier:            1.0,
		AllowImageGeneration:      false,
		AllowBatchImageGeneration: true,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.False(t, repo.created.AllowImageGeneration)
	require.False(t, repo.created.AllowBatchImageGeneration)
	require.False(t, group.AllowBatchImageGeneration)
placeholder

func TestAdminService_CreateGroup_DisablesBatchImageForNonGeminiPlatform(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                      "openai-image",
		Description:               "OpenAI image group",
		Platform:                  PlatformOpenAI,
		RateMultiplier:            1.0,
		AllowImageGeneration:      true,
		AllowBatchImageGeneration: true,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.True(t, repo.created.AllowImageGeneration)
	require.False(t, repo.created.AllowBatchImageGeneration)
	require.False(t, group.AllowBatchImageGeneration)
placeholder

// TestAdminService_UpdateGroup_WithImagePricing 测试更新分组时 ImagePrice 字段正确更新
func TestAdminService_UpdateGroup_WithImagePricing(t *testing.T) {
	existingGroup := &Group{
		ID:       1,
		Name:     "existing-group",
		Platform: PlatformAntigravity,
		Status:   StatusActive,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	price1K := 0.12
	price2K := 0.18
	price4K := 0.36

	input := &UpdateGroupInput{
		ImagePrice1K: &price1K,
		ImagePrice2K: &price2K,
		ImagePrice4K: &price4K,
placeholder

	group, err := svc.UpdateGroup(context.Background(), 1, input)
placeholder
	require.NotNil(t, group)

	// 验证 repo 收到了更新后的字段
	require.NotNil(t, repo.updated)
	require.NotNil(t, repo.updated.ImagePrice1K)
	require.NotNil(t, repo.updated.ImagePrice2K)
	require.NotNil(t, repo.updated.ImagePrice4K)
	require.InDelta(t, 0.12, *repo.updated.ImagePrice1K, 0.0001)
	require.InDelta(t, 0.18, *repo.updated.ImagePrice2K, 0.0001)
	require.InDelta(t, 0.36, *repo.updated.ImagePrice4K, 0.0001)
placeholder

func TestAdminService_UpdateGroup_WithVideoPricing(t *testing.T) {
	existingGroup := &Group{
		ID:       1,
		Name:     "existing-grok",
		Platform: PlatformGrok,
		Status:   StatusActive,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	price480P := 0.09
	price720P := 0.13
	price1080P := 0.19
	videoMultiplier := 0.6
	independent := true

	input := &UpdateGroupInput{
		VideoRateIndependent: &independent,
		VideoRateMultiplier:  &videoMultiplier,
		VideoPrice480P:       &price480P,
		VideoPrice720P:       &price720P,
		VideoPrice1080P:      &price1080P,
placeholder

	group, err := svc.UpdateGroup(context.Background(), 1, input)
placeholder
	require.NotNil(t, group)

	require.NotNil(t, repo.updated)
	require.True(t, repo.updated.VideoRateIndependent)
	require.InDelta(t, 0.6, repo.updated.VideoRateMultiplier, 1e-12)
	require.InDelta(t, 0.09, *repo.updated.VideoPrice480P, 0.0001)
	require.InDelta(t, 0.13, *repo.updated.VideoPrice720P, 0.0001)
	require.InDelta(t, 0.19, *repo.updated.VideoPrice1080P, 0.0001)
placeholder

// TestAdminService_UpdateGroup_PartialImagePricing 测试仅更新部分 ImagePrice 字段
func TestAdminService_UpdateGroup_PartialImagePricing(t *testing.T) {
	oldPrice2K := 0.15
	existingGroup := &Group{
		ID:           1,
		Name:         "existing-group",
		Platform:     PlatformAntigravity,
		Status:       StatusActive,
		ImagePrice2K: &oldPrice2K, // 已有 2K 价格
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	// 只更新 1K 价格
	price1K := 0.10
	input := &UpdateGroupInput{
		ImagePrice1K: &price1K,
		// ImagePrice2K 和 ImagePrice4K 为 nil，不更新
placeholder

	group, err := svc.UpdateGroup(context.Background(), 1, input)
placeholder
	require.NotNil(t, group)

	// 验证：1K 被更新，2K 保持原值，4K 仍为 nil
	require.NotNil(t, repo.updated)
	require.NotNil(t, repo.updated.ImagePrice1K)
	require.InDelta(t, 0.10, *repo.updated.ImagePrice1K, 0.0001)
	require.NotNil(t, repo.updated.ImagePrice2K)
	require.InDelta(t, 0.15, *repo.updated.ImagePrice2K, 0.0001) // 原值保持
	require.Nil(t, repo.updated.ImagePrice4K)
placeholder

func TestAdminService_UpdateGroup_PreservesImageGenerationControlsWhenOmitted(t *testing.T) {
	imageMultiplier := 0.5
	existingGroup := &Group{
		ID:                   1,
		Name:                 "existing-group",
		Platform:             PlatformOpenAI,
		Status:               StatusActive,
		AllowImageGeneration: true,
		ImageRateIndependent: true,
		ImageRateMultiplier:  imageMultiplier,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	updatedDesc := "updated"
	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		Description: &updatedDesc,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.True(t, repo.updated.AllowImageGeneration)
	require.True(t, repo.updated.ImageRateIndependent)
	require.InDelta(t, 0.5, repo.updated.ImageRateMultiplier, 1e-12)
placeholder

func TestAdminService_UpdateGroup_DisablesBatchImageWhenImageGenerationDisabled(t *testing.T) {
	existingGroup := &Group{
		ID:                        1,
		Name:                      "existing-gemini",
		Platform:                  PlatformGemini,
		Status:                    StatusActive,
		AllowImageGeneration:      true,
		AllowBatchImageGeneration: true,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder
	disabled := false

	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		AllowImageGeneration: &disabled,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.False(t, repo.updated.AllowImageGeneration)
	require.False(t, repo.updated.AllowBatchImageGeneration)
	require.False(t, group.AllowBatchImageGeneration)
placeholder

func TestAdminService_UpdateGroup_DisablesBatchImageWhenPlatformChangesFromGemini(t *testing.T) {
	existingGroup := &Group{
		ID:                        1,
		Name:                      "existing-gemini",
		Platform:                  PlatformGemini,
		Status:                    StatusActive,
		AllowImageGeneration:      true,
		AllowBatchImageGeneration: true,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		Platform: PlatformOpenAI,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, PlatformOpenAI, repo.updated.Platform)
	require.False(t, repo.updated.AllowBatchImageGeneration)
	require.False(t, group.AllowBatchImageGeneration)
placeholder

func TestAdminService_UpdateGroup_ClearsDescriptionWhenEmptyString(t *testing.T) {
	existingGroup := &Group{
		ID:          1,
		Name:        "existing-group",
		Description: "Auto-created default group",
		Platform:    PlatformOpenAI,
		Status:      StatusActive,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	empty := ""
	_, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		Description: &empty,
placeholder)
placeholder
	require.NotNil(t, repo.updated)
	require.Equal(t, "", repo.updated.Description, "empty string should clear description")
placeholder

func TestAdminService_UpdateGroup_PreservesDescriptionWhenNil(t *testing.T) {
	existingGroup := &Group{
		ID:          1,
		Name:        "existing-group",
		Description: "keep me",
		Platform:    PlatformOpenAI,
		Status:      StatusActive,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		Description: nil,
placeholder)
placeholder
	require.NotNil(t, repo.updated)
	require.Equal(t, "keep me", repo.updated.Description, "nil should preserve existing description")
placeholder

func TestAdminService_UpdateGroup_RejectsNegativeImageRateMultiplier(t *testing.T) {
	existingGroup := &Group{
		ID:                  1,
		Name:                "existing-group",
		Platform:            PlatformOpenAI,
		Status:              StatusActive,
		ImageRateMultiplier: 1,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder
	negative := -0.1

	_, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		ImageRateMultiplier: &negative,
placeholder)
placeholder
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_CreateGroup_BatchImagePricingSettings(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder
	discount := 0.8
	hold := 0.9

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                         "batch-image-pricing",
		Platform:                     PlatformGemini,
		RateMultiplier:               1,
		BatchImageDiscountMultiplier: &discount,
		BatchImageHoldMultiplier:     &hold,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.InDelta(t, 0.8, repo.created.BatchImageDiscountMultiplier, 1e-12)
	require.InDelta(t, 0.9, repo.created.BatchImageHoldMultiplier, 1e-12)
placeholder

func TestAdminService_CreateGroup_RejectsHoldBelowDiscount(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder
	discount := 0.8
	hold := 0.6

	// hold < discount 时，成功率足够高的批量任务实际成本会超过冻结额，
	// 结算永远失败，必须在配置入口拒绝。
	_, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                         "batch-image-pricing-invalid",
		Platform:                     PlatformGemini,
		RateMultiplier:               1,
		BatchImageDiscountMultiplier: &discount,
		BatchImageHoldMultiplier:     &hold,
placeholder)
placeholder
	require.Nil(t, repo.created)
placeholder

func TestAdminService_GroupBatchImagePricingValidation(t *testing.T) {
	tests := []struct {
		name  string
		input *CreateGroupInput
placeholder{
		{
			name: "negative_discount",
			input: func() *CreateGroupInput {
				v := -0.1
				return &CreateGroupInput{Name: "bad-discount", RateMultiplier: 1, BatchImageDiscountMultiplier: &vplaceholder
		placeholder(),
	placeholder,
		{
			name: "negative_hold",
			input: func() *CreateGroupInput {
				v := -0.1
				return &CreateGroupInput{Name: "bad-hold", RateMultiplier: 1, BatchImageHoldMultiplier: &vplaceholder
		placeholder(),
	placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &groupRepoStubForAdmin{placeholder
			svc := &adminServiceImpl{groupRepo: repoplaceholder

			_, err := svc.CreateGroup(context.Background(), tt.input)
		placeholder
			require.Nil(t, repo.created)
	placeholder)
placeholder
placeholder

func TestAdminService_UpdateGroup_RejectsNegativeVideoRateMultiplier(t *testing.T) {
	existingGroup := &Group{
		ID:                  1,
		Name:                "existing-group",
		Platform:            PlatformGrok,
		Status:              StatusActive,
		VideoRateMultiplier: 1,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder
	negative := -0.1

	_, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		VideoRateMultiplier: &negative,
placeholder)
placeholder
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_UpdateGroup_InvalidatesAuthCacheOnRPMLimitChange(t *testing.T) {
	existingGroup := &Group{
		ID:       1,
		Name:     "existing-group",
		Platform: PlatformAnthropic,
		Status:   StatusActive,
		RPMLimit: 10,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	invalidator := &authCacheInvalidatorStub{placeholder
	svc := &adminServiceImpl{
		groupRepo:            repo,
		authCacheInvalidator: invalidator,
placeholder

	rpmLimit := 60
	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		RPMLimit: &rpmLimit,
placeholder)
placeholder
	require.NotNil(t, group)
	require.Equal(t, 60, repo.updated.RPMLimit)
	require.Equal(t, []int64{1placeholder, invalidator.groupIDs, "分组 RPMLimit 写入 auth snapshot，变更后必须失效 API Key 认证缓存")
placeholder

func TestAdminService_UpdateGroup_ReasoningEffortMappingsTriState(t *testing.T) {
	tests := []struct {
		name  string
		input *UpdateGroupInput
		want  []ReasoningEffortMapping
placeholder{
		{
			name:  "nil preserves existing mappings",
			input: &UpdateGroupInput{placeholder,
			want:  []ReasoningEffortMapping{{From: "max", To: "xhigh"placeholderplaceholder,
	placeholder,
		{
			name: "empty array clears mappings",
			input: func() *UpdateGroupInput {
				empty := []ReasoningEffortMapping{placeholder
				return &UpdateGroupInput{ReasoningEffortMappings: &emptyplaceholder
		placeholder(),
			want: []ReasoningEffortMapping{placeholder,
	placeholder,
		{
			name: "non empty array replaces and canonicalizes mappings",
			input: func() *UpdateGroupInput {
				replacement := []ReasoningEffortMapping{{From: " X-HIGH ", To: " high "placeholderplaceholder
				return &UpdateGroupInput{ReasoningEffortMappings: &replacementplaceholder
		placeholder(),
			want: []ReasoningEffortMapping{{From: "xhigh", To: "high"placeholderplaceholder,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			existing := &Group{
				ID:                      1,
				Name:                    "openai-group",
				Platform:                PlatformOpenAI,
				Status:                  StatusActive,
				ReasoningEffortMappings: []ReasoningEffortMapping{{From: "max", To: "xhigh"placeholderplaceholder,
		placeholder
			repo := &groupRepoStubForAdmin{getByID: existingplaceholder
			svc := &adminServiceImpl{groupRepo: repoplaceholder

			_, err := svc.UpdateGroup(context.Background(), existing.ID, tt.input)

		placeholder
			require.Equal(t, tt.want, repo.updated.ReasoningEffortMappings)
	placeholder)
placeholder
placeholder

func TestAdminService_UpdateGroup_RejectsInvalidReasoningEffortMappings(t *testing.T) {
	existing := &Group{
		ID:               1,
		Name:             "openai",
		Platform:         PlatformOpenAI,
		SubscriptionType: SubscriptionTypeStandard,
		RateMultiplier:   1,
		Status:           StatusActive,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{groups: map[int64]*Group{existing.ID: existingplaceholderplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder
	invalid := []ReasoningEffortMapping{
		{From: "max", To: "xhigh"placeholder,
		{From: " MAX ", To: "high"placeholder,
placeholder

	_, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		ReasoningEffortMappings: &invalid,
placeholder)

placeholder
	require.Contains(t, err.Error(), "duplicate reasoning effort mapping source")
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_UpdateGroup_ClearsReasoningPolicyForUnsupportedPlatform(t *testing.T) {
	existing := &Group{
		ID:                      1,
		Name:                    "openai-group",
		Platform:                PlatformOpenAI,
		Status:                  StatusActive,
		MaxReasoningEffort:      "medium",
		ReasoningEffortMappings: []ReasoningEffortMapping{{From: "max", To: "xhigh"placeholderplaceholder,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{Platform: PlatformAnthropicplaceholder)

placeholder
	require.Empty(t, repo.updated.MaxReasoningEffort)
	require.Empty(t, repo.updated.ReasoningEffortMappings)
placeholder

func TestAdminService_UpdateGroup_ClearsPeakRateWhenChangingToStandard(t *testing.T) {
	existingGroup := &Group{
		ID:                 1,
		Name:               "existing-group",
		Platform:           PlatformOpenAI,
		Status:             StatusActive,
		SubscriptionType:   SubscriptionTypeSubscription,
		PeakRateEnabled:    true,
		PeakStart:          "14:00",
		PeakEnd:            "18:00",
		PeakRateMultiplier: 3,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		SubscriptionType: SubscriptionTypeStandard,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, SubscriptionTypeStandard, repo.updated.SubscriptionType)
	require.False(t, repo.updated.PeakRateEnabled)
	require.Equal(t, "", repo.updated.PeakStart)
	require.Equal(t, "", repo.updated.PeakEnd)
	require.Equal(t, 1.0, repo.updated.PeakRateMultiplier)
placeholder

func TestAdminService_CreateGroup_NormalizesMessagesDispatchModelConfig(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:           "dispatch-group",
		Description:    "dispatch config",
		Platform:       PlatformOpenAI,
		RateMultiplier: 1.0,
		MessagesDispatchModelConfig: OpenAIMessagesDispatchModelConfig{
			OpusMappedModel:   " gpt-5.4-high ",
			SonnetMappedModel: " gpt-5.3-codex ",
			HaikuMappedModel:  " gpt-5.4-mini-medium ",
			ExactModelMappings: map[string]string{
				" claude-sonnet-4-5-20250929 ": " gpt-5.2-high ",
		placeholder,
	placeholder,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.Equal(t, OpenAIMessagesDispatchModelConfig{
		OpusMappedModel:   "gpt-5.4",
		SonnetMappedModel: "gpt-5.3-codex",
		HaikuMappedModel:  "gpt-5.4-mini",
		ExactModelMappings: map[string]string{
			"claude-sonnet-4-5-20250929": "gpt-5.2",
	placeholder,
placeholder, repo.created.MessagesDispatchModelConfig)
placeholder

func TestAdminService_UpdateGroup_NormalizesMessagesDispatchModelConfig(t *testing.T) {
	existingGroup := &Group{
		ID:       1,
		Name:     "existing-group",
		Platform: PlatformOpenAI,
		Status:   StatusActive,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		MessagesDispatchModelConfig: &OpenAIMessagesDispatchModelConfig{
			SonnetMappedModel: " gpt-5.4-medium ",
			ExactModelMappings: map[string]string{
				" placeholder ": " gpt-5.4-mini-high ",
		placeholder,
	placeholder,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, OpenAIMessagesDispatchModelConfig{
		SonnetMappedModel: "gpt-5.4",
		ExactModelMappings: map[string]string{
			"placeholder": "gpt-5.4-mini",
	placeholder,
placeholder, repo.updated.MessagesDispatchModelConfig)
placeholder

func TestAdminService_CreateGroup_ClearsMessagesDispatchFieldsForNonOpenAIPlatform(t *testing.T) {
	repo := &groupRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                  "anthropic-group",
		Description:           "non-openai",
		Platform:              PlatformAnthropic,
		RateMultiplier:        1.0,
		AllowMessagesDispatch: true,
		DefaultMappedModel:    "gpt-5.4",
		MessagesDispatchModelConfig: OpenAIMessagesDispatchModelConfig{
			OpusMappedModel: "gpt-5.4",
	placeholder,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.False(t, repo.created.AllowMessagesDispatch)
	require.Empty(t, repo.created.DefaultMappedModel)
	require.Equal(t, OpenAIMessagesDispatchModelConfig{placeholder, repo.created.MessagesDispatchModelConfig)
placeholder

func TestAdminService_UpdateGroup_ClearsMessagesDispatchFieldsWhenPlatformChangesAwayFromOpenAI(t *testing.T) {
	existingGroup := &Group{
		ID:                    1,
		Name:                  "existing-openai-group",
		Platform:              PlatformOpenAI,
		Status:                StatusActive,
		AllowMessagesDispatch: true,
		DefaultMappedModel:    "gpt-5.4",
		MessagesDispatchModelConfig: OpenAIMessagesDispatchModelConfig{
			SonnetMappedModel: "gpt-5.3-codex",
	placeholder,
placeholder
	repo := &groupRepoStubForAdmin{getByID: existingGroupplaceholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), 1, &UpdateGroupInput{
		Platform: PlatformAnthropic,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, PlatformAnthropic, repo.updated.Platform)
	require.False(t, repo.updated.AllowMessagesDispatch)
	require.Empty(t, repo.updated.DefaultMappedModel)
	require.Equal(t, OpenAIMessagesDispatchModelConfig{placeholder, repo.updated.MessagesDispatchModelConfig)
placeholder

func TestAdminService_ListGroups_WithSearch(t *testing.T) {
	// 测试：
	// 1. search 参数正常传递到 repository 层
	// 2. search 为空字符串时的行为
	// 3. search 与其他过滤条件组合使用

	t.Run("search 参数正常传递到 repository 层", func(t *testing.T) {
		repo := &groupRepoStubForAdmin{
			listWithFiltersGroups: []Group{{ID: 1, Name: "alpha"placeholderplaceholder,
			listWithFiltersResult: &pagination.PaginationResult{Total: 1placeholder,
	placeholder
		svc := &adminServiceImpl{groupRepo: repoplaceholder

		groups, total, err := svc.ListGroups(context.Background(), 1, 20, "", "", "alpha", nil, "", "")
	placeholder
		require.Equal(t, int64(1), total)
		require.Equal(t, []Group{{ID: 1, Name: "alpha"placeholderplaceholder, groups)

		require.Equal(t, 1, repo.listWithFiltersCalls)
		require.Equal(t, pagination.PaginationParams{Page: 1, PageSize: 20placeholder, repo.listWithFiltersParams)
		require.Equal(t, "alpha", repo.listWithFiltersSearch)
		require.Nil(t, repo.listWithFiltersIsExclusive)
placeholder)

	t.Run("search 为空字符串时传递空字符串", func(t *testing.T) {
		repo := &groupRepoStubForAdmin{
			listWithFiltersGroups: []Group{placeholder,
			listWithFiltersResult: &pagination.PaginationResult{Total: 0placeholder,
	placeholder
		svc := &adminServiceImpl{groupRepo: repoplaceholder

		groups, total, err := svc.ListGroups(context.Background(), 2, 10, "", "", "", nil, "", "")
	placeholder
		require.Empty(t, groups)
		require.Equal(t, int64(0), total)

		require.Equal(t, 1, repo.listWithFiltersCalls)
		require.Equal(t, pagination.PaginationParams{Page: 2, PageSize: 10placeholder, repo.listWithFiltersParams)
		require.Equal(t, "", repo.listWithFiltersSearch)
		require.Nil(t, repo.listWithFiltersIsExclusive)
placeholder)

	t.Run("search 与其他过滤条件组合使用", func(t *testing.T) {
		isExclusive := true
		repo := &groupRepoStubForAdmin{
			listWithFiltersGroups: []Group{{ID: 2, Name: "beta"placeholderplaceholder,
			listWithFiltersResult: &pagination.PaginationResult{Total: 42placeholder,
	placeholder
		svc := &adminServiceImpl{groupRepo: repoplaceholder

		groups, total, err := svc.ListGroups(context.Background(), 3, 50, PlatformAntigravity, StatusActive, "beta", &isExclusive, "", "")
	placeholder
		require.Equal(t, int64(42), total)
		require.Equal(t, []Group{{ID: 2, Name: "beta"placeholderplaceholder, groups)

		require.Equal(t, 1, repo.listWithFiltersCalls)
		require.Equal(t, pagination.PaginationParams{Page: 3, PageSize: 50placeholder, repo.listWithFiltersParams)
		require.Equal(t, PlatformAntigravity, repo.listWithFiltersPlatform)
		require.Equal(t, StatusActive, repo.listWithFiltersStatus)
		require.Equal(t, "beta", repo.listWithFiltersSearch)
		require.NotNil(t, repo.listWithFiltersIsExclusive)
		require.True(t, *repo.listWithFiltersIsExclusive)
placeholder)
placeholder

func TestAdminService_ValidateFallbackGroup_DetectsCycle(t *testing.T) {
	groupID := int64(1)
	fallbackID := int64(2)
	repo := &groupRepoStubForFallbackCycle{
		groups: map[int64]*Group{
			groupID: {
				ID:              groupID,
				FallbackGroupID: &fallbackID,
		placeholder,
			fallbackID: {
				ID:              fallbackID,
				FallbackGroupID: &groupID,
		placeholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	err := svc.validateFallbackGroup(context.Background(), groupID, fallbackID)
placeholder
	require.Contains(t, err.Error(), "fallback group cycle")
placeholder

type groupRepoStubForFallbackCycle struct {
	groups map[int64]*Group
placeholder

func (s *groupRepoStubForFallbackCycle) Create(_ context.Context, _ *Group) error {
	panic("unexpected Create call")
placeholder

func (s *groupRepoStubForFallbackCycle) Update(_ context.Context, _ *Group) error {
	panic("unexpected Update call")
placeholder

func (s *groupRepoStubForFallbackCycle) GetByID(ctx context.Context, id int64) (*Group, error) {
	return s.GetByIDLite(ctx, id)
placeholder

func (s *groupRepoStubForFallbackCycle) GetByIDLite(_ context.Context, id int64) (*Group, error) {
	if g, ok := s.groups[id]; ok {
		return g, nil
placeholder
	return nil, ErrGroupNotFound
placeholder

func (s *groupRepoStubForFallbackCycle) Delete(_ context.Context, _ int64) error {
	panic("unexpected Delete call")
placeholder

func (s *groupRepoStubForFallbackCycle) DeleteCascade(_ context.Context, _ int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
placeholder

func (s *groupRepoStubForFallbackCycle) List(_ context.Context, _ pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *groupRepoStubForFallbackCycle) ListWithFilters(_ context.Context, _ pagination.PaginationParams, _, _, _ string, _ *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *groupRepoStubForFallbackCycle) ListActive(_ context.Context) ([]Group, error) {
	panic("unexpected ListActive call")
placeholder

func (s *groupRepoStubForFallbackCycle) ListActiveByPlatform(_ context.Context, _ string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
placeholder

func (s *groupRepoStubForFallbackCycle) ExistsByName(_ context.Context, _ string) (bool, error) {
	panic("unexpected ExistsByName call")
placeholder

func (s *groupRepoStubForFallbackCycle) GetAccountCount(_ context.Context, _ int64) (int64, int64, error) {
	panic("unexpected GetAccountCount call")
placeholder

func (s *groupRepoStubForFallbackCycle) DeleteAccountGroupsByGroupID(_ context.Context, _ int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
placeholder

func (s *groupRepoStubForFallbackCycle) BindAccountsToGroup(_ context.Context, _ int64, _ []int64) error {
	panic("unexpected BindAccountsToGroup call")
placeholder

func (s *groupRepoStubForFallbackCycle) GetAccountIDsByGroupIDs(_ context.Context, _ []int64) ([]int64, error) {
	panic("unexpected GetAccountIDsByGroupIDs call")
placeholder

func (s *groupRepoStubForFallbackCycle) UpdateSortOrders(_ context.Context, _ []GroupSortOrderUpdate) error {
	return nil
placeholder

type groupRepoStubForInvalidRequestFallback struct {
	groups  map[int64]*Group
	created *Group
	updated *Group
placeholder

func (s *groupRepoStubForInvalidRequestFallback) Create(_ context.Context, g *Group) error {
	s.created = g
	return nil
placeholder

func (s *groupRepoStubForInvalidRequestFallback) Update(_ context.Context, g *Group) error {
	s.updated = g
	return nil
placeholder

func (s *groupRepoStubForInvalidRequestFallback) GetByID(ctx context.Context, id int64) (*Group, error) {
	return s.GetByIDLite(ctx, id)
placeholder

func (s *groupRepoStubForInvalidRequestFallback) GetByIDLite(_ context.Context, id int64) (*Group, error) {
	if g, ok := s.groups[id]; ok {
		return g, nil
placeholder
	return nil, ErrGroupNotFound
placeholder

func (s *groupRepoStubForInvalidRequestFallback) Delete(_ context.Context, _ int64) error {
	panic("unexpected Delete call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) DeleteCascade(_ context.Context, _ int64) ([]int64, error) {
	panic("unexpected DeleteCascade call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) List(_ context.Context, _ pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected List call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) ListWithFilters(_ context.Context, _ pagination.PaginationParams, _, _, _ string, _ *bool) ([]Group, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) ListActive(_ context.Context) ([]Group, error) {
	panic("unexpected ListActive call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) ListActiveByPlatform(_ context.Context, _ string) ([]Group, error) {
	panic("unexpected ListActiveByPlatform call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) ExistsByName(_ context.Context, _ string) (bool, error) {
	panic("unexpected ExistsByName call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) GetAccountCount(_ context.Context, _ int64) (int64, int64, error) {
	panic("unexpected GetAccountCount call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) DeleteAccountGroupsByGroupID(_ context.Context, _ int64) (int64, error) {
	panic("unexpected DeleteAccountGroupsByGroupID call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) GetAccountIDsByGroupIDs(_ context.Context, _ []int64) ([]int64, error) {
	panic("unexpected GetAccountIDsByGroupIDs call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) BindAccountsToGroup(_ context.Context, _ int64, _ []int64) error {
	panic("unexpected BindAccountsToGroup call")
placeholder

func (s *groupRepoStubForInvalidRequestFallback) UpdateSortOrders(_ context.Context, _ []GroupSortOrderUpdate) error {
	return nil
placeholder

func TestAdminService_CreateGroup_InvalidRequestFallbackRejectsUnsupportedPlatform(t *testing.T) {
	fallbackID := int64(10)
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			fallbackID: {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                            "g1",
		Platform:                        PlatformOpenAI,
		RateMultiplier:                  1.0,
		SubscriptionType:                SubscriptionTypeStandard,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.Contains(t, err.Error(), "invalid request fallback only supported for anthropic or antigravity groups")
	require.Nil(t, repo.created)
placeholder

func TestAdminService_CreateGroup_InvalidRequestFallbackRejectsSubscription(t *testing.T) {
	fallbackID := int64(10)
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			fallbackID: {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                            "g1",
		Platform:                        PlatformAnthropic,
		RateMultiplier:                  1.0,
		SubscriptionType:                SubscriptionTypeSubscription,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.Contains(t, err.Error(), "subscription groups cannot set invalid request fallback")
	require.Nil(t, repo.created)
placeholder

func TestAdminService_CreateGroup_InvalidRequestFallbackRejectsFallbackGroup(t *testing.T) {
	tests := []struct {
		name        string
		fallback    *Group
		wantMessage string
placeholder{
		{
			name:        "openai_target",
			fallback:    &Group{ID: 10, Platform: PlatformOpenAI, SubscriptionType: SubscriptionTypeStandardplaceholder,
			wantMessage: "fallback group must be anthropic platform",
	placeholder,
		{
			name:        "antigravity_target",
			fallback:    &Group{ID: 10, Platform: PlatformAntigravity, SubscriptionType: SubscriptionTypeStandardplaceholder,
			wantMessage: "fallback group must be anthropic platform",
	placeholder,
		{
			name:        "subscription_group",
			fallback:    &Group{ID: 10, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeSubscriptionplaceholder,
			wantMessage: "fallback group cannot be subscription type",
	placeholder,
		{
			name: "nested_fallback",
			fallback: &Group{
				ID:                              10,
				Platform:                        PlatformAnthropic,
				SubscriptionType:                SubscriptionTypeStandard,
				FallbackGroupIDOnInvalidRequest: func() *int64 { v := int64(99); return &v placeholder(),
		placeholder,
			wantMessage: "fallback group cannot have invalid request fallback configured",
	placeholder,
placeholder

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fallbackID := tc.fallback.ID
			repo := &groupRepoStubForInvalidRequestFallback{
				groups: map[int64]*Group{
					fallbackID: tc.fallback,
			placeholder,
		placeholder
			svc := &adminServiceImpl{groupRepo: repoplaceholder

			_, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
				Name:                            "g1",
				Platform:                        PlatformAnthropic,
				RateMultiplier:                  1.0,
				SubscriptionType:                SubscriptionTypeStandard,
				FallbackGroupIDOnInvalidRequest: &fallbackID,
		placeholder)
		placeholder
			require.Contains(t, err.Error(), tc.wantMessage)
			require.Nil(t, repo.created)
	placeholder)
placeholder
placeholder

func TestAdminService_CreateGroup_InvalidRequestFallbackNotFound(t *testing.T) {
	fallbackID := int64(10)
	repo := &groupRepoStubForInvalidRequestFallback{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                            "g1",
		Platform:                        PlatformAnthropic,
		RateMultiplier:                  1.0,
		SubscriptionType:                SubscriptionTypeStandard,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.Contains(t, err.Error(), "fallback group not found")
	require.Nil(t, repo.created)
placeholder

func TestAdminService_CreateGroup_InvalidRequestFallbackAllowsAntigravity(t *testing.T) {
	fallbackID := int64(10)
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			fallbackID: {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                            "g1",
		Platform:                        PlatformAntigravity,
		RateMultiplier:                  1.0,
		SubscriptionType:                SubscriptionTypeStandard,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.Equal(t, fallbackID, *repo.created.FallbackGroupIDOnInvalidRequest)
placeholder

func TestAdminService_CreateGroup_InvalidRequestFallbackClearsOnZero(t *testing.T) {
	zero := int64(0)
	repo := &groupRepoStubForInvalidRequestFallback{placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.CreateGroup(context.Background(), &CreateGroupInput{
		Name:                            "g1",
		Platform:                        PlatformAnthropic,
		RateMultiplier:                  1.0,
		SubscriptionType:                SubscriptionTypeStandard,
		FallbackGroupIDOnInvalidRequest: &zero,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.created)
	require.Nil(t, repo.created.FallbackGroupIDOnInvalidRequest)
placeholder

func TestAdminService_UpdateGroup_InvalidRequestFallbackPlatformMismatch(t *testing.T) {
	fallbackID := int64(10)
	existing := &Group{
		ID:                              1,
		Name:                            "g1",
		Platform:                        PlatformAnthropic,
		SubscriptionType:                SubscriptionTypeStandard,
		Status:                          StatusActive,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			existing.ID: existing,
			fallbackID:  {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		Platform: PlatformOpenAI,
placeholder)
placeholder
	require.Contains(t, err.Error(), "invalid request fallback only supported for anthropic or antigravity groups")
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_UpdateGroup_InvalidRequestFallbackSubscriptionMismatch(t *testing.T) {
	fallbackID := int64(10)
	existing := &Group{
		ID:                              1,
		Name:                            "g1",
		Platform:                        PlatformAnthropic,
		SubscriptionType:                SubscriptionTypeStandard,
		Status:                          StatusActive,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			existing.ID: existing,
			fallbackID:  {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		SubscriptionType: SubscriptionTypeSubscription,
placeholder)
placeholder
	require.Contains(t, err.Error(), "subscription groups cannot set invalid request fallback")
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_UpdateGroup_InvalidRequestFallbackClearsOnZero(t *testing.T) {
	fallbackID := int64(10)
	existing := &Group{
		ID:                              1,
		Name:                            "g1",
		Platform:                        PlatformAnthropic,
		SubscriptionType:                SubscriptionTypeStandard,
		Status:                          StatusActive,
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			existing.ID: existing,
			fallbackID:  {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	clear := int64(0)
	group, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		Platform:                        PlatformOpenAI,
		FallbackGroupIDOnInvalidRequest: &clear,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Nil(t, repo.updated.FallbackGroupIDOnInvalidRequest)
placeholder

func TestAdminService_UpdateGroup_InvalidRequestFallbackRejectsFallbackGroup(t *testing.T) {
	fallbackID := int64(10)
	existing := &Group{
		ID:               1,
		Name:             "g1",
		Platform:         PlatformAnthropic,
		SubscriptionType: SubscriptionTypeStandard,
		Status:           StatusActive,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			existing.ID: existing,
			fallbackID:  {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeSubscriptionplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	_, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.Contains(t, err.Error(), "fallback group cannot be subscription type")
	require.Nil(t, repo.updated)
placeholder

func TestAdminService_UpdateGroup_InvalidRequestFallbackSetSuccess(t *testing.T) {
	fallbackID := int64(10)
	existing := &Group{
		ID:               1,
		Name:             "g1",
		Platform:         PlatformAnthropic,
		SubscriptionType: SubscriptionTypeStandard,
		Status:           StatusActive,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			existing.ID: existing,
			fallbackID:  {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, fallbackID, *repo.updated.FallbackGroupIDOnInvalidRequest)
placeholder

func TestAdminService_UpdateGroup_InvalidRequestFallbackAllowsAntigravity(t *testing.T) {
	fallbackID := int64(10)
	existing := &Group{
		ID:               1,
		Name:             "g1",
		Platform:         PlatformAntigravity,
		SubscriptionType: SubscriptionTypeStandard,
		Status:           StatusActive,
placeholder
	repo := &groupRepoStubForInvalidRequestFallback{
		groups: map[int64]*Group{
			existing.ID: existing,
			fallbackID:  {ID: fallbackID, Platform: PlatformAnthropic, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: repoplaceholder

	group, err := svc.UpdateGroup(context.Background(), existing.ID, &UpdateGroupInput{
		FallbackGroupIDOnInvalidRequest: &fallbackID,
placeholder)
placeholder
	require.NotNil(t, group)
	require.NotNil(t, repo.updated)
	require.Equal(t, fallbackID, *repo.updated.FallbackGroupIDOnInvalidRequest)
placeholder

func TestAdminService_CreateCompositeRoute_RejectsNonCompositeGroup(t *testing.T) {
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{ID: 7, Platform: PlatformOpenAIplaceholder,
placeholder
	routeRepo := &compositeRouteRepoStubForAdmin{placeholder
	svc := &adminServiceImpl{groupRepo: groupRepo, compositeRouteRepo: routeRepoplaceholder

	_, err := svc.CreateCompositeRoute(context.Background(), 7, CompositeRouteInput{
		PublicModel:    "router/gpt-5",
		TargetPlatform: PlatformOpenAI,
		Enabled:        true,
placeholder)

placeholder
	require.ErrorContains(t, err, "not a composite group")
	require.Nil(t, routeRepo.created)
placeholder

func TestAdminService_CreateCompositeRoute_NormalizesAndPersists(t *testing.T) {
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{ID: 7, Platform: PlatformCompositeplaceholder,
placeholder
	routeRepo := &compositeRouteRepoStubForAdmin{nextID: 99placeholder
	svc := &adminServiceImpl{groupRepo: groupRepo, compositeRouteRepo: routeRepoplaceholder

	route, err := svc.CreateCompositeRoute(context.Background(), 7, CompositeRouteInput{
		PublicModel:    " router/gpt- ",
		MatchType:      CompositeRouteMatchPrefix,
		TargetPlatform: PlatformOpenAI,
		Endpoint:       CompositeRouteEndpointResponses,
		Enabled:        true,
		Notes:          " route note ",
placeholder)

placeholder
	require.NotNil(t, route)
	require.Equal(t, int64(99), route.ID)
	require.Equal(t, "router/gpt-", route.PublicModel)
	require.Equal(t, CompositeRouteMatchPrefix, route.MatchType)
	require.Equal(t, PlatformOpenAI, route.TargetPlatform)
	require.Equal(t, "router/gpt-", route.UpstreamModel)
	require.Equal(t, CompositeRouteEndpointResponses, route.Endpoint)
	require.Equal(t, 100, route.Priority)
	require.True(t, route.Enabled)
	require.Equal(t, "route note", route.Notes)
	require.Equal(t, route, routeRepo.created)
placeholder

func TestAdminService_UpdateAndDeleteCompositeRouteRequireRouteOwnership(t *testing.T) {
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{ID: 7, Platform: PlatformCompositeplaceholder,
placeholder
	routeRepo := &compositeRouteRepoStubForAdmin{
		routes: []CompositeModelRoute{
			{ID: 11, GroupID: 7, PublicModel: "router/gpt-5", TargetPlatform: PlatformOpenAI, Enabled: trueplaceholder,
			{ID: 12, GroupID: 8, PublicModel: "router/other", TargetPlatform: PlatformGemini, Enabled: trueplaceholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: groupRepo, compositeRouteRepo: routeRepoplaceholder

	updated, err := svc.UpdateCompositeRoute(context.Background(), 7, 11, CompositeRouteInput{
		PublicModel:    "router/gpt-5",
		TargetPlatform: PlatformGemini,
		UpstreamModel:  "gemini-2.5-pro",
		Endpoint:       CompositeRouteEndpointChatCompletions,
		Priority:       3,
		Enabled:        true,
placeholder)
placeholder
	require.Equal(t, int64(11), updated.ID)
	require.Equal(t, PlatformGemini, updated.TargetPlatform)
	require.Equal(t, "gemini-2.5-pro", updated.UpstreamModel)
	require.Equal(t, updated, routeRepo.updated)

	err = svc.DeleteCompositeRoute(context.Background(), 7, 12)
	require.ErrorIs(t, err, ErrCompositeRouteNotFound)
	require.Empty(t, routeRepo.deleted)

	err = svc.DeleteCompositeRoute(context.Background(), 7, 11)
placeholder
	require.Equal(t, []int64{11placeholder, routeRepo.deleted)
placeholder

func TestAdminService_PreviewCompositeRouteUsesExplicitRoutes(t *testing.T) {
	groupRepo := &groupRepoStubForAdmin{
		getByID: &Group{ID: 7, Platform: PlatformCompositeplaceholder,
placeholder
	routeRepo := &compositeRouteRepoStubForAdmin{
		routes: []CompositeModelRoute{
			{
				ID:             11,
				GroupID:        7,
				PublicModel:    "openrouter/claude",
				MatchType:      CompositeRouteMatchExact,
				TargetPlatform: PlatformAnthropic,
				UpstreamModel:  "claude-sonnet-4-6",
				Endpoint:       CompositeRouteEndpointMessages,
				Priority:       100,
				Enabled:        true,
		placeholder,
	placeholder,
placeholder
	svc := &adminServiceImpl{groupRepo: groupRepo, compositeRouteRepo: routeRepoplaceholder

	decision, err := svc.PreviewCompositeRoute(context.Background(), 7, CompositeRoutePreviewRequest{
		Model:    "openrouter/claude",
		Endpoint: CompositeRouteEndpointMessages,
placeholder)

placeholder
	require.NotNil(t, decision)
	require.True(t, decision.Matched)
	require.Equal(t, CompositeRouteSourceExplicit, decision.Source)
	require.Equal(t, PlatformAnthropic, decision.TargetPlatform)
	require.Equal(t, "claude-sonnet-4-6", decision.UpstreamModel)
	require.NotNil(t, decision.Route)
	require.Equal(t, int64(11), decision.Route.ID)
placeholder
