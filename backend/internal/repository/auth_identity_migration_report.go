package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AuthIdentityMigrationReport struct {
	ID         int64
	ReportType string
	ReportKey  string
	Details    map[string]any
	CreatedAt  time.Time
placeholder

type AuthIdentityMigrationReportQuery struct {
	ReportType string
	Limit      int
	Offset     int
placeholder

type AuthIdentityMigrationReportSummary struct {
	Total  int64
	ByType map[string]int64
placeholder

func (r *userRepository) ListAuthIdentityMigrationReports(ctx context.Context, query AuthIdentityMigrationReportQuery) ([]AuthIdentityMigrationReport, error) {
	exec := txAwareSQLExecutor(ctx, r.sql, r.client)
	if exec == nil {
		return nil, fmt.Errorf("sql executor is not configured")
placeholder

	limit := query.Limit
	if limit <= 0 {
		limit = 100
placeholder
	rows, err := exec.QueryContext(ctx, `
SELECT id, report_type, report_key, details, created_at
FROM auth_identity_migration_reports
WHERE ($1 = '' OR report_type = $1)
ORDER BY created_at DESC, id DESC
LIMIT $2 OFFSET $3`,
		strings.TrimSpace(query.ReportType),
		limit,
		query.Offset,
	)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	reports := make([]AuthIdentityMigrationReport, 0)
	for rows.Next() {
		report, scanErr := scanAuthIdentityMigrationReport(rows)
		if scanErr != nil {
			return nil, scanErr
	placeholder
		reports = append(reports, report)
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return reports, nil
placeholder

func (r *userRepository) GetAuthIdentityMigrationReport(ctx context.Context, reportType, reportKey string) (*AuthIdentityMigrationReport, error) {
	exec := txAwareSQLExecutor(ctx, r.sql, r.client)
	if exec == nil {
		return nil, fmt.Errorf("sql executor is not configured")
placeholder

	rows, err := exec.QueryContext(ctx, `
SELECT id, report_type, report_key, details, created_at
FROM auth_identity_migration_reports
WHERE report_type = $1 AND report_key = $2
LIMIT 1`,
		strings.TrimSpace(reportType),
		strings.TrimSpace(reportKey),
	)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	if !rows.Next() {
		return nil, sql.ErrNoRows
placeholder
	report, err := scanAuthIdentityMigrationReport(rows)
	if err != nil {
		return nil, err
placeholder
	return &report, rows.Err()
placeholder

func (r *userRepository) SummarizeAuthIdentityMigrationReports(ctx context.Context) (*AuthIdentityMigrationReportSummary, error) {
	exec := txAwareSQLExecutor(ctx, r.sql, r.client)
	if exec == nil {
		return nil, fmt.Errorf("sql executor is not configured")
placeholder

	rows, err := exec.QueryContext(ctx, `
SELECT report_type, COUNT(*)
FROM auth_identity_migration_reports
GROUP BY report_type
ORDER BY report_type ASC`)
	if err != nil {
		return nil, err
placeholder
	defer func() { _ = rows.Close() placeholder()

	summary := &AuthIdentityMigrationReportSummary{
		ByType: make(map[string]int64),
placeholder
	for rows.Next() {
		var reportType string
		var count int64
		if err := rows.Scan(&reportType, &count); err != nil {
			return nil, err
	placeholder
		summary.ByType[reportType] = count
		summary.Total += count
placeholder
	if err := rows.Err(); err != nil {
		return nil, err
placeholder
	return summary, nil
placeholder

func scanAuthIdentityMigrationReport(scanner interface{ Scan(dest ...any) error placeholder) (AuthIdentityMigrationReport, error) {
	var (
		report  AuthIdentityMigrationReport
		details []byte
	)
	if err := scanner.Scan(&report.ID, &report.ReportType, &report.ReportKey, &details, &report.CreatedAt); err != nil {
		return AuthIdentityMigrationReport{placeholder, err
placeholder
	report.Details = map[string]any{placeholder
	if len(details) > 0 {
		if err := json.Unmarshal(details, &report.Details); err != nil {
			return AuthIdentityMigrationReport{placeholder, err
	placeholder
placeholder
	return report, nil
placeholder
