//go:build unit

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// --- helpers ---

func testTimePtr(t time.Time) *time.Time { return &t placeholder

func makeAccWithLoad(id int64, priority int, loadRate int, lastUsed *time.Time, accType string) accountWithLoad {
	return accountWithLoad{
		account: &Account{
			ID:          id,
			Priority:    priority,
			LastUsedAt:  lastUsed,
			Type:        accType,
			Schedulable: true,
			Status:      StatusActive,
	placeholder,
		loadInfo: &AccountLoadInfo{
			AccountID:          id,
			CurrentConcurrency: 0,
			LoadRate:           loadRate,
	placeholder,
placeholder
placeholder

// --- sortAccountsByPriorityAndLastUsed ---

func TestSortAccountsByPriorityAndLastUsed_ByPriority(t *testing.T) {
	now := time.Now()
	accounts := []*Account{
		{ID: 1, Priority: 5, LastUsedAt: testTimePtr(now)placeholder,
		{ID: 2, Priority: 1, LastUsedAt: testTimePtr(now)placeholder,
		{ID: 3, Priority: 3, LastUsedAt: testTimePtr(now)placeholder,
placeholder
	sortAccountsByPriorityAndLastUsed(accounts, false)
	require.Equal(t, int64(2), accounts[0].ID, "优先级最低的排第一")
	require.Equal(t, int64(3), accounts[1].ID)
	require.Equal(t, int64(1), accounts[2].ID)
placeholder

func TestSortAccountsByPriorityAndLastUsed_SamePriorityByLastUsed(t *testing.T) {
	now := time.Now()
	accounts := []*Account{
		{ID: 1, Priority: 1, LastUsedAt: testTimePtr(now)placeholder,
		{ID: 2, Priority: 1, LastUsedAt: testTimePtr(now.Add(-1 * time.Hour))placeholder,
		{ID: 3, Priority: 1, LastUsedAt: nilplaceholder,
placeholder
	sortAccountsByPriorityAndLastUsed(accounts, false)
	require.Equal(t, int64(3), accounts[0].ID, "nil LastUsedAt 排最前")
	require.Equal(t, int64(2), accounts[1].ID, "更早使用的排前面")
	require.Equal(t, int64(1), accounts[2].ID)
placeholder

func TestSortAccountsByPriorityAndLastUsed_PreferOAuth(t *testing.T) {
	accounts := []*Account{
		{ID: 1, Priority: 1, LastUsedAt: nil, Type: AccountTypeAPIKeyplaceholder,
		{ID: 2, Priority: 1, LastUsedAt: nil, Type: AccountTypeOAuthplaceholder,
placeholder
	sortAccountsByPriorityAndLastUsed(accounts, true)
	require.Equal(t, int64(2), accounts[0].ID, "preferOAuth 时 OAuth 账号排前面")
placeholder

func TestSortAccountsByPriorityAndLastUsed_StableSort(t *testing.T) {
	accounts := []*Account{
		{ID: 1, Priority: 1, LastUsedAt: nil, Type: AccountTypeAPIKeyplaceholder,
		{ID: 2, Priority: 1, LastUsedAt: nil, Type: AccountTypeAPIKeyplaceholder,
		{ID: 3, Priority: 1, LastUsedAt: nil, Type: AccountTypeAPIKeyplaceholder,
placeholder
	sortAccountsByPriorityAndLastUsed(accounts, false)
	// 稳定排序：相同键值的元素保持原始顺序
	require.Equal(t, int64(1), accounts[0].ID)
	require.Equal(t, int64(2), accounts[1].ID)
	require.Equal(t, int64(3), accounts[2].ID)
placeholder

func TestSortAccountsByPriorityAndLastUsed_MixedPriorityAndTime(t *testing.T) {
	now := time.Now()
	accounts := []*Account{
		{ID: 1, Priority: 2, LastUsedAt: nilplaceholder,
		{ID: 2, Priority: 1, LastUsedAt: testTimePtr(now)placeholder,
		{ID: 3, Priority: 1, LastUsedAt: testTimePtr(now.Add(-1 * time.Hour))placeholder,
		{ID: 4, Priority: 2, LastUsedAt: testTimePtr(now.Add(-2 * time.Hour))placeholder,
placeholder
	sortAccountsByPriorityAndLastUsed(accounts, false)
	// 优先级1排前：nil < earlier
	require.Equal(t, int64(3), accounts[0].ID, "优先级1 + 更早")
	require.Equal(t, int64(2), accounts[1].ID, "优先级1 + 现在")
	// 优先级2排后：nil < time
	require.Equal(t, int64(1), accounts[2].ID, "优先级2 + nil")
	require.Equal(t, int64(4), accounts[3].ID, "优先级2 + 有时间")
placeholder

// --- selectByCallCount ---

func TestSelectByCallCount_Empty(t *testing.T) {
	result := selectByCallCount(nil, nil, false)
	require.Nil(t, result)
placeholder

func TestSelectByCallCount_Single(t *testing.T) {
	accounts := []accountWithLoad{
		makeAccWithLoad(1, 1, 50, nil, AccountTypeAPIKey),
placeholder
	result := selectByCallCount(accounts, map[int64]*ModelLoadInfo{1: {CallCount: 10placeholderplaceholder, false)
	require.NotNil(t, result)
	require.Equal(t, int64(1), result.account.ID)
placeholder

func TestSelectByCallCount_NilModelLoadFallsBackToLRU(t *testing.T) {
	now := time.Now()
	accounts := []accountWithLoad{
		makeAccWithLoad(1, 1, 50, testTimePtr(now), AccountTypeAPIKey),
		makeAccWithLoad(2, 1, 50, testTimePtr(now.Add(-1*time.Hour)), AccountTypeAPIKey),
placeholder
	result := selectByCallCount(accounts, nil, false)
	require.NotNil(t, result)
	require.Equal(t, int64(2), result.account.ID, "nil modelLoadMap 应回退到 LRU 选择")
placeholder

func TestSelectByCallCount_SelectsMinCallCount(t *testing.T) {
	accounts := []accountWithLoad{
		makeAccWithLoad(1, 1, 50, nil, AccountTypeAPIKey),
		makeAccWithLoad(2, 1, 50, nil, AccountTypeAPIKey),
		makeAccWithLoad(3, 1, 50, nil, AccountTypeAPIKey),
placeholder
	modelLoad := map[int64]*ModelLoadInfo{
		1: {CallCount: 100placeholder,
		2: {CallCount: 5placeholder,
		3: {CallCount: 50placeholder,
placeholder
	// 运行多次确认总是选调用次数最少的
	for i := 0; i < 10; i++ {
		result := selectByCallCount(accounts, modelLoad, false)
		require.NotNil(t, result)
		require.Equal(t, int64(2), result.account.ID, "应选择调用次数最少的账号")
placeholder
placeholder

func TestSelectByCallCount_NewAccountUsesAverage(t *testing.T) {
	accounts := []accountWithLoad{
		makeAccWithLoad(1, 1, 50, nil, AccountTypeAPIKey),
		makeAccWithLoad(2, 1, 50, nil, AccountTypeAPIKey),
		makeAccWithLoad(3, 1, 50, nil, AccountTypeAPIKey),
placeholder
	// 账号1和2有调用记录，账号3是新账号（CallCount=0）
	// 平均调用次数 = (100 + 200) / 2 = 150
	// 新账号用平均值 150，比账号1(100)多，所以应选账号1
	modelLoad := map[int64]*ModelLoadInfo{
		1: {CallCount: 100placeholder,
		2: {CallCount: 200placeholder,
		// 3 没有记录
placeholder
	for i := 0; i < 10; i++ {
		result := selectByCallCount(accounts, modelLoad, false)
		require.NotNil(t, result)
		require.Equal(t, int64(1), result.account.ID, "新账号虚拟调用次数(150)高于账号1(100)，应选账号1")
placeholder
placeholder

func TestSelectByCallCount_AllNewAccountsFallToAvgZero(t *testing.T) {
	accounts := []accountWithLoad{
		makeAccWithLoad(1, 1, 50, nil, AccountTypeAPIKey),
		makeAccWithLoad(2, 1, 50, nil, AccountTypeAPIKey),
placeholder
	// 所有账号都是新的，avgCallCount = 0，所有人 effectiveCallCount 都是 0
	modelLoad := map[int64]*ModelLoadInfo{placeholder
	validIDs := map[int64]bool{1: true, 2: trueplaceholder
	for i := 0; i < 10; i++ {
		result := selectByCallCount(accounts, modelLoad, false)
		require.NotNil(t, result)
		require.True(t, validIDs[result.account.ID], "所有新账号应随机选择")
placeholder
placeholder

func TestSelectByCallCount_PreferOAuth(t *testing.T) {
	accounts := []accountWithLoad{
		makeAccWithLoad(1, 1, 50, nil, AccountTypeAPIKey),
		makeAccWithLoad(2, 1, 50, nil, AccountTypeOAuth),
placeholder
	// 两个账号调用次数相同
	modelLoad := map[int64]*ModelLoadInfo{
		1: {CallCount: 10placeholder,
		2: {CallCount: 10placeholder,
placeholder
	for i := 0; i < 10; i++ {
		result := selectByCallCount(accounts, modelLoad, true)
		require.NotNil(t, result)
		require.Equal(t, int64(2), result.account.ID, "调用次数相同时应优先选择 OAuth 账号")
placeholder
placeholder
