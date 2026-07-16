package securityaudit

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeLegacyEngine struct {
	decision *LegacyDecision
	err      error
	calls    atomic.Int64
placeholder

func (f *fakeLegacyEngine) Check(context.Context, Request) (*LegacyDecision, error) {
	f.calls.Add(1)
	return f.decision, f.err
placeholder

type fakePromptEngine struct {
	mode      Mode
	decision  *PromptDecision
	err       error
	enqueues  atomic.Int64
	evaluates atomic.Int64
placeholder

func (f *fakePromptEngine) EffectiveMode() Mode { return f.mode placeholder
func (f *fakePromptEngine) Enqueue(context.Context, Request) error {
	f.enqueues.Add(1)
	return f.err
placeholder
func (f *fakePromptEngine) Evaluate(context.Context, Request) (*PromptDecision, error) {
	f.evaluates.Add(1)
	return f.decision, f.err
placeholder

func TestCoordinatorModesAndPriority(t *testing.T) {
	tests := []struct {
		name           string
		mode           Mode
		legacy         *LegacyDecision
		prompt         *PromptDecision
		promptErr      error
		wantKind       DecisionKind
		wantCode       string
		wantEnqueue    int64
		wantEvaluation int64
placeholder{
		{name: "off", mode: ModeOff, wantKind: DecisionAllowplaceholder,
		{name: "async only enqueues", mode: ModeAsync, wantKind: DecisionAllow, wantEnqueue: 1placeholder,
		{name: "prompt block", mode: ModeBlocking, prompt: &PromptDecision{Kind: DecisionBlockplaceholder, wantKind: DecisionBlock, wantCode: ErrorCodeBlocked, wantEvaluation: 1placeholder,
		{name: "prompt unavailable", mode: ModeBlocking, promptErr: errors.New("down"), wantKind: DecisionUnavailable, wantCode: ErrorCodeUnavailable, wantEvaluation: 1placeholder,
		{name: "legacy wins both block", mode: ModeBlocking,
			legacy: &LegacyDecision{Blocked: true, StatusCode: http.StatusForbidden, ErrorCode: "content_policy_violation", Message: "legacy"placeholder,
			prompt: &PromptDecision{Kind: DecisionBlockplaceholder, wantKind: DecisionBlock, wantCode: "content_policy_violation", wantEvaluation: 1placeholder,
placeholder
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			legacy := &fakeLegacyEngine{decision: tt.legacyplaceholder
			prompt := &fakePromptEngine{mode: tt.mode, decision: tt.prompt, err: tt.promptErrplaceholder
			decision := NewCoordinator(legacy, prompt).Check(context.Background(), Request{Body: []byte(`{placeholder`)placeholder)
			require.Equal(t, tt.wantKind, decision.Kind)
			require.Equal(t, tt.wantCode, decision.ErrorCode)
			require.Equal(t, int64(1), legacy.calls.Load())
			require.Equal(t, tt.wantEnqueue, prompt.enqueues.Load())
			require.Equal(t, tt.wantEvaluation, prompt.evaluates.Load())
	placeholder)
placeholder
placeholder

func TestCoordinatorDoesNotMutateRequestBody(t *testing.T) {
	body := []byte(`{"messages":[{"role":"user","content":"hello"placeholder]placeholder`)
	original := append([]byte(nil), body...)
	prompt := &fakePromptEngine{mode: ModeAsyncplaceholder
	decision := NewCoordinator(&fakeLegacyEngine{placeholder, prompt).Check(context.Background(), Request{Body: bodyplaceholder)
	require.True(t, decision.AllowNextStage)
	require.Equal(t, original, body)
placeholder

func TestCoordinatorBlockingPriorityCoversBothEngineDecisionMatrix(t *testing.T) {
	legacyCases := []struct {
		name     string
		decision *LegacyDecision
placeholder{
		{name: "allow", decision: &LegacyDecision{Allowed: true, StatusCode: http.StatusOK, Action: "allow"placeholderplaceholder,
		{name: "flag", decision: &LegacyDecision{Allowed: true, Flagged: true, StatusCode: http.StatusOK, Action: "flag"placeholderplaceholder,
		{name: "block", decision: &LegacyDecision{Blocked: true, StatusCode: http.StatusForbidden, ErrorCode: "legacy_exact_code", Message: "legacy exact message", Action: "block"placeholderplaceholder,
placeholder
	promptCases := []struct {
		name     string
		decision *PromptDecision
		wantKind DecisionKind
		wantCode string
placeholder{
		{name: "allow", decision: &PromptDecision{Kind: DecisionAllow, AllowNextStage: trueplaceholder, wantKind: DecisionAllowplaceholder,
		{name: "flag", decision: &PromptDecision{Kind: DecisionFlag, AllowNextStage: trueplaceholder, wantKind: DecisionFlagplaceholder,
		{name: "block", decision: &PromptDecision{Kind: DecisionBlockplaceholder, wantKind: DecisionBlock, wantCode: ErrorCodeBlockedplaceholder,
		{name: "unavailable", decision: &PromptDecision{Kind: DecisionUnavailable, ErrorCode: ErrorCodeUnavailableplaceholder, wantKind: DecisionUnavailable, wantCode: ErrorCodeUnavailableplaceholder,
		{name: "invalid", decision: &PromptDecision{Kind: DecisionInvalid, ErrorCode: ErrorCodeInvalidResponseplaceholder, wantKind: DecisionInvalid, wantCode: ErrorCodeInvalidResponseplaceholder,
placeholder

	for _, legacyCase := range legacyCases {
		for _, promptCase := range promptCases {
			t.Run(fmt.Sprintf("legacy_%s_prompt_%s", legacyCase.name, promptCase.name), func(t *testing.T) {
				legacy := &fakeLegacyEngine{decision: legacyCase.decisionplaceholder
				prompt := &fakePromptEngine{mode: ModeBlocking, decision: promptCase.decisionplaceholder
				decision := NewCoordinator(legacy, prompt).Check(context.Background(), Request{placeholder)

				require.Same(t, legacyCase.decision, decision.Legacy)
				require.Same(t, promptCase.decision, decision.Prompt)
				require.Equal(t, int64(1), legacy.calls.Load())
				require.Equal(t, int64(1), prompt.evaluates.Load())
				if legacyCase.name == "block" {
					require.Equal(t, DecisionBlock, decision.Kind)
					require.Equal(t, "legacy_exact_code", decision.ErrorCode)
					require.Equal(t, "legacy exact message", decision.ClientMessage)
					require.False(t, decision.AllowNextStage)
					return
			placeholder
				require.Equal(t, promptCase.wantKind, decision.Kind)
				require.Equal(t, promptCase.wantCode, decision.ErrorCode)
				require.Equal(t, promptCase.decision.AllowNextStage, decision.AllowNextStage)
		placeholder)
	placeholder
placeholder
placeholder

func TestCoordinatorPreservesIndependentEngineFactsAndMapsOnlyGatewayOutcome(t *testing.T) {
	legacyDecision := &LegacyDecision{
		Allowed: true, Flagged: true, Message: "legacy finding", StatusCode: http.StatusAccepted,
		ErrorCode: "legacy_observation", Action: "legacy_action",
placeholder
	promptResult := &NormalizedResult{
		Decision: EventCritical, RiskLevel: RiskCritical, Action: ActionBlock,
		Categories: []string{"pii"placeholder, ScannerScores: map[string]float64{"pii": 1placeholder,
placeholder
	promptDecision := &PromptDecision{Kind: DecisionBlock, Result: promptResultplaceholder
	decision := NewCoordinator(
		&fakeLegacyEngine{decision: legacyDecisionplaceholder,
		&fakePromptEngine{mode: ModeBlocking, decision: promptDecisionplaceholder,
	).Check(context.Background(), Request{placeholder)

	require.Same(t, legacyDecision, decision.Legacy)
	require.Same(t, promptDecision, decision.Prompt)
	require.Same(t, promptResult, decision.Prompt.Result)
	require.Equal(t, "legacy finding", decision.Legacy.Message)
	require.Equal(t, []string{"pii"placeholder, decision.Prompt.Result.Categories)
	require.Equal(t, ErrorCodeBlocked, decision.ErrorCode)
placeholder

func TestCoordinatorAsyncEnqueueFailuresNeverChangeResponseOrDownstreamDispatch(t *testing.T) {
	for _, enqueueErr := range []error{ErrQueueFull, ErrQueueAdmissionBusy, errors.New("redis unavailable"), errors.New("publish failed")placeholder {
		prompt := &fakePromptEngine{mode: ModeAsync, err: enqueueErrplaceholder
		decision := NewCoordinator(&fakeLegacyEngine{decision: &LegacyDecision{Allowed: trueplaceholderplaceholder, prompt).Check(context.Background(), Request{placeholder)
		downstreamDispatches := 0
		status := http.StatusOK
		responseBody := "unchanged-upstream-response"
		if decision.AllowNextStage {
			downstreamDispatches++
	placeholder else {
			status = decision.HTTPStatus
			responseBody = decision.ClientMessage
	placeholder
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "unchanged-upstream-response", responseBody)
		require.Equal(t, 1, downstreamDispatches)
		require.Equal(t, int64(1), prompt.enqueues.Load())
		require.Zero(t, prompt.evaluates.Load())
placeholder
placeholder
