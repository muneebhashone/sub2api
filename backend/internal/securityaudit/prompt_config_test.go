package securityaudit

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/stretchr/testify/require"
)

type prefixEncryptor struct{placeholder

func (prefixEncryptor) Encrypt(value string) (string, error) { return "enc:" + value, nil placeholder
func (prefixEncryptor) Decrypt(value string) (string, error) { return value[4:], nil placeholder

func TestDefaultConfigIsOff(t *testing.T) {
	storage, err := ParseStorageConfig("")
placeholder
	require.False(t, storage.Enabled)
	active, err := ActiveFromStorage(storage, true, prefixEncryptor{placeholder)
placeholder
	require.Equal(t, ModeOff, active.EffectiveMode())
	require.Equal(t, AllScannerIDs, storage.Scanners)
	publicJSON, err := json.Marshal(PublicFromStorage(storage, true))
placeholder
	require.Contains(t, string(publicJSON), `"group_ids":[]`)
	require.Contains(t, string(publicJSON), `"endpoints":[]`)
placeholder

func TestConfigRejectsBlockingWithoutAudit(t *testing.T) {
	storage := DefaultStorageConfig()
	storage.BlockingEnabled = true
	require.Error(t, validateStorageConfig(storage))
placeholder

func TestPublicConfigNeverMarshalsToken(t *testing.T) {
	storage := DefaultStorageConfig()
	storage.Endpoints = []StorageEndpoint{{ID: "one", Name: "One", Protocol: "openai_compatible", BaseURL: "http://127.0.0.1:8080", Model: DefaultGuardModel, TokenCiphertext: "GUARD_TOKEN_CANARY_SECRET", TimeoutMS: 1000, InputLimit: 1000, Enabled: trueplaceholderplaceholder
	public := PublicFromStorage(storage, true)
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

func TestBuildNextStoragePreserveReplaceAndClearToken(t *testing.T) {
	manager := &ConfigManager{encryptor: prefixEncryptor{placeholderplaceholder
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
