package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func (s *OpsService) ListAlertRules(ctx context.Context) ([]*OpsAlertRule, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return []*OpsAlertRule{placeholder, nil
placeholder
	return s.opsRepo.ListAlertRules(ctx)
placeholder

func (s *OpsService) CreateAlertRule(ctx context.Context, rule *OpsAlertRule) (*OpsAlertRule, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if rule == nil {
		return nil, infraerrors.BadRequest("INVALID_RULE", "invalid rule")
placeholder

	created, err := s.opsRepo.CreateAlertRule(ctx, rule)
	if err != nil {
		return nil, err
placeholder
	return created, nil
placeholder

func (s *OpsService) UpdateAlertRule(ctx context.Context, rule *OpsAlertRule) (*OpsAlertRule, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if rule == nil || rule.ID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_RULE", "invalid rule")
placeholder

	updated, err := s.opsRepo.UpdateAlertRule(ctx, rule)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, infraerrors.NotFound("OPS_ALERT_RULE_NOT_FOUND", "alert rule not found")
	placeholder
		return nil, err
placeholder
	return updated, nil
placeholder

func (s *OpsService) DeleteAlertRule(ctx context.Context, id int64) error {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return err
placeholder
	if s.opsRepo == nil {
		return infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if id <= 0 {
		return infraerrors.BadRequest("INVALID_RULE_ID", "invalid rule id")
placeholder
	if err := s.opsRepo.DeleteAlertRule(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return infraerrors.NotFound("OPS_ALERT_RULE_NOT_FOUND", "alert rule not found")
	placeholder
		return err
placeholder
	return nil
placeholder

func (s *OpsService) ListAlertEvents(ctx context.Context, filter *OpsAlertEventFilter) ([]*OpsAlertEvent, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return []*OpsAlertEvent{placeholder, nil
placeholder
	return s.opsRepo.ListAlertEvents(ctx, filter)
placeholder

func (s *OpsService) GetActiveAlertEvent(ctx context.Context, ruleID int64) (*OpsAlertEvent, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if ruleID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_RULE_ID", "invalid rule id")
placeholder
	return s.opsRepo.GetActiveAlertEvent(ctx, ruleID)
placeholder

func (s *OpsService) GetLatestAlertEvent(ctx context.Context, ruleID int64) (*OpsAlertEvent, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if ruleID <= 0 {
		return nil, infraerrors.BadRequest("INVALID_RULE_ID", "invalid rule id")
placeholder
	return s.opsRepo.GetLatestAlertEvent(ctx, ruleID)
placeholder

func (s *OpsService) CreateAlertEvent(ctx context.Context, event *OpsAlertEvent) (*OpsAlertEvent, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
placeholder
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if event == nil {
		return nil, infraerrors.BadRequest("INVALID_EVENT", "invalid event")
placeholder

	created, err := s.opsRepo.CreateAlertEvent(ctx, event)
	if err != nil {
		return nil, err
placeholder
	return created, nil
placeholder

func (s *OpsService) UpdateAlertEventStatus(ctx context.Context, eventID int64, status string, resolvedAt *time.Time) error {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return err
placeholder
	if s.opsRepo == nil {
		return infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if eventID <= 0 {
		return infraerrors.BadRequest("INVALID_EVENT_ID", "invalid event id")
placeholder
	if strings.TrimSpace(status) == "" {
		return infraerrors.BadRequest("INVALID_STATUS", "invalid status")
placeholder
	return s.opsRepo.UpdateAlertEventStatus(ctx, eventID, status, resolvedAt)
placeholder

func (s *OpsService) UpdateAlertEventEmailSent(ctx context.Context, eventID int64, emailSent bool) error {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return err
placeholder
	if s.opsRepo == nil {
		return infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
placeholder
	if eventID <= 0 {
		return infraerrors.BadRequest("INVALID_EVENT_ID", "invalid event id")
placeholder
	return s.opsRepo.UpdateAlertEventEmailSent(ctx, eventID, emailSent)
placeholder
