//go:build unit

package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/antigravity"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

type settingUpdateRepoStub struct {
	updates        map[string]string
	setMultipleErr error
placeholder

func (s *settingUpdateRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *settingUpdateRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
placeholder

func (s *settingUpdateRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *settingUpdateRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *settingUpdateRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	s.updates = make(map[string]string, len(settings))
	for k, v := range settings {
		s.updates[k] = v
placeholder
	return s.setMultipleErr
placeholder

func (s *settingUpdateRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *settingUpdateRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

type settingGetAllRepoStub struct {
	values map[string]string
placeholder

func (s *settingGetAllRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *settingGetAllRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	panic("unexpected GetValue call")
placeholder

func (s *settingGetAllRepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *settingGetAllRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *settingGetAllRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *settingGetAllRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	out := make(map[string]string, len(s.values))
	for key, value := range s.values {
		out[key] = value
placeholder
	return out, nil
placeholder

func (s *settingGetAllRepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

type forwardedIPMigrationRepoStub struct {
	values         map[string]string
	updates        map[string]string
	getMultipleErr error
	setMultipleErr error
placeholder

func (s *forwardedIPMigrationRepoStub) Get(context.Context, string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *forwardedIPMigrationRepoStub) GetValue(_ context.Context, key string) (string, error) {
	value, ok := s.values[key]
	if !ok {
		return "", ErrSettingNotFound
placeholder
	return value, nil
placeholder

func (s *forwardedIPMigrationRepoStub) Set(context.Context, string, string) error {
	panic("unexpected Set call")
placeholder

func (s *forwardedIPMigrationRepoStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	if s.getMultipleErr != nil {
		return nil, s.getMultipleErr
placeholder
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := s.values[key]; ok {
			result[key] = value
	placeholder
placeholder
	return result, nil
placeholder

func (s *forwardedIPMigrationRepoStub) SetMultiple(_ context.Context, values map[string]string) error {
	if s.setMultipleErr != nil {
		return s.setMultipleErr
placeholder
	s.updates = make(map[string]string, len(values))
	for key, value := range values {
		s.values[key] = value
		s.updates[key] = value
placeholder
	return nil
placeholder

func (s *forwardedIPMigrationRepoStub) GetAll(context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *forwardedIPMigrationRepoStub) Delete(context.Context, string) error {
	panic("unexpected Delete call")
placeholder

type settingAntigravityUARepoStub struct {
	values map[string]string
placeholder

func (s *settingAntigravityUARepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	panic("unexpected Get call")
placeholder

func (s *settingAntigravityUARepoStub) GetValue(ctx context.Context, key string) (string, error) {
	if value, ok := s.values[key]; ok {
		return value, nil
placeholder
	return "", ErrSettingNotFound
placeholder

func (s *settingAntigravityUARepoStub) Set(ctx context.Context, key, value string) error {
	panic("unexpected Set call")
placeholder

func (s *settingAntigravityUARepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	panic("unexpected GetMultiple call")
placeholder

func (s *settingAntigravityUARepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	panic("unexpected SetMultiple call")
placeholder

func (s *settingAntigravityUARepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	panic("unexpected GetAll call")
placeholder

func (s *settingAntigravityUARepoStub) Delete(ctx context.Context, key string) error {
	panic("unexpected Delete call")
placeholder

type defaultSubGroupReaderStub struct {
	byID  map[int64]*Group
	errBy map[int64]error
	calls []int64
placeholder

func TestSettingService_AffiliateAdminRechargeSetting(t *testing.T) {
	t.Run("missing value defaults to disabled", func(t *testing.T) {
		svc := NewSettingService(&settingGetAllRepoStub{values: map[string]string{placeholderplaceholder, &config.Config{placeholder)

		settings, err := svc.GetAllSettings(context.Background())
	placeholder
		require.False(t, settings.AdminRechargeRebateEnabled)
placeholder)

	t.Run("explicit value is parsed", func(t *testing.T) {
		svc := NewSettingService(&settingGetAllRepoStub{values: map[string]string{
			SettingKeyAffiliateAdminRechargeEnabled: "true",
placeholder &config.Config{placeholder)

		settings, err := svc.GetAllSettings(context.Background())
	placeholder
		require.True(t, settings.AdminRechargeRebateEnabled)
placeholder)

	t.Run("value is persisted", func(t *testing.T) {
		repo := &settingUpdateRepoStub{placeholder
		svc := NewSettingService(repo, &config.Config{placeholder)

		err := svc.UpdateSettings(context.Background(), &SystemSettings{
			AdminRechargeRebateEnabled: true,
	placeholder)
	placeholder
		require.Equal(t, "true", repo.updates[SettingKeyAffiliateAdminRechargeEnabled])
placeholder)
placeholder

func (s *defaultSubGroupReaderStub) GetByID(ctx context.Context, id int64) (*Group, error) {
	s.calls = append(s.calls, id)
	if err, ok := s.errBy[id]; ok {
		return nil, err
placeholder
	if g, ok := s.byID[id]; ok {
		return g, nil
placeholder
	return nil, ErrGroupNotFound
placeholder

func TestSettingService_UpdateSettings_DefaultSubscriptions_ValidGroup(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	groupReader := &defaultSubGroupReaderStub{
		byID: map[int64]*Group{
			11: {ID: 11, SubscriptionType: SubscriptionTypeSubscriptionplaceholder,
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)
	svc.SetDefaultSubscriptionGroupReader(groupReader)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		DefaultSubscriptions: []DefaultSubscriptionSetting{
			{GroupID: 11, ValidityDays: 30placeholder,
	placeholder,
placeholder)
placeholder
	require.Equal(t, []int64{11placeholder, groupReader.calls)

	raw, ok := repo.updates[SettingKeyDefaultSubscriptions]
	require.True(t, ok)

	var got []DefaultSubscriptionSetting
	require.NoError(t, json.Unmarshal([]byte(raw), &got))
	require.Equal(t, []DefaultSubscriptionSetting{
		{GroupID: 11, ValidityDays: 30placeholder,
placeholder, got)
placeholder

func TestSettingService_UpdateSettings_DefaultSubscriptions_RejectsNonSubscriptionGroup(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	groupReader := &defaultSubGroupReaderStub{
		byID: map[int64]*Group{
			12: {ID: 12, SubscriptionType: SubscriptionTypeStandardplaceholder,
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)
	svc.SetDefaultSubscriptionGroupReader(groupReader)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		DefaultSubscriptions: []DefaultSubscriptionSetting{
			{GroupID: 12, ValidityDays: 7placeholder,
	placeholder,
placeholder)
placeholder
	require.Equal(t, "DEFAULT_SUBSCRIPTION_GROUP_INVALID", infraerrors.Reason(err))
	require.Nil(t, repo.updates)
placeholder

func TestSettingService_UpdateSettings_DefaultSubscriptions_RejectsNotFoundGroup(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	groupReader := &defaultSubGroupReaderStub{
		errBy: map[int64]error{
			13: ErrGroupNotFound,
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)
	svc.SetDefaultSubscriptionGroupReader(groupReader)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		DefaultSubscriptions: []DefaultSubscriptionSetting{
			{GroupID: 13, ValidityDays: 7placeholder,
	placeholder,
placeholder)
placeholder
	require.Equal(t, "DEFAULT_SUBSCRIPTION_GROUP_INVALID", infraerrors.Reason(err))
	require.Equal(t, "13", infraerrors.FromError(err).Metadata["group_id"])
	require.Nil(t, repo.updates)
placeholder

func TestSettingService_UpdateSettings_DefaultSubscriptions_RejectsDuplicateGroup(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	groupReader := &defaultSubGroupReaderStub{
		byID: map[int64]*Group{
			11: {ID: 11, SubscriptionType: SubscriptionTypeSubscriptionplaceholder,
	placeholder,
placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)
	svc.SetDefaultSubscriptionGroupReader(groupReader)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		DefaultSubscriptions: []DefaultSubscriptionSetting{
			{GroupID: 11, ValidityDays: 30placeholder,
			{GroupID: 11, ValidityDays: 60placeholder,
	placeholder,
placeholder)
placeholder
	require.Equal(t, "DEFAULT_SUBSCRIPTION_GROUP_DUPLICATE", infraerrors.Reason(err))
	require.Equal(t, "11", infraerrors.FromError(err).Metadata["group_id"])
	require.Nil(t, repo.updates)
placeholder

func TestSettingService_UpdateSettings_DefaultSubscriptions_RejectsDuplicateGroupWithoutGroupReader(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		DefaultSubscriptions: []DefaultSubscriptionSetting{
			{GroupID: 11, ValidityDays: 30placeholder,
			{GroupID: 11, ValidityDays: 60placeholder,
	placeholder,
placeholder)
placeholder
	require.Equal(t, "DEFAULT_SUBSCRIPTION_GROUP_DUPLICATE", infraerrors.Reason(err))
	require.Equal(t, "11", infraerrors.FromError(err).Metadata["group_id"])
	require.Nil(t, repo.updates)
placeholder

func TestSettingService_UpdateSettings_RegistrationEmailSuffixWhitelist_Normalized(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		RegistrationEmailSuffixWhitelist: []string{"example.com", "@EXAMPLE.com", " @foo.bar ", "*.EDU.CN"placeholder,
placeholder)
placeholder
	require.Equal(t, `["@example.com","@foo.bar","*.edu.cn"]`, repo.updates[SettingKeyRegistrationEmailSuffixWhitelist])
placeholder

func TestSettingService_UpdateSettings_RegistrationEmailSuffixWhitelist_Invalid(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		RegistrationEmailSuffixWhitelist: []string{"@invalid_domain"placeholder,
placeholder)
placeholder
	require.Equal(t, "INVALID_REGISTRATION_EMAIL_SUFFIX_WHITELIST", infraerrors.Reason(err))
placeholder

func TestParseDefaultSubscriptions_NormalizesValues(t *testing.T) {
	got := parseDefaultSubscriptions(`[{"group_id":11,"validity_days":30placeholder,{"group_id":11,"validity_days":60placeholder,{"group_id":0,"validity_days":10placeholder,{"group_id":12,"validity_days":99999placeholder]`)
	require.Equal(t, []DefaultSubscriptionSetting{
		{GroupID: 11, ValidityDays: 30placeholder,
		{GroupID: 11, ValidityDays: 60placeholder,
		{GroupID: 12, ValidityDays: MaxValidityDaysplaceholder,
placeholder, got)
placeholder

func TestSettingService_UpdateSettings_TablePreferences(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		TableDefaultPageSize: 50,
		TablePageSizeOptions: []int{20, 50, 100placeholder,
placeholder)
placeholder
	require.Equal(t, "50", repo.updates[SettingKeyTableDefaultPageSize])
	require.Equal(t, "[20,50,100]", repo.updates[SettingKeyTablePageSizeOptions])

	err = svc.UpdateSettings(context.Background(), &SystemSettings{
		TableDefaultPageSize: 1000,
		TablePageSizeOptions: []int{20, 100placeholder,
placeholder)
placeholder
	require.Equal(t, "1000", repo.updates[SettingKeyTableDefaultPageSize])
	require.Equal(t, "[20,100]", repo.updates[SettingKeyTablePageSizeOptions])
placeholder

func TestSettingService_UpdateSettings_PaymentVisibleMethodsAndAdvancedScheduler(t *testing.T) {
	resetOpenAIAdvancedSchedulerSettingCacheForTest()
	defer resetOpenAIAdvancedSchedulerSettingCacheForTest()

	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		PaymentVisibleMethodAlipaySource:                   "alipay",
		PaymentVisibleMethodWxpaySource:                    "easypay",
		PaymentVisibleMethodAlipayEnabled:                  true,
		PaymentVisibleMethodWxpayEnabled:                   false,
		OpenAILowUpstreamRatePriorityEnabled:               true,
		OpenAIOAuthSchedulingRateMultiplier:                0.05,
		OpenAIAdvancedSchedulerEnabled:                     true,
		OpenAIAdvancedSchedulerStickyWeightedEnabled:       true,
		OpenAIAdvancedSchedulerSubscriptionPriorityEnabled: true,
		OpenAIAdvancedSchedulerLBTopK:                      " 3 ",
		OpenAIAdvancedSchedulerWeightPriority:              "2.50",
		OpenAIAdvancedSchedulerWeightLoad:                  "0",
		OpenAIAdvancedSchedulerWeightQueue:                 "0.75",
		OpenAIAdvancedSchedulerWeightErrorRate:             "1.25",
		OpenAIAdvancedSchedulerWeightTTFT:                  "0.5",
		OpenAIAdvancedSchedulerWeightReset:                 "",
		OpenAIAdvancedSchedulerWeightQuotaHeadroom:         "0.2",
		OpenAIAdvancedSchedulerWeightUpstreamCost:          "1.5",
		OpenAIAdvancedSchedulerWeightPreviousResponse:      "8",
		OpenAIAdvancedSchedulerWeightSessionSticky:         "4",
placeholder)
placeholder
	require.Equal(t, VisibleMethodSourceOfficialAlipay, repo.updates[SettingPaymentVisibleMethodAlipaySource])
	require.Equal(t, VisibleMethodSourceEasyPayWechat, repo.updates[SettingPaymentVisibleMethodWxpaySource])
	require.Equal(t, "true", repo.updates[SettingPaymentVisibleMethodAlipayEnabled])
	require.Equal(t, "false", repo.updates[SettingPaymentVisibleMethodWxpayEnabled])
	require.Equal(t, "true", repo.updates[SettingKeyOpenAILowUpstreamRatePriorityEnabled])
	require.Equal(t, "0.05", repo.updates[SettingKeyOpenAIOAuthSchedulingRateMultiplier])
	require.Equal(t, "true", repo.updates[openAIAdvancedSchedulerSettingKey])
	require.Equal(t, "true", repo.updates[SettingKeyOpenAIAdvancedSchedulerStickyWeightedEnabled])
	require.Equal(t, "true", repo.updates[SettingKeyOpenAIAdvancedSchedulerSubscriptionPriorityEnabled])
	require.Equal(t, "3", repo.updates[SettingKeyOpenAIAdvancedSchedulerLBTopK])
	require.Equal(t, "2.5", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightPriority])
	require.Equal(t, "0", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightLoad])
	require.Equal(t, "0.75", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightQueue])
	require.Equal(t, "1.25", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightErrorRate])
	require.Equal(t, "0.5", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightTTFT])
	require.Equal(t, "", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightReset])
	require.Equal(t, "0.2", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightQuotaHeadroom])
	require.Equal(t, "1.5", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightUpstreamCost])
	require.Equal(t, "8", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightPreviousResponse])
	require.Equal(t, "4", repo.updates[SettingKeyOpenAIAdvancedSchedulerWeightSessionSticky])
placeholder

func TestSettingService_UpdateSettingsRejectsInvalidOpenAIOAuthSchedulingRateMultiplier(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	for _, rate := range []float64{-0.01, math.NaN(), math.Inf(1)placeholder {
		err := svc.UpdateSettings(context.Background(), &SystemSettings{OpenAIOAuthSchedulingRateMultiplier: rateplaceholder)
	placeholder
placeholder
placeholder

func TestSettingService_UpdateSettings_OpenAIAdvancedSchedulerWeightSums(t *testing.T) {
	maxFloat := strconv.FormatFloat(math.MaxFloat64, 'g', -1, 64)
	tests := []struct {
		name    string
		weights SystemSettings
		wantErr bool
placeholder{
		{
			name: "reset only base is valid",
			weights: SystemSettings{
				OpenAIAdvancedSchedulerWeightPriority:         "0",
				OpenAIAdvancedSchedulerWeightLoad:             "0",
				OpenAIAdvancedSchedulerWeightQueue:            "0",
				OpenAIAdvancedSchedulerWeightErrorRate:        "0",
				OpenAIAdvancedSchedulerWeightTTFT:             "0",
				OpenAIAdvancedSchedulerWeightReset:            "1",
				OpenAIAdvancedSchedulerWeightQuotaHeadroom:    "0",
				OpenAIAdvancedSchedulerWeightUpstreamCost:     "0",
				OpenAIAdvancedSchedulerWeightPreviousResponse: "0",
				OpenAIAdvancedSchedulerWeightSessionSticky:    "0",
		placeholder,
	placeholder,
		{
			name: "base sum overflow is rejected",
			weights: SystemSettings{
				OpenAIAdvancedSchedulerWeightPriority: maxFloat,
				OpenAIAdvancedSchedulerWeightLoad:     maxFloat,
		placeholder,
			wantErr: true,
	placeholder,
		{
			name: "sticky total sum overflow is rejected",
			weights: SystemSettings{
				OpenAIAdvancedSchedulerWeightPriority:         maxFloat,
				OpenAIAdvancedSchedulerWeightPreviousResponse: maxFloat,
		placeholder,
			wantErr: true,
	placeholder,
placeholder

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSettingService(&settingUpdateRepoStub{placeholder, &config.Config{placeholder)
			err := svc.UpdateSettings(context.Background(), &tt.weights)
			if tt.wantErr {
			placeholder
				return
		placeholder
		placeholder
	placeholder)
placeholder
placeholder

func TestSettingService_ParseSettingsDefaultsOpenAIOAuthSchedulingRateMultiplier(t *testing.T) {
	svc := NewSettingService(&settingUpdateRepoStub{placeholder, &config.Config{placeholder)

	require.Equal(t, 1.0, svc.parseSettings(map[string]string{placeholder).OpenAIOAuthSchedulingRateMultiplier)
	require.Equal(t, 0.05, svc.parseSettings(map[string]string{SettingKeyOpenAIOAuthSchedulingRateMultiplier: "0.05"placeholder).OpenAIOAuthSchedulingRateMultiplier)
placeholder

func TestSettingService_GetAllSettings_OpenAIAdvancedSchedulerEffectiveValuesUseConfig(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.Gateway.OpenAIWS.LBTopK = 13
	cfg.Gateway.OpenAIWS.SchedulerScoreWeights = config.GatewayOpenAIWSSchedulerScoreWeights{
		Priority:         2,
		Load:             3,
		Queue:            4,
		ErrorRate:        5,
		TTFT:             6,
		Reset:            7,
		QuotaHeadroom:    8,
		UpstreamCost:     9,
		PreviousResponse: 10,
		SessionSticky:    11,
placeholder
	svc := NewSettingService(&settingGetAllRepoStub{values: map[string]string{
		SettingKeyOpenAIAdvancedSchedulerLBTopK:              "3",
		SettingKeyOpenAIAdvancedSchedulerWeightPriority:      "99",
		SettingKeyOpenAIAdvancedSchedulerWeightSessionSticky: "88",
placeholderplaceholder, cfg)

	settings, err := svc.GetAllSettings(context.Background())
placeholder
	require.Equal(t, "3", settings.OpenAIAdvancedSchedulerLBTopK)
	require.Equal(t, "99", settings.OpenAIAdvancedSchedulerWeightPriority)
	require.Equal(t, "88", settings.OpenAIAdvancedSchedulerWeightSessionSticky)
	require.Equal(t, "13", settings.OpenAIAdvancedSchedulerEffectiveLBTopK)
	require.Equal(t, "2", settings.OpenAIAdvancedSchedulerEffectiveWeightPriority)
	require.Equal(t, "3", settings.OpenAIAdvancedSchedulerEffectiveWeightLoad)
	require.Equal(t, "9", settings.OpenAIAdvancedSchedulerEffectiveWeightUpstreamCost)
	require.Equal(t, "11", settings.OpenAIAdvancedSchedulerEffectiveWeightSessionSticky)
placeholder

func TestSettingService_UpdateSettings_AntigravityUserAgentVersion(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		AntigravityUserAgentVersion: "1.23.2",
placeholder)
placeholder
	require.Equal(t, "1.23.2", repo.updates[SettingKeyAntigravityUserAgentVersion])
placeholder

func TestSettingService_InitializeDefaultSettingsPersistsConfiguredForwardedClientIPHeaders(t *testing.T) {
	repo := &forwardedIPMigrationRepoStub{values: map[string]string{placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.SetForwardedClientIPSettings(true, []string{"X-Cdn-Ip", "True-Client-Ip"placeholder)
	svc := NewSettingService(repo, cfg)

	require.NoError(t, svc.InitializeDefaultSettings(context.Background()))
	require.JSONEq(t, `["X-Cdn-Ip","True-Client-Ip"]`, repo.values[SettingKeyForwardedClientIPHeaders])
placeholder

func TestSettingService_UpdateSettings_APIKeyACLTrustForwardedIPRefreshesConfig(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	cfg := &config.Config{placeholder
	svc := NewSettingService(repo, cfg)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		APIKeyACLTrustForwardedIP: true,
		ForwardedClientIPHeaders:  []string{" x-cdn-ip ", "X-CDN-IP", "true-client-ip"placeholder,
placeholder)
placeholder
	require.Equal(t, "true", repo.updates[SettingKeyAPIKeyACLTrustForwardedIP])
	require.JSONEq(t, `["X-Cdn-Ip","True-Client-Ip"]`, repo.updates[SettingKeyForwardedClientIPHeaders])
	runtimeSettings := cfg.ForwardedClientIPSettings()
	require.True(t, runtimeSettings.TrustForwardedIP)
	require.Equal(t, []string{"X-Cdn-Ip", "True-Client-Ip"placeholder, runtimeSettings.Headers)

	runtimeSettings.Headers[0] = "X-Mutated"
	require.Equal(t, []string{"X-Cdn-Ip", "True-Client-Ip"placeholder, cfg.ForwardedClientIPSettings().Headers)
placeholder

func TestSettingService_UpdateSettings_RejectsInvalidForwardedClientIPHeadersWithoutRefreshing(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	cfg := &config.Config{placeholder
	cfg.SetForwardedClientIPSettings(true, []string{"X-Existing-IP"placeholder)
	svc := NewSettingService(repo, cfg)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		ForwardedClientIPHeaders: []string{"X Invalid"placeholder,
placeholder)

placeholder
	require.Nil(t, repo.updates)
	runtimeSettings := cfg.ForwardedClientIPSettings()
	require.True(t, runtimeSettings.TrustForwardedIP)
	require.Equal(t, []string{"X-Existing-IP"placeholder, runtimeSettings.Headers)
placeholder

func TestSettingService_UpdateSettings_WriteFailureDoesNotRefreshForwardedIPRuntime(t *testing.T) {
	repo := &settingUpdateRepoStub{setMultipleErr: errors.New("database unavailable")placeholder
	cfg := &config.Config{placeholder
	cfg.SetForwardedClientIPSettings(false, []string{"X-Existing-IP"placeholder)
	svc := NewSettingService(repo, cfg)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		APIKeyACLTrustForwardedIP: true,
		ForwardedClientIPHeaders:  []string{"X-New-IP"placeholder,
placeholder)

	require.ErrorContains(t, err, "database unavailable")
	runtimeSettings := cfg.ForwardedClientIPSettings()
	require.False(t, runtimeSettings.TrustForwardedIP)
	require.Equal(t, []string{"X-Existing-IP"placeholder, runtimeSettings.Headers)
placeholder

func TestSettingService_ParseSettings_APIKeyACLTrustForwardedIPFallsBackToConfigWhenMissing(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.Security.TrustForwardedIPForAPIKeyACL = true
	svc := NewSettingService(&settingUpdateRepoStub{placeholder, cfg)

	got := svc.parseSettings(map[string]string{placeholder)

	require.True(t, got.APIKeyACLTrustForwardedIP)
placeholder

func TestSettingService_ParseSettings_APIKeyACLTrustForwardedIPUsesStoredValue(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.SetTrustForwardedIPForAPIKeyACL(true)
	svc := NewSettingService(&settingUpdateRepoStub{placeholder, cfg)

	got := svc.parseSettings(map[string]string{SettingKeyAPIKeyACLTrustForwardedIP: "false"placeholder)

	require.False(t, got.APIKeyACLTrustForwardedIP)
placeholder

func TestSettingService_ParseSettings_ForwardedClientIPHeaders(t *testing.T) {
	cfg := &config.Config{placeholder
	cfg.SetForwardedClientIPSettings(true, []string{"X-Config-IP"placeholder)
	svc := NewSettingService(&settingUpdateRepoStub{placeholder, cfg)

	t.Run("stored value is normalized", func(t *testing.T) {
		got := svc.parseSettings(map[string]string{
			SettingKeyForwardedClientIPHeaders: `[" x-cdn-ip ","X-CDN-IP","true-client-ip"]`,
	placeholder)
		require.Equal(t, []string{"X-Cdn-Ip", "True-Client-Ip"placeholder, got.ForwardedClientIPHeaders)
placeholder)

	t.Run("missing value falls back to config", func(t *testing.T) {
		got := svc.parseSettings(map[string]string{placeholder)
		require.Equal(t, []string{"X-Config-IP"placeholder, got.ForwardedClientIPHeaders)
placeholder)

	t.Run("malformed value disables forwarded trust", func(t *testing.T) {
		got := svc.parseSettings(map[string]string{
			SettingKeyAPIKeyACLTrustForwardedIP: "true",
			SettingKeyForwardedClientIPHeaders:  `{"not":"an array"placeholder`,
	placeholder)
		require.False(t, got.APIKeyACLTrustForwardedIP)
		require.Empty(t, got.ForwardedClientIPHeaders)
placeholder)
placeholder

func TestSettingService_LoadForwardedClientIPSettingsMigration(t *testing.T) {
	tests := []struct {
		name                   string
		values                 map[string]string
		trustedProxiesSet      bool
		configDefault          bool
		wantEnabled            bool
		wantForwardedIPUpdate  string
		wantMigrationMarkerSet bool
placeholder{
		{
			name:                   "missing setting follows configured default",
			values:                 map[string]string{placeholder,
			configDefault:          true,
			wantEnabled:            true,
			wantMigrationMarkerSet: true,
	placeholder,
		{
			name:                   "legacy false without proxy config migrates to compatibility",
			values:                 map[string]string{SettingKeyAPIKeyACLTrustForwardedIP: "false"placeholder,
			wantEnabled:            true,
			wantForwardedIPUpdate:  "true",
			wantMigrationMarkerSet: true,
	placeholder,
		{
			name:                   "legacy false with explicit proxy config stays secure",
			values:                 map[string]string{SettingKeyAPIKeyACLTrustForwardedIP: "false"placeholder,
			trustedProxiesSet:      true,
			wantEnabled:            false,
			wantMigrationMarkerSet: true,
	placeholder,
		{
			name: "completed migration preserves later false choice",
			values: map[string]string{
				SettingKeyAPIKeyACLTrustForwardedIP: "false",
				settingKeyForwardedClientIPModeV2:   "true",
		placeholder,
			wantEnabled: false,
	placeholder,
placeholder

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &forwardedIPMigrationRepoStub{values: test.valuesplaceholder
			cfg := &config.Config{Server: config.ServerConfig{TrustedProxiesConfigured: test.trustedProxiesSetplaceholderplaceholder
			cfg.Security.TrustForwardedIPForAPIKeyACL = test.configDefault
			svc := NewSettingService(repo, cfg)

			require.NoError(t, svc.LoadForwardedClientIPSettings(context.Background()))
			require.Equal(t, test.wantEnabled, cfg.TrustForwardedIPForAPIKeyACL())
			require.Equal(t, test.wantForwardedIPUpdate, repo.updates[SettingKeyAPIKeyACLTrustForwardedIP])
			require.JSONEq(t, `[]`, repo.updates[SettingKeyForwardedClientIPHeaders])
			if test.wantMigrationMarkerSet {
				require.Equal(t, "true", repo.updates[settingKeyForwardedClientIPModeV2])
		placeholder else {
				require.NotContains(t, repo.updates, settingKeyForwardedClientIPModeV2)
		placeholder
	placeholder)
placeholder
placeholder

func TestSettingService_LoadForwardedClientIPSettingsLoadsHeaders(t *testing.T) {
	repo := &forwardedIPMigrationRepoStub{values: map[string]string{
		SettingKeyAPIKeyACLTrustForwardedIP: "true",
		SettingKeyForwardedClientIPHeaders:  `[" x-cdn-ip ","true-client-ip"]`,
		settingKeyForwardedClientIPModeV2:   "true",
placeholderplaceholder
	cfg := &config.Config{placeholder
	svc := NewSettingService(repo, cfg)

	require.NoError(t, svc.LoadForwardedClientIPSettings(context.Background()))
	runtimeSettings := cfg.ForwardedClientIPSettings()
	require.True(t, runtimeSettings.TrustForwardedIP)
	require.Equal(t, []string{"X-Cdn-Ip", "True-Client-Ip"placeholder, runtimeSettings.Headers)
	require.Nil(t, repo.updates)
placeholder

func TestSettingService_LoadForwardedClientIPSettingsMalformedHeadersDisablesCustomTrust(t *testing.T) {
	repo := &forwardedIPMigrationRepoStub{values: map[string]string{
		SettingKeyAPIKeyACLTrustForwardedIP: "true",
		SettingKeyForwardedClientIPHeaders:  `["X Invalid"]`,
placeholderplaceholder
	cfg := &config.Config{placeholder
	svc := NewSettingService(repo, cfg)

	err := svc.LoadForwardedClientIPSettings(context.Background())

	require.ErrorContains(t, err, "load forwarded client ip headers")
	runtimeSettings := cfg.ForwardedClientIPSettings()
	require.False(t, runtimeSettings.TrustForwardedIP)
	require.Empty(t, runtimeSettings.Headers)
	require.Equal(t, "true", repo.updates[settingKeyForwardedClientIPModeV2])
	require.NotContains(t, repo.updates, SettingKeyAPIKeyACLTrustForwardedIP)
placeholder

func TestSettingService_LoadForwardedClientIPSettingsBackfillsConfigHeaders(t *testing.T) {
	repo := &forwardedIPMigrationRepoStub{values: map[string]string{
		settingKeyForwardedClientIPModeV2: "true",
placeholderplaceholder
	cfg := &config.Config{placeholder
	cfg.SetForwardedClientIPSettings(false, []string{"X-Config-IP"placeholder)
	svc := NewSettingService(repo, cfg)

	require.NoError(t, svc.LoadForwardedClientIPSettings(context.Background()))
	require.JSONEq(t, `["X-Config-IP"]`, repo.updates[SettingKeyForwardedClientIPHeaders])
	require.Equal(t, []string{"X-Config-IP"placeholder, cfg.ForwardedClientIPSettings().Headers)
placeholder

func TestSettingService_LoadForwardedClientIPSettingsReadFailureFailsClosed(t *testing.T) {
	repo := &forwardedIPMigrationRepoStub{
		getMultipleErr: errors.New("database unavailable"),
placeholder
	cfg := &config.Config{placeholder
	cfg.SetTrustForwardedIPForAPIKeyACL(true)
	svc := NewSettingService(repo, cfg)

	err := svc.LoadForwardedClientIPSettings(context.Background())

	require.ErrorContains(t, err, "get forwarded client ip settings")
	runtimeSettings := cfg.ForwardedClientIPSettings()
	require.False(t, runtimeSettings.TrustForwardedIP)
	require.Empty(t, runtimeSettings.Headers)
placeholder

func TestSettingService_LoadForwardedClientIPSettingsWriteFailureUsesComputedMode(t *testing.T) {
	tests := []struct {
		name              string
		trustedProxiesSet bool
		wantEnabled       bool
placeholder{
		{name: "compatibility migration remains effective", wantEnabled: trueplaceholder,
		{name: "explicit proxy policy remains secure", trustedProxiesSet: true, wantEnabled: falseplaceholder,
placeholder

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := &forwardedIPMigrationRepoStub{
				values:         map[string]string{SettingKeyAPIKeyACLTrustForwardedIP: "false"placeholder,
				setMultipleErr: errors.New("database unavailable"),
		placeholder
			cfg := &config.Config{Server: config.ServerConfig{TrustedProxiesConfigured: test.trustedProxiesSetplaceholderplaceholder
			svc := NewSettingService(repo, cfg)

			err := svc.LoadForwardedClientIPSettings(context.Background())

			require.ErrorContains(t, err, "migrate forwarded client ip setting")
			require.Equal(t, test.wantEnabled, cfg.TrustForwardedIPForAPIKeyACL())
	placeholder)
placeholder
placeholder

func TestSettingService_GetAntigravityUserAgentVersion_Precedence(t *testing.T) {
	t.Run("后台设置优先", func(t *testing.T) {
		svc := NewSettingService(&settingAntigravityUARepoStub{values: map[string]string{
			SettingKeyAntigravityUserAgentVersion: "1.24.0",
placeholder &config.Config{placeholder)

		require.Equal(t, "1.24.0", svc.GetAntigravityUserAgentVersion(context.Background()))
placeholder)

	t.Run("空值回退配置默认值", func(t *testing.T) {
		svc := NewSettingService(&settingAntigravityUARepoStub{values: map[string]string{
			SettingKeyAntigravityUserAgentVersion: "",
placeholder &config.Config{placeholder)

		require.Equal(t, antigravity.GetDefaultUserAgentVersion(), svc.GetAntigravityUserAgentVersion(context.Background()))
placeholder)

	t.Run("缺失回退配置默认值", func(t *testing.T) {
		svc := NewSettingService(&settingAntigravityUARepoStub{values: map[string]string{placeholderplaceholder, &config.Config{placeholder)

		require.Equal(t, antigravity.GetDefaultUserAgentVersion(), svc.GetAntigravityUserAgentVersion(context.Background()))
placeholder)
placeholder

func TestSettingService_UpdateSettings_RejectsInvalidPaymentVisibleMethodSource(t *testing.T) {
	repo := &settingUpdateRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	err := svc.UpdateSettings(context.Background(), &SystemSettings{
		PaymentVisibleMethodAlipaySource: "not-a-provider",
placeholder)
placeholder
	require.Equal(t, "INVALID_PAYMENT_VISIBLE_METHOD_SOURCE", infraerrors.Reason(err))
	require.Nil(t, repo.updates)
placeholder

func TestSettingService_PasskeySwitchPersistsAndDefaultsToConfigured(t *testing.T) {
	cfg := &config.Config{WebAuthn: config.WebAuthnConfig{
		Enabled:   true,
		RPID:      "sub3.nebula-spaces.com",
		RPOrigins: []string{"https://sub3.nebula-spaces.com"placeholder,
placeholderplaceholder
	runtimeRepo := &forwardedIPMigrationRepoStub{values: map[string]string{placeholderplaceholder
	runtimeService := NewSettingService(runtimeRepo, cfg)

	enabled, err := runtimeService.PasskeyEnabled(context.Background())
placeholder
	require.True(t, enabled)

	updateRepo := &settingUpdateRepoStub{placeholder
	updateService := NewSettingService(updateRepo, cfg)
	require.NoError(t, updateService.UpdateSettings(context.Background(), &SystemSettings{
		PasskeyEnabled: false,
placeholder))
	require.Equal(t, "false", updateRepo.updates[SettingKeyPasskeyEnabled])

	runtimeRepo.values[SettingKeyPasskeyEnabled] = "false"
	enabled, err = runtimeService.PasskeyEnabled(context.Background())
placeholder
	require.False(t, enabled)
	publicSettings, err := runtimeService.GetPublicSettings(context.Background())
placeholder
	require.False(t, publicSettings.PasskeyEnabled)
placeholder

// 移除 WebAuthn 配置后，残留的 passkey_enabled="true" 不得再让 GetAllSettings
// 报告开关开启：admin 更新门控以此为准，一旦误报为 true 会拒绝所有设置保存，
// 而此时前端开关处于禁用态，管理员无法在 UI 里自救。
func TestSettingService_StalePasskeyTrueWithoutConfigReportsDisabled(t *testing.T) {
	repo := &settingGetAllRepoStub{values: map[string]string{
		SettingKeyPasskeyEnabled: "true",
placeholderplaceholder
	service := NewSettingService(repo, &config.Config{placeholder)

	settings, err := service.GetAllSettings(context.Background())
placeholder
	require.False(t, settings.PasskeyEnabled)
placeholder
