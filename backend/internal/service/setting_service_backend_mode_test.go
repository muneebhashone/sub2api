//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type bmRepoStub struct {
	getValueFn func(ctx context.Context, key string) (string, error)
	calls      int
placeholder

func (s *bmRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *bmRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	s.calls++
	if s.getValueFn == nil {
		panic("unexpected GetValue call")
placeholder
	return s.getValueFn(ctx, key)
placeholder

func (s *bmRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *bmRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *bmRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *bmRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *bmRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

type bmUpdateRepoStub struct {
	updates    map[string]string
	getValueFn func(ctx context.Context, key string) (string, error)
placeholder

func (s *bmUpdateRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *bmUpdateRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if s.getValueFn == nil {
		panic("unexpected GetValue call")
placeholder
	return s.getValueFn(ctx, key)
placeholder

func (s *bmUpdateRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *bmUpdateRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *bmUpdateRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	s.updates = make(map[string]string, len(settings))
	for k, v := range settings {
		s.updates[k] = v
placeholder
	return nil
placeholder

func (s *bmUpdateRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *bmUpdateRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

func resetBackendModeTestCache(t *testing.T) {
placeholder

	backendModeCache.Store((*cachedBackendMode)(nil))
	t.Cleanup(func() {
		backendModeCache.Store((*cachedBackendMode)(nil))
placeholder)
placeholder

func TestIsBackendModeEnabled_ReturnsTrue(t *testing.T) {
	resetBackendModeTestCache(t)

	repo := &bmRepoStub{
		getValueFn: func(ctx context.Context, key string) (string, error) {
			require.Equal(t, SettingKeyBackendModeEnabled, key)
			return "true", nil
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	require.True(t, svc.IsBackendModeEnabled(context.Background()))
	require.Equal(t, 1, repo.calls)
placeholder

func TestIsBackendModeEnabled_ReturnsFalse(t *testing.T) {
	resetBackendModeTestCache(t)

	repo := &bmRepoStub{
		getValueFn: func(ctx context.Context, key string) (string, error) {
			require.Equal(t, SettingKeyBackendModeEnabled, key)
			return "false", nil
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	require.False(t, svc.IsBackendModeEnabled(context.Background()))
	require.Equal(t, 1, repo.calls)
placeholder

func TestIsBackendModeEnabled_ReturnsFalseOnNotFound(t *testing.T) {
	resetBackendModeTestCache(t)

	repo := &bmRepoStub{
		getValueFn: func(ctx context.Context, key string) (string, error) {
			require.Equal(t, SettingKeyBackendModeEnabled, key)
			return "", ErrSettingNotFound
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	require.False(t, svc.IsBackendModeEnabled(context.Background()))
	require.Equal(t, 1, repo.calls)
placeholder

func TestIsBackendModeEnabled_ReturnsFalseOnDBError(t *testing.T) {
	resetBackendModeTestCache(t)

	repo := &bmRepoStub{
		getValueFn: func(ctx context.Context, key string) (string, error) {
			require.Equal(t, SettingKeyBackendModeEnabled, key)
			return "", errors.New("db down")
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	require.False(t, svc.IsBackendModeEnabled(context.Background()))
	require.Equal(t, 1, repo.calls)
placeholder

func TestIsBackendModeEnabled_CachesResult(t *testing.T) {
	resetBackendModeTestCache(t)

	repo := &bmRepoStub{
		getValueFn: func(ctx context.Context, key string) (string, error) {
			require.Equal(t, SettingKeyBackendModeEnabled, key)
			return "true", nil
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	require.True(t, svc.IsBackendModeEnabled(context.Background()))
	require.True(t, svc.IsBackendModeEnabled(context.Background()))
	require.Equal(t, 1, repo.calls)
placeholder

func TestUpdateSettings_InvalidatesBackendModeCache(t *testing.T) {
	resetBackendModeTestCache(t)

	backendModeCache.Store(&cachedBackendMode{
		value:     true,
		expiresAt: time.Now().Add(backendModeCacheTTL).UnixNano(),
placeholder)

	repo := &bmUpdateRepoStub{
		getValueFn: func(ctx context.Context, key string) (string, error) {
			require.Equal(t, SettingKeyBackendModeEnabled, key)
			return "true", nil
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		BackendModeEnabled: false,
placeholder)
placeholder
	require.Equal(t, "false", repo.updates[SettingKeyBackendModeEnabled])
	require.False(t, svc.IsBackendModeEnabled(context.Background()))
placeholder
