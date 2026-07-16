package securityaudit

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type scriptedScanner struct {
	mu      sync.Mutex
	calls   []string
	block   <-chan struct{placeholder
	entered chan<- struct{placeholder
placeholder

func (s *scriptedScanner) Scan(ctx context.Context, endpoint ActiveEndpoint, _ string, _ []string) (*NormalizedResult, error) {
	s.mu.Lock()
	s.calls = append(s.calls, endpoint.ID)
	s.mu.Unlock()
	if s.entered != nil {
		select {
		case s.entered <- struct{placeholder{placeholder:
		default:
	placeholder
placeholder
	if s.block != nil {
		select {
		case <-s.block:
		case <-ctx.Done():
			return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Timeout: true, Cause: ctx.Err()placeholder
	placeholder
placeholder
	if endpoint.ID == "bad" {
		return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: trueplaceholder
placeholder
	if endpoint.ID == "invalid" {
		return nil, &GuardError{Code: ErrorCodeInvalidResponseplaceholder
placeholder
	return &NormalizedResult{Decision: EventPass, RiskLevel: RiskLow, Action: ActionAllow, Safety: "Safe", ScannerScores: map[string]float64{placeholder, ScannerEvidence: map[string]string{placeholder, GuardEndpointID: endpoint.IDplaceholder, nil
placeholder

func guardConfig(endpoints ...ActiveEndpoint) ActiveConfig {
	return ActiveConfig{RiskControlEnabled: true, Enabled: true, BlockingEnabled: true, ConfigVersion: 2, Scanners: AllScannerIDs, Endpoints: endpointsplaceholder
placeholder

func TestGuardEvaluatorOrderedFailoverAndInvalidTerminal(t *testing.T) {
	scanner := &scriptedScanner{placeholder
	metrics := NewAtomicMetrics()
	evaluator := newGuardEvaluator(scanner, nil, metrics, 4, 2)
	snapshot := PromptSnapshot{RequestID: "r", ScanText: "hello", PromptLength: 5placeholder
	decision, err := evaluator.Evaluate(context.Background(), guardConfig(
		ActiveEndpoint{ID: "bad", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder,
		ActiveEndpoint{ID: "good", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder,
	), snapshot)
placeholder
	require.Equal(t, DecisionAllow, decision.Kind)
	require.Equal(t, int64(1), metrics.Snapshot().Failovers)
	_, err = evaluator.Evaluate(context.Background(), guardConfig(
		ActiveEndpoint{ID: "invalid", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder,
		ActiveEndpoint{ID: "good", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder,
	), snapshot)
	var guardErr *GuardError
	require.ErrorAs(t, err, &guardErr)
	require.Equal(t, ErrorCodeInvalidResponse, guardErr.Code)
	snapshotMetrics := metrics.Snapshot()
	require.Equal(t, int64(2), snapshotMetrics.Total)
	require.Equal(t, int64(1), snapshotMetrics.Allowed)
	require.Equal(t, int64(1), snapshotMetrics.Invalid)
placeholder

func TestGuardEvaluatorGlobalBulkheadIsNonBlocking(t *testing.T) {
	release := make(chan struct{placeholder)
	entered := make(chan struct{placeholder, 1)
	scanner := &scriptedScanner{block: release, entered: enteredplaceholder
	metrics := NewAtomicMetrics()
	evaluator := newGuardEvaluator(scanner, nil, metrics, 1, 1)
	cfg := guardConfig(ActiveEndpoint{ID: "good", Enabled: true, TimeoutMS: 2000, InputLimit: 100placeholder)
	done := make(chan error, 1)
	go func() {
		_, err := evaluator.Evaluate(context.Background(), cfg, PromptSnapshot{ScanText: "one", PromptLength: 3placeholder)
		done <- err
placeholder()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("first evaluation did not enter scanner")
placeholder
	start := time.Now()
	_, err := evaluator.Evaluate(context.Background(), cfg, PromptSnapshot{ScanText: "two", PromptLength: 3placeholder)
placeholder
	require.Less(t, time.Since(start), 200*time.Millisecond)
	require.Equal(t, int64(1), metrics.Snapshot().BulkheadFull)
	close(release)
	require.NoError(t, <-done)
	snapshotMetrics := metrics.Snapshot()
	require.Equal(t, int64(2), snapshotMetrics.Total)
	require.Equal(t, int64(1), snapshotMetrics.Allowed)
	require.Equal(t, int64(1), snapshotMetrics.Unavailable)
placeholder

func TestGuardEvaluatorPerNodeBulkheadIsNonBlocking(t *testing.T) {
	release := make(chan struct{placeholder)
	entered := make(chan struct{placeholder, 1)
	scanner := &scriptedScanner{block: release, entered: enteredplaceholder
	metrics := NewAtomicMetrics()
	evaluator := newGuardEvaluator(scanner, nil, metrics, 2, 1)
	cfg := guardConfig(ActiveEndpoint{ID: "same-node", Enabled: true, TimeoutMS: 2000, InputLimit: 100placeholder)
	done := make(chan error, 1)
	go func() {
		_, err := evaluator.Evaluate(context.Background(), cfg, PromptSnapshot{ScanText: "one", PromptLength: 3placeholder)
		done <- err
placeholder()
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("first evaluation did not enter scanner")
placeholder
	started := time.Now()
	_, err := evaluator.Evaluate(context.Background(), cfg, PromptSnapshot{ScanText: "two", PromptLength: 3placeholder)
placeholder
	require.Less(t, time.Since(started), 200*time.Millisecond)
	require.GreaterOrEqual(t, metrics.Snapshot().BulkheadFull, int64(1))
	close(release)
	require.NoError(t, <-done)
placeholder

func TestGuardEvaluatorLastChunkFailureNeverAllows(t *testing.T) {
	call := 0
	scanner := PromptScannerFunc(func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error) {
		call++
		if call == 2 {
			return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Cause: errors.New("down")placeholder
	placeholder
		return &NormalizedResult{Decision: EventPass, RiskLevel: RiskLow, Action: ActionAllow, ScannerScores: map[string]float64{placeholder, ScannerEvidence: map[string]string{placeholderplaceholder, nil
placeholder)
	metrics := NewAtomicMetrics()
	evaluator := newGuardEvaluator(scanner, nil, metrics, 2, 2)
	_, err := evaluator.Evaluate(context.Background(), guardConfig(ActiveEndpoint{ID: "one", Enabled: true, TimeoutMS: 1000, InputLimit: 3placeholder), PromptSnapshot{ScanText: "abcdef", PromptLength: 6placeholder)
placeholder
placeholder

func TestGuardEvaluatorBlockStopsRemainingChunksButReportsPlannedTotal(t *testing.T) {
	calls := 0
	scanner := PromptScannerFunc(func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error) {
		calls++
		return &NormalizedResult{
			Decision: EventCritical, RiskLevel: RiskCritical, Action: ActionBlock, Safety: "Unsafe",
			Categories: []string{"jailbreak"placeholder, MatchedScanners: []string{"jailbreak"placeholder,
			ScannerScores: map[string]float64{"jailbreak": 1placeholder, ScannerEvidence: map[string]string{"jailbreak": "Jailbreak"placeholder,
	placeholder, nil
placeholder)
	metrics := NewAtomicMetrics()
	evaluator := newGuardEvaluator(scanner, nil, metrics, 2, 2)
	decision, err := evaluator.Evaluate(context.Background(), guardConfig(
		ActiveEndpoint{ID: "one", Enabled: true, TimeoutMS: 1000, InputLimit: 3placeholder,
	), PromptSnapshot{ScanText: "abcdefghi", PromptLength: 9placeholder)
placeholder
	require.Equal(t, DecisionBlock, decision.Kind)
	require.Equal(t, 1, calls)
	require.Equal(t, 3, decision.Result.ChunkTotal)
	require.Equal(t, int64(1), metrics.Snapshot().Blocked)
placeholder

func TestGuardEvaluatorFlagSharedDeadlineFailClosedAndContextCancel(t *testing.T) {
	t.Run("flag allows next stage", func(t *testing.T) {
		metrics := NewAtomicMetrics()
		evaluator := newGuardEvaluator(PromptScannerFunc(func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error) {
			return &NormalizedResult{Decision: EventFlag, RiskLevel: RiskMedium, Action: ActionWarn, Safety: "Controversial", Categories: []string{"violent"placeholder, MatchedScanners: []string{"violent"placeholder, ScannerScores: map[string]float64{"violent": .5placeholder, ScannerEvidence: map[string]string{"violent": "Violent"placeholderplaceholder, nil
	placeholder), nil, metrics, 2, 2)
		decision, err := evaluator.Evaluate(context.Background(), guardConfig(ActiveEndpoint{ID: "one", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder), PromptSnapshot{ScanText: "review", PromptLength: 6placeholder)
	placeholder
		require.Equal(t, DecisionFlag, decision.Kind)
		require.True(t, decision.AllowNextStage)
		require.Equal(t, int64(1), metrics.Snapshot().Flagged)
placeholder)

	t.Run("all failovers share first endpoint deadline", func(t *testing.T) {
		calls := 0
		scanner := PromptScannerFunc(func(ctx context.Context, endpoint ActiveEndpoint, _ string, _ []string) (*NormalizedResult, error) {
			calls++
			if endpoint.ID == "first" {
				select {
				case <-time.After(35 * time.Millisecond):
					return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: trueplaceholder
				case <-ctx.Done():
					return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Timeout: true, Cause: ctx.Err()placeholder
			placeholder
		placeholder
			<-ctx.Done()
			return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Timeout: true, Cause: ctx.Err()placeholder
	placeholder)
		metrics := NewAtomicMetrics()
		evaluator := newGuardEvaluator(scanner, nil, metrics, 2, 2)
		started := time.Now()
		_, err := evaluator.Evaluate(context.Background(), guardConfig(
			ActiveEndpoint{ID: "first", Enabled: true, TimeoutMS: 70, InputLimit: 100placeholder,
			ActiveEndpoint{ID: "second", Enabled: true, TimeoutMS: 500, InputLimit: 100placeholder,
		), PromptSnapshot{ScanText: "deadline", PromptLength: 8placeholder)
		elapsed := time.Since(started)
	placeholder
		require.Equal(t, 2, calls)
		require.Less(t, elapsed, 180*time.Millisecond)
		require.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
		require.Equal(t, int64(1), metrics.Snapshot().Failovers)
		require.Equal(t, int64(1), metrics.Snapshot().Timeouts)
placeholder)

	t.Run("canceled parent never allows", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		evaluator := newGuardEvaluator(PromptScannerFunc(func(ctx context.Context, _ ActiveEndpoint, _ string, _ []string) (*NormalizedResult, error) {
			<-ctx.Done()
			return nil, &GuardError{Code: ErrorCodeUnavailable, Retryable: true, Cause: ctx.Err()placeholder
	placeholder), nil, NewAtomicMetrics(), 2, 2)
		decision, err := evaluator.Evaluate(ctx, guardConfig(ActiveEndpoint{ID: "one", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder), PromptSnapshot{ScanText: "cancel", PromptLength: 6placeholder)
	placeholder
		require.Nil(t, decision)
placeholder)
placeholder

func TestGuardEvaluatorRecordsExistingResultOnceAndRecordFailureDoesNotChangeDecision(t *testing.T) {
	for _, recordErr := range []error{nil, errors.New("database unavailable")placeholder {
		repo := &fakeJobRepository{recordBlockingErr: recordErrplaceholder
		metrics := NewAtomicMetrics()
		scannerCalls := 0
		evaluator := newGuardEvaluator(PromptScannerFunc(func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error) {
			scannerCalls++
			return &NormalizedResult{Decision: EventCritical, RiskLevel: RiskCritical, Action: ActionBlock, Safety: "Unsafe", Categories: []string{"pii"placeholder, MatchedScanners: []string{"pii"placeholder, ScannerScores: map[string]float64{"pii": 1placeholder, ScannerEvidence: map[string]string{"pii": "PII"placeholderplaceholder, nil
	placeholder), repo, metrics, 2, 2)
		decision, err := evaluator.Evaluate(context.Background(), guardConfig(ActiveEndpoint{ID: "one", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder), PromptSnapshot{ScanText: "raw prompt", RedactedPreview: "raw***", PromptLength: 10placeholder)
	placeholder
		require.Equal(t, DecisionBlock, decision.Kind)
		require.Equal(t, 1, scannerCalls)
		require.Equal(t, 1, repo.recordBlockingCalls)
		require.Empty(t, repo.recordBlockingSnapshot.ScanText)
		require.Same(t, decision.Result, repo.recordBlockingResult)
		if recordErr != nil {
			require.Equal(t, int64(1), metrics.Snapshot().RecordFailed)
	placeholder else {
			require.Zero(t, metrics.Snapshot().RecordFailed)
	placeholder
placeholder
placeholder

func TestGuardEvaluatorNilResultAndScannerPanicBecomeStableFailures(t *testing.T) {
	tests := []struct {
		name string
		scan PromptScannerFunc
		code string
placeholder{
		{name: "nil result", scan: func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error) { return nil, nil placeholder, code: ErrorCodeInvalidResponseplaceholder,
		{name: "panic", scan: func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error) {
			panic("raw prompt canary")
	placeholder, code: ErrorCodeUnavailableplaceholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator := newGuardEvaluator(tt.scan, nil, NewAtomicMetrics(), 2, 2)
			_, err := evaluator.Evaluate(context.Background(), guardConfig(ActiveEndpoint{ID: "one", Enabled: true, TimeoutMS: 1000, InputLimit: 100placeholder), PromptSnapshot{ScanText: "input", PromptLength: 5placeholder)
			var guardErr *GuardError
			require.ErrorAs(t, err, &guardErr)
			require.Equal(t, tt.code, guardErr.Code)
			require.NotContains(t, err.Error(), "canary")
	placeholder)
placeholder
placeholder

type PromptScannerFunc func(context.Context, ActiveEndpoint, string, []string) (*NormalizedResult, error)

func (f PromptScannerFunc) Scan(ctx context.Context, endpoint ActiveEndpoint, chunk string, scanners []string) (*NormalizedResult, error) {
	return f(ctx, endpoint, chunk, scanners)
placeholder
