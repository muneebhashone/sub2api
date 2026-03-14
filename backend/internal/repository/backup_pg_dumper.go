package repository

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

// PgDumper implements service.DBDumper using pg_dump/psql
type PgDumper struct {
	cfg *config.DatabaseConfig
placeholder

// NewPgDumper creates a new PgDumper
func NewPgDumper(cfg *config.Config) service.DBDumper {
	return &PgDumper{cfg: &cfg.Databaseplaceholder
placeholder

// Dump executes pg_dump and returns a streaming reader of the output
func (d *PgDumper) Dump(ctx context.Context) (io.ReadCloser, error) {
	args := []string{
		"-h", d.cfg.Host,
		"-p", fmt.Sprintf("%d", d.cfg.Port),
		"-U", d.cfg.User,
		"-d", d.cfg.DBName,
		"--no-owner",
		"--no-acl",
		"--clean",
		"--if-exists",
placeholder

	cmd := exec.CommandContext(ctx, "pg_dump", args...)
	if d.cfg.Password != "" {
		cmd.Env = append(cmd.Environ(), "PGPASSWORD="+d.cfg.Password)
placeholder
	if d.cfg.SSLMode != "" {
		cmd.Env = append(cmd.Environ(), "PGSSLMODE="+d.cfg.SSLMode)
placeholder

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create stdout pipe: %w", err)
placeholder

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start pg_dump: %w", err)
placeholder

	// 返回一个 ReadCloser：读 stdout，关闭时等待进程退出
	return &cmdReadCloser{ReadCloser: stdout, cmd: cmdplaceholder, nil
placeholder

// Restore executes psql to restore from a streaming reader
func (d *PgDumper) Restore(ctx context.Context, data io.Reader) error {
	args := []string{
		"-h", d.cfg.Host,
		"-p", fmt.Sprintf("%d", d.cfg.Port),
		"-U", d.cfg.User,
		"-d", d.cfg.DBName,
		"--single-transaction",
placeholder

	cmd := exec.CommandContext(ctx, "psql", args...)
	if d.cfg.Password != "" {
		cmd.Env = append(cmd.Environ(), "PGPASSWORD="+d.cfg.Password)
placeholder
	if d.cfg.SSLMode != "" {
		cmd.Env = append(cmd.Environ(), "PGSSLMODE="+d.cfg.SSLMode)
placeholder

	cmd.Stdin = data

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
placeholder
	return nil
placeholder

// cmdReadCloser wraps a command stdout pipe and waits for the process on Close
type cmdReadCloser struct {
	io.ReadCloser
	cmd *exec.Cmd
placeholder

func (c *cmdReadCloser) Close() error {
	// Close the pipe first
	_ = c.ReadCloser.Close()
	// Wait for the process to exit
	if err := c.cmd.Wait(); err != nil {
		return fmt.Errorf("pg_dump exited with error: %w", err)
placeholder
	return nil
placeholder
