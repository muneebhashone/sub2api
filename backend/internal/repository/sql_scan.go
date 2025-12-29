package repository

import (
	"context"
	"database/sql"
)

type sqlQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
placeholder

// scanSingleRow executes a query and scans the first row into dest.
// If no rows are returned, sql.ErrNoRows is returned.
// 设计目的：仅依赖 QueryContext，避免 QueryRowContext 对 *sql.Tx 的强绑定，
// 让 ent.Tx 也能作为 sqlExecutor/Queryer 使用。
func scanSingleRow(ctx context.Context, q sqlQueryer, query string, args []any, dest ...any) error {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return err
placeholder
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
	placeholder
		return sql.ErrNoRows
placeholder
	if err := rows.Scan(dest...); err != nil {
		return err
placeholder
	return rows.Err()
placeholder
