package gatewayhook

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// stubHook 是可编程的测试钩子，记录自己是否被执行。
type stubHook struct {
	id       string
	decision *Decision
	err      error
	panicVal any
	calls    *[]string
placeholder

func (h *stubHook) HookID() string { return h.id placeholder

func (h *stubHook) CheckPreFlight(ctx context.Context, req *Request) (*Decision, error) {
	if h.calls != nil {
		*h.calls = append(*h.calls, h.id)
placeholder
	if h.panicVal != nil {
		panic(h.panicVal)
placeholder
	return h.decision, h.err
placeholder

func TestChainRun_ExecutesHooksInOrder(t *testing.T) {
	var calls []string
	chain := NewChain(zap.NewNop(),
		&stubHook{id: "first", calls: &callsplaceholder,
		&stubHook{id: "second", calls: &callsplaceholder,
		&stubHook{id: "third", calls: &callsplaceholder,
	)

	decision := chain.Run(context.Background(), &Request{placeholder)

	require.Nil(t, decision, "全部放行时 Run 必须返回 nil")
	require.Equal(t, []string{"first", "second", "third"placeholder, calls, "钩子必须按注册顺序执行")
placeholder

func TestChainRun_FirstBlockedShortCircuits(t *testing.T) {
	var calls []string
	blocked := &Decision{Blocked: true, StatusCode: 403, ErrorType: "content_policy_violation", Message: "blocked"placeholder
	chain := NewChain(zap.NewNop(),
		&stubHook{id: "pass", calls: &callsplaceholder,
		&stubHook{id: "block", decision: blocked, calls: &callsplaceholder,
		&stubHook{id: "never", calls: &callsplaceholder,
	)

	decision := chain.Run(context.Background(), &Request{placeholder)

	require.Same(t, blocked, decision, "必须原样返回首个 Blocked Decision")
	require.Equal(t, []string{"pass", "block"placeholder, calls, "Blocked 之后的钩子不得执行")
placeholder

func TestChainRun_NonBlockedDecisionContinues(t *testing.T) {
	var calls []string
	chain := NewChain(zap.NewNop(),
		&stubHook{id: "flag-only", decision: &Decision{Blocked: false, Message: "flagged"placeholder, calls: &callsplaceholder,
		&stubHook{id: "after", calls: &callsplaceholder,
	)

	decision := chain.Run(context.Background(), &Request{placeholder)

	require.Nil(t, decision, "非 Blocked 的 Decision 视为放行")
	require.Equal(t, []string{"flag-only", "after"placeholder, calls)
placeholder

func TestChainRun_PanicIsolatedAndFailOpen(t *testing.T) {
	core, logs := observer.New(zap.ErrorLevel)
	var calls []string
	blocked := &Decision{Blocked: true, StatusCode: 403, Message: "blocked"placeholder
	chain := NewChain(zap.New(core),
		&stubHook{id: "boomer", panicVal: "boom", calls: &callsplaceholder,
		&stubHook{id: "blocker", decision: blocked, calls: &callsplaceholder,
	)

	decision := chain.Run(context.Background(), &Request{placeholder)

	require.Same(t, blocked, decision, "panic 钩子 fail-open 后链必须继续执行后续钩子")
	require.Equal(t, []string{"boomer", "blocker"placeholder, calls)

	entries := logs.FilterMessage("gatewayhook.hook_panic").All()
	require.Len(t, entries, 1, "panic 必须记一条 error 日志")
	fields := entries[0].ContextMap()
	require.Equal(t, "boomer", fields["hook_id"])
	require.Equal(t, "boom", fields["panic_value"])

	// 仅 panic 钩子的链：整体放行。
	onlyPanic := NewChain(zap.New(core), &stubHook{id: "boomer2", panicVal: errors.New("kaboom")placeholder)
	require.Nil(t, onlyPanic.Run(context.Background(), &Request{placeholder))
placeholder

func TestChainRun_HookErrorFailOpen(t *testing.T) {
	core, logs := observer.New(zap.WarnLevel)
	var calls []string
	chain := NewChain(zap.New(core),
		// 即使 error 钩子同时返回了 Blocked Decision，error 优先按 fail-open 丢弃。
		&stubHook{id: "broken", decision: &Decision{Blocked: trueplaceholder, err: errors.New("hook exploded"), calls: &callsplaceholder,
		&stubHook{id: "after", calls: &callsplaceholder,
	)

	decision := chain.Run(context.Background(), &Request{placeholder)

	require.Nil(t, decision, "钩子 error 必须 fail-open 放行")
	require.Equal(t, []string{"broken", "after"placeholder, calls, "error 钩子之后链必须继续")

	entries := logs.FilterMessage("gatewayhook.hook_error").All()
	require.Len(t, entries, 1)
	require.Equal(t, "broken", entries[0].ContextMap()["hook_id"])
placeholder

func TestChainRun_EmptyChain(t *testing.T) {
	require.Nil(t, NewChain(zap.NewNop()).Run(context.Background(), &Request{placeholder))
	require.True(t, NewChain(zap.NewNop()).IsEmpty())

	var nilChain *Chain
	require.True(t, nilChain.IsEmpty(), "nil 链视为空链")
	require.Nil(t, nilChain.Run(context.Background(), &Request{placeholder), "nil 链 Run 必须安全放行")

	// nil 钩子在构造时被过滤。
	require.True(t, NewChain(zap.NewNop(), nil, nil).IsEmpty())
placeholder
