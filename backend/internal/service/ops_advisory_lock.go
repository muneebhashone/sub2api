package service

import (
	"context"
	"database/sql"
	"hash/fnv"
	"time"
)

func hashAdvisoryLockID(key string) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(key))
	return int64(h.Sum64())
placeholder

func tryAcquireDBAdvisoryLock(ctx context.Context, db *sql.DB, lockID int64) (func(), bool) {
	if db == nil {
		return nil, false
placeholder
	if ctx == nil {
		ctx = context.Background()
placeholder

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, false
placeholder

	acquired := false
	if err := conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock($1)", lockID).Scan(&acquired); err != nil {
		_ = conn.Close()
		return nil, false
placeholder
	if !acquired {
		_ = conn.Close()
		return nil, false
placeholder

	release := func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, _ = conn.ExecContext(unlockCtx, "SELECT pg_advisory_unlock($1)", lockID)
		_ = conn.Close()
placeholder
	return release, true
placeholder
