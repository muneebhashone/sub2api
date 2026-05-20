package service

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotificationEmailPreviewEscapesHTMLAndSanitizesSubject(t *testing.T) {
	ctx := context.Background()
	svc := NewNotificationEmailService(newNotificationEmailMemorySettingRepo(), nil)

	preview, err := svc.PreviewTemplate(ctx, NotificationEmailPreviewInput{
		Event:   NotificationEmailEventBalanceLow,
		Locale:  "en-US,en;q=0.9",
		Subject: "Low balance for {{recipient_nameplaceholderplaceholder\r\nInjected",
		HTML:    `<p>{{recipient_nameplaceholderplaceholder</p><a href="{{recharge_urlplaceholderplaceholder">Recharge</a>`,
		Variables: map[string]string{
			"recipient_name": `<script>alert("x")</script>`,
			"recharge_url":   `javascript:alert(1)`,
	placeholder,
placeholder)
placeholder
	require.NotContains(t, preview.Subject, "\r")
	require.NotContains(t, preview.Subject, "\n")
	require.Contains(t, preview.Subject, `Low balance for <script>alert("x")</script>Injected`)
	require.Contains(t, preview.HTML, `&lt;script&gt;alert(&#34;x&#34;)&lt;/script&gt;`)
	require.NotContains(t, preview.HTML, `javascript:alert(1)`)
	require.Contains(t, preview.HTML, `href=""`)
placeholder

func TestNotificationEmailTemplateOverrideAndRestore(t *testing.T) {
	ctx := context.Background()
	repo := newNotificationEmailMemorySettingRepo()
	svc := NewNotificationEmailService(repo, nil)

	official, err := svc.GetTemplate(ctx, NotificationEmailEventBalanceRechargeSuccess, "en")
placeholder
	require.False(t, official.IsCustom)

	updated, err := svc.UpdateTemplate(
		ctx,
		NotificationEmailEventBalanceRechargeSuccess,
		"zh-Hans",
		"充值完成：{{recharge_amountplaceholderplaceholder",
		"<p>{{recipient_nameplaceholderplaceholder 已充值 {{recharge_amountplaceholderplaceholder</p>",
	)
placeholder
	require.True(t, updated.IsCustom)
	require.Equal(t, "zh", updated.Locale)
	require.Equal(t, "充值完成：{{recharge_amountplaceholderplaceholder", updated.Subject)
	require.NotNil(t, updated.UpdatedAt)

	restored, err := svc.RestoreOfficialTemplate(ctx, NotificationEmailEventBalanceRechargeSuccess, "zh")
placeholder
	require.False(t, restored.IsCustom)
	require.NotEqual(t, updated.Subject, restored.Subject)
	_, err = repo.GetValue(ctx, notificationEmailTemplateKey(NotificationEmailEventBalanceRechargeSuccess, "zh"))
	require.ErrorIs(t, err, ErrSettingNotFound)
placeholder

func TestNotificationEmailTemplateRejectsUnsupportedPlaceholder(t *testing.T) {
	ctx := context.Background()
	svc := NewNotificationEmailService(newNotificationEmailMemorySettingRepo(), nil)

	_, err := svc.UpdateTemplate(
		ctx,
		NotificationEmailEventSubscriptionPurchaseSuccess,
		"en",
		"Purchased {{not_allowedplaceholderplaceholder",
		"<p>{{subscription_groupplaceholderplaceholder</p>",
	)
placeholder
	require.Contains(t, err.Error(), "unsupported placeholder")
placeholder

func TestNotificationEmailUnsubscribeOnlyAllowsOptionalEvents(t *testing.T) {
	ctx := context.Background()
	svc := NewNotificationEmailService(newNotificationEmailMemorySettingRepo(), nil)

	token, err := svc.createUnsubscribeToken(ctx, "User@Example.com", NotificationEmailEventBalanceLow)
placeholder
	result, err := svc.Unsubscribe(ctx, token)
placeholder
	require.True(t, result.Done)
	require.Equal(t, NotificationEmailEventBalanceLow, result.Event)
	unsubscribed, err := svc.IsUnsubscribed(ctx, "user@example.com", NotificationEmailEventBalanceLow)
placeholder
	require.True(t, unsubscribed)

	transactionalToken, err := svc.createUnsubscribeToken(ctx, "user@example.com", NotificationEmailEventBalanceRechargeSuccess)
placeholder
	_, err = svc.Unsubscribe(ctx, transactionalToken)
placeholder
	require.Contains(t, err.Error(), "transactional")
placeholder

func TestNotificationEmailLocaleMemoryNormalizesAcceptLanguage(t *testing.T) {
	ctx := context.Background()
	svc := NewNotificationEmailService(newNotificationEmailMemorySettingRepo(), nil)

	svc.RememberRecipientLocale(ctx, 42, "User@Example.com", "zh-CN,zh;q=0.9,en;q=0.8")
	require.Equal(t, "zh", svc.ResolveRecipientLocale(ctx, 42, "user@example.com"))
	require.Equal(t, "zh", svc.ResolveRecipientLocale(ctx, 0, "user@example.com"))
placeholder

type notificationEmailMemorySettingRepo struct {
	mu     sync.RWMutex
	values map[string]string
placeholder

func newNotificationEmailMemorySettingRepo() *notificationEmailMemorySettingRepo {
	return &notificationEmailMemorySettingRepo{values: make(map[string]string)placeholder
placeholder

func (r *notificationEmailMemorySettingRepo) Get(_ context.Context, key string) (*Setting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.values[key]
	if !ok {
		return nil, ErrSettingNotFound
placeholder
	return &Setting{Key: key, Value: valueplaceholder, nil
placeholder

func (r *notificationEmailMemorySettingRepo) GetValue(ctx context.Context, key string) (string, error) {
	setting, err := r.Get(ctx, key)
	if err != nil {
		return "", err
placeholder
	return setting.Value, nil
placeholder

func (r *notificationEmailMemorySettingRepo) Set(_ context.Context, key, value string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[key] = value
	return nil
placeholder

func (r *notificationEmailMemorySettingRepo) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if value, ok := r.values[key]; ok {
			out[key] = value
	placeholder
placeholder
	return out, nil
placeholder

func (r *notificationEmailMemorySettingRepo) SetMultiple(_ context.Context, settings map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, value := range settings {
		r.values[key] = value
placeholder
	return nil
placeholder

func (r *notificationEmailMemorySettingRepo) GetAll(_ context.Context) (map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]string, len(r.values))
	for key, value := range r.values {
		out[key] = value
placeholder
	return out, nil
placeholder

func (r *notificationEmailMemorySettingRepo) Delete(_ context.Context, key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.values[key]; !ok {
		return ErrSettingNotFound
placeholder
	delete(r.values, key)
	return nil
placeholder

func TestNotificationEmailMemorySettingRepoSatisfiesInterface(t *testing.T) {
	var _ SettingRepository = (*notificationEmailMemorySettingRepo)(nil)
	require.False(t, strings.Contains(notificationEmailPreferenceKey(NotificationEmailEventBalanceLow, "User@Example.com"), "User@Example.com"))
placeholder
