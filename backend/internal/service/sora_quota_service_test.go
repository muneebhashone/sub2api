//go:build unit

package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

// ==================== Stub: GroupRepository (用于 SoraQuotaService) ====================

var _ GroupRepository = (*stubGroupRepoForQuota)(nil)

type stubGroupRepoForQuota struct {
	groups map[int64]*Group
placeholder

func newStubGroupRepoForQuota() *stubGroupRepoForQuota {
	return &stubGroupRepoForQuota{groups: make(map[int64]*Group)placeholder
placeholder

func (r *stubGroupRepoForQuota) GetByID(_ context.Context, id int64) (*Group, error) {
	if g, ok := r.groups[id]; ok {
		return g, nil
placeholder
	return nil, fmt.Errorf("group not found")
placeholder
func (r *stubGroupRepoForQuota) Create(context.Context, *Group) error { return nil placeholder
func (r *stubGroupRepoForQuota) GetByIDLite(_ context.Context, id int64) (*Group, error) {
	return r.GetByID(context.Background(), id)
placeholder
func (r *stubGroupRepoForQuota) Update(context.Context, *Group) error { return nil placeholder
func (r *stubGroupRepoForQuota) Delete(context.Context, int64) error  { return nil placeholder
func (r *stubGroupRepoForQuota) DeleteCascade(context.Context, int64) ([]int64, error) {
	return nil, nil
placeholder
func (r *stubGroupRepoForQuota) List(context.Context, pagination.PaginationParams) ([]Group, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (r *stubGroupRepoForQuota) ListWithFilters(context.Context, pagination.PaginationParams, string, string, string, *bool) ([]Group, *pagination.PaginationResult, error) {
	return nil, nil, nil
placeholder
func (r *stubGroupRepoForQuota) ListActive(context.Context) ([]Group, error) { return nil, nil placeholder
func (r *stubGroupRepoForQuota) ListActiveByPlatform(context.Context, string) ([]Group, error) {
	return nil, nil
placeholder
func (r *stubGroupRepoForQuota) ExistsByName(context.Context, string) (bool, error) {
	return false, nil
placeholder
func (r *stubGroupRepoForQuota) GetAccountCount(context.Context, int64) (int64, error) {
	return 0, nil
placeholder
func (r *stubGroupRepoForQuota) DeleteAccountGroupsByGroupID(context.Context, int64) (int64, error) {
	return 0, nil
placeholder
func (r *stubGroupRepoForQuota) GetAccountIDsByGroupIDs(context.Context, []int64) ([]int64, error) {
	return nil, nil
placeholder
func (r *stubGroupRepoForQuota) BindAccountsToGroup(context.Context, int64, []int64) error {
	return nil
placeholder
func (r *stubGroupRepoForQuota) UpdateSortOrders(context.Context, []GroupSortOrderUpdate) error {
	return nil
placeholder

// ==================== Stub: SettingRepository (用于 SettingService) ====================

var _ SettingRepository = (*stubSettingRepoForQuota)(nil)

type stubSettingRepoForQuota struct {
	values map[string]string
placeholder

func newStubSettingRepoForQuota(values map[string]string) *stubSettingRepoForQuota {
	if values == nil {
		values = make(map[string]string)
placeholder
	return &stubSettingRepoForQuota{values: valuesplaceholder
placeholder

func (r *stubSettingRepoForQuota) Get(_ context.Context, key string) (*Setting, error) {
	if v, ok := r.values[key]; ok {
		return &Setting{Key: key, Value: vplaceholder, nil
placeholder
	return nil, ErrSettingNotFound
placeholder
func (r *stubSettingRepoForQuota) GetValue(_ context.Context, key string) (string, error) {
	if v, ok := r.values[key]; ok {
		return v, nil
placeholder
	return "", ErrSettingNotFound
placeholder
func (r *stubSettingRepoForQuota) Set(_ context.Context, key, value string) error {
	r.values[key] = value
	return nil
placeholder
func (r *stubSettingRepoForQuota) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, k := range keys {
		if v, ok := r.values[k]; ok {
			result[k] = v
	placeholder
placeholder
	return result, nil
placeholder
func (r *stubSettingRepoForQuota) SetMultiple(_ context.Context, settings map[string]string) error {
	for k, v := range settings {
		r.values[k] = v
placeholder
	return nil
placeholder
func (r *stubSettingRepoForQuota) GetAll(_ context.Context) (map[string]string, error) {
	return r.values, nil
placeholder
func (r *stubSettingRepoForQuota) Delete(_ context.Context, key string) error {
	delete(r.values, key)
	return nil
placeholder

// ==================== GetQuota ====================

func TestGetQuota_UserLevel(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 10 * 1024 * 1024, // 10MB
		SoraStorageUsedBytes:  3 * 1024 * 1024,  // 3MB
placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, int64(10*1024*1024), quota.QuotaBytes)
	require.Equal(t, int64(3*1024*1024), quota.UsedBytes)
	require.Equal(t, "user", quota.Source)
placeholder

func TestGetQuota_GroupLevel(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 0, // 用户级无配额
		SoraStorageUsedBytes:  1024,
		AllowedGroups:         []int64{10, 20placeholder,
placeholder

	groupRepo := newStubGroupRepoForQuota()
	groupRepo.groups[10] = &Group{ID: 10, SoraStorageQuotaBytes: 5 * 1024 * placeholder
	groupRepo.groups[20] = &Group{ID: 20, SoraStorageQuotaBytes: 20 * 1024 * placeholder

	svc := NewSoraQuotaService(userRepo, groupRepo, nil)
	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, int64(20*1024*1024), quota.QuotaBytes) // 取最大值
	require.Equal(t, "group", quota.Source)
placeholder

func TestGetQuota_SystemLevel(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageQuotaBytes: 0, SoraStorageUsedBytes: 512placeholder

	settingRepo := newStubSettingRepoForQuota(map[string]string{
		SettingKeySoraDefaultStorageQuotaBytes: "104857600", // 100MB
placeholder)
	settingService := NewSettingService(settingRepo, &config.Config{placeholder)
	svc := NewSoraQuotaService(userRepo, nil, settingService)

	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, int64(104857600), quota.QuotaBytes)
	require.Equal(t, "system", quota.Source)
placeholder

func TestGetQuota_NoLimit(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageQuotaBytes: 0, SoraStorageUsedBytes: 0placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, int64(0), quota.QuotaBytes)
	require.Equal(t, "unlimited", quota.Source)
placeholder

func TestGetQuota_UserNotFound(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	svc := NewSoraQuotaService(userRepo, nil, nil)

	_, err := svc.GetQuota(context.Background(), 999)
placeholder
	require.Contains(t, err.Error(), "get user")
placeholder

func TestGetQuota_GroupRepoError(t *testing.T) {
	// 分组获取失败时跳过该分组（不影响整体）
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID: 1, SoraStorageQuotaBytes: 0,
		AllowedGroups: []int64{999placeholder, // 不存在的分组
placeholder

	groupRepo := newStubGroupRepoForQuota()
	svc := NewSoraQuotaService(userRepo, groupRepo, nil)

	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, "unlimited", quota.Source) // 分组获取失败，回退到无限制
placeholder

// ==================== CheckQuota ====================

func TestCheckQuota_Sufficient(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 10 * 1024 * 1024,
		SoraStorageUsedBytes:  3 * 1024 * 1024,
placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.CheckQuota(context.Background(), 1, 1024)
placeholder
placeholder

func TestCheckQuota_Exceeded(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 10 * 1024 * 1024,
		SoraStorageUsedBytes:  10 * 1024 * 1024, // 已满
placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.CheckQuota(context.Background(), 1, 1)
placeholder
	require.Contains(t, err.Error(), "配额不足")
placeholder

func TestCheckQuota_NoLimit(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 0, // 无限制
		SoraStorageUsedBytes:  1000000000,
placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.CheckQuota(context.Background(), 1, 999999999)
placeholder // 无限制时始终通过
placeholder

func TestCheckQuota_ExactBoundary(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 1024,
		SoraStorageUsedBytes:  1024, // 恰好满
placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	// 额外 0 字节不超
	require.NoError(t, svc.CheckQuota(context.Background(), 1, 0))
	// 额外 1 字节超出
	require.Error(t, svc.CheckQuota(context.Background(), 1, 1))
placeholder

// ==================== AddUsage ====================

func TestAddUsage_Success(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.AddUsage(context.Background(), 1, 2048)
placeholder
	require.Equal(t, int64(3072), userRepo.users[1].SoraStorageUsedBytes)
placeholder

func TestAddUsage_ZeroBytes(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.AddUsage(context.Background(), 1, 0)
placeholder
	require.Equal(t, int64(1024), userRepo.users[1].SoraStorageUsedBytes) // 不变
placeholder

func TestAddUsage_NegativeBytes(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.AddUsage(context.Background(), 1, -100)
placeholder
	require.Equal(t, int64(1024), userRepo.users[1].SoraStorageUsedBytes) // 不变
placeholder

func TestAddUsage_UserNotFound(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.AddUsage(context.Background(), 999, 1024)
placeholder
placeholder

func TestAddUsage_UpdateError(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: 0placeholder
	userRepo.updateErr = fmt.Errorf("db error")
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.AddUsage(context.Background(), 1, 1024)
placeholder
	require.Contains(t, err.Error(), "update user quota usage")
placeholder

// ==================== ReleaseUsage ====================

func TestReleaseUsage_Success(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: 3072placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.ReleaseUsage(context.Background(), 1, 1024)
placeholder
	require.Equal(t, int64(2048), userRepo.users[1].SoraStorageUsedBytes)
placeholder

func TestReleaseUsage_ClampToZero(t *testing.T) {
	// 释放量大于已用量时，应 clamp 到 0
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: 500placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.ReleaseUsage(context.Background(), 1, 1000)
placeholder
	require.Equal(t, int64(0), userRepo.users[1].SoraStorageUsedBytes)
placeholder

func TestReleaseUsage_ZeroBytes(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.ReleaseUsage(context.Background(), 1, 0)
placeholder
	require.Equal(t, int64(1024), userRepo.users[1].SoraStorageUsedBytes) // 不变
placeholder

func TestReleaseUsage_NegativeBytes(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: placeholder
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.ReleaseUsage(context.Background(), 1, -50)
placeholder
	require.Equal(t, int64(1024), userRepo.users[1].SoraStorageUsedBytes) // 不变
placeholder

func TestReleaseUsage_UserNotFound(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.ReleaseUsage(context.Background(), 999, 1024)
placeholder
placeholder

func TestReleaseUsage_UpdateError(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageUsedBytes: placeholder
	userRepo.updateErr = fmt.Errorf("db error")
	svc := NewSoraQuotaService(userRepo, nil, nil)

	err := svc.ReleaseUsage(context.Background(), 1, 512)
placeholder
	require.Contains(t, err.Error(), "update user quota release")
placeholder

// ==================== GetQuotaFromSettings ====================

func TestGetQuotaFromSettings_NilSettingService(t *testing.T) {
	svc := NewSoraQuotaService(nil, nil, nil)
	require.Equal(t, int64(0), svc.GetQuotaFromSettings(context.Background()))
placeholder

func TestGetQuotaFromSettings_WithSettings(t *testing.T) {
	settingRepo := newStubSettingRepoForQuota(map[string]string{
		SettingKeySoraDefaultStorageQuotaBytes: "52428800", // 50MB
placeholder)
	settingService := NewSettingService(settingRepo, &config.Config{placeholder)
	svc := NewSoraQuotaService(nil, nil, settingService)

	require.Equal(t, int64(52428800), svc.GetQuotaFromSettings(context.Background()))
placeholder

// ==================== SetUserSoraQuota ====================

func TestSetUserSoraQuota_Success(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{ID: 1, SoraStorageQuotaBytes: 0placeholder

	err := SetUserSoraQuota(context.Background(), userRepo, 1, 10*1024*1024)
placeholder
	require.Equal(t, int64(10*1024*1024), userRepo.users[1].SoraStorageQuotaBytes)
placeholder

func TestSetUserSoraQuota_UserNotFound(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	err := SetUserSoraQuota(context.Background(), userRepo, 999, 1024)
placeholder
placeholder

// ==================== ParseQuotaBytes ====================

func TestParseQuotaBytes(t *testing.T) {
	require.Equal(t, int64(1048576), ParseQuotaBytes("1048576"))
	require.Equal(t, int64(0), ParseQuotaBytes(""))
	require.Equal(t, int64(0), ParseQuotaBytes("abc"))
	require.Equal(t, int64(-1), ParseQuotaBytes("-1"))
placeholder

// ==================== 优先级完整测试 ====================

func TestQuotaPriority_UserOverridesGroup(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 5 * 1024 * 1024,
		AllowedGroups:         []int64{10placeholder,
placeholder

	groupRepo := newStubGroupRepoForQuota()
	groupRepo.groups[10] = &Group{ID: 10, SoraStorageQuotaBytes: 20 * 1024 * placeholder

	svc := NewSoraQuotaService(userRepo, groupRepo, nil)
	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, "user", quota.Source) // 用户级优先
	require.Equal(t, int64(5*1024*1024), quota.QuotaBytes)
placeholder

func TestQuotaPriority_GroupOverridesSystem(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 0,
		AllowedGroups:         []int64{10placeholder,
placeholder

	groupRepo := newStubGroupRepoForQuota()
	groupRepo.groups[10] = &Group{ID: 10, SoraStorageQuotaBytes: 20 * 1024 * placeholder

	settingRepo := newStubSettingRepoForQuota(map[string]string{
		SettingKeySoraDefaultStorageQuotaBytes: "104857600", // 100MB
placeholder)
	settingService := NewSettingService(settingRepo, &config.Config{placeholder)

	svc := NewSoraQuotaService(userRepo, groupRepo, settingService)
	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, "group", quota.Source) // 分组级优先于系统
	require.Equal(t, int64(20*1024*1024), quota.QuotaBytes)
placeholder

func TestQuotaPriority_FallbackToSystem(t *testing.T) {
	userRepo := newStubUserRepoForQuota()
	userRepo.users[1] = &User{
		ID:                    1,
		SoraStorageQuotaBytes: 0,
		AllowedGroups:         []int64{10placeholder,
placeholder

	groupRepo := newStubGroupRepoForQuota()
	groupRepo.groups[10] = &Group{ID: 10, SoraStorageQuotaBytes: 0placeholder // 分组无配额

	settingRepo := newStubSettingRepoForQuota(map[string]string{
		SettingKeySoraDefaultStorageQuotaBytes: "52428800", // 50MB
placeholder)
	settingService := NewSettingService(settingRepo, &config.Config{placeholder)

	svc := NewSoraQuotaService(userRepo, groupRepo, settingService)
	quota, err := svc.GetQuota(context.Background(), 1)
placeholder
	require.Equal(t, "system", quota.Source)
	require.Equal(t, int64(52428800), quota.QuotaBytes)
placeholder
