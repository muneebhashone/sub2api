package service

import (
	"context"
	"sync"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type panelRateLimitSettingRepo struct {
	mu            sync.Mutex
	values        map[string]string
	getValueErr   error
	getValueCalls int
placeholder

func (r *panelRateLimitSettingRepo) Get(_ context.Context, key string) (*Setting, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	value, ok := r.values[key]
	if !ok {
		return nil, ErrSettingNotFound
placeholder
	return &Setting{Key: key, Value: valueplaceholder, nil
placeholder

func (r *panelRateLimitSettingRepo) GetValue(_ context.Context, key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.getValueCalls++
	if r.getValueErr != nil {
		return "", r.getValueErr
placeholder
	value, ok := r.values[key]
	if !ok {
		return "", ErrSettingNotFound
placeholder
	return value, nil
placeholder

func (r *panelRateLimitSettingRepo) Set(_ context.Context, key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.values == nil {
		r.values = make(map[string]string)
placeholder
	r.values[key] = value
	return nil
placeholder

func (r *panelRateLimitSettingRepo) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			out[key] = value
	placeholder
placeholder
	return out, nil
placeholder

func (r *panelRateLimitSettingRepo) SetMultiple(_ context.Context, settings map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.values == nil {
		r.values = make(map[string]string)
placeholder
	for key, value := range settings {
		r.values[key] = value
placeholder
	return nil
placeholder

func (r *panelRateLimitSettingRepo) GetAll(_ context.Context) (map[string]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make(map[string]string, len(r.values))
	for key, value := range r.values {
		out[key] = value
placeholder
	return out, nil
placeholder

func (r *panelRateLimitSettingRepo) Delete(_ context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.values, key)
	return nil
placeholder

func newPanelRateLimitTestService(repo SettingRepository) *SettingService {
	return NewSettingService(repo, &config.Config{placeholder)
placeholder

func TestGetPanelRateLimitSettingsDefaults(t *testing.T) {
	svc := newPanelRateLimitTestService(&panelRateLimitSettingRepo{placeholder)

	settings, err := svc.GetPanelRateLimitSettings(context.Background())
placeholder
	require.Equal(t, DefaultPanelRateLimitSettings(), settings)
placeholder

func TestGetPanelRateLimitSettingsInvalidJSONFallsBack(t *testing.T) {
	repo := &panelRateLimitSettingRepo{values: map[string]string{
		SettingKeyPanelRateLimitSettings: "{not-json",
placeholderplaceholder
	svc := newPanelRateLimitTestService(repo)

	settings, err := svc.GetPanelRateLimitSettings(context.Background())
placeholder
	require.Equal(t, DefaultPanelRateLimitSettings(), settings)
placeholder

func TestGetPanelRateLimitSettingsNormalizesValues(t *testing.T) {
	repo := &panelRateLimitSettingRepo{values: map[string]string{
		SettingKeyPanelRateLimitSettings: `{"enabled":true,"user_rpm":-5,"heavy_rpm":999999999,"exempt_admin":false,"public_ip_rpm":10placeholder`,
placeholderplaceholder
	svc := newPanelRateLimitTestService(repo)

	settings, err := svc.GetPanelRateLimitSettings(context.Background())
placeholder
	require.True(t, settings.Enabled)
	require.Equal(t, 0, settings.UserRPM)
	require.Equal(t, panelRateLimitRPMMax, settings.HeavyRPM)
	require.Equal(t, 10, settings.PublicIPRPM)
	require.False(t, settings.ExemptAdmin)
placeholder

func TestSetPanelRateLimitSettingsValidation(t *testing.T) {
	svc := newPanelRateLimitTestService(&panelRateLimitSettingRepo{placeholder)

	require.Error(t, svc.SetPanelRateLimitSettings(context.Background(), nil))
	require.Error(t, svc.SetPanelRateLimitSettings(context.Background(), &PanelRateLimitSettings{UserRPM: -1placeholder))
	require.Error(t, svc.SetPanelRateLimitSettings(context.Background(), &PanelRateLimitSettings{HeavyRPM: panelRateLimitRPMMax + 1placeholder))
placeholder

func TestSetPanelRateLimitSettingsRoundTripAndCacheRefresh(t *testing.T) {
	repo := &panelRateLimitSettingRepo{placeholder
	svc := newPanelRateLimitTestService(repo)

	// 先填充缓存（默认值）
	cached := svc.GetPanelRateLimitSettingsCached(context.Background())
	require.Equal(t, *DefaultPanelRateLimitSettings(), cached)

	want := &PanelRateLimitSettings{
		Enabled:     true,
		UserRPM:     120,
		HeavyRPM:    30,
		ExemptAdmin: false,
		PublicIPRPM: 60,
placeholder
	require.NoError(t, svc.SetPanelRateLimitSettings(context.Background(), want))

	// 写入后无需等待 TTL，缓存立即反映新值
	cached = svc.GetPanelRateLimitSettingsCached(context.Background())
	require.Equal(t, *want, cached)

	// DB 中持久化的值可直读
	stored, err := svc.GetPanelRateLimitSettings(context.Background())
placeholder
	require.Equal(t, want, stored)
placeholder

func TestGetPanelRateLimitSettingsCachedAvoidsRepeatedDBReads(t *testing.T) {
	repo := &panelRateLimitSettingRepo{values: map[string]string{
		SettingKeyPanelRateLimitSettings: `{"enabled":true,"user_rpm":100,"heavy_rpm":20,"exempt_admin":true,"public_ip_rpm":50placeholder`,
placeholderplaceholder
	svc := newPanelRateLimitTestService(repo)

	for i := 0; i < 5; i++ {
		settings := svc.GetPanelRateLimitSettingsCached(context.Background())
		require.Equal(t, 100, settings.UserRPM)
placeholder

	repo.mu.Lock()
	calls := repo.getValueCalls
	repo.mu.Unlock()
	require.Equal(t, 1, calls, "TTL 内应只读一次 DB")
placeholder
