// Package pagination provides utilities for handling paginated queries and results.
package pagination

// PaginationParams 分页参数
type PaginationParams struct {
	Page     int
	PageSize int
placeholder

// PaginationResult 分页结果
type PaginationResult struct {
	Total    int64
	Page     int
	PageSize int
	Pages    int
placeholder

// DefaultPagination 默认分页参数
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 20,
placeholder
placeholder

// Offset 计算偏移量
func (p PaginationParams) Offset() int {
	if p.Page < 1 {
		p.Page = 1
placeholder
	return (p.Page - 1) * p.PageSize
placeholder

// Limit 获取限制数
func (p PaginationParams) Limit() int {
	if p.PageSize < 1 {
		return 20
placeholder
	if p.PageSize > 100 {
		return 100
placeholder
	return p.PageSize
placeholder
