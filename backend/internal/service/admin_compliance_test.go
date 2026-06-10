package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
)

type adminComplianceRepoStub struct {
	values map[string]string
placeholder

func (r *adminComplianceRepoStub) Get(ctx context.Context, key string) (*Setting, error) {
	if value, ok := r.values[key]; ok {
		return &Setting{Key: key, Value: valueplaceholder, nil
placeholder
	return nil, ErrSettingNotFound
placeholder

func (r *adminComplianceRepoStub) GetValue(ctx context.Context, key string) (string, error) {
	setting, err := r.Get(ctx, key)
	if err != nil {
		return "", err
placeholder
	return setting.Value, nil
placeholder

func (r *adminComplianceRepoStub) Set(ctx context.Context, key, value string) error {
	if r.values == nil {
		r.values = map[string]string{placeholder
placeholder
	r.values[key] = value
	return nil
placeholder

func (r *adminComplianceRepoStub) GetMultiple(ctx context.Context, keys []string) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder

func (r *adminComplianceRepoStub) SetMultiple(ctx context.Context, settings map[string]string) error {
	return nil
placeholder

func (r *adminComplianceRepoStub) GetAll(ctx context.Context) (map[string]string, error) {
	return map[string]string{placeholder, nil
placeholder

func (r *adminComplianceRepoStub) Delete(ctx context.Context, key string) error {
	delete(r.values, key)
	return nil
placeholder

func TestAdminComplianceStatusRequiresAckWhenMissing(t *testing.T) {
	svc := NewSettingService(&adminComplianceRepoStub{placeholder, &config.Config{placeholder)

	status, err := svc.GetAdminComplianceStatus(context.Background(), 1)
placeholder
	require.True(t, status.Required)
	require.Equal(t, AdminComplianceVersion, status.Version)
	require.Equal(t, AdminComplianceAckPhraseZH, status.AckPhraseZH)
	require.Equal(t, AdminComplianceDocumentPathZH, status.DocumentPathZH)
placeholder

func TestAcceptAdminComplianceRejectsWrongPhrase(t *testing.T) {
	svc := NewSettingService(&adminComplianceRepoStub{placeholder, &config.Config{placeholder)

	_, err := svc.AcceptAdminCompliance(context.Background(), AdminComplianceAcceptInput{
		AdminUserID: 1,
		Language:    "zh",
		Phrase:      "我同意",
placeholder)
placeholder
	require.True(t, errors.Is(err, ErrAdminComplianceInvalidPhrase))
placeholder

func TestAcceptAdminCompliancePersistsCurrentVersion(t *testing.T) {
	repo := &adminComplianceRepoStub{placeholder
	svc := NewSettingService(repo, &config.Config{placeholder)

	status, err := svc.AcceptAdminCompliance(context.Background(), AdminComplianceAcceptInput{
		AdminUserID: 42,
		Language:    "zh-CN",
		Phrase:      AdminComplianceAckPhraseZH,
		IPAddress:   "203.0.113.10",
		UserAgent:   "test-agent",
placeholder)
placeholder
	require.False(t, status.Required)
	require.NotNil(t, status.Acknowledgement)
	require.Equal(t, int64(42), status.Acknowledgement.AdminUserID)
	require.Equal(t, "203.0.113.10", status.Acknowledgement.IPAddress)

	var stored AdminComplianceAcknowledgement
	require.NoError(t, json.Unmarshal([]byte(repo.values[adminComplianceAcknowledgementKey(42)]), &stored))
	require.Equal(t, AdminComplianceVersion, stored.Version)
	require.Equal(t, AdminComplianceDocumentPathZH, stored.DocumentZH)
placeholder

func TestAdminComplianceStatusRequiresAckOnOldVersion(t *testing.T) {
	old, err := json.Marshal(AdminComplianceAcknowledgement{Version: "v2026.01.01"placeholder)
placeholder
	svc := NewSettingService(&adminComplianceRepoStub{
		values: map[string]string{adminComplianceAcknowledgementKey(1): string(old)placeholder,
placeholder, &config.Config{placeholder)

	status, err := svc.GetAdminComplianceStatus(context.Background(), 1)
placeholder
	require.True(t, status.Required)
	require.Nil(t, status.Acknowledgement)
placeholder

func TestAdminComplianceStatusIsPerAdminUser(t *testing.T) {
	current, err := json.Marshal(AdminComplianceAcknowledgement{
		Version:     AdminComplianceVersion,
		AdminUserID: 1,
placeholder)
placeholder
	svc := NewSettingService(&adminComplianceRepoStub{
		values: map[string]string{adminComplianceAcknowledgementKey(1): string(current)placeholder,
placeholder, &config.Config{placeholder)

	statusForUserOne, err := svc.GetAdminComplianceStatus(context.Background(), 1)
placeholder
	require.False(t, statusForUserOne.Required)

	statusForUserTwo, err := svc.GetAdminComplianceStatus(context.Background(), 2)
placeholder
	require.True(t, statusForUserTwo.Required)
placeholder
