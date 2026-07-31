package securityaudit

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

type prefixEncryptor struct{placeholder

func (prefixEncryptor) Encrypt(value string) (string, error) { return "enc:" + value, nil placeholder
func (prefixEncryptor) Decrypt(value string) (string, error) {
	if !strings.HasPrefix(value, "enc:") {
		return "", errors.New("cipher: message authentication failed")
placeholder
	return value[4:], nil
placeholder

// testTotpKeyConfig mirrors a deployment with a fixed TOTP_ENCRYPTION_KEY so
// unit tests may persist endpoint tokens.
func testTotpKeyConfig() *config.Config {
	return &config.Config{Totp: config.TotpConfig{EncryptionKeyConfigured: trueplaceholderplaceholder
placeholder

func TestDefaultConfigIsOff(t *testing.T) {
	storage, err := ParseStorageConfig("")
placeholder
	require.False(t, storage.Enabled)
	require.False(t, storage.BlockingLatestTurnOnly)
	active, err := ActiveFromStorage(storage, true, prefixEncryptor{placeholder)
placeholder
	require.Equal(t, ModeOff, active.EffectiveMode())
	require.Equal(t, AllScannerIDs, storage.Scanners)
	publicJSON, err := json.Marshal(PublicFromStorage(storage, true, nil))
placeholder
	require.Contains(t, string(publicJSON), `"group_ids":[]`)
	require.Contains(t, string(publicJSON), `"endpoints":[]`)
placeholder

func TestBlockingLatestTurnOnlyConfigRoundTrip(t *testing.T) {
	manager := &ConfigManager{encryptor: prefixEncryptor{placeholder, encryptionKeyConfigured: trueplaceholder
	request := UpdateConfigRequest{
		ExpectedConfigVersion: 1, Enabled: true, BlockingEnabled: true, BlockingLatestTurnOnly: true,
		Strategy: "priority", WorkerCount: 1, QueueCapacity: 10, Scanners: []string{"pii"placeholder, AllGroups: true,
		Endpoints: []UpdateEndpoint{{
			ID: "guard-1", Name: "Guard", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080",
			Model: DefaultGuardModel, TimeoutMS: 1000, InputLimit: 1000, Enabled: true,
placeholder
placeholder
	next, err := manager.buildNextStorage(DefaultStorageConfig(), request, 9)
placeholder
	require.True(t, next.BlockingLatestTurnOnly)
	require.Contains(t, changeSummary(next), `"blocking_latest_turn_only":true`)

	active, err := ActiveFromStorage(next, true, prefixEncryptor{placeholder)
placeholder
	require.True(t, active.BlockingLatestTurnOnly)
	public := PublicFromStorage(next, true, nil)
	require.True(t, public.BlockingLatestTurnOnly)
placeholder

func TestConfigRejectsBlockingWithoutAudit(t *testing.T) {
	storage := DefaultStorageConfig()
	storage.BlockingEnabled = true
	require.Error(t, validateStorageConfig(storage))
placeholder

func TestPublicConfigNeverMarshalsToken(t *testing.T) {
	storage := DefaultStorageConfig()
	storage.Endpoints = []StorageEndpoint{{ID: "one", Name: "One", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080", Model: DefaultGuardModel, TokenCiphertext: "GUARD_TOKEN_CANARY_SECRET", TimeoutMS: 1000, InputLimit: 1000, Enabled: trueplaceholderplaceholder
	public := PublicFromStorage(storage, true, nil)
	raw, err := json.Marshal(public)
placeholder
	require.NotContains(t, string(raw), "GUARD_TOKEN_CANARY_SECRET")
	require.NotContains(t, string(raw), "ciphertext")
	require.True(t, public.Endpoints[0].HasToken)
placeholder

func TestConfigRuntimeLoadErrorIsStableBoundedAndSecretFree(t *testing.T) {
	const canary = "CONFIG_LOAD_CANARY_SECRET"
	manager := &ConfigManager{clock: fixedClock{placeholderplaceholder
	manager.recordLoadError(errors.New("decrypt failed for token " + canary + " Authorization: Bearer " + canary))
	_, _, _, message := manager.RuntimeState()
	require.Equal(t, stableErrorMessage("config_load_failed"), message)
	require.NotContains(t, message, canary)
	require.LessOrEqual(t, len([]rune(message)), 160)
placeholder

func TestConfigManagerPublicRequiresSuccessfullyLoadedSnapshot(t *testing.T) {
	t.Run("absent persisted setting is legitimate default", func(t *testing.T) {
		manager := NewConfigManager(nil, staticSettingRepository{values: map[string]string{
			SettingKeyPromptAuditConfig: "",
			SettingKeyRiskControl:       "false",
placeholder nil, prefixEncryptor{placeholder, testTotpKeyConfig())
		require.NoError(t, manager.Reload(context.Background()))

		public, err := manager.Public()
	placeholder
		require.Equal(t, int64(1), public.ConfigVersion)
		require.False(t, public.Enabled)
placeholder)

	t.Run("unparseable persisted config is unavailable", func(t *testing.T) {
		const canary = "persisted-token-canary"
		manager := NewConfigManager(nil, staticSettingRepository{values: map[string]string{
			// Endpoint without id/name fails validation, so no trustworthy
			// snapshot can be installed from this raw value.
			SettingKeyPromptAuditConfig: `{"enabled":true,"config_version":9,"endpoints":[{"token_ciphertext":"` + canary + `"placeholder]placeholder`,
			SettingKeyRiskControl:       "true",
placeholder nil, prefixEncryptor{placeholder, testTotpKeyConfig())
		require.Error(t, manager.Reload(context.Background()))

		public, err := manager.Public()
	placeholder
		require.Empty(t, public)
		require.Equal(t, ErrorCodeConfigUnavailable, infraerrors.Reason(err))
		require.NotContains(t, err.Error(), canary)
placeholder)

	t.Run("reload failure preserves last successfully loaded snapshot", func(t *testing.T) {
		storage := DefaultStorageConfig()
		storage.ConfigVersion = 4
		storage.ChangeSummary = "trusted snapshot"
		raw, err := json.Marshal(storage)
	placeholder
		repository := &switchableSettingRepository{staticSettingRepository: staticSettingRepository{values: map[string]string{
			SettingKeyPromptAuditConfig: string(raw),
			SettingKeyRiskControl:       "false",
	placeholderplaceholderplaceholder
		manager := NewConfigManager(nil, repository, nil, prefixEncryptor{placeholder, testTotpKeyConfig())
		require.NoError(t, manager.Reload(context.Background()))
		repository.loadErr = errors.New("settings unavailable")
		require.Error(t, manager.Reload(context.Background()))

		public, err := manager.Public()
	placeholder
		require.Equal(t, int64(4), public.ConfigVersion)
		require.Equal(t, "trusted snapshot", public.ChangeSummary)
placeholder)
placeholder

// Regression coverage for issue #4887: a persisted config whose endpoint token
// can no longer be decrypted (encryption key changed or auto-generated per
// boot) must stay visible and editable for admins instead of falling back to a
// default v1 config that makes every save fail the CAS version check.
func TestConfigManagerUndecryptableTokenKeepsConfigVisibleAndRecoverable(t *testing.T) {
	const canary = "persisted-token-canary"
	persisted := `{"enabled":true,"blocking_enabled":false,"config_version":9,"endpoints":[{"id":"g1","name":"Guard","protocol":"openai_compatible","base_url":"http://127.0.0.1:8080","model":"m","token_ciphertext":"` + canary + `","timeout_ms":1000,"input_limit":1000,"enabled":trueplaceholder]placeholder`
	manager := NewConfigManager(nil, staticSettingRepository{values: map[string]string{
		SettingKeyPromptAuditConfig: persisted,
		SettingKeyRiskControl:       "true",
placeholderplaceholder, nil, prefixEncryptor{placeholder, testTotpKeyConfig())
	require.NoError(t, manager.Reload(context.Background()), "an undecryptable token must not fail the whole config load")

	public, err := manager.Public()
placeholder
	require.Equal(t, int64(9), public.ConfigVersion, "admins must see the real persisted version so CAS saves can succeed")
	require.Len(t, public.Endpoints, 1)
	require.True(t, public.Endpoints[0].HasToken)
	require.Equal(t, "invalid", public.Endpoints[0].TokenStatus)
	raw, err := json.Marshal(public)
placeholder
	require.NotContains(t, string(raw), canary)

	active, ok := manager.Active()
	require.True(t, ok)
	require.Len(t, active.Endpoints, 1)
	require.False(t, active.Endpoints[0].Enabled, "an endpoint with an undecryptable token must not be used at runtime")
	require.True(t, active.Endpoints[0].TokenInvalid)
	require.Empty(t, active.Endpoints[0].Token)
	require.Empty(t, active.EnabledEndpoints())
	require.Equal(t, []string{"g1"placeholder, active.InvalidTokenEndpointIDs())

	expected, activeVersion, _, _ := manager.RuntimeState()
	require.Equal(t, int64(9), expected)
	require.Equal(t, int64(9), activeVersion)
placeholder

func TestConfigManagerUndecryptableTokenStillFailsClosedForBlockingIntent(t *testing.T) {
	persisted := `{"enabled":true,"blocking_enabled":true,"config_version":9,"endpoints":[{"id":"g1","name":"Guard","protocol":"openai_compatible","base_url":"http://127.0.0.1:8080","model":"m","token_ciphertext":"undecryptable","timeout_ms":1000,"input_limit":1000,"enabled":trueplaceholder]placeholder`
	manager := NewConfigManager(nil, staticSettingRepository{values: map[string]string{
		SettingKeyPromptAuditConfig: persisted,
		SettingKeyRiskControl:       "true",
placeholderplaceholder, nil, prefixEncryptor{placeholder, testTotpKeyConfig())
	require.NoError(t, manager.Reload(context.Background()))
	require.Equal(t, ModeBlocking, manager.EffectiveMode())

	service := &PromptService{config: manager, evaluator: NewGuardEvaluator(NewOpenAICompatibleScanner(), nil, nil)placeholder
	decision, err := service.Evaluate(context.Background(), Request{
		Protocol: "openai_chat_completions",
		Body:     []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
placeholder)
	require.Error(t, err, "blocking intent with no usable endpoint must not let requests pass unaudited")
	require.Nil(t, decision)
	var guardErr *GuardError
	require.ErrorAs(t, err, &guardErr)
	require.Equal(t, ErrorCodeUnavailable, guardErr.Code)
placeholder

func TestBuildNextStoragePreserveReplaceAndClearToken(t *testing.T) {
	manager := &ConfigManager{encryptor: prefixEncryptor{placeholder, encryptionKeyConfigured: trueplaceholder
	current := DefaultStorageConfig()
	current.Endpoints = []StorageEndpoint{{ID: "one", Name: "One", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080", Model: DefaultGuardModel, TokenCiphertext: "enc:old", TimeoutMS: 1000, InputLimit: 1000placeholderplaceholder
	base := UpdateConfigRequest{ExpectedConfigVersion: 1, Strategy: "priority", WorkerCount: 1, QueueCapacity: 10, Scanners: []string{"PII"placeholder, AllGroups: true,
		Endpoints: []UpdateEndpoint{{ID: "one", Name: "One", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080", TimeoutMS: 1000, InputLimit: 1000placeholderplaceholderplaceholder
	preserved, err := manager.buildNextStorage(current, base, 9)
placeholder
	require.Equal(t, "enc:old", preserved.Endpoints[0].TokenCiphertext)
	replacedReq := base
	replacedReq.Endpoints = append([]UpdateEndpoint(nil), base.Endpoints...)
	replacedReq.Endpoints[0].Token = "new"
	replaced, err := manager.buildNextStorage(current, replacedReq, 9)
placeholder
	require.Equal(t, "enc:new", replaced.Endpoints[0].TokenCiphertext)
	clearedReq := base
	clearedReq.Endpoints = append([]UpdateEndpoint(nil), base.Endpoints...)
	clearedReq.Endpoints[0].ClearToken = true
	cleared, err := manager.buildNextStorage(current, clearedReq, 9)
placeholder
	require.Empty(t, cleared.Endpoints[0].TokenCiphertext)
placeholder

// Without a fixed encryption key the per-boot auto-generated key would make a
// freshly saved token undecryptable after the next restart (issue #4887), so
// saving a new token must be rejected with an actionable error. Preserving or
// clearing an existing ciphertext stays allowed so admins can still edit or
// disable the feature.
func TestBuildNextStorageRejectsNewTokenWithoutConfiguredEncryptionKey(t *testing.T) {
	manager := &ConfigManager{encryptor: prefixEncryptor{placeholder, encryptionKeyConfigured: falseplaceholder
	current := DefaultStorageConfig()
	current.Endpoints = []StorageEndpoint{{ID: "one", Name: "One", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080", Model: DefaultGuardModel, TokenCiphertext: "enc:old", TimeoutMS: 1000, InputLimit: 1000placeholderplaceholder
	base := UpdateConfigRequest{ExpectedConfigVersion: 1, Strategy: "priority", WorkerCount: 1, QueueCapacity: 10, Scanners: []string{"PII"placeholder, AllGroups: true,
		Endpoints: []UpdateEndpoint{{ID: "one", Name: "One", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080", TimeoutMS: 1000, InputLimit: 1000placeholderplaceholderplaceholder

	newTokenReq := base
	newTokenReq.Endpoints = append([]UpdateEndpoint(nil), base.Endpoints...)
	newTokenReq.Endpoints[0].Token = "fresh-token"
	_, err := manager.buildNextStorage(current, newTokenReq, 9)
placeholder
	require.Equal(t, ErrorCodeEncryptionKeyRequired, infraerrors.Reason(err))

	preserved, err := manager.buildNextStorage(current, base, 9)
placeholder
	require.Equal(t, "enc:old", preserved.Endpoints[0].TokenCiphertext)

	clearedReq := base
	clearedReq.Endpoints = append([]UpdateEndpoint(nil), base.Endpoints...)
	clearedReq.Endpoints[0].ClearToken = true
	cleared, err := manager.buildNextStorage(current, clearedReq, 9)
placeholder
	require.Empty(t, cleared.Endpoints[0].TokenCiphertext)
placeholder

func TestEffectiveModeTruthTable(t *testing.T) {
	tests := []struct {
		risk, enabled, blocking bool
		want                    Mode
placeholder{
		{false, false, false, ModeOffplaceholder, {false, true, true, ModeOffplaceholder, {true, false, false, ModeOffplaceholder,
		{true, true, false, ModeAsyncplaceholder, {true, true, true, ModeBlockingplaceholder,
placeholder
	for _, tt := range tests {
		cfg := ActiveConfig{RiskControlEnabled: tt.risk, Enabled: tt.enabled, BlockingEnabled: tt.blockingplaceholder
		require.Equal(t, tt.want, cfg.EffectiveMode())
placeholder
placeholder

func TestConfigManagerColdStartOnlyFailsClosedForExplicitBlockingIntent(t *testing.T) {
	manager := &ConfigManager{placeholder

	manager.observeExpectedState(`{"enabled":true,"blocking_enabled":false,"config_version":42placeholder`, true)
	require.Equal(t, int64(42), manager.expected.Load())
	require.Equal(t, ModeOff, manager.EffectiveMode(), "an async config version must not imply blocking")
	require.False(t, manager.BlockingActivationDegraded())

	manager.observeExpectedState(`{"enabled":true,"blocking_enabled":true,"config_version":43placeholder`, false)
	require.Equal(t, ModeOff, manager.EffectiveMode(), "the global risk-control switch still gates blocking")

	manager.observeExpectedState(`{"enabled":true,"blocking_enabled":true,"config_version":44placeholder`, true)
	require.Equal(t, ModeBlocking, manager.EffectiveMode())
	require.True(t, manager.BlockingActivationDegraded())

	manager.observeExpectedState(`{"enabled":true`, true)
	require.Equal(t, ModeBlocking, manager.EffectiveMode(), "undecodable storage must not erase the last known strict intent")
placeholder

func TestConfigManagerStaleWeakerSnapshotFailsClosedWhenBlockingExpected(t *testing.T) {
	manager := &ConfigManager{placeholder
	async := ActiveConfig{RiskControlEnabled: true, Enabled: true, BlockingEnabled: false, ConfigVersion: 1placeholder
	manager.snapshot.Store(&activeConfigSnapshot{active: async, storage: DefaultStorageConfig(), loadedAt: fixedClock{placeholder.Now()placeholder)
	manager.expected.Store(2)
	manager.expectedBlocking.Store(true)

	require.True(t, manager.BlockingActivationDegraded())
	require.Equal(t, ModeBlocking, manager.EffectiveMode())

	service := &PromptService{config: manager, evaluator: NewGuardEvaluator(nil, nil, nil)placeholder
	decision, err := service.Evaluate(context.Background(), Request{Protocol: "openai_chat_completions", Body: []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`)placeholder)
placeholder
	require.Nil(t, decision)
	var guardErr *GuardError
	require.ErrorAs(t, err, &guardErr)
	require.Equal(t, ErrorCodeUnavailable, guardErr.Code)
placeholder

type errorSettingRepository struct{ staticSettingRepository placeholder

func (errorSettingRepository) GetMultiple(context.Context, []string) (map[string]string, error) {
	return nil, errors.New("settings unavailable")
placeholder

type switchableSettingRepository struct {
	staticSettingRepository
	loadErr error
placeholder

func (r *switchableSettingRepository) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	if r.loadErr != nil {
		return nil, r.loadErr
placeholder
	return r.staticSettingRepository.GetMultiple(ctx, keys)
placeholder

func TestConfigManagerStartupLoadFailureDoesNotBlockWhenBlockingNotIntended(t *testing.T) {
	// Settings unavailable and no prior blocking intent: stay ModeOff so the
	// gateway remains usable and admins can still disable/configure Prompt Audit.
	manager := NewConfigManager(nil, errorSettingRepository{placeholder, nil, prefixEncryptor{placeholder, testTotpKeyConfig())
	err := manager.Start(context.Background())
placeholder
	require.True(t, manager.configUntrusted.Load())
	require.False(t, manager.BlockingActivationDegraded())
	require.Equal(t, ModeOff, manager.EffectiveMode())

	service := &PromptService{config: manager, evaluator: NewGuardEvaluator(nil, nil, nil)placeholder
	decision, evalErr := service.Evaluate(context.Background(), Request{
		Protocol: "openai_chat_completions",
		Body:     []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
placeholder)
	require.NoError(t, evalErr)
	require.NotNil(t, decision)
	require.Equal(t, DecisionAllow, decision.Kind)
	require.NoError(t, manager.Shutdown(context.Background()))
placeholder

func TestConfigManagerStartupLoadFailureFailsClosedWhenBlockingIntended(t *testing.T) {
	manager := NewConfigManager(nil, errorSettingRepository{placeholder, nil, prefixEncryptor{placeholder, testTotpKeyConfig())
	// Simulate intent observed before a later load failure (e.g. decrypt error).
	manager.observeExpectedState(`{"enabled":true,"blocking_enabled":true,"config_version":3placeholder`, true)
	manager.markConfigUntrusted()
	require.True(t, manager.BlockingActivationDegraded())
	require.Equal(t, ModeBlocking, manager.EffectiveMode())

	service := &PromptService{config: manager, evaluator: NewGuardEvaluator(nil, nil, nil)placeholder
	decision, err := service.Evaluate(context.Background(), Request{
		Protocol: "openai_chat_completions",
		Body:     []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
placeholder)
placeholder
	require.Nil(t, decision)
	var guardErr *GuardError
	require.ErrorAs(t, err, &guardErr)
	require.Equal(t, ErrorCodeUnavailable, guardErr.Code)
placeholder

func TestConfigManagerUntrustedClearsOnSuccessfulDisable(t *testing.T) {
	// After a degraded fail-closed period, saving disabled config must restore ModeOff.
	manager := &ConfigManager{encryptor: prefixEncryptor{placeholder, clock: fixedClock{placeholderplaceholder
	manager.observeExpectedState(`{"enabled":true,"blocking_enabled":true,"config_version":5placeholder`, true)
	manager.markConfigUntrusted()
	require.Equal(t, ModeBlocking, manager.EffectiveMode())

	// Install a trusted disabled snapshot the same way Save does after commit.
	disabled := DefaultStorageConfig()
	disabled.ConfigVersion = 6
	disabled.Enabled = false
	disabled.BlockingEnabled = false
	active, err := ActiveFromStorage(disabled, true, manager.encryptor)
placeholder
	manager.expected.Store(disabled.ConfigVersion)
	manager.expectedBlocking.Store(false)
	manager.snapshot.Store(&activeConfigSnapshot{storage: disabled, active: active, loadedAt: manager.clock.Now()placeholder)
	manager.configUntrusted.Store(false)

	require.False(t, manager.BlockingActivationDegraded())
	require.Equal(t, ModeOff, manager.EffectiveMode())

	service := &PromptService{config: manager, evaluator: NewGuardEvaluator(nil, nil, nil)placeholder
	decision, evalErr := service.Evaluate(context.Background(), Request{
		Protocol: "openai_chat_completions",
		Body:     []byte(`{"messages":[{"role":"user","content":"hi"placeholder]placeholder`),
placeholder)
	require.NoError(t, evalErr)
	require.Equal(t, DecisionAllow, decision.Kind)
placeholder

func TestConfigManagerUntrustedWithoutBlockingDoesNotForceBlockingMode(t *testing.T) {
	manager := &ConfigManager{placeholder
	manager.observeExpectedState(`{"enabled":true,"blocking_enabled":false,"config_version":2placeholder`, true)
	manager.markConfigUntrusted()
	require.False(t, manager.expectedBlocking.Load())
	require.False(t, manager.BlockingActivationDegraded())
	require.Equal(t, ModeOff, manager.EffectiveMode(), "async intent + untrusted must not force blocking unavailable")
placeholder

func TestParseLegacyConfigDefaultsMissingFieldsWithoutEnablingBlocking(t *testing.T) {
	storage, err := ParseStorageConfig(`{"enabled":false,"config_version":9placeholder`)
placeholder
	require.False(t, storage.BlockingEnabled)
	require.Equal(t, "priority", storage.Strategy)
	require.Equal(t, DefaultWorkerCount, storage.WorkerCount)
	require.Equal(t, DefaultQueueCapacity, storage.QueueCapacity)
	require.Equal(t, AllScannerIDs, storage.Scanners)
	require.True(t, storage.AllGroups)
placeholder

func TestUpdateConfigStrictBoundsAndKnownValues(t *testing.T) {
	valid := promptAuditUpdateRequest(1, 1, "")
	require.NoError(t, validateUpdateConfigRequest(valid))

	tests := []struct {
		name   string
		mutate func(*UpdateConfigRequest)
		reason string
placeholder{
		{name: "strategy", mutate: func(req *UpdateConfigRequest) { req.Strategy = "round_robin" placeholder, reason: "prompt_audit_invalid_strategy"placeholder,
		{name: "worker low", mutate: func(req *UpdateConfigRequest) { req.WorkerCount = 0 placeholder, reason: "prompt_audit_invalid_worker_count"placeholder,
		{name: "worker high", mutate: func(req *UpdateConfigRequest) { req.WorkerCount = MaxWorkerCount + 1 placeholder, reason: "prompt_audit_invalid_worker_count"placeholder,
		{name: "capacity low", mutate: func(req *UpdateConfigRequest) { req.QueueCapacity = 0 placeholder, reason: "prompt_audit_invalid_queue_capacity"placeholder,
		{name: "capacity high", mutate: func(req *UpdateConfigRequest) { req.QueueCapacity = MaxQueueCapacity + 1 placeholder, reason: "prompt_audit_invalid_queue_capacity"placeholder,
		{name: "unknown scanner", mutate: func(req *UpdateConfigRequest) { req.Scanners = []string{"made_up"placeholder placeholder, reason: "prompt_audit_invalid_scanner"placeholder,
		{name: "group required", mutate: func(req *UpdateConfigRequest) { req.AllGroups = false; req.GroupIDs = nil placeholder, reason: "prompt_audit_groups_required"placeholder,
		{name: "group positive", mutate: func(req *UpdateConfigRequest) { req.AllGroups = false; req.GroupIDs = []int64{0placeholder placeholder, reason: "prompt_audit_invalid_group"placeholder,
		{name: "timeout low", mutate: func(req *UpdateConfigRequest) { req.Endpoints[0].TimeoutMS = MinTimeoutMS - 1 placeholder, reason: "prompt_audit_invalid_timeout"placeholder,
		{name: "timeout high", mutate: func(req *UpdateConfigRequest) { req.Endpoints[0].TimeoutMS = MaxTimeoutMS + 1 placeholder, reason: "prompt_audit_invalid_timeout"placeholder,
		{name: "input low", mutate: func(req *UpdateConfigRequest) { req.Endpoints[0].InputLimit = MinInputLimit - 1 placeholder, reason: "prompt_audit_invalid_input_limit"placeholder,
		{name: "input high", mutate: func(req *UpdateConfigRequest) { req.Endpoints[0].InputLimit = MaxInputLimit + 1 placeholder, reason: "prompt_audit_invalid_input_limit"placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := valid
			req.Scanners = append([]string(nil), valid.Scanners...)
			req.GroupIDs = append([]int64(nil), valid.GroupIDs...)
			req.Endpoints = append([]UpdateEndpoint(nil), valid.Endpoints...)
			tt.mutate(&req)
			err := validateUpdateConfigRequest(req)
		placeholder
			require.Equal(t, tt.reason, infraerrors.Reason(err))
	placeholder)
placeholder
placeholder
