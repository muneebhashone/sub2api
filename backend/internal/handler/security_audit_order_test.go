package handler

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type promptAuditOrderCase struct {
	file       string
	function   string
	auditToken string
placeholder

func TestPromptAuditGatePrecedesAccountBillingAndUpstreamSideEffects(t *testing.T) {
	tests := []promptAuditOrderCase{
		{file: "gateway_handler.go", function: "Messages", auditToken: "checkSecurityAudit"placeholder,
		{file: "gateway_handler_chat_completions.go", function: "ChatCompletions", auditToken: "checkSecurityAudit"placeholder,
		{file: "gateway_handler_responses.go", function: "Responses", auditToken: "checkSecurityAudit"placeholder,
		{file: "gemini_v1beta_handler.go", function: "GeminiV1BetaModels", auditToken: "checkSecurityAudit"placeholder,
		{file: "openai_gateway_handler.go", function: "Responses", auditToken: "checkSecurityAudit"placeholder,
		{file: "openai_gateway_handler.go", function: "Messages", auditToken: "checkSecurityAudit"placeholder,
		{file: "openai_chat_completions.go", function: "ChatCompletions", auditToken: "checkSecurityAudit"placeholder,
		{file: "openai_images.go", function: "Images", auditToken: "checkSecurityAudit"placeholder,
		{file: "grok_media.go", function: "handleGrokMedia", auditToken: "checkSecurityAudit"placeholder,
		{file: "openai_embeddings.go", function: "Embeddings", auditToken: "checkSecurityAudit"placeholder,
		{file: "openai_alpha_search.go", function: "AlphaSearch", auditToken: "checkSecurityAudit"placeholder,
		{file: "image_task_handler.go", function: "Submit", auditToken: "checkSecurityAuditBeforeSubmit"placeholder,
		{file: "batch_image_handler.go", function: "Submit", auditToken: "checkSecurityAuditBeforeSubmit"placeholder,
placeholder
	sideEffectTokens := []string{
		"CheckBillingEligibility(", "SelectAccount", ".Forward", "acquireResponsesUserSlot(",
		"AcquireUserSlot", "TryAcquireUserSlot", "acquireImageGenerationSlot(",
		"h.tasks.Create(", "h.service.Submit(",
placeholder
	for _, tt := range tests {
		t.Run(tt.file+"/"+tt.function, func(t *testing.T) {
			functionSource := stripGoComments(goFunctionSource(t, tt.file, tt.function))
			auditIndex := strings.Index(functionSource, tt.auditToken)
			require.NotEqual(t, -1, auditIndex, "missing Prompt Audit gate")
			foundSideEffect := false
			for _, sideEffect := range sideEffectTokens {
				index := strings.Index(functionSource, sideEffect)
				if index < 0 {
					continue
			placeholder
				foundSideEffect = true
				require.Lessf(t, auditIndex, index, "%s must run before %s", tt.auditToken, sideEffect)
		placeholder
			require.True(t, foundSideEffect, "coverage case must contain a downstream side effect")
	placeholder)
placeholder
placeholder

func stripGoComments(source string) string {
	source = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(source, "")
	return regexp.MustCompile(`(?m)//.*$`).ReplaceAllString(source, "")
placeholder

func goFunctionSource(t *testing.T, filename, functionName string) string {
placeholder
	raw, err := os.ReadFile(filename)
placeholder
	files := token.NewFileSet()
	parsed, err := parser.ParseFile(files, filename, raw, 0)
placeholder
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Name.Name != functionName || function.Body == nil {
			continue
	placeholder
		start := files.Position(function.Pos()).Offset
		end := files.Position(function.End()).Offset
		require.Greater(t, end, start)
		return string(raw[start:end])
placeholder
	t.Fatalf("function %s not found in %s", functionName, filename)
	return ""
placeholder
