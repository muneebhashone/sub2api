package repository

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/servertiming"
	"github.com/redis/go-redis/v9"
)

func TestServerTimingRedisHookRecordsCommands(t *testing.T) {
	collector := servertiming.New(time.Now())
	ctx := servertiming.WithCollector(context.Background(), collector)
	hook := serverTimingRedisHook{placeholder

	process := hook.ProcessHook(func(context.Context, redis.Cmder) error {
		time.Sleep(time.Millisecond)
		return errors.New("redis failure")
placeholder)
	if err := process(ctx, redis.NewStringCmd(ctx, "get", "sensitive-key")); err == nil {
		t.Fatal("ProcessHook did not return the underlying error")
placeholder

	pipeline := hook.ProcessPipelineHook(func(context.Context, []redis.Cmder) error {
		time.Sleep(time.Millisecond)
		return nil
placeholder)
	commands := []redis.Cmder{
		redis.NewStringCmd(ctx, "get", "first-secret"),
		redis.NewStringCmd(ctx, "get", "second-secret"),
		redis.NewStatusCmd(ctx, "set", "third-secret", "value"),
placeholder
	if err := pipeline(ctx, commands); err != nil {
		t.Fatal(err)
placeholder

	header := collector.HeaderValue(time.Now(), "bypass")
	if !strings.Contains(header, `commands=4`) {
		t.Fatalf("header %q does not report one command and a three-command pipeline", header)
placeholder
	if strings.Contains(header, "secret") || strings.Contains(header, "get") {
		t.Fatalf("Redis command details leaked into header: %q", header)
placeholder
placeholder

func TestServerTimingRedisHookSkipsInactiveContext(t *testing.T) {
	called := false
	hook := serverTimingRedisHook{placeholder
	process := hook.ProcessHook(func(context.Context, redis.Cmder) error {
		called = true
		return nil
placeholder)
	ctx := context.Background()
	if err := process(ctx, redis.NewStringCmd(ctx, "ping")); err != nil {
		t.Fatal(err)
placeholder
	if !called {
		t.Fatal("inactive Redis command did not reach the next hook")
placeholder
placeholder
