package service

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	pluginv1 "github.com/Wei-Shaw/sub2api/pkg/pluginapi/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type pluginTokenRepository struct {
	PluginRepository
	installation *PluginInstallation
	listErr      error
placeholder

func (r *pluginTokenRepository) List(context.Context) ([]*PluginInstallation, error) {
	if r.listErr != nil {
		return nil, r.listErr
placeholder
	if r.installation == nil {
		return nil, nil
placeholder
	copy := *r.installation
	return []*PluginInstallation{&copyplaceholder, nil
placeholder

func (r *pluginTokenRepository) GetByID(context.Context, int64) (*PluginInstallation, error) {
	if r.installation == nil {
		return nil, errors.New("插件不存在")
placeholder
	copy := *r.installation
	return &copy, nil
placeholder

type pluginTokenEncryptor struct{placeholder

func (pluginTokenEncryptor) Encrypt(plaintext string) (string, error) {
	return "ENC:" + plaintext, nil
placeholder

func (pluginTokenEncryptor) Decrypt(ciphertext string) (string, error) {
	if !strings.HasPrefix(ciphertext, "ENC:") {
		return "", errors.New("密文无效")
placeholder
	return strings.TrimPrefix(ciphertext, "ENC:"), nil
placeholder

func TestPluginUIAssetTokenCanBeResolvedByAnotherInstance(t *testing.T) {
	repo := &pluginTokenRepository{installation: &PluginInstallation{ID: 42placeholderplaceholder
	first := &PluginManager{repo: repo, encryptor: pluginTokenEncryptor{placeholderplaceholder
	second := &PluginManager{repo: repo, encryptor: pluginTokenEncryptor{placeholderplaceholder

	token, expires, err := first.CreateUIAssetToken(context.Background(), 42, 30*time.Minute)
placeholder
	require.WithinDuration(t, time.Now().Add(30*time.Minute), expires, time.Second)

	pluginID, err := second.ResolveUIAssetToken(token)
placeholder
	require.Equal(t, int64(42), pluginID)

	decoded, err := base64.RawURLEncoding.DecodeString(token)
placeholder
	decoded[len(decoded)-1] ^= 1
	_, err = second.ResolveUIAssetToken(base64.RawURLEncoding.EncodeToString(decoded))
placeholder
placeholder

func TestPluginUIAssetTokenRejectsOtherEncryptedPayloads(t *testing.T) {
	repo := &pluginTokenRepository{installation: &PluginInstallation{ID: 42placeholderplaceholder
	manager := &PluginManager{repo: repo, encryptor: pluginTokenEncryptor{placeholderplaceholder

	encrypted, err := manager.encryptor.Encrypt(`{"version":1,"plugin_id":42,"expires":4102444800placeholder`)
placeholder
	token := base64.RawURLEncoding.EncodeToString([]byte(encrypted))

	_, err = manager.ResolveUIAssetToken(token)
	require.ErrorContains(t, err, "会话无效")
placeholder

func TestPluginReconcileFailsClosedWhenDesiredStateCannotBeRead(t *testing.T) {
	manager := &PluginManager{
		repo:               &pluginTokenRepository{listErr: errors.New("数据库不可用")placeholder,
		runtimes:           make(map[int64]*pluginRuntime),
		localInstallations: make(map[int64]*PluginInstallation),
placeholder

	err := manager.reconcileOnce(context.Background())
	require.ErrorContains(t, err, "读取插件启用状态")
	require.True(t, manager.ShouldRouteOpenAIOAuth(&Account{
		ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
placeholder))

	request, requestErr := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://example.com/v1/responses", nil)
	require.NoError(t, requestErr)
	_, handled, routeErr := manager.RoundTripOpenAIOAuth(context.Background(), request, "", &Account{
		ID: 1, Platform: PlatformOpenAI, Type: AccountTypeOAuth,
placeholder)
	require.True(t, handled)
	require.ErrorContains(t, routeErr, "插件不可用")
placeholder

type normalizingPluginClient struct {
	pluginv1.TransportPluginClient
	normalized []byte
	applied    []byte
placeholder

type pluginConfigRepository struct {
	PluginRepository
	installation *PluginInstallation
	encrypted    string
placeholder

func (r *pluginConfigRepository) GetByID(context.Context, int64) (*PluginInstallation, error) {
	copy := *r.installation
	return &copy, nil
placeholder

func (r *pluginConfigRepository) UpdateConfig(_ context.Context, _ int64, encrypted, expectedBinarySHA256 string) error {
	if expectedBinarySHA256 != r.installation.BinarySHA256 {
		return ErrPluginStateChanged
placeholder
	r.encrypted = encrypted
	return nil
placeholder

func (c *normalizingPluginClient) ValidateConfig(context.Context, *pluginv1.ValidateConfigRequest, ...grpc.CallOption) (*pluginv1.ValidateConfigResponse, error) {
	return &pluginv1.ValidateConfigResponse{Valid: true, NormalizedConfigJson: c.normalizedplaceholder, nil
placeholder

func (c *normalizingPluginClient) ApplyConfig(_ context.Context, request *pluginv1.ApplyConfigRequest, _ ...grpc.CallOption) (*pluginv1.ApplyConfigResponse, error) {
	c.applied = append([]byte(nil), request.ConfigJson...)
	return &pluginv1.ApplyConfigResponse{Applied: trueplaceholder, nil
placeholder

func TestPluginRuntimeReturnsAndAppliesNormalizedConfig(t *testing.T) {
	client := &normalizingPluginClient{normalized: []byte(`{"z":2,"a":1placeholder`)placeholder
	runtime := &pluginRuntime{api: clientplaceholder

	normalized, err := runtime.validateAndApplyNormalizedConfig(context.Background(), []byte(`{"input":trueplaceholder`))
placeholder
	require.JSONEq(t, `{"a":1,"z":2placeholder`, string(normalized))
	require.Equal(t, normalized, client.applied)
placeholder

func TestPluginRuntimeRejectsInvalidNormalizedConfig(t *testing.T) {
	client := &normalizingPluginClient{normalized: []byte(`{"broken"`)placeholder
	runtime := &pluginRuntime{api: clientplaceholder

	_, err := runtime.validateAndApplyNormalizedConfig(context.Background(), []byte(`{placeholder`))
	require.ErrorContains(t, err, "规范化配置")
	require.Empty(t, client.applied)
placeholder

func TestPluginManagerPersistsPluginNormalizedConfig(t *testing.T) {
	installation := &PluginInstallation{ID: 9, BinarySHA256: strings.Repeat("a", 64)placeholder
	repo := &pluginConfigRepository{installation: installationplaceholder
	client := &normalizingPluginClient{normalized: []byte(`{"timeout":30,"enabled":trueplaceholder`)placeholder
	manager := &PluginManager{
		repo: repo, encryptor: pluginTokenEncryptor{placeholder,
		runtimes: map[int64]*pluginRuntime{9: {installation: installation, api: clientplaceholderplaceholder,
placeholder

	saved, err := manager.SaveConfig(context.Background(), 9, []byte(`{"enabled":falseplaceholder`))
placeholder
	require.JSONEq(t, `{"enabled":true,"timeout":30placeholder`, string(saved))
	plaintext, err := (pluginTokenEncryptor{placeholder).Decrypt(repo.encrypted)
placeholder
	require.JSONEq(t, string(saved), plaintext)
placeholder

func TestPluginRequestSentErrorDoesNotFailOver(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("POST", "/v1/responses", nil)
	account := &Account{ID: 7, Name: "oauth", Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	transportErr := &PluginTransportError{Code: "UPSTREAM_EOF", Message: "eof", RequestSent: trueplaceholder

	result := (&OpenAIGatewayService{placeholder).handleOpenAIUpstreamTransportError(context.Background(), c, account, transportErr, true)

	require.Same(t, transportErr, result)
	var failover *UpstreamFailoverError
	require.False(t, errors.As(result, &failover))
placeholder

func TestPluginRPCAmbiguityPreventsReplayAfterMetadataDelivery(t *testing.T) {
	err := normalizePluginRPCError(context.Background(), "接收插件响应头", errors.New("连接已断开"), true)
	var transportErr *PluginTransportError
	require.ErrorAs(t, err, &transportErr)
	require.True(t, transportErr.RequestSent)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("POST", "/v1/responses", nil)
	account := &Account{ID: 12, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	result := (&OpenAIGatewayService{placeholder).handleOpenAIUpstreamTransportError(context.Background(), c, account, err, true)
	require.Same(t, err, result)
placeholder

func TestPluginRPCFailureBeforeStreamCreationAllowsFailover(t *testing.T) {
	err := normalizePluginRPCError(context.Background(), "创建插件转发流", errors.New("连接失败"), false)
	var transportErr *PluginTransportError
	require.ErrorAs(t, err, &transportErr)
	require.False(t, transportErr.RequestSent)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest("POST", "/v1/responses", nil)
	account := &Account{ID: 13, Platform: PlatformOpenAI, Type: AccountTypeOAuthplaceholder
	result := (&OpenAIGatewayService{placeholder).handleOpenAIUpstreamTransportError(context.Background(), c, account, err, true)
	var failover *UpstreamFailoverError
	require.ErrorAs(t, result, &failover)
placeholder

func TestNormalizePluginRPCErrorPreservesCallerCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := normalizePluginRPCError(ctx, "接收响应", errors.New("rpc error: code = Canceled"), true)
	require.ErrorIs(t, err, context.Canceled)
placeholder

func TestPluginStartingStateUsesBoundedCrashRecoveryWindow(t *testing.T) {
	manager := &PluginManager{cfg: &config.Config{Plugins: config.PluginConfig{StartTimeoutSeconds: 15placeholderplaceholderplaceholder

	require.False(t, manager.startingStateExpired(&PluginInstallation{UpdatedAt: time.Now().Add(-30 * time.Second)placeholder))
	require.True(t, manager.startingStateExpired(&PluginInstallation{UpdatedAt: time.Now().Add(-2 * time.Minute)placeholder))
placeholder
