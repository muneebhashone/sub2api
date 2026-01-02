package service

import (
	"context"
	"time"
)

// ErrorLog represents an ops error log item for list queries.
//
// Field naming matches docs/API-运维监控中心2.0.md (L3 根因追踪 - 错误日志列表).
type ErrorLog struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`

	Level        string `json:"level,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	AccountID    string `json:"account_id,omitempty"`
	APIPath      string `json:"api_path,omitempty"`
	Provider     string `json:"provider,omitempty"`
	Model        string `json:"model,omitempty"`
	HTTPCode     int    `json:"http_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`

	DurationMs *int `json:"duration_ms,omitempty"`
	RetryCount *int `json:"retry_count,omitempty"`
	Stream     bool `json:"stream,omitempty"`
placeholder

// ErrorLogFilter describes optional filters and pagination for listing ops error logs.
type ErrorLogFilter struct {
	StartTime *time.Time
	EndTime   *time.Time

	ErrorCode *int
	Provider  string
	AccountID *int64

	Page     int
	PageSize int
placeholder

func (f *ErrorLogFilter) normalize() (page, pageSize int) {
	page = 1
	pageSize = 20
	if f == nil {
		return page, pageSize
placeholder

	if f.Page > 0 {
		page = f.Page
placeholder
	if f.PageSize > 0 {
		pageSize = f.PageSize
placeholder
	if pageSize > 100 {
		pageSize = 100
placeholder
	return page, pageSize
placeholder

type ErrorLogListResponse struct {
	Errors   []*ErrorLog `json:"errors"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
placeholder

func (s *OpsService) GetErrorLogs(ctx context.Context, filter *ErrorLogFilter) (*ErrorLogListResponse, error) {
	if s == nil || s.repo == nil {
		return &ErrorLogListResponse{
			Errors:   []*ErrorLog{placeholder,
			Total:    0,
			Page:     1,
			PageSize: 20,
	placeholder, nil
placeholder

	page, pageSize := filter.normalize()
	if filter == nil {
		filter = &ErrorLogFilter{placeholder
placeholder
	filter.Page = page
	filter.PageSize = pageSize

	items, total, err := s.repo.ListErrorLogs(ctx, filter)
	if err != nil {
		return nil, err
placeholder
	if items == nil {
		items = []*ErrorLog{placeholder
placeholder

	return &ErrorLogListResponse{
		Errors:   items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
placeholder, nil
placeholder
